package managercli

import (
	"os"
	"strings"

	"github.com/schemahero/schemahero/pkg/apis"
	"github.com/schemahero/schemahero/pkg/config"
	databasecontroller "github.com/schemahero/schemahero/pkg/controller/database"
	migrationcontroller "github.com/schemahero/schemahero/pkg/controller/migration"
	tablecontroller "github.com/schemahero/schemahero/pkg/controller/table"
	"github.com/schemahero/schemahero/pkg/logger"
	"github.com/schemahero/schemahero/pkg/version"
	"github.com/schemahero/schemahero/pkg/webhook"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "run",
		Short:         "runs the schemahero manager",
		Long:          `...`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Infof("Starting schemahero version %+v", version.GetBuild())

			v := viper.GetViper()

			if v.GetString("log-level") == "debug" {
				logger.SetDebug()
			}

			// Get a config to talk to the apiserver
			cfg, err := config.GetRESTConfig()
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			// Create a new Cmd to provide shared dependencies and start components
			options := manager.Options{
				MetricsBindAddress: v.GetString("metrics-addr"),
			}

			if v.GetString("namespace") != "" {
				options.Namespace = v.GetString("namespace")
			}

			mgr, err := manager.New(cfg, options)
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			// Setup Scheme for all resources
			if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			if v.GetBool("enable-database-controller") {
				if err := databasecontroller.Add(mgr, v.GetString("manager-image"), v.GetString("manager-tag")); err != nil {
					logger.Error(err)
					os.Exit(1)
				}
			}

			if len(v.GetStringSlice("database-name")) > 0 {
				if err := databasecontroller.AddForDatabaseSchemasOnly(mgr, v.GetStringSlice("database-name")); err != nil {
					logger.Error(err)
					os.Exit(1)
				}

				if err := tablecontroller.Add(mgr, v.GetStringSlice("database-name")); err != nil {
					logger.Error(err)
					os.Exit(1)
				}

				if err := migrationcontroller.Add(mgr, v.GetStringSlice("database-name")); err != nil {
					logger.Error(err)
					os.Exit(1)
				}
			}

			if err := webhook.AddToManager(mgr); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			// Start the Cmd
			if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().String("metrics-addr", ":8088", "The address the metric endpoint binds to.")

	cmd.Flags().Bool("enable-database-controller", false, "when set, the database controller will be active")
	cmd.Flags().StringSlice("database-name", []string{}, "when present (and not set to *), the controller will reconcile tables and migrations for the specified database")
	cmd.Flags().String("manager-image", "schemahero/schemahero-manager", "the schemahero manager image to use in the controller")
	cmd.Flags().String("manager-tag", "0.12.6", "the tag of the schemahero manager image to use")
	cmd.Flags().String("namespace", "", "when set, limit rbac permissions for watches to this namespace")

	return cmd
}

func defaultManagerTag() string {
	tag := version.Version()
	if strings.HasPrefix(tag, "v") {
		tag = strings.TrimPrefix(tag, "v")
	}

	return tag
}
