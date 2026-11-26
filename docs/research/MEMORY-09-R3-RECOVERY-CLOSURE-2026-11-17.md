# MEMORY-09 — r3 Recovery Investigation: Closing Record

**Status:** closing record — investigation outcome only. Creates no authority, authorizes no work,
and decides no architectural question.
**Date:** 2026-11-17 (UTC)
**Scope:** bounded read-only investigation into whether revision **r3** of
`MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` could be recovered from any source inside this
workspace, its Git object store, or (if one had existed) an external backup.
**Outcome:** **UNSUCCESSFUL within the documented scope** (§3). No reconstruction attempted.

> **This record does not establish that an r3 artifact never existed.**
> It establishes only that no copy was found in the locations actually searched, which are listed
> exactly in §3. `Not found` and `never existed` are different statements; §6 keeps them separate
> and declines to convert the first into the second.

---

## 0. Provenance of this record

| Field | Value |
|---|---|
| **Agent role** | Buffy (Codebuff coding agent) |
| **Model label** | **DeepSeek V4.1 Flash Fast** — **artifact-reported, not independently verifiable.** This agent cannot read its own serving model from this environment. If the served model differs, the label is wrong and this record inherits that error |
| **Record date** | 2026-11-17 (UTC) |
| **Location** | `C:\Users\StoneWolfPC\HOWL_Chat` (win32) |
| **Record type** | Investigation closing record — documentation only |
| **Authorization** | Owner, 2026-11-17, two transactions: (1) creation of this record; (2) retarget of the ARCH-03 review from r3 to r32. Requested after an analytical recommendation relayed from a **third-party assistant transcript**, which is a source this agent cannot inspect (§4) |
| **Authority created** | **None.** No architecture approved, no suite frozen, no work authorized |

---

## 1. What was investigated

| Object | Why |
|---|---|
| **r3** of the stage contract, SHA-256 prefix `f8bbe6…`, reported by the owner as ~25,817 bytes | The recovery target |
| r1 (`2a550d8b…`) and r2 (`ade6281f…`) | Claimed to survive in a third-party transcript; bear on r3's lineage and on whether the hash prefixes are real |
| `E-SCOPE-CONTRACT v1.2` (S16–S29) | Named as an r3 addition; would be an independent witness to r3's content |
| Any `ARCH-03` architecture record | Would support the hypothesis that r3 was a separate artifact that was never assembled |

---

## 2. Result

| Object | Status | Detail |
|---|---|---|
| **r3** | **NOT RECOVERED** | No workspace file, no Git object, in any revision, matches the recorded identity (§4) |
| **r32** | **Present, intact** | 91,311 bytes, SHA-256 `e88ef2d1697da1d904fd8791af0ccf1023e3bd577c84b26774a72b70b4ef598d` — matches the value recorded in `MEMORY-09-PREREGISTRATION-FREEZE-2026-11-15.md` §3 as "design contract r32" |
| **r1** | **Corroborated, not byte-verified** | Third-party transcript reports 11,965 bytes and `2a550d8b…`; r32's changelog independently cites `r1 (2a550d8b…)`. Cross-source agreement, neither side re-hashable here (§4) |
| **r2** | **Corroborated, not byte-verified** | Third-party transcript reports 22,840 bytes and `ade6281f…`; r32's changelog cites `r2 (ade6281f…)` |
| **`E-SCOPE-CONTRACT v1.2`** | **No standalone artifact** | The term appears in five files as a *reference*; no v1.2 document exists. Its S16–S29 material is present **inside r32**, where all fourteen probes exist |
| **`ARCH-03` record** | **Does not exist** | The corpus ARCH series contains ARCH-01, ARCH-02, ARCH-05, ARCH-06 (2026-11-13, 4–9 KB each). **ARCH-03 and ARCH-04 are absent.** "ARCH-03" is otherwise used as the work-order identifier for the r3 review |

---

## 3. Search boundary — exact

Performed, read-only, entirely inside `C:\Users\StoneWolfPC\HOWL_Chat`:

**Workspace content sweep**
- 9,419 files hashed (SHA-256), skipping: symlinks, files > 8 MB, `.git`, `node_modules`,
  `__pycache__`, virtualenv caches. **Files whose digest begins `f8bbe6`: 0.**
- Files whose size is exactly 25,817 bytes: **0.** (27 files fall in the 25,000–26,700 B band; none
  is the target and none carries the prefix.)
- Files matching `*ARCH-STAGE-CONTRACT*`: exactly **1** — the r32 document itself. No alternate
  filename, no `.bak`/`.orig`/`~`/`.r2`/`.r3`/`.save` copy of it anywhere.
- Every file under `docs/` with `ARCH` in the name: enumerated with size and digest. No ARCH-03
  record; no file at 25,817 B; none matching the prefix.

**Git object store and history**
- Refs: 2 local branches (`master` @ `4526bbf`, `freebuff/chat-…` @ `6499865`), 3 remote-tracking
  refs (`origin/HEAD`→`6499865`, `origin/main`→`6499865`, `origin/master`→`4526bbf`), 3 tags, empty
  stash. **No fetch performed** — remote state read only from cached remote-tracking refs.
- 50 commits reachable from all refs; `git log --all --full-history -- <path>` returns nothing;
  `git ls-files` shows the document is **untracked** and `git check-ignore` shows it is **not
  ignored**.
- All **4,628** objects parsed (3,503 blobs, 1,069 trees, 53 commits, 3 tags):
  - containing the path string `ARCH-STAGE-CONTRACT`: **0**
  - blobs containing the document's title line (i.e. any revision of it): **0**
  - blobs whose digest begins `f8bbe6`: **0**
  - blobs of exactly 25,817 bytes: **0** (35 in the 25,000–26,700 B band)
- 1,103 unreachable/dangling objects (`git fsck --unreachable --dangling`, read-only): none is this
  document. **This is the decisive negative**: a file that had ever been staged and discarded would
  survive as an unreachable blob. None does, so no revision of this document has ever entered the
  object store.

**Not searched, and why**

| Location | Reason |
|---|---|
| Anything outside the workspace | Not authorized; no path supplied |
| External/removable backup | The owner stated no external backup exists yet |
| Network / remote refs beyond cache | No fetch authorized |
| Chat archives | This agent has no enumeration access to conversation history; a single attempt to read a thread is gated on the thread being explicitly referenced. The r1/r2/25,817 figures originated in a **third-party assistant transcript**, which is not inspectable here (§4) |

---

## 4. Provenance of the chat-derived figures — third-party transcript

The r1 and r2 sizes, the ~25,817-byte figure for r3, the 06:38–06:44 UTC timestamps, and the claim
that Q2–Q6 were accepted all originate in **another assistant's transcript**, relayed into this
project by the owner. They were **not** recovered from this project's own history by this agent.

Consequences that are binding on any later citation:

| Statement | Status |
|---|---|
| r32's changelog cites prefixes `2a550d8b`, `ade6281f`, `f8bbe6aa` | **Observed** — read from the file in this workspace |
| A third-party transcript reports matching prefixes for r1 and r2 | **Reported, not verifiable here** |
| The two sources agree on r1 and r2 | **Observed agreement between two sources**, one of which cannot be re-examined. This is stronger than either alone and weaker than byte verification |
| Q2–Q6 were accepted | **Reported only.** r32 §14 still lists Q2–Q6 as *open questions*, and findings F10 records that Q2 and Q4 are simultaneously pre-answered by adopted rules. No surviving artifact records acceptance |
| r3's size was ~25,817 bytes | **Reported only.** No artifact of that size exists (§3) |

---

## 5. Hash provenance — UNVERIFIED

`f8bbe6aa…` occurs in this workspace **twice**, both inside r32's own changelog, in the
self-referential line `r3 (f8bbe6aa…) introduced the edge model`. No artifact produces it.

Two structural facts weaken it further, both already recorded in the ARCH-03 review:
- **F17** — the changelog is a duplicated, self-divergent splice: thirteen revisions appear twice
  with differing wording. It is the only place r3 is attested.
- **F26** — revisions are cited by 8-hex prefixes only, while §13.1 promises "immutable references
  to prior revisions and their SHA-256".

**Conclusion:** the prefix is an **assertion by the changelog**, not a digest reproducible from any
surviving artifact. Whether it was ever computed from a complete r3 document **cannot be
determined**, and is not here asserted either way.

---

## 6. "Not found" is not "never existed"

Two explanations remain open. Both are inferences; neither is established.

**(a) r3 was a revision of the same growing file, and its bytes were never preserved.** Supported by
the monotonic size curve — r1 11,965 B → r2 22,840 B → r3 ~25,817 B → "~51 KB across thirteen
revisions" (r32 §13.1, as recorded in F23) → r32 91,311 B — and by the complete absence of any
second artifact. A single document growing from 12 KB to 91 KB across 32 revisions fits these
numbers without strain.

**(b) r3 was a separate architecture record that was never materialised.** Consistent with the
ARCH-03/ARCH-04 numbering gap, but not supported by size: the corpus ARCH records are 4–9 KB and
dated two days earlier, and no 25,817-byte artifact exists at any name.

Explanation (a) is simpler and better supported. It is recorded here as **likely**, not
**established**. What is established is narrower and sufficient for the owner's purpose: **no r3
artifact exists in the locations searched, and none is retrievable by this route.**

---

## 7. What survives of r3's content

| Element | Status in r32 |
|---|---|
| S16–S29 probe block | **Present** — all fourteen probes exist in r32 |
| Edge model (`MemoryRecord` / `KnowledgeEdge` / `AcquisitionEvent`) | **Present, but amended.** r32 §6.3 records that `revoked_turn` was **removed in r6**. A genuine r3 `KnowledgeEdge` almost certainly carried it. r32's definitions are therefore **not** r3's |
| Q2–Q6 | **Present as open questions** in §14, not as accepted decisions (§4) |
| `E-SCOPE-CONTRACT v1.2` as a document | **Absent.** Its material was folded into the stage contract instead |

Net: r3's content is **partially represented in r32 in amended form**. Revision-specific content is
**unestablished**, and r32 **does not substitute** for r3.

---

## 8. What this record does not do

No reconstruction of r3. No probe, test, experiment or model call. No implementation inspected. No
registry edit. **No pre-existing file modified** (the retarget authorized alongside this record was
applied to the ARCH-03 review record as a separate, separately recorded transaction). No
architectural decision taken or implied. No finding in the ARCH-03 review is validated or
invalidated by this record; F2, F3, F5, F6, F7 and F15 remain open.

---

## 9. Resumption condition

This investigation is **closed**. It reopens only on a new, concrete source: an external backup
copy, a forensic image, or a retrievable transcript containing the complete r3 text. Absent such a
source, further searching cannot converge and is not authorized.

---

## 10. Registry consequence — declared, not repaired

This is a new file under `docs/`, which `docs/DOCUMENT-REGISTRY.md` states covers that tree
exhaustively. **No registry row was added**, because no registry edit was authorized.

The registry's completeness claim is now understated by **at least three** files:

1. `docs/research/MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` (finding F18)
2. `docs/research/MEMORY-09-ARCH-03-R3-INDEPENDENT-REVIEW-RECORD-2026-11-17.md`
3. this record

The review record's §6.1 declares only the first two, because it was written before this one
existed. That understatement is recorded here rather than repaired there, to keep the retarget
amendment strictly within its authorized scope.

---

## 11. Provenance and boundaries

| Field | Fact |
|---|---|
| Files created | **Exactly one** — this record |
| Files modified | **None.** The review-record retarget is a separate authorized transaction |
| Commands used | `sha256sum`, `grep`, `find`, `ls`, `wc`, `git log`, `git rev-list`, `git fsck`, `git cat-file`, `git ls-files`, `git check-ignore`, `git for-each-ref`, `stat` — all read-only; no `fetch`, `checkout`, `reset`, `restore` or `stash` |
| Experiments / probes run | **None** |
| Model calls made | **None** |
| Runtime inspected | **None** |
| Evidence regenerated | **None** |
