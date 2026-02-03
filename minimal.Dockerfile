FROM golang:tip-trixie AS tpmsfe

ARG TARGETARCH
ARG TARGETOS

RUN apt-get update && apt-get install -y gcc libc6-dev

WORKDIR /app
COPY . .

ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-lm"

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o tpmsfe .

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=tpmsfe /app/tpmsfe /usr/local/bin/tpmsfe
