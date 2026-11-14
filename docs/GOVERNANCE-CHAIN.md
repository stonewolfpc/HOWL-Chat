# HOWL Chat — Reconstructed Governance and Provenance Chain

**Status:** Historical reconstruction and navigation only. This document creates no decision, contract, freeze, acceptance, or implementation authorization.
**Date:** 2026-11-12
**Basis:** The artifacts listed in `docs/DOCUMENT-REGISTRY.md`, read in their current working-tree form.
**Rule:** This chain connects existing records. It does not reopen a settled decision, reinterpret a frozen contract, or replace an owner decision.

---

## 1. Why this chain exists

The project accumulated five days of governance work across lorebooks, worlds/scenarios, chat persistence, memory, TTS, audio, and requirements reconciliation. The records are individually traceable but were not collected into one chain. The result was predictable: the same ground could be re-investigated because the controlling record was not obvious from the folder.

This document reconstructs the **dependency path**, not a new policy. It answers four questions for each stage:

1. What problem forced this record?
2. What did the record establish?
3. What did it hand to the next stage?
4. What is its status now?

`docs/DOCUMENT-REGISTRY.md` is the completeness check. If a document is not in this narrative, it still has a registry row and a family home.

---

## 2. Chain at a glance

| Stage | Question | Controlling outcome | Current status |
|---|---|---|---|
| 0. Content ownership | Who owns lorebook content, and can the browser remain authoritative? | One canonical lorebook store; browser becomes a rebuildable cache; explicit import | Implemented evidence; historical handoff to F3-C |
| 1. F3-B lorebook cutover | How do stable identity, single writer, and explicit import actually work? | F3-B-Q1–Q5 answered; per-entry files; create-only journaled import; rollback | Implementation complete; native-window acceptance pending at record time |
| 2. F3-C decisions | Can the same model apply to worlds/scenarios? | Six decisions accepted; nine compatibility findings; F3-C unblocked but not started | Decision record |
| 3. F3-C data/ownership freeze | What exactly is owned, imported, snapshotted, and migrated? | R1–R3, document/identity/path/reference/update/import/snapshot contracts | **Frozen** |
| 4. Root authority and legacy migration | What happens when multiple roots exist, and what happens to old name folders? | Roots are equivalent search locations; indeterminate copies are not editable; LF-1 one-time migration | Resolved decision records |
| 5. Single-writer/import audits | Does the implementation actually enforce one writer and safe import? | A10 source/native pass; export bypass closed; import findings F-01–F-09 | Audit evidence; later native A9/A10/A12 pass |
| 6. F3-C implementation evidence | What was built and verified? | Store foundation, import, chat consumer, native verification | Evidence; not acceptance |
| 7. Architecture freeze and Gate v2 | What is the stable implementation authority and how is it verified? | Consolidated master contract; architecture freeze; adopted Gate v2 | **Frozen architecture**; Gate v2 adopted |
| 8. Root-ownership investigation | Does the implementation actually carry `(rootID, chatID)`? | No single canonical root; empty `SetChatID` cannot clear; delete/wipe gaps | Open defects; D-ROOT unresolved |
| 9. Gate B / ROOT-01 | What does the current dirty tree actually do? | No root ID; timestamp collisions; unvalidated path; memory scope gaps; frontend fallback IDs | Current source evidence; runtime unestablished |
| P1. Amendment 02 | What conversation-continuity requirements govern the future architecture? | Revision 9 closed requirements baseline | Closed, not architecture approval |
| P2. Master plan/requirements | What is the whole product, and what is missing? | Recovery baseline, reconciliation, complete plan, adjudication, corrected plan rev 3 | Proposed, not adopted |
| P3. TTS/audio | What is authorized for TTS and audio research? | Stage 0 frozen protocol; A2 adopted; T2.9.5B r14 frozen decision-control | Design authority only; stage not authorized |

---

## 3. Origin: content ownership and lorebooks

The trail does not start with F3-C. It starts with the recognition that lorebook content had no single owner and that browser storage could remain silently authoritative.

### 3.1 F2-C content-ownership cutover

`docs/review/ARCHITECTURE-F2C-CONTENT-OWNERSHIP-CUTOVER-2026-11-09.md` is the content-ownership architecture that precedes the F3 sequence. Its role is to establish the cutover framing: canonical content, browser cache, explicit ownership, no silent dual writer.

This is the conceptual ancestor of the F3-B/F3-C pattern. It is not a substitute for the later frozen contracts.

### 3.2 F3-A ownership and safety

`docs/review/ENGINEERING-F3A-OWNERSHIP-SAFETY-MIRROR-SEMANTICS-2026-11-09.md` and `docs/review/ENGINEERING-F3A-R1-SAFETY-REPAIRS-2026-11-09.md` are the safety/mirror stage.

F3-B's own starting-point section records the handoff: F3-A-R1 was accepted with documented limitations. It removed unlabelled automatic mirror writes, made replacement and wipe safer, and left lorebooks **browser-owned with no file identity**. That is the gap F3-B was created to close.

### 3.3 F3-B lorebook ownership cutover

`docs/review/ENGINEERING-F3B-LOREBOOK-OWNERSHIP-2026-11-10.md`

F3-B resolved five questions:

| ID | Decision |
|---|---|
| F3-B-Q1 | One JSON document per entry: `<root>/<id>.json`, envelope `howl-lorebook-v1` |
| F3-B-Q2 | One canonical store; frontend repository derives the browser cache; chat refreshes canonical before prompt assembly |
| F3-B-Q3 | Preserve existing identities on import; random `lb_` identities for new entries; duplicates reassigned and reported |
| F3-B-Q4 | Explicit, create-only, idempotent by identity + content hash, journaled, reversible import |
| F3-B-Q5 | No propagation subsystem; canonical delete removes the file; rename retains identity and path |

Why this matters to the whole project: F3-B established the reusable ownership grammar that F3-C later applied to worlds and scenarios:

- one canonical writer;
- identity names the path, not the title;
- explicit import, never a silent migration;
- cutover marker;
- rollback;
- unavailable is not absence;
- browser storage is a rebuildable cache, not a second owner.

F3-B's status is implementation complete, automated tests passed, UI harness integration verified, **native-window acceptance pending**. It is implementation evidence, not acceptance of F3-C.

---

## 4. F3-C: from decision to freeze

### 4.1 Decision and compatibility review

`docs/review/REVIEW-F3C-DECISIONS-AND-COMPATIBILITY-2026-11-10.md`

This review accepted six recommendations and recorded nine compatibility findings. The six decisions became the F3-C spine:

| ID | Decision |
|---|---|
| F3-C-Q1 | Folder per item: `<root>/<id>/world.json` or `scenario.json` |
| F3-C-Q2 | Stable, opaque identity; preserve usable legacy IDs; no title authorization |
| F3-C-Q3 | Explicit, journaled, create-only, reversible import |
| F3-C-Q4 | Pinned chat snapshots; explicit refresh; visible staleness |
| F3-C-Q5 | Preserve unavailable; never substitute |
| F3-C-Q6 | Stable IDs for world/scenario lore; rename never changes ownership |

The compatibility findings are the reason the later governance is so strict. In particular:

- **C1:** world/scenario IDs have mixed formats and can collide.
- **C2:** legacy chat metadata can put a display name in an `ID` field.
- **C3:** the chat service is a **second writer keyed by name**.
- **C6:** the editor drops unknown fields on save.
- **C7:** folder-per-item assets conflict with F3-F media ownership.

These are not open questions now. They are recorded constraints that the freeze and audits had to resolve.

### 4.2 Dependency and acceptance review

`docs/review/REVIEW-F3C-DEPENDENCY-AND-ACCEPTANCE-2026-11-10.md`

At that point F3-C was **unblocked and not started**. The review found that its dependency chain was sound because F3-B/W27 had already exercised the ownership contracts, registry roles, isolation, and native evidence model. It also found the actual gap: no canonical world/scenario store existed; worlds and scenarios were still browser-owned with write-only native copies, name-keyed chat reads, and snapshots that carried no revision.

This is the handoff from "the model works for lorebooks" to "apply it to worlds/scenarios."

### 4.3 Frozen data and ownership contracts

`docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md`
SHA-256: `aa11be0047c9410a72a37334d758dbd4d05dfb60cd318477cd5b103157df40ca`

The freeze resolves three cross-cutting decisions:

- **R1:** document refresh and scenario cast sync are separate, explicit operations.
- **R2:** canonical lore belongs to the item; generated lore belongs to the chat until explicitly promoted. Legacy `automated_*` files are imported once as `legacy_generated_unattributed`; no chat attribution is invented.
- **R3:** item folders contain the document and references only; no binaries.

It also freezes the document envelope, identity, paths, references, update, import/rollback, and snapshot contracts. Change control is explicit: any change to a frozen item requires a new dated decision record.

### 4.4 Root authority and LF-1

`docs/review/DECISION-F3C-ROOT-AUTHORITY-2026-11-10.md`

The root-authority decision rejects the idea that registration order is authority. Roots are equivalent search locations. A found copy is authoritative only when every registered root was searched and exactly one copy exists. Otherwise the result is **found, authority indeterminate**: readable for display, not editable. Indeterminate is never treated as absence.

`docs/review/DECISION-F3C-LEGACY-FOLDER-MIGRATION-2026-11-10.md` records **LF-1**: one narrowly scoped, explicit, journaled, reversible migration of legacy name-folder content. Old name folders remain never-read at runtime, never a fallback, never a permanent import source, never written to, and never deleted except by the explicit migration's own rules. Unresolved content is quarantined; a display-name match alone is never sufficient.

`docs/review/F3C-LEGACY-LORE-CLASSIFICATION-2026-11-10.md` separates three populations that had been conflated:

1. legacy automated files already inside a canonical item folder — read-only, classified in memory, prompt-eligible;
2. legacy automated files in old name-keyed folders — inert at runtime; only the explicit migration may read them;
3. newly generated conversation discovery — chat-owned; no canonical origin assigned.

This classification is the direct answer to "where did the lorebook trail go?" It went into the immutable item/chat ownership distinction.

### 4.5 Single-writer audit

`docs/review/AUDIT-F3C-SINGLE-WRITER-2026-11-10.md`

The audit traced every production read/write path for world/scenario content, documents, and item folders. It removed or blocked the remaining legacy name-keyed surfaces and closed a real second-writer bypass: `ExportContentItem` previously accepted worlds/scenarios and wrote a `world.json`/`scenario.json` envelope into a registered root. The binding now refuses both categories. World/scenario export is a deliberate capability gap until a store-backed export exists.

The audit's A10 result is source-level and native-runtime acceptance pass. It is not F3-C acceptance.

### 4.6 Import-boundary audit and later native evidence

`docs/review/AUDIT-F3C-IMPORT-BOUNDARIES-2026-11-10.md`

The import audit found unresolved contract/boundary failures: reference normalization, registered-root rejection, media adoption, owner stamping, staging cleanup, Wails binding exposure, and conflict labeling. At audit time A9/A10/A12 were open.

Later records resolve the evidence sequence:

- `docs/review/VERIFICATION-F3C-NATIVE-A9-A10-2026-11-10.md` — native A9/A10 pass;
- `docs/review/VERIFICATION-F3C-NATIVE-A12-2026-11-10.md` — native A12 pass;
- `docs/review/CORRECTION-BASELINE-F3C-IMPORT-2026-11-10.md` — corrected baseline;
- `docs/review/FINAL-VERIFICATION-F3C-IMPORT-2026-11-10.md` — final import verification.

The audit findings remain the audit-time record even where later evidence changed the gate status. That is why the registry distinguishes **finding** from **current status**.

### 4.7 Implementation records

- `docs/review/IMPLEMENTATION-F3C-STORE-FOUNDATION-2026-11-10.md`
- `docs/review/IMPLEMENTATION-F3C-IMPORT-2026-11-10.md`
- `docs/review/IMPLEMENTATION-F3C-CHAT-CONSUMER-2026-11-10.md`

These are implementation evidence. They do not accept F3-C and do not authorize P1a. The chat-consumer report is explicitly limited by the later `F3C-LEGACY-LORE-CLASSIFICATION` clarification: "legacy generated lore remains prompt-eligible" is true only for the canonical-folder population, not for old name folders or newly generated chat discoveries.

---

## 5. Architecture freeze and Gate v2

### 5.1 Consolidated master contract

`docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`
SHA-256: `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`

This is the candidate implementation authority. It supersedes the isolated Master-1 through Master-3c amendment sequence as the consolidated implementation authority. Its core invariants are:

- chat identity is `(rootID, chatID)`, never a display name;
- chat IDs are not reused within a root incarnation;
- wipe creates a new root incarnation;
- every durable revision has one commit node and one commit record;
- every commit node has one valid receipt, or a specifically permitted reconstructable missing receipt;
- commit and receipt inventories must both be complete and mutually consistent before verified status;
- publication occurs only after durable commit;
- historical and current state are never conflated;
- every authoritative mutation caller must be inventoried.

The contract explicitly says it does **not** accept F3-C, authorize P1a, or claim implementation evidence.

### 5.2 Architecture freeze

`docs/review/FREEZE-F3C-ARCHITECTURE-2026-11-11.md`
SHA-256: `2d9220e8a9cfeb15e621836892d89b5bae23bec4b7ce7c7555d5af0d62fb315a`

The freeze record freezes the architecture only. Frozen scope includes chat identity/root incarnation/revisions, commit records/receipts/successor links/fork inventories, canonicalization and target hashes, target set and ownership, locks/publication/operation map, creation/Save As/append/autosave/switch/delete/cast, journal recovery, tombstones, memory/history migration, wipe staging/tree hashing/root identity commit, and Gate v2.

It explicitly does **not** freeze:

- F3-C implementation completeness;
- P1a authorization;
- F3-D;
- cast-confirmation implementation;
- A9/A10/A12 evidence for the new persistence contract.

### 5.3 Gate v2

`docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`
SHA-256: `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`

Gate v2 is the adopted implementation-verification protocol. It tests both the implementation and the frozen contract. It does not grant P1a, F3-C acceptance, deployment, or F3-D.

---

## 6. Root ownership: from contract to contradiction

### 6.1 The owner decision request

`docs/review/OWNER-DECISION-REQUEST-F3C-FOUNDATIONS-2026-11-11.md` is the request that opened D-ROOT and D-STORAGE. It is a request, not a decision.

### 6.2 Root inventory and D-ROOT closure

`docs/review/INVENTORY-F3C-ROOT-OWNERSHIP-2026-11-11.md`
SHA-256: `343104b528818df699d752c419770fdc1bfd7fe71bc1f249e0dc6847f308c6b8`

The inventory establishes that there is **no single canonical F3-C root**. It identifies at least four relevant surfaces, plus isolation and content-library surfaces:

1. AppData legacy chat store;
2. `StateRoot/Chats` memory store;
3. `cwd/Chats` candidate/continuity store;
4. `StateRoot/Chats/<chatName>` native-service context;
5. isolated `<state>/persisted-chats`;
6. content-library `Documents\Howl\Chats`.

It also records the alias problem: no resolver calls `EvalSymlinks`, a Windows final-path API, volume/file-ID lookup, or a reparse-point check.

`docs/review/INVESTIGATION-F3C-DROOT-CLOSURE-2026-11-11.md`
SHA-256: `494539b9dde27a78c9c4c091131a93dce99c1debd81d11fea911e6d411e2e6a6`

The D-ROOT closure investigation adds:

- candidate/continuity have no production binding/caller found; they remain latent, not proven dead;
- `chat.Service` startup reads `StateRoot/Chats/memories_snapshot.json` with an empty name and can fail construction;
- `SetChatID("")` is a no-op;
- delete covers the legacy chat tree; wipe covers the legacy chat tree and gallery; memory/candidate/continuity/native context are not covered;
- the six defects remain open.

### 6.3 Disposition request

`docs/review/DISPOSITION-REQUEST-F3C-DEFECTS-2026-11-11.md`
SHA-256: `ce32ac0bedbe216b55c065be073b2165209e0ebad4cb435897b829336d29d2c5`

This is a formal owner disposition request, not a decision. It lists the six defects and requires six independent closures:

- `DEF-F3C-ROOT-01` — split chat/memory ownership;
- `DEF-F3C-WIPE-02` — incomplete wipe;
- `DEF-F3C-IDENTITY-03` — empty `SetChatID` cannot clear identity;
- `DEF-F3C-REACH-04` — latent `cwd/Chats` candidate/continuity roots;
- `DEF-F3C-STARTUP-05` — startup read can fail and alter wiring;
- `DEF-F3C-ALIAS-06` — no canonical path identity or alias rejection.

### 6.4 Gate B

`docs/research/GATE-B-SOURCE-INVESTIGATION-2026-11-11.md`
`docs/research/GATE-B-RECONCILIATION-AGAINST-EXISTING-ARTIFACTS-2026-11-11.md`
`docs/research/GATE-B-TO-F3C-DEFECT-WORK-QUEUE-2026-11-11.md`

Gate B moved from the prior reports to the current tree. It established the work queue and handed ROOT-01 priority 1. It did not select a root or close a defect.

### 6.5 ROOT-01

`docs/research/ROOT-01-PRODUCTION-ROOT-OWNERSHIP-AND-CONVERSATION-IDENTITY-2026-11-11.md`
SHA-256: `1a94c7f53a4e5809f34e76128175fabe3dc74131a3c212ff1d430be65dca65ac`

ROOT-01 is the current source investigation. It establishes that the implementation does **not** carry `(rootID, chatID)`:

- no root ID exists in `ChatMeta`, `ChatData`, `AppConfig`, `ChatSummary`, or `MemoryRecord`;
- chat IDs are second-resolution timestamps with no collision check;
- `NewChat` does not reject an existing directory;
- `SaveChatAs` can rewrite the source instead of copying;
- chat IDs are unvalidated path components;
- `SaveChatData` accepts a frontend-supplied `meta.id`;
- `GetChatList`, `LoadChat`, and `GetActiveChatMeta` use different ID fallbacks;
- memory records have no root/chat field; `SetChatID("")` is a no-op; initial scope is `"default"`;
- delete/wipe leave memory and latent roots;
- `SaveChatAs` copies the wrong `memories.json`;
- `chat.Service` is name-keyed and starts empty;
- backend memory IDs collide at one-second resolution;
- the frontend memory store is global;
- the frontend save fallback can invent `Date.now()` IDs;
- candidate/continuity resolve `os.Getwd()`;
- character identity is mixed name/ID and has no continuity ID;
- RAG has chat scope fields but no root identity.

This is source evidence. Runtime behavior remains unestablished. ROOT-01 does not authorize implementation or architecture selection.

**The contradiction is now explicit:** the frozen master contract requires `(rootID, chatID)`; the current implementation does not have a root ID at all. That is the bridge from the governance chain to the actual work queue.

---

## 7. Parallel chain: Amendment 02

Amendment 02 is a separate owner-requirements chain, not part of F3-C. It must not be used to reopen F3-C or vice versa.

| Node | Artifact | Status |
|---|---|---|
| Requirements amendment | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` | Revision 9; closed requirements baseline; planning only |
| Independent review | `docs/review/REVIEW-AMENDMENT-02-REVISION-3-MUSE-2026-11-11.md` | PASS for requirement capture; not architecture approval |
| Owner decisions | `docs/review/AMENDMENT-02-OWNER-DECISION-REQUEST-A2-08-2026-11-11.md` | Resolved by OD-02/03/04 |
| Revision 6/7/8 reviews | `docs/review/REVIEW-AMENDMENT-02-REVISION-6-OD-04-CONSISTENCY-ADDENDUM-2026-11-11.md`, `…-REVISION-7-MUSE-DELTA…`, `…-REVISION-8-MUSE-DELTA…` | Historical review evidence |
| Revision 9 review | `docs/review/REVIEW-AMENDMENT-02-REVISION-9-MUSE-DELTA-2026-11-11.md` | PASS WITH REQUIRED CHANGES for M6-01–M6-03 |
| Closure | `docs/review/AMENDMENT-02-REVISION-CYCLE-CLOSURE-2026-11-11.md` | Closed at rev 9; no rev 10; M6 obligations remain |
| Verification record | `docs/architecture/AMENDMENT-02-VERIFICATION-RECORD-2026-11-11.md` | Documentation/static verification only |

The amendment's own header says it does not replace Amendment 01 and is not implementation authorization. That status is preserved.

---

## 8. Parallel chain: requirements and master plan

This chain is about the whole product, not F3-C. It must not be used to smuggle F3-C acceptance or implementation authorization.

| Node | Artifact | Role / status |
|---|---|---|
| Requirements ledger | `PROJECT_ISSUES_AND_EXPANSION.md` | Authoritative requirements/expansion ledger per the planning documents |
| Recovery baseline | `docs/review/REQUIREMENTS-RECOVERY-OWNER-ADDENDUM-HOWL-2026-11-11.md` | Owner-characterized baseline; SHA-256 `c9f8dc1c…` |
| Reconciliation | `docs/review/REQUIREMENTS-RECONCILIATION-HOWL-2026-11-11.md` | Reconciliation deliverable; not the corrected plan |
| Complete master plan | `docs/review/HOWL-COMPLETE-PRODUCT-ARCHITECTURE-REQUIREMENTS-DEPENDENCY-MASTER-PLAN-2026-11-11.md` | Planning inventory; not a contract/freeze/acceptance |
| Independent review | `docs/review/REVIEW-INDEPENDENT-HOWL-COMPLETE-PRODUCT-ARCHITECTURE-MASTER-PLAN-2026-11-11.md` | Review evidence |
| Adjudication | `docs/review/ADJUDICATION-MATRIX-HOWL-MASTER-PLAN-MUSE-REVIEW-2026-11-11.md` | Read-only adjudication |
| Correction disposition | `docs/review/MUSE-CORRECTION-DISPOSITION-HOWL-2026-11-11.md` | MUSE-01–MUSE-12 corrections |
| Corrected plan rev 3 | `docs/review/CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md` | Proposed, not adopted; SHA-256 `d642039c…` |
| Verification | `docs/review/VERIFICATION-HOWL-CORRECTED-MASTER-PLAN-2026-11-11.md` | Consistency verification, not adoption |

The corrected plan's own status is **PROPOSED — REVISION 3; NOT ADOPTED; NO IMPLEMENTATION AUTHORIZATION**. The verification record confirms section/row completeness; it does not change that status.

---

## 9. Parallel chain: TTS and audio

This chain is separate from F3-C. It has its own authorization register and its own frozen protocol.

| Node | Artifact | Role / status |
|---|---|---|
| T2.9.5B authorization | `docs/t2.9.5b-authorization-decision-register-r14.md` | Exact r14 register; frozen by its freeze record |
| T2.9.5B freeze | `docs/t2.9.5b-authorization-decision-register-r14.freeze.md` | FROZEN decision-control revision; stage still NOT AUTHORIZED |
| T2.9.5B review index | `docs/review/README.md` | T2.9.5B evidence sections; directory-wide header points to this chain |
| TTS governing design | `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` | Governing architecture; Stage 0 frozen in §17 |
| TTS §17.19 draft | `docs/review/AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md` | Historical draft; content adopted, file preserved |
| TTS adoption | `docs/review/ADOPTION-howl-tts-capability-policy-c1-17.19-2026-11-07.md` | Adoption record; design authority only |
| C1/P17 assessment/review set | `docs/review/ASSESSMENT-howl-tts-stage3-c1-*`, `REVIEW-howl-tts-stage3-c1-*`, `DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md` | Evidence and decision records; no stage authorization |
| Audio architecture | `docs/audio-architecture/README.md` and `01–03` | Implementation architecture; not T2.9.5B authorization |

The T2.9.5B register's own header says **DRAFT — REVISION 14 — NOT FROZEN — NOT AUTHORIZED**, but the freeze record explicitly freezes the exact r14 artifact and sidecar. The freeze record controls the register-freeze state. The stage remains not authorized.

The TTS §17.19 draft's own header says **NOT ADOPTED**, but the adoption record adopted Revision 6 content and incorporated it into the governing design. The draft remains the historical standalone artifact; it is not a second governing authority.

---

## 10. Hash anchors for the main chain

These are identity anchors, not a complete hash list. The registry contains every file's hash.

| Artifact | SHA-256 |
|---|---|
| F3-B lorebook ownership | `a87c2682c4f869861e5f5c0e07d6324ef8ebf447fe9afe09b3b6d07bd7786443` |
| F3-C decisions/compatibility | `86050d38851c603b99d73c4e0c7c91ca7e78b981846af5f550626566bad30e85` |
| F3-C dependency/acceptance | `f48e098926df190ab6d7a86d48dbb00efe8d0ae82d3f3566ac23db9992f93fb4` |
| F3-C frozen data/ownership | `aa11be0047c9410a72a37334d758dbd4d05dfb60cd318477cd5b103157df40ca` |
| F3-C root authority decision | `8ea62621e5f3bfd545f1302feb8e800d9dc64535c459cdb470da4bfa25907a19` |
| F3-C LF-1 decision | `472f4a338106d6186923c49a27941b57dad21cc7859aa6926be7b4826e4354db` |
| F3-C legacy lore classification | `2a230658ea70484c6e550f1bd3bc57b24e7693afd3e2b8ac7d1a0c200422aa0f` |
| F3-C single-writer audit | `5cb8fccccda0f1f8e2b3d8840c207999c53eb4f580df437ddbe491a13aacbad4` |
| F3-C import-boundary audit | `0159ea96ab3322787f9a20a7152e4f9de4d0f1e9ee7e59576a2e516fade585ad` |
| F3-C store foundation | `170b768c91cb5dc0cc340d10b1baca91ac70120d6e58f3ef61e787be9ee28062` |
| F3-C import implementation | `e93886d9d8d510110dc5c91dc954f31de079917d182b91026ca62e4a3c7e8362` |
| F3-C chat consumer | `7ea9381817d5f8ae358395b4d9e4130adc2ed4ccfcc603669154faef3de31bbc` |
| F3-C native A9/A10 | `ebdc6193b28806b005e5920c2044c24544cc820e6b431c2d48383cae711dd850` |
| F3-C native A12 | `3613a3f0d196c63eec8c8ed84118a4584d854dab01562fc55e31130e05bb003f` |
| F3-C master architecture | `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31` |
| F3-C architecture freeze | `2d9220e8a9cfeb15e621836892d89b5bae23bec4b7ce7c7555d5af0d62fb315a` |
| F3-C Gate v2 | `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a` |
| F3-C root inventory | `343104b528818df699d752c419770fdc1bfd7fe71bc1f249e0dc6847f308c6b8` |
| F3-C D-ROOT closure | `494539b9dde27a78c9c4c091131a93dce99c1debd81d11fea911e6d411e2e6a6` |
| F3-C defect disposition request | `ce32ac0bedbe216b55c065be073b2165209e0ebad4cb435897b829336d29d2c5` |
| Gate B source investigation | `3ad665497a9f16fd75bcefaee3e2734f7495addc43c84c6767c5dd3ea7aca17f` |
| Gate B reconciliation | `c75483fab28b743349e1025bda03351795eb82cd714bf41d96ca2972914568db` |
| Gate B queue | `d55aa7282059a48ab63679770357aade7db6fd1308578c34d45cf15e2ad5145c` |
| ROOT-01 | `1a94c7f53a4e5809f34e76128175fabe3dc74131a3c212ff1d430be65dca65ac` |
| Amendment 01 | `f3a0f8b9c5d19129db85d559f1880f21451cc7aa70596d37a408952c85c3f4ea` |
| Amendment 02 rev 9 | `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640` |
| Amendment 02 closure | `f828fd0469cef0508e1e2c39482875ae90665bcc07d96821df293a9746863050` |
| Requirements recovery baseline | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` |
| Requirements reconciliation | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` |
| Corrected master plan rev 3 | `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` |
| TTS governing design | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` |
| TTS §17.19 adoption | `9da36691054df40d98c772f7ba9643bc959416e39a7f2d30681d9404cd16f043` |
| TTS §17.19 draft | `dc5852adbc70641dbe8f233851eaec44266330fce197b78a0ab2641e7f73f784` |
| T2.9.5B r14 register | `bd10519bd5e3611fd984c0835f3eac74b04e67818caa8055675967ec5f94cc44` |
| T2.9.5B r14 freeze record | `e1d1df8180ecf266799a1a50040b96d71bae4a5516c776b23c22b76d802206d0` |

---

## 11. Current contradictions and how they are handled

| Contradiction | Resolution without rewriting history |
|---|---|
| r14 register header says `DRAFT — NOT FROZEN`; freeze record says frozen | The freeze record is the later, specific control. `docs/FROZEN-ARTIFACTS.md` records both hashes and the controlling rule. |
| TTS §17.19 draft header says `NOT ADOPTED`; adoption record says adopted | The adoption record is the later, specific control. The draft remains historical. `docs/AMENDMENT-REGISTER.md` records the relationship. |
| Amendment 02 rev 9 is closed; rev 6/7/8 still exist | Rev 6/7/8 are archived revisions. The amendment register marks them historical and points to rev 9. |
| Corrected master plan rev 3 says technical scope/dependencies are "frozen" but is not adopted | "Frozen" there describes internal revision scope, not owner adoption. The plan remains proposed. This distinction is recorded in the registry and README. |
| `docs/review/README.md` described only T2.9.5B audio | The README now begins with a directory-wide scope notice and points to this chain. |
| DROOT header cites inventory hash `cd722…`; actual inventory is `343104…` | The registry records the actual hash. The DROOT file is evidence and is not edited, because the disposition request cites its current `494539…` hash. This is a recorded provenance discrepancy, not a reopened decision. |

No frozen body was edited to resolve these. The overlay documents carry the administrative status.

---

## 12. Current state and remaining gaps

### Established

- F3-C architecture and data/ownership contracts are frozen.
- F3-C implementation is not accepted.
- ROOT-01 establishes source-level non-conformance with `(rootID, chatID)`.
- The six F3-C defects are open.
- Amendment 02 is closed at revision 9, not architecture approval.
- T2.9.5B r14 is the frozen decision-control register; the stage is not authorized.
- The corrected master plan is proposed, not adopted.

### Open

- D-ROOT has no owner-selected canonical root.
- The six defects need six independent dispositions.
- Runtime evidence for packaged roots, aliases, collisions, cross-chat memory, and ID reuse is not established.
- Candidate/continuity reachability is unestablished.
- Amendment 02 `M6-01`–`M6-03` are unresolved evidence obligations.
- The corrected master plan is not adopted.
- T2.9.5B stage authorization is not granted.

### Deliberately outside this chain

- **Howl-specific memory-system evaluation.** `docs/research/MEMORY-01-RESEARCH-SYNTHESIS-2026-11-12.md` and its source register are post-chain external evidence. They do not select an architecture or close runtime gaps. This chain remains local project evidence and governance only.
- Implementation authorization for any queue item.
- A new architecture selection.

---

## 13. How to use the chain

- To answer "why is this rule here?", start at the stage and follow the predecessor.
- To answer "is this current?", check the registry's current-relevance column and the amendment/freeze overlay.
- To answer "can I change this?", check `docs/FROZEN-ARTIFACTS.md`. If frozen, no.
- To answer "where is the hash?", check the registry row or the hash anchors above.
- To answer "what is missing?", read section 12 and `docs/README.md` §8.

The chain is complete enough to trace the actual work. It is not a substitute for the work, and it does not convert evidence into authorization.
