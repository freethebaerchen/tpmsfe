package onepasswordmanager

import "os"

// Config holds configuration for 1Password authentication
type Config struct {
	Method              string // "desktop", "connect", or "service-account"
	AccountName         string // For desktop authentication
	ConnectEndpoint     string // For Connect authentication
	ConnectToken        string // For Connect authentication
	ServiceAccountToken string // For service-account authentication
	Vault               string // Vault name
	Title               string // Item title
}

// ValidateAndSetDefaults validates the config and sets defaults from environment variables
func (c *Config) ValidateAndSetDefaults() error {
	if c.Method == "" {
		c.AccountName = os.Getenv("TPMSFE_OP_AUTH_METHOD")
		if c.AccountName == "" {
			c.AccountName = os.Getenv("OP_AUTH_METHOD")
			if c.AccountName == "" {
				return ErrAuthMethodRequired
			}
		}
	}

	switch c.Method {
	case "desktop":
		if c.AccountName == "" {
			c.AccountName = os.Getenv("TPMSFE_OP_ACCOUNT_NAME")
			if c.AccountName == "" {
				c.AccountName = os.Getenv("OP_ACCOUNT_NAME")
				if c.AccountName == "" {
					return ErrAccountNameRequired
				}
			}
		}
		return nil

	case "connect":
		if c.ConnectToken == "" {
			c.ConnectToken = os.Getenv("TPMSFE_OP_CONNECT_TOKEN")
			if c.ConnectToken == "" {
				c.ConnectToken = os.Getenv("OP_CONNECT_TOKEN")
				if c.ConnectToken == "" {
					return ErrConnectTokenRequired
				}
			}
		}
		if c.ConnectEndpoint == "" {
			c.ConnectEndpoint = os.Getenv("TPMSFE_OP_CONNECT_HOST")
			if c.ConnectEndpoint == "" {
				c.ConnectEndpoint = os.Getenv("OP_CONNECT_HOST")
				if c.ConnectEndpoint == "" {
					c.ConnectEndpoint = "http://localhost:8080"
				}
			}
		}
		return nil

	case "service-account":
		if c.ServiceAccountToken == "" {
			c.ServiceAccountToken = os.Getenv("TPMSFE_OP_SERVICE_ACCOUNT_TOKEN")
			if c.ServiceAccountToken == "" {
				c.ServiceAccountToken = os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")
				if c.ServiceAccountToken == "" {
					return ErrServiceAccountTokenRequired
				}
			}
		}
		return nil

	default:
		return ErrAuthMethodRequired
	}
}
