package schemaherokubectlcli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ApproveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "approve",
		Short:         "approves a change to the database, such as a schema migration",
		Long:          `...`,
		Args:          cobra.MinimumNArgs(1),
		SilenceErrors: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(ApproveMigrationCmd())

	return cmd
}
