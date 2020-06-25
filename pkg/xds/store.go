package xds

import "google.golang.org/protobuf/proto"

// ResourceStore represents a store of Envoy resources. The resources
// are indexed and referred to by globally unique names. Resource stores
// are expected to map resources to Envoy API versions internally, if necessary.
type ResourceStore interface {
	UpdateResource(string, proto.Message)
	DeleteResource(string)
}
