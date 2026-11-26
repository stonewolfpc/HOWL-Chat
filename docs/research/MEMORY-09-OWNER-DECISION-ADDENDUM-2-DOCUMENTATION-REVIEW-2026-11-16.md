# MEMORY-09 — Owner Decision Addendum 2 (documentation review dispositions)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-ODA-2` · **Version:** 1 · **Date:** 2026-11-16 · **Status:** **RECORDED — documentation direction only.**
**Type:** append-only addendum to `MEMORY-09-ODR-1`. It **does not amend, supersede, renumber or rewrite** that record, its addendum 1, or any artifact it references.
**Authority:** owner work order *MEMORY-09 GOVERNANCE DECISIONS AND DOCUMENTATION*, 2026-11-15. Six findings accepted, five accepted outright and one accepted with reframing.
**Affected artifacts:** the five reconstruction artifacts and the fixed review record, all identified by hash in §8.
**Effect:** records decisions and the authorised future documentation workflow. It **authorises no** correction to any existing artifact, no implementation, no execution, and does not authorise contract approval, G4 readiness or conformance. The replacement set and the supersession notice each require a subsequent work order.

---

## 1. Scope

This addendum records the owner's dispositions on the six findings of the documentation review of 2026-11-15, and the naming, scope and ordering rules that any correction must follow.

### 1.1 Identifier decision

**No `OD-M09-nn` identifier is assigned by this addendum.** The dispositions keep the review's own identifiers, `R-01` … `R-06`.

The reason is not convenience. The `OD-M09-19`–`26` band is already occupied by proposals in flight: `OD-M09-19` is proposed elsewhere for a definition of "cue-level", `OD-M09-20`–`22` are left free for records already in flight, and `OD-M09-23`–`26` are proposed and unissued. Issuing `OD-M09-19` here would collide with a live instrument for a different purpose.

Assigning a number the owner has not assigned would be the same class of error the review itself identified: a plausible identifier standing in for a real one. Nothing is renumbered, and no historical identifier is reassigned.

### 1.2 What this addendum is not

It is not a correction. It does not edit, replace or supersede any of the five reconstruction artifacts, and it does not touch the fixed review record. It records what the owner decided, so that the correction can later be made under a separate dated authorisation.

## 2. Decision register

| ID | Subject | Disposition |
|---|---|---|
| `R-01` | Third identifier example (`A30.6.8`) in artifact 1 | **ACCEPTED** — cross-version contamination from superseded rev.1 material |
| `R-02` | P1 instance inventory provenance in artifact 2 | **ACCEPTED WITH REFRAMING** — count retained, attribution corrected |
| `R-03` | "Answers them" wording in artifact 3 | **ACCEPTED** — `ADDRESSED-PENDING-REVIEW` establishes no adequacy |
| `R-04` | P1/P3 rationale attribution in artifact 4 | **ACCEPTED** — subject labels cannot support a rationale attributed to missing prose |
| `R-05` | "There is no implementation" in artifact 5 | **ACCEPTED** — too broad given repository harness code |
| `R-06` | Integrity wording | **ACCEPTED AS A CONTINUING CONSTRAINT** — see §3.6 |

## 3. Finding dispositions

### 3.1 `R-01` — accepted

`A30.6.8` is **cross-version contamination** from the superseded rev.1 handoff. The third identifier example and its associated count are **one defect, not two**.

Consequences, binding on any correction:

- The example is removed; the count is corrected to two.
- The identifier is **not** retained under a repaired count.
- Whether nested `A30` numbering existed in the contract remains `[TRANSCRIPT RECONCILIATION REQUIRED]`.

### 3.2 `R-02` — accepted with reframing

The four-instance count for P1 is **retained**. Its **provenance is corrected**: it comes from an earlier conversational review of the superseded handoff, **not** from the governing brief.

**Qualification, recorded as the owner stated it:** the earlier exchange cannot be independently recovered or re-verified from the five documents and handoff rev.3.1 alone. The correction must therefore **identify** that exchange as the source and preserve the distinction, rather than implying the count was present in the governing brief. The instances themselves remain unrecovered.

### 3.3 `R-03` — accepted

`ADDRESSED-PENDING-REVIEW` does **not** establish that a response adequately answers a finding. The reviewed artifact's "answers them" wording overstates the supplied status and must be replaced with a statement that responses are associated with the three headings and their adequacy is unreviewed.

### 3.4 `R-04` — accepted

P1/P3 **subject labels** cannot establish a rationale attributed to their **missing principle prose**. A label is not a statement. Any correction must present the prohibition as a supplied governance rule whose consistency with P1 and P3 is an inference from labels, pending transcript reconciliation.

### 3.5 `R-05` — accepted

The claim that no implementation exists is **too broad**. The repository contains evaluation harness code. What is supported, and all that may be claimed:

> No implementation conformance to the durability contract has been established, and the reconstruction authorised no implementation work.

### 3.6 `R-06` — accepted as a continuing constraint

Every artifact in this chain, including this addendum, must limit its integrity claim to exactly:

> The byte sequence reviewed is the byte sequence present at creation.

No artifact in the chain may describe any file as **unchanged**, **unmodified** or **uncorrupted** on the basis of two hash observations. Two observations establish equality of the observed byte sequences at those two points; they cannot exclude a change-and-restore between them.

## 4. Replacement naming convention

Replacements append `-REV-2` **immediately before `.md`**, preserving the original date.

| Original | Replacement |
|---|---|
| `MEMORY-09-DURABILITY-CONTRACT-2026-11-15.md` | `MEMORY-09-DURABILITY-CONTRACT-2026-11-15-REV-2.md` |
| `MEMORY-09-GOVERNING-PRINCIPLES-2026-11-15.md` | `MEMORY-09-GOVERNING-PRINCIPLES-2026-11-15-REV-2.md` |
| `MEMORY-09-REVIEW-LEDGER-2026-11-15.md` | `MEMORY-09-REVIEW-LEDGER-2026-11-15-REV-2.md` |
| `MEMORY-09-OPEN-DECISIONS-2026-11-15.md` | `MEMORY-09-OPEN-DECISIONS-2026-11-15-REV-2.md` |
| `MEMORY-09-DOCUMENTATION-RECONSTRUCTION-RECORD-2026-11-15.md` | `MEMORY-09-DOCUMENTATION-RECONSTRUCTION-RECORD-2026-11-15-REV-2.md` |

This is the **required** convention for all five replacements. It is required precisely because the repository convention is topic plus ISO date with no revision channel; without the token, a same-date revision would collide with the path it revises, and the review's subject set would describe files that no longer exist.

## 5. Replacement scope and status

| Requirement | Rule |
|---|---|
| Scope | A **complete second set of five**. No mixture of original and replacement versions |
| Original identification | Each replacement states its original path and **full SHA-256** |
| Finding identification | Each replacement lists the **one** finding ID it answers |
| Boundedness | Each states that the correction is bounded to that finding and that no other change was made |
| Status | Each is **`DRAFT — UNREVIEWED (REV-2)`** |
| Inherited approval | None. No replacement inherits review approval, completeness or conformance from its predecessor |

A replacement is a new artifact with its own identity. The REV-2 set will require its own identity capture and its own independent review; the supersession notice carries its hashes, so no separate manifest is required.

## 6. Supersession notice

To be created only after all five replacements exist:

`MEMORY-09-DOCUMENTATION-SUPERSESSION-NOTICE-2026-11-15.md`

It must identify:

1. The fixed review record by **full SHA-256**.
2. All five original artifacts by path and full SHA-256.
3. All five replacements by path, byte count and full SHA-256.
4. The corrections associated with `R-01` … `R-05`.
5. This addendum as the authorising record.
6. The date of authorisation.

It must **not** imply that the fixed review record was changed, revised or repopulated. The review record's supersession register stays empty permanently; supersession history lives in the notice.

## 7. Required ordering

1. **This addendum** — records the decision.
2. **The complete five-document REV-2 replacement set.**
3. **The supersession notice** — which cannot precede the replacements, because it must carry their actual hashes.

The notice therefore records an outcome; the authorisation lives here, not in the notice.

## 8. Preservation obligations

### 8.1 The fixed review record

| Field | Value |
|---|---|
| Path | `docs/research/MEMORY-09-DOCUMENTATION-REVIEW-2026-11-15.md` |
| SHA-256 | `52f627909615579cea28fb263e7e1cc8b5cd0215c62e74412311ce7f821b1b67` |
| Bytes | 16,252 |
| Status | `REVIEW RECORD — FINDINGS FIXED` |
| Obligation | Permanently fixed. Contents and hash must not change. Its supersession register remains empty |

### 8.2 The five original reconstruction artifacts

| # | Path | SHA-256 |
|---|---|---|
| 1 | `MEMORY-09-DURABILITY-CONTRACT-2026-11-15.md` | `1f2e3a2d64a366247e515867fc6f4d21927e5d637c42ccecd35faefa55243020` |
| 2 | `MEMORY-09-GOVERNING-PRINCIPLES-2026-11-15.md` | `0899e012cb4a13e207816ad5d741d9e86ebe946e4138197d9d760540e1dcff7b` |
| 3 | `MEMORY-09-REVIEW-LEDGER-2026-11-15.md` | `071530d0e996eba3c0b96a57e1316391504638ca86dddb9fe50e93c547d23bc4` |
| 4 | `MEMORY-09-OPEN-DECISIONS-2026-11-15.md` | `41858a23b20fac56eb31e3d35f9b4de4fe4181cac89772ae0e6c69e51bc57c25` |
| 5 | `MEMORY-09-DOCUMENTATION-RECONSTRUCTION-RECORD-2026-11-15.md` | `636ea2ceeb04fcc9a8b74ce6fb7fb42fb2ce8195f396901c2f374b2c4eef3ce7` |

**Obligation:** do not overwrite, edit, rename or delete the originals. They remain the reviewed subject set. Corrections produce new artifacts; they do not repair old ones.

At the time of filing, the byte sequences read at verification matched these values — see §10, which states the limit of that claim.

## 9. Authorised and not authorised

**Authorised by this addendum:** recording the six dispositions; recording the naming convention, replacement rules, notice requirements and ordering; the subsequent creation of the five REV-2 artifacts and the supersession notice **under a separate work order**; the identity capture and independent review of the REV-2 set **under separate authorisation**.

**Not authorised by this addendum:** creating the REV-2 set; creating the supersession notice; editing, renaming or deleting any original artifact or the fixed review record; editing source, harness, runtime or configuration; running tests, model calls or experiments; closing any `U-` item; reconciling mem1/S3; contract approval; implementation conformance; G4 execution or the clearing of any G4 blocker.

## 10. Provenance

| Layer | Fact |
|---|---|
| Owner dispositions | Given 2026-11-15; recorded 2026-11-16, in the work order recorded in §Authority |
| Recorded by | Buffy / Codebuff repository worker |
| Identity | Artifact-reported only; labels are provenance data, not lineage proof |
| Base record | `MEMORY-09-ODR-1` — not edited by this addendum |
| Addendum 1 | `MEMORY-09-ODA-1` — not edited, not superseded |
| Model executed | **No.** No inference, benchmark, model call or experiment was performed in producing this record |
| Contract status | Revision 23.8 remains `DRAFT-UNREVIEWED`; no clause of it was amended |
| Integrity | The byte sequences of the artifacts named in §8 were read at verification and matched the values recorded there. That establishes equality of the observed byte sequences at those two observations only; per `R-06` and §3.6, it does **not** establish that any file was unchanged between them |
| Registry | Append-only convention: `docs/DOCUMENT-REGISTRY.md`. **Row not written** — see §11 |

## 11. Deviations and open questions

| # | Item | State |
|---|---|---|
| 1 | `docs/DOCUMENT-REGISTRY.md` row for this addendum | **Not written.** The established convention appends a registry row for each artifact, but the work order authorises creating this addendum only and does not authorise editing existing documents. The obligation is recorded here rather than discharged silently. A follow-up authorisation is needed, and the row can then carry this file's hash and byte count |
| 2 | Addendum numbering in the filename | This addendum is scoped by subject rather than by an `OD-M09` range, because no range can be claimed without collision (§1.1). If the owner prefers range-based naming, the file would need renaming under a fresh authorisation — not a content change |
| 3 | `R-02` source verification | The earlier exchange is identified as the source on the owner's statement. It cannot be re-verified from the five artifacts and handoff rev.3.1, and this record does not present it as governing-brief material |
| 4 | REV-2 status vocabulary | `DRAFT — UNREVIEWED (REV-2)` is recorded as the required status; no review of the replacement set exists or is implied |
| 5 | Recording date | The work order is dated 2026-11-15. This file was created at `2026-11-16 00:03:46 -0600`, after midnight. It was therefore named and dated **2026-11-16** rather than backdated, which differs from the date in the work order. The correction was made before this record was reported, and the affected line is the `Identity` date in the header. Reverting to the work order's date is a one-command rename if the owner prefers it |
