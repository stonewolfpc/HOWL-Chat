# Howl Chat — M0A Release Truth Acceptance Review

## Review `HCDM-M0A-AR-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Acceptance-criteria review; documentation only |
| Program map | `HCDM-2026-11-21-01` |
| Starting scope baseline | `HFS-AUDIT-2026-11-21-01` |
| Related | `HFR-CP-2026-11-21-01`; `HF1-PLAN-2026-11-21-01` |
| Scope reviewed | M0A: T0.1, T0.2, T0.5 |
| Authorization | **None**. This review defines a gate; it does not authorize implementation |
| Tests executed by this review | **0** |
| Production changes | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Verdict

**M0A-AR-01 (Calculated — recommendation).** M0A is the correct first implementation package **only if** its acceptance criteria include:

1. the actual frontend send path, not only a direct backend call;
2. the unawaited frontend `AppendChatMessage` / backend `recordUserTurn` race;
3. a real two-instance process-safety check for T0.5;
4. an explicit stop rule that prevents T0.2 from becoming a chat-architecture redesign.

With those conditions, M0A is narrow, evidence-based, and worth authorizing after owner review. Without them, it is too weak to close the release-truth questions and too easy to turn into an early rewrite of Phase 1.

**M0A-AR-02 (Calculated — no authorization).** This review does **not** authorize implementation, code changes, test changes, a model run, or a two-instance run.

---

## 2. Why M0A comes before Phase 1

**M0A-AR-03 (Observed/Source-derived).**

- T0.1 is an owner-reported failure: a sent message does not appear until restart.
- T0.2 is an owner-reported failure: the model receives the user message twice.
- T0.5 is a process-ownership safety gate: concurrent instances must not kill each other’s model children.
- Phase 1 will change chat state, prompt assembly, persistence and frontend actions. If the send/render path is still untrusted, Phase 1 can “fix” state while leaving a layer that hides or duplicates the user’s message.

**M0A-AR-04 (Calculated — not a redesign).** M0A is a behavioral truth pass. It must not introduce stable state IDs, branch semantics, token budgets, multi-character isolation, or a new persistence model.

---

## 3. What the current code actually does

### 3.1 Frontend send path

**M0A-AR-05 (Observed).** In `frontend/dist/js/chat-input.js`:

- The active chat is created before the bubble is rendered (comment around lines 218–233).
- `addUserMessage(message)` renders the user bubble at line 235.
- `window.go.main.App.AppendChatMessage('user', message, 'User').catch(...)` is invoked at lines 246–247 **without `await`**.
- The stream starts through `streamMultiCharacterResponses(message)` at line 253 or the image path before it.
- `SaveChat()` is invoked at lines 257–258, also without `await`.

**M0A-AR-06 (Observed — backend).** In `app-chat.go`:

- `SendMessage` calls `recordUserTurn(message)` at line 175.
- `SendMessageStream` delegates to `SendMessageStreamWithImage`.
- In `app-prompt.go`, `SendMessageStreamWithImage` calls `recordUserTurn(message)` at line 58.

**M0A-AR-07 (Observed — guard).** `recordUserTurn` (`app-chat.go:360–376`) skips a new user turn only when the **immediately preceding** in-memory message is a user turn with the identical content.

**M0A-AR-08 (Observed — append).** `AppendChatMessage` (`app-prompt.go:277–310`) appends to `a.chatMessages` and then writes a `ChatHistory` entry to `history.ndjson`.

### 3.2 The race that M0A must test

**M0A-AR-09 (Observed/Inferred — T0.2 risk).** The frontend does not await `AppendChatMessage` before starting the stream. If `SendMessageStreamWithImage` reaches `recordUserTurn` before the append call has updated `a.chatMessages`, the guard sees no preceding identical turn and appends a second copy. The later append then adds its own copy. The guard is correct for a sequential duplicate; it is not a proof against the asynchronous UI path.

**M0A-AR-10 (Calculated — acceptance consequence).** A test that calls `AppendChatMessage` synchronously and then `SendMessageStream` in one goroutine may pass while the real UI path still duplicates. M0A must exercise the **unawaited** call ordering, or reproduce the race deterministically with a controllable append delay.

### 3.3 T0.1 render path

**M0A-AR-11 (Observed).** `addUserMessage(message)` runs after active-chat creation and before `AppendChatMessage`. That ordering is intended to prevent the active-chat creation from clearing the bubble. It does not prove same-session rendering under a real send, a failed append, a failed stream, or a chat switch. No current test named for T0.1 user-message rendering was found in the inspected test set.

---

## 4. Exact M0A acceptance criteria

**M0A-AR-12 (Calculated — all criteria are required).** A criterion passes only with the stated evidence. “Looks right” and a single screenshot do not close it.

### A1 — Baseline single-chat send

| Field | Requirement |
|---|---|
| Setup | Packaged app; one active chat; generic AI path (no multi-character scenario); no image attachment; a loaded model or an approved deterministic inference double |
| Action | Send a unique token such as `HCDM-M0A-<run-id>` once |
| Pass | The user bubble appears in the same session, before restart; the input clears once; no duplicate user bubble appears |
| Evidence | Screenshot/log of the session; active chat ID; timestamped message |

### A2 — Persistence count is exactly one

| Field | Requirement |
|---|---|
| Setup | Same run as A1 |
| Action | Read the active chat’s `history.ndjson`; do not clear caches |
| Pass | Exactly **one** user entry with the unique token; exactly the expected assistant entry count; no duplicate user entry hidden in `history.json` |
| Evidence | Saved file excerpt/record count; chat ID; no cache-clearing step |

### A3 — Chat switch and restart

| Field | Requirement |
|---|---|
| Setup | Same run as A1 |
| Action | Switch to another chat and back without restarting; then restart the app and open the same chat |
| Pass | Exactly one user bubble before and after restart; no duplicate; no disappearance; no need to clear caches |
| Evidence | Screenshots before/after switch and restart; saved history count |

### A4 — Generic prompt contains the token once

| Field | Requirement |
|---|---|
| Setup | The existing `recordingInference` seam (`app_prompt_boundary_test.go:14–58`) or equivalent inference-boundary recorder |
| Action | Simulate the frontend sequence: `AppendChatMessage('user', token, 'User')` then `SendMessageStreamWithImage(token, "")`; also test a direct `SendMessageStreamWithImage` without a prior append |
| Pass | In both paths, the final `[]llama.ChatMessage` contains the token exactly once in a `user` role; the assistant reply does not create a second user occurrence |
| Evidence | Captured message list; role/content counts; test name and result |

### A5 — Async append race

| Field | Requirement |
|---|---|
| Setup | A controllable append seam that can delay the `AppendChatMessage` completion while the stream begins, or a deterministic equivalent of the unawaited frontend call |
| Action | Start the stream before the append finishes; repeat with the append completing before the stream reaches `recordUserTurn` |
| Pass | Each ordering produces exactly one user occurrence in the prompt; the in-memory history and `history.ndjson` contain exactly one user entry |
| Evidence | Deterministic race test; both orderings recorded; no timing-only “usually passes” claim |
| Stop rule | If the only safe fix requires changing the chat state contract or stable-ID design, **stop M0A** and hand the issue to Phase 1A. Do not redesign chat state inside M0A |

### A6 — Real request check

| Field | Requirement |
|---|---|
| Setup | One existing local model; one controlled request; no model download |
| Action | Send one user turn and capture the llama-server request/prompt log or the app’s exact outbound request |
| Pass | One user occurrence per turn; no duplicate role/content pair |
| Evidence | Request log or exact boundary capture; model identity; timestamp |
| Limitation | If no model can be run, A6 is **not passed**; the review must say so. A4/A5 can still prove the app boundary, but the release-truth claim remains limited |

### A7 — Two independent instances, normal startup

| Field | Requirement |
|---|---|
| Setup | Two independent packaged instances, A and B; each has an owned model child or an explicitly approved owned test child started through the production process path |
| Action | Start A; wait for its child; start B; wait for B’s child; observe A’s child for a bounded period |
| Pass | A’s child remains alive after B starts; B does not kill, replace, or reap A’s child; no path-based sweep is used |
| Evidence | PIDs, parent/child relation, job ownership, process table snapshots before/after B starts, app logs |

### A8 — Abrupt exit and normal close

| Field | Requirement |
|---|---|
| Setup | Same two instances as A7 |
| Action | Terminate A’s parent abruptly; observe A’s child and B’s child; then close B normally |
| Pass | A’s child exits as a consequence of A’s job closing; B’s child remains alive after A dies; B’s child exits after B closes; no unrelated process is killed |
| Evidence | Before/after PID/job observations; exit codes; no pattern-based kill command |

### A9 — Scope guard

| Field | Requirement |
|---|---|
| Prohibited | No stable-ID/branch/token redesign, no character-count fallback, no image-generation/TTS/audio-generation work, no visual-character stream, no MEMORY-09 integration, no schema/API migration, no cache clearing as a fix |
| Pass | Diff is limited to the send/render/append/process-ownership path required by A1–A8, or the package closes as verification-only |
| Evidence | Diff summary; explicit out-of-scope statement |

---

## 5. What counts as closure

**M0A-AR-13 (Calculated — closure states).**

| Outcome | Closure |
|---|---|
| A1–A3 pass on the current build | T0.1 closes as **verified**, with no code change |
| A4–A6 pass | T0.2 closes as **verified**, with no code change |
| A4/A5 fail because of the unawaited append race | A minimal append/guard fix is in M0A; no state contract change |
| A5 requires a new state/identity contract | **Stop**; M0A does not redesign; Phase 1A owns it |
| A7/A8 pass | T0.5 closes as **release-safety verified** |
| A7/A8 fail | Fix process ownership only; do not redesign resource management |
| A6 unavailable | Report T0.2 as **partially verified at the app boundary**; do not claim real-request closure |

**M0A-AR-14 (Observed — no test run).** This review did not execute A1–A9. They are proposed acceptance criteria, not results.

---

## 6. Specific risks and stop conditions

**M0A-AR-15 (Calculated).**

| Risk | Why it matters | Stop/response |
|---|---|---|
| **Unawaited append race** | A synchronous test can pass while the real UI duplicates | A5 must test both orderings; fix only append/guard |
| **`addUserMessage` appears fixed while persistence is wrong** | A bubble without a saved turn is still a release failure | A2/A3 require saved-history evidence |
| **Duplicate guard masks a second append** | Guard only checks consecutive identical in-memory turns | A5 must inspect both memory and disk |
| **Test double does not represent the real UI ordering** | False confidence | Use a controllable append seam; A6 real request where possible |
| **Two-instance test itself reaps children** | Harness contaminates the result | Use recorded PIDs/jobs; never pattern-kill |
| **T0.5 model-gated test mistaken for two-instance proof** | `TestD3ChatServerSpawnOwnership` is sequential and model-gated | A7/A8 require independent instances |
| **M0A expands into Phase 1** | Delays usable delivery and creates architectural debt | A9 scope guard; stop at A5 contract boundary |
| **Screenshot-only closure** | Does not prove prompt/history truth | A2/A4/A5 require counts and boundary evidence |

---

## 7. Relationship to Phase 1 and M0B

**M0A-AR-16 (Calculated).**

- Phase 1 remains **proposed, not authorized**.
- Phase 1A (identity/restore) begins after M0A, or after M0A establishes that the send path is not being rewritten underneath it.
- M0B (ghost and ownership truth) can run after M0A in a separate UI lane. It should not be bundled into M0A.
- M0A must not pre-implement Phase 1 IDs or token accounting.

**M0A-AR-17 (Calculated — M0B boundary).** M0B should:
- close T2.2 ghost IDs;
- triage `DEF-FE-01`;
- inventory frontend/backend write authority for chat, lorebook, character and legacy memory;
- fix only local, safe ownership inconsistencies;
- leave state/identity redesign to Phase 1A and lorebook cutover to M5A.

---

## 8. Recommendation

**M0A-AR-18 (Calculated — recommendation only).** M0A is ready for **owner review** with criteria A1–A9. If the owner accepts these criteria, the package is narrow enough to authorize as the next implementation step.

It should **not** be expanded into a general chat cleanup. Its purpose is to establish:

1. the message appears once, immediately, and persistently;
2. the prompt carries it once, including through the unawaited frontend path;
3. two instances do not kill each other’s children.

**M0A-AR-19 (Observed — no authorization).** No implementation is authorized by this review. No code, tests, schema, model, process run, or protected record was changed.

---

## 9. Verification record

**M0A-AR-20 (Observed).**

| Check | Result |
|---|---|
| Production changes | **0** |
| Test changes | **0** |
| Tests run | **0** |
| Model/process run | **0** |
| Protected/canonical changes | **0** |
| This file | New, untracked; acceptance review |
| Canonical hashes | r33 `dfb28b60…` · v9 `698254ca…` · ledger `e5008020…` · map `b4ba33cc…` |

**M0A-AR-21 (Observed — evidence paths).**

| Claim | Path |
|---|---|
| Frontend send sequence | `frontend/dist/js/chat-input.js:215–258` |
| Add user bubble | `frontend/dist/js/chat-input.js:235` |
| Unawaited append | `frontend/dist/js/chat-input.js:246–247` |
| Stream start | `frontend/dist/js/chat-input.js:253` |
| Backend send | `app-chat.go:168–198` |
| Record user turn | `app-chat.go:360–376` |
| Append history | `app-prompt.go:277–310` |
| Stream record | `app-prompt.go:58` |
| Inference recorder | `app_prompt_boundary_test.go:14–58` |
| Process ownership | `internal/proc/job_windows.go` |
| Ownership tests | `internal/proc/job_windows_test.go:40–112`; `d3_chat_ownership_test.go:38` |

**End of review.** M0A should be judged on exact behavioral evidence, not on a code diff or a screenshot. The most important addition to the existing M0A definition is the asynchronous frontend append test: it is the difference between closing T0.2 and merely closing a simpler backend duplicate.
