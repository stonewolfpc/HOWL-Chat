# Review: image-generation cancellation repair (read-only)

Date: 2026-11-02
Mode: **read-only**. No product code, tests, config, artifacts or records were
modified. No builds, tests, server requests, generations or process actions were
run. Claims below come from reading current source, not from the repair author's
summary.

Verdict: **PASS WITH FINDINGS — not freeze-ready.** The ownership core is sound,
but two defects at the same boundary make the stage record's headline guarantee
("no second generation can start on that process while the backend's state is
unknown") false as written. F1 and F2 must be corrected before freezing; F3/F4 are
pre-existing limits the repair did not introduce and should be recorded as
follow-ups, not folded into this stage.

## Evidence key

Observed · Source-derived · Test-derived · Inferred · Unestablished · Superseded.

## Findings, by severity

### F1 — HIGH: every failed generation is reported as a client cancellation

**Source-derived, confirmed by reading the call order.**

- `internal/image/engine/sdcpp/engine.go:296` → `e.endGeneration(h)`, then
  `:298` → `classifyGenError(h, err)` (txt2img; identical at `:318`/`:320` for
  img2img and `:338`/`:340` for inpaint).
- `internal/image/engine/sdcpp/generation.go:116` — `endGeneration` ends with an
  unconditional `h.cancel()`.
- `generation.go:103` — `abandoned := h.ctx.Err() != nil` is read *before* that
  cancel, so the slot transition is right, but the handle's context is cancelled
  by the time `classifyGenError` runs.
- `classifyGenError` therefore takes its `h.ctx.Err() != nil` branch for **every**
  error, not only cancellations.

Consequence: an HTTP 500, `"no images in txt2img response"`, a connection refusal,
a JSON decode failure, or an async `job poll timeout` all surface as
`client request cancelled: sd-server may still be running the abandoned work…`.
The user is told they cancelled when they did not, and `errors.Is(err,
errClientCancelled)` is meaningless to any caller. This contradicts the repair
record's claim that `classifyGenError` "turns a request error into an accurate
statement", and it means requirement 2 ("report this accurately as client-request
cancellation") is met only in the cancel case and inverted in every other.

The defect is invisible to the new tests: all six engine tests exercise either a
clean completion or a real cancellation. No test drives a non-cancel failure
through `GenerateTxt2Img`.

### F2 — HIGH: a non-cancel failure releases the slot while the backend may still be working

**Source-derived.**

`endGeneration` (`generation.go:96-120`) sets `genState = genIdle` whenever the
request returned with `h.ctx.Err() == nil`. The error-returning, non-cancelled
paths on a real engine include:

| Path | Location | Backend state after return |
|---|---|---|
| txt2img HTTP 500 | `client.go` ~`558` | unconfirmed — the server answered with an error body |
| txt2img decoded-but-empty | `client.go` ~`566` | unconfirmed |
| async `pollJob` poll timeout | `client.go:618+` | job may still be running |
| transport error (reset/refused) mid-request | `client.go` ~`556` | unconfirmed |
| `no images in job result` after recovery | `client.go` pollJob completed branch | unconfirmed |

In those cases the next `beginGeneration` is admitted. That is exactly the overlap
state the audit describes, and the live 500 in
`issue-cancelled-generation-reuse.md` was a server-side error *after*
`generate_image completed` — i.e. an F2-shaped return.

So the guarantee holds for **cancellation** (the case the stage set out to fix) but
not for **failure**. The record's verdict sentence overstates it, and the audit's
Option C ("two requests can never overlap") is not yet fully true.

### F3 — MEDIUM: `CancelJob` routing drops a generation cancel, and video cancel has no owner

**Source-derived.** `engine.go:389-399`:

```
if e.cancelVideo() { return nil }              // global: whichever render is running
if e.loadInProgress() { e.abortLoad(jobID); return nil }
e.cancelGeneration(jobID)
```

- If job B is loading (or queued behind `loadRun`) while job A is generating,
  `CancelJob(jobA)` takes the load branch. `abortLoad("jobA")` is a no-op because
  the load owner is B, and the function returns without ever reaching
  `cancelGeneration` — A's cancel is silently lost.
- `cancelVideo()` returns true for *any* active sd-cli render, so cancelling job A
  kills job B's video and skips A's own cancellation.

The repair scoped `abortLoad` and `cancelGeneration` by owner, but not the routing
that chooses between them. Stage 3's review flagged the same class
("cancellation is engine-global"); this stage narrowed it without closing it.

### F4 — MEDIUM: video generation bypasses the ownership boundary and breaks single-`Wait` ownership

**Source-derived, pre-existing code the repair did not touch.**

- `video.go:139-181`: `generateVideo` never calls `beginGeneration`; it takes
  `runMu` only.
- `video.go:181`: it calls `stopServer()`, which kills the sd-server child
  unconditionally — including while an image generation is active on it, or while
  its slot is `abandoned`. A video job can therefore terminate the engine of an
  unrelated image job, and the abandoned-slot refusal does not cover video at all.
- `video.go:116-124`: `stopServer` reads and writes `e.serverProc` **without**
  `procMu` (racy against `getServerProc`/`setServerProc`/`Shutdown`), and calls
  `os.Process.Wait()` (`:122`) while `startProcessReaper` owns `cmd.Wait()`.
  Two waits on one child can make the reaper record a wait error instead of the
  real exit — and `processConfirmedGone()` plus `beginGeneration`'s self-heal both
  depend on that reaper signal.

Consequence for this stage: the repair's statement "`cmd.Wait()` ownership is
unchanged" is true, and its ownership logic does not depend on the no-op
`KillTree`, but the repository as a whole does **not** currently satisfy "exactly
one component owns `cmd.Wait()`". That is a review finding about the system, not a
regression introduced here.

### F5 — LOW/MEDIUM: an empty `JobID` still means "whoever is active"

**Source-derived.**

- `app-image-generate.go:66-78` — `GenerateImage` is a direct, non-job entry point;
  it calls `GenerateImg2Img`/`GenerateInpaint`/`GenerateVideo`/`GenerateTxt2Img`
  with `JobID == ""`, so the slot owner is `""`.
- `app-image-model.go:131` — `SetActiveImageModel` → `LoadModel` with
  `ModelLoadOptions.JobID == ""` → `beginLoad("")`.
- `cancelGeneration("")` and `abortLoad("")` would match those owners.

No current caller sends an empty id — `Manager.Cancel` always passes `job.ID`
(`job.go:113-119`) — so this is latent, not live. It nevertheless leaves the
audit's requirement ("empty ids must not mean cancel whichever job is active")
unenforced. `GenerateImage` also has no cancellation path at all, and if invoked
while a job generates it is refused with `errEngineBusy`.

### F6 — LOW: an abandoned async request can clear the next job's progress id

**Source-derived.** `client.go:893` and `:1088` set `activeJobID`; the matching
`resetActiveJob()` calls sit just after `pollJob` returns. A request abandoned and
still unwinding when a child-exit reset admits the next job can clear that job's
id, degrading progress to the job-agnostic endpoint (`progress()`,
`client.go:1146+`). Progress display only; no ownership effect.

### F7 — LOW: an unload reports itself as a client cancellation

`Shutdown` → `resetGeneration()` cancels the in-flight request's context, so the
job reports `errClientCancelled` ("…abandoned work…") for what the user did by
unloading. Message accuracy only, and a sub-case of F1.

## Ownership transition analysis — core PASS

| Transition | Implementation | Verdict |
|---|---|---|
| `idle → active` | `beginGeneration` claims under `genMu`, bumps `genSeq`, returns token (`generation.go:66-95`) | correct |
| `active → idle` | `endGeneration`, token match, no ctx cancellation (`generation.go:96-120`) | correct for clean completion; see F2 for error returns |
| `active → abandoned` | same, `h.ctx.Err() != nil` | correct; release is deliberately withheld |
| `abandoned → idle` on confirmed exit | `beginGeneration` self-heal via `processConfirmedGone()` on the reaper's `procDone` | correct, platform-neutral |
| `→ idle` after `LoadModel` kills the previous child | `engine.go:150-155` | correct order (kill, reap, reset) |
| `→ idle` after `Shutdown` | reset after kill + `reapStoppedProcess` (`engine.go:400-415`) | correct (the ordering hole found during implementation is closed) |
| stale handle writing after reset | `genSeq++` in `resetGeneration` invalidates the token | correct |

All three HTTP generation paths (`engine.go:282`, `:306`, `:326`) claim the slot.
Sequential batches run inside a single claim (`client.go` `txt2img` →
`generateSequentialBatch`), so a batch cannot interleave with another request.
Video is the exception, and it is a process-level bypass rather than an HTTP one
(F4).

## Lock-order analysis — PASS, no cycle found

Acquisition order observed:

```
loadRun  →  genMu      (LoadModel → resetGeneration)
loadRun  →  procMu     (LoadModel → get/setServerProc, reapStoppedProcess)
loadMu   →  procMu     (abortLoad: lock released before Kill)
genMu, procMu          never held together (generationBusy() returns first)
videoMu                leaf (recordVideoProgress, videoProgressInfo, cancelVideo)
m.mu (job manager)     released before any engine call (job.go:113-119)
```

Checked interleavings: no goroutine waits for a request while holding a lock that
cancellation needs (`cancelGeneration` takes only `genMu` briefly); cancellation
never waits for process exit while holding a lock the reaper needs
(`abortLoad` releases `loadMu` before `Kill`; `reapStoppedProcess` is only called
by `LoadModel`/`Shutdown`, not from a cancel path); unload races are serialized by
the token bump; no path takes `loadRun` while holding `genMu`, so load and
generation ownership cannot deadlock against each other. `loadRun` is held across
a multi-minute load, but nothing cancellation needs sits behind it.

Residual risk worth a comment, not a change: `beginGeneration` → `resetGeneration`
→ `cancel()` runs context callbacks while no lock is held, which is correct today
and would break if a future cancel callback re-entered `genMu`.

## Job-identity gaps

Present and correct: `job.go:113-119` (id captured under `m.mu`, engine call
outside it), `job.go:237-239` (`opts.JobID = job.ID`), `app-image-model.go:74`
(id set on the per-request copy, not on `lastImageModelOpts`).

Gaps: F5 (empty ids are accepted as owners; two direct entry points use them), F3
(routing ignores ownership when a load is in progress; video cancel has no owner).

## HTTP cancellation semantics — PASS, with F1 as the stated defect

`NewRequestWithContext` is used on `txt2img` (`client.go` ~`547`), `pollJob`
(`:604+`), and both `img_gen` posts; `pollJob` selects on `ctx.Done()`/deadline/
ticker. Distinctions hold in the code that makes claims: a closed request is
reported as exactly that (`errClientCancelled`), the only backend-stop claim is
`pollJob`'s server-reported `"cancelled"` status, and nothing infers a stopped
backend from a disconnect. A late response cannot be attributed to a later job's
*result* (each call reads its own body); F6 is the progress-id variant.

## Recovery and UI behaviour

Traced: cancel with the backend alive → job `cancelled`, slot `abandoned`; retry →
`errEngineUnavailable`, no process spawned; unload → kill + reap + reset; reload →
fresh child, generation succeeds; process exit while abandoned → self-heal on the
next claim. The refusal cannot be bypassed through the job manager, and the direct
`GenerateImage` path is refused by the same slot. The record's live sequence
matches this, and its own `Aborted` / "unavailable until it is restarted" /
`UnloadImageModel ok 992ms` / `generate_image completed in 105.29s` lines are
consistent with the transitions above (Observed in the stage record; not re-run
here).

## Portability

- Confirmed child exit is read from the reaper's channel, not from a platform
  helper; `proc.KillTree` is not the relied-upon stop mechanism.
- `proc/job_other.go` remains a set of no-ops off Windows; `Shutdown` and
  `cancelVideo` still call `KillTree`, but each also calls `Process.Kill`, so the
  stop itself is portable (pre-existing).
- No hardware, VRAM, model, drive or duration assumption was introduced; refusal
  is state-based, not time-based. Only Windows has been exercised, and the record
  says so.
- F4's `os.Process.Wait()` in `stopServer` is platform-neutral code with a
  Windows-independent defect.

## Corrections needed before freezing

1. **F1** — read the cancellation state before `endGeneration` (capture
   `cancelled := h.ctx.Err() != nil`, or classify before ending), so a real error
   is reported as a real error.
2. **F2** — decide the slot policy for error-returning requests. Either release
   only on confirmed clean completion and mark unconfirmed errors
   `abandoned`/unavailable, or introduce a typed "backend state unknown" class for
   the 500 / poll-timeout / transport paths. As written, an error silently returns
   the engine to `idle`.
3. **Tests for both** — a server 500 must not produce `errClientCancelled`, and the
   slot's state after that error must be asserted. A `CancelJob` for a generating
   job while another job loads must still cancel that job.
4. **F3** — route cancellation by ownership instead of by "what is currently
   running"; video needs its own owner id.
5. **Record wording** — the repair record's verdict sentence needs narrowing in a
   later stage (this review does not edit it): the guarantee is cancellation-safe,
   not yet failure-safe.

Follow-ups, not blockers for this stage: F4 (video slot + reaper ownership),
F5 (empty-owner ids), F6 (progress id), F7 (unload message).

## What remains unproven

- **Backend interruption is still unestablished.** No interrupt/DELETE path exists
  in the client, and nothing in this review establishes that this sd-server build
  honours one, or that it stops work on client disconnect.
- **The 500 is still unexplained.** The repair removes one overlap mechanism for
  the cancel case (F2 leaves the failure case partly open); it does not establish
  that overlap produced that 500. The earlier record's explanation is not rewritten
  and is not upgraded here (Unestablished).
- **Concurrency under real load is untested.** No test or observation drives two
  live jobs, a load overlapping a generation, or a video run overlapping either.
- **Cross-platform behaviour** of the stop paths is untested; only Windows is
  claimed.
- The stage record's test counts (46 engine, 8 job) match the test functions
  present in the two files reviewed, and its other claims are consistent with the
  source — except the verdict sentence qualified under F2.
