# MEMORY-09 — P-1 Formation Baseline

## The baseline is already in the formation record; the missing piece is a typed formation event

**Status:** representation-question resolution. **No schema, code, migration, test, requirement, map change or
canonical capture.** P-1 remains the established logical representation. The formation baseline is resolved for P-1's
purposes; the initial state's lack of a typed event remains a pre-existing gap.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | does the existing formation record contain the initial `scope` and `knowledge_class`, so P-1 can build on it without a new kind? |
| Finding | **yes** — formation assigns both on `MemoryRecord` (§4, §6.1, §8.1) |
| Consequence | P-1 uses the formation values as its baseline; **no new formation kind is needed** |
| Remaining gap | the initial `provisional` has a producer but **no typed event representation** (A2-43; ECI-17) — pre-existing, not P-1 |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Task and boundary

**P1F-01 (Observed — the owner's sequence).** Resolve the formation baseline from existing canonical material before
completeness and recovery work. If the existing formation record already contains the initial scope and knowledge
class, P-1 can build on it without inventing another event kind; if it does not, there is a genuine representation gap.

**P1F-02 (Observed — the boundary).** No schema, field definition, code, migration, test, requirement, map change or
canonical capture. P-1 remains the established representation; this record resolves its baseline, not its
implementation.

---

## 2. Canonical evidence

**P1F-03 (Source-derived — §4 Formation, `:279`).** Formation "assign[s] owner chat, **scope**, **knowledge class**,
attribution, time, status" and outputs "Records + events (+ edge projection)".

**P1F-04 (Source-derived — §6.1 MemoryRecord, `:340`).** The record fields include `scope`, `knowledge_class` and
`status`. These are the record's fields, not a separate formation event.

**P1F-05 (Source-derived — §8.1, `:404`).** "Formation writes `scope=world`, `knowledge_class=experiential`,
`status=active`, `established_turn=40`." The worked case confirms both values are written at formation.

**P1F-06 (Source-derived — §11, `:604`).** "evidence turn → **provisional** (status, class, attribution set)." The
diagram names status, class and attribution; §4 and §8.1 independently include `scope`.

**P1F-07 (Source-derived — A2-42/A2-43).** The status-vocabulary reconciliation lists `provisional` as a §5 status
value with **no event**; A2-43 part (i) is superseded only as to the **producer**: `provisional` does have a producer —
formation at §11 `:604` — but **what is missing is a typed event representation**.

**P1F-08 (Source-derived — ECI-17).** "The initial `provisional` **must be representable on the lifecycle side**. It
cannot be inferred from the acquisition structure." This is the constraint the missing typed event leaves open.

---

## 3. Resolution

**P1F-09 (Calculated — the formation record contains the initial values).** Yes. §4 assigns `scope` and
`knowledge_class`; §6.1 lists them as `MemoryRecord` fields; §8.1 shows them written at formation. The owner's
condition for building on the existing record is met.

**P1F-10 (Calculated — P-1 can use the record fields as the baseline).** P-1's projection starts from the formation
value and applies `scope_changed` / `class_changed` entries. No new formation kind is needed for that use.

**P1F-11 (Calculated — the baseline is a record field, not a stream entry).** Formation writes the fields on the
record; there is no typed formation event. Strictly, the baseline is an **initial condition** written by formation, not
a projection input from the lifecycle stream.

**P1F-12 (Calculated — this is a pre-existing gap, not a P-1 gap).** A2-43 and ECI-17 already record that the initial
`provisional` has a producer but no typed event representation. P-1 did not create that gap and does not need to close
it to represent **post-formation** changes. The S-3 acceptance already noted the formation question as pre-existing
(S3A-09).

**P1F-13 (Calculated — the strict reading, stated fairly).** If K-R12's "projection of the stream" is read to require
the **initial** value to be stream-derived too, the current record-field baseline does not satisfy it. Closing that
would require a typed formation event (or an equivalent initial-state entry). That is a **separate representation
decision**, not an implementation detail and not authorized by P-1.

**P1F-14 (Calculated — P-1's scope is post-formation change).** P-1 established `scope_changed` and `class_changed`
for changes **after** formation. It did not establish a formation event. The baseline can be the formation record's
values without extending P-1's kinds.

---

## 4. What this means for P-1

**P1F-15 (Calculated — the baseline is resolved for P-1).**

| Question | Answer |
|---|---|
| does the formation record contain initial scope/class? | **yes** — §4, §6.1, §8.1 |
| can P-1 build on it without a new kind? | **yes** |
| is the baseline stream-derived? | **no** — it is a record field; no typed formation event |
| does that invalidate P-1? | **no** — P-1 governs post-formation changes; the initial-state gap is pre-existing |
| what would full stream-derivation require? | a typed formation event or equivalent — a **separate decision** |

**P1F-16 (Calculated — the owner's condition).** The owner said: if the formation record already contains the initial
scope and class, P-1 can build on it without inventing another event kind. It does, so P-1 can. The "genuine
representation gap" the owner anticipated does **not** arise for the baseline values themselves; the separate gap is
the **typed event** for the initial state, which is pre-existing.

**P1F-17 (Calculated — no new formation kind authorized).** Adding a `memory_formed` kind (or similar) is not part of
P-1. It would be a new representation element requiring a separate owner decision. This record does not propose or
authorize it.

---

## 5. The remaining gap

**P1F-18 (Observed — the initial state's representation).** `provisional` has a producer (formation) but no typed
event (A2-43). ECI-17 says it must be representable on the lifecycle side. This is unresolved.

**P1F-19 (Calculated — the overlap with N-b).** The `provisional` representation is part of N-b (`active`/`established`;
`provisional`), which remains separate and open (NSD-06/07; C23-22). This record does **not** decide N-b or the
formation event; it records that both are pre-existing questions.

**P1F-20 (Calculated — what an implementation authorization would need).** If implementation is to proceed with the
record-field baseline, the authorization request must state that the initial state is the formation record's values,
not a typed event. If full stream-derivation is required first, a separate representation decision is needed.

---

## 6. What this record does not do

**P1F-21 (Observed).**

| Untouched | State |
|---|---|
| **P-1 representation** | unchanged; baseline resolved for its use |
| **new formation kind** | not proposed or authorized |
| **N-b / `provisional`** | not decided |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |

---

## 7. Traceability

**P1F-22 (Observed).**

| What | Where |
|---|---|
| Formation assigns scope/class/status | §4 `:279`; §8.1 `:404` |
| MemoryRecord fields | §6.1 `:340` |
| §11 provisional entry | `:604` |
| `provisional` producer, no typed event | A2-42/A2-43 |
| Initial state must be lifecycle-side | ECI-17 |
| S-3 noted the pre-existing question | S3A-09 |
| P-1 baseline | P1D-06 |
| N-b remains separate | NSD-06/07; C23-22 |
| Map unchanged | 118 rows |

---

## 8. Record status

**P1F-23 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| New formation kind | **0** |
| Implementation / tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **representation resolution** — no canonical capture |

**End of resolution.** The existing formation record **does** contain the initial `scope` and `knowledge_class`
(§4, §6.1, §8.1), so P-1 can use the formation values as its baseline **without a new formation kind** (P1F-09/10).
The baseline is a record field, not a typed stream entry; the initial `provisional` has a producer but **no typed event
representation** (A2-43; ECI-17) — a **pre-existing gap**, not a P-1 failure (P1F-11/12). A fully stream-derived
baseline would require a separate formation-event decision, which P-1 did not authorize (P1F-13/17). **N-b and
`provisional` remain separate** (P1F-19). **No schema, code, requirement, map change or canonical capture follows.**
