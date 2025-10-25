package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stefanistkuhl/gns3util/cmd/auth"
	"github.com/stefanistkuhl/gns3util/cmd/class"
	"github.com/stefanistkuhl/gns3util/cmd/exercise"
	"github.com/stefanistkuhl/gns3util/pkg/config"
)

var (
	server   string
	keyFile  string
	insecure bool
	raw      bool
	noColor  bool
	version  bool
)

var Version = "1.2.8"

var Foo bool

var rootCmd = &cobra.Command{
	Use:           "gns3util",
	Short:         "A utility for GNS3v3",
	Long:          `A utility for GNS3v3 for managing GNS3v3 projects and devices.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		viper.SetEnvPrefix("GNS3")
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.AutomaticEnv()

		_ = viper.BindPFlag("server", cmd.Flags().Lookup("server"))
		_ = viper.BindPFlag("insecure", cmd.Flags().Lookup("insecure"))
		_ = viper.BindPFlag("raw", cmd.Flags().Lookup("raw"))
		_ = viper.BindPFlag("no-color", cmd.Flags().Lookup("no-color"))
		_ = viper.BindPFlag("key-file", cmd.Flags().Lookup("key-file"))

		insecure := insecure || viper.GetBool("insecure")
		raw := raw || viper.GetBool("raw")
		noColor := noColor || viper.GetBool("no-color")
		keyFile = viper.GetString("key-file")

		if cmd.Name() == "completion" || (cmd.Parent() != nil && cmd.Parent().Name() == "completion") {
			return nil
		}

		if version {
			return nil
		}

		if err := validateGlobalFlags(); err != nil {
			return err
		}

		skipServer := false
		if f := cmd.Flags().Lookup("cluster"); f != nil {
			if v, _ := cmd.Flags().GetString("cluster"); v != "" {
				skipServer = true
			}
		}
		if !skipServer {
			if f := cmd.InheritedFlags().Lookup("cluster"); f != nil {
				if v, _ := cmd.InheritedFlags().GetString("cluster"); v != "" {
					skipServer = true
				}
			}
		}
		if !skipServer {
			if err := validateRequiresServer(); err != nil {
				return err
			}
		}

		opts := config.GlobalOptions{
			Server:   server,
			Insecure: insecure,
			KeyFile:  keyFile,
			Raw:      raw,
			NoColors: noColor,
		}
		ctx := config.WithGlobalOptions(cmd.Context(), opts)
		cmd.SetContext(ctx)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Printf("gns3util version %s\n", Version)
			os.Exit(0)
		}
		_ = cmd.Help()
	},
}

func init() {
	cobra.OnFinalize()
	rootCmd.PersistentFlags().StringVarP(&server, "server", "s", "", "GNS3v3 Server URL (required for non cluster commands) ENV: GNS3_SERVER")
	rootCmd.PersistentFlags().StringVarP(&keyFile, "key-file", "k", "", "Set a location for a keyfile to use")
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "i", false, "Ignore unsigned SSL-Certificates")
	rootCmd.PersistentFlags().BoolVarP(&raw, "raw", "", false, "Output all data in raw json")
	rootCmd.PersistentFlags().BoolVarP(&noColor, "no-color", "", false, "Output all data in raw json and dont use a colored output")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information")

	rootCmd.AddCommand(auth.NewAuthCmdGroup())

	rootCmd.AddCommand(class.NewClassCmdGroup())
	rootCmd.AddCommand(exercise.NewExerciseCmdGroup())

	rootCmd.AddCommand(NewProjectCmdGroup())
	rootCmd.AddCommand(NewNodeCmdGroup())
	rootCmd.AddCommand(NewLinkCmdGroup())
	rootCmd.AddCommand(NewDrawingCmdGroup())
	rootCmd.AddCommand(NewTemplateCmdGroup())
	rootCmd.AddCommand(NewComputeCmdGroup())
	rootCmd.AddCommand(NewApplianceCmdGroup())
	rootCmd.AddCommand(NewImageCmdGroup())
	rootCmd.AddCommand(NewSymbolCmdGroup())

	rootCmd.AddCommand(NewUserCmdGroup())
	rootCmd.AddCommand(NewGroupCmdGroup())
	rootCmd.AddCommand(NewRoleCmdGroup())
	rootCmd.AddCommand(NewAclCmdGroup())

	rootCmd.AddCommand(NewPoolCmdGroup())
	rootCmd.AddCommand(NewSnapshotCmdGroup())

	rootCmd.AddCommand(NewSystemCmdGroup())

	rootCmd.AddCommand(NewRemoteCmdGroup())

	rootCmd.AddCommand(NewClusterCmdGroup())
	rootCmd.AddCommand(NewShareCmdGroup())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func validateGlobalFlags() error {
	if noColor && !raw {
		return fmt.Errorf("--no-color can only be used when --raw is also used")
	}
	return nil
}

func validateRequiresServer() error {
	if server == "" {
		server = viper.GetString("server")
	}
	if server == "" {
		return fmt.Errorf("required flag(s) \"server\" not set")
	}
	return nil
}
