# HOWL CHAT — COMPLETE PRODUCT ARCHITECTURE, REQUIREMENTS & DEPENDENCY MASTER PLAN

**Date:** 2026-11-11  
**Mission type:** Read-only architectural inventory and planning  
**Status:** **PLANNING ARTIFACT ONLY — NOT A CONTRACT, NOT A FREEZE, NOT AN ACCEPTANCE, NOT AN IMPLEMENTATION AUTHORIZATION, NOT R6**  
**Implementation authorization:** NONE  
**Architecture-freeze authorization:** NONE  
**R6 authorization:** NONE  
**F3-C acceptance:** NONE  

**Authority for requirements:** `PROJECT_ISSUES_AND_EXPANSION.md` remains the single authoritative requirements and expansion ledger, as recorded in `docs/review/ARCHITECTURE-F2C-CONTENT-OWNERSHIP-CUTOVER-2026-11-09.md` §9. This document inventories that ledger; it does not replace it.  
**Authority for F3-C:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` remains controlling. This document does not amend it.  
**Source evidence:** static documentation and source inspection only. No test, build, model load, native run, migration, or destructive operation was performed for this artifact.

---

## 1. Executive architectural overview

Howl Chat is a local-first, single-user Windows desktop creative application built as a Wails v2 shell around a Go backend and an embedded `frontend/dist` WebView2 UI. One `App` object is bound to Wails; the UI is page-based rather than router-based. The current product surface is eight HTML pages: chat (`index.html`), character, world, scenario, lorebooks, image generator, audio settings, and story tools.

The intended product is broader than the current release surface. The authoritative ledger describes a world-generation and character-chat product whose differentiators are application-owned facts, continuity, local model orchestration, and creative production across text, image, audio, video, and eventually 3D avatars. The architecture therefore has to support capabilities that are deliberately not implemented yet. Freezing only today's feature set would create avoidable rework in exactly the areas named by the F3-C review: durable chat identity, model provenance, per-chat continuation, attachments, and hot-swap behavior.

The most important architectural finding is that Howl currently has **several concurrent ownership models**, not one product-wide foundation:

- **Application lifecycle:** improved since the 2026-11-09 assessment. `main.go` now registers `OnShutdown: app.OnShutdown`, and shutdown is idempotent with an admission barrier. The older finding that shutdown is unregistered is **superseded**.
- **Content ownership:** F3-B makes lorebooks canonical in `internal/lorebookstore`; F3-C implements a tested world/scenario store foundation, but F3-C Phase 0 is open, acceptance is not granted, and the six D-ROOT defects remain open. Characters, chats, images, and videos retain their pre-cutover ownership.
- **Chat assembly:** the live streaming path in `app-prompt.go` builds prompts directly from `App` fields; `chat.Service` is a separate canonical consumer used by F3-C and native acceptance seams. RAG is attached to `chat.Service`, not demonstrably consumed by the live `App` prompt path. This is a product-path split, not merely an implementation detail.
- **Persistence:** chats are under `%APPDATA%/HOWLChat/chats`; memory has a separate `StateRoot/Chats` surface; latent `cwd/Chats` candidate/continuity roots exist; characters/worlds/scenarios still have browser and native/legacy surfaces; images and videos use a gallery/output directory. The F3-C defect inventory documents the resulting split ownership and wipe/identity consequences.
- **Runtime authority:** `resource.Manager` records loaded models, but registration is not a verified readiness boundary, image auto-load does not participate consistently, TTS is outside admission, and hot-swapping has no in-flight-stream coordination or rollback after the old server is stopped.
- **Provenance:** chat history stores text, role, sender, timestamp, character ID, and a `hasImage` boolean only. It does not store the model, engine, template fingerprint, effective generation parameters, seed mode/actual seed, or durable attachment identity. The F3-C Phase 0 assessment records this as a **contract gap**, not merely a missing implementation.

The product can remain deferred feature-by-feature, but the shared contracts cannot be designed as if those features will never arrive. The expensive decisions are: chat serialization, root identity, media/attachment identity, model identity and readiness, prompt/template ownership, and one generation/cancellation authority.

**This plan does not promote any deferred feature, close any defect, alter any gate, or authorize implementation.**

---

## 2. Method, evidence classes, and limits

### 2.1 Sources recovered

The inventory was recovered from, at minimum:

- `PROJECT_ISSUES_AND_EXPANSION.md` — Tiers 0–4, E.1–E.8, F2-B/F2-C, F3-A/F3-A-R1, deferred feature register, FUT-HW1/FUT-HW2, FUT-PROMPT1, F3-B/F3-C records.
- `docs/tier4-avatar.md` — reactive 3D avatar, body library, skin pipeline, graceful tiers, licensing, open questions.
- `docs/deferred-model-ecosystem.md` — dependency-aware model resolution, role-based identity, curated collection, offline-first connectivity, per-connection disclosure.
- `docs/review/ASSESSMENT-howl-application-reality-check-shared-foundations-2026-11-09.md` — 15-subsystem baseline, shared-foundation map, risks, engineering order.
- `docs/review/ARCHITECTURE-F2C-CONTENT-OWNERSHIP-CUTOVER-2026-11-09.md` — current ownership map, dual-storage analysis, media ownership, cutover sequence.
- F3-B/F3-C records, including `FREEZE-F3C-*`, `ASSESSMENT-F3C-PHASE-0-FINAL-REQUIREMENTS-IMPACT-2026-11-11.md`, `DISPOSITION-REQUEST-F3C-DEFECTS-2026-11-11.md`, and `GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`.
- `STAGE_STATUS.md`, `TECH_DEBT.md`, `CHECKLIST.md`, `docs/audio-architecture/*`, and current source files under `app-*.go`, `internal/**`, and `frontend/dist/**`.

The repository contains a large review corpus. The inventory below is complete for the **named and discoverable product requirements in the authoritative ledger and the separate avatar/model-ecosystem notes**. It is not a claim that every prose sentence in every review artifact has been converted into a requirement.

### 2.2 Evidence classifications

| Classification | Meaning in this document |
|---|---|
| **Source-confirmed** | The capability, field, binding, or absence is directly supported by inspected source. It does not prove runtime behavior. |
| **Runtime-confirmed** | Evidence from an actual running packaged/native application exists in the cited record. This mission did not create new runtime evidence. |
| **Test-confirmed** | A cited test or verification record reports a passing result. This mission did not rerun it. |
| **Inferred** | The conclusion follows from source structure or documented behavior but was not directly executed. |
| **Unknown** | Evidence is insufficient to classify. |
| **Planned** | A requirement exists in documentation; no implementation is claimed. |
| **Superseded** | An earlier statement is contradicted by current source or a later accepted record. |

A source path is evidence of design intent, not proof that a feature works in production. This distinction is especially important for chat, image, audio, video, 3D, Android/QR, and all model-backed workflows.

### 2.3 Preserved current position

The following position is reproduced from the mission and the F3-C records, and is not changed by this plan:

| Gate / item | Current position |
|---|---|
| F3-C Phase 0 | **OPEN** |
| R6 | **UNAUTHORIZED** |
| F3-C acceptance | **NOT GRANTED** |
| F3-D | **BLOCKED** |
| D-ROOT | **OPEN** |
| D-STORAGE | **BLOCKED** |
| V20 | **BLOCKED pending F3-F** |
| Gate v2 V1–V17 | **ADOPTED** |
| Gate v2 V18–V24 | **PROPOSED — NOT ADOPTED** |
| A9/A10/A12/A13 | **CLOSED; not reopened by this plan** |
| Six F3-C defects | **ALL OPEN** |
| Architecture freeze | **GRANTED for F3-C; does not equal acceptance or implementation authorization** |

Six independent defects remain open. This plan does not close, merge, downgrade, or reassign them:

| Defect | Scope, as recorded | Status |
|---|---|---|
| `DEF-F3C-ROOT-01` | Split chat/memory ownership | **OPEN — Critical** |
| `DEF-F3C-WIPE-02` | Incomplete content wipe across F3-C surfaces | **OPEN — Critical** |
| `DEF-F3C-IDENTITY-03` | Empty `SetChatID` cannot clear memory identity | **OPEN — High** |
| `DEF-F3C-REACH-04` | Latent `cwd/Chats` candidate/continuity roots | **OPEN — Medium** |
| `DEF-F3C-STARTUP-05` | Startup read can nil service and skip RAG/content-library wiring | **OPEN — Medium** |
| `DEF-F3C-ALIAS-06` | No canonical path identity or alias rejection | **OPEN — High** |

---

## 3. Complete subsystem inventory and current implementation baseline

This table is the **subsystem ownership map** and the **implementation baseline** in one place. “Primary status” is deliberately conservative and separates source presence from verification.

| ID | Subsystem | Backend owner / locations | Frontend owner / locations | Primary status and evidence |
|---|---|---|---|---|
| S1 | Shell, lifecycle, process supervision | `main.go`; `app.go`; `app-lifecycle.go`; `internal/proc/*`; `app-llama.go`; `internal/image/engine/sdcpp`; audio supervisors | Wails window only; no shared app shell | **Source-confirmed, partial.** `OnShutdown` is registered; shutdown admission and owned-process job model exist. Full packaged lifecycle acceptance remains **Unknown**. The old “OnShutdown never wired” finding is **Superseded**. |
| S2 | Chat, conversation state, streaming, prompt assembly | `app-chat.go`; `app-prompt.go`; `app-think*.go`; `internal/backend/chat/*`; `internal/backend/llama/*`; `internal/backend/prompt/builder.go` | `chat.js`, `chat-stream.js`, `chat-session.js`, `chat-actions.js`, `chat-state-bridge.js`, `data-sync.js` | **Source-confirmed, incomplete.** Live path streams through `App` and `llama.HTTPClient`; `chat.Service` is separate. No per-turn model/seed provenance; no per-chat model; image attachment is transient. Runtime end-to-end remains **Unknown**. |
| S3 | Character management | `app-characters.go`; `internal/backend/chat/character*.go`; `internal/backend/types/entities/entities.go` | `character*.js`; `character-form.js`; `character-grid.js`; `character-storage.js`; IndexedDB | **Source-confirmed, implemented/unverified.** Characters use browser storage for the grid plus a native `Characters/<name>/character.json` fallback; F3-D is not authorized. AI-assisted fields are **Planned**, not implemented. |
| S4 | Worlds and scenarios | `internal/worldstore/*`; `internal/contentimport/*`; `app-libraries.go`; `app-imports.go`; `app-characters.go` | `world*.js`; `scenario*.js`; browser `howl-worlds` / `howl-chat-scenarios` | **Mixed.** World/scenario store foundation is **Test-confirmed** for its own tests. F3-C Phase 0 is open; native acceptance and consumer completeness are not established. Browser editors still write the explicit import source, not the canonical store. |
| S5 | Lorebooks, triggers, conflict resolution | `internal/lorebookstore/*`; `internal/backend/lorebook/*`; `app-lorebooks-store.go` | `lorebook-repository.js`; `lorebooks.js`; `lorebook-manager.js`; cache `howl-chat-lorebooks` | **Source-confirmed and Test-confirmed (reported).** F3-B is **CLOSED / ACCEPTED WITH DOCUMENTED LIMITATIONS**. Prompt synchronization with a running model remains unverified. |
| S6 | Memory and continuity | `internal/backend/memory/*`; `app-memory.go` | `memory-store.js`; `memory-candidates*.js`; `character-memories.js` | **Source-confirmed, incomplete.** Manual/automated memory, pinning, decay, and candidate review exist. Structured inventory, clothing/equipment state, and application-owned continuity are **Planned** (`FUT-MEM1`), not implemented. |
| S7 | RAG / document knowledge | `internal/backend/rag/*`; `app-rag-upload.go`; `chat.Service` RAG methods | `rag-knowledge.js`; story-tools knowledge panel | **Source-confirmed, integration gap.** `chat.Service.SetRAGService` is wired, but the live `App.buildPromptMessages` path does not demonstrably call it. Source presence is not proof that RAG affects replies. |
| S8 | Image generation and editing | `internal/image/engine/sdcpp/*`; `internal/image/job/*`; `internal/image/gallery/*`; `internal/image/metadata/*`; `app-image*.go` | `image-generator*.js`; `image-workflows.js`; `ig-pose.js`; `image-recall.js` | **Source-confirmed, implemented/unverified.** txt2img, img2img, inpaint, outpaint, upscale, ControlNet, LoRA, workflows, pose, batch, seed, and gallery paths exist. Documented T0.3/T0.4 and runtime acceptance remain unresolved. |
| S9 | Video generation and processing | `internal/image/engine/sdcpp/video.go`; `GenerateVideo` in `app-image-generate.go` and job manager | `image-generator-ui.js`; video mode/families in `image-model-dirs.js` | **Source-confirmed, incomplete.** SVD/WAN/LTX/Hunyuan-style families and video options exist. Video has no gallery sidecar/catalogue/delete contract comparable to images. F3-G owns the future media contract. |
| S10 | ASR / speech recognition | `internal/audio/asr/*` (Whisper, Qwen, llama adapter, fallback); `internal/audio/runtime/*`; `app-audio.go` | `audio-settings.js`; `audio-service.js`; `audio-wav.js`; chat microphone path | **Source-confirmed, partial verification.** Provider architecture, capability reporting, and transcription bindings exist. Remote-engine availability is explicitly `remote-unverified`. |
| S11 | TTS / voice synthesis | `internal/audio/tts/*` (Piper, Kokoro, MeloTTS, helper, cache); `internal/audio/*`; `app-audio.go`; `app-tts-catalog.go` | `audio-settings.js`; per-character voice controls | **Source-confirmed, governance unresolved.** TTS contract/lifecycle decisions remain governed elsewhere; this plan does not alter them. No claim of accepted end-to-end operation. |
| S12 | Audio generation / music | No production generator found | No dedicated page/control found | **Planned only** (`T3.8`; future audio expansion in `F2-B7`). Do not describe as implemented. |
| S13 | Vision and multimodal input | `internal/backend/llama/http_client.go`; `app-prompt.go` image path; `internal/backend/vision/*`; model projector pairing | Chat image attach in `chat-input.js` | **Source-confirmed, current-turn only.** One image can be sent when the loaded model/profile advertises vision. Image cues, multi-image context, persistence, and model-aware anchoring are **Planned/Unknown**. |
| S14 | Model discovery, registry, intelligence | `internal/modelintel/*`; `internal/backend/gguf/*`; `internal/image/registry/*`; `app-model-*.go`; `app-model-intelligence.go` | `chat-model.js`; `model-directory.js`; `image-model-dirs.js`; `model-settings.js` | **Source-confirmed, implemented/unverified.** Registry, evidence-based profiles, template inspection, image-family detection, and discovery records exist. Filename fallback and unverified compatibility remain risks. |
| S15 | Runtime model loading and hot-swap | `app-model.go`; `internal/backend/model/manager.go`; `internal/backend/llama/*`; `app-llama.go` | `chat-model.js`; model browser/settings | **Source-confirmed, partial.** Manual load/unload and pre-stop admission checks exist. `LoadModel` does not coordinate an active stream; after the old server stops, spawn/startup failure has no old-model rollback. Hot-swap is **not an accepted product behavior**. |
| S16 | Hardware monitoring and resource orchestration | `internal/resource/*`; `app-resource.go`; `app-vram-arbitration.go`; `internal/audio/hw/*` | VRAM dialogs/status; image settings | **Source-confirmed, incomplete.** Resource manager and chat/image arbitration exist. Admission is advisory; image auto-load registration is inconsistent; TTS is outside admission; no verified unload/resource-release proof. `FUT-HW1` is **Deferred**. |
| S17 | Content libraries, import/export, ownership | `internal/library/*`; `internal/contentimport/*`; `internal/paths/*`; `app-libraries.go`; `app-imports.go`; `app-wipe.go` | `content-libraries.js`; `content-save.js`; import UI | **Source-confirmed, partial.** Registry roles, inventory, explicit world/scenario import, lorebook import, and staged wipe exist. Characters/chats/images/videos are not cut over. Export is category-limited. |
| S18 | Settings and configuration | `internal/settings/*`; `internal/image/settings/*`; `app-settings-reset.go`; `app-family-defaults.go` | `settings-state.js`; `real-settings.js`; `settings-controls.js`; `settings-actions.js` | **Source-confirmed, partial.** Unified and image-specific stores exist with mixed atomicity and mixed consumption. Audio settings are restored without activation by design. No restart round-trip was performed here. |
| S19 | UI navigation, theme, tutorials | None; static assets | `frontend/dist/*.html`; `navigation.js`; `theme.js`; `tutorials-*`; `menu*.js` | **Source-confirmed, incomplete.** Page navigation and shared scripts exist. Earlier light-theme, ghost-ID, and duplicate-element findings remain relevant unless separately closed. No visuals were reviewed for this artifact. |
| S20 | Story tools and session utilities | `app-rag-upload.go`; memory/export bindings; no cohesive owner | `story-tools*.js`; `story-tools.html` | **Unknown.** Page and scripts exist; intended committed scope is not established. Do not invent scope from controls. |
| S21 | Errors, logging, diagnostics, recovery | `internal/errors/*`; `app-errors.go`; `internal/backend/loader` failure records; `app-image-diagnostics.go` | `error-*.js`; notifications | **Source-confirmed, partial.** Structured logging/reporting and loader diagnostics exist; many paths still log/ignore rather than surface one recovery contract. |
| S22 | Desktop avatars / 3D characters | No production 3D implementation found; image pipeline and TTS are dependencies | Only 2D avatar rendering in chat/character UI | **Planned only** (`docs/tier4-avatar.md`, T4). Tier B/C/D are future; UltraTex is external; licensing flag applies. |
| S23 | Android remote access / QR connection | No server, pairing, transport, or auth implementation found | No page/control found | **Unestablished.** It is a discovery target from the mission, not a recovered project requirement. |
| S24 | Creative studio / general asset generation | No dedicated subsystem found | No dedicated page found | **Unestablished.** Do not treat as committed scope. |

### 3.1 Baseline summary

| Status class | Subsystems | Meaning |
|---|---|---|
| Runtime-confirmed in cited records | F3-B (with limitations); selected F3-C native tests | Evidence exists in project records, not from this mission |
| Test-confirmed in cited records | F3-B, F3-C store/import, selected lorebook/memory/chat tests | Tests exist; acceptance may still be open |
| Source-confirmed but not runtime-verified | S2, S3, S4, S6, S7, S8, S9, S10, S11, S13, S14, S15, S16, S17, S18, S19, S21 | Code/design exists; production behavior not proven |
| Planned only | S12, S22, most future requirements | Documentation exists; no implementation claim |
| Unestablished | S23, S24 | No authoritative project requirement found |
| Superseded | “OnShutdown unregistered” from the 2026-11-09 assessment | Current `main.go` registers it |

**No subsystem is classified WORKING/VERIFIED by this mission.** The only accepted-with-limitations subsystem is F3-B lorebook ownership, and that acceptance is for its defined ownership scope, not for model inference or all REG-1 behavior.

---

## 4. Complete recovered requirements and additions inventory

### 4.1 Release blockers and core mechanics (Tiers 0–2)

These are existing product requirements, not future additions. Their inclusion prevents the additions list from being mistaken for the whole backlog.

| ID | Original requirement, compressed but meaning-preserving | Recovered status |
|---|---|---|
| T0.1 | User messages do not appear until restart; must diagnose duplicate send vs render desync without clearing caches | **UNVERIFIED** in ledger; source has a duplicate-user guard, but no runtime closure |
| T0.2 | Model receives the user message twice; test independently from T0.1 | **UNVERIFIED** |
| T0.3 | Image model will not load; validate stored paths, add model-directory chooser, show real errors | **UNVERIFIED / open blocker** |
| T0.4 | Image preview is broken; prompt must produce a correctly previewed image | **UNVERIFIED / open blocker** |
| T1.1 | State ownership: tab switching must not unload models or drop connections; settings persist visually | Ledger **UNVERIFIED**; `STAGE_STATUS` records a 2026-10-31 closure. **Current F3-C evidence shows ownership is still split**, so the old closure is not evidence of product-wide completeness. |
| T1.2 | VRAM arbitration between chat and image with confirmation, re-check, unload/load in both directions | Ledger **UNVERIFIED**; `STAGE_STATUS` records closure. Current resource authority remains incomplete, so classify **partial**. |
| T1.3 | Think tags/reasoning: on/off/hidden, intensity where supported, hidden not counted, thoughts discarded from history | Ledger **UNVERIFIED**; `STAGE_STATUS` records closure. Source has `app-think*.go` and splitter; runtime not verified here. |
| T1.4 | Family-specific sampler defaults and automatic prompt formatting; no ignored controls | Ledger **UNVERIFIED**; `STAGE_STATUS` records closure. Current model/template evidence remains partial. |
| T1.5 | Header name “Howl Chat” | Ledger **UNVERIFIED**; no independent runtime check |
| T2.1 | Theme unification, especially light theme | **UNVERIFIED** |
| T2.2 | Ghost elements/stale markup and duplicate IDs swept | **UNVERIFIED** |
| T2.3 | Image page organized toward Invoke AI: preview, grouped controls, history; depends on T0.4 | **UNVERIFIED** |
| T2.4 | Speaker identity: thumbnail or initial; resized thumbnail, fixed size, theme-safe; prerequisite for multi-character legibility | **UNVERIFIED**; source has 2D avatar rendering, but the complete presentation contract is not runtime-confirmed |
| T2.5 | Character/world/scenario image sets up to ~10 with default image; investigate data/UI/lifecycle regression | **UNVERIFIED**; source contains multi-image storage paths, but regression status is not established |

### 4.2 Tier 3 — the original additions list

> The ledger is explicit: **“Nothing here starts until Tiers 0–2 are closed.”** These are not current implementation commitments.

| ID | Original addition, preserving qualifiers | Status |
|---|---|---|
| T3.1 | Context handling: first verify whether context flush exists; measure prefill; test `cache_prompt` / `--cache-reuse`; full model unloads are not acceptable | **UNVERIFIED / research first** |
| T3.2 | Load speed: learn from LM Studio what is portable; possible major quality-of-life win | **UNVERIFIED** |
| T3.3 | Scoped lorebook + memory: world/scenario/character scope; primary/secondary keys; word boundaries; scan depth 2–3; recursion capped at 2; no-recursion flag; fixed budget; precedence character → scenario → world; dedupe; lore budget ~15–20%; per-character memory; AI-authored lore with specificity and frequency checks | **UNVERIFIED**; substantial lorebook foundations exist, but the full addition is not complete |
| T3.4 | Per-field AI character sheets: per-field generate button, context from filled fields, target purpose, truncated siblings, return only that field; draft/accept/regenerate/lock; grammar-constrained JSON; no preset voices; authored exemplar dialogue | **UNVERIFIED / Planned**; no AI field-generation path established |
| T3.5 | Multi-character chat: 2+ characters fully separate, appropriate system prompts, context flushed between speakers; root cause is flat context | **UNVERIFIED**; partial multi-character scaffolding exists, but isolation/flush is not established |
| T3.6 | Ternary models as first-class: load without new dependencies by pointing at the file | **UNVERIFIED / Planned** |
| T3.7 | Audio and TTS end to end through the UI; headphones required for verification | **UNVERIFIED**; audio architecture exists, acceptance not established |
| T3.8 | Audio generation: newer small audio model, 8GB short clips; play generated audio in app | **UNVERIFIED / Planned**; no generator found |
| T3.9 | Diffusion version check: confirm latest patched Stable Diffusion and update cost | **UNVERIFIED** |
| T3.10 | LLM writes image prompt — research gated: quality/cost/swap; possible resident vision inspection | **UNVERIFIED / research only** |
| T3.11 | Compatibility for everyone else: generalise auto context/layering, ship beta, fix from reports | **UNVERIFIED** |
| T3.12 | Whole-app optimisation pass: separate stage after everything above | **UNVERIFIED / unscheduled** |

### 4.3 Tier 4 and E.1–E.8 — parked concepts

| ID | Original parked concept | Status |
|---|---|---|
| T4.1 | Reactive 3D character: AI-generated skin, prebuilt skeleton/mesh, emotional animation, mouth driven by TTS amplitude; enormous VRAM cost, likely incompatible with multi-character chat; no Live2D hand animation; revisit only if foundations strong | **PLANNED / PARKED / UNVERIFIED** — detailed in `docs/tier4-avatar.md` |
| T4.2 | Multi-character with avatars: sub-case of 3D avatar with the same resource problem | **PLANNED / PARKED** |
| E.1 | Model-aware image anchoring: world/scenario/character reference images in context; 1,000–1,500 visual tokens per high-res image; three images ~3,000–4,500 tokens; dilution risk; portrait/world at initialization or scene change, action images on new environment/event, downscale to 512/768 | **PARKED / UNVERIFIED**, behind T3.3/T3.5 and vision model situation |
| E.2 | Multi-image vision support as a capability question, not a standalone promise; belongs with E.1 | **PARKED / UNVERIFIED** |
| E.3 | Reference tools only: LM Studio (Hugging Face, VRAM fit calculator), Open WebUI (feature density, auth/RAG/audio/web search/tools), Jan (offline, Apache 2.0, no telemetry), AnythingLLM (local document RAG) | **REFERENCE ONLY**, not scope |
| E.4 | Friction points to avoid: silent failures/frozen UI; hidden paths/locked-in data; silent context truncation; rigid system prompts; fragile ports/dependencies. Recorded fix directions only | **UNVERIFIED**, check against current app |
| E.5 | “Fake strengths” to avoid: massive plugin ecosystems; universal system prompts; endless sampler sliders; unrestricted full-context re-injection. Recorded defenses only | **DESIGN GUARDRAILS**, not features |
| E.6 | Custom workspace layout: text under Documents, images under Pictures, `Pictures\Howl Chat\...` tree | **PARKED / UNVERIFIED** |
| E.7 | “Curtain” vault: hidden folder + PIN + recovery path + optional byte obfuscation; explicitly not cryptographic armor | **PARKED / UNVERIFIED / not assumed wanted** |
| E.8 | Image staging board structure mirrored into the chosen image directory | **PARKED / UNVERIFIED** |

### 4.4 F2-B deferred future requirements

These are the 16 preserved F2-B rows. They remain deferred and assigned to their original roadmap slots.

| ID | Requirement | Eventual owner / gate |
|---|---|---|
| FUT-C1 | Up to five chat message cues in composer, persisted with chat | Chat UI/prompt assembly; schema versioning, cue mapping, user control; TTS contract untouched |
| FUT-C2 | Image cue support in chat; final count established by research | Chat UI/multimodal prompt assembly; attach/read pipeline, capability evidence, bounded image research |
| FUT-C3 | llama.cpp feature exposure: Flash Attention, GPU offload, basic multi-GPU, hardware-specific optimizations | Chat runtime/llama.cpp adapter; verified flags, hardware probe, safe fallback |
| FUT-C4 | Architecture/family/model-era-specific options, capability-gated | Model intelligence + chat settings; reliable evidence, allowlist, never filename-only certainty |
| FUT-C5 | Comprehensive sampling controls including XTC and research into three unidentified samplers | Chat sampling; runtime support matrix and named/versioned controls |
| FUT-ST1 | SillyTavern card import from PNG metadata and JSON | Character import adapter; read-only discovery exists; explicit preview and copy-on-import |
| FUT-ST2 | Intelligent mapping of supplied card fields into Howl fields | Import adapter; mapping table, conflict policy, provenance |
| FUT-ST3 | Preserve and parse information dumped into Description | Import adapter/character form; conservative parser, original retained, user review |
| FUT-ST4 | Character LoRA support | Character model config/llama.cpp adapter; discovery, validation, runtime loading, resource accounting |
| FUT-ST5 | Lorebook import and compatible configuration adoption | Lorebook import; preview, validation, conflict handling, no silent overwrite |
| FUT-CH1 | Chat-scoped Author's Note with popup/floating/hide/side behavior | Chat composer/settings; per-chat persistence, insertion order, visible enable state |
| FUT-CH2 | Single-character chats retain configured context capacity before rollover/memory strategies | Context manager/memory; token accounting, response reserve, observable diagnostics |
| FUT-CH3 | Multi-character context isolation/continuity with flushing or stronger alternative | Multi-character engine; per-character boundaries, shared scene state, testable isolation |
| FUT-IMG1 | Hardware-aware RAM/VRAM placement for image generation | Image resource manager; probe, capability matrix, fallback, measurement |
| FUT-IMG2 | Safe CPU/system-memory offloading where supported | Image resource manager; backend-supported only, guardrails |
| FUT-IMG3 | Maximum practical OOM prevention | Image/resource coordination; preflight, bounded retries, cancellation, error reporting preserving output |

### 4.5 Deferred UI, identity, memory, lore, and principle register

| ID | Requirement | Status |
|---|---|---|
| FUT-CHUI1 | Restore cyan halo around user bubble; consistent gap between bubbles; preserve neon/TRON readability | **DEFERRED — DOCUMENTED ONLY** |
| FUT-CHUI2 | User appearance control with ~20 coordinated colors affecting border/halo/font; independent per-character settings in multi-character chat | **DEFERRED — DOCUMENTED ONLY** |
| FUT-CHTOK1 | Independent context-usage percentage, user-message token count, response token count; no false precision | **DEFERRED — DOCUMENTED ONLY** |
| FUT-CHID1 | Character thumbnail/fallback identity; circular initial fallback; per-character border color; group distinguishability | **DEFERRED — DOCUMENTED ONLY** |
| FUT-CHAIU1 | Optional AI-as-user/persona takeover; explicit activation/deactivation; isolation/flushing; no personality or knowledge bleed; separate feasibility review | **DEFERRED — MAJOR OPTIONAL CAPABILITY** |
| FUT-MEM1 | Application-side structured memory and continuity: inventory, item ownership, acquisition/transfer/consumption/use, health potion state, clothing/outfit continuity, equipment state changes, broader session facts; scope-aware prompt selection | **DEFERRED — MAJOR CAPABILITY**; no structured inventory implementation found |
| FUT-LORE1 | Automated lore/memory intelligence overhaul: natural-language triggers, durable fact extraction, review/edit/correct/override, conflict/duplicate checks against manual lore, no silent overwrite, unobtrusive during roleplay, long-session reliability | **DEFERRED — DOCUMENTED ONLY** |
| FUT-PRINCIPLE1 | Governing principle: **immersion is the goal; continuity is the capability that supports it** | **RECORDED PRINCIPLE — NOT A FEATURE** |

### 4.6 Deferred hardware and prompt/model-format intelligence

| ID | Requirement | Status |
|---|---|---|
| FUT-HW1 | Unified application-wide hardware-aware resource orchestration: identification, estimation, admission, placement, live monitoring, workload scheduling, unload/offload, process cleanup, verified transitions; graceful CPU/RAM layering; broad hardware matrix; power-user overrides; acceptance AC-HW1-1..12 | **DEFERRED — DOCUMENTED ONLY, NOT IMPLEMENTED, NOT AUTHORIZED** |
| FUT-HW2 | Possible hardware monitor presentation over FUT-HW1 telemetry; placement/metrics undecided | **POSSIBLE / UNDECIDED**; not a second resource manager |
| FUT-PROMPT1 | Expanded model-family support with evidence-based automatic prompt/chat formatting: architecture recognition, templates, Jinja/ChatML, reasoning tags, verified/inferred/conflicting/unknown states, user overrides, no marketing-name guessing; AC-PROMPT1-1..8 | **DEFERRED — DOCUMENTED ONLY** |

### 4.7 F3-C proposal-only additions (not adopted)

The F3-C Phase-0 assessment and amendment proposals identify product behaviors that the frozen target set cannot preserve. They are **review proposals**, not adopted requirements and not implementation scope. They matter because a final chat schema made without them would require another migration.

| ID (planning label only) | Proposal | Current status |
|---|---|---|
| PROP-F3C-01 | Per-turn generation provenance: actual model, engine/backend, relevant generation parameters, template fingerprint, seed state that produced each assistant turn | **PROPOSED — NOT ADOPTED**; frozen target set lacks fields |
| PROP-F3C-02 | Per-chat continuation model preference, separate from historical provenance and runtime loaded state | **PROPOSED — NOT ADOPTED**; current selection is global `localStorage`/settings |
| PROP-F3C-03 | Explicit-versus-random seed state, including the actual seed and valid explicit zero | **PROPOSED — NOT ADOPTED**; current global `-1` cannot distinguish random from unknown |
| PROP-F3C-04 | Durable image/attachment references; `hasImage` boolean is insufficient; missing references must remain unavailable, never substituted | **PROPOSED — NOT ADOPTED**; transient base64/current-turn path |
| PROP-F3C-05 | Hot-swap during an active conversation with readiness/failure/rollback rules and no false “loaded” publication | **ADJACENT PRODUCT BEHAVIOR, NOT AN F3-C TARGET-SET REQUIREMENT**; manual path exists, stream coordination/rollback absent |

### 4.8 Mission discovery targets not established in the recovered ledger

These were requested as investigation targets but were **not found as authoritative project requirements**. They must not be silently converted into commitments.

| Target | Finding |
|---|---|
| Android remote access | **Unestablished.** No Android server, pairing, QR, transport, auth, or remote-session requirement found in the ledger or source. The only “android” string found in frontend JS is a user-agent check in `error-system.js`. |
| QR-based connection | **Unestablished.** No QR generation, pairing token, or remote-origin contract found. |
| Desktop character avatars | **Partially established** as the parked reactive 3D avatar in `docs/tier4-avatar.md`; current 2D avatars are source-confirmed. Do not conflate 2D avatar rendering with the 3D desktop-avatar proposal. |
| Additional 3D mesh generation | **Partially established** only as Tier D (TripoSR-class generated geometry) inside the parked avatar ladder. No independent mesh-generation requirement found. |
| Creative studio / general asset generation | **Unestablished.** No dedicated requirement, page, service, or owner found. |
| AI-assisted world/scenario creation | **Partially established.** AI-authored lore (T3.3) and per-field character sheets (T3.4) are documented; a general world/scenario generator is not. |
| Image context-menu operations | **Not formally specified.** Image recall/import-seed/use-as-init/use-as-mask controls exist in source; no authoritative context-menu requirement was recovered. |
| Additional audio capabilities | **Established** as T3.7/T3.8 and F2-B7; music/audio expansion remains deferred until subsystems are ready. |

---

## 5. Shared architectural foundations

The same foundation is used by multiple subsystems. The names below are planning names, not new code packages.

| ID | Foundation | Current state | Duplication / incompleteness / contradiction | Dependents |
|---|---|---|---|---|
| FND-01 | Storage and persistence ownership | F3-B canonical lorebooks; F3-C world/scenario store foundation; chats/memory/images/videos still split | **Contradictory.** AppData chats, `StateRoot/Chats` memory, latent `cwd/Chats` candidate/continuity, native `chat.Service` context, browser stores, and gallery roots coexist. ROOT-01, WIPE-02, REACH-04, STARTUP-05. | Chat, characters, worlds, scenarios, lore, memory, media, settings, import/export, wipe |
| FND-02 | Stable identity and references | Lorebook IDs stable post-import; world/scenario typed `{kind,id}` refs in store; chat history IDs are currently index-derived; `ChatMeta` has IDs; `Lorebooks []interface{}` remains untyped | **Incomplete and duplicated.** Typed refs exist only in some categories; cross-category refs are not uniform; display-name fallbacks remain a risk. | All content, attachments, prompts, provenance, rename/delete, import/export |
| FND-03 | Model identity, capability discovery, lifecycle | `internal/modelintel`, GGUF metadata, image registry, projector pairing, evidence states | **Partial.** Identity is not yet a complete package/role model; filenames still act as hints; runtime capability evidence is incomplete; no single model-state authority. | Chat, image, audio, TTS, video, prompt templates, hardware, provenance, hot-swap |
| FND-04 | Engine abstraction and adapters | `llama.HTTPClient`, `sd.cpp` engine, audio pipeline/providers, TTS providers, video path | **Partial.** Each engine has a different lifecycle, cancellation, and status contract; no common operation envelope. | All model-backed features |
| FND-05 | Hardware detection and resource orchestration | `resource.Manager`, allocator, VRAM arbitration, audio hardware helpers | **Incomplete.** Advisory only; chat/image arbitration only; TTS outside; no reservation; no verified release. FUT-HW1 deferred. | Chat, image, audio, TTS, vision, video, future 3D |
| FND-06 | Concurrency, cancellation, operation coordination | Stream token exists (`beginStream`/`endStream`), image jobs, audio generations/leases, per-process worldstore lock | **Split.** Chat has one active stream slot; image and audio have richer per-job ownership; no application-wide operation ID/queue. Hot-swap does not coordinate with streams. | Chat, image, audio, video, model load, import, 3D |
| FND-07 | Generation provenance and reproducibility | Image metadata records prompt/model/options/seed; chat history does not record model/seed/template/effective params; `types.Message.Metadata` exists in a separate path | **Missing for chat; partial for image.** No durable per-turn identity, no model-removal behavior, no seed mode. | Trust/restore, hot-swap, per-chat model, debugging, audit, export |
| FND-08 | Media and attachment ownership | Gallery images + sidecars; video output has no catalogue; chat images are transient; image sources can point to world/scenario/character | **Incomplete/inconsistent.** No stable attachment identity; no unavailable-reference semantics; video lacks F3-F/F3-G contract. | Chat, image, video, characters, worlds, scenarios, avatars, 3D |
| FND-09 | Memory and context isolation | Lore resolver, memory service, multi-character scaffolding, per-character lore blocks | **Incomplete.** No application-owned inventory/clothing state; no proven cross-character isolation/flush; RAG/live prompt split; context truncation is character-count based. | Chat, multi-character, memory, lore, prompt, continuation |
| FND-10 | Prompt construction and model-specific templates | `app-prompt.go` system/context builder, `internal/backend/prompt/builder.go`, auto-template/chat policy, family defaults | **Duplicate/partial.** Live `App` builder and `chat.Service` prompt path differ; no complete Jinja/template matrix; token counting is approximate; prompt/anchor injection lacks a unified evidence model. | Chat, characters, worlds, lore, memory, RAG, vision, cues |
| FND-11 | Error handling, recovery, user-visible failure | `internal/errors`, loader diagnostics, structured operation results in newer stores, gallery partial-result flag | **Partial.** Some paths report structured failures; others log or ignore. No single recovery contract across engines. | Every subsystem |
| FND-12 | Settings and configuration ownership | Unified settings, image settings, sampler/model/audio sections, defaults | **Mixed atomicity and consumption.** Audio restores without activation by design; not every persisted control is consumed; no complete restart round-trip. | UI, model, image, audio, lifecycle, resource |
| FND-13 | UI/backend communication and state projection | All `App` methods bound; named Wails events; frontend keeps own chat/context state | **Two-way projection.** No request/generation ID in all events; frontend and backend both own chat state; page-local dialogs/notifications. | Chat, image, audio, settings, errors |
| FND-14 | Import/export, migration, rollback | F3-B/F3-C explicit journaled import; `ExportContentItem`; library registry | **Partial.** Only some categories support canonical import/export; characters/chats/images/videos lack equivalent; legacy stores remain recovery sources. | Content, chat, media, wipe, backup |
| FND-15 | Security and local/offline boundaries | Local-first; gallery path containment added; model downloads/document imports are local paths; no telemetry found | **Incomplete for future network features.** No implemented per-connection disclosure, offline switch, connection history, or remote pairing boundary. | Model downloads, curated collection, remote access, future network features |
| FND-16 | Application lifecycle and process ownership | Wails `OnShutdown` registered; shutdown coordinator; job object for children; idempotent stages | **Improved but not fully verified.** Packaged shutdown, multi-instance, alias/cross-process locks, and every child type are not proven. | All engines, persistence, audio, media, future services |

### 5.1 Foundation priority

The broadest foundations are FND-16 (lifecycle), FND-01 (ownership), FND-02 (identity), FND-03 (model identity), FND-05 (resource), FND-06 (operation ownership), and FND-10 (prompt/context). Feature work that bypasses them creates another owner or another state projection.

---

## 6. Data and persistence ownership map

| Data category | Authoritative owner today | Other readers/copies | Identity / format | Status and gate |
|---|---|---|---|---|
| Chats | `AppData/Roaming/HOWLChat/chats/<id>/` via `chat.PersistenceManager` | Browser current-chat cache; `chat.Service` context/snapshot; native test seam; transcript copies in frontend | Chat ID + folder; `meta.json`, history, context files | **Split ownership.** ROOT-01 open; F3-E future |
| Memory | `StateRoot/Chats` via memory service; also `memories.json` in chat persistence | `chat.Service` `memories_snapshot.json` path; frontend memory store | Chat ID + memory IDs | **Split/contradictory.** ROOT-01, IDENTITY-03, STARTUP-05; no production snapshot writer found |
| Candidate/continuity | Latent `cwd/Chats` candidate/continuity roots | No production caller found | Unknown/latent | **REACH-04 open.** Do not treat as supported |
| Characters | Browser `howl-chat-characters` + IndexedDB images; native `Characters/<name>/character.json` fallback | F2-B library JSON copies (non-authoritative); chat snapshots; external card discovery | Local character ID / display name folder | **F3-D future**; no canonical cutover |
| Worlds | Browser `howl-worlds`; canonical `worldstore` world envelopes for F3-C path; legacy `Worlds/<name>` folders | Chat `world.json` snapshots; F2-B copies; editor browser records | `wld_` / preserved legacy ID; typed refs | **F3-C store foundation; acceptance not granted.** Browser editors still write import source |
| Scenarios | Browser `howl-chat-scenarios`; canonical `worldstore` scenario envelopes; legacy `Scenarios/<name>` folders | Chat `scenario.json` snapshots; F2-B copies; editor records | `scn_` / preserved legacy ID | **F3-C store foundation; acceptance not granted** |
| Lorebooks | `internal/lorebookstore` under registered roots | Rebuildable browser cache; per-chat `lorebooks.json` snapshots | `lb_` / preserved legacy ID; `howl-lorebook-v1` per entry | **F3-B accepted with limitations** |
| Images | Gallery root + PNG + JSON sidecar; image library/config; `ImageMetadata` includes full options | Chat `hasImage` flag only; character/world/scenario image refs; F2-B copies | Absolute path + sidecar filename; no stable media catalogue identity | **F3-F future**; T0.3/T0.4 unresolved; atomicity claims need runtime proof |
| Videos | `EngineConfig.OutputDir` direct output + base64 result to UI | No gallery sidecar/catalogue/history owner | Output filename only | **F3-G future**; no delete/export contract |
| Settings | Unified settings directory; image settings store; legacy sampler/model files | Frontend settings state; per-page controls | Sectioned JSON; mixed atomicity | **Partial**; no full restart round-trip here |
| Model registry/cache | Settings directory; LLM/image registries; `modelintel` cache | Frontend dropdowns; image model picker | Path + kind + cached metadata | **Source-confirmed**; filename/capability evidence incomplete |
| Audio configuration | Audio settings in unified store; manager desired/applied generations | Frontend audio state; per-character voice config | Desired vs applied vs capabilities | **Contracted/gated**; TTS governance unresolved |
| Import journals/markers | Settings directory / hidden staging under roots | Import UI; rollback | Operation ID, plan hash, journal | **Implemented for F3-B/F3-C scopes** |
| Wipe/cleanup | `WipeAllData` stages chat root and gallery root; registry-aware inventory | F3-C surfaces not all covered | Root paths | **WIPE-02 open**; settings reset separate |
| Avatars | Planned `Documents\Howl Chat\Avatars\<character>\` `.glb` + maps | Image pipeline; future WebView 3D | Character ID + generated asset | **Planned only** |

### 6.1 Ownership conclusions

1. **F3-B is the only accepted canonical content owner.** It does not make chats, characters, worlds, scenarios, images, or videos canonical.
2. **F3-C is a store foundation, not a completed ownership cutover.** The phase gate, native evidence, consumer wiring, and six defects remain open.
3. **Chat data is still multi-root.** A chat can have AppData history, StateRoot memory, latent `cwd` continuity, native context, and browser projections. The plan must not treat “chat” as one file.
4. **Media has no product-wide identity.** Images have a sidecar path identity; videos have almost none; chat images have only a boolean.
5. **Wipe and settings reset are different contracts.** Neither should be expanded by inference from a UI label.

---

## 7. Model and hardware dependency map

### 7.1 Current model families and formats

| Modality | Current formats / engines | Current evidence | Future requirement |
|---|---|---|---|
| Chat LLM | GGUF/`.bin` through `llama-server`; projector (`mmproj`) pairing; chat template/family defaults | Source-confirmed; model intelligence/registry; readiness monitor; no runtime verification here | Additional formats, ternary (`T3.6`), architecture/era options (`FUT-C4`), prompt/template intelligence (`FUT-PROMPT1`) |
| Vision | llama.cpp multimodal projector; one current-turn image | Source-confirmed; capability gate uses projector metadata | Multi-image (`E.2`), image anchoring (`E.1`), durable image cues (`FUT-C2`) |
| Image | sd.cpp checkpoints/GGUF/safetensors; modular VAE/CLIP/T5/LLM components; LoRA; ControlNet; upscale; inpaint; workflows | Source-confirmed; T0.3/T0.4 open; runtime unverified | Hardware placement (`FUT-IMG1/2/3`), flash attention/offload, context menu/workflow improvements, editing |
| Video | SVD/WAN/LTX/Hunyuan-class families; frame/FPS/motion/loop options | Source-confirmed; no gallery catalogue | F3-G media contract; video generation requirement still needs formal scope |
| ASR | Whisper, Qwen ASR, llama adapter, fallback; optional remote engines | Source-confirmed; remote availability explicitly unverified | Additional audio, streaming, hardware placement |
| TTS | Piper, Kokoro, MeloTTS, helper/cache; per-character voices | Source-confirmed; governance unresolved | F2-B7 says no new TTS feature/contract change in F2; additional audio deferred |
| Audio generation | None found | **Planned only** | T3.8 |
| 3D | None in production; UltraTex external ComfyUI; prebuilt MakeHuman bust/MakeHuman-derived mesh | **Planned only** | T4.1/T4.2; Tier B app-side projected portrait; Tier C external; Tier D TripoSR-class |

### 7.2 Model identity and lifecycle states

The product needs to distinguish, per model component:

- exact model/variant, architecture, container/format, quantisation;
- component role: backbone/thinker, vision projector, audio projector, encoder, decoder, talker, vocoder, tokenizer, processor, text encoder, UNet, VAE;
- capability enabled and capabilities not enabled;
- evidence state: Verified, Source-derived, Inferred, Unknown, Conflicting;
- runtime state: requested, loading, ready, failed, unloading, unloaded;
- dependency state: Satisfied, Missing, Incompatible, Ambiguous, Optional, Unknown.

Current source has partial versions of these ideas (model intelligence evidence labels, projector pairing, loader stages), but not one complete package/role resolver. `docs/deferred-model-ecosystem.md` is the authoritative deferred statement of that requirement.

### 7.3 Hardware/resource dependency findings

- `resource.Manager` stores GPU/RAM estimates and loaded model records, but it is **advisory**.
- Chat/image arbitration exists through `CheckModelFit` / `FreeVRAMFor`; the frontend can continue when `can_fit` is false if no blocking entry is returned.
- Image auto-load does not reliably register; chat registration occurs before asynchronous readiness; unload can unregister all LLM entries rather than the specific model.
- TTS is not an admission participant; TTS resource type exists but no arbitration kind.
- Hot-swap stops the existing llama server without checking an active stream; after stop, spawn/startup failure has no rollback to the old model.
- No evidence that unload actually removed weights/child resources; `FUT-HW1` AC-HW1-4/5/7 are unmet.
- The compatibility matrix required by FUT-HW1 does not exist. No hardware combination may be called “supported” without measured evidence.

### 7.4 Future compatibility analysis

| Future capability | Required existing owners | Current interface fit | Architectural change if frozen without it | Can remain deferred? |
|---|---|---|---|---|
| Model hot-swap during active conversation | Model lifecycle, resource admission, stream ownership, chat state | **Poor.** Manual load exists; no stream coordination/rollback | Chat schema needs continuation/provenance; runtime needs state machine; UI needs requested/active/saved distinction | Yes, but the **chat schema decision cannot ignore it** |
| Per-chat continuation model | Chat metadata, model identity, settings | **Absent.** Global last-used only | Add explicit per-chat preference separate from provenance and runtime; migration for existing chats | Yes, but schema decision now |
| Per-turn provenance | Chat history, model identity, engine config, seed | **Absent in persisted history**; separate type exists unused | Add durable provenance object; all assistant paths must use one helper; legacy must be unknown not fabricated | Yes, but schema decision now |
| Message/image cues | Chat schema, prompt assembly, tokenizer, model capability | **Absent.** One transient image | Persist cues and attachment refs; token budget; capability gating | Yes |
| Memory inventory/clothing | Structured state, scope-aware retrieval, prompt selection, per-character isolation | **Absent.** Current memory is prose records; lorebook is not inventory | Add typed state and selection policy; do not encode inventory as untyped lore | Yes, but prompt/memory architecture must allow it |
| 3D avatar | Media ownership, image pipeline, TTS amplitude, WebView performance, external ComfyUI | **Absent**; 2D avatars exist | Must not make core chat depend on 3D; store `.glb` as media with stable character identity; external failure must fall back | Yes; Tier B can be optional |
| Additional model formats/ternary | Model identity, loader, hardware | **Partial.** GGUF/chat and image formats exist; ternary not established | Do not hardcode “GGUF chat” into model identity; capability/role model required | Yes |
| llama.cpp features/multi-GPU/offload | Hardware probe, model capability, loader | **Settings exist**, but verified support matrix does not | Admission must be backend-driven; never expose a control the backend ignores | Yes |
| AI-assisted character/lore/world creation | Prompt templates, model capability, structured schemas, review UX | **Partial** (lorebook/character forms; no field generator) | Needs grammar-constrained output and per-field draft/accept; do not bind to one model family | Yes |
| Android remote/QR | Server/transport, auth, session, network disclosure, resource coordination | **Absent/unestablished** | Would introduce a remote boundary and security model not present in the local-only shell | Yes; requirement must be established first |
| Creative studio/asset generation | Media identity, jobs, provenance, resource orchestration | **Absent/unestablished** | Needs a catalog and job model before any UI | Yes; requirement must be established first |
| Video workflows | Media catalogue, sidecar, gallery, resource | **Partial** generation, no catalogue | F3-G must precede a real video feature; otherwise clips are orphaned | Yes, but media schema decision now |

**Expensive rework if deferred without architectural allowance:** chat schema, root/path identity, media/attachment identity, model identity/readiness, prompt/template ownership, and one operation/cancellation model. These are the decisions that become migrations later.

---

## 8. Dependency graph and implementation sequence

### 8.1 Logical dependency graph

```text
FND-16 Lifecycle / process ownership
  ├── FND-05 Resource admission / release
  ├── FND-06 Operation / cancellation ownership
  └── all model-backed engines

FND-01 Storage ownership
  ├── FND-02 Stable identity / references
  ├── FND-08 Media / attachments
  ├── FND-09 Memory / context isolation
  ├── FND-14 Import/export/migration
  └── wipe / recovery

FND-03 Model identity / capability
  ├── FND-10 Prompt / template / tokenizer
  ├── FND-05 Resource planning
  ├── FND-07 Provenance / reproducibility
  └── hot-swap / continuation

FND-04 Engine adapters
  ├── chat, image, ASR, TTS, video, future 3D
  └── FND-11 Error/recovery contract

FND-13 UI/backend events
  ├── chat streaming and cancellation
  ├── image/audio job status
  └── settings/error projection
```

### 8.2 Dependency register

| Order | Component | Direct dependencies | Why it must precede dependents |
|---:|---|---|---|
| 1 | F3-C Phase 0 disposition and the six defects | Frozen contract, owner decisions | Every chat/root/wipe/identity claim downstream is unreliable while the target set and ownership boundary are undecided |
| 2 | Canonical root/path identity | F3-C D-ROOT; ALIAS-06 | A path string is not identity; locks, wipe, import, and duplicate resolution depend on physical identity |
| 3 | Durable chat/attachment/model schema decision | F3-C amendment, FND-02, FND-08, FND-07 | Chat serialization is expensive to migrate after F3-E; provenance/continuation/seed/attachment gaps are already documented |
| 4 | F3-D characters | F3-C acceptance; identity; media refs | Characters are needed by chat, multi-character, avatars, and import |
| 5 | F3-E chats | F3-D; root identity; chat schema | Chat is the consumer of all context, memory, model, and media identity |
| 6 | F3-F images / F3-G videos | F3-E; media identity; gallery | Media identity and sidecars must exist before attachments, avatars, and video catalogues |
| 7 | F3-H wipe/cleanup | F3-E/F/G; surface inventory | Destructive lifecycle must enumerate every retained surface; WIPE-02 is already open |
| 8 | FND-16 lifecycle verification | Current lifecycle source; packaged runtime | All engines and final saves depend on a reliable shutdown/ownership boundary |
| 9 | FND-03 model identity/capability + FND-10 prompt | GGUF/metadata; tokenizer; template evidence | Prompt/template and hardware planning cannot be correct if model identity is filename-only |
| 10 | FND-05 resource admission + FND-06 operation ownership | Model identity; lifecycle; engine adapters | Hot-swap, image/audio coexistence, and multi-character concurrency need one admission/operation authority |
| 11 | FND-07 provenance | Chat schema; model identity; engine config | Provenance must capture the effective request at attempt start, not reconstruct it later |
| 12 | FND-09 memory/context isolation | Chat schema; prompt; lore; model capability | Multi-character and inventory/clothing need typed, scope-aware state and prompt selection |
| 13 | Image/audio/video completion | F3-F/G; resource; cancellation; error contract | These are model-backed workloads with real VRAM and lifecycle risk |
| 14 | UI/prompt aids, AI-assisted creation | Prompt/model capability; structured schemas; review UX | Do not build generators without a durable field/schema and review/accept boundary |
| 15 | 3D avatar | Media, image, TTS, resource, WebView performance | It is a consumer of several foundations and must degrade to 2D/image fallback |
| 16 | Remote/QR and creative studio | Explicit requirements, security boundary, media/job model | No authoritative requirement exists; building architecture for them now would be speculation |

### 8.3 Planning sequence, preserving the existing F3 order

This is a dependency recommendation, **not authorization** and **not a reordering of the accepted F3 sequence**.

1. **Preserve the current gate position.** Do not start F3-D, R6, implementation, migration, or cleanup. Resolve the F3-C Phase 0 amendment question and the six defect dispositions first.
2. **Decide the chat target-set boundary.** If per-turn provenance, seed state, attachment references, and continuation preference are required product behaviors, decide them before F3-E serializes chat. If they remain out of scope, record the later migration cost explicitly.
3. **Close D-ROOT and its prerequisites.** Root ownership, wipe surface, identity clear, latent candidate roots, startup degradation, and alias identity. No destructive lifecycle or cross-process lock claim before these.
4. **Complete F3-D → F3-E → F3-F → F3-G → F3-H in the existing order.** F3-D is blocked until F3-C acceptance; each stage must preserve REG-1.
5. **Harden the shared runtime.** Lifecycle verification, model identity/capability, resource admission, operation/cancellation IDs, and provenance. These precede hot-swap, multi-character concurrency, and reliable image/audio coexistence.
6. **Complete context/memory/prompt foundations.** One prompt path, tokenizer/context accounting, model-family templates, multi-character isolation, structured memory/inventory, cue/attachment selection.
7. **Finish media and creative workflows.** Image editing/workflows, video catalogue, attachments, then avatar/3D and optional creative studio if separately authorized.
8. **Audio/TTS only under existing governance.** TTS contract remains frozen; no audio generation or contract change by implication.
9. **Deferred intelligence/UI work.** FUT-HW1/FUT-HW2, FUT-PROMPT1, FUT-C*/ST*/CH*/IMG*, lore/memory intelligence, and Tier 3 features remain deferred until their own gates.

### 8.4 Circular and conflicting dependencies

- **Resource cycle:** model loading needs admission; admission needs to know what is loaded; current loading registers before readiness. The cycle is broken only by a verified state machine, not another map.
- **Chat/prompt cycle:** `App` owns live chat state; `chat.Service` owns canonical context; persistence owns history; frontend owns a projection. One path must be declared authoritative before F3-E.
- **Lore cycle:** canonical lore store → cache → chat snapshot → prompt. F3-B defines the canonical store, but the prompt path still has copies and sync timing.
- **Media cycle:** image generation produces media; gallery owns it; chat attaches it; character/world/scenario reference it; F3-F/F3-G must define the identity before attachments and avatars depend on it.
- **3D cycle:** avatar needs image generation and TTS amplitude; both are model-backed and resource-constrained; avatar must therefore be optional and fallback-safe.

---

## 9. Release versus deferred scope

### 9.1 Existing release requirements

These are current product/release obligations, independent of future additions:

- Tier 0 release blockers T0.1–T0.4.
- Tier 1 core mechanics T1.1–T1.5, subject to current evidence and the F3-C contradictions.
- Tier 2 UI/readability T2.1–T2.5.
- F3-A/F3-A-R1 safety repairs.
- F3-B lorebook ownership, accepted with documented limitations.
- F3-C world/scenario ownership: Phase 0 open, six defects open, acceptance not granted.
- F3-D through F3-H, in order, only after their gates.
- Current chat, character, world/scenario, lorebook, image, audio, settings, and lifecycle functionality to the extent independently verified.

### 9.2 Shared foundations required for the intended product

These are not optional features; they are the architecture future features depend on:

- FND-16 lifecycle/process ownership and verified shutdown.
- FND-01 storage/content ownership and registry-aware wipe.
- FND-02 stable identity and typed references.
- FND-03 model identity, role, capability, and readiness.
- FND-05 hardware/resource admission and verified release.
- FND-06 operation/cancellation ownership.
- FND-07 generation provenance and reproducibility.
- FND-08 media/attachment identity.
- FND-09 memory/context isolation.
- FND-10 prompt/template/tokenizer ownership.
- FND-11 error/recovery contract.
- FND-12 settings ownership.
- FND-13 UI/backend state projection.
- FND-14 import/export/migration.
- FND-15 security/local/offline boundary.

### 9.3 Deferred feature implementation

- T3.1–T3.12, except where a current release requirement independently requires part of the behavior.
- FUT-C1–C5, FUT-ST1–ST5, FUT-CH1–CH3, FUT-IMG1–IMG3.
- FUT-CHUI1/2, FUT-CHTOK1, FUT-CHID1, FUT-CHAIU1, FUT-MEM1, FUT-LORE1.
- FUT-HW1, FUT-HW2, FUT-PROMPT1.
- T4.1/T4.2 3D/multi-character avatars.
- Additional audio generation/music.
- Video workflow completion beyond the current source path.
- AI-assisted character/lore/world creation beyond currently implemented forms.
- Import/export for categories not yet covered by F3.

### 9.4 Optional future capabilities

- Reactive 3D avatar Tier C/D and generated mesh (TripoSR-class).
- Desktop character avatar as a persistent interactive surface.
- Android remote access and QR pairing (unestablished requirement).
- Creative studio/general asset generation (unestablished requirement).
- AI-as-user/persona takeover (separate feasibility study).
- Curtain vault/private boards.
- Multi-user, cloud sync, telemetry, or online services — none is part of the current local-first architecture and none should be inferred.

### 9.5 Scope protection rule

**A future capability may influence architecture without authorizing implementation.** The plan may reserve a field, an identity, a state, or a boundary; it may not claim the feature, add a migration, or alter a release gate merely because the field exists.

---

## 10. Accepted decisions and protected boundaries

| Decision / boundary | Status | Planning consequence |
|---|---|---|
| F2-C1: one authoritative owner per content item; libraries are address/permission registrations | **DESIGN RECORDED** | Never create a second canonical copy; F3-B/F3-C adapters must preserve this |
| F2-C2: F2-B JSON copies are untrusted export/compatibility snapshots | **DESIGN RECORDED** | Do not read them back automatically |
| F2-C3: browser storage becomes a rebuildable cache after cutover | **DESIGN RECORDED** | Do not make localStorage canonical again |
| F2-C4: chats retain multi-file folder format; register existing root in place | **DESIGN RECORDED** | Chat is a folder, not a single JSON document |
| F2-C5: images/videos are binary + sidecar with stable identity and atomic commit | **DESIGN RECORDED** | F3-F/F3-G must define the media catalogue before attachments/avatars |
| F2-C6: Save is create-only; replacement requires Write/explicit collision choice | **DESIGN RECORDED** | Do not permit silent overwrite |
| F2-C7: no automatic migration/background sync; explicit, verified, journaled, reversible | **DESIGN RECORDED** | No “helpful” migration to solve a root defect |
| F2-C8: disconnected roots retain registrations and unavailable inventory | **DESIGN RECORDED** | Missing is not deleted; no substitution |
| F3-A1–A8 | **IMPLEMENTED / deferred as recorded** | Export semantics, collision safety, gallery containment, wipe staging, no legacy read-back |
| F3-A-R1-1–8 | **IMPLEMENTED / one unresolved audio flake** | Preserve unsafe-target rejection, rollback, identity-based export, staged wipe, gallery containment |
| F3-B1–B11 | **IMPLEMENTED; F3-B accepted with limitations** | Lorebooks canonical; cache rebuildable; REG-1 preserved |
| F3-C R1–R3 and frozen contracts | **FROZEN; implementation incomplete; acceptance not granted** | World/scenario envelopes, typed refs, revisions, import/rollback, snapshot staleness |
| F3-C Phase 0 | **OPEN** | No F3-D; no implementation authorization |
| Six D-ROOT defects | **OPEN** | No bulk closure; no canonical-root migration by implication |
| D-STORAGE | **BLOCKED** | No storage migration/ownership cutover |
| V20 | **BLOCKED pending F3-F** | No premature media acceptance |
| R6 | **UNAUTHORIZED** | Do not author or imply it |
| REG-1 | **PROTECTED REGRESSION CONSTRAINT** | Worlds/scenarios/characters/lorebooks remain reusable, independently usable, scoped, and hot-swappable; no future ownership change may weaken it |
| TTS contract / F2-B7 | **FROZEN / NO NEW TTS FEATURE IN F2** | Audio architecture exists, but do not change the TTS contract by inference |
| Gate v2 V1–V17 | **ADOPTED** | Preserve; V18–V24 remain proposed |
| A9/A10/A12/A13 | **CLOSED** | Do not reopen without concrete contradictory evidence |
| Local-first product boundary | **Current product principle** | No cloud inference, telemetry, hidden network activity, or automatic upload |

**Protected artifact preservation:** this plan does not modify the F3-C freeze, master contract, adjudication package, owner-decision request, inventory, investigation, defect disposition request, or any accepted decision record.

---

## 11. Contradiction and unresolved-decision register

| ID | Contradiction / unresolved question | Evidence | Planning impact | Disposition |
|---|---|---|---|---|
| C-01 | `STAGE_STATUS` closes T1.1–T1.4, while F3-C still finds split chat/memory roots and open ownership defects | `STAGE_STATUS.md`; F3-C defect inventory | “Closed” stage status is not product-wide verification; treat state ownership as partial | **Unresolved / do not reopen accepted work; add current evidence** |
| C-02 | `chat.Service` is a canonical consumer, but the live chat stream builds prompts directly from `App` and does not demonstrate RAG use | `app-prompt.go:buildPromptMessages`; `app.go:SetRAGService`; `chat.Service` | Two prompt paths; RAG/lore can be configured without affecting the live reply | **Unresolved architecture split** |
| C-03 | F3-C freeze is strong for text/context but cannot preserve per-turn provenance, explicit seed state, attachment refs, or per-chat continuation | F3-C Phase-0 assessment | Chat schema is incomplete for intended product; later migration if not decided | **Contract gap; amendment proposed, not adopted** |
| C-04 | “Architecture freeze granted” can be misread as “F3-C accepted” | F3-C records | Freeze is not acceptance, verification, or implementation authorization | **Clarified; gates preserved** |
| C-05 | F3-B is accepted, but REG-1 hot-swap/activation and model inference were not all exercised | F3-B acceptance limitations | Accepted ownership does not prove the creative workflow | **Bounded limitation; not a reopening** |
| C-06 | F2-C says chats retain a multi-file folder and register AppData in place; current memory/candidate roots add more surfaces | F2-C §5.5; F3-C root inventory | “Chat root” is not one root; wipe and identity must enumerate surfaces | **ROOT-01/WIPE-02 open** |
| C-07 | Hot-swap can stop the llama server while a stream is active, and has no rollback after stop | `app-model.go`; F3-C Phase-0 assessment | Model switching is not safe as a product behavior; per-chat continuation needs readiness/failure rules | **Unresolved / proposed only** |
| C-08 | Model selection is global (`howl-last-model-path`) while the product needs per-chat continuation | `chat-model.js`; `ChatMeta` has no model field | Restoring a chat cannot restore its last-used model; hot-swap/provenance retrofit | **Proposed amendment** |
| C-09 | Chat history stores `hasImage` but no durable image reference; current image is transient base64 | `ChatHistory`; `chat-input.js`; Phase-0 assessment | Reload cannot restore the visible image; `hasImage` is not ownership | **Proposed amendment / F3-F dependency** |
| C-10 | Chat history IDs are index-derived while stable identity is a stated requirement | Phase-0 assessment; `app-chat.go`; `app-prompt.go` | Editing/regeneration/reload can collide or reorder identity | **Implementation/contract gap** |
| C-11 | `ChatData.Lorebooks` remains `[]interface{}` while F3-C introduces typed refs | `internal/backend/chat/types.go`; F3-C frozen contracts | Cross-category typed references are inconsistent; later migration risk | **Unresolved schema debt** |
| C-12 | Resource manager exists, but registration is not readiness, image auto-load is inconsistent, TTS is outside admission, and release is unverified | `app-model.go`; `app-image-model.go`; `app-vram-arbitration.go`; assessment | “Model loaded” and “fits” are not authoritative; no admission boundary for future hot-swap | **Unresolved foundation** |
| C-13 | `WipeAllData` now correctly stages chat and gallery roots, but F3-C surfaces (memory, candidate/continuity, native context, active identity) are not all enumerated | `app-wipe.go`; WIPE-02 | Wipe can report success while retained surfaces remain | **WIPE-02 open** |
| C-14 | Root authority treats roots as equivalent search locations, but physical path identity/alias rejection is not solved | F3-C root authority; ALIAS-06 | Two paths can address one store; locks/duplicate resolution may be wrong | **ALIAS-06 open** |
| C-15 | Startup can read a corrupt `<StateRoot>/Chats/memories_snapshot.json`, nil `chatService`, and skip RAG/content-library wiring | STARTUP-05 | The app can start degraded without an explicit policy | **STARTUP-05 open** |
| C-16 | Tier 3 is “only once the app is operational”, but several Tier 3 behaviors are already partially present in code/UI | `PROJECT_ISSUES...`; source | Source presence must not be read as completion; requirements remain deferred | **Preserve original ordering** |
| C-17 | T4 avatar assumes TTS amplitude, but TTS lifecycle/governance is unresolved | `docs/tier4-avatar.md`; TTS records | Avatar mouth animation has an unresolved dependency; 3D must be optional/fallback | **Deferred dependency** |
| C-18 | The mission names Android/QR, desktop avatar, 3D mesh, and creative studio, but only the 3D avatar/partial mesh is documented | Ledger/avatar note vs source search | Do not invent requirements or architecture for unestablished features | **Unestablished** |
| C-19 | F2-C reported non-atomic gallery image/sidecar and suppressed delete errors; F3-A6/A7 claim atomic image+sidecar and containment | F2-C; F3-A records; current `gallery` source | The later implementation may have fixed it, but runtime/edge-case proof is not established here | **Verify before relying; do not use old finding as current fact** |
| C-20 | `OnShutdown` was unregistered in the 2026-11-09 assessment; current `main.go` registers it | `main.go`; assessment | The old lifecycle finding is superseded; do not repeat it as current | **Superseded** |
| C-21 | Gate v2 V18–V24 are proposed, not adopted, while F3-C Phase 0 already depends on provenance/continuation questions | Gate v2; Phase-0 assessment | Do not cite proposed gate rows as accepted authority | **Proposed only** |
| C-22 | Story tools page exists but committed scope is unknown | `story-tools*.js`; assessment | Do not infer features from controls or tutorials | **Unknown / owner input needed** |

### 11.1 Decisions genuinely required before architecture is final

1. **Chat durable target set:** adopt, reject, or defer per-turn provenance, seed mode, attachment references, and per-chat continuation before F3-E serialization.
2. **Canonical root and path identity:** one owner per surface or explicit per-surface ownership; physical identity/alias policy.
3. **Wipe semantics:** what “Wipe All Data” covers; what is preserved; how every F3-C surface is enumerated.
4. **Startup degradation:** fail-closed, explicit degraded mode, or removal for corrupt/unreadable root state.
5. **One live prompt path:** `App` vs `chat.Service`; where RAG, lore sync, context accounting, and character isolation live.
6. **Resource authority:** one model state machine and admission boundary; no second manager.
7. **Media identity:** stable attachment/image/video identity and unavailable-reference behavior.
8. **TTS governance:** unchanged; no feature decision by implication.
9. **Story tools scope:** owner must establish it.
10. **Unestablished future features:** Android/QR and creative studio require explicit requirements before architecture work.

---

## 12. Evidence gaps and unknowns

| Gap | Why it matters | Current evidence |
|---|---|---|
| No runtime in this mission | Source presence cannot establish working behavior | Static only |
| No packaged/native Windows run | WebView2, path layout, file identity, shutdown, and permissions can differ | F3-C records explicitly say native acceptance/runtime is incomplete |
| No model load or inference | Chat, RAG, image, audio, video, hot-swap, and 3D behavior unverified | No model loaded |
| No concurrency/race test | Split state and one-stream cancellation may fail under overlap | Source-level only |
| No filesystem alias matrix | Junctions, symlinks, 8.3, case, mapped drives, UNC can defeat root identity | ALIAS-06 open |
| No two-process lock test | Cross-process guarantee cannot be claimed | ALIAS-06 / D-STORAGE |
| No chat reload with attachment | `hasImage` vs durable reference gap | Phase-0 assessment |
| No per-turn provenance/seed verification | No historical model identity exists in persisted history | Phase-0 assessment |
| No hot-swap stream/rollback test | Product behavior and safety unknown | Source/assessment |
| No image OOM/VRAM placement matrix | FUT-IMG1/2/3 and FUT-HW1 unverified | Resource source only |
| No audio/TTS end-to-end acceptance | TTS governance unresolved; remote engines unverified | Audio records |
| No image/video gallery sidecar runtime acceptance | F3-F/F3-G not complete | F2-C + current source |
| No 3D performance measurement | 30fps/viewport/token-speed assumptions are third-party | `docs/tier4-avatar.md` explicitly unverified |
| No Android/QR or creative-studio requirement | Cannot plan or implement responsibly | Not found |
| No exhaustive line-by-line audit of all 240+ Markdown docs | Some prose requirements may exist outside the authoritative ledger | Named corpus recovered |

**Unknown is not a negative result.** It means the evidence does not exist yet.

---

## 13. Recommended next investigations

These are planning/read-only next steps. They are not implementation authorizations.

1. **Adjudicate the F3-C chat target-set gap.** Decide whether the four product behaviors (provenance, seed mode, attachment refs, continuation) are in F3-C, a controlled amendment, or an explicitly accepted later migration. This must precede treating F3-E as final.
2. **Disposition the six defects individually.** ROOT-01, WIPE-02, IDENTITY-03, REACH-04, STARTUP-05, ALIAS-06. No bulk closure.
3. **Define physical root identity and alias behavior.** Source-level path equality is not enough; specify the Windows evidence required for junctions, symlinks, case, 8.3, mapped drives, UNC, and changed cwd.
4. **Reconcile the two chat/prompt paths.** Establish which path owns prompt assembly, RAG, lore sync, token accounting, and character isolation; do not create a third path.
5. **Inventory model identity/capability evidence.** Turn filename/projector hints into a role-aware package model with Verified/Source-derived/Inferred/Unknown/Conflicting states.
6. **Define the resource state machine and evidence.** Requested → loading → ready → failed → unloading → unloaded; verify release; no “loaded” from a request alone.
7. **Define media/attachment identity.** Image/video sidecar, stable ID, missing-reference semantics, chat attachment reference, and F3-F/F3-G boundary.
8. **Recover the unestablished requirements from the owner.** Android remote/QR, creative studio/asset generation, and any desktop-avatar interaction scope. Do not infer them from the mission's discovery list.
9. **Establish story-tools scope and TTS governance.** Both are unresolved and should not be filled by UI inference.
10. **Reconcile old stage closures with current evidence.** Treat `STAGE_STATUS` closures as historical unless current source/runtime evidence supports them; do not reopen accepted work without concrete contradiction.

---

## 14. Source-reference register

| Short reference | Path | Used for |
|---|---|---|
| Ledger | `PROJECT_ISSUES_AND_EXPANSION.md` | Tiers 0–4, E.1–E.8, F2-B/C, F3-A/B/C, deferred register, HW1/HW2, PROMPT1 |
| Avatar note | `docs/tier4-avatar.md` | T4.1/T4.2, body/skin/emotion pipeline, tiers, licensing, unknowns |
| Model ecosystem note | `docs/deferred-model-ecosystem.md` | Role-aware model identity, dependency resolution, offline/network disclosure |
| App assessment | `docs/review/ASSESSMENT-howl-application-reality-check-shared-foundations-2026-11-09.md` | Subsystem baseline, risks, foundation order |
| F2-C architecture | `docs/review/ARCHITECTURE-F2C-CONTENT-OWNERSHIP-CUTOVER-2026-11-09.md` | Ownership map, dual storage, media, cutover sequence |
| F3-C master contract | `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` | Controlling F3-C authority |
| F3-C Phase 0 | `docs/review/ASSESSMENT-F3C-PHASE-0-FINAL-REQUIREMENTS-IMPACT-2026-11-11.md` | Provenance/seed/attachment/continuation gaps; hot-swap adjacency |
| F3-C defect request | `docs/review/DISPOSITION-REQUEST-F3C-DEFECTS-2026-11-11.md` | Six open defects; gate effects |
| F3-C Gate v2 | `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md` | Verification status and authority rules |
| F3-B evidence | `PROJECT_ISSUES...` F3-B section; `docs/review/ENGINEERING-F3B-LOREBOOK-OWNERSHIP-2026-11-10.md` | Accepted lorebook ownership and limitations |
| Status | `STAGE_STATUS.md` | Historical stage closures and deferred requirements |
| Debt | `TECH_DEBT.md` | Known open control/state issues |
| Lorebook checklist | `CHECKLIST.md` | Lorebook features/tests claimed |
| Audio architecture | `docs/audio-architecture/*` | ASR/TTS/model loading/pipeline design |
| Source | `app-*.go`, `internal/**`, `frontend/dist/**` | Presence/absence and current ownership evidence |

---

## 15. Final self-review

- [x] The original additions list was recovered from `PROJECT_ISSUES_AND_EXPANSION.md` and the separate avatar/model-ecosystem notes, not reconstructed from memory.
- [x] Every added requirement above carries a source section or is explicitly marked **Unestablished**.
- [x] Subsystem coverage includes application lifecycle, chat, characters, worlds, scenarios, lorebooks, memory, RAG, image, video, ASR, TTS, audio generation, vision, model management, hardware/resource, content libraries, settings, UI, story tools, errors, 3D/avatar, remote/QR, and creative studio.
- [x] Shared foundations are named once and cross-referenced rather than duplicated.
- [x] Deferred features remain deferred; no requirement is promoted by inclusion.
- [x] F3-B accepted work, F3-C open gates, six defects, A9/A10/A12/A13, D-STORAGE, V20, R6, and Gate v2 status are preserved.
- [x] Hot-swap, per-chat continuation, per-turn provenance, seed state, and attachment references are labeled **proposed/not adopted**, except where the ledger records an existing requirement.
- [x] Tables, IDs, and dependency links were checked for internal consistency.
- [x] This artifact is documentation-only. No source, test, build, migration, model, or user-data action was performed.

**Status of this artifact:** complete for planning review; not adopted as authority; independently reviewable.  
**Next authority step:** Muse review, then owner decisions on the unresolved items in §11.1. This document does not begin those decisions.

---

## 16. Preservation verification record

The following protected F3-C artifacts were hashed after this planning artifact was written. The values match the hashes recorded before this task, so none was modified by it.

| Protected artifact | SHA-256 |
|---|---|
| `ADJUDICATION-PACKAGE-F3C-OD18-INTEGRATION-2026-11-11.md` | `85e1968bac4eb70f489bbd7dc89057debd17e7b3bc2336aa69c036a4352f83e8` |
| `OWNER-DECISION-REQUEST-F3C-FOUNDATIONS-2026-11-11.md` | `ed8bf2e2326a2ec74838dc0c8cd4de3332b8807d2d0a2647a18cfda060bc9c7c` |
| `INVENTORY-F3C-ROOT-OWNERSHIP-2026-11-11.md` | `343104b528818df699d752c419770fdc1bfd7fe71bc1f249e0dc6847f308c6b8` |
| `INVESTIGATION-F3C-DROOT-CLOSURE-2026-11-11.md` | `494539b9dde27a78c9c4c091131a93dce99c1debd81d11fea911e6d411e2e6a6` |
| `DISPOSITION-REQUEST-F3C-DEFECTS-2026-11-11.md` | `ce32ac0bedbe216b55c065be073b2165209e0ebad4cb435897b829336d29d2c5` |

**Reference check:** every explicit `docs/review/*.md` path and every root document named in §14 exists.  
**Table check:** requirement, foundation, ownership, dependency, scope, contradiction, and evidence-gap tables use stable IDs and were checked for duplicate foundation names.  
**Placeholder check:** no unresolved placeholder marker remains.  
**Write boundary:** this task created only this planning artifact; no source, contract, freeze, acceptance, or decision record was edited. The artifact's own SHA-256 is reported at handoff because a document cannot contain its own stable hash.
