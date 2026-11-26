# Howl Foundation Phase 1 — Chat State and Context Integrity

## Work-package definition against `HFR-CP-2026-11-21-01`

**Status:** work-package definition and authorization request. **No implementation, production API change, schema
migration, model run, G2, BGE or MEMORY-09 integration is authorized by this record.**

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Plan ID | **`HF1-PLAN-2026-11-21-01`** |
| Readiness baseline | **`HFR-CP-2026-11-21-01`** — remains the current baseline; not overwritten |
| Program direction | [`HOWL-MEMORY-09-PROGRAM-DIRECTION-2026-11-21.md`](HOWL-MEMORY-09-PROGRAM-DIRECTION-2026-11-21.md) |
| Readiness audit | [`HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md`](HOWL-FOUNDATION-READINESS-AND-MEMORY-09-DECISION-REGISTER-2026-11-21.md) |
| Scope | Stable conversation identity; branches/restoration; world/scenario/character composition; real token accounting; context flushing; model-change preservation; long/multi-character tests |
| Not in scope | MEMORY-09 capture/classification/lifecycle; production memory integration; hardware lease implementation; inference capability matrix; performance benchmarking |
| Existing defect | `DEF-FE-01` remains open and separate; not absorbed into this package |
| Authorization | Separate implementation authorization required |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Why this package

**HF1-01 (Observed — the readiness baseline).** `HFR-CP-2026-11-21-01` found:

- A is **not ready**: history/ownership is E2; the state machine is E0.
- B is **partial/unverified**: queue and lore privacy exist; orchestration is E1.
- C and D are **partial/unverified**; E is **unestablished**.
- No subsystem has end-to-end verification.
- The combined system is not ready for real-workload acceptance.

**HF1-02 (Calculated — the first dependency is conversation state).** Memory capture, retrieval and lifecycle depend on
knowing which conversation state exists, which message produced a memory, and what context was actually available. A
system that cannot answer those questions cannot preserve provenance.

**HF1-03 (Calculated — dependency priority, not a five-project sequence).** The priority order is:

| Priority | Area | Reason |
|---|---|---|
| **1** | Chat state identity and context budget | Everything else depends on reliable conversation state |
| **2** | Multi-character orchestration | Needed for realistic memory testing |
| **3** | Hardware resource management | Prevents instability under actual workloads |
| **4** | Real inference and KV-cache verification | Establishes what the engine really supports |
| **5** | Performance benchmarking | Tells us which optimizations are worth implementing |

**HF1-04 (Observed — parallel work is allowed).** Hardware and inference investigation can proceed alongside chat work.
This package does not block them, and they do not block the identity/token work.

---

## 2. Non-disappearing findings

**HF1-05 (Observed — context budget is not implemented).** `TrimToMaxTokens` is a placeholder; `Message.Tokens` is
always `0`; the app path has no real token budget. Long-conversation behaviour cannot be trusted.

**HF1-06 (Observed — multi-character contamination is untested).** A queue exists, but no test establishes that
characters with conflicting knowledge, different histories, different voices and different memory access remain
isolated.

**HF1-07 (Observed — hardware management is advisory).** `RequestAllocation` returns advice; no lease enforces it. No
test covers telemetry loss, allocation conflicts, model swaps, VRAM pressure or recovery.

**HF1-08 (Calculated — carry-forward register).**

| Finding | Where it is handled | Status |
|---|---|---|
| Context budget | **This Phase 1 package** | Open |
| Multi-character contamination | **This Phase 1 test plan**, then integration | Open |
| Hardware fault testing | **Phase 3 dependency**; not implemented here | Open |
| Existing frontend failure | **Separate defect `DEF-FE-01`** | Open |

---

## 3. Separate defect — do not absorb

**HF1-09 (Observed — `DEF-FE-01`).**

| Field | Value |
|---|---|
| Test | `tests/frontend/character-configuration-integrity.test.js:409` |
| Name | “an explicit voice-model clear works while an unrepresentable stored voice is preserved” |
| Result on current tree | **FAIL** — actual `custom_voice_not_in_the_fixed_select`, expected `''` |
| Relation to Phase 1 | Separate frontend character-configuration defect |
| Disposition | **Open; not fixed by this package; not hidden by a green-suite claim** |

**HF1-10 (Calculated).** Phase 1 acceptance may proceed with `DEF-FE-01` still open, but must report it explicitly.
Phase 1 may not claim the frontend suite is green while this failure remains.

---

## 4. Scope

### 4.1 In scope

**HF1-11 (Observed — proposed scope).**

- Stable **state**, **message** and **turn** identities.
- A non-destructive branch/alternative model with explicit current-state selection.
- Revert/restore that changes the current pointer, not the retained material.
- Explicit deletion semantics: only deletion removes material; deletion names its set.
- World/scenario/character composition by stable ID; names are display labels only.
- Real token counting through the model tokenizer, with an explicit unavailable state.
- A real context budget covering system/context/history/generation reserve.
- Predictable context flushing between character responses.
- Preservation of authoritative history across model changes.
- Tests for long conversations, branching, restart and multiple characters.

### 4.2 Out of scope

**HF1-12 (Observed).**

- MEMORY-09 capture, classification, lifecycle, retrieval or candidate UI.
- Production memory integration or frontend memory authority changes.
- Hardware lease implementation or OOM recovery.
- Real inference capability matrix or KV-cache optimization.
- Performance benchmarking.
- G2, BGE, corpus or model fetch.
- Fixing `DEF-FE-01`.
- Canonical r33/v9/ledger/map changes.

---

## 5. Completion criteria

**HF1-13 (Calculated — completion means evidence, not code volume).** The package is complete only when all named
criteria have evidence at the stated level. A criterion at E1 is not complete.

### 5.1 Stable identity

**HF1-14 (Calculated — required).**

- Every state has a stable `state_id` that survives save/load, edits, regeneration, branch switching, reorder and
  deletion of other states.
- Every message has a stable `message_id` assigned at append; it is never derived from list position.
- Every turn has a stable `turn_id`/sequence for local ordering.
- `ChatMeta` identifies the current state; branches retain their parent/alternative relationships.
- Legacy records remain **unresolved**; no ID is invented for them.
- IDs are unique within a chat and deterministic only where explicitly specified; no index-derived identity.

**Acceptance evidence:** an append → save → restart → branch → revert → restart test in which all retained IDs are
unchanged and no legacy record gains a synthetic ID.

### 5.2 Non-destructive branches and restoration

**HF1-15 (Calculated — required).**

- Producing an alternative response creates a distinct state; the prior state remains retained.
- Regeneration creates an alternative state; it does not silently overwrite the prior material.
- Revert selects a retained state as current; it creates, deletes and rewrites nothing.
- Branch switching changes the current pointer and the context composition; it never deletes memory or history.
- Deletion is the only destructive operation and names the removed set.

**Acceptance evidence:** a branch/alternative/revert scenario with a before/after inventory proving retained material is
unchanged; a deletion scenario with an explicit removed-set receipt.

### 5.3 World/scenario/character composition

**HF1-16 (Calculated — required).**

- Composition uses stable world/scenario/character IDs from `ChatMeta` and loaded `ChatData`.
- Display names are labels; a name never substitutes for a missing ID.
- Switching world/scenario/character changes the active composition without mutating the underlying authored data.
- The prompt contains the correct active world, scenario and character blocks in the established order.
- Restoring a chat restores the same composition and the same current state.

**Acceptance evidence:** a two-world/two-character switch test that asserts active IDs, prompt composition, and
restoration after restart; a missing-ID case that remains unresolved rather than falling back to a name.

### 5.4 Real token counting and budgeting

**HF1-17 (Observed — the tokenizer path exists, but not at the app seam).** `llama.HTTPClient.Tokenize` calls
llama-server's `POST /tokenize`; `TokenizerOptions` and `GetContextSize` exist. However, `App.inferenceClient`
(`app.go`) does **not** include `Tokenize`, so the app cannot use it through the current interface. Phase 1 must extend
the interface or add a tokenizer seam, and update test fakes accordingly. This is an available mechanism, not yet a
context-budget implementation.

**HF1-18 (Calculated — required).**

- Token counts come from the loaded model's tokenizer, not a character-count estimate.
- `Message.Tokens` is populated; persisted token counts are invalid after a model change and are recomputed.
- The budget includes system/instructions, active world/scenario/character context, recent history and a generation
  reserve.
- Over-budget behaviour is explicit: retained material is not silently rewritten; the prompt working set is reduced
  with a recorded receipt of what was included, omitted or summarized.
- If the tokenizer is unavailable, budget state is **unavailable**; the system does not substitute an estimate and does
  not claim compliance.

**Acceptance evidence:** a prompt whose token count is measured by the real tokenizer; a deliberately over-budget
conversation; a model-swap case with recomputed counts; a tokenizer-unavailable case that reports unavailable.

### 5.5 Context flushing between character responses

**HF1-19 (Calculated — required).**

- Each character response receives a bounded working set: identity, active world/scenario, relevant scene/context,
  recent transcript within budget, and any explicitly scoped test fixture.
- Flushing is a **prompt working-set** operation, not a memory-lifecycle event. It never deletes authoritative history
  and never changes memory status. (S2H-12 remains intact.)
- The policy is deterministic and observable: the same state and budget produce the same inclusion/omission decision.
- One character's private context does not enter another character's working set.
- Prior responses included as shared scene context are labelled and bounded.

**Acceptance evidence:** a two-character scene with conflicting private facts and different histories; each prompt
contains only its character's private context; the shared scene remains available; the authoritative history is
unchanged.

### 5.6 Model-change preservation

**HF1-20 (Calculated — required).**

- Unloading/reloading or switching models does not alter `ChatHistory`, state IDs, branch structure or current-state
  selection.
- Memory is not required for this test; the test may use synthetic state fixtures.
- Token counts are recomputed for the new tokenizer; old counts are not reused as authority.
- No behaviour depends on a surviving KV cache.

**Acceptance evidence:** save → unload → reload/switch → load; IDs and history unchanged; token counts recomputed; context
composition restored.

### 5.7 Long and multi-character tests

**HF1-21 (Calculated — required).**

- A long conversation that exceeds the budget exercises real flushing and preserves authoritative history.
- Multiple characters with conflicting knowledge, different histories, different voices and different memory access
  remain isolated.
- Restart reconstructs the same current state and branch structure.
- The tests must be end-to-end at the chat/state boundary, not only unit tests of a helper.

**Acceptance evidence:** named E3 scenarios in §6, with retained logs/artifacts and no ID or content mutation.

---

## 6. Acceptance scenarios

**HF1-22 (Calculated — all must pass for Phase 1).**

| ID | Scenario | Required result | Evidence level |
|---|---|---|---|
| **AC-1** | Append, save, restart | All message/state/turn IDs unchanged; legacy entries unresolved | **E3** |
| **AC-2** | Branch, alternative, revert | Prior state retained; current pointer changes; no deletion | **E3** |
| **AC-3** | Explicit deletion | Named set removed; retained branches/history intact; receipt recorded | **E3** |
| **AC-4** | World/scenario/character switch | Correct active IDs and prompt composition; restoration matches | **E3** |
| **AC-5** | Over-budget long conversation | Real token count; explicit flushing receipt; authoritative history unchanged | **E3** |
| **AC-6** | Two characters, conflicting private facts | No private context or voice contamination; shared scene remains | **E3** |
| **AC-7** | Model unload/reload or switch | History/state IDs unchanged; token counts recomputed | **E3** |
| **AC-8** | Tokenizer unavailable | Budget reports unavailable; no estimate is presented as compliance | **E2/E3** |

**HF1-23 (Calculated — E3 is required for AC-1…AC-7).** A unit test of an ID helper is not AC-1. A prompt-string test is
not AC-4. A tokenizer unit test is not AC-5. E3 means the scenario crosses the relevant boundaries and the evidence is
retained.

---

## 7. Proposed implementation slices — not authorized

**HF1-24 (Calculated — slices are design inputs, not a build order).**

| Slice | Purpose | Likely touch points |
|---|---|---|
| **P1.1 Identity** | State/message/turn IDs; branch/parent model; legacy unresolved | `chat/types.go`, `chat/persistence.go`, `app-prompt.go`, `app-chat.go` |
| **P1.2 Token accounting** | Real tokenizer counts; budget state; recompute on model change | `types/context.go`, `types/message.go`, `llama/http_client.go`, `app-prompt.go` |
| **P1.3 State model** | Current state, alternatives, revert, deletion set | `chat/types.go`, `chat/persistence.go`, frontend chat state/actions |
| **P1.4 Composition and flushing** | Stable-ID world/scenario/character composition; bounded working set | `app-prompt.go`, `chat-stream.js`, `chat-state-bridge.js` |
| **P1.5 Test harness** | Long/multi-character/restart/model-change scenarios | chat/app tests; frontend tests |

**HF1-25 (Calculated — likely files, no edits).**

| File | Expected change |
|---|---|
| `internal/backend/chat/types.go` | State/message/turn IDs; branch/parent/current fields |
| `internal/backend/chat/persistence.go` | Persist/load IDs; no legacy backfill |
| `internal/backend/types/context.go` | Real budget state; remove placeholder-only trim |
| `internal/backend/types/message.go` | Populate token count; model identity for recomputation |
| `internal/backend/llama/http_client.go` | Use existing tokenizer path for counts; no new estimate |
| `app.go` | Extend `inferenceClient` with tokenization or add a tokenizer seam; update test fakes |
| `app-prompt.go` | Stable append IDs; budget-aware working set; composition |
| `app-chat.go` | Preserve IDs in build/save/load; current-state handling |
| `frontend/dist/js/chat-actions.js` | Branch/revert/delete semantics; no destructive rewrite for non-destructive actions |
| `frontend/dist/js/chat-stream.js` | Bounded per-character working set; no cross-character private context |
| `frontend/dist/js/chat-session.js` | Persist current state/branch metadata |
| `frontend/dist/js/chat-state-bridge.js` | Restore current state and composition |
| `frontend/wailsjs/go/main/App.d.ts` | Regenerated if bindings change; never hand-edited |

---

## 8. Test plan

**HF1-26 (Calculated — required tests).**

| Layer | Test | What it proves |
|---|---|---|
| Unit | ID assignment and persistence | No positional identity; legacy unresolved |
| Unit | Branch/alternative/current-state model | Non-destructive transitions |
| Unit | Tokenizer count and budget arithmetic | Real counts; unavailable state; reserve |
| Integration | Append → save → restart | IDs and composition survive |
| Integration | Branch → alternative → revert → restart | Retained material and current selection |
| Integration | World/scenario/character switch | Stable-ID composition; no name fallback |
| Integration | Over-budget conversation | Flushing receipt; authoritative history unchanged |
| Integration | Two characters with conflicting private facts | No private context/voice contamination |
| Integration | Model unload/reload | State survives independently of inference |
| Regression | Existing chat/character/prompt tests | No unrelated breakage |
| Separate | `DEF-FE-01` | Tracked, not absorbed |

**HF1-27 (Calculated — no MEMORY-09 integration test in Phase 1).** The multi-character tests use explicit test
fixtures for private facts/history. They do not require the memory service and must not create a hidden memory
integration path.

---

## 9. Acceptance gate

**HF1-28 (Calculated — Phase 1 is complete only when):**

1. AC-1…AC-7 are E3 with retained evidence.
2. AC-8 reports unavailable correctly.
3. No message/state/turn ID is derived from list position.
4. No budget decision uses a character-count estimate as authority.
5. No non-destructive transition deletes authoritative history.
6. No character receives another character's private context.
7. Model unload/reload does not change state identity.
8. Existing unrelated tests have no new failures.
9. `DEF-FE-01` remains visible and separate; the suite is not reported green while it fails.
10. No MEMORY-09 integration, G2, BGE, production API or canonical change was introduced.

**HF1-29 (Calculated — stop rather than claim).** If any criterion is met only by a mock that bypasses the real
boundary, it is not met. If the tokenizer is unavailable, the budget is unavailable, not estimated.

---

## 10. Dependencies and sequencing

**HF1-30 (Calculated).**

| Dependency | Effect |
|---|---|
| **A state identity** | Phase 1 owns it; no MEMORY-09 dependency |
| **B multi-character** | Phase 1 tests can begin after identity; orchestration tests are part of AC-6 |
| **C hardware** | Phase 1 does not implement it; model-change test may use load/unload stubs until C exists |
| **D inference** | Phase 1 uses the existing tokenizer; real-server tokenization is required for AC-5/AC-7 |
| **M1 capture** | Phase 1 output is an input, not a consumer |
| **M8 contract** | Phase 1 state/branch semantics feed the future contract |
| **Hardware fault testing** | Carry-forward to Phase 3; not a Phase 1 completion criterion |
| **Performance** | Carry-forward to Phase 5; no benchmark claim here |

**HF1-31 (Calculated — parallel, not blocked).** C/D investigation can proceed alongside Phase 1. Phase 1 can proceed
without memory integration. Integration waits for both.

---

## 11. Risks and stop conditions

**HF1-32 (Calculated).**

| Risk | Handling |
|---|---|
| Tokenizer unavailable | Report `unavailable`; do not estimate; stop budget-dependent acceptance |
| Model tokenizer changes | Recompute counts; never reuse persisted counts as authority |
| Branch model scope creep | Start with state/message/turn IDs, parent/alternative, current pointer; defer richer lineage |
| Legacy records | Leave unresolved; no invented IDs or owners |
| Context flush deletes history | Flush only the prompt working set; authoritative history unchanged |
| Cross-character leakage | Private context excluded per character; test with conflicting facts |
| Frontend test failure absorbed | Keep `DEF-FE-01` separate; no green-suite claim |
| Hardware dependency | Do not claim hardware readiness from Phase 1; carry forward |
| Production API creep | Separate authorization; no binding changes in this definition |

**HF1-33 (Calculated — stop and return if):**

- accurate token counting requires a model/tokenizer that is not available and an estimate would be used instead;
- the branch model requires an owner retention decision not yet made;
- context flushing would require deleting or rewriting authoritative history;
- the work requires changing memory capture/classification/lifecycle;
- the work requires a production API/schema change beyond the package’s authorization;
- `DEF-FE-01` is being fixed inside this package instead of tracked separately.

---

## 12. Authorization

**HF1-34 (Observed — no implementation authorization).** This record defines the package and requests a separate
implementation decision. It does not authorize code, schema, migration, bindings, model runs or tests that mutate
production data.

**Suggested authorization boundary, if approved later:**

- Phase 1 files and tests only;
- no MEMORY-09 service changes;
- no production capture/classification path;
- no hardware lease or inference optimization;
- no canonical edits;
- `DEF-FE-01` remains separate;
- evidence at AC-1…AC-8 required before completion.

---

## 13. Traceability

**HF1-35 (Observed).**

| Source | Use |
|---|---|
| `HFR-CP-2026-11-21-01` | Readiness baseline; A not ready, B partial, context placeholder |
| HSL-1…HSL-11 | State ownership, identity, branches, revert, deletion, retention, notification |
| `HFR-12` | Positional IDs and missing state machine |
| `HFR-13` | `TrimToMaxTokens` placeholder and `Message.Tokens = 0` |
| `HFR-18` | Multi-character orchestration E1 |
| `HFR-23` | Hardware advisory/lease gap |
| `HFR-42` | Integration scenarios not end-to-end |
| `PCB-25…31` | Stable source identity and append-boundary proposal |
| `DEF-FE-01` | Separate frontend defect |

---

## 14. Record status

**HF1-36 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Production code/schema edited | **0** |
| Requirements changed | **0** — map **118** |
| L-2 / Decision A amended | **0** |
| G2 / BGE started | **0** |
| Models fetched / corpus created | **0** |
| Tests executed | **0** — planning record only |
| Readiness baseline | `HFR-CP-2026-11-21-01` unchanged |
| This file | new, untracked; **work-package definition**, not implementation |

**End of package.** Phase 1 is the bounded path from the readiness baseline to stable conversation state and a real
context budget. It does not integrate MEMORY-09, does not claim hardware/inference readiness, and does not absorb the
existing frontend defect. Completion requires E3 evidence for AC-1…AC-7, an honest unavailable state for AC-8, and a
visible separate `DEF-FE-01` (HF1-22…HF1-29). Implementation requires separate authorization (HF1-34).
