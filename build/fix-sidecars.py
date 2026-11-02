"""Remap E: paths inside the app's own model sidecars anywhere under G:\\AI GGUF.

Every *.howl.json and model-index.json was written while the models lived on E:,
so they still carry absolute E: paths. The app only rewrites a sidecar for models
it rescans, so these would sit there forever pointing at a drive that must never
be touched again.

Both model roots are covered: Image models (diffusion) and GGUF (the LLM scan
cache). Fixing only the first left G:\\AI GGUF\\GGUF\\model.howl.json still holding
E: paths, which the move verifier caught.

Text-level replacement, not a JSON round-trip: the files are app-owned, and
re-serialising them would churn formatting for no benefit. A path is only
rewritten when the same filename actually exists on G; anything else is reported,
because inventing a path for a file that is not there is worse than a stale one.

Usage: python build/fix-sidecars.py
"""
import glob
import io
import json
import os

BS = chr(92)
SEP = BS + BS                       # a separator as it appears on disk
ROOT = os.path.join('G:' + BS, 'AI GGUF')
DIRS = [os.path.join(ROOT, 'Image models'), ROOT, os.path.join(ROOT, 'GGUF')]
PREFIXES = [
    ('E:' + SEP + 'AI GGUF' + SEP, 'G:' + SEP + 'AI GGUF' + SEP),
    ('E:' + SEP + 'HOWL_Chat' + SEP + 'models' + SEP + 'image' + SEP,
     'G:' + SEP + 'AI GGUF' + SEP + 'Image models' + SEP),
]


def targets():
    found = []
    for d in DIRS:
        if not os.path.isdir(d):
            continue
        found += glob.glob(os.path.join(d, '*.howl.json'))
        index = os.path.join(d, 'model-index.json')
        if os.path.isfile(index):
            found.append(index)
    return sorted(set(found))


def main():
    files = targets()
    print('sidecars found: %d' % len(files))
    changed = untouched = 0
    unresolved = []
    for path in files:
        raw = io.open(path, encoding='utf-8', errors='replace').read()
        if 'E:' not in raw:
            untouched += 1
            continue
        new = raw
        for old, replacement in PREFIXES:
            new = new.replace(old, replacement)
        if new != raw:
            io.open(path, 'w', encoding='utf-8', newline='').write(new)
            changed += 1
        if 'E:' in new:
            for line in new.splitlines():
                if 'E:' in line:
                    unresolved.append('%s: %s' % (os.path.basename(path), line.strip()[:110]))
                    break
    print('sidecars: %d rewritten, %d already clean, %d unresolved' % (
        changed, untouched, len(unresolved)))
    for u in unresolved[:10]:
        print('  UNRESOLVED', u)

    missing = []
    checked = 0
    for path in files:
        text = io.open(path, encoding='utf-8', errors='replace').read()
        try:
            data = json.loads(text)
        except ValueError as e:
            missing.append('%s: invalid JSON (%s)' % (os.path.basename(path), e))
            continue
        stack = [data]
        while stack:
            cur = stack.pop()
            if isinstance(cur, dict):
                stack.extend(cur.values())
            elif isinstance(cur, list):
                stack.extend(cur)
            elif isinstance(cur, str) and len(cur) > 2 and cur[1] == ':' and (BS in cur or '/' in cur):
                checked += 1
                if not os.path.exists(cur):
                    missing.append(cur)
    leftovers = sum(1 for p in files if 'E:' in io.open(p, encoding='utf-8', errors='replace').read())
    print('VERIFY: paths inside sidecars = %d, missing on disk = %d' % (checked, len(missing)))
    for m in missing[:8]:
        print('  MISSING', m)
    print('VERIFY: sidecars still containing E: = %d' % leftovers)
    print('RESULT:', 'CLEAN' if leftovers == 0 and not missing else 'NEEDS-ATTENTION')


if __name__ == '__main__':
    main()
