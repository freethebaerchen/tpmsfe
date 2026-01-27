package keepassxcmanager

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// fetchCLI fetches a secret using KeePassXC CLI
func fetchCLI(config Config) (string, error) {
	// Build kpcli command
	// Format: kpcli show -a <field> <database> <entry> [-k <keyfile>]
	args := []string{"show", "-a", config.FieldName, config.DatabasePath, config.EntryPath}

	cmd := exec.Command("kpcli", args...)

	// If password is provided, pass it via stdin (kpcli reads password from stdin)
	var stdin bytes.Buffer
	if config.Password != "" {
		stdin.WriteString(config.Password + "\n")
		cmd.Stdin = &stdin
	}

	// Execute command
	output, err := cmd.Output()
	if err != nil {
		// Try to get more specific error message
		if exitErr, ok := err.(*exec.ExitError); ok && len(exitErr.Stderr) > 0 {
			return "", fmt.Errorf("failed to fetch secret from KeePassXC CLI: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to fetch secret from KeePassXC CLI: %w", err)
	}

	// Clean up output (remove trailing newline)
	result := strings.TrimSpace(string(output))
	if result == "" {
		return "", fmt.Errorf("empty result from KeePassXC for entry '%s' field '%s'", config.EntryPath, config.FieldName)
	}

	return result, nil
}
