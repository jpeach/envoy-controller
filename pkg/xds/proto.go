package xds

import (
	protov1 "github.com/golang/protobuf/proto" // nolint(staticcheck)
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/proto"
)

// ProtoV1 converts a V2 message to V1.
func ProtoV1(message proto.Message) protov1.Message {
	return protov1.MessageV1(message)
}

// ProtoV2 converts a V1 message to V2.
func ProtoV2(message protov1.Message) proto.Message {
	return protov1.MessageV2(message)
}

// MarshalAny marshals a message into an any.Any type.
func MarshalAny(message proto.Message) (*any.Any, error) {
	return ptypes.MarshalAny(ProtoV1(message))
}
