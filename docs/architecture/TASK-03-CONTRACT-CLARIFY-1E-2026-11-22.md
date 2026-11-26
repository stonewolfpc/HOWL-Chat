# TASK-03 — clarification 1e: identity-bound deletion, corrected state machine, absent-target stability

| Field | Value |
|---|---|
| Purpose | Record owner decision D-10, correct the deletion state machine, assess Windows feasibility |
| Status | **Contract text. Nothing adopted. No implementation authorized. Deletion remains blocked** |
| Relationship | **Additive.** None of rows 528–533 is edited |
| Supersedes | P-06, P-07, P-09 (path_recheck rejected), P-03 rows for `removing`/`removed` (identity verification changes them), L-01/L-02 (ordering and journal schema), K-01/L-20 (path stored in journal) |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1E-2026-11-22` |

Prefixes: `Q-` normative, `QT-` acceptance test.

## 0. Evidence labels

Every non-trivial claim below is labelled. Nothing in this document was produced by running code.

| Label | Meaning |
|---|---|
| **Observed** | Seen in this project's own files or environment |
| **Source-derived** | Taken from vendor/protocol documentation cited inline |
| **Calculated** | Arithmetic or logical consequence of stated facts |
| **Inferred** | Judgement from evidence, stated as such |
| **Unestablished** | Not proven by available evidence; the required evidence is named |
| **Superseded** | Replaced by a later clause; the earlier text is preserved |

---

## 1. Owner decision D-10 is recorded

**D-10 — decided by the owner.** Conversation deletion must remain bound to the exact physical directory identity acquired for that operation. Path-based deletion verified by re-checking is **rejected** for Phase 1.

**Q-01 — Phase 1 deletion is disabled until identity-bound removal is demonstrated.**

> No Phase 1 path may delete a conversation directory until an engineering spike has demonstrated identity-bound removal end to end and the result is recorded as evidence. Until then the deletion operation returns `deletion_unavailable` and changes nothing.

This is the direct consequence of D-10 plus an `Unestablished` feasibility verdict (§4). It is not a new decision; it is what D-10 means while the mechanism is unproven.

**Q-02 — The rejected basis, recorded.** `path_recheck` is removed from the vocabulary. `removal_basis` has exactly one legal value, `handle_identity`. A receipt carrying any other basis is invalid.

**Q-03 — No silent fallback.** If identity-bound removal is unavailable for a target — unsupported filesystem, unopenable handle, missing privilege, unresolvable identity — the operation fails closed with `deletion_unavailable` or `identity_unresolved`. `RemoveAll`, `DeleteFile`, `rd /s`, and any path-based recursive delete are **prohibited in the deletion path without exception**, including cleanup of a partially deleted tree.

---

## 2. Requirements the mechanism must satisfy

**Q-04 — Identity is retained across the destructive phase.** The operation holds the acquired identity — `(volume_serial, file_id_128)` — durably and uses it, not a path, for every destructive and verification step.

**Q-05 — Enumeration is identity-bound.** Children are enumerated through a handle on the acquired directory. Each child is identified by the file ID returned by that enumeration, not by a re-resolved path.

**Q-06 — Each child is removed by identity.** A child is removed through a handle obtained from its ID, on the same volume, with reparse points not followed.

**Q-07 — Replacement cannot redirect deletion.** If the original path is replaced by a different directory, the operation does not follow it: it operates on the identity it acquired, or fails.

**Q-08 — Reparse points are removed as objects, never traversed.** A junction, symbolic link or mount point encountered as a child is deleted as the link object. Its target — on any volume — is never entered or removed.

**Q-09 — Outcome is verified against identity, not path.** `removed` is written only after the acquired identity is shown to be absent. Path absence alone never advances the state machine and never supports `complete`.

**Q-10 — Partial failure is reported as partial, resumably.** If removal stops partway, the operation records `partial`, keeps the identity, and a later explicit attempt may continue from the same identity. Nothing is rolled back (rollback is impossible) and nothing is reported as complete.

**Q-11 — Every step that cannot be established fails closed.** There is no best-effort mode. An unopenable child, a volume boundary, an unsupported filesystem, or an unexpected reparse interaction aborts the operation with `partial` or `deletion_unavailable` — never with a path-based fallback.

---

## 3. Consequence: the journal must store identity, not a path

**Finding, forced by D-10.** The journal currently stores `directory_path` (K-01, L-02, L-20). A path cannot identify: after a replacement, the path names a different directory, and verification against "the path" is exactly what D-10 forbids.

**Q-12 — Journal schema v2.**

| Field | Type | Note |
|---|---|---|
| `journal_v` | int | `2`; a `1` entry is `journal_unversioned` for deletion purposes |
| `operation_id` | UUIDv4 | joins to the receipt |
| `volume_serial` | uint64 | from `FILE_ID_INFO` at acquisition |
| `file_id` | 128-bit | from `FILE_ID_INFO` at acquisition |
| `entries_at_acquisition` | uint64 | child count baseline; makes `partial` calculable |
| `state` | `prepared`/`pending`/`removing`/`removed` | progress barrier |
| `updated_at` | RFC3339 UTC | progression evidence |

**Q-13 — The path is no longer durable.** No journal field, receipt field, log line or error message contains the conversation path. The path exists only in the caller's request and the in-memory call.

*Consequences:*
- **Privacy improves** — the residual risk recorded in K-07/L-20 (a path surviving a crash) disappears. What persists is an identity that becomes unresolvable once the directory is deleted.
- **Reconciliation needs no path**: it locates the volume by serial, opens a volume handle, and resolves the identity with `OpenFileById`. A volume that is no longer mounted yields `indeterminate`/`volume_unavailable`, not a scan.
- `absent_prior` needs no entry at all (P-02), so the absent case never requires durable identity.

**Q-14 — Ownership comparison uses identity.** A new request resolves the target's identity and compares `(volume_serial, file_id)` with non-terminal entries. String comparison remains prohibited (N-12, P-14). *Calculated:* two paths that resolve to one directory can no longer acquire two entries, and neither can a path and its junction alias.

---

## 4. Engineering feasibility: Windows identity-bound directory removal

### 4.1 What the platform documents

| Fact | Evidence | Label |
|---|---|---|
| `FILE_DISPOSITION_INFORMATION_EX { ULONG Flags }` with `DELETE 0x1`, `POSIX_SEMANTICS 0x2`, `FORCE_IMAGE_SECTION_CHECK 0x4`, `ON_CLOSE 0x8`, `IGNORE_READONLY_ATTRIBUTE 0x10` | `_FILE_DISPOSITION_INFORMATION_EX` (ntddk.h), MS Learn | **Source-derived** |
| `FILE_DISPOSITION_DELETE` requires **DELETE** access on the handle | same page, Remarks | **Source-derived** |
| POSIX-style delete removes the link from the visible namespace "as soon as the POSIX delete handle is closed"; data streams stay accessible to other open handles | MS-FSCC `FileDispositionInformationEx` | **Source-derived** |
| `FILE_ID_INFO { ULONGLONG VolumeSerialNumber; FILE_ID_128 FileId; }`; serial + ID uniquely identify a file on one computer; compare both to decide whether two handles are the same file | `FILE_ID_INFO` (winbase.h); `BY_HANDLE_FILE_INFORMATION` Remarks | **Source-derived** |
| `FileIdInfo` information class (0x12) requires Windows 8 / Server 2012 or later; the 64-bit ID is not guaranteed unique on ReFS | `GetFileInformationByHandleEx`; `BY_HANDLE_FILE_INFORMATION` Remarks | **Source-derived** |
| Directory enumeration on a handle: `FileIdBothDirectoryInfo` (0xA) / `FileIdExtdDirectoryInfo` (0x13), returning file IDs | `GetFileInformationByHandleEx` | **Source-derived** |
| `OpenFileById(hVolumeHint, FILE_ID_DESCRIPTOR, …)` opens by ID; "will either open the file or the reparse point, depending on the use of `FILE_FLAG_OPEN_REPARSE_POINT`" | `OpenFileById`; `Reparse Points and File Operations` | **Source-derived** |
| `FILE_ID_DESCRIPTOR` supports `FileIdType`, `ObjectIdType`, `ExtendedFileIdType` (128-bit, Win8+) | `FILE_ID_DESCRIPTOR` | **Source-derived** |
| Per-directory case sensitivity: `FILE_CASE_SENSITIVE_INFORMATION { ULONG Flags }`, `FILE_CS_FLAG_CASE_SENSITIVE_DIR 0x1`, info class `FileCaseSensitiveInformation` = 71, directories only; available from Windows 10 build 17107 | `_FILE_CASE_SENSITIVE_INFORMATION` (ntifs.h); Windows Command Line blog; WSL docs | **Source-derived** for the structure and flag; **Unestablished** as a Win32-documented contract, since class 71 is absent from the published `FILE_INFO_BY_HANDLE_CLASS` enum and is reached through `GetFileInformationByHandleEx` with a numeric class |
| Directory-relative opens exist at the native layer: `OBJECT_ATTRIBUTES.RootDirectory` + relative `ObjectName` | `_OBJECT_ATTRIBUTES`; `NtCreateFile` | **Source-derived** |
| `RootDirectory` is documented as "a handle to a directory obtained by a preceding call to **NtCreateFile**" | `NtCreateFile` (winternl.h) | **Source-derived** — and therefore *not* documented as accepting a `CreateFileW` handle |
| `FILE_OPEN_REPARSE_POINT` is not listed among options compatible with `FILE_DIRECTORY_FILE` | Q&A analysis citing the documented compatibility list | **Unestablished** whether the combination is a supported contract |

### 4.2 The sequence the contract would require

```
1  resolve target identity via PathService            (classification only)
2  open root handle: DELETE | FILE_LIST_DIRECTORY | READ_ATTRIBUTES | WRITE_ATTRIBUTES,
      share = READ | WRITE, EXISTING, BACKUP_SEMANTICS | OPEN_REPARSE_POINT
   — omitting FILE_SHARE_DELETE is the mechanism that denies others DELETE access
3  GetFileInformationByHandleEx(FileIdInfo) → compare with acquired identity; mismatch → abort
4  enumerate: FileIdExtdDirectoryInfo / FileIdBothDirectoryInfo on the root handle
5  for each child entry, by ID:
      open by ID on the volume, OPEN_REPARSE_POINT, DELETE access
      if directory and not a reparse point: recurse (must be empty at its own level)
      if reparse point: delete the link object only
      SetFileInformationByHandle(FileDispositionInfoEx, DELETE|POSIX_SEMANTICS|IGNORE_READONLY_ATTRIBUTE)
      close → link removed
6  root: SetFileInformationByHandle(FileDispositionInfoEx) → close
7  verify: OpenFileById(acquired ID) must fail (not found) → write `removed`
8  terminal receipt → cleanup
```

**Steps 2–3 and 5–7 are the load-bearing ones**, and each contains at least one `Unestablished` dependency: whether a `CreateFileW` handle is a valid `RootDirectory`; whether `FILE_DIRECTORY_FILE` + `FILE_OPEN_REPARSE_POINT` is supported; whether share-mode exclusion actually prevents a directory from being renamed or converted to a reparse point while held (a Q&A observation suggests conversion with only `FILE_WRITE_ATTRIBUTES` may be possible); whether `OpenFileById` failing actually means "this identity no longer exists" on every supported filesystem rather than "not reachable by that route".

### 4.3 Verdict

> **UNESTABLISHED.** The building blocks exist and are individually documented, but no published contract establishes the complete sequence, and three of its steps depend on undocumented or version-specific behaviour.

**Evidence required to move this to Established:** a spike on the target platform that demonstrates, with recorded output:

1. identity acquisition and re-acquisition by `OpenFileById` across a process restart;
2. deletion of a nested tree containing a file, a subdirectory, a file symlink, a directory junction, and a mount point — with proof the link targets survive;
3. a replacement race: replace the directory between acquisition and removal, and show the replacement is untouched;
4. behaviour when the directory is renamed while held, and when a child is created mid-deletion;
5. behaviour on NTFS, ReFS, a FAT/exFAT volume and an SMB share (or explicit refusal on each);
6. verification semantics: what `OpenFileById` returns for a deleted identity on each filesystem;
7. privilege requirements (`SeBackupPrivilege`/`SeRestorePrivilege`) and behaviour without them;
8. minimum OS build behaviour for class 71 and `FileDispositionInfoEx`.

Until all eight are recorded, D-10 cannot be satisfied and Q-01 keeps deletion disabled.

---

## 5. Capability boundary

**Q-15 — A distinct responsibility is required, and it is not PathService.**

| Existing | Responsibility | Not its job |
|---|---|---|
| PathService (TASK-01) | path classification, containment, canonical physical identity of an existing path | destructive operations; holding handles across a lifecycle; journal coordination |
| ContainmentChecker | grant/containment decisions | removal |
| **Deletion authority** (new, contract-level) | acquire identity, hold handles, enumerate and remove by identity, verify by identity, coordinate journal states | path classification — it *consumes* PathService output |

Justification (the demonstrated need required by the review): D-10's guarantee needs a component that holds OS handles across a multi-step lifecycle and performs destruction. PathService is shared with read paths and must remain non-destructive; folding deletion into it would give every reader a destructive capability. That is a correctness boundary, not code reuse.

**Q-16 — No duplicate persistence.** The deletion authority owns no ledger, no receipt store and no second journal. It writes only the journal states of §6 and the receipt records of C-33.

---

## 6. Corrected state machine and crash semantics

**Q-17 — Replaces L-01's ordering.**

```
prepared → intent → pending → removing → [destructive removal] → verify identity absent → removed → terminal → cleanup
```

`removing` is durably written **before** the destructive call. It marks crossing the pre-removal boundary — it does **not** prove the call ran.

**Q-18 — Removal evidence, derived and never invented.**

| Value | Derivation |
|---|---|
| `not_started` | durable state is `prepared` or `pending` — no destructive call could have been made |
| `possible` | durable state is `removing` — the boundary was crossed; whether the call ran is unknowable after a crash |
| `verified` | durable state is `removed` — identity confirmed absent after removal |

**Q-19 — State meanings.**

| State | Meaning |
|---|---|
| `prepared` | entry exists; no intent; removal cannot have been invoked |
| `pending` | intent exists; destructive phase not entered; removal cannot have been invoked |
| `removing` | destructive phase entered; removal may or may not have been invoked |
| `removed` | removal completed **and the acquired identity verified absent** |

### 6.1 Crash-window table

Ten windows. `IE` = removal evidence. Attribution is `operation` only where `IE = verified` and identity was verified absent.

| # | Crash point | Durable state | Removal invoked? | Target could be | IE | Attribution possible | Reconciliation may conclude | Must remain indeterminate |
|---|---|---|---|---|---|---|---|---|
| 1 | before `prepared` | none | no | present (as observed by caller) | `not_started` | none | nothing to reconcile — no operation exists | whether a request was ever made |
| 2 | after `prepared`, before intent | `prepared`, no intent | no | present | `not_started` | none | `abandoned_no_intent`; remove entry; **no receipt** | nothing (no deletion was requested) |
| 3 | after intent, before `pending` | `prepared` + intent | no | present | `not_started` | none | `aborted`, `removed 0`; then cleanup | — |
| 4 | after `pending`, before `removing` | `pending` | no | present / absent | `not_started` | none | present → `aborted`; absent → `absent_after_acquisition` | who removed it if absent |
| 5 | after `removing`, before the call | `removing` | **unknown** | present / absent / partial | `possible` | none | present+intact → `aborted`; absent → `absent_unverified`; partial → `partial` | whether the call ran, and who caused absence |
| 6 | during removal | `removing` | yes | partial (typical), absent, present | `possible` | none | same as #5, with `partial` most likely | how much was removed by us |
| 7 | after the call returns, before verification | `removing` | yes | absent (probably) | `possible` | none | `absent_unverified` | whether absence is ours |
| 8 | after verification, before `removed` | `removing` | yes | absent | `possible` | none | `absent_unverified` — an in-memory verification that never reached disk is **not evidence** | whether verification happened |
| 9 | after `removed`, before terminal | `removed` | yes | absent (verified) | `verified` | `operation` | `complete`, write terminal, cleanup | — |
| 10 | after terminal, before cleanup | `removed` + terminal | yes | absent | `verified` | `operation` | return the receipt unchanged; remove entry; never re-delete | — |

**Q-20 — Two rules the table exists to enforce.**

> Absence never proves removal happened. A pre-removal state never proves removal did not happen — only a state whose ordering forbids the call does that.

Window 5 is the one that makes this concrete: after `removing` is durable, the process may have crashed before or during the call, and no post-crash evidence distinguishes them.

**Q-21 — New outcomes the correction forces.**

| Outcome | Meaning |
|---|---|
| `absent_unverified` | the target is absent after the destructive boundary was crossed, but identity-absent was never verified — so `complete` is unavailable |
| `deletion_unavailable` | identity-bound removal cannot be performed on this target/platform |

`absent_unverified` **replaces** `complete`/`unknown` for `removing` + absent (P-03 row 6). Under D-10, path absence can no longer support completion at all.

---

## 7. Audit of the state/presence table

The table in P-03 was read from the document on disk, not from memory. Findings:

**Finding 1 — incomplete by one row.** P-03 enumerates 11 `(state, presence)` combinations for a present intent, but the matrix is 4 states × 3 presence values = **12**. The missing combination is `removed` + *partly emptied*, which must exist even if only to be declared impossible.

**Finding 2 — `removing` + absent was wrong under D-10.** P-03 gave `complete`/`unknown`. With identity verification required (Q-09), absence of the path cannot support `complete`. Corrected to `absent_unverified`.

**Finding 3 — `removed` + present was misframed.** P-03 called path presence `state_inconsistent`. Under identity-bound deletion a replacement at the path is *expected* and irrelevant: the receipt is about the acquired identity. Path occupancy by another directory is not an inconsistency. The real inconsistency is the **identity** still existing after a verified removal.

**Finding 4 — presence must be evaluated by identity.** All rows are re-expressed against identity presence, with path occupancy recorded only as a diagnostic.

**Q-22 — Corrected table. Replaces P-03.**

Presence below means **presence of the acquired identity**. `entries < baseline` means the child count is below `entries_at_acquisition`.

| # | State | Intent | Identity presence | IE | Identity verified absent | Outcome | Attribution |
|---|---|---|---|---|---|---|---|
| 1 | — | — | absent at acquisition | `not_started` | — | `absent_prior` (no entry, no receipt-of-deletion) | `prior_absence` |
| 2 | `prepared` | absent | any | `not_started` | — | `abandoned_no_intent` — entry removed, no receipt | — |
| 3 | `prepared` | present | absent | `not_started` | no | `absent_after_acquisition` | `unknown` |
| 4 | `prepared` | present | present, intact | `not_started` | no | `aborted`, removed 0 | — |
| 5 | `prepared` | present | present, entries < baseline | `not_started` | no | `partial` | `unknown` |
| 6 | `pending` | present | absent | `not_started` | no | `absent_after_acquisition` | `unknown` |
| 7 | `pending` | present | present, intact | `not_started` | no | `aborted`, removed 0 | — |
| 8 | `pending` | present | present, entries < baseline | `not_started` | no | `partial` | `unknown` |
| 9 | `removing` | present | absent | `possible` | no | **`absent_unverified`** | `unknown` |
| 10 | `removing` | present | present, intact | `possible` | no | `aborted`, removed 0 | — |
| 11 | `removing` | present | present, entries < baseline | `possible` | no | `partial` | `unknown` |
| 12 | `removed` | present | absent | `verified` | yes | `complete` | `operation` |
| 13 | `removed` | present | present (identity still resolvable) | `verified` | claimed | `indeterminate`, `identity_present_after_verified_removal` | `unknown` |
| 14 | `removed` | present | present, entries < baseline | `verified` | claimed | `indeterminate`, `state_inconsistent` | `unknown` |
| 15 | `pending`/`removing`/`removed` | absent | any | — | — | `indeterminate`, `receipt_missing`; entry preserved | `unknown` |
| 16 | any | any | any, terminal exists | — | — | return terminal unchanged, then cleanup | — |

**Row 13 note.** A *different* directory occupying the old path is not row 13 — that is row 12 with a path note, and it must not downgrade the outcome. Row 13 is only for the acquired identity itself resolving after a verified removal.

**Q-23 — Impossible combinations, declared rather than omitted.** `prepared` with `IE = possible`, `pending` with `verified`, `removed` with `not_started`/`possible`, and any row where `attribution = operation` without `state = removed` are structurally impossible; a persisted entry exhibiting them is `state_inconsistent` and is reported, never repaired.

---

## 8. Absent-target observation

**Q-24 — The race.** Resolving the parent, observing child absence, then having the parent renamed or replaced yields an `absent_prior` that describes a directory that is no longer the same object.

**Q-25 — Two admissible bases, strongest first.**

| Basis | Mechanism | Status |
|---|---|---|
| `handle_relative` | child checked relative to a retained parent handle (`RootDirectory`) | **Unestablished** — documented for `NtCreateFile`, not guaranteed for a `CreateFileW` handle, and `FILE_DIRECTORY_FILE` + `FILE_OPEN_REPARSE_POINT` compatibility is undocumented |
| `pinned_parent` | parent handle held **without `FILE_SHARE_DELETE`**, then the child checked by path | **Source-derived** mechanism, **Unestablished** as an absolute guarantee (share modes gate DELETE access; a Q&A observation suggests a directory may be convertible to a reparse point with only `FILE_WRITE_ATTRIBUTES`) |

**Q-26 — If neither basis is available, absence is not definitive.**

> The result is `absence_unverified` with `observation_basis: none`. No receipt is written, no ownership is taken, and the caller is told absence could not be established safely. It is **not** reported as `absent_prior`.

**Q-27 — What `absent_prior` requires.** A retained parent handle at observation time, the child name validated (Q-28), absence determined relative to that handle or under a pinned parent, and `observation_basis` recorded. Anything less is `absence_unverified`.

**Q-28 — Child-name validation.** Exactly one component: no separators, no `.`/`..`, no trailing dot or space, no reserved device name (`CON`, `PRN`, `AUX`, `NUL`, `COM1`–`COM9`, `LPT1`–`LPT9`, including superscript variants). A name failing validation is `invalid_child_name`, not "absent".

---

## 9. Name comparison

**Q-29 — The case rule is per-directory, never global.** Windows supports per-directory case sensitivity (`FILE_CS_FLAG_CASE_SENSITIVE_DIR`, class 71, Windows 10 17107+). A filesystem-wide assumption is wrong.

**Q-30 — The rule is read from the resolved parent handle.** Query class 71 on the parent handle at the moment the composite identity is formed. The result is recorded with the observation.

**Q-31 — If the rule cannot be established, fail closed.** A query failure, an unsupported filesystem, or an inconclusive answer yields `identity_unresolved`/`case_rule_unknown`. Guessing case-insensitive is prohibited: it would merge two distinct children on a case-sensitive directory.

**Q-32 — Comparison is ordinal UTF-16 under the directory's own rule.**

| Aspect | Rule |
|---|---|
| Base comparison | ordinal UTF-16 code-unit comparison |
| Case | folded only if the parent's own rule says case-insensitive |
| Normalization | **never** applied by this contract. No NFC/NFD/NFKC step. NTFS does not normalize names; applying an application-level normalization would change identity, which P-14/Q-14 forbid |
| Unicode equivalence | names that are canonically equivalent but not byte-identical are **different** names. Stated as a limitation, not a bug |
| Reserved names | rejected at validation (Q-28) |
| Trailing dot/space | rejected; Win32 strips them, so they cannot be addressed unambiguously |
| Unrepresentable names | if a name cannot be represented in UTF-16 comparison without loss, `identity_unresolved` |

---

## 10. Acceptance tests

**Crash windows (one per window):**

- **QT-01** W1: killed before `prepared` → no entry, no intent, no receipt; a later request behaves as a first request.
- **QT-02** W2: `prepared` with no intent → `abandoned_no_intent`, entry removed, receipt store byte-identical.
- **QT-03** W3: `prepared` + intent → `aborted`, `removed 0`, `removal_evidence: not_started`.
- **QT-04** W4: `pending` → present gives `aborted`; absent gives `absent_after_acquisition`; neither yields `complete`.
- **QT-05** W5: forced crash between the `removing` write and the call → `removal_evidence: possible`; presence-dependent outcome; never `complete`, never a claim that the call did not run.
- **QT-06** W6: crash during removal → `partial` (or `absent_unverified` if target gone), attribution `unknown`.
- **QT-07** W7: crash after return before verification → `absent_unverified`, not `complete`.
- **QT-08** W8: crash after in-memory verification but before the `removed` write → `absent_unverified`; the lost verification is not evidence.
- **QT-09** W9: `removed` durable → `complete` with `attribution: operation`; terminal written exactly once.
- **QT-10** W10: terminal durable, entry stale → receipt returned unchanged, entry removed, no second receipt, no filesystem change in the conversation directory.

**State/presence matrix:**

- **QT-11** All 16 rows of Q-22 are reachable in a harness and produce the stated outcome; no row produces `complete` without `state = removed` and `identity_verified = true`.
- **QT-12** Row 12 with a *different* directory now at the path → still `complete`; the replacement is untouched.
- **QT-13** Row 13 (acquired identity resolvable after verified removal) → `indeterminate`, `identity_present_after_verified_removal`.
- **QT-14** Every structurally impossible combination (Q-23) is rejected as `state_inconsistent`, never repaired.
- **QT-15** `removal_evidence` is always derivable from durable state alone; no code path sets it from memory.
- **QT-16** `partial` is calculated from `entries_at_acquisition`, not estimated.

**Identity-bound removal:**

- **QT-17** Replacement between acquisition and removal: the replacement survives, outcome is `aborted`/`deletion_unavailable`, never a deletion of the replacement.
- **QT-18** Nested tree with a file symlink, a directory junction and a mount point: links removed as objects, targets intact, verified by content hash.
- **QT-19** Cross-volume junction target: never entered; no file outside the conversation directory is touched.
- **QT-20** Verification is by `OpenFileById` on the acquired identity; no path-based check may set `identity_verified`.
- **QT-21** `RemoveAll`/`DeleteFile`/`rd /s` appear nowhere in the deletion path (call-graph assertion).
- **QT-22** Unsupported filesystem or missing privilege → `deletion_unavailable`; nothing deleted, nothing written.

**Absent target:**

- **QT-23** Parent renamed between resolution and child check → `absence_unverified`, not `absent_prior`.
- **QT-24** Parent handle cannot be held without share-delete → `absence_unverified`.
- **QT-25** Genuine absence under a pinned parent → `absent_prior` with `observation_basis: pinned_parent`.

**Name comparison:**

- **QT-26** Case-insensitive directory: `Foo` and `foo` compare equal; case-sensitive directory: distinct — both asserted on one machine with both directory kinds.
- **QT-27** Class 71 unavailable → `identity_unresolved`/`case_rule_unknown`; no case-insensitive default.
- **QT-28** Reserved device name, trailing dot, trailing space, separator, `..` → `invalid_child_name`, each distinct.

---

## 11. What remains unproven

| # | Question | Evidence needed |
|---|---|---|
| E-1 | Is a `CreateFileW` directory handle accepted as `RootDirectory`? | Documented statement or spike across builds |
| E-2 | Is `FILE_DIRECTORY_FILE` + `FILE_OPEN_REPARSE_POINT` supported? | Documented compatibility list entry or spike |
| E-3 | Does omitting `FILE_SHARE_DELETE` prevent rename **and** reparse conversion of a directory? | Adversarial spike |
| E-4 | Does `OpenFileById` failing prove identity absence on NTFS, ReFS, FAT, SMB? | Per-filesystem spike |
| E-5 | Mount-point deletion: does removing the link leave the volume intact? | Destructive spike on a scratch volume |
| E-6 | Minimum OS build for class 71 and `FileDispositionInfoEx` semantics | Build matrix |
| E-7 | Privilege requirements without `SeBackupPrivilege` | Spike as a non-elevated user |
| E-8 | Behaviour of POSIX delete on a directory with open child handles | Spike |
| E-9 | Exact upcase-table semantics for ordinal comparison | Documentation or per-filesystem spike |
| E-10 | Whether a reparse point's file ID identifies the link or the target on each filesystem | Spike |

**Feasibility verdict: `UNESTABLISHED`.** Per Q-01, deletion stays disabled.

---

## 12. Owner decisions

| # | Decision | Status |
|---|---|---|
| 1 | Retention Option A | **Decided** |
| 2 | Envelope with required predecessor | Open |
| 3 | Five ledger + two evidence kinds | Open |
| 4 | `FC1-02` linear | Open |
| 5 | `FC1-07` legacy read-only | Open |
| 6 | `FC1-14b` excluded/unadopted | Open |
| 7 | `FC1-14c` one message per turn | Open |
| 8 | Receipt without the five excluded items | Open |
| 9 | Per-conversation deletion proof | Open |
| 10 | Identity-bound deletion mandatory | **Decided** |
| 11 | Enable deletion before the spike completes, or keep it disabled until E-1…E-10 are recorded | **Open — new** |

---

## 13. Authorization and evidence

No implementation, migration, ledger creation, persistence adoption, deletion enablement, or FC-1 approval. No Go file was modified. No test of any kind was executed: the content above is contract text plus **source-derived** API facts, and the feasibility verdict is `UNESTABLISHED` by design. TASK-03 is not accepted and not closed.
