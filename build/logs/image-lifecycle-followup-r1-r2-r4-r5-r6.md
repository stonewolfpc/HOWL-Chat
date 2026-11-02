# Stage: lifecycle follow-up — R1, R2, R4, R5, R6

Date: 2026-11-02
Status: **implemented, unit-tested, build-verified** — pending independent review (not frozen)
Basis: `build/logs/image-generation-ownership-followup-review.md` (findings R1–R6),
`build/logs/image-ownership-provisional-freeze.md` (the frozen contract)
Scope honoured: no video code, no F4, no R3, no live run, no Git.

## Findings addressed

| Finding | Defect | Resolution |
|---|---|---|
| R1 | img2img's missing/unreadable init image returned plain errors, so a local input failure was classified `ambiguous` and reserved the generation slot for work that never existed | both wrapped with `notDispatched`, matching inpaint |
| R2 | `describeGenError` read engine-wide `genClose` **after** the slot was released, so a request that started in that window could rewrite an unload as a client cancellation | the close reason now lives on the generation handle and is captured while the slot is still held |
| R4 | `notDispatched` and the uncertain/stopped descriptions used `%v`, so the underlying file/HTTP/transport error was unreachable | `%w` throughout those wrappers |
| R5 | the legacy direct callers were silently non-cancellable | documented at the API boundary, and an empty `CancelImageJob` id now states why it cannot cancel |
| R6 | sequential-batch local failure was unreviewed and untested | two tests: first-item failure is `notDispatched` with nothing dispatched; later-item failure keeps the earlier result and frees the slot |

## Source and test files changed

| File | Change |
|---|---|
| `internal/image/engine/sdcpp/client.go` | `img2img`: `opts.InitImagePath == ""` and `os.ReadFile` failures wrapped with `notDispatched(...)` (lines ~766-773) |
| `internal/image/engine/sdcpp/generation.go` | `genHandle.close genClose` added; `beginGeneration` stores the handle on the engine; `endGeneration` returns the reason captured under `genMu`; `settleGeneration` passes it on; `describeGenError` is now `describeGenError(outcome, reason, err)` and no longer reads engine state; `resetGeneration(reason)` records the reason on the owned handle and clears `genHandle`; engine-wide `genClose` and `closeReason()` removed; `notDispatched` and every `describeGenError` branch wrap with `%w` |
| `internal/image/engine/sdcpp/engine.go` | `genClose genClose` field replaced by `genHandle *genHandle` |
| `app-image-generate.go` | `CancelImageJob` rejects a blank id with an explanatory error; doc comments on `CancelImageJob` and `GenerateImage` state that the direct path is not cancellable and why |
| `app-image-model.go` | doc comment on `SetActiveImageModel` states the same limitation |
| `internal/image/engine/sdcpp/followup_test.go` (**new**) | 7 tests (R1 contract, R2 mechanism ×2, R4 ×2, R6 ×2) |

No change to: the slot state machine, the release policy, dispatch or cancellation
semantics, process supervision, the reaper, or any video code.

## R1 — shared pre-dispatch contract

`TestPreDispatchInputFailuresShareOneContract` runs four equivalent failures against
one engine and asserts the same contract for each: `errors.Is(err, errNotDispatched)`,
**not** `errEngineUncertain`, **not** `errClientCancelled`, and
`GenerationState() == "idle"` afterwards.

```
--- PASS: TestPreDispatchInputFailuresShareOneContract (0.00s)
    img2img without an init image
    img2img with an unreadable init image
    inpaint without an init image
    inpaint without a mask
```

This is the test the earlier suite lacked: it asserts ownership and classification,
not a message substring, so an equivalent path cannot drift again silently.

## R2 — mechanism and regression test

Mechanism: the reason is written under `genMu` onto the handle being closed
(`resetGeneration`), read under `genMu` in `endGeneration` **before** the slot is
released, and returned to `settleGeneration` as a value. `describeGenError` takes it
as a parameter, so no engine-wide mutable state is read after release.

`TestSettledGenerationKeepsItsOwnCloseReason` reproduces the exact window the review
described, deterministically:

1. A begins; `resetGeneration(closeForShutdown)` closes A and frees the slot.
2. B begins — the interleaving that used to rewrite A's reason.
3. A settles with a transport-ish error.

Assertion: A's error satisfies `errEngineStopped` and **not** `errClientCancelled`.
Against the previous implementation this test fails at step 3, because B's
`beginGeneration` reset the engine-wide reason to `closeNone`.

`TestClientCancellationStaysDistinctFromAnEngineStop` guards the opposite direction:
an owner cancel plus a settled request still reports `errClientCancelled` and not
`errEngineStopped`.

**Shutdown classification preserved:** `TestShutdownIsNotReportedAsAClientCancellation`
was not modified and still asserts `errors.Is(err, errEngineStopped)`.

## R4 — unwrapping

`notDispatched(err)` is now `fmt.Errorf("%w: %w", errNotDispatched, err)`, and every
`describeGenError` branch wraps the inner error with `%w` as well (two `%w` verbs are
supported by `go 1.23.8`, this module's toolchain).

- `TestUncertainErrorKeepsTheTransportError`: a dead listener still yields
  `errors.Is(err, errEngineUncertain)` **and** `errors.As(err, *net.OpError)`.
- `TestNotDispatchedErrorKeepsTheFileError`: inpaint with a missing mask still
  yields `errors.Is(err, errNotDispatched)` **and** `errors.As(err, *os.PathError)`.

Wrapping does not change classification: `classifyGeneration` tests the sentinels,
and the tests above assert `errEngineUncertain` is absent where it should be.

## R5 — legacy non-cancellable callers

Remaining non-cancellable operations and why:

| Caller | Operation | Why it cannot be cancelled |
|---|---|---|
| `App.GenerateImage` | direct synchronous generation | carries no job identity; nothing addresses a cancel to it |
| `App.SetActiveImageModel` | direct model load | same; a cancellable load runs inside a job via `StartImageJob` |

How the limitation is surfaced:
- doc comments on both methods state it at the point of use;
- `CancelImageJob("")` now returns
  `"CancelImageJob needs a job id: generation started directly with GenerateImage cannot be cancelled"` instead of "job not found", so an empty id is explained rather than silently ignored;
- no synthetic job id was invented, and no behavior of the two callers changed.

## R6 — batch behaviour

- `TestBatchFirstItemLocalFailureIsNotDispatched`: a `BatchCount: 2` img2img whose
  first item cannot read its init image returns `errNotDispatched` (it survives the
  batch's `"batch item %d/%d failed: %w"` wrapper), **zero** HTTP requests reach the
  server, and the slot is `idle`.
- `TestBatchLaterLocalFailureKeepsTheEarlierResultAndTheSlot`: item 1 succeeds; the
  server deletes the init image while answering it; item 2 fails locally. The batch
  returns the first item's image with no error, exactly one generation was
  dispatched, and the slot is `idle` — not abandoned for work that never existed.

Batch semantics were **not** changed. The existing behaviour (keep partial results,
log a later-item interruption) is now pinned by tests so a future change to it is
visible rather than accidental.

## Verification (actual commands and results)

```
go build ./...                                        exit 0
go test ./internal/image/engine/sdcpp/ ./internal/image/job/   ok, ok (62 and 8 tests)
go test -race -count=1 <same two packages>            ok, ok (clean)
go vet ./...                                          exit 0
go test -count=1 ./...                                exit 0, 17 packages ok, 0 FAIL
  (log: build/logs/full-suite-r1r2.log)
CGO_ENABLED=0 go build -tags "desktop,production" ... exit 0
go build -tags uiharness ...                          exit 0
gofmt -l <changed dirs/files>                         empty (clean)
```

New tests, all passing:

```
--- PASS: TestPreDispatchInputFailuresShareOneContract
--- PASS: TestSettledGenerationKeepsItsOwnCloseReason
--- PASS: TestClientCancellationStaysDistinctFromAnEngineStop
--- PASS: TestUncertainErrorKeepsTheTransportError
--- PASS: TestNotDispatchedErrorKeepsTheFileError
--- PASS: TestBatchFirstItemLocalFailureIsNotDispatched
--- PASS: TestBatchLaterLocalFailureKeepsTheEarlierResultAndTheSlot
```

Correction to an earlier record: the follow-up stage record says "18 packages"; the
suite actually reports **17** `ok` packages (the miscount came from counting the
full-suite block by eye). No result changes — the suite is green either way.

No live generation, harness run, UI exercise or backend probe was performed; none
was authorized for this stage, so no new live evidence is claimed.

## Limitations

- The conservative post-dispatch policy is unchanged: an HTTP error status, empty
  result, transport loss or poll timeout still leaves the engine unavailable until
  unload/reload. R1 only removes the *local* failures from that set.
- R3 (video cancel lost between `cmd.Start()` and owner assignment) and F4 (video
  bypasses the slot; `stopServer` waits beside the reaper) remain open and untouched.
- `GenerateImage`/`SetActiveImageModel` remain non-cancellable by design.
- Only Windows was built and tested; nothing new depends on platform process
  behaviour.

## Frozen-boundary confirmation

- Guarantees 1–6 of `image-ownership-provisional-freeze.md` are unchanged in
  behaviour and structure: no change to `beginGeneration`'s claim logic,
  `endGeneration`'s release policy, `cancelGeneration`'s identity rules, the
  `CancelJob` routing order, the progress token, or supervision.
- Only F7's *mechanism* changed, as authorized: the reason moved from engine state
  to the generation handle. Its guarantee and its test are intact.
- `internal/image/engine/sdcpp/video.go` was not modified in this stage; F4 and R3
  remain as recorded.
- The provisional freeze record and all earlier records were not edited.
