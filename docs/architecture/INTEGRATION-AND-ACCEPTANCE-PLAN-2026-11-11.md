# Howl Chat — integration and acceptance plan

**Date:** 2026-11-11  
**Status:** Planning package; no architecture approval, implementation, migration, runtime, model, benchmark, or network authorization  
**Planning baseline:** Revision 3, frozen at `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874`  
**Corrected design package:** See §2; remains proposed  
**Independent re-review:** `docs/review/INDEPENDENT-ARCHITECTURE-RE-REVIEW-2026-11-11.md`, `88b12137baca5b27bb733de911c5fb69dc0fb21395123fe0e21dcbcc0fb413af`

## 1. Executive summary

The design package is ready for **integration and acceptance planning**, not for building. The independent re-review passed the corrected documentation and identified two minor items to carry forward:

- **N-01:** disposition the apparently unused legacy `buildPrompt` string builder.
- **N-02:** distinguish the existing lorebook-store ownership cutover from the future prompt-assembler cutover.

This plan turns those observations and the unresolved architecture decisions into a gated program:

1. research the field before committing to methods;
2. resolve owner-controlled and architecture decisions with evidence;
3. define integration seams and migration behavior;
4. define acceptance evidence before implementation;
5. obtain separate architecture and implementation authorization;
6. only then implement and validate.

No proposed method in this document is approved because it appears in the design package. The package supplies candidate architecture; the research dossiers and decision register must still test it against Howl's actual constraints.

The first product stage remains **Chat**, including character voice/TTS and Narration. **Image** is next. Music and 3D remain future dedicated studios. Broad UI unification remains later and does not block usability/accessibility fixes.

## 2. Baseline and review status

| Artifact | Identity | Status |
|---|---|---|
| Revision 3 plan | `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` | Frozen; not edited |
| Requirements reconciliation | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` | Frozen; 66 rows |
| Recovery addendum | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` | Frozen |
| Owner-decision addendum | `d228c4ac8a1acfabc7164ab5b277e9c20e8e2c33522554ced4896ec02b41e5ac` | Planning authority only |
| Memory design | `2e8ec935843bfd6b110feb16a626ce856cff60dd9ec1f570df0d21b21ef57650` | Proposed |
| Hardware design | `9348fa0dc297db356f8391f66318e98f8310b3ee6900b7116137ac23f852b0c3` | Proposed |
| RAG design | `2a80de38267131c7c4557729b845ac0bd43b3e345ad395020b93deb9e738e7e6` | Proposed |
| Android research/design | `c25e45b511b883c26d44ac8d3b61d66d6397a12a13f120f035a73f75474e33a2` | Research/design only |
| Reconciled plan | `58e33ed918e22628854a299380b9178b3cf61481b3a82cdebda94c459ee10040` | Proposed |
| Verification record | `0871bd4ace4dc48567f2d338310a464aae13698719ab61dcefaae0244cbdd1d1` | Package record |
| Prior independent review | `f1756fb5e6c21b8dce5ed30e677426f581ecfc732463a30eb4a630d34d8d5eee` | Pre-correction findings superseded |
| Independent re-review | `88b12137baca5b27bb733de911c5fb69dc0fb21395123fe0e21dcbcc0fb413af` | **PASS for planning only** |

The re-review's conclusion is narrow: the corrected documentation is coherent enough to plan integration and acceptance. It does **not** approve the architecture, select a physical store or transport, validate a model, or authorize a gate.

## 3. Review-observation disposition

### N-01 — legacy `buildPrompt`

**Observed:** Static search found `func (a *App) buildPrompt() string` in `app-prompt.go`, but no direct Go call site in the inspected tree. `internal/backend/prompt/builder.go` also defines `BuildPrompt`, with no direct call site found. The live generic path uses `app-prompt.go:buildPromptMessages`; the backend chat service contains a separate `Service.buildPrompt` used by its `SendMessage`/`SendMessageStream` methods, whose call-site reachability still needs confirmation.

**Disposition:** Both unused-looking builders are **retirement candidates**, not proven dead code. No deletion or modification is authorized. The prompt cutover must:

1. inventory build-tag, test, reflection, plugin, or dynamically linked callers before retirement;
2. preserve equivalence fixtures for the live generic and character paths;
3. retire each legacy builder only after the unified assembler passes the agreed fixtures;
4. keep the existing `buildPromptMessages` and chat-service path working until the cutover is explicitly switched;
5. record the removal evidence in the future integration package.

**Evidence class:** Observed source symbols and static call search; not a compiler/linker dead-code proof.

### N-02 — two different cutovers

The word “cutover” must not be reused without a qualifier.

| Term | Owner | What it changes | Current evidence |
|---|---|---|---|
| **Lorebook canonical-store cutover** | `internal/lorebookstore` and `app-lorebooks-store.go` | Transfers authored lorebook content ownership from the browser/legacy store to canonical files after explicit import | Marker file, import journal, `ErrNotCutover`, read-only listing, ownership migration, rollback bindings |
| **Prompt-assembler cutover** | Future `CWP-02`/one assembler | Replaces multiple prompt contributors and call sites with one precedence/budget/provenance contract | Design-only; current app and backend paths remain separate |

The lorebook cutover supplies canonical content to the prompt path. It does **not** unify prompt assembly, and prompt assembly does not own lorebook files. Acceptance criteria for each are separate.

### N-03 — frozen `UD-05` text

Revision 3 still contains the old `UD-05` status. That is controlled baseline preservation, not a current decision. The corrected owner addendum is the controlling planning record: character voice/TTS belongs to Chat, Image is next, and Music/3D are future dedicated studios.

## 4. Planning rules

1. **Do not build from a design slogan.** A design assertion becomes an approved choice only after the decision register has evidence and the owner/architecture gate is recorded.
2. **Do not replace working code by default.** Extending or adapting a working seam is preferred unless evidence shows replacement is necessary.
3. **Current behavior and target behavior are separate evidence.** A design target is not a runtime result.
4. **Unknown is a value.** Missing telemetry, ownership, retrieval evidence, or research is recorded as unknown, not filled by a plausible default.
5. **No silent mutation.** State, memory, lore, and citations retain provenance and an auditable correction path.
6. **One authority per concern.** No second prompt path, second memory authority, second resource allocator, or RAG write path.
7. **Research before selection.** Memory, RAG, hardware, local-model fit, Android, and Story Tools each require a dossier before their architecture decisions are treated as evidence-backed.
8. **Owner decisions are not inferred.** Product scope, privacy policy, deletion/retention, and transport/operator choices require explicit owner approval.
9. **Acceptance precedes implementation authorization.** Tests and failure criteria are designed before code is authorized, not after a feature exists.
10. **No runtime work without a separate grant.** `OA-09` covers runtime investigation; `OA-06` covers network work; `OA-01` and F3-C gates cover migration/storage.

## 5. Research-first program

### 5.1 Common research protocol

Every R-workstream must answer, in this order:

1. What methods have previously been attempted?
2. Under what conditions did they work?
3. What failed, and what limitations were observed?
4. What do actual users complain about?
5. What do users praise, and why?
6. What newer research or practice may improve the approach?
7. What benchmarks or independent evaluations exist?
8. Is it suitable for small local models?
9. What are the compute, memory, storage, latency, and context costs?
10. Can methods be combined into a Howl-specific approach that is better than any single method?

The dossier must distinguish:

- **Observed** — directly inspected code or authorized runtime evidence;
- **Source-derived** — primary documentation or peer-reviewed work;
- **Calculated** — arithmetic from stated specifications;
- **Inferred** — reasoned conclusion from stated evidence;
- **Unestablished** — no sufficient evidence;
- **Superseded** — older claim replaced by better evidence.

A benchmark improvement is not proof of a better writing experience. A result from a frontier model is not evidence for Gemma E2B, TinyLlama, or another constrained local model.

### 5.2 Workstreams

| ID | Workstream | Core questions | Minimum output | Decisions informed |
|---|---|---|---|---|
| **R-01** | Memory architecture | Episodic/semantic/procedural/working memory; hierarchy; extraction/consolidation; event/time; state transitions; contradiction and supersession; retrieval/relevance; compression; editable automatic memory; long-form roleplay; small-model memory | Comparative dossier, failure catalogue, small-model evidence, candidate architecture | `DA-01`, `DA-02`, `DA-06`, `DA-M1`–`DA-M6` |
| **R-02** | RAG and creative retrieval | Lexical/dense/hybrid; fusion/reranking; query expansion/decomposition; structure/hierarchy/graph; adaptive retrieval; context compression; chunking; citations; local embeddings/rerankers; creative writing use | Retrieval comparison, complaint/failure map, citation audit, small-model cost model | `DA-04`, `DA-R1`–`DA-R6` |
| **R-03** | Hardware/resource orchestration | Admission/estimation; VRAM reservation; KV/context; CPU/RAM/GPU offload; multi-GPU; switching/unloading; concurrency; vendor telemetry; unavailable telemetry | Provider/fit evidence matrix, contention model, failure/recovery test design | `DA-03`, `DA-H1`–`DA-H5` |
| **R-04** | Local model capability | Small instruction, quantized, dense/MoE, recurrent/state-space, multimodal, embedding/reranking models; deterministic support and prompt design | Capability matrix with named model/quant/hardware conditions; limits of generalization | Hardware suitability, memory extraction, RAG, prompt budget |
| **R-05** | Android remote security | Pairing/mutual auth; E2E; key/session rotation; NAT/rendezvous; P2P; relay; metadata; revocation; device identity; offline/reconnect | Threat model, transport comparison, privacy review, test design | `DA-05`, `DA-A1`–`DA-A7` |
| **R-06** | Story Tools and Narration | Creative workspace workflows; shared state/memory/RAG use; attributed AI takeover; reversibility; small-model usability | Capability map, ownership map, acceptance scenarios | `DA-08`, narration semantics |

R-01 and R-02 are the highest research risks because a plausible architecture can still produce poor continuity or poor retrieval. R-03 is the highest safety risk because a wrong admission decision can cause resource exhaustion or data loss.

### 5.3 Dossier template

Each dossier must include:

- research question and scope;
- methods and prior art;
- sources with date, publication type, and link;
- benchmark configuration: model, quantization, hardware, context, dataset, metric;
- observed results and reproduction status;
- failure modes and user complaints;
- small-model suitability;
- compute/memory/storage/latency costs;
- candidate approaches and combinations;
- recommendation only where evidence supports one;
- confidence, unknowns, and limitations;
- decisions it does and does not settle.

### 5.4 Runtime-experiment request

If a question cannot be answered statically, write an experiment request instead of running it. The request must state:

- question and hypothesis;
- why static evidence is insufficient;
- authorization needed (`OA-09`, `OA-06`, `OA-01`, or other);
- hardware, models, quantizations, context, and data;
- whether user content is involved and how it is avoided/minimized;
- procedure, metrics, baselines, and stop conditions;
- risks, resource ceilings, and cleanup;
- artifacts to retain and privacy treatment.

No experiment in this plan has been performed.

### 5.5 Research stop conditions

A dossier is not decision-ready when:

- it relies only on vendor claims;
- benchmark hardware/model conditions are absent;
- it generalizes frontier-model behavior to small local models;
- it cannot distinguish retrieval failure from generation failure;
- it treats user anecdotes as a controlled study;
- it proposes a store, model, parser, reranker, or transport without comparing alternatives;
- its recommendation would require user-content inspection or runtime work without authorization.

## 6. Unified decision register

The register is the single list for approval tracking. A design recommendation is not an owner decision.

### 6.1 Core decisions

| ID | Decision | Options/unknowns | Required evidence | Dependencies | Blocks approval | Blocks implementation | Runtime grant |
|---|---|---|---|---|---|---|---|
| `DA-01` | Physical store and migration path for canonical state/events | Transactional embedded DB vs versioned files; legacy JSON bridge | R-01; crash/atomicity evidence; F3-C/D-ROOT/D-STORAGE constraints | `DA-M1`; `DA-M5`; `DA-M6` | Yes | Yes | `OA-01` for migration |
| `DA-02` | Extraction model and structured-output contract | Model size, schema, validation, retry/abstain behavior | R-01/R-04; small-model extraction quality; failure rates | `DA-01`; `DA-M2` | Yes | Yes | `OA-09` for model tests |
| `DA-03` | GPU telemetry provider order and packaging | NVML, ADL/PDH, Level Zero, DXGI/WMI, process accounting | R-03; provider availability and confidence | `DA-H1` | Yes | Yes | `OA-09` for telemetry |
| `DA-04` | RAG index/embedding store and parser order | Embedded FTS/vector vs file hybrid; parser support order | R-02; deletion/atomicity/citation evidence | `DA-R1`–`DA-R3` | Yes | Yes | `OA-10` for indexing |
| `DA-05` | Android transport after comparative evidence | Direct LAN, WebRTC/ICE/TURN, WireGuard overlay, libp2p | R-05; NAT matrix, metadata, Android lifecycle | `DA-A1`–`DA-A7` | Yes | Yes | `OA-06`/`OA-09` |
| `DA-06` | Conflict/automation policy by field | Auto-commit thresholds; review; locks; ambiguity | R-01; user workflow study; failure taxonomy | `DA-02`; `DA-M3`; `DA-M6` | Yes | Yes | `OA-09` if measured |
| `DA-07` | Model compatibility matrix source | Curated evidence vs discovered metadata; maintenance | R-04; named model/hardware evidence | `DA-H4` | Yes | Yes | `OA-09` for runs |
| `DA-08` | Story Tools capability map | Existing controls vs new workflows; ownership | R-06; user workflow evidence | `DA-01`/`DA-04` consumer contracts | Yes | Yes | None for planning |

### 6.2 Memory decisions

| ID | Decision | Required evidence | Relationship |
|---|---|---|---|
| `DA-M1` | Physical store and migration mechanics | R-01; atomic staging; rollback | Alias of `DA-01` |
| `DA-M2` | Extraction schema and model contract | R-01/R-04; structured-output reliability | Alias of `DA-02` |
| `DA-M3` | Automation threshold by field | R-01; user review tolerance | Part of `DA-06` |
| `DA-M4` | Continuity branching/merge semantics | R-01; long-session and retcon scenarios | Blocks state/event schema |
| `DA-M5` | Deletion/forget policy, retention, tombstones, receipts, export | Privacy review; export/import proof; stale-cache proof | **Blocking pre-implementation** |
| `DA-M6` | Conflict/automation policy by field | R-01; conflict taxonomy; user control | **Blocking pre-implementation** |

### 6.3 Hardware decisions

| ID | Decision | Required evidence | Relationship |
|---|---|---|---|
| `DA-H1` | Provider order and dynamic-link packaging | R-03; provider availability | Part of `DA-03` |
| `DA-H2` | Minimum reserve and runtime-overhead model | Authorized measurements; otherwise unknown | `10%` is a default, not a measured constant |
| `DA-H3` | Multi-GPU support mode by backend | Backend evidence; no sum-as-pool assumption | Blocks multi-GPU claims |
| `DA-H4` | Compatibility-matrix source and maintenance | R-04; named evidence | Part of `DA-07` |
| `DA-H5` | Default low-end policy: CPU fallback vs refusal | Product/UX decision plus capability evidence | Owner-visible choice |

### 6.4 RAG decisions

| ID | Decision | Required evidence | Relationship |
|---|---|---|---|
| `DA-R1` | Embedded FTS/vector vs file-backed hybrid | R-02; atomic rebuild/deletion; scale | Part of `DA-04` |
| `DA-R2` | Embedding identity/dimension policy | R-02; local model matrix | Blocks index identity |
| `DA-R3` | Parser support order | R-02; structure fidelity | Part of `DA-04` |
| `DA-R4` | Reranker use and budget | R-02; latency/quality tradeoff | Optional; no default |
| `DA-R5` | Citation granularity | R-02; user verification workflow | Passage vs scene |
| `DA-R6` | Rights/consent UX for reference material | Legal/product review | Blocks external corpus |

### 6.5 Android decisions

| ID | Decision | Required evidence | Relationship |
|---|---|---|---|
| `DA-A1` | Direct LAN plus optional remote stack | R-05 | Overall architecture |
| `DA-A2` | WebRTC/ICE/TURN vs WireGuard vs libp2p | R-05; authorized network spikes | Part of `DA-05` |
| `DA-A3` | Rendezvous design/operator | R-05; metadata/abuse review | No default operator |
| `DA-A4` | Relay policy/limits | R-05; privacy/bandwidth | No content storage |
| `DA-A5` | Pairing UX/fingerprint | R-05; usability/security review | Security control |
| `DA-A6` | E2E handshake/cipher suite | R-05; protocol/library review | No speculative choice |
| `DA-A7` | Session lifetime/re-auth | R-05; threat model | No default interval |

### 6.6 Integration decisions introduced by this plan

| ID | Decision | Required evidence | Dependencies | Blocks |
|---|---|---|---|---|
| `DA-I1` | One prompt assembler contract and legacy retirement | Prompt inventory; golden fixtures; exactly-once proof | `CWP-02`; R-01/R-02 | Prompt implementation |
| `DA-I2` | Lorebook canonical content as prompt contributor, not prompt owner | Existing lorebook cutover evidence; ownership tests | `internal/lorebookstore`; `DA-I1` | Prompt/lore integration |
| `DA-I3` | Memory snapshot writer and authoritative read path | `memories_snapshot.json` writer evidence; migration tests | `DA-01`; `DA-I1` | Memory integration |
| `DA-I4` | Durable lease ledger and consumer drain boundaries | R-03; crash/contention design | `DA-03`; `CWP-15a/b` | Resource enforcement |
| `DA-I5` | RAG migration and citation lifecycle | R-02; deletion receipts; index manifest | `DA-04`; `DA-R1`–`DA-R6` | RAG integration |
| `DA-I6` | Story Tools state/memory/RAG consumer contract | R-06; shared authority map | `DA-08`; `DA-I1` | Story Tools implementation |
| `DA-I7` | Acceptance evidence ownership and review sign-off | Test design; evidence retention | All decisions | Any gate claim |
| `DA-I8` | Product staging and UI unification boundary | Owner staging; usability criteria | `UD-05`; `DA-08` | Sequencing |

## 7. Existing-to-target integration map

This section records **current static behavior** and the proposed transition. It does not prescribe replacement where adaptation is sufficient.

### 7.1 Prompt assembly

| Current seam | Current behavior | Target | Migration/rollback |
|---|---|---|---|
| `app-prompt.go:buildPromptMessages` | Live application-level system/world/character/scenario/lore/history assembly; lore resolved against latest user turn | One assembler with declared precedence, budgets, and provenance | Golden prompt fixtures; adapter first; legacy path retained for rollback until cutover proof |
| `app-prompt.go:buildPrompt` | ChatML string builder; no direct call site observed | Retire after confirmed dead, or adapt as a test fixture only | No deletion until call-site/build/test proof and prompt cutover acceptance |
| `internal/backend/prompt/builder.go:BuildPrompt` | Template/memory builder; no direct call site observed | Retire or fold into assembler behind one contract | Same as above; no second template authority |
| `internal/backend/chat/service.go:buildPrompt` | Backend chat-service prompt path; character/world/scenario/memory, then direct RAG string injection, then history/task | Contributor adapter; no direct RAG injection | Shadow/equivalence fixtures; preserve current behavior behind a switch |
| Frontend `chat-stream.js:buildCharacterSpecificMessage` | Builds character identity, recent events, group context, backend lore block, frontend character memories, user turn | Submit structured character identity/context; assembler owns order and injection | Character-path parity fixtures; don't remove until new path proves exact-once |
| `app-prompt.go:GenerateCharacterResponse` | Receives frontend-composed character message and sends it; comment explicitly warns against double-prompting | Single assembler request with contributor provenance | Preserve current contract until cutover, then retire frontend composition |
| `retrieveRAG` in chat service | Retrieves and writes `[RELEVANT DOCUMENTS]` directly into prompt | `EvidencePackage` contributor, labeled non-canon | Cut over through assembler; no direct string injection |
| Lore paths | General and character lore resolved by shared policy | Scoped lore contributor; canonical content source remains lorebook store | Keep lorebook cutover separate; test policy parity |

**Required prompt tests (future):** each contributor appears once; precedence matches contract; budget/truncation is deterministic; missing memory/lore/RAG is explicit; RAG evidence is labeled and cannot substitute for canon; character path does not duplicate identity; legacy paths are not called after retirement.

### 7.2 Memory and continuity

| Current seam | Current behavior | Target | Migration/rollback |
|---|---|---|---|
| `Chats/<chatID>/memories.json` | Read-modify-write JSON; in-place edit/delete; placeholder backup | Revisioned event/projection store behind one writer | Snapshot, staged import, count/hash verification, rollback |
| `candidates.json` | Separate review queue with its own read/write | Candidate proposals in the same revision/audit model | Preserve candidate status/provenance; quarantine ambiguity |
| `continuity_config.json` | Separate cross-scenario settings | Continuity policy stored with the continuity scope | Export/import with version; no silent default |
| `memories_snapshot.json` | Chat service reads it; prior assessment found no production writer | Explicit snapshot projection with an identified writer | `DA-I3` must resolve before migration; no silent fallback |
| Frontend `howl-chat-memories` | `localStorage` primary/mirror via `memory-store.js`; backend sync is best-effort | Cache/read model only; backend authority | Stale-cache guard, delete tombstones, explicit sync state |
| Character ownership | Stable `character_id`; legacy name/empty owner fails closed | Same, enforced everywhere | No name-based reassignment; quarantine |
| Event history | No general revision/as-of/conflict model | Auditable append-oriented history with correction/supersession/deletion | Migration receipts; no canonical fact invented from prose |

**Compatibility:** Existing exported memories and continuity data must remain readable or fail visibly. No migration may delete legacy files until rollback is proven.

### 7.3 RAG and creative retrieval

| Current seam | Current behavior | Target | Migration/rollback |
|---|---|---|---|
| `RAG/` document JSON + text + `.emb` | File-backed documents; in-memory index loaded from documents | Corpus/source/manifest model with persistent hybrid index | Side-by-side manifest build; old index retained; rollback |
| Embeddings | Local embedder, identity not persisted with index | Embedding identity/version/dimension on every index/result | New index identity; never mix vectors |
| Retrieval | Cosine + dedup + token cap; mutates `LastUsed` during read | Lexical/semantic hybrid, rerank/eval policy, read-only retrieval | Degraded mode; no write under read lock |
| Chat RAG | Direct prompt string injection; no observable empty-history guard | Evidence package; explicit empty-history refusal | Cut over through assembler; safety tests |
| Deletion | File removal + in-memory clear; no durable receipt | Durable deletion receipt, citation tombstone, cache invalidation | Two-phase delete; rollback only before receipt |
| Corpus | world/scenario/chat scopes; no rights separation | User writing/reference/craft/private-note corpora with rights | Quarantine unknown rights; separate indexes |

**Compatibility:** A new index must not silently serve old citations. Existing RAG files may remain until an authorized migration preserves identity and deletion state.

### 7.4 Hardware and resource orchestration

| Current seam | Current behavior | Target | Migration/rollback |
|---|---|---|---|
| `internal/resource` | Advisory snapshot and model registry; no durable leases | One authority with inventory, fit, leases, release, reconciliation | Enforce only behind an implementation gate; keep advisory read model during transition |
| Allocator | 90% of reported free VRAM; priority eviction suggestions | Named 10% reserve, runtime overhead, atomic reservation | Shadow fit reports first; then enforcement; rollback to advisory |
| Model loading | Registration after load; `force_load` bypass; estimates can be file/profile based | Reserve before load, commit-time recheck, audited override | No forced load without warning; preserve prior model if possible |
| Image engine | Registered after load with fixed `4000/2000` estimate; engine state read separately | One lease and live estimate; registered consumer | Register/observe before enforcement; reconcile orphans |
| TTS/ASR | Technical consumer outside complete admission | Lease per workload; no product-studio implication | Add to registry before enforcing; no voice UI change required |
| Unload/release | Fixed 1.5 s sleep before recheck; process exit not proof | Verified/unverified receipt with bounded reconciliation | Remove fixed-sleep assumption only with gate; preserve safe refusal |
| Context/KV/offload | Estimator limited; changes do not force one global fit decision | Any material change invalidates fit; re-plan/reload | No in-place admission patch |
| Existing allocations | Non-Howl usage appears in live free memory; registry does not know it | Treat as used; do not assume registry completeness | No auto-eviction of unregistered work |

**Compatibility:** Existing model selections and known configurations must remain loadable or produce an explicit fit report. Enforcement cannot silently disable an existing consumer.

### 7.5 Lorebook canonical store

| Current seam | Current behavior | Target | Migration/rollback |
|---|---|---|---|
| `internal/lorebookstore` | Canonical files, marker, journal, conflict reporting, owner migration | Keep as canonical authored-lore source | Do not conflate with prompt cutover |
| `app-lorebooks-store.go` | `SaveLorebook` fails before explicit import; `mergeCanonicalLore` overlays canonical text; chat snapshot retains selection | Prompt path consumes resolved canonical entries | Separate cutover evidence; no prompt change required to validate store |
| Legacy browser store | Derived cache/source for explicit import | Read-only compatibility source until retirement | Import/rollback bindings already exist; no silent migration |
| Chat snapshot | Selection and stale content can diverge from canonical | Selection only; canonical text resolved at prompt time | Stale/incomplete root handling; no resurrection |

**Acceptance:** canonical entry identity, ownership, conflict reporting, stale-root behavior, import rollback, and no silent browser-to-file overwrite. This is a precondition for reliable lore contribution, not for the prompt assembler's existence.

### 7.6 Story Tools and Narration

| Current seam | Current behavior | Target | Migration/rollback |
|---|---|---|---|
| Story Tools frontend | `story-tools-data.js` uses `MemoryStore`/localStorage as primary source | Consume canonical state/memory/RAG APIs | Adapter/read model; no parallel store |
| Memory editing | Story Tools and character memory UIs share localStorage semantics | Use revisioned memory API with visible diffs | Preserve current editor behavior while switching authority |
| Narration | No distinct Go narration path found in inspected tree; design is chat-only | Attributed AI takeover in same chat pipeline | Add as assembler consumer; no prompt fork |
| TTS | Character-tab voice catalog; resource consumer | Character voice inside Chat; resource lease | No separate studio or prompt path |

**Acceptance:** Story Tools edits preserve revisions and ownership; Narration messages are attributed and reversible; neither creates a second state store.

### 7.7 Android remote access

| Current seam | Current behavior | Target | Migration/rollback |
|---|---|---|---|
| Remote gateway | No remote Android gateway found in inspected code; model “pairing” symbols refer to multimodal projectors | Off-by-default desktop gateway with explicit device pairing | Greenfield; no data migration |
| Transport | Unselected | Compare direct LAN, WebRTC/ICE/TURN, WireGuard overlay, libp2p | No selection before R-05 and authorized spikes |
| Auth/E2E | Design-only | Mutual device auth, E2E, revocation, rotation | No pairing/listener/test yet |
| Policy | Desktop Local-Only Mode must remain authoritative | Per-request desktop policy; no remote bypass | Fresh opt-in; no legacy remote state |

**Important term collision:** existing `gguf` “pairing” means model/projector pairing. It must never be confused with Android device pairing in code, tests, logs, or acceptance language.

## 8. Migration and rollback planning

### 8.1 Migration principles

- Preserve source data until the destination is verified.
- Prefer a staged, side-by-side build over destructive replacement.
- Use a single writer and revision/CAS for mutable state.
- Never mix index identities or prompt contributors silently.
- Deletion is two-phase: tombstone/receipt first, payload cleanup second.
- Every cutover has an explicit rollback and a "no model/retrieval available" degraded state.
- A migration that loses, duplicates, resurrects, or misattributes data is failed even if the new store starts.

### 8.2 Memory migration phases (future, gated)

| Phase | Action | Evidence before next phase |
|---|---|---|
| M0 | Inventory `memories.json`, `candidates.json`, continuity config, snapshots, legacy ownership, frontend caches, backups | Manifest with counts, hashes, ownership classes, unknowns |
| M1 | Read-only mapping and quarantine of ambiguous records | Mapping report; no source mutation |
| M2 | Build staging store with durable IDs and provenance | Count/hash reconciliation; round-trip export/import |
| M3 | Verify projections, tombstones, conflicts, undo/as-of, stale-cache guard | Acceptance suite passes on synthetic fixtures |
| M4 | Cut over reads, then writes, under a single writer | Rollback tested; no dual writer |
| M5 | Retire legacy path only after stable operation and explicit approval | Receipt and owner sign-off |

**Rollback:** restore the prior JSON/snapshot source if staging or projection verification fails. **Blocked by:** `DA-01`, `DA-M1`, `DA-M5`, `DA-M6`, `OA-01`, F3-C Phase 0/D-ROOT/D-STORAGE.

### 8.3 RAG migration phases (future, gated)

| Phase | Action | Evidence before next phase |
|---|---|---|
| R0 | Inventory sources, chunks, embeddings, rights, scopes, caches, deletions | Manifest and identity report |
| R1 | Build new manifest/index beside old | Source hashes, dimensions, counts, sample retrieval |
| R2 | Verify citations, deletion receipts, corpus separation, stale cache | Acceptance suite passes |
| R3 | Switch retrieval by capability | No mixed vectors; degraded fallback labeled |
| R4 | Retire old index only after rollback window | Owner/implementation approval |

**Rollback:** retain old index and manifest until the new index passes; a failed build is `degraded`, not partially blended. **Blocked by:** `DA-04`, `DA-R1`–`DA-R6`, `OA-10`.

### 8.4 Prompt cutover phases (future, gated)

| Phase | Action | Evidence before next phase |
|---|---|---|
| P0 | Inventory all contributors/callers; confirm `buildPrompt`/`BuildPrompt` status | Static call graph, build-tag/test audit, ownership map |
| P1 | Define assembler contract and golden fixtures | Contributor identity, precedence, budget, exactly-once tests |
| P2 | Wrap current paths as adapters; shadow compare | Equivalence on generic/character/narration/lore/RAG fixtures |
| P3 | Switch one capability at a time | No duplicate injection; legacy path remains fallback |
| P4 | Retire legacy builders/paths | No live caller; tests pass; explicit approval |

**Rollback:** retain the current path behind a gated switch until all fixtures pass. **Blocked by:** `DA-I1`, `CWP-02` design authority.

### 8.5 Lorebook canonical-store migration (already separate)

The existing lorebook cutover is not re-planned here. Acceptance remains: explicit import, no write before cutover, canonical IDs, conflict reporting, stale-root preservation, owner migration, and rollback. Prompt planning may depend on its outputs but may not silently modify its store.

### 8.6 Resource enforcement transition (future, gated)

1. Observe and reconcile registered/unregistered consumers.
2. Produce shadow fit/lease reports without enforcement.
3. Add durable reservation semantics and diagnostics.
4. Enable enforcement for one consumer at a time.
5. Keep a visible forced-override path and advisory fallback.
6. Remove fixed-sleep/release assumptions only after verified receipts.

**Rollback:** return to advisory allocation without changing model selection or deleting model files. **Blocked by:** `DA-03`, `DA-I4`, `DA-H1`–`DA-H5`.

### 8.7 Android rollout (future, gated)

No existing remote state is migrated. A future rollout starts with a fresh, explicit pairing flow, off-by-default, and no public listener. Each transport candidate must pass R-05 and authorized network evidence before selection. **Blocked by:** `DA-05`, `DA-A1`–`DA-A7`, `OA-06`/`OA-09`.

## 9. Acceptance criteria framework

### 9.1 Evidence levels

| Level | Evidence | What it can prove | What it cannot prove |
|---|---|---|---|
| E0 | Static inspection | Symbols, calls, data shapes, obvious conflicts | Runtime behavior |
| E1 | Deterministic unit test | Pure logic, schemas, resolver policy | Integration or model quality |
| E2 | Integration test with fakes/fixtures | Seam behavior, rollback, exactly-once | Real model/driver behavior |
| E3 | Authorized runtime test | Actual local process/memory/telemetry behavior | Broad quality or security |
| E4 | Authorized security/network test | Threat controls under named configuration | Universal security |
| E5 | Authorized quality/benchmark study | Named quality/performance comparison | Generalization beyond configuration |

A unit test alone never establishes system correctness. A design approval never establishes implementation correctness.

### 9.2 Memory acceptance

Functional:
- automatic capture produces evidence-backed candidates, never canonical facts without a source;
- visible diff, edit, undo, lock, correction, retraction, and as-of view work;
- correction/supersession/deletion/forget change the projection without silently rewriting history;
- character ownership uses stable IDs; legacy/name-owned records fail closed.

Data/failure:
- concurrent writes cannot lose or overwrite a newer revision;
- crash/restart preserves committed state and never leaves a partial record;
- stale frontend cache cannot resurrect a deleted memory;
- export/import preserves provenance, tombstones, conflict state, and receipts;
- ambiguous imports are quarantined, not guessed.

Long-session:
- synthetic long-session fixture includes changing clothing/equipment/location/activity/relationships, corrections, retcon/branch, forget, restart, and export/import;
- exact turn count and scenario content belong to the approved test-design package, not this plan.

Evidence threshold: E1/E2 for logic and migration; E3 only under separate runtime authorization. No user chat content is used without explicit permission.

### 9.3 RAG acceptance

Functional:
- exact invented names resolve before thematic near-matches;
- semantic-only matches are labeled as such;
- every citation resolves to a passage hash and source span;
- creative suggestions are separated from cited evidence;
- reference material never changes canon.

Lifecycle/security:
- empty/invalid history is refused;
- retrieval is read-only under shared locks;
- deleted sources produce durable receipts and cannot reappear from cache;
- changed source hash/embedding identity marks citations stale;
- unknown-rights material is refused/quarantined;
- migrated indexes preserve corpus separation and identity.

Quality/performance:
- lexical fallback is labeled;
- token budget includes evidence and response reserves;
- retrieval quality is evaluated on realistic scene/book tasks, not only synthetic query pairs;
- latency/resource claims require E3/E5 evidence.

### 9.4 Prompt and lore acceptance

- every contributor appears exactly once;
- precedence and truncation match the contract;
- RAG is labeled non-canon and cannot substitute for canonical events;
- generic and character paths produce equivalent contributor content when inputs match;
- frontend composition does not duplicate backend identity/memory/lore;
- lorebook canonical content, selection, ownership, stale-root, and deleted-entry behavior remain correct;
- legacy builders are not called after retirement;
- narration uses the same assembler and remains attributed.

Evidence threshold: E2 fixtures before any runtime switch; E3 only after authorization.

### 9.5 Hardware/resource acceptance

- absent telemetry is `unknown`, not zero;
- 10% reserve is visible and cannot be spent by default;
- concurrent acquire/commit cannot double-spend;
- expired/crashed leases are reclaimed without freeing another lease;
- stale fit is rejected/replanned;
- failed load restores prior safe model or reports unloaded;
- release is verified by observation, not fixed sleep;
- context/KV/offload changes force re-fit;
- every consumer is a lease or reported orphan;
- busy models are drained/refused by policy, not silently killed;
- existing non-Howl allocations are treated as used;
- forced override is logged and does not become default.

Evidence threshold: E1/E2 for state machine and arithmetic; E3/E5 only under explicit hardware authorization.

### 9.6 Story Tools and Narration acceptance

- Story Tools reads/writes the shared memory/state authority, not a parallel store;
- edits preserve revisions and ownership;
- narration is attributed, reversible, chat-only, and does not fork prompt use;
- TTS remains character voice in Chat, not a separate studio;
- usability/accessibility fixes are not blocked by later UI unification.

Evidence threshold: E1/E2 for data contracts; E3 for user workflow only when authorized.

### 9.7 Android acceptance

- remote disabled produces no listener;
- unpaired device gets no chats/images;
- pairing requires desktop approval and short-lived ticket;
- revoked device cannot reconnect;
- relay sees no plaintext content;
- replay is rejected;
- remote cannot bypass Local-Only Mode;
- session expiry/re-auth works;
- desktop remains fully usable offline;
- logs contain no message content or long-term secrets.

Evidence threshold: E2 for protocol/state design; E4 only after explicit network/security authorization.

## 10. Dependency graph and sequencing

### 10.1 Research-to-release graph

```text
R-01..R-06 dossiers
        │
        ▼
DA-01..DA-08 + DA-M/DA-H/DA-R/DA-A + DA-I decisions
        │
        ▼
Corrected architecture reconciliation
        │
        ▼
Integration contracts + migration/rollback plan
        │
        ▼
Acceptance-test design + evidence ownership
        │
        ▼
Owner/architecture approval
        │
        ▼
Implementation authorization
        │
        ▼
Implementation
        │
        ▼
Migration + authorized runtime validation
        │
        ▼
Release acceptance
```

This graph is a dependency graph, not a schedule. It does not resurrect the `DWP` numbers as build order.

### 10.2 Product staging

| Stage | Scope | Entry condition |
|---|---|---|
| **Chat** | Prompt cutover, state/memory, lore integration, resource authority, TTS/ASR and Narration inside Chat, local-only boundary | Chat research/decisions and acceptance design accepted; separate implementation authorization |
| **Image** | Image generation/editing as next major capability, registered with the resource authority | Chat stage gates; image design/quality/resource gates accepted |
| **Expansion** | Music and 3D as future dedicated studios; other modalities later | Separate design, evidence, owner approval |
| **Later UI unification** | Shared presentation across Chat, Image, Story Tools, and future studios | Does not block usability, accessibility, or functional corrections |

DWP order remains **design-input order**. No DWP is an implementation package.

### 10.3 Gate sequence

| Gate | Purpose | Current state |
|---|---|---|
| G0 | Baseline/package verification and independent review | PASS for planning only |
| G1 | Research plan/dossier acceptance | Not started |
| G2 | Architecture decisions resolved with evidence | Open |
| G3 | Design reconciliation and integration contracts | Not complete |
| G4 | Acceptance-test design accepted | Not complete |
| G5 | Owner architecture approval | Not granted |
| G6 | Implementation authorization | Not granted |
| G7 | Implementation complete under scope | Not started |
| G8 | Migration/runtime validation | Not authorized |
| G9 | Release acceptance | Not started |

No gate beyond G0 is passed. G0 is not architecture approval.

## 11. Authorization gates

### 11.1 Current authorization

| Activity | Authorized? |
|---|---|
| Read-only source/document inspection | Yes |
| External research | Yes, subject to source quality |
| Planning documents | Yes |
| Source code changes | **No** |
| Migrations/persistent data changes | **No** |
| Model loads/benchmarks | **No** |
| Runtime/telemetry experiments | **No** |
| Network listeners/pairing/packet capture | **No** |
| Acceptance gate claims | **No** |

### 11.2 Required grants for future work

| Work | Grant |
|---|---|
| Runtime investigation, model load, benchmark | `OA-09` |
| Android network/listener/pairing experiment | `OA-06` and `OA-09` as applicable |
| Durable-state migration | `OA-01` plus F3-C Phase 0/D-ROOT/D-STORAGE |
| RAG indexing/runtime | `OA-10` and relevant storage/rights approval |
| Implementation | Separate implementation authorization after G5/G6 |
| Release | Separate release acceptance; no implied deployment |

A successful document review does not imply any of these grants.

## 12. Contradictions and risks discovered

| ID | Finding | Planning consequence |
|---|---|---|
| C-01 | Verification record says F-01–F-09 had not received independent re-review; the re-review now exists. | Do not edit the old verification record; cite the re-review as the controlling status. |
| C-02 | `memories_snapshot.json` is read by chat service, but prior assessment found no production writer. | `DA-I3` must establish the writer/authority before memory migration. |
| C-03 | `buildPrompt` and `prompt.Builder.BuildPrompt` have no direct call site, but no compiler/linker proof of dead code. | Retire only after P0 audit and cutover evidence. |
| C-04 | Generic app path, backend chat-service path, and frontend character path all build prompts differently. | `DA-I1` is a real cutover, not a refactor label. |
| C-05 | Lorebook “cutover” and prompt “cutover” mean different migrations. | Use qualified terms everywhere; separate acceptance. |
| C-06 | RAG retrieval mutates document metadata during read; direct prompt injection bypasses an assembler. | `DA-I5` and `DA-I1` must be designed together; no runtime index build yet. |
| C-07 | Resource registry is advisory; image/TTS consumers are inconsistently registered; release uses a fixed wait. | `DA-I4` covers durable leases, consumer drain, and verified release. |
| C-08 | Existing Android “pairing” symbols are model/projector pairing, not device pairing. | Mandatory terminology guard in Android design/tests. |
| C-09 | Frozen Revision 3 still displays old `UD-05`; addendum resolves it. | Preserve baseline; do not treat frozen text as current decision. |
| C-10 | No R-01–R-06 dossiers exist yet. | The plan establishes research gates; it does not claim methods are validated. |
| C-11 | Story Tools currently uses localStorage/MemoryStore as primary source. | Migration must move it to shared authority without losing existing entries. |
| C-12 | Later UI unification could be mistaken for a prerequisite. | It is explicitly later; usability/accessibility fixes continue. |

## 13. Recommended next actions

1. **Accept this plan as the planning baseline only.** It does not approve architecture or implementation.
2. **Commission R-01 and R-02 first.** They resolve the largest product-quality unknowns and feed `DA-01`/`DA-02`/`DA-06` and `DA-04`.
3. **Commission R-03 in parallel.** It produces the telemetry, contention, reserve, and release evidence needed before resource enforcement can be safe.
4. **Resolve blocking memory decisions:** `DA-M5` (deletion/forget/tombstones/receipts/export) and `DA-M6` (conflict/automation) before any migration design is approved.
5. **Resolve `DA-I3`:** identify the `memories_snapshot.json` writer and the authoritative memory read path.
6. **Commission R-05 but do not select a transport.** It must compare direct LAN, WebRTC/ICE/TURN, WireGuard overlay, and libp2p under the threat model; no listener or pairing test is authorized.
7. **Define prompt fixtures and lorebook-store acceptance separately.** Then produce the unified assembler contract.
8. **Produce acceptance-test designs before implementation authorization.** Include synthetic long-session, failure, contention, deletion, and migration scenarios.
9. **Re-review the integration/acceptance package.** The current PASS is for planning readiness, not architecture approval.
10. **Return to the owner for architecture and gate decisions.** No implementation starts from this plan.

## 14. Source and evidence register

### 14.1 Design and review sources

The sources below are documents, not runtime evidence. Their evidence class is **Source-derived** unless a row says otherwise.

| Source | Use |
|---|---|
| `docs/review/INDEPENDENT-ARCHITECTURE-RE-REVIEW-2026-11-11.md` | PASS, N-01/N-02/N-03, remaining blockers |
| `docs/review/INDEPENDENT-ARCHITECTURE-REVIEW-2026-11-11.md` | Pre-correction findings, superseded where corrected |
| `docs/review/OWNER-DECISION-ADDENDUM-REVISION-3-2026-11-11.md` | Product scope, TTS destination, delta requirements |
| `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md` | State/event semantics, migration bridge |
| `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md` | Inventory, reserve, leases, release, failure modes |
| `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md` | Corpus, hybrid retrieval, citations, migration |
| `docs/architecture/RESEARCH-ANDROID-REMOTE-SECURITY-2026-11-11.md` | Threat model, transport comparison, Android acceptance |
| `docs/architecture/RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md` | Ownership, DWP boundaries, product staging |
| `docs/architecture/VERIFICATION-DESIGN-PACKAGE-2026-11-11.md` | Package integrity record; stale re-review line noted as C-01 |

External sources already recorded by the design package (not independently reproduced here):

- NVIDIA NVML API and WDDM memory-reporting caveat.
- AMD ADL/PDH and Windows GPU-memory counter discussion.
- Intel Level Zero Sysman as a provider candidate.
- WireGuard known limitations and protocol documentation.
- RFC 8656 and RFC 5128 for TURN/NAT traversal.
- libp2p hole-punching and Circuit Relay v2 documentation.
- Tailscale NAT-traversal discussion.

These support candidate comparisons only. They do not validate Howl behavior, model suitability, or a transport choice.

### 14.2 Code seams inspected

Static inspection only; no user chat content was read as evidence.

| File | Why inspected | SHA-256 |
|---|---|---|
| `app-prompt.go` | Generic/character prompt paths, lore resolution, legacy `buildPrompt` | `93c9e6b8759936194a06ee5a49a3c2dd1f16054b20670738c45c42c5940d4128` |
| `app-chat.go` | Generic send/save path, `buildPromptMessages` use | `f08037a8e3aaf32d5e70cad8fc25ecedefeddc3bd1a34a07e6a33d10256fea05` |
| `app-memory.go` | Memory CRUD bindings and identity diagnostics | `7e7a8e83724f3c06e11a22193a1cb246ac89608669b51d15acdaaf1cbb427f60` |
| `app-resource.go` | Resource snapshot, registry, allocation bindings | `6c1c5d2ad59de229b01f5e3fc62e1c64c1995cde30bdaf6750871dd5598b8a38` |
| `app-model.go` | Model load, force-load bypass, unload/registry | `6dd43a1fb0aada7a7283df32463b00aa19510e237cce3b2f75980734cee6fd12` |
| `app-vram-arbitration.go` | Fixed release wait, fit checks, image/chat arbitration | `a266d04b2070bbfb1b5c9f05329c417446202e697ec9e7d7acf61fffe743d4ac` |
| `app-image-model.go` | Image load and fixed resource registration | `6feb61ad50ef283e270bc521235a8fda5145e62d52acc6200aa441c638577de7` |
| `app-image-generate.go` | Image job/lazy-load seam | `87cccad8982656ebadfbb269d9b390d31024d0c4068536ba1708f01ae0e4d770` |
| `app-lorebooks-store.go` | Lorebook canonical-store cutover and overlay | `52a82d6ed053d538182794ea0f8f7287e357756eaa611c165e421b8b2438ea0a` |
| `app-tts-catalog.go` | Character-tab TTS voice catalog | `8d88518c8fd5ac2f4648f1191b21af9c03b1fee22c0acf813907c1a608c55db5` |
| `internal/backend/prompt/builder.go` | Legacy template builder and memory injection | `558c1dcd5ba7a36bb6ced65001ceb9e712f369e5fd8b7470a82608228100bf64` |
| `internal/backend/chat/service.go` | Backend prompt, memory snapshot, RAG path; call-site reachability still to confirm | `1f068fc15f0c89c6ee46a26a202e84f5dda505608c5acfe454f75c7585d4ae8c` |
| `internal/backend/chat/persistence.go` | Chat/memory snapshot persistence | `ad9663a48278cff0ee35eb48d28799702cdbfa7754f9f71b7541a3a6762a87c9` |
| `internal/backend/memory/service.go` | Memory storage, save/backup, IDs | `5f6c7d7824099e6999997afdd6b793bda422b360e73529c17a6f33e4f8837ed2` |
| `internal/backend/memory/candidates.go` | Candidate queue persistence | `fe72182aef09498b08f2aec33382599e78ae614b3b49ad21e750930a4a5750f7` |
| `internal/backend/memory/continuity.go` | Continuity config persistence | `c4c9ce17e3e3d2fee7228d7f7c3c4253e224be961cdbd8cd8aede54d51a425e7` |
| `internal/backend/rag/service.go` | In-memory index, retrieval, delete, LastUsed | `ac37f1c173b72e54d6fdf2b75dadebfe2bd8922f8182a7e70d7d5c7f1938c918` |
| `internal/backend/rag/document.go` | Flat-file document/chunk persistence | `2899641bb49f389ff2925f53bff937dae6b253fb169ee32f9d7c10608ef18c9b` |
| `internal/backend/rag/types.go` | Document/chunk/query schema | `75c9bd44487ec2328f0888a7990fbd5d91205da1b8c773436e6db6d43aecdf76` |
| `internal/lorebookstore/store.go` | Canonical lorebook marker/journal/write rules | `646e915a8e99afa6a14c46130ed8fd1f40e8e3e955adc4de14d6c238ab5fccf9` |
| `internal/resource/allocator.go` | 90% heuristic and eviction advice | `ffed5833a2bbab054c2254c765da469e85c6ee4013be9ca6e4132ff034c93398` |
| `internal/resource/manager.go` | Advisory registry and hardware snapshot | `34261573e18c84a3c9d0cd5e3c3cd63cd35edbd965935ff9f3bc29ff51e392ef` |
| `internal/resource/types.go` | Resource schema | `99be5d86468ddbeb00a9a6155dd0ee9b95525617b82e54de7fe7b66e1a40f9fc` |
| `frontend/dist/js/memory-store.js` | localStorage primary/mirror and backend sync | `78a8f6b8c76d435ccdf4b58614b5692e7bad19d8c9cc6505436973d24176a4cd` |
| `frontend/dist/js/chat-stream.js` | Frontend character prompt composition and sends | `2106be99e088476d58274b10d610034167bb3a6168a064371c5a2367656c377a` |
| `frontend/dist/js/character-memories.js` | Character memory UI reads cache | `f1b02bb846d90e91a05595ab4c212d09e9375d2e1fe3ae61c7fec5c45ee2a6e8` |

## 15. Final authorization statement

This document is a **planning package**. It does not:

- approve the architecture;
- choose a physical store, embedding model, parser, telemetry provider, Android transport, relay, or cipher suite;
- authorize source changes, migrations, model loads, benchmarks, runtime measurements, listeners, pairing, packet capture, or network tests;
- pass any implementation or release gate;
- modify Revision 3, the reconciliation, the recovery addendum, the independent reviews, or the existing design package.

The next legitimate state is **research and decision work**, followed by integration/acceptance review and explicit owner authorization. It is not implementation.
