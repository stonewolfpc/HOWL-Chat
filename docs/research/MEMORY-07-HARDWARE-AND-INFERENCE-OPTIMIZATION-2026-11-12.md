# HOWL Chat — Hardware and Inference Optimization: Research Assessment

**Status:** Research and option analysis only. **No implementation authorization. No benchmark executed. No model loaded.**
**Date:** 2026-11-12
**Why this document exists:** Owner hypothesis `H1` cannot be evaluated without knowing what the hardware can actually host, and the work order requires removing any implication that 3–14B defines Howl's model envelope. This document records verified inference-engine behaviour, the resource arithmetic, offload and quantization options, the reserve policy, and the measurement plan. It is separate from `MEMORY-01` because it is an inference-engineering domain with its own verified upstream facts and its own failure modes.
**Companions:** `MEMORY-01-RESEARCH-SYNTHESIS-2026-11-12.md` (rev 2), `MEMORY-03-EVALUATION-PROTOCOL-2026-11-12.md` (§4, `D-RES-1`), `MEMORY-02-SOURCE-REGISTER-2026-11-12.md` (S28–S32, S43–S45).
**Controlling records:** `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`, `docs/architecture/AMENDMENT-01-UNIFIED-HARDWARE-INTELLIGENCE-2026-11-11.md`.

**Authority:** Research evidence. It does not authorize changing runtime settings, loading models, or shipping resource policy.

---

## 1. Verified upstream facts

These were checked against primary documentation. They correct several assumptions that are common in local-AI projects, including some implied by Howl's current settings.

| # | Fact | Source | Why it matters |
|---|---|---|---|
| 1 | The `--fit` path uses a **fixed `--fit-target` default of 1024 MiB**, not a percentage rule | S32 | A "10% VRAM" policy is a **Howl policy**, not an engine default. Conflating them produces wrong estimates on small and large GPUs alike. |
| 2 | `--fit off` disables the automatic fitting behaviour | S28 | Howl already passes `--fit off` in one path; automatic placement is then the app's responsibility |
| 3 | Context shift is **off by default** since PR #15416; `--no-context-shift` exists | S28, S31 | A design that assumes shift happens is measuring nothing |
| 4 | Quantized V-cache **requires Flash Attention** | S32 | `kv_cache_type` and `flash_attention` are coupled, not independent settings |
| 5 | KV-cache quantization is a memory/speed trade with a quality cost that must be measured per model | S32, S44–S46 | Cannot be assumed free |
| 6 | Multi-GPU supports row-split and layer-split strategies, with different balance characteristics | S31 | "Multi-GPU" is not one configuration |
| 7 | MoE models can offload expert tensors to CPU while keeping attention on GPU (`--n-cpu-moe`, `-ot` overrides) | S28, S31 | A large MoE can run on modest VRAM; a dense model of the same parameter count cannot |
| 8 | The server accepts `-b`/`-ub` batch sizes; prefill memory scales with them | S28 | Batch size is a memory knob, not only a speed knob |
| 9 | Multiple parallel slots **share** the context budget; per-slot context is `-c` divided by slots | S28 | A "20,480 context" claim can mean ~20,480 total, not per conversation |
| 10 | Effective context is often well below the claimed window | S15 | Large `-c` is not memory; it is headroom |
| 11 | Flash Attention changes memory and speed and is required for quantized V | S41, S32 | Must be a recorded, not incidental, setting |
| 12 | OOM prevention cannot be guaranteed in general: allocation depends on driver, other processes, fragmentation, and prompt length | S28, S31 | **No universal OOM-prevention claim is admissible** |

---

## 2. Howl's current configuration, and what is off about it

Source-inspected from `config/model-settings.json` and `app-llama.go`:

| Setting | Current value | Assessment |
|---|---|---|
| `context_size` | 20,480 | Large; effective context will be lower (S15) |
| `gpu_layers` | 22 | Partial offload; the split must be recorded per model, not fixed |
| `kv_cache_size` | 512 | **Ambiguous unit** — 512 MiB, 512 tokens, or a slot count. An ambiguous memory setting cannot be reasoned about; it must be defined before any benchmark |
| `kv_cache_type` | `flash-attention` | **Contradicts** `flash_attention: false` (below) |
| `flash_attention` | `false` | With quantized V-cache this combination is invalid (fact 4); with `kv_cache_type: flash-attention` it is internally inconsistent |
| `offload_kqv` | `true` | Offloads KQV to GPU; interacts with the KV type |
| `context_rollover_mode` | `memory-guided` | Exposed in settings; not implemented in the chat path (`MEMORY-05`) |
| `custom_jinja_template` | unset | Falls back to the server's template selection (`MEMORY-06`) |
| `--fit` | passed as `off` in one path | Automatic placement disabled; the app owns the split |

**Finding (source-level, not runtime):** the KV/attention settings are mutually inconsistent and the KV size unit is undefined. Any benchmark run on this configuration would produce numbers that cannot be attributed to a specific memory policy. That is a **measurement blocker**, not merely a cleanup item.

**Recommendation class:** this document recommends only that the configuration be *defined and recorded* before execution. It does not recommend specific values, because that requires measurement.

---

## 3. Resource arithmetic (calculated, not measured)

### 3.1 Weights

```text
weight_bytes ≈ parameters × bits_per_weight ÷ 8        (plus per-tensor overhead)
```

A K-quant is not a uniform bit rate; block scales and mins add overhead. Approximate effective bits per weight for common GGUF K-quants are **not** the nominal number in the filename. Any manifest must record the **actual file size**, not a computed estimate.

### 3.2 KV cache

```text
kv_bytes_per_token = 2 × n_layers × n_kv_heads × head_dim × bytes_per_element
kv_total           = kv_bytes_per_token × context_tokens × slots
```

- The factor 2 is K and V.
- `n_kv_heads` is the **KV** head count (GQA models have far fewer than attention heads).
- `bytes_per_element` is 2 for f16, 1 for q8_0, ~0.5 for q4_0.
- **Slots multiply the total.** This is the most commonly missed term.

**This is `Calculated`, and it must be checked against the engine's reported allocation.** A calculation that disagrees with observed VRAM by a large margin usually means a wrong head count, a wrong slot count, or a compute buffer not accounted for.

### 3.3 Compute buffers and fragmentation

Prefill buffers scale with batch size, micro-batch, and sequence length. Fragmentation and driver reservation consume VRAM that no formula predicts. **The only trustworthy number is post-load free VRAM measured on the actual machine.**

### 3.4 The reserve policy

| Element | Rule |
|---|---|
| Default on `HW-1` | **10% of total VRAM free after model load**, as an intended default |
| On other hardware | A declared policy with a percentage, a floor, and a ceiling, recorded per machine |
| Not an engine default | Upstream `--fit-target` is 1024 MiB (fact 1); the two must never be conflated in a report |
| User override | Allowed; recorded as a configuration, not a violation |
| Claim limit | **No universal OOM prevention.** The system may refuse to load when a pre-flight estimate exceeds the policy; it may not claim to prevent OOM in general |
| Measurement | Post-load free VRAM recorded on every scored run; a breached run is flagged, not averaged (`D-RES-1`) |

**Why 10% is defensible as a starting default:** it leaves room for compute buffers, fragmentation, driver overhead, and a longer-than-expected prompt without pretending to be a guarantee. It is a policy chosen for margin, not a derived constant. It should be revisited once `HW-1` measurements exist.

---

## 4. Placement and offload options

Each option is recorded against the ten evidence questions in compressed form.

| Option | What it does | Solves | Limitations | Cost | Complements | Conflicts with | Howl-compatible | Architectural change | Status |
|---|---|---|---|---|---|---|---|---|---|
| **Full GPU offload** | All layers on GPU | Speed; simplest KV behaviour | Requires weights + KV + compute to fit | VRAM | Flash Attention | Large models on small VRAM | Yes | No | **Established** |
| **Partial layer offload** (`-ngl`) | Some layers on GPU, rest on CPU | Runs models larger than VRAM | Per-layer transfer cost; slower prefill and decode; balance is model-specific | Speed | Any model size | Latency targets | Yes (current `gpu_layers: 22`) | No | **Established** |
| **CPU-only** | No GPU layers | Portability; `HW-2` | Very slow for large models | Time | Small models | Interactive latency | Yes | No | **Established** |
| **MoE expert offload** (`--n-cpu-moe`, `-ot`) | Experts on CPU, attention on GPU | Large MoE on modest VRAM | Speed depends on expert routing; per-model tuning | Time | MoE families | Dense-model assumptions | Yes, model-dependent | No | **Strongly supported** |
| **Multi-GPU row split** | Tensors split by row across GPUs | Uses multiple GPUs for one model | Inter-GPU traffic; balance issues | Bandwidth | Small-to-mid models | Single-GPU simplicity | Yes | No | **Established** |
| **Multi-GPU layer split** | Layers distributed across GPUs | Large models across GPUs | Pipeline bubbles; unbalanced splits | Bandwidth | Large models | Low-latency single-stream | Yes | No | **Established** |
| **KV-cache quantization** | f16 → q8_0/q4_0 | Frees VRAM for longer context | Quality cost; **requires Flash Attention** | Quality | Larger context | Flash Attention off | Yes | No | **Conditional** (must measure quality) |
| **Flash Attention** | Fused attention kernel | Lower memory, faster attention | Backend/GPU support varies | Compatibility | Quantized V; long context | Some older GPUs/backends | Yes | No | **Established** |
| **mmap / no-mmap** | Map weights from disk vs load | Faster load; lower peak RSS with mmap | mmap can hurt if the file is paged out | Load time | Large models | Predictable latency | Yes | No | **Established** |
| **Model unload/reload** | Free one model before loading another | Avoids two-model OOM | Load latency; cache loss | Time | Single-GPU operation | Fast switching | Yes | Yes (lifecycle policy) | **Established** |
| **Adaptive allocation** | Choose offload/quant/context from detected hardware | Portability across machines | Estimation error; needs measurement feedback | Complexity | Preflight | Deterministic reproducibility | Conceptually | Yes | **Promising, unproven** |
| **Speculative decoding** | Draft model proposes tokens | Speed | Extra VRAM for the draft model; acceptance rate model-dependent | VRAM + complexity | Speed targets | Tight VRAM | Backend-dependent | Yes | **Conditional** |

---

## 5. Pre-flight and OOM behaviour

| Step | What it does | Failure if omitted |
|---|---|---|
| Read hardware facts | VRAM, free VRAM, driver, RAM, page file | Policy applied to the wrong baseline |
| Estimate weights from the **actual file size** | No filename-based guessing | Underestimate → OOM at load |
| Estimate KV from the formula, with slot count | Includes GQA head count | Underestimate → OOM at long prompt |
| Add compute-buffer allowance | Declared, hardware-specific | Underestimate → OOM during prefill |
| Apply the reserve policy | Refuse or reduce offload/context if the result exceeds it | Reserve breached silently |
| Record the decision | Manifest entry | Unreproducible run |
| Measure after load | Actual free VRAM | Policy never validated |

**Documented OOM modes:**

1. **Load-time OOM** — weights + KV reservation exceed available VRAM.
2. **Prefill OOM** — a long prompt creates a compute spike; the model loaded fine and fails on the first big turn.
3. **Slot multiplication** — parallel slots each reserve context; N slots multiply KV.
4. **Two-model OOM** — a companion or embedding model loaded alongside the chat model.
5. **Fragmentation** — repeated load/unload leaves unusable VRAM; a load that worked before now fails.
6. **Driver reset / TDR** — a long kernel triggers a driver timeout; the process dies with no clean error.
7. **Silent spill** — the driver pages VRAM to system memory; the model "runs" at a fraction of the speed.

**Howl must report these distinctly.** "The model failed" is not a diagnosis.

---

## 6. Meeting points with the memory design

| Memory design element | Hardware consequence |
|---|---|
| Structured state card in every prompt | Raises steady-state prompt tokens; prefill cost per turn |
| Evidence retrieval at 1k–6k tokens | Raises prefill; the 8B cliff is around 3k (S1) |
| Summaries and consolidation | Extra model calls; may need a second model resident, or serialized calls |
| Embedding index | Either a resident embedding model (VRAM) or a CPU embedding path (latency) |
| Reranker | CPU cost (S13) or VRAM |
| Local lorebook activation | Token cost; no VRAM |
| Context flushing | Prefill spikes at reconstruction; cache invalidation |

**Therefore:** the memory architecture cannot be evaluated independently of a declared residency plan (which models are resident, which are loaded on demand, and what the reserve is). That is why `MEMORY-03` requires the hardware manifest before any run.

---

## 7. Measurement plan

| Measure | Method | Notes |
|---|---|---|
| `post_load_free_vram` | Engine/driver report after load | Compared to the declared reserve |
| `peak_vram` | Sampled during prefill and decode | Prefill is usually the peak |
| `tokens_per_sec` (prefill, decode) | Timed runs, fixed prompt/output lengths | Separate the two; never report one number |
| `load_time` | Cold load, and reload after unload | Cold vs warm (page cache) |
| `peak_rss` | Process and system | Detects spill and mmap behaviour |
| `oom_events` | Counted and classified per §5 | Never merged into a single "failed" count |
| `swap_activity` | Page-file counters during decode | Detects silent spill |
| `cache_hit_rate` | Server-reported prompt reuse where available | Otherwise declared unmeasured |
| `quality_at_quant` | Task accuracy at the declared quant | A memory condition is only valid at a pinned quant |

All of the above are `Unestablished` for Howl until executed.

---

## 8. Failure modes and mitigations

| Failure | Cause | Mitigation | Residual risk |
|---|---|---|---|
| OOM at load | Underestimated weights + KV | Actual file size + formula + compute allowance | Medium |
| OOM at long prompt | Prefill spike | Reduce batch/micro-batch, cap prompt, lower context | Medium |
| Reserve breach | Policy not applied or not measured | Pre-flight + post-load check | Low if enforced |
| Silent spill | VRAM exhausted into system memory | Swap activity monitoring; refuse rather than spill | Medium |
| Two-model OOM | Companion/embedding model resident | Explicit residency plan | Low |
| Fragmentation | Repeated load/unload | Unload-before-load; restart hint | Medium |
| Driver reset | Long kernel | Batch limits; watchdog awareness | Medium |
| Quantized-V invalid config | Flash Attention off | Refuse the combination | Low if validated |
| Non-reproducible placement | Auto-fit with no record | Record the resolved split | Low |

---

## 9. Unresolved questions

| ID | Question | Status |
|---|---|---|
| `Q-HW1` | What is the actual VRAM of `HW-1`, and what is free after load at each sampled quant? | Unmeasured; required before any benchmark |
| `Q-HW2` | What does `kv_cache_size: 512` mean, and what is the correct unit? | Undefined in settings; must be resolved |
| `Q-HW3` | Is `kv_cache_type: flash-attention` with `flash_attention: false` reachable, and what does the engine do? | Inconsistent; must be tested |
| `Q-HW4` | Is the 10% reserve right for `HW-1`, or should it be a floor in MiB? | Depends on `Q-HW1` |
| `Q-HW5` | What context size is actually usable at each quant before quality degrades? | Unmeasured (S15 says claimed ≠ effective) |
| `Q-HW6` | Does a multi-slot configuration multiply KV as expected on this backend version? | Documented; needs verification |
| `Q-HW7` | What is the prefill cost of the proposed memory package at 1k/3k/6k tokens? | Unmeasured |
| `Q-HW8` | Can an embedding model and a chat model be co-resident under the reserve policy? | Unmeasured |
| `Q-HW9` | Does MoE expert offload make large MoE models practical on `HW-1`? | Untested |
| `Q-HW10` | What is the cost of repeated consolidation calls in VRAM/time? | Unmeasured |
| `Q-HW11` | Is there a hardware configuration below which Howl should warn rather than refuse? | Product policy, undefined |
| `Q-HW12` | How should multi-GPU be detected and balanced without user tuning? | Open |

---

## 10. What this document does not do

- It does not authorize loading a model, changing settings, or running a benchmark.
- It does not change `context_size`, `gpu_layers`, KV type, or Flash Attention.
- It does not declare any model unsupported because of its size.
- It does not claim OOM prevention is achievable.
- It recommends no specific quant, offload split, or reserve value — those require measurement.

**The one firm conclusion: Howl's current KV/attention settings are internally inconsistent and the KV size unit is undefined, so no hardware or memory benchmark can be attributed to a known configuration until that is fixed.**
