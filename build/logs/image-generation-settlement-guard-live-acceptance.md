# Live acceptance: settlement guard in the real Howl Chat application

Date: 2026-11-02, 19:11–19:27 local. Application: the shipped desktop build already running
as `howl-chat.exe` PID 18884, window title "HOWL Chat", 1400x900 at screen 873,422.

## How the application was driven (and what was not used)

The app is a Wails/WebView2 desktop window, not a served page, so the browser preview tools
cannot attach to it. I drove the **running real application** through Windows UI Automation
(`UIAutomationClient`): the WebView2 accessibility tree exposes the actual DOM (buttons,
comboboxes, sliders, status text), and actions were real invocations (`InvokePattern`,
`SelectionItemPattern`) or, for elements exposing no action pattern (the left-rail nav
icon), a real mouse click at the element's own rectangle.

Captures are genuine window pixels via `PrintWindow(..., PW_RENDERFULLCONTENT)`, not
reconstructions. The first attempt used `CopyFromScreen` and produced a desktop composite of
overlapping windows; that image (`01-app-initial.png`) is kept only as the reason for
switching methods and **nothing in this record relies on it**.

Not used, per instruction: the UI harness, the synchronous bindings, direct HTTP calls, and
any mock. No process was killed by name; no service was restarted by hand. The only
process-level facts recorded are observed `tasklist`/`Get-Process` values.

Helper: `build/tmp/uia-probe.ps1` (drive + capture). No application source, test or earlier
record was modified during this acceptance.

## Workflow exercised, click by click

1. **Opened Image Studio.** Left rail icon "Image Generator" (no UIA action pattern) clicked
   at its centre; page changed from the chat view to "Image Generator" with tabs Text to
   Image / Image to Image / Inpaint / Outpaint / Upscale / Video.
   Shot: `03-image-studio.png`.
2. **Resource/model state shown by the UI.** Main model combobox empty, header "NO MODEL
   LOADED", status "Scanning models…" then "39 model(s) • 5 VAE • 14 LoRA discovered";
   two model folders (`…/AI GGUF/Image models` 16 assets, `…/InvokeAI/models` 92 assets);
   gallery populated from earlier sessions. No VRAM figure is exposed anywhere in this UI
   (searching the tree for "VRAM" returns nothing), so resource state could only be read as
   the model-selection state plus the spawned process's working set.
3. **Selected a validated model.** Opened the combobox (`ig-model`) and selected
   `z-image-turbo-Q4_K_M (z_image) [gguf]` from the list. The UI responded by switching the
   resolution to 1024×1024 and the sampler to "Euler", i.e. the family defaults applied.
   Shots: `04-model-list.png`, `05-model-open.png`, `06-model-selected.png`.
4. **Submitted a normal generation** with the prompt already in the box (Starship
   Enterprise-D …), 1024×1024. Observed: button became **ABORT**, spinner, status
   "Loading model…", and `sd-server.exe` PID 20892 appeared.
   Shot: `07-generating.png`.
5. **Loading → progress → completion.** Progress showed a live preview and "Step N of 8";
   completion switched the button back to GENERATE, showed "Generated" with the image, a
   metadata panel, and status "Done!".
   Shots: `08-after-generation.png` (result rendering), `09-completed.png` (final).
   Artifact on disk: `images/gallery/img_20261102_011843_800.png` + `.json`, whose metadata
   matches the run — model `z-image-turbo-Q4_K_M`, family `z_image`, 1024×1024, steps 8,
   seed 1655275034, sampler `euler`, cfg 1, `txt2img`, prompt identical to the UI.
6. **Cancellation through the UI.** Second run submitted, then the button (labelled ABORT)
   invoked ~3.5 s later. Status changed to "Aborted", button returned to GENERATE, spinner
   gone. The UI said "Aborted" — it did **not** claim the backend stopped computing.
   Shots: `10-second-active.png`, `11-aborted.png`.
7. **Behaviour immediately after the abort (discrepancy — see below).** The next Generate
   click did not show a refusal; it started a new run at "Step 1 of 8" and a **new**
   `sd-server.exe` PID 14320 (start time 19:19:14) replaced PID 20892.
   Shot: `12-refused.png`.
8. **That run failed with a backend 500.** Status line:
   `Error: sd-server txt2img error 500: {"error":"generate_image returned no results"}`.
   No gallery artifact was written (gallery stayed at 21 entries, newest still the first
   run). Shot: `13-third-complete.png`, `14-after-500.png`.
9. **Another Generate click after the 500.** No refusal text appeared and no new server
   started; the status line stayed on the same 500 text. I could not distinguish "click
   ignored" from "same failure repeated", so this is recorded as **inconclusive**, not as
   either behaviour.
10. **Recovery through the documented controls.** Unload, then Reload. A fresh
    `sd-server.exe` PID 11920 started at 19:22:17, the panel showed "Model ready:
    z-image-turbo-Q4_K_M", and a generation ran with live preview at "Step 4/5 of 30" with
    `sd-server` at ~12.8 GB. The UI stayed responsive throughout. I aborted that run to bound
    GPU time; status "Aborted" again.
    Shots: `15-unloaded.png`, `16-reloading.png`, `17-reloaded.png`,
    `18-idle-after-reload.png`, `19-recovery-start.png`, `20-recovery-result.png`,
    `21-aborted-recovery.png`.

## UI-observed discrepancies (not explained away)

1. **No refusal was surfaced after a cancelled generation.** The frozen contract describes a
   cancelled post-dispatch generation as leaving the engine unavailable until it is
   restarted, with a second request refused. What the UI did instead was start a new run on
   a **new process**. Source shows `CancelImageJob → Manager.Cancel → Engine.CancelJob →
   cancelGeneration`, which only closes the request — it does not unload or restart. **Why a
   new `sd-server` appeared is not established by this test.** Either the aborted work took
   the server down and the auto-load path started a new one, or another path restarted it;
   no log records it (`howl-errors-*.jsonl` contains no entry for these events, and the
   engine's stdout is not persisted).
2. **A 500 followed the restart.** `generate_image returned no results` is a live
   reproduction of the failure family the project has never explained. It occurred on the
   third run, after an abort and a silent restart — a sequence this project has repeatedly
   said must not be trusted. It does not establish the cause of the historical 500, and it
   does not show that cancellation stops computation.
3. **The post-500 state was visibly ambiguous.** No "unavailable until restarted" message,
   no visible rejection, and no distinguishable new attempt. The engine's slot state after a
   500 could not be read from the UI at all.

None of these were repaired here: this stage is acceptance only.

## What is proven, and by what

- **UI-proven:** Image Studio opens; model discovery reports 39 models; model selection
  applies family defaults; a real generation loads a model, streams progress, completes and
  saves a correct artifact; Abort produces an "Aborted" state without claiming the backend
  stopped; Unload/Reload restores a working engine that generates again; the UI remains
  responsive and reports outcomes.
- **Established only by source and unit tests:** the double-settlement invariant (handle
  invalidated under `genMu`, second settlement inert and reported as
  `errGenerationSettled`). No UI path can invoke it, and this acceptance makes no claim that
  it did.
- **Observed but unexplained:** the silent engine restart after abort, and the 500 that
  followed it.

## Not exercised / blocked

- Img2img, inpaint, outpaint, upscale and video tabs; batch counts > 1; cancellation during
  model loading (cancellation happened during generation); the synchronous `GenerateImage`
  path (no UI uses it); any VRAM planning readout (the UI exposes none).
- The E: drive was not touched, scanned or mounted. No GPU allocation was forced and no
  parallel generation was run; each run was started only after the previous one ended.
- No other process was disturbed; only UI invocations and read-only process queries were
  used.

## Screenshot inventory

All under `build/logs/screenshots/`: `02-app-window.png`, `03-image-studio.png`,
`04-model-list.png`, `05-model-open.png`, `06-model-selected.png`, `07-generating.png`,
`08-after-generation.png`, `09-completed.png`, `10-second-active.png`, `11-aborted.png`,
`12-refused.png`, `13-third-complete.png`, `14-after-500.png`, `15-unloaded.png`,
`16-reloading.png`, `17-reloaded.png`, `18-idle-after-reload.png`, `19-recovery-start.png`,
`20-recovery-result.png`, `21-aborted-recovery.png`. `01-app-initial.png` is the discarded
desktop composite described above.
