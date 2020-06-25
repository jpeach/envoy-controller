package cli

import (
	"github.com/jpeach/envoy-controller/controllers"
	"github.com/jpeach/envoy-controller/pkg/kubernetes"
	"github.com/jpeach/envoy-controller/pkg/must"
	"github.com/jpeach/envoy-controller/pkg/util"
	"github.com/jpeach/envoy-controller/pkg/xds"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
)

// NewRunCommand returns a command that runs the controller.
func NewRunCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "run [OPTIONS]",
		Short: "Run the Envoy controller",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			xdsServer := xds.NewServer(grpc.MaxConcurrentStreams(1 << 20))
			xdsListener, err := util.NewListener(must.String(cmd.Flags().GetString("xds-address")))
			if err != nil {
				return ExitErrorf(EX_CONFIG, "invalid xDS listener address %q: %w",
					must.String(cmd.Flags().GetString("xds-address")), err)
			}

			mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
				Scheme:             kubernetes.NewScheme(),
				MetricsBindAddress: must.String(cmd.Flags().GetString("metrics-address")),
				LeaderElection:     must.Bool(cmd.Flags().GetBool("enable-leader-election")),
				LeaderElectionID:   "06187118.projectcontour.io",
			})
			if err != nil {
				return ExitErrorf(EX_FAIL, "unable to start manager: %w", err)
			}

			for _, k := range xds.Kinds() {
				r := controllers.New(k, mgr.GetClient(), mgr.GetScheme())
				if err := r.SetupWithManager(mgr); err != nil {
					return ExitErrorf(EX_FAIL, "unable to create %q reconciler: %w", k, err)
				}
			}

			errChan := make(chan error)
			stopChan := ctrl.SetupSignalHandler()

			go func() {
				if err := xdsServer.Start(xdsListener, stopChan); err != nil {
					errChan <- ExitErrorf(EX_FAIL, "xDS server failed: %w", err)
				}

				errChan <- nil
			}()

			go func() {
				if err := mgr.Start(stopChan); err != nil {
					errChan <- ExitErrorf(EX_FAIL, "controller manager failed: %w", err)
				}

				errChan <- nil
			}()

			select {
			case err := <-errChan:
				return err
			case <-stopChan:
				return nil
			}
		},
	}

	cmd.Flags().String("metrics-address", ":8080", "The address the metric endpoint binds to.")
	cmd.Flags().String("xds-address", ":8080", "The address the xDS endpoint binds to.")
	cmd.Flags().Bool("enable-leader-election", false,
		"Enable leader election to ensure there is only one active controller.")

	return &cmd
}
