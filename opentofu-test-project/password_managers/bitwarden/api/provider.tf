terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "api",
        "--bw-api-url", "",
        "--bw-client-id", "",
        "--vault", "TPMSFE-Test",
        "--title", "terraform-state-secret"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "bitwarden",
        "--bw-auth", "api",
        "--bw-api-url", "",
        "--bw-client-id", "",
        "--vault", "TPMSFE-Test",
        "--title", "terraform-state-secret"
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
