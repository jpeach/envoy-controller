package controllers

import (
	"context"
	"fmt"
	"path"
	"strings"

	envoyv1alpha1 "github.com/jpeach/envoy-controller/api/v1alpha1"
	"github.com/jpeach/envoy-controller/pkg/must"
	"github.com/jpeach/envoy-controller/pkg/xds"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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

func resourceOf(name types.NamespacedName, gvk schema.GroupVersionKind) string {
	return strings.ToLower(path.Join(name.Namespace, gvk.Kind, name.Name))
}

func anyOf(m envoyv1alpha1.Message) xds.Any {
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
func (e *EnvoyReconciler) Reconcile(req ctrl.Request, obj runtime.Object, gvk schema.GroupVersionKind) (ctrl.Result, error) {
	ctx := context.Background()
	log := e.Log.WithValues(
		"kind", gvk.Kind,
		"name", req.NamespacedName,
		"resource", resourceOf(req.NamespacedName, gvk),
	)

	if err := e.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("deleting resource")
			e.ResourceStore.DeleteResource(resourceOf(req.NamespacedName, gvk))
			return ctrl.Result{}, nil
		}

		log.Error(err, "failed to fetch object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var any xds.Any

	switch obj := obj.(type) {
	case *envoyv1alpha1.Listener:
		any = anyOf(obj.Spec.Listener)
	case *envoyv1alpha1.Cluster:
		any = anyOf(obj.Spec.Cluster)
	case *envoyv1alpha1.RouteConfiguration:
		any = anyOf(obj.Spec.RouteConfiguration)
	case *envoyv1alpha1.ScopedRouteConfiguration:
		any = anyOf(obj.Spec.ScopedRouteConfiguration)
	case *envoyv1alpha1.Secret:
		any = anyOf(obj.Spec.Secret)
	case *envoyv1alpha1.Runtime:
		any = anyOf(obj.Spec.Runtime)
	case *envoyv1alpha1.VirtualHost:
		any = anyOf(obj.Spec.VirtualHost)
	default:
		log.Info("invalid resource type", "type", fmt.Sprintf("%T", obj))
		// TODO(jpeach) set error status.
	}

	// Verify that the type URL is acceptable for the kind.
	if xds.KindForTypename(any.TypeUrl) != gvk.Kind {
		// TODO(jpeach) set error status.
		log.Error(fmt.Errorf("type %s is not valid for a resource of kind %q", any.TypeUrl, gvk.Kind),
			"invalid spec.type")
		return ctrl.Result{}, nil
	}

	resource, err := xds.UnmarshalAny(&any)
	if err != nil {
		/// TODO(jpeach) set error status.
		log.Error(err, "failed to unmarshal Envoy resource")
		return ctrl.Result{}, nil
	}

	log.Info("updated resource")
	// TODO(jpeach) update status

	e.ResourceStore.UpdateResource(resourceOf(req.NamespacedName, gvk), resource)

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
