# Closure Record — Multi-Directory Model Storage

Date: 2026-10-31
Scope: model directory add / refresh / remove, and strict LLM-vs-image separation in the menus.

Restrictions honoured while producing this record: no feature work, no architectural change, no refactor,
no edits to the delivered fixes, no evidence deleted, and no native verification claimed from mocked bindings.

---

## 1. Verified

### 1.1 Frontend behaviour — real `frontend/dist`, mocked Wails bindings

Method: the shipped `frontend/dist` was served over loopback with one injected mock script that stubs the
Wails bindings; every interaction below was a genuine click/typing action in Chromium, not a code reading.

| Step | Result |
| --- | --- |
| Open Model Browser via "Browser" | modal opens, "No models found. Scan a directory first." |
| Add directory with a pasted path | directory appears in the dirs dropdown, list repopulates, path box clears |
| Click "+ Browse for Directory..." while the picker returns `""` | notice "No folder selected. Paste the path instead." is raised and the path input takes focus |
| Close dropdown | populated list stays populated |

DOM assertions taken alongside the clicks: `items: 8`, `dirOptions: ["E:\\AI GGUF\\GGUF"]`,
`firstItem: "Captain-Eris_Violet-V0.420-12B.Q4_K_M"`, `emptyState: false`.

Limit: the bindings are stubs. This proves UI logic, refresh sequencing and the call into the binding —
not that the bridge resolves or that the real scan returns these rows.

### 1.2 Backend directory scanning

- `ScanLLMDir("E:\AI GGUF\GGUF")` returned **25** GGUF files.
- Classification authority is `general.architecture` from the GGUF header (`gguf.KindOf` / `KindOfArch`),
  replacing the filename/size heuristics that caused the bleed.
- The probe test used for this measurement was deleted afterwards.

### 1.3 Registry contents

- `config/model-scan-dirs.json` (single directory ledger): `{"E:\\AI GGUF\\GGUF": "llm", "E:\\HOWL_Chat\\models": "image"}`.
- `config/model-registry.json`: 29 entries — 25 `llm`, all under `E:\AI GGUF\GGUF`; 3 `diffusion`
  under `E:\HOWL_Chat\models` and `...\models\video`; 1 `lora` (see 4.1).
- Load-time heal removed two mis-filed entries (`Captain-Eris`, `Qwen3VL`) from the image list and persisted
  the correction once.

### 1.4 Error handling

- Empty picker result no longer silently returns (`if (!dir) return;` removed): the user is told and the
  text input is focused.
- Adding a directory that is already registered under another kind is refused rather than re-typed.
- Removal unregisters the directory (`UnregisterDir`, case-insensitive) and never touches the filesystem.

### 1.5 UI refresh behaviour

Adding a directory refreshes all three surfaces — the dirs dropdown, the model list (previously never
repopulated, which is why the folder looked like it had not been added) and the LLM dropdown.

### 1.6 Menu separation

Directory kind (`llm` / `tts` / `image`) is the menu boundary, and every read path gates on it:
`ListImageModels`, `GetModelIndex`, `GetCompatibleAssets`, `GetRegisteredModels`, `modelEntryByName`,
`findCompanionPath`, `GetLLMModels`, `GetTTSModels`, `GetRegisteredLLMDirs`.

### 1.7 Checks run against the shipped code

`gofmt` clean · `go build ./...` exit 0 · `go vet` exit 0 ·
`go test ./internal/image/... ./internal/backend/gguf/` ok · `node --check frontend/dist/js/chat-model.js` ok.
Release binary `build/bin/howl-chat.exe` is newer than the newest embedded frontend source and contains the
new UI markup, so the running build includes these changes. No code was edited while writing this record.

---

## 2. Not yet verified

- **Native Wails bridge** — runtime binding resolution inside WebView2.
- **Native folder picker** — the real Windows dialog, including its search-filter behaviour.
- **Actual Windows window rendering** — agent tooling cannot click native controls (WebView2 refused to open
  a CDP port: `ECONNREFUSED 127.0.0.1:9222`).
- **End-to-end model selection and loading** — select → Load → generation in the real app.
- **Remove / Rescan clicked in the real window.**

---

## 3. How to close the remaining gaps (~5 minutes, no code)

1. Launch HOWL Chat → Model Browser → confirm 25 rows and `E:\AI GGUF\GGUF` in the dropdown.
2. "+ Browse for Directory..." → pick a folder containing GGUFs → rows appear; then cancel the dialog →
   notice appears and the path box is focused.
3. Remove `E:\AI GGUF\GGUF` → rows disappear, entry leaves the dropdown, files remain on disk; re-add to restore.
4. Open the Image Generator's model menu → no chat GGUF present; and the reverse for the chat menu.
5. Select a model → Load → send one prompt.

---

## 4. Known leftovers (deliberately untouched)

1. `E:\HOWL_Chat\models\ggml-base.bin` (Whisper ASR) is still registered as `lora`; it is non-GGUF and needs
   magic-based detection. Cosmetic — it appears in neither the LLM nor the image menu.
2. `scripts/kokoro-tts.py` still shows `D:\AI GGUF\TTS\...` in its usage docstring; the actual defaults are
   project-relative/environment-driven. Documentation only.
3. `.freebuff/backup/app.go` retains the pre-change hardcoded `D:\` dataset path. It is a tool snapshot,
   ignored by the Go toolchain and not part of any build.

---

## 5. Evidence and artifacts

- Retained: `screenshots/01-added-directory.png`, `screenshots/02-empty-picker-guard.png`,
  `screenshots/03-populated-list.png`.
- The release binary contains the new UI markup (grep for `mb-add-dir-path` inside the executable:
  3 hits), confirming the running build includes the shipped frontend.
- Removed during testing: `tmpcdp/`, `tmpuiharness/`, the temporary scan probe, and the dead
  `remove_chroma.py` / `go_list.json` / `inspect_cgo_env.bat`.
- Nothing is committed; all changes are working-tree only.

## 6. Working-tree changes for this work

New: `app-model-directories.go`, `app-image-model.go`, `internal/backend/gguf/kind.go` (+test),
`internal/image/registry/kinds.go` (+test), `internal/image/registry/heal.go`.
Modified: `app-image.go`, `app-image-scan.go`, `app-model-registry.go`, `cmd/test-models/main.go`,
`internal/image/registry/registry.go`, `scripts/kokoro-tts.py`, `frontend/dist/index.html`,
`frontend/dist/js/chat-model.js`, `frontend/dist/css/model-browser.css`, `.gitignore`.

Note: the working tree was already dirty before this task (~540 entries, including unrelated in-progress work
such as the spawn-validation rework in `internal/backend/gguf/memory.go` and changes across many other
`frontend/dist` files). Those are **not** part of this change set and were not touched here.
