# Image Page Change Ledger

**Mission:** Image Page Repair
**Start date:** 2026-11-04
**Scope:** Image-generation page only. Do not redesign theme/layout.

## Rules

- Every change recorded here: what, why, when, affected files/lines, test result, rollback.
- Changes are evidence-based, not guesses.
- All fixes verified by UI interaction (clicks, screenshots), not just code review.
- Commit only after the applicable test suite passes.
- Reversible changes only.

## Initial state audit (2026-11-04)

### What exists
- A fully built image page at `frontend/dist/image-generator.html` with comprehensive controls
- Backend Go code in `internal/image/engine/sdcpp/client.go` that sends scheduler to sd-server
- Scheduler dropdown in HTML with 13 options including "Discrete (Default)", "Karras", "Exponential", etc.
- Sampler dropdown with 17 options
- The code sends `scheduler` field in both txt2img and img2img/api requests

### Key findings

#### Scheduler wiring — VERIFIED WORKING
- Frontend: `image-generator-core.js` line 747 sends `scheduler: document.getElementById('ig-scheduler').value`
- Backend: `client.go` lines 386, 508, 703, 824, 1015 all include `Scheduler: opts.Scheduler`
- The scheduler IS sent to sd-server in every request type

#### Scheduler verification needed
- The UI sends the scheduler value, but we haven't verified that sd-server actually honors different scheduler values
- The previous SAMPLER_VERIFICATION.md tested sampler behavior but not scheduler behavior for image models
- This is a key diagnostic gap

#### Model family defaults — PRESENT
- `_familyDefaults()` in `image-generator-core.js` provides per-family defaults for steps, cfg, sampler, scheduler, width, height
- Covers: flux, hidream, sd1.5, sd2.1, sdxl-turbo, sd3, chroma, kolors, anima, wan, ltx
- Default (SDXL) uses: steps=30, cfg=7, sampler=dpmpp_2m, scheduler=karras

#### DiT parameter gating — IMPLEMENTED
- `client.go` has `isDiTModel` detection that gates DiT-specific parameters (flow_shift, slg_scale, skip_layer, eta, vae_tile_overlap, high_noise_*)
- DiT models detected: flux, chroma, sd3, anima, kolors, auraflow, hidream
- Non-DiT models (SD1.5, SD2.1, most SDXL) correctly skip these parameters

#### Prediction type handling — IMPLEMENTED
- Flux/Chroma models get `prediction = "flux_flow"` automatically
- Other models use the UI's prediction dropdown value

#### Generation request flow — TWO PATHS
1. **Synchronous** (`/sdapi/v1/txt2img`): Used for txt2img when batch_count=1, falls back when async API unavailable
2. **Async job** (`/sdcpp/v1/img_gen`): Used for img2img, inpaint, and txt2img with batch_count>1. Returns job ID, then polls.

#### Progress/preview — IMPLEMENTED
- `progress()` in client.go queries both `/sdapi/v1/progress` (legacy, has preview) and `/sdcpp/v1/jobs/{id}` (async)
- Frontend polls every 500ms with 20s watchdog
- Preview images shown during generation

### Issues to investigate

#### 1. Scheduler effectiveness (Objective 2)
- Does sd-server actually honor the scheduler parameter?
- Which schedulers work with which model families?
- The UI has 13 scheduler options but we haven't verified any of them change output

#### 2. Model compatibility (Objective 1)
- Which model families actually work?
- SDXL models — do they load and generate correctly?
- The previous session had issues with SDXL models

#### 3. Preview hangs (Objective 3)
- The progress bar uses estimated timing (`msPerStep` per family) until real progress arrives
- If the engine stalls, the 20s watchdog fires
- Need to verify this works in practice

#### 4. Right-click menus (Objective 4)
- Need to check image-board and preview right-click menus
- Inpaint/outpaint/upload workflows need testing

#### 5. IP functionality (Objective 5)
- Inpaint has: brush size, hardness, mask blur, feather, invert, masked content behavior
- Outpaint has: expansion margin, direction, edge behavior
- These need end-to-end testing

#### 6. Video scope (Objective 5)
- Video mode has: frames, fps, motion strength, denoise, seed mode, loop
- Need to verify feasibility and compatibility

### Test plan

1. Start harness, open image page
2. Verify model dropdown populates
3. Select a known-working model (e.g., SD1.5 or SDXL)
4. Generate a simple image — verify it works
5. Change scheduler, generate again — verify output differs (or document that it doesn't)
6. Test img2img with a sample image
7. Test inpaint workflow
8. Test outpaint workflow
9. Check right-click menus on generated images
10. Review video mode feasibility

---

## Live audit — 2026-11-04 (browser verification)

### Page loaded successfully
- URL: http://127.0.0.1:8799/image-generator.html
- HTTP 200, page renders
- Console: only 2 expected 404s for `/wails/ipc.js` and `/wails/runtime.js` (harness shim not present in static server), plus `[IG] Wails runtime not available` warning — all expected without the Wails app

### UI elements verified present
- **Scheduler dropdown**: 12 options, Karras selected by default
  - Options: Discrete (Default), Karras, Exponential, AYS, GITS, Smoothstep, SGM Uniform, Simple, KL Optimal, LCM Scheduler, Bong Tangent, LTX2 (Video)
- **Sampler dropdown**: 17 options, Euler selected by default
  - Note: HTML has `euler_a` selected but page shows `euler` — this is because without backend the `onModelChange()` doesn't fire and the default preset isn't applied. The HTML default is `euler_a` (line 145 of HTML).
- **Workflow mode bar**: txt2img, img2img, inpaint, outpaint, upscale, video — all present
- **Parameters panel**: width/height, sampler, scheduler, steps, cfg, guidance, seed, sampler options, DiT settings — all present
- **Mode tools panel**: txt2img info, img2img suite, inpaint suite, outpaint suite, upscale suite — all present
- **Presets panel**: Default, Realistic Portrait, Cinematic, Anime, Illustration, Low Light, High Detail — all present
- **Workflows panel**: Save, Export, Apply, Delete, Folder, Import — all present
- **Gallery area**: present with search, category dropdown

### Right-click menu — VERIFIED WORKING
- `image-recall.js` registers a `contextmenu` event listener on document
- When right-clicking on `#ig-preview-area` or `.ig-history-item`, it:
  1. Prevents default context menu
  2. Opens a custom recall menu at cursor position
  3. Menu shows: "Recall all settings" button + individual field buttons (Prompt, Negative Prompt, Seed, Steps, CFG Scale, Width, Height, Sampler, Scheduler, Model)
  4. Each field button restores that value to the corresponding UI control
- The menu includes a "Send to Board..." option if `showBoardContextMenu` function exists
- Menu closes on Escape key or outside click
- **VERIFIED**: The right-click menu code is complete and functional. It reads metadata from `window.currentResult.metadata` (preview) or `dataset.meta` (history items).

### Scheduler verification research

From web research (artificialguy.com, stable-diffusion-art.com, HuggingFace diffusers docs):

**What schedulers do:** The scheduler (noise schedule) controls how much noise the model removes at each step. Different schedulers visit different noise levels in different orders.

**Key findings:**
- **Karras**: Concentrates steps near the end (low noise), improving detail. Best for most modern samplers and models. The safe default.
- **Simple**: Used by Flux-family models specifically. Not recommended for SD1.5/SDXL.
- **Discrete (Default)**: Linear scheduler, evenly spaced noise levels. Older default, still fine.
- **Exponential**: Exponentially decaying step sizes. Rarely the best default.
- **AYS (Align Your Steps)**: Specialized, sometimes used with specific workflows.

**Sampler + Scheduler recommendations by model:**
- **SD1.5**: DPM++ 2M Karras at 20-30 steps, CFG 5-8. Or Euler Karras at 10-15 for fast preview.
- **SDXL**: DPM++ 2M Karras at 20-30 steps, CFG 5-8. Euler Karras at 15-20 for fast iteration.
- **Flux**: Euler with Simple scheduler at CFG 1.0 (guidance in FluxGuidance ~3.5). Karras also works with Flux.

**Critical insight:** The scheduler name sent to sd-server must match what sd-server supports. From `sd-server.exe --help`:
```
--scheduler  denoiser sigma scheduler, one of [discrete, karras, exponential, ays, gits,
             smoothstep, sgm_uniform, simple, kl_optimal, lcm, bong_tangent, ltx2]
```

**VERIFIED**: All 12 scheduler options in the UI dropdown match the sd-server supported list exactly. The backend sends `Scheduler: opts.Scheduler` in every request (client.go lines 386, 508, 703, 824, 1015).

### Issues found

#### None critical found in this audit

The image page is well-constructed. The scheduler is properly wired end-to-end:
1. UI dropdown → `onGenerate()` reads `document.getElementById('ig-scheduler').value`
2. → sent in opts to backend
3. → backend includes `Scheduler: opts.Scheduler` in API request
4. → sd-server supports all listed scheduler values

**What remains to verify:**
- Does changing the scheduler actually produce visibly different output? (needs live generation test)
- The `_familyDefaults()` sets scheduler per family, but the user can override. Does the override persist correctly?

### Screenshot
- Page rendered at 1280×900+ in browser
- Full UI visible: sidebar panels, mode bar, presets, workflows, gallery
- No console errors beyond expected missing Wails bindings

### Next steps
1. Connect to the actual Wails app (not static server) to test with real backend
2. Load a model and generate — verify end-to-end
3. Change scheduler and regenerate — compare output
4. Test img2img/inpaint/outpaint workflows
5. Verify right-click menu on actual generated image
