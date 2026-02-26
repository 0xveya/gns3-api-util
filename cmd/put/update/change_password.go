package update

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stefanistkuhl/gns3util/pkg/api/schemas"
	"github.com/stefanistkuhl/gns3util/pkg/config"
	"github.com/stefanistkuhl/gns3util/pkg/fuzzy"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
	"github.com/stefanistkuhl/gns3util/pkg/utils/messageUtils"
	"github.com/tidwall/gjson"
)

var (
	passwordFlag string
	useFuzzy     bool
)

func NewChangePasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-password [user-name/id]",
		Short: "Change a user's password",
		Long:  `Change a user's password with interactive password input and validation`,
		Example: `gns3util -s https://controller:3080 user change-password my-user
gns3util -s https://controller:3080 user change-password -f
gns3util -s https://controller:3080 user change-password my-user -p "newpassword123"`,
		Args: func(cmd *cobra.Command, args []string) error {
			if useFuzzy {
				if len(args) > 1 {
					return fmt.Errorf("at most 1 positional arg allowed when --fuzzy is set")
				}
				return nil
			}
			if len(args) != 1 {
				return fmt.Errorf("requires 1 arg [user-name/id] when --fuzzy is not set")
			}
			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.SetEnvPrefix("GNS3")
			viper.AutomaticEnv()

			_ = viper.BindPFlag("password", cmd.Flags().Lookup("password"))

			if !cmd.Flags().Changed("password") {
				passwordFlag = viper.GetString("password")
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}

			var userID string
			var username string

			if useFuzzy {
				var rawData []byte
				rawData, _, err = utils.CallClient(cfg, "getUsers", nil, nil)
				if err != nil {
					return fmt.Errorf("error getting users: %w", err)
				}

				result := gjson.ParseBytes(rawData)
				if !result.IsArray() {
					return fmt.Errorf("expected array response, got %s", result.Type)
				}

				var apiData []gjson.Result
				var usernames []string

				result.ForEach(func(_, value gjson.Result) bool {
					apiData = append(apiData, value)
					if val := value.Get("username"); val.Exists() {
						usernames = append(usernames, val.String())
					}
					return true
				})

				if len(usernames) == 0 {
					return fmt.Errorf("no users found")
				}

				selected := fuzzy.NewFuzzyFinder(usernames, false)
				if len(selected) == 0 {
					return fmt.Errorf("no user selected")
				}

				for _, data := range apiData {
					if usernameField := data.Get("username"); usernameField.Exists() && usernameField.String() == selected[0] {
						userID = data.Get("user_id").String()
						username = selected[0]
						break
					}
				}
			} else {
				userID = args[0]
				username = args[0]

				if !utils.IsValidUUIDv4(args[0]) {
					userID, err = utils.ResolveID(cfg, "user", args[0], nil)
					if err != nil {
						return err
					}
				}
			}

			var newPassword string

			if passwordFlag != "" {
				if !utils.ValidatePassword(passwordFlag) {
					return fmt.Errorf("password must be at least 8 characters with at least 1 number and 1 lowercase letter")
				}
				newPassword = passwordFlag
			} else {
				fmt.Printf("Changing password for user: %s\n", messageUtils.Bold(username))
				newPassword, err = utils.GetPasswordFromInput()
				if err != nil {
					return fmt.Errorf("failed to get password: %w", err)
				}
			}

			userUpdate := schemas.UserUpdate{
				Password: &newPassword,
			}

			data, err := json.Marshal(userUpdate)
			if err != nil {
				return fmt.Errorf("failed to marshal user update: %w", err)
			}

			var payload map[string]any
			if err := json.Unmarshal(data, &payload); err != nil {
				return fmt.Errorf("failed to prepare payload: %w", err)
			}

			utils.ExecuteAndPrintWithBody(cfg, "updateUser", []string{userID}, payload)
			return nil
		},
	}

	cmd.Flags().StringVarP(&passwordFlag, "password", "p", "", "New password (env: GNS3_PASSWORD)")
	cmd.Flags().BoolVarP(&useFuzzy, "fuzzy", "f", false, "Use fuzzy search to select a user")

	return cmd
}
