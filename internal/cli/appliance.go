package cli

import (
	"github.com/0xveya/gns3util/internal/cli/cmds/get"
	"github.com/spf13/cobra"
)

func NewApplianceCmdGroup() *cobra.Command {
	applianceCmd := &cobra.Command{
		Use:   "appliance",
		Short: "Appliance operations",
		Long:  `Get and manage GNS3 appliances.`,
	}

	applianceCmd.AddCommand(get.NewGetAppliancesCmd())
	applianceCmd.AddCommand(get.NewGetApplianceCmd())

	return applianceCmd
}
