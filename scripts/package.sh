#!/usr/bin/env bash
set -euo pipefail

APP_NAME=${1:-uit}
DIST=dist
ARCHIVES=$DIST/archives
PLATFORMS=("linux-amd64" "linux-arm64" "darwin-amd64" "darwin-arm64")

mkdir -p "$ARCHIVES"

for platform in "${PLATFORMS[@]}"; do
  binary="$DIST/$APP_NAME-$platform"

  if [[ ! -f "$binary" ]]; then
    echo "!! Skipping $platform — binary not found"
    continue
  fi

  echo "-> Packaging $platform"
  mkdir -p "$DIST/tmp/$platform"
  cp "$binary" "$DIST/tmp/$platform/$APP_NAME"
  cp README.md "$DIST/tmp/$platform/README.md"
  tar -czf "$ARCHIVES/$APP_NAME-$platform.tar.gz" -C "$DIST/tmp/$platform" "$APP_NAME" README.md
  rm -rf "$DIST/tmp/$platform"
done
