terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "cli",
        "--vault", "Test",
        "--title", "OpenTofu-Encryption-Key"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "cli",
        "--vault", "Test",
        "--title", "OpenTofu-Encryption-Key"
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

