# HOWL CHAT — REQUIREMENTS RECOVERY AND OWNER-DECISION ADDENDUM

**Date:** 2026-11-11  
**Mission type:** Read-only historical recovery, requirements reconciliation, and architecture analysis  
**Status:** **RECOVERY BASELINE FOR THE CORRECTED MASTER PLAN — OWNER-CHARACTERIZED 2026-11-11; NOT THE CORRECTED MASTER PLAN, NOT A CONTRACT, NOT A FREEZE, NOT AN ACCEPTANCE, NOT AN IMPLEMENTATION AUTHORIZATION, NOT R6, NOT AN F3-C AMENDMENT, NOT F3-D**  
**Implementation authorization:** NONE  
**Architecture-freeze authorization:** NONE  
**Build / test / model-load / migration authorization:** NONE  
**Capability-mission status:** Research, evidence, compatibility, selection, validation, and implementation authorization remain separate stages. This baseline does not enable any technique, backend feature, or optimization.  
**Network mission status:** Network boundaries, external integrations, and user consent are a first-class product and architecture mission. Local inference and creative work remain offline; every external destination and purpose requires documented, purpose-specific authorization. §23 records the existing and proposed network-path inventory; no external path is validated or authorized by this document.  

**Authority for requirements:** the owner decisions in the 2026-11-11 recovery instruction, plus the previously accepted records identified in §1.3. This addendum does not replace the authoritative ledger or any accepted decision.  
**Authority for F3-C:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` remains controlling. This addendum creates no amendment to it.  
**Purpose:** recover what Howl is actually intended to become before a corrected master-plan revision is prepared. This artifact intentionally does **not** prepare that revision.

---

## 1. Artifact identity and evidence sources

### 1.1 Principal artifact hashes

| Artifact | SHA-256 | Role |
|---|---|---|
| `PROJECT_ISSUES_AND_EXPANSION.md` | `8156318ec8855ddbab1f218927352e59de6d60e23a3c81f11f85a56c61213491` | Authoritative historical requirements and expansion ledger |
| `STAGE_STATUS.md` | `33a16fdbca509fc4407bf840deee1dd622a0fe5a1bd272ac4c1ca1179dbe9a62` | Historical implementation/verification record for T1/T2 and chat UI observations |
| `docs/tier4-avatar.md` | `be04f0253e61c0fc28c245edbd8eb7f3ad753b8bdf152ceb8308cace23be04c9` | Parked reactive-3D-avatar concept; now governed by the explicit desktop-avatar exclusion |
| `docs/review/INVESTIGATION-character-voice-architecture-2026-11-05.md` | `0236e2300ef4cd704a946035bce5c943a407e33d3bf589214dc6927a8b85ea6e` | Chat bubble, character identity, voice, and message-ownership investigation |
| `docs/review/HOWL-COMPLETE-PRODUCT-ARCHITECTURE-REQUIREMENTS-DEPENDENCY-MASTER-PLAN-2026-11-11.md` | `e39cc396fa20adf2656a1b4e650a8bdc6c9a49cb46c3ddf95309ff144229ef8c` | Prior master plan under correction; unmodified |
| `docs/review/REVIEW-INDEPENDENT-HOWL-COMPLETE-PRODUCT-ARCHITECTURE-MASTER-PLAN-2026-11-11.md` | `af493c0b6118ebede556d4927e8ab93b9378ce46c8298f3bb9e38554f6265e9d` | Independent Muse review, findings F-01–F-15; unmodified |
| `docs/review/ADJUDICATION-MATRIX-HOWL-MASTER-PLAN-MUSE-REVIEW-2026-11-11.md` | `a9fdaeefdc4ea57c28a706cc09e57a77566f3ed6150d2e4e305904de7b6049ed` | Prior adjudication and COR-01–COR-10; unmodified |
| `docs/review/OWNER-DECISION-PACKET-HOWL-MASTER-PLAN-2026-11-11.md` | `6721d770694cda2c7b132fdc15bbd8d19dd6de0502c4f3e5f3ca47ce620ff4f3` | Prior owner packet; unmodified |
| `docs/review/ENGINEERING-F2-STORAGE-MODEL-INTELLIGENCE-2026-11-09.md` | `0f40221f544d04365b784f103d6befc4230b7e2f1699f81d9002265169ac76d1` | F2 storage/model-intelligence engineering report; 767 lines |
| `docs/review/ENGINEERING-F2B-STORAGE-MODEL-INTEGRATION-2026-11-09.md` | `95e24ccdcdc033c8b050b0f435df8b3eee85dc2e2dd6dab7c3e0bd6d451bee4a` | F2-B storage/model-intelligence integration report; 161 lines |

At the start of this recovery, the master plan, Muse review, adjudication, and owner packet matched their recorded identities. The five protected F3-C artifact hashes are verified at the end of this addendum.

### 1.2 Source hierarchy used

This addendum distinguishes four source levels:

1. **Current owner decision.** The 2026-11-11 recovery instruction is authoritative where it speaks directly. It can supersede a historical status without rewriting the historical document.
2. **Accepted project record.** A frozen contract, acceptance record, protected boundary, or prior explicit owner decision remains controlling in its own scope.
3. **Historical requirement.** A ledger entry, stage requirement, parked design note, or investigation finding. Its original status is preserved even when the current status changes.
4. **Current implementation evidence.** Static source inspection only. It establishes what code exists, not that the product behavior is correct at runtime.

When these disagree, the addendum states the disagreement. It does not silently harmonize them.

### 1.3 Evidence classes

| Label | Meaning |
|---|---|
| **Owner-decided** | Explicit in the 2026-11-11 instruction or a previously accepted owner record |
| **Historically accepted** | Accepted in a recorded stage, gate, contract, freeze, or adjudication |
| **Source-confirmed** | Directly visible in inspected source |
| **Test-confirmed (reported)** | A cited project record reports a passing test; this mission did not rerun it |
| **Runtime-confirmed (reported)** | A cited project record reports a native/runtime observation; this mission did not reproduce it |
| **Inferred** | Follows from source structure or documented behavior but was not executed |
| **Unestablished** | No reliable source or evidence recovered |
| **Superseded** | A later record or current owner decision controls |
| **Proposed** | A technical or product proposal that has not been adopted |
| **Deferred** | A real requirement whose implementation is not authorized or scheduled |
| **Parked** | A design concept with no current stage or commitment |
| **Excluded** | Explicitly removed from product scope by the owner |

No claim in this addendum is runtime-confirmed unless a cited project record supplies that evidence. This mission did not run the application, a test, a build, a model, a migration, or a destructive operation.

### 1.4 Required sources recovered

The following sources were read directly:

- `PROJECT_ISSUES_AND_EXPANSION.md`: Tiers 0–4, E.1–E.8, the F2-A/F2-B/F2-C ledgers, F3-A/F3-A-R1 decisions, REG-1, FUT-C1–C5, FUT-CH1–CH3, FUT-CHUI1/CHUI2/CHTOK1/CHID1/CHAIU1, FUT-MEM1, FUT-LORE1, FUT-HW1/FUT-HW2, FUT-PROMPT1, HW-C1–C3, MOB-1–MOB-4, W1–W28c, and the F3-B/F3-C status records.
- `STAGE_STATUS.md`: chat-bubble border/halo backlog, T1.4 family/template work, group-chat verification, theme work, and the chat UI observations.
- `docs/tier4-avatar.md`: the complete parked reactive-3D-avatar design.
- `docs/review/INVESTIGATION-character-voice-architecture-2026-11-05.md`: message ownership, character IDs, bubble actions, and the absence of a per-character voice/overflow menu at that time.
- `docs/review/ENGINEERING-F2-STORAGE-MODEL-INTELLIGENCE-2026-11-09.md`: content/model-intelligence foundations, evidence classes, candidate runtimes, and explicit non-completion limits.
- `docs/review/ENGINEERING-F2B-STORAGE-MODEL-INTEGRATION-2026-11-09.md`: bounded integration, model-browser metadata, deferred tokenizer/template/capability negotiation, and the warning that mirror writes are not an ownership migration.
- The prior master plan, Muse review, adjudication, and owner packet.
- Current source for chat history, prompt assembly, model loading, model templates, memory, image generation, UI rendering, samplers, and runtime flags.

### 1.5 Sources not recovered

Two artifacts named by the owner were searched for by exact filename, likely variants, content phrase, and git history:

- `possible expansion of howl chat for review.txt`
- A separate original chat UI planning document

**Neither exact artifact is present in the accessible workspace or its git history.** No claim is made that they never existed outside this checkout. Their recoverable content survives in the ledger and `STAGE_STATUS`; §5 records exactly what was and was not recovered.

The two F2 engineering reports and the ledger are present, but they are **leads**, not proof that every original requirement survived. They are different artifacts from the missing UI plan and expansion text. The exact original document set and count are not established, so this addendum does not claim that “three” or “four” source documents have been recovered.

---

## 2. Executive summary

### 2.1 The corrected baseline in one paragraph

Howl remains a local-first, single-user Windows desktop creative application. Its product baseline is broader than the prior master plan's “current surface plus deferred options” framing. The owner has now fixed four major directions: desktop character avatars are **excluded**; active-conversation model hot-swap is **required**; a general local creative studio is an **approved direction** without immediate modality implementation; and a substantial application-owned memory system is **required**. Durable chat media is approved. Multi-character inference context must be flushed after each character response, while single-character context must not be flushed per turn. Per-turn provenance is required as a concept, but the exact record remains an engineering recommendation. Seed preservation is desired only where it is technically meaningful; no exact-reproduction guarantee is authorized.

Separately, inference capability is now a first-class architecture mission. Howl must inventory techniques across engines, map them by model architecture and backend/version/hardware, record interactions and conflicts, automate safe selection, preserve informed manual override, and benchmark results rather than infer success from parameter acceptance. The current pinned-runtime inventory is a first-pass seed, not a frozen catalogue. Conversation context rollover (capacity management) and multi-character context flushing (speaker isolation) are explicitly separate features.

Network boundaries are also a first-class mission. Howl is a local application that may occasionally use the internet with the user's permission; it is not a cloud-connected AI service. The current implementation does not satisfy that boundary. Image scans can call Civitai automatically with a model hash, remote DashScope ASR and a remappable MeloTTS endpoint are exposed without operation-specific consent, and Local-Only Mode only wraps renderer `fetch`/XHR while Go code remains outside it. §23 is the authoritative mission, path inventory, and acceptance test for this boundary.

### 2.2 The most important corrections

1. **The prior master plan is wrong to call model hot-swap “unaccepted.”** It is now a required product behavior. It is still **not** an F3-C target-set requirement and creates no F3-C amendment. It is also still distinct from REG-1 content hot-swap.
2. **The prior master plan is wrong to call a general creative studio “unestablished.”** The general direction is approved. F2-B5 remains the approved modality-workspace direction; E.6/E.8 remain parked; each new modality still needs its own scoped requirements.
3. **The prior master plan is wrong to retain desktop character avatars as an ordinary parked future.** They are explicitly excluded. Chat portrait thumbnails remain required and are a different feature.
4. **The prior owner packet's OD-01 is no longer wholly open.** Provenance, seed preservation, durable attachments, model switching, and context flushing are now owner-clarified. The per-chat continuation-model preference remains a separate unresolved question. No specific provenance field set or exact-reproduction promise is thereby approved.
5. **The memory requirement is materially larger than the ledger's prior “deferred feature” treatment.** FUT-MEM1 is now a required product subsystem, even though implementation remains unauthorized.
6. **Chat UI requirements were omitted or scattered.** Bubble spacing/borders, per-character colors, portrait thumbnails, narration, a hidden-by-default token counter, and message-triggered image generation are owner-confirmed requirements. The exact historical UI document has not been recovered.
7. **The implementation baseline does not meet the new requirements.** Current chat history stores no model, engine, template, sampler, or seed provenance; no per-chat model; no durable image/attachment identity; no token counter; no message-triggered image generation; no explicit per-character appearance state; no demonstrated per-character context flush; and no coordinated hot-swap/rollback protocol.
8. **Source presence is still not runtime correctness.** The existing model loader has admission checks and a health check, but it stops the old server before the replacement is ready, does not verify weight release beyond process exit, and does not coordinate an active stream. The current `Health()` check is not a model-identity/readiness proof.
9. **Inference optimization is not a backend footnote.** The prior plan treats it mainly as Flash Attention/offload exposure. The required mission is hardware-aware and architecture-aware: technique inventory, architecture mapping, backend/version support, conflict analysis, automatic selection, manual override, and benchmarks. No capability catalogue is frozen.
10. **Context rollover and multi-character flushing are being conflated.** Conversation context rollover manages capacity through shifting, summarization, retention, and rebuilding. Multi-character flushing isolates one character's active inference context from the next speaker. Neither is model hot-swap.
11. **The Author's Note requirement is more specific than a one-line feature row.** The recovered presentation and scope are a click-open flyout, floating or hidden presentation, optional side display, chat-only scope, and a blank note for a new chat, while an existing chat retains its own note.
12. **The local-first claim is not currently true at the network boundary.** The application has automatic Civitai hash lookup during image scans, remote ASR and remappable remote TTS paths, a frontend-only Local-Only Mode, no destination registry, no connection audit, and no backend egress enforcement. A local-first claim is a release property, not a label; §23 defines the evidence required to earn it.

### 2.3 What this addendum does not do

It does not:

- select the durable chat schema;
- approve a specific provenance field set;
- guarantee response reproduction from a seed;
- define an F3-C amendment;
- authorize R6 or F3-D;
- adopt the corrected master plan;
- implement any model switch, memory, studio modality, narration, UI, media, or persistence change;
- close W4 or reinterpret W8;
- reopen accepted F3-B/F3-C decisions;
- convert parked concepts into commitments.

### 2.4 Owner characterization of the recovery baseline

The owner characterized this revision as the new recovery baseline and identified six properties to preserve. The baseline records them as controlling design constraints for the corrected plan, not as implementation authorization.

| Preserved property | Baseline treatment |
|---|---|
| The capability catalogue remains unfrozen | The llama.cpp inventory is a first-pass seed; multi-engine/architecture research and benchmarks remain required |
| Selection is an architecture responsibility | Backend acceptance is not proof of correctness; Howl must reason about requirements, conflicts, fit, and benefit |
| Manual overrides have consequences | Overrides must expose conflicts/risks and cannot silently defeat resource protections |
| Benchmarks include output quality | Speed or memory gains that damage coherence, accuracy, or continuity are not automatically improvements |
| The two context mechanisms stay separate | Conversation rollover must not inherit multi-character isolation behavior |
| Author's Note persistence is scoped | Existing conversations retain their own notes; new conversations start blank |

The owner further distinguished the capability system Howl must possess from a mandate to enable every technique or implement every backend feature. The corrected master plan must consume this baseline systematically; research, compatibility evidence, selection design, validation, and implementation authorization are not interchangeable.

---

## 3. Explicit owner decisions

### 3.1 Decision register

| ID | Owner decision | Status | Consequence |
|---|---|---|---|
| ODR-01 | Desktop character avatars are explicitly excluded | **EXCLUDED** | No architecture, roadmap, implementation, or ordinary future item for a persistent desktop avatar. Chat portraits remain required. |
| ODR-02 | Active-conversation model hot-swap is required | **REQUIRED PRODUCT BEHAVIOR** | Must switch models in an existing conversation without losing history, without simultaneous-residency assumption, and without sending to the wrong model. Not an F3-C amendment. |
| ODR-03 | General local creative studio is an approved product direction | **APPROVED DIRECTION** | Music, 3D mesh, texture, and future local modalities are in the intended product. No modality is thereby implemented or release-committed. |
| ODR-04 | A substantial application-side memory system is required | **REQUIRED MAJOR CAPABILITY** | Inventory, possessions, clothing/appearance, facts, events, relationships, clips, retrieval, summarization, relevance, reinjection, and lifecycle are product requirements. |
| ODR-05 | Preserve the actual seed used when technically meaningful and beneficial | **OWNER DIRECTION; ENGINEERING RECOMMENDATION REQUIRED** | Research actual runtime availability and utility. Do not force manual random/explicit management. Do not promise reproduction. |
| ODR-06 | Per-turn provenance is required as a concept | **REQUIRED CONCEPT; RECORD SHAPE OPEN** | Record which model/configuration produced each response. Not context flushing. Exact fields remain a recommendation, not an owner-approved schema. |
| ODR-07 | Durable chat attachments/media are approved | **APPROVED** | Images/audio remain associated with their messages; display is distinct from model input; lifecycle behavior must be defined. |
| ODR-08 | Multi-character context is flushed after each character response; single-character context is not flushed per turn | **REQUIRED** | Active inference context isolation only. Must not delete history, lore, memory, or conversation state. |
| ODR-09 | The chat UI requirements in §4 are owner-confirmed | **CONFIRMED** | Bubble spacing/borders, per-character colors, portraits, narration, token counter, and message-triggered image generation are product requirements. |
| ODR-10 | Inference capability and optimization are a first-class, hardware-aware and architecture-aware architecture mission | **OWNER-DIRECTED MISSION** | Technique inventory across engines, architecture/backend mapping, conflict analysis, automatic selection, informed manual override, and benchmarks. No individual optimization or catalogue is thereby adopted. |
| ODR-11 | Conversation context rollover and multi-character context flushing are separate features | **CLARIFIED** | Rollover manages conversation capacity; multi-character flushing isolates active inference context after each character response. Neither is model hot-swap. |
| ODR-12 | Network boundaries, external integrations, and user consent are a first-class product and architecture requirement | **REQUIRED BOUNDARY** | Local inference and creative work remain offline; Android remote is explicit and secured; external model information and downloads require purpose-specific permission. No silent telemetry, background model lookup, automatic provider connection, or unexpected update check. Every proposed destination/purpose requires a documented decision. §23 governs the inventory and acceptance criteria; it does not authorize any external path. |

### 3.2 Boundary clarification: avatar exclusion vs 3D assets

The desktop-avatar exclusion is an exclusion of the **interactive desktop character surface** described in `docs/tier4-avatar.md`: a persistent 3D character, reactive emotion, TTS-driven mouth, and desktop interaction. It does **not** prohibit 3D mesh or texture generation as a creative-studio modality.

The distinction is:

| Item | Status |
|---|---|
| Chat portrait thumbnail | **Required** |
| Character card images / image sets | Existing product requirement; not an avatar surface |
| Persistent desktop 3D character | **Excluded** |
| Multi-character desktop avatars | **Excluded** |
| 3D mesh generation as a creative asset | **Approved direction**; separately scoped, not an avatar feature |
| Texture generation as a creative asset | **Approved direction**; separately scoped |
| Using generated 3D assets to reconstruct the excluded desktop avatar | **Not authorized by the studio direction** |

`docs/tier4-avatar.md` is retained as historical evidence. Its status is now **EXCLUDED / HISTORICAL**, not “parked pending foundations.” It must not be carried as an ordinary roadmap row.

### 3.3 Boundary clarification: content hot-swap vs model hot-swap

| Term | Meaning | Status |
|---|---|---|
| **Content hot-swap** | Activation/replacement of world, scenario, character, and lorebook context under REG-1 | **PROTECTED EXISTING BEHAVIOR** |
| **Model hot-swap** | Replacing the loaded inference model during an existing conversation | **REQUIRED PRODUCT BEHAVIOR**, separate from REG-1 |
| **Per-chat continuation-model preference** | A chat remembering which model should continue it | **STILL UNRESOLVED**; do not infer from hot-swap |
| **Per-turn model provenance** | Historical record of the model that produced a response | **REQUIRED CONCEPT**; exact record open |

REG-1 does not authorize, define, or protect model switching. The model-switch requirement does not reopen REG-1 or alter any F3-C contract.

### 3.4 Boundary clarification: seed and provenance

The owner has decided the **concept**, not the exact schema:

- The actual seed used should be preserved if the runtime exposes it and the record is useful.
- Explicit-versus-random seed behavior is not itself a product requirement; it follows research and engineering best practice.
- The likely direction is a per-turn record containing the actual seed, its origin/status, and relevant generation configuration.
- This is not approval of exact response reproduction. Model version, runtime version, hardware, template, context, samplers, kernel nondeterminism, and backend differences can prevent reproduction.
- Image-generation seed handling remains a separate media-generation concern.

---

## 4. Recovered historical requirements

### 4.1 Reading rule

The table below preserves the original requirement, its source, and its original status. “Current” is not a claim that the requirement was always accepted; it is the status after the 2026-11-11 owner instruction and the existing accepted records.

### 4.2 Chat UI and interaction requirements

| ID | Historical requirement | Source | Original status | Current status |
|---|---|---|---|---|
| REC-UI-01 | Restore the missing cyan halo around the user message bubble; add a consistent gap so adjacent bubble borders do not touch or visually overlap; preserve the neon/TRON language and readability | `PROJECT_ISSUES_AND_EXPANSION.md` §FUT-CHUI1; `STAGE_STATUS.md` “Polish backlog” items 1–2 | Deferred / polish backlog | **Owner-confirmed requirement**; implementation not established |
| REC-UI-02 | User appearance control with coordinated colors; selected color affects border, halo, and font; independent per-character presentation in multi-character chat | `PROJECT_ISSUES_AND_EXPANSION.md` §FUT-CHUI2 | Deferred; approximately 20 coordinated choices | **Owner-confirmed**, with the count clarified to approximately ten predefined choices for each configurable property: user bubble, character bubble, halo, text. The exact mapping and theme interaction remain design work. |
| REC-UI-03 | Character image thumbnail beside the speaker name on every message; no image → deterministic fallback; fixed size; resized thumbnail; theme-safe; group distinguishability | `PROJECT_ISSUES_AND_EXPANSION.md` T2.4 and §FUT-CHID1 | T2.4 unverified; FUT-CHID1 deferred | **Owner-confirmed**; chat portrait thumbnails remain required. Exact fallback styling is subordinate to the requirement. |
| REC-UI-04 | Narration feature comparable in purpose to SillyTavern narration | Owner recovery instruction §4.3; exact historical artifact not recovered | **Unestablished historically** | **Owner-confirmed requirement**; semantics, message representation, multi-character interaction, and prompt construction are still recovery/design work |
| REC-UI-05 | Token counter, hidden by default, with a click control to show/hide | Owner recovery instruction §4.4; historical related control at `PROJECT_ISSUES_AND_EXPANSION.md` §FUT-CHTOK1 | FUT-CHTOK1 deferred; independent context %, user-message tokens, response tokens | **Owner-confirmed requirement**; whether one control serves context meter + token display is an implementation-design question |
| REC-UI-06 | Option to initiate the complete image-generation orchestration from a chat message | Owner recovery instruction §4.5 | No exact prior requirement recovered | **Owner-confirmed requirement**; must not bypass orchestration, safeguards, or context assembly |

### 4.3 Historical chat behavior and per-chat requirements

| ID | Historical requirement | Source | Original status | Current status |
|---|---|---|---|---|
| REC-CH-01 | Up to five chat message cues in the composer, persisted with the chat | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-C1 | Deferred real requirement | **Deferred real requirement**; must be included in the per-chat schema union; TTS contract untouched |
| REC-CH-02 | Chat-scoped Author's Note with click-open flyout, floating or hidden presentation, optional side display, chat-only scope, and a blank note when starting a new chat | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-CH1; owner recovery confirmation 2026-11-11 | Deferred real requirement | **Deferred real requirement**; recoverable details now include flyout/hidden/floating/side presentation, chat-only scope, and new-chat clearing. An existing chat retains its own note across reload; insertion order and visible enable state remain required. |
| REC-CH-03 | Single-character chats retain configured context capacity before rollover/memory strategies engage | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-CH2 | Deferred real requirement | **Deferred real requirement**; token accounting, response reserve, observable diagnostics; whether the setting is persisted per chat remains unresolved |
| REC-CH-04 | Multi-character context isolation and continuity, flushing between turns or a demonstrably stronger alternative | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-CH3; T3.5 | Deferred/unverified | **Now explicitly required** for multi-character chats: flush/rebuild after each character response. Single-character chats must not flush after each response. |
| REC-CH-05 | Verify whether context flush exists before optimizing; test `cache_prompt` / `--cache-reuse`; full model unloads are not acceptable as an optimization | `PROJECT_ISSUES_AND_EXPANSION.md` T3.1 | Unverified/research first | **Open technical investigation**; not satisfied by the new model-switch requirement |
| REC-CH-06 | Multi-character chat: 2+ characters fully separate, appropriate system prompts, context flushed between speakers, root cause flat context | `PROJECT_ISSUES_AND_EXPANSION.md` T3.5 | Unverified | **Required** (see REC-CH-04); implementation not established |
| REC-CH-07 | Character-scoped and shared lore, scoped memory, trigger specificity, budget, precedence, conflict handling | `PROJECT_ISSUES_AND_EXPANSION.md` T3.3; REG-1 | Unverified/protected relationship model | **Required product behavior**; substantial lorebook foundation exists, full memory/continuity is not implemented |
| REC-CH-08 | Optional AI-as-user / persona takeover with explicit control and isolation | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-CHAIU1 | Deferred major optional capability | **Deferred**; not authorized; must use the required context isolation |
| REC-CH-09 | Application-side structured memory and continuity: inventory, item ownership, acquisition/transfer/consumption/use, clothing/outfit continuity, equipment state, broader narrative facts | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-MEM1 | Deferred major capability | **Required major capability** under ODR-04; implementation unauthorized |
| REC-CH-10 | Automated lore/memory intelligence overhaul: durable fact extraction, review/correction/override, conflict and duplicate checks, no silent overwrite, unobtrusive review | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-LORE1 | Deferred | **Required as part of the memory direction**, with staged scope; not a license to replace manual lore |
| REC-CH-11 | Conversation context rollover: manage a conversation exceeding available context through context shifting, summarization, retention, rebuilding, and memory reinjection | `PROJECT_ISSUES_AND_EXPANSION.md` T3.1/FUT-CH2; current Context Rollover UI | Unverified/research first | **Clarified as a separate required behavior** from multi-character flushing; current runtime semantics and tokenizer accuracy are unverified |

### 4.4 Historical model and runtime requirements

| ID | Historical requirement | Source | Original status | Current status |
|---|---|---|---|---|
| REC-M-01 | Evidence-based model/template detection; never turn marketing names into facts; verified/source-derived/inferred/unknown/conflicting states | `PROJECT_ISSUES_AND_EXPANSION.md` F2-B4 | Approved | **Accepted**; does not by itself provide per-turn generation provenance |
| REC-M-02 | Expanded model-family support and automatic prompt/chat formatting, Jinja/ChatML, reasoning tags, safe fallback, user override, no filename-only certainty | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-PROMPT1; T1.4 | Deferred documented-only | **Required direction** under the owner recovery instruction; implementation coverage remains partial |
| REC-M-03 | llama.cpp feature exposure: Flash Attention, GPU offload, basic multi-GPU, hardware-specific optimizations | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-C3 | Deferred | **Required investigation/exposure**, capability-gated |
| REC-M-04 | Comprehensive sampling controls including XTC, min-P, quadratic/other relevant options | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-C5; owner instruction §6.5 | Deferred/research | **Required inventory and capability-gated exposure** against the pinned runtime |
| REC-M-05 | Learn from LM Studio's load speed and portable runtime/resource optimizations | `PROJECT_ISSUES_AND_EXPANSION.md` T3.2; owner instruction §6.4 | Unverified/research | **Required research**, no parity claim |
| REC-M-06 | Useful Kobold-style offloading and context features, without assuming direct transferability | Owner instruction §6.2–6.3; `PROJECT_ISSUES_AND_EXPANSION.md` research queue | Unverified/research | **Required research and gap inventory** |
| REC-M-07 | Unified application-wide hardware-aware model resource orchestration with verified transitions and release | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-HW1 | Deferred/documented-only | **Deferred major foundation**, but now a prerequisite for safe model hot-swap and memory-context accounting |
| REC-M-08 | Hardware compatibility matrix, low-end graceful degradation, measurements on named hardware | `PROJECT_ISSUES_AND_EXPANSION.md` HW-C1–C3 | Deferred/not started | **Deferred validation targets**; not satisfied by desktop inference |
| REC-M-09 | Mobile/remote: desktop fully usable without phone/network/companion; mobile never required; phone extends capability; Android test platform when authorized | `PROJECT_ISSUES_AND_EXPANSION.md` MOB-1–MOB-4 | Binding conditional constraints; no implementation authorized | **Binding constraints**; remote phase remains unauthorized and unstarted |
| REC-M-10 | Content reuse/hot-swap and scope relationships for worlds, scenarios, characters, lorebooks | `PROJECT_ISSUES_AND_EXPANSION.md` REG-1 | Protected regression constraint | **Protected**; do not use “hot-swap” unqualified |
| REC-M-11 | Comprehensive local-inference technique inventory across CPU, RAM, VRAM, storage, cache, context, attention, batching, multi-GPU, and model-specific techniques | Owner architecture expansion 2026-11-11 | **Newly explicit first-class mission** | **Required research**; current llama.cpp inventory is a first-pass seed only |
| REC-M-12 | Map techniques to dense, MoE, sliding-window, recurrent/state-space, hybrid, multimodal, and other architecture families without treating them as interchangeable | Owner architecture expansion 2026-11-11 | **Newly explicit** | **Required architecture matrix**; no universal dense-transformer assumption |
| REC-M-13 | Map techniques to backends with actual support, version, hardware requirements, and limitations | Owner architecture expansion 2026-11-11; F2 model intelligence | **Newly explicit** | **Required backend/version/platform matrix**; candidate runtime is not load compatibility |
| REC-M-14 | Record optimization interactions and conflicts, including combinations requiring tests | Owner architecture expansion 2026-11-11 | **Newly explicit** | **Required conflict graph**; “both accepted” is not compatible |
| REC-M-15 | Automatically detect capabilities and hardware, estimate resources, select safe combinations, and explain unavailable optimizations | Owner architecture expansion 2026-11-11 | **Newly explicit** | **Required design**; conservative, explainable, no false precision |
| REC-M-16 | Preserve advanced manual override with clear warnings when an override creates risk | Owner architecture expansion 2026-11-11 | **Newly explicit** | **Required design**; overrides scoped, validated, reversible, and recorded |
| REC-M-17 | Benchmark performance, memory, stability, and output-quality effects rather than assuming an accepted parameter is successful | Owner architecture expansion 2026-11-11 | **Newly explicit** | **Required benchmark protocol**; named hardware and baseline required |

### 4.5 Historical creative/product direction

| ID | Historical requirement | Source | Original status | Current status |
|---|---|---|---|---|
| REC-C-01 | Separate modality-specific creative workspaces; shared catalog must not merge chat and image interfaces; capability filtering | `PROJECT_ISSUES_AND_EXPANSION.md` F2-B5 | Approved | **Accepted**; implementation/UI filtering follow-up remains |
| REC-C-02 | Custom workspace layout under Documents/Pictures | `PROJECT_ISSUES_AND_EXPANSION.md` E.6 | Parked/unverified | **Still parked**; the general studio approval does not promote E.6 |
| REC-C-03 | Image staging-board structure mirrored into the image directory | `PROJECT_ISSUES_AND_EXPANSION.md` E.8 | Parked/unverified | **Still parked**; not part of the approved general studio by implication |
| REC-C-04 | Music/audio expansion deferred until subsystems ready | `PROJECT_ISSUES_AND_EXPANSION.md` F2-B7; T3.8 | Deferred/frozen | **Deferred**, but included in the approved general-studio direction; no implementation |
| REC-C-05 | Video generation exists as an incomplete subsystem; no video sidecar/catalogue/delete contract | Master plan S9; F3-G scope | Source-confirmed/incomplete | **Deferred/incomplete**; not made complete by the studio direction |
| REC-C-06 | Reactive 3D desktop avatar and multi-character avatars | `docs/tier4-avatar.md`; Tiers 4.1/4.2 | Parked/unscheduled | **Excluded** by ODR-01; retained only as historical evidence |
| REC-C-07 | 3D mesh generation and texture generation | Owner instruction §1.3 | **Newly approved direction** | **Approved direction**; separate modality requirements; excludes desktop avatar integration |
| REC-C-08 | General local creative studio beyond chat/images/audio | Owner instruction §1.3 | Previously unestablished | **Approved direction**; no immediate implementation of every modality |

### 4.6 Historical media and persistence requirements

| ID | Historical requirement | Source | Original status | Current status |
|---|---|---|---|---|
| REC-P-01 | Durable image/attachment references; `hasImage` boolean insufficient; unavailable references must remain unavailable | F3-C Phase-0 gap; prior PROP-F3C-04 | Proposed/unadopted in F3-C | **Approved product requirement** under ODR-07; still not an F3-C amendment; schema work remains |
| REC-P-02 | Saved audio remains loadable/playable from its chat context | Owner instruction §3 | **Newly explicit** | **Approved**; implementation not established |
| REC-P-03 | Displayed media is distinct from media supplied as model input | Owner instruction §3 | **Newly explicit** | **Required**; no re-sending on every turn merely because it remains visible |
| REC-P-04 | Delete/move/export/import/restore behavior for chat media must be defined and must not leave broken references | Owner instruction §3 | **Newly explicit** | **Required**; contract not yet defined |
| REC-P-05 | Image metadata records prompt/model/options/seed; chat history does not | Current source; F3-C gap | Source-confirmed | **Existing image behavior**; not chat provenance and not a chat attachment contract |
| REC-P-06 | Imported-field provenance | `PROJECT_ISSUES_AND_EXPANSION.md` FUT-ST2 | Deferred | **Deferred**; separate from per-turn generation provenance |

### 4.7 Original separate UI planning document

The prior investigation did not conclusively recover the separate UI document, and this recovery did not find it. The exact artifact is **not established**.

Searches performed:

- exact filename and likely UI-planning filename variants;
- repository-wide content search for “possible expansion of howl chat for review”;
- git history filename enumeration;
- checks for deleted or renamed UI planning artifacts;
- direct review of `STAGE_STATUS.md` and the ledger's UI sections.

No file was found whose content can be identified as that separate document. The requirements in §4.2 are therefore sourced to the owner instruction and to the historical entries that do survive. Any claim that the original UI document has been recovered would be false.

---

## 5. Exact recovery of the separate chat UI planning document

### 5.1 Result

**NOT RECOVERED.**

This is a factual evidence gap, not a conclusion that the requirements are new. The ledger's UI entries and `STAGE_STATUS.md` demonstrate that several requirements predate the recovery instruction. The exact separate planning document could not be identified by filename, content, or repository history. The two F2 engineering reports are present but are not that document; the original predecessor-document set and count are not established.

### 5.2 What was recovered instead

**Bubble appearance and polish.** `STAGE_STATUS.md` records two deferred polish items from the group-chat review: the user bubble's border appearing to graze neighbouring character bubbles, and the user bubble lacking the matching user-colour halo that character bubbles have. It explicitly says to check at real window sizes before changing spacing because the preview viewport may be responsible. `PROJECT_ISSUES_AND_EXPANSION.md` FUT-CHUI1 independently requires the cyan user-bubble halo, a consistent gap, and preservation of the neon/TRON design and readability.

**Per-character appearance.** FUT-CHUI2 records a user-section appearance control, approximately 20 coordinated colour choices affecting border, halo, and font, and independent per-character settings in multi-character chat. The owner's 2026-11-11 instruction changes the count to approximately ten predefined choices for each configurable property (user bubble, character bubble, halo, text). The later owner instruction controls the count; the historical requirement controls the four independently configurable properties and per-character scope.

**Portrait identity.** T2.4 requires a character-card thumbnail beside the speaker name, a fallback initial on a dark-blue field, fixed size, theme-safe treatment, resized thumbnails, and clear group distinguishability. FUT-CHID1 adds a circular fallback using the character's selected border colour and the first initial in white. These are compatible; the later details refine the fallback rather than replace the thumbnail requirement.

**Token visibility.** FUT-CHTOK1 requires independent display options for context usage as a percentage, token count of the newly entered user message, and token count of a specific generated response, with no false precision. The 2026-11-11 owner instruction requires a token counter, hidden by default, with a click control. The historical and current requirements can be served by one visibility control only if the underlying values remain independently meaningful.

**Narration.** No historical requirement for a SillyTavern-like AI narration feature was recovered. The source contains a different concept: system/narrator messages used for character step-out, death, return, revival, and similar events. That is an event-notice mechanism, not a narrator mode that writes prose, controls perspective, or participates in multi-character turn construction. The two must not be conflated.

**Author's Note.** The follow-up owner recovery confirms the earlier deferred requirements and adds presentation/scope detail: a click-open flyout; floating or hidden presentation; optional side display; chat-only scope; and a note cleared when starting a new chat. These details are consistent with FUT-CH1's chat-scoped persistence, insertion order, and visible enable state if “per-chat persistence” means an existing chat retains its own note across reload while a new chat starts blank. The corrected plan must state that reconciliation explicitly rather than treating “persisted” and “cleared on new chat” as contradictory. The exact source artifact for these details is still not recovered; the owner follow-up and ledger are the evidence.

**Message-triggered image generation.** No exact prior requirement was recovered. The ledger has adjacent but different concepts: T3.10 asks whether a local LLM can write an image prompt; E.1 discusses passing reference images to a vision-language model; image context menus support recall/import-seed/use-as-init. None of those is the owner-confirmed “generate an image from this chat message” workflow.

### 5.3 Consequence for the corrected plan

The corrected master plan must cite the owner instruction for the newly explicit UI requirements and the ledger for the historically established ones. It must not claim that a separate recovered UI document supplied wording it did not supply.

---

## 6. Missing requirements and misclassifications in the current master plan

| ID | Current master-plan treatment | Correct treatment | Evidence / consequence |
|---|---|---|---|
| GAP-01 | Model hot-swap is “unaccepted product behavior” / adjacent to F3-C | **Required product behavior**, separate from F3-C and REG-1 | ODR-02; master plan S15, §4.7 PROP-F3C-05, §7.3, §8.2, §9 |
| GAP-02 | General creative studio is “unestablished” | **Approved direction** beyond F2-B5/E.6/E.8, with each modality separately scoped | ODR-03; master plan S24, §4.8, §9 |
| GAP-03 | Desktop avatar is a parked/optional future | **Explicitly excluded**; historical only | ODR-01; master plan S22, T4.1/T4.2, §4.3, §4.8 |
| GAP-04 | Durable attachments/seed/provenance are proposal-only owner questions | **Attachments approved; provenance required concept; seed preservation directed**; exact schema/recommendation still open | ODR-05–07; master plan §4.7, FND-07/FND-08, §8.2 |
| GAP-05 | Memory is a deferred feature row (FUT-MEM1) and not a major architecture driver | **Required major subsystem**, implementation unauthorized | ODR-04; master plan S6, FND-09, §4.5, §8.3 |
| GAP-06 | Multi-character isolation is “unverified” with no explicit single-character rule | **Flush after each character in multi-character; do not flush per turn in single-character** | ODR-08; T3.5, FUT-CH3 |
| GAP-07 | Chat UI requirements are scattered or absent | Add bubble spacing/border, per-character four-property colors, portraits, narration, hidden token counter, message-triggered image generation | ODR-09; §4.2 |
| GAP-08 | Narration is not represented | Add as owner-confirmed requirement with semantics explicitly unrecovered | §5.2 |
| GAP-09 | Message-triggered image generation is not represented | Add as owner-confirmed orchestration path, not a shortcut | §4.2 |
| GAP-10 | Model-family automation is a deferred generic foundation | Required direction with a family/confidence matrix and hot-swap transition behavior | REC-M-02; owner §6.1 |
| GAP-11 | Sampler/offloading/context-shift/LM Studio coverage is not inventoried against the pinned runtime | Add capability-gated inventory and gaps | §14 |
| GAP-12 | Per-chat continuation-model preference is conflated with historical provenance or hot-swap | Keep unresolved and separate; do not infer restoration behavior | ODR-02/06; prior PROP-F3C-02 |
| GAP-13 | “Creative studio” can absorb E.6/E.8 or F2-B5 by implication | Preserve F2-B5 approved, E.6/E.8 parked, general studio separately approved | REC-C-01–03 |
| GAP-14 | W4/W8 carry-forward | W4 historical/superseded by W5–W7; W8 VERIFIED | Adjudication F-07/F-14 |
| GAP-15 | The separate UI document and expansion text are implied to be recovered | Record exact recovery as not established; F2 reports and ledger are leads, not proof of complete survival | §1.5, §5 |
| GAP-16 | Inference optimization is a short backend/Flash-Attention note | Make technique/architecture/backend research a first-class architecture mission | ODR-10; §15 |
| GAP-17 | Conversation context rollover and multi-character context flushing can be read as one feature | Define two separate mechanisms, triggers, scopes, and invariants | ODR-11; §12, §15.9 |
| GAP-18 | Author's Note is a one-line deferred feature | Add recovered flyout/floating/hidden/side presentation, chat-only scope, and new-chat clearing | REC-CH-02; §5.2 |
| GAP-19 | The first-pass llama.cpp inventory can be mistaken for a completed capability survey | Label it a seed; require other engines, architectures, versions, and benchmarks before freezing | REC-M-11–17; §15 |
| GAP-20 | Presence of two F2 reports and the ledger is treated as proof that all original requirements survived | Treat them as evidence leads; keep the missing UI/expansion artifacts explicit | §1.5, §5 |
| GAP-21 | Network policy is a privacy setting or a generic “local-first” claim, not a first-class architecture/release requirement | Add the destination registry, six-question consent, purpose-scoped permissions, audit, backend offline enforcement, failure visibility, and release acceptance criterion | ODR-12; §23; N-01–N-14; V-NET-01–V-NET-28 |

---

## 7. Current implementation status and supporting evidence

### 7.1 Status summary

| Capability | Current status | Source evidence | Limit |
|---|---|---|---|
| Chat history model/seed/template provenance | **Absent** | `internal/backend/chat/types.go` `ChatHistory` has ID, sender, text, role, timestamp, `hasImage`, `characterId` only | Static source; no runtime schema migration |
| Per-chat model preference | **Absent** | `ChatMeta` has no model field; model selection uses `localStorage` and global settings | Does not imply desired restoration semantics |
| Active-conversation model switching | **Partial/manual only** | `app-model.go` `LoadModel` stops an existing server, then spawns the target; UI dropdown has an explicit load action | No stream coordination, readiness identity proof, rollback, or weight-release verification |
| Load-on-typing | **Absent** | Model loading is initiated from the model picker/load path, not the composer | Owner requires either Load or beginning to type |
| Same-model skip guard | **Source-confirmed** | `LoadModel` checks the loaded model ID before stopping | Relies on `/v1/models`; not a full identity/version check |
| Replacement readiness | **Partial** | Monitor calls `HTTPClient.Health()` and emits `model:loaded` | Health is not model-identity, tokenizer, template, or generation-readiness proof |
| Old-model weight release | **Not verified** | `stopLlamaServer` kills the process tree and waits for monitor exit | Process exit does not prove driver/VRAM release; no measured assertion |
| Stream coordination during switch | **Absent** | `beginStream`/`cancelActiveStream` exist, but `LoadModel` does not call them or reject a switch while streaming | A load can stop the server while a request is active |
| Rollback to previous model | **Absent** | Old server is stopped before target spawn; no stored reload plan | Failure leaves no active model unless the user retries |
| Prompt template transition | **Partial** | `spawnSettingsFor` applies the target model's native template per load and preserves a user override | No fingerprint or per-turn template provenance; no hot-swap test |
| Model-family detection | **Partial** | GGUF `general.architecture`, metadata, family table, model-intelligence profiles | Coverage is incomplete; filenames still act as hints; no proof for obscure families |
| Seed configuration | **Global, ambiguous** | Global sampler `seed` default `-1`; UI HTML shows a numeric seed field; no per-turn actual seed | No actual seed returned/parsed; no origin/status; no reproduction evidence |
| Per-turn provenance | **Absent** | No fields in persisted `ChatHistory`; `types.Message.Metadata` exists in a separate path | `types.Message` is not the persisted chat shape |
| Durable chat image reference | **Absent/transient** | `ChatHistory.HasImage` only; image path uses current-turn base64; tutorials state previous images are not seen unless attached | No attachment ID, ownership, unavailable-reference semantics, or restore contract |
| Saved audio in chat | **Not established** | Audio/TTS pipeline exists; no chat-media association contract found | Do not infer from TTS playback |
| Chat portrait thumbnail | **Partial** | `chat-render.js` renders resolved character image in `.character-thumbnail`; CSS is 24px | Streaming templates at `chat-stream.js` create sender text without the thumbnail; no complete avatar/fallback contract |
| Bubble gap/border repair | **Not established** | `.chat-message` has internal 12px gap; `.chat-bubble` uses 1px border and 10px glow; no fix for border grazing | Cannot confirm visual defect or repair without runtime/visual evidence |
| Per-character bubble colours | **Absent** | Global `--user-*` / `--ai-*` role tokens; no per-character style state found | Owner requires four properties × ~10 choices per character |
| Token counter | **Absent** | Context-window setting exists; tutorials mention “Context Usage”; no token counter control found | No tokenizer-backed display; no user-message/response counts |
| Narration mode | **Absent** | System “Narrator” messages and `queueNarratorEvent` provide side-effect notices | Not a SillyTavern-style narration feature |
| Message-triggered image generation | **Absent** | No chat-bubble control calling image orchestration found | Adjacent image context-menu features are not this requirement |
| Structured memory/inventory | **Partial/absent** | `memory.Service` stores prose-like records with scope/owner/tags/confidence; candidate queue; continuity export/import | No typed inventory, possessions, clothing state, relationships, validity intervals, or retrieval planner |
| Memory prompt injection | **Split** | `chat.Service` loads memory snapshots and builds `[RELEVANT MEMORIES]`; live `App.buildPromptMessages` does not demonstrably consume the memory service | One live prompt path is unresolved |
| Multi-character flush | **Not established** | Each character call is a new HTTP request carrying a frontend-composed prompt; no character-scoped context reset endpoint found | Server KV/cache reuse and actual isolation are not proven |
| Single-character no-flush rule | **Compatible by default, not formalized** | No per-turn flush path exists | Absence is not a designed guarantee; future changes must preserve it |
| Model-family auto template | **Partial** | Native GGUF template plus `spawnSettingsFor`; Qwen/Gemma family defaults | Manual override and fallback exist; incomplete matrix |
| Sampler UI | **Partial** | Temperature, top-p, top-k, min-p, repeat, presence/frequency, typical-p, Mirostat, dynamic temp, DRY multiplier/allowed/base, seed | Missing XTC, top-n-sigma, adaptive-p, TFS exposure, DRY last-n/sequence breaker, sampler sequence, and conflict/capability metadata |
| Runtime offload / context shift | **Partial/implicit** | GPU layers, Flash Attention, tensor split are passed; `--context-shift` is not explicitly configured, so runtime default applies | No UI/status, no SWA/context-shift interaction handling, no measured cache reuse |
| General studio | **Absent as a unified subsystem** | Image and video engines exist; no music/3D/texture generation subsystem | Approved direction does not equal implementation |

### 7.2 What the source confirms — and what it does not

The following are **source-confirmed absences**, not runtime failures:

- no model/seed/template fields in the persisted chat history;
- no chat-level model field;
- no durable media identifier in `ChatHistory`;
- no character-scoped context reset API;
- no per-turn provenance object;
- no per-character appearance object;
- no token counter binding;
- no narration-mode prompt path;
- no chat-message image-generation action.

The following are **not established**:

- whether the current delimiter of character contexts is sufficient at the model/server level;
- whether the old model's VRAM is fully released after process exit on each backend;
- whether `Health()` is sufficient for the pinned server's readiness;
- whether a seed saved from the OpenAI-compatible endpoint reproduces a response;
- whether existing chat media references survive export/import/restore;
- whether current chat UI border overlap is a real layout defect or a preview artefact;
- whether the live memory/lore path actually affects every character reply.

### 7.3 Historical implementation records are not current verification

`STAGE_STATUS.md` records T1.1–T1.4 closures, group-chat runs, and theme verification. Those records remain historically valid for what they exercised. They do not establish current product-wide verification for chat ownership, memory, model switching, media, or character isolation. The F3-C defects remain open and are not closed by any stage closure.

---

## 8. Seed and per-turn provenance recommendation

### 8.1 Seed: what is actually available

The pinned runtime reports:

```
version: 0.8.0-dev (build 12015, commit a34de8618)
```

llama.cpp exposes `--seed` and a random sentinel of `-1`. The upstream server API has historically distinguished the **requested** seed from the **actual** seed: `seed` echoes the request while `seed_cur` reports the seed used when the request was random. The relevant upstream change is llama.cpp #9398 (“Add seed_cur to completion response”).

However:

- Howl's current chat client uses the OpenAI-compatible `/v1/chat/completions` path.
- Its response struct parses choices and usage only. It does **not** parse `seed` or `seed_cur`.
- The current persisted chat record has no seed field.
- The current global sampler default is `-1`, while the UI seed input is a numeric field; there is no explicit origin/status model.

**Therefore:** actual-seed recording is technically plausible but not currently implemented, and whether the pinned build exposes the actual seed on the exact endpoint Howl uses must be tested before it becomes a formal guarantee.

### 8.2 Recommendation

**Recommended, not owner-approved as a schema:**

1. Treat each assistant turn as the provenance unit.
2. Request generation with the user's explicit seed when one is set; otherwise use the runtime's random sentinel and let the runtime choose.
3. Capture the **actual** seed when the backend exposes it. Never infer it from the request.
4. Store a seed status alongside it:
   - `explicit` — the user supplied this value;
   - `random` — the runtime generated this value;
   - `unavailable` — the runtime did not expose the actual value;
   - `legacy_unknown` — an older turn has no record.
5. Preserve explicit zero as a valid seed; do not use zero as “unset.”
6. Also retain a last-used seed pointer for composer convenience, but do not treat it as per-turn provenance.
7. Do not require the user to manage a random-versus-explicit toggle unless testing shows it is necessary.
8. Do not promise exact reproduction. A seed is one input among many.

### 8.3 Why a last-used seed alone is insufficient

A last-used seed answers “what should the composer show next?” It does not answer “what produced this old response?” A conversation may switch models, sampler settings, templates, and context. Per-turn provenance requires the seed that belonged to that turn. A chat-level last seed would be overwritten and would make historical debugging ambiguous.

### 8.4 Per-turn provenance record

The minimum meaningful record should be immutable per assistant turn and should survive model uninstall or replacement. Recommended fields:

| Group | Suggested fields | Requirement status |
|---|---|---|
| Model identity | Display name, file basename, stable descriptor (size/modtime/hash where available), architecture, family, quantization/version | Concept required; exact fields proposed |
| Engine | Backend/engine name and version; server build/commit where available | Concept required; proposed |
| Template | Template identity, source (model-native/user/fallback), fingerprint/hash, length, format family | Concept required; proposed |
| Generation | Effective sampler values, stop strings, max tokens, context size, reasoning mode/budget, seed + seed status | Concept required; proposed |
| Runtime | Relevant hardware/offload/context-shift/cache settings that materially affect output | Proposed; do not over-record unrelated host data |
| Outcome | Completion status, finish reason, prompt/completion token counts where exposed | Proposed; useful for audit and context accounting |

Do **not** store the full prompt by default. Provenance is about generation identity, not a transcript duplicate. Prompt retention is a separate privacy/storage decision.

### 8.5 Reproduction honesty

A saved seed does not guarantee the same text. Reproduction can vary with:

- model file/version/quantization;
- runtime build and kernels;
- hardware, drivers, and GPU/CPU split;
- prompt template and exact message serialization;
- context contents and ordering;
- sampler sequence and effective values;
- nondeterministic parallel execution or backend behavior.

The correct user-facing claim is “the actual seed and configuration were recorded,” not “this response can be reproduced.” A reproducibility test may be offered as evidence, but failure of exact reproduction must not invalidate the provenance record.

### 8.6 Relationship to existing provenance requirements

| Provenance family | What it proves | What it does not prove |
|---|---|---|
| F2-B4 model/template detection | How the app classified a model and selected a template | That a particular response was generated by that model/config |
| FUT-C5 sampler research/provenance | Which sampler capabilities exist and what was studied | The effective sampler state of a turn |
| FUT-ST2 import provenance | Where imported card fields came from | Chat generation provenance |
| Image metadata | Prompt/model/options/seed for an image | Chat provenance |
| Per-turn chat provenance | The model/config that produced this assistant turn | Exact reproducibility or user intent |
| Seed status | Origin/availability of the seed | That the same output will recur |

Adopting one family must not be described as adopting another.

---

## 9. Active-conversation model hot-swap lifecycle requirements

### 9.1 Product behavior required

The owner requires this sequence:

1. The user selects a different model.
2. Loading can start from the Load control or by beginning to type a message.
3. The application safely unloads the previous model completely.
4. Previous model weights are no longer resident.
5. The replacement model loads.
6. Howl verifies that the replacement is actually ready.
7. The model-family prompt template is selected automatically where needed.
8. The pending message is sent using the active cue and existing conversation.
9. The user receives clear success or failure feedback.

The conversation, relevant context, multi-character isolation, and per-character state must survive.

### 9.2 Required lifecycle states

A model switch must have one explicit state machine, not a progress percentage plus a registry entry:

```
idle
  → requested
  → quiescing (no new generation; cancel or await active stream)
  → unloading
  → unloaded_verified
  → loading
  → probing_readiness
  → ready

failure branches:
  unload_failed
  load_failed
  readiness_failed
  identity_mismatch
  cancelled
```

Each state must publish:

- target model identity;
- prior model identity (for recovery messaging);
- whether generation is admitted;
- whether a pending composer message exists;
- the user-visible failure reason;
- the safe next action (retry, choose model, unload, inspect diagnostics).

### 9.3 Required invariants

| ID | Invariant |
|---|---|
| MHS-01 | A model is not reported “active” until readiness and identity checks pass. |
| MHS-02 | No message is sent while the target model is loading or while identity is unverified. |
| MHS-03 | A pending message is staged before unload and is not silently discarded by a failure. |
| MHS-04 | A switch never assumes both models can reside simultaneously. |
| MHS-05 | The previous model is stopped and its process tree is gone before the target is declared active. |
| MHS-06 | Release of memory/VRAM is verified by an external measurement or explicit backend signal, not process exit alone. |
| MHS-07 | The target template comes from the target model, not the previous model, unless a recorded user override applies. |
| MHS-08 | The target model identity is compared with the user's selection before send. |
| MHS-09 | An active stream is cancelled, completed, or explicitly waited for before unload; a switch cannot silently interrupt it. |
| MHS-10 | Loading on typing is a request to switch, not permission to send to the wrong model. |
| MHS-11 | Failure leaves a clear state and a recoverable path; it does not leave a false “loaded” label. |
| MHS-12 | Conversation history, character identity, lore, memory, and pending cues are not deleted by a switch. |
| MHS-13 | Multi-character context isolation rules continue to apply after the switch. |
| MHS-14 | Per-turn provenance records the model that actually produced the response. |
| MHS-15 | The switch is idempotent: repeated requests cannot spawn duplicate servers. |
| MHS-16 | Shutdown during a switch is safe and leaves no owned child process. |
| MHS-17 | A failed unload does not proceed to load and pretend the old model is gone. |
| MHS-18 | A failed target load does not auto-send the pending message or claim success. |
| MHS-19 | Old model weights are not retained merely to make rollback convenient; any old-model reload is a new, explicit load. |
| MHS-20 | The UI exposes the target, state, failure, and recovery action without requiring log inspection. |

### 9.4 Readiness verification

A single `/health` response is insufficient. The switch should verify, as capabilities permit:

1. process alive and owned;
2. HTTP health;
3. server build/version;
4. `/v1/models` or equivalent model identity matches the requested file;
5. model metadata/architecture matches the admission record;
6. context capacity is reported or bounded;
7. template was applied from the target model or a recorded override;
8. a minimal generation or tokenization probe succeeds where safe;
9. the server is not still holding a previous instance on the same port;
10. resource release/usage is consistent with the new model.

Items 4–8 are the difference between “server answered” and “correct model can generate.”

### 9.5 Failure and recovery semantics

- **Unload failure:** do not start the target. Report that the old model could not be confirmed gone. Offer retry/force-stop diagnostics; do not claim memory was freed.
- **Spawn/load failure:** no active model; keep the pending message and conversation. Offer retry or another model. Do not silently reload the prior model unless the user asks; that is another switch and may fail for the same reason.
- **Readiness/identity failure:** stop the target process, do not send, preserve the pending message, and show the mismatch.
- **User cancellation:** stop the target load safely; preserve the conversation and composer. State whether any model is active.
- **Shutdown during switch:** cancel and join the load, stop owned children, and report incomplete cleanup if it cannot be verified.

### 9.6 Current-code delta

| Requirement | Current position |
|---|---|
| Stop old model before loading new | Present |
| Admission before stopping old model | Present for GGUF architecture/task/projector checks |
| Per-target template selection | Present via `spawnSettingsFor` |
| Server health check | Present |
| Model identity verification after load | Not established |
| Stream quiescence | Absent |
| Weight-release verification | Absent |
| Pending-message staging across switch | Absent |
| Automatic load on typing | Absent |
| Rollback/explicit recovery contract | Absent |
| Per-turn provenance of actual model | Absent |
| Switch state machine | Absent; progress/stage strings only |

This is a product and architecture gap, not evidence that manual model loading is wholly broken. Manual model switching may work in ordinary conditions; it has not been shown safe for the owner-required lifecycle.

---

## 10. Durable attachment and media requirements

### 10.1 Product contract

Chat media must be part of conversation continuity without becoming an unconditional part of every model prompt.

The contract must distinguish:

- **display reference** — the image/audio is shown or playable in the chat;
- **input reference** — the media was supplied to a model on a specific turn;
- **stored object** — the durable bytes and metadata owned by the application.

One media item may have both a display and an input reference, but neither implies the other. A displayed image is not re-sent merely because it remains visible.

### 10.2 Required media record

A durable attachment record should carry, where applicable:

| Field | Purpose |
|---|---|
| Stable media ID | Survives rename, move, and display changes |
| Kind | Image, audio, video, future asset |
| MIME type and extension | Correct rendering/playback |
| Relative locator plus owning root | Portable within the local library; no absolute-path lock-in |
| Content hash and byte size | Identity, deduplication, missing/corrupt detection |
| Source | Generated, imported, copied from character/world/scenario, TTS output |
| Created/updated timestamps | Lifecycle and conflict resolution |
| Message/turn association | Which message displays it and which turn consumed it as input |
| Input role and turn | Vision input, audio input, reference-only, none |
| Availability state | Available, missing, corrupt, external-root disconnected |
| Ownership/cleanup policy | Chat-owned, library-owned, shared, or explicitly retained |

A boolean `hasImage` cannot carry this contract.

### 10.3 Lifecycle requirements

| ID | Requirement |
|---|---|
| MED-01 | Message history retains a stable reference even when the media is not in the current prompt. |
| MED-02 | A missing or unavailable object renders as unavailable; it is never silently substituted with another image/audio file. |
| MED-03 | Deleting a message must not silently delete a media object still referenced elsewhere. |
| MED-04 | Deleting a chat defines whether chat-owned media is retained, moved to a library, or removed; the result must be visible. |
| MED-05 | Moving a media library or changing a registered root preserves identity or reports the relocation; it cannot produce an unrelated replacement. |
| MED-06 | Export includes enough metadata to resolve media references, or explicitly marks them external. |
| MED-07 | Import validates references and marks unresolved media unavailable; it does not invent or rebind to same-named files. |
| MED-08 | Restore/reload preserves display associations and input-history truth independently. |
| MED-09 | Wipe/cleanup follows explicit ownership; it must not delete external registered libraries. |
| MED-10 | Generated image/audio is committed to a stable object before the chat references it. |
| MED-11 | Image input is recorded as input for the turn; continuing to display it does not re-send it. |
| MED-12 | Audio remains loadable/playable after chat reload and after the generating process has exited. |
| MED-13 | The local-first privacy boundary is preserved; no cloud upload is implied by attachment. |
| MED-14 | Media identity and provenance are separate: a reference says which object; generation metadata says how it was made. |

### 10.4 Interaction with existing storage

The existing surfaces are not one system:

- gallery images and sidecars (`internal/image/gallery`, image metadata);
- image output directories and generated file names;
- chat folder snapshots and `history.ndjson`;
- IndexedDB/localStorage character and chat state;
- future audio/TTS outputs and video outputs;
- content-library registrations and external roots.

F2-C5 already records images and videos as “binary-plus-sidecar generated media with stable identity and atomic commit.” F3-A6 records atomic image+sidecar commit and containment. Those are useful foundations, but they do not define chat attachment ownership. F3-G is the future video media contract; the absence of a video catalogue means video attachments have no planning basis yet.

### 10.5 Current delta

- Chat persistence stores only `hasImage`; no media ID, path, hash, availability, or input-use record.
- Chat image input is current-turn base64; the tutorial explicitly says previous images are not seen unless still attached.
- No saved-audio-in-chat association contract was found.
- Export/import/restore behavior for chat media is not established.

The owner's approval makes the contract a product requirement. It does not make the current implementation durable.

---

## 11. Application-owned memory-system architecture

### 11.1 What the memory system is — and is not

The memory system is **not** a renamed lorebook. It is the application-owned layer that records what has happened and what is currently true, then selects a small relevant subset for a turn.

| Layer | Owner | Nature |
|---|---|---|
| Canonical transcript | Chat persistence | Immutable/ordered record of messages and turns |
| Authored world/scenario/character canon and lore | Existing content ownership and lorebook system | Author-established rules, facts, relationships, and triggers |
| Durable memory/event state | New memory subsystem | Learned facts, state changes, inventory, clothing, relationships, episodes, promises |
| Derived summaries/embeddings | Memory subsystem, rebuildable | Compression and retrieval aids; never the sole authority |
| Turn memory packet | Prompt assembly | The selected, budgeted memory text for this generation |

The model may propose extraction and summaries. The application owns validation, identity, scope, persistence, conflict resolution, retrieval, and reinjection.

### 11.2 Memory record model

A durable memory item should support:

- memory ID and schema version;
- kind: fact, state, inventory item, possession, clothing/appearance, relationship, event, promise, location, preference, episode, summary;
- subject and owner IDs (character, chat, world, scenario, shared scene);
- predicate/attribute and value;
- validity interval (`fromTurn`, optional `toTurn`) and supersession link;
- source turn/message IDs;
- confidence and evidence class;
- status: candidate, active, superseded, archived, rejected;
- salience/importance and recency;
- privacy/visibility scope;
- user-edit/lock state;
- optional tags and retrieval terms;
- provenance of extraction (manual/model/import) without treating model output as authority.

Inventory requires quantity, owner, location/container, state, and acquisition/transfer/consumption events. Clothing/appearance requires garment/item identity, worn/removed state, visible description, and change history. Relationship state requires participants, relation type, direction where relevant, and the event/turn that changed it.

### 11.3 Lifecycle

```
turn observed
  → candidate extraction
  → schema validation
  → scope/identity resolution
  → conflict and duplicate check against active memory + manual lore
  → pending review or auto-accept according to an explicit policy
  → durable commit (append event + update projection)
  → retrieval eligibility
  → summary/compression when needed
  → archive/decay/delete by explicit policy
```

Rules:

- Never silently overwrite an authored lore fact.
- A memory update is an event plus a new current projection, not destructive replacement of history.
- Deleting a memory must be distinguishable from marking it no longer current.
- A model-extracted candidate that cannot be scoped to a stable owner fails closed for injection.
- Review must be unobtrusive during roleplay; batched review is preferable to interrupting every turn.

### 11.4 Retrieval and prompt reinjection

Recommended retrieval pipeline:

1. Build a query from the current user turn, active scene, active character, participants, and explicit references.
2. Filter by visibility, character ownership, chat/world/scenario scope, status, and validity at the current turn.
3. Retrieve candidates using deterministic lexical/tag/entity matching first; add semantic embedding retrieval only when evidence shows it improves recall without unacceptable cost.
4. Rerank by relevance, recency, importance, confidence, and unresolved-change priority.
5. Deduplicate and resolve conflicts; authored canon wins where applicable.
6. Budget tokens against the active model/context and reserved response space.
7. Compress selected items into a compact memory packet with source/validity cues.
8. Inject the packet as application-owned context, separate from authored lore and from the recent transcript.
9. Record which memory IDs were injected for the turn (provenance/debugging), without making the packet a second transcript.

A compact clip should answer “what must this character/scene know now?” It should not restate the character card or the entire lorebook.

### 11.5 Scope and isolation

| Memory scope | Visibility |
|---|---|
| Character-private | Only that stable character ID |
| Chat-private | The current conversation only |
| Shared scene | All active participants in the current scene, subject to privacy rules |
| World | All chats using that world |
| Scenario | Chats using that scenario |
| Global/system | Only where explicitly applicable; not a default dumping ground |

Character-specific memory must not leak to another character merely because both are in the same group chat. Shared events may be visible to all participants, but private thoughts, relationship state, and character-only facts require explicit scope.

### 11.6 Relationship to existing requirements

- **Chat history:** the transcript remains canonical; memory is a derived, selective continuity layer.
- **Character cards:** static identity/appearance baseline; memory records deltas and learned facts.
- **Lorebooks:** authored world knowledge and triggers; memory is dynamic state. Automated lore extraction may produce memory candidates but must not silently become manual canon.
- **World/scenario:** shared context and scope boundaries; memory cannot bypass REG-1 relationships.
- **Automated lore triggers:** remain the trigger/retrieval mechanism for authored lore; memory retrieval uses its own relevance policy.
- **Multi-character isolation:** memory filtering must use stable character IDs and the same isolation boundary as prompt construction.
- **Inventory/clothing:** first-class typed state, not prose hidden in a generic lore entry.
- **Durable conversation data:** every memory item must point to source turns/messages and survive reload/export/import.
- **Context limits:** tokenizer/context accounting is mandatory; no character-count approximation as the final policy.

### 11.7 App versus model responsibilities

| Responsibility | Application | Model |
|---|---|---|
| Stable identity and scope | Owns | May suggest, never authorizes |
| Schema validation | Owns | Produces candidates |
| Canonical fact/state | Owns | Cannot overwrite |
| Extraction | Validates/commits | May propose |
| Summarization | Controls policy and storage | May draft |
| Retrieval | Owns query/filter/rank/budget | No |
| Prompt reinjection | Owns | Receives packet |
| Conflict resolution | Owns | May flag ambiguity |
| Deletion/retention | Owns | No |

This division is required for smaller local models: the app must not rely on the model to remember what the app can store and retrieve deterministically.

### 11.8 Current delta

Current `memory.Service` provides chat-scoped records, candidate review, confidence, tags, pinning, scopes, stable character ownership, and continuity import/export. This is a useful foundation. It is not the required subsystem because it lacks typed inventory/clothing/relationship state, validity/supersession semantics, a full retrieval/reranking pipeline, token-budgeted prompt packets, and a demonstrated single live prompt path for injection.

---

## 12. Multi-character and single-character context behavior

### 12.1 Required rule

This section defines **multi-character context flushing only**. Conversation context rollover — managing a conversation that exceeds its available context through shifting, summarization, retention, and rebuilding — is a separate feature defined in §15.9. Neither is model hot-swap.

| Conversation shape | Required behavior |
|---|---|
| Multi-character | After each character produces a response, flush/rebuild that character's active inference context before another character generates. |
| Single-character | Do not flush after each response. Preserve normal continuation. |

The rule is about **active inference context**, not durable data.

### 12.2 What “flush” must mean

A compliant flush must:

- invalidate or rebuild the model/server context used for the previous character's prompt;
- prevent one character's active prompt, hidden reasoning, or transient context from being reused as another character's generation prefix;
- preserve the canonical transcript, character state, lore, memory, and pending cues;
- use stable character identity, not display-name matching;
- be observable in tests or diagnostics, not inferred from a new HTTP request.

A new HTTP request is not automatically a flush. The server may reuse a KV cache or a prompt prefix. Conversely, throwing away all history is not an acceptable flush.

### 12.3 Baseline turn sequence for multi-character

1. Resolve the speaker and canonical history.
2. Build that character's scoped prompt from shared transcript + character state + authorized lore/memory.
3. Generate.
4. Commit the response to canonical history with character ID.
5. Invalidate/flush the character-specific inference context.
6. Build the next character's prompt from the canonical state, not from the previous character's transient buffer.
7. Repeat.

### 12.4 Single-character rule

A single-character conversation keeps its normal continuation. A flush after every turn would add cost, may degrade coherence, and is not authorized by this requirement. Context-shift, summarization, and model switching remain separate mechanisms.

### 12.5 Distinctions that must survive

| Mechanism | Purpose | Not the same as |
|---|---|---|
| Character context flush | Prevent cross-character prompt/context bleed | Context shift, summarization, model unload, memory deletion |
| Context shift | Drop/shift old KV context when the window is full | Character identity isolation |
| Summarization/compression | Reduce long history into a smaller representation | A mandatory per-turn flush |
| Model hot-swap | Replace the inference model | Character context flush |
| Memory retrieval | Select durable facts/state for a turn | Full transcript retention |
| REG-1 content hot-swap | Activate/reuse content relationships | Any inference lifecycle event |

### 12.6 Current delta

No character-scoped context-reset API was found. `GenerateCharacterResponse` receives a frontend-composed prompt as one user message and invokes a new HTTP request. `ClearChat` is a global in-memory reset, not a character isolation primitive. The server's cache behavior is not controlled or measured here. Therefore the requirement is **not established**, even if ordinary group replies appear separated.

---

## 13. Model-family and prompt automation coverage

### 13.1 Current evidence

The current source uses GGUF metadata (`general.architecture`, model/tokenizer metadata, native chat template) and a small set of explicit family defaults. The model-family UI currently offers Auto, LLaMA, Mistral, Phi, Qwen, LFM, and Gemma. The model-intelligence layer records evidence classes. The prior stage record verifies Qwen and Gemma template separation on real models.

This is meaningful but not the required comprehensive matrix. It does not prove support for obscure families, Jinja edge cases, reasoning-tag variants, or correct automatic behavior after a hot-swap.

### 13.2 Required detection hierarchy

1. **Runtime/model metadata** — GGUF architecture, chat template, tokenizer metadata, special tokens, model config.
2. **Template content** — parse the native Jinja/ChatML structure and fingerprint it.
3. **Consistent family evidence** — architecture and tokenizer agree.
4. **Filename/marketing name** — weakest hint only; never sufficient when stronger evidence exists.
5. **Unknown** — safe generic fallback, visibly labelled, with manual override.

### 13.3 Evidence and confidence states

| State | Detection meaning | UI wording | Behavior |
|---|---|---|---|
| Verified | Metadata/template exercised on a real model | “Verified” | Automatic selection allowed |
| Inferred | Consistent metadata but not exercised | “Inferred” | Automatic selection allowed with warning |
| Conflicting | Sources disagree | “Conflict — review required” | Do not silently choose; require override or safe fallback |
| Unknown | No reliable evidence | “Unknown — default applied” | ChatML/generic fallback, clearly labelled |

A “confidently incorrect” template is worse than a labelled fallback. The design must prefer refusal/uncertainty over silent wrong formatting.

### 13.4 Compatibility matrix shape

Every supported family should have:

- recognition rule and evidence source;
- architecture identifiers/aliases;
- template source and fingerprint;
- required special tokens;
- reasoning-tag behavior;
- stop sequences;
- tool/function-call support if applicable;
- tokenizer caveats;
- verified model examples and test dates;
- fallback and override behavior;
- hot-swap behavior.

Candidate coverage to investigate includes Llama 2/3.x, Mistral/Mixtral, Qwen/Qwen2/2.5/3/3.5, Gemma 1–4, Phi-3/4, DeepSeek-R1 distills, LFM, ChatML-generic, Command R, Yi, Falcon, MPT, GLM, InternLM, Baichuan, EXAONE, OLMo, Nemotron, and emerging architectures. **This is a research target list, not a claim of support.**

### 13.5 Hot-swap transition

A switch must re-run detection for the target before load and verify the applied template after readiness. A cached template or family value from the old model is a correctness defect. A user override must survive reload but must be visibly associated with the model/family to which it applies; it must not silently follow the user to an incompatible model.

---

## 14. Sampler, offloading, context-shift, and optimization inventory

### 14.1 Runtime baseline

The pinned server is `llama.cpp 0.8.0-dev, build 12015, commit a34de8618`. The following options are present in that binary's help and are therefore the valid capability baseline for this inventory. Availability in the binary does not prove Howl passes or exposes them.

This is a **first-pass, single-backend evidence seed**, not a comprehensive survey and not a frozen capability catalogue. The full mission — other engines, architecture families, backend versions, conflicts, automatic selection, manual override, and benchmarks — is §15. The current research must not be treated as complete until that mission is executed.

### 14.2 Sampler inventory

| Capability | Pinned runtime | Howl UI | Howl request path | Status |
|---|---|---|---|---|
| Temperature | Yes | Yes | Yes | Exposed |
| Top-K | Yes | Yes | Yes | Exposed |
| Top-P | Yes | Yes | Yes | Exposed |
| Min-P | Yes | Yes | Yes | Exposed |
| Typical-P | Yes | Yes | Yes | Exposed |
| Repeat penalty / last-N | Yes | Yes | Yes | Exposed |
| Presence penalty | Yes | Yes | Yes | Exposed |
| Frequency penalty | Yes | Yes | Yes | Exposed |
| Mirostat 1/2, tau, eta | Yes | Yes | Yes | Exposed |
| Dynamic temperature range/exponent | Yes | Yes | Yes | Exposed |
| DRY multiplier/base/allowed length | Yes | Yes | Yes | Exposed |
| DRY last-N | Yes | No | No | Missing |
| DRY sequence breakers | Yes | No | No | Missing |
| XTC probability/threshold | Yes | No | No | Missing |
| Top-n-sigma | Yes | No | No | Missing |
| Adaptive-P target/decay | Yes | No | No | Missing |
| TFS | No explicit CLI flag recovered | No | Struct field only, not mapped | Unestablished / dead field |
| Smoothing factor/curve | Not confirmed as server CLI | No | Struct only, not sent | Unestablished / do not expose |
| Top-A | Not confirmed as server CLI | No | Struct only, not sent | Unestablished / do not expose |
| Epsilon/Eta cutoff | Not confirmed in pinned help | No | Struct only, not sent | Unestablished |
| No-repeat n-gram | Not confirmed in pinned help | No | Struct only, not sent | Unestablished |
| Encoder repeat penalty | Not confirmed in pinned help | No | Struct only, not sent | Unestablished |
| Sampler sequence/order | Yes (`--samplers`) | No | No | Missing |
| Seed | Yes | Numeric field | Yes, without actual-seed capture | Partial |
| Grammar / JSON schema | Yes | No chat control found | No | Missing/unscoped |

“Quadratic penalization” did not appear in the pinned binary's help under that name. It must not be claimed as supported without a runtime check. XTC and min-P are confirmed available in the pinned binary; XTC is not exposed by Howl.

Additional sampler work must be capability-gated by the actual backend/model/runtime. Controls should carry conflict metadata and should not imply that every filter is meaningful when Mirostat or another overriding mode is active.

### 14.3 Offloading and runtime controls

| Capability | Pinned runtime | Current Howl status | Assessment |
|---|---|---|---|
| GPU layer offload (`-ngl`) | Yes | Stored and passed | Implemented path |
| Flash Attention | Yes | Stored and passed | Implemented path; model/backend support not proven |
| Tensor split | Yes | Stored and passed | Basic multi-GPU path |
| Main GPU | Yes | Stored but not passed by `buildLlamaServerArgs` | Dead setting risk |
| KV offload | Yes | Stored but not passed | Dead setting risk |
| mmap/mlock/load mode | Yes | Stored but not passed; runtime default applies | Not controlled |
| KV cache type/quantization | Yes | Not exposed | Missing |
| SWA full cache | Yes | Not exposed | Missing/capability-gated |
| CPU MoE offload (`--n-cpu-moe`) | Yes | Not exposed | Candidate optimization |
| Tensor override (`--override-tensor`) | Yes | Not exposed | Advanced candidate |
| Device selection | Yes | Not exposed | Candidate |
| `--fit` / automatic fit | Yes | Used conditionally; disabled for SWA models | Partial; no evidence of correctness across families |
| Context shift | Yes | Not explicitly set or surfaced | Runtime-dependent; no product control |
| Cache reuse (`--cache-reuse`) | Yes | Not exposed | Missing; performance research required |
| Context checkpoints/SWA checkpoints | Yes | Not exposed | Candidate; not assumed applicable |

A setting appearing in `UpdateRuntimeSettings` is not evidence that it reaches the server. Current `buildLlamaServerArgs` passes only a subset. The corrected plan must distinguish **stored**, **sent**, **accepted by server**, and **verified effective** states.

### 14.4 Kobold-style capabilities

KoboldCpp documents fast-forwarding/prompt-cache reuse, context shift, SWA handling, quantized KV cache, `--lowvram`-style placement, smart-context/smart-cache concepts, and GPU layer offload. These are not proprietary requirements. The transferable pieces must be mapped to actual llama.cpp flags:

- prompt-prefix reuse → `cache_prompt` / `--cache-reuse` where supported;
- context shift → `--context-shift`, with SWA/position-encoding constraints;
- KV quantization → `--cache-type-k/v`, with Flash Attention requirements;
- offload → `--n-cpu-moe`, `--override-tensor`, `--device`, placement controls;
- smart-cache-style snapshots → not assumed to exist in the pinned server; possible app-side extension requiring separate design.

Do not copy a Kobold feature name into Howl and claim parity. Each item needs a pinned-runtime test.

### 14.5 Context shift versus Howl's current truncation

Current `truncateMessagesToFitContext` uses a rough four-characters-per-token estimate and drops older messages from the middle while keeping the system prompt and recent messages. It is not tokenizer-accurate and it is not context shift.

Context shift is a backend KV-cache mechanism for dropping/shifting old context when the window fills. It must not be presented as memory summarization, character isolation, or durable continuity. The app needs a policy that states which mechanism acts first and how the user is told.

### 14.6 LM Studio optimization research

No claim is made about LM Studio's internal implementation. The portable candidates worth measuring against Howl are:

- model load mode (mmap/mlock/direct I/O) and file-cache behavior;
- GPU layer placement and automatic fit;
- Flash Attention and KV-cache dtype;
- prompt-prefix/cache reuse;
- model warmup and readiness rather than process-start;
- keeping a model resident without reloading per turn;
- memory estimation that accounts for weights, KV cache, context, and compute buffers;
- avoiding duplicate model copies during a switch.

These are hypotheses for benchmarking, not verified LM Studio behaviour and not authorization to change the loader.

### 14.7 What this inventory does not establish

- It does not survey KoboldCpp, LM Studio internals, Ollama, vLLM, SGLang, TensorRT-LLM, MLC-LLM, ExLlamaV2/TabbyAPI, or other engines.
- It does not map techniques to dense, MoE, SWA, recurrent/SSM, hybrid, or multimodal architectures.
- It does not establish interactions or safe combinations.
- It does not measure performance, memory, stability, or output quality.
- It does not provide automatic selection, resource estimation, or manual-override warnings.
- It does not define conversation context rollover as distinct from multi-character flushing.

The catalogue must remain open until §15's research outputs exist.

---

## 15. Inference technique, architecture, and capability-research mission

### 15.1 Status and correction

Hardware-aware and architecture-aware inference is a **first-class architecture mission**, not an appendix under Flash Attention or a list of convenient llama.cpp flags.

The current inventory in §14 is a **first-pass pinned-runtime seed**, not a completed survey. It covers one binary, one point in time, and a subset of techniques. It must not be frozen as Howl's capability catalogue.

The mission is required; no individual optimization is thereby adopted. “Supported,” “safe,” “faster,” and “better” are separate claims requiring separate evidence.

This is **not** a mandate to enable every technique or implement every backend. It is a mandate to possess the capability to investigate, select, explain, and validate safely. The staged separation is defined in §15.12.

The product goal is not to expose every backend option. It is for Howl to:

- discover the machine, backend, model architecture, and workload constraints;
- know which techniques are actually applicable;
- choose a safe combination with an explanation;
- let advanced users override with informed warnings;
- measure the result rather than trusting parameter acceptance.

### 15.2 Comprehensive technique inventory

The inventory must span the full local-inference stack, not just samplers and GPU layers.

| Domain | Techniques and constraints to inventory |
|---|---|
| CPU execution | Thread count, affinity, NUMA, instruction sets, CPU-only operation, hybrid CPU/GPU placement |
| GPU execution | Device selection, main GPU, layer offload, full offload, unified memory, driver/runtime constraints |
| Multi-GPU | Tensor split, pipeline/parallel placement, interconnect, per-device memory, topology |
| Model memory | mmap, mlock, direct I/O, page cache, pinned memory, model residency, duplicate-copy avoidance |
| Quantization | Weight formats, mixed quantization, KV-cache dtype, kernel support, quality effects, architecture compatibility |
| Attention | Flash Attention, paged/contiguous KV, sliding-window attention, attention sinks, RoPE scaling, context checkpoints |
| Context | Context size, context shift, cache reuse, prompt-prefix reuse, SWA full cache, recurrent state, position handling |
| Batching | Batch size, micro-batch, parallel sequences/slots, continuous batching, prefill/decode tradeoffs |
| Model architecture | Dense transformer, MoE, sliding-window, recurrent/state-space, hybrid, multimodal projector, embedding/reranker |
| Model-specific features | Draft model/speculative decoding, grammar constraints, tool-call templates, reasoning modes, multimodal encoders |
| Storage/loading | Split model files, sidecars, adapters/LoRA, model hash, cold/warm load, disk throughput, network roots |
| Workload/lifecycle | Chat latency, long-context prefill, multi-character switching, image/video coexistence, hot-swap, idle release |
| Quality | Output degradation, sampler interactions, quantization effects, context loss, reasoning instability |

Each technique record must include:

| Field | Required content |
|---|---|
| Technique ID and name | Stable identity independent of a marketing name |
| Engine and version | Exact binary/fork/version/commit; “llama.cpp-compatible” is insufficient |
| Platform | OS, CPU architecture, GPU vendor/runtime, driver constraints |
| Model applicability | Architectures, quantization formats, modalities, context shapes |
| Support state | Native, app-level, emulated, version-dependent, unsupported, or unknown |
| Required settings | Flags/API fields and any prerequisite combination |
| Resource effect | Weights, KV/state, compute buffers, host RAM, VRAM, disk/page cache |
| Quality effect | Known or suspected output impact, including “not measured” |
| Interactions/conflicts | What it requires, disables, invalidates, or makes unsafe |
| Evidence | Source/docs, runtime probe, benchmark, or none |
| Limitations | Hardware, model, context, concurrency, stability, and version caveats |

### 15.3 Architecture mapping

Architectures are not interchangeable. A technique that is correct for a dense transformer can be meaningless or harmful for a recurrent, state-space, sliding-window, or hybrid model.

| Architecture family | What must be mapped | Primary risks of a dense-transformer assumption |
|---|---|---|
| Dense transformer | Full KV growth, layer offload, KV quantization, Flash Attention, context shift, batching | Treating these as universal |
| Mixture-of-Experts (MoE) | Total vs active parameters, router/expert placement, expert offload, tensor override, per-layer memory | Estimating memory or speed from active parameters alone; offloading only “layers” |
| Sliding-window attention (SWA) | Window size, layer pattern, full-cache option, checkpoints, position handling, cache shape | Assuming a full-context KV cache; applying context shift without SWA semantics |
| Recurrent/state-space | Recurrent state size, chunking, state carry, batching, context-extension method | Assuming transformer KV cache or attention flags apply |
| Hybrid attention/recurrent | Per-layer mixture, which layers hold KV vs state, offload granularity | Applying one cache/context policy to the whole model |
| Multimodal | Projector/vision/audio encoder memory, separate preprocess, image-token cost, backend support | Counting only language-model weights/context |
| Embedding/reranker | Pooling, sequence length, batch, non-generative runtime, vector-store interaction | Loading it as a chat model or applying sampler settings |

Model classification must use model metadata and architecture evidence, not filename alone. The existing F2 model-intelligence evidence classes are a foundation, but they do not prove runtime compatibility, kernel support, or memory fit.

### 15.4 Backend and version mapping

Support must be recorded **per backend, version, platform, architecture, quantization, and feature**. A capability is not portable because two engines share a model format or because one wraps another.

Candidate survey set — not a claim that Howl supports these engines:

| Backend family | Why it belongs in the survey | Required caution |
|---|---|---|
| llama.cpp | Current Howl backend; broad local formats and offload controls | Pinned build/commit matters; flags and defaults change |
| KoboldCpp | Local inference UI with context/cache/offload techniques | Bundles/forks components; verify its exact build and semantics; no feature-name parity |
| LM Studio | Relevant load-speed, placement, cache, and UX behavior | Proprietary implementation; observable behavior only, no internal assumptions |
| Ollama | Popular local runtime and model packaging | Wrapper/runtime version and hidden defaults; server API may not expose controls |
| vLLM | High-throughput serving, batching, paged attention | Different hardware/OS assumptions; server-oriented behavior |
| SGLang | Serving/runtime optimizations and structured generation | Version/feature support and hardware constraints |
| TensorRT-LLM | NVIDIA-specific compiled engines and quantization | Compilation/version/GPU coupling; not a drop-in local chat runtime |
| MLC-LLM | Compiled cross-platform runtimes | Model compilation and device-specific artifacts |
| ExLlamaV2 / TabbyAPI | GPU-local inference with EXL2 quantization | Format/GPU constraints; not a generic GGUF path |
| LocalAI and similar wrappers | Aggregated local backends and APIs | Wrapper abstraction can hide which engine/version actually executes |
| Other local engines | Future architectures or hardware-specific paths | Must be evaluated, not assumed replaceable |

A backend support record must state:

- engine, exact version/commit, and execution API;
- OS/CPU/GPU/driver requirements;
- model formats, architectures, quantization, and multimodal support;
- feature availability and prerequisites;
- whether the control is native, app-level, or only a wrapper claim;
- memory/performance behavior on named hardware;
- known limitations and failure modes;
- whether effective application was verified, not merely accepted.

The F2 reports establish a useful pattern: identification, candidate runtimes, and capabilities are separate from load compatibility. They explicitly leave tokenizer execution, template application, request-level capability negotiation, and real model-load verification incomplete. That is a foundation for this mission, not evidence that the mission is done.

### 15.5 Interactions and conflicts

Optimizations cannot be selected independently. The catalogue must record a conflict graph and a tested-combination state.

| ID | Interaction/conflict | Required treatment |
|---|---|---|
| C-01 | Flash Attention and KV-cache quantization | Record prerequisite/version constraints and quality effects; do not assume both are always compatible |
| C-02 | Context shift and SWA/checkpoint behavior | Determine whether shift is valid with the model's window pattern; do not apply a full-cache assumption |
| C-03 | Prompt-cache reuse and character-context flush / prefix mutation | Define invalidation; a reused prefix must not defeat multi-character isolation |
| C-04 | mmap, mlock, direct I/O, and page cache | Measure RAM, load time, eviction, and interaction with model switching |
| C-05 | GPU layer offload and context/KV/batch size | Fit is multidimensional; layers alone do not determine VRAM |
| C-06 | Tensor split, main GPU, and interconnect | Topology can make a nominal fit slower or unstable |
| C-07 | MoE expert offload and dense layer/tensor offload | Architecture-specific; estimate total weights and expert placement separately |
| C-08 | Batch/ubatch and parallel sequences | Context memory multiplies; throughput gains can increase latency |
| C-09 | Speculative decoding and draft-model compatibility | Draft vocabulary/template/architecture must match; extra residency consumes memory |
| C-10 | Quantization and kernels/Flash Attention | A quantized format can lose fast kernels or change quality |
| C-11 | Automatic fit and compute-buffer/concurrency reserve | Underestimating non-weight memory causes OOM after apparent fit |
| C-12 | Model hot-swap and residency/release | A new model must not be admitted on an unproven release; simultaneous residency is not assumed |
| C-13 | Cache reuse and single/multi-character chat | Single-character continuation may benefit; multi-character isolation may invalidate the cache |
| C-14 | Sampler stack, grammar, and reasoning mode | Individually valid controls can conflict or distort output; quality must be measured |

For every pair, the catalogue must state: compatible, conditional, conflicting, unknown, or untested. “Both flags are accepted” is not a compatibility result.

### 15.6 Automatic capability selection

Howl should select a safe configuration rather than require the user to understand the entire stack.

Required inputs:

- hardware: CPU, RAM, GPU devices, VRAM, driver/runtime, storage, available memory, OS constraints;
- backend: available binaries, versions, kernels, supported formats, known defects;
- model: architecture, quantization, parameter count, context limit, multimodal components, adapters;
- workload: chat, long context, multi-character, image/video coexistence, hot-swap, concurrency;
- policy: quality floor, latency/throughput preference, safety reserve, local-only requirement.

Selection procedure:

1. Resolve hard incompatibilities first (unsupported architecture/format/kernel/platform).
2. Estimate weights, KV/state, compute buffers, multimodal components, and runtime overhead.
3. Apply a safety reserve for the OS, driver, other modalities, and context growth.
4. Filter combinations that cannot fit or violate a quality/safety policy.
5. Rank remaining candidates by predicted latency, throughput, memory headroom, and stability.
6. Prefer the least complex configuration that meets the workload.
7. Load and verify effective settings, model identity, and readiness.
8. Fall back in a defined order on failure; do not silently leave the prior model half-released.
9. Record the selected configuration and the reason for exclusions.
10. Expose an explanation when an expected optimization is unavailable.

The selector may use conservative estimates, but it must label uncertainty. It must not claim exact VRAM/RAM prediction before measurement.

### 15.7 Manual control and override

Advanced users must be able to override automatic choices, but not silently bypass safety or invent support.

| Requirement | Behavior |
|---|---|
| Scope | Overrides apply to a specific model/backend/hardware profile, not globally by accident |
| Visibility | Show auto-selected value, override value, and source |
| Risk warning | Explain instability, OOM, quality loss, unsupported combination, or slow path |
| Validation | Reject values the backend/model cannot accept; do not send dead parameters |
| Effective verification | Show whether the backend applied the setting, not just that it was requested |
| Reset | Restore automatic selection without losing the model |
| Provenance | Record overrides used for a generation |
| Safety | An override cannot bypass hot-swap readiness, resource release, or data-integrity rules |
| Persistence | Preserve valid overrides per model/backend where appropriate; invalidate on changed fingerprint |

A warning is not a guarantee. If an override is known to corrupt output or crash the backend, Howl should refuse it rather than merely warn.

### 15.8 Benchmark rather than assume

A technique is not successful because a backend accepts its parameter. It must be measured against the unoptimized baseline on named hardware.

Minimum measurements:

| Dimension | Measurements |
|---|---|
| Load | Cold load, warm load, time to first token, failure rate |
| Prefill | Prompt lengths, long-context prefill, cache reuse effect |
| Decode | Tokens/second, latency distribution, batch/parallel effect |
| Memory | Peak host RAM, peak VRAM, KV/state share, page-cache behavior, headroom |
| Stability | OOM, crash, hang, repeated-switch reliability, thermal/clock effects |
| Quality | Controlled prompts, regression suite, blinded comparison where feasible, reasoning/tool-call degradation |
| Switching | Unload/release, load, readiness, first generation, rollback |
| Coexistence | Chat plus image/video/TTS or studio jobs, if those are resident |

The benchmark record must include hardware, OS, driver, backend version/commit, model hash, quantization, context, batch/parallelism, sampler settings, and exact test procedure. Report distributions and repeated runs, not a single best number. Preserve negative and inconclusive results.

Output quality is part of the result. A faster configuration that breaks character consistency, instruction following, structured output, or long-context recall is not an unconditional improvement.

### 15.9 Context rollover versus multi-character context flushing

These are separate features. They are related through context and memory management, but neither is a synonym for the other.

| Dimension | Conversation context rollover | Multi-character context flushing |
|---|---|---|
| Problem | The conversation exceeds available context | Character A's inference context could bleed into character B |
| Trigger | Capacity/rollover threshold | After each character response in a multi-character chat |
| Scope | One conversation's history/prompt budget | Active inference context for the next speaker |
| Mechanisms | Context shift, summarization, retention policy, rebuilding, memory reinjection | Flush/invalidate/rebuild the character-scoped inference context |
| Frequency | At or near context pressure | Every character turn in multi-character chat |
| Single-character | Applies according to configured capacity/rollover policy | **Must not flush after every response** |
| Durable state | Transcript, memory, lore, and character state must survive | Transcript, memory, lore, and character state must survive |
| Not the same as | Character isolation, model hot-swap, transcript deletion | Context shift, summarization, memory deletion, model unload |

The required pipeline distinction is:

```
turn start
  → build token budget
  → retrieve lore/memory for this character/scene
  → assemble canonical transcript + scoped state
  → if multi-character, ensure a fresh character context
  → generate
  → commit response
  → if multi-character, flush/invalidate that character's active context
  → next character

separately, when conversation context is near capacity:
  → apply conversation rollover policy (shift/summarize/retain/rebuild)
```

Current source exposes a “Context Rollover Mode” setting and tutorial text describing Memory-Guided, Standard, and Sliding Window behavior. The actual implementation was not runtime-verified for this addendum, and the existing truncation path is character-count-based rather than tokenizer-accurate. The UI label and tutorial claims must not be treated as proof that rollover semantics are implemented. The multi-character requirement is independently unestablished.

### 15.10 Required research outputs

The mission must produce, before a capability catalogue is frozen:

1. a technique register with evidence and limitations;
2. an architecture applicability matrix;
3. a backend/version/platform support matrix;
4. an interaction/conflict graph with tested combinations;
5. an automatic-selection design with resource estimation and explanation;
6. a manual-override and warning specification;
7. a benchmark protocol with named hardware and quality checks;
8. a negative-results register;
9. a context-rollover design distinct from multi-character flushing;
10. a traceability link from every exposed control to evidence and provenance.

No catalogue entry may be marked “supported” solely because one engine accepts a flag. No “optimization” may be marked successful without a measured baseline. No architecture-specific technique may be generalized from a dense transformer without evidence.

### 15.11 Current foundation and gap

The F2 model-intelligence layer is a useful start: it records profile evidence, capabilities, candidate runtimes, sidecar/template presence, and conflicts, and it separates identification from runtime compatibility. The F2-B report confirms that the Model Browser exposes this metadata, while tokenizer execution, template application, request-level capability negotiation, and live model-load verification remain deferred.

Current Howl passes only a subset of stored runtime settings to the server, has no architecture-aware estimator, no effective-capability verification, no conflict graph, no benchmark harness, and no automatic selector. The pinned-runtime inventory in §14 is the first seed for this mission, not its conclusion.

### 15.12 Capability mission is not an enablement mandate

The mission requires Howl to possess a **capability-management system**. It does not require Howl to expose or enable every technique, and it does not authorize implementation of every surveyed backend feature.

The stages are separate and ordered:

| Stage | Question | Output | Does not imply |
|---|---|---|---|
| Research | What exists, in which engine/version, for which architecture/hardware? | Technique/backend evidence records | Support or enablement |
| Evidence | What actually happens on named hardware/models? | Runtime probes and benchmarks, including quality | Automatic selection or release |
| Compatibility | Which combinations are valid, conflicting, or unknown? | Conflict graph and tested-combination states | That a safe default is implemented |
| Selection | Which safe configuration should be chosen and why? | Estimator, ranking, fallback, explanation | That the user may not override |
| Validation | Does the selected/overridden configuration work and preserve quality? | Effective-setting, stability, and quality checks | Acceptance of the feature |
| Authorization | Is implementation of this technique/backend approved? | Explicit scope/decision | Any other technique or backend |

A technique appearing in the catalogue, being accepted by a backend, or passing one benchmark does not advance through all six stages automatically. No individual optimization may be enabled merely because the architecture mission exists. No implementation may leapfrog the recovery-baseline reconciliation.

---

## 16. General creative-studio modality architecture

### 16.1 Approved direction

Howl must not assume that chat, images, and audio are the permanent limit. Music generation, 3D mesh generation, texture generation, and future local generative modalities are within the approved product direction. Approval of the direction is not approval of any immediate implementation.

### 16.2 Shared foundations versus modality-specific work

| Foundation | Shared? | Why |
|---|---|---|
| Asset/media identity and catalog | Yes | Every modality produces durable objects |
| Job/operation contract | Yes | Request, progress, cancel, failure, result, recovery |
| Resource admission and placement | Yes | One desktop GPU/RAM budget |
| Engine adapter boundary | Yes as an interface, not one engine | Different engines and formats |
| Provenance record | Yes as a pattern, modality-specific fields | Model/version/seed/settings differ |
| Workspace/UI shell | Shared patterns, separate modality workspaces | F2-B5 requires separate interfaces |
| Library roles/export | Yes | Files must be findable and portable |
| Model dependency resolution | Yes | Missing VAE/CLIP/encoder/LoRA assets must fail clearly |
| Prompt/template logic | Modality-specific | Text, image, music, and 3D prompts are not interchangeable |
| Sampler/decoder controls | Modality-specific | Diffusion, audio codec, and mesh decoders differ |
| Output validation | Modality-specific | Image dimensions, audio duration, mesh topology, texture UVs |

### 16.3 Modality scope table

| Modality | Direction | Required first scope | Explicit non-claim |
|---|---|---|---|
| Image | Existing | Preserve current txt2img/img2img/inpaint/outpaint/upscale, workflows, gallery | Existing runtime acceptance is still incomplete |
| Video | Existing/incomplete | Define sidecar/catalogue/delete contract (F3-G) before chat attachment | Not complete because controls exist |
| Music/audio generation | Approved direction | Separate from TTS; model format, duration, hardware, playback, export, provenance | No engine selected; F2-B7 remains frozen for TTS |
| 3D mesh | Approved direction | Input format, topology/scale validation, preview, export, storage | Not a desktop avatar; no automatic character rigging |
| Texture | Approved direction | Source mesh/UV validation, reference image, output maps, preview/export | Not proof that an external texturing service is bundled |
| Other local modalities | Approved direction | Capability discovery, model/resource requirements, workspace, persistence | No generic “plugin” promise |

### 16.4 Studio constraints

- Local-first operation remains mandatory.
- Each modality must declare its model formats, dependencies, VRAM/RAM profile, and failure behavior.
- No modality may assume that the chat inference backend is the same backend.
- No modality may silently merge into the chat interface; F2-B5's separate workspace direction remains.
- No modality is added to the release gate merely because the direction is approved.
- 3D asset generation must not be used to resurrect the excluded desktop-avatar surface.
- E.6/E.8 remain parked and are not automatically promoted into the studio.

### 16.5 Current delta

Image and video have real engine code. Music, 3D mesh, and texture generation have no production subsystem found. The studio is an approved architectural direction, not an implemented product surface.

---

## 17. Corrected product-versus-contract classification

### 17.1 The classification rule

A product decision does not automatically become:

- an F3-C target-set field;
- an F3-C amendment;
- an accepted schema;
- a migration;
- a release commitment;
- an implementation authorization.

Conversely, an F3-C proposal appearing in a review does not become a product requirement merely because the review discussed it. The 2026-11-11 owner instruction now supplies product authority for several concepts. That authority changes their **product** classification, not their **contract** adoption.

### 17.2 Corrected classification table

| Item | Product classification | Contract classification | Implementation classification |
|---|---|---|---|
| Desktop character avatar | **Excluded** | No contract | No work |
| Chat portrait thumbnail | Required | UI/chat presentation; no F3-C schema requirement | Partial source, not complete |
| Model hot-swap | **Required product behavior** | Not an F3-C target-set requirement; no amendment | Manual partial; unsafe/incomplete |
| Per-chat continuation-model preference | **Unresolved product question** | Not adopted; separate from provenance | Absent |
| Per-turn provenance | **Required concept** | Exact schema not adopted; no amendment | Absent |
| Actual seed preservation | **Owner-directed if meaningful** | Exact fields not adopted | Absent |
| Durable attachments | **Approved product requirement** | Exact schema/media contract not adopted | Absent/transient |
| Multi-character flush | **Required behavior** | F3-C explicitly does not define/change it | Not established |
| Single-character no-flush | **Required behavior** | F3-C does not change it | No per-turn flush exists; must remain |
| Memory subsystem | **Required major capability** | No F3-C adoption | Partial foundation only |
| General creative studio | **Approved direction** | No contract | Absent as unified subsystem |
| F2-B5 separate workspaces | **Approved** | Existing approved direction | Partial/deferred UI filtering |
| E.6/E.8 | Parked | No contract | None |
| REG-1 content hot-swap | Protected behavior | Existing regression constraint | Existing behavior |
| Inference capability catalogue | **Required first-class research mission** | No contract; explicitly unfrozen | First-pass pinned-runtime evidence only |
| Network boundary and offline enforcement | **Required first-class product/architecture boundary** | No F3-C contract | Partial controls; backend enforcement missing |
| External destination registry, consent, and audit | **Required** | No contract | Missing |
| Automatic Civitai model-hash lookup | **Prohibited in the local baseline unless separately authorized** | No contract | Implemented and called from image scan; non-compliant |
| Remote ASR/TTS as part of local inference | **Not part of the local baseline** | No contract | Remote paths exist without operation-specific authorization |
| Android remote | **Permitted only if explicitly configured and secured; currently unauthorized** | MOB constraints preserved | Not implemented |
| Automatic inference selection | **Required architecture behavior** | No contract | Absent |
| Manual inference override | **Required advanced-user control** | No contract | Partial raw settings only |
| Conversation context rollover | **Separate required behavior** | No contract | UI/tutorial claims; runtime unverified |
| Author's Note | Deferred real per-chat requirement | No F3-C contract | Absent |
| F2 model intelligence | Accepted foundation | Identification is not runtime compatibility | Partial |
| Individual optimization/backend feature | Not adopted by the mission | Requires evidence, compatibility, selection, validation, and authorization | None implied |
| Chat bubble/theme/portrait UI | Owner-confirmed requirement | No F3-C schema | Partial or absent |
| Narration | Owner-confirmed requirement; semantics unrecovered | No contract | Absent as a feature |
| Token counter | Owner-confirmed requirement | No contract | Absent |
| Message-triggered image generation | Owner-confirmed requirement | No contract | Absent |
| One live prompt path | Architecture prerequisite | F3-E cannot be complete with split owners | Split `App`/`chat.Service` |
| FUT-CH1/FUT-CH2/FUT-C1 | Deferred real per-chat requirements | Must be considered in durable schema | Not implemented |
| F2-B4 model/template evidence | Accepted requirement | Existing model-intelligence contract | Partial |
| FUT-C5 sampler provenance/research | Deferred requirement | No contract | Partial sampler UI |
| FUT-ST2 import provenance | Deferred requirement | No contract | Partial import work |
| Image metadata | Existing behavior | No chat contract | Source-confirmed |
| W4 | Historical blocker | Superseded in substance by W5–W7 | No live blocker; formal closure evidence gap |
| W8 | Verified | Accepted record | Verified |
| MOB-1–MOB-4 | Binding conditional constraints | Protected product boundary | No implementation |
| HW-C1–C3 | Deferred validation targets | No contract | Not started |
| FUT-HW1 | Deferred foundation | No contract | Partial resource manager |
| FUT-PROMPT1 | Required direction | No contract | Partial |

### 17.3 Specific classification corrections

**Model hot-swap.** Product: required. F3-C: separate. Contract: no change. Implementation: not safe. The corrected plan must not repeat “unaccepted” as the product status, but must repeat “not an F3-C amendment” as the contract status.

**Per-chat model preference.** The owner has required switching during a conversation, but has not decided whether each chat stores a preferred continuation model, whether every assistant turn carries the model forward, or what happens on reopen. Historical evidence does not settle it. It remains a distinct unresolved question. Active switching is not the same as a per-chat default.

**Provenance and seed.** Product: required concepts. Contract: no adopted fields. Engineering: recommendation required. The owner did not approve all possible fields, a particular storage format, or exact reproduction. The prior “proposal-only, owner decision needed” label is superseded at the concept level but remains accurate at the schema level.

**Attachments.** Product: approved. Contract: no adopted attachment schema and no F3-C amendment. Engineering: durable identity, input/display distinction, and lifecycle required.

**Memory.** Product: required major capability. Contract: no F3-C adoption. Engineering: architecture and staged implementation required. FUT-MEM1 is no longer merely an optional deferred idea, but it is still not implemented or authorized.

**Creative studio.** Product: approved direction. Contract: no contract. Engineering: no modality implementation. F2-B5 remains approved; E.6/E.8 remain parked; the general direction is separate from both.

**Desktop avatar.** Product: excluded. Historical `docs/tier4-avatar.md` remains evidence but is not a roadmap item. Chat portraits are a different requirement.

**Chat UI.** Product: owner-confirmed. Contract: no F3-C schema. Engineering: partial or absent; the exact separate UI planning artifact is not recovered.

**Inference capability mission.** Product: required first-class architecture mission. Contract: no contract and no frozen catalogue. Engineering: a first-pass llama.cpp inventory plus model-intelligence metadata, not a hardware/architecture-aware selector or benchmark system.

**Conversation context rollover.** Product: required separately from multi-character flushing. Contract: no contract. Engineering: UI/tutorial claims exist; tokenizer-accurate rollover behavior is unverified.

**Author's Note.** Product: deferred real per-chat requirement with recovered presentation details. Contract: no F3-C schema. Engineering: absent.

---

## 18. Dependency and implementation-order implications

### 18.1 Dependency graph

```
Stable identity + canonical transcript
        ├── durable chat schema (incl. per-chat union)
        ├── media identity + attachment contract
        └── per-turn provenance

One live prompt path
        ├── tokenizer/context accounting
        ├── lore + RAG + memory selection
        ├── character context isolation/flush
        └── model-family template selection

Model identity/template evidence
        ├── family detection
        ├── hot-swap target preflight
        └── readiness/identity verification

Resource orchestration
        ├── unload/release verification
        ├── load admission
        └── hot-swap safety + future modalities

Memory schema/lifecycle
        ├── retrieval/ranking/compression
        ├── prompt reinjection
        └── multi-character scope isolation

Media catalog/identity
        ├── chat image references
        ├── chat audio references
        ├── export/import/restore
        └── studio output library

General studio shell
        ├── modality workspaces
        ├── shared job/result contract
        ├── model/dependency resolution
        └── resource admission per modality

Network policy gate
        ├── destination registry
        ├── six-question consent and scoped permissions
        ├── local connection audit
        ├── backend/sidecar offline enforcement
        └── local-versus-external classification
```

Capability research is also a cross-cutting dependency chain: model architecture and backend/version evidence feed the conflict graph, which feeds automatic selection, which feeds benchmarks. It cannot be completed after hot-swap merely by documenting the flags the loader happened to use.

Network policy is likewise not a late privacy pass. It is a prerequisite for every feature that can reach a destination, because consent, audit, cache, failure, and offline behavior change the feature architecture rather than decorating it.

### 18.2 Required-before-work ordering

The following are planning prerequisites. They do not authorize implementation and do not reorder accepted F3 work.

1. **Preserve the protected baseline.** Do not close F3-C Phase 0, R6, F3-D, D-ROOT, D-STORAGE, V20, or any acceptance gate to make room for the expanded product.
2. **Recover the one live prompt path.** RAG, lore, memory, tokenizer/context, character isolation, and template selection cannot have two owners. F3-E cannot be complete while `App` and `chat.Service` remain separate prompt-relevant authorities.
3. **Define stable chat/turn identity.** Durable provenance and attachment references require stable turn/message identity; index-derived IDs are insufficient.
4. **Define the per-chat schema union.** FUT-CH1, FUT-CH2, FUT-C1, the unresolved model preference, provenance, and attachments must be considered together before migration.
5. **Define media identity and lifecycle.** Chat references must point to durable objects; display and input must be separate.
6. **Define memory identity/scope/schema.** Retrieval cannot be designed on untyped prose alone.
7. **Define tokenizer/context accounting.** The memory/lore budget, token counter, and single-character context capacity all depend on it.
8. **Define model-family evidence and template fingerprints.** Hot-swap must not reuse the old family's format.
9. **Define model lifecycle/readiness/release.** Hot-swap cannot be accepted on process exit and `/health` alone.
10. **Define the shared studio job/media contract.** A new modality must not invent a parallel gallery/media universe.
11. **Establish capability evidence before exposing controls.** Architecture, backend/version, conflicts, resource estimation, and effective-setting verification must precede a frozen catalogue or broad advanced controls.
12. **Keep research, selection, validation, and authorization separate.** The capability mission is not permission to enable a technique, and no backend feature is adopted by appearing in the catalogue.
13. **Define the network boundary before adding any external path.** The destination registry, six-question consent, purpose scope, audit, cache/failure behavior, and backend offline policy are architectural prerequisites. They cannot be retrofitted after a provider is already wired into a feature.

### 18.3 Parallel work that does not require F3 resequencing

- UI visual repairs can be investigated without altering persistence.
- Model-family and sampler inventories can be evidence work in parallel.
- LM Studio/Kobold comparisons can be measured in isolation.
- Memory schema research can proceed without implementing it.
- Studio modality research can proceed without selecting release scope.
- The shutdown/lifecycle guard can be verified in parallel; it cannot be deferred until after evidence that depends on it.
- Engine and architecture surveys can run in parallel with the current llama.cpp inventory; they must not be represented as complete until benchmarked.
- Author's Note UI/persistence design can proceed from the recovered presentation requirements without changing F3-C.
- Network-path inventory, consent wording, and audit design can proceed in parallel; no external destination is thereby authorized.

### 18.4 Cycles and breakers

| Cycle | Breaker |
|---|---|
| Hot-swap needs resource release; resource manager needs lifecycle authority | Define one model-operation owner and one state machine before either is complete |
| Memory needs tokenizer budget; tokenizer needs prompt path | Designate prompt path first, then token accounting |
| Media needs chat identity; chat schema needs media fields | Define stable turn/media IDs together |
| UI needs per-character settings; schema needs settings ownership | Decide character presentation state as part of the character/chat contract |
| Studio needs media catalog; media catalog is scoped to chat/images | Define shared media identity before modality expansion |
| Model family needs template evidence; template needs runtime metadata | Use evidence states and a safe fallback; never filename certainty |

---

## 19. Owner decisions still unresolved

Only the following remain genuine owner/product decisions that evidence cannot settle. They are not blockers to recording the baseline.

| ID | Question | Why it remains open | Safe default if deferred |
|---|---|---|---|
| UD-01 | Does each chat store a preferred continuation model, and how is it restored on reopen? | Active hot-swap is decided; per-chat preference/restoration is not. Historical PROP-F3C-02 is a proposal, not an accepted requirement. | Do not store a per-chat pointer yet; use the global selection and record actual per-turn models. |
| UD-02 | Is an Android/remote phase wanted; if so, what pairing/auth/API contract applies? | MOB-1–MOB-4 constrain any future phase but do not authorize it; no protocol exists. | Keep remote unauthorized and optional; desktop remains complete. |
| UD-03 | What is the committed Story Tools scope? | Page and controls exist; intended committed scope is not established. | Preserve existing controls; no new commitment. |
| UD-04 | What exact narration behavior is intended? | Owner confirmed SillyTavern-like purpose, but the historical artifact was not recovered and current “narrator events” are different. | Do not invent a narration system; recover/confirm semantics before design. |
| UD-05 | What release priority/scope applies to each approved studio modality? | Direction is approved; modality requirements, hardware, and sequencing are not. | No modality implementation; plan shared foundations only. |
| UD-06 | What retention/review policy should memory use for model-extracted candidates? | Required capability is decided; how intrusive review is, auto-accept thresholds, and retention are product policy. | Default to reviewable candidates, no silent lore overwrite, conservative retention. |

### 19.1 Questions that are not owner decisions

These are engineering/evidence questions and must not be sent to the owner as product choices:

- whether llama.cpp exposes the actual seed on Howl's endpoint;
- which prompt path is technically live;
- whether a given sampler is supported by the pinned build;
- how to fingerprint a template;
- which media hash/locator format to use;
- whether a process exit released VRAM;
- how to test character-context isolation;
- the tokenizer implementation;
- the exact cache keys and job IDs;
- whether a UI control should be a slider or numeric input.

---

## 20. Technical investigations still required

| ID | Investigation | Blocks | Required evidence |
|---|---|---|---|
| TI-01 | Recover the separate historical UI planning document and expansion text, or establish that they are unavailable | Exact UI wording/traceability | Artifact identity or documented absence; no invented quote |
| TI-02 | Define narration semantics in Howl: representation, speaker, display, prompt role, multi-character behavior | Narration design | Recovered source or explicit owner clarification |
| TI-03 | Establish one live prompt-assembly path and owner for system, lore, RAG, memory, cues, tokenizer, and character isolation | F3-E completeness; memory; context | Source trace and runtime/prompt-boundary evidence |
| TI-04 | Implement/verify tokenizer-accurate context accounting and token counter values | FUT-CHTOK1; memory budget; FUT-CH2 | Tokenizer evidence per active model; no false precision |
| TI-05 | Define durable chat/turn schema and migration/rollback | Provenance, attachments, per-chat settings | Adopted schema decision; migration tests |
| TI-06 | Define media identity, ownership, unavailable-reference, export/import/restore, and cleanup semantics | Durable attachments; studio | Contract and failure tests |
| TI-07 | Prove hot-swap concurrency, readiness, identity, unload/release, failure, and recovery | Model hot-swap | Runtime tests with two models and active/pending message |
| TI-08 | Determine actual-seed availability on Howl's endpoint and test reproduction limits | Seed/provenance | Runtime probe on pinned build; model variance results |
| TI-09 | Define model identity/fingerprint and template fingerprint that survive uninstall | Provenance; hot-swap | Stable identity design; unavailable-model behavior |
| TI-10 | Build the model-family/template compatibility matrix with confidence/fallback | FUT-PROMPT1 | Metadata/template evidence and real-model tests |
| TI-11 | Complete the pinned sampler capability matrix and conflict model | FUT-C5; owner sampler requirement | Backend/UI/request/effective verification |
| TI-12 | Measure offload, context shift, cache reuse, LM Studio/Kobold-portable optimizations | Performance; context | Benchmarks on named hardware; no parity claims |
| TI-13 | Design/validate memory extraction, retrieval, compression, reinjection, and review | Memory subsystem | Schema, retrieval quality, budget, isolation tests |
| TI-14 | Prove per-character context isolation/flush and single-character no-flush | Character continuity | Multi-character and single-character generation traces |
| TI-15 | Define per-character appearance state and theme interaction | Chat UI | Persistence/state contract, default/override behavior |
| TI-16 | Define message-triggered image-generation orchestration and media association | Chat image generation | End-to-end request using normal image path and safeguards |
| TI-17 | Define shared studio job/media/resource contracts and modality requirements | Studio | Modality matrix; no single-backend assumption |
| TI-18 | Reconcile all new product requirements with F3-C/F3-E without amendment | Corrected plan | Explicit mapping; no contract drift |
| TI-19 | Survey llama.cpp, KoboldCpp, LM Studio, Ollama, vLLM, SGLang, TensorRT-LLM, MLC-LLM, ExLlamaV2/TabbyAPI, and other relevant local engines | Capability catalogue | Per-engine version/platform evidence; no unsupported claims |
| TI-20 | Build the architecture applicability matrix for dense, MoE, SWA, recurrent/SSM, hybrid, multimodal, and embedding/reranker models | Capability catalogue | Architecture-specific cache/state/offload/context rules |
| TI-21 | Record backend/version/hardware support and limitations for every technique | Capability catalogue | Native/app-level/unsupported/unknown states; effective verification |
| TI-22 | Build and test the optimization interaction/conflict graph | Selection; loader | Tested combinations, invalidations, quality effects |
| TI-23 | Design automatic capability detection, resource estimation, candidate ranking, fallback, and explanation | Hardware-aware Howl | Conservative estimates and observable reasons for exclusions |
| TI-24 | Design scoped manual override, risk warnings, validation, reset, and provenance | Advanced controls | No bypass of safety/readiness; no dead parameters |
| TI-25 | Build the benchmark protocol for load, prefill, decode, memory, stability, switching, coexistence, and output quality | Capability claims | Named hardware, baselines, distributions, negative results |
| TI-26 | Define and verify conversation context rollover separately from multi-character flushing | Context architecture | Tokenizer budget, shift/summarize/retain/rebuild, observable semantics |
| TI-27 | Define Author's Note flyout/hidden/floating/side presentation, persistence, insertion order, enable state, and new-chat clearing | Chat UI/per-chat schema | Source-supported behavior and persistence tests |
| TI-28 | Complete the destination registry and static egress inventory, including proxy, sidecar, browser-open, and webview paths | Network mission | Every production path mapped to a registry entry or documented local/test exclusion |
| TI-29 | Design and verify backend offline enforcement, scoped consent, cache, audit, and failure behavior | Network mission | Egress-blocked runtime tests; no renderer-only control |
| TI-30 | Validate the release acceptance criterion on a release build | Network release gate | Offline local workflows plus allowed/denied/failed path evidence |

### 20.1 Existing evidence limitations

- T2.9.5B F1/F2 definitions were not independently re-read for this addendum.
- W4 formal closure remains an evidence gap even though W5–W7 supersede its practical blocker.
- No current runtime test was performed for this addendum.
- The exact UI planning document and expansion text are not present.
- The live prompt path split is source-confirmed but its runtime effect on RAG/memory is not measured.
- Model hot-swap has not been exercised with an active stream.
- Weight release has not been measured.
- Narration has no recovered authoritative semantics.
- The inference technique survey is a first-pass llama.cpp inventory; other engines and architecture families were not benchmarked.
- No optimization conflict combination or output-quality effect was measured.
- Automatic selection, resource estimation, and manual-override warnings are designs, not implementation evidence.
- The exact original document set for the recovered Author's Note requirements is not established.
- No network request was made, no packet capture was taken, and no offline, consent, audit, or remote-engine test was executed.
- The network inventory is a static source scan of inspected production files; test harnesses and generated/build files were not treated as product paths.

---

## 21. Explicit exclusions and protected boundaries

### 21.1 Explicit exclusions

| Excluded item | Scope of exclusion |
|---|---|
| Desktop character avatars | Persistent/reactive 3D desktop character surfaces, desktop interaction, multi-character avatar surfaces |
| Use of 3D asset generation to recreate that surface | The approved studio direction does not authorize it |
| F3-C amendment from this addendum | None is created |
| R6 | Not authorized |
| F3-D | Remains blocked |
| Implementation/migration/testing/build/model load | Not authorized by this mission |
| Cloud inference, cloud model providers, remote pairing | Not authorized; local-first remains. Remote ASR/TTS are not part of the local baseline and require separate owner authorization. |
| Automatic model lookup or background provider connection | Prohibited in the local baseline; a user-initiated, purpose-specific metadata permission is required before any replacement |
| Treating a renderer-only Local-Only Mode as backend enforcement | Prohibited as a compliance claim; see §23.4 |
| New TTS feature/contract change | F2-B7/TTS governance remains frozen |
| A frozen inference capability catalogue derived from the first pass | Research mission remains open until multi-engine/architecture evidence and benchmarks exist |
| Any optimization adopted merely because a backend accepts its parameter | No support/performance/quality claim without effective and measured evidence |
| Author's Note implementation or schema adoption | Requirement recovered; implementation unauthorized and no F3-C contract |
| Automatic selection as a safety bypass | Manual or automatic settings cannot bypass hot-swap readiness/release or data integrity |

Chat portrait thumbnails are **not** excluded. Character card images, world/scenario images, and generated images in chat are **not** desktop avatars.

### 21.2 Protected boundaries

| Boundary | Status |
|---|---|
| F3-C architecture freeze | **PRESERVED**; freeze ≠ acceptance/authorization |
| F3-C Phase 0 | **OPEN** |
| R6 | **UNAUTHORIZED** |
| F3-C acceptance | **NOT GRANTED** |
| F3-D | **BLOCKED** |
| D-ROOT / six defects | **OPEN** |
| D-STORAGE | **BLOCKED** |
| V20 | **BLOCKED pending F3-F** |
| A9/A10/A12/A13 | **CLOSED** |
| Gate v2 V1–V17 | **ADOPTED** |
| Gate v2 V18–V24 | **PROPOSED — NOT ADOPTED** |
| TTS / F2-B7 | **FROZEN** |
| REG-1 | **PROTECTED**; content hot-swap only |
| F2-B5 | **APPROVED** direction |
| E.6/E.8 | **PARKED** |
| MOB-1–MOB-4 | **BINDING CONDITIONAL CONSTRAINTS** |
| W4 | Historical/superseded in substance by W5–W7; formal closure gap |
| W8 | **VERIFIED** |
| Local-first | **REQUIRED; CURRENTLY CONTRADICTED BY UNSCOPED PATHS** — see §23.4 |

### 21.3 Artifact protection

No edit was made to:

- the master plan;
- the Muse review;
- the adjudication;
- the owner-decision packet;
- F3-C contracts, freezes, acceptance records, or protected artifacts;
- application source, tests, configuration, or data.

---

## 22. Exact corrections required in the master plan and owner-decision packet

### 22.1 Master-plan corrections

| ID | Section(s) | Required correction |
|---|---|---|
| MP-01 | S15, §4.7, §7.3, §8.2, §9 | Replace “model hot-swap is unaccepted” with “model hot-swap is a required product behavior, separate from F3-C and REG-1; implementation is unauthorized and unsafe today.” Keep stream/readiness/rollback gaps explicit. |
| MP-02 | S22, §4.3 T4.1/T4.2, §4.8, §9, §12 | Reclassify desktop character avatars as **EXCLUDED**. Retain `docs/tier4-avatar.md` only as historical evidence. Remove it from ordinary roadmap/deferred-future rows. Distinguish chat portraits. |
| MP-03 | S24, §4.8, §9, §12 | Replace “general creative studio unestablished” with “approved product direction; each modality separately scoped; no modality implemented.” Split F2-B5, E.6/E.8, and general studio. |
| MP-04 | S6, FND-09, §4.5, §8.3, §12 | Promote application-side memory from a deferred optional row to a required major subsystem, while retaining implementation-unauthorized status. Include typed inventory/clothing/relationship/validity/retrieval/reinjection. |
| MP-05 | §4.7 PROP-F3C-01/03/04; FND-07/FND-08; §8.2 | Update product classification: per-turn provenance required concept; actual seed preservation directed if useful; durable attachments approved. Do **not** adopt specific fields or an F3-C amendment. |
| MP-06 | S2, FND-09/FND-10, §4.1/§4.5, §8.3 | Add the UI requirements: bubble border/gap, four-property per-character palette (~10 choices each), portrait thumbnails, narration, hidden-by-default token counter, message-triggered image generation. |
| MP-07 | §4.1 T3.5, FND-09, FUT-CH3, §12 | State the exact context rule: flush after each character response in multi-character; never flush per turn in single-character. Preserve canonical history/lore/memory. |
| MP-08 | §4.7, §8.2, §9, §11.1 | Keep per-chat continuation-model preference explicitly unresolved and separate from active hot-swap and per-turn provenance. Add the FUT-CH1/CH2/C1 union to the schema prerequisite, not only a model pointer. |
| MP-09 | §4.4 FUT-C5, §13, S14/S15 | Add the pinned-runtime sampler/offload/context-shift inventory; distinguish stored, sent, accepted, and effective controls. Do not claim unsupported quadratic/TFS/Top-A features. |
| MP-10 | §4.4 FUT-PROMPT1, S14 | Replace generic “expanded families” with a required evidence matrix, confidence/fallback model, and hot-swap template transition. |
| MP-11 | S23, §4.8, §8.2, FND-15, §9.4 | Preserve MOB-1–MOB-4; keep remote implementation unstarted/unauthorized; do not design a remote boundary. |
| MP-12 | §4.7/§4.8, §8.2, §9.5 | Add the studio shared-foundation versus modality-specific split; preserve local-first; no assumption of one backend. |
| MP-13 | §4.3 T4, §9.4/§9.5 | Clarify that 3D mesh/texture generation is an approved studio direction, but using it to create the excluded desktop avatar is not authorized. |
| MP-14 | §2.1, §4.7, §12 | Record the separate UI planning document and expansion text as unrecovered; do not imply exact source recovery. |
| MP-15 | §8.2, §12 | Carry W4 as historical/superseded by W5–W7 with formal-closure gap; W8 as VERIFIED. Do not reintroduce stale blocker/open labels. |
| MP-16 | §9.1/§9.5, §17 | Add per-category gates and the rule that reserving a field/identity/boundary creates no migration, backfill, compatibility promise, or release change. |
| MP-17 | S14/S15, §4.4, §8.2, §9 | Replace the Flash-Attention/offload footnote with a first-class inference capability mission: technique inventory, architecture mapping, backend/version support, conflicts, automatic selection, manual override, and benchmarks. |
| MP-18 | §4.1 T3.1/T3.5, §8.2, §12 | Separate conversation context rollover from multi-character context flushing in definitions, dependencies, UI wording, tests, and failure behavior. |
| MP-19 | §4.1 FUT-CH1, §7.4, §8.2 | Add Author's Note flyout/floating/hidden/side presentation, chat-only scope, existing-chat persistence, insertion order, enable state, and blank-on-new-chat behavior. |
| MP-20 | §1.5, §2.1, §4.7, §12 | Record that the two F2 reports and ledger are evidence leads, not proof that all original requirements survived; the exact UI/expansion artifacts and source-document count remain unestablished. |
| MP-21 | S14/S15, §9.1, §12 | State that the first-pass inventory is not a frozen capability catalogue; no optimization is adopted or benchmark-successful by inclusion. |
| MP-22 | S14/S15, §9.1, §12 | Add the owner's non-mandate boundary and the ordered research/evidence/compatibility/selection/validation/authorization stages. |
| MP-23 | S2, S23, §4.8, §8.2, §9.5, §12 | Add **Network Boundaries, External Integrations, and User Consent** as a first-class product and architecture mission, not a privacy sub-setting. Preserve the three-category policy and the no-silent-communication rule. |
| MP-24 | §4.8, §8.2, §9.4, §9.5, §12 | Add the backend-enforced offline boundary, destination registry, six-question consent, purpose-scoped permissions, local audit, cache/retention controls, failure visibility, and no-permission-creep rule. Do not describe the renderer-only toggle as complete. |
| MP-25 | §4.4, §4.8, §9.1, §12 | Record the complete existing/proposed network-path inventory in §23, classify remote ASR/TTS and automatic Civitai lookup as outside the local baseline unless separately authorized, and add V-NET-01–V-NET-28 plus the release acceptance criterion. |
| MP-26 | §2.1, §4.1, §9.5, §12 | Correct the claim that Local-Only Mode protects the application. It currently wraps renderer `fetch`/XHR only; `disableTelemetry` only sets a renderer flag. Backend/sidecar enforcement is missing. |
| MP-27 | Entire corrected plan | Require bidirectional traceability. Every corrected-plan entry and work package names the reconciliation ID(s) it serves; every reconciliation row names its planned disposition, work package, or explicit unresolved status. No orphan work package and no silently dropped requirement. The original 66 rows plus NET-01–NET-20 must survive the plan. |
| MP-28 | Entire corrected plan | Do not assign implementation estimates until the architecture and existing-capability investigations establish what is present. Dependency and sequencing analysis may proceed; effort, duration, and staffing estimates must be marked provisional and evidence-dependent, not presented as commitments. |

### 22.2 Owner-decision packet corrections

| ID | Packet section | Required correction |
|---|---|---|
| OP-01 | OD-01 | Split the old four-part question. Mark provenance as required concept with record recommendation; seed preservation as directed if meaningful; attachments approved; per-chat model preference still unresolved. Do not present all four as one open owner choice. |
| OP-02 | OD-02 model hot-swap | Remove “deferred product decision.” It is now required. Preserve “no F3-C amendment” and list lifecycle evidence as technical. |
| OP-03 | OD-04 general studio | Remove “unestablished/owner must define whether.” Record approved direction, with modality scopes unresolved. Keep F2-B5/E.6/E.8 separate. |
| OP-04 | OD-05 desktop-avatar interaction | Change to **EXCLUDED**; do not offer retain/defer options. |
| OP-05 | OD-06 Story Tools | Keep open/unresolved; no change from evidence alone. |
| OP-06 | OD-03 Android/remote | Keep unresolved; MOB constraints preserved, no remote work. |
| OP-07 | T-01–T-09 | Add hot-swap lifecycle, actual-seed availability, media lifecycle, memory retrieval, model-family matrix, and per-character context test as technical investigations. |
| OP-08 | COR-01–COR-10 | Add corrections for UI omission, narration unrecovered, token counter, message-triggered image generation, memory promotion, studio approval, avatar exclusion, and model hot-swap product status. |
| OP-09 | Final register | Add: owner decisions now resolved; remaining UD-01–UD-06; technical TI-01–TI-18; corrections that proceed without owner input; explicit exclusions. |
| OP-10 | Artifact identity | If the owner packet is later edited, its hash changes. Do not continue citing the current `6721d7…` hash as if the edited file still matches. This addendum does not edit it. |
| OP-11 | Architecture scope | Add the first-class hardware/architecture-aware inference mission; do not reduce it to Flash Attention, GPU layers, or a fixed llama.cpp menu. |
| OP-12 | Source recovery | Add the F2 engineering reports as leads and state explicitly that they do not establish complete survival of the missing UI/expansion documents. |
| OP-13 | Capability mission | Preserve the distinction between possessing the capability system and enabling every technique/backend feature; keep research, selection, validation, and authorization separate. |
| OP-14 | Network mission | Add the first-class network boundary, the complete path inventory, the existing frontend-only/offline contradiction, and the release criterion. Do not treat remote ASR/TTS, automatic Civitai lookup, or Android access as authorized by this addendum. |
| OP-15 | Corrected-plan adoption | Require every plan entry/work package to reference its reconciliation ID(s), preserve all 66 original rows plus NET-01–NET-20, and label every effort estimate provisional until the investigations establish the implementation baseline. |

### 22.3 Master-plan adoption

None of these corrections adopts the corrected master plan. The plan must be revised, reviewed against the source hashes and this recovery baseline, and separately authorized by the owner before adoption.

The revised plan must preserve the reconciliation as a two-way map: start from any requirement and find its planned disposition, or start from any work package and find the requirement that justifies it. A plan entry without a reconciliation identifier is an orphan. A reconciliation row without a disposition, work package, or explicit unresolved status is a silent drop. The plan is not complete until it consumes the original 66 requirement rows, the additive NET-01–NET-20 family, all six reconciliation outcomes, the seven-state capability register, owner decisions/exclusions, unresolved owner decisions, and the network negative tests.

Effort estimates are not part of this recovery baseline. Dependency and sequencing analysis may begin, but implementation estimates remain provisional until the architecture and existing-capability investigations establish what is actually present. An estimate before that evidence is a planning assumption, not a commitment.

---

## 23. Network Boundaries, External Integrations, and User Consent

### 23.1 Status and controlling boundary

This is a **first-class product and architecture mission**, not a privacy sub-setting and not an implementation detail. It governs the desktop application, every sidecar it starts, every model and dependency workflow, remote access, logging, error handling, updates, and any future integration. It is a release criterion in §23.10.

The controlling product statement is:

> Howl Chat is a local application that may occasionally use the internet with the user's permission. It is not a cloud-connected AI service.

The policy has three categories:

| Category | Policy |
|---|---|
| Local inference and creative work | Offline. No cloud inference or cloud model providers. |
| Android mobile remote | Permitted network functionality, explicitly configured and secured. |
| External model information and downloads | User-initiated, purpose-specific permission required. |

The following are prohibited in the local baseline:

- silent telemetry;
- background model lookups;
- automatic provider connections;
- unexpected update checks;
- a new destination or purpose added without a documented product decision and an appropriate permission path.

Existing optional remote engines do not become permitted merely because they are optional. Remote ASR through DashScope and a remotely configured MeloTTS endpoint are cloud model services. They are currently exposed through the audio surface without the authorization and disclosure this mission requires. The source evidence is recorded in §23.3; the contradiction is recorded in §23.4.

This section distinguishes four different questions:

| Question | What it establishes |
|---|---|
| Does a network path exist? | Static source evidence only |
| Is it reachable in the local baseline? | Runtime behavior, not established here |
| Is it disclosed and authorized? | Consent and governance review |
| Is it safe and worthwhile? | Validation and owner authorization, neither granted here |

Source presence is not validation. A permission control is not proof that no other path exists. A global offline toggle is not effective if background code can bypass it.

### 23.2 Network requirements

| ID | Requirement | Current outcome | Evidence and limit |
|---|---|---|---|
| NET-01 | Local creative work and local inference remain fully usable with external connectivity disabled. | **Unresolved** | Local backends exist, but no offline end-to-end test was run; Civitai and remote engine paths complicate the claim. |
| NET-02 | No cloud inference or cloud model provider is part of the local baseline. | **Contradicted** | Remote DashScope ASR and remotely configurable MeloTTS are exposed without separate authorization. |
| NET-03 | No background model lookup, telemetry, provider connection, or update check occurs without a user action and permission. | **Contradicted** | Image scans can call Civitai automatically from `scanFile`; see N-01. |
| NET-04 | Every external destination is recorded in a destination registry with purpose, feature scope, data classes, and authorization mode. | **Missing** | No registry exists. |
| NET-05 | Every connection prompt answers the six questions in §23.6 before any request is sent. | **Partially satisfied** | The image-dependency modal gives one generic sentence; the error reporter has a toggle and preview but no full disclosure; audio has labels, not consent. |
| NET-06 | Opening a website, downloading a file, and allowing Howl to organize or track that file are separate actions and permissions. | **Partially satisfied** | Image dependencies distinguish a manual external link from an automatic download, but not with a full disclosure; no separate indexing permission exists. |
| NET-07 | Permission is purpose-specific; metadata retrieval does not imply download, indexing, background refresh, or another provider. | **Missing** | No scoped permission model exists. |
| NET-08 | A global offline control is enforced in the backend and in every sidecar, not only in the renderer. | **Contradicted** | `localOnlyMode` wraps frontend `fetch`/XHR only; Go code and Wails methods are unaffected. |
| NET-09 | A human-readable local connection audit records allowed and denied attempts. | **Missing** | No audit history was found. |
| NET-10 | Cache and retention behavior is disclosed; cached metadata is identifiable, clearable, and does not silently refresh. | **Partially satisfied** | Civitai names are cached in manifests; no user-facing retention or refresh control exists. |
| NET-11 | A failed lookup or interrupted download is visible; Howl never silently substitutes unverified information. | **Partially satisfied** | Dependency downloads report errors; Civitai lookup errors are silently ignored and the scan proceeds. |
| NET-12 | Android remote access is separately configured, authenticated, encrypted, revocable, and never enabled by default. | **Missing** | Remote access is unimplemented and unauthorized; MOB constraints remain binding. |
| NET-13 | External model information and downloads are user-initiated and purpose-specific. | **Partially satisfied** | Dependency downloads are user-initiated; the automatic Civitai lookup is not. |
| NET-14 | Adding a destination, purpose, or data class requires a documented product decision; permission does not creep. | **Missing** | No change-control or registry gate exists. |
| NET-15 | The release acceptance criterion in §23.10 passes on the release build. | **Missing** | No network acceptance test or runtime evidence exists. |
| NET-16 | Remote destinations use authenticated TLS; loopback may use plain HTTP; custom URLs are validated and cannot silently downgrade. | **Partially satisfied** | MeloTTS accepts any parsed `http(s)` URL; DashScope uses HTTPS; custom endpoint validation is incomplete. |
| NET-17 | Credentials and API keys are not logged, included in reports, or sent to any destination other than their documented service. | **Partially satisfied** | Keys are environment or config inputs; no secret-handling test or redaction proof was found. |
| NET-18 | Enabling offline mode blocks external egress without breaking local chat, image generation, local ASR/TTS, scanning, or settings. | **Missing** | Only renderer `fetch`/XHR is wrapped; no backend enforcement or runtime test exists. |
| NET-19 | Every proposed path has a destination, purpose, disclosure, cache rule, failure rule, and authorization decision before implementation. | **Missing** | Proposed paths are listed in §23.5 but not authorized. |
| NET-20 | No automatic update check exists; any future check is explicit, signed, and separately permissioned. | **Satisfied** (absence only) | No application update checker was found; this must remain true. |

### 23.3 Existing external network-path register

Scope: production and shipping code inspected for this addendum. Test servers, the native test seam, and development harnesses are excluded in §23.4. This is a complete register **against the inspected source**, not a packet-capture proof and not a statement that a path is reachable on a given machine.

| ID | Path and destination | Trigger | Information sent | Information received | Current consent and control | Cache, retention, failure | Outcome |
|---|---|---|---|---|---|---|---|
| N-01 | `GET https://civitai.com/api/v1/model-versions/by-hash/{sha256}` | Image scan where the model name is generic and no cached Civitai name exists (`internal/image/scanner/scanner.go`, `civitai.go`) | Model file SHA-256, ordinary HTTP metadata, IP address, proxy environment if configured | Model metadata used to replace a fallback display name | **None.** No prompt, setting, or audit. The caller ignores lookup errors. | Civitai name and hash cached in the model manifest; no refresh or clear control; failure is silent | **Non-compliant** |
| N-02 | `GET https://huggingface.co/.../resolve/main/...` for image companion files | User accepts the missing-dependency modal while loading a split model | File URL path, ordinary HTTP metadata, IP address | Model weight files | One generic consent sentence plus a total-size and save-directory note; no destination, data-received, recurrence, or retention detail | Files saved in the model directory; `.part` file removed on failure; failure shown with manual links; no checksum proof | **Partially compliant** |
| N-03 | `GET https://api.github.com/repos/ggml-org/whisper.cpp/releases/latest`; then the selected release asset | Explicit whisper installer/acquisition path (`downloader_acquire.go`) | Ordinary HTTP metadata, IP address, GitHub API request | Release metadata and a Windows ZIP or equivalent asset | No pre-request destination/purpose disclosure was found; the install action is user-invoked but not a six-question consent flow | Files installed under `bin/whisper` or `HOWL_WHISPER_HOME`; download errors are returned; no audit | **Partially compliant** |
| N-04 | `git clone https://github.com/ggml-org/whisper.cpp.git` | Unix whisper source-build acquisition path | Ordinary Git/HTTP metadata, IP address | Source repository | No separate disclosure or permission; user-invoked installer only | Source tree retained for the build; failure returned; no audit | **Partially compliant** |
| N-05 | `POST https://howl-error-reporter.vercel.app/api/report` | User enables error reporting and clicks **Send Report** (`app-errors.go`, `internal/errors/reporter.go`) | App version, OS, architecture, Go version, CPU count, up to 50 error entries with code, severity, operation, location, time, and an arbitrary `Context` map; IP address and HTTP metadata | HTTP response only; the response body is not used | Default disabled. Toggle and preview exist. The UI claims “error codes and line numbers” but the `Context` contents are not proven to match that claim. | Upload is asynchronous; network errors are silently discarded; no audit or server-retention disclosure | **Partially compliant** |
| N-06 | `POST https://dashscope-intl.aliyuncs.com/api/v1/services/audio/asr/transcription` | Selecting a Qwen3 ASR remote engine and transcribing (`internal/audio/asr/qwen/client.go`) | API key, model identifier, local audio filesystem path submitted as `file_url`, request options; IP address and HTTP metadata | Task identifier and status | No consent or disclosure beyond the engine label “DashScope API, remote”; no audit | Task polling continues until completion or timeout; no local audio bytes are uploaded by this implementation, so remote transcription may fail independently of privacy | **Non-compliant** |
| N-07 | `GET https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id}` | Polling a Qwen ASR task | Task identifier, API key, IP address and HTTP metadata | Task status and transcription result | Same as N-06 | Polling is repeated; cancellation and failure are not surfaced as a network-consent event | **Non-compliant** |
| N-08 | `GET` the result URL returned by DashScope | Fetching a completed remote ASR result (`client.go` result download path) | Provider-issued URL, API headers where applicable, IP address | Transcript or result payload | Same as N-06; the final host may be a provider CDN rather than the API host | Result retrieval failure is returned; no destination registry entry for the returned host | **Non-compliant** |
| N-09 | `POST` to the configured MeloTTS base URL, first `/synthesize`, then `/v1/audio/speech` on 404/405 | Selecting MeloTTS and synthesizing speech (`internal/audio/tts/melotts/client.go`) | Text to speak, language, voice, speed, optional bearer token, IP address and HTTP metadata | Audio bytes | Default `http://127.0.0.1:8888` is local. A remote URL can be supplied through `HOWL_MELOTTS_URL` or `TTSConfig.ModelPath`; the UI says “Remote API” but there is no per-use consent or remote-destination disclosure. | Audio may be cached by the TTS layer; HTTP errors are returned; a remote endpoint is not distinguished from local in the consent model | **Non-compliant when remote; compliant only while loopback** |
| N-10 | TCP dial to the configured MeloTTS host and port | Capability/probe refresh (`internal/audio/capability/runtime.go`) | A TCP connection attempt, DNS lookup if a hostname is used, source IP | TCP acceptance or refusal | No consent; a remote configured URL is dialed during status/probe paths | Result cached for about five seconds; failure shown as unreachable; no audit | **Non-compliant when remote** |
| N-11 | System browser navigation through `BrowserOpenURL` to a dependency `Source` URL | User clicks **Get it here** in the missing-dependency modal (`frontend/dist/js/image-deps.js`) | The browser sends its normal request metadata, IP address, cookies for that site, and any site-specific state; Howl does not add app data | The website's response in the system browser | The click is separate from download, but there is no pre-open warning that the site will see the visit and connection metadata | Browser history and site cookies are outside Howl's control; no Howl audit | **Partially compliant** |
| N-12 | Tutorial links rendered as `<a target="_blank">` to GitHub, Hugging Face, Civitai, and other publishers | User clicks a link in tutorial content (`frontend/dist/js/tutorials-ui.js`) | Normal browser/webview request metadata and IP address | The linked page | No warning or destination registry; links are user-initiated only | Depends on webview/browser behavior; no Howl audit | **Partially compliant** |
| N-13 | Optional proxy route through `HTTP_PROXY`/`HTTPS_PROXY` environment variables | Any Go HTTP client using the default transport, and explicitly the Civitai client via `http.ProxyFromEnvironment` | All of the above requests may traverse a user-configured proxy, which can observe destinations and metadata | Proxied response | No visibility in the app; the proxy is outside Howl's disclosure | Proxy behavior and retention are unknown; offline enforcement must account for it | **Unresolved** |
| N-14 | External source URLs stored in model profiles, ASR/TTS catalogues, and image metadata | Display or copy only; no request found in the inspected production path | Nothing until the user opens the URL through a browser | Nothing | None needed while display-only | Not a network path by itself; must not be counted as one | **Not a path** |

No application analytics, crash-upload SDK, cloud sync, automatic updater, or background model-download scheduler was found in the inspected production source. The only remote inference clients found are N-06–N-10: DashScope ASR and a remotely configured MeloTTS endpoint. Absence in this inspection is not a permanent guarantee; it is a requirement to keep the registry and static egress check complete.

### 23.4 Current contradictions and local/loopback paths

#### 23.4.1 Current contradictions

These are source-confirmed policy failures, not merely missing UI polish. Each must be resolved before the local-first claim can be treated as earned.

| Contradiction | Current evidence | Required correction |
|---|---|---|
| The product claims local-first, but an image scan can contact Civitai | `scanFile` calls `IdentifyFromFile` for a generic name and no cached name; `civitai.go` sends the model SHA-256 to `civitai.com` without a prompt or setting | Disable the automatic lookup by default, require a purpose-specific metadata permission, or remove it. A scan must never be the consent trigger for an undisclosed provider. |
| Local-Only Mode is presented as an offline protection, but it only wraps the renderer | `frontend/dist/js/settings-actions.js` replaces `window.fetch` and `XMLHttpRequest.open`; Go `net/http`, `net.Dial`, sidecars, Wails calls, and browser opens are outside it | Enforce offline mode in the backend transport/dialer and sidecar policy, then prove it with an egress-blocked test. |
| Disable-telemetry is presented as a privacy control, but it only sets a renderer flag | `window.__telemetryDisabled` is assigned; no inspected Go consumer reads it | Remove the misleading control or make it govern every telemetry/error upload path in the backend. |
| Remote model services are exposed inside the local audio surface | DashScope ASR is selectable as a remote engine; MeloTTS accepts a remote URL through `HOWL_MELOTTS_URL` or `TTSConfig.ModelPath`; neither has an operation-specific consent flow | Treat remote ASR/TTS as external integrations requiring separate owner authorization and disclosure, or remove them from the local baseline. |
| The image-dependency consent is too broad to be informed consent | The modal shows one sentence plus size/directory; it does not name the destination, data received, recurrence, or retention | Replace it with the six-question, purpose-scoped prompt and a distinct download permission. |
| There is no authoritative destination registry or connection audit | No registry, allowlist, connection history, or denied-attempt record was found | Add the registry, local audit, and no-permission-creep gate described in §23.6. |
| The error-report UI claims no model names or personal data, but the payload includes an arbitrary context map | `Report` includes `ErrorPayload.Context`; no redaction proof was found | Verify and constrain the payload before repeating the claim; show the exact preview and destination. |
| Manual source links and downloads are adjacent enough to be confused | The dependency modal opens a publisher page with `BrowserOpenURL` and downloads files through a separate path, but does not fully explain the browser's own network metadata | Keep browser open, download, and local indexing as three separately recorded actions with separate disclosures. |

#### 23.4.2 Local and loopback paths

These are not external integrations, but they are part of the network-boundary inventory because a loopback-only promise must be enforced and verified.

| ID | Local path | Evidence | Boundary requirement |
|---|---|---|---|
| L-01 | llama.cpp chat server at `127.0.0.1:8080` | `app-llama.go`; local HTTP client | Must bind loopback only; no remote host setting may silently turn it into a LAN service. |
| L-02 | llama.cpp ASR server at an ephemeral `127.0.0.1` port | `internal/audio/asr/llama/llama.go` | Ephemeral loopback bind; lifecycle and port ownership must be verified. |
| L-03 | sd-server image backend at `http://localhost:<port>` | `internal/image/engine/sdcpp/engine.go` | Local process only; port selection and process ownership must be verified. |
| L-04 | MeloTTS default `http://127.0.0.1:8888` | `internal/audio/tts/melotts/client.go` | Local default; any non-loopback URL becomes N-09 and requires the external policy. |
| L-05 | Wails/frontend asset origin | Local app runtime | Local UI loading; not a cloud service. |

Excluded from the product register:

- `app_native_test_seam.go` and other test-only seams;
- `harness_main.go` and test harness servers;
- `httptest` servers in test files;
- `tools/*.js` development utilities;
- local fixture servers and build probes.

These may bind ports or accept test traffic, but they are not shipping product network policy and must not be used to justify a production path.

### 23.5 Proposed network paths

These paths are **not implemented or authorized**. They are the paths implied by the current product direction, the owner's examples, and the corrected-plan requirements. Each must pass the registry, consent, security, failure, and verification requirements before any code is written.

| ID | Proposed path | Purpose | Required authorization and disclosure |
|---|---|---|---|
| P-01 | Model metadata/compatibility lookup to a named publisher such as Hugging Face or Civitai | Retrieve model metadata and compatibility information for a user-selected model | User-initiated; show publisher, model identifier or hash sent, metadata received, whether the result is cached, and whether the lookup can recur. A metadata permission must not permit a download. |
| P-02 | General model download manager from a named publisher | Fetch a user-selected model or dependency | Separate download permission; size, destination, source, checksum or signature status, retry policy, cancellation, and no silent mirror substitution. |
| P-03 | Application update check | Discover a newer Howl release | Explicit permission or a manual check; no startup check; disclose endpoint, version data sent, update metadata received, signature verification, and whether the check recurs. |
| P-04 | Remote ASR provider | Transcribe audio outside the device | Requires an explicit owner decision separate from this baseline. Disclose provider, audio or path data sent, transcript received, retention, credential use, and why local ASR is insufficient. It is not part of the local-inference promise. |
| P-05 | Remote TTS provider | Synthesize speech outside the device | Same separate authorization as P-04; disclose text, voice, speed, provider, audio received, retention, and whether generated audio is cached. |
| P-06 | Android remote access | Control or use Howl from a configured phone | Pairing, authentication, transport encryption, bind scope, revocation, session timeout, and no cloud relay by default. MOB-1–MOB-4 remain binding. |
| P-07 | External publisher/source link | Open a model, dependency, or documentation page in the browser | Tell the user that the site will see the visit and ordinary connection metadata; keep this distinct from downloading and from Howl indexing the result. |
| P-08 | Provenance/licensing metadata lookup | Retrieve license or provenance data for a model | Same metadata rules as P-01; no upload of model weights or user content; cache and attribution behavior disclosed. |

Explicitly **not** proposed: cloud chat inference, cloud image generation, cloud memory, cloud conversation storage, silent telemetry, or a background model-indexing service. Those conflict with NET-02 and require a product-boundary change, not a settings toggle.

### 23.6 Destination registry, permission, cache, and audit model

#### 23.6.1 Required destination registry

Every external destination must have a record with at least:

| Field | Meaning |
|---|---|
| Destination ID | Stable internal identifier |
| Scheme, host, port, and path scope | Exact allowed destination, not a wildcard without justification |
| Purpose | The user-visible reason for the connection |
| Feature scope | The feature and action allowed to use it |
| Direction and data classes | What is sent and received; whether prompts, audio, images, paths, identifiers, IP, or credentials are involved |
| Authorization mode | One-time, session, persistent, or prohibited |
| Credential source | None, environment, OS secret store, or user-entered; never logged |
| Cache and retention | What is stored locally, where, for how long, and how to clear it |
| Failure behavior | User-visible result and whether local work continues |
| Audit behavior | What is recorded locally, with sensitive values excluded |
| Owner/decision record | Who approved the destination and purpose |
| Verification status | Static, mocked, runtime, or release-verified |

A destination not in the registry must be blocked, not merely logged. The registry is a product and architecture artifact; this section does not create it in code.

#### 23.6.2 Six-question consent

Every prompt must answer, in plain language:

1. **Who** is Howl connecting to?
2. **Why** is Howl connecting?
3. **What** information is sent?
4. **What** information is received?
5. **What** will Howl do with the result?
6. **Can this happen again** without asking?

The prompt must name the publisher or service, the action, the data classes, the result, the local storage/use, and the permission lifetime. It must not require the user to understand HTTP, headers, metadata schemas, or TLS. It must include the ordinary connection-metadata disclosure: the destination will see the user's IP address and normal connection information, just as a website would.

For a browser-open action, the disclosure must say that Howl is handing the URL to the browser and cannot promise that the site receives nothing. For a download, it must state the file, approximate size, destination, and validation. For a metadata lookup, it must state the model identifier or hash and whether the result is cached. These are different permissions.

#### 23.6.3 Permission states and cache

Permissions must be scoped to a destination, purpose, and action. The minimum states are:

- **Denied** — no request; local operation continues.
- **Allow once** — one request or one bounded operation.
- **Allow for this model/session** — bounded to the named model or current session.
- **Allow and remember** — persistent for that exact destination and purpose, revocable at any time.
- **Prohibited** — a product-level block that no override can bypass.

A persistent metadata permission must not authorize a download, a second publisher, a background refresh, or indexing. Cache entries must record source, model identity, timestamp, and permission scope. Clearing a cache entry must make the next lookup ask again unless a current persistent permission explicitly covers it. A background refresh must not be hidden inside a cache miss.

#### 23.6.4 Audit history

Howl must maintain a local, human-readable connection audit. At minimum it should record:

- timestamp;
- feature and action;
- destination ID and human-readable host;
- purpose;
- permission decision (allowed, denied, expired, revoked);
- data classes sent and received;
- approximate bytes or file size where applicable;
- result (success, failure, cancelled, offline);
- whether a cache was used.

The audit must not store prompts, conversation text, audio content, model contents, API keys, or full local paths by default. It must be viewable, exportable, and clearable locally. Denied attempts should be recorded so permission creep is observable. The audit itself must not be uploaded.

### 23.7 Global offline enforcement

Offline enforcement must be a backend property, not a renderer trick.

The current frontend `applyPrivacySettings` wrapper intercepts `fetch` and `XMLHttpRequest` in the webview. It does not intercept Wails method calls, Go `net/http`, `net.Dial`, sidecar processes, `BrowserOpenURL`, or code that runs before the renderer loads. The separate `disableTelemetry` flag only sets `window.__telemetryDisabled`; no backend consumer of that flag was found. These controls therefore do not establish the promised boundary.

A compliant design needs at least:

- one policy-enforcing HTTP transport and dialer used by all production clients;
- explicit loopback exceptions for L-01–L-05;
- a deny-by-default destination registry;
- no raw `http.DefaultClient`, `net.Dial`, or equivalent bypass in production code;
- a sidecar policy so a bundled engine cannot contact the internet on Howl's behalf;
- a renderer guard for links and webview navigation;
- a visible offline indicator and a clear statement of what remains available;
- immediate cancellation or refusal of in-flight requests when offline mode is enabled;
- test coverage proving the backend, not just the UI, is blocked.

Offline mode must not disable local inference, local image generation, local ASR/TTS, local scanning, local metadata display, or local file management. A remote-only engine may be unavailable; the UI must say so rather than silently substituting another provider.

### 23.8 Failure, security, and retention rules

- **TLS:** remote destinations require HTTPS with normal certificate validation. Loopback HTTP is allowed only for local services. A custom remote URL must not silently downgrade to HTTP.
- **Credentials:** API keys are secrets. They must not appear in logs, error reports, audit entries, URLs, or model metadata. The current environment-variable pattern is acceptable only with redaction and documentation; a future secret store is preferable.
- **No silent substitution:** a failed metadata lookup, download, ASR request, or TTS request must not be replaced by an unverified result from another host or a stale cache presented as current. The user must be told what failed and what local capability remains.
- **Download validation:** downloads must use a temporary `.part` file, verify length, and verify a published checksum or signature where available. No automatic mirror fallback. A failed or truncated file must not be used.
- **Cancellation:** every user-initiated network operation must be cancellable. Cancellation must be visible and must leave local data usable.
- **Retention:** download files remain until the user deletes them; metadata caches have a documented lifetime and clear action; audits are local and clearable; error-report server retention must be disclosed before upload.
- **Browser opens:** the app must not claim that opening a browser transmits nothing. The site receives the normal visit metadata, IP address, and any site cookies. This is a separate user action from a download or an indexing decision.
- **Permission creep:** adding a new host, path scope, data class, or purpose requires a documented product decision and a registry update. A generic “allow internet” switch is not sufficient.

### 23.9 Verification criteria

These criteria are release gates, not suggestions. Each must have repeatable evidence. Until then, the network mission is not satisfied.

| ID | Verification | Pass condition |
|---|---|---|
| V-NET-01 | Static egress inventory | Every production `net/http`, `net.Dial`, HTTP client, `BrowserOpenURL`, `window.open`, `fetch`, XHR, WebSocket, updater, and sidecar launch is enumerated and mapped to a registry entry or an explicit local/test exclusion. |
| V-NET-02 | Startup and idle capture | With external egress blocked, app startup, settings, model scan, model list, and idle operation produce no outbound connection attempt. |
| V-NET-03 | Offline local workflow | With external egress blocked, local chat, model load, image generation, local ASR, local TTS, scanning, and settings complete without external access. |
| V-NET-04 | Backend enforcement | Enabling offline mode blocks a deliberately injected Go HTTP/dial path and a Wails-triggered network action, not just renderer `fetch`. |
| V-NET-05 | Six-question prompt | For every external path, the UI answers all six questions with the actual destination, data classes, result use, and recurrence. |
| V-NET-06 | Decline path | Declining any prompt sends no request; the feature fails visibly and local work continues. |
| V-NET-07 | Permission scope | Metadata-only permission cannot download; download permission cannot index; one publisher cannot authorize another; revoke takes effect immediately. |
| V-NET-08 | Action separation | Opening a link, downloading a file, and indexing/organizing it are three separately recorded actions with distinct consent states. |
| V-NET-09 | Civitai scan | A generic-name scan with no metadata permission performs no network request. With permission, the request contains only the documented hash/metadata and the result is cached and auditable. A failed lookup shows “metadata unavailable,” not a silent name. |
| V-NET-10 | Dependency download | Cancellation, offline failure, truncated body, HTTP error, and checksum mismatch are visible; no partial file is used; manual links remain available; no mirror is substituted. |
| V-NET-11 | Remote ASR | No remote ASR request occurs without a specific prompt naming DashScope, the data sent, the result, retention, and recurrence. If remote ASR is not authorized, the option is not offered as part of local inference. |
| V-NET-12 | Remote TTS | A non-loopback MeloTTS URL is treated as an external destination with a prompt; plaintext remote HTTP is rejected; credentials are redacted; local Piper/Kokoro remain usable. |
| V-NET-13 | Audit | Allowed, denied, cancelled, and failed attempts appear in a local human-readable audit with no prompts, audio, model contents, keys, or full paths. |
| V-NET-14 | Cache/retention | Cache clear causes the next lookup to re-ask under the current permission; a stale cache is not presented as a fresh result; download retention is visible. |
| V-NET-15 | Failure behavior | Every external failure is visible, has a local fallback or a clear explanation, and never silently substitutes another provider or unverified data. |
| V-NET-16 | Android remote | Remote access is off by default, requires explicit pairing/auth, uses encrypted transport, binds only the configured scope, expires sessions, and can be revoked. |
| V-NET-17 | Error reporting | Default is off; preview matches the actual payload; the user can see destination, data, and retention; upload failure is visible; no automatic send occurs. |
| V-NET-18 | Browser links | The UI warns that the site will see the visit/IP; the action is recorded separately from download/indexing; offline mode blocks or clearly warns about external navigation. |
| V-NET-19 | Update behavior | No startup or periodic update request occurs; any manual check names the endpoint, sent data, received data, signature check, and recurrence. |
| V-NET-20 | Registry completeness | A production build fails its network acceptance check if any external host or dial path lacks a registry entry, an owner decision, or a documented test-only exclusion. |
| V-NET-21 | Offline egress negative test | With external connectivity disabled, no external destination receives traffic. Verify at the application boundary with network capture, firewall/proxy evidence, and backend instrumentation, not only a renderer mock. |
| V-NET-22 | Unauthorized scan lookup | A generic-name scan with no metadata permission performs no outbound request; a denied prompt performs no request and leaves a visible “metadata unavailable” result. |
| V-NET-23 | Local model load | Loading a local model without external permission produces no provider connection, DNS lookup, browser handoff, or model-download request. |
| V-NET-24 | Permission fail-closed | A failed, expired, revoked, or malformed permission cannot fall through to another route, mirror, provider, cached remote result, or renderer path. |
| V-NET-25 | Browser/download separation | A browser handoff is recorded as a browser action, not an application download; it writes no file and does not index the resulting page or file. |
| V-NET-26 | Revocation | Revoking a permission prevents subsequent requests for that destination/purpose, leaves local workflows usable, and requires a new decision before re-allow. |
| V-NET-27 | Android isolation | Android remote remains separately disabled/enabled; revoking it does not change desktop policy, and desktop offline mode cannot be bypassed through remote access. |
| V-NET-28 | Boundary-level test harness | All negative tests run at the application boundary — Go transport/dialer, sidecar, and browser handoff — not only through renderer mocks. The release check fails if a production bypass exists. |

A mocked local HTTP server may be used to exercise the success, failure, timeout, cancellation, redirect, and oversized-response branches. A passing mock does not prove the real publisher’s behavior; it proves Howl’s policy code.

### 23.10 Release acceptance criterion

The release cannot be accepted unless all of the following are true:

> A user must be able to operate Howl Chat's local creative and inference capabilities without an internet connection. Every external connection must have a documented purpose, destination, information disclosure, and user authorization appropriate to that action. No undisclosed external communication is permitted.

The criterion is testable:

1. **Offline operation** — local capabilities pass V-NET-03 with egress blocked.
2. **Complete disclosure** — every external path passes V-NET-05 and V-NET-06.
3. **Purpose-specific authorization** — every path passes V-NET-07 and V-NET-08.
4. **No silent communication** — startup, idle, scan, load, and local generation pass V-NET-02.
5. **Visible failure** — every path passes V-NET-15.
6. **Auditability** — every path passes V-NET-13.
7. **Registry completeness** — V-NET-01 and V-NET-20 pass on the release build.
8. **Negative-path enforcement** — V-NET-21–V-NET-28 pass at the application boundary, not only in renderer tests.

A feature is not compliant merely because a prompt exists. A prompt that permits a broader action than it describes is a permission failure. A toggle that only blocks the renderer is not offline enforcement. A network path that exists but is “optional” is still part of the disclosure and authorization boundary.

### 23.11 Implementation status and non-authorization

This section records requirements and a source inventory. It does not:

- implement a registry, consent UI, audit, transport gate, offline mode, or remote-access system;
- remove or authorize Civitai, DashScope, MeloTTS, GitHub, Hugging Face, or any other destination;
- authorize remote ASR/TTS, Android access, updates, telemetry, or model downloads;
- establish that any existing path is safe, reachable, or validated;
- modify source, tests, builds, migrations, configuration, contracts, freezes, acceptance records, or protected artifacts.

The current evidence is static-source only. No network request was made, no packet capture was taken, no offline test was run, and no consent or audit behavior was exercised.

---

## 24. Requirement-by-requirement traceability matrix

### 24.1 Matrix rules

**Bidirectional traceability is mandatory.** The identifiers in this matrix are stable join keys for the corrected master plan. Every plan entry and work package must name the reconciliation ID(s) it serves. Every requirement row must name a planned disposition, work package, or explicit unresolved status. A plan entry without a source ID is an orphan; a row with no disposition is a silent drop. The original 66 rows and the additive NET-01–NET-20 family must both survive.

- “Status” is the current product/requirement status, not implementation completion.
- “Implementation” is static-source evidence only unless explicitly marked as a reported runtime/test record.
- “Plan correction” names the required treatment, not an authorization to implement.
- Historical status is preserved in §4 even where the current status is stronger.

### 24.2 Chat UI and interaction

| ID | Source | Current status | Implementation evidence | Plan correction |
|---|---|---|---|---|
| REC-UI-01 | FUT-CHUI1; STAGE_STATUS polish backlog | Owner-confirmed | Global bubble border/glow exists; no verified gap repair or user halo | Add requirement; investigate visually before changing spacing |
| REC-UI-02 | FUT-CHUI2; owner instruction §4.1 | Owner-confirmed; count superseded to ~10 per property | Global role tokens only; no per-character four-property state | Add per-character user bubble, character bubble, halo, text controls |
| REC-UI-03 | T2.4; FUT-CHID1; owner §4.2 | Owner-confirmed | `chat-render.js` has 24px character thumbnail; streaming templates omit it | Require thumbnail + deterministic fallback in every assistant message |
| REC-UI-04 | Owner §4.3; historical artifact unrecovered | Owner-confirmed; semantics open | System narrator event notices exist; no narration mode | Add requirement; do not conflate with narrator events; recover semantics |
| REC-UI-05 | FUT-CHTOK1; owner §4.4 | Owner-confirmed; hidden by default | No token counter control or tokenizer-backed display found | Add hidden-by-default show/hide control; reconcile with context meter |
| REC-UI-06 | Owner §4.5 | Owner-confirmed | No chat-message action invoking image orchestration found | Add full-orchestration entry point, not a bypass |

### 24.3 Chat, per-chat, context, and memory

| ID | Source | Current status | Implementation evidence | Plan correction |
|---|---|---|---|---|
| REC-CH-01 | FUT-C1 | Deferred real requirement | No persisted composer cues found in chat schema | Include in durable schema union; keep TTS untouched |
| REC-CH-02 | FUT-CH1; owner recovery confirmation | Deferred real requirement; presentation details recovered | No per-chat Author's Note field found | Include flyout/hidden/floating/side, chat-only scope, reload retention, and new-chat clearing in schema/UI |
| REC-CH-03 | FUT-CH2 | Deferred real requirement; persistence semantics open | Global context-size/rollover controls exist; no per-chat behavior proven | Include in schema union; define persistence |
| REC-CH-04 | T3.5; FUT-CH3; owner §5 | Required | No character-scoped flush endpoint; new HTTP request per turn | Require flush/rebuild after each character response |
| REC-CH-05 | T3.1 | Open technical investigation | Approximate char-based truncation; no `cache-reuse` control | Investigate context handling separately from hot-swap |
| REC-CH-06 | T3.5 | Required | Multi-character sequencing exists; isolation not proven | Require isolation tests and canonical-history preservation |
| REC-CH-07 | T3.3; REG-1 | Required; protected relationship model | Lore resolver/store exists; full memory/lore behavior unverified | Preserve REG-1; add memory and scope integration |
| REC-CH-08 | FUT-CHAIU1 | Deferred optional | No committed AI-as-user feature | Keep deferred; require isolation if ever pursued |
| REC-CH-09 | FUT-MEM1; owner §1.4 | Required major capability | Prose-like memory records/candidates; no typed inventory/clothing/relationships | Promote to architecture driver; implementation unauthorized |
| REC-CH-10 | FUT-LORE1 | Required within memory direction | Candidate/extraction foundations; review and conflict behavior incomplete | Scope staged memory/lore intelligence; no silent overwrite |
| REC-CH-11 | T3.1; FUT-CH2; current Context Rollover UI | Separate required behavior; runtime unverified | UI/tutorial modes exist; truncation is approximate character count | Add a distinct conversation-rollover design and tests |
| ODR-08 | Owner §5 | Required behavior | No explicit single/multi flush contract | Add exact rule to plan and tests |
| ODR-11 | Owner expansion | Clarified separation | Rollover and flush remain conflated in planning language | Define and test them independently |
| ODR-04 | Owner §1.4 | Required | Partial memory foundation | Add memory subsystem to architecture and dependencies |

### 24.4 Model, provenance, hot-swap, and runtime

| ID | Source | Current status | Implementation evidence | Plan correction |
|---|---|---|---|---|
| REC-M-01 / F2-B4 | F2 ledger | Accepted | Model intelligence/profiles/templates exist | Preserve; do not claim per-turn proof |
| REC-M-02 / FUT-PROMPT1 | FUT-PROMPT1; T1.4; owner §6.1 | Required direction | Native template plus small family table; manual family select | Add full family/confidence/fallback matrix |
| REC-M-03 / FUT-C3 | FUT-C3 | Required investigation/exposure | Some flags passed; several stored settings unused | Add stored/sent/effective distinction |
| REC-M-04 / FUT-C5 | FUT-C5; owner §6.5 | Required inventory | Partial sampler UI; XTC/top-n-sigma/adaptive missing | Complete pinned-runtime matrix and conflict model |
| REC-M-05 / T3.2 | T3.2; owner §6.4 | Required research | No measured LM Studio comparison | Add benchmark plan; no parity claim |
| REC-M-06 | Owner §6.2–6.3 | Required research | Runtime flags exist; no product control | Add Kobold/context-shift gap inventory |
| REC-M-07 / FUT-HW1 | FUT-HW1 | Deferred prerequisite | Advisory resource manager; no verified release | Promote hot-swap dependency; keep implementation unauthorized |
| REC-M-08 / HW-C1–C3 | Ledger | Deferred validation | No measured matrix | Carry as deferred targets |
| REC-M-09 / MOB-1–4 | Ledger | Binding constraints | No remote implementation | Preserve; remote unauthorized |
| REC-M-10 / REG-1 | Ledger | Protected | Existing content relationships | Preserve; never call model switch “hot-swap” unqualified |
| ODR-02 / model hot-swap | Owner §1.2 | Required product behavior | Manual stop/spawn; no quiescence/readiness identity/rollback/release proof | Replace “unaccepted”; add lifecycle requirements |
| ODR-06 / provenance | Owner §2.2 | Required concept; schema open | No persisted per-turn fields | Add concept and recommendation; no F3-C amendment |
| ODR-05 / seed | Owner §2.1 | Directed if meaningful | Global seed only; no actual seed capture | Add evidence-based recommendation; no reproduction promise |
| UD-01 / per-chat model | Historical proposal; owner not decided | Unresolved | No `ChatMeta` model field | Keep separate from hot-swap/provenance |
| REC-M-11 / technique inventory | Owner architecture expansion 2026-11-11 | Required first-class mission | First-pass llama.cpp inventory only | Add multi-engine/domain inventory; do not freeze |
| REC-M-12 / architecture mapping | Owner architecture expansion 2026-11-11 | Required | Family table exists; architecture-specific mapping absent | Add dense/MoE/SWA/recurrent/hybrid/multimodal matrix |
| REC-M-13 / backend mapping | Owner architecture expansion 2026-11-11; F2-B4 | Required | Model intelligence candidate runtimes only; no load guarantee | Add backend/version/platform support matrix |
| REC-M-14 / interaction conflicts | Owner architecture expansion 2026-11-11 | Required | No conflict graph | Add tested combination states |
| REC-M-15 / automatic selection | Owner architecture expansion 2026-11-11 | Required | Absent | Add estimator, ranking, fallback, explanation |
| REC-M-16 / manual override | Owner architecture expansion 2026-11-11 | Required | Partial raw controls | Add scoped overrides, warnings, validation, reset |
| REC-M-17 / benchmark evidence | Owner architecture expansion 2026-11-11 | Required | No benchmark harness or quality measurements | Add protocol, baselines, negative results |
| ODR-10 | Owner architecture expansion 2026-11-11 | Required mission | Not treated as first-class | Replace appendix treatment |
| REC-M-18 / mission boundary | Owner baseline characterization 2026-11-11 | Required system capability, not enablement mandate | No selector; first-pass evidence only | Add staged research/evidence/compatibility/selection/validation/authorization gates |

### 24.5 Media and attachments

| ID | Source | Current status | Implementation evidence | Plan correction |
|---|---|---|---|---|
| REC-P-01 / ODR-07 | F3-C gap; owner §3 | Approved product requirement | `hasImage` only; transient image input | Add durable reference contract; no F3-C amendment |
| REC-P-02 | Owner §3 | Approved | No saved-audio chat association found | Add audio ownership/playback/restore |
| REC-P-03 | Owner §3 | Required | Tutorial says previous images not seen unless attached | Separate display from model input; no re-send by visibility |
| REC-P-04 | Owner §3 | Required | No delete/move/export/import/restore media contract | Add lifecycle contract and unavailable semantics |
| REC-P-05 | Current image metadata | Existing | Image prompt/model/options/seed recorded | Preserve; do not borrow as chat provenance |
| REC-P-06 / FUT-ST2 | F2 ledger | Deferred | Import work partial | Preserve separately |
| F2-C5 / F3-A6 | Ledger | Accepted foundations | Image+sidecar atomic commit/containment | Reuse, but do not equate with chat attachment |

### 24.6 Creative studio, avatars, and workspaces

| ID | Source | Current status | Implementation evidence | Plan correction |
|---|---|---|---|---|
| REC-C-01 / F2-B5 | F2 ledger | Approved | Separate chat/image surfaces; filtering pending | Preserve as approved; complete capability filtering |
| REC-C-02 / E.6 | Ledger | Parked | Directory work partial; not promoted | Keep parked |
| REC-C-03 / E.8 | Ledger | Parked | Image staging not a committed structure | Keep parked |
| REC-C-04 / music | F2-B7; T3.8; owner §1.3 | Approved direction; implementation deferred | No music generator | Add modality scope; no engine commitment |
| REC-C-05 / video | Current source; F3-G | Incomplete | Video engine exists; no media catalogue/sidecar/delete contract | Add F3-G corollary; not complete |
| REC-C-06 / desktop avatar | Tier4; owner §1.1 | Excluded | No production 3D desktop avatar | Remove from ordinary roadmap |
| REC-C-07 / 3D mesh + texture | Owner §1.3 | Approved direction | No production 3D mesh/texture subsystem | Add modality requirements; no avatar integration |
| REC-C-08 / general studio | Owner §1.3 | Approved direction | No unified studio subsystem | Add approved direction and shared foundations |
| ODR-01 | Owner §1.1 | Excluded | Current 2D portraits only | Record exclusion; keep portraits |
| ODR-03 | Owner §1.3 | Approved direction | Image/video only | Add studio architecture; no immediate modalities |

### 24.7 Historical status and protected records

| ID | Source | Current status | Evidence | Plan correction |
|---|---|---|---|---|
| W4 | Ledger; adjudication | Historical/superseded in substance by W5–W7 | `main.go` isolation profile set; W5–W7 verified | Do not list as live blocker; note closure gap |
| W8 | Ledger | Verified | Later record at `:1864` | Do not list as open |
| STAGE T1.1–T1.4 | STAGE_STATUS | Historical closures | Stage evidence exists | Do not treat as product-wide current verification |
| F3-C Phase 0 | Protected | Open | No change | Preserve |
| R6 | Protected | Unauthorized | No change | Preserve |
| F3-C acceptance | Protected | Not granted | No change | Preserve |
| F3-D | Protected | Blocked | No change | Preserve |
| D-ROOT | Protected | Six open | No change | Preserve |
| D-STORAGE | Protected | Blocked | No change | Preserve |
| V20 | Protected | Blocked pending F3-F | No change | Preserve |
| Gate v2 V1–V17 | Protected | Adopted | No change | Preserve |
| Gate v2 V18–V24 | Protected | Proposed | No change | Preserve |
| A9/A10/A12/A13 | Protected | Closed | No change | Preserve |
| TTS/F2-B7 | Protected | Frozen | No change | Preserve |
| F2-B5 | Approved | Approved direction | No change | Preserve |
| E.6/E.8 | Parked | Parked | No change | Preserve |
| REG-1 | Protected | Protected | No change | Preserve |
| MOB-1–MOB-4 | Binding | Conditional constraints | No remote implementation | Preserve |
| HW-C1–C3 | Deferred | Not started | No hardware matrix | Carry |

### 24.8 Source dates and dependency keys

Several historical ledger sections are undated. Their date is recorded as **unestablished in the ledger**, not invented. The applicable source-date key is:

| Source | Historical date |
|---|---|
| Tier 0–4, E.1–E.8, T2.4, T3.1/T3.3/T3.5 | Undated in the ledger; predates the dated F2 ledger |
| F2-A/F2-B/F2-C, FUT-CH1–CH3, FUT-C1–C5, FUT-CHUI/CHTOK/CHID/CHAIU/MEM/LORE | 2026-11-09 |
| F3-A/F3-A-R1 | 2026-11-09 |
| HW-C1–C3, MOB-1–MOB-4, W1–W28c | 2026-11-10 |
| REG-1 | 2026-11-09 (protected record) |
| Tier-4 avatar note | Undated in the note |
| T1.1–T1.4 stage records | 2026-10-31 |
| Group-chat and theme records | 2026-10-31 to 2026-11-01 |
| F3-B and F3-C records | 2026-11-10 to 2026-11-11 |
| `ENGINEERING-F2-STORAGE-MODEL-INTELLIGENCE-2026-11-09.md` | 2026-11-09 |
| `ENGINEERING-F2B-STORAGE-MODEL-INTEGRATION-2026-11-09.md` | 2026-11-09 |
| 2026-11-11 owner recovery instruction | 2026-11-11 |
| Owner architecture-expansion follow-up | 2026-11-11 |

Dependency keys used by the matrix:

| Key | Dependency |
|---|---|
| D1 | Stable chat/turn/character/media identity |
| D2 | One live prompt path and prompt assembly ownership |
| D3 | Tokenizer and context accounting |
| D4 | Durable media identity and ownership |
| D5 | Model/resource lifecycle, admission, readiness, release |
| D6 | Model identity, family detection, template fingerprint |
| D7 | Memory schema, retrieval, compression, reinjection |
| D8 | Per-character presentation state and theme integration |
| D9 | Image-generation orchestration and gallery/media association |
| D10 | Shared studio job/media/resource contract |
| D11 | Protected F3/REG-1 boundaries and gate evidence |
| D12 | Inference technique/architecture/backend catalogue and conflict graph |
| D13 | Benchmark harness, named hardware, and output-quality evidence |
| D14 | Destination registry, purpose-scoped consent, local audit, and backend/sidecar egress policy |

Applying the keys to the requirement groups:

| Requirement group | Dependencies |
|---|---|
| REC-UI-01/02/03 | D8, D1 |
| REC-UI-04 | D2, D3, D7 |
| REC-UI-05 | D3, D2 |
| REC-UI-06 | D9, D4 |
| REC-CH-01/02/03 | D1, D2, D3 |
| REC-CH-04/06 | D2, D3, D7 |
| REC-CH-05 | D2, D3, D5 |
| REC-CH-07/09/10 | D7, D1, D2, D11 |
| REC-CH-08 | D2, D7 |
| REC-CH-11 | D2, D3, D7 |
| REC-M-01/M-02 | D6, D2 |
| REC-M-03/M-04/M-05/M-06 | D5, D6 |
| REC-M-11/M-12/M-13/M-14 | D12, D6 |
| REC-M-15/M-16 | D5, D12, D13 |
| REC-M-17 | D13, D5, D12 |
| REC-M-18 | D12, D13, D5 |
| REC-M-07/M-08 | D5, D11 |
| REC-M-09 | D11 |
| REC-M-10 | D1, D11 |
| ODR-02/model hot-swap | D1, D5, D6, D3 |
| ODR-06/provenance | D1, D6, D3, D5 |
| ODR-05/seed | D6, D3, D5 |
| UD-01/per-chat model | D1, D6 |
| REC-P-01–P-04 | D1, D4, D9 |
| REC-P-05 | D4, D9 |
| REC-P-06 | D1, D6 |
| REC-C-01 | D10, D4 |
| REC-C-02/C-03 | D10 |
| REC-C-04/C-05/C-07/C-08 | D10, D4, D5 |
| REC-C-06 | D11 (excluded) |
| NET-01–NET-20 | D14, D5, D6, D11 |

### 24.9 Prior review, adjudication, correction, and decision-ID reconciliation

#### Muse findings F-01–F-15

| Finding | Prior adjudication | Effect of the 2026-11-11 recovery |
|---|---|---|
| F-01 Mobile/remote | Valid | Carried as MOB-1–MOB-4 constraints; remote phase remains UD-02 |
| F-02 Hot-swap terminology | Valid | Strengthened: model hot-swap is now required product behavior; content hot-swap remains REG-1 |
| F-03 Per-chat mapping | Valid | Expanded: FUT-CH1/CH2/C1 union remains; UD-01 remains for model preference |
| F-04 Provenance/seed/attachments | Valid | Product concepts clarified; schema and exact fields remain proposed |
| F-05 STAGE closures | Valid control confirmation | Preserved; no closure is product-wide verification |
| F-06 One prompt path | Valid | Still an open prerequisite (TI-03) |
| F-07 OnShutdown / W4 / W8 | Partially valid | W4 historical/superseded; W8 VERIFIED; packaged lifecycle remains open |
| F-08 Gallery/video | Valid control + minor | Video gap preserved; no media catalogue claim |
| F-09 Creative studio | Valid | Now an approved direction; F2-B5/E.6/E.8 remain distinct |
| F-10 Avatar/3D | Valid control + minor | Desktop avatar now explicitly excluded; 3D assets remain a separate studio modality |
| F-11 Dependency order | Partially valid | No F3 reorder; parallel lifecycle guard retained; remote gate retained |
| F-12 Release taxonomy | Valid | Per-category gates and no-migration-by-inclusion retained |
| F-13 Protected boundaries | Valid critical control confirmation | No change; all boundaries preserved |
| F-14 Completeness | Partially valid | HW-C1–C3 and E.4/E.5 carried; W4/W8 corrected |
| F-15 FND-15 | Partially valid | Local-only wording retained; no remote boundary designed |

#### Adjudication COR-01–COR-10

| Correction | Disposition in this addendum |
|---|---|
| COR-01 MOB classification | Carried; MP-11 and REC-M-09 |
| COR-02 Hot-swap split | Carried and expanded; MP-01, ODR-02 |
| COR-03 Per-chat union | Carried; MP-08, UD-01 |
| COR-04 Provenance taxonomy | Carried; MP-05, §8 |
| COR-05 STAGE rule + one path | Carried; §7.3, TI-03 |
| COR-06 W4/W8, HW-C, E.4/E.5 | Carried; MP-15, REC-M-08 |
| COR-07 Video-catalogue corollary | Carried; REC-C-05 |
| COR-08 F2-B5 vs E.6/E.8 vs studio; avatar tier | Carried and strengthened by ODR-01/03 |
| COR-09 Lifecycle guard; MOB gate; resource state machine | Carried; §9.2, §18 |
| COR-10 Per-category gates; no migration by inclusion; FND-15 | Carried; MP-16 |

#### Prior owner-packet OD-01–OD-06

**Important:** the owner packet and the adjudication use overlapping OD numbers for different questions. The mapping below is to the **owner packet**, whose OD-01 was the durable chat truth set.

| Prior packet ID | Prior question | Current treatment |
|---|---|---|
| OD-01 | Durable chat truth set: provenance, seed, attachments, per-chat continuation model | **Split:** provenance = required concept; seed = directed if useful; attachments = approved; per-chat model preference = UD-01 |
| OD-02 | Model hot-swap | **Resolved as required product behavior** (ODR-02); no F3-C amendment |
| OD-03 | Android/remote phase | **Still unresolved** (UD-02); MOB constraints preserved |
| OD-04 | General creative studio | **Resolved as approved direction** (ODR-03); modality scopes open |
| OD-05 | Desktop-avatar interaction | **Resolved as excluded** (ODR-01) |
| OD-06 | Story Tools scope | **Still unresolved** (UD-03) |

#### Adjudication OD-01–OD-10 ID collision

The adjudication's register uses a different numbering. To prevent cross-document mistakes:

| Adjudication OD | Question | Current treatment |
|---|---|---|
| OD-01 | Remote/mobile phase | UD-02; no work authorized |
| OD-02 | General creative studio | ODR-03; approved direction |
| OD-03 | Durable F3-E target set | Per-chat union; UD-01 remains; no adopted schema |
| OD-04 | One live prompt path | TI-03; architecture prerequisite |
| OD-05 | F3 sequence reordering | No reordering is made by this addendum; remains an owner decision only if actual reordering is proposed |
| OD-06 | Desktop-avatar scope | ODR-01; excluded |
| OD-07 | Root/path/wipe/startup defects | Protected; unchanged |
| OD-08 | TTS governance / Story Tools | TTS frozen; Story Tools UD-03 |
| OD-09 | Media/attachment identity | ODR-07 approved product requirement; contract/schema open |
| OD-10 | Adopt corrected master plan | Not adopted; still the next authority step after review |

This collision is itself a documentation defect. Future artifacts must name the source register, not just “OD-01.”

### 24.10 Network boundaries

| ID | Requirement | Outcome | Evidence | Dependencies |
|---|---|---|---|---|
| NET-01 | Local creative work and inference remain usable offline | **Unresolved** | Local backends exist; no egress-blocked end-to-end test | D14, D5 |
| NET-02 | No cloud inference or cloud model provider in the local baseline | **Contradicted** | DashScope ASR and remappable MeloTTS are exposed without separate authorization | D14, D5 |
| NET-03 | No background model lookup, telemetry, provider connection, or update check | **Contradicted** | Civitai hash lookup is called during image scan; N-01 | D14, D6 |
| NET-04 | Destination registry with purpose, feature, data classes, and authorization | **Missing** | No registry or allowlist found | D14 |
| NET-05 | Six-question consent before every request | **Partially satisfied** | Generic image-dependency consent; error toggle/preview; audio labels only | D14 |
| NET-06 | Separate browser-open, download, and organize/track permissions | **Partially satisfied** | Manual links separated from download; no indexing permission or full disclosure | D14 |
| NET-07 | Purpose-specific permission; no implied download/refresh/indexing | **Missing** | No scoped permission state | D14 |
| NET-08 | Backend/sidecar global offline enforcement | **Contradicted** | Local-Only Mode wraps renderer `fetch`/XHR only; no Go/dialer gate | D14 |
| NET-09 | Human-readable local connection audit | **Missing** | No audit history found | D14 |
| NET-10 | Disclosed cache/retention with clear and refresh control | **Partially satisfied** | Civitai name cached in manifest; no user control; download retention implicit | D14, D6 |
| NET-11 | Visible failure; no silent unverified substitution | **Partially satisfied** | Dependency failures visible; Civitai errors ignored | D14 |
| NET-12 | Android remote explicit, authenticated, encrypted, revocable, off by default | **Missing** | Remote unimplemented; MOB constraints remain | D14, D11 |
| NET-13 | External model information/downloads user-initiated | **Partially satisfied** | Dependency downloads user-initiated; Civitai lookup is automatic | D14 |
| NET-14 | No permission creep; new destination/purpose requires documented decision | **Missing** | No change-control gate | D14 |
| NET-15 | Release acceptance criterion passes | **Missing** | No network acceptance test or runtime evidence | D14 |
| NET-16 | Authenticated TLS for remote; loopback HTTP exception; no silent downgrade | **Partially satisfied** | DashScope HTTPS; MeloTTS accepts arbitrary HTTP remote URL | D14 |
| NET-17 | Credentials not logged or sent to undocumented destinations | **Partially satisfied** | Environment keys exist; no redaction/secret-handling proof | D14 |
| NET-18 | Offline mode blocks egress without breaking local workflows | **Missing** | No backend enforcement or runtime test | D14 |
| NET-19 | Every proposed path has disclosure, cache, failure, and authorization rules before implementation | **Missing** | P-01–P-08 are proposed, not authorized | D14 |
| NET-20 | No automatic update check; future check explicit and signed | **Satisfied** (absence only) | No update checker found; must remain absent | D14 |

### 24.11 Coverage conclusion

The matrix covers every major historical requirement named by the recovery mission: the UI requirements, the per-chat union including Author's Note details, T3.1/T3.3/T3.5, FUT-MEM1/FUT-LORE1, F2-B4/F2-B5, FUT-C3/C5, FUT-PROMPT1, FUT-HW1, HW-C1–C3, MOB-1–MOB-4, REG-1, E.6/E.8, the creative studio direction, the model hot-swap and provenance decisions, media durability, the two separate context-management mechanisms, the first-class capability mission, the mission-versus-enablement boundary, the first-class network boundary and its complete existing/proposed path inventory, and protected F3 status. It also records the requirements that are **not** recovered rather than fabricating them.

---

## 25. Verification and limitations

### 25.1 Checks performed

At completion, the artifact was checked for:

- exact SHA-256 identity;
- presence of every owner-decision and requirement-family identifier, including ODR-12, NET-01–NET-20, V-NET-01–V-NET-28, N-01–N-14, L-01–L-05, and P-01–P-08;
- required section coverage (§1–§25);
- absence of placeholder markers;
- unchanged principal source hashes;
- unchanged protected F3-C hashes;
- network inventory cross-checks every inspected production external path against a register entry, a loopback classification, or an explicit test/dev exclusion;
- no application source, test, build, migration, model load, configuration, master-plan, adjudication, owner-packet, contract, freeze, or acceptance edit.

The five protected F3-C hashes verified unchanged are:

| Artifact | SHA-256 |
|---|---|
| `ADJUDICATION-PACKAGE-F3C-OD18-INTEGRATION-2026-11-11.md` | `85e1968bac4eb70f489bbd7dc89057debd17e7b3bc2336aa69c036a4352f83e8` |
| `OWNER-DECISION-REQUEST-F3C-FOUNDATIONS-2026-11-11.md` | `ed8bf2e2326a2ec74838dc0c8cd4de3332b8807d2d0a2647a18cfda060bc9c7c` |
| `INVENTORY-F3C-ROOT-OWNERSHIP-2026-11-11.md` | `343104b528818df699d752c419770fdc1bfd7fe71bc1f249e0dc6847f308c6b8` |
| `INVESTIGATION-F3C-DROOT-CLOSURE-2026-11-11.md` | `494539b9dde27a78c9c4c091131a93dce99c1debd81d11fea911e6d411e2e6a6` |
| `DISPOSITION-REQUEST-F3C-DEFECTS-2026-11-11.md` | `ce32ac0bedbe216b55c065be073b2165209e0ebad4cb435897b829336d29d2c5` |

The final SHA-256 is reported in the completing response. A file cannot contain its own final content hash without changing that hash; it is therefore recorded externally rather than embedded as a self-referential value.

### 25.2 Limitations

- No runtime, UI, model-load, generation, hot-swap, memory-retrieval, media-restore, or isolation test was executed.
- Current implementation conclusions are static-source observations.
- The separate UI planning document and named expansion text were not recovered.
- Narration semantics remain unrecovered.
- The actual seed may not be exposed on Howl's current endpoint; that requires a pinned-runtime probe.
- Weight release after process exit is not proven.
- T2.9.5B F1/F2 definitions were not independently re-read here.
- W4 has no explicit formal closure record; W5–W7 supersede its practical condition.
- The live prompt path split is source-confirmed but its runtime effect on RAG and memory is unmeasured.
- The general studio is an approved direction, not an implementation plan.
- The capability mission is a research and architecture requirement, not a completed survey. No other engine or architecture family was benchmarked here.
- The Author's Note presentation details come from the owner follow-up and the ledger; the exact original UI document is still missing.
- The presence of two F2 reports does not prove the full original document set or every original requirement survived.
- No capability technique or backend feature has passed the research/evidence/compatibility/selection/validation/authorization sequence; the baseline does not authorize any of them.
- The network inventory is static-source evidence. No outbound request, packet capture, offline run, cancellation test, consent/audit test, remote ASR/TTS test, or Android test was performed.
- Local-Only Mode was not demonstrated to block Go `net/http`, `net.Dial`, sidecar egress, or Wails-mediated actions; its inspected implementation is renderer-only.

### 25.3 Final status

This addendum establishes a traceable product baseline and the corrections needed before a new master-plan revision. It does **not** adopt that revision, amend F3-C, authorize implementation, or convert any approved direction into a release commitment. It also does not freeze the first-pass capability inventory: hardware- and architecture-aware inference is a required first-class mission whose research outputs and measurements remain outstanding. Possessing the capability system is not permission to enable every technique or implement every backend feature. The network boundary is likewise required but not satisfied: the existing path inventory is not a validation record, and no external destination is authorized by this document. The corrected plan must preserve bidirectional traceability and treat estimates as provisional until the investigations establish the implementation baseline. The reconciliation remains the next authority step.
