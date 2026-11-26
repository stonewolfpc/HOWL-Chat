# HOWL Foundation Readiness and MEMORY-09 Decision Register

## Consolidated readiness audit — checkpoint `HFR-CP-2026-11-21-01`

**Status:** readiness audit only. **No** production code, schema, API, model, GPU stress test, corpus, G2 or BGE work.
This record does not authorize implementation and does not reopen settled MEMORY-09 decisions.

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Audit scope | Howl foundation A–E; MEMORY-09 decision register; integration readiness |
| Checkpoint | **`HFR-CP-2026-11-21-01`** |
| Program checkpoint | `CP-2026-11-21-01` remains the program-level recovery point |
| Program direction | [`HOWL-MEMORY-09-PROGRAM-DIRECTION-2026-11-21.md`](HOWL-MEMORY-09-PROGRAM-DIRECTION-2026-11-21.md) |
| Canon | r33, v9, owner ledger and conformance map unchanged |
| Map | **118** requirements |
| Tests executed | Focused Go and Node tests, listed in §3; **no model or GPU runtime test** |
| Authorization | None for implementation, production integration, G2 or BGE |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Executive readiness

**HFR-01 (Observed — overall verdict).** Howl has significant working code, but **no subsystem in A–E has end-to-end
verification**, and the combined Howl + MEMORY-09 system is **not ready for real-workload acceptance**. This is a
readiness statement, not a defect report.

| Subsystem | Highest evidence today | Readiness | Principal gap |
|---|---|---|---|
| **A. Chat/state machine** | **E2** for history persistence and character ownership; **E0** for state machine | **Not ready** | No stable state identity, branch/alternative set, revert, lifecycle notification or non-destructive retention |
| **B. Multi-character orchestration** | **E2** for ownership/lore privacy; **E1** for queue orchestration | **Partial / unverified** | No end-to-end multi-character test; context flushing, voice separation, awareness and memory boundaries unverified |
| **C. Hardware-aware intelligence** | **E1** for resource manager/allocator/VRAM arbitration | **Partial / unverified** | No resource tests; no lease authority; no OOM/telemetry-loss recovery evidence |
| **D. Inference engine** | **E2** for props/admission/loader checks; **E1** for streaming/load | **Partial / unverified** | No real llama-server run, load/unload test, KV-cache/context matrix or model-specific acceptance |
| **E. Performance** | **E0** | **Not established** | No benchmark harness, baseline, tokens/sec, latency or measured optimization evidence |

**HFR-02 (Calculated — the strongest evidence is narrow).** A passing ownership test is not evidence that the chat state
machine works. A passing prompt-boundary test is not evidence that multi-character orchestration is ready. A model
loader test is not evidence that hardware-aware inference is ready.

---

## 2. Evidence standard

**HFR-03 (Calculated — four levels).**

| Level | Meaning | What it permits us to say |
|---|---|---|
| **E0 — Not established** | No implementation, or implementation contradicts the claimed behaviour | “We do not know” or “this is absent” |
| **E1 — Implemented** | Code exists on a live path, but no automated or controlled exercise demonstrates it | “It is written” |
| **E2 — Working in a limited test** | A narrow automated test or controlled local check passes | “This specific path works under test conditions” |
| **E3 — End-to-end verified** | A real user-relevant scenario passes across the relevant boundaries | “This scenario works” |

**HFR-04 (Calculated — no level promotion by adjacency).** A component at E2 does not raise its subsystem to E2. The
subsystem's readiness is limited by its weakest critical capability. A model loader at E2 does not make hardware-aware
inference E2. A history test at E2 does not make the chat state machine E2.

**HFR-05 (Calculated — “working” is reserved).** “Working” means an identified test or run exercised the named path.
It never means “the code looks correct.”

---

## 3. What was actually executed

**HFR-06 (Observed — commands and results, current tree, 2026-11-21).**

| Command | Result | What it exercises | What it does not prove |
|---|---|---|---|
| `go test ./internal/backend/chat` | **PASS** | History, character ownership, canonical context, parser | Branching, revert, lifecycle notification |
| `go test ./internal/bench` | **PASS** | Save/load invariants, stress rounds, concurrent saves | Performance; it is an invariant/chaos scanner, not a benchmark |
| `go test ./internal/modelintel` | **PASS** | Model profile/caching | Hardware fit or actual inference |
| `go test ./internal/backend/llama` | **PASS** | `/props` context parsing | Real server, KV cache, sampling matrix |
| `go test ./internal/resource` | **No test files** (exit 0) | Nothing | Resource manager, allocator, VRAM arbitration, OOM |
| `go test . -run 'Test(CharacterOwnership\|ReplaceChatHistory\|LegacyHistory\|ConvertFrontendHistory\|CharacterLorePrivacy\|BuildCharacterLoreBlock\|CharacterResponseBoundary\|InferenceBoundary\|BuildPromptMessages)' -count=1` | **PASS** | App history/character ownership, prompt boundary, lore privacy | Full UI/backend/model turn |
| `go test . -run 'Test(LoadModel\|LLMInfos\|InspectModel\|PreferModelPath\|ModelCatalog\|CharacterVoice)' -count=1` | **PASS** | Loader rejection, discovery, model path | Real load, unload, generation |
| `go test ./internal/backend/gguf` | **PASS** | Admission, pairing, memory estimation | Runtime fit or generation |
| `node --test tests/frontend/*.js` | **15 pass, 1 fail** | Frontend identity, lore ownership, memory ownership, character config | Multi-character orchestration; suite is not green |

**HFR-07 (Observed — the frontend failure).** `tests/frontend/character-configuration-integrity.test.js:409`
(“an explicit voice-model clear works while an unrepresentable stored voice is preserved”) fails on the current tree:
actual `custom_voice_not_in_the_fixed_select`, expected `''`. It is outside the memory work and was **not repaired by
this audit**. It means the frontend suite cannot be reported as green.

**HFR-08 (Observed — no runtime evidence).** No llama-server or other model was started; no real generation, model
load/unload, VRAM pressure, GPU telemetry, thermal event or OOM recovery was exercised. No benchmark was run.

**HFR-09 (Observed — engine connection).** `eval/mem1/m09/` is **1,951 lines** (38 + 236 + 1,209 + 115 + 252 + 101) and
has **no caller outside the package** in `eval/mem1`. The 2026-11-18 audit's integration finding still holds on the
current tree.

---

## 4. A — Chat and state machine

### 4.1 What exists

**HFR-10 (Observed).**

- `ChatHistory` with stable character ownership field; append/load/save/rewrite in
  `internal/backend/chat/persistence.go`.
- `App.appendChatMessage`, `ReplaceChatHistory`, `ClearChatHistory`, `buildChatHistory` in `app-prompt.go` /
  `app-chat.go`.
- Streaming generic and character paths in `app-prompt.go`; character-scoped events.
- Frontend chat actions: copy, edit, delete, regenerate, continue, undo/prune, prev/next scrolling.
- Chat state bridge and session save/load.
- HSL-1…HSL-11 requirement record; no state-machine implementation.

### 4.2 What has been tested

**HFR-11 (Observed — E2 for narrow history/ownership).**

- Character ownership survives append → save → restart.
- `ReplaceChatHistory` preserves character owner and metadata.
- Legacy history remains unassigned; no character ID is synthesized.
- Frontend conversion preserves character owner.
- Prompt order and lore privacy at the prompt boundary.
- These are narrow tests. None exercises a state machine.

### 4.3 What remains incomplete

**HFR-12 (Observed — E0 for the state machine).**

| HSL requirement | Current evidence |
|---|---|
| HSL-1 state ownership | No state model; one mutable list |
| HSL-2 lifecycle notification | No interface |
| HSL-3 stable state identity | Positional IDs: `len(a.chatMessages)-1`, `len(history)`, `i` |
| HSL-4 alternative selection | No alternative set; Prev/Next scrolls |
| HSL-5 regeneration | Removes DOM reply and re-streams; no alternative retention |
| HSL-6 revert | No revert action |
| HSL-7 deletion only destructive | Undo/edit prune and rewrite; previous state unrecoverable |
| HSL-8 deletion names its set | No deletion report |
| HSL-9 dependency determination | Contract deferred |
| HSL-10 retention across non-destructive transitions | Not implemented |
| HSL-11 no unilateral decisions | No boundary exists |

**HFR-13 (Observed — context budget is a placeholder).** `TrimToMaxTokens` is explicitly a placeholder in
`internal/backend/types/context.go`; `Message.Tokens` is initialised to 0 and is never populated by a tokenizer
(`types/message.go`). The trim therefore cannot enforce a real token budget. The app path also has no token-based
history trim. This is a direct blocker for the “long conversation with context flushing” scenario.

**HFR-14 (Calculated — A is the first memory dependency).** Stable state identity and source message identity are
prerequisites for the production capture boundary (`PD-M01`). Without them, memory provenance depends on list position.

### 4.4 Dependencies and independence

**HFR-15 (Calculated).**

| Question | Answer |
|---|---|
| Depends on MEMORY-09? | **No** for the state machine itself |
| Blocks MEMORY-09? | **Yes** — stable source/state identity |
| Can progress independently? | **Yes** — Howl-owned design and verification |
| Evidence that would demonstrate readiness | Branch → alternative → revert → delete scenario with stable IDs across save/load, explicit lifecycle notifications, and retained memory state; plus a real token-budget/context-flush test |

---

## 5. B — Multi-character orchestration

### 5.1 What exists

**HFR-16 (Observed).**

- `frontend/dist/js/character-queue.js`: FIFO queue, inactive/dead filtering, priority metadata sort, pause/resume,
  responded tracking.
- `streamMultiCharacterResponses` in `chat-stream.js`: loops the queue, passes prior responses as group context,
  requests a per-character lore block, calls `GenerateCharacterResponse(message, character.id)`, collects responses.
- `GenerateCharacterResponse` appends the assistant message with stable `CharacterID`.
- `buildCharacterSpecificMessage` includes character sheet, recent system events, prior responses, inactive notes,
  backend lore block and frontend `MemoryStore.getCharacterMemories` entries.
- Character-scoped TTS and streaming events.

### 5.2 What has been tested

**HFR-17 (Observed — E2 for narrow ownership/lore, not orchestration).**

- Character identity survives persistence and restart.
- Aria/Bram lore privacy: own private lore included; foreign, unowned-private and hidden lore excluded; switching
  characters does not reuse the previous block.
- Character response boundary sends the frontend-composed prompt unchanged.
- No test references `CharacterQueue`, `streamMultiCharacterResponses` or `buildCharacterSpecificMessage` in
  `tests/`; the frontend suite has no multi-character scenario.

### 5.3 What remains incomplete

**HFR-18 (Observed — E1 for orchestration).**

| Capability | State |
|---|---|
| Response sequencing | Implemented FIFO/priority loop; no automated test |
| Context flushing between characters | No real token budget; no flush test |
| Voice separation | Prompt-level instruction and per-character TTS; no contamination test |
| Shared scenes | Prior responses injected as text; no shared scene state |
| Character awareness | Prior responses and inactive notes; no knowledge-boundary test |
| Memory access boundaries | Frontend `MemoryStore` filter; no backend memory integration at this boundary |
| Abort/error recovery | Queue stops on null; no automated test |
| Concurrency | Sequential loop; no concurrent-character test |

**HFR-19 (Observed — a boundary risk).** The character prompt reads memories from the frontend `MemoryStore`, while
the backend candidate/memory service is the intended authority. This is the same bypass identified in
`PCB-03…04`: a character can receive frontend-local memories without the backend capture envelope.

### 5.4 Dependencies and independence

**HFR-20 (Calculated).**

| Question | Answer |
|---|---|
| Depends on MEMORY-09? | **No** for queue/prompt mechanics; **yes** for memory access boundaries |
| Depends on A? | **Yes** — state identity, branching, context ownership |
| Can progress independently? | **Yes** — with A; memory integration waits for `PD-M01`/`PD-M08` |
| Evidence that would demonstrate readiness | 2–3 characters, distinct voices, shared scene, private memory isolation, context flush, abort/recovery, restart; no knowledge or voice contamination |

---

## 6. C — Hardware-aware intelligence

### 6.1 What exists

**HFR-21 (Observed).**

- `internal/resource`: singleton `ResourceManager`, hardware snapshot, model registration, loaded-model queries,
  allocation advisor, priority eviction, 90% VRAM safety and 95% RAM safety.
- `app-resource.go`: snapshot, loaded models, register/unregister, allocation advice, summary.
- `app-vram-arbitration.go`: `CheckModelFit`, `FreeVRAMFor`, LLM/image estimates, live VRAM re-read after unload.
- `app-model-intelligence.go`: model catalog/profile.
- GPU detection on Windows (`nvidia-smi`, registry fallback) and Unix (`nvidia-smi`, AMD sysfs).

### 6.2 What has been tested

**HFR-22 (Observed — E1, no resource tests).**

- `go test ./internal/resource` reports **no test files**.
- `internal/modelintel` tests pass, but they exercise profile/caching, not hardware fit.
- No test covers `RequestAllocation`, `CheckModelFit`, `FreeVRAMFor`, eviction, telemetry loss, OOM or release timing.

### 6.3 What remains incomplete

**HFR-23 (Observed).**

| Capability | State |
|---|---|
| Hardware inventory | Implemented with multiple providers and fallbacks; no confidence/source field |
| One lease authority | **Absent** — allocation is advice; callers can ignore it |
| Reservation semantics | A 10% margin exists (`vramSafetyPct=0.90`); no measured reserve/overhead model |
| OOM prevention | Estimates and eviction advice; no admission enforcement or OOM test |
| OOM recovery | Not established |
| Telemetry loss | Last-known-good fallback exists; no fault-injection test |
| Async release | 1.5 s sleep then re-read; no test that the driver has reclaimed memory |
| Multi-GPU | Summed VRAM; no multi-GPU mode/claim |
| Thermal awareness | Not established in this audit |

**HFR-24 (Calculated — the “working model loader” trap).** A successful model load would prove that one model loaded
on one machine. It would not prove an authoritative resource system, safe admission, telemetry confidence or recovery
under pressure.

### 6.4 Dependencies and independence

**HFR-25 (Calculated).**

| Question | Answer |
|---|---|
| Depends on MEMORY-09? | **No** |
| Depends on D? | **Yes** — model footprint and context/KV behavior |
| Can progress independently? | **Yes** — resource tests and fault injection |
| Evidence that would demonstrate readiness | Measured telemetry under load; admission refusal/eviction; telemetry-loss behavior; OOM recovery; verified release; known model-fit report with confidence |

---

## 7. D — Inference engine and Kobold-like capabilities

### 7.1 What exists

**HFR-26 (Observed).**

- `internal/backend/llama/options.go`: context size, GPU layers, batch size, threads, sampling controls (temperature,
  top-p, top-k, repeat penalty, TFS, typical-p, mirostat, top-a), stop strings.
- `http_client.go`: OpenAI-compatible requests, stop strings, `/props` `n_ctx` parsing, streaming callbacks.
- `internal/backend/gguf`: admission, pairing, quantization metadata, memory/KV-cache estimation, usable-context
  calculation.
- Loader integrity checks and model discovery.
- `app-model`/`app-loader` paths for load rejection and registration.

### 7.2 What has been tested

**HFR-27 (Observed — E2 for narrow parsing/admission).**

- `TestGetModelInfoReadsProps`, `TestGetContextSizeFallsBackToTopLevel` pass.
- GGUF admission/pairing/memory tests pass.
- Loader tests reject vocoders, unknown architectures, task-ineligible models and component binaries before spawning.
- Model discovery/path tests pass.

These tests do not start llama-server, stream a real response, load/unload a model or measure KV-cache behavior.

### 7.3 What remains incomplete

**HFR-28 (Observed).**

| Capability | State |
|---|---|
| KV cache management | Options/estimation only; no cache lifecycle or reuse test |
| Context sizing | `/props` read; no model-specific behavior test |
| Sampling controls | Fields and request serialization; no matrix test |
| Stop sequences | Serialized; no end-to-end stop test |
| Batching | Batch options exist; no batching execution path verified |
| Speculative decoding | Not found in inspected code |
| Quantization | Metadata/registry support; no runtime capability matrix |
| Model-specific optimizations | Not consolidated |
| Load/unload | Code paths exist; no end-to-end test |
| GPU offload | Options and estimates; no measured placement test |

**HFR-29 (Calculated — D is not “Kobold-like” yet).** The existence of sampling fields and an OpenAI-compatible client
is not a capability claim. D readiness requires a named model/backend/hardware matrix with observed behavior.

### 7.4 Dependencies and independence

**HFR-30 (Calculated).**

| Question | Answer |
|---|---|
| Depends on MEMORY-09? | **No** |
| Depends on C? | **Yes** — load/fit/telemetry |
| Can progress independently? | **Yes** — with C |
| Evidence that would demonstrate readiness | Real server run; context sizes; sampling/stop matrix; load/unload/reload; KV-cache reuse; quantization and backend support; named model matrix |

---

## 8. E — Performance and hardware tricks

### 8.1 What exists

**HFR-31 (Observed).**

- Some measured constants/notes in code (image VRAM floor, overhead, 90% margin).
- `internal/bench` is an invariant/stress scanner, not a performance benchmark.
- No `func Benchmark` exists in the repository.
- No reproducible harness for tokens/sec, latency, VRAM peak, load time or optimization comparisons.

### 8.2 What has been tested

**HFR-32 (Observed — E0).** No benchmark was run. `go test ./internal/bench` passing says the persistence invariants
held; it says nothing about performance.

### 8.3 What remains incomplete

**HFR-33 (Observed).**

| Area | State |
|---|---|
| CPU/GPU layer placement | Estimates; no measured comparison |
| Cache reuse | Not established |
| Memory mapping | Not established |
| Flash attention | Option/flag present; no measured gain |
| Async loading | Not established |
| Parallelism | Not established |
| Model-specific acceleration | Not established |
| Benchmark discipline | No harness; no baseline; no “measure, do not assume” gate |

**HFR-34 (Calculated — E cannot be claimed from constants).** A recorded constant is not a benchmark. Performance
readiness requires a repeatable workload, named hardware/model, baseline, metric and variance.

### 8.4 Dependencies and independence

**HFR-35 (Calculated).**

| Question | Answer |
|---|---|
| Depends on MEMORY-09? | **No** |
| Depends on C/D? | **Yes** — stable load/run behavior |
| Can progress independently? | **Measurement only**, once C/D can run a controlled model |
| Evidence that would demonstrate readiness | Reproducible benchmark set with baseline and measured gains; no optimization accepted from intuition |

---

## 9. Cross-cutting blockers

**HFR-36 (Observed — context flushing is not real).** `TrimToMaxTokens` is a placeholder and `Message.Tokens` is never
populated. Long-conversation memory cannot be accepted while the context budget is unenforced.

**HFR-37 (Observed — source identity is positional).** `appendChatMessage` uses `len(a.chatMessages)-1`;
`buildChatHistory` uses `len(history)`; `ReplaceChatHistory` uses `i`. A branch/revert/delete changes identities.

**HFR-38 (Observed — memory authority is split).** Frontend `MemoryStore`/localStorage is used for character memories;
backend candidate/memory service is a separate path with no App bindings. This is the production capture break in
`PCB-01…04`.

**HFR-39 (Observed — no combined harness).** No test drives UI → App → chat persistence → memory candidate → retrieval
under one scenario. The `m09` engine has no caller.

**HFR-40 (Observed — no real runtime evidence).** No model, GPU pressure, telemetry loss, restart-under-load or
multi-character real generation was exercised.

**HFR-41 (Observed — frontend suite not green).** One frontend character-configuration test fails on the current tree
(HFR-07). It is unrelated to MEMORY-09 but is part of foundation test health.

---

## 10. Integration scenario readiness

**HFR-42 (Calculated — none of the nine scenarios is end-to-end verified).**

| Scenario | Current readiness | Blocking prerequisite |
|---|---|---|
| **Character switches** | Partial: character IDs and lore privacy at E2; no memory activation/isolation | A state identity; B orchestration; M1/M8 |
| **World/scenario switches** | Partial: context persistence and stable IDs exist; no memory selection/restoration test | A context ownership; M1/M5/M8 |
| **Character forgets something** | Not established: no lifecycle event path | M3/M4/M5 |
| **Character learns something** | Not established: no production capture/acquisition path | M1/M4 |
| **Long conversation** | Not established: token budget is a placeholder | A real context budget; M5 |
| **Model unload/reload** | Not established: load/unload untested; memory independence untested | C/D; M7 |
| **Multiple characters** | Partial: queue and lore privacy exist; no contamination/E2E test | A/B; M1/M8 |
| **Hardware pressure** | Not established: no resource tests or OOM/telemetry evidence | C/E |
| **Rebuild/restart** | Partial: history persistence tested; memory reconstruction not integrated | A; M5/M7/M8 |

**HFR-43 (Calculated — the integration gate is evidence, not a meeting).** A scenario is ready to test only when its
prerequisite seam exists. A scenario is passed only when the named end-to-end evidence exists.

---

## 11. MEMORY-09 decision register — reconciled

**HFR-44 (Calculated — status reconciliation only; no decision changed).** The register preserves `PD-M01…PD-M10` from
the program direction and adds the readiness evidence.

| ID | Decision | Current position | Evidence from this audit | Status |
|---|---|---|---|---|
| **PD-M01** | Production capture boundary | Server-side capture at append recommended | Stable IDs absent; frontend bypass confirmed; proposal exists | **Open** |
| **PD-M02** | Autonomous classification | Rules-first; hybrid/BGE candidate only | G2 blocked; no corpus; no class experiment | **Blocked** |
| **PD-M03** | Unresolved-class lifecycle | Preserve uncertainty; no class-sensitive use | L-2/Decision A amendment not authorized | **Blocked** |
| **PD-M04** | Acquisition/lifecycle events | A2 ordering selected; carrier incomplete | No production event path; HSL contract deferred | **Partially decided** |
| **PD-M05** | Ordering/authority/provenance/reconstruction/integrity | 57 rules / 118 requirements | No runtime evidence; conformance mapping stale for current tree | **Open** |
| **PD-M06** | Conformance coverage | 118 defined; 0 executed | Map 118; G1 fixture only; m09 unconnected | **Open** |
| **PD-M07** | Engine ↔ harness integration | No caller from harness to m09 | 1,951-line package; no external caller observed | **Open** |
| **PD-M08** | Howl ↔ MEMORY-09 contract | HSL §7 deferred | HSL requirements recorded; no state machine | **Open** |
| **PD-M09** | BGE embedding experiment | Separately gated | No fetch; reranker is not an embedding model | **Deferred** |
| **PD-M10** | Real-workload acceptance | Not defined | A–E prerequisites unmet | **Deferred** |

**HFR-45 (Observed — settled decisions remain settled).** L-2, Decision A, G0, the G1 fixture contract and P-1/S-3 are
not reopened. This audit only changes readiness evidence, not the decisions.

---

## 12. What can progress independently

**HFR-46 (Calculated — independent work is allowed; integration claims are not).**

| Workstream | Can proceed now? | Next evidence work | Depends on |
|---|---|---|---|
| **A. State machine** | **Yes** | State identity, branch/alternative, revert, notification, retention test design | None from MEMORY-09 |
| **B. Multi-character** | **Yes** | Queue/orchestration tests, contamination harness, context flush | A |
| **C. Hardware** | **Yes** | Resource tests, telemetry-loss/OOM fault injection, lease design | D for footprint |
| **D. Inference** | **Yes** | Real server matrix, load/unload, KV/context/sampling tests | C |
| **E. Performance** | **Measurement only** | Benchmark harness and baseline | C/D |
| **M1 capture** | **Design only** | Interface design against stable IDs/state identity | A |
| **M2 classification** | **No** | — | M1, corpus, G2 |
| **M3 unresolved lifecycle** | **Analysis only** | Decision analysis; no class-sensitive use | Owner/L-2/Decision A |
| **M4/M5 events** | **Research/design** | Event/ordering closure analysis | M1, M8 |
| **M8 contract** | **Draft after A** | State/transition/deletion contract | HSL |

**HFR-47 (Calculated — no global barrier).** A, B, C, D and E can move without MEMORY-09. M1 design can move without
Howl code, but its implementation waits for stable source/state identity. Integration waits for both sides.

---

## 13. Checkpoint `HFR-CP-2026-11-21-01`

**HFR-48 (Observed — resume here).**

| Field | State |
|---|---|
| Audit status | **Complete** |
| Overall readiness | **Not ready for integrated real-workload acceptance** |
| A. Chat/state | **Not ready** — history at E2; state machine at E0 |
| B. Multi-character | **Partial / unverified** — ownership/lore at E2; orchestration at E1 |
| C. Hardware | **Partial / unverified** — resource code at E1; no tests |
| D. Inference | **Partial / unverified** — parsing/admission at E2; end-to-end at E0/E1 |
| E. Performance | **Not established** — E0 |
| Tests run | Focused Go tests pass; `internal/resource` has no tests; Node 15/16 pass, 1 fail |
| Production changes | **0** |
| MEMORY-09 decisions changed | **0** |
| Authorization | No implementation, G2, BGE, production API or corpus |

### 13.1 Recommended next evidence work

**HFR-49 (Calculated — highest leverage first).**

1. **A state identity + context budget:** define stable state/message IDs and a real token-budget path. This unblocks
   M1 and the long-conversation scenario.
2. **C resource tests:** fault-inject telemetry loss, low VRAM, eviction, release and OOM; decide whether allocation is
   advisory or authoritative.
3. **D real-server smoke matrix:** load, stream, stop, unload, reload, context-size read; no performance claims yet.
4. **B orchestration harness:** two/three characters, private lore/memory leakage, abort/recovery.
5. **E benchmark harness:** only after C/D can run a controlled model.

### 13.2 Invariants

**HFR-50 (Observed).**

- No E3 claim without an end-to-end scenario.
- No E2 claim from a code path alone.
- No class from incomplete capture.
- No fabricated knowledge to fill a missing field.
- No `told` inference from character speech.
- No authored-world eligibility from subject matter.
- No memory state dependent on inference state.
- No performance claim without measurement.
- No reopening L-2, Decision A or settled owner decisions.

### 13.3 Resume instruction

**HFR-51 (Calculated).** On return:

1. Read the program direction §2 and §8, then this audit §1, §10 and §13.
2. Re-run the focused tests in §3; do not assume the results still hold.
3. Check whether A's state identity or C's resource tests have moved; if so, update this register with new evidence.
4. Do not start integration verification until the specific seam exists.
5. If a decision is needed, record it; do not infer it from a convenient subsystem.

---

## 14. Limits

**HFR-52 (Observed).**

- This is a **static audit plus focused tests**, not a full code review.
- No real model, GPU, llama-server, OOM, thermal or restart-under-load test was run.
- The `m09` line count and caller check were re-verified; other 2026-11-18 figures are dated.
- `go test` was run on selected packages, not the whole repository.
- The frontend suite is not green (HFR-07).
- No UI was rendered; no browser interaction was performed.
- No production code, schema, API or canonical document was changed.

---

## 15. Record status

**HFR-53 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Production code/schema edited | **0** |
| Requirements changed | **0** — map **118** |
| L-2 / Decision A amended | **0** |
| G2 / BGE started | **0** |
| Models fetched / corpus created | **0** |
| Tests executed | Focused Go + Node tests, §3 |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **readiness audit**, not canonical capture |

**HFR-54 (Observed — evidence paths).**

| Evidence | Path |
|---|---|
| History/state append and rewrite | [`app-prompt.go`](../../app-prompt.go) · [`app-chat.go`](../../app-chat.go) · [`internal/backend/chat/types.go`](../../internal/backend/chat/types.go) |
| Frontend destructive actions and regenerate | [`frontend/dist/js/chat-actions.js`](../../frontend/dist/js/chat-actions.js) |
| Multi-character queue and orchestration | [`frontend/dist/js/character-queue.js`](../../frontend/dist/js/character-queue.js) · [`frontend/dist/js/chat-stream.js`](../../frontend/dist/js/chat-stream.js) |
| Resource manager and VRAM arbitration | [`internal/resource/`](../../internal/resource/) · [`app-resource.go`](../../app-resource.go) · [`app-vram-arbitration.go`](../../app-vram-arbitration.go) |
| Inference options and GGUF estimates | [`internal/backend/llama/`](../../internal/backend/llama/) · [`internal/backend/gguf/`](../../internal/backend/gguf/) |
| Context-budget placeholder | [`internal/backend/types/context.go`](../../internal/backend/types/context.go) |
| MEMORY-09 engine | [`eval/mem1/m09/`](../../eval/mem1/m09/) |
| G1 fixture harness | [`internal/backend/memory/capture_contract_test.go`](../../internal/backend/memory/capture_contract_test.go) |

**End of audit.** A is not ready; B–D are partial and unverified; E is unestablished. The strongest existing evidence is
narrow history/ownership/prompt-boundary testing (E2). No subsystem has E3 evidence, and the combined system has no
end-to-end verification (HFR-01…HFR-08). MEMORY-09 readiness is unchanged: capture boundary open, G2 blocked,
unresolved lifecycle blocked, engine unconnected (HFR-44). Next evidence work is state identity + context budget,
resource fault injection, and a real-server smoke matrix (HFR-49). **Resume at `HFR-CP-2026-11-21-01`.**
