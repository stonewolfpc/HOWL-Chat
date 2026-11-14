# Architecture F2-C — Content Ownership, Persistence and Cutover Design

**Date:** 2026-11-09  
**Assignment type:** Read-only architecture and evidence. No migration, content movement, code change, or user-data modification was performed.  
**Status:** Design recorded; **F2-C does not implement the cutover**.

Evidence labels: **Source-confirmed**, **Test-confirmed**, **Runtime-confirmed**, **Inferred**, **Unknown**, **Deferred**.

---

## 1. Existing ownership map

**Source-confirmed:** Howl currently has four different persistence owners, not one content system:

1. Browser localStorage/IndexedDB for characters, worlds, scenarios, and lorebooks.
2. `%USERPROFILE%\AppData\Roaming\HOWLChat` for chats.
3. A gallery directory for generated images, defaulting to `paths.ProjectRoot()/images/gallery`.
4. The content-library registry, which writes visible JSON copies but is not the reader for any of those JSON files.

| Category | Current authoritative store | Other copies or readers | Current delete owner | Source identity today |
|---|---|---|---|---|
| Characters | `howl-chat-characters` localStorage metadata plus `CharacterStorage` IndexedDB images | F2-B `character.json` library copy; native `Characters/` legacy directory for a prompt fallback; external card discovery | `CharacterStorage.remove` removes local metadata and images only | Local record, or external card path; the F2-B library copy is not read back |
| Worlds | `howl-worlds` localStorage plus `HOWLStorage` `world-images` IndexedDB | F2-B `world.json` copy; chat session `world.json` snapshot | `HOWLStorage.deleteWorld` removes local record and images | None on the world itself beyond its local ID; F2-B copy has no read path |
| Scenarios | `howl-chat-scenarios` localStorage plus `HOWLStorage` `scenario-images` IndexedDB | F2-B `scenario.json` copy; chat session `scenario.json` snapshot | `HOWLStorage.deleteScenario` removes local record and images | Same limitation as worlds |
| Lorebooks | `howl-chat-lorebooks` localStorage | F2-B `lorebook.json` copy; chat `lorebooks.json` snapshot | Local array filter and `saveLorebooksToStorage` | No file source identity |
| Chats | `%APPDATA%\HOWLChat\chats\<id>` (`meta.json`, `history.ndjson`, `history.json`, context files) | `HOWLChatStorage` `chat-data` cache under key `howl-chat-current-chat`; localStorage fallback | `chat.PersistenceManager.DeleteChat` removes the AppData chat folder | Chat ID and AppData folder |
| Images | Gallery directory (`EngineConfig.OutputDir`; default project-relative `images/gallery`) | PNG plus JSON sidecar; `ListGalleryImages`; `SaveImageAs` exports a user-selected copy | `DeleteGalleryImage`/gallery `Delete` | Recorded absolute path plus sidecar filename; no library registration |
| Videos | The same image output directory, written directly by `sdcpp.generateVideo` | Base64 `VideoData` returned to the UI; no sidecar and no gallery listing | No gallery-owned delete path for the returned clip | Output filename only; no stable catalogue entry |

Additional source-confirmed ownership details:

- `CharacterStorage` uses a separate IndexedDB database, `howl-character-images`, while worlds and scenarios use the `HOWLChatStorage` database with `world-images` and `scenario-images` stores.
- `chat.PersistenceManager` calls `getAppDataDir`, which hardcodes `%USERPROFILE%/AppData/Roaming/HOWLChat`. `NewPersistenceManagerWithDir` exists for tests and shows that the manager can accept an arbitrary root.
- The native `Characters`, `Worlds`, and `Scenarios` directories are resolved from the current working directory or executable directory. `CreateCharacter` writes `Characters/<name>/character.json`; `loadCharacterFromDiskByID` scans it as a prompt fallback. The character grid does not use that native directory as its primary catalog.
- `paths.ProjectRoot()` searches upward from the executable/project markers. The default gallery path therefore depends on where the app is launched/installed, while chat data depends on the user profile.
- The library registry stores roots, roles, origin, format, and one default save root per category. Removing a registration deletes only the registration. It does not own the content files.

**Requirement consequence:** the current system cannot answer “who owns this item?” consistently. A character can exist in localStorage, in a native `Characters` directory, in a library JSON file, and inside a chat snapshot at the same time.

---

## 2. Dual-storage analysis

### 2.1 The F2-B JSON copy is a write-only mirror

**Source-confirmed:** `frontend/dist/js/content-save.js` calls `SaveContentItem` after the existing page store succeeds. `SaveContentItem` writes a JSON file through `library.Registry.WriteFile` to the category’s default save root.

**Source-confirmed:** None of the four page loaders reads the library JSON back:

- Characters load from `CharacterStorage.loadAll()` plus external character discovery.
- Worlds load from `howl-worlds`.
- Scenarios load from `howl-chat-scenarios`.
- Lorebooks load from `howl-chat-lorebooks`.

The library file is therefore **not authoritative content, not an import source, and not a persistent mirror used by the UI**. It is, at best, an **unlabelled export snapshot**. Because it is written automatically, it is also a compatibility liability.

### 2.2 Lifecycle by category

| Operation | Characters | Worlds | Scenarios | Lorebooks |
|---|---|---|---|---|
| Normal load | Local metadata + IndexedDB images | `howl-worlds` + IndexedDB images | `howl-chat-scenarios` + IndexedDB images | localStorage |
| Edit | Local record; library file is overwritten only if the save succeeds locally | Local record; library file overwritten by title | Local record; library file overwritten by title | Local record; library file overwritten by title |
| Create | Local save, then library write | Local save, then library write | Local save, then library write | Local save, then fire-and-forget library write |
| Delete | Local record/images only | Local record/images only | Local record/images only | Local record only |
| Read library copy | Never for local cards; external discovery is a separate path | Never | Never | Never |
| Copy/duplicate | `CharacterStorage` upsert by ID; no library copy operation | `copyWorld` creates local copy, then later saves can overwrite the same title path | No dedicated copy path found | No dedicated copy path found |

**Source-confirmed:** The write order is existing store first, library file second. There is no transaction spanning the two.

### 2.3 Failure and disagreement behavior

**Source-confirmed:** A library-write failure cannot roll back the browser save:

- Character save: the local record is saved, then `ContentSave.persist` is attempted. The helper reports a notification, but the character save continues to completion.
- World save: `storeWorld` succeeds, then the library write is attempted. A library failure is caught and ignored by the caller; the world still closes and the success alert is shown.
- Scenario save: local storage succeeds, then the library write is attempted. The library failure is caught by the caller and the form still closes.
- Lorebook save: local storage succeeds, then the library write is started without awaiting it and its rejection is ignored by the caller.

**Source-confirmed:** If the local save fails, the library write is not attempted. If the local save succeeds and the library write fails, the UI has a local success and a partial library failure. The UI does not present one transactional result, and there is no retry queue.

**Source-confirmed:** If the local record is deleted, the library JSON is not deleted. It remains as an orphan. If the local record is renamed, the old library JSON remains and a new one may be created. The copies can disagree indefinitely.

**Source-confirmed:** Because no loader reads the library file, an older library copy cannot currently overwrite a newer local record. That safety is accidental, not an ownership rule. If a future loader starts reading the library file without a revision/conflict contract, stale-content overwrite becomes possible.

### 2.4 Chat dual storage

**Source-confirmed:** The chat backend is the primary chat store. The frontend also keeps a minimized current-chat cache in `HOWLChatStorage` (`chat-data` store; localStorage fallback).

**Source-confirmed:** `saveChatData` writes the IndexedDB cache first and then invokes `SaveChatData` on the backend without awaiting it; when that fails, it logs a fallback attempt and does not surface a single authoritative failure to the user. `dataSync` awaits its own save path but its autosave catches and logs errors.

**Source-confirmed:** `LoadChat` reads the AppData chat; `LoadChatData` returns in-memory backend state or falls back to the persisted chat. The browser cache is then hydrated and validated against `GetChatList`.

**Source-confirmed:** A saved chat embeds `world.json`, `scenario.json`, `characters.json`, `lorebooks.json`, and `memories.json`. Those are chat-scoped snapshots, not references to the live creative catalog. Editing a world or character later does not update an existing chat’s embedded copy. The chat can therefore display an older world/character than the one in the creative page.

**Source-confirmed:** `history.ndjson` is documented as the true history source; `history.json` is a fallback snapshot. `SaveChat` writes both. The frontend can build history from DOM text, and the backend has guards, but multiple representations still exist.

### 2.5 Generated media

**Source-confirmed:** `gallery.Service.Save` writes an image file, then a JSON sidecar, then creates a thumbnail directory. It returns an error if the metadata write fails, but the image has already been written. This is not an atomic pair.

**Source-confirmed:** Both `GenerateImage` and the job completion callback only print a warning when gallery saving fails. A generation can be shown as successful while its gallery persistence failed.

**Source-confirmed:** `gallery.Delete` ignores removal errors and always returns nil. The UI can report a successful delete when files remain.

**Source-confirmed:** Gallery listing is driven by JSON sidecars; an image without a sidecar is invisible to the gallery. `resolveImagePath` repairs stale absolute paths by preferring the sidecar’s basename, but the sidecar still records the original absolute path.

**Source-confirmed:** Video generation writes `video_<nano>.<ext>` directly to `EngineConfig.OutputDir` and returns `VideoData`. The gallery save path only handles `ImageData`/`Images`; it does not accept `VideoData`. The video preview is a data URL, and there is no gallery sidecar or history entry for it. Deleting “current” only operates on `currentResult.metadata.file`, which a video result does not have.

**Source-confirmed:** `GetGalleryImage(filePath)` reads the supplied path without a gallery containment check. This is an ownership/security boundary that must be tightened when media roots are registered.

### 2.6 Wipe and destructive-operation ownership

**Source-confirmed:** `WipeAllData` deletes `exeDir/chats`, but production chats are under `%APPDATA%/HOWLChat/chats`. The backend wipe therefore does not target the actual chat root. `app-wipe.js` separately clears browser keys and IndexedDB and then reloads.

**Source-confirmed:** `WipeAllData` also deletes executable-relative `Characters`, `Worlds`, `Scenarios`, and `images/*` directories. Those are not the primary character/world/scenario stores, and they may not be the configured image output directory. The operation's scope does not match its UI promise.

**Inferred:** A future root-based cutover must make wipe/cleanup registry-aware. Deleting every registered external root would be unacceptable; the operation must delete only Howl-owned content after explicit confirmation, and must never delete an external registration's files.

---

## 3. Risks and inconsistencies

| Finding | Classification | Consequence |
|---|---|---|
| Library JSON is written automatically but never read | Source-confirmed | Users see files that look authoritative; they are not. Deletions and renames leave stale copies. |
| No single save result across two stores | Source-confirmed | Partial success is possible; world/character success messaging can hide a failed library write. |
| Save-path title collision replaces existing file | Source-confirmed | Two different items with the same title can overwrite one library file. `RoleSave` currently also permits replacement. |
| `RoleSave` is treated as sufficient to replace an existing file | Source-confirmed | A save-only library can be modified; the roles are not independently meaningful for replacement. |
| External registrations disappear from the grid when the root is unavailable | Source-confirmed | `GetContentLibraryCharacters` skips unreadable roots and keeps no last-known item cache. Registration survives, but the visible card does not. |
| Chat embeds full creative snapshots | Source-confirmed | Chat context can be stale relative to the live world/character catalog; the data is duplicated in AppData. |
| Chat cache and backend save are asynchronous | Source-confirmed | The cache can show a saved state before the backend has accepted it; a backend failure is not transactional. |
| Gallery image and metadata are not one atomic unit | Source-confirmed | Orphan images or orphan sidecars can result. |
| Gallery save failure is non-fatal at the app boundary | Source-confirmed | The generation UI can imply persistence that did not happen. |
| Gallery delete suppresses errors | Source-confirmed | UI reports success for a failed delete. |
| Video output has no catalogue entry | Source-confirmed | A generated video can be unreachable from history and has no delete contract. |
| Wipe targets the wrong chat root | Source-confirmed | “Wipe all” may leave actual AppData chats behind. |
| Content and gallery paths depend on project/executable location | Source-confirmed | Moving the install or changing launch directory can strand content. `paths.ProjectRoot` is not a user-content contract. |
| `GetGalleryImage` accepts arbitrary paths | Source-confirmed | The gallery read boundary is not enforced. |
| No user-facing import/export for chats, characters, worlds, scenarios, or lorebooks | Source-confirmed from inspected bindings | Existing content is trapped in browser stores; only `SaveImageAs` exports a generated image. |
| F2-B copied full image data into JSON files | Source-confirmed | Duplicated large blobs, slower saves, and a second copy that is never maintained. |

**Unknown:** The exact packaged/install behavior of `paths.ProjectRoot`, Windows known-folder redirection, and externally removed roots was not exercised at runtime in F2-C.

---

## 4. Recommended canonical ownership architecture

### 4.1 Core rule

**Source-confirmed problem:** Howl currently has multiple stores that each believe they own the same content.

**Recommendation:** One content item has exactly one authoritative owner. The content-library registry is the address/permission system, not an additional copy of the content. Browser localStorage/IndexedDB becomes a rebuildable cache/index once a category is cut over.

### 4.2 Ownership classes

| Class | Owner | Permitted operations | Must not |
|---|---|---|---|
| Howl-managed | Howl, inside a registered `howl` root | list, read, create, edit, delete through the category adapter | be silently copied into a second Howl store |
| Imported | Howl, copied into a `howl` or approved save root | list, read, edit, delete; record original source | modify or delete the source file |
| External read-only | The external library | list, read, use; show unavailable state | write, rename, delete, or imply editability |
| External read/write | The external library, by explicit role | list, read; explicit edit/delete only when the role and user confirmation permit it | assume that `Read` implies `Write`; delete the root |
| Generated | Howl generation pipeline | create, list, read, delete, export | be edited as if it were a source document |
| Exported copy | The user-chosen destination | user-initiated export only | be indexed, watched, re-imported automatically, or treated as live |
| Internal application state | Howl AppData/config | settings, pointers, caches, migration journals | be presented as user content or deleted by a content-root cleanup |

### 4.3 Role semantics to implement before cutover

**Source-confirmed:** `RoleRead`, `RoleWrite`, and `RoleSave` exist, but `WriteFile` currently allows `RoleSave` to replace an existing file.

Recommended contract:

- `Read`: enumerate and load.
- `Write`: modify or delete existing items, with atomic replacement and an explicit user action.
- `Save`: create a new item at a new path; do not silently replace an existing item.
- `DefaultSave`: exactly one per category for ordinary creation. Multiple readable roots are allowed. A save operation that targets an existing item must either use `Write` or produce a collision error/save-as choice.
- Root removal never deletes files.
- External `Read` is the default. External `Write`/`Save` require explicit selection and user confirmation.

### 4.4 Source identity and availability

Every listed item should carry:

- `owner`: Howl, imported, external, generated, or exported.
- `rootId`: the library registration that supplied it.
- `relativePath`: stable path inside the root.
- `format` and `revision` (content-derived or metadata-derived).
- `available`: false when the root cannot be reached.
- `lastSeenAt`: for index display, not authority.
- For imports: original source path and import time.
- For exports: no registration; the destination is recorded in the export result only.

**Recommendation:** The application index caches last-known external items so a disconnected drive greys them out instead of making them vanish. The cache is disposable and rebuilt from the root; it is never an alternate save location.

### 4.5 Save protocol

Every save returns a structured result, not a bare promise:

```text
{ category, id, owner, target, status: saved | saved-with-warning | failed,
  canonicalPath, cacheStatus, error }
```

Rules:

- The UI must not show plain success when the authoritative write failed.
- A cache write is not a save.
- A library export is not a canonical save.
- Retry is explicit; no silent background synchronization.
- Writes use temp file + verification + atomic rename where the filesystem permits.
- A partial multi-file save (chat, image plus sidecar) keeps a journal and surfaces recovery.

### 4.6 What F2-B's files should mean now

**Recommendation:** Reclassify the current automatic JSON files as **untrusted compatibility snapshots**, not canonical content and not user-requested exports. Do not read them back automatically. Before any migration, stop the automatic dual write and expose an explicit “Export a copy” action, or leave the files untouched but clearly mark the feature as export-only.

This avoids building F3 on top of write-only files that look authoritative.

---

## 5. Per-category integration design

### 5.1 Characters

**Current:** localStorage metadata, `CharacterStorage` IndexedDB images, external card discovery, and a write-only library JSON copy.

**Recommended canonical shape:** one folder per Howl-managed character, containing `character.json` and its images/assets. An external SillyTavern JSON/PNG card remains in place and is indexed read-only. A Howl-managed library root may contain the same Howl format already used by the external scanner.

Design:

- The character catalogue reads registered roots and displays `owner`, root, format, and availability.
- Local storage/IndexedDB becomes a cache of the last catalogue snapshot and image references, never the owner.
- External cards are displayed and usable without conversion; edit/delete is disabled unless the root has an explicit write role.
- Import is a copy into a Howl save root, with source metadata and a preview. The original is never changed.
- Images for a Howl-managed character live beside its card; a sidecar or manifest links them. Do not store full blobs in a second JSON index.
- A Howl-managed character delete uses the registry writer with explicit confirmation; external delete remains unavailable.
- The native `Characters/` prompt fallback must be retired or registered as a legacy read-only root after evidence of its use; it must not remain a hidden second source.

### 5.2 Worlds

**Current:** `howl-worlds` localStorage plus IndexedDB images; write-only `world.json`; embedded chat snapshot.

**Recommended canonical shape:** one Howl-managed `world.json` per world in a registered root, with image assets in the item folder or a dedicated media root and stable references. External/imported worlds use the same adapter but a different owner.

Design:

- The world grid lists from registered roots through a world adapter, not localStorage.
- Save uses `Save` for a new world and `Write` for an existing world.
- Images are references, not base64 in every save payload; the renderer resolves them through the media store.
- Chat stores a context snapshot plus the world source ID and revision. It must show that the chat is using a snapshot when the source has changed.
- Do not create a separate world JSON copy in the chat folder without a defined snapshot policy.

### 5.3 Scenarios

**Current:** `howl-chat-scenarios` localStorage plus `scenario-images`; write-only `scenario.json`; embedded chat snapshot.

**Recommended canonical shape:** one `scenario.json` per scenario, with character/lorebook references by stable ID and image assets resolved through the media store. Scenario is closely analogous to world but has more references.

Design:

- Scenario ID is the stable identity; title changes must not change the path.
- Reference resolution is explicit: active/removed/added character and lorebook IDs belong to the chat/session state, not duplicated full objects.
- The scenario adapter must reject dangling references or display them as unavailable, not silently substitute.
- Import/export is explicit and versioned.

### 5.4 Lorebooks

**Current:** `howl-chat-lorebooks` localStorage; write-only `lorebook.json`; chat `lorebooks.json` snapshot.

**Recommended canonical shape:** one JSON file per lorebook entry, or one versioned lorebook document per collection if a collection concept is added. Entries remain addressable by stable ID.

Design:

- This is the lowest-risk first cutover because there are no image blobs and the schema is already JSON.
- The lorebook page reads the registered root; localStorage becomes a cache.
- Chat stores only the selected lorebook IDs plus any frozen injected text needed for historical reproduction. It must not duplicate the full lorebook collection as a second editable source.
- Import/export follows an explicit format version, with collision handling and no automatic overwrite.

### 5.5 Chats

**Current:** AppData folder per chat; append log plus JSON snapshots; frontend cache; embedded full context.

**Recommended canonical shape:** retain the folder-per-chat format, but make the folder root a registered chat library. The folder is the canonical unit because a chat is a multi-file record, not a single JSON document.

Canonical files:

```text
<chat-root>/<chat-id>/
  meta.json          # identity, name, timestamps, references
  history.ndjson     # authoritative conversation append log
  history.json       # optional fallback snapshot, never preferred
  context.json       # frozen chat context snapshot + source revisions
  memories.json      # chat-owned memory state
```

Design:

- Register the existing AppData chat root as the first `howl` chat root, in place. Do not move user chats.
- New chats go to the selected default save root. Multiple chat roots are readable; the active chat pointer stores `(rootId, chatId)`.
- `history.ndjson` remains the only authoritative history. `history.json` is a recovery artifact, not a second source.
- Context snapshots are explicit and versioned. The UI shows when a chat uses an older snapshot than the live world/character.
- The frontend `HOWLStorage` current-chat cache is disposable and cannot be used to declare a save successful.
- `SaveChatData` must await the backend and return a structured failure; the UI must not write “saved” from the cache only.
- Delete removes one chat folder through the registry; it never touches the chat root or external chats.
- The AppData config file keeps only the active-chat pointer, schema version, and migration journal. It is internal state, not a chat.

### 5.6 Images

**Current:** gallery directory; PNG plus JSON sidecar; absolute path recorded in metadata; no library registration; delete suppresses errors.

**Recommended canonical shape:** a generated-media root registered for the Images category. Each generated image is a pair:

```text
<image-root>/<id>.png
<image-root>/<id>.json
```

The sidecar contains stable ID, prompt, model, parameters, source references, created time, and a relative filename. The absolute path may be retained only as diagnostic data.

Design:

- The gallery lists across registered image roots; the selected save root is the only write target.
- A new generation writes to a temp file, verifies the image, writes the sidecar, then commits both. A failed sidecar is a failed save and the UI says so.
- Delete uses the sidecar identity and reports actual errors.
- `GetGalleryImage` resolves only inside a registered image root; arbitrary paths are rejected.
- Read-only image roots can be browsed and used as inputs; they cannot be deleted or overwritten.
- `SaveImageAs` remains an explicit export to a user-selected path and is not indexed.

### 5.7 Videos

**Current:** direct clip write to the image output directory; base64 returned to the UI; no sidecar, gallery listing, or delete contract.

**Recommended canonical shape:** a generated-video root registered for the Videos category, using the same paired media pattern:

```text
<video-root>/<id>.<container>
<video-root>/<id>.json
```

Design:

- `generateVideo` must not write an untracked file. It should return bytes/path to a media service that commits clip plus sidecar atomically.
- Video metadata includes container, duration/frames if known, model, prompt, generation options, and source references.
- The history/gallery surfaces videos as first-class generated content.
- Video export is an explicit copy. Video deletion uses the media service.
- Images and videos may share a media index, but their roots, schemas, and write/delete behavior stay distinct enough to avoid forcing image assumptions onto video.

### 5.8 Why not one identical store for every category

- A character is a folder plus assets.
- A world/scenario/lorebook is a document with references.
- A chat is an append log plus metadata and snapshots.
- An image/video is a binary plus sidecar.
- Internal state is not user content.

The unified part is the **address, permission, identity, and save-result contract**, not one file format.

---

## 6. Migration requirements and rollback strategy

### 6.1 Is migration necessary?

| Category | Migration recommendation | Reasoning |
|---|---|---|
| Chats | Register existing AppData root in place; no file move initially | Current format is already durable and multi-file. Moving first adds risk without archival benefit. |
| Images | Register the existing gallery directory in place | PNG/JSON pairs already exist; ownership can be added without moving files. |
| Videos | No migration until current outputs are inventoried; new outputs use the media contract | Existing videos have no sidecars and may be orphaned. Do not guess ownership. |
| Characters | One-time explicit import from browser storage and, if used, from the legacy native `Characters/` directory | Browser records are not user-visible files. Preserve source until verification. |
| Worlds/scenarios/lorebooks | One-time explicit migration from browser storage to canonical JSON files | These are browser-only today; no safe filesystem source exists. |

No migration is authorized or performed by F2-C.

### 6.2 Preflight

1. Resolve every selected root and confirm it exists, is a directory, and is writable only if it is a save/write target.
2. Reject a destination inside an external read-only root.
3. Check free space and path length; never scan E: or unrelated drives.
4. Enumerate source records without modifying them.
5. Verify IDs, required fields, schema version, duplicate IDs, and image references.
6. Detect collisions using case-insensitive normalized paths on Windows.
7. Build a plan and show a count plus a per-item list before any write.
8. Record a journal under AppData, not inside the destination root.

### 6.3 Identity and collision

- Source identity is `(category, source kind, source path/key, record ID)`.
- Destination identity is `(rootID, category, stable item ID)`.
- If the same ID and same verified content exist: skip and report.
- If the same ID but different content exists: do not overwrite. Require a user choice: keep existing, import as a new ID, or cancel.
- If two source records have the same title: never use title alone as the path; include a stable ID or suffix.
- External roots are never migration destinations.

### 6.4 Verification

- JSON: parse, verify required fields, verify ID and schema version.
- Images: verify the file exists, is non-empty, and decodes as the expected format where practical; verify the sidecar points to that file.
- Chats: verify `meta.json`, count `history.ndjson` lines, compare with a parsed snapshot, and verify no unknown required file is missing.
- Record a content hash for JSON and small media; do not hash huge model files. Hashing is migration-only and bounded to the item being copied.
- A migration is not marked complete if any item is unverified; report partial completion explicitly.

### 6.5 Failure recovery and rollback

- Write each new item to a staging directory under the destination.
- Verify the staged item.
- Commit with atomic rename(s).
- Update the migration journal only after verification.
- On failure, delete only staged files created by this migration; never delete source records or pre-existing destination items.
- Rollback restores the active pointer to the legacy store and leaves copied destination items untouched for inspection or explicit cleanup.
- The old browser store stays readable for at least one release. Do not clear localStorage/IndexedDB automatically.
- After a successful migration, mark the old store as legacy/read-only; delete it only through a separate, explicit cleanup operation.

### 6.6 Old-location handling

- **Chats:** AppData remains a registered root; no automatic move/delete.
- **Images:** existing gallery remains registered; no automatic move/delete.
- **Characters/worlds/scenarios/lorebooks:** keep the old browser data as a legacy backup; do not delete until the user has verified the new library and a separate cleanup flow is approved.
- **Videos:** leave unknown files untouched and classify them as unindexed legacy media until a user explicitly adopts or exports them.

---

## 7. Runtime acceptance plan

These are manual Wails procedures. **F2-C did not perform them; none is Runtime-confirmed.**

| # | Test | Steps | Expected pass | Failure condition |
|---|---|---|---|---|
| 1 | Add external library | Add a temporary character folder as an external readable root | Root appears with source/origin and read-only role | Root missing, path lost after restart, or UI implies write access |
| 2 | Remove registration | Remove the external root | Registration disappears; underlying files still exist | Any file is deleted or root remains registered |
| 3 | External discovery | Place Howl JSON, SillyTavern JSON, and PNG cards in the root; open Characters | Cards appear with source/format and read-only state | Cards vanish, duplicate under overlapping roots, or become editable local cards |
| 4 | Read-only enforcement | Open an external card; attempt edit/save/delete | Fields disabled; no save/delete action; original unchanged | Any write reaches the external file |
| 5 | Save destination | Select a save root for each category; save one item | Item reloads from the page store and the selected destination receives the expected file | Wrong root, silent fallback, or success shown after failed write |
| 6 | Delete after save | Delete a local item that has a library copy | Local item disappears; library copy handling matches the documented mirror/export policy | Orphan copy is silently presented as live content |
| 7 | Model intelligence | Select models with and without sidecar metadata | Evidence, capabilities, runtimes, tokenizer, source, notes display; unknown remains unknown | Filename-only certainty, merged chat/image lists, or unhandled conflict |
| 8 | Reset settings | Reset Application Settings with image settings active | Settings and image settings reset; content, registrations, models, and generated files survive | User content lost, registration removed, or image settings silently retained |
| 9 | Failed save | Make the selected destination read-only/unavailable, then save | Local record remains; UI clearly reports the failed destination | Plain success, silent fallback, or lost local record |
| 10 | Unavailable drive | Add a removable-root fixture, disconnect it, refresh the list | Registration retained; known items visible as unavailable if indexed | Items disappear without explanation or registration is removed |
| 11 | Reconnection | Restore the same root, refresh | Items return with no duplicates or identity changes | New IDs, duplicate cards, or lost registrations |
| 12 | Chat lifecycle | New chat, send messages, autosave, Save As, load, delete | History and context match; failures are surfaced; delete removes only that chat | Cache-only success, lost history, wrong chat deleted, or AppData folder left inconsistent |
| 13 | Image lifecycle | Generate an image, list it, reopen it, delete it | PNG and sidecar agree; delete reports errors | Image visible without sidecar, success on failed delete, or arbitrary path read |
| 14 | Video lifecycle | Generate WebM/AVI/WebP, then find and delete it | Clip and sidecar are catalogued and deletable | Clip only in transient preview, no history entry, or untracked file |
| 15 | Wipe/cleanup | Run the explicitly named wipe/reset | Only intended Howl-owned data is removed; external roots and files survive | AppData chats remain or external/user content is removed |

---

## 8. Dependency map

```text
internal/paths
  DefaultContentDir / known folders
        |
        v
internal/library.Registry
  roots, roles, origin, format, default save, source identity
        |
        +--> character adapter ----> character page / external cards
        +--> world adapter --------> world page
        +--> scenario adapter -----> scenario page
        +--> lorebook adapter -----> lorebook page
        +--> chat root adapter ----> chat persistence manager
        +--> media root adapter ---> gallery / video service
        |
        v
AppData/config
  active-chat pointer, indexes, migration journal, settings,
  model cache, temp state
```

Key dependencies and blockers:

- **Source-confirmed:** `paths.ProjectRoot` is executable/project-relative and cannot be the long-term owner of user content.
- **Source-confirmed:** `chat.NewPersistenceManager` hardcodes AppData, while `NewPersistenceManagerWithDir` proves a root parameter is feasible. The production constructor must accept a registered root before chat cutover.
- **Source-confirmed:** `gallery.Service` and `sdcpp.EngineConfig.OutputDir` must be kept in sync; currently video writes directly to the engine output while images are catalogued by the gallery.
- **Source-confirmed:** `SaveContentItem` has no chat/image/video implementation; those categories need dedicated adapters, not a generic JSON writer.
- **Source-confirmed:** `content-save.js` is a browser bridge, not a canonical store. It must not remain the long-term save owner.
- **Source-confirmed:** `app-wipe.go` and `app-wipe.js` use different roots and different data sets; they need one registry-aware operation.
- **Source-confirmed:** model intelligence is independent of content ownership and should not be folded into this cutover.
- **Source-confirmed:** the TTS contract and audio pipeline are outside this design and must remain untouched.

---

## 9. Requirements ledger reconciliation

**Source-confirmed:** `PROJECT_ISSUES_AND_EXPANSION.md` remains the single authoritative requirements ledger. The 16 F2-B future-requirement rows were preserved.

Roadmap ownership already recorded there:

- `FUT-C1`–`FUT-C5`: chat cues, image cues, llama.cpp features, architecture/family/era options, and sampling research — future chat capability phase.
- `FUT-ST1`–`FUT-ST5`: SillyTavern character import/mapping/Description/LoRA/lorebook configuration — future compatibility phase.
- `FUT-CH1`–`FUT-CH3`: Author’s Note, single-character context capacity, multi-character isolation — future chat behavior phase.
- `FUT-IMG1`–`FUT-IMG3`: RAM/VRAM placement, safe offload, OOM prevention — future image resource phase.

**F2-C decisions to record:**

| ID | Decision | Status |
|---|---|---|
| F2-C1 | One authoritative owner per content item; the library registry is an address/permission layer, not a second content store. | DESIGN RECORDED |
| F2-C2 | F2-B automatic JSON copies are untrusted export/compatibility snapshots, not canonical content; they must not be read back automatically. | DESIGN RECORDED |
| F2-C3 | Browser localStorage/IndexedDB becomes a rebuildable cache after each category cutover; never the canonical owner. | DESIGN RECORDED |
| F2-C4 | Chats retain their multi-file folder format; the existing AppData root is registered in place before any move. | DESIGN RECORDED |
| F2-C5 | Images and videos are binary-plus-sidecar generated media; both need stable catalogue identity and atomic commit. | DESIGN RECORDED |
| F2-C6 | `Save` is create-only; replacing an existing item requires `Write` or an explicit collision choice. | DESIGN RECORDED |
| F2-C7 | No automatic migration or background synchronization; migration is explicit, verified, journaled, and reversible. | DESIGN RECORDED |
| F2-C8 | Disconnected roots retain registrations and last-known index entries marked unavailable. | DESIGN RECORDED |

No future requirement was implemented or silently promoted to complete during F2-C.

---

## 10. Recommended implementation sequence

### Smallest safe next slice

**F3-A — Correct the mirror semantics before any ownership move.**

- Stop treating the F2-B JSON writes as saves.
- Mark them as export/compatibility snapshots or replace them with an explicit “Export a copy” action.
- Do not read them back.
- Add a structured save result so a failed destination write cannot be presented as plain success.
- Add a read-only inventory for registered roots and last-known unavailable items.

This is smaller and safer than starting a migration: it removes the false dual-source signal without moving any user data.

### Then, in order

1. **F3-B — Lorebooks:** canonical JSON file store, explicit import from localStorage, read-back through the library adapter, legacy store retained.
2. **F3-C — Worlds and scenarios:** canonical JSON plus media references; explicit migration; chat snapshots given a visible revision policy.
3. **F3-D — Characters:** canonical folder plus images; external read-only cards; explicit import; no automatic copy of external cards.
4. **F3-E — Chats:** register AppData in place, support multiple chat roots, fix the active pointer, make history/context ownership explicit, and surface save failures.
5. **F3-F — Images:** register the gallery in place, tighten path containment, make PNG+sidecar commit atomic, and fix delete reporting.
6. **F3-G — Videos:** introduce the media sidecar contract, catalogue clips, and provide export/delete.
7. **F3-H — Wipe/cleanup:** registry-aware, root-aware, explicit, and separate from Reset Application Settings.

### Acceptance gate for the first cutover

Do not remove localStorage/IndexedDB or delete any legacy source until:

- the canonical read path is proven,
- the copy/migration verification passes,
- the runtime acceptance tests pass,
- external registrations survive disconnect/reconnect,
- and a rollback has been exercised.

**F2-C is complete as a design only. It does not claim that any cutover, migration, or runtime acceptance test has occurred.**
