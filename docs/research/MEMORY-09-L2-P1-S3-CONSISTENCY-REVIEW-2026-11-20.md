# MEMORY-09 — L-2 / P-1 / S-3 Consistency Review

## L-2 is consistent with P-1 and S-3, with three explicit guardrails

**Status:** read-only logical consistency review. **No design round, no schema, code, migration, test, requirement,
map change or canonical capture.** P-1 and S-3 remain closed; Decision B remains independent.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Objects | L-2 formation representation; P-1 post-formation changes; S-3 acceptance |
| Result | **consistent**, with three guardrails: baseline-vs-change distinction; status-slot exclusion; provenance outside reconstruction |
| Not done | no new design, no kind/payload/order selection, no implementation |
| P-1 / S-3 | **closed**; not reopened |
| Decision B | **untouched** |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Task and boundary

**LCR-01 (Observed — the task).** Final logical consistency review of L-2 against P-1 and S-3 before any
implementation-authorization request. Read-only; not another design round.

**LCR-02 (Observed — the checks).**

1. both P-1 dimensions start from the corresponding L-2 baseline;
2. no P-1 rule introduces cross-dimension ordering;
3. S-3’s rebuild and integrity checks remain valid with a formation entry as the baseline source;
4. the unresolved status slot is not treated as a valid status value;
5. formation authority, N-b, N-c and completeness assurance remain outside.

**LCR-03 (Observed — the boundary).** No schema, code, tests, requirements, map change or canonical capture. P-1 and
S-3 remain closed; Decision B is untouched.

---

## 2. The objects

**LCR-04 (Observed — L-2).** One formation entry with `scope_baseline`, `class_baseline` and `status_slot`; independent
local precedence per dimension; no universal/cross-dimension/cross-stream order (L2D-01…L2D-19).

**LCR-05 (Observed — P-1).** `scope_changed` / `class_changed`; independent `scope_version` / `class_version`;
formation value baseline; first change version 1; contiguous; same-dimension serialized; cross-dimension independent;
projection from the lifecycle stream (P1D-04…P1D-11).

**LCR-06 (Observed — S-3).** P-1 passes K-R12, K-R14, §6.4, K-R13, K-R15’s local-order boundary, K-R2’s integrity
principle and S2H-12, with the missing-latest limitation and open completeness/recovery/concurrency questions preserved
(S3A-27).

---

## 3. Check 1 — both P-1 dimensions start from the L-2 baseline

**LCR-07 (Calculated — the mapping).**

| P-1 element | L-2 element | Consistency |
|---|---|---|
| `scope` formation baseline (P1D-06) | `scope_baseline` | **match** |
| `class` formation baseline (P1D-06) | `class_baseline` | **match** |
| projection starts from the formation value (P1D-09) | baseline component is the start | **match** |
| first change is version 1 (P1D-07) | baseline precedes version 1 locally | **match** |
| no change entry for the baseline (P1D-06) | the formation entry is **not** a `scope_changed`/`class_changed` entry | **match, with clarification** |

**LCR-08 (Calculated — the reconstruction paths).**

- **Scope:** `scope_baseline` → scope changes; no class or status input.
- **Class:** `class_baseline` → class changes; no scope or status input.

P-1’s per-dimension projection and L-2’s per-dimension components are the same paths. **Pass.**

**LCR-09 (Observed — the clarification).** P1D-06 says the formation value “has no change entry.” That means **no
`scope_changed`/`class_changed` entry**, not “no lifecycle entry.” L-2 adds a formation entry; it is the baseline, not
a P-1 change. The intended reading is the former.

**LCR-10 (Calculated — the guardrail).** The P-1 fold must distinguish the **formation entry** from **change entries**.
A generic fold over “lifecycle entries” must not treat the formation entry as a change or a status event. The fold’s
change domain is `scope_changed`/`class_changed`; the formation entry supplies the baseline.

**LCR-11 (Calculated — S3A-09 is resolved, not contradicted).** S3A-09 said: “If formation is not itself a lifecycle
entry, the baseline is the record’s initial condition…” L-2 makes formation a lifecycle entry, so the conditional’s
“if” is now false. S-3’s pass remains; the baseline source is explicit rather than a record field. **Pass, strengthened.**

---

## 4. Check 2 — no cross-dimension ordering

**LCR-12 (Calculated — P-1’s rules are per dimension).**

- P1D-05: a scope change advances only `scope_version`; a class change advances only `class_version`.
- P1D-08: same-dimension writes serialize; cross-dimension writes are independent.
- P1D-13: the two histories remain independent.
- P1D-27: the version sequences order only the change entries; status ordering remains separate.
- S3A-19: N-c remains open.
- S3A-22: `class_changed` is separate from the status kinds.

**LCR-13 (Calculated — L-2 does not add one).** L-2’s two precedence relations are local to their dimensions. It
asserts no relation between a scope change and a class change, and no universal sequence. **Pass.**

**LCR-14 (Calculated — no accidental cross-dimension rule).** The only shared element is the formation entry’s
identity; identity is not an ordering relation. The per-dimension components carry the precedence. **Pass.**

---

## 5. Check 3 — S-3 rebuild and integrity remain valid

**LCR-15 (Calculated — the criteria).**

| S-3 criterion | With L-2 | Result |
|---|---|---|
| K-R12 locus | formation entry and P-1 changes are lifecycle-side | **pass** |
| K-R12 domain | scope and class each have a baseline and change kind | **pass** |
| K-R12 projection | baseline + highest version per dimension | **pass** |
| K-R12 non-mutation | record fields remain caches | **pass** |
| K-R14 independent rebuild | lifecycle stream alone: formation entry + changes | **pass** |
| K-R14 no stored join | optional provenance ref is not a join | **pass with guardrail** |
| §6.4 deterministic check | fold change entries, compare to cache | **pass** |
| §6.4 mismatch/duplicate | integrity error | **pass** |
| K-R13 no status on edge | unaffected | **pass** |
| K-R15 local order | versions order only change entries | **pass** |
| K-R2 integrity principle | unreconstructable projection is failure | **pass** |
| S2H-12 context/status | unaffected | **pass** |

**LCR-16 (Calculated — the rebuild).** A rebuild reads the formation entry for the baseline, then the highest
version per dimension. It does not read the record, the acquisition stream or a stored join. This is exactly S3A-10’s
“baseline + highest version,” with the baseline now a lifecycle entry. **Pass.**

**LCR-17 (Calculated — the optional provenance reference).** If present, it may identify supporting evidence but must
not be dereferenced for values and must not create a stored cross-stream join. It sits outside the reconstruction
path. **Pass with guardrail.**

**LCR-18 (Calculated — S-3’s verdict stands).** S3A-27 passes; L-2 supplies the baseline S3A-09 left conditional. No
re-acceptance is needed; the acceptance is strengthened, not reopened.

---

## 6. Check 4 — the status slot is not a status value

**LCR-19 (Calculated — the distinction).** `status_slot` is an **unresolved representation component**, not
`provisional`, `active`, or any other status value (L2D-05/11). It is not a status transition and does not order
status events.

**LCR-20 (Calculated — P-1 does not use it).** P-1’s projection is over `scope_changed`/`class_changed`; the status slot
is outside that domain. P-1 does not fold it as a status. **Pass.**

**LCR-21 (Calculated — S-3 does not use it).** S-3’s criteria are scope/class; K-R13 keeps status off the edge;
S3A-19 keeps N-c open. The slot is lifecycle-side but outside the P-1 projection. **Pass.**

**LCR-22 (Calculated — the guardrail).** Any consumer of the formation entry must treat `status_slot` as **unresolved**
— not as a valid status, not as `provisional`/`active`, not as a status transition. The slot preserves Decision A’s
obligation; N-b supplies the value. **Pass with guardrail.**

---

## 7. Check 5 — outside decisions remain outside

**LCR-23 (Calculated — the boundary).**

| Outside item | Evidence | Result |
|---|---|---|
| **formation authority** | L2D-20/22; P1D-24; S3A-21 | **not decided** |
| **N-b status value** | L2D-05/11; P1D-24 | **not decided** |
| **N-c general lifecycle order** | L2D-06/16; P1D-26; S3A-19 | **not decided** |
| **completeness assurance** | P1D-17; S3A-25; Decision B | **not decided / untouched** |
| **§8.1 collision** | L2D-20; P1D-24 | **not decided** |
| **implementation / schema / canon** | L2D-21; P1D-25; S3A-29 | **not authorized** |

**LCR-24 (Calculated — no accidental capture).** L-2 does not say who may form; the status slot does not select a
status; the per-dimension precedence does not create a general order; no completeness marker is added. **Pass.**

**LCR-25 (Observed — the missing-latest limitation remains).** L-2 does not detect a missing latest P-1 event; that
limitation is preserved (P1D-17; S3A-25). Completeness assurance remains Decision B.

---

## 8. Verdict

**LCR-26 (Calculated — consistent).** L-2 is consistent with P-1 and S-3:

- both P-1 dimensions start from the corresponding baseline;
- no P-1 rule introduces cross-dimension ordering;
- S-3’s rebuild and integrity checks remain valid, with the formation entry as the baseline source;
- the status slot is not a valid status value;
- authority, N-b, N-c and completeness remain outside.

**LCR-27 (Observed — the three guardrails).**

1. **Baseline vs change:** P1D-06’s “no change entry” means no `scope_changed`/`class_changed` entry; the formation
   entry is the baseline, and the fold must distinguish it from change entries.
2. **Status slot:** exclude it from the scope/class projection; never treat it as a valid status.
3. **Provenance:** the optional reference is outside reconstruction; no value dereference, no stored join, no
   cross-stream chronology.

**LCR-28 (Calculated — what the pass means).** The logical representation is internally consistent and ready for a
**separately scoped implementation-authorization request**. This review is not that request and authorizes nothing.

---

## 9. What remains open

**LCR-29 (Calculated).**

| Open | State |
|---|---|
| **implementation authorization** | separate; not requested here |
| **kind name / payload / encoding** | separate design |
| **ordering implementation** | separate design |
| **N-b / N-c / §8.1** | open |
| **formation authority** | open |
| **completeness assurance** | Decision B; untouched |
| **missing-latest limitation** | preserved |
| **schema / code / tests / requirements / canonical** | not authorized |

---

## 10. What this record does not do

**LCR-30 (Observed).**

| Untouched | State |
|---|---|
| **P-1 / S-3** | closed; not reopened |
| **L-2** | adopted; not redesigned |
| **implementation** | not authorized |
| **N-b / N-c / §8.1** | not decided |
| **Decision B** | untouched |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |

---

## 11. Traceability

**LCR-31 (Observed).**

| What | Where |
|---|---|
| L-2 decision | L2D-01…L2D-19 |
| P-1 baseline/projection | P1D-06…P1D-11 |
| P-1 independence | P1D-05/08/13 |
| P-1 limitations | P1D-16…P1D-21 |
| S-3 criteria | S3A-04 |
| S-3 formation baseline | S3A-09 |
| S-3 rebuild | S3A-10/14 |
| S-3 verdict | S3A-27 |
| K-R12/K-R14/§6.4/K-R13/K-R15/K-R2 | r33 `:111`–`:113`, `:366`–`:367`, `:83` |
| S2H-12 | addendum |
| Map unchanged | 118 rows |

---

## 12. Record status

**LCR-32 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Consistency result | **passes**, with three guardrails |
| Implementation authorized | **0** |
| N-b / N-c / §8.1 decided | **0** |
| Implementation / tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **consistency review** — no canonical capture |

**End of review.** L-2 is **consistent** with P-1 and S-3 (LCR-26). Both P-1 dimensions start from the corresponding
baseline; no P-1 rule introduces cross-dimension ordering; S-3’s rebuild and integrity checks remain valid; the status
slot is not a valid status value; authority, N-b, N-c and completeness remain outside (LCR-07…LCR-25). The three
guardrails are the baseline/change distinction, the status-slot exclusion, and provenance outside reconstruction
(LCR-27). The logical representation is ready for a **separately scoped implementation-authorization request**, which
this review does not draft or authorize (LCR-28/30).
