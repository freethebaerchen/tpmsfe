package onepasswordmanager

import (
	"context"
	"fmt"
	"os"

	onepassword "github.com/1password/onepassword-sdk-go"
)

func fetchServiceAccount(serviceAccountToken, vault, title string) (string, error) {
	serviceAccountTokenEnv := os.Getenv("TPMSFE_OP_SERVICE_ACCOUNT_TOKEN")
	if serviceAccountToken == "" {
		serviceAccountToken = serviceAccountTokenEnv
	}
	if serviceAccountToken == "" {
		return "", fmt.Errorf("service account token required for service-account authentication")
	}
	ctx := context.Background()
	client, err := onepassword.NewClient(
		ctx,
		onepassword.WithServiceAccountToken(serviceAccountToken),
		onepassword.WithIntegrationInfo("TPMSFE", "1.0.0"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to initialize 1Password service account client: %w", err)
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