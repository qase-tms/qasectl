#!/usr/bin/env bash
# Usage: generate.sh [GEN...] — regenerates fixture inputs; no args = all generators.
set -uo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
testdata="$here/../testdata"
ALL="junit-pytest junit-maven junit-jest allure-pytest maestro"

gens="$*"
[ -n "$gens" ] || gens="$ALL"

tools_for() {
  case "$1" in
    junit-pytest|allure-pytest) echo "python3" ;;
    junit-maven) echo "mvn java" ;;
    junit-jest) echo "node npm" ;;
    maestro) echo "maestro xcrun" ;;
    *) echo "" ;;
  esac
}

have_tool() {
  command -v "$1" >/dev/null 2>&1 && return 0
  [ "$1" = "maestro" ] && [ -x "$HOME/.maestro/bin/maestro" ] && return 0
  return 1
}

skipped=0
failed=0

for gen in $gens; do
  script="$here/$gen/generate.sh"
  if [ ! -f "$script" ]; then
    echo "skip $gen: generator not found"
    skipped=1
    continue
  fi

  missing=""
  for tool in $(tools_for "$gen"); do
    have_tool "$tool" || { missing="$tool"; break; }
  done
  if [ -n "$missing" ]; then
    echo "skip $gen: $missing not found"
    skipped=1
    continue
  fi

  out="$(mktemp -d)"
  if [ "$gen" = "maestro" ] && ! command -v maestro >/dev/null 2>&1; then
    export PATH="$HOME/.maestro/bin:$PATH"
  fi

  echo "generate $gen" >&2
  produced="$("$script" "$out")" || { echo "fail $gen: generator exited with an error" >&2; failed=1; rm -rf "$out"; continue; }
  "$here/shrink.sh" "$out" || { echo "fail $gen: shrink failed" >&2; failed=1; rm -rf "$out"; continue; }

  for scenario in $produced; do
    src="$out/$scenario/input"
    dst="$testdata/$scenario/input"
    if [ ! -d "$src" ]; then
      echo "fail $gen: $scenario produced no input/ directory" >&2
      failed=1
      continue
    fi
    rm -rf "$dst"
    mkdir -p "$dst"
    cp -R "$src/." "$dst/"
    echo "updated $scenario"
  done
  rm -rf "$out"
done

[ "$skipped" -eq 0 ] && [ "$failed" -eq 0 ]
