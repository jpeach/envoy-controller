package bootstrap

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jpeach/envoy-controller/pkg/xds"

	envoy_config_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/proto"
)

// Bootstrap ...
type Bootstrap = envoy_config_bootstrap_v3.Bootstrap

// Address ...
type Address = envoy_config_core_v3.Address

// APIVersion ...
type APIVersion = envoy_config_core_v3.ApiVersion

const (
	ApiVersion_AUTO APIVersion = envoy_config_core_v3.ApiVersion_AUTO //nolint
	ApiVersion_V2   APIVersion = envoy_config_core_v3.ApiVersion_V2   //nolint
	ApiVersion_V3   APIVersion = envoy_config_core_v3.ApiVersion_V3   //nolint
)

// Option ...
type Option func(*Bootstrap)

// NodeID sets the envoy node ID.
func NodeID(s string) Option {
	return func(b *Bootstrap) {
		b.Node.Id = s
	}
}

// NodeCluster sets the envoy node Cluster name.
func NodeCluster(s string) Option {
	return func(b *Bootstrap) {
		b.Node.Cluster = s
	}
}

// ResourceVersion sets the default resource API version Envoy will ask for.
func ResourceVersion(vers APIVersion) Option {
	return func(b *Bootstrap) {
		b.DynamicResources.LdsConfig.ResourceApiVersion = vers
		b.DynamicResources.CdsConfig.ResourceApiVersion = vers
	}
}

// SetNodeOnFirstMessageOnly tells Envoy to only send the Node message once.
func SetNodeOnFirstMessageOnly(value bool) Option {
	return func(b *Bootstrap) {
		b.DynamicResources.AdsConfig.SetNodeOnFirstMessageOnly = value
	}
}

// AdminAccessLog set the access log path for the admin endpoint.
func AdminAccessLog(path string) Option {
	return func(b *Bootstrap) {
		b.Admin.AccessLogPath = path
	}
}

// AdminAddress sets the address the admin server will listen on.
func AdminAddress(addr *Address) Option {
	return func(b *Bootstrap) {
		b.Admin.Address = addr
	}
}

// ManagementAddress sets the address to connect to the xDS management server.
func ManagementAddress(addr *Address) Option {
	return func(b *Bootstrap) {
		for _, c := range b.StaticResources.Clusters {
			// Find the xDS cluster.
			if c.LoadAssignment.ClusterName != "xds" {
				continue
			}

			// Stuff the address into a host endpoint.
			ep := &envoy_config_endpoint_v3.LbEndpoint{
				HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
					Endpoint: &envoy_config_endpoint_v3.Endpoint{
						Address: addr,
					},
				},
			}

			// Add  the endpoint to the cluster.
			c.LoadAssignment.Endpoints[0].LbEndpoints = append(
				c.LoadAssignment.Endpoints[0].LbEndpoints, ep,
			)
		}
	}
}

// EnableIncrementalDiscovery enables incremental (Delta) xDS.
func EnableIncrementalDiscovery() Option {
	return func(b *Bootstrap) {
		b.DynamicResources.AdsConfig.ApiType = envoy_config_core_v3.ApiConfigSource_DELTA_GRPC
	}
}

// New returns a new bootstrap protobuf.
func New(options ...Option) (proto.Message, error) {
	type Admin = envoy_config_bootstrap_v3.Admin                                 //nolint
	type ApiConfigSource = envoy_config_core_v3.ApiConfigSource                  //nolint
	type Cluster = envoy_config_cluster_v3.Cluster                               //nolint
	type ConfigSource = envoy_config_core_v3.ConfigSource                        //nolint
	type DynamicResources = envoy_config_bootstrap_v3.Bootstrap_DynamicResources //nolint
	type GrpcService = envoy_config_core_v3.GrpcService                          //nolint
	type Node = envoy_config_core_v3.Node                                        //nolint
	type StaticResources = envoy_config_bootstrap_v3.Bootstrap_StaticResources   //nolint

	b := &Bootstrap{
		Node: &Node{
			Id:       "",
			Cluster:  "",
			Metadata: nil,
			Locality: nil,
		},
		Admin: &Admin{
			AccessLogPath: "",
			ProfilePath:   "",
			Address:       nil,
			SocketOptions: nil,
		},
		StaticResources: &StaticResources{
			Listeners: nil,
			Clusters: []*Cluster{
				//nolint(gofmt)
				&Cluster{
					Name:                 "xds",
					ConnectTimeout:       ptypes.DurationProto(time.Second * 10),
					Http2ProtocolOptions: &envoy_config_core_v3.Http2ProtocolOptions{},
					ClusterDiscoveryType: &envoy_config_cluster_v3.Cluster_Type{
						Type: envoy_config_cluster_v3.Cluster_STATIC,
					},
					LoadAssignment: &envoy_config_endpoint_v3.ClusterLoadAssignment{
						ClusterName: "xds",
						Endpoints: []*envoy_config_endpoint_v3.LocalityLbEndpoints{
							&envoy_config_endpoint_v3.LocalityLbEndpoints{
								LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{},
							}},
					},
				},
			},
		},
		DynamicResources: &DynamicResources{
			CdsConfig: &ConfigSource{
				ResourceApiVersion: ApiVersion_AUTO,
				ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_Ads{
					Ads: &envoy_config_core_v3.AggregatedConfigSource{},
				},
			},
			LdsConfig: &ConfigSource{
				ResourceApiVersion: ApiVersion_AUTO,
				ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_Ads{
					Ads: &envoy_config_core_v3.AggregatedConfigSource{},
				},
			},
			AdsConfig: &ApiConfigSource{
				ApiType:                   envoy_config_core_v3.ApiConfigSource_GRPC,
				TransportApiVersion:       ApiVersion_V3,
				ClusterNames:              []string{},
				RefreshDelay:              nil,
				RequestTimeout:            nil,
				RateLimitSettings:         nil,
				SetNodeOnFirstMessageOnly: false,
				GrpcServices: []*GrpcService{
					//nolint(gofmt)
					&GrpcService{
						TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{
								ClusterName: "xds",
							},
						},
					},
				},
			},
		},
		ClusterManager:             nil,
		HdsConfig:                  nil,
		FlagsPath:                  "",
		StatsSinks:                 nil,
		StatsConfig:                nil,
		StatsFlushInterval:         nil,
		Watchdog:                   nil,
		Tracing:                    nil,
		LayeredRuntime:             nil,
		OverloadManager:            nil,
		EnableDispatcherStats:      false,
		StatsServerVersionOverride: nil,
		UseTcpForDnsLookups:        false,
	}

	for _, o := range options {
		o(b)
	}

	if err := b.Validate(); err != nil {
		return nil, err
	}

	return xds.ProtoV2(b), nil
}

// NewAddress parses the addr string into a Envoy Address that can
// subsequently be used in an Option. If the address contains ":", it
// is assumed to be a socket "address:port" spec, otherwise is it the
// path to a pipe.
func NewAddress(addr string) (*Address, error) {
	address := Address{}

	if strings.Contains(addr, ":") {
		parts := strings.SplitN(addr, ":", 2)

		addr := parts[0]
		if addr == "" {
			addr = "0.0.0.0"
		}

		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid socket address %q: %w", addr, err)
		}

		address.Address = &envoy_config_core_v3.Address_SocketAddress{
			SocketAddress: &envoy_config_core_v3.SocketAddress{
				Address: addr,
				PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
					PortValue: uint32(port),
				},
			},
		}
	} else {
		address.Address = &envoy_config_core_v3.Address_Pipe{
			Pipe: &envoy_config_core_v3.Pipe{
				Path: addr,
				Mode: 0640,
			},
		}
	}

	return &address, nil
}
