#!/usr/bin/env bash
# Usage: generate.sh <out-dir>; prints produced scenarios to stdout, tool output goes to stderr.
set -uo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
out="$1"

# The failing test is intentional, so Maven must not fail the build.
(cd "$here" && mvn -B -q -Dmaven.test.failure.ignore=true clean test >&2) || exit 1

mkdir -p "$out/junit/maven/input"
cp "$here"/target/surefire-reports/TEST-*.xml "$out/junit/maven/input/" || exit 1

# Surefire dumps every JVM system property (user, home, tmp paths), so replace them with fixed neutral ones.
for f in "$out"/junit/maven/input/TEST-*.xml; do
  awk '
    /<properties>/ {
      print "  <properties>"
      print "    <property name=\"java.version\" value=\"17\"/>"
      print "    <property name=\"file.encoding\" value=\"UTF-8\"/>"
      print "    <property name=\"maven.version\" value=\"3.9.9\"/>"
      print "  </properties>"
      skip = 1
    }
    !skip { print }
    /<\/properties>/ { skip = 0 }
  ' "$f" > "$f.tmp" && mv "$f.tmp" "$f"
done

echo "junit/maven"
