#!/usr/bin/env bash
# copy-modpack.sh — last file to leave E:. Mirrors the archive to the same
# relative path on G so build/verify-move.py verifies it with no special case.
set -u

REL='[ScrewThisNoise][OPTIONAL][KK-KKS][20250319] Sideloader Modpack Studio'
LOG="build/logs/modpack.log"

log() { printf '%s  %s\n' "$(date +%H:%M:%S)" "$*" >>"$LOG"; }

log "modpack: robocopy start ($REL)"

# Directory form, no file list: robocopy takes the whole subtree. /J unbuffered
# so a dropped read fails this file instead of stalling the transfer, /XJ keeps
# junctions out, and /R:2 gives up on a bad sector quickly rather than retrying
# a million times on a drive that is already dying.
MSYS2_ARG_CONV_EXCL="*" robocopy.exe \
	"E:\\AI GGUF\\$REL" \
	"G:\\AI GGUF\\$REL" \
	/E /J /COPY:DAT /DCOPY:DAT /R:2 /W:2 /XJ /NP /NFL /NDL >>"$LOG" 2>&1
rc=$?
log "modpack: robocopy rc=$rc (0-7 success, 8+ failure)"

src_bytes=$(stat -c%s "/e/AI GGUF/$REL/$REL.7z" 2>/dev/null || echo 0)
dst_bytes=$(stat -c%s "/g/AI GGUF/$REL/$REL.7z" 2>/dev/null || echo 0)
log "modpack: size E=$src_bytes G=$dst_bytes"

python build/verify-move.py --hash 2>&1 | tee -a "$LOG"
log "modpack: done"
