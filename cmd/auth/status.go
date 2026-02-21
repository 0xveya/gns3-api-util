package auth

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/pkg/api/schemas"
	"github.com/stefanistkuhl/gns3util/pkg/authentication"
	"github.com/stefanistkuhl/gns3util/pkg/config"
	"github.com/stefanistkuhl/gns3util/pkg/utils/messageUtils"
)

func NewAuthStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check the current status of your Authentication",
		Long:  `Check the current status of your Authentication`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var user schemas.User

			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}

			keys, err := authentication.LoadKeys(cfg.KeyFile)
			if err != nil {
				return fmt.Errorf("failed to load keys: %w", err)
			}

			userData, err := authentication.TryKeys(keys, cfg)
			if err != nil {
				return err
			}

			err = json.Unmarshal([]byte(userData), &user)
			if err != nil {
				return fmt.Errorf("error unmarshaling JSON: %w", err)
			}
			fmt.Printf("%s logged in as user %s", messageUtils.SuccessMsg("Logged in as user"), messageUtils.Bold(*user.Username))
			return nil
		},
	}
	return cmd
}
