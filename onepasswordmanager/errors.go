package onepasswordmanager

import "fmt"

var (
	ErrAccountNameRequired         = fmt.Errorf("1Password account name is required for desktop authentication. Provide via -op-account flag or TPMSFE_OP_ACCOUNT_NAME env var")
	ErrConnectTokenRequired        = fmt.Errorf("1Password Connect token is required for Connect authentication. Provide via -op-connect-token flag or TPMSFE_OP_CONNECT_TOKEN env var")
	ErrServiceAccountTokenRequired = fmt.Errorf("1Password service account token is required for service-account authentication. Provide via -op-service-account-token flag or TPMSFE_OP_SERVICE_ACCOUNT_TOKEN env var")
	ErrUnsupportedMethod          = fmt.Errorf("unsupported 1Password authentication method. Supported methods are: desktop, connect, service-account")
)
