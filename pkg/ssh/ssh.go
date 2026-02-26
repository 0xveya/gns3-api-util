package ssh

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type SSHClient struct {
	client   *ssh.Client
	config   *ssh.ClientConfig
	hostname string
	username string
	port     int
	verbose  bool
}

type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Success  bool
}

func (c *SSHClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *SSHClient) ExecuteCommand(command string) (*CommandResult, error) {
	if c.client == nil {
		return nil, fmt.Errorf("SSH client not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func() {
		if closeErr := session.Close(); closeErr != nil {
			fmt.Printf("failed to close session: %v", closeErr)
		}
	}()

	var stdout, stderr strings.Builder
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)

	result := &CommandResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0,
		Success:  err == nil,
	}

	if err != nil {
		exitError := &ssh.ExitError{}
		if errors.As(err, &exitError) {
			result.ExitCode = exitError.ExitStatus()
		} else {
			result.ExitCode = -1
		}
	}

	if c.verbose {
		fmt.Printf("Command: %s\n", command)
		fmt.Printf("Exit Code: %d\n", result.ExitCode)
		if result.Stdout != "" {
			fmt.Printf("Stdout: %s\n", result.Stdout)
		}
		if result.Stderr != "" {
			fmt.Printf("Stderr: %s\n", result.Stderr)
		}
	}

	return result, nil
}

func ConnectWithKeyOrPassword(hostname, username string, port int, customPrivateKeyPath string, verbose bool) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // #nosec G106
		Timeout:         30 * time.Second,
	}

	if authMethods := getSSHAgentAuth(); len(authMethods) > 0 {
		config.Auth = append(config.Auth, authMethods...)
		if verbose {
			fmt.Println("Using SSH agent authentication")
		}
	}

	keyPaths := getPrivateKeyPaths(customPrivateKeyPath)
	for _, keyPath := range keyPaths {
		if authMethod, err := getPrivateKeyAuth(keyPath); err == nil {
			config.Auth = append(config.Auth, authMethod)
			if verbose {
				fmt.Printf("Added private key: %s\n", keyPath)
			}
		}
	}

	config.Auth = append(config.Auth, ssh.PasswordCallback(func() (string, error) {
		fmt.Printf("Enter password for %s@%s: ", username, hostname)
		var password string
		_, _ = fmt.Scanln(&password)
		return password, nil
	}))

	address := fmt.Sprintf("%s:%d", hostname, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	if verbose {
		fmt.Printf("Successfully connected to %s@%s:%d\n", username, hostname, port)
	}

	return &SSHClient{
		client:   client,
		config:   config,
		hostname: hostname,
		username: username,
		port:     port,
		verbose:  verbose,
	}, nil
}

func (c *SSHClient) ExecuteScript(scriptContent, remotePath string) (bool, error) {
	createScriptCmd := fmt.Sprintf(`cat > %s << 'SCRIPT_EOF'
%s
SCRIPT_EOF`, remotePath, scriptContent)

	result, err := c.ExecuteCommand(createScriptCmd)
	if err != nil || !result.Success {
		return false, fmt.Errorf("failed to create script: %w", err)
	}

	chmodResult, err := c.ExecuteCommand(fmt.Sprintf("chmod +x %s", remotePath))
	if err != nil || !chmodResult.Success {
		return false, fmt.Errorf("failed to make script executable: %w", err)
	}

	execResult, err := c.ExecuteCommand(fmt.Sprintf("bash %s", remotePath))
	if err != nil {
		return false, fmt.Errorf("failed to execute script: %w", err)
	}

	cleanupResult, _ := c.ExecuteCommand(fmt.Sprintf("rm -f %s", remotePath))
	if c.verbose && cleanupResult != nil && !cleanupResult.Success {
		fmt.Printf("Warning: failed to clean up script file: %s\n", cleanupResult.Stderr)
	}

	// Consider exit code 143 (SIGTERM) as success since it commonly occurs
	// when scripts complete successfully but get terminated by parent process
	return execResult.Success || execResult.ExitCode == 143, nil
}

func (c *SSHClient) CheckPrivileges() error {
	uidResult, err := c.ExecuteCommand("id -u")
	if err != nil {
		return fmt.Errorf("failed to check user ID: %w", err)
	}

	if uidResult.Stdout == "0\n" {
		if c.verbose {
			fmt.Println("User is root")
		}
		return nil
	}

	sudoResult, err := c.ExecuteCommand("sudo -n true")
	if err != nil {
		return fmt.Errorf("user lacks passwordless sudo access: %w", err)
	}

	if !sudoResult.Success {
		return fmt.Errorf("user lacks passwordless sudo access")
	}

	if c.verbose {
		fmt.Println("User has passwordless sudo access")
	}

	return nil
}

func getPrivateKeyPaths(customPath string) []string {
	var paths []string

	if customPath != "" {
		paths = append(paths, customPath)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return paths
	}

	sshDir := filepath.Join(homeDir, ".ssh")

	keyTypes := []string{"id_rsa", "id_dsa", "id_ecdsa", "id_ed25519", "id_rsa_ed25519"}

	for _, keyType := range keyTypes {
		keyPath := filepath.Join(sshDir, keyType)
		if _, err := os.Stat(keyPath); err == nil {
			paths = append(paths, keyPath)
		}
	}

	return paths
}

func loadPrivateKey(keyPath string) (ssh.Signer, error) {
	key, err := os.ReadFile(keyPath) // #nosec G304
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return signer, nil
}

func getPrivateKeyAuth(keyPath string) (ssh.AuthMethod, error) {
	signer, err := loadPrivateKey(keyPath)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func getSSHAgentAuth() []ssh.AuthMethod {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var d net.Dialer
	sshAgent, err := d.DialContext(ctx, "unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil
	}
	defer func() {
		if closeAgentErr := sshAgent.Close(); closeAgentErr != nil {
			fmt.Printf("failed to close SSH agent connection: %v", closeAgentErr)
		}
	}()

	agentClient := agent.NewClient(sshAgent)
	signers, err := agentClient.Signers()
	if err != nil {
		return nil
	}

	if len(signers) == 0 {
		return nil
	}

	return []ssh.AuthMethod{ssh.PublicKeys(signers...)}
}
