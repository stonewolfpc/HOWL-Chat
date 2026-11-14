# OWNER DECISION REQUEST — F3C FOUNDATIONAL AUTHORITY

**Status:** **DECISION REQUEST ONLY — NOT R6, NOT A CONTRACT REVISION, NOT A DECISION RECORD, NOT FREEZE, NOT ACCEPTANCE, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Companion adjudication:** `docs/review/ADJUDICATION-PACKAGE-F3C-OD18-INTEGRATION-2026-11-11.md` (`85e1968bac4eb70f489bbd7dc89057debd17e7b3bc2336aa69c036a4352f83e8`)

**Purpose:** request the five foundational owner decisions that must exist before the state machine and lock proof can be finalized, and before R6 can be considered.

**Rule:** no option in this document is approved by being recommended, listed, or technically attractive. A decision exists only when the responsible owner records it. A recorded decision is not a formal amendment until the required authority accepts it. Neither is implementation evidence.

---

## 1. Decision D-ROOT — Canonical F3-C root ownership and legacy-root status

### 1.1 Why this must be decided first

The source currently exposes more than one root resolution surface:

- `internal/paths/paths.go` — `paths.StateRoot()`; defaults to the installation/project root when isolation is off.
- `internal/backend/chat/persistence.go` — legacy chat root at `%USERPROFILE%\AppData\Roaming\HOWLChat` when isolation is off.
- `internal/backend/chat/service.go` — uses `paths.StateRoot()` for chat/service state.

A lock over one root does not protect another. Until root ownership is recorded, OD-18 has no defined lock scope.

### 1.2 Required findings

The owner decision must answer all four:

1. **Which root owns F3-C chat state?**
2. **Does the legacy `%USERPROFILE%\AppData\Roaming\HOWLChat` path remain active for any F3-C read or write?**
3. **Can both paths resolve to the same underlying data through different routes?**
4. **Does every relevant caller use the same canonical root identity, rather than merely the same path string?**

### 1.3 Options

| Option | Meaning | Consequence |
|---|---|---|
| **A. Single canonical state root; legacy quarantined** *(recommended for evaluation)* | `paths.StateRoot()` or one explicitly configured state root owns F3-C. Legacy AppData is a migration source only; no F3-C reads or writes. | Clearest lock scope, but requires migration and caller audit. |
| **B. Legacy AppData remains canonical** | The AppData root owns F3-C; `paths.StateRoot()` is not authoritative for chat state. | Requires proof that current `service.go` and all callers are redirected. |
| **C. Dual-root with explicit root identity** | Both roots are recognized, each with a distinct identity and lock, and callers select one explicitly. | Only acceptable if no caller can cross roots accidentally and no alias can bypass a lock. |
| **D. Undetermined** | No canonical root recorded. | **Blocks OD-18, state-machine recovery, and all lock-proof work.** |

### 1.4 Evidence required before the decision can be accepted

- A complete inventory of F3-C root resolutions and callers.
- A canonical identity for each root: volume identity plus file identity, not only a normalized string.
- Resolution of reparse points, junctions, symlinks, 8.3 names, case variants, drive-letter changes, and UNC/mapped-drive aliases.
- Proof that the selected canonical root cannot be reached through a second unprotected route.
- A migration/quarantine rule for any legacy F3-C state.
- A fail-closed rule for ambiguous identity, missing roots, or mixed-root state.

### 1.5 Decision record

```
D-ROOT:
Selected option:
Canonical root identity:
Legacy-root disposition:
Callers covered:
Alias/reparse treatment:
Fail-closed conditions:
Owner:
Date:
Authority record:
```

**Status:** **UNDECIDED — BLOCKS OD-18.**

---

## 2. Decision D-STORAGE — Supported Windows storage boundary

### 2.1 Question

Which Windows storage configurations may run F3-C with an enforceable single-writer claim, and which must fail closed as unsupported?

### 2.2 Candidate boundary

**Candidate only:** local fixed NTFS volume, one authority process per root, with an OS-enforced exclusive handle or lock held for process lifetime.

**Excluded until separately proved and explicitly approved:**

- SMB/network shares;
- ReFS;
- removable media;
- OneDrive, Dropbox, or other cloud-synced folders;
- redirected or roaming profiles that can synchronize;
- containers and virtual machines;
- snapshot-restored, cloned, or duplicated volumes;
- split mounts or alternate paths to the same storage;
- any configuration whose filesystem/share type cannot be determined.

“Local NTFS looks plausible” is not evidence. No unsupported configuration may receive an optimistic fallback.

### 2.3 Required Windows evidence

No test in this matrix has been run for this decision request. All are **required before an OD-18 enforcement claim**.

| Test | Scenario | Required result |
|---|---|---|
| T1 | Two independent Howl processes attempt authority. | Exactly one acquires; the other performs no ordinary F3-C access. |
| T2 | Holder is terminated normally and forcibly. | OS releases authority; no second writer overlaps the terminated holder. |
| T3 | Stale, inherited, or duplicated handles exist. | No process gains authority from a stale handle; ambiguity fails closed. |
| T4 | Lock file/path is deleted, renamed, or replaced. | No deletion/recreation grants authority; ordinary access is refused. |
| T5 | Suspend/resume, reboot, PID reuse, process incarnation change. | No process assumes authority from a prior incarnation; ambiguity fails closed. |
| T6 | Storage error, access denied, timeout, or ambiguous acquisition. | No optimistic fallback; ordinary access is refused. |
| T7 | A second process attempts access during startup recovery. | Recovery is protected by the same authority; second process is blocked. |
| T8 | Junction, symlink, 8.3 name, or mapped path reaches the same root. | Same canonical identity is detected, or access fails closed. |
| T9 | Shutdown/release races with a new process. | No overlap; successor acquires only after confirmed release. |

### 2.4 Required enforcement semantics

If D-STORAGE selects local NTFS:

- acquire authority before any F3-C storage access, including startup recovery and read-only inspection;
- hold it for the authority process lifetime;
- never steal or clear a lock by PID, timestamp, file content, or “looks stale” inference;
- never proceed after a retry budget;
- block reads, writes, repair, deletion, wipe, and successor projection on ambiguity;
- record root identity and acquisition outcome without recording secret material.

### 2.5 Decision record

```
D-STORAGE:
Supported configurations:
Excluded configurations:
Enforcement primitive:
Acquisition point:
Release point:
Failure behavior:
Evidence references:
Owner:
Date:
Authority record:
```

**Status:** **UNDECIDED — BLOCKS OD-18.**

---

## 3. Decision D-LAYOUT — Is frozen §4.1 exhaustive?

### 3.1 Question

Does frozen §4.1 permit only the enumerated layout, or may F3-C add authority/evidence paths outside that enumeration?

### 3.2 Options

| Option | Meaning | Consequence |
|---|---|---|
| **A. Exhaustive** | §4.1 is a closed list. No additional authority or evidence path is permitted. | CR-01 and related additions remain blocked. |
| **B. Enumerated additions** | Additional paths are permitted only by formal amendment that names each path and its ownership. | Requires an amendment; no implicit path creation. |
| **C. Undetermined** | No ruling recorded. | **Blocks CR-01 and any layout-dependent proof.** |

A decision that §4.1 is exhaustive does not authorize a new path. A decision that it is non-exhaustive does not enumerate one.

### 3.3 Decision record

```
D-LAYOUT:
Selected option:
Additional paths, if any:
Owning decision:
Formal amendment required:
Owner:
Date:
Authority record:
```

**Status:** **UNDECIDED — CR-01 BLOCKED.**

---

## 4. Decision D-COMPAT — Target replacement and v2 record authority

### 4.1 Question

What formal route authorizes repeated target replacement, multi-target installation, and F3C-TARGET-2 record/target field additions against the frozen contract?

### 4.2 Required separation

| Matter | Required authority |
|---|---|
| Per-file target replacement mechanics | OD-02a |
| Durable recovery, layout, and mutable-evidence authorization | OD-02b |
| V2 target/record compatibility or version ruling | OD-19 (proposed) plus OD-12 |
| Frozen create-only installation change | Formal amendment, not an implementation note |

**OD-02a must remain inert without OD-02b.** Neither may be treated as accepted because the other is proposed.

### 4.3 Options

| Option | Meaning | Consequence |
|---|---|---|
| **A. Formal amendment** | Amend the frozen create-only and schema provisions with explicit replacement and v2 rules. | Clearest authority; requires owner action. |
| **B. Compatibility ruling** | The owner rules that existing provisions permit the proposed behavior without amendment. | Must be recorded and must cover every CR-02/03/05/06/07/11/12 obligation. |
| **C. Undetermined** | No route recorded. | **Blocks R6 and implementation.** |

### 4.4 Decision record

```
D-COMPAT:
Target replacement route:
V2 target/record route:
OD-02a status:
OD-02b status:
OD-19 status:
Formal amendment reference:
Owner:
Date:
Authority record:
```

**Status:** **UNDECIDED — CR-02/03/05/06 BLOCKED.**

---

## 5. Decision D-LOCK — Single-writer mechanism and failure contract

### 5.1 Question

What exact mechanism enforces one authority process per canonical root, and what must happen when enforcement fails?

### 5.2 Minimum contract

The decision must name:

- the OS primitive and its sharing/exclusion semantics;
- the canonical root identity it protects;
- the acquisition point before any F3-C access;
- the release point and how release is proved;
- the behavior for stale handles, PID reuse, deletion/replacement, suspend/resume, I/O errors, and ambiguous outcomes;
- the relationship between the storage lock and frozen F3-C lock order;
- the rule for startup recovery, multi-chat recovery, wipe recovery, gate release, and successor projection;
- the rule that verification never acquires an upward lock while holding `chatWriteMu`;
- the rule that no destructive action runs on an unverified package.

### 5.3 Required state/event table fields

Any future transition table must specify, for every `(state, event)` pair:

| Field | Meaning |
|---|---|
| Preconditions | What must be true before the event |
| Durable evidence | What must already exist |
| Action | What recovery or runtime may do |
| Next state | Exactly one deterministic destination |
| Failure outcome | What happens if the action cannot complete |
| Access policy | Whether readers and writers remain blocked |

The table must distinguish:

- **absent record** from **corrupt record**;
- **absent package** from **unverifiable package**;
- **completed operation** from **visible-but-not-durable completion**.

No state may be exited by “best effort,” timeout, or a plausible-looking file.

### 5.4 Decision record

```
D-LOCK:
Mechanism:
Canonical root identity:
Acquisition point:
Release proof:
Stale-handle rule:
Deletion/replacement rule:
Error/ambiguous rule:
Interaction with frozen lock order:
Wipe/recovery rule:
Owner:
Date:
Authority record:
```

**Status:** **UNDECIDED — LOCK PROOF AND R6 BLOCKED.**

---

## 6. CR milestone ladder

The following are four distinct milestones. None implies the next.

| Milestone | Meaning |
|---|---|
| **A — Owner assigned** | A named owner has accepted responsibility for resolving the CR. A proposal in an unaccepted document is not assignment. |
| **D — Owner decision recorded** | The owner has chosen an option and recorded the rationale and scope. |
| **F — Formal amendment accepted** | The decision has the architectural authority required by the contract. |
| **I — Implementation verified** | The resulting behavior has source/runtime evidence and the required tests pass. |

### 6.1 CR register

| CR | Proposed owner | A | D | F | I |
|---|---|---|---|---|---|
| CR-01 | OD-17 | Proposed | No | No | No |
| CR-02 | AM-1 / OD-02a | Proposed | No | No | No |
| CR-03 | OD-02a + OD-02b | Proposed | No | No | No |
| CR-04 | OD-18 | Proposed | No | No | No |
| CR-05 | OD-19 (proposed) + OD-12 | Proposed | No | No | No |
| CR-06 | OD-19 (proposed) + OD-12 | Proposed | No | No | No |
| CR-07 | OD-01 + OD-02b | Proposed | No | No | No |
| CR-08 | OD-14 | Proposed | No | No | No |
| CR-09 | OD-15 | Proposed | No | No | No |
| CR-10 | None — scope prohibition | N/A | N/A | N/A | N/A |
| CR-11 | OD-02b | Proposed | No | No | No |
| CR-12 | OD-02b | Proposed | No | No | No |

**No CR has reached D, F, or I.** CR-10 is not a decision; it is a prohibition against citing the worldstore audit as chat-store evidence.

---

## 7. Required decision order

1. **D-ROOT** — canonical F3-C root and legacy-root disposition.
2. **D-STORAGE** — supported Windows configurations and fail-closed boundary.
3. **D-LAYOUT** — §4.1 exhaustiveness and permitted authority/evidence paths.
4. **D-COMPAT** — target replacement and v2 record authority.
5. **D-LOCK** — enforcement mechanism and failure contract.
6. **Transition table** — total deterministic state/event mapping.
7. **Whole-lifecycle lock proof** — against the selected root, storage boundary, and lock contract.
8. **R6** — only after the above.

The state machine and lock proof are not independent of D-ROOT, D-STORAGE, D-LAYOUT, D-COMPAT, or D-LOCK. Finalizing them before those decisions would silently decide them.

---

## 8. Status

- **This request:** open; no decision recorded.
- **D-ROOT:** undecided; blocks OD-18.
- **D-STORAGE:** undecided; local NTFS candidate only.
- **D-LAYOUT:** undecided; CR-01 blocked.
- **D-COMPAT:** undecided; CR-02/03/05/06 blocked.
- **D-LOCK:** undecided; lock proof and R6 blocked.
- **R6:** not authorized; not written.
- **F3-C phase 0:** not closed.
- **F3-C acceptance:** not granted.
- **F3-D:** blocked.
- **V20:** blocked pending F3-F evidence.
- **Gate v2:** V1–V17 adopted only.
- **A14 / V24:** open / blocked.

No source, tests, builds, migrations, model loads, destructive tests, application data, authority records, freeze, acceptance, or implementation action was performed.

---

## Appendix A — Identity

| Artifact | SHA-256 |
|---|---|
| Adjudication package | `85e1968bac4eb70f489bbd7dc89057debd17e7b3bc2336aa69c036a4352f83e8` |
| Integration review | `b29bb562c2a5b596baa3bb5689d6aaaf0adf9c580ac18ab72a99eae8c13f03f6` |
| Independent integration review | `636a5f8d91ba569fd0c5fed45bd00621a1e2d0aafc6fbb4bfbb902e93ff9374a` |
| R5 proposal | `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e` |
| Frozen Master | `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31` |
| Gate v2 | `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a` |

This request's own SHA-256 is recorded in the completion report rather than embedded here, to avoid a self-referential file hash.
