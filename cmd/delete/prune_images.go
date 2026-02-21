package delete

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/pkg/config"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
)

func NewDeletePruneImagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "prune",
		Short:   "Prune unused images",
		Long:    `Delete unused images from the GNS3 server to free up disk space.`,
		Example: `gns3util -s https://controller:3080 image prune`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}

			utils.ExecuteAndPrint(cfg, "deletePruneImages", nil)
			return nil
		},
	}

	return cmd
}
