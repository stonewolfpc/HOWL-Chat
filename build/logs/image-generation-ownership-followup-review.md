# Review: generation ownership follow-up (read-only)

Date: 2026-11-02
Mode: **read-only**. No product code, tests, config or records were modified, and
no build/test/run was executed. Claims below come from reading the current source
and the records, not from the stage author's summary.

Target: `build/logs/image-generation-ownership-followup-stage.md` (F1, F2, F3, F5,
F6, F7 of `build/logs/image-cancellation-repair-review.md`).

**Verdict: PASS WITH FINDINGS — freeze-ready for the ownership guarantees, with one
incomplete policy application and one F7 race that should be corrected in the next
bounded change.** The central guarantees hold in source: an ordinary error can no
longer become `errClientCancelled`, an unconfirmed outcome no longer frees the slot,
a load no longer swallows another job's generation cancel, empty ids cancel nothing,
progress ids are token-owned, and unload/reload/exit report their own reason. No
finding below invalidates those; F1 and F2's core are correctly implemented.

## Findings

### R1 — MEDIUM: img2img's pre-dispatch failures are not marked `notDispatched`, so a user typo wedges the engine

**Source-proven** (`internal/image/engine/sdcpp/client.go:766-773`):

```go
if opts.InitImagePath == "" {
    return nil, fmt.Errorf("img2img requires an init image path")     // plain error
}
imgData, err := os.ReadFile(opts.InitImagePath)
if err != nil {
    return nil, fmt.Errorf("read init image: %w", err)                // plain error
}
```

Execution path: `GenerateImg2Img` → `settleGeneration` → `classifyGeneration` sees
no cancellation and no `errNotDispatched` → `genOutcomeAmbiguous` → slot
`abandoned` → every later request is refused with `errEngineUnavailable` until the
user unloads.

Consequence: a mistyped or moved init image is reported as "the engine may still be
working on it", and generation is blocked until an unload. The stage record claims
"Local failures that never left the process (marshal, request build, missing init
image or mask) are wrapped with errNotDispatched" — that is true for inpaint
(`client.go:1022-1035`, all four checks wrapped) and for img2img's marshal/build,
but **not** for these two img2img checks. The asymmetry is the defect: identical
inputs behave differently depending on the mode.

Required correction: wrap both with `notDispatched(...)` exactly as inpaint does,
and add a test asserting `GenerateImg2Img` with no init image leaves the slot
`idle` and the error is not `errEngineUncertain`. Not corrected here (read-only).

### R2 — LOW/MEDIUM: the close reason is read after cleanup, so a new request can rewrite an unload as a cancellation

**Source-proven interleaving** (`generation.go`, `settleGeneration` → `endGeneration`
→ `describeGenError` → `closeReason`).

`settleGeneration` calls `endGeneration` (which sets the slot `idle`) and only then
`describeGenError`, which reads `e.genClose`. `beginGeneration` sets
`e.genClose = closeNone`. So:

1. `Shutdown` → `resetGeneration(closeForShutdown)` → request A's context closes.
2. A unwinds; `endGeneration` marks the slot idle and clears `genCancel`.
3. A fresh generation calls `beginGeneration` and sets `genClose = closeNone`.
4. A's `describeGenError` now reads `closeNone` → reports
   `errClientCancelled` ("client request cancelled") for an unload.

Consequence: F7's guarantee is correct in the single-request case (which the test
covers) but can be inverted under concurrency. The window is small (between
`endGeneration` and `describeGenError`).

Required correction: capture the reason into the handle, or have
`describeGenError` receive the reason captured in the same critical section that
decides the outcome. Alternatively classify and describe before releasing, using a
local snapshot.

### R3 — LOW/MEDIUM: a video cancel can be lost between `cmd.Start()` and owner assignment

**Source-proven** (`internal/image/engine/sdcpp/video.go:337-344`):

```go
if err := cmd.Start(); err != nil { ... }
proc.Attach(cmd)

e.videoMu.Lock()
e.videoProc = cmd.Process
e.videoOwner = owner     // set only here
```

A `CancelJob(jobID)` landing after `Start` but before this assignment finds
`videoProc == nil` → `cancelVideoOwned` false → `abortLoad` false →
`cancelGeneration` false → the cancel is silently dropped while the render runs to
completion and the job is already marked `cancelled`.

Consequence: a cancelled video keeps the GPU while the UI says it was cancelled —
the same "cancel did not stop the work" class, at a narrower window than before.
Pre-existing in shape (the old `cancelVideo` had the same window), not introduced by
this stage, but this stage is where video cancel identity was added.

Required correction (when F4 is taken): register the process and owner before or
atomically with `Start`'s success handling, or accept a pre-registered owner id that
`cancelVideoOwned` can honour before the process exists.

### R4 — LOW: `notDispatched` and `errEngineUncertain`/`errEngineStopped` drop the underlying error for `errors.Is`/`As`

`notDispatched` uses `fmt.Errorf("%w: %v", errNotDispatched, err)`, and
`describeGenError` uses `%v` for the wrapped error. Messages survive, but callers
cannot unwrap to the transport/HTTP error (e.g. `*url.Error`, status codes).
Consequence: diagnostic-only. Correction optional; if done, use `%w` for the inner
error and confirm the job layer does not double-report.

### R5 — LOW: non-cancellable legacy callers are silent

`CancelJob("")` returns `nil` (`engine.go`, `CancelJob`), and the two direct callers
(`app-image-generate.go` `GenerateImage`, `app-image-model.go` `SetActiveImageModel`)
never set `JobID`. Nothing tells the user those operations cannot be cancelled; only
code comments record it. F5's requirement is met (an empty id cannot target
anything), but the review asked whether the behaviour is surfaced — it is not.
Correction optional; documented contract is acceptable for a legacy path.

### R6 — LOW: sequential-batch partial dispatch is classified by the last item

`generateSequentialBatch` returns the failing item's error. If item 2 fails locally
(errNotDispatched) after item 1 completed, the whole request is classified
`notDispatched` and the slot is released. Item 1's response was already read, so no
work is outstanding — **inferred safe**, not proven. No test covers it.

## Confirmed passing (source-derived)

- **F1**: cancellation is sampled before cleanup (`settleGeneration` computes
  `h.ctx.Err()` first), `endGeneration(h, outcome)` takes it as an argument, and
  `describeGenError` never re-derives it. A non-cancel error cannot reach
  `errClientCancelled`; the 500 / empty-result / transport / poll-timeout tests
  assert exactly that (`ownership_test.go`).
- **F2 (core)**: `endGeneration` retains ownership for `cancelled` and `ambiguous`,
  releasing only for `completed`/`notDispatched` or a confirmed child exit
  (`processConfirmedGone` on the reaper's `procDone`). HTTP error statuses and poll
  timeouts are treated as unconfirmed, matching the review's instruction not to
  assume a response proves backend completion.
- **F3**: `CancelJob` routes video → load → generation by owner, each step returning
  whether it owned the operation; `abortLoad` requires `loading && loadOwner == jobID`;
  the generating-job cancel reaches `cancelGeneration` while another job loads
  (test drives this deterministically). Video cancel is owner-scoped.
- **F5**: `cancelGeneration("")` and `abortLoad("")` return false; `CancelJob("")`
  does nothing; the direct path survives an empty-id cancel (test).
- **F6**: `setActiveJob` returns a token; `resetActiveJob(token)` only clears its
  own; `resetState` bumps the token so a pre-shutdown owner cannot clear a
  post-shutdown one. txt2img does not use ids (no async job), batches use one token
  per item, double-reset is idempotent.
- **F7 (single-request)**: `resetGeneration(reason)` records the reason before
  cancelling; `beginGeneration` clears it; `Shutdown` reports `errEngineStopped`
  (test).
- **Lock order** unchanged and still acyclic: `loadRun → genMu/procMu`,
  `loadMu → procMu`, `genMu`/`procMu` never nested (`endGeneration` reads
  `processConfirmedGone()` before locking `genMu`), `videoMu` leaf.
- **Records**: the audit, both prior reviews, the frozen stage records and the
  repair record are unmodified; the follow-up record's test counts match the tests
  present (55 engine, 8 job) and its claims match the source except R1's wording.

## F4 boundary — confirmed untouched

`generateVideo` still never calls `beginGeneration`; `stopServer`
(`video.go:116-124`) still calls `proc.KillTree`, then `os.Process.Wait()`, while
`startProcessReaper` owns `cmd.Wait()` for sd-server; `serverProc` access there is
still not `procMu`-guarded. The only video change is cancellation identity
(`videoOwner`, `cancelVideoOwned`, owner parameter on `runVideoCLI`). The record's
claim that video process ownership was not modified is accurate; the underlying
double-`Wait` and unconditional `stopServer` kill remain open exactly as F4 recorded
them.

## Test adequacy

Meaningful and behavioural: the F1/F2 tests use real `httptest` servers and assert
both the sentinel and whether a second generation is admitted; F3's test is the
interleaving itself (it would hang before the fix); F5 proves an empty id leaves the
work running.

Gaps, none of which hide a guarantee failure except R1's:
- no test for a *local* img2img failure (which would have caught R1);
- no test for concurrent unload with completion (R2's window);
- no test for two job ids cancelling different operations simultaneously (F3 covers
  sequential non-owner calls only);
- no video cancellation test (R3);
- F6 is proven at the client primitive, not through a full HTTP unwind (the record
  states this);
- no test for a transport failure *after* partial transmission.

## Portability

No Windows-only assumption was added: release-on-exit reads the reaper's channel,
cancellation is HTTP context based, and `KillTree` is not a correctness dependency.
Only Windows was built/tested; the video kill paths and `Shutdown`'s kill remain
unverified off Windows, unchanged by this stage. No hardware, model, drive or
duration assumption appears in the change.

## What remains unproven

- Whether sd-server stops work on client disconnect: **unestablished**, unchanged.
  The code never claims it does.
- Whether an unconfirmed failure leaves work running: **unestablished**; the policy
  retains ownership precisely because of that, which is why R1's misclassification
  matters (it reserves the slot for work that never existed).
- The historical HTTP 500's cause: **still unexplained**; nothing here is evidence
  about it.
- Backend interruption: unsupported, no endpoint added.
