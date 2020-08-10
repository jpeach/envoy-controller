package kubernetes

import (
	"github.com/jpeach/envoy-controller/api/v1alpha1"
	"github.com/jpeach/envoy-controller/pkg/must"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Condition aliases v1alpha1.Condition until metav1.Condition is
// available, which should be  Kubernetes 1.19.
type Condition = v1alpha1.Condition

// NewAcceptedCondition returns a *v1alpha1.Condition initialized for the given runtime.Object.
func NewAcceptedCondition(obj runtime.Object) *Condition {
	m := must.Object(meta.Accessor(obj))
	c := v1alpha1.Condition{
		Type:               "Accepted",
		Status:             metav1.ConditionTrue,
		ObservedGeneration: m.GetGeneration(),
		LastTransitionTime: metav1.Now(),
		Reason:             "",
		Message:            "",
	}

	return &c
}

// AcceptanceError captures the reason that a resource update was not accepted.
type AcceptanceError struct {
	Reason  string
	Message string
}
