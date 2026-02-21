package get

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/pkg/config"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
)

func NewGetStatisticsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "statistics",
		Short:   "Get the statistics of the GNS3 Server",
		Long:    `Get the statistics of the GNS3 Server`,
		Example: "gns3util -s https://controller:3080 get statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}
			utils.ExecuteAndPrint(cfg, "getStatistics", nil)
			return nil
		},
	}
	return cmd
}
