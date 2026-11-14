# DISPOSITION MATRIX — Independent R5 Review (Muse)

**Status:** **ADJUDICATION ONLY — NOT R6, NOT A CONTRACT REVISION, NOT FREEZE, NOT ACCEPTANCE, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Reviewed proposal:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md`
**Proposal SHA-256:** `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`
**Independent review:** `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R5-2026-11-11.md`
**Review SHA-256:** `a8441f9dc9b30eddfc9493f11d4dca2d0252da3f89bab735bedef957d4499c2d`
**Review size:** 560 lines; 32,238 bytes.
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` (`61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`)
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md` (`7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`)

**Scope:** documentation-only adjudication. This file does not amend R5, create R6, adopt V18–V24, freeze anything, or authorize implementation. No source, tests, builds, migrations, model loads, destructive tests, application data, authority records, P1a, or F3-D work was performed.

---

## 1. Adjudication rule

Muse's findings are treated as review claims to be dispositioned, not as instructions to implement. Each finding is assigned exactly one of:

- **Accept** — the finding is correct and a future R6 must incorporate the correction.
- **Accept; split** — the finding is correct, but the correction requires two separate owner decisions.
- **Keep open** — the finding cannot be closed without missing baseline evidence or an owner decision.
- **Reject** — not used in this matrix.

A disposition is not a freeze resolution. R5 remains proposal-only, and its own A1–A14 labels remain proposal-level.

**Overall adjudication:** accept Muse's no-freeze recommendation. The six blocking findings (C1–C6) are valid. The seven exactness findings (H1–H7) are valid as future-R6 requirements. A14 remains open. No finding is rejected.

---

## 2. Blocking findings

| ID | Finding | Disposition | Required change in a future R6 | Owner decision and frozen impact |
|---|---|---|---|---|
| **C1** | R5's universal “all new evidence files remain create-only” claim contradicts its mutable batch, authority, lease, attempt, and candidate records, and its new roots are not in the frozen layout. | **Accept; split** | Replace the universal sentence with an explicit inventory of immutable versus mutable records. Split OD-02 into OD-02a (target replacement only) and OD-02b (new roots plus mutable operation evidence). Enumerate permitted states, transitions, write protocol, recovery, inventory, and wipe treatment for every mutable record. | **OD-02a + OD-02b.** Master §4.1 layout and §6.9 create-only installation. AM-1 must not absorb OD-02b. |
| **C2** | Lease expiry can lapse after validation and before the record installation that is the commit point. | **Accept** | Use a monotonic clock for lease lifetime, with persisted wall-clock timestamps only for audit. Repeat the complete commit proof under `chatWriteMu` immediately before the create-only record installation. Permit no I/O between that final proof and the record install. Treat target projection and receipt installation after the record as post-commit apply, not authorization. Bound the final critical section. | **OD-04** for timings; no new lock and no change to the frozen `chatWriteMu` order. |
| **C3** | Slot abandonment is not a fencing event and does not automatically revoke a later slot-scoped lease. | **Accept** | Add batch, plan, participant index, slot terminal state, and abandonment to the universal commit proof. Make explicit slot abandonment, slot removal, and batch invalidation fencing events. Revoke affected slot-scoped leases or increment `fenceEpoch` under `chatWriteMu`. Define whether abandonment is itself a durable commit and order it against an in-flight candidate. | **OD-06** and **OD-12**; no frozen-lock change. V22 and V24 remain proposed. |
| **C4** | Seed presence rules conflict: one rule requires omission for unavailable seeds, while the defaults table writes `null`; therefore canonical bytes and fingerprints can diverge. | **Accept** | Define one canonical representation. Recommended: unavailable values are **omitted**, never `null`; `seed` is present only for fixed mode; `effectiveSeed` is present only for random mode with a reported value; `seedStatus` has one closed mapping consistent with those cases. Add golden vectors for fixed, random-reported, random-unreported, and not-applicable. Normalize `explicit` versus `fixed` terminology. | **OD-12**; V19 remains proposed. No frozen target change outside the v2 schema. |
| **C5** | “Same-volume atomic replacement” is true for one file, not for an eight-path target set plus optional-file deletion; recovery does not yet prohibit reads of mixed current targets. | **Accept** | Delete the multi-file atomicity claim. Define a verified multi-step projection protocol with an explicit same-filesystem precondition. Normal reads, publication, migration, wipe-source selection, and further commits must wait for startup repair and manifest verification. Optional-file deletion must be package-derived and idempotent. | **OD-02a** must authorize only per-file replacement from a verified package, not a multi-file atomic primitive. **OD-02b** must authorize any durable projection-status record. Master §6.9. |
| **C6** | Fencing-event coverage is distributed across sections and omits slot-level batch invalidation and recovery-specific authority. | **Accept** | Define one authoritative fencing-event inventory. Include root identity and incarnation, session, epoch, cast, character version, fence epoch, operation/token, process incarnation, barrier/tombstone, batch/slot abandonment or removal, plan invalidation, and recovery incarnation. Every invalidation operation must use that inventory. | **OD-05**, **OD-06**, and **OD-12**; no frozen-lock change. |

---

## 3. Exactness findings

| ID | Finding | Disposition | Required change in a future R6 | Owner decision and Gate impact |
|---|---|---|---|---|
| **H1** | The universal commit checklist omits batch identity, participant index, same-slot terminal state, abandonment, response-plan linkage, and recovery scope. | **Accept** | State the complete proof once in the commit-gate section: operation, lease/token/incarnation, session/epoch/cast/fence, plan hash, batch and slot, attempt/candidate, base/history, package/manifest/record, abandonment and terminal state, recovery scope, and prior commit linkage. | **OD-12**; V22 remains proposed. |
| **H2** | Attempt, candidate, lease, authority, and barrier transitions lack complete state machines. | **Accept** | Define every status value, every legal transition, every forbidden transition, terminal-reason codes, the operation permitted to perform each transition, and crash behavior for a partially applied update. No prose-only state. | **OD-12**; V22 remains proposed. |
| **H3** | Migration lacks a complete v1-to-v2 initialization mapping. | **Accept** | Specify initialization for continuation sequence and preference, v1 history entries to v2 `turn` records, character versions, narrator/user slots, optional world/scenario, absent generation objects, unknown-field preservation, and archive-to-package linkage. | **OD-10** and **OD-12**; V23 remains proposed. |
| **H4** | Preference writes are not expressly bound to the ordinary durable transaction protocol. | **Accept** | Every accepted preference change uses a `preference`-scope lease, package, record, receipt, revision, and recovery path. No direct `meta` write outside the commit protocol. CAS failure or commit failure leaves the previous saved value. | **OD-09** and **OD-12**; V21 remains proposed. |
| **H5** | `samplerOrder = []` does not preserve observed engine order, and “512 plus any effective thinking budget” has no exact formula. | **Accept** | Narrow the claim: record whether sampler order is explicit, engine-default, or engine-opaque. Do not call `[]` the effective engine order. If the order is not observable, record that fact and narrow the “complete effective request” claim. Define the exact thinking-budget formula, or record it as unknown and narrow the claim. | **OD-12** and **OD-13**; V18 and V19 remain proposed. |
| **H6** | Archive and generation-evidence schemas lack exact canonical byte rules, hash domains, and evidence-directory definitions. | **Accept** | Define canonical bytes and hash domains for `archiveHash`, authority, barrier, lease, plan, batch, attempt, and candidate. State exactly which fields are included and excluded from identity versus audit. Define the `evidence/` relative-path manifest, sizes, hashes, and mtime exclusion. | **OD-12**; V22 and V23 remain proposed. |
| **H7** | V2 history does not retain the observed context head or a normative candidate-evidence pointer, weakening A14 auditability. | **Accept** | Add `contextHistoryHead` directly to the durable v2 `turn`, or define one normative pointer to immutable candidate evidence. The history package must not depend on mutable generation state to answer what context a character observed. | **OD-15**; V24 remains blocked. This changes the v2 target schema, not v1. |

---

## 4. A14 disposition

**A14 remains open. It is not closed by `captured`.**

`contextMode = captured` is accepted as the safe interim behavior because it avoids changing the existing prompt builder. It is not evidence that persisted history preserves the context head actually observed.

Required before A14 can close:

1. Capture a pre-amendment two- or three-character baseline covering prompt construction, context-memory selection, cast behavior, lore visibility, response ordering, timing, retry/restart behavior, and partial-batch visibility.
2. Decide whether `captured` is the owner-approved default or whether `prefix_committed` is required.
3. Preserve the observed context head in durable history or through a normative immutable candidate pointer.
4. Define how concurrent `captured` generation may diverge from persisted prefix order without being mistaken for a context change.
5. Run V24. Until then, V24 remains **BLOCKED** and A14 remains **Partially resolved** at best.

---

## 5. AM-1 split

R5's present OD-02 is too broad. It must be split into two decisions. The existing OD-02 is not silently redefined; the split is explicit.

### 5.1 OD-02a — target projection replacement only

OD-02a may authorize only:

- byte-verified replacement of the current eight target paths from an already durable, verified commit package;
- per-file same-volume replacement, with the prior target bytes retained in that package;
- the commit record as the single commit point;
- target-manifest verification, projection-status tracking, read blocking until repair, and recovery from the package;
- optional-file deletion derived from the package manifest, performed idempotently.

OD-02a must not authorize:

- commit records, receipts, tombstones, packages, or predecessor archives to be replaced;
- plans, attempts, candidates, batches, authority records, leases, or barriers to be replaced;
- a multi-file “atomic” primitive;
- any new root or mutable operation-evidence class.

### 5.2 OD-02b — new roots and mutable operation evidence

OD-02b is a separate owner decision. It must authorize and enumerate:

- `f3c-packages/` as immutable, create-only package storage;
- `f3c-authority/` including authority, lease, and barrier records;
- `f3c-generations/` including immutable plan headers and mutable batch, attempt, and candidate state;
- `f3c-predecessors/` as immutable archive storage;
- any projection-status or repair marker introduced by C5.

OD-02b must state, for every mutable file:

- its path domain and identity fields;
- its permitted statuses and transitions;
- its write protocol and lock;
- its recovery and crash behavior;
- its inventory and wipe treatment;
- which fields are hashed, which are excluded from identity, and which are audit-only.

If frozen Master §4.1 is treated as an exhaustive layout, OD-02b requires its own formal amendment. It must not be folded into AM-1.

### 5.3 AM-1 prohibitions

AM-1 must not authorize replacement of commit records, receipts, tombstones, packages, predecessor archives, plans, attempts, candidates, batches, authority records, leases, or barriers. It must not approve “same-volume atomic replacement” as a multi-file primitive.

---

## 6. Final commit-boundary rule

This is a required correction for a future R6, not R6 text.

### 6.1 Clock

- Lease lifetime uses a monotonic clock while the process is alive.
- Persisted `issuedAt`, `lastHeartbeatAt`, `expiresAt`, and `maxExpiresAt` are wall-clock audit values, not the authority comparison.
- The in-memory lease carries `issuedMono`, `deadlineMono`, and `maxDeadlineMono`.
- Expiry is uniform: `nowMono >= deadlineMono` or `nowMono > maxDeadlineMono` fails.
- On restart, the previous process incarnation is not reusable. Existing leases are `orphaned`; no persisted monotonic value is trusted across a process boundary.

### 6.2 Final gate

Under `chatWriteMu`, immediately before the create-only commit-record installation, validate all of:

- lease status is `active`; holder token and process incarnation match;
- current monotonic time is before both deadlines;
- root identity and root incarnation match;
- chat is non-tombstoned and the barrier permits commit;
- `stateSessionId`, `stateEpoch`, `castEpoch`, character version, and `fenceEpoch` match;
- operation identity and kind match;
- base revision and base history head match;
- no conflicting commit exists for the operation;
- plan hash, batch ID, participant index, slot state, and abandonment state match;
- candidate and attempt identities match the slot;
- package hash, target manifest, record body bytes, and record hash match;
- recovery lease scope is valid where recovery is the operation.

A mismatch fails with `ErrMutationFenceStale` or `ErrMutationLeaseInvalid`; no record is installed.

### 6.3 Linearization and post-commit work

- The atomic create-only installation of the commit record is the single commit point.
- No I/O may occur between the final gate and that installation.
- Package construction, flushing, read-back, and hash verification occur before the final gate.
- Target projection and receipt installation occur after the record installation. They are post-commit apply and repair work, not authorization. They may not be used to claim that an expired lease authorized the commit.
- If the lease expires during package construction, the final gate rejects the commit; the package is orphaned evidence and is not adopted.
- The final critical section is bounded by operation class and payload size. The contract must define the bound and the staging cleanup rule.

### 6.4 Heartbeat and process incarnation

- `processIncarnation` is a fresh 128-bit random value allocated at process start, recorded under `chatWriteMu`, and never reused by a later process.
- A heartbeat is a durable validating update, not a blind lifetime extension. It must recheck lease, holder, incarnation, session, epoch, cast, fence, and barrier state.
- Revoked, expired, orphaned, terminal, or barrier-drained leases cannot be renewed.
- A heartbeat must not be required while the holder is inside the final commit section. If staging would exceed the heartbeat cadence, the design must stage outside the final section or fail closed; it must not shorten the authority check to preserve liveness.
- Missed heartbeats may make a lease recoverable as `orphaned` or `expired`; they never authorize deletion of unreconciled evidence.

---

## 7. Mutable evidence inventory required by OD-02b

| Path or class | Classification | Required rule |
|---|---|---|
| `chats/<chatId>/` targets | Replaced only under OD-02a | Per-file replacement from a verified package; no multi-file atomic claim. |
| `f3c-commits/` records | Immutable | Create-only; never replaced. |
| `f3c-receipts/` | Immutable | Create-only; reconstructed only from record body bytes. |
| `f3c-tombstones/` | Immutable | Create-only; never replaced. |
| `f3c-packages/` | Immutable | Create-only package; exact target bytes; retained until deletion or wipe. |
| `f3c-predecessors/` | Immutable | Create-only archive; exact v1 evidence. |
| `f3c-authority/.../authority.json` | Mutable | Explicit transition table; `fenceEpoch`; commit sequence; lease count; recovery incarnation. |
| `f3c-authority/.../leases/` | Mutable lifecycle | Explicit issued/active/committing/committed/aborted/expired/orphaned/revoked transitions. |
| `f3c-authority/.../barriers/` | Mutable lifecycle | Explicit open/closing/sealed/blocked-timeout/completed transitions. |
| `f3c-generations/.../plans/` | Immutable | Create-only plan header. |
| `f3c-generations/.../batches/` | Mutable | Batch and slot state; abandonment and terminal transitions. |
| `f3c-generations/.../attempts/` | Mutable | Attempt lifecycle and terminal reasons. |
| `f3c-generations/.../candidates/` | Mutable | Candidate lifecycle and terminal reasons. |
| Projection status or repair marker | Mutable, if adopted | Must be authorized by OD-02b; blocks reads until repair completes. |
| `f3c-tmp/` | Ephemeral | No user evidence may remain after cleanup; never visible to inventory readers. |

R5's current sentence that “all new evidence files remain create-only” must be removed or replaced by this classification. It cannot remain as a universal claim.

---

## 8. Frozen and Gate conflict check

| Question | Result |
|---|---|
| Does C2 require a new lock or a different lock order? | **No.** The final gate stays under `chatWriteMu`; record installation remains the commit point. |
| Does C3 require a lock outside the frozen order? | **No.** Abandonment and lease revocation remain under `chatWriteMu`, with root journaling only where already required. |
| Does C5 permit a multi-file atomic claim? | **No.** OD-02a must authorize verified per-file replacement and a read/recovery barrier, not a multi-file atomic primitive. |
| Does C1 require frozen Master §4.1 to change? | **Yes, unless the owner treats §4.1 as non-exhaustive.** The new roots and mutable operation evidence need OD-02b, and likely a separate layout amendment. |
| Does H7 change frozen v1? | **No.** It adds a v2 history field or pointer; v1 targets and receipts remain unrehashed and valid. |
| Do any findings adopt V18–V24? | **No.** V18–V24 remain proposed. No Gate row passes. |
| Do any findings reverse A9/A10/A12/A13 closures? | **No.** H3, H5, and H6 are new exactness requirements. They do not reopen the prior findings by implication. |
| Does this matrix authorize implementation? | **No.** Phase 0 remains open; P1a remains unauthorized; F3-C acceptance is not granted; F3-D remains blocked. |

---

## 9. Impact on R5's prior dispositions

R5's prior-finding table is directionally sound but overstates several entries as final contract resolutions. The following are **proposal-level** impacts only.

| R5 finding | Impact of Muse review |
|---|---|
| A1 | Package mechanism retained; C2 and C5 remain. Not freeze-ready. |
| A2 | Lock architecture retained; final commit-gate timing remains incomplete. |
| A3 | Lease shape retained; renewal, incarnation, heartbeat, and expiry-boundary rules remain incomplete. |
| A4 | Session/cast/character fencing retained; batch-slot abandonment remains a gap. |
| A5 | Plan/batch/attempt/candidate model retained; C1, C3, H1, and H2 remain. |
| A6 | Canonicalization split retained; C4 and H5 remain. |
| A7 | No new defect; operation and generation fingerprints remain separate. |
| A8 | One-preference rule retained; H4 durable-transaction binding remains. |
| A9 | Package/archive order retained; H3 and H6 remain. |
| A10 | Inventory retained; mutable evidence classification must be added under OD-02b. |
| A11 | Schema work retained; H3 and H6 remain. |
| A12 | Model/execution split retained; H5 remains. |
| A13 | No change; V18–V24 remain proposed. |
| A14 | H7 adds the missing evidence link; A14 remains open and V24 remains blocked. |

No prior finding is verified. No R5 label is converted into a freeze claim by this matrix.

---

## 10. R6 authorization boundary

**No R6 is authorized or written by this matrix.**

Before an R6 is commissioned, the owner should:

1. Accept or modify this disposition matrix.
2. Decide OD-02a and OD-02b separately.
3. Decide whether OD-02b requires its own formal amendment to frozen §4.1.
4. Accept the final commit-boundary rule in §6, or direct a different safe rule.
5. Confirm that A14 remains open and V24 remains blocked.

An R6, if later authorized, must preserve R5's sound architecture and address every row above. It must not treat this matrix as an R6 draft, a freeze record, or a Gate adoption.

---

## 11. Status

**Matrix:** **ADJUDICATION ONLY; NOT A CONTRACT REVISION.**
**R5:** **PROPOSAL ONLY; NOT FROZEN; NOT ACCEPTED; NOT IMPLEMENTED; NOT VERIFIED.**
**R6:** **NOT AUTHORIZED; NOT WRITTEN.**
**F3-C phase 0:** **NOT CLOSED.**
**P1a:** **UNAUTHORIZED**, downstream of P1c.
**F3-C acceptance:** **NOT GRANTED.**
**F3-D:** **BLOCKED.**
**A9/A10/A12/A13:** closed; not reopened.
**V20:** **BLOCKED** pending F3-F evidence.
**Gate v2:** adopted for V1–V17 only; V18–V24 remain proposed.

No source, test, build, migration, model load, destructive test, application data, authority record, freeze, or implementation action was performed.

---

## Appendix A — Identity of adjudicated artifacts

| Artifact | SHA-256 |
|---|---|
| R5 proposal | `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e` |
| Independent R5 review | `a8441f9dc9b30eddfc9493f11d4dca2d0252da3f89bab735bedef957d4499c2d` |
| Frozen Master | `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31` |
| Gate v2 | `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a` |

The matrix's own SHA-256 is recorded in the completion report rather than embedded here, to avoid a self-referential file hash.
