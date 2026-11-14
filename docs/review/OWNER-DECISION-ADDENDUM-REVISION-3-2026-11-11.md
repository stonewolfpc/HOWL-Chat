# Owner-decision addendum — Revision 3 planning baseline

**Document ID:** `OWNER-DECISION-ADDENDUM-REVISION-3-2026-11-11`  
**Date:** 2026-11-11  
**Status:** **OWNER PRODUCT DECISIONS RECORDED — PLANNING ACCEPTANCE ONLY; NO IMPLEMENTATION AUTHORIZATION**  
**Baseline:** `docs/review/CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md`, Revision 3  
**Baseline SHA-256:** `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` — **frozen, unmodified**  
**Baseline size:** 131,755 bytes / 1,223 lines  

## 1. Authority and boundary

Revision 3 is accepted as the **planning baseline**. Its embedded status line still reads “NOT ADOPTED” because that file is frozen at the revision the owner reviewed. This addendum is the separate record of planning acceptance; it does not rewrite the plan.

Planning acceptance is **not** implementation authorization.

This addendum authorizes:

- read-only inspection of existing code and documentation;
- architecture research, comparison, threat modeling, and design;
- production of proposed interfaces, data contracts, work packages, and acceptance designs;
- no source, test, runtime, migration, prompt, network, model, or settings change.

It does **not** authorize:

- model loads, benchmarks, or runtime resource measurement;
- network experiments, listeners, pairing, or remote access;
- schema migration or durable-state cutover;
- prompt or memory injection changes;
- enablement of RAG, Android, studios, narration, or hardware optimization;
- any claim that a design is implemented, validated, or accepted.

The protected F3-C Phase 0, D-ROOT, and D-STORAGE boundaries remain in force.

## 2. Frozen inputs

| Input | SHA-256 | Handling |
|---|---|---|
| Corrected master plan, Revision 3 | `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` | Frozen; not edited by this addendum |
| Requirements reconciliation | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` | Unchanged; 66 rows preserved |
| Recovery owner addendum | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` | Unchanged |
| Revision 3 verification | `1ccfafd2fc166f700e11e35a14eb1e0a08d5906f833e0925093f4d355e441f92` | Historical verification record |
| MUSE N1–N7 closure review | `docs/review/REVIEW-MUSE-N1-N7-CLOSURE-2026-11-11.md` | Reviewed against `5ff5e1d4…`; O-1–O-3 changes are documentary |

The N1–N7 closure review was not run against Revision 3. Revision 3 is the reviewed content plus the owner-approved O-1–O-3 documentary cleanups. No new technical review is claimed here.

## 3. Owner decisions

The identifiers `UD-01`–`UD-06` are retained so the decision trail remains compatible with Revision 3.

| ID | Owner decision | Status | Consequence |
|---|---|---|---|
| **UD-01** | Each chat remembers its selected model. Hot-swapping is convenience, not a release blocker. | **Resolved** | Per-chat model binding is a product behavior. A chat must reload with the model it selected. Hot-swap must not corrupt state, but release does not wait for seamless active-stream swapping. |
| **UD-02** | Android remote access is committed product scope. | **Resolved as product scope; security design open** | Desktop remains the model host. Android remotely accesses chats, images, and application functions. No phone-hosted model and no cloud inference. Pairing, authentication, encryption, revocation, and transport remain design work. |
| **UD-03** | Story Tools is committed product scope. | **Resolved as product scope; detailed capability set open** | Existing page/scripts are evidence, not an automatic feature list. Story Tools requires a bounded capability map and ownership model, not removal or unexamined expansion. |
| **UD-04** | Narration is chat-only AI takeover of the user's turn; no voice or images. | **Resolved** | AI-authored user-role output must be explicitly attributed and reversible. It must not be presented as the user's own writing. Voice and image generation are out of this scope. |
| **UD-05** | Character voice/TTS belongs to Chat. Image is the next major capability to finalize. Music and 3D are future dedicated studios; other modalities come later. | **Resolved as destination and staging** | TTS is not a separate studio/tab. TTS/ASR remains a technical resource consumer. Image is next; Music and 3D are not authorized. |
| **UD-06** | Memory is automatic, transparent, editable, and robust with small local models. | **Resolved** | Automatic capture and visible state management are required. The model may propose state but may not silently overwrite it. Small-model operation is a design constraint, not an afterthought. |
| **ODA-07** | Local RAG is a vital creative capability, not merely document search. | **Resolved** | RAG is committed as local creative infrastructure: indexing, hybrid retrieval, structure awareness, citations, and exploration/critique workflows. It remains local-only and separate from canonical fictional state. |
| **ODA-08** | `OA-09` clarification: design work is not a runtime investigation. | **Resolved as interpretation** | Read-only code/document inspection, static analysis, web research, comparison, and architecture design may proceed. Model loads, benchmarks, network experiments, user-content inspection, and runtime evidence capture still require explicit `OA-09` authorization. |
| **ODA-09** | Hardware detection and optimization are core functionality, not an optional power-user feature. | **Resolved** | Automatic inventory, suitability assessment, safe admission, and conservative fallback are part of the normal product experience. Power-user overrides exist, with warnings, rather than replacing automatic behavior. |

### 3.1 Scope consequences that must not be overstated

- **UD-01** does not make seamless hot-swap a release blocker. It does require a chat's model selection to persist and be restored.
- **UD-02** commits Android as scope; it does not select WebRTC, WireGuard, libp2p, QUIC, a relay operator, or a public listener.
- **UD-03** commits Story Tools; it does not promote every existing control into a requirement.
- **UD-04** permits AI takeover of the user's turn only in chat, with explicit attribution and no voice/image extension.
- **UD-05** places character voice/TTS in Chat and stages Image next; Music and 3D are future dedicated studios. It authorizes no modality.
- **UD-06** requires automation, but not silent mutation. The state authority remains Howl.
- **ODA-07** resolves the RAG product-scope question. It does not enable RAG, load an embedding model, create an index at runtime, or relax `OA-10` lifecycle evidence.
- **ODA-08** does not turn design into implementation or authorize any runtime investigation.
- **ODA-09** does not authorize telemetry experiments, GPU queries, or overclocking. It requires design for those functions.
- **TTS/ASR remains a technical resource consumer**, not a product studio. Character voice belongs to Chat even though it requests a model/resource lease.

## 4. Controlled delta register

The accepted plan and reconciliation are not rewritten. These are delta product requirements for the design package. They are not inserted into the 66-row reconciliation until a separate reconciliation change is reviewed and accepted.

| Delta ID | Product requirement | Source decision | Design owner |
|---|---|---|---|
| **DPR-01** | A chat persists and restores its selected model; model switching is safe at a turn boundary and does not depend on seamless mid-stream swapping. | UD-01 | Memory/chat identity + resource orchestration |
| **DPR-02** | Android is a remote client of the desktop host, off by default, with explicit pairing, mutual device authentication, end-to-end encryption, revocation, and no cloud chat storage. | UD-02, `NET-12`, `MOB-1`–`MOB-4` | Android security research |
| **DPR-03** | Story Tools is a committed writing workspace with bounded capabilities and explicit ownership. | UD-03 | Reconciled architecture plan |
| **DPR-04** | Narration is a chat-only, explicitly attributed AI takeover of the user's turn, without voice or image output. | UD-04 | Reconciled architecture plan |
| **DPR-05** | Character voice/TTS is part of Chat; Image is the next major capability; Music and 3D are future dedicated studios, with separate gates for each. | UD-05 | Reconciled architecture plan |
| **DPR-06** | Memory is automatic, visible, editable, provenance-bearing, conflict-aware, and operable with small local models. | UD-06, `REC-CH-09`, `REC-CH-10`, `ODR-04` | Memory and continuity design |
| **DPR-07** | RAG is a local creative subsystem: structure-aware indexing, hybrid retrieval, evidence references, and exploration/critique workflows, separate from canonical state. | ODA-07 | Local RAG design |
| **DPR-08** | Hardware inventory, suitability assessment, 10%-free-VRAM reserve, vendor-neutral telemetry, fallback, and explicit overrides are core behavior. | ODA-09, `REC-M-07`, `REC-M-08`, `REC-M-15`, `REC-M-16` | Hardware/resource design |
| **DPR-09** | Design/research is permitted without runtime authorization; gated runtime work remains gated. | ODA-08 | All design packages |

## 5. Effect on existing plan decisions

| Existing item | Effect of owner decisions | Still open |
|---|---|---|
| `UD-01` in `CWP-01`/`CWP-08` | Resolved as per-chat model binding; hot-swap is not a release blocker | Durable schema and migration approval |
| `UD-02` in `CWP-25` | Resolved as committed product scope | Security architecture, `OA-06`, transport selection |
| `UD-03` in `CWP-28` | Resolved as committed scope | Capability map and ownership |
| `UD-04` in `CWP-17` | Resolved as chat-only AI takeover | Semantics, attribution design, implementation gate |
| `UD-05` in `CWP-19` | Resolved as Chat/TTS destination and Image-next staging | Dedicated Music/3D design and authorization |
| `UD-06` in `CWP-06` | Resolved as automatic, transparent, editable memory | Schema, extraction policy, retrieval design |
| `OA-10` | Product scope is resolved in favor of committed local RAG | Implementation, embedding choice, runtime indexing, lifecycle evidence |
| `OA-09` | Interpretation clarified for design/research | Any runtime investigation remains closed |
| `AD-05` | Still recorded at the `15a → 08a` boundary | Resource contract and enforcement design |
| F3-C / D-ROOT / D-STORAGE | Still protected | Owner-controlled phase-0 and migration prerequisites |

## 6. Non-goals

This addendum does not:

- amend F3-C, R6, F3-D, D-STORAGE, V20, REG-1, Gate v2, or any protected artifact;
- create a new accepted plan version;
- change the 66-row reconciliation or 37-row capability register;
- choose the Android transport or relay operator;
- choose the memory database, embedding model, or GPU telemetry library;
- authorize a single source or runtime change;
- claim any design is implemented.

## 7. Design package references

- `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md`
- `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`
- `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md`
- `docs/architecture/RESEARCH-ANDROID-REMOTE-SECURITY-2026-11-11.md`
- `docs/architecture/RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md`

## 8. Acceptance state

| Question | Answer |
|---|---|
| Is Revision 3 the accepted planning baseline? | **Yes, by this addendum.** |
| Was the plan file edited to say so? | **No.** |
| Are the owner product decisions recorded? | **Yes.** |
| Is RAG product scope committed? | **Yes.** |
| Is RAG implementation authorized? | **No.** |
| Is Android product scope committed? | **Yes.** |
| Is Android implementation authorized? | **No.** |
| Are runtime investigations authorized? | **No.** |
| May design work proceed? | **Yes, within this addendum.** |
| Is any work package accepted or gate passed? | **No.** |

No implementation authorization is granted by this document.
