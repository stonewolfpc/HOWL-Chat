# Stable Diffusion.cpp Workflow Orchestration — Plan and Verification Matrix

**Date:** 2026-11-24
**Status:** Research complete. No code changed. No live image-generation experiment run.
**Gate:** Implementation and live verification require owner authorization.

**Scope:** stable-diffusion.cpp `sd-server` commit `3b12781` as shipped in this
repository, and Howl's image execution, admission, placement, and reclamation paths.

**Success target:** one sequential workflow on the GTX 1080 8 GB:

```
base generation -> ESRGAN 2x upscale -> pose/ControlNet img2img -> ESRGAN 2x upscale
```

Each stage waits for its predecessor, obtains admission, executes, reports actual
placement, and transitions safely. Compatible components stay resident when the
measured and projected budget permits. Reclamation remains the existing single owner.

---

## 1. Research result in one page

### 1.1 What the installed engine actually exposes

Installed evidence:

- `build/verification/sd-server-help-installed.txt` (232 lines, commit `3b12781`).
- `build/verification/sd-cli-help-installed.txt` (255 lines, commit `3b12781`).
- `build/verification/dgl-builtins.txt` (complete DGL built-in catalog).
- Local binary route strings include `/sdcpp/v1/dgl/scripts`, `/dgl/script`,
  `/dgl/jobs`, `/dgl/jobs/{id}`, and `/dgl/jobs/{id}/outputs/{index}`.
- The local binary has **no** `/sdcpp/v1/upscale` route. Current upstream `master`
  documents that route, but it is not in `3b12781`. It must not be planned as if it
  exists.

DGL is real and usable in this build. It can load models, generate, apply ControlNet,
upscale, save/encode media, and return tensors. It is not a graph UI in Howl; it is an
HTTP job API plus a scripting language.

### 1.2 The decision that determines the architecture

DGL has its **own** `ModelManager` and resource cache (`src/dgl/context.cpp`:
`Context::models`, `Context::resources`). It is obtained from the server context but
does not share the server's startup `-m` model cache. Therefore:

- Loading the base checkpoint in a DGL script while `sd-server` already has it loaded
  through `-m` would hold two copies of the same 5 GB q8 model.
- On an 8 GB card that is not a safe plan, even with offload.
- All-DGL staging is therefore **rejected** for this PoC.

The selected design is hybrid:

| Stage | Executor | Why |
|---|---|---|
| 1. Base generation | existing native renderer (`GenerateTxt2Img`) | reuses the already-admitted resident model and all existing admission/q8/VAE-tiling policy |
| 2. ESRGAN 2x upscale | DGL job (`load_upscaler` + `upscale`) | the only verified ESRGAN path in `3b12781`; loads only the small upscaler into the DGL cache |
| 3. Pose/ControlNet img2img | existing native renderer (`GenerateImg2Img`) | reuses the same resident checkpoint and the server's ControlNet slot; the txt2img ControlNet path is proven, the img2img variant is verified by I-3 |
| 4. ESRGAN 2x upscale | DGL job, same script/options as stage 2 | reuses the upscaler cached by the DGL context |

This gives one resident `sd-server`, one main model copy, one small DGL upscaler copy,
and no competing GPU work. The DGL job and native generation both take the server's
generation/context lock, so they cannot overlap.

### 1.3 What is supported, what is not

Legend:

- **Supported** — reaches the engine today and is meaningful.
- **DGL-only** — real in this build, but only through DGL; no native route.
- **Engine-supported, Howl-unwired** — the engine can do it, Howl does not send it.
- **Unsupported / incompatible** — cannot be used as advertised.
- **Unverified** — needs an authorized live run; not claimed.

| Capability | Evidence | Howl today | Classification | PoC use |
|---|---|---|---|---|
| Base txt2img | `/sdapi/v1/txt2img`, existing client | yes | Supported | Stage 1 |
| Img2img / inpaint / outpaint | `/sdcpp/v1/img_gen`, existing client | yes | Supported | Stage 3 (img2img) |
| ControlNet single generation | startup `--control-net` slot; txt2img proven in `STAGE_STATUS.md`; img2img variant not yet run | yes, bound to model | Supported (txt2img); img2img unverified | Stage 3 |
| ESRGAN upscale | DGL `load_upscaler` + `upscale` | no; UI "Upscale" is diffusion img2img | **DGL-only** | Stages 2 and 4 |
| HTTP `/sdcpp/v1/upscale` | upstream only; absent in local `3b12781` | no | **Incompatible with installed build** | Rejected |
| DGL script jobs | `/sdcpp/v1/dgl/jobs` in local binary | no | **DGL-only** | Stages 2 and 4 |
| DGL shared component cache | same DGL context across jobs; cache keyed by sources + loader options | no | **DGL-only** | Upscaler reuse |
| DGL ControlNet | `controlnet::load` / `controlnet::apply` | no | Engine-supported, Howl-unwired | Not used in PoC; native path is proven |
| OpenPose pose **extraction** from a photo | no pose estimator in DGL built-ins; no route | no | **Unsupported** | External dependency or prepared pose |
| OpenPose skeleton **drawing** | `frontend/dist/js/ig-pose.js`, COCO-18 | yes | Supported (manual/prepared input) | Pose input for stage 3 |
| Hires fix / second-pass refinement | sdapi `enable_hr`/`hr_upscaler`; DGL latent resize + sigma tail | txt2img only; `img2img` does not send it | Engine-supported, Howl-unwired | Not PoC; ESRGAN is the upscale proof |
| Canny / image preprocessing | DGL `image::canny`; native `image_preprocess` | no | Engine-supported, Howl-unwired | Not PoC |
| Batch | existing sequential batch | yes | Supported | Repeat/batch verification |
| LoRA | existing structured LoRAs | yes | Supported | Unchanged |
| VAE tiling | engine auto-retry + `--vae-tiling` | yes | Supported | Stage 1/3 pressure |
| CPU/CUDA/RAM/disk placement | startup flags; DGL loader flags | yes for native; DGL mapping needs work | Supported, DGL mapping unverified | Constrained cases |
| Concurrent GPU stages | server serializes native jobs and DGL jobs | not allowed | **Must not do** | Sequential only |
| Multi-GPU / RPC | flags exist, no local verification | no | Unverified / out of scope | No |

### 1.4 The two explicit dependencies that must not be hidden

1. **ESRGAN weights are not currently in the image registry.** The registry has
   OpenPose ControlNet entries, but no `upscaler` type. `build/scancheck.txt` shows
   `RealESRGAN_x2plus.pth` and `RealESRGAN_x4plus.pth` under an InvokeAI model folder,
   but the image page's model registry does not expose them. The PoC needs an explicit
   upscaler path (and a scan/registration step, or a workflow-only asset field).
2. **Pose extraction is external.** `ig-pose.js` draws a COCO-18 skeleton; it does not
   extract keypoints from a photo. The engine has no OpenPose/DWPose extractor in its
   DGL built-ins. For an automated photo-to-pose input, an external OpenPose/DWPose
   processor is required. The PoC may use a prepared skeleton from the existing editor;
   that is pose **preparation**, not extraction, and must be labelled that way.

---

## 2. Engine contract details that the implementation must honor

### 2.1 DGL job API (verified from `routes_dgl.cpp` at `3b12781`)

| Method | Path | Request | Response |
|---|---|---|---|
| `POST` | `/sdcpp/v1/dgl/jobs` | `{"source":"<DGL source>"}` | `202` + async job JSON |
| `GET` | `/sdcpp/v1/dgl/jobs/{id}` | — | job status, progress, `result` when completed |
| `GET` | `/sdcpp/v1/dgl/jobs/{id}/outputs/{index}?frame=N` | — | `image/png` for a U8 tensor output |
| `GET` | `/sdcpp/v1/dgl/scripts` | — | saved script names (not needed for direct source) |
| `GET/POST` | `/sdcpp/v1/dgl/script` | name/source | saved script read/write (not needed) |

Requirements enforced by the server:

- Only loopback clients; `Host` must be `localhost`, `127.0.0.1` or `[::1]`, optionally
  with the local port.
- `Origin`, if present, must be same-origin; cross-origin is `403`.
- `POST` content type must be `application/json`; otherwise `415`.
- `source` must be non-empty UTF-8; otherwise `400`.
- Job states: `queued`, `generating`, `completed`, `failed`, `cancelled`.
- Completed DGL result shape:
  `{type, outputs:[{name,dtype,shape,bytes,frames}], media:[{name,mime_type,b64_json}]}`.
- Output index is the position in `result.outputs`; `frame` defaults to `0`.
- Completed jobs are retained by `--completed-ttl-seconds` (600 s default). Fetch the
  output promptly; do not treat the job as a durable artifact store.

### 2.2 DGL model-path access

`allowed_dgl_path` allows model reads from:

- startup model sources (`-m`, `--vae`, `--control-net`, etc.) via `configured_files`;
- any path under a `--model-dir` root;
- `--lora-model-dir`;
- `--hires-upscalers-dir`.

Plain `SD_DGL_FILE_READ` (images) only requires a readable file. Model reads
(`SD_DGL_FILE_READ_MODEL`) require the allowlist.

Consequences:

- The base checkpoint loaded by `-m` is already allowed for DGL, but the PoC does not
  load it again in DGL.
- The ESRGAN upscaler must be under `--hires-upscalers-dir` (or a model root).
  Add that flag from the workflow's upscaler directory.
- The ControlNet used by the native renderer must be startup `--control-net` (the
  server binds it to the model) or under a model root. Add a startup controlnet field
  to `ModelLoadOptions`, emitted as `--control-net` when set.
- `image::save` may write only under `<server cwd>/output`. The PoC should not rely on
  it: return the image tensor and fetch `/outputs/{index}`. Howl owns the artifact.

### 2.3 DGL context lifetime and cache

- All DGL jobs share one DGL context for the server's lifetime.
- Components are cached by weight sources and loader options. The same upscaler path
  and same `tile_size` hit the cache on stage 4.
- The cache can release idle weights under memory pressure.
- `--weight-residency-ttl` can unload DGL components when idle. The Howl default is
  `-1` (disabled). For the reuse proof, keep `-1`. If the user explicitly sets a TTL,
  respect it and record a stage-4 reload rather than overriding user intent.
- `--weight-staging-ttl` is not used; prior measurement showed no generation-boundary
  effect in this build.

### 2.4 Model paths and weight type

The DGL upscaler loader takes `backend`, `params_backend`, `tile_size`, `direct`.
It does not take a weight type. The base checkpoint's q8 choice does not affect the
upscaler. For future DGL generation scripts, the loader `wtype` must match Howl's
admitted `WeightType` (`q8_0` when auto-quantized); loading the same file at f16 in
DGL would defeat the admission budget. This is not needed for the selected hybrid.

---

## 3. Howl architecture mapping

Existing authorities stay exactly as they are:

| Concern | Existing owner | Workflow use |
|---|---|---|
| Admission, 10% preferred reserve + 500 MB soft allowance, lease | `resource.Authority.Admit`, `App.admitImageLoad` | called before each stage that can allocate; no new planner |
| Load policy (auto q8, VAE tiling, max-vram cap, constrained placement) | `app-image-model.go` | reused for stage 1 and stage 3 |
| Actual loaded state | `Image.LoadedOptions` | reused; DGL upscaler residency is recorded separately as orchestration state |
| Placement reporting | `app-image-placement.go` | extended additively with workflow stage rows |
| Reclamation | `app-image-reclaim.go` (`reclaimImageEngineAfterGeneration`) | called at workflow terminal or failure; never a second unload path |
| Hardware snapshot | `resource.GetManager().RefreshSnapshot()` | same source for stage boundaries |
| Job serialization | image job manager + engine generation lock + server queue | workflow holds the image start mutex for the whole sequence |

No new resource authority, no second reclamation mechanism, no direct
`resource.UnregisterModel` call from the workflow.

### 3.1 Why hold the workflow lock for the whole sequence

If another image job starts between stages, it can take the GPU, change the loaded
model, or force a restart. The workflow must own the engine from stage 1 admission to
the terminal report. The public binding acquires `imageStartMu` once and calls an
internal runner that does not re-acquire it. This is the same mutex
`ensureEngineReady` already uses; no new lock order is introduced.

### 3.2 Stage admission in delta form

Free VRAM already reflects the resident model and ControlNet. Charging the full model
again at stage 2 would double-count it and produce false refusals. The runner builds an
additive request for each stage:

| Stage | Resident before stage | Charge at admission |
|---|---|---|
| 1 base | none (or an old model to be replaced) | checkpoint + conditioner + VAE + controlnet + activation |
| 2 upscale | base + controlnet | upscaler bytes + upscale activation; no base charge |
| 3 ControlNet img2img | base + controlnet + upscaler | activation only; no duplicate weights |
| 4 upscale | base + controlnet + upscaler | activation only; upscaler is resident; charge again only if the cache was lost |

If stage 1 fails or the server restarted, the next stage falls back to the existing
full-weight request. This is a caller-side choice between the existing full request and
a new delta request shape, not a new estimate algorithm.

`AdmissionRequest` needs one additive field, `UpscalerBytes`, included in the existing
weight total and component count. That is a measured file-size addition, not a new
placement authority.

Activation is separate by design: `estimateImage` reports `ActivationVRAMMB` but does not
add it to `RequiredVRAMMB` (`internal/resource/admission.go`). The workflow runner must
therefore gate each stage on activation versus post-reserve free memory itself, and use
the existing tiling levers (native VAE tiling; DGL `tile_size`), rather than claiming
`Admit` enforced activation. If a hard admission refusal is required, add an additive
`WorkloadVRAMMB` field and include it in the estimate with tests proving existing image
loads are unchanged; do not silently change the meaning of `ActivationVRAMMB`.

### 3.2.1 Preferred reserve, bounded soft allowance

The 10% free-VRAM reserve remains the target. Admission still plans tiling and
ControlNet CPU placement against it first. When a stage still does not fit, it may use
at most 500 MB of that reserve, and it must keep at least 128 MB free. A request that
needs more than 500 MB of the reserve is refused with the exact required, free,
reserve, available and soft-cap numbers. An explicit measured reserve
(`ReserveOverrideMB`/`ReserveExplicit`, such as a deliberate `--max-vram` reserve) is a
hard requirement and is never softened. A positive `--max-vram` spend cap also remains
hard. This changes no reclamation or admission ownership: it is one final-fit branch
inside the existing `Authority.Admit`.

### 3.3 Reclamation policy for a workflow

- Transient tensors are released by the engine at the end of each native generation and
  by the DGL context at the end of each DGL job.
- Reusable components (base model, ControlNet, upscaler) stay resident.
- The workflow does **not** call a process recycle between successful stages. That would
  destroy the upscaler cache and force a base-model reload for no memory reason.
- At the terminal of the workflow, and on any stage failure, the workflow calls the
  existing `reclaimImageEngineAfterGeneration` exactly once. The existing measurement
  band, reserve branch, and abandoned-engine recovery decide whether a recycle is
  warranted.
- If a stage is refused by admission because free VRAM is already at the reserve, the
  failure path calls the same owner; the workflow does not unload anything directly.

This preserves the single-owner rule and prefers residency.

---

## 4. Proposed workflow contract

### 4.1 Inputs

```go
type ImageWorkflowOptions struct {
    ModelName       string
    Prompt          string
    NegativePrompt  string
    Seed            int64

    BaseWidth       int
    BaseHeight      int
    BaseSteps       int
    BaseCfg         float64
    Sampler         string
    Scheduler       string

    UpscalerPath    string // ESRGAN .pth/.safetensors
    UpscalerTile    int    // default 128
    UpscaleFactor   int    // model-native; 2 for this PoC

    ControlNetPath  string
    PoseImagePath   string
    ControlStrength float64

    RefineWidth     int
    RefineHeight    int
    RefineSteps     int
    RefineStrength  float64

    BatchCount      int    // sequential complete workflows
}
```

Defaults for the PoC: 512×512 base, 4–8 steps; 2x upscale to 1024; ControlNet img2img
at 1024 with strength 0.55; 2x upscale to 2048. If admission refuses the 1024 refine,
the run reports the refusal; a 768 refine is a separate, pre-declared constrained case,
never a silent downscale.

### 4.2 Stage record

Every stage returns:

```go
type ImageWorkflowStageResult struct {
    Index          int
    Name           string
    JobID          string
    Status         string // completed, failed, refused, skipped
    StartedAt      string
    DurationMS     int64
    PIDBefore      int
    PIDAfter       int
    WorkingSetMB   float64
    FreeVRAMMB     int64
    ReserveMB      int64
    Placement      string
    InputPath      string
    OutputPath     string
    Error          string
}
```

The workflow result contains the ordered stage records, the final image path, and the
existing placement/reclaim report. Engine-derived fields are limited to the pinned
build's actual evidence (§9.14): INFO loaded bytes per model component, VERBOSE planned
params target and per-runner compute workspace, DEBUG buffer eviction/release. It does
not claim a live component residency map the engine does not expose.

### 4.3 Stages

**Stage 1 — base generation (native).**

- `Mode = "txt2img"`, prompt, 512×512, seed.
- `ensureEngineReady` admits and loads the base model with the workflow's ControlNet as
  a startup slot.
- Output: PNG bytes from the existing result, written to a workflow temp directory.

**Stage 2 — ESRGAN 2x upscale (DGL).**

- `POST /sdcpp/v1/dgl/jobs` with the script below.
- Poll to `completed`; find the output named `image`; fetch
  `/outputs/{index}?frame=0` as PNG.
- Write `stage2-upscale.png`.

**Stage 3 — pose/ControlNet img2img (native).**

- `Mode = "img2img"`, `InitImagePath = stage2-upscale.png`, `ControlImagePath = pose`,
  `ControlNetPath = workflow controlnet`, `ControlStrength = 0.9`,
  `Width/Height = 1024`, `Strength = 0.55`.
- Existing `GenerateImg2Img` uses `/sdcpp/v1/img_gen`; the server binds the ControlNet
  startup slot, and the request carries the control image. This exact img2img + ControlNet
  combination is unverified in this repo; I-3 is its first acceptance test. Stage 3
  consumes the stage-2 pixels, so it is a real dependency, not a decorative stage.
- Output: `stage3-controlnet.png`.

**Stage 4 — ESRGAN 2x upscale (DGL).**

- Same script and `tile_size` as stage 2, different `input_path`.
- DGL cache should return the already-loaded upscaler; no second `load_upscaler` cost.
- Output: `stage4-final.png` (2048×2048 for a 1024 input).

### 4.4 DGL scripts

Stage 2 and 4 use the same template with different JSON-escaped paths:

```dgl
# version: 0.0.0

var {
    input_path = "<howl temp path, JSON-escaped>"
    upscaler_path = "<upscaler path, JSON-escaped>"
    tile_size = 128
}

fn main() {
    upscaler = load_upscaler(upscaler_path, tile_size = tile_size)
    image = upscale(upscaler, image::load(input_path))
    return { image = image }
}
```

Notes:

- `load_upscaler` has no `family`; it accepts `path`, `tile_size`, `backend`,
  `params_backend`, `direct`.
- `upscale` uses the model's native scale. For a 2x model, no `factor` is needed.
- Return the `image` tensor. The server's output endpoint encodes U8 tensors as PNG;
  no `image::save` path restriction is involved.
- Verify the output dtype in I-2. If `upscale` returns F32 in this build, the script must
  set `enlarged = image::to_u8(enlarged)` before returning it; the output endpoint
  returns 404 for a non-U8 tensor. This is a script-level adjustment, not a second
  pipeline.
- Paths must be generated with a real JSON string encoder, not string concatenation.
  Backslashes in Windows paths are the main injection/escaping risk.

If a future DGL stage needs the base model (for example, to avoid native/DGL
duplication in a different architecture), the loader must set `wtype` to the admitted
weight type and `backend`/`params_backend` to the admitted placement. That is not part
of this PoC.

---

## 5. Implementation plan (authorization-gated)

Ordered slices. Each slice is independently reviewable. No slice starts before approval.

### Slice A — DGL client and script builder (no GPU)

Files:

- `internal/image/engine/sdcpp/dgl.go` (new)
- `internal/image/engine/sdcpp/dgl_test.go` (new)
- `internal/image/engine/sdcpp/client.go` (small shared HTTP/output helper if needed)

Work:

1. `RunDGL(ctx, source)` posts `{"source":...}` with `Content-Type: application/json`,
   parses the `202` job, polls `/sdcpp/v1/dgl/jobs/{id}`, and returns a typed result.
2. `FetchDGLOutput(ctx, jobID, index, frame)` returns PNG bytes.
3. Validate status codes, `failed`/`cancelled` errors, empty outputs, and context
   cancellation. Do not invent a cancel route: the DGL route table at `3b12781` has no
   dedicated cancel handler.
4. Script builder with `json.Marshal` for every path/string. Unit tests cover Windows
   backslashes, quotes, non-ASCII, and a path that would break naive quoting.
5. Fake-server tests for 202 → queued → generating → completed, failed, timeout, 403,
   415, 404 output, and malformed result.

### Slice B — workflow runner and stage admission

Files:

- `app-image-workflow-run.go` (new)
- `app-image-workflow-run_test.go` (new)
- `app-image-model.go` (delta-request helper, startup ControlNet field)
- `internal/image/types/types.go` (`ImageWorkflowOptions`, stage result)
- `internal/resource/admission.go` (`UpscalerBytes` additive field)

Work:

1. Public `RunImageWorkflow(opts)` acquires `imageStartMu` once; internal
   `runImageWorkflowLocked` performs the four stages.
2. Stage 1 calls the existing `ensureEngineReady` path with the workflow's ControlNet
   startup field, then `GenerateTxt2Img`.
3. Stage 2/4 call `RunDGL` with the upscaler script; admission is a delta request using
   `UpscalerBytes` and current free VRAM.
4. Stage 3 calls `GenerateImg2Img` with the stage-2 output and pose image.
5. Write every output under the isolated workflow temp directory and copy the final
   image to the gallery using the existing gallery path.
6. After each stage, record PID, working-set footprint, free VRAM, reserve, and DGL job
   ID. At terminal/failure, call `reclaimImageEngineAfterGeneration` once.
7. Unit tests use a fake engine/client; no real model.

### Slice C — server-start arguments and restart detection

Files:

- `internal/image/engine/sdcpp/engine.go` (`buildServerArgs`)
- `app-image-model.go` (`ModelLoadOptions` population / `imageStartupOptionsDiffer`)
- `internal/image/types/types.go` (`ModelLoadOptions.ControlNetPath`,
  `HiresUpscalersDir`)

Work:

1. Emit `--control-net <path>` from the workflow's ControlNet when the model bundle does
   not already provide one. Include it in startup-difference checks so a changed
   ControlNet is a real restart, not an inert control.
2. Emit `--hires-upscalers-dir <dir>` for the workflow's upscaler directory, so DGL's
   `load_upscaler` model read passes the allowlist.
3. Charge the ControlNet exactly once in `imageAdmissionRequest`: prefer the explicit
   path, then the startup path, then the existing bound-model glob. Do not add the
   glob estimate on top of an explicit file. This is a correctness fix supported by the
   current code path, not a redesign.
4. Tests assert exact argument lists and restart behavior.

### Slice D — placement/report and UI

Files:

- `app-image-placement.go` (additive `workflow` section)
- `frontend/dist/image-generator.html` (workflow panel, existing controls untouched)
- `frontend/dist/js/image-workflow-run.js` (new)
- `frontend/dist/css/ig-controls.css` (panel styling)

Work:

1. Add a "Sequential workflow" panel with base prompt/size/steps, upscaler path, pose
   path/editor handoff, ControlNet path, refine size/strength, and Run.
2. Show one row per stage: status, duration, PID, working set, free VRAM, output
   thumbnail. Show the existing placement card and reclaim summary alongside it.
3. Keep the current single-generation controls and recipe save/load untouched.
4. The panel is a workflow runner, not a new model manager; no new persistence format
   is required for the PoC.
5. Label engine evidence by its actual scope (§9.14): INFO loaded bytes, VERBOSE
   planned target/workspace, DEBUG buffer release/eviction. Never present a planned
   target as actual component residency.

### Slice E — live verification (only after A–D pass unit/static checks)

Run the matrix in §6 with the isolated root
`C:/Users/StoneWolfPC/howl-image-hw-iso-20261124`, read-only model files, and the real
`sd-server.exe`. Record logs, CSVs, screenshots, and output PNGs under
`build/verification/workflow-2026-11-24/`.

---

## 6. Verification matrix

Status now: **not run**. Every row is a planned live check, not a claim.

### 6.1 Independent operations

| ID | Case | Procedure | Evidence | Pass criteria |
|---|---|---|---|---|
| I-1 | Base generation | native txt2img, 512×512, 4 steps, q8/auto | command line, PID, duration, PNG | image saved; admission logged; reserve held |
| I-2 | ESRGAN 2x | DGL job on a fixed 512 PNG | job ID, output shape 1024×1024, output dtype, duration, PNG | PNG byte-valid; output dimensions 2×; DGL status `completed`; if dtype is not U8, the script uses `image::to_u8` and the rerun passes |
| I-3 | ControlNet img2img | native img2img, 1024 input, prepared COCO-18 pose, startup `--control-net` | job ID, command line, controlnet load line, PNG | image saved; control image actually applied; no engine error. This variant has not been run before. |
| I-4 | Second upscale | DGL job on I-3 output | job ID, output shape 2048×2048, duration | PNG valid; dimensions 2×; no OOM |

### 6.2 Complete workflow and repetition

| ID | Case | Procedure | Evidence | Pass criteria |
|---|---|---|---|---|
| W-1 | Adequate VRAM | full four-stage workflow, auto cap, q8, VAE tiling as admitted | per-stage rows; one PID unless a justified recycle; free VRAM after each stage | all four stages complete; 512→1024→1024→2048; the preferred 10% reserve is held at every boundary; any documented soft dip stays within 500 MB and above the 128 MB floor |
| W-2 | Repeat 3× | run W-1 three times sequentially | PIDs; DGL job IDs; stage 2/4 durations; working set | same server PID across repeats; no base reload; upscaler reused; no reclaim in the normal band |
| W-3 | Batch 2 | `BatchCount=2` complete workflows | order, two final PNGs, peak VRAM/RAM, duration | sequential, never concurrent; peak equals one workflow's peak (±measurement noise). Fresh current-policy run: §9.13. |
| W-4 | Constrained 3 GiB | explicit 3 GiB + offload; existing bounded CPU-diffusion mode | actual command line; per-stage placement; peak total GPU; free VRAM | every stage either completes within the measured bounded envelope or the first impossible stage refuses with a resource reason; no OOM/crash/corruption; reserve held; slower time recorded |
| W-5 | CPU-only 2 GiB | `--max-vram 2 --backend cpu`; DGL upscaler backend `cpu` | command line; process working set; total GPU; stage durations | every stage either completes with no GPU weights or refuses cleanly with a resource reason; reserve held; time recorded |

### 6.3 Reuse, pressure, and recovery

| ID | Case | Procedure | Evidence | Pass criteria |
|---|---|---|---|---|
| R-1 | Component reuse | W-1/W-2 with TTL `-1` | DGL job logs, `load_upscaler` occurrences, server `Starting sd-server` count | one server start; upscaler loaded once, reused at stage 4; base model not reloaded between stages |
| R-2 | Necessary reclamation | inject measured growth (`HOWL_IMAGE_RECLAIM_TEST_GROWTH_MB`) or external VRAM pressure | reclaim event, PID before/after, working set, free VRAM, reason | exactly one recycle when the band/reserve is crossed; next stage succeeds; no second reclamation path |
| R-3 | Reserve pressure | external `hold_vram.exe` allocation leaving less than the stage's projection after the preferred reserve | live snapshot, admission decision, reclaim event, post-reclaim free VRAM | stage either fits through the bounded 500 MB soft allowance with at least 128 MB free, or refuses cleanly; no orphan process |
| R-4 | Failure recovery | kill `sd-server` between stages; invalid upscaler path; missing pose | error, PID, logs, next-run result | failed stage is reported; existing owner sees the failure; next workflow starts a clean server and completes |
| R-5 | DGL contract negatives | foreign `Origin`, wrong content type, empty source, bad output index | HTTP status and body | 403 / 415 / 400 / 404 respectively; no job or partial artifact left behind |

### 6.4 Hosts

| ID | Case | Procedure | Evidence | Pass criteria |
|---|---|---|---|---|
| H-1 | Browser harness | `howl-ui-harness.exe`, real bindings, workflow panel | stage rows, placement card, PNG, harness log | all four stages visible; no console errors; output appears |
| H-2 | Packaged WebView2 | `howl-chat-s1.exe`, native clicks | screenshot, log, actual command line, output PNG | same visible workflow; no host-only failure |

### 6.5 Measurements to record for every live row

- Actual OS command line of `sd-server` (not the planned flags).
- PID at each stage; any PID change.
- `nvidia-smi` total used/free at 0.5 s intervals; desktop baseline for process estimate.
- `sd-server` working set after load and after each stage.
- DGL job ID, `created`/`started`/`completed`, progress stage, output shape/bytes.
- Native job ID/status for stages 1 and 3.
- Admission line: admitted, required, activation, reserve, available, warning.
- Reclaim line: baseline, band, measured, growth, reason, PID before/after, free before/after.
- Output artifact paths and dimensions; SHA-256 of final PNGs.
- Engine evidence with its level: `Diffusion/Conditioner/VAE/ControlNet weights` and
  `Upscaler architecture` at INFO; `auto-fit plan` params target, `compute buffer
  size`, and `free compute buffer` at VERBOSE; `model manager evicted/released` at
  DEBUG. Absence of an unload line is not evidence at INFO; see §9.14.

---

## 7. Risks and open items

1. **8 GB at 1024 + ControlNet is the pressure point.** The base q8 model plus the
   OpenPose ControlNet may not fit concurrently at 1024 without CPU-diffusion and VAE
   tiling. This is exactly what W-4 must measure. The plan does not silently reduce
   resolution; a constrained 768 refine would be a separate declared case.
2. **Upscaler asset registration.** The registry has no upscaler type. The workflow
   needs an explicit path and a way to expose it. The PoC can use a workflow-only field;
   a general upscaler registry is out of scope.
3. **ControlNet startup slot and img2img variant.** The existing native path binds
   ControlNet at server start, and the proven run was txt2img. Adding a per-workflow
   startup ControlNet is a small change, but it means the ControlNet is resident for
   stages 1–4. Admission must charge it once, I-3 must prove the img2img variant, and
   W-1 must show whether that residency is acceptable at 512 and 1024.
4. **DGL upscaler path allowlist.** If the upscaler is outside
   `--hires-upscalers-dir` and all model roots, `load_upscaler` fails with a file-access
   error. The server arg must be added before the job, not after a failure.
5. **DGL result retention.** Completed jobs expire after 600 s. The client fetches the
   PNG immediately; it does not treat the job as durable storage.
6. **Engine residency evidence is partial and per-component.** INFO reports loaded
   bytes for the diffusion model, conditioner, VAE and ControlNet, plus the upscaler
   architecture; it does not report their device. VERBOSE reports each DGL resource's
   planned params target and each runner's actual compute workspace; DEBUG reports
   buffer release/eviction. There is no live component residency map or actual
   per-component weights VRAM/RAM split. See §9.14. The report must not collapse these
   into one "confirmed placement" claim.
7. **TTL interaction.** A user-set residency TTL can unload the DGL upscaler between
   stages. The default `-1` is the reuse case. If TTL is set, record the stage-4 reload
   instead of overriding it.
8. **No dedicated DGL cancel handler in `3b12781`.** Do not assume
   `/dgl/jobs/{id}/cancel`. Treat a client timeout as an abandoned job and use the
   existing engine recovery path. Verify the generic cancel endpoint before using it.
9. **Pose extraction remains external.** The PoC's pose stage is allowed to use a
   prepared skeleton. A photo-to-pose claim requires an explicitly installed
   OpenPose/DWPose processor and its own verification row.
10. **Admission double-count risk.** `imageAdmissionRequest` currently adds
    `gen.ControlNetPath` and `boundControlNetBytes`. A workflow that sets both would
    charge the ControlNet twice. The implementation must charge exactly once before
    W-1 is considered valid.

---

## 8. Authorization request

No code has been changed and no image-generation experiment has been run. The next
phase would:

1. implement Slices A–D (DGL client, workflow runner, startup args, panel) with tests;
2. run the static DGL script check and fake-server tests;
3. run the live matrix in §6 on the GTX 1080, starting with I-1 through I-4 and W-1.

Authorize implementation and live verification only if this plan, the hybrid
native+DGL choice, and the explicit external-pose/upscaler dependencies are acceptable.

---

## 9. Implementation and measured verification (2026-11-24)

Authorization was granted after section 8. Slices A–D were implemented and the live
matrix below was run on the GTX 1080 8 GB. This section supersedes section 8's
"no code changed" statement; it does not rewrite the research record.

### 9.1 Owners (unchanged)

- **Admission:** `resource.Authority` via `Admit`. `imageAdmissionRequest` now charges
  the ControlNet once and accepts an `UpscalerBytes` delta for DGL stages.
- **Reclamation:** `app-image-reclaim.go` remains the sole owner. The workflow calls
  the existing locked helper once at terminal outcome; DGL is an executor only.
- **Serialization:** the workflow holds `imageStartMu` for the whole sequence, so no
  two GPU stages can run concurrently.
- **No competing reclamation mechanism was added.** No second admission authority,
  component cache, or post-workflow unload path was introduced.

### 9.2 What was implemented

- `internal/image/engine/sdcpp/dgl.go`: DGL submit/poll/output, script builder,
  structured report, PNG dimensions.
- `client.go`: tolerant job-progress shape, synchronous `/sdapi/v1/img2img` mode.
- `engine.go`: ControlNet startup slot, upscaler roots, `GenerateImg2ImgSync`.
- `app-image-workflow-run.go`: four-stage runner, batch loop, delta admission,
  per-stage evidence, terminal reclaim.
- `app-image-model.go` / `app-image-reclaim.go`: locked ensure-ready refactor,
  ControlNet charge-once fix, reclaim wrapper.
- `frontend/dist/image-generator.html`, `js/image-workflows.js`, `css/ig-pose.css`:
  Sequential Pipeline panel, stage report, pose path, batch pass count.
- Tests: `dgl_test.go` (DGL happy path, progress shapes, sync steps, failed-job
  settlement), `internal/resource/upscaler_test.go` (upscaler admission).

### 9.3 Verification matrix

| Operation | How exercised | Result | Evidence |
|---|---|---|---|
| Base txt2img | Stage 1 of the full workflow, both hosts | **Verified** | stage 1 `sdcpp-native`, 512×512 |
| DGL ESRGAN upscale | Stage 2 and standalone DGL probe | **Verified** | DGL job `completed`, 1024×1024 PNG fetched |
| ControlNet img2img | Stage 3 of the full workflow, native async route with a supplied COCO pose PNG | **Verified after correcting a false-positive** | ControlNet loaded once (`ControlNet weights: 1415.82 MiB`), applied, 4/4 requested steps, 512×512; see §9.8. The earlier sync-route claim in §9.6 was invalid: it produced the identical image with and without a control image. |
| Second upscale | Stage 4, same upscaler, same process | **Verified** | DGL job `completed`, 1024×1024 |
| Pose extraction | External OpenPose/DWPose only | **Not installed / not claimed** | ControlNet conditions on a supplied skeleton; it does not extract pose |
| HTTP `/sdcpp/v1/upscale` | — | **Incompatible with installed `3b12781`** | route absent; DGL used instead |
| `RealESRGAN_x2plus.pth` | Stage 2 | **Unsupported by this build** | `unsupported image upscaler weights`; use GGUF |
| `2x-ESRGAN.gguf` | Stages 2 and 4 | **Verified** | 64 MB GGUF, both upscales completed |
| Batch > 1 (512×512) | Browser harness, `BatchCount=2`, fresh server, revised admission policy, ControlNet CPU, ESRGAN DGL | **Verified fresh under the revised policy** | both items completed; same PID 8,760; one server start, one upscaler architecture/ControlNet weight line and no component reload; unload absence is not observable at INFO (§9.14); boundary and terminal reclaims both `retain`; item paths/hashes and stage admissions in §9.13. |
| Engine placement/memory evidence | Source audit of pinned `3b12781`, existing verbose residency log, and the bounded DEBUG batch capture | **Verified, scoped per component** | INFO loaded bytes for diffusion/conditioner/VAE/ControlNet and upscaler architecture; VERBOSE planned target and per-runner workspace; DEBUG observed all five runner workspaces, zero instrumented weight-release/eviction/unload lines, and no reload across two passes; release lines themselves are not component-named; no live residency map and no actual per-component weights VRAM/RAM split; §9.14. |
| VAE tiling | Stage 1 of the final browser and packaged runs | **Verified** | `--vae-tiling` selected by admission; base peaks 6,095/1,965 and 6,225/1,835; reserve held. |
| ControlNet CPU placement | Stage 3 of the final browser and packaged runs | **Verified** | `--backend controlnet=cpu`; `ControlNet weights: 1415.82 MiB`; stage peaks 6,223/1,837 and 6,273/1,787; no CUDA staging; see §9.9. |
| Reserve across all four stages | Final browser and packaged runs | **Verified** | lowest free 1,787 MB (packaged), above the 819 MB reserve; no component reload; unload absence is not observable at INFO (§9.14); one PID each. |
| Soft reserve allowance | Browser harness with external 1,620 MB holder | **Verified** | 212 MB of a 414 MB allowance used, 330 MB reported residual; measured `total - used` within 14 MB; §9.10. |
| Packaged WebView2 revised policy | Full four-stage packaged run, 7,139 MB startup free | **Verified** | one PID 15,292, ControlNet CPU, all stages above preferred reserve, same stage-4 hash; §9.11. |
| Failure recovery, fresh policy | Existing but invalid `invalid-upscaler.gguf`; stage 2 DGL failure; restore valid GGUF and rerun twice | **Verified fresh** | stage 2 reported `upscaler loading failed` after 679 ms; reclaim retained PID 21,608; next two full runs succeeded with hashes `5165011…e535`; §9.12. |
| Consecutive-run reuse | Two complete four-stage runs back-to-back on the failed run's surviving server | **Verified fresh** | one `Starting sd-server`, one `Upscaler architecture: esrgan`, one `ControlNet weights`, no component reload, same PID 21,608; unload absence is not observable at INFO (§9.14); stage timings sequential; §9.12. |

### 9.4 Measured runs

**Packaged WebView2, final build, one pass (PID 15004):**

| Stage | Executor | Duration | Output | Admission |
|---|---|---:|---|---|
| 1 base | native | 11,914 ms | 512×512 | delta admission after startup |
| 2 upscale | DGL | 24,310 ms | 1024×1024 | required 2111 MB, free 2596 MB, reserve 259 MB |
| 3 pose/ControlNet | native img2img | 4,902 ms | 512×512 | required 2048 MB, free 2564 MB, reserve 256 MB |
| 4 final upscale | DGL | 24,315 ms | 1024×1024 | required 2048 MB, free 3392 MB, reserve 339 MB |

- **Pre-correction note:** this table's stage 3 used the sync route, which did not
  load ControlNet. Its timings and completion remain measurements of the old path;
  they are not ControlNet evidence. The corrected stage 3 is measured in §9.8.
- Final PNG: 1024×1024, displayed in the packaged app with `Mode: workflow`,
  `Steps: 4`.
- Reclaim: `decision=retain`, `growth=0MB`, baseline/measured 423 MB,
  free 3398 MB, reserve 819 MB, `pid=15004->15004`.
- Startup: one `sd-server` with `--control-net`, `--model-dir` and
  `--hires-upscalers-dir`; q8 fallback after the initial f16 request was refused.

**Browser dev host, same build, one pass:** four stages completed, stage 1 4912 ms,
stage 2 23,297 ms, final 1024×1024; reclaim `retain`, `pid=24824->24824`.

**Browser failure then recovery:** stage 2 failed with
`generation_failed: line 10, column 16: upscaler loading failed`; the pipeline showed
the failure and stopped. Replacing the `.pth` with the working GGUF succeeded without
recycling the process (`pid=24824->24824`), proving a definite DGL failure is not
treated as an abandoned engine.

**Batch:** browser 1m49.9s for two passes, both final PNGs present; packaged app also
produced `run-1`/`run-2` stage-4 PNGs with the same PID retained.

### 9.5 Verified limitations

1. **Reserve pressure: verified (see §9.6).** A real CUDA allocation held 157 MB while
   stage 2 was admitted. Stage 2 was refused at 2111 MB required, 2339 MB free, 233 MB
   reserve, 2106 MB available — free VRAM alone was sufficient, so the 10% reserve was
   the deciding constraint. The server PID did not change and no second engine started.
2. **ControlNet on the compatibility route: not applied (corrected 2026-11-24).** The
   temporary probe and the workflow stage used `/sdapi/v1/img2img`; that route accepts
   `control_net_path`/`control_image` without error but never loads or applies the
   ControlNet. The proof is a byte-identical output hash with and without a control
   image (`ca5e4f3b…beaee`). The probe's 6-step result was therefore an img2img result,
   not a ControlNet result. Stage 3 now uses the native async route and is verified in
   §9.8; the old probe remains removed.
3. **Engine residency evidence is partial, not absent.** The engine reports loaded
   bytes per major component at INFO, planned params targets and per-runner workspace
   at VERBOSE, and buffer release at DEBUG. It does not expose a live component map or
   an actual per-component weights VRAM/RAM split; the UI must label that part
   unverified. OS working set and live `nvidia-smi` remain the independent check. See
   §9.14.
4. **DGL upscaler reuse is now component-evidenced.** Four DGL jobs (stages 2 and 4 in
   each of two passes) produced one `Upscaler architecture: esrgan` load line and no
   unload. The per-job `load_upscaler` line is the auto-fit plan phase, not a disk
   reload. The 31.80 MiB adapter stayed resident on CUDA0; each job allocated and then
   freed its 104 MB compute buffer. See §9.8.
5. **Async `/sdcpp/v1/img_gen` ignores the flat `steps` field** (8 requested, 20
   sampled), but honours `sample_params.sample_steps` (4 requested, 4 sampled). The
   workflow now uses the native route with nested sampling fields and base64
   `control_image`; the sync route is no longer used for stage 3. This is an API-shape
   limitation, not a UI choice.
6. **Pose extraction remains external.** A photo-to-pose claim would require an
   installed OpenPose/DWPose processor and its own verification row.
7. **512×512 two-pass batch is verified under the current policy.** The fresh
   browser-harness run in §9.13 completed both items with ControlNet active, one
   server start, one upscaler/ControlNet load each, no component reload, and `retain`
   at both the item boundary and the terminal decision. Its INFO-only log cannot prove
   zero unloads because the engine's unload lines are VERBOSE (§9.14). The earlier five
   packaged runs used the sync route and ignored ControlNet; §9.7 remains
   completion/residency evidence for that older route, not the current-policy proof.

### 9.6 Remaining-gap closure (2026-11-24)

#### Reserve-pressure refusal and recovery

The reserve-pressure experiment used a temporary verification-only pause after
stage 1 so pressure could be applied between stages. That hook has been removed
from the delivered source and is absent from the current normal and packaged
binaries; no admission or reclamation code was bypassed.

- Baseline before pressure: **free 2463 MB**.
- Pressure process: a temporary CUDA driver allocation requested 157 MB; live free
  measured **2198 MB**, then **2173 MB** while the process was held. No Howl code
  allocated or released this memory, and no component was unloaded.
- Stage 2 admission (packaged WebView2, server PID 24272):
  **required 2111 MB, free 2339 MB, reserve 233 MB, available 2106 MB**.
  The request failed by 5 MB *after* the reserve; raw free was 228 MB above the
  requirement, so the refusal is attributable to the 10% reserve, not generic
  insufficiency.
- The UI showed `Pipeline failed: stage 2 upscale: ... needs 2111 MB VRAM; 2339 MB
  free ... 233 MB reserve (10%), 2106 MB available; blocking model ...`.
  Stage 1 completed; stage 2 failed in 83 ms with `admitted=false`.
- No unload/offload occurred: the log contains one `Starting sd-server`, the PID was
  `24272 -> 24272`, and reclamation recorded `decision=retain`.
- After release: free **2438 MB** before rerun; the same stage-2 operation was then
  admitted at **2111 MB required, 3493 MB free, 349 MB reserve, 3144 MB available**.
  The four-stage sequence completed, PID `24272 -> 24272`, final
  1024×1024 PNG:
  `C:\Users\StoneWolfPC\howl-image-hw-iso-20261124\temp\workflow-20261124_184537\run-1-stage-4.png`,
  SHA-256 `4933453faf6961a5694a9f1f48a79ed064ed5d8b84117d8c430513937850715e`.
- Evidence: `build/verification/workflow-2026-11-24/reserve-refusal.png`,
  `reserve-released-success.png`, `reserve-pressure-measurements.txt`,
  `desktop-pressure2.log`, `pressure-target.log`.

#### Standalone ControlNet-conditioned img2img — invalidated and replaced

The temporary `RunControlNetProbe` result recorded here is **not valid evidence that
ControlNet was applied**. The probe used `GenerateImg2ImgSync`
(`/sdapi/v1/img2img`), and a direct A/B test on the installed server proved that route
silently ignores the ControlNet fields: the same seed, init image and request with and
without `control_net_path`/`control_image` produced byte-identical PNGs, SHA-256
`ca5e4f3be91c97e1f521adacb2da74a5bf324aa12edccb3dec1f2a472b4beaee` (585,776 bytes).
The 6-step result and its hash therefore describe an unconditioned img2img run. The
probe remains removed from the delivered source; the corrected stage-3 route is in
§9.8. No conclusion about ControlNet efficacy should cite this subsection.

#### Packaged 2-pass batch label

- First replay: packaged WebView2, `ig-batch=2`, 256 base/refine, 2/1 steps.
- Second replay: packaged WebView2, `ig-batch=2`, 512×256 base/refine, 4 steps.
- Visible result in both: **`Pipeline complete (2 passes). Weights retained unless
  reclamation reported otherwise.`**
- Artifacts: `run-1-stage-4.png` and `run-2-stage-4.png` in
  `workflow-20261124_185628`; the second has SHA-256
  `28d33227a1e3ffbce6e842745d6131d928f125208c846695568e5cba50c0b7fd`.
  The 512×256 replay is in `workflow-20261124_190112`; its `run-2-stage-4.png`
  SHA-256 is `e8817c3ee06d61f6a3282ca7b48b3468ee36856b2e0165046b50f6c763d7fad4`,
  and reclamation retained PID `2364 -> 2364` with 18 MB growth and 3421 MB free.
- Evidence: `build/verification/workflow-2026-11-24/batch-2-pass-label.png`,
  `batch-2-pass-label-512.png`, and the 512×512 runs in §9.7.

### 9.7 512×512 two-pass batch reliability and delivered-tree cleanup

Five full 512×512 two-pass packaged-WebView2 batches completed after the earlier
single stall. Runs A–C were consecutive on one server; D used a freshly rebuilt
packaged binary; E followed a 256×256 two-pass batch on the same server (the
resolution transition that preceded the original stall). Every run used base and
refine 512×512, 4 steps, `ig-batch=2`, seed 9001/9002, and the same ControlNet and
ESRGAN assets. **Those runs used the sync stage-3 route, which ignored the ControlNet;
 they prove completion, batching and residency, not ControlNet conditioning.**

| Run | GPU before / after (used/free MB) | Stage durations (base, upscale, pose/CN, final) | Admission (stage 2 / 3 / 4) | Server | Artifact and SHA-256 |
|---|---|---|---|---|---|
| A, 19:06Z | 1055/7005 → 4844/3216 | 4636, 24333, 4661, 24305 ms | 2111 MB @3382 free/338; 2048 @3350/335; 2048 @3351/335 | 12328→12328 | `workflow-20261124_190650/run-2-stage-4.png`, `28bb2d84…58ac` |
| B, 19:09Z | 4842/3218 → 4844/3216 | 4623, 23800, 4697, 23809 ms | 2111 @3382/338; 2048 @3381/338; 2048 @3381/338 | 12328→12328 | `workflow-20261124_190944/run-2-stage-4.png`, `28bb2d84…58ac` |
| C, 19:12Z | 4841/3219 → 4818/3242 | 4618, 23803, 4671, 23805 ms | 2111 @3382/338; 2048 @3381/338; 2048 @3381/338 | 12328→12328 | `workflow-20261124_191212/run-2-stage-4.png`, `28bb2d84…58ac` |
| D, 19:17Z (rebuilt) | 1056/7004 → 4891/3169 | 4617, 24297, 4710, 25808 ms | 2111 @3390/339; 2048 @3390/339; 2048 @3288/328 | 20116→20116 | `workflow-20261124_191716/run-2-stage-4.png`, `28bb2d84…58ac` |
| E, 19:22Z (after 256) | 1104/6956 → 5689/2371; 512 pre-run 5812/2248 → 4918/3142 | 4816, 24809, 4831, 24815 ms | 2111 @3303/330; 2048 @3307/330; 2048 @3307/330 | 22540→22540 | `workflow-20261124_192249/run-2-stage-4.png`, `28bb2d84…58ac` |

All five hashes are identical because the runs use the same seed and deterministic
sampling (`28bb2d8438899869496bcc51ba3e216ff36a10f9ba14556dcada6f8fc32758ac`);
the distinct workflow directories and both pass artifacts prove separate runs.
`ImageReclaim` recorded `decision=retain` each time. Server working sets after the
runs were 450–477 MB; the packaged app was 77–84 MB. Run E's 256×256 pass and the
following 512×512 pass shared PID `22540` and one `Starting sd-server`, so changing
resolution did not force an unload.

Evidence: `batch512-evidence.txt`, `batch512-runA-ui.txt` through
`batch512-runE-switch-ui.txt`, `desktop-batch512-a.log`,
`desktop-batch512-final.log`, `desktop-res-switch.log`, `batch-512-two-pass-c.png`,
and `batch-512-final.png`.

The original stall log stopped after `sampling completed` for pass-2 stage 3 and
contains no following `computing vae decode graph` line, at 7.6 GB used / 474 MB
free. Its own startup admission line records 7071 MB free, so it was not a leftover
pressure allocation at launch. The five replays all logged the VAE decode and held
3.1–3.3 GB free at the equivalent point, and no pressure-holding CUDA process,
Python process, or second `sd-server` was present in the delivered state. The stall is therefore treated as a one-off low-free-memory event, not a
deterministic workflow defect: the batch path was not changed, and five full
512×512 replays, including the exact 256→512 transition, all completed.

Delivered-tree check after removing the pause hook and probe source:

- No `howl-chat-s1`, `sd-server`, `python`, or `pythonw` process remains.
- `nvidia-smi` shows no Python CUDA holder; desktop GPU use returned to 963 MB used
  / 7097 MB free.
- `grep` over `build/bin/howl-chat-s1.exe` and `build/bin/howl-ui-harness.exe` finds
  none of `HOWL_WORKFLOW_TEST_PAUSE_MS`, `RunControlNetProbe`,
  `controlnet standalone`, or `verification pause`.
- The stale `howl-ui-harness-verify.exe` and the temporary CUDA pressure scripts were
  deleted.

### 9.8 Component residency, DGL reuse, and the corrected ControlNet route (2026-11-24)

#### Defect and correction

Stage 3 previously used `GenerateImg2ImgSync` (`/sdapi/v1/img2img`) because the async
route was believed to ignore requested steps. That route also silently ignores
ControlNet: an A/B request with the same seed and init image, one with
`control_net_path` + `control_image` and one without either, produced the identical
PNG hash `ca5e4f3be91c97e1f521adacb2da74a5bf324aa12edccb3dec1f2a472b4beaee`.
The verbose engine log contained no ControlNet load line for either request.

The async route does work when the API shape is correct:

- `init_image_path` is supplied as a path;
- `control_net_path` is the model path and `control_image` is base64 PNG data
  (`control_image_path` was ignored in this request shape);
- sampling values are nested under `sample_params.sample_steps`,
  `sample_params.sample_method`, `sample_params.scheduler`, and
  `sample_params.guidance.txt_cfg`; the flat `steps` field is ignored.

The corrected stage 3 is `Engine.GenerateImg2ImgNativeControl`, using
`/sdcpp/v1/img_gen`, polling `/sdcpp/v1/jobs/{id}`, and decoding the `b64_json`
image objects. The workflow now calls that method. The compatibility route is
unchanged for callers that do not use ControlNet.

Unit coverage: `TestNativeImg2ImgCarriesControlAndNestedSteps` asserts the native
route, nested step/sampler/scheduler/cfg fields, base64 control image, normalized
paths, and `b64_json` decoding.

Live browser-host verification (rebuilt `howl-ui-harness.exe`, isolated root
`howl-image-hw-iso-20261124`, `ig-batch=2`, base/refine 512×512, 4 steps, seed
9001/9002):

- Stage 3 logged `ControlNet weights: 1415.82 MiB` and
  `sampling started: method=euler_a, steps=4` in both passes.
- Pass 1 stage 3: 6,604 ms, 512×512, SHA-256
  `856e87ef8559ef2df39bf79f26ca331ad8f1384356e7e07d619d11cb3408b08f`.
- Pass 2 stage 3: 512×512, SHA-256
  `3d76eb1730fefd3c3fcd91189c56f8343cc8ec46de38efb61322080592c87ca6`.
- The UI showed `Pipeline complete (2 passes). Weights retained unless reclamation
  reported otherwise.` All eight stage PNGs are present in
  `C:\Users\StoneWolfPC\howl-image-hw-iso-20261124\temp\workflow-20261124_194118`.
- External pose extraction is still not implemented; the stage consumes a supplied
  OpenPose skeleton.

#### Measured residency between stages

The installed server's verbose auto-fit plan exposes each DGL resource's **planned**
params target; the INFO `... weights` lines expose bytes actually registered, and the
verbose compute-buffer lines expose workspace. This is not a live component map:
actual per-component weights allocation is inferred from the plan plus the aggregate
buffer line. In the verbose run the consistent result is that weights stay resident and
only per-job compute workspaces are freed.

| Component | Planned target (V) | Loaded bytes / workspace evidence | Between stages |
|---|---|---|---|
| UNet | CUDA0 | `load_unet: 2910.44 MiB params -> CUDA0`; 1,680 tensors loaded once | Retained; no unload |
| Text encoder | CUDA0 | `load_text_encoder: 925.92 MiB params -> CUDA0`; loaded once | Retained; no unload |
| VAE | CUDA0 | `load_vae: 159.68 MiB params -> CUDA0`; 140 tensors loaded once | Retained; no unload |
| ControlNet (OpenPoseXL2) | weights on CPU, compute on CUDA0 | `ControlNet weights: 1415.82 MiB` once; `controlnet::load: 1415.82 MiB params -> CPU`; 122.42 MB VRAM + 3.67 MB RAM compute buffers freed after the job | Weights retained in host RAM; compute freed |
| ESRGAN 2x | CUDA0 | `Upscaler architecture: esrgan` once; 31.80 MiB params; 104.00 MB VRAM + 0.19 MB RAM compute buffers freed after each job | Adapter retained; compute freed |

No `unload` line appeared for any component between the four stages. The
`load_upscaler` and `controlnet::load` lines that repeat inside each job are
auto-fit plan phases, not disk reloads: the one-time architecture/weights lines are
the actual load evidence.

#### DGL reuse

Two passes, four DGL jobs (stage 2 and stage 4 in each), produced exactly one
`Upscaler architecture: esrgan` line and no unload. Each job re-created only its
104 MB compute workspace and freed it:

```
[I] Upscaler architecture: esrgan, weight type: f16
[I] upscaling from (512 x 512) to (1024 x 1024)
[I] input_image_tensor upscaled, taking 22.74s
...
[V] free compute buffer: 104.00 MB(VRAM) on CUDA0
```

This is real DGL cache reuse: the adapter is loaded once per server process and
shared by both upscale stages and both batch passes. It is not a persistent DGL
script context; each job still runs its own `load_upscaler` plan phase, which does
not reload the already-resident weights.

#### Actual memory, and a reserve limitation

Live `nvidia-smi` sampling at 0.5 s intervals during the two-pass run:

| Point | Used / free MB |
|---|---|
| Before workflow | 1,087 / 6,973 |
| Pass-1 base peak (sampling/decode) | **7,532 / 528** |
| Pass-1 stage 3 ControlNet peak | 7,334 / 726 |
| Pass-1 stage 4 DGL peak | 7,092 / 968 |
| After completion, weights resident | 4,763 / 3,297 |

The server process working set after the run was 1,977,004,032 bytes (~1,885 MB
reported by reclamation); the engine PID remained 25000 from first start to
completion, and reclamation recorded `decision=retain` with `growth=0MB`.

**The 10% free-VRAM reserve was not held at the pass-1 peak.** The peak left 528 MB
free, below the 819 MB reserve. Admission had recorded
`required 5084 MB, 6144 MB available after 814 MB reserve`, so the plan's decode
activation estimate (551 MB) under-predicted the actual peak. This is a real
resource-policy gap, not a residency failure: no OOM occurred and no component was
unloaded, but Howl must not claim the reserve is guaranteed until the peak decode /
load-staging estimate is corrected and re-measured. The peak occurred in the base
stage, before ControlNet was loaded. **Superseded by §9.9:** the activation estimate,
VAE tiling and ControlNet placement were corrected and the reserve was held in both
hosts.

#### Host scope

Verified in the browser dev host with the real backend over HTTP. The packaged
WebView2 binary was rebuilt but **not re-driven after the stage-3 route change**;
the previous five packaged runs exercised the old sync route and therefore do not
verify ControlNet application. Backend code and `frontend/dist` are shared, but the
packaged host still needs a ControlNet-run confirmation. **Closed in §9.9:** the
rebuilt packaged binary was re-driven and logged the ControlNet load and four steps.

Evidence: `harness-control-fixed.log`, `vram-control-fixed.tsv`,
`residency-verbose.log`, `residency-async-control.js`, `residency-dgl.js`, and the
eight PNGs in `workflow-20261124_194118`.

### 9.9 Reserve-safe correction and final host proof (2026-11-24)

#### What was wrong

The 551 MB decode activation was not a measurement of the installed build. A verbose
base-stage run at 512×512 showed **1,920 MB** of VAE compute, and the untiled base peak
was **7,579 MB used / 481 MB free** — a direct reserve breach. Measurement gave
**7,325 MB/MP** for VAE decode (480 MB at 256², 1,920 MB at 512², 4,320 MB at 768²).
The old estimate also missed that sd.cpp stages ControlNet params onto CUDA0 during
sampling even though `controlnet::load` reports CPU: **1,415.92 MB, 844 tensors**.
That made the stage-3 peak **7,339 MB used / 721 MB free**.

#### Correction (single owners preserved)

- `sdActivationMBPerMegapixel = 7325`; tiled decode is **480 MB**; fixed engine
  overhead is **512 MB**.
- `Admit` now projects resident converted weights + selected activation + fixed
  overhead. If the untiled projection exceeds the post-reserve budget and the tiled
  projection fits, it plans `--vae-tiling`.
- If staged ControlNet would cross the post-reserve budget, it plans
  `--backend controlnet=cpu`; diffusion stays on GPU. The measured CPU-ControlNet
  stage peak was **6,099 / 1,961** and the verbose log showed
  `controlnet::load: 1415.82 MiB params -> CPU` with no CUDA staging.
- Workflow deltas with `ResidentLoaded` no longer recharge the 2,048 MB module
  reserve, and they now also charge the 10% reserve (a zero `RequiredMB` no longer
  disables it). Stage 4 had previously refused at free 2,003 MB.
- The reserve percentage, `resource.Authority` and `app-image-reclaim.go` were not
  changed; no second reclamation path was added.

Startup admission before the final browser run:

```
[App] image admission: admitted=false activation=1920 MB ... needs 10478 MB VRAM;
  7170 MB free ... 6453 MB available; projected peak 10478 MB = resident 6615 MB +
  activation 1920 MB + 512 MB engine overhead
[App] image hardware: VAE tiling planned for 512x512: tiled projection 5093 MB fits
[App] image hardware: ControlNet CPU placement planned: staging 1431 MB ...
[App] image admission final: admitted=true projected=5093MB resident=4101MB
  activation=480MB required=5084MB tiling=true controlnet_cpu=true free=7058MB
  reserve=705MB
```

The final startup command line contained both selected placements:

```
--max-vram 6 --backend controlnet=cpu --type q8_0 --vae-tiling --flash-attn auto
```

#### Browser harness, final build, one pass

`howl-ui-harness.exe`, isolated root `howl-image-hw-iso-20261124`, batch 1, base and
refine 512×512, 4 steps, supplied COCO-18 pose, one server PID 22376 from first start
to completion. `nvidia-smi` sampled at 0.1 s; peaks below are bounded by the four
stage PNG completion times.

| Stage | Executor | Duration | Peak used / free MB | Admission |
|---|---|---:|---|---|
| 1 base | native | 15,050 ms | **6,095 / 1,965** | startup projection 5,093 MB, reserve 705 MB, tiling + ControlNet CPU |
| 2 upscale | DGL | 26,344 ms | **6,095 / 1,965** | projected 512 MB, free 2,577, reserve 257, available 2,320 |
| 3 pose/ControlNet | native img2img | 64,659 ms | **6,223 / 1,837** | projected 992 MB, free 2,535, reserve 253, available 2,282 |
| 4 final upscale | DGL | 23,878 ms | **5,822 / 2,238** | projected 512 MB, free 2,480, reserve 248, available 2,232 |

The used/free values are device totals and include the existing desktop baseline
(~1,134 MB used). The workflow-attributable drop was **7,058 − 1,965 = 5,093 MB**,
exactly the projected peak. The packaged run's drop was **6,991 − 1,835 = 5,156 MB**,
63 MB (1.2%) above projection. Stage 3's measured free drop was 698 MB against a
992 MB projection, so its delta admission was conservative.

- One `Starting sd-server`; one `ControlNet weights: 1415.82 MiB`; one
  `Upscaler architecture: esrgan`; two `sampling started` (stages 1 and 3), both
  `steps=4`; zero `unload`/`offload` lines; one `ImageReclaim` line with
  `decision=retain`, `baseline=1941MB band=1572MB measured=1941MB growth=0MB`,
  `pid=22376->22376`.
- Four PNGs in `workflow-20261124_203217`, stage 4 SHA-256
  `f6c95bd05e8838a5f6befdf05b669172fb3e20fdf77dd7e25d723b8dc87c82e0`.
- Evidence: `harness-final.log`, `final-vram.tsv` (browser portion ends 14:34:28).

#### Packaged WebView2, final build, one pass

`build/bin/howl-chat-s1.exe` was rebuilt and driven through its native window with
real mouse/keyboard input. The same isolated root and workflow inputs were used;
Settings were applied at startup, not by a test hook.

| Stage | Executor | Duration | Peak used / free MB |
|---|---|---:|---|
| 1 base | native | 15,679 ms | **6,225 / 1,835** |
| 2 upscale | DGL | ~27,460 ms | **5,956 / 2,104** |
| 3 pose/ControlNet | native img2img | ~65,770 ms | **6,273 / 1,787** |
| 4 final upscale | DGL | ~26,480 ms | **5,893 / 2,167** |

- Startup admission: projected 5,093 MB (resident 4,101 + activation 480 + 512 MB
  overhead), free 6,991 MB, reserve 699 MB, available 6,292 MB, tiling and
  ControlNet CPU planned. Command line included `--backend controlnet=cpu` and
  `--vae-tiling`.
- One server PID 13872 from 14:40:00 to completion; one `Starting sd-server`; one
  `ControlNet weights: 1415.82 MiB`; `sampling started: ... steps=4`; no unload or
  offload; one `ImageReclaim` line `decision=retain` / `pid=13872->13872`.
- Four PNGs in `workflow-20261124_204000`; stage 4 SHA-256 equals the browser run
  `f6c95bd…c82e0`. UI showed `Pipeline complete. Weights retained unless reclamation
  reported otherwise.`
- Evidence: `desktop-final-cpu.log`, `desktop-final-cpu-vram.tsv`,
  `desktop-cpu-result.png`, `desktop-cpu-wf-filled.png`.

#### Result against the policy

The lowest sampled free VRAM was **1,787 MB**, above the 819 MB reserve, while
ControlNet weights stayed resident in host RAM and the ESRGAN adapter stayed resident
on CUDA0. The 512² stage-3 peak fell from 7,339/721 (CUDA-staged ControlNet) to
6,223/1,837 (CPU placement) in the browser and 6,273/1,787 in the packaged host. No
component was unloaded, no second reclamation path was introduced, and the workflow
completed all four stages in both hosts.

### 9.10 Soft-reserve live proof (2026-11-24)

The revised policy was exercised in the browser harness on the GTX 1080. An external
`hold_vram.exe` reservation was applied before the run and released immediately after
the final admission decision; the workflow itself used the normal UI button and the
normal admission/placement/reclamation path.

#### Four-stage run with a soft grant

- Live admission: free **5,423 MB**, preferred reserve **542 MB**, post-reserve
  available **4,881 MB**. The initial untiled request was refused; admission then
  planned **VAE tiling** and **ControlNet CPU placement**, producing a projected peak
  of **5,093 MB** (resident 4,101 + activation 480 + 512 MB overhead).
- The final grant used **212 MB of a 414 MB soft allowance** and reported **330 MB
  free** above the 128 MB floor. The sampler's `total - used` at the same decision
  read **5,409 MB free** (14 MB lower, inside the sampling interval), and the holder
  log recorded 5,565 MB before that read. This is a grant that could not have passed
  the preferred 542 MB reserve.
- All four stages completed: base native 14,892 ms; ESRGAN DGL 23,424 ms to 1024²;
  ControlNet/pose native 122,650 ms; second DGL 22,411 ms to 2048². Stage rows 2–4
  each fitted the preferred reserve (free 3,588 / 3,576 / 3,580 MB; reserve 358 / 357 /
  358 MB).
- One `Starting sd-server` (PID 18,916), one `ControlNet weights: 1415.82 MiB`, no
  unload/offload, one `[ImageReclaim] decision=retain pid=18916->18916`, and four PNGs
  in `workflow-20261124_213256`; stage 4 SHA-256
  `f6c95bd05e8838a5f6befdf05b669172fb3e20fdf77dd7e25d723b8dc87c82e0`.
- Evidence: `harness-soft-live.log`, `soft-live-vram.tsv`,
  `soft-live-holder.log`, `soft-live-admission.txt`, `soft-live-released.txt`.

#### Hard boundaries, measured

A verification-only Go probe called the production `resource.DefaultAuthority().Admit`
against the live hardware monitor; it added no planner or policy code. Both probes used
the real 512² image request shape and a real `hold_vram.exe` reservation.

| Case | Free | Reserve | Available | Projected / need | Result |
|---|---:|---:|---:|---:|---|
| Soft grant | 5,423 | 542 | 4,881 | 5,093 | admitted; 212 MB used, 330 MB residual |
| Over 500 MB | **6,951** | 695 | 6,256 | **6,801** | refused; 545 MB short of available, exactly 500 MB allowance, 150 MB would remain |
| 128 MB floor | **5,172** | 517 | 4,655 | **5,084** | refused; 429 MB short (inside 389 MB allowance), only **88 MB** would remain |

- Over-500 evidence: `soft-probe-over500.log`, `soft-probe-over500-holder.log`.
- Floor evidence: `soft-probe-floor.log`, `soft-probe-floor-holder.log`.
- The soft grant's reported free agrees with the live monitor's `total - used`; both
  hard refusals state the measured free, reserve, available and projected peak in the
  refusal text.

After the proof, the harness, `sd-server`, sampler and VRAM holders were stopped; the
card returned to its idle baseline. No commit was made.

### 9.11 Packaged WebView2 under the revised policy (2026-11-24)

The packaged desktop build (`howl-chat-s1.exe`, WebView2) was driven through the real
window with the same ESRGAN/OpenPose/pose inputs, 512×512 base and refine, 4 steps.
No pressure holder was used for this host run.

- Startup admission: free **7,139 MB**, reserve **713 MB**, available **6,426 MB**;
  initial untiled request refused, then tiling and ControlNet CPU planned and the
  final 5,093 MB projection admitted. Command line: `--max-vram 6`,
  `--backend controlnet=cpu`, `--type q8_0`, `--vae-tiling`.
- One `Starting sd-server`, one ControlNet load (`1415.82 MiB`), server PID 15,292;
  no component reload. The INFO log cannot show unload/offload; see §9.14. Reclaim:
  `decision=retain`, `pid=15292->15292`, growth 0.
- Stage rows: base **13,731 ms**; DGL upscale **25,438 ms**; pose/ControlNet
  **135,802 ms**; final DGL **24,824 ms**. Stage 3 admission read free 2,600 MB,
  reserve 260 MB, projected 992 MB, available 2,340 MB; stage 4 read free 2,580 MB,
  reserve 258 MB, projected 512 MB, available 2,322 MB. Every stage fitted the
  preferred reserve.
- Live sampler peaks: stage 1 **6,075 / 1,985 MB** used/free; stage 2
  **5,709 / 2,351 MB**; stage 3 **6,092 / 1,968 MB**; stage 4 **5,717 / 2,343 MB**.
  Lowest free stayed well above the 819 MB reserve.
- The UI showed `Pipeline complete. Weights retained unless reclamation reported
  otherwise.` Four PNGs are in `workflow-20261124_214632`; stage 4 SHA-256 is
  `f6c95bd05e8838a5f6befdf05b669172fb3e20fdf77dd7e25d723b8dc87c82e0`, identical to
  the browser run.
- Evidence: `desktop-soft-live.log`, `desktop-soft-vram.tsv`,
  `desktop-soft-complete.png`, `desktop-soft-stages3.png`,
  `desktop-soft-stages2.png`.

The packaged host and browser host now both complete the four-stage workflow under the
revised admission policy.

### 9.12 Fresh failure recovery and consecutive-run reuse (2026-11-24)

The browser harness was started fresh and driven through the normal pipeline controls.
The upscaler path was pointed at `workflow-assets/invalid-upscaler.gguf` (a real file
with non-GGUF contents), then restored to `2x-ESRGAN.gguf` for two complete back-to-back
runs. Telemetry sampled VRAM, free system RAM and the `sd-server` working set every
200 ms. All three runs used the same isolated root.

#### Failure row

- Stage 1 completed in **13,330 ms**; startup admission was the untiled q8 request:
  free **7,264 MB**, reserve **726 MB**, projected **6,533 MB**, tiling false (it fitted
the preferred reserve by 5 MB).
- Stage 2 DGL failed in **679 ms** with
  `generation_failed: line 10, column 16: upscaler loading failed`; the UI reported
  `Pipeline failed: stage 2 upscale`. The failed stage still reported admission
  (`projected=512MB`, `free=2683MB`, `reserve=268MB`) and placement working set
  (171 MB after load, 408 MB after generation).
- The existing owner recorded:
  `[ImageReclaim] decision=retain ... measured=408MB growth=237MB free=2683MB
  pid=21608->21608`. There was no component reload, no second engine start and no
  competing reclamation path. The INFO log cannot show an unload; see §9.14.
- Raw telemetry during the failed base stage peaked at **7,429 MB used / 631 MB
  nvidia-smi-free**; the product's `total - used` free was **763 MB** against its
  **726 MB** reserve, a 37 MB margin. Minimum free system RAM was **48,719 MB**.
- Before the first run the host was **956 MB VRAM used / 7,104 MB free** and
  **50,249 MB free system RAM**; immediately after the second run, with the server
  deliberately retained, it was **4,718 MB VRAM used / 3,342 MB free** and
  **48,162 MB free system RAM**. The server was then stopped and the card returned to
  **7,112 MB free**.

#### Two consecutive successful runs

| Run | Base | Upscale | ControlNet | Final upscale | Peak used/free (nvidia-smi) | min free RAM | max server WS |
|---|---:|---:|---:|---:|---|---:|---:|
| 1 | 4,707 ms | 23,540 ms | 115,669 ms | 22,417 ms | 7,461/599; 4,774/3,286; 6,562/1,498; 4,746/3,314 MB | 47,906 MB | 2,220 MB |
| 2 | 4,780 ms | 22,987 ms | 117,177 ms | 22,836 ms | 6,567/1,493; 4,777/3,283; 6,605/1,455; 4,789/3,271 MB | 47,908 MB | 2,144 MB |

- Run 1 ended at 15:57:38; run 2 started at 15:58:01, so the workflows did not overlap.
  Each workflow's stage PNG timestamps are strictly sequential (run 2: 15:58:10,
  15:58:33, 16:00:30, 16:00:53).
- Run 1 stage 1 reused the resident base: working set after generation **409 MB**,
  not a fresh model load; stage 2 logged one `Upscaler architecture: esrgan`; stage 3
  logged one `ControlNet weights: 1415.82 MiB`. Run 2 logged neither a server start
  nor an upscaler/ControlNet load.
- Both runs retained the same PID **21,608**. Counts across failure + both successes:
  `Starting sd-server` = **1**, `Upscaler architecture` = **1**, `ControlNet weights` =
  **1**, component-reload lines = **0**, `[ImageReclaim]` = **3** (all `decision=retain`).
  Engine unload lines are VERBOSE and were not captured; see §9.14.
- Run 2 stage admissions: stage 2 free 3,546/reserve 354 MB; stage 3 free 3,548/
  reserve 354 MB, projected 2,432 MB, activation 1,920 MB; stage 4 free 3,507/
  reserve 350 MB. The final stage's working set was 1,942 MB.
- Artifacts: `workflow-20261124_215450` and `workflow-20261124_215805`; both stage 4
  PNGs hash to
  `5165011ea3bcedbff6d509f48c9261ea9fbe45e45d3006dfef0024440b70e535`.
- Evidence: `harness-failure-reuse.log`, `failure-reuse-telemetry.tsv`,
  `failure-run-start.txt`, `reuse-run1-start.txt`, `reuse-run2-start.txt`.

The invalid fixture was removed after the run. No product code change was needed: the
failure was reported, the one reclamation owner retained the server, and both recovery
runs reused it.

### 9.13 Fresh Batch > 1 proof under the revised policy (2026-11-24)

A fresh browser-harness batch used `BatchCount=2` through the normal **Run pipeline**
button. Inputs: base/refine 512×512, 4 steps, seed 9001 (item seeds 9001/9002),
`2x-ESRGAN.gguf`, `OpenPoseXL2.safetensors`, and the supplied 512×512 COCO-18 pose
`stage3-control.png`. The server was started from idle (1,003 MB used / 7,057 MB free
at 16:35:17); startup admission projected **5,093 MB** against free **7,138 MB** with
reserve **713 MB**, and planned VAE tiling and ControlNet CPU placement.

#### Per-item stage evidence

| Item | Stage | Duration | Admission (free / reserve / projected) | Peak used/free MB | Min free RAM | Max server WS |
|---|---|---:|---|---|---:|---:|
| 1 | base | 15,101 ms | startup 5,093 MB; free 7,138; reserve 713; tiling + ControlNet CPU | 6,038 / 2,022 | 48,563 MB | 1,347 MB |
| 1 | ESRGAN | 23,952 ms | free 2,634 / reserve 263 / projected 512 | 5,698 / 2,362 | 49,505 MB | 444 MB |
| 1 | pose/ControlNet | 83,697 ms | free 2,613 / reserve 261 / projected 992 | 6,067 / 1,993 | 47,663 MB | 2,281 MB |
| 1 | second ESRGAN | 23,091 ms | free 2,605 / reserve 260 / projected 512 | 5,691 / 2,369 | 47,873 MB | 1,957 MB |
| 2 | base | 6,128 ms | free 2,600 / reserve 260 / projected 992 | 6,067 / 1,993 | 48,045 MB | 1,947 MB |
| 2 | ESRGAN | 23,104 ms | free 2,605 / reserve 260 / projected 512 | 5,696 / 2,364 | 47,732 MB | 1,962 MB |
| 2 | pose/ControlNet | 77,662 ms | free 2,600 / reserve 260 / projected 992 | 6,091 / 1,969 | 47,778 MB | 2,144 MB |
| 2 | second ESRGAN | 23,154 ms | free 2,581 / reserve 258 / projected 512 | 5,715 / 2,345 | 47,912 MB | 1,968 MB |

- Item 1 base is covered by the startup admission; every later stage carries its own
  admission row. Item 2 base, which reused the resident engine, was admitted as a
  **992 MB** delta (free 2,600 / reserve 260). This was a real code gap found by this
  run: previously a reused base stage had placement but no fresh admission.
- Stage PNG times are strictly ordered. Item 1: 16:35:32.877, 16:35:57.174,
  16:37:21.067, 16:37:44.296. Item 2: 16:37:50.879, 16:38:14.111, 16:39:31.954,
  16:39:55.267. Item 2's base started at 16:37:44.751, **455 ms** after item 1's final
  artifact; there is no overlap between any stage, predecessor, or peer item.
- The batch took **277.5 s**: item 1 146.5 s, boundary 0.455 s, item 2 130.5 s.
  Minimum free VRAM was **1,969 MB**, 1,150 MB above the 819 MB reserve; minimum free
  system RAM was **47,663 MB**. After the second item, with the server retained, the
  card read **5,658 MB used / 2,402 MB free** and **47,698 MB free RAM**.
- Same server PID **8,760** throughout. Counts: `Starting sd-server` **1**,
  `Upscaler architecture: esrgan` **1**, `ControlNet weights` **1**, four
  `sampling started` (base + pose for each item). The INFO log shows no component
  reload. It cannot show "zero unloads": `DGL unloaded ...` is VERBOSE and this run
  captured no `[V]` lines. A separate DEBUG capture of the same batch shape (§9.14)
  recorded zero instrumented release/eviction/unload lines and zero parameter
  preparations in its second pass; that closes the unload gap for this batch shape
  without changing this run's INFO-only claim. See §9.14.
- The item-1 boundary call to the existing reclamation owner logged
  `decision=retain ... baseline=1933MB measured=1933MB growth=0MB free=2605MB
  pid=8760->8760`. The terminal call logged `decision=retain ... measured=1944MB
  growth=11MB free=2582MB pid=8760->8760`. No other reclamation path exists.

#### Artifacts

Directory `C:\\Users\\StoneWolfPC\\howl-image-hw-iso-20261124\\temp\\workflow-20261124_223517`:

| Item | Final path | Dimensions | SHA-256 |
|---|---|---|---|
| 1 | `run-1-stage-4.png` | 1024×1024 | `42b21711200016fb3c1be43ebfadefeb899727a17d7e3ff420fb39ff564d563b` |
| 2 | `run-2-stage-4.png` | 1024×1024 | `d7326f5ff23d590934491861ab0fb26265b31cfda88feef487d63952dbb1deaa` |

All eight stage PNGs are in that directory; intermediate hashes are
`b9cda305…30c8c` (item 1 base), `f598649f…1ffb` (item 1 upscale),
`cf674cb6…7f0f` (item 1 pose), `dafa8a07…640` (item 2 base),
`d5f4495f…c3c2` (item 2 upscale), and `39798f0c…94d4` (item 2 pose).

#### Code change

`app-image-workflow-run.go` now calls the **same** reclamation owner at the item
boundary (`runIndex < batch-1`) and records that report on the completed
`ImageWorkflowRun`. It also sends a base-stage delta admission when the engine is
already resident, so a reused base reports the same admission/lease path as the
other stages. No authority, policy constant, or competing reclamation path changed.

Evidence: `harness-batch3.log`, `batch-telemetry-3.tsv` (200 ms samples),
`batch3-harness-build.log`, `batch-telemetry-3.tsv`. The obsolete prior-policy batch
row was replaced by this evidence.

### 9.14 Pinned-source audit: engine placement and memory evidence (2026-11-24)

Scope: this is a source audit of the exact commit the installed binary reports,
`3b12781fac1a22eda27ad0f08d0884ef816f24d7` (binary line:
`stable-diffusion.cpp version unknown, commit 3b12781` in `residency-verbose.log`). No
engine, model, or generation process was started for this section. The audit uses that
commit's source tree; the `file:line` references below are from the commit, not from a
rebuilt working tree. The installed binary's `--- file:line` tags are a few lines ahead of the tarball
in some files (e.g. `run_plan.cpp:79` vs 77, `model_manager.cpp:695` vs 691); match the
format string, not the tag.

#### The current log level does discard real evidence

`sd-server` defaults to `SD_LOG_INFO` (`examples/server/runtime.h:33`), and
`log_print` drops every message whose level is numerically lower
(`examples/common/log.cpp:83-89`; enum order is DEBUG, VERBOSE, INFO, WARN, ERROR in
`include/stable-diffusion.h:276-283`). The current workflow startup command in
`harness-batch3.log` contains no `--log-level`/`-v`, and that log contains **0** `[V]`
and **0** `[D]` lines but **4** INFO `weights:` lines.

Therefore the current startup:

- **keeps** loaded parameter bytes for the diffusion model, conditioner, VAE and
  ControlNet, plus the upscaler architecture;
- **discards** every DGL `auto-fit plan` planned target, every `compute buffer size` /
  `free compute buffer` workspace line, every `DGL unloaded` line, and every
  `model manager prepared/staged/evicted/reused` line;
- **does not** discard the INFO `... weights` evidence. The earlier shorthand that all
  placement evidence is verbose-only was too broad.

A single bounded run with `--log-level debug` (most detailed; includes VERBOSE and
INFO) is the only live capture needed to observe the sites below in this binary.

#### Emit-site inventory

`I` = INFO, `V` = VERBOSE, `D` = DEBUG. All paths are in the pinned commit.

| ID | Level | Site | Exact format string | Attribution |
|---|---|---|---|---|
| P1 | I | `src/dgl/models/model_runtime.cpp:483` | `"%s weights: %.2f MiB"` | `model_component_name(component)`; loaded/registered bytes, no device |
| P2 | I | `src/dgl/models/model_runtime.cpp:245` | `"%s row split: %zu tensors (%.1f MB) split across %zu devices (main %s)"` | Component, total MB, device count, main device; no per-device bytes |
| P3 | I | `src/dgl/models/model_runtime.cpp:317` | `"%s graph-cut layer split: deferring %zu tensors across %zu runtime backends until first graph"` | Component, tensor count, backend count; no bytes or device names |
| P4 | V | `src/dgl/run_plan.cpp:73` | `"auto-fit plan: %zu models, %zu phases"` | Plan header |
| P5 | V | `src/dgl/run_plan.cpp:75` | `"  %s: available weight/compute budget %.2f MiB"` | Backend budget, not a component |
| P6 | V | `src/dgl/run_plan.cpp:77` | `"  %s: %.2f MiB params -> %s"` | DGL expression name, params MiB, selected `backend`/`disk`/`(none)`; **planned** target |
| P7 | V | `src/dgl/run_plan.cpp:91` | `"  phase %s: %s"` | Phase and the DGL resources active in it |
| P8 | D | `src/dgl/run_plan.cpp:224` | `"preserved active or split placement for %s"` | DGL expression name when `place_params` did not move it |
| P9 | V | `src/model_manager.cpp:691` | `"model manager prepared params backend buffers (%6.2f MB, %zu tensors, %zu blocks, %s) on %s"` | Actual buffer bytes/tensors/blocks, `RAM`/`VRAM`, buffer/backend name; **no component name** |
| P10 | V | `src/model_manager.cpp:881` | `"model manager staged compute params (%6.2f MB, %zu tensors, %zu blocks) to %s, taking %.2fs"` | Actual staged copy to compute backend; no component name |
| P11 | V | `src/model_manager_files.cpp:102` | `"model manager reused %zu resident weight tensors"` | Count only, no component or bytes |
| P12 | V | `src/model_manager.cpp:1293` | `"model manager evicted %.2f MiB of staged weights on %s for RAM"` | Staged buffer bytes and backend; no component |
| P13 | V | `src/model_manager.cpp:1303` | `"model manager evicted %.2f MiB of weight storage from RAM"` | Bytes only; no component or backend |
| P14 | V | `src/model_manager.cpp:1465` | `"model manager releasing params backend buffer (%6.2f MB, %zu tensors, %s)"` | Failure path only; no component |
| P15 | D | `src/model_manager.cpp:1795` | `"model manager released compute params (%6.2f MB, %zu tensors, %zu blocks) from %s"` | Backend and bytes; no component |
| P16 | D | `src/model_manager.cpp:1804` | `"free params buffer: %.2f MB(%s), %zu tensors on %s"` | Bytes, `RAM`/`VRAM`, tensor count, buffer; no component |
| P17 | D | `src/model_manager.cpp:2614` | `"model manager evicted %.2f MiB of staged weights on %s"` | Bytes and backend; no component |
| P18 | D | `src/model_manager.cpp:2623` | `"model manager evicted %.2f MiB of weight storage on %s"` | Bytes and backend; no component |
| P19 | V | `src/dgl/context.cpp:49` | `"DGL unloaded idle %s"` | DGL value type name; no bytes |
| P20 | V | `src/dgl/context.cpp:63` | `"DGL unloaded %zu cached model components"` | Count only |
| P21 | V | `src/core/execution_engine.cpp:748` | `"%s compute buffer size: %.2f MB(%s) on %s (peak across %zu segment%s)"` | Runner description (`unet`, `vae`, `clip`, `control_net`, `esrgan`), actual peak MB, `RAM`/`VRAM`, backend |
| P22 | V | `src/core/compute_workspace.cpp:314` | `"free compute buffer: %.2f MB(%s) on %s"` | Actual workspace release bytes, `RAM`/`VRAM`, backend; no runner name |
| P23 | I | `src/dgl/models/image_upscale.cpp:95` | `"Upscaler architecture: %s, weight type: %s"` | Architecture/type only; no bytes or device |
| P24 | W | `src/model_manager.cpp:2468` | `"model manager memory on %s: reported free %.2f MB / total %.2f MB, tracked weights %.2f MB / other runtime %.2f MB / current runtime %.2f MB"` | Backend aggregate, only when capacity logging requests details; no component |
| P25 | W | `src/core/execution_engine.cpp:975` | `"execution engine cannot make enough memory available on %s: need %.2f MB device / %.2f MB budget, available %s device / %s budget"` | Failed capacity request, backend aggregate |

Other memory-adjacent lines exist but do not add component attribution:
`src/core/ggml_graph_cut.cpp:261` prints the auto-detected `--max-vram` budget at INFO;
`src/core/ggml_extend_backend.cpp:675-678` lists devices at VERBOSE; and
`src/core/backend_fit.cpp:376` prints a per-component planned table (DiT, Conditioner,
VAE only, no ControlNet or upscaler) at INFO, but `sd-server` does not call
`derive_backend_specs`, so that table is not part of the server path. The server-side
`/sdcpp/v1/models/status` response (`examples/server/model_catalog.cpp:681-701`) has
state, recipe, models, revision and error only; it has no memory or placement field.
The idle staging release path is also quiet: `sd_dgl_ctx_release_staging`
(`src/dgl/dgl.cpp:144-147`) calls `ModelManager::release_idle_gpu_staging()` with no
success log of its own; only the DEBUG buffer lines P15/P16 can appear. The
`--weight-residency-ttl` unload path only produces P19/P20 if `Context::unload_models`
actually removes a model, and the server wrapper logs only failures
(`examples/server/async_jobs.cpp:740`).

#### Per-component result

"Covered" means an emit site exists and was observed in the existing verbose log;
"uncovered" means no emit site names that fact for that component. The `P` identifiers
resolve to the exact pinned-source `file:line` in the inventory above.

| Component | Parameter bytes | Load target backend/device | VRAM vs RAM split | Eviction / unload | Scope of proof |
|---|---|---|---|---|---|
| Base diffusion model (UNet/DiT) | **Covered at INFO**: `Diffusion model weights: 2910.44 MiB` (P1; observed) | **Covered as planned target at VERBOSE**: `load_unet: 2910.44 MiB params -> CUDA0` (P6; observed) | **Weights: uncovered.** Actual buffer line P9 is aggregate only. **Workspace: covered at VERBOSE**: `unet compute buffer size: 125.33 MB(VRAM) on CUDA0` and `unet compute buffer size: 26.30 MB(RAM) on CPU` (P21; observed) | **Uncovered by component.** Only P12-P20 buffer/type/count lines; none names the diffusion model | Bytes: INFO. Target: planned, VERBOSE. Actual weights split and component eviction: not exposed |
| ESRGAN / DGL upscaler | **INFO has architecture only**: `Upscaler architecture: esrgan, weight type: f16` (P23; observed). **Bytes at VERBOSE only**: `load_upscaler: 31.80 MiB params -> CUDA0` (P6; observed). No `ESRGAN weights` line exists because this model does not use `ModelRuntime::rebuild` (P1) | **Covered as planned target at VERBOSE**: `load_upscaler: 31.80 MiB params -> CUDA0` (P6; observed) | **Weights: uncovered.** P9/P10 are aggregate. **Workspace: covered at VERBOSE**: `esrgan compute buffer size: 104.00 MB(VRAM) on CUDA0` and `0.19 MB(RAM) on CPU` (P21; observed) | **Uncovered by component.** A DGL type/count line may appear (P19/P20); no bytes or explicit "ESRGAN" unload | Architecture: INFO. Bytes/target: VERBOSE only. Actual weights split and component eviction: not exposed |
| ControlNet | **Covered at INFO**: `ControlNet weights: 1415.82 MiB` (P1; observed) | **Covered as planned target at VERBOSE**: `controlnet::load: 1415.82 MiB params -> CPU` (P6; observed) | **Weights: uncovered.** P9/P10 aggregate. **Workspace: covered at VERBOSE**: `control_net compute buffer size: 122.42 MB(VRAM) on CUDA0` and `3.67 MB(RAM) on CPU` (P21; observed) | **Uncovered by component.** Generic buffer/type lines only | Bytes: INFO. Target: planned, VERBOSE. Actual weights split and component eviction: not exposed |
| VAE | **Covered at INFO**: `VAE weights: 159.68 MiB` (P1; observed) | **Covered as planned target at VERBOSE**: `load_vae: 159.68 MiB params -> CUDA0` (P6; observed) | **Weights: uncovered.** P9/P10 aggregate. **Workspace: covered at VERBOSE**: `vae compute buffer size: 1920.00 MB(VRAM) on CUDA0` and `0.06 MB(RAM) on CPU` (P21; observed; tiled runs also show 976.50 MB) | **Uncovered by component.** Generic buffer/type lines only | Bytes: INFO. Target: planned, VERBOSE. Actual weights split and component eviction: not exposed |
| Text encoders (SDXL conditioner group) | **Covered at INFO, aggregate**: `Conditioner model weights: 925.92 MiB` (P1; observed). No CLIP-L/CLIP-G split | **Covered as planned target at VERBOSE**: `load_text_encoder: 925.92 MiB params -> CUDA0` (P6; observed) | **Weights: uncovered.** P9/P10 aggregate; no per-encoder split. **Workspace: covered at VERBOSE**: `clip compute buffer size: 1.42 MB(VRAM) on CUDA0` and `0.02 MB(RAM) on CPU` (P21; observed) | **Uncovered by component.** Generic buffer/type lines only | Bytes: INFO aggregate. Target: planned, VERBOSE. Actual weights split, per-encoder split and component eviction: not exposed |
| Compute workspace (all runners) | Not a weight component | Each P21 line names backend and actual peak; no separate "load target" | **Covered at VERBOSE per runner/backend** (P21); release is covered at VERBOSE (P22) | Workspace release is P22; persistent workspace unload is not separately named | Actual peak and release covered, but only while a runner executes; no live current-occupancy query |

Conclusion for the matrix:

- **Proven by this build now:** loaded bytes for base, conditioner, VAE and ControlNet
  (INFO); upscaler architecture (INFO); planned params target for base, upscaler,
  ControlNet, VAE and text encoder (VERBOSE); actual per-runner workspace peak and
  release, split by RAM/VRAM and backend (VERBOSE); buffer-level reuse/eviction/unload
  evidence (VERBOSE/DEBUG).
- **Not proven by this build:** actual per-component weights placement after allocation
  (as opposed to the plan), actual per-component weights VRAM/RAM split, per-component
  eviction bytes, per-encoder text split, and any live residency snapshot. Those need
  either a correlator with a patch or the upstream memory-stats API; they are not
  recoverable from these logs alone.
- **Current startup:** at INFO it **is** discarding planned placement, workspace and
  unload evidence that already exists in this build. It is **not** discarding the
  component loaded-byte lines.

#### Corrections to earlier wording

1. §9.8 previously called the verbose `load_unet`/`load_vae`/`controlnet::load` lines
   "component buffers" and used them as actual placement. They are **planned targets**.
   Actual allocation is the aggregate P9/P10 line; per-component attribution is not
   explicit. The table header there now says "Planned target (V)".
2. §9.3, §9.5, §9.11, §9.12 and §9.13 previously said "zero/no unloads" for
   INFO-only runs. Those logs cannot show unloads because P19/P20 are VERBOSE. The
   correct claim is "no component reload" (one architecture/weight line each), not
   "no unload".
3. "No per-component residency" was too broad. The engine does expose per-component
   loaded bytes and per-runner workspace. What it does not expose is a live component
   map, actual per-component weights split, or per-component eviction.
4. The prior turn's "verbose/debug is the only placement source" was also too broad:
   INFO already exposes loaded bytes. The precise statement is: bytes are INFO,
   planned target and workspace are VERBOSE, buffer eviction/release is DEBUG.

#### Live DEBUG capture (completed 2026-11-24)

The authorized observational run used the existing four-stage workflow and model
configuration unchanged, with `BatchCount=2` and the same 512×512/4-step seeds
9001/9002, `2x-ESRGAN.gguf`, `OpenPoseXL2.safetensors` (CPU) and the prepared COCO-18
pose. The only engine-argument change was `--log-level debug`; it was appended for this
build, then the source was reverted and the normal harness rebuilt. No admission,
placement, residency or reclamation policy changed.

The run completed both passes in **318.606 s** (`RunImageWorkflow ok 5m18.606s`), with
one `Starting sd-server` (PID 13,992) and all eight stage PNGs present. The UI rendered
the final pass as: base **6,211 ms**, ESRGAN **22,978 ms**, pose/ControlNet **101,237
ms**, final ESRGAN **22,529 ms**. Raw evidence is retained under
`build/verification/debug-capture-2026-11-24/`: `harness-debug.log` (182,012 bytes,
SHA-256 `55cc03b0…bf666b`), `debug-telemetry.tsv`, `ui-stages.txt`, and
`capture-summary.txt`.

**Observed results against P1-P25:**

- **All five runner workspaces appear:** `clip` 8 lines (1.42/2.65 MB VRAM, 0.02 MB
  RAM), `unet` 8 (249.41 MB + 1.35 MB plain; 125.33 MB + 26.30 MB with ControlNet),
  `vae` 8 (480.00 MB + 0.02 MB), `control_net` 2 (121.17 MB CPU RAM), `esrgan` 8
  (104.00 MB + 0.19 MB). The component set is complete, but exact workspace sizes are
  configuration-dependent: this run's VAE workspace was **480.00 MB**, not the
  1,920/976.5 MB seen in earlier verbose logs, and ControlNet was CPU-only (**121.17 MB
  RAM**), matching `--backend controlnet=cpu` rather than the earlier GPU-placement
  table. P9 actual parameter preparation appears seven times, all in pass 1: 151.93,
  740.58, 26.25, 2,910.65, 94.47, 31.80 and 1,415.82 MB; only the ControlNet block is
  tagged `RAM`/CPU. P9 has no component name, so the first six are not individually
  attributable.
- **ESRGAN has identifiable execution/workspace lifecycle:** GGUF parse, one
  `Upscaler architecture: esrgan`, one planned `load_upscaler: 31.80 MiB -> CUDA0`, one
  actual 31.80 MB/700-tensor preparation, 100 `esrgan executing segment 1/1` lines
  (4 upscales × 25 tiles), 104.00 MB VRAM + 0.19 MB RAM workspace, matching releases,
  and four `input_image_tensor upscaled` timings (23.18/22.52/22.52/21.94 s). There is
  no explicit `ESRGAN loaded` or `ESRGAN unloaded` line and no weight release for it.
- **Release attribution is partial:** workspace *peak* lines name the runner, but P22
  `free compute buffer` releases give only bytes, RAM/VRAM and backend; attribution is
  by adjacency and matching size, not by a runner field. P15-P18 weight release/eviction
  lines are absent (`free params buffer` 0, `released compute params` 0, `evicted` 0),
  as are P19/P20 DGL unload lines (`unloaded` 0). Howl's only reclamation evidence is
  two aggregate `ImageReclaim` `retain` decisions.
- **No unload or reload was observed across both stages and the item boundary:** same
  PID 13,992 throughout; one server start, one upscaler architecture line, one ControlNet
  weights line; seven P9 preparations in pass 1 and **zero in pass 2**; zero instrumented
  release/eviction/unload lines. Both `ImageReclaim` decisions were `retain`, with
  `pid=13992->13992` and growth 0/8 MB inside the 1,572 MB band. This is an
  absence-of-instrumented-events plus no-reload proof, not a live per-component
  residency map.
- **DEBUG overhead was material in time but not memory:** baseline INFO **277.926 s**
  versus DEBUG **318.606 s** (+40.68 s, +14.6%). Almost all of the delta is the CPU
  ControlNet stage: **101.44 s** and **101.40 s** in the two passes versus 83.697/77.662
  s in the INFO baseline. Base, ESRGAN and final-upscale timings stayed within ~1 s of
  baseline. Peak VRAM was 6,147 MB versus 6,091 MB (+56 MB); sd-server working-set peak
  was 2,208 MB versus 2,281 MB (−73 MB); free-RAM minimum was 46,513 MB versus 47,639 MB
  (system-level). The 182 KB log is not a disk-pressure explanation. Because the capture
  also ran a 250 ms telemetry sampler and ControlNet is CPU-bound, this single
  observational run cannot separate log I/O from CPU contention; DEBUG should remain a
  diagnostic mode, not production telemetry, until a matched control isolates the cause.

This closes the observable gap P1-P25 on the exact binary, but not the per-component
weights VRAM/RAM split or a live residency snapshot; those still require a patch or a
newer engine API. No reclamation or resource-policy change was made.
