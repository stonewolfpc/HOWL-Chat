# Memory Evaluation Corpus — Specification and Index

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**
>
> **Memory should change what the model understands, not force the model to announce what it remembers.** — owner, north star for this evaluation
>
> A memory system must be evaluated not only on what it retrieves, but on what it **changes**, what it **leaves untouched**, and what it **refuses to pretend it knows**. — owner

**Status:** `G1` execution — corpus construction. **45 of ≥80 items authored (batches 1–3). Authoring PAUSED pending independent review** (`P-5`). No model executed, no source changed, no schema changed.
**Review:** `docs/research/corpus/CORPUS-05-INDEPENDENT-REVIEW-PROTOCOL-2026-11-12.md` — protocol defined, **not executed**.
**Owner disposition of batch 1:** accepted as a **format-validation milestone** only — not a completed corpus, and not evidence supporting any architecture choice (owner, 2026-11-12).
**Authority:** `docs/architecture/OWNER-APPROVAL-MEMORY-09-AND-G1-2026-11-12.md` §4 (G1 bounds); `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` §3 (corpus requirements).
**Items:** `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md`, `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md`, `docs/research/corpus/CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` (review protocol: `CORPUS-05`)

---

## 1. What this corpus is for

It is not a demonstration that memory works. It is a set of test items designed so that a configuration can **fail** on them — including failure by remembering too much, too accurately, or at the wrong moment.

Three consequences follow, and they are the reason this file exists before the items:

1. **An item without a reference answer is not an item.** It is removed, not judged later.
2. **Every item names its required and forbidden evidence.** Scope and abstention are properties of the item, not opinions formed after reading a reply.
3. **The four principle cases are mandatory members** (`P-01`–`P-04`). A corpus of "did it recall X" items would reward exactly the behaviour the governing principle forbids.

---

## 2. Item record schema

```text
item_id                     F<family>-<nn> | P-<nn>
family                      F-1 … F-9 | P
title
split                       tuning | heldout
principle_case              (P items only: less-is-better | nothing-is-correct |
                             accurate-but-contextually-wrong | natural-use)
scenario                    world / characters / persona / history shape
conversation_sketch         the turns that carry the information
required_evidence[]         what must exist and reach the reply
forbidden_evidence[]        the leakage set — must never reach this reply
expected_abstention         true | false
expected_character_perspective
reference_answer            exact | acceptable_range | expected_abstention | expected_correction
pre_registered_diagnostics  failure classes expected if the item is failed
notes                       authoring caveats
```

### 2.1 Item ID scheme

`F<family>-<nn>` for family items; `P-<nn>` for principle cases. IDs are permanent. A retired item keeps its ID and is marked `RETIRED` with a reason; IDs are never reused, because a reused ID silently merges two test histories.

### 2.2 Split discipline (`R-11`)

| Rule | Statement |
|---|---|
| `S-1` | Every item declares `tuning` or `heldout`. |
| `S-2` | The held-out set is hash-sealed before any tuning begins. |
| `S-3` | The held-out set is not read by whoever tunes a configuration. |
| `S-4` | Final-corpus target: **≥40% held out**, **≥2 held-out items per family**. |
| `S-5` | Moving an item from held-out to tuning after tuning began is a **contamination event** and is reported as one. |

Current batch-1 split: 14 tuning / 8 held-out (36%). **Below the `S-4` target and declared as such** — the target applies to the final corpus, not to this batch.

### 2.3 Reference-answer types

| Type | Use | Requirement |
|---|---|---|
| `exact` | Name, number, factual event | Canonical answer plus allowed variants |
| `acceptable_range` | Open creative response | Explicit must-mention / must-not-claim criteria |
| `expected_abstention` | Insufficient, unauthorized, or irrelevant | The **reason** is stated per item |
| `expected_correction` | Superseded information | Must state current truth **and** that it changed |

### 2.4 Reference-answer sufficiency (required from batch 2 onward)

A reference answer must separate four things, or it will score style as fact:

| Part | Meaning | Scored? |
|---|---|---|
| `mandatory` | Facts or behaviours that must be satisfied | **Yes** |
| `variation` | Equally correct phrasings and behaviours | No — all pass |
| `never scored` | Length, tone, dialect, humour, confidence, warmth | **No — never a failure signal** |
| `prohibited` | Definite failures | **Yes** |

**Rule:** a creative item that scores a stylistic difference as a factual failure is a **defective item** and is reworked, not reinterpreted. Batch 1 items predate this section; each requires an **annex pass** to add the four parts before the corpus is sealed.

### 2.5 Minimum sufficient response (owner review question, 2026-11-12)

`mandatory` is a **floor, not a checklist to be recited**. An item written as "must mention X, Y, Z" teaches a configuration to say more, which is the behaviour the corpus exists to catch.

| Concept | Definition |
|---|---|
| **Mandatory** | The least that satisfies the item. Never "everything true about it". |
| **Over-exposure** | Additional true information that, in this moment, constitutes failure |

**Rule `M-1`:** every item states its over-exposure boundary — what extra, *true* material would fail it — or declares `over_exposure: none`. **Rule `M-2`:** for restraint items (`P-01`, `F6-04`, `F1-01`, `P-04`), over-exposure is a **scored** part, not a note. A response can satisfy every required fact and still fail by saying too much.

This closes the asymmetry that made the rubric gameable: before `M-1`, a verbose configuration was never worse than a concise one on `mandatory`, only on style — and style was explicitly unscored.

---

## 3. The four principle cases

These exist because every other item type can be passed by injecting more. These cannot.

| Case | Item | The failure it catches |
|---|---|---|
| **Less is better** | `P-01` | A correct, complete memory block that wrecks the scene |
| **Nothing is correct** | `P-02` | Confident fabrication where abstention was the only right answer |
| **Accurate but contextually wrong** | `P-03` | A true fact leaked to a character, persona, or moment that must not have it |
| **Natural use without recitation** | `P-04` | Correct recall delivered as a status report instead of behaviour |

**Scoring note:** `P-01`–`P-04` are scored on the six dimensions like any other item, but each also carries a single **pass/fail gate** on its named failure. A configuration that passes `P-01` cannot be credited for comprehension of the same scene.

### 3.1 Restraint is not silence (owner, 2026-11-12)

Six distinct behaviours are routinely collapsed into "restraint". They are not interchangeable, and a corpus that only ever rewards saying nothing teaches **avoidance**, not judgment.

| Mode | Behaviour | Pass condition | Can silence pass it? | Current items |
|---|---|---|---|---|
| `R-1` | **Suppress mention** — knows the fact, chooses not to name it | Fact unreferenced; nothing implies it was consulted | Yes | `P-03`, `F9-01` |
| `R-2` | **Influence tone** — knowledge changes manner, not content | Behaviour reflects it; the fact is never stated | **No** — a flat, unchanged reply fails | `F1-01`, `F1-02`, `F6-01` |
| `R-3` | **Natural use** — fact used as though simply remembered | Used without announcing that it was remembered | **No** — omission fails | `P-04`, `F6-03` |
| `R-4` | **Recognise irrelevance** — the fact does not apply here | No callback; an ordinary reply | Yes | `F8-01`, `F8-03`, `F9-01` |
| `R-5` | **Respect entitlement** — withholds because the recipient may not know | No leakage, no hint of knowing | Yes | `F2-01`, `F9-02`, `F9-05`, `P-03` |
| `R-6` | **Preserve uncertainty** — refuses to manufacture certainty | Uncertainty is preserved — no false precision is manufactured | **Yes** — silence is an acceptable way to preserve uncertainty; see `A-1` | `F7-01`, `F9-04` |

**Errata — 2026-11-13 (`AM-01`, authority `OD-M09-01`).** `R-6` is **silence-satisfiable**. The pre-amendment row read: pass condition "Uncertainty is conveyed", silence "**No** — saying nothing also fails"; both are superseded above. `A-1` is unchanged and is now the sole anti-avoidance constraint for `R-6`. `check_corpus.py` drops `R-6` from `POSITIVE_MODES` in the same instrument. No item, tag, score or result is changed — `C-1` reported 0 failures before the amendment.

**Rule `A-1` — restraint ≠ silence.** A configuration that passes `R-1`, `R-4`, `R-5`, `R-6` while failing `R-2` and `R-3` has learned **avoidance**, and is reported as avoidance — not as a restraint result. This is checked explicitly: the scorecard reports the six modes separately and refuses to aggregate them into one restraint figure.

**Rule `A-2` — every item is tagged with exactly one mode** (or `none`, if it is not a restraint item). Untagged restraint items are not accepted.

**Rule `A-3` — modes are independent of quality dimensions.** `R-2` can fail while `R-5` passes: a reply may leak nothing and still ignore the emotional context. `R-5` can fail while everything else passes. Each result stands alone; no mode excuses or compensates for another, and the same applies across the six scoring dimensions. The report is a matrix, not a score.

**Operational membership (tag-derived, 2026-11-13).** The item lists in the mode table above are retained as a **historical record** and are superseded for operation by the per-item `Restraint:` tag (`OD-M09-05`). Current tag membership:

| Mode | Tagged items |
|---|---|
| `R-1` | `F6-02`, `P-01` |
| `R-2` | `F1-01`, `F1-02`, `F1-03`, `F6-01` |
| `R-3` | `F3-01`, `F6-03`, `F6-04`, `P-04` |
| `R-4` | `F8-01`, `F8-03`, `F9-01`, `F9-03` |
| `R-5` | `F2-01`, `F2-02`, `F2-03`, `F5-01`–`F5-04`, `F9-02`, `F9-05`, `P-03` |
| `R-6` | `F7-01`, `F9-04`, `R6-01`–`R6-04` incl. `R6-03a`–`d` |

Where the two disagree, the tag governs. **Neither representation establishes that an item substantively exercises its mode** — see `MEMORY-09-RCN-1` §3 on `R-1`.

**Note on duplicated memberships.** This table's item lists contain duplicate family memberships that conflict with the currently governing `A-2` single-tag convention. The available sources do not establish when those memberships originated, whether they reflected an intentional earlier policy, or whether they were carried forward inadvertently. The table is retained as historical evidence, not as operational membership authority. An item appearing under two modes therefore records that both were considered, **not** that it carries two modes; `A-2` governs, and exactly one tag applies per item.

**Coverage (after batch 3, item-heading convention):** `R-1` 2 · `R-2` 4 · `R-3` 4 · `R-4` 4 · `R-5` 10 · `R-6` 6 · `none` 12 = **42 item headings**, corresponding to **45 scored units** (see the accounting basis in §4).

**`R-6` is now covered by four distinct uncertainty mechanisms**, not four variations of one (`CORPUS-04`): conflicting accounts requiring a decision, an incomplete timeline neither explanation can exclude, a corroboration illusion with a single origin, and a challenge that does not settle the accepted account. Four repetitions of "two versions, nobody knows" would have measured one behaviour four times.

**`R-6` now separates four provenance cases** (`R6-03a`–`R6-03d`): three independent firsthand observations, three derived from one origin, one source repeated, and unknown relationship. `R6-03a` is the **anti-hedging control** — hedging is the failure there, which is what prevents universal "I don't know" from passing the mode.

**`R6-03d` states the rule the whole mode depends on:** unknown independence is **not** established corroboration. An incomplete provenance record must be preserved as incomplete, never silently completed.

**Remaining balance target:** `R-2` and `R-3` still want **two more each** in the remaining items. `R-6` no longer needs additions.

**Gate rule `G-1` — no aggregate offsets a gate.** A configuration that leaks one character's private knowledge to another fails, regardless of its scores elsewhere. Blocking-gate items are listed per batch; a gate failure is reported as a gate failure in the first line of the result, before any score.

---

## 4. Coverage — batches 1–2

| Family | Items | Tuning | Held-out | Requirement (final) | Met? |
|---|---:|---:|---:|---:|---|
| `F-1` Long-running relationships | 3 | 1 | 2 | 10 | Partly |
| `F-2` Multi-character boundaries | 3 | 1 | 2 | 10 | Partly |
| `F-3` Long conversations | 4 | 2 | 2 | 10 | Partly |
| `F-4` Changing circumstances | 3 | 1 | 2 | 10 | Partly |
| `F-5` World / scenario changes | 4 | 2 | 2 | 8 | Partly |
| `F-6` Casual conversations | 4 | 2 | 2 | 8 | Partly |
| `F-7` Contradictory information | 11 | 3 | 8 | 8 | **Yes** (surplus) |
| `F-8` Irrelevant information | 4 | 2 | 2 | 8 | Partly |
| `F-9` No-retrieval cases | 5 | 3 | 2 | 8 | Partly |
| `P` Principle cases | 4 | 2 | 2 | 4 (mandatory) | **Yes** |
| **Total** | **45** | **19** | **26** | **≥80** | **56%** |

**Accounting basis (revised 2026-11-13).** Three counts are in use and are **not interchangeable**:

| Count | Value | What it counts | Excludes |
|---|---:|---|---|
| Item headings | **42** | One heading per family-level item | The four `R6-03a`–`d` provenance variants, which sit under the `R6-03` family heading |
| Scored units | **45** | Units carrying a complete scored body | The `R6-03` container, which has no `Required` or `Forbidden` |
| Heading lines / tags | **46** | Every `##`/`###` heading line, each carrying a `Restraint:` tag | Nothing — this is the raw line count, including the container and the four variants |

The §4 table and §7's "scored units" and ratio lines count **units**. The §5 coverage line counts **item headings**. Tag totals exceed scored units by one because the unscored `R6-03` container is itself tagged.

*Superseded:* this note previously read "42 item *headings*, 45 scored *units* — `R6-03a`–`d` are scored separately but share a heading", and recorded that the earlier `18 / 27` split was derived from batch totals rather than per-item tags and was wrong (`CORPUS-06` §2). Both statements remain accurate as history; the `18 / 27` correction is unaffected by this revision.

**Held-out accounting (resolved 2026-11-13, `AM-02`).** `check_corpus.py` measures **19 tuning / 26 held-out across 45 units**, matching the totals above (`TOTALS MATCH`). Each of the four `R6-03` variants now carries its own `**Split:** heldout` line, so all four are counted — 26 − 22 = 4.

*Superseded note, preserved verbatim:* "**Open discrepancy:** `check_corpus.py` measures **19 tuning / 22 held-out** — the four `R6-03` variants carry no per-unit split line. The `26` above assumes inheritance from the parent heading, which is written nowhere. Not reproducible → not compliant. See `CORPUS-06` §4.8." Its diagnosis was accurate when written and is no longer true.

*Basis disclosure:* the per-family figures in this table were re-derived by assigning units to families from their item-ID prefixes (`F7-*`, `R6-*` → `F-7`), because batch-1 items carry no `Family:` field. `F-9`'s row above was transposed (`2 / 3` → `3 / 2`); its units are `F9-01` tuning, `F9-02` tuning, `F9-03` held-out, `F9-04` tuning, `F9-05` held-out. The `Total` row was already correct and is unchanged.

**Split targets met at 47.5% corpus size:** ≥40% held out (**55%** actual) and ≥2 held-out items per family (**met in every family**). The item-count floor is still unmet, and `F-9` remains the thinnest family at 5 of 8.

**This is a numerical target, not a quality claim** (`P-1`). Whether the held-out items are valid, leak-free, and independent is settled by the `P-3` review, which has not happened.

`MEMORY-04` fixtures `E-01`–`E-12` are retained as **required members** and are not counted in the 38; they are scenario specifications, and their conversion into scored items is part of the remaining G1 work.

---

## 5. What G1 still requires

| # | Remaining work | Why it cannot be skipped |
|---|---|---|
| 1 | Author ≥38 more family items to reach the floor (42 of 80) | Coverage |
| 1a | Annex pass: add the four-part reference structure to all 22 batch-1 items | §2.4 |
| 1b | Bring `F-9` from 5 to 8 items | Abstention is the least-covered capability |
| 2 | Convert `E-01`–`E-12` into scored items with reference answers | Required members |
| 3 | Add the `forbidden_evidence` set to every existing fixture | Leakage must be testable |
| 5 | Independent second-pass review of every reference answer | An item whose answer is arguable is not a test item |
| 6 | Seal the held-out set by hash | `S-2` |
| 7 | Pilot on 16 items to estimate variance | `MEMORY-09` §10 |

**Batches 1–2 exist to prove the item schema is workable and that the principle cases can be authored at all.** It is not a corpus that can support a decision, and nothing here may be reported as a finding.

---

## 6. Authoring rules

1. **Synthetic only.** No user chat, character, persona, world, or lorebook is used. Any resemblance is incidental.
2. **One testable thing per item.** If an item tests two things, split it.
3. **The forbidden set is mandatory.** An item with no leakage risk is testing recall only.
4. **No item may require more context than the arm is allowed to carry**, unless the item is explicitly a budget-stress item and labelled as one.
5. **The reference answer is written before the item is accepted**, and never edited after a run.
6. **Neutral naming.** No item is named after the behaviour it forbids.
7. **Over-exposure boundary stated** (`M-1`), or explicitly `none`.

### 6.2 Required-field check — IMPLEMENTED

`docs/research/corpus/check_corpus.py` enforces the check. Run it before announcing a batch and before commissioning a review:

```bash
python docs/research/corpus/check_corpus.py
```

It verifies, **per scored unit**: unique ID at unit level, `Split`, `Restraint`, required evidence, forbidden evidence, and all five reference-answer parts including `over-exposure`. It also recomputes the split and total figures and reports whether the `CORPUS-01` declared totals are **reproducible from the files**.

The check has two layers: **presence** (all required fields) and **consistency** — `C-1` fails a perceptible restraint mode (`R-2`, `R-3`, `R-6`) whose `mandatory` is purely negative; `C-2` warns on a proper noun used in a scored field but absent from the scenario; `C-3` warns on a restraint item declaring `over-exposure: none`.

**The check reports and never repairs** (owner, 2026-11-12). A tool that silently fixed records would erase the difference between authored content and corrected content — which is the distinction the review exists to preserve. Exit code 1 on any defect.

**Why it is required rather than recommended:** the 15 missing boundaries and the false `38 / 38` claim would both have been caught in seconds. The failure was declaring compliance without a mechanism able to contradict the author.

### 6.1 Protections during G1 (owner, 2026-11-12)

| # | Protection | Consequence of forgetting it |
|---|---|---|
| `P-1` | **The held-out percentage is not a quality claim.** 55% satisfies a numerical target; it says nothing about whether those items are valid, leak-free, or genuinely independent. | A number mistaken for a verdict |
| `P-2` | **The batch-1 annex is a blocker.** The four-part standard must apply to all 22 earlier items before sealing. | Two reference-answer standards in one corpus |
| `P-3` | **Review must be genuinely independent.** The reviewer receives the item, its evidence, and the scoring rules — **not** the author's intended answer. Disagreements are recorded, never silently harmonised. | Author assumptions becoming ground truth |
| `P-4` | **Negative outcomes are preserved.** If a configuration over-retrieves, leaks, or manufactures certainty, the result stays. The corpus is not tuned until the failure disappears. | The corpus quietly becoming a demonstration |
| `P-5` | **Review outranks item count.** The next quality milestone is independent review of the 38 existing items, not reaching 80. A large unreviewed corpus is a liability, not progress. | Volume mistaken for rigor |

**`P-3` mechanics:** review is conducted blind to author intent; each reviewer records `agree` / `disagree` / `ambiguous` per item, and `disagree` items enter a recorded reconciliation with the original verdict preserved. An item that cannot be reconciled is **retired**, not forced.

---

## 7. Status

- Item headings: **42** · scored units: **45**
- Items with complete reference answers: **45 / 45**
- Items with four-part reference answers: **45 / 45** (batch 1 annexed 2026-11-12)
- Items with forbidden sets: **45 / 45**
- Held-out: **26 / 45 (58%)**
- Blocking-gate items: **9**
- Families at target: **1** (`F-7`, 11 of 8 — surplus retained)
- Batch-1 annex pass (`P-2`): **complete** — 22/22 items now carry the four parts and an over-exposure boundary
- Over-exposure boundary stated: **30 / 45** — **15 missing** (all batch 2 except `F6-04`). **Previously reported as `38 / 38` in error**; corrected 2026-11-12 per `CORPUS-06` §2
- Independent review (`P-3`): **not started** — the next quality milestone (`P-5`)
- Restraint modes tagged: **45 / 45** (`A-2`)
- Restraint balance: **`R-6` resolved** — eight dedicated fixtures (`R6-01`, `R6-02`, `R6-03` with variants `R6-03a`–`d`, `R6-04`), of which **seven are scored units** (the `R6-03` container is a family heading and is not scored); **ten tagged headings** in total, the eight fixtures plus legacy members `F7-01` and `F9-04`. `R-2`/`R-3` still want two more each
  - *Superseded record:* an earlier figure read "9 items, 4 mechanisms + 4 provenance variants"; its arithmetic yields 8, and the ninth element is not identified in any source.
- Review (`P-3`): **protocol defined, not executed** — the gate on sealing
- Authoring: **paused** until review completes (`P-5`)
- Held-out sealed: **no** (sealing happens once the final corpus exists)
- Model executed: **no**
- Results claimed: **none**
