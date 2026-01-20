package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/freethebaerchen/tpmsfe/bitwardenmanager"
	"github.com/freethebaerchen/tpmsfe/encryptor"
	"github.com/freethebaerchen/tpmsfe/keepassxcmanager"
	"github.com/freethebaerchen/tpmsfe/onepasswordmanager"
)

const (
	magicHeader = "OpenTofu-External-Encryption-Method"
	version     = 1
)

type ExternalRequest struct {
	Payload string `json:"payload"`
	Key     string `json:"key,omitempty"`
}

type ExternalResponse struct {
	Payload string `json:"payload"`
}

func main() {
	provider := flag.String("provider", "", "The provider you use as password manager. Possible values:\n - 1password\n - bitwarden\n - kepassxc.\n Default is none\n")
	vault := flag.String("vault", "", "The vault you want to use inside your password manager. Default is none")
	title := flag.String("title", "", "The title of the password entry you want to use. Default is none")
	opAuthMethod := flag.String("op-auth", "desktop", "The authentication method for 1Password. Possible values:\n - desktop\n - service-account\n - connect\n Default is desktop")
	opAccountName := flag.String("op-account", "", "The 1Password account name. Required, when using 1Password Desktop.")
	opConnectEndpoint := flag.String("op-connect-endpoint", "http://127.0.0.1:8080", "The 1Password Connect endpoint. Required, when using 1Password Connect. Default is http://127.0.0.1:8080")
	opConnectToken := flag.String("op-connect-token", "", "The 1Password Connect token. Required, when using 1Password Connect.")
	opServiceAccountToken := flag.String("op-service-account-token", "", "The 1Password service account token. Required, when using 1Password service account authentication.")
	bwAuthMethod := flag.String("bw-auth", "cli", "The authentication method for Bitwarden. Possible values:\n - cli\n - api\n Default is cli")
	bwSessionKey := flag.String("bw-session", "", "Bitwarden CLI session key. Required when using CLI authentication.")
	bwApiUrl := flag.String("bw-api-url", "https://bitwarden.com", "The Bitwarden API URL. Default is https://bitwarden.com")
	bwAccessToken := flag.String("bw-access-token", "", "The Bitwarden access token. Required for Bitwarden authentication.")
	bwClientId := flag.String("bw-client-id", "", "The Bitwarden organization ID. Required to access organization secrets.")
	bwProjectId := flag.String("bw-project-id", "", "The Bitwarden project ID. Required to access project secrets.")
	kxDatabase := flag.String("kx-database", "", "The KeePassXC database path (.kdbx file). Required for KeePassXC.")
	kxPassword := flag.String("kx-password", "", "The KeePassXC database password. Required if keyfile is not provided.")
	kxFieldName := flag.String("kx-field", "Password", "The KeePassXC field name to retrieve. Default is 'Password'.")
	flag.Parse()

	// Output the magic header for OpenTofu identification (step 1)
	header := map[string]interface{}{
		"magic":   magicHeader,
		"version": version,
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling header: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(headerJSON))

	// Fetch the secret from password manager
	secret, err := fetchSecret(
		*provider, *vault, *title,
		*opAuthMethod, *opAccountName, *opConnectEndpoint, *opConnectToken, *opServiceAccountToken,
		*bwAuthMethod, *bwSessionKey, *bwApiUrl, *bwAccessToken, *bwClientId, *bwProjectId,
		*kxDatabase, *kxPassword, *kxFieldName,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching secret: %v\n", err)
		os.Exit(1)
	}

	// Read JSON request from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	var req ExternalRequest
	if err := json.Unmarshal(input, &req); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON request: %v\n", err)
		os.Exit(1)
	}

	if req.Payload == "" {
		fmt.Fprintf(os.Stderr, "Error: request must contain 'payload'\n")
		os.Exit(1)
	}

	// Decode the payload to check its format
	payloadBytes, err := base64.StdEncoding.DecodeString(req.Payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding payload: %v\n", err)
		os.Exit(1)
	}

	// Determine if we should encrypt or decrypt
	// Our encrypted format: 32 bytes salt + 12 bytes nonce (GCM standard) + ciphertext
	// Minimum encrypted size is 32 + 12 = 44 bytes
	// If payload is smaller or decryption fails, treat as plaintext and encrypt
	var resp ExternalResponse
	const minEncryptedSize = 44 // 32 bytes salt + 12 bytes nonce (minimum)

	if len(payloadBytes) >= minEncryptedSize {
		// Try to decrypt - if it succeeds, this was encrypted data
		decryptedBytes, err := encryptor.DecryptBytes(secret, req.Payload)
		if err == nil {
			// Successfully decrypted - return decrypted payload
			resp.Payload = base64.StdEncoding.EncodeToString(decryptedBytes)
		} else {
			// Decryption failed - treat as plaintext and encrypt
			ciphertextBase64, err := encryptor.EncryptBytes(secret, payloadBytes)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error encrypting: %v\n", err)
				os.Exit(1)
			}
			resp.Payload = ciphertextBase64
		}
	} else {
		// Too small to be encrypted data - treat as plaintext and encrypt
		ciphertextBase64, err := encryptor.EncryptBytes(secret, payloadBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encrypting: %v\n", err)
			os.Exit(1)
		}
		resp.Payload = ciphertextBase64
	}

	// Output JSON response
	respJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(respJSON))
}

func fetchSecret(
	provider, vault, title string,
	opAuthMethod, opAccountName, opConnectEndpoint, opConnectToken, opServiceAccountToken string,
	bwAuthMethod, bwSessionKey, bwApiUrl, bwAccessToken, bwClientId, bwProjectId string,
	kxDatabase, kxPassword, kxFieldName string,
) (string, error) {
	if provider == "" {
		provider = os.Getenv("TPMSFE_PROVIDER")
		if provider == "" {
			return "", fmt.Errorf("--provider flag or Environment variable TPMSFE_PROVIDER is required")
		}
	}
	if provider == "1password" && vault == "" || provider == "bitwarden" && vault == "" {
		vault = os.Getenv("TPMSFE_VAULT")
		if vault == "" {
			return "", fmt.Errorf("--vault flag or Environment variable TPMSFE_VAULT is required")
		}
	} else if kxDatabase == "" && provider == "kepassxc" {
		kxDatabase = os.Getenv("TPMSFE_KX_DATABASE")
		if kxDatabase == "" {
			return "", fmt.Errorf("--kxDatabase flag or Environment variable TPMSFE_KX_DATABASE is required")
		}
	}
	if title == "" {
		title = os.Getenv("TPMSFE_TITLE")
		if title == "" {
			return "", fmt.Errorf("--title flag or Environment variable TPMSFE_TITLE is required")
		}
	}

	switch provider {
	case "1password":
		return onepasswordmanager.CheckMethod(
			opAuthMethod,
			opAccountName,
			opConnectEndpoint,
			opConnectToken,
			opServiceAccountToken,
			vault,
			title,
		)
	case "bitwarden":
		return bitwardenmanager.CheckMethod(
			bwAuthMethod,
			bwSessionKey,
			bwApiUrl,
			vault,
			title,
			bwAccessToken,
			bwClientId,
			bwProjectId,
		)
	case "keepassxc":
		// For KeePassXC: vault is database path, title is entry path
		databasePath := vault
		if kxDatabase != "" {
			databasePath = kxDatabase
		}
		return keepassxcmanager.CheckMethod(
			databasePath,
			kxPassword,
			title,
			kxFieldName,
		)
	default:
		return "", fmt.Errorf("unsupported provider: %s. Supported providers are: 1password, bitwarden, keepassxc", provider)
	}
}
