# Inference Memory and Offload Capability Survey — 2026-11-24

**Status:** research complete. No production behaviour changed by this document.
**Question:** beyond KV-cache sizing and Flash Attention, what memory techniques exist for
running inference on low-VRAM machines, which ones do our shipped engines actually support,
and which are appropriate for Howl?

**Design constraints carried from the owner:**

- Target more than this one machine. The only available test box is an 8 GB NVIDIA card, so
  nothing may be hard-coded to 8 GiB, CUDA, or this GPU generation.
- Real measurements only. No fabricated capacity, no simulated offload.
- Automatic features should activate wherever they can. If an automatic feature cannot be
  applied or cannot be verified, it must not lock the user out or block a model load — it
  must fall back to a manual choice.
- The existing engines are already hardware-aware. Leverage them rather than duplicating
  their placement logic.

---

## 0. Bottom line

1. **The engines already contain most of the low-VRAM machinery.** llama.cpp build 12015 has
   automatic fitting, MoE/FFN CPU offload, per-tensor placement, KV offload, quantized KV,
   load modes, lazy tensors, prompt caching and multi-GPU splitting. stable-diffusion.cpp
   `3b12781` has automatic graph-cut placement across GPU → RAM → other GPU → disk, per-module
   backends, weight-residency expiry, VAE tiling and on-load quantization.
2. **Howl currently uses only the oldest subset of it.** Chat sends a fixed `-ngl`, a fixed
   context, a fixed batch and an optional Flash Attention switch. That makes llama.cpp's
   `--fit` sizing pass inert. Three settings the UI presents (`kv_cache_type`,
   `kv_cache_size`, `offload_kqv`) never reach the engine at all.
3. **The highest-value wins are not exotic.** They are: let layers be placed from measured
   memory instead of a fixed number; use MoE/dense FFN offload when the architecture supports
   it; quantize KV only when it is needed; and report the placement the engine actually chose.
4. **"Never block" requires a policy change.** Today the admission authority refuses to start
   a model when free VRAM cannot be measured or the request does not fit. That is fail-closed.
   The owner's rule is fail-open-to-manual: warn with the measured numbers, offer the load,
   and require an explicit choice. That is not a silent bypass; it is an explicit one.
5. **Two reserve policies currently coexist.** Chat admission keeps 10 % of *free* VRAM.
   Image admission always keeps a fixed 1 GiB (mirroring sd-server's `--max-vram -1`) and
   marks that reserve explicit, so the 10 % rule never applies to image loads. This is a
   deliberate engine-compatibility choice but it is not the agreed default; it needs an
   explicit decision rather than accidental drift.

---

## 1. What was inspected

| Component | Version / pin | Evidence |
|---|---|---|
| llama.cpp server | `version: 0.8.0-dev (build 12015, commit a34de8618)` | `runtime/llama/llama-server.exe --version`; full help captured at `build/tmp/research/llama-help-current.txt` (756 lines) |
| stable-diffusion.cpp server | `master-1122-3b12781`, commit `3b12781` | `sd-server.exe --help`; captured at `build/tmp/research/sd-help-current.txt` (232 lines); pinned in `internal/image/engine/sdcpp/version.go:19-20` |
| Go runtime | `go 1.23.8` | `go.mod` |
| Integration surface | CLI only | Howl starts both engines as child processes and talks HTTP. There are no CGo/C bindings, so **any flag the shipped binary accepts is available to Howl**, and unsupported flags are the only compatibility limit. |
| Whisper | not installed on this machine | `bin/whisper/whisper-cli.exe` absent; audio placement claims below are marked unverified |
| TTS | Piper / ONNX | small CPU-resident models; no GPU placement flags in Howl's path |

The image engine already probes the sd binary's own `--help` and keeps a per-commit capability
set (`internal/image/engine/sdcpp/capability.go`, `ProbeCLI`). llama.cpp has no equivalent
probe; `ValidateLlamaServer` only runs `--version`. That asymmetry matters for future flags.

---

## 2. What Howl sends today

### 2.1 Chat / D&D — llama-server

`buildLlamaServerArgs` (`app-llama.go`) sends:

| Flag | Source | Consequence |
|---|---|---|
| `-m`, `--mmproj` | selected model / verified projector | fixed |
| `-c` context | saved `context_size` (capped to model window in `applySafeDefaults`) | **always explicit** → `--fit` cannot resize context |
| `--threads` | saved `threads` | CPU parallelism |
| `-b` batch | saved `batch_size` | **always explicit** → `--fit` cannot resize batch |
| `-ngl` | saved `gpu_layers` (default `0`) | **always explicit** → `--fit` cannot choose layers |
| `--fit off` | only for sliding-window models | deliberate workaround for `GGML_SCHED_MAX_SPLIT_INPUTS` |
| `--flash-attn` | only when the saved boolean is `true` | `false` sends nothing, which means `auto`, **not off** |
| `--tensor-split` | saved string | multi-GPU only |
| rope / template / reasoning flags | saved | unrelated to memory |

Not sent: `--fit-target`, `--fit-ctx`, `--n-cpu-moe`, `--cpu-moe`, `--moe-cache-mib`,
`--n-cpu-ffn`, `--override-tensor`, `--no-kv-offload`, `--cache-type-k/v`, `--ubatch-size`,
`--cache-ram`, `--cache-reuse`, `--load-mode`, `--lazy-mode`, `--swa-full`,
`--ctx-checkpoints`, `--no-op-offload`, `--parallel`, `--device`, `--main-gpu`,
`--split-mode`.

`--fit` is documented by the binary as *"whether to adjust unset arguments to fit in device
memory (default: on)"*. Because Howl sets every argument fit would adjust, the automatic
sizing pass has nothing left to do. The engine's hardware awareness is present but bypassed.

### 2.2 Image generation — sd-server

`Engine.buildServerArgs` (`internal/image/engine/sdcpp/engine.go`) sends:

- `--max-vram` (default `-1`: snapshot free memory, reserve 1 GiB)
- `--backend`, `--params-backend` (only when the user sets them)
- `--type` (on-load quantization), `--conditioning-cache-size`, `--mmap`
- `--weight-residency-ttl`, `--weight-staging-ttl` (opt-in; defaults disabled)
- `--vae-tiling`, `--temporal-tiling`, `--circular`
- `--flash-attn on` when the checkbox is on; otherwise nothing (`auto`)
- `--offload-to-cpu`, `--cache-mode` / `--cache-option`

Placement is left to sd-server's `--auto-fit`, which is on by default. The app reads the
resulting plan and per-component split out of the server log (`placement.go`) and exposes
`GetImagePlacement` / `GetImageDiagnostics`. This is the correct division already: **admission
decides policy, the engine decides placement, the app reports what actually happened.**

### 2.3 Audio and embeddings

- **ASR:** Whisper through `whisper-cli`, or Qwen3-ASR through a llama-server adapter. The
  adapter inherits llama.cpp's flags, but Howl does not pass any placement flags on that path.
- **TTS:** Piper (ONNX) or Kokoro/Melo; small CPU models. No VRAM placement flags.
- **Embeddings/RAG/memory:** use the already-loaded chat server's `/v1/embeddings`
  (`internal/backend/rag/embedder.go`). There is no separate embedding-model residency, which
  is the right default for a low-VRAM machine: it avoids a second model entirely.
- Audio workloads do pass through the same admission authority (`internal/audio/pipeline.go`),
  so they are budgeted, but their placement is not optimised beyond choosing CPU/GPU.

---

## 3. The list: techniques that exist

Legend for **Our build**: ✅ flag present in the shipped binary's help; ➖ absent/irrelevant;
**Howl fit**: verdict from section 5.

### 3.1 llama.cpp — placing weights

| Technique | Flag | What it does | Cost / limit | Our build | Howl fit |
|---|---|---|---|---|---|
| Full GPU offload | `-ngl all` / layer count | All transformer layers on GPU | Needs the whole model plus KV and compute buffers | ✅ | Automatic when it fits |
| Partial layer offload | `-ngl N` | First/last N layers on GPU, rest on CPU | Per-layer split; VRAM for N layers + KV | ✅ | Automatic fallback |
| Automatic fit | `--fit on\|off` (default on) | Adjusts **unset** args (layers, context) to device memory | Does nothing for args Howl already sets; can shrink context unless `-fitc` floors it | ✅ | **Adopt** — but must stop forcing the args it needs |
| Fit margin | `-fitt, --fit-target MiB` | Per-device headroom target; default 1024 MiB | Larger margin = fewer layers, safer desktop | ✅ | **Adopt** — natural carrier for the 10 % reserve |
| Fit context floor | `-fitc, --fit-ctx N` (default 4096) | Minimum context fit may choose | Prevents fit from collapsing context to nothing | ✅ | Adopt with fit |
| MoE expert offload, all | `-cmoe, --cpu-moe` | Keeps all MoE expert matrices in RAM | MoE only; attention + KV stay on GPU | ✅ | **Adopt** for MoE when VRAM-short |
| MoE expert offload, counted | `-ncmoe, --n-cpu-moe N` | Keeps experts of the first N layers in RAM | Manual N; fit stops moving layers once an override exists | ✅ | **Adopt**, auto-computed + manual |
| MoE hot-expert GPU cache | `--moe-cache-mib N`, `--moe-cache-cpu-share F` | Caches hot CPU-resident experts in VRAM; cache misses can run on CPU | Takes VRAM from layers; must be sized explicitly; new/less proven | ✅ | Adopt later, manual first |
| Dense FFN offload | `-ncffn, --n-cpu-ffn N` | Keeps dense feed-forward weights of first N layers in RAM | Dense only; manual N; less community data than MoE | ✅ | **Adopt**, auto-computed + manual |
| Per-tensor placement | `-ot, --override-tensor pattern=CPU,...` | Arbitrary tensor→buffer routing; the general form of the two flags above | Easy to misconfigure; expert-only knob | ✅ | Manual / advanced |
| KV cache offload | `-nkvo, --no-kv-offload` | Keeps the whole KV cache in RAM | Frees the most VRAM, costs the most speed | ✅ | Manual fallback |
| Multi-GPU split | `-sm layer\|row\|tensor`, `-ts`, `-mg` | Splits layers or rows across GPUs | Needs multiple GPUs; row split needs backend support | ✅ | Use when >1 GPU, not here |
| Device selection | `--device`, `--list-devices` | Restricts offload to named devices | Wrong name = no offload | ✅ | Adopt for multi-GPU/ignore-iGPU |

### 3.2 llama.cpp — cache, context and loading

| Technique | Flag | What it does | Cost / limit | Our build | Howl fit |
|---|---|---|---|---|---|
| Quantized KV cache | `-ctk, --cache-type-k`, `-ctv` | Stores K/V at q8_0/q4_0 etc. | q8_0 is near-lossless; q4 can hurt long-context retrieval. **Quantized V cache requires Flash Attention enabled**: build 12015 rejects it at context creation with `quantized V cache requires flash_attn to be enabled` | ✅ | **Adopt** as a measured fallback step, not a default; fall back to f16 with a warning when Flash Attention is explicitly off |
| SWA window cache | default | Window-sized cache for sliding-window models | Already the memory-saving default | ✅ | Keep default |
| Full SWA cache | `--swa-full` | Full-size SWA cache | More memory, faster | ✅ | **Reject** as automatic |
| Context checkpoints | `-ctxcp N` (default 32/slot), `-cms N` | Per-slot checkpoints for context reuse | Each checkpoint holds state | ✅ | Lower on low-RAM; manual |
| Host prompt cache | `-cram, --cache-ram N` (default 8192 MiB) | Caches prompt state in host RAM | **8 GiB of RAM by default, not counted by Howl's admission** | ✅ | Size down on low-RAM; expose |
| Prompt cache reuse | `--cache-reuse N` | Reuses KV chunks via shifting | Memory-neutral; big prompt-time win | ✅ | Adopt for D&D long prompts |
| Micro-batch | `-ub, --ubatch-size N` (default 512) | Physical batch; drives compute-buffer size | Lower = less peak VRAM, slower prompt | ✅ | **Adopt** as fallback step |
| Parallel slots | `-np, --parallel N` | Server slots; each needs KV | More slots = more KV | ✅ | Pin to 1 for single-user |
| Unified KV per slot | `--kv-unified-per-slot N` | Caps shared KV pool per slot | New; needs a context budget | ✅ | Advanced |
| Load mode | `-lm, --load-mode auto\|none\|mmap\|mlock\|mmap+mlock\|dio` | mmap (default), pin in RAM, DirectIO | mlock pins RAM and can starve the OS; dio needs support | ✅ | Keep `auto`; expose manual only |
| Lazy tensors | `-lzm, --lazy-mode on\|auto\|off` | Reads huge per-layer embedding tensors from disk on demand (mmap) | Niche; disk latency | ✅ | Keep `auto` |
| Host-op offload | `--no-op-offload` | Runs host tensor ops on CPU | Saves some VRAM, costs speed | ✅ | Manual only |
| Speculative decoding | `--model-draft`, `--spec-*` | Draft model proposes tokens | **Adds** a model and memory | ✅ | **Reject** for low-VRAM goals |

### 3.3 stable-diffusion.cpp — placing and releasing

| Technique | Flag | What it does | Cost / limit | Our build | Howl fit |
|---|---|---|---|---|---|
| Automatic placement | `--auto-fit on\|off` (default on) | Places modules GPU → RAM → other GPU → disk; 512 MiB device margin; RAM reserve max(2 GiB, 10 %) | Disabled by any explicit `--params-backend`; compute reserves are estimates | ✅ | **Adopt** — already relied on |
| VRAM budget | `--max-vram` (GiB, per device) | Caps managed weights + runner buffers; negative reserves | Budget only covers managed allocations, not driver contexts | ✅ | Adopt — already wired |
| RAM residency | `--params-backend cpu` / `--offload-to-cpu` | Weights in RAM, staged to GPU | Slower; `--offload-to-cpu` is `params-backend '*=cpu'` | ✅ | Adopt as fallback |
| Disk residency | `--params-backend disk` | Reloads weights from disk on demand | Very slow; never chosen implicitly; last resort | ✅ | **Manual last resort** — prevents hard lockout |
| Per-module placement | `--backend te=cpu,vae=cpu,diffusion=cuda0`; `--params-backend` | Put text encoders/VAE on CPU, diffusion on GPU | Some families have bugs with CPU text encoders (SD3.5 noise) | ✅ | Manual, family-gated |
| On-load quantization | `--type q8_0/q5/q4`; `--tensor-type-rules` | Converts weights while loading; q8_0 measured at 0.62× file size | Quality loss at low types; tensor rules are expert-only | ✅ | Adopt for `--type`; manual for rules |
| VAE tiling | `--vae-tiling`, `--vae-tile-size`, overlap | Decodes in 256×256 tiles | Slower; required for large/Flux decodes on 8 GB | ✅ | **Adopt automatically under pressure** |
| Cheap previews | `--taesd`, `--taesd-preview-only` | Tiny autoencoder for step previews, real VAE for final | Preview quality only | ✅ | Already used when present |
| Idle release | `--weight-residency-ttl`, `--weight-staging-ttl` (seconds, -1 off) | Unloads components / GPU staging copies after idle | Reload latency on next job | ✅ | Adopt with a measured idle policy; never a fixed arbitrary delay |
| No prefetch | `--disable-prefetch` | Disables next-segment weight prefetch | Less overlap memory, slower | ✅ | Manual |
| Segmentation | default on; `--disable-segmented-compute` | Graph cutting to fit segments | Disabling uses more memory | ✅ | **Reject** as automatic |
| Eager load | `--eager-load` | Loads all params up front | More peak memory | ✅ | **Reject** as automatic |
| Flash / Sage attention | `--flash-attn` per module; `--sage-attn` (CUDA) | Attention kernel selection | Sage needs a CUDA build and supported shapes | ✅ | Keep `auto`; manual Sage |
| Direct convolution | `--conv-direct auto\|on\|off` | Implicit-GEMM path per module | Backend-dependent; not a memory guarantee | ✅ | Keep `auto` |
| Resolution / steps / batch | generation options | Dominant driver of activation memory | Not a flag trick; must be in the budget | ✅ | **Adopt** in preflight |

### 3.4 Audio and embeddings

| Workload | Memory reality | Techniques available | Howl fit |
|---|---|---|---|
| Qwen3-ASR via llama-server | a full llama.cpp model plus audio projector | all llama.cpp flags above, including `-ngl`, KV type, fit | reuse the chat plan; do not add a manager |
| Whisper via `whisper-cli` | model is small by comparison; CPU is often fine | whisper.cpp exposes GPU/thread controls, not layer offload | select by measured model size; no new tricks |
| Piper TTS | ONNX on CPU, tens of MB | none needed | keep CPU; never load for VRAM |
| Embeddings / RAG / memory | currently the chat server's embedding endpoint | reuse resident chat model; optional separate embedding model would add residency | keep on the chat server; CPU fallback = chat CPU fallback |

---

## 4. What the field reports actually say

These are community reports, not measurements on this machine. They justify trying a
technique; they do not qualify it for this hardware.

- **MoE expert offload is the biggest small-VRAM lever for MoE models.** A 75.8 GB IQ3_XXS
  model ran on 16 GB VRAM + 32 GB RAM with `-ngl 99 -ncmoe 42`; `-ncmoe 44` was slower,
  `41` OOM'd at 64K context, and `--cpu-moe` was slowest. The same report credits
  `--cache-reuse 256` with reusing 1.08 M prompt tokens.
  <https://github.com/ggml-org/llama.cpp/discussions/32762>
- **The MoE expert cache can nearly double decode** versus a same-path CPU-expert control,
  but it must be sized by the caller; `--fit` counts the cache when it does the sizing.
  <https://github.com/ggml-org/llama.cpp/discussions/24528>,
  <https://inferprice.pages.dev/articles/llama-cpp-moe-expert-cache>
- **Dense models have a sibling lever.** `--n-cpu-ffn N` keeps dense FFN weights in RAM;
  `--no-kv-offload` frees the most VRAM but costs the most speed; `-fitt 2048` reserves more
  desktop headroom because fit measures free memory only once at startup.
  <https://www.sector88.co/blog/llama-cpp-cuda-out-of-memory-fix>
- **`--fit` only adjusts what you leave unset.** Passing explicit values opts out of the
  sizing pass. <https://pipelinemonk.com/local-llms/run-big-model-small-gpu-llama-cpp>
- **WDDM can invalidate an on-paper fit.** A 262K-context q8/q4 KV + 6 GiB expert pool fit on
  paper and looked fine for a 5-token run, then collapsed from 5.5 to 0.71 tok/s on a
  250-token run. Peak VRAM must be proven at real workload, not at a toy prompt.
  <https://github.com/ggml-org/llama.cpp/discussions/28363>
- **sd.cpp auto-fit is a real placement engine**, with an explicit priority ladder and
  margins, and `--params-backend disk` as the genuine last resort. It does not switch a
  module to CPU *compute* merely because its weights exceed VRAM.
  <https://github.com/leejet/stable-diffusion.cpp/blob/master/docs/backend.md>
- **VAE tiling is not optional at 8 GB for heavy models.** A Flux.1-schnell 8 GB Vulkan
  walkthrough reports VAE decode OOM without `--vae-tiling`, with T5/VAE on CPU and the
  diffusion model on GPU. <https://discuss.huggingface.co/t/flux-1-schnell-on-8gb-vram-amd-no-cuda-the-gguf-format-mismatch-that-wastes-hours/177010>
- **Module-level CPU placement can be broken for a family.** SD3.5 Large produced pure noise
  with `--clip-on-cpu` because T5 compute buffers stayed ~12 MB; the same flags worked for
  SD3.5 Medium. <https://github.com/leejet/stable-diffusion.cpp/issues/1578>
- **`--offload-to-cpu` has crashed on some Vulkan builds** while working on others.
  <https://github.com/leejet/stable-diffusion.cpp/issues/1483>
- **sd-server can use more VRAM than sd-cli** for the same model (reported ~2 GB more at
  1152×704). Howl runs the server, so its own overhead belongs in the budget.
  <https://github.com/leejet/stable-diffusion.cpp/issues/1293>

---

## 5. Assessment for Howl

The rule applied to every row: **automatic only when the engine can verify the placement and
the result is observable; otherwise manual; never a hard refusal.**

### 5.1 Adopt as automatic

| Item | Why it fits | Conditions |
|---|---|---|
| Automatic layer/context placement (llama `--fit`, or an explicit plan computed from the same numbers) | the machine-specific answer already exists in the binary; a fixed `-ngl` cannot be right for every card | exactly one planner owns the numbers; the app must not also force `-c/-b/-ngl` |
| Reserve expressed as a per-device margin (`-fitt`, `--max-vram` negative) | keeps the agreed 10 % free-VRAM rule and matches each engine's own budget model | reserve computed from measured free memory, not total |
| MoE expert offload (`-ncmoe`) for MoE architectures | attention + KV stay on GPU; largest win on small cards for modern MoE models | architecture detected from GGUF, not filename; N computed then clamped; fall back to layer offload |
| Dense FFN offload (`-ncffn`) for dense architectures | same idea where `-ncmoe` does not apply | keep behind a conservative threshold until measured |
| Micro-batch reduction (`-ub`) | cheap way to cut peak compute VRAM | only after layers/offload; never below a usable prompt speed without saying so |
| KV quantization (`-ctk/-ctv q8_0`) as a *fallback step* | near-lossless and large saving when context is the pressure | never silently on; report it; q4 only by manual choice |
| VAE tiling under pressure (sd) | prevents a decode OOM that would otherwise lose the image | already engine-retried on failure; surface that it happened |
| sd on-load `--type` when the file type does not fit | the repo records 0.62× for q8_0 on this family | never re-quantize a file that already fits |
| Idle release of image components (`--weight-residency-ttl`) | satisfies "release when no longer needed" without killing active work | measured idle window; opt-in; never while a job is queued/running |
| Actual-placement reporting (both engines) | "loaded" is not "on GPU"; the user must see the truth | parse engine logs; no inferred claim |

### 5.2 Adopt as manual / advanced

| Item | Why manual |
|---|---|
| `--params-backend disk` | genuine last resort, but minutes-per-image slow. It should exist so the user is never locked out, not be chosen automatically. |
| `--override-tensor` (`-ot`) | powerful and easy to misconfigure; expert surface. |
| `--no-kv-offload` (`-nkvo`) | largest VRAM saving, largest speed loss; a deliberate trade. |
| `--load-mode mlock` / `dio` | pinning RAM can destabilise the OS; DirectIO is hardware-dependent. |
| `--moe-cache-mib` | must be sized against a budget; community results are promising but new. |
| Module CPU compute (`te=cpu`, `vae=cpu`) | family-specific bugs exist; only after per-family verification. |
| `--disable-prefetch`, `--no-op-offload`, `--sage-attn` | narrow performance/memory trades, not defaults. |
| Full SWA cache (`--swa-full`) | costs memory; default is the memory-saving choice. |
| Lowering context / batch / resolution | visible quality or capability loss; user's call, with numbers. |

### 5.3 Reject for this mission

| Item | Reason |
|---|---|
| Speculative decoding as a memory technique | adds a draft model and memory; it is a speed feature. |
| `--eager-load`, `--disable-segmented-compute` | increase peak memory; opposite of the goal. |
| Automatic `mlock` | pins the model in RAM and can starve the OS. |
| Automatic disk residency | destroys interactivity; only an explicit escape hatch. |
| Fixed layer counts or fixed 1 GiB reserves for every machine | violates "designed for more than my system". |
| Fabricated capacity when telemetry fails | already forbidden by the mission and still present in fallback paths (`estimateVRAMFromGPUName`). |

---

## 6. Defects found in the current wiring

These are evidence-backed, not stylistic:

1. **`--fit` is inert in Howl.** `app-llama.go` always passes `-c`, `-b` and `-ngl`; the
   binary documents `--fit` as adjusting *unset* arguments. The auto-sizing pass Howl relies
   on cannot change anything.
2. **`ClampToHardware` has no caller.** `internal/backend/gguf/memory.go:120` implements the
   full reduce-layers → reduce-context → reduce-batch → CPU fallback ladder (including the SWA
   partial-offload crash rule), but `applySafeDefaults` deliberately does not call it and
   admission refuses instead. The automatic-fit implementation exists and is unused.
3. **Three UI controls are dead.** `kv_cache_type`, `kv_cache_size` and `offload_kqv` are
   saved by `UpdateRuntimeSettings` and never read by `buildLlamaServerArgs`. Worse,
   `kv_cache_type` offers `standard / flash-attention / paged`; llama.cpp has no "paged" KV
   cache, and Flash Attention is a separate switch. The control is not merely unwired; it is
   misleading.
4. **Unchecking Flash Attention does not turn it off.** Both paths only append the flag when
   the boolean is true; the engine default is `auto`. A user who unchecks it still gets
   Flash Attention when the engine selects it.
5. **The image estimate ignores the generation workload's compute buffers.** Admission sizes
   image loads from module weights plus fixed reserves; width, height, batch and steps are
   carried in the request but not used in `estimateImage`. sd-server's own capacity check
   still catches the failure, but the preflight can say "fits" for a resolution that will not.
6. **The fixed image reserve bypasses the 10 % rule.** `imageVRAMReserveMB(nil)` returns
   1024 MB with `ReserveExplicit: true`, so the 10 % default never applies to image loads.
   On 8 GB that is >10 %, on 24 GB it is <10 %.
7. **llama-server's host RAM cache is not budgeted.** `--cache-ram` defaults to 8192 MiB and
   is never passed or accounted for in admission's RAM estimate. On a 16 GB machine that can
   be the difference between "fits" and swapping.
8. **No placement evidence for chat.** Image parses the engine's auto-fit plan and per-module
   split; chat does not parse "offloaded N/M layers" or per-device buffer sizes, so the app
   cannot say whether a chat model actually ran on GPU, partially offloaded, or on CPU.
   **Resolved 2026-11-24:** chat now requests level-4 load logs and reports the actual layer
   split, context, batch and model/KV/compute buffers after health. See
   `docs/plans/CHAT-PLACEMENT-EVIDENCE-2026-11-24.md`. This is load-time evidence, not
   continuous telemetry or release verification.
9. **Admission is fail-closed with no explicit override.** Missing live VRAM, unknown
   detection, or an over-budget request returns an error and no lease. The owner's rule
   requires a warned, explicit manual path instead.
10. **No CPU, temperature or storage telemetry.** The snapshot has GPU/RAM only. Temperature
    is unavailable on this platform through the current detectors and must be reported as
    unavailable, not guessed. Storage free space matters for disk residency and mmap.
11. **Level-4 load logs contain prompt text.** Placement evidence requires llama-server's
    `load_tensors` lines at verbosity 4, but that same level logs the prepared prompt after
    the model is loaded. A placement parser that forwards the whole stream into a failure
    capture would persist user content. The implemented guard forwards only through
    `model loaded` and clears the capture at health.

---

## 7. Proposed integration with the hardware-aware backend

### 7.1 One planner, two executors

The existing image split is the model to copy:

```
hardware snapshot ──► policy planner (Go, one authority)
                         │  decides: may it start? with what budget?
                         ▼
                    engine placement (llama.cpp fit / sd.cpp auto-fit)
                         │  decides: which exact tensors live where
                         ▼
                    evidence parser (logs)
                         │  reports: actual device, layers, RAM/VRAM split
                         ▼
                    UI: plan + actual + reasons
```

The planner must not re-implement tensor placement, and the engine must not be asked to make
a policy decision. One set of numbers flows down; one observed result flows back up.

### 7.2 Plan shape

A single `InferencePlan` should carry every decision and its reason, so the UI and the tests
read the same object:

```go
type InferencePlan struct {
    Modality      string // llm | image | asr | tts
    Device        int
    GPULayers     int    // -1 = engine auto
    NCPUMoE       int    // MoE expert layers in RAM
    NCPUFFN       int    // dense FFN layers in RAM
    KVTypeK       string // "", q8_0, q4_0 ...
    KVTypeV       string
    Batch         int
    UBatch        int
    Context       int
    LoadMode      string
    ReserveMB     int64  // measured headroom, per device
    RequiredVRAM  int64
    AvailableVRAM int64
    Confidence    float64
    Reason        string
    Fallbacks     []string // ordered alternatives that were considered
}
```

### 7.3 Selection ladder

1. Full GPU if the estimate plus reserve fits.
2. MoE model and short on VRAM → `-ncmoe` with N computed from the measured gap, then clamp
   N to a usable range; dense model → `-ncffn`.
3. Still short → reduce `-ngl`, then micro-batch `-ub`.
4. Still short → `-ctk/-ctv q8_0`, reported explicitly.
5. Still short → CPU-only plan with a clear warning.
6. Still short/unmeasurable → **manual**: show the numbers, keep the user's settings, offer
   "load anyway" and the advanced levers; never a silent fallback and never a bare refusal.

### 7.4 No-lockout semantics

- Auto-plan failure is a *warning*, not an error. The load may proceed with the user's
  settings after an explicit acknowledgement.
- An explicit override is recorded as such (in the plan and in the log), so it is not a
  hidden risk acceptance.
- A model is never rejected because optional acceleration is absent. CPU is always a valid
  plan.
- Telemetry failure degrades to "unknown capacity, user-directed" and says so.

### 7.5 Truthful telemetry

- Report `vram_live`, `vram_source`, snapshot age and per-device values.
- Temperature: report **unavailable** unless a real sensor path exists; never estimate.
- CPU load and storage free space: add real readings; storage is required for disk residency
  and mmap decisions.
- Remove name-based VRAM guessing from any path that presents capacity as fact; keep it only
  as an explicitly labelled estimate, or drop it.

### 7.6 Controls

- Capability-probe the llama binary exactly as sd.cpp already does; only send flags the
  binary advertises.
- Make Flash Attention tri-state (`auto/on/off`), not a boolean that cannot express "off".
- Either wire or remove `kv_cache_type`, `kv_cache_size` and `offload_kqv`. A dead control is
  worse than no control.
- Expose the new levers as advanced/manual first (`-nkvo`, `-lm`, `-cram`, `-ot`,
  `--params-backend disk`), and promote to automatic only after measurements on this box and
  at least one report from a different configuration.

---

## 8. Recommended order of work

| Priority | Slice | Why first | Verification |
|---|---|---|---|
| P0 | Capability probe for llama + truthful tri-state flags; remove or wire dead KV/KQV controls | stops the UI lying; prerequisite for every later flag | unit tests on probe/args; both hosts show the control state matching the binary |
| P0 | One chat planner: measured layers/reserve instead of fixed `-ngl`; plan surfaced | unlocks the engine's existing hardware awareness | load on this box; parse log; compare planned vs actual layers; UI shows both |
| P1 | MoE/dense offload selection (`-ncmoe`/`-ncffn`) with manual override | biggest gain for models larger than VRAM | MoE model load; log shows expert placement; VRAM measured |
| P1 | KV quantization + micro-batch as ordered fallbacks | covers dense models where FFN offload is not enough | forced-short budget; plan walk; speed/VRAM recorded |
| P1 | Non-blocking manual override path in admission | owner rule; prevents lockout when telemetry fails | simulate detection failure; load proceeds only after explicit choice |
| P2 | Image resolution/steps in the preflight; module CPU/disk manual options | closes the "fits but OOMs" gap | preflight vs sd-server capacity check at several resolutions |
| P2 | Placement evidence for chat; telemetry extensions (CPU/storage/temp-unavailable) | makes decisions visible and honest | log/UI comparison; unavailable states shown as unavailable |

Nothing in P0–P2 requires a new inference library, a new model manager, or a new UI page. It
uses the flags the shipped binaries already advertise and the planner/reporting pattern the
image path already follows.

---

## 9. Limitations of this survey

- Every performance figure quoted in section 4 is someone else's machine. None was measured
  here, and none should be presented to a user as this machine's expected speed.
- Whisper is not installed on this box, so its placement options are unverified. Its model
  sizes are small enough that it is not the priority.
- MoE and dense-FFN offload could not be exercised: the available models on this machine are
  not MoE, and no MoE fixture was loaded during this research.
- sd.cpp `--params-backend disk`, module CPU compute, SageAttention and `--moe-cache-mib`
  are documented and supported by the binaries but were not run in this pass.
- The 10 % reserve is implemented for chat admission; the image path's fixed 1 GiB reserve is
  an engine-compatibility choice that still needs an explicit policy decision.
