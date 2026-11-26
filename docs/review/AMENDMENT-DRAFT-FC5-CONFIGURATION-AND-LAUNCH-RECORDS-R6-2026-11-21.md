# AMENDMENT DRAFT — FC-5 Configuration Decision and Launch Construction Records, Revision 6

**Status:** **OWNER-APPROVED AMENDMENT DRAFT — REVISION 6 — WORKING COPY — NOT APPLIED.**
**Date:** 2026-11-21
**Target paper:** `docs/architecture/contracts/FC-5-CAPABILITY-RESOURCE-AND-HANDOFF-CONTRACT-2026-11-21.md` (213 lines, 12,326 bytes)
**Target paper SHA-256 (observed 2026-11-21; identical to the value in register §6):** `c245096fbb58c8832bb9a87a1ebf27d0fa1598963b24e4c3a618025de2fcc828`
**Contract register:** `docs/architecture/contracts/README.md`
**Document registry:** `docs/DOCUMENT-REGISTRY.md` (its FC-5 row states the paper is not revised after the baseline)
**Baseline commit:** `0bc3b091b5023ac867f6d060bd56c481d6f84a3a`
**Scope:** documentation-only working copy. No canonical edit, no register or document-registry edit, no hash, pin, commit reference or status change, no M4 action, no implementation.

**Disposition:** The amendment's substance was approved by owner review of revision 6. This document persists the approved text so that it exists outside conversation history. It applies nothing. The canonical FC-5 paper still lacks every clause below, and its registered hash still describes the unamended text.

---

## 1. Authorization boundary

This working copy is a draft. Its existence establishes only that the approved wording is now recorded in a file.

It does **not**:

- amend `FC-5-CAPABILITY-RESOURCE-AND-HANDOFF-CONTRACT-2026-11-21.md`;
- modify the contract register or `docs/DOCUMENT-REGISTRY.md`;
- change any hash, pin, commit reference or status entry;
- change or evaluate M4 status;
- resolve `FC5-CONFIG-RETENTION`;
- resolve the criterion-ID question inside FC-6;
- modify FC5-01;
- implement configuration records, identity mechanisms or storage infrastructure.

Four separate milestones remain, and none is satisfied by this document: **application** to the canonical paper, **reconciliation** of the register and document registry, any **M4** change, and **implementation**. The runtime defect the amendment addresses — request, default and grant conflation in `applySafeDefaults`, discarded grant justification and unretained resolved arguments — remains documentable and prioritizable independently of all four.

## 2. Source of the approved text

Revisions 1 through 6 existed only as conversation text before this file. Nothing was applied at any revision, and no file recorded any earlier draft.

| Revision | Change | Status of its text |
|---|---|---|
| 1 | Initial amendment: identifier note, §4 lead-in, FC5-27–FC5-31, FC5-26a, placement map | Superseded |
| 2 | Retry wording rewritten; unsupported evidence attribution removed from FC5-30; FC5-31 moved into §4 after FC5-30; identifier note corrected; retention wording separating definition from durability | Superseded |
| 3 | Option 2 derivation provenance with the FC-6 limitation; identifier note separating existing from new suffixes; retention claim scoped | Superseded |
| 4 | Build identity stated as association rather than causation; omission versus substitution separated | Superseded |
| 5 | Equal-default substitution case: substitution recorded as `transformed` with its kind determined by comparing effective values | Superseded by revision 6 in FC5-29 only |
| 6 | FC5-29: identity continuity instead of persistence; failed construction records made subject to a future retention policy rather than asserted as retained | **Approved — the text below** |

Revision 6 changes FC5-29 and nothing else. Every other clause below is revision 5 wording, retained unchanged.

## 3. Front matter addition — unnumbered

Placed immediately after the existing `**Labels:**` line and before the `---` rule that follows it. This mirrors the identifier note at FC-2 line 17.

> **Identifier note:** Existing letter-suffixed clauses include `FC5-22a` and `FC5-22b`. This amendment adds `FC5-26a` as a related clause without renumbering any existing clause, so cross-references from other papers and the register remain valid. Clauses numbered above the highest existing clause (`FC5-27`–`FC5-31`) continue that sequence and are placed by subject rather than by number: numbering follows the highest existing clause while document order follows subject, so FC5-27–FC5-31 appear consecutively in §4 and FC5-26a appears in §10 after them despite its lower number. Identifiers of the form `FC5-CONFIG-RETENTION` identify questions, not clauses, and are cited as such.

## 4. §4 lead-in — unnumbered

Placed under the `## 4. Handoff phases` heading, above `FC5-12`. The heading text is unchanged.

> This section defines the handoff phase sequence and the two configuration records that accompany it: a record of the configuration decision taken for an operation, and a record of what Howl constructed for a launch attempt. The records are defined at the end of this section (FC5-27–FC5-31). They are not phases of the sequence in FC5-12, and no row of that table is satisfied by either record.

## 5. §4 clauses — placed after FC5-15, before §5

> **FC5-27 (Proposed — configuration decision record).**
>
> An operation's configuration is decided in three distinguishable stages, and each stage is recorded:
>
> | Stage | Recorded content |
> |---|---|
> | User intent | What the user asked for, captured at the operation boundary before application defaults are merged |
> | Application defaults | The defaults Howl applied, each named with the identity of the policy that supplied it, distinguishable from user intent |
> | Configuration grant | The configuration values Howl permits for this request, the adjustment made to each value where the grant diverges from the request, the stated reason for each adjustment, and the identity of the granting decision |
>
> User intent is not reconstructed from merged or persisted settings. A stored settings value that already contains an application default is not evidence that the user chose that value. Where intent was not captured — for example because the operation began from persisted settings rather than from a request — the record states that, and does not present a merged value as intent.
>
> The configuration grant is not FC-5 admission. Admission decides whether a workload may proceed under capability and resource policy. The configuration grant decides which configuration values Howl permits for that request. The two decisions may constrain one another, but neither substitutes for the other and neither is evidence for the other. A grant is also not a reservation state: it is not recorded in FC5-09's `state` values, and an existing reservation implies no grant.
>
> The configuration grant is not a claim-ladder level (FC5-22a). It is a configuration decision and carries no evidence maturity of its own.
>
> This clause defines the record's content and the distinction between its stages. It does not specify which component writes the record, where it is stored, or in what format beyond the distinctions above (FC5-30, FC5-31).
>
> ---
>
> **FC5-28 (Proposed — launch construction record).**
>
> For each launch attempt, Howl records what it constructed, expressed per relevant setting or argument:
>
> | Disposition | Required content |
> |---|---|
> | `emitted` | The option and the value emitted |
> | `omitted` | The option and an explicit reason |
> | `transformed` | The option, the requested value, the emitted value, the kind of transformation, and its provenance |
>
> The kind of transformation distinguishes a semantic value change from a representation-only change. A semantic change alters the effective value or its meaning. A representation-only change expresses the same effective value differently — spelling, ordering, formatting, or the field or flag form used — and does not alter what the runtime would be asked to do. The two are recorded as different kinds, because they carry different consequences for a later reader.
>
> An omission reason preserves the distinction that explains it. At minimum the record separates: never requested; requested, not emitted, and left at the application default; conditionally excluded by a stated condition; and requested but dropped as unsupported. These are distinct reasons and are not collapsed into one unspecified reason.
>
> The second of these is an omission only while the default remains in effect and no value is emitted in place of the request. Where Howl substitutes an application default for the requested value and emits a value, the outcome is `transformed`, not `omitted`. Its kind is determined by comparing the effective values, not by the fact of substitution: where the substituted default differs from the requested value in effect, the change is semantic; where the substituted default is equivalent to the requested value, the substitution is still recorded, but the kind is representation-only, because the effective configuration did not change. A substitution is never described as a semantic change on the strength of the substitution alone.
>
> Where construction performs a derivation — applying a default, clamping a value, resolving an unsupported request, conditionally excluding an option — the record states that derivation's provenance: the stage that performed it, the identity of the governing policy where one exists, and the build identity associated with the construction record, reusing FC-6's `build_identity` with the meaning FC6-06 and FC6-07 already give it. Where no governing policy identity exists, the record states that rather than inventing one.
>
> The build identity recorded here is the one associated with the record. Recording it does not assert that this build performed the derivation for this attempt. Where that relationship is established, the record states it at the strength the evidence supports; where it is not, the record states the association only and its provenance is correspondingly weaker. This clause creates no identity for a derivation and no registry of derivation rules.
>
> Record B describes construction only. What that construction does not establish is stated in FC5-30.
>
> ---
>
> **FC5-29 (Proposed — identities, references and retry).**
>
> Each configuration decision record has its own identity. Each launch construction record has its own identity and references the configuration decision record it derives from.
>
> The construction identity is FC-5's own and is capability-agnostic: it applies to text, vision, audio and other operations alike, and there is one identity per construction attempt. It is not `request_id`, `prompt_plan_id` or `reservation_id`, which identify a prompt-dispatch attempt (FC2-05b), a composition artifact and a reservation respectively; none of those scopes is a launch construction attempt.
>
> A reservation reference is recorded when a reservation has been established for the attempt. It is never universally mandatory. FC5-12 places reservation before load and does not state where construction falls relative to it, so a construction record is not invalid for want of a reservation, and an absent reference is not a claim that no reservation existed.
>
> One configuration decision can support more than one construction attempt. The configuration decision identity remains the common reference across those attempts; each construction record retains its own distinct identity. This identity relationship does not establish any retention or persistence guarantee (FC5-31).
>
> A genuine change to the configuration decision produces a new decision identity, and the earlier decision record is not rewritten to represent it. A failed construction record is subject to whatever retention policy is eventually established (FC5-31). This clause does not establish that policy.
>
> ---
>
> **FC5-30 (Proposed — configuration records are not runtime evidence).**
>
> Neither record authorizes, implements or proves any of the following, and neither may be cited as evidence for any of them:
>
> - configuration recording or storage infrastructure;
> - observation or monitoring of a process;
> - collection of runtime telemetry;
> - verification that an argument, option or value was accepted by a runtime;
> - backend or execution identity;
> - allocation, loading, hardware execution, or workload readiness.
>
> A configuration grant is a decision about permitted values. A launch construction record is a description of what Howl constructed. Neither establishes what a runtime accepted, what was allocated, what loaded, what executed, or whether a workload became ready. Those runtime facts remain governed by their applicable contracts. This clause attributes no evidence obligation to any other contract or phase, and does not claim that any phase or paper owes backend identity, hardware-execution evidence or runtime telemetry.
>
> Readiness remains FC-5's own phase claim (FC5-22). This clause does not extend FC-5 to backend identity, hardware execution or runtime telemetry, and neither record becomes readiness evidence because FC-5 owns readiness. The records are not promoted into readiness claims, not even when they are complete and mutually consistent.
>
> ---
>
> **FC5-31 (Unresolved — configuration-record retention).**
>
> This clause defines no retention behaviour and does not qualify FC5-27–FC5-30. Those clauses require the records to be defined and to carry their identities and references; they do not establish that a record is written to durable storage, that it survives a process exit or restart, that it is retained for any duration, where it is held, what owns it, what scope it has, or how it is redacted.
>
> `DQ-HRA-04` asks which **reservation** records survive a crash; it does not answer these questions for configuration records. D2's retention-by-default governs authored content and does not reach application settings or launch records. Configuration-record retention is therefore not established by the provisions governing these records: nothing in FC5-27–FC5-30, `DQ-HRA-04` or D2's retention-by-default requires these records to persist. This clause does not decide whether any other provision constrains where or how they may be stored.
>
> Open question, minted here: `FC5-CONFIG-RETENTION` — what retention, storage, scope and redaction policy applies to configuration decision records and launch construction records?
>
> This clause is not a `DQ-HRA-*` question, is not inherited from `HRA-DB-2026-11-21-01`, and decides none of the points it lists.

## 6. §10 addition — after FC5-26, before the closing line

> **FC5-26a (Proposed — boundary pointer).** FC5-26's exclusion list is not exhaustive. The configuration records defined in §4 introduce boundaries that list does not name: FC5-30 states what those records neither authorize nor prove, and it applies to both records. This clause points to that boundary and does not restate it.

## 7. Placement map

| Location | Addition | Numbered? |
|---|---|---|
| Front matter, after the `Labels` line, before the `---` rule | Identifier-convention note | No |
| §4, under the `## 4. Handoff phases` heading, above FC5-12 | Lead-in paragraph | No |
| §4, after FC5-15, before §5 | FC5-27, FC5-28, FC5-29, FC5-30, FC5-31 | Yes |
| §10, after FC5-26, before the closing line | Boundary pointer | FC5-26a |

Unchanged: FC5-01 through FC5-26 in their entirety, the `## 4. Handoff phases` heading text, §8, §9, and the closing sentence.

## 8. Draft-status and approval boundaries

The amendment is approved **as a draft**. That approval covers the wording in sections 3 through 6 and the placement in section 7. It covers nothing else.

Previously settled and not reopened by this revision: FC-5 owns the amendment through additive clauses; both records sit in §4 after FC5-15 and before §5; the records are defined by numbered clauses with an unnumbered lead-in; the identifier form continues the sequence from the highest existing clause; identity is FC-5-owned and capability-agnostic, one per construction attempt; the grant is a record stage and not a claim-ladder level; non-authorization is clause-local in FC5-30 with a non-restating pointer in FC5-26a; retention is recorded as unresolved only.

## 9. Open at approval — not resolved here

- `FC5-CONFIG-RETENTION` is minted as a question identifier. Its content is a question, and this amendment answers none of it.
- FC6-07's `contract` field requires a paper and a criterion ID. FC5-27–FC5-31 are citable clauses; this amendment makes no claim that a clause ID satisfies that field, and the question stays open inside FC-6.
- FC5-01 is unchanged, including its incompleteness. The second decision it omits is named in FC5-27, not added to FC5-01.
- `DQ-HRA-04` is unchanged and remains a reservation question.

## 10. Downstream consequences — identified, not executed

If the amendment is later applied under its own authorization, these entries would need reconciliation. None has been touched:

- register §1 Status for FC-5;
- register §2 item 5, which describes the paper;
- register §6 hash row for FC-5 (`c245096f…`);
- register §6's sentence stating that FC-3, FC-5 and FC-6 are unchanged from baseline commit `0bc3b091…`;
- the FC-5 row in `docs/DOCUMENT-REGISTRY.md`.

Applying the text would change the paper's byte count and hash. This working copy asserts no new hash and registers nothing.

## 11. Validation of this working copy

- Read-only inspection preceded writing: the contracts directory, the `docs/review` draft conventions, the existing `AMENDMENT-DRAFT-*` and `AMENDMENT-PROPOSAL-*-R5` headers, the FC-5 front matter and §4 heading, the FC-2 identifier-note precedent, and the register's FC-5 rows.
- FC-5 on disk matches the register: 213 lines, 12,326 bytes, SHA-256 `c245096f…`. The canonical paper is unamended and unmodified by this task.
- The approved text is reproduced from revision 5 with revision 6's two FC5-29 corrections applied. Both corrections are present verbatim: the configuration decision identity remains the common reference across attempts with a distinct identity per construction record and no retention or persistence guarantee; and a failed construction record is subject to a future retention policy rather than asserted as retained.
- Every identifier is present: FC5-26a, FC5-27, FC5-28, FC5-29, FC5-30, FC5-31, `FC5-CONFIG-RETENTION`.
- Every cross-reference is preserved: FC5-09, FC5-12, FC5-22, FC5-22a, FC5-26, FC5-30, FC5-31, FC6-06, FC6-07, FC2-05b, `DQ-HRA-04`, `HRA-DB-2026-11-21-01`.
- No existing clause was renumbered and no placement changed.
