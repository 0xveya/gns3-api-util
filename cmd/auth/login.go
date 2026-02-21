package auth

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stefanistkuhl/gns3util/pkg/api/schemas"
	"github.com/stefanistkuhl/gns3util/pkg/authentication"
	"github.com/stefanistkuhl/gns3util/pkg/config"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
	"github.com/stefanistkuhl/gns3util/pkg/utils/messageUtils"
)

var (
	username string
	password string
)

func NewAuthLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in as user",
		Long:  `Log in as a user`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.SetEnvPrefix("GNS3")
			viper.AutomaticEnv()

			_ = viper.BindPFlag("user", cmd.Flags().Lookup("user"))
			_ = viper.BindPFlag("password", cmd.Flags().Lookup("password"))

			if !cmd.Flags().Changed("user") {
				username = viper.GetString("user")
			}
			if !cmd.Flags().Changed("password") {
				password = viper.GetString("password")
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}

			if username == "" || password == "" {
				interactiveUsername, interactivePassword, err := utils.GetLoginCredentials()
				if err != nil {
					return err
				}

				if username == "" {
					username = interactiveUsername
				}
				if password == "" {
					password = interactivePassword
				}
			}

			if username == "" || password == "" {
				return fmt.Errorf("username and password are required")
			}

			credentials := schemas.Credentials{
				Username: username,
				Password: password,
			}

			data, err := json.Marshal(credentials)
			if err != nil {
				return fmt.Errorf("failed to marshal credentials: %w", err)
			}

			var payload map[string]any
			if err := json.Unmarshal(data, &payload); err != nil {
				return fmt.Errorf("failed to prepare payload: %w", err)
			}
			body, status, err := utils.CallClient(cfg, "userAuthenticate", []string{}, payload)
			if err != nil {
				if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "Authentication was unsuccessful") {
					return fmt.Errorf("authentication failed. Please check your username and password")
				}
				return err
			}

			if status == 200 {
				fmt.Printf("%v Successfully logged in as %s\n", messageUtils.SuccessMsg("Success"), messageUtils.Bold(username))
				var token schemas.Token
				marshallErr := json.Unmarshal(body, &token)
				if marshallErr != nil {
					return fmt.Errorf("failed to unmarshall response: %w", marshallErr)
				}
				writeErr := authentication.SaveAuthData(cfg, token, credentials.Username)
				if writeErr != nil {
					return fmt.Errorf("failed to write authentication data to the keyfile: %w", writeErr)
				}
			} else {
				return fmt.Errorf("authentication failed (status: %d)", status)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&username, "user", "u", "", "User to log in as (env: GNS3_USER)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password to use (env: GNS3_PASSWORD)")

	return cmd
}
