terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "api",
        "--bw-api-url", "https://bitwarden.jokester.cloud",
        "--bw-client-id", "user.b067c26c-d221-4340-9b53-4d2af7c629fd",
        "--vault", "user.b067c26c-d221-4340-9b53-4d2af7c629fd",
        "--title", "OpenTofu-Encryption-Key"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "api",
        "--bw-api-url", "https://bitwarden.jokester.cloud",
        "--bw-client-id", "user.b067c26c-d221-4340-9b53-4d2af7c629fd",
        "--vault", "user.b067c26c-d221-4340-9b53-4d2af7c629fd",
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
