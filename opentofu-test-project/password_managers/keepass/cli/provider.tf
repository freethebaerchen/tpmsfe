terraform {
  encryption {
    method "external" "tpmsfe" {
      encrypt_command = [
        "tpmsfe",
        "--provider", "keepassxc",
        "--title", "Your Item",
      ]
      
      decrypt_command = [
        "tpmsfe",
        "--provider", "keepassxc",
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