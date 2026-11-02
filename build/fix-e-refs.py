"""Repair stale E: references in the model registry and drop exact duplicates.

The app's rescan appended entries for models on G without retiring the ones that
still pointed at E:, so the registry kept 12 paths on a drive that must never be
read again.

Two passes:
  1. text-level path rewrite E: -> G:. The file holds JSON escapes, so every
     separator on disk is a doubled backslash; matching a single one finds
     nothing (which is exactly how the first attempt silently did nothing).
  2. drop entries identical in type, family and model path, keeping the first.
     The path is required: keying on fields that simply do not exist for LLM
     entries collapsed 17 distinct models into one in the previous attempt.
"""
import io
import json
import os
import shutil

REG = os.path.join('config', 'model-registry.json')
BACKUP = os.path.join('build', 'logs', 'model-registry.before-e-fix.json')
BS = chr(92)
SEP = BS + BS  # what a path separator looks like on disk


def model_path(entry):
    files = entry.get('files') or {}
    return str(entry.get('path') or files.get('model') or '')


def main():
    raw = io.open(REG, encoding='utf-8').read()
    if not os.path.isfile(BACKUP):
        shutil.copyfile(REG, BACKUP)
        print('backup written:', BACKUP)

    rewrites = [
        ('E:' + SEP + 'AI GGUF' + SEP, 'G:' + SEP + 'AI GGUF' + SEP),
        ('E:' + SEP + 'HOWL_Chat' + SEP + 'models' + SEP + 'image' + SEP,
         'G:' + SEP + 'AI GGUF' + SEP + 'Image models' + SEP),
    ]
    total = 0
    for old, new in rewrites:
        n = raw.count(old)
        if n:
            raw = raw.replace(old, new)
            total += n
            print('rewrote %d path(s): %s -> %s' % (n, old, new))
    print('total path rewrites:', total)
    io.open(REG, 'w', encoding='utf-8', newline='').write(raw)

    rows = json.loads(raw)
    seen = set()
    kept, dropped = [], []
    for r in rows:
        path = model_path(r)
        if not path:
            kept.append(r)          # nothing to key on: never merge
            continue
        key = (r.get('type'), r.get('family'), path.lower())
        if key in seen:
            dropped.append(os.path.basename(path))
            continue
        seen.add(key)
        kept.append(r)
    print('entries: %d -> %d (dropped %d duplicates on the same path)' % (
        len(rows), len(kept), len(dropped)))
    for name in dropped[:12]:
        print('  dropped duplicate:', name)
    io.open(REG, 'w', encoding='utf-8', newline='').write(
        json.dumps(kept, indent=2, ensure_ascii=False) + '\n')

    after = io.open(REG, encoding='utf-8').read()
    leftovers = after.count('E:' + BS)
    rows = json.loads(after)
    missing, checked = [], 0
    for r in rows:
        candidates = [str(r.get('path') or '')] + [str(v) for v in (r.get('files') or {}).values()]
        for path in candidates:
            if len(path) > 2 and path[1] == ':':
                checked += 1
                if not os.path.isfile(path):
                    missing.append('%s' % os.path.basename(path))
    print('VERIFY: E: references remaining = %d' % leftovers)
    print('VERIFY: paths checked = %d, missing on disk = %d' % (checked, len(missing)))
    for m in missing[:8]:
        print('  missing:', m)
    print('RESULT:', 'CLEAN' if leftovers == 0 and not missing else 'NEEDS-ATTENTION')


if __name__ == '__main__':
    main()
