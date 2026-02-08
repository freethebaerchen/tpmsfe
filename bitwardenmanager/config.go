package bitwardenmanager

import "os"

// Config holds configuration for Bitwarden authentication
type Config struct {
	Method      string // "cli" or "api"
	SessionKey  string // For CLI authentication
	ServerURL   string // Bitwarden server URL
	Vault       string // Vault/folder name (for CLI) or organization ID (for API)
	Title       string // Secret/item title
	AccessToken string // For API authentication
	ClientID    string // For API authentication
	ProjectID   string // For API authentication (optional)
}

// ValidateAndSetDefaults validates the config and sets defaults from environment variables.
func (c *Config) ValidateAndSetDefaults() error {
	if v := os.Getenv("TPMSFE_BW_AUTH"); v != "" {
		c.Method = v
	} else if v := os.Getenv("BW_AUTH"); v != "" {
		c.Method = v
	} else if c.Method == "" {
		c.Method = "cli"
	}

	if v := os.Getenv("TPMSFE_BW_API_URL"); v != "" {
		c.ServerURL = v
	} else if v := os.Getenv("BW_API_URL"); v != "" {
		c.ServerURL = v
	} else if c.ServerURL == "" {
		c.ServerURL = "https://bitwarden.com"
	}

	if v := os.Getenv("TPMSFE_BW_CLIENT_ID"); v != "" {
		c.ClientID = v
	} else if v := os.Getenv("BW_CLIENT_ID"); v != "" {
		c.ClientID = v
	}

	if v := os.Getenv("TPMSFE_BW_PROJECT_ID"); v != "" {
		c.ProjectID = v
	} else if v := os.Getenv("BW_PROJECT_ID"); v != "" {
		c.ProjectID = v
	}

	switch c.Method {
	case "cli":
		if c.SessionKey == "" {
			c.SessionKey = os.Getenv("TPMSFE_BW_SESSION")
			if c.SessionKey == "" {
				c.SessionKey = os.Getenv("BW_SESSION")
				if c.SessionKey == "" {
					return ErrSessionKeyRequired
				}
			}
		}
		return nil

	case "api":
		if c.AccessToken == "" {
			c.AccessToken = os.Getenv("TPMSFE_BW_ACCESS_TOKEN")
			if c.AccessToken == "" {
				c.AccessToken = os.Getenv("BW_ACCESS_TOKEN")
				if c.AccessToken == "" {
					return ErrAccessTokenRequired
				}
			}
		}
		if c.ClientID == "" {
			return ErrClientIDRequired
		}
		return nil

	default:
		return ErrUnsupportedMethod
	}
}
