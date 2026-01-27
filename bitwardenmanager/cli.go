package bitwardenmanager

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// fetchCLI fetches a secret using the Bitwarden CLI
func fetchCLI(config Config) (string, error) {
	// Sync vault to ensure we have latest data
	if err := syncVault(config.SessionKey); err != nil {
		return "", fmt.Errorf("failed to sync Bitwarden vault: %w", err)
	}

	// First, try the simple approach: get password directly (works for login items)
	if password, err := getPasswordDirect(config); err == nil {
		return password, nil
	}

	// If that fails, search for the item and get full details (needed for custom fields)
	return getPasswordFromItem(config)
}

// syncVault syncs the Bitwarden vault to ensure we have the latest data
func syncVault(sessionKey string) error {
	cmd := exec.Command("bw", "sync")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("bw sync failed: %w", err)
	}
	return nil
}

// getPasswordDirect tries to get password directly using 'bw get password'
// This is simpler and faster for login items
func getPasswordDirect(config Config) (string, error) {
	cmd := exec.Command("bw", "get", "password", config.Title)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("direct password fetch failed: %w", err)
	}

	password := strings.TrimSpace(string(output))
	if password == "" {
		return "", fmt.Errorf("empty password returned for item '%s'", config.Title)
	}

	return password, nil
}

// getPasswordFromItem searches for the item and extracts password from login or custom fields
func getPasswordFromItem(config Config) (string, error) {
	// Search for items matching the title
	itemID, err := findItemID(config.Title, config.SessionKey)
	if err != nil {
		return "", fmt.Errorf("failed to find item '%s': %w", config.Title, err)
	}

	// Get the full item details
	cmd := exec.Command("bw", "get", "item", itemID, "--session", config.SessionKey)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to fetch item '%s': %w", config.Title, err)
	}

	// Parse JSON response
	var item struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Login *struct {
			Password string `json:"password"`
		} `json:"login"`
		Fields []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"fields"`
	}

	if err := json.Unmarshal(output, &item); err != nil {
		return "", fmt.Errorf("failed to parse Bitwarden CLI response: %w", err)
	}

	// Verify we got the right item
	if item.Name != config.Title {
		return "", fmt.Errorf("item name mismatch: expected '%s', got '%s'", config.Title, item.Name)
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

	return "", fmt.Errorf("neither 'password' nor 'credential' field found in item '%s'", config.Title)
}

// findItemID searches for an item by name and returns its ID
func findItemID(title, sessionKey string) (string, error) {
	cmd := exec.Command("bw", "list", "items", "--session", sessionKey, "--search", title)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to search for items: %w", err)
	}

	var items []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	if err := json.Unmarshal(output, &items); err != nil {
		return "", fmt.Errorf("failed to parse search results: %w", err)
	}

	// Find exact match (case-sensitive)
	for _, item := range items {
		if item.Name == title {
			return item.ID, nil
		}
	}

	// Try case-insensitive match
	for _, item := range items {
		if strings.EqualFold(item.Name, title) {
			return item.ID, nil
		}
	}

	return "", fmt.Errorf("item '%s' not found", title)
}
