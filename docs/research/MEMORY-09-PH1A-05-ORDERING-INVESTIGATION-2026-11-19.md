# MEMORY-09 — PH1A-05 Investigation: K-R15 and the Missing `seq`

**Status:** read-only investigation. **No schema change proposed as selected, no schema change made, no
implementation, no test, no probe, no model call, no corpus, no Git operation.**

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Question | PH1A-05 — can ordering by `(turn, seq)` be satisfied, and what is actually missing? |
| Baseline | r33 `dfb28b6003ca356f` (contains r32) · v9 `698254ca262c0874` |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. Exact source wording

**IN-01 (Source-derived).** **K-R15**, r33 §4 R1c:

> Acquisition events for a (character, memory) pair are totally ordered by `(turn, seq)`, where `seq` is a monotonic per-chat counter. **No wall clock and no insertion order participates.** Events sharing a turn are resolved by `seq` alone.

**IN-02 (Source-derived).** **§6.2 `AcquisitionEvent` (authoritative, append-only)** — the complete field
list, verbatim:

> `id`, `character_id`, `memory_id`, `turn`, `via` ∈ `witnessed | told | read | inferred | verified |
> narrator_grant | revoked`, `from_character_id` (nullable), `source_memory_id` (nullable),
> `confidence_in_source`, `evidence_span`.

**IN-03 (Source-derived).** **§6.3 `KnowledgeEdge` (projection):** `character_id`, `memory_id`,
`acquired_turn`, `via`, `verified`, `derived_from_event_ids[]` — "Rebuilt by folding acquisition events
for the pair (R1c); never edited directly." The projection table in R1c defines `acquired_turn` as
"`(turn, seq)` of the learn event that opened the current episode" and `derived_from_event_ids` as "the
events of the current episode, **in order**".

**IN-04 (Source-derived).** **§6.4 fold determinism:** "replaying the same ordered log produces the same
edge, on any machine, in any insertion order; and a missing edge is not 'unknown state' to be guessed —
it is a determinate consequence of the last event."

**IN-05 (Source-derived).** **K-R17:** "`revoked` is valid only when the pair's ordered log **currently
ends in a learn event**. A revoke with no open episode is **rejected at write**, not silently ignored."

**IN-06 (Source-derived).** **K-R20:** "`revoked_turn` is therefore **not a field**. It is derived on
demand — the timestamp of the last revoke event."

**IN-07 (Source-derived).** **S38:** "a revoke with no open episode is rejected at write; same-turn
learn+revoke resolves deterministically by `seq`, never by insertion order". **S37:** "told at t10,
revoked at t20, witnessed at t40 ⇒ edge `acquired_turn=40` …".

**IN-08 (Observed).** **ARCH-03 F6**, from the independent review record (`:158`):

> **`AcquisitionEvent` has no `seq`** (and no chat field). §6.2 lists `id, character_id, memory_id, turn,
> via, from_character_id, source_memory_id, confidence_in_source, evidence_span`. K-R15 orders by
> per-chat `(turn, seq)`, and S38 asserts same-turn resolution "by `seq`". Same-turn ordering is **not
> representable** as specified.

**IN-09 (Observed).** **ARCH-03 F25** (`:192`): "`acquired_turn` is defined as `(turn, seq)` but named
and asserted as a turn (`acquired_turn=40`, S37)".

---

## 2. The nine questions

### Q1 — What exactly must be ordered?

**IN-10 (Source-derived).** The ordered object is **the acquisition-event log of one `(character,
memory)` pair** — not the global event stream, not the memory stream, not the chat transcript. K-R15
says "for a (character, memory) pair"; the fold in R1c is defined per pair (`known(pair) ⇔
ordered_log(pair)[-1] is a learn event`).

**IN-11 (Source-derived).** Because `revoked` is a member of `via`'s closed set, **a forgetting event is
an `AcquisitionEvent`** and lies in the same per-pair log as the learns it closes.

### Q2 — Total order, per-stream order, or something else?

**IN-12 (Source-derived).** A **total order**, but a **scoped** one: total over each pair's log
("totally ordered"), keyed by `(turn, seq)` lexicographically, where `seq` is **monotonic per chat**.
K-R15 does not require a global total order across pairs, and does not require ordering by insertion
position. Two distinct properties are demanded at once: **a total order within the pair**, and **a key
that is not any ambient ordering** ("no wall clock and no insertion order participates").

**IN-13 (Inferred).** The scoping is mixed and this matters: the *log* is per pair, the *counter* is per
chat. A per-pair log ordered by a per-chat counter means seq must be globally unique within the chat even
though only pair-relative order is used.

### Q3 — Does `(turn)` alone genuinely fail?

**IN-14 (Source-derived).** Yes, and in four independent places:

1. two events at the same `turn` for one pair are **indistinguishable** in §6.2 — no listed field
   discriminates them;
2. `acquired_turn` is specified as a **`(turn, seq)` tuple** but the events carry only `turn`;
3. `derived_from_event_ids[]` must be "in order", and that order must come from the log — which has
   none within a turn;
4. **S38 cannot be executed**: it asserts resolution "by `seq`", and `seq` is absent.

**IN-15 (Inferred — the sharpest consequence).** **K-R17 silently degrades.** "The ordered log currently
ends in a learn event" requires knowing which event is last. With a tie on `turn` and no discriminator,
the check can only be applied at *turn* granularity — i.e. "reject a revoke whose pair has no open
episode *before this turn*". That is a weaker rule than the one written, and it weakens precisely in the
case DN-4 creates.

### Q4 — Can same-turn events be distinguished by any existing field?

**IN-16 (Calculated).** Each candidate fails, for a stated reason:

| Field | Why it cannot order |
|---|---|
| `turn` | ties by construction |
| `id` | K-R16 requires **uniqueness only**; nothing states ids are monotonic, and ordering by them would be an unstated rule |
| `via` | K-R26 orders **evidence strength**, not time; and a weaker event legitimately follows a stronger one (K-R23) |
| `from_character_id`, `source_memory_id` | provenance, not sequence |
| `confidence_in_source`, `evidence_span` | content, not sequence |
| `derived_from_event_ids[]` | lives on the **projection**, only for the **current episode**, and is itself derived — ordering by it is circular |
| memory's `owner_chat_id` | not on the event; and chat scope is not order |

**IN-17 (Inferred).** Chat scope *is* derivable (`memory_id → memory.owner_chat_id`), so the counter's
*namespace* exists. The counter's *value* does not, and no field can stand in for it. **The answer is no.**

### Q5 — Does the answer differ by event kind?

**IN-18 (Source-derived).** **Acquisition and forgetting: identical.** `revoked` ∈ `via`, so both are
acquisition events in one per-pair log under one rule.

**IN-19 (Source-derived).** **Status changes: different, and the difference is an absence.** Retraction,
correction, supersession and authority changes live on the **memory stream** (K-R12). r32 specifies that
each stream must be independently rebuildable (K-R14, §6.4 "two streams, two checks") but specifies **no
ordering key at all** for the memory stream — there is no `seq`, no analogue of K-R15, and no statement
that status is last-write-wins.

**IN-20 (Unestablished).** Whether the memory stream **requires** an ordering rule is not settled by
anything read. It is not "unconstrained"; it is **unspecified**. This is recorded as a newly discovered
dependency (ND-1, §6) and is **not** turned into a requirement here.

**IN-21 (Inferred).** A forget (acquisition stream) and a retraction (memory stream) in the same turn do
**not** need a cross-stream order: §6.4 makes the join read-time and never stored (K-R13). The ordering
problem is intra-stream only.

### Q6 — Schema defect, projection defect, or both?

**IN-22 (Observed).** **Primarily a schema defect**: §6.2's field list omits the key that K-R15 names and
that §6.3's tuple requires. ARCH-03 F6 states this as observed, not inferred.

**IN-23 (Observed).** **Secondarily a specification inconsistency**, independent of any implementation —
three artifacts of r32 disagree about what `acquired_turn` is:

| Location | Says |
|---|---|
| R1c projection table | `(turn, seq)` — a tuple |
| **S37** | `acquired_turn=40` — a scalar turn |
| **K-R20** | `revoked_turn` is "the **timestamp** of the last revoke event" — a third vocabulary, in a document that forbids wall clocks |

**IN-24 (Inferred).** It is **not** a projection-algorithm defect. Given a log with a real order, the
fold in R1c is well defined and deterministic. The projection fails only because its input key is
missing and its own field's type is inconsistently stated.

### Q7 — Can the ordering be guaranteed without changing the schema?

**IN-25 (Inferred).** **No — not within the specified field set.** Every candidate route is excluded by
text: insertion order is explicitly excluded by K-R15 and by §6.4 ("in any insertion order"); wall clock
is excluded; id ordering is unstated and would itself be a new rule; and ordering by a derived array is
circular.

**IN-26 (Inferred).** An implementation could add a column without touching the contract, but that is
**not** a remedy: conformance is assessed against §6.2, so the implementation would be relying on a
field the specification does not carry, and the fold would remain unspecified for anyone else building
it. A private column is a workaround, not a fix.

### Q8 — What exactly makes a change unavoidable?

**IN-27 (Source-derived).** The minimal forcing pair is:

- **K-R15** names `seq` as *the* discriminator and **forbids** the two ambient orders that could
  substitute for it; and
- **§6.2 does not contain `seq`**.

Together these make the requirement and the schema mutually unsatisfiable. §6.3's tuple and **S38**'s
"by `seq`" are corroborating, not independent, constraints.

**IN-28 (Inferred).** There is no third option: either the schema gains an order-bearing field, or K-R15
and S38 are amended. Both are changes; neither is a mere clarification.

### Q9 — What are the remedies, and what does each require?

**IN-29 (Observed/Inferred).** Listed for comparison only. **None is selected.**

| Option | Shape | What it would additionally require | Cost |
|---|---|---|---|
| **O-A** | add `seq` — integer, monotonic per chat, assigned at write | who assigns it; reset/lifetime semantics; uniqueness scope; whether the fold compares `(turn, seq)` lexicographically; migration for existing events | smallest textual change; leaves the pair/chat scoping mixed (IN-13) |
| **O-B** | add one composite `order_key` field instead of `seq` | an encoding rule and a comparison rule; same assignment questions as O-A | one field, but a new type and a new comparison rule |
| **O-C** | make `id` monotonic and order by it | **amends K-R15** (it names `(turn, seq)`), K-R16's semantics, and S38; ids become semantically loaded | normative edit, not schema-only |
| **O-D** | forbid more than one acquisition event per (pair, turn) | **amends K-R15 and voids S38**; contradicts DN-4's own same-turn learn+forget example | removes the capability the decision needs |
| **O-E** | carry order only in `derived_from_event_ids[]` position | — | **non-viable**: K-R2 requires the projection be rebuildable from events; a projection that is the sole carrier of order cannot be rebuilt (circular) |
| **O-F** | add `seq` **and** an explicit chat field | all of O-A plus a rule that the counter is per chat, self-contained rather than derived through the memory | largest change; removes a derivation dependency |

**IN-30 (Inferred).** O-A, O-B and O-F are schema changes and need PH1A-05's answer. O-C and O-D are
**normative** changes and need an owner decision beyond PH1A-05 — they are not remedies for a schema
defect, they are amendments to the requirement.

---

## 3. Same-turn examples

### 3.1 Acquire, then forget, in one turn

**IN-31 (Source-derived).** Required outcome under K-R19/K-R17: the log is `learn(T) → revoke(T)`; the
log ends in a revoke, so **no edge exists**; the history holds two events; "did they ever know" answers
**yes**; `revoked_turn` is derived `(T, seq)` for inspection.

### 3.2 Forget, then reacquire, in one turn

**IN-32 (Source-derived).** A valid revoke requires an open episode, so this case presupposes knowledge
from an **earlier** turn. The log is `… learn(T₀) → revoke(T) → learn(T)`; the log ends in a learn, so
**an edge exists**, `acquired_turn` is the `(T, seq)` of the **reacquiring** event, and `via` comes from
that event — "relearning after forgetting cannot shortcut" (R1c).

### 3.3 The decisive point

**IN-33 (Calculated).** Both cases store **the same unordered multiset**: one learn at `T`, one revoke at
`T`, and a prior learn at `T₀` in the second case only — which does **not** disambiguate, because in case
3.1 the revoke is valid by virtue of the same-turn learn. Strip the order and the two cases are the same
data:

| Stored, order-stripped | Case 3.1 result | Case 3.2 result |
|---|---|---|
| `{learn@T, revoke@T}` (+ prior learn@T₀) | **no edge** | **edge present**, `acquired_turn=(T, seq)` |

**IN-34 (Inferred).** Two different conforming states follow from the same stored data. This is **not** a
taste question about representation; it is a **determinism failure** in the fold, and §6.4 requires the
opposite ("replaying the same ordered log produces the same edge … on any machine"). The events are
ordered in reality and unordered in the schema.

**IN-35 (Inferred).** Note also that the *rejection* rule cannot separate them. In 3.1 the revoke must be
accepted because the same-turn learn precedes it; in 3.2 the revoke must be accepted because an earlier
learn is open. A writer that received the revoke first could not tell whether it should be rejected under
K-R17 or accepted under either case — which is IN-15 restated at the write boundary.

---

## 4. The required separation

| Bucket | Content |
|---|---|
| **What the normative requirement demands** | a total order per (character, memory) pair by `(turn, seq)`, seq monotonic per chat, no wall clock, no insertion order (K-R15); the log's end determines validity (K-R17) and the fold (R1c); the projection carries `(turn, seq)` and an ordered id list (§6.3); replay is deterministic (§6.4); S37/S38 test it |
| **What the schema can represent** | `turn` (rank one only); unique `id`; `via` including `revoked`; **no `seq`**, no chat field, no order-bearing field (IN-02) |
| **What the implementations do** | harness: events in Python lists, ordered by **insertion**, `established_turn`/`revoked_at_turn` at turn granularity, edges emitted 1:1 per event including `revoked` (violates K-R19); engine: **no event log at all**, status as a field |
| **What the harness can reconstruct** | the current status; **not** within-turn order; **not** invalid-revoke rejection; **not** `(turn, seq)`; its edge is a status flip, not a fold |
| **What remains unestablished** | whether the memory/status stream needs an ordering rule (IN-20); whether `character_id` implies a single chat; whether K-R20's "timestamp" means `(turn, seq)`; whether S38 is well-formed at all given the schema |

**IN-36 (Inferred).** No new normative requirement is created by any of the above. The failure is that
existing text cannot be satisfied by existing fields — not that the text asks for too little.

---

## 5. What this investigation does not do

**IN-37 (Observed).** It selects no remedy, proposes no schema edit, and changes no artifact. PH1A-05
remains **open**. Whether the answer is O-A, O-B, O-F, or a normative amendment (O-C/O-D) is an owner
decision, and O-C/O-D would additionally require amending K-R15 — which this investigation does not
authorise.

**IN-38 (Inferred).** DN-6a/DN-6b are untouched: ordering is orthogonal to authority.

---

## 6. Newly discovered dependency (undecided)

**ND-1 — the memory/status stream has no ordering rule.** K-R12 places status on its own stream and
K-R14/§6.4 require independent rebuildability, but r32 specifies no order key for that stream (IN-19).
Two status events in one turn are as indistinguishable as two acquisition events, and no fold rule
("last write wins", "precedence", or otherwise) is stated. Whether this matters depends on whether the
experiment can produce two status events in one turn — **unestablished**. Recorded, not decided, and not
added to any requirement set.

---

## 7. Verification

| Check | Result |
|---|---|
| r33 hash | `dfb28b6003ca356f` — unchanged |
| v9 hash | `698254ca262c0874` — unchanged |
| Schema or code edited | **0** |
| Tests / probes run | **0** |

**End of investigation.**
