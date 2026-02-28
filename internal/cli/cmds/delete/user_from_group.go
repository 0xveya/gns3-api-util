package delete

import (
	"fmt"

	"github.com/0xveya/gns3util/internal/cli/cli_pkg/config"
	"github.com/0xveya/gns3util/internal/cli/cli_pkg/utils"
	"github.com/spf13/cobra"
)

func NewDeleteUserFromGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove-from-group [group-name/id] [user-name/id]",
		Short:   "Remove a user from a group",
		Long:    `Remove a user from a group on the GNS3 server.`,
		Example: "gns3util -s https://controller:3080 group remove-from-group my-group my-user",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID := args[0]
			userID := args[1]
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}

			if !utils.IsValidUUIDv4(groupID) {
				id, err := utils.ResolveID(cfg, "group", groupID, nil)
				if err != nil {
					return err
				}
				groupID = id
			}

			if !utils.IsValidUUIDv4(userID) {
				id, err := utils.ResolveID(cfg, "user", userID, nil)
				if err != nil {
					return err
				}
				userID = id
			}

			utils.ExecuteAndPrint(cfg, "deleteUserFromGroup", []string{groupID, userID})
			return nil
		},
	}

	return cmd
}
