package clustercmd

import (
	"github.com/spf13/cobra"
)

func NewClusterConfigmdGroup() *cobra.Command {
	clusterConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "cluster config operations",
		Long:  `commands to manage your cluster config`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}
	clusterConfigCmd.AddCommand(NewEditConfigCmd())
	clusterConfigCmd.AddCommand(NewSyncClusterConfigCmdGroup())
	clusterConfigCmd.AddCommand(NewApplyConfigCmd())
	clusterConfigCmd.AddCommand(NewPurgeClusterConfigCMD())
	return clusterConfigCmd
}
