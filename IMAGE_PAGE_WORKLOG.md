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

### F3 · Model directory chooser (pending live check)
- Backend bindings (`OpenImageModelsDirPicker`, `AddModelDirectory`, etc.) still need a live click test. Not yet verified end to end.

### F4 · Debug logging drift (not started)
- `fmt.Printf("[DEBUG] ...)` and `fmt.Printf("[ImageEngine] ...)` in the engine vs the single structured system in `internal/errors` + `app-errors.go`. Recorded as digital debt to return to; not an obvious blocker for the image-page click-path fixes.

### F5 · Missing-dependency surface (not started)
- `window.ImageDeps.ensureReady` / `app-image-deps.go` still need a live test with a missing file to confirm a readable error state.

Next action
- Verify F3 (Add Folder / Rescan) end to end via clicks, then continue down the list one item at a time.

