# MEMORY-09 — PH1A-05: Sub-Message Ordering Assessment (Option O-G)

**Status:** read-only assessment. **No schema change selected or made, no implementation, no test, no
probe, no model call, no corpus, no Git operation.**

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Question | Does an ordered sub-message identifier satisfy K-R15/K-R17/K-R20 and the ND-1 concern, without a per-chat counter field? |
| Baseline | r33 `dfb28b6003ca356f` · harness `eval/mem1/model.py`, `eval/mem1/form.py` |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. The proposal, stated precisely (O-G)

**SM-01 (Inferred, from the owner's description).** Order acquisition events by a **sub-message
identifier** — the parent turn/message number plus an ordinal marking the statement's position within
it — rather than by a separate per-chat `seq` counter:

```
38a  ← first memory-bearing statement in message 38
38b  ← second
39   ← next message, after both
```

The parent number remains the turn; the suffix orders statements inside it. No counter namespace, no
new monotonic field — the claim.

**SM-02 (Inferred).** The proposal is **not cosmetic**. It changes what the ordering key *is*: from a
counter assigned at write time to an **identity of the source statement**. That shift is why it can be
right and why it can silently fail, in ways a counter cannot.

---

## 2. The decisive test: is the label data, or insertion order?

**SM-03 (Source-derived).** K-R15: "**No wall clock and no insertion order participates.**" This fires
directly on the proposal, in two opposite ways:

| How the suffix is produced | Status under K-R15 |
|---|---|
| **Persisted** with the statement, then read by the fold | legitimate — a recorded fact about content order |
| **Computed by enumeration** over a storage list at fold time (`enumerate(spans)`) | **violates K-R15** — that is insertion order, whatever it is named |

**SM-04 (Observed — harness).** The harness already contains **both** shapes:

- `AcquisitionEvent.source_span: str` and `MemoryRecord.source_span` — a **persisted** reference to the
  statement that produced the record, set from `sid` in `form.py:232` and carried onto the event in
  `model.py:_append_event`. This is a genuine sub-message identifier, already stored.
- `event_id = f"ev-{len(self._events) + 1:05d}"` — a counter **derived from append position**. It looks
  exactly like `seq` and is precisely what K-R15's insertion-order clause targets.

**SM-05 (Inferred).** So the harness is a working demonstration that the *carrier* for O-G exists, and a
cautionary demonstration that the *naive* `seq` (SM-04, second bullet) is the forbidden thing wearing the
right name.

**SM-06 (Observed).** The `event_id` counter has a second defect independent of ordering: it is
regenerated from list length on every append, so **a rebuilt log gets different ids**. Under K-R2/K-R14
(a projection must be rebuildable and match), a log whose identity changes on rebuild cannot be compared
to anything. The field is monotonic per store and reproducible from nothing.

---

## 3. Test against K-R15

**SM-07 (Source-derived).** Literal text: "Acquisition events for a (character, memory) pair are totally
ordered by `(turn, seq)`, where `seq` is a monotonic per-chat counter."

**SM-08 (Inferred).** A per-turn suffix is **not** a monotonic per-chat counter: it restarts every turn.
Strictly, `(38, 'a'), (38, 'b'), (39, 'a')` is not a sequence of a chat-global counter. So O-G satisfies
K-R15 only under one of two readings:

| Reading | Consequence |
|---|---|
| **(a) Literal** — the key must be a two-part tuple whose second element is chat-global and monotonic | O-G **insufficient**; K-R15 must be amended to "an intra-turn ordinal", which is a textual change |
| **(b) Functional** — the requirement is a persisted, total, content-derived order; the counter is one implementation of it | O-G **sufficient** if the suffix is total and persisted (SM-03) |

**SM-09 (Inferred).** **O-G does not avoid PH1A-05; it changes which horn is taken.** Under (a) it is a
`§6.2` + `K-R15` change; under (b) it is a `§6.2` carrier change plus a wording clarification. Neither
horn is free, but (b) is much cheaper than adding a counter namespace — which is the owner's point, and
it stands.

**SM-10 (Calculated — cross-turn comparison).** The suffix only orders correctly if the parent is
compared **numerically**:

| Comparison | Result |
|---|---|
| `38b` vs `39` | correct under naive string compare (`'8' < '9'`) |
| `310a` vs `39` | **wrong** (`'1' < '9'` makes 310 sort before 39) |

So the key must be parsed as `(int parent, ordinal)`, or zero-padded. A label scheme that inherits string
sorting is silently wrong from message 100 onward — a failure that would appear only in long chats.

**SM-11 (Calculated — uniqueness).** `(turn, suffix)` is unique per chat **only if** suffixes are unique
within each turn. Reusing `a` for two different statements in one message (e.g. two speakers alternating)
collapses the key. Uniqueness must be stated, not assumed.

---

## 4. Test against K-R17

**SM-12 (Source-derived).** K-R17: a revoke is valid only when "the pair's ordered log currently ends in
a learn event", rejected at write.

**SM-13 (Inferred — this is the proposal's strongest result).** O-G **restores K-R17 to full strength.**
My prior finding was that without a discriminator, the check degrades to turn granularity ("reject a
revoke whose pair had no open episode *before this turn*"). With a persisted suffix, "ends in a learn"
becomes decidable **within** the turn: `38b` extends `38a`, so a revoke at `38b` is rejected or accepted
against the state after `38a` — not against the state at the start of the message.

That is exactly the behaviour S38 asserts, and it is achieved **without a counter**.

**SM-14 (Inferred).** Conditional on the suffix being persisted (SM-03) and unique per pair per label
(SM-17).

---

## 5. Test against K-R20

**SM-15 (Source-derived).** K-R20: `revoked_turn` is not a field; it is derived on demand — "the
**timestamp** of the last revoke event".

**SM-16 (Inferred).** O-G is **neutral** here. It does not fix the existing inconsistency, and it adds a
fourth representable shape to a field already stated three ways:

| Source | Shape of the derived value |
|---|---|
| R1c projection table | `(turn, seq)` tuple |
| S37 | `acquired_turn=40` scalar |
| K-R20 | a "timestamp" — in a document that forbids wall clocks |
| **O-G** | `38b`-style label |

K-R20 needs its own wording decision regardless of PH1A-05's outcome. O-G neither causes nor resolves it.

---

## 6. Test against the ND-1 status-stream concern

**SM-17 (Observed).** Retraction records already carry `source_span` (`model.py:retract(...,
source_span, ...)`), so the **mechanism is available** on the status side too.

**SM-18 (Inferred).** The mechanism being available does **not** resolve ND-1. ND-1 asks whether the
memory stream *needs* an ordering rule; r32 states none, and ordering alone would not define the status
fold (what a later status event does to an earlier one is unspecified). O-G would be a carrier, not a
rule.

**SM-19 (Inferred).** If ND-1 is ever answered "yes, ordering is required", O-G generalises to it at low
cost — an advantage over a counter, which would have to be added to a second stream as well.

---

## 7. The real question: does O-G still need a schema change?

**SM-20 (Observed — the pivot).** **`evidence_span` appears exactly once in r32**: in the §6.2 field list.
It is **never defined** — not in §6.2, §6.3, §6.4, or anywhere else in the contract. So the one existing
field that might already carry a statement identifier is **undefined**.

**SM-21 (Inferred).** That single fact decides whether O-G is free or cheap-but-not-free:

| If `evidence_span` is… | Then |
|---|---|
| an **identifier** of the source statement, with an intra-turn ordinal | **no new field is needed at all** — ordering is `(turn, ordinal(evidence_span))`, and PH1A-05 reduces to a **wording** decision (SM-08 reading (a) vs (b)) |
| a **text quotation / character range** | the fold has nowhere to read an ordinal from, and a field is still required — O-G is a *smaller* field, not no field |

**SM-22 (Unestablished).** Which it is cannot be determined from r32 as written. That is an evidence gap
in the contract, not in this assessment, and it is the single highest-value thing to settle: **one
definitional sentence in §6.2** either closes PH1A-05 almost entirely or moves it back to a field
addition.

**SM-23 (Observed).** The harness's own `source_span` is an identifier (`sid`), not a quotation — but the
harness is not the contract, and its naming is not evidence about r32's `evidence_span`.

---

## 8. Where O-G genuinely fails: granularity

**SM-24 (Inferred — the residual tie).** The label's granularity is the **statement**, but K-R15 orders
**events**. These differ when one statement yields two events for the same pair:

| Case | Pair-level events from one statement | Tie? |
|---|---|---|
| two characters learn from one statement | one per pair (different pairs) | no — per-pair logs |
| one character learns two facts | one per pair (different memories) | no |
| **one character, one memory, two `via` kinds** (e.g. told *and* inferred from the same sentence) | **two in one log** | **yes** |

**SM-25 (Inferred).** So O-G is a total order **iff**:

> **at most one acquisition or forgetting event exists per (character, memory) pair per label.**

No artifact states this. It is plausible in practice (the extractor emits at most one candidate per
rule/value per span), but plausible is not guaranteed, and it is an **extractor property**, not a
contract property — so the contract cannot rely on it silently.

**SM-26 (Calculated).** If a tie ever occurs, O-G fails exactly where a counter would not, and the
remedy reverts to an event-level ordinal. The failure is rare but not impossible, and it is
**data-dependent**, which means it would surface in production rather than in review.

**SM-27 (Inferred).** Note the contrast with the hard case that motivated this: learn-then-forget in one
turn is `38a` then `38b` — distinct labels, so O-G handles it. The proposal solves the case that
matters and leaves a narrower one open.

---

## 9. Sufficiency conditions

**SM-28 (Inferred).** O-G satisfies the ordering requirement if and only if all four hold:

1. **Persisted** — the label is stored, never computed by enumeration at fold time (SM-03).
2. **Numeric parent** — comparison is `(int parent, ordinal)`, not raw string order (SM-10).
3. **Unique per turn** — no two statements in one message share a suffix (SM-11).
4. **At most one event per pair per label** (SM-25).

Conditions 1 and 2 are cheap and within an implementation's control. Condition 3 is a naming discipline.
**Condition 4 is the one that is not guaranteed by anything**, and it is the one that decides whether O-G
is sufficient in general or only in practice.

---

## 10. Revised option set

**SM-29 (Superseded).** My prior statement that **O-A is the smallest change** is superseded by O-G,
which is smaller *if* `evidence_span` proves to be an identifier (SM-21).

| Option | Shape | Requires | Status |
|---|---|---|---|
| **O-G** | order by `(turn, sub-statement ordinal)` | a persisted ordinal — possibly already in `evidence_span`; condition 4 | **new — cheapest if SM-21 lands well** |
| **O-A** | add `seq`, per-chat monotonic | new field, assignment and lifetime rules | unchanged |
| **O-B** | one composite `order_key` | new field + encoding rule | unchanged |
| **O-C** | make `id` monotonic, order by it | **amends K-R15** | normative |
| **O-D** | one event per (pair, turn) | **amends K-R15**, voids S38, contradicts DN-4's example | likely unusable |
| **O-E** | order carried only in the projection array | — | non-viable (circular, K-R2) |
| **O-F** | `seq` + explicit chat field | largest change | unchanged |

**SM-30 (Inferred).** O-G and O-A are not exclusive in every respect: O-G's condition 4 could be
satisfied by *also* having an event-level ordinal, which is O-A. The practical question is whether the
contract may rely on statement granularity alone.

---

## 11. Verdict

**SM-31 (Inferred).**

- **K-R15:** satisfied by O-G under the functional reading, subject to SM-28's four conditions; under the
  literal reading, O-G requires K-R15 to be reworded. It does not escape PH1A-05 — it makes it cheaper.
- **K-R17:** **improved** — O-G restores within-turn decidability and removes the degradation I found.
- **K-R20:** unaffected; the tuple/scalar/timestamp/label inconsistency remains its own decision.
- **ND-1:** mechanism becomes available; the requirement question is untouched.
- **Schema:** **may need nothing at all** if `evidence_span` is an identifier; otherwise a smaller field
  than `seq`.

**SM-32 (Observed).** No remedy is selected here. The proposal is **not premature** — it is the most
promising shape found so far, and it is now blocked on a definitional gap in r32 rather than on a design
question.

---

## 12. What would settle it

**SM-33 (Inferred).** Three answers, in order, none of which requires implementation:

1. **Is `evidence_span` an identifier or a quotation?** (r32 §6.2 — one sentence.) This decides
   "no new field" vs "smaller field".
2. **May the contract rely on one event per pair per statement?** (SM-25.) This decides sufficiency in
   general vs in practice.
3. **Is K-R15 literal about the counter, or about the granularity?** (SM-08.) This decides whether a
   wording change accompanies the carrier change.

None is answerable from the artifacts as they stand. All three are cheap to answer.

---

## 13. Limits

**SM-34 (Observed).** Static reading only. No fold was executed, no label scheme was implemented, no
message store was inspected beyond `eval/mem1`. The `evidence_span` question (SM-22) is an **evidence
gap**, not a finding of defect — it may be defined in a document not read here.

**End of assessment.**
