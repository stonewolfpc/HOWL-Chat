# Reconciled architecture and controlled design plan

**Status:** Proposed design plan; no implementation authorization  
**Date:** 2026-11-11  
**Planning baseline:** Revision 3, `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874`  
**Owner decision record:** `docs/review/OWNER-DECISION-ADDENDUM-REVISION-3-2026-11-11.md`  
**Scope:** Reconcile the memory/continuity, hardware/resource, local RAG, and Android security designs into one architecture without duplicate systems or contradictory authority.

## 1. Decision in one sentence

Howl will use one canonical state/event service, one prompt assembler, one resource authority, one local corpus/retrieval subsystem, and one remote-access gateway; each has a single owner, explicit boundaries, and design-only work packages that cannot drift into implementation.

## 2. Baseline and change control

| Artifact | Identity | Rule |
|---|---|---|
| Corrected plan, Revision 3 | `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` | Frozen; not edited |
| Reconciliation | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` | Unchanged; 66 rows |
| Recovery addendum | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` | Unchanged |
| Owner decisions | `OWNER-DECISION-ADDENDUM-REVISION-3-2026-11-11.md` | Planning authority only |
| Design docs | See §10 | Proposed; no runtime evidence |

Revision 3's internal status line still says “NOT ADOPTED.” Planning acceptance is recorded in the separate owner-decision addendum. This preserves the frozen plan while making the acceptance trail explicit.

The 66-row reconciliation is not reopened by these designs. `DPR-01`–`DPR-09` are delta product requirements. RAG remains outside the 66 until a separate reconciliation change is reviewed.

## 3. Architecture at a glance

```text
                         ┌─────────────────────────────┐
                         │        Howl desktop         │
                         │                             │
  Chat / Story Tools ───▶│  One prompt assembler       │
  Narration ────────────▶│    ├─ canon state card      │
                         │    ├─ memory/event context  │
                         │    ├─ RAG evidence         │
                         │    └─ response budget      │
                         │                             │
                         │  Canonical state + events   │
                         │  Memory/retrieval           │
                         │  Reference corpora         │
                         │  Resource authority         │
                         │  Local model backends       │
                         └─────────────┬───────────────┘
                                       │
                         ┌─────────────▼───────────────┐
                         │   Howl Link Gateway         │
                         │   pairing / sessions / E2E  │
                         └─────────────┬───────────────┘
                                       │
                              Android remote client
                         (no models, no cloud inference)
```

## 4. Single-owner boundaries

| System | Single owner | Owns | Must not own |
|---|---|---|---|
| **Canonical state/event service** | Memory/continuity design | Character state, event log, revisions, conflicts, undo | Prompt assembly, reference retrieval, model lifecycle |
| **Memory index** | Memory/continuity design | Searchable distilled memories and event retrieval | Canonical state, external references |
| **Prompt assembler** | `CWP-02` design | Exact-once contributor order and token budget | A second storage path or independent state |
| **Resource authority** | `CWP-15a/15b` | Inventory, suitability, leases, admission, release | Model quality claims, prompt content |
| **Reference corpus/RAG** | Local RAG design | Corpus, structure, hybrid index, citations, creative modes | Canonical state or chat truth |
| **Remote gateway** | Android design + `CWP-20` | Pairing, device auth, E2E sessions, remote policy | Chat storage, model hosting, cloud inference |
| **Story Tools** | `CWP-28` design | Writing workspace capabilities | A separate memory/state store |
| **Narration** | `CWP-17` design | Chat-only AI takeover semantics | Voice/image generation or unattributed user messages |

### 4.1 Non-duplication rules

- There is no second memory file in the renderer. The frontend cache is a read model.
- Memory records are derived from events; they do not become a second canonical state.
- RAG cannot write state or be labeled canon.
- Lorebooks remain their own authored reference/trigger system, not RAG and not character state.
- Chat, image, TTS, ASR, and embedding workloads all use the resource authority.
- The Android gateway does not implement its own model/resource selection.
- Story Tools consumes the same state/memory APIs; it does not create a parallel store.
- Narration produces an attributed message in the same chat pipeline; it does not fork the prompt path.

## 5. Resolved cross-design conflicts

| Conflict | Resolution |
|---|---|
| Frontend `localStorage` memory versus backend memory | Backend is authoritative; frontend is a cache/read model with explicit sync state |
| Prose memory versus typed character state | Both exist: state is canonical; memory is derived, searchable evidence |
| Historical event versus current state | History is auditable and append-oriented; current state is a projection; corrections, supersession, deletion, and as-of queries are explicit |
| Automatic capture versus no silent overwrite | Automatic capture; evidence-backed reversible commits; review for ambiguous/destructive changes |
| Model-selected state versus Howl authority | Model returns proposals only; deterministic application layer commits |
| Global model selection versus per-chat decision (`UD-01`) | Chat stores model binding; no silent substitution; hot-swap not a release blocker |
| Resource manager versus image/TTS/embedding consumers | One lease/admission authority; every consumer registers or is an orphan |
| Existing 90% safety constant versus 10% reserve | Formalize 10% of currently free VRAM as reserve; explicit overhead separate |
| RAG versus canonical fictional world | Separate corpus with `REFERENCE — NOT CANON` labels; RAG cannot write state |
| RAG versus Story Tools | Story Tools may consume RAG/state; it does not own either |
| Android remote versus Local-Only Mode | Remote is off by default; desktop policy is rechecked per request; no public listener |
| Android committed (`UD-02`) versus `OA-06` closed | Product scope committed; implementation and network experiments remain gated |
| Narration (`UD-04`) versus optional persona | Narration is a chat-only AI takeover with explicit attribution; it is not a general AI-as-user permission |
| Product staging (`UD-05`) versus modality authorization | Character voice/TTS is in Chat; Image is next; Music/3D are future dedicated studios; each still needs its own design/authorization |

## 6. Data flow and authority

### 6.1 A chat turn

```text
user turn
  → state/event service records the turn
  → prompt assembler requests:
       canon state card
       relevant memory/events
       optional RAG evidence
       recent turns
  → model produces text + optional state proposals
  → validator applies/queues proposals
  → response is stored with provenance
  → visible diff and undo remain available
```

The model never writes the store directly. The prompt assembler never reads raw files directly.

### 6.2 A creative RAG request

```text
request + mode
  → corpus/structure filters
  → lexical + vector retrieval
  → fusion, dedup, optional rerank
  → cited evidence package
  → creative output with evidence panel
```

RAG may read canon state to perform a continuity check, but it may not mutate it.

### 6.3 A model load

```text
chat model binding
  → hardware inventory
  → suitability report
  → resource lease (10% reserve)
  → load at safe boundary
  → ready / busy / release
  → reconcile actual vs tracked
```

No backend bypasses the lease.

### 6.4 An Android request

```text
paired Android device
  → authenticated E2E envelope
  → desktop gateway policy check
  → authorized chat/image/action
  → response encrypted to device
```

The relay, if any, sees ciphertext only.

### 6.5 Prompt-path cutover

`ADR-INT-01` is a target, not current behavior. The current build still has several prompt contributors and call sites. A cutover must name them and retire them rather than adding one more path.

| Current path | Current role | Cutover action |
|---|---|---|
| `app-prompt.go:buildPromptMessages` | Application-level message assembly and wrapper text | Become an adapter or be retired behind the assembler |
| `internal/backend/prompt/builder.go:BuildPrompt` | Backend prompt/template construction | Consume one assembler request; no independent precedence |
| Lore paths | Triggered lore selection and injection | Contribute a scoped, deduplicated lore package |
| Chat-service snapshots | Memory/lore state assembled near the chat turn | Supply a turn snapshot, not prompt strings |
| `retrieveRAG` | Retrieval and evidence selection | Return an `EvidencePackage`; no direct injection |
| Frontend character messages | Character/system text assembled in the client | Be retired or submit structured character identity to the assembler |
| Narration/Story Tools | Attributed AI takeover using the same chat pipeline | Use the same assembler; no forked prompt path |

**Target contract:** one assembler owns contributor precedence, inclusion, deduplication, token budget, and provenance. A proposed order, subject to `CWP-02` validation, is:

1. safety/system constraints;
2. character identity and voice;
3. world/scene/scenario state;
4. canonical state card;
5. scoped authored lore;
6. memory/event evidence;
7. RAG evidence, labeled `REFERENCE — NOT CANON`;
8. recent transcript;
9. the current user turn.

Exactly-once rules:

- Each contributor has an identity and is included at most once per turn.
- The assembler, not a contributor, enforces global token budget and truncation order.
- Memory/lore/RAG snapshots are fixed for the turn; a contributor cannot mutate shared state while assembling.
- Missing or stale evidence is represented as missing/degraded, never as invented text.
- RAG and canonical events have distinct labels and cannot substitute for each other.
- The model receives no raw storage handles or write authority.

Cutover design requirements: inventory every path, define equivalence fixtures, wrap old paths as adapters, verify no double injection, switch by capability with a visible fallback, then retire legacy paths. A flag or adapter is not a second source of truth. No runtime cutover is authorized here.

## 7. Design work packages

`DWP-*` are **design work packages**, not `CWP-*` implementation packages. They produce decision records, interfaces, threat models, migration plans, and test designs. They do not authorize code, migrations, model loads, or network activity.

The numeric order below is **design-input order**, not a build or delivery sequence. Product staging is defined separately in §7.1.

| ID | Design package | Outputs | Prerequisites | Authorization |
|---|---|---|---|---|
| **DWP-01** | Owner-decision and requirements delta | Addendum, `DPR-01`–`DPR-09`, trace map | Owner direction | **Complete** (planning only) |
| **DWP-02** | Canonical state/event model | State schema, event schema, revision/conflict/undo semantics, as-of contract | DWP-01; F3-C/D-ROOT/D-STORAGE constraints identified | Design only |
| **DWP-03** | Memory, retrieval, and context assembly | Memory/event retrieval contract, state card, small-model budget, prompt precedence | DWP-02; `CWP-02`/`CWP-03` design | Design only |
| **DWP-04** | Hardware inventory and telemetry | Provider interface, vendor matrix, confidence model, inventory schema | DWP-01; no runtime tests | Design only |
| **DWP-05** | Resource authority and model binding | Atomic lease/commit, concurrency/expiry, rollback, verified release, 10% reserve, fit report, per-chat binding | DWP-04; `CWP-15a` design authority | Design only |
| **DWP-06** | Local RAG creative subsystem | Corpus/rights model, structure/chunking, hybrid index, citation contract, creative modes | DWP-02 for canon separation; DWP-05 for embedding resource path | Design only |
| **DWP-07** | Android remote security | Threat model, pairing/auth/E2E design, transport comparison, revocation, test plan; residual Android data-at-rest, key/session rotation, relay metadata, pairing UX, scoped auth, and offline behavior | DWP-01; `CWP-20` control-plane design; `OA-09` for any network spike | Design/research only |
| **DWP-08** | Product surfaces and sequencing | Story Tools capability map, narration semantics, Chat voice/TTS destination, Image-next staging, Music/3D future-studio gates | DWP-01; DWP-02/03 | Design only |
| **DWP-09** | Integration and acceptance architecture | Cross-system interfaces, dependency graph, acceptance matrix, ADR set | DWP-02–DWP-08 drafts | Design only |

### 7.1 Product staging, not DWP order

| Stage | Product scope | Exit condition |
|---|---|---|
| **1. Chat finalization** | Chat prompt/state/memory; character voice/TTS inside Chat; resource authority; local-only boundary | Chat behavior, state/memory, prompt cutover, and resource lifecycle have accepted design and later authorized gates |
| **2. Image finalization** | Image generation/editing as the next major capability; image engine registered with the resource authority | Image design and resource/quality gates accepted separately |
| **3. Expansion** | Music and 3D as future dedicated studios; other modalities later | Each dedicated studio has its own design, evidence, and authorization |
| **Later UI unification** | Shared surfaces across Chat, Image, Story Tools, and future studios | Does not block usability fixes in the current stage |

This preserves Chat → Image → expansion. It does not convert the DWP numbers into a build order, promise dates, or authorize a modality.

## 8. Dependency graph

```text
DWP-01
  ├─▶ DWP-02 ──▶ DWP-03 ──▶ DWP-08
  │       └────▶ DWP-06
  ├─▶ DWP-04 ──▶ DWP-05 ──▶ DWP-06 (resource path)
  │                    └──▶ DWP-07 (host capacity/policy)
  └─▶ DWP-07
DWP-02..DWP-08 ──▶ DWP-09
```

The graph is acyclic and is a **design-input graph**, not a build schedule. It does not reopen the accepted plan's `CWP-08`/`CWP-15` sequence; it supplies design inputs to `CWP-15a`, `CWP-08a`, `CWP-15b`, and `CWP-08b` in that order.

## 9. Architecture decisions to record

| ID | Decision | Current status |
|---|---|---|
| `ADR-INT-01` | One prompt assembler owns contributor order and budget | Proposed; cutover contract in §6.5 |
| `ADR-INT-02` | Howl, not the model, owns canonical state | Proposed |
| `ADR-MEM-01` | Auditable append-oriented history plus revisioned projections | Proposed |
| `ADR-MEM-02` | No silent overwrite; user locks and conflicts | Proposed |
| `ADR-HW-01` | 10% of currently free VRAM is reserved by default | Proposed |
| `ADR-HW-02` | Telemetry has source/confidence; unknown is not zero | Proposed |
| `ADR-HW-03` | Every model consumer uses one resource lease authority | Proposed |
| `ADR-RAG-01` | Reference corpus is separate from canon and cannot write state | Proposed |
| `ADR-RAG-02` | Hybrid lexical/semantic retrieval with source citations | Proposed |
| `ADR-AND-01` | Android is off by default; no unauthenticated public listener | Proposed |
| `ADR-AND-02` | Pairing, mutual auth, E2E, revocation before remote use | Proposed |
| `ADR-STU-01` | Character voice/TTS is in Chat; Image is next; Music/3D are future dedicated studios; no modality is thereby authorized | Proposed |

## 10. Traceability

| Source or decision | Design coverage |
|---|---|
| `REC-CH-09` structured character state | Memory design §4–§6 |
| `REC-CH-10` extraction/review/conflict | Memory design §5, §7 |
| `ODR-04` application-owned memory | Memory design; addendum `DPR-06` |
| `REC-M-07` hardware orchestration | Hardware design §4–§10 |
| `REC-M-08` compatibility/degradation | Hardware design §7, §9 |
| `REC-M-15` automatic capability/resource selection | Hardware design §7–§9 |
| `REC-M-16` overrides/warnings | Hardware design §9.3 |
| `NET-12` Android security | Android design §4–§12 |
| `MOB-1`–`MOB-4` phone constraints | Android design scope and non-goals |
| `CWP-02` one prompt path | Reconciled plan §4, §6.5 |
| `CWP-06` durable memory | Memory design; `DPR-06` |
| `CWP-15a/15b` resource authority | Hardware design; `DWP-05` |
| `CWP-20` network control plane | Android design §7.4, §9 |
| `CWP-25` Android gate | Android design; `DPR-02` |
| `CWP-28` Story Tools | Reconciled plan §4, §7; `DPR-03` |
| `UD-01` per-chat model | Addendum; hardware design §8.3 |
| `UD-02` Android | Addendum; Android design |
| `UD-03` Story Tools | Addendum; `DWP-08` |
| `UD-04` narration | Addendum; `DWP-08` |
| `UD-05` Chat/TTS destination and Image-next staging | Addendum; `DWP-08`; reconciled plan §7.1 |
| `UD-06` memory policy | Addendum; memory design |
| `ODA-07` RAG commitment | RAG design; `DPR-07` |
| `ODA-09` hardware core | Hardware design; `DPR-08` |
| `ODA-08` OA-09 clarification | Addendum §3; all design packages |

## 11. New requirement identity

RAG has no row in the accepted 66. The owner decision makes it a committed product capability, but the accepted plan and reconciliation remain frozen.

Therefore:

- `DPR-07` is a **delta product requirement**, not a 67th reconciliation row.
- `CWP-28` remains conditional in Revision 3; the RAG design does not claim that `CWP-28` owns the full RAG lifecycle.
- A future reconciliation change must add a RAG requirement identity, source, owner, dependencies, and acceptance evidence before implementation planning can claim it.
- `RG-13` remains conditional; no RAG acceptance is claimed.
- `OA-10` product-scope question is resolved in favor of a committed local RAG; implementation and lifecycle authorization remain closed.

This avoids silently changing the accepted 66-row accounting.

## 12. Design-only acceptance

A design package is complete for review when it has:

- named boundaries and a single owner;
- explicit non-goals;
- data contracts and failure behavior;
- threat/failure analysis;
- traceability to the source requirement or owner decision;
- a test design that distinguishes static, runtime, benchmark, and acceptance evidence;
- unresolved decisions recorded rather than guessed;
- no implementation, migration, model load, or network action.

No `DWP` is an implementation package. No `DWP` acceptance is a product gate.

## 13. Open decisions

| ID | Decision | Needed before |
|---|---|---|
| `DA-01` | Physical store and migration path for state/events | Any schema or migration proposal |
| `DA-02` | Extraction model and structured-output contract | Memory automation design |
| `DA-03` | GPU telemetry provider order and packaging | Hardware implementation design |
| `DA-04` | RAG index/embedding store and parser order | RAG implementation design |
| `DA-05` | Android transport after authorized comparison/spikes | Android implementation design |
| `DA-06` | Conflict/automation policy by field (blocking) | Memory UX/acceptance |
| `DA-07` | Model compatibility matrix source | Suitability claims |
| `DA-08` | Story Tools capability map | `DWP-08` completion |

## 14. Risks

| Risk | Mitigation |
|---|---|
| Two memory stores diverge | Backend authority; frontend cache only |
| State history bloats prompts | Compact state cards, retrieval budgets, summaries |
| Small model invents state | Proposal-only, evidence validator, review for ambiguity |
| VRAM telemetry lies or is absent | Confidence labels, unknown state, warned attempt, no fabricated fit |
| Resource authority becomes advisory again | Single lease path and orphan diagnostics |
| RAG contaminates canon | Separate corpus and labels; no write path |
| RAG rights violation | Rights gate; unknown disabled by default |
| Android NAT failure | Relay fallback; do not promise direct connectivity |
| Pairing is phished | Desktop approval, short window, fingerprint |
| Remote access weakens offline guarantee | Off by default; per-request desktop policy |
| Design scope becomes implementation | `DWP` boundary and `OA-09`/`OA-06`/`OA-01` gates; DWP order is design-input order only |

## 15. Non-goals

- Implementing any design.
- Selecting an Android transport or database.
- Loading a model, building an index, or running a benchmark.
- Changing the accepted plan, 66-row reconciliation, or protected artifacts.
- Claiming RAG, Android, memory, hardware optimization, narration, or studios are operational.
- Treating planning acceptance as implementation authorization.

## 16. Review state

An independent architecture review of the pre-correction package is recorded at `docs/review/INDEPENDENT-ARCHITECTURE-REVIEW-2026-11-11.md` (`f1756fb5e6c21b8dce5ed30e677426f581ecfc732463a30eb4a630d34d8d5eee`). It found a completed first design pass, not approved architecture. Its targeted corrections F-01–F-09 have been applied in this package; a re-review and integration/acceptance planning remain pending.

Architecture approval and implementation authorization remain withheld. The package is **not** ready for implementation and does not pass any product or runtime gate.
