# Stage — Generation Settlement Guard (Review-1)

Scope: Review-1 only, from `image-lifecycle-followup-r1-r2-r4-r5-r6-review.md`.
Date of work: 2026-11-01. Nothing was committed; no Git write operations were performed.

## Defect being closed

`endGeneration` cleared `genHandle` and `genCancel` but left the handle's token equal
to `genSeq`, and `genSeq` only moves in `beginGeneration` / `resetGeneration`. A second
settlement of an already-settled handle therefore still matched `e.genSeq == h.token`
and could move an idle engine to `abandoned`.

No production path settles a handle twice today, so the defect was latent, not live.

## Change (exact)

File: `internal/image/engine/sdcpp/generation.go` (only production file touched).

1. `genHandle` gains `settled bool`, guarded by `Engine.genMu` — "this generation has
   already reported its outcome".
2. `endGeneration`, inside the existing `e.genMu` critical section, sets `h.token = 0`
   and `h.settled = true` when it performs the transition. Token invalidation and the
   flag therefore use the same lock that decides settlement; no new lock, no new state
   machine edge.
3. New sentinel `errGenerationSettled = errors.New("generation was already settled")`.
4. `settleGeneration` returns that sentinel immediately for an already-settled handle
   instead of classifying and describing again.

Chosen approach: invalidate the token *and* record settlement. The token alone was
insufficient in practice — see "First attempt and what it exposed" below.

Explicitly not changed: admission and release policy, cancellation routing,
`processConfirmedGone` / `resetGeneration`, progress-token ownership, close-reason
capture, `describeGenError` classification branches, batch semantics, video code.

## First attempt and what it exposed

The minimal token-only change (`h.token = 0`) already kept the slot idle, but the new
regression test failed on the error the second settlement reported:

```
settlement_guard_test.go:37: want errEngineUncertain from the second settlement,
got client request cancelled: sd-server may still be running the abandoned work, so
the engine stays unavailable until it is restarted (late failure)
```

Cause: `endGeneration` calls `h.cancel()` at the end of the first settlement, so on a
second call `classifyGeneration(h.ctx.Err() != nil, err)` sees a cancelled context and
labels an ordinary late error a client cancellation — exactly the mislabelling F1
removed. The state was already safe; the report was wrong. The `settled` flag plus the
sentinel close that, and the final code reports `errGenerationSettled` instead.

This is recorded because the first implementation was not sufficient and the failing
assertion is the evidence that drove the final shape.

## Regression test

New file: `internal/image/engine/sdcpp/settlement_guard_test.go` (2 test functions).

`TestSettledHandleCannotSettleTwice`:
1. `beginGeneration("job1")` claims the slot.
2. `settleGeneration(h, nil)` settles successfully; asserts state `idle`.
3. `settleGeneration(h, errors.New("late failure"))` — the defect path; asserts state
   is still `idle` and the error `errors.Is(..., errGenerationSettled)`.
4. `beginGeneration("job2")` proves the slot is genuinely free, and settling it leaves
   the engine `idle`.

`TestFirstSettlementStillClassifiesEachOutcome` (table, 3 subtests) proves a legitimate
first settlement is unchanged:

| outcome | setup | state after | error |
|---|---|---|---|
| completed | settle with nil | `idle` | nil |
| cancelled | `cancelGeneration("job1")` then settle | `abandoned` | `errClientCancelled` |
| ambiguous | settle with a plain error | `abandoned` | `errEngineUncertain` |

No existing test was modified.

## Verification (actual commands, actual results)

Working directory `C:\Users\StoneWolfPC\HOWL_Chat`.

| command | result |
|---|---|
| `gofmt -w internal/image/engine/sdcpp/generation.go internal/image/engine/sdcpp/settlement_guard_test.go` | applied |
| `go build ./...` | exit 0, no output |
| `go test -count=1 -v -run 'TestSettledHandleCannotSettleTwice\|TestFirstSettlementStillClassifiesEachOutcome' ./internal/image/engine/sdcpp/` | `--- PASS` for both; all 3 subtests PASS |
| `go test -count=1 ./internal/image/engine/sdcpp/` | `ok howl-chat/internal/image/engine/sdcpp 1.695s` |
| `go test -race -count=1 ./internal/image/engine/sdcpp/` | `ok ... 3.262s`, no data race |
| `go vet ./...` | exit 0 |
| `gofmt -l internal/image/engine/sdcpp` | empty (clean) |

Test counts: `internal/image/engine/sdcpp` now has **64** test functions (62 before +
2 new, grep of `^func Test` across `*_test.go`); `internal/image/job` unchanged at 8.

Not run, by instruction: live generation, harness/UI testing, tagged builds. The
`-race` run covers the new lock interaction, which is the only concurrency surface added.

## Frozen boundary confirmation

Guarantees 1–7 of `image-ownership-provisional-freeze.md` are intact:

- **Admission/release:** `beginGeneration` untouched; release still only on confirmed
  result, pre-dispatch failure, or confirmed child exit; retain/`gone` conditions in
  `endGeneration` unchanged.
- **Cancellation routing:** `cancelGeneration`, `CancelJob`, `CancelImageJob`, handle
  `cancel` untouched.
- **Process-exit confirmation:** `processConfirmedGone`, `stopServer` untouched.
- **Progress-token ownership:** progress tokens and `resetActiveJob(token)` untouched.
- **Close-reason capture:** `h.close` set and read exactly as before; the new field is a
  separate boolean.
- **Error classification:** `describeGenError` and `classifyGeneration` unchanged; only a
  new pre-check short-circuits an already-settled handle.
- **Batch semantics:** untouched; the 62 pre-existing tests, including both batch tests
  and the shutdown/`errEngineStopped` test, still pass.

Video code: `internal/image/engine/sdcpp/video.go` mtime `2026-11-01_18:47`, unchanged
through this stage (F4 and R3 remain deferred).

Records untouched: `image-lifecycle-followup-r1-r2-r4-r5-r6-review.md` (19:02) and
`image-ownership-provisional-freeze.md` (18:55) both predate this stage's edits to
`generation.go` and the new test file (19:08).

## Deliberately out of scope

Review-2 (R5 surfaced only at the definition), Review-3 (`errors.Is` observable change),
Review-4 (silent ignore of unreadable img2img control/mask paths) were not addressed.

## Residual risk

- The new path is exercised only by unit tests; no live generation was run, per
  instruction.
- `errGenerationSettled` is currently unreachable in production, by design: it exists to
  make a future double-settlement fail visibly rather than corrupt the slot.
- Backend interruption and the historical 500 remain unestablished, as before.
