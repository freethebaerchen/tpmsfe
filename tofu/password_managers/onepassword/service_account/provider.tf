terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "service-account",
        "--vault", "kube-hetzner",
        "--title", "terraform-state-secret"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "service-account",
        "--vault", "kube-hetzner",
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
