package get

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/pkg/config"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
)

func NewGetAclCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     utils.ListAllCmdName,
		Short:   "Get the acl-rules of the GNS3 Server",
		Long:    `Get the acl-rules of the GNS3 Server`,
		Example: "gns3util -s https://controller:3080 acl ls",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}
			utils.ExecuteAndPrint(cfg, "getAcl", nil)
			return nil
		},
	}
	return cmd
}

func NewGetAceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     utils.ListSingleElementCmdName + " [ace-id]",
		Short:   "Get an ace by id",
		Long:    `Get an ace by id`,
		Example: "gns3util -s https://controller:3080 acl info ace-id",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}
			utils.ExecuteAndPrint(cfg, "getAce", []string{id})
			return nil
		},
	}
	return cmd
}

func NewGetAclEndpointsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "endpoints",
		Short:   "Get the available endpoints for acl-rules",
		Long:    `Get the available endpoints for acl-rules`,
		Example: "gns3util -s https://controller:3080 acl endpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}
			utils.ExecuteAndPrint(cfg, "getAclEndpoints", nil)
			return nil
		},
	}
	return cmd
}
