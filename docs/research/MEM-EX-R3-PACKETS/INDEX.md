# Master index — Round 3 source-preserving reading packets

Generated from each packet's own header. Word counts are approximate.

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

**Total:** 128,206 words across 45 packets.

**First send:** 00, 01, 02A, 02B, 05A, 08A, 08B.

**Verbatim guarantee:** every model response in Packets 03–05 was inserted
directly from `results/*.json` by `make_packets.py`, which asserts that each
stored `content` string appears in the series. Zero responses were missing at
generation time. Nothing was normalised, shortened or reworded.

**Oversize packets:** none.

## Evidence that could not be included or recovered

1. **The Muse audit original file.** Produced in a chat session; no repository
   copy existed. Packet 08 is a transcription of the delivered message, and is
   the only copy in the workspace.
2. **Full `run_code` uniqueness proof across the ablation files** — never
   performed (audit F-13). The records themselves are present; the proof is not.
3. **Primary-source card text at pinned revisions** — citations and retrieval
   dates are in `SOURCES.md`; the pages were not re-fetched for this series.
4. **Chat-template renderings per model** — not dumped; the default-mode
   conflict and the LFM3 off-switch are inferred from observed reasoning fields.
5. **A truncation-free LFM3 corpus** — never run.
6. **Human scores and inter-rater reliability** — never collected; one scorer.
7. **Concurrency isolation matrix** — never run.
8. **The DeepSeek dossier and the owner's adversarial review** — separate tracks,
   not produced, so they are not in this series.
9. **The 9B corpus run** — never performed; only triage and a ceiling probe exist.
10. **Per-tokenizer prompt accounting** — recorded `prompt_tokens` were not
    cross-checked against each model's tokenizer.
11. **Full server logs** — only reset/coexistence logs are quoted; the rest are
    cited by path.
12. **108 ablation generations** — summarised in Packet 07 rather than reproduced;
    the file is cited and the omission is stated in that packet's header.
