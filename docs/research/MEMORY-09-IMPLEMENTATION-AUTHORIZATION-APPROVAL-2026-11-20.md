# MEMORY-09 — Implementation Authorization Approval

## Approved within the IAR scope, with the unresolved-status safeguard

**Status:** **owner approval recorded.** Implementation is authorized **within the scope of IAR-04 and IAR-14**, with
the safeguard below. **No implementation, code, schema, migration, test, requirement capture or canonical capture has
been started by this record.** P-1 and S-3 remain closed; Decision B remains untouched.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Decision | **Approved** |
| Scope | L-2 formation representation and P-1 integration; three guardrails; S-3 rebuild/integrity; status-slot exclusion verification |
| Safeguard | the unresolved status must survive **storage, reconstruction, validation and inspection**; it is not a nullable ordinary status |
| Exclusions | Decision B, N-b, N-c, formation authority, cross-stream ordering, Howl state machine, requirement/canonical capture, deployment |
| P-1 / S-3 | **closed**; not reopened |
| Map | **118** rows |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The approval

**IAA-01 (Observed — the owner’s determination).** The owner **approves** the implementation-authorization request
(MEMORY-09-IMPLEMENTATION-AUTHORIZATION-REQUEST-2026-11-20.md) within its scope and exclusions.

**IAA-02 (Observed — the authorized scope).** Implementation is authorized for:

1. the **L-2 formation representation** — one lifecycle-side entry with `scope_baseline`, `class_baseline`,
   `status_slot`;
2. its integration with the **two P-1 change dimensions** — `scope_changed` / `class_changed`, independent versions,
   same-dimension serialization;
3. the **three consistency guardrails** — baseline/change distinction; status-slot exclusion; provenance outside
   reconstruction;
4. **S-3 rebuild and integrity behavior** — lifecycle-stream-only rebuild; mismatch/duplicate integrity errors;
   missing-latest limitation preserved;
5. **status-slot exclusion verification**.

**IAA-03 (Observed — the safeguard).** The unresolved `status_slot` must not be silently turned into a **nullable
status** that downstream code treats as an ordinary lifecycle value. Its unresolved nature must survive:

| Stage | Requirement |
|---|---|
| **storage** | the slot is stored as unresolved; it is not collapsed to null-as-status |
| **reconstruction** | reconstruction preserves the unresolved component; it does not infer `provisional`/`active` |
| **validation** | validation rejects treating the slot as a valid status value |
| **inspection** | inspection shows the slot as unresolved, not as an ordinary status |

**IAA-04 (Observed — acceptance).** The IAR acceptance criteria AC-1…AC-11 apply, including AC-11: the unresolved
status survives storage, reconstruction, validation and inspection and is not a nullable ordinary status.

---

## 2. What the approval authorizes

**IAA-05 (Calculated — authorized work).** Within the scope, implementation may include the minimum schema, code,
migrations and tests needed to realize L-2 and P-1, subject to the logical constraints and acceptance criteria. The
physical design — field names, encoding, kind names and order-key implementation — is an implementation-plan detail
and must satisfy the constraints; this approval does not select a particular one.

**IAA-06 (Calculated — required behavior).**

- the formation entry is the baseline, not a `scope_changed`/`class_changed` entry;
- the fold distinguishes the formation entry from change entries;
- scope and class reconstruct independently; no cross-dimension ordering;
- the status slot is excluded from the scope/class projection and is never a valid status;
- provenance is outside reconstruction; no value dereference, no stored join, no cross-stream chronology;
- S-3 rebuild uses the lifecycle stream alone; mismatch/duplicate are integrity errors;
- the missing-latest limitation is documented, not “fixed.”

---

## 3. What the approval does not authorize

**IAA-07 (Observed — the exclusions).**

| Excluded | State |
|---|---|
| **Decision B** — completeness assurance | separate; untouched |
| **N-b** — status value/vocabulary | open; the slot stays unresolved |
| **N-c** — general lifecycle ordering | open; no universal sequence |
| **formation authority** | open; separate |
| **cross-stream ordering** | excluded; OD-a separate structures |
| **Howl state machine / UX** | not designed |
| **N-a / status-transition fold** | outside scope |
| **§8.1 / O-8a / v9 domain** | outside scope |
| **requirement capture / canonical capture** | not authorized; map stays 118 |
| **deployment** | separate gate |

**IAA-08 (Calculated — stop conditions remain).** If implementation needs any excluded item, it stops and returns to the
owner. Approval does not widen the scope.

---

## 4. Implementation status

**IAA-09 (Observed — not started).** This approval records authority; it does not itself create code, schema,
migrations, tests, requirements or canonical text. The implementation plan and physical design are the next work under
the authorization.

**IAA-10 (Calculated — the next step).** Prepare the implementation plan within IAR-04/IAR-14, making the status-slot
safeguard concrete across storage, reconstruction, validation and inspection. Any excluded decision returns to the
owner.

---

## 5. Preserved boundaries

**IAA-11 (Observed — P-1 and S-3).** P-1 and S-3 remain closed; the approval implements them, it does not reopen them.

**IAA-12 (Observed — logical validity vs runtime).** S-3’s pass is logical validity, not runtime correctness or
empirical performance. Implementation-stage verification must establish runtime behavior within the scope; no
performance claim without measurement.

**IAA-13 (Observed — no requirement/canonical).** The approval adds no requirement and changes no canonical text. The
map stays **118**.

---

## 6. Traceability

**IAA-14 (Observed).**

| What | Where |
|---|---|
| Authorization request | IAR-01…IAR-24 |
| Scope and acceptance | IAR-04; IAR-14 (AC-1…AC-11) |
| Three guardrails | IAR-07; LCR-27 |
| Status-slot safeguard | IAA-03; AC-11 |
| Exclusions | IAR-10; IAA-07 |
| L-2 decision | L2D-01…L2D-19 |
| P-1 decision | P1D-04…P1D-11 |
| S-3 acceptance | S3A-27 |
| Consistency review | LCR-26/27 |
| Map unchanged | 118 rows |

---

## 7. Record status

**IAA-15 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Authorization | **approved** within IAR scope |
| Implementation started | **0** |
| Tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **owner approval** — no canonical capture |

**End of approval.** The owner approves implementation within the IAR scope (IAA-01/02). The unresolved `status_slot`
must survive **storage, reconstruction, validation and inspection** and must not become a nullable ordinary status
(IAA-03; AC-11). Decision B, N-b, N-c, formation authority, cross-stream ordering, Howl, requirement capture and
canonical capture remain excluded (IAA-07). **No code, schema, migration, test, requirement or canonical change has
been started by this record** (IAA-09/13/15). The next step is the implementation plan under the authorization.
