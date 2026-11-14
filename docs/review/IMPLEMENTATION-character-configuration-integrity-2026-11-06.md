# Implementation — Character Configuration Integrity

**Date:** 2026-11-06  
**Status:** Implemented and locally verified; not committed or pushed.  
**Scope:** Stage 2 only. No character-specific TTS engine selection, voice resolution, synthesis changes, cache redesign, performance benchmarking, new voice settings, Character Tab redesign, or chat-bubble controls were implemented.

## 1. Result

The Character Tab no longer rebuilds a character from only the visible form fields. It now writes an edited patch, merges that patch into the existing persisted character record, and preserves unrelated top-level, metadata, `howl`, nested voice, image, and unknown fields.

The shared `CharacterStorage.save` path also merges against the persisted record. This protects alternate save paths that pass a derived or partial character object, such as the sidebar image actions.

No character schema change was required. Character IDs remain unchanged during edits.

## 2. Files changed

| File | Change |
|---|---|
| `frontend/dist/js/character-form.js` | Collects an edit patch instead of a fresh full object; deep-merges it with the existing record before saving; preserves `scenario`, `tags`, metadata timestamps/seed, world, memories, notes, image style, advanced voice fields, and unknown fields; tracks changed voice controls so an unrepresentable stored value is not blanked accidentally; returns the save promise. |
| `frontend/dist/js/character-storage.js` | Added `mergeRecord`; `save` now merges into the existing persisted record before writing; only replaces IndexedDB image blobs when the incoming record contains an `images` field; exported `mergeRecord` for the form path. |
| `frontend/dist/js/character-state.js` | Added per-field dirty tracking used by the voice controls. |
| `tests/frontend/character-configuration-integrity.test.js` | New focused tests using the actual `saveCharacter` and `CharacterStorage.save`/`loadAll` paths with representative records. |
| `docs/review/IMPLEMENTATION-character-configuration-integrity-2026-11-06.md` | This report. |

No Go/backend, audio, chat-history, or Stage 1 identity files were changed.

## 3. Original data-loss mechanism

**Observed in source.**

`collectCharacterData()` in `frontend/dist/js/character-form.js` constructed a complete new object for every save. On an edit it reset or replaced:

- top-level `scenario` with `''` (there is no `character-scenario` field in the current form);
- `tags` with `[]`;
- `metadata.version`, `metadata.created`, and `metadata.seed` with new values;
- `howl.world` with `''`;
- `howl.memories` with `[]`;
- `howl.notes` with `''`;
- `howl.image_style` with `'fantasy-realistic'`;
- `howl.voice` with only `enabled`, `engine`, `voice`, `speed`, and `volume`.

`saveCharacter()` preserved only the existing character ID before calling `CharacterStorage.save()`.

`CharacterStorage.save()` then replaced the entire localStorage record with that reconstructed object. It also unconditionally wrote `character.howl?.images || []` to IndexedDB, so a partial save without an image intent could clear existing image blobs.

A second problem affected existing voice values. `editCharacter()` loads `voice.voice` into a fixed-option `<select>`. A stored voice ID not present in that select resolves to `''` in the browser, and the old save path then wrote that empty value back.

## 4. Authoritative character representation

**Observed and source-derived.**

The Character Tab's authoritative persisted representation is:

- **Metadata:** localStorage key `howl-chat-characters`, one JSON array of character records.
- **Image blobs:** IndexedDB database `howl-character-images`, object store `images`, keyed by `charId`.
- **Hydrated view:** `CharacterStorage.loadAll()` reads localStorage, loads each character's images from IndexedDB, and returns a combined object. `characterState.characterCache` is a derived cache of that result.

`CharacterStorage.save()` is the shared write point for the loaded Character Tab modules.

Other character concepts were checked:

- `frontend/dist/js/character-manager.js` contains another localStorage-backed manager, but `character.html` does not load it. It is not an active Character Tab path.
- The backend `Characters/<name>/character.json` path exists through `app-characters.go`, but the Character Tab does not call `CreateCharacter`; its active save path is frontend localStorage/IndexedDB. No synchronization between those two concepts was introduced.
- No import/export or preset path was found for the Character Tab.

No conflicting active Character Tab source of truth was found, so the stop condition for conflicting persistence ownership was not triggered.

## 5. Before/after update behavior

### Before

```text
form fields
  -> collectCharacterData() creates a fresh object
     -> resets unrepresented fields and visible-only voice object
  -> CharacterStorage.save() replaces the whole stored record
  -> localStorage metadata + IndexedDB images
```

### After

```text
form fields
  -> collectCharacterData(existing) creates an edit patch
     -> only visible fields + explicit list changes
     -> voice fields only when changed
  -> mergeRecord(existing, patch)
     -> recursive object merge; arrays/scalars replace
  -> CharacterStorage.save(merged)
     -> merges once more against the latest persisted record
     -> writes images only when an images intent is present
  -> localStorage metadata + IndexedDB images
```

The form-level merge keeps the returned/posted character complete for the current page. The storage-level merge protects other callers and uses the latest persisted record rather than trusting a potentially partial caller object.

## 6. Field-preservation strategy

### Objects

`mergeRecord(existing, patch)` recursively merges plain objects. This preserves siblings when one nested property changes:

- `metadata.age` changes without replacing `metadata.created`, `metadata.seed`, or `metadata.extraMetadata`;
- `howl.world` remains an object with its nested properties;
- `howl.voice.noiseScale`, `noiseW`, `sentenceSilence`, and unknown voice properties survive edits to name or visible voice fields.

`null` is distinct from omission. `mergeRecord` skips only `undefined`; an explicitly supplied `null` is written as the new value. The focused test confirms that an omitted nested property keeps its non-null value, while an explicit `null` is persisted and reloads as `null`.

### Arrays and scalar values

Arrays and scalar values are replaced, not recursively merged. This is intentional:

- `howl.images: []` explicitly removes images.
- `howl.lorebooks: []` explicitly removes lorebook references.
- Empty strings clear scalar fields such as description or role.
- `tags` is preserved for edits because it is not exposed in the current form; a new character still starts with `[]`.

### Voice controls

Only voice fields whose controls were actually changed are written into the patch. `characterState` tracks dirty field IDs from input/change events.

This matters because the engine/model selects are fixed option lists. If a stored voice/engine is not represented by an option, the DOM value is empty. Without dirty tracking, merely saving an unrelated field would erase the stored value. With tracking:

- an untouched, unrepresentable stored voice survives;
- choosing the empty “Inherit Global” option marks the field dirty and explicitly clears it;
- editing speed/volume does not erase the untouched engine/model;
- advanced voice parameters not shown in the form survive because the voice object is deep-merged.

### Images

The form now copies full image entries with object spread, preserving fields such as `medium` and `original`, rather than rebuilding only `id`, `name`, `size`, `data`, and `isThumbnail`.

`CharacterStorage.save()` writes full image arrays to IndexedDB when the incoming patch has an `images` property. It does not touch IndexedDB when the caller has no image intent. This prevents the sidebar partial-save path from deleting blobs while still allowing an explicit `images: []` to clear them.

### Identity

`collectCharacterData(existing)` uses `existing.id` when editing. `saveCharacter` also restores it before merging. `CharacterStorage.save` upserts by id. The focused sequential-edit test confirms the ID is unchanged.

## 7. Persistence and reload behavior

- Metadata is written to localStorage after the merged record and images are prepared.
- Full image blobs are written to IndexedDB; localStorage receives metadata-only image entries via `stripImageData`.
- `loadAll()` reads the metadata, loads images by character ID, and restores `howl.images` plus the display `avatar`.
- A fresh page/`CharacterStorage` instance in the tests reads the same localStorage and IndexedDB state and reconstructs the complete character.
- No destructive migration or new required property was introduced. Legacy records without `howl`, `metadata`, or voice objects remain loadable; saving adds only fields the current form actually supplies.

## 8. Tests

### New test file

`tests/frontend/character-configuration-integrity.test.js` exercises the real form save function and the real storage load/save code in a Node VM. The harness provides a DOM-like form, localStorage, and an IndexedDB-compatible in-memory store. It does not just call a merge helper.

Representative fixture fields include:

- identity and descriptive fields;
- `scenario` and `tags`;
- metadata including `created`, `seed`, `age`, `role`, `background`, `priority`, and unknown nested metadata;
- `howl.world` with nested data;
- `howl.memories`;
- `howl.notes`;
- `howl.image_style`;
- full `howl.images` entries with thumbnail, medium, original, and unknown fields;
- voice settings plus advanced `noiseScale`, `noiseW`, `sentenceSilence`, and unknown nested voice data;
- unknown top-level and `howl` properties.

Tests cover:

1. Editing a scalar and reloading the complete record.
2. Editing nested metadata fields.
3. Editing one section while preserving other sections.
4. Explicitly clearing a scalar and an image list.
5. Preserving arrays and nested objects.
6. Preserving unknown properties.
7. Saving and reloading through a fresh storage instance.
8. Editing multiple fields sequentially with an unchanged ID.
9. Loading and editing a legacy record without newer optional properties.
10. Preserving an unrepresentable stored voice when the select displays empty.
11. Explicitly clearing the voice model.
12. A partial sidebar-style save preserving persisted metadata and image blobs.
13. Preserving an omitted nested property while applying an explicit `null` through save and reload.

### Actual results

```text
node --test tests/frontend/character-configuration-integrity.test.js
✔ Character Tab edit preserves unrelated, nested, unknown, and image data through save and reload
✔ a partial sidebar-style storage save preserves the persisted character and image blobs
✔ an omitted nested property is preserved while an explicit null is applied
✔ an explicit image-list clear still removes the persisted image blobs
✔ sequential edits preserve other sections and explicit scalar clears remain authoritative
✔ an explicit voice-model clear works while an unrepresentable stored voice is preserved
✔ legacy records without newer optional fields remain editable without destructive migration
ℹ tests 7
ℹ pass 7
ℹ fail 0
character_tests_exit=0
```

```text
node --test tests/frontend/chat-render-ownership.test.js
✔ explicit character identity selects the owning character
✔ missing character identity does not fall back to the first character
✔ unknown character identity is preserved and does not select another character
ℹ tests 3
ℹ pass 3
ℹ fail 0
stage1_node_tests_exit=0
```

Stage 1 regression coverage:

```text
go test ./internal/backend/chat -count=1
ok  	howl-chat/internal/backend/chat	1.066s
stage1_chat_tests_exit=0

go test . -run 'Test(CharacterOwnership|ReplaceChatHistory|LegacyHistory|ConvertFrontendHistory|TranscribedTextPersistsInChatHistory)' -count=1
ok  	howl-chat	2.268s
stage1_app_tests_exit=0
```

Syntax and changed-source whitespace checks:

```text
node --check frontend/dist/js/character-form.js
node --check frontend/dist/js/character-storage.js
node --check frontend/dist/js/character-state.js
syntax_status=0

git diff --check -- frontend/dist/js/character-form.js frontend/dist/js/character-storage.js frontend/dist/js/character-state.js
scoped_diff_check_exit=0
```

A full-worktree `git diff --check` still exits `2` because of pre-existing trailing whitespace in unrelated modified files, including `IMAGE_PAGE_CHANGE_LEDGER.md` and previously modified generated content in `frontend/wailsjs/go/models.ts`. Those lines were not introduced by Stage 2. The changed-source scoped check is clean.

## 9. Compatibility and protected functionality

- **Stage 1 character identity:** untouched. The Stage 1 frontend ownership tests still pass, and no chat-history or backend identity file was changed.
- **Chat history and restoration:** untouched.
- **Group 0 settings contract:** untouched.
- **Group 1 Whisper decode settings:** untouched.
- **Group 2A TTS speed:** untouched.
- **Accepted WASAPI microphone gate:** untouched.
- **No-auto-load behavior:** untouched.
- **TTS independence from ASR:** untouched.
- **Headphone-only playback:** untouched.
- **No new voice settings or voice routing:** no engine discovery, voice resolution, synthesis, cache, or Character Tab redesign was added.

## 10. Limitations and unverified paths

1. **No browser E2E or native close/relaunch test was run.** The new tests run the actual form/storage functions in a VM with localStorage/IndexedDB-compatible fakes. They simulate a fresh storage/page instance, but not a real Wails process restart or a real browser profile.
2. **The fake IndexedDB exercises the implemented transaction API, not browser-specific storage behavior.**
3. **The legacy `character-manager.js` remains a lossy alternate manager in source.** It is not loaded by `character.html`; re-enabling it would require separate review.
4. **The backend `Characters/<name>/character.json` path remains independent.** It was not changed or synchronized with the Character Tab store.
5. **Top-level legacy `avatar` blobs were not separately modeled.** The active image representation is `howl.images` plus IndexedDB; the tests cover that representation. An old record containing only a top-level avatar data URL was not exercised.
6. **No historical production corpus was available.** Compatibility was tested with representative legacy records and the current storage code.

## 11. Evidence classification

### Observed

- The focused frontend tests pass with the actual `saveCharacter` and `CharacterStorage` functions, including the omitted-versus-explicit-`null` nested-property test.
- The Stage 1 Node and Go tests pass.
- Syntax checks and the changed-source diff check pass.
- The full-worktree diff check reports pre-existing whitespace in unrelated files.
- No commit or push was performed.

### Source-derived

- `howl-chat-characters` plus IndexedDB is the active Character Tab persistence path.
- The old `collectCharacterData` rebuild was the data-loss mechanism.
- The sidebar image path can pass a partial derived character, so storage-level merging is needed in addition to the form-level fix.
- Deep object merging preserves unknown/nested data while array/scalar replacement preserves explicit clears.

### Calculated

- The test suite contains seven Stage 2 cases and three Stage 1 renderer cases; all pass.
- The changed-source scoped diff check returns `0`.

### Unverified

- Real browser localStorage/IndexedDB behavior.
- Native application close/relaunch.
- Every historical character-file variant and every production record shape.
- Visual Character Tab appearance; no UI was redesigned.

## 12. Authorization boundary

Stage 2 only. No backend character voice architecture, voice resolution, synthesis, cache correctness, performance benchmarking, Character Tab voice controls, or chat-bubble controls were implemented.

No commit or push was performed.
