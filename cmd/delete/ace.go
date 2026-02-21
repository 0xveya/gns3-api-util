package delete

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/pkg/config"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
)

func NewDeleteACECmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     utils.DeleteSingleElementCmdName + " [ace-id]",
		Short:   "Delete an ACE",
		Long:    `Delete an Access Control Entry (ACE) from the GNS3 server.`,
		Example: "gns3util -s https://controller:3080 acl delete ace-id",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			aceID := args[0]
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}

			if !utils.IsValidUUIDv4(aceID) {
				return fmt.Errorf("ACE ID must be a valid UUID")
			}

			utils.ExecuteAndPrint(cfg, "deleteACE", []string{aceID})
			return nil
		},
	}

	return cmd
}
