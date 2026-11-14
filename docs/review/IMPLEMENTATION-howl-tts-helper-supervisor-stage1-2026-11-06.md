# Stage 1 — Main-Side TTS Helper Supervisor and Protocol Client

**Date:** 2026-11-06  
**Scope:** Stage 1 only  
**Governing design:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` §17, including A1 and A1-D7  
**Status:** SUBMITTED FOR INDEPENDENT REVIEW — NO COMMIT, PUSH, OR STAGE 2 WORK

---

## 1. Implementation summary

The main application now has an isolated, provider-neutral TTS helper package at
`internal/audio/tts/helper`. It provides:

- ordered lazy discovery;
- Windows job-owned process startup through `internal/proc.Start`;
- exactly one supervisor-owned `Wait()` observer after successful startup;
- Protocol v1.0 NDJSON framing, strict decoding, duplicate-key rejection, version checks,
  request ID generation, and response correlation;
- a dedicated protocol reader with control-message handling while synthesis is active;
- provider/process status without model loading;
- cancellation, acceptance-barrier, timeout, crash, and shutdown paths;
- exact-path WAV validation and cleanup;
- precise internal `helper.Failure` values with an explicit compatibility projection to
  `types.AudioError` (A1-D7).

No audio pipeline, character, cache, renderer, or UI code was changed. The package is not yet
wired into `audio.Manager`, the pipeline, or a Wails binding. It has no native, Kokoro, G2P, ORT,
model, or asset dependency.

## 2. Exact changed files

All implementation files are new; no tracked production file was modified by Stage 1.

| File | Purpose |
|---|---|
| `internal/audio/types/provider.go` | Additive provider inspector/status/catalog types; leaves `TTSSynthesizer` unchanged |
| `internal/audio/tts/helper/errors.go` | Precise `Failure` classification and public `AudioError` compatibility projection |
| `internal/audio/tts/helper/protocol.go` | Protocol envelopes, strict decoding, duplicate-key and frame checks, DTO validation |
| `internal/audio/tts/helper/discovery.go` | Three-rule absolute-path discovery with no PATH fallback |
| `internal/audio/tts/helper/supervisor.go` | Lazy process start, job ownership, one `Wait()` observer, stderr ring, status |
| `internal/audio/tts/helper/client.go` | Protocol client, reader/correlation, operations, acceptance barrier, lifecycle |
| `internal/audio/tts/helper/output.go` | Canonical WAV/path/hash validation and exact-path cleanup |
| `internal/audio/tts/helper/*_test.go` | Fake process/helper transport, protocol, discovery, supervisor, and client tests |
| `docs/review/IMPLEMENTATION-howl-tts-helper-supervisor-stage1-2026-11-06.md` | This report |

The worktree already contained unrelated untracked review documents and other pre-existing changes.
They were not modified by this stage. A scoped `git status` is included in §13.

## 3. Architecture-to-code mapping

| Frozen requirement | Implementation |
|---|---|
| Discovery order: env, app runtime, development runtime | `Discoverer.Discover` in `discovery.go` |
| No PATH/System32/current-directory search | `Discoverer` only constructs the three explicit candidates |
| Helper not found is TTS-only | `Supervisor.Start` records `helper_missing` / `unavailable`; no process is spawned |
| Suspended start and job ownership before resume | `StartProcess` calls `proc.Start`; `internal/proc/job_windows.go` owns suspend/attach/resume |
| One `Wait()` owner after successful `Start()` | `Supervisor.Start` starts only `Supervisor.observeExit`; `osProcess.Wait` also caches |
| Failed start creates no waiter | `StartProcess` returns before observer creation; `proc.Start` owns its failure wait |
| NDJSON strict framing | `readFrame`, `prepareFrame`, `encodeRequestFrame` |
| Duplicate JSON keys fatal | `validateDuplicateKeys`; `prepareFrame` returns `invalid_frame` |
| Unknown core fields rejected | `decodeStrict` with `DisallowUnknownFields` |
| Handshake/version negotiation | `helloParams`, `helloResult`, `validateHello`; exact v1.0 response protocol |
| Request IDs and retired IDs | `newRequestID`, `Client.pending`, `Client.retired` |
| Reader must not block on synthesis | `Client.readLoop` is a dedicated goroutine; control calls do not hold `workMu` |
| One active non-control operation | `Client.workMu` around catalog/load/unload/synthesize |
| One terminal per request | `pendingRequest.terminal`; duplicate terminal is `protocol_violation` |
| Parent acceptance barrier | `requestAcceptance` plus the atomic state check in `Client.Synthesize` |
| Cancellation precedes acceptance | `watchCancellation`, `cancel_requested`, grace, kill/wait, exact cleanup |
| Helper publication barrier | Main only validates after terminal; it does not ask helper to mutate after terminal |
| `protocol_violation` retains precise code | `Failure.Code`; `AudioError.Cause`; A1-D7 mapping is a projection only |
| Provider status without model load | `Supervisor.Status`, `Client.ProviderStatus`; tests confirm zero loads |
| No unsolicited events | `decodeResponseFrame` requires a terminal response shape |

## 4. Process lifecycle and state transitions

The implemented supervisor path is:

```text
absent
  -> starting             (explicit request, path resolved)
  -> handshaking          (proc.Start succeeded; Wait observer started)
  -> unloaded             (hello validated)
  -> loading              (load request sent)
  -> ready                (load result validated)
  -> busy                 (synthesize request sent)
  -> ready / degraded     (terminal accepted / non-fatal failure)
  -> stopping             (shutdown or forced termination)
  -> absent or exited     (normal exit / non-zero exit)
```

Unexpected process exit from any live state is observed by the single waiter and records `crashed` with
`ProcessAlive=false`. Protocol errors, malformed frames, unknown IDs, and duplicate terminals call the
same fatal protocol path and terminate the helper before cleanup.

A replacement is started only by a later explicit client call after the supervisor reports the old
process no longer alive. Stage 1 does not implement the deferred backoff budget, automatic restart
schedule, or idle-close timer; those remain later work and are listed in §11.

## 5. Protocol operation coverage

| Operation | Client coverage | Validation |
|---|---|---|
| `hello` | `Client.Start` / `helloHandshake` | exact v1.0, limits, providers, capabilities, extensions |
| `health` | `Client.Health` | allowed state, zero queue depth, active ID only while synthesizing |
| `catalog` | `Client.Catalog` / `ProviderCatalog` | provider/engine echo, catalog version 1, speed range, voice uniqueness |
| `load` | `Client.Load` | provider/engine/identity echo, loaded=true, runtime present |
| `unload` | `Client.Unload` | `unloaded=true` |
| `synthesize` | `Client.Synthesize` | text/output limits, strict result, WAV, hash, metadata, identity, exact path |
| `cancel` | `Client.Cancel` | separate request ID, target echo, acknowledgement |
| `shutdown` | `Client.Shutdown` | accepted=true, then supervisor termination |

No second protocol, extension operation, unsolicited event, HTTP transport, named pipe, or audio byte
transport was added.

## 6. Cancellation and acceptance evidence

The parent acceptance state is explicit: `active`, `cancel_requested`, `cancelled`, `timeout`,
`validating`, `accepted`, and rejected/retired. A success is returned only after:

1. the success terminal is received;
2. the result DTO and actual WAV are validated;
3. the per-request lock re-checks cancellation and timeout;
4. `accepted` is committed.

Cancellation before the terminal is sent with its own `cancel` request ID. If no terminal arrives
within the configured grace, the supervisor closes stdin, waits for the sole exit observer, kills the
tree if needed, and only then removes the exact final and `.part-*` paths. A late success cannot revive
the request. A success accepted before cancellation remains returned; later cancellation does not
retract it.

The tests `TestCancellationBeforeTerminal`, `TestCancellationAfterTerminalBeforeAcceptance`,
`TestAcceptedResultIsNotRetractedByLaterCancellation`, `TestTimeoutForcesTermination`, and
`TestControlMessagesDuringSynthesis` exercise these paths.

## 7. Test inventory and exact results

Command:

```text
go test ./internal/audio/tts/helper -count=1
```

Result:

```text
ok  howl-chat/internal/audio/tts/helper  2.199s
```

Tests executed:

```text
TestClientLifecycleAndCatalog
TestConstructionDoesNotStartHelper
TestShutdownWithActiveRequest
TestSynthesisSuccess
TestCancellationBeforeTerminal
TestCancellationAfterTerminalBeforeAcceptance
TestAcceptedResultIsNotRetractedByLaterCancellation
TestTimeoutForcesTermination
TestControlMessagesDuringSynthesis
TestMalformedFrameIsFatal
TestUnknownResponseIDIsProtocolViolation
TestMissingHelperIsOptional
TestNoGoroutineLeakAfterShutdown
TestProtocolIncompatibleStopsHandshake
TestDiscoveryUsesEnvironmentFirst
TestDiscoveryUsesApplicationThenDevelopment
TestDiscoveryRejectsRelativeEnvironmentPath
TestDiscoveryMissingDoesNotUsePATH
TestDuplicateKeysAreFatal
TestUnknownResponseFieldIsRejected
TestOversizedFrameIsRejected
TestRequestFrameCarriesExactVersionAndNoExtensions
TestRetiredResponseIsDiscarded
TestDuplicateTerminalIsProtocolViolation
TestSynthesisValidationRejectsEscapingOutput
TestSupervisorUsesOneWaitObserver
TestSupervisorObservesCrash
TestSupervisorFailedStartCreatesNoWaiter
TestHelperProcess
TestStartProcessUsesRealStartupPath
```

`TestHelperProcess` is the child-mode branch of the process test double used by
`TestStartProcessUsesRealStartupPath`; it is not a helper executable and is not part of the production
build.

## 8. Race results

Command:

```text
go test -race ./internal/audio/tts/helper ./internal/audio/types -count=1
```

Result:

```text
ok  howl-chat/internal/audio/tts/helper  3.866s
ok  howl-chat/internal/audio/types       1.486s
```

`go vet ./internal/audio/types ./internal/audio/tts/helper` completed with no diagnostics.

## 9. CGO-disabled main build

Command:

```text
CGO_ENABLED=0 go build -o build/howl-chat-stage1-check.exe .
```

Result:

```text
build_exit=0
```

The temporary build artifact was removed after verification. No ONNX Runtime, Kokoro, G2P, native
helper, model, or asset dependency was introduced.

## 10. Chat regression results

Command:

```text
go test . -count=1
```

Result:

```text
ok  howl-chat  11.062s
```

Command:

```text
go test ./internal/audio/... -count=1
```

Result: all audio packages passed:

```text
howl-chat/internal/audio
howl-chat/internal/audio/asr
howl-chat/internal/audio/asr/fallback
howl-chat/internal/audio/asr/llama
howl-chat/internal/audio/asr/qwen
howl-chat/internal/audio/asr/whisper
howl-chat/internal/audio/capability
howl-chat/internal/audio/hw
howl-chat/internal/audio/runtime
howl-chat/internal/audio/tts
howl-chat/internal/audio/tts/cache
howl-chat/internal/audio/tts/helper
howl-chat/internal/audio/tts/kokoro
howl-chat/internal/audio/tts/melotts
howl-chat/internal/audio/tts/piper
howl-chat/internal/audio/types
```

No application startup path calls the helper. The package is inert until an explicit client operation
is invoked.

## 11. Known limitations and unresolved findings

1. **No pipeline/UI integration.** Stage 1 deliberately does not implement `TTSSynthesizer`, connect
   the helper to `audio.Manager`, change character voice resolution, alter the cache, or modify
   renderer/Wails signatures.
2. **No real helper compatibility evidence.** Every protocol test uses a fake process/transport. The
   tests prove the parent-side contract, not compatibility with a future native helper.
3. **No restart budget/idle close.** A crash is observed and active work fails; a later explicit
   request may start a replacement. The frozen 3/10-minute restart budget, 1/2/5-second backoff
   schedule, and 120-second idle close are not implemented in Stage 1.
4. **Windows reparse-point coverage is partial.** `os.Lstat` rejects symlinks in discovery and output;
   a full `FILE_ATTRIBUTE_REPARSE_POINT` check for every path component is left for the Windows
   hardening pass. The helper itself must also repeat containment checks.
5. **Status is lifecycle-derived.** `ProviderStatus` reports the supervisor's handshake/load state;
   it does not issue a fresh `health` call or merge a helper health snapshot.
6. **Synthesis output allocation is caller-owned.** Stage 1 requires an absolute output directory and
   a validated basename. The pipeline's future adapter must allocate those paths.
7. **No `SynthesizeStream`.** Protocol v1 is file-only; this package does not advertise streaming.
8. **Job-object behavior is inherited, not reimplemented.** `StartProcess` uses the established
   `internal/proc.Start` path; existing `internal/proc` tests cover suspended start, assignment, and
   resume. The Stage 1 test process verifies actual startup/exit through that path but is not a native
   helper.
9. **No commit or push.** All files remain local and uncommitted.

## 12. Artifact hashes

Implementation artifact SHA-256 values at report time:

| Artifact | SHA-256 |
|---|---|
| `internal/audio/types/provider.go` | `cc67286f14406ca1f8420739544924fdd085bfed294f4aa8f967540163f5489e` |
| `internal/audio/tts/helper/client.go` | `f1f798dc1043a0440665af2d35e9d5577fc5705ac2269d804941f9b9129b52c6` |
| `internal/audio/tts/helper/client_test.go` | `d4f6eb085929b1ac429bd733d10d86a77d3f36f37f08a4adeab2a644904937a2` |
| `internal/audio/tts/helper/discovery.go` | `3f6dd9f8395c8a053e241bddd29e62ed5c2f0d5e066811cd2076ff8d5ec4552f` |
| `internal/audio/tts/helper/discovery_test.go` | `e5a6953fd97c6bb57a161ffa957d02ff9ec4e72809e43370887f732f5e32736e` |
| `internal/audio/tts/helper/errors.go` | `faca46d0051ee5a3400f83390e892454ba6a7950fe40e3934115c2b63e17662e` |
| `internal/audio/tts/helper/fake_test.go` | `2091a0906ddf179a20b8fd800991da34d78f4e33c1598ccee73cd00ceaa43f9a` |
| `internal/audio/tts/helper/output.go` | `a2575f53588fc46e150d117cbf006c2940c08bc4b707905edddcf987d31e3575` |
| `internal/audio/tts/helper/protocol.go` | `6c373d4da65ace746d2555932db260c4e42255046871274a6b4f8fb7f0dbac6d` |
| `internal/audio/tts/helper/protocol_test.go` | `5ce7b3ff331d9767f921fadb76badfed8eb6da1dc3aab516f9b4aa240ad70594` |
| `internal/audio/tts/helper/supervisor.go` | `3dd467d3d26323d47afbda1d2b8fec9018e9ab90115c91c05fa2fd92802718f3` |
| `internal/audio/tts/helper/supervisor_test.go` | `5a76a24bb7063c83be0db903fb5519250915efff21270d2d206c87ec77caa0a6` |

The report file's own SHA-256 is recorded in the handoff response, because a file cannot contain a
stable hash of itself.

## 13. Scoped worktree summary

The scoped tracked diff is **not empty**, but none of it was produced by Stage 1. The worktree
already contained these unrelated tracked modifications before this work:

```text
 internal/audio/types/interfaces.go | 30 ++++++++++++++++++++++--
 internal/audio/types/utils.go      | 48 +++++++++++++++++++++++++++++++++++++-
 2 files changed, 75 insertions(+), 3 deletions(-)
```

Those two files were not edited during Stage 1. All Stage 1 implementation files and this report are
untracked new files. The broader worktree was already dirty with unrelated untracked review documents;
those remain untouched. No commit, push, branch change, model/asset change, native build, or Stage 2
artifact was created.

---

**END — STAGE 1 IMPLEMENTATION REPORT**
