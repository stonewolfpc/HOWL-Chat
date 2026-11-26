# MEMORY-09 — S-3 Acceptance Review

## The P-1 logical representation passes the existing memory requirements, with explicit limitations

**Status:** S-3 acceptance review. **Logical evaluation only. No implementation, schema, code, migration, test,
deployment, requirement capture or canonical capture.** Acceptance is not implementation authorization.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Decision under review | **P-1** — `scope_changed` / `class_changed`; independent `scope_version` / `class_version` (P1D-01…P1D-10) |
| Criteria | **K-R12**, **K-R14**, **§6.4**; plus K-R13, K-R15, K-R2's integrity principle, S2H-12, and the authority/validity/cause boundaries |
| Verdict | **passes the logical acceptance criteria**, with the missing-latest limitation and open implementation details |
| Not authorized | schema, code, tests, requirements, canon, Howl/UX |
| Tests executed | **0** — the matrix is a logical walk-through |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Task and boundary

**S3A-01 (Observed — the owner's recommendation).** Proceed to S-3 acceptance: establish whether the P-1 logical
representation satisfies the existing memory requirements, especially K-R12, K-R14 and §6.4. Acceptance is not
implementation authorization.

**S3A-02 (Observed — what is under review).** P-1 as recorded in the owner decision: two dedicated change kinds,
independent per-dimension version sequences, formation baseline, contiguous versions, same-dimension serialization,
projection from the lifecycle stream, and a check that treats mismatch/duplicates as integrity errors
(P1D-04…P1D-10).

**S3A-03 (Observed — what this record does not do).** No schema, field definition, code, migration, test, probe,
model call, deployment, requirement or map change, or canonical capture. The map stays at **118**.

---

## 2. Acceptance criteria

**S3A-04 (Calculated — the tests).**

| # | Criterion | Source |
|---|---|---|
| 1 | the carrier is in the lifecycle stream | K-R12 locus |
| 2 | it covers `scope` and `knowledge_class` changes | K-R12 domain |
| 3 | current values are derived from the stream | K-R12 projection |
| 4 | the value is not a mutable authoritative field | K-R12 non-mutation |
| 5 | the stream is independently rebuildable | K-R14 |
| 6 | no stored cross-stream join | K-R13/K-R14 |
| 7 | the check is deterministic; mismatch is an integrity error | §6.4 |
| 8 | ordering is local; N-c is not settled | K-R15 |
| 9 | an unreconstructable projection is an integrity failure | K-R2 principle |
| 10 | context inactivity is not lifecycle status | S2H-12 |
| 11 | authority, validity and causes are not decided | §11/D6/D5/K-R9/K-R10/I13 |

---

## 3. K-R12 — locus, domain, projection, non-mutation

**S3A-05 (Calculated — locus).** `scope_changed` and `class_changed` are lifecycle-stream entries. **Pass.**

**S3A-06 (Calculated — domain).** The two kinds cover the two dimensions K-R12 names. A change to either dimension
has a dedicated entry. **Pass.**

**S3A-07 (Calculated — projection).** For each dimension, the current value is the formation baseline or the
highest-version change. The value is derived from the lifecycle stream, not assigned by a mutable field. **Pass.**

**S3A-08 (Calculated — non-mutation).** The stored current value is a cache compared by §6.4; it is not authoritative.
The authoritative input is the change entry. **Pass.**

**S3A-09 (Calculated — formation baseline).** The formation value is the baseline (P1D-06). If formation is not itself
a lifecycle entry, the baseline is the record's initial condition, not a post-formation change; the projection
requirement is satisfied for post-formation changes. Whether formation is an event is a pre-existing representation
question (A2-11/A2-12), not introduced by P-1. **Pass with a noted pre-existing question.**

---

## 4. K-R14 — independent rebuild, no stored join

**S3A-10 (Calculated — independent rebuild).** The class/scope projection is reconstructable from lifecycle entries
alone: baseline + highest version per dimension. It does not consult the acquisition stream. **Pass.**

**S3A-11 (Calculated — no stored join).** P-1 stores change entries, not a join of edge × memory. The knowledge
projection and the memory projection remain separate. **Pass.**

**S3A-12 (Calculated — mismatch is an integrity failure).** A stored/computed mismatch is an integrity error surfaced
to `inspection`; it is not silently resolved. This is the K-R14/K-R2 principle applied to the memory projection.
**Pass.**

**S3A-13 (Calculated — K-R2's literal scope).** K-R2's `derived_from_event_ids` clause is edge-specific and is not
imported onto the memory stream. The **integrity principle** is what P-1 satisfies. **Pass.**

---

## 5. §6.4 — the projection check

**S3A-14 (Calculated — determinism).** The projection is deterministic: sort by version per dimension, take the
highest, fall back to the baseline. Replaying the same entries yields the same value. **Pass.** "Cheap enough to run
on write" is a performance property; it is not evaluated here and belongs to implementation.

**S3A-15 (Calculated — mismatch).** Fold the entries, compute the value, compare to the stored cache. A mismatch is an
integrity error surfaced to `inspection`, never silently resolved in favour of either side. **Pass.**

**S3A-16 (Calculated — two streams, two checks).** P-1 adds change kinds to the **memory** stream; it does not create
a third stream. Acquisition events still project to edges; memory events project to status/scope/class. Neither
projection consults the other. **Pass.**

**S3A-17 (Calculated — no status on the edge).** P-1 does not put class/scope on `KnowledgeEdge`. **Pass.**

---

## 6. K-R13, K-R15 and the boundary rules

**S3A-18 (Calculated — K-R13).** Knowledge × memory remains a read-time join; P-1 does not store it. **Pass.**

**S3A-19 (Calculated — K-R15).** P-1's `scope_version` / `class_version` order only the change entries. They do not
order the four status kinds and do not decide whether the status stream needs an order key. **N-c remains open.**
**Pass with the N-c boundary preserved.**

**S3A-20 (Calculated — S2H-12).** P-1 has no context-activation field. A memory removed from a scene is not thereby
retracted, superseded or otherwise status-changed. **Pass.**

**S3A-21 (Calculated — authority, validity, causes).** A change entry records a value; it does not define who may
cause the change, whether it is valid, or which causes are permitted. D5/K-R9/K-R10/I13 remain outside the
representation. **Pass, with the boundary stated.**

**S3A-22 (Calculated — no status/class conflation).** `class_changed` is a separate kind from the four status kinds.
A class change is not disguised as a lifecycle-status transition. **Pass.**

---

## 7. Failure-mode matrix

**S3A-23 (Calculated — logical walk-through; not executed).**

| Scenario | Expected P-1 behavior | Acceptance |
|---|---|---|
| formation baseline | projection = formation value | pass |
| scope change v1 | projection = v1 value | pass |
| class change v1 | class projection = v1; scope unaffected | pass |
| both dimensions change | independent versions; both projections correct | pass |
| stored cache mismatch | integrity error to `inspection` | pass |
| duplicate `(memory, dimension, version)` | integrity error/conflict | pass |
| internal gap (v2 missing, v3 present) | incomplete/untrusted; integrity error if completeness known | pass with limitation |
| out-of-order (v3 before v2) | hold v3; apply when v2 arrives; recoverable | pass |
| missing latest (v3 absent) | undetectable by sequence; projection may be stale | pass with limitation |
| rebuild from stream alone | same projection; no cross-stream input | pass |

**S3A-24 (Observed — no execution).** The matrix is logical. No code, schema, probe or model call was run.

---

## 8. Limitations and open questions

**S3A-25 (Observed — preserved limitations).**

| Limitation | Status |
|---|---|
| missing internal version | detectable as a gap; integrity error only when completeness is known; otherwise incomplete/untrusted |
| missing latest event | **not detectable by sequence alone**; no later version exists |
| out-of-order arrival | recoverable; hold the higher version pending the predecessor |
| duplicate version | integrity error/conflict |
| silent repair | forbidden; untrusted projection is not authoritative |
| completeness marker | **open** — how the system knows a sequence should be complete |
| recovery policy | **open** — how long to wait before escalating |
| concurrency mechanism | **open** — serialization vs compare-and-swap; the invariant is fixed |
| formation as event | **pre-existing question**, not P-1-specific |

**S3A-26 (Calculated — none of these is an acceptance failure).** The limitations are consistent with K-R12/K-R14/
§6.4: the requirements govern the **recorded** stream, not the detection of an unrecorded change. The open items are
implementation/design details, not violations of the logical contract.

---

## 9. Verdict

**S3A-27 (Calculated — the acceptance result).** **P-1 passes the logical acceptance criteria** (S3A-04). It satisfies
K-R12's locus, domain, projection and non-mutation; K-R14's independent rebuild and no-stored-join; §6.4's
deterministic check and integrity-error handling; K-R13; K-R15's local-order boundary; K-R2's integrity principle; and
S2H-12's context/status separation. **The missing-latest limitation and the open completeness/recovery/concurrency
questions are preserved, not waived.**

**S3A-28 (Observed — no requirement capture).** The acceptance does **not** add a row to the 118-row map. It evaluates
the existing K-R12 obligation; it does not create a new conformance requirement.

**S3A-29 (Observed — no implementation authorization).** Acceptance is not permission to build. Schema, fields, code,
migrations, tests, deployment, requirement capture and canonical capture each require separate authorization
(S2R-25; P1D-25).

---

## 10. What acceptance does not authorize

**S3A-30 (Observed).**

| Unauthorized | State |
|---|---|
| physical schema / field definitions | **0** |
| code / migrations / implementation | **0** |
| tests / probes / model calls / deployment | **0** |
| requirement capture / map change | **0**; **118** |
| canonical capture | **0** |
| Howl state machine / UX | **0** |
| authority / validity / causes | not decided |
| N-b / N-c / N-a / S-1 | not decided |
| §11 / O-8a / v9 domain | not decided |

---

## 11. Traceability

**S3A-31 (Observed).**

| What | Where |
|---|---|
| P-1 decision | P1D-01…P1D-10 |
| P-1 limitations | P1D-16…P1D-21 |
| K-R12 | r33 `:111` |
| K-R14 | r33 `:113` |
| §6.4 | r33 `:366`–`:367` |
| K-R13 | r33 `:112` |
| K-R15 | §6.4 |
| K-R2 | r33 `:83` |
| S2H-12 | addendum |
| Failure taxonomy | S2F-06/07 |
| Adversarial equivalence | S2F-18…S2F-21 |
| Acceptance gate | S2R-24; P1D-27 |

---

## 12. Record status

**S3A-32 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Acceptance result | **passes the logical criteria**, with limitations |
| Implementation / tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **acceptance review** — no canonical capture |

**End of review.** **P-1 passes the logical S-3 acceptance criteria** (S3A-27): it satisfies K-R12, K-R14, §6.4,
K-R13, K-R15's local-order boundary, K-R2's integrity principle and S2H-12, with the missing-latest limitation and the
open completeness/recovery/concurrency questions preserved (S3A-23…S3A-26). **Acceptance adds no requirement and
authorizes no implementation** (S3A-28/29). **Schema, code, tests, requirement capture and canonical capture remain
separate authorizations** (S3A-30).
