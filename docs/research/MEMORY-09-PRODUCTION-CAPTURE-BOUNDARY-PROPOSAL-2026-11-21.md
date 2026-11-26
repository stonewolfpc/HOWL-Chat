# MEMORY-09 — Production Capture Boundary Proposal

## What must change before the memory service can receive the G1 envelope

**Status:** **read-only trace and change proposal.** No production API change, schema, code, migration, test, model
fetch, corpus, L-2 amendment or Decision A amendment. G2 remains blocked. The fixture contract passed; production
capture is still unestablished.

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Scope | trace the production path and propose the boundary change |
| Finding | two independent breaks: frontend candidate bindings are absent, and the Go append path has no stable source ID and never calls the memory service |
| Recommendation | server-side capture at the append boundary, with a stable source envelope |
| Production change | **not made**; separate authorization required |
| L-2 / Decision A / BGE | **unchanged** |
| G2 | **blocked** |
| Map | **118** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. The two breaks

**PCB-01 (Observed — frontend break).** `frontend/dist/js/memory-candidates-menu.js` checks
`window.go.main.App.ProcessMessageForEvents` and calls it with `(message, [], [])`. `frontend/wailsjs/go/main/App.d.ts`
has **no** `ProcessMessageForEvents` binding.

**PCB-02 (Observed — fallback bypass).** When the binding is absent, the wrapper falls back to frontend keyword
detection on the **user's outgoing message** and stores a candidate in `localStorage` via
`MemoryCandidates.addLocalCandidate`. The backend candidate queue is not used.

**PCB-03 (Observed — approval bypass).** `memory-candidates.js` checks `App.ApproveCandidate`; `App.d.ts` has no such
binding. The fallback `approveLocalCandidate` calls `MemoryStore.createMemoryEntry`, which writes localStorage and
calls `App.AddMemory`. `App.AddMemory` exists and creates a record directly.

**PCB-04 (Observed — backend candidate API unused).** `memory.Service` has `GetCandidates`, `ApproveCandidate`,
`RejectCandidate` and `GetPendingCandidateCount`, but there are no App bindings for them. The frontend cannot reach the
backend queue.

**PCB-05 (Observed — Go append path).** `App.appendChatMessage(role, content, sender, characterID)` appends a map to
`a.chatMessages`, builds a `ChatHistory` with
`ID: fmt.Sprintf("%d", len(a.chatMessages)-1)`, and calls `AppendHistoryEntry`. It never calls the memory service.

**PCB-06 (Observed — reindexing).** `App.buildChatHistory` assigns `ID: fmt.Sprintf("%d", len(history))` again;
`App.ReplaceChatHistory` assigns `ID: fmt.Sprintf("%d", i)`. Source message IDs are positional, not stable across
edits, deletions, reloads or save snapshots.

**PCB-07 (Observed — no capture call).** The memory service is initialized in `App` and in `chat.Service`, but the
chat send/append path calls only `memoryService.SaveChatData` or nothing; `ProcessMessageForEvents` has no Go caller.

---

## 2. The actual path

**PCB-08 (Observed — generic stream).**

1. `memory-candidates-menu.js` wraps `window.sendMessage`.
2. `App.SendMessageStreamWithImage` calls `recordUserTurn`, streams the response, then
   `AppendChatMessage("assistant", ..., "AI")`.
3. `App.appendChatMessage` appends the map and `ChatHistory` to `a.chatMessages` / `history.ndjson`.
4. No memory capture occurs.

**PCB-09 (Observed — character stream).**

1. The frontend builds a character-specific prompt and calls `App.GenerateCharacterResponse(message, characterID)`.
2. `GenerateCharacterResponse` appends the assistant reply with
   `appendChatMessage("assistant", ..., character.Name, characterID)`.
3. The app has the stable character ID at this point.
4. No memory capture occurs.

**PCB-10 (Observed — user/system turns).** The frontend calls `App.AppendChatMessage("user"|"system", content, sender)`;
`appendChatMessage` receives no character ID. The app has `a.world`, `a.scenario`, `a.characters`, and can load
`ChatData` for the active chat.

**PCB-11 (Observed — candidate path).** Backend `ProcessMessageForEvents` would call `DetectEvents` → `BatchClassify` →
`AddCandidates`, but no App method or chat path calls it. The frontend fallback creates local candidates instead.

---

## 3. Field availability at the boundary

**PCB-12 (Observed).**

| Field | Available at `appendChatMessage`? | Source |
|---|---|---|
| chat ID | yes | `chatPersistence.GetLastChatID()` |
| role | yes | argument |
| content | yes | argument |
| sender | yes | argument |
| character ID | yes for character path; empty for user/system/generic | argument |
| world/scenario IDs | yes via `a.world`/`a.scenario` or loaded `ChatData.Meta` | app/chat data |
| world/scenario/character metadata | yes via `a.world`, `a.scenario`, `a.characters` or `LoadChat` | app/chat data |
| source message ID | **no stable ID** | positional `len(...)` |
| source span | text only; no offsets | message text |
| turn ID | **no** | absent |
| attribution | **no** | absent |
| authored asset reference | **no** | absent |
| event type/evidence | **no** at append | detector only |

**PCB-13 (Calculated — the boundary has enough to start, but not enough to finish).** The append hook can supply
chat/role/content/character/world/scenario. It cannot supply stable source identity, turn, attribution or authored
asset status without additional fields.

---

## 4. Origin distinctions

**PCB-14 (Calculated — source kinds).**

| Source kind | Detection | Class consequence |
|---|---|---|
| `user_message` | role `user`, no character ID | user-source class rule deferred |
| `character_message` | role `assistant`, stable `CharacterID` | eligible for `experiential` / `character_specific` |
| `narrator_message` | role `assistant` or `system`, **no** `CharacterID` | not a character; authored/ambient only with an authored ref |
| `authored_world` / `authored_scenario` | content-library/import source, not chat history | `ambient` / `restricted` if authored |
| `legacy` | missing role/character/source | capture-incomplete |

**PCB-15 (Calculated — no inference from names).** A character message must use `CharacterID`. A display name or
`Sender` is never ownership.

**PCB-16 (Calculated — no ambient from subject matter).** A character message about the world is still a character
message. It gets no `authored_asset_ref` and is not ambient/restricted.

---

## 5. Stable source identity

**PCB-17 (Observed — positional IDs).** `appendChatMessage`, `buildChatHistory` and `ReplaceChatHistory` all assign
IDs from list position. The frontend maps in `a.chatMessages` do not carry a stable `id`; `convertFrontendDataToChatData`
reads an `id` that is usually absent.

**PCB-18 (Calculated — required change).** Assign a unique stable message ID at append, store it in the in-memory map
and `ChatHistory`, and preserve it through `buildChatHistory`, `ReplaceChatHistory`, `convertFrontendDataToChatData`,
`AppendHistoryEntry` and load. Do not derive it from list position.

**PCB-19 (Calculated — turn identity).** Add a persisted `turn_id`/sequence at append. A list index is not a turn
identity after deletion or reorder.

**PCB-20 (Calculated — source span).** For the first boundary, `source_span` can be the full message text plus the
stable message ID. Exact character offsets are a later refinement; the ID must still be stable.

---

## 6. Attribution and authored status

**PCB-21 (Observed — no attribution field).** `ChatHistory`, `types.Message` and `MemoryCandidate` carry no
attribution. The app cannot currently preserve observed/told/believed/asserted/authoritative.

**PCB-22 (Calculated — only role-based assertion is safe).**

| Source | Safe attribution | Not safe |
|---|---|---|
| user message | `asserted` | `observed`, `told`, `believed`, `authoritative` |
| character message | `asserted` | inferring `observed` from first person |
| narrator/system | **unresolved** unless explicitly supplied | `authoritative` by default |
| authored world/scenario | explicit authored source | inferring from chat text |

**PCB-23 (Calculated — the character-tell case cannot be produced from current data).** `told` requires a transfer
relation to a listener; a character saying “I grew up in the north” does not establish `told`. The fixture uses `told`
as an explicit test input; production must not infer it. The character-tell memory can be captured as `asserted` and
left classification-unresolved, or explicitly supplied as `told` by a future transfer source.

**PCB-24 (Calculated — authored asset reference).** `authored_asset_ref` must be set only when capturing an authored
world/scenario/character definition from the content-library/import path. A normal chat message gets none; its
world-related subject matter does not create one.

---

## 7. Proposed boundary

**PCB-25 (Calculated — recommend server-side capture at append).** Add a capture call in
`App.appendChatMessage` after a successful append. The app already has role, content, sender, character ID and chat ID;
it can load `ChatData` for world/scenario metadata and call the memory service with a capture envelope. This avoids the
missing frontend bindings and the empty-array call.

**PCB-26 (Calculated — why not frontend capture).** The frontend path is already broken, has no bindings for the
candidate queue, passes empty history/identities, and duplicates candidates in localStorage. Fixing only the binding
would still leave the source envelope incomplete and the approval path bypassed.

**PCB-27 (Calculated — why not `chat.Service` capture alone).** The App send path uses `a.httpClient` directly, not
`chat.Service.SendMessage`; the chat service context would not see every character/user turn without a larger
refactor. The append boundary is the common path.

**PCB-28 (Calculated — logical API).**

```
CaptureInput{
  chat_id, turn_id, source_message_id, source_span,
  speaker_role, source_character_id, source_character_name,
  source_kind, attribution, authored_asset_ref,
  world_id, scenario_id, event_type, evidence, content
}

CaptureResult{
  state: complete|incomplete,
  missing_fields[],
  candidate_id
}
```

The capture call does **not** assign a class. It writes a candidate with capture state and missing fields.

**PCB-29 (Calculated — candidate queue).** `MemoryCandidate` gains capture state/missing fields, not a class.
`ApproveCandidate` rejects capture-incomplete candidates. `GetCandidates` / `ApproveCandidate` / `RejectCandidate` /
`GetPendingCandidateCount` get App bindings so the frontend uses the backend queue.

**PCB-30 (Calculated — remove the duplicate path).** The frontend keyword fallback and `MemoryStore.createMemoryEntry`
approval path must not create memories independently. Until the backend candidate bindings exist, the fallback should
be disabled or labelled diagnostic; otherwise duplicate and provenance-free records continue.

**PCB-31 (Calculated — `AddMemory` boundary).** `App.AddMemory` remains a manual path, but approved candidates should
not use it directly. If `AddMemory` is retained, it must mark capture state incomplete when the envelope is absent; it
must not imply the memory service captured the source.

---

## 8. Exact files likely to change

**PCB-32 (Calculated — proposal, not applied).**

| File | Proposed change |
|---|---|
| `app-prompt.go` | `appendChatMessage`: assign stable ID/turn; call memory capture after append |
| `app-chat.go` | `buildChatHistory` / `ReplaceChatHistory`: preserve stable IDs and turn |
| `app-chat.go` / `app-memory.go` | App bindings for candidate queue; optional capture diagnostics |
| `internal/backend/chat/types.go` | stable message ID/turn fields on `ChatHistory` |
| `internal/backend/memory/service.go` | `CaptureInput` / `ProcessCapture` |
| `internal/backend/memory/candidates.go` | capture state/missing fields; reject incomplete approval |
| `internal/backend/memory/classify.go` | consume envelope; no class assignment in this step |
| `app-memory.go` | candidate DTO/bindings; `AddMemory` capture-state handling |
| `frontend/wailsjs/go/main/App.d.ts` | regenerated bindings; do not hand-edit |
| `frontend/dist/js/memory-candidates-menu.js` | remove/disable duplicate fallback; use backend |
| `frontend/dist/js/memory-candidates.js` | use backend queue; no localStorage approval bypass |

---

## 9. Migration and legacy

**PCB-33 (Calculated — no backfill).** Existing history entries with positional IDs remain legacy. Do not invent stable
IDs or attribution. New appends get stable IDs; legacy entries are reported as capture-incomplete/unmigrated.

**PCB-34 (Calculated — no class).** This boundary does not add `KnowledgeClass`; class remains the later G2/G4
question. The candidate can be captured, stored and inspected without a class.

---

## 10. Tests and acceptance

**PCB-35 (Calculated — integration tests).**

- `appendChatMessage` for F1–F10 creates the expected envelope; no class.
- Stable message ID survives append → save → load → replace.
- Character message uses `CharacterID`; user/system gets explicit none.
- F8 gets no authored ref and is not ambient/restricted.
- F6 missing attribution is capture-incomplete.
- Character-tell is captured as `asserted`, not inferred `told`; class unresolved/deferred.
- Candidate approval rejects capture-incomplete.
- Frontend no longer creates duplicate local candidates.

**PCB-36 (Calculated — no production change without separate authorization).** This record is a proposal. The test-only
G1 harness remains unchanged; G2 is blocked until the production boundary is implemented and verified.

---

## 11. Risks and stop conditions

**PCB-37 (Calculated — risks).**

| Risk | Handling |
|---|---|
| stable ID migration | preserve existing IDs; new IDs only; no backfill |
| duplicate capture | capture once at append; idempotency by source message ID |
| append succeeds, capture fails | message remains; capture-incomplete logged; no class |
| frontend duplicate path | disable fallback; backend queue only |
| attribution inference | only `asserted` from role; otherwise unresolved |
| ambient inference | no authored ref for chat messages |
| character-tell `told` | cannot infer; explicit source only |
| scope creep into class | stop; class belongs to G2/G4 |

**PCB-38 (Calculated — stop conditions).** Stop and return if the boundary cannot assign stable source IDs, cannot
distinguish character/narrator/authored origins, cannot preserve `asserted` without inference, or requires an L-2/class
decision. A production API change is a separate authorization.

---

## 12. Decision requested

**PCB-39 (Observed).** Authorize a **production-boundary implementation** with this scope:

1. stable source message ID and turn;
2. `CaptureInput` / `ProcessCapture` at the append boundary;
3. capture state and missing fields on candidates;
4. backend candidate bindings and removal of the duplicate frontend path;
5. tests for F1–F10, character-tell, stable IDs and capture-incomplete rejection.

No class, no L-2 amendment, no BGE, no G2. The existing G1 fixture contract stays as the test model.

---

## 13. What this proposal does not do

**PCB-40 (Observed).**

| Untouched | State |
|---|---|
| Production code/API | **unchanged** |
| Schema/migration | **0** |
| L-2 / Decision A | **unchanged** |
| BGE / G2 | **not started** |
| Class assignment | **0** |
| Requirements / map | **0**; **118** |
| Canonical capture | **0** |
| Test-only fixture | unchanged; separate |

---

## 14. Traceability

**PCB-41 (Observed).**

| What | Where |
|---|---|
| Frontend break | `memory-candidates-menu.js`; `App.d.ts`; PCB-01/02 |
| Approval bypass | `memory-candidates.js`; `memory-store.js`; PCB-03/04 |
| Append path | `app-prompt.go`; PCB-05/06 |
| Character path | `GenerateCharacterResponse`; PCB-09 |
| Stable IDs | PCB-17/18 |
| Attribution | PCB-21…23 |
| Authored status | K-R8; PCB-16/24 |
| G1 fixture | GFC-01…31; G1R-22 |
| L-2 / Decision A | L2D-01/08/10; DAO-05 |
| Map | **118** |

---

## 15. Record status

**PCB-42 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Production code/schema edited | **0** |
| Requirements changed | **0** — map **118** |
| L-2 / Decision A amended | **0** |
| BGE / G2 started | **0** |
| Tests executed | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **proposal** — no canonical capture |

**End of proposal.** The frontend calls candidate methods with no App bindings, falls back to localStorage, and
approves through `AddMemory`; the Go append path uses positional IDs and never calls the memory service (PCB-01…11).
The append boundary already has role/content/character/chat context but lacks stable source identity, turn, attribution
and authored status (PCB-12…24). The proposed change is server-side capture at `appendChatMessage` with a stable
envelope, capture state/missing fields, and backend candidate bindings; no class, L-2 amendment, BGE or G2
(PCB-25…36). **No production change follows without separate authorization** (PCB-37…42).
