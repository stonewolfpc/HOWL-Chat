# Packet 01 — Experimental design and protocol

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `harness.py`, `configs.py`, `corpus.json`, `RUBRIC.md`, `SOURCES.md`
**Case / run identifiers covered:** all 18 cases and 51 triage runs
**Raw model output in this packet:** n/a — design only
**Omitted:** nothing omitted

**Labels used throughout:** `OBSERVED` directly present in a stored response or
execution record · `SOURCE-DERIVED` stated in a source or configuration file ·
`CALCULATED` derived by explicit arithmetic · `INFERRED` a reading of the
evidence · `UNESTABLISHED` not established by the available evidence ·
`SUPERSEDED` an earlier claim replaced by later evidence.

**Whose words you are reading:** fenced blocks are the model's own output
(`OBSERVED`, verbatim). Prompt text is `SOURCE-DERIVED` from `corpus.json`.
Scores are `CALCULATED`. Any sentence outside a fence that is not a definition is
Buffy's interpretation unless it is attributed to Muse or to the project owner.

**Caution.** These are exploratory measurements. Nothing here selects a model,
establishes memory quality, or supplies MEMORY-09 evidence. Sample sizes are
1–3 draws per case.

---

## 1. Research questions

The experiment asked four separable questions. They are listed because a later
reader should be able to see which one any number answers.

1. **Quality.** Given one conversation, does a local model produce memory
   candidates that a careful reader would keep — correct content, correct
   decision, correct speaker, correct scope? `OBSERVED` via the 18-case corpus.
2. **Contract.** Does it produce the required four-field line
   (`ID | DECISION | MEMORY | REASON`) without a preamble, and stop cleanly?
   `OBSERVED` via parse status.
3. **Cost.** What latency, RAM and VRAM does that take on this machine, on CPU
   and on one 8 GB GPU? `OBSERVED` from wall time and sampled memory.
4. **Isolation.** Does one extraction job's context affect the next job?
   `OBSERVED` narrowly, by sentinel probe.

It did **not** ask which model should be deployed, whether the memory design is
correct, or how a production system should store records.

## 2. Procedure

`SOURCE-DERIVED` from `harness.py`, `run_stage_a.py`, `run_stage_a2.py`,
`run_stage_b.py`, `run_stage_c.py`, `run_budget.py`, `run_corpus_plan.py`,
`run_gpu.py`, `run_coexist.py`, `run_reset.py`, `run_isolation.py`.

1. **Stage A** — one publisher-preset configuration per model and mode on the
   original six-item prompt, three draws, CPU only.
2. **Stage A2** — add models and modes that Stage A did not cover, same prompt.
3. **Stage B** — bounded sampler ladder on the two weakest useful models:
   `presence_penalty` 2.0 → 0.0 for Qwen3.5-2B; temperature 0.3 / 0.5 / 0.8 for
   LFM2.5-1.2B-Instruct.
4. **Stage C** — rebuild the round-2 sampler chain to test whether the published
   round-2 numbers reproduce. They did not (see §4).
5. **Budget ladder** — `--reasoning-budget` 256 / 512 / 1024 where the build
   supports it.
6. **Corpus** — 18 cases × the selected conditions, 1–3 draws, 235 stored runs.
7. **GPU** — full offload for two models, two draws.
8. **Coexistence** — a 9B chat model resident on the GPU while a 1.2B extractor
   runs on CPU.
9. **Reset / isolation** — the erase endpoint, then a sentinel probe.

## 3. The six-item prompt

`SOURCE-DERIVED`, imported byte-for-byte from
`experiments/lfm25-cpu/run_probe.py` (`PROMPT`) so it cannot drift between
rounds. It asks for one line per item, `ID | DECISION | REASON`, no preamble,
for six candidate memories (`M1`–`M6`) covering canon confirmation, chatter, a
character fact, a contradiction, a preference and an uncertain date.

## 4. The corpus

`SOURCE-DERIVED` from `corpus.json`. Eighteen cases, one per category:

| ID | Category | Length |
|---|---|---|
| C01 | explicit long-term preference | short |
| C02 | implicit preference via behaviour | medium |
| C03 | character fact | short |
| C04 | world / scenario rule | medium |
| C05 | correction of earlier information | short |
| C06 | changing, superseded preference | medium |
| C07 | unresolved contradiction | short |
| C08 | uncertain / speculative claim | short |
| C09 | transient chatter | short |
| C10 | important detail buried in long context | long (~250 words) |
| C11 | context-dependent significance | medium |
| C12 | must not be extracted (voided narration) | short |
| C13 | fiction vs real-world assertion | medium |
| C14 | multiple speakers / attribution | medium |
| C15 | repeated information / duplicates | medium |
| C16 | indirect information | short |
| C17 | adversarial embedded instruction | medium |
| C18 | exact wording and numbers | short |

Each case carries the expected decision per required item, keyword groups, and
forbidden strings that indicate fabrication or reversal. Two cases (C09, C12)
have no required items: the correct answer is to store nothing.

`UNESTABLISHED`: the reference answers were written before the runs according to
the file's own `status` field, but the file is untracked and no hash was
published beforehand, so the **order of writing cannot be proven** from the
repository.

## 5. Configuration actually sent

`SOURCE-DERIVED` from `configs.py` and the per-run `sampling` records. Every
value is traceable to a publisher statement in `SOURCES.md`; `min_p` is always
sent explicitly because llama.cpp defaults it to 0.05.

| Condition | Sampling sent | Max tokens | Thinking control |
|---|---|---|---|
| Qwen3.5-4B non-thinking | T0.7 / top_p0.8 / top_k20 / min_p0 / pres1.5 / rep1.0 | 768 | `enable_thinking=false` |
| Qwen3.5-4B thinking | T1.0 / top_p0.95 / top_k20 / min_p0 / pres1.5 / rep1.0 | 1024 | `enable_thinking=true`, budget 512 |
| Qwen3.5-2B non-thinking | T1.0 / top_p1.0 / top_k20 / min_p0 / pres2.0 / rep1.0 | 2048 | `enable_thinking=false` |
| Qwen3-1.7B thinking | T0.6 / top_p0.95 / top_k20 / min_p0 / pres0 / rep1.0 | 3072 | `enable_thinking=true` |
| LFM2.5-1.2B-Instruct | T0.1 / top_p1.0 / top_k50 / min_p0 / pres0 / rep1.05 | 512 (corpus) | none available |
| LFM2.5-1.2B-Instruct T0.8 | T0.8 / top_p1.0 / top_k50 / min_p0 / pres0 / rep1.05 | 512 | none available |
| LFM3-1.7B thinking | T0.1 / top_p1.0 / top_k50 / min_p0 / pres0 / rep1.1 | 768 / budget 512 | `enable_thinking=true` |
| TinyLlama-1.1B-Chat | T0.7 / top_p0.95 / top_k50 / min_p0 / pres0 / rep1.0 | 768 | none exists (ctx 2048) |

Server launch for every round-3 session, recorded per run in `command`:
`llama-server -m <gguf> -c 4096 -ngl <0|99> -np 1 --host 127.0.0.1 --port 5099
--jinja --reasoning-format deepseek --metrics`.

## 6. Execution and measurement

`SOURCE-DERIVED` from `harness.py`.

- One `llama-server` process per session; one request per sequence (`-np 1`);
  streaming responses; a random 8-character `run_code` per generation so scoring
  can be blind.
- The harness attempts `POST /slots/0?action=erase` before each request and
  records the result. On this build it returns HTTP 501, so no erase occurred.
- Wall time is measured around the HTTP request and **excludes model load**;
  load time is recorded separately as `load_s`.
- Peak RSS is sampled at 100 ms over the server process tree. Peak VRAM is
  `nvidia-smi memory.used`, which includes other desktop consumers.
- Seeds are `4100 + draw`. They are recorded but are not a variance control:
  llama.cpp sampling is not guaranteed reproducible across builds.

## 7. Context isolation procedure

`SOURCE-DERIVED` from `run_reset.py` and `run_isolation.py`.

Two designs. **Reset**: run benign case C03, then adversarial case C17, then C03
again on the same server, then C03 on a fresh server; compare bytes and read
`cache_n` from the server metrics. **Sentinel**: give the model a unique name
("Marrowgale") and code ("7712") in one job, then ask an unrelated extraction
and probe whether the sentinel surfaces; repeat on a fresh server.

## 8. What was not done

No constrained decoding, no grammar, no DRY, no seed sweep, no partial GPU
offload, no concurrency, no LFM3 GPU run, no second scorer, no inter-rater
reliability, no more than three draws per case. Each is a limitation, not an
oversight: they were outside the round's scope.

## 9. The six-item prompt, verbatim

`SOURCE-DERIVED` from `run_probe.py::PROMPT`.

```text
You are triaging candidate memories extracted from a roleplay session before they are written to long-term memory. For each item below, decide KEEP (durable, worth remembering), REJECT (transient or worthless), or HOLD (needs the user's confirmation), and give a reason of no more than twelve words. Watch for duplicates, contradictions, and anything marked as uncertain. Respond with exactly one line per item in the form `ID | DECISION | REASON`, no preamble and no extra commentary. Items: M1 — Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user). M2 — The user said "brb" and returned four minutes later. M3 — Sable is a fire mage who cannot swim. M4 — The user first said Sable has never met Corvin, then later said Corvin trained her. M5 — The user prefers short replies and dislikes long scene descriptions. M6 — A note says Sable's birthday is 3 Frostfell, but the user tagged it "maybe canon".
```


---
