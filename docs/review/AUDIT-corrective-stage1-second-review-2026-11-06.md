# Second Adversarial Review — Stage 1 Corrective Pass

**Date:** 2026-11-06
**Reviewer:** independent second pass (adversarial)
**Subject:** `docs/review/CORRECTIVE-howl-tts-helper-supervisor-stage1-2026-11-06.md`
**Subject SHA-256 (verified before reliance):** `5308ad5700b586ebc6c7764d654fa5bd770b128b52820ca8c9184a1b38a3e350` — MATCH (`sha256sum`, this environment)
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` §17 incl. A1/A1-D7; first adversarial review `AUDIT-howl-tts-helper-stage1-adversarial-2026-11-06.md`

**Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. Not an authorization. Stage 1 remains **HOLD**; Stage 2 is **not authorized**. No commit, push, pipeline, character, cache, UI, or native work was performed by this review.

**Method:** read the full corrective report; read every changed/new source and test file (`client.go`, `output.go`, `protocol.go`, `discovery.go`, `pathcheck_windows.go`, `pathcheck_other.go`, `supervisor.go`, `fake_test.go`, `regression_test.go`, `reparse_test.go`, `golden_test.go`, all 17 fixtures); re-ran the focused tests, full helper suite, race, vet, and a CGO-disabled build; grepped for cross-package wiring. No production file was modified; no test file was modified.

---

## 1. Task 1 — reproduce the reported verification

Environment: `go1.27.3 windows/amd64`, module `howl-chat`. All commands run this pass.

| Reported claim | Independent result |
|---|---|
| Focused set `Test(FailAllFullOutcome\|FailedHandshakeStopsOldReader\|RejectedCleanupTerminates\|GoldenFixtures\|WindowsReparse\|WindowsFinalSymlink)` exit 0 | **Reproduces.** F-1, F-2, F-3 PASS; `TestGoldenFixtures` PASS incl. all 17 subtests; `TestWindowsReparseComponentValidation` PASS; `TestWindowsFinalSymlinkValidation` **SKIP** (`A required privilege is not held by the client`) — exactly as reported |
| Full helper suite `go test ./internal/audio/tts/helper -count=1` ok | **Reproduces** (`ok … 1.298s`) |
| `go test -race ./internal/audio/tts/helper` ok | **Reproduces** (`ok … 3.900s`), no data races |
| `go vet ./internal/audio/tts/helper` exit 0 | **Reproduces** (no output, exit 0) |
| `CGO_ENABLED=0` build | **Reproduces** for the helper package (`go build ./internal/audio/tts/helper` exit 0). Full `go build ./...` was not re-run; the reported root/audio/race results are taken as stated and the helper-scope checks above all pass |

No evidence discrepancies. The report describes the code on disk.

---

## 2. Task 2 — corrected lifecycle paths: deadlock, reader leak, termination order

### F-1 (failAll deadlock) — FIXED, regression is genuine

- `pendingRequest.deliver` (`client.go`) is first-outcome-wins under `deliverMu` with a non-blocking send into the capacity-one channel; both `dispatch` and `failAll` use it. The fatal path can no longer block on a full outcome buffer, so `protocolViolation → Kill` is always reached.
- The regression test is stronger than it looks: the full channel is pre-filled *without* going through `deliver` (so `delivered=false`), forcing the exact hard case — non-blocking drop with the buffered success preserved. The empty channel receives the fatal; the pre-filled success is byte-identical afterward; exactly one `Wait` and a closed process are observed; `readFailure` records `protocol_violation`. The 1-second watchdog converts any regression into a hard failure.
- Residual note (not a finding): when a success was already delivered via the normal path, a later fatal is recorded in `readFailure` but the retired consumer keeps its success. That is the correct precedence (a completed synthesis should not be rewritten by a subsequent session failure) and matches the frozen single-terminal semantics.

### F-2 (failed-handshake reader leak) — FIXED, retry is single-reader

- Every non-incompatible handshake failure now runs the same sequence: `Terminate` (grace) → `waitForReader(readerDone)` → `SetState(unavailable)` → clear `started`/`transport`. Incompatible codes preserve their typed state instead of being overwritten by `observeExit`'s generic absent/crashed diagnostic — ordering is correct.
- The regression test pins all three properties: first process `Wait`ed once and closed, first `readerDone` closed (no competing reader survives), state is `unavailable`, and retry starts exactly one observer on helper 2. Reaching the retry exercises the cleared-transport path, so a second `readLoop` can never share the old stream.
- One environment-shaped observation: with the fake transport, `Terminate → Close` unblocks `Receive` immediately; with the real `osProcess`, `Receive` blocks on stdout until the process actually exits, and `Terminate`'s kill-on-grace-expiry is what releases it. The test therefore proves the *ordering logic*, not the wall-clock behavior against a stubborn native process. The logic is what F-2 challenged, so this is sufficient — recorded here so nobody over-claims it.

### F-3 (cleanup before termination) — FIXED on every path that deletes

- `cleanupRejected(dir, filename, force)` terminates the helper first whenever it is alive and (`force` or request-owned output exists), then removes only the caller-derived exact path plus same-directory `.part-*`. `Terminate` itself waits for the single exit observer under grace with kill fallback, so "exit-before-cleanup" holds through `Supervisor.Terminate` rather than a second ad-hoc wait — correct layering, single waiter preserved.
- Call-site coverage verified by reading `Synthesize`: `force=true` on all no-terminal paths (send failure, timeout, cancellation-grace expiry, stream failure, wire failure) and `force=false` with existence-check on terminal-error/invalid-result/WAV/acceptance paths. The send-failure branch — previously the uncovered hole — now calls `cleanupRejected(..., true)`.
- The wait-hook seam (`fake_test.go`) makes the regression test a true ordering proof: the file is present at the moment `Wait` returns and absent afterward. A delete-before-exit regression fails the test deterministically.
- Residual note: `cleanupRejected` consults `sup.ProcessAlive()`, which can transiently report alive between process death and observer run (first-review F-5, still open as INFO). Worst case is a redundant `Terminate` (nil-op when `proc==nil`) — fail-safe direction, no finding.

### F-4 / F-5 (deferred, informational) — untouched as authorized

- `Shutdown` still bypasses `workMu`; transient liveness misreporting between death and observation is unchanged. Both were explicitly outside the authorized corrective scope and remain open. The new code does not interact with either path in a way that worsens them: `Shutdown`'s call sites are unchanged, and nothing in the corrective diff reads `Status()` for control flow.

---

## 3. Task 3 — every path component, including the final file

**Code verdict: the walk unambiguously covers the final component.** `pathcheck_windows.go:17-51` starts at `filepath.Clean(abs)` — the target itself — and ascends to the volume root, applying *both* `os.Lstat` symlink rejection and `windows.GetFileAttributes → FILE_ATTRIBUTE_REPARSE_POINT` rejection to every component it visits. The final file is the *first* component checked, not an edge case. `pathcheck_other.go` mirrors this with the symlink check for non-Windows builds.

**Call-site verdict: all four deletion/validation gates use it.**

| Gate | Check | Fail-closed? |
|---|---|---|
| `regularHelperFile` (discovery) | `Lstat` + full walk; error → reject | Yes (`reparseErr == nil && !reparse`) |
| `validateOutput` (output dir) | full walk on `Clean(dir)`; error → reject | Yes (`err != nil \|\| reparse` → reject) |
| `validateWAVFile` (actual file) | `Lstat` leaf + full walk on final path; error → reject | Yes |
| `cleanupExact` (pre-delete) | full walk on `Clean(dir)`; error/reparse → return without deleting | Yes (refusal, not traversal) |

**Evidence verdict — the owner's distinction is preserved and correct:**

- **Observed Windows result:** intermediate-junction rejection PASS (`TestWindowsReparseComponentValidation`, reproduced this pass). Five properties proven against a real `mklink /J`: nested-path detection, helper rejection, output-dir rejection, validated-output rejection, and no-delete-traversal.
- **Not observed:** final-file symlink (`TestWindowsFinalSymlinkValidation` SKIP — privilege not held, reproduced this pass). The production final-component checks (`ModeSymlink` + reparse walk from the target itself) remain in the code path, so coverage is by construction, not by observation.
- **Do not upgrade the skip to a pass.** Any status claiming "final-file symlink validated on Windows" would be false. The correct statement is: final-file rejection is implemented on every path and unit-covered only insofar as the shared walk is exercised via the junction test; the Windows final-symlink execution remains open and requires a privileged runner.

Minor robustness notes (no finding): a nonexistent output directory fails `validateOutput` via the walk's `Lstat` error — fail-closed, correct for a directory that must exist; `cleanupExact`'s walk covers the directory rather than re-walking the joined final, but `validateWAVFile` already walked the exact final before acceptance and `cleanupExact` re-derives the final from caller parameters, so no wire-controlled component can enter between the two.

---

## 4. Task 4 — golden fixtures: meaningful and really exercised

- **All 17 files present** (`hello`, `catalog`, `load`, `synthesize`, `health`, `cancel`, `unload`, `shutdown` request+response, plus `error.response`) — the count and the Stage 1 §17.14 deliverable match.
- **Real strict codec, not a parallel schema:** `TestGoldenFixtures` globs `testdata/golden/*.json` (fails on zero files, fails on unrecognized names) and decodes every request through `decodeRequestFrame` and every response through `decodeResponseFrame` — the production `prepareFrame` (cap, NUL, BOM, UTF-8, duplicate-key) + `decodeStrict` (`DisallowUnknownFields`, trailing-data) path. Fixture IDs use the `req-<32hex>-<16digit>` shape the request-ID pattern enforces.
- **Result validators applied where semantics exist:** hello → `validateHello` (limits pinned to 4 MiB / 64 KiB); catalog → `validateCatalog` against piper params; load → `validateLoad` identity match; synthesize → `validateSynthesisResult` (format/codec/mono/rate/path-pin); health → `validateHealth`. Ack/empty shapes (cancel/unload/shutdown) and the `synthesis_failed` error envelope are shape-checked by the strict decoder, which is the correct level for them.
- **Not tautological:** the validators assert concrete values (e.g. synthesize output pinned to `C:/Howl/cache/tts` + `c0ffee.wav`, matching the fixture), so a fixture edit that breaks the contract fails the test rather than following it. The one asymmetry — requests are frame-checked but their `params` are not run through op-specific validators — is appropriate: the fixtures document the wire envelope; op validation already has unit coverage elsewhere.
- Genuine gap, correctly scoped: all fixtures are still parent-side vectors (no native helper exists to interop against). Fixture presence does not equal helper compatibility — the corrective report states this (§10.5), and it remains true.

---

## 5. Task 5 — no regression against the frozen contract or Stage 1

- **Stage 0 / frozen contract:** single-`Wait` ownership intact (exactly one observer per successful start; failed start creates none — unchanged); no autonomous restart (no new spawn path; `observeExit` still only records); session rotation on replacement intact; strict-codec invariants untouched (`protocol.go` diff is additive `validateOutput` hardening only); public error projection unchanged (`protocol_violation` internal, `ErrCodeSynthesisFailed` compatibility under A1-D7 — `deliver` change preserves this). `TTSSynthesizer` interface untouched.
- **Stage 1 boundary:** the package remains unwired — repo-wide grep for `audio/tts/helper` imports finds references only inside the package's own files and the review docs. No manager/pipeline/character/cache/binding caller exists. CGO-disabled build of the package passes; no native, Kokoro, G2P, ORT, model, or asset dependency introduced (import surface unchanged apart from `golang.org/x/sys/windows` behind the `windows` build tag plus stdlib).
- **Full-suite signal:** helper suite, race, and vet all pass this pass (see §1). The broader `go test ./...` / audio-tree / root signals from the corrective report (§8) were not all re-run here; the helper-scope evidence plus unchanged cross-package surface (no new imports, no interface changes) gives no reason to suspect a regression outside the helper, but this review asserts that narrowly — it does not re-certify the audio tree or root suite.

---

## 6. Disposition and remaining items

| Item | Assessment |
|---|---|
| F-1 `failAll` deadlock | **Fixed; regression genuine.** Hold-lifting requirement (a) satisfied pending owner acceptance of this review |
| F-2 handshake reader leak | **Fixed; single-reader retry proven.** Requirement (b) satisfied |
| F-3 cleanup ordering | **Fixed on all deleting paths; ordering proven.** Requirement (c) satisfied |
| Reparse-point validation | **Implemented on every component incl. final; Windows-observed for intermediate junction only.** Final-symlink execution stays open — requires a privileged runner; must not be recorded as passed |
| Golden fixtures | **Supplied; real-codec exercised; meaningful.** §17.14 deliverable closed |
| F-4 / F-5 | Open informational, untouched as authorized |
| Restart/backoff/idle-close | Still deferred, correctly unclaimed |
| Real-helper compatibility | Still unproven by construction (fake transport only) |
| Stage 2 / native / pipeline / character / cache / UI | Untouched; no authorization implied |

**Recommendation: LIFT THE HOLD on Stage 1 — accept the corrective pass — with two carried conditions:** (1) the final-file Windows symlink test remains an open execution item for any privileged CI/runner and must not be reported as passed until then; (2) F-4/F-5, restart/backoff/idle-close ownership, and real-helper compatibility transfer explicitly to their owning later stages. No Stage 2 authorization follows from this review; native provider work, pipeline/character integration, cache changes, and production adoption remain outside authorization until separately granted.

## 7. Boundaries of this review

**Performed:** subject-hash verification; full read of all changed/new helper sources and tests; re-execution of the focused set, full helper suite, race, vet, and CGO-disabled package build; cross-package wiring grep; all 17 fixture contents inspected.

**Not performed:** no implementation change, no fix, no new test, no download, no native build, no pipeline/UI work, no commit or push, no full `go build ./...`, no audio-tree or root-suite re-run (see §5 for the narrowed assertion), no real-helper interop (none exists).
