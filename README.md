# (Open)Tofu Password Manager State File Encryption (TPMSFE)

## 1Password
### Terraform snippet for configuration (Desktop):
```tf
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account name"
      ]
      
      decrypt_command = [
        "tpmsfe",
        "--provider", "1password",
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

### Terraform snippet for configuration (Service Account):
```tf
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "service-account"
        // Not recommended:
        // "--op-service-account-token" "your-sa-token"
        // Consider using either of the environment variables TPMSFE_OP_SERVICE_ACCOUNT_TOKEN or OP_SERVICE_ACCOUNT_TOKEN
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account name"
      ]
      
      decrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "service-account"
        // Not recommended:
        // "--op-service-account-token" "your-sa-token"
        // Consider using either of the environment variables TPMSFE_OP_SERVICE_ACCOUNT_TOKEN or OP_SERVICE_ACCOUNT_TOKEN
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

### Terraform snippet for configuration (Service Account):
```tf
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "connect"
        // Not recommended:
        // "--op-connect-token" "your-connect-token"
        // Consider using either of the environment variables TPMSFE_OP_CONNECT_TOKEN or OP_CONNECT_TOKEN
        "--vault", "Your Vault",
        "--title", "Your Item",
        "--op-account", "Your Company/Your Account name"
      ]
      
      decrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "connect"
        // Not recommended:
        // "--op-connect-token" "your-connect-token"
        // Consider using either of the environment variables TPMSFE_OP_CONNECT_TOKEN or OP_CONNECT_TOKEN
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

## BitWarden
For CLI authentication, you first need to run `bw config (server https://bitwarden.self.hosted)`

Then run `bw login` and copy the export command.

You can now `export BW_SESSION="Your_Session"`, `export TPMSFE_BW_SESSION="Your_Session"` or provide the session via the `--bw-session` but this is not recommended.

### Terraform snippet for configuration (CLI):
```tf
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "cli",
        "--vault", "Your Vault",
        "--title", "Your Item"
      ]
      
      decrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "cli",
        "--vault", "Your Vault",
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

### Terraform snippet for configuration (API):
```tf
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "api",
        // For self-hosted or bitwarden.eu
        // "--bw-api-url", "https://bitwarden.eu",
        "--bw-client-id", "Your Client ID",
        // The client ID can als be provided with the environment variables TPMSFE_BW_CLIENT_ID or BW_CLIENT_ID
        "--vault", "Your Vault",
        "--title", "Your Item"
      ]
      
      decrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "api",
        // For self-hosted or bitwarden.eu
        // "--bw-api-url", "https://bitwarden.eu",
        "--bw-client-id", "Your Client ID",
        // The client ID can als be provided with the environment variables TPMSFE_BW_CLIENT_ID or BW_CLIENT_ID
        "--vault", "Your Vault",
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
## KeePassXC
### Terraform snippet for configuration:
```tf
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "keepassxc",
        // Not recommended:
        // "--kx-password" "your-database-password"
        // Consider using either of the environment variables TPMSFE_KX_PASSWORD or KX_PASSWORD
        "--title", "Your Item",
      ]
      
      decrypt_command = [
        "tpmsfe",
        "--provider", "keepassxc",
        // Not recommended:
        // "--kx-password" "your-database-password"
        // Consider using either of the environment variables TPMSFE_KX_PASSWORD or KX_PASSWORD
        "--title", "Your Item",
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