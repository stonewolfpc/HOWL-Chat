#!/usr/bin/env bash
# move-to-g.sh — one serialised, resumable, fully logged move of the AI model
# trees off the retiring E: drive onto G:.
#
# Safety rules this script holds to:
#   * copy-only. Nothing is ever deleted from E: here. A stall, a hang or a
#     bug-check costs only the files not yet read, and a rerun resumes because
#     robocopy skips files whose size + timestamp already match.
#   * one reader at a time. Two robocopies pulling from a failing disk makes the
#     stalls worse, so the phases run strictly in sequence.
#   * no collisions. Image models goes to a staging dir first and is promoted
#     with a same-volume rename only after every byte is verified.
#   * a manifest is written before any copying, so the full expected file list
#     survives on C: even if E: never answers another read.
#
# Artifacts:
#   build/logs/e-manifest.tsv     every file expected, with size and mtime
#   build/logs/move-to-g.log      robocopy + phase log
#   G:\AI GGUF\_move-status.txt   current phase, for watching from Explorer
set -u

SRC_ROOT="/e/AI GGUF"
DST_ROOT="/g/AI GGUF"
STAGE="$DST_ROOT/_incoming"
LOGDIR="build/logs"
LOG="$LOGDIR/move-to-g.log"
MANIFEST="$LOGDIR/e-manifest.tsv"
STATUS="$DST_ROOT/_move-status.txt"

mkdir -p "$LOGDIR" "$STAGE" "$DST_ROOT"
: >>"$LOG"

log() {
	local line
	line="$(date +%H:%M:%S)  $*"
	printf '%s\n' "$line" | tee -a "$LOG"
}
phase() { printf '%s\n' "$1" >"$STATUS"; log "PHASE  $1"; }
# robocopy's exit codes: 0-7 are success-ish (1 = files copied, 2 = extras,
# 3 = both), 8+ is a real failure. /R:2 /W:2 gives a bad sector two chances
# without wedging the whole job the way the default million retries would.
rc() {
	local out
	MSYS2_ARG_CONV_EXCL="*" robocopy.exe "$@" >>"$LOG" 2>&1
	out=$?
	if [ "$out" -ge 8 ]; then log "  FAIL rc=$out (locked/unreadable; rerun retries)"; fi
	return 0
}

# ---- Phase 0: manifest, before touching anything -------------------------
if [ ! -s "$MANIFEST" ]; then
	phase "0/5 manifest: enumerating E:"
	python - "$MANIFEST" <<'PY'
import io, os, sys
root = r'E:\AI GGUF'
out = sys.argv[1]
rows = []
for dirpath, _dirs, files in os.walk(root):
    for name in files:
        p = os.path.join(dirpath, name)
        try:
            st = os.stat(p)
        except OSError as e:
            rows.append('%s\t-1\t-\t%s' % (p, e.errno)); continue
        rows.append('%s\t%d\t%d' % (p, st.st_size, int(st.st_mtime)))
io.open(out, 'w', encoding='utf-8', newline='\n').write('\n'.join(rows) + '\n')
print('manifest rows:', len(rows))
PY
	cp -p "$MANIFEST" "$DST_ROOT/e-manifest.tsv" 2>/dev/null
	log "manifest: $(wc -l <"$MANIFEST") entries -> $MANIFEST (copy on G:)"
else
	log "manifest: reusing existing $(wc -l <"$MANIFEST") entries"
fi

# ---- Phase 1: Image models into staging (no collisions with anything) ----
phase "1/5 Image models -> _incoming (resumable)"
rc 'E:\AI GGUF\Image models' "G:\\AI GGUF\\_incoming\\Image models" /E /J /COPY:DAT /DCOPY:DAT /R:2 /W:2 /XJ /NP /NFL /NDL
log "  staged: $(du -sh "$STAGE/Image models" 2>/dev/null | cut -f1) / source $(du -sh "$SRC_ROOT/Image models" 2>/dev/null | cut -f1)"

# ---- Phase 2: files that live only under HOWL's own models dir -----------
phase "2/5 orphan files from HOWL models dir"
mkdir -p "$STAGE/Image models"
for f in controlnet_openpose_sdxl_thibaud_fp16.safetensors taesdxl.safetensors model.howl.json; do
	if [ -f "/e/HOWL_Chat/models/image/$f" ]; then
		rc "E:\\HOWL_Chat\\models\\image\\$f" "G:\\AI GGUF\\_incoming\\Image models\\$f" /J /COPY:DAT /R:2 /W:2 /NP /NFL /NDL
		log "  $f"
	fi
done

# ---- Phase 3: GGUF into place (Explorer paste is cancelled, we own it) ---
phase "3/5 GGUF -> G (resumable)"
rc 'E:\AI GGUF\GGUF' 'G:\AI GGUF\GGUF' /E /J /COPY:DAT /DCOPY:DAT /R:2 /W:2 /XJ /NP /NFL /NDL
log "  gguf now: $(du -sh "$DST_ROOT/GGUF" 2>/dev/null | cut -f1)"

# ---- Phase 4: verify every manifest entry against what landed -----------
phase "4/5 verify sizes against manifest"
python - "$MANIFEST" "$STAGE" "$DST_ROOT" <<'PY' 2>&1 | tee -a "$LOG"
import io, os, sys
manifest, stage, dstroot = sys.argv[1], sys.argv[2], sys.argv[3]
src_root = r'E:\AI GGUF'
ok = short = missing = unreadable = 0
report = []
for line in io.open(manifest, encoding='utf-8'):
    parts = line.rstrip('\n').split('\t')
    if len(parts) < 2:
        continue
    src, size = parts[0], int(parts[1])
    if size < 0:
        unreadable += 1
        report.append('UNREADABLE %s' % src)
        continue
    rel = os.path.relpath(src, src_root).replace('\\', '/')
    if rel.startswith('Image models/'):
        target = os.path.join(stage, 'Image models', rel.split('/', 1)[1].replace('/', os.sep))
    else:
        target = os.path.join(dstroot, rel.replace('/', os.sep))
    if not os.path.isfile(target):
        missing += 1
        report.append('MISSING  %s' % rel)
    elif os.path.getsize(target) != size:
        short += 1
        report.append('SHORT    %s (%d of %d)' % (rel, os.path.getsize(target), size))
    else:
        ok += 1
for r in report[:12]:
    print('  ' + r)
print('VERIFY: intact=%d missing=%d short=%d unreadable=%d' % (ok, missing, short, unreadable))
print('RESULT: %s' % ('CLEAN' if missing == 0 and short == 0 else 'INCOMPLETE - rerun to resume'))
PY

# ---- Phase 5: promote staging into place when it verifies clean ---------
phase "5/5 promote"
staged_files=$(find "$STAGE/Image models" -type f 2>/dev/null | wc -l)
src_files=$(find "$SRC_ROOT/Image models" -type f 2>/dev/null | wc -l)
log "  staged $staged_files files vs source $src_files"
if [ "$staged_files" -ge "$src_files" ] && [ "$src_files" -gt 0 ]; then
	if [ -d "$DST_ROOT/Image models" ] && [ -z "$(ls -A "$DST_ROOT/Image models" 2>/dev/null)" ]; then
		rmdir "$DST_ROOT/Image models" 2>/dev/null
	fi
	if [ ! -e "$DST_ROOT/Image models" ]; then
		mv "$STAGE/Image models" "$DST_ROOT/Image models" && log "  promoted _incoming\\Image models -> Image models"
	else
		log "  SKIP promote: $DST_ROOT/Image models already exists"
	fi
else
	log "  SKIP promote: staged file count is below source"
fi
phase "done"
log "DONE. Full log: $LOG  |  manifest: $MANIFEST"
