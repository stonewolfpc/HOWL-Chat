# HOWL CHAT — CORRECTED MASTER PLAN

**Document ID:** `HOWL-CORRECTED-MASTER-PLAN-2026-11-11`  
**Date:** 2026-11-11  
**Status:** **PROPOSED CORRECTED MASTER PLAN — REVISION 3; NOT ADOPTED; NO IMPLEMENTATION AUTHORIZATION**  
**Revision:** 3 — MUSE-01–MUSE-12 correction, N1–N7 consistency closure, and O-1–O-3 documentary cleanup. See `docs/review/MUSE-CORRECTION-DISPOSITION-HOWL-2026-11-11.md`.  
**Revision identity:** Revision 2 recorded the MUSE correction; Revision 3 is the separate N1–N7 closure identity. Technical scope, work packages, and dependencies are frozen except for those consistency corrections.  
**Phase:** Architecture and Requirements Reconciliation  
**Supersedes:** No prior plan. Consumes, without replacing, the existing master plan, owner packet, adjudication register, recovery addendum, and requirements reconciliation.  
**Authority:** `docs/review/REQUIREMENTS-RECOVERY-OWNER-ADDENDUM-HOWL-2026-11-11.md` and `docs/review/REQUIREMENTS-RECONCILIATION-HOWL-2026-11-11.md`.

> This is a planning artifact. It does not authorize implementation, model loading, benchmarking, network experiments, capability enablement, schema migration, or acceptance. No work package in this document is accepted, closed, or funded by appearing here.

## Source verification

The following source identities were verified before this plan was written. The plan relies on these versions, not on older or similarly named copies.

| Source | SHA-256 | Size / lines |
|---|---|---|
| Recovery addendum `REQUIREMENTS-RECOVERY-OWNER-ADDENDUM-HOWL-2026-11-11.md` | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` | 222,234 bytes / 2,418 lines |
| Reconciliation `REQUIREMENTS-RECONCILIATION-HOWL-2026-11-11.md` | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` | 39,612 bytes / 421 lines |
| Existing master plan `HOWL-COMPLETE-PRODUCT-ARCHITECTURE-REQUIREMENTS-DEPENDENCY-MASTER-PLAN-2026-11-11.md` | `e39cc396fa20adf2656a1b4e650a8bdc6c9a49cb46c3ddf95309ff144229ef8c` | 83,354 bytes / 703 lines |
| Owner decision packet `OWNER-DECISION-PACKET-HOWL-MASTER-PLAN-2026-11-11.md` | `6721d770694cda2c7b132fdc15bbd8d19dd6de0502c4f3e5f3ca47ce620ff4f3` | 29,376 bytes / 385 lines |
| Adjudication matrix `ADJUDICATION-MATRIX-HOWL-MASTER-PLAN-MUSE-REVIEW-2026-11-11.md` | `a9fdaeefdc4ea57c28a706cc09e57a77566f3ed6150d2e4e305904de7b6049ed` | 46,357 bytes / 458 lines |
| Muse review `REVIEW-INDEPENDENT-HOWL-COMPLETE-PRODUCT-ARCHITECTURE-MASTER-PLAN-2026-11-11.md` | `af493c0b6118ebede556d4927e8ab93b9378ce46c8298f3bb9e38554f6265e9d` | 33,410 bytes / 250 lines |
| MUSE corrected-plan review `REVIEW-MUSE-CORRECTED-MASTER-PLAN-2026-11-11.md` | `56aa7b715a6d962a2ad2f08db04725c4c9fb328c8275daf827b8b160a80ef942` | 34,301 bytes / 328 lines |
| `STAGE_STATUS.md` | `33a16fdbca509fc4407bf840deee1dd622a0fe5a1bd272ac4c1ca1179dbe9a62` | Current stage record |
| `PROJECT_ISSUES_AND_EXPANSION.md` | `8156318ec8855ddbab1f218927352e59de6d60e23a3c81f11f85a56c61213491` | Requirements ledger |
| `docs/tier4-avatar.md` | `be04f0253e61c0fc28c245edbd8eb7f3ad753b8bdf152ceb8308cace23be04c9` | Historical avatar note |
| F2 model intelligence report `ENGINEERING-F2-STORAGE-MODEL-INTELLIGENCE-2026-11-09.md` | `0f40221f544d04365b784f103d6befc4230b7e2f1699f81d9002265169ac76d1` | Model intelligence foundation |
| F2-B integration report `ENGINEERING-F2B-STORAGE-MODEL-INTEGRATION-2026-11-09.md` | `95e24ccdcdc033c8b050b0f435df8b3eee85dc2e2dd6dab7c3e0bd6d451bee4a` | Browser integration |
| Character voice investigation `docs/review/INVESTIGATION-character-voice-architecture-2026-11-05.md` | `0236e2300ef4cd704a946035bce5c943a407e33d3bf589214dc6927a8b85ea6e` | Voice evidence |

The protected F3-C artifacts remain unchanged under their five recorded hashes:

- `85e1968bac4eb70f489bbd7dc89057debd17e7b3bc2336aa69c036a4352f83e8`
- `ed8bf2e2326a2ec74838dc0c8cd4de3332b8807d2d0a2647a18cfda060bc9c7c`
- `343104b528818df699d752c419770fdc1bfd7fe71bc1f249e0dc6847f308c6b8`
- `494539b9dde27a78c9c4c091131a93dce99c1debd81d11fea911e6d411e2e6a6`
- `ce32ac0bedbe216b55c065be073b2165209e0ebad4cb435897b829336d29d2c5`

---

## 1. Product definition and owner-approved boundaries

Howl Chat is a **single-user, local-first Windows desktop creative environment** for character-based conversation, local text inference, local image and video generation, local speech, and a bounded set of explicitly authorized external actions. The product's differentiator is not a hosted model service. It is the combination of user-owned models, durable character and conversation state, local creative tools, and transparent control over what leaves the machine.

The product definition is constrained by the owner-characterized recovery baseline:

| Boundary | Owner-characterized position | Planning consequence |
|---|---|---|
| Desktop character avatars | **Excluded** (`ODR-01`, `REC-C-06`) | Not a roadmap item. Historical tier-4 notes are evidence only. Chat portraits remain required. |
| General local creative studio | **Approved direction** (`ODR-03`) | The direction is accepted; no individual modality is authorized by that approval. |
| Model hot-swap during active conversation | **Required product behavior** (`ODR-02`) | Implementation remains partial and unvalidated; product acceptance is not implementation acceptance. |
| Durable conversational continuity and memory | **Required** (`ODR-04`, `REC-CH-09`, `REC-CH-10`) | Memory, rollover, and character isolation remain separate responsibilities. |
| Network boundary | **First-class requirement family** (`ODR-12`, `NET-01`–`NET-20`) | Local-first is a required property, not an observed property. The current implementation contradicts it in inspected paths. |
| Android remote | **Conditional, unimplemented** (`UD-02`, `NET-12`) | No remote server, pairing, authentication, or transport exists. Off by default if ever authorized. |
| Local-Only Mode | **Renderer-only today** | It is not an application-wide network boundary and cannot be cited as offline enforcement. |
| Cloud model inference | **Not part of the product** (`NET-02`) | No cloud model provider is an approved fallback for local inference. |
| Requirements reconciliation and capability evidence | **Independent registers** | A requirement can be missing while a related control exists. A control can be implemented while the requirement remains partial. |

**Plan rules.** This plan does not promote a requirement because code exists, a control is visible, a setting is stored, an API call succeeds, or a test name exists. It does not demote an owner decision because the implementation is incomplete. It preserves the progression:

**Research → Evidence → Compatibility → Selection → Validation → Authorization**

That progression is not a conveyor belt. Any stage may stop, reject, return an item for investigation, or require a new owner decision.

## 2. Authoritative requirements and exclusions

### 2.1 Requirement inventory

The corrected plan consumes two requirement families without merging them:

| Family | Count | Source | Treatment |
|---|---:|---|---|
| Original reconciliation rows | **66** | Reconciliation §4 | Every row remains individually traceable. No row is dropped, merged, or silently re-scoped. |
| Network requirements | **20** | Addendum §23; reconciliation §4.7 | `NET-01`–`NET-20` remain a separate first-class family. |
| Existing-capability register | **37 rows** | Reconciliation §5.2 | Reproduced in §4. It remains evidence about existing controls, not a requirement-status table. |
| Unresolved owner questions | **6** | `UD-01`–`UD-06` | Preserved as unresolved. No answer is invented. |
| Owner decisions | **11 matrix rows plus ODR-12** | `ODR-01`–`ODR-11` in the 66; `ODR-12` in the addendum | All preserved. `ODR-12` is controlling governance, not an extra matrix row. |
| Network verification tests | **28** | Addendum §23 (`V-NET-01`–`V-NET-28`) | Individually mapped in §10.6; no test is marked passed. |
| Network path identities | **27** | Addendum §23 (`N-01`–`N-14`, `L-01`–`L-05`, `P-01`–`P-08`) | Individually listed in §10.3–§10.5; `N-14` is display-only, not egress. |

The original 66 rows comprise 49 `REC-*` rows, 11 `ODR-*` rows, and 6 `UD-*` rows. The apparent `ODR-12` gap in the matrix is intentional: `ODR-12` is an addendum-level governance decision and is preserved separately in §7 and §12. It is not fabricated as a 67th matrix row.

### 2.2 Requirement outcomes

The six reconciliation outcomes are preserved exactly:

- **Satisfied**
- **Partially satisfied**
- **Missing**
- **Contradicted**
- **Unresolved**
- **Excluded**

An outcome is not a capability state. “Missing” does not mean no related code exists. “Satisfied” does not mean the implementation was runtime-validated. “Contradicted” means the requirement conflicts with current evidence or an owner decision, not that the requirement has been withdrawn.

### 2.3 Exclusions and parked items

| Item | Status | Planning treatment |
|---|---|---|
| Desktop reactive 3D avatar / multi-character desktop avatars | **Excluded** | No work package. Historical note remains evidence only. |
| `REC-C-02` custom workspace layout | **Satisfied (parked disposition)** | No promotion by the studio direction. No work package. |
| `REC-C-03` image staging-board structure | **Satisfied (parked disposition)** | No promotion by the studio direction. No work package. |
| Native TTS/F2-B7 expansion | **Frozen** | Not reopened by this plan. Any future change requires the governing decision. |
| F3-C acceptance | **Not granted** | No plan section implies acceptance. Release gates remain blocked where the protected record requires it. |
| Cloud model inference/provider support | **Not approved** | `NET-02` remains contradicted by inspected remote ASR/TTS paths until separately resolved. |
| Automatic update checking | **Absent; must remain absent** | `NET-20` is satisfied by absence only. Any future check requires explicit signed design and authorization. |

### 2.4 Non-negotiable preservation rules

1. A requirement outcome and a capability-evidence state remain separate columns and separate reasoning.
2. A work package may address part of a requirement without marking the requirement satisfied.
3. A prerequisite that is not evidenced remains a prerequisite. It is not assumed complete.
4. An owner decision is not inferred from an implementation, a control label, or a test name.
5. A capability is not promoted to **Operational**, **Validated**, **Benchmarked**, or **Accepted** without the corresponding evidence.
6. No implementation estimate is a commitment. This plan contains no duration, staffing, or effort commitments.

---

## 3. Current architecture and responsibility ownership

### 3.1 Architectural baseline

The current application is a Wails desktop shell with a Go backend and static frontend assets under `frontend/dist`. The backend owns process supervision, local model servers, chat transport, storage, and most domain services. The renderer owns presentation, some browser-local state, and parts of import/editing workflows. The architecture is not a single clean service boundary: several workflows have more than one possible owner, and some “state” exists in both browser storage and native files.

This plan does not redesign the application by preference. It records the current ownership baseline and identifies where ownership must be decided before implementation can be ordered safely.

### 3.2 Complete subsystem inventory

The following table preserves the existing `S1`–`S24` inventory. Status is conservative: **source-confirmed** means inspected structure exists, not that it operates correctly.

| ID | Subsystem | Current backend/frontend ownership | Baseline status |
|---|---|---|---|
| S1 | Shell, lifecycle, process supervision | `main.go`, `app.go`, lifecycle/process files, llama/image/audio supervisors; Wails window | Source-confirmed, partial. Shutdown admission and owned-process model exist. Packaged lifecycle acceptance remains unknown. |
| S2 | Chat, conversation state, streaming, prompt assembly | `app-chat.go`, `app-prompt.go`, llama HTTP client, `internal/backend/chat/*`; chat scripts | Source-confirmed, incomplete. Live path streams through `App`/llama; `chat.Service` is separate. No per-turn provenance or per-chat model. |
| S3 | Character management | Character handlers/store/types; browser character UI and IndexedDB plus native `Characters/<name>/character.json` fallback | Source-confirmed, implemented/unverified. Dual ownership remains. F3-D is not authorized. |
| S4 | Worlds and scenarios | World store, content import, library/character handlers; world/scenario browser editors | Mixed. Store foundation is test-confirmed for its own tests; browser editors still use import sources; F3-C Phase 0 is open. |
| S5 | Lorebooks, triggers, conflict resolution | Lorebook store/backend; lorebook UI and cache | Source- and test-confirmed with documented limitations. F3-B is closed/accepted only within that scope. Prompt synchronization with a running model remains unverified. |
| S6 | Memory and continuity | Memory backend and handlers; memory store/candidate/review UI | Source-confirmed, incomplete. Manual/automated memory, pinning, decay, and candidate review exist. Structured inventory and application-owned continuity are planned. |
| S7 | RAG / document knowledge | RAG backend and upload handler; `chat.Service` methods; knowledge UI | Source-confirmed, integration gap. The live prompt path does not demonstrably call RAG. Presence is not proof that replies use it. |
| S8 | Image generation and editing | SD.cpp engine, job/gallery/metadata packages, image handlers; generator/workflow UI | Source-confirmed, implemented/unverified. Many image paths exist; runtime acceptance remains unresolved. |
| S9 | Video generation and processing | Video engine and job manager; video modes in image UI | Source-confirmed, incomplete. No complete gallery sidecar/catalogue/delete contract comparable to images. |
| S10 | ASR / speech recognition | Whisper/Qwen/llama adapters, audio runtime, audio handlers; audio settings/chat microphone | Source-confirmed, partial verification. Remote engine availability is explicitly unverified; remote paths create network-policy work. |
| S11 | TTS / voice synthesis | Piper/Kokoro/MeloTTS/helper/cache, audio handlers, catalog; voice controls | Source-confirmed, governance unresolved. F2-B7 remains frozen. No accepted end-to-end claim. |
| S12 | Audio generation / music | No production generator found | Planned only. Not implemented. |
| S13 | Vision and multimodal input | llama HTTP client, prompt image path, vision support, projector pairing; chat input | Source-confirmed, current-turn only. Multi-image context, persistence, and model-aware anchoring are planned/unknown. |
| S14 | Model discovery, registry, intelligence | Model-intelligence, GGUF, image registry, model handlers; model directory/settings UI | Source-confirmed, implemented/unverified. Evidence-based profiles and template inspection exist; filename fallback and compatibility gaps remain. |
| S15 | Runtime model loading and hot-swap | Model manager, llama lifecycle, model handlers; model browser/settings | Source-confirmed, partial. Load does not coordinate an active stream; failure after old-server stop has no rollback. Hot-swap is required product behavior but unvalidated. |
| S16 | Hardware monitoring and resource orchestration | Resource manager, VRAM arbitration, hardware/audio monitors; VRAM dialogs/settings | Source-confirmed, incomplete. Admission is advisory; image registration is inconsistent; TTS is outside admission; release/reclamation is unproven. |
| S17 | Content libraries, import/export, ownership | Library/content-import/path/wipe handlers; content-library and import UI | Source-confirmed, partial. Registry roles and staged wipe exist; not all content types are cut over; export is category-limited. |
| S18 | Settings and configuration | Unified/image settings, reset/family defaults; settings scripts | Source-confirmed, partial. Mixed atomicity and mixed consumption. Audio settings restore without activation by design. |
| S19 | UI navigation, theme, tutorials | Static frontend assets and shared scripts | Source-confirmed, incomplete. Earlier visual/ghost-ID/duplicate-element findings remain relevant unless separately closed. |
| S20 | Story tools and session utilities | Upload/memory/export bindings, no cohesive owner; story-tools page/scripts | Unknown. Page exists; committed scope is not established. |
| S21 | Errors, logging, diagnostics, recovery | Error/logging packages, loader failure records, image diagnostics; error UI | Source-confirmed, partial. Structured reporting exists; many paths still log/ignore rather than surface one recovery contract. Error-report payload and redaction are not proven. |
| S22 | Desktop avatars / 3D characters | No production 3D implementation; image/TTS dependencies; 2D chat portraits only | **Excluded desktop surface.** Historical notes are not roadmap commitments. |
| S23 | Android remote access / QR connection | No server, pairing, transport, or auth found | Unimplemented. Discovery target only; `NET-12`/`UD-02` control any future work. |
| S24 | Creative studio / general asset generation | No dedicated subsystem; separate image/video/audio surfaces | Approved direction, not an implemented unified system. Each modality remains separately gated. |

### 3.3 Ownership findings that constrain sequencing

1. **Chat has two possible prompt paths.** The live `App`/llama stream and `chat.Service` coexist. Any memory, lore, RAG, rollover, or token-budget work is premature until one live prompt path is authoritative.
2. **Chat and memory roots are split.** The F3-C root defect remains open. Durable continuity cannot be designed on an unproven identity boundary.
3. **“Per-chat” is not a settled schema.** `UD-01` asks whether model selection is global or per-chat. Author's Note, attachments, and target sets depend on the same durable chat identity.
4. **Media identity is incomplete.** `ChatHistory` records `hasImage`; display media, model-input media, and saved media are not one durable lifecycle.
5. **Model selection is global.** Hot-swap exists as manual stop/spawn behavior, not as a conversation-safe lifecycle with readiness and rollback.
6. **Resource admission is advisory.** Image, chat, and TTS do not share one enforced model/resource state machine.
7. **Network enforcement is not centralized.** The observed Local-Only Mode is renderer-only; Go clients, sidecars, and Wails-mediated open actions are not covered.
8. **Browser storage and native files remain dual owners.** Character, world, scenario, lore, and chat state must not be treated as one canonical store until cutover evidence exists.
9. **F3-C is frozen, not accepted.** The architecture freeze preserves a boundary; it does not close root, wipe, identity, reach, startup, or alias defects.
10. **The capability catalogue is not a compatibility proof.** Model intelligence provides profiles and template evidence, but no target matrix demonstrates that a technique activates or helps on a given model/backend/hardware combination.

### 3.4 Protected boundaries

| Boundary | Current position | Plan treatment |
|---|---|---|
| F3-C Phase 0 | **OPEN** | Blocks durable chat schema and root-dependent work. |
| R6 | **UNAUTHORIZED** | Not used or inferred as authorization. |
| F3-C acceptance | **NOT GRANTED** | No gate may report acceptance from this plan. |
| F3-D | **BLOCKED** | Character/workspace work cannot assume it is available. |
| D-ROOT | **OPEN** (six defects) | Root identity is a hard prerequisite for schema and continuity. |
| D-STORAGE | **BLOCKED** | Storage work remains blocked. |
| V20 | **BLOCKED pending F3-F** | Not reopened. |
| Gate v2 V1–V17 | **ADOPTED** | Preserved. |
| Gate v2 V18–V24 | **PROPOSED — NOT ADOPTED** | No adoption inferred. |
| F3-B | **Closed/accepted with documented limitations** | Does not establish product-wide lore correctness. |
| W8 | **Verified** | Preserve as verified; do not reopen without cause. |
| W4 | **Historical/superseded; formal closure gap** | Record as unresolved evidence gap, not an active implementation mission. |

---

## 4. Existing implementation and capability evidence

### 4.1 Evidence model

The seven capability states are independent and ordered only for evidence review, not for automatic promotion:

| State | Required evidence | Not sufficient |
|---|---|---|
| **Implemented** | A control, field, binding, or code path exists | A label, a menu item, or a stored setting alone |
| **Wired** | The control reaches the intended runtime path | A setting being saved; a conditional branch existing |
| **Supported** | The selected backend, model architecture, quantization, and hardware can accept it | A generic option being exposed |
| **Operational** | It actually activates without silent fallback in the target environment | A process starting; a request returning HTTP 200 |
| **Validated** | Behavior is correct against defined expected outcomes | One successful run without a comparison |
| **Benchmarked** | A measured benefit and cost exist for a declared scope | Anecdote, advertised speed, or a single timing |
| **Accepted** | The owner accepts inclusion/use for that scope | The implementation merely existing; F3-C acceptance |

A capability row may be `Implemented = Yes`, `Wired = Partial`, and every later state `Unknown`. That is a valid and common result. It is not a missing-feature classification.

### 4.2 Flash Attention: reference discipline

Flash Attention is the reference example because its checkbox is not evidence of activation.

| Question | Current answer | Why |
|---|---|---|
| Does the control exist? | **Yes** | A setting/option exists. |
| Does configuration reach the runtime? | **Partial** | It reaches a server-argument path conditionally; end-to-end setting-to-runtime proof is absent. |
| Does the model/backend/hardware support it? | **Unknown** | Architecture, quantization, backend build, and hardware dependencies are not resolved. |
| Does it activate? | **Unknown** | No runtime evidence that the backend enabled it rather than falling back. |
| Does it operate correctly? | **Unknown** | No behavioral or numerical validation on target combinations. |
| Does it produce a measurable benefit? | **No benchmark** | No measured speed/memory gain or quality trade-off. |
| Is the benefit worth retaining? | **Unresolved** | No evidence basis for retention beyond control presence. |

The same seven questions apply to GPU offload, tensor split, KV-cache choices, context shift, cache reuse, quantization, multi-GPU, samplers, and every future technique. Investigation work packages exist so that an unverified control is not automatically reclassified as a missing feature or as a completed capability.

### 4.3 Complete existing-capability verification register

The register below is reproduced from reconciliation §5.2. It preserves the seven-state vector `I/W/S/O/V/B/A` and its stated scope. `Y` yes, `P` partial, `N` no, `U` unverified/unknown, `—` not applicable. The vector describes existing capability; it does not replace the requirement outcome in §12.

| Capability | Scope/evidence | I | W | S | O | V | B | A | Notes |
|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| Flash Attention | Setting + server arguments, pinned build | Y | P | U | U | U | N | U | Checkbox/flag existence is not effective activation |
| GPU layer offload | Stored/passed `-ngl` | Y | Y | U | U | U | N | U | Per-model/VRAM fit unverified |
| Tensor split | Stored/passed | Y | Y | U | U | U | N | U | Topology/interconnect untested |
| Main GPU | Stored | Y | N | U | U | U | N | U | Not passed by server-argument builder |
| KV offload | Stored | Y | N | U | U | U | N | U | Not passed |
| mmap/mlock/load mode | Stored | Y | N | U | U | U | N | U | Runtime default applies |
| KV-cache dtype/quantization | No product control found | N | N | U | U | U | N | U | Model/backend dependent |
| SWA full cache | No product control found | N | N | U | U | U | N | U | Applicability requires architecture evidence |
| CPU MoE offload | No product control found | N | N | U | U | U | N | U | MoE-specific |
| Tensor override | No product control found | N | N | U | U | U | N | U | Advanced/version-specific |
| Automatic fit (`--fit`) | Conditional application | P | P | U | U | U | N | U | Disabled for some SWA cases; no broad evidence |
| Context shift | Runtime option; no product control | P | N | U | U | U | N | U | UI rollover claims are not proof |
| Cache reuse | Runtime option; no product control | N | N | U | U | U | N | U | Required research |
| Context checkpoints | Runtime option; no product control | N | N | U | U | U | N | U | Not assumed portable |
| Standard samplers (temp/top-k/top-p/min-p/typical/repeat/presence/frequency) | UI + request fields | Y | Y | P | U | U | N | P | Generic runtime acceptance; quality not measured |
| Mirostat / dynamic temperature | UI + request fields | Y | Y | P | U | U | N | P | Interaction/quality unverified |
| DRY multiplier/base/allowed | UI + request fields | Y | Y | P | U | U | N | P | Partial implementation |
| DRY last-N / sequence breakers | No UI/request mapping | N | N | U | U | U | N | U | Runtime options exist |
| XTC | Runtime option; no UI/request | N | N | U | U | U | N | U | Required inventory |
| Top-n-sigma | Runtime option; no UI/request | N | N | U | U | U | N | U | Required inventory |
| Adaptive-P | Runtime option; no UI/request | N | N | U | U | U | N | U | Required inventory |
| Sampler order | Runtime option; no control | N | N | U | U | U | N | U | No chain control |
| Actual seed capture | Global seed field; no response parsing | P | Y | P | U | U | N | U | Request seed is not actual seed |
| Model hot-swap | Manual stop/spawn; owner requirement | P | P | U | U | U | N | Y* | *Product behavior accepted; implementation/validation not |
| Per-chat model preference | No field | N | N | — | — | — | N | U | Owner decision unresolved |
| One live prompt path | Split `App`/`chat.Service` evidence | P | N | N | U | U | N | U | Architecture prerequisite |
| Chat memory/injection | Memory records + live prompt split | P | P | U | U | U | N | U | No demonstrated single injection path |
| Durable chat attachments | `hasImage` only | N | N | — | — | — | N | U | Contract missing |
| Multi-character flush | No scoped reset API | N | N | U | U | U | N | U | New requests are not proof |
| Single-character no-flush | No per-turn flush code found | P | N | P | U | U | N | U | Ordinary continuation exists; no explicit contract |
| Author's Note | No field/UI | N | N | — | — | — | N | U | Recovered requirement |
| Token counter | No control/tokenizer display | N | N | — | — | — | N | U | Recovered requirement |
| Narration mode | No feature; narrator events differ | N | N | — | — | — | N | U | Semantics unresolved |
| Per-character four-property colors | Global role tokens only | P | N | — | — | — | N | U | State/UI missing |
| Portrait thumbnail | Finalized assistant path partial | P | P | P | U | U | N | U | Streaming path lacks thumbnail |
| Message-triggered image generation | No chat action | N | N | — | — | — | N | U | Orchestration entry missing |
| Conversation context rollover | UI/tutorial + approximate truncation | P | U | U | U | U | N | U | Semantics not verified |

### 4.4 Static versus runtime evidence

This plan was produced without loading a model, exercising the UI, running a build, executing a migration, capturing packets, or performing a benchmark. Therefore:

- **Source-confirmed** statements are structural evidence only.
- **Test-confirmed** statements are cited records; tests were not rerun for this plan.
- **Runtime-confirmed** claims are limited to existing records and are not extended to new configurations.
- **Unknown** remains unknown. It is not filled by inference from a control name.
- A successful API call, if one were observed, would not prove model compatibility, correctness, or user benefit.

No capability in this plan is promoted to **Operational**, **Validated**, **Benchmarked**, or **Accepted** by this planning mission.

---

## 5. Known contradictions and release-readiness risks

### 5.1 Contradictions preserved from the original plan

The existing `C-01`–`C-22` register is preserved. “Corrected” means this plan resolves the contradiction at the planning level only; it does not claim the implementation is fixed.

| ID | Contradiction or prior statement | Corrected planning position |
|---|---|---|
| C-01 | “OnShutdown is never wired.” | **Superseded.** Shutdown registration and owned-process supervision exist; packaged lifecycle acceptance remains unknown. |
| C-02 | Model hot-swap is “unaccepted.” | **Corrected.** Hot-swap is required product behavior (`ODR-02`); implementation is partial, unsafe around active streams, and unvalidated. |
| C-03 | General creative studio is “unestablished.” | **Corrected.** The direction is owner-approved (`ODR-03`); no modality is implemented or authorized by that direction. |
| C-04 | Desktop avatar remains an ordinary parked future item. | **Corrected.** It is excluded (`ODR-01`, `REC-C-06`); only historical evidence remains. |
| C-05 | Context rollover and multi-character flushing are one feature. | **Separated.** Rollover is single-character context capacity; multi-character flush is character-scoped isolation. |
| C-06 | Author's Note is “persisted,” but new chats are blank. | **Reconciled.** Per-chat retention with no inheritance across chats; schema not implemented. |
| C-07 | A first-pass llama.cpp inventory is a capability survey. | **Corrected.** The catalogue is unfrozen and scope-limited; multi-engine research is unstarted. |
| C-08 | “Local-first” is a completed property. | **Contradicted.** Inspected Civitai lookup and remote ASR/TTS paths violate the current claim; local-first is a requirement. |
| C-09 | Local-Only Mode enforces offline operation. | **Corrected.** It wraps renderer `fetch`/XHR only; Go, sidecar, and Wails paths remain outside it. |
| C-10 | Generic consent is informed consent. | **Corrected.** Consent must name destination, purpose, data, recurrence, retention, and refusal effect. |
| C-11 | Error reporting sends “no model names.” | **Unverified.** An arbitrary `Context` map exists; payload and redaction are not proven. |
| C-12 | F3-C freeze implies acceptance. | **Rejected.** Freeze preserves a boundary; F3-C acceptance is not granted. |
| C-13 | F3-B closure proves product-wide lore correctness. | **Scoped.** Closure covers its documented limitations only; live prompt synchronization is unverified. |
| C-14 | One live prompt path exists. | **Contradicted.** `App`/llama and `chat.Service` remain separate evidence paths. |
| C-15 | Chat and memory have one canonical root. | **Contradicted.** The root split remains an open F3-C defect. |
| C-16 | RAG affects live replies. | **Unproven.** RAG service wiring exists, but the live prompt path does not demonstrably consume it. |
| C-17 | Hot-swap implementation is accepted because it is required. | **Rejected.** Product requirement and implementation acceptance are independent. |
| C-18 | Per-chat model selection is settled. | **Unresolved** (`UD-01`). Current selection is global. |
| C-19 | Android remote is conditionally available. | **Unimplemented.** No server, pairing, transport, or auth exists; `NET-12` controls any future work. |
| C-20 | Model compatibility is generic across models/backends/hardware. | **Rejected.** Compatibility is target-specific and requires evidence. |
| C-21 | Durable media identity exists. | **Contradicted.** `ChatHistory` records only `hasImage`; display, input, and stored media are conflated. |
| C-22 | Request seed is the actual generation seed. | **Unproven.** Actual seed capture requires response evidence; a request field is not proof. |

### 5.2 Network contradictions

| ID | Finding | Required correction |
|---|---|---|
| NC-01 | Image scanning calls Civitai automatically with a model hash. | `NET-03`/`NET-13` require user initiation and purpose-specific consent; current behavior contradicts both. |
| NC-02 | Local-Only Mode is renderer-only. | `NET-08`/`NET-18` require backend/sidecar enforcement and an egress-blocked runtime test. |
| NC-03 | Error-report upload has no proven redaction or destination/retention contract. | `NET-05`/`NET-09`/`NET-17` require payload proof, consent, audit, and credential redaction. |
| NC-04 | Remote ASR (DashScope) and remappable MeloTTS are exposed without a separate owner decision. | `NET-02`/`NET-16`/`NET-17` require either local-only operation or explicit, secured, purpose-specific authorization. |
| NC-05 | Generic “allow download” language does not disclose destination, data, recurrence, or retention. | `NET-05`/`NET-07` require the six-question consent model. |
| NC-06 | No destination registry or connection audit exists. | `NET-04`/`NET-09`/`NET-14` require registry, audit, and change control. |
| NC-07 | Browser-open, download, and organize/track actions are easy to conflate. | `NET-06` requires three separate permissions and explicit refusals. |
| NC-08 | Proposed paths `P-01`–`P-08` could be mistaken for approved features. | `NET-19` requires disclosure, cache, failure, and authorization rules before any implementation. |

### 5.3 Release-readiness risks

1. **Foundation risk:** durable chat schema, root identity, and one prompt path are unresolved. Building memory, rollover, or hot-swap on them would encode ambiguity.
2. **Evidence risk:** existing controls are mostly at the “Implemented” stage. Release claims would outrun evidence.
3. **Compatibility risk:** model architecture, backend build, quantization, and hardware combinations are not mapped.
4. **Resource risk:** model loading, image generation, TTS, and chat do not share one enforced admission/release authority.
5. **Continuity risk:** memory, lore, RAG, token budget, and rollover lack a common accounting and injection basis.
6. **Media risk:** there is no durable attachment identity or complete lifecycle contract.
7. **Network risk:** local-first is contradicted; no centralized egress boundary, registry, or audit exists.
8. **Consent risk:** generic labels do not establish informed, purpose-scoped consent.
9. **Remote-inference risk:** remote ASR/TTS sit inside an ostensibly local baseline without a separate owner decision.
10. **Governance risk:** protected gates and frozen records can be accidentally treated as implementation authority or acceptance.

---

## 6. Unresolved questions and investigation requirements

### 6.1 Owner questions preserved

| ID | Question | Current constraint | Related package |
|---|---|---|---|
| UD-01 | Is model selection global or per-chat? | Current implementation is global; no per-chat field exists. | CWP-01 |
| UD-02 | Is Android remote in scope, and under what security model? | Unimplemented; off by default if authorized; no server/pairing/auth. | CWP-25 |
| UD-03 | What is the committed scope of Story Tools? | Page/scripts exist; ownership and scope are unknown. | CWP-28 |
| UD-04 | What are the semantics of narration mode? | No feature; narrator events differ. | CWP-17 |
| UD-05 | Which studio modality, if any, is authorized next? | Direction approved; no modality authorized. | CWP-19 |
| UD-06 | What memory review policy is required? | Manual/automated memory exists; structured inventory and continuity are planned. | CWP-06 |

No answer is supplied here. A work package may investigate the question, but it may not silently choose the owner's answer.

### 6.2 Investigation prerequisites

Investigations are prerequisites, not implementation packages. They may conclude “not viable,” “requires owner decision,” or “return for more evidence.”

| ID | Investigation | Why it blocks |
|---|---|---|
| INV-01 | Build a packaged runtime baseline without changing gates, models, or settings. | No runtime baseline exists for lifecycle, prompt, media, capability, or network claims. |
| INV-02 | Trace the live prompt path from UI action through `App`, `chat.Service`, llama request, and response persistence. | Memory, lore, RAG, token budget, rollover, and provenance cannot share an unproven path. |
| INV-03 | Resolve F3-C root, alias, startup, wipe, and identity behavior. | Durable schema and continuity depend on canonical identity. |
| INV-04 | Define and test the durable chat/attachment schema against unavailable, moved, and deleted media. | `hasImage` is insufficient for display, input, and lifecycle semantics. |
| INV-05 | Inventory model/backend/hardware capability evidence and fallback behavior. | A control name cannot establish support or activation. |
| INV-06 | Map model load, active stream, image job, TTS, and resource admission into one state model. | Hot-swap and resource release cannot be safe without one authority. |
| INV-07 | Trace every outbound destination, caller, payload, trigger, cache, and failure path. | The network boundary cannot be enforced or consented to without a complete path register. |
| INV-08 | Verify error-report payload and redaction, including arbitrary context values. | The UI claim is not evidence. |
| INV-09 | Determine whether remote ASR/TTS are product requirements or removable options. | `NET-02` cannot be resolved by implementation convenience. |
| INV-10 | Establish a target hardware/backend/model compatibility matrix. | Capability selection and benchmarks require declared scope. |
| INV-11 | Inspect F2-B5 workspace state and Story Tools consumers. | Approved direction does not establish committed scope or ownership. |
| INV-12 | Recover or bound the missing F3-B/C and W4 closure evidence. | Historical claims must not be treated as current product-wide verification. |

### 6.3 Source evidence that could not be recovered

- The exact historical UI planning document and full expansion text remain unrecovered.
- The formal W4 closure record may be missing; W4 is historical/superseded, not an open implementation mission.
- The F3-B/C acceptance evidence for every claimed surface was not independently reproduced.
- T2.9.5B `F1`/`F2` definitions were not independently verified.
- F2-B5 runtime UI state was not inspected in this mission.
- No runtime, benchmark, packet-capture, offline, cancellation, consent/audit, remote ASR/TTS, or Android evidence was produced.

These gaps are limitations. They are not filled by assumption.

### 6.4 RAG scope: no requirement identity, explicit owner decision required

The 66 reconciliation rows do **not** contain a RAG requirement. This correction does not invent one. RAG is nevertheless a source-confirmed subsystem with an integration gap (`S7`, `C-16`) and a possible embedding/resource consumer.

| Item | Disposition |
|---|---|
| Requirement identity | **None.** No RAG row is added to the 66. |
| Provisional package owner | `CWP-02` owns prompt-path integration only while scope is unresolved. |
| Full retrieval lifecycle | **Blocked** pending `OA-10`; no embedding model, index, document deletion, retention, rebuild, or resource claim is authorized. |
| Required owner question | Is RAG a committed product capability, an explicitly excluded legacy surface, or a future scope? |
| If authorized later | The owner must create the requirement/reconciliation identity before a package can claim RAG acceptance. |
| Embedding model selection/compatibility | Not authorized. If `OA-10` is granted, `CWP-11` supplies identity evidence and `CWP-15` includes the embedding consumer in resource authority. |
| Index lifecycle/versioning | Not authorized. `CWP-02` would own the retrieval contract; the index must be rebuildable and versioned against the embedding model. |
| Retrieval correctness/relevance | Not authorized. `CWP-02` would own exactly-once injection and retrieval fixtures; no current claim that documents affect replies. |
| Document deletion/retention | Not authorized. `CWP-10` would own durable document/media reference semantics and deletion propagation. |
| Resource consumption and hardware coordination | Not authorized. `CWP-15` would own embedding/index workload admission only if `OA-10` is granted. |
| Validation/acceptance | `RG-13` is conditional and not passed. It requires live effect, identity, index versioning, deletion/rebuild, retention, and resource evidence. |
| Current acceptance | **None.** No claim that documents affect replies. |

---

## 7. Architecture decisions still required

These are technical/architecture decisions, not owner product decisions. They may be investigated and proposed by engineering, but several cannot be settled without the owner questions in §6.1 or the protected F3-C boundary. None is approved by this plan.

| ID | Decision required | Blocked by / affects | Required evidence before decision | Related package |
|---|---|---|---|---|
| AD-01 | Canonical durable chat identity and target-set schema | `UD-01`, F3-C root/alias, attachments, Author's Note, per-chat settings | Root/alias resolution; migration and unavailable-media behavior | CWP-01 |
| AD-02 | Which prompt path is authoritative | Split `App`/`chat.Service`; memory, lore, RAG, rollover, token budget | End-to-end path trace with one injected message set | CWP-02 |
| AD-03 | Root/path identity and alias rejection | `D-ROOT`, F3-C, memory/chat roots | F3-C defect resolution evidence; canonical path tests | CWP-01 |
| AD-04 | Tokenizer-accurate context budget | Model templates, memory, lore, rollover, token counter | Per-model tokenizer evidence and budget calculation | CWP-03 |
| AD-05 | Single model/resource state authority | Hot-swap, image jobs, TTS, VRAM admission, release | Initial evidence: `INV-06` plus the `CWP-15a` resource-authority contract. AD-05 is recorded at the `15a → 08a` boundary. The later `15a → 08a → 15b → 08b` sequence validates the decision; it is not a prerequisite for starting `15a`. | CWP-08, CWP-15 |
| AD-06 | Durable media and attachment identity | `hasImage`, display/input/stored media, delete/move/export | Lifecycle contract and unavailable-reference tests | CWP-10 |
| AD-07 | Capability evidence and selection contract | Techniques, compatibility, overrides, benchmarks | Target-scope vectors; conflict rules; fallback behavior | CWP-12, CWP-13, CWP-14 |
| AD-08 | Central network enforcement and registry | Local-Only Mode, Go/sidecar/Wails egress, consent, audit | Complete path inventory; fail-closed enforcement design | CWP-20 |
| AD-09 | Remote ASR/TTS boundary | `NET-02`, `NET-16`, `NET-17`, product local-first claim | Owner decision; local fallback; TLS/auth/redaction proof | CWP-22 |
| AD-10 | Android remote boundary | `NET-12`, `UD-02`, `MOB-1`–`MOB-4` | Owner authorization; pairing/auth/transport/revocation design | CWP-25 |
| AD-11 | Studio job/media contract | Separate image/video/audio surfaces; shared catalogue | Modality scope; asset identity; delete/export semantics | CWP-19 |
| AD-12 | Optional persona/AI-as-user boundary | `REC-CH-08`; character ownership and attribution | Feasibility study; explicit owner decision | CWP-26 |

**Decision discipline.** An architecture decision record may select an approach only after its required evidence exists. A decision must name the scope it applies to (model, backend, hardware, content type, destination) and the conditions under which it fails closed. A decision that cannot state its failure behavior is not complete.

### 7.1 Authorization gates

These are authorization checkpoints, not owner decisions already made.

| Gate | What must be authorized | Current status |
|---|---|---|
| OA-01 | Any schema migration or durable-state cutover | Not authorized. F3-C/D-STORAGE constraints remain. |
| OA-02 | Any change to the live prompt path or model lifecycle | Not authorized. Investigation only. |
| OA-03 | Enabling, removing, or changing an inference technique | Not authorized. Research/evidence first. |
| OA-04 | Any new external destination or purpose | Not authorized. Registry and consent design first. |
| OA-05 | Remote ASR/TTS retention or removal | Not authorized. Owner decision required. |
| OA-06 | Android remote implementation | Not authorized. `UD-02`/`NET-12` remain unresolved. |
| OA-07 | Any studio modality implementation | Direction approved; modality authorization absent. |
| OA-08 | Any claim of release readiness or acceptance | Not authorized. Release gates remain open. |
| OA-09 | Any investigation that loads a model, runs a benchmark, exercises the network, inspects user content, or captures runtime evidence | Not authorized. Proposed investigations require explicit owner authorization to begin. |
| OA-10 | RAG committed scope and full retrieval lifecycle | Not authorized. No requirement identity exists; `CWP-02` is a provisional prompt-path owner only. |

## 8. Dependency-ordered work packages

### 8.1 Ordering principle

The order below follows architectural dependency, not UI page order, apparent feature size, or implementation convenience. The sequence is:

1. **Preserve the boundary.** Do not disturb protected gates or create new acceptance records.
2. **Establish identity and truth.** Canonical chat/root/media identity before continuity features.
3. **Establish one prompt and accounting path.** One path, one tokenizer budget, one character-scope rule.
4. **Build continuity.** Memory, lore, rollover, provenance, and attachments on those foundations.
5. **Make model operation safe.** Hot-swap, resource authority, compatibility evidence, selection.
6. **Repair presentation.** UI, narration, and chat-to-image orchestration after their data paths exist.
7. **Bound the studio.** Shared foundation before modality implementation.
8. **Bound the network.** Registry, enforcement, consent, audit, and negative tests before new external paths.
9. **Keep remote and optional features gated.** Android and optional persona remain conditional.
10. **Close traceability and validation.** No package is complete because code exists.

A package may begin investigation only when its prerequisites are evidenced and the owner has granted `OA-09`. A package must not silently assume a prerequisite is complete because another package exists.

### 8.2 Work-package format

Each package records:

- unique identifier (`CWP-*`, a new family separate from original requirement IDs);
- exact source requirement/reconciliation identifiers;
- current evidence and known gaps;
- prerequisites;
- dependencies and affected components;
- required investigation;
- expected deliverable;
- validation and acceptance evidence;
- risks and unresolved decisions;
- required owner authorization.

No package includes a duration, staffing, or effort commitment.

### 8.3 Keystone unblocking sequence (CWP-01)

`CWP-01` is the keystone, but it may not absorb a protected decision. The sequence is:

| Order | Gate / authority | Required release evidence | Work permitted before release |
|---|---|---|---|
| 1 | `UD-01` — Project Owner | Explicit global-or-per-chat model decision | Settings inventory; schema options; no model field committed |
| 2 | F3-C Phase 0 / D-ROOT — Project Owner acting through the existing F3-C record, not a new authority | Authorized phase-0 disposition for root, alias, startup, wipe, and identity; `INV-03` supplies evidence, not release | Root/alias evidence collection; no root change; no identity decision |
| 3 | D-STORAGE — existing F3-C storage gate/owner | Authorized storage path and migration conditions; `OA-01` | Schema proposal and migration-safety design; no migration |
| 4 | `CWP-01` design review | All three release records plus `INV-04` | Schema proposal only; no cutover |

`INV-03` is evidence input to the protected F3-C path. It does not resolve F3-C, grant R6, or authorize phase-0 work. A local root/alias decision inside `CWP-01` is prohibited.

### 8.4 Hot-swap / resource handoff: acyclic sequence

The former `CWP-08 ↔ CWP-15` cycle is replaced by a named handoff. `CWP-15a/15b` and `CWP-08a/08b` are phases inside the existing 28 packages, not new package identities.

| Step | Artifact | Owner | May start when |
|---|---|---|---|
| 1 | `CWP-15a` — resource-authority contract: states, consumers, admission, release, failure | CWP-15 | `INV-06`; owner authorization to prepare the contract. `AD-05` is not required to start `15a`. |
| 2 | `CWP-08a` — hot-swap consumer/state contract: drain, stop, spawn, ready, commit, rollback | CWP-08 | `CWP-01`, `CWP-02`, `CWP-15a`, and `AD-05` recorded at the `15a → 08a` boundary |
| 3 | `CWP-15b` — admission/release enforcement **design** | CWP-15 | `CWP-15a` and `CWP-08a` |
| 4 | `CWP-08b` — hot-swap implementation design | CWP-08 | `CWP-08a` and `CWP-15b` |

Edges: `INV-06 + design authorization → CWP-15a → AD-05 → CWP-08a → CWP-15b → CWP-08b`. `AD-05` is recorded after `15a`; it is not an input to `15a`. No edge runs from `CWP-08` back to `CWP-15` before `CWP-15a`. No implementation is authorized at any step. `CWP-15b` is a design artifact, not enforcement evidence.

### 8.5 Investigation authorization

Every investigation in this plan is **proposed**. This unadopted plan does not authorize its own investigations. `OA-09` is the owner gate for any investigation that loads a model, runs a benchmark, exercises the network, inspects user content, or captures runtime evidence. A completed investigation does not authorize implementation or acceptance.


### CWP-01 — Canonical chat truth set and durable schema

- **Sources:** `REC-CH-01`, `REC-CH-02`, `UD-01`
- **Current evidence/gap:** Per-chat Author's Note and composer controls are absent. The durable chat schema is not established; `UD-01` leaves model scope unresolved. Chat and memory roots are split, and F3-C root/alias defects remain open.
- **Prerequisites:** `INV-03` root/alias evidence; `INV-04` schema investigation; owner answer to `UD-01`; F3-C/D-STORAGE release sequence in §8.3. No migration before all gates.
- **Dependencies/components:** S2 chat, S3 characters, S6 memory, S17 content ownership, `ChatHistory`, settings/session state.
- **Required investigation:** Define the smallest durable chat identity that survives rename, move, import, and restart; enumerate every per-chat setting; prove unavailable-media and missing-root behavior; decide what remains global if `UD-01` selects global model scope.
- **Expected deliverable:** A proposed schema, state-ownership record, and migration-safety contract covering chat, memory, character, and world stores: free-space preflight, backup/snapshot with named restore, atomic staging/commit, post-migration integrity verification, resumable or cleanly abortable interruption, single-writer enforcement, and explicit failure behavior. Media migration remains `CWP-10`, dependent on this contract. No migration executes.
- **Validation/acceptance:** Schema review against F3-C root/alias constraints; round-trip design for create, rename, move, import, export, restart, and delete; fault fixtures for power loss, partial write, insufficient space, and concurrent writer across chat, memory, character, and world stores; explicit failure behavior for missing roots. Implementation acceptance requires a later owner-authorized migration and runtime evidence.
- **Risks/unresolved:** `UD-01` unanswered; root/alias defects open; browser/native dual ownership; migration may be blocked by D-STORAGE. No schema may assume `chat.Service` is authoritative.
- **Owner authorization:** **Required before any migration or durable-state cutover.** Investigation is proposed; it requires `OA-09`. Implementation is not authorized.

### CWP-02 — One authoritative live prompt path

- **Sources:** `REC-UI-05`, `REC-M-02`
- **Current evidence/gap:** The live stream uses `App`/llama while `chat.Service` is a separate prompt path. No demonstrated single injection point exists for characters, memory, lore, RAG, attachments, or rollover.
- **Prerequisites:** `INV-02` path trace; `CWP-01` identity decisions for chat-scoped inputs; `INV-03` root evidence.
- **Dependencies/components:** S2 chat/prompt assembly, S5 lore, S6 memory, S7 RAG, S13 vision, llama HTTP client, frontend chat stream.
- **Required investigation:** Trace one user turn end to end; identify every prompt contributor and every persistence write; prove which path owns model parameters, seed, and character context; identify duplicate or dead wiring.
- **Expected deliverable:** A single authoritative prompt-path design with an explicit contributor order, ownership map, and removal/deprecation plan. RAG is a **provisional, blocked sub-scope** under §6.4: only an ownership/decision brief is permitted; no retrieval lifecycle implementation.
- **Validation/acceptance:** A test design that injects a known character, lore, memory, attachment, and rollover state and proves each appears exactly once. RAG acceptance is conditional on `OA-10`; no RAG claim is possible now. Acceptance requires authorized implementation and runtime capture.
- **Risks/unresolved:** Removing the wrong path could break background or story-tool flows; `chat.Service` consumers may be undocumented; RAG may currently be inert. No reconciliation row authorizes RAG scope; `OA-10` is required. Do not “fix” by adding a second injection point.
- **Owner authorization:** **Implementation authorization required.** Investigation is proposed and requires `OA-09`; no code change is authorized.

### CWP-03 — Tokenizer-accurate context accounting and token counter

- **Sources:** `REC-UI-05`, `REC-M-02`, `REC-CH-11`
- **Current evidence/gap:** No product token counter or tokenizer display was found. Context budgeting is approximate and split across prompt, memory, lore, and rollover. Tokenizer behavior is model/template dependent.
- **Prerequisites:** `CWP-02` authoritative prompt path; `CWP-11` model/template evidence; `AD-04`.
- **Dependencies/components:** S2 prompt assembly, S5 lore, S6 memory, S14 model intelligence, chat UI.
- **Required investigation:** Identify the tokenizer/template actually used per model family; measure system, character, lore, memory, attachment, and response reserve costs; define overflow policy without silently truncating required content.
- **Expected deliverable:** A shared accounting contract and proposed UI display, including what is counted, what is reserved, and what is dropped first. No model run is required to document the contract; tokenizer verification requires authorized evidence.
- **Validation/acceptance:** Known fixture prompts across at least the supported template families; expected counts with declared tolerances; overflow and rollover boundary cases. A counter that disagrees with the actual request is a failure.
- **Risks/unresolved:** Tokenizer mismatch, template fallbacks, multimodal tokens, and tool/system overhead. The counter must not be presented as exact until validated per model.
- **Owner authorization:** **Implementation authorization required.** No tokenizer benchmark or model load is authorized here.

### CWP-04 — Character-scoped context isolation and multi-character flush

- **Sources:** `REC-CH-04`, `REC-CH-05`, `REC-CH-06`, `ODR-08`, `ODR-11`
- **Current evidence/gap:** No scoped reset API or demonstrated per-character inference isolation was found. Sequential HTTP calls are not proof of a flushed backend context. Single-character continuation and multi-character flush are distinct requirements.
- **Prerequisites:** `CWP-02` one prompt path; `CWP-01` chat/character identity; `CWP-03` token accounting; `AD-05` state authority.
- **Dependencies/components:** S2 chat, S3 characters, S6 memory, llama context lifecycle, session state.
- **Required investigation:** Determine whether the backend exposes a reliable context-reset or session boundary; define target-set membership and speaker attribution; prove what is retained across character switches and what must be flushed.
- **Expected deliverable:** A character-scope contract covering target set, flush trigger, retained state, and failure behavior, plus a test design that distinguishes transport isolation from inference-context isolation.
- **Validation/acceptance:** Multi-character turn sequence with cross-character leakage checks; single-character continuation without an unnecessary flush; explicit failure if reset cannot be confirmed. A successful HTTP request is not acceptance.
- **Risks/unresolved:** Backend reset semantics may be unavailable or asynchronous; flushing may destroy desired continuity; speaker attribution may be ambiguous. `ODR-11` permits isolated functions but does not prove implementation.
- **Owner authorization:** **Implementation authorization required.** No context-reset experiment is authorized without a controlled test plan.

### CWP-05 — Conversation-context rollover for a single character

- **Sources:** `REC-CH-03`, `REC-CH-11`, `ODR-11`
- **Current evidence/gap:** UI/tutorial language and approximate truncation exist; rollover semantics are not verified. There is no proven boundary between “continue the same conversation,” “summarize,” and “start a new context.”
- **Prerequisites:** `CWP-02`, `CWP-03`; canonical chat identity. Multi-character flush is explicitly **not** a prerequisite.
- **Dependencies/components:** S2 chat, S6 memory, prompt assembly, token budget, UI status.
- **Required investigation:** Define trigger conditions (token threshold, turn count, model change), retained summary, user visibility, and recovery from failed rollover. Separate rollover from multi-character flushing.
- **Expected deliverable:** A rollover state machine and retention contract, including what the model sees before/after rollover and how the user can tell it happened.
- **Validation/acceptance:** Boundary tests at under/at/over budget; continuity checks across rollover; no duplicated or silently dropped messages; recoverable failure. Runtime proof is required before “Operational.”
- **Risks/unresolved:** Summarization quality, tokenizer variance, memory duplication, and user confusion. Rollover must not become an implicit memory rewrite.
- **Owner authorization:** **Implementation authorization required.** `ODR-11` is a product requirement, not implementation acceptance.

### CWP-06 — Durable memory and conversational continuity

- **Sources:** `REC-CH-09`, `REC-CH-10`, `ODR-04`, `UD-06`
- **Current evidence/gap:** Manual/automated memory, pinning, decay, and candidate review exist, but there is no demonstrated single injection path, structured inventory, clothing/equipment state, or application-owned continuity. `UD-06` leaves review policy unresolved.
- **Prerequisites:** `CWP-01`, `CWP-02`, `CWP-03`; owner answer to `UD-06`; root/alias evidence.
- **Dependencies/components:** S6 memory, S2 prompt, S3 characters, S17 content ownership, memory review UI.
- **Required investigation:** Inventory every memory type and writer; define authority, deduplication, expiry, user correction, and injection priority; determine what “durable” means across import, export, move, and model change.
- **Expected deliverable:** A memory ownership and lifecycle contract, including structured inventory, candidate review, injection order, and failure behavior when memory is missing or conflicting.
- **Validation/acceptance:** Deterministic injection fixtures; duplicate/conflict tests; import/export/restart continuity; user correction and deletion propagation. No memory is “working” until its live effect is demonstrated.
- **Risks/unresolved:** Two writers may disagree; memory may be silently truncated; `UD-06` may change review policy; privacy of inferred personal data.
- **Owner authorization:** **Implementation authorization required.** Investigation is proposed and requires `OA-09`; no schema or prompt change is authorized.

### CWP-07 — Lore scope, precedence, and prompt synchronization

- **Sources:** `REC-CH-07`
- **Current evidence/gap:** Lorebook store and trigger/conflict logic are source- and test-confirmed within F3-B's documented limitations. Prompt synchronization with a running model remains unverified. `REC-CH-07` is partially satisfied.
- **Prerequisites:** `CWP-02` one prompt path; `CWP-03` token accounting; F3-B closure limitations treated as scope, not product-wide proof.
- **Dependencies/components:** S5 lorebook store/triggers, S2 prompt assembly, character/world scope, cache.
- **Required investigation:** Prove which lore entries are selected, in what order, for which character/world, and how conflicts are resolved at prompt time; identify stale cache and duplicate-entry paths.
- **Expected deliverable:** A lore precedence and injection contract with deterministic selection fixtures and explicit conflict outcomes.
- **Validation/acceptance:** Same input produces the same selected set; conflict precedence is documented; token overflow behavior is defined; live prompt capture shows the selected entries once. Store tests alone are insufficient.
- **Risks/unresolved:** Trigger semantics may be model-dependent; cache invalidation may be incomplete; lore and memory may compete for the same budget.
- **Owner authorization:** **Implementation authorization required** for prompt-path changes. Investigation is proposed and requires `OA-09`.

### CWP-08 — Conversation-safe model hot-swap

- **Sources:** `ODR-02`
- **Current evidence/gap:** Manual load/unload and pre-stop admission exist, but `LoadModel` does not coordinate an active stream. After the old server stops, spawn/startup failure has no old-model rollback. Hot-swap is required product behavior but not accepted implementation behavior.
- **Prerequisites:** `CWP-01` chat identity, `CWP-02` one prompt path, `CWP-15a` resource-authority contract, `AD-05`; `CWP-15b` before `CWP-08b` per §8.4.
- **Dependencies/components:** S15 model manager, S1 process supervision, S2 active stream, S16 resource manager, model UI/status.
- **Required investigation:** Enumerate every active consumer (chat stream, image job, TTS, embedding/index if `OA-10`, downloads); define drain, stop, spawn, readiness, commit, rollback, and cancellation states; define architecture/tokenizer/template compatibility and re-fit behavior.
- **Expected deliverable:** A hot-swap state machine and rollback contract plus a **conversation-coherence contract**: canonical transcript, re-render under the new template, re-tokenize, re-fit to the new budget, preserve/report what changed, and explicitly refuse an incoherent swap. It must state what happens to the old model before the new one is confirmed ready.
- **Validation/acceptance:** Fault-injected tests for spawn failure, readiness timeout, mid-stream switch, cancellation, and rollback; cross-template and cross-context-size fixtures; a visibly reported truncation/re-fit; no prompt sent to an unready model; no old model silently destroyed before a safe commit point. Runtime proof required.
- **Risks/unresolved:** Shared resource authority is not implemented; rollback may require holding two models and exceeding VRAM; per-chat versus global model scope is unresolved (`UD-01`); templates/tokenizers may not share message rendering. The §8.4 edge is one-way: `15a → 08a → 15b → 08b`.
- **Owner authorization:** **Implementation authorization required after design review.** Product requirement already exists; it does not authorize code.

### CWP-09 — Per-turn model/configuration provenance and actual seed

- **Sources:** `ODR-05`, `ODR-06`
- **Current evidence/gap:** There is no per-turn model/configuration provenance. A request seed field exists, but actual seed capture is unproven. The product cannot reconstruct which model, template, sampler, or settings produced a response.
- **Prerequisites:** `CWP-01` durable chat identity; `CWP-02` authoritative request path; `CWP-11` model/template identity.
- **Dependencies/components:** S2 chat, S14 model intelligence, S18 settings, response parsing, gallery/image metadata (separate scope).
- **Required investigation:** Identify every parameter that affects a turn; determine what the backend returns or exposes; define immutable turn provenance versus mutable settings; prove whether actual seed is recoverable.
- **Expected deliverable:** A provenance record contract and replay-readiness assessment, including explicit “actual seed unavailable” behavior. It must not label the request seed as the actual seed.
- **Validation/acceptance:** Two turns with different models/settings produce distinct provenance; restart retains it; a reconstructed request matches the recorded values; unavailable fields are explicit, not guessed. Runtime evidence required.
- **Risks/unresolved:** The backend may not return actual seed; templates/settings may change underneath a chat; provenance storage may conflict with F3-C schema limits.
- **Owner authorization:** **Implementation authorization required.** Investigation is proposed and requires `OA-09`.

### CWP-10 — Durable attachments and media lifecycle

- **Sources:** `REC-P-01`, `REC-P-02`, `REC-P-03`, `REC-P-04`, `ODR-07`
- **Current evidence/gap:** Chat schema records only `hasImage`; saved audio has no chat association; displayed and model-input media are conflated; delete/move/export/import/restore semantics are missing.
- **Prerequisites:** `CWP-01` durable chat identity, `INV-04`, `AD-06`; F3-C/D-STORAGE authorization before migration.
- **Dependencies/components:** S2 chat, S8 image gallery, S9 video, S10/S11 audio, S17 libraries/import-export, metadata store.
- **Required investigation:** Inventory every media producer and consumer; define stable asset identity, derivatives, source vs display vs model-input references; prove unavailable/moved/deleted behavior; map retention and export.
- **Expected deliverable:** A media lifecycle contract with reference states (`available`, `missing`, `moved`, `deleted`, `regenerable`), ownership, and migration design. It includes the same migration-safety contract as `CWP-01`: free-space preflight, backup/restore, atomic staging/commit, integrity verification, interruption handling, and single-writer enforcement. No file is deleted or moved by the investigation.
- **Validation/acceptance:** Fixtures for each reference state; chat reload with missing media; import/export/restore round-trip; delete does not orphan or silently substitute; model input uses the correct derivative; power-loss, partial-write, insufficient-space, and concurrent-writer fixtures. Runtime proof required.
- **Risks/unresolved:** Existing `hasImage` data cannot be reconstructed into durable references without a migration; image metadata is generation metadata, not chat provenance; storage may be blocked.
- **Owner authorization:** **Required before migration or destructive media semantics.** Design work only.

### CWP-11 — Model intelligence, template identity, and compatibility evidence

- **Sources:** `REC-M-01`, `REC-M-13`
- **Current evidence/gap:** Registry, evidence-based profiles, template inspection, image-family detection, and discovery records exist. Filename fallback and unverified compatibility remain risks. Model identity is not strong enough for per-turn provenance or target-specific technique selection.
- **Prerequisites:** `INV-05`; model files available for authorized inspection only; no model load.
- **Dependencies/components:** S14 model intelligence, S15 loader, S18 settings, llama templates, image registry.
- **Required investigation:** Define model identity (hash/architecture/quantization/template/backend), distinguish verified metadata from filename inference, and enumerate unknown fields explicitly.
- **Expected deliverable:** A model-evidence record and compatibility schema that supports later selection without pretending unknown fields are supported. It must include provenance for every field.
- **Validation/acceptance:** Same model resolves to the same identity across discovery paths; conflicting metadata is surfaced; filename fallback is labelled; no technique is marked supported from a name alone.
- **Risks/unresolved:** Models may be renamed/moved; templates may be embedded or external; backend versions may not expose architecture facts.
- **Owner authorization:** **Implementation authorization required.** No model load or download is authorized.

### CWP-12 — Technique catalogue and compatibility/conflict graph

- **Sources:** `REC-M-11`, `REC-M-12`, `REC-M-14`, `REC-M-18`, `ODR-10`
- **Current evidence/gap:** The llama.cpp inventory is a single-binary subset; the catalogue is unfrozen. No conflict graph exists for techniques that disable, invalidate, or degrade each other. Multi-engine research is unstarted.
- **Prerequisites:** `CWP-11`; `INV-05`; `INV-10` for target scope.
- **Dependencies/components:** S14 model intelligence, S15 loader, inference settings, backend adapters, documentation/research corpus.
- **Required investigation:** Build the research inventory by technique family and backend; define support predicates, conflicts, fallbacks, and evidence requirements; include MoE, SWA, recurrent/hybrid, quantization, multi-GPU, cache, and sampler interactions.
- **Expected deliverable:** A versioned, unfrozen technique catalogue and conflict graph with explicit unknown states and required evidence per edge. It must preserve the research → evidence → compatibility → selection → validation → authorization progression.
- **Validation/acceptance:** Every technique has a scope and evidence vector; every conflict names the affected target; no cycle is silently accepted; adding an unknown backend cannot auto-enable a technique.
- **Risks/unresolved:** Catalogue breadth can outrun evidence; backend versions change; a conflict may be model-specific rather than global.
- **Owner authorization:** **Research/documentation only.** No technique enablement or backend implementation is authorized.

### CWP-13 — Existing-capability verification and benchmark campaign

- **Sources:** `REC-M-03`, `REC-M-04`, `REC-M-05`, `REC-M-06`, `REC-M-17`, `ODR-10`
- **Current evidence/gap:** Most controls are implemented but not proven wired, supported, operational, correct, or beneficial. Sampler completeness, context shift, cache reuse, checkpoints, and actual seed capture are unresolved. No quality benchmark exists.
- **Prerequisites:** `CWP-11`, `CWP-12`; `INV-05`, `INV-10`; `OA-09`; explicit test scope and owner authorization for model loads/benchmarks.
- **Dependencies/components:** S14/S15/S16, inference server, samplers, context options, output-quality fixtures, logs.
- **Required investigation:** Select a declared matrix of model/backend/hardware; define activation evidence, correctness fixtures, quality checks, memory/speed measures, and failure/fallback detection. Separate speed/memory from coherence/continuity.
- **Expected deliverable:** A verification protocol and evidence ledger per capability and scope. Results may be “unsupported,” “failed,” or “inconclusive”; those are valid outcomes. No result is prewritten.
- **Validation/acceptance:** Reproducible commands and environment capture; a control that does not activate is not operational; a benchmark without a comparator is not a benefit claim. Owner authorization is required before any run.
- **Risks/unresolved:** Benchmark cost, hardware variance, model licensing, output-quality subjectivity, and the risk of generalizing one result. This plan performs no benchmark.
- **Owner authorization:** **Explicit authorization required before loading models or running benchmarks.** This package is a protocol, not a run.

### CWP-14 — Automatic capability selection and scoped manual overrides

- **Sources:** `REC-M-15`, `REC-M-16`, `ODR-10`
- **Current evidence/gap:** Automatic selection is missing; manual overrides are partial. There is no evidence sufficient to select techniques safely per model/backend/hardware.
- **Prerequisites:** `CWP-11`, `CWP-12`, `CWP-13`; `AD-07`; owner acceptance of selection policy.
- **Dependencies/components:** S14 model intelligence, S15 loader, S16 resources, inference settings, UI indicators.
- **Required investigation:** Define selection inputs, precedence, explicit user override, scope of an override (turn/chat/model), fallback, and visibility when evidence is insufficient.
- **Expected deliverable:** A deterministic selection contract and override model. Default behavior when evidence is missing must be conservative and visible, not an optimistic enable.
- **Validation/acceptance:** Same evidence produces the same selection; overrides are scoped and expire correctly; unsupported/conflicting selections fail closed; user can see why a technique was not selected.
- **Risks/unresolved:** Automatic selection can hide trade-offs; overrides can create unsafe combinations; compatibility evidence may be too sparse to automate.
- **Owner authorization:** **Implementation authorization required.** No automatic enablement is authorized by this plan.

### CWP-15 — Model/resource state authority and safe release

- **Sources:** `REC-M-07`, `REC-M-08`
- **Current evidence/gap:** Resource manager and chat/image arbitration exist, but admission is advisory, image auto-load registration is inconsistent, TTS is outside admission, and weight release/VRAM reclamation is unmeasured.
- **Prerequisites:** `INV-06`; owner authorization to prepare the resource contract; `CWP-13` scope for any measurement. Phase handoff only: `CWP-15a` starts first; `CWP-08a` consumes `15a`; `CWP-15b` consumes `08a`. There is no reverse package-level dependency.
- **Dependencies/components:** S1 process supervision, S8 image jobs, S10/S11 audio, S15 model manager, S16 resource manager, VRAM dialogs, embedding/index consumer only if `OA-10` authorizes it.
- **Required investigation:** Inventory every model consumer and every owned process; define one state machine for unloaded/loading/ready/busy/swapping/failed/unloading; determine admission and release guarantees; identify advisory-only checks; include embedding/index workload only if authorized.
- **Expected deliverable:** `CWP-15a` resource-authority contract (states, consumers, admission, release, timeout, failure) followed by `CWP-15b` admission/release **enforcement design**. `AD-05` is recorded at the `15a → 08a` boundary; it is not an input to `15a`. It must state what happens when VRAM cannot be reclaimed. No enforcement is claimed to exist.
- **Validation/acceptance:** Concurrent chat/image/TTS scenarios; swap under load; failed load; cancellation; unload while busy; evidence that release occurred or an explicit failure is surfaced; no orphan embedding/index consumer if RAG is authorized. No “freed” claim without measurement.
- **Risks/unresolved:** Two models may not fit simultaneously; process cleanup may be incomplete; image/TTS paths may not accept admission; release measurement requires authorized runtime work; RAG is blocked pending `OA-10`.
- **Owner authorization:** **Implementation authorization required.** No process, VRAM, or model experiment is authorized here.

### CWP-16 — Chat presentation, identity, and per-character visual properties

- **Sources:** `REC-UI-01`, `REC-UI-02`, `REC-UI-03`, `ODR-09`
- **Current evidence/gap:** Some styling/portrait pieces exist; per-character four-property colors and state/UI are incomplete; streaming portraits lack thumbnails. Earlier visual/ghost-ID/duplicate-element findings remain relevant.
- **Prerequisites:** `CWP-01` character/chat identity; `CWP-03` token display contract; `CWP-10` media identity for portraits.
- **Dependencies/components:** Chat UI, character model/store, theme tokens, streaming renderer, portrait asset path.
- **Required investigation:** Inventory every character visual property and its persistence; reproduce the portrait state for finalized and streaming messages; identify global role tokens that incorrectly override character identity.
- **Expected deliverable:** A chat presentation contract covering colors, portraits, streaming placeholders, empty/error states, and accessibility. No visual redesign beyond the stated requirements.
- **Validation/acceptance:** Character A/B isolation; restart persistence; streaming and finalized parity; no duplicate IDs or ghost controls; unavailable portrait fallback. Visual verification requires an actual rendered check in an authorized session.
- **Risks/unresolved:** Browser/native character state may disagree; portrait assets may be missing; chat design changes could be mistaken for studio work.
- **Owner authorization:** **Implementation authorization required.** No visual verification was performed for this plan.

### CWP-17 — Narration mode semantics and implementation boundary

- **Sources:** `REC-UI-04`, `UD-04`, `ODR-09`
- **Current evidence/gap:** No narration feature exists. Narrator events differ from the requested mode. Semantics remain unresolved (`UD-04`).
- **Prerequisites:** Owner answer to `UD-04`; `CWP-02` one prompt path; `CWP-03` context accounting.
- **Dependencies/components:** Chat UI, prompt assembly, character/speaker model, streaming renderer, message schema.
- **Required investigation:** Define narration speaker, attribution, formatting, prompt behavior, persistence, and interaction with multi-character scenes. Distinguish narration from author note and system prompt.
- **Expected deliverable:** A semantics specification and implementation boundary. Until `UD-04` is answered, the deliverable is a decision brief, not a feature.
- **Validation/acceptance:** Speaker attribution tests; narration does not masquerade as a character message; persistence/reload; prompt inclusion occurs once; user can distinguish narration visually and structurally.
- **Risks/unresolved:** `UD-04` may change the requirement; narration can conflict with character isolation and message roles.
- **Owner authorization:** **Owner decision required before design finalization; implementation authorization required separately.**

### CWP-18 — Message-triggered image generation

- **Sources:** `REC-UI-06`, `ODR-09`
- **Current evidence/gap:** No chat action for image generation was found. The image engine and generation metadata exist, but orchestration from a message is missing. `REC-P-05` image metadata is satisfied for image scope, but it is not chat provenance and is not a source of new work here.
- **Prerequisites:** `CWP-02` prompt path, `CWP-10` media lifecycle, `CWP-15` resource admission, `CWP-19` studio boundaries.
- **Dependencies/components:** Chat UI, image job manager, gallery/media store, prompt extraction, resource manager.
- **Required investigation:** Define trigger, prompt source, user confirmation, model/resource admission, progress/error display, result attachment, and cancellation. Decide whether generation is inline or a handoff to the image workspace.
- **Expected deliverable:** An orchestration contract and explicit non-goals. It must not merge chat provenance with image-generation metadata.
- **Validation/acceptance:** Trigger from a message creates one traceable job; failure is visible; cancellation releases resources; result attaches to the correct chat; no duplicate generation on retry.
- **Risks/unresolved:** Prompt extraction may be unreliable; image generation may contend with chat inference; attachment identity depends on `CWP-10`.
- **Owner authorization:** **Implementation authorization required.** Modality direction does not authorize this path.

### CWP-19 — Creative studio shared foundation

- **Sources:** `REC-C-01`, `REC-C-04`, `REC-C-05`, `REC-C-07`, `REC-C-08`, `ODR-03`, `UD-05`
- **Current evidence/gap:** Separate image/video/audio surfaces exist; no unified studio system or complete media contract. Music and 3D mesh/texture generation are missing. Video lacks a complete catalogue/sidecar/delete/attachment contract. `UD-05` leaves modality priority unresolved.
- **Prerequisites:** `CWP-10` media identity, `CWP-15` resource authority, `AD-11`, owner answer to `UD-05` before any modality.
- **Dependencies/components:** S8 image, S9 video, S10/S11 audio, S12 music (planned), S17 library/export, shared catalogue and job/media contracts.
- **Required investigation:** Define shared job, asset, provenance, delete/export, and workspace contracts; inventory modality-specific gaps; identify what can be shared without merging chat and studio UI.
- **Expected deliverable:** A studio foundation specification and modality-readiness matrix. It authorizes no modality and creates no unified system by implication.
- **Validation/acceptance:** Shared asset identity across surfaces; delete/export consistency; unavailable media behavior; no cross-modality data loss. Modality acceptance requires separate owner authorization and evidence.
- **Risks/unresolved:** `UD-05` unanswered; 3D and music have no production subsystem; studio scope could absorb unrelated chat work; licensing/asset ownership unresolved.
- **Owner authorization:** **Direction approved; modality implementation not authorized.**

### CWP-20 — Central network policy, destination registry, consent, and audit

- **Sources:** `NET-01`, `NET-04`, `NET-07`, `NET-08`, `NET-09`, `NET-14`, `NET-18`, `NET-19`, `NET-20`
- **Current evidence/gap:** No destination registry, scoped permission state, or connection audit exists. Local-Only Mode is renderer-only. No backend/sidecar global offline enforcement was found. Proposed paths `P-01`–`P-08` are not authorized.
- **Prerequisites:** `INV-07` complete path inventory; `AD-08`; owner authorization for enforcement design.
- **Dependencies/components:** Go HTTP/dial paths, sidecars, Wails open/download actions, renderer, settings, audit storage, update path.
- **Required investigation:** Enumerate every caller, destination, trigger, payload class, credential, cache, retry, and failure path; define a single policy decision point that all egress paths must use; design fail-closed behavior and human-readable audit.
- **Expected deliverable:** A network control-plane design including destination/purpose registry, six-question consent, deny-by-default policy, offline enforcement, audit format, retention/clear, and change control. It must include immediate in-flight cancellation/refusal when offline is enabled and a renderer/webview navigation guard for unregistered external navigation. No destination is approved by the design.
- **Validation/acceptance:** Boundary-level positive and negative tests; every egress path attempts policy evaluation; unknown destination is denied; offline blocks egress and cancels in-flight work without breaking local work; webview navigation is guarded; audit records match actual connections; consent refusal has a defined effect.
- **Risks/unresolved:** Go/sidecar coverage may require architectural change; Wails/browser actions may bypass a single chokepoint; audit can leak sensitive data; update path must remain absent.
- **Owner authorization:** **Implementation authorization required.** No network experiment or enablement is authorized.

### CWP-21 — Civitai and automatic model-lookup remediation

- **Sources:** `NET-03`, `NET-10`, `NET-11`, `NET-13`
- **Current evidence/gap:** Image scanning calls Civitai automatically with a model hash; the name is cached in the manifest without clear/refresh control; errors are ignored. This contradicts no-background-lookup and user-initiated external information requirements.
- **Prerequisites:** `CWP-20` control plane; owner authorization for any retained lookup feature.
- **Dependencies/components:** Image scanning/registry, metadata cache, Go HTTP client, settings/consent UI, logs.
- **Required investigation:** Trace exact trigger and payload; determine whether lookup is required for local operation; define local-only fallback; design explicit user initiation, disclosure, cache clear/refresh, and visible failure.
- **Expected deliverable:** A remediation design that either removes automatic lookup or places it behind the control plane with purpose-specific consent. It must preserve local image functionality without the lookup.
- **Validation/acceptance:** Offline scan performs no request; no hash leaves the machine without consent; refusal leaves local metadata usable; cache is inspectable/clearable; failure is visible and not silently substituted.
- **Risks/unresolved:** Existing cached names may be stale; users may expect automatic enrichment. `NET-02` is owned by `CWP-22`; Civitai metadata lookup is not cloud inference.
- **Owner authorization:** **Required before retaining any Civitai path.** Removal is also a product change requiring authorization.

### CWP-22 — Remote ASR/TTS boundary and owner decision

- **Sources:** `NET-02`, `NET-05`, `NET-16`, `NET-17`
- **Current evidence/gap:** DashScope ASR and a remappable MeloTTS path are exposed without a separate owner decision. HTTPS exists for one path; an arbitrary remote HTTP URL is accepted elsewhere. Credential logging/redaction is unproven.
- **Prerequisites:** `CWP-20` control plane; `INV-09`; owner authorization (`OA-05`); `AD-09`.
- **Dependencies/components:** S10 ASR, S11 TTS, audio settings, Go HTTP clients, secrets/env handling, consent/audit.
- **Required investigation:** Determine whether remote audio is a product requirement; compare local-only alternatives; trace credentials, payloads, retention, and failure behavior; verify TLS and downgrade handling.
- **Expected deliverable:** A decision brief with a local-only baseline, an optional secured remote design, and an explicit recommendation. No remote path is enabled by the brief.
- **Validation/acceptance:** If remote is retained: authenticated TLS only, no silent downgrade, purpose-specific consent, credential redaction, audit, local fallback, and negative tests. If removed: local workflows remain usable and no remote code path is reachable.
- **Risks/unresolved:** `NET-02` cannot be satisfied while any cloud inference is enabled without an owner exception; remote audio may be the only usable quality path on some hardware; secrets handling is unproven.
- **Owner authorization:** **Owner decision required before any retention, removal, or enablement.** No model or network test is authorized.

### CWP-23 — External path governance: browser, downloads, indexing, error reports, updates

- **Sources:** `NET-05`, `NET-06`, `NET-10`, `NET-11`, `NET-13`, `NET-16`, `NET-17`; `V-NET-01`–`V-NET-28`
- **Current evidence/gap:** Browser-open, download, and organize/track actions are not fully separated. Error-report payload/redaction is unproven. Dependency downloads are user-initiated but not governed by a registry. No update checker exists and it must remain absent.
- **Prerequisites:** `CWP-20`; `INV-07`, `INV-08`; owner authorization for any download/indexing feature.
- **Dependencies/components:** Wails/browser open actions, dependency download UI, model library/indexer, error reporter, update path, audit.
- **Required investigation:** Enumerate all open/download/index/error/update actions; define separate permissions, destinations, payloads, retention, and failure; verify error context redaction; prove update absence.
- **Expected deliverable:** A path-specific governance matrix and consent copy grounded in the six questions. It must treat “open page,” “download file,” and “let Howl index it” as three separate decisions. Downloads use a temporary `.part` file, verify length, and verify a published checksum/signature **where available**; if unavailable, mark the result visibly unverified and non-authoritative. Reject truncation and never substitute a mirror. Browser handoff is recorded and never counted as a download or indexing grant.
- **Validation/acceptance:** Refusing one permission does not grant another; download does not imply indexing; browser open does not imply download; `.part`, length, checksum/signature, truncation, and mirror-failure branches are tested; error report shows exact payload and destination; no update request occurs.
- **Risks/unresolved:** Dependency download may be required to use local features; error context may contain character/chat content; a future signed updater would require a new decision.
- **Owner authorization:** **Required for any new path or retained external behavior.** No experiment is authorized.

### CWP-24 — Network boundary verification and negative-test harness

- **Sources:** `NET-15`; `V-NET-01`–`V-NET-28`
- **Current evidence/gap:** No release acceptance test or runtime network evidence exists. The 28 verification tests are defined but none has been executed for this plan.
- **Prerequisites:** `CWP-20` through `CWP-23` designs; `OA-09`; authorized runtime environment.
- **Dependencies/components:** Egress policy, audit, consent, offline control, local ASR/TTS/image/chat, test harness, packet/log capture.
- **Required investigation:** Map each `V-NET` test individually using §10.6 and its stated boundary; define fixtures for no-consent, refusal, expired permission, offline, destination change, credential exposure, redirect, oversized response, and silent fallback.
- **Expected deliverable:** An executable test specification and evidence ledger covering all 28 `V-NET` rows individually, including redirect and oversized-response harness branches. It must include negative tests that prove denial, not only positive tests that prove a request can be made.
- **Validation/acceptance:** Every test records command/environment, expected result, observed result, and evidence artifact; unknown remains unknown; a passing local request does not prove global enforcement.
- **Risks/unresolved:** Packet capture/tooling may be unavailable; sidecar paths may be hard to observe; test accounts/services may be unavailable. Failures must be reported, not skipped silently.
- **Owner authorization:** **Explicit authorization required before running network experiments or loading services.**

### CWP-25 — Android mobile remote gate

- **Sources:** `NET-12`, `UD-02`, `REC-M-09`
- **Current evidence/gap:** No server, pairing, transport, authentication, encryption, revocation, or off-by-default control exists. Android remote is conditional and unimplemented. `REC-M-09` preserves desktop completeness as a constraint.
- **Prerequisites:** Owner answer to `UD-02`; `OA-06`; `AD-10`; `CWP-20` network control plane.
- **Dependencies/components:** S23 (absent), network boundary, auth/session, QR/pairing, transport, device management, audit.
- **Required investigation:** Define threat model, pairing, authentication, transport encryption, revocation, local-network discovery, and offline behavior; determine whether desktop-only operation is the correct default permanently.
- **Expected deliverable:** A go/no-go decision brief and, only if authorized, a security architecture. No server or pairing implementation is created by this plan.
- **Validation/acceptance:** If authorized: off by default, explicit pairing, authenticated/encrypted transport, revocation, no silent LAN exposure, negative tests for unauthorized devices. If not authorized: no remote listener exists.
- **Risks/unresolved:** `UD-02` unresolved; remote access expands the attack surface; Android work could distract from desktop completeness (`REC-M-09`); no runtime evidence exists.
- **Owner authorization:** **Owner decision and separate implementation authorization required. Not authorized now.**

### CWP-26 — Optional persona / AI-as-user feasibility

- **Sources:** `REC-CH-08`
- **Current evidence/gap:** The requirement is unresolved. No design authority establishes whether the AI may author user-role messages, how attribution works, or what safeguards apply.
- **Prerequisites:** `CWP-01` identity, `CWP-04` character isolation; `AD-12`; owner decision.
- **Dependencies/components:** Chat schema, message roles, character/player attribution, UI, export.
- **Required investigation:** Define the user problem, authorship/attribution, consent, edit/delete behavior, and interaction with character isolation. Determine whether the feature is needed at all.
- **Expected deliverable:** A feasibility and risk brief, not an implementation. It must preserve the distinction between AI-generated character content and AI-authored user content.
- **Validation/acceptance:** If pursued: explicit attribution, no silent user-role fabrication, reversible edits, export clarity. If rejected: requirement remains unresolved/excluded by owner decision, not silently dropped.
- **Risks/unresolved:** Authorship confusion, user trust, export semantics, and interaction with persona play. No owner decision exists.
- **Owner authorization:** **Owner decision required before any design or implementation.**

### CWP-27 — Import/export provenance and round-trip integrity

- **Sources:** `REC-P-06`
- **Current evidence/gap:** Import work exists; provenance coverage is incomplete. Characters, chats, images, and videos are not fully cut over; export is category-limited. Import is a write path for untrusted character cards, JSON, PNG metadata, and possible archives; no bounded parsing or containment responsibility is documented.
- **Prerequisites:** `CWP-01` durable chat schema; `CWP-10` media identity; F3-C/D-STORAGE constraints.
- **Dependencies/components:** S17 libraries/import-export, content registry, chat/character/media stores, wipe/restore.
- **Required investigation:** Enumerate every imported field and its origin; define unknown-field preservation, conflict handling, and round-trip expectations; identify destructive normalization; inventory parser entry points and archive handling.
- **Expected deliverable:** A provenance and round-trip contract plus an import-hardening contract: bounded read sizes, decompression limits, path containment for extracted entries, safe embedded-metadata handling, no in-place rewrite of external files, and malformed-input failure that preserves the source.
- **Validation/acceptance:** Import → export → import preserves declared fields; missing media remains explicit; unknown fields are preserved or visibly rejected; provenance identifies source; hostile fixtures for malformed, oversized, traversal, archive-bomb, and unsafe-reference inputs fail closed without source mutation. Runtime proof required.
- **Risks/unresolved:** Legacy data may lack provenance; export scope may expand beyond evidence; normalization can silently destroy data; parser hardening is not a validated security guarantee until tested.
- **Owner authorization:** **Implementation authorization required.** No import/export operation is performed here.

### CWP-28 — Story Tools scope investigation

- **Sources:** `UD-03`
- **Current evidence/gap:** Story Tools page/scripts exist, but committed scope and ownership are unknown. It is not safe to treat existing controls as a product requirement.
- **Prerequisites:** `INV-11`; owner answer to `UD-03`.
- **Dependencies/components:** S20 story tools, memory/export bindings, chat/character stores; RAG only if `OA-10` is authorized.
- **Required investigation:** Inventory actual consumers and workflows; identify which are obsolete, duplicate, or require chat/memory foundations; define a bounded scope proposal.
- **Expected deliverable:** A scope decision brief with keep/retire/defer options. No page is removed or extended by this plan.
- **Validation/acceptance:** Each proposed capability maps to an owner need and a requirement ID; no orphan controls are promoted to requirements.
- **Risks/unresolved:** Scope may overlap memory/RAG; removing controls could break undocumented workflows; no usage evidence exists.
- **Owner authorization:** **Owner decision required before any implementation or removal.**

---

## 9. Validation, acceptance, and regression requirements

### 9.1 Evidence classes

| Class | Meaning | Cannot establish |
|---|---|---|
| Static source trace | Paths, bindings, schema, and absence are directly inspected | Runtime activation, correctness, or benefit |
| Cited test record | A prior test/verification record reports a result | That the test still passes on the current tree |
| Runtime observation | An authorized run produced an artifact | Generalization to other models/backends/hardware |
| Benchmark | A declared comparison with measured cost/benefit | Product acceptance without owner review |
| Owner acceptance | The owner accepts a defined scope | That unrelated scopes are accepted |

### 9.2 Validation rules

1. Every acceptance claim names its scope: model, backend, quantization, hardware, content type, or destination.
2. Every negative requirement has a negative test. “No request was observed” is not proven by a successful local workflow alone.
3. A control that is present but not reached is **Wired: No/Partial**, not operational.
4. A runtime success without a comparator does not establish benefit.
5. A test that cannot run is reported as blocked, not passed.
6. A failed check is not deleted, skipped, or weakened to make a package look complete.
7. No package is accepted by this plan. Acceptance requires a separate owner record after evidence review.

### 9.3 Regression domains

| Domain | Required regression focus |
|---|---|
| Chat identity | create/rename/move/import/export/delete; missing root; alias rejection; restart |
| Prompt path | exactly-once character, lore, memory, attachment, and rollover injection; RAG only if `OA-10` is authorized |
| Context | tokenizer budget, overflow, rollover boundary, multi-character flush, single-character continuation |
| Model lifecycle | load/ready/busy/swap/fail/rollback/unload; active stream; image/TTS contention |
| Capability | activation vs fallback; target support; conflict combinations; quality and resource measures |
| Media | available/missing/moved/deleted; display vs model input; delete/export/restore |
| Studio | shared asset identity; delete/export consistency; unavailable media |
| Import/export | field provenance; unknown fields; unavailable media; round-trip; malformed/oversized/traversal/archive-bomb inputs; source preservation |
| Durable state | backup/restore; atomic migration; crash consistency; free-space preflight; single writer; partial-write recovery across chat, memory, character, and world stores; media via `CWP-10` |
| RAG (conditional) | if `OA-10` is authorized: exactly-once retrieval, embedding identity, index version, delete/rebuild, resource admission |
| Network | denial without consent; offline and in-flight cancellation; unknown destination; destination change; credential redaction; audit accuracy; `.part`/checksum/truncation; webview guard |
| UI | character isolation; streaming/finalized parity; narration attribution; no ghost/duplicate controls |

### 9.4 Validation gates for work packages

A package may report only one of: **investigation complete**, **design proposed**, **implementation evidence present**, **validation blocked**, **validation failed**, or **owner acceptance recorded**. “Done” is not a valid status. No package may self-accept.

---

## 10. Network and privacy enforcement

### 10.1 Required boundary

The product does not provide cloud model inference or cloud model-provider support (`NET-02`). Permitted network functionality remains explicitly scoped to:

- Android mobile remote, only if separately authorized and secured (`NET-12`);
- user-authorized model metadata or compatibility lookups;
- user-authorized model downloads or browser handoffs;
- other external paths only after documentation, registry entry, purpose-specific consent, and owner authorization (`NET-19`).

The current renderer-only Local-Only Mode is **not** an application-wide network boundary. Go `net/http`, `net.Dial`, sidecars, proxies, Wails-mediated browser/download actions, and future paths must be subject to the same policy decision point.

### 10.2 Required control plane

| Control | Requirement | Evidence needed |
|---|---|---|
| Destination registry | Every destination has owner, purpose, feature, data classes, retention, and authorization state (`NET-04`) | Registry entry and policy evaluation for every path |
| Purpose-specific permission | No implied download, refresh, indexing, or secondary use (`NET-07`) | Denial tests and scoped permission state |
| Six-question consent | Destination, purpose, data, recurrence, retention, refusal effect (`NET-05`) | Human-readable copy plus actual payload proof |
| Separate actions | Browser open, download, and organize/track are distinct (`NET-06`) | Independent grant/refusal tests |
| Fail-closed enforcement | Unknown/ambiguous destination is denied (`NET-08`) | Negative tests at each boundary |
| Offline enforcement | Egress blocked without breaking local workflows (`NET-18`) | Egress-blocked runtime test |
| Audit | Human-readable record of actual connections (`NET-09`) | Audit-to-connection comparison |
| Cache/retention | Disclosed, inspectable, clearable (`NET-10`) | Cache inventory and clear/refresh test |
| Visible failure | No silent fallback or substituted result (`NET-11`) | Failure injection |
| Change control | New destination/purpose requires documented decision (`NET-14`) | Registry diff and approval record |
| Credentials | Never logged or sent to undocumented destinations (`NET-17`) | Redaction and destination tests |
| Updates | No automatic check; future check explicit and signed (`NET-20`) | Absence test now; signed design later |

### 10.3 External path register (N-01–N-14)

These are the inspected production external-path identities. No row is authorized by this plan. “Owner” is the accountable package for policy/design; it is not approval.

| ID | Destination / action | Trigger and data | Current compliance | Owning package | Disposition |
|---|---|---|---|---|---|
| N-01 | Civitai model-version lookup by SHA-256 | Image scan with a generic name and no cached name; sends model hash and ordinary HTTP metadata | **Non-compliant** — automatic, no prompt or audit | CWP-21 | Pending removal or purpose-specific consent; unauthorized now |
| N-02 | Hugging Face dependency file download | Missing-dependency acceptance; sends file path request and receives weights | **Partially compliant** — generic consent, `.part` cleanup, no checksum proof | CWP-23 | Pending registry, consent, validation; unauthorized now |
| N-03 | GitHub latest whisper release and asset | User-invoked whisper acquisition; sends API request, receives release ZIP | **Partially compliant** — no six-question pre-request disclosure | CWP-23 | Pending registry and consent; unauthorized now |
| N-04 | `git clone` whisper source | User-invoked Unix build acquisition; sends Git/HTTP metadata | **Partially compliant** — no separate disclosure or audit | CWP-23 | Pending registry and consent; unauthorized now |
| N-05 | Howl error-report endpoint | User enables reporting and sends report; includes version, OS, CPU count, errors, arbitrary `Context` | **Partially compliant** — default off, preview exists; payload/retention unproven | CWP-23 | Pending payload proof and consent; unauthorized now |
| N-06 | DashScope remote ASR transcription | Selected remote engine; sends API key, model ID, local audio path, options | **Non-compliant** — no operation-specific consent or audit | CWP-22 | Pending owner decision; unauthorized now |
| N-07 | DashScope ASR task polling | Repeated status polling with API key and task ID | **Non-compliant** — no consent event or visible cancellation | CWP-22 | Pending owner decision; unauthorized now |
| N-08 | DashScope result-URL fetch | Provider-issued result URL may be a CDN host | **Non-compliant** — returned host is not registry-controlled | CWP-22 | Pending owner decision and destination proof; unauthorized now |
| N-09 | Remote MeloTTS synthesize endpoint | Non-loopback configured URL; sends text, language, voice, speed, optional token | **Non-compliant when remote** — plaintext URL and no per-use consent | CWP-22 | Pending owner decision; unauthorized now |
| N-10 | TCP dial to configured MeloTTS host | Capability/probe refresh; DNS lookup and TCP attempt | **Non-compliant when remote** — probe occurs outside consent | CWP-22 | Pending owner decision; unauthorized now |
| N-11 | System-browser dependency source link | **Get it here**; browser sends visit/IP/cookies | **Partially compliant** — no pre-open warning or audit | CWP-23 | Pending browser-action disclosure; unauthorized now |
| N-12 | Tutorial publisher links | User clicks a tutorial link; webview/browser sends visit metadata | **Partially compliant** — no registry or warning | CWP-23 | Pending navigation guard and disclosure; unauthorized now |
| N-13 | Environment proxy route | Any default-transport client may traverse `HTTP_PROXY`/`HTTPS_PROXY` | **Unresolved** — proxy is outside Howl’s visibility | CWP-20 | Pending enforcement/audit design; unauthorized now |
| N-14 | Stored external source URLs | Display/copy only; no request found | **Not a path** — must not be counted as egress | CWP-23 | Display-only; preserve |

The 14 rows above include one non-path (`N-14`). The release inventory must still account for it as display-only evidence.

### 10.4 Local/loopback path register (L-01–L-05)

| ID | Local path | Boundary requirement | Owning package | Disposition |
|---|---|---|---|---|
| L-01 | llama.cpp chat server at `127.0.0.1:8080` | Loopback bind only; no silent LAN service | CWP-20 | Pending loopback verification |
| L-02 | llama.cpp ASR server on ephemeral loopback port | Ephemeral bind, lifecycle, and port ownership | CWP-20 | Pending lifecycle verification |
| L-03 | sd-server image backend at `localhost:<port>` | Local process only; port and process ownership | CWP-20 | Pending process verification |
| L-04 | MeloTTS default `http://127.0.0.1:8888` | Local default; any non-loopback URL becomes N-09 | CWP-22 | Pending boundary classification |
| L-05 | Wails/frontend asset origin | Local UI origin, not a cloud service | CWP-20 | Preserve as local |

### 10.5 Proposed paths (P-01–P-08)

These are **not implemented and not authorized**. They appear only so the registry does not silently absorb future work.

| ID | Proposed path | Required authorization and disclosure | Owning package | Disposition |
|---|---|---|---|---|
| P-01 | Model metadata/compatibility lookup | User-initiated; name publisher, hash sent, result, cache, recurrence; metadata ≠ download | CWP-20, CWP-21 | Unauthorized proposal |
| P-02 | General model download manager | Separate download permission; size, source, destination, checksum/signature, cancellation, no mirror | CWP-23 | Unauthorized proposal |
| P-03 | Application update check | Explicit/manual; no startup check; endpoint, sent/received data, signature, recurrence | CWP-20 | Unauthorized proposal |
| P-04 | Remote ASR provider | Separate owner decision; provider, audio/path data, retention, credential use | CWP-22 | Unauthorized proposal |
| P-05 | Remote TTS provider | Separate owner decision; text, voice, provider, retention, cache | CWP-22 | Unauthorized proposal |
| P-06 | Android remote access | Pairing, auth, encrypted transport, bind scope, timeout, revocation, no cloud relay | CWP-25 | Unauthorized proposal |
| P-07 | External publisher/source link | Browser receives visit/IP/cookies; separate from download/indexing | CWP-23 | Unauthorized proposal |
| P-08 | Provenance/licensing lookup | Metadata rules; no weight/user-content upload; cache and attribution disclosed | CWP-20, CWP-21 | Unauthorized proposal |

### 10.6 Individual V-NET verification dispositions

All 28 tests are **not run**. The state is not a pass. “App boundary” means the test must include Go transport/dialer, sidecar, webview, or browser-handoff evidence; a renderer mock alone is insufficient.

| ID | What it proves | Owning package | Prerequisites | Required evidence | State | Negative | App boundary |
|---|---|---|---|---|---|---|---|
| V-NET-01 | Static egress inventory is complete and mapped to registry entries or exclusions | CWP-20 | `INV-07`; `CWP-20` design | Inventory, registry mapping, build check | **Not run** | N | Yes |
| V-NET-02 | Startup, settings, scan, model list, and idle cause no outbound attempt | CWP-20 | `CWP-20`; egress-blocked runtime | Capture and startup log comparison | **Not run** | Y | Yes |
| V-NET-03 | Local chat, model load, image, ASR/TTS, scan, and settings work offline | CWP-20, CWP-24 | `CWP-20`; local fixtures | Egress-blocked end-to-end evidence | **Not run** | N | Yes |
| V-NET-04 | Offline mode blocks injected Go HTTP/dial and Wails network action | CWP-20 | Backend enforcement design | Fault-injected boundary evidence | **Not run** | Y | Yes |
| V-NET-05 | Every external path answers the six consent questions with actual data | CWP-20, CWP-23 | `CWP-20`; `CWP-23` copy | Prompt/payload comparison per path | **Not run** | N | Yes |
| V-NET-06 | Declining consent sends no request and preserves local work | CWP-20, CWP-23 | Consent implementation | Denial/continuity evidence | **Not run** | Y | Yes |
| V-NET-07 | Metadata permission cannot download, index, or authorize another publisher | CWP-20 | Scoped permission state | Cross-scope denial evidence | **Not run** | Y | Yes |
| V-NET-08 | Browser open, download, and indexing are separately recorded actions | CWP-23 | Action-separation design | Independent grant/refusal evidence | **Not run** | Y | Yes |
| V-NET-09 | Generic Civitai scan performs no request without permission and is auditable with permission | CWP-21 | `CWP-20`; `CWP-21` | Scan capture, consent, cache/audit evidence | **Not run** | Y | Yes |
| V-NET-10 | Dependency download handles cancellation, offline, truncation, HTTP error, and checksum mismatch without partial use | CWP-23 | `.part`/length/checksum design | Failure-injection and file-integrity evidence | **Not run** | Y | Yes |
| V-NET-11 | Remote ASR cannot run without a specific DashScope disclosure | CWP-22 | Owner decision; consent design | Prompt/payload/capture evidence | **Not run** | Y | Yes |
| V-NET-12 | Remote MeloTTS requires consent; plaintext HTTP rejected; secrets redacted; local TTS remains usable | CWP-22 | Owner decision; TLS/redaction design | Negative TLS/redaction/fallback evidence | **Not run** | Y | Yes |
| V-NET-13 | Allowed, denied, cancelled, and failed attempts appear in local audit without sensitive content | CWP-20 | Audit design | Audit-to-operation comparison | **Not run** | N | Yes |
| V-NET-14 | Cache clear re-asks; stale cache is not fresh; retention is visible | CWP-21, CWP-23 | Cache/retention design | Cache lifecycle evidence | **Not run** | Y | Yes |
| V-NET-15 | Every external failure is visible and never silently substitutes another provider or data | CWP-20, CWP-23 | Failure-behavior design | Failure-injection matrix | **Not run** | Y | Yes |
| V-NET-16 | Android remote is off by default, paired, authenticated, encrypted, scoped, expiring, revocable | CWP-25 | `UD-02`; `OA-06`; `AD-10` | Security/runtime evidence | **Not run** | Y | Yes |
| V-NET-17 | Error reporting is default-off, preview-accurate, destination/retention-visible, and not automatic | CWP-23 | `INV-08` | Payload, preview, consent, failure evidence | **Not run** | Y | Yes |
| V-NET-18 | Browser links warn about visit/IP; offline blocks or warns; action is separate | CWP-23 | Navigation guard | Browser/webview evidence | **Not run** | Y | Yes |
| V-NET-19 | No startup/periodic update request exists; manual check discloses endpoint/data/signature | CWP-20 | `CWP-20` design | Absence capture; future signed design | **Not run** | Y | Yes |
| V-NET-20 | Release build fails if any external host/dial path lacks registry entry or approved exclusion | CWP-20 | Registry completeness | Build-gate failure evidence | **Not run** | Y | Yes |
| V-NET-21 | With connectivity disabled, no external destination receives traffic at the application boundary | CWP-24 | `CWP-20`–`CWP-23`; authorized runtime | Capture, firewall/proxy, backend instrumentation | **Not run** | Y | Yes |
| V-NET-22 | Unauthorized generic scan performs no request; denial leaves visible “metadata unavailable” | CWP-21, CWP-24 | `CWP-21` | Boundary capture plus UI result | **Not run** | Y | Yes |
| V-NET-23 | Local model load causes no provider connection, DNS lookup, browser handoff, or download | CWP-20, CWP-24 | Local model fixture; egress block | DNS/connection/browser-handoff capture | **Not run** | Y | Yes |
| V-NET-24 | Failed/expired/revoked/malformed permission fails closed with no fallback route or cache | CWP-20, CWP-24 | Permission state machine | Negative state-transition evidence | **Not run** | Y | Yes |
| V-NET-25 | Browser handoff writes no file, downloads nothing, and indexes nothing | CWP-23 | Action separation | Filesystem/browser/index evidence | **Not run** | Y | Yes |
| V-NET-26 | Revocation stops subsequent requests, preserves local work, and requires a new decision | CWP-20 | Revocation design | Revoke-then-request evidence | **Not run** | Y | Yes |
| V-NET-27 | Android remote is separately gated; revocation does not alter desktop policy; remote cannot bypass desktop offline | CWP-25 | `UD-02`; `OA-06`; `CWP-20` | Desktop/remote isolation matrix | **Not run** | Y | Yes |
| V-NET-28 | Boundary harness covers Go dialer, sidecar, webview, and browser handoff, including redirect and oversized response | CWP-24 | `CWP-20`–`CWP-23` | Harness source, run log, failure proof | **Not run** | Y | Yes |

### 10.7 Addendum §23.8 control trace

These are existing addendum controls, not new product requirements.

| Control ID | Addendum §23.8 rule | Owning package | Verification | State |
|---|---|---|---|---|
| RULE-NET-01 | Downloads use a temporary `.part`, verify length, and verify a published checksum/signature **where available**; if no published value exists, mark the file visibly unverified and non-authoritative. Reject truncation and never fall back to a mirror. | CWP-23 | V-NET-10, V-NET-24 | Not implemented / not run |
| RULE-NET-02 | Enabling offline mode cancels or refuses in-flight requests immediately | CWP-20 | V-NET-21, V-NET-24 | Not implemented / not run |
| RULE-NET-03 | Renderer/webview navigation guard prevents unregistered external navigation and records browser handoff | CWP-20 | V-NET-18, V-NET-25 | Not implemented / not run |
| RULE-NET-04 | Failed or truncated downloads/results are never used or silently replaced by unverified data | CWP-23, CWP-24 | V-NET-10, V-NET-15 | Not implemented / not run |

### 10.8 Network completion criterion

`NET-15` is satisfied only when all 28 `V-NET` rows above have executed evidence in an authorized environment, every negative test proves denial at the application boundary, and no undocumented external connection occurs during local workflows. A local request succeeding is not completion. A skipped test is not a pass. A renderer-only mock is not application-boundary proof.

---

## 11. Release-readiness gates

These gates are release-readiness criteria, not acceptance records. A gate may be **open**, **blocked**, **failed**, or **evidence-ready for owner review**. It may not be marked adopted by this plan.

| Gate | Requirement | Current state | Governing packages |
|---|---|---|---|
| RG-01 | Protected F3-C/F3-D/F3-B/W8 records remain unchanged; no unauthorized acceptance | Preserved; F3-C Phase 0 open, D-ROOT open, D-STORAGE blocked | All; CWP-01 |
| RG-02 | Canonical chat, memory, character/world, and media identity is evidenced, including backup/restore, atomic migration, crash consistency, free-space preflight, and single-writer guarantees | Blocked by root/alias defects, `UD-01`, and unproven migration safety | CWP-01, CWP-10 |
| RG-03 | One authoritative prompt path with exactly-once contributors | Not established | CWP-02, CWP-06, CWP-07 |
| RG-04 | Token-accurate context, rollover, and character isolation are validated | Missing/unresolved | CWP-03, CWP-04, CWP-05 |
| RG-05 | Hot-swap is conversation-safe with rollback, cross-template re-render/re-tokenize/re-fit, and enforced resource authority | Partial/unsafe; cycle broken by §8.4 | CWP-08, CWP-15 |
| RG-06 | Existing capabilities have activation/support/correctness evidence for declared scopes | Mostly Implemented only | CWP-11, CWP-12, CWP-13, CWP-14 |
| RG-07 | Durable media and per-turn provenance are reconstructable | Missing/partial | CWP-09, CWP-10 |
| RG-08 | Chat presentation, narration, and image orchestration meet their contracts | Partial/missing/unresolved | CWP-16, CWP-17, CWP-18 |
| RG-09 | Studio foundation is bounded; no modality is implied | Direction only | CWP-19 |
| RG-10 | Network is deny-by-default, consented, audited, offline-proven, and tested individually through all 28 V-NET rows | Contradicted/missing; all 28 not run | CWP-20–CWP-24 |
| RG-11 | Remote ASR/TTS and Android have explicit owner decisions | Unresolved/not authorized | CWP-22, CWP-25 |
| RG-12 | Every requirement and package is bidirectionally traceable; no orphan or silent drop; all 28 V-NET tests, 27 path identities, and 20 network requirements are individually accounted for | Planned in §12; must be verified | All |
| RG-13 | If RAG is authorized by `OA-10`, retrieval demonstrably affects replies exactly once with embedding identity, index versioning, deletion/rebuild, retention, and resource admission | Blocked — no requirement identity or owner authorization | CWP-02, CWP-15 |

**Release statement:** No release-readiness claim is made. The plan is not adopted. The gates above are not passed.

---

## 12. Bidirectional traceability and final status

### 12.1 Forward traceability: requirement to disposition

Every one of the 66 original rows and every `NET-01`–`NET-20` row is listed below. The outcome is the reconciliation outcome, not a new judgment. A work package is the primary package that justifies action; where a requirement spans multiple packages, all are shown. “No work” entries preserve an exclusion, parked disposition, protected constraint, or already-satisfied scope. A requirement with no work package is not silently dropped.

| ID | Reconciliation outcome | Work package(s) | Disposition |
|---|---|---|---|
| `REC-UI-01` | **Partially satisfied** | `CWP-16` | Action required; implementation not authorized. |
| `REC-UI-02` | **Partially satisfied** | `CWP-16` | Action required; implementation not authorized. |
| `REC-UI-03` | **Partially satisfied** | `CWP-16` | Action required; implementation not authorized. |
| `REC-UI-04` | **Missing** | `CWP-17` | Action required; implementation not authorized. |
| `REC-UI-05` | **Missing** | `CWP-02`, `CWP-03` | Action required; implementation not authorized. |
| `REC-UI-06` | **Missing** | `CWP-18` | Action required; implementation not authorized. |
| `REC-CH-01` | **Missing** | `CWP-01` | Action required; implementation not authorized. |
| `REC-CH-02` | **Missing** | `CWP-01` | Action required; implementation not authorized. |
| `REC-CH-03` | **Partially satisfied** | `CWP-05` | Action required; implementation not authorized. |
| `REC-CH-04` | **Missing** | `CWP-04` | Action required; implementation not authorized. |
| `REC-CH-05` | **Unresolved** | `CWP-04` | Action required; implementation not authorized. |
| `REC-CH-06` | **Missing** | `CWP-04` | Action required; implementation not authorized. |
| `REC-CH-07` | **Partially satisfied** | `CWP-07` | Action required; implementation not authorized. |
| `REC-CH-08` | **Unresolved** | `CWP-26` | Action required; implementation not authorized. |
| `REC-CH-09` | **Partially satisfied** | `CWP-06` | Action required; implementation not authorized. |
| `REC-CH-10` | **Partially satisfied** | `CWP-06` | Action required; implementation not authorized. |
| `REC-CH-11` | **Unresolved** | `CWP-05`, `CWP-03` | Action required; implementation not authorized. |
| `ODR-08` | **Missing** | `CWP-04` | Action required; implementation not authorized. |
| `ODR-11` | **Partially satisfied** | `CWP-04`, `CWP-05` | Action required; implementation not authorized. |
| `ODR-04` | **Partially satisfied** | `CWP-06` | Action required; implementation not authorized. |
| `REC-M-01` | **Partially satisfied** | `CWP-11` | Action required; implementation not authorized. |
| `REC-M-02` | **Partially satisfied** | `CWP-02`, `CWP-03` | Action required; implementation not authorized. |
| `REC-M-03` | **Partially satisfied** | `CWP-13` | Action required; implementation not authorized. |
| `REC-M-04` | **Partially satisfied** | `CWP-13` | Action required; implementation not authorized. |
| `REC-M-05` | **Unresolved** | `CWP-13` | Action required; implementation not authorized. |
| `REC-M-06` | **Unresolved** | `CWP-13` | Action required; implementation not authorized. |
| `REC-M-07` | **Partially satisfied** | `CWP-15` | Action required; implementation not authorized. |
| `REC-M-08` | **Unresolved** | `CWP-15` | Action required; implementation not authorized. |
| `REC-M-09` | **Satisfied** (constraint preservation) | `CWP-25` | Preserve constraint; `CWP-25` is a dependency reference, not a work owner. |
| `REC-M-10` | **Satisfied** (protected scope) | — | No work — protected F3-C scope; preserve. |
| `REC-M-11` | **Missing** | `CWP-12` | Action required; implementation not authorized. |
| `REC-M-12` | **Missing** | `CWP-12` | Action required; implementation not authorized. |
| `REC-M-13` | **Partially satisfied** | `CWP-11` | Action required; implementation not authorized. |
| `REC-M-14` | **Missing** | `CWP-12` | Action required; implementation not authorized. |
| `REC-M-15` | **Missing** | `CWP-14` | Action required; implementation not authorized. |
| `REC-M-16` | **Partially satisfied** | `CWP-14` | Action required; implementation not authorized. |
| `REC-M-17` | **Missing** | `CWP-13` | Action required; implementation not authorized. |
| `REC-M-18` | **Partially satisfied** | `CWP-12` | Action required; implementation not authorized. |
| `ODR-02` | **Contradicted** | `CWP-08` | Action required; implementation not authorized. |
| `ODR-05` | **Missing** | `CWP-09` | Action required; implementation not authorized. |
| `ODR-06` | **Missing** | `CWP-09` | Action required; implementation not authorized. |
| `ODR-07` | **Missing** | `CWP-10` | Action required; implementation not authorized. |
| `ODR-09` | **Partially satisfied** | `CWP-16`, `CWP-17`, `CWP-18` | Action required; implementation not authorized. |
| `ODR-10` | **Partially satisfied** | `CWP-12`, `CWP-13`, `CWP-14` | Action required; implementation not authorized. |
| `UD-01` | **Unresolved** | `CWP-01` | Action required; implementation not authorized. |
| `UD-02` | **Unresolved** | `CWP-25` | Action required; implementation not authorized. |
| `UD-03` | **Unresolved** | `CWP-28` | Action required; implementation not authorized. |
| `UD-04` | **Unresolved** | `CWP-17` | Action required; implementation not authorized. |
| `UD-05` | **Unresolved** | `CWP-19` | Action required; implementation not authorized. |
| `UD-06` | **Unresolved** | `CWP-06` | Action required; implementation not authorized. |
| `REC-P-01` | **Missing** | `CWP-10` | Action required; implementation not authorized. |
| `REC-P-02` | **Missing** | `CWP-10` | Action required; implementation not authorized. |
| `REC-P-03` | **Partially satisfied** | `CWP-10` | Action required; implementation not authorized. |
| `REC-P-04` | **Missing** | `CWP-10` | Action required; implementation not authorized. |
| `REC-P-05` | **Satisfied** (image scope) | — | No work — image metadata satisfied; preserve as existing evidence. |
| `REC-P-06` | **Partially satisfied** | `CWP-27` | Action required; implementation not authorized. |
| `REC-C-01` | **Partially satisfied** | `CWP-19` | Action required; implementation not authorized. |
| `REC-C-02` | **Satisfied** (parked disposition) | — | No work — parked disposition; no promotion. |
| `REC-C-03` | **Satisfied** (parked disposition) | — | No work — parked disposition; no promotion. |
| `REC-C-04` | **Missing** | `CWP-19` | Action required; implementation not authorized. |
| `REC-C-05` | **Partially satisfied** | `CWP-19` | Action required; implementation not authorized. |
| `REC-C-06` | **Excluded** | — | No work — owner exclusion (desktop avatars). |
| `REC-C-07` | **Missing** | `CWP-19` | Action required; implementation not authorized. |
| `REC-C-08` | **Missing** | `CWP-19` | Action required; implementation not authorized. |
| `ODR-01` | **Excluded** | — | No work — owner exclusion (desktop avatars). |
| `ODR-03` | **Contradicted** | `CWP-19` | Action required; implementation not authorized. |
| `NET-01` | **Unresolved** | `CWP-20` | Action required; implementation not authorized. |
| `NET-02` | **Contradicted** | `CWP-22` | Action required; `CWP-20` is governance dependency; `CWP-22` is sole work owner. |
| `NET-03` | **Contradicted** | `CWP-21` | Action required; implementation not authorized. |
| `NET-04` | **Missing** | `CWP-20` | Action required; implementation not authorized. |
| `NET-05` | **Partially satisfied** | `CWP-22`, `CWP-23` | Action required; implementation not authorized. |
| `NET-06` | **Partially satisfied** | `CWP-23` | Action required; implementation not authorized. |
| `NET-07` | **Missing** | `CWP-20` | Action required; implementation not authorized. |
| `NET-08` | **Contradicted** | `CWP-20` | Action required; implementation not authorized. |
| `NET-09` | **Missing** | `CWP-20` | Action required; implementation not authorized. |
| `NET-10` | **Partially satisfied** | `CWP-21`, `CWP-23` | Action required; implementation not authorized. |
| `NET-11` | **Partially satisfied** | `CWP-21`, `CWP-23` | Action required; implementation not authorized. |
| `NET-12` | **Missing** | `CWP-25` | Action required; implementation not authorized. |
| `NET-13` | **Partially satisfied** | `CWP-21`, `CWP-23` | Action required; implementation not authorized. |
| `NET-14` | **Missing** | `CWP-20` | Action required; implementation not authorized. |
| `NET-15` | **Missing** | `CWP-24` | Action required; implementation not authorized. |
| `NET-16` | **Partially satisfied** | `CWP-22`, `CWP-23` | Action required; implementation not authorized. |
| `NET-17` | **Partially satisfied** | `CWP-22`, `CWP-23` | Action required; implementation not authorized. |
| `NET-18` | **Missing** | `CWP-20` | Action required; implementation not authorized. |
| `NET-19` | **Missing** | `CWP-20` | Action required; implementation not authorized. |
| `NET-20` | **Satisfied** (absence only) | `CWP-20` | Preserve absence; `CWP-20` is an absence-maintenance check, not a product action. |

### 12.2 Reverse traceability: work package to requirement

Every `CWP-*` package below has at least one exact source requirement. The source IDs are generated from the same forward map, so no package is orphaned and no requirement is linked only by prose.

| Work package | Exact source requirement IDs |
|---|---|
| `CWP-01` | `REC-CH-01`, `REC-CH-02`, `UD-01` |
| `CWP-02` | `REC-UI-05`, `REC-M-02` |
| `CWP-03` | `REC-UI-05`, `REC-CH-11`, `REC-M-02` |
| `CWP-04` | `REC-CH-04`, `REC-CH-05`, `REC-CH-06`, `ODR-08`, `ODR-11` |
| `CWP-05` | `REC-CH-03`, `REC-CH-11`, `ODR-11` |
| `CWP-06` | `REC-CH-09`, `REC-CH-10`, `ODR-04`, `UD-06` |
| `CWP-07` | `REC-CH-07` |
| `CWP-08` | `ODR-02` |
| `CWP-09` | `ODR-05`, `ODR-06` |
| `CWP-10` | `ODR-07`, `REC-P-01`, `REC-P-02`, `REC-P-03`, `REC-P-04` |
| `CWP-11` | `REC-M-01`, `REC-M-13` |
| `CWP-12` | `REC-M-11`, `REC-M-12`, `REC-M-14`, `REC-M-18`, `ODR-10` |
| `CWP-13` | `REC-M-03`, `REC-M-04`, `REC-M-05`, `REC-M-06`, `REC-M-17`, `ODR-10` |
| `CWP-14` | `REC-M-15`, `REC-M-16`, `ODR-10` |
| `CWP-15` | `REC-M-07`, `REC-M-08` |
| `CWP-16` | `REC-UI-01`, `REC-UI-02`, `REC-UI-03`, `ODR-09` |
| `CWP-17` | `REC-UI-04`, `ODR-09`, `UD-04` |
| `CWP-18` | `REC-UI-06`, `ODR-09` |
| `CWP-19` | `UD-05`, `REC-C-01`, `REC-C-04`, `REC-C-05`, `REC-C-07`, `REC-C-08`, `ODR-03` |
| `CWP-20` | `NET-01`, `NET-04`, `NET-07`, `NET-08`, `NET-09`, `NET-14`, `NET-18`, `NET-19`, `NET-20` |
| `CWP-21` | `NET-03`, `NET-10`, `NET-11`, `NET-13` |
| `CWP-22` | `NET-02`, `NET-05`, `NET-16`, `NET-17` |
| `CWP-23` | `NET-05`, `NET-06`, `NET-10`, `NET-11`, `NET-13`, `NET-16`, `NET-17` |
| `CWP-24` | `NET-15` |
| `CWP-25` | `REC-M-09`, `UD-02`, `NET-12` |
| `CWP-26` | `REC-CH-08` |
| `CWP-27` | `REC-P-06` |
| `CWP-28` | `UD-03` |

`CWP-23` and `CWP-24` consume the 28 `V-NET` verification tests. Each is individually mapped in §10.6; `CWP-24` owns the harness. They are verification controls, not additional requirement rows, and are therefore not counted in the 86-row forward matrix. The 14 `N-*`, five `L-*`, and eight `P-*` rows are individually listed in §10.3–§10.5.

### 12.3 Final status

| Item | Status |
|---|---|
| Corrected master plan | **Proposed — not adopted** |
| Original reconciliation rows consumed | **66 of 66** |
| Network requirements consumed | **20 of 20** |
| Existing-capability register | **Complete 37-row reproduction in §4.3** |
| Requirement outcomes preserved | **Yes — six-state reconciliation language retained** |
| Capability evidence separated | **Yes — seven-state vector retained; no promotions claimed** |
| Major architectural missions | **Present as independent packages** |
| Work packages | **28 (`CWP-01`–`CWP-28`)** |
| Network path identities | **27 (`N-01`–`N-14`, `L-01`–`L-05`, `P-01`–`P-08`), individually listed** |
| Network verification tests | **28 (`V-NET-01`–`V-NET-28`), individually mapped, all not run** |
| MUSE findings | **12 of 12 dispositioned; see correction disposition record** |
| Bidirectional traceability | **Generated and checked; see updated verification record** |
| Implementation estimates | **None** |
| Runtime or benchmark claims added | **None** |
| Application source/tests/config/migrations changed | **None** |
| Capabilities enabled/disabled | **None** |
| Protected artifacts or gates changed | **None** |
| Implementation authorized | **No** |
| Plan adopted or accepted | **No** |

### 12.4 Closing statement

This corrected plan is an engineering roadmap, not an implementation authorization. It preserves the owner-characterized recovery baseline, the 66 original requirement rows, the 20 network requirements, the seven-state capability evidence model, the six reconciliation outcomes, the owner exclusions and unresolved questions, and the protected-boundary constraints.

It deliberately leaves the following unresolved rather than inventing answers:

- whether model selection is global or per-chat (`UD-01`);
- whether Android remote is in scope (`UD-02`);
- the committed Story Tools scope (`UD-03`);
- narration semantics (`UD-04`);
- the next studio modality, if any (`UD-05`);
- the memory review policy (`UD-06`);
- whether remote ASR/TTS remain product paths;
- whether any proposed external destination is authorized.

No work package may be treated as accepted, closed, or authorized because it appears here. No runtime behavior, benchmark result, compatibility claim, or network guarantee has been added. The next legitimate step is owner review of this plan and authorization of specific investigations, not implementation.
