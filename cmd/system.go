package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/cmd/get"
	"github.com/stefanistkuhl/gns3util/cmd/post"
	"github.com/stefanistkuhl/gns3util/cmd/put/update"
)

func NewSystemCmdGroup() *cobra.Command {
	systemCmd := &cobra.Command{
		Use:   "system",
		Short: "System operations",
		Long:  `Manage GNS3 system operations and settings.`,
	}

	systemCmd.AddCommand(get.NewGetVersionCmd())
	systemCmd.AddCommand(get.NewGetStatisticsCmd())
	systemCmd.AddCommand(get.NewGetNotificationsCmd())
	systemCmd.AddCommand(get.NewGetIouLicenseCmd())

	systemCmd.AddCommand(post.NewCheckVersionCmd())
	systemCmd.AddCommand(post.NewControllerCmdGroup())

	systemCmd.AddCommand(update.NewUpdateIOULicenseCmd())

	return systemCmd
}
