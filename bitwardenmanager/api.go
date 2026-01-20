package bitwardenmanager

import (
	"fmt"
	"os"

	bitwarden "github.com/bitwarden/sdk-go"
)

// fetchAPI fetches a secret using the Bitwarden SDK/API
func fetchAPI(config Config) (string, error) {
	// Set API URL if provided
	if config.ServerURL != "" {
		os.Setenv("BITWARDEN_API_URL", config.ServerURL)
		os.Setenv("BITWARDEN_IDENTITY_API_URL", config.ServerURL)
	}

	// Initialize Bitwarden client
	client, err := bitwarden.NewBitwardenClient(nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Bitwarden SDK client: %w", err)
	}
	defer client.Close()

	// Authenticate with access token
	err = client.AccessTokenLogin(config.AccessToken, nil)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate with Bitwarden SDK: %w", err)
	}

	// List secrets in client
	secrets, err := client.Secrets().List(config.ClientID)
	if err != nil {
		return "", fmt.Errorf("failed to list secrets in organization '%s': %w", config.ClientID, err)
	}

	// Find secret by key/title
	var secretID string
	for _, secret := range secrets.Data {
		if secret.Key == config.Title {
			secretID = secret.ID
			break
		}
	}

	if secretID == "" {
		return "", fmt.Errorf("no secret found with title '%s' in organization '%s'", config.Title, config.ClientID)
	}

	// Get the secret value
	secret, err := client.Secrets().Get(secretID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch secret '%s' from Bitwarden SDK: %w", config.Title, err)
	}

	return secret.Value, nil
}
