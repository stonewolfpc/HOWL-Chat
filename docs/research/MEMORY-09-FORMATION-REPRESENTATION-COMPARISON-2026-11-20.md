# MEMORY-09 — Formation Representation Comparison

## A single entry with per-dimension baselines satisfies the boundary without a global order

**Status:** bounded logical comparison. **No representation selected. No kind name, final payload, order-key
implementation, schema, code, test, requirement, map change or canonical capture.** P-1 and S-3 remain closed;
Decision B remains independent.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | can a dedicated formation entry satisfy the conditions without a shared envelope, global order or acquisition dependency? |
| Finding | **yes** — a single formation entry with **per-dimension baseline components** satisfies all conditions; a single undifferentiated baseline is ambiguous about universal sequence; separate per-dimension entries split formation |
| Selection | **none** |
| Stop condition | if the relationships need a broader lifecycle ordering rule, stop and identify the **N-c** dependency |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The refined boundary

**FRM-01 (Observed — the owner’s refinement).** “The memory’s first lifecycle entry” is safe only as a
**formation-local relationship**. It must not establish that every lifecycle event belongs to one universal sequence,
and it must not settle **N-c** by implication.

**FRM-02 (Observed — what the representation may establish).** Only these four relationships:

1. a formation entry **exists**;
2. its initial values are **reconstructible from that entry**;
3. P-1 changes **follow the formation baseline for their respective dimensions**;
4. **no cross-stream ordering** is implied.

**FRM-03 (Observed — the stop condition).** If those relationships cannot be represented without a broader lifecycle
ordering rule, **stop and identify the N-c dependency** rather than extending the authorization.

**FRM-04 (Observed — the status condition).** The unresolved status component stays **explicit**. The value waits for
N-b; the obligation that initial status be stream-reconstructible remains.

---

## 2. Conditions the shape must satisfy

**FRM-05 (Calculated — the checklist).**

| # | Condition |
|---|---|
| C1 | reconstructs initial `scope`, `knowledge_class` and the status slot from the entry |
| C2 | formation is distinct from later lifecycle transitions |
| C3 | no shared envelope or shared identifier space |
| C4 | no global lifecycle order; no N-c settlement |
| C5 | no acquisition dependency for values |
| C6 | P-1 changes follow the baseline **per dimension** |
| C7 | status slot explicit; value deferred to N-b |
| C8 | lifecycle-local identity; no acquisition event id |
| C9 | `MemoryRecord` fields checkable as projections |
| C10 | no cross-stream chronology or stored join |

---

## 3. Candidate logical shapes

### 3.1 L-1 — single formation entry, single undifferentiated baseline

**FRM-06 (Calculated).** One lifecycle entry carries `scope`, `knowledge_class` and the status slot. P-1 changes are
said to “follow” it.

- **Strength:** minimal; one entry; one identity.
- **Risk:** a single “follows” relation can be read as a **universal sequence** — the entry preceding every lifecycle
  event. That is exactly the N-c implication the owner warns against.
- **Verdict:** can satisfy **only if** “follows” is explicitly defined per dimension; otherwise it is ambiguous.

### 3.2 L-2 — single formation entry, per-dimension baseline components

**FRM-07 (Calculated).** One lifecycle entry carries three logical components:

| Component | Carries | Precedence |
|---|---|---|
| `scope_baseline` | initial `scope` | the first `scope_changed` version follows it |
| `class_baseline` | initial `knowledge_class` | the first `class_changed` version follows it |
| `status_slot` | unresolved status component | value/vocabulary deferred to N-b |

The entry has **one identity**. Its components are not separate entries. The precedence is **per dimension**: scope
changes are ordered against the scope baseline; class changes against the class baseline; no relation is asserted
between scope and class changes, or between the baseline and status transitions.

- **Strength:** explicit per-dimension precedence; no universal sequence; no global order; one formation fact.
- **Verdict:** satisfies C1–C10.

### 3.3 L-3 — separate per-dimension formation entries

**FRM-08 (Calculated).** Three lifecycle entries: scope baseline, class baseline, status slot. Each has its own
identity; P-1 changes follow the corresponding entry.

- **Strength:** per-dimension precedence is explicit; no single universal baseline.
- **Risk:** “a formation entry exists” becomes **three formation entries**; it may need a formation-group identity to
  tie them, adding identity and ordering complexity. It also raises whether formation is one fact or three.
- **Verdict:** viable, but splits formation and is less faithful to the singular entry the direction adopted.

### 3.4 L-4 — L-2 plus a typed provenance reference

**FRM-09 (Calculated).** L-2 plus an optional typed reference to the formation evidence/acquisition event.

- **Strength:** provenance; OD-a permits typed cross-stream references.
- **Condition:** the reference must be **provenance only**; the lifecycle rebuild must not dereference it for the
  initial values.
- **Verdict:** an optional variant of L-2, not a separate shape.

### 3.5 Excluded shapes

**FRM-10 (Calculated).**

| Shape | Why excluded |
|---|---|
| **first P-1 change carries the baseline** | conflates formation with a change; no baseline if no change occurs |
| **global lifecycle order** | settles N-c; exceeds the authorization |
| **shared envelope / shared id** | fails OD-a’s separate structures |
| **external checkpoint/manifest** | not stream-reconstructible; external baseline authority |
| **record/acquisition dereference for values** | fails Decision A and K-R14 |

---

## 4. Comparison

**FRM-11 (Calculated — the table).**

| Condition | L-1 single baseline | L-2 per-dimension components | L-3 separate entries | L-4 L-2 + typed ref |
|---|---|---|---|---|
| **C1** scope/class/status slot | yes | yes | yes | yes |
| **C2** formation distinct | yes | yes | **tension** — three formation facts | yes |
| **C3** no shared envelope | yes | yes | yes | yes |
| **C4** no global order / N-c | **ambiguous** | **yes** — per-dimension | yes | yes |
| **C5** no acquisition dependency | yes | yes | yes | **conditional** — provenance only |
| **C6** per-dimension precedence | **implicit** | **explicit** | explicit | explicit |
| **C7** status slot | yes | yes | yes | yes |
| **C8** lifecycle-local identity | yes | yes | **multiple** identities; group needed | yes |
| **C9** record fields checkable | yes | yes | yes | yes |
| **C10** no cross-stream chronology/join | yes | yes | yes | **conditional** — ref only |

**FRM-12 (Calculated — L-2 is the least ambiguous).** L-2 makes the per-dimension precedence explicit without
splitting formation. It satisfies C1–C10. L-1 is the same shape with the precedence left implicit, which is where the
universal-sequence reading enters. L-3 is viable but creates multiple formation entries and a grouping problem. L-4 is
L-2 plus an optional provenance reference that must never supply values.

**FRM-13 (Calculated — L-1 and L-2 are not different ontologies).** They differ in **explicitness**, not in the facts
they record: both have one formation entry with scope, class and the status slot. L-2 states the per-dimension
precedence; L-1 leaves it to inference. The explicit form is safer against the N-c implication.

---

## 5. How L-2 satisfies the four relationships

**FRM-14 (Calculated — relationship 1).** A formation entry exists: one lifecycle-local entry per memory. It is not a
later lifecycle transition.

**FRM-15 (Calculated — relationship 2).** Its initial values are reconstructible from the entry: `scope_baseline` and
`class_baseline` carry values; `status_slot` carries the unresolved component. No record or acquisition dereference is
needed.

**FRM-16 (Calculated — relationship 3).** P-1 changes follow the baseline **for their respective dimensions**: the first
`scope_changed` follows `scope_baseline`; the first `class_changed` follows `class_baseline`. No relation is asserted
between the two dimensions.

**FRM-17 (Calculated — relationship 4).** No cross-stream ordering is implied: the entry has a lifecycle-local
identity; a typed provenance reference, if present, carries no chronology.

**FRM-18 (Calculated — no global order).** L-2 does not define a lifecycle sequence. It defines **two local
precedence relations** (scope and class) plus an unresolved status component. It does not order status transitions,
scope against class, or the entry against any acquisition event. N-c remains open.

**FRM-19 (Calculated — the status slot is not a status).** The slot is a representation component, not `provisional` or
`active`. It preserves the obligation; N-b supplies the value. The record’s status cannot be semantically checked until
then; structural presence is checkable now.

**FRM-20 (Calculated — the record fields).** `MemoryRecord.scope` and `knowledge_class` are checked against the
formation entry plus P-1 changes. `MemoryRecord.status` is checked structurally against the slot; semantic checking
waits for N-b. No record field is the reconstruction authority.

---

## 6. Stop conditions

**FRM-21 (Calculated — stop if a broader order is needed).** If the design cannot express “the first scope change
follows the scope baseline” and “the first class change follows the class baseline” without a **general lifecycle
order**, stop. That is an **N-c** dependency; do not extend the authorization to cover it.

**FRM-22 (Calculated — stop if acquisition ordering is borrowed).** If the design uses K-R15’s `(turn, seq)` as the
lifecycle order, stop. That is a cross-stream ordering rule OD-a does not establish.

**FRM-23 (Calculated — stop if a shared envelope is needed).** If the design needs one identifier space or a shared
envelope, stop. OD-a Option B is separate structures.

**FRM-24 (Calculated — stop if values are dereferenced).** If the design needs the record or the acquisition event to
supply the initial values, stop. That is the external baseline Decision A rejects.

**FRM-25 (Calculated — stop if status is forced).** If the design must choose `provisional` or `active` before N-b,
stop. The status slot is the explicit deferral.

---

## 7. What remains open

**FRM-26 (Calculated).**

| Open | State |
|---|---|
| **representation shape** | L-2 is the least ambiguous; **not selected** |
| **kind name** | not selected |
| **final payload** | not selected; at least scope, class, status slot |
| **order-key implementation** | not selected; per-dimension local precedence is the constraint |
| **N-b / status value** | open; slot explicit |
| **N-c / general lifecycle order** | open; L-2 avoids it |
| **formation authority** | open; orthogonal |
| **Decision B** | untouched |
| **schema / code / tests / requirements / canonical** | not authorized |

---

## 8. What this record does not do

**FRM-27 (Observed).**

| Untouched | State |
|---|---|
| **P-1 / S-3** | closed; not reopened |
| **Decision A** | resolved Yes; not reopened |
| **formation representation** | compared; not selected |
| **N-b / N-c / §8.1** | not decided |
| **Decision B** | untouched |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |

---

## 9. Traceability

**FRM-28 (Observed).**

| What | Where |
|---|---|
| Formation-local, not universal | FRM-01/02/03; FRC-13 |
| OD-a Option B; separate identity/ordering; typed refs | ECI-01; A2-45 |
| No lifecycle schema / fold rule | A2-11/A2-12; A2-41 |
| `provisional`/`established`/`active` | A2-42/A2-43; NAE-12/25 |
| Two streams, two checks | K-R11–K-R14; §6.4 |
| Acquisition ordering `(turn, seq)` | K-R15 |
| Copied turn ≠ cross-stream ordering | ECI-34/35 |
| Entry must be lifecycle-side | ECI-17 |
| Decision A owner determination | DAO-01…DAO-14 |
| Existing-entry audit | EPA-11…EPA-31 |
| Coexistence check | FRC-10…FRC-25 |
| Map unchanged | 118 rows |

---

## 10. Record status

**FRM-29 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Representation selected | **0** — L-2 is the least ambiguous, not adopted |
| Kind/payload/order selected | **0** |
| N-b / N-c / §8.1 decided | **0** |
| Implementation / tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **logical comparison** — no canonical capture |

**End of comparison.** A dedicated formation entry **can** satisfy the conditions without a shared envelope, global
order or acquisition dependency: **L-2** — one entry with per-dimension baseline components — satisfies C1–C10
(FRM-07/11/12). L-1 is the same facts with the per-dimension precedence left implicit, which risks the universal-
sequence reading; L-3 splits formation; L-4 is an optional provenance reference that must never supply values
(FRM-06…FRM-13). The representation establishes only the four relationships in FRM-02; it does **not** order status
transitions or settle N-c (FRM-18). If a broader lifecycle order is needed, **stop and identify N-c** (FRM-21). **No
representation, kind, payload or order key is selected; P-1 and S-3 remain closed; Decision B is untouched; no schema,
code, requirement, map change or canonical capture follows** (FRM-26/27/29).
