package onepasswordmanager

import (
	"fmt"
	"os"

	connect "github.com/1Password/connect-sdk-go/connect"
)

func fetchConnect(connectEndpoint, connectToken, vault, title string) (string, error) {
	// Use flag value if provided, else fall back to TPMSFE-prefixed env, then set OP_CONNECT_TOKEN/OP_CONNECT_HOST for SDK
	connectTokenEnv := os.Getenv("TPMSFE_OP_CONNECT_TOKEN")
	if connectToken == "" {
		connectToken = connectTokenEnv
	}
	if connectToken == "" {
		return "", fmt.Errorf("connect token required for connect authentication")
	}
	os.Setenv("OP_CONNECT_TOKEN", connectToken)
	os.Setenv("OP_CONNECT_HOST", connectEndpoint)
	client, err := connect.NewClientFromEnvironment()
	if err != nil {
		return "", fmt.Errorf("failed to initialize 1Password Connect client: %w", err)
	}

	// List all items in the vault to find the item UUID by title
	items, err := client.GetItems(vault)
	if err != nil {
		return "", fmt.Errorf("failed to list items in vault: %w", err)
	}
	var itemUUID string
	for _, item := range items {
		if item.Title == title {
			itemUUID = item.ID
			break
		}
	}
	if itemUUID == "" {
		return "", fmt.Errorf("No item found with title '%s' in vault '%s'", title, vault)
	}
	// Fetch the item by UUID
	item, err := client.GetItem(itemUUID, vault)
	if err != nil {
		return "", fmt.Errorf("failed to fetch item from 1Password Connect: %w", err)
	}
	// Try to get the 'password' field, then 'credential'
	var value string
	for _, field := range item.Fields {
		if field.Label == "password" && field.Value != "" {
			value = field.Value
			break
		}
	}
	if value == "" {
		for _, field := range item.Fields {
			if field.Label == "credential" && field.Value != "" {
				value = field.Value
				break
			}
		}
	}
	if value == "" {
		return "", fmt.Errorf("neither 'password' nor 'credential' field found in item '%s'", title)
	}
	return value, nil
}
