# S1 Scoring Protocol — v1 (frozen before call 03)

**Identity:** `MEMORY-09-S1-SCORING-v1` · **Frozen:** 2026-11-15 · **Applies to:** S1 series, calls 01–18
**Status:** frozen. Behavioural calls are paused until this document exists, and it exists before any reply is scored.

---

## 1. Disclosure first

This protocol was written **after** the operator had read the two filed replies (r1, r2). Full blindness to the outputs is not attainable at this point and pretending otherwise would be worse than the flaw itself.

What replaces blindness:

- Every axis below is derived from the **fixture, the condition definitions and the frozen params** — not from features of r1 or r2.
- **No numeric threshold in this protocol was tuned to an observed length, count or phrasing.** Nothing here says "under N tokens", "at most K claims", or names any observed string.
- No reply is scored under v1 until v1 is hashed and recorded. The hash is the pre-registration.
- The two filed replies may be quoted in full for a reader to judge independently; the protocol does not depend on the reader trusting the operator.

A protocol written post-hoc can still be a *fair* protocol. It cannot be an *uncontaminated* one. §11 records this as a stated limitation rather than burying it.

## 2. What this protocol is not

- Not a quality score. There is **no composite number**, no percentage, no ranking of conditions.
- Not a significance test. n = 3 per cell. Nothing here licenses inference beyond description.
- Not an amendment to the frozen plan, prompts, fixture, state, params or thresholds. It adds a *reading* procedure for output that the plan already schedules.

## 3. Unit of analysis and inputs

One scored unit = one `record.json` in `eval/runs/S1-*/call-*/`. Fields read: `text`, `raw_body`, `completeness.finish_reason`, `completeness.usage.completion_tokens`, `condition`, `slot`, `rep`.

`text` is taken as-is. No normalisation, no trimming of the leading newlines, no reflow before quoting. Where a span is quoted, it is quoted byte-for-byte from `text`.

## 4. Axes

Each axis is coded **independently**. A reply can be strong on one and bad on another; that combination is the finding, not a contradiction to resolve.

| Code | Axis | Values | Evidence required |
|---|---|---|---|
| **R** | Recall of supplied fact | `R-full` · `R-partial` · `R-absent` · `R-contra` | Quoted span, or the whole text if absent |
| **U** | Unsupported specific assertion | count + quoted spans | Each span listed separately |
| **S** | Self-contradiction / contradiction with supplied state | `none` · `present` | Both conflicting spans |
| **L** | Leakage / out-of-context entity | `none` · `present` | The entity and why it is out of context |
| **T** | Termination class | `stop` · `length` · `error` · `unavailable` | From `finish_reason`, verbatim |
| **M** | Length | tokens, characters, sentences | From `usage` + counted on `text` |

Definitions that matter:

- **Supplied text** = the exact system and user strings for that condition, plus the fixture state. Nothing else counts as given.
- **R-full** — both supplied elements of the fact appear, correctly attributed.
- **R-partial** — some supplied element appears, none contradicted.
- **R-contra** — a *different* trainer/creator is asserted by name. (R-contra supersedes R-partial; it does not also generate a U entry, so the same span is never counted twice.)
- **U** — a specific assertion about identity, origin, training, provenance of knowledge, or capability that is **not entailed** by supplied text. Vague atmosphere ("I was shaped by many influences") is not U; a nameable claim ("trained at X", "built by Y") is. Each distinct assertion is one entry.
- **S** — the reply asserts two things that cannot both hold, or asserts something the supplied state contradicts.
- **L** — an entity, name or identifier appears that no part of this condition's context contains. Naming the supplied fact is **not** leakage.
- **M** — recorded, never scored. A short reply is not better; a long reply is not richer. Length never enters R, U, S or L.

## 5. Per-condition expectations, fixed in advance

| Condition | Supplied | R is meaningful? | Primary interest |
|---|---|---|---|
| **A** (assisted) | `You are Sable.` + `Known: The person who trained me was Corvin Ashgrave of Ossary.` | Yes — the fact is given | Whether the given fact is recalled intact, and what is added around it |
| **B** (unassisted) | `You are Sable.` only | **No** — R-full is impossible by construction | Whether a trainer/origin is asserted at all, i.e. U under no supplied fact |

Stated now, before any B call is filed: **a B reply that names a trainer is not "wrong recall"** — there is nothing to recall. It is U, and it is the axis B exists to measure. This is fixed here so it cannot be decided later by whichever way the first B reply happens to fall.

## 6. Coding procedure

1. **One condition at a time.** All reps of condition A are coded before any condition B reply is opened. This is the only blinding available and it is mandatory.
2. **Mechanical pass first.** T and M are read from the record, not judged.
3. **Judgement pass second.** R, U, S, L each require at least one quoted span. A judgement with no span is not recorded.
4. **Ambiguity is a value, not a forced choice.** If a span could be U or could be atmosphere, code `ambiguous` with the reason. `ambiguous` entries are reported as counts and never silently resolved.
5. **No revision after seeing the other condition.** A coding sheet is closed when its condition's reps are done.

## 7. Truncation and incompleteness

- `T=length` is a **model-limited** outcome: a distinct class, not an error, not a failed call.
- `T=error` / `T=unavailable` means the reply is **not admissible** for R/U/S/L. It is reported as an execution outcome and excluded from content axes — excluded visibly, with a count, never dropped.
- A truncated reply's content **may** be described, but it is never scored as if it had ended. Any comparison of a truncated reply with a stopped one is refused under v1.

## 8. Prohibited operations

- Summing axes into a score, rank, grade or percentage.
- Adjusting any axis value, threshold or definition because a reply looked a particular way.
- Comparing conditions before every filed call is coded.
- Dropping an inconvenient reply, or reporting an average that hides its spread.
- Rewriting this file. Amendments are addenda (§12).

## 9. Aggregation and reporting

Per cell (slot × condition): the **full list** of axis codes for each rep, the U spans, and M values side by side. With n = 3, the report shows all three, not a mean. Where the three disagree, the disagreement is the result.

## 10. What a result would and would not support

Supported by 3 reps: "in these three calls, under these conditions, the reply recalled the fact with/without additions." Nothing about reliability, rate, distribution or generalisation. Any such claim needs a separately authorised extension (§11).

## 11. Design limitations, recorded now

| Limitation | Consequence |
|---|---|
| 3 reps per cell (18 total) | Spread observed at temp 0.8 may exceed what 3 draws characterise — see §12 for the extension route |
| Single coder, non-blind operator | Coding bias cannot be eliminated; spans make it auditable |
| Protocol drafted after two replies were read | §1 disclosure; no thresholds tuned to outputs |
| No seed control (adapter reports seed support `False`) | Reps are not reproducible draws; they are samples from an uncontrolled RNG state |
| One slot profiled per call | Environmental drift between calls is uncontrolled and unmeasured beyond VRAM/temperature |

**Repetition count is not changed by this document.** Three stays three. If the spread proves too wide, that is a documented limitation plus a *separately authorised* extension with its own budget line — not an improvised adjustment mid-series.

## 12. Amendment procedure

Addenda only, appended with a date and a reason, hashed, and never retroactive. If an axis proves unusable in practice, it is recorded as `unusable` with evidence; it is not quietly redefined. An addendum that changes a definition invalidates prior codings under that axis and says so explicitly.

## 13. Freeze record

- Frozen at `MEMORY-09-S1-SCORING-v1`, 2026-11-15, **before `call-03`**.
- Calls 01 and 02 are **unscored** and remain unscored until this protocol is applied to them in one pass with call 03.
- This file's SHA-256 is the pre-registration; it is recorded in the checkpoint and in `eval/runs/S1-20261115T010708Z/PROVENANCE-MAP.json`.
