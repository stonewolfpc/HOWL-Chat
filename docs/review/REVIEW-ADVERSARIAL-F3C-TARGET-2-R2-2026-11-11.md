# REVIEW — Fresh Adversarial Review of F3C-TARGET-2 R2

**Status:** **ADVERSARIAL REVIEW — NOT APPROVAL, NOT FREEZE, NOT IMPLEMENTATION.**  
**Date:** 2026-11-11  
**Reviewed artifact:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R2-2026-11-11.md`  
**Reviewed R2 SHA-256:** `ba989dcc87d8788421e5cb554dc893863644a28f121f30dc1d5ede3352e82fd4`  
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**Prior review:** `docs/review/REVIEW-ADVERSARIAL-F3C-TARGET-2-2026-11-11.md`  
**Scope:** read-only review. No source edits, tests, builds, model loads, migrations, destructive tests, implementation, or freeze.

---

## 1. Overall verdict

**R2 is ready for further revision, not for freeze.**

R2 resolves the shape of all ten prior findings, and several resolutions are sound:

- intent ownership now includes deletion, tombstone, wipe, and root incarnation;
- commit linkage fields exist;
- a recovery precedence is stated;
- hot-swap has a real state machine;
- F3-F is scoped to attachment durability rather than all F3-C work;
- V24 is present;
- v2 optional-field rules are substantially clearer.

But R2 still contains **three critical residual defects**:

1. it treats an absent chat as if it were tombstoned, conflicting with the frozen `ErrChatStateUnknown` rule;
2. it accepts “an assistant commit exists” without proving that the commit is the intended attempt;
3. deletion and wipe do not explicitly quiesce an active generation stream before removing or staging state.

Those are not documentation polish. Each admits a different plausible implementation and a different data outcome. The proposal is not safe to freeze until they are corrected.

**Evidence note:** R2 is a proposal, so its rules are level 5 (unverified) by default. Frozen-contract rules are level 1. Current source references are level 2. There are no level 3 test-proven or level 4 runtime/native results in this review.

---

## 2. Evidence classifications

| Classification | Meaning | Present here |
|---|---|---|
| 1. Contractually established | Frozen Master contract rule. | Yes, cited where used. |
| 2. Source-inspected | Current source behavior observed read-only. | Yes, for delete, wipe, stream, model load, character generation, and storage. |
| 3. Test-proven | Reproducible test evidence with positive and negative cases. | None. No tests were run. |
| 4. Runtime/native-verified | Real Wails/WebView2 or media runtime evidence. | None. No model loads or native runs. |
| 5. Unverified/inferred | Proposed rule, design inference, or future requirement. | R2 itself and all Gate rows. |

No proposed R2 rule is promoted to an established fact in this review.

---

## 3. New findings

### N1 — Absent chat is treated as tombstoned, contrary to `ErrChatStateUnknown`

**Severity: Critical. Defect class: contract.**

**Exact R2 rule:** §2.6, recovery table:

> `Tombstoned/absent chat` → “Delete/terminalize intent; no target, no recreation.”

**Invariant:** A missing chat without a tombstone is unknown state, not proof of deletion. Deletion authority must come from a tombstone, not from absence.

**Evidence:** The frozen contract says: “missing chat without tombstone is `ErrChatStateUnknown`” (`F3C-MASTER...:722`). The current `DeleteChat` removes the directory (`persistence.go:471`), but the contract’s delete path requires a tombstone before removal. R2 conflates “absent” with “tombstoned.”

**Adversarial counterexample:** A chat directory is temporarily unavailable because a root is mounted at the wrong path, a previous copy was incomplete, or an offline volume is not yet attached. No tombstone exists. Recovery sees a generation intent, treats the chat as deleted, deletes the intent, and loses the only evidence of the interrupted attempt. If the chat later reappears, the attempt is gone without a deletion record.

**Does the rule survive?** No. It converts an unknown-state condition into an irreversible cleanup.

**Required correction:** Separate the cases:

- tombstone present → clean intents; no recreation;
- chat directory absent and no tombstone → `ErrChatStateUnknown`; no cleanup, no target, no intent deletion;
- root identity/incarnation mismatch → block, not clean (see N12);
- wipe intent present → wipe precedence.

**Gate rows:** V8 (delete/tombstone), V22 (interrupted generation).

---

### N2 — “Assistant commit exists” is not tied to the exact attempt

**Severity: Critical. Defect class: contract.**

**Exact R2 rule:** §2.6:

> “Assistant commit/receipt exists” → “Mark `committed`; clean; never append again.”

R2 introduces `assistantMessageId`, `assistantCommitOperationId`, and `assistantCandidateFingerprint` in §2.2, but the recovery table does not require them to match the found commit.

**Invariant:** Recovery must adopt the exact attempt, not merely any assistant commit that happens to exist.

**Adversarial counterexample:** The intent belongs to attempt A. A later retry or a neighboring character produces attempt B at the same chat/base. Recovery sees B’s receipt and marks A `committed`. A’s candidate is lost or duplicated; the intent’s `assistantMessageId` and fingerprint are ignored. A different implementation searches only by chat/revision and reaches the opposite conclusion.

**Does the rule survive?** No. Existence is not identity.

**Required correction:** Recovery must verify all of:

- root/root incarnation;
- chat ID;
- operation ID and operation type;
- assistant commit operation ID;
- assistant message ID;
- candidate fingerprint;
- base revision and session/epoch;
- receipt/record body and target manifest;
- no conflicting commit for the same attempt.

Mismatch → conflict/block, no adoption. A receipt for a different candidate is not evidence for this attempt.

**Gate rows:** V22, V18, V13.

---

### N3 — Delete and wipe do not quiesce an active generation stream

**Severity: Critical. Defect class: contract/implementation.**

**Exact R2 rule:** §2.1 defines `DeleteChat` as tombstone → terminalize/remove intents → remove directory. It does not require canceling or awaiting an active stream. §6.3 and the media rules do not address wipe while a stream is active. The hot-swap section says wait for a stream, but deletion/wipe are separate operations.

**Invariant:** No generation write may occur after the chat is tombstoned or the root is staged for wipe.

**Evidence:** Current streaming records chunks and appends only after completion (`app-prompt.go:25-146`; character path `:160-244`). Current `WipeAllData` stages chat and gallery roots but does not cancel a stream (`app-wipe.go:1-125`). A late assistant append can occur after deletion/wipe begins.

**Adversarial counterexample:** A stream is active. The user deletes the chat. `DeleteChat` writes a tombstone, removes the intent, and removes the directory. The stream completes and `appendChatMessage`/commit runs against the deleted chat, creating an orphan target or resurrecting a directory. Wipe has the same race: staging the root does not stop the stream, and a late write can land in the staged/removed tree.

**Does the rule survive?** No. Intent cleanup is not stream quiescence.

**Required correction:** Delete and wipe must, under the same locks:

1. reject or cancel the active generation;
2. wait for the generation attempt to reach a terminal intent state;
3. ensure no append/commit can start afterward;
4. then tombstone/delete or stage/wipe.

An active stream cannot be silently completed after deletion. If deletion is requested during generation, report the abort and the unsaved status.

**Gate rows:** V8, V12, V15, V22.

---

### N4 — First-v2 migration installs record/receipt before proving installed targets

**Severity: High. Defect class: contract.**

**Exact R2 rule:** §5.4:

1. v1 remains current;
2. build and stage v2 targets;
3. read back and verify;
4. verify v1 base/receipt/chain;
5. **create record/receipt** with a v1 predecessor and v2 manifest;
6. **install and publish v2**.

The order is ambiguous: a record/receipt created before the v2 targets are installed can prove a manifest that is not yet the installed state. Conversely, installing v2 targets before record/receipt can leave unproven v2 bytes if the process crashes.

**Invariant:** A receipt must prove the exact installed target set, not a staged candidate.

**Adversarial counterexample:** Crash after step 5 but before step 6. The record/receipt says v2; the current target files are still v1. Recovery must either reject a valid-looking receipt or install unproven targets. Another implementation installs targets first and has the opposite crash window. The contract cannot verify both.

**Does the rule survive?** No, as written. The sequence needs a single defined commit point and recovery table.

**Required correction:** Define the exact atomic sequence and crash result, for example:

- stage all v2 targets and verify readback;
- install targets create-only/atomic;
- verify installed bytes;
- write record;
- write receipt;
- publish;
- define pre-install, post-install/pre-record, post-record/pre-receipt, and post-receipt recovery.

No mixed manifest is valid.

**Gate rows:** V4, V13, V23.

---

### N5 — Model identity is still deferred, not frozen

**Severity: High. Defect class: contract.**

**Exact R2 rule:** §4.2:

> “The exact registry/digest implementation must be frozen in the model identity record before TARGET-2 freeze.”

That is a deferral. R2 does not define the algorithm, the identity record, or the fallback when no registry identity exists.

**Invariant:** A historical assistant turn must name a stable model identity that survives moves, display-name changes, and removal.

**Adversarial counterexample:** Two implementations both claim `modelId` compliance: one uses a registry UUID, another uses a path hash; one hashes the GGUF file, another hashes metadata only. The same model produces different IDs, or two different quants share an ID. Historical provenance is not comparable, and “model unavailable” cannot be resolved reliably.

**Does the rule survive?** No. “Must be frozen later” cannot be the freeze condition for the target format.

**Required correction:** Freeze one identity algorithm before TARGET-2 freeze. At minimum define:

- registry identity when present;
- digest identity and canonicalization when absent;
- model file/quant change behavior;
- identity record location and lifecycle;
- behavior when identity cannot be computed (no generation, not a placeholder ID).

**Gate rows:** V18, V21.

---

### N6 — The parameter set may omit actual generation controls

**Severity: High. Defect class: contract.**

**Exact R2 rule:** §4.3 fixes six keys: `temperature`, `topP`, `topK`, `repeatPenalty`, `maxTokens`, `contextTokens`.

**Invariant:** The provenance describes the effective runtime configuration actually used, not merely a convenient subset.

**Adversarial counterexample:** Two turns have identical six recorded values but different presence penalty, frequency penalty, stop sequences, grammar, sampler order, or template behavior. The output differs, but the durable provenance is identical. A reviewer cannot explain or reproduce the turn. Conversely, a control may be absent or not applicable, and R2 provides no absence/default representation.

**Does the rule survive?** No. The six-key set is not proven complete for the actual request surface.

**Required correction:** Freeze the complete effective request identity, or add every generation-relevant control with exact absence/default semantics. A canonical request fingerprint over the actual engine payload is preferable to an incomplete named subset. Leave the choice to the owner, but do not call the six-key object “actual runtime configuration” without the missing-control rule.

**Gate rows:** V18.

---

### N7 — Strict unknown-field blocking can permanently brick legacy chats

**Severity: High. Defect class: contract/owner decision.**

**Exact R2 rule:** §5.2: a v1 file with unknown target fields is read-only/openable but cannot upgrade; `ErrLegacyUnknownField`; a repair/export path “may be defined later.”

**Invariant:** Preserve legacy data without silently dropping it, while allowing a legitimate upgrade.

**Adversarial counterexample:** A harmless legacy key such as `legacyNote` exists in one `meta.json`. Every future write is blocked. The user can open the chat but never save it again; there is no repair/export path in the contract. The strict rule protects the field but destroys the workflow.

**Does the rule survive?** It satisfies “do not silently drop,” but not “complete migration.” It is a deliberate trade-off, not a resolved design.

**Required correction:** Owner decision:

- **lossless preservation where possible:** a defined sidecar/`legacyExtras` or quarantine store, with proof that the target bytes retain the field; or
- **strict block only where preservation cannot be proven**, with a concrete repair/export path; or
- **explicitly accept permanent read-only legacy chats** and remove the upgrade claim.

Do not state “resolved” without choosing one.

**Gate rows:** V23, V4.

---

### N8 — V1 predecessor-byte non-retention leaves rollback undefined

**Severity: High. Defect class: contract/owner decision.**

**Exact R2 rule:** §5.3: after a successful v2 commit, v1 target bytes are not retained; only the receipt/manifest/hash chain remains.

**Invariant:** If the product claims a complete predecessor or rollback, the predecessor must be reconstructable. Hash-only evidence proves identity, not content.

**Adversarial counterexample:** A v2 commit is later found corrupt or semantically wrong. The v1 receipt says what the old bytes hashed to, but those bytes are gone. There is no way to restore the v1 chat. A “rollback” feature can only fail. The same problem applies to a user who wants the pre-v2 file for inspection.

**Does the rule survive?** Only if the owner explicitly accepts no byte-level predecessor and no rollback after v2 commit. It does not satisfy “complete predecessor” as written.

**Required correction:** Owner decision:

- retain an immutable v1 predecessor snapshot outside the target set, with retention and wipe rules; or
- explicitly accept hash/receipt-only history and remove/qualify any complete-predecessor or rollback claim.

**Gate rows:** V23.

---

### N9 — Preference-update semantics remain ambiguous

**Severity: High. Defect class: contract.**

**Exact R2 rule:** §3.1 says only a successful selection may update `savedContinuation`; §3.2 says B “may become saved” and allows an “explicit preference-only action” to save an unavailable model. It does not define the operations, ordering, or whether a normal successful load automatically saves.

**Invariant:** Saved preference, requested model, and active model must be independently represented and never silently change.

**Adversarial counterexample:** B loads successfully. One implementation automatically saves B; another waits for a separate explicit “use for this chat” action. The same user action produces different future chat behavior. Or a preference-only action saves unavailable B while active A is loaded; the send path now tries B and fails, whereas the user expected A to continue.

**Does the rule survive?** No. “May become saved” is not a contract.

**Required correction:** Define exact operations:

- `SelectModel` → active after ready; does it save? choose yes/no;
- `SetContinuationPreference` → explicit preference, allowed unavailable, with no active claim;
- ordering and failure behavior for preference commit;
- send behavior when saved ≠ active and saved is unavailable;
- no automatic fallback.

**Gate rows:** V21.

---

### N10 — Gate mapping claims V1–V24 but only maps a subset

**Severity: High. Defect class: evidence/test.**

**Exact R2 rule:** §9 says “The Gate v2 row inventory becomes V1–V24,” but the table includes only V4, V5, V7, V8, V11, V12, V13, V15, V16, V18–V24. V1, V2, V3, V6, V9, V10, V14, and V17 have no disposition, claim, or evidence update.

**Invariant:** A Gate row that is not explicitly dispositioned is not tested and cannot support a freeze.

**Adversarial counterexample:** A lock-order defect (V6), a symlink/reparse escape (V14), or a Mimo challenge gap (V17) survives while the R2 table claims a complete V1–V24 inventory. The implementation can pass the listed rows and still fail the missing ones.

**Does the rule survive?** No.

**Required correction:** Produce the full V1–V24 table. For unchanged rows, state “unchanged from Gate v2,” the exact claim, and why the amendment does not affect it. For changed rows, provide claim/evidence/counterexample. V17 cannot be omitted.

**Gate rows:** V1–V24.

---

### N11 — One active attempt per chat may change multi-character behavior

**Severity: Medium. Defect class: contract/implementation.**

**Exact R2 rule:** §2.3 serializes generation attempts per chat; §2.7 says character generation uses the same protocol; §7 says multi-character behavior is unchanged and V24 requires parity.

**Invariant:** A serialization rule must not silently remove an existing capability.

**Evidence:** `GenerateCharacterResponse` starts a stream and calls `beginStream` (`app-prompt.go:160-244`); `beginStream` sets a cancellation token but does not itself reject a second generation (`app-lifecycle.go:166`). Multiple character responses can therefore be invoked without an explicit one-active gate in the current code. R2 would impose a new rejection/queue rule.

**Adversarial counterexample:** A multi-character group requests two character responses. R2 rejects or queues the second; current behavior may run both. The user-visible interaction changes, but V24 claims parity.

**Does the rule survive?** Not unless the owner explicitly accepts the behavior change or the existing product is proven to be serialized.

**Required correction:** Decide:

- one active attempt per chat, with V24 updated to acknowledge the deliberate serialization; or
- multiple active attempts keyed by operation/character, with recovery and provenance per attempt.

Either is possible; silence is not. V24 must test the chosen behavior.

**Gate rows:** V15, V24.

---

### N12 — Root-incarnation mismatch cleans evidence instead of failing closed

**Severity: Medium. Defect class: contract.**

**Exact R2 rule:** §2.5 step 1: root identity/incarnation mismatch → “orphan; no adoption”; §2.6: “Root incarnation mismatch | Orphan; clean; no target.”

**Invariant:** A root-incarnation mismatch is an identity failure, not proof that the intent is disposable. Cleaning evidence on a wrong mount can destroy the only record.

**Adversarial counterexample:** The app starts against the wrong root or a copied root. The incarnation does not match. Recovery treats every intent as an orphan and cleans them. The real root’s interrupted attempts are lost, and the mismatch is hidden.

**Does the rule survive?** No, as a blanket cleanup rule.

**Required correction:** Distinguish:

- an intent explicitly quarantined as belonging to an old incarnation may be cleaned under a defined policy;
- a current-root identity/incarnation mismatch blocks with `ErrRootIncarnationChanged` and preserves evidence;
- no automatic cleanup merely because the current root does not recognize the intent.

**Gate rows:** V13, V22.

---

### N13 — Candidate fingerprint and commit operation semantics are undefined

**Severity: Medium. Defect class: contract.**

**Exact R2 rule:** §2.2 requires `assistantCandidateFingerprint` and `assistantCommitOperationId` but does not define their hash domain, canonicalization, allocation, retry, or failure behavior.

**Invariant:** Exact linkage must be reproducible across implementations and recoverable after a failed commit.

**Adversarial counterexample:** One implementation fingerprints the candidate target manifest, another fingerprints the in-memory message; recovery cannot compare. A failed commit reuses the same operation ID, producing an idempotency conflict or a duplicate. The intent is not a reliable anchor.

**Does the rule survive?** No.

**Required correction:** Freeze the fingerprint domain/canonicalization, operation-ID allocation and reuse rules, and the failed-commit retry/abandon semantics. The fingerprint must bind the exact candidate that the receipt proves.

**Gate rows:** V13, V22.

---

### N14 — New-chat and initial-creation format are not explicitly v2

**Severity: Medium. Defect class: contract.**

**Exact R2 rule:** §5.1 says the first new write creates a v2 revision; §5.4 covers legacy v1→v2; §5.6 covers Save As. It does not explicitly state that a brand-new chat created after the amendment is v2, or what manifest/targets creation uses.

**Invariant:** Every new revision has one deterministic target format.

**Adversarial counterexample:** One implementation creates a new chat at v1 because there is no v1 predecessor; another creates v2. Both can claim R2 compliance. A later reader sees different formats for equivalent new chats.

**Does the rule survive?** No.

**Required correction:** State explicitly: new chat creation, Save As, and every post-amendment write use `F3C-TARGET-2`/`tgt2-`; v1 is only for legacy current revisions and historical proof.

**Gate rows:** V4, V9, V23.

---

### N15 — V24 has no baseline evidence

**Severity: Medium. Defect class: evidence/test.**

**Exact R2 rule:** §7 requires “existing character-context and cast suites” and parity, but no baseline artifact, command, or current evidence is identified.

**Invariant:** A regression row must compare against a captured pre-amendment baseline; otherwise it proves nothing.

**Adversarial counterexample:** The named suites do not exist, or they assert a new behavior rather than the old one. V24 passes by label while character context or cast behavior changes.

**Does the rule survive?** No evidence yet. It is a requirement, not a verified result.

**Required correction:** Require a captured baseline (prompts/context/cast outputs) before implementation, exact suite/evidence paths, and a native run. Do not mark V24 pass from a test name.

**Gate rows:** V24, V16.

---

## 4. Disposition of the ten prior findings

| Prior finding | R2 disposition | Residual issue |
|---|---|---|
| R1 — intent outside deletion/wipe | **Partially resolved.** Intent area is root-owned and wipe/tombstone-aware. | N1 absent-chat cleanup, N3 stream quiescence, N12 root mismatch. |
| R2 — no commit linkage/switch | **Partially resolved.** Fields and switch rule added. | N2 recovery does not verify the exact commit; N13 fingerprint semantics. |
| R3 — third journal precedence | **Partially resolved.** Precedence is stated. | N1, N3, N10; ordinary/cast combinations need explicit reconciliation; full Gate mapping missing. |
| R4 — rollback aspiration | **Partially resolved.** State machine and failure matrix added. | N9 preference semantics; N3 delete/wipe; LoadModel readiness-vs-accepted remains a choice. |
| R5 — provenance not actual | **Not resolved.** Capture point is improved. | N5 model identity deferred; N6 parameter set may omit controls. |
| R6 — v1 predecessor/unknown fields | **Partially resolved.** Strict policy and staged write are stated. | N4 install order, N7 deadlock, N8 no predecessor bytes, N14 new-chat format. |
| R7 — media scope | **Resolved in principle.** Core/full split and V20 block are correct. | F3-F evidence absent; V20 remains blocked. |
| R8 — multi-character boundary | **Resolved in principle.** Non-goal and V24 added. | N11 one-active rule may change behavior; N15 no baseline. |
| R9 — optional-byte rules | **Partially resolved.** Optional/unknown rules improved. | N6 parameter completeness; N14 new-chat format; no serializer evidence. |
| R10 — non-streaming/character paths | **Resolved in principle.** All three paths are named. | N11 one-active behavior; N2/N13 recovery linkage still applies. |

No prior finding is marked fully closed merely because R2 names it.

---

## 5. Challenge results A–H

### A. Generation intent, deletion, recovery

**Result: not safe to freeze.**

- Intent creation after user commit is coherent in principle.
- Deletion ownership is improved.
- N1 makes absent-chat recovery destructive.
- N2 makes commit adoption non-specific.
- N3 allows a late stream write after deletion/wipe.
- N12 cleans on root mismatch.
- The precedence is a direction, not yet a total, evidence-backed recovery algorithm.

### B. Hot-swap transaction

**Result: improved, still incomplete.**

- Preflight, capture, stop, start, health, publish, rollback, and unavailable are now named.
- Wait-for-stream is a coherent policy.
- N9 leaves the saved-preference update operation ambiguous.
- The delete/wipe race is N3.
- `LoadModel` returning after ready versus an accepted operation is still an alternative; the contract must choose one.
- No crash/native evidence exists; this remains level 5.

### C. Provenance and generation attempts

**Result: improved, not complete.**

- Capture at attempt start is correct.
- Effective values are the right target.
- N5 leaves model identity undefined.
- N6 leaves the parameter set incomplete for the actual request.
- Retry/partial rules are coherent.
- N11 leaves multi-character concurrency undefined.

### D. V1/V2 migration

**Result: not ready for an owner decision to be inferred.**

- N4 leaves the commit/install order ambiguous.
- N7’s strict unknown-field rule can permanently block writes.
- N8’s no-byte retention means no complete predecessor/rollback.
- N14 leaves new-chat format unstated.
- Old receipt/manifest chain preservation is sound.

### E. Media and F3-F

**Result: boundary is sound; evidence is absent.**

- F3-C-CORE can proceed without F3-F.
- F3-C-FULL and V20 remain blocked.
- References are not durability.
- Missing/corrupt/external media is unavailable, not substituted.
- No F3-F contract or runtime evidence exists in this review.

### F. Multi-character isolation

**Result: boundary is stated, but not fully verified.**

- No cast/prompt/memory/context-flush changes are intended.
- V24 is the right regression shape.
- N11 may change concurrency behavior.
- N15 has no baseline evidence.
- Persistence preserves `characterId`; it does not redefine what a character sees.

### G. Serialization and all generation paths

**Result: improved, not byte-complete.**

- Optional/unknown rules are clearer.
- N6 and N14 leave provenance/new-chat bytes under-specified.
- No serializer evidence exists.
- All three stated paths are covered in principle.

### H. Gate v2 completeness

**Result: incomplete.**

- R2 claims V1–V24 but maps only a subset (N10).
- V17 (Mimo challenge) is omitted.
- V6 (locks), V14 (filesystem), and V9/V10 (creation/memory) have no disposition.
- No row passes; V20 is blocked; V16 and V24 lack native/baseline evidence.

---

## 6. Owner trade-offs

### D1. Unknown-field blocking

R2 chooses strict blocking to avoid silent data loss. That protects the field but can permanently brick a legacy chat if any harmless unknown key exists, and it provides no repair path.

**Recommendation:** Prefer lossless preservation where the target can represent it without weakening validation; block only where preservation cannot be proven, and require a concrete repair/export path. Do not silently accept permanent read-only chats.

**Owner decision required:** strict block, lossless sidecar, or explicit read-only policy.

### D2. V1 predecessor-byte non-retention

R2 explicitly does not retain v1 target bytes. That is honest, but it means:

- “complete predecessor” is not true at byte level;
- rollback after a v2 commit cannot reconstruct v1;
- only the receipt/manifest/hash chain remains.

**Recommendation:** If the trust promise includes rollback or reconstructability, retain an immutable v1 predecessor snapshot outside the target set. If not, explicitly state “hash/receipt-only; no byte-level rollback” and remove any stronger claim.

**Owner decision required:** retain bytes or accept hash-only/no rollback.

---

## 7. Gate v2 V1–V24 audit

No row is marked PASS. This is a contract review, not a test run.

| Row | Disposition | Reason |
|---|---|---|
| V1 | UNCHANGED / NOT TESTED | Canonicalization/hash chain; R2 does not alter; no tests. |
| V2 | UNCHANGED / NOT TESTED | Receipt inventory; no tests. |
| V3 | UNCHANGED / NOT TESTED | Commit fork/orphan; no tests. |
| V4 | UNRESOLVED | N4 install order, N14 new-chat format, exact bytes. |
| V5 | UNRESOLVED | N3 stream/delete, N9 preference, switch/commit. |
| V6 | UNMAPPED | Lock ordering for intent/wipe not defined. |
| V7 | UNRESOLVED | N7 unknown fields, Save As. |
| V8 | UNRESOLVED | N1 absent chat, N3 stream quiescence. |
| V9 | UNMAPPED | Creation journal and new-chat v2 format. |
| V10 | UNCHANGED / NOT TESTED | Memory migration; no tests. |
| V11 | UNRESOLVED | N7 unknown fields, v1/v2 history. |
| V12 | UNRESOLVED | N3 wipe/stream, intent staging. |
| V13 | UNRESOLVED | N4/N12/N13 fault boundaries. |
| V14 | UNMAPPED | Filesystem recovery; no amendment disposition. |
| V15 | UNRESOLVED | N3/N11 races. |
| V16 | UNVERIFIED | No native evidence. |
| V17 | UNMAPPED | Mimo challenge not addressed. |
| V18 | UNRESOLVED | N5 identity, N6 parameters. |
| V19 | UNRESOLVED | Semantics clear, no evidence. |
| V20 | BLOCKED | F3-F media ownership absent. |
| V21 | UNRESOLVED | N9 preference, N5 identity, hot-swap evidence. |
| V22 | UNRESOLVED | N1/N2/N3/N12/N13. |
| V23 | UNRESOLVED | N4/N7/N8/N14. |
| V24 | UNRESOLVED | N11/N15 baseline. |

Gate v2 itself is not amended or frozen by this review. V20 is blocked, not failed.

---

## 8. Decisions requiring explicit owner approval

1. **Absent chat/root mismatch:** fail closed with `ErrChatStateUnknown`/root error, or allow cleanup. Recommended: fail closed.
2. **Unknown fields:** strict block, lossless sidecar, or explicit permanent read-only policy.
3. **V1 predecessor bytes:** retain an immutable snapshot or accept hash-only/no rollback.
4. **Model identity:** exact registry/digest algorithm and fallback.
5. **Effective parameter set:** six keys, full request fingerprint, or expanded set.
6. **Preference semantics:** automatic save after ready versus explicit `SetContinuationPreference`; send behavior when saved ≠ active.
7. **One active attempt per chat:** preserve current multi-character concurrency or deliberately serialize.
8. **New-chat format:** explicitly v2 for all post-amendment creation.
9. **LoadModel API:** return after ready or accepted operation plus `model:ready`.
10. **Gate v2 scope:** full V1–V24 mapping, including V6/V9/V10/V14/V17.

---

## 9. Document and authority integrity

- Frozen authority SHA-256 verified: `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`.
- R2 SHA-256: `ba989dcc87d8788421e5cb554dc893863644a28f121f30dc1d5ede3352e82fd4`.
- R2 line count: 533.
- R2 fenced-block count: 14; balanced.
- No source, test, build, migration, model load, runtime, or application-data changes were made.
- No frozen-contract, R2, prior-review, or source file was edited.

---

## 10. Status

**R2 proposal:** not frozen, not accepted; ready for further revision, not freeze.  
**Frozen contract authority:** verified, unchanged.  
**F3-C phase 0:** remains **NOT CLOSED**.  
**P1a:** remains **UNAUTHORIZED** and downstream of P1c.  
**F3-C acceptance:** **NOT GRANTED**.  
**F3-D:** remains **BLOCKED**.  
**F3-B:** remains closed.  
**A9/A10/A12/A13:** remain closed and are not reopened.

The next action is to revise R2 against N1–N15 and return it for a fresh review. No implementation authorization exists.
