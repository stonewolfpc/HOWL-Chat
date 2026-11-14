# REVIEW — Independent Adversarial Review of F3C-TARGET-2 R5

**Status:** **INDEPENDENT ADVERSARIAL REVIEW — NOT APPROVAL, NOT FREEZE, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Reviewed artifact:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md`
**R5 SHA-256:** `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`
**R5 size:** 1,475 lines; 77,233 bytes.
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`
**Frozen authority size:** 1,091 lines; 41,270 bytes.
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`
**Gate v2 SHA-256:** `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`
**Prior independent review:** `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R4-2026-11-11.md`
**Prior independent review SHA-256:** `7c884af2c44daade26d96ba8bef527a18ac71894b573eda8541af31cf032d964`
**Prior proposal:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R4-2026-11-11.md`
**Prior proposal SHA-256:** `7a2c221147621706e299924c20de1fe29080b7d649fac25264ce212d11e7417a`
**Scope:** documentation-only adversarial proposal review and read-only source trace. No source, test, build, migration, model load, destructive test, application-data, authority-record, freeze, acceptance, P1a, or F3-D action was performed or authorized.

---

## 1. Overall verdict

**R5 is a materially stronger and more coherent proposal than R4, but it is not freeze-ready.**

R5 correctly treats the fourteen R4 findings as symptoms of one connected authority problem. Its most important architectural choices are sound:

- immutable per-commit packages containing recoverable target bytes;
- leases as durable authorization records inside the frozen lock order;
- explicit session, epoch, cast, and character fencing;
- immutable response plans separated from mutable batch/slot state;
- separation of operation and generation fingerprints;
- one authoritative continuation preference;
- Gate v2 V1–V17 preserved as adopted and V18–V24 preserved as proposed;
- explicit acknowledgment that OD-02/AM-1, OD-14, OD-15, and multiple evidence gaps remain unresolved.

Those are genuine improvements. They nevertheless leave **six freeze-blocking contract defects**, several high-severity exactness gaps, and every Gate row unrun. The proposal should remain:

- **PROPOSAL ONLY;**
- **NOT FROZEN;**
- **NOT ACCEPTED;**
- **NOT IMPLEMENTED;**
- **NOT VERIFIED.**

The requested owner decisions must not be inferred from architectural improvement. In particular, OD-02/AM-1 must remain narrowly scoped, OD-14 must not be treated as decided, OD-15 and V24 must remain blocked, and no Gate row may be treated as passed.

---

## 2. Answer to the four review instructions

### 2.1 OD-02/AM-1 is correctly exposed, but its present scope is unsafe

R5 correctly refuses to reinterpret frozen Master §6.9. It states that ordinary v2 revision replacement conflicts with the frozen create-only installation rule and therefore requires formal amendment AM-1. It also correctly preserves create-only semantics for commit packages, commit records, receipts, tombstones, and predecessor archives while retaining the replaced target bytes in the immutable package.

That framing is a substantial improvement over R4. It is not yet narrow enough.

R5 §1.3 simultaneously says:

> all new evidence files remain create-only.

That universal statement is false within R5 itself. R5 expressly creates:

- mutable `batch.json` aggregate and slot state;
- authority records containing mutable fence, barrier, lease-count, and recovery fields;
- lease records with mutable lifecycle statuses;
- attempt and candidate records with multiple nonterminal and terminal statuses.

The frozen layout also does not contain the new R5 roots:

- `f3c-packages/`;
- `f3c-authority/`;
- `f3c-generations/`;
- `f3c-predecessors/`.

Therefore, approving only AM-1 as “target projection replacement” would leave two unacknowledged changes:

1. a new durable-layout expansion;
2. a new class of mutable operation-evidence state transitions.

A future implementation could reasonably claim that OD-02 authorized all of this because R5 presented AM-1 as the only required amendment.

**Required owner treatment:**

- Split the decision.
- OD-02a should authorize only byte-verified target-projection replacement from an already durable commit package, with the prior target bytes retained in that package.
- A separate layout/state decision should authorize every new root and enumerate every mutable operation-evidence file.
- That decision must list permitted status transitions, forbidden transitions, write protocol, recovery treatment, inventory treatment, and wipe treatment for leases, authority records, plans, batches, attempts, and candidates.
- AM-1 must not authorize replacement of commit records, receipts, tombstones, packages, archives, plans, attempts, candidates, or mutable batch/authority records.
- “Same-volume atomic replacement” must not be approved as a multi-file atomic primitive.

### 2.2 The lease model is substantially better, but four stale-authority questions remain open

R5 answers the general shape of the lease questions correctly:

- a revoked or expired lease cannot commit through retry, renewal, replay, or process restart;
- fencing changes must increment `fenceEpoch` under `chatWriteMu`;
- lower-epoch leases must be marked revoked before the lock is released;
- commit validation must include token, incarnation, session, epoch, cast, fence, operation, base, package, and record checks.

The remaining defects concern **when validation happens, which state changes are fencing events, and how process identity and time are established**.

#### Finding C2 — Expiry can lapse between validation and the commit point

R5 validates the lease during the commit sequence, then performs potentially lengthy durable work—building target bytes, writing and flushing a package, reading the package back, installing it, writing and installing a commit record, installing target projections, and installing a receipt—while retaining the frozen locks.

Time can therefore pass after a successful expiry check and before the atomic record installation. No fencing thread can intervene because the same `chatWriteMu` is held, but wall-clock expiry is not a thread-mediated change. The proposal does not require a final validation immediately before record installation and does not name a monotonic clock.

**Concrete sequence:**

1. A holder receives a lease expiring at T+120 seconds.
2. It acquires `chatWriteMu` and passes validation at T+20 seconds.
3. Package construction, flushing, verification, or target-projection installation stalls for more than 100 seconds.
4. Its nominal authority expires while it still holds the lock.
5. It installs the commit record after expiry.
6. The durable history gains a post-expiry commit despite the stated lease lifetime.

The same long critical section also prevents the holder from performing its own heartbeat renewal if renewal requires reentry to the already-held `chatWriteMu`.

**Required correction:**

- specify a monotonic clock source for lease lifetime;
- make the expiry comparison uniform;
- repeat the complete commit-gate validation immediately before atomic record installation;
- prohibit staging, flushing, verification, or projection I/O after that final validation;
- bound the final critical section by size, duration, or operation class;
- define how long staging may remain outside the final gate and how it is inventoried and cleaned.

This is freeze-blocking because it concerns the actual commit boundary, not merely lease issuance.

#### Finding C3 — Slot abandonment is not a lease-fencing event

R5 handles character removal much better than R4 because removal increments both `castEpoch` and `characterVersion`, and those values are carried into candidate, attempt, batch, lease, and commit validation.

Slot abandonment has no equivalent protection.

The lease-revocation trigger list covers root, barrier, session, cast, character-version, fence, operation/token, and process-incarnation changes. It does not include batch-slot abandonment, slot removal, or batch invalidation. The universal commit-time checklist likewise does not require batch, slot, plan, or abandonment validation.

Section 7 says `AbandonBatchSlot` invalidates an attempt, but that is not automatically a write-lease revocation. R5 separately says generation attempts do not hold write leases and that a durable candidate/commit transaction obtains a new write lease. Consequently, abandonment can occur before that new lease is issued.

**Concrete sequence:**

1. Slot 2 has completed generation but has not yet obtained a candidate/commit write lease.
2. An operator or deadline writes `abandoned` into `batch.json`.
3. The attempt is marked invalid.
4. A candidate/commit transaction requests a new write lease.
5. Lease issuance checks chat authority, barrier, incarnation, session, epoch, and cast, but has no slot-abandonment rule.
6. The new lease is issued.
7. The generic commit-gate checklist validates session, cast, base, package, and operation identity, but has no explicit same-slot terminal-state rule.
8. Slot 2 commits despite explicit abandonment.

Even if an implementation adds a coordinator check, the authority model is split: the fencing rule lives in the lease section while batch invalidation lives in the multi-character section. A freeze contract must not rely on an implementer connecting those sections correctly.

**Required correction:**

- make explicit batch-slot abandonment, slot removal, and batch invalidation fencing events or explicit lease-revocation events;
- increment the fence or revoke the affected slot-scoped lease under `chatWriteMu`;
- add plan, batch, participant index, same-slot terminal state, and abandonment checks to the universal commit-gate proof;
- define whether an abandonment operation itself creates a durable commit and how it is ordered against an in-flight candidate lease.

#### Finding C4 — Renewal safety is stated, but renewal validation is incomplete

R5 correctly provides the required safety outcome: a revoked or expired holder cannot regain commit authority by renewing or retrying. That resolves the narrow question asked in the review instruction.

The renewal procedure itself remains incomplete:

- renewal must occur under `chatWriteMu`;
- R5 does not state the full fence/session/operation validation required before updating `expiresAt`;
- revoked and expired leases are not given an explicit immutable-terminal transition table;
- renewal of a nominally active lease can therefore update lifetime state without proving that no intervening fencing event is already durable.

The outcome may remain safe because the later commit revalidates, but the lifecycle contract is not exact.

**Required correction:**

- define renewal as a validating operation, not merely a lifetime extension;
- prohibit lifetime updates to revoked, expired, orphaned, terminal, or barrier-drained leases;
- enumerate all terminal lease transitions;
- require renewal to fail closed where the authority, fence, session, batch, or barrier state cannot be proven current.

#### Finding C5 — Process incarnation and heartbeat identity remain undefined

R5 depends on `ownerProcessIncarnation` to prevent a resumed process from using an old holder token. It does not define:

- how a process incarnation is allocated;
- whether it is random, monotonic, persistent, or ephemeral;
- how collisions are prevented;
- how it survives or is regenerated across a crash;
- how startup distinguishes a delayed legitimate owner from a dead owner;
- what heartbeat record exists;
- how missed heartbeats are distinguished from scheduling delay;
- what force-aborts a stuck lease without deleting evidence.

If incarnation identity can collide or be reconstructed, the holder-token protection is weaker than claimed. If heartbeats require a lock already held by long staging work, the TTL can become a liveness trap rather than a stale-work fence.

**Concrete concern:**

1. A legitimate operation begins a large package write and retains `chatWriteMu`.
2. Staging exceeds the renewal cadence.
3. It cannot heartbeat without reentering the held lock.
4. The lease expires through no fencing change.
5. A barrier or deletion path must then wait for expiry or process termination.
6. The system is safe but needlessly blocked, or an implementer shortens safety checks to preserve liveness.

**Required correction:**

- define process-incarnation allocation, uniqueness, persistence, and reuse rules;
- define the durable heartbeat record and validation protocol;
- define missed-heartbeat and forced-abandonment behavior;
- separate safety expiration from liveness recovery without deleting unreconciled evidence.

### 2.3 A14 is not closed

R5’s treatment is careful and honest: it labels A14 partially resolved, retains `captured` as a safe interim default, leaves `prefix_committed` as an owner choice, and blocks V24 for lack of a baseline.

That is the correct status. It is not closure.

`captured` preserves the existing prompt builder only if the existing behavior is first captured. The proposal contains no pre-amendment two- or three-character baseline for:

- prompt construction;
- context-memory selection;
- cast behavior;
- lore visibility;
- response ordering;
- timing and delay;
- retry and restart behavior;
- partial-batch visibility.

There is also an evidence-linkage gap. Candidate records retain `contextHistoryHead`, but durable v2 history `turn` records retain only:

- `batchId`;
- `participantIndex`;
- `commitSequence`;
- `userCommitOperationId`;
- `characterVersion`;
- `attemptId`.

They do not retain `contextHistoryHead`, base history, plan hash, target-package hash, or generation-result linkage. An auditor must therefore leave the committed history target and consult the separate mutable generations root to determine what context a character actually saw.

That separation may be acceptable, but it weakens the claim that `captured` mode is independently auditable from the history package alone. It also makes V24 depend on both history preservation and generation-root preservation.

**Required correction:**

- retain the V24 blocked status;
- capture the missing baseline before any behavioral claim;
- decide explicitly whether persistent history must carry `contextHistoryHead` or a normative candidate-evidence pointer;
- define how concurrent `captured` generation is allowed to diverge from persisted prefix order without being misread as a context change.

A proposed ordering rule cannot substitute for missing evidence about current behavior.

### 2.4 V18–V24 remain proposed, and no Gate row has passed

R5 handles this correctly. Gate v2 still adopts only V1–V17. R5 presents V18–V24 as proposed additions, requires a new Gate record and owner action for adoption, and expressly states that no row passes.

That boundary must survive independent review. The proposal’s separate A1–A14 acceptance-criteria section should not be mistaken for adopted Gate acceptance. It is useful Gate input, but it does not change Gate v2’s adopted inventory.

---

## 3. Freeze-blocking defects

### C1 — The “all new evidence remains create-only” claim contradicts R5’s mutable state

R5 proposes AM-1 for target projection replacement while stating that commit records, receipts, tombstones, and all new evidence files remain create-only.

R5 itself requires mutable:

- `batch.json`;
- authority fence, barrier, lease-count, and recovery state;
- lease lifecycle state;
- attempt status and terminal-reason state;
- candidate status state.

This is not a drafting nicety. The owner cannot safely approve AM-1 alone while R5 relies on a broader unacknowledged class of mutable operation evidence.

A compliant implementation following the universal create-only sentence would either:

- fail to implement batch, lease, authority, attempt, or candidate lifecycle; or
- violate the sentence by updating those records.

**Required correction:** replace the universal sentence with an explicit mutable-evidence inventory and transition discipline.

### C2 — Final commit-gate validation is not bound to the linearization point

See §2.2. Validation can precede lengthy package, record, projection, and receipt work while the commit lock is held. Expiry is time-based rather than lock-based. R5 needs a final no-I/O validation immediately before atomic record installation.

### C3 — Abandoned batch slots can escape lease revocation

See §2.2. Slot abandonment must become a fencing or lease-revocation event and must appear explicitly in the universal commit proof.

### C4 — Seed presence rules contradict each other

This is a concrete deterministic-identity defect.

One rule says:

- `seed` appears only for explicit seed mode;
- `effectiveSeed` appears only for random mode when reported.

Another table says:

- `seed` is `null` in random mode;
- `effectiveSeed` is otherwise `null`.

R5 separately states that `null` is not an alias for omission. Consequently, two implementations can both claim R5 compliance while producing different canonical bytes and different generation fingerprints for the same logical generation outcome.

**Concrete counterexample:**

- Random seed mode is used.
- The engine does not report an effective seed.
- Serializer A omits both fields.
- Serializer B writes both as `null`.
- Their canonical forms differ.
- Their `generationRequestFingerprint` values differ.
- Provenance, duplicate suppression, and V19 cannot be deterministic.

**Required correction:**

- choose omission or explicit `null` for every seed mode;
- define `seedStatus` consistently with that representation;
- add golden vectors proving that explicit, random-reported, random-unreported, and non-applicable cases serialize deterministically.

### C5 — Multi-file target replacement is not atomic

AM-1 calls target-projection installation “same-volume atomic replacement.” A same-volume rename can be atomic for one file, but R5 has an eight-path target set plus optional-file deletion. There is no multi-file atomic primitive in the proposal.

The verified package and crash-recovery table substantially mitigate this, but the wording still overclaims. More importantly, R5 does not state that live readers and normal chat traffic must wait for startup crash inventory and target repair before relying on mixed current-target state.

**Concrete sequence:**

1. The commit record is installed.
2. Target projection begins replacing multiple target files.
3. The process crashes after some files are replaced and others remain old.
4. Current live targets are a mixture of two revisions.
5. A restarted reader consults current targets before package recovery completes.
6. It observes a manifest/target mismatch or attributes the wrong revision to the mixed files.

The eventual package recovery may be correct, but the intermediate mixed state must never be treated as verified current state.

**Required correction:**

- replace “atomic replacement” with a verified multi-step replacement protocol;
- require the same filesystem/volume precondition explicitly;
- prohibit normal reads, publication, migration, wipe-source selection, or further commits until startup recovery completes;
- define how optional-file deletion and partial replacement are reconciled from the package.

### C6 — Fencing-event coverage is incomplete

Character removal is now covered through paired `castEpoch` and `characterVersion` changes. Root, session, chat, cast, barrier, operation, token, and incarnation changes are also covered. Slot-level batch invalidation is missing, as shown in C3.

The same concern applies to recovery-incarnation and to any future batch-scoped authority event. The proposal should not enumerate fencing changes in one section and add independent invalidation operations elsewhere without connecting them.

**Required correction:** define one authoritative fencing-event inventory and require every invalidation operation to use it.

---

## 4. High-severity exactness gaps

### H1 — The universal commit checklist omits batch and recovery-specific state

The §3.5 checklist does not explicitly require:

- batch identity;
- participant index;
- plan hash;
- same-slot terminal state;
- slot-abandonment state;
- response-plan/user-commit linkage for candidate commits;
- recovery incarnation or recovery-lease scope validation.

Candidate, batch, and plan fields exist elsewhere, but the universal gate does not bind them. A freeze contract should state the complete proof once rather than distribute mandatory checks across unrelated sections.

### H2 — Attempt, candidate, authority, and lease transitions lack a complete state machine

R5 enumerates many statuses but does not define:

- attempt status values;
- all permitted attempt transitions;
- all permitted candidate transitions;
- all permitted lease transitions;
- all permitted authority and barrier transitions;
- terminal-reason codes;
- which operations may perform each transition;
- forbidden transitions;
- crash behavior for a partially applied status update.

This is especially important because R5 rejects a universal create-only model for this state. Mutable authority state without an exact transition table is not enforceable.

### H3 — Migration lacks an exact initialization mapping for new v2 fields

R5 recognizes v1 history as legacy and defines exact v2 schemas, but §11 does not provide a complete v1-to-v2 initialization mapping for:

- continuation sequence and preference;
- v1 history entries versus v2 `turn` records;
- character-version initialization;
- narrator/user slot conventions;
- optional world/scenario initialization;
- generation-object absence;
- unknown-field preservation;
- archive-to-package linkage.

The predecessor archive preserves old bytes, but the ordinary v2 package still needs exact construction rules from the validated v1 state. V23 cannot be meaningful without them.

### H4 — Preference writes need an explicit durable-transaction binding

R5 correctly selects `meta.json` as the sole authoritative saved preference and adds `continuationSeq` CAS. Section 10 nevertheless says that preference operations “write” metadata without expressly binding every accepted preference change to the ordinary lease, package, record, receipt, revision, and recovery protocol.

Because preference state lives inside a versioned chat target, that binding should be explicit. Otherwise, a future reader could treat preference as a lightweight metadata update outside the commit package.

### H5 — `samplerOrder` and thinking-budget normalization are not exact

R5 records `samplerOrder` as `[]` for engine order, but does not preserve the observed engine order. If engine order can affect sampling, recording an empty array loses effective control information.

Likewise, `maxTokens` is described as 512 plus any effective thinking budget without an exact budget formula. An “exact effective value” requirement needs an exact computation rule.

These are high rather than blocking only because they affect provenance completeness rather than basic durability. They still undermine the claim of complete effective-request capture.

### H6 — Archive and generation-evidence schemas need canonical byte rules

The predecessor archive defines fields and immutability but does not state the exact canonical byte order or canonicalization domain used for `archiveHash`. Its `evidence/...` directory is also undefined.

The authority, lease, barrier, attempt, candidate, plan, and batch listings likewise need explicit statements about which bytes are hashed, which bytes are audited, and which metadata fields are excluded from identity. R5 is strong on target and fingerprint domains but weaker on the new operational roots it introduces.

### H7 — A14 evidence linkage should be strengthened

Durable history should either carry `contextHistoryHead` directly or carry a normative pointer sufficient to locate the exact attempt/candidate evidence without relying on mutable generation state. The current history `turn` representation preserves ordering and attempt identity, but not the observed context head needed for the central A14 comparison.

---

## 5. Crash sequences that must be tested after correction

The following sequences are not predictions that R5 will fail them; they are the minimum adversarial tests for the corrected contract.

### Package/commit protocol

1. Crash before package installation.
   - Expected: no commit; prior state current; staging cleaned only when unreferenced.
2. Crash after package installation but before commit-record installation.
   - Expected: orphaned package; no commit; no adoption.
3. Crash after record installation but before target projections complete.
   - Expected: no mixed current-target state treated as verified; recovery installs all targets from the verified package.
4. Crash during multi-file projection.
   - Expected: restart blocks normal reads until package recovery completes.
5. Crash after targets but before receipt reconstruction.
   - Expected: receipt reconstructed only from exact record body bytes.
6. Crash after publication but before lease terminalization.
   - Expected: idempotent recovery without duplicate revisions, receipts, packages, or turns.
7. Record present with missing or corrupt package.
   - Expected: hard block; no guessed bytes.
8. Package present without a record.
   - Expected: orphaned evidence; no commit.

### Stale-writer and fencing protocol

1. Queued autosave resumes after chat switch.
   - Expected: session/epoch mismatch causes `ErrSessionInvalidated`.
2. Queued autosave resumes after root replacement.
   - Expected: incarnation mismatch blocks the write.
3. Lease expires during package staging under the commit lock.
   - Expected: final gate revalidation rejects the commit.
4. Revoked lease attempts renewal.
   - Expected: renewal fails closed; no lifetime extension.
5. Expired lease attempts renewal.
   - Expected: terminal lease remains terminal.
6. Process restarts with an active durable lease.
   - Expected: incarnation mismatch orphans the lease; no adoption of its holder token.
7. Character removed while its candidate is generated.
   - Expected: paired `castEpoch` and `characterVersion` changes reject the stale candidate.
8. Candidate transaction begins after explicit slot abandonment.
   - Expected: slot-abandonment fencing rejects lease issuance or commit.
9. Abandonment occurs after candidate lease issuance but before record installation.
   - Expected: final gate rejects the commit because the slot is no longer committable.
10. Barrier closes while a commit holds `chatWriteMu`.
    - Expected: barrier waits without forcing lock release; no new commit enters after the barrier state is durable.

### Multi-character concurrency

1. Three slots start concurrently; the middle slot never starts.
   - Expected: later slots remain blocked until explicit abandonment or another defined terminal event.
2. Slot 1 commits, slot 2 remains in flight, and the process crashes.
   - Expected: batch remains partially committed; committed turns are visible; uncommitted turns are absent; restart resumes deterministically.
3. Fast slot 3 completes before slow slot 1.
   - Expected: slot 3 does not bypass prefix order; no duplicate or reordered committed history.
4. Middle character is removed during generation.
   - Expected: slot becomes removed; attempt is invalid; stale candidate cannot commit.
5. Slot is abandoned while its candidate commit is in flight.
   - Expected: abandonment wins; no resurrection through lease replay.
6. Same display name is recreated after removal.
   - Expected: new character ID; old attempt/candidate cannot attach.
7. `prefix_committed` is selected.
   - Expected: later slots wait; delay is explicit and observable; prompts are not silently rewritten.
8. `captured` is used concurrently.
   - Expected: every candidate records the actually observed head; history order remains participant order; future baseline comparison can distinguish generation context from persisted order.

---

## 6. Prior-finding assessment

R5’s prior-finding table is directionally accurate but overstates several resolutions as final contract resolutions. Those labels are proposal-level only and must not be carried into a freeze record without the corrections above.

- A1 is substantially improved by immutable packages, but C2 and C5 remain.
- A2 is substantially improved by retaining one lock order, but final gate timing remains incomplete.
- A3 is substantially improved by expiry, heartbeat, token, incarnation, and startup rules, but C4–C6 remain.
- A4 is substantially improved by paired `castEpoch` and `characterVersion` handling.
- A5 is substantially improved by durable plans, batches, attempts, and candidates, but C1, C3, H1, and H2 remain.
- A6 is substantially improved by exact RFC 8785 use and identity/metadata separation, but C4/H5 remain.
- A7 is substantially improved by preserving the frozen operation fingerprint and adding a separate generation fingerprint.
- A8 is substantially improved by one authoritative preference and CAS, but H4 remains.
- A9 is substantially improved by one package/record protocol and archive binding, but H3 and H6 remain.
- A10 is substantially improved by a broader wipe inventory, but mutable operation evidence needs the same inventory discipline.
- A11 is substantially improved by complete target/record/package schemas, but H3 and H6 remain.
- A12 is substantially improved by model/execution separation and parameter defaults, but H5 remains.
- A13 is correctly handled: V18–V24 remain proposed.
- A14 remains correctly unresolved pending baseline and stronger history linkage.

No prior finding is verified. Verification still requires Gate evidence, native evidence, fault-injection evidence, and the staged Mimo challenge.

---

## 7. Required owner decisions

This review makes no owner decision. It preserves the following boundaries:

- **OD-02/AM-1:** awaiting owner; must be split into narrow target-projection replacement and separate layout/mutable-evidence authorization.
- **OD-14:** V18–V24 remain proposed; adoption requires a separate owner decision and new Gate binding.
- **OD-15:** `captured` remains the interim default; `prefix_committed` remains an owner choice; neither is validated without V24 baseline evidence.
- **OD-01, OD-03 through OD-13, and OD-16:** remain proposed unless separately accepted.
- **F3-C phase 0:** remains open.
- **P1a:** remains unauthorized.
- **F3-C acceptance:** not granted.
- **F3-D:** remains blocked.
- **A9/A10/A12/A13:** remain closed and are not reopened by this review.
- **V20:** remains blocked pending F3-F evidence.
- **Gate v2:** remains adopted only for V1–V17.

---

## 8. Integrity and limits

Reviewed-file integrity was established by direct local SHA-256 calculation:

- R5 hash: `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`.
- Frozen-contract hash: `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`.
- Gate v2 hash: `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`.
- R4 proposal hash: `7a2c221147621706e299924c20de1fe29080b7d649fac25264ce212d11e7417a`.

This review checked proposal-to-proposal consistency, proposal-to-frozen-authority compatibility, prior-review closure claims, and adversarial crash/stale/concurrency sequences. It did not execute or authorize:

- source changes;
- tests;
- builds;
- migrations;
- model loading;
- destructive testing;
- native Wails/WebView2 runs;
- application-data access;
- authority-record creation;
- freeze;
- acceptance;
- P1a implementation;
- F3-D work.

No existing review, proposal, contract, assessment, decision-register, source, or test file was modified. This is a new independent-review artifact only.
