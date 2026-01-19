package bitwardenmanager

import (
	"fmt"
	"os"
)

func CheckMethod(method, sessionKey, serverURL, vaultOrFolder, title, accessToken, organizationID, projectID string) (string, error) {
	switch method {
	case "cli":
		// Validate CLI requirements
		if sessionKey == "" {
			sessionKeyEnv := os.Getenv("TPMSFE_BW_SESSION")
			if sessionKeyEnv == "" {
				return "", fmt.Errorf("Bitwarden session key is required for CLI authentication. Run 'bw unlock' and provide via -bw-session flag or TPMSFE_BW_SESSION env var")
			}
			sessionKey = sessionKeyEnv
		}
		return fetchBitwardenCLI(sessionKey, serverURL, vaultOrFolder, title)

	case "api":
		// Validate API requirements
		if accessToken == "" {
			accessTokenEnv := os.Getenv("TPMSFE_BW_ACCESS_TOKEN")
			if accessTokenEnv == "" {
				return "", fmt.Errorf("Bitwarden access token is required for API authentication. Provide via -bw-access-token flag or TPMSFE_BW_ACCESS_TOKEN env var")
			}
			accessToken = accessTokenEnv
		}
		if organizationID == "" {
			return "", fmt.Errorf("Bitwarden organization ID is required for API authentication via -bw-org-id flag")
		}
		return fetchBitwardenAPI(serverURL, accessToken, organizationID, projectID, title)

	default:
		return "", fmt.Errorf("unsupported Bitwarden authentication method: %s. Supported methods are: cli, api", method)
	}
}
