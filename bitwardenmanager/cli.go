package bitwardenmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func fetchBitwardenCLI(sessionKey, serverURL, folderName, title string) (string, error) {
	// Use bw CLI with session key (from `bw unlock`)
	if sessionKey == "" {
		sessionKey = os.Getenv("TPMSFE_BW_SESSION")
	}
	if sessionKey == "" {
		return "", fmt.Errorf("Bitwarden session key required. Run 'bw unlock' and provide session key")
	}

	// Set server URL if provided
	if serverURL != "" && serverURL != "https://vault.bitwarden.com" {
		configCmd := exec.Command("bw", "config", "server", serverURL)
		if err := configCmd.Run(); err != nil {
			return "", fmt.Errorf("failed to set Bitwarden server URL: %w", err)
		}
	}

	// Get item by title
	cmd := exec.Command("bw", "get", "item", title, "--session", sessionKey)
	output, err := cmd.Output()
	if err != nil {
		// Try to get more specific error message
		if exitErr, ok := err.(*exec.ExitError); ok && len(exitErr.Stderr) > 0 {
			return "", fmt.Errorf("failed to fetch item from Bitwarden: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to fetch item from Bitwarden: %w", err)
	}

	// Parse JSON response
	var item struct {
		Name   string `json:"name"`
		Login  *struct {
			Password string `json:"password"`
		} `json:"login"`
		Fields []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"fields"`
	}

	if err := json.Unmarshal(output, &item); err != nil {
		return "", fmt.Errorf("failed to parse Bitwarden response: %w", err)
	}

	// Check if we got the right item
	if item.Name != title {
		return "", fmt.Errorf("item name mismatch: expected '%s', got '%s'", title, item.Name)
	}

	// Try to get password from login field first
	if item.Login != nil && item.Login.Password != "" {
		return item.Login.Password, nil
	}

	// Then check custom fields for 'password' or 'credential'
	for _, field := range item.Fields {
		fieldName := strings.ToLower(field.Name)
		if (fieldName == "password" || fieldName == "credential") && field.Value != "" {
			return field.Value, nil
		}
	}

	return "", fmt.Errorf("neither 'password' nor 'credential' field found in item '%s'", title)
}
