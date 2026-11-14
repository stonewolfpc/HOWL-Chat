# Adversarial Review — Stage 1 TTS Helper Supervisor and Protocol Client

**Date:** 2026-11-06
**Reviewer:** Muse (independent adversarial pass)
**Subject:** `docs/review/IMPLEMENTATION-howl-tts-helper-supervisor-stage1-2026-11-06.md`
**Subject SHA-256 (verified):** `1e8e6553a6e5c2f71f391940fe6ef670fbcf73d4c9381b23b066eeb639705d0a`
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` §17 as clarified by A1 (§17.18)

**Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. Not an authorization. Stage 1 remains **HOLD**;
Stage 2 is **not authorized**. Recommended disposition is in §8.

**Method:** inspected every implementation file and test (not the report's summary), re-ran all reported
verification commands, verified all 12 artifact hashes, confirmed the worktree claims, and executed one
temporary adversarial probe (since removed) to convert the principal finding from reading into
demonstrated fact. No production file was modified; the probe file was deleted after its run.

---

## 1. Evidence verification — all reported claims reproduce

| Reported claim | Independent result |
|---|---|
| Report SHA-256 | **Matches exactly** (`sha256sum` of the reviewed file) |
| All 12 artifact hashes (§12) | **All 12 match** the files on disk |
| `go test ./internal/audio/tts/helper -count=1` → ok | **Reproduces** (`ok … 1.339s`) |
| `go test -race` helper + types → ok | **Reproduces** (3.118s / 1.460s), no data races detected |
| `go vet` → no diagnostics | **Reproduces** (exit 0) |
| `CGO_ENABLED=0 go build` → exit 0 | **Reproduces**; temp artifact removed afterward |
| `go test .` → ok | **Reproduces** (`ok … 10.044s`) |
| `go test ./internal/audio/...` → all pass | **Reproduces** (all 16 packages ok) |
| 30 tests listed in §7 | **All 30 exist** in the test files; none omitted, none extra |
| "All implementation files are new" (§2, §13) | **Confirmed**: `git status` shows `internal/audio/tts/helper/` and `provider.go` as untracked (`??`) only |
| Tracked diff is exactly `interfaces.go` + `utils.go`, 75+/3− (§13) | **Confirmed byte-for-byte** via `git diff --stat`; the diff contains **no helper/supervisor/inspector/provider references** — it is unrelated character/pipeline work, consistent with "not edited during Stage 1" |
| `TTSSynthesizer` unchanged | **Confirmed** (`interfaces.go:126` interface intact) |

No evidence discrepancies found. The report describes the code that is on disk.

---

## 2. What the implementation gets right (load-bearing strengths)

These are not courtesies; they are the reason the findings below are bounded rather than structural.

- **The acceptance barrier is correctly built.** The final check in `Synthesize` (`client.go`, post-validation
  block) re-reads cancellation state under the same `state.mu` that `watchCancellation` and the `ctx.Done`
  branch write under. The two orderings the contract cares about — cancel-then-validate and
  validate-then-cancel — are both tested (`TestCancellationAfterTerminalBeforeAcceptance`,
  `TestAcceptedResultIsNotRetractedByLaterCancellation`), including via a real validation-window hook.
  The tricky third case (Go `select` randomly picking an already-ready success over an already-fired
  `ctx.Done`) is still safe, because the post-validation re-check observes `cancel_requested` regardless
  of which branch won the select.
- **No autonomous restart exists.** `observeExit` only records; nothing re-spawns. Restart happens solely
  through an explicit later client call, matching §17.11.5's "restart only on the next explicit TTS request".
- **Single-`Wait()` ownership holds.** Exactly one observer goroutine per successful start; `osProcess.Wait`
  is additionally `waitOnce`-guarded; the failed-start path creates no waiter (tested).
- **Cleanup paths are caller-parameter-derived, never wire-derived.** `cleanupExact` and the validation
  path build the final path from the caller's `OutputDirectory`/`Filename`, so a malicious helper cannot
  redirect deletion by lying in `output_path`. `validateSynthesisResult` independently pins the wire path
  to the requested directory and basename.
- **The filename allowlist excludes every glob metacharacter** (`^[A-Za-z0-9][A-Za-z0-9._-]{0,123}\.wav$`,
  `protocol.go`), so the `<final>.part-*` glob in cleanup cannot be widened through the filename.
- **Strict decoding is thorough in both directions**: duplicate keys, unknown fields, NUL, BOM, non-UTF-8,
  and the 4 MiB cap on send and receive.
- **Retired-ID handling matches A1-D5**, including session rotation on process replacement (`client.go`
  `Start`), so a late frame from a dead session cannot complete a new request.

---

## 3. Question 1 — cancellation races: SOUND, with one classification caveat

**Verdict: cancellation cannot lose to an arrived-but-unaccepted success.** The barrier resolves the race
atomically: exactly one of `accepted` / `cancelled` commits under `state.mu`, and the tests prove both
orders. A1-D1/D2 are correctly implemented.

Caveat (minor, diagnostic only): when the request context fires while the outcome is a *failure* (e.g. a
crash raced with cancellation), the code returns `cancelled` and discards the crash diagnostic
(`helper_crashed` is retryable; `cancelled` is not). Cancellation winning is defensible; losing the
crash signal is not ideal for restart-budget accounting later. Record, do not block on.

---

## 4. Question 2 — process lifecycle: one REAL finding, two minor leaks, two observations

### F-1 — `failAll` can block forever on a full outcome channel, skipping `Kill` (MEDIUM — proven by probe)

`Client.failAll` (`client.go`) delivers to each pending channel with an unconditional blocking send:

```go
for _, req := range pending {
    req.ch <- pendingOutcome{failure: f}
}
```

Each channel has capacity 1. If it already holds an unconsumed success terminal — e.g. `dispatch` delivered
it, but `Synthesize`'s `select` chose the timeout (or cancel-grace-expiry) branch instead — the send blocks
**forever**, because the consumer has retired the request and will never drain the channel. Consequences,
in order: the readLoop goroutine is stuck; `protocolViolation`'s `_ = c.sup.Kill(...)` (which runs *after*
`failAll`) never executes; a live-but-violating helper is never terminated; the reader goroutine leaks.

**Proven, not inferred.** A temporary in-package probe (full channel → `failAll` in goroutine → 2 s
watchdog) confirmed the permanent block; the empty-channel control returned instantly. The probe file was
deleted afterward; the package directory is back to its 12 files.

Reachability: benign helpers need a near-simultaneous timeout plus a second fatal trigger (duplicate
terminal, malformed frame, unknown ID, EOF) — low likelihood. A *malicious* helper (success frame, then
garbage) reaches it deterministically — and §17.10.2(4) explicitly places a compromised helper in scope.
Impact is bounded: `readFailure` is assigned *before* the send loop, so later operations fail fast rather
than hanging, and the job object still reaps the process at app exit. But the contract's kill/wait-before-
cleanup ordering (§17.11.6 rule 4, A1-D6) is violated on exactly the path built for adversarial input.

Fix direction (for the implementing stage, not this review): non-blocking delivery with a drain-on-retire
invariant, or close-based outcome signalling. **This should be fixed before Stage 2**, because every later
stage inherits this client.

### F-2 — non-incompatible handshake failure leaks a reader and reuses an unvalidated process (LOW)

If `helloHandshake` fails with anything other than `protocol_incompatible`/`protocol_minor_unsupported`,
`Start` resets `started`/`transport` but does **not** terminate the process and does **not** stop the
already-spawned `readLoop`. The next `Start` reuses the same process (supervisor still holds it) and spawns
a **second** `readLoop` on the same `bufio.Reader`. The two readers split frames between them; dispatch is
ID-based on shared maps so correlation survives, but one goroutine is leaked per failed handshake and the
supervisor's state sits in `handshaking` indefinitely. Fail-closed (no wrong result), but a real leak on a
path a version-skewed helper would hit repeatedly.

### F-3 — the `outcome.failure` path cleans up without kill-observed-exit (LOW)

On a stream failure or helper-reported error, `Synthesize` calls `cleanupExact` immediately with no
`Terminate` — while, on the protocol-violation sub-path, the kill is only fired asynchronously inside the
readLoop. A1-D6 requires kill-and-observed-exit before deleting files the helper could still be writing.
On Windows an in-use file mostly fails `Remove` (fail-safe direction → stale `.part` rather than data
loss), so this is a cleanup-completeness gap rather than a safety hole, but it is a contract-ordering gap
on the adversarial path. F-1's fix and this path should be addressed together.

### F-4 — `Shutdown` bypasses `workMu`; concurrent Shutdown + operation is untested (INFO)

`Shutdown` sends its frame and terminates without holding the operation serialization lock, so an
in-flight `Synthesize` fails with `helper_crashed` while `Shutdown` returns nil. No corruption (single
terminal per request is preserved), and single-caller use is safe. The future adapter — the only
concurrent caller — must own this serialization; record it as an adapter-stage requirement.

### F-5 — transient liveness misreporting between death and observation (INFO)

`SetState` flips `processState` to alive whenever `proc != nil`, so between process death and the observer
running, `ProcessAlive()`/`Status()` can report a live process. Self-healing once `observeExit` runs;
do not build restart-budget counting on `Status()` snapshots.

**On concurrent control requests and goroutine leaks generally:** `begin`/`retire`/`dispatch` are
fully mutex-guarded; channel sends in `dispatch` occur after unlock into capacity-1 buffers with
single-terminal semantics; `watchCancellation` always exits via `state.done`; the only unbounded block
found is F-1. `TestNoGoroutineLeakAfterShutdown` covers the clean path; F-1/F-2 are the uncovered ones.

---

## 5. Question 3 — WAV cleanup security: leaf checks STRONG, one contract rule NOT met

**What holds:** absolute-path + `filepath.Rel` containment (no string-prefix check), caller-derived exact
paths, leaf `Lstat` symlink rejection in both discovery and output validation, regular-file enforcement,
byte-identity checks (`bytes` vs actual length, canonical 44-byte header, PCM16/mono, rate range,
duration ±1 ms, SHA-256 recompute, identity echo, `.part` absence) — all performed on the bytes actually
read. `validateSynthesisResult` pins the wire `output_path` to the requested directory and basename, so
cleanup can never be steered by helper output. The escaping-output test proves the pin.

**What does not hold — G-1 (contract non-conformance, must stay open):** §17.10.2 rule 5 and §17.11.1
require **reparse-point/symlink rejection in every path component** on Windows. The implementation checks
only the leaf (`ModeSymlink` + `IsRegular`). The report admits this (§11.4: "partial"), which is honest,
but it means the output contract is *not* fully implemented: a reparse point in an intermediate component
of the output directory or helper path passes validation. Stage 1 cannot close this without the
`FILE_ATTRIBUTE_REPARSE_POINT` per-component walk; the helper side cannot repeat any check because no
helper exists. End-to-end containment is therefore **not established** — correctly scoped as future work,
but it must be tracked as a failed rule, not a limitation note.

Minor robustness notes: `filepath.Glob` is fed a directory that may contain `[` (only the filename is
pattern-constrained) — fail-closed in validation, silent part-cleanup skip in `cleanupExact`; and the
`Lstat`→`ReadFile` sequence is a standard TOCTOU window, low severity here since all acceptance checks
run on the in-memory bytes.

---

## 6. Question 4 — restart/backoff/idle-close against the frozen contract

The contract distinguishes mechanism (frozen, §17.11.5) from values (provisional, §17.8). The Stage 1
boundary (§17.14) authorizes "the main-side fake transport, supervisor, adapter tests, and golden
fixtures" plus proof of "absent/crash/hang behavior with the main build still `CGO_ENABLED=0`".

| Contract item | Stage 1 state | Assessment |
|---|---|---|
| Restart only on explicit request | **Implemented** (no auto-restart anywhere) | Required behavior present |
| Crash observed, active work failed, Chat unaffected | **Implemented + tested** | Required behavior present |
| Backoff 1/2/5 s schedule | Not implemented | **Acceptable omission**: no autonomous restart loop exists, so the schedule has no attachment site. Deferred *mechanism*, not deferred tuning — must arrive with the restart owner |
| Restart budget 3/10 min + no-restart list | Not implemented | Same as above; note the no-restart list (`protocol_incompatible`, `helper_missing`, …) is partially honored by structure (those paths `Terminate`/never start) but not as an explicit policy |
| Idle close 120 s | Not implemented | **Acceptable omission**: needs a timer owner that only exists at adapter/manager stage |
| Golden fixtures (Stage 1 deliverable, §17.14) | **Absent** — no `testdata/` fixtures; DTO duplication is real-encode/real-decode tested but not fixture-pinned | **Omission against the Stage 1 boundary**. Low risk (fake exercises the real codec paths), but the boundary names fixtures explicitly |

The report's §11.3 lists the omissions without conceding any required behavior as done. That treatment is
correct and should be preserved: nothing silently accepted. For the record: **no restart/backoff/idle-close
value in §17.8 may be cited as measured or guaranteed** — the provisional classification is intact.

---

## 7. Disposition and Stage 2 gate

| Area | Assessment |
|---|---|
| Main executable remains pure Go | **Pass** — reproduced `CGO_ENABLED=0` build |
| Lazy helper startup | **Pass** |
| Single process waiter | **Pass** |
| Protocol framing and correlation | **Pass in fake tests** |
| Cancellation/acceptance ordering | **Pass in fake tests** — barrier sound (§3) |
| WAV validation and cleanup | **Implemented; leaf checks strong; component-level reparse rule open (G-1)** |
| Existing Chat regression | **Pass** — reproduced |
| Real helper compatibility | **Not established** (report agrees, §11.2) |
| Restart/backoff/idle close | **Acceptably deferred as mechanism; fixtures missing (§6)** |
| Windows reparse-point hardening | **Incomplete — admitted; stays open** |

**Recommendation: CONDITIONAL ACCEPTANCE of Stage 1, gate stays HOLD.**

- **Required before the HOLD lifts:** fix F-1 (fatal-path delivery must not block; prove with a regression
  test that `failAll` with a full outcome channel terminates the session and still kills).
- **Required before lifting, low cost:** fix F-2 or record it as accepted behavior with a regression test
  pinning the chosen behavior; order the F-3 cleanup after kill-observed-exit on the violation path.
- **Required at a later explicit stage (not Stage 1):** per-component reparse-point checks (main side),
  golden fixtures, backoff/budget/idle-close with their owning stage, Shutdown/workMu serialization.
- **No change to:** MisakiSharp HOLD, no native provider, no pipeline integration, no Stage 2
  authorization. Real-helper compatibility remains unproven by construction — every protocol test uses the
  fake transport, and that is exactly what Stage 1 was scoped to prove.

---

## 8. Boundaries of this review

**Performed:** full read of all 12 Stage 1 files and all tests; `internal/proc` `Start`/`KillTree`
existence check (both present, `job_windows.go:84/166`); frozen-contract cross-check (§§17.8–17.11,
17.14–17.18, A1-D1–D7); hash and worktree verification; re-execution of every reported verification
command; one temporary adversarial probe, removed afterward (package directory verified back to 12 files).

**Not performed:** no implementation change, no fix, no new permanent test, no helper/fake modification,
no download, no native build, no pipeline or UI work, no commit or push.

**Not re-verified:** the pre-existing tracked diff (`interfaces.go`/`utils.go`) was confirmed unrelated
to Stage 1 by content, but its own correctness is outside this review's scope.
