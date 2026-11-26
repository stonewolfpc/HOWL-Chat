# TASK-01 — PathService: ownership inventory and completion record

| Field | Value |
|---|---|
| Task | `TASK-01` — PathService substrate (plan `HFI-EP-2026-11-22-01` §9) |
| Date | 2026-11-22 |
| Status | **Implemented and tested. Awaiting owner audit.** Not advanced to `TASK-02` |
| Authorization | Owner-authorized 2026-11-22 for **this task's scope only** — not for Phase 1, not for DurableStore, not for the backend rebuild |
| Scope guard | No existing call site, legacy function or storage location changed |
| Recovery token | **`TASK-01-PATHSERVICE-2026-11-22`** |

---

## 1. Safeguard A — ownership established before implementation

### 1.1 What actually resolves paths today

Three rules, all live, in one package and outside it.

| Rule | Definition | Behaviour |
|---|---|---|
| **R-1** `ProjectRoot()` | `internal/paths/paths.go:29` | Searches **upward from the executable**, then the working directory, for `wails.json`/`go.mod`, depth ≤ 6; cached process-wide |
| **R-2** `StateRoot()` | `internal/paths/isolation.go:185` | `Isolation.StateRoot` when isolated, otherwise **`ProjectRoot()`** |
| **R-3** relative literals | 7 `filepath.Join("…")` sites plus bare strings outside `internal/paths` | Resolved against the **working directory** |

Derived functions: `StatePath` (`isolation.go:193`), `InternalConfigDir` = `StateRoot()/config` (`content.go:74`),
`DefaultContentDir` (known-folder API, `content.go:25`), `TempDir`, `WebViewDataDir`, `IsolationActive`.

### 1.2 Call-site census (non-test, excluding `build/tmp` scratch tooling)

| Function | Call sites | What the callers actually do |
|---|---:|---|
| `ProjectRoot()` | 21 external + 1 internal | **Resource** lookups: `runtime/llama/llama-server.exe`, `sd-server.exe`, `sd-cli.exe`, `models/`, kokoro search roots, audio capability roots. No caller uses it to write state |
| `StateRoot()` | 9 | Storage: loader logs, RAG service, chat service (×2), memory service (×2), `InternalConfigDir` |
| `StatePath(` | 14 | Storage: characters, image gallery/entities, workflows, TTS cache, persisted chats, native-test token, startup logs |
| `InternalConfigDir()` | 2 | Settings directory, isolation probe |
| `DefaultContentDir(` | 3 | Library root creation/reporting |
| `TempDir()` | 12 | Scratch audio/image |
| `IsolationActive()` | 8 | Branch selection between isolated and production paths |

**The census changes the reading of the defect.** `ProjectRoot()`'s 21 external calls are *legitimate read-only
resource resolution*. The storage error is not that callers use the installation root directly; it is that
`StateRoot()` **defaults to it**. That is a single decision in one function, not 21 scattered bugs.

### 1.3 Disposition — retain, repair, wrap, replace

| Implementation | Disposition | Reason |
|---|---|---|
| `ProjectRoot()` | **Retain**, reclassified | Correct for R1/R4 resource resolution. Must never be treated as storage authority |
| `StateRoot()` | **Repair** (deferred to the migration task) | Its default is the defect. Flipping it without ADR-001 §12 S0–S4 would strand existing data |
| `StatePath()`, `InternalConfigDir()` | Retain | Already root-relative; they inherit the repair |
| `DefaultContentDir()` | **Retain**, now surfaced through `RootSet` | Known-folder resolution is already correct, including the OneDrive/redirected case |
| `containsPath()` | **Superseded for authorization**, retained for display | Text comparison cannot settle a junction. Authorization now uses filesystem identity |
| Relative literals (7 + bare strings) | **Replace** (deferred) | Working-directory storage. Listed in §7, not fixed here |
| External-root classification | **Absent entirely** | Nothing distinguished R4 from R3 before this task |

---

## 2. What was built

Additive only. Four new files in `internal/paths`; **no existing file modified**.

| File | Lines | Purpose |
|---|---:|---|
| `roots.go` | 385 | `Roots()`, `RootSet`, `Root`/`RootClass`/`Mode`, `CheckContained`, `ExternalRoot`, `CheckGrant`, typed `RootError` |
| `identity_windows.go` | 47 | Volume serial + file index via `GetFileInformationByHandle`; `FILE_FLAG_BACKUP_SEMANTICS` to open directories |
| `identity_other.go` | 32 | Device + inode for non-Windows builds |
| `roots_test.go` | 363 | 13 tests, listed in §4 |

`Roots()` returns R1 (read-only installation), R2 (`%AppData%\Howl` installed; isolation root isolated), R5
(temporary), and all seven content categories as R3. `Describe()` renders mode and roots deterministically.

---

## 3. Safeguard B — acceptance criteria and measured results

| # | Criterion | Test / evidence | Result |
|---|---|---|---|
| A-1 | Layout is independent of the working directory | `TestRootsAreWorkingDirectoryIndependent` — two different cwds, `Describe()` byte-identical | **PASS** |
| A-2 | No storage location derives from the executable's location | `TestRootsInstalledMode` — R2 is `UserConfigDir()/Howl`; R1 is `ProjectRoot()` and read-only | **PASS** |
| A-3 | Installed: R2 = `%AppData%\Howl`; isolated: R2/R3/R5 inside the isolation root | `TestRootsInstalledMode`, `TestRootsIsolatedMode` | **PASS** |
| A-4 | Category resolution uses the known-folder API | `content_windows.go` calls `windows.KnownFolderPath`; `TestRootsCoversEveryContentCategory` asserts absolute, classed, writable | **PASS** (API path) — **redirected-profile run not performed** |
| A-5 | Containment refuses `..` traversal, prefix-sharing siblings, relative candidates | `TestCheckContainedRejectsTraversalAndSiblings`, `TestCheckContainedRejectsRelativeAndEmpty` | **PASS** |
| A-6 | Containment decides by **identity**, not text | `TestFileIdentityIsStableAndDistinct`; `CheckContained` compares volume serial + file index when both endpoints exist | **PASS** (identity exercised) — **no real junction created** |
| A-7 | Resolution performs zero writes | `TestRootsPerformsNoWrites` — directory listing before/after three resolutions | **PASS** |
| A-8 | R1 is read-only and never a write target | `TestRootsInstalledMode` asserts `Writable: false` | **PASS** |
| A-9 | External paths are read-only and confer no permission | `TestExternalRootIsReadOnlyAndNeverGrants` — writable `false`, relative and missing roots refused | **PASS** |
| A-10 | Grant containment is strict and refuses the destination root itself | `TestCheckGrantRequiresStrictlyInside` | **PASS** |
| A-11 | Failure is typed, names the category/root and a remedy, and relocates nothing | `TestCheckContainedRejectsRelativeAndEmpty`, `TestExternalRootIsReadOnlyAndNeverGrants`, `TestRootSetCheckReportsUnknownCategory` | **PASS** |
| A-12 | **Compatibility with existing storage locations** | No existing `.go` file modified (`git diff --stat internal/paths/` empty); the four pre-existing isolation/content tests still pass; `go build ./...` passes | **PASS** |
| A-13 | Unknown category is an error, never a default | `TestRootSetCheckReportsUnknownCategory` | **PASS** |

**Verification actually run:** `gofmt -l` (clean), `go vet ./internal/paths/` (clean),
`go test ./internal/paths/` (all entries pass), `go build ./...` (passes).

**Known limitations, stated rather than implied**

1. A-4's redirected-profile case was not run; the known-folder call is exercised, redirection is not.
2. A-6's junction-escape case was not run — creating a junction needs privileges this task did not take.
   Identity comparison is exercised, junction refusal is not yet demonstrated.
3. A-1 is proven within one process by `chdir`. The literal "three separate launches" form of the criterion
   has not been performed.
4. Isolation setup (`newIsolation`) still creates directories during resolution; see I-02.

---

## 4. DurableStore — dependency made explicit

`PathService` **writes nothing**. It resolves, classifies and checks. It has **no compile-time dependency** on
DurableStore and none is created by this task.

| Relationship | Statement |
|---|---|
| PathService → DurableStore | None. PathService answers *where*; it never persists |
| DurableStore → PathService | **Required later.** Any store that writes must resolve its root through `Roots()` before opening a document |
| Boundary | `Roots()` being writable is **not** a durability guarantee. DurableStore owns atomicity, pre-images and the single-writer lock |
| Not built | DurableStore, in whole or part. It is `TASK-02` and is not authorized by this record |

---

## 5. What TASK-01 deliberately did not do

- **Did not change `StateRoot()`'s default.** The application still stores state in the installation directory.
  Flipping it without the migration sequence would orphan real data; that is the S0–S4 task.
- **Did not migrate or move any file.** No adoption records, no relocation, no quarantine.
- **Did not rewire any call site.** The 21 `ProjectRoot()` and 25 `StateRoot`/`StatePath`/`InternalConfigDir`
  callers are untouched; migrating them belongs with the migration task.
- **Did not touch pickers, navigation, characters, chats or any service.**
- **Did not implement DurableStore, grants enforcement, or containment on live write paths.**

`TASK-01` therefore has **no runtime effect on the shipped application**. That is the intended evidence and the
guard the owner asked for against a quiet storage rebuild.

---

## 6. Findings raised by this task

| # | Finding | Status |
|---|---|---|
| I-01 | `containsPath` is text-only and cannot refuse a junction that resolves outside the root | **Repaired in the new API**; legacy function retained for display |
| I-02 | `newIsolation` calls `MkdirAll` during resolution, though `PathService` is declared to have no writers | **Pre-existing, not changed** — altering it would change isolation test expectations; recorded as a follow-up |
| I-03 | `startup_log.go:24–27` writes `config/logs` **relative to the working directory** in normal mode, but `StateRoot/config/logs` under isolation | **Not fixed** — third-rule defect; belongs to the call-site migration |
| I-04 | `internal/audio/pipeline.go` `defaultTTSCacheDir` returns `"cache/tts"` cwd-relative when isolation is off | **Not fixed** — same class |
| I-05 | Whisper recognizer, `app-tts-catalog.go` and `voiceModelPath` build model paths relative to the working directory | **Not fixed** — same class |
| I-06 | `internal/backend/chat/service.go:148` names its variable `projectRoot` while assigning `paths.StateRoot()` | **Recorded** — naming evidence of the conflation the rebuild removes |
| I-07 | Nothing distinguished external (R4) data from user content (R3) before this task | **Addressed** — `ExternalRoot` classifies and withholds permission |
| I-08 | A writable root was previously indistinguishable from an authorized write | **Addressed** — `Writable` is a classification; `CheckGrant` is the authorization test |

---

## 7. Handoff — the five required fields

**Starting point.** `TASK-01-PATHSERVICE-2026-11-22`; plan §9; ADR-001 §4, contracts A and B. Specified and
already delivered: `internal/paths/roots.go`, `identity_windows.go`, `identity_other.go`, `roots_test.go`.

**Ownership.** PathService owns the root set, mode, category→location, class labels and containment checks.
Writers: none — resolution is pure. Readers: every store and service. It must **not** own write permission,
grants, policy, engine lifecycle, or any document's durability. Failure: typed error naming root/category and
remedy; never relocates.

**Acceptance criteria.** §3, A-1…A-13, all observable.

**Evidence.** `go test ./internal/paths/` (all pass, 13 tests added), `go vet` clean, `gofmt` clean,
`go build ./...` passes, `git diff --stat internal/paths/` empty. Limitations: no junction escape test, no
redirected-profile run, no three-process launch run.

**Next action.** Owner audit of this record. On approval: `TASK-02` — DurableStore, bounded to atomic replacement,
single-writer lock, pre-image retention and the four outcome states, with the reader-sharing question (`UN` in
ADR-001 §8.1) resolved first. Do not begin `1B` (ledger) until `DurableStore` passes.

---

## 8. Authorization

This record documents a completed, authorized task. It authorizes nothing further: not `TASK-02`, not the
migration sequence, not call-site rewiring, not Phase 1B, not any contract approval. FC-5, FC5-31 and R6 statuses
are untouched.
