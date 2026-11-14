# Engineering Report — F1 Lifecycle and Process Ownership

**Project:** Howl Chat  
**Phase:** Shared Foundations, F1  
**Date:** 2026-11-09  
**Baseline report:** `docs/review/ASSESSMENT-howl-application-reality-check-shared-foundations-2026-11-09.md`  
**Baseline SHA-256:** `fb0c05ced2619339e01ef712f850f96efbd756bfb5c21c347baec94aac2233c5`  
**Baseline length:** 554 lines (source/document inspection; not modified by F1)

**Status:** F1 source and unit-test work is complete. F1 runtime behavior is **not** runtime-confirmed. No claim in this report should be read as proof that a real Wails window close, a real model/image engine, or two simultaneous Howl instances were exercised.

---

## 1. Starting Working-Tree Condition

**Source-confirmed**

At the start of F1 the repository was at commit `9c7d2a1` on `master`, and the working tree was already dirty. It contained pre-existing modified tracked files and untracked files, including audio, UI, configuration, documentation, build/test artifacts, and review material. The tree was not clean, and no assumption was made that it was.

The final inspection still showed a dirty tree with **158 entries: 64 modified tracked files and 94 untracked files**. This count includes F1 changes and the pre-existing work. It is recorded as an observation, not as a claim that every entry belongs to F1.

No `git reset`, `git checkout`, `git clean`, `git restore`, stash, or other destructive Git operation was used. No unrelated file was removed to make the tree or tests appear clean.

The authoritative baseline report was left untouched. Its recorded length and SHA-256 above are the values supplied with the assignment and were not altered by this work.

---

## 2. Files Inspected

**Source-confirmed**

Primary files inspected during F1:

- `main.go`
- `app.go`
- `app-chat.go`
- `app-prompt.go`
- `app-model.go`
- `app-llama.go`
- `app-image-generate.go`
- `app-image-model.go`
- `app-audio.go`
- `app-lifecycle.go` (new)
- `app-lifecycle_test.go` (new)
- `app-lifecycle_windows_test.go` (new)
- `internal/proc/job_windows.go`
- `internal/proc/job_other.go`
- `internal/proc/process_windows.go` and related process helpers
- `internal/image/job/job.go`
- `internal/image/job/job_test.go`
- `internal/image/engine/sdcpp/engine.go`
- `internal/image/engine/sdcpp/video.go`
- `internal/image/engine/sdcpp/ownership_test.go` (existing shutdown/ownership coverage)
- `internal/audio/manager.go`
- `internal/audio` request/process handling relevant to shutdown
- `internal/errors/types.go`
- `frontend/dist/js/settings-state.js`
- `frontend/dist/index.html`
- `docs/review/ASSESSMENT-howl-application-reality-check-shared-foundations-2026-11-09.md`

Related callers and tests were inspected where needed to establish ownership, dependency order, and failure behavior. This was a focused F1 inspection, not a second application-wide audit.

---

## 3. Files Modified and Why

**Implemented changes**

| File | Change | Reason |
|---|---|---|
| `main.go` | Registered `OnShutdown: app.OnShutdown` | The Wails shutdown callback existed but was not connected; normal window close could bypass application cleanup. |
| `app-lifecycle.go` | Added the shutdown coordinator, work gate, stream ownership helpers, owned llama stop path, and stage failure isolation | One authoritative, idempotent lifecycle sequence was missing. |
| `app.go` | Added lifecycle/process/stream fields; removed the path-based orphan reaper call | Shared ownership state and removal of unsafe cross-instance cleanup. |
| `app-chat.go` | `ExitApp` now uses the shared shutdown; `SendMessage` checks the work gate; final-save errors are propagated | Explicit Exit and window close must not diverge; shutdown must stop new work; final-save failure must be visible. |
| `app-prompt.go` | Added shutdown gates and request-generation tokens around the existing single stream cancel function | Prevent new generations after shutdown starts and prevent an older generation from clearing a newer cancel handle. Full request-ID architecture remains F3. |
| `app-model.go` | Added lifecycle gates; made the llama monitor the single `Wait` owner; bounded the owned stop; rechecked shutdown after spawn/adopt/retry | Prevent a model load or retry from starting after shutdown began, and remove ambiguous `Wait` ownership. |
| `app-llama.go` | Uses `proc.Start`; waits through `processMonitor` | Process ownership must be established at creation, not after an unowned start. |
| `app-image-generate.go`, `app-image-model.go` | Added shutdown checks around image work and model loading | Image work could otherwise be admitted during application shutdown. |
| `app-audio.go` | Added shutdown gates to managed-audio and settings operations | New audio operations must not start after shutdown begins. |
| `internal/image/job/job.go` | `CancelAll` closes admission and cancels every non-terminal job; `StartImageJob` rejects closed admission | In-flight image jobs and model loads must not outlive the window. |
| `internal/image/job/job_test.go` | Added `CancelAll` coverage | Verify the new admission/cancellation behavior. |
| `internal/image/engine/sdcpp/engine.go` | sd-server starts through `proc.Start` | Remove the unowned-start/attach race. |
| `internal/image/engine/sdcpp/video.go` | Video starts through `proc.Start` | Use the same ownership policy for the video path. |
| `internal/proc/job_windows.go`, `internal/proc/job_other.go` | Removed path-based reaping and name-based tree killing; retained job ownership and explicit `KillTree` for an owned PID | Executable-path matching cannot prove ownership and could kill another Howl instance's child. |
| `internal/errors/types.go` | Added `SYSShutdownFail = "SYS-SHUTDOWN"` | Shutdown failures need a stable, visible error category. |
| `frontend/dist/js/settings-state.js` | `resetAllSettings` is asynchronous, backend-first, and restores prior values on failure | A reset must restore defaults through the authoritative settings write and must not claim success when persistence failed. |
| `frontend/dist/index.html` | Renamed the menu/button to “Reset Settings to Defaults”; confirmation text states what is preserved | The previous “Reset Application” wording implied a nuclear data wipe. |
| `app-lifecycle_test.go`, `app-lifecycle_windows_test.go` | Added lifecycle and real-child shutdown tests | Establish the claimed coordinator and owned-stop behavior. |

No new process manager, persistence layer, settings framework, or UI subsystem was introduced. The existing process manager and existing persistence mechanisms were retained.

---

## 4. Existing Defects Confirmed by Source Inspection

**Source-confirmed**

1. **The Wails shutdown callback was not registered.**  
   `App.OnShutdown` existed in `app.go`, but `main.go` did not pass it to the Wails application options. Normal window shutdown could therefore bypass the application-level cleanup path.

2. **There was no single authoritative shutdown sequence.**  
   Explicit exit and window close did not share one idempotent coordinator. Cleanup could be skipped, repeated, or ordered inconsistently.

3. **Path-based orphan reaping was unsafe across instances.**  
   The process layer contained path/name-based orphan cleanup. Matching an executable path does not establish ownership. A second Howl instance could own a process with the same executable path and be terminated by the first instance's cleanup.

4. **sd-server ownership was established after startup.**  
   The image server was started with `exec.Command` and attached to process management afterward. That creates a race: the child exists before ownership is recorded, and shutdown/crash handling can miss it.

5. **The video path used a different termination policy.**  
   Video-related execution did not follow the same creation-time ownership policy as the other managed processes.

6. **llama-server `Wait` ownership was ambiguous.**  
   The monitor and shutdown paths could both be understood as waiting for the same command. That is a structural double-`Wait`/lost-exit risk. A process command must have exactly one waiter.

7. **New application work was not gated during shutdown.**  
   Send, stream, model-load, image-load, and audio paths could begin after cleanup had started, racing teardown.

8. **Image jobs were not cancelled as a group during shutdown.**  
   The job manager could retain non-terminal jobs while the engine was being stopped.

9. **Audio shutdown was not part of the application lifecycle.**  
   The existing audio manager had a bounded close/drain operation, but the application-level shutdown sequence did not invoke it.

10. **Final-save failure could be lost.**  
    The final save path did not consistently surface its error to the lifecycle caller, so shutdown could appear successful while durable state failed to save.

11. **“Reset Application” implied a nuclear wipe.**  
    The UI wording described an application/data reset even though the intended operation is a settings reset. That is both a product contradiction and a risk of user data loss by expectation.

These are source-inspection findings. They are not claims that every defect was reproduced against a live application.

---

## 5. Defects Reproduced by Tests or Runtime Verification

**Test-confirmed**

No pre-fix failing test was run against the old implementation. The tests below were added with the repair and confirm the intended behavior after the change. They must not be described as a reproduction of the original failure.

The following behavior is **test-confirmed**:

- Eight concurrent shutdown callers execute the shutdown body exactly once and all receive the same result.
- A repeated shutdown returns the first stored error rather than running cleanup again.
- A panic in the shutdown body is converted into a returned error, and later callers observe that stored result instead of hanging.
- Shutdown with nil components is idempotent and closes the work gate.
- `main.go` contains the registered `OnShutdown: app.OnShutdown` callback (static source assertion).
- A real bounded child started for the test is stopped through `stopLlamaServer` and the monitor, without a second `Wait`; the child is no longer alive and the server reference is released.
- `CancelAll` closes image-job admission and cancels active jobs.

**Runtime-confirmed:** none.

**Blocked:** full Wails window-close verification, real model/image process verification, two-instance verification, and live settings-reset verification were not performed. See Sections 12 and 13.

---

## 6. Architectural Decisions and Rationale

**Implemented decisions**

1. **Keep the existing process manager and correct its ownership model.**  
   Replacing it would be a rewrite without a demonstrated need. The existing Windows Job Object is retained as the crash backstop.

2. **Establish ownership at creation on Windows.**  
   `proc.Start` starts the child suspended, assigns it to the job, then resumes it. This removes the start-then-attach race for llama, sd-server, and sd-video.

3. **Remove path/name-based orphan cleanup.**  
   Ownership is now based on the job/process handle and the recorded PID, not on an executable path. No broad process-name killing is used.

4. **Use one `Wait` owner for llama-server.**  
   The process monitor owns `Wait`; shutdown asks the monitor to stop and waits for its result. The normal path does not call `Wait` from a second goroutine.

5. **Centralize shutdown in one coordinator.**  
   `App.shutdown` is the single entry point for both Wails `OnShutdown` and explicit `ExitApp`. A `sync.Once`-style coordinator stores the result and makes repeated calls safe.

6. **Use a work gate rather than deleting entry points.**  
   `shuttingDown` is checked at user-work boundaries. Existing operations keep their behavior before shutdown and fail visibly after it.

7. **Cancel before save, then tear down in dependency order.**  
   New work is rejected; active stream/HTTP/image work is cancelled; durable state is saved; then owned model, image, audio, and job resources are released. This preserves the last usable chat state while preventing new work from racing teardown.

8. **Isolate stage failures and report all of them.**  
   Each stage is recovered independently and its error is accumulated. One failed component cannot skip the remaining cleanup. Errors carry the `SYS-SHUTDOWN` category and stage context.

9. **Bound the shutdown waits.**  
   llama shutdown has a 10-second bound; image shutdown uses its existing 5-second reaping bound; audio uses its existing bounded close/drain policy. A timeout is reported, not silently converted to success.

10. **Treat reset as settings-only.**  
    The UI performs a backend settings write first and only then updates localStorage. On failure it restores the previous in-memory settings, so a partial reset is not presented as successful. User content and model files are outside this operation.

11. **Do not redesign persistence or settings for F1.**  
    The existing chat persistence and settings mechanisms are used. Limitations that belong to F2 are documented rather than silently expanded into an F1 migration.

---

## 7. Process Ownership Model Before and After

**Source-confirmed**

### Before

| Process/resource | Creation/ownership | Weakness |
|---|---|---|
| llama-server | Existing `proc.Start` + Windows job; monitor existed | `Wait` ownership was ambiguous; shutdown could wait through a different path. |
| sd-server | `exec.Command` start, then process attachment | Child could exist before ownership was recorded; crash/shutdown race. |
| sd-video | Separate start/termination path | Inconsistent ownership and tree cleanup. |
| Orphan cleanup | Executable path/name matching | Could terminate another Howl instance's child. |
| Image jobs | Manager-owned jobs; no group cancellation | Jobs could remain non-terminal while the engine stopped. |
| Audio manager | Existing manager with close/drain support | Not called by one authoritative application shutdown. |
| TTS request processes | Existing request-scoped execution | Left under the frozen TTS boundary; not redesigned by F1. |

### After

| Process/resource | Ownership and cleanup policy |
|---|---|
| llama-server | Created through `proc.Start` (suspended → job assignment → resume). `processMonitor` is the sole `Wait` owner. Shutdown takes the recorded command, kills the owned tree/PID, waits through the monitor with a 10-second bound, and releases the reference. |
| sd-server | Created through `proc.Start`. Owned by the image engine; `Engine.Shutdown` cancel/kill/reset path reaps the recorded server process and resets engine state. |
| sd-video | Created through `proc.Start`. Owned by the video path; shutdown cancels video work and uses the owned-process cleanup path. |
| Image jobs | Manager-owned. `CancelAll` closes admission and marks/cancels every non-terminal job. `StartImageJob` refuses new work after admission closes. |
| Audio manager | Application-owned. `Close(ctx)` is invoked by shutdown with the existing bounded drain policy; errors are reported. |
| TTS/ASR request processes | Existing request-scoped context cancellation and wait behavior, unchanged. TTS authorization and contract remain frozen. |
| Orphan cleanup | Removed. No executable-path proof of ownership and no broad name-based kill. Crash cleanup relies on the Windows Job Object established at process start. |

**Cross-instance conclusion (source-confirmed, not runtime-tested):** the path-based reaper that could kill another instance's child has been removed. Cleanup is now scoped to the job/process handle and the recorded owned PID. A two-instance runtime test was not executed, so the absence of cross-instance interference is not runtime-confirmed.

---

## 8. Shutdown Sequence Before and After

### Before

**Source-confirmed**

- `main.go` registered only `OnStartup`; the `OnShutdown` method was not wired to Wails.
- Normal window close could therefore skip the application-level cleanup sequence.
- Explicit `ExitApp` did not share one authoritative idempotent lifecycle with window close.
- There was no unified work gate, no group cancellation of image jobs, no single llama wait owner, and no guaranteed application-level audio close.
- Path-based orphan reaping existed outside an ownership boundary.

### After

**Implemented; source- and test-confirmed at the component level**

The actual sequence in `app-lifecycle.go` is:

1. Set `shuttingDown` and reject new application work.
2. Cancel the active stream, cancel the HTTP request, and call `jobManager.CancelAll()`.
3. Attempt final chat persistence through the existing `performFinalSave` path; record any error.
4. Stop the owned llama-server, kill its tree/PID, and wait through the single monitor with a bounded timeout.
5. Shut down the image engine, including its owned server/video processes.
6. Close/drain the audio manager with its existing bounded policy.
7. Close the process job as the final backstop.
8. Return all accumulated stage failures with stage context; `OnShutdown` reports them visibly.

Each stage is independently panic-isolated. A stage failure does not prevent later stages from running. Final save is attempted before model/image teardown so durable chat state is not dependent on child-process shutdown succeeding.

The actual Wails callback is registered in `main.go` as `OnShutdown: app.OnShutdown`; explicit `ExitApp` calls the same `a.shutdown(nil)` entry point.

**Not runtime-confirmed:** the sequence has not been observed through a real Wails window close, and the real child processes have not been observed disappearing in a packaged application run.

---

## 9. Tests Added or Changed

**Test-confirmed**

Added:

- `app-lifecycle_test.go`
  - `TestShutdownCoordinatorRunsOnceUnderConcurrency`
  - `TestShutdownCoordinatorReturnsTheSameError`
  - `TestShutdownCoordinatorRecoversPanicAndLetsWaitersFinish`
  - `TestShutdownWithNilComponentsIsIdempotent`
  - `TestMainRegistersShutdownCallback`

- `app-lifecycle_windows_test.go` (Windows build tag)
  - `TestStopLlamaServerKillsAndWaitsThroughMonitor`

- `internal/image/job/job_test.go`
  - `TestCancelAllClosesAdmissionAndCancelsActiveJob`

The tests use bounded waits and a short-lived local child (`cmd /c ping`) for the owned-process test. They do not launch a real model, image engine, audio device, or second Howl instance. No unrelated process is killed by the tests.

---

## 10. Exact Test Commands and Results

**Test-confirmed** where marked PASS. Commands were executed during F1.

```text
go test -v -run 'TestShutdown|TestStopLlamaServer|TestMainRegisters' .
PASS
```

The regex `TestStopLlamaServer` matches `TestStopLlamaServerKillsAndWaitsThroughMonitor`.

```text
go test -v ./internal/proc
PASS
```

```text
go test -v -run TestCancelAllClosesAdmissionAndCancelsActiveJob ./internal/image/job
PASS
```

```text
go test . ./internal/proc ./internal/image/job ./internal/image/engine/sdcpp
PASS
```

```text
go test -race -run 'TestShutdown|TestStopLlama|TestMainRegisters' .
ok  	howl-chat	3.795s
```

```text
go test -race ./internal/image/job
ok  	howl-chat/internal/image/job	3.372s
```

```text
node --check frontend/dist/js/settings-state.js
PASS (no syntax error)
```

**Non-F1 whitespace check:**

```text
git diff --check
exit code 2
```

The reported whitespace errors were in pre-existing files outside the F1 change set, specifically `IMAGE_PAGE_CHANGE_LEDGER.md` and `frontend/wailsjs/go/models.ts`. This was not treated as an F1 test pass or failure; those files were not modified to hide the condition.

**Not run:** `go vet ./...`, the complete repository test suite, shell/UI integration tests, and a packaged Wails runtime test. Their absence is a limitation, not a pass.

---

## 11. Build Results

**Test-confirmed**

```text
go build ./...
PASS
```

This establishes compilation of the touched Go packages and their dependencies. It does not establish runtime lifecycle correctness.

`node --check frontend/dist/js/settings-state.js` also passed. It establishes JavaScript syntax only, not browser execution or backend persistence behavior.

---

## 12. Runtime Verification Results

**Unverified**

No application runtime verification was performed. Specifically, F1 did not exercise:

- a real Wails window close through the registered `OnShutdown` callback;
- a real llama-server launch and shutdown;
- a real sd-server or sd-video launch and shutdown;
- a real image generation cancellation;
- audio draining with a real output device;
- the settings reset through the live UI and backend;
- two simultaneous Howl instances and cross-instance process isolation;
- a crash/restart path that validates Job Object cleanup.

The Windows child test is a real bounded process test, but it is a unit/integration-style test of the owned stop path, not a packaged application runtime test.

Therefore this report does **not** claim that no Howl-owned child remains after a real shutdown. That requirement remains runtime-unverified.

---

## 13. Known Limitations

**Source-confirmed or inferred as noted**

1. **TTS request-scoped processes are not job-owned.**  
   Existing Piper/Kokoro request processes remain governed by the frozen TTS implementation. F1 did not redesign TTS, change its contract, or alter its authorization status. Their lifecycle is request-scoped and context-cancelled, but they are not part of the application Job Object policy in the same way as persistent model servers. This boundary must remain explicit.

2. **Image shutdown cancels jobs but does not prove every job goroutine has returned.**  
   `CancelAll` closes admission and asks the engine to cancel each job. `Engine.Shutdown` kills/reaps the server with a bounded wait and resets client state. The tests verify admission/cancellation behavior, not that every generation goroutine has completed before the process exits.

3. **Audio close is bounded, not unconditional.**  
   The existing manager may return a timeout/drain error. Shutdown reports it and continues. A real audio device was not tested.

4. **Settings reset is not a complete reset of every settings store.**  
   The reset writes the unified settings contract through the backend and restores the prior in-memory state on failure. Separate `config/image-settings.json` state and active image/audio engine state are not reset by this UI operation. This is a documented F2/F5 settings-ownership gap, not a silent claim of total reset.

5. **Settings storage-location defaults may change paths without moving data.**  
   The defaults include storage-location strings. Resetting them does not move existing user content, which is correct for a settings-only operation, but it may change which location the application will use on the next operation. The final data-location policy remains F2 work.

6. **`WipeAllData` / “Clear All Data” remains a separate explicit operation.**  
   It was not converted into the settings reset and was not implemented as the future comprehensive cleanup feature. Its presence is a product/UI review item because it is a nuclear operation distinct from “Reset Settings to Defaults.”

7. **Full request identity is not implemented.**  
   F1 added only a generation token around the existing single stream cancel handle. A complete request-ID architecture is F3.

8. **No two-instance runtime proof.**  
   The unsafe path reaper was removed and ownership is handle/PID scoped, but two live instances were not run together.

9. **No packaged shutdown proof.**  
   The Wails callback is registered and a static test checks the registration, but the callback was not observed in a real window-close run.

10. **The dirty tree was preserved, not normalized.**  
    Pre-existing changes remain. This is intentional and prevents F1 from destroying unrelated work.

---

## 14. Remaining Risks

**Inferred / Unverified**

- A real model or image engine may have a shutdown path that blocks longer than the configured bound or leaves a descendant not covered by the job. This is not runtime-tested.
- A real Wails close may invoke `OnShutdown` with a context whose cancellation interacts with audio close differently than the unit path. This is not runtime-tested.
- Image job cancellation marks jobs and asks the engine to cancel, but a job goroutine could still be finishing when the process exits.
- Request-scoped TTS helper processes remain outside the job ownership model; crash cleanup for those processes is not covered by F1.
- The settings reset may leave image-specific settings or engine state inconsistent with the unified settings file until F2/F5 addresses the ownership split.
- The final chat save is attempted before child teardown, but a save failure is only reported; there is no retry or alternate persistence target. F2 owns the persistence contract.
- The process-wide Job Object is the correct Windows crash backstop, but no crash/restart test proves its behavior in this build.
- The working tree contains extensive pre-existing changes; a later build or test failure may originate outside F1 and must not be attributed automatically to this work.

---

## 15. F2 Dependencies Discovered

**Source-confirmed / Inferred**

F1 intentionally did not move data or redesign persistence. F2 must resolve:

1. The authoritative distinction between installation files, user-created content, application settings/state, downloaded models/assets, and logs/temporary files.
2. The relationship between the unified settings file and separate image settings (and any other settings stores).
3. Whether the settings reset must also reset image-engine settings, and how to do so without deleting mixed user content.
4. Atomic or recoverable settings writes, including what “partially reset” means across multiple stores.
5. The final data-location policy and migration/approval boundary.
6. The product decision around the separate nuclear “Clear All Data” operation versus a future comprehensive cleanup feature.
7. Any persistence retry/repair behavior needed when final save fails during shutdown.

F1 does not make F2 unsafe to begin. It does make the settings-reset scope and the final-save failure policy dependencies that F2 must own.

---

## 16. Preservation of Unrelated User Changes

**Source-confirmed**

Unrelated work was preserved.

- No destructive Git command was used.
- No `git clean`, reset, checkout, restore, or stash operation was used.
- The pre-existing dirty tree was left dirty.
- The baseline assessment was not modified.
- No unrelated source or documentation file was rewritten to make tests pass.
- The only intentional UI text changes were the reset label and its confirmation copy, which are part of the assigned settings-reset clarification.

The final tree still contains pre-existing modified and untracked work in addition to F1 changes. That condition is expected, not evidence of a clean-room rewrite.

---

## 17. User Content and Downloaded Models

**Source-confirmed**

No user content or downloaded model was deleted by F1.

- No chat, character, world, scenario, lore, or memory data was deleted.
- No model file, image asset, audio asset, or downloaded model directory was deleted.
- The settings reset preserves user content and model files by design; it targets settings, not the data directories.
- No data migration or directory move was performed.
- No installation directory was removed.

This report does not claim that every storage path is correctly separated; that is the F2 contract. It claims only that F1 did not delete user content or downloaded models.

---

## 18. F1 Completion Status

**Source/test status: complete for the implemented F1 scope.**  
**Runtime status: unverified.**

The completion gate is met at the source and unit-test level:

- The actual Wails shutdown callback is registered: **source-confirmed and statically test-confirmed**.
- Shutdown is single-entry, idempotent, dependency-aware, panic-isolated, and bounded: **source-confirmed and test-confirmed** for the coordinator behavior.
- Persistent child processes have an explicit creation-time ownership and cleanup policy: **source-confirmed** for llama, sd-server, and sd-video; TTS request processes remain under the frozen TTS boundary.
- Path-based cross-instance termination was removed: **source-confirmed**; two-instance runtime behavior is **unverified**.
- Partial initialization and shutdown failures are handled: **test-confirmed** for nil components, repeated calls, concurrency, and panic recovery.
- Final persistence is attempted through the existing mechanism and failures are propagated: **source-confirmed**; live failure behavior is **unverified**.
- Audio shutdown is addressed without changing TTS authorization: **source-confirmed**; real-device drain is **unverified**.
- Relevant tests passed: **test-confirmed** (Section 10).
- Unrelated work was preserved: **source-confirmed**.
- No user content or downloaded models were deleted: **source-confirmed**.

F1 is therefore **engineering-complete but not runtime-confirmed**. The remaining verification is not a reason to repeat the implementation or to begin a broad redesign; it is a bounded acceptance step: run the packaged application, close the window, observe the registered shutdown path, verify the owned child processes are gone, run two instances, and exercise the settings reset.

**F2 can begin.** Persistence/data-location work does not depend on F1's coordinator internals, but F2 must carry forward the settings-store split, final-save failure policy, and data-location limitations documented above. F1 should not be reported as runtime-proven until the bounded acceptance step in Section 12 is performed.
