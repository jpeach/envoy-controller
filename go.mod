module github.com/jpeach/envoy-controller

go 1.14

require (
	github.com/envoyproxy/go-control-plane v0.9.5
	github.com/go-logr/logr v0.1.0
	github.com/golang/protobuf v1.4.1
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/mattn/go-isatty v0.0.8
	github.com/spf13/cobra v0.0.5
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	google.golang.org/protobuf v1.24.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	k8s.io/apimachinery v0.18.3
	k8s.io/cli-runtime v0.18.3
	k8s.io/client-go v0.18.3
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.3.0
)
