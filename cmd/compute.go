package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/cmd/delete"
	"github.com/stefanistkuhl/gns3util/cmd/get"
	"github.com/stefanistkuhl/gns3util/cmd/post/create"
	"github.com/stefanistkuhl/gns3util/cmd/put/update"
)

func NewComputeCmdGroup() *cobra.Command {
	computeCmd := &cobra.Command{
		Use:   "compute",
		Short: "Compute operations",
		Long:  `Create, manage, and manipulate GNS3 computes.`,
	}

	computeCmd.AddCommand(create.NewCreateComputeCmd())

	computeCmd.AddCommand(get.NewGetComputesCmd())
	computeCmd.AddCommand(get.NewGetComputeCmd())
	computeCmd.AddCommand(get.NewGetComputeDockerImagesCmd())
	computeCmd.AddCommand(get.NewGetComputeVirtualboxVMSCmd())
	computeCmd.AddCommand(get.NewGetComputeVmWareVMSCmd())

	computeCmd.AddCommand(update.NewUpdateComputeCmd())

	computeCmd.AddCommand(delete.NewDeleteComputeCmd())

	return computeCmd
}
