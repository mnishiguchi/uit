#!/usr/bin/env bash
set -e

# build.sh — Build binaries for multiple platforms
#
# Usage:
#   ./scripts/build.sh <app_name> <output_dir> [version]
#
# Example:
#   ./scripts/build.sh uit dist v2025.03.28
#
# Arguments:
#   <app_name>   The name of the binary (e.g., "uit")
#   <output_dir> Output directory (e.g., "dist")
#   [version]    Optional version string (default: "dev")

APP_NAME="${1:-uit}"
DIST="${2:-dist}"
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

