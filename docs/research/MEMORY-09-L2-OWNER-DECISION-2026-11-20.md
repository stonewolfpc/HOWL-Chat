# MEMORY-09 — L-2 Owner Decision

## One formation entry with independent per-dimension baselines

**Status:** **owner decision recorded** at the logical-representation level. No physical schema, field encoding, kind
name, ordering implementation, code, test, requirement, map change or canonical capture. P-1 and S-3 remain closed;
Decision B remains independent.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Decision | **L-2 adopted** as the logical representation of formation |
| Shape | one formation entry with `scope_baseline`, `class_baseline`, `status_slot` |
| Conditions | independent per-dimension precedence; status slot unresolved until N-b; no universal/cross-dimension/cross-stream order |
| Not selected | physical schema, field encoding, kind name, ordering implementation |
| P-1 / S-3 | **closed**; not reopened |
| Decision B | **untouched** |
| Map | **118** rows |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The decision

**L2D-01 (Observed — the owner’s decision).** The owner adopts **L-2** as the logical representation: a single
formation entry carrying two baseline components and an unresolved status component.

**L2D-02 (Observed — the level).** This is a **logical** decision. It does not select a physical schema, field
encoding, kind name, ordering implementation, or any code.

**L2D-03 (Calculated — why L-2).** L-2 records the same facts as L-1 but makes the **two independent per-dimension
precedence relationships explicit**. That explicitness is the architectural benefit: each dimension reconstructs
without inventing an ordering relationship between dimensions.

---

## 2. The three conditions

**L2D-04 (Observed — condition 1).** The two baseline components have **independent local precedence** over their
corresponding P-1 changes:

- `scope_baseline` → that memory’s scope changes;
- `class_baseline` → that memory’s class changes.

Neither baseline orders the other dimension’s changes.

**L2D-05 (Observed — condition 2).** The `status_slot` remains **unresolved until N-b**, without weakening Decision A.
The obligation that initial status be stream-reconstructible remains; only the value/vocabulary is pending.

**L2D-06 (Observed — condition 3).** No **universal lifecycle sequence**, **cross-dimension ordering**, or
**cross-stream chronology** is created.

---

## 3. What L-2 establishes

**L2D-07 (Calculated — one entry).** A formation entry exists per memory. It is lifecycle-side; it is not a later
lifecycle transition.

**L2D-08 (Calculated — the components).**

| Component | Carries | Local precedence |
|---|---|---|
| `scope_baseline` | initial `scope` | the first `scope_changed` version follows it |
| `class_baseline` | initial `knowledge_class` | the first `class_changed` version follows it |
| `status_slot` | unresolved status component | value/vocabulary deferred to N-b |

The component labels are **logical**, not physical field names.

**L2D-09 (Calculated — scope reconstruction).** `scope_baseline` → that memory’s scope changes. A rebuild reads the
baseline, then the scope version sequence. It does not read class or status.

**L2D-10 (Calculated — class reconstruction).** `class_baseline` → that memory’s class changes. A rebuild reads the
baseline, then the class version sequence. It does not read scope or status.

**L2D-11 (Calculated — status reconstruction).** `status_slot` is present and explicit. The concrete value/vocabulary
waits for N-b. The stream preserves the obligation; it does not yet reconstruct a concrete status.

**L2D-12 (Calculated — acquisition separation).** Acquisition history remains separate, with **no ordering dependency**.
The formation entry carries its own values; it does not require an acquisition event for reconstruction.

**L2D-13 (Calculated — optional provenance).** A typed provenance reference may identify supporting evidence, but it
sits **outside the reconstruction path**. The formation entry contains everything needed to reconstruct its values.

**L2D-14 (Calculated — record fields).** `MemoryRecord.scope` and `knowledge_class` are checkable projections against
the baseline plus the corresponding P-1 changes. `MemoryRecord.status` is structurally checkable against the slot;
semantic checking waits for N-b. No record field is the reconstruction authority.

---

## 4. What L-2 does not create

**L2D-15 (Calculated — no cross-dimension ordering).** A scope change is not ordered against a class change. The two
precedence relations are independent.

**L2D-16 (Calculated — no universal lifecycle sequence).** The entry is not “first” in a global lifecycle order. It is
the local baseline for its own dimensions. N-c remains open.

**L2D-17 (Calculated — no status ordering).** Status transitions are not ordered by `scope_baseline` or
`class_baseline`. Their semantics and ordering remain N-a/N-b/N-c.

**L2D-18 (Calculated — no cross-stream chronology).** Nothing compares a lifecycle position with an acquisition
position. A copied `established_turn` may be a payload value, not a cross-stream ordering rule.

**L2D-19 (Calculated — no join).** No stored cross-stream join is created; acquisition and lifecycle remain separate
structures (OD-a Option B).

---

## 5. What is not selected

**L2D-20 (Observed — the non-selections).**

| Not selected | Why |
|---|---|
| **physical schema** | separate implementation/design gate |
| **field encoding** | separate design choice |
| **kind name** | separate design choice |
| **ordering implementation** | separate design choice |
| **N-b status value** | open; the slot is explicit |
| **N-c general lifecycle order** | open; L-2 avoids it |
| **§8.1 collision** | open; not decided by L-2 |
| **formation authority** | open; orthogonal |
| **completeness marker** | Decision B; untouched |

**L2D-21 (Calculated — no implementation authorization).** L-2 is a logical representation. It does not authorize
schema, code, tests, migrations, requirement capture or canonical capture.

---

## 6. What remains open

**L2D-22 (Calculated).**

| Open | State |
|---|---|
| **kind name / final payload / encoding** | open; separate design |
| **ordering implementation** | open; per-dimension local precedence is the constraint |
| **N-b** | open; status value pending |
| **N-c** | open; no general lifecycle order |
| **§8.1** | open |
| **formation authority** | open |
| **Decision B** | untouched |
| **schema / code / tests / requirements / canonical** | not authorized |

---

## 7. What this record does not do

**L2D-23 (Observed).**

| Untouched | State |
|---|---|
| **P-1 / S-3** | closed; not reopened |
| **Decision A** | resolved Yes; not reopened |
| **physical design** | not selected |
| **N-b / N-c / §8.1** | not decided |
| **Decision B** | untouched |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |

---

## 8. Traceability

**L2D-24 (Observed).**

| What | Where |
|---|---|
| L-2 comparison | FRM-07/11/12/14…FRM-20 |
| Formation-local boundary | FRM-01/02/03; FRC-13 |
| OD-a Option B; separate structures | ECI-01; A2-45 |
| Two streams, two checks | K-R11–K-R14; §6.4 |
| Acquisition ordering `(turn, seq)` | K-R15 |
| Copied turn ≠ cross-stream ordering | ECI-34/35 |
| Entry must be lifecycle-side | ECI-17 |
| Decision A owner determination | DAO-01…DAO-14 |
| Existing-entry audit | EPA-11…EPA-31 |
| Coexistence check | FRC-10…FRC-25 |
| Map unchanged | 118 rows |

---

## 9. Record status

**L2D-25 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Logical representation | **L-2 adopted** |
| Physical design selected | **0** |
| N-b / N-c / §8.1 decided | **0** |
| Implementation / tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **owner decision** — no canonical capture |

**End of decision.** The owner adopts **L-2** as the logical formation representation: one entry with
`scope_baseline`, `class_baseline` and `status_slot` (L2D-01/07/08). The two baselines have **independent local
precedence** over their corresponding P-1 changes; the status slot stays unresolved until N-b without weakening
Decision A; no universal lifecycle sequence, cross-dimension ordering or cross-stream chronology is created
(L2D-04…L2D-06/15…L2D-19). No physical schema, field encoding, kind name or ordering implementation is selected
(L2D-20/21). **P-1 and S-3 remain closed; Decision B is untouched; no schema, code, requirement, map change or
canonical capture follows** (L2D-22/23/25).
