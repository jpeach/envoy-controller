package xds

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// EnvoyVersion is the version of the Envoy API.
type EnvoyVersion string

// EnvoyVersionUnknown ...
const EnvoyVersionUnknown = ""

// EnvoyVersion2 ...
const EnvoyVersion2 = "v2"

// EnvoyVersion3 ...
const EnvoyVersion3 = "v3"

// Kinds returns the names of the Envoy API kinds.
func Kinds() []string {
	return []string{
		"Cluster",
		"ClusterLoadAssignment",
		"Listener",
		"RouteConfiguration",
		"Runtime",
		"ScopedRouteConfiguration",
		"Secret",
		"VirtualHost",
	}
}

// ProtobufForKind returns the versioned MessageType for the given resource kind.
func ProtobufForKind(vers EnvoyVersion, kind string) (protoreflect.MessageType, error) {
	var protoForKind map[string]string

	switch vers {
	case EnvoyVersion2:
		protoForKind = map[string]string{
			"Cluster":                  "envoy.api.v2.Cluster",
			"ClusterLoadAssignment":    "envoy.api.v2.ClusterLoadAssignment",
			"Listener":                 "envoy.api.v2.Listener",
			"RouteConfiguration":       "envoy.api.v2.RouteConfiguration",
			"Runtime":                  "envoy.service.discovery.v2.Runtime",
			"ScopedRouteConfiguration": "envoy.api.v2.ScopedRouteConfiguration",
			"Secret":                   "envoy.api.v2.auth.Secret",
			"VirtualHost":              "envoy.api.v2.route.VirtualHost",
		}
	case EnvoyVersion3:
		protoForKind = map[string]string{
			"Cluster":                  "envoy.config.cluster.v3.Cluster",
			"ClusterLoadAssignment":    "envoy.config.endpoint.v3.ClusterLoadAssignment",
			"Listener":                 "envoy.config.listener.v3.Listener",
			"RouteConfiguration":       "envoy.config.route.v3.RouteConfiguration",
			"Runtime":                  "envoy.service.runtime.v3.Runtime",
			"ScopedRouteConfiguration": "envoy.config.route.v3.ScopedRouteConfiguration",
			"Secret":                   "envoy.extensions.transport_sockets.tls.v3.Secret",
			"VirtualHost":              "envoy.config.route.v3.VirtualHost",
		}
	default:
		return nil, errors.New("unsupported Envoy API version")
	}

	messageType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(protoForKind[kind]))
	if err != nil {
		return nil, fmt.Errorf("protobuf message type %q: %s", protoForKind[kind], err)
	}

	return messageType, nil
}

// KindForTypename returns the Kubernetes resource kind for the
// given type name, which may be either a protobuf.Fullname, or a
// any.Any type URL.
func KindForTypename(typeName string) string {
	messageKinds := map[string]string{
		"envoy.api.v2.Cluster":                             "Cluster",
		"envoy.api.v2.ClusterLoadAssignment":               "ClusterLoadAssignment",
		"envoy.api.v2.Listener":                            "Listener",
		"envoy.api.v2.RouteConfiguration":                  "RouteConfiguration",
		"envoy.api.v2.ScopedRouteConfiguration":            "ScopedRouteConfiguration",
		"envoy.api.v2.auth.Secret":                         "Secret",
		"envoy.api.v2.route.VirtualHost":                   "VirtualHost",
		"envoy.config.cluster.v3.Cluster":                  "Cluster",
		"envoy.config.endpoint.v3.ClusterLoadAssignment":   "ClusterLoadAssignment",
		"envoy.config.listener.v3.Listener":                "Listener",
		"envoy.config.route.v3.RouteConfiguration":         "RouteConfiguration",
		"envoy.config.route.v3.ScopedRouteConfiguration":   "ScopedRouteConfiguration",
		"envoy.config.route.v3.VirtualHost":                "VirtualHost",
		"envoy.extensions.transport_sockets.tls.v3.Secret": "Secret",
		"envoy.service.discovery.v2.Runtime":               "Runtime",
		"envoy.service.runtime.v3.Runtime":                 "Runtime",
	}

	// If we got a TypeURL, strip the prefix.
	typeName = strings.Replace(typeName, "type.googleapis.com/", "", 1)

	return messageKinds[typeName]
}

// VersionForMessage returns the Envoy API version that matches the given message type.
func VersionForMessage(m protoreflect.MessageDescriptor) EnvoyVersion {
	messageVersions := map[protoreflect.FullName]EnvoyVersion{
		"envoy.api.v2.Cluster":                             EnvoyVersion2,
		"envoy.api.v2.ClusterLoadAssignment":               EnvoyVersion2,
		"envoy.api.v2.Listener":                            EnvoyVersion2,
		"envoy.api.v2.RouteConfiguration":                  EnvoyVersion2,
		"envoy.api.v2.ScopedRouteConfiguration":            EnvoyVersion2,
		"envoy.api.v2.auth.Secret":                         EnvoyVersion2,
		"envoy.api.v2.route.VirtualHost":                   EnvoyVersion2,
		"envoy.service.discovery.v2.Runtime":               EnvoyVersion2,
		"envoy.config.cluster.v3.Cluster":                  EnvoyVersion3,
		"envoy.config.endpoint.v3.ClusterLoadAssignment":   EnvoyVersion3,
		"envoy.config.listener.v3.Listener":                EnvoyVersion3,
		"envoy.config.route.v3.RouteConfiguration":         EnvoyVersion3,
		"envoy.config.route.v3.ScopedRouteConfiguration":   EnvoyVersion3,
		"envoy.config.route.v3.VirtualHost":                EnvoyVersion3,
		"envoy.extensions.transport_sockets.tls.v3.Secret": EnvoyVersion3,
		"envoy.service.runtime.v3.Runtime":                 EnvoyVersion3,
	}

	if vers, ok := messageVersions[m.FullName()]; ok {
		return vers
	}

	return EnvoyVersionUnknown
}
