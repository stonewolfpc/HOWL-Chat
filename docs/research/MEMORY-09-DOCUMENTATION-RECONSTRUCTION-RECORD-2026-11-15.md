# MEMORY-09 — Documentation Reconstruction Record

**Identity:** `MEMORY-09-DOCUMENTATION-RECONSTRUCTION-RECORD-2026-11-15` · **Set:** artifact 5 of 5 created 2026-11-15 under owner authorisation

**Status: DRAFT — UNREVIEWED. Reconstructed from a review handoff, not from the source transcript. Textual review only; no execution, no implementation, no G4 work performed or authorised. Marked items were not supplied and must not be invented.**

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

These are the values recorded in the freeze record and re-observed here. Re-observation does not establish that the files were never changed and restored at some point between the freeze and this reading; it establishes that the byte sequence read now matches the recorded value.

## 6. Discrepancies found against the brief

Two, both reported rather than reconciled.

| # | Discrepancy | Handling |
|---|---|---|
| 1 | The brief names the protected file `docs/research/OWNER-DECISION-RECORD-OD-M09-2026-11-13.md`. No file exists at that path. The file exists as `docs/research/MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md`, and its addendum as `docs/research/MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md` | The brief omits the `MEMORY-09-` prefix. The repository is authoritative. No file was renamed, created or modified to make the brief correct, and the two real paths are cited in §7 and in artifact 4 |
| 2 | The brief proposes the date 2026-11-16; the actual local date is 2026-11-15 | The actual date is used throughout, per §3 |

Neither discrepancy affects the content of the reconstruction. Both are recorded because a silent correction would have made the set disagree with the brief that authorised it.

## 7. Protected material

Confirmed unmodified by this work.

| Item | State |
|---|---|
| Frozen preregistration set (4 files) | Unchanged; hashes match §5 |
| `MEMORY-09-CHECKPOINT-2026-11-15.md` | Unchanged |
| `MEMORY-09-CONTRACT-NORMATIVE-2026-11-15.md` | Unchanged |
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` | Unchanged |
| `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | Unchanged |
| `MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md` | Unchanged |
| `MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md` | Unchanged |
| Every other dated `.md` under `docs/research/` and `docs/research/corpus/` | Unchanged |
| `eval/**` including receipts, sidecars and fixtures | Untouched; no command wrote to it and no instrument was executed |
| Files written | **Exactly five**, all named in §1 |

Pre-existing working-tree modifications that predate this work (`docs/research/OWNER-ENVIRONMENT-HISTORY-2026-11-14.md`, `eval/harness/validate.py`, `eval/harness/__pycache__/validate.cpython-312.pyc`) were left as found.

## 8. What this set establishes

- A structured, dated, reviewable repository representation of the durability contract review.
- An explicit boundary between what was supplied, what was derived, and what is missing.
- The preflight state of the repository at creation, including the index position and the frozen hashes.
- A record of the two brief discrepancies and how they were handled.

## 9. What this set does not establish

- Conformance of any implementation to the contract. There is no implementation.
- Execution. Nothing was run. No model call was made.
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
