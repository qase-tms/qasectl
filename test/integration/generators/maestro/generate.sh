#!/usr/bin/env bash
# Usage: generate.sh <out-dir>; prints produced scenarios to stdout, tool output goes to stderr.
# Simulators: FIXTURE_DEVICE (default "iPhone 16 Pro") and FIXTURE_SHARD_DEVICE (default "iPhone 16"), name or UDID.
set -uo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
out="$1"
device="${FIXTURE_DEVICE:-iPhone 16 Pro}"
shard_device="${FIXTURE_SHARD_DEVICE:-iPhone 16}"

command -v maestro >/dev/null 2>&1 || export PATH="$HOME/.maestro/bin:$PATH"
export MAESTRO_CLI_NO_ANALYTICS=1
export MAESTRO_CLI_ANALYSIS_NOTIFICATION_DISABLED=true

udid="$("$here/setup.sh" "$device")" || exit 1

run="$out/maestro/run/input"
mkdir -p "$run" "$out/junit/maestro/input"
# Negative flows fail on purpose, so Maestro's exit code is ignored.
(cd "$here" && maestro test --device "$udid" \
  --debug-output "$run" --flatten-debug-output \
  --format junit --output "$run/report.xml" .maestro/ >&2) || true
[ -f "$run/report.xml" ] || { echo "maestro produced no $run/report.xml" >&2; exit 1; }
cp "$run/report.xml" "$out/junit/maestro/input/report.xml"
echo "maestro/run"
echo "junit/maestro"

shard_udid="$("$here/setup.sh" "$shard_device")" || exit 1
shards="$out/maestro/shards/input"
mkdir -p "$shards"
(cd "$here" && maestro test --device "$udid,$shard_udid" --shard-all 2 \
  --debug-output "$shards" --flatten-debug-output \
  --format junit --output "$shards/report.xml" .maestro/flows/smoke/01-launch.yaml >&2) || true
xcrun simctl shutdown "$shard_udid" >&2 || true
[ -f "$shards/report.xml" ] || { echo "maestro produced no $shards/report.xml" >&2; exit 1; }

# The xcodebuild runner log is pure machine noise and qasectl never reads it; other logs carry temp-dir paths.
# Simulator UDIDs stay: commands.json and report.xml must agree on them for the device mapping.
rm -f "$run"/xctest_runner_*.log "$shards"/xctest_runner_*.log
find "$out" -type f -name '*.log' | while IFS= read -r f; do
  LC_ALL=C sed -i.bak -E -e "s|$out|/tmp/out|g" -e 's|/var/folders/[^ ]*|/tmp|g' "$f"
  rm -f "$f.bak"
done
echo "maestro/shards"
