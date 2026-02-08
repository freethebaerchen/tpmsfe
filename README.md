# TPMSFE – (Open)Tofu Password Manager State File Encryption

TPMSFE is an external encryption helper for **OpenTofu** and **Terraform** state and plan files. It uses a secret stored in a password manager (1Password, Bitwarden, or KeePassXC) to derive an encryption key and encrypt/decrypt payloads with AES-256-GCM. OpenTofu invokes the same binary for both encrypt and decrypt via the `encryption { method "external" ... }` block.

## Features

- Works with OpenTofu (and Terraform) external encryption API: single binary for both encrypt and decrypt.
- **1Password:** Desktop app, Connect server, or Service Account.
- **Bitwarden:** CLI session, compatible with self-hosted/Vaultwarden for CLI. API is, as of right now, unsupported, since Vaultwarden doesn't support the Client API.
- **KeePassXC:** Local `.kdbx` database; secret read via KeePassXC CLI.
- Encryption: AES-256-GCM; key derived with PBKDF2 (100,000 iterations, SHA-256).
- Configuration via command-line flags and/or environment variables; no secrets in `.tf` when using env vars.

## Quick start

1. **Install** the `tpmsfe` binary (see [Installation](docs/installation.md)).
2. **Create a secret** in your password manager (e.g. a Secure Note in 1Password, an item in Bitwarden, or an entry in KeePassXC). This value will be the encryption key material.
3. **Add the encryption block** to your Terraform/OpenTofu root module. Example for 1Password desktop:

```hcl
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account"
      ]
    }

    state {
      method   = method.external.tpmsfe
      enforced = true
    }

    plan {
      method   = method.external.tpmsfe
      enforced = true
    }
  }
}
```

4. Run `tofu init`, `tofu plan`, `tofu apply` as usual. State and plan files are encrypted using the secret from your password manager.

## Documentation

Full documentation is in the [docs](docs/README.md) directory:

| Document | Description |
|----------|-------------|
| [Installation](docs/installation.md) | Build from source, Docker images |
| [Configuration](docs/configuration.md) | Encryption block, flags, how encryption works |
| [Docker](docs/docker.md) | Image variants and usage |
| [1Password](docs/providers/1password.md) | Desktop, Connect, Service Account |
| [Bitwarden](docs/providers/bitwarden.md) | CLI (and API) |
| [KeePassXC](docs/providers/keepassxc.md) | Database path, password, field |

## Requirements

- **OpenTofu** or **Terraform** with support for external encryption (e.g. OpenTofu 1.8+).
- **Go 1.25+** for building from source; `CGO_ENABLED=1` for KeePassXC support.
- One of: **1Password** (CLI/Connect/Service Account), **Bitwarden** (CLI or API), or **KeePassXC** (CLI and `.kdbx`).
