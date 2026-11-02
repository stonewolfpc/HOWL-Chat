# Review: generation settlement guard (Review-1) — read-only source inspection

Date: 2026-11-02. Mode: read-only for source. No production file, test, config or earlier
record was modified. Evidence categories are separated: claims marked **source** come from
reading the current code; claims marked **reported** come from the stage record and were
not executed here.

Target: `build/logs/image-generation-settlement-guard-stage.md` (19:09)
Contract: `build/logs/image-ownership-provisional-freeze.md` (18:55)
Predecessor review: `build/logs/image-lifecycle-followup-r1-r2-r4-r5-r6-review.md` (19:02)

**Verdict: PASS.** Review-1 is closed correctly. The handle is invalidated inside the same
critical section that decides the transition, a second settlement can neither move the slot
nor be reported as a cancellation, and no frozen path, lock order or branch changed.
Four informational notes follow; none requires a change.

Boundary evidence: `generation.go` and `settlement_guard_test.go` are both 19:08, i.e.
before the review began; `client.go` 18:57, `engine.go` 18:59, `video.go` 18:47, the stage
record 19:09, the freeze 18:55. No source file carries a later timestamp.

## What was verified

### 1. Invalidation happens inside the deciding boundary — PASS (source)

`endGeneration` (`generation.go:174-207`) performs, under one `e.genMu.Lock()`:
the transition (`genAbandoned`/`genIdle` + owner), clearing `genCancel`/`genHandle`,
`h.settled = true`, and `h.token = 0`. There is no window in which the slot has been
released but the handle can still match it, and no second lock is introduced.

### 2. A second settlement cannot change the slot — PASS (source)

`settleGeneration` (`:210-228`) returns before classifying when `handleSettled` (`:230-234`)
is true. Even if it did not, `endGeneration`'s guard is now `e.genSeq == h.token`, and
`genSeq` only advances in `beginGeneration`/`resetGeneration` while `h.token` is zero.

Three independent layers therefore protect the slot: the settled check, the zeroed token,
and the sequence comparison. The test asserts the outcome (`GenerationState() == "idle"`)
and that the slot is still claimable afterwards.

### 3. Not misreported as cancellation or uncertainty — PASS (source)

The settled branch returns `errGenerationSettled` (`:89-91`), immediately, before
`classifyGeneration` reads `h.ctx.Err()`. This is what closes the second half of the
defect: the stage record's failing assertion
(`want errEngineUncertain … got client request cancelled`) shows the first
implementation left `h.cancel()`'s effect visible to classification. Now the cancelled
context is never consulted for a settled handle. `errGenerationSettled` matches neither
`errClientCancelled` nor `errEngineUncertain`, which the new test checks with `errors.Is`.

### 4. A settled handle cannot affect a later generation — PASS (source)

`beginGeneration` allocates a fresh handle with `token: e.genSeq` after incrementing; a
stale handle's `token` is 0 or behind, and its `settled` flag belongs to that handle only.
The test claims `job2` and settles it, asserting `idle` at the end. Live behaviour of slot
reuse is covered separately in the acceptance record.

### 5. First settlements retain their behaviour — PASS (source + reported tests)

`classifyGeneration`, `describeGenError` and the `retain`/`gone` logic are unchanged, so
completed → `idle`, cancelled → `abandoned` + `errClientCancelled`, ambiguous →
`abandoned` + `errEngineUncertain`. The new table test pins all three; the 62 pre-existing
engine tests (including the shutdown/`errEngineStopped` and batch tests) were reported
passing by the stage and are not re-run here.

### 6. Repeated settlement, shutdown and unwinding — PASS (source), with note N2

- Shutdown/unload/reload call `resetGeneration`, which bumps `genSeq` and records the close
  reason on the handle. The unwinding request's `settleGeneration` then finds no token
  match: no transition, but `describeGenError` still reports `errEngineStopped` with the
  generation's own reason. That is the R2/F7 behaviour, unchanged.
- A second settlement of such a reset-stale handle is also inert (`token` mismatch), and it
  describes again rather than returning the sentinel — see N2.

### 7. No lock-order or synchronisation regression — PASS (source)

`handleSettled` takes `genMu` and releases it before `endGeneration` takes it again: no
nested acquisition, no path holds `genMu` while calling out. `processConfirmedGone()` is
still called before locking, and `h.cancel()` after unlocking. `h.settled` and `h.token`
are only read/written under `genMu`. The race detector result for the stage is **reported**
(`ok … 3.262s`), not re-run here.

### 8. The seven provisional guarantees — PASS (source)

1 admission/refusal — `beginGeneration` untouched; self-heal via `processConfirmedGone`
unchanged. 2 release policy — same `retain` set, same confirmed-exit release. 3 close
reported as close — `describeGenError` unchanged. 4 no failure reported as cancellation —
strengthened: the settled branch cannot reach the cancellation text. 5 routed by owner,
empty id inert — `cancelGeneration`/`CancelJob` untouched. 6 progress tokens — untouched.
7 unload/reload/exit report their own reason — `resetGeneration` and the reason field
untouched; the new flag is separate from `close`.

## Informational notes (no change requested)

- **N1 — TOCTOU between the settled check and the transition (concurrent only).** Two
  goroutines settling the same handle at once could both pass `handleSettled` before either
  sets the flag; the loser then reaches `endGeneration`, finds no token match, and — since
  it classified from the now-cancelled context — could return `errClientCancelled`. The slot
  stays correct (token guard). No production path settles one handle from two goroutines;
  every `Generate*` calls `settleGeneration` exactly once on its own stack. Recorded, not
  fixed: making the check and the transition one critical section would be a larger change
  than Review-1 authorises.
- **N2 — reset-stale handles are not marked settled.** `h.settled = true` is set only when
  the token matched. After a reset the request still needs to describe its own outcome, so
  this is deliberate; the consequence is that such a handle can describe twice. No state
  transition is possible either time.
- **N3 — nil handle semantics changed slightly.** `settleGeneration(nil, err)` now returns
  `err` unchanged instead of describing it as completed. No caller passes nil (all four
  call sites use a handle from `beginGeneration`), so this is unobservable.
- **N4 — the late error is embedded, not wrapped.** `errGenerationSettled` includes the
  later error with `%v`. Intentional: this is a programming-error signal, and exposing a
  second `%w` would invite callers to treat it as a real result.

One earlier-stage observation still stands and was not touched: `errGenerationSettled` is
unreachable from current production code, by design.

## Not established here

- No test, build, vet or race run was performed in this review; those results are **reported
  evidence** from the stage record. Test counts are independently checkable and consistent:
  64 engine test functions (62 + 2) and 8 job tests, matching the record.
- Whether the double-settlement condition is reachable through the UI: it is not. It is an
  internal invariant, and the live record treats it that way.
- Backend interruption and the historical 500 remain unestablished by this review; the live
  acceptance record separately documents a 500 observed during UI testing.
