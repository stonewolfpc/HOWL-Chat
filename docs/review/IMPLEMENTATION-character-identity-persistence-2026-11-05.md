# Implementation — Character Identity and Persistence Foundation

**Date:** 2026-11-05  
**Status:** Implemented and locally verified; not committed or pushed.  
**Scope:** The identity/persistence foundation only. No character voice routing, voice resolution, synthesis, cache redesign, Character tab UI, chat-bubble controls, character-sheet save-integrity work, performance benchmarking, or later-stage work was implemented.

## 1. Result

Assistant messages now carry an explicit character ID from the character response path through in-memory history, snapshot persistence, NDJSON history, restart restoration, and frontend rendering.

The identity is never derived from array position, the first character, the current character selection, or the renderer's label. Messages without a reliable persisted identity remain unassigned.

No migration fabricates IDs for historical messages. Existing records without a `characterId` field load with an empty ID, which is the explicit unknown-ownership value.

## 2. Files changed

### Backend and persistence

| File | Change |
|---|---|
| `internal/backend/chat/types.go` | Added `ChatHistory.CharacterID string` with `json:"characterId,omitempty"`. |
| `app-prompt.go` | Added internal `appendChatMessage(role, content, sender, characterID)`; kept public `AppendChatMessage(role, content, sender)` unchanged for user/system turns; character responses now pass the generating character ID; `ReplaceChatHistory` now preserves owner, timestamp, and `hasImage`; system timestamps use RFC3339. |
| `app-chat.go` | `buildChatHistory` preserves character ID, timestamp, and `hasImage`; `updateAppStateFromChat` restores those fields; `convertFrontendDataToChatData` parses character ID; non-character message creation uses RFC3339 timestamps so they survive later history rebuilds. |
| `internal/backend/chat/history_character_test.go` | New persistence tests for JSON round-trip, append/load/rewrite, legacy unknown ownership, repeatable rewrite, and `SaveChatAs`. |
| `app_chat_character_test.go` | New app-level tests for generation append, save/restart restoration, final-save rewrite, replacement history, legacy app save, and frontend conversion. |

### Frontend rendering and save paths

| File | Change |
|---|---|
| `frontend/dist/js/chat-render.js` | `addAssistantMessage` records `characterId` explicitly (including `''`); removed the `cd.characters[0]` fallback; only an explicit ID selects a character for display. |
| `frontend/dist/js/sidebar-session.js` | History reconstruction restores `characterId`, timestamp, and `hasImage` on the rendered element; unknown ownership is stored as an empty ID rather than omitted or replaced. |
| `frontend/dist/js/chat-session.js` | Alternate rebuild/collect paths preserve character ID, timestamp, and `hasImage`. |
| `frontend/dist/js/chat-state-bridge.js` | Canonical `collectForSave` payload includes character ID, persisted timestamp, and attachment flag. |
| `frontend/dist/js/chat-actions.js` | Replacement-history entries include character ID, timestamp, and attachment flag. |
| `frontend/wailsjs/go/models.ts` | Generated `chat.ChatHistory` model exposes `characterId`. |
| `tests/frontend/chat-render-ownership.test.js` | New dependency-free Node tests for explicit, missing, and unknown character identity rendering. |

### Documentation

| File | Change |
|---|---|
| `docs/review/IMPLEMENTATION-character-identity-persistence-2026-11-05.md` | This report. |

## 3. Data model

### Before

`internal/backend/chat/types.go` stored:

```go
type ChatHistory struct {
    ID        string `json:"id"`
    Sender    string `json:"sender"`
    Text      string `json:"text"`
    Role      string `json:"role"`
    Timestamp string `json:"timestamp"`
    HasImage  bool   `json:"hasImage"`
}
```

There was no persisted owner. `AppendChatMessage(role, content, sender)` had no character parameter, and the character response path therefore recorded only the display name.

### After

```go
type ChatHistory struct {
    ID        string `json:"id"`
    Sender    string `json:"sender"`
    Text      string `json:"text"`
    Role      string `json:"role"`
    Timestamp string `json:"timestamp"`
    HasImage  bool   `json:"hasImage"`
    // Empty means ownership is not established (legacy or non-character).
    CharacterID string `json:"characterId,omitempty"`
}
```

`omitempty` is intentional:

- Assigning an empty ID does not add a fabricated field to old or non-character records.
- Loading a record that never had the field produces the same empty value.
- There is one representation for unknown ownership: `CharacterID == ""`.

## 4. Ownership path

### Live character response

1. `GenerateCharacterResponse(message, characterID)` already received the generating character's stable ID.
2. It now calls `appendChatMessage(..., characterID)` instead of the public three-argument binding.
3. `appendChatMessage` stores the ID in the in-memory message and in the `chat.ChatHistory` entry appended to `history.ndjson`.
4. The public `AppendChatMessage(role, content, sender)` remains available for user/system turns and delegates with an empty character ID. This preserves its Wails signature and its non-character semantics.

The ID is established at the point where the character response is generated; it is not inferred later.

### Save and replacement

- `chat-state-bridge.collectForSave` reads `data-character-id` from each rendered message.
- `chat-actions.js` includes the same value in entries sent to `ReplaceChatHistory`.
- `ReplaceChatHistory` accepts `characterId` (and the generated `CharacterId` form) and writes it to both history and the in-memory message.
- `buildChatHistory` carries the ID into snapshot history.
- `convertFrontendDataToChatData` carries it through frontend-to-backend conversion.
- `updateAppStateFromChat` restores it when a chat is loaded.
- `SaveChatAs` copies the history through the normal persistence path; the focused test confirms the copied chat retains the ID.

### Rendering

- An explicit ID selects only the matching character.
- A missing or unknown ID renders the stored sender and assigns no other character.
- The previous first-character fallback was removed.
- The rendered element retains the empty/explicit `characterId` value so a later save does not silently reassign it.

## 5. Persistence and migration behavior

### Existing formats

The existing persistence model already supports additive JSON fields:

- `history.ndjson` is the append log and is preferred when loading.
- `history.json` is the snapshot form.
- `SaveChatAs` creates a new chat and copies history.
- Final save rewrites `history.ndjson` from the restored state.

Adding `characterId` does not change the file layout, chat directory layout, or load precedence.

### What happens to legacy records

A legacy entry such as:

```json
{"id":"0","sender":"Aria","text":"legacy line","role":"assistant","timestamp":"2026-01-01T00:00:00Z"}
```

loads with `CharacterID == ""`. It is not assigned to "Aria", to the first character in the chat, or to the currently selected character.

On save/rewrite, the entry remains unassigned. Because the field is `omitempty`, no `characterId` is added to the legacy record. This is deliberate: unknown ownership stays unknown and historical records are not destructively rewritten.

### Migration decision

No schema version bump or one-time migration was introduced.

- The change is additive and backward-compatible.
- Historical messages cannot be reliably attributed from the data that exists, so no IDs are invented.
- The behavior is deterministic and repeatable: loading and rewriting the same legacy record produces the same unassigned record.
- New messages get an explicit owner only when the generating character ID is available.

## 6. Verification

### Commands and actual results

Run from `C:\Users\StoneWolfPC\HOWL_CHAT`.

```text
gofmt -w app-chat.go app-prompt.go internal/backend/chat/types.go \
  internal/backend/chat/history_character_test.go app_chat_character_test.go
gofmt_exit=0
# gofmt -l produced no file names
```

```text
go test ./internal/backend/chat -count=1
ok  	howl-chat/internal/backend/chat	2.193s
chat_tests_exit=0
```

```text
go test . -run 'Test(CharacterOwnership|ReplaceChatHistory|LegacyHistory|ConvertFrontendHistory|TranscribedTextPersistsInChatHistory)' -count=1
ok  	howl-chat	1.306s
app_tests_exit=0
```

The selected app tests were also run verbosely and passed individually:

```text
--- PASS: TestTranscribedTextPersistsInChatHistory
--- PASS: TestCharacterOwnershipSurvivesAppendSaveRestart
--- PASS: TestReplaceChatHistoryPreservesCharacterOwner
--- PASS: TestLegacyHistoryRemainsUnassignedThroughAppSave
--- PASS: TestConvertFrontendHistoryPreservesCharacterOwner
PASS
ok  	howl-chat	1.327s
```

```text
node --test tests/frontend/chat-render-ownership.test.js
✔ explicit character identity selects the owning character
✔ missing character identity does not fall back to the first character
✔ unknown character identity is preserved and does not select another character
ℹ tests 3
ℹ pass 3
ℹ fail 0
node_test_exit=0
```

```text
node --check frontend/dist/js/chat-render.js
node --check frontend/dist/js/sidebar-session.js
node --check frontend/dist/js/chat-session.js
node --check frontend/dist/js/chat-state-bridge.js
node --check frontend/dist/js/chat-actions.js
node_check_exit=0
```

```text
git diff --check -- app-chat.go app-prompt.go internal/backend/chat/types.go \
  frontend/dist/js/chat-render.js frontend/dist/js/sidebar-session.js \
  frontend/dist/js/chat-session.js frontend/dist/js/chat-state-bridge.js \
  frontend/dist/js/chat-actions.js
scoped_git_diff_check_exit=0
```

A full-worktree `git diff --check` still exits `2` because of pre-existing trailing whitespace in unrelated modified files, including `IMAGE_PAGE_CHANGE_LEDGER.md` and generated `frontend/wailsjs/go/models.ts` content. Those whitespace lines were not introduced by this task. The changed-source scoped check is clean.

### What the tests establish

- **Creation:** `appendChatMessage(..., "char_aria")` stores the owner in memory and in the append log.
- **Multiple characters:** two assistant messages from different IDs remain distinct after save and restart.
- **Save/restart:** a new `App` instance calling `LoadChat` restores both IDs and populates `chatMessages` with the same owners.
- **Final-save rewrite:** `performFinalSave` rewrites `history.ndjson` without losing owners.
- **Replacement:** `ReplaceChatHistory` preserves IDs, the supplied timestamp, and `hasImage`; a subsequent save/load does not lose them.
- **Frontend conversion:** frontend `characterId` is parsed; a user message without one stays empty.
- **Historical compatibility:** a legacy NDJSON record loads with an empty ID and survives two rewrite passes with content and timestamp intact.
- **Repeatability:** the legacy record gains no `characterId` field after repeated rewrites.
- **Save As:** the copied chat retains the explicit character ID.
- **Rendering:** explicit, missing, and unknown IDs do not fall back to the first character.

## 7. Historical compatibility findings

**Observed in source and tests.**

- Legacy records have no character ID.
- The zero value of the new field is the intended unknown state.
- JSON unmarshalling of old records succeeds without a migration.
- `omitempty` prevents an empty owner from being materialized into rewritten legacy text.
- Existing fields (`ID`, `Sender`, `Text`, `Role`, `Timestamp`, `HasImage`) remain unchanged.
- Existing chat IDs, message IDs, content, timestamps, and attachment flags are not intentionally changed by this work.
- The change does not resend messages, replay inference, or load models during restoration.

**Source-derived conclusion.**

The previous renderer fallback could visually attribute an unowned message to the first character. Removing that fallback is necessary for unknown ownership to remain visually and semantically unknown; such a message displays its stored sender instead.

**Unverified assumption.**

No historical production conversation corpus was available for inspection. Legacy behavior was verified with representative old-format NDJSON records and the existing persistence code, not with every possible historical file variant.

## 8. Protected functionality

No changes were made to:

- Group 0 settings contract.
- Group 1 Whisper decode settings.
- Group 2A TTS speed.
- The accepted WASAPI microphone gate.
- No-auto-load behavior.
- TTS independence from ASR.
- Headphone-only playback.
- Character voice routing or voice resolution.
- TTS synthesis or cache behavior.
- Character tab UI or chat-bubble controls.

The only message-path changes are identity and metadata preservation, plus timestamp serialization in in-memory message maps to RFC3339 so timestamps can survive history rebuilds. The public `AppendChatMessage` Wails signature is unchanged.

## 9. Remaining limitations and decisions

1. **No historical ID recovery.** This is intentional. The available legacy data does not reliably identify the generating character, so old messages remain unknown.
2. **Character ID is only as stable as the ID already supplied to `GenerateCharacterResponse`.** The implementation does not create or modify character-sheet IDs.
3. **No full application restart or browser E2E was run.** The restart test constructs a fresh `App` with the same persistence manager and calls `LoadChat`/`performFinalSave`; that exercises the actual persistence and restoration code, but not a process launch.
4. **Frontend save collection is source-verified, not browser-integration-tested.** The renderer behavior is covered by Node tests, and the backend replacement/save path is covered by Go tests; a full browser flow was not run.
5. **The public three-argument `AppendChatMessage` intentionally records no character owner.** Future code that appends assistant messages through that binding must use the internal ID-aware helper or add an explicitly authorized contract.
6. **The legacy fallback path in `chat-actions.js` for a backend without `ReplaceChatHistory` cannot preserve the new field.** The current backend provides `ReplaceChatHistory`, so this fallback is not on the active path.

## 10. Evidence classification

- **Observed:** test execution results, source changes, JSON/persistence behavior, rendering assertions, and the existing repository cleanliness limitations above.
- **Source-derived:** the ownership path, the absence of first-character inference after the change, the additive compatibility argument, and the conclusion that no migration is needed.
- **Unverified:** production historical conversation variants, full Wails/browser restart behavior, and any UI visual appearance not exercised by the Node renderer tests.

## 11. Authorization boundary

This implementation stops at the identity and persistence foundation. It does not authorize or implement character voice routing, character-specific engine selection, voice resolution, synthesis, cache changes, Character tab UI, chat-bubble controls, character-sheet save-integrity changes, or performance benchmarking.

No commit or push was performed.
