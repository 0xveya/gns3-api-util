package clustercmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/pkg/cluster"
	"github.com/stefanistkuhl/gns3util/pkg/cluster/db"
	"github.com/stefanistkuhl/gns3util/pkg/cluster/db/sqlc"
	"github.com/stefanistkuhl/gns3util/pkg/utils/messageUtils"
)

var (
	name string
	desc string
)

func NewCreateClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a cluster",
		Long:  `create a cluster`,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := db.Init()
			if err != nil {
				return fmt.Errorf("failed to init db: %w", err)
			}

			ctx := context.Background()
			exists, err := store.CheckIfClusterExists(ctx, name)
			if err != nil {
				return fmt.Errorf("failed to check cluster existence: %w", err)
			}
			if exists == 1 {
				return fmt.Errorf("a cluster with the name %s already exists", name)
			}
			var insertDesc sql.NullString
			if desc != "" {
				insertDesc = sql.NullString{String: desc, Valid: true}
			} else {
				insertDesc = sql.NullString{String: "", Valid: false}
			}
			_, err = store.CreateCluster(ctx, sqlc.CreateClusterParams{Name: name, Description: insertDesc})
			if err != nil {
				return fmt.Errorf("failed to create cluster: %w", err)
			}

			cfg, cfgErr := cluster.LoadClusterConfig()
			if cfgErr != nil {
				if errors.Is(cfgErr, cluster.ErrNoConfig) {
					cfg = cluster.NewConfig()
				} else {
					return fmt.Errorf("failed to load config: %w", cfgErr)
				}
			}

			cfg.Clusters = append(cfg.Clusters, cluster.Cluster{
				Name:        name,
				Description: desc,
			})
			if writeErr := cluster.WriteClusterConfig(cfg); writeErr != nil {
				return fmt.Errorf("failed to write to config file: %w", writeErr)
			}

			fmt.Printf("%s created new empty cluster %s\n", messageUtils.SuccessMsg("Success"), name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "name for the cluster")
	cmd.Flags().StringVarP(&desc, "description", "d", "", "description for the cluster")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	return cmd
}
