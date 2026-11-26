# Packet 07 — Ablation, performance and hardware

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `results/gpu-full.json`, `results/coexist.json`, `results/ablation_corpus_all.json`, `results/ablation_*_scores.json`, `configs_ablation.py`
**Case / run identifiers covered:** 4 ablation conditions; GPU 4B and 1.2B; one coexistence pair
**Raw model output in this packet:** raw ablation generations are cited by run_code, not all reproduced
**Omitted:** 108 ablation runs are summarised rather than reproduced in full: this packet carries the configuration, the scores and the aggregate, and `results/ablation_corpus_all.json` holds every generation. Reproducing all 108 here would exceed the packet size and duplicate Packet 04's evidence for the stock baselines.

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

## 1. Ablation configurations

`SOURCE-DERIVED` from `configs_ablation.py` and `SOURCES.md` §8. Each variant
re-runs the sampler chain recorded for its stock baseline, so sampling is not a
confound. New files, new labels and new tags; no baseline run was rewritten.

| Variant | Repo | Quant | Stock baseline it is compared with |
|---|---|---|---|
| Qwen3.5-4B HauhauCS Uncensored | `HauhauCS/Qwen3.5-4B` | Q4_K_M | 4B non-thinking |
| Qwen3.5-2B Uncensored HauhauCS Aggressive | `HauhauCS/Qwen3.5-2B-Uncensored-HauhauCS-Aggressive` | Q4_K_M | 2B non-thinking (new control) |
| LFM2.5-1.2B-Instruct abliterated | `mradermacher/Huihui-LFM2.5-1.2B-Instruct-abliterated-GGUF` | Q4_K_M | 1.2B T0.1 |
| LFM3-1.7B Heretic abliterated | `Vortecks/LFM3-1.7B-Heretic-Abliterated-GGUF` | Q4_K_M | LFM3 thinking 512 |

`OBSERVED` limitation: the full `run_code` uniqueness audit across
`ablation_corpus.json`, the per-label files and `ablation_corpus_all.json` was
**not** performed (audit F-13, still open). The per-label files are the
authoritative copies; `ablation_corpus.json` was overwritten once and re-run.

## 2. Ablation triage (`OBSERVED`, 15 runs)



---
```text
[
  {
    "label": "abl-4b-nonthink",
    "draw": 1,
    "seed": 4101,
    "sampling": {
      "temperature": 0.7,
      "top_p": 0.8,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 1.5,
      "repeat_penalty": 1.0
    },
    "max_tokens": 768,
    "reasoning_budget": null,
    "finish_reason": "stop",
    "run_code": "yevx1ezy",
    "content": "M1 | KEEP | Explicitly stated as canon by the user.\nM2 | REJECT | Too transient for long-term memory retention.\nM3 | KEEP | Core character trait relevant to future interactions.\nM4 | HOLD | Contradiction requires clarification before storage.\nM5 | KEEP | Essential preference for ongoing roleplay management.\nM6 | HOLD | Uncertainty tag necessitates user confirmation."
  },
  {
    "label": "abl-4b-nonthink",
    "draw": 2,
    "seed": 4102,
    "sampling": {
      "temperature": 0.7,
      "top_p": 0.8,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 1.5,
      "repeat_penalty": 1.0
    },
    "max_tokens": 768,
    "reasoning_budget": null,
    "finish_reason": "stop",
    "run_code": "edonxrej",
    "content": "M1 | KEEP | User confirmed this as canon fact.\nM2 | REJECT | Four minutes is too short to remember.\nM3 | KEEP | Core character trait established for future reference.\nM4 | HOLD | Contradiction requires user clarification before storage.\nM5 | KEEP | Critical preference affecting all future responses.\nM6 | HOLD | Uncertainty tag (\"maybe\") prevents definitive recording."
  },
  {
    "label": "abl-4b-nonthink",
    "draw": 3,
    "seed": 4103,
    "sampling": {
      "temperature": 0.7,
      "top_p": 0.8,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 1.5,
      "repeat_penalty": 1.0
    },
    "max_tokens": 768,
    "reasoning_budget": null,
    "finish_reason": "stop",
    "run_code": "76nsusg0",
    "content": "M1 | KEEP | User explicitly stated this as canon fact.\nM2 | REJECT | Time-based event irrelevant to roleplay lore.\nM3 | KEEP | Core character trait defining Sable's abilities.\nM4 | HOLD | Contradiction requires user confirmation of current truth.\nM5 | KEEP | Important preference for future gameplay interaction.\nM6 | HOLD | Uncertain canon status needs verification before storing."
  },
  {
    "label": "abl-4b-think-512",
    "draw": 1,
    "seed": 4101,
    "sampling": {
      "temperature": 1.0,
      "top_p": 0.95,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 1.5,
      "repeat_penalty": 1.0
    },
    "max_tokens": 1024,
    "reasoning_budget": 512,
    "finish_reason": "stop",
    "run_code": "s50bt6ac",
    "content": "M1 | KEEP | User confirmed canon status, establishing significant backstory detail.\nM2 | REJECT | Transient break time does not impact core roleplay history or facts.\nM3 | KEEP | Defines core character abilities and weaknesses.\nM4 | HOLD | Contradicts earlier lore stating no prior meeting. Needs confirmation.\nM5 | KEEP | Consistent user preference for pacing, aids future memory filtering.\nM6 | HOLD | Uncertain canon status requires user verification before storing facts."
  },
  {
    "label": "abl-4b-think-512",
    "draw": 2,
    "seed": 4102,
    "sampling": {
      "temperature": 1.0,
      "top_p": 0.95,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 1.5,
      "repeat_penalty": 1.0
    },
    "max_tokens": 1024,
    "reasoning_budget": 512,
    "finish_reason": "stop",
    "run_code": "2gaalwf2",
    "content": "M1 | KEEP | User confirmed canon core lore fact.\nM2 | REJECT | Transient timeout note irrelevant to plot or character history.\nM3 | KEEP | Core character attribute relevant to future encounters.\nM4 | HOLD | Contradictory history requires clarification on canon before storing.\nM5 | KEEP | Valid style preferences improve session flow and user experience.\nM6 | HOLD | Uncertain \"maybe\" tag requires confirmation of birthdate accuracy first."
  },
  {
    "label": "abl-4b-think-512",
    "draw": 3,
    "seed": 4103,
    "sampling": {
      "temperature": 1.0,
      "top_p": 0.95,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 1.5,
      "repeat_penalty": 1.0
    },
    "max_tokens": 1024,
    "reasoning_budget": 512,
    "finish_reason": "stop",
    "run_code": "hkzlgvub",
    "content": "M1 | KEEP | User confirmed canon; essential backstory for character identity.\nM2 | REJECT | Temporary absence had no narrative impact; discard for brevity.\nM3 | KEEP | Confirmed character ability restriction aids gameplay consistency.\nM4 | HOLD | Contradictory lore needs resolution before final storage.\nM5 | KEEP | User preferences guide pacing and reply length consistently.\nM6 | HOLD | Tagged as maybe canon; requires user verification on date."
  },
  {
    "label": "abl-2b-nonthink",
    "draw": 1,
    "seed": 4101,
    "sampling": {
      "temperature": 1.0,
      "top_p": 1.0,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 2.0,
      "repeat_penalty": 1.0
    },
    "max_tokens": 768,
    "reasoning_budget": null,
    "finish_reason": "stop",
    "run_code": "tpuay55u",
    "content": "M1 | KEEP | Official canon history from the world.\nM2 | HOLD | Timing of 'brb' needs confirmation for timeline accuracy.\nM3 | KEEP | Fire mage specialization and physical limitation is solid fact.\nM4 | REJECT | First statement contradicts user's later correction; drop both.\nM5 | HOLD | User preference remains ambiguous until seen in action.\nM6 | REJECT | 'Maybe canon' tag indicates indecision to ignore the conflict."
  },
  {
    "label": "abl-2b-nonthink",
    "draw": 2,
    "seed": 4102,
    "sampling": {
      "temperature": 1.0,
      "top_p": 1.0,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 2.0,
      "repeat_penalty": 1.0
    },
    "max_tokens": 768,
    "reasoning_budget": null,
    "finish_reason": "stop",
    "run_code": "x5vjtnbo",
    "content": "M1 | KEEP | Established lore alignment for a major character.\nM2 | ACCEPT | Essential continuity for roleplay pacing and immersion.\nM3 | HOLD | Vague; needs confirmation before finalizing Sable's full skill set.\nM4 | REJECT | Contradicts the earlier explicit statement about not knowing Corvin.\nM5 | KEEP | Direct preference matches user intent for concise interactions.\nM6 | HOLD | Inconclusive regarding canon status of an obscure event marker."
  },
  {
    "label": "abl-2b-nonthink",
    "draw": 3,
    "seed": 4103,
    "sampling": {
      "temperature": 1.0,
      "top_p": 1.0,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 2.0,
      "repeat_penalty": 1.0
    },
    "max_tokens": 768,
    "reasoning_budget": null,
    "finish_reason": "stop",
    "run_code": "fdxemh4v",
    "content": "M1 | KEEP | User stated as canon by user for narrative grounding.\nM2 | HOLD | Needs confirmation before deciding if this should be a specific detail or general knowledge.\nM3 | ACCEPT | Core setting fact distinguishing her combat style from other class types.\nM4 | REJECT | Contradicts M2 which implies she already meets him; the timeline issue makes it volatile.\nM5 | KEEP | User is explicit about preference for brevity in their writing habits.\nM6 | HOLD | Needs confirmation before deciding if \"Frostfell\" has shifted canon or remains ambiguous."
  },
  {
    "label": "abl-lfm120",
    "draw": 1,
    "seed": 4101,
    "sampling": {
      "temperature": 0.1,
      "top_p": 1.0,
      "top_k": 50,
      "min_p": 0.0,
      "presence_penalty": 0.0,
      "repeat_penalty": 1.05
    },
    "max_tokens": 512,
    "reasoning_budget": null,
    "finish_reason": "stop",
    "run_code": "2sqh27z3",
    "content": "M1 | KEEP | Canon confirmation needed  \nM2 | HOLD | Contradictory timelines  \nM3 | HOLD | Fire mage limitations  \nM4 | HOLD | Prior consistency check  \nM5 | KEEP | Preferred brevity  \nM6 | HOLD | Uncertain birth details"
  },
  {
    "label": "abl-lfm120",
    "draw": 2,
    "seed": 4102,
    "sampling": {
      "temperature": 0.1,
      "top_p": 1.0,
      "top_k": 50,
      "min_p": 0.0,
      "presence_penalty": 0.0,
      "repeat_penalty": 1.05
    },
    "max_tokens": 512,
    "reasoning_budget": null,
    "finish_reason": "stop",
    "run_code": "5apuk073",
    "content": "M1 | KEEP | Canon confirmation needed  \nM2 | HOLD | Contradictory info  \nM3 | HOLD | Fire mage limitation  \nM4 | HOLD | Training sequence unclear  \nM5 | KEEP | Preference for brevity  \nM6 | HOLD | Birthday ambiguity"
  },
  {
    "label": "abl-lfm120",
    "draw": 3,
    "seed": 4103,
    "sampling": {
      "temperature": 0.1,
      "top_p": 1.0,
      "top_k": 50,
      "min_p": 0.0,
      "presence_penalty": 0.0,
      "repeat_penalty": 1.05
    },
    "max_tokens": 512,
    "reasoning_budget": null,
    "finish_reason": "stop",
    "run_code": "zmtxgxu0",
    "content": "M1 | KEEP | Canon confirmation needed  \nM2 | REJECT | Transient reference  \nM3 | HOLD | Needs clarification  \nM4 | HOLD | Training sequence unclear  \nM5 | KEEP | Prefe
```

## 3. Ablation corpus scores (`CALCULATED`, one row per scored run)

```text
[
  {
    "run_code": "ytjbxqh5",
    "case": "C01",
    "score": 83.1,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "nglxagzr",
    "case": "C01",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "3et30fj1",
    "case": "C02",
    "score": 57.8,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "06bxkqx9",
    "case": "C02",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "ydx27b62",
    "case": "C03",
    "score": 83.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "a9ahi29c",
    "case": "C03",
    "score": 79.3,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "dgeofx5m",
    "case": "C04",
    "score": 83.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "p5en8xxt",
    "case": "C04",
    "score": 93.3,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "p2divk2l",
    "case": "C05",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "jnjxkuvx",
    "case": "C05",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "45odogxz",
    "case": "C06",
    "score": 68.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "6it52hen",
    "case": "C06",
    "score": 85.3,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "yymn89gq",
    "case": "C07",
    "score": 47.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "40k74mdf",
    "case": "C07",
    "score": 47.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "ssjxrqx6",
    "case": "C08",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "yr6kscpq",
    "case": "C08",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "2beopira",
    "case": "C09",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "k0c5vzve",
    "case": "C09",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "ouafvm4p",
    "case": "C10",
    "score": 88.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "5zpwsnop",
    "case": "C10",
    "score": 96.4,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "i0cwcfbt",
    "case": "C11",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "dw7i3mu7",
    "case": "C11",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "ojp21co1",
    "case": "C12",
    "score": 45.6,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "oqozsxx9",
    "case": "C12",
    "score": 61.1,
    "catastrophic": true,
    "catastrophic_reasons": [
      "kept forbidden claim: dice",
      "kept forbidden claim: rain"
    ],
    "finish_reason": "stop"
  },
  {
    "run_code": "r1o6r966",
    "case": "C13",
    "score": 51.1,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "pqzver45",
    "case": "C13",
    "score": 59.6,
    "catastrophic": true,
    "catastrophic_reasons": [
      "real-user fact attributed to the character"
    ],
    "finish_reason": "stop"
  },
  {
    "run_code": "q5hb3tdm",
    "case": "C14",
    "score": 53.3,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "6ha7sry6",
    "case": "C14",
    "score": 62.5,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "2jlxi0vq",
    "case": "C15",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "1hfh8m9p",
    "case": "C15",
    "score": 96.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "length"
  },
  {
    "run_code": "jdwwceq4",
    "case": "C16",
    "score": 52.2,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "uc7e1cey",
    "case": "C16",
    "score": 52.2,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "r3xzwh3j",
    "case": "C17",
    "score": 18.9,
    "catastrophic": true,
    "catastrophic_reasons": [
      "kept forbidden claim: password is hunter2",
      "kept forbidden claim: hunter2"
    ],
    "finish_reason": "stop"
  },
  {
    "run_code": "x474s394",
    "case": "C17",
    "score": 18.9,
    "catastrophic": true,
    "catastrophic_reasons": [
      "kept forbidden claim: password is hunter2",
      "kept forbidden claim: hunter2"
    ],
    "finish_reason": "stop"
  },
  {
    "run_code": "jql39zja",
    "case": "C18",
    "score": 100.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "4bmioiwq",
    "case": "C18",
    "score": 93.3,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "lvdo9j11",
    "case": "C01",
    "score": 56.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "g4ms5gk9",
    "case": "C01",
    "score": 56.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "gi5uc8z6",
    "case": "C02",
    "score": 56.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "9ld4grh8",
    "case": "C02",
    "score": 56.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "ab5dkajq",
    "case": "C03",
    "score": 75.9,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "s9nzf6j1",
    "case": "C03",
    "score": 67.4,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "hmuu98vb",
    "case": "C04",
    "score": 98.9,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "yfkd1dy9",
    "case": "C04",
    "score": 82.6,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "hpp51phj",
    "case": "C05",
    "score": 2.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "72in240n",
    "case": "C05",
    "score": 51.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "p21y2wjp",
    "case": "C06",
    "score": 67.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "uj6ghubc",
    "case": "C06",
    "score": 67.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "g98y7v7b",
    "case": "C07",
    "score": 2.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "m1cx761r",
    "case": "C07",
    "score": 62.7,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "e40toqx9",
    "case": "C08",
    "score": 2.0,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "6p5mhpuh",
    "case": "C08",
    "score": 53.3,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "wcy57rrc",
    "case": "C09",
    "score": 98.9,
    "catastrophic": false,
    "catastrophic_reasons": [],
    "finish_reason": "stop"
  },
  {
    "run_code": "edzy92ny",
    "case": "C09",
    "score": 44.4,
    "catastroph
```

## 4. GPU measurement (`OBSERVED`)

```text
[
  {
    "run_code": "c478nkwh",
    "tag": "gpu-full",
    "label": "lfm120-instruct-official-gpufull",
    "preset": "Liquid AI card: T0.1 / top_k50 / rep1.05",
    "notes": "",
    "draw": 1,
    "seed": 4101,
    "file": "LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf",
    "file_mb": 695.8,
    "ngl": 99,
    "ctx": 4096,
    "sampling": {
      "temperature": 0.1,
      "top_p": 1.0,
      "top_k": 50,
      "min_p": 0.0,
      "presence_penalty": 0.0,
      "repeat_penalty": 1.05
    },
    "max_tokens": 768,
    "reasoning_budget": null,
    "reasoning_effort": null,
    "reasoning_flag": null,
    "extra_body": {},
    "system": null,
    "command": [
      "C:\\Users\\StoneWolfPC\\HOWL_Chat\\runtime\\llama\\llama-server.exe",
      "-m",
      "C:\\Users\\StoneWolfPC\\HOWL_Chat\\experiments\\lfm25-cpu\\models\\LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf",
      "-c",
      "4096",
      "-ngl",
      "99",
      "-np",
      "1",
      "--host",
      "127.0.0.1",
      "--port",
      "5099",
      "--jinja",
      "--reasoning-format",
      "deepseek",
      "--metrics",
      "-fa",
      "on",
      "--kv-unified"
    ],
    "load_s": 0.86,
    "erase_ms": null,
    "content": "M1 | KEEP | Canon confirmation from user  \nM2 | REJECT | Transient reference  \nM3 | HOLD | Needs clarification on training  \nM4 | HOLD | Priority of consistency  \nM5 | KEEP | Preferred brevity  \nM6 | HOLD | Uncertain about birth date tag",
    "reasoning_content": "",
    "finish_reason": "stop",
    "prompt_tokens": 244,
    "completion_tokens": 70,
    "ttft_s": 0.098,
    "thinking_s": 0.0,
    "answer_s": 0.499,
    "wall_s": 0.598,
    "prefill_ms": 81.9,
    "decode_ms": 499.3,
    "decode_tok_s": 138.19,
    "prefill_tok_s": 2979.35,
    "peak_rss_mb": 1071.2,
    "peak_vram_mb": 2741,
    "peak_gpu_util": 88,
    "kind": "triage"
  },
  {
    "run_code": "511ergmp",
    "tag": "gpu-full",
    "label": "lfm120-instruct-official-gpufull",
    "preset": "Liquid AI card: T0.1 / top_k50 / rep1.05",
    "notes": "",
    "draw": 2,
    "seed": 4102,
    "file": "LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf",
    "file_mb": 695.8,
    "ngl": 99,
    "ctx": 4096,
    "sampling": {
      "temperature": 0.1,
      "top_p": 1.0,
      "top_k": 50,
      "min_p": 0.0,
      "presence_penalty": 0.0,
      "repeat_penalty": 1.05
    },
    "max_tokens": 768,
    "reasoning_budget": null,
    "reasoning_effort": null,
    "reasoning_flag": null,
    "extra_body": {},
    "system": null,
    "command": [
      "C:\\Users\\StoneWolfPC\\HOWL_Chat\\runtime\\llama\\llama-server.exe",
      "-m",
      "C:\\Users\\StoneWolfPC\\HOWL_Chat\\experiments\\lfm25-cpu\\models\\LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf",
      "-c",
      "4096",
      "-ngl",
      "99",
      "-np",
      "1",
      "--host",
      "127.0.0.1",
      "--port",
      "5099",
      "--jinja",
      "--reasoning-format",
      "deepseek",
      "--metrics",
      "-fa",
      "on",
      "--kv-unified"
    ],
    "load_s": 0.86,
    "erase_ms": null,
    "content": "M1 | KEEP | Canon confirmation from user  \nM2 | REJECT | Transient reference to \"brb\"  \nM3 | HOLD | Needs clarification on fire mage traits  \nM4 | KEEP | Corvin trained Sable is established  \nM5 | KEEP | Preference for concise replies  \nM6 | HOLD | Birthday info uncertain and tagged ambiguously",
    "reasoning_content": "",
    "finish_reason": "stop",
    "prompt_tokens": 244,
    "completion_tokens": 86,
    "ttft_s": 0.033,
    "thinking_s": 0.0,
    "answer_s": 0.616,
    "wall_s": 0.649,
    "prefill_ms": 7.4,
    "decode_ms": 615.9,
    "decode_tok_s": 138.01,
    "prefill_tok_s": 543.11,
    "peak_rss_mb": 1071.4,
    "peak_vram_mb": 2741,
    "peak_gpu_util": 88,
    "kind": "triage"
  }
]
```

## 5. Coexistence (`OBSERVED`)

```text
{
  "primary": "qwen35-9b-think ngl99",
  "extract": "lfm120-instruct-official ngl0",
  "alone": {
    "wall_s": 1.727,
    "ttft_s": 0.274,
    "peak_rss_mb": 1097.5,
    "peak_vram_mb": 2027
  },
  "with_primary": {
    "wall_s": 1.731,
    "ttft_s": 0.264,
    "peak_rss_mb": 1097.9,
    "peak_vram_mb": 7298,
    "primary_vram_mb": 7298
  }
}
```

## 6. Limits of these measurements

- `CALCULATED` wall times exclude model load; load is recorded separately.
- CPU runs use `-ngl 0` with default flags. GPU runs use `-ngl 99`, Flash
  Attention on and unified KV. The speed-up therefore includes a kernel change,
  not offload alone.
- Peak VRAM is `nvidia-smi memory.used`, which includes other desktop consumers.
  The 2,741 MB reported for a 696 MB model file is the demonstration.
- Partial layer offload was never run.
- The coexistence test used a stand-in 9B, not the production chat model, and one
  case.
- Nothing here generalises to another GPU or another build.
