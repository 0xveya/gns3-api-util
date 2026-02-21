package post

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/pkg/config"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
)

func NewLockProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock-project",
		Short: "Lock a project by id or name",
		Long:  `Lock a project by id or name`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}
			if !utils.IsValidUUIDv4(args[0]) {
				id, err = utils.ResolveID(cfg, "project", args[0], nil)
				if err != nil {
					return err
				}
			}
			utils.ExecuteAndPrint(cfg, "lockProject", []string{id})
			return nil
		},
	}
	return cmd
}
