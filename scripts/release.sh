#!/usr/bin/env bash
set -euo pipefail

# release.sh — Tag and push a version to trigger GitHub Release
#
# Usage:
#   ./scripts/release.sh vYYYY.MM.DD
#
# This script:
# - Ensures you’re on main branch
# - Ensures working tree is clean
# - Shows changelog since last tag
# - Tags and pushes the release tag

VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 vYYYY.MM.DD[-X]"
  exit 1
fi

branch=$(git rev-parse --abbrev-ref HEAD)
if [[ "$branch" != "main" ]]; then
  echo "❌ You must be on 'main' branch (currently on '$branch')"
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "❌ Working directory is not clean. Commit or stash changes first."
  exit 1
fi

echo "📋 Changelog since last tag:"
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [[ -n "$LAST_TAG" ]]; then
  git log "$LAST_TAG"..HEAD --pretty=format:"- %s (%h)"
else
  git log --pretty=format:"- %s (%h)"
fi

echo
echo "🏷️ Tagging version: $VERSION"
git tag "$VERSION"
git push origin main --tags

echo
echo "🚀 Release triggered! GitHub Actions will publish binaries."

