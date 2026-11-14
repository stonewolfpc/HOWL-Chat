# `GGUF-FILE-PROVENANCE-MATRIX`

**Identity:** `HOWL-PROV-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-02` Track 1
**Status:** reconciliation only. **No volume accessed, no handle opened, no digest computed, no model file touched.**
**Inputs:** the writer audit, the integrity record, the execution audit, the deferral record, `eval/runs/H-MIGRATION-RECORD-2026-11-14.json`, and the filed qualification records.

**Rule applied:** a fact is recorded with the class of its support. Nothing is filled in from plausibility. Where the record is silent, the entry says **unknown** — including where a likely answer exists and nobody wrote it down.

**Fact classes:** **DOCUMENTED** (written in a filed record) · **REMEMBERED** (owner/session statement, not filed) · **INFERRED** (follows from documented facts) · **UNKNOWN** (not in the record).

---

## 1. The two events that must not be merged

The single most common error available here is to treat "the file has digest X" as one fact. It is at least four:

| Event | What it is | Evidence available | Status |
|---|---|---|---|
| **E-CREATE** | the file was first written to `G:\AI GGUF\GGUF\` | mtime `2026-09-04 20:53:48` on `S1-B`; nothing else | **largely unknown** |
| **E-COPY** | copies were made from `G:` to `H:` on 2026-11-14 | `H-MIGRATION-RECORD`, robocopy 149 files / 428.931 GiB / 1 h 30 m / 0 failed | **documented** |
| **E-READ** | the file was hashed, repeatedly, across sessions | four distinct digests for `S1-B`; a baseline match later for `S1-A` | **documented** |
| **E-WRITE?** | *if it happened at all*, something changed bytes after creation | inferred only from E-READ disagreement | **not established** |

Every conclusion in this investigation depends on which of these a given observation belongs to, and the record does not always allow the assignment. §4 is where that bites.

---

## 2. Provenance matrix — the three candidates

| Field | `S1-A` | `S1-B` | `S1-C` |
|---|---|---|---|
| **File name** | `Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M.gguf` | `Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M(1).gguf` | `LFM2.5-8B-A1B-heretic.Q4_K_M.gguf` |
| **Size** | 5,627,044,224 B | 5,335,285,728 B | 5,155,565,440 B |
| **Location in use** | `G:\AI GGUF\GGUF\` | same | same |
| **Original source (publisher, repo, URL)** | **unknown** | **unknown** | **unknown** |
| **Downloaded by which application** | **unknown** | **unknown** | **unknown** |
| **Download date** | **unknown** | **unknown**; mtime `2026-09-04 20:53:48` is the only timestamp on record | **unknown** |
| **Quantiser / uploader** | **INFERRED** from the filename: a community re-quant, not an upstream release | same | same (`heretic` suffix, also community) |
| **Recorded baseline digest** | `2ca636d9…` | `b5b1913e…` | `f820ef3f…` |
| **Baseline creation time and method** | 13:45, harness baseline over the `G:` path | same | same |
| **Written directly or via temp+rename** | **unknown** — no record of the creating write | **unknown** | **unknown** |
| **Copied after creation?** | yes — `G:`→`H:` on 2026-11-14 (robocopy) | yes | yes |
| **`H:` copy digest, migration time** | `de3b8d28…` | `6ff1b33f…` (stable over 3 reads) | `7fca804a…` |
| **`G:` == baseline at migration** | **no** (`0e124fab…`) | **no** (`290f7ead…`) | **yes** |
| **Later `G:` reading** | `2ca636d9…` = **baseline** | `341fb28c…` ×2 | `f820ef3f…` |
| **Later `H:` reading** | `cd4b4430…` | `01c357cd…` | `7fca804a…` |
| **Functional state** | loads, generates on both copies | same | same |

### 2.1 What the naming pattern does and does not tell us

All three are quantised GGUF builds by community uploaders rather than first-party releases. That is an **inference from filenames**, and it has one investigative consequence and no others: a community re-quant may have been uploaded in a broken state, so "the file was already wrong when it arrived" is a live hypothesis that no digest can exclude, because the baseline was created *from that same file* at 13:45 on 2026-11-14, weeks after the recorded mtime.

**The baseline is not a provenance anchor.** It is a reading of the file as it existed that morning. If the file arrived damaged on 2026-09-04, every "match" since then certifies internal consistency, not correctness.

---

## 3. Access surface over the file

| Question | Answer | Class |
|---|---|---|
| Which applications could read the path? | any process running as the user; the volume is a fixed local NTFS volume with default ACLs | **INFERRED** |
| Which application *did* read it? | KoboldCpp 1.120 (every launch, `--usemmap`), the harness (`file_digest`, `sha256sum`), robocopy during migration | **DOCUMENTED** |
| Which application could **write** it? | any process running as the user with write access to `G:\AI GGUF` — the volume's ACLs were not recorded | **UNKNOWN** |
| Was a handle held open at any observed moment? | only the loader's mapping, for observations 2–4; **nothing** for 5–7 | **DOCUMENTED** |
| Was a write-capable handle ever observed? | **never** — handle enumeration was never performed (`handle.exe` unavailable) | **DOCUMENTED as a gap** |
| Backup / sync / AV / indexer interaction with `G:` | OneDrive **excluded** (root on `C:`); Defender detections were on `C:` executables, not models; Controlled Folder Access **denies** rather than substitutes; no backup product identified; volume label is *"Second backup"*, untraced | **DOCUMENTED for the exclusions; UNKNOWN for the rest** |
| Scheduled task referencing `G:` | none found | **DOCUMENTED** |

---

## 4. The reconciliation problem this matrix exposes

The record contains **two readings of the same `H:` paths that disagree**:

| Path | Migration-time reading | Later reading | Consistent? |
|---|---|---|---|
| `S1-A` on `H:` | `de3b8d28…` | `cd4b4430…` | **no** |
| `S1-B` on `H:` | `6ff1b33f…` (recorded stable over 3 reads) | `01c357cd…` | **no** |
| `S1-C` on `H:` | `7fca804a…` | `7fca804a…` | yes |

Two readings of an `H:` file that nobody is recorded as having re-copied cannot both be right. There are exactly three explanations, and **the record does not choose between them**:

| # | Explanation | What it requires | Can the record settle it? |
|---|---|---|---|
| **X-1** | The `H:` files were **re-copied** between the two readings, and the second copy captured a different read of `G:` | an undocumented re-copy | **no** — no re-copy is filed. Resolvable from `H:` file mtimes/creation times, which requires `H:` access |
| **X-2** | The `H:` files were **not** re-copied, so `H:` reads also varied over time | a second volume producing varying digests | **no** — resolvable only by a new read, which is gated |
| **X-3** | The two readings are of **different files** (path, name or directory changed between them) | a rename or a second directory | **no** — resolvable from directory metadata, gated |

**Why this matters more than anything else in the matrix.** Under **X-1**, every anomaly belongs to the `G:` read path and the copy that serialised it — one suspect volume. Under **X-2**, two physically different drives (a WD SATA HDD on `G:`, a Crucial MX300 SATA SSD on `H:`) both returned varying digests for files that were stable by size and mtime, which points away from any single device and toward the host — memory, the read path, or whatever produced the bytes both times. Those two worlds imply different next steps, and the record cannot currently distinguish them.

**Therefore:** no mechanism may be ranked until X-1/X-2/X-3 is resolved, and resolving it is a metadata question, not a payload question (see `STORAGE-DIAGNOSTIC-STAGING-PLAN` step S-3).

### 4.1 A second, independent anomaly of the same shape

During the same migration, `AI apps/llama.cpp-prism-b9596-cuda12/runtime/cublas64_12.dll` — 78,259 files copied, **one** mismatch — arrived with **the same size (100,033,536 B), the same mtime, and a different SHA-256**, and **both sides were stable over three reads**. It was re-copied, and the second copy matched.

That is the same signature as the models: *metadata identical, content different, reads internally stable*. It occurred on a different file type, a different destination, and — critically — it was **resolved by a second copy**, which means at least one of the two robocopy passes read the source differently. A single mis-read during a copy is the simplest explanation, and it is a **read-path** phenomenon, not a writer.

Two same-shape anomalies, both involving reads of `G:` producing bytes that disagree with another read of the same path, with metadata untouched. That is the strongest pattern in the record — and it is a pattern about **reads**, not about writes.

---

## 5. What is established, and what is not

| Statement | Status |
|---|---|
| Neither the publisher nor the downloader of any S1 file is recorded anywhere in this repository | **established (by absence)** |
| The files were copied `G:`→`H:` by robocopy on 2026-11-14 after `H:` was quick-formatted to NTFS | **established** |
| HOWL Chat wrote none of them | **established** (writer audit, source-visible) |
| A second copy is enough to change which bytes land on the destination | **established** (`cublas64_12.dll`) |
| Reads of `G:` have produced at least four different contents for one path | **established** |
| Reads of `H:` have produced two different contents for one path | **apparent, conditional on X-1** |
| Something wrote to the file after 2026-09-04 | **not established** |
| Something *did not* write to it | **not established** |
| The files are correct as published | **unknown — and unknowable from local digests alone** |

### 5.1 The clean path that already exists

The migration record already names it: *"a fresh official download with its hash recorded as the new baseline is the clean path if one is ever needed"*. That is the only route that establishes provenance, because it introduces an external reference point — the publisher's own hash — which no local read can supply.

---

## 6. Ranked provenance questions

| Rank | Question | Why it is ranked here | Answerable without storage access? |
|---:|---|---|---|
| 1 | Were the `H:` copies re-made between the two recorded readings? (X-1/X-2/X-3) | splits the entire hypothesis space in two | **no** — needs `H:` metadata |
| 2 | Which program downloaded the S1 files, and is its hash known? | the only external anchor available | **partly** — a Hugging Face page/file card is public, needs network, not storage |
| 3 | Was the file already wrong when it arrived on 2026-09-04? | decides whether "baseline" means anything | **no** — needs question 2 answered first |
| 4 | Do the volume's ACLs let a non-owner process write? | bounds the writer population | **no** — volume metadata |
| 5 | Did any backup product ever target the "Second backup" volume? | cheapest remaining writer hypothesis | **partly** — installed-software inventory is not storage access |

---

## 7. Claim ceiling

**Establishes:** the provenance record as it actually stands, with every gap marked; the separation of creation, copy, read and possible write; and the X-1/X-2/X-3 ambiguity that currently prevents any mechanism from being weighted.

**Does not establish:** who wrote the files, when, with what tool, or whether they were correct on arrival. It does not establish that nothing wrote to them, and it does not establish that something did.
