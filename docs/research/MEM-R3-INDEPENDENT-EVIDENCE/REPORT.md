# Report — Round 3 independent model evidence

**Prepared:** 2026-11-19
**Scope:** the model evidence produced by the Round 3 memory-extraction harness
**Companion files:** `INDEX.md`, `REFERENCE-EXECUTION-METADATA.md`,
`REFERENCE-MODEL-IDENTITY.md`, `PACKET-00`–`PACKET-03M`
**Builder:** `build_evidence.py` (reads result files; writes this series; starts no
server and makes no network request)

This report describes what was collected and what could not be. It contains no
evaluation of model quality and no recommendation.

---

## 1. Source files that contain the required evidence

| Source | What it supplies | Coverage |
|---|---|---|
| `experiments/mem-extraction/round3/corpus.json` | the 18 case conversations, verbatim, and each case's `required`/`forbidden` grading key (key not reproduced here) | all 379 corpus-prompt runs |
| `experiments/mem-extraction/round3/harness.py` | `EXTRACTION_INSTRUCTION`; the composition expression at line 409; the `run_corpus` and `run_triage` record writers | all runs |
| `experiments/lfm25-cpu/run_probe.py` | `PROMPT`, the six-item prompt, imported by the harness | all 99 six-item runs |
| 19 generation result files under `results/` | stored `content`, `reasoning_content`, identifiers, source batch | 478 unique generations |
| `results/model_manifest.json` | byte size and SHA-256 for each model file | all 15 model files used |
| `raw/<tag>-<label>/server.log` | launch banner, including `n_ctx_slot` | per configuration |
| 15 `results/*_scores.json` files | `run_code` and derived scores only; **no response text** | not used by this series |

The 19 generation files are:

`corpus.json` (235), `all_triage.json` (51), `ablation_corpus_all.json` (108),
`ablation_triage.json` (15), `stageB.json` (18), `stageC.json` (6),
`tinyllama.json` (3), `budget.json` (2), `gpu-full.json` (2),
`gpu-full-qwen35-4b-nonthink-gpufull.json` (2), `control_corpus_2b.json` (36),
`ablation_corpus.json` (36), `ablation_corpus_abl-2b-nonthink.json` (36),
`ablation_corpus_abl-4b-nonthink.json` (36), `corpus_tinyllama.json` (36),
`stageA.json` (16), `stageA2a.json` (24), `stageA2b.json` (1),
`stageA2b2.json` (6).

Counts in that list are stored records, not unique runs; several files are
duplicates of runs also present in a merged file (see §4.1).

The 478 unique runs comprise 379 on the corpus prompt and 99 on the six-item
prompt, across 43 stored identifiers and 15 model files.

## 2. Which runs have complete prompt, model identity and response records

| Component | Runs with a complete record | Basis |
|---|---:|---|
| Stored response | 478 of 478 | `content` present; `reasoning_content` present separately (empty in 388, non-empty in 90) |
| Model identity | 478 of 478 | stored `label` and `file`; byte size and SHA-256 available for all 15 files |
| Recorded prompt | **0 of 478** | no request body, prompt string or prompt hash is stored |
| Reconstructable prompt | 478 of 478 | 379 from instruction + case text; 99 from one constant |

No run is missing a response or a model identifier. Every run is missing a
recorded prompt.

What "complete model identity" means here: the stored identifier and file name,
plus the manifest's size and hash. Quantisation is not a separate field and is
visible only in the file name (`Q4_K_M`, `QAD-Q4_0`). A repository identifier was
recorded only for the five ablation conditions, in `notes`; for baseline
conditions it is absent from the run record.

What "reconstructable prompt" means, and does not mean: the instruction text and
the composition expression are read from `harness.py`, and the case text from
`corpus.json`. Applying them produces a candidate prompt. It is labelled
`SOURCE-DERIVED` in every packet. It is not a recorded request, and its equality
with the bytes the server received is `UNESTABLISHED` (Packet 01 §5).

## 3. Runs with missing or incomplete components

| Gap | Runs affected | Detail |
|---|---:|---|
| No stored prompt / request body | 478 | `prompt_tokens` is the only prompt-side record |
| No `ctx` field | 379 | all corpus-prompt runs; context appears only in `raw/corpus-*/server.log` |
| No `command` field | 379 | all corpus-prompt runs |
| No `tag` field | 235 | the runs held in `results/corpus.json` |
| No stored `system` field | 379 | corpus path sets no system message; this is read from code, not recorded |
| Empty stored `content` | 2 | `yjyha4dh` (`qwen35-2b-think`, six-item, `length`) and `t1hggto0` (`corpus-lfm3-nonthink`, C01, `length`) |
| Truncated completions (`finish_reason: length`) | 27 | `corpus-tinyllama` 17, `corpus-lfm3-think-512` 6, `corpus-lfm3-nonthink` 1, `abl-4b-nonthink` 1, `qwen35-2b-think` 1, `qwen35-9b-think-rb1024` 1 |
| Planned but absent | 53 | `plan.json` schedules 3 draws × 18 cases for `corpus-lfm3-nonthink`; 1 run is stored |
| One stored draw only | 3 identifiers | `qwen35-2b-think` (1 run), `qwen35-9b-think-rb1024` (1 run), `corpus-4b-think-512` (1 draw × 18 cases = 18 runs) |
| Prompt-token verification | 478 | `prompt_tokens` is recorded for all runs but was not cross-checked against any tokenizer |

`corpus-lfm3-nonthink` is the clearest gap: 54 runs were planned, one was stored,
and that one returned an empty `content` with `finish_reason: length` at a
512-token cap. The condition is therefore present in the evidence as a single
failed run, not as a completed condition.

The corpus prompt was run with eleven identifiers (seven stock, three ablation,
one control) and the six-item prompt with 36; four identifiers are used for both
prompt types. Model files that appear only in the six-item set (the 0.8B, 350M,
2B, Qwen3-1.7B, SmolLM3, 9B, LFM2.5-1.2B-Thinking and abl-LFM3 files) have no
corpus-prompt generation at all. This is observable coverage, not an omission
from this series.

## 4. Source conflicts and discrepancies discovered

### 4.1 Duplicate copies of 191 runs — no disagreement

191 of the 478 `run_code`s appear in two generation files. The pairs are:
108 ablation corpus runs (each per-condition file plus
`ablation_corpus_all.json`), 36 TinyLlama corpus runs (`corpus.json` plus
`corpus_tinyllama.json`), and 47 six-item runs (`stageA`, `stageA2a`, `stageA2b`,
`stageA2b2` plus `all_triage.json`). Comparing every copy pair on `content` and
`reasoning_content` produced **0 disagreements**. The index names the copy used
and the other file holding the same record.

The 15 score files also repeat `run_code`s, but carry no `content`; they are
derived records and were not treated as evidence.

### 4.2 Identifiers reused across prompt types

Four identifiers name both corpus-prompt and six-item runs:
`abl-2b-nonthink`, `abl-4b-nonthink`, `abl-lfm120`, `qwen35-2b-nonthink`. The
string alone is therefore not unique; `kind` and the source file disambiguate.
This is a naming collision, not a conflict in the data.

### 4.3 Recorded sampler values versus one derived reading utility

`qual_extract.py` (a reading generator, not the harness) describes four corpus
conditions with values that differ from both `configs.py`/`plan.json` and the
`sampling` object stored in every run:

| Identifier | `qual_extract.py` prose | Stored `sampling` |
|---|---|---|
| `corpus-4b-think-512` | T0.6 / top_p0.95 | temperature 1.0 / top_p 0.95 |
| `corpus-lfm120-instruct` | T0.1 / top_p0.1 | temperature 0.1 / top_p 1.0 |
| `corpus-lfm120-t08` | T0.8 / top_p0.95 | temperature 0.8 / top_p 1.0 |
| `corpus-lfm3-think-512` | T0.1 / top_p0.1 | temperature 0.1 / top_p 1.0 |

The stored values agree with `configs.py` and `plan.json`. The conflict is
confined to the prose in `qual_extract.py`; it is recorded here rather than
corrected, and none of its text is used in this series.

A related case that is **not** a conflict: `configs.lfm3(False)` declares
`max_tokens=3072`, while the stored `corpus-lfm3-nonthink` run has 512.
`plan.json` overrides the cap to 512 for that label, and `run_corpus_plan.py`
applies the override. The record matches the plan, not the base spec.

### 4.4 Plan versus stored runs

`plan.json` schedules 252 corpus runs (54 + 54 + 36 + 18 + 36 + 54). 199 of them
are stored: all of the first five conditions in full, and 1 of 54 for
`corpus-lfm3-nonthink`. The remaining 36 stored corpus runs
(`corpus-tinyllama`) are not in `plan.json`; they were produced by a separate
script with a 2,048-token context.

### 4.5 Context size is absent from corpus records

No corpus record carries `ctx`. The server logs are the only place the value
appears: `n_ctx_slot = 2048` for `corpus-tinyllama` and `4096` for the other
corpus configurations. A reviewer comparing corpus runs without consulting the
logs would not see that one condition had half the context.

### 4.6 Manifest written after the runs

`results/model_manifest.json` has modification time 2026-11-19 02:12; the latest
generation file is 02:11. The hashes identify the files present when the manifest
was written. They do not prove the model files were unchanged while the runs were
produced.

### 4.7 The corpus's own provenance claim cannot be verified

`corpus.json` states that its reference answers were "written BEFORE any round-3
model output was generated". There is no pre-run hash or commit, so the ordering
cannot be proved from the workspace. This affects the grading key, not the model
responses, and the key is not part of this series.

### 4.8 No conflicts found in response text

No two stored copies of a run disagree, and no run appears under two different
model files.

## 5. Packet structure and run coverage

| Series | Packets | Runs | Structure |
|---|---|---:|---|
| Front matter | 00, 01 | 0 | scope and prompt record |
| Corpus by case | 02A–02R (36 files) | 379 | one case per packet, split at a run boundary when over ~3,400 words |
| Six-item by configuration | 03A–03M (13 files) | 99 | grouped by stored identifier, split at run boundaries |
| **Total** | **51 files** | **478** | 118,904 words |

Coverage: C01 has 22 stored runs; each of C02–C18 has 21. Every stored run for
every case appears exactly once in the series, and the six-item series covers all
99 runs.

Ordering is alphabetical by stored identifier and then by draw. No ordering by
any measured quantity was used, deliberately, so that the sequence of packets
does not imply a ranking.

Each corpus packet prints the reconstructed prompt for its case once, then the
runs. Each six-item packet points to the single prompt in Packet 01 and prints
only responses and identifiers. The `system` value is printed where a run
recorded one.

`INDEX.md` lists every run with its `run_code`, kind, case, identifier, draw,
packet, primary source and array index, and any other file holding the same
generation.

Size: the largest packet is 3,561 words; the smallest evidence packet is 510.
No packet exceeds 4,200 words.

## 6. Preservation and integrity checks actually performed

| Check | Method | Result |
|---|---|---|
| Result files parse | JSON parse of all 34 result files that hold run records (19 generation, 15 score) | 0 failures |
| Unique runs | set of `run_code` across generation files | 478; every one has a record containing `content` |
| Duplicate agreement | compare `content` and `reasoning_content` across every copy | 191 duplicated codes; 0 disagreements |
| Response preservation | write the decoded string as UTF-8 bytes; re-read the packet and test the string as an exact byte substring | 478/478 `content`, 90/90 non-empty `reasoning_content` |
| Packet attribution | every `run_code` a packet claims must appear in that file | 0 failures |
| Size bound | word count per packet | max 3,561; 0 packets over 4,200 |
| Excluded material | strip fenced blocks, scan remaining text for scorer artefacts (`v1 score`, `v2 \``, `MALFORMED`, `VALID_*`, `catastrophic`, event codes) | 0 hits in the 49 evidence packets |
| Score files unused | builder ingests only records containing `content`; score files contain none | asserted at load |
| Source files untouched | modification times of the result files before and after the build; builder writes only inside this directory | `results/corpus.json` 2026-11-19 01:08, `results/all_triage.json` 2026-11-18 23:52, `results/ablation_corpus_all.json` 02:11 — unchanged. `experiments/` is untracked in this workspace (`git status` shows one `?? experiments/` entry), so git cannot corroborate preservation there |
| Prompt-source identity | SHA-256 of `harness.py`, `corpus.json`, `run_probe.py` recorded in Packet 01 | recorded |

Not performed, and therefore not claimed:

1. Byte-for-byte equality between a packet and a JSON file's raw bytes. The
   claim is narrower: the decoded string is written unchanged and verified as an
   exact substring.
2. Recovery or verification of the prompt bytes the server received.
3. Per-tokenizer re-tokenisation of the reconstructed prompts against the
   recorded `prompt_tokens`.
4. Verification that any recorded sampler value was applied by the server.
5. Verification of the corpus's pre-registration ordering.

The series is generated by `build_evidence.py` from the stored files; rerunning
it reproduces the packets and repeats every check above.

## 7. Standing limitation

The response side of this evidence set is complete and traceable. The prompt side
is not recorded and cannot be verified from the workspace; it is reconstructed
from the code that built the requests and is labelled as such throughout. Any
review that needs the exact prompt bytes must treat that as an open gap rather
than a recovered fact.
