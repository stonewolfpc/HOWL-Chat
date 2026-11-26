# Round 3 instrument corrections

**Deliverable:** correction pass on the round-3 memory-extraction instrument,
responding to the Muse adversarial audit (`F-01`–`F-14`) plus one defect found
during this correction pass. Written 2026-11-19.

**Corrected instrument:** `experiments/mem-extraction/round3/instrument_v2/`,
version identifier **MEMX-EVAL 2.0.0**.

**Status:** implemented and unit-tested. **Not validated as a memory-quality
instrument by those tests.** Nothing was run against a model; nothing was
committed, staged, registered, deployed or integrated. No historical round-3
output, score or report was modified.

---

## 1. What was and was not done

**Done.** A separate evaluator that distinguishes output validity from semantic
content, measures record atomicity, separates a semantic composite from a
contract composite, treats ownership/scope contamination as a critical category,
accounts for one incident once, keeps truncation visible, and enforces claim
language. Thirty-two behavioural tests cover every case the brief requires.
A separate recalculation over the preserved runs was produced and written to
`instrument_v2/out/corpus_v2.json`.

**Not done.** No model inference. No new corpus run. No re-scoring of the
published v1 figures. No edit to
`docs/research/MEM-EXTRACTION-ROUND3-2026-11-19.md` — its overstatements are
recorded in Appendix C as errata with exact quotes, not rewritten in place.
No MEMORY-09 or production decision.

**Why the report was not edited.** It is a dated deliverable that other
documents may cite. Editing it in place would destroy the ability to tell what
was claimed when. Appendix C gives the correction for each sentence; applying
those corrections to the report requires the owner's approval (§9, item 1).

---

## 2. Finding disposition

| Finding | Severity | Disposition | Where |
|---|---|---|---|
| **F-01** atomicity unmeasured | High | **Corrected.** Per-record verdict `ATOMIC` / `LEGITIMATE_COMPOUND` / `MERGED_DISTINCT` / `UNRELATED` / `UNCERTAIN`; merged distinct facts earn half credit each and one atomicity violation, not full recall. Declared joint facts are accepted, so no arbitrary splitting rule is imposed. | `atomicity.py`, `rubric_v2.py` |
| **F-02** restraint scored as emptiness | High | **Corrected.** Explicit outcome taxonomy. Absence-based credit requires a valid status. A header-only or label-only table is `MALFORMED` and scores 0 on both composites. | `outcomes.py`, `rubric_v2.py` |
| **F-03** semantics conflated with format | Medium | **Corrected.** Content is matched in any text-bearing row, contract-invalid included. Decision correctness and contract validity are separate fields. Two composites; `production_usable` stays false for malformed output. | `matching.py`, `rubric_v2.py` |
| **F-04** fidelity + false-memory double penalty | Medium | **Corrected.** One event per incident with one `primary_dimension`. `false_memory` is no longer a second scored dimension; it appears as a secondary impact. Events are reported outside both composites. | `rubric_v2.py` |
| **F-05** thin samples; executive claim exceeds evidence | High | **Corrected in tooling, recorded as erratum for the report.** `reporting.sample_qualifier` and `reporting.compare` make a thin-sample comparison return `UNESTABLISHED` with the wording "no measured difference … not evidence that a difference does not exist". §1.1's "does not beat" is quoted and corrected in Appendix C. | `reporting.py` |
| **F-06** round-2 rebuild not like-for-like | High | **Documented, cause left unknown.** Exact configuration deltas tabulated (§4). This pass does **not** attribute the non-reproduction to variance or to a harness bug; both remain open. | this document §4 |
| **F-07** sampler ablations as point estimates | Medium | **Corrected in tooling.** `reporting.effect` returns a directional statement with n and range and refuses a causal reading. | `reporting.py` |
| **F-08** truncation folded into quality | Medium | **Corrected.** `truncated` flag, `INCOMPLETE` status, `quality_eligible=false`; aggregates report eligible and all-runs means side by side with a truncation count. Nothing is dropped. | `rubric_v2.aggregate` |
| **F-09** isolation proven narrowly | Low | **Preserved and narrowed.** `reporting.ISOLATION_ALLOWED` states only the serial `-np 1` observation and the HTTP 501; a prohibited-phrase list blocks the general reading. | `reporting.py` |
| **F-10** CPU/GPU not apples-to-apples | Low | **Preserved and qualified.** `reporting.TIMING_ALLOWED` records that wall time excludes load, that the flag sets differ, and that VRAM includes desktop consumers. | `reporting.py` |
| **F-11** rubric integrity; pre-registration unverified | Low | **Preserved as unverified.** `reporting.PREREGISTRATION` records status `UNESTABLISHED` and what evidence would be needed. No claim of pre-registration is made anywhere in the v2 material. | `reporting.py` |
| **F-12** TinyLlama inadmissible is genuine | Low | **Preserved.** Not rescored, not reclassified. `configs.py` and `score_triage.py` untouched. | — |
| **F-13** ablation controls clean; dedup gap remains | Low | **Preserved, gap open.** No ablation re-analysis was performed. The full `run_code` uniqueness audit across `ablation_corpus.json`, the per-label files and `ablation_corpus_all.json` is listed as outstanding (§8). | — |
| **F-14** C17 reclassification is a judgment change | Low | **Preserved.** v2 still treats a stored forbidden string as an event, now marked `policy_sensitive: true`, so the policy question stays visible instead of being silently decided. | `rubric_v2.py` |
| **D-15** negated forbidden phrase counted as stored | — | **Found in this pass; corrected.** `"double the toll (not triple)"` was a false positive in v1 (run `jjf67xzt`) and in the first v2 build. Negation-aware matching now returns no event, and 16 of 28 previously counted stored-forbidden incidents disappear. | `rubric_v2.forbidden_occurrences` |
| **D-16** one incident emitted as several events | — | **Found in this pass; corrected.** The C17 thinking run produced two events for one credential disclosure. Events are now aggregated per incident and list all phrases once. | `rubric_v2.stored_forbidden` |

Two audit statements required correction against the stored record, and are
reported rather than accepted:

- The audit says the 1.2B's correct C13 attribution "scores 17". The stored
  record shows draw 1 at **47.8** and draws 2–3 at **2.0**; 17.3 is the
  three-draw case mean, not a single run. Both facts are true, but the sentence
  is imprecise and is quoted in Appendix C.
- The audit says the round-2 rebuild changed `min_p` from "implicit 0.05" to
  explicit `0.0`. That is correct, and the reason is documented in
  `configs.py`: llama.cpp defaults `min_p` to 0.05, so leaving it unset silently
  changes the chain.

---

## 3. Files

**New (no existing file modified):**

```
experiments/mem-extraction/round3/instrument_v2/
  README.md                    package description and method
  outcomes.py                  validity taxonomy and parser
  matching.py                  content matcher
  atomicity.py                 atomicity judgement
  ownership.py                 ownership / scope contamination
  rubric_v2.py                 scoring, incidents, aggregation
  reporting.py                 claim-language rules and scanner
  evaluate.py                  CLI with write guards
  tests_v2.py                  32 tests
  cases_adversarial_v2.json    unrun ownership reference cases
  out/corpus_v2.json           separate recalculation (235 runs)
docs/research/MEM-EXTRACTION-ROUND3-INSTRUMENT-CORRECTIONS-2026-11-19.md
```

`git status` still reports the same six pre-existing modified files and no
modification under `experiments/mem-extraction/round3/`.

---

## 4. Round 2 versus round 3 configuration (F-06)

Round 2 is `experiments/lfm25-cpu/run_round2.py` importing `run_probe.py`.
Round 3 is `harness.py` with `configs.py`. Both send the **same prompt** — round
3 imports `PROMPT` from `run_probe.py` rather than retyping it — but almost
nothing else is equal.

| Parameter | Round 2 | Round 3 | Effect on comparability |
|---|---|---|---|
| context | `-c 2048` | `-c 4096` | different KV budget; changes truncation headroom |
| `min_p` | not sent → server default **0.05** | explicit **0.0** | removes a tail cut-off; changes the sampled distribution |
| `top_p` | 0.95 | 0.95 (0.8 for 4B non-thinking) | same for the 2B default chain, not for 4B |
| `top_k` | 40 | 20 | narrower candidate set |
| temperature | 0.8 | 0.7 / 1.0 by preset | different chain |
| repeat penalty | 1.0 | 1.0 | same |
| max tokens | 400 | 768–3072 | round 2 truncated; round 3 usually does not |
| `--jinja` | absent | present | template path differs |
| `--reasoning-format deepseek` | absent | present | separates `reasoning_content`; changes what `content` contains |
| `--metrics` | absent | present | measurement only |
| `-np 1` | present | present | same |
| thinking control | `chat_template_kwargs.enable_thinking` | same, plus `--reasoning-budget` where supported | budget is new |
| prompt | `PROMPT` | identical import of `PROMPT` | same |
| triage rubric | 12-point | identical 12-point | same |

**Conclusion, stated at the strength the evidence supports:** round 2's 9.33 and
9.00 did not reproduce under round 3's harness (3.67 and 6.67). The cause is
**unknown**. Variance at n=3 on a 12-point scale is plausible. A template or
harness effect is also plausible and cannot be excluded from these records.
Neither explanation is established, and no future document should present one as
the answer without a controlled re-run.

---

## 5. Before and after, on the preserved evidence

Recalculated with `MEMX-EVAL 2.0.0`; originals untouched. v1 and v2 totals are
**not on the same scale** and must not be substituted for one another. The
comparison shows what each defect did, not a "corrected score".

| Defect demonstrated | Run | v1 | v2 | What changed |
|---|---|---:|---:|---|
| Compound record rewarded (F-01) | `da7y2ahx` C01, 4B non-thinking d1 | 100.0 | semantic 75.0, contract 75.0 | four separately required rules in one row: recall 30 → 15, atomicity 25 → 0, one `ATOMICITY_MERGED_DISTINCT` event |
| Label table scored as restraint (F-02) | `frmuuwk5` C12, 1.2B d1 | 97.8 | semantic 0.0, contract 0.0, `MALFORMED` | three column-label rows, no memory; absence credit now requires a valid status |
| Right fact, broken format (F-03) | `luqe6lqe` C13, 1.2B d1 | 47.8 | semantic 43.1, contract 35.0, `MALFORMED` | content credit kept for `"The user is an ER nurse, not a healer"` (recall 10/30); usability refused |
| Contentless but well-formed | `zijgcap7` C13, 1.2B d2 | 2.0 | 0.0 | the model inverted the attribution (`"I was a healer, not an ER nurse"`); format validity does not rescue it |
| Negated forbidden phrase (D-15) | `jjf67xzt` C18, 4B d2 | 61.1, flagged catastrophic | semantic 75.0, fidelity 18/18, no forbidden event | `"exactly double the toll (not triple)"` no longer counts as storing "triple" |
| Credential disclosure (F-04, D-16) | `jo21wmi5` C17, 4B thinking d1 | 18.9 | semantic 50.0, contract 100.0, one event | both `hunter2` and `password is hunter2` reported once, `policy_sensitive: true` |
| Truncation visible (F-08) | all LFM3 thinking runs | folded into mean | `truncated=6`, excluded from the eligible mean, all-runs mean still reported | 6 of 36 runs were cap-truncated |

**Aggregate recalculation (separate from §9 of the report; different scale):**

| Condition | v1 mean | v2 semantic (eligible) | v2 contract | eligible / runs | truncated | production-usable |
|---|---:|---:|---:|---:|---:|---:|
| 4B non-thinking | 77.09 | 63.2 | 77.3 | 54 / 54 | 0 | 12 / 54 |
| 4B thinking 512 | 75.87 | 59.0 | 81.8 | 18 / 18 | 0 | 7 / 18 |
| LFM3 thinking 512 | 69.35 | 62.7 | 69.6 | 31 / 36 | 6 | 5 / 36 |
| 1.2B T0.1 | 27.22 | **not reported** | 13.7 | 0 / 54 | 0 | 0 / 54 |
| 1.2B T0.8 | 19.89 | **not reported** | 10.0 | 0 / 36 | 0 | 0 / 36 |
| TinyLlama | 21.67 | 42.2 | 19.1 | 19 / 36 | 17 | 1 / 36 |

The 1.2B rows are the clearest behavioural change: v1 produced a semantic-looking
number for a model that emitted no contract-valid row in any run. v2 reports no
semantic mean at all. That is the correct answer — the measurement does not
exist — and it is the single strongest justification for the rewrite.

Outcome census over 235 runs: `VALID_KEEP` 81, `MALFORMED` 89, `UNPARSEABLE` 40,
`VALID_NONE` 10, `VALID_HOLD` 9, `INCOMPLETE` 4, `VALID_REJECT` 2. Events:
`ATOMICITY_MERGED_DISTINCT` 76, `FORBIDDEN_CLAIM_STORED` 12 (down from 28 before
the negation fix).

---

## 6. Tests

`python -m unittest tests_v2 -v` → **32 passed, 0 failed**. Every test asserts an
expected outcome; none merely executes code.

| Required case | Test |
|---|---|
| four facts in four records | `test_four_facts_in_four_records_get_full_recall` |
| four facts merged into one record | `test_four_facts_merged_into_one_record_do_not_get_full_recall` |
| valid NONE | `test_valid_none_scores_full_on_a_restraint_case` |
| malformed header-only | `test_malformed_header_only_gets_no_restraint_credit` |
| correct fact, invalid format | `test_correct_fact_in_invalid_format_keeps_semantic_but_not_contract` |
| incorrect attribution | `test_incorrect_attribution_is_an_event_and_zeroes_the_dimension` |
| belief promoted to world truth | `test_character_belief_promoted_to_world_truth_is_critical` |
| two contradictory beliefs | `test_two_contradictory_beliefs_collapsed_into_one_claim` |
| private fact shared | `test_private_fact_shared_with_another_character` |
| truncated output | `test_truncated_output_is_flagged_and_not_quality_eligible` |
| missing output | `test_missing_output` |
| execution failure | `test_execution_failure_is_not_scored` |
| one incident, several dimensions | `test_one_incident_reduces_one_dimension_and_stays_visible` |
| additional | valid REJECT, valid HOLD, legitimate compound, off-topic padding, silence on a required case, negation, multi-phrase incident, scenario leak, unscored ownership, input immutability, four reporting-language tests, two write-guard tests |

Two of these failed on first run and the implementation was changed, not the
test: invalid-format content was being zeroed (defeating F-03), and a valid
`NONE` could not reach the usability threshold because schema and atomicity were
still being required of a response that has no rows.

**What the tests do not establish.** They show that the evaluator implements its
own rules on synthetic cases. They say nothing about whether those rules track
what a reader would call a good memory. No inter-rater reliability exists, and
the ownership detector is declarative rather than inferential.

---

## 7. Remaining limitations

1. **Atomicity is lexical.** Clause splitting plus a verb lexicon plus declared
   joint groups. It will miss a merge written as two sentences without a
   conjunction and can flag a long single fact. Low-confidence verdicts are
   listed in `atomicity_review`.
2. **Ownership is only detectable where declared.** Without an expectation
   block, ownership is unscored. The adversarial cases in
   `cases_adversarial_v2.json` are reference definitions, not measurements.
3. **`production_usable` is a strict conjunction** (contract ≥ 80, valid
   status, no truncation, no critical event). It produces 27/235 on historical
   runs, which is a statement about those outputs under this contract, not a
   deployability estimate.
4. **The semantic and contract weights are a design choice**, not a validated
   psychometric model. They are fixed in code and documented so that a later
   pass can challenge them explicitly rather than discover them.
5. **No confidence intervals are computed in v2.** The historical bootstrap CIs
   remain the only interval estimates, and they are over 18 cases.
6. **The v1 scorer is still the only scorer for triage.** The six-item instrument
   was not rewritten; the audit found it adequate and TinyLlama's inadmissibility
   genuine.
7. **The ablation dedup audit (F-13) is still open.** No `run_code` uniqueness
   proof across the four ablation JSON files was produced in this pass.

---

## 8. Open items that need a decision, not more code

1. **Apply Appendix C to the report in place, or publish it as an addendum?**
   Recommended: addendum. The report is a dated artifact; rewriting it erases
   the distinction between what was claimed and what is now known.
2. **C17 forbidden-string policy.** v2 keeps the v1 rule and marks the event
   `policy_sensitive`. §26's argument (retention is not the failure) is a
   policy judgment about MEMORY-09 authority, not a scoring defect. Changing the
   rule would change every C17 number. Owner decision.
3. **Should a three-field row (`ID | DECISION | MEMORY`, no REASON) be
   contract-valid?** v2 says no. This materially affects the 4B's C14 draws,
   which state the right content without a reason field.
4. **Ownership detection scope.** Declarative blocks only (current), or an
   attempt at free-text inference with a much higher false-positive rate?
5. **Re-run the corpus under v2?** Not authorized here and not needed to prove
   the instrument works. It is needed to make any *quality* claim on v2 terms.
6. **Publish a correction note for the v1 C18 false positive?** The historical
   score and catastrophic flag remain as published; the errata records the
   defect.

---

## 9. Preservation and scope confirmation

- **No historical round-3 evidence was modified.** `results/*.json`, `raw/`,
  `corpus.json`, `RUBRIC.md`, `SOURCES.md`, `qualitative/` and the report are
  byte-identical to their state before this pass. `evaluate.py` refuses to write
  into `results/`, `raw/` or `models/`.
- **No retroactive rescoring.** `results/corpus_scores.json` is unchanged. The
  v2 recalculation is a new file and is labelled as a different instrument.
- **No inference.** No `llama-server` was started; no model file was read; the
  only inputs were stored JSON and Markdown.
- **No staging, commit, push, registration or deployment.**
- **No Howl Chat, MEMORY-09, production or registry change.** No new MEMORY-09
  architecture was designed from the ownership examples; they are measurement
  requirements only.
- **Six pre-existing modified files remain the only tracked modifications**, as
  at the start of the session.

---

## Appendix A — validity taxonomy

| Status | Meaning | Absence credit | Usable |
|---|---|---|---|
| `VALID_KEEP` | ≥1 contract-valid KEEP row | yes | if contract ≥ 80 |
| `VALID_HOLD` | contract-valid rows, no KEEP, ≥1 HOLD | yes | if contract ≥ 80 |
| `VALID_REJECT` | contract-valid rows, all REJECT | yes | if contract ≥ 80 |
| `VALID_NONE` | literal `NONE` | yes | on restraint cases |
| `VALID_EMPTY` | explicit empty declaration, not the literal token | yes | on restraint cases |
| `MALFORMED` | headers/labels only, no data row | no | no |
| `UNPARSEABLE` | data rows, none satisfies the contract | no | no |
| `MISSING` | empty content, not a cap effect | no | no |
| `INCOMPLETE` | cap reached before a usable row | no | no |
| `EXECUTION_FAILURE` | the run failed | no | no |

## Appendix B — incident overlap

| Event | Primary dimension | Secondary (visible, not deducted) | Critical |
|---|---|---|---|
| `FORBIDDEN_CLAIM_STORED` | fidelity | false-memory, restraint | yes |
| `ATTRIBUTION_ERROR` | attribution | ownership | yes |
| `PROMOTION` | ownership | semantic | yes |
| `INHERITANCE` | ownership | semantic | yes |
| `CONTRADICTION_COLLAPSE` | ownership | semantic | yes |
| `SCENARIO_LEAK` | ownership | semantic | yes |
| `ATOMICITY_MERGED_DISTINCT` | atomicity (contract) | recall (halved per item) | no |
| `EXECUTION_FAILURE` | none | none | yes |

No event reduces two dimensions. Where a behaviour legitimately affects two, the
second is recorded as an impact string and is not subtracted.

## Appendix C — language errata (report left unchanged)

| § | Quoted text | Problem | Correction |
|---|---|---|---|
| 1.1 | "the thinking variant, at a 512-token runtime budget, does not beat it (see §9)" | n=1 cannot establish "does not beat" | "showed no measured gain at n=1, at roughly 14× the wall time" |
| 1.2 | "so that one parameter costs about three points" | point estimate on n=3, 12-point scale | "was associated with a 3-point difference in this session; not an effect size" |
| 21.4 | "The likely cause is high variance at T0.8 with n=3" | attributes cause | "the cause is unknown; variance is plausible, so is a harness or template effect" |
| 16.1 | "Per-request state **is** cleared by the runtime" | mechanism inferred from output identity | "each job produced output identical to a fresh process, and `cache_n` was 0; no reset API was exercised (HTTP 501)" |
| 13 | "Both fit fully in 8 GB." | card-wide claim from two models | "Both fit in 8 GB on this card in this configuration." |
| 24.4 | "There is no niche for it." | absence of evidence read as absence | "No niche was measured for it on this corpus at n=2." |
| 25 | "the same model's one correct attribution … scores 17" | 17.3 is the 3-draw case mean; draw 1 is 47.8 | "the case mean is 17.3 across three draws; the best draw, which states the attribution correctly, is 47.8" |
| 27 | "77.19 vs 77.09 is noise" | "noise" asserts equivalence at n=2 | "not distinguishable at n=2; the difference is within the observed spread" |
| 27 | "At 4B, uncensoring is free on this task." | "free" is a capability claim | "No capability loss was measured at 4B on this task at n=2." |
| 20 | "Pareto set" | depends on v1 means that F-01/F-02 inflate | retain as a cost-aware shortlist, not a Pareto proof |

Scan result over the report: 12 hits — 1 `PROOF`, 1 `SELECTION`, 10 `RANKING`.
The `PROOF` hit is a negated sentence ("This does **not** prove …") and is a
false positive of the scanner, recorded here so the count is not misread.
