# HOWL Chat — Beta Stage Status

Source of truth for the numbered list: `Howl Chat issues to correct -.txt`.
Repairs first; additions only with approval.

## ControlNet (repaired — it had never worked)

Three separate breaks, all fixed:

1. **Model path never sent.** `types.ImageOptions.ControlNetPath` existed and the UI
   populated it, but `sdcpp.txt2imgRequest` / `img2imgRequest` had no
   `control_net_path` field, and none of the three request builders assigned it. The
   engine therefore received a control image with no ControlNet to apply it. Added the
   field to both structs and the assignment to all three builders.
2. **Control image dropped.** `generate()` saved the uploaded control image to a temp
   path, then never put it in the request. Now sent as `control_image_path`.
3. **Data-URL decode failure.** `SaveTempImage` fed the whole `data:image/png;base64,…`
   string to `base64.StdEncoding`, so anything produced by a canvas editor failed with
   `illegal base64 data at input byte 4`. New `decodeBase64Image` accepts both bare
   base64 and full data URLs.

A control image with no ControlNet model selected now stops before generation with a
clear message instead of burning a full sampling pass.

## Workflows (new)

InvokeAI-style saved recipes, backed by disk files rather than browser storage.
`sdcpp` is not a graph engine, so a workflow is the complete UI state, not a node graph:
mode, every panel control, the LoRA stack and an optional pose.

- Go: `app-image-workflows.go` — `ListImageWorkflows`, `SaveImageWorkflow`,
  `GetImageWorkflow`, `DeleteImageWorkflow`, `OpenImageWorkflowsFolder`.
  Files live in `%AppData%\HOWL_Chat\workflows\*.json`, written temp-then-rename,
  names sanitised so they cannot escape the folder, 2 MiB cap, version checked.
- UI: Workflows panel with save / apply / delete / export / import / open folder.
- Verified live: saved `ZZ Test Recipe`, listed it, changed the prompt, applied it —
  **100 settings restored**, prompt returned to the saved value, pose restored, then
  deleted (list back to 0).

## Pose editor (new)

`frontend/dist/js/ig-pose.js` — COCO-18 OpenPose skeleton on black, exactly the palette
and layout OpenPose ControlNets were trained on. Drag any joint, six presets
(standing, T-pose, sitting, walking, arms up, running), flip, reset, 512/768/1024 export.
Pose is persisted and travels inside a workflow.

Verified: canvas 478 px box at 1384×798, **8,901 lit pixels** with red limb pixels
present, preset switch redraws (7,843 px), drag reports "Moving nose", export is a
38,654-byte PNG at 768 that saves through `SaveTempImage` to
`images/tmp/init_20261031_231618.png`. Dialog fits at 900×620 (760×572, internal scroll).

## Generation completion watchdog (hardening)

While checking that the request changes had not broken generation, a run finished
(backend RPC confirmed `"status":"completed","progress":1`) but the UI stayed on
"Generating…" with an Abort button because the status poll stalled. A 20 s watchdog now
aborts such a run with "Lost contact with the image engine… check the gallery" instead
of spinning forever.

## App-wide adversarial check

- All 9 pages return 200 and **410 local asset references resolve**; the only two
  broken refs are `/wails/ipc.js` and `/wails/runtime.js`, which the harness shim
  replaces and which do not exist outside WebView.
- Image page console: no script errors from the new modules.
- One full run after the engine edits: `sampling completed, taking 9.89s`, result saved
  to gallery, TAESD previews 0.17 → 1.00, no ControlNet errors.
- Mode scoping and model-family filtering still behave (video hides the SDXL checkpoint,
  still modes hide the Wan models).

## ControlNet end-to-end (installed and proven)

Model: `E:\HOWL_Chat\models\image\controlnet_openpose_sdxl_thibaud_fp16.safetensors`
(2,502,139,104 bytes, fp16, 844 tensors, diffusers layout, verified byte-exact).
Source: `https://huggingface.co/ckpt/controlnet-sdxl-1.0` — `thibaud_xl_openpose.safetensors`,
the OpenPose v2 SDXL ControlNet trained on SDXL base 1.0. Chosen over the 5 GB fp32
`OpenPoseXL2.safetensors` because the GTX 1080 has 8 GB.

Two further breaks had to be fixed before it could work:

4. **Wrong CLI flag.** The engine passed `--control_net`; this sd-server build only accepts
   `--control-net`. The rejected argument made sd-server refuse to start at all
   (`error: unknown argument: --control_net`). Fixed in `buildServerArgs`.
5. **TAESD stole the model slot.** In a bundle directory the first unrecognised file became
   `files["model"]`, and `taesdxl.safetensors` sorts before the checkpoint, so the engine
   was about to load a 9 MB preview decoder as the diffusion model. TAESD now has its own
   `taesd` slot and the largest file wins `model` (latent bug, only visible after a rescan).

Verified run (job `job_1793497153218097700`, 512×512, 4 steps):
server launched with `--control-net …thibaud_fp16.safetensors`, log shows
`controlnet 2387.61MB (VRAM)` alongside `diffusion_model 4900.07MB (RAM)`, the request
carries `control_image` plus `control_strength: 0.9`, and the job reached `completed 4/4`.
GPU sat at ~2.1 GB of 8 GB; the rest is offloaded to RAM by `--max-vram -1`.

Note on architecture: a ControlNet is a **server-side slot** in sd.cpp ("a `control_net`
slot enables `controlnet`"), so the file must live in the checkpoint's bundle directory to
be attached at load. The UI reflects this: the ControlNet panel states
"Bound to this model: … — pose applies automatically", and a pose no longer requires a
dropdown pick.


## T1.1 - State ownership (closed 2026-10-31)

Page switching is `window.location.href = page`, a full document reload, so every
JS value dies and must rehydrate. It now does.

**A fatal emit bug had to be fixed first.** Wails v2.12.0 `runtime.EventsEmit` calls
`log.Fatalf` (os.Exit) when `ctx.Value("events")` is nil. `app-model.go` emitted
`model:progress` from its load monitor and `app-prompt.go` emitted `chat:chunk` per
token, so any emit outside a live webview terminated the whole process - the UI harness
died twice exactly this way. New `App.emitEvent()` in `app.go` drops the event when
there is no dispatcher; 11 model emits, 13 chat emits and the 3 emits in
`internal/errors/wails.go` now route through a guarded helper. This also removes a real
crash class: an emit from a goroutine that outlives window teardown could otherwise exit
the app in production.

Evidence (real runs, harness on 127.0.0.1:8799):

| Check | Result |
|---|---|
| baseline VRAM | 1414 MiB |
| LoadModel returns | 1.9 s (async by design: spawn then monitor) |
| model ready | 6.0 s, `IsModelLoaded` true, stage `Ready` |
| VRAM while loaded | 6366 MiB, llama-server pid 25140 |
| navigate to image page | VRAM 6407, same pid |
| navigate back | VRAM 6352, same pid, name `LFM2.5-8B-...gguf` shown |
| streamed reply after return | content tokens produced (reply "SWITCH TEST OK") |
| navigate away *during* load | model still reached ready, 0 fatal lines |
| sampler set via UI | temperature 1.37, top_k 77 in `config/settings.json` |
| process restart -> reload | slider AND number field show 1.37 / 77 |
| restored defaults | 0.7 / 40 written back and shown on reload |

Known harness limit, not an app bug: streamed tokens arrive through
`window.runtime.EventsOn('chat:chunk')` and backend-originated events cannot reach the
harness shim, so a streamed reply shows an empty bubble there. `SendMessage` (the
synchronous binding) is the reliable answer check in the harness.

New UI finding, logged not fixed (T2.1 / fluid layout): at ~428 px width the header
row overlaps - "No model loaded", the page title and nav items stack on each other.
At 1440 px there is no horizontal overflow and the layout is clean.

## T1.2 - VRAM arbitration between chat and image models (closed 2026-10-31)

One 8 GB card cannot hold a chat model and a diffusion checkpoint at once, so the app
now says so before a load fails quietly. New bindings in `app-vram-arbitration.go`:

- `CheckModelFit(kind, path)` -> `can_fit`, `required_mb`, `free_vram_mb`,
  `total_vram_mb`, `blocking[]`, `message`. Free VRAM is re-read live through the
  resource manager (nvidia-smi) on every call, never inferred from the registry.
- `FreeVRAMFor(kind, path)` -> unloads whatever is in the way, waits 1.5 s for the
  driver to release the allocation, then re-reports the fit.

The image requirement is measured, not invented: 1600 MB base (text encoders, VAE,
sampling buffers) + 500 MB overhead + the bound ControlNet's file size. Observed on
this machine: sd-server held 2092-2180 MB for a 4-step SDXL run, and ~2.1 GB with the
2.4 GB ControlNet bound in the earlier ControlNet run.

UI gates: `chat-model.js` before `LoadModel`, `image-generator-core.js` inside
`onGenerate` before `StartImageJob`. Both put the backend's sentence in `confirm()`
and call `FreeVRAMFor` only on OK.

| Check | Result |
|---|---|
| chat model loaded | 6326-6382 MiB |
| fit(image) | can_fit=false, needs 2100, free 1847-1870, blocking chat 4028 |
| dialog text (real click) | names the chat model at 4028 MB, needs 2100 MB, 1847 MB free of 8192 MB, ends "Unload it and continue?" |
| Cancel chosen | status "Cancelled - not enough VRAM with the chat model loaded.", chat still loaded, VRAM unchanged |
| OK chosen | chat unloaded, sd-server pid 9536 took the card, status "Generating...", button -> Abort |
| job result | `sampling completed, taking 8.04s`, PNG `img_20261101_021037_800.png` |
| after FreeVRAMFor(image) | free 6749, chat unloaded, VRAM 1443 |
| reverse: image loaded (2092) + 20B chat | can_fit=false, needs 7591, free 6100, blocking image 2100, mirrored wording |
| reverse FreeVRAMFor(llm) | sd-server gone, VRAM 1439 |
| both fit (SDXL 2180 + 8B 4028 in 8192) | no dialog - correct, the card really does fit both |

The confirm dialog was stubbed in the harness (`window.confirm`), because a native
modal blocks the automation. The gate, message text, cancel path and accept path are
all real code; the dialog itself was not clicked by hand in the Wails window.

Two UI defects found while wiring this, logged not fixed:
- The image model dropdown lists "(2 incompatible files hidden)" as a **selectable
  option** whose value is that literal string; picking it fails with
  "model not found". It should be a disabled label.
- `onLoadModelClick` is called from the Ctrl+L shortcut at
  `image-generator-ui.js:1603` but is defined nowhere in the codebase. The image page
  has no Load button at all; the model loads implicitly on Generate.

## T1.3 - Think tags / reasoning (closed 2026-10-31)

Researched against llama.cpp's own documentation and issue tracker before writing
anything. Three findings shaped this:

1. `--reasoning-budget N` exists in this binary: -1 unrestricted, 0 ends thinking
   immediately, N>0 is a hard cap. Confirmed from `llama-server --help` (build 12015,
   commit a34de8618), not from memory.
2. `--reasoning-budget-message` is not optional in practice. The llama.cpp maintainer
   states that without an injected message the model, once the budget runs out, keeps
   thinking *inside the answer* instead of answering (issue #22717). That is exactly the
   loop the user asked to fix, so budget and message are always set together.
3. Reasoning arrives as a separate `reasoning_content` delta field (verified in this
   app's own stream log for LFM2.5), so it can be diverted without touching content.

Implementation:

- `app-llama.go` builds the reasoning flags from `model.think_mode`:
  off -> `--reasoning off --reasoning-budget 0`; on/hidden ->
  `--reasoning on --reasoning-format deepseek --reasoning-budget <budget>` plus the
  budget message. Default is now **hidden with a 2048 budget**: the previous behaviour
  passed no flags at all, which means llama-server's default of -1 (unrestricted) - the
  loop-prone setting.
- `internal/backend/llama/http_client.go` parses `reasoning_content` and delivers it to
  an `OnReasoning` sink; `app-prompt.go` re-emits it as `chat:reasoning`, never into
  the content callback, so it can never become the reply or enter history.
- `app-think.go`: `stripThinkBlocks()` removes inline blocks (think/thinking/pipe forms,
  unterminated included) from the assistant text before it is appended to history.
- Thinking gets its own headroom: when a mode is on, `n_predict` is raised by the
  budget, so a thought cannot eat the answer's token allowance.
- UI: `index.html` gains a Reasoning card (Off / On - show the thought / Hidden),
  backed by `chat-think.js`, which streams reasoning into a collapsible block
  (`.chat-think` + existing `updateThinkArea`/`toggleThink`/`cc-thinking.css`).

Evidence:

| Check | Result |
|---|---|
| hidden mode server args | `--reasoning on --reasoning-format deepseek --reasoning-budget 1024 --reasoning-budget-message "I have thought enough. Time to answer."` |
| hidden mode stream | `Stream reasoning: 131 chunks, 570 chars (never fed back)` |
| reply with thinking on | `
Blue` - clean answer, no thought text |
| off mode server args | `--reasoning off --reasoning-budget 0` |
| off mode stream | `Stream reasoning: 0 chunks, 0 chars` |
| strip unit test | `go test -run TestStripThinkBlocks` -> ok (7 cases: complete, only-block, multiline, pipe tags, unterminated) |
| UI controls | `#think-mode` + `#think-budget` present in index.html, `chat-think.js` loaded |

Limits of this evidence:
- The reasoning spoiler could not be visually verified: backend-originated Wails events
  never reach the harness shim, so nothing renders from `chat:reasoning` there. The
  event, the sink and the render path exist; the rendered block has not been seen.
- The harness does not persist chats (no `chats/*/history.ndjson` written from RPC
  turns), so "reasoning is absent from the saved transcript" rests on the code path plus
  the strip unit test, not on an inspected history file.
- Mode changes take effect on the next model load, because the budget is a server
  start flag. The UI says so; it does not force a reload.

## T1.4 - Per-family sampler defaults and prompt format (closed 2026-10-31)

Only families with a *published* recommendation get numbers, so nothing here is
invented. Sources are in the comments next to each table row in
`app-family-defaults.go`:

| Family | Values | Source |
|---|---|---|
| Qwen3/3.5 | temp 0.6, top_p 0.95, top_k 20, min_p 0 | Qwen3-8B model card (thinking mode) |
| Gemma 3/4 | temp 1.0, top_p 0.95, top_k 64 | Gemma 4 model card; Gemma 3 ships same top_k/top_p |
| DeepSeek-R1 distills | temp 0.6, top_p 0.95 | DeepSeek-R1 model card (0.5-0.7) |
| Llama 3.x, Mistral-Nemo, others | generic defaults | no vendor recommendation exists - deliberately not invented |

Family defaults are the *base*: `buildInferenceOptions()` starts from them and then
applies the user's sampler settings, so an explicit choice always wins.

Three real bugs were found while proving this:

1. **Sticky chat template.** The first model's template was written into the saved
   settings behind an "is it already set?" check, so every later family was loaded
   with the *first* family's template. Resolution now happens per load in
   `spawnSettingsFor()` (`app-model-template.go`); a template the user typed is kept,
   an auto-injected one is replaced.
2. **Fallback template matched by exact family name.** The switch tested `case "qwen"`
   while the real architectures are `qwen35` and `gemma4`, so both fell through to one
   generic `role: content` template - both came out at exactly 95 bytes, which is how
   it was caught. Now prefix-matched, with real Qwen (`<|im_start|>`) and Gemma
   (`<start_of_turn>`) formats.
3. **`SetSettings` merged sections instead of replacing them.** A cleared sampler
   section left the old values alive and silently overrode the family defaults.

Evidence (fresh sampler settings, two families, same binary):

| Model | family | temperature | top_p | top_k | min_p | template |
|---|---|---|---|---|---|---|
| Qwen3.5-9B | qwen35 | 0.60 | 0.95 | 20 | 0.00 | 170 bytes |
| Gemma4-12B | gemma4 | 1.00 | 0.95 | 64 | 0.05 | 133 bytes |
| before the template fix | - | - | - | - | - | both 95 bytes (identical) |

Also verified: `n_predict` 2560 = 512 + the 2048 thinking budget, so the thought does
not eat the answer's allowance; the non-streaming `SendMessage` path strips reasoning
before returning; and every sampler card in the UI has a consumer - top-p, min-p,
typical-p, frequency-penalty, presence-penalty, mirostat-mode, dyn-temp-range all map
through `buildInferenceOptions` into `buildRequestFromMessages`, and the Reasoning card
drives the server flags. No card is ignored by the backend.

Known issue, logged not fixed: Qwen3.5 returned an empty reply on one run - with
reasoning enabled the thought can still consume the whole allowance on that model, even
with the budget message. Worth watching in beta.

Limit: the UI now displays effective values via `GetEffectiveSampler()`, but that was
not re-checked visually in the browser after the change.

## Generation settlement guard (closed 2026-11-02)

Closes the image-generation settlement item opened by the ownership/lifecycle reviews.
No product code changed in the closing pass: the guard was already in the working tree
(`sdcpp/generation.go`, `sdcpp/settlement_guard_test.go`, 19:08); what was missing was
evidence from the real application.

| Gate | Result | Evidence |
|---|---|---|
| Source review | PASS | `image-generation-settlement-guard-review.md` |
| Regression tests | PASS (reported in stage record) | `image-generation-settlement-guard-stage.md` |
| Independent review | PASS | `image-lifecycle-followup-r1-r2-r4-r5-r6-review.md` |
| Real application acceptance | PASS on a verified build | `image-generation-abort-500-discrepancy-investigation.md`, addendum |

Live checks on the verified binary (sha256 `d7e8aacd…`, PID 9252, start 20:20:42, path
`build/bin/howl-chat.exe`): cold auto-load spawned the child and reached "Step 1 of 8";
abort reported "Aborted" and **left the child alive**; the next request was **refused** with
`errEngineUnavailable` naming the abandoned owner; Unload closed the child and Reload
brought up a new one ("Model ready"); a generation was then admitted on that child and
aborted, leaving the app idle with no orphans.

**Build-identity lesson (process, not code).** The first acceptance ran an 18:59 binary
that predated the guard by nine minutes, so its results described the old behaviour (the
post-abort request was *admitted*, which is how it reached the backend and returned the
500). Any future acceptance must record revision + worktree diff hash + binary hash +
running PID/start time before interpreting a failure.

**Deliberately still open, not part of this item:** the status-line elision (UI repair
pass), the unexplained sd-server exit (logging/diagnostics audit), video `stopServer`
shared-process ownership (separate defect), and Image Studio's broader lifecycle and
model-compatibility work. The settlement guard is closed; Image Studio is not.

## UI conventions (apply to every future change)

- Menu and panel icons must be **real icons from the theme**, never ASCII/emoji text.
  Use the existing `<span class="icon icon-NAME"></span>` markup and only names the
  CSS actually defines (`frontend/dist/css/*.css`); check before using one.
- An icon must match its neighbours in weight, colour and size - the surrounding
  cards are the reference, not a new invention.

### Polish backlog (deferred - deliberately not started)

Captured from the group-chat review. These are visual polish, not defects, and
must not be touched until the earlier tiers are closed.

1. **User bubble borders must not bump.** The user message bubble's border appears
   to graze the neighbouring character bubbles. This may simply be the preview
   viewport width rather than a real layout fault, so check at real window sizes
   (wide and narrow) before changing any spacing.
2. **Matching user-colour halo.** Character bubbles carry a halo tinted to that
   character's colour; the user bubble should carry the same treatment in the
   user's own colour instead of the current neutral edge.

## Still open
- Other pages (chat, characters, worlds, lorebooks, settings) still need the fluid-layout
  pass.
- Inpaint / Outpaint / Upscale end-to-end runs.

## Group chat verification - three characters, click-built (2026-10-31)

Created three characters through the real creation form with no images
(Yara / dockhand, Fen / archivist, Kesh / gate guard), then started a new chat,
added all three with Add Character, and sent one message.

Working:

- All three are created, persisted to IndexedDB and mirrored to localStorage, and
  appear in the character grid and the chat picker.
- A single send produces three sequential replies in list order through
  CharacterQueue: Yara 8.8s, Fen 12.2s, Kesh 2.8s. Each answers in its own voice.
- A second message repeats the 1-2-3 rotation correctly.
- Reasoning renders as a collapsible 'Show reasoning' block on each character
  reply, with no raw tags in the visible text.

Fixed during this test:

- `GenerateCharacterResponse` never set `opts.OnReasoning`, and no path filtered
  inline thinking out of the content stream, so Qwen3.5's native `think` tags and
  the whole thought were rendered into the reply and saved into history. Added
  `thinkSplitter` (`app-think-stream.go`, 6-case test) used by both the generic and
  the character streaming paths. Before: raw `...` visible, Reasoning card empty,
  `Stream reasoning: 0 chunks`. After: no tags, card populated.

Open defects found (not fixed, need a decision):

1. First group turn after a page load labels all three bubbles with the first
   character's name. The text is correct ('Fen: ...', 'Kesh: ...') but the bubble
   header reads Yara for all three. A second turn labels correctly.
2. Regenerate rewrites the reply as **'AI'**, losing the character, and always
   regenerates the *last* assistant message regardless of which bubble's button
   was clicked (`regenerateMessage` removes the last assistant node and calls
   `streamAssistantResponseFromBackend`). In a group chat that is the wrong turn
   and the wrong speaker.
3. The Next (variant) button leaves the text unchanged in this test - no alternate
   response was produced or shown.
4. The new chat wrote no `chats/<id>/history.ndjson`, so turn persistence could not
   be verified in the harness. Chat state itself does reload (characters intact
   after a restart), so this is the append log specifically.
## T2.1 - Theme unification (STAGE A COMPLETE, B-D PENDING)

### Root cause (evidence-backed, not inferred)

The light theme defined a full token vocabulary (`--surface-*`, `--border-*`,
`--text-*`, `--input-bg`, 32 tokens total). **The dark theme defined none of
them.** Consequences:

- Any rule written as `var(--surface-primary)` resolved to nothing in dark.
- Light mode could not be a token flip, so the codebase grew parallel override
  universes: theme-light-layout (8 var refs), theme-light-components (153),
  theme-overrides (20), plus per-feature `*-light-theme.css` files.
- Every feature stylesheet invented its own literals for the same ideas:
  `rgba(100,100,100,0.3)` borders, `rgba(212,212,212,x)` text,
  `rgba(21,21,31,x)` panels, `rgba(10,10,15,x)` strips.

Two self-inflicted defects found during the audit, both now fixed:

1. Four **cyclic** custom properties in `theme.css :root`
   (`--shell-nav-bg: var(--shell-nav-bg)` and friends) left by a bulk
   substitution that ran over variable definitions. A cyclic custom property is
   invalid at computed-value time, so every `background: var(--shell-panel-bg)`
   collapsed to transparent.
2. My earlier measurement was also wrong: reading `getComputedStyle`
   immediately after a theme class flip can return the pre-transition value.
   Audits now disable transitions first.

### Process finding (explains "none of your fixes show up")

`main.go` uses `//go:embed all:frontend/dist`, so the packaged exe contains the
CSS **as of build time**. Screenshots taken from the app after a CSS edit show
stale styles until it is rebuilt:

    CGO_ENABLED=0 GOFLAGS=-buildvcs=false go build -tags "desktop,production"       -ldflags "-w -s -H windowsgui" -o build/bin/howl-chat.exe .

### What was done

- Dark now defines the full semantic vocabulary in `theme.css :root`: 10 surface
  /border/text/input tokens plus 31 accent, status and component tokens
  (`--danger`, `--accent-primary`, `--button-*`, `--input-border`, shadows...).
  Only `--topbar-height` is still referenced without a dark definition (layout,
  not theme).
- 577 literals converted to tokens across the dark stylesheets: 479 first pass,
  38 alpha-agnostic pass, 52 theme-dependent colour/hex pass, 8 stragglers.
  Remaining residue in the four literal families outside token definitions: **0**.
- Light canvas unified to a single value: `--app-bg: #F5F4F0` used by body,
  `.main-content` and page containers, so they cannot disagree again.
- 110 border-token-as-background rules repaired to `--surface-hover`; 56 of those
  were then correctly returned to the border colour because the selector is a
  1px divider (`.menu-separator`, `.info-separator`, scrollbar thumbs).
- One theme implementation: `menu-view.js setTheme` delegates to `theme.js
  setThemeMode`; `real-settings.html` and `audio-settings.html` now load
  `theme.js`. Verified by `--theme-glow` flipping purple -> teal.
- All CSS `@import`s version-busted; harness serves every asset `no-store`.

### Verification (measured, not eyeballed)

Automated contrast audit (16px+, text nodes, resolved ancestor background,
difference < 60 in luminance), transitions disabled:

| Page | Light problems | Dark problems |
|---|---|---|
| index (chat) | 0 | 0 |

Screenshots: `build/t21/chat-light-stageA.png`, `build/t21/chat-dark-stageA2.png`
(also rendered inline in the session).

### Still open

1. Stages B/C: run the same audit on the other 8 pages in both themes
   (`character`, `world`, `scenario`, `lorebooks`, `story-tools`,
   `image-generator`, `audio-settings`, `real-settings`).
2. Now that dark owns the vocabulary, the parallel override sheets
   (theme-light-components 153 refs, theme-overrides 20, per-feature
   `*-light-theme.css`) are candidates for deletion - measure before removing.
3. Header collision at narrow widths (`Help` overlapping `NO MODEL LOADED` at
   ~780px) - layout, not theme.
4. Image page joined the same system in Stage B - the exemption came from a stale comment.

### Note on an earlier wrong claim

An intermediate note said 86 stylesheets were "loaded by no page". The import
graph resolves **122/122 reachable** - they are pulled in via `@import`. Nothing
is dead; the layering itself is the debt.

### Role accents (fixed after Stage A review)

The user bubble hardcoded teal, while the assistant bubble took its border/glow
from `--theme-glow`, which the theme cycle rewrites (teal -> cyan -> purple).
Both bubbles therefore went teal whenever the cycle stood on index 0.

Role identity is now tokenised and rotation-independent:
`--user-accent/-bg/-border/-glow(+ -strong)` and `--ai-accent/-border/-glow(+ -strong)`,
defined in `theme.css :root` (dark) and `body.theme-light` (softened values).
`cc-messages.css` and `cc-light-theme.css` read only those tokens.

Measured (transitions disabled, chat page):

| | border | glow |
|---|---|---|
| dark assistant | rgba(168,85,247,0.45) | rgba(168,85,247,0.22) |
| dark user | rgba(20,184,166,0.3) | rgba(20,184,166,0.15) |
| light assistant | rgba(126,34,206,0.35) | shadow-card |
| light user | rgba(15,118,110,0.45) | rgba(20,184,166,0.1) |

Re-setting `--theme-color`/`--theme-glow` to the teal theme leaves the assistant
bubble at rgba(168,85,247,0.45) - the reported bug no longer reproduces.

Still deferred by request: user-bubble halo and the border "bump" refinement.

## T2.1 Stage B - Backgrounds and cross-page unification (IN PROGRESS)

### Light was "blinding" for a specific reason

The canvas had been unified, but `--surface-primary` was `#FAFAF8` and
`--surface-secondary` `#F5F4F0` - i.e. every card, panel and input sat at
maximum luminance on top of near-white paper. A light theme needs three
separable levels, not two almost-whites.

New light palette (token values in `theme-light.css`):

| token | value | role |
|---|---|---|
| --app-bg | #E9E6DF | canvas |
| --shell-panel-bg | rgba(247,245,240,.97) | cards / panels |
| --shell-aside-bg | rgba(233,230,223,.97) | side rails |
| --shell-nav-bg | rgba(238,235,228,.97) | top bar / nav |
| --surface-primary | #F7F5F0 | raised card |
| --surface-secondary | #EFEBE3 | recessed band |
| --surface-tertiary | #E4E0D6 | chrome / track |
| --border-primary | #CFC9BE | hairlines |
| --text-primary | #262420 | body text (not pure black) |

Accent darkened to `#0d9488` in light (teal `#14b8a6` sat at ~2:1 on paper),
and status colours deepened (`#c62828`, `#15803d`, `#b45309`, `#0369a1`) since
the dark-theme values are luminous by design and unreadable on light.

### The blockers found by adversarial scan

1. **JS-built UI bypassed the theme entirely.** Modals, context menus and
   hover states were assembled with literal colours in `style` attributes:
   `rgba(0,0,0,0.7)` scrims, `rgba(255,255,255,0.06)` hovers.
   67 literals in `context-selectors.js`, 55 in `sidebar-characters.js`,
   48 in `memory-candidates.js`. No CSS fix can reach an inline style.
   Inline styles *do* resolve custom properties, so 119 of them were rewritten
   to `var(--surface-hover)`, `var(--overlay-scrim)`, `var(--surface-subtle)`,
   `var(--border-secondary)`, `var(--success-bg)` - line-scoped, and skipped
   canvas/`fillStyle` colour paths.
2. **`theme-light-layout.css` repainted every page region** with
   `--surface-secondary` while `body` painted `--app-bg`, so canvas and content
   disagreed by design. Dark has no such rule. Removed (both the content block
   and the sidebar/chrome block); those regions now read their shell tokens.
3. **Page-level canvases repeated the same mistake**: `.lorebooks-page`,
   `.story-tools-page`, `.story-tools-grid`. Removed. The image page keeps its
   own sheets untouched, by instruction.
4. Dark-canvas literals left in CSS outside token definitions: **0**
   (`var(--app-bg, #0a0a0f)` fallbacks excluded).

### New tokens (defined in both modes)

`--overlay-scrim`, `--surface-subtle`, `--surface-hover-strong`, `--success-bg`.

### Cache-busting

`?v=7` -> `?v=8` across 20 html/css files, plus the harness already sends
`no-store`. Any future CSS edit needs a version bump or the browser reuses the
old sheet.

### Verified so far

| page | light canvas | stuck-dark elements | low-contrast |
|---|---|---|---|
| character | corrected to --app-bg | 0 | 0 |

### Open

- Remaining pages to verify: index, world, scenario, lorebooks, story-tools,
  image-generator, audio-settings, real-settings.
- `world.html` hangs the browser tab on load (navigation never fires, all later
  evaluations queue behind it). Investigate before trusting any world-page
  screenshot: a blocked main thread would also prevent the theme flip.
- Multi-page iframe auditing is not viable: 22 stylesheets per page starves the
  10s evaluation cap. Verify page by page.

### Stage B measured result (transitions disabled, 600ms settle)

| | light | dark |
|---|---|---|
| body | rgb(233,230,223) | rgb(10,10,15) |
| .main-content | rgb(233,230,223) | rgb(10,10,15) |
| .nav-sidebar | rgba(238,235,228,.97) | rgba(10,10,15,.95) |
| .world-content | rgba(233,230,223,.97) | rgba(15,21,31,.95) |
| stuck-dark areas | 0 | 5 (expected) |
| low-contrast text | 0 | 0 |

Two more duplicate painters were found while measuring, both removed:

- `layout.css` carried its own "LIGHT THEME OVERRIDES" tail that repainted
  `.main-content`, `.nav-sidebar`, `.character-sidebar`, `.chat-panel`,
  `.world-sidebar`, `.input-bar`. 0 background declarations remain in it.
- `theme-light-components.css` + 6 other sheets repainted the page content
  regions (`*-content`) through comma-separated selector lists. 7 blocks
  cleaned; **0** light-mode content repaints remain outside the image page.

Measuring note, relearned the hard way: `theme.js` runs a 4s custom-property
transition, so a `getComputedStyle` read milliseconds after the flip returns the
*previous* theme's colour. One intermediate reading reported `rgb(10,10,15)` for
a light body purely because of that. Every number above was taken after
injecting `*{transition:none!important;animation:none!important}` and waiting
600ms.

### Still open

- Page-by-page verification remaining: index, scenario, lorebooks, story-tools,
  audio-settings, real-settings (character, world and image-generator done).
- World page clears its own `theme-light` repaints, but check its modals at
  runtime - the JS-built dialogs were tokenised, not visually confirmed.

### Stage B final numbers (world.html, transitions disabled)

| measure | light | dark |
|---|---|---|
| body canvas | rgb(233,230,223) = --app-bg | rgb(10,10,15) = --app-bg |
| .main-content | rgb(233,230,223) | rgb(10,10,15) |
| .nav-sidebar | rgba(238,235,228,.97) | rgba(10,10,15,.95) |
| .world-content | rgba(233,230,223,.97) | rgba(15,21,31,.95) |
| primary button | #0d9488 fill / #fff label (delta .536) | theme-glow fill / #d4d4d4 |
| stuck-dark areas in light | 0 | - |
| low-contrast text | 0 | 0 |

### The button bug the audit could not see

The ancestor-background audit skipped labels inside child spans, so it reported
0 offenders while `Create New World` rendered dark grey on teal. Found by
looking at the screenshot. Cause: `world-header.css` set a primary button to
`background: var(--button-primary-bg)` with `color: var(--text-secondary)`.

New token `--on-accent` (white, both modes) now labels every accent-filled
button; 13 blocks corrected across 10 sheets. Check for
`background: var(--button-primary-bg)` + `color: var(--text-*)` returns 0.

Lesson for the remaining pages: the numeric audit is necessary but not
sufficient - a filled child element needs its own passed-in colour check, and a
screenshot per page is still required.

### Image generator joins the system (exemption retracted)

The "do not touch the image page palette" note was not a design decision - it was
a comment left behind during an earlier revert, written as a reminder to an agent
that would not otherwise pay attention. Removed, along with the page-level canvas
repaints in `ig-light-theme.css` / `igp-light-theme.css`.

| measure | light | dark |
|---|---|---|
| body | rgb(233,230,223) | rgb(10,10,15) |
| .ig-layout | rgb(228,224,214) = --surface-tertiary | rgb(26,26,36) = --surface-tertiary |
| .ig-panel | rgb(247,245,240) = --shell-panel-bg | rgb(21,21,31) = --shell-panel-bg |
| stuck-dark areas | 0 | 5 (correct) |
| low-contrast text | 0 | 0 |

The page keeps its character (image grid, canvases, its own red danger states)
while sourcing every surface from the shared tokens, so it flips with everything
else.

### Undefined token names (found while unifying the image page)

`igp-layout.css` referenced `--bg-secondary`, `--bg-hover` and `--accent-color`.
None of those names exists anywhere in the sheets, so every use silently fell
back to its hardcoded dark literal - a page-cannot-flip bug hiding in plain
sight. 3 sheets fixed, aliases mapped:

| used | real |
|---|---|
| --bg-primary / secondary / tertiary | --surface-primary / secondary / tertiary |
| --bg-hover | --surface-hover |
| --text-color | --text-primary |
| --border-color | --border-primary |
| --panel-bg | --shell-panel-bg |
| --accent-color | --accent-primary |

The image sheets also carried the pre-unification teal `#14b8a6` (21 uses),
`#4ade80`/`#22c55e` success greens, `#ef4444` danger reds and white labels.
All mapped to `--accent-primary`, `--success`, `--danger`, `--warning`,
`--info`, `--on-accent`. Hardcoded hex left in the image sheets: 23, all
decorative (model-badge pastel hues, SVG escape arrows).

### Image page after unification

| measure | light | dark |
|---|---|---|
| body | rgb(233,230,223) | rgb(10,10,15) |
| .ig-layout | rgb(228,224,214) | rgb(26,26,36) |
| .ig-panel | rgb(247,245,240) | rgb(21,21,31) |
| page heading | rgb(13,148,136) = --accent-primary | rgb(212,212,212) |
| stuck-dark areas | 0 | 5 (correct) |
| low-contrast text | 0 | 0 |

## T2.1 Stage C - remaining pages verified (2026-11-01)

Verification pass only: no CSS, JS or markup was changed. The instruction for
this stage was to verify and correct the existing product, not to redesign it,
so nothing was "improved" on the way past.

### Scope correction

The Stage B/C worklists name six remaining pages including `real-settings`. That
page no longer exists: it was the orphan copy of the settings surface deleted by
the T2.2 sweep (`index.html` holds the live inline modal). The real remaining set
is five pages, and all five are now measured.

### Method

`build/t21/audit-stageC.js` holds the probe verbatim with its rationale. Three
things about it matter more than the numbers:

1. Theme flips go through `theme.js setThemeMode()`, never through a body-class
   toggle, so a missing entry point fails loudly instead of passing silently.
2. Transitions are disabled (`*{transition:none!important}`) and the read waits
   **at least 1200 ms**. Earlier figures in this file used a 600 ms settle, and the
   1200 ms wait is kept as a margin, not because a long transition is known to be
   in flight: there is no `@property` registration anywhere in the CSS, so the 4 s
   custom-property transition declared on `:root` cannot animate these properties
   in Chromium.
3. Background resolution is memoised per element. Recomputing the ancestor chain
   per node blew the 10 s evaluation budget on `index` and `story-tools`.

One tab at a time: each harness page holds an open `/__events` SSE stream
(`harness/shim.js` opens one `EventSource` per page), and six tabs exhaust the
browser's six-connection budget, so later pages never reach
`readyState: complete`. Two "page never loads" results were this, not a page bug.

**Scope of the numeric gate.** The probe only evaluates text nodes whose element
font-size is 16 px or larger (`audit-stageC.js` skips `fontSize < 16`), and its
offender test is `|dL| * 255 < 60` on relative luminance - a documented
luminance-difference heuristic, not a contrast ratio and not a WCAG or
accessibility result (a delta of 60 sits near 1.6-2.0:1). Foreground alpha is not
composited: colours below `a = 0.3` are skipped, everything above is scored as if
opaque. In particular the 14 px primary-button labels (`scenario-header.css:32`,
`input-components.css:39`) are outside the numeric gate; their appearance was
assessed visually only.

### Measured result (transitions disabled, 1200 ms settle, text >=16 px only)

| page | light canvas | light low-contrast | light stuck-dark | dark low-contrast | dark stuck-dark |
|---|---|---|---|---|---|
| index | rgb(233,230,223) = --app-bg | 0 | 0 | 0 | 13 (expected) |
| scenario | rgb(233,230,223) | 0 | 0 | 0 | 5 (expected) |
| lorebooks | rgb(233,230,223) | 0 | 0 | 0 | 6 (expected) |
| story-tools | rgb(233,230,223) | 0 | 0 | 0 | 11 (expected) |
| audio-settings | rgb(233,230,223) | 0 | 0 | 0 | 4 (expected) |

Column notes: "low-contrast" is the heuristic above, applied to text >=16 px only;
"stuck-dark" counts matching elements (own background opacity > 0.5, area
> 6000 px^2, resolved luminance < 25), so nested panels can count more than once
and the numbers are element counts rather than distinct regions.

Visual pass per page covered the filled-child case the numeric gate does not
reach: the teal primary buttons checked on scenario, lorebooks, story-tools and
the audio footer render white labels, so the world-page "dark grey on teal"
failure does not repeat. Qualification: that is a per-button visual observation,
not a token sweep. `--on-accent` is what the components checked use
(`scenario-header.css:248-252`, `theme-light-components.css:159-163`), but not
every filled button does - `theme-light-components.css:145` labels `.send-button`
with the literal `#FFFFFF` and `st-memories.css:141` labels
`.memories-manage-btn` with `color: white`. Both render white today, so the claim
is "these buttons render white", not "every filled button carries the token".

### Observed, deliberately not changed

- **`Manage All Memories` is cyan, not teal.** It uses `--accent-secondary`,
  which is `var(--theme-secondary)`, defined as `#06b6d4` in both modes
  (`theme.css:30,77`; `theme-light.css:18,47`). It is therefore deliberately
  constant cyan in light and dark - not a token that changes value with the theme,
  and not a dark value leaking into light. A secondary accent doing its job;
  changing it would be a redesign, so it was left alone and recorded instead.
- **Top-bar collision at narrow widths.** At a 689 px viewport
  `.loading-progress-container` (x=437, w=236) sits on top of `.menu-bar`
  (x=181, w=492), so `NO MODEL LOADED` runs into the window edge and the menu
  labels are overprinted. Positions and colours are identical in both modes, so
  this is the already-open layout item (Stage A, "Still open" 3), not a theme
  defect. Not fixed here; it belongs to the layout pass.
- **`audio-settings` body has no background in dark** (`rgba(0,0,0,0)`), while
  every other page paints `rgb(10,10,15)`. The overlay and its own containers
  cover the viewport, so no white area was observed, but nothing guarantees that
  during load. Recorded for the layout pass rather than patched here.

### Corrected claim (mine, retracted)

An intermediate reading of this pass reported the memories flyout as "stuck light
in dark": panel `rgb(247,245,240)`, inputs `rgb(251,250,246)`, and three
low-contrast labels. The later measurement invalidates that reading. At a
800-1200 ms settle the same flyout measures panel `rgb(21,21,31)` = dark
`--surface-secondary`, inputs `rgb(21,21,31)`, label `rgb(212,212,212)`, and 0
offenders at the probe's heuristic. The flyout is clean.

What is *not* established is why the first read saw light colours. The cause
offered at the time - theme.js's 4 s custom-property transition still running -
is unsupported and contradicted: there is no `@property` registration anywhere in
the CSS (0 matches), so those custom properties are not animatable in Chromium,
and the probe overrides `transition` with `!important` on every element before
flipping. **Cause: Unestablished.** No replacement explanation is offered here.
What stands is the later measurement and the retraction itself.

### Ownership defect found in the evidence path (open, not fixed, not safe)

`build/t21/harness-stageC.log` records `[proc] reaped 1 orphaned model
process(es)` at 20:51:06 during harness startup. That reaping event is
artifact-backed. The identity of the victim - `sd-server.exe` PID 24780 as a child
of the running app `howl-chat.exe` PID 9252 - is author-observed process
inspection from that session; the saved log names no PID, so the attribution is
not independently established by the artifact.

The mechanism is source-derived, and it is broader than this record first said.
`internal/proc/job_windows.go:110-136` (`ReapOrphans`) enumerates
`llama-server.exe` / `sd-server.exe` through CIM and force-stops every match whose
`ExecutablePath` begins with the project root. It consults no parent PID, no child
list and no recorded owner, so it does not establish orphan status at all - it
kills by executable path, whoever owns the process. The earlier wording,
"orphaned was decided from recorded state", was wrong; there is no recorded-state
check to blame.

Consequences, stated plainly:

- The running app was left idle with no image child until a Reload; the audit did
  not touch generation, so no work was lost.
- This is the same ownership class as video `stopServer`, and worse: ownership is
  inferred from where an executable lives rather than from who started it.
- Neither the video defect nor "shared infrastructure is safe" may be claimed on
  the strength of this stage. Parent-identity verification belongs with the
  diagnostics/ownership work, not here. **This stays open and separate from the
  theme stage.**

### Evidence limitations

- **No per-page result artifact was persisted.** `build/t21` holds the probe, the
  harness log and older worklists, but no JSON or CSV of the five pages, so the
  table above cannot be re-derived from disk. The counts and canvas values are
  author-reported probe output, not saved artifacts.
- **Screenshots were inspected inline only.** Every save attempt was refused by
  the capture path ("use a new workspace-relative .png path"), so no image of
  these five pages exists in the repo. The visual checks are author observations.
- **The session was not read-only with respect to application data.** The harness
  runs the real backend against the real profile; `harness-stageC.log` shows
  `SaveChatData` (20:55:29) and `SaveChat` (20:56:00). No source file was modified
  by the pass, but chat state was written through the app's own bindings.

### Still open after Stage C

1. Delete the parallel override sheets (`theme-light-components`, `theme-overrides`,
   per-feature `*-light-theme.css`) - measured before removal, as before.
2. Top-bar/menu collision and the fluid-layout pass (layout, not theme).
3. `audio-settings` dark body background, decided in the layout pass.
4. Process reaping must verify parent identity (ownership work).
5. Stale `real-settings.html` references, cleanup only: `frontend/dist/js/theme.js:231`
   (shipped comment), `build/t21/audit-pages.txt` ("PAGES: 9" plus its row) and
   `build/t21/audit-worklist.txt` (its row).
6. Persist per-page audit output (results and screenshots) so the next theme pass
   is reproducible rather than testimony.

### Closure

Independent review: **PASS WITH FINDINGS, no blocking findings** (2026-11-01).
Stage C theme verification is closed on that verdict, with no blocking findings
and nothing to repair before closure.

The review findings are preserved here as recorded limitations, not converted
into passes: the numeric gate covers text >=16 px and is a luminance heuristic,
the flyout cause is Unestablished, the visual checks are unpersisted observations,
and the process-ownership defect stays open as cross-cutting work outside this
stage. No product code, test, configuration or frontend asset changed in this
stage or in this correction.

## The real reason pages stayed dark (root cause, fixed 2026-11-01)

Every page looked themed and every token was correct - and six of the nine pages
still opened dark. Two independent restore paths were fighting, and only one of
them could ever win.

1. `theme.js` restored from `window.settingsState.current.general.theme`. Only
   **index, audio-settings and real-settings load settings-state.js**. On the
   other six the store is absent, so the restore found nothing and the legacy
   fallback read a key that had already been migrated away and deleted.
2. `menu-view.js` then restored the theme again, from a **third** key
   (`howl-chat-theme`) with a hardcoded `'dark'` fallback. On any page without
   the store that is always `'dark'`, so it overrode whatever `theme.js` had
   done - and wrote `dark` back into the preference.

Chat flipped because it is one of the three pages that loads the store, so both
paths agreed there.

### Fix (single owner, durable mirror)

- `theme.js` is the only writer and only reader of `howl-theme-mode`
  (localStorage). Every mode change mirrors there synchronously.
- Restore order: unified store first, durable mirror second, with `'system'`
  resolved per load via matchMedia. Restoring no longer writes back, so a
  'system' preference is not collapsed into a fixed mode.
- `menu-view.js` restore block removed (its comment now says why); it keeps only
  the live View-menu handler, which delegates to `theme.js`.
- Init runs as soon as `theme.js` is reachable instead of waiting for
  `DOMContentLoaded` (the script sits near the end of the body on most pages).

### Evidence

| step | body class | canvas |
|---|---|---|
| index, View > Light | theme-light | wrote mirror=light |
| world.html (no settings store), navigation only | theme-light | rgb(233,230,223) = #E9E6DF |
| scenario.html (no settings store), navigation only | theme-light | rgb(233,230,223), content matches canvas |
| back to dark, then navigate | theme-dark | rgb(10,10,15) |

Writers of `howl-theme-mode`: 2 (both in theme.js). Readers: 1. Restore paths
outside theme.js: 0. `node --check frontend/dist/js/theme.js`: OK.

### Audit lesson

My harness probe called `setTheme` when present and otherwise toggled the body
class itself. The fallback silently produced a passing result on exactly the
pages that were broken. Any probe must exercise the app's own path and fail loudly
when the path is missing.

## Live step preview during generation (verified 2026-11-01)

It exists, it was already wired end to end, and it now has direct evidence
instead of assumption.

Chain: `buildServerArgs` passes `--taesd <file>` + `--taesd-preview-only` when the
bundle has a TAESD decoder -> the job poller calls `engine.GetProgress()` every
500 ms -> the engine merges the stdout step (sd-server rewrites its progress line
with CR) with the frame from `/sdapi/v1/progress` (`current_image`) -> the frame
becomes `job.preview_image` -> `image-generator-core.js` paints it into
`#ig-preview-area` on every poll.

### Backend probe (real run, GTX 1080, SDXL bundle)

Launched sd-server with exactly the engine's flags and ran 512x512 / 8 steps:

- server log: `loading tae from taesdxl.safetensors`, `using TAE for preview`
- `/sdapi/v1/progress` returned **42 frames** across 68 polls, **8 distinct**
- sizes 1360-2884 bytes, raw base64 starting `/9j/` (JPEG, no `data:` prefix)
- progress reported 0.125 -> 0.875 -> 1.0 while frames arrived
- final image 695,620 bytes in 17.7 s

So stable-diffusion.cpp does support intermediate previews and needs no sidecar.
The probe is kept at `build/preview_probe.py`.

### Two snags found

1. `showPreviewImage` sent every frame as `image/png`. sd-server's frame is JPEG;
   the renderer now sniffs the base64 magic (`/9j/`, `UklGR`) instead.
2. Nothing constrained the preview `<img>`, so a 64-128 px frame painted at
   natural size in the middle of the canvas with black bands around it. Scoped
   rule `.ig-canvas img.ig-live-preview` now fills the canvas with
   `object-fit: contain`; result images keep their existing rules.

### What still needs the desktop app

The Go->UI hop (`GetImageJobStatus` carrying `preview_image` into that poller)
cannot be exercised from the browser harness, which has no Wails bindings. The
poller and renderer are verified on their own, and the backend is verified live;
a click-to-generate run in the real exe is the remaining end-to-end check.

## T2.2 · Settings surface sweep (first pass, evidence below)

Swept the settings UI the way the ghost-element hunt requires: every control
checked against its binding, and every binding checked against a reader.

### What was there

The settings modal existed **twice**: the live copy inline in `index.html`
(the one users see) and a standalone `real-settings.html` page that nothing
referenced. They had already drifted — the orphan had an `android-companion`
checkbox that exists nowhere in code, the live copy had a whole Network panel
the orphan lacked. The orphan has been deleted; `index.html` is the only copy.

Live modal before: 5 panels, 39 controls, 58 setting-cards.
Live modal after: 4 panels, 20 controls, 17 setting-cards.

### Removed because nothing could ever use them

| Removed | Why |
|---|---|
| Network & Remote panel (8 controls) | No listener, no state key, no reader anywhere: endpoint, both timeouts, retries, delay, proxy host/port, remote model |
| `android-companion` (orphan only) | Disabled placeholder; no code path in the repo |
| 6 storage-location cards | Fields were `readonly` (no input listener) and the Browse buttons used a `webkitdirectory` file input whose `first.path` does not exist in WebView2, so they wrote a bare folder name into keys nothing reads. Storage is fixed to AppData\HOWL_Chat in this build |
| `ui-scale` | Its one CSS consumer is a single `font-size` on one element — it did not scale the interface |
| `disable-telemetry` | Set `window.__telemetryDisabled`, which nothing reads; the app has no telemetry |
| Manage / Reset Permissions | `prompt()` theatre writing `privacy.permissions`, which nothing reads |
| `default-tts-voice` | Options were invented (`male-1`, `female-2`); the real Piper voice list lives on the Audio settings page |
| `tts-autoplay-message` | Duplicate of `tts-autoplay`. The General one (`general.ttsAutoplay`) is the one audio-service/chat-stream read; this one was never read |
| `custom` theme option | `theme.js` resolves only light/dark/system; choosing it produced `theme-custom` and no styles |

Dead code behind them went too: `_browseDirectory`, `_resetLocation`, the twelve
`browse*/reset*` one-liners, `managePermissions`, `resetPermissions`,
`applyUIScale` — 4.3 KB from settings-actions.js, plus their bindings.

### Two real defects fixed, not just dead paint

1. **Settings' theme picker bypassed the single theme owner.** `applyTheme()`
   reimplemented the body-class swap and wrote the settings store directly,
   skipping the `howl-theme-mode` mirror the other eight pages restore from.
   It also called `setThemeMode(theme)` with no `persistStore`, so the choice
   was never saved, and passed `'system'` straight into a function that only
   accepts light/dark (it logged an error and did nothing). It now resolves
   `system` itself, calls `setThemeMode(effective, false)`, and stores the
   preference separately.
2. **"Clear All Data" cleared localStorage only.** Chats, worlds, characters,
   scenarios and lorebooks are JSON files under AppData, removed by the
   backend — so the button deleted almost nothing while claiming everything
   was gone. It now delegates to `resetApplication()` (backend `WipeAllData`,
   IndexedDB, localStorage, with its own two-step confirmation).

### Verified

- 9 pages: **0 duplicate element ids**.
- Live modal: **20 controls, 0 unbound**; 59 inline handlers, **all resolve**
  to a defined function.
- Modal region div balance 74/74 after surgery (the page-wide 628/629 imbalance
  is pre-existing — HEAD reads 746/747 — and is unrelated).
- Theme path through the real select: dark → body `theme-dark`, mirror `dark`,
  store `dark`; light → all three `light`; system → resolved to a concrete mode
  with the preference kept in the store.
- `node --check` clean on all three touched files; exe rebuilt (20,314,624 bytes).

### Found in passing, not fixed here

- **JS was loaded with no cache-busting at all** — only CSS carried `?v=`. Every
  `<script src="js/*.js">` now has `?v=15` across all 8 pages (277 tags). This
  is why the first applyTheme test appeared to fail: the tab was running the
  pre-fix file.
- **20 function names are defined in two or more JS files**, so load order
  decides which one wins silently (`escapeHtml` in 6 files, `setTheme` in 2,
  `toggleSidebar`, `deleteChatSession`, `updateModelDisplay`, …). Needs its own
  pass; fixing it by load order would be a guess.
- Storage is not relocatable. If that is wanted it is an addition: a native
  folder picker (`runtime.OpenDirectoryDialog` already used for the image models
  dir) plus backend storage that actually consumes the chosen path.

## T2.2 · Duplicate function hunt (second pass)

Twenty function names were defined in more than one JS file. Defining them is
not the same as them mattering: a duplicate only has power when two files that
define the same name load on the same page, because the later one silently
wins. Measured that split first.

**Before:** 20 duplicated names — **10 real shadows** (co-loaded, last wins) and
10 cross-page twins (never co-loaded, organisational only).

**After:** 20 → 10 duplicated names — **0 real shadows**.

### The ten real shadows, and what each one was doing

| Function | Who won | What that meant |
|---|---|---|
| `setTheme` | menu-view.js over theme.js on 6 pages | A second theme implementation, exactly the class of bug T2.1 removed. It also left a dangling `if (theme === 'system')` listener that referenced a deleted variable, so every one of those six pages threw `ReferenceError` on load |
| `escapeHtml` | 4 different bodies across 6 files | The textContent versions do not escape quotes; the regex one did. Whether `escapeHtml` was attribute-safe depended on script order on the page |
| `deleteChatSession` | sidebar-session.js over chat-session.js | 2.6 KB of dead implementation, including a `confirm()` path that ignored the app's own delete dialog |
| `showLoading` | rag-knowledge.js over story-tools-ui.js on story-tools | Two incompatible meanings: `showLoading('Uploading...')` (text) and `showLoading(true)` (toggle). Whichever file loaded last decided. On story-tools, `showLoading(false)` could show the overlay with the text `false` and never hide it |
| `setLoadingStatus` | error-handler.js over loading-progress.js | The wrapper forwarded to `window.setLoadingStatusOriginal`, which nothing ever assigned, and skipped the progress-bar state classes entirely. Progress states (loading/ready/error) have been dead on index |
| `toggleMetadata` | image-generator.js over image-generator-core.js | Near-identical; the loser lacked the "no metadata" status message |
| `updateModelDisplay` | error-handler.js over chat-model.js | Identical behaviour, so nothing broke — it was pure noise in the model-load path |
| `viewCharacterMemories` | character-memories.js over sidebar-characters.js | Real implementation beat a delegating wrapper that also worked; one of the two was unnecessary |
| `updateContextDisplay` | sidebar-panels.js over sidebar-context.js | sidebar-panels.js re-declared shims on purpose, with a comment saying the other modules "overwrite context-display.js shims". A workaround for a workaround |
| `updateLorebookDisplay` | sidebar-panels.js over sidebar-context-lorebooks.js | Same re-shim clash; 2 KB of real render code in the loser that no page could reach |

### Fixes

- `escapeHtml`: one owner, `js/html-utils.js`, loaded first on index, story-tools
  and lorebooks. Escapes `& < > " '` so it is safe in text and in attributes.
- `setTheme`: theme.js owns it again, resolves `system` itself, delegates to
  `setThemeMode`, and stores the user's choice (`system` stays `system` so
  Settings shows what was picked) — the deleted copy was the only one that saved
  the preference.
- `system` theme now follows the OS while the app runs: the listener moved into
  theme.js, where the preference lives.
- `showLoading`: one implementation in story-tools-ui.js that accepts both call
  styles; both already targeted `#loading-overlay`.
- `setLoadingStatus`, `updateModelDisplay`, `deleteChatSession`,
  `toggleMetadata`, `viewCharacterMemories`, `updateContextDisplay`,
  `updateLorebookDisplay`, `updateWorldDisplay`, `updateScenarioDisplay`: one
  owner each, losing copies deleted. The four panel updaters are owned by
  context-display.js (which instantiates ContextDisplay); story-tools.html loads
  the sidebar modules without that orchestrator and has no world, scenario or
  lorebook panel, so it gets explicit no-ops instead of undefined functions.

### Verified (runtime, not just grep)

- index.html: every consolidated global is a function; `setTheme('light')` moves
  body class, the `howl-theme-mode` mirror **and** the stored preference; panel
  refresh routes to `contextDisplay` and does not throw.
- story-tools.html: all four panel updaters exist and are callable with no
  context panel present; `showLoading('text')` and `showLoading(false/true)`
  both behave on the shared overlay; `escapeHtml` comes from html-utils.
- lorebooks.html: `escapeHtml` present and quote-escaping; the seven functions
  that were nearly lost in this pass are all intact.
- All 100 JS files parse; exe rebuilt at 20,310,016 bytes; assets `?v=18`.

### Near miss worth recording

The first deletion run used a brace counter that desynced on a regex literal
containing a quote (`/'/g`). It removed eight functions instead of one from
lorebooks.js (`populateForm`, `editLorebook`, `deleteLorebook`,
`exportLorebooksForBackend`, …) in a single 3.7 KB cut. The diff showed zero
insertions against HEAD, which proved the file was clean and the restore was
exact. `build/drop_fn.py` now snapshots the file, checks the syntax, compares
the function inventory (exactly one name may disappear), and restores itself on
any mismatch. It caught the same desync on the retry and refused.

### Not touched: the ten cross-page twins

`toggleDropdown` x3, `toggleLorebookSelection` x3, `resolveLorebooksFromRefs`
x3, `lorebookKey`/`lorebookDisplayName` x3, `editCharacter`, `loadAvailableCharacters`,
`removeImage`, `toggleSidebar`, `updateSelectedLorebooksList`. None of these pairs
ever load on the same page, so no page's behaviour depends on load order. They
are copy-paste between the character/scenario/world page families, and merging
them means parameterising three sets of near-identical DOM flows — a refactor with
real regression risk and no bug behind it. Listed rather than guessed at.

### Correction: the first scan was declaration-only

The "0 real shadows" figure above counts `function name(` declarations only.
That undercounts, and it hid a live conflict: sidebar-context-world.js declares
a bare `async function updateWorldDisplay()` and loads at line 108, after
context-display.js at line 100, so it silently replaces that module's
`window.updateWorldDisplay` shim. Same for `updateScenarioDisplay`.

`build/dup_scan2.py` now counts declarations, `window.NAME =` assignments, bare
assignments and `ContextDisplay.prototype.NAME`. Corrected numbers:

- **33 duplicate names** (not 20)
- **22 names defined by two files that load on the same page** (not 10)
- of those, the ones that are genuinely the *same namespace* (global vs global,
  not a prototype method beside a global) are still open:

| Name | Co-loaded on | Conflict |
|---|---|---|
| `updateWorldDisplay` | index | sidebar-context-world.js declaration beats context-display.js shim, and the two renderers differ (the sidebar one also writes `#world-description-display`) |
| `updateScenarioDisplay` | index | same shape |
| `updateCharacterPanel` | index | four definitions: prototype, context-display shim, sidebar-characters declaration, sidebar-panels declaration |
| `deleteMemory` | index | story-tools-data.js vs story-tools-ui.js, both declarations |
| `handleBackendError` | index | error-system.js declaration vs error-handler.js window assignment |
| `loadChatData` / `saveChatData` | index | sidebar-session.js vs sidebar-panels.js (re-shim again) |
| `loadWorlds` | story-tools | world-grid.js vs sidebar-context.js |
| `resetToDefaultState` | index | context-display.js window vs sidebar-state.js declaration |
| `switchMemoriesTab` | story-tools | story-tools.js declaration vs rag-knowledge.js window |
| `toggleLorebookSelection` | story-tools | world-dropdowns.js vs sidebar-context-lorebooks.js |
| `logChatContract` | index | sidebar-state.js declaring and also assigning window |
| `pendingNarratorEvents` | index | a *state variable* assigned in both chat-constants.js and chat-input.js, not a function |

Two of the twenty-two need no action and should not be "fixed": `selectWorld`,
`selectScenario`, `openWorldSelector`, `openScenarioSelector`,
`removeCharacterFromChat` and `createCharacterItem` pair a
`ContextDisplay.prototype` method with a global of the same name. Those are
different namespaces and the app calls them deliberately. Fixing them by
deletion would be a guess dressed as a cleanup.

`sidebar-panels.js`'s guarded no-ops also show up in this scan; they only run
when `window.contextDisplay` is undefined, so they do not shadow anything on
index.

Next pass should merge the two renderers (world/scenario) rather than delete
either one, then take the sidebar re-shim family (`loadChatData`, `saveChatData`,
`updateCharacterPanel`) together, as one unit, because they exist for the same
reason.

### Final sub-pass: the twelve same-namespace conflicts -> zero

Each of these had one name reachable by two or more files on the same page.
Resolution rule: keep whichever one the page could actually reach, merge any
behaviour the loser carried, delete the rest. No overrides added.

| Name | Resolution |
|---|---|
| `loadChatData`, `saveChatData` | Deleted the one-line delegating copies in sidebar-panels.js. sidebar-session.js owns the load sequence and the save contract; the copies only forwarded to it |
| `logChatContract` | Deleted the wrapper in sidebar-panels.js. sidebar-state.js declares it and exposes it; the wrapper existed solely to re-expose it |
| `resetToDefaultState` | Deleted the global shim in context-display.js. sidebar-state.js owns the state reset; the shim was already being overwritten on every page |
| `updateWorldDisplay`, `updateScenarioDisplay` | context-display.js stopped publishing globals for these. The sidebar modules are the richer renderers (they also write `#world-description-display` / `#scenario-description-display`) and they load later, so publishing a second global only created the appearance that the module version worked |
| `updateCharacterPanel` | Deleted the dead 180-line renderer in sidebar-characters.js and the re-shim in sidebar-panels.js. context-display.js is the owner; the sidebar copy used `.character-list` and had not rendered a panel since the re-shim landed. Its 5 call sites and its API export now go through `window.updateCharacterPanel` |
| `deleteMemory` | Merged. story-tools-ui.js now takes `(memoryId, memoryTitle)` *or* nothing (editor button calls `deleteMemory()`), uses the passed id or `storyToolsState.currentMemoryId`, calls `MemoryStore` **and** the backend, then re-renders both lists. Deleted the no-argument copy in story-tools-data.js |
| `switchMemoriesTab` | Merged into story-tools.js (which owns the tabs), including the knowledge-panel show/hide that lived in rag-knowledge.js's wrapper. Deleted the wrapper and its captured "original" |
| `handleBackendError` | error-handler.js used to monkey-patch error-system.js after load to add recovery. error-system.js now calls `onErrorReported()` itself; the patch block is gone. Recovery no longer depends on load order |
| `toggleLorebookSelection` | These were two different functions sharing a name: the world/scenario dropdowns take `(lorebook, element)`, the chat-level one takes an id and syncs the backend. Renamed the chat-level one to `toggleChatLorebook` (and its inline onclick and `sidebarContext` entry) |
| `loadWorlds` | Deleted the stub in sidebar-context.js that only logged an array. world-grid.js owns loading, image rehydration and grid rendering. The stub also shipped a `loadWorlds();` call that would have shadowed it on story-tools |
| `pendingNarratorEvents` | Left alone deliberately: one array, initialised in chat-constants.js, guarded in chat-render.js, drained and cleared in chat-input.js. Three mutators, one owner, no shadow |

### Two live bugs this found

1. **`window.sidebarContext` never existed on story-tools.html.** That page loads
   the sidebar modules without context-display.js and loads sidebar-context.js
   *before* sidebar-panels.js, so its API object referenced globals that did not
   exist yet — a `ReferenceError` that aborted the assignment. The page then ran
   with no orchestrator API at all. Every entry is now resolved at call time
   (`_callGlobal(name)`), 19 keys before and after.
2. **The world lorebook dropdown on story-tools matched the wrong function.** Two
   `toggleLorebookSelection` definitions loaded on that page with incompatible
   signatures; the chat-level one won, so a dropdown row was passed as an id and
   the toggle returned silently. Renaming removed the collision; the dropdown
   keeps its 2-argument form (verified: `toggleLorebookSelection.length === 2`).

### Verified

- index.html: all seven owners resolve; `updateWorldDisplay` still contains the
  description write (i.e. the rich renderer survived); `resetToDefaultState`
  still resets `currentChatData`; `loadChatData` still contains
  `GetActiveChatMeta`; the character panel renders its 3 characters; panel
  refreshes all return `ok`; `sidebarContext` exposes 19 callable entries.
- story-tools.html: `sidebarContext` is an object with 19 keys and its lazy
  entries call cleanly; `switchMemoriesTab('knowledge')` reveals the knowledge
  panel and hides the memory cards, `'all'` restores them; `deleteMemory()` with
  no selection reports instead of throwing; console is free of `ReferenceError`.
- All 100 JS files parse; exe rebuilt at 20,310,528 bytes; assets `?v=20`.

### What remains on the scanner, and why it is not a defect

13 names still appear "co-loaded", and the assignment-vs-declaration conflict
list is now empty:

- `selectWorld`, `selectScenario`, `openWorldSelector`, `openScenarioSelector`,
  `removeCharacterFromChat`, `createCharacterItem`, `updateCharacterPanel`,
  `updateContextDisplay`, `updateLorebookDisplay`, `updateWorldDisplay`,
  `updateScenarioDisplay` — a `ContextDisplay.prototype` method beside a global
  of the same name. Different namespaces; the app calls both on purpose. The
  prototype world/scenario renderers are now unreachable from the global path and
  are the natural deletion target when context-panels.js is next touched.
- `pendingNarratorEvents` — one array, three intentional mutators (above).
- `activeLoRAs` — three `window.activeLoRAs =` assignments on image-generator.html
  (image-generator-core.js, image-generator-ui.js, image-workflows.js). This is a
  genuine state-ownership problem of the same shape, but it belongs to the image
  page and therefore to T2.3; listed here rather than half-fixed blind.

## UI pass · Sampler settings: narrow popover -> full-screen sheet

Reported as "beautiful but cramped", wanting one place to see more options without
scrolling. Measuring it found a structural fault, not just tight styling.

### The fault

`index.html` opened the Advanced Controls category header and then the cards
directly, with **no `<div class="settings-grid">` wrapper**. Nine cards were
therefore children of the flex-column `.flyout-content` and stacked one per row
at full width, and the stray `</div>` at the end of the flyout was closing the
sheet itself rather than a grid. That is also the source of the page-wide div
imbalance that has been on the books for a while: `index.html` is now
**629 open / 629 close** (was 628/629).

### What changed (scoped to `#sampler-settings-flyout` only)

- Sheet, not popover: `width: min(1560px, 96vw)`, height `calc(100vh - 24px)`,
  header pinned, only the content pane scrolls. Also overrides the
  `ss-responsive.css` cap that pinned every `.flyout` to `max-width: 400px`,
  which was cropping the sheet on small windows.
- Grid: `repeat(auto-fit, minmax(250px, 1fr))` -> 5 columns at 1600px, 1 at 520.
- Reset moved to the card's top-right corner; the "Enabled" toggle joined it
  there instead of owning a row (label kept for screen readers).
- The three "Good for / Example / Warning" lines became a hover/focus popover, so
  nothing reflows and no information is lost. On touch (`hover: none`) they stay
  inline.
- Descriptions capped at two lines so one long sentence cannot stretch a row; the
  full sentence is attached as a tooltip by `sampler-settings.js` (verified: only
  one description is actually clamped, the thinking-models one, and it keeps its
  full text).

### Measured at 1600x1000 (dark theme)

| | before | after |
|---|---|---|
| content height | 2074 px | **917 px** |
| visible without scrolling | 807 px (overflowing) | **917 px (0 overflow)** |
| columns | 2 (advanced cards: 1 full-width) | **5** |
| card height | up to 220 px | 104-136 px |
| cards inside a grid | 10 of 19 | **19 of 19** |

At 520x820: sheet 499 px wide inside the viewport, 1 column, toggle inline,
page scroll width identical with the sheet open (523) and closed (523), so the
residual 3 px page overflow is pre-existing and unrelated.

Exe rebuilt at 20,316,672 bytes; assets `?v=27`. Note: the browser screenshot
surface is narrower than an emulated wide viewport, so captures show the
one-column state; the five-column result is verified by computed layout.

### Deferred, deliberately

`.setting-card` and `.category-header` are still overridden from two layers
(`theme-settings.css`, `ss-light-theme.css`, `theme-overrides.css`). They now
agree, but the clean end state is one sampler stylesheet; that is a T2.3-scale
reshape, not a drive-by.

## UI pass · Model settings: 760px modal -> full-screen sheet

Same request as the sampler sheet, applied to `#model-settings-flyout`. Unlike the
sampler this one had no markup fault: all four categories already had their
`.settings-grid` wrapper. The problem was density, not structure.

### What changed (scoped to `.model-settings-flyout` only)

- `ms-flyout.css`: flex-column sheet, `min(1560px, 96vw)` x `calc(100vh - 24px)`,
  pinned header, only `.flyout-content` scrolls (`max-height: none`, `flex: 1`).
- `ms-categories.css`: `repeat(auto-fit, minmax(250px, 1fr))` with 10px gaps,
  cards at `8px 10px 7px` and 10px radius, category headers at 14px.
- `ms-controls.css`: the three "Good for / Example / Warning" lines are now a
  hover/focus popover (was ~50px of every card), descriptions clamped to two
  lines with the full sentence kept as a tooltip via `model-settings.js`, and
  controls tightened (64px number inputs, 6px selects, 64px minimum textarea).
- `ms-reset.css`: Reset is absolutely positioned in the card corner, so it no
  longer adds a row to all 16 cards that have one; labels reserve 56px so they
  cannot run under it.
- Toggles were left inline on the two cards that have one: a bare 48px switch in
  a corner is less legible than the labelled row, and there was height to spare.

### Measured at 1600x1000, before/after on the same page

The "before" column is measured by re-injecting the old geometry as a temporary
stylesheet in the live page (760px sheet, 230px columns, inline details, Reset row
placeholders), not estimated.

| | before | after |
|---|---|---|
| sheet | 760 x 849 (85vh) | **1536 x 976 (full height)** |
| columns | 2 | **5** |
| card height | 169-245 px | 84-160 px |
| content | 2890 px in 770 px | **897 px in 897 px** |
| overflow | 2120 px of scrolling | **0** |
| cards / controls | 21 / 28 | 21 / 28, all in one view |

### Other checks

- Hover popover verified through `:focus-within`: `display: flex`, 3 detail lines,
  270px wide, positioned over the row below rather than reflowing it.
- Header stays pinned while the content pane scrolls; sheet sits inside the
  viewport at 1600x1000 and at 473x799 (one column, 454px wide, nothing clipped).
- Light theme sampled live: popover `rgb(239,235,227)` on border
  `rgb(207,201,190)`, detail text `rgb(87,83,75)` on card `rgb(247,245,240)`.
  Nothing in the new rules hardcodes a dark colour; both sheets use theme vars.
- Assets `?v=29`, exe rebuilt.

### Known limitation

The browser screenshot surface here is ~480px wide regardless of the emulated
viewport, so captures render the one-column layout and cannot show the five-column
result. The five-column claim rests on computed layout, which is what is recorded
above.

## T2.2 · Ghost-ID inventory (evidence revision, read-only, 2026-11-02)

Nothing under `frontend/dist` was modified. No build, test or app launch. This
section **supersedes figures from the earlier inventory pass, which had not
previously been recorded in these project documents**; the 18 ghost ids and their
30 reference sites survive it, the population counts around them do not.

All numbers below describe **this working copy at `1477530`** and are derived from
source text only. They are machine-specific measurements. The structural findings
are candidates about the product, not measurements, and none has been repaired.

### Stage status

**T2.2 is OPEN.** The inventory exercise is complete; the stage is not. Thirty
reference sites remain. This record counts *references*, not read-versus-write
operations, so the stage's acceptance criterion — zero write-to-missing-ID —
**cannot yet be scored**. Nothing below implies that all 30 sites are writes or
runtime failures, and nothing here authorises a repair.

### Units, stated explicitly

A count without a unit is what produced the first, wrong table. Three different
things were being called "ids": element ids passed to `getElementById`, selector
strings that happen to start with `#`, and template-built id prefixes.

| Unit | Count |
|---|---|
| Distinct ids named by literal `getElementById('X')` | **372** |
| Distinct ids present in static page markup (`id="X"` in the 9 pages) | **601** |
| Lookup ids that markup provides | **325** |
| Lookup ids markup does not provide | **47** |
| …of those 47: created by script before use | **29** |
| …of those 47: never created anywhere — **true ghosts** | **18** |
| Distinct ids created in JS (`id="X"` or `.id = 'X'`) | **50** |
| …of those 50: never looked up by id | **21** |
| Distinct selector-based lookups — raw strings / distinct selectors / distinct leading tokens | **33 / 31 / 27** |
| …of those 27 leading tokens: real markup ids / runtime-created | **12 / 15** |
| Reference **sites** for the 18 ghosts | **30** |
| Template-literal lookups `getElementById(\`…\`)` — excluded from every count above | 16 occurrences |

Identities that hold exactly: **325 + 47 = 372**, **29 + 18 = 47**, and
**15 + 15 = 30 sites** (6 ghost ids / 15 sites in files with no load path, 12 ids /
15 sites in loaded files).

A "site" is an occurrence of the id as a literal `getElementById` argument, or
`#id` followed by `'`, `,` or `)`. It is a text measure, not an execution count.

### Why the first numbers could not sum

An earlier pass reported 75 absent ids, split into 18 ghosts + 51 runtime-created +
4 dynamic prefixes + 10 compound fragments — which is 83, not 75. All five figures
are withdrawn, for three independent reasons:

- **51 was impossible on its face.** Only 50 ids are created anywhere in the JS,
  and only 29 of those are both created and looked up. A subset cannot exceed its
  superset.
- **The 14 artifacts were never ids.** Compound entries were whole selector strings
  (`#character-dropdown .world-dropdown-sort`); tokenising yields
  `character-dropdown`, which markup does provide. Prefix entries (`dep-bar-${i}`)
  are template fragments that never match a literal-string extractor.
- **75 counted strings, not ids.** 47 + 28 ≈ 75: the gap matches selector strings
  counted as if they were element ids.

75 cannot be regenerated — the original regex was not retained — so it is recorded
here as unreproducible and replaced by 47 rather than reconciled against it.

Separately, the first re-run reported **0 selector tokens**. That reading was a
failure, not a finding: its `sed` exited with "unterminated `s' command". The
figures above come from an `awk` tokeniser. Do not reuse the broken method.

### The 18 ghosts

Source-level findings. **Thirteen** of the eighteen have a guard, so most are
silent; the other **five are unguarded** and would throw if executed. An unguarded
reference is still a source-level finding, not a demonstrated runtime failure:
none of the eighteen was observed failing, and three of the five sit in a file no
page loads.

| id | sites | file | page loads it? | guard | effect |
|---|---|---|---|---|---|
| `world-image-preview` | 4 | world-manager.js | no | **none** — `.src`, `.style` | would throw if executed |
| `world-image-placeholder` | 2 | world-manager.js | no | **none** — `.style` | would throw if executed |
| `world-image-remove` | 2 | world-manager.js | no | **none** — `.style` | would throw if executed |
| `character-backstory` | 3 | character-manager.js | no | `if (el)` | read no-ops |
| `selected-character-lorebooks` | 1 | character-manager.js | no | early return | read no-ops |
| `scenario-objective` | 3 | scenario-manager.js | no | `if (el)` | read no-ops |
| `chat-input` | 1 | audio-service.js | yes | `if (input)` | transcript written nowhere |
| `character-intro` | 1 | character-form.js | yes | `?.` | `first_mes` always `''` |
| `character-scenario` | 1 | character-form.js | yes | `?.` | `scenario` always `''` |
| `settings-panel` | 2 | error-settings.js | yes | early return | section never injected |
| `sidebar-menu` | 1 | memory-candidates-menu.js | yes | early return | memory menu never mounts |
| `image-status` | 3 | error-handler.js | yes | `if (st)` | image messages never shown |
| `ig-history-status` | 1 | image-generator.js | yes | guarded | filter feedback dropped |
| `world-selector` | 1 | context-data.js | yes | `if (el)` | click never binds |
| `scenario-selector` | 1 | context-data.js | yes | `if (el)` | click never binds |
| `character-selector` | 1 | context-data.js | yes | `if (el)` | click never binds |
| `world-auto-switch` | 1 | sidebar-context-world.js | yes | **none** — `.checked` | would throw; no caller found |
| `scenario-auto-switch` | 1 | sidebar-context-scenario.js | yes | **none** — `.checked` | would throw; no caller found |

### Finding 1 · The audio defect is conditional, so it is not a one-line repair

- Lookup: `audio-service.js:158`, inside `insertTranscriptIntoChat(text)`, guarded
  `if (input)`; sets `.value`, calls `.focus()`, dispatches `new Event('input',
  { bubbles: true })` with the comment "Trigger input event for any listeners".
- The composer: `index.html:1164` `<input ... id="message-input" class="input-field">`.
  No page defines `id="chat-input"`, no script creates that id, and there is exactly
  one lookup of that id — this one. The other occurrences of the string refer to
  the `chat-input.js` module (`index.html:60`, `story-tools.html:50`), its
  `[chat-input]` log tag (`chat-input.js:227`), or comments (`chat.js:9`), which is
  why the name reads as if it were in use.
- Intent is not inferred from the name alone: `chat-input.js:207` reads the composer
  as `document.querySelector('.input-field')`, and `index.html:1157` already uses
  `getElementById('message-input')` for TTS. The transcription code meant the same
  element, by id instead of class.
- **The blocking unknown is upstream.** The caller at `:129` runs after
  `window.go.main.App.TranscribeAudio(tempPath)`, where `tempPath` comes from
  `saveBlobToTempFile`, which returns `URL.createObjectURL(blob)` with the comments
  "Wails runtime doesn't give us direct file access from JS" and "Backend will need
  to handle this". A `blob:` URL is not a filesystem path. If that boundary fails,
  no transcript is ever returned and the id mismatch is unreachable.

Evidence class: id absence is source-derived and certain; the user-visible failure
is an inference from guards and comments, untested. Changing the id first would
risk reporting a repair for a path that still cannot produce input.

### Finding 2 · Two character properties are persisted, written-only, and empty

- Written at `character-form.js:61-62` in `collectCharacterData()`:
  `scenario: getElementById('character-scenario')?.value || ''` and
  `first_mes: getElementById('character-intro')?.value || ''`.
- Restored at `character-grid.js:127` as `set('character-intro', …)`, where the
  helper at `:118` is `const el = document.getElementById(id); if (el) el.value = val`.
  Both directions target an element that does not exist.
- `character.html` offers name, age, role, priority, description, personality,
  background and dialogue-examples. No intro, no scenario.
- **No consumer:** `first_mes` occurs exactly twice in the frontend — this write and
  this restore. Every `.scenario` reference in `js/` is the chat's scenario object
  (`chat-state-bridge.js`, `chat-session.js:37`, `context-data.js:164`), not a
  character record's string field.

Recorded as two write-only properties that persist as `''`. Not deleted, not
mapped to description or personality, no replacement field invented: whether they
are still part of the intended character schema is a product decision.

### Finding 3 · The manager modules look legacy; that is not proof they are unused

Checks run, all negative for loading:

| Mechanism | Result |
|---|---|
| Literal `<script src>` in the 9 pages | absent |
| ESM `import` / `import()` / `require(` | 0 in `dist/js` — all classic scripts |
| `createElement('script')`, `.src = '…js'`, `document.write`, `getScript`, `loadScript`, `insertAdjacentHTML` | 0 |
| Service worker / `new Worker(` | 0 |
| Global registration | top-level `const xManager = new XManager()` + `DOMContentLoaded` init; **no `window.` assignment** |
| Cross-file use of `characterManager`/`worldManager`/`scenarioManager`/class names | only inside the three files |
| Inline `on*=` handler names in all 9 pages → defining file | none resolves to a `*-manager.js` |
| Top-level function names defined in the three files, referenced elsewhere | none |
| Repo-wide (Go, JSON, YAML, scripts, docs) | only self-references plus `T1_NOTES.md:45` |
| Other copies of the files | one each, under `frontend/dist/js` |
| Packaging | `main.go:13` `//go:embed all:frontend/dist` — shipped in the binary, never executed |

Two apparent holes in that audit are artifacts, recorded so they are not mistaken
for defects: the inline-handler map's unresolved rows for `AudioService`,
`document`, `event`, `if` and `window` are attribute-expression fragments, and
`newChat` is real (`chat.js:29` assigns `window.newChat`).

What this supports: *no load path found*. What it does not support: redundancy.
The manager versions carry selector-modal behaviour (`removeSelectedItem`,
`toggleItemSelection`, `confirmItemSelection`) with no live counterpart found by
name, and their element contracts (`world-image-*`) do not match the live forms.
Whether that behaviour was abandoned or never migrated is unestablished, so the
files stay.

### Finding 4 · The form/manager pairs are parallel implementations (candidate only)

`world-form.js` (loaded by `world.html` and `story-tools.html`) and
`world-manager.js` (never loaded) both define `clearWorldForm`, `saveWorld` and
`populateEditForm`, but one is procedural over `window.worldState` and
`world-images-list` while the other is class-based over `this.*` and the ghost
`world-image-*` ids. `character-form`/`character-manager` and
`scenario-form`/`scenario-manager` have the same shape. Overlapping names in
parallel implementations are not evidence that either side is redundant; the live
forms use `character-lorebooks.js` / `world-dropdowns.js` / `scenario-dropdowns.js`
for selection instead. No behavioural equivalence audit was performed.

### Limitations that stay with this record

1. **16 template-literal lookups are unclassified** (`getElementById(\`…\`)`). A
   ghost built from a template is invisible to this method.
2. **The selector-token figure has one valid method.** The earlier run returned 0
   because its `sed` failed; that method must not be reused or cited.
3. **Creation detection covers two patterns only** — `id="X"` and `.id = 'X'`.
   Creations via `setAttribute('id', …)` or a template literal would be miscounted
   as ghosts.
4. **Reachability is static.** No page was loaded, so "never executed" means "no
   load path found", not "observed unexecuted". A Go-side or string-concatenated
   script injection would evade every pattern used.
5. **The prior figures are unreproducible**, so this correction rests on the new
   measurements plus the impossibility argument (51 > 50), not on recovering the
   old pipeline.

### Next decisions, none authorised

1. **Documentation (this entry)** — done.
2. **Audio handoff contract audit** — trace recording → blob → temp-file/backend →
   transcription → composer, read-only, before authorising any id change.
3. **Manager retirement review** — behavioural inventory of the three modules first;
   no deletion until their missing functionality is accounted for.

Checkpoint `1477530` remains the restore point and the working tree was not changed
by this audit.

