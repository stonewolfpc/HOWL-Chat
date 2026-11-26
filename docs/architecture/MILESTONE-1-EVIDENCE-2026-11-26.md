# Milestone 1 evidence — automatic chat placement and KV accounting

Date: 2026-11-26
Commit: `b6274b0` ("Plan chat GPU layers upward under the admission authority")
Hosts: packaged desktop (WebView2) and browser dev host, same isolated root
(`C:\HOWL-Test\hra-image-20261125`, `.howl-isolated-root` present)

## Owners after this milestone

| State | Owner |
|---|---|
| Chat layer/context plan and admission decision | `resource.Authority.Admit` (single policy authority) |
| Upward layer search | `resource.PlanChatLayers`, called only by `Admit` |
| KV cost per element | `resource.KVCacheBytesPerElement`, consumed by `estimateChat` |
| Requested placement | `settings.gpu_layers` / `context_size` / `batch_size` |
| Observed placement | `llamabackend.Placement`, reported by `chatPlacementReport`; never fed back into admission |

Competing owners found: none added. The reserve policy is still duplicated between
`internal/resource/allocator.go` and admission — deliberately untouched in this
milestone and still the Milestone 2 reconciliation target. No deletions were made;
`allocator.go` eviction and `gguf.ValidateSpawn` keep their policy and validation
logic until the centralized authority is proven across every call path.

## What was implemented

- Automatic requests (`gpu_layers = -1`) run the existing estimator through
  `Admit` and return a concrete layer count in `AdmissionAdjustment{Auto: true}`.
  Explicit `0` remains CPU-only; saved zeroes are not migrated.
- KV cost is charged per element for the live cache type (f16/unknown 2, q8_0 1,
  q4_0 1) instead of assuming f16.
- A chat load that wants the GPU frees a resident image model first
  (`FreeVRAMFor(..., true)`); a blocking automatic request that fits zero layers
  refuses instead of demoting to CPU.

## Evidence, per host

Both hosts used the same model (`Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M`),
the same isolated root, and the same committed binary sources.

| | browser dev host | packaged desktop |
|---|---|---|
| admission | `6057 MB VRAM required, 6468 MB available after 718 MB reserve; fitted plan context=3096 gpu_layers=32 batch=512` | `6057 MB VRAM required, 6159 MB available after 684 MB reserve; fitted plan context=3096 gpu_layers=32 batch=512` |
| argv | `-c 3096 -b 512 -ngl 32 --fit off` | `-c 3096 -b 512 -ngl 32 --fit off` |
| engine context | `llama_context: n_ctx = 3328` | `llama_context: n_ctx = 3328` |
| engine placement | `layers=32/33 vram=4677MB ram=679MB context=3328 matches_plan=true` | `layers=32/33 vram=4677MB ram=679MB context=3328 matches_plan=true` |
| visible result | "Model loaded as planned: planned 32 GPU layer(s), context 3096, batch 512; … Engine placed 32/33 GPU layers …" | same message in the model panel |
| GPU after load | 6229 MiB used (≈600 idle) | 6553 MiB used (≈600 idle) |

Logs: `build/tmp/m1-auto/harness.out`, `build/tmp/m1-auto/packaged3.out`.

## Detail 1 — the browser host drove the page handler, not a physical click

In the browser dev host the Load control sits below the panel's visible area, and
synthetic pointer events did not reach it. The run therefore invoked the page's
own `loadSelectedModel()` through the browser tooling. This exercises the real
frontend load path against the real backend over the harness RPC bridge, but it is
**not** a physical UI click. The packaged-host run is the physical-input evidence:
in the maximized 2576×1456 window the model was selected through Browser → model
row → **Load Model** with real mouse events.

The packaged picker itself could not be driven by synthetic clicks: the custom
dropdown's option rows are laid out outside the list's scroll viewport, so a click
closed the menu without changing the selection. The Browser-detail route was used
instead and reaches the same `LoadModel` binding with the same `-1` default.

## Detail 2 — requested context 3096 vs effective context 3328

This is expected engine normalization, not a plan departure, and it is now
recorded as such rather than hidden behind the match flag.

- Requested: `-c 3096` (the admission plan's context).
- Effective: `llama_context: n_ctx = 3328` — llama.cpp pads the requested value up
  to a 256-token block: `cparams.n_ctx = GGML_PAD(cparams.n_ctx, 256)` in
  `src/llama-context.cpp` (upstream `ggml-org/llama.cpp`, present in every
  published revision checked). `ceil(3096 / 256) = 13`, `13 × 256 = 3328`, so the
  engine gains 232 tokens; `n_ctx_seq` and the KV cache follow at 3328 cells
  (`llama_kv_cache: size = 104.00 MiB (3328 cells, 8 layers)`).
- Vocabulary: **requested context** is what HOWL asks for and what admission
  costs; **effective context** is what the engine allocates. `matches_plan` treats
  `effective == requested` or `effective > requested` by at most 256 tokens as a
  match; a smaller context or a different layer count is a real departure.
- The distinction is user-visible, not only internal: the report reads
  "…llama-server rounded the requested context 3096 up to its token block".

Known consequence, recorded not fixed: admission estimates KV at the requested
context, while the engine allocates the padded one. KV is therefore understated by
up to 256 tokens per sequence (here 104 MiB allocated against a 3096-token
estimate). It is bounded, fits inside the post-reserve headroom in both runs, and
belongs to the Milestone 3 capacity work rather than to this one.

## Verified limitations

- No chat turn was generated; placement was verified, inference was not.
- The browser-host leg is a handler-path test, as described above.
- Image optimization, audio, and the deferred deletions were not touched.
- `GetChatPlacement` was read through the placement log and the on-screen message;
  the RPC read-back was not separately exercised.

## Next

Milestones 2–7 as agreed, in the frozen plan's dependency order: authority and
lifecycle reconciliation, hardware intelligence and reserve decisions, model
intelligence across architectures, bounded image placement (HiDream-class),
cross-workload coordination including audio's CPU/RAM-first behavior, then full
acceptance across both hosts.
