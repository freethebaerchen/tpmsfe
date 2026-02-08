# Bitwarden

## Overview

TPMSFE can use Bitwarden (or a compatible server such as Vaultwarden) as the source for the encryption secret. Supported authentication methods:

- **cli** – Bitwarden CLI (`bw`). You must unlock the vault and provide a session key.
- **api** – Bitwarden API with client ID and access token. Note: the API path is intended for Bitwarden's own API; self-hosted/Vaultwarden may only support the client API. This mode is not tested.

You must have an item or secret whose title matches `--title`. The same title must be used for both `encrypt_command` and `decrypt_command`.

## Requirements

- Bitwarden account (or self-hosted/Vaultwarden instance).
- For **CLI:** Bitwarden CLI (`bw`) installed. For self-hosted: `bw config server https://your-bitwarden.example`. Log in with `bw login`, then unlock and export the session (export command is shown after login; with `--raw` only the session token is shown). Run `bw sync`.
- For **API:** Organization client ID and an API access token.

## CLI authentication

Uses the Bitwarden CLI. Provide the session key via environment variable or flag.

### Terraform snippet

```hcl
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "cli",
        "--title", "Your Item"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "cli",
        "--title", "Your Item"
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

### Options

- Session key: `TPMSFE_BW_SESSION` or `BW_SESSION` (preferred over `--bw-session`).

## API authentication

Uses the Bitwarden SDK/API with client ID and access token. For self-hosted or Bitwarden EU, set the API URL with `--bw-api-url` (default: `https://bitwarden.com`). The client ID is the organization ID used to list and fetch secrets. Prefer environment variables for the access token.

### Terraform snippet

```hcl
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "api",
        "--bw-client-id", "Your Client ID",
        "--title", "Your Item"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "api",
        "--bw-client-id", "Your Client ID",
        "--title", "Your Item"
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

### Options

- API URL: `--bw-api-url` or `TPMSFE_BW_API_URL` or `BW_API_URL`. Default is `https://bitwarden.com`.
- Access token: `TPMSFE_BW_ACCESS_TOKEN` or `BW_ACCESS_TOKEN` (preferred over `--bw-access-token`).
- Project ID: `--bw-project-id` or `TPMSFE_BW_PROJECT_ID` or `BW_PROJECT_ID` (optional).

## Environment variables and flags

| Environment variable(s) | Flag |
|-------------------------|------|
| `TPMSFE_PROVIDER` | `--provider` |
| `TPMSFE_VAULT` | `--vault` |
| `TPMSFE_TITLE` | `--title` |
| `TPMSFE_BW_AUTH`, `BW_AUTH` | `--bw-auth` |
| `TPMSFE_BW_SESSION`, `BW_SESSION` | `--bw-session` |
| `TPMSFE_BW_API_URL`, `BW_API_URL` | `--bw-api-url` |
| `TPMSFE_BW_ACCESS_TOKEN`, `BW_ACCESS_TOKEN` | `--bw-access-token` |
| `TPMSFE_BW_CLIENT_ID`, `BW_CLIENT_ID` | `--bw-client-id` |
| `TPMSFE_BW_PROJECT_ID`, `BW_PROJECT_ID` | `--bw-project-id` |

When an environment variable is set, it overrides the corresponding flag. Default for `--bw-auth` is `cli`; default for `--bw-api-url` is `https://bitwarden.com`.
