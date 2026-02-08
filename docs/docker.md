# Docker

TPMSFE is built and published as several Docker image variants. All are built for `linux/amd64` and `linux/arm64`.

## Image variants

| Image tag | Dockerfile | Description |
|-----------|------------|-------------|
| `bubatzlegal/tpmsfe:minimal` | `minimal.Dockerfile` | Single binary (`tpmsfe`) in a distroless image. Use for building your custom docker container. |
| `bubatzlegal/tpmsfe:tofu` | `tofu.Dockerfile` | OpenTofu + tpmsfe in a Debian slim image. Default entrypoint is `tofu`. Suitable for CI or local runs that only need OpenTofu and tpmsfe. |
| `bubatzlegal/tpmsfe:custom` | `custom.Dockerfile` | OpenTofu + tpmsfe + jq, talosctl, packer. Broader tooling for custom pipelines. |
| `bubatzlegal/tpmsfe:cicd` | `cicd.Dockerfile` | OpenTofu + 1Password CLI (`op`) + Bitwarden CLI (`bw`) KeePassXC CLI (`kpcli`) + Minio client (`mc`) + Node (for `bw`), plus build tools. Intended for CI that tests multiple providers in this repository. |

## Using the minimal image

The minimal image only contains `/usr/local/bin/tpmsfe`. You must provide OpenTofu and the password manager integration yourself. Example with a volume and env:

```Dockerfile
FROM bubatzlegal/tpmsfe:minimal AS tpmsfe

FROM alpine:3.20

# Copy the tpmsfe binary from the minimal image
COPY --from=tpmsfe /usr/local/bin/tpmsfe /usr/local/bin/tpmsfe

# Add any other tools or dependencies you need
RUN apk add --no-cache git curl

# Your application setup
WORKDIR /workspace
```

In practice you would typically use this image to copy `tpmsfe` into an image that already has `tofu` and your other tools.

## Using the tofu image

The `tofu` image runs OpenTofu by default. Mount your Terraform/OpenTofu project and set env (or use a backend that does not require encryption to be available at plan time if needed):

```bash
docker run --rm -it \
  -e TPMSFE_PROVIDER=1password \
  -e TPMSFE_VAULT=MyVault \
  -e TPMSFE_TITLE=state-secret \
  -e OP_SERVICE_ACCOUNT_TOKEN="..." \
  -v "$(pwd)":/workdir -w /workdir \
  bubatzlegal/tpmsfe:tofu \
  tofu plan
```

Ensure `encrypt_command` and `decrypt_command` in your `.tf` files reference `tpmsfe` and that `tpmsfe` is on the image PATH (it is in the `tofu` and `custom` images).

## Building locally

From the repo root:

```bash
# Minimal (tpmsfe only)
docker build -f minimal.Dockerfile -t tpmsfe:minimal .

# OpenTofu + tpmsfe
docker build -f tofu.Dockerfile -t tpmsfe:tofu .

# Custom (tofu + jq + talosctl + packer)
docker build -f custom.Dockerfile -t tpmsfe:custom .

# CI (tofu + op + bw + mc + kpcli + Node)
docker build -f cicd.Dockerfile -t tpmsfe:cicd .
```

Build args:

- `tofu.Dockerfile`, `custom.Dockerfile`, `cicd.Dockerfile`: `UID`, `GID` (default `60000`) for the nonroot user.
