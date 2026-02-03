FROM alpine:latest AS opentofu

WORKDIR /app
RUN apk add --no-cache curl cosign

RUN curl --proto '=https' --tlsv1.2 -fsSL https://get.opentofu.org/install-opentofu.sh -o install-opentofu.sh && \
    chmod +x install-opentofu.sh && \
    ./install-opentofu.sh --install-method standalone

FROM debian:trixie AS certificates

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

FROM golang:tip-trixie AS tpmsfe

ARG TARGETARCH
ARG TARGETOS

RUN apt-get update && apt-get install -y gcc libc6-dev

WORKDIR /app
COPY . .

ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-lm"

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o tpmsfe .

FROM debian:stable-slim

ARG UID=60000
ARG GID=60000

WORKDIR /workdir

RUN groupadd -g ${GID} nonroot && \
    useradd -u ${UID} -g ${GID} -m -s /bin/bash nonroot

USER nonroot

COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=jq /usr/local/bin/jq /usr/local/bin/jq
COPY --from=talosctl /usr/local/bin/talosctl /usr/local/bin/talosctl
COPY --from=packer /usr/bin/packer /usr/local/bin/packer
COPY --from=tpmsfe /app/tpmsfe /usr/local/bin/tpmsfe
COPY --from=opentofu /opt/opentofu/tofu /usr/local/bin/tofu

ENTRYPOINT ["/usr/local/bin/tofu"]
