package keepassxcmanager

// FetchSecret fetches a secret from KeePassXC using the configured method
func FetchSecret(config Config) (string, error) {
	if err := config.ValidateAndSetDefaults(); err != nil {
		return "", err
	}

	// Currently only CLI method is supported
	return fetchCLI(config)
}

// CheckMethod is a convenience function that maintains backward compatibility
// with the old API signature (if needed in the future)
func CheckMethod(databasePath, password, entryPath, fieldName string) (string, error) {
	config := Config{
		DatabasePath: databasePath,
		Password:     password,
		EntryPath:    entryPath,
		FieldName:    fieldName,
	}
	return FetchSecret(config)
}
