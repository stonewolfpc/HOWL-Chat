# Audit: cancellation contract for synchronous image generation

Date: 2026-11-02
Mode: **read-only**. No product code, tests, config, artifacts or existing
records were modified. No builds, tests, generations or server requests were run.
Question: can an abandoned generation be safely followed by another request?

## Evidence classification key

Observed (runtime record) · Source-derived (current source) ·
Source/doc-derived (bundled backend material) · Inferred · Unestablished.

## Files and functions inspected

| File | Region |
|---|---|
| `internal/image/job/job.go` | `Start`, `run`, `transition`, `Cancel`, `finish` |
| `internal/image/engine/sdcpp/engine.go` | `GenerateTxt2Img` (~278), `CancelJob` (~353), `Shutdown`, `loadRun` (~124) |
| `internal/image/engine/sdcpp/client.go` | `txt2img` (411–560), `cancelJob` (581–598), `resetCancel` (604), `pollJob` (609–660), `img2img` (~860), `inpaint` (~1055), `progress` (1150+) |
| `internal/image/engine/sdcpp/video.go` | `SelectModel` (128), `generateVideo` (139+, `runMu.TryLock` at 154) |
| `internal/proc/job_windows.go`, `job_other.go` | process helper implementations per platform |
| `stable-diffusion.cpp.txt` | bundled backend material (license only) |
| Records | `issue-cancelled-generation-reuse.md`, `stage3-review.md`, `image-loading-stage.md` |

## 1. The synchronous request path

```
Manager.run (job.go)                       status generating
  -> engine.GenerateTxt2Img                engine.go:278
     -> client.txt2img                     client.go:411
        -> json.Marshal(request)           client.go:~536
        -> c.http.Post("/sdapi/v1/txt2img") client.go:545
        -> io.ReadAll(resp.Body); Close()  client.go:548-550
```

- **Context: none.** `http.Post` builds the request on `context.Background()`.
  There is no `NewRequestWithContext`, so nothing in the process can abort the
  in-flight request. (Source-derived)
- **Timeout: none.** The shared client is `&http.Client{Timeout: 0}` — "No
  timeout for long generations" (`client.go:~199`). (Source-derived)
- **Client-side cancellation point: none on this path.** `txt2img` never reads
  `c.cancelled`; it blocks in `Post` until the server answers. (Source-derived)
- **Connection closure:** impossible from our code — with no context, the only
  way the socket closes is process exit or server completion. (Source-derived)

## 2. What `Cancel` actually does

`Engine.CancelJob` (`engine.go:353`):

```go
if e.cancelVideo() { return nil }          // video only
if e.loadInProgress() { e.abortLoad(); return nil }   // LOAD only (Stage 3)
if e.httpClient == nil { return nil }
return e.httpClient.cancelJob("")          // generation
```

`client.cancelJob("")` (`client.go:581`):

```go
if jobID == "" { jobID = c.activeJobID }   // txt2img never sets activeJobID
c.cancelled = true
if jobID == "" { return nil }              // NO HTTP request is sent
```

`activeJobID` is assigned **only** by `img2img` (`client.go:888`) and `inpaint`
(`client.go:1081`), both of which use the async `/sdcpp/v1/img_gen` + `pollJob`
flow. The primary text-to-image flow uses `/sdapi/v1/txt2img` and therefore
always takes the empty-`jobID` branch.

**Conclusion (source-derived):** for synchronous txt2img — the main generation
path — a cancel is **local bookkeeping only**. No interrupt reaches sd-server,
no connection is closed, and the server-side generation is not stopped. Our
observed 230.82 s completion after a cancel is consistent with this, and is the
only runtime evidence that the server kept working.

Corollary: the UI's "Aborted" reflects the *job's* state, not the engine's work.

## 3. Backend interruption contract

**No contract is documented in this repository.** The only bundled backend file
is `stable-diffusion.cpp.txt`, which contains the MIT license text and nothing
else. There is no vendored sd.cpp source, no OpenAPI/spec, no `sd-server --help`
capture, and no README describing endpoints. (Source/doc-derived: absence)

Endpoints the client *already uses* (source-derived): `GET /sdapi/v1/samplers`,
`GET /sdapi/v1/progress`, `POST /sdapi/v1/txt2img`, `POST /sdcpp/v1/img_gen`,
`GET /sdcpp/v1/jobs/{id}`, `DELETE /sdcpp/v1/jobs/{id}` (reachable only with a
job id, i.e. img2img/inpaint).

**Unestablished** (cannot be settled without execution or upstream source):
- whether `DELETE /sdcpp/v1/jobs/{id}` actually interrupts a running job;
- whether a `/sdapi/v1/interrupt`-style endpoint exists in this build;
- whether sd-server cancels work when the HTTP client disconnects;
- whether sd-server serialises concurrent generations internally, queues them,
  or interleaves them;
- what it returns when two requests overlap.

Contradictory hint worth noting: `progress_stdout.go:3` states "sd-server does
not implement /sdapi/v1/progress", yet `client.progress()` calls it first. Our
assumptions about this API surface are evidently not reliable enough to design
interruption on without probing.

## 4. Lifecycle after a cancel — can two generations overlap?

- **Caller returns immediately:** yes for the RPC/UI (job becomes `cancelled`
  at once); **no** for the job goroutine, which stays blocked in `Post` until the
  server eventually answers. That goroutine leaks for the duration. (Source-derived)
- **Server keeps generating:** inferred from the mechanism above and consistent
  with the observed 230.82 s completion. Not directly proven.
- **A second request can reach the same process:** **yes.** `GenerateTxt2Img` and
  `txt2img` take no lock. `loadRun` covers loads only; `runMu` covers video only
  (`video.go:154`); `srvMu` guards stdout progress; `procMu` guards the child
  handle. Nothing serialises or drains generations. (Source-derived)
- **Ownership/drain mechanism:** none exists.
- **Can the next request receive the previous request's response?** No — each
  request reads its own `resp.Body`; there is no shared response buffer.
  (Source-derived)
- **Can it receive an empty result or error?** Yes, empirically: HTTP 500
  `generate_image returned no results` (Observed, `issue-cancelled-…md`).

**Overlap is therefore possible by construction.** Whether overlap is *what
produced* the 500 remains **Inferred**, not a finding. The alternative — sd.cpp
mishandling its own state after an abandoned request — is equally consistent
with the evidence gathered so far.

## 5. Secondary source-derived defect found during this audit

`resetCancel()` (`client.go:604`) has **no call sites**. `c.cancelled` is set
`true` by `cancelJob` and otherwise cleared only by `resetState()` (reached from
`Shutdown`). Consequences:
- async path (img2img/inpaint): after one cancel, `pollJob` (`client.go:614`)
  returns `job cancelled` at its first tick for **every later job** until the
  engine is shut down;
- `generateSequentialBatch` (`client.go:232`) also reads the stale flag.

Not observed live (the failures we saw were on the sync path), but it follows
directly from the source and is a second, independent cancellation defect.

## 6. Smallest defensible repair options

### Option A — cancellable request (responsiveness only)

`http.NewRequestWithContext` plus a per-generation `context.CancelFunc` the
engine can invoke; on cancel, close the body and return a cancellation error.
- Guarantees: the job goroutine stops waiting; the socket closes; no goroutine
  leak.
- Does **not** guarantee the server stops working. Disconnect-cancel behaviour is
  unestablished.
- Tests: fake server that holds a response; cancel; assert prompt return, assert
  the server observes its request context cancelled, assert no overlap is
  *implied* to be solved (it is not).

### Option B — backend interruption (only if established)

Send `DELETE /sdcpp/v1/jobs/{id}` or an interrupt endpoint, if and only if a
probe proves it works on this build.
- Blocked: evidence absent. Must not be assumed from another version or from
  A1111 compatibility folklore.

### Option C — generation slot / no-overlap floor (recommended floor)

One in-flight generation per engine process, enforced in `Engine` (not the job
manager): a second `Generate*` either waits for the slot or fails fast with a
clear message. A cancelled generation marks the slot "abandoned" and the engine
reports honestly that the work continues in the background.
- Guarantees: **two requests can never overlap**, so the empty-result collision
  class disappears regardless of what sd-server does internally.
- Limits: the next generation may be delayed by the abandoned one; needs an
  explicit policy (wait with a bound, or restart the child — see portability).
- Tests: fake server holding request 1; start request 2; assert it does not
  begin until request 1 completes (or is refused), and that cancel does not
  silently release the slot.

### Recommended boundary

**C as the correctness floor, A for responsiveness, B only after evidence.**
Plus the one-line class fix for §5 (call `resetCancel` at the start of every
generation, and set/clear `activeJobID` deliberately).

## 7. Acceptance-test plan (for the repair stage)

| # | Test | Kind | Asserts |
|---|---|---|---|
| 1 | Cancel a held request | unit (httptest) | client returns promptly, connection closed, no goroutine leak |
| 2 | Second request during the first | unit | no overlap: queued/refused, never concurrent |
| 3 | Cancelled flag does not leak | unit | next request proceeds (fixes §5) |
| 4 | Empty-result 500 is surfaced as a typed error | unit | job fails with the server message, not a generic one |
| 5 | Cancel then generate on the same engine | integration | artifact produced, no 500 (needs execution; **not** in this audit) |
| 6 | Backend interrupt probe | integration | DELETE/interrupt effects recorded, or recorded as absent |
| 7 | Disconnect behaviour | integration | what the server does when the client closes the socket |
| 8 | Slot release on engine death | unit/integration | crash does not wedge the slot forever |

## 8. Portability notes (public-product constraint)

- The cancellation design above is HTTP-level and platform-neutral.
- **Landmine:** `internal/proc/job_other.go` is a set of **no-ops** off Windows —
  `KillTree`, `Attach`, `CloseJob`, `HideWindow` do nothing. Any repair that
  falls back to "restart the child on cancel" via `proc.KillTree` would silently
  do nothing on macOS/Linux. A portable fallback must use `os.Process.Kill`
  (available everywhere) or add a real non-Windows `KillTree`.
- No hardware, VRAM, drive, model, or duration assumption appears in this design;
  none should be introduced by the repair.

## 9. Unresolved questions

1. Does `DELETE /sdcpp/v1/jobs/{id}` interrupt a *running* job on this build?
2. Does this sd-server implement `/sdapi/v1/progress` at all (see §3 conflict)?
3. Does the server cancel work on client disconnect?
4. Does it serialise concurrent generations, and what does it return to the
   loser? (Directly distinguishes the two hypotheses for the 500.)
5. Was the observed 500 caused by overlap or by server state after an abandoned
   request? **Not decided by this audit.**
6. Does the same defect exist for img2img/inpaint, which do have a job id and a
   DELETE path? Unverified.

## Status

Root cause of the 500: **not established.** Mechanism for "cancel does not stop
server work on the sync path": **established from source.** Overlap possible:
**established.** Overlap caused the 500: **inferred, unproven.**
