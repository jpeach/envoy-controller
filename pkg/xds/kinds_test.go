package xds

import (
	"strings"
	"testing"

	envoyv1alpha1 "github.com/jpeach/envoy-controller/api/v1alpha1"
	"github.com/jpeach/envoy-controller/pkg/must"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestKindsForScheme(t *testing.T) {
	s := runtime.NewScheme()
	must.Must(envoyv1alpha1.AddToScheme(s))

	kindNames := map[string]bool{}

	for _, k := range Kinds() {
		kindNames[k] = true
	}

	// Make sure every listed kind is a legit CRD type.
	for _, k := range Kinds() {
		gvk := envoyv1alpha1.GroupVersion.WithKind(k)
		if !s.Recognizes(gvk) {
			t.Errorf("unrecognized kind %q", k)
		}
	}

	// Make sure that we have all the kinds listed.
	for kindName, typeInfo := range s.KnownTypes(envoyv1alpha1.GroupVersion) {
		if strings.Contains(kindName, "List") {
			continue
		}

		if !strings.Contains(typeInfo.PkgPath(), "envoy-controller/api/v1alpha1") {
			continue
		}

		assert.Contains(t, kindNames, kindName)
	}
}

func TestProtoForKind(t *testing.T) {
	for _, vers := range []EnvoyVersion{EnvoyVersion2, EnvoyVersion3} {
		for _, k := range Kinds() {
			msg, err := ProtobufForKind(vers, k)
			require.NoErrorf(t, err, "no Envoy %s protobuf for kind %q", vers, k)

			msgName := msg.Descriptor().FullName()
			assert.Equalf(t, k, KindForTypename(string(msgName)),
				"no kind for %q message", msgName)

			assert.Equalf(t, vers, VersionForMessage(msg.Descriptor()),
				"no Envoy version for %q message", msgName)
		}
	}
}
