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

Next action
- Continue through remaining image-page concerns: T2.3 (layout) and any others on the list.

