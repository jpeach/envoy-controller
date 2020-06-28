/*
Copyright 2020 VMware, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jpeach/envoy-controller/pkg/cli"
	"github.com/jpeach/envoy-controller/pkg/must"
	"github.com/jpeach/envoy-controller/pkg/util"
	"github.com/jpeach/envoy-controller/pkg/version"

	// Ensure that xDS protobuf types are always registered.
	_ "github.com/jpeach/envoy-controller/pkg/xds"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	// Initialize all known client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	root := cli.Defaults(&cobra.Command{
		Use:     version.Progname,
		Short:   "Kubernetes controller for the Envoy proxy",
		Version: fmt.Sprintf("%s/%s, built %s", version.Version, version.Sha, version.BuildDate),
	})

	// Initialize logging in the pre-run so that we can parse
	// the debug flag globally. Presumably, no sub-command will
	// ever use a pre-run.
	root.PersistentPreRun = func(*cobra.Command, []string) {
		opts := []zap.Opts{
			zap.UseDevMode(isatty.IsTerminal(os.Stdout.Fd())),
		}

		// We have to use root here, because the local cmd
		// var is a the Command instance for the subcommand.
		if must.Bool(root.PersistentFlags().GetBool("debug")) {
			opts = append(opts, zap.Level(util.ZapEnableDebug()))
		}

		ctrl.SetLogger(zap.New(opts...))
	}

	root.PersistentFlags().Bool("debug", false, "Enable debug logging and behavior.")

	root.AddCommand(cli.Defaults(cli.NewRunCommand()))
	root.AddCommand(cli.Defaults(cli.NewCreateCommand()))
	root.AddCommand(cli.Defaults(cli.NewBootstrapCommand()))

	if err := root.Execute(); err != nil {
		if msg := err.Error(); msg != "" {
			fmt.Fprintf(os.Stderr, "%s: %s\n", version.Progname, msg)
		}

		var exit *cli.ExitError
		if errors.As(err, &exit) {
			os.Exit(int(exit.Code))
		}

		os.Exit(int(cli.EX_FAIL))
	}
}
