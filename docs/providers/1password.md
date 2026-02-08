# 1Password

## Overview

TPMSFE can use 1Password as the source for the encryption secret. Supported authentication methods:

- **desktop** – 1Password desktop app (default). Requires `op` CLI and an active desktop session.
- **connect** – 1Password Connect server. Requires a Connect API token and endpoint.
- **service-account** – 1Password service account. Requires a service account secret token.

You must create a vault and an item (e.g. a Secure Note or password item) that will hold the secret used to derive the encryption key. The same vault and item title must be used for both `encrypt_command` and `decrypt_command`.

## Requirements

- 1Password account with a vault and an item to store the encryption secret.
- For **desktop:** 1Password desktop application and `op` CLI with SDK integration enabled ([1Password documentation](https://developer.1password.com/docs/sdks/desktop-app-integrations/)).
- For **connect:** 1Password Connect server and an API token.
- For **service-account:** A 1Password service account and its secret token.

## Desktop

Uses the 1Password desktop app. The account name is usually of the form `Your Company/Your Account`.

### Terraform snippet

```hcl
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "desktop",
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account name"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "desktop",
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account name"
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

- Account name: `--op-account` or `TPMSFE_OP_ACCOUNT_NAME` or `OP_ACCOUNT_NAME`.

## Service account

For headless or CI use. Create a service account in 1Password and use its secret token. Prefer environment variables over putting the token in `.tf` or the command line.

### Terraform snippet

```hcl
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "service-account",
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account name"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "service-account",
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account name"
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

- Service account token: `TPMSFE_OP_SERVICE_ACCOUNT_TOKEN` or `OP_SERVICE_ACCOUNT_TOKEN` (preferred over `--op-service-account-token`).

## Connect

Uses a 1Password Connect server (self-hosted). You need a Connect API token and the server URL.

### Terraform snippet

```hcl
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "connect",
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account name"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "connect",
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account name"
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

- Connect endpoint: `--op-connect-endpoint` or `TPMSFE_OP_CONNECT_HOST` or `OP_CONNECT_HOST`. Default is `http://localhost:8080`.
- Connect token: `TPMSFE_OP_CONNECT_TOKEN` or `OP_CONNECT_TOKEN` (preferred over `--op-connect-token`).

## Environment variables and flags

| Environment variable(s) | Flag |
|-------------------------|------|
| `TPMSFE_PROVIDER` | `--provider` |
| `TPMSFE_VAULT` | `--vault` |
| `TPMSFE_TITLE` | `--title` |
| `TPMSFE_OP_ACCOUNT_NAME`, `OP_ACCOUNT_NAME` | `--op-account` |
| `TPMSFE_OP_CONNECT_TOKEN`, `OP_CONNECT_TOKEN` | `--op-connect-token` |
| `TPMSFE_OP_CONNECT_HOST`, `OP_CONNECT_HOST` | `--op-connect-endpoint` |
| `TPMSFE_OP_SERVICE_ACCOUNT_TOKEN`, `OP_SERVICE_ACCOUNT_TOKEN` | `--op-service-account-token` |
