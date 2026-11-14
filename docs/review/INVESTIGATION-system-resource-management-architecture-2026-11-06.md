# INVESTIGATION — System Resource Management and VRAM/RAM Architecture

**Date:** 2026-11-06  
**Repository:** `C:\Users\StoneWolfPC\HOWL_Chat`  
**Branch / HEAD:** `master` / `9c7d2a1`  
**Mode:** READ-ONLY source investigation. No application launch, build, test, benchmark, model load, configuration write, download, or production repair was performed.  
**Scope:** Chat, Image, Audio, and Video where shared infrastructure or existing code touches it.  
**Deliverable status:** Investigation only. This document does **not** authorize implementation.

---

## Evidence vocabulary

This report separates what was read from what would require execution:

| Grade | Meaning in this report |
|---|---|
| **Source-observed** | Directly visible in the cited source path and symbol. This is evidence of code structure, not proof of runtime behavior. |
| **Source-derived** | A conclusion that follows from combining several source-observed facts. Runtime behavior is still unproven. |
| **Inferred** | A plausible mechanism or consequence that the source suggests but does not establish. |
| **Unestablished** | Requires a runtime measurement, external artifact, or product decision not available in this read-only audit. |
| **Not applicable** | The item does not exist in the current code path, or the question does not apply to the modality as implemented. |

“Observed” below never means “measured on this machine.” No runtime measurements were taken. Historical numbers supplied by the user are treated as observations with no located original artifacts.

---

## 1. Executive findings

1. **The claimed “98% of free VRAM” ceiling is not present in the allocation path.** The chat load path uses **85% of free VRAM** and **90% of free RAM** (`app-model.go:250-251`, `app-model.go:734-735`). The separate allocator uses **90% of free VRAM** and **95% of free RAM** (`internal/resource/allocator.go:21-22`). Spawn validation permits an estimate up to **110%** of reported free memory (`internal/backend/gguf/memory.go`, `ValidateSpawn`). No `0.98` allocation constant or 98% policy exists in the inspected Go source. The only `0.98` match in the relevant search was an unrelated Whisper parser value. The historical “98%” should therefore be treated as an observation or an external/backend measurement, not as Howl’s implemented policy.

2. **Howl has a resource observer and advisor, not a single admission authority.** `internal/resource.ResourceManager` stores a hardware snapshot and a map of models, and `internal/resource.Allocate` can recommend eviction or CPU fallback. But there is no background monitor, no reservation, no queue, and no central lock around model starts. The production load paths call the manager after or around their own decisions; they do not ask it for permission in a way that can atomically reserve memory.

3. **Chat, Image, Audio, and Video launch independent backends with different ownership rules.** Chat owns a llama-server on fixed port 8080 (`app-llama.go`, `buildLlamaServerArgs`). Image owns sd-server on port 8081 and may launch sd-cli for video (`internal/image/engine/sdcpp/engine.go`, `video.go`). Audio owns a separate llama-server on an ephemeral port (`internal/audio/asr/llama/llama.go`). There is no cross-backend startup lock.

4. **Chat’s automatic optimization is real in code but conservative, silent, and pre-launch only.** `App.applySafeDefaults` computes a memory budget, calls `MemoryModel.ClampToHardware`, and writes clamped `context_size`, `gpu_layers`, and `batch_size` back into `modelSettings` (`app-model.go:726-810`). The adjustment strings are printed to the console and then discarded; they are not emitted to the UI. The in-memory values are not saved by `applySafeDefaults`, so the original persisted values return on restart.

5. **The runtime backend may change the placement after Howl’s estimate.** For non-SWA models the chat path leaves llama.cpp parameter fitting enabled (`fitParams := true`; `--fit off` only for sliding-window models). “Fit” is a backend behavior, not Howl’s estimator. Howl registers the pre-launch estimate and does not read actual per-layer placement for chat. `HTTPClient.GetUsedMemory` and `GetTotalMemory` are stubs returning zero (`internal/backend/llama/http_client.go`). Image is the only modality with a measured placement parser, and that data is a binding (`GetImagePlacement`) with no non-generated frontend caller found.

6. **Resource registration is not a reliable lifecycle record.** Chat registers the model immediately after spawning, before the asynchronous health monitor has established readiness (`app-model.go:338`, `:489-512`). A later load failure does not unregister it. Switching chat models kills the old process but does not unregister the old model; a different model name leaves both entries. Image’s normal generation path (`StartImageJob` → `ensureEngineReady` → `engine.EnsureReady`) does not register at all; only the manual `SetActiveImageModel` path registers, and it records fixed `4000 MB` GPU / `2000 MB` RAM with no measurement (`app-image-model.go:143-150`). Image unload and reinitialization do not unregister (`app-image.go:176-181`; `reinitEngine`).

7. **The frontend preflight is advisory and can be bypassed by its own control flow.** Chat calls `CheckModelFit` before a UI-initiated load; image calls it before generation. But both frontends immediately proceed when `can_fit` is false **and** `blocking.length === 0` (`frontend/dist/js/chat-model.js:332-347`, `frontend/dist/js/image-generator-core.js:592-607`). The backend `LoadModel` and `StartImageJob` paths do not repeat a mandatory resource admission check. `RequestAllocation` exists as a Wails binding but no production caller was found outside generated bindings.

8. **Image ignores the modeled VRAM cap.** `ModelLoadOptions.MaxVRAM` is parsed (`app-image-model.go:211-212`) and the type documents it (`internal/image/types/types.go:141-142`), but `buildServerArgs` always emits `--max-vram -1` (`internal/image/engine/sdcpp/engine.go:498-502`). The comment in the engine says a negative value reserves that much free VRAM (i.e. keeps 1 GiB back); the type comment says `-1` means “use all available.” That conflict is itself an implementation/documentation defect.

9. **Audio is outside the shared resource architecture.** The local Qwen3-ASR llama-server is started with `-ngl 99`, context 2048, and no `--fit` (`internal/audio/asr/llama/llama.go:93`). It does not register with `ResourceManager`, does not consult free VRAM, and is not considered by chat/image arbitration. `BaseRecognizer.CheckHardwareCompatibility` checks total GPU memory and available RAM, not free VRAM (`internal/audio/asr/asr.go:190-220`), and the local Qwen llama adapter does not call it. The audio manager has strong generation/lease/cleanup semantics, but those semantics are internal to audio.

10. **“Context rollover” is presented in the UI but not implemented in the chat path.** The UI persists `context_rollover_mode`, `memory_integration`, `system_prompt_cache`, and `tokenizer_override` (`frontend/dist/js/model-settings.js`, `frontend/dist/index.html`). No backend reference to those four keys was found. The actual prompt path calls `truncateMessagesToFitContext`, which estimates `contextSize * 4` characters and keeps the system prompt plus the newest messages (`app-prompt.go:572-646`, called at `:771`). It has no tokenizer, no memory-guided policy, and no rollover event.

11. **Howl can detect readiness, not acceptable performance.** Chat’s monitor stops when llama-server answers `/health` after the spawn; it does not measure tokens/second, CPU/GPU utilization, or actual offload. Image has a stdout placement parser but no performance gate. A CPU-only or RAM-spilled model can pass health and be “loaded” while being unusably slow. This directly matches the previously observed chat GPU-offload degradation: the architecture has no acceptance test for runtime performance.

12. **No host-wide OOM has been established.** The risks below are architectural; actual OOM, fragmentation, and driver outcomes require runtime evidence. The most dangerous controllable gap is the combination of independent backends, advisory preflight, and stale/absent registration, not a single proven crash.

---

## 2. Architecture and ownership map

### 2.1 Shared infrastructure

| Concern | Owner | Exact source | What it actually does |
|---|---|---|---|
| Hardware snapshot | `ResourceManager` singleton | `internal/resource/manager.go`, `GetManager`, `HardwareSnapshot`, `RefreshHardware`, `refreshHardware` | Calls `audio/hw.DetectHardware` for RAM/GPU presence, then `queryGPUs` for live GPU totals/usage. Stores `snapshot` and `lastScan`. Singleton initialization refreshes once. |
| GPU discovery | Platform-specific `queryGPUs` | `internal/resource/gpu_info_windows.go`, `detectWindowsGPUs`, `detectViaNvidiaSMI`, `detectViaRegistry`, `detectViaWMI` | On Windows: nvidia-smi first, registry second, WMI last. nvidia-smi supplies total and used; registry/WMI fallback estimates used VRAM from tracked models. |
| System RAM discovery | `audio/hw.DetectHardware` | `internal/audio/hw/detect.go`, `DetectHardware`, `SystemRAM` | `gopsutil/mem.VirtualMemory`; `FreeRAMBytes = Available`. The same hardware package serves resource and audio. |
| Model footprint registry | `ResourceManager` | `internal/resource/manager.go`, `RegisterModel`, `UnregisterModel`, `LoadedModels`, `ModelsByType`, `usedVRAMLocked` | Map keyed `"type:name"`; stores estimated GPU/RAM, priority, last-used, and `IsGPUOffload`. Does not measure. `Loaded` is always true on registration. |
| Allocation advice | `resource.Allocate` | `internal/resource/allocator.go`, `Allocate`, `evictionCandidates` | Computes safe budgets, recommends evictions and CPU fallback. Exposed through `App.RequestAllocation`. No production caller found. |
| Chat memory estimator | `gguf.MemoryModel` | `internal/backend/gguf/memory.go`, `Estimate`, `ClampToHardware`, `ValidateSpawn` | Architecture-aware estimate, layer/context/batch clamping, spawn validation. Used by chat; not by image or audio. |
| Chat/image preflight | `App.CheckModelFit` / `FreeVRAMFor` | `app-vram-arbitration.go` | Refreshes hardware, estimates `llm` or `image`, identifies the other modality as blocking, optionally unloads it and waits 1500 ms. Only two kinds. |
| Settings | App settings store + legacy files + frontend | `app.go` (`NewApp`, `appSettings`, `modelSettings`), `app-model.go` (`loadModelSettingsFromFile`, `saveModelSettingsToFile`), `frontend/dist/js/settings-state.js`, `frontend/dist/js/model-settings.js` | Unified `config/settings.json` and legacy `config/model-settings.json`. Frontend defaults and Wails bindings. |
| Failure diagnostics | `internal/backend/loader` | `app-model.go` failure records, `loader.WriteFailureRecord` | Chat load failures produce records with stage, args, exit status, and tail. Image has its own errors/logging. |

### 2.2 Modality-specific ownership

| Modality | Process / backend | Owner | Registration with `ResourceManager` | Resource gate |
|---|---|---|---|---|
| Chat | `llama-server.exe`, fixed `127.0.0.1:8080` | `App.serverProc`; `app-model.go`, `app-llama.go` | Yes, after spawn/estimate (`app-model.go:504`), unregister only on explicit `UnloadModel` | Clamp + validation; frontend `CheckModelFit` optional |
| Image | `sd-server.exe`, default port 8081 | `sdcpp.Engine`, `job.Manager`, `App.imageEngine` | Only manual `SetActiveImageModel`; **not** normal job auto-load | Frontend `CheckModelFit('image')` optional; engine always `--max-vram -1` |
| Video | `sd-cli.exe` | `sdcpp.Engine` (`videoProc`, `runMu`) | No | No `ResourceManager` path; image preflight does not cover video |
| Audio ASR (Qwen local) | separate `llama-server.exe`, ephemeral port | `audio.Pipeline`, `asr/llama.Recognizer` | No | No shared gate; `-ngl 99` |
| Audio ASR (Whisper) | in-process/backend adapter | `audio.Pipeline`, `asr` registry | No | `CheckHardwareCompatibility` in Whisper recognizer |
| Audio TTS | Piper/Kokoro local; MeloTTS HTTP | `audio.Pipeline`, `tts` adapters | No | No shared gate; TTS profiles mostly CPU/HTTP |

### 2.3 What is shared versus duplicated

Shared:
- hardware detection (`internal/audio/hw`),
- the `ResourceManager` and its Wails bindings,
- `app-vram-arbitration` for chat/image only,
- the GGUF estimator,
- process ownership through `internal/proc` for several backends,
- settings persistence.

Duplicated or independent:
- hardware budgets: chat’s 85/90, allocator’s 90/95, image’s fixed floor, audio’s profile checks;
- process lifecycle: chat `serverProc`, image `serverProc`/`videoProc`, audio `Recognizer.cmd`;
- unload logic: chat unregister-all, image engine shutdown, audio manager shutdown;
- readiness: chat `/health` poll, image probe/exit watcher, audio `/health` poll;
- diagnostics: chat failure records, image stderr/placement, audio tail buffer.

---

## 3. Resource lifecycle: request → load → execute → release

### 3.1 Chat

**Request / preflight**

1. The UI may call `App.CheckModelFit('llm', path)` (`frontend/dist/js/chat-model.js`, `requestVRAMArbitration`). This calls `ResourceManager.RefreshHardware()` and estimates the selected path with the *current* `context_size`/`gpu_layers` (`app-vram-arbitration.go:49-92`).
2. If the estimate does not fit and an image model is loaded, the UI asks for confirmation and calls `FreeVRAMFor('llm', path)`, which unloads the image model, sleeps 1500 ms, and rechecks (`app-vram-arbitration.go:97-120`).
3. If `can_fit` is false but no blocking model is listed, `requestVRAMArbitration` returns `true` and the load proceeds. There is no backend rejection at this point.
4. The UI guard `isModelLoading()` is only a UI guard; it does not serialize Wails calls from other callers.

**Load**

5. `App.LoadModel` (`app-model.go:28`) validates path, runs GGUF architecture admission (`gguf.AdmitModel`), chat task eligibility (`chatpolicy.Evaluate`), and projector identity resolution before stopping the current server.
6. It checks for an existing `serverProc`; if the same model is reported by the HTTP client it returns early; otherwise it kills the existing process and waits (`app-model.go:151-173`). It then checks port 8080, validates `llama-server.exe`, and reads the model profile.
7. `applySafeDefaults(profile)` (`app-model.go:726`) reads `ResourceManager.HardwareSnapshot()` **without refreshing it**, computes:
   - `freeVRAM = FreeVRAMBytes * 0.85`,
   - `freeRAM = FreeRAMBytes * 0.90`,
   caps context to `profile.UsableContext`, builds a fallback `ArchitectureSpec` if the profile has none, and calls `ClampToHardware`.
8. `ClampToHardware` may reduce GPU layers, context, batch, and move layers to CPU. The returned adjustments are printed by `applySafeDefaults`; they are not emitted to the UI, and the mutated settings are not saved at this point.
9. If `profile.ArchitectureSpec != nil`, `ValidateSpawn` runs with the same 85/90 budgets and a 110% tolerance. If the original profile has no `ArchitectureSpec`, this validation block is skipped entirely (`app-model.go:247-306`, `else` at the end).
10. `spawnLlamaServerWithModel` builds args (`app-llama.go:97-155`): `-m`, `-c`, `--threads`, `-b`, `-ngl`, optional `--fit off`, `--mmproj`, rope, flash attention, tensor split only if a string value is present, template, and reasoning flags. It calls `cmd.Start()`, `proc.Attach(cmd)`, and returns.
11. `App.LoadModel` assigns `a.serverProc`, starts `waitMonitor`, and launches the 60-second monitor goroutine. **It then registers the model with `ResourceManager` and returns**, before the monitor has established health (`app-model.go:489-512`).
12. The monitor polls process liveness and `/health` every 500 ms. On success it emits `model:loaded`; on death it writes a failure record and emits `model:error`, but it does not unregister the model.

**Execution**

13. Chat requests go through `HTTPClient` to the llama-server. `buildPromptMessages` assembles system/character/scenario/lore and calls `truncateMessagesToFitContext`. `GetContextSize` reads `/props` `n_ctx`; if that fails, the truncator falls back to 2048 tokens.
14. `truncateMessagesToFitContext` uses `maxChars = contextSize * 4`, keeps the system message and the newest messages that fit by character count, and drops the middle. It does not use a tokenizer, does not account for message framing or special tokens, and does not consult `context_rollover_mode` or `memory_integration`.
15. There is no per-request admission check, no reservation, no queue, and no update of `LastUsed` on requests. Multiple concurrent chat requests can target the same server; the server/backend decides queueing.

**Release**

16. `UnloadModel` kills `serverProc` (direct `Process.Kill`, not `KillTree`), waits, clears `serverProc`, resets loading state, emits `model:unloaded`, and unregisters every `ResourceLLM` entry (`app-model.go:518-535`).
17. `LoadModel` switching kills the old process but does not unregister the old entry. If the new model has a different basename, the old LLM entry remains in the registry. If `LoadModel` fails after killing the old process, the old entry also remains.
18. `OnShutdown` uses `proc.KillTree` for `serverProc`, shuts down image, closes audio, and calls `proc.CloseJob` (`app.go:283-315`). `NewApp` schedules `proc.ReapOrphans` for servers left by a previous crashed instance.

### 3.2 Image

**Request / preflight**

1. The image page may call `CheckModelFit('image', '')` before generation (`frontend/dist/js/image-generator-core.js`, `requestImageVRAMArbitration`). This uses a fixed `1600 MB` base plus `500 MB` overhead plus the largest `*control*net*.safetensors` in the selected/last model directory (`app-vram-arbitration.go:38-44`, `estimateImageVRAM`, `boundControlNetBytes`).
2. The blocking list is populated only if `App.IsModelLoaded()` is true, and it records one LLM entry using the maximum tracked LLM VRAM (`app-vram-arbitration.go:75-89`). Audio and video are not considered.
3. As with chat, a false fit with no blocking entry proceeds. `FreeVRAMFor('image')` calls `UnloadModel()` (which unregisters all LLM entries), sleeps 1500 ms, and rechecks.

**Load**

4. `StartImageJob` (`app-image-generate.go:28`) starts a job goroutine. The job calls `ensureEngineReady` (`app-image-model.go:23`), which initializes the engine if needed, resolves the model/companions, and calls `imageEngine.EnsureReady`.
5. `sdcpp.Engine.LoadModel` (`internal/image/engine/sdcpp/engine.go:129`) holds `loadRun` (serializing loads), begins load ownership, resolves the target, kills/reaps the previous server, validates the checkpoint, resets placement, builds args, spawns `sd-server`, attaches the process, starts the reaper and stdout pump, sets `currentModel = target`, and waits for readiness or process exit.
6. `buildServerArgs` always appends `--max-vram -1` (`engine.go:498-502`), regardless of `opts.MaxVRAM`. It emits `--offload-to-cpu` when selected and `--backend`/`--params-backend` when supplied. It does not emit a free-VRAM admission value.
7. `currentModel` is set **before** `waitReadyOrExit`. A concurrent `EnsureReady` for the same model can see a non-nil matching `currentModel` and return before the server is ready; `LoadModel`'s `loadRun` serializes load calls, but it does not cover this readiness window. The job manager starts each job in its own goroutine (`internal/image/job/job.go`, `Start`), so overlap is possible.

**Execution / measurement**

8. sd-server owns the actual graph-cut placement. Its stdout is parsed for `total params memory size = ... (VRAM ..., RAM ...)`, per-component `RAM|VRAM|N/A`, and the `auto-fit` backend plan (`internal/image/engine/sdcpp/placement.go`).
9. `GetImagePlacement` exposes the measured placement (`app-image-diagnostics.go`). No non-generated frontend caller was found; the image UI calls `GetAppDiagnostics`, which reports only model name/family and Go heap, not placement.
10. Generation is serialized by `genMu` inside the engine; cancellation is job-owned. A second job can still start a load or hit the readiness window described above.

**Release**

11. `UnloadImageModel` calls `imageEngine.Shutdown()` (`app-image.go:176-181`). `Shutdown` cancels video, kills the server tree, waits up to 5 seconds, resets generation state, clears the HTTP client, and sets `currentModel = nil` (`engine.go:413-428`). `reapStoppedProcess` ignores its own timeout result, so `Shutdown` can return nil even if it did not observe the child exit.
12. `UnloadImageModel` does **not** call `ResourceManager.UnregisterModel(ResourceImage, ...)`. `reinitEngine` also shuts down the engine without unregistering. The normal job path never registered, so the registry may be empty, stale, or contain a fixed estimate from an earlier manual load.
13. Video (`engine.generateVideo`) resolves a model with `SelectModel` (no sd-server), stops the image server, and runs `sd-cli` with `runMu.TryLock`. It is not registered with `ResourceManager` and is not covered by chat/image fit checks.

### 3.3 Audio

**Request / activation**

1. Audio is not applied at startup. `NewApp` restores the desired audio settings but leaves `audioManager` untouched (`app.go` initialization; `setAudioDesired`). `ApplyAudioSettings` is the explicit activation path (`app-audio.go:387`).
2. `applyAudioSettings` validates the desired contract, checks Windows microphone availability when ASR is requested, then calls `audioManager.Apply`. The manager builds a new pipeline generation through `NewProductionFactory`, and the pipeline initializes selected ASR/TTS providers.

**Load**

3. Local Qwen3-ASR uses `internal/audio/asr/llama.Recognizer.Initialize`. It chooses a free ephemeral port, starts `llama-server` with `-m`, `--mmproj`, `--host 127.0.0.1`, `--port`, `-c 2048`, `-ngl 99`, `--no-warmup`, and `-t` (`internal/audio/asr/llama/llama.go:78-98`).
4. It starts under `proc.Start`/job ownership, waits up to 120 seconds for `/health`, and releases on failure/cancel. It does not call `CheckHardwareCompatibility`, does not query free VRAM, and does not register with `ResourceManager`.
5. `BaseRecognizer.CheckHardwareCompatibility` (`internal/audio/asr/asr.go:190-220`) checks `HasGPU`, **total** `GPUMemoryBytes` against the profile's `VRAMBytes`, and `AvailableRAMBytes` with 15% headroom. It is used by the Whisper recognizer (`internal/audio/asr/whisper/recognizer.go:271`), not by the local Qwen llama adapter.

**Execution / release**

6. `audio.Manager` provides operations, leases, retired generations, bounded drain, cancellation, and per-capability readiness. `Pipeline.Shutdown` unloads every loader-tracked ASR/TTS provider. This is the strongest lifecycle implementation in the audited area, but it is audio-internal.
7. The ASR process is killed by `proc.KillTree` plus `Process.Kill`, and `stopLocked` waits up to 5 seconds. `audioManager.Close(ctx)` is called from `OnShutdown`. A failure to exit within 5 seconds is reported as an error.
8. No audio path informs `ResourceManager` or `app-vram-arbitration` that a 4 GB-class ASR model is resident.

### 3.4 Video

1. `ensureEngineReady` returns `eng.SelectModel(loadOpts)` for `Mode == "video"` instead of launching sd-server (`app-image-model.go:72-77`).
2. `GenerateVideo` calls `sdcpp.Engine.generateVideo`, which stops sd-server, builds `sd-cli` args, runs the CLI under job ownership, parses progress, and kills only when the owning job cancels.
3. There is no `ResourceManager` registration and no `CheckModelFit` kind for video. The image preflight does not know a video render is running. This is outside the initial Chat/Image release scope but is shared image-engine infrastructure.

---

## 4. Current allocation and safety rules

### 4.1 Chat

**Budgets**

```text
freeVRAMMB = floor(HardwareSnapshot.FreeVRAMBytes * 0.85 / 1024 / 1024)
freeRAMMB  = floor(HardwareSnapshot.FreeRAMBytes  * 0.90 / 1024 / 1024)
```

Source: `app-model.go:250-251` (validation) and `app-model.go:734-735` (`applySafeDefaults`). These are percentages of **reported free** memory at the time of the snapshot, not total memory. The snapshot is read without a refresh in `applySafeDefaults`.

**Suspicious-detection fallback**

If total VRAM is over 1 GiB and computed free VRAM is under 1024 MB, `LoadModel` replaces the budget with `TotalVRAMBytes * 0.70` (`app-model.go:255-262`). This is a heuristic, not a measurement of free memory. It can authorize a load when free detection failed.

**Estimator**

`gguf.MemoryModel.Estimate` (`internal/backend/gguf/memory.go`) computes:

- `weightMB` from the model file/quantization;
- `kvMB` from architecture and context;
- `actMB` from architecture and batch;
- GPU: `weightMB * gpuLayers/(Layers+1)`, `kvMB * gpuLayers/Layers`, half activations, `400 MB` overhead;
- RAM: residual weights, residual KV, half activations, `500 MB` overhead;
- MoE reduction (active experts);
- sliding-window KV reduction when `SlidingWindowSize < ctxSize`.

The estimate is an analytical model, not a driver measurement. Its confidence and notes are descriptive; no calibration data is attached to the result.

**Clamp order**

`ClampToHardware` (`internal/backend/gguf/memory.go`) does, in order:

1. Clamp `gpuLayers` to `[0, Layers]`; negative means all layers.
2. Reduce GPU layers stepwise until `Estimate.GPUMB <= freeVRAMMB`; if zero, mark CPU-only.
3. Reduce context by halving until `Estimate.RAMMB <= freeRAMMB`; floor `2048`, or lower only when the model's own context is below 2048 (down to 512).
4. Halve batch while RAM is exceeded and batch > 128.
5. Move more layers to CPU while RAM is exceeded and layers > 0.
6. For sliding-window models, force all layers on GPU if all fit, otherwise CPU-only, because partial offload is treated as unsafe for this backend.
7. Return the final estimate and adjustment strings.

`ReduceGPULayers` steps by 16/8/4/2/1 depending on current count (`memory.go`).

**Spawn validation**

`ValidateSpawn` compares the estimate with the same `freeVRAMMB`/`freeRAMMB`:

- critical if `estimatedGPU > freeVRAM * 1.1`;
- critical if `estimatedRAM > freeRAM * 1.1`;
- warning instead of critical when `freeVRAM < 256 && freeRAM < 512` (suspicion that detection is wrong);
- `force_load=true` bypasses a critical result (`app-model.go:277-300`).

Because clamping targets 85%/90%, a critical result should normally be reachable only when the clamp could not fit (e.g. context floor, batch floor, zero layers, or an invalid/absent spec). `force_load` does **not** bypass `ClampToHardware`; it only bypasses the validation result. This is an important semantic mismatch: a user who sets `force_load` to preserve requested layers/context still receives clamped settings.

**Actual backend behavior**

`buildLlamaServerArgs` passes `-ngl`, `-c`, `-b`, `--threads`, and so on. For non-SWA models, `fitParams=true`, so no `--fit off` is passed; llama.cpp may perform its own parameter fitting/placement. Howl does not read the resulting placement for chat. The registration estimate is therefore a pre-launch estimate, not an observed allocation.

**Registration**

```text
gpuMB = EstimateMemory(profile, gpuLayers, ctxSize).GPUMB
ramMB = EstimateMemory(profile, gpuLayers, ctxSize).RAMMB
isGPUOffload = gpuLayers > 0
```

Source: `app-model.go:489-511`. If `profile == nil`, it falls back to `gpuLayers*100 MB` GPU and `4000 MB` RAM. `gpuLayers` is read from `modelSettings` with an `int` type assertion; a value that has not yet been normalized to `int` (for example, a raw JSON `float64`) yields zero and makes `isGPUOffload` false.

### 4.2 Resource allocator (advisory)

`internal/resource/allocator.go`:

```text
freeVRAM = FreeVRAMBytes * 0.90
freeRAM  = FreeRAMBytes  * 0.95
```

It:

- returns `UseGPU=true`, `GPULayers=-1` by default;
- if GPU is preferred and VRAM is sufficient, succeeds;
- otherwise considers lower-priority, different-type GPU-offloaded models for eviction (LLM priority 0, TTS 5, Image 8); never evicts the same type;
- if eviction would free enough, returns names and confidence 0.75;
- otherwise falls back to CPU if `AllowCPU`, setting GPU memory to zero and RAM to `RequiredVRAM + RequiredRAM`, confidence 0.70;
- if RAM is still over budget, sets `ReduceBatch` and `ReduceContext`, halves the reported RAM need, and reduces confidence.

It does not unload anything itself, does not reserve, and does not update a central state. `RequestAllocation` is exposed (`app-resource.go:123-154`), but no production caller was found outside generated Wails bindings.

### 4.3 Image

**Preflight estimate**

```text
requiredMB = 1600 (imageBaseVRAMMB)
           + 500  (imageOverheadMB)
           + largest *control*net*.safetensors in the model directory (if any)
```

Source: `app-vram-arbitration.go:38-44`, `estimateImageVRAM`, `boundControlNetBytes`. This is a machine-specific floor derived from an SDXL/GTX 1080 observation in `STAGE_STATUS.md`, not an architecture/model-specific estimator. It does not inspect the checkpoint size, family, resolution, text encoders, or actual offload plan. It also ignores the passed model path when the engine has a last model, because `boundControlNetBytes` falls back to `lastImageModel`/configured models directory.

**Admission**

No backend admission. The engine always uses `--max-vram -1`, never `opts.MaxVRAM`. The image engine relies on sd-server’s auto-fit to place weights across VRAM/RAM/disk.

**Registration**

Manual `SetActiveImageModel`:

```text
GPUMemoryMB = 4000
SystemRAMMB = 2000
isGPUOffload = !OffloadToCPU
```

Source: `app-image-model.go:143-150`. These constants are not derived from the model, the placement, or the free budget. Normal generation through `StartImageJob`/`GenerateImage` does not register at all.

### 4.4 Audio

- Qwen3-ASR local: fixed `-ngl 99`, `-c 2048`, `-t` default 6, no `--fit`, no free-VRAM gate.
- ASR profile Qwen3-ASR 1.7B local declares `VRAMBytes = 4 GiB` (`internal/audio/types/asr_models.go`).
- `CheckHardwareCompatibility` uses total GPU memory, not free; it is not called by the local Qwen adapter.
- TTS profiles: Piper/Kokoro local with zero VRAM profiles; MeloTTS HTTP. They are not registered with `ResourceManager`.

### 4.5 What the “98%” claim maps to

No allocation code uses 98%. The plausible explanations for a historical 98% observation are:

- an external GPU utilization reading (task manager/nvidia-smi) at a moment when Howl plus other applications occupied the card;
- a backend placement decision by llama.cpp `--fit` or sd.cpp auto-fit;
- a stale/incorrect detection reading, not an allocation limit;
- a different build or setting than the current tree.

This cannot be resolved from source alone. The implemented limits are 85/90 for chat, 90/95 for the unused allocator, 110% tolerance for validation, and an unquantified `-1` for image.

---

## 5. Cross-modality interaction map

### 5.1 Chat ↔ Image

- `app-vram-arbitration.CheckModelFit("llm")` sees a loaded image model by asking the image engine directly (`loadedImageModelInfo`), not by trusting `ResourceManager` (`app-vram-arbitration.go:126-138`).
- `CheckModelFit("image")` sees a chat model through `IsModelLoaded()` and the maximum tracked `ResourceLLM` estimate.
- `FreeVRAMFor` unloads the other modality and waits 1500 ms, but does not verify actual freed VRAM, does not retry with a timeout, and does not reserve the freed space.
- Only UI-initiated loads/generations call this. Direct Wails calls do not.
- Image’s normal auto-load path is invisible to `ResourceManager`, so `CheckModelFit("llm")` still asks the engine and can see it; but `GetLoadedModels`, `Summary`, and allocator advice cannot.

### 5.2 Chat/Image ↔ Audio

- No interaction. Audio starts a separate llama-server with `-ngl 99` and never registers or checks fit. Chat/image preflight does not know it exists. Audio manager ownership does not include GPU arbitration.
- A local Qwen3-ASR load can therefore occur while chat or image is resident. Whether it OOMs depends on actual VRAM, driver behavior, and model placement; that is unestablished read-only.

### 5.3 Image ↔ Video

- Video shares `sdcpp.Engine` state and the image process family. `generateVideo` stops sd-server before running sd-cli, so it frees the image server intentionally, but no `ResourceManager` update occurs and no video job is registered.
- `runMu` prevents two video renders; it does not coordinate with chat/audio.

### 5.4 Shared process ownership

- `proc.Attach`/job-object ownership is used for chat llama-server, image sd-server/sd-cli, and audio ASR. `OnShutdown` calls `proc.CloseJob`. This is a strong crash-cleanup mechanism.
- `proc.ReapOrphans` is scheduled at app construction for servers left by a prior crashed instance.
- These mitigations reduce orphan risk but do not provide admission, accounting, or performance control.

---

## 6. Historical observations reconciled

The following are user observations, not audit measurements. No original artifacts establishing the numbers were located in this audit. They remain **unestablished** unless separately measured.

### 6.1 Howl ~60% VRAM at 22K context vs SillyTavern ~90% at 20K

**What source can explain**

- Howl’s chat path caps its estimate at 85% of **free** VRAM and then clamps layers/context/batch. At 22 layers and 22K context, the estimator may produce a lower GPU requirement than a configuration that offloads every layer; the actual `-ngl` setting and llama.cpp `--fit` determine placement.
- `app-model.go`’s current persisted settings are `gpu_layers: 22`, `context_size: 20480`. The estimator’s layer fraction is `gpuLayers/(Layers+1)` for weights and `gpuLayers/Layers` for KV, so a partial layer count can leave substantial VRAM unused even when the context is large.
- `CheckModelFit` uses the same estimator and calls it “required” memory; it is not a driver measurement.
- SillyTavern’s number could reflect a different layer count, KV cache type, batch/context, backend flags, or measurement point.

**What is unestablished**

- The two products’ exact launch flags, quantization, KV cache type, context, batch, offload settings, and measurement methodology.
- Whether Howl’s clamp actually reduced the user’s 22K request, and to what.
- Whether llama.cpp `--fit` changed placement after Howl’s estimate.

**Conclusion:** source mechanisms are compatible with a conservative Howl number, but the 60/90 comparison cannot be attributed to a specific Howl constant or proven from source.

### 6.2 “98% of free VRAM” cap

**Source finding**

- No 98% allocation constant or policy exists in the inspected Go code.
- The implemented budgets are 85/90 (chat), 90/95 (allocator), 110% validation tolerance, and `--max-vram -1` (image).
- The 98% claim is therefore **not source-supported** as a Howl allocation policy. It may be a utilization measurement or an external/backend behavior.

**Status:** contradicted as an implementation claim; the origin of the historical 98% number is unestablished.

### 6.3 Qwen 27B with ~4 GB RAM, high CPU, heavy GPU, thinking under five minutes

**What source can explain**

- `ClampToHardware` can reduce context/batch and move layers to CPU, including CPU-only, when the RAM estimate exceeds 90% of reported available RAM.
- `--fit` is left enabled for non-SWA models; llama.cpp may place weights in RAM and stream them.
- The health monitor only checks `/health`; it does not measure CPU/GPU utilization or tokens/second. A CPU-heavy/RAM-spilled configuration can be accepted.
- The estimator uses `AvailableRAMBytes`, not a reservation; another process can consume the remaining RAM after the check.

**What is unestablished**

- The actual 4 GB figure, model quantization, context, layers, page-file/mmap behavior, and whether the run was CPU-bound or GPU-spilling.
- The five-minute timing and whether it was measured end-to-end.
- Whether Howl logged any clamp adjustment for that run. Adjustments are console-only and are not persisted.

**Conclusion:** the architecture can produce a slow CPU/RAM-spilled run that still becomes “ready.” The specific Qwen observation is not verified.

### 6.4 Three-character trial exceeded the expected 4K budget without rollover

**What source can explain**

- The app’s current persisted settings are `context_size: 20480`, not 4K (`config/settings.json`, `config/model-settings.json`). The UI default is 4096 (`frontend/dist/js/model-settings.js`), but the unified/legacy settings override it.
- `truncateMessagesToFitContext` reads `n_ctx` from llama-server `/props`; if the server was started with 20480, the character budget is `20480 * 4`, not 4K.
- The truncator drops oldest messages by character count; it has no tokenizer, no rollover event, and no memory-guided policy. `context_rollover_mode` is not consulted anywhere in Go.
- A three-character scenario can therefore carry more than 4K worth of content if the loaded server context is larger and the rough character estimate remains under the limit.

**What is unestablished**

- The actual token count, the server’s reported `n_ctx`, and whether the user’s expectation of 4K came from the UI default or a prior configuration.
- Whether the additional content was accepted because of context size, because the truncator undercounted, or because the model tolerated it.

**Conclusion:** the source strongly suggests the “4K budget” was not an enforced 4K limit in the current configuration, and “rollover” is not implemented. The exact observation remains unverified.

### 6.5 Chat GPU-offload degradation was CPU-bound and slow

**Source finding**

- Chat `LoadModel` accepts health, not performance. `IsModelLoaded` is an HTTP state check.
- `HTTPClient.GetUsedMemory`/`GetTotalMemory` return zero; there is no chat placement parser or tokens/second metric.
- `ClampToHardware` can set `gpu_layers=0`; `--fit` may also offload. Both can yield a slow CPU-bound run that passes health.
- No user-visible runtime performance warning or automatic rollback exists.

**Conclusion:** the architecture cannot distinguish “loaded and acceptable” from “loaded but unusably slow.” This is a confirmed source-level capability gap, regardless of the historical run’s exact cause.

---

## 7. Defect and risk register

Severity reflects potential user impact if the source-derived behavior occurs at runtime; it is not a claim that every path has been reproduced. “Grade” uses the vocabulary in §0.

| ID | Severity | Grade | Finding | Evidence and consequence |
|---|---|---|---|---|
| R1 | **High** | Source-derived | No global admission/reservation across chat, image, audio, and video. | Separate backends and no shared lock. Frontend preflight is advisory; audio bypasses it. Concurrent starts can overcommit VRAM/RAM. Actual OOM unestablished. |
| R2 | **High** | Source-derived | Chat load is not serialized; `serverProc` is unsynchronized. | `LoadModel` reads/writes `a.serverProc` without a mutex and has no load mutex (`app-model.go:151-173`, `:338`, `:518-524`). Two calls can both pass the existing-process check, kill/overwrite each other, or leave an orphan. UI `isModelLoading()` is not a backend guard. |
| R3 | **High** | Source-derived | Model registry lifecycle is incomplete and stale. | Chat switching does not unregister the old LLM; image normal job loads do not register; image unload/reinit do not unregister; registration values are estimates/constants. `GetLoadedModels`, `Summary`, and allocator advice can be wrong. |
| R4 | **High** | Source-derived | Chat registers before health and never unregisters on load failure. | `RegisterModel` at `app-model.go:504` occurs before the monitor’s `/health` success. The monitor’s failure path writes diagnostics but does not call `UnregisterModel`. Failed/partial loads remain “loaded.” |
| R5 | **High** | Source-derived | Audio ASR is outside resource arbitration. | Local Qwen ASR uses `-ngl 99`, no free-VRAM check, no registration (`internal/audio/asr/llama/llama.go:93`). Chat/image preflight ignores it. Co-residency can overcommit; actual outcome unestablished. |
| R6 | **High** | Source-derived | CPU/RAM spillover has no performance acceptance gate. | Health-only monitors; no tokens/sec, CPU/GPU, or placement telemetry for chat. `ClampToHardware` can produce CPU-only. A technically successful but unusable load is accepted. |
| R7 | **Medium** | Source-derived | Preflight can proceed despite `can_fit=false`. | Both frontends return true when there is no blocking model (`chat-model.js:332-347`; `image-generator-core.js:592-607`). The backend does not re-check. “May fall back to RAM” is shown only in the message; the user is not forced to confirm. |
| R8 | **Medium** | Source-derived | Automatic clamping is silent and `force_load` has misleading semantics. | `applySafeDefaults` logs adjustments and discards them; no event/UI. `force_load` bypasses validation but not clamping. User cannot see what changed or why; requested values may not be used. |
| R9 | **Medium** | Source-derived | Image ignores its own VRAM cap setting. | `MaxVRAM` parsed but `buildServerArgs` always uses `--max-vram -1` (`app-image-model.go:211`; `engine.go:502`). The type and engine comments disagree about `-1` semantics. |
| R10 | **Medium** | Source-derived | Image fit estimate is a fixed floor, not a model estimate. | 1600 + 500 + directory ControlNet; no checkpoint/family/resolution/encoder/placement analysis. Manual registration uses fixed 4000/2000. Can under- or over-state. |
| R11 | **Medium** | Source-derived | Context truncation is approximate and “rollover” is not wired. | `contextSize*4` characters, no tokenizer, no special-token accounting; `context_rollover_mode`/`memory_integration` only in frontend. Could exceed real context or discard more than intended. |
| R12 | **Medium** | Source-derived | Hardware snapshot is stale unless explicitly refreshed. | `GetResourceSnapshot` returns stored snapshot; `applySafeDefaults` reads without refresh; no background loop. A long-running app can admit against old free memory. |
| R13 | **Medium** | Source-derived | GPU fallback accounting is circular. | When live query fails, `estimateUsedVRAM` sums the same tracked models the manager is supposed to supervise; stale/missing registrations feed back into free-VRAM estimates. |
| R14 | **Medium** | Source-derived | Several persisted settings are not applied. | `buildLlamaServerArgs` does not emit `kv_cache_size`, `kv_cache_type`, `offload_kqv`, `use_mmap`, `use_mlock`, or `main_gpu`; UI `context_rollover_mode`, `memory_integration`, `system_prompt_cache`, `tokenizer_override`, `model_family`, `model_variant` have no backend use. `tensor_split` is only read as a string, while the UI supplies a number. |
| R15 | **Medium** | Source-derived | Image shutdown does not prove the child exited. | `reapStoppedProcess` ignores a 5-second timeout and `Shutdown` returns nil; `UnloadImageModel` does not unregister. A still-live process can be reported as unloaded. |
| R16 | **Low** | Source-derived | Chat unload does not cancel the in-flight stream. | `UnloadModel` kills the process but does not call `streamCancel`; the request sees a transport failure rather than a clean cancellation. `streamCancel` is only used by `AbortStream` (`app-prompt.go:421`). |
| R17 | **Low** | Inferred | Image can expose a model as current before readiness. | `currentModel` is set before `waitReadyOrExit`; a second `EnsureReady` can see it. `loadRun` serializes loads but does not cover this window; job starts are concurrent. |
| R18 | **Low** | Source-derived | Chat and image use fixed ports/hard-coded process assumptions. | Chat port 8080 and image 8081; one chat server. Port conflicts are reported, but there is no dynamic chat port or multi-instance coordination. |
| R19 | **Low** | Source-derived | Loaded state is not persisted or reconstructed. | ResourceManager is in-memory; image persists `lastImageModel` but not a loaded process. After restart the registry is empty and no model is auto-loaded. Settings persist, loaded state does not. |
| R20 | **Unestablished** | Unestablished | VRAM fragmentation and allocation failure behavior. | The app does not allocate tensors; llama.cpp/sd.cpp own allocation. No fragmentation telemetry or failure taxonomy exists. Cannot be established read-only. |

**Dangerous-issue disposition:** R1/R5 are the potentially host-wide risks, but no runtime OOM was observed or reproduced in this read-only audit. They are documented, not repaired.

---

## 8. Missing capabilities, ranked by practical impact

1. **A single cross-modality resource coordinator with real admission.** One authority should serialize or queue backend starts, reserve a budget, account for chat/image/audio/video, and own the decision to unload. The current manager is a passive registry; the allocator is unused.
2. **Truthful lifecycle accounting.** Every successful load must be registered with a measured or clearly labeled estimate; every failure, switch, unload, reinit, cancel, and shutdown must unregister exactly once. Registration must follow health, not spawn.
3. **Post-load placement and runtime performance acceptance.** Read actual VRAM/RAM/offload from the backend where possible, compare with the estimate, and reject or warn on an unusable CPU-spilled configuration. Chat needs the equivalent of image’s placement parser, or a documented limitation.
4. **User-visible automatic changes.** Emit clamp adjustments, show resulting layers/context/batch and why, distinguish “requested” from “applied,” and make `force_load` semantics honest.
5. **Image and audio integration.** Remove the fixed image cap/registration constants; honor `MaxVRAM` or delete it; register auto-loaded images; include audio in fit calculations or declare audio mutually exclusive.
6. **Exact context accounting.** Tokenizer-based truncation, special-token and message-framing accounting, and a real implementation (or removal) of the advertised rollover modes.
7. **Persistence/reconciliation after restart.** Persist or reconstruct what was actually loaded and clearly reset stale state after a crash. Do not imply a model is resident when only its name is remembered.
8. **Observability and a reproducible test harness.** A read-only diagnostic report of snapshot, registered models, actual processes, placement, and clamp decisions would make every later claim testable.

---

## 9. Recommended smallest next stage

### Stage R1 — Read-only runtime validation (recommended first)

**Purpose:** establish which source-derived risks are real before changing architecture.

**Prerequisites**

- Owner authorization for runtime validation, separate from this read-only audit.
- A test machine with a known discrete GPU and no unrelated GPU workloads.
- A small, fixed model set (one chat GGUF, one image checkpoint, one local ASR model/projector), with recorded hashes and sizes.
- A clean settings backup and a disposable user-data profile.
- Permission to build/run the app and inspect `nvidia-smi`, process lists, and logs.
- No TTS/G2P artifacts; do not touch the TTS branch.

**Minimum test matrix**

1. **Chat load:** record `nvidia-smi` before, during, and after `LoadModel`; record the chosen `gpu_layers`/context/batch and the actual llama-server args; compare actual placement/VRAM with the registered estimate.
2. **Chat switch/failure:** load model A, then B; unload; force a startup failure; inspect `GetLoadedModels`, process list, and GPU memory after each step.
3. **Concurrent loads:** issue two `LoadModel` calls without the UI guard; verify one server, one registration, no orphan, and predictable error behavior.
4. **Image auto-load:** generate through `StartImageJob` without `SetActiveImageModel`; verify whether the image model appears in `GetLoadedModels`, and whether unload clears it.
5. **Audio co-residency:** with chat loaded, apply local Qwen ASR; record VRAM before/after, whether ASR starts, and whether chat remains healthy. Repeat with image loaded.
6. **Preflight bypass:** force a `can_fit=false` result with no blocking model; verify the frontend proceeds and what the backend does.
7. **Context:** with 4K and 20K settings, measure actual token counts and the truncator’s behavior; verify whether `context_rollover_mode` changes anything.
8. **Restart/crash:** kill the app during a load and after a load; verify `ReapOrphans`, job-object cleanup, and registry reconstruction.

**Acceptance criteria for R1**

- Every scenario produces a recorded before/after snapshot and an actual process/VRAM measurement, not just a UI state.
- No test may be recorded as “passed” if it only observed a health endpoint; performance/placement claims require measured data.
- The report identifies which risks are reproduced, which are source-only, and which remain unestablished.
- No production source, settings, models, or user data are modified; all artifacts are under a new investigation directory with hashes.

### Stage R2 — Smallest implementation candidate (only after R1 and explicit authorization)

Do **not** start this from the current audit. The smallest coherent change would be a single load coordinator shared by chat/image/audio:

- one mutex/queue around backend starts;
- a reservation object with explicit release;
- register-after-health and unregister on every terminal path;
- unregister on image shutdown/reinit and chat switch;
- a post-load reconciliation hook that records measured placement where the backend exposes it;
- no UI redesign, no new model formats, no changes to TTS/G2P.

**Acceptance criteria for R2**

- Concurrent starts cannot create two owners for one port or two registrations for one process.
- Every successful load has exactly one registry entry, and every unload/failure removes it.
- Image auto-load and audio ASR are either represented or explicitly rejected by the coordinator.
- The UI can display requested vs applied vs measured values without a new page.
- Existing model load, generation, and shutdown tests pass unchanged; new tests cover the added lifecycle transitions.
- A full app shutdown leaves no owned backend process and no registry entry.

This is a proposal for a later authorized stage, not an implementation performed here.

---

## 10. Source index for the most important findings

| Finding | Exact source |
|---|---|
| Chat load lifecycle | `app-model.go`, `LoadModel`, `UnloadModel`; `app-llama.go`, `buildLlamaServerArgs`, `spawnLlamaServerWithModel` |
| Chat 85/90 budgets | `app-model.go:250-251`, `:734-735`; `applySafeDefaults` |
| Chat estimator / clamp / validation | `internal/backend/gguf/memory.go`, `MemoryModel.Estimate`, `ClampToHardware`, `ValidateSpawn` |
| Chat registration before health | `app-model.go:338`, `:489-512`; monitor goroutine starts at `:343` |
| Resource manager | `internal/resource/manager.go`, `GetManager`, `RegisterModel`, `UnregisterModel`, `refreshHardware`, `fallbackGPUValues`, `usedVRAMLocked` |
| Allocator 90/95 | `internal/resource/allocator.go:21-22`, `Allocate` |
| Resource bindings | `app-resource.go`, `GetResourceSnapshot`, `GetLoadedModels`, `RequestAllocation` |
| Chat/image preflight | `app-vram-arbitration.go`, `CheckModelFit`, `FreeVRAMFor`, `estimateLLMVRAM`, `estimateImageVRAM` |
| Frontend preflight bypass | `frontend/dist/js/chat-model.js:332-347`; `frontend/dist/js/image-generator-core.js:592-607` |
| Image engine load / shutdown | `internal/image/engine/sdcpp/engine.go`, `LoadModel`, `Shutdown`, `buildServerArgs` |
| Image hard-coded registration | `app-image-model.go:143-150`; `ensureEngineReady` has no registration |
| Image unload no unregister | `app-image.go:176-181`; `reinitEngine` at `:157-171` |
| Image `--max-vram -1` | `internal/image/engine/sdcpp/engine.go:498-502`; `MaxVRAM` parse at `app-image-model.go:211-212` |
| Image placement | `internal/image/engine/sdcpp/placement.go`; `app-image-diagnostics.go`, `GetImagePlacement` |
| Audio manager lifecycle | `internal/audio/manager.go`, `Apply`, `Acquire`, `Close`; `internal/audio/pipeline.go`, `Shutdown` |
| Audio ASR launch | `internal/audio/asr/llama/llama.go:78-98` (`-ngl 99`) |
| Audio hardware check | `internal/audio/asr/asr.go:190-220`; call site `internal/audio/asr/whisper/recognizer.go:271` |
| Context truncation | `app-prompt.go:572-646`, call at `:771`; `internal/backend/llama/http_client.go`, `GetContextSize` |
| Settings persistence | `app.go`, `NewApp`; `app-model.go`, `loadModelSettingsFromFile`, `saveModelSettingsToFile`, `SetModelOptions`; `config/settings.json`; `config/model-settings.json` |
| Settings UI | `frontend/dist/index.html`; `frontend/dist/js/model-settings.js` |
| Video path | `internal/image/engine/sdcpp/video.go`, `SelectModel`, `generateVideo`, `runVideoCLI` |
| Process ownership | `app.go`, `OnShutdown`; `internal/proc` usage in `app-llama.go`, image engine, audio ASR |

---

## 11. Read-only boundaries, changed files, and hashes

### 11.1 What was and was not done

Performed: directory/file reads, targeted `grep`/`rg` searches, `git rev-parse`, `git status`, and `sha256sum`.

Not performed: application launch, build, `go test`, benchmark, model load, generation, audio run, image run, network download, dependency installation, configuration write, settings write, model write, user-data write, commit, push, or production source modification.

The TTS/G2P feasibility artifacts under `build/tmp/kokoro-g2p-feasibility/` were not read, modified, or used as evidence for this report. The existing TTS branch was not altered.

### 11.2 Changed-file inventory

**Created by this audit (only):**

- `docs/review/INVESTIGATION-system-resource-management-architecture-2026-11-06.md`

**Not modified by this audit:**

- `app-model.go`
- `app-llama.go`
- `app-image.go`
- `app-image-model.go`
- `app-audio.go`
- `internal/resource/manager.go`
- `internal/resource/allocator.go`
- `internal/backend/gguf/memory.go`
- `internal/audio/asr/llama/llama.go`
- all other production source, configuration, models, assets, and user data

**Pre-existing worktree condition:** the repository already contained extensive modifications and untracked artifacts before this audit, including a modified `app-audio.go` and many unrelated files. Those are not attributable to this audit. The audit added only the report path above.

### 11.3 Hashes captured for key production files

The following SHA-256 values were captured while preparing this report. They establish the audit-time content of the principal files cited above; they are not a substitute for a pre-audit baseline, because the worktree already had modifications.

```text
54ab5901ad3876c7e05e51340abfa66aec40c1633ff6c6f513074a5c9db384eb  app-model.go
5ee8531edd7a54937297c53a56cf8b9ace40cbfe86e54caacbd9238108978331  app-llama.go
1f0e364954002d31106a89e4029ea81e332822de971faa26ed50436d0f23a201  app-image.go
c1ff3a1783a877bb0f7d22619618f5f15b6e0db21c9e269b49c60e8e5642f90c  app-image-model.go
ff5041f19c7264cbc3cfae4baa91b49d1f58189c021a9c8155bb1833713cc9d9  app-audio.go
34261573e18c84a3c9d0cd5e3c3cd63cd35edbd965935ff9f3bc29ff51e392ef  internal/resource/manager.go
ffed5833a2bbab054c2254c765da469e85c6ee4013be9ca6e4132ff034c93398  internal/resource/allocator.go
aed6b125c95809b44b257c4db20ab159f56a349594df18956684be8c0930149e  internal/backend/gguf/memory.go
eda8974ac87895e9c22a28485f37eda2b0f28fa303852aedaf2ef6f37351c5da  internal/audio/asr/llama/llama.go
```

### 11.4 What could not be established read-only

- Actual VRAM/RAM usage, OOM occurrence, allocation failure, fragmentation, or driver recovery.
- Whether llama.cpp `--fit` changed Howl’s clamped values at runtime.
- The actual placement of chat weights (VRAM vs RAM vs CPU); chat has no placement reader.
- The actual effect of image `--max-vram -1`, and whether the engine comment (“reserves 1 GiB”) or the type comment (“use all available”) is correct for the shipped binary.
- Whether audio co-residency with chat/image causes OOM, and how `-ngl 99` behaves under contention.
- Whether the historical 60/90, 98%, Qwen 27B, or 4K observations are reproducible; their original artifacts were not located.
- Whether the current frontend’s preflight is the only path users take; direct Wails or scripted calls bypass it.
- Whether persisted settings produce the same effective values after restart, given that `applySafeDefaults` mutates memory but does not save.
- Whether the UI exposes `GetResourceSnapshot`, `GetLoadedModels`, `RequestAllocation`, `GetResourceSummary`, or `GetImagePlacement` anywhere outside generated bindings; no non-generated caller was found.

---

## 12. Bottom line

Howl currently has a useful **resource-awareness layer**: hardware discovery, a shared snapshot, architecture-aware estimation, conservative chat clamping, an unused eviction advisor, and a chat/image preflight. It does **not** yet have a resource-management architecture in the strict sense: there is no single admission controller, no reservation, no cross-modality lock, no reliable lifecycle registry, no measured chat placement, and no performance acceptance test.

The 98% ceiling should not be treated as implemented. The implemented numbers are 85/90 for chat, 90/95 for the allocator, 110% for validation tolerance, and an unquantified image `-1`. The most important next step is not a redesign; it is a bounded runtime validation of load ownership, actual placement, registration truth, and cross-modality contention. Until that evidence exists, the safest statement is: **Howl can often keep a single chat or image load within a conservative estimate, but it cannot yet guarantee which backend owns the GPU, what was actually placed where, or that a technically loaded model will perform acceptably.**
