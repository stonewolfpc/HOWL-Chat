# MEMORY-09 — Cross-Stream Temporal Claim: Focused Clarification

**Status:** clarification addendum to `MEMORY-09-A2-ORDERING-AND-EVENT-SCHEMA-RECONCILIATION-2026-11-19.md`.
Read-only. **A2 not reopened. No schema, code, test, probe, model call, canonical edit or ledger entry.**

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Trigger | Owner challenge to finding **A2-24**: "H10/K16 prove A2 sufficient for cross-stream temporal questions" |
| Outcome | **A2-24 is superseded.** A2's cross-stream sufficiency is **Unestablished**, not proven. A2 itself is unaffected |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. The claim that was challenged, verbatim

> **A2-24.** "After" in H10 and K16 is a comparison of **each stream's own position** … Neither requires a
> shared sequence … **This is the strongest evidence that A2 is sufficient** rather than merely convenient.

**XT-01 (Source-derived — owner).** The challenge: this "works when both events have comparable turn
information", but D-3b permits events with **no statement and potentially no turn**. A2 ordinals alone
cannot order such events across streams. The report "should not accidentally imply that A2, by itself,
solves every cross-stream temporal question".

**XT-02 (Observed).** The challenge is correct, and the finding contains **two separate defects** — one in
what I claimed was needed, one in what I claimed supplies it.

---

## 2. Three classes of temporal question

**XT-03 (Inferred).** The original finding collapsed three questions that behave differently:

| | Question | Answered by |
|---|---|---|
| **Q1** | Which event precedes which **within one stream**? | **A2** — a persisted per-stream position |
| **Q2** | Which of two events **across streams** came first, when both carry a shared coordinate? | an **independent comparandum** (e.g. `turn`) — **not A2** |
| **Q3** | Which of two events **across streams** came first, when one carries no shared coordinate? | **nothing currently defined** |

**XT-04 (Inferred — the load-bearing point).** **A2's ordinals are deliberately not a cross-stream
comparandum.** Two positions drawn from independent sequences are not comparable; comparing them
numerically would manufacture exactly the global order A2 excluded (A2-02). So A2 answers **Q1 only**, and
the credit for any Q2 answer belongs to the shared coordinate, not to A2.

**XT-05 (Inferred).** Therefore the correct statement is: **A2 is necessary for Q1 in both streams, and
silent on Q2 and Q3.** The original finding credited A2 with a capability that belongs to `turn`.

---

## 3. Does `turn` supply Q2? Partly — and it is narrower than assumed

**XT-06 (Observed — nullability is annotated, and `turn` is not annotated).** r33 §6.2 marks nullable fields
**explicitly**: `from_character_id` (nullable), `source_memory_id` (nullable); §6.1 marks `asset_ref`
(nullable). Three annotations in §6, and `turn` and `evidence_span` carry none.

**XT-07 (Inferred).** By that convention, **`turn` is declared non-nullable**, and so is `evidence_span`.
D-3b permits events with no statement — which contradicts `evidence_span`'s non-nullability outright, and
puts `turn` in question for the same class of event. §6.2 was written for events that follow utterances.

**XT-08 (Inferred — comparability rests on an unstated assumption).** Q2 via `turn` presumes both streams
use **the same turn coordinate** (a chat turn). §6.2's `turn` and §6.1's `established_turn` are never stated
to be the same coordinate system, and neither declares a scope (per chat? per character?). The assumption is
reasonable and **unstated**.

**XT-09 (Unestablished).** Whether a turn-less event can be assigned the turn in which the system acted.
That value would not be a fabricated *statement position* (SD-04 forbids those), but it is still an inferred
placement in time, and nothing authorises it.

---

## 4. H10 and K16 re-examined

**XT-10 (Observed — the rows, split into setup and assertion).**

| Row | Scenario (setup) | Assertion (what is checked) | Failure |
|---|---|---|---|
| **H10** | source retracted **after** disclosure | account and historical event **persist**; other lifecycle operations not precluded | account removed **because** the source was retracted |
| **K16** | basis retracted **after** promotion | promotion **stands**; no automatic undo | promotion silently reverts |
| **A2a** | memory **later** retracted | character still holds what it acquired | knowledge disappears with status |

**XT-11 (Inferred — the defect in my reading).** In all three the temporal word sits in the **scenario
construction**, not in the assertion. The assertions are **non-cascade properties**: they require that a
later event does *not* retroactively remove an earlier fact. **None requires computing which event came
first.** The ordering is arranged by the test, then assumed.

**XT-12 (Calculated — no ratified row asks for the comparison).** 19 of the 118 rows contain a temporal word
(`after`, `before`, `later`, `prior`, `then`): A2a, B2, B5, B6, C1, D1, D5, F3, G2, H10, K1, K7, K10, K11,
K16, K18, K23, K26, L10. Every one places the word in the scenario column. **No expected column asserts an
ordering; no row mentions cross-stream comparison at all.**

**XT-13 (Inferred — the corrected verdict).** So H10/K16 are **not** evidence that A2 suffices, because they
do not require a cross-stream temporal computation. They are evidence of something else, and weaker:
**the architecture's non-cascade rules are designed to make the comparison unnecessary.**

**XT-14 (Inferred).** This both **narrows** and **strengthens** the position: narrower, because I claimed
proof of sufficiency and have none; stronger, because A2's scope claim no longer rests on a contested
comparison — it rests on Q1 alone, which A2 satisfies by construction.

---

## 5. The turn-less case, stated exactly

**XT-15 (Inferred — what is undefined).** For two events `a` (acquisition stream) and `s` (status stream)
where `s` has no turn:

1. `pos(a)` and `pos(s)` are drawn from different sequences and are **incomparable** (XT-04);
2. `turn(s)` is **absent**, so the shared coordinate is unavailable (XT-06/07);
3. therefore **no defined answer exists** to "did `a` precede `s`?" — the question is not false, it is
   **undefined**.

**XT-16 (Inferred — four possible resolutions, none selected).**

| | Resolution | Cost |
|---|---|---|
| **R1** | the operation does not need the comparison — non-cascade rules already cover it | requires showing no operation needs it; XT-12 suggests this is currently true, but "currently true" is not "never needed" |
| **R2** | a **causal** relation establishes dependence (`caused_by`, `supersedes_event_id`) | establishes *which* produced *which*, **not** order — and v1's link may point backwards (XT-18) |
| **R3** | a **recorded cross-stream comparandum** (shared coordinate on every event) | this is the thing A2 excluded if it becomes a shared sequence; a non-sequential comparandum would need its own definition and its own "no fabricated position" guarantee |
| **R4** | a **turn bucket** for turn-less events | an inferred placement in time, not a statement position — admissible only if explicitly decided, never defaulted |

**XT-17 (Source-derived — owner's boundary, adopted).** R2 does not discharge Q3. **A causal relationship may
establish that one event depends on another; it does not establish which came first.** XT-18 gives the
concrete case.

---

## 6. Why cause cannot stand in for order

**XT-18 (Inferred).** Consider a forgetting event caused by a lifecycle sweep that ran **after** the
memory's last acquisition but was **recorded** against an earlier-cause chain, and a v1-style
`supersedes_event_id` pointing at an event later in the log. The causal edge is correct while the
chronological relation is the opposite of what the edge suggests. **A back-pointing cause is normal, not
corrupt**: correction and supersession are exactly that shape.

**XT-19 (Inferred).** A turn-less status event can therefore have **a fully specified cause and no
chronological position at all**. That is the precise sense in which the two are independent, and it is the
case the original finding failed to consider.

---

## 7. Consequence for the gating decisions

**XT-20 (Inferred).** Q3 becomes a **third gating item**, sequenced *after* OD-a and alongside OD-b:

| Item | Why it is where it is |
|---|---|
| **OD-a** (envelope vs separate) | decides whether a shared coordinate is natural at all — a shared envelope makes a shared turn field cheap; separate structures make it a typed cross-reference |
| **OD-b** (`provisional`: truth or review) | decides what the status stream must represent |
| **OD-g** (Q3 comparandum, *new*) | cannot be answered before OD-a; does not block OD-a or OD-b; **must be answered before any lifecycle schema is drafted**, because a schema may otherwise bake in an assumption that Q3 is always answerable |

**XT-21 (Observed).** OD-c/d/e keep their previous positions. **A2 stays settled. D-3b stays settled.
DN-6a/DN-6b stay open. ND-1 stays open and widened. The 118 are unchanged.** No implementation follows from
this addendum.

---

## 8. Correction record

**XT-22 (Superseded).** In the reconciliation report:

- **A2-23** — the label "genuinely temporal across streams" is replaced by "temporal in scenario setup";
- **A2-24** — the sufficiency claim is **superseded** by XT-11–XT-14 and XT-03–XT-05.

**XT-23 (Observed).** Both were corrected **in place** with a pointer to this addendum, so the earlier
report no longer asserts a proof it does not have. A2-25 (causal ≠ chronological) stands unchanged and is
strengthened by XT-18.

---

## 9. Limits

**XT-24 (Observed).** The 19-row scan is term-based (`after|before|later|prior|then|precede`) over row text;
a row requiring a comparison without those words would not be caught. The nullability convention
(XT-06/07) is inferred from three annotations, not from a stated rule that unannotated means non-nullable.

**XT-25 (Unestablished).** Whether `turn` is a chat-scoped coordinate shared by both streams; whether a
turn-less event may receive a turn bucket; whether any *unratified* future operation would need Q3.

**End of clarification.**
