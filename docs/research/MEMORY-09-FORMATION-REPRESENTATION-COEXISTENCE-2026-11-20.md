# MEMORY-09 — Formation Representation Coexistence

## The dedicated direction fits OD-a’s separate structures, if it stays lifecycle-local

**Status:** direction adopted by the owner; **coexistence check only**. No kind name, final payload, order-key
implementation, schema, code, test, requirement, map change or canonical capture. Decision B remains independent;
P-1 and S-3 remain closed.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Direction | **dedicated formation representation**, adopted as the preferred design |
| Check | can it coexist with the separate acquisition and lifecycle structures established by OD-a and A2? |
| Finding | **yes, conditionally** — lifecycle-side placement, lifecycle-local identity/order, values carried in the entry, no cross-stream dependency, explicit unresolved status slot |
| Not selected | kind name, final payload, order-key implementation |
| P-1 / S-3 | **closed**; not reopened |
| Decision B | **untouched** |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Owner direction and boundaries

**FRC-01 (Observed — the adopted direction).** The owner adopts the **dedicated formation representation** as the
preferred design, subject to the boundaries below. This is a **direction**, not a kind, payload or order-key
selection.

**FRC-02 (Observed — the three values stay distinct).**

| Value | Requirement |
|---|---|
| **`scope`** | **must be carried** by the formation representation; it cannot be recovered from `MemoryRecord` alone |
| **`knowledge_class`** | **must be carried** by the formation representation |
| **`status`** | the stream must **ultimately** reconstruct it under Decision A; its initial value **cannot be selected until N-b** resolves the status interpretation |

**FRC-03 (Observed — the deferral limit).** Deferring the status **value** to N-b must not defer or weaken the
already-adopted obligation that initial status be **stream-reconstructible**. The slot and the obligation remain; only
the value is pending.

**FRC-04 (Observed — the boundaries).**

1. record formation as **distinct from later lifecycle transitions**;
2. carry `scope` and `knowledge_class` in the formation representation;
3. preserve an **explicit unresolved status component** pending N-b;
4. establish identity and deterministic ordering **before the first P-1 change**, without borrowing an unapproved
   cross-stream ordering rule;
5. treat `MemoryRecord` fields as **checkable projections**, never as reconstruction authority.

**FRC-05 (Observed — not selected).** No kind name, final payload or order-key implementation is selected here.

---

## 2. The architecture it must coexist with

**FRC-06 (Source-derived — OD-a, settled Option B).** “Acquisition events and memory-lifecycle/status events remain
**separate structures**.” “Each structure owns its **own event identity and occurrence-ordering responsibility**,
consistent with A2.” “Cross-stream references, where needed, must be **typed** rather than implying shared identity or
chronology.” “Separate structures do **not** establish cross-stream chronological ordering” (ECI-01; A2-45; OD-a
decision record).

**FRC-07 (Source-derived — A2).** The lifecycle stream is authoritative by declaration and **unspecified by schema**
(A2-11/A2-12); no fold rule exists (A2-41); the status vocabularies do not reconcile (A2-42); `provisional` has a
producer but no typed event, and `established`/`active` have no stated mapping (A2-43).

**FRC-08 (Source-derived — the two streams).** K-R11: the edge fold depends only on acquisition events. K-R12: memory
lifecycle is its own authoritative stream; status/class/scope are its projection. K-R13: the edge carries no truth
status. K-R14: each stream is rebuilt and checked independently; no stored join. §6.4: two streams, two checks, neither
projection consults the other.

**FRC-09 (Source-derived — acquisition ordering).** K-R15's `(turn, seq)` order is the **acquisition** stream’s ordering
rule. It is not a lifecycle order key; the lifecycle stream has no order key (A2-41; DS-09).

---

## 3. Coexistence checks

**FRC-10 (Calculated — the table).**

| Test | Requirement | Dedicated formation direction | Result |
|---|---|---|---|
| **Placement** | lifecycle side, not acquisition | a dedicated formation entry on the lifecycle structure | **pass** |
| **Identity** | each structure owns its own identity | lifecycle-local entry id; typed cross-ref allowed | **pass if designed** |
| **Ordering** | own occurrence-ordering; no cross-stream chronology | lifecycle-local entry order; not K-R15 `(turn, seq)`; entry precedes the first P-1 change by construction | **pass if entry-local** |
| **Values** | lifecycle rebuild from its own entries | `scope`/`class` carried; status slot; no fetch from acquisition or record | **conditional** |
| **No join** | no stored cross-stream join | no edge/join; typed ref only | **pass** |
| **Acquisition purity** | no status/class/scope on acquisition | the entry is lifecycle-side; acquisition unchanged | **pass** |
| **Knowledge separation** | formation does not create knowledge | no acquisition event or edge created | **pass** |
| **OD-a envelope** | separate structures; typed refs | lifecycle identity; typed ref; no shared chronology | **pass** |
| **A2 schema** | no claim A2 settled a schema | the entry is a new representation, not an A2 schema; N-a/N-c stay open | **pass if scoped** |
| **N-b** | no silent status choice | explicit unresolved status slot; value deferred | **conditional** |
| **Decision A** | status obligation preserved | slot plus obligation; concrete value pending N-b | **conditional, not waived** |

---

## 4. Conditions

### 4.1 Lifecycle-local identity and ordering

**FRC-11 (Calculated — identity).** The formation entry must use a **lifecycle-local identity**. It must not reuse an
acquisition event id as its own id. A typed reference to the formation evidence/acquisition event is allowed by OD-a,
but it is **provenance**, not identity or chronology (ECI-15).

**FRC-12 (Calculated — ordering).** The entry must have a **lifecycle-local** occurrence position. It must not borrow
K-R15's `(turn, seq)` as the lifecycle order key. A copied `established_turn` may be a payload value — “one fact
recorded twice,” not two positions compared (ECI-34/35) — but it is not the lifecycle ordering rule.

**FRC-13 (Calculated — entry-local, not a general sequence).** “Before the first P-1 change” can be established by the
entry’s **local position** (the formation entry is the memory’s first lifecycle entry; P-1 versions begin after it).
It must **not** claim a general lifecycle sequence. A global lifecycle order key is **N-c** territory; it is not
decided here.

**FRC-14 (Calculated — if a global order is required, stop).** If the design cannot place the entry before P-1’s
version 1 without a general lifecycle order, the direction needs an N-c decision first. It must not borrow acquisition
ordering to avoid that.

### 4.2 Values carried in the entry

**FRC-15 (Calculated — `scope`).** The entry must carry `scope`. The `MemoryRecord.scope` is a cache; it cannot be the
source. A rebuild reads the formation entry, then applies P-1’s `scope_changed` entries.

**FRC-16 (Calculated — `knowledge_class`).** The entry must carry `knowledge_class`. The record field is a cache; the
stream is the source.

**FRC-17 (Calculated — status slot).** The entry must carry an **explicit unresolved status component**. It is a
representation **slot**, not a third status value; it must not be folded as `provisional` or `active`. N-b supplies
the value/vocabulary. Until then, the stream preserves the obligation but cannot reconstruct a concrete status.

**FRC-18 (Calculated — the record check is conditional for status).** `scope` and `knowledge_class` can be checked
against the stream now. `status` can be checked **structurally** (the slot exists) but not semantically until N-b. That
is an explicit, bounded deferral, not a waiver of Decision A.

### 4.3 No cross-stream dependency or join

**FRC-19 (Calculated — values, not dereferences).** The lifecycle rebuild must obtain the initial values from the
formation entry itself, not by dereferencing the acquisition event or the record. A typed cross-ref may say “entry
grounded in this evidence”; it must not be required to read the entry’s values.

**FRC-20 (Calculated — no stored join).** The formation entry is lifecycle-side; it creates no stored edge/join and no
status/class/scope on the acquisition stream (K-R11/K-R13/K-R14; §6.4).

**FRC-21 (Calculated — formation is not knowledge).** Creating the formation entry does not create an acquisition event
or edge; a memory’s existence does not imply character knowledge (I9; K-R11).

### 4.4 OD-a and A2 compatibility

**FRC-22 (Calculated — OD-a).** The direction uses the lifecycle structure’s own identity and occurrence-ordering
responsibility. It does not require a shared envelope, a shared identifier space or cross-stream chronology. It is
compatible with OD-a Option B and does not reopen it.

**FRC-23 (Calculated — A2).** The formation entry is a **new lifecycle-side representation**, not a claim that A2
settled the lifecycle schema. It must not assert a fold rule, a transition set or a general order key; A2-41…A2-43 and
N-a/N-c remain open.

**FRC-24 (Calculated — the representation is narrow).** It records the initial state; it does not define what later
status events do to state (N-a), how the status stream is ordered (N-c), or who may form (authority).

---

## 5. What this check finds

**FRC-25 (Calculated — compatible, conditionally).** The dedicated direction can coexist with OD-a’s separate
structures and A2 **if** it is lifecycle-local, carries its values, uses an entry-local order, avoids cross-stream
dereferences and keeps the status slot explicit. It fails compatibility if it:

- places status/class/scope on the acquisition stream;
- reuses an acquisition event id as its own identity;
- borrows `(turn, seq)` as the lifecycle order;
- requires a general lifecycle order before N-c;
- dereferences the record or acquisition event for the initial values;
- forces `provisional` or `active` before N-b.

**FRC-26 (Calculated — no compatibility failure found in the direction itself).** The failures are design choices to
avoid, not properties of the dedicated direction. The direction is the right shape; the conditions are the guardrails.

**FRC-27 (Calculated — the status condition is the only obligation-bearing deferral).** The entry can carry `scope`
and `knowledge_class` now. The status slot is present, but the concrete value waits for N-b. Decision A’s obligation
is preserved, not weakened (FRC-03/17/18).

---

## 6. What remains open

**FRC-28 (Calculated).**

| Open | State |
|---|---|
| **kind name** | not selected |
| **final payload** | not selected; at least `scope`, `knowledge_class`, status slot |
| **order-key implementation** | not selected; entry-local position is the constraint |
| **N-b / status value** | open; the slot is explicit |
| **N-c / general lifecycle order** | open; entry-local position avoids it |
| **formation authority** | open; orthogonal |
| **Decision B** | untouched |
| **schema / code / tests / requirements / canonical** | not authorized |

---

## 7. What this record does not do

**FRC-29 (Observed).**

| Untouched | State |
|---|---|
| **P-1 / S-3** | closed; not reopened |
| **Decision A** | resolved Yes; not reopened |
| **formation kind/payload/order** | not selected |
| **N-b / N-c / §8.1** | not decided |
| **Decision B** | untouched |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |

---

## 8. Traceability

**FRC-30 (Observed).**

| What | Where |
|---|---|
| OD-a Option B; separate structures; own identity/ordering; typed refs | ECI-01; A2-45 |
| No lifecycle schema; no fold rule | A2-11/A2-12; A2-41 |
| Status vocabularies; `provisional`/`established`/`active` | A2-42/A2-43; NAE-12/25 |
| Two streams, two checks | K-R11–K-R14; §6.4 |
| Acquisition ordering `(turn, seq)` | K-R15 |
| Copied turn ≠ cross-stream ordering | ECI-34/35 |
| Entry must be lifecycle-side | ECI-17 |
| Decision A owner determination | DAO-01…DAO-14 |
| Existing-entry audit | EPA-11…EPA-31 |
| Formation resolution | P1F-09…P1F-17 |
| Map unchanged | 118 rows |

---

## 9. Record status

**FRC-31 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Direction adopted | **yes** — dedicated formation representation |
| Kind/payload/order selected | **0** |
| N-b / N-c / §8.1 decided | **0** |
| Implementation / tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **coexistence check** — no canonical capture |

**End of check.** The dedicated formation direction **can coexist** with OD-a’s separate structures and A2, provided
it stays lifecycle-local: own identity, own entry-local order (not `(turn, seq)`), values carried in the entry, no
cross-stream dereference or join, and an explicit unresolved status slot (FRC-10…FRC-25). `scope` and
`knowledge_class` can be carried now; the status **value** waits for N-b, while the **obligation** remains (FRC-03/17/18).
No kind, payload or order-key implementation is selected (FRC-05/28). **P-1 and S-3 remain closed; Decision B is
untouched; no schema, code, requirement, map change or canonical capture follows** (FRC-29/31).
