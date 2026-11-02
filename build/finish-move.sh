#!/usr/bin/env bash
# finish-move.sh — the tail end of the E: -> G: move, run unattended.
#
# The main script's orphan phase used robocopy file-to-file, which robocopy
# refuses (exit 16, nothing copied), so those two files never landed. This waits
# for the GGUF copy to finish — never two readers on the dying drive at once —
# then copies them with the directory form, promotes the staged Image models
# folder, and re-runs the verifier.
set -u

LOG="build/logs/move-to-g.log"
log() { printf '%s  %s\n' "$(date +%H:%M:%S)" "$*" | tee -a "$LOG"; }

log "finish: waiting for the main copy to go idle"
for _ in $(seq 1 360); do
	if [ "$(tasklist 2>/dev/null | grep -ci robocopy)" -eq 0 ]; then
		log "finish: drive is idle, starting"
		break
	fi
	sleep 15
done

if [ "$(tasklist 2>/dev/null | grep -ci robocopy)" -ne 0 ]; then
	log "finish: main copy still running after 90 minutes, leaving it alone"
	exit 0
fi

STAGE="/g/AI GGUF/_incoming/Image models"
FINAL="/g/AI GGUF/Image models"
mkdir -p "$STAGE"

# Directory form with an explicit file list: robocopy rejects a file-to-file
# destination, which is what the first attempt did.
for f in controlnet_openpose_sdxl_thibaud_fp16.safetensors taesdxl.safetensors; do
	if [ -f "/e/HOWL_Chat/models/image/$f" ]; then
		MSYS2_ARG_CONV_EXCL="*" robocopy.exe 'E:\HOWL_Chat\models\image' "$(cygpath -w "$STAGE")" "$f" \
			/J /COPY:DAT /R:2 /W:2 /NP /NFL /NDL >>"$LOG" 2>&1
		log "finish: orphan $f rc=$? -> $(stat -c%s "$STAGE/$f" 2>/dev/null || echo missing) bytes"
	fi
done

# model.howl.json is deliberately not copied: the Image models folder has its
# own 368-byte sidecar and the 536-byte one from HOWL's models dir would clobber
# it with metadata about a directory that is being retired.

if [ ! -d "$FINAL" ] && [ -d "$STAGE" ]; then
	mv "$STAGE" "$FINAL" && log "finish: promoted _incoming/Image models -> Image models"
else
	log "finish: promote not needed (final exists: $([ -d "$FINAL" ] && echo yes || echo no))"
fi

python build/verify-move.py --hash 2>&1 | tee -a "$LOG"
log "finish: done"
