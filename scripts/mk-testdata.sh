#!/usr/bin/env bash
set -euo pipefail

# mk-testdata.sh — Generate test input directories for uit CLI tests
#
# Usage:
#   ./scripts/mk-testdata.sh
#
# This creates test input data under internal/cli/testdata/input/

ROOT="internal/cli/testdata/input"
CASES=(
  "default"
  "max-lines"
  "no-tree"
  "no-content"
  "filter"
  "copy"
  "binary"
)

write_text_file() {
  local file="$1"
  cat <<EOF >"$file"
line 1
line 2
line 3
line 4
line 5
EOF
}

echo "🔧 Generating test input directories under $ROOT..."

for case in "${CASES[@]}"; do
  dir="$ROOT/$case"
  mkdir -p "$dir"
  echo "📁 Created $dir"

  mkdir -p "$dir/sub"

  case "$case" in
  "filter")
    write_text_file "$dir/a.txt"
    write_text_file "$dir/b.txt"
    write_text_file "$dir/sub/c.txt"
    ;;
  "binary")
    printf '\x00This is binary data' >"$dir/a.txt"
    write_text_file "$dir/b.txt"
    write_text_file "$dir/sub/c.txt"
    ;;
  *)
    write_text_file "$dir/a.txt"
    write_text_file "$dir/b.txt"
    write_text_file "$dir/sub/c.txt"
    ;;
  esac
done

echo "✅ Done generating test inputs."
