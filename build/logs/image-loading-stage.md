# Stage 3: engine loading inside the image job lifecycle

Date: 2026-11-01/02
Status: **implemented, tested, live-exercised** — pending independent review (not frozen)
Supersedes: nothing. Builds on the frozen `build/logs/image-supervisor-stage.md`.

## Objective

`ensureEngineReady` ran inside the `StartImageJob` RPC, so a slow load blocked the
caller, could not be observed as part of the job, and could not be cancelled.
Loading is now a phase of the job.

## Files changed

| File | Change |
|---|---|
| `internal/image/job/job.go` | `Start(opts, prepare)`; `loading` phase; `transition()` compare-and-set so a cancel cannot be overwritten; poller starts only after generation begins; `Status()` returns a snapshot |
| `internal/image/job/job_test.go` | 6 tests (see below) |
| `internal/image/engine/sdcpp/engine.go` | `loadRun` serializes loads; `beginLoad`/`endLoad`; `CancelJob` aborts a load by killing the child; `serverProc` guarded by `procMu`; abort-aware error paths |
| `internal/image/engine/sdcpp/supervisor.go` | `errLoadCancelled`; `setServerProc`/`getServerProc`; `abortLoad`; abort checked at loop top *and* on the exit/timeout branches of `waitReadyOrExit` |
| `internal/image/engine/sdcpp/supervisor_test.go` | 2 new tests incl. cancel-during-readiness |
| `internal/image/types/types.go` | `ImageJob.Status` documents `loading` |
| `app-image-generate.go` | `StartImageJob` returns immediately; loading runs in a `prepare` closure |
| `frontend/dist/js/image-generator-core.js` | poll handler renders `loading` as "Loading model..." |

## Observed (live, harness `http://127.0.0.1:8799`)

- `StartImageJob ok 1ms` and `ok 0s` — the RPC returns before loading. (Stage 2
  baseline for the same call was `FAILED 10m0.405s` on a bad model, `ok` only
  after the whole load on a good one.)
- Loading is visible: status went `Generating...` at 208–416 ms →
  `Loading model...` at 817–825 ms, repeated on three separate attempts.
- Cancel during loading: UI showed `Loading model...`, then `Aborted` 213 ms after
  the click; `CancelImageJob ok 0s` in the RPC log.
- Known-good generation on a fresh engine: `decode_first_stage completed,
  taking 13.52s`, `generate_image completed in 94.55s`, artifact
  `images/gallery/img_20261102_000805_600.png` written at 18:08.
- Untouched load path: `total params memory size = 12619.23MB (VRAM 0.00MB,
  RAM 12619.23MB)` — same placement as the frozen stage.

## Test-derived

`go test ./internal/image/job/` (6 tests) and `./internal/image/engine/sdcpp/`
(39 tests) pass, including the new:

```
--- PASS: TestStartReturnsWhileLoadingIsStillRunning (0.01s)
--- PASS: TestLoadingSucceedsAndAdvancesToGeneration (0.01s)
--- PASS: TestLoadingFailureIsTerminal (0.01s)
--- PASS: TestCancelDuringLoadingPreventsGeneration (0.00s)
--- PASS: TestCancelDuringGenerationKeepsTerminalState (0.07s)
--- PASS: TestJobFailsWhenTheEngineDiesEvenIfProgressErrors (1.01s)
--- PASS: TestCancelJobStopsALoadInProgress (0.03s)
--- PASS: TestWaitReadyOrExitStopsWhenTheLoadIsCancelled (0.16s)
```

Race detector (`CGO_ENABLED=1 go test -race`) is **clean** on
`internal/image/job`, `internal/image/engine/sdcpp` and `.` (root).
It found one real defect during this stage: `Status()` returned the live
`*ImageJob` while its goroutine mutated it, so the JSON encode raced the job.
`Status()` now returns a copy.

`go build ./...` exit 0, `go vet ./internal/image/... .` exit 0, `go test ./...`
exit 0, gofmt clean in every touched package.

## Source-derived (what the code now guarantees)

- A job is registered before its goroutine starts, so the returned id is always
  findable.
- `queued → loading → generating` goes through `transition()`, a compare-and-set
  under the manager mutex. A cancel that lands first changes the status, the
  transition then fails, and generation is never started for an abandoned job.
- The progress poller starts only after the job enters `generating`: it treats a
  silent backend as a dead one, which during loading would mean a false failure.
- `finish()` treats `cancelled` and `completed` as final in both directions.
- A cancel during loading kills the child (`abortLoad`), and `waitReadyOrExit`
  reports it as `errLoadCancelled`, not as a crash.
- `loadRun` serializes `LoadModel`, so two jobs cannot replace each other's
  process mid-read.
- `Cancel` releases the manager mutex before touching the engine, so stopping a
  load cannot stall status polling.

## Inferred (reasoned, not directly tested)

- Serializing loads prevents a concurrent-start race; no test forces two
  simultaneous `LoadModel` calls through a real process.
- Killing the child is a portable stop signal; only Windows was exercised.

## Unverified / not tested

- **Live cancellation during a genuinely slow (cold) load.** The live attempts
  raced a cached load: `sd-server` pid 9908 was created 17:57:33 and had finished
  loading by 17:57:34 (placement line + `listening on`), five seconds before the
  17:57:39 cancel, so the surviving process was a *loaded* server serving HTTP
  200, not an abandoned load. The cancel-during-load behaviour is covered by
  unit tests only.
- macOS/Linux: the process tests skip (`cmd.exe`); kill semantics there are
  untested. No Windows-only API was introduced beyond the existing
  `os.Process.Kill` and the `proc` job object that already guarded children.
- Wails desktop shell: the frontend `loading` label was exercised through the
  harness shim only.
- CPU-only and non-NVIDIA configurations: nothing in this change inspects
  hardware, VRAM, CUDA, model family, drives or memory size, but no such
  configuration was run.

## Observed regression candidate (NOT caused by this stage, not fixed)

On an engine instance that had a *cancelled generation* earlier in the session,
the next Z-Image run logged `decode_first_stage completed` and
`generate_image completed in 230.82s`, then answered
`500 {"error":"generate_image returned no results"}`. A fresh engine with the
same model and parameters produced the artifact normally (94.55s). This is a
reused-engine state problem in the generation path, outside this stage's scope,
and is recorded here rather than fixed.

## Still open from earlier stages

`Manager.List()` ordering, repeated `Cancel` on an already-cancelled job, the
`stderrBuf` `strings.Builder` race, krea2's conditioner-tensor mismatch, and
near-duplicate model names in the dropdown. None are touched here.

## Acceptance matrix

| Scenario | Required result | Status |
|---|---|---|
| Start with slow loading | RPC returns a job id promptly | **PASS** (unit + `ok 1ms` live) |
| Loading in progress | Job reports loading | **PASS** (unit + live UI) |
| Loading succeeds | Advances to generation | **PASS** (unit) |
| Loading fails | Terminal failure | **PASS** (unit + live krea2/flux failures) |
| Cancel during loading | Terminal cancelled | **PASS** (unit); live raced a cached load |
| Engine exits during loading | Supervisor reports promptly | **PASS** (frozen stage; still passing) |
| Late readiness after cancel | Cannot revive the job | **PASS** (unit) |
| Late completion after cancel | Cannot overwrite cancellation | **PASS** (unit) |
| Concurrent start/cancel | No race, no duplicate ownership, no orphan | **PARTIAL** — race detector clean, loads serialized, no live concurrency exercised |
| Successful generation | Still completes with its artifact | **PASS** on a fresh engine; **FAILS** on a reused post-cancel engine (see above) |

## Portability review

No hardcoded hardware, VRAM, drive letters, model names, install paths, CUDA
assumptions or load durations were added. Loading is driven by the existing
engine abstraction; cancellation stops whatever child the engine owns. The only
platform-specific code touched is the existing Windows job-object helper, which
was already in place.
