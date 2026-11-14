# Stage 2 — Fake Helper / Real-Process Interoperability

**Date:** 2026-11-06  
**Scope:** Stage 2 only  
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` §17, including A1/A1-D7  
**Accepted Stage 1 baseline:** `docs/review/IMPLEMENTATION-howl-tts-helper-supervisor-stage1-2026-11-06.md` and `docs/review/CORRECTIVE-howl-tts-helper-supervisor-stage1-2026-11-06.md`  
**Status:** STAGE 2 ACCEPTED FOR ITS REAL-PROCESS INTEROPERABILITY MILESTONE, WITH THE S2-1/S2-2 CORRECTION APPLIED. STAGE 3 IS NOT AUTHORIZED. NO COMMIT OR PUSH.

---

## 1. Scope and authorization boundary

This pass implements only the authorized Stage 2 objective: a provider-neutral fake helper executable and real OS-process interoperability tests against the accepted Stage 1 supervisor.

The fake helper:

- is a separate command, `cmd/howl-tts-helper-fake`;
- is launched by the actual `Supervisor`/`StartProcess` path;
- receives parent-owned stdin/stdout/stderr pipes;
- speaks Protocol v1 NDJSON with the real Stage 1 client;
- creates deterministic synthetic PCM16 mono WAV output;
- is a test implementation, not a production TTS provider.

No native Kokoro/G2P implementation, model, voice, ONNX Runtime, CGO, Python, eSpeak, DLL, production provider registration, pipeline/character/cache/UI integration, download, packaging, or release change was made. The main application remains `CGO_ENABLED=0`.

The owner accepted Stage 2 for its real-process interoperability milestone and authorized a narrow correction for S2-1 and S2-2. That correction is recorded in `docs/review/CORRECTIVE-howl-tts-helper-capabilities-no-overwrite-2026-11-06.md`. Stage 2 is accepted for the milestone, not as a claim that every contract detail was conformant when first submitted. Stage 3 is not authorized. No commit or push occurred.

---

## 2. Exact changed files

All additions are local and uncommitted.

| File | Change |
|---|---|
| `cmd/howl-tts-helper-fake/main.go` | New fake-helper executable entry point |
| `internal/audio/tts/fakehelper/protocol.go` | Duplicated strict Protocol v1 DTOs, framing, duplicate-key and envelope validation |
| `internal/audio/tts/fakehelper/server.go` | Reader/worker dispatch, operations, state, cancellation, shutdown, test modes |
| `internal/audio/tts/fakehelper/wav.go` | Deterministic synthetic PCM16 mono WAV creation and validation |
| `internal/audio/tts/fakehelper/publish.go` | No-replace final publication (correction) |
| `internal/audio/tts/fakehelper/pathcheck_windows.go` | Fake-helper per-component reparse validation |
| `internal/audio/tts/fakehelper/pathcheck_other.go` | Non-Windows symlink validation |
| `internal/audio/tts/fakehelper/server_test.go` | In-process protocol, strictness, busy, WAV, capability, and path tests |
| `internal/audio/tts/helper/protocol.go` | Capability registry and mandatory service validation (correction) |
| `internal/audio/tts/helper/client.go` | Service/provider capability gates and `unsupported_capability` (correction) |
| `internal/audio/tts/helper/protocol_test.go` | Capability-registry and mandatory-capability tests (correction) |
| `internal/audio/tts/helper/client_test.go` | Client capability-gate tests (correction) |
| `internal/audio/tts/helper/fake_test.go` | Conformant five-service-capability fake hello (correction) |
| `internal/audio/tts/helper/testdata/golden/hello.response.json` | Corrected registry-native golden hello (correction) |
| `internal/audio/tts/helper/interop_test.go` | Real-process test harness, interoperability matrix, and strict handshake capability assertion |
| `internal/audio/tts/helper/interop_process_windows_test.go` | Windows PID existence/exit observation |
| `internal/audio/tts/helper/interop_process_other_test.go` | Non-Windows PID fallback |
| `docs/review/IMPLEMENTATION-howl-tts-helper-fake-process-stage2-2026-11-06.md` | This report, including the post-review correction record |
| `docs/review/CORRECTIVE-howl-tts-helper-capabilities-no-overwrite-2026-11-06.md` | S2-1/S2-2 correction report |

No accepted Stage 1 supervisor, ownership, discovery, output, path-check, or contract file was edited. The authorized S2-1 correction changes the accepted main-side `protocol.go` and `client.go` capability validation/gating; no process-lifecycle or output-validation file changed. The pre-existing broader worktree changes were not touched.

---

## 3. Baseline and resulting hashes

The Stage 2 baseline is the accepted Stage 1/corrective state. The hashes below are the starting hashes; `client.go` and `protocol.go` were subsequently changed only by the authorized S2-1 correction. Their corrected hashes are below and in the correction report.

| Baseline file | SHA-256 |
|---|---|
| `internal/audio/tts/helper/supervisor.go` | `3dd467d3d26323d47afbda1d2b8fec9018e9ab90115c91c05fa2fd92802718f3` |
| `internal/audio/tts/helper/client.go` | `e2b2e713f3e44d0e93392aa2508ef9ff512cc0a625848664f7a949dff1e72e77` |
| `internal/audio/tts/helper/protocol.go` | `8a9f710b4b141de342be88dafb4e046cb6c6b095507a70faa58db7fac7f48780` |
| `internal/audio/tts/helper/output.go` | `c50fce52ab901f6a16500522132cea58000b3d119b8dd8bb9a67186a608c10c0` |
| `internal/audio/tts/helper/discovery.go` | `af785db5a7b4bcc802cadfacc8ad60c57dc3f0210745651086f9cfd336369732` |
| `internal/audio/tts/helper/pathcheck_windows.go` | `a3747e78f211acd3a5ec1b492877511673f7b47c4410a032d5a1b68d135e66ab` |
| `internal/audio/tts/helper/pathcheck_other.go` | `e799d9c9520cdcb4de1fdd80c93dd2f38d5a6132bddd00d8d8a845738415b8fa` |
| `internal/proc/job_windows.go` | `1edef75a7fc753b949707cb1066056ef137cfc511afca9fea53f28ecee663c46` |
| `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` | `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393` |
| `docs/review/IMPLEMENTATION-howl-tts-helper-supervisor-stage1-2026-11-06.md` | `1e8e6553a6e5c2f71f391940fe6ef670fbcf73d4c9381b23b066eeb639705d0a` |
| `docs/review/CORRECTIVE-howl-tts-helper-supervisor-stage1-2026-11-06.md` | `5308ad5700b586ebc9c7764d654fa5bd770b128b52820ca8c9184a1b38a3e350` |
| `docs/review/AUDIT-howl-tts-helper-stage1-adversarial-2026-11-06.md` | `03ed0a4676089a0a3086bbbfeaf36f9159e6cca3bd9ec9b6d63121cfa68680f0` |

Resulting Stage 2 artifacts:

| Artifact | SHA-256 |
|---|---|
| `cmd/howl-tts-helper-fake/main.go` | `50b61853126cd4bc655cd7bfc76e1e032d73ba724c79caf16b58a0014d2404de` |
| `internal/audio/tts/fakehelper/protocol.go` | `ece9097423381653a999a126fe5e4bcbff084e4d80a63a9fef77890ea8dd79ea` |
| `internal/audio/tts/fakehelper/server.go` | `c4c613306569f1f34e516084ef0b90199e2fdece673433dbb3f62d9045f2767e` |
| `internal/audio/tts/fakehelper/publish.go` | `f6844a5c54c3a8489d1fdeaf0deaf5325db92dedcc67d0fbc0c2cd2c4c7c986d` |
| `internal/audio/tts/fakehelper/wav.go` | `87106725704a670335fc72d28eb9f0493cbe2a9dcdbdd161630939c92915f0e7` |
| `internal/audio/tts/fakehelper/pathcheck_windows.go` | `7bd034ca91cf7d24530b32c86b9d29d4ca519e4a73e7a01f286f8abaff6fd78e` |
| `internal/audio/tts/fakehelper/pathcheck_other.go` | `6ed0d5eea5df03da99b009c5aafe425925dba42f0159e9d451356a98de5c5b07` |
| `internal/audio/tts/fakehelper/server_test.go` | `1f2234bd2720231f0915a2d6f1b009ba3f093012241c0b6c242632fadf0ca846` |
| `internal/audio/tts/helper/protocol.go` | `3a0b1cd57655b6731641b2e7f6a5f0dad42dd45504fd503c2bc1c0bb1ad30204` |
| `internal/audio/tts/helper/client.go` | `b33c9f5151ba36f4b1363d19b0c8875a6211295077f40d2e413665866577160a` |
| `internal/audio/tts/helper/protocol_test.go` | `2c3b8e1e2470f83f23f92dda3849fd07fd7d55b93e93bfe875a913fff1dedb8f` |
| `internal/audio/tts/helper/client_test.go` | `ff347711e76b78d0b1c8a3c8b3aa378acb535356b8a7d7745211737d00ce0a6d` |
| `internal/audio/tts/helper/fake_test.go` | `c5d2a857fb2d70c23e3e7d6d41c47438eb96443f1da01003cf505e55a5b3a3a1` |
| `internal/audio/tts/helper/interop_test.go` | `28ed6d11b551fc2aa86c4bec02be4e7311002bb45f52e654578d75b3d8cf203d` |
| `internal/audio/tts/helper/testdata/golden/hello.response.json` | `e0d5ded1065a8e7c5a3a870d138d7c74bb05012237ac33f0b5b72951c9414709` |
| `internal/audio/tts/helper/interop_process_windows_test.go` | `0737e1bbad6c0e2dfffd6e794687b11fdb0fc75c01081f0ec16fa71bd72f4399` |
| `internal/audio/tts/helper/interop_process_other_test.go` | `1f2ea9ff045f48b0d4b6a6be2f513ca595253e8b6d8cac6af385aac64f642f18` |

This report's own SHA-256 is recorded in the handoff response.

---

## 4. Fake-helper architecture and protocol coverage

### 4.1 Separate executable, no shared helper import

`cmd/howl-tts-helper-fake` imports only `internal/audio/tts/fakehelper`. The fake package duplicates the small wire DTOs and does not import the accepted main-side helper package. This preserves the frozen ownership rule: duplicated hand-written DTOs, no shared module and no main-side helper import.

### 4.2 Actual process path

The interoperability tests set `HOWL_TTS_HELPER_PATH` to the built fake executable and use the accepted `Discoverer`, `Supervisor`, `StartProcess`, `proc.Start`, `proc.HideWindow`, `proc.KillTree`, `osProcess`, pipes, stderr ring, and single supervisor `Wait()` observer. The test factory only appends `HOWL_FAKE_HELPER_MODE` and `HOWL_FAKE_HELPER_MARKER_DIR` to `LaunchSpec.Env`; it still calls the real `StartProcess`.

The helper process is created suspended, assigned to the existing Windows job object, and resumed only after ownership succeeds. There is no shell invocation, TCP listener, HTTP transport, named pipe, or audio-byte IPC.

### 4.3 Reader/worker separation

The helper has one protocol reader and one operation worker:

- the reader parses and strictly validates frames, writes all responses, and handles `health`, `cancel`, and `shutdown` while a worker operation is active;
- the worker executes `catalog`, `load`, `unload`, and `synthesize`;
- a second non-control operation receives immediate `busy` and is not queued;
- `queue_depth` is always zero;
- one active synthesis is enforced by the worker gate;
- all stdout writes are serialized and flushed one complete JSON line at a time.

The normal path produced no stderr diagnostics. `TestRunHelloHealthAndShutdown` also asserts empty diagnostics for the in-process fake session.

### 4.4 Frozen operations

| Operation | Fake-helper behavior |
|---|---|
| `hello` | Exact v1.0, client/version/range validation, provider and service capabilities, extensions, frozen limits |
| `health` | Deterministic process/state/readiness/accepting-work snapshot; active synthesis ID only while synthesizing |
| `catalog` | Deterministic `fake-go` / `fake` catalog with `fake-voice`; optional unavailable voice; no load |
| `load` | Validates provider, identity, absolute clean asset root, contained regular assets; deterministic loaded state; idempotent same-identity load; different identity rejected |
| `unload` | Idempotent state release; no model |
| `synthesize` | Validates loaded identity, text, voice, language, speed, output path; writes synthetic WAV; cancellation checks; terminal publication barrier |
| `cancel` | Separate request ID; target acknowledgement; active synthesis cancellation; active-load `request_not_cancellable` |
| `shutdown` | Stops acceptance, cancels active work, acknowledges, and exits 0 |

The fake parses envelopes with `DisallowUnknownFields`, a duplicate-key scanner, UTF-8/BOM/NUL checks, the 4 MiB frame cap, exact v1.0 checks, request-ID validation, extension rejection, and per-operation strict parameter decoding. Unknown operations return `unsupported_operation`; malformed/duplicate-key frames are fatal with a best-effort frame error and exit 4.

The post-review capability correction is enforced on both sides of the boundary. The parent validates every advertised ID against the nine-ID v1 registry, rejects a wrong version, and requires the five service capabilities (`catalog`, `load`, `unload`, `synthesize`, `cancel`) during `hello`. Before an operation it checks the corresponding service capability, and for synthesis it checks the provider's `synthesize`, `speed`, and `language` capabilities when those features are requested. A missing capability returns `unsupported_capability` without sending the operation frame. The fake helper applies the same service/provider gate; `partial-capabilities` is a test-only mode for exercising the optional-feature rejection.

### 4.5 Deterministic test modes

These are test-harness modes, not new protocol operations or production behavior:

- `normal`;
- `hold-synth` — synthesis waits for cancel or a marker release;
- `hang-synth` — synthesis ignores cancel so the parent grace/kill path can be exercised;
- `crash-on-health` — process exits instead of answering health;
- `malformed-hello` — writes non-protocol stdout;
- `malformed-health` — writes non-protocol stdout during a healthy session;
- `incompatible` — returns `protocol_incompatible` and exits;
- `bad-wav` — returns a deliberately incorrect SHA-256 for parent validation;
- `partial-capabilities` — advertises synthesis without optional speed/language capabilities.

No request or synthesis text is written to stderr.

---

## 5. Real-process interoperability matrix

Environment: `go1.27.3 windows/amd64`, module `howl-chat`, actual child process creation through `internal/proc.Start`.

| Test | Observed evidence | Result |
|---|---|---|
| `TestStage2RealProcessHandshakeHealthCatalog` | Environment discovery; child start; hello; protocol `1.0`; helper `fake-1.0`; five service capabilities; health `unloaded`; catalog without load | PASS |
| `TestStage2RealProcessLoadSynthesizeUnload` | Load state `ready`; synthesis result; 2,444-byte canonical WAV; 24,000 Hz / mono / 16-bit; no `.part`; unload returns `unloaded`; shutdown | PASS |
| `TestStage2CancellationAndControlResponsiveness` | Held synthesis; health while active reports `synthesizing`, not accepting work, active ID; cancel acknowledged; cancelled synthesis has no final or `.part`; helper remains alive after cooperative cancel | PASS |
| `TestStage2ForcedCancellationKillsAndCleans` | Context cancellation with a helper that ignores cancel; grace expiry kills; PID gone; no output | PASS |
| `TestStage2TimeoutKillsAndCleans` | 100 ms operation timeout; helper killed; PID gone; no final or `.part` | PASS |
| `TestStage2ShutdownObservesExit` | `shutdown` acknowledged; helper exits; status no longer alive | PASS |
| `TestStage2UnexpectedTermination` | `crash-on-health`; health fails; supervisor reaches `crashed`/`exited`; PID gone | PASS |
| `TestStage2MalformedHelperResponseFailsSafely` | Non-protocol stdout during handshake; `invalid_frame`; helper terminated; state `unavailable`; PID gone | PASS |
| `TestStage2IncompatibleHelperFailsSafely` | Wire `protocol_incompatible`; state `incompatible`; helper terminated; PID gone | PASS |
| `TestStage2ParentRejectsBadWAVHash` | Helper returns a bad SHA-256; parent rejects with `output_failed`; rejected file removed; process terminated | PASS |
| `TestStage2MissingHelperIsNonfatal` | No helper path; `helper_missing`; process not alive; Chat path is not involved | PASS |
| `TestStage2ReaderDoesNotLeakAfterShutdown` | Shutdown, PID gone, bounded goroutine count back to baseline tolerance | PASS |
| `TestServerRejectsSecondNonControlOperation` | One active non-control operation; second rejected; work accepted after finish | PASS |
| `TestWindowsFinalSymlinkValidation` | Symlink creation denied by Windows privilege; final-file execution still unverified | SKIP |

Eleven of the twelve `TestStage2...` tests launch a real process; `TestStage2MissingHelperIsNonfatal` deliberately launches none. All twelve passed in the final helper-package run. None uses a fake transport or fake process. The only test-specific seam is the environment appended to the actual launch specification; actual `StartProcess` and job-owned process startup remain in use.

---

## 6. WAV generation and parent validation

The fake writes a deterministic, clearly synthetic 24,000 Hz PCM16 mono WAV with 1,200 samples (0.05 seconds, 2,444 bytes) and no real speech or external asset. It constructs the canonical 44-byte RIFF/WAVE/`fmt `/`data` header, writes a low-amplitude repeating ramp, hashes the bytes, creates `<final>.part-<16 hex>` with exclusive creation, validates the bytes, and publishes with `publishNoReplace`. That primitive creates the final directory entry with a hard link, which fails if the destination exists, then unlinks the part; it does not call `os.Rename`. The pre-publication `Lstat` check remains, but the no-overwrite guarantee no longer depends on it alone.

`TestStage2RealProcessLoadSynthesizeUnload` verifies:

- the returned absolute path equals the request-owned final path;
- metadata is 24,000 Hz, mono, 16-bit, positive duration;
- the file is 2,444 bytes with canonical RIFF/WAVE markers;
- the file hash is 64 lowercase hex characters;
- no `.part-*` remains.

`TestStage2ParentRejectsBadWAVHash` proves the parent does not trust the helper metadata: the fake publishes a valid WAV but reports a wrong SHA-256; the actual WAV/hash check rejects it as `output_failed` and cleans the exact requested file.

---

## 7. Cancellation, malformed protocol, crash, and shutdown evidence

- Cooperative cancellation: `TestStage2CancellationAndControlResponsiveness` holds synthesis, confirms health remains responsive, sends cancellation, receives the terminal error, and verifies no output remains.
- Forced cancellation: `TestStage2ForcedCancellationKillsAndCleans` uses a helper that ignores cancel, forces the parent grace expiry, and verifies the real PID is gone.
- Timeout: `TestStage2TimeoutKillsAndCleans` uses a short operation deadline and verifies the real PID is gone with no output.
- Malformed stdout: `TestStage2MalformedHelperResponseFailsSafely` observes `invalid_frame`, termination, `unavailable`, and PID absence.
- Incompatible helper: `TestStage2IncompatibleHelperFailsSafely` observes the typed wire error, `incompatible` state, and PID absence.
- Crash: `TestStage2UnexpectedTermination` observes the authoritative `Wait` result as `crashed`/`exited`, not a health poll.
- Shutdown: `TestStage2ShutdownObservesExit` observes the acknowledgement, exit, and dead PID.
- Bad terminal result: `TestStage2ParentRejectsBadWAVHash` exercises the terminal-then-validation rejection and exact output cleanup.

The forced-cancel and timeout tests assert the end state—dead PID, no final, no `.part`. They do not distinguish terminate → observe exit → delete from delete-then-kill by observation at the OS-child boundary. The ordering is structural in `cleanupRejected` and `Terminate`, and Stage 1 proved it against the in-memory process; it is inferred, not observed, for the real child. No claim of an observed real-process ordering is made.

No unsolicited stdout message or second terminal was emitted by the fake on the exercised paths. stdout was protocol-only; stderr was empty on the normal path.

---

## 8. Process, waiter, and reader leak checks

- Every real-process test records the actual child PID and waits for it to disappear after shutdown, cancellation, timeout, crash, malformed handshake, or incompatible handshake.
- `TestStage2ReaderDoesNotLeakAfterShutdown` compares a pre-start goroutine baseline to the post-shutdown count and allows only a small bounded tolerance for test-runtime goroutines.
- The single `Wait()` owner remains the accepted `Supervisor.observeExit`; `osProcess.Wait` is `sync.Once`-guarded, and failed starts create no observer. Stage 2 did not add a second waiter or an autonomous restart loop.
- No child process artifact was written into the repository: the helper is built into a temporary directory by `TestMain` and removed after the package test run.
- The test suite does not claim a precise waiter counter; the no-duplicate-waiter property is structural in the accepted supervisor and is supported by the observed single exit path and dead PIDs.

---

## 9. Build, test, race, vet, and Chat results

Go version: `go1.27.3 windows/amd64`; module `howl-chat`.

Fake-helper package and accepted helper package:

```text
go test ./internal/audio/tts/fakehelper ./internal/audio/tts/helper -count=1
ok  howl-chat/internal/audio/tts/fakehelper  1.378s
ok  howl-chat/internal/audio/tts/helper      4.322s
exit 0
```

Race:

```text
go test -race ./internal/audio/tts/fakehelper ./internal/audio/tts/helper -count=1
ok  howl-chat/internal/audio/tts/fakehelper  2.566s
ok  howl-chat/internal/audio/tts/helper      6.812s
exit 0
```

Vet:

```text
go vet ./internal/audio/tts/fakehelper ./internal/audio/tts/helper
exit 0
```

No-CGO build:

```text
CGO_ENABLED=0 go build ./...
exit 0
```

Audio tree:

```text
go test ./internal/audio/... -count=1
exit 0
```

All audio packages passed, including `internal/audio` at 3.443s, `internal/audio/tts/fakehelper`, and `internal/audio/tts/helper`.

Root Chat:

```text
go test . -count=1
ok  howl-chat  12.633s (root package in the correction full run)
```

Full repository, correction run:

```text
go test ./... -count=1
exit 0
```

The correction run passed every package; `internal/image/scanner` took 232.315s.

**Full-suite nondeterminism:** the earlier Stage 2 full run exited 1 in the unrelated pre-existing `internal/audio` test `TestConcurrentDifferentConfigurationsRemainIsolated` (`hits 5 misses 2`, want `0/2`). The independent audit reproduced the same failure. The test passed focused 1× and 3× and passed in the audio-tree run; no Stage 2 or correction file touched `internal/audio/pipeline_cache_test.go` or the cache implementation. The correction run passed, but the repository suite is **not deterministically green** and must not be recorded as reliably passing.

Skipped test:

```text
go test ./internal/audio/tts/helper -run '^TestWindowsFinalSymlinkValidation$' -v
--- SKIP: TestWindowsFinalSymlinkValidation
reason: A required privilege is not held by the client
```

This remains an execution limitation, not a pass.

---

## 10. Skipped tests and limitations

1. **Final-file Windows symlink execution remains unverified.** The helper's production and fake-helper path walks cover the final component by construction, and the intermediate-junction Stage 1 test passes, but this account cannot create the final symlink. Do not describe final symlink as executed.
2. **No real native helper.** Stage 2 establishes interoperability with the fake helper only. It does not establish Kokoro/G2P, ONNX Runtime, model, voice, native cancellation, or real TTS performance compatibility.
3. **Synthetic audio is not speech.** The WAV is a deterministic ramp for protocol/file testing only. No real speech, recording, or external audio asset was used.
4. **No restart budget, backoff, or idle close.** These remain deferred. The provisional 1/2/5-second schedule, 3/10-minute budget, and 120-second idle close are not measured or implemented guarantees.
5. **F-4 and F-5 remain carried.** `Shutdown`/`workMu` serialization and transient liveness observation were not changed.
6. **Duplicate-waiter evidence is structural, not counter-instrumented.** The accepted single observer and `waitOnce` remain; real PIDs are checked for survival. No new instrumentation was added to `osProcess`.
7. **Reader-leak evidence is bounded goroutine counting.** It is not a registry of every reader goroutine.
8. **No review index entry was added.** The existing index does not maintain a Stage 2 entry, and no unrelated review entry was altered.
9. **Kill-before-delete ordering is inferred, not observed, at the real-process boundary.** S2-3 is an evidence limitation. The end-state assertions are real; the ordering proof remains the accepted Stage 1 in-memory test plus the structural `cleanupRejected` gate.

---

## 11. Carried Stage 1 conditions

The Stage 1 acceptance conditions remain in force:

- final-file Windows symlink execution remains open;
- F-4 and F-5 remain open informational findings;
- restart budget/backoff/idle-close remain deferred and unmeasured;
- real native-helper compatibility remains unproven;
- no native provider, production pipeline, character, cache, UI, packaging, download, or licensing work is authorized by Stage 2.

---

## 12. Explicit non-scope confirmation

No native Kokoro or G2P implementation, ONNX Runtime, native DLL, Python, eSpeak, model/voice asset, production TTS provider, production audio-manager wiring, pipeline, character voice, cache, UI, system resource, download, license, packaging, or release change was made.

The fake helper is not registered as a production provider and is not a user-facing TTS option. It exists only as a command and test fixture.

---

## 13. Worktree and authorization status

- `cmd/howl-tts-helper-fake/` and `internal/audio/tts/fakehelper/` are untracked new directories.
- `internal/audio/tts/helper/` remains untracked from Stage 1 and now contains the Stage 2 interop tests plus the authorized S2-1 correction to `protocol.go` and `client.go`; no supervisor, ownership, discovery, output-validation, or path-check source was edited.
- The broader worktree was already dirty with unrelated pre-existing changes; none were modified by Stage 2.
- The temporary helper binary is built under the system temp directory and removed after tests. No repository binary artifact remains.
- **No commit, push, branch change, staging, or deployment occurred.**
- **Stage 2 is accepted for its real-process interoperability milestone. The S2-1/S2-2 correction is applied; Stage 2 is not reopened wholesale. Stage 3 is not authorized.**

---

**END — STAGE 2 IMPLEMENTATION REPORT**
