# MEMORY-09 — PH1A-05 Decision Brief

**Status:** decision brief. Read-only. **No schema change, no implementation, no test, no probe, no model
call, no corpus, no Git operation.** Nothing is selected here; each item states what a decision would
settle and what it would not.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Scope | Four items: `evidence_span` definition · K-R15 wording · event provenance · reconstruction check |
| Carried separately | ND-1 (status fold) · the `chat_id` drop (`MEMORY-09-CHAT-ID-DROP-RECORD-2026-11-19.md`) |
| Baseline | r33 `dfb28b6003ca356f` · `DESIGN-MEMORY-CONTINUITY-2026-11-11.md` · `eval/mem1/` |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 0. What changed in the diagnosis

**DB-01 (Inferred).** Two passes ago this looked like a **missing ordering subsystem**. It is not. There
is an evidence model — `source_message_id`, `turn_id`, `evidence_span` in the v1 design, and `turn`,
`source_span` in the harness — whose **contract is incomplete and whose use is inconsistent**. That is a
reconciliation problem, not a construction problem, and the brief is scoped to it.

**DB-02 (Observed).** Nothing in this brief proposes a schema edit. Three of the four items may resolve
without one; the fourth is an implementation defect.

---

## 1. Item D-1 — Define `evidence_span`

### 1.1 Why it is the blocker

**DB-03 (Observed).** `evidence_span` appears in four schemas and in **zero definitions**. In the v1
`Event` it sits beside `source_message_id` and `turn_id`, so it is a span *within* a message — but
whether it is an **identifier** or a **text range** is stated nowhere.

**DB-04 (Inferred).** Every downstream question inherits that ambiguity: whether ordering is
`(turn, ordinal)` or `(turn, start offset)`; whether ids are stable enough to order by; whether two
events can share a span.

### 1.2 What a definition must settle — four properties

| Property | Why it must be stated | Present today? |
|---|---|---|
| **Identity** | what the value denotes — a statement, a character range, a quotation | **unstated** |
| **Position** | whether it carries, or implies, a position within the turn | implied by both readings, stated by neither |
| **Stability** | whether it survives rebuild, reload and re-derivation | **unstated**; depends on construction |
| **Comparison** | the rule that orders two spans in one turn | **absent** |

**DB-05 (Observed — a precedent worth copying).** The harness already enforces *resolvability*: it raises
`ProvenanceFabricationError` when a span does not resolve. So "a span must resolve against its message" is
an existing, respected invariant. The missing properties are **stability** and **comparison**.

### 1.3 The options

| Option | Shape | Gains | Costs |
|---|---|---|---|
| **E-1** | opaque identifier per statement | reference and dedupe already work | carries no position unless ids are built to encode one — which makes ids semantic and load-bearing |
| **E-2** | text range (start/end) within `source_message_id` | position is intrinsic; comparison is integer; stability follows from immutable text | a second field or a structured value; spans must be non-overlapping or the order is partial |
| **E-3** | both — id for reference, range for order | each used for what it is good at | two representations of one thing, which is the defect shape r32 warns about elsewhere ("must not become competing sources of truth") |

**DB-06 (Inferred).** E-2 is the only option where **position is a property of the evidence itself**
rather than of how it was stored. E-1 can be made to work, but only by making identifiers order-bearing,
which is the same move it appears to avoid.

**DB-07 (Unestablished).** Which r32 intends. One definitional sentence plus one comparison clause closes
this item. **Decision needed:** E-1, E-2 or E-3, and the four properties stated explicitly.

---

## 2. Item D-2 — K-R15's wording

### 2.1 The tension, stated exactly

**DB-08 (Source-derived).** The **purpose** is stated separately from the rule:

> "`revoked_turn` must be **entirely reproducible** from the acquisition stream, including ordering and
> multiple learn/forget cycles." — R1c review concern
> "replaying the same ordered log produces the same edge, **on any machine, in any insertion order**."
> — §6.4

**DB-09 (Source-derived).** The **rule** is mechanism-level: "totally ordered by `(turn, seq)`, where
`seq` is a monotonic per-chat counter … Events sharing a turn are resolved by `seq` alone."

**DB-10 (Inferred).** The goal is determinism; the counter is one means. Ordering by a persisted span
position meets the goal and does not match the sentence.

### 2.2 The options

| Option | Shape | Protects | Costs |
|---|---|---|---|
| **K-1** | keep the literal counter; add `seq` | exact literal conformance | new field and namespace for a property an existing carrier may already supply |
| **K-2** | reword to the property: "totally ordered by a persisted, content-derived position within the turn; no wall clock, no insertion order; events sharing a turn resolved by that position alone" | determinism preserved verbatim in intent; span position satisfies it | an amendment to K-R15 and to S38's wording ("by `seq`" → "by position") |
| **K-3** | keep `seq` optional, fall back to span position | nothing breaks immediately | **two sources of order** — the exact "competing sources of truth" pattern r32 rejects for acquisition vs edge |

**DB-11 (Inferred).** K-3 is the worst of the three and should be named so it is not chosen by default:
it makes "which event is last" depend on which field happens to be populated.

**DB-12 (Inferred).** K-2 is an amendment, and amendments are the owner's to make — but it is a
**one-clause** amendment that keeps the guarantee that matters. K-1 is a **new field** that duplicates a
capability the evidence model already has.

**DB-13 (Unestablished).** Whether the contract may rely on position **alone**. If any event can exist
without a source statement (§4), position is incomplete and K-2 needs a fallback rule — which is K-3 by
another name, with the same hazard. **This dependency is why D-2 should be decided after D-3.**

---

## 3. Item D-3 — Event provenance (the priority)

### 3.1 The defect

**DB-14 (Observed — `model.py:303-312`).** `revoke(record_id, turn, reason)` takes no span and builds its
event with `source_span=rec.source_span` — **the span of the statement that created the record**.

**DB-15 (Source-derived).** This is not only an ordering problem. §6.2 gives **every** `AcquisitionEvent` an
`evidence_span`, including events with `via=revoked`. An event that inherits another statement's span does
not satisfy its own schema field. Consequences compound:

| Rule | How it fails |
|---|---|
| **K-R3** (forgetting is a new event) | the event exists but does not record *its own* occurrence |
| **K-R90** ("how did they learn this?" from events) | the log now answers a question about learning with a forgetting event carrying the learning's location |
| **I11** (acquisition history authoritative) | the authoritative record of the loss points at the wrong moment |
| **§6.4 fold determinism** | ordering by span position puts the loss at the learning's position |

**DB-16 (Inferred).** A fact learned in turn 12 and forgotten in turn 40 would, under span ordering,
reconstruct as lost at turn 12's statement. The history is not merely mis-ordered; it is **false**.

### 3.2 The second half, which the defect exposes

**DB-17 (Inferred — this is the architecturally important part).** The fix "record the forgetting
statement's span" **presupposes that a forgetting statement exists.** It may not:

- forgetting could be caused by the character's model, by the operator, or by a system action;
- **DN-6a is open** — who may cause forgetting is undecided;
- nothing in T2 requires an event to originate in a message.

**DB-18 (Inferred).** So there are two decisions here, not one, and they are of different kinds:

| | Question | Kind |
|---|---|---|
| **D-3a** | must every acquisition-family event carry the span of the statement that caused *it*? | fix to an implementation defect |
| **D-3b** | may an event exist with **no** source statement at all? | architecture decision, dependent on DN-6a |

**DB-19 (Inferred).** D-3b is decisive for D-2. If span-less events are possible, then span position
cannot be the universal key, and K-2 needs a companion rule — which reintroduces K-3's hazard unless the
rule is "a cause position is assigned at write and stored", i.e. a counter under another name. **This
should be settled before K-R15 is reworded.**

**DB-20 (Observed).** Nothing here authorises the code fix. D-3a is an implementation defect; repairing it
is separate work with its own authorisation.

---

## 4. Item D-4 — The reconstruction check

### 4.1 What it must establish

**DB-21 (Inferred).** Four histories must be reproducible **regardless of insertion order**:

| History | Rebuild target | Ordering needed? |
|---|---|---|
| acquisition | the edge (present / absent, `via`, `verified`, `acquired_turn`) | yes |
| forgetting | the same edge, from the other direction | yes |
| retraction | the memory's status projection | ND-1 — **unspecified** |
| status | the same projection | ND-1 — **unspecified** |

**DB-22 (Observed — a prerequisite defect).** `event_id = f"ev-{len(self._events) + 1:05d}"` is derived
from append position, so **a rebuilt log is renumbered**. Comparison against the original is therefore
impossible until identity is stable — which is D-1's **stability** property, not a separate item.

**DB-23 (Observed).** Neither implementation has a rebuild routine for the acquisition log: the harness
projects on write and never folds, and the engine keeps a status field. **D-4's check is currently
unrunnable**, and would remain unrunnable after any schema decision until a fold exists.

**DB-24 (Inferred).** So D-4 is not a test to run now; it is a **property to state**, and its statement
depends on D-1 (stability) and D-2 (what order means). It cannot be closed first.

---

## 5. Item order

**DB-25 (Inferred).** The dependencies run:

```
D-3b (can an event lack a statement?)  ─┐
D-1  (what is a span; is it stable?)   ─┼─→ D-2 (reword K-R15 how?)
D-3a (fix the reuse defect)            ─┘        │
                                                └─→ D-4 (state the reproducibility property)
```

**DB-26 (Inferred).** D-3b first is the counter-intuitive but correct order: it is the only item whose
answer can **invalidate** the span-position approach. D-3a is independent and can proceed as a defect
fix under separate authorisation.

---

## 6. ND-1 — kept separate

**DB-27 (Source-derived).** The status stream has an ordered carrier available (retraction records carry a
span) and **no fold rule**. Ordering and lifecycle are different questions: knowing *which* status event
is later does not say what a later status event *does* to an earlier one. ND-1 remains **open** and is not
touched by D-1…D-4.

---

## 7. What this brief does not do

**DB-28 (Observed).** It selects no option, proposes no schema edit, authorises no code change, and does
not reopen DN-1, DN-4, NEW-1 or O-14. It does not restore `chat_id` (see the separate record). It does not
decide DN-6a/DN-6b, and it does not claim the ordering problem is solved.

---

## 8. Verification

r33 `dfb28b6003ca356f` · v9 `698254ca262c0874` · ledger `e5008020ff68bd21` · map `b4ba33ccbfe8f8c1` —
all unchanged. Tests executed: **0**.

**End of brief.**
