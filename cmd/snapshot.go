package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/cmd/delete"
	"github.com/stefanistkuhl/gns3util/cmd/get"
	"github.com/stefanistkuhl/gns3util/cmd/post/create"
)

func NewSnapshotCmdGroup() *cobra.Command {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Snapshot operations",
		Long:  `Create, manage, and manipulate GNS3 snapshots.`,
	}

	snapshotCmd.AddCommand(create.NewCreateSnapshotCmd())

	snapshotCmd.AddCommand(get.NewGetSnapshotsCmd())

	snapshotCmd.AddCommand(delete.NewDeleteSnapshotCmd())

	return snapshotCmd
}
