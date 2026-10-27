#!/usr/bin/env bash
# Usage: shrink.sh <dir> — shrinks generated reports in place to keep fixtures small and machine-independent.
set -euo pipefail

dir="${1:?usage: shrink.sh <dir>}"
PNG_B64='iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVR4nGP4//8/AAX+Av4N70a4AAAAAElFTkSuQmCC'

sed_inplace() {
  if sed --version >/dev/null 2>&1; then sed -i "$@"; else sed -i '' "$@"; fi
}

b64_decode() {
  if base64 -D </dev/null >/dev/null 2>&1; then base64 -D; else base64 -d; fi
}

tmp_png="$(mktemp)"
trap 'rm -f "$tmp_png"' EXIT
printf '%s' "$PNG_B64" | b64_decode >"$tmp_png"

find "$dir" -type f \( -iname '*.png' -o -iname '*.jpg' -o -iname '*.jpeg' \) | while IFS= read -r f; do
  cp "$tmp_png" "$f"
done

find "$dir" -type f -name '*.log' | while IFS= read -r f; do
  head -n 5 "$f" >"$f.tmp" && mv "$f.tmp" "$f"
done

find "$dir" -type d -name 'screen-hierarchy' -prune -exec rm -rf {} +

repo_root="$(git rev-parse --show-toplevel)"
find "$dir" -type f \( -name '*.xml' -o -name '*.json' -o -name '*.txt' \) | while IFS= read -r f; do
  sed_inplace -e "s|$repo_root|/workspace|g" -e "s|$HOME|/workspace|g" "$f"
done
