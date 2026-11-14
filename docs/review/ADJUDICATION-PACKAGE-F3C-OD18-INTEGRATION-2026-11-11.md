# ADJUDICATION PACKAGE — F3C OD-18, State Machine, Lock Proof, and CR Ownership

**Status:** **READ-ONLY ADJUDICATION — NOT R6, NOT A CONTRACT REVISION, NOT FREEZE, NOT ACCEPTANCE, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Integration review:** `docs/review/INTEGRATION-REVIEW-F3C-TARGET-2-R5-2026-11-11.md` (`b29bb562c2a5b596baa3bb5689d6aaaf0adf9c580ac18ab72a99eae8c13f03f6`)
**Independent integration review:** `docs/review/REVIEW-INDEPENDENT-INTEGRATION-REVIEW-F3C-TARGET-2-R5-2026-11-11.md` (`636a5f8d91ba569fd0c5fed45bd00621a1e2d0aafc6fbb4bfbb902e93ff9374a`)
**R5 proposal:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md` (`bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`)
**Original matrix:** `docs/review/DISPOSITION-MATRIX-F3C-TARGET-2-R5-MUSE-2026-11-11.md` (`061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6`)
**Precision addendum 1:** `docs/review/ADDENDUM-PRECISION-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md` (`d059054b678f95d8697b8987bf856035a20a6415bd61828a940ebbd99d7549b0`)
**Precision addendum 2:** `docs/review/ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md` (`be399e5bfda3fec62e4356ad8e0dadde11a643881bd74fc314ba2bf9b116e403`)
**Frozen Master:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` (`61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`)
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md` (`7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`)

**Scope:** documentation and read-only source inspection. No source, tests, builds, migrations, model loads, destructive tests, application data, authority records, freeze, acceptance, or implementation changes were made. No R6 was drafted.

**Limits:** source inspection was limited to root/lock surfaces relevant to OD-18. This is not a complete source audit and not a Windows runtime proof. The integration review's own source-level entry-point limit remains.

---

## 1. Verdict

| Question | Answer |
|---|---|
| Can OD-18 be enforced on the currently supported Windows configurations? | **No. Supported configurations are undefined, and no cross-process authority primitive exists or has been proved.** |
| Is the state machine total and deterministic? | **No. It has the right states but not a total event mapping or defined exits for every blocked/destructive state.** |
| Is the lock proof complete across every inventoried entry point and failure path? | **No. It is conditional and leaves startup, multi-chat, wipe, gate-release, verification, and OD-18 placement gaps.** |
| Does every CR entry have an explicit owner, disposition, and required authority? | **No in the current artifacts. This package assigns proposed owners, but no owner has accepted them and no CR is closed.** |
| Which requirements remain mutually incompatible or unproven? | **Several frozen-authority conflicts and unproven Windows/runtime properties remain blockers.** |

**Overall:** no R6. The unresolved items are blockers, not invitations to invent a compromise.

---

## 2. Q1 — Can OD-18 be enforced on supported Windows configurations?

### 2.1 Answer

**No, not as currently evidenced.**

The integration review recommends one authority process per root plus a storage-level exclusive lock, with split views outside the supported deployment. That is a coherent direction, but it is not an enforceable claim yet:

- no supported Windows storage configuration list exists;
- no cross-process authority primitive is implemented for the F3-C chat store;
- no Windows runtime proof exists for lock acquisition, stale handles, deletion, network shares, removable media, cloud sync, snapshots, or duplicated environments;
- the current source has separate root resolvers, so the OD-18 scope itself is not settled.

### 2.2 Source and root evidence

Read-only source observations:

- `internal/paths/paths.go` resolves the project root and `paths.StateRoot()` defaults to the installation/project root when isolation is off.
- `internal/backend/chat/persistence.go` resolves a legacy chat root at `%USERPROFILE%\AppData\Roaming\HOWLChat` when isolation is off.
- `internal/backend/chat/service.go` uses `paths.StateRoot()` for chat/service state.
- No `LockFileEx`, `CreateFile` sharing-mode authority, `flock`, or equivalent F3-C chat-store lock was found in the inspected source.
- The only inspected `O_EXCL` uses are unrelated audio/content-import paths, not a chat-store authority lock.

This means OD-18 has two unresolved questions before any Windows proof:

1. **Which root is the authority root?** `paths.StateRoot()` or the legacy `AppData\Roaming\HOWLChat` root?
2. **Must one lock cover both roots,** or is the legacy root excluded from F3-C until migration?

Until those are answered, a lock cannot be scoped to "the storage root."

### 2.3 Candidate supported-set decision

**Recommended supported set, if the owner chooses enforcement rather than an explicit bounded assumption:**

- **Local fixed NTFS volume only.**
- One authority process per root.
- An OS-enforced exclusive handle or lock held for the process lifetime, acquired before any F3-C storage access.
- Fail closed if the primitive is unavailable, denied, timed out, or cannot be proved to exclude a second writer.

**Unsupported until separately proved and approved:**

- SMB/network shares;
- ReFS;
- removable USB or other removable media;
- OneDrive, Dropbox, or other cloud-synced folders;
- redirected/roaming profiles that can sync;
- containers, virtual machines, snapshot restores, and duplicated volumes;
- any split mount or second path to the same storage that can bypass the lock;
- any configuration where the app cannot determine the filesystem/share type.

An unlisted configuration must not be treated as locked. It must fail closed or be explicitly declared unsupported with owner-accepted residual risk.

### 2.4 Required fail-closed behavior

If OD-18 is adopted, every F3-C storage access denied by lock failure must:

- block ordinary reads and writes;
- block startup recovery and repair;
- distinguish "held by a live authority" from "lock subsystem broken";
- never proceed after a timeout or retry budget;
- never delete a lock file or steal a stale handle by PID, timestamp, or content inference;
- release only through OS process termination or a proven authority handoff.

### 2.5 Verdict

**OD-18 remains unresolved.** The local-NTFS direction is feasible in principle but unproved. SMB, removable, synced, snapshot, container, and duplicated configurations cannot be claimed safe without per-configuration evidence. **No cross-process single-writer claim is allowed.**

---

## 3. Q2 — Is the state machine total and deterministic?

### 3.1 Answer

**No.**

The integration state machine correctly separates `S0_STABLE`, `S1_ORPHAN_PACKAGE`, `S2_RECORD_INDETERMINATE`, `S3_REPAIR_REQUIRED`, `S4_REPAIRING`, `S5_VERIFIED_PENDING_RELEASE`, `S6_RELEASED`, `S7_BLOCKED`, `S8_ORPHAN_MARKER`, `D1_DELETING`, and `W1_WIPING`. It also correctly requires no ordinary access before `S6`.

It is not total. The following event classes are missing or underdetermined:

| State | Missing / underdetermined event |
|---|---|
| `S0` | New wipe intent, tombstone arrival, lock failure, storage disappearance, unexpected package/record appearance. |
| `S1` | Record appears, torn record appears, duplicate package, marker appears, retention expiry, lock failure. |
| `S2` | Package disappears, marker appears, record torns after observation, duplicate/conflicting record, lock loss. |
| `S3` | Package disappears, marker loss, repair attempt failure, wipe intent, tombstone. |
| `S4` | Marker loss, package unreadable mid-repair, lock loss, power loss, wipe intent, tombstone. |
| `S5` | Successor commit begins, power loss after verified state but before release, reader admission, release write failure, lock loss. |
| `S6` | Successor projection starts, tombstone/delete, wipe intent, lock loss, out-of-band target corruption. |
| `S7` | No defined legal exit; operator disposition, wipe, and destructive repair are not ordered. |
| `S8` | Record appears late, marker cleanup, package orphan, retention expiry. |
| `D1` | Crash mid-delete, tombstone loss, directory removal ambiguity, completion/failure exits. |
| `W1` | Crash mid-wipe, wipe intent loss, partial new-root state, recovery exclusion, completion/failure exits. |

The critical `S5 → S6` release step is not yet proven atomic with ordinary reader admission and successor projection. `S6 → S2` is stated, but the exact instant ordinary reads stop and the old release becomes unobservable is not defined.

**Verdict:** the state machine is a required shape, not a total deterministic contract. R6 must publish the full `(state, event) → outcome` table, including blocked/deleting/wiping exits and the release race.

---

## 4. Q3 — Is the lock proof complete across every inventoried entry point and failure path?

### 4.1 Answer

**No.**

The integration review correctly preserves the frozen lock order and marks the proof conditional. It does not prove the whole lifecycle. The following are still open:

1. **Startup recovery sequence:** exact acquisition/release order, root-journal placement, and prohibition on publishing before release.
2. **Multi-chat recovery:** stable `(rootId, chatId)` order is stated but not proved across repair, marker writes, retries, and error paths.
3. **Wipe/recovery exclusion across restarts:** tree-Lock vs tree-R is not enough if crash recovery can begin old-root repair while wipe recovery is pending.
4. **Gate release:** exact reader lock sequence and serialization against successor projection are not defined.
5. **Verification purity:** no upward lock acquisition under `chatWriteMu` remains a rule, not a structural proof.
6. **Operation-map integration:** registration, duplicate waiters, cancellation, historical results, and wipe invalidation are not bound to lock ownership.
7. **OD-18 placement:** the storage-level lock is outside the frozen order, but its acquisition/release relationship to F3-C locks is not proved.
8. **Publication:** durable state before publication is stated, but `chatStateMu` leaf behavior and release publication are not fully sequenced.

The contract-level entry-point inventory is broad but source-level completeness is unverified. A lock proof cannot be called complete while the set of callers is not closed.

**Verdict:** the happy-path order is preserved; the whole-lifecycle proof is incomplete. **No lock-inversion-free claim is allowed yet.**

---

## 5. Q4 — Does every CR entry have an explicit owner, disposition, and required authority?

### 5.1 Answer

**No in the current artifacts.** The integration review leaves several entries as open, compatible, or scope mismatch without an accepted owner decision. This package proposes owners; none is accepted by the owner.

| CR | Proposed owner | Disposition | Required authority | Status |
|---|---|---|---|---|
| CR-01 | Owner / OD-17 | Open | Formal layout amendment if §4.1 exhaustive; OD-02b enumeration | **Open** |
| CR-02 | AM-1 / OD-02a | Open | Formal amendment or explicit compatibility ruling | **Open** |
| CR-03 | OD-02a + OD-02b | Open | Per-file protocol; no multi-file atomic claim | **Open** |
| CR-04 | OD-18 | Open | Storage-level authority decision | **Open / critical** |
| CR-05 | OD-19 (proposed) + OD-12 | Open | Formal v2 target-schema compatibility/version ruling | **Open** |
| CR-06 | OD-19 (proposed) + OD-12 | Open | Formal v2 record compatibility/version ruling | **Open** |
| CR-07 | OD-01 + OD-02b | Open | Reference-aware retention and cleanup policy | **Open** |
| CR-08 | OD-14 | Compatible — hold | Gate record and owner adoption | **Hold** |
| CR-09 | OD-15 | Compatible — open | V24 baseline and context-mode decision | **Open** |
| CR-10 | None — prohibition | Scope mismatch | Do not cite worldstore audit as chat-store proof | **Closed as prohibition only** |
| CR-11 | OD-02b | Open | Wipe inventory for every new path | **Open** |
| CR-12 | OD-02b | Open | Retention linkage; forbid prefix/age deletion | **Open** |

**OD-19 is a proposed tracking ID** for the v2 compatibility ruling; final numbering remains the owner's.

No CR is closed by this package. CR-10 is closed only as a scope prohibition, not as an implementation or evidence claim.

---

## 6. Q5 — Which requirements remain mutually incompatible or unproven?

### 6.1 Incompatible under current frozen authority

These are not proven impossible; they cannot both be true without an owner decision or formal amendment.

1. **Frozen §4.1 exhaustive layout vs new roots** — CR-01.
2. **Frozen §6.9 create-only installation vs repeated target replacement** — CR-02.
3. **Frozen “same-volume atomic move” vs eight-target replacement plus optional deletion** — CR-03.
4. **Frozen in-process locks vs cross-process storage authority** — CR-04.
5. **Frozen target/record field sets vs F3C-TARGET-2 additions** — CR-05/CR-06.
6. **Retained evidence while a chain depends on it vs bounded cleanup/wipe** — CR-07/CR-11/CR-12.

### 6.2 Unproven, not necessarily incompatible

- OD-18 enforcement on local NTFS.
- OD-18 behavior on SMB, ReFS, removable, OneDrive/synced, container, snapshot, or duplicated volumes.
- stale-handle and lock-file deletion/recreation races;
- state-machine totality and all blocked/deleting/wiping exits;
- lock proof across startup, multi-chat recovery, wipe restart, gate release, verification, and operation map;
- source-level completeness of the storage-entry inventory;
- reference-aware cleanup and bounded A14 evidence retention;
- power-loss durability and package/record/directory synchronization;
- A14 baseline, OD-15 decision, V24;
- F3-F ownership and V20;
- Gate V18–V24 adoption and any PASS claim.

### 6.3 Verdict

The incompatible items require named authority. The unproven items require evidence. Neither may be resolved by R6 prose alone or by an optimistic implementation.

---

## 7. Blocking decisions and dispositions

| ID | Required decision / evidence | State |
|---|---|---|
| OD-17 | Frozen §4.1 exhaustiveness and other layout-adjacent authority | **Awaiting owner** |
| OD-02a | Per-file target replacement mechanics; inert without OD-02b | **Awaiting owner** |
| OD-02b | Durable recovery/layout/mutable-evidence authorization | **Awaiting owner** |
| OD-18 | Storage-level single-writer authority and supported Windows set | **Awaiting owner; critical** |
| OD-19 (proposed) | V2 target/record compatibility ruling | **Proposed; no owner acceptance** |
| OD-14 | V18–V24 Gate adoption | **Awaiting owner** |
| OD-15 | Multi-character context mode | **Open** |
| A14 / V24 | Baseline and parity evidence | **Open / blocked** |
| State machine | Total transition table and release proof | **Missing** |
| Lock proof | Whole-lifecycle acquisition and exclusion proof | **Missing** |
| Entry-point inventory | Source-level caller binding | **Unverified** |
| Retention lifecycle | Reference-aware cleanup and bounded A14 retention | **Missing** |

---

## 8. Final status

**This package:** **READ-ONLY ADJUDICATION; NOT R6; NOT A CONTRACT REVISION.**
**R5:** **PROPOSAL ONLY; NOT FROZEN; NOT ACCEPTED; NOT IMPLEMENTED; NOT VERIFIED.**
**Original matrix:** preserved.
**Precision addenda 1 and 2:** preserved.
**Integration review:** preserved as composition/proof-obligation authority, not proof.
**Independent integration review:** accepted as findings; not approval.
**R6:** **NOT AUTHORIZED; NOT WRITTEN.**
**F3-C phase 0:** **NOT CLOSED.**
**P1a:** **UNAUTHORIZED**, downstream of P1c.
**F3-C acceptance:** **NOT GRANTED.**
**F3-D:** **BLOCKED.**
**A9/A10/A12/A13:** closed; not reopened.
**V20:** **BLOCKED** pending F3-F evidence.
**Gate v2:** adopted for V1–V17 only; V18–V24 remain proposed.
**A14/V24:** open / blocked.

**No R6 until:**

1. OD-17, OD-02a, OD-02b, and OD-18 are separately decided.
2. Supported Windows configurations and fail-closed behavior are fixed.
3. The state machine is total and deterministic.
4. The lock proof covers startup, recovery, ordinary reads, multi-chat, wipe, gate release, verification, and operation-map paths.
5. Every CR has an accepted owner, disposition, and authority path.
6. The source-level entry-point inventory and retention lifecycle are complete.

No source, tests, builds, migrations, model loads, destructive tests, application data, authority records, freeze, acceptance, or implementation action was performed.

---

## Appendix A — Artifact identity

| Artifact | SHA-256 |
|---|---|
| Integration review | `b29bb562c2a5b596baa3bb5689d6aaaf0adf9c580ac18ab72a99eae8c13f03f6` |
| Independent integration review | `636a5f8d91ba569fd0c5fed45bd00621a1e2d0aafc6fbb4bfbb902e93ff9374a` |
| R5 proposal | `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e` |
| Original matrix | `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6` |
| Precision addendum 1 | `d059054b678f95d8697b8987bf856035a20a6415bd61828a940ebbd99d7549b0` |
| Precision addendum 2 | `be399e5bfda3fec62e4356ad8e0dadde11a643881bd74fc314ba2bf9b116e403` |
| Frozen Master | `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31` |
| Gate v2 | `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a` |

This package's own SHA-256 is recorded in the completion report rather than embedded here, to avoid a self-referential file hash.
