#!/usr/bin/env bash
set -euo pipefail

# empty-release-commit.sh — Create an empty commit for a release
#
# Usage:
#   ./scripts/empty-release-commit.sh vYYYY.MM.DD
#
# Creates a placeholder commit for marking a release point.

VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 vYYYY.MM.DD[-X]"
  exit 1
fi

git commit --allow-empty -m "$VERSION release"

