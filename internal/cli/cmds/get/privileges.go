package get

import (
	"fmt"

	"github.com/0xveya/gns3util/internal/cli/cli_pkg/config"
	"github.com/0xveya/gns3util/internal/cli/cli_pkg/utils"
	"github.com/spf13/cobra"
)

func NewGetPrivilegesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "privileges-all",
		Short:   "Get all privileges",
		Long:    `Get all privileges from the GNS3 server.`,
		Example: `gns3util -s https://controller:3080 role privileges-all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}

			utils.ExecuteAndPrint(cfg, "getPrivileges", nil)
			return nil
		},
	}

	return cmd
}
