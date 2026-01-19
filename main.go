package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/freethebaerchen/tpmsfe/bitwardenmanager"
	"github.com/freethebaerchen/tpmsfe/encryptor"
	"github.com/freethebaerchen/tpmsfe/onepasswordmanager"
)

func main() {
	var secret string
	var err error
	provider := flag.String("provider", "", "The provider you use as password manager. Possible values:\n - 1password\n - bitwarden\n - kepassxc.\n Default is none\n")
	vault := flag.String("vault", "", "The vault you want to use inside your password manager. Default is none")
	title := flag.String("title", "", "The title of the password entry you want to create. Default is none")
	direction := flag.String("direction", "encrypt", "The direction you want to use. Possible values:\n - encrypt: Encrypt the terraform/tofu state file\n - decrypt: the terraform/tofu state file\n Default is encrypt")
	statePath := flag.String("path", "./terraform.tfstate", "(optional) The path to the terraform/tofu state file you want to encrypt/decrypt. Default is current directory './terraform.tfstate'\n")
	authMethod := flag.String("op-auth", "desktop", "The authentication method for 1Password. Possible values:\n - desktop\n - service-account\n - connect\n Default is desktop")
	opAccountName := flag.String("op-account", "", "The 1Password account name. Required, when using 1Password Desktop.")
	opConnectEndpoint := flag.String("op-connect-endpoint", "http://127.0.0.1:8080", "The 1Password Connect endpoint. Required, when using 1Password Connect. Default is http://127.0.0.1:8080")
	opConnectToken := flag.String("op-connect-token", "", "The 1Password Connect token. Required, when using 1Password Connect.")
	opServiceAccountToken := flag.String("op-service-account-token", "", "The 1Password service account token. Required, when using 1Password service account authentication.")
	bwAuthMethod := flag.String("bw-auth", "cli", "The authentication method for Bitwarden. Possible values:\n - cli\n - api\n Default is cli")
	bwSessionKey := flag.String("bw-session", "", "Bitwarden CLI session key. Required when using CLI authentication.")
	bwApiUrl := flag.String("bw-api-url", "http://localhost", "The Bitwarden API URL. Default is http://localhost")
	bwAccessToken := flag.String("bw-access-token", "", "The Bitwarden access token. Required for Bitwarden authentication.")
	bwOrgId := flag.String("bw-org-id", "", "The Bitwarden organization ID. Required to access organization secrets.")
	bwProjectId := flag.String("bw-project-id", "", "The Bitwarden project ID. Required to access project secrets.")
	flag.Parse()
	if *provider == "" {
		fmt.Fprintln(os.Stderr, "Error: -provider flag is required")
		os.Exit(1)
	} else if *vault == "" {
		fmt.Fprintln(os.Stderr, "Error: -vault flag is required")
		os.Exit(1)
	} else if *title == "" {
		fmt.Fprintln(os.Stderr, "Error: -title flag is required")
		os.Exit(1)
	} else if *direction != "encrypt" && *direction != "decrypt" {
		fmt.Fprintln(os.Stderr, "Error: -direction flag must be either 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	switch *provider {
	case "1password":
		secret, err = onepasswordmanager.CheckMethod(
			*authMethod,
			*opAccountName,
			*opConnectEndpoint,
			*opConnectToken,
			*opServiceAccountToken,
			*vault,
			*title,
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching item from 1Password:", err)
			os.Exit(1)
		}
	case "bitwarden":
		secret, err = bitwardenmanager.CheckMethod(
			*bwAuthMethod,
			*bwSessionKey,
			*bwApiUrl,
			*vault,
			*title,
			*bwAccessToken,
			*bwOrgId,
			*bwProjectId,
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching item from Bitwarden:", err)
			os.Exit(1)
		}
	// case "keepassxc":
	//  handleKeePass(*vault, *title, *direction, *statePath)
	default:
		fmt.Fprintln(os.Stderr, "Error: Unsupported provider. Supported providers are: 1password, bitwarden, keepassxc")
		os.Exit(1)
	}
	switch *direction {
	case "decrypt":
		err = encryptor.DecryptFileContent(secret, *statePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error decrypting file:", err)
			os.Exit(1)
		}
	case "encrypt":
		err = encryptor.EncryptFileContent(secret, *statePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error encrypting file:", err)
			os.Exit(1)
		}
	}
}
