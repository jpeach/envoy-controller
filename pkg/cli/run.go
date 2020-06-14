package cli

import (
	envoyv1alpha1 "github.com/jpeach/envoy-controller/api/v1alpha1"
	"github.com/jpeach/envoy-controller/controllers"
	"github.com/jpeach/envoy-controller/pkg/must"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

// NewRunCommand returns a command that runs the controller.
func NewRunCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "run [OPTIONS]",
		Short: "Run the Envoy controller",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var kinds = []string{
				"Listener",
				"Cluster",
				"RouteConfiguration",
				"ScopedRouteConfiguration",
				"Secret",
				"Runtime",
				"VirtualHost",
			}

			s := runtime.NewScheme()

			must.Must(scheme.AddToScheme(s))
			must.Must(envoyv1alpha1.AddToScheme(s))

			mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
				Scheme:             s,
				MetricsBindAddress: must.String(cmd.Flags().GetString("metrics-address")),
				LeaderElection:     must.Bool(cmd.Flags().GetBool("enable-leader-election")),
				LeaderElectionID:   "06187118.projectcontour.io",
			})
			if err != nil {
				return ExitErrorf(EX_FAIL, "unable to start manager: %w", err)
			}

			for _, k := range kinds {
				r := controllers.New(k, mgr.GetClient(), mgr.GetScheme())
				if err := r.SetupWithManager(mgr); err != nil {
					return ExitErrorf(EX_FAIL, "unable to create %q reconciler: %w", k, err)
				}
			}

			if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
				return ExitErrorf(EX_FAIL, "problem running manager: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().String("metrics-address", ":8080", "The address the metric endpoint binds to.")
	cmd.Flags().Bool("enable-leader-election", false,
		"Enable leader election to ensure there is only one active controller.")

	return &cmd
}
