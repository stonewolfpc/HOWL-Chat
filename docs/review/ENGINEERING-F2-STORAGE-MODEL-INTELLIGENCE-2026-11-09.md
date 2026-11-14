# HOWL CHAT — F2 STORAGE, CONTENT LIBRARIES & MODEL INTELLIGENCE

**Project:** Stone Wolf Systems — Howl Chat  
**Phase:** F2  
**Report date:** 2026-11-09  
**Repository:** `C:\Users\StoneWolfPC\HOWL_Chat`  
**Baseline HEAD:** `9c7d2a1` (`master`)

---

## 1. Starting working-tree condition

The working tree was already dirty before F2 and was treated as project context, not
as something to clean. No `git reset`, `git clean`, `git checkout --`, revert, file
deletion, or destructive migration was performed.

At the F2 baseline recorded during this work:

- Tracked modifications: 64
- Untracked entries: 94
- Total visible tree entries recorded by the bounded inventory: 158
- F1 implementation and its report were present and intentionally preserved.
- F2 changes are uncommitted, alongside the pre-existing work.

The exact pre-existing file list was not regenerated as a new commit; the tree was
inspected and left intact. F2 only added files and made targeted edits to existing
files listed in section 16.

**Evidence:** source-confirmed for the inspected tree and commands; the tree remains
dirty by design.

---

## 2. Files inspected

The bounded inventory covered the storage, settings, model, library, and relevant UI
surfaces rather than re-auditing the whole application.

Backend and storage:

- `internal/paths/paths.go`
- `internal/paths/content.go`
- `internal/paths/content_windows.go`
- `internal/paths/content_other.go`
- `internal/library/registry.go`
- `internal/library/cards.go`
- `internal/library/registry_test.go`
- `internal/library/cards_test.go`
- `internal/modelintel/modelintel.go`
- `internal/modelintel/modelintel_test.go`
- `internal/image/registry/registry.go`
- `internal/image/registry/kinds.go`
- `internal/image/settings/*`
- `internal/settings/*`
- `app.go`
- `app-libraries.go`
- `app-model-intelligence.go`
- `app-settings-reset.go`
- `app-image-scan.go`
- `app-model-registry.go`
- `app-tts-catalog.go`
- `internal/image/scanner/*`
- `internal/resource/*`
- `internal/image/engine/*` (model-loading and settings boundaries only)

Frontend and documentation:

- `frontend/dist/index.html`
- `frontend/dist/js/settings-state.js`
- `frontend/dist/js/real-settings.js`
- `frontend/dist/js/model-directory.js`
- `frontend/dist/js/character-grid.js`
- `frontend/dist/js/character-storage.js`
- `frontend/dist/character.html`
- `frontend/dist/css/rs-content.css`
- `frontend/dist/css/rs-actions.css`
- `frontend/dist/css/rs-container.css`
- `PROJECT_ISSUES_AND_EXPANSION.md`
- `ISSUES.md`
- `docs/deferred-model-ecosystem.md`
- `docs/review/ENGINEERING-F1-LIFECYCLE-PROCESS-OWNERSHIP-2026-11-09.md`

The inventory did not attempt to re-document every page or every historical issue.

---

## 3. Existing path-resolution architecture

Before F2, `internal/paths.ProjectRoot()` was the main backend path anchor. It searches
upward from the executable and then the working directory for `wails.json` or `go.mod`,
bounded to six parent levels. That supports both `wails dev` and a built executable.

The important consequence is that the existing internal state, logs, model folders,
and much of the visible content resolution are tied to the project/installation tree
or to older settings defaults. That is not an acceptable final layout for visible user
creative content, and it is not safe to assume that a future installer will place the
executable in a writable directory.

F2 added a separate content-path resolver in `internal/paths/content.go`:

- Chats: `Documents\Howl\Chats`
- Characters: `Documents\Howl\Characters`
- Worlds: `Documents\Howl\Worlds`
- Scenarios: `Documents\Howl\Scenarios`
- Lorebooks: `Documents\Howl\Lorebooks`
- Images: `Pictures\Howl`
- Videos: `Videos\Howl`

Windows resolution uses `windows.KnownFolderPath` with `FOLDERID_Documents`,
`FOLDERID_Pictures`, and `FOLDERID_Videos` and `KF_FLAG_DEFAULT`; it does not assume
that Documents is literally below the user profile. `content_other.go` provides a
portable fallback for non-Windows builds.

`InternalConfigDir()` deliberately remains `ProjectRoot()/config`. Moving it is a
migration question and was not performed silently. New F2 state files use that
existing internal location.

**Evidence:** source-confirmed for the resolver and known-folder calls; the actual
known-folder locations on this host are not runtime-confirmed in a packaged app.

---

## 4. Existing storage and settings ownership

The bounded inventory found these ownership groups:

1. **Unified application settings** — `internal/settings.Store`, persisted through the
   app's settings directory. This owns the general application contract, including
   sampler/model/audio sections used by the existing UI.
2. **Image settings** — `internal/image/settings.Store`, a separate persistable
   contract containing the image engine configuration and last-model information.
3. **Legacy sampler and model settings files** — loaded as migration fallbacks into
   the unified store when the corresponding unified sections are absent.
4. **Frontend configuration state** — localStorage is used for some settings and
   character-cache compatibility data. Not all of it is application settings; user
   characters and chat content must not be treated as resettable settings.
5. **Image/model registry** — `internal/image/registry` persists known model entries
   and registered scan directories.
6. **Gallery/job state** — image-generated content and job state, separate from
   settings.
7. **Logs and diagnostics** — owned by the error/diagnostic paths, not creative
   content.
8. **Content folders** — previously resolved through older app paths and frontend
   storage rather than one user-owned content-library contract.

F2 did not replace these stores. It added two explicit shared layers beside them:

- `internal/library` for content-root registrations and card discovery.
- `internal/modelintel` for shared model profiles and the identification cache.

**Evidence:** source-confirmed for the inspected stores; ownership of every legacy
frontend key is not fully enumerated.

---

## 5. Existing model-discovery architecture

The existing model manager was retained. It is the source of discovered model entries
for the image/model registry and the modality-specific runtime paths.

The relevant behavior found:

- `internal/image/registry` persists entries and scan directories.
- Its `All()` path provides the catalog view; `ByType` and `ByTypeInDirs` provide
  runtime-specific selections.
- Image scanning uses `internal/image/scanner` and syncs into the registry.
- `GetRegisteredModels` previously described pruning deleted files; that behavior was
  not compatible with the requirement that an unplugged drive must not erase a
  registration.
- LLM and TTS catalogs are exposed through app bindings backed by the same registry.
- Model execution remains in the existing llama.cpp/image/TTS paths. F2 did not force
  a universal engine.

F2 changed availability handling rather than replacing discovery:

- Missing entries stay persisted.
- `Entry.Available` is computed from filesystem presence.
- `All()` retains unavailable entries for reconnect and catalog history.
- `ByType`, `ByTypeInDirs`, and `ToModelIndex` exclude unavailable entries from
  picker/runtime-facing results.
- `GetRegisteredModels` now reports available models and no longer claims to prune.

**Evidence:** source-confirmed; registry behavior is test-confirmed for a missing
entry. Actual unplug/reconnect behavior is runtime-unverified.

---

## 6. Existing model-identification capabilities

Before F2, model entries carried type, family, format, path, size, and some scanner
metadata. There was no single shared evidence/provenance/capability profile spanning
chat and image models, and no shared identification cache.

F2 added `internal/modelintel` with:

- `Input` and `Profile`.
- Evidence levels: `verified`, `source-derived`, `inferred`, `unknown`, `conflicting`.
- Capabilities: chat, image, video, TTS, ASR, embedding, unknown.
- Runtimes: llama.cpp, sd-server, sd-cli, Piper, Kokoro, unknown.
- `Classify`, `Profile.Supports`, and a persisted `Cache`.

The current classifier uses the registry's declared type, family, format, and
architecture where present. It also records a conflicting result when a filename
suggests an image model but the catalog type is not image, and when architecture and
family disagree. A cached profile is reused only when both file size and modification
time still match.

This is a foundation, not full artifact inspection. Tokenizer inspection, chat-template
extraction, config parsing, model-card retrieval, and trusted online provenance were
not implemented. Filename evidence is never promoted to a verified fact by itself.

**Evidence:** source-confirmed and test-confirmed for the implemented classification
and cache rules; tokenizer/template claims are explicitly unimplemented.

---

## 7. Confirmed defects and inconsistencies

The following defects or inconsistencies were confirmed by source inspection:

1. **Settings reset did not cover the separate image settings store.** The unified
   settings reset path could leave image settings and last-image-model state intact.
   F2 added `ResetApplicationSettings`, which resets both and clears the in-memory
   last-image state.
2. **Missing model files could be described as pruned.** The registry's public
   description and the desired unavailable-drive behavior disagreed. F2 preserves
   missing entries, marks them unavailable, and filters them from picker results.
3. **Visible content had no shared multi-root contract.** Content roots were not
   modeled as registrations with independent roles and origins. F2 added that registry.
4. **The existing settings panel had no real content-library management surface.**
   It had storage-related settings and backup/cache actions but no registered-root
   list. F2 added a real manager backed by the registry, with no fake controls.
5. **Path ownership is still split.** F2 introduced the approved visible-content
   resolver and left existing internal paths in place for compatibility; the old
   paths are not silently migrated. This is a known limitation, not a completed
   migration.

No claim is made that all pre-existing defects in the repository were found.

**Evidence:** source-confirmed for each item; items 1–3 are also covered by tests or
build checks as noted below.

---

## 8. Recommended architecture

The recommended separation is four planes:

1. **Installation plane** — executable, Wails resources, and installed runtime files.
   Read-only or conventionally installed; never used for frequently written user data.
2. **Visible creative-content plane** — user-owned Documents/Pictures/Videos folders
   and additional registered roots. Owned by the content-library registry, with
   per-root category, role, origin, format, availability, and save-destination state.
3. **Internal application-state plane** — settings, model-intelligence cache, library
   registrations, logs, indexes, and runtime state. F2 keeps the existing internal
   config location and adds new files there rather than mixing them into creative
   folders.
4. **Large assets and execution plane** — model roots, runtimes, generated media, and
   process ownership. Model registration is independent of the creative workspaces;
   runtime loading remains modality-specific.

The central rule is: **registration is not migration**. A user can point Howl at an
existing folder without Howl moving, rewriting, or importing it. A save destination is
an explicit role, not an assumption that every readable folder is writable.

---

## 9. Directory registry design

`internal/library.Registry` persists `content-libraries.json` in the internal config
directory. Each registration has:

- a stable ID derived from category plus normalized path;
- category;
- absolute, cleaned path;
- role bits (`read`, `write`, `save`);
- origin (`howl`, `external`, `imported`);
- format hint (`howl`, `sillytavern`, `unknown`);
- default-save flag;
- creation timestamp.

Design properties:

- Multiple roots per category are allowed.
- One physical directory can be registered independently for more than one category;
  registrations are keyed by category + path.
- Duplicate registration of the same category/path updates the existing registration
  rather than creating a second entry.
- A save role is required for a default save destination.
- Only one default save destination exists per category.
- Removing a registration deletes only the JSON association. It never calls
  `os.RemoveAll`, moves files, or rewrites content.
- Availability is computed at read time from `os.Stat`; a disconnected root remains
  registered.
- `WriteNewFile` rejects read-only roots, absolute paths, `..` traversal, and existing
  targets; it writes a temporary file and renames it into place.

**Evidence:** source-confirmed; multi-root, role, removal-preserves-file, duplicate
category, traversal, and read-only write behavior are test-confirmed.

---

## 10. Model catalog and intelligence design

The design keeps one discovery/registry source and adds a shared intelligence view:

```text
existing model roots / scanner
          |
          v
internal/image/registry (persisted entries + directories)
          |
          v
app.modelProfile -> internal/modelintel.Classify
          |
          v
config/model-intelligence.json (size + modtime keyed cache)
          |
          +--> GetModelCatalog / GetModelIntelligence
          |
          +--> modality-specific selectors and runtime loading (existing code)
```

Identification and runtime compatibility are separate fields/concerns:

- A profile can identify a chat model and list `llama.cpp` as a candidate runtime
  without proving that this machine has enough VRAM or a compatible binary.
- An image asset can be identified as image-capable without becoming visible in the
  chat model selector.
- The cache is metadata, not a loaded-model registry and not a claim that execution
  succeeded.

The cache uses file size and modification time, not a hash of the model. This satisfies
the bounded-discovery requirement for ordinary use. A missing file is marked
unavailable and is not offered by the picker.

No cloud lookup, telemetry, or background network activity was added.

**Evidence:** source-confirmed; classification, cache reuse, size invalidation, and
conflict detection are test-confirmed. Actual runtime load compatibility is not tested.

---

## 11. External-library compatibility design

The character-card scanner supports:

- Howl directory cards: `<root>/<card-dir>/character.json`.
- SillyTavern JSON cards: `<root>/<name>.json`.
- SillyTavern PNG cards: PNG `tEXt` chunk named `chara`, base64-decoded and parsed.

The scanner is bounded to one directory level and rejects cards larger than 16 MiB.
It reads card data; it does not rewrite the source, copy it into a Howl store, or import
it automatically.

`GetContentLibraryCharacters` aggregates readable character roots, deduplicates by
lowercased cleaned path, and retains the source registration metadata. An item reports:

- source library ID;
- category;
- format;
- origin;
- `external` status;
- `editable` status based on write role and known format;
- availability.

The settings UI can register an external folder read-only and identifies it as
external. The existing character grid does **not** yet consume this aggregation. That
is the main integration gap for F2-A; it is not described as complete.

Compatibility language is intentionally separate:

- **Discoverable:** scanner can see and parse the card.
- **Displayable:** the character browser can render it.
- **Usable:** chat can select it without copying it.
- **Importable:** a deliberate import creates a Howl-owned copy.
- **Editable:** write role plus a supported edit path.
- **Exportable:** a deliberate adapter writes a compatible representation.

Only discovery and source/format metadata are implemented. Import/export adapters and
editable external-card behavior remain future work.

**Evidence:** source-confirmed; JSON/PNG discovery and deduplication are test-confirmed
for the parser/registry level. UI display and editing are not runtime-confirmed.

---

## 12. Reset-settings completeness

`app.ResetApplicationSettings(defaults)` now:

1. Rejects nil defaults or unavailable stores.
2. Reads the image settings store before changing anything.
3. Writes the caller-supplied complete unified defaults through `SetSettings`.
4. Shuts down an active image engine if present.
5. Writes zero-value image settings.
6. Clears `lastImageModel` and `lastImageModelOpts`.
7. Rolls the unified settings back if image shutdown/write fails.

The frontend reset path calls `ResetApplicationSettings(currentSettings)` when the
binding exists, with a fallback to the older unified-only write for older builds. It
updates localStorage only after the backend call succeeds.

The reset does **not** delete chats, characters, worlds, scenarios, lorebooks,
generated images/videos, downloaded models, library registrations, or the executable.
It is not a cleanup operation and is not an uninstall operation.

Known limitation: the complete frontend/localStorage settings contract is not fully
enumerated, so the report cannot claim that every possible localStorage preference is
reset. The authoritative backend unified and image stores are covered.

**Evidence:** source-confirmed; both-store reset and creative-file preservation are
test-confirmed. Live UI reset and engine-shutdown behavior are runtime-unverified.

---

## 13. Migration strategy

No destructive migration was implemented or authorized.

The migration direction is:

1. **Keep existing data in place.** Do not move chats, cards, worlds, scenarios,
   lorebooks, images, videos, or model files as part of F2.
2. **Register, do not relocate.** Add existing folders to the content registry. A
   registration is reversible.
3. **Prefer path remapping to copying.** If a legacy directory is discovered, record
   its registration and source metadata.
4. **Detect conflicts before merging.** Duplicate paths, overlapping roots, and
   same-name cards need explicit policy rather than deletion or silent overwrite.
5. **Keep recovery copies.** Any future migration must retain originals and provide a
   rollback path.
6. **Do not migrate large model assets** merely to make the catalog uniform.

F2 intentionally did not implement automatic migration. The new default locations are
available through an explicit “Add default” action, and existing user folders can be
registered without being moved. The internal config location remains unchanged.

**Evidence:** source-confirmed by absence of migration code and presence of
registration-only operations; this is a deliberate limitation.

---

## 14. Network/offline findings

F2 added no network activity. The content registry, card scanner, model catalog, and
model-intelligence cache are local-only.

The new model-intelligence source is `local filesystem`; there is no online lookup,
upload, telemetry, or cloud inference path. The existing application's legitimate
model/runtime behavior was not removed. The frozen TTS contract was not changed.

Any future optional online metadata lookup would need explicit disclosure, user
control, a description of transmitted data, a prohibition on prompts/conversations/
characters/lore/generated content, local caching, and an offline fallback. None of that
was implemented in F2.

**Evidence:** source-confirmed for the F2 files; no runtime network capture was
performed.

---

## 15. Requirements and decision-ledger updates

The authoritative working ledger is the existing project document
`PROJECT_ISSUES_AND_EXPANSION.md`. F2 appended a dedicated section:

`## F2 Requirements Ledger — 2026-11-09`

It captures the required items:

- multiple directories per content category;
- read-only/read-write/save roles;
- external visibility without forced migration;
- SillyTavern character compatibility;
- Documents/Pictures/Videos defaults;
- separate future cleanup;
- unified model directory registry;
- unified identification across modalities;
- cached identification and revalidation;
- evidence/provenance and template detection;
- separate creative workspaces;
- shared model intelligence/resource coordination;
- deferred music/audio expansion.

The ledger also records these F2 decisions:

- internal config remains in the existing project-root `config` location;
- missing model entries remain registered but unavailable;
- removing a library registration never deletes files;
- no destructive migration, telemetry, cloud inference, or background network work.

No duplicate competing ledger was created.

---

## 16. Files modified and why

New files:

- `internal/paths/content.go` — approved visible-content category defaults and
  `InternalConfigDir`.
- `internal/paths/content_windows.go` — Windows known-folder resolution.
- `internal/paths/content_other.go` — portable fallback.
- `internal/paths/content_test.go` — default-shape and unknown-category tests.
- `internal/library/registry.go` — multi-root registrations, roles, persistence,
  availability, safe new-file writes, non-deleting removal.
- `internal/library/cards.go` — Howl/SillyTavern JSON and PNG card discovery.
- `internal/library/registry_test.go` — role/removal/duplicate/traversal tests.
- `internal/library/cards_test.go` — card discovery/format tests.
- `internal/modelintel/modelintel.go` — profiles, evidence, capabilities, runtimes,
  cache, classification.
- `internal/modelintel/modelintel_test.go` — classification, conflict, cache tests.
- `internal/image/registry/missing_test.go` — missing file retained but unavailable.
- `app-libraries.go` — content-library Wails bindings, character aggregation, native
  content-folder picker.
- `app-model-intelligence.go` — shared catalog/profile bindings.
- `app-settings-reset.go` — complete unified + image settings reset.
- `app-settings-reset_test.go` — reset and creative-file preservation tests.
- `frontend/dist/js/content-libraries.js` — real settings manager for registrations,
  roles, default save, source/format/availability, and non-deleting removal.
- `frontend/dist/css/content-libraries.css` — manager styling using the existing
  settings tokens and visual language.

Modified files:

- `app.go` — added content-library and model-intelligence fields/initialization.
- `app-image-scan.go` — corrected `GetRegisteredModels` behavior/comment for
  availability-aware catalog output.
- `internal/image/registry/registry.go` — retained missing entries, added
  `Available`, filtered unavailable entries from type/index views.
- `internal/image/registry/kinds.go` — filtered unavailable entries from type/dir
  queries.
- `frontend/dist/index.html` — added the Content Libraries card, CSS, and script.
- `frontend/dist/js/real-settings.js` — refresh the manager when settings open.
- `frontend/dist/js/settings-state.js` — reset now calls the application-wide reset
  binding with a legacy fallback.
- `PROJECT_ISSUES_AND_EXPANSION.md` — appended the F2 requirements ledger.

No pre-existing unrelated files were reverted, deleted, or cleaned.

---

## 17. Tests and exact results

Run from the repository root on 2026-11-10 (the working session following the
2026-11-09 phase assignment):

```text
set -o pipefail
go build ./... 2>&1
b=$?; echo "GO_BUILD_EXIT=$b"
go test ./internal/library ./internal/modelintel ./internal/image/registry ./internal/paths 2>&1
t=$?; echo "PACKAGE_TESTS_EXIT=$t"
go test -run 'TestResetApplicationSettings' . 2>&1
a=$?; echo "APP_RESET_TESTS_EXIT=$a"
exit $((b||t||a))
```

Exact result:

```text
GO_BUILD_EXIT=0
ok  	howl-chat/internal/library	(cached)
ok  	howl-chat/internal/modelintel	(cached)
ok  	howl-chat/internal/image/registry	2.012s
ok  	howl-chat/internal/paths	1.166s
PACKAGE_TESTS_EXIT=0
ok  	howl-chat	1.893s
APP_RESET_TESTS_EXIT=0
```

Earlier focused checks in the same F2 implementation also passed:

- `go test ./internal/library`
- `go test ./internal/modelintel`
- `go test ./internal/image/registry`
- `go test -run 'TestShutdownWithNilComponents|TestMainRegistersShutdownCallback' .`
- `go build ./...`

A first build attempt failed because `app-libraries.go` passed the string `"medium"` to
`logErrorSilent`; it was corrected to `apperrors.Medium`, after which the build passed.

Frontend syntax checks:

```text
node --check frontend/dist/js/content-libraries.js
node --check frontend/dist/js/settings-state.js
node --check frontend/dist/js/real-settings.js
```

Result: exit code 0, no syntax errors.

What the tests establish:

- multiple roots and independent categories;
- read-only write rejection and traversal rejection;
- registration removal does not remove the file;
- missing model entry remains registered but unavailable;
- default content-path shape and rejection of `music`;
- evidence/capability separation and conflict detection;
- cache reuse and size-based invalidation;
- unified reset writes defaults and clears image settings;
- reset leaves a creative file in place.

What they do **not** establish:

- live Wails UI behavior;
- actual Windows known-folder resolution in a packaged install;
- external character-grid rendering or editing;
- unavailable-drive reconnection in the packaged app;
- model runtime load success or VRAM compatibility;
- full localStorage reset coverage;
- visual quality of the new settings manager.

---

## 18. Build results

`go build ./...` completed with exit code 0 after the F2 changes, including the native
folder-picker binding and the registry availability changes.

The frontend is static Wails assets; no frontend bundler build was required for these
files. `node --check` passed on all three edited/added scripts. The HTML and CSS were
inspected for placement and selector consistency, but no browser visual verification
was performed.

---

## 19. Runtime verification results

No live Wails runtime verification was performed in this session.

Therefore this report does **not** claim runtime-confirmed behavior for:

- opening the native content-folder picker;
- adding a Windows known-folder default;
- role changes persisting across application restart;
- removing a registration and observing the folder survive;
- the character grid displaying external cards;
- settings reset with a live image engine;
- model-cache reuse in the packaged application;
- unplugging and reconnecting a model drive;
- any model loading or hardware compatibility.

The source and tests support the implemented paths, but unit tests are not runtime
evidence. The UI is source-integrated, not visually reviewed.

---

## 20. Remaining limitations

1. **Content browsers are not wired to the registry.** The settings manager can
   register roots, but existing chat/character/world/scenario/lorebook browsers still
   use their existing stores. External character cards are discoverable through the
   backend binding but are not shown in `character-grid.js`.
2. **Read-only external editing is not implemented.** External cards must not be
   edited as if they were local copies. The data model marks editability, but the grid
   does not yet enforce a separate external-item flow.
3. **No import/export adapters.** SillyTavern cards can be discovered, but deliberate
   import/export transformations are future work.
4. **No automatic migration.** Existing content and internal paths are preserved;
   users must register folders explicitly.
5. **Model intelligence is foundational, not artifact-complete.** Tokenizer,
   chat-template, config, and trusted online metadata inspection are not implemented.
6. **No unified model-management UI.** The shared catalog bindings exist, but the
   model settings UI does not yet present one cross-modality catalog with evidence and
   runtime compatibility.
7. **Audio/TTS is intentionally untouched.** The frozen contract remains frozen, and
   no new audio functionality was added.
8. **No packaged-app verification.** Known-folder redirects, drive loss/reconnect,
   and settings reset with a real engine remain unverified.
9. **The reset contract does not enumerate every localStorage preference.** The
   backend stores are covered; the complete frontend configuration inventory remains
   future work.

These limitations are material. F2 must not be reported as fully complete while they
remain.

---

## 21. Deferred requirements

Recorded in the F2 ledger:

- Comprehensive cleanup of selected Howl-generated content, separately named and
  explicitly confirmed.
- Full external-library import/export and editable external cards.
- Migration tooling and conflict resolution for legacy content locations.
- Unified model-root management UI.
- Tokenizer, chat-template, model-card, and optional trusted online provenance
  inspection.
- Shared resource coordination beyond the current metadata cache.
- Audio/music expansion when those subsystems exist.
- Any new TTS implementation or TTS contract change remains out of scope.

These are not silently dropped; they are captured with dependencies and statuses in
`PROJECT_ISSUES_AND_EXPANSION.md`.

---

## 22. F3/F4 dependencies

F3 and later work depend on the following F2 foundations:

- **F3 content workflows:** must consume `internal/library` registrations rather than
  adding another directory setting. External items need explicit source and editability
  behavior before editing is exposed.
- **F3 model selectors:** must use capability filtering from `internal/modelintel` so
  chat, image, and TTS selectors do not display incompatible models.
- **F3 runtime compatibility:** must remain separate from identification. A profile is
  not a load guarantee.
- **F4 resource/lifecycle work:** can build on the existing process ownership and
  shared metadata, but must not replace working modality engines without evidence.
- **Future audio/TTS:** must use the shared catalog and preserve the frozen contract;
  no modality should create a contradictory identification database.

The unresolved F2 UI integration is a blocking dependency for claiming end-to-end
content-library behavior, but not for the backend contracts now present.

---

## 23. F2 completion status

**F2 is not fully complete.**

### F2-A — Content and storage

| Gate | Status |
|---|---|
| Coherent path ownership | Partial: visible-content resolver added; internal paths deliberately retained |
| Default creative locations | Implemented and unit-tested; packaged known-folder behavior unverified |
| Multiple directories per category | Backend + settings UI implemented; browsers not integrated |
| Roles enforced | Registry/write enforcement implemented; existing save paths not fully routed |
| External content discoverable/usable | Character cards discoverable in backend; UI use not integrated |
| Content sources distinguishable | Metadata model and settings labels implemented; grid display not wired |
| Removal does not delete content | Implemented and test-confirmed |
| Reset covers complete settings contract | Backend unified + image stores covered; full localStorage enumeration incomplete |
| Existing content/models safe | No move/delete/migration performed; reset preservation test-confirmed |

### F2-B — Model intelligence

| Gate | Status |
|---|---|
| Multiple model roots | Existing registry retained and made availability-aware; unified manager UI absent |
| One shared identification metadata owner | `internal/modelintel` added; not yet the sole UI-facing catalog |
| Identification vs runtime compatibility distinct | Represented in the profile/design; runtime verification absent |
| Capability filtering preserves workspace separation | Capability model exists; selectors not yet switched to it |
| Cache + incremental revalidation | Implemented and tested for size/modtime; packaged behavior unverified |
| Existing modality engines remain functional | Existing execution paths were not replaced; not runtime-tested in this session |
| Supports future modality integration | Design permits it; no premature TTS/audio implementation |

### Bottom line

F2 delivered a real, tested foundation:

- user-owned visible content defaults;
- a non-destructive multi-root content registry with roles;
- Howl/SillyTavern card discovery with source metadata;
- availability-aware model persistence;
- a shared, cached model-identification layer with evidence and capability separation;
- a settings reset that covers the unified and image stores;
- a real settings UI for library registration and permissions.

It did **not** complete the end-to-end integration of those foundations into every
content browser and model selector, and it did not perform live runtime verification.
Those gaps are listed above rather than hidden behind a “complete” label.
