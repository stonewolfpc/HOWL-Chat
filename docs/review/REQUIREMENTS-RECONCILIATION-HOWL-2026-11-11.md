# HOWL CHAT — REQUIREMENTS RECONCILIATION AND EXISTING-CAPABILITY VERIFICATION

**Date:** 2026-11-11  
**Baseline authority:** `docs/review/REQUIREMENTS-RECOVERY-OWNER-ADDENDUM-HOWL-2026-11-11.md`  
**Baseline SHA-256:** `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056`  
**Status:** **RECONCILIATION DELIVERABLE — NOT THE CORRECTED MASTER PLAN, NOT AN IMPLEMENTATION AUTHORIZATION, NOT A CONTRACT, NOT A FREEZE, NOT AN ACCEPTANCE, NOT R6, NOT AN F3-C AMENDMENT, NOT F3-D**  
**Implementation authorization:** NONE  
**Runtime / benchmark authority:** NONE

This document reconciles the owner-characterized recovery baseline against existing architecture and implementation evidence. It does not replace the baseline, adopt a corrected plan, authorize a capability, or claim that a requirement is operational because code or a control exists.

---

## 1. Purpose and authority chain

The governing distinction is:

| Layer | Responsibility | Does not do |
|---|---|---|
| Recovery baseline | Defines what Howl must be capable of managing | Does not prove implementation or authorize enablement |
| Corrected master plan | Reconciles those requirements against architecture, evidence, and dependencies | Does not silently turn a gap into a commitment |
| Investigation | Establishes what is technically possible and supported | Does not establish safe or worthwhile operation |
| Validation | Establishes whether a capability works correctly and produces a worthwhile result | Does not authorize inclusion or use |
| Owner authorization | Governs implementation and enablement | Does not retroactively validate prior work |

The six-stage progression is **not** an automatic conveyor belt:

```text
research → evidence → compatibility → selection → validation → authorization
```

Each stage may stop, fail, invalidate an earlier conclusion, or require a narrower scope. Passing research does not authorize evidence work to become implementation. Passing compatibility does not authorize selection. Passing validation does not authorize enablement. No stage grants general permission for a backend, model family, or class of techniques.

A capability appearing in the catalogue, being exposed in the UI, or being accepted by a backend is not evidence that it is correct, active, beneficial, or approved.

---

## 2. Two different status systems

### 2.1 Seven-state capability status

This status vector answers **how far a capability has been demonstrated**. It is not a release status.

| Status | Meaning | Evidence required |
|---|---|---|
| **Implemented** | Code or UI exists | Source location; existence only |
| **Wired** | Configuration reaches the intended component | Trace from setting to request/argument |
| **Supported** | The specific model/backend/hardware combination supports it | Versioned compatibility evidence |
| **Operational** | Runtime demonstrably activates it | Runtime observation; no silent fallback |
| **Validated** | Tests establish correct behavior | Repeatable behavioral tests |
| **Benchmarked** | Measured against a baseline | Named hardware, procedure, distributions, quality check |
| **Accepted** | Owner approves inclusion/use | Explicit owner decision for that scope |

The statuses are independent. A capability may be implemented and wired but unsupported; supported but not operational; operational but incorrect; correct but not beneficial; beneficial but not accepted.

A capability can occupy different statuses for different:

- model architectures;
- quantization formats;
- backends and backend versions;
- operating systems and drivers;
- GPU/CPU configurations;
- context sizes and concurrency levels;
- modalities and coexistence workloads.

A single global checkmark is therefore invalid. The verification register in §5 records the scope of each conclusion and uses `Y`, `P` (partial), `N` (no), `U` (unverified/unknown), or `—` (not applicable).

### 2.2 Six reconciliation outcomes

This matrix answers **whether the requirement is met**, not merely whether code exists.

| Outcome | Meaning |
|---|---|
| **Satisfied** | The required behavior exists and has adequate evidence for the stated scope |
| **Partially satisfied** | Some required behavior/evidence exists, but material parts are missing or unverified |
| **Missing** | No adequate implementation or evidence was found |
| **Contradicted** | The requirement conflicts with recorded plan text, current behavior, or another controlling decision |
| **Unresolved** | The requirement, scope, semantics, or evidence is not settled enough to classify as met or missing |
| **Excluded** | The owner explicitly removed the item from product scope |

“Partially satisfied” is not a soft pass. A partial implementation may still fail the requirement operationally.

---

## 3. Evidence rules

| Label | Meaning |
|---|---|
| **Source-confirmed** | Directly visible in inspected source |
| **Report-confirmed** | A cited project report records the result; this reconciliation did not rerun it |
| **Test-confirmed (reported)** | A cited report records a test pass; not rerun here |
| **Runtime-confirmed (reported)** | A cited report records runtime evidence; not reproduced here |
| **Inferred** | Follows from source structure but was not executed |
| **Unestablished** | No reliable source or evidence recovered |

No conclusion in this document is runtime-confirmed by this reconciliation. No model was loaded, no benchmark was run, no UI was exercised, and no migration was attempted.

Dependencies use the baseline keys:

| Key | Dependency |
|---|---|
| D1 | Stable chat/turn/character/media identity |
| D2 | One live prompt path and prompt-assembly ownership |
| D3 | Tokenizer and context accounting |
| D4 | Durable media identity and ownership |
| D5 | Model/resource lifecycle, admission, readiness, release |
| D6 | Model identity, family detection, template fingerprint |
| D7 | Memory schema, retrieval, compression, reinjection |
| D8 | Per-character presentation state and theme integration |
| D9 | Image-generation orchestration and gallery association |
| D10 | Shared studio job/media/resource contract |
| D11 | Protected F3/REG-1 boundaries and gate evidence |
| D12 | Inference technique/architecture/backend catalogue and conflict graph |
| D13 | Benchmark harness, named hardware, and output-quality evidence |
| D14 | Destination registry, purpose-scoped consent, local audit, and backend/sidecar egress policy |

---

## 4. Requirement reconciliation matrix

### 4.1 Chat UI and interaction

| ID | Requirement | Outcome | Evidence | Dependencies |
|---|---|---|---|---|
| REC-UI-01 | Repair user bubble border/gap and restore matching halo without losing theme/readability | **Partially satisfied** | Global bubble border/glow exists; no verified repair, gap policy, or matching user halo | D8 |
| REC-UI-02 | Per-character colors for user bubble, character bubble, halo, and text (~10 choices each) | **Partially satisfied** | Global role-based tokens exist; no per-character four-property state | D8, D1 |
| REC-UI-03 | Character portrait thumbnail beside speaker, with fallback, on every message | **Partially satisfied** | Finalized assistant rendering has a thumbnail; streaming templates omit it; fallback not complete | D1, D8 |
| REC-UI-04 | Narration comparable in purpose to SillyTavern narration | **Missing** | No narration mode found; system narrator events are different | D2, D7 |
| REC-UI-05 | Hidden-by-default token counter with show/hide control | **Missing** | No token counter or tokenizer-backed display found | D3, D2 |
| REC-UI-06 | Initiate complete image-generation orchestration from a chat message | **Missing** | No chat-message action invoking the image orchestration path found | D9, D4 |

### 4.2 Chat, per-chat, context, and memory

| ID | Requirement | Outcome | Evidence | Dependencies |
|---|---|---|---|---|
| REC-CH-01 | Up to five composer cues persisted with the chat | **Missing** | No persisted composer-cue field found in the chat schema | D1 |
| REC-CH-02 | Chat-scoped Author's Note with click-open flyout, floating/hidden presentation, side display, chat-only scope, reload retention, and blank new chat | **Missing** | No per-chat Author's Note field or presentation path found; requirement details are recovered | D1, D2 |
| REC-CH-03 | Single-character chats retain configured context capacity before rollover/memory strategies | **Partially satisfied** | Global context/rollover settings exist; per-chat behavior is not established | D3, D1 |
| REC-CH-04 | Flush/rebuild active inference context after each character response in multi-character chat | **Missing** | No character-scoped flush API; a new HTTP request is not proof of flush | D2, D3, D7 |
| REC-CH-05 | Verify whether flush existed before optimizing; test cache reuse; no full unload as optimization | **Unresolved** | T3.1 remains a research question; current truncation is approximate | D2, D3, D5 |
| REC-CH-06 | Keep 2+ characters fully separate with appropriate prompts and shared scene continuity | **Missing** | Sequential generation exists; isolation and scoped state are not established | D2, D3, D7 |
| REC-CH-07 | Character-scoped/shared lore, scoped memory, triggers, budget, precedence, conflict handling | **Partially satisfied** | Lore resolver/store and candidate memory exist; full scope and precedence are unverified | D7, D2 |
| REC-CH-08 | Optional AI-as-user/persona takeover with explicit control and isolation | **Unresolved** | Deferred optional capability; no committed implementation | D2, D7 |
| REC-CH-09 | Structured memory: inventory, ownership, acquisition/transfer/use, clothing, continuity facts | **Partially satisfied** | Prose-like memory/candidate foundation exists; typed state and full lifecycle do not | D7, D1 |
| REC-CH-10 | Automated lore/memory extraction with review, correction, conflict/duplicate checks, no silent overwrite | **Partially satisfied** | Extraction/candidate foundations exist; review/conflict behavior is incomplete | D7, D2 |
| REC-CH-11 | Conversation context rollover: shift, summarize, retain, rebuild, reinject | **Unresolved** | UI/tutorial claims exist; runtime semantics and tokenizer accuracy are unverified | D2, D3, D7 |
| ODR-08 | Multi-character flush required; single-character must not flush per turn | **Missing** | No explicit character-scoped flush contract; ordinary single-character continuation exists | D2, D3 |
| ODR-11 | Rollover and multi-character flushing remain distinct | **Partially satisfied** | Distinction is documented in the baseline; no independent implementation or tests | D2, D3, D7 |
| ODR-04 | Substantial application-side memory system | **Partially satisfied** | Records, candidates, review, scopes, persistence exist; typed inventory/retrieval/reinjection do not | D7, D1 |

### 4.3 Model, inference, provenance, and runtime

| ID | Requirement | Outcome | Evidence | Dependencies |
|---|---|---|---|---|
| REC-M-01 | Evidence-based model/template detection; never filename-only certainty | **Partially satisfied** | Model intelligence profiles and evidence/conflict states exist; full artifact/template coverage absent | D6 |
| REC-M-02 | Broad family support, automatic formatting, Jinja/ChatML, reasoning tags, fallback, override | **Partially satisfied** | Native templates and a small family table exist; comprehensive matrix absent | D6, D2 |
| REC-M-03 | Capability-gated Flash Attention, offload, multi-GPU, hardware optimization | **Partially satisfied** | Some settings are stored/passed; stored/sent/supported/effective states are not distinguished | D5, D12 |
| REC-M-04 | Comprehensive sampling including XTC/min-P/top-n-sigma/adaptive/DRY and others | **Partially satisfied** | Partial UI/request support; XTC, top-n-sigma, adaptive-P, DRY last-N/breakers, sampler order absent | D12, D13 |
| REC-M-05 | Learn from LM Studio load/resource optimizations | **Unresolved** | No measured comparison or verified internal behavior | D5, D13 |
| REC-M-06 | Useful Kobold-style offload/context features without assuming parity | **Unresolved** | Runtime flags exist; no transfer matrix or benchmark | D12, D13 |
| REC-M-07 | Application-wide hardware-aware resource orchestration with verified release | **Partially satisfied** | Advisory resource manager and lifecycle code exist; release/identity/readiness not proven | D5 |
| REC-M-08 | Hardware compatibility matrix, low-end degradation, named measurements | **Unresolved** | HW-C1–C3 not started; no hardware benchmark evidence | D13 |
| REC-M-09 | Desktop remains complete without phone/network; remote optional and unstarted | **Satisfied** (constraint preservation) | Binding constraints preserved; no remote implementation was found or added | D11 |
| REC-M-10 | REG-1 content reuse/hot-swap relationships remain protected | **Satisfied** (protected scope) | Protected record remains unchanged; this reconciliation does not test runtime | D11 |
| REC-M-11 | Comprehensive technique inventory across compute, memory, storage, cache, context, attention, batching, multi-GPU, models | **Missing** | Only a first-pass llama.cpp inventory exists | D12 |
| REC-M-12 | Map techniques to dense/MoE/SWA/recurrent/hybrid/multimodal architectures | **Missing** | No architecture applicability matrix | D12, D6 |
| REC-M-13 | Map techniques to backend/version/hardware with actual support and limits | **Partially satisfied** | Model-intelligence candidate runtimes exist; no execution/compatibility matrix | D12, D6, D5 |
| REC-M-14 | Record optimization interactions and conflicts | **Missing** | No conflict graph or tested-combination states | D12, D13 |
| REC-M-15 | Automatically detect capabilities, estimate resources, select safe combinations, explain failures | **Missing** | No architecture-aware estimator/selector | D5, D12, D13 |
| REC-M-16 | Manual override with scoped warnings, validation, reset, and provenance | **Partially satisfied** | Raw settings exist; no scoped conflict warnings/effective verification | D12, D5 |
| REC-M-17 | Benchmark performance, memory, stability, and output quality against baseline | **Missing** | No benchmark harness or quality measurement | D13, D12 |
| REC-M-18 | Capability system is not an enablement mandate; stages remain separate | **Partially satisfied** | Boundary is documented; no operational stage-gate process exists | D12, D13, D11 |
| ODR-02 | Active-conversation model hot-swap is required | **Contradicted** | Prior master-plan classification says “unaccepted”; owner requires it. Current implementation is only partial: manual stop/spawn exists, but quiescence, readiness identity, rollback, and release proof do not. | D1, D5, D6, D3 |
| ODR-05 | Preserve actual seed when meaningful; no reproduction promise | **Missing** | Global seed field exists; actual seed is not captured/parsed | D6, D3 |
| ODR-06 | Per-turn provenance is required as a concept | **Missing** | Preserved `ChatHistory` has no per-turn model/config/seed record | D1, D6, D3 |
| ODR-07 | Durable chat attachments/media | **Missing** | `hasImage` only; no durable media identity/reference/input-use record | D1, D4 |
| ODR-09 | Chat UI requirements are owner-confirmed | **Partially satisfied** | Some visual elements are partial; narration, token counter, and message image generation are absent | D8, D9, D3 |
| ODR-10 | Capability management is a first-class architecture mission | **Partially satisfied** | Mission and governance defined; multi-engine/architecture survey and benchmarks not executed | D12, D13 |
| UD-01 | Per-chat continuation-model preference and restore behavior | **Unresolved** | No `ChatMeta` model field; owner decision not made | D1, D6 |
| UD-02 | Android/remote phase, pairing/auth/API | **Unresolved** | MOB constraints exist; no protocol or authorization | D11 |
| UD-03 | Committed Story Tools scope | **Unresolved** | Page/controls exist; intended scope not established | D2 |
| UD-04 | Exact narration semantics | **Unresolved** | Purpose confirmed; source/behavior representation not recovered | D2, D7 |
| UD-05 | Studio modality release priority/scope | **Unresolved** | General direction approved; modality requirements and sequencing not set | D10, D5 |
| UD-06 | Memory review/retention policy | **Unresolved** | Required capability; review intrusiveness and auto-accept policy unset | D7 |

### 4.4 Media and attachments

| ID | Requirement | Outcome | Evidence | Dependencies |
|---|---|---|---|---|
| REC-P-01 | Durable image/attachment references; unavailable media remains unavailable | **Missing** | Chat schema has only `hasImage`; no stable media record | D1, D4 |
| REC-P-02 | Saved audio remains loadable/playable from chat | **Missing** | No saved-audio chat association found | D1, D4 |
| REC-P-03 | Displayed media and model-input media are distinct | **Partially satisfied** | Current image input is transient; durable display/reference separation is absent | D4, D9 |
| REC-P-04 | Delete/move/export/import/restore semantics for chat media | **Missing** | No lifecycle contract or unavailable-reference behavior | D4, D1 |
| REC-P-05 | Image metadata records prompt/model/options/seed | **Satisfied** (image scope) | Image-generation metadata exists; it is not chat provenance | D4, D9 |
| REC-P-06 | Imported-field provenance | **Partially satisfied** | Import work exists; provenance coverage is incomplete | D1, D6 |

### 4.5 Studio, avatars, and workspaces

| ID | Requirement | Outcome | Evidence | Dependencies |
|---|---|---|---|---|
| REC-C-01 | Separate modality-specific creative workspaces; shared catalog without merging chat/image UI | **Partially satisfied** | Separate surfaces exist; capability filtering and full integration remain incomplete | D10, D4 |
| REC-C-02 | E.6 custom workspace layout | **Satisfied** (parked disposition) | Remains parked; no promotion by the studio direction | D10 |
| REC-C-03 | E.8 image staging-board structure | **Satisfied** (parked disposition) | Remains parked; no promotion by the studio direction | D10 |
| REC-C-04 | Music/audio expansion within approved studio direction | **Missing** | No music generator; F2-B7/TTS remains frozen | D10, D5 |
| REC-C-05 | Video generation as a complete media subsystem | **Partially satisfied** | Video engine exists; no complete catalogue/sidecar/delete/attachment contract | D4, D10 |
| REC-C-06 | Reactive desktop 3D avatar / multi-character avatars | **Excluded** | Owner excluded the desktop surface; historical note retained only as evidence | D11 |
| REC-C-07 | 3D mesh and texture generation as creative assets | **Missing** | Approved direction; no production subsystem | D10, D4, D5 |
| REC-C-08 | General local creative studio beyond current surfaces | **Missing** | Approved direction; no unified studio system | D10, D4, D5 |
| ODR-01 | Desktop character avatars are excluded | **Excluded** | Owner exclusion controls; historical tier-4 note is evidence only, not a roadmap item | D11 |
| ODR-03 | General local creative studio is an approved direction | **Contradicted** | Prior master plan calls it “unestablished”; owner approves the direction. Implementation of individual modalities remains missing. | D10, D4, D5 |

### 4.6 Governance and protected boundaries

| ID | Requirement/decision | Outcome | Evidence | Dependencies |
|---|---|---|---|---|
| REC-UI overall | Owner-confirmed chat UI requirements are satisfied | **Partially satisfied** | Some styling/portrait pieces exist; narration, token counter, image action absent | D8, D9, D3 |
| F3-C Phase 0 | Remains open | **Unresolved** | Protected record unchanged; no acceptance | D11 |
| R6 | Unauthorized | **Satisfied** (constraint preservation) | No authorization inferred from this reconciliation | D11 |
| F3-C acceptance | Not granted | **Satisfied** (constraint preservation) | No acceptance record created | D11 |
| F3-D | Blocked | **Satisfied** (constraint preservation) | No unblocking action taken | D11 |
| D-ROOT | Six defects open | **Unresolved** | Protected status unchanged | D11 |
| D-STORAGE | Blocked | **Satisfied** (constraint preservation) | No storage work attempted | D11 |
| V20 | Blocked pending F3-F | **Satisfied** (constraint preservation) | No gate change | D11 |
| Gate v2 V1–V17 | Adopted | **Satisfied** (record preservation) | No change | D11 |
| Gate v2 V18–V24 | Proposed, not adopted | **Satisfied** (record preservation) | No adoption inferred | D11 |
| TTS / F2-B7 | Frozen | **Satisfied** (constraint preservation) | No TTS change attempted | D11 |
| F2-B5 | Approved direction | **Satisfied** (record preservation) | Workspace direction remains approved | D11 |
| E.6/E.8 | Parked | **Satisfied** (record preservation) | No promotion | D11 |
| REG-1 | Protected | **Satisfied** (record preservation) | Content hot-swap remains distinct | D11 |
| MOB-1–MOB-4 | Binding conditional constraints | **Satisfied** (constraint preservation) | No remote implementation | D11 |
| W4 | Historical/superseded; formal closure gap | **Unresolved** | W5–W7 supersede practical blocker; no explicit closure record | D11 |
| W8 | Verified | **Satisfied** (record preservation) | Verified record remains unchanged | D11 |

### 4.7 Network boundaries

| ID | Requirement | Outcome | Evidence | Dependencies |
|---|---|---|---|---|
| NET-01 | Local creative work and inference remain usable offline | **Unresolved** | Local backends exist; no egress-blocked end-to-end test | D14, D5 |
| NET-02 | No cloud inference or cloud model provider in the local baseline | **Contradicted** | DashScope ASR and remappable MeloTTS are exposed without separate authorization | D14, D5 |
| NET-03 | No background model lookup, telemetry, provider connection, or update check | **Contradicted** | Civitai hash lookup is called during image scan; addendum N-01 | D14, D6 |
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
| NET-19 | Every proposed path has disclosure, cache, failure, and authorization rules before implementation | **Missing** | Addendum P-01–P-08 are proposed, not authorized | D14 |
| NET-20 | No automatic update check; future check explicit and signed | **Satisfied** (absence only) | No update checker found; must remain absent | D14 |

---

## 5. Existing-capability verification register

### 5.1 Scope warning

This register is about **existing capability**, not future research. It uses a single source-level/pinned-runtime scope. A row does not establish behavior on every model, backend, quantization, or machine.

Status vector key: `I` implemented, `W` wired, `S` supported, `O` operational, `V` validated, `B` benchmarked, `A` accepted. `Y` yes, `P` partial, `N` no, `U` unverified/unknown, `—` not applicable.

“Accepted” means the owner approves inclusion/use for that capability scope. It does not mean the current implementation is accepted as correct, and it is not F3-C acceptance or a freeze.

### 5.2 Register

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

### 5.3 Flash Attention worked example

The seven questions produce a deliberately unflattering result for the current evidence:

| Question | Current answer | Reason |
|---|---|---|
| Implemented? | **Yes** | A Flash Attention setting/option exists |
| Wired? | **Partial** | It reaches a server-argument path conditionally; the exact setting-to-runtime path is not proven end-to-end |
| Supported? | **Unknown** | Depends on model architecture, quantization, backend build, and hardware |
| Activated? | **Unknown** | No runtime proof that the backend enabled it rather than falling back |
| Correct? | **Unknown** | No behavioral or numerical validation on target combinations |
| Beneficial? | **No benchmark** | No measured speed/memory gain or trade-off |
| Worth retaining? | **Unresolved** | No evidence basis for inclusion decisions beyond a control existing |

The checkbox answers only the first question. The same discipline applies to GPU offload, context shift, cache reuse, quantization, multi-GPU, and every future technique.

---

## 6. Contradictions and high-risk gaps

### 6.1 Contradictions requiring plan correction

| Contradiction | Controlling evidence | Required reconciliation |
|---|---|---|
| Master plan calls model hot-swap “unaccepted” | Owner decision makes it required product behavior; current implementation is partial | Plan must separate product requirement from implementation/F3-C status |
| Master plan calls general creative studio “unestablished” | Owner approved general direction | Plan must record approved direction without authorizing modalities |
| Master plan retains desktop avatar as ordinary parked future | Owner explicitly excluded it | Plan must mark excluded; historical note only |
| Context rollover and multi-character flushing can be read as one feature | Owner definitions make them separate | Plan must separate triggers, scope, mechanisms, and tests |
| Author's Note “persisted” appears to conflict with “cleared on new chat” | Existing chat retention and new-chat blank state can coexist | Plan must define per-chat retention and no inheritance |
| First-pass llama.cpp inventory can read as a capability survey | Only one binary and subset were inspected | Plan must keep the catalogue unfrozen and identify research gaps |
| “Local-first” is stated as a completed property | Image scans call Civitai automatically with a model hash; remote ASR/TTS paths exist | Plan must state the boundary is required but currently contradicted and add the complete path register |
| Local-Only Mode is treated as enforcing offline operation | The only inspected implementation wraps renderer `fetch`/XHR; Go, sidecar, and Wails paths remain outside it | Plan must require backend/sidecar enforcement and a runtime offline test |
| Generic consent is treated as informed consent | The image-dependency sentence omits destination, received data, recurrence, and retention; audio has labels only | Plan must require the six-question, purpose-scoped model |
| Error-report UI claims “no model names” | The report includes an arbitrary `Context` map; no redaction proof was found | Plan must verify the actual payload, not repeat the UI claim |

### 6.2 High-risk gaps

1. **No effective-capability verification.** A control may be present, wired, supported, operational, correct, beneficial, and accepted at different levels; Howl currently records mostly the first.
2. **No existing-capability campaign.** The missing-feature inventory is not enough; controls already present need the same seven-question verification.
3. **No architecture-aware selection.** Dense-transformer assumptions can mis-handle MoE, SWA, recurrent, and hybrid models.
4. **No conflict graph.** Individually valid options can disable, invalidate, or degrade each other.
5. **No output-quality benchmark.** Speed/memory claims are incomplete without coherence, accuracy, continuity, and structured-output checks.
6. **No per-turn provenance or actual seed capture.** The product cannot reconstruct which configuration produced a response.
7. **No durable attachment identity.** Display, model input, and stored object are conflated.
8. **No safe hot-swap lifecycle.** Streams, readiness, rollback, and release are not coordinated.
9. **No tokenizer-accurate context budget.** Memory, lore, token counter, and rollover lack a common accounting basis.
10. **No demonstrated character-context isolation.** Sequential HTTP calls do not prove a flushed inference context.
11. **No backend network boundary.** The renderer-only toggle does not constrain Go HTTP clients, DNS/dial calls, sidecars, proxies, or browser-open actions.
12. **No destination registry or connection audit.** There is no authoritative record of what may be contacted, why, with what data, or what happened.
13. **Remote inference sits inside an ostensibly local baseline.** Remote ASR/TTS are exposed without an owner decision separating them from local inference or a purpose-specific consent flow.
14. **External actions are easy to conflate.** Opening a publisher page, downloading a file, and allowing Howl to index it are not yet three separately governed actions.

---

## 7. Dependency and execution order

The reconciliation order is:

1. Preserve protected F3/REG-1 boundaries and record no implied authorization.
2. Reconcile requirement outcomes and dependencies against the baseline.
3. Establish one live prompt path and stable identity/turn schema.
4. Establish tokenizer/context accounting.
5. Reconcile media identity and attachment lifecycle.
6. Reconcile memory schema, retrieval, reinjection, and scope.
7. Reconcile model identity, family/template evidence, and hot-swap lifecycle.
8. Establish the existing-capability verification campaign and seven-state register.
9. Execute the multi-engine/architecture/backend survey and conflict graph.
10. Benchmark safe combinations, including output quality, before enabling.
11. Obtain explicit owner authorization for implementation/enablement.
12. Establish the destination registry, backend offline enforcement, six-question consent, audit, cache/failure behavior, and release acceptance evidence before adding or enabling any external integration.
13. Require bidirectional traceability: every corrected-plan entry/work package must reference its reconciliation ID(s), and every reconciliation row must name a disposition, work package, or explicit unresolved status.
14. Keep effort estimates provisional until the architecture and existing-capability investigations establish what is present; do not turn planning assumptions into commitments.

No item in this order is implemented by this document.

---

## 8. Verification and limitations

### 8.1 Verification performed

- The original 66 requirement rows are preserved; NET-01–NET-20 are additive and do not replace any row.
- The matrix covers REC-UI-01–06, REC-CH-01–11, REC-M-01–18, REC-P-01–06, REC-C-01–08, ODR-01–12, NET-01–NET-20, UD-01–06, and the protected-boundary register.
- The network register covers N-01–N-14, L-01–L-05, and proposed paths P-01–P-08; each inspected production external path is classified as an external path, a local/loopback path, or a test/dev exclusion.
- Every outcome uses one of the six reconciliation categories.
- Every capability row carries an explicit seven-state vector and scope.
- The baseline hash is recorded in the header.
- No source, test, configuration, protected F3 artifact, master plan, adjudication, owner packet, or acceptance record was modified.
- No placeholder markers are used.

### 8.2 Limitations

- This is a static reconciliation, not runtime verification.
- No model was loaded, no UI exercised, no test/build/migration run, and no benchmark performed.
- Capability status is not generalized across models, backends, hardware, or versions.
- The exact UI planning document and expansion text remain unrecovered.
- Narration semantics remain unresolved.
- Actual seed availability on Howl's endpoint remains unproven.
- Weight release and VRAM reclamation remain unmeasured.
- The capability catalogue remains unfrozen and the multi-engine survey unstarted.
- The network inventory is static-source evidence. No outbound request, packet capture, offline run, cancellation test, consent/audit test, remote ASR/TTS test, or Android test was performed.
- Local-Only Mode was not demonstrated to block Go `net/http`, `net.Dial`, sidecar egress, or Wails-mediated actions; its inspected implementation is renderer-only.
- The corrected master plan has not been prepared or adopted.

### 8.3 Final status

This reconciliation establishes what is satisfied, partial, missing, contradicted, unresolved, or excluded without equating implementation presence with verified capability. The network boundary is now included as a first-class requirement family and as a source-confirmed contradiction of the current local-first claim. This document is an input to the corrected master plan, not a substitute for it and not an implementation authorization; it does not authorize or validate any external destination.

The corrected plan is not complete until it consumes the original 66 rows plus NET-01–NET-20, preserves the seven-state capability register and all six outcomes, and maintains bidirectional traceability: every plan entry maps back to a reconciliation ID, and every row maps forward to a disposition, work package, or explicit unresolved status. Effort estimates remain provisional until the architecture and existing-capability investigations establish the implementation baseline. No decision, estimate, or work package in this reconciliation is an implementation authorization.
