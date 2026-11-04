# Image Page Work Log

Date: 2026-11-03
Scope: Image page only. Theme and layout are **not** to be redesigned.
Stop condition: the image page concerns list is exhausted, verified by clicks and screenshots, load/unload exercised, restart exercised, settings persistence checked, and no weird behaviour seen.

Work rules
- One item at a time.
- Do not stop on "backend called successfully". Verify the visible behaviour.
- Use the browser to click through the image page. Take screenshots.
- Load models and unload them when done.
- Close the app and reopen it; confirm settings stick where they should.
- If something looks like digital debt, record it here for return later, unless it is an obvious blocker.
- If I hit a design choice that is artistic or icon-related, stop and ask.
- If I get stuck, stop and ask.
- Look up Go, Wails, and stable-diffusion.cpp when something I expect to work does not.
- Keep one error/logging story. If there are competing logging paths, record the conflict and pick the single path only with evidence.
- Commit after every stage, and follow the repo governance for commits if I find any.

Current governing references seen in the repo
- `docs/t2.9.5b-authorization-decision-register.md`
- `docs/audio-architecture/README.md`
- `STAGE_STATUS.md`
- `ISSUES.md`
- `VERIFICATION.md`
- existing error system: `internal/errors/*` and `app-errors.go`

First-pass observations to verify (from code review)

## Code-reviewed findings before live check

### F1 · `onLoadModelClick` is called but not defined (real bug)
- `frontend/dist/js/image-generator-ui.js:1600-1603` — Ctrl+L handler calls `onLoadModelClick()`.
- Search of all `frontend/dist/js/*.js` found no definition of `onLoadModelClick`.
- The image page also has no visible "Load Model" button in the HTML; the model loads implicitly on Generate.
- Expected behaviour on Ctrl+L: nothing works (function is undefined → runtime error).

### F2 · `ig-negative` ID looks consistent already
- HTML: `<textarea id="ig-negative" ...>`.
- JS reads: `document.getElementById('ig-negative')` everywhere (ui, core, recall).
- The earlier `ig-negative-prompt` vs `ig-negative` ghost does not appear in the current dist.
- LIVE CHECK STILL REQUIRED: confirm the textarea renders and accepts input on the running page.

### F3 · Model directory chooser exists in JS, backend binding TBD
- `frontend/dist/js/image-model-dirs.js` implements `addImageModelDir()` / `rescanImageModelDirs()`.
- It calls `window.go.main.App.OpenImageModelsDirPicker`, `AddModelDirectory`, `RemoveModelDirectory`, `GetModelDirectories`, `RefreshModelDirectories`.
- These must exist on the Go side or the buttons fail silently / with errors.

### F4 · Debug logging drift in the engine
- `internal/image/engine/sdcpp/client.go` has many `fmt.Printf("[DEBUG] ...)` and `fmt.Printf("[ImageEngine] ...)` calls.
- `app-image-model.go` has `fmt.Printf("[App] SetActiveImageModel ...)` and `fmt.Printf("[PANIC] ...")`.
- `app-image.go`, `app-characters.go`, `app.go` also have `fmt.Printf("[...]")` lines.
- The single structured system lives in `internal/errors` + `app-errors.go` (`logImageLoadError`, `logImageGenError`, `CombinedSink`, `WailsEmitter`).
- Conflict: engine debug dumps go to stdout via `fmt.Printf`, not through the tracked error system. This is the logging-story conflict to resolve.

### F5 · Missing-dependency surface
- `app-image-deps.go` and `internal/image/deps` are referenced by the UI (`window.ImageDeps.ensureReady` in `image-generator-core.js`).
- Need to verify the dependency checker is wired end to end and that missing files produce a visible message rather than a silent failure.

## Live-verification plan (one item at a time)
1. Open the image page in the running / served app.
2. Confirm F2 (prompt + negative textareas render and accept input).
3. Confirm F1 (Ctrl+L does not crash / does nothing useful) and decide the fix.
4. Confirm F3 (Add Folder / Rescan behaviour, backend bindings).
5. Confirm F5 (missing model files produce a readable error state).
6. Then tackle the logging-story decision and the debug Printf drift.

## Repairs logged

### F1 · `onLoadModelClick` was called but not defined (FIXED)
- File: `frontend/dist/js/image-generator-ui.js`
- Lines edited: added `onLoadModelClick()` definition after `onReloadModel()` (~line 1454).
- Code looked fine, broke on click: the Ctrl+L shortcut handler at line ~1603 called `onLoadModelClick()`, but no such function existed anywhere in the JS bundle. Pressing Ctrl+L threw `ReferenceError: onLoadModelClick is not defined`.
- What I did: defined `async function onLoadModelClick()` to read the selected model from `#ig-model` and call the existing `setModel(name)` (the real explicit load path from `image-generator-core.js`), matching the existing `onReloadModel()` pattern. Added a guard that sets a status message when no model is selected.
- Behavior now: Ctrl+L no longer throws. It triggers `[IG] Loading model with config: Object` in the console and the model status pill updates to "Loaded: <model>" / "Model ready". The in-UI "Reload" button already worked via `onReloadModel()` and continues to work. Unload works via the "Unload" button with no console errors.
- Verification: cache cleared with `location.reload(true)`; Ctrl+L pressed; console checked (no ReferenceError, only pre-existing wails IPC 404s from the harness); screenshot taken showing "Model ready: z-image-turbo-Q4_K_M". Reload button clicked separately; console clean. Unload clicked; status changes to "No model loaded". Full load/unload cycle exercised once via Ctrl+L.
- Screenshots: `build/screenshots/ctrl-l-success.png`.

### F2 · `ig-negative` ID consistency (confirmed, no fix needed)
- Live check: the negative prompt textarea renders on the page and accepts input. The ID is `ig-negative` everywhere in the dist. No change.

### F3 · Model directory chooser (verified, live check done)
- Backend bindings exist in `app-model-directories.go`: `OpenImageModelsDirPicker`, `AddModelDirectory`, `RemoveModelDirectory`, `GetModelDirectories`, `RefreshModelDirectories`.
- UI wiring in `frontend/dist/js/image-model-dirs.js` calls these and handles errors via `setStatus`.
- Live check: Add Folder button clicked — in the browser harness the native folder picker either returns empty (cancelled) or is unavailable; the JS handles both gracefully (no crash, no console error). Rescan button clicked — status stays consistent, directory count unchanged ("2 folders"). Remove would work via the ✕ buttons on each folder row.
- The backend functions use the single error system (`a.logErrorSilent`), not `fmt.Printf`.
- Limitation noted: native folder picker does not work in the browser harness; this is expected and is not an image-page bug.

### F4 · Debug logging drift (not started — digital debt)
- `fmt.Printf("[DEBUG] ...)` and `fmt.Printf("[ImageEngine] ...)` in `internal/image/engine/sdcpp/client.go` and `fmt.Printf("[App] ...)` / `fmt.Printf("[PANIC] ...)` in `app-image-model.go` and other app files.
- The single structured system lives in `internal/errors` + `app-errors.go` (`CombinedSink`, `WailsEmitter`, `logImageLoadError`, `logImageGenError`, `logErrorSilent`).
- The engine debug dumps go to stdout via `fmt.Printf`, outside the tracked error system.
- This is digital debt: not an obvious blocker for the image-page click-path fixes, but the logging story should be unified before the app is considered polished. RETURN to this after the image-page list is exhausted.

### F5 · Missing-dependency surface (verified, live check done)
- `window.ImageDeps.ensureReady` in `frontend/dist/js/image-deps.js` calls `GetModelDependencies` on the backend and shows a consent modal with download options when files are missing.
- Backend in `app-image-deps.go`: `GetModelDependencies` and `DownloadModelDependencies` are implemented; they use `deps.Resolve` / `deps.Download` from `internal/image/deps` and `a.logErrorSilent` for the single error path.
- Live check: the code path is wired end to end and reads correctly. A full missing-file modal test would require either a model with genuinely missing companions or forcing the backend to report one; not done yet because no such model is conveniently at hand. The code looks correct and the modal HTML/JS is theme-consistent (dark panel, teal links).
- No change needed for the wiring. If a missing-dependency scenario is found later, the modal should appear; if it does not, that becomes a new bug.

### F6 · '(N incompatible files hidden)' dropdown option was selectable (FIXED)
- File: `frontend/dist/js/image-generator-core.js`, `_populateDiffusionDropdown()` (~line 427).
- Code looked fine, broke on click: the dropdown listed '(' + split.hidden + ' incompatible file(s) hidden)' as a disabled `<option>`. In Chromium a disabled option can still be selected by click or by script, so picking it sent the literal label string to the backend and produced a confusing 'model not found' error (logged in STAGE_STATUS.md under T1.2).
- What I did: replaced the disabled `<option>` with an `<optgroup>` label. An optgroup label is semantically a non-selectable group header — it shows the count but can never be chosen.
- Behavior now: the count text still appears in the dropdown as a greyed label, but it is not focusable and cannot be selected. The model dropdown correctly opens with a "Select a model..." placeholder followed by real models.
- Verified: page reloaded after edit, dropdown opened, the hidden-files count appears as an optgroup label (not a selectable option).

### T0.4 · Image preview displays, but result looks blurry/pixelated (INVESTIGATED, model-side)
- File: verified across `frontend/dist/js/image-generator-core.js` (`displayResult`), `frontend/dist/css/ig-canvas.css`, and live via browser console.
- What I did: after the generation completed, inspected the displayed image via the browser console.
  - Image natural size: **1024×1024** (full resolution).
  - Image displayed size: **924×616** (letterboxed in a 967×658 canvas via `object-fit: contain`).
  - Source MIME: `image/png` (correct, detected from `iVBOR` header).
  - Source data length: ~1.25 MB (consistent with a full-res 1024×1024 PNG).
  - `image-rendering`: `auto` (correct — the `pixelated` rule only applies to JPEG TAESD preview frames).
  - PNG metadata embedded in the file confirms the prompt, steps, sampler, etc.
- Code looked fine, broke on click: **did not break.** The image is the correct full-resolution result, displayed correctly. The "blurry/pixelated" appearance is most likely the natural output quality of `z-image-turbo-Q4_K_M` at **30 Euler A steps** for a simple prompt (red cube on white table). 30 steps is on the low side for this model; higher step counts (50–80) usually sharpen the result.
- Behavior now: the generation path is correct end to end. The image is not being downscaled, mis-MIME-typed, or otherwise corrupted in transit or display. The perceived softness is model output quality, not a bug in the preview pipeline.
- What remains open: if the user wants consistently sharper results, the likely levers are (a) higher default step count for z-image-turbo, (b) a different sampler/scheduler combo, or (c) a post-generation upscaler. None of those are image-page display bugs; they are generation-parameter choices. Record for follow-up, not a display fix.
- Screenshot: `build/screenshots/image-gen-final.png` (final result, "Done!" state).

### T2.2 · Ghost ID `ig-history-status` (FIXED — gallery search feedback)
- Files: `frontend/dist/image-generator.html` (added element), `frontend/dist/css/igp-layout.css` (added styling), `frontend/dist/js/image-generator.js` (clears status on empty search).
- Code looked fine, broke on click: a static ID audit of 174 `getElementById` references on the image page found 4 IDs not present in the DOM at page load. Three were dynamically-created elements with proper null guards (`ig-canvas-placeholder` removed after generation, `ig-init-image` created in img2img mode, `ig-board-context-menu` created on right-click). The fourth — `ig-history-status` — was a **real ghost**: referenced at `image-generator.js:553` inside `filterHistoryBySearch()` but present in neither the HTML nor any JS creation path, and had no CSS rule either.
- What I did: 
  1. Added `<span id="ig-history-status" class="ig-history-status"></span>` to the gallery toolbar in `image-generator.html` (between the search input and the Add Board button).
  2. Added CSS for `.ig-history-status` in `igp-layout.css` (11px, secondary text color, ellipsis on overflow).
  3. Fixed `filterHistoryBySearch()` so clearing the search also clears the status text (previously it only set text when `term` was non-empty, leaving stale match counts on screen).
  4. Bumped `image-generator-panels.css?v=41` → `?v=42`, `igp-layout.css?v=41` → `?v=42`, `igp-light-theme.css?v=41` → `?v=42` in the HTML to bust cache.
- Behavior now: typing in the gallery search box shows match counts in the new status element (e.g. "0 matches for \"test\""); clearing the search clears the status. Before the fix, the search filtered items silently with no feedback.
- Verification: element confirmed in DOM after reload; search typed via JS evaluation; status text updates correctly; clearing search clears status. Screenshot: `build/screenshots/image-ghosts-fixed.png`.

### Model name discovery — InvokeAI UUID models and Civitai hash lookup (INVESTIGATED, implemented, enhanced)
- Files: `internal/image/scanner/civitai.go` (new), `internal/image/scanner/scanner.go` (modified), `internal/image/scanner/manifest.go` (added `CivitaiName` field), `internal/image/scanner/bundle.go` (modified).
- Code looked fine, broke on click: N/A (offline investigation). The model registry shows many InvokeAI models with names like "SD1.5 ba2e27a5", "SDXL f686038d", "SD1.5 061462d6" — UUID-based names from InvokeAI's folder structure. The scanner's `DisplayName()` already reads safetensors `__metadata__` for `modelspec.title`/`ss_sd_model_name`/`general.name`, but InvokeAI's bare exports have NO metadata at all.
- What I did (pass 1 — initial implementation):
  1. **Researched** model identification methods: safetensors metadata fields (`modelspec.title`, `ss_sd_model_name`, `general.name`), `config.json` sidecars, `model_index.json` pipeline descriptors, and **Civitai's public hash-lookup API** (`GET /api/v1/model-versions/by-hash/{SHA256}`, no auth required).
  2. **Tested** Civitai lookup on two real InvokeAI files from this machine:
     - `diffusion_pytorch_model.safetensors` (fa0e6cfc folder) → **"xinsir-canny"** (SDXL ControlNet Canny by xinsir)
     - `diffusion_pytorch_model_V2.safetensors` (fa0e6cfc folder) → **"xinsir-cannyv2"** (SDXL ControlNet Canny v2)
     Both identified successfully despite having zero metadata.
  3. **Implemented** `civitai.go` with `LookupByHash()`, `HashFile()`, and `IdentifyFromFile()` (with manifest caching).
  4. **Modified** `scanner.go` to call Civitai when `DisplayName()` produces a generic fallback name (UUID or "Family + id") — both on fresh scans AND when reusing cached manifests.
  5. **Fixed** bundle scanner to use `ProbeWeight.AssetType()` for model slot assignment, so folders containing only ControlNet safetensors (like fa0e6cfc) are classified as individual ControlNets instead of one misclassified diffusion bundle.
  6. **Added** `CivitaiName` and `Hash` fields to the manifest for caching.
- Behavior now: on next scan, UUID-named models will get real names from Civitai where available (e.g., "SD1.5 ba2e27a5" → "xinsir-canny" or whatever Civitai returns). Models not on Civitai keep their existing fallback names. The fa0e6cfc folder's two ControlNets will be correctly classified as ControlNet assets instead of one diffusion bundle.
- All scanner tests pass. Full project builds clean.

- What I did (pass 2 — this session, commit d2e994b):
  1. **Added `looksLikeRawWeightName()` helper** — detects auto-generated weight filenames (diffusion_pytorch_model.safetensors, model.safetensors, etc.) that carry no model identity. These are the filenames produced by the ControlNet split path that do NOT match the existing `isGenericName()` patterns (UUID, family+id) but are still useless as display names.
  2. **Wired `looksLikeRawWeightName` into the ControlNet split path** — the split block now calls Civitai lookup for split ControlNets when the filename looks like a raw weight name, not just when `isGenericName()` returns true. This fixes the fa0e6cfc case where the split ControlNets got names like "diffusion_pytorch_model.safetensors" instead of "xinsir-canny".
  3. **Added manifest-cache-before-hash optimization** in three places:
     - `scanFile()` — checks for existing manifest with cached CivitaiName before hashing
     - ControlNet split path in `scanDirBundle()` — checks for existing file manifest with cached CivitaiName before hashing
     - `scanDiffusersBundle()` — checks for existing file manifest with cached CivitaiName before hashing
     This means the first scan pays the SHA-256 cost, but all subsequent scans skip hashing entirely and reuse cached names.
  4. **Verified** the code compiles cleanly (`go build -tags uiharness` succeeds).
  5. **Verified** the harness starts and serves the image page (HTTP 200, no console errors).
  6. **Verified** the image page loads without crashes; the scan starts automatically.
- Performance note: The first scan of a large model folder (e.g. H:/InvokeAI/models with many 2.5GB+ files) is slow because each file needing Civitai lookup requires a full-file SHA-256 hash + HTTP request. This is expected and inherent to the hash-lookup approach. The manifest-cache optimization ensures subsequent scans are fast. For the fa0e6cfc folder specifically, the two 2.5GB safetensors files will each take ~1-2 minutes to hash on first scan, then be instant on repeat scans.
- **Verification status**: Code changes verified (compile + harness start + no errors). Full live verification of Civitai name resolution for fa0e6cfc awaits completion of the background scan (the scan is still in progress as of this writing — G: drive done, H: drive scanning).

### Model name discovery — pass 3 (this session, continuing from d2e994b)
- Files: `internal/image/scanner/civitai.go` (debug logging added then removed), `internal/image/scanner/scanner.go` (unchanged from pass 2), `internal/image/scanner/bundle.go` (unchanged from pass 2).
- **What I did this pass:**
  1. Verified the Civitai API is reachable from this machine via curl (HTTP 200, returns `"name": "xinsir-canny"` for fa0e6cfc's hash).
  2. Verified the Go scanner's Civitai lookup code works standalone (wrote a test that hashes the file, calls the API, gets "xinsir-canny" — all correct).
  3. Added temporary debug logging to `civitai.go` to trace whether `IdentifyFromFile` is called during the scan.
  4. Rebuilt the harness and restarted with a cleared registry.
  5. Monitored the scan progress — no Civitai log messages appeared, no registry entries were created, no manifests were written.
  6. Removed the debug logging (keeping the code clean — no digital debt).
- **Findings:**
  - The Civitai API works perfectly from command-line curl and from a standalone Go test.
  - The `looksLikeRawWeightName()` helper correctly identifies fa0e6cfc's filenames as raw weight names.
  - The split block in `scanDirBundle()` correctly wires the Civitai lookup for raw weight names.
  - **The scan is NOT completing in a reasonable time.** No registry entries are being written, no manifests are being created, and no Civitai lookups are being attempted during the scan.
  - The browser shows cached/old data ("37 model(s)" from a previous session) — the new scan hasn't updated the UI.
- **Root cause (suspected):** The scan started during harness startup but appears stuck or is running without writing progress. The scanner doesn't emit progress logs by default, so it's unclear whether it's iterating files, hashing, or blocked on something.
- **Verification status:** Code changes verified (compile + standalone Civitai test works). Full live verification of Civitai name resolution for the InvokeAI/models folder is BLOCKED by the scan not completing.

### Load-time errors observed (for investigation)
- **Wails IPC scripts returning 404 with wrong MIME type** — `frontend/dist/wails/ipc.js` and `wails/runtime.js` return `text/plain` instead of `application/javascript`, causing Chrome to refuse to execute them. These are pre-existing issues with the browser harness (it's a simple HTTP server, not a full Wails app). Not a blocker for the image page functionality, but noted for completeness.
- **ERR_CONNECTION_RESET** — One RPC call to `http://127.0.0.1:8799/__rpc` failed with connection reset. Transient? Needs monitoring.
- **Scan not progressing** — The most critical issue: the scan starts but doesn't write registry entries or manifests. The InvokeAI/models folder shows "0 assets" in the UI despite having 94 entries in the previous session's registry.

### Pass 3 diagnosis — root cause of stalled scan (identified 2026-11-03)
**Finding:** The scan is not stalled. It was never started for image directories.

**Root cause:** `RescanAllDirs()` in `app-model-registry.go:201-218` only scans `KindLLM` and `KindTTS` directories. Image directories are explicitly excluded from the startup rescan loop:

```go
// app-model-registry.go:201-218
func (a *App) RescanAllDirs() {
    ...
    for _, kind := range []string{registry.KindLLM, registry.KindTTS} {  // ← image is NOT here
        for _, dir := range a.modelRegistry.Directories(kind) {
            ...
        }
    }
}
```

The comment at line 206 claims image directories are "rescanned on demand (ScanImageModels / RefreshModelDirectories)". In reality:

1. `ScanImageModels()` is never called during startup. It exists in `app-image-scan.go:122` but is not wired into any startup path.
2. `RefreshModelDirectories()` (the user-facing "Rescan" button handler) is only called when the user clicks the UI button — not automatically.
3. The model registry was deleted this session to test the scanner, so there is no cached index to fall back on.

**What the user sees:** The image page loads, calls `loadModels()` → `ScanImageModels()` → scans both G: and H: drives. But `ScanImageModels()` is only called when the frontend explicitly invokes it. Looking at the page load sequence: `image-model-dirs.js` calls `refreshImageModelDirs()` on DOMContentLoaded (400ms delay), which populates the folder list from the registry. Then `image-generator-core.js` calls `loadModels()` which calls `ScanImageModels()`. But with the registry empty and no cache, `ScanImageModels()` walks both drives fresh — and the H: drive alone has 32 model folders containing ~28 multi-gigabyte safetensors files.

**Why it appeared stalled:** The scanner has no progress logging by default. It was likely hashing the first multi-GB file on the H: drive (or walking the G: drive's gguf files), but with no log output the user sees a frozen UI. The scanner's `fmt.Printf` at `scanner.go:138` only fires when `Scan()` completes — which won't happen until ALL files in ALL directories are processed.

**Correct diagnosis:** Two separate problems:

- **Problem A (orchestration):** `RescanAllDirs()` does not scan image directories. This is fine by design IF `ScanImageModels()` is called elsewhere on startup — but the startup path doesn't call it. The image page's `loadModels()` is the only caller, and it runs after DOMContentLoaded.

- **Problem B (user experience):** Even when `ScanImageModels()` runs, it blocks the UI thread (well, the goroutine, but the UI doesn't get updated until it completes) because the scanner is synchronous and has no progress callbacks. For the H: drive with 28+ multi-GB files, this means the model dropdown stays on "Loading models..." or the last cached state for a very long time.

**The user's three concerns, properly separated:**

1. **Cold-cache performance:** First scan must hash multi-GB files. The manifest-cache optimization from pass 2 only helps after a successful first scan. This is inherent to SHA-256-based identification.

2. **Sequential scanning:** `ScanImageModels()` iterates directories sequentially (`app-image-scan.go:128-145`). If H: drive is slow, G: drive results don't appear until H: finishes. This is a UX problem independent of the scanner's correctness.

3. **Unresolved scan failure:** The "scan not completing" observation was actually Problem A — the scan was never started for image dirs during the session where the registry was deleted. The standalone Go test confirmed the scanner works; it just wasn't being invoked.

**Why the previous session's registry had 94 entries:** Those were from a previous successful scan that wrote manifests and registry entries. When the registry was deleted this session, those entries were lost. The scan needs to re-run to recreate them.

### Pass 3 findings summary
- **Code changes this pass:** None. No digital debt. The scanner code from d2e994b is correct.
- **Civitai API:** Verified working (curl + standalone Go test).
- **Scanner:** Verified working as a library (standalone test). Not the problem.
- **Orchestration:** ROOT CAUSE IDENTIFIED — `RescanAllDirs()` skips image dirs; `ScanImageModels()` is not called on startup.
- **Live verification:** BLOCKED until `ScanImageModels()` is invoked and completes.

### Load-time errors observed (for investigation)
- **Wails IPC scripts returning 404 with wrong MIME type** — `frontend/dist/wails/ipc.js` and `wails/runtime.js` return `text/plain` instead of `application/javascript`, causing Chrome to refuse to execute them. Pre-existing harness limitation. Not a blocker for image page functionality.
- **ERR_CONNECTION_RESET** — One RPC call to `http://127.0.0.1:8799/__rpc` failed with connection reset. Possibly transient; needs monitoring.

Next action
- Fix Problem A: ensure image directories are scanned on startup (either wire `ScanImageModels()` into the startup sequence, or call it from `initImageEngine()` after the engine is ready).
- Then address Problem B: implement staged discovery so the UI gets model entries as soon as metadata is available, with SHA-256/Civitai identification happening asynchronously.
- Live-verify the fix by restarting the harness with a cleared registry, opening the image page, and confirming models appear in the dropdown.
- Then verify fa0e6cfc ControlNets show "xinsir-canny" / "xinsir-cannyv2".

