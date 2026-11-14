# `S3-METADATA-TEST-RECONCILIATION`

**Identity:** `HOWL-S3REC-1` · **Date:** 2026-11-14 · **Scope:** documentation-only review of staging step S-3 against the actual migration flags and the available historical metadata.
**Status:** reconciliation. **No volume accessed, no file read, no metadata query, no test run.**
**Outcome in one line:** S-3 as written **cannot discriminate** X-1 from X-2 for the tool that actually performed the copy, and the repository contains a **fifth digest dated 2026-11-05** that changes the incident's timeline.

---

## 1. The copy tool is identified, and so are its flags

The models were not placed on `G:` by a download client. They were **copied from a retiring `E:` drive** by a script in this repository.

| Evidence | Location | Content |
|---|---|---|
| Copy script | `build/move-to-g.sh` L92 | `robocopy 'E:\AI GGUF\GGUF' 'G:\AI GGUF\GGUF' /E /J /COPY:DAT /DCOPY:DAT /R:2 /W:2 /XJ /NP /NFL /NDL` |
| Script purpose | `build/move-to-g.sh` L1–L4 | "one serialised, resumable, fully logged move of the AI model trees off the retiring `E:` drive onto `G:`" |
| Same pattern elsewhere | `build/copy-modpack.sh` L13–L20 | `/E /J /COPY:DAT /DCOPY:DAT /R:2 /W:2 /XJ /NP /NFL /NDL`, with the comment "**/J unbuffered**" |
| App-tree copy pattern | `build/finish-move.sh` L36–L37 | `/J /COPY:DAT /R:2 /W:2 /NP /NFL /NLD` |

So the flags are **recorded**, contrary to what the staging plan asserted ("the migration's flags are not recorded"). Four consequences follow, and each one narrows what metadata can prove.

### 1.1 `/J` means the ingestion writer was unbuffered

`/J` is robocopy's unbuffered-I/O switch. The recorded migration therefore ran in unbuffered mode — the same access mode whose behaviour is at issue in T-15. That does not make the copy a suspect for a November divergence, but it does mean the machine's own tooling is configured to write large files through the unbuffered path, and reading D-1 (stale cache after an unbuffered write) as exotic would be wrong.

**Qualification, retained deliberately:** `/J` establishes the *configured* robocopy mode. It does not establish that every operation in the ingestion chain bypassed every cache — a filter, a driver or a second tool could sit between — and it does not establish that the resulting data was correct. Unbuffered copying is a statement about the requested access mode, not a guarantee about the bytes or about what every layer did.

### 1.2 `/COPY:DAT` preserves the modification time, so mtime cannot date the copy

Every model on `G:` carries an mtime inherited from `E:`. `S1-B`'s `2026-09-04 20:53:48` is therefore **not** when the file arrived on `G:`. It is a timestamp the destination inherited from `E:`. `/COPY:DAT` explains its preservation; it does not establish that the timestamp corresponds to the last actual content-changing write, because timestamps can also be preserved or restored by earlier operations. The staging plan already said the baseline is not a provenance anchor; this makes the same point about the mtime itself.

Whether `/COPY:DAT` also preserves **creation** time is version- and flag-specific. The plan must not assume it either way: measuring it is the only way to know, and that measurement cannot distinguish the cases below regardless.

### 1.3 A robocopy re-run is metadata-invisible, which is the fatal objection

The script documents its own resumption rule, in its own words:

> "a rerun resumes because **robocopy skips files whose size + timestamp already match**"

Robocopy's change detection is size + timestamp. It does not hash. Therefore:

- a re-run over an unchanged-looking file **does nothing at all** — no write, no metadata change, no log line beyond a skip;
- a re-run **skips** a source/destination pair whose size and timestamp match, **regardless of whether their contents match**. If a re-copy does occur because the metadata comparison does not match, `/COPY:DAT` can preserve the source timestamp on the destination. Consequently, matching metadata cannot establish whether an earlier re-copy occurred.

**This is the crux.** The one re-copy hypothesis that matters (X-1: the `H:` or `G:` file was written twice) produces exactly the metadata signature of X-2 (no re-copy at all). Creation time, mtime, size and directory timestamps are all identical in both worlds. So S-3, as a *discriminator*, cannot work against this tool.

### 1.4 The copy log exists in principle and is absent in fact

`build/move-to-g.sh` logs to `build/logs/move-to-g.log` and writes `build/logs/e-manifest.tsv`, plus `G:\AI GGUF\_move-status.txt`. **`build/logs/` does not exist in this repository** (glob returns nothing), so the documentary route — "did a second robocopy run?" — is closed from here. The manifest that would have listed every file's size and mtime at copy time is likewise gone.

If `e-manifest.tsv` survives anywhere — it was copied to `G:\AI GGUF\e-manifest.tsv` by the script itself — it would give the *expected* size and mtime at ingestion for all 149 files. That is a metadata read of one small text file, and it is a stronger artefact than anything S-3 proposed to collect.

---

## 2. A fifth digest, nine days earlier

`build/unit-c/chat-model-verification.md` L47–L51 records, for the **same path with matching reported metadata** (that record carries no file ID, so identity across the two dates is not established):

| Field | Value |
|---|---|
| File | `G:\AI GGUF\GGUF\Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M(1).gguf` |
| Date | **2026-11-05** (Unit C, owner-accepted, closed) |
| Size | `5,335,285,728` B |
| mtime | `2026-09-04 20:53:48 -0500` |
| SHA-256 | **`7096c0dd243da9762247386d476801714971e53e6894b4001710287a522fea1b`** |
| Structural check | 720 tensors, embd 2560, ftype 15, `data_ok: true`, `problems: []` |

Compare the November 14 baseline for the same file: `b5b1913e…ede0`, same size, same mtime, same structure.

**This is the most consequential artefact found in this review.** It establishes:

1. **Two different digest observations are recorded nine days apart.** `7096c0dd…` (Nov 5) ≠ `b5b1913e…` (Nov 14), with the reported path, size and mtime matching. The interval between them is unobserved; the records do not establish when, or whether, the on-disk content changed.
2. **The earlier observation predates every November 14 actor.** No qualification run, server load, migration or baseline edit is recorded for Nov 5. That dates the observation and nothing more. The earlier phrasing "nothing in the November 14 workflow caused it" **is withdrawn**: the Nov 5 record does not tell us what caused either reading.
3. **The "baseline" was one reading among at least seven** for this file — `7096c0dd`, `b5b1913e`, `f46b1049`, `e6dca6eb`, `eac8960a`, `341fb28c`, `d0027dd3` — and no two of them can all be the file's content.
4. **The structural checks passed for the November 5 byte stream** (720 tensors, `embd 2560`, `ftype 15`, `data_ok: true`, `problems: []`). That is compatible with differences that preserve the checked GGUF structure. It does not locate the differing bytes, and it does not exclude replacement by a structurally compatible GGUF.
5. **What it does not establish.** Two observations nine days apart. It does not show that the file changed on 6 November, that the read path was stable in between, or that the same mechanism produced both readings. The interval is unknown; only its endpoints are recorded. Every downstream summary must say "two observations nine days apart", never "the file was different for nine days".

The staging plan's framing of the incident as a November 14 event is therefore **too narrow**. The two observations are nine days apart, and the intervening period is unobserved, so neither persistence nor the timing of any transition between returned byte streams is established. The hypothesis table should accommodate observations separated by days, without treating that interval as evidence of continuous behaviour.

### 2.1 Why this does not establish a writer

A writer remains one possible explanation, but the seven recorded digests do not establish seven content changes, or any content change at all. They are **digest observations, not a write history**. The observations made while the machine was reportedly idle do not establish that a write occurred during those periods. A writer hypothesis must therefore remain open without being assigned a number of writes or a required write cadence. The available evidence does not distinguish writes from read-path divergence.

---

## 3. Revised S-3

| Aspect | Original S-3 | Corrected |
|---|---|---|
| Claim | creation time post-dating a reading ⇒ a re-copy occurred | **withdrawn as a discriminator.** For a robocopy re-run, creation time, mtime and size are all consistent with *both* "re-copy happened" and "no re-copy happened" |
| What it can still do | — | report the metadata that exists, and compare `H:` against the `D:` staging tree and against `G:\AI GGUF\e-manifest.tsv` if that file survives |
| Classification | automatic proof of X-1 or X-2 | **supports X-1 / supports X-2 / supports X-3 / inconclusive**, with inconclusive as the expected outcome |
| Best available artefact | none proposed | `G:\AI GGUF\e-manifest.tsv` (written by the script, size + mtime per file at ingestion) and `G:\AI GGUF\_move-status.txt` |

**What would actually discriminate X-1 from X-2** — none of it is metadata:

| Route | Why it works | Gate |
|---|---|---|
| A publisher hash for each candidate | external reference; independent of every local read | none — public metadata, no storage |
| A content-derived marker that survives copying but changes on rewrite (e.g. a second hash over a *different* range) | still a read; does not distinguish | none |
| USN history for the path over the window | records writes, if the journal retains them | **closed gate** |
| A known-pattern scratch file written and re-read (the V-B probe in `S9A-RANGE-PROBE-SPECIFICATION`) | gives ground truth for the *mechanism*, not for this file's history | **gate 1** |
| Robocopy logs / manifest comparing expected vs actual at ingestion | documentary, if they survive | none |

The honest conclusion: **X-1 versus X-2 is not resolvable by metadata on this machine.** It is resolvable by an external hash, by journal history, or not at all. The plan must say that rather than schedule a step that cannot answer its own question.

---

## 4. September 4 — what the record now says

| Question | Answer | Basis |
|---|---|---|
| Were the three GGUFs on `G:` before the Windows reinstall? | **unknown.** No install date is recorded anywhere in the repository | absence |
| Where did they come from? | **`E:\AI GGUF\GGUF`**, a retiring drive, copied to `G:` by `build/move-to-g.sh` | script L92 |
| When was the copy made? | **on or before 2026-11-05** — Unit C audited `G:` and hashed the file that day | Unit C record |
| What does `2026-09-04 20:53:48` date? | the **recorded** September 4, 2026 mtime, inherited from `E:` and preserved through `/COPY:DAT`. Whether it corresponds to the last actual content-changing write is unknown | flag semantics |
| Was the file downloaded, or restored from a backup? | **unknown.** No downloader, client, URL or archive is recorded for any of the three | absence |
| Does the CUDA reinstall bear on this? | **no evidence either way.** CUDA working proves the toolkit is functional, not that a given file or read path is | reasoning |

One observation is available and worth stating: the recorded mtime is roughly **nine weeks before** the 5 November observation and **ten weeks before** the 14 November one (62 and 71 days), and the copy preserved it. The recorded mtime is September 4, 2026, inherited from `E:`; `/COPY:DAT` explains its preservation on `G:`, and the record does not establish whether that timestamp corresponds to the last actual content-changing write. A Windows reinstall on `C:` in that window would not, by itself, touch a data drive's files; a reinstall that reformatted or restored `E:` would. **Nothing in the record says which happened.**

---

## 5. What this review changes

| Change | Consequence |
|---|---|
| The copy tool and flags are **known**, not unknown | the staging plan's "flags are not recorded" is wrong; correct it before the plan is used |
| `/J` = unbuffered | the machine's normal bulk-write path is unbuffered; D-1 is not an exotic mechanism here |
| robocopy skips on size+time | **metadata cannot detect a re-copy**; S-3 is reclassified as non-discriminating |
| `build/logs/` absent | the documentary route is closed from this repository; `G:\AI GGUF\e-manifest.tsv` may survive and is the best remaining artefact |
| `7096c0dd…` dated Nov 5 | two digest observations nine days apart; the hypothesis table and resumption criteria must accommodate a timescale of days rather than hours |
| No install date recorded | the September 4 question cannot be answered from documents; it needs either an owner statement or an OS install-date query, which is a live system query and not authorised |

---

## 6. Claim ceiling

**Establishes:** the exact ingestion path and flags; that metadata cannot distinguish a robocopy re-copy from no re-copy; that a fifth digest was recorded nine days earlier for the same path with matching reported size, mtime and structure; and what the record does and does not say about September 4.

**Does not establish:** which of the seven digests is the file's content, whether the on-disk content changed at any point between 5 and 14 November, whether any write occurred after early September, whether the two dated records concern the same file identity (no file ID was recorded), or whether the files predate the reinstall. It removes a planned step from the diagnostic path rather than adding one.
