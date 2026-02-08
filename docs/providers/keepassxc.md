# KeePassXC

## Overview

TPMSFE can use a KeePassXC database (`.kdbx`) as the source for the encryption secret. It uses the KeePassXC CLI to read a single entry's field (by default the "Password" field).

You must have a `.kdbx` file and an entry whose title (or path) matches `--title`. The same database path and title must be used for both `encrypt_command` and `decrypt_command`.

## Requirements

- KeePassXC CLI tools available on the system (e.g. `kpcli`).
- Path to the `.kdbx` database file and the master password (or keyfile).

## Configuration

Create an entry in your KeePassXC database to hold the encryption secret. Reference it by title in the Terraform block. Prefer environment variables for the database path and password in shared or CI environments.

### Terraform snippet

```hcl
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "keepassxc",
        "--kx-database", "/path/to/database.kdbx",
        "--title", "Your Item"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "keepassxc",
        "--kx-database", "/path/to/database.kdbx",
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

- Database path: `--kx-database` or `TPMSFE_KX_DATABASE` or `TPMSFE_KEEPASSXC_DATABASE` or `KEEPASSXC_DATABASE`.
- Password: `--kx-password` or `TPMSFE_KEEPASSXC_PASSWORD` or `KEEPASSXC_PASSWORD`.
- Field name: `--kx-field` (default: `Password`). No environment variable; must be set in the encrypt/decrypt command if changed.

## Environment variables and flags

| Environment variable(s) | Flag |
|-------------------------|------|
| `TPMSFE_PROVIDER` | `--provider` |
| `TPMSFE_TITLE` | `--title` |
| `TPMSFE_KX_DATABASE`, `TPMSFE_KEEPASSXC_DATABASE`, `KEEPASSXC_DATABASE` | `--kx-database` |
| `TPMSFE_KEEPASSXC_PASSWORD`, `KEEPASSXC_PASSWORD` | `--kx-password` |

Flags with no environment variable: `--kx-field`.
