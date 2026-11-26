# M5 — Bounded image placement

Date: 2026-11-26
Isolated root: `C:\HOWL-Test\hra-image-20261125` (`.howl-isolated-root`)

A model whose full resident set is larger than the card used to have one honest
outcome: refusal. M5 lets Howl bound the start instead, and proves on real
hardware that the bound is respected through a completed generation.

## What changed

`planBoundedPlacement` is now the single placement planner. It takes a budget
and selects the measured bounded mode that fits it:

- **cpu-diffusion** — diffusion and CLIP on CPU, VAE on CUDA and tiled, weights
  staged from system RAM (`--backend diffusion=cpu,clip=cpu,vae=cuda0`,
  `--offload-to-cpu`, `--vae-tiling`);
- **cpu-only** — below the measured GPU-assisted process need, no GPU weights or
  compute at all.

Previously this was only reachable when the user had set an explicit
`--max-vram`. Now, when admission refuses a full residency, Howl re-plans inside
the budget the authority just computed (the post-reserve available figure it
reported), sets that budget as an automatic `--max-vram` cap, and re-admits the
bounded placement. The budget is consumed, never recalculated: no reserve math
was touched.

Guards stay: an explicit `Backend`/`ParamsBackend` is the user's placement and
is never rewritten; an explicit positive `--max-vram` was already planned before
admission; `imageBoundedPlacementEligible` rejects both.

## Evidence: a model larger than the card now runs

Flux-dev Q3_K_M with its T5 encoder: resident need 10,754 MB against 8,192 MB.

```
[App] image admission: admitted=false ... needs 11761 MB VRAM; ... 6275 MB available
[App] image hardware: VRAM Budget 6 GiB is below the 8311 MB resident weights ...
      Automatic partial placement: diffusion and CLIP on CPU, VAE on CUDA and tiled.
[App] bounded placement admitted: mode=cpu-diffusion budget=6241MB
Starting sd-server ... --diffusion-model flux1-dev-Q3_K_M.gguf --t5xxl t5xxl_fp8_e4m3fn.safetensors
      --max-vram 6 --backend diffusion=cpu,clip=cpu,vae=cuda0 --vae-tiling --offload-to-cpu
[I] Detected model family: Flux (flux1-dev-Q3_K_M.gguf)
```

A real 512×512 generation then **completed** (`GetImageJobStatus` →
`completed`), and the post-generation measurement confirms the process stayed
far inside the bound:

```
[Residency] kind=image model="flux1-dev-Q3_K_M" planned=2048MB engine=0MB device_delta=1836MB ledger=1836MB within plan
[ImageReclaim] decision=retain ... baseline=10754MB ... free=6356MB
```

10.7 GB of weights live in RAM; the device delta is 1.8 GB against a 6 GiB cap.
That is bounded placement working, not an estimate of it.

## HiDream: classification corrected

The HiDream-class case does **not** fail on memory on this engine build. Both
HiDream GGUFs are rejected by the engine itself, before any sampling:

```
{"error":"cannot identify model family from 'H:\\AI GGUF\\Image models\\hidream-i1-fast-Q2_K.gguf'"}
{"error":"cannot identify model family from 'H:\\AI GGUF\\Image models\\hidream-i1-dev-q2_k.gguf'"}
```

Verified directly against `sd-server` (3b12781), with `-m` and with
`--diffusion-model`, both files. Through the app, admission and the bounded
placement both succeed and the process starts inside the bound; the failure
arrives as a surfaced `txt2img error 400`, and the reclaimer marks the engine
abandoned and restarts it rather than pretending the job ran.

So the previous "unsafe estimate/launch combination" reading is superseded: the
verified fact is **engine-capability-blocked**, not a VRAM shortfall. M6/M7 must
not claim HiDream placement proof on this build. The bounded path is proven on
Flux instead, which is the same shape of problem (weights ≫ device) with a
family the engine can identify.

## Owners

Unchanged, deliberately:

- admission (`resource.Authority`) still decides *whether* the workload starts;
- `planBoundedPlacement` decides *where* modules go, from the budget it is given;
- sd.cpp's `--auto-fit`/`--backend` remains the engine's placement executor;
- the registry remains the resident ledger, and M4's observation floor still
  applies independently of the placement chosen.

No reserve was consolidated or recalculated: R1–R4 are untouched.

## Verified limitations

- **Packaged-host image UI not clicked.** Evidence again comes from the browser
  dev host driving the same `App` bindings.
- **One step.** The Flux generation used 1 step at 512×512: it proves execution
  and memory behaviour, not image quality or throughput.
- **CPU diffusion is slow.** Not benchmarked here; M6 owns cross-workload
  performance evidence.
- **Bounded mode is a session decision.** Existing rules carry it forward for a
  resident workflow; it is not persisted as user intent.
- **Not every engine failure is diagnosable from the 400.** Howl surfaces the
  engine's message; it does not yet classify "engine cannot identify this
  architecture" separately from other generation errors.
