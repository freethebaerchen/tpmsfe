package onepasswordmanager

import (
	"context"
	"fmt"
	"os"

	onepassword "github.com/1password/onepassword-sdk-go"
)

func fetchDesktop(accountName, vault, title string) (string, error) {
	accountNameEnv := os.Getenv("TPMSFE_OP_ACCOUNT_NAME")
	if accountName == "" {
		accountName = accountNameEnv
	}
	if accountName == "" {
		return "", fmt.Errorf("account name required for desktop authentication")
	}
	ctx := context.Background()
	client, err := onepassword.NewClient(
		ctx,
		onepassword.WithDesktopAppIntegration(accountName),
		onepassword.WithIntegrationInfo("TPMSFE", "1.0.0"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to initialize 1Password desktop client: %w", err)
	}

	secret, err := client.Secrets().Resolve(ctx, fmt.Sprintf("op://%s/%s/password", vault, title))
	if err == nil && secret != "" {
		return secret, nil
	}
	// Try credential field if password is not found
	secret, err = client.Secrets().Resolve(ctx, fmt.Sprintf("op://%s/%s/credential", vault, title))
	if err == nil && secret != "" {
		return secret, nil
	}
	return "", fmt.Errorf("neither 'password' nor 'credential' field found in item '%s'", title)
}
