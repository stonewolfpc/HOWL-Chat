# Image Hardware Optimization Plan — 2026-11-24

**Status:** implemented and exercised on 2026-11-24. Both supported hosts were driven
against the real `sd-server` on the final tree: the browser dev host through its image
controls, and the packaged WebView2 app with native mouse clicks. The final backend/TTL
audit spawned real engines for every exposed value, read the
actual OS command line, and removed Vulkan and staging TTL because this engine does not
honor them. §6 records five sequential images in one process, a two-image batch,
512/1024/1344×1024 resolutions, 3 GiB partial placement, 2 GiB automatic CPU-only,
and recovery back to 3 GiB in the same session. It also records the pinned engine's
missing per-component residency line and the OS working-set measurements that replace
the former permanent “unverified” card. Automatic placement is now re-derived when the
budget or model changes instead of being carried as saved intent. No commits were made.

**Scope:** stable-diffusion.cpp `sd-server` as shipped in this repository
(`master-1122-3b12781`, commit `3b12781`) and Howl's image execution path.

**Authorities:**

- Admission policy: `internal/resource` (`estimateImage`, `Authority.Admit`).
- Engine placement: sd.cpp auto-fit, graph-cut and VAE tiling.
- Observed placement/events: `internal/image/engine/sdcpp/placement.go` and the
  engine's own log.
- Model load flags: `internal/image/engine/sdcpp/engine.go`.
- User intent: the image page controls and `ModelLoadOptions`.

No second planner is introduced. The app decides policy and passes a budget; sd.cpp
decides which module lives where. The one exception is the measured constrained mode,
where the app explicitly assigns diffusion to CPU and the VAE to CUDA because no
auto-placement or ordinary offload combination holds a small budget (§6.4).

---

## 0. Implementation report (owners and changes)

**Current owners**

- Admission and reserve policy: `resource.AdmissionRequest` / `Authority.Admit`.
- Runtime image start options: `ModelLoadOptions`, built once by
  `imageLoadOptionsForGeneration` and passed to both preflight and the spawn.
- Actual loaded state: `Image.LoadedOptions` is the runtime authority for the
  process that is running. It is not saved user intent.
- Startup-flag reconciliation: `imageStartupOptionsDiffer` and
  `preserveImageLoadedOptions` in `app-image-model.go`.
- Engine placement and its VAE retry: sd.cpp, observed through
  `GetImagePlacement`.

**Competing owners found and removed**

- The saved request was being treated as the loaded reality. An auto-q8 process was
  credited at q8 but charged again at f16 on the next generation, which refused a
  512×512 run by 6 MB. `LoadedOptions` is now the authority for a resident process.
- A generation against a resident model silently ignored changed startup controls
  (offload, weight type, backend, FA, tiling, max-vram). Those options now force a
  real `LoadModel` restart; unchanged options still reuse the process.
- Replacement credit used a **registered** image model even after its server was
  killed, choosing an 8 GiB cap on an 8 GiB card. Replacement credit now requires
  `Engine.IsRunning()`.
- Automatic placement was stored alongside the saved request, so changing 3 GiB to
  2 GiB reused the 3 GiB partial plan, skipped re-planning, and admitted at 2048 MB
  instead of selecting CPU-only. It is now kept only when the same model is running at
  the same budget; every other request re-derives it.
- `--offload-to-cpu` was charged as a fixed 2 GiB even though sd.cpp opportunistically
  loads modules into VRAM. Ordinary offload still cannot honor a sub-model budget: the
  direct `--max-vram 3 --offload-to-cpu` run peaked at 6.37 GiB total. The per-module
  CPU-diffusion mode does bound it, so `applyConstrainedOffload` selects that mode when
  an explicit offload budget is below the resident need and refuses below its measured
  conservative 3 GiB floor after the measured 2636 MB process need. The UI warning and
  §6.4 state the tradeoff.
- The explicit positive `--max-vram` path had disabled the 10% reserve. The reserve is
  now unconditional for image admission; `--max-vram` only limits spend.
- Batch metadata was never populated, so a two-image batch wrote only the first image
  to the gallery. Each sequential item now carries its own metadata row.
- Empty Backend/Params Backend and a blank residency field were treated as “keep the
  previous value”, so the UI could show Auto/blank while the resident engine still ran
  `disk` and `--weight-residency-ttl 3`. Generation requests now assign empty placement
  through to the engine, and a blank TTL sends `-1` (disabled).

**What stays and why**

- sd.cpp remains the placement and VAE-retry authority.
- Howl remains the policy and admission authority.
- Auto q8_0, auto tiling and the offload cap are runtime decisions only; none are
  written back as user intent.

---

## 1. Measured evidence

### 1.1 Engine and model

- `sd-server --help`: 232 lines, captured at `build/tmp/research/sd-help-current.txt`.
- Model: `Illustrious-XL-v2.0.safetensors`, 6616.63 MiB file.
- Engine-reported weights:
  - diffusion 4900.07 MiB
  - conditioner 1564.36 MiB
  - VAE 159.68 MiB
- Safetensors header split (measured, header only):
  - diffusion 4897.05 MiB
  - conditioner 1559.68 MiB
  - VAE 159.56 MiB

The header split matches the engine's component report within ~1 %, so model
components can inform admission before the process starts.

### 1.2 Workload measurement on the 8 GiB GTX 1080

Run: `--max-vram -1`, Flash Attention auto, two sampling steps.

| Workload | Result | Engine evidence |
|---|---|---|
| 512×512, batch 1 | success | no VAE retry |
| 1024×1024, batch 1 | success after auto-retry | `VAE decode compute failed`; `VAE decode ran out of memory; retrying with spatial tiling`; tiles reduced 1024→512; 9 tiles |
| 512×512, batch 2 | success | Howl's client runs batch items sequentially; no concurrent two-image peak |

- `--max-vram -1` measured 6.99 GiB free, reserved 1.00 GiB, used 5.99 GiB.
- At 1024×1024 the decode asked for 8798.47 MB device / 7774.47 MB budget
  against 6596.71 MB available, then tiled and completed.
- At 512×512 the same model decoded without the retry.

This is evidence for two decisions:

1. sd.cpp's automatic VAE tiling is real and is the last-resort executor. Howl must
   not replace it with a second retry implementation.
2. Howl can pre-empt the known failure by enabling `--vae-tiling` when its
   workload estimate says the decode will exceed the remaining budget. That saves
   the failed attempt without changing the engine's behaviour.

### 1.3 Reclamation API and execution backends

- The shipped server exposes no public unload/free endpoint. The HTTP surface is
  `/sdapi/v1/*` and `/sdcpp/v1/*`; strings show internal
  `sd_dgl_ctx_release_staging` / `sd_dgl_ctx_unload_models`, not routes.
- `/sdapi/v1/options` returns only `samples_format` and `sd_model_checkpoint`, so
  the TTL values cannot be changed after launch.
- **Residency TTL is real.** With `--weight-residency-ttl 2`, a 64×64 generation
  finished at working set 4357 MiB; eight idle seconds later the same process was
  361 MiB (private 6679 → 2362 MiB). `-1` kept 4356 MiB resident for eight seconds.
  `0` was already at 360 MiB when the response returned. A following generation at
  the same TTL reloaded tensors and succeeded.
- **Staging TTL is not observable at a generation boundary.** `--weight-staging-ttl
  2` and `0` both finished with 4230 MiB GPU in use and 4335 MiB working set, unchanged
  after eight idle seconds, with or without `--offload-to-cpu`. sd.cpp already
  releases staged GPU copies when a component's execution ends, so a post-generation
  staging timer has nothing left to reclaim. The UI control and the automatic
  `--weight-staging-ttl 0` were removed rather than presented as memory savings.
- **Execution backends measured through Howl's spawn path, actual OS command line
  and engine logs:**
  - compute `--backend cuda`: accepted; `flash attention auto: CUDA0`; generation
    succeeded.
  - compute `--backend cpu`: accepted; `flash attention auto: CPU`; generation
    succeeded.
  - compute `--backend vulkan`: rejected by this binary (`cannot load component`),
    because only CUDA and CPU backend DLLs are shipped. The Vulkan option is gone.
  - params `--params-backend cpu`: accepted; generation succeeded. It is not a
    no-VRAM guarantee: the sustained max3 run peaked at 6.18 GiB total (§6.4).
  - params `--params-backend cuda`: accepted; generation succeeded.
  - params `--params-backend disk`: accepted; auto-fit disabled, and the engine
    logged 1680/2515 tensors loaded from the model file during sampling before the
    image completed. It remains a manual last resort.
- **Between-generation reclamation is now measurement-gated.** The pinned build
  still has no unload/free endpoint, so the strongest real lever is a controlled
  process recycle followed by a reload of the already-admitted options. It is not
  an unconditional per-generation restart: the normal measured band retains the
  resident weights and the same PID. The final policy and evidence are in §6.11.
  A recycle is heavier than an ideal unload endpoint; that cost is paid only when
  the working set or the 10% reserve says it is warranted.

### 1.4 Verification on this machine (real engines, isolated root)

The 512/1024/batch/offload runs used `Illustrious-XL-v2.0` and two sampling steps.
The backend/TTL audit used the same checkpoint at 64×64 and one step to isolate
backend acceptance and timer behaviour. All runs used isolated gallery
`C:/Users/StoneWolfPC/howl-image-hw-iso-20261124/images` and the real `sd-server`.

| Case | Observed command / result |
|---|---|
| 512×512, sequential twice | First run auto-refused f16 and loaded q8_0. Second run reused the same process and succeeded; no second `Starting sd-server` line. |
| 512×512 on resident q8 | `CheckImageFit` reported `required_mb=5084`, not the saved f16 6945. |
| 1024×1024 on resident q8 | Admission admitted; engine logged `VAE decode ran out of memory; retrying with spatial tiling`; image saved. UI reported `vaeTiled=true` with the engine line. |
| 1024×1024 fresh load | Auto q8_0 plus pre-enabled tiling: `--max-vram -1 --type q8_0 --vae-tiling --flash-attn auto`. Image saved without a Howl retry. |
| Batch 2 at 512×512 | `Batch: generating 2 images sequentially`; two gallery PNG/JSON pairs written after the per-item metadata fix. |
| Offload, no explicit cap | `--max-vram 6 --flash-attn auto --offload-to-cpu`; admission logged 2048 MB required, 718 MB reserve, and the auto-cap warning. Image saved. Post-generation GPU fell to 2.62 GiB. |
| Offload + explicit 3 GiB cap, ordinary path | Direct engine probe: `--max-vram 3 --type q8_0 --offload-to-cpu` peaked at 6.37 GiB total (process ≈5.47 GiB), leaving 1.82 GiB free. This path alone does not honor the small budget. |
| Offload + explicit 3 GiB cap, constrained mode | `--max-vram 3 --backend diffusion=cpu,clip=cpu,vae=cuda0 --type q8_0 --vae-tiling --offload-to-cpu`; 512×512 succeeded in 23.51 s, peak total 2089 MiB (process ≈1197 MiB), 6103 MiB free. Batch 2 and 1024×1024 peaked at 2087/2083 MiB; repeat reused PID 25268. See §6.4. |
| Offload + 2 GiB cap (earlier pass) | Refused in 0.04 s before any spawn: `VRAM Budget 2 GiB is below the measured 2636 MB CPU-diffusion process need (conservative 3072 MB admission floor...`; the next 3 GiB request started and succeeded. The final ladder replaces this with automatic CPU-only execution; see §6.10. |
| Max VRAM = 0 | Admission kept the 10% reserve, warned that 0 means live free memory, and spawned `--max-vram 6 --type q8_0` instead of the engine's uncapped 0. Image saved. |
| FA control | Changing Auto to On forced a restart and produced `--flash-attn on`; it was not inert. |
| Forced recovery | Killed `sd-server.exe` (PID confirmed), then generated again. The next start succeeded with a 6 GiB auto cap, not the previous dead model's credit. |
| Post-generation residency | Default (no TTL) held weights in RAM for reuse; GPU use stayed near the cap. `--weight-residency-ttl` was never set automatically. |
| Backend/TTL audit | Real `sd-server` commands: CUDA, CPU, params CPU/CUDA/disk succeeded; Vulkan failed; residency TTL 0/2 unloaded and reloaded; staging TTL 0/2 changed nothing. Details in §1.3. |
| UI residency TTL = 2 | Actual OS command line: `--max-vram 6 --backend cuda --params-backend cpu --weight-residency-ttl 2 --flash-attn auto`. Image saved; after idle the process was 341 MiB working set. The next generation reloaded tensors and saved without restarting the server. |
| UI TTL 2 → 3 | New `Starting sd-server` with `--weight-residency-ttl 3` and a new PID; image saved. A startup-flag change is a real restart. |
| UI clears Params/TTL | Params empty and TTL blank restarted to `--backend cuda --type q8_0 --weight-residency-ttl -1` with no `--params-backend`; image saved. Empty means Auto, blank means disabled, not “keep the old value”. |
| UI compute CPU | Actual command line `--max-vram -1 --backend cpu --weight-residency-ttl -1 --flash-attn auto`; engine logged `flash attention auto: CPU`; image saved. |
| UI params disk | Restarted to `--backend cuda --params-backend disk --weight-residency-ttl 3`; image saved. |
| Final params CPU/CUDA | `--params-backend cpu` and `--params-backend cuda` each spawned, generated and saved on the final build. |
| Packaged WebView2 | Native clicks opened the image page and Generate. Log: `--max-vram 6 --type q8_0 --flash-attn auto`; UI showed `Done!` and the 512×512 result; gallery PNG/JSON written. |

Measured live values during these runs: 8 GiB GTX 1080; an idle desktop baseline of
about 892 MiB; 6.47 GiB total at the 512×512/6 GiB admission peak; 2.08 GiB total
at the constrained 3 GiB CPU-diffusion peak (1.19 GiB process, 6.11 GiB free). The
explicit 3 GiB flag is a per-device managed-weight/runner budget, not a hard process
total for ordinary GPU/offload work (6.37 GiB peak); the bounded mode comes from the
explicit per-module CPU backend plus VAE tiling. See §6.4.

---

## 2. Capability matrix

Legend:

- **Used** — already reaches the engine and is meaningful.
- **Implement** — shipped capability Howl can use, with a clear benefit.
- **Engine-owned** — the engine already performs it; Howl must observe/report, not duplicate.
- **Manual** — exposed only as an explicit trade, not automatic.
- **Reject** — evaluated and deliberately not implemented; reason recorded.

| Capability / flag | Decision | Reason and Howl behaviour |
|---|---|---|
| Auto-fit / graph-cut (`--auto-fit`) | Used | Default on. Admission passes the budget; sd.cpp chooses placement. |
| VRAM budget (`--max-vram`) | Used, not a hard total alone | Wired as sd.cpp's per-device budget for managed weights and runner buffers. A value above the post-reserve budget is clamped with a warning; blank Auto is also capped post-reserve (sd-server's own -1 reserves only 1 GiB). Ordinary GPU/offload work can still exceed it (6.37 GiB at max3); the constrained mode below is what bounds the process. |
| Automatic offload budget | Implement as a spend hint | With `--offload-to-cpu`/`--params-backend cpu` and no explicit cap, set a positive `--max-vram` from post-reserve free VRAM. It does not bound the process total; offload loads modules into VRAM when needed. |
| Constrained low-VRAM mode | Implement, automatic ladder | With Backend/Params on Auto and a budget below resident weights + decode need, Howl selects the measured bounded mode at ≥3 GiB: `--backend diffusion=cpu,clip=cpu,vae=cuda0`, VAE tiling, offload. Below the 3072 MB conservative floor it selects **CPU-only** (`--backend cpu`) instead of refusing, so any positive budget starts without GPU weights; expected cost is roughly 18 s per 512 one-step image on this box. Explicit Backend/Params are never rewritten. |
| Actual allocation report | Implemented; engine residency unavailable on this build | The pinned `3b12781` binary does not print its `total params memory size` line. Howl records the sd-server process working set after readiness and after every generation (`internal/proc/memory_windows.go`, `footprint.go`) and shows it beside the plan. It explicitly says component residency is unverified rather than inventing it. Measured 171 MB after readiness (lazy load) and ~7.0 GB after generation for this SDXL model. |
| 10 % free-VRAM reserve | Implement, unconditional | Image admission always keeps 10% of effective free VRAM. Explicit `--max-vram` limits spend but cannot cancel the reserve. |
| RAM residency (`--params-backend cpu`, `--offload-to-cpu`) | Used / automatic fallback | User choice remains; auto-fit already stages modules in RAM. Post-generation working set falls to ~361 MiB, but generation peak still includes the modules sd.cpp loads into VRAM. |
| Disk residency (`--params-backend disk`) | Manual | Last resort, minutes per image; expose in Params Backend, never choose automatically. |
| Compute backend (`--backend cpu\|cuda`) | Manual | Both were spawned and generated successfully; CPU is slow but truthful. No automatic CPU encoder because family-specific defects exist. |
| Compute backend (`--backend vulkan`) | Removed | This binary has no Vulkan backend; the engine rejected the flag with `cannot load component`. The UI option was removed. |
| Parameter backend (`--params-backend cpu\|cuda\|disk`) | Manual | All three were spawned and generated successfully. Disk disables auto-fit and reloads from the model file during sampling; it is labelled last resort. |
| On-load quantization (`--type`) | Implement automatic q8_0 fallback | Only when the stored type does not fit, only when the user left Weight Type on Auto, only q8_0. Never re-quantize an explicit user type. Warn. |
| VAE tiling (`--vae-tiling`) | Implement pre-emption + keep controls | Engine already retries automatically. Howl pre-enables tiling when its workload estimate exceeds headroom, and surfaces the engine retry when it happens anyway. |
| VAE tile size / relative size | Reject automatic | No measured reason to override sd.cpp's 256×256 default. Keep engine default. |
| VAE tile overlap | Used | Per-request option already reaches the engine for DiT families. |
| Cheap previews (`--taesd`, `--taesd-preview-only`) | Used | Already attached when present. |
| Idle staging release (`--weight-staging-ttl`) | Not exposed | Measured 0 and 2 with and without offload: no post-generation GPU or RAM change. sd.cpp releases staged copies when a component finishes, so the control was removed. |
| Idle residency release (`--weight-residency-ttl`) | Manual, verified | 2 s dropped the process from 4357 to 361 MiB working set and the next generation reloaded; 0 was immediate; -1 kept weights resident. Never automatic: unloading destroys reuse. |
| Weight type explicit (`f16`, q8, q4, q5) | Used | User choice always wins, including f16. |
| Tensor type rules | Reject | Expert surface; no family-independent rule. |
| Flash Attention (`--flash-attn on\|off\|auto`) | Implement truthful tri-state | Current checkbox cannot express off; unchecked means auto. Make Auto/On/Off explicit and pass it. |
| Sage attention (`--sage-attn`) | Reject automatic | CUDA-specific, shape-dependent, speed feature not a memory guarantee. |
| Conv direct (`--conv-direct`) | Reject automatic | Engine's auto chose off/on per module; no memory guarantee. |
| Disable prefetch (`--disable-prefetch`) | Reject automatic | Reduces overlap memory but costs speed; no measured win on this box. Not exposed. |
| Disable segmented compute | Reject | Increases peak memory; opposite of goal. |
| Eager load | Reject | Increases peak memory; opposite of goal. |
| Batched CFG (`--batched-cfg`) | Reject automatic | Engine default; no measured memory reason to change. |
| Denoise cache (`--cache-mode`, `--cache-option`) | Used | Speed/cache controls already wired; not a residency guarantee. |
| Conditioning cache size | Used | Already wired. |
| mmap | Used | Already wired. |
| Multi-GPU split / RPC | Reject automatic | This target has one GPU; no verified multi-GPU behaviour. |
| `--max-pending-jobs`, completed TTL | Not applicable | Server-only queue knobs; Howl owns its job queue. |
| `/sdapi/v1/options` TTL mutation | Rejected by engine | Endpoint exposes no TTL fields; cannot be used for reclamation. |

---

## 3. Implementation shape

### 3.1 Truthful controls

1. Replace the image Flash Attention checkbox with `Auto / On / Off`.
   - Legacy `flash_attention: true` maps to `on`; `false` maps to `auto` (existing
     behaviour) only when no mode is supplied.
2. Send `--flash-attn <mode>` explicitly, including `auto`.
3. Wire `weight_residency_ttl` through:
   - frontend read/save/restore,
   - `ImageOptions` and `ModelLoadOptions`,
   - `parseModelLoadOpts`,
   - `imageLoadOptionsForGeneration`,
   - `buildServerArgs`.
   Staging TTL is not modelled or forwarded: it has no generation-boundary effect.
4. Expose Params Backend CPU/CUDA/disk and compute Backend CPU/CUDA; do not expose
   Vulkan in this CUDA/CPU build.
5. Rename the VAE tiling checkbox to state that the engine also tiles automatically
   under pressure; the checkbox means "force tiling".

### 3.2 Component-aware admission

1. Add a header-only safetensors component reader.
   - Reads the 8-byte header length, then the JSON header.
   - Sums tensor byte ranges by role:
     - `model.diffusion_model.*` → diffusion
     - `conditioner.*`, `text_encoders.*`, `te.*` → conditioner
     - `first_stage_model.*`, `vae.*` → VAE
     - anything else → other (charged to the checkpoint module).
   - No tensor data is read.
2. `imageAdmissionRequest` uses the split when available; otherwise keeps the
   whole-file fallback. This makes SDXL/SD1.5-style full checkpoints charge their
   real largest module instead of the whole file.
3. GGUF remains whole-file until a tensor-table parser is verified; this is a
   documented limitation, not a fabricated split.

### 3.3 Workload-aware execution plan

1. Extend `AdmissionEstimate` with `ActivationVRAMMB`.
   - Measured coefficient: ~2.1 MB per megapixel for this SDXL VAE decode family.
   - Batch items are sequential in Howl (`generateSequentialBatch`), so peak uses
     one image; batch count affects time, not peak.
   - Steps affect time, not peak.
2. After weights are admitted, compute remaining headroom:
   `available - required`.
3. If activation exceeds headroom and the user did not force tiling:
   set `loadOpts.VAETiling = true` and emit `image:warning`.
4. The engine remains the fallback authority. If it still retries, parse
   `VAE decode ran out of memory; retrying with spatial tiling` and surface it as
   an observed event, not as a Howl retry.

### 3.5 Automatic quantization

1. Only on admission refusal.
2. Only when Weight Type is empty (Auto).
3. Only for non-GGUF stored checkpoints (safetensors/ckpt); do not re-quantize a
   file whose format is already quantized.
4. Try `q8_0` only. Do not silently drop to q4/q5.
5. If q8_0 admits, use it for this load, emit a visible warning, and do not persist
   it as the user's saved choice.
6. If q8_0 still fails, report the original refusal plus the q8_0 refusal.

### 3.6 Reclamation

1. Wire the residency TTL field only; it is a launch flag, so a resident engine restarts
   when the value changes.
2. A blank field sends `-1` explicitly. Clearing a previously saved timer disables it;
   it must not leave the resident engine running the old timer.
3. Empty Backend and Params Backend are assigned through to the request as Auto. They
   must clear a previous explicit placement rather than silently retaining it.
4. Do not change residency TTL automatically; unloading destroys reusable weights. The
   automatic policy uses a measured process recycle instead, and only when its band or
   the 10% reserve is crossed (§6.11).
5. Do not forward or expose staging TTL. The engine already frees staged copies when a
   component ends, and measured post-generation memory did not change for 0 or 2 s.
6. Document that transient compute buffers are freed by sd.cpp at the end of each
   request; no additional app endpoint exists.
7. One decision point only: `App.reclaimImageEngineAfterGeneration`, called after the
   synchronous generation and from the job manager's all-outcome callback. It reuses
   `internal/resource` estimates, `Image.LoadedOptions`, and the sd-server footprint;
   it does not admit, budget, or choose placement.

---

## 4. Test plan

### 4.1 Unit

- Safetensors component reader: split, unknown tensors, malformed header, no data read.
- `estimateImage`: 10 % reserve; component split; activation scaling; batch peak = 1.
- Policy: auto q8 only on refusal and only Auto type; explicit f16 never changes.
- Policy: auto tiling only when activation exceeds headroom; forced tiling is untouched.
- Args: tri-state Flash Attention, TTL flags, disk Params Backend.
- UI state: read/save/restore of mode and TTL.

### 4.2 Integration / deterministic

- Fake sd-server: sequential generation, batch splitting, and no regression in
  ownership/cancellation tests.
- Arg tests assert the exact spawn command for q4/q8/tiling/TTL/FA.

### 4.3 Real engine, sustained

Use the configured SDXL checkpoint read-only, isolated output/gallery state.

1. Load at 512; generate two sequential images.
2. Generate 1024; verify either pre-enabled tiling or the engine's auto-retry line.
3. Generate a two-image batch; verify sequential execution and both outputs.
4. Constrain with `--max-vram`/offload; verify load and generation still succeed.
5. Trigger a failed retry/recovery path without claiming success; verify next load
   is clean.
6. Confirm no listeners remain.

### 4.4 Hosts

- Browser host (`uiharness`): drive the image page controls and inspect warnings.
- Packaged host: build with `desktop,production`; exercise the same visible state
  where the native surface allows. Record any host-only difference.

---

## 5. Limits that will remain after this slice

- GGUF component splitting is not implemented until its tensor table is parsed and
  verified; admission charges the whole file for GGUF.
- SageAttention, conv-direct, prefetch, batched CFG and module CPU compute are
  deliberately not automated; their tradeoffs are family- or shape-dependent.
- Disk residency remains manual. It disables auto-fit, so Howl's budget cannot place individual
  modules; the engine reloads weights from the model file during sampling. On this box a
  64×64 one-step run took 11 s versus 7–8 s for RAM/VRAM residency; real resolutions are
  much slower.
- The 10 % reserve is app policy; sd-server's `--max-vram -1` retains its own 1 GiB
  floor. The engine remains more conservative at startup.
- Multi-GPU and RPC offload are unverified on this single-GPU box.
- Compute backend choices are limited to the backends this build ships. CUDA and CPU were
  exercised; Vulkan was rejected and its UI option removed. A user-supplied build with a
  Vulkan DLL would need the option restored and measured again.
- Residency TTL demonstrably released RAM (4357 → 361 MiB working set) and the next
  generation reloaded successfully, but the CUDA allocator kept its VRAM reservation in
  `nvidia-smi` (2840 MiB unchanged). TTL is not presented as a VRAM-floor guarantee.
- Staging TTL is not exposed: 0 and 2 s produced no post-generation GPU or RAM change,
  because staged copies are released when a component finishes.
- The 2.1 MB/MP activation coefficient is measured on one SDXL-family checkpoint and
  one 8 GiB card; it is a planning estimate, not a hardware guarantee.
- This sd-server build does not publish per-component residency. The report's actual
  side is the OS working set of the child plus live device use, not the engine's own
  module split. Component-level verification remains impossible until a build prints
  `total params memory size`; the card says so instead of claiming agreement.
- The process working set after readiness was only 171 MB because sd.cpp loads tensors
  lazily; the resident cost becomes visible after the first generation (~7.0 GB here).
  A post-load-only reading is therefore not the model's resident footprint.
- A resolution change on an already-loaded model cannot change `--vae-tiling`; the
  engine's own retry handles the decode. The flag is applied when the process starts,
  which is why the fresh 1024 run had it pre-enabled and the resident 512→1024 run used
  the engine retry instead. This is a real limitation of launch-time flags, not a second
  retry path.
- Automatic VAE tiling is sticky while the process stays loaded. Changing back to a
  resolution that does not need it does not force a restart solely to clear the flag;
  the next real reload re-evaluates it. The user can force tiling, but cannot force it
  off without reloading.
- `--offload-to-cpu` and `--params-backend cpu` are opportunistic: sd.cpp still loads
  modules into VRAM when they fit. Ordinary offload measured 6.37 GiB total peak with
  `--max-vram 3`; CLIP-only CPU 6364, VAE-only CPU 6036, params CPU 6358, VAE tiling
  6030, disable-prefetch 6358, eager-load 6359, disk params 6043. None bounds a small
  GPU diffusion budget. The bounded mode is the explicit per-module CPU assignment: at
  512 no tiling it peaked at 3526 MiB total; with Howl's tiling it peaked at 2089 MiB
  for 512, 2087 MiB for batch and 2083 MiB at 1024. Below the 3072 MiB conservative
  floor the automatic ladder selects **all-CPU** instead of refusing: measured 953 MiB
  total GPU (desktop baseline, no process allocation), 6893 MB process working set,
  18.23 s sampling for 512 one step. This is slower than CPU diffusion + GPU VAE, and
  the VAE stays on CUDA whenever the budget can hold that measured mode. Both modes are
  slower than ordinary GPU diffusion and are named in the plan report.
- Changing a startup option while the previous process is still releasing VRAM can
  make the automatic cap conservative (observed 5 GiB plus pre-enabled tiling at 512
  during one handoff). It is safe and the next generation succeeds, but it can cost
  time. A cold host starts at the full post-reserve cap (observed 6 GiB).
- Automatic and explicit budgets are whole GiB, rounded down from the post-reserve
  budget. A too-large explicit cap is clamped and reported, so it cannot spend the
  reserve. This can leave more free VRAM than the reserve, never less.
- Replacement credit uses the registered model footprint when a process is running,
  and no credit when it is not. It is not a live per-process VRAM measurement; a
  preflight can therefore be conservative relative to what killing the process actually
  releases.
- Batch persistence was fixed in this slice (two images now produce two gallery
  entries); the batch still runs sequentially, so a batch's peak memory is one image
  and its wall time is N images.
- The packaged WebView2 host was driven natively: image page opened, Generate clicked,
  q8_0 + 6 GiB auto cap spawned, image displayed, and a gallery PNG was written. Native
  screenshots were captured under `build/tmp/` for the verification record; they are
  artifacts, not product files.
- No commits, pushes or pull requests were made.

---

## 6. Sustained real-engine verification (final build)

**How it was run.** The browser dev host (`build/bin/howl-ui-harness.exe`) served the
real `App` bindings from an isolated root,
`C:/Users/StoneWolfPC/howl-image-hw-iso-20261124`. Every generation went through
`GenerateImage` to the real `sd-server.exe` (`master-1122-3b12781`). The harness sampled
`nvidia-smi --query-gpu=memory.used` and the `sd-server` working set every 0.5 s while
it generated, read the actual OS command line with `Win32_Process`, and used total GPU
use minus the measured desktop baseline (about 892 MiB) as the process figure. Windows
nvidia-smi does not expose per-process `used_memory` for this WDDM card, so `[N/A]`
there is expected; the total and working set are the measured series. Model:
`Illustrious-XL-v2.0`, q8_0, 512×512 unless stated. The gallery is the isolated
`images` folder.

### 6.1 Five sequential generations in one resident process

| Generation | Wall time | Peak total GPU | Working set after |
|---|---:|---:|---:|
| 1 (cold load) | 13.94 s | 7369 MiB | 407 MiB |
| 2 | 6.99 s | 7369 MiB | 407 MiB |
| 3 | 6.94 s | 6475 MiB | 407 MiB |
| 4 | 6.98 s | 6475 MiB | 408 MiB |
| 5 | 6.94 s | 6475 MiB | 408 MiB |

These are the final-build runs with Max VRAM left at Auto (blank), which now receives
the post-reserve cap: `--max-vram 6 --backend cuda --type q8_0 --weight-residency-ttl -1 --flash-attn auto`.
All five succeeded in one PID (`26020`) with one `Starting sd-server` line. No
`loading tensors completed` block appeared between generations 2–5; the resident weights
were reused. The 6.94–6.99 s steady time is not a reload loop. Peak total 7369 MiB on
an 8192 MiB card leaves 823 MiB free, which is 10.0% of the card, so the established
reserve held at the admission peak.

### 6.2 Batch request

`batch_count=2`, 4 steps: 8.47 s, two images returned, two PNG/JSON pairs written.
Howl logged `Batch: generating 2 images sequentially`, `accumulated image 1/2` and
`2/2`. Peak total GPU was 6470 MiB, the same one-image peak as a single 512 run; batch
count changed wall time (about 2× a one-image 4-step run), not peak. No second process
was started and admission never multiplied a non-existent concurrent peak.

### 6.3 Resolution and admission math

| Request | Activation estimate | Wall time | Peak total GPU | Engine evidence |
|---|---:|---:|---:|---|
| 512×512, 4 steps | 551 MB | 4.29 s | 6470 MiB | no VAE retry |
| 1024×1024, 4 steps | 2202 MB | 26.93 s | 6472 MiB | `VAE decode ran out of memory; retrying with spatial tiling`; tiles reduced 1024→512 |

The admission log now prints the workload component:
`image admission: admitted=true activation=551 MB ...` and `activation=2202 MB ...`.
The weight requirement stayed 5084 MiB because it is the model, while the decode
estimate quadrupled with megapixels and drove the tiling decision/retry. This is the
visible workload-aware part; it is not a fixed constant.

The larger request dropped GPU use to 1642 MiB after the engine's tiled decode, with
the process still alive. The next 512×512 generation reloaded tensors (3.97 s), reused
the same PID and succeeded in 8.48 s. No server restart was needed.

A fresh 2048×2048 run on the same build (2 steps) also succeeded in 125.7 s; the spawn
line included `--vae-tiling` because the 4.2 MP decode estimate exceeded the startup
headroom. That is the pre-emption path, and the engine remains the fallback when the
request reuses an already-loaded 512 process.

### 6.4 Constrained VRAM: measured matrix and bounded mode

> The 2 GiB outcome in this section is from the earlier pass. The final ladder in
> §6.10 replaces that refusal with automatic all-CPU execution; the 3 GiB measured
> bound and the 2636 MB/3072 MB thresholds remain valid.

Direct `sd-server` probes, 512×512, 1 step, q8_0, `--max-vram 3`, 892 MiB desktop
baseline. Process MB is peak total minus that baseline; free is 8192 − peak.

| Flags | Peak total | Process (approx) | Free | Wall |
|---|---:|---:|---:|---:|
| `--offload-to-cpu` | 6364 | 5472 | 1828 | 12.5 s |
| `--backend clip=cpu` + offload | 6364 | 5472 | 1828 | 12.6 s |
| `--backend vae=cpu` + offload | 6036 | 5144 | 2156 | 26.8 s |
| `--params-backend cpu` + offload | 6358 | 5466 | 1834 | 11.7 s |
| `--vae-tiling` + offload | 6030 | 5138 | 2162 | 12.7 s |
| `--disable-prefetch` + offload | 6358 | 5466 | 1834 | 11.8 s |
| `--eager-load` + offload | 6359 | 5467 | 1833 | 9.4 s |
| `--params-backend diffusion=disk,clip=cpu` + offload | 6043 | 5151 | 2149 | 11.7 s |
| `--backend diffusion=cpu,clip=cpu,vae=cuda0` (no offload) | 3526 | 2636 | 4666 | 22.4 s |
| same + `--vae-tiling`, 1024×1024 | 2075 | 1185 | 6117 | 79.4 s |
| `--backend diffusion=cpu,clip=cpu,vae=cpu` (all CPU) | 857 | 0 | 7335 | 37.4 s |

No GPU-diffusion combination bounded a 3 GiB budget; offload still loads fitting
modules into VRAM, and CLIP/VAE/params/tiling/prefetch variations do not change that.
The explicit per-module CPU assignment is the only bounded GPU-assisted execution
class. All-CPU is also bounded (no process VRAM), but 37 s versus 23 s for CPU
diffusion + GPU VAE at 512, and much worse at 1024. The final policy therefore keeps
the GPU-assisted bounded mode at ≥3 GiB and selects all-CPU below it, rather than
refusing a budget that CPU-only can honor. This is an evidenced tradeoff, not a hidden
inability; §6.10 records the CPU-only run and the recovery back to 3 GiB.

Howl selected that mode automatically in this pass when the explicit budget was below
the resident weights plus decode need. Actual command:

```
sd-server.exe --listen-port 8092 --threads 16 -m "…/Illustrious-XL-v2.0.safetensors"
  --max-vram 3 --backend diffusion=cpu,clip=cpu,vae=cuda0 --type q8_0
  --weight-residency-ttl -1 --vae-tiling --flash-attn auto --offload-to-cpu
```

Howl's samples over the whole run, same process (`PID 25268`):

| Request | Wall | Peak total | Free at peak | Result |
|---|---:|---:|---:|---|
| 2 GiB request | 0.04 s | — | — | refused before spawn; no PID |
| 512×512, 1 step, next request | 23.51 s | 2089 | 6103 | image, PID 25268 |
| 512×512 repeat | 15.87 s | 2091 | 6101 | image, same PID |
| 512×512 batch 2 | 31.99 s | 2087 | 6105 | two images, same PID |
| 1024×1024 | 71.81 s | 2083 | 6109 | image, same PID |

Fresh-host 2 GiB check: refused in 0.04 s with **no `sd-server` process or listener**
(`sd_after_refusal: []`), then the 3 GiB request started one process, peaked at
2089 MiB and succeeded in 23.51 s.

Resident recovery on the final binary: PID 17524 generated at 3 GiB (23.61 s, peak
2089 MiB); a 2 GiB request was then refused in 0.03 s while that PID stayed alive and
unchanged; the next 3 GiB request succeeded in 15.84 s on the same PID, again peaking
at 2089 MiB.

**Verdict for the done criterion:** at 3 GiB the process is 1197 MiB (total 2089),
below the configured 3072 MiB, and the device retains 6103 MiB (74%), far above the
required 819 MiB (10%). Batch, 1024 and repeat use the same command and stay at or below
2091 MiB. In this earlier pass the 2 GiB request was refused before any process started;
the error named the 2636 MB process need and 3072 MB conservative floor, the 5084 MB
resident weights, the 551 MB decode estimate and the 2 GiB budget. The final ladder in
§6.10 changes that below-floor outcome to automatic all-CPU execution, while keeping
the 3 GiB bounded mode and its measured 2636/3072 MB thresholds.

**Tradeoff/limitation:** the bounded mode runs diffusion on CPU. It is 16–24 s for one
512 step versus about 7 s for GPU diffusion, and 73 s for 1024; that is the cost of the
hard bound. Ordinary GPU diffusion + offload is still not a hard total, but Howl no
longer uses it for a budget below the resident need. The 2636 MB need and conservative
3 GiB floor come from this q8 model and this card; a different model/card needs its own
measurement rather than an inherited guarantee.

### 6.5 Reclamation and retention

- Default `--weight-residency-ttl -1`: the five-generation run held the same PID,
  weights and 408–410 MiB working set; nothing was unloaded between generations.
- `--params-backend cpu --max-vram 3 --weight-residency-ttl 0`: after generation the
  working set fell from 4360 MiB to **361 MiB**; idle samples 4 s later were still 361
  MiB and 2620 MiB total GPU. The same process then generated again in 11.89 s,
  reloading tensors without a restart.
- `--weight-residency-ttl 2` was verified separately in §1.3: 4357 → 361 MiB after
  idle, next generation reloaded and succeeded.

Model components/RAM are released; the CUDA context and its cached allocator remain in
`nvidia-smi` until the process exits. Reclamation is therefore a working-set/weights
release, not a return to the 892 MiB desktop baseline. That is documented in §5.

### 6.6 Stress and recovery

A deliberate oversized request, 4096×4096 at `--max-vram 3` (q8 weights need 5084 MiB),
was refused before any engine start:

```
image admission: admitted=false activation=35232 MB
cannot start image "Illustrious-XL-v2.0": needs 5084 MB VRAM;
8688 MB free on device 0, 868 MB reserve (10%), 3072 MB available
```

It returned in 0.12 s with a clear error, no hang and no process change. The next
512×512 request succeeded in 2.49 s on the same PID. This is the required clean refusal.

Before the fixes, the same 4096 request was admitted at 2048 MB because a previous
`--params-backend cpu` was silently retained; the engine then returned
`diffusion output shape [] does not match input [512, 512, 4, 2]` and Howl marked the
engine unavailable. The rerun above is the repaired behavior.

### 6.7 Defects found and fixed by this run

1. **Stale auto-cap flag.** `imageLoadOptionsForGeneration` overwrote an explicit
   `MaxVRAM` but left `AutoMaxVRAM=true` from the previous automatic startup. The next
   blank Max VRAM copied the old 3 GiB cap forward; a 6 GiB request was refused with
   `3072 MB available`. Fixed by clearing `AutoMaxVRAM` whenever the generation panel
   supplies the cap. Verified: params CPU/max3 → blank/Auto restarted a new PID with
   `--max-vram 6` and no `--params-backend`, generated in 10.0 s.
2. **Blank Backend/Params did not restart a resident engine.** Empty placement compared
   equal to “no change,” so the UI could show Auto while the old process kept running
   `--params-backend cpu`. `imageStartupOptionsDiffer` now restarts when either side is
   empty and the other is an explicit flag. Verified with the same new-PID run.
3. **Max VRAM changes were not startup-differed.** A changed cap could sit on the
   resident process. `imageStartupOptionsDiffer` now compares `MaxVRAM`.
4. **Admission telemetry did not show the workload component.** The log now prints the
   activation estimate, which is what made the 512/1024 difference auditable.
5. **Truthful UI text.** The VRAM Budget and Offload controls now state that
   `--max-vram` is a per-device managed budget and that offload peak can still load
   modules into VRAM; neither claims a hard total.
6. **Blank Max VRAM used the engine's 1 GiB reserve.** `nil` reached sd-server as `--max-vram -1`, which reserves 1 GiB rather than 10% of the card. Blank Auto now receives the same post-reserve automatic cap as Max VRAM 0; a cold blank run starts at `--max-vram 6` on this card. The placeholder now says `auto (10% reserve)`.
7. **Replacement credit could exceed the card.** During a blank-Max-VRAM restart the
   old image process was still alive; free memory plus the registered replacement
   footprint reported 8693 MB on an 8192 MB card and an automatic 7 GiB cap. The
   estimate is now clamped to the device total. The observed peak still left 823 MiB
   free (10.0%), but the counter itself was wrong before the clamp.
8. **A small explicit offload budget admitted over-budget GPU diffusion.** At
   `--max-vram 3 --offload-to-cpu` the engine peaked at 6364 MiB total (process 5472),
   silently exceeding the configured cap. Direct probes showed no GPU-diffusion flag
   combination bounds it. Howl now enters the explicit CPU-diffusion + VAE-tiling mode
   for such budgets and refuses below the conservative 3072 MB floor after measuring
   a 2636 MB process need, before leasing VRAM or spawning. The 3 GiB path peaks at
   2089 MiB (512), 2087 MiB (batch 2) and 2083 MiB (1024); 2 GiB is refused with no
   process, and the next request succeeds.

`go test ./...` passed after these changes, including the new focused admission test
(`TestImageConstrainedOffloadBoundsOrRefusesAtAdmission`). Both hosts rebuilt:
`build/bin/howl-ui-harness.exe` (uiharness) and `build/bin/howl-chat-s1.exe`
(desktop,production). The rebuilt harness then re-ran the constrained matrix in §6.4.
No commit, push or PR was made; the harness and engine were shut down and no listener
remained.

### 6.8 UI-driven budget/offload control (post-run verification)

This pass did not call the binding directly for the budget decision. It drove the app's
own `#ig-max-vram` control in the browser dev host (`http://127.0.0.1:8799`), which was
serving the real `App` bindings from the same isolated root
(`C:/Users/StoneWolfPC/howl-image-hw-iso-20261124`) against the real `sd-server.exe`.
Model `Illustrious-XL-v2.0`, 512×512, 1 step, Discrete scheduler, Offload to CPU checked.

**Source of truth.** `frontend/dist` is hand-maintained, not generated. There is no
`package.json`, frontend source tree, or frontend build config; the only image page is
`frontend/dist/image-generator.html` plus its JS/CSS. Production embeds that directory
with `//go:embed all:frontend/dist` (`main.go` line 18), and the harness serves the same
directory with `http.FileServer(http.Dir("frontend/dist"))` (`harness_main.go` line 120).
After the final rebuild, the production binary contained both exact hint strings
(`grep -a -F`): `conservative 3 GiB floor (measured need 2.6 GiB)` and
`a constrained VRAM budget runs diffusion on CPU`. A fresh build therefore cannot
revert this wording; it is the shipped source.

**Earlier 2 GiB policy (superseded by §6.10).** Setting the budget field to `2` and
clicking Generate on that build produced a persistent visible status (not only an
alert), no image, and no new engine process. The status element wrapped its full text (550.5 px wide, 66.64 px high,
`scrollHeight` 65 = `clientHeight` 65); the exact text was:

```
Error: Could not verify available VRAM: VRAM Budget 2 GiB is below the measured 2636 MB CPU-diffusion process need (conservative 3072 MB admission floor; resident weights 6945 MB plus 551 MB decode); refusing rather than starting a process over budget
```

The final ladder no longer refuses this case: it starts `--backend cpu`, as recorded in
§6.10. The refusal text above is historical evidence for the earlier policy, not current
behavior.

This was on the resident PID `19804`; the harness log still held exactly one
`Starting sd-server` line. The earlier refusal also happened before any spawn on a
fresh host, so the no-spawn result is not an artefact of the resident process.

**3 GiB through the UI control.** With the budget field at `3` and Offload to CPU
checked, the visible status became `Done!` and a 512×512 image appeared in
`#ig-preview-area`. The actual spawned command line, read from the harness log, was:

```
sd-server.exe --listen-port 8092 --threads 16 -m "G:\AI GGUF\Image models\Illustrious-XL-v2.0.safetensors"
  --max-vram 3 --backend diffusion=cpu,clip=cpu,vae=cuda0
  --weight-residency-ttl -1 --vae-tiling --flash-attn auto --offload-to-cpu
```

That is one spawn, PID `19804` (started 07:54:11). A separate `nvidia-smi` series for
the whole UI-driven run (`build/verification/gpu-ui-3g-sustained.csv`, 32 samples at
1 s from 07:57:37.968 to 07:58:09.343) measured a peak of **2188 MiB total** and
**6004 MiB free (73.29%)**. The configured budget was 3072 MiB and the required
reserve was 819 MiB, so the run stayed inside the budget and above 10% free. The
adjacent earlier UI run sampled 2125 MiB peak with 6067 MiB free, consistent with the
direct 2089–2091 MiB series.

**Resident recovery.** After the 3 GiB success, the same UI control was set to `2` and
Generate was clicked: the same visible refusal appeared while PID `19804` stayed alive
and unchanged. Setting it back to `3` and clicking Generate completed with `Done!` on
that same PID. The spawn count remained one throughout; the resident process survived
the refusal and served the next request.

**Hint/command consistency.** The served performance panel said that a budget below the
resident need with Offload enabled “switches to the bounded CPU-diffusion mode”, that
budgets below the “conservative 3 GiB floor (measured need 2.6 GiB)” are refused before
start, and that Offload “Stages weights in RAM and reloads them on demand” and runs
diffusion on CPU below the resident-plus-decode need. Those statements match the
actual command line above (`diffusion=cpu,clip=cpu`, `--max-vram 3`, `--vae-tiling`,
`--offload-to-cpu`). The generation status while offload was active also said
`CPU offload is active; a constrained VRAM budget runs diffusion on CPU (slower),
otherwise weights are staged from RAM.` No hint/command mismatch was found in this run.

**Residual limitation found by this pass (not repaired here).** A UI request at 3 GiB
with scheduler `Karras` and `steps=1` spawned the same constrained engine but returned
`sd-server txt2img error 500: {"error":"line 65, column 13: line 65, column 88: number
must be finite"}` and Howl marked the engine unavailable. Changing only the scheduler to
`Discrete` at the same 1 step succeeded on a fresh engine. This is an engine/scheduler
edge case, not a budget or offload failure: the budget control did not cause it, and the
budget/offload flags were received correctly. The UI still offers Karras with 1 step, so
the combination remains a documented limitation; this pass did not change scheduler
semantics. The sampler is 1 s for this UI series (0.5 s elsewhere), so a shorter spike
than that could be missed; nvidia-smi reports total device use, including the desktop.

### 6.9 Verification limits in this run

- Windows nvidia-smi reports `[N/A]` for per-process compute memory on this WDDM GPU.
  Process VRAM is the measured total minus the 892 MiB desktop baseline; a 0.5 s
  sampler can miss a shorter spike.
- The packaged WebView2 host was re-driven natively in §6.10 on the final binary: the
  image page opened, Generate was clicked, the real engine spawned with the constrained
  command, an image appeared, the placement card showed the measured footprint, and the
  gallery PNG was written.
- The ordinary 6.37 GiB offload peak and the 3 GiB CPU-diffusion floor are one
  model/card measurement (q8 Illustrious-XL, 8 GiB GTX 1080). They are not generalized
  to larger cards or other models without re-measurement.
- The bounded mode deliberately runs diffusion on CPU. It is much slower than GPU
  diffusion; that tradeoff is shown in §6.4, not hidden.
- Batch is sequential by Howl design, so there is no concurrent multi-image memory
  peak to measure or admit.

---

## 6.10 Final automatic-placement and footprint pass (2026-11-24)

**Build.** `build/bin/howl-ui-harness.exe` and `build/bin/howl-chat-s1.exe`, both built
from the final tree after the automatic-placement reset and footprint code. The browser
host used isolated root `C:/Users/StoneWolfPC/howl-image-hw-iso-20261124`; the packaged
host used the same root. Model `Illustrious-XL-v2.0`, 512×512 otherwise, 1 step, Euler A,
Discrete, stored f16 (no `--type`), Offload unchecked, Backend/Params Auto.

**The final automatic command at 3 GiB.** A request that used to require the Offload
checkbox now produces, without it:

```
sd-server.exe --listen-port 8092 --threads 16 -m "…/Illustrious-XL-v2.0.safetensors"
  --max-vram 3 --backend diffusion=cpu,clip=cpu,vae=cuda0
  --weight-residency-ttl -1 --vae-tiling --flash-attn auto --offload-to-cpu
```

**The final automatic command at 2 GiB.** The same controls at `2` produce no GPU
weights or compute, and no refusal:

```
sd-server.exe --listen-port 8092 --threads 16 -m "…/Illustrious-XL-v2.0.safetensors"
  --max-vram 2 --backend cpu --weight-residency-ttl -1 --vae-tiling --flash-attn auto
```

Admission logged `0 MB VRAM required, 2048 MB available after 0 MB reserve`; the first
run was PID 9652, sampling 18.23 s, working set 6893 MB, total GPU 953 MiB (desktop
baseline), 7107 MiB free. The reserve was held by a wide margin. On the final binary a
fresh 2 GiB UI request emitted exactly one visible warning — `VRAM Budget 2 GiB selected
automatic CPU-only placement: no GPU weights or compute. This honors the budget but is
much slower.` — then sampled in 18.00 s, finished with 6891 MB working set and 951 MiB
total GPU (7109 MiB free), displayed the image, and updated the card. The preflight did
not duplicate the warning.

**Sustained partial-placement run.** One resident process, PID 17796, five sequential
512×512 images. No second `Starting sd-server` line appeared. Working set is the OS
measurement after each generation, not a process-VRAM estimate.

| Generation | Sampling | Working set after | Total GPU | Free |
|---|---:|---:|---:|---:|
| 1 (cold, lazy load) | 18.52 s | 7024 MB | 1666 MiB | 6394 MiB |
| 2 | 15.89 s | 7367 MB | 1663 MiB | 6397 MiB |
| 3 | 16.05 s | 7031 MB | 1662 MiB | 6397 MiB |
| 4 | 16.11 s | 7035 MB | 1662 MiB | 6397 MiB |
| 5 | 16.14 s | 7037 MB | 1662 MiB | 6398 MiB |

After the first generation the weights stayed resident: steady sampling fell to about
16 s and the working set did not grow across four more requests (7031–7037 MB). That is
the between-generation reclamation/retention result: transient buffers were released,
reusable weights were not reloaded. No `loading tensors completed` block appeared
between generations 2–5.

**Batch and larger resolutions.** Batch 2 at 512 completed two images in the same PID
(sampling 16.35 s and 16.49 s), working set 7045 MB, total GPU 2147 MiB, free 5913 MiB.
1024×1024 took 66.36 s sampling; the observed peak working set was 7543 MB, total GPU
2145 MiB, free 5915 MiB, and the post-generation working set was 7062 MB. 1344×1024
took 93.16 s on the final binary; a 32-sample trace (`build/verification/gpu-1344-peak2.csv`)
peaked at **8334 MB working set, 2297 MiB total GPU, 5763 MiB free**, then returned to
7030 MB working set. No OOM, no VAE retry, and the PID stayed 18424.

**Recovery in one session.** Starting from the 2 GiB CPU-only PID, setting the budget
back to 3 and clicking Generate restarted to PID 18424 with the partial command above;
sampling took 17.87 s, working set 7023 MB, total GPU 1697 MiB, free 6363 MiB. Changing
back to 2 restarted to PID 25400 with `--backend cpu` and generated successfully. A
budget change is re-planned, not carried as saved intent.

**Visible card.** Browser host, after the final generation:

```
Placement planned: automatic bounded placement: diffusion and CLIP on CPU, VAE on CUDA and tiled.
This sd-server build did not report per-component residency, so the plan is unverified.
Measured sd-server process working set: 171 MB after load, 7023 MB after generation.
```

The details panel separately shows theoretical weights (6.5 GB), decode activation
(551 MB), resident need (6.8 GB), the plan flags, the live device reserve, and the
measurement source. The packaged WebView2 host showed the same card and status text after
its own Generate click; its actual command line was the 3 GiB partial command, one
`image hardware` warning was logged, sampling took 18.68 s, and the new image was
written at `…/images/img_20261124_154428_33300.png`. The 171 MB post-readiness value is real and
useful: this engine loads tensors lazily, so the post-generation value is the resident
footprint. The card does not claim the engine confirmed component placement, because it
did not.

**What this pass fixed.** The final 2 GiB request first inherited the saved 3 GiB
automatic placement, skipped re-planning, and was refused by admission. Automatic
placement is now re-derived whenever the model or budget differs; the 1 → 2 → 3 GiB
sequence above is the verification. A preflight warning is no longer duplicated by the
start path.

**Verification status.** `go test ./...` and `node --check
frontend/dist/js/image-generator-core.js` pass on the final tree; both hosts build. The
browser host was driven through its real controls; the packaged host was driven with
native mouse clicks. Remaining unverified: component-level residency (engine does not
publish it), q8 CPU-only RAM reduction (not selected because it is a quality tradeoff,
not a fit requirement), and other model families/cards (thresholds are measured only for
this SDXL checkpoint on this 8 GiB card). No commit or push was made.

---

## 6.11 Measurement-gated between-generation reclamation (2026-11-24)

**One owner.** `App.reclaimImageEngineAfterGeneration` is the only reclamation decision.
It is called from the synchronous `GenerateImage` path and from
`job.Manager.SetOnFinished`, which now runs for every generation attempt before the job
reaches its terminal status. It does not admit, budget, or place a model: it reads the
already-admitted `Image.LoadedOptions` and the existing `internal/resource` estimate,
and it leaves sd.cpp's placement/VAE-retry decisions untouched. No second budget or
admission path was added.

**Baseline.** The readiness footprint is recorded but is only a lazy-loader floor on this
engine (171 MB). The first successful generation becomes the resident baseline
(7022–7024 MB here). Without that calibration every first generation would look like a
multi-gigabyte leak. The report names the baseline source, so this is visible rather than
implicit.

**Band and trigger.** After every generation, including failures, the measured OS
working set is compared with that baseline. The workload-aware band is:

```
max(512 MB,
    0.50 * decode activation
  + 0.02 * component bytes + 32 MB per active component
  + 128 MB per extra batch item (capped at 512 MB)
  + placement allowance: 256 MB GPU-resident, 512 MB partial/offload, 768 MB CPU-only)
```

For the measured SDXL run at 512×512 this is 1016 MB for one GPU/partial image and
1144 MB for batch 2. A reclaim is requested when growth exceeds the band, when free
VRAM is at or below the established 819 MB (10%) reserve, or when a failed generation
left the engine state abandoned (recovery, not a memory claim). The report records the
precise reason and measurement.

**Lever.** This pinned sd-server has no unload/free route. The strongest real lever is a
controlled process recycle: `Engine.LoadModel` kills the child and reloads the same
admitted options. This is not a new placement or budget decision; `--weight-residency-ttl`
is not changed. The next generation pays the lazy tensor load again, which is why this
path is gated rather than run every time.

**Observability.** The placement report now carries `reclaim`: policy, baseline and its
source, band, threshold, measured/growth, reclaim count, last decision, every reclaim
event, and the result (PID before/after, working set before/after, free VRAM before/after,
reserve, reload error). Logs emit one `[ImageReclaim]` decision line and a second line
for an actual recycle. The placement card adds the same policy/baseline/threshold/event
line; no user control was added, because there is no separate engine setting to expose.
The only test hook is `HOWL_IMAGE_RECLAIM_TEST_GROWTH_MB`, consumed once and never shown
as a product setting. The reserve branch did not use a test hook: it was forced by
`build/verification/hold_vram.exe`, a separate driver-API allocation on the same card,
while the production policy read the real `RefreshSnapshot` value.

**Evidence (real sd-server, isolated root
`C:/Users/StoneWolfPC/howl-image-hw-iso-20261124`).**

| Case | Measured result |
|---|---|
| Five sequential 512×512 images, one resident PID 16776 | Working set 7024, 7026, 7029, 7034, 7036 MB; zero reclamations; same PID; no second `Starting sd-server`. Baseline 7024 MB, band 1016 MB, threshold 8039 MB. Log: `build/verification/reclaim-sustained.log`. |
| Batch 2 at 512×512, same PID 16776 | Band 1144 MB, threshold 8167 MB, post-generation 7040 MB, growth 17 MB; zero reclamations. Log: same. |
| 2 GiB CPU-only, PID 12092 | `--max-vram 2 --backend cpu`; baseline 6887 MB, band 1272 MB, subsequent 6894 MB (growth 7 MB), free 7249→7241 MB, reserve 819 MB held; zero reclamations. Log: `build/verification/reclaim-forced.log`. |
| Forced growth, final binary, PID 18424 | Baseline 7022 MB, band 1016 MB, measured 7026 MB, growth 4 MB + injected 4096 MB = 4100 MB; exactly one reclaim to PID 2264, working set 7026→171 MB, free 6425→7216 MB; next generation succeeded on PID 2264 at 7022 MB with no second reclaim. Log: `build/verification/reclaim-final.log`. |
| Failure then success | Karras + 1 step returned `sd-server txt2img error 500 ... number must be finite` and left the engine abandoned. The policy recycled PID 20536→5008 (working set 7023→171 MB, free 6399→7142 MB); the next Discrete 1-step generation succeeded on PID 5008 at 7022 MB. Log: `build/verification/reclaim-forced.log`. |
| **Reserve branch, real external pressure, final binary, PID 7124** | CPU-only 2 GiB engine loaded at baseline 6890 MB, band 1272 MB. A separate CUDA helper allocated 6500 MB on the same card; the helper's driver reading was 685 MB free and `nvidia-smi` reported 343 MiB free. The post-generation policy snapshot read **441 MB free**, below the 819 MB reserve, while measured working set was 6894 MB and growth was only 4 MB (inside the band). It reclaimed for the reserve reason, PID 7124→22832, working set 6894→171 MB, post-reclaim free 7055 MB. The next generation succeeded on PID 22832 at 6891 MB with 7077 MB free and zero further reclaim. Helper released when the old PID exited. Logs: `build/verification/reclaim-reserve-final.log`, `build/verification/hold-vram-reserve-final.log`. |
| Retention after reserve reclaim, PID 22832 | CPU-only generation after the recycle: baseline 6891 MB, band 1272 MB, post-generation 6891 MB, free 7077 MB, zero reclaim. |
| Regression: sequential and batch, partial placement, PID 21968 | Two sequential 512×512 images: 7023→7028 MB, growth 5 MB, zero reclaim. Batch 2: band 1144 MB, post-generation 7031 MB, growth 8 MB, zero reclaim. Same PID throughout. Log: `build/verification/reclaim-reserve.log`. |
| Regression: forced growth after the reserve run, PID 25436→26172 | Baseline 7023 MB, band 1016 MB, measured 7026 MB, growth 3 MB + injected 4096 MB = 4099 MB; exactly one reclaim, working set 7026→171 MB, free 6247→6967 MB; next generation succeeded on PID 26172 at 7023 MB. Log: `build/verification/reclaim-growth-final.log`. |
| Packaged desktop host, 3 GiB, one 512×512 generation | Real Generate click (native input) spawned the partial command, image appeared, card was visible, sampling 16.43 s, baseline 7024 MB, band 1016 MB, free 5973 MB; zero reclaim. Screenshot: `build/verification/desktop-reclaim-small.png`; log: `build/verification/reclaim-desktop.log`. |
| Packaged desktop host, 2 GiB CPU-only, final binary, PID 17480 | Real Generate click (native input) spawned `--max-vram 2 --backend cpu`; image appeared, card was visible, sampling 16.27 s, baseline 6891 MB, band 1272 MB, free 6846 MB; zero reclaim. Screenshot: `build/verification/desktop-reclaim-final-small.png`; log: `build/verification/reclaim-desktop-final.log`. |

The final card/report after the reserve run read: `Baseline: 6.8 GB (first successful
generation (resident weights)); band 1272 MB; threshold 8163 MB; reclaim events 1` and
retained the event with reason `free VRAM 441MB fell to or below the 819MB (10%) reserve`,
working set `6894 -> 171 MB`, and free `441 -> 7055 MB`. The growth run separately
reported `working-set growth 4099MB exceeded the workload band 1016MB`.

**Residual limitations with evidence.**

- The first successful generation calibrates the baseline, so a leak that happens during
  that first generation is not distinguishable from the known lazy resident load. This
  is why readiness (171 MB) is not used as the comparison baseline.
- sd-server `3b12781` publishes no per-component residency; the measurement is the OS
  working set, not a component map. The card says so.
- A recycle is a process restart, not an in-process unload. It returns the process to
  171 MB and the next generation reloads weights; no measured in-process free surface
  exists in this build. The cost is paid only on a band/reserve crossing or abandoned
  engine recovery.
- The band constants are derived from this SDXL checkpoint and 8 GiB GTX 1080. Another
  model/card needs its own measurement; they are planning allowances, not a capacity
  guarantee.
- The free-VRAM snapshot is taken after completion, so a spike shorter than the sampler
  interval can be missed. The reserve branch itself was forced with a real 6500 MB
  external CUDA allocation, not a synthetic reading: the policy observed 441 MB free
  against the 819 MB reserve. Post-reclaim it observed 7055 MB free after the helper
  released when the old sd-server PID exited. This is a point-in-time reading, not a
  claim about sub-second transients.
- An abandoned engine is recycled even when the working set is inside the band. That is
  failure recovery, and the report states the failure reason rather than calling it
  memory growth.
- The reserve case was forced on the CPU-only engine so the 6500 MB external allocation
  could not make the generation itself OOM. That isolates the post-generation reserve
  decision from growth; it does not claim that a GPU/partial generation can run under
  only 343 MB free. The shared production decision and snapshot are the same for both
  placement modes. The helper's `cuMemGetInfo` read 685 MB, `nvidia-smi` read 343 MB,
  and the policy snapshot read 441 MB at decision time; all are below the 819 MB
  reserve, with the difference attributable to driver/display accounting timing.
- Reclaim history in the report is capped at the last 10 events; the log is the full
  session record.

**Verification.** `go test ./...`, `go build -tags uiharness`, `go build -tags
"desktop,production"`, and `node --check frontend/dist/js/image-generator-core.js` pass
on the final tree. The browser host was driven through its real image page and bindings;
the packaged desktop host was also driven with native input: its image page opened, a
real Generate click produced a visible image and the placement/reclamation card
(`build/verification/desktop-reclaim-small.png`), and the log recorded the same
retention decision (`build/verification/reclaim-desktop.log`). The browser card and the
JSON report were read back, not just the log. The reserve run was separately driven with
the external CUDA pressure helper and the same production path; it attributed the
reclaim to `free VRAM 441MB fell to or below the 819MB (10%) reserve`, not to growth.
No commit or push was made.
