#!/usr/bin/env bash
# Usage: setup.sh [device name or UDID] — boots an iOS simulator and installs the Wikipedia sample app.
# Prints the simulator UDID on stdout; everything else goes to stderr. Safe to run repeatedly.
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
device="${1:-iPhone 16 Pro}"
cache="$here/.cache"
app="$cache/app/Wikipedia.app"

if [ ! -d "$app" ]; then
  mkdir -p "$cache"
  if [ ! -f "$cache/samples.zip" ]; then
    # The archive behind `maestro download-samples`; its built-in HTTP timeout is too short for ~300 MB.
    echo "downloading Maestro sample apps (~300 MB)" >&2
    curl -fsSL --connect-timeout 20 --max-time 900 -o "$cache/samples.zip.part" \
      "https://storage.googleapis.com/mobile.dev/samples/samples.zip" >&2
    mv "$cache/samples.zip.part" "$cache/samples.zip"
  fi
  tmp="$(mktemp -d)"
  unzip -o -q "$cache/samples.zip" wikipedia.zip -d "$tmp"
  unzip -o -q "$tmp/wikipedia.zip" -d "$tmp/ios"
  mkdir -p "$cache/app"
  cp -R "$tmp"/ios/*.app "$app"
  rm -rf "$tmp"
fi

udid="$(xcrun simctl list devices available -j | DEVICE="$device" python3 -c '
import json, os, sys
q = os.environ["DEVICE"]
devs = [x for v in json.load(sys.stdin)["devices"].values() for x in v]
print(next((x["udid"] for x in devs if q in (x["udid"], x["name"])), ""))')"
if [ -z "$udid" ]; then
  echo "simulator \"$device\" not found" >&2
  exit 1
fi

echo "booting $device ($udid)" >&2
xcrun simctl boot "$udid" 2>/dev/null || true
xcrun simctl bootstatus "$udid" -b >&2
xcrun simctl install "$udid" "$app" >&2
echo "$udid"
