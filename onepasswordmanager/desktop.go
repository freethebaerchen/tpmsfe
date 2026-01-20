package onepasswordmanager

import (
	"context"
	"fmt"

	onepassword "github.com/1password/onepassword-sdk-go"
)

// fetchDesktop fetches a secret using 1Password Desktop integration
func fetchDesktop(config Config) (string, error) {
	ctx := context.Background()
	client, err := onepassword.NewClient(
		ctx,
		onepassword.WithDesktopAppIntegration(config.AccountName),
		onepassword.WithIntegrationInfo("TPMSFE", "1.0.0"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to initialize 1Password desktop client: %w", err)
	}

	// Try password field first
	secret, err := client.Secrets().Resolve(ctx, fmt.Sprintf("op://%s/%s/password", config.Vault, config.Title))
	if err == nil && secret != "" {
		return secret, nil
	}

	// Try credential field if password is not found
	secret, err = client.Secrets().Resolve(ctx, fmt.Sprintf("op://%s/%s/credential", config.Vault, config.Title))
	if err == nil && secret != "" {
		return secret, nil
	}

	return "", fmt.Errorf("neither 'password' nor 'credential' field found in item '%s'", config.Title)
}
