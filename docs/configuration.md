# Configuration

## Overview

TPMSFE acts as an **external encryption program** for OpenTofu (and Terraform) state and plan files. You configure it via the `terraform { encryption { ... } }` block: `encrypt_command` and `decrypt_command` both run the same `tpmsfe` binary with the same arguments. OpenTofu sends a JSON object on stdin with a `payload` (base64-encoded) and expects a JSON object on stdout with an (optionally re-encoded) `payload`.

- **Encryption:** The encryption key is derived from a secret fetched from your password manager (1Password, Bitwarden, or KeePassXC). Plaintext is encrypted with AES-256-GCM; the key is derived from that secret using PBKDF2 (100,000 iterations, SHA-256, 32-byte key).
- **Decryption:** The same secret is fetched; if the payload looks like previously encrypted data (format: 32-byte salt + 12-byte nonce + ciphertext), it is decrypted; otherwise it is treated as plaintext and encrypted (so the same command handles both directions).

You must use the same secret (same vault/item and credentials) for encrypt and decrypt so that state and plans remain readable across runs.

## Terraform / OpenTofu block

All providers follow this pattern:

```hcl
terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "<provider>",
        # ... provider-specific flags (vault, title, auth, etc.)
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "<provider>",
        # ... same flags as encrypt_command
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

Replace `<provider>` with one of: `1password`, `bitwarden`, `keepassxc`. The exact flags and any required environment variables depend on the provider; see:

- [1Password](providers/1password.md)
- [Bitwarden](providers/bitwarden.md)
- [KeePassXC](providers/keepassxc.md)

## Security notes

- Prefer environment variables for tokens and passwords instead of putting them in `.tf` or command lines.
- Ensure the password-manager entry used for the encryption secret is restricted and audited like other state secrets.
- The same secret must be available wherever OpenTofu/Terraform runs (CI, other machines); otherwise state cannot be decrypted.
