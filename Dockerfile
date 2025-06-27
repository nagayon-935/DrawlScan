FROM golang:1.24.3-bullseye AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y libpcap-dev

ARG CGO_ENABLED=1
ENV CGO_ENABLED=${CGO_ENABLED}

COPY . .
RUN go mod download
RUN go build -o drawlscan ./cmd/main/drawlscan.go ./cmd/main/version.go

FROM debian:stable-slim

ARG VERSION=0.3.5

RUN apt-get update && apt-get install -y libpcap-dev

# メタデータ（OCIラベル）
LABEL org.opencontainers.image.source="https://github.com/nagayon-935/DrawlScan" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.title="DrawlScan" \
      org.opencontainers.image.description="CLI-based network scanner with ASCII and GeoIP rendering"

# 実行用ユーザと作業ディレクトリ
RUN mkdir -p /workdir

# バイナリとCA証明書の配置
COPY --from=builder /app/drawlscan /opt/drawlscan/drawlscan

WORKDIR /workdir
