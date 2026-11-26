# ADR-002 — Storage Locations, Picker Navigation, and Content Routing

**Status:** Accepted as policy (owner-authorized documentation, 2026-11-22).
**Implementation:** **Not authorized.**
**Refines:** `ADR-001-FOUNDATIONAL-BACKEND-REBUILD-2026-11-22.md` decision D1 (chat location).
**Owns:** per-category storage locations, picker navigation state, content routing, and the browser-surface rule (N1–N19).
**Does not own:** path-resolution mechanics, configuration authority, persistence durability, runtime resolution, or ownership classes — those belong to ADR-001. This record refines exactly one ADR-001 decision (D1, chats) and settles nothing else there.
**Discoverability:** linked from `HOWL-MEMORY-09-PROGRAM-DIRECTION-2026-11-21.md` §2.3 and from the HCDM architectural contract register (`HCDM-2026-11-21-01` §F). **A link establishes discoverability, not authority.**
**Frozen boundaries untouched:** FC-5 revision 3, FC5-31 (`Unresolved`), R6 draft (unapplied).

## Evidence labels

`SC` source-established · `AO` archived observation · `RC` runtime-confirmed · `UN` unestablished.

---

## 1. Context — the actual file-selection surface

A policy covering only the native-dialog API would leave most of the surface ungoverned.

| Surface | Count | Evidence |
|---|---|---|
| Native dialogs (Wails `runtime.*Dialog`) | 6 | `app-model.go:851`, `app-model-registry.go:232`, `app-libraries.go:203`, `app-image-diagnostics.go:198,211`, `app-image-generate.go:171` `SC` |
| Browser `<input type="file">` | 6 | `character.html:151`, `scenario.html:146`, `world.html:146`, `story-tools.html:263`, `image-generator.html:1059,1094` `SC` |
| Drag-and-drop | ≥1 | `image-generator-ui.js:570–578` `SC` |
| Any future File System Access API | 0 | none present; governed prospectively (§7) |

Related findings:

- **No picker navigation state exists in Go.** Every `lastDir`/`lastPath`/`recentDir` search returns zero `SC`.
- The model picker's start directory is a **storage** value: `DefaultDirectory: a.GetModelsDir()` (`app-model.go:853`), and `GetModelsDir()` runs `MkdirAll` before returning (`:845`) `SC` — so opening the picker creates the app's models directory.
- The generic directory picker passes **no** `DefaultDirectory`, surrendering navigation to the shell's own MRU, shared across unrelated dialogs (`app-model-registry.go:232–234`) `SC`.
- `howl-last-model-path` (`chat-model.js:393`) is the **selected model**, not a dialog directory, and is a single global key `SC`.
- `settings-controls.js:125` records that a broken `webkitdirectory` picker was removed rather than left pretending `SC`.

**This policy is new behaviour, built once, correctly.** There is no legacy navigation state to migrate.

---

## 2. Decision — three state families, no cross-writes

| Family | Question | Class | Owner | Store |
|---|---|---|---|---|
| **Category storage location** | *Where does this content live?* | app-state (R2) | PathService + ConfigAuthority | `storage-locations.json` + adoption records |
| **Picker navigation** | *Where did I last browse for this purpose?* | app-state (R2) | ConfigAuthority (NavigationState) | `navigation.json` |
| **External resource selection** | *Which model / library / file is chosen?* | external reference (R2 metadata, R4 data) | ModelConfigService, library registry | existing registries |

No new service: NavigationState and the authority record are stores owned by ConfigAuthority.

**Product guarantee:** loading a model, browsing for a preset, importing a character and choosing an output folder never hijack one another's locations.

---

## 3. Category storage policy

| Category | Default | Class | Notes |
|---|---|---|---|
| chats | `Documents\Howl\Chats` (new) / adopted root (existing) | user-content | §4 |
| characters | adopted legacy root | user-content | ADR-001 §11.2 |
| worlds, scenarios, lorebooks, text projects | `Documents\Howl\<Category>` | user-content | known-folder resolver |
| images | `Pictures\Howl` | user-content | |
| videos | `Videos\Howl` | user-content | |
| settings, registries, logs, caches, journals, navigation | `%APPDATA%\Howl` | app-state | not user-selectable |
| models, image models, libraries | user-chosen | external | read-only; grant required to write |

Defaults resolve through the **Windows known-folder API**, never by string concatenation, so redirected and OneDrive-backed folders work without special cases (`internal/paths/content.go:21–24,60–64` `SC`).

### 3.1 Presets are three kinds

| Kind | Lives in | Browsable | Rule |
|---|---|---|---|
| **Managed presets + associations** | R2 application data | no — selected from app UI | stable IDs; never treated as files |
| **Imported preset files** | wherever the user chose | via `preset.import` | source only; reference recorded, no write |
| **Exported preset documents** | `Documents\Howl\Presets` (default) | via `preset.export` | output; user-content |

Only import/export have picker memory. Managed presets have no picker. Presets are **not** a model-like external category.

---

## 4. Chats — two populations (settles ADR-001 D1)

D1 was a compatibility decision, not a product position. Both hold when the rule is conditioned on whether content already exists:

| Condition | Chat root |
|---|---|
| `NeverInitialized` **and** no user data of any kind (no adopted roots, no content, no chats) | `Documents\Howl\Chats` |
| Initialized installation (witness + record) with chats at `%APPDATA%\HOWLChat` | adopted legacy root, in place |
| **Initialized installation with no chats yet but other user data** | **adopt the legacy root — not Documents.** The installation is not fresh; only its chat store happens to be empty |
| `RecoveryRequired` | no chat root chosen; safe mode (ADR-001 §6.2) |

Consequences:

1. **Two chat roots will exist in the wild indefinitely.** PathService must support **multiple adopted roots per category** — a generalisation of the adoption record already required for characters.
2. **Chats are user-content in both populations.** `%APPDATA%\HOWLChat` being app-relative does not make it disposable: reset, wipe and migration protections apply identically.
3. Relocation of existing chats is **explicit only** — plan, journal, confirmation. Never automatic.

---

## 5. Native dialog navigation

### 5.1 Every dialog receives an explicit, validated start directory (N13)

| Order | Candidate | Validation |
|---|---|---|
| 1 | `navigation.json[purpose].lastDir` | exists **and** identity matches §5.3 |
| 2 | purpose's category root | exists; **existence check only** — never a call that creates it (`GetModelsDir()` disqualified, `app-model.go:845`) |
| 3 | Windows **Documents** known folder | exists |
| 4 | Windows **user profile** known folder | exists |
| 5 | purpose-specific validated fallback supplied by PathService | exists |

- **The installation root is not a universal fallback.** It may be read-only, on removable media, or hold resources inappropriate to browse. It is used only where a purpose explicitly permits it, and no current purpose does.
- `%USERPROFILE%\Documents` resolves through the **known-folder API**, not a string join.
- If every candidate fails (a broken profile), PathService reports a diagnostic and the dialog opens at the shell default — a **documented degradation with a report**, never silent.

### 5.2 Purposes

`model.file` · `model.dir` · `image.model.file` · `image.output` · `video.output` · `preset.import` · `preset.export` · `character.import` · `world.import` · `library.dir`.

Each purpose has its own key. No shared key may exist.

### 5.3 Identity, not path strings (N15)

Record: `{lastDir, resolvedPath, volumeSerial, fileID, updatedAt}`.

| Condition | Behaviour |
|---|---|
| identity matches | open there |
| directory missing / volume offline | fall back per §5.1, notify, **record unchanged** |
| **identity differs** (directory replaced; reparse point now resolves elsewhere) | **do not open the new destination.** Open the fallback, show a notice naming the change, update the record only when the user picks again |

### 5.4 Concurrency (N14)

`navigation.json` is strict read-modify-write **under the DurableStore lock**, never from a cached copy:

```
lock → read current document → set one purpose key → write → unlock
```

In-memory caching of the whole navigation document is prohibited — that is how one picker erases another's update. Per-purpose keys mean a concurrent `preset.import` update cannot affect `model.file`; the lock means neither update is lost.

### 5.5 Save failure must be reported honestly (N19)

If `navigation.json` cannot be written (lock timeout, sharing violation, disk error):

- the dialog **completes** — navigation persistence never blocks the user's work;
- the user is told the location **was not remembered**, with the reason;
- Howl never states or implies it will be remembered next time;
- the failure is reported once per session per purpose, not on every dialog.

---

## 6. Browser surfaces

### 6.1 What the evidence establishes

| Surface | Established | Not established |
|---|---|---|
| Character images (`character-images.js:15–33`) | file content converted to data URLs, pushed into `characterState.characterImages` | full persistence path of the owning character record — **`UN`** |
| World images (`world-form.js:111–129`) | same pattern into `worldState.worldImages` | full persistence path of the owning world record — **`UN`** |
| Image-generator upload (`image-generator-ui.js:1243–1256`) | content becomes `initImageData` in memory | that it is persisted at all — **`UN`** |
| Knowledge upload (`rag-knowledge.js:33–39`) | text sent via `UploadRAGDocumentText`; source comment states the input gives content, not a path | — (backend import path evidenced) |

So the central claim is supported: **browser mechanisms receive content, not a usable destination path, and therefore cannot establish a storage destination.** Only the RAG handler is evidenced as completing a backend persistence step.

### 6.2 N17 (normative)

1. Browser file inputs select **source files only**.
2. Drag-and-drop selects **source files only**.
3. Neither establishes, changes or implies a storage destination.
4. Imported content passes through the appropriate import and validation workflow.
5. A browser-selected item never becomes a persistent content root, and never supplies a path later treated as one.
6. Any future **File System Access API** integration obeys the same purpose-specific navigation and authorization rules as native dialogs.

### 6.3 WebView limitation, documented rather than papered over

The WebView owns the `<input type="file">` dialog, so Howl **cannot** reliably supply its starting directory. Browser pickers therefore **must not claim navigation parity**: no `navigation.json` purpose is written for them, and no UI text implies a remembered location. This is documented behaviour, not a defect.

---

## 7. Normative rules

- **N1** One default resolver per content category.
- **N2** One configurable storage location per user-content category.
- **N3** One last-used directory per picker purpose.
- **N4** Picker memory and storage settings live in different stores with different writers; a picker write never modifies a storage location.
- **N5** Choosing a **source** never changes an **output** location.
- **N6** Changing a storage location does **not** move existing content; moving requires a separate explicit action with its own plan and journal.
- **N7** An invalid, disconnected or unavailable **storage** location refuses the save with a clear reason. No silent fallback: a save never lands in a directory the user did not choose.
- **N8** An invalid **navigation** location produces a notice and a default start point — a different rule, because nothing is written.
- **N9** Settings, registries, journals, logs and caches are never mixed into Documents or Pictures.
- **N10** A **model directory is an external resource location, not Howl-owned storage.** Selecting it grants read access only; no manifest, index, sidecar or reorganisation without an explicit per-directory grant (ADR-001 §10).
- **N11** Opening any picker performs no filesystem creation or modification.
- **N12** Picker default directories are navigation; a storage location may seed them only as a read-only hint, never by calling a function that creates it.
- **N13** Every native dialog receives an explicit, validated start directory (§5.1); the shell MRU is never relied upon.
- **N14** Navigation updates are serialized read-modify-write; no cached whole-document writes.
- **N15** A changed directory identity produces a notice and a fallback, never a silent redirect.
- **N16** Presets are managed records, imported sources, or exported documents — never one undifferentiated category.
- **N17** Browser file inputs and drag-drop are sources only (§6.2).
- **N18** Chat-root selection is an **adoption decision** based on installation identity and prior initialization, evaluated once and recorded; a first save never re-evaluates it.
- **N19** Navigation persistence failure never blocks the dialog and is always reported honestly (§5.5).

---

## 8. Interaction matrix — who may write what

| User action | storage-locations | navigation | registries | content |
|---|---|---|---|---|
| Browse model file | — | `model.file` | — | — |
| Load model | — | — | model admission | — |
| Add scan directory | — | `model.dir` | scan dirs | — |
| Change image output | image output | `image.output` | — | — |
| Import world | — (unless promoted to a location) | `world.import` | library registry | new files at destination, via grant |
| Change characters location | characters (separate decision) | — | — | **no move** (N6) |
| Save image | — | — | image registry | output file at configured root |
| Browser image / knowledge input | — | **no purpose** (§6.3) | — | owning record or RAG import |

---

## 9. Acceptance criteria

1. Browse `D:\AI\Models` for a model, then `G:\Howl Presets` for a preset; reopening the model picker lands on `D:\AI\Models`.
2. Model picker with no models directory present: opens at a fallback **and creates nothing** (N11, §5.1).
3. **Model-directory read-only check without recursive hashing:** snapshot depth-1 entries only (name, size, mtime) of the chosen directory before and after picker + load; assert no change. No recursive walk, no content hashing, no traversal of model trees.
4. Image output set to a disconnected drive: save refused with a named reason; no file in Documents or any fallback (N7).
5. Navigation record for a deleted directory: opens at fallback with a notice; `navigation.json` byte-identical.
6. Two simultaneous picker updates to different purposes both survive (N14).
7. Remembered directory replaced by a junction to elsewhere: notice shown; new destination not silently used (N15).
8. Fresh profile: first chat save creates `Documents\Howl\Chats`. Existing profile: chats stay at `%APPDATA%\HOWLChat`, no new directory created (§4).
9. Corrupt `navigation.json`: pickers still open via fallback; `storage-locations.json` unaffected.
10. Changing a category location leaves existing files in place; listing identical (N6).
11. **Browser surface containment.** For each of the six `<input type="file">` elements and the image-generator drag-drop zone, perform a selection, then assert:
    - **(a)** no `navigation.json` purpose changed;
    - **(b)** no entry in `storage-locations.json` changed;
    - **(c)** **the browser picker or selection mechanism itself creates no directories.** Downstream creation by an approved import workflow is governed by that workflow's own rules and is not a failure of this criterion;
    - **(d)** the selected content is consumed **only by its owning workflow**: persisted through the appropriate record or RAG import path where applicable, **or retained as transient input where the workflow requires it** (image-to-image initialisation). **No browser-selected path becomes a destination.**

    No recursive hashing: the check is (a)–(c) by document comparison plus a depth-1 listing of the category roots.

---

## 10. Migration

- Chats: new installs → Documents; existing installs → adopted root. Relocation only on explicit request (§4).
- Adopted-root mechanism generalised from characters to chats, and available to any category.
- No legacy navigation state exists (`SC`: zero matches).
- `howl-last-model-path` remains **selected-model memory**, never migrated into navigation.

---

## 11. Non-goals

Relocating existing chats (§4) · relocating or merging characters (ADR-001 §11.2) · portable mode (ADR-001 D2) · moving content on a location change (N6) · writing into external model directories (N10) · controlling the shell's own dialog MRU, which is bypassed rather than governed.

---

## 12. Open register

| Item | Status |
|---|---|
| Persistence path for character image data URLs | `UN` |
| Persistence path for world image data URLs | `UN` |
| Whether image-generator uploads are ever persisted | `UN` |
| Portable mode | deferred (ADR-001 D2) |
| Chat relocation for existing installs | explicit only; not scheduled |
| Character permanent root | needs installed-layout evidence |
| Largest chat/content document size (for D3) | `UN` |

All three persistence questions are open and **non-blocking** for this policy.

---

## 13. Authorization

This record fixes policy only. It authorizes no implementation, no code, schema, API, migration, test, build, launch or process change, and no amendment to any frozen contract. FC-5 revision 3 remains frozen, FC5-31 remains `Unresolved`, R6 remains unapplied.
