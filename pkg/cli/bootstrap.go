package cli

import (
	"os"

	"github.com/jpeach/envoy-controller/pkg/bootstrap"
	"github.com/jpeach/envoy-controller/pkg/must"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// NewBootstrapCommand returns a command that writes an Envoy bootstrap file.
func NewBootstrapCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "bootstrap [OPTIONS]",
		Short: "Write an Envoy bootstrap configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			vers := bootstrap.ApiVersion_AUTO
			if must.Bool(cmd.Flags().GetBool("2")) && must.Bool(cmd.Flags().GetBool("3")) {
				return ExitErrorf(EX_USAGE, "multiple Envoy API versions specified")
			}
			if must.Bool(cmd.Flags().GetBool("2")) {
				vers = bootstrap.ApiVersion_V2
			}
			if must.Bool(cmd.Flags().GetBool("3")) {
				vers = bootstrap.ApiVersion_V3
			}

			xdsAddr, err := bootstrap.NewAddress(must.String(cmd.Flags().GetString("xds-address")))
			if err != nil {
				return ExitErrorf(EX_CONFIG, "invalid xDS address: %s", err)
			}

			adminAddr, err := bootstrap.NewAddress(must.String(cmd.Flags().GetString("admin-address")))
			if err != nil {
				return ExitErrorf(EX_CONFIG, "invalid admin address: %s", err)
			}

			out := cmd.OutOrStdout()

			if path := must.String(cmd.Flags().GetString("filename")); path != "-" {
				file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0640) //nolint(gosec)
				if err != nil {
					return ExitError{EX_FAIL, err}
				}

				defer func() {
					must.Must(file.Sync())
					must.Must(file.Close())
				}()

				out = file
			}

			opts := []bootstrap.Option{
				bootstrap.NodeCluster(must.String(os.Hostname())),
				bootstrap.NodeID(must.String(os.Hostname())),
				bootstrap.ResourceVersion(vers),
				bootstrap.ManagementClusterName(must.String(cmd.Flags().GetString("xds-clustername"))),
				bootstrap.ManagementAddress(xdsAddr),
				bootstrap.AdminAddress(adminAddr),
				bootstrap.AdminAccessLog(must.String(cmd.Flags().GetString("admin-accesslog"))),
			}

			if must.Bool(cmd.Flags().GetBool("xds-incremental")) {
				opts = append(opts, bootstrap.EnableIncrementalDiscovery())
			}

			boot, err := bootstrap.New(opts...)
			if err != nil {
				return ExitError{EX_FAIL, err}
			}

			_, err = out.Write(
				must.Bytes(protojson.MarshalOptions{
					Multiline: true,
					Indent:    "  ",
				}.Marshal(boot)),
			)

			return err
		},
	}

	cmd.Flags().String("admin-address", ":8080", "The address the Envoy admin endpoint binds to.")
	cmd.Flags().String("admin-accesslog", "/dev/null", "Path for the Envoy admin endpoint access log.")
	cmd.Flags().String("xds-address", "/var/run/xds.sock", "The address the xDS endpoint binds to.")
	cmd.Flags().String("xds-clustername", "envoy-controller", "The name to use for the xDS management cluster.")
	cmd.Flags().Bool("xds-incremental", false, "Enable the incremental (delta) xDS protocol.")

	cmd.Flags().StringP("filename", "f", "-", "Filename used to create the resource.")
	cmd.Flags().BoolP("3", "3", false, "Bootstrap Envoy to default to the v3 API.")
	cmd.Flags().BoolP("2", "2", false, "Bootstrap Envoy to default to the v2 API.")

	return &cmd
}
