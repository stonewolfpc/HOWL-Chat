# S1-B Read-Stability Probe — Specification

**Identity:** `MEMORY-09-S1B-STAB-1` · **Status:** **SPECIFICATION ONLY — not authorised, not executed.**
**Question:** does the current `S1-B` file produce a **stable digest across independently initiated reads**, and does that depend on the read path?
**Gate:** this is a **new operation on `G:`** (repeated reads of a model file). It is not gate 1 (scratch file), not gate 2 (event log), not gate 3 (USN). It needs its own authorisation.

---

## 1. Why this exists

The readiness report recommended re-baselining `S1-B` to `d0027dd3…6e96`, because one cache-bypassing read matched the repository-hosted LFS pointer — a repository-verified reference, not publisher-verified provenance. That recommendation is **withdrawn as premature** by the owner's objection, and the objection is correct: seven digests were recorded for one file identity, and **one matching read does not establish stability**. Re-baselining on a single observation would calibrate the experiment to an input whose behaviour is still unknown.

**This probe measures stability. It decides nothing.** Changing the expected hash and qualifying the slot remain separate, later decisions.

---

## 2. What counts as a read here

| Label | Tool | Path | What it may be called |
|---|---|---|---|
| **B** | `sha256sum` (GNU coreutils) | normal cached I/O | **buffered read** |
| **C** | `eval/harness/unbuffered_digest.py`, sector 4096, block 4 MiB | `FILE_FLAG_NO_BUFFERING` body + buffered tail | **unbuffered-body composite read** |
| **B*** | `sha256sum` after an intervening `C` | normal cached I/O | buffered read, cache state made explicit |

**No read in this probe may be described as "device-level."** `NO_BUFFERING` bypasses the system file cache. It does not bypass the filter stack — `WdFilter`, `FileInfo`, `Wof`, `bfs`, `gameflt` and `UCPD` are attached to `G:` — nor the class/port drivers. The strongest available description is *cache-bypassing*, and for tool C the digest is a **composite** of an unbuffered body and a buffered 4,064-byte tail.

**Independence, defined honestly.** Each read runs in a **separate process invocation**, started after the previous one exits. That is independent *initiation*. It is **not** independent cache state: the system cache persists across processes, and no cache purge is performed or possible. Any claim of "independent reads" here means independent process starts, nothing more.

---

## 3. Procedure

Ten reads, fixed order, one process each, no overlap.

| # | Kind | Note |
|---|---|---|
| 1 | B | first touch in this session |
| 2 | C | — |
| 3 | B | after a cache-bypassing read |
| 4 | C | repeat |
| 5 | B | — |
| 6 | C | — |
| 7 | B | after a **60 s idle** — time-dependence check |
| 8 | C | immediately after 7 |
| 9 | B | — |
| 10 | C | final |

Per read, capture:

- **identity before**: `size`, `mtime_ns`, `volume_serial`, `file_id`, `in_use`
- digest, kind (B / C), process start time, elapsed
- **identity after**; abort the probe if anything changed
- for C: the tool's own accounting (`aligned_len`, `tail_len`, `unbuffered_bytes`, `buffered_bytes`, `read_calls`, `short_sector_aligned_reads`)

Output: `eval/runs/S1B-STAB-<UTCstamp>/read-<nn>.json` plus `summary.tsv` and `summary.md`. Write-once; the directory name carries the stamp.

**Cost:** 10 × 5.34 GB ≈ **53 GB of I/O**, roughly 15–25 minutes. No GPU use, no server, no VRAM.

---

## 4. Decision table — fixed before the data exists

Let `B₁…B₅` be the five buffered digests and `C₁…C₅` the five composite digests.

| Outcome | Reading | Consequence for S1-B |
|---|---|---|
| All `B` equal, all `C` equal, `B = C` | **stable and path-independent, in this session** | still not evidence about the past; manifest decision remains open |
| All `B` equal, all `C` equal, `B ≠ C` | **two stable streams, path-dependent** | reproduces the 19:16 shape; the slot is unusable until one stream is justified |
| Any `B` varies, or any `C` varies | **unstable input** | slot cannot be qualified; population decision required |
| Identity changes at any point | probe aborts; a writer is active | stop, report |
| Read error / short read / CRC | probe aborts | stop, report |

**None of these outcomes establishes which digest is the file's content.** A stable `b5b1913e` would not make it correct; a stable `d0027dd3` would not by itself validate the re-baseline, though it would remove the objection that the input is unstable.

---

## 5. Stop conditions

Abort, preserve what exists, do not retry if: identity changes; any read errors, short-reads, or returns `ERROR_CRC` / `ERROR_DEVICE_HARDWARE_ERROR`; a read exceeds 180 s; total elapsed exceeds 30 min; more than 12 reads would be required; or the file becomes `in_use`.

A stop is a result. A second attempt is a new authorisation.

---

## 6. Prohibitions

No manifest edit. No re-baseline. No slot qualification. No model load. No server start. No prompt. No cache purge. No USN query. No event-log query. No handle enumeration. No write of any kind to `G:`. No change to `eval/` source. No corpus, standard, tag or score change.

---

## 7. What the probe can and cannot establish

**Establishes:** whether the current file yields repeatable digests within one session, separately for a cached path and a cache-bypassing path, with identity pinned at each read.

**Cannot establish:** which stream is the artifact's true content; whether the earlier seven digests were reads of a changing file or a changing path; whether stability now implies stability on 5 or 14 November; or that the file is safe to use as experimental input. A stable result is a **precondition** for that argument, not the argument.

---

## 8. Authorisation

One sentence is enough: *run the S1-B stability probe*. On that, I run §3 exactly, file the evidence, and report §4's outcome without deciding the manifest question.
