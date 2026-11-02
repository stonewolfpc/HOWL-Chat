# Provisional freeze: image-generation ownership guarantees

Date: 2026-11-02
Status: **provisionally frozen** — the ownership core only, with findings R1–R6
retained as open work. This record freezes a *contract*, not a subsystem: the video
process lifecycle and the six findings below remain outside it.

Basis:
`build/logs/image-generation-ownership-followup-stage.md` (implementation),
`build/logs/image-generation-ownership-followup-review.md` (independent review,
verdict PASS WITH FINDINGS), `build/logs/image-cancellation-repair-review.md`
(where F1–F7 came from).

## What is frozen

Seven properties. Each is asserted by a test that exercises the behaviour, not a
helper, and each holds in the current source.

| # | Guarantee | Anchor |
|---|---|---|
| 1 | One generation request at a time per engine process; a second is refused, never overlapped | `generation.go` `beginGeneration`/`endGeneration`; `TestSecondGenerationIsRefusedWhileTheFirstIsActive` |
| 2 | Ownership is released only on a confirmed result, a pre-dispatch failure, or a confirmed child exit — never on an error or a stopped wait | `generation.go` `endGeneration`; `processConfirmedGone`; `TestOwnershipIsReleasedAfterAConfirmedProcessExit` |
| 3 | Closing the client request is reported as exactly that; no code claims the backend stopped | `describeGenError`; `TestTxt2ImgIsCancelledThroughItsRequestContext` |
| 4 | An ordinary failure is never reported as a client cancellation, directly or via wrapping | `classifyGeneration` + `settleGeneration`; the 500 / empty-result / transport / poll-timeout tests |
| 5 | Cancellation is routed by owner (video → load → generation); an empty id cancels nothing | `CancelJob`, `abortLoad`, `cancelGeneration`; `TestCancellingAGenerationIsNotSwallowedByAnotherJobsLoad`, `TestEmptyJobIdentityCannotCancelWork` |
| 6 | Progress ownership is tokenised: a late request clears only its own id | `client.setActiveJob`/`resetActiveJob`; `TestStaleProgressOwnerCannotClearANewerJob` |
| 7 | Unload, reload and process exit report their own reason, not a cancellation (single-request case) | `genClose` + `resetGeneration(reason)`; `TestShutdownIsNotReportedAsAClientCancellation` |

Evidence at freeze: 55 engine tests and 8 job tests passing, race detector clean on
both packages, `go vet` and the full suite green, both tagged builds green. The
review found no defect in guarantees 1–6; guarantee 7 holds in the tested case and
carries R2 below.

## What is explicitly not frozen

- **Video process ownership (F4)**: `generateVideo` still does not claim the
  generation slot; `stopServer` still kills sd-server unconditionally, calls
  `os.Process.Wait()` beside the reaper's `cmd.Wait()`, and reads `serverProc`
  without `procMu`. Untouched by design.
- **The lifecycle around the core**, i.e. findings R1–R6.

## Open findings carried forward

| Finding | Severity | Location | Disposition |
|---|---|---|---|
| R1 — img2img pre-dispatch errors not marked `notDispatched` | medium | `client.go:766-773` | fix in the next bounded change, with a test |
| R2 — close reason read after cleanup; a new request can rewrite an unload as a cancellation | low/medium | `settleGeneration` → `describeGenError` → `closeReason` | fix by carrying the reason on the generation handle |
| R3 — video cancel can be lost between `cmd.Start()` and owner assignment | low/medium | `video.go:337-344` | **carried into the video ownership stage as an acceptance criterion; must not be lost** |
| R4 — inner errors not unwrappable through the new sentinels | low | `notDispatched`, `describeGenError` | fix with R1/R2 |
| R5 — legacy non-cancellable callers are silent | low | `CancelJob`, `GenerateImage`, `SetActiveImageModel` | document or surface explicitly |
| R6 — sequential-batch local failure classified by the failing item | low | `client.generateSequentialBatch` | focused test |

## Standing constraints

- **Backend interruption remains unestablished.** The conservative refusal and the
  explicit unload/reload recovery are a safety guarantee: they prevent two requests
  overlapping work whose state is unknown. They are **not** evidence that
  cancellation stops computation, and no record here may be read that way.
- **The historical HTTP 500 remains unexplained.** Nothing frozen here establishes
  its cause.
- Only Windows has been built and tested. Nothing in the frozen core depends on
  Windows process behaviour; the video kill paths (F4) remain unverified elsewhere.

## Reopening rule

This core is reopened only to close a listed finding or a demonstrated defect in
guarantees 1–7. Cosmetic cleanup, broader refactoring, or rewording of the frozen
claims is not a reason to reopen it. A change that touches these paths re-tests the
affected behaviour; unrelated suites need not be re-run to restate what this freeze
already established.
