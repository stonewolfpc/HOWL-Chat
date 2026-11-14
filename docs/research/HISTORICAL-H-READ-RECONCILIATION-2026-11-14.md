# `HISTORICAL-H-READ-RECONCILIATION`

**Identity:** `HOWL-HREC-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-03` Track 1
**Status:** reconciliation only. **No volume accessed, no file read, no digest computed.**
**Sources:** `eval/runs/` raw artefacts (`H-MODEL-MIGRATION-CANDIDATES`, `h-migration-candidates.log`, `H-MIGRATION-RECORD`, `H-RESTORE-VERIFY`, `H-MIGRATION-VERIFY*`, `geometry-G/H.txt`, `handle-enumeration-2026-11-14.txt`, `filter-inventory-2026-11-14.txt`, `S1B-G-BUFFERED-VS-UNBUFFERED-*`), plus the integrity and deferral records.

**The finding that governs this document:** the deferral record's `H:` "earlier → later" pairs have **no supporting run record**. The only filed `H:` measurements are one run for `S1-A`, **three identical** runs for `S1-B`, and one run for `S1-C` — each recorded `H_stable: true`. The later values `cd4b4430…` and `01c357cd…` appear in exactly one place: the deferral record's table. They are not in any log, JSON artefact, or timeline.

---

## 1. Chronology

Identity columns are populated only where a filed artefact recorded them. `—` means the observation was recorded without it; that absence is itself a finding, because it makes every transition below less classifiable than it should be.

| # | When | Path | Size (B) | mtime | Volume serial / file id | SHA-256 | What it was | Copy between? |
|---:|---|---|---:|---|---|---|---|---|
| 1 | 13:03:38 | all six `G:`/`H:` candidates | — | — | — | — | **exclusive-open sharing probe** (`dwShareMode = 0`), 6/6 opened, 0 blocked | — |
| 2 | 13:45 | `G:` S1-B | 5,335,285,728 | 2026-09-04 20:53:48 | — | `b5b1913e…ede0` | baseline creation, server idle | — |
| 3 | 13:50 | `G:` S1-B | same | same | — | `b5b1913e…ede0` | qualification v1 — **matched baseline** | — |
| 4 | 14:03 | `G:` S1-B | same | same | — | `f46b1049…9dcd2` | requalification, **server running** | — |
| 5 | 14:03 | `G:` S1-B | same | same | — | `e6dca6eb…51e23` | second read, server running | — |
| 6 | 14:04 | `G:` S1-B | same | same | — | `e6dca6eb…51e23` | server **stopped**, read | — |
| 7 | 14:04 | `G:` S1-B | same | same | — | `e6dca6eb…51e23` | server stopped, read again | — |
| 8 | 14:05 | `G:` S1-B | same | same | — | `eac8960a…2e3d0` | server stopped, read | — |
| 9 | after #2 (baseline referenced as `recorded_prefix`) | `G:` S1-A | 5,627,044,224 | — | — | `0e124fab…683f` | migration-time read of source | copy `G:`→`H:` running |
| 10 | same window | `G:` S1-B | 5,335,285,728 | — | — | `290f7ead…2d4` | migration-time read of source | copy running |
| 11 | same window | `G:` S1-C | 5,155,565,440 | — | — | `f820ef3f…5c7e` | migration-time read of source — **equals baseline** | copy running |
| 12 | same window | `H:` S1-A | 5,627,044,224 | — | H: serial `3c3871ba…` (`geometry-H`) | `de3b8d28…53b4` | post-copy verification, **1 run** | copy complete |
| 13 | same window | `H:` S1-B | 5,335,285,728 | — | same | `6ff1b33f…f2f6` ×3 | post-copy verification, **3 runs, identical** | copy complete |
| 14 | same window | `H:` S1-C | 5,155,565,440 | — | same | `7fca804a…2895` | post-copy verification, 1 run | copy complete |
| 15 | later | `D:` backup of apps | — | — | — | — | restore source (78,259 files, 0 mismatches vs restored `H:`) | restore `D:`→`H:` |
| 16 | 19:16:03–19:17:12 | `G:` S1-B | 5,335,285,728 | mtime_ns `1788573228448965700` | **`880ecbb30ecb9918` / `914f0100…0000`** | buffered `341fb28c…7c35`, unbuffered `d0027dd3…6e96` | **buffered vs unbuffered, one run; file not in use** | none |

`geometry-G`: logical sector **512**, physical **4096**, serial `0x880ecbb30ecb9918`, NTFS 3.1. `geometry-H`: logical 512, **physical 512**, serial `0x3c3871ba3871742c` — a *different* serial from `G:`, so the two paths are different volumes by construction.

**Ordering evidence.** Rows 9–14 post-date row 2 because the migration log compares against `recorded_prefix` values that only exist once the baseline was written. Rows 12–14 therefore represent reads made *after* the 13:45 baseline. That places the migration inside the same day but **after** the 14:05 readings — the record contains no timestamp for it, which is a gap, not an ordering.

---

## 2. Transitions, classified

| # | Transition | Class | Basis |
|---:|---|---|---|
| T-1 | #2 → #3 (`b5b1913e` = `b5b1913e`) | **2. Same-file, same-read** | repeat matched |
| T-2 | #3 → #4 (`b5b1913e` → `f46b1049`) | **2. Same-file, different read** | same path, same size, same mtime; loader held the file |
| T-3 | #4 → #5 (`f46b1049` → `e6dca6eb`, same minute) | **2. Same-file, different read** | consecutive reads, server still running |
| T-4 | #5 → #6 → #7 (`e6dca6eb` ×3) | **2. Same-file, same read** | three consecutive matches, idle |
| T-5 | #7 → #8 (`e6dca6eb` → `eac8960a`, idle) | **2. Same-file, different read** | nothing attached; the decisive idle transition |
| T-6 | #2 → #9 (`b5b1913e` → `0e124fab`, `G:` S1-A baseline → migration read) | **2. Same-file, different read** | baseline digest reproduced nowhere in the migration window |
| T-7 | #2 → #10 (`b5b1913e` → `290f7ead`, `G:` S1-B) | **2. Same-file, different read** | same path, size constant |
| T-8 | #11 (`f820ef3f` at migration = baseline) | **2. Same-file, same read** | `S1-C` consistent throughout on `G:` |
| T-9 | #9 → #12 (`0e124fab` → `de3b8d28`, `G:`→`H:`) | **4. Unresolved** | source and destination digests differ. A copy that read correctly would transfer `0e124fab`; a copy that read badly would transfer something else. Both are consistent with the record. No file id or creation time was captured on either side |
| T-10 | #10 → #13 (`290f7ead` → `6ff1b33f` ×3) | **4. Unresolved** | same reasoning; `H:` stable over three runs only proves the destination was stable *after* the copy |
| T-11 | #11 → #14 (`f820ef3f` → `7fca804a`) | **4. Unresolved**, but bounded | `cmp` localises the first difference at **byte 58,444,208** — a real content difference between the volumes, in a file whose source digest matched the baseline |
| T-12 | #13 → deferral's `H:` "later" `01c357cd` | **4. Unresolved — and unsupported** | no filed run produces `01c357cd`; `H-MODEL-MIGRATION-CANDIDATES` reports `H_stable: true` from three identical runs |
| T-13 | #12 → deferral's `H:` "later" `cd4b4430` | **4. Unresolved — and unsupported** | same: one filed run only (`de3b8d28`), no second run recorded |
| T-14 | #8 → #16 buffered (`eac8960a` → `341fb28c`) | **2. Same-file, different read** | hours apart, buffered path |
| T-15 | #16 buffered vs #16 unbuffered (`341fb28c` vs `d0027dd3`, **same minute, same identity, not in use**) | **2½ — a new subtype: same-file, same-moment, different access path** | this is not "different read"; it is two read *mechanisms* disagreeing at one instant |

**No transition in the record is class 1 (confirmed re-copy or rewrite) or class 3 (confirmed path/identity change).** Nothing was observed being written; nothing was observed being replaced.

### 2.1 Why T-15 is a different animal

Every other transition can be explained by "the state changed between two moments". T-15 cannot: the file id (`914f0100…`), volume serial, size, mtime and not-in-use flag were identical before, between and after both reads, and the two reads were ~19 s and ~49 s apart in the same minute. One of the two paths returned bytes the other did not have, **at the same time, from the same file object.**

That is a statement about the storage stack, not about the file's history.

### 2.2 Caveat on T-15's attribution

The comparison is not perfectly clean: the unbuffered digest covers 5,335,281,664 bytes read unbuffered **plus 4,064 trailing bytes read buffered** (the file is not a multiple of the 4096-byte block). The tool accounts for this (`accounting: … COMPLETE`, `short_sector_aligned_reads: 0`), so the coverage is total — but a difference confined to the final 4,064 bytes would be attributed to the wrong path. 4,064 bytes is 0.00008 % of the file; it is a real confound in *attribution*, not in *existence*.

Alignment is sound: `G:` has logical sector 512 and physical 4096; the tool used 4096-byte blocks and reported zero short sector-aligned reads. #16's unbuffered path therefore met the documented alignment requirement.

---

## 3. What each observation's missing evidence would settle

| Missing evidence | Would settle | Obtainable without reading payload? |
|---|---|---|
| **Creation time of each `H:` copy** | whether the `H:` models were written once or twice — i.e. X-1 vs X-2 | **yes**, `H:` metadata only |
| A second `H:` model-copy run record | directly the same question, from the documents | yes — but the record contains none |
| Provenance of the deferral table's later `H:` values | whether T-12/T-13 are observations at all or transcription | yes — a records question, not a storage one |
| File id / volume serial for observations #2–#14 | whether every read was of the same object (identity confusion, hard links, restored files) | no — requires re-observation; only #16 carries identity |
| Timestamp for the migration run | exact ordering of #9–#14 against #4–#8 | yes, if any log carries it (none found) |
| Publisher hashes for `S1-A`/`S1-B`/`S1-C` | whether the baseline ever matched the file as distributed | **yes** — public metadata, no local storage |

---

## 4. What the reconciliation establishes

| Statement | Class |
|---|---|
| `G:` S1-B produced **six distinct digests** across 13:45–19:16 with constant size and mtime | **confirmed observation** (filed runs) |
| `G:` S1-A produced at least three distinct digests (`b5b1913e`-era baseline `2ca636d9`, migration `0e124fab`, later `2ca636d9`) | **confirmed observation** |
| `S1-C` on `G:` has been stable at its baseline in every filed observation | **confirmed observation** |
| The same `S1-B` file returned **different bytes through buffered and unbuffered reads within one minute**, with identical file identity and no holder | **confirmed observation** |
| The `H:` copies differ from their `G:` sources at migration time | **confirmed observation** |
| The `H:` copies were produced by a copy that serialised a different byte stream than the concurrent `G:` read | **consistent, not established** — the source digest at copy time is not the same as the digest recorded during the copy; a single robocopy report cannot distinguish "read badly" from "wrote badly" |
| The `H:` digests later changed | **not supported by any filed run** |
| Any write, re-copy or replacement occurred after 2026-09-04 | **not established** |

## 5. Claim ceiling

**Establishes:** the full chronology as the artefacts support it, the classification of every transition, the three unresolved `H:` transitions, and the new subtype (T-15) that no file-history explanation covers.

**Does not establish:** why any digest differed; whether the `H:` copies or the `G:` originals are faithful; whether the deferral table's later `H:` values are observations. The `H:` contradiction remains **unresolved**, and now more narrowly: it is either an undocumented second copy run (X-1) or an unsourced table entry — with the filed evidence favouring the latter and nothing proving it.
