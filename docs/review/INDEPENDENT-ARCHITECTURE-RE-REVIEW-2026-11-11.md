# Independent targeted re-review — corrected Howl architecture package

**Date:** 2026-11-11
**Review type:** Independent verification of targeted documentation corrections
**Authorization:** Read-only inspection and review only
**Status:** Previous independent review preserved unchanged; this is a separate re-review
**Implementation authorization:** NONE
**Runtime/network/migration authorization:** NONE

## 1. Executive verdict

**PASS — READY FOR INTEGRATION AND ACCEPTANCE PLANNING**

The corrected documentation package substantively resolves the prior review’s F-01–F-09 correction areas, remains internally coherent, and does not introduce a new blocking architectural problem.

This PASS means only that the corrected package is sufficiently coherent to proceed to integration and acceptance planning. It does **not** approve the architecture, pass any implementation gate, authorize code/runtime/migration/model/network work, or establish that any behavior works at runtime.

Two minor non-blocking documentation clarifications should be carried into integration planning:

- N-01: explicitly disposition the apparently unused legacy `buildPrompt` string builder alongside the prompt-path cutover inventory.
- N-02: disambiguate the existing lorebook-store “cutover” term from the future prompt-assembler cutover.

## 2. Scope and artifacts inspected

Re-review basis:

- Previous independent review: `docs/review/INDEPENDENT-ARCHITECTURE-REVIEW-2026-11-11.md`
- Corrected owner addendum: `docs/review/OWNER-DECISION-ADDENDUM-REVISION-3-2026-11-11.md`
- Corrected memory design: `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md`
- Corrected hardware/resource design: `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`
- Corrected RAG design: `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md`
- Android research/design: `docs/architecture/RESEARCH-ANDROID-REMOTE-SECURITY-2026-11-11.md`
- Corrected reconciliation: `docs/architecture/RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md`
- Corrected verification record: `docs/architecture/VERIFICATION-DESIGN-PACKAGE-2026-11-11.md`
- External digest sidecar: `docs/architecture/VERIFICATION-DESIGN-PACKAGE-2026-11-11.md.sha256`

Implementation seams were inspected read-only where necessary to test design claims about current behavior, including memory persistence/identity, RAG retrieval/index behavior, resource registration/allocation/VRAM arbitration, model-load bypass behavior, prompt construction paths, lore canonicalization, character-message paths, and narration surfaces. Local user chat content was not used as evidence.

No document or source file was modified. No test, build, migration, model load, benchmark, listener, pairing flow, packet capture, index build, or hardware experiment was performed.

## 3. Hash and frozen-baseline verification

Independently recomputed digests match the expected identities supplied for this re-review.

| Artifact | Expected SHA-256 | Recomputed | Match |
|---|---|---|---|
| Owner-decision addendum, Revision 3 | `d228c4ac8a1acfabc7164ab5b277e9c20e8e2c33522554ced4896ec02b41e5ac` | `d228c4ac8a1acfabc7164ab5b277e9c20e8e2c33522554ced4896ec02b41e5ac` | Yes |
| Memory design | `2e8ec935843bfd6b110feb16a626ce856cff60dd9ec1f570df0d21b21ef57650` | `2e8ec935843bfd6b110feb16a626ce856cff60dd9ec1f570df0d21b21ef57650` | Yes |
| Hardware/resource design | `9348fa0dc297db356f8391f66318e98f8310b3ee6900b7116137ac23f852b0c3` | `9348fa0dc297db356f8391f66318e98f8310b3ee6900b7116137ac23f852b0c3` | Yes |
| RAG design | `2a80de38267131c7c4557729b845ac0bd43b3e345ad395020b93deb9e738e7e6` | `2a80de38267131c7c4557729b845ac0bd43b3e345ad395020b93deb9e738e7e6` | Yes |
| Android research/design | `c25e45b511b883c26d44ac8d3b61d66d6397a12a13f120f035a73f75474e33a2` | `c25e45b511b883c26d44ac8d3b61d66d6397a12a13f120f035a73f75474e33a2` | Yes |
| Reconciled plan | `58e33ed918e22628854a299380b9178b3cf61481b3a82cdebda94c459ee10040` | `58e33ed918e22628854a299380b9178b3cf61481b3a82cdebda94c459ee10040` | Yes |
| Verification record | `0871bd4ace4dc48567f2d338310a464aae13698719ab61dcefaae0244cbdd1d1` | `0871bd4ace4dc48567f2d338310a464aae13698719ab61dcefaae0244cbdd1d1` | Yes |

Observed sizes and line counts:

| Artifact | Bytes | Lines |
|---|---|---|
| Owner-decision addendum | 11939 | 141 |
| Memory design | 20764 | 421 |
| Hardware/resource design | 20240 | 392 |
| RAG design | 16394 | 361 |
| Android research/design | 15641 | 318 |
| Reconciled plan | 22008 | 354 |
| Verification record | 13494 | 178 |

The external sidecar contains the verification record digest and matches the recomputed verification-record hash. The verification record no longer embeds a self-hash.

Frozen baseline verification:

| Frozen artifact | Expected | Recomputed | Match |
|---|---|---|---|
| Corrected master plan, Revision 3 | `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` | `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` | Yes |
| Requirements reconciliation | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` | Yes |
| Recovery owner addendum | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` | Yes |

The prior independent review remains byte-identical at:

- `f1756fb5e6c21b8dce5ed30e677426f581ecfc732463a30eb4a630d34d8d5eee`

That review examined the pre-correction package; its F-01–F-09 findings are therefore marked **Superseded** where the corrections resolve them, not deleted or rewritten.

## 4. Disposition table

| ID | Prior finding | Disposition |
|---|---|---|
| F-01 | TTS destination/product staging | **PASS** |
| F-02 | Event-history semantics | **PASS** |
| F-03 | Memory migration bridge | **PASS** |
| F-04 | Resource lifecycle | **PASS** |
| F-05 | Prompt-path cutover | **PASS**, with minor N-01/N-02 clarifications |
| F-06 | RAG migration/retrieval | **PASS** |
| F-07 | DWP order/product staging | **PASS** |
| F-08 | Android residuals | **PASS** |
| F-09 | Verification clarity | **PASS** |

No item is marked FAIL or UNVERIFIABLE. The PASS dispositions below are documentation dispositions, not runtime findings.

### F-01 — TTS destination and product staging: PASS

**Observed:** The corrected addendum now states that character voice/TTS belongs to Chat, Image is next, Music and 3D are future dedicated studios, and TTS/ASR remains a technical resource consumer rather than a product studio. The reconciled plan carries the same destination into product staging, `DWP-08`, and `ADR-STU-01`. The hardware design explicitly calls TTS/ASR a resource consumer, not a product studio. The verification record separates the prior TTS sequencing from the corrected destination.

**Source-derived:** No residual corrected-package language routes TTS after Image as a separate studio. Older UD-05 “unresolved” wording remains only in frozen or historical records that the package is required to preserve byte-identical; the corrected addendum is the controlling planning record for destination and staging.

**Inferred:** The product/resource distinction is now coherent: Chat owns the voice destination, while the future resource authority governs TTS/ASR workload admission.

**Remaining boundary:** The correction authorizes no modality. Image, Music, and 3D still require separate design, evidence, and authorization.

### F-02 — Event-history semantics: PASS

**Observed:** Memory §1, §3, and §4.2 replace event-log absolutism with auditable append-oriented history. The design defines correction, supersession, retraction, retcon/branch, deletion/privacy removal, forget/tombstone, derived-evidence invalidation, and effective/recorded-time as-of reconstruction. `ADR-MEM-01` uses the same terminology.

**Source-derived:** The semantics can coexist with auditability at the documentation level because normal history remains attributable, while forget removes payload/index content and retains a non-content tombstone and deletion receipt for referential integrity.

**Unestablished:** Whether any particular future tombstone/metadata retention rule satisfies privacy or export obligations. That is correctly left to blocking `DA-M5`.

**Inferred:** The prior terminology contradiction is resolved. The remaining privacy question is a blocking pre-implementation decision, not a reason to withhold integration planning.

### F-03 — Memory migration bridge: PASS

**Observed:** Memory §9.1 inventories the existing surfaces and constraints: per-chat memories, candidates, continuity configuration, legacy/name-owned ownership, frontend caches, placeholder backups, imports/exports, second-resolution IDs, read-modify-write persistence, and absent revision/conflict/undo/as-of modeling. It then requires durable IDs, atomic staging, single-writer compare-and-swap, tombstones, deletion receipts, locks/conflicts/attribution preservation, stale-cache guards, rollback, quarantine of ambiguous legacy records, and continued blocking behind `OA-01`, F3-C Phase 0/D-ROOT, and D-STORAGE.

**Source-derived:** Representative implementation seams confirm the bridge’s premises: `saveMemories` performs read-modify-write with a placeholder backup and no lock/revision/atomic replace; memory IDs use second-resolution timestamps; candidate storage is a separate JSON queue; continuity configuration is a separate chat file; frontend keys for memories and candidates exist as a separate renderer surface.

**Inferred:** This is a genuine migration-safety design, not a relabeling of current JSON. It does not authorize migration or select the physical store.

**Remaining boundary:** `DA-M5` and `DA-M6` are correctly identified as blocking. `DA-M1`–`DA-M4` also remain open.

### F-04 — Resource lifecycle: PASS

**Observed:** Hardware §8.4–§8.7 now specifies atomic acquire/commit, idempotency keys, bounded lease lifetime and renewal, expiration/reclamation, crash reconciliation, commit-time revision checks, concurrency/contention behavior, busy-model draining, failure rollback, verified versus unverified release, live context/KV/offload invalidation, existing non-Howl allocations, and lease-or-orphan coverage for chat, image, TTS/ASR, embeddings, Android-requested work, and future studios. Failure table §11 and acceptance design §13 repeat the same guarantees as unexecuted test designs.

**Source-derived:** The design’s treatment matches the observed current limitations it is meant to replace: advisory allocation, fixed image-registration estimates, image-engine registry gaps, fixed post-unload sleep, force-load bypass behavior, and process exit being insufficient proof of driver/VRAM reclamation.

**Inferred:** The proposal addresses the race and release-verification concerns at the documentation level: reservations are visible under the authority’s write lock, stale fits are rejected or re-planned, unverified release retains the ledger entry, failed loads reconcile actual usage, and existing allocations are treated as used rather than assumed absent.

**Unestablished:** No concurrency, release, telemetry, KV, compatibility, benchmark, or cross-consumer behavior was demonstrated. Consumer-specific drain boundaries and durable lease-record details remain integration-design work.

**Remaining boundary:** This PASS does not make the resource authority implemented, measured, or safe to enforce. It makes the lifecycle proposal complete enough to plan integration and acceptance evidence.

### F-05 — Prompt-path cutover: PASS, with N-01/N-02 clarifications

**Observed:** Reconciliation §6.5 labels `ADR-INT-01` as a target rather than current behavior; inventories the application message assembler, backend template builder, lore paths, chat-service snapshots, `retrieveRAG`, frontend character messages, and narration/Story Tools; proposes contributor precedence; states exactly-once contributor, budget, snapshot, provenance, and labeling rules; and requires equivalence fixtures, adapters, capability-switched migration, fallback, and legacy retirement. RAG §12.1 separately requires an `EvidencePackage` and exactly-once labeled injection.

**Source-derived:** The current code confirms the need for this cutover: application-level message assembly, a separate backend template builder, lore resolution, chat-service canonical/memory snapshot handling with direct RAG string injection, and a character path that accepts frontend-composed messages while warning against duplicate backend prompting.

**Inferred:** The correction is substantive and does not merely assert “one assembler.” The legacy-path retirement requirement is explicit, and no runtime cutover is authorized.

**Minor follow-ups:**

- N-01: `app-prompt.go` also contains a legacy ChatML `buildPrompt` string builder with no observed Go call site in this review. Integration planning should explicitly confirm whether it is dead and retire or remove it.
- N-02: “Cutover” is already used for lorebook-store canonical ownership. Integration planning should disambiguate that existing lorebook cutover from the future prompt-assembler cutover.

Neither changes the F-05 disposition.

### F-06 — RAG migration and retrieval: PASS

**Observed:** RAG §2 records the current gaps directly, including absent index/embedding identity, read-path mutation, absent empty-history guard, and non-durable deletion handling. RAG §12.1 requires identity persistence, atomic staged replacement, corpus separation, rights/parser carry-through, durable deletion receipts, stale-citation invalidation, stale-cache guards, blocked migration until authorization, and read-only retrieval. Acceptance criteria expressly cover empty history, read-lock safety, deletion receipts, stale evidence, corpus separation, and exactly-once labeled injection.

**Source-derived:** The implementation confirms the stated seams: retrieval embeds the query, searches under a read lock, mutates document `LastUsed` files during retrieval, lacks embedding identity/version in the document/chunk model, clears only the in-memory index on delete, and the chat-service RAG helper indexes the latest message without first guarding against an empty history.

**Inferred:** The correction identifies the unsafe seams rather than concealing them. Canon protection is architecturally explicit: reference material is evidence only, canonical continuity must use canonical events, and prompt labels separate canon, memory, reference, and draft material.

**Unestablished:** Parser behavior, embedding choice, store selection, reranking, citation granularity, rights UX, retrieval quality, performance, and migration safety remain untested and unauthorized.

### F-07 — DWP order and product staging: PASS

**Observed:** Reconciliation §7 states that DWP numbering is design-input order, not build or delivery sequence. Section 7.1 separately stages Chat finalization first, Image finalization next, expansion afterward, and later UI unification as a presentation pass that does not block current-stage usability fixes. The dependency graph is labeled as a design-input graph.

**Inferred:** A reader can still mistake any numbered work-package list for a schedule, but the document now affirmatively blocks that reading in the relevant sections.

### F-08 — Android residuals: PASS

**Observed:** The Android document remains unchanged, technology-neutral, off-by-default, mutually authenticated, E2E at the application boundary, revocable, and non-experimental. Reconciliation `DWP-07` now explicitly retains data-at-rest, key/session rotation, relay metadata, pairing UX, scoped authorization, and offline behavior as residual design matters.

**Source-derived:** The Android digest is unchanged and matches the expected corrected-package identity. No transport, relay operator, cipher suite, session lifetime, pairing UX, or implementation claim has been smuggled into the residual list.

**Inferred:** Android remains correctly uncommitted on implementation while retaining the security considerations needed for future integration planning.

### F-09 — Verification clarity: PASS

**Observed:** The verification record separates current behavior from design targets, labels resource and prompt-assembler coverage as targets, preserves all frozen baseline identities, lists applied corrections by section, discloses that F-01–F-09 corrections had not yet received independent re-review, and moves its own digest to the external `.sha256` sidecar.

**Calculated:** All seven package digests, the frozen baseline digests, sizes, line counts, and sidecar digest were independently recomputed and match.

**Inferred:** The prior self-hash concern is resolved without creating a new self-referential digest problem.

## 5. Cross-document consistency

The corrected package is coherent as one connected architecture:

- Memory owns canonical character/world state and auditable history; memory items are derived evidence, not a second canon.
- RAG is confined to user-writing/reference corpora, labeled non-canon, and prohibited from writing state.
- Lorebooks remain an authored reference/trigger system distinct from both canonical state and RAG.
- Prompt assembly is assigned one future owner and precedence model, with current contributor paths named for retirement.
- Resource authority is assigned one future owner across chat, image, TTS/ASR, embedding, Android-requested, and future-studio workloads.
- Story Tools consumes shared state/memory/RAG services rather than owning a parallel store.
- Narration remains a chat-pipeline-local, attributed AI takeover without voice/image extension or prompt-path fork.
- Android remains a desktop-policy-checked remote client without its own model/resource authority.
- Product staging is Chat, then Image, then expansion; Music/3D remain future dedicated studios.
- RAG remains outside the frozen 66-row reconciliation as `DPR-07`.

No circular ownership, duplicate canonical authority, hidden modality authorization, or conflicting TTS destination was found in the corrected package.

## 6. Newly discovered issues

| ID | Severity | Issue |
|---|---|---|
| N-01 | Minor | The prompt-cutover inventory does not explicitly disposition the legacy `buildPrompt` ChatML string builder in `app-prompt.go`. No Go call site was observed in this review. Integration planning should confirm dead status and retire it. |
| N-02 | Minor | “Cutover” now has two nearby meanings: existing lorebook canonical-store cutover and future prompt-assembler cutover. Integration planning should use distinct terms. |
| N-03 | Informational | Frozen Revision 3 still marks `UD-05` unresolved in its preserved text, while the corrected addendum resolves destination and staging. This is controlled baseline preservation, not a package contradiction, but the owner should keep the distinction explicit during integration planning. |

No new blocking architectural defect was found.

## 7. Remaining blocking decisions

The following must remain closed until their separate authorized decisions, designs, evidence, and gates are complete. This list is not exhaustive, but these are the integration-planning blockers most directly connected to the re-review:

- Memory physical store and migration path: `DA-01`, `DA-M1`.
- Memory extraction model and structured-output contract: `DA-02`, `DA-M2`.
- Memory automation/conflict policy by field: `DA-06`, `DA-M3`, `DA-M6`.
- Memory deletion/forget policy, retention, tombstones, receipts, and export obligations: `DA-M5`.
- Continuity branching/merge semantics: `DA-M4`.
- Any schema migration or durable-state cutover: `OA-01`, F3-C Phase 0/D-ROOT/D-STORAGE.
- Telemetry provider order/packaging and minimum-reserve/overhead calibration: `DA-03`, `DA-H1`, `DA-H2`.
- Model compatibility-matrix source and maintenance: `DA-07`, `DA-H4`.
- Multi-GPU behavior and default low-end CPU-fallback policy: `DA-H3`, `DA-H5`.
- Consumer-specific busy-drain boundaries and durable lease-record integration behavior.
- RAG physical store, embedding identity/version, parser order, reranker policy, citation granularity, and rights UX: `DA-04`, `DA-R1`–`DA-R6`.
- RAG lifecycle evidence and runtime indexing authorization: `OA-10`, `RG-13`.
- Android transport, rendezvous/relay operator, cipher/handshake suite, session lifetime, pairing UX, device data-at-rest, and scoped authorization: `DA-05`, `DA-A1`–`DA-A7`.
- Story Tools capability map and ownership: `DA-08`.
- Narration semantics, attribution implementation, and gate.
- Music/3D studio design and authorization.
- Any runtime investigation, model load, benchmark, network experiment, pairing test, packet capture, or remote-access test: `OA-06`, `OA-09`.
- Any implementation or product/runtime gate: all CWP acceptance remains withheld.

## 8. Issues that can safely remain open for later investigation

These are genuine unknowns, but they do not prevent integration and acceptance planning:

- Hardware provider implementation order and dynamic-link packaging.
- Exact minimum reserve and runtime-overhead calibration.
- Backend-specific multi-GPU behavior.
- Compatibility-matrix sourcing and maintenance.
- Default low-end CPU-fallback versus refusal UX.
- RAG parser coverage order and reranker thresholds.
- Citation-granularity and rights-consent UX.
- Android transport selection after authorized spikes.
- Session/key rotation intervals and relay metadata minimization.
- Long-term UI-unification aesthetics.
- Performance, quality, security-validation, and benchmark evidence, all of which belong to later authorized stages.

## 9. Explicit authorization boundary

This re-review:

- Creates no new baseline.
- Changes no frozen plan, reconciliation, recovery addendum, or protected artifact.
- Preserves the prior independent review unchanged.
- Approves no architecture for implementation.
- Passes no CWP or product/runtime gate.
- Authorizes no source change.
- Authorizes no migration or durable-state cutover.
- Authorizes no model load, benchmark, runtime measurement, network listener, pairing test, packet capture, index build, or user-content inspection.
- Treats successful documentation correction as readiness for integration and acceptance planning only.

Integration and acceptance planning may identify decisions, dependencies, evidence requirements, and approval gates. It may not start implementation.

## 10. Review provenance

- Prior-review disposition: pre-correction findings F-01–F-09 are superseded by the dispositions in §4.
- Evidence classes used: Observed, Source-derived, Calculated, Inferred, Unestablished, Superseded.
- Hash evidence: Calculated and independently reproduced in §3.
- Current-behavior evidence: Observed in code and static documents; no runtime behavior was exercised.
- Architectural judgments: Inferred from the corrected documents and seam inspection.
- Privacy/runtime-sensitive unknowns: marked Unestablished where evidence would require unauthorized execution or user-content inspection.
