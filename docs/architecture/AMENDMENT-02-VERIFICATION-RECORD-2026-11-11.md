# Amendment 02 verification record

**Date:** 2026-11-11  
**Scope:** Documentation and static source inspection only  
**Verdict:** Amendment revision 9 remains the closed requirements baseline; Muse's revision-3 requirement-capture PASS accepted; `M2-01` and `M2-02` are closed; `OD-02`, `OD-03`, and `OD-04` are recorded; the M3, M4, and M5 corrections are applied; Muse's revision-9 review returned PASS WITH REQUIRED CHANGES for `M6-01`–`M6-03`, which are recorded as unresolved evidence obligations `EO-01`–`EO-03`; no revision 10 is created; requirements are not validated architecture; architecture not approved; implementation and research execution not authorized; acceptance not passed.

## 1. Artifact identity

| Artifact | Value |
|---|---|
| Amendment | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` |
| SHA-256 | `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640` |
| Size | `88,499 B / 1,303 lines` |
| Sidecar | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md.sha256` |
| Sidecar verification | `OK` |
| Reviewed revision 3 | `fa766f76779e2517b31bd55da993f872525a52941b348fc8f78bbeb1d9382c44` (`61,111 B / 1,122 lines`) |
| Muse review | `docs/review/REVIEW-AMENDMENT-02-REVISION-3-MUSE-2026-11-11.md`; `0c82874b8fc5c0b82aa3a0e42b83802a8e21b2a2384b453758a4c5a429b1acf3` |
| PASS acceptance | `docs/review/AMENDMENT-02-MUSE-PASS-ACCEPTANCE-2026-11-11.md` |
| Owner decision record | `docs/review/AMENDMENT-02-OWNER-DECISION-REQUEST-A2-08-2026-11-11.md` |
| Revision 6 archive | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-6-2026-11-11.md`; `ef2d5dd88fc30cf49424f6bcf20f50aca5c599ae6d9a8f76a33bc65673370433` |
| M3 audit | `docs/review/REVIEW-AMENDMENT-02-REVISION-6-OD-04-CONSISTENCY-ADDENDUM-2026-11-11.md`; `b7be1d46abd2c82fd2aa6edda824b9ddf4c6f1f3bf90df500806d342b91f625d` |
| Revision 7 archive | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-7-2026-11-11.md`; `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1` |
| Revision 7 delta review | `docs/review/REVIEW-AMENDMENT-02-REVISION-7-MUSE-DELTA-2026-11-11.md`; `8e7f248b07d4a6e5de6e46f8c509a2e685b58e0ba4181118a631a67507baf23c` |
| Revision 7 request | `docs/review/AMENDMENT-02-REVISION-7-MUSE-DELTA-REVIEW-REQUEST-2026-11-11.md` |
| Revision 8 archive | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-8-2026-11-11.md`; `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f` |
| Revision 8 review | `docs/review/REVIEW-AMENDMENT-02-REVISION-8-MUSE-DELTA-2026-11-11.md`; `0c67713a62059db099add6b8891aa50e01ba99ca73298293b6941368b551e291` |
| Revision 8 request | `docs/review/AMENDMENT-02-REVISION-8-MUSE-DELTA-REVIEW-REQUEST-2026-11-11.md` |
| Revision 9 delta review | `docs/review/REVIEW-AMENDMENT-02-REVISION-9-MUSE-DELTA-2026-11-11.md`; `381fac1ce84e91f4e20d4649538598513f3f4de6e5bf0809bcb4d3574806a6e8` |
| Revision-cycle closure | `docs/review/AMENDMENT-02-REVISION-CYCLE-CLOSURE-2026-11-11.md` |
| Prior plan | `1e92a54434e8e3269015c47374ad57bd6ae755e64e27a07954c083ab952510cf` |
| Prior plan sidecar | `OK`, byte-identical |
| Amendment 01 | `f3a0f8b9c5d19129db85d559f1880f21451cc7aa70596d37a408952c85c3f4ea` |
| Amendment 01 sidecar | `OK`, byte-identical |

## 2. Preserved baselines verified

The following checks passed against their recorded SHA-256 digests. The plan and Amendment 01 also passed their sidecar checks with `sha256sum -c`.

- `docs/review/CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md`
- `docs/review/REQUIREMENTS-RECONCILIATION-HOWL-2026-11-11.md`
- `docs/review/REQUIREMENTS-RECOVERY-OWNER-ADDENDUM-HOWL-2026-11-11.md`
- `docs/review/OWNER-DECISION-ADDENDUM-REVISION-3-2026-11-11.md`
- `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md`
- `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`
- `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md`
- `docs/architecture/RESEARCH-ANDROID-REMOTE-SECURITY-2026-11-11.md`
- `docs/architecture/RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md`
- `docs/review/INDEPENDENT-ARCHITECTURE-REVIEW-2026-11-11.md`
- `docs/review/INDEPENDENT-ARCHITECTURE-RE-REVIEW-2026-11-11.md`

No preserved baseline document was modified. Amendment 02 itself was revised from revision 6 to revision 7 after the accepted Muse PASS; revision 3 remains the independently reviewed identity and is not claimed as current. Revision 6 is preserved byte-for-byte in the archive.

## 3. Content coverage verified

The amendment contains the required controlling sections:

- canonical conversation history separate from model input and persistent state;
- durable conversation history and media;
- protected initial conversation message;
- scenario, world, character, and lore transitions;
- contextual image generation from Chat;
- message deletion and persistent-information provenance;
- multi-character Public and Private routing;
- backend-enforced audience boundaries;
- resolved owner decision `OD-01` that Go/Wails remains the technology foundation, with no speculative rewrite and concrete evidence as the only reopening condition;
- `OD-01` cross-references in the Amendment-01 cross-reference, backend/frontend ownership, integration-map, decision-register, sequencing, governance, and migration sections;
- resolved owner decision `OD-02` on world/scenario/character/chat layering and continuity, including chat-scoped Step Out, death, return, and ghost participation;
- resolved owner decision `OD-03` on automatic versus confirmed transitions, undo, and partial-failure recovery;
- resolved owner decision `OD-04` on explicit-change scope, preserving the change conversation-locally and requiring explicit promotion for character or scenario scope;
- confirmation gate and decline behavior that prevent unapproved transitions from reaching active state, first messages, participants, memory, or model context;
- conversation-scoped memory and automatic-lore corrections with explicit promotion;
- an explicit proposal, pending, committed, and rejected/interrupted transition-state model with one durable applied marker as the commit point;
- restart and partial-failure behavior that cannot turn a pending record into a half-transition;
- atomic promotion and explicit demotion/revocation semantics, including derived indexes and previously consumed information;
- chat-scoped lifecycle isolation across canonical records, summaries, derived memories, retrieval, receipts, caches, and newly created conversations;
- canonical commit state and ownership in the event model, with uncommitted records excluded from projection and context;
- atomic authorization commits separated from asynchronous materialization, with no broader-scope artifact served;
- four distinct revocation operations with read-time authorization and historical-belief preservation;
- an observed non-receipt test for lifecycle facts in existing and newly created conversations;
- auditable change records;
- `AC-C67`–`AC-C77` as the enforcement criteria for the OD-02/OD-03/OD-04 policies;
- R-07 through R-11;
- `DA-C1` through `DA-C12`;
- `O-01` through `O-11`, with `O-11` resolved by `OD-04` and no longer open;
- stable research-question IDs `RQ-07.1` through `RQ-11.6`, including receipt privacy in `RQ-07.4` and `RQ-10.4`;
- stable acceptance IDs `AC-C01` through `AC-C77`, including `AC-C49`, the receipt criterion `AC-C66`, and the OD-04 enforcement criteria `AC-C67`–`AC-C77`;
- migration, compatibility, rollback, and recovery;
- source/evidence register and change log, including `SRC-09`–`SRC-14`, `CH-07`–`CH-13`;
- closure of the revision cycle at revision 9, with `M6-01`–`M6-03` carried as `EO-01`–`EO-03` and no revision 10;
- the Muse independent-review gate, including `RG-13` for the Go/Wails constraint;
- the Muse PASS disposition, including `M2-01`, `M2-02`, and `M2-03`;
- conflicts and deliberately unadopted assumptions;
- sequencing, governance, and authorization limits.

## 4. Static source seams

These are direct observations from the checked source at the time of this record. They are not runtime findings.

| Seam | Observation |
|---|---|
| `internal/backend/chat/types.go` | `ChatHistory` carries `ID`, sender, text, role, timestamp, `HasImage`, and `CharacterID`; it has no media identity, content blocks, audience, or source provenance. |
| `app-prompt.go` | `appendChatMessage` creates a message with text fields and appends a `ChatHistory`; the event ID is derived from the in-memory index. |
| `app-prompt.go` | `ReplaceChatHistory` rebuilds IDs from the supplied array index and persists a rewritten history. |
| `app-prompt.go` | `SendMessageStreamWithImage` decodes transient base64 and passes image bytes directly to one inference call. |
| `app-chat.go` | `buildChatHistory` carries only the `hasImage` flag, not the image reference or asset identity. |
| `frontend/dist/js/chat-input.js` | A chosen image is held in a transient `currentImage` data URL, then reduced to base64 for one request. |
| `frontend/dist/js/chat-input.js` | Audio and file attachment handlers are explicit stubs. |
| `frontend/dist/js/chat-session.js` | Reloaded history is rebuilt from text and `hasImage`; no image source is restored. |
| `frontend/dist/js/chat-actions.js` | `deleteMessage` removes the DOM node and immediately rewrites history without a confirmation prompt or provenance check. |
| `internal/backend/chat/persistence.go` | History is append-log primary with JSON fallback; rewrite replaces the log from supplied entries. |
| `internal/backend/types/message.go` | The lower-level `Metadata` type has image/audio/document references, but `ChatHistory` does not use them. |
| `internal/backend/types/memory.go` | `MemoryRecord` has no source-message field; `MemoryCandidate.SourceMessage` is free text. |
| `app-memory.go` | The frontend DTO always returns an empty `SourceMessages` slice. |
| `internal/backend/memory/candidates.go` | Candidate deduplication uses `SourceMessage` text, not a durable message ID. |
| `app-prompt.go` | `AddSystemMessage` appends a system entry to in-memory history; it does not establish a typed transition event. |
| `frontend/dist/js/chat-state-bridge.js` | `applyScenarioChange` exists, but a system line is emitted only when characters were removed. |
| `frontend/dist/js/context-selectors.js` | Scenario selection updates scenario data and can batch-update characters, but no scenario first message/transition transaction is present. |
| `frontend/dist/js/chat-stream.js` | Multi-character orchestration is frontend code; it accumulates prior responses as group context. |
| `frontend/dist/js/chat-input.js` | Character-based sending calls the frontend multi-character orchestration; no audience is sent to the backend. |
| `internal/backend/chat/service.go` | The service builds a generic/character prompt; no per-recipient Public/Private audience contract is present. |
| `internal/backend/chat/service.go` | RAG retrieval is scope/context based and injects chunks into the prompt without per-message audience enforcement. |
| `app-image-generate.go` | Image generation saves successful output to the gallery; no conversation/message association is written. |
| `frontend/dist/js/image-generator-core.js` | Image controls and job invocation exist on the Image page; no assistant-message Generate Image flyover exists. |
| `app-lorebooks-store.go` | Canonical personal lorebook ownership exists separately from automatic lore. |
| `internal/backend/chat/service.go` | Automated lore is loaded separately from manual lore; message-level source IDs are not present. |

## 5. Checks performed

- Created the Amendment 02 sidecar.
- Verified the sidecar with `sha256sum -c`.
- Verified the original plan and Amendment 01 hashes and their sidecars.
- Verified frozen plan/reconciliation/recovery hashes.
- Verified the corrected design and review package hashes.
- Verified required section and decision/workstream terms, including `OD-01` and its backend/migration cross-references.
- Verified `SRC-09`–`SRC-14`, `CH-07`–`CH-13`, and `RG-13`–`RG-16` are present and consistent with the revision-9 change log.
- Verified the revision-9 review hash and size against the review artifact.
- Verified the revision-cycle closure records `M6-01`–`M6-03` as `EO-01`–`EO-03` and creates no revision 10.
- Verified that `OD-01`, `OD-02`, `OD-03`, and `OD-04` are recorded as resolved while `O-01`–`O-10` remain open.
- Verified that `OD-02` records the layered continuity policy, `OD-03` the confirmation threshold, and `OD-04` the explicit-change scope.
- Verified that `O-11` is recorded as resolved by `OD-04` rather than silently removed.
- Verified that `AC-C67`–`AC-C77` cover override preservation, confirmation ordering, the commit boundary, canonical commit-state ownership, atomic authorization, revocation operations, memory scope, derived-surface lifecycle isolation, observable non-receipt, and auditable change history.
- Verified that §4 carries `commitState` and `appliedMarkerRef` (or named equivalent), §4.1 excludes uncommitted records from projection/context, §4.2 assigns commit and pending records, and the model requires one store or explicit two-phase discipline.
- Verified that §7.1 contains the confirmation gate, decline path, proposal/pending/committed/interrupted states, the applied marker as sole commit point, and restart behavior; `AC-C23` no longer implies pre-confirmation application.
- Verified that §9.3.1 separates atomic authorization from asynchronous materialization, forbids broader-scope serving, names the four revocation operations, and requires read-time enforcement rather than cache invalidation alone.
- Verified that §9.3 makes Chat-derived entries conversation-scoped by default and that R-09 is explicitly barred from weakening `OD-02`, `OD-03`, or `OD-04`.
- Verified that `AC-C71` exercises existing Conversation B and new Conversation C with forced retrieval, summary generation, assembled model input, and receipt inspection.
- Verified that the revision-6, revision-7, and revision-8 archive hashes, the M3/M4/M5 review hashes, and the M6 revision-9 review hash match their artifacts.
- Verified the final revision-9 amendment hash after the delta review remains `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640`; no revision 10 was created.
- Verified Muse's review digest and size against the review artifact.
- Verified that `M2-01` is closed in §5.3, §10.4, `AC-C49`, `AC-C66`, `RQ-07.4`, and `RQ-10.4`.
- Verified that `M2-02` is recorded as closed by `OD-02`, `OD-03`, and `OD-04`, and that field taxonomy is explicitly left to `R-09`/`DA-C4`.
- Verified that `A2-08` records the owner-policy resolution, and that `DA-C4` still blocks implementation of the transition model.
- Verified that `M2-03` is recorded as informational and §12.1 remains the dependency authority.
- Inspected the source seams listed above.

## 6. Activities not performed

- No application source code was modified.
- No model files, runtime configuration, settings, or databases were changed.
- No migration was run.
- No model was loaded.
- No benchmark or runtime probe was run.
- No network or Android pairing test was run.
- No Public/Private routing was implemented.
- No contextual image workflow was implemented.
- No production conversation behavior was changed.
- No architecture approval or acceptance gate was claimed.

## 7. Review handoff

Muse reviewed revision 3 through:

`docs/review/AMENDMENT-02-MUSE-REVIEW-REQUEST-2026-11-11.md`

The review produced **PASS** for requirement capture. That PASS is accepted in `docs/review/AMENDMENT-02-MUSE-PASS-ACCEPTANCE-2026-11-11.md` and remains limited to revision 3. It does not authorize R-07–R-11 research, implementation, migration, runtime investigation, model loading, benchmarking, network testing, Android pairing, architecture approval, or acceptance.

Revision 4 was the controlled post-review delta that closed `M2-01` and recorded `M2-02` as `O-11`. Revision 5 recorded the owner's broad policy as `OD-02` and `OD-03`. Revision 6 recorded `OD-04`, resolving `O-11`; the decision record is `docs/review/AMENDMENT-02-OWNER-DECISION-REQUEST-A2-08-2026-11-11.md`. Revision 7 applied the M3 enforcement audit. Muse reviewed that delta in `docs/review/REVIEW-AMENDMENT-02-REVISION-7-MUSE-DELTA-2026-11-11.md` and returned PASS WITH REQUIRED CHANGES for `M4-01`–`M4-03`. Revision 8 applied those corrections. Muse reviewed revision 8 in `docs/review/REVIEW-AMENDMENT-02-REVISION-8-MUSE-DELTA-2026-11-11.md` and returned PASS WITH REQUIRED CHANGES for `M5-01`–`M5-04`. Revision 9 applied those corrections. Muse reviewed revision 9 in `docs/review/REVIEW-AMENDMENT-02-REVISION-9-MUSE-DELTA-2026-11-11.md`; `M6-01`–`M6-03` are recorded as evidence obligations in `docs/review/AMENDMENT-02-REVISION-CYCLE-CLOSURE-2026-11-11.md`. The field taxonomy remains design work under `R-09`/`DA-C4`.

## 8. Outstanding research and decisions

R-07 through R-11 remain research plans, not completed dossiers. Owner decisions `O-01` through `O-10` remain open. `OD-01`, `OD-02`, `OD-03`, and `OD-04` are resolved. `OD-01` constrains backend and migration work to Go/Wails unless the owner reopens it with concrete engineering evidence. `DA-C1` through `DA-C12` remain open and block implementation for the affected areas.

`M2-01` is closed as a documentation requirement. `M2-02` is closed at the owner-policy level by `OD-02`, `OD-03`, and `OD-04`; the exact field taxonomy and promotion mechanics remain `R-09`/`DA-C4` design work.

Muse’s M3 delta review confirmed `M3-01`–`M3-05` and treated `M3-06` as informational, while requiring `M4-01`–`M4-03`. Revision 8 applied the commit boundary and restart rule, atomic promotion and demotion/revocation semantics, and the expanded derived-surface/new-conversation isolation. Muse’s M4 review then required `M5-01`–`M5-04`. Revision 9 applied the canonical commit-state model, authorization/materialization separation, four revocation operations, and observed non-receipt test. Muse’s revision-9 review confirmed those changes and identified `M6-01`–`M6-03`; the owner closed the revision cycle at revision 9 and recorded them as `EO-01`–`EO-03`.

The owner decision is recorded. Architecture approval, research, and implementation remain unauthorized.
