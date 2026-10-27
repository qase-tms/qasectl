#!/usr/bin/env bash
# Usage: generate.sh <out-dir>; prints produced scenarios to stdout, tool output goes to stderr.
set -uo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
out="$1"

if [ ! -x "$here/.venv/bin/pytest" ]; then
  python3 -m venv "$here/.venv" >&2 || exit 1
  "$here/.venv/bin/pip" install -q -r "$here/requirements.txt" >&2 || exit 1
fi

mkdir -p "$out/junit/pytest/input"
# Failures and errors are the point of this fixture, so pytest's exit code is ignored.
(cd "$here" && .venv/bin/pytest tests --junitxml="$out/junit/pytest/input/report.xml" -o junit_logging=all -p no:cacheprovider >&2) || true

# The host name is machine-specific.
sed -i.bak -E 's/ hostname="[^"]*"/ hostname="localhost"/' "$out/junit/pytest/input/report.xml"
rm -f "$out/junit/pytest/input/report.xml.bak"

echo "junit/pytest"
