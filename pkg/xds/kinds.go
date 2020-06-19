package xds

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type EnvoyVersion string

// EnvoyVersionUnknown ...
const EnvoyVersionUnknown = ""

// EnvoyVersion2 ...
const EnvoyVersion2 = "v2"

// EnvoyVersion3 ...
const EnvoyVersion3 = "v3"

// ProtobufForKind returns the versioned MessageType for the given resource kind.
func ProtobufForKind(vers EnvoyVersion, kind string) (protoreflect.MessageType, error) {
	var protoForKind map[string]string

	switch vers {
	case EnvoyVersion2:
		protoForKind = map[string]string{
			"Listener":                 "envoy.api.v2.Listener",
			"Cluster":                  "envoy.api.v2.Cluster",
			"RouteConfiguration":       "envoy.api.v2.RouteConfiguration",
			"ScopedRouteConfiguration": "envoy.api.v2.ScopedRouteConfiguration",
			"Secret":                   "envoy.api.v2.auth.Secret",
			"Runtime":                  "envoy.service.discovery.v2.Runtime",
			"VirtualHost":              "envoy.api.v2.route.VirtualHost",
		}
	case EnvoyVersion3:
		protoForKind = map[string]string{
			"Listener":                 "envoy.config.listener.v3.Listener",
			"Cluster":                  "envoy.config.cluster.v3.Cluster",
			"RouteConfiguration":       "envoy.config.route.v3.RouteConfiguration",
			"ScopedRouteConfiguration": "envoy.config.route.v3.ScopedRouteConfiguration",
			"Secret":                   "envoy.extensions.transport_sockets.tls.v3.Secret",
			"Runtime":                  "envoy.service.runtime.v3.Runtime",
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
