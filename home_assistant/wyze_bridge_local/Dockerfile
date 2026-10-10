# Builds the bridge from source at a pinned commit (the upstream add-on
# copies prebuilt binaries from a :latest image instead). go2rtc and its
# player script are checksum-verified.
ARG BUILD_FROM=ghcr.io/home-assistant/base:3.21

FROM alpine:3.21 AS go2rtc-fetch
ARG BUILD_ARCH
ARG GO2RTC_VERSION=1.9.14
ARG GO2RTC_SHA256_aarch64=359fabade8a7a51e81a55fe6df6b0ef81764a5e1d63179577534eaaa71904b50
ARG GO2RTC_SHA256_amd64=32d616af226bd731678ffde328b94cfb94e30339bfefc469cfb76323144615a6
ARG VIDEO_RTC_SHA256=d48ce627baf7c341a92c0f5844a3c546431f9db873ff21489671aba2ecfe64fb
RUN apk add --no-cache curl ca-certificates && \
    case "${BUILD_ARCH}" in \
      aarch64) arch=arm64; sum="${GO2RTC_SHA256_aarch64}" ;; \
      amd64)   arch=amd64; sum="${GO2RTC_SHA256_amd64}" ;; \
      *) echo "unsupported arch: ${BUILD_ARCH}" >&2; exit 1 ;; \
    esac && \
    curl -fsSL -o /go2rtc \
      "https://github.com/AlexxIT/go2rtc/releases/download/v${GO2RTC_VERSION}/go2rtc_linux_${arch}" && \
    echo "${sum}  /go2rtc" | sha256sum -c - && chmod +x /go2rtc && \
    curl -fsSL -o /video-rtc.js \
      "https://github.com/AlexxIT/go2rtc/raw/v${GO2RTC_VERSION}/www/video-rtc.js" && \
    echo "${VIDEO_RTC_SHA256}  /video-rtc.js" | sha256sum -c -

FROM golang:1.26.4-alpine AS builder
ARG BRIDGE_REPO=https://github.com/ssw/docker-wyze-bridge
ARG BRIDGE_COMMIT=a7518b343e001287c587a1279a599c3af954991d
ARG BRIDGE_VERSION=4.5.0-local.2
RUN apk add --no-cache git ca-certificates
WORKDIR /src
RUN git init -q && git remote add origin "${BRIDGE_REPO}" && \
    git fetch -q --depth 1 origin "${BRIDGE_COMMIT}" && git checkout -q FETCH_HEAD && \
    test "$(git rev-parse HEAD)" = "${BRIDGE_COMMIT}"
RUN go mod download && go mod verify
COPY --from=go2rtc-fetch /video-rtc.js ./internal/webui/static/video-rtc.js
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=${BRIDGE_VERSION}" \
    -o /wyze-bridge ./cmd/wyze-bridge

FROM ${BUILD_FROM}
# ffmpeg is only needed by go2rtc's JPEG snapshot endpoint.
RUN apk add --no-cache ca-certificates tzdata ffmpeg && \
    mkdir -p /config /media/wyze_bridge/snapshots /media/wyze_bridge/recordings
COPY --from=builder      /wyze-bridge /usr/local/bin/wyze-bridge
COPY --from=go2rtc-fetch /go2rtc      /usr/local/bin/go2rtc
COPY run.sh /
RUN chmod +x /run.sh
CMD [ "/run.sh" ]
