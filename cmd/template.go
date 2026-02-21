package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/cmd/delete"
	"github.com/stefanistkuhl/gns3util/cmd/get"
	"github.com/stefanistkuhl/gns3util/cmd/post"
	"github.com/stefanistkuhl/gns3util/cmd/post/create"
	"github.com/stefanistkuhl/gns3util/cmd/put/update"
)

func NewTemplateCmdGroup() *cobra.Command {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Template operations",
		Long:  `Create, manage, and manipulate GNS3 templates.`,
	}

	templateCmd.AddCommand(create.NewCreateTemplateCmd())

	templateCmd.AddCommand(get.NewGetTemplatesCmd())
	templateCmd.AddCommand(get.NewGetTemplateCmd())

	templateCmd.AddCommand(post.NewDuplicateTemplateCmd())

	templateCmd.AddCommand(update.NewUpdateTemplateCmd())

	templateCmd.AddCommand(delete.NewDeleteTemplateCmd())

	return templateCmd
}
