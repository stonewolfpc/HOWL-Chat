# REVIEW — Adversarial Review of F3C-TARGET-2 Proposal

**Status:** **ADVERSARIAL REVIEW — NOT APPROVAL, NOT FREEZE, NOT IMPLEMENTATION.**  
**Date:** 2026-11-11  
**Reviewed artifact:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-2026-11-11.md`  
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**Verification discipline:** Gate v2, contract-first review.  
**Scope:** document review and source trace only. No implementation, tests, builds, model loads, source changes, or freeze.

---

## 1. Verdict

**The proposal is directionally sound but not ready for freeze.**

The four ownership decisions are correct in principle:

- F3-C owns durable chat identity and references.
- F3-F owns media bytes.
- Hot-swap is runtime orchestration, not target state.
- Stream completion is not durable commit.

The proposal still leaves five areas under-specified enough that a Gate v2 reviewer could not distinguish a correct implementation from a plausible one. The two most serious are:

1. the generation-intent journal is outside the deletion/tombstone/wipe transaction and has no explicit commit linkage;
2. the hot-swap “restore the old model” path is a desired outcome, not a defined transition with a recoverable state.

The proposal should be revised before freeze. No P1a work, no tests, and no implementation should start from it in its current form.

---

## 2. What holds

These parts should survive revision:

- **Version coexistence.** `F3C-TARGET-1` remains valid; `tgt1-` receipts are not invalidated or rehashed. A v1 chat can open unchanged, and the first new write creates a v2 revision with a v1 predecessor.
- **Media separation.** F3-C stores references only; F3-F owns bytes. A reference does not prove durability. Gate V20 is correctly blocked until F3-F evidence exists.
- **Readiness before publication.** A selected model must not become the active model until health is established. The current `LoadModel` return-before-health window is a real false-success risk.
- **Stream versus commit.** A partial assistant response must not be classified as completed. A durable-success event must be distinct from a stream-end event.
- **No runtime state in the target.** `loading`, `starting`, `ready`, and `unavailable` are runtime states, not chat fields.
- **Existing preflight preserved.** GGUF admission, chat-task eligibility, and projector pairing checks happen before the current server is stopped. Those checks should not be rewritten.

---

## 3. Findings

### R1 — Generation-intent journal is outside the deletion and wipe ownership boundary

**Severity: Critical.**

The proposal places generation intents under a root-level operation area:

```text
f3c-generations/<F3C-PATH-GEN-1 hash>/<operationID>.json
```

It says the intent is not a target and not a receipt, but it does not say how it is owned by deletion, tombstone, chat switching, or wipe.

The frozen contract already defines:

- deletion as a tombstone before directory removal (`F3C-MASTER...:690-716`);
- a tombstone inventory with exactly one tombstone per chat;
- wipe staging that includes `chats/`, `f3c-commits/`, `f3c-receipts/`, `f3c-tombstones/`, `f3c-creation-intents/`, `f3c-memory-migrations/`, the cast journal if present, and the gallery root (`:771-805`).

The proposed `f3c-generations/` is not in that wipe list. Current `PersistenceManager.DeleteChat` removes only the chat directory (`internal/backend/chat/persistence.go:471`). A root-level intent can therefore outlive the chat it describes.

**Risk:** recovery can encounter an intent for a deleted or tombstoned chat. Depending on recovery order, it can block the chat, report a phantom interrupted generation, or attempt to attach a commit to a chat that no longer exists. The proposal does not define a tombstone/existence check, so “resurrection” is an undefined behavior, not a safe one.

**Required revision:**

- generation intents must be under a root/chat ownership path covered by tombstone and wipe;
- delete must terminalize or remove intents, or recovery must treat any intent whose chat is tombstoned/absent as an orphan and clean it without creating a target or receipt;
- wipe must stage and hash the generation-intent root explicitly;
- recovery must check `chatID`, `rootID`, base revision, and tombstone state before any adoption;
- a deleted chat’s intent must never be adopted, reported as live, or used to create a revision.

**Gate rows:** V12 (wipe), V22 (interrupted generation), plus deletion recovery under V8.

---

### R2 — The intent has no explicit commit linkage or switch policy

**Severity: Critical.**

The proposal’s intent fields are:

```text
version, rootId, chatId, operationId, operationType,
baseRevision, userMessageId, userMessageHash,
modelId, startedAt, status
```

`status = completed` is not evidence of a commit. Nothing links the intent to:

- the user-turn commit operation;
- the assistant commit operation;
- the receipt/record that proves the assistant turn;
- the state session/epoch that was active at start;
- the switch/discard decision that may have invalidated it.

A crash after the assistant commit but before the intent is marked `completed` is indistinguishable from a crash after a partial stream unless recovery searches receipts by convention. A crash after a chat switch can leave an intent for the old chat and a new active session.

**Required revision:**

Add at least:

```text
stateSessionId, baseEpoch, userCommitOperationId,
assistantCommitOperationId, terminalReason
```

and define:

- the intent is bound to one chat/root/session and one base revision;
- a switch cancels the stream and terminalizes/marks the intent for the outgoing chat;
- no intent is reused across chats;
- completion requires a matching assistant receipt/record, not just a status field;
- recovery distinguishes “commit exists,” “commit missing,” “base stale,” “session invalidated,” and “chat deleted”;
- ordinary autosave commits that advance the base do not silently adopt the intent.

**Gate rows:** V5 (append/autosave/switch), V22 (interrupted generation), V15 (races).

---

### R3 — Generation intent creates a third journal without a recovery-order rule

**Severity: High.**

The contract’s recovery rule says:

- validate ordinary and cast journals;
- both nonterminal → `ErrDualJournalAmbiguous`;
- then choose recovery under the ordinary phase table (`F3C-MASTER...:825-875`).

A generation intent is a third journal. The proposal does not say whether it is:

- subordinate to an ordinary autosave journal;
- a separate operation type with its own recovery phase table;
- checked before or after the ordinary journal;
- allowed to coexist with a pending ordinary save;
- an error when both are nonterminal.

Without an order, two reviewers can implement different recovery behavior and both claim the contract.

**Required revision:**

Define a single recovery order, for example:

1. validate root and tombstones;
2. validate ordinary save/append journal;
3. validate generation intent;
4. validate cast journal;
5. refuse ambiguous combinations;
6. only then adopt or clean.

The intent must not be a competing “commit” journal. It is an attempt journal whose terminal state is resolved against the ordinary commit/receipt evidence.

**Gate rows:** V13 (fault injection), V22 (interrupted generation), V3 (fork/orphan).

---

### R4 — Hot-swap rollback is not a defined transaction

**Severity: Critical.**

The proposal says:

> Attempt one restore of the previous model; if restore succeeds, old selection remains active. If restore fails, enter `unavailable`.

That is the right desired outcome, but it is not enough to implement or verify. Missing details:

- What exactly is saved before stop? The previous model path, projector, spawn settings, admission/eligibility evidence, and resource estimate are all needed. A model path alone is not enough.
- What happens if the old server did not stop cleanly? The current code prints a warning and continues; the proposal says stop failure must not start a new server, but does not define the old state.
- What happens if the new server spawned but never became healthy? It must be killed before restore; otherwise port/process conflicts make the restore path nondeterministic.
- What happens if the old restore also fails? `unavailable` is correct, but the proposal must say whether the old preference is retained, whether the UI shows “requested but not loaded,” and whether the global last-model path is updated.
- What happens on a crash during the transition? On restart there is no active model and no durable transition state. The proposal does not define the recovery/UI state.
- What if the active stream is cancelled versus allowed to finish? The proposal says “cancel or wait,” which is two contracts, not one.
- What if `LoadModel` returns “accepted” rather than “ready”? The proposal allows either. A freeze needs one API/event contract.

**Required revision:**

Define a runtime-only transition record containing the previous and requested model identity/settings, transition ID, and phase. It is not target state. Then define a single failure matrix:

| Condition | Required result |
|---|---|
| Preflight fails | Old model untouched; no new selection; no stop. |
| Stop fails | Do not start replacement; health-check old; report stop failure; do not claim new. |
| Spawn fails | Kill any partial process; attempt one restore with saved settings; otherwise unavailable. |
| Health fails/times out | Stop failed process; attempt one restore; otherwise unavailable. |
| Restore fails | `unavailable`; retain last preference; no fallback model; no false loaded. |
| Crash mid-transition | No durable active selection; startup discovers no healthy server; no automatic substitution. |
| Active stream | Choose cancel-and-record-aborted, or wait. Do not leave both. |

Also choose one API contract: either `LoadModel` returns only after health, or it returns an accepted operation and a separate `model:ready` event is the only publication signal. The frontend must not set `howl-last-model-path` before `ready`.

`App.serverProc`, the HTTP client, `resource.Manager`, and `model.Manager` must be reconciled to one runtime truth after every branch. “Reconcile” needs an owner and an order, not just a sentence.

**Gate rows:** V15 (races), V16 (native), V21 (continuation/hot swap), V13 (fault injection).

---

### R5 — Per-turn provenance does not yet prove actual runtime configuration

**Severity: High.**

The proposal correctly says “actual request/generation result, not the picker.” But it leaves the critical details open:

- `modelId` is still “stable registry/digest identity”; no algorithm is frozen. Two models with the same display name, a moved file, and an updated quant can all be represented differently by different implementers.
- `parameters` is “relevant generation controls” without an exact field set or ordering. A target hash cannot be stable if the set is implementation-defined.
- It does not distinguish **requested** parameters from **effective** parameters after model defaults, clamping, template injection, or backend normalization.
- It does not say when the snapshot is captured. Capturing at save time can record a later global setting; capturing before the model resolves defaults can record a request that was not used.
- Partial generations: the proposal says no assistant entry, but does not define whether a partial attempt gets a provenance record. It should not get a committed assistant record; the intent can carry the attempted model.
- Retries: a retry of the same user turn must create a new attempt/operation and a new assistant turn. The proposal does not prohibit merging or overwriting the prior attempt.
- Mid-stream model switch: the in-flight generation must keep the model it started with. The proposal’s transition rules prohibit stopping during a stream, but provenance should explicitly freeze the model at attempt start.
- Character responses: the proposal does not enumerate every generation path. `GenerateCharacterResponse` and non-streaming `SendMessage` need the same attempt/provenance rules; otherwise multi-character turns can bypass persistence.

**Required revision:**

- freeze a concrete `modelId` algorithm (registry identity preferred; otherwise a documented digest over identity-bearing model artifacts, not path/display name);
- freeze the exact `parameters` key set, value types, canonical order, and whether values are requested/effective;
- record effective values after defaults/clamping and before the first token;
- bind the provenance snapshot to the attempt at start, not at save;
- define retry as a new attempt with a new operation ID and no overwrite;
- state that partial attempts have no committed assistant turn;
- cover `SendMessage`, `SendMessageStreamWithImage`, and character generation paths.

**Gate rows:** V18 (provenance), V19 (seed), V21 (mid-stream), V24 (multi-character regression, proposed).

---

### R6 — V1-to-V2 migration does not define unknown-field or predecessor-byte preservation

**Severity: High.**

The proposal says v1 remains valid and old receipts are not rehashed. That is necessary but not sufficient for the owner’s “complete V1 predecessor and all unknown fields” requirement.

Open questions:

- The v1 target has a fixed field set. If a legacy `meta.json` or `history.json` has unknown keys, the current decoder drops them. The proposal does not say whether v2 must preserve them, quarantine them, or reject them.
- “Complete V1 predecessor” can mean at least three different things:
  1. the v1 receipt/record/hash chain remains valid;
  2. the v1 target **bytes** are retained somewhere;
  3. all unknown source fields are preserved.
  The proposal only clearly promises (1).
- The first v2 write replaces target files. If it is interrupted, the old current v1 state must remain recoverable. The proposal says ordinary write, but does not define a v1→v2 migration phase, staging, or rollback. Existing commit rules may cover it, but Gate V23 needs an explicit claim.
- `Save As` from v1 creates v2; no receipt is copied as v2 proof. That is correct, but unknown fields and source binding need a stated policy.
- The proposal does not define whether a v2 writer may rewrite old v1 history entries to add `attachments: []` or must leave them untouched. Byte-level compatibility depends on this.

**Required revision:**

Choose and freeze one policy:

- **Strict:** unknown target fields are rejected/quarantined, not silently dropped; or
- **Lossless:** a defined `legacyExtras`/sidecar preserves unknown fields; this is a v2 schema addition and must be explicit.

Also state whether historical v1 target bytes are retained. If not, say “historical proof is hash/receipt-only; predecessor bytes are not retained.” Do not let the phrase “complete predecessor” imply something the storage model does not provide.

Define migration as a staged first-v2 commit with the same create-only/readback/receipt protocol, and specify the recoverable outcome for interruption at each point.

**Gate rows:** V4 (target bytes), V11 (history migration), V13 (fault injection), V23 (v1/v2 compatibility).

---

### R7 — Media dependency scope is still ambiguous

**Severity: High.**

The proposal says F3-C cannot claim attachment durability before F3-F. It also says F3-F is a prerequisite in `OD2`. Those are not the same scope.

The correct distinction is:

- F3-F is **not** a prerequisite for the entire F3-C implementation;
- F3-F **is** a prerequisite for attachment durability and for any final F3-C acceptance that includes the promise that images come back;
- F3-C can implement text, context, provenance, continuation, and interrupted recovery while attachment references are `reference-only` or `unavailable`;
- Gate V20 must be `BLOCKED`, not `PASS`, until F3-F provides ownership evidence;
- F3-C must not claim “complete workspace persistence” while V20 is blocked.

The proposal’s current wording can be read as blocking unrelated work. That would be an unnecessary dependency inversion.

**Required revision:**

Split acceptance:

1. **F3-C text/context/provenance acceptance** — may proceed without F3-F, provided attachments are explicitly unavailable/reference-only.
2. **F3-C full workspace acceptance** — requires F3-F for attachment durability and V20.
3. **F3-C claim language** — no “attachments survive” claim until F3-F.

Do not block the target-format freeze on F3-F; block only the attachment acceptance claim. This keeps the persistence work moving without repeating the false-durability mistake.

**Gate rows:** V20 (blocked), V16 (native), plus the overall acceptance record.

---

### R8 — Multi-character isolation is not stated as an explicit non-goal and regression

**Severity: High.**

The proposal changes `ChatHistory`, `ChatMeta`, generation events, and recovery. Those are shared by multi-character flows. The owner requires no changes to:

- character-specific context flushing;
- character memory isolation;
- character-specific prompt construction;
- cast synchronization;
- multi-character interaction behavior.

The proposal says F3-C owns persistence, but it does not state this boundary in its scope or acceptance. An implementer could “helpfully” use the new `characterId`/provenance fields to alter cast sync or prompt assembly.

**Required revision:**

Add an explicit non-goal:

> F3C-TARGET-2 preserves the information character systems need; it does not define or change what each character sees, how character context is flushed, how character memory is isolated, how character prompts are built, or how cast synchronization behaves.

Add a regression row to Gate v2 for multi-character context/prompt parity. The new `generation` record is per assistant turn; `characterId` remains the existing owner field. The continuation model remains per chat, not per character. No cast semantics are amended.

**Gate rows:** V24 (proposed multi-character regression), V16 (native), existing character-context suites.

---

### R9 — V2 exact-byte rules are incomplete for optional fields

**Severity: Medium.**

The proposal gives field order but leaves several byte-level questions:

- Is `generation` omitted or `null` when legacy-unknown? The proposal says omitted, which is correct, but the serializer must enforce it.
- Is `attachments` always `[]` for v2, including migrated v1 entries? The proposal says yes for “v2 entries,” but does not say whether a migrated v1 entry becomes a v2 entry in the same commit.
- Are `continuationModel` fields all optional except `modelId`? What is the exact omission rule for `sourceHint` and `displayName`?
- Does `parameters` allow unknown keys? If yes, target bytes vary by implementation. If no, an unknown parameter must fail.
- Are `seed` and `effectiveSeed` integers with a bounded range? What is the encoding for a random seed reported as a string?
- Does a v2 entry with `hasImage: true` and `attachments: []` require `legacy_unknown` marker, or is the inconsistency itself the marker?
- Does `history.ndjson` include the same optional fields as `history.json` with identical presence rules?

Without answers, two implementations can both claim `F3C-TARGET-2` and produce different bytes.

**Required revision:** a precise optional-field table and serializer tests in the amended contract, not just prose.

**Gate rows:** V4 (target bytes), V18/V20 (field semantics), V23 (migration).

---

### R10 — Generation intent is not covered for non-streaming and character paths

**Severity: Medium.**

The proposal’s attempt journal is described around streaming. Current source has at least:

- `SendMessage` synchronous generation (`app-chat.go:170`);
- `SendMessageStreamWithImage` (`app-prompt.go:25`);
- `GenerateCharacterResponse` (`app-prompt.go:160`).

All three can produce an assistant turn and all three can fail or be interrupted. If only the streaming path writes intents, a non-streaming crash or character-generation failure is invisible. If all three write intents, the operation/recovery rules must be shared.

**Required revision:** require one generation-attempt protocol for every assistant-producing path, with the same intent, commit, and recovery semantics. Character-specific prompt/cast behavior remains outside the amendment.

**Gate rows:** V18, V22, V24.

---

## 4. Answers to the five challenges

### 4.1 Generation-intent journal: commits, switching, deletion, wipe

**As proposed: incomplete and unsafe to freeze.**

- It is outside the tombstone/wipe inventory (R1).
- It has no commit linkage (R2).
- It is a third journal without a recovery order (R3).
- It does not define switch cancellation or orphan cleanup.

**Required before freeze:** intent must be chat/root/session-bound, tombstone-aware, wipe-staged, operation-map-integrated, and linked to the user/assistant commit operations. Recovery of a deleted chat must clean the intent, not reconstruct the chat.

### 4.2 Model transition rollback

**As proposed: desired outcome, not a complete transaction.**

The old model can be restored only if its path, projector, settings, and admission evidence are saved before stop. If stop succeeds and start fails, one restore may work; if restore fails, `unavailable` is the only honest state. There is no durable rollback guarantee and should not be one in the target.

**Required before freeze:** choose cancel-or-wait for active streams, define the runtime transition record, kill failed processes before restore, reconcile all runtime authorities, and never publish a selection before health. A crash during transition must leave no false active model.

### 4.3 Per-turn provenance: actual runtime configuration, partials, retries

**As proposed: directionally correct, not reproducible.**

The record must be captured at attempt start from the effective configuration actually handed to the engine, not from the frontend settings at save time. It needs a concrete model identity and parameter set. Partial attempts should have an intent but no committed assistant turn. Retries must be new attempts with new operation IDs and no overwrite. Mid-stream model changes must not affect the in-flight attempt. Character and non-streaming paths must be included.

**Required before freeze:** freeze identity, parameter schema, capture point, and attempt/retry semantics.

### 4.4 V1-to-V2 migration

**As proposed: preserves the receipt chain, not necessarily the predecessor bytes or unknown fields.**

Old receipts and `tgt1-` hashes are correctly preserved. The proposal does not define unknown-field handling, historical target-byte retention, or an interrupted first-v2 migration. The phrase “complete V1 predecessor” is therefore too strong.

**Required before freeze:** choose strict/lossless unknown-field policy, state whether predecessor bytes are retained, and define a staged, recoverable first-v2 commit. Gate V23 must test interruption at every phase.

### 4.5 Media dependency

**As proposed: too broad if read as “F3-F blocks all F3-C.”**

F3-F is required for attachment durability and for full workspace acceptance. It is not required for text/context/provenance/continuation persistence. The right split is:

- F3-C text/context acceptance can proceed with attachments explicitly unavailable;
- F3-C full acceptance and V20 remain blocked until F3-F;
- no attachment-durability claim may be made before F3-F.

This preserves unrelated progress without weakening the attachment promise.

---

## 5. Multi-character boundary

The amendment must preserve, not redefine:

- character-specific context flushing;
- character memory isolation;
- character-specific prompt construction;
- cast synchronization semantics;
- multi-character interaction behavior.

The new persistence fields are additive:

- `characterId` remains the existing owner field;
- `generation` is per assistant turn;
- `attachments` is per message;
- `continuationModel` is per chat, not per character;
- no generation-intent field may alter cast state.

The review should be able to say: “F3-C preserves the inputs character systems need; it does not decide what a character sees.” Gate v2 should add a regression row for multi-character prompt/context parity, not a new cast contract.

---

## 6. Gate v2 disposition

No row is passing. This is a contract review, not a verification run.

| Row | Disposition | Reason |
|---|---|---|
| V4 target bytes | **UNRESOLVED** | Optional-field rules and v1/v2 exact bytes need freezing. |
| V5 append/autosave/switch | **UNRESOLVED** | Intent/switch/commit relationship undefined. |
| V7 Save As | **UNRESOLVED** | v1→v2 source binding and unknown fields need definition. |
| V8 deletion | **UNRESOLVED** | Generation intents are outside tombstone ownership. |
| V11 history migration | **UNRESOLVED** | v1/v2 field presence and unknown fields. |
| V12 wipe | **UNRESOLVED** | `f3c-generations` missing from staged inventory. |
| V13 fault injection | **UNRESOLVED** | Intent phases and migration interruption not defined. |
| V15 race testing | **UNRESOLVED** | Stream/switch/transition policy is ambiguous. |
| V16 native | **UNRESOLVED** | No native evidence; not a code review claim. |
| V18 provenance | **UNRESOLVED** | Model identity, effective parameters, capture point. |
| V19 seed | **UNRESOLVED** | Effective random seed and encoding. |
| V20 attachments/F3-F | **BLOCKED** | No F3-F contract/evidence; reference is not durability. |
| V21 continuation/hot swap | **UNRESOLVED** | Rollback/crash/stream contract incomplete. |
| V22 interrupted generation | **UNRESOLVED** | Intent ownership, commit link, recovery order. |
| V23 v1/v2 compatibility | **UNRESOLVED** | Unknown fields, predecessor bytes, interrupted migration. |
| V24 multi-character regression | **PROPOSED** | Boundary must be explicit; no cast/prompt changes. |

Gate v2’s rule still applies: no row passes on a test name, a positive case alone, or a description of intent. V20 cannot pass without F3-F ownership evidence.

---

## 7. Required revisions before freeze

1. Make generation intents tombstone/wipe/operation-map citizens.
2. Add user/assistant commit linkage and switch/session/epoch binding.
3. Define the generation-intent recovery order relative to ordinary and cast journals.
4. Define a concrete hot-swap transition record, rollback matrix, crash behavior, and one stream policy.
5. Freeze model identity, effective parameter schema, capture point, partial/retry rules, and all generation paths.
6. Freeze unknown-field and predecessor-byte policy for v1→v2.
7. Define a staged, recoverable first-v2 migration.
8. Split F3-C text/context acceptance from full attachment acceptance.
9. Add the multi-character non-goal and regression row.
10. Freeze v2 optional-field serialization rules.
11. Update Gate v2 V4/V5/V7/V8/V11/V12/V13/V15/V16 and add V18–V24.
12. Re-run the adversarial review after revision; only then consider freeze.

---

## 8. Final status

**Adversarial verdict: directionally sound; not ready for freeze.**

No implementation, tests, builds, model loads, source changes, or contract changes were performed. This review is not an approval and does not authorize P1a.

Phase 0 remains **not closed**. P1a remains unauthorized and downstream of P1c. F3-C is not accepted. F3-D remains blocked. A9/A10/A12/A13 remain closed and are not reopened.

---

## Appendix — evidence index

### Proposal

- `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-2026-11-11.md`
  - target model: §4
  - provenance/seed: §5
  - continuation: §6
  - hot-swap: §7
  - media/F3-F: §8
  - interrupted generation: §9
  - legacy: §10
  - affected surfaces: §11
  - Gate v2: §12
  - decisions: §13

### Frozen contract

- `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`
  - §8.3 full/compound saves, `:670`
  - §8.4 append/replace/clear, `:678`
  - §8.5 streaming/dirty, `:686`
  - §8.6 delete/tombstone, `:694`
  - §8.10 wipe, `:771`
  - §9.1 journal recovery, `:825`
  - §9.2 ordinary phases, `:840`
  - §11 acceptance, `:946`
  - §12 known gaps, `:1018`

### Current source

- `internal/backend/chat/persistence.go:471` — current `DeleteChat` removes the chat directory; no generation intent ownership.
- `app-wipe.go:1-125` — current wipe stages chat and gallery roots; no `f3c-generations`.
- `app-chat.go:96-111` — load continues after failed outgoing save.
- `app-chat.go:170` — synchronous `SendMessage`.
- `app-chat.go:314-327` — `performSave`.
- `app-chat.go:355-381` — `recordUserTurn` in-memory only.
- `app-prompt.go:25` — streaming path.
- `app-prompt.go:138-146` — assistant append, `chat:complete`, async save.
- `app-prompt.go:160` — character generation path.
- `app-model.go:28-174` — preflight then stop.
- `app-model.go:317` — spawn.
- `app-model.go:345-512` — monitor/register/return before readiness.
- `app-model.go:537` — unload.
- `app-lifecycle.go:221` — stop server.
- `internal/backend/chat/types.go:50`, `:72` — no continuation model; history lacks provenance.
- `internal/backend/types/message.go:31-43` — metadata exists separately.
- `frontend/dist/js/chat-model.js:347-393` — selection after load return; global last path.
- `frontend/dist/js/chat-input.js:20`, `:50`, `:240` — transient image base64.
- `frontend/dist/js/chat-actions.js:62` / `chat-session.js:146` — only `hasImage`.
- `internal/image/gallery/gallery.go:94-128` — PNG+sidecar writer.

### Gate v2

- `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md:201-245` — V1–V17, row evidence, isolated destructive testing.
- V18–V24 are proposed by this review and are not frozen.
