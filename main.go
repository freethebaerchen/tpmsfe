package main

import (
	"flag"
	"fmt"
	"os"

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
	authMethod := flag.String("auth", "desktop", "The authentication method for 1Password. Possible values:\n - desktop\n - service-account\n - connect\n Default is desktop")
	opAccountName := flag.String("account", "", " The 1Password account name. Required, when using 1Password Desktop.")
	opConnectEndpoint := flag.String("connect-endpoint", "http://127.0.0.1:8080", " The 1Password Connect endpoint. Required, when using 1Password Connect. Default is http://127.0.0.1:8080")
	opConnectToken := flag.String("connect-token", "", " The 1Password Connect token. Required, when using 1Password Connect.")
	opServiceAccountToken := flag.String("service-account-token", "", " The 1Password service account token. Required, when using 1Password service account authentication.")
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
	// case "bitwarden":
	//  handleBitwarden(*vault, *title, *direction, *statePath)
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
