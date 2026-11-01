# T1 stage notes — working checklist

Scratch file for the T1.1 -> T1.4 run. Delete lines when a test closes.

## Ground truth (verified, not assumed)

- Page switch is `navigateTo(page) { window.location.href = page }` — a full document
  reload. Go side + llama-server survive; all JS state dies and must rehydrate.
- `LoadModel` is **async by design**: spawns llama-server, returns in ~2 s, then a monitor
  goroutine polls health for up to 60 s and emits `model:progress` / `model:loaded`.
  So `IsModelLoaded()==false` immediately after `LoadModel` is correct, not a bug.
- Wails `runtime.EventsEmit` calls `log.Fatalf` (os.Exit) when `ctx.Value("events")==nil`
  (v2.12.0 `pkg/runtime/runtime.go:47-61`). Any emit without a live webview **kills the
  process**. This took the UI harness down on the first `model:progress`.
- Fix: `App.emitEvent()` in `app.go` — drops the event when there is no dispatcher.
  11 emits in `app-model.go` rerouted through it. Other files still call
  `runtime.EventsEmit` directly (28 app-wide, production ctx is valid).
- Streaming chat cannot be verified in the harness: tokens arrive via
  `window.runtime.EventsOn('chat:chunk')` and backend-originated events cannot reach the
  shim. Use `SendMessage` (synchronous binding) as the "still answering" proof.

## T1.1 evidence so far

| Check | Result |
|---|---|
| baseline VRAM | 1414 MiB |
| LoadModel return | 1.9 s (async) |
| model ready | 6.0 s, `IsModelLoaded` true, stage `Ready` |
| VRAM loaded | 6366 MiB |
| navigate image-generator.html | VRAM 6407, llama-server pid 25140 alive |
| navigate back to index.html | VRAM 6352, same pid 25140 |
| model name on return | `modelDisplay` = LFM2.5-8B-A1B-heretic.Q4_K_M.gguf |
| chat turn after return | see SendMessage result |

Remaining: sampler persistence across app restart.

## T1.2 · VRAM arbitration — CLOSED (evidence in STAGE_STATUS.md)
Needs: confirm dialog both directions, unload-then-recheck-then-load, VRAM readings
before/after. Extend emitEvent coverage to any remaining direct emits first.
## T1.3 · think tags — OPEN (full spec from user, not implemented)
## T1.4 · per-family sampler defaults — CLOSED (evidence in STAGE_STATUS.md)

## Ghost IDs (T2.2, bundle later)
Functional: `chat-input` (audio-service), `character-selector` (context-data),
`world-image-preview`/`-placeholder`/`-remove` (world-manager; real ids `world-image-input`,
`world-images-list`). Dead: ig-history-status, image-status, settings-panel, sidebar-menu,
scenario-objective, character-intro, character-backstory, character-scenario,
world-auto-switch, scenario-auto-switch, selected-character-lorebooks, scenario-selector,
world-selector.
