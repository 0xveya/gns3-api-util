package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/cmd/add"
	"github.com/stefanistkuhl/gns3util/cmd/delete"
	"github.com/stefanistkuhl/gns3util/cmd/get"
	"github.com/stefanistkuhl/gns3util/cmd/post/create"
	"github.com/stefanistkuhl/gns3util/cmd/put/update"
)

func NewGroupCmdGroup() *cobra.Command {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Group operations",
		Long:  `Create, manage, and manipulate GNS3 groups.`,
	}

	groupCmd.AddCommand(create.NewCreateGroupCmd())

	groupCmd.AddCommand(get.NewGetGroupCmd())
	groupCmd.AddCommand(get.NewGetGroupsCmd())
	groupCmd.AddCommand(get.NewGetGroupMembersCmd())

	groupCmd.AddCommand(update.NewUpdateGroupCmd())

	groupCmd.AddCommand(delete.NewDeleteGroupCmd())

	groupCmd.AddCommand(add.NewAddGroupMemberCmd())

	return groupCmd
}
