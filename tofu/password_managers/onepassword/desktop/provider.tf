terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "desktop",
        "--vault", "kube-hetzner",
        "--title", "terraform-state-secret",
        "--op-account", "Unsere Familie"
      ]

      decrypt_command = [
        "tpmsfe",
        "--provider", "1password",
        "--op-auth", "desktop",
        "--vault", "kube-hetzner",
        "--title", "terraform-state-secret",
        "--op-account", "Unsere Familie"
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
