# D3 — ASR process ownership and clean shutdown: lifecycle map

Date: 2026-11-05. Scope: local ASR `llama-server` child only. TTS, chat model,
and Unit C are out of scope.

Status: **ACCEPTED / CLOSED — OWNER DISPOSITION (2026-11-05).** Superseded status,
retained for the record: implemented and verified at process level, pending owner
review. The acceptance record is §7; the protected-data deviation it accepts is
recorded in §6.6 and is not softened anywhere in this document.

## 1. Ownership today

| Stage | Code | Ownership |
|---|---|---|
| Launch | `internal/audio/asr/llama/llama.go` `Initialize` → `exec.Command(serverPath, args...)`, `cmd.Start()` | **None.** No job attachment; no tree handle. |
| Health checks | `healthy()` polls `http://127.0.0.1:<port>/health` until 200 or 120 s | Process held only by `*exec.Cmd` in the `Recognizer`. |
| Wait/reap | goroutine `waitCh <- cmd.Wait()` | Reaps only if the server exits on its own. |
| Normal release | `Recognizer.Release` → `stopLocked`: `Process.Kill()`, wait ≤5 s | Kills only the direct child, no descendants, no job membership. |
| Manager close | `Manager.Close` retires the slot, starts `runShutdown` (goroutine), drains ≤5 s (`defaultCloseTimeout`), `buildCancel`/`ops.cancelAll`/`interruptAll`; returns `ErrDrainTimeout` on expiry | Release happens only if `runShutdown` reaches `Pipeline.Shutdown` → `loader.Unload` → `Release`. |
| App shutdown | `App.OnShutdown` (`app.go:254`): final save → chat server `proc.KillTree` → image engine → `audioManager.Close(ctx)` → `proc.CloseJob()` | `proc.CloseJob()` kills the job tree — which contains the chat server, but **not** the ASR server. |
| Crash/force-kill | none | `KILL_ON_JOB_CLOSE` covers job members only; ASR is not a member. |

`internal/proc/job_windows.go` is the established mechanism: one job with
`JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`, `Attach(cmd)` right after `Start`, tree kill
via `taskkill /T /F`, and `CloseJob()` at shutdown. The chat `llama-server`
(`app-llama.go`) uses it; the ASR `llama-server` does not.

## 2. Paths that can leave the child alive

- **F1 — hard exit / crash before `OnShutdown` completes (confirmed failure).**
  The ASR child is not in the job, so nothing terminates it when the app process
  disappears. This is the reported orphan.
- **F2 — `Manager.Close` drain timeout.**
  `Close` returns `ErrDrainTimeout` after 5 s; `OnShutdown` logs it, then calls
  `proc.CloseJob()`. Because the child is unattached, the job close does not
  reach it. If `runShutdown` was never scheduled or `loader.Unload` fails, the
  process survives.
- **F3 — release kills the direct child only.**
  `Process.Kill()` leaves any descendant of `llama-server` running, and a failed
  kill is not observed (wait result ignored, `stopLocked` always returns nil).
- **F4 — unexpected server exit with live descendants.**
  the `waitCh` branch in `Initialize` clears the handle without a tree kill.
- Not a path: `Pipeline.Shutdown` does reach `Release` on the normal path, and
  `initializeProvider` releases on initialization failure, so ordinary success
  and ordinary failure already call `Release`; the defect is that `Release`
  is neither crash-safe nor tree-complete.

Whisper (`whisper-cli`) is a short-lived per-request child, started and waited
synchronously inside `Transcribe`; it is not the long-lived server in the
reported failure and is not changed here.

## 3. Fix

1. `llama.Initialize`: `proc.Attach(cmd)` immediately after `Start` (same
   established pattern as the chat server), so `KILL_ON_JOB_CLOSE` owns it for
   the process lifetime and covers F1/F2.
2. `stopLocked`: `proc.KillTree(pid)` before `Process.Kill()` (F3), and report a
   bounded-exit failure instead of always returning nil.
3. Unexpected-exit branch: tree-kill the pid before clearing handles (F4).
4. No manager, app, timeout, or chat-model change.

## 4. Verification plan

Real child processes, real PIDs, OS-level exit checks
(`OpenProcess`/`WaitForSingleObject`), no mocked cleanup:

1. Normal `Initialize` → `Release`: server **and** spawned descendant exit.
2. Cancelled initialization after `Start`: server and descendant exit.
3. `Release` while a transcription request is in flight: no orphan; repeated
   `Release` safe.
4. Abrupt owning-process exit (helper subprocess calls `os.Exit` after
   `proc.Attach`): job kill removes server and descendant.
5. `proc.CloseJob()` while the owning process is still alive: server and
   descendant exit.

## 5. Implementation and results (2026-11-05)

### 5.1 A second, blocking defect in the ownership primitive

The first test run failed the two job tests. A direct Windows API probe
(`build/tmp/jobprobe`, deleted after use; Windows 10.0 build 26200) showed:

```
pid-as-handle assign=The handle is invalid.
pid-as-handle alive after job close: true
real-handle   assign=<nil>
real-handle   alive after job close: false
```

`proc.Attach` called `AssignProcessToJobObject(job, windows.Handle(pid))` — a
PID is not a process handle, so **every attachment failed and the error was
discarded**. The job was therefore empty for the chat `llama-server` as well;
`KILL_ON_JOB_CLOSE` never protected either server. This is the mechanism D3
must integrate with, so it was fixed rather than bypassed: `Attach` now opens
the process with `PROCESS_SET_QUOTA|PROCESS_TERMINATE`, assigns that handle,
closes it, and writes a diagnostic to stderr on failure instead of hiding it.

### 5.2 Changes

- `internal/proc/job_windows.go` — `Attach` uses a real process handle;
  assignment failures are reported, not swallowed.
- `internal/audio/asr/llama/llama.go` — `proc.Attach(cmd)` immediately after
  `Start`; `stopLocked` tree-kills (`proc.KillTree`) before `Process.Kill` and
  returns an error when the child is still alive after the 5 s bound; the
  unexpected-exit branch tree-kills too.
- `internal/audio/asr/llama/lifecycle_windows_test.go` — new real-process tests.

No manager, app, timeout, chat-model, or frontend change. The chat server's
ownership behaviour changes only in that its job attachment now actually
succeeds (the explicit `KillTree` at shutdown was already there).

### 5.3 Test results

`go test ./internal/audio/asr/llama/ -count=1 -v` → **PASS** (16.2 s),
exit 0:

| Test | Result | Evidence |
|---|---|---|
| `TestReleaseTerminatesServerAndDescendant` | PASS | server + spawned descendant gone after `Release` |
| `TestCancelledStartupTerminatesStartedServer` | PASS | startup cancelled after `Start`; server gone |
| `TestReleaseDuringActiveWorkAndRepeatedRelease` | PASS | in-flight `Transcribe` completes, server gone, second `Release` clean |
| `TestAbruptOwnerExitKillsServer` | PASS | owner subprocess `os.Exit(0)`; job kills server + descendant |
| `TestCloseJobKillsAttachedServer` | PASS | `CloseJob` while owner alive; server + descendant gone |
| `TestAbruptExitHelper` / `TestCloseJobHelper` | SKIP by design | not independent tests; launched as subprocesses by the two job tests |
| 3 existing transcript tests | PASS | unchanged |

Exit evidence is OS-level: `OpenProcess(SYNCHRONIZE)` + `WaitForSingleObject`
on the recorded PIDs (polled to a deadline), never a mocked cleanup call or a
return value. The test server is a purpose-built child that also spawns a
long-lived grandchild, so tree termination is observable.

`gofmt -l` clean; `go build ./...` exit 0.

### 5.4 Remaining limitations

- The **packaged application was not launched**; this is verified at process
  level, not through a packaged shutdown.
- The chat `llama-server` path now gets a working job attachment, but that was
  not exercised at runtime here; it is covered by the same primitive probe and
  by compilation only.
- Windows-only. `job_other.go` is still a no-op, so the guarantee does not
  extend to macOS/Linux.
- A crash in the microseconds between `Start` and `Attach` can still orphan the
  child (no suspended-start handshake). Narrow, but real.
- The manager's 5 s drain timeout still exists; the job now covers the case
  where that timeout would otherwise leave the server behind.

Observed after the fixed run: two `fakeasr.exe` processes (PIDs 296, 6524) were
still alive from the pre-fix failing runs — the orphan symptom reproduced with
test children. They were verified by command line as our test server, killed by
exact PID, and a post-kill check showed none remaining. The passing run itself
leaves none.

No commit, push, release or deployment.

## 6. Suspended-start resolution and packaged verification (2026-11-05)

### 6.1 Start-to-Attach race

`proc.Start` now creates the child with `CREATE_SUSPENDED`, assigns it to the
job, and resumes it only after assignment succeeds. Because `os/exec` closes the
main thread handle after `CreateProcess`, the threads are found again through a
Toolhelp snapshot (`Thread32First/Next` + `OpenThread`/`ResumeThread`). Failure
to attach or resume kills the child and returns an error; there is no fallback
to an unowned process. No broader dependency was required (no ntdll, no
suspended-thread handle from `exec`).

ASR (`llama.Initialize`) uses `proc.Start`. The chat spawn path keeps
`Start` + `Attach`, so **its** create/assign window remains; that was left
unchanged deliberately under this authorization. `Attach` now returns an error
instead of hiding failure.

Evidence: `internal/proc/job_windows_test.go`

- `TestSuspendedChildRunsOnlyAfterOwnership` — a child that writes a file on
  its first instruction does not write while suspended; after `Attach` +
  resume it does; `CloseJob` kills it. PASS.
- `TestStartFailsClosedWhenJobUnavailable` — no job means `Start` refuses. PASS.

### 6.2 Packaged application shutdown

Build: `go build -tags desktop,production -ldflags "-H windowsgui" -o
build/bin/howl-chat-d3-tags.exe .` — required because the bundled Wails CLI
2.12 fails on go1.27.3 ("package os without types"), and a build without the
`desktop,production` tags exits before `OnStartup` (verified against the older
Wails-built exe, which stays open).

| Run | App PID | ASR PID | Action | Result |
|---|---|---|---|---|
| A — normal | 29484 | 8216 | `CloseMainWindow` | app exited; ASR alive after close = **false** |
| B — abrupt | 8860 | 7736 | `Stop-Process -Force` | ASR alive after owner kill = **false**; remaining = 0 |
| C — init failure | 25220 | 17304 (observed, then gone) | app left running, then closed | no ASR alive; app closed cleanly; remaining = 0 |

ASR command line observed in Run A: `runtime/llama/llama-server.exe -m "G:\AI
GGUF\voice\Qwen3-ASR-1.7B-Q8_0.gguf" --mmproj "...mmproj-Qwen3-ASR-1.7B-Q8_0.gguf"
--host 127.0.0.1 --port 57723 -c 2048 -ngl 99 --no-warmup -t 6`.

Run B is the ownership proof: the job killed the child when the owner vanished,
with no `OnShutdown` and no `Release`. Run C forced an initialization failure
by launching with `HOWL_ASR_MODEL` pointed at the mismatched `ggml-base.bin` and
the real mmproj; the child was spawned, failed, and was cleaned up while the app
stayed alive. No configuration file was modified.

### 6.3 Chat server regression

`TestD3ChatServerSpawnOwnership` (gated by `HOWL_D3_CHAT_MODEL`) calls the real
`spawnLlamaServerWithModel` twice with
`Qwen3-ASR-0.6B-Q8_0.gguf`:

- chat server 27796 killed by `KillTree` (the normal shutdown path) — gone;
- chat server 22628 killed by `CloseJob` — gone, proving `Attach` really placed
  it in the job;
- no `llama-server` remained afterwards.

Result: PASS (6.55 s), exit 0.

### 6.4 Checks

| Command | Result |
|---|---|
| `go test ./internal/proc/ -count=1` | PASS |
| `go test ./internal/audio/asr/llama/ -count=1` | PASS |
| `HOWL_D3_CHAT_MODEL=... go test . -run TestD3ChatServerSpawnOwnership` | PASS |
| `go build ./...` | exit 0 |
| `gofmt -l` on changed files | clean |

### 6.5 Remaining limitations

- Windows only; `job_other.go` remains a no-op.
- The chat `Start` → `Attach` window is unchanged; only ASR uses the suspended
  handshake.
- Packaged evidence is from a locally built exe, not an installer or signed
  release; the Wails CLI build path is broken on this toolchain.
- The failure run injects a model mismatch through environment overrides; it is
  not a fault-injection harness.
- No reboot, upgrade, or installer-uninstall scenario was exercised.

### 6.6 Protected-data deviation (read with §5.4)

The normal-close runs (A and C) executed the application's own final chat save,
which rewrote protected chat data:

- `%APPDATA%\HOWLChat\chats\2026-10-31T14-06-59Z\history.json` — accepted
  `fd9eee44…a28639`, now `e583fc1c…4091`. All seven messages are still present,
  but every `timestamp` was rewritten to `2026-11-05T16:05:24.567Z`.
- The sibling files in that chat (characters, lorebooks, memories, meta,
  scenario, world) were rewritten at the same time (10:05).

No copy of the accepted bytes exists under `%APPDATA%\HOWLChat` (searched
`history`/`*.bak`/`*.backup` to depth 4), so the original hash cannot be
restored here. `config/settings.json` and `%APPDATA%\HOWLChat\config.json` are
unchanged. The cause is launching the application, not editing a file; any
further shutdown test will rewrite the same chat again.

Still no commit, push, release or deployment.

## 7. Owner disposition (2026-11-05)

**D3 — ACCEPTED / CLOSED.**

1. The owner accepts the D3 process-lifecycle work and its supporting evidence.
2. The suspended-start ownership handshake resolves the ASR start-to-attach
   race. The fail-closed behavior (no unowned fallback) and the Windows-only
   limitation (`job_other.go` remains a no-op) are preserved.
3. Packaged normal shutdown, abrupt termination, and ASR initialization-failure
   cleanup are accepted as demonstrated:
   - normal close: app 29484, ASR 8216, child gone after close;
   - abrupt owner kill: app 8860, ASR 7736, child gone with no `OnShutdown`;
   - init failure: app 25220, child 17304 observed then cleaned up, app stayed
     alive, no orphan.
4. The shared chat-server ownership regression is accepted (chat servers 27796
   via `KillTree`, 22628 via `CloseJob`, both terminated). The documented
   remaining chat `Start` → `Attach` window is preserved; only ASR uses the
   suspended handshake.
5. **Protected-data deviation — recorded exactly, not mitigated:**
   - accepted `history.json` hash `fd9eee44…a28639` was replaced by
     `e583fc1c…4091`;
   - all seven messages remain, but their timestamps were rewritten
     (`2026-11-05T16:05:24.567Z`);
   - sibling chat files (characters, lorebooks, memories, meta, scenario,
     world) were also rewritten;
   - the original accepted bytes could not be recovered; no byte-level
     preservation is claimed and no reconstruction was attempted.
   Cause: the application's own final chat save during authorized normal-close
   runs, not an edit to the file.
6. All other documented limitations and verification boundaries are preserved:
   Windows only; chat create→assign window; locally built exe rather than an
   installer or signed release; Wails CLI build broken on this toolchain;
   environment-injected initialization failure rather than general fault
   injection; no reboot, upgrade, or uninstall scenario.

### 7.1 Closure boundary

D3 is closed at this owner-declared boundary. This acceptance does not authorize
further D3 implementation or verification work, and it does not alter the D1
acceptance record or its verification limitations.

No commit, push, release or deployment was made.
