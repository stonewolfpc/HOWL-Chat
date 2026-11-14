# Stage 1 Corrective Pass — TTS Helper Supervisor and Protocol Client

**Date:** 2026-11-06  
**Scope:** Authorized Stage 1 corrective pass only  
**Governing design:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` §17, including A1 and A1-D7  
**Supersedes for corrective evidence:** the corrective portions of `docs/review/IMPLEMENTATION-howl-tts-helper-supervisor-stage1-2026-11-06.md`  
**Adversarial review addressed:** `docs/review/AUDIT-howl-tts-helper-stage1-adversarial-2026-11-06.md`  
**Audit SHA-256:** `03ed0a4676089a0a3086bbbfeaf36f9159e6cca3bd9ec9b6d63121cfa68680f0`  
**Status:** CORRECTIVE PASS SUBMITTED FOR INDEPENDENT REVIEW — STAGE 1 REMAINS ON HOLD. NO COMMIT, PUSH, OR STAGE 2 WORK.

---

## 1. Scope and disposition

This pass is limited to the findings authorized by the owner:

| Finding | Disposition in this pass | Status |
|---|---|---|
| F-1: `failAll` deadlock on a full outcome channel | Non-blocking first-outcome delivery; deterministic full-channel regression; termination and classified delivery verified | Corrected, pending independent review |
| F-2: failed-handshake reader leak | Terminate the helper, observe exit, stop the existing reader, clear transport before retry | Corrected, pending independent review |
| F-3: cleanup before termination/observed exit | Cleanup helper terminates and waits before exact-path deletion when output may exist or no terminal arrived | Corrected, pending independent review |
| Windows reparse-point validation | Per-component `FILE_ATTRIBUTE_REPARSE_POINT` walk on Windows; shared checks use it | Corrected, pending independent review |
| Golden fixtures | 17 protocol fixtures under `testdata/golden`, decoded by the real strict codec | Supplied, pending independent review |

No Stage 2 work, native provider, Kokoro/G2P/ORT dependency, pipeline integration, character integration, cache change, or UI change was performed. The main executable remains `CGO_ENABLED=0` and contains no native helper dependency.

The restart/backoff and idle-close mechanisms remain deferred. This report does not claim that the provisional 1/2/5-second backoff, 3/10-minute restart budget, or 120-second idle close are measured, implemented, or guaranteed. The existing F-4 (`Shutdown`/`workMu` serialization) and F-5 (transient liveness observation) informational findings also remain outside this corrective pass.

---

## 2. F-1 — non-blocking fatal delivery

### Implementation

`pendingRequest.deliver` now publishes only the first outcome and never blocks the protocol reader:

```go
func (r *pendingRequest) deliver(out pendingOutcome) {
    if r == nil {
        return
    }
    r.deliverMu.Lock()
    defer r.deliverMu.Unlock()
    if r.delivered {
        return
    }
    r.delivered = true
    select {
    case r.ch <- out:
    default:
    }
}
```

Both `dispatch` and `failAll` use `deliver`. The outcome channel remains capacity one; the correction preserves the concurrency invariant rather than increasing capacity. `Client.readFailure` is recorded before delivery, so later operations still fail fast even when a terminal already occupies the buffer. `protocolViolation` calls `failAll` and then reaches `Supervisor.Kill` without an unbounded send.

The public error projection is unchanged: `protocol_violation` remains the precise internal `Failure.Code`; `ErrCodeSynthesisFailed` remains only the compatibility projection under A1-D7.

### Regression evidence

`TestFailAllFullOutcomeDoesNotBlockOrSkipKill` in `regression_test.go`:

1. starts a fake helper through the supervisor;
2. fills a capacity-one request channel with a success outcome;
3. installs a second pending request with an empty channel;
4. triggers `protocolViolation` as a second fatal event;
5. requires the call to return within one second;
6. verifies exactly one `Process.Wait` call and a closed process;
7. verifies `readFailure.Code == "protocol_violation"`;
8. verifies the empty request receives the fatal outcome;
9. verifies the already-buffered success is preserved.

This is a deterministic in-package regression, not an inference from channel capacity.

---

## 3. F-2 — failed-handshake reader lifecycle

### Implementation

`Client.Start` now handles every failed `helloHandshake` the same way:

1. terminate the existing helper with the configured shutdown grace;
2. wait for the existing `readerDone` observer;
3. preserve `protocol_incompatible`/`protocol_minor_unsupported` as `incompatible`, or record other handshake failures as `unavailable`;
4. clear `started` and `transport` only after termination and reader observation.

The previous behavior left a non-incompatible failed handshake alive in `handshaking`; a retry could then create a second reader over the same stream. The new ordering prevents that. State preservation occurs after termination so `observeExit` cannot replace an incompatibility diagnostic with a generic absent state.

### Regression evidence

`TestFailedHandshakeStopsOldReaderAndAllowsRetry`:

- uses two independent fake helpers and a factory that returns a new process for the second start;
- makes helper 1 reject `hello` with `invalid_request`;
- verifies one `Wait` on helper 1 and a closed process;
- verifies the first reader's `readerDone` channel is closed;
- verifies the failed-handshake state is `unavailable`;
- retries `Start` and verifies helper 2 starts and has exactly one `Wait` observer.

The test therefore covers the failed handshake, old-reader shutdown, and successful retry without competing readers.

---

## 4. F-3 — termination before rejected-output cleanup

### Implementation

`Client.cleanupRejected(dir, filename, force)` enforces the ordering:

```text
if helper is alive and (force || request-owned output exists):
    terminate helper
    wait for the single process exit observer
then:
    remove only the caller-derived exact output and .part-* paths
```

`force` is true when no terminal arrived (send failure, timeout, cancellation grace expiry, stream failure), because the helper could still create output. For terminal errors and invalid results, `force` is false, but termination still occurs when request-owned output exists.

The helper is now used by the send-failure branch, timeout, cancellation, stream failure, wire error, decode failure, synthesis-result validation failure, WAV validation failure, and acceptance-rejection paths.

### Regression evidence

`TestRejectedCleanupTerminatesBeforeRemovingOutput` starts a real fake process through the supervisor, creates a valid rejected WAV, installs a `Wait` observer hook, and calls `cleanupRejected(..., false)`. The hook proves the file is still present at the moment process exit is observed; the test then verifies the file is removed. The test fails if cleanup deletes the file before `Wait` returns.

---

## 5. Windows reparse-point validation

### Implementation

`pathcheck_windows.go` performs a component-by-component walk from the target path to the volume root. For each existing component it:

- rejects `os.ModeSymlink` via `Lstat`;
- calls `windows.GetFileAttributes`;
- rejects `windows.FILE_ATTRIBUTE_REPARSE_POINT`.

It therefore covers intermediate junctions, final reparse points, and ordinary symlinks. `pathcheck_other.go` provides the corresponding non-Windows symlink walk so shared code compiles everywhere without importing Windows-only packages.

The validation is used by:

- `regularHelperFile` in discovery;
- `validateOutput` for the output directory;
- `validateWAVFile` for the actual output path;
- `cleanupExact` before any deletion.

The `filepath.Rel` containment check and caller-derived cleanup path are unchanged; a helper cannot redirect deletion through wire-supplied `output_path`.

### Windows evidence

`TestWindowsReparseComponentValidation` created a real junction with:

```text
cmd /c mklink /J <temp>\junction <temp>\target
```

It then verified that:

- a nested path beneath the junction is rejected by `hasReparseComponent`;
- a helper beneath the junction is rejected by `regularHelperFile`;
- an output directory beneath the junction is rejected by `validateOutput`;
- output validation rejects a file beneath the intermediate junction;
- exact cleanup does not traverse the junction to delete a managed target file.

Result: **PASS** on `go1.27.3 windows/amd64`.

`TestWindowsFinalSymlinkValidation` attempted to create a final-file symlink. This account does not hold the Windows symlink privilege, so the test **SKIPPED** with `A required privilege is not held by the client`. The existing final-component `os.ModeSymlink` check and the per-component reparse walk remain in the production path; the final-file symlink case was not executed in this environment and is not claimed as observed evidence. The intermediate-junction test does run and is the new Windows path-validation evidence.

---

## 6. Golden fixtures

The Stage 1 deliverable required by §17.14 is now present under:

```text
internal/audio/tts/helper/testdata/golden/
```

The set contains 17 fixtures covering request and response envelopes:

```text
hello.request.json
hello.response.json
catalog.request.json
catalog.response.json
load.request.json
load.response.json
synthesize.request.json
synthesize.response.json
health.request.json
health.response.json
cancel.request.json
cancel.response.json
unload.request.json
unload.response.json
shutdown.request.json
shutdown.response.json
error.response.json
```

`TestGoldenFixtures` reads every fixture through the real strict `decodeRequestFrame`/`decodeResponseFrame` path and additionally applies the relevant hello, catalog, load, synthesis, and health result validators. The fixtures are not a parallel codec or test-only schema.

---

## 7. Exact changed files in this corrective pass

All files below are either modified relative to the original Stage 1 submission or newly added by this pass. The helper package is untracked in the worktree; no commit or push occurred.

| File | Change |
|---|---|
| `internal/audio/tts/helper/client.go` | F-1 non-blocking delivery; F-2 reader lifecycle and state preservation; F-3 `cleanupRejected` and send-failure use |
| `internal/audio/tts/helper/output.go` | `requestOutputExists`; component reparse checks in WAV validation and exact cleanup |
| `internal/audio/tts/helper/protocol.go` | `validateOutput` rejects reparse components |
| `internal/audio/tts/helper/discovery.go` | `regularHelperFile` rejects reparse components |
| `internal/audio/tts/helper/pathcheck_windows.go` | New Windows per-component reparse walk |
| `internal/audio/tts/helper/pathcheck_other.go` | New non-Windows symlink walk |
| `internal/audio/tts/helper/fake_test.go` | `waitHook` seam for exit-ordering evidence |
| `internal/audio/tts/helper/regression_test.go` | New F-1, F-2, and F-3 regression tests |
| `internal/audio/tts/helper/reparse_test.go` | New Windows junction and final-symlink tests |
| `internal/audio/tts/helper/golden_test.go` | New real-codec fixture test |
| `internal/audio/tts/helper/testdata/golden/hello.request.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/hello.response.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/catalog.request.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/catalog.response.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/load.request.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/load.response.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/synthesize.request.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/synthesize.response.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/health.request.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/health.response.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/cancel.request.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/cancel.response.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/unload.request.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/unload.response.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/shutdown.request.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/shutdown.response.json` | New golden fixture |
| `internal/audio/tts/helper/testdata/golden/error.response.json` | New golden fixture |
| `docs/review/CORRECTIVE-howl-tts-helper-supervisor-stage1-2026-11-06.md` | This report |

No production file outside `internal/audio/tts/helper` was modified by this corrective pass. The pre-existing tracked diff in `internal/audio/types/interfaces.go` and `internal/audio/types/utils.go` remains unrelated character/pipeline work and was not edited.

---

## 8. Verification results

Environment: `go1.27.3 windows/amd64`, module `howl-chat`.

Focused corrective tests:

```text
go test ./internal/audio/tts/helper -run 'Test(FailAllFullOutcome|FailedHandshakeStopsOldReader|RejectedCleanupTerminates|GoldenFixtures|WindowsReparse|WindowsFinalSymlink)' -count=1 -v
```

Result: **exit 0**. F-1, F-2, F-3, and `TestGoldenFixtures` passed; the intermediate-junction test passed; the final-file symlink test skipped for the privilege reason above.

Full helper suite:

```text
go test ./internal/audio/tts/helper -count=1
ok  howl-chat/internal/audio/tts/helper  2.342s
```

Race:

```text
go test -race ./internal/audio/tts/helper -count=1
ok  howl-chat/internal/audio/tts/helper  4.112s
```

Vet:

```text
go vet ./internal/audio/tts/helper
exit 0
```

No-CGO build:

```text
CGO_ENABLED=0 go build ./...
exit 0
```

Audio regression:

```text
go test ./internal/audio/... -count=1
```

All 16 audio packages passed, including `internal/audio/tts/helper`.

Root regression:

```text
go test . -count=1
ok  howl-chat  11.547s
```

A broader `go test ./... -count=1` run was also completed during this corrective pass with exit 0; the final state-preservation adjustment was then followed by the helper, race, vet, no-CGO build, audio, and root checks above.

---

## 9. Updated artifact hashes

Implementation and test artifacts:

| Artifact | SHA-256 |
|---|---|
| `internal/audio/tts/helper/client.go` | `e2b2e713f3e44d0e93392aa2508ef9ff512cc0a625848664f7a949dff1e72e77` |
| `internal/audio/tts/helper/output.go` | `c50fce52ab901f6a16500522132cea58000b3d119b8dd8bb9a67186a608c10c0` |
| `internal/audio/tts/helper/protocol.go` | `8a9f710b4b141de342be88dafb4e046cb6c6b095507a70faa58db7fac7f48780` |
| `internal/audio/tts/helper/discovery.go` | `af785db5a7b4bcc802cadfacc8ad60c57dc3f0210745651086f9cfd336369732` |
| `internal/audio/tts/helper/pathcheck_windows.go` | `a3747e78f211acd3a5ec1b492877511673f7b47c4410a032d5a1b68d135e66ab` |
| `internal/audio/tts/helper/pathcheck_other.go` | `e799d9c9520cdcb4de1fdd80c93dd2f38d5a6132bddd00d8d8a845738415b8fa` |
| `internal/audio/tts/helper/fake_test.go` | `9ec259bbd339f02cccc92496a896c45c8f5bbc2a43e337c2ba8f9a4bc0ef4f25` |
| `internal/audio/tts/helper/regression_test.go` | `27be858a52a62bfec4b7cd244bb75486535fcba5e706aec22e2ec4fbb0514e4e` |
| `internal/audio/tts/helper/reparse_test.go` | `866885890279e29848068af240561bfa93ee6e3e170b65b3c8930b666a62ff61` |
| `internal/audio/tts/helper/golden_test.go` | `4cbecad7ee372f13f7916329aa2bbb630c26d5f1e2a2976bb791606ebc2eba09` |

Golden fixtures:

| Artifact | SHA-256 |
|---|---|
| `testdata/golden/cancel.request.json` | `bbdde39d46aaf55a70e56fd94f3e804b32f11fb6770ee043865b37ecaf5c9cab` |
| `testdata/golden/cancel.response.json` | `504c44d3e140cb3e9238a520ad8da1188933dbb435c738fc10d1fe3fe1e1431a` |
| `testdata/golden/catalog.request.json` | `b73b5c2f5b3b886eb8965fa1dd6ad9e56a024da5820bd77c8f1d687e1f7151c9` |
| `testdata/golden/catalog.response.json` | `34b2eb7e46abf1275e2f42e76dce0d7d5f275133f0a0cde648f15d52e3d1d747` |
| `testdata/golden/error.response.json` | `0df9a5ac852b642eaba183b24e1f7fa0ecbc37019ef044ac4d63d7a58ff51a77` |
| `testdata/golden/health.request.json` | `544905aefa8fb0b9548ab727ca383798e4f1beaa04832e990928e1b95e4bc8c9` |
| `testdata/golden/health.response.json` | `ddaa47cd4b71433110f07e5719633059fc8dc2aba9c9c87b005f7b9727fe3417` |
| `testdata/golden/hello.request.json` | `aa5bc590c97034463978ded4dc92094a49cf13e17a9c6df33c9e8ac4b1bf758d` |
| `testdata/golden/hello.response.json` | `bcb195f0eead06e3a01ec56b613370a21a848f2844d841b0d4512f5bcade7fc3` |
| `testdata/golden/load.request.json` | `01ce6edef087acd3140b7afd7db052ffcb1aee0afc0d82dc51af1d457c262e38` |
| `testdata/golden/load.response.json` | `fc0f556089e99f63f2789bcc29d21c4b5eae0db31512f0178de166abf6b8fca2` |
| `testdata/golden/shutdown.request.json` | `e2f5b7a8b7033a01b15219212d3e2ba16731d592c32e28b7d75a001f15a28ffe` |
| `testdata/golden/shutdown.response.json` | `cc463977b2208f62456a0b12ffe84f3d10646455cf015c870a17da06f4a1cd72` |
| `testdata/golden/synthesize.request.json` | `dc3a7e397c017d33c466ef7a5b32a64af8db15ddd3ee8e2328abc8199fda7156` |
| `testdata/golden/synthesize.response.json` | `da7ed4db61ac10036c42d494fecac0deeb6794d4097dcbbf367e16a5b658b0f4` |
| `testdata/golden/unload.request.json` | `5c46f5a77e834b122246f939fa2f6a6e5e4f3cffa3b64c19684c1e214c87bf19` |
| `testdata/golden/unload.response.json` | `e59f3166248a972d81414d98c36bc5752117993bad66d477b1e506269395d0c3` |

The original implementation report remains `IMPLEMENTATION-howl-tts-helper-supervisor-stage1-2026-11-06.md` with SHA-256 `1e8e6553a6e5c2f71f391940fe6ef670fbcf73d4c9381b23b066eeb639705d0a`. This report's own SHA-256 is recorded in the handoff response.

---

## 10. Remaining limitations and review status

1. **Stage 1 is not accepted.** The corrective pass is submitted for independent review; the HOLD remains in force.
2. **Restart/backoff and idle close remain deferred.** The 1/2/5-second schedule, 3/10-minute budget, and 120-second idle close are not implemented and must not be described as measured or guaranteed behavior.
3. **Final-file symlink test was not executed.** Windows denied the test symlink privilege; intermediate-junction validation and production final-component checks passed. This is an environment limitation, not a claim of observed final-file symlink evidence.
4. **F-4 and F-5 remain open informational findings.** `Shutdown`/`workMu` serialization and transient liveness observation were outside this authorized correction and were not changed.
5. **No real helper compatibility evidence.** All protocol tests still use the fake process/transport, as intended for Stage 1.
6. **No commit or push.** All helper files and review documents remain local and uncommitted.

**Requested next step:** independent review of this corrective report and the exact files above. Do not begin Stage 2, native provider work, or pipeline/character integration.

---

**END — STAGE 1 CORRECTIVE REPORT**
