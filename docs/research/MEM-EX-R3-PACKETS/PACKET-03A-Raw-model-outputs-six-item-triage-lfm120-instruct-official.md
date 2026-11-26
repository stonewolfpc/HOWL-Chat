# Packet 03A — Raw model outputs — six-item triage: lfm120-instruct-official

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `results/all_triage.json`, `results/*_scores.json`
**Case / run identifiers covered:** lfm120-instruct-official
**Raw model output in this packet:** yes, verbatim
**Omitted:** none — every stored triage run appears in this series

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

This packet reproduces the six-item prompt runs exactly as stored in
`results/all_triage.json`. The prompt itself is quoted in full at the top of the
first triage packet and in Packet 01; it is not repeated per run.

Two things to watch, both `OBSERVED`:

1. **Reasoning leakage.** Some models put working notes (`Analysis:`,
   `Decision:`) into `content`. The v1 triage scorer marks such runs
   inadmissible; their scores are shown but should not be read as verdicts.
2. **Non-binding budgets.** Where `reasoning_budget` is set but
   `completion_tokens` exceeds it substantially, the budget did not bind. The
   report treats those as protocol failures, not model quality.

The 12-point triage score is the historical instrument and is unchanged in v2.


---
## The six-item prompt, verbatim

`SOURCE-DERIVED` from `experiments/lfm25-cpu/run_probe.py` (`PROMPT`), which round 3 imports rather than retyping.

```text
You are triaging candidate memories extracted from a roleplay session before they are written to long-term memory. For each item below, decide KEEP (durable, worth remembering), REJECT (transient or worthless), or HOLD (needs the user's confirmation), and give a reason of no more than twelve words. Watch for duplicates, contradictions, and anything marked as uncertain. Respond with exactly one line per item in the form `ID | DECISION | REASON`, no preamble and no extra commentary. Items: M1 — Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user). M2 — The user said "brb" and returned four minutes later. M3 — Sable is a fire mage who cannot swim. M4 — The user first said Sable has never met Corvin, then later said Corvin trained her. M5 — The user prefers short replies and dislikes long scene descriptions. M6 — A note says Sable's birthday is 3 Frostfell, but the user tagged it "maybe canon".
```

## lfm120-instruct-official

`SOURCE-DERIVED`: 3 run(s) on the six-item prompt. Sampler values are in each block.

#### lfm120-instruct-official — draw 1

* ``ez3jx5az` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 78`
* `wall 2.788 s`
* `v1 score 8`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference to "brb"  
M3 | HOLD | Needs clarification on fire mage traits  
M4 | HOLD | Conflicting training details  
M5 | KEEP | Preferred concise format  
M6 | HOLD | Birthday ambiguity requires verification
```

*`reasoning_content`: empty.*

#### lfm120-instruct-official — draw 2

* ``xswefny0` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 77`
* `wall 2.55 s`
* `v1 score 8`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference to "brb"  
M3 | HOLD | Needs clarification on mage abilities  
M4 | HOLD | Conflicting training details  
M5 | KEEP | Prefers concise responses  
M6 | HOLD | Birthday ambiguity requires verification
```

*`reasoning_content`: empty.*

#### lfm120-instruct-official — draw 3

* ``f5fp8urt` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 78`
* `wall 2.551 s`
* `v1 score 8`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference to "brb"  
M3 | HOLD | Needs clarification on fire mage traits  
M4 | HOLD | Conflicting training details  
M5 | KEEP | Preferred concise format  
M6 | HOLD | Birthday ambiguity requires verification
```

*`reasoning_content`: empty.*
