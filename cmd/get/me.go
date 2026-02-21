package get

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/pkg/config"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
)

func NewGetMeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "me",
		Short:   "Display info about the currently logged in user on the GNS3 Server",
		Long:    `Display info about the currently logged in user on the GNS3 Server`,
		Example: "gns3util -s https://controller:3080 get me",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}
			utils.ExecuteAndPrint(cfg, "getMe", nil)
			return nil
		},
	}
	return cmd
}
