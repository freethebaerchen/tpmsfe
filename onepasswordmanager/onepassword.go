package onepasswordmanager

// FetchSecret fetches a secret from 1Password using the configured method
func FetchSecret(config Config) (string, error) {
	if err := config.ValidateAndSetDefaults(); err != nil {
		return "", err
	}

	switch config.Method {
	case "desktop":
		return fetchDesktop(config)
	case "connect":
		return fetchConnect(config)
	case "service-account":
		return fetchServiceAccount(config)
	default:
		return "", ErrAuthMethodRequired
	}
}

// CheckMethod is a convenience function that maintains backward compatibility
// with the old API signature
func CheckMethod(authMethod, accountName, connectEndpoint, connectToken, serviceAccountToken, vault, title string) (string, error) {
	config := Config{
		Method:              authMethod,
		AccountName:         accountName,
		ConnectEndpoint:     connectEndpoint,
		ConnectToken:        connectToken,
		ServiceAccountToken: serviceAccountToken,
		Vault:               vault,
		Title:               title,
	}
	return FetchSecret(config)
}
