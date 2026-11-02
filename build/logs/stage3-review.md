# Stage 3 review (read-only)

Date: 2026-11-02
Reviewer: read-only source review against `build/logs/image-loading-stage.md`
Verdict: **freeze as evidence-qualified, with two recorded follow-ups**
Product code changed during this review: none.

## 1. Process ownership during cancellation — PASS, one scope flaw

One reaper owns `Wait()` (`startProcessReaper`), and every reader of the child
handle goes through `getServerProc`/`setServerProc` under `procMu`. Kill sites:
`abortLoad`, `LoadModel`'s pre-load stop, `LoadModel`'s error path, `Shutdown`.
None of them calls `Wait()`. That is consistent with the frozen stage's
single-owner contract.

**Flaw found — cancellation is engine-global, not job-scoped.**
`Manager.Cancel` calls `eng.CancelJob()` with no job identity. `CancelJob` then
asks the engine `loadInProgress()`, which is true for *whatever* load is running.
With two jobs in flight (job A loading, job B queued behind `loadRun`), cancelling
B sets the abort flag and kills A's child. A then fails with
`errLoadCancelled` and is recorded `failed`. Not reachable from the current UI
(`generationInFlight` serializes), but reachable from any other caller.
Severity: moderate. Not a blocker for freezing, is a blocker for calling
"concurrent start/cancel" solved.

## 2. Lock ordering around `loadRun` and `procMu` — PASS

Acquisition order observed in source:

```
loadRun (engine.go:124)  ->  loadMu (beginLoad/endLoad/loadAborted)
loadMu                   ->  procMu (abortLoad -> getServerProc)
procMu                   ->  (leaf: no other engine lock)
```

No path takes `procMu` and then `loadMu`/`loadRun`; no path takes `loadMu` and
then `loadRun`. `waitReadyOrExit` reads `processDone()` and `loadAborted()`
sequentially, not nested. `Shutdown` releases `procMu` before
`reapStoppedProcess`. Therefore no lock cycle is possible among these three,
and `-race` is clean on the package.

Residual risk: `loadRun` is held for the whole load (minutes by design). Any
future caller that takes `loadRun` and then waits on a job or on `abortFlow`
would deadlock. Worth a comment, not a change.

## 3. Can cancellation interrupt every loading path? — NO, two gaps

| Path | Interruptible | Evidence |
|---|---|---|
| Before the child starts (resolve/validate) | yes | `loadAborted()` check before `buildServerArgs` |
| During readiness | yes | abort checked at loop top, and on the exit + timeout branches |
| After readiness, before generation | n/a by design | `transition` refuses `loading -> generating` |
| **Video model selection** (`ensureEngineReady` -> `SelectModel`) | **no** | `SelectModel` never calls `LoadModel`, so `loading` is false and `CancelJob` does nothing |
| Engine lazy-init (`lazyInitEngine`) | no | short, but uncancellable |

The video gap matters because `ensureEngineReady` routes `Mode == "video"` to
`SelectModel` — a job can sit in `loading` with no way to stop it. Recorded as a
follow-up; video was outside this stage's stated scope.

## 4. Is a cancelled engine reusable? — PASS by source, PARTIAL live

After `errLoadCancelled`: `LoadModel`'s error path sets `serverProc = nil` and
`currentModel = nil`, `endLoad` clears `loading` (deferred), `loadRun` is
released (deferred), and `beginLoad` clears `aborting` at the next load. So the
next `EnsureReady` starts a fresh load rather than reusing half-built state.

Live evidence is indirect: the successful 94.55 s Z-Image generation ran after
several cancellations, but also after a full restart. Not a clean
"cancel, then generate without restart" observation.

## 5. Are terminal transitions consistently protected? — PASS for what is claimed

- `Cancel` sets `cancelled` under `m.mu`.
- `transition` is a compare-and-set, so `loading -> generating` loses to a cancel.
- `finish` treats `cancelled` and `completed` as final in both directions, and
  keeps an existing `failed` against a late success.
- A load that fails after the cancel does not resurrect the job.

Known holes, both pre-existing and documented rather than fixed:
`Cancel` on an already-`cancelled` job returns true and re-invokes the engine;
`List()` returns live pointers (`Status` was fixed, `List` was not).

## Verdict

The stage's claims hold against the source, with the qualification that
"concurrent start/cancel" is **partial**: the manager can protect its own state,
but engine-level cancellation has no notion of which job it belongs to. That is
a design limit of the existing single-engine model, not a defect introduced by
moving loading into the job.

**Freeze.** Follow-ups, in order: (a) the cancelled-generation reuse failure
(`issue-cancelled-generation-reuse.md`), (b) job-scoped cancellation, (c) video
loading cancellation.
