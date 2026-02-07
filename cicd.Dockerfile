FROM debian:trixie AS onepassword

ARG DEBIAN_FRONTEND=noninteractive
ARG TARGETARCH

RUN apt-get update && \
    apt-get install -y --no-install-recommends git ca-certificates wget unzip && \
    rm -rf /var/lib/apt/lists/*

RUN wget "https://cache.agilebits.com/dist/1P/op2/pkg/v2.32.1/op_linux_${TARGETARCH}_v2.32.1.zip" -O op.zip && \
    unzip -d op op.zip && \
    mv op/op /usr/local/bin

FROM alpine:latest AS opentofu

WORKDIR /app
RUN apk add --no-cache curl cosign

RUN curl --proto '=https' --tlsv1.2 -fsSL https://get.opentofu.org/install-opentofu.sh -o install-opentofu.sh && \
    chmod +x install-opentofu.sh && \
    ./install-opentofu.sh --install-method standalone

FROM node:trixie-slim

ARG UID=60000
ARG GID=60000

WORKDIR /workdir

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates golang kpcli && \
    rm -rf /var/lib/apt/lists/*

RUN npm install -g @bitwarden/cli

RUN groupadd -g ${GID} nonroot && \
    useradd -u ${UID} -g ${GID} -m -s /bin/bash nonroot

USER nonroot

COPY --from=opentofu /opt/opentofu/tofu /usr/local/bin/tofu
COPY --from=onepassword /usr/local/bin/op /usr/local/bin/op
COPY --from=minio/mc:latest /bin/mc /usr/local/bin/mc

ENTRYPOINT ["/usr/local/bin/tofu"]
