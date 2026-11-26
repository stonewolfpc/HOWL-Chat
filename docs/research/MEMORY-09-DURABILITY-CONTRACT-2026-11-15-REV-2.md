# MEMORY-09 — Durability Contract (reconstruction, REV-2)

**Identity:** `MEMORY-09-DURABILITY-CONTRACT-2026-11-15-REV-2` · **Version:** REV-2 · **Set:** artifact 1 of 5, replacement revision

**Supersedes:** `docs/research/MEMORY-09-DURABILITY-CONTRACT-2026-11-15.md`, SHA-256 `1f2e3a2d64a366247e515867fc6f4d21927e5d637c42ccecd35faefa55243020`. The original is **preserved and not modified**; it remains the reviewed subject set.

**Finding answered:** `R-01` — a third identifier-variant example (`A30.6.8`) was recorded as handoff-supplied when the governing brief supplied only `A22.1′` and `A19.4″`. The example and its count were one defect, not two.

**Corrections applied:** `R-01` — §3 count corrected from three to two; the `A30.6.8` row **removed** from the §5 mapping table with no renumbering and no replacement; pending mappings reduced from three to two; the §4.1 paragraph replaced so that the nested form is attributed to an earlier superseded draft rather than to the brief. `M-1` — the §9 bullet rewritten so that it claims no implementation conformance has been established rather than claiming no implementation exists, with execution and testing scoped to this reconstruction.

**Boundedness:** the corrections are limited to `R-01` and `M-1`. No other clause, mapping, registry entry, cross-reference or statement is changed. Documentation correction only — **no contract approval, implementation, execution or G4 authority is created.**

**Status: DRAFT — UNREVIEWED (REV-2). Reconstructed from a review handoff, not from the source transcript. Textual review only; no execution, no implementation, no G4 work performed or authorised. Marked items were not supplied and must not be invented.**

**Implementation status: NOT STARTED / NOT AUTHORIZED**

**Construction basis:** reconstruction handoff rev.3.1, supplied through conversation and having no repository presence; plus read-only repository preflight performed 2026-11-15. The larger part of the historical contract is not available to this reconstruction — see §3.

---

## 1. What this document is

A structured inventory of the durability contract that was reviewed between revision 1 and revision 23.8 in a conversation. It records the clause families, the identifiers that were supplied, and the registries that were designed — and it marks, explicitly, everything that could not be recovered.

It is **not** the contract. No clause in this document carries prose unless that prose was supplied. It is not a normative source, it is not approved, and it does not supersede or amend anything.

Revision 23.8, the newest revision of the contract, remains **`DRAFT-UNREVIEWED`**. Nothing here upgrades it.

## 2. Markers

| Marker | Meaning |
|---|---|
| `[PROSE PENDING]` | The clause, statement or row exists in design; its wording was not supplied to this reconstruction |
| `[TRANSCRIPT RECONCILIATION REQUIRED]` | The item requires the original transcript to be established at all |
| `[ID MAPPING PENDING]` | The identifier's historical form or its stable successor is not established |

A marker is a reconciliation task, not a defect of the design. It is also not a finding that the item is absent from the design.

## 3. Source classification

| Class | Content in this document |
|---|---|
| Handoff-supplied | Clause family names and supplied ranges; the eight supplied A30 topics; the four registry structures and the expected total of 16; two identifier-variant examples; three supplied corrections |
| Repository-derived | Nothing that bears on the durability contract. The architecture-stage contract (`MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md`) and its generated extract are a **different artifact set**; their tables were not substituted here |
| Missing | Clause prose for every clause; per-clause status; the complete identifier history; all registry rows; the three findings raised against rev.23.8 |

The absence of repository-derived contract content is a deliberate restraint, not an oversight. Substituting the architecture-stage contract's tables for the durability contract's registries would have produced a document that looked substantiated and was not.

## 4. Clause family inventory

Eleven families were supplied. No clause prose was supplied for any of them.

| Family | Supplied range | Supplied subject | Prose |
|---|---|---|---|
| `A19` | `A19.1`–`A19.4` plus amendment variants | Ordinal namespace and partition | `[PROSE PENDING]` |
| `A22` | `A22.1`–`A22.9` plus amendment variants | Admission accounting | `[PROSE PENDING]` |
| `A23` | `A23.1`–`A23.3` | Reconciliation generation | `[PROSE PENDING]` |
| `A24` | `A24.1`–`A24.6` | Immutability and derived void | `[PROSE PENDING]` |
| `A25` | `A25.1`–`A25.3` | Reconciliation applicability | `[PROSE PENDING]` |
| `A26` | `A26.1`–`A26.9` | Classification, snapshots, episodes | `[PROSE PENDING]` |
| `A28` | `A28.1`–`A28.3` | Audit journal and witnesses | `[PROSE PENDING]` |
| `A29` | `A29.1`–`A29.5` | State machines and legal pairings | `[PROSE PENDING]` |
| `A30` | `A30.1`–`A30.23` | Durability (see §4.1) | `[PROSE PENDING]` |
| `C3` | Not supplied | Generation chain | `[PROSE PENDING]` |
| `C4` | Not supplied | Canonicalization, digest roles, integer domain | `[PROSE PENDING]` |

### 4.1 Supplied A30 topics

Eight topics were supplied as subclause subjects. They support the claim that the durability family covers these areas; they do not establish which number carries which topic.

| Topic | Clause number |
|---|---|
| Append behaviour | `[ID MAPPING PENDING]` |
| Operation log | `[ID MAPPING PENDING]` |
| Epochs | `[ID MAPPING PENDING]` |
| Blocking-cause log | `[ID MAPPING PENDING]` |
| Commit registration | `[ID MAPPING PENDING]` |
| Evidence | `[ID MAPPING PENDING]` |
| Provenance | `[ID MAPPING PENDING]` |
| Receipt index | `[ID MAPPING PENDING]` |

"A30" was supplied with a stated range only. No nested subclause numbering was supplied. A nested form appeared in an earlier superseded draft of the handoff and is not reproduced here; whether nested numbering existed in the contract is `[TRANSCRIPT RECONCILIATION REQUIRED]`.

## 5. Identifier mapping

No stable identifier has been assigned. Assigning one from a family range and a topic name would manufacture traceability that the material does not support, and renumbering without the historical text is how distinctions get erased.

Verified mappings: **0**. Pending: **2**, listed below; family ranges are not per-clause identifiers and are not mapped.

| Historical ID | Stable ID | Revision | Disposition | Evidence |
|---|---|---|---|---|
| `A22.1′` | `[ID MAPPING PENDING]` | Unknown | Unresolved | Supplied as an example of prime notation; exact historical form unverified, possibly transcription noise |
| `A19.4″` | `[ID MAPPING PENDING]` | Unknown | Unresolved | As above, for double-prime notation |

Recorded disposition of the notation itself: prime and double-prime forms appear across `A22` and `A30`, and no key to their meaning was supplied. They may denote amendments, insertions, or nothing at all. That question is `[TRANSCRIPT RECONCILIATION REQUIRED]`.

## 6. Registry appendices

These four registries were designed during the review. Each is reported as verified rows against expected rows, and the two numbers are not combined.

### 6.1 Artifact identity/integrity registry

| Measure | Value |
|---|---|
| Expected rows (design claim) | 16 |
| Verified rows | **0** |
| Missing rows | 16 — all `[TRANSCRIPT RECONCILIATION REQUIRED]` |

The expected total of 16 is a design claim, not evidence that sixteen entries were ever populated or can be recovered.

### 6.2 Witness table

| Measure | Value |
|---|---|
| Expected rows | Not supplied |
| Verified rows | **0** |
| Missing rows | `[TRANSCRIPT RECONCILIATION REQUIRED]` |

### 6.3 Failure-state table

| Measure | Value |
|---|---|
| Expected rows | Not supplied |
| Verified rows | **0** |
| Missing rows | `[TRANSCRIPT RECONCILIATION REQUIRED]` |

### 6.4 Operation-class × scope matrix

| Measure | Value |
|---|---|
| Expected rows / columns | Not supplied |
| Verified cells | **0** |
| Missing cells | `[TRANSCRIPT RECONCILIATION REQUIRED]` |

The four registries were designed to be the single authority for their own contents. With zero verified rows, this document establishes that the registries were specified, and nothing about what they contain.

## 7. Supplied corrections, and what they do and do not establish

Three corrections were supplied. Each is a change made in the review; none is a reviewer-confirmed closure.

| Correction | Affected family | Clause location | Closure |
|---|---|---|---|
| Journal sequence upper bound changed from `2^63−1` to `2^53−1` | `A28` (journal), probably | `[ID MAPPING PENDING]` | Not established — see review ledger `L-04` |
| `all_appends` renamed `all_state_changes`, with enumerated permitted sets | `A19` or `A22`, unidentified | `[ID MAPPING PENDING]` | Not established — `L-05` |
| `read_back_verified_at` moved out of commit registration into `CommitVerification` | `A30` (commit registration) | `[ID MAPPING PENDING]` | Not established — `L-06` |

The family attributions above are inferences from the subject matter of each correction. They are the only inferences in this document, they are marked as such here, and they should be checked against the transcript rather than relied upon.

## 8. Clause status vocabulary

Only three states are used in this document:

| State | Meaning |
|---|---|
| Supplied | The item's existence and subject came from the handoff |
| Pending | The item exists in design; its content was not supplied |
| Unverified | An identifier, mapping or attribution was supplied in a form that could not be confirmed |

No clause is recorded as accepted, closed, or superseded. Those states require either the transcript or a named reviewer decision, and neither is available. A clause that appears to be superseded by a later one is not recorded as superseded here: **the existence of a later clause is not evidence that an earlier one was withdrawn.**

## 9. What this document does not establish

- That any clause says what its title suggests.
- That the clause numbering was ever stable.
- That rev.23.8 is complete, correct, reviewed, or freeze-ready.
- That any registry was populated, or that its design was sound.
- That the contract was implemented or that any implementation conforms to it. No implementation conformance to the durability contract has been established, and this reconstruction authorised no implementation work. The durability contract was not executed or tested as part of this reconstruction. No durability-contract instrument was run for it, and no model call was made in producing it; earlier evaluation activity recorded elsewhere in the repository is outside this claim.
- That the interface and storage work implied by the contract is authorised. It is not.

## 10. Cross-references

| Document | Relationship |
|---|---|
| `MEMORY-09-GOVERNING-PRINCIPLES-2026-11-15.md` | Artifact 2 of 5 — the principles this contract applies; referenced by ID only, not restated |
| `MEMORY-09-REVIEW-LEDGER-2026-11-15.md` | Artifact 3 of 5 — the findings these clauses answer |
| `MEMORY-09-OPEN-DECISIONS-2026-11-15.md` | Artifact 4 of 5 — what remains undecided |
| `MEMORY-09-DOCUMENTATION-RECONSTRUCTION-RECORD-2026-11-15.md` | Artifact 5 of 5 — provenance, limits and dates for the set |
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` | Protected. A different artifact; cited by section where relevant, never summarised here |
| `MEMORY-09-CONTRACT-NORMATIVE-2026-11-15.md` | Protected. Generated extract of the above; **not** an extract of this contract |
| `MEMORY-09-PREREGISTRATION-FREEZE-2026-11-15.md` | Protected. Establishes which artifacts are frozen; this document is outside that set |
