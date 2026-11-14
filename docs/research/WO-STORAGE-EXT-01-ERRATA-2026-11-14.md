# WO-STORAGE-EXT-01 — Errata Record

**Identity:** `HOWL-EXT-HASH-01-ERRATA` · **Date:** 2026-11-14 · **Mode:** documentation only
**Target of record:** `docs/research/WO-STORAGE-EXT-01-EXTERNAL-HASH-RESEARCH-2026-11-14.md`, identity `HOWL-EXT-HASH-01`, **FROZEN**, sha256 `cd7366fe82620201284db037dce90a6f9eb77c349732b3704bb3d21049253ff6`.

**The frozen document has not been edited.** These are candidate corrections held for a possible `HOWL-EXT-HASH-02`, recorded so they are not lost and not silently applied.

---

## 1. Why an errata file rather than an edit

A frozen artefact that is quietly amended is not frozen, and a correction that exists only in a conversation is not recorded. Both failure modes are worse than a short separate note. This file therefore:

- takes no position on whether a revision will ever be made;
- does not restate the findings, only the wording defects;
- leaves the frozen hash valid and verifiable.

**Neither item changes a single value in the report.** The matched digest, the byte size, the provenance class, and the three separated conclusions are unaffected.

---

## 2. E-1 — the matching observation is a composite read, and three short references say otherwise

The frozen report states the composite nature correctly in §4.4 and in §6.1: the body (5,335,281,664 B, 1,273 sector-aligned calls, zero short reads) was unbuffered, and the final 4,064 B were read through the buffered path because the file is not a whole multiple of the 4,096-byte sector. Three shorter references omit that qualification.

| Location | Current text | Correction for a future revision |
|---|---|---|
| §4.3 table, row 1 | "Unbuffered `d0027dd3…` (19:16)" | "**19:16:03 composite read (unbuffered body, buffered 4,064-byte tail)** — `d0027dd3…`" |
| §4.4 bullet | "it establishes that one **unbuffered read**, on one file, at one time, returned the published bytes" | "it establishes that one **composite read** — unbuffered body, buffered tail — on one file, at one time, returned the published bytes" |
| §7 claim ceiling | "that one local observation — the **19:16:03 unbuffered read** — equals it" | "that one local observation — the **19:16:03 composite read (unbuffered body, buffered tail)** — equals it" |

**Why it matters.** "Unbuffered read" is the phrase that would let a later reader conclude the unbuffered path was validated in isolation. It was not, and no controlled experiment has yet separated the paths.

---

## 3. E-2 — "cannot exist" is stronger than the research supports

Two locations state an impossibility. The research established an absence in the repositories examined on one date; it did not establish that a publisher-originated GGUF reference cannot exist anywhere or at any future time.

| Location | Current text | Correction for a future revision |
|---|---|---|
| §2.2, "Implication" row | "no publisher-originated GGUF hash **can exist**" | "**No publisher-originated GGUF reference was found as of 2026-11-14**" |
| §6.2 | "`Qwen/Qwen3.5-9B` ships safetensors only, so **no publisher GGUF hash can exist**" | "`Qwen/Qwen3.5-9B` ships safetensors only; **no publisher-originated GGUF reference was found as of 2026-11-14**" |

A publisher could ship a GGUF later, or publish a checksum elsewhere; neither is excluded by what was checked. The classification of both hashes as **third-party reference** is unaffected and remains correct for the retrieval date.

---

## 4. Position at freeze

| Question | Finding |
|---|---|
| Is there an independent reference for S1-B? | **Yes** |
| Does a local observation match it? | **Yes — composite read** |
| Are the other six observations consistent with it? | **No** |
| Is the buffered path proven defective? | **No** |
| Is the disk proven to contain canonical bytes? | **No** |
| Is root cause identified? | **No** |
| Have any investigation gates been opened? | **No** |

**Phrase preserved for the investigation summary:** *anchored, not explained.*

---

## 5. Disposition

External-hash research **stops here**. It achieved its objective: a repository-verified reference for S1-B, no authoritative reference for S1-A or S1-C, and no root cause.

The next meaningful question is no longer "which hash should we believe?" but **"where did the non-matching byte streams originate?"** That requires a separately authorised controlled experiment, not further research and not another round of reads.

**Reopening rule:** these errata are applied only if `HOWL-EXT-HASH-02` is opened for a substantive reason. A revision made solely to fix wording would spend the freeze for no gain, and the frozen text plus this record is already unambiguous to a careful reader.

**Claim ceiling:** this record establishes two wording defects and their corrections, and the position at freeze. It establishes nothing about the file, the read paths, or the cause of the disagreement, and it authorises no access of any kind.
