package controllers

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler interface {
	Reconcile(req ctrl.Request) (ctrl.Result, error)
	SetupWithManager(mgr ctrl.Manager) error
}

func New(kind string, c client.Client, s *runtime.Scheme) Reconciler {
	switch kind {
	case "Listener":
		return &ListenerReconciler{
			Client: c,
			Scheme: s,
			Log:    ctrl.Log.WithName("controllers").WithName(kind),
		}

	case "Cluster":
		return &ClusterReconciler{
			Client: c,
			Scheme: s,
			Log:    ctrl.Log.WithName("controllers").WithName(kind),
		}

	case "RouteConfiguration":
		return &RouteConfigurationReconciler{
			Client: c,
			Scheme: s,
			Log:    ctrl.Log.WithName("controllers").WithName(kind),
		}

	case "ScopedRouteConfiguration":
		return &ScopedRouteConfigurationReconciler{
			Client: c,
			Scheme: s,
			Log:    ctrl.Log.WithName("controllers").WithName(kind),
		}

	case "Secret":
		return &SecretReconciler{
			Client: c,
			Scheme: s,
			Log:    ctrl.Log.WithName("controllers").WithName(kind),
		}

	case "Runtime":
		return &RuntimeReconciler{
			Client: c,
			Scheme: s,
			Log:    ctrl.Log.WithName("controllers").WithName(kind),
		}

	case "VirtualHost":
		return &VirtualHostReconciler{
			Client: c,
			Scheme: s,
			Log:    ctrl.Log.WithName("controllers").WithName(kind),
		}

	default:
		panic(fmt.Sprintf("no reconciler for kind %q", kind))
	}
}
