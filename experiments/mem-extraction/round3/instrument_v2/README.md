# MEMX-EVAL 2.0.0 — corrected instrument

Separate, versioned evaluator for the round-3 memory-extraction evidence. It
exists because an independent audit found that the v1 scorer measured keyword
presence and format together. It does not replace, reinterpret or rescore
anything in `../results/`; the v1 numbers stay exactly as published.

**No inference is performed by anything in this directory.** It reads stored
JSON and computes.

## Layout

| File | Role |
|---|---|
| `outcomes.py` | output-validity taxonomy and parser (F-02) |
| `matching.py` | stem-tolerant content matcher; decision and contract kept separate (F-03) |
| `atomicity.py` | per-record atomicity judgement (F-01) |
| `ownership.py` | world/scenario/character/disclosure scope and contamination (task C) |
| `rubric_v2.py` | semantic and contract composites, incident accounting, aggregation (F-04, F-08) |
| `reporting.py` | permitted/forbidden claim language and a document scanner (F-05, F-07, F-09, F-10, F-11) |
| `evaluate.py` | CLI with write guards |
| `tests_v2.py` | 32 behavioural tests, including all cases required by the brief |
| `cases_adversarial_v2.json` | ownership/atomicity reference cases — **unrun**, no model output |
| `out/corpus_v2.json` | separate recalculation over the preserved round-3 runs |

## Run

```
python -m unittest tests_v2 -v
python evaluate.py --runs ../results/corpus.json --corpus ../corpus.json \
                   --out out/corpus_v2.json
```

The CLI refuses to write inside `results/`, `raw/` or `models/`, and refuses to
overwrite an existing output without `--force`.

## Outcome taxonomy

`VALID_KEEP`, `VALID_HOLD`, `VALID_REJECT`, `VALID_NONE`, `VALID_EMPTY`,
`MALFORMED`, `UNPARSEABLE`, `MISSING`, `INCOMPLETE`, `EXECUTION_FAILURE`.

`VALID_REJECT` is an addition to the minimum set in the brief: a REJECT-only
answer is a legitimate decision list and would otherwise have to be called
malformed. Truncation is a flag as well as a status; `quality_eligible` is false
for every truncated run, and aggregates report both the eligible and all-runs
figures.

## Two composites

| Semantic (what was understood) | Weight | Contract (is it usable) | Weight |
|---|---:|---|---:|
| recall | 30 | parseability | 30 |
| fidelity | 18 | four-field schema | 25 |
| significance | 12 | atomicity | 25 |
| temporal | 10 | termination | 10 |
| attribution | 10 | uniqueness | 10 |
| ownership / scope | 15 | | |
| uncertainty | 5 | | |

Restraint cases use `restraint 60 / fidelity 25 / ownership 15`. A dimension
that does not apply to a case is excluded and the remaining weights rescale; it
is never awarded for free. `production_usable` requires contract ≥ 80, a
decisional or valid-empty status, no truncation, and no critical event.

## One incident, one deduction

Each failure is an event with one `primary_dimension`. Only that dimension is
reduced; other effects appear in `secondary_impacts` and are not deducted again.
A forbidden claim stored in a KEEP row, for example, zeroes `fidelity` once and
emits one event, no matter how many forbidden phrases or rows describe it.

## What this instrument does not do

- It cannot infer ownership from free text. Ownership is scored only when a case
  supplies a declarative expectation block; otherwise it reports `unscored`.
- Atomicity is lexical (clause boundaries, verb lexicon, declared joint groups)
  and marks low-confidence merges for review.
- It does not establish memory quality. Passing tests establish that the
  evaluator implements its own rules, nothing more.
- It does not select a model, change MEMORY-09, or touch production.
