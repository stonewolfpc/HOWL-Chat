# TASK-03 — OQ-6: settlement evidence

| Field | Value |
|---|---|
| Purpose | What durable evidence of a completed generation transition would have to be |
| Status | **Design investigation. Proposal only. No implementation authorized.** |
| Subject | V3.1 rows 16–19, clauses V3-17, V3-31, V3-37 |
| Recovery token | `TASK-03-OQ6-SETTLEMENT-2026-11-22` |

Prefix: `S-` proposal clauses. None is adopted; V3.1 is unchanged by this document.

---

## 1. The question, restated precisely

**S-01 — The marker already answers "is the current state settled?" by observation.** V3-31 defines settlement as *marker `(cleared,G)` and both representations valid and void*. That is checkable at the moment it is true.

**S-02 — The gap is the past, not the present.** Once the marker advances to `(committed,G+1)`, the observation that made G settled is gone: V3-17 records that the marker has no history. Rows 18–19 exist because a reader cannot reconstruct that moment.

So the question is not "how do we know a clear finished?" but **"how does a later reader know it finished before the next content was written?"**

**S-03 — A settlement record must prove completion at a point the append could not precede.** Two facts, and the second is the hard one:

| # | Fact | Sufficient? |
|---|---|---|
| 1 | both representations were void at time T | no — says nothing about when the append began |
| 2 | no append could begin before T | **required** |

Fact 2 is not a property of a record. It is a property of the **writer** — which is OQ-7. **Settlement evidence alone can never establish purity**; it establishes only that the transition completed.

---

## 2. Where the evidence should live

**S-04 — The recommendation is to extend the marker document, not to create a second record.**

| Option | Consequence |
|---|---|
| **A. fields in the existing marker** | one document, so the settlement and the pair are replaced by the same rename; atomicity, CAS, pre-images and crash tests are inherited unchanged; no new crash domain |
| B. a separate settlement file | a second document with its own lifecycle, its own crash points, and a two-document consistency problem — precisely the ambiguity A-07 warns against |
| C. settlement inferred by the scanner | forbidden: the scanner would be asserting completion of an operation it did not perform (V3-56) |

**S-05 — Proposed shape**, additive to the three existing fields:

| Field | Meaning |
|---|---|
| `settled_generation` | the highest generation whose clear is recorded as completed |
| `settled_at` | bookkeeping only; never evidence |

Written in the **same CAS replacement** as the pair, so there is no ordering window between "the pair advanced" and "the settlement exists".

**S-06 — Write timing.** The writer closes generation G: rewrite log → rewrite snapshot → both observed valid and void → CAS the marker to `(cleared,G)` with `settled_generation = G`. The next generation's first write CASes `(pending,G+1)` and **must require** `settled_generation ≥ G` in the compared value.

**S-07 — The CAS comparison is what makes the field load-bearing.** If settlement is merely stored and not required by the transition, it becomes a comment. The transition's compare-and-swap must include it: `(state, generation, settled_generation)` is the compared triple.

---

## 3. Crash behaviour

**S-08 — No new crash points are introduced.** Because the settlement rides in the marker document, its crash outcomes are the document's: `Saved`, `Recovered`, or absent (V3-19).

**S-09 — Staleness fails closed.** A `Recovered` marker carries its *older* `settled_generation`, so it under-certifies: a transition that needed `settled ≥ G` sees a smaller value and refuses. There is no direction in which a recovered marker over-certifies. **This is the property that makes option A safe and option B not obviously safe.**

**S-10 — Death points, exhaustively:**

| Death after | Marker | Content | Result |
|---|---|---|---|
| both rewrites, before CAS | `(cleared,G)`, `settled < G` | void | Pending `unattributed`-family; retry re-runs the same clear (V3-33) |
| CAS, before next generation's write | `(cleared,G)`, `settled = G` | void | Healthy `closed`; the next append is now permitted |
| next generation's CAS, before append | `(pending,G+1)`, `settled = G` | void | Pending `awaiting_content`; the settlement survives, so the transition is provable later |
| append torn | `(pending,G+1)`, `settled = G` | invalid log | Damaged `truncated`; the settlement is not in question |

**S-11 — The record is monotonic.** `settled_generation` never decreases except by `Recovered` staleness, which is the fail-closed direction. No clause may lower it deliberately.

---

## 4. What this does and does not buy

**S-12 — Purity requires both records:**

```
purity(G ≥ 2) = settlement_record(G−1)  ∧  enforcement_evidence
```

Neither term alone is sufficient (S-03), and neither may be derived from the other (A-05). A conversation can be settled-but-ungated, or gated-but-unknown-settled; both remain Pending.

**S-13 — No backfill.** A settlement record must never be created by inference — not by the scanner, not by a migration, not by "the marker looks cleared and the content looks consistent". Backfilling would manufacture exactly the evidence the field exists to provide. **Historical generation ≥2 conversations therefore stay Pending even after this field ships**; that is the honest boundary and it does not weaken V3-37.

**S-14 — The field's trustworthiness equals its writer's.** Without OQ-7, `settled_generation` records what a writer asserted. With OQ-7, it records what an enforced boundary could not have violated. This is not a defect of the field — it is why the two questions stay separate.

**S-15 — Interaction with OQ-5.** Adding fields changes the marker schema. Manifest versioning, older builds reading a newer marker, and the refusal path (`RefusedUnknownSchema`) are OQ-5's problem and must be settled before any implementation. Noted, not decided.

---

## 5. Decision-table effect, if adopted

**S-16 — Two rows change, and only these two:**

| Row | Now | With settlement evidence |
|---|---|---|
| 16 / 17 | `gate_enforced` ∧ content cannot predate enforcement → Healthy | additionally require `settled_generation ≥ G−1` |
| 18 / 19 | otherwise → Pending `attribution_unestablished` | unchanged |

**S-17 — Rows 16–17 remain inactive** until *both* the gate and the settlement field exist and the enforcement record (OQ-7) is positive. Their test coverage stays declared-inactive under capability `settlement_and_gate`.

**S-18 — No other row moves.** Generation 1 (rows 14–15) is untouched: a first commit has no preceding generation to settle.

---

## 6. Tests this would require

- **S-T1** Settlement is written only after both representations are observed valid and void.
- **S-T2** A `Recovered` marker with an older `settled_generation` refuses the transition that needs a newer one.
- **S-T3** The next generation's CAS compares `settled_generation`; a value one lower is refused, not accepted.
- **S-T4** No code path — scanner, migration or repair — writes `settled_generation`; an instrumented assertion.
- **S-T5** Death after both rewrites but before the CAS leaves the conversation retryable, never Healthy-with-content.
- **S-T6** A conversation whose settlement is recorded but whose enforcement is absent is Pending, and a conversation whose enforcement is proven but whose settlement is absent is Pending — the two-record rule, asserted both ways.

---

## 7. Recommendation

**S-19 — Adopt option A if OQ-6 is pursued:** a `settled_generation` field on the existing marker, CAS-compared, monotonic, never backfilled, never written by the scanner.

**S-20 — Do not pursue it before OQ-7.** A settlement record without an enforced boundary proves what a writer said, and the product would be entitled to no more confidence than it has today — while appearing to have more, which is worse than Pending.

**S-21 — Nothing here is authorized.** No marker field, no schema change, no migration, no backfill, no code. V3.1 remains the baseline; this is a design answer for the owner to accept, reject or redirect.
