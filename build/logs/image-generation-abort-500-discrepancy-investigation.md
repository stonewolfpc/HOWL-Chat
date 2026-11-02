# Investigation: abort → process replacement → HTTP 500 → dead-looking UI

Date: 2026-11-02. Mode: **read-only**. No product file, test, configuration or frozen record was
modified; no build, no test run, no new generation. The live application was not touched during
this investigation.

Inputs read in full: `image-generation-settlement-guard-review.md`,
`image-generation-settlement-guard-live-acceptance.md`,
`image-ownership-provisional-freeze.md`,
`image-lifecycle-followup-r1-r2-r4-r5-r6-review.md`,
`image-generation-settlement-guard-stage.md`, and the sources listed in §2. Screenshots
`12-refused.png`, `13-third-complete.png`, `14-after-500.png` were re-inspected at full size.

Evidence tags: **[Observed]** captured in the live run · **[Source]** established by reading the
current implementation · **[Test]** established by an existing automated test · **[Inferred]**
supported but unproven · **[Unestablished]** evidence insufficient.

---

## 1. Executive verdict

**Nothing in the abort path kills or restarts `sd-server`. The PID change is explained by a lazy
relaunch, which only happens because the original child was already gone — and the app itself has
no path that ends that child on a generation cancel.** So the observed replacement is real, but its
*trigger* (why 20892 exited) is **unestablished**, and no evidence in the app records it.

The 500 is a genuine backend error that the app classifies correctly *in current source*
(`errEngineUncertain`, slot → `abandoned`). What the UI actually displayed, however, was the raw
transport text with **no uncertainty prefix** — which current source cannot produce. The running
binary (`build/bin/howl-chat.exe`, 18:59) predates the 19:08 settlement-guard edit and may also
predate the classification shape the current source shows, so **the live observation and the current
source disagree, and the disagreement cannot be resolved without instrumenting a run.** It is
recorded, not explained away.

The post-500 click is genuinely ambiguous in the product: after a 500 the engine is supposed to be
unavailable until restarted, but the UI has no surface for engine availability, and its only channel
—a single-line status field—**truncates with CSS ellipsis** (`.ig-status-text`: `white-space:
nowrap; overflow: hidden; text-overflow: ellipsis`), so the recovery instruction is exactly the part
that gets cut off. That is a real usability defect independent of the 500's cause.

Three separate things are being conflated in the live report and should not be: (a) the *restart* is
lazy and legitimate, (b) the *500* is a backend failure with unknown cause, (c) the *UI's silence* is
a presentation defect.

---

## 2. Source paths and relevant functions

| Layer | File | Function / line |
|---|---|---|
| UI | `frontend/dist/js/image-generator-core.js` | generate handler; abort branch ~894–905; `resetGenerateButton` 991; `setStatus` 1354 |
| UI style | `frontend/dist/css/ig-feedback.css` | `.ig-status-text` 40–52 (nowrap + ellipsis) |
| Binding | `app-image-generate.go` | `StartImageJob` ~30, `CancelImageJob` ~50, `GenerateImage` 73 |
| Binding | `app-image-model.go` | `ensureEngineReady` ~25, `EnsureReady` call 83, `SetActiveImageModel` ~110 |
| Job | `internal/image/job/job.go` | `Cancel` ~110, `run` ~160, poller liveness 201, `finish` ~250 |
| Engine | `internal/image/engine/sdcpp/engine.go` | `LoadModel` ~140, `IsRunning` 239, `EnsureReady` 259, `GenerateTxt2Img` 288, `CancelJob` 395, `Shutdown` 413 |
| Slot | `internal/image/engine/sdcpp/generation.go` | `beginGeneration` 132, `endGeneration` 174, `settleGeneration` 210, `cancelGeneration` 277, `resetGeneration` 290 |
| HTTP | `internal/image/engine/sdcpp/client.go` | `txt2img` ~419, request 588, 500 branch 593 |
| Supervision | `internal/image/engine/sdcpp/supervisor.go` | `abortLoad` 88, `startProcessReaper`, `waitReadyOrExit` |
| Shared | `internal/image/engine/sdcpp/video.go` | `stopServer` 116, `generateVideo` → `stopServer` 181 |
| Shared | `app-vram-arbitration.go`, `app-image.go` | `FreeVRAMFor` ~95, `UnloadImageModel` 176, `reinitEngine` 157 |

---

## 3. End-to-end abort trace and the next request

### What Abort cancels

1. UI abort branch (~894–905): `await CancelImageJob(activeJobId)`, clears the poll interval,
   rejects the pending promise, `resetGenerateButton()`, `setStatus('Aborted')`. **[Source]**
2. `CancelImageJob` rejects a blank id, then `Manager.Cancel(id)`; failure of `Cancel` (unknown or
   finished job) is swallowed by the UI's `catch { /* ignore */ }`. **[Source]**
3. `Manager.Cancel` sets `job.Status = statusCancelled` **before** touching the engine, then calls
   `eng.CancelJob(jobID)` outside the manager lock. **[Source]** Consequence: the job is terminal
   immediately, so whatever the engine later reports for that request is discarded by `finish()`
   (`if job.Status == statusCancelled || statusCompleted { return }`). **[Source]**
4. `CancelJob` routes by owner: video → load → generation. For a running generation it reaches
   `cancelGeneration(jobID)`, which calls the request's `genCancel()` and nothing else — no kill, no
   unload, no reset, no HTTP cancel endpoint (the client deliberately has none). **[Source]**

**So Abort closes the client request only.** It does not stop sd-server, does not unload the model,
does not release the slot, and does not claim the backend stopped. That matches the frozen contract
and the "Aborted" label in the UI. **[Source]**

### How the next Generate is admitted

- `StartImageJob` → `Manager.Start` → goroutine: `prepare` (`ensureEngineReady`) → `EnsureReady`
  **(Auto-load on generate)** → `beginGeneration`. **[Source]**
- `EnsureReady`: if `IsRunning()` is false, or the loaded model name differs, it calls `LoadModel`,
  which starts a new child. That single line is the **only** path that can start a process without
  the user pressing Load/Reload. **[Source]**
- `beginGeneration`: on a busy slot it self-heals **iff** `processConfirmedGone()` — i.e. the reaper
  saw the child exit — then refuses an active slot (`errEngineBusy`) or an abandoned slot
  (`errEngineUnavailable`). **[Source]** **[Test]** (`TestHTTPErrorIsNotReportedAsAClientCancellation`
  asserts the second request is refused after a 500.)

Therefore a *new* generation can only begin in one of two ways: the slot was idle, or the child had
demonstrably exited. There is no override, no "retry anyway" flag, and no automatic unload/reload/
engine-recreation path on the cancel or error route. **[Source]**

### What could unload/replace the engine (complete list)

- `LoadModel`'s `prev.Kill()` — only when a load is requested (model switch, or not running).
- `abortLoad` — only when the *cancelling job owns an in-progress load*, in which case it kills the
  loading child.
- `Shutdown` — explicit unload, app exit, `reinitEngine` (changing models dir / sd-server path).
- `FreeVRAMFor("llm")` → `UnloadImageModel` — when the user is making room for a chat model.
- `video.stopServer()` — before an sd-cli video render.

None of these ran during the test: no unload, reload, reinit, chat-model load, or video request
occurred between the abort and the third Generate. **[Observed]**

---

## 4. Process identity and ownership analysis

Facts:

- PID 20892 served run 1 and the start of run 2; RSS climbed ~5 GB → ~13 GB. **[Observed]**
- Run 2 was aborted ~3.5 s in; UI showed "Aborted". **[Observed]**
- The next Generate reached **"Step 1 of 8"** and a **new** `sd-server` PID 14320 existed with a
  start time of 19:19:14. **[Observed]**
- No Load/Reload click, no model change, no unload occurred. **[Observed]**

Analysis:

- A generation can only have started if the slot was free. Post-abort the slot is `abandoned`, and
  `beginGeneration` releases it only on `processConfirmedGone()`. **[Source]**
- A request can only reach "Step 1 of 8" if `EnsureReady` returned a live server; with PID 14320
  fresh, that server was started by `LoadModel` — which `EnsureReady` only calls when
  `IsRunning() == false`. **[Source]** `IsRunning` is a *process* question answered by the reaper's
  closed channel. **[Source]**
- **Therefore PID 20892 had already exited before the third Generate was submitted.** **[Inferred,
  high confidence]** That is the only combination consistent with both the self-healed slot and the
  new process.
- **Why it exited is unestablished.** No app path we can point at kills it on cancel; `abortLoad`
  requires an owned in-flight load (none: the model was already loaded and `EnsureReady` returned
  immediately), `LoadModel` was never invoked with a live `prev`, and no shutdown path ran. The
  remaining candidate is that **sd-server itself terminated** after its client disconnected
  mid-generation — a backend behaviour this project has never established, and exactly the kind of
  thing the "backend interruption is unestablished" clause forbids assuming. **[Unestablished]**
- No record of the exit survives: sd-server's stdout is pumped for progress only, `howl-errors-
  *.jsonl` has no entry for these events, and `procExit` is never persisted. **[Source]**
- It is **not** established that the app "replaced" the process, and it is equally not established
  that the old one "crashed". A changed PID alone proves only that the old process was gone.

Ownership note: the abandoned handle from run 2 is settled by the unwinding request when its HTTP
call returns; with the child already gone, `gone = true` in `endGeneration`, so the slot goes
**idle** rather than abandoned. **[Source]** That is guarantee 2 behaving as designed — but it means
the *next* request is admitted silently, which is exactly what the UI showed and why no refusal text
appeared. **[Inferred]** The refusal the contract describes requires a *live* engine with an
abandoned slot; that condition never existed here.

---

## 5. HTTP 500 classification and settlement state

### Did the request reach the backend?

Yes. **[Observed]** The UI advanced through sampling ("Step 1 of 8"), the message body
`{"error":"generate_image returned no results"}` is sd-server's own payload, and the status came from
`client.txt2img`'s non-200 branch (`sd-server txt2img error %d: %s`, client.go:593), which is only
reachable after `http.Do` returned a response. **[Source]** So this is **not** a not-dispatched
failure.

### How current source classifies it

`classifyGeneration(cancelled=false, err)` → `genOutcomeAmbiguous` → `endGeneration` retains the
slot → `genOutcomeAbandoned`, owner kept → `describeGenError` prepends
`errEngineUncertain` ("image generation failed and the engine may still be working on it, so it
stays unavailable until it is restarted"). **[Source]** **[Test]**
(`TestHTTPErrorIsNotReportedAsAClientCancellation` asserts `errEngineUncertain`, state `abandoned`,
and that a second generation is refused with `errEngineUnavailable`.)

### Contract match — policy: yes; observed text: no

- **Policy** matches the freeze: an unconfirmed backend failure must not hand the slot to the next
  request, and must not be called a cancellation. **[Source]**
- **The displayed text does not match.** The status line read
  `Error: sd-server txt2img error 500: {"error":"generate_image ret…` — the prefix is at the
  *start* of the string, and the clipping is at the *end*, so the missing `errEngineUncertain`
  prefix is a real absence, not a truncation artifact. **[Observed]** (screenshots 13 and 14, both
  re-inspected.)
- Current source cannot produce that string from the settle path. **[Source]** Two candidate
  explanations remain, and neither is proven:
  1. **Version skew** — the running exe is `build/bin/howl-chat.exe` built **18:59**, i.e. *before*
     `generation.go`'s 19:08 edit. It contains "the engine may still be working on it" and "client
     request cancelled" but **not** "generation was already settled", confirming it predates the
     settlement guard. It may also predate the exact describe/classification shape the current
     source shows. **[Observed]** (binary string check; a read-only `grep -a` over the running exe.)
  2. **A path that bypasses `settleGeneration`** — no such path was found in current source: every
     `Generate*` settles once, exactly, on its own stack. **[Source]**
- The **root cause of the 500 itself remains unestablished**, and this investigation does not
  propose one. Two observations are worth carrying: it occurred on a *freshly restarted* process,
  which had answered readiness but had not yet produced a single image; and a 500 mid-generation on
  this backend accompanies "no results", i.e. the server ran and produced nothing.

### Slot after the 500 (source expectation)

`abandoned`, owner retained, and a subsequent request refused with `errEngineUnavailable` until an
explicit unload/reload (or a confirmed child exit). **[Source]** What the UI then showed is analysed
in §6 — it does **not** match that expectation, and the evidence is insufficient to say whether the
slot, the request, or the UI was the deviating layer.

---

## 6. Frontend status and recovery path

- **Single channel.** `setStatus` writes `textContent` of `#ig-status`; the failure branch renders
  `'Error: ' + (err.message || err)`. **[Source]** There is no separate "engine available /
  unavailable" indicator, no job-state field, and no recovery action attached to the message.
- **The channel truncates.** `.ig-status-text` is `white-space: nowrap; overflow: hidden;
  text-overflow: ellipsis`. **[Source]** The contract's recovery sentence ("…stays unavailable until
  it is restarted") is therefore rendered as `…generate_image ret…`. **[Observed]** seen in
  screenshots 13 and 14. This is a **presentation defect**: even a correct refusal from the engine
  would be cut off before the instruction.
- **Button state after failure.** The `finally` block always calls `resetGenerateButton()` (clears
  the abort `onclick`, re-enables, restores "Generate"), so the button is consistent with "no job
  running". **[Source]**
- **Why the post-500 click was indistinguishable.** With the slot expected `abandoned` and the
  process alive, a click should have produced a *fast* `errEngineUnavailable` — a status text
  beginning "Error: image engine is unavailable until it is restarted". That text was never seen;
  what was seen was the identical 500 text, an idle GENERATE button, and no new process.
  **[Observed]** Possible readings, none proven:
  1. the new job never started (click not delivered / handler not bound);
  2. the new job failed instantly with a message so similar it looked unchanged — implausible, since
     the unavailable/wrapped texts differ at the start;
  3. the engine's failure was swallowed and the job was reported with a stale error.
  Since `Manager.Cancel`/`finish` can discard a late engine failure for a terminal job, and the UI
  only ever sees one string, **the user cannot tell which of these happened** — that is the finding.
  **[Unestablished]**
- **Consequence for the user:** after a 500 the app offers no statement of engine availability and
  no instruction, so "click again" is the only visible option; the only reliable recovery the test
  found was explicit **Unload → Reload**. **[Observed]** (that recovery did work in the same run.)

---

## 7. Shared image/video infrastructure

These paths can end the image server while a generation is live; none was exercised, and none was
changed. **[Source]**

1. **Video.** `generateVideo` calls `stopServer()`, which kills `serverProc` unconditionally
   (`proc.KillTree` + `Kill` + `os.Process.Wait()`), reads `e.serverProc` without `procMu`, and calls
   `Wait()` beside the reaper's own `cmd.Wait()`. It also does **not** call `resetGeneration`, so a
   live image generation would die with a transport error and settle as `abandoned` with no close
   reason. This is the known F4/R3 area and remains deferred.
2. **VRAM arbitration.** `FreeVRAMFor("llm")` unloads the image model (`UnloadImageModel` →
   `Shutdown`), which closes a live generation as `errEngineStopped` ("engine unloaded"). That is
   correct reporting, but it is a second, user-triggered route that ends an image generation
   mid-flight and is not surfaced as such beforehand.
3. **Engine re-init.** `SetImageModelsDir` / `SetSDServerPath` → `reinitEngine` shuts the engine
   down and swaps engines while a job may be running; the running job's engine reference is the old
   one, and `jobManager.SetEngine` only affects later jobs.
4. **Slot is not process-wide.** Ownership is per `Engine` struct, so two engines (before/after a
   reinit) each have their own slot. No evidence this occurred here; noted because it makes "the
   engine" ambiguous in any future multi-engine diagnosis.

Nothing in the abort path touches video state, and nothing above was triggered during the live test.
**[Observed]**

---

## 8. Evidence table

| # | Claim | Class |
|---|---|---|
| 1 | Abort calls `CancelImageJob` → `Manager.Cancel` → `CancelJob` → `cancelGeneration`, which only cancels the request context | **Source** |
| 2 | Manager marks the job `cancelled` before the engine is called, so later engine results are discarded | **Source** |
| 3 | No cancel/unload/reset/kill path runs on generation cancel | **Source** |
| 4 | A new generation starts only if the slot is idle or the child is confirmed gone; otherwise refused | **Source**, **Test** |
| 5 | The UI showed "Aborted" and no claim the backend stopped | **Observed** (shot 11) |
| 6 | The next Generate started at "Step 1 of 8" with PID 14320 replacing 20892 | **Observed** (shot 12) |
| 7 | No user action between abort and that request could have unloaded/restarted the engine | **Observed** |
| 8 | The old child had therefore exited before the third request | **Inferred** (high confidence) |
| 9 | Why it exited (crash vs. self-termination on client disconnect) | **Unestablished** — no log survives |
| 10 | The 500 request did reach sd-server (steps advanced; sd-server's own error body) | **Observed** + **Source** |
| 11 | Current source classifies a 500 as `ambiguous` → `abandoned` → `errEngineUncertain` | **Source**, **Test** |
| 12 | The UI displayed the raw transport text with **no** `errEngineUncertain` prefix | **Observed** (shots 13, 14) |
| 13 | Current source cannot produce that raw text on the settle path | **Source** |
| 14 | The running binary is 18:59, predates the 19:08 guard edit, and lacks the guard's sentinel string | **Observed** (binary string check) |
| 15 | Whether version skew explains #12 | **Unestablished** (no reproduction permitted) |
| 16 | Post-500 click produced no new process, no ABORT state, and unchanged text | **Observed** (shot 14) |
| 17 | The status field truncates with ellipsis, cutting recovery instructions | **Source** + **Observed** |
| 18 | The UI exposes no engine-availability state or recovery action | **Source** |
| 19 | Unload → Reload restores a working engine that generates again | **Observed** (shots 15–21) |
| 20 | Video `stopServer` kills the image server unconditionally without resetting the slot | **Source** (F4, deferred) |
| 21 | `FreeVRAMFor("llm")` unloads the image model, ending a live generation as "engine unloaded" | **Source** |
| 22 | Root cause of the historical 500 | **Unestablished** — unchanged by this investigation |

---

## 9. Exact unresolved questions

1. **Why did PID 20892 exit?** Did sd-server terminate on client disconnect, abort internally, or
   crash? No exit status, stderr tail, or event was persisted, and `procExit` is discarded.
2. **Why did the UI show an unwrapped 500?** Which build produced it, and does that build's
   `settleGeneration`/`describeGenError` differ from the current source? The 18:59 binary predates
   the guard; whether it also predates the wrapping is not established.
3. **What did the post-500 click actually do?** No job state was observable, and the only channel is
   a truncated single line.
4. **Was the slot `idle` or `abandoned` after the 500?** Source says `abandoned` while the process
   lives; the UI cannot show it, and no refusal text appeared.
5. **Does sd-server answer readiness before it can serve a txt2img?** A fresh process answering
   `probe()` and then failing with "no results" would make readiness a false promise. Untested.
6. **Would the same abort→generate sequence reproduce?** Not tested: reproducing requires new
   generations, which this investigation is not authorized to run.

---

## 10. Recommended next action (smallest justified scope)

**Instrument, do not repair.** One bounded, read-only-to-behaviour change is justified next, because
every remaining question is an observability gap rather than a policy disagreement:

1. **Persist child exits.** When the reaper completes, log one line with PID, exit code, elapsed
   time, and the last ~40 lines of sd-server stderr. This alone answers Q1 and Q2's "did the backend
   die" half, and costs no behaviour change.
2. **Log the settled outcome at the job boundary.** One line per job: outcome class, close reason,
   final error string, slot state. That answers Q2 and Q4 without guessing.
3. **After that, one controlled reproduction** of abort → generate on a build that includes the
   logging, with a single generation at a time and no parallel work; compare the recorded exit with
   the observed 500.

Only then is a fix justified, and only for whatever the logs establish. Separately — and independent
of the 500 — the **status line truncation** (§6) is a confirmed defect worth its own small change:
the one message that explains an unavailable engine is the one the UI cuts off.

Not recommended now: touching the settle/ownership policy, adding retries, adding an auto-restart
after failures, or changing video/server sharing. Each would add a layer on top of an unmeasured
cause, which the project's standing constraints forbid.

**Stop.** No repair was made.

---

# Addendum (2026-11-02, later): live acceptance on the correctly built binary

This addendum records the build-identity correction and the bounded live checks that followed it.
It supersedes the live-evidence rows and the §5/§6 conclusions of the body; the source analysis in
§2–§4 stands unchanged. No frozen record was edited.

## A1. Build identity (established)

| Item | Value |
|---|---|
| Revision | `b2fa517` (master, 2026-11-01 09:40), 113 modified files, worktree diff hash `51c19042e76591a2afcbc666f8120fb9f63343b7` |
| Sources newer than the old binary | exactly two: `sdcpp/generation.go`, `sdcpp/settlement_guard_test.go` (both 19:08) |
| Old binary (what the first acceptance tested) | `build/bin/howl-chat.exe`, 18:59, sha256 `b7781762…90520`, **no** "generation was already settled" string |
| New binary (built this pass) | recorded workflow `CGO_ENABLED=0 GOFLAGS=-buildvcs=false go build -tags "desktop,production" -ldflags "-w -s -H windowsgui"` → exit 0, sha256 `d7e8aacd…91866`, guard sentinel present |
| Tree hash before/after build | unchanged `51c19042…` — the build did not alter sources |
| Process actually running the new build | **PID 9252**, path `build/bin/howl-chat.exe`, start 20:20:42, verified by path + start time + file hash |
| Old binary retained | `build/bin/howl-chat-1859.exe.bak` |

The discrepancy was a **version gap, not a policy failure**: the first acceptance exercised a build
that predates the settlement guard by nine minutes and lacks its sentinel string.

## A2. Bounded live checks on PID 9252 (observed; screenshots 22–39)

| # | Action | Result |
|---|---|---|
| 1 | Generate from idle, no model loaded | `ABORT` state, "Loading model...", child **PID 27340** (20:22:09) → "Step 1 of 8" — auto-load on generate works |
| 2 | Abort mid-sampling | status "Aborted"; **PID 27340 still alive** — abort does not kill the backend |
| 3 | Immediately submit again | **Refused**: `Error: image engine is unavailable until it is restarted: the previous generation was abandoned and its completion could not be confirmed (owner "job_1793586129940319900")`; no new child, no steps |
| 4 | Unload | child gone |
| 5 | Reload | new child **PID 24780** (20:23:37), "Model ready: z-image-turbo-Q4_K_M" |
| 6 | Generate on restarted engine | admitted on the same child ("Generating...", no re-load) — reset released the slot |
| 7 | Abort | "Aborted", app left idle, no orphan processes |

Three generation attempts total, two aborted, no images produced, gallery unchanged, no stress run.

## A3. Corrections this forces

1. **The abandoned-slot contract is live on the current build.** A client-side abort leaves the engine
   unavailable until restart — `errEngineUnavailable`, naming the abandoned owner. On the 18:59 binary
   the same sequence **admitted the next request** (it reached "Step 1 of 8" on a new PID), which is
   how the original run 3 reached the backend and returned the 500 at all. The 500 sequence as first
   observed therefore **cannot occur on the current source** in that form: the next request is refused
   before dispatch. This substantially strengthens the version-gap explanation for the unwrapped 500
   text, though it still does not identify why sd-server exited.
2. **Status truncation is narrower than first written.** The refusal *does* render the actionable
   clause ("…unavailable until it is restarted") and clips only the trailing owner detail. The earlier
   claim that the recovery instruction is cut off is **withdrawn**; the remaining defect is that the
   diagnostic tail is elided and that one line is the only channel. Still a UI-pass item, smaller.
3. **Two failure modes remain unexercised**: a real backend 500 on this build (no fault was induced or
   authorized), and the `errGenerationSettled` sentinel, which is unit-only and unreachable from the UI.

## A4. Boundary of this pass

The settlement guard **passes source review and now passes the bounded live checks that do not require
a backend fault**. Its fault-path acceptance remains open, deliberately: inducing a 500 would mean
fabricating a failure, which the project's evidence rules forbid.

Unchanged and still open: **why PID 20892 exited** (deferred to the logging/diagnostics audit; no
instrumentation was added, by decision), **video `stopServer` process ownership** (separate defect, not
exercised, not declared safe), and **the status-line elision** (deferred to the UI repair pass).

No product file, test, configuration or frozen record was modified in this addendum; the only writes
were the new binary, its `.bak`, the screenshots, the probe dumps under `build/tmp`, and this record.
