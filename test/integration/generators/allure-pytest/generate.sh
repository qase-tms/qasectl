#!/usr/bin/env bash
# Usage: generate.sh <out-dir>; prints produced scenarios to stdout, tool output goes to stderr.
set -uo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
out="$1"

if [ ! -x "$here/.venv/bin/pytest" ]; then
  python3 -m venv "$here/.venv" >&2 || exit 1
  "$here/.venv/bin/pip" install -q -r "$here/requirements.txt" >&2 || exit 1
fi

input="$out/allure/pytest/input"
mkdir -p "$input"
# Failures and errors are the point of this fixture, so pytest's exit code is ignored.
(cd "$here" && .venv/bin/pytest tests --alluredir="$input" -p no:cacheprovider >&2) || true

# host and thread labels are machine-specific.
for f in "$input"/*-result.json; do
  sed -i.bak -E \
    -e 's/("name": ?"host", ?"value": ?")[^"]*"/\1localhost"/' \
    -e 's/("name": ?"thread", ?"value": ?")[^"]*"/\1main"/' "$f"
  rm -f "$f.bak"
done

echo "allure/pytest"
