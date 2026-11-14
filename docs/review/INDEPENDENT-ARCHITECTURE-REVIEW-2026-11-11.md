# Independent architecture review — six-artifact design package

**Date:** 2026-11-11
**Type:** Independent technical review
**Status:** **First architecture-design pass reviewed; NOT approved architecture; NO implementation authorization**
**Scope:** Six new architecture/design artifacts, checked against Revision 3 planning baseline, Revision 3 owner decisions, existing implementation seams, migration risk, and dependency/staging risk.

## 0. Verdict

The package is a **completed first design pass**, not yet an approved architecture.

- **Do not discard it.**
- **Do not send it back for wholesale redesign.**
- **Do not implement from it.**
- Require the targeted corrections in §7 before treating any part as implementation-ready.

| Area | Independent assessment |
|---|---|
| Baseline preservation | Good; no plan/baseline edit found |
| Implementation authorization | Correctly withheld |
| Memory architecture | Promising, but terminology and migration gaps must be fixed |
| Resource/hardware architecture | Highest priority; directionally correct but lifecycle unproven |
| RAG architecture | Promising; lifecycle/storage/migration gaps must be fixed |
| Android security | Correctly left open; direction sound |
| Cross-system reconciliation | Needs explicit build/integration staging clarification |
| Product scope | **TTS destination correction required** |
| Verification record | Accurate as an integrity/coverage record, with wording fixes needed |

## 1. Review method and boundaries

This was a read-only architecture review:

- Compared the six design artifacts with each other.
- Checked them against `docs/review/CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md`, Revision 3.
- Checked them against `docs/review/OWNER-DECISION-ADDENDUM-REVISION-3-2026-11-11.md`.
- Inspected existing implementation seams for memory, RAG, resource/VRAM arbitration, model loading, prompt assembly, character voice/TTS, and shutdown/stream behavior.
- Did not change source, tests, runtime behavior, schemas, migrations, prompts, models, networks, settings, baselines, or protected artifacts.
- Did not run a model, benchmark, migration, network listener, pairing flow, index build, or hardware experiment.

## 2. Artifact integrity

Independently verified SHA-256 and sizes at review time:

| Artifact | SHA-256 | Size |
|---|---|---|
| `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md` | `79fed0c08cf1ba12190ddf831364fb4e46ec1600008cd8162d32103357cf821d` | 16286 B / 368 lines |
| `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md` | `1b697fb73b7c954251d86b39be5032a6dbe9a61e37ed59ea9fbefc1f6cc971d7` | 14455 B / 318 lines |
| `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md` | `82b6eef993bb74763c480a183603bc6ede99416b4d5db699e908979bd9ff12de` | 12709 B / 318 lines |
| `docs/architecture/RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md` | `a0d52c10f8cd4bf00770a362b3ece2744aa34d754a3797010f048e919d031bed` | 17421 B / 302 lines |
| `docs/architecture/RESEARCH-ANDROID-REMOTE-SECURITY-2026-11-11.md` | `c25e45b511b883c26d44ac8d3b61d66d6397a12a13f120f035a73f75474e33a2` | 15641 B / 318 lines |
| `docs/architecture/VERIFICATION-DESIGN-PACKAGE-2026-11-11.md` | `eeed8ab6c13e19f20543f22e5fb0d38a0258933025dc7151f0a0dd957a3f4ba5` | 7392 B |
| Revision 3 corrected master plan | `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` | 131755 B |
| Revision 3 owner-decision addendum | `d5fc5f95eecb2d8be8f5c903e4f05df9f14c4705534e08c6ffe4e71079938b3c` | 11564 B |

The six design hashes match the package verification record. The Revision 3 baseline hash matches the owner addendum.

## 3. What the package gets right

- It preserves Revision 3 instead of rewriting the accepted plan.
- It keeps RAG outside the frozen 66-row reconciliation as `DPR-07`, rather than silently creating a 67th row.
- It extends existing seams instead of inventing parallel systems:
  - memory candidates/ownership;
  - RAG chunk/index/retrieval;
  - resource allocator/telemetry/VRAM arbitration;
  - desktop policy for Android.
- It withholds implementation, runtime, migration, model-load, benchmark, and network authorization.
- It leaves Android transport, memory physical store, RAG store/embedding choice, telemetry provider order, and compatibility-matrix sourcing explicitly open.
- It correctly distinguishes roadmap order from modality authorization.

Those are the right disciplines for a first design pass.

## 4. Findings

Severity means architectural importance, not implementation urgency. Nothing below authorizes implementation.

| ID | Severity | Finding |
|---|---|---|
| F-01 | **High — product scope** | `UD-05` still places TTS as a separate studio/tab after Image; owner correction places character voice/TTS in Chat |
| F-02 | **High — terminology/privacy** | “Immutable event log” conflicts with required correction, supersession, retraction, deletion, and forget semantics |
| F-03 | **High — migration/safety** | Current memory persistence cannot support the proposed event/projection model without an explicit bridge |
| F-04 | **Highest — foundational risk** | Resource-authority lifecycle is directionally correct but incomplete and unproven for real concurrency, release, failure, and competing-consumer behavior |
| F-05 | **High — integration risk** | “One prompt assembler” is proposed, but multiple live prompt paths remain and no cutover is defined |
| F-06 | **Medium-High — lifecycle risk** | RAG design identifies the right lifecycle, but current RAG storage/index semantics and prompt seam need explicit migration and hardening |
| F-07 | **Medium — staging risk** | DWP dependency order is a design-input order, not an implementation/build/integration sequence |
| F-08 | **Low-Medium — residual security work** | Android direction is sound and correctly unselected, but device/session/data-at-rest details remain open |
| F-09 | **Low — record clarity** | Verification record is accurate but contains wording that can be misread as a current operational claim |

### F-01 — TTS studio destination must be corrected

**Evidence:**

- Owner addendum `UD-05`/`DPR-05` still records: `Chat → Image → TTS → Music → 3D → other modalities`.
- Reconciled plan `DWP-08`, `ADR-STU-01`, and verification wording say studio order is only a roadmap, but do not fix the TTS destination.
- Hardware design correctly treats `TTS/ASR` as a workload/resource consumer, but that is not the same as making TTS a product studio/tab.
- Existing code already treats TTS voice selection as part of the Character surface: `app-tts-catalog.go` describes `TTSVoiceCatalog` as “the engine-aware voice menu used by the Character tab.”

**Required correction:**

- Change the product destination to:
  - **Chat:** includes character voice/TTS functionality.
  - **Image:** next existing major capability to finalize.
  - **Music:** future dedicated studio.
  - **3D:** future dedicated studio.
  - **Other modalities:** evaluated later.
- Preserve the separate technical statement:
  - TTS/ASR remains a resource consumer under the future resource authority.
  - Character voice remains a Chat product destination, not a separate studio.
- Update `UD-05`, `DPR-05`, `DWP-08`, `ADR-STU-01`, hardware §8.1 wording, and any studio-order diagram/table that still implies a TTS studio/tab.

### F-02 — Replace “immutable” with auditable append-oriented history

**Evidence:**

- Memory design correctly requires:
  - append-only events;
  - correction/retraction as new events;
  - `supersedes_event_id`;
  - `disputed`, `superseded`, `user-locked` states;
  - forget as payload/index removal with a non-content tombstone;
  - as-of/effective/recorded-time queries.
- But the same package also uses the blunt phrase “immutable event log,” and reconciliation proposes `ADR-MEM-01: Immutable events plus revisioned projections`.

That wording invites the exact misreading the owner flagged: that nothing can ever be corrected or deleted.

**Required correction:**

- Replace “immutable” with **auditable, append-oriented history**.
- Define immutability narrowly:
  - no silent overwrite;
  - no hidden mutation;
  - every correction is explicit and attributable.
- Define the non-immutability operations explicitly:
  - correction;
  - supersession;
  - retcon/branch;
  - deletion/privacy removal;
  - forget/tombstone;
  - invalidation of stale derived state.
- Keep the memory design’s tombstone and referential-integrity behavior, but move deletion/forget policy from an open UX issue into a blocking pre-implementation requirement.

### F-03 — Current memory implementation needs an explicit migration bridge

The design correctly says current JSON remains the compatibility source until an authorized migration. But the gap is larger than “prose records versus canonical state.”

**Observed current seams:**

- `internal/backend/memory/service.go`:
  - CRUD operates directly on per-chat `memories.json`;
  - edits overwrite records in place;
  - deletes remove records without tombstones;
  - `saveMemories` does read-modify-write with no lock, transaction, revision check, or atomic replace;
  - the “backup” is currently a placeholder write of `{}` to `memories_backup.json`.
- `internal/backend/types/memory.go`:
  - memory IDs use second-resolution timestamps, creating a collision/duplication risk for rapid or batched capture.
- Memory filtering is title/summary-oriented and does not provide the event/projection retrieval contract.
- Frontend `memory-store.js` remains a second memory surface that can diverge.
- Candidate queue, continuity export/import, and memories are related but not unified into revisions, conflicts, undo, or as-of state.

**Required before any migration proposal:**

- A current-data inventory:
  - memories;
  - candidates;
  - continuity configs;
  - legacy owner/name-owned records;
  - frontend-cached records.
- Atomic durable-write semantics for any replacement store.
- Unique event/memory identity beyond second-resolution timestamps.
- Revision/compare-and-swap semantics.
- Tombstones and retention/expiry rules.
- A stale-cache guard so a forgotten fact cannot be re-ingested from cache.
- Export/import that preserves provenance, locks, conflicts, revisions, and deletion state—or fails visibly.
- Explicit treatment of `DA-M5` deletion/forget policy and `DA-M6` conflict/automation policy.

This does not require selecting the database now. It requires preventing a future migration from silently losing, duplicating, resurrecting, or misattributing memories.

### F-04 — Resource authority is the highest-risk item

The hardware design identifies the right target:

- vendor-neutral inventory;
- confidence-labeled telemetry;
- explicit 10% free-VRAM reserve;
- suitability as a report, not a guarantee;
- lease/admission/release;
- orphan diagnostics;
- verified versus unverified release;
- no fabricated fit when telemetry is unknown.

But the existing implementation shows why this must remain design-only and why the lifecycle needs more detail.

**Observed current seams and gaps:**

- `internal/resource/manager.go`:
  - thread-safe registry, but advisory: it records estimates and recommendations; it does not enforce admission.
- `internal/resource/allocator.go`:
  - 90%-of-reported-free heuristic;
  - eviction candidates;
  - no persistent lease ledger;
  - no atomic commit path;
  - no crash-recovery semantics.
- `app-resource.go`:
  - allocation advice is callable, but low confidence is only logged, not enforced as admission.
- `app-model.go`:
  - registers the chat model using memory estimates;
  - unregisters all tracked LLM models after stopping `llama-server`;
  - includes a `force_load` bypass path.
- `app-image-model.go`:
  - registers image models with fixed `4000/2000` MB values rather than measured backend/model-specific usage.
- `app-vram-arbitration.go`:
  - explicitly says the image engine is not consistently tracked by the registry;
  - unloads the competing model, waits a fixed 1.5 seconds, then rereads VRAM rather than verifying release through a release receipt/reconciliation loop.
- Shutdown/stream code:
  - orderly shutdown and stream cancellation exist;
  - per-chat concurrent inference, lease ownership across chats/jobs, and busy-model preemption policy are not defined.

**Missing lifecycle requirements:**

1. **Atomicity:**
   - stale fit rejection at commit;
   - no double-spending of the same free bytes;
   - lease expiry/renewal;
   - crash recovery for a dead owner/process.
2. **Concurrency:**
   - simultaneous Chat/Image/TTS/ASR/embedding/Android-job requests;
   - same-chat versus cross-chat contention;
   - queue/refuse/preempt policy;
   - starvation avoidance.
3. **Release verification:**
   - process exit is not proof of driver/VRAM reclamation;
   - define re-poll/reconciliation windows;
   - define “release unverified”;
   - do not promote fixed sleeps into release proofs.
4. **Failure and rollback:**
   - failed admission;
   - failed load after admission;
   - partial backend state;
   - orphaned subprocesses;
   - restoration of the prior safe model where applicable.
5. **Context/KV changes:**
   - context size, batch, cache type, quantization/bytes-per-element, attention variant, offload placement, encoder/VAE/control/audio overhead;
   - family-specific adapter or explicit `unknown`.
6. **Existing GPU allocations:**
   - non-Howl processes;
   - WDDM versus TCC reporting differences;
   - weak AMD/Intel live-used signals;
   - process-attribution limits.
7. **Competing consumers:**
   - TTS/ASR;
   - embeddings/index work;
   - Android-requested jobs;
   - future studios;
   - every consumer must be a lease or a diagnosed orphan.
8. **Busy-model policy:**
   - `busy cannot be evicted without drain` needs an operational definition:
     - drain boundary;
     - user-visible forced stop;
     - queued request behavior;
     - per-chat binding interaction.

None of this is a reason to abandon the design. It is the reason the resource authority must be treated as foundational and unproven.

### F-05 — Define the prompt-path cutover

The reconciled target is one prompt assembler. The current system has several prompt contributors/paths:

- `app-prompt.go:buildPromptMessages`;
- `internal/backend/prompt/builder.go:BuildPrompt`;
- lore resolution and lore prompt blocks;
- chat-service lore/memory snapshot handling;
- RAG retrieval paths;
- character-specific generation using frontend-built messages.

The designs acknowledge that memory/RAG injection is not yet proven exactly-once or authoritative. The missing piece is a cutover plan.

**Required:**

- Enumerate every current prompt contributor.
- Define canonical contributor order and token budget.
- Define exactly-once memory/RAG/state injection.
- Define legacy-path retirement.
- Prevent double prompting from backend plus frontend-built character messages.
- Prevent RAG/reference material from becoming canon.
- Prevent narration or character paths from forking prompt assembly.

`ADR-INT-01` should remain proposed until that cutover is specified.

### F-06 — RAG needs storage, deletion, retrieval, and seam hardening

The RAG direction is sound:

- local-only;
- structure-aware;
- hybrid lexical/semantic;
- citations;
- creative modes;
- canon/reference separation;
- degraded-state labeling.

**Observed current gaps:**

- `internal/backend/rag/service.go` and `document.go`:
  - working ingestion/chunk/index/retrieve seams;
  - no persistent hybrid index;
  - no embedding model identity/dimension/version persistence;
  - keyword fallback behavior can vary by corpus;
  - delete handling is oriented around clearing/rebuilding the in-memory index;
  - no durable deletion receipt/versioned rebuild manifest;
  - retrieval updates document `LastUsed` files while holding the service read lock, creating avoidable IO under lock and write amplification.
- `internal/backend/chat/service.go:retrieveRAG`:
  - accesses the last message without an empty-history guard.
- Scope model is world/scenario/chat and does not yet separate user canon drafts from external reference material.
- Rights metadata, edition/page mapping, structure parsing, and parser order remain open.

**Required:**

- Persist embedding/index identity and version.
- Define migration from current flat RAG files to manifests/versioned indexes.
- Make deletion durable and receipted, not merely an in-memory rebuild.
- Define stale-citation behavior after source edits.
- Separate user writing corpus from reference corpus.
- Harden retrieval IO/locking and empty-history behavior.
- Preserve `DA-R1` through `DA-R6` as pre-implementation decisions.

### F-07 — Clarify DWP order versus staged product/integration order

The DWP graph is acyclic and explicitly design-only. It does not authorize implementation or reopen `CWP-08`/`CWP-15`.

But a reader could mistake DWP input order for a build order and accidentally schedule four designs as four immediate implementation projects.

**Required clarification:**

- DWP order is the order for producing compatible designs.
- It is not a code/build/migration sequence.
- Preserve the intended product/integration staging:
  1. Finalize Chat behavior and prompt/state/memory seams.
  2. Finalize Image behavior and engine/resource registration.
  3. Then expand to further modalities and broader unification.
- Comprehensive UI unification remains a later artistic pass.
- It must not block targeted usability fixes.
- Future integration planning needs explicit integration-test gates:
  - Chat finalization;
  - Image finalization;
  - resource-authority enforcement;
  - memory/RAG prompt integration;
  - Android gateway policy checks.

### F-08 — Android is correctly uncommitted on transport, with residual items open

The Android research/design package has the right posture:

- desktop remains model host;
- no phone model;
- no cloud inference;
- off by default;
- no public listener;
- explicit pairing;
- mutual authentication;
- revocable scoped credentials;
- E2E at the application boundary;
- relay/rendezvous untrusted;
- no transport selected;
- no network experiment authorized.

**Residual pre-implementation issues already visible in the document:**

- Android-side cached chat/image data at rest.
- Session/key storage, rotation, expiry, and re-authentication defaults.
- Relay operator/metadata-retention policy.
- Pairing UX as a security control.
- Scope-specific authorization for chats, images, jobs, cancellation, model selection, and destination-policy enforcement.
- Offline/desktop-only behavior.

No change in direction is needed. These items should remain attached to `DWP-07` and future transport/security selection.

### F-09 — Verification-record wording

The verification record correctly establishes:

- planning acceptance versus implementation authorization;
- no runtime/model/network/migration activity;
- frozen baseline accounting;
- RAG delta identity;
- open decisions.

Two clarity fixes are needed:

1. Its final self-hash entry should be completed or explicitly marked as a packaging artifact, not left as “see final report.”
2. Statements such as “one resource authority covers…” should be labeled as the reconciled design target, not as current operational behavior.

## 5. Contradictions found

| Contradiction/tension | Resolution required |
|---|---|
| Owner addendum still sequences TTS as a separate studio after Image versus owner correction placing TTS/character voice in Chat | Correct `UD-05`/`DPR-05` destination |
| “Immutable events” versus correction, supersession, retcon, deletion, and forget | Adopt append-oriented/auditable terminology |
| “One resource authority” versus current advisory registry, unregistered image-engine state, TTS outside admission, and unregistered embeddings | Label authority as target; define lease/commit/release/orphan lifecycle |
| “One prompt assembler” versus multiple live prompt construction paths | Define contributor order and legacy cutover |
| Memory design’s no-migration rule versus future need to migrate current JSON/candidates/frontend cache | Require explicit current-data migration bridge |
| RAG canon/reference separation versus current scope model and prompt seams | Require corpus separation and durable lifecycle migration |
| Verification “covers TTS/ASR” versus product-surface uncertainty | Distinguish workload coverage from TTS product destination |
| DWP design-input order versus potential misreading as implementation schedule | Add explicit non-implementation and staging language |

No contradiction requires abandoning the package. All are resolvable with targeted revisions.

## 6. Missing requirements and decisions

These must remain open and must not be silently chosen later:

- TTS/character-voice product destination and Chat scope.
- Memory physical store and migration path.
- Memory deletion/forget policy, retention, export obligations, and stale-cache controls.
- Conflict/automation policy by field.
- Extraction model and structured-output contract.
- Resource provider order/packaging.
- Minimum reserve, runtime-overhead model, and compatibility-matrix source.
- Default low-end behavior: CPU fallback versus refusal.
- Busy-model drain/preemption and cross-consumer contention policy.
- Lease persistence, expiry, crash recovery, and orphan remediation.
- RAG physical store, embedding identity/version policy, parser order, reranker policy, citation granularity, rights UX.
- Android transport, rendezvous/relay operator, cipher/handshake suite, session lifetime, pairing UX, device data-at-rest.
- Story Tools capability map and ownership.
- Narration semantics, attribution, and implementation gate.
- Studio modality gates after the TTS destination correction.
- Integration-test staging for Chat, Image, resource authority, memory/RAG prompt path, and Android policy.

## 7. Required corrections before approval

No implementation work may begin on the basis of this package. The next revision should make only documentation corrections and explicitly remain design/research-only.

1. **Correct UD-05/TTS destination.**
   - Owner addendum, reconciled plan, hardware workload wording, verification coverage, and studio roadmap must consistently place character voice/TTS in Chat.

2. **Fix event-log terminology.**
   - Replace “immutable” with auditable append-oriented semantics.
   - Specify correction, supersession, retraction, deletion, forget, tombstone, retention, and as-of behavior.

3. **Specify resource-authority lifecycle.**
   - Atomic lease/commit.
   - Concurrency and contention policy.
   - Failure rollback.
   - Verified versus unverified release.
   - Orphan reconciliation.
   - KV/context/offload matrix or explicit unknown handling.
   - Competing-consumer coverage.

4. **Specify prompt-path cutover.**
   - Current contributors, canonical order, budgets, exactly-once injection, legacy retirement.

5. **Specify memory and RAG migration bridges.**
   - Current-data inventory.
   - Atomicity and identity.
   - Tombstones and deletion receipts.
   - Corpus separation.
   - Stale-cache and stale-citation controls.

6. **Clarify staging.**
   - DWP order is design-input order.
   - Preserve Chat → Image → expansion integration staging.
   - Keep broad UI unification later without blocking usability fixes.

7. **Clean up the verification record.**
   - Resolve the self-hash packaging note.
   - Label design targets as targets.

## 8. Recommended review outcome

- **Architecture package:** accepted as a completed first design pass.
- **Architecture approval:** **withheld**.
- **Implementation authorization:** **withheld**.
- **Runtime/network/model/migration authorization:** **withheld**.
- **Next step:** targeted documentation revision addressing §7.
- **Step after that:** integration and acceptance planning, not immediate implementation.
- **No redesign:** do not return the package for unstructured rewriting.

## 9. Authorization state after this review

| Question | Answer |
|---|---|
| Is this an approved architecture? | **No** |
| May implementation begin? | **No** |
| May runtime investigation begin? | **No** |
| May Android/RAG/memory/hardware changes begin? | **No** |
| May design corrections proceed? | **Yes, documentation-only** |
| Is the baseline changed by this review? | **No** |
| Are protected artifacts changed by this review? | **No** |

This review creates no new baseline, reconciliation row, product authorization, or implementation gate.
