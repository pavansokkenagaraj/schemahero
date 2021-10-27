package schemaherokubectlcli

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	schemasclientv1alpha4 "github.com/schemahero/schemahero/pkg/client/schemaheroclientset/typed/schemas/v1alpha4"
	"github.com/schemahero/schemahero/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	kuberneteserrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ApproveMigrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "migration",
		Short:         "approve a migration, which will allow changes to be automatically be applied to the database",
		Long:          `approve a migration, which will allow changes to be automatically be applied to the database`,
		Args:          cobra.MinimumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()
			ctx := context.Background()
			migrationName := args[0]

			cfg, err := config.GetRESTConfig()
			if err != nil {
				return err
			}

			client, err := kubernetes.NewForConfig(cfg)
			if err != nil {
				return err
			}

			schemasClient, err := schemasclientv1alpha4.NewForConfig(cfg)
			if err != nil {
				return err
			}

			namespaceNames := []string{}

			if viper.GetBool("all-namespaces") {
				namespaces, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
				if err != nil {
					return err
				}

				for _, namespace := range namespaces.Items {
					namespaceNames = append(namespaceNames, namespace.Name)
				}
			} else {
				if v.GetString("namespace") != "" {
					namespaceNames = []string{v.GetString("namespace")}
				} else {
					namespaceNames = []string{"default"}
				}
			}

			for _, namespaceName := range namespaceNames {
				migration, err := schemasClient.Migrations(namespaceName).Get(ctx, migrationName, metav1.GetOptions{})
				if kuberneteserrors.IsNotFound(err) {
					// continue to the next namespace
					continue
				}
				if err != nil {
					return err
				}

				// rejected migrations cannot be approved
				if migration.Status.RejectedAt > 0 {
					return errors.Errorf("Rejected migration %s cannot be approved\n", migrationName)
				}

				migration.Status.ApprovedAt = time.Now().Unix()
				if _, err := schemasClient.Migrations(namespaceName).Update(ctx, migration, metav1.UpdateOptions{}); err != nil {
					return err
				}

				fmt.Printf("Migration %s approved\n", migrationName)
				return nil
			}

			err = errors.Errorf("migration %q not found", migrationName)
			return err
		},
	}

	cmd.Flags().Bool("all-namespaces", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")

	return cmd
}
