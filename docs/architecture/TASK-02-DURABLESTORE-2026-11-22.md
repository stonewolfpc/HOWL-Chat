# TASK-02 — DurableStore

| Field | Value |
|---|---|
| Scope | Implement the `DurableStore` primitive from ADR-001 §8. **Adoption into call sites is not in scope** |
| Authorized | 2026-11-22 (owner: TASK-01 closure and TASK-02 start approved) |
| Decision of record | `TASK-02-READER-SHARING-DECISION-2026-11-22.md` — resolves ADR-001 §8.1 (`UN`) |
| Status | **Complete for the primitive contract.** Not adopted by any call site; adoption is a separate unauthorized step |
| Recovery token | `TASK-02-DURABLESTORE-2026-11-22` |

---

## 1. Step 1 — ADR-001 §8.1 resolved

§8.1 asked which readers may coexist with a writer and by what mechanism, and labelled the answer `UN`.
Resolved as **short-lived read + bounded writer retry**, on measured evidence rather than convenience:

| Operation, with a reader holding the live document via `os.Open` | Measured result |
|---|---|
| `os.Rename(next, live)` | **`Access is denied`** |
| `os.WriteFile(live, …)` | **succeeds** |

Go's `os.Open` on Windows shares read and write but not delete, so replacement is blocked by a held reader — and
direct writes are observably torn. The reader contract (open, read, close within one call) plus a 3 × 50 ms rename
retry inside the lock is the ADR's first listed protocol, needs no `UN` capability, and works with every existing
reader. Delete-sharing readers were rejected as requiring `NtCreateFile` (`UN`) and a rewrite of every call site;
shared-lock readers as turning reads into failures; snapshot pointers as creating a second "which document is
live" authority.

**No ADR amendment is required.** The decision supplies a value for an `UN` field without altering any normative
statement. The exact wording proposed if the owner later folds it in is in §6 of the decision record.

## 2. Step 2 — Inventory

### 2.1 Verified implementation facts (read from source, this tree)

| Location | Behaviour |
|---|---|
| `internal/settings/store.go:47–61` | `Save` writes with a **direct `os.WriteFile`** — truncate in place, no temp, no flush, no rename, no backup. `Load` is a single `os.ReadFile` |
| `internal/settings/store.go:47–56` | The comment justifies this as avoiding "temp-file rename dances that cause Windows file-locking errors". **Falsified by §1:** rename fails because a reader holds the file, not because rename is fragile |
| `internal/backend/chat/persistence.go:571–583` | `saveJSON` writes `<file>.tmp` then `os.Rename` — atomic, but **no `Sync`**, a fixed temp name (collides between writers), no lock, no pre-image, no manifest |
| `internal/image/gallery/gallery.go:118–155` | `writeFileAtomic` — temp + `Sync` + rename, mode-preserving, cleans up on error. **No lock, no pre-image, no checksum, no schema check, no reader coordination** |
| `internal/contentimport/journal.go:296` | Journal entries created with `O_EXCL`, i.e. append-only uniqueness rather than replacement |
| Whole repository | **No `LockFileEx`, no `flock`, no file-based lock anywhere.** The only mutex is an in-process `TryLock` (`internal/image/engine/sdcpp/video.go:154`) |
| Whole repository | 20+ raw `os.WriteFile` sites, several of them on settings/config paths |

### 2.2 Documentation claims vs verified facts

| Claim | Status |
|---|---|
| ADR-001 §8: `os.WriteFile` is "a direct, non-atomic write — no temp, no flush, no rename, no backup" | **Verified** (`settings/store.go:60`) |
| ADR-001 §8: `DurableStore` "must replace this" | Requirement; **not yet true of any call site** |
| `settings` comment: temp+rename "causes Windows file-locking errors" | **Falsified.** The block is a held reader; §1 measures both directions |
| ADR-001 §8: single-writer enforcement | **Unestablished in the tree before this task.** No lock mechanism existed |

### 2.3 Competing implementations and authorities

Three replacement strategies now exist in the tree: `gallery.writeFileAtomic`, `chat.saveJSON`, and this store.
That is a defect the rebuild exists to remove, and **TASK-02 did not remove it** — removing it means migrating call
sites, which this task is explicitly not authorized to do. The primitive is therefore **not yet authoritative**.
Recorded as deferred finding **D-1**.

## 3. Step 3 — Acceptance criteria and their outcomes

| # | Criterion | Test | Outcome |
|---|---|---|---|
| A-1 | Single-writer enforcement across processes | `TestBusyAcrossProcesses` — a second process holds the lock, this one attempts a write | **PASS** — `Busy`, and the same write succeeds after release |
| A-2 | Selected reader-sharing behaviour | `TestSaveReportsRenameDeniedWhenReaderHoldsLive` | **PASS** — `NotSaved(rename_denied_by_open_reader)`, previous document intact |
| A-3 | `Busy` with no bypass | `TestSaveReportsBusyWhenLockHeld`, `TestBusyAcrossProcesses` | **PASS** — document unchanged after `Busy` |
| A-4 | Atomic replacement | `TestReaderNeverObservesPartialDocument` — 20 commits under continuous concurrent reads | **PASS** — no read observed a partial or unknown document |
| A-5 | Backup / rotation | `TestSaveRotatesRetainedPreImages` — 4 commits, retention 2 | **PASS** — `v1` = 3rd, `v2` = 2nd, `v3` absent |
| A-6 | Process-crash durability | `TestCrashBeforeReplaceLeavesPreviousDocument` — child process exits between manifest commit and rename | **PASS** — previous document readable; reported `Recovered(interrupted_commit_rolled_back)` |
| A-7 | Recovery | `TestLoadRecoversFromPreImageWhenLiveCorrupt`; leftover-temp sweep in the crash test | **PASS** — recovers the newest hash-matching pre-image and names it |
| A-8 | Error propagation / fail-closed | `TestLoadRefusesNewerSchema`, `TestSaveRefusesNewerSchema` | **PASS** — refused, and the document is not rewritten |
| A-9 | Round-trip correctness | `TestSaveThenLoad`, `TestLoadAbsentWhenNoDocument` | **PASS** |
| A-10 | **Power-loss durability** | — | **`UN`. Not claimed, not tested** |

## 4. Step 4 — Implementation

New package `internal/durable`, four files, **no existing file modified**:

| File | Lines | Content |
|---|---:|---|
| `store.go` | 606 | `Outcome`, `Result`, `Options`, `Store`, `Save`, `Load`, manifest, pre-image rotation, temp sweep, bounded replacement retry |
| `lock_windows.go` | 62 | `LockFileEx` exclusive range lock, `Busy` on timeout, released on process exit |
| `lock_other.go` | 45 | advisory `flock` fallback so the package builds and tests off Windows |
| `store_test.go` | 430 | 11 tests plus 2 process-helper entry points |

Shape of a commit: lock → sweep dead temps and staged files → write temp → flush → **stage** pre-images →
commit manifest → replace by rename (3 × 50 ms) → **promote** staged generations → drop pre-images beyond N →
release. **The live path is never opened for writing**; its only mutation is a rename.

> **Corrected by the TASK-02 audit (`TASK-02-AUDIT-2026-11-22.md`).** This section originally described rotating
> pre-images directly, before the rename. That ordering was defective: see F-A2 there. Pre-images are now staged and
> promoted only after a successful replacement.

Three defects were found by these tests during implementation and fixed rather than relaxed:

- **Manifest replacement had the same reader-denial problem** as the document. It now uses the same bounded retry
  and reports `rename_denied_by_open_reader` rather than a generic error.
- **A reader's document+manifest pair could straddle a commit**, which would look like corruption. `Load` now
  re-reads the pair a bounded number of times before any recovery, and distinguishes a rollback (live matches a
  retained generation) from corruption (live matches nothing) — `interrupted_commit_rolled_back` versus
  `live_hash_mismatch`.

A third defect, **F-A2**, was found by the later audit rather than here: rotation overwrote `v1` before the rename, so
**a failed commit destroyed the newest retained generation** and left the restored manifest describing files that no
longer held those bytes. Recovery after a routine reader-contention failure reported `corrupt_no_valid_preimage`.
Corrected by staging and promotion; `TestFailedCommitPreservesPreImageChain` reproduces it and
`TestFailedCommitLeavesNoStagedArtifacts` pins the cleanup. This is why the ordering above now reads stage → commit →
rename → promote.

## 5. Step 5 — Verification

| Command | Result |
|---|---|
| `gofmt -l internal/durable/` | no output |
| `go vet ./internal/durable/` | clean |
| `go test ./internal/durable/ -race -count=2 -v` | **13 PASS, 2 SKIP (helpers), exit 0** |
| `go test ./internal/durable/ -count=3` | all pass |
| `go build ./...` | passes |
| `git status` on `internal/durable` | only new untracked files |

### 5.1 What each test demonstrates

| Test | Demonstrates |
|---|---|
| `TestSaveThenLoad` | Commit and read back; manifest written |
| `TestLoadAbsentWhenNoDocument` | Absence is a distinct, non-error outcome |
| `TestSaveRotatesRetainedPreImages` | Rotation order and the retention bound |
| `TestReaderNeverObservesPartialDocument` | The atomicity claim under real concurrency; also that writer-failure under reader contention is reported, not retried forever |
| `TestSaveReportsBusyWhenLockHeld` | In-process contention produces `Busy` and no write |
| `TestBusyAcrossProcesses` | **Cross-process** mutual exclusion — the actual single-writer requirement |
| `TestSaveReportsRenameDeniedWhenReaderHoldsLive` | The reader protocol's consequence, and that a failed commit does not leave a misleading manifest |
| `TestLoadRecoversFromPreImageWhenLiveCorrupt` | Recovery selects a hash-matching generation |
| `TestLoadRefusesNewerSchema` / `TestSaveRefusesNewerSchema` | Unknown newer schema is refused and never rewritten |
| `TestCrashBeforeReplaceLeavesPreviousDocument` | A real process death mid-commit leaves the previous document intact and the outcome honest |
| `TestFailedCommitPreservesPreImageChain` | A failed commit leaves the retained chain untouched and recovery still works (F-A2) |
| `TestFailedCommitLeavesNoStagedArtifacts` | A failed commit leaves no prepared generation behind |

### 5.2 Limitations and untested conditions

- **Power-loss durability is `UN`.** Flush ordering exists; NTFS metadata durability after power loss is not
  demonstrated and is not claimed.
- **Single-machine, single-filesystem only.** No network filesystem was exercised; SMB semantics for
  `LockFileEx` and rename-over are not established here.
- **Reader long-holds are not enforced**, only surfaced as `rename_denied`. Enforcement is a review obligation
  (decision record §5).
- **The crash test kills between manifest commit and rename.** A crash between temp write and manifest, or during
  pre-image rotation, is not separately injected; those points leave the live document untouched by construction.
- **Helper tests skip in normal runs** (`TestCrashChildDurable`, `TestLockHolderChild`); they only execute as
  driven child processes.
- **No call site uses this store yet**, so no application behaviour changed and none was regression-tested.

## 6. Step 6 — Status

**TASK-02 is complete for its bounded objective:** the primitive exists, is exercised by 13 passing tests including
cross-process exclusion, real process-death recovery and failed-commit chain preservation, and has been audited
(`TASK-02-AUDIT-2026-11-22.md`) — one blocking defect (F-A2) was found there and corrected. **Acceptance remains the
owner's decision; this record does not accept itself.** It is **not complete as an
application capability**, because nothing uses it. That distinction is deliberate: adopting it means migrating
call sites, which this task was not authorized to do.

### Deliberately deferred

- **D-1** `gallery.writeFileAtomic`, `chat.saveJSON` and the direct `settings.Save` remain in place; three
  replacement strategies coexist until migration.
- **D-2** The 20+ raw `os.WriteFile` sites are untouched.
- **D-3** `settings.Save`'s falsified rationale comment remains; it is fixed by adoption, not by editing the
  comment in isolation.
- **D-4** Power-loss durability remains `UN` pending independent demonstration.

### Next step

Adoption begins with `internal/settings` (smallest document, one writer, one reader, already identified as the
non-atomic write in ADR-001 §8) — **not started, not authorized by this task**.

## 7. Authorization

Implements one primitive and documents its decision. Authorizes nothing further: not call-site migration, not
ConfigAuthority, not ModelConfigService, not generalized persistence migration, not any contract approval.

---

## 8. Formal closure — TASK-02 (appended 2026-11-22, after owner acceptance)

**Owner decision:** TASK-02 is **accepted as complete for its bounded objective**, with the limitations in §8.3
carried forward. This section is a closure record appended to the original; **nothing above it was rewritten as
though F-A2 had never happened.**

### 8.1 F-A2 is preserved as a corrected blocking defect

The defect stays visible in three places, and none of them was edited to remove it:

| Where | What it preserves |
|---|---|
| `TASK-02-AUDIT-2026-11-22.md` §2.2 | The full analysis: rotation overwrote `v1` before the rename, a failed commit desynchronised the chain from the restored manifest, and recovery then reported `corrupt_no_valid_preimage` after a routine reader-contention event |
| `HEKI-CP-2026-11-22-09` | The pre-audit claims, **including the one F-A2 falsified** ("a failed commit never rotates"), exactly as stated when the implementation was reported |
| §4 above | The original ordering sentence **retained in place**, with a visible correction note immediately beneath it saying what it used to say and why it was wrong |

**Reproduction evidence.** The failing test that demonstrated the defect, run against the delivered implementation:

```
audit_test.go:44: after failed commit, v1 = "d3", want the untouched d2
```

**Regression evidence.** After the correction, with pre-images staged as `.staged` and promoted only after a
successful replacement:

| Test | What it prevents from returning |
|---|---|
| `TestFailedCommitPreservesPreImageChain` | The reproduction itself: a failed commit must leave `v1`/`v2` untouched, `Load` must still return the live document, and recovery after later corruption must yield a real generation |
| `TestFailedCommitLeavesNoStagedArtifacts` | A failed commit leaving prepared generations behind, and a subsequent write being confused by them |

The defect was found by the audit, not by the implementation report. That sequence is the one that has to survive
into the record.

### 8.2 The reader-consistency condition is met by construction

Acceptance was conditional on `Recovered` being observable and not mistakable for an ordinary current read. The
result type now says so explicitly:

- `Outcome` is a distinct value (`Recovered`, never `Saved`), and `Reason`/`From` explain it.
- `OK()` means *usable bytes* — it includes `Recovered`.
- **`Current()`** means *the current committed document* — it is true only for `Saved`, so a caller that must not act
  on a superseded document has a predicate to use instead of `OK()`.
- `TestRecoveredIsDistinguishableFromCurrentRead` pins it: a recovered generation reports `OK()=true` and
  `Current()=false` and names its source.

This is an addition made **after** the audit, in response to the acceptance condition, which is why the audit record
still reports 13 tests and this record reports 14. The audit record is history and was not updated to match.

### 8.3 Limitations carried forward (accepted, not resolved)

| # | Limitation |
|---|---|
| 1 | **Power-loss durability is `UN`.** Flush ordering exists; durability of the rename itself across a power cut is not demonstrated and is not claimed |
| 2 | **Five interruption points have reasoned coverage, not crash-test coverage:** after temp write, after staging, after rename before promote, during promote, during trim. Only the point between manifest commit and rename is injected with a real process death |
| 3 | **Reader contention can return a previous, complete generation marked `Recovered`** — one generation behind under sustained commits, never torn |
| 4 | **Existing persistence implementations were not migrated.** `settings.Save`, `chat.saveJSON` and `gallery.writeFileAtomic` remain; three replacement strategies still coexist |
| 5 | **DurableStore is a completed primitive, not the application's persistence authority.** No call site uses it, so no application behaviour changed and none was regression-tested |

### 8.4 Verification at closure

`gofmt` clean · `go vet` clean · `go test ./internal/durable/ -race -count=2 -v` → **14 PASS, 2 helper SKIP** ·
`go build ./...` passes · **no pre-existing Go file modified**. Registry hashes verified after the final edits; the
actual results are recorded in `HEKI-CP-2026-11-22-11`.

### 8.5 Branch status

`TASK-01` closed. `TASK-02` closed. **This branch is complete.**

**Adoption is not authorized and was not started.** It is a separate task requiring its own inventory, scope,
acceptance criteria and authorization. The next branch is selected from
`HOWL-FOUNDATIONAL-IMPLEMENTATION-EXECUTION-PLAN-2026-11-22.md` **after** this closure checkpoint — not chosen here.
