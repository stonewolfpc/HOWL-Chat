#!/usr/bin/env bash
# Engine train — capture and diff a stable-diffusion.cpp release before any
# binary in the tree is replaced.
#
# Usage:
#   scripts/engine-train.sh <tag> <asset-zip-url> [<cudart-zip-url>]
#
# What it does, in order:
#   1. downloads the release asset into build/tmp/engine-train/<tag>/
#   2. unpacks it and prints the commit the binary reports
#   3. captures `sd-server --help` as build/engine-baseline/sd.cpp.<tag>.help.txt
#   4. diffs the flag surface against the newest committed baseline
#
# It never touches the binaries in the project root. Replacing them is a
# separate, deliberate step after this diff has been read.
set -euo pipefail

TAG="${1:?usage: engine-train.sh <tag> <asset-zip-url> [<cudart-zip-url>]}"
ASSET="${2:?missing asset zip url}"
CUDART="${3:-}"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORK="$ROOT/build/tmp/engine-train/$TAG"
BASE="$ROOT/build/engine-baseline"
mkdir -p "$WORK" "$BASE"

echo "== downloading $TAG"
curl -fsSL -o "$WORK/sd.zip" "$ASSET"
unzip -oq "$WORK/sd.zip" -d "$WORK/dist"
if [ -n "$CUDART" ]; then
  curl -fsSL -o "$WORK/cudart.zip" "$CUDART"
  unzip -oq "$WORK/cudart.zip" -d "$WORK/dist"
fi

SERVER="$(find "$WORK/dist" -name sd-server.exe | head -1)"
[ -n "$SERVER" ] || { echo "no sd-server.exe in $WORK/dist"; exit 1; }

echo "== binary reports:"
"$SERVER" --help 2>&1 | head -1

"$SERVER" --help > "$BASE/sd.cpp.$TAG.help.txt" 2>&1
echo "== captured $BASE/sd.cpp.$TAG.help.txt"

PREV="$(ls -1t "$BASE"/sd.cpp.*.help.txt | grep -v "$TAG" | head -1 || true)"
if [ -n "$PREV" ]; then
  echo "== flag diff against $(basename "$PREV")"
  python - "$PREV" "$BASE/sd.cpp.$TAG.help.txt" <<'PY'
import re, sys
def flags(p):
    return set(re.findall(r'(?<![\w-])(--[a-z0-9][a-z0-9_.-]+)', open(p, encoding='utf-8', errors='replace').read()))
old, new = flags(sys.argv[1]), flags(sys.argv[2])
print("  gained:", " ".join(sorted(new - old)) or "(none)")
print("  lost  :", " ".join(sorted(old - new)) or "(none)")
PY
else
  echo "== no previous baseline to diff against"
fi

echo
echo "Next: review 'lost'. Every flag HOWL emits that appears there fails the"
echo "upgrade before it starts (internal/image/engine/sdcpp.UnsupportedFlags)."
echo "Only then copy the binaries from $WORK/dist over the project root."
