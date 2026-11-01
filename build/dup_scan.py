"""Duplicate-definition scanner for frontend/dist.

Reports, for every function name defined at top level in two or more JS files:
  - the defining files and whether the bodies are identical or divergent
  - which pages load 2+ of those files (i.e. a real shadow, last one wins)
  - which pages load only one (harmless duplication, organisational only)
"""
import collections
import glob
import hashlib
import json
import os
import re

ROOT = "frontend/dist"


def strip_code(text):
    """Remove comments and string literals so brace matching is reliable."""
    out, i, n = [], 0, len(text)
    while i < n:
        c = text[i]
        if c == "/" and i + 1 < n and text[i + 1] == "*":
            j = text.find("*/", i + 2)
            i = n if j < 0 else j + 2
            continue
        if c == "/" and i + 1 < n and text[i + 1] == "/":
            j = text.find("\n", i)
            i = n if j < 0 else j
            continue
        if c in "\"'`":
            quote = c
            i += 1
            while i < n:
                if text[i] == chr(92):
                    i += 2
                    continue
                if text[i] == quote:
                    i += 1
                    break
                i += 1
            out.append('""')
            continue
        out.append(c)
        i += 1
    return "".join(out)


def top_level_defs(path):
    """name -> normalised body text for every `function name(` at column 0."""
    code = strip_code(open(path, encoding="utf-8", errors="replace").read())
    found = {}
    for m in re.finditer(r"^function\s+([A-Za-z_$][\w$]*)\s*\(", code, flags=re.M):
        name = m.group(1)
        brace = code.find("{", m.end())
        if brace < 0:
            continue
        depth, j = 0, brace
        while j < len(code):
            if code[j] == "{":
                depth += 1
            elif code[j] == "}":
                depth -= 1
                if depth == 0:
                    break
            j += 1
        body = code[brace:j + 1]
        found[name] = re.sub(r"\s+", " ", body).strip()
    return found


def main():
    pages = {}
    for p in sorted(glob.glob(os.path.join(ROOT, "*.html"))):
        s = open(p, encoding="utf-8", errors="replace").read()
        srcs = [os.path.basename(x.split("?")[0])
                for x in re.findall(r'<script[^>]+src="js/([^"]+)"', s)]
        pages[os.path.basename(p)] = srcs

    print("=== page -> script count ===")
    for pg, srcs in pages.items():
        print("  %-24s %2d" % (pg, len(srcs)))

    by_name = collections.defaultdict(dict)  # name -> {file: body}
    for f in sorted(glob.glob(os.path.join(ROOT, "js", "*.js"))):
        base = os.path.basename(f)
        for name, body in top_level_defs(f).items():
            by_name[name][base] = body

    dups = {k: v for k, v in by_name.items() if len(v) > 1}

    real, harmless = [], []
    for name, defs in sorted(dups.items()):
        hashes = {hashlib.md5(b.encode()).hexdigest()[:8] for b in defs.values()}
        shadows = []
        for pg, srcs in pages.items():
            hit = [f for f in srcs if f in defs]
            if len(hit) >= 2:
                shadows.append((pg, hit))
        rec = {"name": name, "files": sorted(defs), "variants": len(hashes),
               "shadows": shadows, "bodies": defs}
        (real if shadows else harmless).append(rec)

    print("\n=== duplicate names: %d  (real shadows: %d, cross-page only: %d) ==="
          % (len(dups), len(real), len(harmless)))

    print("\n--- REAL SHADOWS: two definers load on one page, later file wins ---")
    for rec in real:
        kind = "identical" if rec["variants"] == 1 else "DIVERGENT x%d" % rec["variants"]
        print("\n%s   [%s]" % (rec["name"], kind))
        for pg, order in rec["shadows"]:
            print("   %-22s %s" % (pg, "  ->  ".join(order)))

    print("\n--- duplicate but never co-loaded (organisation only) ---")
    for rec in harmless:
        kind = "identical" if rec["variants"] == 1 else "divergent x%d" % rec["variants"]
        print("  %-28s %-14s %s" % (rec["name"], kind, ", ".join(rec["files"])))

    json.dump({"real": [{k: v for k, v in r.items()} for r in real],
               "harmless": harmless},
              open("build/dupdefs.json", "w"), indent=1)
    print("\nwrote build/dupdefs.json")


if __name__ == "__main__":
    main()
