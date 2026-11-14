# `STORAGE-INVESTIGATION-CLAIM-CORRECTIONS`

**Identity:** `HOWL-CLAIMCOR-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-03` Track 5
**Status:** correction memorandum. **No existing report was edited.** Every document named below stands as filed; this memorandum states what is wrong and what should replace it.

**Why a memorandum rather than edits:** the reports are inputs to an owner decision, and silently rewriting an input is how an audit trail loses its value. Several corrections below are to documents produced earlier in this same engagement, including two that are wrong because I asserted something the filed artefacts contradict.

---

## 1. Corrections

### C-01 — "the behaviour has not been measured on this machine"

| Field | Content |
|---|---|
| **Document** | `WINDOWS-CACHE-HYPOTHESIS-REVIEW-2026-11-14.md` §4, final paragraph ("The behaviour has not been measured on this machine and cannot be asserted either way"), and §2 step S-9a's negative column |
| **Claim** | that a buffered-versus-unbuffered comparison had never run |
| **Status** | **factually wrong** |
| **Actual** | it ran on 2026-11-14 at **19:16:03–19:17:12**. Result: **DIFFERENT** — buffered `341fb28c…7c35`, unbuffered `d0027dd3…6e96`, identical file id (`914f0100…`), identical volume serial (`880ecbb3…`), size and mtime unchanged, `in_use: false`, byte accounting COMPLETE. Filed as `eval/runs/S1B-G-BUFFERED-VS-UNBUFFERED-20261114T191603Z.txt`. An earlier attempt at 19:15:48 **aborted** on a missing intermediate file and is correctly recorded as `RESULT: ABORTED` |
| **Correction** | The documented claim — that unbuffered I/O reaches the physical disk — **held on this machine for this file**: the two paths diverged. The practitioner claim that a resident page would satisfy an unbuffered request was **not** observed. The test is therefore informative, not uninformative, and it is one of the two most important results in the record |
| **Basis** | the cited run file, read in full |

### C-02 — the DLL incident's source volume

| Field | Content |
|---|---|
| **Document** | `GGUF-FILE-PROVENANCE-MATRIX-2026-11-14.md` §4.1 and the §5 table row "A second copy is enough to change which bytes land on the destination — established (`cublas64_12.dll`)"; also `WINDOWS-CACHE-HYPOTHESIS-REVIEW` §6's DLL row and its §10 "directional reading" |
| **Claim** | the DLL came *from `G:`*, and therefore that two same-shape events both concerned reads of `G:`, and that this leaned the evidence "away from a deliberate writer and toward a read-path phenomenon" |
| **Status** | **factually wrong** on the source volume; the leaning built on it is **withdrawn** |
| **Actual** | the source was `H:\AI apps` **before that volume was reformatted** — exFAT, `HealthStatus: Warning`, `OperationalStatus: Full Repair Needed`. The destination was `D:\H-MIGRATION-2026-11-14`. `G:` was not involved. Verified against `H-MIGRATION-VERIFY-2026-11-14.json` (`"source": "H:\\AI apps"`) |
| **Correction** | The DLL incident is an unexplained copy discrepancy from a volume that was itself reporting a damaged filesystem. It demonstrates the *signature* (same size, same mtime, different content, both sides stable) outside a model path; it does not implicate `G:`, and it is not evidence about the model files' mechanism. The "directional reading" in `WINDOWS-CACHE-HYPOTHESIS-REVIEW` §10 must be struck |
| **Basis** | `H-MIGRATION-VERIFY*`, `H-MIGRATION-RECORD`, `H-RESTORE-VERIFY` |

### C-03 — the `H:` "later reading" contradiction

| Field | Content |
|---|---|
| **Document** | `GGUF-FILE-PROVENANCE-MATRIX` §4 (X-1/X-2/X-3) and `HISTORICAL-H-READ-RECONCILIATION` §2 (T-12/T-13) |
| **Claim as filed** | two readings of the same `H:` paths disagree, which is a contradiction inside the record with three possible resolutions |
| **Status** | **overstated** — the *later* values are not backed by any filed measurement |
| **Actual** | the only filed `H:` model measurements are one run for `S1-A`, **three identical** runs for `S1-B`, and one run for `S1-C`, each recorded `H_stable: true` (`H-MODEL-MIGRATION-CANDIDATES-2026-11-14.json`). The values `cd4b4430…` and `01c357cd…` appear **only** in the deferral record's table — no log, no JSON, no timeline |
| **Correction** | Reclassify T-12/T-13 from "conflicting observations" to **"unsourced table entries"**. The X-1/X-2/X-3 question survives but changes shape: it is now "were the `H:` copies written more than once, or does the deferral table contain values that were never measured?" Both are answerable from metadata/records, and neither is answered |
| **Basis** | `h-migration-candidates.log`, `H-MODEL-MIGRATION-CANDIDATES`, deferral record |

### C-04 — "a handle was never tested"

| Field | Content |
|---|---|
| **Document** | `GGUF-FILE-PROVENANCE-MATRIX` §3, row "Was a handle held open at any observed moment?"; `WINDOWS-11-STORAGE-EVIDENCE-MATRIX` §5 rank 1 |
| **Claim** | only the loader's mapping was ever present, and handle presence was never tested |
| **Status** | **incomplete** — a sharing probe exists |
| **Actual** | `handle-enumeration-2026-11-14.txt` records an **exclusive-open probe** (`GENERIC_READ`, `dwShareMode = 0`) on all six `G:`/`H:` candidate paths at **13:03:38**: 6/6 opened, 0 blocked |
| **Correction** | The record does contain one handle test. It shows that **at 13:03:38 no other handle was open** on any candidate. It does not show handle state at 13:45, 14:03–14:05, or 19:16, and an exclusive-open success is not enumeration (it proves absence by exclusion, reports no access mask, and briefly holds the file itself) |
| **Basis** | the probe file |

### C-05 — "filter drivers: unknown"

| Field | Content |
|---|---|
| **Document** | `WINDOWS-11-STORAGE-EVIDENCE-MATRIX` §1.2 row "Filter drivers in the stack"; §5 rank 4 ("Filter driver list") |
| **Claim** | the filter list is unknown and should be collected |
| **Status** | **wrong** — it was already collected |
| **Actual** | `filter-inventory-2026-11-14.txt`: on `G:` — `WdFilter` (Defender, altitude 328010), `FileInfo` (40500), `Wof` (40700), `bfs` (150000), `gameflt` (189850), `UCPD` (385250.5); `CldFlt` on `C:` only; **no third-party AV or backup minifilter present** |
| **Correction** | The filter list is **known**, and the plausible-agent set is narrower than the matrix implied: five Microsoft filters, none of which is documented to substitute file content. What remains unknown is their **configuration and versions**, and whether any of them is active against `G:` during idle periods |
| **Basis** | the inventory file |

### C-06 — the September-integrity inference

| Field | Content |
|---|---|
| **Document** | `docs/research/corpus/MEMORY-09-S1-DEFERRAL-DIGEST-REPRODUCIBILITY-2026-11-14.md` §2 ("The originals on `G:` are the more likely intact copies"); `MEMORY-09-INTEGRITY-STOP-S1B-DIGEST` §9.1 |
| **Claim** | baseline-matching digests indicate intact files |
| **Status** | **overstated** — conflates baseline agreement with integrity |
| **Actual** | the baseline was created at 13:45 on 2026-11-14 from the file as it stood that morning. The file's mtime is 2026-09-04. Nothing in the record establishes that the 2026-09-04 content was correct, and no publisher hash was ever obtained |
| **Correction** | A digest matching the baseline establishes **consistency with a 2026-11-14 reading**, not integrity and not provenance. "The more likely intact copies" should read "the copies whose digests are consistent with the only reference this machine has" |
| **Basis** | baseline creation record; absence of any external hash in the corpus |

### C-07 — Defender as an excluded mechanism

| Field | Content |
|---|---|
| **Document** | `eval/runs/H-MIGRATION-RECORD-2026-11-14.json`, `defender.conclusion`: "Defender is not the mechanism behind the byte differences" |
| **Claim** | Defender is excluded |
| **Status** | **overstated** |
| **Actual** | the evidence supports a narrower conclusion: no Defender detection or Controlled Folder Access block is recorded for these paths on the relevant dates, and CFA **denies** access rather than substituting content. `WdFilter` is nonetheless attached to `G:` and to `H:` at altitude 328010 |
| **Correction** | "No Defender detection or CFA block is recorded for these paths; CFA denies rather than substitutes; `WdFilter` remains attached and is not excluded as an observer" |
| **Basis** | `defender` block + filter inventory |

### C-08 — "no other mechanism has been tested, so this one leads"

| Field | Content |
|---|---|
| **Document** | `GGUF-FILE-PROVENANCE-MATRIX` §5 ("the strongest pattern in the record … a pattern about reads"), `WINDOWS-CACHE-HYPOTHESIS-REVIEW` §3.2 of the earlier matrix ("the only mechanism that fits") |
| **Claim** | a read-path reading leads because the alternatives are untested |
| **Status** | **invalid reasoning**, even though the conclusion is now partly supported by C-01 |
| **Actual** | the read-path reading is supported by a *positive observation* (O-10, the layer divergence), not by the failure of its competitors. Those are different arguments and only one is sound |
| **Correction** | State the support as: "the layers of the read path diverged at one instant; this establishes that the phenomenon is at least partly in the read stack, and identifies nothing further". Do not present the absence of writer evidence as support for a read-path cause |
| **Basis** | `GGUF-DLL-COMPETING-HYPOTHESES` §5 |

### C-09 — "control file stable therefore the path is healthy"

| Field | Content |
|---|---|
| **Document** | `MEMORY-09-INTEGRITY-STOP-S1B-DIGEST` §3 and §9.1 (correctly qualified there); `WINDOWS-11-STORAGE-EVIDENCE-MATRIX` §3.1 row 4 |
| **Claim** | the `LFM2.5-8B` control read consistently, so the hashing path is not at fault |
| **Status** | **correctly narrow in the source, at risk of being read wider** |
| **Actual** | a single control file reading consistently excludes "broken hashing tool" and "wrong file". It cannot exclude a localised bad region, a path-specific agent, page-local corruption, or a per-file cache state |
| **Correction** | Keep the narrow reading and add explicitly: the control certifies the tool and the path resolution, not the volume, not the device, not the cache |
| **Basis** | integrity record §9.1 (already says this) — restated so downstream summaries cannot widen it |

### C-10 — staging-plan steps that are already done

| Field | Content |
|---|---|
| **Document** | `STORAGE-DIAGNOSTIC-STAGING-PLAN-2026-11-14.md` S-2 (collect filter list), S-4 preamble ("`handle.exe` is not installed"), S-9a ("worth running eventually") |
| **Claim** | those steps are future work |
| **Status** | **superseded by filed artefacts** |
| **Correction** | S-2's filter list is **already collected** (C-05) — the remaining unknown is configuration, not membership. S-4 remains undone as *enumeration*, but a sharing probe exists (C-04). S-9a has **already run** with a positive result (C-01) and should be reclassified from "proposed" to "performed; result decisive for the family, not the member" |
| **Basis** | the three artefacts |

---

## 2. Statements that survive review unchanged

Not everything needs correcting, and saying so is part of the reconciliation:

| Statement | Stands because |
|---|---|
| "HOWL Chat wrote none of the S1 files" | source-visible; no chat-model write path exists |
| "Size and mtime constant while digests varied" | recorded in every observation |
| "The change occurred while nothing in the workflow was running" | observations 5→7, idle |
| "`NO_BUFFERING` does not affect the device's own cache or memory-mapped files" | `CreateFile` documentation |
| "A rename without `Sync` does not establish durability" | `FlushFileBuffers` documentation, and the writer audit's D-3 |
| "The `F-1` handle leak is real" | `usn_reader.py` contains no `CloseHandle` |
| "No acknowledged Windows 11 defect matches this signature" | KB5063878 was investigated and not confirmed; the device and symptom differ |
| "One copy pass cannot certify a tree on this machine" | `H-MIGRATION-VERIFY-CORRECTION`'s own note, confirmed by its own data |

---

## 3. Pattern in the corrections

Six of the ten corrections share one cause: **a claim was made from the narrative documents without checking the raw artefacts those documents summarised.** The narrative record is a compression, and three compressions were lossy in ways that changed conclusions — the unbuffered test disappeared from the story, the DLL's source volume was misremembered, and the filter inventory was treated as uncollected.

The procedural consequence, stated for future work orders: **every quantitative claim in a synthesis document should cite a file under `eval/runs/`, not another synthesis document.** Where no such file exists — as with the deferral table's later `H:` values — the claim should be marked unsourced at the point of writing, not discovered later.

---

## 4. Claim ceiling

**Establishes:** ten specific corrections with document, location, original wording, status, replacement and basis; and eight statements that survive unchanged.

**Does not establish:** that the corrected set is complete — it covers the documents named in the work order and the artefacts found while checking them. A further pass could find more, and the pattern in §3 suggests it would.
