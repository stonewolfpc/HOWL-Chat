"""Duplicate scanner that also sees assignment-style definitions.

The first pass only counted `function name(` declarations, which hid a real
conflict: sidebar-context-world.js declares a bare `async function
updateWorldDisplay()` and loads after context-display.js, whose
`window.updateWorldDisplay = () => ...` it silently overwrites. This pass counts
declarations, `window.NAME =` assignments, bare `NAME =` assignments and
prototype methods.
"""
import collections
import glob
import os
import re

ROOT = "frontend/dist"


def definitions(text):
    out = collections.defaultdict(set)  # name -> set(kind)

    for m in re.finditer(r"^function\s+([A-Za-z_$][\w$]*)\s*\(", text, flags=re.M):
        out[m.group(1)].add("decl")
    for m in re.finditer(r"^async\s+function\s+([A-Za-z_$][\w$]*)\s*\(", text, flags=re.M):
        out[m.group(1)].add("decl")
    for m in re.finditer(r"^\s*window\.([A-Za-z_$][\w$]*)\s*=", text, flags=re.M):
        out[m.group(1)].add("window")
    for m in re.finditer(r"^ContextDisplay\.prototype\.([A-Za-z_$][\w$]*)\s*=", text, flags=re.M):
        out[m.group(1)].add("proto")
    for m in re.finditer(r"^([A-Za-z_$][\w$]*)\s*=\s*(?:async\s*)?(?:function|\()", text, flags=re.M):
        out[m.group(1)].add("assign")
    return out


pages = {}
for p in sorted(glob.glob(os.path.join(ROOT, "*.html"))):
    s = open(p, encoding="utf-8", errors="replace").read()
    pages[os.path.basename(p)] = [os.path.basename(x.split("?")[0])
                                  for x in re.findall(r'<script[^>]+src="js/([^"]+)"', s)]

by_name = collections.defaultdict(dict)
for f in sorted(glob.glob(os.path.join(ROOT, "js", "*.js"))):
    base = os.path.basename(f)
    for name, kinds in definitions(open(f, encoding="utf-8", errors="replace").read()).items():
        by_name[name][base] = kinds

dups = {k: v for k, v in by_name.items() if len(v) > 1}

print("=== duplicate names (declarations + assignments + prototypes): %d ===" % len(dups))

shadows = []
for name, defs in sorted(dups.items()):
    for pg, srcs in pages.items():
        hit = [f for f in srcs if f in defs]
        if len(hit) >= 2:
            kinds = ["%s[%s]" % (f, ",".join(sorted(defs[f]))) for f in hit]
            shadows.append((name, pg, kinds))
            break

print("\n=== names defined by 2+ files that load together on one page: %d ===" % len(shadows))
for name, pg, kinds in shadows:
    print("  %-24s %-20s %s" % (name, pg, "  ->  ".join(kinds)))

print("\n=== assignment-vs-declaration conflicts (same name, different mechanism) ===")
for name, defs in sorted(dups.items()):
    kinds = set()
    for v in defs.values():
        kinds |= v
    if len(kinds) > 1:
        print("  %-24s %s" % (name, ", ".join("%s(%s)" % (f, ",".join(sorted(k))) for f, k in defs.items())))
