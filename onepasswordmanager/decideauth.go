package onepasswordmanager

import (
	"fmt"
)

// Check1Password fetches the value of the specified item in the specified vault using the correct authentication method.
func CheckMethod(authMethod, accountName, connectEndpoint, connectToken, serviceAccountToken, vault, title string) (string, error) {
	switch authMethod {
	case "desktop":
		return fetchDesktop(accountName, vault, title)
	case "connect":
		return fetchConnect(connectEndpoint, connectToken, vault, title)
	case "service-account":
		return fetchServiceAccount(serviceAccountToken, vault, title)
	default:
		return "", fmt.Errorf("unsupported authentication method: %s", authMethod)
	}
}
