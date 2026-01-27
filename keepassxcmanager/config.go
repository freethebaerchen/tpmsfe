package keepassxcmanager

import "os"

// Config holds configuration for KeePassXC authentication
type Config struct {
	DatabasePath string // Path to .kdbx database file
	Password     string // Database password (optional if keyfile is used)
	EntryPath    string // Entry path (e.g., "Group/Entry" or just entry title)
	FieldName    string // Field name to retrieve (default: "Password")
}

// ValidateAndSetDefaults validates the config and sets defaults from environment variables
func (c *Config) ValidateAndSetDefaults() error {
	if c.DatabasePath == "" {
		c.DatabasePath = os.Getenv("TPMSFE_KEEPASSXC_DATABASE")
		if c.DatabasePath == "" {
			c.DatabasePath = os.Getenv("KEEPASSXC_DATABASE")
			if c.DatabasePath == "" {
				return ErrDatabasePathRequired
			}
		}
	}

	// Password or keyfile must be provided
	if c.Password == "" {
		c.Password = os.Getenv("TPMSFE_KEEPASSXC_PASSWORD")
		if c.Password == "" {
			c.Password = os.Getenv("KEEPASSXC_PASSWORD")
			if c.Password == "" {
				return ErrPasswordOrKeyfileRequired
			}
		}
	}

	// Default field name is "Password"
	if c.FieldName == "" {
		c.FieldName = "Password"
	}

	return nil
}
