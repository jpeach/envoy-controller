package xds

import (
	"google.golang.org/protobuf/proto"
)

// ResourceVersion is a globally unique identifier for a xDS resource.
// xDS resource names may be re-used, in which case the version is
// used to distinguish between them.
type ResourceVersion struct {
	Identifier string
	Version    string
}

// ResourceName is a globally unique name for an xDS resource.
type ResourceName string

// ResourceStore represents a store of Envoy resources. The resources
// are indexed and referred to by globally unique names. Resource stores
// are expected to map resources to Envoy API versions internally, if necessary.
type ResourceStore interface {
	UpdateResource(ResourceName, ResourceVersion, proto.Message)
	DeleteResource(ResourceName)
}
