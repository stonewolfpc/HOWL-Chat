# Gate B to F3-C Work Queue: Defects, R5, and Current-Source Verification

**Date:** 2026-11-11
**Type:** Read-only prioritized research queue
**Status:** No owner decision, architecture selection, implementation, migration, test execution, runtime work, freeze, acceptance, or gate closure
**Inputs:** Gate B source report, Gate B reconciliation, six OPEN F3-C defects, live R5 proposal thread, D-ROOT evidence, Phase 0 assessment
**Method:** documentation comparison and source-locator mapping only

The final handoff message carries this file’s SHA-256, byte size, and line count. No digest is embedded because embedding it would change the digest.

## 1. Queue logic

Work is ordered by evidence leverage, not by engineering preference:

1. Questions where Gate B already supplies current-tree source evidence.
2. Questions where older F3-C evidence exists but needs a dirty-tree version check.
3. Questions that are genuinely unanswered and need a separate decision or separately authorized evidence.

“Needs current-source verification” means targeted read-only tracing. “Needs runtime evidence” means a separately authorized test that has not been performed.

## 2. Six-defect mapping

### P1 — ROOT-01: split chat/memory ownership

Gate B evidence already available:

- G-05/G-06: frontend memories are global and character-filtered; no chat predicate reaches the request.
- Q-09: backend chat-scoped memory APIs exist but are not used by live prompt composition.
- D-B3: `ChatMeta.ID`, frontend `sessionId`, backend `currentChatId`, and Save As behavior are not mapped to one canonical conversation key.

Existing evidence:

- `INVENTORY-F3C-ROOT-OWNERSHIP-2026-11-11.md`: no single canonical F3-C root.
- `INVESTIGATION-F3C-DROOT-CLOSURE-2026-11-11.md`: memory, candidate, continuity, native-service, AppData, and isolated roots.
- `FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md`: chat-owned generation and stable identity rules.

Still needed, read-only:

- Verify the current dirty-tree roots for production chat history, backend memory, frontend `sessionId`, `Meta.ID`, and Save As identity.
- Verify every `SetChatID` caller and the absence of a clear/reset operation.
- Map each Gate B memory path to its physical root.

Priority: highest. This defect directly controls memory scope and conversation identity for the overhaul.

### P2 — IDENTITY-03: empty `SetChatID` cannot clear identity

Gate B evidence already available:

- Backend memory scope is switched by explicit chat ID.
- No clear/reset path was found in the inspected composer/memory flow.
- Death and swap memories persist under stable or chat-labeled ownership without a discovered revocation path.

Existing evidence:

- D-ROOT investigation: `SetChatID("")` leaves state unchanged; delete and wipe do not clear the selected memory identity.
- Disposition request: separate lifecycle prerequisite, coordinated with WIPE-02.

Still needed, read-only:

- Trace delete-chat, wipe, new-chat, Save As, and session-switch paths for memory-identity behavior.
- Distinguish frontend store clearing from backend selected-identity clearing.
- Do not assume a frontend cache clear equals backend identity reset.

Priority: highest for destructive-lifecycle work; high for prompt isolation.

### P3 — WIPE-02: incomplete content wipe

Gate B evidence:

- Limited. Gate B traced save/load, clear-world/clear-scenario persistence, NDJSON versus snapshot behavior, and local/global memory stores.
- Gate B did not trace wipe coverage.

Existing evidence:

- D-ROOT investigation: delete covers legacy chat tree; wipe covers legacy chat plus gallery; memory, candidate, continuity, native context, active identity, recovery, and retained evidence are not covered.
- Disposition request: separate prerequisite under OD-02b; must not be folded into root migration.

Still needed, read-only:

- Inventory current wipe/delete entry points and every retained surface.
- Check whether Gate B’s stale-clear finding has a wipe analogue.
- Record wipe scope as unknown until traced.

Priority: critical for destructive behavior, but separate from prompt-path isolation. Do not bundle it into ROOT-01 closure.

### P4 — STARTUP-05: startup read can alter service wiring

Gate B evidence already available:

- G-01/G-02: `chat.NewService` is constructed with empty identifiers; RAG/content-library wiring follows.
- D-ROOT evidence says a corrupt/unreadable state-root snapshot can nil the service and skip RAG wiring.
- Gate B found no production generation caller, but startup construction still matters because it can change wiring and Degrade behavior.

Existing evidence:

- D-ROOT investigation startup section.
- Disposition request: separate startup/degradation prerequisite.

Still needed:

- Read-only: confirm current construction, failure, nil-service, RAG-wiring, and content-library behavior.
- Runtime: corrupt/unreadable snapshot and packaged-startup behavior remain separately authorized and unperformed.

Priority: medium for prompt correctness; high before any production-startup acceptance claim.

### P5 — REACH-04: latent working-directory candidate/continuity roots

Gate B evidence already available:

- G-01/G-03 pattern: constructed or stored paths without discovered production callers.
- G-09: chat-scoped swap memories are written but have no prompt consumer.
- D-ROOT continuity/candidate paths were already found latent through missing bindings, callers, and startup registration.

Existing evidence:

- D-ROOT investigation §§2.1–2.3.
- Disposition request: explicit exclusion or ownership ruling required; no silent deferral.

Still needed, read-only:

- Re-verify no new binding, caller, startup registration, or generated Wails method appeared in the dirty tree.
- If excluded, record the guard needed against future accidental use.
- If included, expand wipe, lock, identity, and migration scope accordingly.

Priority: medium. It is a future-ownership hazard, not current prompt-path behavior.

### P6 — ALIAS-06: no canonical path identity or alias rejection

Gate B evidence:

- None. Gate B did not inspect Windows path identity, junctions, symlinks, 8.3 names, case variants, mapped drives, UNC paths, or duplicated volumes.

Existing evidence:

- OD-18 adjudication: no supported-configuration list, cross-process primitive, or Windows proof.
- Disposition request: direct OD-18 blocker and prerequisite to D-STORAGE.

Still needed:

- Separate Windows identity and concurrency evidence.
- No prompt-path finding can substitute for it.

Priority: high for lock/storage claims; out of the prompt-path workstream.

## 3. Live R5 thread mapping

R5 is proposal-only. Its blocking findings and A14 remain open. Gate B does not close any of them, but it supplies current-source context.

### R5-A — A14 and V24: observed context head

This is the most prompt-relevant open item.

Gate B evidence:

- G-16/G-19: generic and character prompts use materially different context.
- Q-09: flattened facts lose memory IDs, scope, chat identity, lifecycle source, lore ownership, and prior-response identity.
- G-07/G-10/G-13: recent events, lore routing, and private world/scenario behavior are structurally visible.
- G-14/G-15: persisted snapshots and NDJSON behavior lose or omit provenance-relevant structure.

Existing requirement:

- H7/A14: durable history must preserve the observed context head or supply a normative immutable candidate pointer.
- V24 remains blocked.
- A14 requires a pre-amendment multi-character baseline covering prompt construction, context-memory selection, cast behavior, lore visibility, ordering, timing, retry/restart, and partial-batch visibility.

Next work:

1. Treat Gate B as the source map for that baseline, not the baseline itself.
2. Define the baseline capture separately; do not claim Gate B observed runtime behavior.
3. Keep `captured` as interim behavior only where already accepted; do not treat it as A14 closure.

Priority: highest for the architectural overhaul’s prompt assembly work.

### R5-B — C1/C5 and OD-02a/OD-02b: replacement, mutable records, and new roots

Gate B evidence:

- G-11/G-12/G-14/G-15: save, clear, snapshot, identity, and history-stream behavior.
- `updateChatDataFromAppState`, `updateAppStateFromChat`, `collectForSave`, lorebook snapshot reduction, and NDJSON behavior.

Existing issue:

- C1: mutable batch, authority, lease, attempt, candidate, and new-root records contradict a universal create-only rule.
- C5: multi-file atomic replacement is unproven; recovery must handle mixed current targets.
- OD-02 must be split into target replacement and new-root/mutable-evidence work.

Next work:

- Map the current chat files to the proposed v2 target set.
- Identify which current save/clear behavior would violate the future commit/projection rules.
- Do not rewrite current save behavior as an R6 implementation.

Priority: high.

### R5-C — C4 and seed/provenance exactness

Gate B and Phase 0 evidence:

- Per-turn model, seed mode, actual seed, engine, parameters, and attachment identity are absent from persisted chat history.
- Richer message metadata exists separately but is not the persisted chat shape.
- Global settings can retroactively change the apparent configuration of old turns.

Existing issue:

- C4: unavailable seed representation conflicts with canonical bytes and fingerprints.
- H5: sampler/engine order and thinking-budget exactness remain unresolved.
- V18/V19 remain proposed.

Next work:

- Inventory every provenance-relevant field touched by Gate B and Phase 0.
- Keep schema exactness separate from prompt-path behavior.
- No seed/model semantics may be invented.

Priority: high for durability and trust; medium for immediate prompt isolation.

### R5-D — C2/C3/C6 and H1/H2: leases, fencing, batches, slots, and state machines

Gate B evidence:

- Multi-character queue, priority ordering, prior-response accumulation, retry/abort behavior, and same facts reaching every recipient.
- No lease, fence epoch, batch, slot, abandonment, or durable operation state was found in the live composer.

Existing issue:

- Slot abandonment, batch invalidation, recovery authority, and mutable-record transitions lack complete state machines.
- V22 remains proposed.

Next work:

- Treat the live queue as behavior to be fenced and modeled, not as an implementation of R5 concepts.
- Do not retrofit lease/fence language onto current frontend orchestration.

Priority: medium now; high during durable multi-character design.

### R5-E — H3/H4/H6: migration, preferences, and canonical bytes

Gate B and F3 evidence:

- Save/load, import, legacy attribution, preference storage, stable IDs, and index-derived unstable IDs.
- F3-B import/rollback and F3-C migration constraints.

Next work:

- Map v1-to-v2 initialization only as a documentation exercise.
- Keep preference commits, migration, archive hashes, and evidence manifests separate from prompt behavior.
- H3/V23, H4/V21, and H6/V22/V23 remain proposed.

Priority: medium.

### R5-F — CR ownership and OD-18

Gate B evidence:

- None for cross-process locks, storage roots beyond prompt/memory paths, or Windows alias behavior.

Existing status:

- CR-01 through CR-05 remain open.
- OD-18 remains unresolved.
- No cross-process single-writer claim is allowed.

Next work:

- Keep OD-18 and D-STORAGE separate from prompt assembly.
- Do not let prompt-path progress imply storage-lock progress.

Priority: high for storage; separate workstream.

## 4. Prioritized queue

| Order | Work item | Why now | Evidence status |
|---|---|---|---|
| 1 | ROOT-01 ownership and conversation identity | Controls memory scope for the overhaul | Gate B plus D-ROOT evidence; needs dirty-tree root/identity map |
| 2 | A14 observed-context baseline definition | Controls prompt assembly and auditability | Gate B source map; runtime baseline separately authorized and unperformed |
| 3 | IDENTITY-03 clear/reset semantics | Blocks destructive-lifecycle acceptance | Gate B plus D-ROOT evidence; needs delete/wipe/session-switch trace |
| 4 | R5 C1/C5 and OD-02a/OD-02b mapping | Controls target replacement and mutable evidence | Gate B save/load evidence; needs v2 target mapping |
| 5 | WIPE-02 surface inventory | Critical destructive defect | Mostly unanswered in Gate B; needs separate trace |
| 6 | R5 C4/H5 provenance exactness | Durable trust and canonical bytes | Phase 0 plus Gate B; schema decisions remain open |
| 7 | STARTUP-05 degradation policy | Startup correctness and RAG wiring | Gate B plus D-ROOT; runtime proof separately authorized |
| 8 | REACH-04 inclusion/exclusion | Future ownership hazard | Existing D-ROOT evidence; needs dirty-tree recheck |
| 9 | R5 C2/C3/C6/H1/H2 fencing and state machines | Multi-character durability | Unanswered as implementation; Gate B supplies behavior map |
| 10 | ALIAS-06 and OD-18 | Storage lock prerequisite | Unanswered; separate Windows evidence required |

## 5. Explicit non-goals

- No settled owner decision is reopened.
- R1–R4 are not revived as live proposals.
- R5 is not treated as frozen authority.
- No R6 is drafted.
- No implementation, migration, wipe, lock test, model run, benchmark, network test, or destructive test is authorized.
- Audio/TTS remains separate unless a concrete chat dependency is established.
- Gate B remains a source-only report; this queue does not upgrade its findings into runtime proof.

## 6. Limitations

- This queue maps documentation and source-locator evidence; it does not re-execute the cited reviews.
- Dirty-tree version checks for every older implementation record remain future work.
- Runtime, Windows, storage, wipe, migration, and model behavior remain unproven where marked unestablished.
- No defect, gate, decision, or acceptance item is closed by this file.

## 7. Prohibited actions avoided

- No source, configuration, schema, migration, dependency, deployment, or documentation-history change was made except creating this new research file.
- No application launch, model load, benchmark, database/user-data access, network activity, packet capture, pairing, transport test, or Git mutation was performed.
- No frozen, proposal, review, defect, decision, or master-plan artifact was edited, renamed, replaced, or regenerated.
- No owner decision was recorded and no architecture was selected.
