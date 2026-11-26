# Round 3 quality rubric — pre-registered 2026-11-19

Written **before** any round-3 generation was produced, and before the winning
configuration is known. The 12-point triage rubric from rounds 1–2 is kept
unchanged for historical comparability; this 100-point rubric is the expanded
standard for the corpus. Neither replaces the other.

## 1. What is scored

Two separate instruments:

| Instrument | Input | Scale | Comparability |
|---|---|---|---|
| **Legacy triage score** | the original six-item prompt, unchanged | 12 points (2/1/0 × 6) | identical to rounds 1–2 |
| **Corpus quality score** | the 18 cases in `corpus.json` | 100 points per case | new in round 3 |

The legacy score is reported for every configuration that was also run on the
triage prompt. It is **not** used to select finalists on its own.

## 2. Corpus dimensions and weights

Weights are fixed now and are not changed after seeing results.

| # | Dimension | Weight | What earns the points |
|---|---|---|---|
| 1 | Critical-detail recall | 20 | Each required fact present, with the right decision |
| 2 | Factual fidelity | 20 | Meaning preserved; no distortion, no reversal, no invented specifics |
| 3 | Contextual significance | 15 | Correct KEEP/REJECT/HOLD and a reason that names the actual ground |
| 4 | False-memory prevention | 15 | No fabricated entity, cause, name, number or certainty |
| 5 | Temporal accuracy | 10 | Corrections and superseded values represented as such |
| 6 | Uncertainty handling | 5 | "Maybe canon" stays uncertain; HOLD where confirmation is needed |
| 7 | Attribution | 5 | Speaker/character/real-user attribution preserved |
| 8 | Appropriate restraint | 5 | Nothing durable extracted from chatter or voided narration |
| 9 | Completeness | 3 | Details needed for later use are retained, not just headlines |
| 10 | Structured-output reliability | 2 | Parseable, one line per candidate, no preamble, no truncation |

**Repeatability is reported separately** as a consistency measure (mean
per-case standard deviation across draws, and the fraction of cases whose
decision vector is identical across all draws). It is not folded into the 100.

## 3. How the 100 is computed per case

1. **Recall (20).** `20 × (found required items / total required items)`. For a
   case with no required items (C09, C12), recall is awarded in full only if no
   KEEP appears; otherwise it is 0.
2. **Fidelity (20).** Start at 20. −6 per required fact whose meaning is
   distorted but still recognisable; −20 if a required fact is reversed
   (e.g. "can swim", "wards hold at dusk"). Floor 0.
3. **Significance (15).** Per required item: 3 points for the correct decision,
   plus up to 2 for a reason that names the real ground (durability,
   transience, contradiction, uncertainty, world-rule). A correct decision with
   a canned or absent reason gets 3.
4. **False memory (15).** Start at 15. −5 per fabricated atomic fact or
   unsupported causal claim that appears in a KEEP/HOLD line; −15 for a
   fabricated named entity or number. Floor 0.
5. **Temporal (10).** Only cases C05, C06, C07, C10 are temporally loaded. For
   them: 10 if the new/current value and the superseded/conflicting value are
   both represented correctly; 5 if only the current value is kept and the
   history is silently dropped; 0 if the obsolete value is presented as
   current. Cases with no temporal content receive 10 automatically and are
   marked `n/a` in the per-case table.
6. **Uncertainty (5).** 5 if every uncertain claim is marked uncertain or held;
   2 if the claim is kept but the uncertainty is mentioned in the reason; 0 if
   uncertainty is erased.
7. **Attribution (5).** 5 = every attributed item has the right speaker;
   0 = any required attribution is wrong; 2 = attribution absent but not wrong.
8. **Restraint (5).** 5 = no durable candidate from chatter; −2 per
   over-extracted transient item, floor 0. Cases C09 and C12 are pass/fail:
   5 or 0.
9. **Completeness (3).** 3 = all required facts plus their necessary qualifiers
   (numbers, names, relationships); 1 = facts present but qualifiers lost.
10. **Structure (2).** 2 = every line parses as `ID | DECISION | MEMORY |
    REASON` (or a clean `NONE`); 1 = parseable with noise; 0 = unusable.

Case score = sum, capped at 100. Corpus score = mean of the 18 case scores,
reported with per-case and per-category breakdowns.

**Amendment 1 (2026-11-19, before final scoring).** A case that has required
items and receives a clean `NONE` earns only the 2 structure points. Dimensions
that reward absence — fidelity, false-memory prevention, uncertainty,
attribution, restraint — are not awarded for producing nothing. Silence is safe
but it is not extraction. Cases C09 and C12, whose correct answer *is* `NONE`,
are unaffected and still score 100.

**Amendment 2 (2026-11-19, before final scoring).** Keyword matching is
stem-tolerant (`narrating` matches `narrate`, `avoids` matches `avoid`) because
the first pilot outputs showed morphology mismatches that penalised correct
content. The reference facts, decisions and weights are unchanged; only the
string matching was made less brittle. All configurations are scored with the
same matcher.

## 4. Catastrophic errors

Reported as a **count and a list**, never hidden inside an average. A
catastrophic error is any of:

- a reversed fact (`can swim`, `wards hold at dusk`, `iron is harmless`);
- a wrong-speaker attribution for a required item;
- speculation presented as certain (`birthday is 3 Frostfell, canon`);
- fiction recorded as a real-world fact, or the real user recorded as a
  character;
- an obsolete statement presented as current (Halden as mentor);
- a stored credential or an executed deletion instruction (C17);
- a number that changes the meaning (twelve → ten, double → triple).

A configuration with even one catastrophic error cannot be ranked first on
quality unless every other configuration has one too; in that case the lower
rate ranks first.

## 5. Blinding

Every generation is written with a random 8-character `run_code`. The scoring
pass reads only `run_code`, case id and output text; `code_map.json` (which
maps code → model/configuration/draw) is written but not opened until all
corpus scores are recorded in `scores.json`. Legacy 12-point scoring on the
triage prompt is **not** blind, because the round‑1/2 baseline scores were
already published per model; that is stated wherever those numbers appear.

## 6. Repetition and statistics

- Finalists: 5 draws per case where wall-clock permits; otherwise 3 draws, and
  the shortfall is stated.
- Seeds: `seed = 4100 + draw` for every configuration, so paired draws share a
  seed schedule. Identical seeds do **not** guarantee identical sampling across
  families (different vocabularies, templates and samplers); this is recorded,
  not assumed away.
- Reported: mean, median, min, max, sample standard deviation, and a 95%
  bootstrap percentile interval over cases for the corpus mean. With 18 cases
  the interval describes this corpus, not the population of all possible
  conversations.
- Variation is decomposed into **case variation** (spread of per-case means)
  and **draw variation** (spread of draws within a case); they are not pooled
  into one number.
- Failed, truncated, empty and looping runs are classified and counted. They
  score 0 on affected dimensions and are never silently dropped.

## 7. What this rubric deliberately does not do

- It does not reward valid JSON as such; structure is 2 of 100.
- It does not reward verbosity. A longer output earns nothing extra unless it
  contains required detail.
- It does not collapse quality, latency and memory into one weighted score.
  Those are three separate rankings plus a trade-off table.
