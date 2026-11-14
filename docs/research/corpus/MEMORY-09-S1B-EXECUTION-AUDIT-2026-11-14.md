# MEMORY-09 — `S1-B` Execution-Path Audit (Phase 1–2, read-only)

**Status:** audit complete for phases 1–2. **No new hash, no model load, no disk scan, no hardware intervention.**
**Date:** 2026-11-14 · **Subject:** the four distinct digests of one GGUF file, and whether the AI workflow caused them.

---

## 1. Phase 1 — what the workflow actually did

### 1.1 Question by question

| # | Question | Answer | Basis |
|---|---|---|---|
| 1 | What was instructed immediately before each hash? | (1) baseline manifest creation; (2) `qualify --slot S1-B` at 36 layers; (3–4) `qualify --slot S1-B` at 20 layers — **refused**; (5–6) post-termination re-read; (7) third post-termination read inside evidence preservation | `timeline.json`, command history |
| 2 | Was Gemma loaded in KoboldCpp during the observations? | **Yes** for observations 2, 3, 4 · **No** for 1, 5, 6, 7 | timeline contexts + server log timestamps |
| 3 | Was the file memory-mapped or held open? | **Yes while loaded** — every launch used `--usemmap`; the mapping dies with the process. **No** for observations 5–7 | launch arguments in both `kobold-S1B*.log` files |
| 4 | Did anything in the workflow have write access to the GGUF path? | **No** | see §1.2 |
| 5 | Hashes taken loaded, unloading, or idle? | 2–4 loaded · 5–6 idle (identical) · 7 idle (**different**) | timeline |

### 1.2 Write-access audit

| Check | Result |
|---|---|
| `file_digest()` | `path.open("rb")` — read-only, streamed in 1 MiB blocks |
| Harness write calls | all targets are inside `eval/` (`record.json`, `record.sha256`, `raw_body.json`, `verification.json`, `preflight.json`) — **no model path is ever a write target** |
| Copy / move / replace / rename over a model | none in `eval/harness/` |
| Scheduled tasks referencing `G:` | **none** |
| OneDrive sync root | `C:\Users\StoneWolfPC\OneDrive` — **not** `G:` |

**So the workflow did not modify the file.** A tool that only ever opens a path with `"rb"`, into a directory it never writes, cannot have changed its contents.

### 1.3 The detail that matters most

Observations **5 and 6 were identical** (`e6dca6eb…`), taken seconds apart with the **server terminated** and no harness process touching the file.

Observation **7 differed** (`eac8960a…`), also idle.

**The content changed at least once while nothing in this workflow was running.** That is the crux: a loaded/mmap'd model, a race with loading, or a write by the harness are all excluded as explanations for that change.

---

## 2. Phase 2 — what the existing evidence shows

| Item | Finding |
|---|---|
| Read/write operations against the Gemma path | reads only, by me, at the timestamps in `timeline.json` |
| File replacement / temp-file behaviour | **none observed** — same inode-level path, size constant at 5,335,285,728 B, mtime still `2026-09-04 20:53:48` |
| Memory-mapped access | present during loads 2–4 (`--usemmap`); absent for 5–7 |
| Path confusion between model and evidence output | **none** — evidence writes go to `eval/runs/…`; the digest subject is the absolute `G:\AI GGUF\GGUF\…` path recorded in `timeline.json` |
| Concurrent access during hashing | only the server's read-only mapping (observations 2–4); observations 5–7 had **no** concurrent access, and still diverged |
| **Server-side structural consistency** | the 13:50 and 14:03 loads both report `48 key-value pairs and 720 tensors`, `GGUF V3`, `n_layer = 42`, `n_embd = 2560` — **structurally identical**, and total buffer size sums equal (5,598 MiB both times, split differently by layer count) |

**Interpretation.** Structure stable, content not. A wholesale file swap would normally change size or metadata; neither happened. That pattern is more consistent with **localised content change** — a writer touching payload regions, or a storage path returning different bytes for some regions — than with someone replacing the whole model.

---

## 3. What this audit rules out

| Hypothesis | Verdict |
|---|---|
| The harness wrote to the model | **excluded** — read-only, writes confined to `eval/` |
| Concurrent access during hashing caused it | **excluded** — divergence continued with nothing attached |
| Model-output nondeterminism was being hashed | **excluded** — digests are of file bytes (§10 of the integrity record) |
| OneDrive | **excluded as syncing `G:`** — its folder is on `C:` |
| A scheduled task targeting `G:` | **not found** — no task references the volume |
| Whole-file replacement | **not supported** — size, mtime and GGUF structure all unchanged |

## 4. What remains (unranked, still unresolved)

1. **A writer outside this workflow** — a backup product (the volume is labelled *"Second backup"*), a launcher, an indexer, or a mod manager with write access to `G:`. Not visible in command lines or scheduled tasks; would show in **handle enumeration**.
2. **Storage-path instability** — a failing region, cable, or controller returning inconsistent bytes for part of the file. Also unresolved, but this audit deliberately did not test hardware.
3. **Filesystem-level behaviour** — NTFS cache/metadata faults; least likely, not excluded.

**The decisive unperformed check is open-handle enumeration.** Sysinternals `handle.exe` is not installed here, so no process was ever shown to hold the file — and equally, no process was shown *not* to. That gap is the single largest one in this investigation.

---

## 5. Recommended order (unchanged in spirit, now evidence-backed)

| Phase | Action | Requires |
|---|---|---|
| 3 | Enumerate handles on the exact file and directory; if nothing holds it, the writer hypothesis weakens sharply and storage becomes primary | a handle tool, plus permission to run it |
| 3 | Identify which product owns the "Second backup" volume and whether it writes on a schedule | read-only: backup software inventory |
| 4 | Obtain a **trusted copy on `D:`** (1.9 TB free), verify it three times **there**, and leave `G:` untouched as evidence | owner decision |
| 5 | Qualify `S1-B` and `S1-C` at S1-v1.1 against the `D:` copy; keep `S1-A` as-is; keep the v1 baseline untouched | qualification authorisation |

Behavioral S1 stays closed until all three candidates pass qualification at the profile.

---

## 6. Boundaries observed

No model loaded, no digest recomputed, no CHKDSK, no surface scan, no repair, no manifest change, no file copied. Everything above is read from filed evidence, server logs, source code, or read-only system metadata.
