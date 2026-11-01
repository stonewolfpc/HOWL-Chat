"""Print shadowed duplicate definitions with string literals intact."""
import glob
import json
import os
import re

ROOT = "frontend/dist"
data = json.load(open("build/dupdefs.json"))


def raw_bodies(path):
    text = open(path, encoding="utf-8", errors="replace").read()
    found = {}
    for m in re.finditer(r"^function\s+([A-Za-z_$][\w$]*)\s*\(", text, flags=re.M):
        name = m.group(1)
        brace = text.find("{", m.end())
        if brace < 0:
            continue
        depth, j, n = 0, brace, len(text)
        while j < n:
            c = text[j]
            if c == '"' or c == "'" or c == "`":
                q = c
                j += 1
                while j < n:
                    if text[j] == chr(92):
                        j += 2
                        continue
                    if text[j] == q:
                        break
                    j += 1
            elif c == "{":
                depth += 1
            elif c == "}":
                depth -= 1
                if depth == 0:
                    break
            j += 1
        found[name] = text[m.start():j + 1]
    return found


want = {rec["name"] for rec in data["real"]}
cache = {}
for name in sorted(want):
    print("=" * 78)
    print(name)
    for f in sorted(glob.glob(os.path.join(ROOT, "js", "*.js"))):
        base = os.path.basename(f)
        if base not in cache:
            cache[base] = raw_bodies(f)
        if name in cache[base]:
            body = cache[base][name]
            print("-" * 78)
            print("[%s]" % base)
            print(body if len(body) <= 1600 else body[:1600] + "\n    ... (%d chars)" % len(body))
    print()
