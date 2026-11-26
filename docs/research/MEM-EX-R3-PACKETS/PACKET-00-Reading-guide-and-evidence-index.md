# Packet 00 — Reading guide and evidence index

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Packets:** 45 files · **Total:** 128,206 words
**Master table:** the same table is repeated in `INDEX.md`.

## What this series is

Round 3 compared local language models on one narrow task: given a roleplay
conversation, produce memory candidates as `ID | DECISION | MEMORY | REASON`.
The series contains the raw generations, the references they were scored
against, the scoring code, the hardware measurements, and two independent
readings of the result. It contains no production decision.

## Source hierarchy

Read in this order of authority; later layers never override earlier ones.

1. **Raw generations** — `results/corpus.json` (235 runs), `results/all_triage.json`
   (51 runs), `results/ablation_corpus_all.json` (108 runs). These are the
   evidence. Nothing in the series rewrites them.
2. **Execution records** — the `command`, `sampling`, `seed`, `finish_reason`,
   token counts and timings stored beside each generation; `raw/**/server.log`.
3. **Frozen references** — `corpus.json` (18 cases with expected decisions and
   forbidden strings), `RUBRIC.md`, `SOURCES.md`.
4. **Derived scores** — `results/*_scores.json` (historical, v1) and
   `instrument_v2/out/corpus_v2.json` (separate corrected instrument, v2).
5. **Interpretation** — `docs/research/MEM-EXTRACTION-ROUND3-2026-11-19.md`,
   `qualitative/EVALUATION-BUFFY.md`, Muse's audit, Packet 09.

A derived score is not evidence that the model did something. It is evidence
that a scorer computed something. Where the two disagree, the raw text wins.

## Identifiers

- **`run_code`** — random 8-character id per generation, used for blind scoring.
  It is the join key between every packet and every JSON file.
- **`label`** — the configuration name (`corpus-4b-nonthink`, …). Written to the
  run record; *not* read by the v1 scorer, which is what made scoring blind.
- **`tag`** — the run batch that produced the file.
- **Case id** — `C01`–`C18` in `corpus.json`; `M1`–`M6` in the six-item prompt.
- **Item id** — `C1`, `C2`, … inside a response, one per candidate memory.
- **`draw`** — repetition index within the same case and configuration; seeds
  are `4100 + draw`.
- **`finish_reason`** — `stop` (model ended) or `length` (cap reached).
- **Status (v2)** — `VALID_KEEP`, `VALID_HOLD`, `VALID_REJECT`, `VALID_NONE`,
  `VALID_EMPTY`, `MALFORMED`, `UNPARSEABLE`, `MISSING`, `INCOMPLETE`,
  `EXECUTION_FAILURE`.
- **Event codes (v2)** — `FORBIDDEN_CLAIM_STORED`, `ATTRIBUTION_ERROR`,
  `PROMOTION`, `INHERITANCE`, `CONTRADICTION_COLLAPSE`, `SCENARIO_LEAK`,
  `ATOMICITY_MERGED_DISTINCT`, `EXECUTION_FAILURE`.

## Labels used in the packets

`OBSERVED` directly present in a stored response · `SOURCE-DERIVED` stated in a
source or config · `CALCULATED` explicit arithmetic · `INFERRED` a reading ·
`UNESTABLISHED` not established · `SUPERSEDED` replaced by later evidence.

**Whose interpretation:** fenced text is the model's. Packet 08 is Muse's.
Packet 09 and any unfenced commentary are Buffy's. The project owner has not
written a reading in this series.

## Raw, derived, interpreted — how to tell

| You are looking at | Layer | Trust it as |
|---|---|---|
| A fenced block after "verbatim `content`" | raw | what the model emitted |
| A `v1 score` in a run header | derived | the historical 100-point rubric |
| A `v2` line in a run header | derived | a different, corrected instrument |
| Packet 02B code excerpt | source | the actual scoring implementation |
| Packet 08 | interpretation | one adversarial reviewer's judgment |
| Packet 09 | interpretation | questions, not findings |

## Known missing evidence

1. **The Muse audit original.** No repository copy exists; Packet 08 is a
   transcription of the delivered message.
2. **Full `run_code` uniqueness proof** across the four ablation JSON files
   (audit F-13, open).
3. **Primary-source card re-fetch** at pinned revisions — sampler claims rest on
   `SOURCES.md` URLs and retrieval dates.
4. **Per-model template rendering dumps** for the Qwen default-mode conflict and
   the LFM3 off-switch.
5. **A truncation-free LFM3 corpus run** to separate quality from the 512-token
   cap.
6. **Inter-rater reliability** — one scorer, no second human.
7. **Concurrency / cross-seed isolation matrix** — never run.
8. **The DeepSeek dossier and the owner's adversarial review** — separate tracks,
   not produced.
9. **9B corpus run** — never performed.
10. **Prompt-token accounting per tokenizer** — recorded `prompt_tokens` not
    cross-checked against each tokenizer.

## Caveats before you read any score

- Draw counts are 1–3. A packet header saying `n=1` means one observation, not a
  ranking.
- The v1 and v2 totals are different instruments. Do not compare them directly
  or substitute one for the other.
- v2's tests prove the evaluator follows its own rules. They do not prove the
  rules track human judgment.
- The corpus has one author and 18 cases. Confidence intervals, where present,
  describe this corpus only.
- Nothing here is MEMORY-09 evidence or a production recommendation.

## Master index

| Packet | Title | Sources | Cases / runs | Words | Raw output | Omitted |
|---|---|---|---|---:|---|---|
| 01 | Experimental design and protocol | `harness.py`, `configs.py`, `corpus.json`, `RUBRIC.md`, `SOURCES.md` | all 18 cases and 51 triage runs | 1,481 | n/a — design only | nothing omitted |
| 02A | The rubric as written | `experiments/mem-extraction/round3/RUBRIC.md` (143 lines, verbatim) | not case-specific; applies to all 18 corpus cases and the six-item prompt | 1,483 | none — the rubric contains no model output | nothing omitted; the file is reproduced in full |
| 02B | How the scorer implements the rubric | `score_corpus.py` (295 lines) and `score_triage.py` (106 lines), excerpts | scoring functions, not case-specific | 1,047 | none — no model output | The excerpts below are the parser, the fidelity/false-memory block and the |
| 03A | Raw model outputs — six-item triage: lfm120-instruct-official | `results/all_triage.json`, `results/*_scores.json` | lfm120-instruct-official | 807 | verbatim | none — every stored triage run appears in this series |
| 03B-1 | Raw model outputs — six-item triage: lfm120-think-official (part 1) | `results/all_triage.json`, `results/*_scores.json` | lfm120-think-official | 1,961 | verbatim | none — every stored triage run appears in this series |
| 03B-2 | Raw model outputs — six-item triage: lfm120-think-official (part 2) | `results/all_triage.json`, `results/*_scores.json` | lfm120-think-official | 2,236 | verbatim | none — every stored triage run appears in this series |
| 03B-3 | Raw model outputs — six-item triage: lfm120-think-official (part 3) | `results/all_triage.json`, `results/*_scores.json` | lfm120-think-official | 1,478 | verbatim | none — every stored triage run appears in this series |
| 03C-1 | Raw model outputs — six-item triage: lfm120-think-official-rb1024 (part 1) | `results/all_triage.json`, `results/*_scores.json` | lfm120-think-official-rb1024 | 2,937 | verbatim | none — every stored triage run appears in this series |
| 03C-2 | Raw model outputs — six-item triage: lfm120-think-official-rb1024 (part 2) | `results/all_triage.json`, `results/*_scores.json` | lfm120-think-official-rb1024 | 1,861 | verbatim | none — every stored triage run appears in this series |
| 03D | Raw model outputs — six-item triage: lfm3-17b-nonthink–lfm350-official | `results/all_triage.json`, `results/*_scores.json` | lfm3-17b-nonthink–lfm350-official | 3,256 | verbatim | none — every stored triage run appears in this series |
| 03E | Raw model outputs — six-item triage: qwen3-17b-nonthink–qwen35-2b-nonthink | `results/all_triage.json`, `results/*_scores.json` | qwen3-17b-nonthink–qwen35-2b-nonthink | 2,888 | verbatim | none — every stored triage run appears in this series |
| 03F | Raw model outputs — six-item triage: qwen35-2b-think-rb1024–qwen35-4b-nonthink | `results/all_triage.json`, `results/*_scores.json` | qwen35-2b-think-rb1024–qwen35-4b-nonthink | 3,059 | verbatim | none — every stored triage run appears in this series |
| 03G | Raw model outputs — six-item triage: qwen35-4b-think-rb1024–smollm3-3b-nonthink | `results/all_triage.json`, `results/*_scores.json` | qwen35-4b-think-rb1024–smollm3-3b-nonthink | 2,884 | verbatim | none — every stored triage run appears in this series |
| 03H | Raw model outputs — six-item triage: smollm3-3b-think-rb1024–qwen35-4b-think-rb512 | `results/all_triage.json`, `results/*_scores.json` | smollm3-3b-think-rb1024–qwen35-4b-think-rb512 | 3,058 | verbatim | none — every stored triage run appears in this series |
| 03I | Raw model outputs — six-item triage: qwen35-2b-think | `results/all_triage.json`, `results/*_scores.json` | qwen35-2b-think | 2,394 | verbatim | none — every stored triage run appears in this series |
| 03J | Raw model outputs — six-item triage: qwen35-9b-think-rb1024 | `results/all_triage.json`, `results/*_scores.json` | qwen35-9b-think-rb1024 | 2,335 | verbatim | none — every stored triage run appears in this series |
| 04A | Raw model outputs — corpus by case: C01 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C01 | 3,959 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04B | Raw model outputs — corpus by case: C02 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C02 | 3,798 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04C | Raw model outputs — corpus by case: C03 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C03 | 3,269 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04D | Raw model outputs — corpus by case: C04 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C04 | 3,537 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04E | Raw model outputs — corpus by case: C05 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C05 | 3,536 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04F | Raw model outputs — corpus by case: C06 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C06 | 3,673 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04G | Raw model outputs — corpus by case: C07 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C07 | 3,580 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04H | Raw model outputs — corpus by case: C08 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C08 | 3,493 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04I | Raw model outputs — corpus by case: C09 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C09 | 3,412 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04J-1 | Raw model outputs — corpus by case: C10 (part 1) | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C10 | 2,460 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04J-2 | Raw model outputs — corpus by case: C10 (part 2) | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C10 | 2,770 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04K | Raw model outputs — corpus by case: C11 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C11 | 3,349 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04L | Raw model outputs — corpus by case: C12 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C12 | 3,981 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04M | Raw model outputs — corpus by case: C13 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C13 | 3,630 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04N | Raw model outputs — corpus by case: C14 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C14 | 3,660 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04O | Raw model outputs — corpus by case: C15 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C15 | 3,387 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04P | Raw model outputs — corpus by case: C16 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C16 | 3,575 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04Q | Raw model outputs — corpus by case: C17 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C17 | 3,327 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 04R | Raw model outputs — corpus by case: C18 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C18 | 3,306 | verbatim | none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries. |
| 05A | Adversarial and failure cases: C01 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C01 | 4,019 | verbatim | none for these five cases; every draw is present. Other cases are in Packet 04. |
| 05B | Adversarial and failure cases: C12 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C12 | 4,039 | verbatim | none for these five cases; every draw is present. Other cases are in Packet 04. |
| 05C | Adversarial and failure cases: C13 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C13 | 3,686 | verbatim | none for these five cases; every draw is present. Other cases are in Packet 04. |
| 05D | Adversarial and failure cases: C17 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C17 | 3,384 | verbatim | none for these five cases; every draw is present. Other cases are in Packet 04. |
| 05E | Adversarial and failure cases: C18 | `corpus.json`, `results/corpus.json`, `results/corpus_scores.json` | C18 | 3,365 | verbatim | none for these five cases; every draw is present. Other cases are in Packet 04. |
| 06 | Context isolation and reset | `run_reset.py`, `run_isolation.py`, `results/reset.json`, `results/isolation.json`, `raw/reset/server.log` | C03 (benign), C17 (adversarial), sentinel probes | 1,126 | n/a — probe outputs, reproduced from JSON | Server logs are excerpted to the request lines and cache metrics; the full logs remain at the paths above. |
| 07 | Ablation, performance and hardware | `results/gpu-full.json`, `results/coexist.json`, `results/ablation_corpus_all.json`, `results/ablation_*_scores.json`, `configs_ablation.py` | 4 ablation conditions; GPU 4B and 1.2B; one coexistence pair | 2,759 | raw ablation generations are cited by run_code, not all reproduced | 108 ablation runs are summarised rather than reproduced in full: this packet carries the configuration, the scores and the aggregate, and `results/ablation_corpus_all.json` holds every generation. Reproducing all 108 here would exceed the packet size and duplicate Packet 04's evidence for the stock baselines. |
| 08A | Muse audit: verdict, findings register, configuration matrix | the Muse adversarial audit as delivered in session on | all | 2,699 | none — this is a review document, not model | audit sections D–H are in Packet 08B; nothing else is omitted |
| 08B | Muse audit: isolation, validity matrices, disposition | the Muse adversarial audit as delivered in session on | all | 1,235 | none — this is a review document | audit sections A–C are in Packet 08A; nothing else is omitted |
| 09 | Outstanding questions and review checklist | synthesis of all packets | all | 1,027 | n/a — no model output | nothing omitted; this packet is argument, not evidence |

## Recommended order for a first read

1. **00** (this packet) and **01** — what was asked and what was measured.
2. **02A** and **02B** — the rules and their known defects.
3. **05A** — C01, where a 100-point record is four rules in one line.
4. **08A** and **08B** — Muse's audit, including the missing-evidence list.
5. **04A**–**04R** as needed when a specific case is disputed.

Send **00, 01, 02A, 02B, 05A, 08A, 08B** first. That is about
11,964
words and is enough for an informed first verdict; the remaining packets are the
evidence behind individual claims.

## Size note

Target was 2,500–4,000 words per packet. 0 packet(s) exceed 4,200
words after splitting at run boundaries, because a single generation with several
thousand tokens of reasoning cannot be cut without breaking the verbatim
guarantee. They are listed in `INDEX.md` with their sizes.
