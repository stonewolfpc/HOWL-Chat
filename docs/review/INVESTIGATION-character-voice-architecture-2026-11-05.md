# Group 2B — Character Voice Architecture Investigation

**Date:** 2026-11-05

**Authorization:** Read-only investigation. No implementation authorized.

**Status:** Investigation complete and accepted on 2026-11-06. Group 2B is closed as an accepted investigation. Implementation remains blocked on character identity and persistence prerequisites. No production code, test code, settings, models, builds, or audio artifacts were changed or produced.

## 1. Executive answer

A per-character voice path already exists, but it is a renderer-side partial implementation rather than a complete character-voice system.

- `audio-service.js` reads `character.howl.voice`, merges it with global audio settings, and calls `SynthesizeSpeechWithVoice`.
- `app-audio.go` parses that JSON into `VoiceOverride`.
- `Pipeline` applies only `Voice` and `Speed` from the override.
- `VoiceOverride.Engine` is parsed but never consumed. A character cannot select a different engine from the active global engine.
- Piper voice switching is implemented by lazily replacing the pipeline's synthesizer with another Piper instance, regardless of the active engine.
- The character ID never reaches the Go synthesis request; only the resolved voice string does.
- The Character tab already contains a voice block, but it is static and its save path reconstructs the character object, dropping unrelated and advanced fields.
- Historical assistant messages do not persist a character ID. On reload, unowned assistant messages are assigned to the first character.

The architecture can support multiple voices on one active engine with a bounded backend change. It cannot safely support independent per-character engines, full voice-specific parameters, reliable historical ownership, or portable character voice persistence until the prerequisites in this report are completed.

Three stop conditions are met:

1. **Historical bubble ownership is not reliable.** `ChatHistory` has no character-ID field, and the rebuild path falls back to the first character.
2. **The Character tab save path cannot safely preserve a character while editing voice.** It reconstructs `howl` and loses fields such as `world`, `notes`, `memories`, `image_style`, and advanced voice parameters.
3. **Per-character engine selection is incompatible with the current single-engine pipeline.** `VoiceOverride.Engine` is ignored and `ensureTTSVoice` always attempts a Piper load.

No UI or Group 2C implementation should begin until these are addressed.

## 2. Evidence classification

This report uses the following labels:

- **Observed:** directly verified in source or checked artifact.
- **Source-derived:** established by tracing the implementation, without executing it.
- **Measured:** backed by an existing measurement.
- **Inferred:** reasoned conclusion not directly demonstrated.
- **Unestablished:** requires a benchmark, live service, or product decision.

**Measured:** No synthesis latency, peak memory, VRAM, voice-switch time, or concurrency measurement exists in the inspected artifacts for this feature. The model-profile figures cited in section 9.2 are declared registry estimates, not measurements.

No synthesis or benchmark was run for this investigation. No MeloTTS service was contacted.

## 3. Existing Character tab architecture and UI inventory

### 3.1 Page structure

**Observed.** `frontend/dist/character.html` is a single Character page, not a tabbed editor. Its structure is:

- top menu and navigation sidebar;
- page header with Create New and User Details;
- character grid (`#character-grid`);
- character creation/edit overlay (`#character-creation-overlay`);
- one scrolling character form with identity, description, personality, background, dialogue examples, lorebooks, images, and voice;
- a separate User Details overlay;
- a tutorials overlay.

There are no Character sub-tabs and no route-level character editor.

### 3.2 Existing voice controls

**Observed.** The voice section already exists at `frontend/dist/character.html:226-266`. It contains:

- `#char-voice-enabled` checkbox;
- `#char-voice-engine` select with `""` (inherit global), `piper`, `kokoro`, `melotts`, and `browser`;
- `#char-voice-model` select with five hard-coded Piper names:
  - `en_US-lessac-medium`
  - `en_US-lessac-low`
  - `en_US-ryan-medium`
  - `en_US-amy-medium`
  - `en_GB-northern_english_male-medium`
- `#char-voice-speed`, range `0.5–2.0`;
- `#char-voice-volume`, range `0–1`.

There is no dynamic voice catalog, no MeloTTS speaker list, no Kokoro voice list, no language control, no Piper noise/width/sentence-silence controls, and no browser voice list.

**Source-derived.** The static model list is not populated from `GetTTSAvailableVoices` or `GetTTSModels`. The only dynamic Piper registry dropdown is the global audio-settings dropdown in `frontend/dist/js/audio-settings.js:970-1000`.

### 3.3 Edit and save flow

**Observed.** `frontend/dist/js/character-grid.js:114-140` loads an existing character into the form. It sets the voice fields directly:

```js
set('char-voice-model', voice?.voice || '');
```

Because `#char-voice-model` is a fixed `<select>`, a stored voice ID that is not one of the five hard-coded options cannot be represented. The browser silently resolves the select to an empty value, and a subsequent save can overwrite the stored voice with `""`.

**Observed.** `frontend/dist/js/character-form.js:51-101` rebuilds the entire character object on save. It creates a fresh `howl` object and sets:

- `world: ''`
- `lorebooks` from current selection
- `memories: []`
- `notes: ''`
- `image_style: 'fantasy-realistic'`
- `images` from current state
- `voice` with only `enabled`, `engine`, `voice`, `speed`, and `volume`

This is not a voice-only defect. It means editing a character's voice through the existing form can discard unrelated character data. It also discards advanced voice fields that are present in the backend `VoiceConfig` type but absent from the form.

**Observed.** `saveCharacter` uses `window.CharacterStorage.save`, then posts `character-created` or `character-updated` to `window.parent` only when the page is running inside a different window. The standalone Character page saves to localStorage and IndexedDB, then reloads the grid.

**Observed.** `frontend/dist/js/sidebar-panels.js:90-110` listens for those messages, but only when the current path includes `character.html`. The `character.html` script list does not load `sidebar-panels.js`, and no iframe host for `character.html` was found. The live chat page therefore has no working path to learn that a character voice changed without a reload.

### 3.4 Deep-link behavior

**Observed.** `app-characters.go:22-24` implements `OpenCharacterSheet(characterId)` by navigating to:

```text
character.html?id=<characterId>
```

**Observed.** `frontend/dist/js/character.js:22-26` loads the character grid on `DOMContentLoaded` but never reads the `id` query parameter and never calls `editCharacter`. Opening a character from the sidebar therefore lands on the character page without opening that character's sheet.

This matters for the proposed bubble menu: a "Voice settings" action cannot safely deep-link to the editor until this behavior is fixed.

### 3.5 Persistence and duplication

**Observed.** The Character tab's primary save path is frontend-only:

- metadata: localStorage key `howl-chat-characters`;
- images: IndexedDB `howl-character-images`, store `images`;
- `CharacterStorage.save` strips image data from localStorage and stores full image blobs in IndexedDB.

**Observed.** A separate backend path exists in `app-characters.go`:

- `CreateCharacter` writes `Characters/<name>/character.json`;
- `GetCharacters` reads those directories;
- `OpenCharacterSheet` navigates the UI.

No frontend call to `CreateCharacter` was found. The Character tab uses `CharacterStorage`, not the backend character-directory API. The checked `Characters/Aria/character.json` and `Characters/Doren/character.json` contain `{}`.

**Source-derived.** There are two character persistence concepts with no synchronization contract: localStorage/IndexedDB for the Character tab, and `Characters/<name>/character.json` for backend character loading. This is a pre-existing architectural split that a character-voice feature must not deepen.

### 3.6 Presets, imports, exports, and editing workflows

**Observed.** No Character import, export, or preset workflow was found in the Character tab. The tab supports image selection, lorebooks, manual creation, grid editing, and deletion. Image handling and lorebook references are the only existing attachment workflows.

**Source-derived.** If a future character preset/import/export feature is added, it must include the canonical voice object. The current save reconstruction is already lossy, so it cannot be used as a format authority.

### 3.7 Unused/legacy character manager

**Observed.** `frontend/dist/js/character-manager.js` implements another localStorage-backed character manager with its own load/save/edit methods. It is not loaded by `character.html`'s script list. It should not be used as a second source of truth for voice.

### 3.8 Correct placement for voice controls

**Source-derived.** The existing voice section is the natural place for character voice configuration. It already sits inside the character sheet and is loaded by the existing edit flow.

The correct change is to make that block engine-aware and data-driven, not to add a new global audio panel or a separate chat-only workaround. Before changing it, the save path must preserve the rest of the character object and the select must be populated from backend voice discovery.

## 4. Character-sheet data model and persistence

### 4.1 Frontend character shape

**Observed.** The Character tab saves objects shaped approximately as:

```js
{
  id, name, description, personality, scenario, first_mes, mes_example,
  avatar, tags[], metadata{...},
  howl: {
    world, lorebooks[], memories[], notes, image_style, images[],
    voice: { enabled, engine, voice, speed, volume }
  }
}
```

The active per-character voice is `character.howl.voice`, not a top-level `voice` field.

### 4.2 Backend character shape

**Observed.** `internal/backend/types/entities/entities.go:19-48` declares `VoiceConfig` with:

- `Enabled`
- `Engine`
- `Voice`
- `Speed`
- `Volume`
- `NoiseScale`
- `NoiseW`
- `SentenceSilence`

`entities.Character` stores `VoiceConfig map[string]interface{}` under JSON key `voice_config`.

**Source-derived.** `ParseCharacterContextFromMap` only handles image extraction and metadata hoisting. It does not map `howl.voice` into `VoiceConfig`. No production code populates `Character.VoiceConfig` from the frontend. The field is effectively unused.

### 4.3 Two incompatible JSON shapes

**Observed.** The frontend voice object uses camelCase keys such as `noiseScale`, `noiseW`, and `sentenceSilence`.

`VoiceOverride` uses snake_case tags:

```go
NoiseScale      float64 `json:"noise_scale"`
NoiseW          float64 `json:"noise_w"`
SentenceSilence float64 `json:"sentence_silence"`
```

`SynthesizeSpeechWithVoice` unmarshals the renderer JSON into that struct. The camelCase fields do not match those tags, so they do not populate. The pipeline also does not consume them.

**Source-derived.** A typed voice contract is required. A map plus two different naming schemes is not safe for persisted character data.

### 4.4 Character identity

**Observed.** Character IDs are generated in `collectCharacterData` as:

```text
char_<Date.now()>_<random>
```

`migrateCharactersWithIds` adds IDs to legacy characters. Character IDs are used by sidebar panels, context selection, and the chat message DOM (`data-character-id`).

**Observed.** `ChatHistory` in `internal/backend/chat/types.go:42-49` contains only:

- `ID`
- `Sender`
- `Text`
- `Role`
- `Timestamp`
- `HasImage`

It has no `CharacterID`.

**Observed.** `AppendChatMessage(role, content, sender)` and `ReplaceChatHistory` both construct `chat.ChatHistory` without a character ID.

### 4.5 Chat snapshot persistence

**Observed.** `PersistenceManager.SaveChat` writes `characters.json` from `[]*CharacterContext` (`internal/backend/chat/persistence.go:176-215`). `CharacterContext` embeds `entities.Character`, so it has no `howl` field. A frontend `character.howl.voice` is dropped when the chat is saved through the Go model.

**Observed.** `frontend/dist/js/chat-state-bridge.js:685-698` enriches a backend character with a matching localStorage record:

```js
return Object.assign({}, match, character, {
  image: character.image || match.image || match.avatar || null
});
```

This can recover `howl.voice` from localStorage because the backend object usually has no `howl`, but only when the same character still exists in that browser profile. It is not portable character persistence and it is not a chat-file contract.

**Source-derived.** A character's voice can appear to survive a restart on the same machine because localStorage enrichment restores it, not because the chat preserved it. Moving the chat to another profile or machine, or deleting/replacing localStorage, loses it.

### 4.6 Edit propagation and cache invalidation

**Observed.** `audio-service.js` caches merged voice configs by character ID in `_voiceConfigCache`. `_invalidateVoiceCache` exists but has no call site.

**Source-derived.** After a Character tab edit, the open chat can continue using the old merged voice config until a page reload. The localStorage record is updated, but the current `currentChatData` character object and the voice-config cache are not invalidated.

**Source-derived.** This is a correctness issue for the proposed bubble menu and for any inline character-voice edit. The update path needs one authoritative character-voice mutation that updates persistence, the active chat object, and the renderer cache.

## 5. Chat bubble and message architecture

### 5.1 Current bubble structure

**Observed.** Assistant messages are rendered in `frontend/dist/js/chat-render.js:62-115` and in the streaming templates in `frontend/dist/js/chat-stream.js`.

Each assistant bubble currently has a `.chat-bubble-buttons` group with:

- regenerate
- previous
- next
- continue
- edit
- copy
- undo
- delete

There is no three-vertical-dots button and no voice/TTS control in the bubble. `chat-constants.js` has no overflow/dots icon.

### 5.2 Live multi-character ownership

**Observed.** `streamSingleCharacterResponse` in `chat-stream.js:308-382` sets:

```js
messageDiv.dataset.characterId = character.id;
```

On completion it calls:

```js
window.AudioService.playTTS(streamContent, character.id);
```

This is the strongest existing association: a live streamed multi-character response knows its character.

### 5.3 Generic and image streams

**Observed.** The generic backend stream and the image stream do not set a character ID on the message. Their completion paths call:

```js
window.AudioService.playTTS(mainContent, null);
```

The `CHARACTER_INFO:` chunk updates the sender display, but not `messageDiv.dataset.characterId`. As a result, those responses cannot select a per-character voice during autoplay.

### 5.4 Historical reconstruction

**Observed.** `chat-state-bridge.collectForSave` writes history entries with only id, sender, text, role, and timestamp. It does not include `characterId`.

**Observed.** `ChatHistory` has no character-ID field, so even the `characterId` collected by `data-sync.js:230-250` is dropped by the backend model and append/rewrite paths.

**Observed.** `rebuildChatFromHistory` in `sidebar-session.js:300-330` reads `msg.characterId || msg.CharacterId`, then calls `addAssistantMessage(sender, text, characterId)`. For backend history that value is always absent.

**Observed.** The alternate `rebuildChatFromLoadedHistory` path in `chat-session.js` also does not attach a character ID to restored assistant messages. It is not a fallback path that restores ownership.

**Observed.** `addAssistantMessage` then does this:

```js
character = characterId
  ? cd.characters.find(c => c.id === characterId)
  : cd.characters[0] || null;
```

and later:

```js
messageDiv.dataset.characterId = characterId || (character ? character.id : null);
```

Therefore every unowned assistant message in a multi-character chat is displayed and tagged as the first character after reload.

This is a direct stop condition for the proposed bubble menu. A bubble cannot reliably know its owning character for historical messages, and assigning the first character is silent voice misattribution.

### 5.5 Existing context-menu pattern

**Observed.** `sidebar-characters.js:122-166` has a fixed-position `character-context-menu` pattern with outside-click dismissal and a menu of Edit Character, Replace Image, View Memories, and Remove from Chat. It is a sidebar character menu, not a chat-bubble menu.

**Source-derived.** Its visual/interaction pattern can inform a bubble menu, but the data source must be the bubble's owning character ID, not the sidebar's current character.

### 5.6 Proposed bubble control placement

**Source-derived.** A three-dots button can be added to `.chat-bubble-buttons` without redesigning the bubble, but it must not be implemented until:

- every assistant message has a persisted, non-inferred character ID;
- unknown-owner messages remain explicitly unknown instead of falling back to the first character;
- the character voice update path is authoritative and invalidates the renderer cache.

A safe menu would:

- read `data-character-id` from the owning `.chat-message`;
- resolve that ID against the current chat's character state;
- show the assigned voice and engine;
- open the existing character sheet at that character, or expose character-scoped quick controls;
- never write a global `ttsVoice` or global engine;
- show "Voice unavailable" when the character or voice asset is missing, not a different character's voice.

No placeholder menu should be added before those conditions are met.

## 6. Per-engine voice identity and selection contracts

### 6.1 Summary table

| Engine | Voice identity | Discovery | Asset form | Voice switch mechanism | Multiple voices on one active engine | Current adapter behavior |
|---|---|---|---|---|---|---|
| Piper | Model name, e.g. `en_US-lessac-medium` | Static profile plus model registry for `.onnx` files | Separate `.onnx` model file | Replace Piper synthesizer model path | Possible sequentially, not resident concurrently | `GetVoices` returns profile list; `runPiper` ignores `req.Voice`; `ensureTTSVoice` replaces the synthesizer |
| Kokoro | Voice archive key, e.g. `af_heart` | Python bridge `--list-voices` | One voices archive plus model | Per-request `--voice` argument | Yes, same model/archive can serve different voices | `GetVoices` lists archive voices; `runKokoro` passes `req.Voice` |
| MeloTTS | Speaker ID, e.g. `EN`, `ZH`, `JP` | Static profile list | Remote server voices | Per-request `speaker_id` | Yes, remote server handles speaker selection | `GetVoices` returns profile list; HTTP payload uses `normalizeSpeakerID(req.Voice)` |

### 6.2 Piper

**Observed.** `internal/audio/tts/piper/local.go:157-175` returns `profile.SupportedVoices` from the model registry. It does not scan installed files.

**Observed.** `runPiper` at `piper/local.go:189-200` builds the command with `--model s.modelPath`. It does not use `req.Voice` at all. The voice is a property of the loaded synthesizer instance, not the synthesis request.

**Observed.** `ensureTTSVoice` in `internal/audio/pipeline.go:485-522`:

- compares `p.loadedVoice`;
- resolves `models/<voiceName>.onnx` or a matching `HOWL_TTS_PIPER_MODEL`;
- creates another `piper` synthesizer through the loader;
- replaces `p.ttsSynth` and `p.loadedVoice`.

It does not check the selected `TTSEngine`. It always loads Piper.

**Source-derived.** If a Piper voice file is missing, `ensureTTSVoice` returns an error. `ProcessAudioMessage` logs `voice switch failed ... using current voice` and continues. `synthesizeWithCache` still receives `override.Voice` and includes it in the cache key, while `runPiper` actually uses the previous model path. This is silent substitution plus a cache key that claims the requested voice.

**Source-derived.** If the active engine is Kokoro or Melo and a character requests a voice name that happens to exist as a Piper model, `ensureTTSVoice` replaces the active Kokoro/Melo synthesizer with Piper. Subsequent characters in the same generation can then be spoken by Piper. This violates the requirement that changing one character's voice must never silently change another character's voice.

### 6.3 Kokoro

**Observed.** `internal/audio/tts/kokoro/local.go:203-221` lists voices by invoking the Python bridge with `--list-voices` and returns each output line as an ID. It does not validate that a requested voice is in the list before synthesis.

**Observed.** `runKokoro` at `kokoro/local.go:120-137` passes `--voice <req.Voice>` and `--speed <req.Speed>` on every synthesis. The model and voices archive remain the initialized Kokoro assets.

**Source-derived.** Multiple characters can use different Kokoro voices with one initialized synthesizer because voice selection is a request parameter. A missing or invalid Kokoro voice is expected to fail in the Python bridge rather than fall back silently; that failure behavior is not live-tested.

**Observed.** `GetVoices` serializes with synthesis through `LockCalls`, so voice enumeration and synthesis contend on the same synthesizer lock.

### 6.4 MeloTTS

**Observed.** `internal/audio/tts/melotts/client.go:181-195` returns `profile.SupportedVoices`, currently `EN`, `ZH`, `JP`, with language hints.

**Observed.** `normalizeSpeakerID` at `melotts/client.go:340-350` passes any non-empty value through unchanged, and falls back to the first supported voice only when the value is empty. It does not validate against the server's actual speakers.

**Observed.** `normalizeMeloLanguage` at `melotts/client.go:352-378` maps non-empty unsupported language prefixes to `EN`.

**Observed.** `synthesizeWithCache` never sets `SynthesisRequest.Language`. The MeloTTS request therefore receives the zero value, and `normalizeMeloLanguage("")` returns `EN`.

**Source-derived.** A character can request a `ZH` or `JP` speaker while the request still carries the default `EN` language. The voice identity and language are separate dimensions, but the current model conflates them and the pipeline cannot set language at all.

**Observed.** The legacy `/synthesize` payload uses `speaker_id: normalizeSpeakerID(req.Voice)` and `speed`; the OpenAI-compatible payload uses `voice: normalizeSpeakerID(req.Voice)` and its existing speed clamp. Neither path discovers server voices.

### 6.5 Voice contract limitations

**Observed.** `types.Voice` contains only `ID`, `Name`, `Language`, `Gender`, and `Description`. It has no engine field, no availability state, no asset identity, and no capability record.

**Observed.** `GetTTSAvailableVoices` in `app-audio.go:350-360` returns only IDs from the current pipeline's synthesizer. It has no engine argument, no error, and no availability detail.

**Source-derived.** `GetTTSAvailableVoices` cannot serve an engine-aware Character tab because changing the engine changes the valid catalog. The Character tab currently uses a hard-coded list instead.

## 7. Full backend source-to-invocation trace

### 7.1 Current per-character path

The actual path for an autoplayed multi-character response is:

1. `chat-stream.js` completes a streamed character response and calls `AudioService.playTTS(text, character.id)`.
2. `audio-service.js:431-452` calls `_getMergedVoiceConfig(characterId)`.
3. The config is read from `currentChatData.characters.find(c => c.id === characterId).howl.voice`.
4. If `voice.enabled` is false, the function returns `null`; otherwise it merges explicit fields with global audio settings.
5. `playTTS` chooses the browser path when the merged engine is `browser`; otherwise it calls `SynthesizeSpeechWithVoice(text, JSON.stringify(voiceCfg))`.
6. `app-audio.go:185-205` unmarshals the JSON into `audiotypes.VoiceOverride`.
7. `AudioMessage.VoiceOverride` is attached only when `override.Voice != "" || override.Engine != ""`.
8. `processManagedAudio` obtains an operation and generation lease from the audio manager.
9. `Pipeline.ProcessAudioMessage` calls `ensureTTSVoice(ctx, msg.VoiceOverride.Voice)` if the voice is non-empty.
10. `ensureTTSVoice` always attempts a Piper load, regardless of `VoiceOverride.Engine`.
11. `ttsSpeedForConfig` resolves the global pipeline speed, then `synthesizeWithCache` sets the request voice from the override.
12. `GenerateCacheKey(text, voice, speed)` is computed and checked.
13. `synth.Synthesize` invokes the selected adapter.

### 7.2 What is lost between character and engine

**Observed.** The character ID is not passed into Go. `SynthesizeSpeechWithVoice` receives only text and voice JSON. If the voice JSON is invalid, the function logs a warning and proceeds with a zero-value override instead of returning an error.

**Observed.** `VoiceOverride.Engine` is parsed but `Pipeline.ProcessAudioMessage` never reads it. The active engine is fixed in `PipelineConfig.TTSEngine` from the global Apply.

**Observed.** `VoiceOverride.Volume` is parsed but not used by the pipeline. The renderer applies volume when playing the returned audio.

**Observed.** `VoiceOverride.NoiseScale`, `NoiseW`, and `SentenceSilence` are tagged with snake_case JSON names, while the renderer sends camelCase names. They are not parsed and are not consumed even when present.

**Observed.** `VoiceOverride` has no `Language` field. `SynthesisRequest` has a `Language` field, but `synthesizeWithCache` does not populate it.

**Observed.** `SynthesisRequest` contains only `Text`, `Voice`, `Speed`, and `Language`. There is no engine, noise, sentence-silence, volume, character ID, or voice revision field.

### 7.3 Global voice path

**Observed.** `DesiredSettings` contains `TTSVoice`, but `Backend()` excludes it. `Settings` and `PipelineConfig` have no voice field.

**Observed.** `internal/audio/settings_contract.go:736` reports `ttsVoice` as:

```text
This TTS option is not wired in this build.
```

**Observed.** `SynthesizeSpeech` (the path without a character override) sends no voice. `synthesizeWithCache` chooses `firstVoiceOrDefault(synth, "default")`, so the pipeline uses the synthesizer's first voice, not the persisted global `ttsVoice`.

**Source-derived.** The global `ttsVoice` selector is presently a renderer-side default/merge value, not an applied backend setting. A character voice system must not build on the assumption that global voice selection already works.

### 7.4 Silent substitution and cross-character contamination

**Source-derived.** Two concrete failure modes violate the product requirement:

1. **Missing Piper asset:** `ensureTTSVoice` fails, the pipeline logs and continues; `runPiper` uses the old model while the cache key uses the requested voice. The user hears a different voice than the requested one.
2. **Engine contamination:** a character's Piper-compatible voice causes `ensureTTSVoice` to replace the active Kokoro or Melo synthesizer with Piper for the whole pipeline generation. A later character can then be synthesized by the wrong engine.

Both must fail closed or be handled by an explicit engine router. Neither may report the requested voice as applied.

### 7.5 Character persistence path

**Observed.** `SaveChatData` calls `applyChatDataToState`, which parses characters through `ParseCharacterContextFromMap`, and then `convertFrontendDataToChatData`, which parses them again. The typed `CharacterContext` does not carry `howl`.

**Observed.** `SaveChat` writes `characters.json` from that typed data. The chat snapshot therefore loses `howl.voice`.

**Observed.** `LoadChat` returns typed characters. The renderer's `_enrichCharacterFromLocalStorage` can merge a matching localStorage record back in, but only on the same browser profile and only if the character still exists there.

**Observed.** History is persisted without a character ID. The `data-sync.js` DOM scraper does collect `dataset.characterId`, but the canonical `chatBridge.collectForSave` path used by the main save does not, and `ChatHistory` has no field for it.

## 8. Configuration ownership and recommended precedence

### 8.1 Current ownership

| Layer | Current behavior |
|---|---|
| Character sheet | Stores `howl.voice` with enabled, engine, voice, speed, volume; advanced fields are not exposed |
| Global audio settings | Stores `ttsEngine`, `ttsVoice`, `ttsSpeed`, volume, noise, sentence silence; backend applies engine and speed, reports voice/noise/sentence as unsupported |
| Renderer merge | `_getMergedVoiceConfig` merges character fields over global fields |
| Backend | Applies only global engine and global speed; per-request override applies voice and speed |
| Engine | Resolves whether a requested voice is valid only implicitly; Piper switches globally; Kokoro and Melo use request parameters |

### 8.2 Current precedence

For an enabled character voice:

```text
character.howl.voice.engine
  ?? global.ttsEngine
character.howl.voice.voice
  ?? global.ttsVoice
character.howl.voice.speed
  ?? global.ttsSpeed
character.howl.voice.volume
  ?? global.ttsVolume
```

The merged engine selects the browser path or the backend path. The merged voice and speed are sent as an override. The merged engine is not sent to Go in a way the pipeline honors.

For a character without an enabled voice, `playTTS` passes `null` as the config and calls `SynthesizeSpeech` with no voice override. The backend then uses its current synthesizer's first voice.

### 8.3 Recommended ownership

- **Character sheet:** owns explicit `enabled`, engine, voice ID, language, speed, volume, and engine-supported voice parameters.
- **Global audio settings:** owns defaults and general audio behavior, not a single flat voice ID for every engine.
- **Engine profile/catalog:** owns which voice IDs exist, which languages and parameters are valid, and availability.
- **Runtime resolver:** resolves the effective voice per request and returns an explicit applied/rejected result.
- **Chat history:** owns the character ID needed to identify the speaker. It should not be the authority for the character's current voice.

### 8.4 Recommended precedence

1. If the character has no explicit voice configuration, use the global default only when the active engine can resolve it.
2. If the character has explicit configuration, use its engine and voice only if that engine is available and the voice is valid.
3. If the character inherits a field, resolve that field from the global default for the same engine.
4. If the active engine is automatic/unresolved, do not claim a character voice is applied. Resolve the engine first or report `pending/unavailable`.
5. If a voice asset is missing or incompatible, fail closed with the requested voice and the reason. Do not substitute the first voice or the previous character's voice.
6. If a character engine differs from the active global engine, route the request to that engine only if the backend has an explicit engine router. Otherwise report unsupported; do not silently use the global engine.
7. Keep character ID and resolved voice separate. Character ID selects configuration; engine plus voice plus synthesis parameters select audio.

### 8.5 Required data-model change

**Source-derived.** A single `ttsVoice` string is inadequate. The canonical character voice needs at least:

```text
enabled
engine
voice_id
language
speed
volume
noise_scale
noise_w
sentence_silence
voice_revision / asset_identity
```

The exact schema should be a typed Go struct and one JSON contract, not a `map[string]interface{}` plus a differently tagged request struct.

## 9. Multi-character performance and resource feasibility

### 9.1 What the source establishes

**Observed.** Piper, Kokoro, and Melo each serialize synthesis with `BaseSynthesizer.LockCalls`. Concurrent requests against one synthesizer instance queue on that lock.

**Observed.** Piper and Kokoro invoke an external process for each synthesis. They do not keep an in-process model resident between calls. MeloTTS is an HTTP client; its model residency is on the remote server.

**Observed.** `streamMultiCharacterResponses` processes characters sequentially. When a character completes, `playTTS` stops any current audio before starting the next. Autoplay therefore does not overlap character speech.

**Observed.** `AudioService.playTTS` also stops the previous audio before starting a new request.

**Observed.** `ensureTTSVoice` calls `ModelLoader.LoadTTS` and replaces `p.ttsSynth`. `trackTTSLoaded` assigns a unique loader ID for every TTS load, and `ensureTTSVoice` never unloads the replaced synthesizer. The loader's loaded-model registry grows with voice changes. `Pipeline.Shutdown` later unloads all tracked IDs.

### 9.2 Declared profile estimates

These are source-derived declarations, not measurements:

| Engine | Declared RAM | Declared VRAM | Real-time factor | Backend |
|---|---:|---:|---:|---|
| Piper | 200 MB | 0 | 0.1 | local CLI |
| Kokoro | 600 MB | 0 | 0.3 | local Python/ONNX |
| MeloTTS | 500 MB | 1 GB | 0.5 | HTTP |

Source: `internal/audio/types/model_profiles.go:171-265`.

### 9.3 Feasibility conclusion

- **Same engine, multiple voices:** feasible with a bounded backend change for all three engines. Kokoro and Melo switch voices by request parameter. Piper switches model files and currently replaces the pipeline synthesizer; it needs a safe cache/switch policy and no cross-engine replacement.
- **Different engines per character:** not supported by the current architecture. There is one active `Pipeline.ttsSynth`, one `TTSEngine`, and no engine field in the effective request. Supporting this requires an engine router or multiple synthesizer instances, plus resource management for resident models.
- **Concurrent character speech:** currently serialized, not concurrent. Autoplay is sequential. Any future overlap needs a queue/priority policy and per-engine limits.
- **Resident voice assets:** Piper voices are separate model files; Kokoro uses one voices archive; Melo uses server speakers. No evidence establishes a need for multiple resident Piper models, but no benchmark has measured switch latency or peak memory.
- **Resource-management branch:** not required for same-engine, one-request-at-a-time voices if the adapters remain process-per-call or remote. It is required before per-character engines, concurrent model residency, or aggressive voice caching.
- **Verdict:** character voices can work safely before the resource-management branch only under a single active engine and a bounded voice-switch policy. The current `ensureTTSVoice` behavior is not that policy.

### 9.4 Minimum benchmark required

No latency or memory figure should be claimed until at least this matrix is measured on the target Windows machine:

1. Two characters alternating 10 short utterances on one engine.
2. Cold first synthesis versus warm subsequent synthesis.
3. Voice switch every message versus voice held constant.
4. Piper with three installed voices.
5. Kokoro with at least three voices from the archive.
6. Melo with at least three server speakers, if a live endpoint is available.
7. Peak RSS and VRAM per engine; process count per synthesis; loader tracked-ID count.
8. Two and four simultaneous synthesis requests, including cancellation.
9. Cache hit rate and correctness for the same text with different voices and engines.

Report the measurement environment, model files, voice IDs, text lengths, and whether the service was local or remote. Do not generalize from registry estimates.

**Unestablished:** the following cannot be decided from source alone and require the benchmark or a live service:

- whether Piper model startup or Kokoro process startup dominates per-message latency;
- whether one Kokoro archive can serve concurrent voices without memory or bridge contention;
- whether a MeloTTS server accepts every configured speaker ID and language combination;
- whether a per-character engine router can remain within the release memory/VRAM budget;
- whether cached audio remains valid across engine or voice-asset upgrades.

## 10. Cache correctness

### 10.1 Current key

**Observed.** `types.GenerateCacheKey` hashes only:

```text
voice | text | speed
```

Source: `internal/audio/types/utils.go:266-270`.

### 10.2 Correctness risks

- **Engine is absent.** If two engines can produce the same voice ID, or if the pipeline changes engine while the cache directory is shared, a cached file can be reused for the wrong engine.
- **Language is absent.** Melo voice/language combinations can collide.
- **Advanced synthesis parameters are absent.** Noise scale, noise width, and sentence silence would affect audio but are not in the key.
- **Piper request voice can disagree with the actual model.** On a failed switch, the key says the requested voice while the audio uses the old model.
- **Cache invalidation is absent.** Changing a character voice creates a new key but does not remove old entries; that is acceptable for correctness if the key is complete, but not for storage control.
- **Character identity is absent.** That is correct for sharing identical synthesis output, but only if character-specific style is fully encoded in the synthesis parameters.

### 10.3 Recommended key

Use a canonical synthesis contract:

```text
engine | voice_id | language | speed | noise_scale | noise_w |
sentence_silence | text | request_revision
```

Include character ID only if the character contributes synthesis behavior not already represented, or as diagnostic metadata outside the hash. Do not use character ID as a substitute for voice/engine identity.

Changing a character's voice should produce a new key. It should not invalidate another character's entries unless the underlying voice asset changed. Asset changes should include an asset identity or revision in the key.

## 11. Proposed backend architecture

This is a proposal for a separately authorized backend stage, provisionally Group 2C. Nothing below was implemented.

### 11.1 Canonical character voice

Create one typed character voice object on the character record. It should be serialized with the character, not reconstructed from a map:

```text
enabled: bool
engine: "piper" | "kokoro" | "melotts" | "browser" | ""
voice_id: string
language: string
speed: float
volume: float
noise_scale: float       // Piper only
noise_w: float           // Piper only
sentence_silence: float  // Piper only
asset_revision: string
```

The existing `howl.voice` data should be migrated into this contract. The unused `VoiceConfig` map should not remain a parallel authority.

### 11.2 Voice catalog

Add an engine-aware catalog:

```text
ListVoices(ctx, engine) -> []VoiceDescriptor
ValidateVoice(ctx, engine, voiceID) -> result
```

`VoiceDescriptor` should include engine, ID, display name, language, availability, asset identity, and supported parameters. Implementation belongs in the shared audio/capability layer:

- Piper: scan registered `.onnx` files and report installed/missing.
- Kokoro: query the voices archive through the bridge.
- Melo: query the server when an endpoint exists; otherwise report configured speakers and their unverified state.
- Browser: enumerate renderer voices only in the renderer, not in the Go catalog.

### 11.3 Resolver

Add a runtime resolver:

```text
ResolveCharacterVoice(characterVoice, globalDefaults, activeEngine) -> ResolvedVoice | explicit error
```

The resolver must never silently choose another voice. It must return the requested identity and the reason it cannot be applied. Automatic engine selection must be resolved before a voice can be reported applied.

### 11.4 Synthesis request

Extend the request contract:

```text
Engine
VoiceID
Language
Speed
NoiseScale
NoiseW
SentenceSilence
CharacterID       // observability and routing, not a global setting
```

The pipeline must set every field from the resolved voice. Adapters must reject unsupported fields rather than ignore them.

### 11.5 Engine router and lifecycle

If per-character engine is a product requirement, add an engine-aware router:

- one synthesizer instance per engine, or a bounded pool;
- a selected engine per request;
- explicit availability and failure results;
- no cross-engine replacement;
- ownership and release for every loaded synthesizer;
- a policy for resident voices and eviction.

Until that exists, the product must either restrict character voice to the active engine or report per-character engine changes as unsupported.

### 11.6 Character and chat persistence

- Persist the typed voice object with the character record.
- Persist the character ID on every assistant history entry.
- Keep the chat's character reference by ID; do not use the chat snapshot as the source of truth for voice.
- Decide whether historical playback uses the character's current voice or a per-message voice snapshot. If snapshots are required, add a voice revision to the message.

### 11.7 Cache

Use the expanded synthesis key in section 10.3. Keep the cache shared across characters when the resolved contract is identical. Do not key by character name or sender label.

## 12. Proposed Character tab integration

**Source-derived.** Use the existing voice block in `character.html:226-266`. Do not create a new global voice selector.

The block should:

- keep the enable checkbox;
- offer engine choices that are actually available;
- populate voice options from `ListVoices(engine)`, not the five hard-coded Piper names;
- show availability and missing-asset status;
- expose only parameters the selected engine supports:
  - Piper: voice, speed, noise scale, noise width, sentence silence;
  - Kokoro: voice, speed;
  - Melo: speaker, language, speed;
  - browser: browser voice, rate, volume;
- show "Inherit global" explicitly for each inherited field;
- validate and report incompatibility in place;
- preserve all unrelated `howl` fields on save;
- preserve unknown voice fields through an edit rather than reconstructing the object;
- invalidate the active renderer voice cache and chat character object after save.

The existing deep link must be fixed so `character.html?id=...` opens the named character. The bubble menu depends on that route if it uses "Open voice settings."

No controls should be added as placeholders before the backend catalog and resolver exist.

## 13. Proposed per-character chat bubble controls

**Source-derived.** A three-dots button can be the entry point, not the storage location.

Proposed behavior:

- Add a `.bubble-button` to assistant bubbles with a vertical-dots icon.
- Read the owning `.chat-message[data-character-id]`.
- If the ID is missing, show "Character unknown for this message" and disable voice actions. Never default to the first character.
- Resolve the character from the current chat, then show:
  - current engine and voice;
  - voice availability;
  - a link to the character sheet at that character;
  - optional enable/disable speech for that character, if the existing TTS architecture supports it;
  - optional character-scoped quick controls, using the same mutation path as the character sheet.
- Make every mutation write to the character voice object, not to global `ttsVoice`, `ttsEngine`, or global speed.
- Invalidate the renderer cache and update the open character object after mutation.
- Use the existing menu dismissal/positioning pattern only as a visual precedent; do not reuse sidebar character selection as the data source.

**Prerequisite blocker:** historical assistant messages have no persisted character ID. Until that is fixed and the first-character fallback is removed, the menu cannot safely operate on older conversations.

## 14. Shared infrastructure boundaries

| Responsibility | Correct owner |
|---|---|
| Character voice fields, defaults, enable state, per-character overrides | Character configuration and persistence |
| Voice discovery, installed-asset checks, engine capabilities | Shared audio capability/catalog layer |
| Engine-specific voice semantics, validation, request payloads | TTS engine adapters |
| Effective voice resolution, request construction, cache, scheduling | Audio pipeline |
| Resident model/voice lifetime, VRAM/memory arbitration, eviction | Resource-management branch |
| Chat message ownership and character-ID persistence | Chat persistence and message model |
| Character tab controls and bubble menu | Renderer UI, consuming the shared contracts |

The Character tab must not own model scanning. The chat bubble must not own voice validation. The audio pipeline must not infer character identity from sender names.

## 15. Exact implementation sequence and prerequisites

### Phase 0 — product decisions

Resolve the unanswered questions in section 16 before schema work. The engine-scope and fallback decisions change the backend design.

### Phase 2C.1 — character voice schema and persistence

- Define the typed character voice contract.
- Migrate `howl.voice`.
- Preserve unknown/unrelated character fields on save.
- Persist character ID on assistant history.
- Make chat snapshots reference character identity rather than duplicate voice state.
- Fix `character.html?id=...` if the bubble menu will deep-link.

### Phase 2C.2 — engine-aware voice catalog and resolver

- Add `ListVoices(engine)` and `ValidateVoice`.
- Add installed-asset/availability reporting.
- Add the resolver and fail-closed error contract.
- Add language and advanced parameter support.
- Remove the first-character and previous-voice fallback paths.

### Phase 2C.3 — adapters, pipeline, and cache

- Make `ensureTTSVoice` engine-aware or replace it with an engine router.
- Stop loading Piper for Kokoro/Melo requests.
- Populate the full synthesis request.
- Expand the cache key.
- Release or bound every lazy-loaded synthesizer.
- Add backend tests for missing assets, wrong-engine requests, per-character isolation, and cache isolation.

### Phase 2D — Character tab

- Replace the static select with backend catalogs.
- Add engine-specific controls and availability states.
- Preserve all character fields.
- Invalidate the active chat voice cache after save.

### Phase 2E — bubble menu

- Add the three-dots entry point.
- Require a persisted character ID.
- Show unknown-owner state rather than first-character fallback.
- Route edits through the same character voice mutation as the Character tab.

### Phase 2F — benchmark and release gate

Run the benchmark in section 9.4 before claiming acceptable multi-character performance.

### 15.1 Protected behavior

Group 2C must not regress the accepted work in these areas:

- **Group 0:** desired-settings contract, accurate ownership/status reporting, persistence, and no-auto-load behavior.
- **Group 1:** Whisper decode settings.
- **Group 2A:** TTS speed, including Piper inverse length scale, Kokoro direct speed, MeloTTS payload correction, validation, and cache correctness.
- **Audio lifecycle:** WASAPI microphone gate and lifecycle behavior.
- **Separation:** TTS independence from ASR.
- **Playback:** headphone-only audio playback.
- **Conversation state:** conversation persistence and restoration.

Character voice must use the existing Group 2A speed resolver rather than introducing a second speed path. Character-specific speed, if approved, must remain an explicit override with a reported applied/unavailable status, not a silent replacement of the global value.

No protected behavior was changed during this investigation.

## 16. Unresolved product decisions

1. Is a character allowed to select its own engine, or only a voice within the globally active engine? The current UI implies independent engine selection; the backend cannot honor it.
2. Does a new character default to "inherit global voice" or "character voice disabled"?
3. When a voice is unavailable, should synthesis fail with a clear error, use a designated global default, or use a designated per-character fallback? Silent substitution is prohibited either way.
4. Does global `ttsVoice` remain a default, become engine-scoped, or disappear from the product?
5. Is language a character property, a voice property, or both?
6. Should historical playback use the character's current voice or the voice revision that existed when the message was generated?
7. Should the bubble menu deep-link to the character sheet or offer inline quick controls?
8. Is character volume a synthesis parameter or a playback-only control?
9. Which browser TTS voices are supported, and how are they identified across platforms?
10. How should existing `howl.voice` values with unsupported voice IDs be migrated?
11. Can per-character speed override the global Group 2A speed? If so, how should status report both?
12. What happens to authored character voice data when a character is removed from a chat, moved to another machine, or loaded without its local character record?

## 17. Files and source areas inspected

### Character UI and persistence

- `frontend/dist/character.html`
- `frontend/dist/js/character.js`
- `frontend/dist/js/character-form.js`
- `frontend/dist/js/character-grid.js`
- `frontend/dist/js/character-state.js`
- `frontend/dist/js/character-storage.js`
- `frontend/dist/js/character-manager.js`
- `frontend/dist/js/sidebar-characters.js`
- `frontend/dist/js/sidebar-panels.js`
- `frontend/dist/js/sidebar-session.js`
- `frontend/dist/js/chat-state-bridge.js`
- `frontend/dist/js/data-sync.js`

### Chat bubbles and association

- `frontend/dist/js/chat-render.js`
- `frontend/dist/js/chat-stream.js`
- `frontend/dist/js/chat-actions.js`
- `frontend/dist/js/chat-constants.js`
- `frontend/dist/js/chat-session.js`
- `frontend/dist/css/cc-messages.css`
- `frontend/dist/css/cc-buttons.css`
- `frontend/dist/js/chat-input.js`

### Renderer audio

- `frontend/dist/js/audio-service.js`
- `frontend/dist/js/audio-settings.js`
- `frontend/dist/audio-settings.html`
- `frontend/dist/js/settings-state.js`

### Backend character and chat persistence

- `app-characters.go`
- `app-chat.go`
- `app-prompt.go`
- `internal/backend/chat/types.go`
- `internal/backend/chat/parser.go`
- `internal/backend/chat/persistence.go`
- `internal/backend/types/entities/entities.go`
- `internal/backend/types/context.go`

### Backend audio

- `app-audio.go`
- `app-model-registry.go`
- `internal/audio/manager.go`
- `internal/audio/pipeline.go`
- `internal/audio/runtime/loader.go`
- `internal/audio/settings_contract.go`
- `internal/audio/types/interfaces.go`
- `internal/audio/types/model_profiles.go`
- `internal/audio/types/utils.go`
- `internal/audio/tts/tts.go`
- `internal/audio/tts/piper/local.go`
- `internal/audio/tts/kokoro/local.go`
- `internal/audio/tts/melotts/client.go`

### Existing review artifacts

- `docs/review/AUDIT-audio-settings-contract-2026-11-05.md`
- `docs/review/BLOCKER-audio-tts-speed-group2a-2026-11-05.md`
- `docs/review/IMPLEMENTATION-audio-tts-speed-group2a-2026-11-05.md`

## 18. Conclusion and stop conditions

The same-engine, per-character voice concept is feasible, but the current implementation is not safe enough to build UI on. It is a renderer-side partial path with silent Piper fallback, an ignored character engine, an unvalidated voice catalog, an incomplete cache key, and no persisted message ownership.

Implementation should not proceed until a backend-first Group 2C resolves:

- typed character voice ownership;
- character ID persistence in chat history;
- engine-aware catalog and resolver;
- engine-correct synthesis requests and cache keys;
- Character tab save integrity and cache invalidation.

No Group 2C work, Character tab changes, bubble controls, tests, builds, synthesis, model loads, settings changes, commits, or pushes were performed for this investigation.

## 19. Review disposition

**Accepted:** 2026-11-06

**Group 2B status:** Closed as an accepted investigation.

**Implementation status:** Blocked. No Group 2C, Character tab, or chat-bubble implementation is authorized.

The independent review accepted the investigation and agreed that three findings are architectural blockers:

- historical messages lack reliable character ownership;
- character-sheet saves can discard unrelated fields;
- engine selection is not genuinely character-specific and can affect other characters.

It also classified the cache omissions as high severity and the absence of performance measurements as open.

The accepted implementation order is:

1. **Character identity and persistence foundation.** Establish reliable character ownership for messages, including historical conversations and restoration.
2. **Character configuration integrity.** Make character-sheet updates preserve every existing field and persist voice configuration safely.
3. **Backend character voice architecture.** Resolve character-specific engine and voice settings per synthesis request without unintended global engine replacement.
4. **Cache correctness and performance.** Separate incompatible synthesis outputs and benchmark multi-character switching, concurrency, latency, and memory.
5. **Character tab UI.** Expose voice configuration through the existing character editor.
6. **Chat bubble UI.** Add the per-character three-dot menu connected to the owning character sheet.

These stages are proposed and accepted for planning only. They are not additional authorizations. The next gate is explicit authorization of stage 1.

The following architectural rule is confirmed:

- **Character:** owns voice preferences.
- **Audio settings:** provide defaults and general playback behavior.
- **Engine:** determines compatible voice capabilities.
- **Resource manager:** determines what can be loaded and when.
- **Synthesis request:** resolves the effective configuration for that particular character.

This separation is required so that one character's voice selection cannot change the next character's output. Where this disposition refines the phase sequence in section 15, this section governs.
