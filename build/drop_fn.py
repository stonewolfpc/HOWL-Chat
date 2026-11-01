"""Delete a top-level `function name(...) { ... }` from a JS file, safely.

Safety net learned the hard way: a naive brace counter desyncs on a regex
literal containing a quote (`/'/g`), which silently swallowed six neighbouring
functions. So every drop now:

  1. snapshots the file to build/backup/
  2. removes the definition
  3. checks the result with `node --check`
  4. compares the top-level function inventory: exactly one name may disappear,
     and it must be the requested one
  5. restores the snapshot and reports FAIL if anything else changed
"""
import os
import re
import shutil
import subprocess
import sys

BACKUP = os.path.join("build", "backup")


def function_names(text):
    return re.findall(r"^function\s+([A-Za-z_$][\w$]*)\s*\(", text, flags=re.M)


def find_definition(text, name):
    """Return (start, end) for the single top-level definition of `name`."""
    matches = list(re.finditer(r"^function\s+%s\s*\(" % re.escape(name), text, flags=re.M))
    if len(matches) != 1:
        return None
    m = matches[0]
    brace = text.find("{", m.end())
    if brace < 0:
        return None
    depth, j, n = 0, brace, len(text)
    while j < n:
        c = text[j]
        if c in "\"'`":
            quote = c
            j += 1
            while j < n:
                if text[j] == chr(92):
                    j += 2
                    continue
                if text[j] == quote:
                    break
                j += 1
        elif c == "/" and j + 1 < n and text[j + 1] == "/":
            nl = text.find("\n", j)
            j = n if nl < 0 else nl
            continue
        elif c == "/" and j + 1 < n and text[j + 1] == "*":
            close = text.find("*/", j + 2)
            j = n if close < 0 else close + 1
        elif c == "{":
            depth += 1
        elif c == "}":
            depth -= 1
            if depth == 0:
                break
        j += 1
    return m.start(), j + 1


def main():
    path, name = sys.argv[1], sys.argv[2]
    text = open(path, encoding="utf-8", errors="replace", newline="").read()
    before = function_names(text)

    span = find_definition(text, name)
    if not span:
        print("SKIP   %-30s %s not found exactly once" % (name, path))
        return 1

    os.makedirs(BACKUP, exist_ok=True)
    snap = os.path.join(BACKUP, os.path.basename(path) + ".bak")
    if not os.path.exists(snap):
        shutil.copy2(path, snap)

    start, end = span
    while end < len(text) and text[end] in "\r\n":
        end += 1
    result = text[:start] + text[end:]
    open(path, "w", encoding="utf-8", newline="").write(result)

    after = function_names(result)
    lost = [f for f in before if f not in after]
    added = [f for f in after if f not in before]
    syntax = subprocess.run(["node", "--check", path], capture_output=True, text=True)

    if syntax.returncode != 0 or lost != [name] or added:
        shutil.copy2(snap, path)
        print("FAILED %-30s restored. lost=%s added=%s syntax=%s"
              % (name, lost, added, syntax.returncode))
        return 2

    print("removed %-30s from %-28s %5d bytes" % (name, os.path.basename(path), end - start))
    return 0


if __name__ == "__main__":
    sys.exit(main())
