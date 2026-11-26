# MEMORY-09 — Documentation Reconstruction Record (REV-2)

**Identity:** `MEMORY-09-DOCUMENTATION-RECONSTRUCTION-RECORD-2026-11-15-REV-2` · **Version:** REV-2 · **Set:** artifact 5 of 5, replacement revision · **Created:** 2026-11-16

**Supersedes:** `docs/research/MEMORY-09-DOCUMENTATION-RECONSTRUCTION-RECORD-2026-11-15.md`, SHA-256 `636ea2ceeb04fcc9a8b74ce6fb7fb42fb2ce8195f396901c2f374b2c4eef3ce7`. The original is **preserved and not modified**; it remains the reviewed subject set.

**Finding answered:** `R-05` — §9 stated "There is no implementation", which is broader than the evidence supports given that the repository contains evaluation harness code.

**Corrections applied:** `R-05` — the §9 implementation bullet states that no implementation conformance to the durability contract has been established, that this reconstruction authorised no implementation work, and that the existence of evaluation harness code does not establish implementation of, or conformance to, the contract. `M-1b` — the §9 execution bullet is narrowed to this documentation set, excluding earlier evaluation activity recorded elsewhere in the repository. `M-2` — this revision's own creation date is recorded, without changing any date belonging to the original set. `L-4` — the referent of the §5 verification statement is clarified as the creation of this record, not of the frozen files.

**Boundedness:** the corrections are limited to `R-05`, `M-1b`, `M-2` and `L-4`. `M-2` records this revision's own creation date and leaves every date belonging to the original set untouched; no preflight observation, discrepancy or reconciliation task is changed. The work order's `R-06` constraint is applied in §5 and §7: identity is stated only as *the byte sequence read at verification is the byte sequence present at creation*, and no file is described as unchanged, unmodified or uncorrupted. Documentation correction only — **no contract approval, implementation, execution or G4 authority is created.**

**Status: DRAFT — UNREVIEWED (REV-2). Reconstructed from a review handoff, not from the source transcript. Textual review only; no execution, no implementation, no G4 work performed or authorised. Marked items were not supplied and must not be invented.**

**Implementation status: NOT STARTED / NOT AUTHORIZED**

---

## 1. What this record covers

Provenance, dates, preflight observations and limits for the five-artifact reconstruction set created on 2026-11-15. The other four artifacts are:

1. `MEMORY-09-DURABILITY-CONTRACT-2026-11-15.md`
2. `MEMORY-09-GOVERNING-PRINCIPLES-2026-11-15.md`
3. `MEMORY-09-REVIEW-LEDGER-2026-11-15.md`
4. `MEMORY-09-OPEN-DECISIONS-2026-11-15.md`

## 2. Provenance

| Fact | Statement |
|---|---|
| Origin | The contract was produced in a review conversation, rev.1 → rev.23.8, not from repository source files |
| Route | Reconstructed from a documentation brief (handoff rev.3.1) supplied in that conversation |
| Repository presence of the brief | None. The brief is not a repository artifact and had no hash when this set was created |
| Source of record | **The transcript remains the historical source of record.** This set is not a substitute for it |
| Lossiness | The brief is lossy at clause-prose level. No clause text was supplied, so no clause text appears in the set |
| Marker meaning | Every pending marker is a reconciliation task for the owner, not a statement that the item is unimportant or absent from the design |

The transcript and this set are **not interchangeable** until the owner reconciles them. Until then, citing this set as the contract would cite a reconstruction of a review of a draft.

## 3. Dates

| Item | Value |
|---|---|
| Proposed date (in the brief) | 2026-11-16 |
| Actual local date at creation | **2026-11-15** — preflight clock read `2026-11-15T23:35:43-0600` |
| Date used in all five filenames and status blocks | **2026-11-15** |
| Reason | The brief requires the actual local calendar date of completion; the proposed date was not the actual date, so it was not used and no reference to it survives in the set other than this table |

**This REV-2 revision.** Created 2026-11-16, local time `00:27:16-0600`. The filename retains the original's date because a revision names what it supersedes; its own creation date is recorded here, not in the filename.

No UTC timestamp is asserted anywhere in the set. The clock value above is local time, labelled with its offset.

## 4. Preflight observations

Read-only inspection performed before any file was written.

| Check | Result |
|---|---|
| Four frozen preregistration hashes vs the freeze record | **All four match exactly** (see §5) |
| Frozen preregistration byte counts | 9,235 · 9,800 · 8,717 · 6,803 |
| `MEMORY-09-CHECKPOINT-2026-11-15.md` present | Yes |
| `MEMORY-09-CONTRACT-NORMATIVE-2026-11-15.md` present | Yes |
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` present | Yes |
| `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` present | Yes |
| Index position, top level of `docs/research/` | **No index exists.** No file matching index, readme or manifest naming |
| Index position, `docs/research/corpus/` | One index: `CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md`. It is a dated artifact and therefore protected |
| Index edited or created | **No.** No index was edited and none was created |
| Target filenames free before writing | Yes — none of the five existed |
| Git tracking of the protected MEMORY-09 documents | The MEMORY-09 research documents, including the frozen set, are **untracked** in git at the time of this record. Their identity rests on the hashes, not on version control |

## 5. Frozen set, as re-verified

| Document | SHA-256 |
|---|---|
| `MEMORY-09-PREREGISTRATION-2026-11-15.md` | `8c788a96ea6a866d4fd06e13496ddbfa46b649da5a5e0c8ef9bace960844ff12` |
| `MEMORY-09-PREREGISTRATION-AMENDMENT-1-2026-11-15.md` | `975033401ae1beaa9caaa5d5ae55aab0ca3ddb4f5159eb41e9feb8e01a4c9b1a` |
| `MEMORY-09-PREREGISTRATION-AMENDMENT-2-2026-11-15.md` | `c3ffd7ed735ed4cf68f8b455555af2a1d10eb38273b07a756f269076cd917f8f` |
| `MEMORY-09-PREREGISTRATION-AMENDMENT-3-2026-11-15.md` | `c5af036265b03662b26feaf40d358ad08c3063fa0055c72d2fb0c2f41c236281` |

These are the values recorded in the freeze record and re-observed at the creation of this record. The byte sequence read at verification is the byte sequence present at creation of this record — not at the creation of the frozen files. That establishes equality only between the freeze record's value and the value observed here; it does not establish that the files were never changed and restored in between.

## 6. Discrepancies found against the brief

Two, both reported rather than reconciled.

| # | Discrepancy | Handling |
|---|---|---|
| 1 | The brief names the protected file `docs/research/OWNER-DECISION-RECORD-OD-M09-2026-11-13.md`. No file exists at that path. The file exists as `docs/research/MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md`, and its addendum as `docs/research/MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md` | The brief omits the `MEMORY-09-` prefix. The repository is authoritative. No file was renamed, created or modified to make the brief correct, and the two real paths are cited in §7 and in artifact 4 |
| 2 | The brief proposes the date 2026-11-16; the actual local date is 2026-11-15 | The actual date is used throughout, per §3 |

Neither discrepancy affects the content of the reconstruction. Both are recorded because a silent correction would have made the set disagree with the brief that authorised it.

## 7. Protected material

No file in this table was written to by this work. The second column states the evidence available for each item; a missing recorded hash is not evidence that a file was unchanged.

| Item | Evidence |
|---|---|
| Frozen preregistration set (4 files) | Not written by this work; byte sequences read at verification matched the values in §5 |
| `MEMORY-09-CHECKPOINT-2026-11-15.md` | Not written by this work; no hash recorded in this document |
| `MEMORY-09-CONTRACT-NORMATIVE-2026-11-15.md` | Not written by this work; no hash recorded in this document |
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` | Not written by this work; no hash recorded in this document |
| `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | Not written by this work; no hash recorded in this document |
| `MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md` | Not written by this work; no hash recorded in this document |
| `MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md` | Not written by this work; no hash recorded in this document |
| Every other dated `.md` under `docs/research/` and `docs/research/corpus/` | Not written by this work; not enumerated or hashed here |
| `eval/**` including receipts, sidecars and fixtures | Not written by this work; no instrument executed |
| Files written | **Exactly five**, all named in §1 |

Pre-existing working-tree modifications that predate this work (`docs/research/OWNER-ENVIRONMENT-HISTORY-2026-11-14.md`, `eval/harness/validate.py`, `eval/harness/__pycache__/validate.cpython-312.pyc`) were left as found.

## 8. What this set establishes

- A structured, dated, reviewable repository representation of the durability contract review.
- An explicit boundary between what was supplied, what was derived, and what is missing.
- The preflight state of the repository at creation, including the index position and the frozen hashes.
- A record of the two brief discrepancies and how they were handled.

## 9. What this set does not establish

- Conformance of any implementation to the contract. No implementation conformance to the durability contract has been established, and this reconstruction authorised no implementation work. The repository contains evaluation harness code; its existence does not establish implementation of, or conformance to, the durability contract.
- Execution of the durability contract. Nothing was run and no model call was made in producing this set; earlier evaluation activity recorded elsewhere in the repository is outside this claim.
- Implementation authorisation. None was granted and none is implied.
- Approval, adoption or freeze of the contract. Revision 23.8 remains `DRAFT-UNREVIEWED`.
- Historical accuracy. Clause prose, finding texts, identifier history and all four registries are missing, and their absence is marked rather than filled.
- That this set may be cited as the contract. It may be cited as a reconstruction of a review of a draft contract, with the transcript as the source of record.

## 10. Reconciliation tasks handed to the owner

| # | Task | Covers |
|---|---|---|
| 1 | Supply clause prose for the eleven families | Artifact 1 §4 |
| 2 | Establish the identifier mapping, including the meaning of prime and double-prime forms | Artifact 1 §5 |
| 3 | Recover the four registry contents, including the 16 artifact identity rows | Artifact 1 §6 |
| 4 | Recover principle statements and the P1, P2–P5 and P7 instance inventories | Artifact 2 §2, §4 |
| 5 | Recover the 18 counted finding texts and the three `L-02a`–`L-02c` texts | Artifact 3 §3 |
| 6 | Record subjects and closure conditions for the 19 open owner decisions | Artifact 4 §3 |

Until these are done, this set is a skeleton with an honest label, and it should be read as exactly that.
