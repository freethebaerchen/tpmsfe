# Installation

## Requirements

- **Go 1.25 or later** (for building from source). The project uses `CGO_ENABLED=1` and links against `libc` (for KeePassXC CLI via `kpcli`).
- **OpenTofu or Terraform** with support for external encryption methods (e.g. OpenTofu 1.8+).
- A supported **password manager**: 1Password, Bitwarden, or KeePassXC. Provider-specific CLI or API tools may be required; see the provider docs.

## Build from source

Clone the repository and build:

```bash
git clone https://github.com/freethebaerchen/tpmsfe.git
cd tpmsfe
CGO_ENABLED=1 go build -o tpmsfe .
```

Ensure `tpmsfe` is on your `PATH` or use an absolute path in your Terraform/OpenTofu `encrypt_command` and `decrypt_command`.

## Docker

Pre-built images are published for multiple Dockerfile variants. See [Docker](docker.md) for image names and usage. You can also build locally:

```bash
docker build -f minimal.Dockerfile -t tpmsfe:minimal .
```

## Verification

After installation, OpenTofu will invoke the binary and expect a JSON line with a magic header on stdout. You can sanity-check the binary by running:

```bash
echo '{"payload":"'$(echo -n "test" | base64)'"}' | tpmsfe --provider keepassxc --title "Your Item" ...
```

(Replace `...` with the flags required for your provider; see [Configuration](configuration.md) and the provider docs.)
