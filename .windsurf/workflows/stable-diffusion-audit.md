---
description: Stable Diffusion Integration — Audit Results & Task Breakdown (May 2026)
---

# HOWL Chat — Stable Diffusion Integration Audit

**Plan:** `stable-diffusion-integration-064de1.md`  
**Audit Date:** May 18, 2026  
**Files Audited:** `internal/image/*`, `app.go`, `frontend/dist/image-generator.*`, `frontend/dist/js/image-generator*.js`, `frontend/dist/css/image-generator*.css`

---

## Executive Summary

| Phase | Status | Notes |
|-------|--------|-------|
| 0 — Foundations | **COMPLETE** | All types, interface, storage layout defined |
| 1 — SD.cpp Engine (txt2img) | **COMPLETE** | Subprocess + HTTP client, model loading, generation working |
| 2 — Frontend txt2img + Gallery | **COMPLETE** | Full 3-col UI, collapsible panels, history, diagnostics, presets |
| 3 — img2img Support | **NOT STARTED** | Backend stub exists; frontend button is decorative |
| 4 — Inpainting | **NOT STARTED** | Backend stub exists; frontend button is decorative |
| 5 — Job System & Progress | **NOT STARTED** | ImageJob type exists; no queue, worker, or bindings |
| 6 — World/Scenario/Character Integration | **NOT STARTED** | No bindings, no frontend actions |
| 7 — Advanced Polish | **PARTIAL** | Model presets, style presets done; batch/highres/LoRA partially wired |

**Critical Bugs Found:**
1. `app.go:GenerateImage` ignores `opts.Mode` — always calls `GenerateTxt2Img` regardless of requested mode
2. `client.go:txt2img` ignores `enable_hr`, `loras`, `tiling`, and other advanced opts in the HTTP payload
3. No `UnloadImageModel` Wails binding exists, but frontend `image-generator-ui.js:179` calls it
4. Workflow mode buttons (img2img, inpaint, outpaint, upscale) are UI-only — no mode state propagation to generation

---

## Detailed Phase-by-Phase Audit

### Phase 0 — Foundations and Contracts ✅ COMPLETE

| Item | File | Status |
|------|------|--------|
| ImageOptions struct | `internal/image/types/types.go:4` | ✅ All fields present including Mode, Strength, InitImagePath, MaskPath, BatchCount, LoRAs |
| ModelInfo struct | `internal/image/types/types.go:24` | ✅ Name, Path, Family, Capabilities, ResolutionHint, Files |
| ImageResult struct | `internal/image/types/types.go:36` | ✅ Path, ThumbPath, Metadata |
| ImageMetadata struct | `internal/image/types/types.go:51` | ✅ All fields including File, Mode, Source |
| ImageSource struct | `internal/image/types/types.go:44` | ✅ WorldID, ScenarioID, CharacterID |
| EngineConfig struct | `internal/image/types/types.go:69` | ✅ BinaryPath, ModelsDir, OutputDir, TempDir, Port |
| ImageJob struct | `internal/image/types/types.go:78` | ✅ ID, Status, Progress, Options, Result, Error, CreatedAt |
| ImageEngine interface | `internal/image/engine/engine.go:8` | ✅ Initialize, ListModels, ListLoRAs, LoadModel, GenerateTxt2Img, GenerateImg2Img, GenerateInpaint, GetLoadedModel, Shutdown |
| Gallery Service | `internal/image/gallery/gallery.go` | ✅ Save, List, Delete, BaseDir |

### Phase 1 — SD.cpp Engine Adapter (txt2img) ✅ COMPLETE

| Item | File | Status |
|------|------|--------|
| Engine.Initialize | `internal/image/engine/sdcpp/engine.go:33` | ✅ Validates binary, creates dirs, sets up HTTP client |
| Engine.ListModels | `internal/image/engine/sdcpp/engine.go` | ✅ Scans models/, detects FLUX bundles vs single-file, returns ModelInfo slices |
| Engine.LoadModel | `internal/image/engine/sdcpp/engine.go:70` | ✅ Spawns sd-server subprocess, waits for ready, captures stderr for diagnostics |
| client.txt2img | `internal/image/engine/sdcpp/client.go:60` | ✅ HTTP POST to /sdapi/v1/txt2img, returns ImageResult with metadata |
| client.waitForReady | `internal/image/engine/sdcpp/client.go:27` | ✅ Polls /sdapi/v1/samplers until responsive |
| Wails: ListImageModels | `app.go:3024` | ✅ |
| Wails: ListImageLoRAs | `app.go:3039` | ✅ |
| Wails: SetActiveImageModel | `app.go:3054` | ✅ With panic recovery |
| Wails: GenerateImage | `app.go:3068` | ⚠️ **BUG**: Always calls GenerateTxt2Img regardless of opts.Mode |
| Wails: ListGalleryImages | `app.go:3093` | ✅ |
| Wails: DeleteGalleryImage | `app.go:3101` | ✅ |
| Wails: GetGalleryImage | `app.go:3109` | ✅ |
| Wails: SaveImageAs | `app.go:3119` | ✅ |
| Wails: DebugImageEngine | `app.go:3156` | ✅ |
| Wails: GetImageEngineConfig | `app.go:3194` | ✅ |
| Wails: SetImageModelsDir | `app.go:3202` | ✅ |
| Wails: SetSDServerPath | `app.go:3225` | ✅ |
| Wails: OpenImageModelsDirPicker | `app.go:1795` | ✅ |
| Wails: OpenSDServerFilePicker | `app.go` | ✅ |

### Phase 2 — Frontend: txt2img UI + Gallery ✅ COMPLETE

| Item | File | Status |
|------|------|--------|
| Workflow mode bar | `image-generator.html:107` | ✅ 5 buttons (txt2img/img2img/inpaint/outpaint/upscale) with key hints |
| 3-column layout | `image-generator.html:115` | ✅ Left sidebar (380px), center canvas (flex), right panel (260px) |
| Collapsible panels | `image-generator.html` | ✅ Model, Prompt, Parameters, Advanced, Presets, Engine Paths |
| Model panel | `image-generator.html:118` | ✅ Dropdown, Load/Unload/Reload buttons, status indicator |
| Prompt panel | `image-generator.html:135` | ✅ Prompt + negative prompt textareas, weight tools ( ), +, - |
| Parameters panel | `image-generator.html:151` | ✅ Width/Height + aspect buttons, sampler, scheduler, steps slider, CFG slider, seed row |
| Advanced panel | `image-generator.html:225` | ✅ Batch count, clip skip, tiling checkbox, placeholders for LoRA/VAE/ControlNet/Upscaler |
| Presets panel | `image-generator.html:253` | ✅ 7 built-in presets + save custom |
| Engine Paths panel | `image-generator.html:276` | ✅ Models dir + sd-server binary with browse buttons |
| Auto-sizing canvas | `image-generator.html:305` | ✅ ig-canvas-wrapper with scroll, placeholder, metadata, action bar |
| Action bar | `image-generator.html:316` | ✅ Save As, Metadata toggle, Delete |
| Generate button | `image-generator.html:322` | ✅ Footer button with icon, text, Ctrl+Enter hint |
| Right panel: Model Info | `image-generator.html:333` | ✅ Status, Type, Family, Format, Backend |
| Right panel: Diagnostics | `image-generator.html:343` | ✅ Load Time, Last Gen, Errors with color coding |
| Right panel: History | `image-generator.html:351` | ✅ Scrollable thumbnail list with delete, click-to-view |
| Panel styles | `css/image-generator-panels.css` | ✅ All panel, info, history, workflow styles |
| Base styles | `css/image-generator.css` | ✅ Inputs, canvas, actions, responsive |
| JS: Core | `js/image-generator-core.js` | ✅ Init, paths, model loading, generation, gallery helpers, actions |
| JS: UI | `js/image-generator-ui.js` | ✅ Events, state persistence, presets, seed tools, shortcuts, prompt weighting |
| JS: Main | `js/image-generator.js` | ✅ Orchestrator, history panel, diagnostics, model info, metadata viewer |
| Keyboard shortcuts | `js/image-generator-ui.js:312` | ✅ Ctrl+Enter (gen), Ctrl+S (save), Ctrl+L/U/R (load/unload/reload), Ctrl+1-5 (modes), R (random seed) |
| Auto-save | `js/image-generator-ui.js:9` | ✅ localStorage key `howl-image-gen-state` |

### Phase 3 — img2img Support ❌ NOT STARTED

| Item | Status | Notes |
|------|--------|-------|
| Backend: client.img2img | ❌ STUB | `client.go:118` returns `fmt.Errorf("img2img not yet implemented")` |
| Backend: engine.GenerateImg2Img | ⚠️ DELEGATES | Calls client.img2img (stub) |
| Backend: app.go mode routing | ❌ MISSING | `GenerateImage` always calls `GenerateTxt2Img` regardless of `opts.Mode` |
| Frontend: img2img mode tab | ⚠️ DECORATIVE | Button exists but only calls `setStatus('img2img mode selected — coming soon')` |
| Frontend: image upload/drop zone | ❌ MISSING | No file input, no drag-and-drop handler |
| Frontend: strength slider | ❌ MISSING | No UI element; ImageOptions.Strength field exists but unused |
| Frontend: init image path handling | ❌ MISSING | No temp file upload to backend |

### Phase 4 — Inpainting Support ❌ NOT STARTED

| Item | Status | Notes |
|------|--------|-------|
| Backend: client.inpaint | ❌ STUB | `client.go:122` returns `fmt.Errorf("inpaint not yet implemented")` |
| Backend: engine.GenerateInpaint | ⚠️ DELEGATES | Calls client.inpaint (stub) |
| Frontend: inpaint mode tab | ⚠️ DECORATIVE | Button exists but only shows "coming soon" |
| Frontend: base image + mask canvas | ❌ MISSING | No drawing canvas, no brush tools |
| Frontend: mask export as PNG | ❌ MISSING | No canvas-to-blob logic |
| Frontend: brush size / erase / clear | ❌ MISSING | No controls |

### Phase 5 — Job System, Progress, and UX Polish ❌ NOT STARTED

| Item | Status | Notes |
|------|--------|-------|
| ImageJob struct | ✅ EXISTS | Defined in `types.go:78` but never instantiated |
| Async worker goroutine | ❌ MISSING | No job queue, no worker pool |
| Wails: StartImageJob | ❌ MISSING | No binding |
| Wails: GetImageJobStatus | ❌ MISSING | No binding |
| Wails: CancelImageJob | ❌ MISSING | No binding |
| Frontend: progress bar | ❌ MISSING | No polling for generation progress |
| Frontend: Cancel button | ❌ MISSING | No job cancellation UI |
| Frontend: spinner only | ⚠️ PARTIAL | Simple CSS spinner during generation; no progress percentage |

### Phase 6 — Integration with Worlds, Scenarios, Characters, Chat ❌ NOT STARTED

| Item | Status | Notes |
|------|--------|-------|
| Wails: SetAsCharacterPortrait | ❌ MISSING | No binding |
| Wails: AttachToWorld | ❌ MISSING | No binding |
| Wails: AttachToScenario | ❌ MISSING | No binding |
| Wails: SendToChat | ❌ MISSING | No binding |
| Frontend: "Send to Chat" action | ❌ MISSING | Removed during restructure; old gallery had placeholder stub |
| Frontend: "Set as Character Portrait" | ❌ MISSING | No UI, no binding |
| Frontend: "Attach to World/Scenario" | ❌ MISSING | No UI, no binding |
| Backend: copy to images/characters/ | ❌ MISSING | No filesystem copy logic |
| Backend: copy to images/worlds/ | ❌ MISSING | No filesystem copy logic |
| Backend: copy to images/scenarios/ | ❌ MISSING | No filesystem copy logic |
| Backend: insert into chat | ❌ MISSING | No chat message injection with image |

### Phase 7 — Advanced Polish ⚠️ PARTIAL

| Item | Status | Notes |
|------|--------|-------|
| Per-model presets (FLUX low steps) | ✅ DONE | `applyModelPreset()` in `image-generator-core.js:226` handles flux/sd15/sdxl/sd3 |
| Style presets | ✅ DONE | 7 presets in `image-generator-ui.js:205` + save custom |
| Batch generation | ⚠️ UI ONLY | `ig-batch` input exists (1-8), sent as `batch_count` in opts, mapped to `BatchSize` in txt2imgRequest. Needs verification if sd-server supports batch_size > 1 |
| Highres fix | ⚠️ UI ONLY | `ig-hires-enable` checkbox referenced in JS but element doesn't exist in HTML. Backend `txt2imgRequest` has no `enable_hr` field. |
| LoRA discovery and loading | ⚠️ PARTIAL | `ListImageLoRAs` binding works. Frontend `ig-lora` select removed from HTML. `loadLoRAs()` has null guard. Needs re-add to UI. |
| VAE Selector | ❌ MISSING | Placeholder in Advanced panel only |
| ControlNet | ❌ MISSING | Placeholder in Advanced panel only |
| Upscaler | ❌ MISSING | Placeholder in Advanced panel only |

---

## Task Breakdown — Execution Order

### Section A: Critical Fixes (Do First)

**A1. Fix `GenerateImage` mode routing in `app.go`**
- `app.go:3077` currently always calls `a.imageEngine.GenerateTxt2Img(opts)`
- Must switch on `opts.Mode`: `"txt2img"` → `GenerateTxt2Img`, `"img2img"` → `GenerateImg2Img`, `"inpaint"` → `GenerateInpaint`
- **File:** `app.go`  
- **Lines:** ~3077  
- **Est. lines changed:** 10

**A2. Fix `UnloadImageModel` missing Wails binding**
- Frontend `image-generator-ui.js:179` calls `window.go.main.App.UnloadImageModel()`
- No binding exists in `app.go`
- Must add `UnloadImageModel()` that calls `a.imageEngine.Shutdown()` or kills sd-server process
- **File:** `app.go`  
- **Est. lines added:** 15

**A3. Fix `client.go:txt2img` to pass advanced opts**
- `txt2imgRequest` struct is missing `enable_hr`, `denoising_strength`, and other fields the frontend sends
- Map `opts.EnableHR`, `opts.DenoisingStrength`, etc. into the JSON payload
- **File:** `internal/image/engine/sdcpp/client.go`  
- **Lines:** ~40-75  
- **Est. lines changed:** 15

---

### Section B: img2img Implementation (Phase 3)

**B1. Backend: Implement `client.img2img()`**
- Build HTTP POST to `/sdapi/v1/img2img` with `init_images`, `prompt`, `negative_prompt`, `strength`, etc.
- Requires base64 init image or file path handling
- **File:** `internal/image/engine/sdcpp/client.go`  
- **Est. lines added:** 60

**B2. Backend: Add init image temp file handling**
- `ImageOptions.InitImagePath` needs to be validated
- May need base64-to-temp-file conversion if frontend sends data URI
- **File:** `internal/image/engine/sdcpp/engine.go` or new helper  
- **Est. lines added:** 30

**B3. Frontend: Add img2img upload zone to HTML**
- Replace canvas placeholder with file drop zone when img2img mode is active
- Add Strength slider (0.0–1.0)
- **File:** `frontend/dist/image-generator.html`  
- **Est. lines added:** 40

**B4. Frontend: Add img2img mode switching logic**
- `setWorkflowMode()` should toggle visibility of txt2img vs img2img controls
- Upload handler: read file → base64 → send to backend temp endpoint
- **File:** `frontend/dist/js/image-generator-ui.js`  
- **Est. lines added:** 80

**B5. Frontend: Update `onGenerate()` for img2img**
- When mode is img2img, include `init_image_path` and `strength` in opts
- **File:** `frontend/dist/js/image-generator-core.js`  
- **Est. lines changed:** 20

**B6. Build + test**

---

### Section C: Inpainting Implementation (Phase 4)

**C1. Backend: Implement `client.inpaint()`**
- HTTP POST to `/sdapi/v1/img2img` with `init_images`, `mask`, `prompt`, etc.
- stable-diffusion.cpp may use img2img endpoint with mask for inpainting
- **File:** `internal/image/engine/sdcpp/client.go`  
- **Est. lines added:** 50

**C2. Frontend: Add mask drawing canvas**
- HTML5 canvas overlay on uploaded image
- Brush size slider, erase mode, clear mask button
- Export mask as PNG/dataURI
- **File:** `frontend/dist/image-generator.html`, new JS file  
- **Est. lines added:** 150 (may need new file `image-generator-canvas.js`)

**C3. Frontend: Add inpaint mode switching**
- Show/hide canvas + mask tools when inpaint mode active
- **File:** `frontend/dist/js/image-generator-ui.js`  
- **Est. lines added:** 40

**C4. Build + test**

---

### Section D: Job System (Phase 5)

**D1. Backend: Create `internal/image/job/` package**
- `JobQueue` struct with channels
- `StartJob(opts) → jobID`
- `GetJobStatus(jobID) → status`
- `CancelJob(jobID)` kills sd-server subprocess
- Worker goroutine that calls engine methods
- **File:** `internal/image/job/job.go` (new)  
- **Est. lines:** 200

**D2. Backend: Add job Wails bindings**
- `StartImageJob`, `GetImageJobStatus`, `CancelImageJob`
- **File:** `app.go`  
- **Est. lines added:** 40

**D3. Frontend: Poll for job progress**
- Replace synchronous `GenerateImage` call with `StartImageJob` + polling loop
- Progress bar or percentage display
- Cancel button that calls `CancelImageJob`
- **File:** `frontend/dist/js/image-generator-core.js`  
- **Est. lines changed:** 60

**D4. Build + test**

---

### Section E: World/Scenario/Character Integration (Phase 6)

**E1. Backend: Add entity image bindings**
- `SetAsCharacterPortrait(imagePath, characterID)` → copy to `images/characters/<id>/`
- `AttachImageToWorld(imagePath, worldID)` → copy to `images/worlds/<id>/`
- `AttachImageToScenario(imagePath, scenarioID)` → copy to `images/scenarios/<id>/`
- `SendImageToChat(imageData)` → reuse vision pipeline
- **File:** `app.go`  
- **Est. lines added:** 120

**E2. Frontend: Add action buttons to image result**
- "Set as Character Portrait" → dropdown of characters
- "Attach to World" → dropdown of worlds
- "Attach to Scenario" → dropdown of scenarios
- "Send to Chat" → inserts into current chat
- **File:** `frontend/dist/image-generator.html`, `image-generator.js`  
- **Est. lines added:** 80

**E3. Build + test**

---

### Section F: Advanced Polish (Phase 7)

**F1. Re-add LoRA selector to HTML**
- Add `<select id="ig-lora" multiple>` to Advanced panel
- Wire `loadLoRAs()` (already implemented, null guard removed)
- **File:** `frontend/dist/image-generator.html`  
- **Est. lines added:** 10

**F2. Verify batch generation end-to-end**
- Confirm sd-server `/sdapi/v1/txt2img` supports `batch_size > 1`
- If yes, `GenerateImage` should return multiple ImageResults
- If no, handle gracefully or generate sequentially
- **File:** `internal/image/engine/sdcpp/client.go`  
- **Est. lines changed:** 20

**F3. Add highres fix wiring**
- Add `ig-hires-enable` checkbox and `ig-hires-scale` slider to HTML (or verify they exist)
- Pass `enable_hr`, `hr_scale`, `hr_upscaler` to backend
- Backend: add fields to `txt2imgRequest`
- **File:** `frontend/dist/image-generator.html`, `client.go`, `core.js`  
- **Est. lines changed:** 30

**F4. Add outpaint placeholder handling**
- Similar to inpaint but with canvas expansion logic
- May require different sd-server arguments
- **Status:** Deferred — outpaint is least common use case

**F5. Add upscale placeholder handling**
- May use sd-server's upscaler endpoint or external upscaler
- **Status:** Deferred

**F6. Build + test**

---

## File Size Tracking (500-line limit compliance)

| File | Current Lines | After Section A | After Section B | After Section C | After Section D | Risk |
|------|---------------|-----------------|-----------------|-----------------|-----------------|------|
| `app.go` | 3246 | +25 | +25 | +25 | +65 | N/A (already large) |
| `client.go` | 145 | +15 | +75 | +125 | +125 | LOW |
| `engine.go` | 244 | +0 | +30 | +30 | +30 | LOW |
| `image-generator.html` | 370 | +0 | +40 | +190 | +190 | MED — may need splitting |
| `image-generator-core.js` | 387 | +20 | +40 | +40 | +100 | MED |
| `image-generator-ui.js` | 337 | +15 | +95 | +135 | +135 | MED |
| `image-generator.js` | 170 | +0 | +0 | +0 | +0 | LOW |
| `job.go` (new) | 0 | 0 | 0 | 0 | +200 | N/A |
| `image-generator-canvas.js` (new) | 0 | 0 | 0 | +150 | +150 | N/A |

**Mitigation:** Any JS file approaching 500 lines should be split. Consider:
- `image-generator-canvas.js` for mask drawing (Section C)
- `image-generator-job.js` for job polling (Section D)
- `image-generator-integrations.js` for world/scenario/character actions (Section E)

---

## Recommended Execution Order

1. **Section A** (Critical Fixes) — Unblocks everything else. ~2 hours.
2. **Section B** (img2img) — High user value, builds on Section A. ~4 hours.
3. **Section F1-F3** (LoRA, batch, highres) — Quick wins, improves txt2img experience. ~2 hours.
4. **Section D** (Job System) — Enables progress bars and cancellation. ~6 hours.
5. **Section C** (Inpainting) — Complex canvas work. ~8 hours.
6. **Section E** (Integrations) — Depends on chat/character UI state. ~6 hours.
7. **Section F4-F5** (Outpaint, Upscale) — Lowest priority.

Total estimated effort: **~28 hours** of focused development.
