package kubernetes

import (
	envoyv1alpha1 "github.com/jpeach/envoy-controller/api/v1alpha1"
	"github.com/jpeach/envoy-controller/pkg/must"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// NewScheme returns a *runtime.Scheme with the envoy-controller types registered.
func NewScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	must.Must(scheme.AddToScheme(s))
	must.Must(envoyv1alpha1.AddToScheme(s))

	return s
}

// NewClient does the needful.
func NewClient() (client.Client, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	return client.New(restConfig, client.Options{
		Scheme: NewScheme(),
	})
}

type CreateOptionFunc func(*client.CreateOptions)

func (f CreateOptionFunc) ApplyToCreate(o *client.CreateOptions) {
	if f != nil {
		f(o)
	}
}

var _ client.CreateOption = CreateOptionFunc(nil)
