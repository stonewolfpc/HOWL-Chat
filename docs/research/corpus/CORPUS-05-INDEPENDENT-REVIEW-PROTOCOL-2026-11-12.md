# Independent Reference-Answer Review — Protocol

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** Review **protocol defined**; review **not executed**. No item has been reviewed. Authoring is **paused** until review completes (`P-5`).
**Applies to:** 45 items across `CORPUS-02`, `CORPUS-03`, `CORPUS-04`.
**Authority:** `CORPUS-01` §6.1 `P-3` and `P-5`; `MEMORY-09` §5.

---

## 1. Why this exists now, and not later

The corpus authors can no longer see their own items clearly. After forty-five items the author knows what each one is *supposed* to test, which is exactly the knowledge that prevents noticing when it tests something else.

Review is therefore not a polish pass. It is the step that decides whether these are test items or merely careful documents.

**Order is deliberate:** review of 45 items before authoring the remaining 35 (`P-5`). A corpus that grows faster than it is checked accumulates unexamined assumptions at the same rate as coverage.

---

## 2. What a reviewer receives — and what they do not

| Provided | Withheld |
|---|---|
| The item: scenario, turns, required and forbidden evidence | The author's intended answer or rationale |
| The reference answer and its five parts | The author's notes on what the item was designed to catch |
| The six scoring dimensions and the rubric anchors | Which batch the item came from |
| The restraint taxonomy (`R-1`–`R-6`) and the item's tag | The `tuning` / `heldout` split label |
| The failure-class definitions | Any prior reviewer's verdict, before the reviewer submits |

**The critical omission is the intent.** A reviewer told "this tests corroboration illusion" will find corroboration illusion in it. Told only "is the mandatory section supported by the evidence," they may find that it is not.

---

## 3. The five required challenges

Every item is answered against all five, per reviewer. Each answer is `pass`, `fail`, or `ambiguous`, with a one-line reason. `ambiguous` is a legitimate verdict and is never rounded to `pass`.

| # | Challenge | What a `fail` means |
|---|---|---|
| `C-1` | **Does the item actually test its named mechanism?** (for `R-6` items: the named epistemic mechanism; for others: the named failure class) | The item measures something adjacent. Rewrite or retire. |
| `C-2` | **Does the reference answer give away the author's preferred interpretation?** | The item rewards a specific phrasing rather than a correct behaviour. Strip the tell. |
| `C-3` | **Is the over-exposure boundary objectively scorable?** | Two competent raters would disagree about whether the boundary was crossed. Sharpen or convert to `none`. |
| `C-4` | **Can the item be passed through blanket silence or blanket abstention?** | The item rewards avoidance. It needs a perceptible requirement or it should be tagged as a silence-passing mode and checked against `A-1`. |
| `C-5` | **Does the recorded evidence support every supposedly mandatory fact?** | A mandatory part rests on something not in the history. This is the most serious failure: it makes the item unanswerable except by invention. |

**`C-5` is checked first**, because an item that fails it cannot be repaired by rewording — only by adding evidence or removing the fact.

---

## 4. Per-item record

```text
item_id
reviewer_id
review_date
C-1  : pass | fail | ambiguous   reason:
C-2  : pass | fail | ambiguous   reason:
C-3  : pass | fail | ambiguous   reason:
C-4  : pass | fail | ambiguous   reason:
C-5  : pass | fail | ambiguous   reason:
verdict      : accept | rework | retire
author_intent_guess: (optional) what the reviewer believes the item was meant to test
notes
```

**`author_intent_guess` is optional and recorded before any reconciliation**, because a reviewer's guess at intent is itself evidence: if it differs from the actual intent, the item is ambiguous even when every challenge passed.

---

## 5. Reconciliation — preserving the original verdict

The original verdict is **never edited**. Reconciliation is a separate record that sits beside it.

| Step | Rule |
|---|---|
| 1 | Two independent reviewers where feasible; one reviewer is **not** independence and is labelled as lower confidence |
| 2 | Any `fail` on `C-5` routes the item to **retire** unless the missing evidence is added and the item re-reviewed |
| 3 | Disagreement between reviewers is recorded as disagreement, with both verdicts intact |
| 4 | An item that cannot be reconciled is **retired**, not forced to consensus |
| 5 | A retired item keeps its ID and is marked `RETIRED` with the reason; IDs are never reused |
| 6 | The author may propose a rework, but the reworked item is **new evidence** and is re-reviewed from scratch |

**Rule:** consensus is not correctness. If two reviewers disagree because the item is ambiguous, agreeing manufactures confidence the item never earned.

---

## 6. Reconciliation record

```text
item_id
reviewer_a_verdict  (verbatim, preserved)
reviewer_b_verdict  (verbatim, preserved)
disagreement        : none | substantive | minor
resolution          : accept | rework | retire
resolution_reason
re_review_required  : yes | no
```

---

## 7. Independence requirements

| Requirement | Reason |
|---|---|
| The reviewer must not have authored the item | Familiarity with intent defeats `C-2` |
| The reviewer must not be told the split label | Held-out knowledge can bias leniency or severity |
| Conflicts of interest are declared | An author reviewing their own work is not a reviewer |
| Review is recorded, not summarised | A summary loses the disagreements that matter |

**Independence threshold (owner, 2026-11-12).** One non-human reviewer is **evidence, not validation**. Sealing requires convergence from **two** independent reviewers, preferably using different processes or model families. The second reviewer receives the corrected items **without** the first reviewer's verdicts, and a disagreement is preserved rather than resolved by the weight of a single report.

**Whole-record rule (owner, 2026-11-12).** A repair that changes one field must be reviewed against the **complete record**, not the changed field. The `F6-04` case is the model: `mandatory` was corrected while `Required`, `never scored`, and `Notes` still asserted the opposite, and the record passed every mechanical check at the time. Reviewers must read `Required`, `Forbidden`, `Variation`, `Prohibited`, `never scored`, `over-exposure`, `If failed`, and `Notes` together, and ask whether they describe **one** item.

**Rule `P-6` — registration is not inspection** (owner, 2026-11-12). Registering a file records its existence and hash; it says nothing about its content. A follow-up written from an instruction's description is not a substitute for reading the artifact the instruction refers to. Before acting on a reference, open it. This rule exists because a reviewer report sat filed and registered in this workspace for one cycle while a follow-up document asserted its findings were unavailable. **Scope:** the rule is general research governance, not a corpus convention. It applies to corpus files, source code, test reports, registry entries, hashes, and execution summaries alike. A hash establishes byte identity; a registry establishes recorded existence and location; neither establishes that anyone inspected the content or that a report accurately represents it.

**This protocol cannot be executed by the corpus author.** The items above were authored in this workspace, so an independent reviewer must be a different party — the owner, or a separately authorized agent given only §2 material. Self-review would satisfy the paperwork and fail the purpose.

---

## 8. Outputs

| Output | Use |
|---|---|
| Per-item verdicts (45 records) | Decides which items can be sealed |
| Retire list with reasons | Removes unanswerable or ambiguous items |
| `C-3` failures | Rewrites over-exposure boundaries that cannot be scored |
| `C-4` failures | Identifies items that reward avoidance |
| `C-5` failures | The serious category: mandatory facts without support |
| Spec-change proposals | Amendments to `CORPUS-01`, dated, never silent |

**Review does not edit the corpus.** It produces verdicts; edits are a separate, recorded step.

---

## 9. What review does not establish

- It does not establish that the corpus is adequate in size; 45 items remain below the floor.
- It does not establish that any configuration works; nothing has been run.
- It does not substitute for the pilot; variance is still unmeasured.
- It does not authorise `G2` or any model execution.

**After review:** seal the surviving items, then resume authoring toward the floor with the `R-2`/`R-3` balance target, then pilot.
