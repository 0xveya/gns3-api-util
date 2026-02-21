package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/cmd/delete"
	"github.com/stefanistkuhl/gns3util/cmd/get"
	"github.com/stefanistkuhl/gns3util/cmd/post/create"
	"github.com/stefanistkuhl/gns3util/cmd/put/update"
)

func NewDrawingCmdGroup() *cobra.Command {
	drawingCmd := &cobra.Command{
		Use:   "drawing",
		Short: "Drawing operations",
		Long:  `Create, manage, and manipulate GNS3 drawings.`,
	}

	drawingCmd.AddCommand(create.NewCreateDrawingCmd())

	drawingCmd.AddCommand(get.NewGetDrawingsCmd())
	drawingCmd.AddCommand(get.NewGetDrawingCmd())

	drawingCmd.AddCommand(update.NewUpdateDrawingCmd())

	drawingCmd.AddCommand(delete.NewDeleteDrawingCmd())

	return drawingCmd
}
