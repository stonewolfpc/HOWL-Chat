# Stage: image-engine process supervision

Date: 2026-11-01
Status: **closed with evidence** (frozen; do not expand retroactively)

## Why this stage existed

sd-server could die and the app would never say so. The UI sat on
"Generating... Step 8 of 8" indefinitely, and a failed model load cost the full
ten-minute timeout even when the child had printed the real error in seconds.

The read-only audit found three defects, not a missing architecture:

1. `LoadModel` started the child and never called `Wait()`. The exit status was
   never observed; `IsRunning()` guessed from HTTP instead.
2. The job poller set `failed = true` on the first progress error and then
   `continue`d past its own `IsRunning()` check on every later tick, so the one
   mechanism that could end the job was disabled by the first symptom.
3. The generation HTTP client had `Timeout: 0`, and progress requests shared it.

## What shipped

| File | Change |
|---|---|
| `internal/image/engine/sdcpp/supervisor.go` (new) | `startProcessReaper`, `processDone`, `exitError`, `waitReadyOrExit`, `reapStoppedProcess`, sentinel `errProcessExited` |
| `internal/image/engine/sdcpp/engine.go` | reaper owns `Wait()`; readiness waits on the exit channel; `IsRunning()` reads the reaper's verdict; pre-load kill and `Shutdown` wait for the reaper instead of calling `Wait()` again |
| `internal/image/engine/sdcpp/client.go` | `controlTimeout` (3s) + `httpShort` for `probe`/`progress`; generation keeps the unbounded client |
| `internal/image/job/job.go` | liveness checked on every tick, before progress; `failed` only gates progress |
| `internal/image/engine/sdcpp/supervisor_test.go` (new) | 3 regression tests |
| `internal/image/job/job_test.go` (new) | 1 regression test with a fake engine |

## Evidence

Live, same model, same failure, real binaries:

```
before:  2026/11/01 ... StartImageJob FAILED 10m0.405s  engine ready failed: sd-server did not become ready in 10m0s
after:   2026/11/01 17:42:03 [rpc] StartImageJob FAILED 858ms
         engine ready failed: krea2_turbo-Q8_0 failed to load:
         engine process exited while loading: exit status 1
```

The UI now shows the engine's own words without opening a log:

```
Error: engine ready failed: krea2_turbo-Q8_0 failed to load: engine process exited
while loading: exit status 1 Last output from sd-server: ... [E] Conditioner model
tensor 'text_encoders.llm.visual.merger.linear_fc1.weight' not in model metadata
```

Tests:

```
--- PASS: TestWaitReadyOrExitReturnsTheExitWhenTheChildDies (0.03s)
--- PASS: TestWaitReadyOrExitSucceedsWhenTheServerAnswers (0.00s)
--- PASS: TestIsRunningIsAProcessQuestion (0.03s)
--- PASS: TestJobFailsWhenTheEngineDiesEvenIfProgressErrors (1.01s)
```

Checks: `go build ./...` 0, `go vet ./internal/image/... .` 0, gofmt clean on
every touched file, all pre-existing sdcpp and app tests green.
Logs: `build/logs/harness-run10.log` (845ms measurement),
`build/logs/harness-run11.log` (858ms + accurate message).

The reaper race was caught by its own test: the first version registered the
done-channel inside the goroutine, so `waitReadyOrExit` could read nil and miss
the exit; the test failed at exactly 30s before the fix shipped.

## Acceptance matrix (honest)

| Requirement | Status |
|---|---|
| Child exits during startup -> readiness returns the actual exit promptly | **PASS** (unit 0.03s; live 858ms vs 10m0.405s) |
| Child exits during generation -> job reaches terminal failure | **PASS** (unit only; not re-run live after the fix) |
| Progress endpoint errors -> liveness monitoring continues | **PASS** (unit) |
| HTTP request stalls -> bounded request returns control | **PARTIAL** — 3s control timeout implemented; no test asserts the stall path |
| Process reaped exactly once | **BY CONSTRUCTION** — single `Wait()` caller; not asserted by a test |
| stderr capture cannot block the child | **N/A** — `cmd.Stderr` is an `io.Writer`, not a pipe, so there is no pipe to fill |
| User cancels during loading | **NOT APPLICABLE YET** — loading is not a job phase (see below) |
| Late success cannot overwrite a failure | **PRE-EXISTING GUARD** in `finish()`; not covered by a new test |

## Known open items (do not fix inside this stage)

1. `stderrBuf` is a `strings.Builder` written by the stderr-copying goroutine and
   read on the error path: a real data race, pre-existing, left in place.
2. `ensureEngineReady` still runs inside the `StartImageJob` RPC. Loads now fail
   fast, but a successful slow load is still invisible and uncancellable.
3. `Manager.List()` reverses unordered map iteration while claiming
   "newest-first"; `Cancel()` treats an already-`cancelled` job as cancellable.
4. krea2 variants (3 registry entries: `Q8_0` 13.71GB on H:, `Q4_K_M` twice —
   H: and G:) all fail with the same conditioner-tensor mismatch against this
   sd.cpp build (`c6dfa9d`). Requires a text-only Qwen2.5-VL encoder or a newer
   sd.cpp; not a code defect.
5. The dropdown still lists near-duplicate names (`krea2_turbo-Q8_0` vs
   `krea2_turbo-Q4_K_M` x2); `preferModelPath` only helps exact-name collisions.

## Next stage contract (agreed, not started)

- `StartImageJob` returns a job id promptly.
- Loading becomes an observable phase (`loading`).
- Cancellation works during loading and terminates the owned process tree.
- Startup failure becomes a terminal job failure, not an RPC error.
- Late readiness or completion cannot overwrite cancellation or failure.
- Z-Image's working generation path is unchanged.

Placement planning stays out of that stage.
