package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	envoyv1alpha1 "github.com/jpeach/envoy-controller/api/v1alpha1"
	"github.com/jpeach/envoy-controller/pkg/kubernetes"
	"github.com/jpeach/envoy-controller/pkg/must"
	"github.com/jpeach/envoy-controller/pkg/xds"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/printers"
)

func envoyVersion(flags *pflag.FlagSet) xds.EnvoyVersion {
	if must.Bool(flags.GetBool("2")) {
		return xds.EnvoyVersion2
	}

	if must.Bool(flags.GetBool("3")) {
		return xds.EnvoyVersion3
	}

	return xds.EnvoyVersion3
}

// NamespaceOrDefault returns the namespace ns, or "default" if ns is empty.
func NamespaceOrDefault(ns string) string {
	if ns != "" {
		return ns
	}

	return metav1.NamespaceDefault
}

func NewCreateCommand() *cobra.Command {
	var kinds = []string{
		"Listener",
		"Cluster",
		"RouteConfiguration",
		"ScopedRouteConfiguration",
		"Secret",
		"Runtime",
		"VirtualHost",
	}

	cmd := cobra.Command{
		Use:   "create RESOURCE NAME [OPTIONS]",
		Short: "Create an Envoy resource from a file or stdin",
	}

	for _, k := range kinds {
		k := k
		kindCmd := &cobra.Command{
			Use:   fmt.Sprintf("%s NAME [OPTIONS]", strings.ToLower(k)),
			Short: fmt.Sprintf("Create an Envoy %s resource from a file or stdin", k),
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				name := types.NamespacedName{
					Namespace: NamespaceOrDefault(must.String(cmd.Flags().GetString("namespace"))),
					Name:      args[0],
				}

				var input []byte
				var err error

				if must.Bool(cmd.Flags().GetBool("2")) && must.Bool(cmd.Flags().GetBool("3")) {
					return ExitErrorf(EX_USAGE, "multiple Envoy API versions specified")
				}

				mtype, err := xds.ProtobufForKind(envoyVersion(cmd.Flags()), k)
				if err != nil {
					return &ExitError{EX_CONFIG, err}
				}

				if fname := must.String(cmd.Flags().GetString("filename")); fname != "-" {
					input, err = ioutil.ReadFile(fname)
				} else {
					input, err = ioutil.ReadAll(os.Stdin)
				}

				if err != nil {
					return &ExitError{Code: EX_DATAERR, Err: err}
				}

				obj, err := createResourceV3(k, name, input, mtype)
				if err != nil {
					return &ExitError{Code: EX_FAIL, Err: err}
				}

				if must.String(cmd.Flags().GetString("output")) == "" {
					return createResource(obj)
				} else {
					return formatResource(obj, must.String(cmd.Flags().GetString("output")))
				}
			},
		}

		cmd.AddCommand(Defaults(kindCmd))
	}

	cmd.PersistentFlags().StringP("namespace", "n", "", "The namespace in which to create the resource.")
	cmd.PersistentFlags().StringP("filename", "f", "-", "Filename used to create the resource.")
	cmd.PersistentFlags().StringP("output", "o", "", "Output the object as YAML or JSON instead of creating it.")
	cmd.PersistentFlags().BoolP("3", "3", false, "Output the object as YAML or JSON instead of creating it.")
	cmd.PersistentFlags().BoolP("2", "2", false, "Output the object as YAML or JSON instead of creating it.")

	return &cmd
}

func createResource(obj runtime.Object) error {
	client, err := kubernetes.NewClient()
	if err != nil {
		return &ExitError{EX_CONFIG, err}
	}

	if err := client.Create(context.Background(), obj, kubernetes.CreateOptionFunc(nil)); err != nil {
		return &ExitError{EX_FAIL, err}
	}

	return nil
}

func formatResource(obj runtime.Object, format string) error {
	switch format {
	case "yaml":
		p := printers.YAMLPrinter{}
		return p.PrintObj(obj, os.Stdout)
	case "json":
		p := printers.JSONPrinter{}
		return p.PrintObj(obj, os.Stdout)
	default:
		return ExitErrorf(EX_USAGE, "unsupported output format %q", format)
	}
}

func createResourceV3(kind string, name types.NamespacedName, in []byte, mtype protoreflect.MessageType) (runtime.Object, error) {
	// Unmarshal the JSON into an instance of the message type.
	protoMessage := mtype.New().Interface()
	if err := protojson.Unmarshal(in, protoMessage); err != nil {
		return nil, err
	}

	//  TODO(jpeach): if the protobug object has a "name" field,
	// force it to match the fully qualified Kubernetes resource
	// name.

	// Marshal the Any message payload.
	anyMessage, err := xds.MarshalAny(protoMessage)
	if err != nil {
		return nil, err
	}

	objectMeta := metav1.ObjectMeta{
		Name:              name.Name,
		Namespace:         name.Namespace,
		CreationTimestamp: metav1.Now(),
	}

	message := envoyv1alpha1.Message{
		Type:  anyMessage.GetTypeUrl(),
		Value: anyMessage.GetValue(),
	}

	var obj runtime.Object

	switch kind {
	case "Listener":
		obj = &envoyv1alpha1.Listener{
			ObjectMeta: objectMeta,
			Spec:       envoyv1alpha1.ListenerSpec{Listener: message},
		}

	case "Cluster":
		obj = &envoyv1alpha1.Cluster{
			ObjectMeta: objectMeta,
			Spec:       envoyv1alpha1.ClusterSpec{Cluster: message},
		}
	case "RouteConfiguration":
		obj = &envoyv1alpha1.RouteConfiguration{
			ObjectMeta: objectMeta,
			Spec:       envoyv1alpha1.RouteConfigurationSpec{RouteConfiguration: message},
		}

	case "ScopedRouteConfiguration":
		obj = &envoyv1alpha1.ScopedRouteConfiguration{
			ObjectMeta: objectMeta,
			Spec:       envoyv1alpha1.ScopedRouteConfigurationSpec{ScopedRouteConfiguration: message},
		}

	case "Secret":
		obj = &envoyv1alpha1.Secret{
			ObjectMeta: objectMeta,
			Spec:       envoyv1alpha1.SecretSpec{Secret: message},
		}

	case "Runtime":
		obj = &envoyv1alpha1.Runtime{
			ObjectMeta: objectMeta,
			Spec:       envoyv1alpha1.RuntimeSpec{Runtime: message},
		}

	case "VirtualHost":
		obj = &envoyv1alpha1.VirtualHost{
			ObjectMeta: objectMeta,
			Spec:       envoyv1alpha1.VirtualHostSpec{VirtualHost: message},
		}

	default:
		panic(fmt.Sprintf("invalid kind %q", kind))
	}

	// YAML output requires us to set the GVK explicitly.
	obj.GetObjectKind().SetGroupVersionKind(envoyv1alpha1.GroupVersion.WithKind(kind))

	return obj, nil
}
