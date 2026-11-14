# `MEMORY-09` — Model Population and Hardware Plan

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-MPHP-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **PLAN — deliberate non-exhaustive population; nothing loaded or measured.**
**Authority:** `MEMORY-09-WO-ARCH-06` §5.

---

## 1. Hardware envelope (owner-provided — not independently verified)

| Component | Value |
|---|---|
| CPU | Intel i7-10700K |
| RAM | 64 GB DDR4 |
| GPU | GTX 1080, 8 GB VRAM |
| OS | Windows 11 Pro |

**Treat as stated, not measured.** Any figure in this plan derived from these numbers is an **estimate**, labelled as such, pending measurement (§5).

## 2. Eligibility rule

> **If it can be made to load and run on this machine, it is eligible.**

Explicitly **not** eligibility conditions:

- parameter count (no arbitrary ceiling),
- generation speed (slow models are eligible and often informative),
- category fit.

**Only genuine blockers exclude a model:** it will not load, it crashes, or it produces no usable output at all. Failure to load is recorded as a *result about this hardware*, not as a model's deficiency.

## 3. Why a non-exhaustive population is correct here

The design goal is **variance across kinds of model**, not coverage. Ten near-identical 7B general models would cost ten times the effort of three and add nothing but noise. The categories below are chosen because each plausibly changes *what the architecture can do* — and that is the variable under test.

| # | Category | Why it belongs |
|---|---|---|
| 1 | Creative-writing fine-tunes | Primary category: the product's stated purpose is companion/roleplay quality |
| 2 | General instruct/chat | The reference class most users would load |
| 3 | Reasoning / thinking | Tests whether long internal deliberation helps or harms in-character persistence |
| 4 | Non-thinking conversational | The contrast case for 3, at similar size where possible |
| 5 | Different families/architectures | Prompt-format and sampler sensitivity are family effects, not model-quality effects |
| 6 | Coding-oriented (lower priority) | Included, not excluded — a weak companion model with strong instruction-following isolates architecture from prose flair |
| 7 | Hardware-heavy | At least one that meaningfully challenges the envelope, to test whether architecture still helps when the model is strained |

## 4. Candidate matrix

**No specification below is asserted.** Parameters, quantisation and context are read from the GGUF metadata at load time and written into the record; "unknown" is the honest value until then.

| Slot | Candidate (family) | Category | CW relevance | Params/quant | Context | Runtime | Expected demand | Fit status | Inclusion reason |
|---|---|---|---|---|---|---|---|---|---|
| M-01 | Qwen3 series (app already carries published defaults, `app-family-defaults.go`) | 2 / 3 | Medium | *read metadata* | *read metadata* | llama.cpp | Moderate | **NOT TESTED** | Family with **published** sampling numbers — a clean baseline test |
| M-02 | Gemma series (published top-k 64 / top-p .95 in app table) | 2 | Medium | *read metadata* | *read metadata* | llama.cpp | Moderate | **NOT TESTED** | Second family with published numbers; different sampler shape |
| M-03 | DeepSeek-R1 distil (published T 0.5–0.7) | 3 | Low-Medium | *read metadata* | *read metadata* | llama.cpp | Moderate | **NOT TESTED** | Reasoning model with a published range, versus M-04 |
| M-04 | Llama 3.x instruct | 2 / 4 | Medium | *read metadata* | *read metadata* | llama.cpp | Moderate | **NOT TESTED** | Family with **no published numbers** — tests the "generic default is the honest choice" rule (`app-family-defaults.go` comment) |
| M-05 | Mistral / Nemo instruct | 2 / 4 | Medium | *read metadata* | *read metadata* | llama.cpp | Moderate | **NOT TESTED** | Second no-published-numbers family; long-context reputation to verify, not assume |
| M-06 | Creative-writing fine-tune, mid-size | **1** | **High** | *read metadata* | *read metadata* | llama.cpp | Moderate-High | **NOT TESTED** | Primary category; prose quality is the product surface |
| M-07 | Creative-writing fine-tune, small | **1** | **High** | *read metadata* | *read metadata* | llama.cpp | Low | **NOT TESTED** | Cheap control for M-06: does architecture compensation shrink as prose quality rises? |
| M-08 | **Midnight Rose / Midnight Miqu** (owner-reported ~40B, ~8K context) | **1 / 7** | **High** | 40B class, quant *read metadata* | ~8K *(owner-reported)* | llama.cpp | **Heavy** — CPU offload expected | **PREVIOUSLY RAN** *(owner-reported)* | The hardware-heavy representative that already ran once; ~8K context makes memory pressure real |
| M-09 | Coding-oriented instruct | 6 | Low | *read metadata* | *read metadata* | llama.cpp | Low-Moderate | **NOT TESTED** | Instruction-following without prose flair — separates compliance from style |
| M-10 | A second ~30–40B quantised, different family | 7 | Medium | *read metadata* | *read metadata* | llama.cpp | Heavy | **NOT TESTED** | Second heavy point, to see whether "heavy" effects are model-specific or general |

**Late additions are expected and allowed.** The matrix records *why* each model is present, so a later reader can tell a deliberate contrast from a convenient extra.

## 5. Fit estimation — arithmetic, not assertion

Rough weight footprint: **bytes ≈ params × bits/8**.

| Class | ~Q4_K_M footprint | Fits 8 GB VRAM? | Plausible path |
|---|---|---|---|
| 7–8B | ~4.5–5 GB | **Yes** | Full GPU |
| 12–14B | ~7–9 GB | Marginal | Partial offload |
| 30–34B | ~18–21 GB | No | CPU/GPU split (slow) |
| 40B | ~23–25 GB | No | CPU-heavy split (very slow) |
| 70B | ~40–43 GB | No | RAM-feasible, speed-poor |

**These are estimates for planning only.** Actual fit depends on quantisation, KV-cache at the chosen context, runtime overhead and offload behaviour — all of which §6 measures rather than predicts.

## 6. Measurement plan (bounded, safe)

Per model that loads: load time; time to first token; generation speed (tokens/sec, short sample); peak system RAM; peak VRAM; GPU utilisation where available; effective context capacity; stability across a **short** extended session; and unload/switch behaviour if relevant.

**Safety limits:** no sustained saturation, no multi-hour stress runs, no thermal endurance testing. A measurement that risks the machine is not worth the number.

## 7. What this plan must not be read as claiming

- That any candidate fits — **fit is unknown until loaded** (§5 estimates notwithstanding).
- That Midnight Miqu's ~40B / ~8K figures are verified — they are **owner-reported**.
- That the population is representative of "models in general" — it is a deliberate, documented, non-exhaustive sample.
- That slow models are failures — speed is measured, never a gate (§2).
