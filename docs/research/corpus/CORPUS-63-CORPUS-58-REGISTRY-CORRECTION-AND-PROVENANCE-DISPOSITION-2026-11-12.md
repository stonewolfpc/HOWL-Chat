# `CORPUS-58` Registry Correction and Provenance Disposition

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-63` · **Version:** 1 · **Date:** 2026-11-12.
**Authorizing decisions:** `CORPUS-62` Decisions 1 and 2 (owner-approved 2026-11-12, arising from `CORPUS-61`).
**Assigned authoring model identity (owner-supplied metadata):** DeepSeek. **Not** inferable from workspace evidence and **not** claimed on the basis of any self-description.
**Executing agent:** Codebuff coding agent under the framework persona label "Buffy". The persona label is a framework identity, not model identity.
**Execution:** fresh conversation.
**Status:** one **administrative registry correction applied** (row 393, authorized by Decision 1) plus a **provenance disposition recorded** (Decision 2). **NOT operative as to `F3-01`, `F6-01`, `R-2`, or the checker.** No corpus item, label, standard, schema, runtime, model, or checker was changed.

---

## 1. What this document is for

`CORPUS-61` §5 established three distinct evidentiary objects around `CORPUS-58`, and refused to substitute one for another:

| Object | Description | State at `CORPUS-61` |
|---|---|---|
| **O-1** | The repo-root file `CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md` | Present, 4,311 bytes, 80 lines, hash `292befdd…`; untracked; content = assignment brief (work order), **no findings** |
| **O-2** | The registry's linked corpus path `research/corpus/CORPUS-58-…md` | **Absent** — no file at that path |
| **O-3** | The substantive review text with findings, furnished in other conversations | **Not in this workspace and not furnished in this execution** |

Decision 1 disposes of O-1 and O-2. Decision 2 disposes of O-3.

This document records both dispositions, quotes the before/after registry row verbatim, and states exactly what remains unverifiable.

---

## 2. Verified facts about O-1 (the root artifact)

Verified this session against the working tree, not against any prior report's summary.

| Attribute | Value |
|---|---|
| Path | `CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md` (repository root) |
| Bytes | **4,311** |
| Lines | **80** |
| SHA-256 | `292befdd5a4529cef79aa7cb10b7e7ea5aa5fe99e15cdb21f5c1205d47cfad56` |
| Git state | untracked (`??`) |
| mtime | 2026-11-12 12:47 |
| Content role | **Assignment brief / work order** for an independent review |
| Contains findings? | **No** |
| Contains verdicts? | **No** |
| Contains source citations? | **No** |

**Content structure read in full:** recipient and review target; objective; seven required challenges; evidence requirements; prohibited actions; required deliverable naming this same filename; and the statement that a chat response alone is not the deliverable.

It is the instruction *for* the review. It is not the review. A byte-identical file cannot satisfy the brief's own "Required deliverable" condition, because the deliverable was to contain findings, and this file contains none.

**Byte-for-byte preservation:** this correction did not read-write, rewrite, re-encode, move, or rename the artifact. Its 4,311 bytes and hash are unchanged by this filing. Only the registry row that points at it was corrected.

---

## 3. The original registration (preserved, not erased)

Registry row 393 read, verbatim, before this correction:

```text
| 393 | [CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md](<research/corpus/CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md>) | Research | Adversarial Review | Review | `292befdd5a4529cef79aa7cb10b7e7ea5aa5fe99e15cdb21f5c1205d47cfad56` | 4311 | 80 |
```

What was wrong with it, precisely — and what was not:

| Claim in the row | Finding |
|---|---|
| Hash `292befdd…` | **Correct.** Matches O-1 exactly. |
| Bytes `4311` | **Correct.** Matches O-1 exactly. |
| Lines `80` | **Correct.** Matches O-1 exactly. |
| Linked path `<research/corpus/…>` | **Wrong.** No file exists at that path (O-2). |
| Role "Adversarial Review", placement "Review" | **Wrong.** The artifact is an assignment brief. Naming a file after the review it requests does not make it the review. |

So the defect is **linkage and role**, not metadata. The row's numbers identified a real artifact all along; `CORPUS-59`'s "dangling / unbacked metadata" characterisation was corrected by `CORPUS-61` §5.3 on exactly this ground.

---

## 4. The correction applied (Decision 1, required actions 3–6)

One row was rewritten. Its **number (393) is preserved**; no new row was created for the same artifact, because a second row pointing at the same file would introduce a duplicate target and make the registry ambiguous about which row controls. The corrected row is:

```text
| 393 | [CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md](<../CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md>) | Research | Assignment brief (work order) | CORRECTED 2026-11-12 (`CORPUS-62` D1, recorded in `CORPUS-63`) — this file is the **assignment brief** for the CORPUS-58 review, at the repository root, not the review. Original registration recorded it as Role "Adversarial Review" / placement "Review" at `research/corpus/…`, where no file exists. Hash/bytes/lines are unchanged and now resolve. The substantive CORPUS-58 review remains **unfiled**; see `CORPUS-63` §§5–7. | `292befdd5a4529cef79aa7cb10b7e7ea5aa5fe99e15cdb21f5c1205d47cfad56` | 4311 | 80 |
```

Correction rules observed:

| Requirement (Decision 1) | How it is satisfied |
|---|---|
| 1. Preserve the root brief byte-for-byte | No write to the artifact; hash re-verified after filing (§8). |
| 2. Preserve its original hash and metadata | The three numbers in the row are unchanged: `292befdd…`, 4311, 80. |
| 3. Auditable correction identifying the actual artifact, location, role | The row now resolves to the root path and names the role "Assignment brief (work order)". |
| 4. Preserve the fact that row 393 originally misidentified the artifact | The correction note quotes the original role and the original (wrong) path in the row itself, and §3 of this document quotes the row verbatim. The audit trail survives both a file change and a registry regeneration. |
| 5. Do not represent the brief as the completed review | Role changed to "Assignment brief"; the note states the review remains unfiled. |
| 6. Do not erase the original registration history | The row keeps number 393; nothing was deleted; the original claim is recorded as text in the row and in §3 here. |

**Scope line.** The registry's declared scope was "Every file under `docs/` plus the root requirements ledger". Because the corrected row now points outside `docs/`, the scope line records that second root-level exception explicitly rather than leaving the row contradicting the declared scope.

**What the correction deliberately does not do:**

- It does not make the brief the review.
- It does not imply the review was ever filed, completed, or received.
- It does not attach findings to the brief.
- It does not confer independent-reviewer status, convergence credit, or any evidentiary weight on the brief — a work order is an instruction, not evidence about the work.
- It does not alter rows 391, 392, 394–396 or any other row.

---

## 5. Decision 2 — no reconstruction was created, and why

Decision 2 permits a reconstruction **if a reconstruction is useful**, subject to five conditions (new identifier, explicit reconstruction label, source identified, unverifiable provenance stated, no byte-identity claim, not convergence evidence).

**No reconstruction artifact was created in this filing.** The reason is factual, not discretionary:

- The substantive `CORPUS-58` text is **not present in this workspace** (`CORPUS-61` §3.2) and was **not furnished in this execution**. No turn of this conversation contained it.
- A reconstruction must identify the conversation-supplied text as its source and reproduce its substance. Without the text, any reconstruction would be written from other reports' *descriptions* of it. That would make it a third-order paraphrase — `CORPUS-59`/`CORPUS-60` describing a review neither had on disk — and presenting it as a reconstruction of the review would be exactly the fabrication the governing rules forbid.
- Writing it would also require inventing the very metadata Decision 2 says must not be invented: original filename, byte count, line count, hash, registry provenance, filing date.

**Therefore the correct action under Decision 2 is to record the disposition and the blockers, not to manufacture an artifact.** The conversation-supplied text, where it exists, remains **conversation evidence**: it may be read, quoted as "reported in another conversation", and reasoned about with that qualification. It may not be cited as a filed artifact, and no finding may rest on it as though it had been verified at source.

### 5.1 Provenance record for a possible future reconstruction

This is the record a future authorized reconstruction would have to fill in. Fields marked **UNVERIFIABLE** must stay unverifiable rather than be completed by assumption.

| Field | Required content | Present state |
|---|---|---|
| Proposed identifier | A **new** identifier, never `CORPUS-58` and never a reuse of any issued ID | Not issued — no artifact exists |
| Explicit label | The words "RECONSTRUCTION — NOT THE ORIGINAL" in the title and status line | — |
| Source | "Conversation-supplied text, furnished in conversation(s) on [date(s)], not present in the workspace" | Date(s) **UNVERIFIABLE** here |
| Source location | Which conversation(s), by the owner's record | **Not furnished in this execution** |
| Original filename | Unknown | **UNVERIFIABLE** — must be recorded as unknown, not guessed |
| Original SHA-256 | Unknown | **UNVERIFIABLE** |
| Original byte count | Unknown | **UNVERIFIABLE** — the claimed 4,311 describes the brief (O-1), not the review |
| Original line count | Unknown | **UNVERIFIABLE** — the claimed 80 describes the brief (O-1) |
| Original registry provenance | Unknown | **UNVERIFIABLE** — row 393 never registered the review |
| Authorship / model identity | Owner-supplied at the time, if any | **Owner metadata only; not workspace-verifiable** |
| Independence | Must be recorded as *not established* | **NOT ESTABLISHED** |
| Convergence eligibility | Must be recorded as *zero* | **0 of 2 stands** |

Two prohibitions attach to any future reconstruction, and are recorded here so they bind in advance:

1. **No byte-identity claim.** It may not be described as identical, equivalent-in-bytes, or a faithful copy of the original. At most: "a reconstruction from conversation-supplied text, whose original is unavailable for comparison."
2. **No convergence credit.** `CORPUS-05` §7 requires two independent reviewers, and a reconstruction of an unavailable review is not an independent review. A reconstruction may not be counted toward the 0-of-2 convergence figure.

---

## 6. States after this filing

| Object | State | Who may rely on it |
|---|---|---|
| **O-1** root brief | Filed, unmodified, hash-verified, registered at row 393 with corrected link and role | Anyone needing to know what the reviewer was *asked* to do |
| **O-2** corpus-tree path | **Still absent.** No file was created there. The row no longer points at it. | Nobody — it does not exist |
| **O-3** substantive review | **Still unfiled.** Not reconstructed, not registered, not citable at source | Nobody, except as "conversation-reported", with that qualification stated every time |

**Consequence for earlier conclusions.** Any claim of the form "`CORPUS-58` found X" remains **UNVERIFIABLE at source** in this workspace. `CORPUS-60`'s citation checks and `CORPUS-61`'s 60-4/60-12 items keep that status. Nothing in this correction upgrades them.

---

## 7. Registry state after the correction

| Check | Result |
|---|---|
| Row count after adding this filing's six new artifacts | 402 rows (see `CORPUS-62` §11) |
| Row 393 target resolves | **Yes** — root file, 4,311 bytes |
| Row 393 hash matches | **Yes** — `292befdd…` |
| Row 393 role accurate | **Yes** — assignment brief |
| Duplicate targets introduced | **None** |
| Original registration history | **Preserved** in the row's correction note and §3 above |

**Known gaps that this correction did not touch, and was not authorized to touch:** two `docs/` files (`CORPUS-26`, `CORPUS-32`) have no registry row. They are recorded as a pre-existing coverage gap in `CORPUS-62` §9.4 and are not repaired here.

---

## 8. Verification performed

| Item | Method | Result |
|---|---|---|
| O-1 unchanged | `sha256sum` and byte count before/after this filing | `292befdd…`, 4,311 bytes, 80 lines — unchanged |
| Registry row 393 correct | Row parsed, path resolved from `docs/`, hash/bytes/lines recomputed | All three match; target exists |
| No artifact created at O-2 | Directory listing of `docs/research/corpus/` | No `CORPUS-58*` file there |
| No other row changed | Full-registry target resolution sweep | Only row 393 changed by this work |

---

## 9. Unresolved questions (owner decisions)

1. Whether the substantive `CORPUS-58` review should be furnished and filed at all, and if so under which new identifier and with which of the §5.1 fields honestly marked unknown.
2. Whether the untracked root brief should be committed to version control. This filing left its git state untouched.
3. Whether the brief's filename — which asserts "INDEPENDENT ADVERSARIAL REVIEW" for a document that is a work order — should be preserved as-is for identity reasons. It **was** preserved as-is; renaming it would break the hash identity that row 393 and every prior report cite. This document recommends preservation and records the mismatch as a naming defect rather than a filing defect.

---

## 10. Statement of no changes

Changed by this filing: **registry row 393** (authorized administrative correction, quoted in §4) and the registry's declared scope and row-count lines. Added: **this document and its registry row**.

**Not** changed: the root `CORPUS-58` brief; `CORPUS-01`, `CORPUS-02`, `CORPUS-03`, `CORPUS-04`, `CORPUS-05`, `CORPUS-07`, `MEMORY-09`; `check_corpus.py`; any corpus item, label, fixture, acceptance condition, classification, or score; any runtime, schema, or model. `AM-01` remains unissued. `A-2a-R1` remains NOT IN FORCE. `D1` remains unadopted. Annotation remains blocked.

`G1` open · corpus frozen · expansion frozen · convergence 0 of 2.
