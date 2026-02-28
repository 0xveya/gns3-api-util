package clustercmd

import (
	"fmt"

	"github.com/0xveya/gns3util/internal/cli/cli_pkg/cluster"
	"github.com/0xveya/gns3util/internal/cli/cli_pkg/utils"
	"github.com/0xveya/gns3util/internal/cli/cli_pkg/utils/messageUtils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func NewEditConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "edit your configuration with your $EDITOR",
		Long:  `edit your configuration with your $EDITOR`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgLoaded, err := cluster.LoadClusterConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			res, marshallErr := toml.Marshal(&cfgLoaded)
			if marshallErr != nil {
				return fmt.Errorf("failed to marshal config: %w", marshallErr)
			}
			str, editErr := utils.EditTextWithEditor(string(res), "toml")
			if editErr != nil {
				return fmt.Errorf("failed to edit config: %w", editErr)
			}
			var cfgNew cluster.Config
			unMarshallErr := toml.Unmarshal([]byte(str), &cfgNew)
			if unMarshallErr != nil {
				return fmt.Errorf("failed to unmarshal config: %w", unMarshallErr)
			}
			writeErr := cluster.WriteClusterConfig(cfgNew)
			if writeErr != nil {
				return fmt.Errorf("failed to write edited config to the config file: %w", writeErr)
			}
			fmt.Printf("%s wrote new config to the config file.", messageUtils.SuccessMsg("Success"))
			return nil
		},
	}

	return cmd
}
