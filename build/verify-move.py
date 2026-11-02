"""Verify the model move against the manifest WITHOUT ever reading E:.

The E: drive is a failing NVMe: it logged 198 paging errors in one burst and the
machine bugchecked with KERNEL_DATA_INPAGE_ERROR (0x7A) shortly after. A verifier
that stats or hashes files on that drive can trip the same fault, so this script
does not open, stat or hash a single path on E:. The manifest written before any
copying is the reference; G: is the only side touched.

States per manifest row:
  SAFE          present on G at exactly the recorded size
  IN-PROGRESS   written in the last two minutes (a copy is running now)
  SIZE-DIFF     present but a different size (needs a re-copy, or it is being
                written right now)
  MISSING       not on G at all
  METADATA-OK   an app-generated sidecar: valid JSON, no E: reference inside
  METADATA-BAD  sidecar missing, unreadable, or still pointing at E:

ORPHANS are checked by size in the promoted Image models folder.

Usage: python build/verify-move.py
"""
import glob
import io
import json
import os
import sys
import time

MANIFEST = os.path.join('build', 'logs', 'e-manifest.tsv')
BS = chr(92)
DST_ROOT = os.path.join('G:' + BS, 'AI GGUF')
STAGE_IMG = os.path.join(DST_ROOT, '_incoming', 'Image models')
FINAL_IMG = os.path.join(DST_ROOT, 'Image models')

# Expected sizes, recorded from E: before the move (never re-read from E:).
ORPHANS = [
    ('controlnet_openpose_sdxl_thibaud_fp16.safetensors', 2502139104),
    ('taesdxl.safetensors', 9793292),
]


def image_dir():
    return FINAL_IMG if os.path.isdir(FINAL_IMG) else STAGE_IMG


def dest_for(rel):
    top = rel.split(os.sep)[0]
    if top == 'Image models':
        return os.path.join(image_dir(), rel.split(os.sep, 1)[1])
    return os.path.join(DST_ROOT, rel)


def metadata_state(path):
    name = os.path.basename(path)
    if not os.path.isfile(path):
        return 'METADATA-BAD', '%s is missing' % name
    try:
        text = io.open(path, encoding='utf-8', errors='replace').read()
        json.loads(text)
    except (ValueError, OSError) as e:
        return 'METADATA-BAD', '%s is not readable JSON: %s' % (name, e)
    if 'E:' in text:
        return 'METADATA-BAD', '%s still references E:' % name
    return 'METADATA-OK', ''


def main():
    if not os.path.isfile(MANIFEST):
        print('no manifest at %s' % MANIFEST)
        return 1

    tally = {}
    problems = []
    for line in io.open(MANIFEST, encoding='utf-8'):
        parts = line.rstrip('\n').split('\t')
        if len(parts) < 2:
            continue
        src, want = parts[0], int(parts[1])
        rel = os.path.relpath(src, os.path.join('E:' + BS, 'AI GGUF'))
        dst = dest_for(rel)
        base = os.path.basename(rel).lower()

        if base == 'model-index.json' or base.endswith('.howl.json'):
            state, note = metadata_state(dst)
            if note:
                problems.append(note)
        else:
            if not os.path.isfile(dst):
                state = 'MISSING'
                problems.append('%s is missing on G' % rel)
            else:
                got = os.path.getsize(dst)
                fresh = (time.time() - os.path.getmtime(dst)) < 120
                if got == want:
                    state = 'SAFE'
                elif fresh:
                    state = 'IN-PROGRESS'
                else:
                    state = 'SIZE-DIFF'
                    problems.append('%s is %d of %d bytes' % (rel, got, want))
        tally[state] = tally.get(state, 0) + 1

    ok_orphans = 0
    for name, size in ORPHANS:
        dst = os.path.join(image_dir(), name)
        if os.path.isfile(dst) and os.path.getsize(dst) == size:
            ok_orphans += 1
        else:
            problems.append('orphan %s missing or wrong size' % name)

    for p in problems[:10]:
        print('  PROBLEM %s' % p)
    print('TOTALS', ' '.join('%s=%d' % (k, v) for k, v in sorted(tally.items())),
          '| orphans_ok=%d/2' % ok_orphans)
    blocking = sum(tally.get(k, 0) for k in ('MISSING', 'SIZE-DIFF', 'IN-PROGRESS', 'METADATA-BAD'))
    settled = blocking == 0 and ok_orphans == 2 and not problems
    print('RESULT:', 'SETTLED - verified from G: alone, E: never opened' if settled else 'INCOMPLETE')
    print('NOTE: this run touched no path on E:.')
    return 0


if __name__ == '__main__':
    sys.exit(main())
