#!/usr/bin/env bash
set -e

APP_NAME="$1"
DIST="$2"
VERSION="${3:-dev}"

PLATFORMS=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
)

mkdir -p "$DIST"

for platform in "${PLATFORMS[@]}"; do
  set -- $platform
  GOOS=$1
  GOARCH=$2
  echo "  -> Building for $GOOS/$GOARCH"
  CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH \
    go build -ldflags="-s -w -X main.version=$VERSION" \
    -o "$DIST/$APP_NAME-$GOOS-$GOARCH" ./cmd/$APP_NAME
done

