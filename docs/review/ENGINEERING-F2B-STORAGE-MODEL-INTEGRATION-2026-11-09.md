# Engineering F2-B — Storage and Model-Intelligence Integration

**Date:** 2026-11-09  
**Scope:** Connect the F2 content-library and model-intelligence foundations to existing Howl workflows without rewriting working subsystems.  
**Status:** **Partially complete. Not accepted as full F2 completion.** The character, save-destination, model-browser, settings-reset, and discovery integrations are implemented and tested at source/unit level. Chat, image, and video storage remain on their established owners; live Wails interaction was not performed.

Evidence labels used below: **Source-confirmed**, **Test-confirmed**, **Runtime-confirmed**, **Not verified**, **Deferred**, **Blocked**.

## 1. Starting working-tree condition and safeguards

- **Source-confirmed:** Work began on branch `master` at `HEAD 9c7d2a1` with a heavily dirty working tree containing the earlier F2 foundation changes. No reset, clean, checkout, revert, or destructive Git operation was run.
- **Source-confirmed:** Existing F2 files were edited in place; user content and model files were not moved, renamed, deleted, or reorganized.
- **Source-confirmed:** No broad drive inventory was performed. Inspection stayed inside application-owned code and test fixtures.
- **Source-confirmed:** No cloud inference, telemetry, hidden upload, remote processing, TTS contract change, new TTS feature, or music feature was introduced.
- **Source-confirmed:** Existing execution owners were retained: llama.cpp, image generation, gallery, chat persistence, and audio paths were not replaced.
- **Not verified:** No live Wails UI session was run against the modified frontend or bindings.

## 2. Existing architecture inspected

- **Source-confirmed:** Content registration and roles live in `internal/library`; Windows known-folder defaults live in `internal/paths`.
- **Source-confirmed:** Character cards are discovered from registered roots by `library.ScanCharacters`; discovery is read-only and supports Howl directory cards plus JSON/PNG SillyTavern cards.
- **Source-confirmed:** The character grid previously read only `CharacterStorage` (IndexedDB/localStorage) and had no external-card presentation path.
- **Source-confirmed:** The Model Browser previously fetched only `GetLLMModelDetails`; the shared intelligence profile was not shown there.
- **Source-confirmed:** Chat persistence writes under `%USERPROFILE%\AppData\Roaming\HOWLChat`; gallery output is owned by the image pipeline. Those paths were not silently repointed.
- **Source-confirmed:** Settings reset already split application and image stores; the UI confirmation text was weaker than the actual reset scope.

## 3. Exact files and packages changed

**Backend**

- `internal/library/registry.go` — safe `WriteFile` replacement using staged `.bak` handling; save-only destinations can write.
- `internal/library/cards.go` — external card items now carry description, personality, scenario, first message, and example dialogue.
- `internal/library/registry_test.go` — save-only creation/replacement coverage.
- `internal/modelintel/modelintel.go` — bounded sidecar inspection, tokenizer/chat-template metadata, metadata fingerprint, and conflict/evidence handling.
- `internal/modelintel/modelintel_test.go` — sidecar parsing, source-derived evidence, fingerprint invalidation, and cache mismatch coverage.
- `app-model-intelligence.go` — passes the sidecar fingerprint into the profile cache.
- `app-libraries.go` — `SaveContentItem` routing through the selected save destination; safe content path segment handling.
- `app-libraries_test.go` — save-destination and no-destination behavior.

**Frontend**

- `frontend/dist/js/character-grid.js` — loads external cards alongside local cards, renders source/read-only badges, and opens a read-only view.
- `frontend/dist/js/character-form.js` — external-card read-only state, disabled fields, hidden save action, and write guard.
- `frontend/dist/js/content-save.js` — bridge from visible page saves to the selected library destination; reports write failures.
- `frontend/dist/js/chat-model.js` — Model Browser populates evidence, capabilities, candidate runtimes, tokenizer, source, and notes.
- `frontend/dist/js/world-form.js`, `scenario-form.js`, `lorebooks.js` — save a visible copy through the bridge while retaining their existing stores.
- `frontend/dist/js/settings-state.js` — confirmation now states that image settings are included and that content/libraries/models survive.
- `frontend/dist/css/content-integration.css` — external-card and model-intelligence styling using existing tokens.
- Page shell wiring: `character.html`, `world.html`, `scenario.html`, `lorebooks.html`, `story-tools.html`, `index.html`.

## 4. Content-library integration results by category

| Category | Result | Evidence |
|---|---|---|
| Characters | External Howl/SillyTavern cards appear in the existing grid with source and read-only labeling; external cards cannot enter local edit/delete flows. | **Source-confirmed**, **Test-confirmed** (discovery/card tests), **Not verified** in live UI |
| Characters — save | Existing IndexedDB/localStorage remains the working store; a visible JSON copy is written when a save destination is selected. | **Source-confirmed**, **Test-confirmed** (save binding) |
| Worlds | Existing storage remains authoritative; visible JSON copy routed when a save destination exists. | **Source-confirmed**, **Not verified** live |
| Scenarios | Existing storage remains authoritative; visible JSON copy routed when a save destination exists. | **Source-confirmed**, **Not verified** live |
| Lorebooks | Existing storage remains authoritative; visible JSON copy routed when a save destination exists. | **Source-confirmed**, **Not verified** live |
| Chats | Chat persistence remains the AppData owner. No mirror was added because it would create a second, potentially stale chat copy without changing load authority. | **Blocked** for source-of-truth integration; **Deferred** |
| Images | Gallery remains the image-output owner; no path rewrite was attempted. | **Blocked** for source-of-truth integration; **Deferred** |
| Videos | Video output remains with the generation pipeline; no path rewrite was attempted. | **Blocked** for source-of-truth integration; **Deferred** |

**Source-confirmed:** Read-only roots are skipped for discovery where not readable and `WriteFile` checks the role. Removing a registration remains association-only. Missing roots do not delete registrations. External cards are displayed from the original path and are never copied by discovery.

**Not verified:** No controlled drive-disconnect/reconnect runtime test was performed; the existing registry persistence behavior is covered only by source-level and focused test evidence.

## 5. Model catalog and identification results

- **Source-confirmed / Test-confirmed:** `InspectArtifacts` reads only `config.json`, `tokenizer_config.json`, `generation_config.json`, and `model_index.json`, capped at 1 MiB each. It does not hash model weights.
- **Source-confirmed / Test-confirmed:** Profiles now expose architecture, family, tokenizer, chat-template presence, fingerprint, capabilities, candidate runtimes, evidence, source, and notes.
- **Source-confirmed:** A changed sidecar fingerprint invalidates the cached profile even when file size and modification time are unchanged. A missing sidecar produces no fingerprint and falls back to size/modtime, which is a documented residual risk.
- **Source-confirmed:** The existing Model Browser displays the shared profile without merging chat and image workspaces. Chat model listing still comes from the existing chat model registry.
- **Source-confirmed:** Evidence labels are exposed as Verified, Source-derived, Inferred, Unknown, or Conflicting; candidate runtimes remain separate from capabilities.
- **Not verified:** No live model-directory selection was exercised in the Wails UI. No model-loading regression test was run against a real model.
- **Deferred:** Trusted online metadata lookup; request-level capability negotiation; tokenizer execution; template application; audio modality expansion.

## 6. Settings reset results

- **Source-confirmed:** Reset remains split between application settings and image settings. The image engine is shut down/cleared as part of the existing reset path, and content/library/model files are preserved.
- **Source-confirmed:** UI confirmation now explicitly names image-generation settings and lists preserved content, registered libraries, and model files.
- **Test-confirmed:** Existing `TestResetApplicationSettings` passed after the integration edits.
- **Not verified:** Live reset with the image engine running, including rollback behavior if one store fails.
- **Source-confirmed:** Comprehensive cleanup and uninstall remain distinct future operations and were not implemented as reset.

## 7. Tests, builds, and exact outcomes

Commands were run from the project root on the supplied Windows host.

| Check | Result |
|---|---|
| `gofmt -w` on changed Go files | **Passed** (`GOFMT_EXIT=0`) |
| `go build ./...` | **Passed** (`GO_BUILD_EXIT=0`) |
| `go test ./internal/library ./internal/modelintel ./internal/image/registry ./internal/paths` | **Passed** (all `ok`) |
| `go test -run 'TestResetApplicationSettings|TestSaveContent' .` | **Passed** (`ok howl-chat`) |
| `go test -run 'TestSaveContentItem' .` | **Passed** (`ok howl-chat`) |
| `go test -race ./internal/library ./internal/modelintel` | **Passed** (`RACE_EXIT=0`) |
| `node --check` on character grid/form, chat model, content save, world/scenario/lorebooks | **Passed** |
| `git diff --check` on changed tracked files | **Passed** (`DIFF_CHECK_EXIT=0`; Git emitted only LF/CRLF warnings) |

`git diff --check` does not inspect untracked new files. `frontend/dist/js/content-save.js` and `frontend/dist/css/content-integration.css` are new; they were syntax-checked only where applicable (JavaScript), not checked as Git-tracked diffs.

## 8. Runtime checks actually performed

- **None of the modified Wails screens were opened in a live session.**
- No character grid was visually inspected with a real external card.
- No save-destination write was performed from the running UI.
- No model was selected in the Model Browser, and no model was loaded.
- No settings reset was run against live stores.
- No reconnect test was performed with a removable or redirected library root.

These are **Not verified**, not passes.

## 9. Runtime checks not performed and why

- **Blocked:** Live UI verification requires an interactive Wails session and user-visible test data. The assignment required documenting the procedure instead of claiming it passed.
- **Blocked:** Chat/image/video source-of-truth integration requires changing their established persistence and output ownership. That would be a migration/ownership change, not a safe integration slice, and could break existing sessions or generation output.
- **Not verified:** Missing-drive reconnection and overlapping-root duplicate behavior were not exercised against real removable media; the E: drive was intentionally not touched.
- **Not verified:** Settings reset rollback across both stores with a running image engine.
- **Not verified:** Model cache behavior across an application restart and packaged installation.

### Required manual runtime procedure

1. Register a character folder containing one Howl card and one SillyTavern JSON/PNG card. Open Characters and confirm both appear, with the external badge and original source path.
2. Click the external card; confirm fields are visible, all inputs are disabled, the save action is absent, and closing/reopening a local card restores the normal editable state.
3. Select a save-destination library for worlds/scenarios/lorebooks/characters, save one item of each type, and confirm the JSON file appears under that root while the existing page still reloads its item.
4. Disconnect/rename the external root while the app is open, refresh, then restore it; confirm discovery recovers and no registration is deleted.
5. Open Model Browser, select a model with sidecar metadata, and confirm evidence, capabilities, runtimes, tokenizer, source, and notes are populated and that no model-loading action is implied.
6. Run Reset Application Settings with image settings active; confirm the confirmation text, image engine shutdown, preserved content/libraries/models, and a second reset succeeds.
7. Verify no read-only card can reach the save path and no failed library write is reported as success.

## 10. Remaining defects and risks

- **Blocked — ownership:** chats remain AppData-owned; images/videos remain pipeline/gallery-owned. The selected library is not yet their source of truth.
- **Risk — duplicate visible copies:** local page stores remain authoritative while a visible library copy is written. A later edit can leave the visible copy stale until the next save; no reconciliation job exists.
- **Risk — external card scope:** discovery is one level deep, matching current formats. Deeply nested collections are not discovered.
- **Risk — model fingerprint:** models without sidecar metadata still rely on size/modtime. A same-size, same-mtime replacement can reuse a stale profile.
- **Risk — sidecar conflicts:** sidecar evidence is presented conservatively, but no user-facing conflict-resolution workflow exists.
- **Not verified — UI:** external card layout, disabled-field behavior, model metadata wrapping, and error notifications have not been visually or interactively checked.
- **Not verified — reset:** rollback and image-engine shutdown were not exercised live.

## 11. Requirements ledger changes

- **Source-confirmed:** The existing F2 ledger remains the single authoritative requirements document in `PROJECT_ISSUES_AND_EXPANSION.md`.
- **Source-confirmed:** A new `F2-B Deferred Future Requirements — 2026-11-09` section records chat cues/image cues, llama.cpp features, capability-gated model-era options, comprehensive sampling including XTC and the three unidentified samplers, SillyTavern import/field mapping/Description/LoRA/lorebook work, Author's Note placement, single-character context capacity, multi-character isolation/flushing, and hardware-aware image RAM/VRAM/offload/OOM controls.
- **Source-confirmed:** Each deferred item names an eventual owner/roadmap slot and dependency rather than creating a competing ledger.
- **Deferred:** All items in that new section remain unimplemented by F2-B.

## 12. Dependencies and recommended next engineering assignment

**F2-B is not complete as a full storage migration.** It is complete as a bounded integration slice: the foundations are now connected to character discovery, visible content saves, the Model Browser, and settings-reset communication, with focused tests and a clean build.

Recommended next assignment: **F3-A — Content ownership cutover design and migration safety**, covering:

1. Decide and document the source-of-truth owner for chats, images, and videos.
2. Add a reversible, previewable import/cutover path for chat sessions and generated media without silent moves.
3. Add a reconciliation strategy for the visible library copies written by F2-B.
4. Add live Wails acceptance tests for external-card read-only behavior, missing-root reconnection, save destinations, and reset rollback.
5. Only then expose broader model capability controls, starting with evidence-gated llama.cpp flags and sampling controls.

No F3 work should assume that the F2-B mirror writes are a completed ownership migration.
