# Review: image lifecycle follow-up R1/R2/R4/R5/R6 (read-only)

Date: 2026-11-02
Mode: **read-only**. No code, tests, config or records were modified, and no
build/test/run was executed. Claims come from reading the current source and the
saved suite log, not from the stage record's prose.

Target: `build/logs/image-lifecycle-followup-r1-r2-r4-r5-r6.md`
Contract: `build/logs/image-ownership-provisional-freeze.md`

**Verdict: PASS WITH FINDINGS — the authorized repair is correct and R1–R6 are
addressed.** One latent weakness exists in the settle path that R2 touched
(harness Review-1), plus two documentation/observability notes. None invalidates
the frozen guarantees or the findings as repaired; Review-1 should be closed the
next time this mechanism is touched.

## Findings

### Review-1 — LOW (latent): a settled handle can be settled again and re-transition the slot

**Source-proven** (`generation.go`, `endGeneration` and `settleGeneration`).

`endGeneration` clears `e.genHandle` and `e.genCancel`, but neither invalidates
`h.token` nor advances `e.genSeq`. The guard is `if e.genSeq == h.token`, and the
sequence only moves in `beginGeneration` and `resetGeneration`.

Execution path:

1. `beginGeneration("A")` → `genSeq = N`, `h.token = N`.
2. `settleGeneration(h, nil)` → completed → state `idle`; `genHandle` cleared;
   `genSeq` still `N`.
3. `endGeneration(h, genOutcomeAmbiguous)` called again with the same handle →
   `e.genSeq == h.token` is **still true** → `retain && !gone` → state
   `abandoned`, owner `"A"`.

Consequence: a second settlement of the same generation marks a free engine
unavailable, which then requires an unload. No production path settles twice today
(each `Generate*` calls `settleGeneration` exactly once), so this is latent rather
than live. It matters because the guard's original meaning — "this handle still
identifies the current owner" — became false once settlement stopped moving the
sequence.

Required correction: invalidate the handle on settlement (`h.token = 0`) or compare
`e.genHandle == h` instead of the sequence. A test that settles the same handle
twice and asserts the state stays `idle` would pin it.

### Review-2 — LOW: R5 is documented at the definition, not surfaced to a direct caller

**Source-proven** (`app-image-generate.go` `GenerateImage`, `CancelImageJob`;
`app-image-model.go` `SetActiveImageModel`).

The doc comments state the limitation accurately, and `CancelImageJob("")` now
returns an explicit message. But a caller that invokes `GenerateImage` or
`SetActiveImageModel` gets no runtime indication that the operation cannot be
cancelled — the explanation is only visible in Go source, not through the binding.
The review asked for communication "at the actual caller/API boundary"; this is
partial.

Consequence: a user on the direct path can believe the Cancel button applies. No
correctness impact: an empty-id cancel is inert and cannot touch another operation
(verified in `cancelGeneration`/`abortLoad`/`CancelJob`).

Correction: surface it in the UI copy for the direct path, or note it as an
accepted documentation-only closure. Do not invent a job id.

### Review-3 — INFO: R4 makes `context.Canceled` reachable from closed outcomes

`describeGenError` now wraps the inner error with `%w` in every cancellation branch,
so a shutdown-closed request satisfies both `errEngineStopped` and
`errors.Is(err, context.Canceled)`. That is accurate (the context was cancelled),
and no caller inspects it — `job.finish` stores `err.Error()` only. Recorded because
it is a real observable change that the stage record does not mention; no correction
warranted.

### Review-4 — INFO: optional img2img inputs still fail silently

`client.img2img` ignores read errors for `ControlImagePath` and `MaskPath`
(`if data, err := os.ReadFile(...); err == nil`). That is not a classification
inconsistency — no error is produced at all — so R1's contract does not cover it,
but a user can get an image generated without a mask they thought they supplied.
Pre-existing and outside this stage; noted for the record only.

## R1 — shared pre-dispatch contract: PASS

Verified in source:

- `client.go:766-773` now wraps both img2img failures with `notDispatched(...)`,
  matching inpaint's four checks (`InitImagePath`, `MaskPath`, init read, mask read)
  and both marshal/build paths in each function. No equivalent pre-dispatch failure
  in `txt2img`, `img2img` or `inpaint` remains unwrapped.
- `notDispatched` uses `%w: %w`, so `classifyGeneration` still finds
  `errNotDispatched` and the inner error survives.
- The test (`followup_test.go:23`) drives the production `GenerateImg2Img` /
  `GenerateInpaint` entry points (not helpers), covers all four cases, and asserts
  `errNotDispatched`, absence of `errEngineUncertain` and `errClientCancelled`, and
  `GenerationState() == "idle"`.

It would have failed before R1 for the two img2img cases (state would be
`abandoned`), which is the behaviour that needed pinning.

## R2 — generation-specific close reason: PASS (with Review-1)

- Reason initialisation: `genHandle` is created only in `beginGeneration`, with
  `close: closeNone`. No other construction path exists.
- Capture before release: `endGeneration` reads `reason := h.close` under `genMu`
  **before** transitioning and clears the handle; it returns the value.
- Stability after release: `describeGenError(outcome, reason, err)` is a free
  function taking the value; engine-wide close state and `closeReason()` are gone
  (grep confirms no `genClose genClose` field and no `closeReason` method).
- A subsequent generation cannot rewrite an earlier reason: proven by construction
  and by `TestSettledGenerationKeepsItsOwnCloseReason` (`followup_test.go:80`), which
  performs exactly the old race (reset A → begin B → settle A) and asserts
  `errEngineStopped`, not `errClientCancelled`. Against the previous implementation
  B's `beginGeneration` reset the engine-wide reason to `closeNone`, so this test
  would have failed — it is a genuine regression test.
- Reverse direction: `TestClientCancellationStaysDistinctFromAnEngineStop`
  (`:110`) keeps a user cancel as `errClientCancelled` and not `errEngineStopped`.
- Shutdown assertion preserved: `TestShutdownIsNotReportedAsAClientCancellation`
  (`ownership_test.go:259`) is unchanged and still asserts
  `errors.Is(err, errEngineStopped)`.
- Locking: `h.close` is written under `genMu` (`resetGeneration`) and read under
  `genMu` (`endGeneration`); no new nesting, `processConfirmedGone()` is still
  called before taking `genMu`, and `cancel()` is still invoked after unlocking.
  No lock-order cycle was introduced.

Residual weakness: Review-1 (double settlement).

## R4 — unwrapping: PASS

- Sentinels remain discoverable: `errors.Is(err, errNotDispatched)`,
  `errEngineUncertain`, `errClientCancelled`, `errEngineStopped` all still match
  because each wrapper carries its own `%w`.
- Underlying errors remain inspectable: `TestUncertainErrorKeepsTheTransportError`
  (`:131`) reaches `*net.OpError`; `TestNotDispatchedErrorKeepsTheFileError`
  (`:151`) reaches `*os.PathError` through both wrappers.
- Classification did not shift: `classifyGeneration` keys on `errNotDispatched`
  only, and the tests assert `errEngineUncertain` is absent where it should be.
- No caller depended on the old formatting: the job layer stores `err.Error()`, and
  no source checks for the previous `%v` text shape (grep of `errClientCancelled`
  users shows tests and the job layer only).

## R5 — legacy non-cancellable callers: PASS WITH THE Review-2 NOTE

- Empty identity is inert end to end: `CancelImageJob` rejects a blank id before
  the manager is touched; `Manager.Cancel` cannot reach the engine with `""`;
  `CancelJob("")` returns nil without routing; `cancelGeneration("")` and
  `abortLoad("")` return false.
- No other operation can be targeted: routing requires `genOwner == jobID` or
  `loadOwner == jobID`.
- The limitation is stated at both definitions; the cancel boundary explains itself
  at runtime. `GenerateImage` itself does not (Review-2).
- No synthetic identity was introduced.

## R6 — sequential batch: PASS

- Production path exercised: both tests call `GenerateImg2Img` with
  `BatchCount: 2`, which enters `generateSequentialBatch` and the real HTTP client.
- First-item failure: `TestBatchFirstItemLocalFailureIsNotDispatched` (`:178`)
  asserts `errNotDispatched` survives `"batch item %d/%d failed: %w"`, **zero**
  requests reached the server, and the slot is `idle`.
- Later-item failure: `TestBatchLaterLocalFailureKeepsTheEarlierResultAndTheSlot`
  (`:204`) removes the init image inside the first `img_gen` handler, so item 2
  fails at `os.ReadFile`. Asserts the first item's image is returned, exactly one
  generation was dispatched, and the slot is `idle` — not abandoned.
- Item numbering and wrapping in `generateSequentialBatch` are unchanged and
  correct (`i+1`, `opts.BatchCount`); the tests depend on that wrapper, so a silent
  change to it would break them.

Residual product note: a later-item failure is still only logged, not reported to
the caller (Review-4 adjacent). Semantics were correctly left alone.

## Freeze boundary: PASS

- Slot state machine and release policy: unchanged. Same `retain` set
  (`cancelled`, `ambiguous`), same confirmed-exit release, same self-heal in
  `beginGeneration`.
- Cancellation routing: unchanged this stage (`CancelJob` still video → load →
  generation with ownership checks).
- Process-exit confirmation: unchanged (`processConfirmedGone` on the reaper
  channel, used in the same two places).
- Progress token ownership: untouched (`setActiveJob`/`resetActiveJob(token)`).
- Job identity: only the authorized R5 clarification (`CancelImageJob` blank-id
  rejection); `JobID` propagation and `cancelGeneration` rules unchanged.
- Batch behaviour: no semantic change.
- Video: `video.go` was not edited in this stage (file mtime 18:47, before the
  stage's edits at ~18:56–19:00); `cancelVideoOwned`, `runVideoCLI(..., owner)` and
  the untouched `stopServer`/`os.Process.Wait` remain as the F4 record describes.
  F4 and R3 stay deferred.

## Verification record accuracy

Independently checkable:

- **62 engine tests**: `grep -c '^func Test'` over the package's `*_test.go` returns
  62. Consistent.
- **8 job tests**: `job_test.go` contains 8 test functions. Consistent.
- **17 packages, exit 0**: `build/logs/full-suite-r1r2.log` contains 17 `ok` lines,
  no `FAIL`, and only `ok`/`?` package lines. Consistent, and the record's
  correction of the earlier "18" is right.
- **New tests present and named as reported**: verified in `followup_test.go`
  (7 tests at lines 23, 80, 110, 131, 151, 178, 204).

Reported-only (not executable under this read-only review): the two `-race` runs,
`go vet`, the two tagged builds, and the gofmt check. They are internally consistent
with the source (no unused fields or imports are visible, and the record's command
list matches the project's conventions), but this review treats them as **reported
evidence**, not independently confirmed.

## Remaining unproven

- Backend interruption: still unestablished. Nothing in this stage touches it, and
  no claim is made.
- The historical HTTP 500: still unexplained.
- Live behaviour of this stage: not exercised (not authorized). The conservative
  post-dispatch refusal policy is unchanged.
- Cross-platform behaviour: only Windows was built; nothing here adds a platform
  dependency, and the video kill paths remain unverified elsewhere.
