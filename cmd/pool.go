package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/cmd/add"
	"github.com/stefanistkuhl/gns3util/cmd/delete"
	"github.com/stefanistkuhl/gns3util/cmd/get"
	"github.com/stefanistkuhl/gns3util/cmd/post/create"
	"github.com/stefanistkuhl/gns3util/cmd/put/update"
)

func NewPoolCmdGroup() *cobra.Command {
	poolCmd := &cobra.Command{
		Use:   "pool",
		Short: "Resource pool operations",
		Long:  `Create, manage, and manipulate GNS3 resource pools.`,
	}

	poolCmd.AddCommand(create.NewCreatePoolCmd())

	poolCmd.AddCommand(get.NewGetPoolsCmd())
	poolCmd.AddCommand(get.NewGetPoolCmd())
	poolCmd.AddCommand(get.NewGetPoolResourcesCmd())

	poolCmd.AddCommand(update.NewUpdatePoolCmd())

	poolCmd.AddCommand(delete.NewDeletePoolCmd())
	poolCmd.AddCommand(delete.NewDeletePoolResourceCmd())

	poolCmd.AddCommand(add.NewAddToPoolCmd())

	return poolCmd
}
