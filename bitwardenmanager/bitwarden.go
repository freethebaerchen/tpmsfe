package bitwardenmanager

// FetchSecret fetches a secret from Bitwarden using the configured method
func FetchSecret(config Config) (string, error) {
	if err := config.ValidateAndSetDefaults(); err != nil {
		return "", err
	}

	switch config.Method {
	case "cli":
		return fetchCLI(config)
	case "api":
		return fetchAPI(config)
	default:
		return "", ErrUnsupportedMethod
	}
}

// CheckMethod is a convenience function that maintains backward compatibility
// with the old API signature
func CheckMethod(method, sessionKey, serverURL, vaultOrFolder, title, accessToken, organizationID, projectID string) (string, error) {
	config := Config{
		Method:         method,
		SessionKey:     sessionKey,
		ServerURL:      serverURL,
		Vault:          vaultOrFolder,
		Title:          title,
		AccessToken:    accessToken,
		OrganizationID: organizationID,
		ProjectID:      projectID,
	}
	return FetchSecret(config)
}
