terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "keepassxc",
        "--title", "terraform-state-secret",
        "--kx-database", "/kp/Passwords.kdbx"
      ]
      
      decrypt_command = [
        "tpmsfe",
        "--provider", "keepassxc",
        "--title", "terraform-state-secret",
        "--kx-database", "/kp/Passwords.kdbx"
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