package xds

import (
	protov1 "github.com/golang/protobuf/proto" // nolint(staticcheck)
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/proto"
)

// Any ...
type Any = any.Any

// ProtoV1 converts a V2 message to V1.
func ProtoV1(message proto.Message) protov1.Message {
	return protov1.MessageV1(message)
}

// ProtoV2 converts a V1 message to V2.
func ProtoV2(message protov1.Message) proto.Message {
	return protov1.MessageV2(message)
}

// MarshalAny marshals a message into an Any type.
func MarshalAny(message proto.Message) (*Any, error) {
	return ptypes.MarshalAny(ProtoV1(message))
}

// UnmarshalAny unmarshals an Any message.
func UnmarshalAny(anyMessage *Any) (proto.Message, error) {
	var x ptypes.DynamicAny

	if err := ptypes.UnmarshalAny(anyMessage, &x); err != nil {
		return nil, err
	}

	return ProtoV2(x.Message), nil
}

// Validate calls the `Validate` method of the proto.Message, if it has one.
func Validate(message proto.Message) error {
	if v, ok := interface{}(message).(interface{ Validate() error }); ok {
		return v.Validate()
	}

	return nil
}
