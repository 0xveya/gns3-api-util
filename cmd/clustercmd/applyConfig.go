package clustercmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/pkg/cluster"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
	"github.com/stefanistkuhl/gns3util/pkg/utils/messageUtils"
)

func NewApplyConfigCmd() *cobra.Command {
	var noConfirm bool
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply your config file to the local database",
		Long:  `apply your config file to the local database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgLoaded, err := cluster.LoadClusterConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if !noConfirm {
				if !utils.ConfirmPrompt(fmt.Sprintf("%s do you want to apply this config to the Database?", messageUtils.WarningMsg("Warning")), false) {
					return nil
				}
			}

			applyErr := cluster.ApplyConfig(cfgLoaded)
			if applyErr != nil {
				return fmt.Errorf("error applying config: %w", applyErr)
			}
			fmt.Printf("%s applied config to the Database.\n", messageUtils.SuccessMsg("Success"))
			return nil
		},
	}
	cmd.Flags().BoolVarP(&noConfirm, "no-confirm", "n", false, "Run the command without asking for confirmations.")

	return cmd
}
