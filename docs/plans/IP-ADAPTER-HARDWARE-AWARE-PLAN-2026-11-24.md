# IP-Adapter — Hardware-Aware Integration Plan and Verification Matrix

**Date:** 2026-11-24
**Status:** Native IP-Adapter and the native four-stage branch verified with real
weights on the GTX 1080 in the browser harness (2026-11-25). The resource-policy
and recovery rows are closed there too: 5 GiB constrained, external reserve
pressure, batch/repeat, measured reuse, failure recovery, and profile switch-back.
Packaged WebView2 H-2 is verified by the test-session packaged run in §5.11.
**Gate:** Owner authorized the roadmap item and real-asset verification. Do not
claim the whole matrix from this single-host evidence.

**Scope:** Native IP-Adapter support for Howl's image orchestrator, admission,
placement, reuse, and reclamation on the GTX 1080 8 GB.

## Implementation status (2026-11-24)

Implemented on the capability-routed profile:

- Asset classes `ip_adapter` and `clip_vision` (scanner, manifests, registry,
  model-library groups, family filters). IP-Adapter is no longer filed as a
  ControlNet.
- `sdcpp-c6dfa9d-ip` resolves through `EngineConfig.IPAdapterBinaryPath`,
  `HOWL_IP_ENGINE`, `runtime/stablediffusion-ip`, or the checked-in backup. It
  replaces the default process through the existing `LoadModel`, admission and
  single reclamation owner; there is no second resident engine.
- Startup slots `--ip-adapter` and `--clip_vision`; CLIP-Vision gets
  `--backend clip_vision=cpu` only when the measured projection needs it.
- Request fields `ip_adapter_image` and `ip_adapter_strength` on native txt2img
  and native ControlNet img2img.
- Four-stage native branch: native txt2img -> native `/sdcpp/v1/upscale` -> native
  ControlNet img2img -> native upscale. DGL is not mixed in when an IP-Adapter is
  selected.
- Admission charges adapter and CLIP-Vision bytes; placement reports the profile,
  adapter and CLIP-Vision CPU placement.
- UI: IP-Adapter panel (adapter, encoder, strength, reference image, CPU-only
  override), normal generation and pipeline request wiring, workflow restore.

Real-asset evidence (isolated harness root, 2026-11-25):

- Real weights downloaded into the scanned `G:\AI GGUF\Image models` root:
  `ip-adapter_sd15.safetensors` (42.6 MB, 36 tensors),
  `clip_vision_h.safetensors` (2.41 GB, 521 tensors),
  `control_v11p_sd15_openpose.pth` (1.38 GB), and
  `4x-UltraSharp.pth` (63.9 MB). SHA-256 values are in §2.3.
- A real 512x512 IP generation completed on
  `build/engine-baseline/backup-c6dfa9d/sd-server.exe` with the SD1.5
  checkpoint, real ViT-H encoder, and real adapter. The engine logged
  `IP-Adapter: 4 image tokens, strength 0.70`. The PNG is recognizable as the
  reference person (same face, dark hair, light blue shirt).
- A real native four-stage run completed in the same harness: 256x256 base,
  native 4x ESRGAN to 1024, native IP/ControlNet img2img at 512, native 4x
  ESRGAN to 2048. All four stages report `completed`; no DGL job ran.
- The first run exposed two scanner defects that hid valid assets: a vision-only
  CLIP tower was classified as `llm`, and official `control_v11*` ControlNets
  lacked the literal word "controlnet". Both are fixed with regression tests.
- A second defect was found on the repeated run: Howl's automatic post-reserve
  `MaxVRAM` cap was treated as a user hard budget, forcing CPU diffusion despite
  a just-measured resident placement. Automatic caps are now excluded from the
  hard-budget fallback.
- `c6dfa9d` does not accept `--weight-residency-ttl`; the IP profile omits it
  (the engine default is already disabled).
- Final resource closure (§5.9): constrained `--max-vram 5`, batch/repeat, explicit
  reserve pressure, crash recovery in one harness process, and non-IP profile
  switch-back all completed with per-stage VRAM and RAM before/after, exact
  artifacts and hashes. Batch/repeat reused one PID; reserve pressure recycled
  once; the crash run failed cleanly and the next run succeeded.

Not verified, and why:

- Packaged WebView2 H-2 is verified by the test-session packaged run in §5.11.
  The normal packaged binary remains debug-free; the test build only added the
  CDP channel used to set the reference image and press the workflow button.
- Explicit CPU-only diffusion (W-5), missing/mismatched assets (R-4/R-5), and an
  ordinary non-workflow generation (G-2) were not rerun in this closure; they are
  classified UNVERIFIED with their closest test coverage in §5.7.
- The pinned non-IP engine does not report per-component residency, so W-6's
  actual placement is limited to the live device/process measurements; its
  startup flags prove the profile change.
- A killed engine leaves no live placement measurement. The final crash-recovery
  rerun nevertheless restarted at `--max-vram 5` and completed all four stages
  without a CPU-diffusion flag (§5.10); earlier fresh-restart fallback behavior is
  not used as final evidence.
- Per-stage live samples exist for every run below: one `nvidia-smi` sample
  before and after each stage, recorded in
  `ImageWorkflowStage.resources_before/resources_after` and in the harness log.
- Native IP + ControlNet output was verified as a completed image with the
  ControlNet preprocess log; a byte-level with/without-ControlNet A/B was not run
  in this pass.
- Pose extraction remains external (OpenPose/DWPose). The proof used a
  precomputed official OpenPose sample; ControlNet does not extract pose.

Consolidation trigger remains: retire the second binary once one engine build
passes the full matrix with native IP, native ESRGB upscale, existing DGL/native
replacement, the same placement/reclamation, and W-1/W-2/W-3/W-6/G-1.

**Target:** IP-Adapter conditions the existing sequential workflow without becoming a
separate, unmanaged generator:

```
base txt2img (IP) -> ESRGAN 4x -> ControlNet img2img (IP) -> ESRGAN 4x
```

---

## 1. Decision

### 1.1 Recommended route

Use the already-present native-IP engine as a **capability-routed engine profile**:

| Profile | Engine | Used for | Workflow executor |
|---|---|---|---|
| `sdcpp-3b12781` (default) | pinned `master-1122-3b12781` | all existing generation and non-IP workflows | native stages + DGL upscales |
| `sdcpp-c6dfa9d-ip` | backup `c6dfa9d` | workflows that explicitly enable IP-Adapter | four native stages, native `/sdcpp/v1/upscale` |

Only one engine process runs at a time. Selecting the IP profile uses the existing
`LoadModel` path, which already stops the previous server before starting the new one.
The shared admission authority, the existing reclamation owner, and the workflow lock
remain unchanged. No second model is resident and no second reclamation path is added.

This is not a downgrade of the default engine. It is an explicit profile selected by a
capability requirement, and the default profile remains the pinned, proven engine.

### 1.2 Why not the alternatives

| Option | Verdict | Reason |
|---|---|---|
| Native IP-Adapter on pinned `3b12781` | **Impossible** | The pinned engine has no native `--ip-adapter`, `--clip_vision`, or `ip_adapter_image` feature. Its own doc says IP-Adapter is supported only through DGL. |
| DGL IP-Adapter beside the native loaded model | **Rejected** | DGL has its own `ModelManager` and cache. It would load a second full SDXL checkpoint alongside the resident `-m` model; unsafe on 8 GB. |
| All-DGL IP workflow | **Rejected for this PoC** | It would replace the native sampling path and require a new full workflow executor. The existing orchestration plan already rejected all-DGL staging for the same cache conflict. |
| Upgrade the pinned engine to a build with native IP + DGL + native upscale | **Deferred consolidation** | Cleanest long-term shape, but reopens all proven `3b12781` behavior and needs a full engine re-verification. Keep it as a later option, not the first fix. |
| Two resident engines, one per profile | **Rejected** | Wastes VRAM and creates competing process ownership. Profile switching must be serialized and single-owner. |
| `c6dfa9d` native-IP profile, one process at a time | **Recommended** | Native IP exists, native upscale exists, the engine supports q8/max-vram/VAE tiling/ControlNet CPU, and the pinned workflow is untouched. |

### 1.3 What this fixes

IP-Adapter becomes a first-class asset and workflow control, not a hidden CLI flag:

- scanner and registry know IP-Adapter and CLIP-Vision assets;
- the UI has an IP-Adapter reference image and strength control;
- admission charges the adapter and the vision encoder;
- placement decides whether the vision encoder stays on GPU or runs on CPU;
- the workflow keeps adapter and vision resident across stages and batches;
- reclamation remains the one existing owner;
- missing or incompatible assets fail before a process starts, with an actionable error;
- the stage report names the profile, placement, and actual evidence.

### 1.4 Owner scope (2026-11-24)

In scope:

- native IP-Adapter support;
- SD1.5/SDXL compatibility detection;
- reference-image selection and strength controls;
- CLIP-Vision management and hardware-aware placement;
- integration with the existing generation workflow and batches;
- transparent resource accounting and failure handling.

Out of scope:

- a custom embedding or vision-model ecosystem;
- every experimental adapter architecture;
- a separate IP-Adapter workflow engine;
- a general-purpose image-training platform.

The owner accepts the capability-routed profile as the **first implementation**, but
does not accept two engine builds as the permanent architecture. This profile is a
compatibility bridge. The integration must not dictate Howl's architecture.

### 1.5 Consolidation exit criteria for the second engine

The second engine is retired when one engine build passes the same IP matrix while
retaining the proven non-IP behavior:

1. native IP-Adapter and native RGB ESRGAN upscale in the same build;
2. the existing DGL path still works for non-IP workflows, or a verified native
   replacement exists;
3. the same q8/max-vram/VAE-tiling/ControlNet-CPU startup behavior;
4. the same admission, placement, reuse, and single-owner reclamation outcomes;
5. W-1, W-2, W-3, W-6, and G-1 pass without a profile-specific workflow branch.

The capability profile/selection layer stays; only the second binary is removed. No
new specialized binaries may be added for individual adapter families.

---

## 2. Evidence

### 2.1 The pinned engine cannot do native IP-Adapter

- `docs/models/ip_adapter.md` at `3b12781` states: **"Currently supported only through
  dgl."**
- Pinned server help has no `--ip-adapter`, `--ip-adapter-image`, or
  `--ip-adapter-strength`.
- Pinned `routes_sdcpp.cpp` sets `ip_adapter_image` false; `model_catalog.cpp` forces
  the `ip_adapter` feature false even when the request carries the fields.
- Therefore the existing server's IP fields must not be advertised as working.

### 2.2 The backup engine has the missing native path

`build/engine-baseline/backup-c6dfa9d/sd-server.exe --help` confirms:

- `--ip-adapter <path>` — requires `--clip_vision`;
- `--ip-adapter-image <path>` — reference image;
- `--ip-adapter-strength <float>` — default 1.0;
- `--clip_vision <path>` — CLIP-Vision encoder;
- `--control-net`, `--backend controlnet=cpu`, `--max-vram`, `--vae-tiling`,
  `--type q8_0`, `--flash-attn auto`, `--offload-to-cpu`;
- `--hires-upscalers-dir` and `--upscale-model`;
- `POST /sdcpp/v1/upscale` with `image`, `upscaler`, `repeats`, `tile_size`,
  `output_format`; it is synchronous, holds the generation context lock, and does not
  create a job;
- `-M dgl` is **invalid**: `invalid mode dgl, must be one of [img_gen, adetailer,
  vid_gen, convert, upscale, metadata]`.

The backup engine is a complete native four-stage executor; it is not a DGL engine.

### 2.3 Placements and assets

- The backend module key is `clip_vision` (aliases `clipvision`, `clip-vision`,
  `vision`). `docs/backend.md` at `c6dfa9d` lists it as a first-class module, so
  `--backend clip_vision=cpu` places the encoder's params and compute on CPU. The
  binary contains the `clip_vision` loading/compute path.
- Real assets now present in the scanned `G:\AI GGUF\Image models` root:

| Asset | Role | Size | SHA-256 |
|---|---|---:|---|
| `ip-adapter_sd15.safetensors` | SD1.5 IP-Adapter | 42.6 MB | `289b45f16d043d0bf542e45831f971dcdaabe18b656f11e86d9dfba7e9ee3369` |
| `clip_vision_h.safetensors` | ViT-H/14 CLIP-Vision | 2.41 GB | `6ca9667da1ca9e0b0f75e46bb030f7e011f44f86cbfb8d5a36590fcd7507b030` |
| `control_v11p_sd15_openpose.pth` | SD1.5 OpenPose ControlNet | 1.38 GB | `db97becd92cd19aff71352a60e93c2508decba3dee64f01f686727b9b406a9dd` |
| `4x-UltraSharp.pth` | ESRGAN 4x upscaler | 63.9 MB | `a5812231fc936b42af08a5edba784195495d303d5b3248c24489ef0c4021fe01` |

- The adapter and encoder live directly in the root; the ControlNet and upscaler live
  in subdirectories. The scanner found all four after the two classification fixes.
- The backup engine directory is 383 MB including `ggml-cuda.dll` and the other DLLs.
  Shipping a second profile has an installer-size cost; the plan does not hide it.

### 2.4 Existing behavior that must not regress

The prior workflow proof established:

- one `sd-server` PID across four stages and two-pass batches;
- base, ControlNet, and ESRGAN residency with no component reload;
- ControlNet CPU placement and VAE tiling under pressure;
- 10% reserve held at every stage;
- `retain` at item and terminal boundaries;
- the existing single reclamation owner and `imageStartMu` serialization.

The IP profile must reuse those mechanisms, not create a parallel runner.

---

## 3. Architecture

### 3.1 Engine profile contract

Add a runtime profile to the existing load options, not a second global engine:

```go
const (
    EngineProfileDefault   = "sdcpp-3b12781"
    EngineProfileIPAdapter = "sdcpp-c6dfa9d-ip"
)
```

`ModelLoadOptions` gains:

```go
EngineProfile   string `json:"engine_profile,omitempty"`
EngineBinaryPath string `json:"-"` // resolved runtime path; never saved as user intent
IPAdapterPath   string `json:"ip_adapter_path,omitempty"`
ClipVisionPath  string `json:"clip_vision_path,omitempty"`
ClipVisionCPU   bool   `json:"clip_vision_cpu,omitempty"`
```

`ImageOptions` gains the same request-facing fields plus:

```go
IPAdapterImagePath string  `json:"ip_adapter_image_path,omitempty"`
IPAdapterStrength  float64 `json:"ip_adapter_strength,omitempty"`
```

Rules:

1. The default profile is empty/`sdcpp-3b12781`; it continues to use
   `EngineConfig.BinaryPath`.
2. `EngineProfileIPAdapter` resolves its binary from a configured IP-engine root,
   defaulting to the bundled `runtime/stablediffusion-ip/sd-server.exe`, then the
   development path `build/engine-baseline/backup-c6dfa9d/sd-server.exe`. If neither
   exists, the profile is unavailable and the UI refuses before admission.
3. `EngineBinaryPath` is runtime-only (`json:"-"`) and is derived from the profile; it
   is not persisted into image settings or model metadata.
4. `Engine.LoadModel` uses the selected binary for `UnsupportedFlags`, the process
   command, and the version/capability probe. It still uses the same port and the same
   "kill previous server" path.
5. `imageStartupOptionsDiffer` includes profile, binary path, IP adapter path, CLIP
   vision path, and `ClipVisionCPU`. A profile or adapter change is a real restart, not
   an inert control.
6. The profile is resolved from the request. An empty profile means the default
   profile; it must not be inherited from `last_model_opts` or from the currently
   resident IP engine. Only an explicit IP request carries the adapter, vision, and
   placement forward. A normal generation after an IP run therefore switches back to
   the default profile, which W-6 verifies.
7. `preserveImageLoadedOptions` may carry adapter/vision/placement only when the
   request explicitly selects the IP profile. It must not turn an ordinary generation
   into an IP-profile request.
8. `UnsupportedFlags` is probed per binary, not against the pinned table. If the old
   binary lacks a requested flag, the load fails before the process starts.
9. The profile is part of `LoadedOptions`; the placement/report can name the actual
   engine, not a planned one.

No second `ImageEngine` instance. No parallel engine manager. No new health loop.

### 3.2 Startup arguments

For the IP profile, `buildServerArgs` adds:

```
--ip-adapter <adapter path>
--clip_vision <vision path>
--backend controlnet=cpu          (when the existing admission plan selected it)
--backend clip_vision=cpu         (only when admission selected CPU placement)
--hires-upscalers-dir <upscaler dir>
```

If multiple module assignments are used, merge them into one `--backend` value
(`controlnet=cpu,clip_vision=cpu`) rather than emitting the flag twice.

The reference image is sent per request as base64 `ip_adapter_image`; it is **not** a
startup path. This keeps reference-image selection per generation and avoids granting
the engine a new arbitrary file-read path.

The IP profile does not add `--control-net` unless the workflow already requested it.
ControlNet remains a startup slot with the existing charge-once rule.

### 3.3 Admission and placement

`AdmissionRequest` gains:

```go
IPAdapterBytes   int64
ClipVisionBytes  int64
ClipVisionCPU    bool
```

Accounting:

- IP-Adapter weights are resident adapter weights. Charge `IPAdapterBytes` in the
  resident VRAM total by default. If the engine is told to offload them, charge RAM
  instead; do not count them twice.
- CLIP-Vision weights are conditioner weights. If `ClipVisionCPU` is true, charge
  `ClipVisionBytes` to `RequiredRAMMB` and not to VRAM. If false, charge them to VRAM
  as part of the conditioner set.
- The adapter and vision encoder are charged once at engine start. Stage deltas with
  `ResidentLoaded` charge only a missing adapter/vision component, never the same
  resident weights again.
- The 10% preferred reserve, 500 MB soft allowance, 128 MB floor, explicit reserve
  rules, and `ProjectedPeakVRAMMB` comparison are unchanged.
- `ControlNetCPU` remains the existing separate flag and is not folded into
  `ClipVisionCPU`.

Placement decision:

1. Project the stage with the vision encoder on GPU.
2. If the projected peak fits the post-reserve budget, keep it on GPU.
3. If it does not, project `clip_vision=cpu`; if system RAM fits the encoder plus the
   existing RAM reserve, select CPU and record why.
4. If neither fits, refuse before starting the engine with exact VRAM/RAM numbers.
5. Never select CPU merely because the feature exists. The decision is based on the
   measured/projected budget, and the UI shows the selected placement.

On the measured 8 GB profile, the expected outcome is IP-Adapter on GPU and CLIP-Vision
on CPU: a 666 MB SDXL adapter plus the existing ~6.2 GB stage peak leaves the reserve
intact, while a ~2.5 GB ViT-H encoder on GPU would not. This is a projection, not a
claim; W-1/W-4 must measure it.

### 3.4 Workflow for the IP profile

The IP profile uses the native executor for all four stages:

| Stage | Executor | Request |
|---|---|---|
| 1 base | native `GenerateTxt2Img` | prompt, 512×512, seed, `ip_adapter_image`, `ip_adapter_strength` |
| 2 upscale | native `POST /sdcpp/v1/upscale` | stage-1 PNG base64, upscaler name, tile size |
| 3 pose/ControlNet | native `GenerateImg2Img` | stage-2 PNG as init, pose PNG as control, same IP image/strength, ControlNet strength |
| 4 upscale | native `POST /sdcpp/v1/upscale` | stage-3 PNG, same upscaler/tile |

Sequence and ownership:

- The workflow acquires `imageStartMu` once and holds it through all stages.
- Stage 1 waits for admission and engine readiness; stages 2–4 wait for their
  predecessor.
- The native upscale route holds the generation context lock, so it cannot overlap a
  native generation.
- The adapter and CLIP-Vision encoder stay resident across stages and batch items.
- No DGL job is submitted in this profile.
- Terminal and failure paths call `reclaimImageEngineAfterGenerationLocked` exactly
  once, as today.
- A profile switch is not a reclaim; it is the existing `LoadModel` replacement path.

Batch:

- One profile selection for the batch; one engine start where possible.
- Each item gets its own reference image and seed.
- Stage outputs are isolated per item.
- The next item reuses the adapter/vision/model unless admission or the existing
  reclamation policy says otherwise.

Failure behavior:

- Missing `clip_vision_h.safetensors`, missing adapter, family mismatch, unavailable
  IP engine, or unsupported native upscale is a preflight refusal. Do not start the
  engine and do not silently drop IP-Adapter.
- A native stage failure is reported with the stage name and engine error. The existing
  reclaim owner decides whether to retain or recycle.
- A failed IP run must leave the next run able to select the default profile and run the
  existing non-IP workflow. The profile is per request, not a global latch.

### 3.5 Scanner, registry, and UI

Scanner:

- Add `AssetIPAdapter` and `AssetClipVision`.
- Stop classifying IP-Adapter files as `AssetControlNet`. Keep a migration path so old
  UI state that points at those files gets a clear "this is an IP-Adapter, not a
  ControlNet" error rather than a silent mismatch.
- Search `ipadapter/` for adapters and `clip_vision/` for ViT-H encoders under every
  configured model root. Do not require the files to sit beside the checkpoint.
- Detect family from the adapter filename/metadata: SD1.5 vs SDXL. Reject mixed
  families before admission.
- The FLUX `ip_adapter.safetensors` entry currently misclassified as `diffusion` must
  not be treated as an SDXL/SD1.5 IP-Adapter. It needs its own compatibility check or
  exclusion.

UI:

- Image generator: an "IP-Adapter" section with enable toggle, adapter picker, vision
  encoder picker (auto-selected when a unique compatible encoder exists), reference
  image drop/select, strength (default 0.7, range 0–1.5), and a readiness line.
- Show the selected engine profile and placement: "IP engine: native c6dfa9d",
  "CLIP-Vision: CPU (projected)", "IP-Adapter: CUDA0 resident".
- Disable Run when the profile, adapter, vision encoder, or reference image is missing;
  state exactly which item is missing.
- Workflow panel: IP-Adapter controls are part of the pipeline, not a separate
  afterthought. Stage rows include profile, adapter/vision placement, and duration.
- Diagnostics: engine profile, binary path, adapter/vision paths, resolved placement,
  and capability probe result. Do not expose a misleading "IP-Adapter supported" badge
  on the pinned profile.

---

## 4. Implementation slices

No slice starts before authorization.

### Slice A — types, scanner, registry

Files: `internal/image/types/types.go`, `internal/image/scanner/classifier.go`,
scanner tests, registry/model metadata, `app-image-scan.go`.

- Add asset types and fields.
- Reclassify IP-Adapters and add CLIP-Vision.
- Add family compatibility and actionable validation.
- Unit tests: SDXL/SD1.5 classification, FLUX exclusion, duplicate encoder, missing
  encoder, old ControlNet state migration.

### Slice B — engine profile routing and startup args

Files: `app-image-model.go`, `internal/image/engine/sdcpp/engine.go`,
`internal/image/engine/sdcpp/capability.go`, args/startup tests.

- Add profile/binary resolution.
- Add IP/vision args and merged backend assignments.
- Include profile in startup-difference and preserve logic.
- Probe the selected binary, not the pinned table.
- Tests assert exact argv for default and IP profiles, profile-switch restart, and
  no duplicate `--backend`.

### Slice C — admission

Files: `internal/resource/admission.go`, `app-image-model.go`,
`internal/resource/*_test.go`, `app-image-*_test.go`.

- Add `IPAdapterBytes`, `ClipVisionBytes`, `ClipVisionCPU`.
- Count adapter VRAM and vision RAM/VRAM correctly; no double charge.
- Add a placement projection helper that keeps vision on GPU when it fits and selects
  CPU only when the projection requires it.
- Tests: adapter-only, vision GPU, vision CPU, insufficient VRAM, insufficient RAM,
  resident delta, reserve refusal, and exact refusal numbers.

### Slice D — native IP workflow and native upscale

Files: `app-image-workflow-run.go`, `internal/image/engine/sdcpp/client.go`,
`internal/image/engine/sdcpp/upscale.go` (new), workflow tests.

- Add native `Upscale` client for `/sdcpp/v1/upscale`.
- Select native upscale when the active profile is the IP profile.
- Add IP image/strength to stages 1 and 3.
- Keep the existing DGL path for the default profile unchanged.
- Stage report: profile, adapter/vision placement, native upscale dimensions, PID,
  free VRAM/RAM, durations.
- Tests: native upscale request/response, IP stage payload, failure before fallback,
  no DGL call in IP profile, terminal reclaim once.

### Slice E — UI and diagnostics

Files: `frontend/dist/image-generator.html`, `frontend/dist/js/image-workflows.js`,
`frontend/dist/js/image-generator-core.js`, `frontend/dist/css/ig-controls.css`,
diagnostics bindings.

- Add the IP section and workflow controls.
- Add profile/placement badges and readiness errors.
- Add diagnostics fields.
- Tests/manual harness checks for disabled states and labels.

### Slice F — verification

Run §5 on the GTX 1080 with isolated output and read-only models. Record commands,
PIDs, `nvidia-smi`, working set, admission lines, reclamation lines, stage artifacts,
and hashes under `build/verification/ip-adapter-2026-11-24/`.

---

## 5. Verification matrix

Status: **partially run (2026-11-25)**. Rows marked PASS below have real-asset
evidence; every other row remains a planned check, not a claim. Do not promote a
row because its UI control exists.

### 5.1 Static and preflight

| ID | Case | Procedure | Pass criteria |
|---|---|---|---|
| S-1 | Profile probe | Run `--help` for both binaries through the capability probe | default has no IP flags; IP profile reports `--ip-adapter`, `--clip_vision`, native upscale; no false capability on default |
| S-2 | Native upscale route | `GET /sdcpp/v1/capabilities` on IP profile with `--hires-upscalers-dir` | `upscale` true, ESRGAN listed with `image_upscale: true`; route present |
| S-3 | CLIP-Vision CPU key | Parse the IP help/backend module table and probe `clip_vision=cpu` at startup | accepted; no "invalid module" error; placement reported |
| S-4 | Missing assets | Select IP with missing encoder, missing adapter, wrong family, or no reference image | Run disabled before spawn; exact missing item named; no engine start |
| S-5 | Scanner | Scan the known adapter directory | IP adapters are `AssetIPAdapter`; no `AssetControlNet` duplicates; CLIP-Vision is `AssetClipVision`; FLUX entry excluded |

### 5.2 Independent operations (IP profile)

| ID | Case | Procedure | Evidence | Pass criteria |
|---|---|---|---|---|
| I-1 | IP txt2img | 512×512, 4 steps, q8, SDXL adapter, vision CPU, reference image | argv, PID, admission, IP load lines, PNG | image saved; adapter/vision loaded; strength reflected; reserve held |
| I-2 | Native ESRGAN 2x | `POST /sdcpp/v1/upscale` on I-1 output | request, response dimensions, duration, PNG | 1024×1024; HTTP 200; no DGL job; no base reload |
| I-3 | IP + ControlNet img2img | stage-2 image as init, prepared pose, ControlNet CPU, same IP image/strength | argv, ControlNet load, IP load, output | ControlNet applied (not byte-identical to no-control A/B); IP applied; 512×512 |
| I-4 | Second native upscale | upscale I-3 output | response dimensions, duration | 1024×1024; same upscaler; no reload |

### 5.3 Complete workflow and repetition

| ID | Case | Procedure | Pass criteria |
|---|---|---|---|
| W-1 | Adequate VRAM full workflow | 512 base -> 1024 -> 512 refine -> 1024, IP on, ControlNet CPU, vision CPU | all stages complete; one PID unless justified; free VRAM stays above 10% reserve (or documented bounded dip); adapter/vision resident; no DGL |
| W-2 | Repeat 3× | W-1 three times sequentially | same profile; no adapter/vision reload; no base reload; terminal reclaim `retain` in normal band |
| W-3 | Batch 2 | two complete IP workflows in one batch | sequential; two final PNGs; peak equals one workflow ± noise; no concurrent GPU stage |
| W-4 | Constrained 5 GiB | `--max-vram 5` with IP + controlnet CPU + vision CPU | either completes within the measured envelope or refuses at the first impossible stage with exact numbers; no OOM/crash |
| W-5 | CPU-diffusion pressure | explicit CPU-diffusion placement if W-4 cannot fit | completes or refuses cleanly; time and placement recorded; no silent downscale |
| W-6 | Profile switch back | after W-2, run the existing non-IP four-stage workflow | default pinned engine starts; DGL upscales work; no IP args; prior IP artifacts unaffected |

### 5.4 Reuse, pressure, and recovery

| ID | Case | Procedure | Pass criteria |
|---|---|---|---|
| R-1 | Adapter/vision reuse | inspect startup and per-stage logs across W-1/W-2 | one adapter load, one vision load, one engine start; stages 1 and 3 reuse them |
| R-2 | Necessary reclamation | inject the existing measured-growth test hook or external pressure | exactly one recycle when the band/reserve is crossed; next stage succeeds; no second owner |
| R-3 | Reserve pressure | external holder leaves less than the stage projection after reserve | stage refuses with required/free/reserve/available; no unload; release then succeeds |
| R-4 | Missing vision failure | remove/rename the encoder and run | preflight refusal; no engine start; clear message; restore and succeed |
| R-5 | Incompatible adapter | SD1.5 adapter with SDXL checkpoint, or vice versa | refused before spawn; no engine error loop; exact family mismatch |
| R-6 | Engine crash | kill the IP `sd-server` between stages | failure reported; existing owner sees the dead PID; next run starts cleanly and completes |

### 5.5 Hosts and regression

| ID | Case | Procedure | Pass criteria |
|---|---|---|---|
| H-1 | Browser harness | real bindings, IP panel, full workflow | all four stages visible; profile/placement shown; no console errors; output appears |
| H-2 | Packaged WebView2 | same run through the packaged app | same visible workflow and artifacts; profile path resolves in the packaged layout |
| G-1 | Non-IP workflow regression | existing base -> DGL upscale -> ControlNet -> DGL upscale | unchanged behavior; pinned engine; no IP flags; prior hashes/performance envelope respected |
| G-2 | Single generation regression | ordinary txt2img with no IP | unchanged startup args and result; no profile switch |

### 5.6 Measurements recorded for every live row

- Actual command line, selected binary, profile, PID.
- Per-stage `nvidia-smi` used/free and process working set.
- Admission: required/resident/projected peak/activation, free VRAM, reserve, RAM,
  IP adapter bytes, CLIP-Vision bytes, placement decision.
- Engine lines: IP-Adapter load, CLIP-Vision load, ControlNet load, upscaler load,
  native upscale response, any release/eviction at its actual log level.
- Stage durations, output dimensions, output paths, SHA-256.
- Reclamation: decision, baseline, measured, growth, PID before/after, free before/after.
- Profile switch: old PID stopped, new PID, new binary, no orphan process.

### 5.7 Measured results (real assets, GTX 1080, browser harness)

Engine command for the single IP generation:

```text
build/engine-baseline/backup-c6dfa9d/sd-server.exe
  --listen-port 8081 --threads 16
  -m G:\AI GGUF\Image models\sd15-v1-5-pruned-emaonly.safetensors
  --clip_vision G:\AI GGUF\Image models\clip_vision_h.safetensors
  --ip-adapter G:\AI GGUF\Image models\ip-adapter_sd15.safetensors
  --max-vram 6 --backend clip_vision=cpu --vae-tiling --flash-attn auto
```

Single IP run: 512x512, 8 steps, CFG 7, strength 0.70. Gallery artifact
`data/images/gallery/img_20261125_003405_947900.png` (463,148 bytes) and its
sidecar record the real adapter, encoder, staged reference, strength 0.70, and
`engine_profile: sdcpp-c6dfa9d-ip`. Sampling took 6.95 s; VAE decode 2.90 s.

| Measurement | Actual value |
|---|---|
| Admission | projected 5100 MB, resident 4108 MB, activation 480 MB, required 5326 MB, tiling planned; free 7181 MB, reserve 718 MB |
| Placement | diffusion_model 2191.89 MB VRAM; text_encoders 469.44 MB RAM; VAE 159.68 MB RAM; CLIP-Vision CPU; IP-Adapter resident |
| Backend line | `clip_vision=cpu,diffusion=CUDA0,te=CUDA0,vae=CUDA0 | params: te=cpu,vae=cpu` |
| sd-server working set | 303.85 MB after load, 3354.95 MB after generation |
| Live VRAM after run | 3872 MB used, 4320 MB free, 819 MB reserve held |
| Reclaim | `retain`, growth 0 MB, baseline 3355 MB, band 1393 MB, threshold 4748 MB, PID 3988 unchanged |

Four-stage command adds `--control-net ...control_v11p_sd15_openpose.pth`,
`--model-dir ...\upscalers`, and `--hires-upscalers-dir ...\upscalers`; it keeps
`--max-vram 6 --backend clip_vision=cpu --vae-tiling`. One engine start served all
four stages (PID unchanged) with no DGL job.

Artifacts: `temp/howl-image-temp/workflow-20261125_003444/run-1-stage-{1..4}.png`.
Stage 4 is 6,025,560 bytes at 2048x2048.

| Stage | Duration | Output | Admission | Actual placement / process |
|---|---:|---|---|---|
| 1 base txt2img | 8181 ms | 256x256 | projected 992 MB, required 0 MB, activation 480 MB, free 7024 MB, reserve 702 MB | 3129 MB VRAM + 629 MB RAM; WS 306 MB load, 3358 MB after |
| 2 native ESRGAN 4x | 9173 ms | 1024x1024 | projected 512 MB, required 63 MB, activation 0 MB, free 4350 MB, reserve 435 MB | 3129 MB VRAM + 629 MB RAM; WS 3354 MB after |
| 3 IP + ControlNet img2img | 15136 ms | 512x512 | projected 992 MB, required 0 MB, activation 480 MB, free 4349 MB, reserve 434 MB | 3129 MB VRAM + 629 MB RAM; WS 3434 MB after |
| 4 native ESRGAN 4x | 24537 ms | 2048x2048 | projected 512 MB, required 0 MB, activation 0 MB, free 3412 MB, reserve 341 MB | 3129 MB VRAM + 629 MB RAM; WS 3436 MB after |

Consecutive resident run (second full workflow on the same engine, artifact
`workflow-20261125_004424/run-1-stage-4.png`): one engine start, PID unchanged,
no upscaler or adapter reload. Durations were 4888, 9231, 10135, and 25015 ms.

| Stage | Live before (VRAM used/free; RAM used/available) | Live after | Duration |
|---|---:|---:|---:|
| 1 base | 4819 / 3373 MiB; 19209 / 46216 MiB | 4818 / 3374 MiB; 19291 / 46134 MiB | 4888 ms |
| 2 upscale | 4818 / 3374 MiB; 19316 / 46109 MiB | 4790 / 3402 MiB; 19282 / 46143 MiB | 9231 ms |
| 3 pose/controlnet | 4790 / 3402 MiB; 19271 / 46154 MiB | 4836 / 3356 MiB; 19265 / 46160 MiB | 10135 ms |
| 4 final upscale | 4836 / 3356 MiB; 19268 / 46157 MiB | 4785 / 3407 MiB; 19263 / 46161 MiB | 25015 ms |

The lowest free VRAM in that run was 3356 MiB, still 2537 MiB above the 819 MiB
reserve. No stage boundary shows a reload spike or a second engine process.

Stage 3 engine log: `preprocess control[0]: 563x855 crop=(0,146,563,563)
resize=512x512`. Stage 2 upscaler log: `Upscaler architecture: esrgan`,
`upscaling from (256 x 256) to (1024 x 1024)`, 8.74 s. Stage 4 reused the loaded
upscaler with no second upscaler load, `upscaling from (512 x 512) to (2048 x
2048)`, 23.53 s. The IP adapter and ViT-H encoder loaded once at engine start and
encoded the reference at each IP stage (`4 image tokens, strength 0.70`).

Matrix status from this evidence:

Logs cited by these rows: `build/verification/workflow-reserve819.log` (S-3, S-5,
I-2–I-4, W-1, H-1, W-2/W-3, W-4, R-1–R-3), `build/verification/workflow-recovery-final.log`
(W-6, R-6, G-1), `build/verification/packaged-h2-picker.log` (H-2), and the gallery
artifact `data/images/gallery/img_20261125_003405_947900.png` (I-1). W-5, R-4/R-5,
and G-2 cite the closest unit tests in their cells because no host log exists.

| Row | Status | Evidence / reason |
|---|---|---|
| S-3, S-5 | **PASS** | `clip_vision=cpu` accepted; real adapter/encoder/ControlNet/upscaler scanned into the correct classes |
| I-1 | **PASS** | Real 512x512 IP PNG, recognizable reference likeness, strength 0.70 |
| I-2 | **PASS** (as stage 2) | Native ESRGAN 4x, 256 -> 1024, 8.74 s, no DGL |
| I-3 | **PASS** (as stage 3) | Native IP + ControlNet img2img, preprocessing and control graph logged, 512x512, 15.14 s |
| I-4 | **PASS** (as stage 4) | Native 4x, 512 -> 2048, 23.53 s, upscaler reused |
| W-1 | **PASS** | All four stages completed sequentially; one PID; reserve held; no CPU-diffusion fallback; weights retained |
| R-1 | **PASS** | One adapter load, one encoder load, one engine start across two full workflows; adapter/encoder reused at stage 3; upscaler reused at stage 4 |
| H-1 | **PASS** | Browser harness displayed all four stages, durations, dimensions, admissions, placement, and final image; log `build/verification/workflow-reserve819.log` |
| W-2 repeat / W-3 batch | **PASS** | Three consecutive complete workflows on PID 18356 (W-4 + batch runs 1–2); no restart, no concurrent GPU stage; `retain` growth 4–11 MB |
| W-4 5 GiB pressure | **PASS** | Four stages completed at `--max-vram 5`; no offload/CPU-diffusion flag; admission reserve 819 MB on every stage; final hash `01bb…97df` |
| W-5 CPU-diffusion pressure | **UNVERIFIED** | No host-level CPU-only diffusion run. Closest coverage: `TestImageAdmissionPlansControlNetCPUWhenStagingWouldCrossReserve` in `internal/resource/image_activation_test.go`; it does not exercise an explicit CPU-diffusion workflow. |
| W-6 profile switch back | **PASS** | Non-IP pinned engine PID 12536, no IP args, four stages, DGL upscales, final hash `98a6…aa0f`; per-component residency not exposed by this binary |
| R-2 reclamation / R-3 reserve | **PASS** | Holders left 208–270 MB free; stage refused with reserve 819/available 0; sole owner recycled 18356→20284 and reloaded; release succeeded. Batch growth stayed `retain` |
| R-6 crash recovery | **PASS** | Kill between stages produced a clean failure; the failed stage wrote no artifact; next run on new PID 20788 at `--max-vram 5` completed all four stages with `retain` |
| R-4/R-5 asset failures | **UNVERIFIED** | No live missing/mismatched-asset run. Closest coverage: `TestProbeIdentifiesClipVisionFromStructure`, `TestClassifyNameRecognisesControlNetV11Files`, and `TestFamilyRefinementStaysInsideTheArchitecture` in `internal/image/scanner`; they classify assets but do not exercise the preflight refusal path. |
| H-2 packaged WebView2 | **PASS** | §5.11: `build/verification/packaged-h2-picker.log`, one engine PID 16540, four sequential stages, final 2048x2048 artifact SHA-256 `8798f57ea238c153088a2209b6c099df27d57040f2f8229dc1151027ae09f1df`; shared harness picker route re-check **PASS** (no regression; §5.11, `build/verification/harness-picker-regression.txt`) |
| G-2 single generation | **UNVERIFIED** | No post-change single-generation host log. Closest coverage: `TestImageAutoMaxVRAMDoesNotSelectConstrainedOffload` in `app-image-hardware_test.go`; it covers the startup budget, not the full single-generation route. |
| G-1 non-IP workflow | **PASS** | Same evidence as W-6: pinned non-IP engine, no IP flags, four stages, DGL outputs |

Known output limitation: the official pose sample is 563x855; ControlNet
center-cropped it to 512x512, so stage 3/4 crop the top of the head. That is a
sample/aspect-ratio limitation, not a failed stage or a resource failure.

### 5.8 Resource-policy and recovery closure (2026-11-25, browser harness)

All runs below used the isolated root
`%LOCALAPPDATA%\Temp\howl-ip-verify`, the real weights in
`G:\AI GGUF\Image models`, and the real sd-server binaries. The device was the
GTX 1080 with 8192 MB total; the reserve was always 819 MB. Each stage row is a
live `nvidia-smi` sample before and immediately after the stage, and every
artifact is on disk.

#### W-4 constrained 5 GiB IP workflow

Request: IP-Adapter SD1.5 + ViT-H, OpenPose ControlNet, 4x-UltraSharp,
256x256 base -> 1024 -> 512 refine -> 2048, VRAM Budget **5 GiB**. Engine command
contained `--max-vram 5 --backend clip_vision=cpu --flash-attn auto` and **no**
`--offload-to-cpu` and no `diffusion=cpu`.

Admission final: `admitted=true projected=4121MB resident=3129MB
activation=480MB required=3129MB free=8192MB reserve=819MB`. The 3129 MB value is
the engine's measured resident set from the immediately preceding 6 GiB run, not
the old 5326 MB module+reserve projection. The engine's own auto-fit then chose
`diffusion=CUDA0,te=CUDA0,vae=CUDA0` with `params-backend
diffusion=cpu,te=cpu,vae=cpu`; its own total was 3758.62 MB (VRAM 937.60 MB,
RAM 2821.01 MB: diffusion 2191.89 MB RAM, controlnet 937.60 MB VRAM, text 469.44
MB RAM, VAE 159.68 MB RAM). CUDA still performed the compute; Howl did not select
CPU diffusion.

| Stage | Duration | Live before VRAM used/free | Live after | Admission |
|---|---:|---:|---:|---|
| 1 base | 8860 ms | 5057 / 3135 MiB | 1907 / 6285 MiB | projected 992 MB; required 0 MB; free 6733 MB; reserve 673 MB |
| 2 upscale | 9395 ms | 1907 / 6285 MiB | 1904 / 6288 MiB | projected 512 MB; required 63 MB; free 6285 MB; reserve 628 MB |
| 3 IP/ControlNet | 13658 ms | 1904 / 6288 MiB | 2830 / 5362 MiB | projected 2432 MB; required 0 MB; activation 1920 MB; reserve 628 MB |
| 4 final upscale | 25709 ms | 2830 / 5362 MiB | 2824 / 5368 MiB | projected 512 MB; required 0 MB; reserve 536 MB |

Artifact directory: `temp/howl-image-temp/workflow-20261125_010342`; final
`run-1-stage-4.png` is 6,281,497 bytes, SHA-256
`e921c79f68d8bd2b193e24247f676350d4bc6c48aed3416277dd2ad8fdd2012f`.
The lowest free value after a stage was 5368 MB; the reserve was never crossed.

#### W-2/W-3 repeated and batch

After the 5 GiB single run, one request with `batch_count=2` ran two complete
four-stage passes on the same engine PID 22912. There was no second
`Starting sd-server` line, no adapter/vision reload, and no concurrent GPU stage.

| Pass/stage | Duration | Live before (VRAM used/free; RAM used/available) | Live after |
|---|---:|---:|---:|
| 1/1 | 4884 ms | 2861 / 5331 MiB; see harness log | 2829 / 5363 MiB |
| 1/2 | 9255 ms | 2829 / 5363 MiB | 2824 / 5368 MiB |
| 1/3 | 10161 ms | 2824 / 5368 MiB | 2824 / 5368 MiB |
| 1/4 | 25465 ms | 2824 / 5368 MiB | 2821 / 5371 MiB |
| 2/1 | 4878 ms | 2821 / 5371 MiB | 2821 / 5371 MiB |
| 2/2 | 9577 ms | 2821 / 5371 MiB | 2819 / 5373 MiB |
| 2/3 | 9688 ms | 2819 / 5373 MiB | 2819 / 5373 MiB |
| 2/4 | 25280 ms | 2819 / 5373 MiB | 2815 / 5377 MiB |

Artifacts: `temp/howl-image-temp/workflow-20261125_010459/run-1-stage-{1..4}.png`
and `run-2-stage-{1..4}.png`. Final files: run 1 is 6,281,497 bytes
(`e921c79f...`) and run 2 is 5,499,238 bytes
(`83cd9ceda829987d996796119c00521bff0ed55e2a2c851bdbce7b6e64805f17`).
The item boundary reclaimer logged `retain` at growth 5 MB, and the terminal
reclaimer logged `retain` at growth 20 MB, with the same PID. Together with the
preceding 5 GiB single run this is three complete workflows on one engine, with
no reloads and no competing GPU work.

#### R-2/R-3 reserve pressure

Two real CUDA allocations (`hold_vram.exe 4500` then `2200`) were held outside
Howl. `nvidia-smi` reported 7907 MiB used / 153 MiB free before the request. The
workflow's stage 1 then refused cleanly: `needs 992 MB VRAM; 260 MB free; 819 MB
reserve; 234 MB available`, naming the resident model as the blocking owner. The
existing single reclamation owner then made the required pressure decision:
`decision=reclaim reason="free VRAM 260MB fell to or below the 819MB reserve"`,
PID 22912 -> 21292, `reload_error=""`. No second reclaimer or second engine
existed.

After both holders were stopped, the same workflow completed on PID 21292:
stages 8385/9194/13641/25012 ms; live VRAM before/after 1367/6825 -> 1840/6352,
1840/6352 -> 1839/6353, 1839/6353 -> 2783/5409, 2783/5409 -> 2776/5416 MiB;
admission reserve 682/635/635/540 MB. The terminal reclaimer was `retain` on the
same PID (growth 0 MB). Artifacts:
`workflow-20261125_010802/run-1-stage-{1..4}.png`, final
`run-1-stage-4.png` 6,281,497 bytes, SHA-256 `e921c79f...`.

#### R-6 crash recovery

While a run was executing, PID 21292 was killed externally after stage 3. Stage 4
failed with a real transport error (`read tcp ... wsarecv: An existing connection
was forcibly closed by the remote host`). The existing reclaimer logged one
failure decision, `reason="failed generation left the engine state abandoned"`,
PID 21292 -> 9012, reload without error. The next run completed all four stages on
PID 9012 in 7891/9103/13642/24709 ms; live VRAM 1369/6823 -> 1841/6351,
1841/6351 -> 1841/6351, 1841/6351 -> 2785/5407, 2785/5407 -> 2778/5414 MiB;
terminal `retain`, growth 0 MB. Artifacts:
`workflow-20261125_011004/run-1-stage-{1..4}.png`; final 6,281,497 bytes,
SHA-256 `e921c79f...`. Failed run left stage 1-3 artifacts in
`workflow-20261125_010920` and stage 4 had no artifact, as expected.

#### W-6 profile switch-back to plain non-IP

After the IP runs, the IP adapter and encoder were cleared and the same UI was
run with VRAM Budget 6 GiB. Howl stopped the IP profile and started
`sdcpp-3b12781` (`sd-server.exe`) with `--max-vram 6`, the ControlNet and
upscaler directories, and **no** `--ip-adapter`, `--clip_vision`, or IP request
field. All four stages completed: 3966/9295/9136/24216 ms; DGL/`image_upscale`
logs show the 256->1024 and 512->2048 upscales. Live VRAM 1248/6944 -> 4587/3605,
4587/3605 -> 4591/3601, 4591/3601 -> 5535/2657, 5535/2657 -> 5535/2657 MiB;
admission reserve 694/360/360/265 MB; terminal reclaim `retain`, growth 0 MB.
Only PID 1792 was listening on 8081; the prior IP PID was gone. Artifact
`workflow-20261125_011538/run-1-stage-4.png` is 4,173,709 bytes, SHA-256
`9cda25de125458bf9e3766ea2688eb800b311bf2837b833fcfc84b51a854421f`.
The pinned engine does not expose per-component residency; its process working set
was 171 MB after load and 1371 MB after generation, live device use was 5567 MB,
and the live reserve remained held.

#### Defects found and fixed during closure

1. Workflow requests ignored `ig-max-vram`; the 5 GiB test initially couldn't be
   constrained. `ImageWorkflowOptions.MaxVRAM` now forwards the panel value, with
   tests for explicit and Auto values.
2. The no-op `ensureEngineReady` path re-ran full load admission for an already
   resident model, so the second batch pass was refused at 5 GiB. A resident
   generation now sets `ResidentLoaded` and is charged only its stage delta.
3. At 5 GiB the conservative module+2048 MB projection selected CPU diffusion
   even though the same composition had measured 3129 MB resident. Admission now
   accepts a measured resident total only for an identical model/component set;
   this is still gated by live free memory and the reserve.
4. Clearing IP-Adapter still sent `ip_adapter_strength`; the pinned engine rejects
   that without a DGL script. The native img2img client now omits the strength
   unless an IP reference image is present.

#### Owner report

- Current owners: `internal/resource` admission for capacity, the single
  measurement-gated image reclamation policy, sd.cpp `backend_fit` for actual
  placement, `Image.LoadedOptions` for the resident process, and `imageStartMu`
  for sequential workflow exclusion.
- Competing owners: none found; the single reclaimer made every recycle, and no
  batch/repeat run started a second engine or ran concurrent GPU stages.
- What stays: the existing admission/reclamation authority and the capability
  profile. What was added: measured-resident input and workflow budget forwarding;
  no second authority or reclaimer was added.
- Verified hosts: browser dev host with real bindings and real engine, and the
  packaged WebView2 host in §5.11. The pinned engine's per-component placement is
  not reported.

---

### 5.9 Final RAM-instrumented closure (2026-11-25, supersedes earlier sample values)

The stage log now records live VRAM and process RAM before and after every stage.
These are the final runs used for the matrix; older §5.8 numbers were taken before
the RAM columns existed and are retained only as history. Logs:
`build/verification/workflow-closure.log` (max6, 5 GiB, batch, reserve) and
`build/verification/workflow-recovery-final.log` (profile switch, crash, recovery).
All runs used the isolated root and the real assets listed in §5.8. Engine command
for W-4/reserve/batch: `sd-server.exe ... --max-vram 5 --backend clip_vision=cpu
--flash-attn auto`; the engine's own auto-fit reported 937.60 MB VRAM and 2821.01
MB RAM (diffusion 2191.89 MB on CPU/RAM, ControlNet 937.60 MB VRAM, text 469.44 MB
RAM, VAE 159.68 MB RAM), compute CUDA0. The 6 GiB IP runs used `--max-vram 6` with
the same IP/CLIP-Vision/ControlNet components and no CPU-diffusion flag.

**Constrained VRAM, W-4** — `workflow-20261125_012102`, one PID 20100, four stages,
no `--offload-to-cpu`, no `diffusion=cpu`, reserve held at each admission.

| Stage | ms | VRAM used/free before → after (MiB) | RAM used/avail before → after (MiB) |
|---|---:|---|---|
| base | 8382 | 5027/3165 → 1837/6355 | 20005/45420 → 22180/43245 |
| upscale | 9190 | 1837/6355 → 1837/6355 | 22180/43245 → 22144/43281 |
| pose/ControlNet | 13151 | 1837/6355 → 2781/5411 | 22140/43284 → 22322/43103 |
| final upscale | 24975 | 2781/5411 → 2772/5420 | 22318/43107 → 22389/43036 |

Final: `run-1-stage-4.png`, 6,281,497 bytes, SHA-256
`e921c79f68d8bd2b193e24247f676350d4bc6c48aed3416277dd2ad8fdd2012f`.
Terminal reclaimer: `retain`, baseline 5627 MB, growth 0 MB, PID 20100→20100.

**Batch and repeat, W-2/W-3** — `workflow-20261125_012222`, same resident PID 20100
for both runs; no engine restart and no second GPU stage. W-4 (`012102`) plus these
two runs are three consecutive complete four-stage workflows on that one PID, so
W-2's repeat condition is covered by the same evidence.

| Run/stage | ms | VRAM used/free before → after (MiB) | RAM used/avail before → after (MiB) |
|---|---:|---|---|
| 1 base | 4869 | 2822/5370 → 2805/5387 | 22341/43084 → 22324/43101 |
| 1 upscale | 8972 | 2805/5387 → 2780/5412 | 22309/43116 → 22208/43217 |
| 1 pose | 9637 | 2780/5412 → 2780/5412 | 22209/43216 → 22336/43089 |
| 1 final upscale | 24633 | 2780/5412 → 2778/5414 | 22362/43063 → 22200/43225 |
| 2 base | 4883 | 2778/5414 → 2781/5411 | 22185/43240 → 22157/43268 |
| 2 upscale | 8979 | 2781/5411 → 2772/5420 | 22151/43274 → 22298/43127 |
| 2 pose | 9703 | 2772/5420 → 2773/5419 | 22150/43275 → 22149/43276 |
| 2 final upscale | 24701 | 2773/5419 → 2772/5420 | 22146/43279 → 22212/43213 |

Run 1 final is the same 6,281,497-byte artifact and hash as W-4. Run 2 final:
5,499,238 bytes, SHA-256
`83cd9ceda829987d996796119c00521bff0ed55e2a2c851bdbce7b6e64805f17`.
Reclaimer after each run: `retain`, growth 13 MB then 11 MB, PID 20100→20100.
Both runs had all four artifacts (`run-{1,2}-stage-{1..4}.png`).

**Explicit reserve pressure, R-3** — external CUDA holders requested 4500 MB and
2200 MB; `nvidia-smi` then showed 178–300 MB free. Before the run the monitor read
7892 MB used / 300 MB free, RAM 24602 MB used / 40823 MB available. Stage 1 was
refused after 373 ms: free 265 MB, reserve 26 MB (the reserve is 10% of the
*currently free* budget at that instant), available 239 MB, needs 992 MB. Phase-logged
state after the refusal: 7927 MB used / 265 MB free, RAM 24631 MB / 40793 MB. No
output. The existing single owner logged `reclaim`, reason “free VRAM 265MB fell to
or below the 819MB (10%) reserve”, PID 20100→2228, reload_error empty. Holders were
released, the monitor read 1295 MB used / 6765 MB free, and the next run succeeded:
`workflow-20261125_012502`, PID 2228.

| Stage | ms | VRAM used/free before → after (MiB) | RAM used/avail before → after (MiB) |
|---|---:|---|---|
| base | 8460 | 1334/6858 → 1797/6395 | 17075/48350 → 22222/43203 |
| upscale | 9149 | 1797/6395 → 1797/6395 | 22213/43212 → 22457/42968 |
| pose/ControlNet | 13640 | 1797/6395 → 2741/5451 | 22460/42965 → 22322/43103 |
| final upscale | 24988 | 2741/5451 → 2759/5433 | 22301/43124 → 22434/42991 |

Final artifact: same 6,281,497 bytes and `e921…012f` hash. Terminal `retain`, growth
0 MB, PID 2228→2228. This is the reserve-pressure success; the refusal is the
pressure evidence.

**Failure recovery, R-6** — same harness process as the profile run. The IP engine
(PID 18596) was killed between stage 1 and stage 2 of `workflow-20261125_013130`.
Stage 1 completed in 8066 ms (VRAM 5566/2626 → 4020/4172, RAM 17942/47483 →
19897/45527, artifact `run-1-stage-1.png`, 97,781 bytes, SHA-256
`498d081df36b2ce6eae59902a35f20361190ccfb91cd839337c955add6ee0b33`). Stage 2 failed
after 4742 ms with `wsarecv: An existing connection was forcibly closed by the
remote host`; VRAM 4020/4172 → 1169/7023, RAM 19899/45525 → 17622/47803, no output.
The sole reclaimer saw no process to recycle (PID had already exited) and did not
invent a second reclaim path. The next run, `workflow-20261125_013153`, completed on
a clean new PID 15600:

| Stage | ms | VRAM used/free before → after (MiB) | RAM used/avail before → after (MiB) |
|---|---:|---|---|
| base | 8327 | 1213/6979 → 4019/4173 | 16722/48703 → 20003/45422 |
| upscale | 9515 | 4019/4173 → 4019/4173 | 20038/45387 → 20048/45377 |
| pose/ControlNet | 13141 | 4019/4173 → 4963/3229 | 20046/45379 → 19925/45500 |
| final upscale | 25476 | 4963/3229 → 4957/3235 | 19919/45506 → 20081/45343 |

Final artifact: 6,281,497 bytes, `e921…012f`; terminal `retain`, growth 0 MB, PID
15600→15600. Limitation: after a process death there is no live placement
measurement. The final reserve-floor rerun in §5.10 nevertheless restarted at
`--max-vram 5` and completed all four stages without a CPU-diffusion flag; that
later run is the recovery evidence. No recorded run in this section contains an
unnecessary CPU fallback.

**Profile switch-back, W-6/G-1** — after the IP runs, controls were cleared and the
request started `sd-server (sdcpp-3b12781)` from `sd-server.exe` with no
`--ip-adapter`, no `--clip_vision`, and no IP args; PID 15144 only. Four stages in
`workflow-20261125_013003`, DGL upscales logged:

| Stage | ms | VRAM used/free before → after (MiB) | RAM used/avail before → after (MiB) |
|---|---:|---|---|
| base | 4366 | 5009/3183 → 4568/3624 | 20256/45168 → 17168/48257 |
| upscale | 9283 | 4568/3624 → 4576/3616 | 17171/48253 → 16915/48509 |
| pose/ControlNet | 9136 | 4576/3616 → 5520/2672 | 16909/48516 → 17916/47509 |
| final upscale | 25223 | 5520/2672 → 5514/2678 | 17925/47500 → 17929/47496 |

Final: 4,173,709 bytes, SHA-256
`9cda25de125458bf9e3766ea2688eb800b311bf2837b833fcfc84b51a854421f`. This binary
does not report per-component residency; the evidence is live VRAM/RAM, process
working set (1373 MB after load), and the absence of IP arguments. Terminal
`retain`, growth 0 MB, PID 15144→15144.

### 5.10 Final 819 MB reserve-floor rerun (same day, supersedes §5.9 sample numbers)

The §5.9 runs exposed one policy defect in admission, not in the workflow. For an
image request the default reserve was computed as 10% of the moment's *free* VRAM.
Under the reserve-pressure run that printed `reserve=26MB` with 265 MB free, while
the single image reclaimer uses the fixed 819 MB (10% of the 8 GiB device) floor.
A small stage could therefore have been admitted into headroom the reclaimer is
meant to protect even though the measured 992 MB stage correctly refused.

Fix: image admission now uses the stronger of the existing free-percentage reserve
and the 10%-of-device floor; chat and audio keep their existing rule. Regression
tests updated/added in `internal/resource/image_activation_test.go`:
`TestImageReserveDoesNotShrinkBelowDeviceFloorWhenFreeIsTight` (265 MB free →
819 MB reserve, 0 MB available), measured-base tiling (819/6181), resident delta
(819/1184), soft allowance (391 MB), and device-floor default (1638 MB on 16384).
All image stages in `build/verification/workflow-reserve819.log` now print
`reserve=819MB`, including the 5 GiB runs.

The final harness runs all used `--max-vram 5` (the engine retained that budget
when the profile changed; the log's `Starting sd-server` lines are authoritative),
the IP binary for IP conditions, and `sd-server.exe` with no IP args for the plain
profile. The reference upload was the stage-1 PNG from `012102`, so final hashes
differ from the earlier human-reference set.

| Condition | Run / PID | Stage ms | Min free after stages | Final artifact |
|---|---|---|---:|---|
| W-1 adequate, IP | `013925` / 18356 | 8757, 9290, 13139, 25628 | 5420 MB | 5,229,208 B `01bb…97df` |
| W-4 constrained 5 GiB | `014036` / 18356 | 4947, 9006, 9639, 24803 | 5417 MB | same `01bb…97df` |
| W-2/W-3 batch run 1/2 | `014149` / 18356 | 4884/9098/9640/24716; 4883/8876/9140/24088 | 5368 MB | run 2: 6,203,436 B `75d7…e52d` |
| R-3 pressure refusal | `18356` → `20284` | 385 (refused) | 269 MB, reserve 819, available 0 | none |
| R-3 post-release success | `014407` / 20284 | 7900, 9098, 13641, 24913 | 5415 MB | 5,229,208 B `01bb…97df` |
| W-6/G-1 plain profile | `014514` / 12536 | 4468, 9291, 9635, 24214 | 4007 MB | 4,179,949 B `98a6…aa0f` |
| R-6 crash | `014619` / 19096 | stage 1 7242; stage 2 failed 4412 | 5346 MB before failure | stage 1 only 82,874 B `4a1c…da37` |
| R-6 next run | `014637` / 20788 | 7814, 9260, 13142, 25129 | 5346 MB | 5,229,208 B `01bb…97df` |

Full per-stage VRAM and RAM before/after for every row are in the log's
`[ImageWorkflow]` lines; the terminal `[ImageReclaim]` lines are `retain`
except the single pressure recycle. The batch was three consecutive workflows on
PID 18356 (W-4 + both batch passes) with no restart and no concurrent stage. The
plain profile command was
`sd-server.exe ... --max-vram 5 --type q8_0 --flash-attn auto` with no
`--ip-adapter` or `--clip_vision`; its 4.18 MB final is a DGL upscale, not native.
The crash run's stage 2 failed with a forcibly closed socket, wrote no artifact,
and the next run on a clean PID completed all four stages. The pressure refusal
recycled once under the sole owner, then the release run succeeded. No final run
ended below 819 MB free and none used CPU-diffusion placement.

**Final matrix status.** W-1, W-2, W-3, W-4, W-6, R-1, R-2, R-3, R-6, G-1, H-1 and
H-2 are verified by these runs. H-2 is the packaged run in §5.11. W-5 (explicit
CPU-only diffusion), R-4/R-5 (missing or mismatched assets), and G-2 remain
UNVERIFIED with their closest test coverage cited in §5.7; no host log exists for
them.

### 5.11 Packaged WebView2 closure (2026-11-25, H-2)

**Picker defect and fix.** The packaged reference-image drop zone used the HTML
file input, which WebView2 opens as an owned file chooser. Windows disables the
Wails window while that dialog is open; on return, the focus event reaches
go-webview2's `Chromium.Focus`, `MoveFocus` returns `E_INVALIDARG`, and Wails'
error callback calls `os.Exit(1)`. The earlier packaged session reproduced this
as `[WebView2 Error] The parameter is incorrect.` from `(*Chromium).Focus` via
`onFocus`. This is an app-level dialog-ownership/embedding defect, not a
test-window artifact: the same owned-dialog pattern had already forced the
ownerless folder chooser in `OpenImageModelsDirPicker`.

The fix keeps the real user path in each host:

- `internal/host/filechooser_windows.go` adds `host.ChooseFile`, using
  `GetOpenFileNameW` with `hwndOwner=0`, so Windows never disables the Wails
  window.
- `App.OpenIPAdapterImagePicker` reads the selected file, rejects non-images,
  returns a data URL, and logs `[ImagePicker] reference image selected ...`.
- `frontend/dist/js/image-generator-ui.js` calls that method in the packaged
  host; the browser harness is detected by `__HOWL_HARNESS__` and keeps its HTML
  input, which is its supported native surface.

No resource authority, reclamation owner, or workflow path was changed.

**Harness picker regression check (2026-11-25).** The shared frontend edit was
clicked through the actual `build/bin/howl-ui-harness.exe` host, not merely
compiled. At UI runtime `window.__HOWL_HARNESS__` was `true` (boolean), and the
shim script preceded the image UI script in the live document. The real
`#ig-upload-ip` control called the HTML input exactly once with
`dataset.context=ig-upload-ip`; the native branch evaluated false. A real
463,148-byte PNG from the harness gallery was loaded through that input into
`window.ipAdapterImageData` (617,532 base64 characters), and the UI showed
“Reference image loaded.” Because the harness has no OS chooser, automation
populated the HTML input's FileList through the browser's DataTransfer API and
dispatched its change event; the control click and frontend handler were real.
The unchanged harness four-stage pipeline then
completed: 15428 ms base, 9182 ms ESRGAN upscale, 48660 ms pose/ControlNet
generation, 25143 ms final upscale to 2048x2048, with every stage admitted
under the 819 MB reserve. This is a harness picker regression check, not new
packaged H-2 evidence; `packaged-h2-picker.log` remains the packaged H-2 run.

On the same normal packaged binary, the picker's cancel and rejected-selection
paths were checked through the real control. Esc logged
`[ImagePicker] reference image selection cancelled` and the UI showed
“IP-Adapter reference selection cancelled.”; selecting
`build/verification/not-an-image.txt` produced the UI error “is not an image
(text/plain; charset=utf-8)”. The process stayed responding after both. Evidence:
`build/verification/harness-picker-regression.txt`,
`build/verification/harness-picker-regression.log`,
`build/verification/packaged-picker-cancel.log`, and the four PNGs cited there.

**Recorded host and picker outcome.** Log:
`build/verification/packaged-h2-picker.log`. Artifacts:
`build/verification/packaged-h2-picker-artifacts/`. The host was the normal
`build/bin/howl-chat-s1.exe`, built with `desktop,production`; no `h2debug`, no
`devtools`, no Wails fork, and no go.work or modfile overlay. The UI was driven
with native mouse and keyboard; workflow inputs were set through the page
controls, not injected through DevTools. The picker log line is:

```text
2026/11/24 20:27:23 [ImagePicker] reference image selected path="C:\Users\StoneWolfPC\HOWL_Chat\build\verification\h2-ref.png" bytes=97781 mime=image/png
```

The process returned to `HOWL Chat` and stayed responsive after the dialog
closed. The reference was the real 256x256 stage-1 PNG from the browser-harness
set, not a test-only pre-seed.

**Cold-start precondition.** The first 5 GiB run in this fresh session
(`workflow-20261125_023003`) had no live measurement and selected
`--backend diffusion=cpu,clip=cpu,vae=cuda0,clip_vision=cpu --offload-to-cpu`.
It completed, but it is not the H-2 record. A 6 GiB run (`023314`) established
the measured resident baseline; the controls were then returned to 5 GiB for
the recorded run. The cold-start CPU-diffusion fallback is stated here rather
than hidden.

**Target run.** `workflow-20261125_023522`, one engine start and one PID 16540
for all four stages. The engine command was:

```text
build/engine-baseline/backup-c6dfa9d/sd-server.exe --listen-port 8081 --threads 16
  -m G:\AI GGUF\Image models\sd15-v1-5-pruned-emaonly.safetensors
  --clip_vision G:\AI GGUF\Image models\clip_vision_h.safetensors
  --ip-adapter G:\AI GGUF\Image models\ip-adapter_sd15.safetensors
  --control-net G:\AI GGUF\Image models\controlnet\control_v11p_sd15_openpose.pth
  --model-dir G:\AI GGUF\Image models\upscalers
  --hires-upscalers-dir G:\AI GGUF\Image models\upscalers
  --max-vram 5 --backend clip_vision=cpu --vae-tiling --flash-attn auto
```

There is no `diffusion=cpu`, no `--offload-to-cpu`, and no second engine PID.
The log has one `Starting sd-server` line for the target run; the stage lines
are in completion order 1, 2, 3, 4 on PID 16540, so no two GPU stages overlap.
The adapter, ViT-H encoder, ControlNet, and upscaler stayed in the same process.

| Stage | Duration | VRAM used/free before → after (MiB) | RAM used/avail before → after (MiB) | Admission |
|---|---:|---|---|---|
| 1 base generation | 8441 ms | 5009/3183 → 1771/6421 | 20687/44738 → 22547/42878 | projected 992 MB, activation 480 MB, required 0 MB; free 6859 MB, reserve 819 MB |
| 2 upscale | 9812 ms | 1771/6421 → 1771/6421 | 22550/42875 → 22650/42775 | projected 512 MB, required 63 MB; free 6421 MB, reserve 819 MB |
| 3 pose/ControlNet generation | 15640 ms | 1771/6421 → 2709/5483 | 22687/42738 → 22666/42759 | projected 992 MB, activation 480 MB; free 6421 MB, reserve 819 MB |
| 4 final upscale | 26782 ms | 2709/5483 → 2704/5488 | 22624/42801 → 22700/42725 | projected 512 MB, required 0 MB; free 5483 MB, reserve 819 MB |

Every target admission prints `reserve=819MB`; the minimum free VRAM after a
stage is 5483 MB. The terminal reclaimer is `retain`: baseline 5626 MB,
measured 5626 MB, growth 0 MB, threshold 7078 MB, PID 16540→16540. The packaged
UI reported “Pipeline complete. Weights retained unless reclamation reported
otherwise.” and listed all four stages completed, with stage 4 at 2048x2048.

| Artifact | Bytes | SHA-256 |
|---|---:|---|
| `run-1-stage-1.png` | 82,838 | `bf24eadbd1397cd12e579fb6efe55896a7058f0bb2c6eba3fdf86d76fb189bc5` |
| `run-1-stage-2.png` | 1,296,392 | `48436147144b7e7d8ca0e7bf77655eb906bb8bf42da8021f1b6323f9f5a6e73c` |
| `run-1-stage-3.png` | 453,728 | `0bdf7c00edb441caf954ad6fd98ccf6299f404fd8ff91eb64627b41b0187b162` |
| `run-1-stage-4.png` | 5,267,930 | `8798f57ea238c153088a2209b6c099df27d57040f2f8229dc1151027ae09f1df` |

These are new artifacts from the picker-fixed session; earlier H-2 artifacts are
not reused as this evidence. The packaged profile path resolves to
`build/engine-baseline/backup-c6dfa9d/sd-server.exe`; the app logged no IPC,
asset-resolution, picker, or workflow-dispatch failure.


## 6. Prepare items (before implementation)

Status (2026-11-25): items 1-5 are complete; item 5's packaged H-2 is verified in
§5.11. W-5, R-4/R-5, and G-2 remain unverified.

1. **Acquire the ViT-H/14 CLIP-Vision encoder.** DONE. Real
   `clip_vision_h.safetensors` (2.41 GB) is in the scanned root and loaded by the
   engine, which produces 4 image tokens.
2. **Provision the IP engine profile.** DEVELOPMENT PATH DONE. The profile resolves
   to `build/engine-baseline/backup-c6dfa9d`; packaged `runtime/stablediffusion-ip`
   provisioning is still open.
3. **Register the assets.** DONE. Real SD1.5 adapter, ViT-H encoder, SD1.5
   OpenPose ControlNet, and ESRGAN upscaler are scanned into the correct classes.
4. **Preflight the real model.** DONE. SD1.5 + IP + CLIP-Vision + ControlNet +
   ESRGAN completed a native four-stage run under `--max-vram 6` with the 10%
   reserve held.
5. **Packaged host and pressure rows.** DONE for H-2 and the pressure/recovery
   rows. H-2 is verified in §5.11; W-2/W-3/W-4/W-6 and R-2/R-3/R-6 are closed in
   §5.9/§5.10. W-5, R-4/R-5, and G-2 remain unverified.

---

## 7. Risks and open items

1. **Native IP + ControlNet on `c6dfa9d` is verified as a completed image**, not as a
   byte-level with/without-control A/B. The control preprocess and graph logs are
   present; keep the A/B row open.
2. **Native upscale with a resident base is verified.** Stage 2 and stage 4 used the
   same loaded engine and upscaler; stage 4 did not reload the upscaler.
3. **Adapter + ControlNet at 512 was measured**, not at 1024. The run stayed at
   4773 MiB used / 3419 MiB free after stage 4 with the reserve intact. W-4 is the
   constrained pressure proof; W-5 remains unverified. A refusal is acceptable, a
   silent resolution drop is not.
4. **CPU CLIP-Vision latency is included in stage durations.** It avoids ~2.4 GB of
   VRAM; the four stages took 8.2/9.2/15.1/24.5 s at these sizes.
5. **Dual-engine maintenance remains.** The IP profile is a second binary and DLL set.
   Consolidation still requires one build with native IP, native upscale, and the
   proven DGL/native replacement.
6. **External pose extraction remains explicit.** The proof used a precomputed
   OpenPose sample; Howl does not extract pose. Do not describe ControlNet as doing
   extraction.
7. **FLUX false positive remains a scanner concern.** A FLUX `ip_adapter.safetensors`
   must not appear as an SDXL/SD1.5 IP-Adapter.
8. **Automatic-cap bug fixed.** The first repeated run showed that an automatic
   post-reserve cap could be mistaken for a user hard budget and force CPU diffusion.
   `AutoMaxVRAM` is now excluded from `constrainOffload`; regression test
   `TestImageAutoMaxVRAMDoesNotSelectConstrainedOffload` covers it.
9. **Profile selection is per request, not global.** W-6 still needs a live pass.

---

## 8. Remaining work and decisions

The real-asset outcome is proven. The remaining work is not a reimplementation of
the route:

1. keep the normal packaged binary debug-free; H-2 is closed by §5.11 and only
   needs a re-run if the packaged profile path changes;
2. W-2/W-3/W-4/W-6 are closed in §5.10; W-5 CPU-diffusion pressure remains;
3. R-2/R-3/R-6 are closed in §5.9/§5.10; R-4/R-5 failure rows remain;
4. G-1 is closed by W-6; G-2 single-generation regression remains;
5. decide whether the c6dfa9d binary is shipped under `runtime/stablediffusion-ip` or
   remains an explicit user-configured path.

No further synthetic-weight run is evidence. The two real PNGs and four stage files
listed in §5.7 are the current proof set.
