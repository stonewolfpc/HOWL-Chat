# HOWL Chat — Debt Sweep

Date: 2026-10-31 (overnight pass)
Method: repo-wide sweep for TODO / FIXME / "not implemented" / "for now" / "coming soon" / simulated / dummy /
hardcoded markers across Go, JS, HTML, CSS and Python; each hit triaged by hand, then fixed or logged.
Every fix below was followed by `go build ./...`, `gofmt`, `node --check` and the relevant package tests.

Sweep totals: 310 raw marker hits, of which ~270 were false positives (`placeholder="…"` attributes in HTML,
template placeholders in `vision.go`, tutorial prose). 34 were real; 10 are fixed.

---

## 1. Fixed in this pass

| # | Where | Was | Now |
|---|-------|-----|-----|
| 1 | `frontend/dist/js/image-generator-core.js` | `hideSpinner() {}` — a no-op, so any failed generation left a spinner in the preview area forever | removes the `.ig-spinner` node; every error path repaints cleanly |
| 2 | `internal/backend/chat/service.go` | `"seed": 987654321, // Placeholder for now` written to every `chat.json` | records `s.currentSeed()` — the pinned seed, or `-1` for random (new `chat/seed.go`) |
| 3 | `frontend/dist/js/menu-edit.js` | Edit → Undo/Redo showed "coming soon" while the shortcut dialog advertised Ctrl+Z/Y | real undo/redo against the focused text field, with a clear message when no field is focused |
| 4 | `internal/backend/llama/http_client.go` | `GetModelInfo`, `Tokenize`, `Detokenize`, `GetContextSize` all returned "not implemented" | real `/props`, `/tokenize`, `/detokenize` calls + `n_ctx` extraction; covered by `http_client_props_test.go` |
| 5 | `internal/backend/memory/service.go` | `GetAutoCaptureConfig` always returned defaults, ignoring the file `SetAutoCaptureConfig` writes | loads `Chats/<id>/auto_capture_config.json`, falls back to defaults only when absent; `config_test.go` |
| 6 | `frontend/dist/js/data-sync.js` | `hasUnsavedChanges()` returned `true` unconditionally ("always assume changes exist") — autosaved every 30s regardless | compares a cheap chat signature (world/scenario/characters/lorebooks/message count/last message) and stores it after each save |
| 7 | `frontend/dist/js/rag-knowledge.js` + `app-rag-upload.go` | upload handed Go a bare filename as a file path, so every browser upload failed | new `UploadRAGDocumentText` binding ingests the text the browser already has; registered in `frontend/wailsjs` |
| 8 | `frontend/dist/js/chat-stream.js` | `streamAssistantResponse()` — a fake character-by-character streamer with no callers | deleted (the real backend streaming path is untouched) |
| 9 | `frontend/dist/js/sidebar-panels-old.js` | unreferenced duplicate panel script containing "Load/Save chat coming soon" stubs | deleted |
| 10 | 6 page headers (`index`, `story-tools`, `character`, `image-generator`, `lorebooks`, `world`, `scenario`) | dead menu items "Advanced / Experimental" (handler only showed a toast) and "Documentation" (no handler at all) | removed, along with the `openAdvancedSettings` function and export |

## 2. Checked and found sound (no change)

- **Layout / resize.** Probed every page at 1024×700, 820×620 and 1920×1040: `scrollWidth == clientWidth` in all
  cases — no horizontal overflow, no scrollbars. The flagged `div.hexagon-grid` is a decorative rotated layer with
  `pointer-events: none` inside a fixed `.ambient-bg`, clipped by `body { overflow: hidden }`. Not a defect.
- **`internal/backend/vision/vision.go`** — "placeholder" hits are template substitution (`{{image}}`), not stubs.
- **`internal/backend/llama/http_client.go` GetLoadedModel / IsLoaded** were already real `/v1/models` calls.

## 3. Open — needs your call (not stubs I can safely finish alone)

1. **Active character assumption** — `internal/backend/chat/service.go:469,546` pass `s.characterNames[0]` as the
   acting character for character-scoped lore. The chat meta already stores `ActiveCharacters`; the missing piece is
   a `SetActiveCharacter` binding called when the sidebar selection changes, plus reading it on load. Correct fix needs
   the UI's notion of "who is speaking" — worth doing deliberately.
2. **Prompt summaries** — `service.go:745,762,787,825` build "simplified summary for now" world/scenario/memory blocks
   and fall back to a generic `"Assistant"` sender. These are prompt-design decisions, not bugs; changing them changes
   model behaviour, so I left them.
3. **`Lorebooks []interface{}`** (`internal/backend/chat/types.go:88`) — intentionally untyped "will address later".
   Needs a typed lorebook reference model and a persistence migration; too invasive for an unattended change.
4. **HuggingFace repo-link support** — `app-chat.go:15` and `internal/image/engine/sdcpp/engine.go:21` TODOs. That is a
   new feature (download from a repo id), not debt to repair.
5. **`GetUsedMemory()` returns 0** for the HTTP client — llama-server exposes no memory endpoint. Real values need
   process RSS sampling from the spawned server; `/props` cannot answer it.
6. **RAG upload covers text-like files only.** `rag.ExtractTextFromFile` (Go-side, path based) still handles richer
   formats; the browser path now sends `File.text()`. PDF/DOCX from the browser would need a Go extraction binding.
7. **EXIF `UserComment`** (`internal/image/metadata/extractor.go:144`) — A1111 writes PNG, so JPEG EXIF parsing is low
   value. Left documented.
8. **`internal/backend/memory/cluster.go`** is not gofmt-clean (pre-existing; untouched by this pass).
9. **Repo-root build litter** — ~25 stale `howl-chat-*.exe` builds, loose `.dll`s and `sd-*.txt`/`server-test*.log`
   dumps sit next to the source. Not deleted: they are yours to keep or purge, and they are already ignored for builds.

## 4. UI polish — measured, not yet applied

The resize contract holds, so the remaining work is visual consistency rather than breakage. Concrete candidates:

- `.hexagon-grid` / ambient layers are duplicated per page; a single shared include would stop six copies drifting.
- Typography: 120 CSS files with ad-hoc `font-size` values; card headers (character / lorebook / world / scenario)
  use different paddings and radii. A shared card token set (`--card-radius`, `--card-pad`, `--font-card-title`)
  would unify them without touching markup.
- 820×620 still shows all three columns on the chat page — dense but not overflowing. A collapse breakpoint for the
  context panel would be the natural next step if you want small-window comfort.

## 5. Verification run in this pass

`gofmt` clean on all touched Go files · `go build ./...` exit 0 · `node --check` on every edited JS file ·
`go test ./internal/backend/memory/ ./internal/backend/llama/ ./internal/image/... ./internal/backend/gguf/` all ok ·
release binary rebuilt to `build/bin/howl-chat.exe`.
