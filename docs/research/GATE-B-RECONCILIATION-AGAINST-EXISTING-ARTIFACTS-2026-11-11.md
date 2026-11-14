# Gate B Reconciliation Against Existing Artifacts

**Date:** 2026-11-11
**Type:** Read-only documentation reconciliation
**Status:** Research only; no architecture approval, implementation, migration, test, runtime, freeze, acceptance, or owner decision
**Companion source report:** `docs/research/GATE-B-SOURCE-INVESTIGATION-2026-11-11.md`
**Method:** filename inventory, status-header search, targeted cross-reference search, and sampled reads across `docs/review`, `docs/architecture`, and `docs/research`

The final handoff message carries this file’s SHA-256, byte size, and line count. No digest is embedded because embedding it would change the digest.

## 1. What this reconciliation does

Gate B traced the present working tree. The documentation corpus already contains policy, design, frozen contracts, implementation records, reviews, and open defect requests. This file reconciles those two bodies without reopening settled owner decisions.

It uses the classification requested for the next useful work:

- Previously established policy/design.
- Previously identified implementation gap.
- Current source confirmation.
- New implementation finding.
- Contradiction requiring reconciliation.
- Unestablished.

“Not found in this investigation” is not treated as “does not exist.”

## 2. Controlling authority map

The following remain the controlling authorities encountered in this audit:

- F3-C implementation authority:
  - `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`
  - SHA-256: `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`
- F3-C verification protocol:
  - `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`
- F3-C architecture freeze:
  - `docs/review/FREEZE-F3C-ARCHITECTURE-2026-11-11.md`
- F3-C data/ownership freeze:
  - `docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md`
- Continuity policy:
  - `docs/review/AMENDMENT-02-OWNER-DECISION-REQUEST-A2-08-2026-11-11.md`
  - Amendment 02 revision 9 and its closure/verification records.
- Master-plan requirements authority:
  - `PROJECT_ISSUES_AND_EXPANSION.md` remains the ledger cited by the master-plan records.
- Audio/TTS governing contract:
  - `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md`, including frozen Stage 0.

Everything else below is proposal, review, assessment, investigation, disposition request, or historical lineage unless its own header says otherwise.

## 3. `docs/review` cluster inventory

### 3.1 F3-C ownership and compatibility foundation

Authoritative or accepted foundation:

- `docs/review/REVIEW-F3C-DECISIONS-AND-COMPATIBILITY-2026-11-10.md`
  - Accepted six F3-C decisions.
  - Recorded compatibility conflicts C1–C9, including name-backed IDs, second-writer risk, name fallback, pinned snapshots without revision, lossy editors, and media ownership.
- `docs/review/REVIEW-F3C-DEPENDENCY-AND-ACCEPTANCE-2026-11-11.md`
  - Declared F3-C unblocked but not started.
  - Inherited F3-B registry, import, cutover, availability, atomic-write, and native-evidence contracts.
  - Listed six pre-implementation decisions.
- `docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md`
  - Frozen snapshot refresh versus cast sync.
  - Frozen chat-owned generated lore and explicit promotion.
  - Frozen stable identity and unavailable-reference behavior.
- `docs/review/DECISION-F3C-ROOT-AUTHORITY-2026-11-10.md`
  - Resolved root authority; store behavior implemented and tested.
- `docs/review/DECISION-F3C-LEGACY-FOLDER-MIGRATION-2026-11-10.md`
  - Amendment recorded; import not implemented.

Implementation records:

- `docs/review/IMPLEMENTATION-F3C-STORE-FOUNDATION-2026-11-10.md`
- `docs/review/IMPLEMENTATION-F3C-IMPORT-2026-11-10.md`
- `docs/review/IMPLEMENTATION-F3C-CHAT-CONSUMER-2026-11-10.md`
  - Converted `chat.Service` to canonical stable-ID consumption.
  - Stopped writing generated world/scenario lore from chat saves.
  - Explicitly deferred promotion, refresh, cast sync, media manifests, and F3-D character lore.
- `docs/review/DESIGN-F3C-IMPORT-2026-11-10.md`
- `docs/review/F3C-LEGACY-LORE-CLASSIFICATION-2026-11-10.md`
- `docs/review/AUDIT-F3C-SINGLE-WRITER-2026-11-10.md`
- `docs/review/AUDIT-F3C-IMPORT-BOUNDARIES-2026-11-10.md`
  - Contains FAIL and OPEN rows; it is an audit result, not a second authority.
- `docs/review/VERIFICATION-F3C-NATIVE-A9-A10-2026-11-10.md`
- `docs/review/VERIFICATION-F3C-NATIVE-A12-2026-11-10.md`

### 3.2 F3-C target-2 proposal lineage

This is proposal history, not frozen authority:

- `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-2026-11-11.md`
- `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R2-2026-11-11.md`
- `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R3-2026-11-11.md`
- `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R4-2026-11-11.md`
- `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md`
  - Current proposal in the lineage.
  - Expressly supersedes R4 for review; R1–R4 remain historical.
  - Requires formal amendment AM-1 and owner decisions before conformance can be claimed.

Reviews and adjudication:

- `docs/review/REVIEW-ADVERSARIAL-F3C-TARGET-2-2026-11-11.md`
- `docs/review/REVIEW-ADVERSARIAL-F3C-TARGET-2-R2-2026-11-11.md`
- `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R3-2026-11-11.md`
- `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R4-2026-11-11.md`
- `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R5-2026-11-11.md`
- `docs/review/REVIEW-INDEPENDENT-INTEGRATION-REVIEW-F3C-TARGET-2-R5-2026-11-11.md`
- `docs/review/INTEGRATION-REVIEW-F3C-TARGET-2-R5-2026-11-11.md`
- `docs/review/DISPOSITION-MATRIX-F3C-TARGET-2-R5-MUSE-2026-11-11.md`
  - Accepts no-freeze; blocking findings remain valid; A14 remains open.
- `docs/review/ADDENDUM-PRECISION-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md`
- `docs/review/ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md`
- `docs/review/ADJUDICATION-PACKAGE-F3C-OD18-INTEGRATION-2026-11-11.md`
  - OD-18, state machine, lock proof, and change-record ownership remain unresolved or unproven.

None of these creates R6, freezes a contract, grants acceptance, or authorizes implementation.

### 3.3 D-ROOT, storage roots, and open F3-C defects

- `docs/review/INVENTORY-F3C-ROOT-OWNERSHIP-2026-11-11.md`
  - Finds no single canonical F3-C root.
  - Inventories AppData, state-root memory, working-directory candidate/continuity, native-service context, isolated persisted chats, and content-library destinations.
- `docs/review/INVESTIGATION-F3C-DROOT-CLOSURE-2026-11-11.md`
  - Candidate and continuity storage have no discovered production reachability.
  - Startup corruption can nil the chat service and skip RAG wiring.
  - Empty `SetChatID` does not clear memory identity.
  - Delete and wipe do not cover all discovered trees.
  - D-ROOT is not ready to close.
- `docs/review/DISPOSITION-REQUEST-F3C-DEFECTS-2026-11-11.md`
  - Six defects remain OPEN:
    - split chat/memory ownership;
    - incomplete wipe;
    - uncleared memory identity;
    - latent working-directory roots;
    - startup degradation;
    - missing canonical path identity/alias rejection.
- `docs/review/OWNER-DECISION-REQUEST-F3C-FOUNDATIONS-2026-11-11.md`
  - The decision request behind the root investigation.

These directly overlap Gate B’s memory-scope, chat-identity, dormant-path, lifecycle, and wipe-adjacent findings.

### 3.4 Amendment 02 continuity lineage

- `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`
  - Revision 9 policy and acceptance criteria.
- Revision archives:
  - `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-6-2026-11-11.md`
  - `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-7-2026-11-11.md`
  - `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-8-2026-11-11.md`
- Delta reviews:
  - `docs/review/REVIEW-AMENDMENT-02-REVISION-6-OD-04-CONSISTENCY-ADDENDUM-2026-11-11.md`
  - `docs/review/REVIEW-AMENDMENT-02-REVISION-7-MUSE-DELTA-2026-11-11.md`
  - `docs/review/REVIEW-AMENDMENT-02-REVISION-8-MUSE-DELTA-2026-11-11.md`
  - `docs/review/REVIEW-AMENDMENT-02-REVISION-9-MUSE-DELTA-2026-11-11.md`
- Closure and verification:
  - `docs/review/AMENDMENT-02-REVISION-CYCLE-CLOSURE-2026-11-11.md`
  - `docs/architecture/AMENDMENT-02-VERIFICATION-RECORD-2026-11-11.md`
- Owner policy:
  - `docs/review/AMENDMENT-02-OWNER-DECISION-REQUEST-A2-08-2026-11-11.md`

This lineage is internally well linked. Its missing link is implementation conformance against the current dirty tree, which is exactly what Gate B began.

### 3.5 Memory, RAG, resource, reconciliation, and master-plan lineage

- `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md`
- `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md`
- `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`
- `docs/architecture/RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md`
- `docs/architecture/RESEARCH-ANDROID-REMOTE-SECURITY-2026-11-11.md`
- `docs/architecture/VERIFICATION-DESIGN-PACKAGE-2026-11-11.md`
- `docs/review/INDEPENDENT-ARCHITECTURE-REVIEW-2026-11-11.md`
- `docs/review/INDEPENDENT-ARCHITECTURE-RE-REVIEW-2026-11-11.md`
  - Corrected package is coherent enough for integration and acceptance planning.
  - Does not approve architecture or authorize implementation.
  - Carries two minor carry-forward clarifications, including unused-`buildPrompt` disposition and lorebook-cutover terminology.
- Master-plan lineage:
  - `docs/review/HOWL-COMPLETE-PRODUCT-ARCHITECTURE-REQUIREMENTS-DEPENDENCY-MASTER-PLAN-2026-11-11.md`
  - `docs/review/ADJUDICATION-MATRIX-HOWL-MASTER-PLAN-MUSE-REVIEW-2026-11-11.md`
  - `docs/review/OWNER-DECISION-PACKET-HOWL-MASTER-PLAN-2026-11-11.md`
  - `docs/review/CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md`
  - `docs/review/REQUIREMENTS-RECONCILIATION-HOWL-2026-11-11.md`
  - `docs/review/REQUIREMENTS-RECOVERY-OWNER-ADDENDUM-HOWL-2026-11-11.md`

### 3.6 F3-B, F2-C, and shared-ownership lineage

- `docs/review/ENGINEERING-F3B-LOREBOOK-OWNERSHIP-2026-11-10.md`
- `docs/review/ARCHITECTURE-F2C-CONTENT-OWNERSHIP-CUTOVER-2026-11-09.md`
- `docs/review/ENGINEERING-F3A-OWNERSHIP-SAFETY-MIRROR-SEMANTICS-2026-11-09.md`
- `docs/review/ENGINEERING-F3A-R1-SAFETY-REPAIRS-2026-11-09.md`
- `docs/review/ENGINEERING-F2-STORAGE-MODEL-INTELLIGENCE-2026-11-09.md`
- `docs/review/ENGINEERING-F2B-STORAGE-MODEL-INTEGRATION-2026-11-09.md`

These are relevant shared foundations, not stranded material. F3-C explicitly inherits F3-B registry, import, cutover, availability, atomic-write, and native-evidence patterns.

### 3.7 Audio/TTS lineage

The TTS/audio corpus is large but mostly outside the chat-architecture overhaul:

- Governing contract:
  - `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md`
- Research, audits, reconciliation, capability policy, implementation, lifecycle, and assessment records under matching TTS/audio filenames.
- `docs/review/README.md` is principally a T2.9.5B audio-research index, not a whole-repository map.

These should not be treated as stranded chat-architecture work. They are a separate domain with their own authorization register and unresolved owner items.

## 4. Gate B reconciliation

| Gate B item | Existing record | Classification |
|---|---|---|
| Conversation-local memory default and explicit promotion | A2-08 §4; Amendment 02 `AC-C70`; F3-C freeze R2 | Previously established policy/design; Gate B is current source confirmation that no promotion mechanism was found in the live composer |
| Lifecycle belongs to the chat, not the global character | A2-08 §§2–4; Amendment 02 `AC-C71` | Previously established policy; Gate B is current source confirmation of a mismatch: DOM/roster lifecycle plus global death memory |
| Stable identity over display names | F3-C Q2/Q6; F3-C chat-consumer conversion; memory design | Previously established policy/design; Gate B confirms primary lore routing uses stable character IDs |
| Chat-owned generated lore and snapshots | F3-C freeze R1/R2; F3-C master target set | Previously established contract; Gate B confirms related save/load behavior but does not prove the newer F3-C target set is implemented in the working tree |
| One prompt assembler | Memory design; reconciled plan; independent review F-05 | Previously established design direction; Gate B is current source confirmation that generic, character, and dormant service paths remain distinct |
| Dormant secondary service | F3-C chat consumer; D-ROOT investigation | Previously identified implementation/latency concern; Gate B is current source confirmation for generation-path non-reachability, not a contradiction |
| Character memory without chat scope | A2-08 character-memory rows; memory design continuity scope | Previously identified gap/design requirement; Gate B is current source confirmation |
| Backend lore pool not hydrated on normal load | F3-B cache/hydration rules; Phase 0 source/tree caveat | New implementation finding requiring reconciliation against F3-B and the current tree |
| Clear-world/clear-scenario stale persistence | F3-C snapshot semantics; Phase 0 durability concerns | New implementation finding requiring reconciliation |
| Lorebook snapshot drops scope/owner/visibility | Phase 0 representation limits; F3-C lorebook snapshot rules | Previously identified gap plus current source confirmation |
| RAG attached to dormant service | Master-plan product-path split; memory design | Current source confirmation |
| `hasImage` without durable attachment identity | Phase 0 §3.3; owner packet OD-01 | Previously identified gap |
| Unstable message identity | Phase 0 §3.1 | Previously identified gap |
| Exact runtime behavior of dormant paths | D-ROOT residual unknowns; Gate B limits | Unestablished |

## 5. Interlink verification

### 5.1 Links that are present

- Amendment 02 revisions 6–8 link forward to revision 9 through closure and verification records.
- Delta reviews link each revision to its predecessor.
- A2-08 links continuity policy to `OD-02`, `OD-03`, `OD-04`, and downstream `R-09`/`DA-C4` work.
- F3-C decisions link to compatibility findings, dependency review, frozen contracts, implementation records, audits, and native verification.
- R5 links backward to R4 and subordinate reviews, and forward to required future R6 decisions.
- The master-plan lineage links owner packet, adjudication matrix, corrected plan, reconciliation, recovery addendum, and independent reviews.
- Gate B links backward to R-07, Muse review, erratum, register, Q-06, Amendment 02 revision 9, closure, and verification artifacts.

### 5.2 Links that are missing or weak

1. Gate B does not yet cite the F3-C contract lineage.
   - This reconciliation supplies that bridge, but Gate B itself remains a source-only report.
2. Earlier F3-C and Phase 0 source findings do not necessarily describe the current dirty tree.
   - Q-06 established 258 dirty entries and a non-provenance-proven working baseline.
   - Any claim that an older implementation record proves current behavior needs a source-version check.
3. F3-C target-2 intermediates R1–R4 are historical.
   - They must remain readable history, but the overhaul should treat R5 plus its disposition/adjudication records as the live proposal thread.
4. Target-2 reviews, addenda, disposition matrix, integration reviews, and OD-18 adjudication are numerous and easy to misread as competing authorities.
   - Only the frozen master contract and Gate v2 control implementation verification.
5. Audio/TTS docs are easy to mistake for stranded chat work because of their volume.
   - They belong to a separate authorization and product domain.
6. The old `buildPrompt` string builder and the word “cutover” have more than one meaning across lorebook, prompt-assembler, and ownership contexts.
   - The independent re-review already flags both as carry-forward clarifications.

### 5.3 Stranded-document risk assessment

- High risk if ignored:
  - F3-C defect request and D-ROOT inventory/investigation.
  - Phase 0 assessment.
  - Target-2 R5 disposition and adjudication records.
  - F3-B ownership and F3-C chat-consumer records.
  - Amendment 02 revision lineage and A2-08.
- Medium risk:
  - Older target-2 revisions and superseded reconciliations.
  - Compatibility and dependency reviews that examined older source versions.
  - Independent master-plan reviews that preserve corrections.
- Low risk for the chat overhaul:
  - Audio/TTS research and stage reviews, except where character voice or lifecycle semantics intentionally intersect chat.

## 6. What the overhaul must carry forward

1. Conversation scope and promotion remain policy; implementation mechanics remain design work.
2. Stable identity remains controlling; display names must not authorize access.
3. Chat-owned generation remains the frozen ownership rule; promotion remains explicit.
4. Snapshots remain pinned; refresh, cast sync, migration, wipe, and lock behavior remain separate decisions and evidence gates.
5. One prompt assembler remains the direction; parallel live paths remain a defect to retire.
6. Memory, lore, RAG, transcript, receipts, and persistence remain separate authorities with explicit turn snapshots.
7. Visible chat versus persisted chat remains a known gap; model, seed, attachment, continuation, and stable-ID fields cannot be assumed.
8. D-ROOT defects remain open and block canonical-root, wipe, identity-clearing, and lock claims.
9. Target-2 R5 remains proposal-only; blocking findings and A14 remain open.
10. Source-version checks are required before reusing any older implementation record against the current tree.

## 7. Limitations

- This is a documentation audit, not a re-execution of the cited reviews.
- It does not establish that every prose sentence in `docs/review` has been mapped.
- It does not change the status of any frozen contract, proposal, review, defect, or owner decision.
- It does not select architecture, authorize implementation, or close any gate.
- Source-version reconciliation against every older report remains future work.

## 8. Prohibited actions avoided

- No source, configuration, schema, migration, dependency, deployment, or documentation-history change was made except creating this new reconciliation file.
- No application launch, model load, benchmark, database/user-data access, network activity, packet capture, pairing, transport test, or Git mutation was performed.
- No frozen R-07, Muse, erratum, register, Q-06, Amendment 02, F3-C, or master-plan artifact was edited, renamed, replaced, or regenerated.
- No owner decision was recorded and no architecture was selected.
