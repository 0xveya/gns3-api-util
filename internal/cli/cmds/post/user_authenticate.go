package post

import (
	"fmt"

	"github.com/0xveya/gns3util/internal/cli/cli_pkg/config"
	"github.com/0xveya/gns3util/internal/cli/cli_pkg/utils"
	"github.com/spf13/cobra"
)

func NewUserAuthenticateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user-authenticate",
		Short:   "Authenticate as a user",
		Long:    `Authenticate as a user on the GNS3 server.`,
		Example: `gns3util -s https://controller:3080 post user-authenticate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}

			utils.ExecuteAndPrint(cfg, "userAuthenticate", nil)
			return nil
		},
	}

	return cmd
}
