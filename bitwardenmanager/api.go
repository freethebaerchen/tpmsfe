package bitwardenmanager

import (
	"fmt"
	"os"

	bitwarden "github.com/bitwarden/sdk-go"
)

func fetchBitwardenAPI(apiURL, accessToken, organizationID, projectID, title string) (string, error) {
	// Set API URL if provided
	if apiURL != "" {
		os.Setenv("BITWARDEN_API_URL", apiURL)
		os.Setenv("BITWARDEN_IDENTITY_API_URL", apiURL)
	}

	// Initialize Bitwarden client
	client, err := bitwarden.NewBitwardenClient(nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Bitwarden client: %w", err)
	}
	defer client.Close()

	// Authenticate with access token
	err = client.AccessTokenLogin(accessToken, nil)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate with Bitwarden: %w", err)
	}

	// List secrets in organization
	secrets, err := client.Secrets().List(organizationID)
	if err != nil {
		return "", fmt.Errorf("failed to list secrets: %w", err)
	}

	// Find secret by key/title
	var secretID string
	for _, secret := range secrets.Data {
		if secret.Key == title {
			secretID = secret.ID
			break
		}
	}
	
	if secretID == "" {
		return "", fmt.Errorf("No secret found with title '%s' in organization '%s'", title, organizationID)
	}

	// Get the secret value
	secret, err := client.Secrets().Get(secretID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch secret from Bitwarden: %w", err)
	}

	return secret.Value, nil
}
