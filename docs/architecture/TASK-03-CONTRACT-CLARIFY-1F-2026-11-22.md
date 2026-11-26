# TASK-03 — clarification 1f: volume reacquisition, absent-target gate, absence semantics

| Field | Value |
|---|---|
| Purpose | Close the six review findings against Clarification 1e |
| Status | **Contract text. Nothing adopted. No implementation authorized. Deletion unavailable** |
| Relationship | **Additive.** None of rows 528–534 is edited |
| Supersedes | Q-12 (schema + locator), Q-01/decision-11 framing, 1e row 1, 1e rows 12–14, §8's `absent_prior` reachability |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1F-2026-11-22` |

Prefixes: `R-` normative, `RT-` acceptance test.

## 0. Evidence labels

**Observed** (this project's files/environment) · **Source-derived** (vendor/protocol documentation, cited) · **Calculated** (logical consequence) · **Inferred** (judgement, stated) · **Unestablished** (not proven; required evidence named) · **Superseded** (replaced, earlier text preserved).

Nothing below was produced by running code. No Windows behaviour has been empirically tested.

---

## 1. Volume reacquisition after process death

### 1.1 The gap

Q-12 stores `(volume_serial, file_id_128, entries_at_acquisition)`. That identifies an object **while a handle exists**, but after restart there is no handle, and `OpenFileById` needs `hVolumeHint` — "a handle to any file on a volume or share" (Source-derived: `OpenFileById`, winbase.h). Nothing in 1e says how that handle is obtained without a path. The gap is real: reconciliation after restart cannot currently execute.

### 1.2 What the platform provides

| Fact | Evidence | Label |
|---|---|---|
| Volume GUID path form `\\?\Volume{GUID}\`; a volume is not required to have a drive letter; letter assignments change as volumes are added/removed | *Naming a Volume* | **Source-derived** |
| "a volume GUID path can refer only to one volume", but "a volume can have more than one volume GUID path" | *Naming a Volume* | **Source-derived** |
| The OS assigns a volume GUID path when the volume is mounted; paths are obtained with `GetVolumeNameForVolumeMountPoint` or by enumeration with `FindFirstVolume`/`FindNextVolume`/`FindVolumeClose` | *Naming a Volume*, *Enumerating Volumes* | **Source-derived** |
| `CreateFile` on a volume GUID path with trailing backslash opens the volume root | *Naming a Volume* | **Source-derived** |
| Volume serial number is written when the volume is formatted, and third-party tools can change it | *Volume Serial Numbers* (Microsoft/Sysinternals VolumeID reference) | **Source-derived** (vendor doc); **Inferred** that it is not a stable locator |
| `OpenFileById` takes "a handle to any file on a volume or share"; access 0 permits attribute query | `OpenFileById` | **Source-derived** |

**Conclusion (Calculated).** The GUID path is the locator; the serial is a **corroborator**, because it can change independently of the volume's identity and cannot itself locate a volume.

### 1.3 The recovery process

**R-01 — Journal schema v3.** Adds `volume_guid_path` (`\\?\Volume{GUID}\`, with trailing backslash) and `filesystem` (e.g. `NTFS`, `ReFS`). Retains `volume_serial`, `file_id_128`, `entries_at_acquisition`, `state`, `updated_at`, `operation_id`, `journal_v = 3`.

*Privacy:* a volume GUID identifies a volume, not a conversation, and appears in no receipt. Q-13's prohibition on durable **conversation** paths stands; this is a volume locator and is stated as such.

**R-02 — Reacquisition, in order.**

```
1  open \\?\Volume{GUID}\                     → failure: volume_unavailable
2  query volume serial, filesystem            → serial mismatch: volume_identity_changed
                                              → filesystem mismatch: volume_identity_changed
3  OpenFileById(volume_handle, ExtendedFileIdType, stored file_id,
                FILE_READ_ATTRIBUTES, share = READ|WRITE|DELETE, OPEN_REPARSE_POINT)
   → not found: identity_absent
4  on success: GetFileInformationByHandleEx(FileIdInfo)
   → serial and file_id must equal the stored pair exactly, else identity_mismatch
   → confirm the object is a directory, else identity_mismatch
5  use the reopened handle only for verification or continuation; close it before any
   absence claim, then re-open fresh (R-14)
```

**R-03 — Failure taxonomy, all indeterminate and entry-preserving.** `volume_unavailable` (GUID resolves to nothing — volume absent, disconnected, reformatted, or renamed away), `volume_identity_changed` (serial or filesystem differs), `volume_unopenable` (access denied, device error), `identity_absent`, `identity_mismatch`, `identity_ambiguous` (two volumes in the enumeration answering to the same GUID — pathological, but must not be resolved by picking one).

**R-04 — No path fallback.** If the GUID path no longer resolves, reconciliation does **not** search for the directory by path, by name, by counts, or by timestamp. It reports `indeterminate` and stops.

**R-05 — Stability, stated honestly.**

| Scenario | Behaviour | Label |
|---|---|---|
| Reboot, same volume | GUID path unchanged; locator works | **Source-derived** (paths are mount-manager assigned and letter-independent) |
| Drive letter changed | GUID path unchanged | **Source-derived** |
| Volume unplugged | `volume_unavailable` | **Calculated** |
| Same volume re-plugged later | GUID path unchanged → reconciliation may resume | **Inferred**; needs the spike (E-11) |
| Volume reformatted | new GUID → `volume_unavailable`; old identity gone with the data | **Inferred** (a reformat destroys the directory anyway) |
| Volume cloned byte-for-byte | identical GUID, serial and IDs on two volumes — **indistinguishable** | **Unestablished** and not resolvable by this contract |

**R-06 — What a cloned volume does to the claim (limitation, stated).** Because the claim is scoped to "the volume reached by the stored locator", a clone presenting the same GUID is not distinguishable. A receipt therefore must never be read as "no copy of this conversation survives anywhere". This is the exact over-claim risk that `path_recheck` had, narrowed but not eliminated; it is recorded rather than denied.

### 1.4 Filesystems without reliable 128-bit identity

**R-07 — Identity support is a precondition, not a fallback.**

| Filesystem | 128-bit ID | Deletion |
|---|---|---|
| NTFS | `FileId128` = 48-bit MFT index (low) + 16-bit sequence (high), high 64 unused (Source-derived: MS-FSA Appendix A) | **Eligible**, pending spike |
| ReFS | 64-bit parent-directory index + 64-bit per-directory ordinal (Source-derived: MS-FSA) | **Eligible**, pending spike |
| FAT/exFAT | IDs derived from directory-entry position; defrag or rename can change them (Source-derived: `BY_HANDLE_FILE_INFORMATION` Remarks) | **Refused**: `identity_unsupported` |
| Anything reporting `FileId == 0` | "does not support a 128-bit file ID … must be ignored" (Source-derived: `FILE_ID_INFORMATION`) | **Refused**: `identity_unsupported` |
| Network shares | UNESTABLISHED | **Refused** until E-13 records evidence |

A refused filesystem means **deletion is unavailable for conversations stored there** — stated plainly, because D-10 leaves no alternative.

### 1.5 File-ID reuse

**R-08 — Reuse cannot manufacture a false `complete`.**

| Filesystem | Reuse behaviour | Label |
|---|---|---|
| NTFS | Freeing a record increments its 16-bit sequence, so a reused record gets a **different** ID; a stored ID stops matching at deletion and only matches again if the sequence wraps (2^16 reuses of that record) | **Source-derived** (MS-FSA) + **Calculated** |
| ReFS | IDs described as monotonically increasing and not reused for the volume's lifetime | **Source-derived** (community technical reference; vendor confirmation pending E-15) |
| FAT | ID can change without deletion | **Source-derived** — another reason FAT is refused |

**Direction of error (Calculated, important):** reuse makes an absent identity appear *present*. That yields an under-claim — `indeterminate`, `aborted`, or `identity_mismatch` — never a false `complete`. The only over-claim vector is the cloned volume of R-06, which is a different mechanism and is stated as a limitation.

**R-09 — Corroboration may only reduce confidence.** When a reopened identity's creation time is later than the journal's `updated_at`, the result is downgraded to `indeterminate`/`identity_reuse_suspected`. Corroborating evidence is never used to *strengthen* an identity claim.

---

## 2. The `absent_prior` contradiction

**Finding, accepted.** 1e's table row 1 emits `absent_prior`/`prior_absence` for an absent target at acquisition, while 1e §8 requires a demonstrated stable-parent basis and otherwise `absence_unverified`. Row 1 and §8 disagree; §8 is correct.

**R-10 — `absent_prior` is gated, and currently unreachable.**

> `absent_prior` may be emitted only when **all** of the following hold: the parent's physical identity was established from a retained handle; the observation was made relative to that parent; the child name passed validation; the parent identity was re-verified unchanged after the observation; and the observation basis is recorded as `handle_relative` or `pinned_parent` **and that basis has been demonstrated by recorded evidence**.

**R-11 — Until then, the outcome is `absence_unverified`.** A result code, not a receipt state: no receipt is written, no ownership is taken, no entry is created, nothing is deleted. The caller is told absence could not be established. This is fail-closed, not a degraded success.

**R-12 — `pinned_parent` is not proven by opening a handle.** Opening a parent without `FILE_SHARE_DELETE` gates DELETE access (Source-derived: share modes govern DELETE), but 1e already recorded a Q&A observation that a directory may be convertible to a reparse point with only `FILE_WRITE_ATTRIBUTES`. `pinned_parent` therefore remains **Unestablished** until E-3 is demonstrated; the contract must not treat a handle as proof of stability.

**R-13 — What each basis must demonstrate.**

| Element | `handle_relative` | `pinned_parent` |
|---|---|---|
| Parent identity | `FILE_ID_INFO` from the held parent handle | same |
| Parent stability | `RootDirectory` resolves the child under the held object | share-mode exclusion blocks rename **and** reparse conversion — E-3 |
| Child-name validity | R-14 | R-14 |
| Absence relative to parent | native relative open returns not-found | path probe while the prefix is pinned |
| Persistable / reconcilable | observation basis + parent identity recorded | same |

**R-14 — Child-name validation (carried, unchanged in substance).** Exactly one component; no separators; no `.`/`..`; no trailing dot or space; no reserved device name (`CON`, `PRN`, `AUX`, `NUL`, `COM1`–`COM9`, `LPT1`–`LPT9`, superscript variants). Failure is `invalid_child_name`, never "absent".

**R-15 — Vocabulary note.** `absent_prior` stays defined so that the vocabulary does not churn, but it is **unreachable in Phase 1**. A receipt carrying it while the gate is closed is invalid.

---

## 3. Crash windows against external filesystem change

**Finding, accepted.** 1e's crash table described "absent / present / partially removed" loosely and omitted the replacement case and the prepared/pending partial columns that the state table allows. The two tables did not agree.

**R-16 — Four target conditions, evaluated against the acquired identity.**

| Code | Condition |
|---|---|
| **A** | acquired identity not reopenable (absent) |
| **B** | acquired identity present, `entries == entries_at_acquisition` (unchanged) |
| **C** | acquired identity present, `entries < entries_at_acquisition` (partially changed) |
| **D** | condition A **and** the original path is now occupied by a different identity |

**R-17 — Attribution rule, restated because the matrix depends on it.** `attribution: operation` requires `state = removed`, identity verified absent, and a removal performed by this operation. **Every other cell is `unknown` or empty.** An observed partial directory is never evidence that this operation deleted part of it — including at `removing`, where the call may not have run.

### 3.1 Complete crash-window table

Ten windows × four conditions. `IE` = removal evidence. Reaching B/C requires reacquisition by identity (R-02); "path" is never consulted for the outcome, only to distinguish D.

| # | Crash point | Durable state | Removal invoked? | IE | A → outcome / attribution | B → outcome | C → outcome | D → outcome |
|---|---|---|---|---|---|---|---|---|
| 1 | before `prepared` | none | no | `not_started` | no operation exists; a later request acquires fresh | — | — | fresh acquisition would target the **replacement**, which is a different conversation — refused unless the caller names it |
| 2 | after `prepared`, before intent | `prepared`, no intent | no | `not_started` | `abandoned_no_intent`, entry removed, **no receipt** | same | same | same — no deletion was requested |
| 3 | after intent, before `pending` | `prepared` + intent | no | `not_started` | `absent_after_acquisition` / `unknown` | `aborted`, 0 | `partial` / `unknown` | `absent_after_acquisition` / `unknown` |
| 4 | after `pending`, before `removing` | `pending` | no | `not_started` | `absent_after_acquisition` / `unknown` | `aborted`, 0 | `partial` / `unknown` | `absent_after_acquisition` / `unknown` |
| 5 | after `removing`, before the call | `removing` | **unknown** | `possible` | `absent_unverified` / `unknown` | `aborted`, 0 | `partial` / `unknown` | `absent_unverified` / `unknown` |
| 6 | during removal | `removing` | yes | `possible` | `absent_unverified` / `unknown` | `aborted`, 0 | `partial` / `unknown` | `absent_unverified` / `unknown` |
| 7 | after return, before verification | `removing` | yes | `possible` | `absent_unverified` / `unknown` | `aborted`, 0 | `partial` / `unknown` | `absent_unverified` / `unknown` |
| 8 | after verification, before `removed` | `removing` | yes | `possible` | `absent_unverified` / `unknown` — a verification lost before the durable write is not evidence | `aborted`, 0 | `partial` / `unknown` | `absent_unverified` / `unknown` |
| 9 | after `removed`, before terminal | `removed` | yes | `verified` | `complete` / `operation` | `indeterminate` / `identity_present_after_verified_removal` | `indeterminate` / `state_inconsistent` | `complete` / `operation` — the replacement is untouched and irrelevant |
| 10 | after terminal, before cleanup | `removed` + terminal | yes | `verified` | receipt returned unchanged; entry removed | receipt returned unchanged; inconsistency **reported, never rewritten** | same | same |

**R-18 — Two readings the table forbids.** Absence never proves removal happened (windows 5–8, column A). A pre-removal state never proves it did not — only `prepared`/`pending`, whose ordering forbids the call (windows 2–4).

**R-19 — Column C at `prepared`/`pending` (windows 3–4).** A partly changed directory while no call could have run is an **external** change. `partial` is the honest effect report; the attribution is `unknown` and the receipt must not suggest otherwise.

---

## 4. "Verified absent", precisely

**R-20 — Four conditions the word "absent" can mean.**

| Code | Condition | Evidence |
|---|---|---|
| **L1** | identity still linked in the namespace | fresh `OpenFileById` succeeds and matches |
| **L2** | identity unlinked but reachable through a pre-existing handle | POSIX-semantics delete removes the link when the delete handle closes, while data stays reachable through other open handles (Source-derived: MS-FSCC) |
| **L3** | identity no longer reopenable through supported identity mechanisms | fresh `OpenFileById` fails after all operation handles are closed |
| **L4** | a different identity occupies the original path | path probe differing from the acquired identity |

**R-21 — `removed` requires L3, and only L3.**

> `removed` is written only when a **fresh** `OpenFileById` on the acquired identity fails after every deletion handle has been closed. A retained handle, an open file object, or a cached result is **not** evidence — under POSIX semantics the object can remain readable through handles that already exist while the namespace entry is gone.

**R-22 — `complete` claims exactly this:** the stored identity was not reopenable on the volume reached by the stored locator, at verification time, after this operation's removal. It does **not** claim: that the data is erased; that no process still holds it open; that no copy exists on a clone; that the path is empty (L4 may hold); or that the volume's other contents were unaffected beyond the enumerated tree.

**R-23 — Row 13 clarified (1e's "identity still resolvable").** It means a **fresh** open of the identity succeeding — a live namespace object — not this process still holding an earlier handle. An old handle proves L2, which is compatible with a successful deletion, so it must never be read as "still present".

**R-24 — Verification order.** Close all operation handles → fresh reacquisition of the volume (R-02) → fresh `OpenFileById` → only then write `removed`. Verification that reuses a handle opened before the delete is invalid and is a test failure.

---

## 5. Matrix corrections

**R-25 — Row 1 (absent at acquisition).** Gated by R-10/R-11: `absent_prior` only under a demonstrated basis; otherwise `absence_unverified` as a result code.

**R-26 — Row 12 (`removed`, identity absent).** `complete` / `operation`, **whether or not** the path is occupied by a replacement (condition D). 1e's row 12 was right; the replacement case is explicitly absorbed here so it cannot be mistaken for inconsistency.

**R-27 — Row 13 (`removed`, identity present).** `indeterminate` / `identity_present_after_verified_removal`. Candidate explanations, none of which the contract may choose between: the earlier verification was wrong (unsupported filesystem semantics), the ID was reused after sequence wrap, the volume was cloned, or the tree was restored from elsewhere. The receipt is **not** rewritten and the entry is **not** repaired; the condition is reported.

**R-28 — Row 14 (`removed`, identity present, `entries < baseline`).** Same class as R-27 with additional evidence that the directory lost entries. Outcome `indeterminate` / `state_inconsistent`. It cannot be `complete` (identity is present), and it cannot be `aborted` (that would assert nothing was removed while the tree is demonstrably smaller). Both assertions would be unsupported; `indeterminate` is the only honest cell.

**R-29 — Corrected matrix.** Replaces 1e's Q-22. Presence = presence of the **acquired identity**; "path" appears only to distinguish D.

| # | State | Intent | Identity | Path | IE | Verified absent | Outcome | Attribution |
|---|---|---|---|---|---|---|---|---|
| 1 | — | — | absent at acquisition | — | `not_started` | — | `absent_prior` **gated (R-10)**; otherwise `absence_unverified`, no receipt | `prior_absence` only when gated |
| 2 | `prepared` | absent | any | any | `not_started` | — | `abandoned_no_intent`, no receipt | — |
| 3 | `prepared` | present | absent | empty/replacement | `not_started` | no | `absent_after_acquisition` | `unknown` |
| 4 | `prepared` | present | present, intact | — | `not_started` | no | `aborted`, 0 | — |
| 5 | `prepared` | present | present, `entries <` | — | `not_started` | no | `partial` | `unknown` |
| 6 | `pending` | present | absent | empty/replacement | `not_started` | no | `absent_after_acquisition` | `unknown` |
| 7 | `pending` | present | present, intact | — | `not_started` | no | `aborted`, 0 | — |
| 8 | `pending` | present | present, `entries <` | — | `not_started` | no | `partial` | `unknown` |
| 9 | `removing` | present | absent | any | `possible` | no | `absent_unverified` | `unknown` |
| 10 | `removing` | present | present, intact | — | `possible` | no | `aborted`, 0 | — |
| 11 | `removing` | present | present, `entries <` | — | `possible` | no | `partial` | `unknown` |
| 12 | `removed` | present | absent | **replacement allowed** | `verified` | yes | `complete` | `operation` |
| 13 | `removed` | present | present (fresh open succeeds) | — | `verified` | claimed | `indeterminate` / `identity_present_after_verified_removal` | `unknown` |
| 14 | `removed` | present | present, `entries <` | — | `verified` | claimed | `indeterminate` / `state_inconsistent` | `unknown` |
| 15 | `pending`/`removing`/`removed` | absent | any | any | — | — | `indeterminate` / `receipt_missing`, entry preserved | `unknown` |
| 16 | any | any | any | any | — | — | terminal returned unchanged, then cleanup | — |

**Impossible combinations retained (R-30).** `prepared`+`possible`; `pending`+`verified`; `removed`+`not_started`; `removed`+`possible`; `attribution: operation` outside `state = removed`; a receipt with `removal_basis` other than `handle_identity`; `absent_prior` while the R-10 gate is closed. Each is `state_inconsistent`, reported, never repaired.

---

## 6. Owner decisions

**R-31 — 1e's decision 11 is withdrawn.** Framing deletion's availability as an owner choice contradicted D-10: identity-bound deletion was *decided*, and "enable it before the mechanism exists" is not an alternative — it would be D-10 not implemented. What remains is an engineering gate, not a decision.

> **Engineering gate G-1.** Deletion stays unavailable until the identity-bound mechanism is technically established and its acceptance evidence passes. This is a consequence of D-10, not a product choice to be made again.

**Remaining owner decisions (nine open, one decided):**

| # | Decision | Status |
|---|---|---|
| 1 | Retention Option A | **Decided** |
| 2 | Envelope with required predecessor | Open |
| 3 | Five ledger + two evidence kinds | Open |
| 4 | `FC1-02` linear | Open |
| 5 | `FC1-07` legacy read-only | Open |
| 6 | `FC1-14b` excluded / unadopted | Open |
| 7 | `FC1-14c` one message per turn | Open |
| 8 | Receipt without the five excluded items | Open |
| 9 | Per-conversation deletion proof | Open |
| 10 | Identity-bound deletion | **Decided** |

No decision asks the owner to choose a Windows API or to weaken identity-bound deletion.

---

## 7. Remaining feasibility gaps

Carried from 1e: **E-1** `CreateFileW` handle as `RootDirectory`; **E-2** `FILE_DIRECTORY_FILE`+`OPEN_REPARSE_POINT`; **E-3** share-mode exclusion vs rename **and** reparse conversion; **E-4** `OpenFileById`-failure semantics per filesystem; **E-5** mount-point deletion leaves the volume intact; **E-6** minimum build for class 71 and `FileDispositionInfoEx`; **E-7** privilege requirements; **E-8** POSIX delete with open child handles; **E-9** upcase-table semantics; **E-10** reparse-point ID identifies the link, not the target.

New in 1f: **E-11** GUID-path stability across reboot, letter change and removable reconnection, with the same volume; **E-12** minimal `hVolumeHint` access and whether a volume-root handle is accepted; **E-13** behaviour on SMB and other network filesystems; **E-14** FAT/exFAT refusal path (`FileId == 0` detection); **E-15** ReFS non-reuse over the volume lifetime, vendor-confirmed; **E-16** NTFS sequence-wrap behaviour on ID reuse; **E-17** whether a fresh `OpenFileById` fails immediately after POSIX delete while other handles remain open — the semantics R-21 depends on; **E-18** whether a volume GUID can be duplicated by a clone and how enumeration then behaves.

**Destructive tests required** (scratch volume, not a user's data): **DT-1** nested tree with file, subdirectory, file symlink, directory junction, mount point — link targets verified intact by hash; **DT-2** replacement race between acquisition and removal; **DT-3** rename while held; **DT-4** child created mid-deletion; **DT-5** crash injection at each of the ten windows; **DT-6** verification semantics matrix (NTFS/ReFS/FAT/SMB).

---

## 8. Acceptance tests

**Volume and identity recovery:**
- **RT-01** Restart between acquisition and reconciliation: volume reopened by GUID, serial+ID matched, identity resolved — no path used.
- **RT-02** GUID path absent → `volume_unavailable`; entry preserved; no scan performed.
- **RT-03** Serial or filesystem differs → `volume_identity_changed`; nothing deleted.
- **RT-04** Stored ID not found → `identity_absent`.
- **RT-05** ID opens but serial/ID differ, or object is not a directory → `identity_mismatch`.
- **RT-06** Two volumes answering the same GUID → `identity_ambiguous`, never a pick.
- **RT-07** FAT/exFAT or `FileId == 0` → `identity_unsupported`; deletion unavailable; nothing written.
- **RT-08** Creation time later than `updated_at` → `identity_reuse_suspected`, downgraded to `indeterminate`.

**Absent target:**
- **RT-09** With the R-10 gate closed, an absent target yields `absence_unverified` with no receipt and no entry.
- **RT-10** A receipt carrying `absent_prior` while the gate is closed is rejected as invalid.
- **RT-11** Parent renamed between observation and re-verification → `absence_unverified`.
- **RT-12** Invalid child name → `invalid_child_name`, distinct from absence.

**Crash × external change:**
- **RT-13** Windows 3 and 4, column C: a partly changed directory with no possible call is `partial`/`unknown`.
- **RT-14** Windows 5–8, column A: `absent_unverified`, never `complete`.
- **RT-15** Window 9, column D: `complete` with a replacement present; replacement byte-identical afterwards.
- **RT-16** Window 10, column B/C: receipt returned unchanged; inconsistency reported; no rewrite.
- **RT-17** No cell outside R-29's table is producible; every impossible combination is rejected.

**Verified absence:**
- **RT-18** A retained handle open across the delete does not prevent `removed`; verification uses a fresh open (E-17).
- **RT-19** Verification reusing a pre-delete handle is rejected as invalid.
- **RT-20** `complete` carries the R-22 claim and none of the stronger ones (string assertion on user-facing text).
- **RT-21** Row 13 with a fresh open succeeding → `indeterminate`, receipt unchanged.
- **RT-22** Row 14 → `indeterminate`/`state_inconsistent`; never `aborted`, never `complete`.

---

## 9. Verification and authorization

- Rows 528–534 remain byte-identical; row 534 (`0529728ef74f…`) is unchanged by this document.
- 1f is registered as row 535; corpus counts updated; no historical registry entry altered.
- **Zero Go files modified. No ledger created, no migration, no deletion enabled.**
- No test of any kind executed. No Windows behaviour is claimed as observed; API facts are source-derived, and every guarantee that depends on behaviour rather than documentation is labelled `Unestablished`.

TASK-03 remains in contract preparation. It is not accepted and not closed.
