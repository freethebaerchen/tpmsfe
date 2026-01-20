package keepassxcmanager

import "fmt"

var (
	ErrDatabasePathRequired      = fmt.Errorf("KeePassXC database path is required. Provide via -kx-database flag or TPMSFE_KEEPASSXC_DATABASE env var")
	ErrPasswordOrKeyfileRequired = fmt.Errorf("KeePassXC password or keyfile is required. Provide via -kx-password/-kx-keyfile flags or TPMSFE_KEEPASSXC_PASSWORD/TPMSFE_KEEPASSXC_KEYFILE env vars")
	ErrUnsupportedMethod        = fmt.Errorf("unsupported KeePassXC method. Currently only CLI is supported")
)
