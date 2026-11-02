# Stage: generation ownership follow-up (F1, F2, F3, F5, F6, F7)

Date: 2026-11-02
Status: **implemented, unit-tested, build-verified** — pending independent review (not frozen)
Basis: `build/logs/image-cancellation-repair-review.md` (the findings), `build/logs/cancelled-generation-contract-audit.md`
Builds on: `build/logs/image-cancellation-repair-stage.md`
F4 (video process ownership) was explicitly **not** changed; see the last section.
No live generation was run in this stage. No Git operations were performed.

## Defects addressed

| Finding | Defect | Resolution |
|---|---|---|
| F1 | `classifyGenError` ran after `endGeneration`, which cancels the handle's context, so **every** failure was reported as a client cancellation | cancellation is captured **before** cleanup and passed in; `describeGenError` never re-derives it |
| F2 | any non-cancel error released the generation slot, admitting a second request while the backend might still be working | explicit outcome classification: only a confirmed completion or a confirmed process exit releases ownership |
| F3 | `CancelJob` returned after `abortLoad` whenever any load was in progress, swallowing a generation cancel; `cancelVideo` had no owner | routing is by ownership: video (owned) → load (owned) → generation (owned); each step reports whether it owned the operation |
| F5 | an empty `JobID` could act as "whichever job is active" | identity is required for cancellation: empty ids cancel and abort nothing |
| F6 | an unwinding request could clear a newer job's progress id | progress ownership is tokenised; a cleanup only clears its own token |
| F7 | an unload/reload/exit was reported as a client cancellation | the close **reason** survives cleanup and is reported instead |

## Source paths changed

| File | Change |
|---|---|
| `internal/image/engine/sdcpp/generation.go` | `genOutcome` (`completed`/`notDispatched`/`cancelled`/`ambiguous`); `genClose` (`none`/`reload`/`shutdown`/`processExited`); `settleGeneration`, `classifyGeneration`, `describeGenError`, `closeReason`; `endGeneration(h, outcome)`; `resetGeneration(reason)`; `cancelGeneration` requires identity; new sentinels `errEngineUncertain`, `errEngineStopped`, `errNotDispatched` + `notDispatched()` |
| `internal/image/engine/sdcpp/engine.go` | `genClose` and `videoOwner` fields; three `Generate*` paths use `settleGeneration`; `CancelJob` routes by ownership and ignores empty ids; `LoadModel` resets with `closeForReload`; `Shutdown` with `closeForShutdown` |
| `internal/image/engine/sdcpp/supervisor.go` | `abortLoad(jobID) bool` — empty id is a no-op, and it only aborts the load it owns |
| `internal/image/engine/sdcpp/video.go` | `runVideoCLI(..., owner)` records/clears `videoOwner`; `cancelVideoOwned(jobID)`. **`stopServer`, `KillTree`, and `os.Process.Wait()` are unchanged** |
| `internal/image/engine/sdcpp/client.go` | `pollTimeout` field (600 s default, overridable in tests); `setActiveJob`/`resetActiveJob(token)`/`activeJob`; `resetState` bumps the token; local pre-dispatch failures wrapped with `notDispatched` (txt2img marshal/build, img2img marshal/build, inpaint path/mask/read/marshal/build) |
| `internal/image/engine/sdcpp/ownership_test.go` | 9 new tests |
| `internal/image/engine/sdcpp/generation_test.go` | updated to the new `endGeneration`/`resetGeneration` signatures |

## Ownership state transitions

```
idle ──beginGeneration──▶ active ──completed──────▶ idle
                            │      ──notDispatched──▶ idle   (nothing was sent)
                            │      ──cancelled──────▶ abandoned
                            │      ──ambiguous──────▶ abandoned
                            │      ──error + child gone──▶ idle
abandoned ──beginGeneration (child confirmed gone)──▶ idle, then active
any ──LoadModel kills child / Shutdown reaps child / reset──▶ idle (reason recorded)
```

- `endGeneration` still never releases an abandoned slot on its own; it releases
  only when the outcome proves completion (or the child is confirmed gone).
- `resetGeneration(reason)` bumps `genSeq` (a stale handle cannot write state back)
  and records the reason **before** cancelling the context, which is what lets the
  request report "unloaded" instead of "cancelled".
- `beginGeneration` clears the reason, so it cannot leak into the next request.

## Error classification policy

| Class | Trigger | Slot | Reported as |
|---|---|---|---|
| completed | `err == nil` (a result was read) | released | — |
| notDispatched | `errors.Is(err, errNotDispatched)`: marshal, request build, missing/unreadable input | released | the underlying error, unchanged |
| cancelled | `h.ctx.Err() != nil` before cleanup | retained (`abandoned`) | `errClientCancelled` + "may still be running"; or `errEngineStopped` when the close reason is unload/reload/exit |
| ambiguous | anything else after a request was sent: HTTP error status, empty/malformed result, transport loss, poll timeout | retained (`abandoned`) | `errEngineUncertain` + the backend's own message |

Evidence behind each choice:
- **completed** — a successful response was read; the handler returned with results.
- **notDispatched** — no HTTP request existed, so no backend work could be running.
  Wrapping these explicitly is what keeps a typo (missing init image) from wedging
  the engine.
- **cancelled** — the audit established that closing the request is the only thing
  this client can do, and that backend interruption is unestablished.
- **ambiguous** — the conservative choice the stage asked for: the client cannot
  distinguish "the server finished and errored" from "the server is still working",
  so ownership is retained. Deliberate consequence: a 400/500 keeps the engine
  unavailable until the child is unloaded or replaced. That is honest but stricter
  than before, and is the main behavioural trade recorded here.

Backend interruption remains unsupported: no interrupt/DELETE path was added.

## Tests and actual results

`go test ./internal/image/engine/sdcpp/` — **55 tests**, 9 new, all pass:

```
--- PASS: TestHTTPErrorIsNotReportedAsAClientCancellation      (500 → errEngineUncertain, refused next, not errClientCancelled)
--- PASS: TestEmptyResultKeepsTheEngineUnavailable             (200 + no images → ambiguous, refused next)
--- PASS: TestTransportFailureKeepsTheEngineUnavailable        (closed listener → ambiguous)
--- PASS: TestPollTimeoutKeepsTheEngineUnavailable             (pollTimeout=100ms → ambiguous)
--- PASS: TestCancellingAGenerationIsNotSwallowedByAnotherJobsLoad (F3, deterministic)
--- PASS: TestEmptyJobIdentityCannotCancelWork                 (F5; work survives CancelJob(""), B's load untouched)
--- PASS: TestStaleProgressOwnerCannotClearANewerJob           (F6)
--- PASS: TestShutdownIsNotReportedAsAClientCancellation       (F7 → errEngineStopped)
--- PASS: TestLocalFailureReleasesTheSlotWithoutBlame          (notDispatched → idle, plain error)
```

Each proves both the classification **and** whether another generation is admitted
(except the local-failure case, which asserts the slot is free by claiming it
directly — the test server answers 404, so actually generating there would land in
the ambiguous branch).

`go test ./internal/image/job/` — **8 tests**, all pass (unchanged; the job layer
consumes the new errors without modification).

One failure occurred during the stage and was in the test, not the product: the
last assertion of `TestLocalFailureReleasesTheSlotWithoutBlame` expected a real
generation to succeed against a `http.NotFoundHandler()` server, which the
conservative policy correctly refuses. The assertion now claims the slot directly.

## Concurrency and race evidence

- `go test -race` on `internal/image/engine/sdcpp` and `internal/image/job`: clean.
- `go vet ./...`: exit 0. `go test ./...`: 18 packages, exit 0. gofmt: clean in both
  touched packages. Both tagged builds exit 0 (`desktop,production` and `uiharness`).
- F3's test is the interleaving itself: job B holds a load while job A generates,
  and A's cancel must still arrive. Before the fix this test hangs by construction.
- Lock order is unchanged: `loadRun → genMu`, `loadRun → procMu`, `loadMu → procMu`,
  `genMu`/`procMu` never nested (`endGeneration` reads `processConfirmedGone()`
  *before* taking `genMu`). `videoMu` remains a leaf. No new nesting was added.
- F6 is proven at the client level with explicit tokens (`old` then `new`, then a
  stale `resetActiveJob(old)`), because a full HTTP interleaving of that race is
  not deterministic.

## Remaining limitations

- **F4 is untouched and still open**: `generateVideo` does not claim the generation
  slot, and `stopServer()` still kills sd-server unconditionally (including under an
  active generation), reads `serverProc` without `procMu`, and calls
  `os.Process.Wait()` alongside the reaper's `cmd.Wait()`. This stage only stopped
  one job's cancel from killing another job's video render; the process-level
  collision remains exactly as the review recorded it.
- The conservative policy makes any post-dispatch failure end in "unavailable until
  restarted". Recovery is unload/reload, which the UI already offers.
- `progress()` still reads a single `activeJobID` for the whole client. The token
  prevents stale clears, but two simultaneous async requests would share it — they
  cannot occur while the slot is enforced.
- The two legacy direct callers (`GenerateImage`, `SetActiveImageModel`) submit with
  no job id and are therefore not cancellable. That is now an explicit, documented
  contract rather than an ambiguity.
- Backend interruption: **unestablished**, unchanged. The historical HTTP 500:
  **still unexplained**; this stage narrows the guarantee but establishes no cause.
- No live generation, harness run, or manual UI exercise was performed in this
  stage; every claim above is source-derived or test-derived.

## Platform-specific uncertainty

Only Windows was built and tested. Nothing here depends on Windows process
behaviour: the release-on-exit path reads the reaper's `procDone` channel, the
cancellation paths are HTTP contexts, and `KillTree` is not used as a correctness
guarantee. `proc/job_other.go` remains a no-op set off Windows, so the video kill
paths (F4) and `Shutdown`'s kill remain unverified there — unchanged by this stage.

## Confirmation

Video **process** ownership was not modified: `stopServer`, its `proc.KillTree`,
`os.Process.Wait()` call, and the reaper were left exactly as they were. The only
video change is cancellation identity (`videoOwner`, `cancelVideoOwned`).
Earlier records were not modified.
