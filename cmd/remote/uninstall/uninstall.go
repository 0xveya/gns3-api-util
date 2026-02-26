package uninstall

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stefanistkuhl/gns3util/pkg/config"
	"github.com/stefanistkuhl/gns3util/pkg/ssh"
	"github.com/stefanistkuhl/gns3util/pkg/utils"
	"github.com/stefanistkuhl/gns3util/pkg/utils/gns3"
	"github.com/stefanistkuhl/gns3util/pkg/utils/messageUtils"
	"github.com/stefanistkuhl/gns3util/pkg/utils/ssl"
)

func NewUninstallCmdGroup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall services from remote GNS3 servers",
		Long:  "Uninstall various services and configurations from remote GNS3 servers",
	}

	// Add subcommands
	cmd.AddCommand(NewUninstallHTTPSCmd())
	cmd.AddCommand(NewUninstallGNS3Cmd())

	return cmd
}

func NewUninstallHTTPSCmd() *cobra.Command {
	var (
		port             int
		privateKeyPath   string
		reverseProxyPort int
		gns3Port         int
		domain           string
		subject          string
		firewallAllow    string
		firewallBlock    bool
		verbose          bool
		interactive      bool
	)

	cmd := &cobra.Command{
		Use:   "https [user]",
		Short: "Uninstall SSL reverse proxy setup from remote GNS3 server",
		Long: `Uninstall the SSL reverse proxy setup from a remote GNS3 server.

This command will:
- Stop and remove Caddy reverse proxy
- Remove SSL certificates
- Remove firewall rules
- Clean up systemd services
- Remove configuration files

The GNS3 server itself will remain running on its original port.

If a state file is found from a previous installation, all configuration values
will be automatically loaded and command line flags will be ignored.`,
		Args: cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if interactive {
				fmt.Println("=== Interactive SSL Uninstallation Setup ===")
				fmt.Println("Press Enter to use default values (shown in brackets)")
				fmt.Println()

				fmt.Printf("Reverse proxy port to uninstall [%d]: ", reverseProxyPort)
				var input string
				_, _ = fmt.Scanln(&input)
				if input != "" {
					if p, err := strconv.Atoi(input); err == nil {
						reverseProxyPort = p
					}
				}

				fmt.Printf("GNS3 server port [%d]: ", gns3Port)
				_, _ = fmt.Scanln(&input)
				if input != "" {
					if p, err := strconv.Atoi(input); err == nil {
						gns3Port = p
					}
				}

				fmt.Printf("Domain that was used (leave empty if none) [%s]: ", domain)
				_, _ = fmt.Scanln(&input)
				if input != "" {
					domain = input
				}

				fmt.Printf("SSL certificate subject that was used [%s]: ", subject)
				_, _ = fmt.Scanln(&input)
				if input != "" {
					subject = input
				}

				fmt.Printf("Firewall allow subnet that was used (leave empty if none) [%s]: ", firewallAllow)
				_, _ = fmt.Scanln(&input)
				if input != "" {
					firewallAllow = input
				}

				fmt.Printf("Were firewall rules configured? (y/N): ")
				_, _ = fmt.Scanln(&input)
				firewallBlock = strings.EqualFold(input, "y") || strings.EqualFold(input, "yes")

				fmt.Println()
				fmt.Println("=== Uninstallation Configuration Summary ===")
				fmt.Printf("Reverse proxy port: %d\n", reverseProxyPort)
				fmt.Printf("GNS3 server port: %d\n", gns3Port)
				fmt.Printf("Domain: %s\n", domain)
				fmt.Printf("Subject: %s\n", subject)
				fmt.Printf("Firewall allow: %s\n", firewallAllow)
				fmt.Printf("Firewall block: %t\n", firewallBlock)
				fmt.Println()
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			user := args[0]

			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}

			hostname, sshPort := ssl.ParseServerURLForSSH(cfg.Server, port)

			stateManager, err := ssl.NewStateManager()
			var state *ssl.ServerState
			if err == nil {
				if localState, stateErr := stateManager.LoadState(hostname); stateErr == nil {
					state = localState
					fmt.Printf("%s Loaded state from local machine\n", messageUtils.SuccessMsg("Loaded state from local machine"))
				}
			}

			sslArgs := ssl.InstallSSLArgs{
				FirewallAllow:    firewallAllow,
				FirewallBlock:    firewallBlock,
				ReverseProxyPort: reverseProxyPort,
				Domain:           domain,
				GNS3Port:         gns3Port,
				Subject:          subject,
				Verbose:          verbose,
			}

			if state != nil {
				sslArgs.FirewallAllow = state.FirewallAllow
				sslArgs.FirewallBlock = state.FirewallBlock
				sslArgs.ReverseProxyPort = state.ReverseProxyPort
				sslArgs.Domain = state.Domain
				sslArgs.GNS3Port = state.GNS3Port
				sslArgs.Subject = "/CN=localhost"
				fmt.Printf("%s Using saved configuration: RP=%d, GNS3=%d, Firewall=%t\n",
					messageUtils.InfoMsg("Using saved configuration"), state.ReverseProxyPort, state.GNS3Port, state.FirewallBlock)
			} else {
				fmt.Printf("%s No state found, using command line flags or defaults\n", messageUtils.WarningMsg("No state found"))
				if sslArgs.ReverseProxyPort == 443 && sslArgs.GNS3Port == 3080 && !sslArgs.FirewallBlock {
					fmt.Printf("%s Using default values: RP=443, GNS3=3080, Firewall=false\n", messageUtils.InfoMsg("Using default values"))
				}
			}

			if validationErr := ssl.ValidateInstallSSLInput(&sslArgs); validationErr != nil {
				return fmt.Errorf("validation error: %w", validationErr)
			}

			fmt.Printf("%s %s\n", messageUtils.Bold("🗑️"), messageUtils.Bold("GNS3 SSL Uninstallation"))
			fmt.Printf("%s\n", messageUtils.Separator(strings.Repeat("─", 50)))
			fmt.Println()

			fmt.Printf("%s Connecting to remote server...\n", messageUtils.InfoMsg("Connecting to remote server"))
			sshClient, err := ssh.ConnectWithKeyOrPassword(hostname, user, sshPort, privateKeyPath, verbose)
			if err != nil {
				return fmt.Errorf("failed to connect via SSH: %w", err)
			}
			defer func() {
				if sshClient != nil {
					_ = sshClient.Close()
				}
			}()
			fmt.Printf("%s Connected successfully\n", messageUtils.SuccessMsg("Connected successfully"))

			fmt.Printf("%s Checking user privileges...\n", messageUtils.InfoMsg("Checking user privileges"))
			if checkPrivErr := sshClient.CheckPrivileges(); checkPrivErr != nil {
				return fmt.Errorf("privilege check failed: %w", checkPrivErr)
			}
			fmt.Printf("%s Privileges verified\n", messageUtils.SuccessMsg("Privileges verified"))

			fmt.Printf("%s Preparing SSL uninstall script...\n", messageUtils.InfoMsg("Preparing SSL uninstall script"))
			script := ssl.GetUninstallScript()
			editedScript := ssl.EditUninstallScriptWithFlags(script, &sslArgs)
			fmt.Printf("%s Script prepared\n", messageUtils.SuccessMsg("Script prepared"))

			fmt.Printf("%s Uninstalling Caddy reverse proxy...\n", messageUtils.InfoMsg("Uninstalling Caddy reverse proxy"))
			success, err := sshClient.ExecuteScript(editedScript, "/tmp/gns3_ssl_uninstall.sh")
			if err != nil {
				return fmt.Errorf("failed to execute uninstall script: %w", err)
			}

			if !success {
				return fmt.Errorf("uninstall script failed")
			}
			fmt.Printf("%s Uninstall completed\n", messageUtils.SuccessMsg("Uninstall completed"))

			if stateManager != nil {
				if err := stateManager.DeleteState(hostname); err != nil {
					fmt.Printf("%s Warning: failed to delete local state: %v\n", messageUtils.WarningMsg("Warning: failed to delete local state"), err)
				} else {
					fmt.Printf("%s Local state cleaned up\n", messageUtils.SuccessMsg("Local state cleaned up"))
				}
			}

			fmt.Printf("\n%s Successfully uninstalled Caddy reverse proxy\n", messageUtils.SuccessMsg("Successfully uninstalled Caddy reverse proxy"))
			fmt.Printf("%s GNS3 server is now accessible on port %d\n", messageUtils.InfoMsg("GNS3 server is now accessible"), sslArgs.GNS3Port)
			return nil
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 22, "SSH port")
	cmd.Flags().StringVarP(&privateKeyPath, "key", "", "", "Path to a custom SSH private key file")
	cmd.Flags().IntVarP(&reverseProxyPort, "reverse-proxy-port", "r", 443, "Port for the reverse proxy that was used")
	cmd.Flags().IntVarP(&gns3Port, "gns3-port", "g", 3080, "Port of the GNS3 Server")
	cmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain that was used for the reverse proxy")
	cmd.Flags().StringVarP(&subject, "subject", "", "/CN=localhost", "Subject that was used for the SSL certificate")
	cmd.Flags().StringVarP(&firewallAllow, "firewall-allow", "a", "", "Firewall allow subnet that was used. Example: 10.0.0.0/24")
	cmd.Flags().BoolVarP(&firewallBlock, "firewall-block", "b", false, "Whether firewall rules were configured")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	cmd.Flags().BoolVarP(&interactive, "interactive", "t", false, "Set the options for this command interactively")

	return cmd
}

func NewUninstallGNS3Cmd() *cobra.Command {
	var (
		port           int
		privateKeyPath string
		homeDir        string
		gns3Port       int
		verbose        bool
		interactive    bool
		preserveData   bool
	)

	cmd := &cobra.Command{
		Use:   "gns3 [user]",
		Short: "Uninstall GNS3 server from remote machine",
		Long: `Uninstall GNS3 server from a remote machine.

This command will:
- Stop and remove GNS3 server service
- Remove GNS3 configuration files
- Remove GNS3 home directory (unless --preserve-data is used)
- Remove GNS3 packages
- Clean up systemd services

Note: The GNS3 user account is preserved to avoid file ownership issues.

The command will attempt to load installation state from a previous installation
to determine the correct configuration values. If no state is found, command line
flags or defaults will be used.

Use --preserve-data to keep the GNS3 home directory with projects and user data.

Note: This will NOT remove virtualization packages (QEMU, Docker, VirtualBox, etc.)
as they may be used by other applications.`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if interactive {
				editedText, err := utils.EditTextWithEditor(gns3.UninstallInteractiveOptionsText, "txt")
				if err != nil {
					return fmt.Errorf("failed to edit options: %w", err)
				}

				interactiveArgs, err := gns3.ParseInteractiveOptions(editedText)
				if err != nil {
					return fmt.Errorf("failed to parse interactive options: %w", err)
				}

				homeDir = interactiveArgs.HomeDir
				gns3Port = interactiveArgs.GNS3Port
				verbose = interactiveArgs.Verbose
				preserveData = interactiveArgs.PreserveData
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			user := args[0]

			cfg, err := config.GetGlobalOptionsFromContext(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get global options: %w", err)
			}

			hostname, sshPort := gns3.ParseServerURLForSSH(cfg.Server, port)

			stateManager, err := gns3.NewStateManager()
			var state *gns3.GNS3ServerState
			if err == nil {
				if localState, loadStateErr := stateManager.LoadState(hostname); loadStateErr == nil {
					state = localState
					fmt.Printf("%s Loaded state from local machine\n", messageUtils.SuccessMsg("Loaded state from local machine"))
				}
			}

			gns3Args := gns3.InstallGNS3Args{
				Username:     "gns3",
				HomeDir:      homeDir,
				ListenHost:   "0.0.0.0",
				GNS3Port:     gns3Port,
				Verbose:      verbose,
				PreserveData: preserveData,
			}

			if state != nil {
				gns3Args.Username = state.Username
				gns3Args.HomeDir = state.HomeDir
				gns3Args.GNS3Port = state.GNS3Port
				fmt.Printf("%s Using saved configuration: User=%s, Home=%s, Port=%d\n",
					messageUtils.InfoMsg("Using saved configuration"), state.Username, state.HomeDir, state.GNS3Port)
			} else if !interactive {
				fmt.Printf("%s No state found, using command line flags or defaults\n", messageUtils.WarningMsg("No state found"))
				if gns3Args.Username == "gns3" && gns3Args.HomeDir == "/opt/gns3" && gns3Args.GNS3Port == 3080 {
					fmt.Printf("%s Using default values: User=gns3, Home=/opt/gns3, Port=3080\n", messageUtils.InfoMsg("Using default values"))
				}
			}

			if validateErr := gns3.ValidateUninstallGNS3Input(&gns3Args); validateErr != nil {
				return fmt.Errorf("validation error: %w", validateErr)
			}

			fmt.Printf("%s %s\n", messageUtils.Bold("🗑️"), messageUtils.Bold("GNS3 Server Uninstallation"))
			fmt.Printf("%s\n", messageUtils.Separator(strings.Repeat("─", 50)))
			fmt.Println()

			fmt.Printf("%s Connecting to remote server...\n", messageUtils.InfoMsg("Connecting to remote server"))
			sshClient, err := ssh.ConnectWithKeyOrPassword(hostname, user, sshPort, privateKeyPath, verbose)
			if err != nil {
				return fmt.Errorf("failed to connect via SSH: %w", err)
			}
			defer func() {
				if sshClient != nil {
					_ = sshClient.Close()
				}
			}()
			fmt.Printf("%s Connected successfully\n", messageUtils.SuccessMsg("Connected successfully"))

			fmt.Printf("%s Checking user privileges...\n", messageUtils.InfoMsg("Checking user privileges"))
			if checkPrivErr := sshClient.CheckPrivileges(); checkPrivErr != nil {
				return fmt.Errorf("privilege check failed: %w", checkPrivErr)
			}
			fmt.Printf("%s Privileges verified\n", messageUtils.SuccessMsg("Privileges verified"))

			fmt.Printf("%s Preparing GNS3 uninstall script...\n", messageUtils.InfoMsg("Preparing GNS3 uninstall script"))
			script := gns3.GetUninstallScript()
			editedScript := gns3.EditUninstallScriptWithFlags(script, &gns3Args)
			fmt.Printf("%s Script prepared\n", messageUtils.SuccessMsg("Script prepared"))

			fmt.Printf("%s Uninstalling GNS3 server...\n", messageUtils.InfoMsg("Uninstalling GNS3 server"))
			success, err := sshClient.ExecuteScript(editedScript, "/tmp/gns3_uninstall.sh")
			if err != nil {
				return fmt.Errorf("failed to execute uninstall script: %w", err)
			}

			if !success {
				return fmt.Errorf("uninstall script failed")
			}
			fmt.Printf("%s Uninstall completed\n", messageUtils.SuccessMsg("Uninstall completed"))

			if stateManager != nil {
				if err := stateManager.DeleteState(hostname); err != nil {
					fmt.Printf("%s Warning: failed to delete local state: %v\n", messageUtils.WarningMsg("Warning: failed to delete local state"), err)
				} else {
					fmt.Printf("%s Local state cleaned up\n", messageUtils.SuccessMsg("Local state cleaned up"))
				}
			}

			fmt.Printf("\n%s GNS3 uninstallation completed successfully\n", messageUtils.SuccessMsg("GNS3 uninstallation completed successfully"))
			fmt.Printf("%s If GNS3 was installed, it has been removed. If not found, no action was needed.\n", messageUtils.InfoMsg("If GNS3 was installed, it has been removed. If not found, no action was needed."))
			return nil
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 22, "SSH port")
	cmd.Flags().StringVarP(&privateKeyPath, "key", "", "", "Path to a custom SSH private key file")
	cmd.Flags().StringVarP(&homeDir, "home-dir", "", "/opt/gns3", "GNS3 home directory to remove")
	cmd.Flags().IntVarP(&gns3Port, "gns3-port", "g", 3080, "Port of the GNS3 Server")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	cmd.Flags().BoolVarP(&interactive, "interactive", "t", false, "Set the options for this command interactively")
	cmd.Flags().BoolVarP(&preserveData, "preserve-data", "", false, "Keep GNS3 home directory and user data")

	return cmd
}
