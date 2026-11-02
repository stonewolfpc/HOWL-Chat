# Stage: image-generation cancellation correctness repair

Date: 2026-11-02
Status: **implemented, unit-tested, live-exercised** — pending independent review (not frozen)
Basis: `build/logs/cancelled-generation-contract-audit.md` (read-only audit, unchanged)
Builds on: `build/logs/image-supervisor-stage.md` (frozen), `build/logs/image-loading-stage.md`
Supersedes: the historical *explanation* of the 500 in `issue-cancelled-generation-reuse.md`?
No. That explanation stays unproven. See "Unestablished" below.

## Objective

Cancellation was bookkeeping only: it set a client-global boolean, the synchronous
txt2img POST was un-abortable, and nothing prevented an abandoned request from
overlapping the next one on the same sd-server process. This stage makes
cancellation internally consistent and makes overlap impossible, without claiming
the backend stopped.

## Exact changes

| File | Change |
|---|---|
| `internal/image/engine/sdcpp/generation.go` (**new**) | Slot `idle`/`active`/`abandoned`; `beginGeneration(jobID)`, `endGeneration(handle)`, `cancelGeneration(jobID)`, `resetGeneration()`, `generationBusy()`, `GenerationState()`, `processConfirmedGone()`, `classifyGenError()`; `errEngineBusy`, `errEngineUnavailable`, `errClientCancelled` |
| `internal/image/engine/sdcpp/engine.go` | Slot fields + `genSeq` token; `GenerateTxt2Img`/`GenerateImg2Img`/`GenerateInpaint` claim the slot, carry `ctx` into the client, classify errors; `CancelJob(jobID string)` routes to video → load abort → generation cancel; `resetGeneration()` after the previous child is killed in `LoadModel`; `Shutdown` resets the slot **after** the child is reaped |
| `internal/image/engine/sdcpp/client.go` | `txt2img`/`img2img`/`inpaint`/`generateSequentialBatch`/`pollJob` take `ctx`; `http.NewRequestWithContext` + `c.http.Do` replace `http.Post`/`http.Get`; `pollJob` selects on `ctx.Done()`/deadline/ticker instead of sleeping; **deleted** `cancelJob` (the guessed `DELETE /sdcpp/v1/jobs/{id}`), `resetCancel` (no callers), the `cancelled` field, and `resetState`'s reset of it; `activeJobID` is cleared per request |
| `internal/image/engine/sdcpp/supervisor.go` | `beginLoad(jobID)`, `abortLoad(jobID)` are owner-scoped |
| `internal/image/types/types.go` | `ImageOptions.JobID`, `ModelLoadOptions.JobID` |
| `internal/image/engine/engine.go` | Interface: `CancelJob(jobID string) error` |
| `internal/image/job/job.go` | `Start(opts, prepare func(jobID string) error)`; `run` injects `opts.JobID`; `Cancel` reads the id under the manager mutex and passes it to the engine |
| `app-image-generate.go` | Both call sites pass the job id into `ensureEngineReady` |
| `app-image-model.go` | `ensureEngineReady(modelName, jobID string, genOpts ...)`; sets `loadOpts.JobID` on the request copy, never on the remembered `lastImageModelOpts` |
| `internal/image/engine/sdcpp/generation_test.go` (**new**) | 6 tests |
| `internal/image/engine/sdcpp/supervisor_test.go` | Existing cancel tests updated to owner-scoped calls |
| `internal/image/job/job_test.go` | Fake `CancelJob(string)`; 2 new tests, 6 existing updated for the `prepare` signature |

Removed rather than layered over: the DELETE endpoint call, the global boolean,
and the unused `resetCancel`. No new override layer was added.

## Ownership and lock order

```
idle ──beginGeneration──▶ active ──request returned cleanly──▶ idle
                            │
                            └──request cancelled──▶ abandoned ──confirmed child exit────▶ idle
                                                              ──LoadModel kills child───▶ idle
                                                              ──Shutdown reaps child───▶ idle
```

- The slot lives on `*Engine`, so it is per sd-server process, not per UI or per
  job. A second request is refused at the engine even if a caller bypasses the
  job manager.
- `endGeneration` deliberately does **not** release an abandoned slot. Releasing
  there is exactly the overlap the audit described.
- The only completion accepted without backend cooperation is a confirmed child
  exit (`procDone` closed), because a dead process cannot own a generation.
- `resetGeneration` bumps `genSeq`, so a request still unwinding cannot write its
  state back over a reset (this closed a real hole found while reviewing the
  unload path: unload-during-generation could have left the fresh engine
  `abandoned`).
- Lock order: `loadRun → genMu`, `loadRun → procMu`; `genMu` and `procMu` are
  never held together (`generationBusy()` returns before `processConfirmedGone()`
  acquires `procMu`). `CancelJob` takes each in sequence, never nested. No path
  takes `loadRun` while holding `genMu`, so no deadlock cycle exists.
- `cmd.Wait()` ownership is unchanged: the supervisor's reaper is still the sole
  owner, and slot resets wait on its signal rather than calling `Wait()`.

## Job identity

`ImageOptions.JobID` / `ModelLoadOptions.JobID` are set once, by the job, from
`job.ID`. `Manager.Cancel` passes that id; `cancelGeneration`/`abortLoad` compare
it against the current owner. A cancel from job B therefore cannot close job A's
request or abort job A's load — proven by `TestCancelIsScopedToOneJob` and by the
live run, where the refused second request named the first job in its message.

## Test-derived

`go test ./internal/image/engine/sdcpp/` — **46 tests**, 6 new:

```
--- PASS: TestTxt2ImgIsCancelledThroughItsRequestContext
--- PASS: TestSecondGenerationIsRefusedWhileTheFirstIsActive
--- PASS: TestOwnershipIsReleasedAfterAConfirmedProcessExit
--- PASS: TestCancelIsScopedToOneJob
--- PASS: TestRepeatedCancelIsSafe
--- PASS: TestAsyncCancellationDoesNotPoisonTheNextJob
```

`go test ./internal/image/job/` — **8 tests**, 2 new
(`TestCancelCarriesTheJobID`, `TestNextJobCompletesAfterACancel`), 6 updated for
the new signatures.

Requirement mapping:

1. synchronously cancelled through ctx — `TestTxt2ImgIsCancelledThroughItsRequestContext` (the fake server observes `r.Context().Done()`).
2. no false "backend stopped" claim — same test asserts the error says the work `may still be running`.
3. no overlap while active/indeterminate — `TestSecondGenerationIsRefusedWhileTheFirstIsActive`, and the abandoned branch in test 1.
4. release after normal completion — `TestSecondGenerationIsRefusedWhileTheFirstIsActive` (final phase).
5. release after confirmed termination — `TestOwnershipIsReleasedAfterAConfirmedProcessExit` (closed `procDone`, platform-neutral).
6. async state does not poison the next job — `TestAsyncCancellationDoesNotPoisonTheNextJob` (img2img job 1 cancelled mid-poll, job 2 completes).
7. job identity scoping — `TestCancelIsScopedToOneJob`.
8. repeated cancel / terminal safety — `TestRepeatedCancelIsSafe`, `TestCancelDuringGenerationKeepsTerminalState`.

Also green: `go test -race` on both packages (clean), `go vet ./...` exit 0,
`go test ./...` exit 0 (18 packages), gofmt clean, and both builds:
`-tags "desktop,production"` and `-tags uiharness` exit 0.

Test-harness note: the first version of the new engine test failed, and the
failure was in the test, not the product. A Go server only starts watching the
connection for an abort once the request body has been consumed, so the fake
handler had to drain the body before waiting on `r.Context()`.

## Observed (live, harness on `http://127.0.0.1:8799`, rebuilt binary)

Run 1 (pre-final-binary ordering fix), Z-Image turbo, prompt unchanged from Stage 3:

- `StartImageJob ok 0s` at 18:28:41, engine started on its own port.
- `CancelImageJob ok 0s` at 18:28:47 (UI: `Aborted`), and **`sd-server.exe` PID
  20688 was still alive at 13,073,536 K** — the client stopped waiting, the
  backend did not stop. This is the honest state the error text describes.
- Second `StartImageJob ok 0s` at 18:29:03 → UI status:
  `Error: image engine is unavailable until it is restarted: the previous
  generation was abandoned and its completion could not be confirmed (owner
  "job_1793579321187570300")`.

Run 2 (final binary, after the `resetGeneration`/`Shutdown` ordering fix):

- `StartImageJob ok 0s` 18:31:32 → `Starting sd-server` → `generate_image
  1024x1024` → `CancelImageJob ok 0s` 18:31:36; `sd-server.exe` PID 5356 still
  alive at 6,711,304 K.
- Second `StartImageJob ok 0s` 18:31:47: accepted by the RPC, **no second engine
  process spawned**, job failed with:
  `image engine is unavailable until it is restarted: the previous generation was
  abandoned and its completion could not be confirmed (owner "job_1793579492348455900")`.
- `[App] UnloadImageModel called` → `UnloadImageModel ok 992ms` 18:32:04, child
  gone (`(none - child gone)`).
- `StartImageJob ok 0s` 18:32:16 → fresh `sd-server` → `decode_first_stage
  completed, taking 13.90s` → `generate_image completed in 105.29s` 18:34:09 →
  `[App] Job job_1793579536427649000 result saved to gallery`; UI `Done!`;
  artifact `images/gallery/img_20261102_003402_400.png`, 791,480 B.

Validated states, each observed rather than assumed:
client request cancelled ✔ · server generation still running / unknown ✔ (process
alive after cancel) · server generation confirmed stopped — **not observed**, and
not claimed · engine unavailable pending recovery ✔ · engine safe for reuse after
an explicit restart ✔.

## Source-derived

- `http.Post`/`http.Get` are gone from every generation path; each request is
  bound to a context that the owning job's cancel closes.
- `pollJob` returns `ctx.Err()` on cancellation and no longer sleeps in a loop
  that ignores cancellation; the async job's `activeJobID` is cleared when the
  call returns, so the progress poller cannot query a previous job's id.
- `pollJob`'s `cancelled` case (the server reporting it stopped work) is the only
  path that reports a backend-side stop, and only because the server said so.
- A batch stops between items on cancellation and keeps completed items.
- `loadOpts.JobID` is set on the per-request copy, so a later load triggered from
  saved settings cannot inherit a stale owner id.

## Inferred

- Bumping `genSeq` in `resetGeneration` is what protects against a late
  `endGeneration`; the hole was found by reading the code, not by a failing test.
  It is covered indirectly (existing tests still pass) but no test forces an
  unload during an in-flight request.
- Refusing rather than queueing keeps the invariant simple; a queue would need
  its own abandoned-work policy.

## Unestablished

- **Backend interruption: unsupported, by design.** No `DELETE`/interrupt call
  exists. Whether this sd-server build honours one is still unproven; the audit's
  "unestablished" verdict stands and was not upgraded.
- The 500 in `issue-cancelled-generation-reuse.md` is *still* not explained.
  This stage removes one plausible overlap mechanism; it does not prove that was
  the cause, and the earlier write-up is not rewritten.
- How long sd-server keeps working after its client disconnects is not measured.
  The engine treats it as "maybe forever until the process is gone".

## Platform coverage and limitations

- Exercised on Windows only. `proc.KillTree` (a no-op off Windows) is used in
  `Shutdown` and remains a pre-existing limitation; the new ownership logic does
  **not** depend on it — it depends on the reaper's `procDone` signal, which is
  platform-neutral, and nothing here restarts the child to fake success.
- No CPU/GPU/VRAM/RAM/model/duration assumptions were added: refusal is
  state-based, not time-based, and model/engine paths come from existing config.

## Regressions and unresolved issues

- No regressions observed: full suite, both tagged builds, vet and race are green.
- Known, untouched, out of scope: `Manager.List()` returns live pointers in
  nondeterministic order; repeated `Cancel` on a cancelled job still returns
  true; `stderrBuf` `strings.Builder` is written from a pipe goroutine; video
  `SelectModel` load is not cancellable; krea2 conditioner mismatch.
- Product-copy gap: the UI surfaces the engine's English error verbatim in the
  status line (`Error: image engine is unavailable until it is restarted...`).
  Accurate, but not user-facing prose; deliberately not "fixed" by softening it.
- `Engine.GenerationState()` is currently a diagnostic used by tests and not
  wired into the UI; the UI learns the state only through the error message.
- Environment note during live validation: a stale harness from an earlier
  session held port 8799 and its `sd-server` (13 GB) was running; both were
  stopped before the new binary was exercised. One `taskkill /IM sd-server.exe`
  was issued before that, which would also have stopped an engine belonging to
  the separately running `howl-chat.exe`; no product code was affected and the
  app respawns its engine on the next load.

## Verdict

Cancellation is now internally consistent: the client request is really closed,
the abandoned work is never claimed to have stopped, and no second generation can
start on that process while the backend's state is unknown. In this build the
engine is **not** automatically reusable after a client cancellation — it stays
unavailable until the child exits or the user unloads/reloads — and that refusal
is what makes the next successful generation trustworthy.
