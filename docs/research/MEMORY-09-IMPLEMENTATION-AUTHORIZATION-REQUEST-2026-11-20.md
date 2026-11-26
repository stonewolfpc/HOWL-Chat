# MEMORY-09 — Implementation Authorization Request

## A bounded request for L-2 and P-1 integration; readiness is not authorization

**Status:** **implementation-authorization request only.** No implementation, schema change, code, migration, test,
requirement capture or canonical capture is authorized by this document. P-1 and S-3 remain closed; Decision B remains
untouched.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Request | authorize implementation of the L-2 formation representation and its P-1 integration, within the scope below |
| Basis | L-2 adopted; P-1/S-3 closed; consistency review passed with three guardrails |
| Not authorization | **yes** — the owner must approve; this record only requests |
| Excluded | Decision B, N-b, N-c, formation authority, cross-stream ordering, Howl state machine |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The request

**IAR-01 (Observed — what is requested).** The owner is asked to authorize implementation of:

1. the **L-2 formation representation**;
2. its integration with the **two P-1 change dimensions**;
3. the **three consistency guardrails**;
4. the **S-3 rebuild and integrity behavior**;
5. **verification of the unresolved-status exclusion**.

**IAR-02 (Observed — why now).** L-2 is adopted (L2D-01); P-1 and S-3 are closed; the consistency review passes with
the baseline/change distinction, the status-slot exclusion and the provenance guardrail (LCR-26/27).

**IAR-03 (Observed — readiness is not authorization).** The logical decisions are complete enough to **request**
implementation authority. They do not grant it. No implementation starts until the owner approves this request or a
narrowed version of it.

---

## 2. Scope — what the authorization would cover

**IAR-04 (Observed — the included set).**

| # | Included |
|---|---|
| 1 | **L-2 formation representation** — one lifecycle-side formation entry with `scope_baseline`, `class_baseline`, `status_slot` |
| 2 | **P-1 integration** — `scope_changed` / `class_changed`; independent `scope_version` / `class_version`; first change version 1; contiguous; same-dimension serialization |
| 3 | **Three guardrails** — baseline/change distinction; status-slot exclusion; provenance outside reconstruction |
| 4 | **S-3 rebuild and integrity** — lifecycle-stream-only rebuild; mismatch/duplicate integrity errors; missing-latest limitation preserved |
| 5 | **Status-exclusion verification** — tests that the slot is not treated as a valid status or folded into scope/class |

**IAR-05 (Calculated — L-2 implementation).** The implementation would add a lifecycle-side formation entry carrying
the initial `scope` and `knowledge_class`, with an explicit unresolved status component. The entry is the baseline for
the two P-1 dimensions; it is not a `scope_changed`/`class_changed` entry.

**IAR-06 (Calculated — P-1 integration).** The implementation would add the two change kinds and their independent
version sequences, the projection (baseline + highest version per dimension), and the check (mismatch/duplicate =
integrity error). Same-dimension writes serialize; cross-dimension writes do not contend.

**IAR-07 (Calculated — guardrails).**

1. **Baseline vs change:** the fold distinguishes the formation entry from change entries; the formation entry is not
   folded as a change or status event.
2. **Status slot:** excluded from the scope/class projection; never treated as `provisional`, `active` or any valid
   status. Its unresolved nature must survive **storage, reconstruction, validation and inspection**; it must not become
   a nullable status that downstream code treats as an ordinary lifecycle value.
3. **Provenance:** an optional typed reference may identify evidence, but is outside the reconstruction path; no value
   dereference, no stored join, no cross-stream chronology.

**IAR-08 (Calculated — S-3 behavior).** The rebuild reads the formation entry and the highest version per dimension; it
does not read the record, the acquisition stream or a stored join. A mismatch or duplicate is an integrity error
surfaced to `inspection`; no silent repair. The missing-latest limitation is documented, not “fixed.”

**IAR-09 (Calculated — status exclusion verification).** Implementation-stage tests must show that the status slot is
not a valid status value, is not folded into the scope/class projection, and is not exposed as `provisional`/`active`.

---

## 3. Exclusions — what the authorization would not cover

**IAR-10 (Observed — the excluded set).**

| Excluded | Why |
|---|---|
| **Decision B** — completeness assurance | separate decision; missing-latest detection remains unsolved |
| **N-b** — status value/vocabulary | open; the status slot stays unresolved |
| **N-c** — general lifecycle ordering | open; L-2 creates no universal sequence |
| **formation authority** | separate; representation does not decide who may form |
| **cross-stream ordering** | OD-a separate structures; no acquisition ordering, no shared chronology |
| **Howl state machine / UX** | separate; not designed here |
| **N-a / status-transition fold** | outside this scope; not decided |
| **§8.1 collision / O-8a / v9 domain** | separate axes; not decided |
| **requirement capture / canonical capture** | separate gates; map stays 118 |
| **deployment** | separate gate |

**IAR-11 (Calculated — no scope expansion).** If implementation needs any excluded item, it stops and returns for a
separate decision. Readiness does not justify widening the request.

---

## 4. Physical design is not selected here

**IAR-12 (Observed — no physical selection).** This request does not select a physical schema, field names, encoding,
kind names or order-key implementation. Those are implementation-plan details within the authorized scope, and they
must respect the logical constraints.

**IAR-13 (Calculated — the constraints the implementation must preserve).**

- lifecycle-side placement and identity;
- lifecycle-local ordering; no `(turn, seq)` borrowing;
- values carried in the formation entry; no record/acquisition dereference;
- explicit unresolved status slot;
- per-dimension precedence; no cross-dimension/global order;
- optional provenance outside reconstruction.

---

## 5. Acceptance criteria for the authorization

**IAR-14 (Calculated — what “done” would mean).**

| # | Criterion |
|---|---|
| AC-1 | the formation entry exists and carries initial `scope` and `knowledge_class`; the status slot is explicit |
| AC-2 | scope and class reconstruct independently from their baselines and P-1 changes |
| AC-3 | no P-1 rule introduces cross-dimension ordering |
| AC-4 | the fold distinguishes the formation entry from change entries |
| AC-5 | the status slot is excluded from scope/class projection and is never a valid status |
| AC-6 | the provenance reference is outside reconstruction; no join or chronology |
| AC-7 | S-3 rebuild from the lifecycle stream alone; mismatch/duplicate integrity errors |
| AC-8 | the missing-latest limitation is documented; no false completeness claim |
| AC-9 | no Decision B, N-b, N-c, authority, cross-stream ordering or Howl work |
| AC-10 | no requirement capture or canonical capture |
| AC-11 | the unresolved status slot survives storage, reconstruction, validation and inspection; it is not a nullable ordinary status |

**IAR-15 (Calculated — logical validity vs runtime correctness).** S-3’s pass establishes **logical validity**, not
runtime correctness or empirical performance. The implementation-stage verification must establish runtime behavior
within the scope; it must not claim empirical performance without measurement.

---

## 6. Verification plan (implementation stage, not now)

**IAR-16 (Calculated — the checks).**

| Check | What it verifies |
|---|---|
| formation baseline | scope/class reconstruct from the entry |
| P-1 scope change | scope advances; class unaffected |
| P-1 class change | class advances; scope unaffected |
| both dimensions | independent versions; both projections correct |
| guardrail 1 | the formation entry is not folded as a change |
| guardrail 2 | the status slot is not a status; excluded from projection |
| guardrail 3 | provenance is not used for values; no join |
| rebuild | lifecycle stream alone; no record/acquisition input |
| mismatch/duplicate | integrity error to `inspection` |
| missing latest | limitation documented; no false completeness |

**IAR-17 (Observed — no execution now).** This request runs no tests. The verification plan is for the implementation
stage after authorization.

---

## 7. Risks and preserved limitations

**IAR-18 (Observed — the limitations).**

| Limitation | Status |
|---|---|
| missing latest P-1 event | **undetectable by sequence**; preserved; Decision B |
| status value | **unresolved**; N-b |
| general lifecycle order | **open**; N-c |
| formation authority | **open** |
| completeness assurance | **open**; Decision B |
| runtime correctness/performance | **not established**; implementation-stage verification |

**IAR-19 (Calculated — implementation risk).** The most likely error is treating the formation entry as an ordinary
change or the status slot as a valid status — including representing it as a nullable ordinary status that downstream
code treats as a lifecycle value. The guardrails, AC-4/AC-5 and AC-11 exist to catch that.

---

## 8. What this request does not do

**IAR-20 (Observed).**

| Untouched | State |
|---|---|
| **P-1 / S-3** | closed; not reopened |
| **L-2** | adopted; not redesigned |
| **implementation** | **not authorized** |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Decision B / N-b / N-c / authority / Howl** | untouched |

---

## 9. Decision requested

**IAR-21 (Observed — the owner’s choice).**

| Option | Meaning |
|---|---|
| **Approve** | authorize implementation within IAR-04 and IAR-14, with the exclusions in IAR-10 |
| **Approve with changes** | narrow or add conditions; return the revised scope |
| **Decline** | do not authorize; the logical decisions remain closed |

**IAR-22 (Observed — approval is the authorization).** Only the owner’s approval creates implementation authority. This
document is a request, not a grant.

---

## 10. Traceability

**IAR-23 (Observed).**

| What | Where |
|---|---|
| L-2 decision | L2D-01…L2D-19 |
| P-1 decision | P1D-04…P1D-11 |
| S-3 acceptance | S3A-27 |
| Consistency review | LCR-26/27 |
| Three guardrails | LCR-27; IAR-07 |
| Missing-latest limitation | P1D-17; S3A-25; IAR-18 |
| OD-a separate structures | ECI-01; A2-45 |
| P-1/S-3 closed | P1D-01; S3A-27 |
| Map unchanged | 118 rows |

---

## 11. Record status

**IAR-24 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Authorization granted | **0** — request only |
| Implementation / tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **authorization request** — no canonical capture |

**End of request.** The owner is asked to authorize implementation of **L-2 and its P-1 integration**, the three
guardrails, S-3 rebuild/integrity behavior and status-slot exclusion verification, within IAR-04 and IAR-14
(IAR-01/14). The request explicitly excludes Decision B, N-b, N-c, formation authority, cross-stream ordering, Howl,
requirement capture and canonical capture (IAR-10). **Readiness is not authorization; no implementation, schema change,
code, test, requirement capture or canonical capture follows from this document** (IAR-03/20/24).
