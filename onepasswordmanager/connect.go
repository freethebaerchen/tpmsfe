package onepasswordmanager

import (
	"fmt"
	"os"

	connect "github.com/1Password/connect-sdk-go/connect"
)

// fetchConnect fetches a secret using 1Password Connect
func fetchConnect(config Config) (string, error) {
	// Set environment variables for Connect SDK
	os.Setenv("OP_CONNECT_TOKEN", config.ConnectToken)
	os.Setenv("OP_CONNECT_HOST", config.ConnectEndpoint)

	client, err := connect.NewClientFromEnvironment()
	if err != nil {
		return "", fmt.Errorf("failed to initialize 1Password Connect client: %w", err)
	}

	// List all items in the vault to find the item UUID by title
	items, err := client.GetItems(config.Vault)
	if err != nil {
		return "", fmt.Errorf("failed to list items in vault '%s': %w", config.Vault, err)
	}

	var itemUUID string
	for _, item := range items {
		if item.Title == config.Title {
			itemUUID = item.ID
			break
		}
	}

	if itemUUID == "" {
		return "", fmt.Errorf("no item found with title '%s' in vault '%s'", config.Title, config.Vault)
	}

	// Fetch the item by UUID
	item, err := client.GetItem(itemUUID, config.Vault)
	if err != nil {
		return "", fmt.Errorf("failed to fetch item '%s' from 1Password Connect: %w", config.Title, err)
	}

	// Try to get the 'password' field first
	for _, field := range item.Fields {
		if field.Label == "password" && field.Value != "" {
			return field.Value, nil
		}
	}

	// Try 'credential' field if password is not found
	for _, field := range item.Fields {
		if field.Label == "credential" && field.Value != "" {
			return field.Value, nil
		}
	}

	return "", fmt.Errorf("neither 'password' nor 'credential' field found in item '%s'", config.Title)
}
