#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 vYYYY.MM.DD[-X]"
  exit 1
fi

git commit --allow-empty -m "$VERSION release"
