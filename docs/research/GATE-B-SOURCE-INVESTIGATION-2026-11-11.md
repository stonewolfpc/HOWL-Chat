# Gate B Source Investigation — Howl Chat Working Tree

**Date:** 2026-11-11
**Gate:** B, source-level evidence only
**Status:** Read-only investigation completed; no implementation selected
**Baseline note:** All behavior below is working-tree behavior unless stated otherwise. HEAD is not treated as behavior.

The final “Report handoff” numbers in the closing assistant message carry the report file’s SHA-256, byte size, and line count. No digest is embedded here because embedding it would change the digest.

## 1. Executive findings and scope

### Scope completed

The following authorized Gate B questions were traced end to end at source level:

- Q-07: whether the secondary `chat.Service` memory/RAG prompt path is live.
- Q-01: the complete frontend character-composition contract.
- M-01: conversation-scope handling for memory.
- M-02: lifecycle/inactive-character material.
- Q-09: the existing composition contract and change surface.
- Q-02: world/scenario authority and retention behavior.

No Gate C owner decision and no Gate D runtime experiment were made or executed.

### Main findings

- F-01: The secondary `chat.Service` prompt path is constructed but not demonstrably production-reachable. `chat.NewService` is called with empty identifiers in `app.go:170`, RAG is attached in `app.go:186-187`, but no production Go path calls `SendMessage` or `SendMessageStream`. If the process never invokes those methods, memory, RAG, world/scenario, and prompt generation in that service cannot reach the model.
- F-02: The dormant secondary path and the live primary path use different model boundaries. The secondary path builds one legacy string prompt through `internal/backend/chat/service.go:612-690`, `694-780`, and `899-1026`. The live generic path builds a structured message array through `app-prompt.go:671-799`; the live character path builds one composed user string through `frontend/dist/js/chat-stream.js:121-176` and `app-prompt.go:154`.
- F-03: Character memories are character-scoped but not conversation-scoped. They are stored globally in localStorage, filtered only by stable character identity, composed into plain text, and sent without a chat identifier. The same stable character can therefore expose the same memories in another chat.
- F-04: There is no discovered enforcement layer that supplies conversation ownership between localStorage retrieval and model invocation. Backend chat-scoped memory exists for storage/editing APIs, but the primary character composer does not read from it.
- F-05: Recent lifecycle material is current-turn material from the rendered DOM and the current chat roster. Death/step-out facts additionally persist globally as character memories, so they can cross a chat boundary if the same character is used elsewhere.
- F-06: The lorebook routing layer does not compare entry origin against the active world or scenario. World and scenario entries are appended unconditionally, and `worldName`/`scenarioName` are unused in `RouteResolvedEntries`.
- F-07: A clear-world or clear-scenario action only clears the in-memory pointer before saving. The ordinary save path preserves a previously persisted `ChatData.World` or `ChatData.Scenario`, so a stale object can return after reload.
- F-08: Chat-scoped character-swap memories are written but have no prompt-facing retrieval path. They are therefore stored but not injected by the inspected composer.
- F-09: The normal chat-load path does not hydrate the backend lore pool. After `LoadChat`, world, scenario, characters, and transcript are set, but `a.lorebooks` is not.
- F-10: `RouteResolvedEntries` checks stable ownership only for character entries. It does not check ownership for world/scenario entries, and private world/scenario entries are not excluded by routing.

These are source-derived exposure and contract findings. They are not runtime observations that a particular leakage event occurred.

## 2. Verified input artifact identities

All frozen artifacts were hashed before reliance. Every digest matched the Gate B mission.

| Artifact | Expected SHA-256 | Result |
|---|---|---|
| Original R-07 evidence report | `e3860e3d5d6ca71724d31d0021897819c968afc1b8e1aad09fc011f8d9a67134` | matched |
| Independent Muse review | `4bdb399cfcb56c75e32d3103a97fd57b06db8895ea3fd8dfef16feba5a98d792` | matched |
| R-07 erratum and supplement | `930a4ea34b02d4d497417bf565c98ebb0752d3e540f9749dbdfe60ac` | matched |
| R-07 evidence register | `1f297d63903f0f386c97943e81e77ac5db3c3dc4b335968471595687df406ac6` | matched |
| Q-06 provenance report | `8b0ba58356e776c01e661865b619d64eabf2216f65fc65f4386c6ee60c4639b4` | matched |
| Q-06 baseline patch | `832d08e8857599fbe3b21965ceb75ada2c8f37d235c346307119` | matched |
| Amendment 02 revision 9 | `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640` | matched |
| Amendment 02 closure | `f828fd0469cef0508e1e2c39482875ae90665bcc07d96821df293a9746863050` | matched |
| Amendment 02 verification record | `aa48d8a529ba0d6dd9e7d39136fe9e9b72d81c6e9ee60ab0cf1d3` | matched |

No frozen artifact was edited, renamed, replaced, or regenerated.

## 3. Working-tree baseline caveat

Q-06 was accepted as the baseline. The working tree was inspected rather than HEAD.

Observed baseline before creating this report:

- 258 dirty Git entries.
- 107 tracked-modified files.
- 91 modified tracked source-like files.
- HEAD: `9c7d2a17dd7f702d7db6d95d9316e7640e34bf6f`.
- The five review-critical source files remain unstaged tracked modifications.
- Their provenance remains timestamp-supported, not provenance-proven.

Working-tree hashes observed during inspection:

- `app-prompt.go`: `93c9e6b8759936194a06ee5a42c5940d4128`
- `frontend/dist/js/chat-stream.js`: `2106be99e088476d58274b610d034167bb3a0643716a064c377a`
- `frontend/dist/js/memory-store.js`: `78a8f6b8c76d435ccdf4b58614d9c69552e7bad19d8c9cc65043`
- `internal/backend/chat/persistence.go`: `ad9663a48278cff0ee35eb48d28799702cdbfa7754f71a3a676a87c`
- `internal/backend/lorebook/lorebook.go`: `1f6b5d5fa18c9aeaa3e6bbe542e3c7dc5574d65ab972e524b069e`

These working hashes are identifying measurements, not evidence that committed history contains the same behavior.

## 4. Q-07 secondary `chat.Service` reachability analysis

### 4.1 Construction and RAG attachment

The service is created in `app.go:170`:

```go
chatService, err := chat.NewService(modelManager, client, "", "", []string{}, "")
```

RAG is attached in `app.go:186-187`:

```go
chatService.SetRAGService(ragService)
chatService.SetRAGEnabled(true)
```

Content libraries are later attached in `app.go:249-250`.

**Classification:** Observed.

This establishes that a service object exists and that its RAG pointers/flags are set. It does not establish that the service’s prompt methods execute.

### 4.2 Production calls not found

A targeted source search found no production caller for either generation method.

- `internal/backend/chat/service.go:612` defines `SendMessage`.
- `internal/backend/chat/service.go:694` defines `SendMessageStream`.
- No production Go path invokes either method.
- `app-characters.go:345` defines `syncAppStateFromChatService`, but it has no caller.
- The only identified `LoadChatData` uses are construction with empty identifiers and the `native_test_seam` build-tagged seam.
- Production builds use the no-op seam in `app_native_test_seam_off.go`.
- The service is not exposed as a Wails binding; it is an internal App field.

The search covered `app.go`, all `app_*.go` files, and `internal/backend/chat`. It cannot exclude a hypothetical unexamined dynamic route, but no such route was visible and the service object itself is not registered for external invocation.

**Classification:**

- Service instantiated: Observed.
- RAG enabled on the object: Observed.
- No production generation call: Source-derived.
- Dormant/legacy characterization: Inferred.
- Runtime execution: Unestablished.

### 4.3 Prompt flow if the service were invoked

If `SendMessage` or `SendMessageStream` were invoked, the service would:

1. Append the user message to its internal context at `service.go:616-624` or `698-708`.
2. Combine `manualLore` and `memorySnapshot` into one candidate pool.
3. Call `lorebook.Resolve` with `ActiveWorld`, `ActiveScenario`, and the first active character.
4. Call `RouteResolvedEntries` using the stable character ID.
5. Call `buildPrompt` to emit world, scenario, active character, other participants, memories, optional RAG text, history, and task text.
6. Use the legacy string-prompt methods on the llama client.

Memory enters at `service.go:631-632` and RAG enters in `service.go:985-994`. The RAG query uses the last context message plus `WorldName` and `ScenarioName` at `service.go:1028-1052`.

**Classification:** Source-derived conditional.

This is what the code would do. It is not evidence that the code does it in production.

### 4.4 Context condition at construction

`NewService` receives empty world, scenario, character, and chat identifiers. `LoadChatData` is not called again on a discovered production path. `SetContentLibraries` supplies registry pointers, but the memory snapshot and canonical world/scenario fields are not refreshed without another `LoadChatData` call.

**Classification:** Source-derived.

This is important: even if a generation method were invoked, the constructed empty state would supply less context than the method’s own logic suggests unless another caller supplied identifiers.

### 4.5 Q-07 answer in evidence terms

- The secondary path exists.
- It has memory support.
- It has RAG support.
- It is not connected to the production text, character, or image invocation paths inspected for this report.
- Its separate model boundary is therefore dormant by source evidence, not absent from the repository.
- A runtime trace or model invocation would be needed to prove execution, and neither was authorized or performed.

## 5. Q-01 end-to-end frontend composition

### 5.1 Normal load path

The normal load path is:

1. `frontend/dist/js/sidebar-session.js:504-510` uses `syncAllDataToBackend` only for the fallback path.
2. Normal chat selection uses `LoadChat`.
3. `app-chat.go:92-120` loads the persisted chat, calls `updateAppStateFromChat`, sets the active chat ID, and switches backend memory scope.
4. `app-chat.go:500-525` restores transcript, world, scenario, characters, and triggered lore.
5. That function does not assign `a.lorebooks`.
6. `frontend/dist/js/chat-state-bridge.js:287-300` re-syncs world, scenario, and characters to backend, but not lorebooks.
7. `sidebar-session.js:229-231` explicitly avoids a full backend resync on load to preserve chat-specific state.

The consequence is conditional but direct: the backend lore pool is whatever it was before the load unless another action syncs lorebooks.

**Classification:** Source-derived.

### 5.2 Live generic path

The generic text path uses:

- `app-prompt.go:25-100` for the image-capable user-turn entry point.
- `app-prompt.go:277-305` for message append/replacement.
- `app-prompt.go:671-799` for structured message construction.
- Included material includes the system prompt, world fields, scenario fields, character sheets, resolved lore text, and the full retained transcript.
- World and scenario text are inserted as plain model-visible prose.

**Classification:** Source-derived.

### 5.3 Live character path

The character path uses:

- `chat-stream.js:244-327` to sequence active characters.
- `chat-stream.js:108-119` to request a backend lore block.
- `chat-stream.js:121-176` to compose a single character-turn message.
- `app-prompt.go:857-925` to resolve lore with a six-message history window.
- `app-prompt.go:154` receives only the composed message and character ID.
- The final message contains the character card, three scraped system events, earlier responses this turn, inactive-character notes, backend lore text, as many as four character memories, and the user message.
- It does not contain the full world and scenario structures or the full transcript.

**Classification:** Source-derived.

### 5.4 What survives the boundary

Generic structured messages retain role/content boundaries, but not provenance for individual facts. The character path retains even less: all selected facts become one user string.

The backend can know:

- which character ID was requested;
- the current in-memory App world/scenario/characters for the generic path;
- that lore text was produced by the shared resolver.

It cannot recover:

- the memory record from which a memory sentence came;
- the DOM event from which a recent-event sentence came;
- the lorebook owner behind a sentence;
- the chat from which any composed sentence originated;
- the scope or visibility of a flattened sentence.

**Classification:** Source-derived.

## 6. M-01 conversation-scope trace

### 6.1 Storage

`memory-store.js:14` uses the global key:

```js
const MEMORY_STORAGE_KEY = 'howl-chat-memories';
```

`createMemoryEntry` at `memory-store.js:57-104` stores every entry in that global list. It does not add a conversation or chat field.

It supplies `chatID` only as the first argument to backend `AddMemory`. That backend scoping does not become a property used when the frontend composer selects memories.

**Classification:** Observed.

### 6.2 Retrieval

`getCharacterMemories` at `memory-store.js:203-208` filters by:

- `scope === 'character'`;
- `ownerKind === 'character_id'`;
- exact `ownerId`.

There is no chat predicate. Other local functions expose all scenario or global entries, but the character composer only uses the character function.

**Classification:** Observed.

### 6.3 Point of scope loss

Conversation identity is available in several places:

- `sidebarState.currentChatData.sessionId`;
- `dataSync.currentChatId`;
- persisted `ChatMeta.ID`;
- `_state.sessionId`.

The composer does not retrieve using any of them. Scope is therefore lost before retrieval, not merely flattened afterward. Composition further flattens the surviving memory into:

```text
[MEMORIES]
- title: summary/content
```

before sending one user string.

**Classification:** Source-derived.

### 6.4 Backend storage does not save the primary path

Backend memory APIs accept a chat ID at `app-memory.go:67-145` and use `SetChatID`. But the character composer reads localStorage directly and does not call backend retrieval for prompt composition.

Backend chat-scoped storage therefore exists for memory APIs and is not a substitute for the missing frontend prompt scope.

**Classification:** Source-derived.

### 6.5 M-01 answer

No conversation-scope term exists in character-memory retrieval, composition, or the generated request. Character identity is the only effective boundary.

This supports the structural concern behind M-01. It does not prove that any two specific chats leaked data at runtime.

## 7. M-02 lifecycle and inactive-character trace

### 7.1 Sources

Lifecycle material comes from three distinguishable layers:

1. Rendered system events.
2. Current chat roster status.
3. Persistent character memory.

`getRecentSystemMessages` at `chat-stream.js:59-70` scrapes the last three `.chat-message.system` DOM nodes. Those events are composed into `[RECENT EVENTS]` in `chat-stream.js:130-136`.

Inactive characters are identified by `isCharacterActive` and composed into presence notes in `chat-stream.js:145-153`.

Death memories are created in `sidebar-characters.js:402-435` with stable character ownership. They are then eligible for ordinary character-memory retrieval.

**Classification:** Observed.

### 7.2 Recipient filtering

The multi-character loop sends each active character its own composed message. However:

- the same recent system events are supplied to all recipients;
- the same inactive roster is supplied to all recipients;
- prior same-turn responses accumulate into later recipients;
- there is no per-character lifecycle authorization;
- there is no audience receipt or separate structure.

**Classification:** Source-derived.

### 7.3 Persistence

`AddSystemMessage` at `app-prompt.go:410-417` appends only to `a.chatMessages`. It does not append a canonical NDJSON history entry. A snapshot save can preserve the message, but the NDJSON stream does not receive a corresponding lifecycle fact through that method.

**Classification:**

- In-memory behavior: Observed.
- Reload consequence: Inferred from the separate NDJSON and snapshot paths.

### 7.4 M-02 answer

Lifecycle facts can reach a model in two ways:

- transiently, as the last three rendered system events;
- persistently, as a global character death/swap memory.

The transient path is limited to the currently rendered chat. The persistent death path is not limited to that chat if the same stable character is used elsewhere.

## 8. Context inventory used by the live paths

| ID | Fact/content class | Authoritative source | Chat scoping | Selection/retrieval | Filtering | Reduction to model text | Request field |
|---|---|---|---|---:|---|---|---|
| C-01 | Active world | Persisted `ChatData.World`; active `a.world` | Current chat object | Load or world-change action | None for authorization | World prose in generic system prompt | Generic system/content |
| C-02 | Active scenario | Persisted `ChatData.Scenario`; active `a.scenario` | Current chat object | Load or scenario-change action | None for authorization | Scenario prose in generic system prompt | Generic system/content |
| C-03 | Characters and sheets | Persisted/chat roster; active `a.characters` | Current chat roster | Load or character action | Roster membership only | Character prose | Generic system; character card in character user string |
| C-04 | World lore entries | Lore library/backend pool | Pool-dependent, not guaranteed chat-selected | Resolver trigger match | Enabled/stale/hidden; no active-world match | Resolved lore prose | Generic or character request text |
| C-05 | Scenario lore entries | Lore library/backend pool | Pool-dependent | Resolver trigger match | Enabled/stale/hidden; no active-scenario match | Resolved lore prose | Generic or character request text |
| C-06 | Character lore entries | Lore library/backend pool | Stable-owner based | Resolver trigger match | Owner ID or shareable-unowned legacy rule | Resolved lore prose | Generic or character request text |
| C-07 | Automated character memories | Global localStorage | Character-global, not chat-scoped | Stable character ID | No chat check; four-item cap | Memory prose | Character user string |
| C-08 | Explicit chat-swap memory | `data-sync.js:527-547`, scope `chat` | Chat-labeled but consumerless | No prompt retrieval path found | None | Not injected | None found |
| C-09 | Death memory | Death path plus global localStorage | Character-global | Same as C-07 | Same as C-07 | Memory prose | Character user string |
| C-10 | Last three system events | Rendered DOM | Currently rendered chat only | DOM scrape, last three | Count only | Event prose | Character user string |
| C-11 | Inactive status/reasons | Current chat roster | Current chat roster | Current roster | Active/inactive test only | Presence prose | Character user string |
| C-12 | Same-turn prior answers | Current orchestration loop | Current turn only | Accumulated responses | None | Quoted model output | Character user string |
| C-13 | Full generic transcript | Persisted history plus live messages | Current chat conversation | Explicit full-history composition | None found | Message array | Generic messages |
| C-14 | Six-message lore history | Live `a.chatMessages` | Current backend session | Last six non-empty messages | Empty-text skip only | Resolver match input, not directly emitted | Character resolution; only matched lore is emitted |
| C-15 | System/developer framing | Settings and prompt builder | App-level | Every prompt | Roleplay presence only | System prose | Generic system |
| C-16 | User message | Current UI | Current chat turn | Direct dispatch | Empty-message guards | Final user text | User field/string |
| C-17 | Image/attachment | Image path input | Current chat turn | Separate image route | Not fully characterized here | Separate multimodal input | Separate binding; transport remains EO-03 |
| C-18 | Secondary RAG documents | Dormant service | No live scope | No live retrieval path | None; also no audience model | Would be document prose if invoked | Dormant service prompt text |
| C-19 | Consolidated summaries | None found in live prompt paths | Unestablished | No prompt call found | Unestablished | Unestablished | None found |
| C-20 | Promotion/change record | None found | Absent | Absent | Absent | Absent | Absent |

## 9. Q-09 frontend contract and change surface

### 9.1 Actual contract

The current contract is not formally declared. The observed division of responsibility is:

| Layer | Retrieves | Filters | Composes | Sends |
|---|---|---|---|---|
| Frontend state | World/scenario/characters/history/session | Roster membership and panel selection | Panels, save payload, character message | Backend bindings and save data |
| Character composer | Character memories, DOM events, backend lore, prior answers | Stable character identity only | One flat user string | `GenerateCharacterResponse` |
| Backend lore resolver | Active backend lore pool | Enabled/stale/hidden and character ownership | One lore text block | Returned to frontend for flattening |
| Primary generic backend | World/scenario/characters/transcript | No provenance or per-fact authorization | Structured message array | Chat transport |
| Secondary service | Nothing in production | Nothing in production | Nothing in production | No production call |

Frontend code therefore performs retrieval, most effective filtering, and composition. Backend code resolves lore correctly for the one character supplied, but receives no structured information about the other flattened facts.

### 9.2 Affected implementation surface

A provenance-preserving change would have to touch at least:

- `frontend/dist/js/chat-stream.js`
- `frontend/dist/js/memory-store.js`
- `frontend/dist/js/chat-render.js`
- `frontend/dist/js/chat-state-bridge.js`
- `frontend/dist/js/context-selectors.js`
- `frontend/dist/js/sidebar-session.js`
- `app-prompt.go`
- `app-chat.go`
- `app-characters.go`
- `app-memory.go`
- `internal/backend/chat/service.go`
- `internal/backend/chat/types.go`
- `internal/backend/lorebook/lorebook.go`
- `internal/backend/memory/service.go`
- Wails bindings under `frontend/wailsjs/go/main`
- persisted `ChatData`, history, localStorage, and affected tests

Likely affected tests include prompt ordering/boundary tests, lore privacy tests, character prompt tests, memory identity tests, lorebook resolver tests, and chat/HTML integration checks. Compatibility work would involve localStorage migration, persisted history fields, Wails model changes, and the distinction between legacy flattened records and structured scoped records.

No replacement contract is proposed here.

## 10. Q-02 active world/scenario authority trace

### 10.1 Representation

World and scenario authority is split among several structures:

- `a.world` and `a.scenario` are the live backend pointers.
- `ChatData.World` and `ChatData.Scenario` are persisted objects.
- `ChatMeta.WorldID` and `ChatMeta.ScenarioID` exist in `internal/backend/chat/types.go:50-68`.
- Frontend `_state.world` and `_state.scenario` are the UI/session states.
- The secondary service has separate `worldID`, `scenarioID`, `worldName`, and `scenarioName` fields.
- Persisted lore snapshots retain only ID, title, content, and tags.

`updateChatDataFromAppState` at `app-chat.go:417-468` does not update `Meta.WorldID` or `Meta.ScenarioID`. `updateAppStateFromChat` at `app-chat.go:500-525` does not reconstruct lore state. Persisted IDs and live objects can therefore diverge.

### 10.2 All discovered `RouteResolvedEntries` uses

The routing function is defined in `internal/backend/lorebook/lorebook.go:350-378`.

Discovered production-relevant uses are:

- the primary character/general resolver in `app-prompt.go:881-925`;
- the dormant secondary service send paths in `internal/backend/chat/service.go:612-780`;
- lorebook tests.

No other production implementation of active-context comparison was found.

### 10.3 Why the active names do not control routing

`scenarioName` and `worldName` are parameters of the function but are never read inside it. World entries are appended at `lorebook.go:356-357`; scenario entries at `358-359`. `canConsider` at `249-269` checks enabled state, stale roots, hidden visibility, direction, and prior triggering. It does not perform:

- world ownership comparison;
- scenario ownership comparison;
- conversation ownership comparison;
- private-visibility exclusion for world/scenario entries.

`shareableWithoutOwner` applies only to unowned character entries. A world or scenario entry marked private therefore passes routing without an active-context recipient check.

The secondary service sets `ActiveWorld` and `ActiveScenario` before calling `Resolve`, but `canConsider` ignores those fields. They are therefore query metadata in that path, not an authorization decision.

### 10.4 Transitions and retention

- `AddWorldToChat` replaces the live world pointer and saves; it does not clear the lore pool or adjudicate retained lore.
- `AddScenarioToChat` behaves the same way for scenario state.
- `ClearWorldFromSession` sets `a.world=nil`, then calls `performSave`. Because `updateChatDataFromAppState` only overwrites a persisted world when the live world is non-nil, the older persisted `ChatData.World` remains and is reloaded later.
- `ClearScenarioFromSession` has the same stale-object consequence.
- `SetLorebooks` replaces the pool wholesale.
- `UpdateLorebookContext` toggles enablement but does not determine provenance or import authorization.
- The persisted lore snapshot drops scope, owner, visibility, enablement, and stale status.

The source therefore cannot answer whether retained lore was intentionally imported or accidentally left selected. It contains no import marker and no transition authorization.

## 11. Source-level data-flow diagram

```text
Persisted chat
  ChatData.World/Scenario/Characters/History
  ChatMeta.WorldID/ScenarioID
        |
        v
LoadChat / updateAppStateFromChat
  restores world/scenario/characters/transcript
  does not restore backend lore pool
        |
        v
Live App state
  a.world / a.scenario / a.characters
  a.chatMessages
  a.lorebooks, only if separately synced
        |
        +-----------------------------+-----------------------------+
        v                                                           v
Generic chat path                                    Character chat path
buildPromptMessages                                  DOM events + roster
  world/scenario/characters                          getCharacterMemories
  + backend lore + full history                      + backend lore block
  + settings framing                                 + prior responses
        |                                            composed into ONE user string
        v                                                      |
Structured chat messages                                     v
GenerateChatStream                              GenerateCharacterResponse
                                                              no chat ID,
                                                              no scope,
                                                              no memory ID

Global MemoryStore                     Backend memory APIs
character ID only                      chat ID accepted for CRUD
        |                              but not used by live composer
        +---------- flat text ---------+

Dormant chat.Service
constructed empty; RAG flags set;
no production SendMessage/Stream call.
```

Transport-byte proof remains EO-03 and is outside this report.

## 12. Complete inventory of chat-scoped facts that can reach model context

“Can reach” means the source path exists. It does not mean leakage was observed.

| ID | Fact | Live path | Chat boundary preserved? | Cross-chat path by source |
|---|---|---|---|---|
| C-01 | Active world prose | Generic | No; flattened | Only through shared or retained world selection |
| C-02 | Active scenario prose | Generic | No; flattened | Only through shared or retained scenario selection |
| C-03 | Character sheets | Generic and character | Roster remains, provenance does not | Same character/template reused elsewhere |
| C-04 | World lore text | Both live paths, pool permitting | No active-world check | Enabled global or retained entries |
| C-05 | Scenario lore text | Both live paths, pool permitting | No active-scenario check | Enabled global or retained entries |
| C-06 | Character lore text | Both live paths, pool permitting | Stable owner enforced, chat not enforced | Same stable character elsewhere |
| C-07 | Automated character memories | Character | Character only, not conversation | Same stable character elsewhere |
| C-08 | Chat-swap memory | None found | Labeled chat scope but unused | No prompt injection found |
| C-09 | Death memory | Character | Character only, not conversation | Same stable character elsewhere |
| C-10 | Recent system events | Character | Current DOM only | No cross-chat path found |
| C-11 | Inactive notes/status | Character | Current roster only | No cross-chat path found |
| C-12 | Same-turn prior responses | Character | Current turn only | No cross-chat path found |
| C-13 | Full transcript | Generic | Current conversation structurally included | None found beyond chat load/save |
| C-14 | Six-message lore context | Character resolver input | Used for matching only | No direct emission found |
| C-15 | Developer/system framing | Generic | App-level | Not chat-private by design |
| C-16 | User message | All | Current turn | None found |
| C-17 | Image/attachment | Separate image path | Current input | Persistence/transport question remains EO-03 |
| C-18 | RAG material | Dormant only | No live boundary | None live |
| C-19 | Consolidated summaries | None found | Unestablished | Unestablished |
| C-20 | Promotion/change record | None found | Absent | No path exists |

## 13. Findings register

| ID | Finding | Classification | Principal evidence | Established / not established |
|---|---|---|---|---|
| G-01 | Secondary chat service has no production generation caller. | Source-derived | `app.go:170`; `service.go:612,694`; no caller found | Existence established; runtime non-execution not proven |
| G-02 | RAG is attached to the dormant service but cannot inject. | Source-derived | `app.go:186-187`; `service.go:985-994,1028-1052` | Attached but unreachable in production |
| G-03 | `syncAppStateFromChatService` has no caller. | Observed | `app-characters.go:345` | Present but unused |
| G-04 | Normal chat load does not hydrate backend lore pool. | Source-derived | `app-chat.go:500-525`; `chat-state-bridge.js:287-300` | Conditional pool gap; exact runtime pool depends on prior UI actions |
| G-05 | Character memory is globally stored and character-filtered only. | Source-derived | `memory-store.js:14,57-104,203-208` | Scope absence established; particular leakage not observed |
| G-06 | Chat ID is absent from character-memory retrieval and final request. | Source-derived | `chat-stream.js:121-176`; `app-prompt.go:154` | Missing field established |
| G-07 | Recent events are scraped DOM text, capped at three. | Observed | `chat-stream.js:59-70,130-136` | Mechanism established |
| G-08 | Death memories use stable global character ownership. | Observed | `sidebar-characters.js:402-435` | Cross-chat eligibility established by identity reuse |
| G-09 | Chat-scoped swap memories have no prompt consumer. | Source-derived | `data-sync.js:527-547` | Written but not injected |
| G-10 | World/scenario routing ignores active identity. | Observed | `lorebook.go:350-378` | Function behavior established |
| G-11 | Clear-world/clear-scenario can leave stale persisted objects. | Source-derived | `app-chat.go:133-149,417-468` | Logic gap established; actual user impact not runtime-tested |
| G-12 | Backend save does not maintain Meta world/scenario IDs. | Source-derived | `app-chat.go:417-468`; `types.go:50-68` | Missing maintenance established |
| G-13 | Private world/scenario entries are not excluded by routing. | Observed | `lorebook.go:249-269,350-385` | Routing behavior established |
| G-14 | Lorebook save snapshots lose scope/owner/visibility. | Observed | `app-chat.go:434-447`; `types.go:86-94` | Structural loss established |
| G-15 | System lifecycle additions do not append NDJSON history. | Observed/inferred | `app-prompt.go:410-417` | Method behavior observed; reload impact inferred |
| G-16 | Generic and character prompts use materially different context. | Source-derived | `app-prompt.go:671-799`; `chat-stream.js:121-176` | Path divergence established |
| G-17 | No promotion, atomic-scope, or scope-history fields were found. | Observed absence | Memory and lore types; no promotion API | Absence in inspected paths; not a repository-wide impossibility proof |
| G-18 | RAG has no audience/user filtering in the inspected service. | Source-derived | `internal/backend/rag/service.go:119-169`; `service.go:1028-1052` | Dormant design gap; no live exposure |
| G-19 | Active and inactive character facts are flattened without recipient control. | Source-derived | `chat-stream.js:130-153,244-327` | Same facts reach every recipient |

The earlier M-01 and M-02 structural concerns are corroborated. The new material is principally the lore-hydration gap, stale clear behavior, consumerless chat-scoped memory, and the explicit dormant-service reachability conclusion. No frozen report was contradicted on object existence; reachability and scope claims are sharpened rather than assumed.

## 14. Unresolved questions and owner-decision briefs

### D-B1: Active world/scenario authority for lore

Exact question: should world/scenario identity authorize lore access, or merely supply default prompt framing?

Source cannot resolve this because routing ignores both identities, persisted objects and Meta IDs can diverge, and retained lore has no provenance marker.

Alternatives:

- Strict active ownership: only entries originating from the active world/scenario are eligible; transitions remove prior lore unless explicitly imported.
- Selection-based chat scope: explicitly selected chat entries remain eligible even after a transition; active context only supplies defaults.
- Retention with provenance: prior lore may remain, but only as a marked retained/imported entry after confirmation.
- Advisory context: enabled library entries remain eligible regardless of active identity; this matches the current code.

Strict ownership maximizes isolation but harms continuity. Selection preserves user control but can preserve stale lore. Retention with provenance matches AC-C67 best but requires new fields, confirmation, and migration. Advisory context preserves current behavior but conflicts with transition-preservation requirements.

This blocks final Q-02 authorization semantics and any contract that claims active-context filtering.

### D-B2: Memory promotion and default scope

Exact question: is character memory conversation-local by default, and what explicit action promotes it to a broader scope?

Source cannot resolve this because `AC-C70` requires conversation-local default and explicit promotion, while the inspected implementation stores character memories globally and contains no promotion mechanism.

Alternatives:

- Per-entry explicit promotion with recorded prior scope and reason.
- Whole-chat promotion.
- No promotion; all character memory remains persistent.
- Hierarchical promotion across character, scenario, and world scopes.

Per-entry promotion best matches AC-C70 but has the largest UI/API/migration cost. Whole-chat promotion is simpler but coarse. No promotion preserves behavior and violates AC-C70. Hierarchy is expressive but introduces the hardest revocation problem.

This blocks Q-01, M-01, Q-09, and lifecycle-memory isolation.

### D-B3: Canonical conversation identity

Exact question: which identifier defines a conversation for memory and provenance—persisted `ChatMeta.ID`, frontend `sessionId`, backend `currentChatId`, or a new conversation-scoped key after Save As?

Source cannot resolve this because all identifiers exist, Save As changes persisted identity, and no memory field maps one identifier to the others.

Alternatives:

- Persisted chat ID is canonical.
- Frontend session ID is canonical.
- A separate immutable conversation key survives renames and clones.
- Different boundaries are used for persistence, memory, and receipts.

This blocks memory scoping, receipt interpretation, Save As semantics, and migration design.

### D-B4: Lifecycle representation

Exact question: should death, Step Out, return, and ghost participation be structured durable events or transient system prose?

Source cannot resolve this because lifecycle state is split among DOM text, roster status, and global character memories, while AC-C71 requires cross-chat isolation plus originating-chat history.

Alternatives:

- Structured canonical lifecycle events with scope and provenance.
- Transient system messages plus durable roster state.
- Global character memory, as now.
- Conversation-local lifecycle memory with explicit promotion.

Structured events best support AC-C71 but require schema and migration changes. Transient text preserves current UX but cannot prove isolation. Global memory preserves character continuity but exposes lifecycle facts through character reuse. Conversation-local memory matches AC-C70 but requires promotion semantics.

This blocks M-02 closure and AC-C71 acceptance.

## 15. Separately authorized follow-up tests, not executed

No test below was executed.

| ID | Objective | Authorization needed | Scope/data | Risk/cost | Source question it would resolve |
|---|---|---|---|---|
| T-B1 | Determine whether `chat.Service` generation can execute in production. | Runtime/model invocation authorization | Instrumented local launch; temporary trace; no user data | Model load and local app runtime | Q-07 runtime reachability |
| T-B2 | Test whether reused-character memories appear in a second chat. | Runtime/data authorization | Two local chats; synthetic memories only | Local storage and model/data exposure | M-01 actual exposure |
| T-B3 | Test AC-C71 lifecycle isolation and receipts. | Runtime/model authorization | Synthetic death/event; old and new chats; assembled payload inspection | Model/runtime cost; synthetic data only | M-02 and lifecycle isolation |
| T-B4 | Test world/scenario transition with overlapping lore. | Runtime/data authorization | Synthetic worlds/scenarios; payload/receipt inspection | Local runtime and data handling | Q-02 transition/retention |
| T-B5 | Separately verify transport bytes and receipts. | EO-03-specific authorization | Packet/transport observation; no authorization-behavior change | Network capture and data handling | EO-03 only; not Gate B |
| T-B6 | Test Save As/clone identity and memory scope. | Runtime/data authorization | Synthetic chat clone | Local data handling | D-B3 |
| T-B7 | Test dormant RAG audience behavior if separately enabled. | Runtime/data/model authorization | Synthetic RAG corpus; audience separation | Index/runtime/model cost | Q-07/RAG exposure |

## 16. Limitations and unestablished claims

- No runtime, model, benchmark, database, network, packet, transport, or Git-state experiment was performed.
- Production non-execution of `chat.Service` is source-derived; absolute runtime impossibility is unestablished.
- Cross-chat leakage is an identified structural path for reused character identities, not an observed event.
- The lore-pool consequence after any particular load depends on prior UI actions; it is conditional, not a claim that lore is always absent.
- Clear-world and NDJSON consequences follow directly from the inspected save paths, but user-visible reload behavior was not executed.
- Amendment 02 criteria remain unvalidated design targets; this report does not treat them as implemented behavior.
- Exact line numbers refer to the dirty working tree and can shift; the baseline hashes identify the inspected state.
- Negative caller searches covered the relevant Go application and service directories, not an unbounded proof across every conceivable execution mechanism.
- No performance, migration, security-exploit, or replacement-architecture conclusion is made.

## 17. Confirmation of prohibited actions avoided

- No source, configuration, schema, migration, dependency, deployment, or documentation-history change was made except creating this new Gate B report.
- No application, model, benchmark, database, user-data, network, packet-capture, Android, transport, or Git-state operation was performed.
- No frozen R-07, Muse review, erratum, register, or Q-06 artifact was edited, renamed, replaced, or regenerated.
- Read-only Git inspection used status, diff-name, hash, and line-number commands only.
- This report records findings and unresolved owner decisions; it selects no architecture and implements nothing.
