package xds

import (
	"context"
	"fmt"
	"net"

	clusterserviceV2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	endpointserviceV2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	listenerserviceV2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	routeserviceV2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	clusterserviceV3 "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoveryserviceV2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	runtimeserviceV2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	secretserviceV2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	discoveryserviceV3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointserviceV3 "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerserviceV3 "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeserviceV3 "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeserviceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretserviceV3 "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	cacheV2 "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	cacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/log"
	serverV2 "github.com/envoyproxy/go-control-plane/pkg/server/v2"
	serverV3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Server is a handle to a GRPC server that implements the xDS v2 and v3 protocols.
type Server struct {
	v2   serverV2.Server
	v3   serverV3.Server
	grpc *grpc.Server
}

var _ ResourceStore = &Server{}

// NewServer returns a new xDS server for both the v2 and v3 Envoy
// API. The internal cache always uses the identity node hash since
// the deployment model is one envoy-controller for each Envoy server.
func NewServer(options ...grpc.ServerOption) *Server {
	logger := ctrl.Log.WithName("xds")

	// The logr API just plain sucks. Why on earth it is baked
	// into controller-runtime is beyond me.
	l := log.LoggerFuncs{
		DebugFunc: func(format string, args ...interface{}) {
			logger.V(1).Info(fmt.Sprintf(format, args...))
		},
		InfoFunc: func(format string, args ...interface{}) {
			logger.Info(fmt.Sprintf(format, args...))
		},
		WarnFunc: func(format string, args ...interface{}) {
			logger.Info(fmt.Sprintf(format, args...))
		},
		ErrorFunc: func(format string, args ...interface{}) {
			logger.Info(fmt.Sprintf(format, args...))
		},
	}

	srv := Server{
		v2: serverV2.NewServer(
			context.Background(),
			cacheV2.NewSnapshotCache(true /* ads */, cacheV2.IDHash{}, l),
			nil, /* callbacks */
		),

		v3: serverV3.NewServer(
			context.Background(),
			cacheV3.NewSnapshotCache(true /* ads */, cacheV3.IDHash{}, l),
			nil, /* callbacks */
		),

		grpc: grpc.NewServer(options...),
	}

	clusterserviceV3.RegisterClusterDiscoveryServiceServer(srv.grpc, srv.v3)
	discoveryserviceV3.RegisterAggregatedDiscoveryServiceServer(srv.grpc, srv.v3)
	endpointserviceV3.RegisterEndpointDiscoveryServiceServer(srv.grpc, srv.v3)
	listenerserviceV3.RegisterListenerDiscoveryServiceServer(srv.grpc, srv.v3)
	routeserviceV3.RegisterRouteDiscoveryServiceServer(srv.grpc, srv.v3)
	runtimeserviceV3.RegisterRuntimeDiscoveryServiceServer(srv.grpc, srv.v3)
	secretserviceV3.RegisterSecretDiscoveryServiceServer(srv.grpc, srv.v3)

	// go-control-plane doesn't support scoped routes or virtualhosts:
	// 	https://github.com/envoyproxy/go-control-plane/issues/310
	// 	https://github.com/envoyproxy/go-control-plane/issues/309

	clusterserviceV2.RegisterClusterDiscoveryServiceServer(srv.grpc, srv.v2)
	discoveryserviceV2.RegisterAggregatedDiscoveryServiceServer(srv.grpc, srv.v2)
	endpointserviceV2.RegisterEndpointDiscoveryServiceServer(srv.grpc, srv.v2)
	listenerserviceV2.RegisterListenerDiscoveryServiceServer(srv.grpc, srv.v2)
	routeserviceV2.RegisterRouteDiscoveryServiceServer(srv.grpc, srv.v2)
	runtimeserviceV2.RegisterRuntimeDiscoveryServiceServer(srv.grpc, srv.v2)
	secretserviceV2.RegisterSecretDiscoveryServiceServer(srv.grpc, srv.v2)

	return &srv
}

// Serve serves the GRPC endpoint on the given listener, returning only when it is stopped.
func (srv *Server) Serve(listener net.Listener) error {
	return srv.grpc.Serve(listener)
}

// Stop requests a graceful stop of the xDS GRPC server.
func (srv *Server) Stop() {
	srv.grpc.GracefulStop()
}

// Start runs the xDRS GRPC server on the given listener, return
// only when the server is stopped by stopchan closing, or when it
// fails with an error.
func (srv *Server) Start(listener net.Listener, stopChan <-chan struct{}) error {
	errChan := make(chan error)

	go func() {
		errChan <- srv.grpc.Serve(listener)
	}()

	select {
	case err := <-errChan:
		return err
	case <-stopChan:
		srv.Stop()
		return nil
	}
}

// UpdateResource ...
func (srv *Server) UpdateResource(name ResourceName, vers ResourceVersion, message proto.Message) {
	// TODO(jpeach) Enforce the invariant that names are globally unique.
	switch VersionForMessage(message.ProtoReflect().Descriptor()) {
	case EnvoyVersion2:
	case EnvoyVersion3:
	}
}

// DeleteResource ...
func (srv *Server) DeleteResource(name ResourceName) {
	// name is globally unique, so we can safely delete the
	// corresponding entry from both the v2 and v3 resources.
}
