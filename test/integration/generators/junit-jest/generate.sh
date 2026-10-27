#!/usr/bin/env bash
# Usage: generate.sh <out-dir>; prints produced scenarios to stdout, tool output goes to stderr.
set -uo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
out="$1"

[ -d "$here/node_modules" ] || (cd "$here" && npm ci --silent >&2) || exit 1

mkdir -p "$out/junit/jest/input"
# The failing test is intentional, so jest's exit code is ignored.
(cd "$here" && JEST_JUNIT_OUTPUT_DIR="$out/junit/jest/input" npx jest >&2) || true

echo "junit/jest"
