# Owner Decision Record — Decisions D1–D7

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Identifier:** `CORPUS-69` · **Version:** 1 · **Date:** 2026-11-12.
**Record type:** consolidated owner-decision instrument. **Not a proposal. Not an authorization to implement.**
**Assignment model identity (owner-supplied):** DeepSeek. Framework persona label: Buffy — an execution fact, not model identity.
**Filed under:** the owner work order "`CORPUS-68` Remediation and Reconciliation", §3.
**Status:** all seven decisions were **approved by the owner on 2026-11-12** and have **not been applied**. This document records them; it changes no corpus item, standard, tag, or checker.

---

## 0. Provenance — read this before relying on anything below

The seven decisions were **approved by the owner in conversation on 2026-11-12**. Three artefacts now carry them, and they are not interchangeable:

| Layer | Artifact | What it is |
|---|---|---|
| **Owner approval** | the conversation itself | The act. Not filed, not reproducible from the workspace. |
| **Agent transcription (1)** | `CORPUS-62` §3 (`35ab7b9f…`, row 397) | Detailed record with rationale, scope, consequences per decision. |
| **Agent transcription (2)** | this file | Consolidated instrument, aligned to the work order's own D1–D7 numbering. |

**What is established:** the decisions were stated, approved, and transcribed on 2026-11-12.
**What is not established, and is not invented here:** no signature, no owner-signed file, no filing timestamp proving ownership, no cryptographic attestation, and no verification mechanism of any kind. A reader cannot confirm from the workspace that the owner — rather than an agent reporting the owner — approved these. That limitation is permanent for this record and is **not** a defect to be repaired by adding a signature line.

**Consequence, stated plainly:** every downstream proposal (`CORPUS-63`–`67`, and the revisions in `CORPUS-71`) inherits this provenance limit. If a decision below misstates owner intent, the error propagates and nothing in the workspace will catch it.

**These decisions are not reopened.** `CORPUS-68`'s S-1 finding concerns *verifiability*, not approval. Where `CORPUS-68` raised a genuine ambiguity in **how a decision should be executed**, it is recorded separately in `CORPUS-70` §3 (S-1) and `CORPUS-72`, with the decision text left intact.

---

## 1. Decision D1 — `CORPUS-58` registry identity

**Approved.** The root-level `CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md` is a **real historical artifact and an assignment brief**, not the substantive review. Preserve it **byte-for-byte**. Correct its **registry identity and path** while retaining the original registration history. Do not represent the brief as the completed review.

**Limits.** One registry row (393) plus the scope/row-count lines that follow from it. No rename, no second registration of the same file, no reconstruction of the review, no edits to `CORPUS-59`/`60`/`61`.

**Status.** **Applied** by `CORPUS-63` §4. Row 393 now links `../CORPUS-58-…` and resolves; hash `292befdd…`, 4,311 bytes, 80 lines unchanged; original role and wrong path preserved in-row and quoted in `CORPUS-63` §3.

---

## 2. Decision D2 — the missing substantive review

**Approved.** Do **not** reconstruct the absent original review and present it as authentic. Conversation-supplied material may be cited as conversation-supplied. Any future reconstruction must carry a **distinct identifier**, an explicit "RECONSTRUCTION — NOT THE ORIGINAL" label, and disclose that the original file, hash, and provenance are **unavailable**; no byte-identity claim; no convergence credit.

**Limits.** Provenance disposition only.

**Status.** **Disposition recorded; no reconstruction created.** Provenance schema with fields left UNVERIFIABLE at `CORPUS-63` §5.1; the two prohibitions recorded there so they bind in advance.

---

## 3. Decision D3 — `F3-01` intended construct

**Approved.** `F3-01` is intended to test **memory-dependent recognition of the persona's own historical experience**, not generic conversational appropriateness.

**Limits.** `F3-01` only. States an intent; establishes no diagnosticity, adopts no pass condition, and does not make the item Type P.

**Status.** Construct carried into `CORPUS-64` §4; revised in `CORPUS-71` §2.

---

## 4. Decision D4 — `F3-01` adjudication

**Approved.** `CORPUS-07` §2's existing disjunction is **insufficient to establish diagnosticity**. Do not adopt it as written. Develop a historically attributable recognition requirement and test memory-absent alternatives. **Preserve the prior adjudication visibly in the historical record.**

**Limits.** Preparation of a proposal only. `CORPUS-07` §2 remains an unapplied adjudication and is not edited.

**Status.** `CORPUS-64` v1 prepared; revised to v1.1 in `CORPUS-71` §2. `CORPUS-07` unchanged (`72f6e2a2…`).

---

## 5. Decision D5 — `R-2` discrepancy

**Approved.** Treat the conflict among `CORPUS-01` §3.1's explicit item list, its stated count, the item tags, and `F1-03`'s omission as a **specification inconsistency**. Do **not** silently change a tag, remove `F1-03`, or infer historical intent. Prepare a **dated amendment proposal**.

**Limits.** `CORPUS-01` §3.1, one row, plus a version notice. Not `MEMORY-09`. The tag-side reading stays recorded and unproposed.

**Status.** `A-R2-01` prepared in `CORPUS-65`; revised in `CORPUS-71` §3. Not applied; `CORPUS-01` unchanged (`be1f11b4…`).

**Execution ambiguity, recorded not resolved (see `CORPUS-70` §6.1).** Adopting `A-R2-01` *is itself* the determination that the row was stale rather than the tag wrong. D5 presumes the specification-side reading; the amendment therefore cannot be presented as neutral. `CORPUS-65` §1 says this; `CORPUS-72` OD-5 makes it a decision the owner must make explicitly.

---

## 6. Decision D6 — `F6-01`

**Approved.** The current scenario **does not establish the relevant cake contents** and therefore cannot support a valid positive test in its present form. Prepare a **narrowly scoped scenario proposal** that establishes the necessary restriction while **preserving the adjudicated observable accommodation and the fact that silence fails**. Add no unsupported assumptions.

**Limits.** Scenario/evidence preparation only. The adjudicated threshold is preserved, not re-adjudicated. The `R-2` tag is not changed.

**Status.** `CORPUS-66` v1 prepared; revised in `CORPUS-71` §4. `CORPUS-02` and `CORPUS-07` unchanged.

**Execution ambiguity, recorded (see `CORPUS-70` §5.1).** D6 says "preserve the adjudicated … accommodation". Specifying *what the warning must be about* ("the cake contains almonds") adds content the adjudicated text ("warn", no object) does not contain. `CORPUS-66` §6 asserted the threshold was "preserved in substance" while adding that object. That is a **tightening**, and D6 alone does not authorize it. It is raised as a separate decision (OD-3), not folded into D6.

---

## 7. Decision D7 — registry validator

**Approved.** A future integrity validator should detect **missing and mislinked targets** and check **identity, path, hash, bytes, lines, duplicates, and root-scope consistency**. It must be **read-only** and must **never auto-repair**. **No checker implementation is authorized during `G1`.**

**Limits.** A specification. Nothing implemented; `check_corpus.py` untouched.

**Status.** `CORPUS-67` v1 prepared; revised in `CORPUS-71` §5. Checker unchanged (`be3e317c…`).

**Execution ambiguity, recorded (see `CORPUS-70` §7.2–7.3).** D7 lists capability classes but not severity policy. Whether an `UNREGISTERED_IN_SCOPE` finding fails the run is a policy choice `CORPUS-67` leaves indeterminate and that the approved decision does not settle.

---

## 8. What this record does not do

- It does not apply, adopt, implement, or authorize any proposal.
- It does not convert any reviewer finding into an owner decision, or any owner decision into a reviewer conclusion.
- It does not reopen D1–D7, and it does not treat the four recorded execution ambiguities as absent approval.
- It does not touch `CORPUS-01`, `02`, `03`, `04`, `05`, `07`, `MEMORY-09`, `check_corpus.py`, or any corpus item, tag, split, or acceptance condition.
- It runs no model, performs no inference, measures nothing.

## 9. Status

| Field | State |
|---|---|
| Decisions recorded | **7 of 7** |
| Provenance | Conversation, 2026-11-12 — owner approval; transcription by DeepSeek |
| Independent verification | **None available. Not invented.** |
| Applied | **Only D1** (`CORPUS-63`) |
| Governing standards | Unchanged |
| `G1` | Open · corpus frozen · `AM-01` unissued · `A-2a-R1` NOT IN FORCE · `D1` unadopted · convergence **0 of 2** |
