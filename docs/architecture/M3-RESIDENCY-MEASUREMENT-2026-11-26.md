# M3 — What the engines actually consume

Date: 2026-11-26
Isolated root: `C:\HOWL-Test\hra-image-20261125` (`.howl-isolated-root`)

M3 replaces the planning bound with a measured figure where one exists, and
says so where none exists. The ledger still plans from the admission estimate;
it now also records what the engine and the device reported, and takes the
conservative of the two.

## Method

Each workload load is bracketed by two reads of the same snapshot the admission
gate uses (`TotalVRAM − FreeVRAM`), one at decision time and one after the
engine reports its health. That difference is the device's own statement about
what this start cost. The engine's statement comes from its placement: chat from
the llama.cpp placement report, image from `sdcpp` placement.

The ledger policy is one line: **the larger of (engine-reported total,
device delta)**, falling back to the plan only when neither was measured.
`internal/resource/residency.go` owns that policy; `app-residency.go` owns the
recording and the single measured-update path into the registry.

## Chat (llama-server, Qwen3.5-9B Q4_K_M, 32/33 layers)

Packaged host, auto-load from `general.defaultModel = "last-used"`:

```
[LoadModel] admission: admitted=true manual=false admitted: 6057 MB VRAM required,
            6309 MB available after 701 MB reserve; fitted plan context=3096 gpu_layers=32 batch=512
[Residency] kind=llm model="Qwen3.5-9B-..." planned=6057MB engine=4887MB device_delta=5222MB ledger=5222MB within plan
[LoadModel] actual placement: known=true layers=32/33 vram=4677MB ram=679MB context=3328 matches_plan=true
```

Finding: the engine's VRAM total (4887 MB) is **335 MB below** what the device
actually lost (5222 MB). The difference is CUDA context, allocator slack and the
process's own overhead, none of which llama.cpp counts in its reported total.
The plan (6057 MB) is above the device truth, so the reservation stays the
planning bound; the ledger now carries 5222 MB instead of the plan.

Reproduced twice on the same build (5236 and 5222 MB), so the gap is stable
within ~15 MB, not a sampling artefact.

## Image (sd-server master-1122-3b12781, Illustrious-XL-v2.0, q8_0)

Load plus a real 512×512 / 8-step txt2img generation, harness RPC driving the
same binding the UI calls:

```
[App] image admission final: admitted=true projected=4613MB resident=4101MB activation=0MB
      required=5084MB ... free=7063MB reserve=819MB
[Residency] kind=image model="Illustrious-XL-v2.0" planned=5093MB engine=0MB device_delta=5679MB ledger=5679MB above plan
[ImageReclaim] decision=retain reason="resident baseline established; weights retained" ...
```

Two findings, both load-bearing:

1. **This sd-server build reports no placement at all.** `GetImagePlacement`
   returns `known=false`, `vramMB=0`, empty components, with the explicit note
   that the build did not report per-component placement. Its only measurements
   are process working sets (170 MB after load, 409 MB during generation).
   Engine-reported image VRAM is therefore zero *by fact*, not by omission, and
   the device delta is the only image evidence available.
2. **The image plan understates residency.** The admitted projection is
   5093 MB committed; the device lost 5679 MB — 586 MB more. Weights at q8_0
   plus the VAE/CLIP/tiling working set exceed what the projection reserves.
   The ledger now holds 5679 MB (`GetLoadedModels` → `gpu_memory_mb: 5679`),
   and `GetResidencyReport` reports `agrees_with_plan: false`, which is the
   honest answer.

A third finding shaped the code, not just the numbers: **image residency cannot
be measured at load time.** sd-server loads weights lazily; a load-time
observation returned `engine=0, device_delta=0` and would have been an invented
figure. The load now records a baseline (planned MB and pre-load device usage)
and the observation is taken after the first generation, where both the engine's
placement and the device's usage are real. No observation is recorded before
that point.

## What changed

- `internal/resource/residency.go` — observation type, conservative ledger
  policy, plan agreement, negative measurements discarded.
- `AdmissionEstimate.CommittedVRAMMB()` — exported so code that reasons about a
  running engine uses the projection rather than the resident-only figure.
- `app-residency.go` — recorder, `[Residency]` evidence line, baseline holding,
  `GetResidencyReport`, and the one path that rewrites a registry entry from a
  measurement (registry remains the single ledger owner).
- `app-model.go` — chat measured across the health transition; admission now
  registers the observed figure.
- `app-image-model.go` — image registers the committed projection and holds the
  measurement baseline; `app-image-reclaim.go` records the post-generation
  observation and reconciles the ledger.

## Explicitly not done

- **Reserves untouched.** R1 admission, R2 allocator, R3 image reclaim report and
  R4 `--max-vram` are still four computations. Nothing consolidated: the
  measurements above are the input that work would need, not a licence to start.
- **RAM is still planned, not measured.** Only VRAM was bracketed. The image
  entry keeps the planned 659 MB system RAM; no delta was read for it.
- **HiDream not exercised here.** Its classification stays as M2 left it
  (verified unsafe estimate/launch combination, inferred historical mechanism);
  bounded HiDream placement is M5.
- **Host difference recorded.** The packaged image panel exposes no accessible
  nav targets (icon-only elements with empty names), so the image leg was driven
  through the harness RPC on the same App bindings. The generated image and the
  engine logs are real; the packaged-host image *UI* path remains unexercised
  and is not claimed.
