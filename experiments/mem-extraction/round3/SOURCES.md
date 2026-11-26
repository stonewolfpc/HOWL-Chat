# Model inventory and sources — round 3

All URLs retrieved **2026-11-19**. "Official" means the publisher (model card,
publisher docs, or the publisher's own GGUF repo). Anything else is marked
**community**. Where two official pages disagree, both are recorded and the
conflict is left visible rather than resolved by preference.

Runtime for every run: `runtime/llama/llama-server.exe`, llama.cpp `0.8.0-dev`,
build **12015**, commit **a34de8618**, Clang 20.1.8, win-x64. That build exposes
`--reasoning-budget N`, `--reasoning-effort LEVEL`, `-rea on|off`,
`--reasoning-format deepseek` and `--chat-template-kwargs`; all were verified
with `llama-server --help` on this machine, not assumed.

## 1. Qwen

### Qwen3.5-0.8B / 2B / 4B / 9B

| Claim | Value | Source | Status |
|---|---|---|---|
| Identity | dense text+vision family, 0.8B/2B/4B/9B | `https://huggingface.co/Qwen/Qwen3.5-2B`; `https://huggingface.co/Qwen/Qwen3.5-4B` | official |
| Context | 262,144 native; 4B extensible to ~1,010,000 | Qwen3.5-4B card | official |
| Thinking control | `enable_thinking` boolean only; Qwen3.5 **does not** support Qwen3's `/think` `/nothink` soft switch | `https://huggingface.co/Qwen/Qwen3.5-2B` ("Qwen3.5 does not officially support the soft switch of Qwen3") | official |
| Effort levels | none documented for the open small models; `thinking_budget` is a hosted-API parameter (Chat Completions/DashScope), not a model control | `https://docs.qwencloud.com/developer-guides/text-generation/thinking` | official (API) |
| Default mode, 0.8B/2B | non-thinking | `https://huggingface.co/Qwen/Qwen3.5-0.8B` ("operates in non-thinking mode by default") | official |
| Default mode, 4B/9B | **conflicting**: Qwen card says thinking by default; the Unsloth GGUF guide says reasoning is disabled by default for all four | `https://huggingface.co/Qwen/Qwen3.5-4B`; `https://unsloth.ai/docs/models/qwen3.5` | official vs community, unresolved |
| Sampling, 0.8B/2B non-thinking text | T 1.0, top_p 1.00, top_k 20, min_p 0.0, presence 2.0, rep 1.0 | Qwen3.5-0.8B/2B cards | official |
| Sampling, 0.8B/2B thinking text | T 1.0, top_p 0.95, top_k 20, min_p 0.0, presence 1.5, rep 1.0 | same | official |
| Sampling, 4B/9B non-thinking general | T 0.7, top_p 0.8, top_k 20, min_p 0.0, presence 1.5, rep 1.0 | Qwen3.5-4B card | official |
| Sampling, 4B/9B thinking general | T 1.0, top_p 0.95, top_k 20, min_p 0.0, presence 1.5, rep 1.0 | Qwen3.5-4B card | official |
| Known issue | thinking can consume the whole completion budget and return empty `content` | community (round‑2 run on this machine) + `https://lyn.one/reasoning-control-flow` | community/observed |
| GGUF | `unsloth/Qwen3.5-*-GGUF`, Q4_K_M files used: 0.8B 532 MB, 2B 1.28 GB, 4B 2.74 GB, 9B 5.68 GB | HF API listing | official-adjacent (Unsloth conversion) |

### Qwen3-1.7B

| Claim | Value | Source | Status |
|---|---|---|---|
| Thinking | **default True**, `enable_thinking` switch, ` thinking...</think>` with token 151668 ending the block | `https://huggingface.co/Qwen/Qwen3-1.7B` | official |
| Sampling, thinking | T 0.6, top_p 0.95, top_k 20, min_p 0 | same card | official |
| Sampling, non-thinking | T 0.7, top_p 0.8, top_k 20, min_p 0 | same card | official |
| Effort levels | none; Qwen3 has only thinking on/off | same card | official |
| GGUF | `unsloth/Qwen3-1.7B-GGUF` Q4_K_M 1.11 GB | HF API | official-adjacent |

## 2. Liquid AI

### LFM2.5-1.2B-Instruct / 1.2B-Thinking / 350M

| Claim | Value | Source | Status |
|---|---|---|---|
| Identity | 1.17B, 16 layers (10 double-gated conv + 6 GQA), 32,768 ctx, vocab 65,536 | `https://huggingface.co/LiquidAI/LFM2.5-1.2B-Instruct`; `.../LFM2.5-1.2B-Thinking` | official |
| Instruct sampling | T 0.1, top_k 50, rep 1.05 | Instruct card | official |
| Thinking sampling | T 0.05, top_k 50, rep 1.05 | Thinking card | official |
| 350M / 230M sampling | T 0.1, top_k 50, rep 1.05 | `https://github.com/vllm-project/recipes/blob/main/LiquidAI/LFM2.5.md` (recipes table attributed to the cards) | official-adjacent |
| Thinking switch | LFM2.5-1.2B-Thinking: thinking **on, cannot switch off, no thinking levels** | `https://docs.liquid.ai/lfm/key-concepts/thinking` | official |
| Instruct has no reasoning mode | yes — the docs point at Instruct/350M/230M for direct answers | same page | official |
| Context cap | 32K, i.e. advertised but not the 262K of Qwen3.5 | cards | official |
| Known caveat | Liquid advises leaving `max_tokens` unset for reasoning models because a cap truncates the chain of thought | Liquid vLLM recipe | official-adjacent |
| GGUF | QAD Q4_0 used for 1.2B/350M in rounds 1–2 (vendor QAD checkpoint); Thinking Q4_K_M 731 MB | vendor GGUF repos | official |

**Do not treat Instruct and Thinking as interchangeable.** They are separate
checkpoints with different temperatures (0.1 vs 0.05), different defaults and
different failure modes; round 3 tests them as separate candidates.

### LFM3-1.7B (added during research)

| Claim | Value | Source | Status |
|---|---|---|---|
| Identity | 1.73B, tri-hybrid (13 short conv + 5 sliding-window attn + 6 global attn), 256K ctx | `https://huggingface.co/LiquidAI/LFM3-1.7B` | official |
| Reasoning | on by default, ChatML ` thinking`; `enable_thinking=False` available but the card warns it can be inconsistent | card + `https://goremy.ai/blog/articles/lfm3-1-7b-liquid-ai-local-agent-model` | official / community |
| Sampling general | T 0.1, top_k 50, rep 1.1 | card | official |
| Sampling math/STEM | T 1.0, top_p 0.95, top_k 20 | card | official |
| Runtime requirement | llama.cpp b11960 (v0.8.0)+ — this machine has build 12015 | card + `https://www.modelscope.cn/models/LiquidAI/LFM3-1.7B-GGUF` | official |
| Stated limits | not recommended for complex multi-step reasoning; positioned for structured output/extraction | card | official |
| Size | Q4_K_M 1.10 GB | card / HF | official |

## 3. HuggingFace SmolLM

### SmolLM3-3B

| Claim | Value | Source | Status |
|---|---|---|---|
| Identity | 3B, GQA + NoPE (3:1), 64K trained, 128K with YaRN | `https://huggingface.co/HuggingFaceTB/SmolLM3-3B` | official |
| Thinking | dual mode, `/think` and `/no_think` in the **system prompt**; template `enable_thinking` also controls it | card + `chat_template.jinja` | official |
| Sampling | T 0.6, top_p 0.95 | card | official |
| Effort levels | none; extended thinking is binary | card | official |
| GGUF | `ggml-org/SmolLM3-3B-GGUF` Q4_K_M 1.92 GB; thinking requires `--jinja` | GGUF card | official-adjacent |
| Known issue | long reasoning traces are the default; no budget control | card | official |

### SmolLM2-1.7B-Instruct
Non-thinking control, retained from round 2 (4.67/12 under the shared
configuration). No reasoning mode. `https://huggingface.co/HuggingFaceTB/SmolLM2-1.7B-Instruct`.

## 4. IBM Granite 4.0

| Claim | Value | Source | Status |
|---|---|---|---|
| Variants | H-Small 32B/9B active; H-Tiny 7B/1B active MoE; H-Micro 3B dense hybrid; H-1B; H-350M; plus non-hybrid Micro/1B/350M for runtimes without Mamba2 | `https://www.ibm.com/granite/docs/models/granite4-0` | official |
| Reasoning | **no thinking/reasoning mode documented** for 4.0 Nano/H-350M/H-Micro; these are instruct models | model cards (`ibm-granite/granite-4.0-h-350m`) | official |
| Consequence | Granite cannot enter the thinking matrix. Round‑2's `granite-4.0-350m` (0.33/12, collapsed output) remains a non-thinking historical control only | round‑2 record | observed |
| GGUF | `ibm-granite/granite-4.0-h-350m-GGUF` Q4_K_M 223 MB; H-Micro Q4_K_M 1.94 GB (not downloaded — no reasoning mode to test) | HF API | official |

## 5. Historical controls (not re-tuned)

| Model | Why retained | Round‑2 result |
|---|---|---|
| FunctionGemma-270M-it | function-calling specialist, not a text extractor | 0/12, token loop |
| Gemma-3-270M-it | smallest Gemma, no reasoning mode | 0/12, no decision column |
| Granite-4.0-350M | non-reasoning, collapsed output | 0.33/12 |
| Qwen2.5-0.5B / 1.5B | pre-reasoning generation | 2.0 / 4.0 of 12 |

These are reported for continuity only. Per the brief, no further compute was
spent forcing them onto a task their training does not target.

## 6. Sampler-chain facts that affect fairness

- llama.cpp's default `min_p` is **0.05**, while every publisher preset above
  states `min_p = 0.0`. Round 3 therefore sends `min_p` explicitly on every
  request. Leaving it unset would have silently made the "official" preset
  stricter than the publisher's.
- `repeat_penalty` is applied by llama.cpp over the last 64 tokens by default.
  Qwen recommends 1.0 (i.e. off); Liquid recommends 1.05 (LFM2.5) or 1.1
  (LFM3). These are sent per request and recorded.
- Qwen3.5's hosted `thinking_budget` is **not** the same control as llama.cpp's
  `--reasoning-budget`. The runtime flag forces the end-of-thinking token when
  the budget is hit; it is a harness-side cap. Results from it are labelled
  "runtime-imposed budget", never "model effort level".
- Identical seed values across families do not imply identical sampling: the
  vocabularies, templates and sampler chains differ. Seeds are recorded for
  reproducibility, not treated as a controlled variable.

## 7. Addendum: TinyLlama-1.1B-Chat-v1.0 (added 2026-11-19, after scoring)

- **Latest official release is still v1.0.** The TinyLlama paper's "v1.1" is a
  *pretraining* checkpoint (`TinyLlama-1.1B-intermediate-*`); no v1.1 chat model
  and no TinyLlama-2 exist. A Hugging Face search for `TinyLlama-1.1B-Chat-v1.1`
  returns nothing. Source: `https://huggingface.co/TinyLlama/TinyLlama-1.1B-Chat-v1.0`,
  `https://arxiv.org/abs/2401.02385` (§2.5), retrieved 2026-11-19.
- **Architecture and limits:** Llama-2 architecture and tokenizer, 1.1B params,
  BF16 originals, Apache-2.0, **2,048-token context**. Trained 3T tokens, then
  SFT on UltraChat and DPO on UltraFeedback following the Zephyr recipe. No
  reasoning mode, no thinking switch, no effort control. Source: model card,
  retrieved 2026-11-19.
- **Publisher sampling** (the card's own `pipeline` example): `do_sample=True`,
  `temperature=0.7`, `top_k=50`, `top_p=0.95`. The card sets no repetition
  penalty, so round 3 sent `repeat_penalty=1.0` rather than inventing one.
- **GGUF:** `TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF`,
  `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`, 668,788,096 bytes,
  SHA-256 `9fecc3b3cd76bba89d504f29b616eedf7da85b96540e490ca5824d3f7d2776a0`.
  No first-party GGUF repo exists under the TinyLlama org. Community
  re-quantization, labelled as such.
- **Note on `ggml-org-dev/tiny-llamas`:** despite the name, that repo holds
  Karpathy's `stories15M`/`stories260K` toy models, not TinyLlama. It was
  checked and rejected as a source.

## 8. Addendum: abliterated / uncensored variants (2026-11-19)

All four were added on request after the main round. Each run reuses the exact
sampler chain recorded for its baseline, so configuration is not a confound.
Retrieved 2026-11-19.

| Variant | Repository | Claim on the card | Verified? |
|---|---|---|---|
| Qwen3.5-4B Uncensored, HauhauCS Aggressive | `HauhauCS/Qwen3.5-4B-Uncensored-HauhauCS-Aggressive` | "0/465 refusals … zero capability loss"; recommends Qwen's own presets (thinking T0.6/top_p0.95/top_k20; non-thinking T0.7/top_p0.8/top_k20), and says to keep ≥128K context for thinking | card read; the capability claim is tested here at n=2 corpus draws, not accepted on trust |
| Qwen3.5-2B Uncensored, HauhauCS Aggressive | `HauhauCS/Qwen3.5-2B-Uncensored-HauhauCS-Aggressive` | same family claim | card not separately read; same lineage |
| LFM2.5-1.2B-Instruct abliterated | `mradermacher/Huihui-LFM2.5-1.2B-Instruct-abliterated-GGUF` (weights: huihui-ai) | no method detail verified beyond repo metadata | **not verified** — labelled community abliteration |
| LFM3-1.7B Heretic abliterated | `Vortecks/LFM3-1.7B-Heretic-Abliterated-GGUF` | no method detail verified beyond repo metadata | **not verified** — labelled community abliteration |

Files used (Q4_K_M, SHA-256):

- `abl-4b-q4km.gguf` — 2,707,513,696 B — `79e28ecacf84e75b6056cf4059636d435aa9eb67795780f7b7dbc7d32a962741`
- `abl-2b-q4km.gguf` — 1,270,808,032 B — `be3ccca13a9d1bc8b67165ea80bd103cd6151a37ad131183cc7ca388f78d9517`
- `abl-lfm120-q4km.gguf` — 730,896,000 B — `11731cb404b6267b63dac25710379425e2c1f42b41d0f24b7002fa6f5175530a`
- `abl-lfm3-q4km.gguf` — 1,099,627,040 B — `b145feb7fcc28fd0b4d0c2ef10a0f9ba4cbde17920ef2eb717afc8e8c6209e25`

"Uncensored" and "abliterated" are vendor/community terms, not standard
measurements. Nothing in this round verifies refusal removal directly; the
corpus contains no request a stock model would refuse, so the variants were
tested only for side effects on formatting, attribution and judgement.
