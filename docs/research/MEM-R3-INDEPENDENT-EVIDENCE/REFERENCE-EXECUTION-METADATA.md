# Reference index — execution metadata

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19

This file is the separate reference index required by the evidence-preparation
brief. It exists so that sampler parameters, context settings, token counts,
timings and hardware measurements stay out of the reading packets while
remaining independently inspectable.

Nothing here is evidence that a configuration was applied. A recorded value is a
statement about what the harness intended to send and what the server reported;
it is not a measurement of the server's effective sampler chain.

## 1. Where each value is stored

| Value | JSON key | Present in |
|---|---|---|
| Sampler chain | `sampling` (object) | all generation records |
| Output cap | `max_tokens` | all generation records |
| Reasoning budget | `reasoning_budget` | all generation records (may be null) |
| Reasoning effort / flag | `reasoning_effort`, `reasoning_flag` | six-item records |
| Context size | `ctx` | six-item records only |
| GPU layers | `ngl` | all generation records |
| Seed | `seed` | all generation records |
| Server launch command | `command` | six-item records only |
| Prompt / completion tokens | `prompt_tokens`, `completion_tokens` | all generation records |
| Timings | `ttft_s`, `thinking_s`, `answer_s`, `wall_s`, `load_s` | all generation records |
| Throughput | `prefill_ms`, `decode_ms`, `decode_tok_s`, `prefill_tok_s` | six-item records |
| Peak memory | `peak_rss_mb`, `peak_vram_mb`, `peak_gpu_util` | all generation records |
| Server log | `raw/<tag>-<label>/server.log` | files on disk |

`OBSERVED` gaps in the records themselves:

- Corpus records store no `ctx` and no `command`. Context size is therefore not
  recorded per corpus run; the corpus launch logs under `raw/corpus-<label>/`
  contain the server banners and are the place to check it.
- Corpus records store no `tag`; the batch is recoverable only from which result
  file holds the record (see `INDEX.md`).
- Timings are wall-clock around the HTTP request and exclude model load, which is
  stored separately in `load_s` where recorded.

## 2. Per-configuration recorded values

One row per stored model identifier. Where a value varies within a label, every
distinct value is listed with a count.

### `abl-2b-nonthink`

* Runs: 39 (corpus 36, triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 2.0, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 1.0}`×39
* Max tokens: `768`×39
* Reasoning budget: `None`×39
* `ngl`: `0`×39
* `ctx`: `4096`×3, `None`×36
* System message: `None`×39
* `command` recorded: 3 of 39 runs
* Server log(s): `raw/ablation_corpus_abl-2b-nonthink-abl-2b-nonthink/server.log`, `raw/ablation_triage-abl-2b-nonthink/server.log`

### `abl-4b-nonthink`

* Runs: 39 (corpus 36, triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 0.7, "top_k": 20, "top_p": 0.8}`×39
* Max tokens: `768`×39
* Reasoning budget: `None`×39
* `ngl`: `0`×39
* `ctx`: `4096`×3, `None`×36
* System message: `None`×39
* `command` recorded: 3 of 39 runs
* Server log(s): `raw/ablation_corpus_abl-4b-nonthink-abl-4b-nonthink/server.log`, `raw/ablation_triage-abl-4b-nonthink/server.log`

### `abl-4b-think-512`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 0.95}`×3
* Max tokens: `1024`×3
* Reasoning budget: `512`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/ablation_triage-abl-4b-think-512/server.log`

### `abl-lfm120`

* Runs: 39 (corpus 36, triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×39
* Max tokens: `512`×39
* Reasoning budget: `None`×39
* `ngl`: `0`×39
* `ctx`: `4096`×3, `None`×36
* System message: `None`×39
* `command` recorded: 3 of 39 runs
* Server log(s): `raw/ablation_corpus-abl-lfm120/server.log`, `raw/ablation_triage-abl-lfm120/server.log`

### `abl-lfm3-think-512`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.1, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `1024`×3
* Reasoning budget: `512`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/ablation_triage-abl-lfm3-think-512/server.log`

### `corpus-4b-nonthink`

* Runs: 54 (corpus 54)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 0.7, "top_k": 20, "top_p": 0.8}`×54
* Max tokens: `768`×54
* Reasoning budget: `None`×54
* `ngl`: `0`×54
* `ctx`: `None`×54
* System message: `None`×54
* `command` recorded: 0 of 54 runs
* Server log(s): `raw/corpus-corpus-4b-nonthink/server.log`

### `corpus-4b-think-512`

* Runs: 18 (corpus 18)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 0.95}`×18
* Max tokens: `1024`×18
* Reasoning budget: `512`×18
* `ngl`: `0`×18
* `ctx`: `None`×18
* System message: `None`×18
* `command` recorded: 0 of 18 runs
* Server log(s): `raw/corpus-corpus-4b-think-512/server.log`

### `corpus-lfm120-instruct`

* Runs: 54 (corpus 54)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×54
* Max tokens: `512`×54
* Reasoning budget: `None`×54
* `ngl`: `0`×54
* `ctx`: `None`×54
* System message: `None`×54
* `command` recorded: 0 of 54 runs
* Server log(s): `raw/corpus-corpus-lfm120-instruct/server.log`

### `corpus-lfm120-t08`

* Runs: 36 (corpus 36)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.8, "top_k": 50, "top_p": 1.0}`×36
* Max tokens: `512`×36
* Reasoning budget: `None`×36
* `ngl`: `0`×36
* `ctx`: `None`×36
* System message: `None`×36
* `command` recorded: 0 of 36 runs
* Server log(s): `raw/corpus-corpus-lfm120-t08/server.log`

### `corpus-lfm3-nonthink`

* Runs: 1 (corpus 1)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.1, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×1
* Max tokens: `512`×1
* Reasoning budget: `None`×1
* `ngl`: `0`×1
* `ctx`: `None`×1
* System message: `None`×1
* `command` recorded: 0 of 1 runs
* Server log(s): `raw/corpus-corpus-lfm3-nonthink/server.log`

### `corpus-lfm3-think-512`

* Runs: 36 (corpus 36)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.1, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×36
* Max tokens: `1024`×36
* Reasoning budget: `512`×36
* `ngl`: `0`×36
* `ctx`: `None`×36
* System message: `None`×36
* `command` recorded: 0 of 36 runs
* Server log(s): `raw/corpus-corpus-lfm3-think-512/server.log`

### `corpus-tinyllama`

* Runs: 36 (corpus 36)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 0.7, "top_k": 50, "top_p": 0.95}`×36
* Max tokens: `768`×36
* Reasoning budget: `None`×36
* `ngl`: `0`×36
* `ctx`: `None`×36
* System message: `None`×36
* `command` recorded: 0 of 36 runs
* Server log(s): `raw/corpus-corpus-tinyllama/server.log`

### `lfm120-ins-r2cap`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.05, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 0.8, "top_k": 40, "top_p": 0.95}`×3
* Max tokens: `400`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `2048`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageC-lfm120-ins-r2cap/server.log`

### `lfm120-ins-t03`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.3, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageB-lfm120-ins-t03/server.log`

### `lfm120-ins-t05`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.5, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageB-lfm120-ins-t05/server.log`

### `lfm120-ins-t08`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.8, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageB-lfm120-ins-t08/server.log`

### `lfm120-ins-t08-rep10`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 0.8, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageB-lfm120-ins-t08-rep10/server.log`

### `lfm120-instruct-official`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA-lfm120-instruct-official/server.log`

### `lfm120-instruct-official-gpufull`

* Runs: 2 (triage 2)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×2
* Max tokens: `768`×2
* Reasoning budget: `None`×2
* `ngl`: `99`×2
* `ctx`: `4096`×2
* System message: `None`×2
* `command` recorded: 2 of 2 runs
* Server log(s): `raw/gpu-full-lfm120-instruct-official-gpufull/server.log`

### `lfm120-think-official`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.05, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `3072`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA-lfm120-think-official/server.log`

### `lfm120-think-official-rb1024`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.05, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `3072`×3
* Reasoning budget: `1024`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA2a-lfm120-think-official-rb1024/server.log`

### `lfm3-17b-nonthink`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.1, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `3072`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA2a-lfm3-17b-nonthink/server.log`

### `lfm3-17b-think-rb1024`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.1, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `3072`×3
* Reasoning budget: `1024`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA2a-lfm3-17b-think-rb1024/server.log`

### `lfm3-17b-think-rb512`

* Runs: 2 (triage 2)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.1, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×2
* Max tokens: `3072`×2
* Reasoning budget: `512`×2
* `ngl`: `0`×2
* `ctx`: `4096`×2
* System message: `None`×2
* `command` recorded: 2 of 2 runs
* Server log(s): `raw/budget-lfm3-17b-think-rb512/server.log`

### `lfm350-official`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.05, "temperature": 0.1, "top_k": 50, "top_p": 1.0}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA-lfm350-official/server.log`

### `qwen3-17b-nonthink`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 0.7, "top_k": 20, "top_p": 0.8}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA2a-qwen3-17b-nonthink/server.log`

### `qwen3-17b-think-rb1024`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 0.6, "top_k": 20, "top_p": 0.95}`×3
* Max tokens: `3072`×3
* Reasoning budget: `1024`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA2a-qwen3-17b-think-rb1024/server.log`

### `qwen35-08b-nonthink`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 2.0, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 1.0}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA-qwen35-08b-nonthink/server.log`

### `qwen35-2b-nonthink`

* Runs: 39 (triage 3, corpus 36)
* Sampling: `{"min_p": 0.0, "presence_penalty": 2.0, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 1.0}`×39
* Max tokens: `2048`×39
* Reasoning budget: `None`×39
* `ngl`: `0`×39
* `ctx`: `4096`×3, `None`×36
* System message: `None`×39
* `command` recorded: 3 of 39 runs
* Server log(s): `raw/control_corpus_2b-qwen35-2b-nonthink/server.log`, `raw/stageA-qwen35-2b-nonthink/server.log`

### `qwen35-2b-nonthink-pres0`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 1.0}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageB-qwen35-2b-nonthink-pres0/server.log`

### `qwen35-2b-nonthink-r2cap`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.05, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 0.8, "top_k": 40, "top_p": 0.95}`×3
* Max tokens: `400`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `2048`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageC-qwen35-2b-nonthink-r2cap/server.log`

### `qwen35-2b-nonthink-round2`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.05, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 0.8, "top_k": 40, "top_p": 0.95}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageB-qwen35-2b-nonthink-round2/server.log`

### `qwen35-2b-think`

* Runs: 1 (triage 1)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 0.95}`×1
* Max tokens: `3072`×1
* Reasoning budget: `None`×1
* `ngl`: `0`×1
* `ctx`: `4096`×1
* System message: `None`×1
* `command` recorded: 1 of 1 runs
* Server log(s): `raw/stageA-qwen35-2b-think/server.log`

### `qwen35-2b-think-rb1024`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 0.95}`×3
* Max tokens: `3072`×3
* Reasoning budget: `1024`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA2a-qwen35-2b-think-rb1024/server.log`

### `qwen35-4b-nonthink`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 0.7, "top_k": 20, "top_p": 0.8}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA2a-qwen35-4b-nonthink/server.log`

### `qwen35-4b-nonthink-gpufull`

* Runs: 2 (triage 2)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 0.7, "top_k": 20, "top_p": 0.8}`×2
* Max tokens: `768`×2
* Reasoning budget: `None`×2
* `ngl`: `99`×2
* `ctx`: `4096`×2
* System message: `None`×2
* `command` recorded: 2 of 2 runs
* Server log(s): `raw/gpu-full-qwen35-4b-nonthink-gpufull-qwen35-4b-nonthink-gpufull/server.log`

### `qwen35-4b-think-rb1024`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 0.95}`×3
* Max tokens: `3072`×3
* Reasoning budget: `1024`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA2a-qwen35-4b-think-rb1024/server.log`

### `qwen35-4b-think-rb256`

* Runs: 2 (triage 2)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 0.95}`×2
* Max tokens: `3072`×2
* Reasoning budget: `256`×2
* `ngl`: `0`×2
* `ctx`: `4096`×2
* System message: `None`×2
* `command` recorded: 2 of 2 runs
* Server log(s): `raw/budget-qwen35-4b-think-rb256/server.log`

### `qwen35-4b-think-rb512`

* Runs: 2 (triage 2)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 0.95}`×2
* Max tokens: `3072`×2
* Reasoning budget: `512`×2
* `ngl`: `0`×2
* `ctx`: `4096`×2
* System message: `None`×2
* `command` recorded: 2 of 2 runs
* Server log(s): `raw/budget-qwen35-4b-think-rb512/server.log`

### `qwen35-9b-think-rb1024`

* Runs: 1 (triage 1)
* Sampling: `{"min_p": 0.0, "presence_penalty": 1.5, "repeat_penalty": 1.0, "temperature": 1.0, "top_k": 20, "top_p": 0.95}`×1
* Max tokens: `3072`×1
* Reasoning budget: `1024`×1
* `ngl`: `0`×1
* `ctx`: `4096`×1
* System message: `None`×1
* `command` recorded: 1 of 1 runs
* Server log(s): `raw/stageA2b-qwen35-9b-think-rb1024/server.log`

### `smollm3-3b-nonthink`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 0.6, "top_k": 0, "top_p": 0.95}`×3
* Max tokens: `2048`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `/no_think`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA2b2-smollm3-3b-nonthink/server.log`

### `smollm3-3b-think-rb1024`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 0.6, "top_k": 0, "top_p": 0.95}`×3
* Max tokens: `3072`×3
* Reasoning budget: `1024`×3
* `ngl`: `0`×3
* `ctx`: `4096`×3
* System message: `/think`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/stageA2b2-smollm3-3b-think-rb1024/server.log`

### `tinyllama-11b-chat`

* Runs: 3 (triage 3)
* Sampling: `{"min_p": 0.0, "presence_penalty": 0.0, "repeat_penalty": 1.0, "temperature": 0.7, "top_k": 50, "top_p": 0.95}`×3
* Max tokens: `768`×3
* Reasoning budget: `None`×3
* `ngl`: `0`×3
* `ctx`: `2048`×3
* System message: `None`×3
* `command` recorded: 3 of 3 runs
* Server log(s): `raw/tinyllama-tinyllama-11b-chat/server.log`

## 3. Per-run values

| run_code | label | case | draw | seed | finish | prompt_tok | compl_tok | cap | budget | wall_s | rss_mb | vram_mb | source |
|---|---|---|---:|---:|---|---:|---:|---:|---:|---:|---:|---:|---|
| `00gx1xfl` | `abl-lfm120` | C09 | 1 | 4101 | `stop` | 225 | 43 | 512 | None | 1.647 | 1148.9 | 1917 | `results/ablation_corpus.json`[16] |
| `02ezb90h` | `corpus-4b-nonthink` | C06 | 1 | 4101 | `stop` | 233 | 50 | 768 | None | 7.611 | 3556.4 | 1935 | `results/corpus.json`[15] |
| `03b4a74n` | `corpus-lfm120-instruct` | C02 | 2 | 4102 | `stop` | 311 | 51 | 512 | None | 1.684 | 1098.9 | 1891 | `results/corpus.json`[58] |
| `05f8vfej` | `corpus-lfm120-instruct` | C17 | 1 | 4101 | `stop` | 248 | 49 | 512 | None | 1.791 | 1111.3 | 2117 | `results/corpus.json`[102] |
| `06bxkqx9` | `abl-4b-nonthink` | C02 | 2 | 4102 | `stop` | 300 | 113 | 768 | None | 16.433 | 3400.9 | 1959 | `results/ablation_corpus_abl-4b-nonthink.json`[3] |
| `09z3csf9` | `lfm120-ins-t08-rep10` | — | 2 | 4102 | `stop` | 244 | 73 | 2048 | None | 2.569 | 1098.3 | 2025 | `results/stageB.json`[10] |
| `0bil3r3p` | `abl-lfm3-think-512` | — | 1 | 4101 | `stop` | 238 | 634 | 1024 | 512 | 30.993 | 1558.6 | 1886 | `results/ablation_triage.json`[12] |
| `0cpv8l91` | `corpus-tinyllama` | C16 | 2 | 4102 | `length` | 230 | 768 | 768 | None | 23.441 | 1520.4 | 2076 | `results/corpus.json`[230] |
| `0cyhnkqq` | `lfm3-17b-nonthink` | — | 2 | 4102 | `stop` | 239 | 107 | 3072 | None | 5.462 | 1559.3 | 1918 | `results/stageA2a.json`[16] |
| `0ewnc62x` | `corpus-lfm3-think-512` | C17 | 1 | 4101 | `stop` | 245 | 616 | 1024 | 512 | 30.667 | 2148.3 | 2080 | `results/corpus.json`[140] |
| `0i6ybmif` | `corpus-4b-nonthink` | C11 | 3 | 4103 | `stop` | 241 | 50 | 768 | None | 6.836 | 3681.5 | 1935 | `results/corpus.json`[32] |
| `0keg5xbt` | `corpus-4b-nonthink` | C06 | 3 | 4103 | `stop` | 233 | 47 | 768 | None | 6.41 | 3556.4 | 1935 | `results/corpus.json`[17] |
| `0nepnebi` | `lfm120-ins-r2cap` | — | 3 | 4103 | `stop` | 244 | 74 | 400 | None | 2.45 | 1070.8 | 2027 | `results/stageC.json`[5] |
| `0oyt8v77` | `corpus-4b-nonthink` | C13 | 2 | 4102 | `stop` | 251 | 56 | 768 | None | 7.658 | 3681.7 | 1935 | `results/corpus.json`[37] |
| `0rhp3xav` | `qwen35-2b-nonthink` | — | 3 | 4103 | `stop` | 230 | 110 | 2048 | None | 6.569 | 1821.1 | 1895 | `results/stageA.json`[14] |
| `0vrcdu0c` | `corpus-4b-nonthink` | C08 | 1 | 4101 | `stop` | 196 | 55 | 768 | None | 8.248 | 3556.4 | 1935 | `results/corpus.json`[21] |
| `0yds3kny` | `corpus-lfm3-think-512` | C02 | 2 | 4102 | `stop` | 306 | 665 | 1024 | 512 | 32.607 | 1613.4 | 2040 | `results/corpus.json`[111] |
| `12fsrtxj` | `corpus-lfm120-t08` | C03 | 1 | 4101 | `stop` | 214 | 109 | 512 | None | 3.748 | 1103.1 | 2027 | `results/corpus.json`[166] |
| `16857ex3` | `corpus-lfm120-t08` | C11 | 1 | 4101 | `stop` | 245 | 68 | 512 | None | 2.418 | 1109.9 | 2029 | `results/corpus.json`[182] |
| `1fpf8tid` | `corpus-4b-nonthink` | C09 | 2 | 4102 | `stop` | 218 | 2 | 768 | None | 0.324 | 3556.4 | 1935 | `results/corpus.json`[25] |
| `1gsjt1so` | `corpus-4b-nonthink` | C15 | 1 | 4101 | `stop` | 249 | 48 | 768 | None | 7.323 | 3681.7 | 1935 | `results/corpus.json`[42] |
| `1hbt39r2` | `corpus-lfm120-t08` | C08 | 1 | 4101 | `stop` | 201 | 60 | 512 | None | 2.128 | 1103.3 | 2027 | `results/corpus.json`[176] |
| `1hfh8m9p` | `abl-4b-nonthink` | C15 | 2 | 4102 | `length` | 249 | 768 | 768 | None | 111.666 | 4122.3 | 2027 | `results/ablation_corpus_abl-4b-nonthink.json`[29] |
| `1lc3q5h8` | `corpus-lfm120-instruct` | C03 | 1 | 4101 | `stop` | 214 | 86 | 512 | None | 3.137 | 1103.9 | 1891 | `results/corpus.json`[60] |
| `1m7pzcv9` | `corpus-4b-nonthink` | C14 | 3 | 4103 | `stop` | 256 | 2 | 768 | None | 0.372 | 3681.7 | 1935 | `results/corpus.json`[41] |
| `1ozvnggg` | `smollm3-3b-nonthink` | — | 3 | 4103 | `stop` | 287 | 59 | 2048 | None | 5.005 | 2620.2 | 1964 | `results/stageA2b2.json`[5] |
| `1u67za0x` | `corpus-lfm3-think-512` | C14 | 1 | 4101 | `length` | 258 | 1024 | 1024 | 512 | 50.88 | 2040.5 | 2080 | `results/corpus.json`[134] |
| `1w471qee` | `abl-lfm120` | C02 | 1 | 4101 | `stop` | 311 | 75 | 512 | None | 2.813 | 1143.1 | 1917 | `results/ablation_corpus.json`[2] |
| `1zc1ycfe` | `corpus-4b-think-512` | C15 | 1 | 4101 | `stop` | 247 | 538 | 1024 | 512 | 81.788 | 5294.2 | 2193 | `results/corpus.json`[158] |
| `21wf8org` | `abl-lfm120` | C17 | 2 | 4102 | `stop` | 248 | 79 | 512 | None | 2.707 | 1165.7 | 1919 | `results/ablation_corpus.json`[33] |
| `23c75e3c` | `corpus-lfm120-t08` | C15 | 2 | 4102 | `stop` | 256 | 61 | 512 | None | 2.011 | 1110.0 | 2029 | `results/corpus.json`[191] |
| `26ukqanl` | `corpus-lfm120-t08` | C02 | 2 | 4102 | `stop` | 311 | 45 | 512 | None | 1.495 | 1097.8 | 2027 | `results/corpus.json`[165] |
| `29asv4kf` | `corpus-lfm3-think-512` | C10 | 1 | 4101 | `length` | 467 | 1024 | 1024 | 512 | 51.475 | 1874.3 | 2080 | `results/corpus.json`[126] |
| `2b15p3wn` | `corpus-4b-nonthink` | C10 | 1 | 4101 | `stop` | 462 | 79 | 768 | None | 11.707 | 3556.4 | 1935 | `results/corpus.json`[27] |
| `2beopira` | `abl-4b-nonthink` | C09 | 1 | 4101 | `stop` | 218 | 2 | 768 | None | 0.989 | 3636.6 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[16] |
| `2gaalwf2` | `abl-4b-think-512` | — | 2 | 4102 | `stop` | 228 | 612 | 1024 | 512 | 86.159 | 3533.4 | 1940 | `results/ablation_triage.json`[4] |
| `2hmmuyhd` | `corpus-4b-nonthink` | C18 | 3 | 4103 | `stop` | 212 | 54 | 768 | None | 7.416 | 3681.7 | 1935 | `results/corpus.json`[53] |
| `2jlxi0vq` | `abl-4b-nonthink` | C15 | 1 | 4101 | `stop` | 249 | 48 | 768 | None | 7.537 | 4122.3 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[28] |
| `2ktblesx` | `corpus-tinyllama` | C07 | 1 | 4101 | `length` | 265 | 768 | 768 | None | 23.882 | 1209.0 | 1901 | `results/corpus.json`[211] |
| `2l1gi5ta` | `corpus-lfm120-instruct` | C16 | 2 | 4102 | `stop` | 198 | 54 | 512 | None | 1.788 | 1111.3 | 2117 | `results/corpus.json`[100] |
| `2mc6gqwt` | `lfm350-official` | — | 1 | 4101 | `stop` | 244 | 16 | 2048 | None | 0.26 | 610.7 | 1836 | `results/stageA.json`[3] |
| `2mltxy3h` | `corpus-4b-nonthink` | C11 | 2 | 4102 | `stop` | 241 | 51 | 768 | None | 7.002 | 3681.5 | 1935 | `results/corpus.json`[31] |
| `2rfx87bo` | `corpus-lfm3-think-512` | C16 | 2 | 4102 | `stop` | 195 | 582 | 1024 | 512 | 28.928 | 2131.8 | 2080 | `results/corpus.json`[139] |
| `2sqh27z3` | `abl-lfm120` | — | 1 | 4101 | `stop` | 244 | 68 | 512 | None | 2.533 | 1136.7 | 1892 | `results/ablation_triage.json`[9] |
| `2t24lbnf` | `abl-2b-nonthink` | C13 | 2 | 4102 | `stop` | 251 | 55 | 768 | None | 3.294 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[25] |
| `2wd8k7zq` | `corpus-4b-nonthink` | C09 | 1 | 4101 | `stop` | 218 | 2 | 768 | None | 0.964 | 3556.4 | 1935 | `results/corpus.json`[24] |
| `2yxoeoxi` | `qwen35-2b-think-rb1024` | — | 1 | 4101 | `stop` | 228 | 1111 | 3072 | 1024 | 67.165 | 1822.0 | 1931 | `results/stageA2a.json`[18] |
| `34x9syv7` | `corpus-4b-think-512` | C13 | 1 | 4101 | `stop` | 249 | 558 | 1024 | 512 | 79.277 | 5030.0 | 2193 | `results/corpus.json`[156] |
| `37s3bhg9` | `corpus-lfm120-instruct` | C07 | 3 | 4103 | `stop` | 230 | 41 | 512 | None | 1.4 | 1104.0 | 2117 | `results/corpus.json`[74] |
| `3ad3exbj` | `qwen35-2b-nonthink` | C06 | 2 | 4102 | `stop` | 233 | 52 | 2048 | None | 3.136 | 1912.9 | 1929 | `results/control_corpus_2b.json`[11] |
| `3c2rkui7` | `corpus-lfm120-instruct` | C11 | 3 | 4103 | `stop` | 245 | 62 | 512 | None | 2.031 | 1111.3 | 2117 | `results/corpus.json`[86] |
| `3cozjltm` | `qwen35-2b-nonthink-round2` | — | 3 | 4103 | `stop` | 230 | 100 | 2048 | None | 6.222 | 1821.5 | 2027 | `results/stageB.json`[17] |
| `3et30fj1` | `abl-4b-nonthink` | C02 | 1 | 4101 | `stop` | 300 | 116 | 768 | None | 16.684 | 3400.9 | 1959 | `results/ablation_corpus_abl-4b-nonthink.json`[2] |
| `3f92l6va` | `corpus-4b-nonthink` | C12 | 3 | 4103 | `stop` | 236 | 2 | 768 | None | 0.344 | 3681.5 | 1935 | `results/corpus.json`[35] |
| `3nbtpxch` | `abl-lfm3-think-512` | — | 2 | 4102 | `stop` | 238 | 629 | 1024 | 512 | 30.751 | 1577.7 | 1886 | `results/ablation_triage.json`[13] |
| `3o05kymx` | `qwen35-2b-nonthink` | C13 | 1 | 4101 | `stop` | 251 | 107 | 2048 | None | 6.713 | 2053.0 | 1930 | `results/control_corpus_2b.json`[24] |
| `3penv5qb` | `qwen35-2b-nonthink` | — | 2 | 4102 | `stop` | 230 | 100 | 2048 | None | 5.977 | 1821.0 | 1895 | `results/stageA.json`[13] |
| `3pqxkyyl` | `corpus-lfm120-instruct` | C18 | 2 | 4102 | `stop` | 215 | 49 | 512 | None | 1.621 | 1111.3 | 2117 | `results/corpus.json`[106] |
| `3sd6o8he` | `corpus-tinyllama` | C11 | 1 | 4101 | `length` | 288 | 768 | 768 | None | 23.883 | 1348.8 | 2076 | `results/corpus.json`[219] |
| `3v6qp23o` | `corpus-lfm3-think-512` | C12 | 1 | 4101 | `stop` | 238 | 586 | 1024 | 512 | 28.981 | 1960.2 | 2080 | `results/corpus.json`[130] |
| `3y1lkt18` | `lfm3-17b-think-rb512` | — | 2 | 4102 | `stop` | 238 | 628 | 3072 | 512 | 32.139 | 1577.3 | 2019 | `results/budget.json`[1] |
| `3z5j5ig5` | `abl-lfm120` | C05 | 2 | 4102 | `stop` | 233 | 40 | 512 | None | 1.374 | 1148.9 | 1917 | `results/ablation_corpus.json`[9] |
| `40k74mdf` | `abl-4b-nonthink` | C07 | 2 | 4102 | `stop` | 226 | 70 | 768 | None | 9.718 | 3636.6 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[13] |
| `41ewsrna` | `corpus-lfm120-instruct` | C01 | 2 | 4102 | `stop` | 236 | 43 | 512 | None | 1.418 | 1098.0 | 1891 | `results/corpus.json`[55] |
| `41w4gihk` | `corpus-4b-nonthink` | C02 | 1 | 4101 | `stop` | 300 | 51 | 768 | None | 7.812 | 3438.6 | 1929 | `results/corpus.json`[3] |
| `42g6gw7q` | `corpus-lfm3-think-512` | C02 | 1 | 4101 | `stop` | 306 | 584 | 1024 | 512 | 28.669 | 1595.1 | 2040 | `results/corpus.json`[110] |
| `45odogxz` | `abl-4b-nonthink` | C06 | 1 | 4101 | `stop` | 233 | 42 | 768 | None | 6.462 | 3636.5 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[10] |
| `45tdof1p` | `abl-2b-nonthink` | C15 | 2 | 4102 | `stop` | 249 | 66 | 768 | None | 3.966 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[29] |
| `47monrko` | `corpus-tinyllama` | C08 | 2 | 4102 | `stop` | 235 | 724 | 768 | None | 22.084 | 1264.7 | 1901 | `results/corpus.json`[214] |
| `4bmioiwq` | `abl-4b-nonthink` | C18 | 2 | 4102 | `stop` | 212 | 48 | 768 | None | 6.518 | 4261.0 | 2027 | `results/ablation_corpus_abl-4b-nonthink.json`[35] |
| `4btjqk9n` | `abl-lfm120` | C08 | 2 | 4102 | `stop` | 201 | 3 | 512 | None | 0.129 | 1148.9 | 1917 | `results/ablation_corpus.json`[15] |
| `4dd1exgl` | `qwen35-2b-nonthink` | C03 | 2 | 4102 | `stop` | 207 | 67 | 2048 | None | 4.041 | 1866.7 | 1927 | `results/control_corpus_2b.json`[5] |
| `4hfrgqwl` | `qwen3-17b-nonthink` | — | 1 | 4101 | `stop` | 231 | 53 | 2048 | None | 2.927 | 1965.6 | 1925 | `results/stageA2a.json`[9] |
| `4j2xw619` | `corpus-lfm120-instruct` | C12 | 3 | 4103 | `stop` | 241 | 50 | 512 | None | 1.643 | 1111.3 | 2117 | `results/corpus.json`[89] |
| `4posvu7n` | `corpus-lfm120-instruct` | C14 | 1 | 4101 | `stop` | 261 | 53 | 512 | None | 1.913 | 1111.3 | 2117 | `results/corpus.json`[93] |
| `4qbwuujq` | `lfm120-think-official` | — | 3 | 4103 | `stop` | 244 | 1717 | 3072 | None | 61.711 | 1213.1 | 1956 | `results/stageA.json`[8] |
| `4vacdiea` | `corpus-lfm120-instruct` | C16 | 3 | 4103 | `stop` | 198 | 54 | 512 | None | 1.756 | 1111.3 | 2117 | `results/corpus.json`[101] |
| `4x6m1xv8` | `abl-lfm120` | C01 | 1 | 4101 | `stop` | 236 | 100 | 512 | None | 3.642 | 1136.7 | 1917 | `results/ablation_corpus.json`[0] |
| `504okf6p` | `abl-lfm3-think-512` | — | 3 | 4103 | `stop` | 238 | 650 | 1024 | 512 | 32.151 | 1595.5 | 2074 | `results/ablation_triage.json`[14] |
| `511ergmp` | `lfm120-instruct-official-gpufull` | — | 2 | 4102 | `stop` | 244 | 86 | 768 | None | 0.649 | 1071.4 | 2741 | `results/gpu-full.json`[1] |
| `54hcnvhy` | `abl-2b-nonthink` | C14 | 1 | 4101 | `stop` | 256 | 39 | 768 | None | 2.711 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[26] |
| `54qq3v3w` | `corpus-lfm3-think-512` | C01 | 2 | 4102 | `stop` | 233 | 605 | 1024 | 512 | 29.528 | 1577.5 | 2038 | `results/corpus.json`[109] |
| `58zrn1ip` | `corpus-4b-nonthink` | C05 | 2 | 4102 | `stop` | 232 | 55 | 768 | None | 7.558 | 3556.4 | 1935 | `results/corpus.json`[13] |
| `5apuk073` | `abl-lfm120` | — | 2 | 4102 | `stop` | 244 | 67 | 512 | None | 2.242 | 1136.8 | 1892 | `results/ablation_triage.json`[10] |
| `5gxyuv5e` | `corpus-lfm3-think-512` | C04 | 2 | 4102 | `length` | 229 | 1024 | 1024 | 512 | 50.138 | 1683.6 | 2040 | `results/corpus.json`[115] |
| `5iqroi0i` | `lfm3-17b-nonthink` | — | 3 | 4103 | `stop` | 239 | 106 | 3072 | None | 5.379 | 1559.4 | 1918 | `results/stageA2a.json`[17] |
| `5ufc25wv` | `corpus-4b-nonthink` | C16 | 2 | 4102 | `stop` | 192 | 55 | 768 | None | 7.54 | 3681.7 | 1935 | `results/corpus.json`[46] |
| `5vfwqmxm` | `corpus-lfm120-t08` | C14 | 1 | 4101 | `stop` | 261 | 97 | 512 | None | 3.398 | 1110.0 | 2029 | `results/corpus.json`[188] |
| `5ybm0tgq` | `qwen35-2b-nonthink` | C12 | 1 | 4101 | `stop` | 236 | 101 | 2048 | None | 6.55 | 2052.5 | 1930 | `results/control_corpus_2b.json`[22] |
| `5zpwsnop` | `abl-4b-nonthink` | C10 | 2 | 4102 | `stop` | 462 | 117 | 768 | None | 16.561 | 3636.6 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[19] |
| `62c7wa5p` | `abl-lfm120` | C03 | 1 | 4101 | `stop` | 214 | 113 | 512 | None | 4.042 | 1148.8 | 1917 | `results/ablation_corpus.json`[4] |
| `6bwwbtbx` | `corpus-lfm120-t08` | C07 | 2 | 4102 | `stop` | 230 | 65 | 512 | None | 2.153 | 1103.3 | 2027 | `results/corpus.json`[175] |
| `6fjzbkzn` | `corpus-lfm3-think-512` | C03 | 2 | 4102 | `stop` | 211 | 587 | 1024 | 512 | 28.716 | 1648.8 | 2040 | `results/corpus.json`[113] |
| `6ha7sry6` | `abl-4b-nonthink` | C14 | 2 | 4102 | `stop` | 256 | 87 | 768 | None | 13.057 | 4004.7 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[27] |
| `6hqepjgc` | `corpus-4b-nonthink` | C02 | 3 | 4103 | `stop` | 300 | 56 | 768 | None | 7.69 | 3438.6 | 1929 | `results/corpus.json`[5] |
| `6it52hen` | `abl-4b-nonthink` | C06 | 2 | 4102 | `stop` | 233 | 52 | 768 | None | 7.158 | 3636.5 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[11] |
| `6ixb61t6` | `corpus-lfm120-t08` | C04 | 1 | 4101 | `stop` | 232 | 61 | 512 | None | 2.199 | 1103.2 | 2027 | `results/corpus.json`[168] |
| `6kq22wrl` | `corpus-lfm120-instruct` | C05 | 3 | 4103 | `stop` | 233 | 46 | 512 | None | 1.527 | 1104.0 | 1891 | `results/corpus.json`[68] |
| `6oxshvbw` | `abl-2b-nonthink` | C11 | 1 | 4101 | `stop` | 241 | 60 | 768 | None | 3.892 | 2037.4 | 1930 | `results/ablation_corpus_abl-2b-nonthink.json`[20] |
| `6p5mhpuh` | `abl-2b-nonthink` | C08 | 2 | 4102 | `stop` | 196 | 96 | 768 | None | 5.852 | 1988.7 | 1930 | `results/ablation_corpus_abl-2b-nonthink.json`[15] |
| `6y83mbtz` | `lfm120-ins-t08` | — | 1 | 4101 | `stop` | 244 | 88 | 2048 | None | 3.333 | 1097.5 | 2025 | `results/stageB.json`[6] |
| `6yiyxbb7` | `corpus-lfm120-instruct` | C14 | 3 | 4103 | `stop` | 261 | 50 | 512 | None | 1.683 | 1111.3 | 2117 | `results/corpus.json`[95] |
| `72in240n` | `abl-2b-nonthink` | C05 | 2 | 4102 | `stop` | 232 | 86 | 768 | None | 5.155 | 1898.5 | 1929 | `results/ablation_corpus_abl-2b-nonthink.json`[9] |
| `758iu45k` | `corpus-lfm120-t08` | C10 | 1 | 4101 | `stop` | 470 | 96 | 512 | None | 3.449 | 1103.3 | 2029 | `results/corpus.json`[180] |
| `76nsusg0` | `abl-4b-nonthink` | — | 3 | 4103 | `stop` | 230 | 89 | 768 | None | 12.506 | 3400.1 | 1934 | `results/ablation_triage.json`[2] |
| `7byj4wm6` | `abl-2b-nonthink` | C17 | 2 | 4102 | `stop` | 239 | 59 | 768 | None | 3.541 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[33] |
| `7cec29gh` | `qwen35-2b-nonthink` | C17 | 2 | 4102 | `stop` | 239 | 116 | 2048 | None | 6.976 | 2142.8 | 1930 | `results/control_corpus_2b.json`[33] |
| `7faz3cal` | `corpus-4b-nonthink` | C02 | 2 | 4102 | `stop` | 300 | 47 | 768 | None | 6.454 | 3438.6 | 1929 | `results/corpus.json`[4] |
| `7fb3h9z1` | `corpus-lfm120-t08` | C14 | 2 | 4102 | `stop` | 261 | 52 | 512 | None | 1.715 | 1110.0 | 2029 | `results/corpus.json`[189] |
| `7l1ntpqd` | `corpus-lfm3-think-512` | C15 | 2 | 4102 | `stop` | 253 | 567 | 1024 | 512 | 28.259 | 2099.0 | 2080 | `results/corpus.json`[137] |
| `7n68wfwm` | `corpus-lfm120-instruct` | C09 | 1 | 4101 | `stop` | 225 | 46 | 512 | None | 1.677 | 1104.0 | 2117 | `results/corpus.json`[78] |
| `7oci2ogc` | `corpus-lfm120-t08` | C01 | 2 | 4102 | `stop` | 236 | 39 | 512 | None | 1.286 | 1097.7 | 2027 | `results/corpus.json`[163] |
| `7pvmpy3z` | `corpus-lfm120-instruct` | C02 | 1 | 4101 | `stop` | 311 | 72 | 512 | None | 2.547 | 1098.9 | 1891 | `results/corpus.json`[57] |
| `7rzgyjuu` | `abl-lfm120` | C11 | 2 | 4102 | `stop` | 245 | 159 | 512 | None | 5.599 | 1155.4 | 1919 | `results/ablation_corpus.json`[21] |
| `7tkftty2` | `qwen35-2b-nonthink` | C14 | 1 | 4101 | `stop` | 256 | 63 | 2048 | None | 4.088 | 2098.0 | 1930 | `results/control_corpus_2b.json`[26] |
| `80m5wks3` | `qwen35-2b-nonthink-round2` | — | 1 | 4101 | `stop` | 230 | 108 | 2048 | None | 7.114 | 1821.3 | 2027 | `results/stageB.json`[15] |
| `82ozq269` | `qwen35-2b-nonthink` | C02 | 1 | 4101 | `stop` | 300 | 71 | 2048 | None | 4.595 | 1821.9 | 1925 | `results/control_corpus_2b.json`[2] |
| `86sbp4tw` | `qwen35-2b-think-rb1024` | — | 3 | 4103 | `stop` | 228 | 1102 | 3072 | 1024 | 66.237 | 1936.3 | 1931 | `results/stageA2a.json`[20] |
| `89i18iwr` | `qwen35-2b-nonthink` | C08 | 1 | 4101 | `stop` | 196 | 56 | 2048 | None | 3.632 | 1959.7 | 1930 | `results/control_corpus_2b.json`[14] |
| `8akn5dd9` | `corpus-tinyllama` | C18 | 1 | 4101 | `length` | 254 | 768 | 768 | None | 23.793 | 1566.9 | 2076 | `results/corpus.json`[233] |
| `8gyag2hk` | `corpus-lfm120-instruct` | C17 | 3 | 4103 | `stop` | 248 | 46 | 512 | None | 1.528 | 1111.3 | 2117 | `results/corpus.json`[104] |
| `8kwk538v` | `corpus-4b-nonthink` | C06 | 2 | 4102 | `stop` | 233 | 42 | 768 | None | 5.802 | 3556.4 | 1935 | `results/corpus.json`[16] |
| `8q1maguo` | `corpus-lfm3-think-512` | C05 | 2 | 4102 | `stop` | 230 | 623 | 1024 | 512 | 30.539 | 1723.2 | 2040 | `results/corpus.json`[117] |
| `8urvt69a` | `qwen35-4b-think-rb1024` | — | 3 | 4103 | `stop` | 228 | 1154 | 3072 | 1024 | 163.861 | 3738.5 | 1919 | `results/stageA2a.json`[5] |
| `8wwyrst4` | `qwen35-2b-nonthink` | C05 | 2 | 4102 | `stop` | 232 | 60 | 2048 | None | 3.649 | 1912.9 | 1929 | `results/control_corpus_2b.json`[9] |
| `93ffxfwx` | `lfm120-ins-t08-rep10` | — | 1 | 4101 | `stop` | 244 | 69 | 2048 | None | 2.653 | 1098.2 | 2025 | `results/stageB.json`[9] |
| `9j8duvoz` | `corpus-tinyllama` | C13 | 2 | 4102 | `length` | 300 | 768 | 768 | None | 23.632 | 1422.5 | 2076 | `results/corpus.json`[224] |
| `9ld4grh8` | `abl-2b-nonthink` | C02 | 2 | 4102 | `stop` | 300 | 70 | 768 | None | 4.238 | 1807.6 | 1926 | `results/ablation_corpus_abl-2b-nonthink.json`[3] |
| `9ml5rc53` | `corpus-4b-nonthink` | C14 | 1 | 4101 | `stop` | 256 | 47 | 768 | None | 7.123 | 3681.7 | 1935 | `results/corpus.json`[39] |
| `9ojbk8ek` | `corpus-lfm120-t08` | C12 | 1 | 4101 | `stop` | 241 | 57 | 512 | None | 2.053 | 1109.9 | 2029 | `results/corpus.json`[184] |
| `9qox7nau` | `abl-lfm120` | C01 | 2 | 4102 | `stop` | 236 | 111 | 512 | None | 3.786 | 1136.9 | 1917 | `results/ablation_corpus.json`[1] |
| `9vq2j9ou` | `corpus-lfm120-instruct` | C02 | 3 | 4103 | `stop` | 311 | 73 | 512 | None | 2.449 | 1098.9 | 1891 | `results/corpus.json`[59] |
| `9xvqyyas` | `lfm120-ins-t03` | — | 3 | 4103 | `stop` | 244 | 81 | 2048 | None | 2.89 | 1098.0 | 2025 | `results/stageB.json`[2] |
| `9xyh7sj0` | `corpus-4b-think-512` | C14 | 1 | 4101 | `stop` | 254 | 568 | 1024 | 512 | 82.73 | 5161.8 | 2193 | `results/corpus.json`[157] |
| `a69orocq` | `qwen35-2b-nonthink` | C09 | 1 | 4101 | `stop` | 218 | 83 | 2048 | None | 5.395 | 1959.7 | 1930 | `results/control_corpus_2b.json`[16] |
| `a72l8hia` | `abl-lfm120` | C12 | 2 | 4102 | `stop` | 241 | 41 | 512 | None | 1.417 | 1160.6 | 1919 | `results/ablation_corpus.json`[23] |
| `a87200a4` | `qwen35-2b-nonthink` | C16 | 2 | 4102 | `stop` | 192 | 60 | 2048 | None | 3.632 | 2142.8 | 1930 | `results/control_corpus_2b.json`[31] |
| `a9ahi29c` | `abl-4b-nonthink` | C03 | 2 | 4102 | `stop` | 207 | 42 | 768 | None | 5.803 | 3520.5 | 1963 | `results/ablation_corpus_abl-4b-nonthink.json`[5] |
| `ab5dkajq` | `abl-2b-nonthink` | C03 | 1 | 4101 | `stop` | 207 | 158 | 768 | None | 9.669 | 1852.7 | 1928 | `results/ablation_corpus_abl-2b-nonthink.json`[4] |
| `am5dyzl7` | `corpus-tinyllama` | C17 | 2 | 4102 | `length` | 288 | 768 | 768 | None | 23.616 | 1543.0 | 2076 | `results/corpus.json`[232] |
| `am9g546q` | `corpus-lfm120-instruct` | C09 | 2 | 4102 | `stop` | 225 | 46 | 512 | None | 1.529 | 1104.0 | 2117 | `results/corpus.json`[79] |
| `ampeuiri` | `corpus-tinyllama` | C07 | 2 | 4102 | `length` | 265 | 768 | 768 | None | 23.642 | 1232.3 | 1901 | `results/corpus.json`[212] |
| `apc56f1h` | `corpus-lfm120-t08` | C18 | 1 | 4101 | `stop` | 215 | 41 | 512 | None | 1.515 | 1110.0 | 2029 | `results/corpus.json`[196] |
| `arqhftrm` | `corpus-4b-nonthink` | C17 | 2 | 4102 | `stop` | 239 | 2 | 768 | None | 0.339 | 3681.7 | 1935 | `results/corpus.json`[49] |
| `awm6ptwx` | `corpus-tinyllama` | C03 | 2 | 4102 | `stop` | 247 | 370 | 768 | None | 11.046 | 1121.0 | 1901 | `results/corpus.json`[204] |
| `axfzx87w` | `abl-lfm120` | C18 | 2 | 4102 | `stop` | 215 | 135 | 512 | None | 4.584 | 1165.7 | 1919 | `results/ablation_corpus.json`[35] |
| `b0elqo5f` | `corpus-4b-think-512` | C04 | 1 | 4101 | `stop` | 226 | 577 | 1024 | 512 | 82.261 | 3836.6 | 2125 | `results/corpus.json`[147] |
| `b8u2405u` | `qwen35-2b-nonthink` | C06 | 1 | 4101 | `stop` | 233 | 78 | 2048 | None | 4.995 | 1912.9 | 1929 | `results/control_corpus_2b.json`[10] |
| `b9fbvwtq` | `qwen35-2b-nonthink` | C13 | 2 | 4102 | `stop` | 251 | 114 | 2048 | None | 6.869 | 2053.0 | 1930 | `results/control_corpus_2b.json`[25] |
| `bcsr81h1` | `corpus-4b-nonthink` | C15 | 3 | 4103 | `stop` | 249 | 54 | 768 | None | 7.393 | 3681.7 | 1935 | `results/corpus.json`[44] |
| `bg2hck2t` | `corpus-4b-think-512` | C05 | 1 | 4101 | `stop` | 230 | 563 | 1024 | 512 | 79.965 | 3968.4 | 2125 | `results/corpus.json`[148] |
| `blxhrh3v` | `lfm120-ins-t03` | — | 2 | 4102 | `stop` | 244 | 74 | 2048 | None | 2.622 | 1097.9 | 2025 | `results/stageB.json`[1] |
| `bm99eisc` | `corpus-lfm120-instruct` | C11 | 1 | 4101 | `stop` | 245 | 52 | 512 | None | 2.054 | 1111.3 | 2117 | `results/corpus.json`[84] |
| `bmsswywm` | `qwen35-4b-think-rb512` | — | 2 | 4102 | `stop` | 228 | 629 | 3072 | 512 | 88.709 | 3570.7 | 2033 | `results/all_triage.json`[50] |
| `boa2v2mz` | `corpus-4b-nonthink` | C04 | 1 | 4101 | `stop` | 228 | 98 | 768 | None | 14.201 | 3556.4 | 1935 | `results/corpus.json`[9] |
| `boyxujrl` | `qwen35-2b-nonthink` | C07 | 1 | 4101 | `stop` | 226 | 291 | 2048 | None | 17.735 | 1912.9 | 1930 | `results/control_corpus_2b.json`[12] |
| `bq3oq9l1` | `corpus-lfm120-t08` | C03 | 2 | 4102 | `stop` | 214 | 83 | 512 | None | 2.735 | 1103.2 | 2027 | `results/corpus.json`[167] |
| `bxtos15y` | `abl-lfm120` | C02 | 2 | 4102 | `stop` | 311 | 131 | 512 | None | 4.474 | 1143.1 | 1917 | `results/ablation_corpus.json`[3] |
| `bz375bog` | `corpus-lfm3-think-512` | C18 | 1 | 4101 | `stop` | 212 | 602 | 1024 | 512 | 29.464 | 2182.8 | 2135 | `results/corpus.json`[142] |
| `c04enw6f` | `corpus-lfm3-think-512` | C11 | 2 | 4102 | `stop` | 242 | 633 | 1024 | 512 | 31.452 | 1942.4 | 2080 | `results/corpus.json`[129] |
| `c1xrqyca` | `corpus-4b-nonthink` | C10 | 3 | 4103 | `stop` | 462 | 140 | 768 | None | 19.832 | 3556.4 | 1935 | `results/corpus.json`[29] |
| `c3e7yd80` | `corpus-lfm120-t08` | C16 | 1 | 4101 | `stop` | 198 | 49 | 512 | None | 1.803 | 1110.0 | 2029 | `results/corpus.json`[192] |
| `c478nkwh` | `lfm120-instruct-official-gpufull` | — | 1 | 4101 | `stop` | 244 | 70 | 768 | None | 0.598 | 1071.2 | 2741 | `results/gpu-full.json`[0] |
| `c47i09o7` | `corpus-4b-nonthink` | C11 | 1 | 4101 | `stop` | 241 | 45 | 768 | None | 6.951 | 3681.5 | 1935 | `results/corpus.json`[30] |
| `c9r75fdr` | `corpus-lfm120-t08` | C06 | 1 | 4101 | `stop` | 238 | 39 | 512 | None | 1.458 | 1103.3 | 2027 | `results/corpus.json`[172] |
| `cejjp61m` | `corpus-lfm120-instruct` | C08 | 3 | 4103 | `stop` | 201 | 44 | 512 | None | 1.438 | 1104.0 | 2117 | `results/corpus.json`[77] |
| `chkhhwv7` | `corpus-4b-nonthink` | C03 | 3 | 4103 | `stop` | 207 | 57 | 768 | None | 7.993 | 3556.4 | 1933 | `results/corpus.json`[8] |
| `chknnyf3` | `corpus-lfm120-t08` | C05 | 1 | 4101 | `stop` | 233 | 78 | 512 | None | 2.738 | 1103.2 | 2027 | `results/corpus.json`[170] |
| `ckr2t2al` | `smollm3-3b-think-rb1024` | — | 2 | 4102 | `stop` | 467 | 403 | 3072 | 1024 | 35.09 | 2687.4 | 1964 | `results/stageA2b2.json`[1] |
| `clo369ah` | `lfm350-official` | — | 2 | 4102 | `stop` | 244 | 16 | 2048 | None | 0.189 | 610.8 | 1836 | `results/stageA.json`[4] |
| `clpr1c77` | `abl-2b-nonthink` | C13 | 1 | 4101 | `stop` | 251 | 52 | 768 | None | 3.406 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[24] |
| `cq2nde1d` | `corpus-lfm120-t08` | C04 | 2 | 4102 | `stop` | 232 | 55 | 512 | None | 1.811 | 1103.2 | 2027 | `results/corpus.json`[169] |
| `cz5xa27b` | `abl-lfm120` | C06 | 1 | 4101 | `stop` | 238 | 37 | 512 | None | 1.466 | 1148.9 | 1917 | `results/ablation_corpus.json`[10] |
| `d74gtao9` | `corpus-lfm3-think-512` | C04 | 1 | 4101 | `stop` | 229 | 636 | 1024 | 512 | 31.181 | 1665.7 | 2040 | `results/corpus.json`[114] |
| `d8bdt78q` | `lfm120-ins-t05` | — | 2 | 4102 | `stop` | 244 | 70 | 2048 | None | 2.491 | 1098.0 | 2025 | `results/stageB.json`[4] |
| `da7y2ahx` | `corpus-4b-nonthink` | C01 | 1 | 4101 | `stop` | 232 | 62 | 768 | None | 9.231 | 3437.6 | 1929 | `results/corpus.json`[0] |
| `dc4xdwh8` | `abl-2b-nonthink` | C18 | 1 | 4101 | `stop` | 212 | 99 | 768 | None | 6.173 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[34] |
| `dcjcz8cq` | `abl-lfm120` | C07 | 2 | 4102 | `stop` | 230 | 32 | 512 | None | 1.093 | 1148.9 | 1917 | `results/ablation_corpus.json`[13] |
| `deah77im` | `abl-2b-nonthink` | C11 | 2 | 4102 | `stop` | 241 | 62 | 768 | None | 3.731 | 2037.4 | 1930 | `results/ablation_corpus_abl-2b-nonthink.json`[21] |
| `dgeofx5m` | `abl-4b-nonthink` | C04 | 1 | 4101 | `stop` | 228 | 63 | 768 | None | 9.34 | 3520.5 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[6] |
| `dgth442y` | `corpus-lfm120-t08` | C01 | 1 | 4101 | `stop` | 236 | 45 | 512 | None | 1.712 | 1097.6 | 2027 | `results/corpus.json`[162] |
| `dhwdvlcb` | `qwen35-2b-nonthink-r2cap` | — | 2 | 4102 | `stop` | 230 | 143 | 400 | None | 8.429 | 1794.1 | 2035 | `results/stageC.json`[1] |
| `dj2g470d` | `corpus-tinyllama` | C15 | 2 | 4102 | `stop` | 295 | 673 | 768 | None | 20.503 | 1475.8 | 2076 | `results/corpus.json`[228] |
| `djgsdmpd` | `corpus-lfm120-t08` | C13 | 1 | 4101 | `stop` | 258 | 36 | 512 | None | 1.465 | 1110.0 | 2029 | `results/corpus.json`[186] |
| `dmsfwlqb` | `lfm120-ins-t08` | — | 2 | 4102 | `stop` | 244 | 65 | 2048 | None | 2.291 | 1097.6 | 2025 | `results/stageB.json`[7] |
| `dmv1b02u` | `corpus-lfm120-instruct` | C04 | 2 | 4102 | `stop` | 232 | 63 | 512 | None | 2.089 | 1103.9 | 1891 | `results/corpus.json`[64] |
| `dnqqbgtd` | `abl-lfm120` | C15 | 1 | 4101 | `stop` | 256 | 86 | 512 | None | 3.128 | 1165.7 | 1919 | `results/ablation_corpus.json`[28] |
| `do80gjrz` | `abl-2b-nonthink` | C10 | 1 | 4101 | `stop` | 462 | 180 | 768 | None | 11.338 | 1989.5 | 1930 | `results/ablation_corpus_abl-2b-nonthink.json`[18] |
| `dq2we2kg` | `corpus-lfm120-instruct` | C03 | 2 | 4102 | `stop` | 214 | 84 | 512 | None | 2.742 | 1103.9 | 1891 | `results/corpus.json`[61] |
| `dsq5e8vd` | `tinyllama-11b-chat` | — | 1 | 4101 | `stop` | 268 | 305 | 768 | None | 9.315 | 1018.7 | 1901 | `results/tinyllama.json`[0] |
| `dtx2e3tv` | `abl-lfm120` | C13 | 2 | 4102 | `stop` | 258 | 132 | 512 | None | 4.487 | 1160.6 | 1919 | `results/ablation_corpus.json`[25] |
| `dw7i3mu7` | `abl-4b-nonthink` | C11 | 2 | 4102 | `stop` | 241 | 47 | 768 | None | 6.504 | 3760.9 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[21] |
| `dwzc2jgw` | `lfm3-17b-think-rb1024` | — | 2 | 4102 | `stop` | 238 | 1122 | 3072 | 1024 | 55.169 | 1584.5 | 1919 | `results/stageA2a.json`[13] |
| `dzeoymkr` | `corpus-lfm120-t08` | C16 | 2 | 4102 | `stop` | 198 | 41 | 512 | None | 1.362 | 1110.0 | 2029 | `results/corpus.json`[193] |
| `e2rnuks4` | `lfm120-think-official` | — | 2 | 4102 | `stop` | 244 | 2949 | 3072 | None | 109.477 | 1173.3 | 1956 | `results/stageA.json`[7] |
| `e3uhlrx6` | `qwen35-4b-think-rb1024` | — | 2 | 4102 | `stop` | 228 | 1144 | 3072 | 1024 | 162.379 | 3588.0 | 1919 | `results/stageA2a.json`[4] |
| `e40toqx9` | `abl-2b-nonthink` | C08 | 1 | 4101 | `stop` | 196 | 78 | 768 | None | 4.925 | 1988.7 | 1930 | `results/ablation_corpus_abl-2b-nonthink.json`[14] |
| `e5la6skq` | `qwen35-2b-nonthink` | C09 | 2 | 4102 | `stop` | 218 | 83 | 2048 | None | 4.989 | 1959.7 | 1930 | `results/control_corpus_2b.json`[17] |
| `e9hpmbq5` | `corpus-tinyllama` | C11 | 2 | 4102 | `stop` | 288 | 348 | 768 | None | 10.39 | 1372.7 | 2076 | `results/corpus.json`[220] |
| `ea4l84s1` | `abl-2b-nonthink` | C18 | 2 | 4102 | `stop` | 212 | 99 | 768 | None | 5.915 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[35] |
| `edonxrej` | `abl-4b-nonthink` | — | 2 | 4102 | `stop` | 230 | 86 | 768 | None | 11.987 | 3400.1 | 1934 | `results/ablation_triage.json`[1] |
| `edzy92ny` | `abl-2b-nonthink` | C09 | 2 | 4102 | `stop` | 218 | 60 | 768 | None | 3.615 | 1988.7 | 1930 | `results/ablation_corpus_abl-2b-nonthink.json`[17] |
| `ee3ukpn6` | `corpus-lfm120-instruct` | C10 | 3 | 4103 | `stop` | 470 | 108 | 512 | None | 3.749 | 1104.0 | 2117 | `results/corpus.json`[83] |
| `emuu8vl2` | `abl-2b-nonthink` | C16 | 1 | 4101 | `stop` | 192 | 159 | 768 | None | 9.744 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[30] |
| `ep2h4gc3` | `abl-lfm120` | C14 | 1 | 4101 | `stop` | 261 | 79 | 512 | None | 2.89 | 1165.7 | 1919 | `results/ablation_corpus.json`[26] |
| `er7kqvb1` | `corpus-lfm120-instruct` | C05 | 2 | 4102 | `stop` | 233 | 47 | 512 | None | 1.527 | 1103.9 | 1891 | `results/corpus.json`[67] |
| `euibewmb` | `corpus-lfm120-instruct` | C01 | 1 | 4101 | `stop` | 236 | 42 | 512 | None | 1.554 | 1097.9 | 1891 | `results/corpus.json`[54] |
| `evp6j9xa` | `corpus-lfm120-t08` | C09 | 2 | 4102 | `stop` | 225 | 21 | 512 | None | 0.693 | 1103.3 | 2027 | `results/corpus.json`[179] |
| `ewktjzvh` | `qwen35-4b-nonthink-gpufull` | — | 2 | 4102 | `stop` | 230 | 92 | 768 | None | 2.206 | 3265.7 | 4833 | `results/gpu-full-qwen35-4b-nonthink-gpufull.json`[1] |
| `ez3jx5az` | `lfm120-instruct-official` | — | 1 | 4101 | `stop` | 244 | 78 | 2048 | None | 2.788 | 1097.9 | 1870 | `results/stageA.json`[0] |
| `f27b1ggf` | `corpus-lfm120-t08` | C07 | 1 | 4101 | `stop` | 230 | 11 | 512 | None | 0.55 | 1103.3 | 2027 | `results/corpus.json`[174] |
| `f4rz3v6g` | `tinyllama-11b-chat` | — | 3 | 4103 | `stop` | 268 | 464 | 768 | None | 13.935 | 1050.2 | 1901 | `results/tinyllama.json`[2] |
| `f5fp8urt` | `lfm120-instruct-official` | — | 3 | 4103 | `stop` | 244 | 78 | 2048 | None | 2.551 | 1098.0 | 1870 | `results/stageA.json`[2] |
| `f79n26h8` | `abl-lfm120` | C05 | 1 | 4101 | `stop` | 233 | 131 | 512 | None | 4.66 | 1148.8 | 1917 | `results/ablation_corpus.json`[8] |
| `f7i6ucx9` | `corpus-4b-nonthink` | C07 | 2 | 4102 | `stop` | 226 | 68 | 768 | None | 9.309 | 3556.4 | 1935 | `results/corpus.json`[19] |
| `fa6xrslj` | `corpus-lfm120-instruct` | C04 | 3 | 4103 | `stop` | 232 | 97 | 512 | None | 3.166 | 1103.9 | 1891 | `results/corpus.json`[65] |
| `fdxemh4v` | `abl-2b-nonthink` | — | 3 | 4103 | `stop` | 230 | 123 | 768 | None | 7.349 | 1807.7 | 1894 | `results/ablation_triage.json`[8] |
| `fh4qbfv2` | `corpus-lfm120-t08` | C09 | 1 | 4101 | `stop` | 225 | 26 | 512 | None | 1.038 | 1103.3 | 2027 | `results/corpus.json`[178] |
| `fjqmajpi` | `qwen35-2b-nonthink` | C08 | 2 | 4102 | `stop` | 196 | 79 | 2048 | None | 4.738 | 1959.7 | 1930 | `results/control_corpus_2b.json`[15] |
| `fkl1qe4g` | `corpus-4b-nonthink` | C04 | 3 | 4103 | `stop` | 228 | 69 | 768 | None | 9.451 | 3556.4 | 1935 | `results/corpus.json`[11] |
| `flyb1ycp` | `corpus-tinyllama` | C04 | 2 | 4102 | `stop` | 271 | 141 | 768 | None | 4.159 | 1135.2 | 1901 | `results/corpus.json`[206] |
| `fmgwwq04` | `abl-lfm120` | C10 | 2 | 4102 | `stop` | 470 | 30 | 512 | None | 1.031 | 1148.9 | 1919 | `results/ablation_corpus.json`[19] |
| `fo96ue2z` | `corpus-4b-nonthink` | C10 | 2 | 4102 | `stop` | 462 | 124 | 768 | None | 17.199 | 3556.4 | 1935 | `results/corpus.json`[28] |
| `fqu0r73o` | `corpus-lfm3-think-512` | C16 | 1 | 4101 | `stop` | 195 | 586 | 1024 | 512 | 28.925 | 2115.6 | 2080 | `results/corpus.json`[138] |
| `frmuuwk5` | `corpus-lfm120-instruct` | C12 | 1 | 4101 | `stop` | 241 | 49 | 512 | None | 1.819 | 1111.3 | 2117 | `results/corpus.json`[87] |
| `fsots5p7` | `qwen3-17b-think-rb1024` | — | 3 | 4103 | `stop` | 227 | 719 | 3072 | 1024 | 37.907 | 2112.5 | 1994 | `results/stageA2a.json`[8] |
| `g3cfm9hp` | `corpus-lfm120-instruct` | C08 | 2 | 4102 | `stop` | 201 | 44 | 512 | None | 1.482 | 1104.0 | 2117 | `results/corpus.json`[76] |
| `g4ms5gk9` | `abl-2b-nonthink` | C01 | 2 | 4102 | `stop` | 232 | 37 | 768 | None | 2.224 | 1807.5 | 1926 | `results/ablation_corpus_abl-2b-nonthink.json`[1] |
| `g4wiytnc` | `abl-2b-nonthink` | C14 | 2 | 4102 | `stop` | 256 | 55 | 768 | None | 3.323 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[27] |
| `g98y7v7b` | `abl-2b-nonthink` | C07 | 1 | 4101 | `stop` | 226 | 96 | 768 | None | 6.016 | 1942.8 | 1929 | `results/ablation_corpus_abl-2b-nonthink.json`[12] |
| `gbkapfjb` | `qwen35-2b-nonthink` | C18 | 2 | 4102 | `stop` | 212 | 48 | 2048 | None | 2.907 | 2187.6 | 1930 | `results/control_corpus_2b.json`[35] |
| `gfbwiuwi` | `corpus-lfm120-instruct` | C16 | 1 | 4101 | `stop` | 198 | 54 | 512 | None | 1.923 | 1111.3 | 2117 | `results/corpus.json`[99] |
| `ghoae6xk` | `abl-2b-nonthink` | C17 | 1 | 4101 | `stop` | 239 | 69 | 768 | None | 4.397 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[32] |
| `gi5uc8z6` | `abl-2b-nonthink` | C02 | 1 | 4101 | `stop` | 300 | 73 | 768 | None | 4.821 | 1807.5 | 1926 | `results/ablation_corpus_abl-2b-nonthink.json`[2] |
| `gktos44v` | `qwen35-2b-nonthink` | C17 | 1 | 4101 | `stop` | 239 | 50 | 2048 | None | 3.398 | 2142.8 | 1930 | `results/control_corpus_2b.json`[32] |
| `gp26oj0o` | `qwen35-2b-think-rb1024` | — | 2 | 4102 | `stop` | 228 | 1137 | 3072 | 1024 | 68.273 | 1879.0 | 1931 | `results/stageA2a.json`[19] |
| `gtdf0ftp` | `corpus-lfm120-t08` | C18 | 2 | 4102 | `stop` | 215 | 49 | 512 | None | 1.625 | 1110.0 | 2029 | `results/corpus.json`[197] |
| `gwku9k8h` | `corpus-4b-nonthink` | C05 | 1 | 4101 | `stop` | 232 | 53 | 768 | None | 7.912 | 3556.4 | 1935 | `results/corpus.json`[12] |
| `h3247ejy` | `corpus-4b-think-512` | C08 | 1 | 4101 | `stop` | 194 | 564 | 1024 | 512 | 79.383 | 4363.4 | 2193 | `results/corpus.json`[151] |
| `h44826yb` | `corpus-lfm120-t08` | C12 | 2 | 4102 | `stop` | 241 | 55 | 512 | None | 1.839 | 1110.0 | 2029 | `results/corpus.json`[185] |
| `h72bu67o` | `corpus-4b-nonthink` | C14 | 2 | 4102 | `stop` | 256 | 53 | 768 | None | 7.467 | 3681.7 | 1935 | `results/corpus.json`[40] |
| `h7efkw2d` | `corpus-lfm120-instruct` | C18 | 1 | 4101 | `stop` | 215 | 49 | 512 | None | 1.768 | 1111.3 | 2117 | `results/corpus.json`[105] |
| `hb10z3tx` | `corpus-lfm3-think-512` | C18 | 2 | 4102 | `stop` | 212 | 613 | 1024 | 512 | 29.87 | 2199.5 | 2135 | `results/corpus.json`[143] |
| `heqk20kd` | `corpus-tinyllama` | C17 | 1 | 4101 | `stop` | 288 | 219 | 768 | None | 6.787 | 1543.0 | 2076 | `results/corpus.json`[231] |
| `hih53m75` | `corpus-tinyllama` | C01 | 1 | 4101 | `length` | 278 | 768 | 768 | None | 23.928 | 1019.0 | 1901 | `results/corpus.json`[199] |
| `hjrhoayt` | `corpus-lfm120-instruct` | C15 | 1 | 4101 | `stop` | 256 | 53 | 512 | None | 1.92 | 1111.3 | 2117 | `results/corpus.json`[96] |
| `hkzlgvub` | `abl-4b-think-512` | — | 3 | 4103 | `stop` | 228 | 613 | 1024 | 512 | 86.279 | 3666.4 | 1940 | `results/ablation_triage.json`[5] |
| `hmuu98vb` | `abl-2b-nonthink` | C04 | 1 | 4101 | `stop` | 228 | 91 | 768 | None | 5.719 | 1898.5 | 1929 | `results/ablation_corpus_abl-2b-nonthink.json`[6] |
| `hnw1x1og` | `corpus-lfm120-instruct` | C07 | 1 | 4101 | `stop` | 230 | 41 | 512 | None | 1.519 | 1104.0 | 2117 | `results/corpus.json`[72] |
| `hnwenpqt` | `lfm3-17b-think-rb1024` | — | 1 | 4101 | `stop` | 238 | 1130 | 3072 | 1024 | 55.603 | 1559.1 | 1918 | `results/stageA2a.json`[12] |
| `hpp51phj` | `abl-2b-nonthink` | C05 | 1 | 4101 | `stop` | 232 | 106 | 768 | None | 6.597 | 1898.5 | 1929 | `results/ablation_corpus_abl-2b-nonthink.json`[8] |
| `hrckg6hd` | `corpus-4b-nonthink` | C01 | 3 | 4103 | `stop` | 232 | 57 | 768 | None | 8.033 | 3437.7 | 1929 | `results/corpus.json`[2] |
| `ht15ha3v` | `qwen35-2b-nonthink` | C14 | 2 | 4102 | `stop` | 256 | 104 | 2048 | None | 6.255 | 2098.0 | 1930 | `results/control_corpus_2b.json`[27] |
| `hw5ojioj` | `corpus-lfm120-instruct` | C03 | 3 | 4103 | `stop` | 214 | 65 | 512 | None | 2.141 | 1103.9 | 1891 | `results/corpus.json`[62] |
| `i01u65j3` | `tinyllama-11b-chat` | — | 2 | 4102 | `stop` | 268 | 519 | 768 | None | 15.675 | 1032.2 | 1901 | `results/tinyllama.json`[1] |
| `i0cwcfbt` | `abl-4b-nonthink` | C11 | 1 | 4101 | `stop` | 241 | 55 | 768 | None | 8.471 | 3760.9 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[20] |
| `i2k0v8s0` | `corpus-4b-think-512` | C07 | 1 | 4101 | `stop` | 224 | 579 | 1024 | 512 | 81.693 | 4231.7 | 2193 | `results/corpus.json`[150] |
| `i59ithrx` | `corpus-tinyllama` | C05 | 2 | 4102 | `length` | 270 | 768 | 768 | None | 23.644 | 1160.4 | 1901 | `results/corpus.json`[208] |
| `i9cqrswv` | `qwen35-2b-nonthink` | C12 | 2 | 4102 | `stop` | 236 | 63 | 2048 | None | 3.796 | 2053.0 | 1930 | `results/control_corpus_2b.json`[23] |
| `ifc11h5r` | `corpus-lfm120-t08` | C06 | 2 | 4102 | `stop` | 238 | 52 | 512 | None | 1.714 | 1103.3 | 2027 | `results/corpus.json`[173] |
| `iho1g9oz` | `qwen35-4b-nonthink` | — | 1 | 4101 | `stop` | 230 | 90 | 2048 | None | 13.071 | 3437.7 | 1913 | `results/stageA2a.json`[0] |
| `ij22jeoc` | `corpus-lfm120-instruct` | C12 | 2 | 4102 | `stop` | 241 | 51 | 512 | None | 1.666 | 1111.3 | 2117 | `results/corpus.json`[88] |
| `ijocqvql` | `corpus-lfm120-instruct` | C06 | 1 | 4101 | `stop` | 238 | 43 | 512 | None | 1.75 | 1104.0 | 2117 | `results/corpus.json`[69] |
| `il3bbxtt` | `abl-2b-nonthink` | C12 | 2 | 4102 | `stop` | 236 | 67 | 768 | None | 4.014 | 2037.4 | 1930 | `results/ablation_corpus_abl-2b-nonthink.json`[23] |
| `ipg1mn60` | `corpus-lfm3-think-512` | C06 | 2 | 4102 | `stop` | 233 | 598 | 1024 | 512 | 29.453 | 1757.8 | 2040 | `results/corpus.json`[119] |
| `ipsh93qy` | `qwen35-08b-nonthink` | — | 3 | 4103 | `stop` | 230 | 120 | 2048 | None | 3.441 | 1068.8 | 1877 | `results/stageA.json`[11] |
| `iss7y0d2` | `qwen35-2b-nonthink` | C11 | 2 | 4102 | `stop` | 241 | 99 | 2048 | None | 6.028 | 2007.9 | 1930 | `results/control_corpus_2b.json`[21] |
| `iws648vk` | `corpus-lfm3-think-512` | C13 | 1 | 4101 | `length` | 255 | 1024 | 1024 | 512 | 50.994 | 2000.0 | 2080 | `results/corpus.json`[132] |
| `j06olf6e` | `corpus-4b-nonthink` | C12 | 2 | 4102 | `stop` | 236 | 156 | 768 | None | 21.917 | 3681.5 | 1935 | `results/corpus.json`[34] |
| `j4x0yqp5` | `corpus-lfm3-think-512` | C11 | 1 | 4101 | `stop` | 242 | 588 | 1024 | 512 | 29.251 | 1925.2 | 2080 | `results/corpus.json`[128] |
| `j7378njc` | `corpus-lfm3-think-512` | C13 | 2 | 4102 | `stop` | 255 | 598 | 1024 | 512 | 29.763 | 2022.8 | 2080 | `results/corpus.json`[133] |
| `ja69ru5g` | `corpus-lfm120-instruct` | C05 | 1 | 4101 | `stop` | 233 | 47 | 512 | None | 1.727 | 1103.9 | 1891 | `results/corpus.json`[66] |
| `jbkngh4s` | `corpus-tinyllama` | C13 | 1 | 4101 | `length` | 300 | 768 | 768 | None | 23.773 | 1398.4 | 2076 | `results/corpus.json`[223] |
| `jc7tq2qb` | `corpus-4b-nonthink` | C01 | 2 | 4102 | `stop` | 232 | 52 | 768 | None | 7.266 | 3437.7 | 1929 | `results/corpus.json`[1] |
| `jdwwceq4` | `abl-4b-nonthink` | C16 | 1 | 4101 | `stop` | 192 | 50 | 768 | None | 7.937 | 4261.0 | 2027 | `results/ablation_corpus_abl-4b-nonthink.json`[30] |
| `jex9gov2` | `corpus-tinyllama` | C02 | 1 | 4101 | `length` | 357 | 768 | 768 | None | 24.402 | 1056.6 | 1901 | `results/corpus.json`[201] |
| `jjf67xzt` | `corpus-4b-nonthink` | C18 | 2 | 4102 | `stop` | 212 | 76 | 768 | None | 10.38 | 3681.7 | 1935 | `results/corpus.json`[52] |
| `jjt5amq4` | `corpus-lfm120-t08` | C17 | 2 | 4102 | `stop` | 248 | 35 | 512 | None | 1.18 | 1110.0 | 2029 | `results/corpus.json`[195] |
| `jlvi1cs4` | `abl-lfm120` | C16 | 1 | 4101 | `stop` | 198 | 54 | 512 | None | 2.031 | 1165.7 | 1919 | `results/ablation_corpus.json`[30] |
| `jlvnuykn` | `abl-lfm120` | C15 | 2 | 4102 | `stop` | 256 | 70 | 512 | None | 2.376 | 1165.7 | 1919 | `results/ablation_corpus.json`[29] |
| `jnjxkuvx` | `abl-4b-nonthink` | C05 | 2 | 4102 | `stop` | 232 | 51 | 768 | None | 7.043 | 3636.5 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[9] |
| `jo21wmi5` | `corpus-4b-think-512` | C17 | 1 | 4101 | `stop` | 237 | 581 | 1024 | 512 | 88.285 | 5555.9 | 2193 | `results/corpus.json`[160] |
| `jput80yt` | `corpus-4b-think-512` | C01 | 1 | 4101 | `stop` | 230 | 599 | 1024 | 512 | 86.114 | 3437.7 | 2119 | `results/corpus.json`[144] |
| `jql39zja` | `abl-4b-nonthink` | C18 | 1 | 4101 | `stop` | 212 | 52 | 768 | None | 7.66 | 4261.0 | 2027 | `results/ablation_corpus_abl-4b-nonthink.json`[34] |
| `jr3wb6d0` | `qwen35-2b-nonthink-pres0` | — | 1 | 4101 | `stop` | 230 | 93 | 2048 | None | 6.166 | 1820.9 | 2027 | `results/stageB.json`[12] |
| `jrij0wqw` | `abl-lfm120` | C09 | 2 | 4102 | `stop` | 225 | 43 | 512 | None | 1.464 | 1148.9 | 1917 | `results/ablation_corpus.json`[17] |
| `jrl6f05f` | `corpus-lfm120-instruct` | C01 | 3 | 4103 | `stop` | 236 | 40 | 512 | None | 1.315 | 1098.0 | 1891 | `results/corpus.json`[56] |
| `jwq9efbf` | `abl-lfm120` | C12 | 1 | 4101 | `stop` | 241 | 35 | 512 | None | 1.575 | 1160.6 | 1919 | `results/ablation_corpus.json`[22] |
| `jy7w18t5` | `corpus-tinyllama` | C18 | 2 | 4102 | `stop` | 254 | 292 | 768 | None | 8.695 | 1590.0 | 2076 | `results/corpus.json`[234] |
| `k0c5vzve` | `abl-4b-nonthink` | C09 | 2 | 4102 | `stop` | 218 | 2 | 768 | None | 0.335 | 3636.6 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[17] |
| `k2f51lf8` | `corpus-lfm120-instruct` | C13 | 3 | 4103 | `stop` | 258 | 64 | 512 | None | 2.132 | 1111.3 | 2117 | `results/corpus.json`[92] |
| `k2wujq4c` | `corpus-4b-nonthink` | C18 | 1 | 4101 | `stop` | 212 | 53 | 768 | None | 7.962 | 3681.7 | 1935 | `results/corpus.json`[51] |
| `k4bewsj2` | `qwen35-4b-nonthink-gpufull` | — | 1 | 4101 | `stop` | 230 | 99 | 768 | None | 2.556 | 3265.5 | 4833 | `results/gpu-full-qwen35-4b-nonthink-gpufull.json`[0] |
| `k9wx79t9` | `abl-lfm120` | C16 | 2 | 4102 | `stop` | 198 | 104 | 512 | None | 3.538 | 1165.7 | 1919 | `results/ablation_corpus.json`[31] |
| `kbyukf95` | `corpus-lfm3-think-512` | C14 | 2 | 4102 | `stop` | 258 | 617 | 1024 | 512 | 30.718 | 2063.5 | 2080 | `results/corpus.json`[135] |
| `kg1kmeh8` | `qwen35-2b-nonthink` | C10 | 2 | 4102 | `stop` | 462 | 172 | 2048 | None | 10.365 | 1959.7 | 1930 | `results/control_corpus_2b.json`[19] |
| `kl458eev` | `corpus-4b-nonthink` | C03 | 1 | 4101 | `stop` | 207 | 66 | 768 | None | 9.677 | 3556.4 | 1933 | `results/corpus.json`[6] |
| `ktklchvc` | `abl-lfm120` | C06 | 2 | 4102 | `stop` | 238 | 21 | 512 | None | 0.736 | 1148.9 | 1917 | `results/ablation_corpus.json`[11] |
| `kug7gup3` | `corpus-lfm3-think-512` | C03 | 1 | 4101 | `stop` | 211 | 595 | 1024 | 512 | 29.174 | 1632.2 | 2040 | `results/corpus.json`[112] |
| `kzqpgvun` | `qwen35-4b-think-rb256` | — | 1 | 4101 | `stop` | 228 | 363 | 3072 | 256 | 50.654 | 3437.7 | 1964 | `results/all_triage.json`[47] |
| `l5jq1mb7` | `corpus-tinyllama` | C09 | 1 | 4101 | `length` | 259 | 768 | 768 | None | 24.261 | 1286.6 | 2076 | `results/corpus.json`[215] |
| `l8cdqh92` | `corpus-4b-nonthink` | C16 | 3 | 4103 | `stop` | 192 | 52 | 768 | None | 7.159 | 3681.7 | 1935 | `results/corpus.json`[47] |
| `lju6j28m` | `corpus-4b-nonthink` | C12 | 1 | 4101 | `stop` | 236 | 99 | 768 | None | 14.339 | 3681.5 | 1935 | `results/corpus.json`[33] |
| `lol858ws` | `corpus-lfm120-instruct` | C06 | 2 | 4102 | `stop` | 238 | 43 | 512 | None | 1.531 | 1104.0 | 2117 | `results/corpus.json`[70] |
| `lowm1zik` | `corpus-lfm3-think-512` | C05 | 1 | 4101 | `stop` | 230 | 591 | 1024 | 512 | 28.997 | 1705.9 | 2040 | `results/corpus.json`[116] |
| `lsm32mct` | `qwen3-17b-think-rb1024` | — | 1 | 4101 | `stop` | 227 | 396 | 3072 | 1024 | 20.427 | 1966.4 | 1877 | `results/stageA2a.json`[6] |
| `lt7jeu9e` | `corpus-tinyllama` | C14 | 2 | 4102 | `stop` | 300 | 372 | 768 | None | 11.239 | 1460.5 | 2076 | `results/corpus.json`[226] |
| `luqe6lqe` | `corpus-lfm120-instruct` | C13 | 1 | 4101 | `stop` | 258 | 63 | 512 | None | 2.249 | 1111.3 | 2117 | `results/corpus.json`[90] |
| `lvdo9j11` | `abl-2b-nonthink` | C01 | 1 | 4101 | `stop` | 232 | 198 | 768 | None | 12.205 | 1807.3 | 1926 | `results/ablation_corpus_abl-2b-nonthink.json`[0] |
| `m04nahm9` | `lfm120-think-official-rb1024` | — | 2 | 4102 | `stop` | 244 | 1118 | 3072 | 1024 | 39.777 | 1176.0 | 1925 | `results/stageA2a.json`[22] |
| `m1cx761r` | `abl-2b-nonthink` | C07 | 2 | 4102 | `stop` | 226 | 217 | 768 | None | 12.936 | 1942.8 | 1929 | `results/ablation_corpus_abl-2b-nonthink.json`[13] |
| `m1qgbanq` | `abl-lfm120` | C04 | 1 | 4101 | `stop` | 232 | 134 | 512 | None | 4.74 | 1148.8 | 1917 | `results/ablation_corpus.json`[6] |
| `m5abmcdd` | `qwen35-08b-nonthink` | — | 1 | 4101 | `stop` | 230 | 134 | 2048 | None | 4.091 | 1068.6 | 1877 | `results/stageA.json`[9] |
| `maft6z9n` | `qwen35-2b-nonthink-pres0` | — | 3 | 4103 | `stop` | 230 | 80 | 2048 | None | 4.966 | 1821.1 | 2027 | `results/stageB.json`[14] |
| `mebd8kgd` | `corpus-4b-think-512` | C03 | 1 | 4101 | `stop` | 205 | 584 | 1024 | 512 | 83.126 | 3704.9 | 2123 | `results/corpus.json`[146] |
| `mg9ai5c9` | `abl-lfm120` | C18 | 1 | 4101 | `stop` | 215 | 93 | 512 | None | 3.372 | 1165.7 | 1919 | `results/ablation_corpus.json`[34] |
| `mgpxa7zp` | `corpus-lfm3-think-512` | C09 | 2 | 4102 | `stop` | 222 | 517 | 1024 | 512 | 25.795 | 1857.7 | 2080 | `results/corpus.json`[125] |
| `moyk7h7z` | `corpus-4b-nonthink` | C04 | 2 | 4102 | `stop` | 228 | 61 | 768 | None | 8.367 | 3556.4 | 1935 | `results/corpus.json`[10] |
| `mu08h6b0` | `corpus-4b-think-512` | C11 | 1 | 4101 | `stop` | 239 | 591 | 1024 | 512 | 83.917 | 4764.3 | 2193 | `results/corpus.json`[154] |
| `muheoyey` | `corpus-tinyllama` | C15 | 1 | 4101 | `stop` | 295 | 258 | 768 | None | 7.854 | 1475.8 | 2076 | `results/corpus.json`[227] |
| `muymlpe9` | `corpus-tinyllama` | C09 | 2 | 4102 | `stop` | 259 | 628 | 768 | None | 19.144 | 1310.1 | 2076 | `results/corpus.json`[216] |
| `n20tc159` | `qwen35-2b-nonthink-r2cap` | — | 1 | 4101 | `stop` | 230 | 108 | 400 | None | 6.695 | 1794.1 | 2035 | `results/stageC.json`[0] |
| `nak777ef` | `corpus-tinyllama` | C01 | 2 | 4102 | `stop` | 278 | 295 | 768 | None | 8.802 | 1043.2 | 1901 | `results/corpus.json`[200] |
| `nbpyv0nw` | `lfm120-think-official-rb1024` | — | 3 | 4103 | `stop` | 244 | 2274 | 3072 | 1024 | 83.432 | 1193.6 | 2019 | `results/stageA2a.json`[23] |
| `nds7hqaw` | `corpus-lfm3-think-512` | C08 | 1 | 4101 | `stop` | 198 | 559 | 1024 | 512 | 27.444 | 1809.2 | 2040 | `results/corpus.json`[122] |
| `nglxagzr` | `abl-4b-nonthink` | C01 | 2 | 4102 | `stop` | 232 | 57 | 768 | None | 7.839 | 3399.9 | 1959 | `results/ablation_corpus_abl-4b-nonthink.json`[1] |
| `ngrqugjm` | `corpus-lfm3-think-512` | C08 | 2 | 4102 | `stop` | 198 | 553 | 1024 | 512 | 27.052 | 1825.1 | 2040 | `results/corpus.json`[123] |
| `nj7firx7` | `qwen35-2b-nonthink-r2cap` | — | 3 | 4103 | `stop` | 230 | 100 | 400 | None | 5.898 | 1794.2 | 2035 | `results/stageC.json`[2] |
| `njo3oggi` | `corpus-lfm120-instruct` | C10 | 1 | 4101 | `stop` | 470 | 113 | 512 | None | 3.983 | 1104.0 | 2117 | `results/corpus.json`[81] |
| `nls49335` | `abl-lfm120` | C07 | 1 | 4101 | `stop` | 230 | 140 | 512 | None | 4.935 | 1148.9 | 1917 | `results/ablation_corpus.json`[12] |
| `nuo0753y` | `qwen35-4b-nonthink` | — | 3 | 4103 | `stop` | 230 | 87 | 2048 | None | 11.884 | 3437.7 | 1913 | `results/stageA2a.json`[2] |
| `nxhqtfvv` | `corpus-lfm120-t08` | C10 | 2 | 4102 | `stop` | 470 | 57 | 512 | None | 1.89 | 1103.3 | 2029 | `results/corpus.json`[181] |
| `nxvtp78k` | `corpus-lfm3-think-512` | C06 | 1 | 4101 | `stop` | 233 | 599 | 1024 | 512 | 29.332 | 1740.6 | 2040 | `results/corpus.json`[118] |
| `nzj9rd8f` | `lfm120-think-official-rb1024` | — | 1 | 4101 | `stop` | 244 | 2710 | 3072 | 1024 | 100.026 | 1139.2 | 1925 | `results/stageA2a.json`[21] |
| `o0xw1hlf` | `corpus-4b-think-512` | C10 | 1 | 4101 | `stop` | 460 | 621 | 1024 | 512 | 89.092 | 4623.2 | 2193 | `results/corpus.json`[153] |
| `o112u0gv` | `corpus-lfm120-instruct` | C07 | 2 | 4102 | `stop` | 230 | 41 | 512 | None | 1.426 | 1104.0 | 2117 | `results/corpus.json`[73] |
| `o3o5hn6j` | `corpus-4b-think-512` | C16 | 1 | 4101 | `stop` | 190 | 571 | 1024 | 512 | 86.295 | 5425.4 | 2193 | `results/corpus.json`[159] |
| `o4i38ub5` | `lfm3-17b-nonthink` | — | 1 | 4101 | `stop` | 239 | 108 | 3072 | None | 5.548 | 1557.9 | 1918 | `results/stageA2a.json`[15] |
| `o8klmnkl` | `qwen3-17b-think-rb1024` | — | 2 | 4102 | `stop` | 227 | 418 | 3072 | 1024 | 21.391 | 2038.1 | 1877 | `results/stageA2a.json`[7] |
| `oa02qie2` | `qwen35-2b-nonthink` | C15 | 1 | 4101 | `stop` | 249 | 95 | 2048 | None | 6.004 | 2142.8 | 1930 | `results/control_corpus_2b.json`[28] |
| `oa5tl97s` | `corpus-4b-nonthink` | C07 | 3 | 4103 | `stop` | 226 | 83 | 768 | None | 11.362 | 3556.4 | 1935 | `results/corpus.json`[20] |
| `oazcoaw8` | `corpus-lfm3-think-512` | C01 | 1 | 4101 | `stop` | 233 | 658 | 1024 | 512 | 32.193 | 1558.0 | 2038 | `results/corpus.json`[108] |
| `oevxxlvc` | `abl-lfm120` | C11 | 1 | 4101 | `stop` | 245 | 43 | 512 | None | 1.67 | 1155.4 | 1919 | `results/ablation_corpus.json`[20] |
| `oh2kp9sa` | `corpus-lfm3-think-512` | C07 | 2 | 4102 | `stop` | 227 | 655 | 1024 | 512 | 31.949 | 1791.7 | 2040 | `results/corpus.json`[121] |
| `ojp21co1` | `abl-4b-nonthink` | C12 | 1 | 4101 | `stop` | 236 | 85 | 768 | None | 12.354 | 3760.9 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[22] |
| `oqozsxx9` | `abl-4b-nonthink` | C12 | 2 | 4102 | `stop` | 236 | 80 | 768 | None | 11.021 | 3760.9 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[23] |
| `oru0lclb` | `corpus-lfm120-t08` | C02 | 1 | 4101 | `stop` | 311 | 64 | 512 | None | 2.295 | 1097.7 | 2027 | `results/corpus.json`[164] |
| `ouafvm4p` | `abl-4b-nonthink` | C10 | 1 | 4101 | `stop` | 462 | 110 | 768 | None | 16.087 | 3636.6 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[18] |
| `ov6pcv3m` | `abl-2b-nonthink` | C16 | 2 | 4102 | `stop` | 192 | 53 | 768 | None | 3.163 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[31] |
| `ovvbwplt` | `qwen35-2b-nonthink` | C18 | 1 | 4101 | `stop` | 212 | 73 | 2048 | None | 4.681 | 2187.6 | 1930 | `results/control_corpus_2b.json`[34] |
| `p21y2wjp` | `abl-2b-nonthink` | C06 | 1 | 4101 | `stop` | 233 | 85 | 768 | None | 5.462 | 1942.8 | 1929 | `results/ablation_corpus_abl-2b-nonthink.json`[10] |
| `p2divk2l` | `abl-4b-nonthink` | C05 | 1 | 4101 | `stop` | 232 | 50 | 768 | None | 7.778 | 3636.5 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[8] |
| `p5en8xxt` | `abl-4b-nonthink` | C04 | 2 | 4102 | `stop` | 228 | 96 | 768 | None | 13.305 | 3520.5 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[7] |
| `pa1h9xep` | `qwen35-2b-nonthink` | C03 | 1 | 4101 | `stop` | 207 | 79 | 2048 | None | 5.02 | 1866.7 | 1927 | `results/control_corpus_2b.json`[4] |
| `pgp54tpp` | `corpus-4b-nonthink` | C13 | 3 | 4103 | `stop` | 251 | 2 | 768 | None | 0.353 | 3681.7 | 1935 | `results/corpus.json`[38] |
| `piu6p69i` | `lfm120-ins-r2cap` | — | 1 | 4101 | `stop` | 244 | 69 | 400 | None | 2.487 | 1070.6 | 2027 | `results/stageC.json`[3] |
| `pova8bpn` | `corpus-4b-nonthink` | C05 | 3 | 4103 | `stop` | 232 | 46 | 768 | None | 6.47 | 3556.4 | 1935 | `results/corpus.json`[14] |
| `pqzver45` | `abl-4b-nonthink` | C13 | 2 | 4102 | `stop` | 251 | 59 | 768 | None | 8.972 | 3876.6 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[25] |
| `q019y7xm` | `corpus-lfm120-t08` | C13 | 2 | 4102 | `stop` | 258 | 49 | 512 | None | 1.635 | 1110.0 | 2029 | `results/corpus.json`[187] |
| `q0v9bbyb` | `corpus-lfm3-think-512` | C17 | 2 | 4102 | `stop` | 245 | 577 | 1024 | 512 | 28.961 | 2166.0 | 2135 | `results/corpus.json`[141] |
| `q1ue40uk` | `qwen35-9b-think-rb1024` | — | 1 | 4101 | `length` | 228 | 3072 | 3072 | 1024 | 758.263 | 5817.0 | 2085 | `results/stageA2b.json`[0] |
| `q575o5l9` | `corpus-4b-think-512` | C02 | 1 | 4101 | `stop` | 298 | 586 | 1024 | 512 | 83.452 | 3570.5 | 2119 | `results/corpus.json`[145] |
| `q5hb3tdm` | `abl-4b-nonthink` | C14 | 1 | 4101 | `stop` | 256 | 440 | 768 | None | 63.464 | 3876.8 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[26] |
| `q7ykffw0` | `abl-lfm120` | C10 | 1 | 4101 | `stop` | 470 | 79 | 512 | None | 2.973 | 1148.9 | 1919 | `results/ablation_corpus.json`[18] |
| `qahlyez3` | `corpus-4b-nonthink` | C08 | 3 | 4103 | `stop` | 196 | 47 | 768 | None | 6.431 | 3556.4 | 1935 | `results/corpus.json`[23] |
| `qe2f5tkm` | `abl-2b-nonthink` | C15 | 1 | 4101 | `stop` | 249 | 63 | 768 | None | 4.064 | 2037.4 | 1931 | `results/ablation_corpus_abl-2b-nonthink.json`[28] |
| `qehp1jfg` | `lfm350-official` | — | 3 | 4103 | `stop` | 244 | 16 | 2048 | None | 0.19 | 610.9 | 1836 | `results/stageA.json`[5] |
| `qf3azlw0` | `corpus-4b-nonthink` | C16 | 1 | 4101 | `stop` | 192 | 60 | 768 | None | 8.78 | 3681.7 | 1935 | `results/corpus.json`[45] |
| `qjmhxi9k` | `corpus-tinyllama` | C02 | 2 | 4102 | `length` | 357 | 768 | 768 | None | 23.928 | 1082.2 | 1901 | `results/corpus.json`[202] |
| `ql5zrcmb` | `qwen35-2b-nonthink` | C11 | 1 | 4101 | `stop` | 241 | 96 | 2048 | None | 6.07 | 2007.9 | 1930 | `results/control_corpus_2b.json`[20] |
| `qx8hu3am` | `qwen35-2b-nonthink` | C01 | 1 | 4101 | `stop` | 232 | 100 | 2048 | None | 6.342 | 1820.9 | 1925 | `results/control_corpus_2b.json`[0] |
| `qzgwlj6w` | `lfm120-think-official` | — | 1 | 4101 | `stop` | 244 | 2453 | 3072 | None | 91.051 | 1139.2 | 1956 | `results/stageA.json`[6] |
| `r07g101o` | `qwen35-2b-nonthink` | C01 | 2 | 4102 | `stop` | 232 | 64 | 2048 | None | 3.874 | 1821.0 | 1925 | `results/control_corpus_2b.json`[1] |
| `r1o6r966` | `abl-4b-nonthink` | C13 | 1 | 4101 | `stop` | 251 | 41 | 768 | None | 6.606 | 3876.6 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[24] |
| `r3xzwh3j` | `abl-4b-nonthink` | C17 | 1 | 4101 | `stop` | 239 | 92 | 768 | None | 13.189 | 4261.0 | 2027 | `results/ablation_corpus_abl-4b-nonthink.json`[32] |
| `r75w0fsb` | `qwen35-2b-nonthink` | — | 1 | 4101 | `stop` | 230 | 136 | 2048 | None | 8.476 | 1820.9 | 1895 | `results/stageA.json`[12] |
| `readpkz1` | `corpus-lfm120-instruct` | C17 | 2 | 4102 | `stop` | 248 | 65 | 512 | None | 2.154 | 1111.3 | 2117 | `results/corpus.json`[103] |
| `rfjqowmq` | `abl-lfm120` | C08 | 1 | 4101 | `stop` | 201 | 3 | 512 | None | 0.292 | 1148.9 | 1917 | `results/ablation_corpus.json`[14] |
| `ribqbcco` | `qwen35-2b-nonthink` | C04 | 1 | 4101 | `stop` | 228 | 46 | 2048 | None | 3.042 | 1866.7 | 1929 | `results/control_corpus_2b.json`[6] |
| `rn4xe7ln` | `qwen3-17b-nonthink` | — | 2 | 4102 | `stop` | 231 | 106 | 2048 | None | 5.267 | 1965.7 | 1925 | `results/stageA2a.json`[10] |
| `rn821lbj` | `corpus-tinyllama` | C12 | 1 | 4101 | `stop` | 276 | 272 | 768 | None | 8.252 | 1387.2 | 2076 | `results/corpus.json`[221] |
| `rodx7mvy` | `lfm120-ins-r2cap` | — | 2 | 4102 | `stop` | 244 | 73 | 400 | None | 2.42 | 1070.7 | 2027 | `results/stageC.json`[4] |
| `rq0xspl5` | `abl-2b-nonthink` | C10 | 2 | 4102 | `stop` | 462 | 147 | 768 | None | 8.796 | 1989.5 | 1930 | `results/ablation_corpus_abl-2b-nonthink.json`[19] |
| `rxeowzey` | `corpus-4b-nonthink` | C07 | 1 | 4101 | `stop` | 226 | 43 | 768 | None | 6.531 | 3556.4 | 1935 | `results/corpus.json`[18] |
| `rxt9ji93` | `corpus-tinyllama` | C10 | 1 | 4101 | `stop` | 535 | 445 | 768 | None | 14.38 | 1330.6 | 2076 | `results/corpus.json`[217] |
| `s50bt6ac` | `abl-4b-think-512` | — | 1 | 4101 | `stop` | 228 | 618 | 1024 | 512 | 86.883 | 3400.4 | 1940 | `results/ablation_triage.json`[3] |
| `s5ffsmb5` | `corpus-tinyllama` | C08 | 1 | 4101 | `length` | 235 | 768 | 768 | None | 23.72 | 1242.0 | 1901 | `results/corpus.json`[213] |
| `s7zlwy1x` | `corpus-tinyllama` | C04 | 1 | 4101 | `stop` | 271 | 103 | 768 | None | 3.205 | 1135.1 | 1901 | `results/corpus.json`[205] |
| `s9nzf6j1` | `abl-2b-nonthink` | C03 | 2 | 4102 | `stop` | 207 | 162 | 768 | None | 9.66 | 1852.7 | 1928 | `results/ablation_corpus_abl-2b-nonthink.json`[5] |
| `sd2biwzn` | `qwen35-4b-think-rb1024` | — | 1 | 4101 | `stop` | 228 | 1128 | 3072 | 1024 | 160.483 | 3438.0 | 1919 | `results/stageA2a.json`[3] |
| `sdgh63sv` | `corpus-4b-nonthink` | C08 | 2 | 4102 | `stop` | 196 | 37 | 768 | None | 5.048 | 3556.4 | 1935 | `results/corpus.json`[22] |
| `sdxbfnx3` | `corpus-tinyllama` | C16 | 1 | 4101 | `length` | 230 | 768 | 768 | None | 23.502 | 1497.7 | 2076 | `results/corpus.json`[229] |
| `sekax9ya` | `corpus-tinyllama` | C06 | 2 | 4102 | `stop` | 273 | 347 | 768 | None | 10.422 | 1207.7 | 1901 | `results/corpus.json`[210] |
| `smr5b5lo` | `corpus-lfm120-instruct` | C15 | 3 | 4103 | `stop` | 256 | 56 | 512 | None | 1.834 | 1111.3 | 2117 | `results/corpus.json`[98] |
| `so8fot6n` | `abl-lfm120` | C13 | 1 | 4101 | `stop` | 258 | 159 | 512 | None | 5.622 | 1160.6 | 1919 | `results/ablation_corpus.json`[24] |
| `ssjxrqx6` | `abl-4b-nonthink` | C08 | 1 | 4101 | `stop` | 196 | 56 | 768 | None | 8.396 | 3636.6 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[14] |
| `sw9hssvq` | `qwen35-2b-nonthink` | C15 | 2 | 4102 | `stop` | 249 | 66 | 2048 | None | 3.984 | 2142.8 | 1930 | `results/control_corpus_2b.json`[29] |
| `t1hggto0` | `corpus-lfm3-nonthink` | C01 | 1 | 4101 | `length` | 233 | 512 | 512 | None | 25.081 | 1557.8 | 2021 | `results/corpus.json`[198] |
| `t5qe1gww` | `corpus-lfm120-instruct` | C18 | 3 | 4103 | `stop` | 215 | 49 | 512 | None | 1.623 | 1111.3 | 2117 | `results/corpus.json`[107] |
| `t9s0ve09` | `corpus-lfm120-instruct` | C10 | 2 | 4102 | `stop` | 470 | 103 | 512 | None | 3.453 | 1104.0 | 2117 | `results/corpus.json`[82] |
| `taqhpkpn` | `corpus-lfm120-t08` | C08 | 2 | 4102 | `stop` | 201 | 54 | 512 | None | 1.781 | 1103.3 | 2027 | `results/corpus.json`[177] |
| `te14qihn` | `corpus-tinyllama` | C14 | 1 | 4101 | `stop` | 300 | 313 | 768 | None | 9.543 | 1446.6 | 2076 | `results/corpus.json`[225] |
| `tg5uc595` | `qwen35-2b-nonthink` | C02 | 2 | 4102 | `stop` | 300 | 53 | 2048 | None | 3.216 | 1821.9 | 1925 | `results/control_corpus_2b.json`[3] |
| `tgk8076d` | `smollm3-3b-think-rb1024` | — | 1 | 4101 | `stop` | 467 | 471 | 3072 | 1024 | 41.7 | 2617.8 | 1964 | `results/stageA2b2.json`[0] |
| `tkt52gqc` | `lfm120-ins-t05` | — | 1 | 4101 | `stop` | 244 | 72 | 2048 | None | 2.786 | 1097.9 | 2025 | `results/stageB.json`[3] |
| `tmne7m5l` | `corpus-4b-think-512` | C18 | 1 | 4101 | `stop` | 210 | 567 | 1024 | 512 | 85.941 | 5688.1 | 2193 | `results/corpus.json`[161] |
| `tpuay55u` | `abl-2b-nonthink` | — | 1 | 4101 | `stop` | 230 | 101 | 768 | None | 6.481 | 1807.5 | 1894 | `results/ablation_triage.json`[6] |
| `tqykil1b` | `qwen35-2b-nonthink` | C16 | 1 | 4101 | `stop` | 192 | 59 | 2048 | None | 3.796 | 2142.8 | 1930 | `results/control_corpus_2b.json`[30] |
| `tyi4az0y` | `corpus-lfm3-think-512` | C12 | 2 | 4102 | `length` | 238 | 1024 | 1024 | 512 | 51.076 | 1977.3 | 2080 | `results/corpus.json`[131] |
| `u1evflgw` | `corpus-4b-think-512` | C09 | 1 | 4101 | `stop` | 216 | 516 | 1024 | 512 | 73.119 | 4493.7 | 2193 | `results/corpus.json`[152] |
| `u4203avw` | `corpus-lfm120-instruct` | C09 | 3 | 4103 | `stop` | 225 | 44 | 512 | None | 1.433 | 1104.0 | 2117 | `results/corpus.json`[80] |
| `uc7e1cey` | `abl-4b-nonthink` | C16 | 2 | 4102 | `stop` | 192 | 55 | 768 | None | 7.654 | 4261.0 | 2027 | `results/ablation_corpus_abl-4b-nonthink.json`[31] |
| `uh4hiwm3` | `qwen35-2b-nonthink-round2` | — | 2 | 4102 | `stop` | 230 | 143 | 2048 | None | 8.849 | 1821.4 | 2027 | `results/stageB.json`[16] |
| `uiu3s599` | `corpus-lfm3-think-512` | C10 | 2 | 4102 | `length` | 467 | 1024 | 1024 | 512 | 51.423 | 1899.9 | 2080 | `results/corpus.json`[127] |
| `uj6ghubc` | `abl-2b-nonthink` | C06 | 2 | 4102 | `stop` | 233 | 64 | 768 | None | 3.837 | 1942.8 | 1929 | `results/ablation_corpus_abl-2b-nonthink.json`[11] |
| `uszmntyz` | `qwen35-08b-nonthink` | — | 2 | 4102 | `stop` | 230 | 134 | 2048 | None | 3.865 | 1068.8 | 1877 | `results/stageA.json`[10] |
| `uv6q42gs` | `qwen35-2b-nonthink` | C05 | 1 | 4101 | `stop` | 232 | 159 | 2048 | None | 10.035 | 1912.9 | 1929 | `results/control_corpus_2b.json`[8] |
| `uxjxuguf` | `corpus-lfm120-t08` | C11 | 2 | 4102 | `stop` | 245 | 41 | 512 | None | 1.373 | 1109.9 | 2029 | `results/corpus.json`[183] |
| `uzo09tjl` | `corpus-lfm120-t08` | C05 | 2 | 4102 | `stop` | 233 | 45 | 512 | None | 1.487 | 1103.2 | 2027 | `results/corpus.json`[171] |
| `v1ztel52` | `abl-lfm120` | C03 | 2 | 4102 | `stop` | 214 | 111 | 512 | None | 3.768 | 1148.8 | 1917 | `results/ablation_corpus.json`[5] |
| `v9eda8nv` | `abl-lfm120` | C04 | 2 | 4102 | `stop` | 232 | 95 | 512 | None | 3.24 | 1148.8 | 1917 | `results/ablation_corpus.json`[7] |
| `vhs9701u` | `qwen35-4b-think-rb256` | — | 2 | 4102 | `stop` | 228 | 371 | 3072 | 256 | 52.993 | 3562.2 | 1964 | `results/all_triage.json`[48] |
| `vkoaxitr` | `corpus-4b-nonthink` | C17 | 1 | 4101 | `stop` | 239 | 2 | 768 | None | 1.172 | 3681.7 | 1935 | `results/corpus.json`[48] |
| `vl8gze0f` | `corpus-lfm120-instruct` | C08 | 1 | 4101 | `stop` | 201 | 44 | 512 | None | 1.616 | 1104.0 | 2117 | `results/corpus.json`[75] |
| `vller9bp` | `qwen35-4b-think-rb512` | — | 1 | 4101 | `stop` | 228 | 619 | 3072 | 512 | 88.751 | 3437.5 | 2033 | `results/all_triage.json`[49] |
| `vsj9dgnk` | `corpus-4b-think-512` | C06 | 1 | 4101 | `stop` | 231 | 576 | 1024 | 512 | 81.954 | 4099.8 | 2125 | `results/corpus.json`[149] |
| `vtl5jphd` | `lfm120-ins-t08-rep10` | — | 3 | 4103 | `stop` | 244 | 75 | 2048 | None | 2.624 | 1098.3 | 2025 | `results/stageB.json`[11] |
| `vyz7jhl9` | `abl-lfm120` | C14 | 2 | 4102 | `stop` | 261 | 68 | 512 | None | 2.33 | 1165.7 | 1919 | `results/ablation_corpus.json`[27] |
| `w4nf53gq` | `corpus-lfm120-instruct` | C14 | 2 | 4102 | `stop` | 261 | 50 | 512 | None | 1.633 | 1111.3 | 2117 | `results/corpus.json`[94] |
| `w7wmdv9r` | `corpus-lfm120-instruct` | C11 | 2 | 4102 | `stop` | 245 | 52 | 512 | None | 1.714 | 1111.3 | 2117 | `results/corpus.json`[85] |
| `w8n04jpv` | `corpus-4b-nonthink` | C09 | 3 | 4103 | `stop` | 218 | 2 | 768 | None | 0.325 | 3556.4 | 1935 | `results/corpus.json`[26] |
| `wcy57rrc` | `abl-2b-nonthink` | C09 | 1 | 4101 | `stop` | 218 | 91 | 768 | None | 5.708 | 1988.7 | 1930 | `results/ablation_corpus_abl-2b-nonthink.json`[16] |
| `we9aofzt` | `corpus-lfm120-t08` | C15 | 1 | 4101 | `stop` | 256 | 122 | 512 | None | 4.222 | 1110.0 | 2029 | `results/corpus.json`[190] |
| `weg5us35` | `qwen35-2b-nonthink` | C07 | 2 | 4102 | `stop` | 226 | 72 | 2048 | None | 4.328 | 1959.7 | 1930 | `results/control_corpus_2b.json`[13] |
| `wkp6s8fw` | `corpus-lfm120-instruct` | C15 | 2 | 4102 | `stop` | 256 | 53 | 512 | None | 1.743 | 1111.3 | 2117 | `results/corpus.json`[97] |
| `wne5ktq9` | `abl-lfm120` | C17 | 1 | 4101 | `stop` | 248 | 21 | 512 | None | 0.906 | 1165.7 | 1919 | `results/ablation_corpus.json`[32] |
| `wvu2586s` | `corpus-lfm3-think-512` | C09 | 1 | 4101 | `stop` | 222 | 576 | 1024 | 512 | 28.241 | 1841.1 | 2040 | `results/corpus.json`[124] |
| `wzx1pw1l` | `corpus-4b-nonthink` | C13 | 1 | 4101 | `stop` | 251 | 60 | 768 | None | 9.005 | 3681.7 | 1935 | `results/corpus.json`[36] |
| `x29ccll5` | `corpus-lfm3-think-512` | C15 | 1 | 4101 | `stop` | 253 | 603 | 1024 | 512 | 29.834 | 2081.4 | 2080 | `results/corpus.json`[136] |
| `x474s394` | `abl-4b-nonthink` | C17 | 2 | 4102 | `stop` | 239 | 42 | 768 | None | 5.713 | 4261.0 | 2027 | `results/ablation_corpus_abl-4b-nonthink.json`[33] |
| `x5qcg9r2` | `qwen3-17b-nonthink` | — | 3 | 4103 | `stop` | 231 | 128 | 2048 | None | 6.386 | 1965.8 | 1925 | `results/stageA2a.json`[11] |
| `x5vjtnbo` | `abl-2b-nonthink` | — | 2 | 4102 | `stop` | 230 | 103 | 768 | None | 6.202 | 1807.6 | 1894 | `results/ablation_triage.json`[7] |
| `x6dduhk4` | `smollm3-3b-nonthink` | — | 2 | 4102 | `stop` | 287 | 61 | 2048 | None | 5.223 | 2620.0 | 1964 | `results/stageA2b2.json`[4] |
| `x7rxe9fw` | `qwen35-2b-nonthink` | C10 | 1 | 4101 | `stop` | 462 | 300 | 2048 | None | 18.42 | 1959.7 | 1930 | `results/control_corpus_2b.json`[18] |
| `xfiqvtbb` | `corpus-lfm3-think-512` | C07 | 1 | 4101 | `stop` | 227 | 590 | 1024 | 512 | 29.255 | 1774.8 | 2040 | `results/corpus.json`[120] |
| `xgza5wzr` | `corpus-4b-nonthink` | C03 | 2 | 4102 | `stop` | 207 | 56 | 768 | None | 7.662 | 3556.4 | 1933 | `results/corpus.json`[7] |
| `xifq9fpk` | `lfm120-ins-t05` | — | 3 | 4103 | `stop` | 244 | 80 | 2048 | None | 2.814 | 1098.1 | 2025 | `results/stageB.json`[5] |
| `xjbgxiyy` | `corpus-tinyllama` | C06 | 1 | 4101 | `length` | 273 | 768 | 768 | None | 24.026 | 1184.0 | 1901 | `results/corpus.json`[209] |
| `xrg5k7bf` | `corpus-tinyllama` | C05 | 1 | 4101 | `stop` | 270 | 384 | 768 | None | 11.654 | 1144.7 | 1901 | `results/corpus.json`[207] |
| `xswefny0` | `lfm120-instruct-official` | — | 2 | 4102 | `stop` | 244 | 77 | 2048 | None | 2.55 | 1098.0 | 1870 | `results/stageA.json`[1] |
| `xvbw9h3m` | `smollm3-3b-think-rb1024` | — | 3 | 4103 | `stop` | 467 | 545 | 3072 | 1024 | 47.527 | 2688.6 | 1964 | `results/stageA2b2.json`[2] |
| `xy7owo3j` | `corpus-lfm120-instruct` | C06 | 3 | 4103 | `stop` | 238 | 43 | 512 | None | 1.448 | 1104.0 | 2117 | `results/corpus.json`[71] |
| `y13js660` | `corpus-tinyllama` | C12 | 2 | 4102 | `length` | 276 | 768 | 768 | None | 23.569 | 1387.7 | 2076 | `results/corpus.json`[222] |
| `y7yno9z7` | `corpus-lfm120-t08` | C17 | 1 | 4101 | `stop` | 248 | 65 | 512 | None | 2.335 | 1110.0 | 2029 | `results/corpus.json`[194] |
| `ydx27b62` | `abl-4b-nonthink` | C03 | 1 | 4101 | `stop` | 207 | 60 | 768 | None | 9.04 | 3520.5 | 1963 | `results/ablation_corpus_abl-4b-nonthink.json`[4] |
| `yevx1ezy` | `abl-4b-nonthink` | — | 1 | 4101 | `stop` | 230 | 86 | 768 | None | 12.425 | 3400.1 | 1934 | `results/ablation_triage.json`[0] |
| `yfkd1dy9` | `abl-2b-nonthink` | C04 | 2 | 4102 | `stop` | 228 | 88 | 768 | None | 5.25 | 1898.5 | 1929 | `results/ablation_corpus_abl-2b-nonthink.json`[7] |
| `yfrcuwqo` | `corpus-lfm120-instruct` | C04 | 1 | 4101 | `stop` | 232 | 45 | 512 | None | 1.652 | 1103.9 | 1891 | `results/corpus.json`[63] |
| `yjyha4dh` | `qwen35-2b-think` | — | 1 | 4101 | `length` | 228 | 3072 | 3072 | None | 190.28 | 1823.7 | 1899 | `results/stageA.json`[15] |
| `yl5aneya` | `qwen35-2b-nonthink-pres0` | — | 2 | 4102 | `stop` | 230 | 185 | 2048 | None | 11.418 | 1821.0 | 2027 | `results/stageB.json`[13] |
| `yr6kscpq` | `abl-4b-nonthink` | C08 | 2 | 4102 | `stop` | 196 | 62 | 768 | None | 8.601 | 3636.6 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[15] |
| `yr7vog1a` | `abl-2b-nonthink` | C12 | 1 | 4101 | `stop` | 236 | 131 | 768 | None | 8.08 | 2037.4 | 1930 | `results/ablation_corpus_abl-2b-nonthink.json`[22] |
| `yteygv3u` | `lfm120-ins-t03` | — | 1 | 4101 | `stop` | 244 | 78 | 2048 | None | 2.995 | 1097.8 | 2025 | `results/stageB.json`[0] |
| `ytjbxqh5` | `abl-4b-nonthink` | C01 | 1 | 4101 | `stop` | 232 | 60 | 768 | None | 9.095 | 3399.9 | 1959 | `results/ablation_corpus_abl-4b-nonthink.json`[0] |
| `yymn89gq` | `abl-4b-nonthink` | C07 | 1 | 4101 | `stop` | 226 | 64 | 768 | None | 9.461 | 3636.6 | 1965 | `results/ablation_corpus_abl-4b-nonthink.json`[12] |
| `z231wxrn` | `corpus-tinyllama` | C03 | 1 | 4101 | `stop` | 247 | 327 | 768 | None | 10.005 | 1107.8 | 1901 | `results/corpus.json`[203] |
| `z4avor13` | `corpus-4b-nonthink` | C17 | 3 | 4103 | `stop` | 239 | 34 | 768 | None | 4.693 | 3681.7 | 1935 | `results/corpus.json`[50] |
| `z8dtj0hr` | `lfm3-17b-think-rb512` | — | 1 | 4101 | `stop` | 238 | 618 | 3072 | 512 | 31.592 | 1558.2 | 2019 | `results/budget.json`[0] |
| `z8ggsop3` | `lfm3-17b-think-rb1024` | — | 3 | 4103 | `stop` | 238 | 1157 | 3072 | 1024 | 56.894 | 1608.4 | 1919 | `results/stageA2a.json`[14] |
| `zct8gk6f` | `lfm120-ins-t08` | — | 3 | 4103 | `stop` | 244 | 73 | 2048 | None | 2.595 | 1097.6 | 2025 | `results/stageB.json`[8] |
| `zhsjs4wh` | `corpus-tinyllama` | C10 | 2 | 4102 | `stop` | 535 | 263 | 768 | None | 8.081 | 1330.7 | 2076 | `results/corpus.json`[218] |
| `zhtq5j8z` | `qwen35-4b-nonthink` | — | 2 | 4102 | `stop` | 230 | 79 | 2048 | None | 10.801 | 3437.7 | 1913 | `results/stageA2a.json`[1] |
| `zijgcap7` | `corpus-lfm120-instruct` | C13 | 2 | 4102 | `stop` | 258 | 64 | 512 | None | 2.107 | 1111.3 | 2117 | `results/corpus.json`[91] |
| `zju2qklf` | `corpus-4b-nonthink` | C15 | 2 | 4102 | `stop` | 249 | 71 | 768 | None | 9.721 | 3681.7 | 1935 | `results/corpus.json`[43] |
| `zmtxgxu0` | `abl-lfm120` | — | 3 | 4103 | `stop` | 244 | 64 | 512 | None | 2.149 | 1136.9 | 1892 | `results/ablation_triage.json`[11] |
| `zqovi3yt` | `qwen35-2b-nonthink` | C04 | 2 | 4102 | `stop` | 228 | 233 | 2048 | None | 13.972 | 1866.7 | 1929 | `results/control_corpus_2b.json`[7] |
| `zvtwnqsg` | `smollm3-3b-nonthink` | — | 1 | 4101 | `stop` | 287 | 87 | 2048 | None | 7.896 | 2619.5 | 1964 | `results/stageA2b2.json`[3] |
| `zylyqaux` | `corpus-4b-think-512` | C12 | 1 | 4101 | `stop` | 234 | 609 | 1024 | 512 | 86.582 | 4897.0 | 2193 | `results/corpus.json`[155] |

## 4. What these values do not show

- Whether the server applied the requested sampler chain. llama.cpp accepts the
  values; no echo of the effective chain is stored.
- Whether `max_tokens`, `reasoning_budget` or `ctx` bound a particular run. The
  recorded `finish_reason` and `completion_tokens` are the only direct indicators.
- Whether measured memory belongs to the model alone. `peak_vram_mb` is
  `nvidia-smi memory.used`, which includes other consumers on the desktop.
