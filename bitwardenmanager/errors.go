package bitwardenmanager

import "fmt"

var (
	ErrSessionKeyRequired     = fmt.Errorf("Bitwarden session key is required for CLI authentication. Run 'bw unlock' and provide via -bw-session flag or TPMSFE_BW_SESSION env var")
	ErrAccessTokenRequired    = fmt.Errorf("Bitwarden access token is required for API authentication. Provide via -bw-access-token flag or TPMSFE_BW_ACCESS_TOKEN env var")
	ErrOrganizationIDRequired = fmt.Errorf("Bitwarden organization ID is required for API authentication via -bw-org-id flag")
	ErrUnsupportedMethod     = fmt.Errorf("unsupported Bitwarden authentication method. Supported methods are: cli, api")
)
