package controllers

import (
	"context"
	"fmt"
	"path"
	"strings"

	envoyv1alpha1 "github.com/jpeach/envoy-controller/api/v1alpha1"
	"github.com/jpeach/envoy-controller/pkg/kubernetes"
	"github.com/jpeach/envoy-controller/pkg/must"
	"github.com/jpeach/envoy-controller/pkg/xds"

	"github.com/go-logr/logr"
	"google.golang.org/protobuf/proto"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var factories = []func() runtime.Object{
	func() runtime.Object { return &envoyv1alpha1.Listener{} },
	func() runtime.Object { return &envoyv1alpha1.Cluster{} },
	func() runtime.Object { return &envoyv1alpha1.RouteConfiguration{} },
	func() runtime.Object { return &envoyv1alpha1.ScopedRouteConfiguration{} },
	func() runtime.Object { return &envoyv1alpha1.Secret{} },
	func() runtime.Object { return &envoyv1alpha1.Runtime{} },
	func() runtime.Object { return &envoyv1alpha1.VirtualHost{} },
}

func resourceOf(name types.NamespacedName, gvk schema.GroupVersionKind) xds.ResourceName {
	return xds.ResourceName(
		strings.ToLower(path.Join(name.Namespace, gvk.Kind, name.Name)),
	)
}

func versionOf(obj runtime.Object) xds.ResourceVersion {
	metaObj := must.Object(meta.Accessor(obj))

	return xds.ResourceVersion{
		Identifier: string(metaObj.GetUID()),
		Version:    metaObj.GetResourceVersion(),
	}
}

func anyOf(m *envoyv1alpha1.Message) xds.Any {
	return xds.Any{
		TypeUrl: m.Type,
		Value:   m.Value,
	}
}

// EnvoyReconciler reconciles a Listener object.
type EnvoyReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	ResourceStore xds.ResourceStore
}

// AcceptResource decides whether the given Envoy resource sould be accepted.
func AcceptResource(obj envoyv1alpha1.Object, gvk schema.GroupVersionKind) (proto.Message, *kubernetes.AcceptanceError) {
	any := anyOf(obj.GetSpecMessage())

	// Verify that the type URL is acceptable for the kind.
	if xds.KindForTypename(any.TypeUrl) != gvk.Kind {
		return nil, &kubernetes.AcceptanceError{
			Reason:  "TypeAmbiguity",
			Message: fmt.Sprintf("invalid type %q for resource kind %q", any.TypeUrl, gvk.Kind),
		}
	}

	resource, err := xds.UnmarshalAny(&any)
	if err != nil {
		return nil, &kubernetes.AcceptanceError{
			Reason:  "InvalidFormat",
			Message: err.Error(),
		}
	}

	// Run protobuf validation for the resource.
	if err := xds.Validate(resource); err != nil {
		return nil, &kubernetes.AcceptanceError{
			Reason:  "FailedValidation",
			Message: fmt.Sprintf("protobuf validation error: %s", err),
		}
	}

	return resource, nil
}

// nolint(lll)
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=clusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=listeners,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=listeners/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=routeconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=routeconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=runtimes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=runtimes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=scopedrouteconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=scopedrouteconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=secrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=virtualhosts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=envoy.projectcontour.io,resources=virtualhosts/status,verbs=get;update;patch

// Reconcile ...
//
// nolint(lll)
func (e *EnvoyReconciler) Reconcile(req ctrl.Request, o runtime.Object, gvk schema.GroupVersionKind) (ctrl.Result, error) {
	ctx := context.Background()
	log := e.Log.WithValues(
		"kind", gvk.Kind,
		"name", req.NamespacedName,
		"resource", resourceOf(req.NamespacedName, gvk),
	)

	if err := e.Get(ctx, req.NamespacedName, o); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("deleting resource")
			e.ResourceStore.DeleteResource(resourceOf(req.NamespacedName, gvk))
			return ctrl.Result{}, nil
		}

		log.Error(err, "failed to fetch object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Convert to and Envoy resource object so we can access fields generically.
	obj, ok := o.(envoyv1alpha1.Object)
	if !ok {
		log.Error(fmt.Errorf("resource is not an Envoy object"), "invalid resource type",
			"type", fmt.Sprintf("%T", obj))
		return ctrl.Result{}, nil
	}

	accepted := kubernetes.NewAcceptedCondition(obj)

	// Do initial acceptance validation.
	resource, err := AcceptResource(obj, gvk)
	if err != nil {
		accepted.Status = metav1.ConditionFalse
		accepted.Reason = err.Reason
		accepted.Message = err.Message
	}

	var conditions []envoyv1alpha1.Condition

	// Preserve all conditions except "Accepted".
	for _, c := range obj.GetStatusConditions() {
		if c.Type != "Accepted" {
			conditions = append(conditions, c)
		}
	}

	conditions = append(conditions, *accepted)
	obj.SetStatusConditions(conditions)

	// Update the status condition on this object. The default
	// for new "Accepted" conditions is "True", switching to "False"
	// if we reject for any reason.
	if err := e.Client.Status().Update(ctx, obj); err != nil {
		// Requeue (rate-limited) if we lost an update race.
		if apierrors.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}

		log.Error(err, "failed to update .Status.Conditions")
		return ctrl.Result{}, err
	}

	switch accepted.Status {
	case metav1.ConditionFalse:
		log.Info("rejected resource", "reason", accepted.Reason, "message", accepted.Message)
		return ctrl.Result{}, nil
	default:
		log.Info("accepted resource")

	}

	log.Info("", "resource", resource)
	e.ResourceStore.UpdateResource(resourceOf(req.NamespacedName, gvk), versionOf(obj), resource)

	return ctrl.Result{}, nil
}

// SetupWithManager ...
func (e *EnvoyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	for _, factory := range factories {
		factory := factory
		gvk := must.GroupVersionKind(apiutil.GVKForObject(factory(), e.Scheme))

		if err := ctrl.NewControllerManagedBy(mgr).
			For(factory()).
			Complete(reconcile.Func(
				func(req ctrl.Request) (ctrl.Result, error) {
					obj := factory()
					return e.Reconcile(req, obj, gvk)
				})); err != nil {
			return fmt.Errorf("failed to set up %q reconciliation: %w",
				strings.ToLower(gvk.Kind), err)
		}
	}

	return nil
}
