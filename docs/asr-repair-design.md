# T2.9 — ASR repair design and implementation sequence

**Status:** design complete and implementation-ready. This document is
design-only; it does not authorize implementation.

Inputs: the T2.8 audit, the accepted T2.9 lifecycle amendments, the `audio.*`
settings schema, and the connection-transparency requirements in
`docs/deferred-model-ecosystem.md` §3.

T2.2 and T0.5 remain open and unchanged. No stage acceptance criteria are
altered by this document.

---

## 0. Scope

### 0.1 In scope

- Make the existing local Whisper path reachable, correct, and honest.
- Make ASR and TTS coexist instead of displacing each other.
- Replace the pipeline transactionally and safely for concurrent readers.
- Carry recorded audio across the WebView boundary in a form the backend accepts.
- Insert the transcript where the composer actually is.
- Report capability state and failures truthfully.

### 0.2 Out of scope

- Downloading, curating, verifying, or licensing ASR models and binaries. That
  is the deferred model-ecosystem work.
- Streaming or partial transcription, VAD, diarization, translation.
- Remote ASR (Qwen). It requires the disclosure and consent flow in
  `docs/deferred-model-ecosystem.md` §3.3–3.5 before it may be offered at all.
- Any change to stage acceptance criteria or to T2.2's status.

### 0.3 Design invariants

1. Persisted settings express **desired** state; status expresses **actual**
   state. The two are never conflated in UI or API.
2. No capability is reported ready until its provider is initialized and usable.
3. No network access occurs as a side effect of enabling a local capability.
4. A failed change leaves previously working capabilities untouched.
5. Every resource acquired during a failed initialization is released or its
   leak is reported.
6. No panic escapes a manager or binding boundary.
7. One owner per resource and per temporary directory.

### 0.4 Resolved product decisions

These are settled and are design inputs, not open questions.

| Decision | Resolution | Consequence |
| --- | --- | --- |
| Implicit TTS on transcription | **Remove** | `audio_input` returns text only; speech requires an explicit `audio_output` request. T2.9.4 deletes the implicit branch and asserts zero synthesizer calls during transcription. Explicit synthesis is unchanged. |
| Push-to-talk | **Implement** | Pointer-down starts capture; pointer-up or pointer-cancel stops and transcribes; a release before the minimum duration yields a truthful "too short" message; blur or lost capture cannot leave recording active; a second start while recording is ignored. Implemented in T2.9.7. |
| Automatic model order | **Leave unconfigured** | `auto` returns a policy-not-configured error. No implicit selection or fallback. Because the persisted default is `asrEngine: "auto"`, enabling ASR does not initialize until the user selects an explicit engine; settings must direct them to do so, and startup reconciliation surfaces the same error. Explicit engines are unaffected. |

### 0.5 Boundaries

- This document does not authorize implementation.
- The sequence and gates in §5 remain in force.
- No unrelated stage status or project state is changed by this design.

---

## 1. Evidence found while designing

These findings come from source inspection. No runtime claim is made except
where explicitly labelled.

### 1.1 The fallback chain re-initializes recognizers with an empty config

`Pipeline.initASR` initializes each recognizer with a real `ModelPath`, then
calls:

```go
chain.Initialize(ctx, types.ASRConfig{Language: "auto"})
```

`FallbackChain.Initialize` loops over its recognizers and calls
`recognizer.Initialize(ctx, config)` again
(`internal/audio/asr/fallback/chain.go:50-58`). `whisper.Initialize` rejects an
empty `ModelPath` (`internal/audio/asr/whisper/recognizer.go:37-46`).

Consequence: even with a valid model, a valid `whisper-cli`, and a wired
toggle, initialization fails at `chain.Initialize` after the recognizers were
loaded. `p.asrRecognizer` is assigned only after that call
(`internal/audio/pipeline.go:236-242`), leaving loaded providers tracked by a
loader while the pipeline has no recognizer — and `EnableAudioASR` has already
shut down the previous pipeline.

This is source-derived, not observed at runtime. It makes the
double-initialization fix a prerequisite for any end-to-end test.

### 1.2 The validator rejects the format the browser records

`MediaRecorder` produces `audio/webm` (`frontend/dist/js/audio-service.js:77`).
`DetectAudioFormat` recognizes only RIFF/WAVE, MP3, OggS, and fLaC headers
(`internal/audio/types/utils.go:14-62`); the supported list is WAV, MP3, OGG,
FLAC, M4A (`:176-190`). A WebM temp file therefore fails validation even with
a real path.

The dead `ensureWAVFormat` (`recognizer.go:218`) would not help: it depends on
`ffmpeg`/`ffprobe` (`types/utils.go:66-140`), which are not bundled or verified
present, and `Transcribe` never calls it (`recognizer.go:107-109`).

Transport and format are one decision.

### 1.3 The recognizer temp directory is shared

`whisper.Initialize` uses `os.TempDir()/howl-whisper` (`recognizer.go:64`) and
`Release` removes the whole directory (`:161-170`). Multiple recognizers share
it, so releasing one can delete another's working files.

### 1.4 `Pipeline.Initialize` holds the write lock across provider I/O

`Initialize` holds `p.mu` while loading models and possibly acquiring binaries
(`pipeline.go:49-90`). Status blocks for the duration and a second enable
waits. Providers must be built outside the pipeline lock.

### 1.5 Binary discovery is relative and network-capable during enable

`findOrInstallWhisperCLI` searches relative paths and then downloads from
GitHub Releases (`recognizer.go:179-211`). `whisperInstallDir` is relative
(`:174-179`). Enabling a toggle must not download or depend on CWD.

### 1.6 Failed initialization leaks loaded providers

After a later error, `Pipeline.Initialize` returns without releasing what it
already loaded (`pipeline.go:56-79`). A discarded pipeline's loader is never
shut down.

### 1.7 Every reader of `a.audioPipeline` is unsynchronized

Readers: `TranscribeAudio`, `SynthesizeSpeech`, `SynthesizeSpeechWithVoice`,
`GetAudioStatus`, `GetTTSAvailableVoices`, `EnableAudioASR`, `EnableAudioTTS`
(`app-audio.go`), and `OnShutdown` (`app.go:291-293`). Nothing serializes them
against assignment (`app.go:264`, `app-audio.go:130,156`). A mutex around the
assignment would not prevent a swap from shutting down a recognizer in use.

### 1.8 Provider instances have no concurrency contract

`ProcessAudioMessage` snapshots its providers under `RLock` and calls them with
the lock released (`pipeline.go:86-96`). No synchronization primitive exists in
`internal/audio/asr/**` or `internal/audio/tts/**`; the only mutexes in the
package are `Pipeline.mu` and the TTS cache's.

- Two transcriptions can run through one recognizer concurrently. Whisper
  derives its output path from the input path (`input.wav` →
  `input.wav.json`, `recognizer.go:241-245`), so equal input paths collide.
- Voice enumeration and synthesis on one synthesizer are not serialized, and
  `ensureTTSVoice` mutates synthesizer state inside the message path
  (`pipeline.go:124-130`).
- A chain holds one set of recognizer instances; concurrent chain calls share
  them.

Leases protect lifetime, not provider safety.

### 1.9 Provider release ownership is ambiguous

`FallbackChain.Release` releases every recognizer (`chain.go:96-112`), while
the loader that initialized them also releases them on unload. Any call to
both paths double-releases. Whisper's `Release` is accidentally idempotent
because it clears `tempDir`, but that is not a contract.

Intended contract: `ModelLoader` is the sole release owner for providers it
constructs. `FallbackChain` borrows providers and must not initialize or
release them. Manager and pipeline may request disposal but must not call
provider `Release` directly. T2.9.3 establishes this; T2.9.1 only records it in
the deferred tests described in §5.

### 1.10 Builds cannot currently be cancelled

`ModelLoader.LoadASR/LoadTTS` initialize with `context.Background()`
(`internal/audio/runtime/loader.go`), and `Pipeline.Shutdown(ctx)` ignores its
context (`pipeline.go:186-206`). Until those change, cancellation is advisory
and a hung build or shutdown can only be timed out, not interrupted.

---

## 2. Decisions

### D1 — Transport: browser encodes 16 kHz mono WAV; Go owns a per-request temp directory

The recorder keeps WebM in memory, then the frontend decodes with
`AudioContext.decodeAudioData`, resamples through an `OfflineAudioContext` to
16 kHz mono, and encodes 16-bit PCM WAV. One request struct crosses the
binding:

```go
type TranscriptionRequest struct {
    WAVBase64     string `json:"wavBase64"`
    DurationMs    int    `json:"durationMs"`     // advisory; Go recomputes
    CorrelationID string `json:"correlationId"`
}

type TranscriptionResult struct {
    Text          string `json:"text"`
    DurationMs    int    `json:"durationMs"`     // authoritative, from PCM
    CorrelationID string `json:"correlationId"`
}

func (a *App) TranscribeAudio(req TranscriptionRequest) (TranscriptionResult, error)
```

The Go side base64-decodes into memory, validates, writes `input.wav` into a
fresh `os.MkdirTemp("", "howl-asr-*")`, transcribes, and removes the directory
on every exit path, including the provider's `input.wav.json` sidecar.

The legacy `TranscribeAudioFile(path string)` remains for future imported
files and rejects non-local inputs (`filepath.IsAbs` plus a scheme check), so
a `blob:` URL can never reach `os.Stat`.

**Imported-file contract (T2.9.5B).** The frozen local-file contract guarantees
local-path handling; it does not promise every container format. T2.9.5B
accepts only canonical PCM WAV imports in the same 16 kHz mono 16-bit form
produced for recordings, with no optional chunks beyond `fmt ` and `data`. It
validates the RIFF/WAVE structure, derives duration from the
data chunk, and applies `MaxImportedDurationMs = 600000` and
`MaxImportedWAVBytes = 19,200,044` before acquiring a manager lease. A
600-second file is accepted; a file over 600 seconds, over the byte cap, or in
MP3/OGG/FLAC/M4A is rejected with a specific error and no provider call.
Non-WAV imports remain unsupported until a bundled or otherwise verified
duration probe exists; they must never bypass the duration ceiling. The
15-minute transcription timeout applies to both recorded and imported
transcription. This narrowing is compatible with the frozen contract because
the contract guarantees local-path handling, not multi-format decoding.

**Why not a WebM path.** It fails `ValidateAudioFile` on format (§1.2).

**Why base64 rather than `[]byte`.** Encoding, size capping, and error
reporting stay in one place and do not depend on generated binding
marshalling. A `[]byte` signature may replace it once verified against the
packaged Wails runtime.

**Validation is authoritative in Go.** Browser checks are UX only and are
repeated server-side.

| Check | Browser | Go |
| --- | --- | --- |
| capture duration ≤ cap | timer stops recording | request `durationMs` sanity only |
| encoded Blob ≤ cap | before decode | — |
| decoded duration ≤ cap | after decode, before encode | — |
| WAV bytes ≤ cap | after encode | decoded length ≤ cap |
| base64 length ≤ cap | before invoke | string length ≤ cap |
| `RIFF` at 0, `WAVE` at 8 | — | required |
| RIFF size field equals `len(buf) - 8` | — | required; mismatch rejects even if the data chunk parses |
| chunk walk advances by `8 + size + (size & 1)` | — | required |
| exactly one `fmt ` and one `data` chunk | — | required |
| `fmt ` length ≥ 16, PCM format 1, mono, 16 kHz, 16-bit, block align 2 | — | required |
| `data` length even, `dataOffset + dataLen ≤ len(buf)` | — | required |
| no bytes after the final chunk except one pad byte | — | required |
| duration = `dataLen / 32000`, within cap and within 500 ms of `durationMs` | advisory | authoritative |

**Ordered pre-decode limits.** Constants are provisional until the packaged
runtime is measured:

1. `MaxRecordingMs` enforced by a capture timer (proposed 600,000).
2. `MaxEncodedBlobBytes` checked before `decodeAudioData` (proposed 64 MiB),
   because decode allocates before any duration check can run.
3. `MaxDecodedSeconds` checked immediately after decode.
4. `MaxWAVBytes` after encoding (proposed 19,200,044 for 600 s mono 16-bit at
   16 kHz plus header).
5. `MaxBase64Chars` before invoke (`ceil(4/3 × MaxWAVBytes)`).
6. The binding payload ceiling is measured in T2.9.5B. If it is below
   `MaxBase64Chars` with a 2× margin, reduce `MaxRecordingMs`; never raise the
   ceiling without measurement.

Each rejection is surfaced at its stage with a specific message. A
`decodeAudioData` failure is reported as a decode failure, never as "no
speech". Audio bytes are never logged; only lengths, duration, and correlation
ID are.

### D2 — Coexistence: one pipeline, per-capability initialization, two apply modes

Keep one `Pipeline`, but make initialization capability-independent: initialize
every requested capability, record per-capability readiness and error, and
never reinitialize with the other capability disabled.

A capability is **requested** when desired settings enable it. It is **ready**
when the active pipeline initialized it.

**Strict replacement (an active pipeline exists).** A candidate must
initialize every requested capability. If any fails:

- the candidate is not activated;
- its successfully initialized providers are released;
- the existing pipeline stays active and keeps serving;
- status keeps the old ready flags, sets `phase: "failed"`, and reports the
  candidate failure by capability and stage;
- persisted desired state is unchanged, so requested and ready visibly differ.

Example: old ASR + TTS ready, replacement requests both, TTS initializes and
ASR fails. The old pipeline remains, TTS keeps working, and the ASR error is
reported. Degraded replacement is not permitted while an older generation
serves.

**Degraded initial activation (no active pipeline).** When nothing is active —
startup or first enable — a candidate containing only the capabilities that
initialized is activated, because otherwise one broken capability blocks the
other forever. Status shows `phase: "degraded"`, the failed capability as
requested-but-not-ready with its error, and the ready capability as ready. If
no requested capability initializes, the manager stays empty and reports
`phase: "failed"`.

**Explicit disable is not a failure.** A capability absent from desired is
shut down with the retired generation.

**No implicit TTS.** `audio_input` returns text only. The `ttsSynth != nil`
branch inside the transcription case (`pipeline.go:104-113`) is removed;
speech is requested explicitly via `audio_output`. This is a resolved product
decision, not a pending option.

**Chain fix.** The chain must not initialize providers the loader already
initialized, and must not release them. It validates readiness and orchestrates
only. Providers receive the full `ASRConfig` built from settings at load time;
the chain never overwrites it. Provider-specific initialization errors are
preserved per candidate.

### D3 — Lifecycle: leased generations with tracked shutdown

**Slot.** Every initialized pipeline is wrapped in a slot that stays
manager-tracked until its shutdown has returned:

```go
type slot struct {
    p   ManagedPipeline
    gen uint64

    refs            int           // guarded by Manager.mu
    retired         bool          // guarded by Manager.mu
    shutdownStarted bool          // guarded by Manager.mu
    shutdownDone    chan struct{} // closed after Shutdown returns
    shutdownErr     error
}

type Manager struct {
    mu         sync.Mutex
    active     *slot
    retired    map[*slot]struct{} // includes slots whose Shutdown is running
    builds     int                // construction and disposal in progress
    closing    bool
    drained    chan struct{}      // one-shot; closed when the drain holds
    closeErrs  []error            // permanent for this manager
    closeTimeouts int
    terminal   *closeResult       // immutable snapshot once finalized
    desired    Settings
    desiredGen uint64
    appliedGen uint64
    status     Status
    lastBuildFailure *BuildFailure

    applyMu    sync.Mutex
    buildCtx   context.Context
    buildCancel context.CancelFunc
    ops        *OpRegistry
    factory    PipelineFactory
}
```

**Manager-owned pipeline and factory.** T2.9.2 does not store a concrete
`*audio.Pipeline`. It stores a manager-owned interface and receives candidates
through an injected factory:

```go
type ManagedPipeline interface {
    ProcessAudioMessage(context.Context, types.AudioMessage) (*types.AudioResponse, error)
    GetStatus() types.PipelineStatus
    GetTTSVoices() []types.Voice
    Shutdown(context.Context) error
}

type Settings struct {
    ASREnabled bool
    TTSEnabled bool
}

type CapabilityResult struct {
    Ready bool
    Err   error
}

type Candidate struct {
    Pipeline ManagedPipeline
    ASR      CapabilityResult
    TTS      CapabilityResult
}

type PipelineFactory func(ctx context.Context, desired Settings, candidate *Candidate) error

type Capability struct {
    Requested bool
    Ready     bool
    Error     string
}

type Status struct {
    ASR      Capability
    TTS      Capability
    InFlight int
    Builds   int
    Phase    string
}
```

The production `*audio.Pipeline` satisfies `ManagedPipeline`; no change to
`pipeline.go` is required. T2.9.2 tests supply fake implementations. The
minimal `Settings` and `Status` above are the T2.9.2 D3 types. The full
`Settings` and `Status` in §4 remain the T2.9.6 target; T2.9.6 extends these
types for D5. T2.9.2 does not implement D5 reconciliation, model path/backend
reporting, policy errors, or the UI-facing status contract.

**Factory contract.** The manager allocates a zero `Candidate` before each
build and passes it to the factory. The factory transfers ownership of every
resource it creates through that candidate:

- It sets `candidate.Pipeline` immediately after acquiring a pipeline handle,
  before any later operation that can fail or panic. If no handle exists,
  `Pipeline` stays nil.
- It sets `candidate.ASR` and `candidate.TTS` for requested capabilities only.
  `Ready` means that capability initialized successfully; `Err` is non-nil
  when it did not. Unrequested capabilities are not initialized and remain
  `{Ready:false, Err:nil}`.
- A nil return means the factory completed its attempt and the per-capability
  results are authoritative, even when one requested capability has an error
  (the M3 degraded and M3b zero-ready cases). A non-nil return means the attempt failed as a whole;
  the manager must not activate it, but must still dispose any non-nil
  `candidate.Pipeline`.
- The factory must not call `Shutdown` or otherwise release the candidate.
  Ownership transfers to the manager when the factory returns or panics.
- When at least one capability is requested, a pipeline with no ready requested
  capability is not activatable. The manager disposes it and reports the
  capability error(s), or the top-level factory error when present. When no
  capability is requested, a nil pipeline with zero results is the
  explicit-disable outcome, not a failure.
- When at least one capability is requested, a nil return with no pipeline and
  no ready requested capability is a failed attempt. When none is requested,
  the manager retires the old generation without activating a replacement.
- A candidate that reports any requested capability `Ready` with a nil
  `Pipeline` is invalid and is treated as a failed attempt.
- A panic after `candidate.Pipeline` is set leaves that handle available to the
  manager. The manager recovers, disposes the handle exactly once, and reports
  `BuildFailure`; a panic before any handle exists reports `BuildFailure` with
  no cleanup claim.

When the desired set requests at least one capability and an active generation
exists, every requested capability must be ready and the factory must return
nil; otherwise the candidate is disposed and the old pipeline remains active.
When no capability is requested, the manager retires any active generation and
does not build or activate a replacement; this is not a failure.
When at least one capability is requested and no active generation exists, only
ready requested capabilities are activated. If at least one requested
capability is ready and another requested capability failed, the phase is
`degraded`. If every requested capability is ready, activation is not degraded.
If no requested capability is ready, the candidate is rejected and disposed,
the manager remains empty, and the manager reports phase `failed`.

**States.** Candidate (not visible) → Active → Retired (serving in-flight
leases and/or shutting down) → Finished (removed after `Shutdown` returns).

**Acquire.** Under `mu`: reject if closing or `active == nil`; otherwise
`active.refs++` and return a lease bound to that slot. A lease never follows a
swap.

**Release.** Idempotent via `sync.Once`; decrements only its own slot. If the
count reaches zero and the slot is retired and shutdown has not started, mark
`shutdownStarted` and become the initiator. Outside `mu`, the initiator calls
`Shutdown`, records the outcome, then removes the slot from `retired` and
re-evaluates the drain. Another path that already started shutdown does
nothing further and never removes the slot.

**Every retirement path checks references.** A single helper, called under
`mu`, marks the slot retired, adds it to `retired`, and — if `refs == 0` —
claims shutdown. Apply, Close, and Release all use it, so a slot retired with
zero references starts shutdown immediately, and a slot retired with references
waits for the last Release.

**Shutdown outcome.** A slot is `NotStarted` (references remain), `Running`,
`Succeeded`, or `Failed`. Completion is not success: the slot leaves `retired`
in both completed states, but a failure appends `fmt.Errorf("generation %d
shutdown: %w", gen, err)` to `closeErrs` before removal.

**Apply and generations.** Apply stores the desired settings and increments
`desiredGen` under `mu`, then acquires `applyMu`. It re-reads the latest
desired and generation, and builds that. A candidate is swapped only if, under
`mu`, the manager is not closing and `desiredGen` still equals the candidate's
generation. Otherwise it is obsolete: it is disposed, then the builder loops
to build the latest desire. Queued callers therefore coalesce onto the newest
state rather than replaying intermediate toggles. A superseded caller returns
the latest applied status with `Superseded: true`.

**Candidate accounting invariant.** `builds` increments exactly once under
`mu` before the factory is called, and decrements exactly once only after the
candidate is installed as active (under `mu`), after its `Shutdown` returns or
panics, or after the manager determines that no handle exists. No path
decrements earlier: not on factory return, not on rejection, not on close, not
when disposal fails or panics. A build owns its `Candidate` and
`ManagedPipeline` from the factory call, so a partially initialized candidate
is disposed through the same `Shutdown` path.

**Close.** Each call performs one bounded wait; a timeout is that call's result,
not a terminal manager state.

1. Under `mu`: if already terminal, return the recorded result. If not already
   closing, set `closing`, move `active` into `retired` (so `active == nil`),
   and capture `drained`.
2. Cancel `buildCtx` and all registered operation contexts; interrupt
   subprocess providers that expose `Interrupt()`.
3. Wait on `drained` against this call's deadline, never holding `mu`.
4. Drain condition:
   `closing && active == nil && builds == 0 && len(retired) == 0 && ops.Len() == 0`.
   A slot leaves `retired` only after `Shutdown` returns, and `len(retired) ==
   0` implies no outstanding references.
5. On drain: finalize once under `mu` into an immutable `closeResult`
   capturing the outcome and `closeErrs`; return it. Success requires both the
   drain and `len(closeErrs) == 0`.
6. On timeout: increment `closeTimeouts`, return `ErrDrainTimeout` with
   outstanding generations; remain `closing`; do not finalize; do not close
   `drained`.
7. A later Close waits again. If the drain now holds without errors, it
   finalizes as `closed` and returns nil; with errors, `closed_with_errors` and
   the aggregate. The earlier call's timeout result is unchanged.
8. Repeated Close after finalization returns the recorded result without
   recomputation, slot traversal, or provider calls.

**Drain signal.** `drained` is created under `mu` whenever the condition is
false and a waiter may need it. Exactly one holder of `mu` closes it when the
condition becomes true and sets it nil, so it is closed once. Waiters capture
it under `mu` and wait outside. A timeout never closes it. Every state change
that can affect the condition re-evaluates it: close start, retirement,
shutdown completion, build increment/decrement, lease release, operation end.

**Operation registry and ordering.**

```go
op, ctx := ops.Begin(appCtx, opTimeout) // register cancellation first
defer op.End()                          // runs last

lease, err := mgr.Acquire(ctx)
if err != nil { return }
defer lease.Release()                   // runs before op.End()
```

A cancelled operation stays counted until it has stopped using its provider
and returned its lease. Drain requires both `refs == 0` and `ops.Len() == 0`.
Pure status reads take a manager snapshot and register nothing.

**Panic boundary.** Builds run synchronously on the caller's goroutine inside
`Apply`; no unguarded build goroutine is permitted. Recovered panics never
propagate.

```go
type BuildFailure struct {
    Gen        uint64
    Recovered  any
    CleanupErr error
}

func (e *BuildFailure) Error() string
func (e *BuildFailure) Unwrap() error // CleanupErr or nil
```

- Construction panic before a handle exists: recover, decrement `builds`,
  return `BuildFailure{Gen, Recovered, nil}`; no cleanup is claimed.
- Initialization panic with a handle: recover, attempt `Shutdown` exactly once
  under its own recovery guard, decrement only after that attempt returns or
  panics, and return `BuildFailure{Gen, Recovered, CleanupErr}`. Never swap.
- Factory error (not panic): treat the attempt as failed. Dispose any handle
  exactly once; return `errors.Join(factoryErr, cleanupErr)`; if cleanup fails,
  append it to `closeErrs`; decrement only after cleanup concludes. A nil
  factory error with a per-capability `Err` follows the strict/degraded rules
  above; if the candidate is rejected, dispose it and preserve the
  capability error(s) in status and the returned result.
- Disposal or retired shutdown panic: record `ShutdownPanic{Gen, Recovered}` in
  `closeErrs`, finish accounting, continue. No retry.
- A future goroutine implementation must recover and deliver the same typed
  result; a panic escaping manager-owned work is a defect.

**Error retention.** `closeErrs` belongs to the manager, not to a slot or
candidate. Every failed or panicked shutdown appends under `mu` before the slot
leaves `retired` or `builds` drops, so a concurrent or later Close cannot
observe an empty drain without also observing the error. A candidate that
vanishes after a panicked disposal still forces `closed_with_errors`.

**Provider concurrency (D8).** Refcounting protects lifetime, not provider
thread safety; see D8.

**Lock boundaries.** `mu` guards fields and counters only, never provider
calls, waits, initialization, or shutdown. `applyMu` is never acquired while
holding `mu`.

**Status.** Derived from active, desired, applied/desired generation, builds,
closing, and the last errors. `inFlight` sums references over active and
retired slots plus registered operations.

### D4 — Model readiness: resolve, preflight, report; no acquisition, no default policy

Resolution order, rooted at `paths.ProjectRoot()`
(`internal/paths/paths.go:29`) and never CWD:

1. `HOWL_ASR_MODEL` (explicit file) if set;
2. `<ProjectRoot>/models/whisper/ggml-<model>.bin` (canonical, nested);
3. `<ProjectRoot>/models/ggml-<model>.bin` (compatibility with the flat layout);
4. nothing else.

A relative override resolves against `ProjectRoot`; an absolute override is
used as given. Once resolution reaches the filesystem — for an explicit engine,
or for `auto` after a policy is present — a configured override that is missing
or unreadable fails closed with `ErrModelNotFound` and its absolute path. It
never falls through to the canonical layouts and is never handed to the policy
as if absent. `auto` with no policy reports `ErrPolicyNotConfigured` first and
performs no filesystem read.

**Model identity.** `ModelFile.Type` is never inferred from an arbitrary path.
The two modes differ:

- **Explicit engine.** The selected engine is authoritative. If
  `HOWL_ASR_MODEL` is set, it replaces the file for that engine and the
  returned `Type` is the explicit engine. The resolver does not inspect the
  file's contents. The caller contract is explicit: whoever sets the override
  is responsible for the file matching the selected engine. A mismatch is not
  detectable by path inspection and surfaces only when the provider fails or
  produces wrong output.
- **`auto`.** Typed candidates come only from the canonical layouts, where the
  `ggml-<name>.bin` filename supplies a known identity. The override
  participates as one additional *untyped* candidate (`Type` empty,
  `Source` env). A policy may assign a known engine type to that candidate; the
  assignment is the policy's assertion, not a resolver inference. It may not
  assign a different type to a canonical candidate, whose
  `ggml-<name>.bin` identity is fixed. Selecting the untyped override without
  an explicit known-type assignment is invalid. The resolver never guesses
  identity from a filename, header, or size.

`whisper-cli` resolves through an env override, `<ProjectRoot>/bin/whisper/`,
then PATH. **No download is attempted during apply, enable, status, or
transcription.** A missing dependency fails closed with its path.

The resolver accepts an explicit policy:

```go
type SelectionPolicy interface {
    Candidates(installed []ModelFile) ([]ModelFile, error)
}

type ResolverConfig struct {
    Root   string
    Getenv func(string) string
    Stat   func(path string) (os.FileInfo, error) // nil means os.Stat
}

func (r *Resolver) Resolve(engine string, policy SelectionPolicy) (ModelFile, error)
```

`Stat` is a single function field, not an interface or virtual filesystem.
Every path existence check goes through `cfg.stat`, which calls `cfg.Stat`
when non-nil and `os.Stat` otherwise. The resolver does not call `ReadDir`,
`Glob`, `EvalSymlinks`, or `os.Stat` anywhere else. Candidate discovery stats
the fixed set of known canonical filenames rather than enumerating directories,
so `Stat` is the resolver's only filesystem primitive. Production leaves
`Stat` nil; tests inject a recorder.

No ordering is built in. `auto` returns `ErrPolicyNotConfigured` until the user
chooses a policy, and status reports that. Because the persisted default is
`asrEngine: "auto"`, enabling ASR reports that error until an explicit engine
is selected; the settings UI must direct the user to choose one. There is no
fallback to base, small, tiny, or any other model.

For `auto`, policy presence is checked before any filesystem access. With a
policy present, a configured-but-missing override fails before discovery. The
candidate set is deterministic: canonical typed candidates ordered by source
precedence (nested, then flat) and then lexically, deduplicated by absolute
path, plus the untyped override when it exists. A policy may filter or reorder
that set, but not introduce a path or change a typed candidate's identity. An
empty candidate set is `ErrModelNotFound`; an empty policy result is
`ErrNoModelSelected`.

| Error | Condition |
| --- | --- |
| `ErrPolicyNotConfigured` | `auto` was requested and no policy was supplied. No filesystem read occurs. |
| `ErrNoModelSelected` | A policy ran and returned no candidates. |
| `ErrPolicyInvalidSelection` | A policy returned a path outside the candidate set, retyped a typed candidate, duplicated an entry with conflicting identity, or selected the untyped override without a known type. |

The three errors are distinct: absence of a policy, an empty decision, and an
invalid decision. Status reports the specific error rather than collapsing
them.

Explicit engine selections map to that exact model file and fail if absent.
The policy is not consulted. Remote engines (Qwen) are excluded while the
offline policy is in force and are not offered until the disclosure flow
exists.

**Network boundary.** `findOrInstallWhisperCLI` must not run during
initialization or apply; a resolver that only looks replaces it. Acquiring a
missing binary or model is a separate, user-initiated operation outside this
repair. When built, it must use the disclosure, consent, activity indicator,
and history required by `docs/deferred-model-ecosystem.md` §3.3–3.5, following
the consent wording pattern in `internal/image/deps/deps.go:51`. Local
enablement, settings application, status, and transcription perform zero
network I/O.

### D5 — Settings: one mapping, one reconciliation call

```go
func (a *App) ApplyAudioSettings(s audio.Settings) (ApplyResult, error)
func (a *App) GetAudioStatus() audio.Status
```

Fields that exist today (`settings-state.js:52-79`):

| Setting | Effect |
| --- | --- |
| `asrEnabled` / `ttsEnabled` | requested capabilities |
| `asrEngine` | explicit model selection; `auto` = unconfigured policy error |
| `asrLanguage` | `ASRConfig.Language` |
| `asrDevice` | `ASRConfig.Device` (accepted; provider may ignore) |
| `asrThreads`, `asrBeamSize`, `asrBestOf`, `asrTemperature` | passed through to `ASRConfig` |
| `asrTranslateToEnglish`, `asrNoTimestamps`, `asrFlashAttention` | not supported by the current provider; disabled in UI with a reason |
| `vad*` | out of scope; disabled with a reason |
| `pushToTalk` | implemented per §0.4; see D6 |
| `audioInputDevice` | out of scope; browser device selection is independent |

Reconciliation rules:

- Startup loads persisted settings and applies once asynchronously after the
  window is up. Status carries `starting` until the apply finishes.
- Saving settings calls `ApplyAudioSettings` once with the full desired set;
  toggles do not call separate enable bindings.
- Apply is idempotent: if the normalized desired set equals the active set and
  the active pipeline is healthy, return current status without rebuilding.
- Rapid changes coalesce by generation (D3); a superseded caller returns the
  latest status.
- The toggle reflects requested state immediately; readiness reflects actual
  state and the error that prevented it.
- With `auto` unconfigured, the ASR row shows "automatic model selection is not
  configured" and directs the user to select an engine.

### D6 — Truthful UI: insertion, notification, mic, push-to-talk

- Target the composer's real selector (`.input-field` / `#message-input`,
  `frontend/dist/index.html:1164`, `frontend/dist/js/chat-input.js:207`).
- `insertTranscript(text)` returns a boolean; the success notification appears
  only when it returns true. A missing composer is an error, logged and
  surfaced.
- After insertion, dispatch `input` and focus the composer. The send path reads
  the value at click time, so no extra contract is required; confirm during the
  manual test.
- The mic is disabled unless ASR is requested **and** ready, with the reason in
  its tooltip. Recording cannot start otherwise.
- **Push-to-talk (implemented):** pointer-down on the mic starts capture;
  pointer-up or `pointercancel` stops it and transcribes; a release before the
  minimum duration (proposed 300 ms) shows "recording too short" and does not
  transcribe; window blur or `visibilitychange` stops capture and discards it
  with a cancellation message; a second start while active is ignored; keyboard
  activation provides an equivalent toggle. Pointer capture is used so a
  release outside the button still ends recording.
- The success path remains: verified insertion, then notification.

### D7 — Errors, status, and diagnostics

- Failures carry capability, stage (`resolve`, `load`, `validate`,
  `transcribe`, `insert`), and a user-facing sentence.
- Status includes per-capability requested/ready/error, resolved model path,
  backend, in-flight count, build count, close state, last apply error, and
  `lastBuildFailure` (message only; the recovered value stays in logs).
- `phase`: `idle`, `starting`, `applying`, `degraded`, `failed`, `closing`,
  `closed`, `closed_with_errors`.
- No notification appears for an operation that did not complete.
- Logs carry a correlation ID shared between the JS call and the Go handler;
  audio bytes are never logged.

### D8 — Provider contract: serialization, temp ownership, release ownership

**Serialization is per underlying resource, not per interface value.**

- Each recognizer/synthesizer has a `callMu` covering `Transcribe`,
  `Synthesize`, `GetVoices`, and voice changes, including lazy loads.
  Configuration becomes immutable after `Initialize`; the initialized flag is
  guarded.
- `GetVoices` and `Initialize` use the same `callMu` as `Synthesize`; no
  state-mutex/callMu lock-order inversion is permitted.
- A fallback chain has its own `callMu` so attempts for one input do not
  interleave with another chain call.
- `SynthesizeStream` must not recursively acquire `callMu` through
  `Synthesize`. It may call the locked `Synthesize` and then read the returned
  file only when that file is uniquely owned by the call and cannot be
  overwritten, replaced, or removed by another call or by `Release` before the
  read. Otherwise it must use one internal unlocked synthesis path and hold
  `callMu` through the read. T2.9.3 tests the chosen property.
- Whisper adds a package-level limiter at default concurrency 1 because each
  call spawns a process that loads the model. This is a memory/resource limit,
  not a correctness lock. T2.9.3 implements the limiter object, the default,
  an explicit startup-level setter, and an isolated test constructor. The
  limit must be configured before recognizers begin using the limiter; a
  change while calls are active must be rejected or have explicitly defined
  behavior. Tests must use independent limiter instances and must not mutate
  the package default. Settings/UI wiring, if any, remains T2.9.6.
- ASR and TTS hold different locks and may run concurrently. Voice enumeration
  and synthesis on one synthesizer serialize.

**Two temporary-directory lifetimes, one owner each.**

- *Recognizer-lifetime directory* (T2.9.3): `os.MkdirTemp("",
  "howl-whisper-*")` created in `Initialize`, stored on the instance, used only
  for provider-internal conversion artifacts. `Release` removes only this
  directory and clears the field. Initialization failure after creation cleans
  it up by defer. The shared `os.TempDir()/howl-whisper` constant is removed.
- *Per-transcription directory* (T2.9.5B): created by the binding before writing
  `input.wav`, removed by the binding on every exit path, including the
  provider's `input.wav.json` sidecar. The recognizer writes that sidecar next
  to its input; ownership and cleanup stay with the binding.
  `cleanupTempFile` deletes only inside the recognizer's own lifetime
  directory.

**Release ownership.** The loader is the single owner. The chain neither
initializes nor releases providers. The loader must accept the caller's
context (`LoadASR(ctx, ...)`, `LoadTTS(ctx, ...)`) so builds can be cancelled;
until then cancellation is advisory.

**Lazy TTS voice replacement.** A replacement voice is an additional load, not
an in-place provider mutation. The loader tracks each loaded synthesizer under
a unique id so a replacement cannot overwrite an earlier release function.
Superseded synthesizers remain tracked for the pipeline lifetime and are
released by the loader at shutdown. This deliberately retains their memory and
file resources until shutdown; it is not an optimization. T2.9.4 tests that
repeated replacement leaves every provider tracked, releases each exactly once,
and never creates an untracked provider.

**Initialization failure cleanup.** `Pipeline.Initialize` releases every
provider it loaded on any failure path, using the loader's tracked set, and
preserves the original error with `errors.Join`. No partial loader survives a
discarded candidate.

---

## 3. Target architecture

```
recorder (MediaRecorder, WebM in memory)
  │  decodeAudioData → OfflineAudioContext(1ch, 16 kHz) → WAV encoder
  ▼
AudioService.transcribeBlob()
  │  TranscriptionRequest: base64 + duration + correlation id
  ▼
App.TranscribeAudio ──► Manager.Acquire ──► Pipeline.ProcessAudioMessage
  │                        │                     │
  │                        │                     ├─ ASR chain (pre-initialized, serialized)
  │                        │                     └─ no implicit TTS
  │                        ▼
  │                   leases / generations / drain
  ▼
per-request temp dir (input.wav) ──► whisper-cli ──► input.wav.json ──► text
  │                                                              │
  └──────────── removed on every exit ───────────────────────────┘
                                                                 ▼
                                                 insertTranscript() → composer
                                                      │
                                                      └─ success only if found
```

Responsibilities:

- **Recorder**: capture and encode; knows nothing about model paths.
- **AudioService**: decode, enforce limits, base64, invoke, verify insertion.
- **Bindings**: thin; register operation, acquire lease, translate errors.
- **Manager**: desired vs actual, leases, generations, builds, drain, errors.
- **Pipeline**: capability composition only; no I/O under its lock.
- **Resolver/preflight**: paths, policy, readiness; no acquisition.
- **Providers**: serialization, temp ownership, no download on initialize.

---

## 4. Proposed interfaces

```go
// internal/audio
type Settings struct {
    ASREnabled bool
    TTSEnabled bool
    Engine     string // explicit; "auto" is unconfigured
    Language   string
    Device     string
    Threads    int
    BeamSize   int
    BestOf     int
    Temperature float64
    PushToTalk bool
}

type Capability struct {
    Requested bool
    Ready     bool
    Error     string
    Detail    string
}

type Status struct {
    ASR        Capability
    TTS        Capability
    ModelPath  string
    Backend    string
    InFlight   int
    Builds     int
    Phase      string
    PolicyErr  string
}

type ApplyResult struct {
    Status     Status
    Superseded bool
}

type BuildFailure struct {
    Gen        uint64
    Recovered  any
    CleanupErr error
}

type ShutdownPanic struct {
    Gen       uint64
    Recovered any
}

func (a *App) TranscribeAudio(req TranscriptionRequest) (TranscriptionResult, error)
func (a *App) TranscribeAudioFile(path string) (string, error) // local files only
func (a *App) ApplyAudioSettings(s audio.Settings) (ApplyResult, error)
func (a *App) GetAudioStatus() audio.Status
```

```js
// frontend/dist/js/audio-service.js
async function transcribeBlob(blob, correlationId) // → { text, durationMs }
function insertTranscript(text)                    // → boolean
async function applyAudioSettings(settings)        // → ApplyResult
```

Generated Wails bindings are rebuilt, never hand-edited.

The `Settings`, `Capability`, `Status`, and `ApplyResult` types above are the
D5/T2.9.6 target. T2.9.2 declares only the minimal manager-local `Settings`,
`Capability`, and `Status` subset in `internal/audio/manager.go`; T2.9.6
extends them. `ApplyResult` and the full status fields are not implemented in
T2.9.2.

---

## 5. Implementation sequence

Each step is independently reviewable, leaves the app runnable, and has a
measurable exit gate. T2.9.2 uses fake pipelines only. T2.9.4 is the real
pipeline gate. No step claims a property owned by a later step.

| Step | Owns | Prereq | File scope | Tests | Exit gate |
| --- | --- | --- | --- | --- | --- |
| **T2.9.1** Characterization and resolver | Resolution precedence; model identity; policy injection and its distinct errors | none | `internal/audio/model_resolver.go` + test, `types/utils_test.go`, `fallback/chain_test.go` | V1–V4 | Resolver tests pass with injected roots, env, engines, and policies; `auto` returns `ErrPolicyNotConfigured`; policy errors are distinct; chain single-init and loader-only-release assertions are present and skipped with the exact `T2.9.3 (P5)` prerequisite; resolver is not wired, so no production call site changes; no behavior change; no network |
| **T2.9.2** Manager: leases, generations, builds, drain, panics | D3 in full; manager-owned `ManagedPipeline` and factory | 1 | `internal/audio/manager.go`, `internal/audio/manager_test.go` | M1–M14, all `-race` | Manager semantics and accounting pass with fake pipelines; manager-owned interface/factory; no App wiring; no claim about real initialization or cleanup |
| **T2.9.3** Provider contract | D8 provider serialization, temp ownership, and release ownership; concrete provider conformance tests | 2 | `asr/asr.go`, `asr/asr_test.go`, `asr/whisper/recognizer.go`, `asr/whisper/temp_test.go`, `asr/whisper/concurrency_test.go`, `asr/qwen/client.go`, `asr/qwen/client_test.go`, `asr/fallback/chain.go`, `asr/fallback/chain_test.go`, `tts/tts.go`, `tts/tts_test.go`, `tts/piper/local.go`, `tts/piper/local_test.go`, `tts/kokoro/local.go`, `tts/kokoro/local_test.go`, `tts/melotts/client.go`, `tts/melotts/client_test.go` | P1–P5 | Per-resource serialization at the concrete provider boundary; distinct temp dirs; loader-only release; no download on initialize; T2.9.1 ownership assertions strengthened, unskipped, and passing; no ctx-aware loader claim; no P6 claim |
| **T2.9.4** Real pipeline capability and transaction | D2; ordinary-error and panic cleanup gates; context-aware loader; P6; lazy TTS provider ownership | 2, 3 | `internal/audio/types/interfaces.go`, `internal/audio/runtime/loader.go`, `internal/audio/runtime/loader_test.go`, `internal/audio/pipeline.go`, `internal/audio/pipeline_test.go`, `internal/audio/asr/fallback/chain_test.go` (signature-helper only) | R1–R8, P6 | Real pipeline matrix passes; strict replacement preserves old readiness; degraded startup truthful; exactly-once release on error and panic; no implicit TTS; context-aware loader; P6 passes; lazy TTS replacement keeps every provider tracked and released exactly once |
| **T2.9.5A** Production manager foundation and ownership migration | App-owned manager; production pipeline factory; migration of all direct App pipeline access; temporary enable/disable/status bridge | 2, 4 | `app.go`, `app-audio.go`, `internal/audio/manager_factory.go` (new), `internal/audio/manager_factory_test.go` (new), `app_audio_ownership_test.go` (new) | A1–A8 | One App manager; no direct pipeline owner; all existing audio bindings migrated; factory candidate-before-init; TTS/voices/old ASR leases; shutdown calls existing `Manager.Close` and reports failure; temporary bridge only; no D5 |
| **T2.9.5B** Transport, limits, WAV validation, binding temp | D1; per-transcription temp; App request access through manager leases; imported-file policy | 5A | `app-audio.go`, `frontend/dist/js/audio-service.js`, `frontend/dist/js/wav-encode.js` (new), `frontend/dist/index.html`, `app_audio_transport_test.go` (new), `frontend/wailsjs/**` (generated) | T1–T6, L1–L7 | Request contract works; request access uses manager leases; ordered limits; full WAV validation; imported files capped at 600 s and limited to PCM WAV; temp dir removed on every path; 15-minute timeout only for transcription; payload ceiling measured or gate blocked |
| **T2.9.6** Settings, status, startup reconciliation | D5; App apply/status/enable/disable and startup reconciliation; replacement of the temporary 5A bridge | 2, 4, 5A, 5B | `internal/audio/manager.go` (D5 type extension only), `app-audio.go`, `app.go`, `audio-settings.js`, `settings-state.js`, `audio-settings.html` | S1–S5 | Toggles never claim readiness from request; `auto` unconfigured is surfaced truthfully; failed apply preserves old readiness; enable/disable use manager Apply; 5A bridge replaced; D5 types extended without changing T2.9.2 lifecycle semantics |
| **T2.9.7** Insertion, notification, mic, push-to-talk | D6 | 6 | `audio-service.js`, `index.html`, CSS | U1–U5 | Transcript inserted before success; missing composer errors; mic gated; push-to-talk contract holds |
| **T2.9.8** End-to-end matrix and documentation | Verification | 5B, 6, 7 | docs only | Manual matrix | Every row passes or is recorded blocked with cause |

T2.9.5A introduces no 15-minute request timeout. TTS, voice enumeration, status,
enable/disable, and the temporary ASR bridge use manager state, a lease, or a
zero-timeout operation; `Manager.Close` is their cancellation boundary. The
15-minute `audioTranscriptionTimeout` is added by T2.9.5B and applies only to
`TranscribeAudio` and `TranscribeAudioFile`. TTS must not inherit it.

The 5A enable/disable/status bridge is temporary. T2.9.5B must not expand it
into D5 settings, persistence, full status, `ApplyResult`, or startup
reconciliation. T2.9.6 replaces it.

T2.9.5A preserves `Manager.Close(ctx context.Context) error` unchanged. With no
caller deadline it supplies the existing five-second default; a blocked drain
returns `ErrDrainTimeout`, leaves the manager `closing`, and does not force a
lease release. The App reports the error rather than recording shutdown as
successful. A7 is the execution test for this behavior; it is not changed to
make shutdown appear successful.

T2.9.3 removes implicit Whisper acquisition. A missing whisper-cli is a
model-not-found error with installation instructions; initialization never
downloads or builds it. Between T2.9.3 and T2.9.4, provider initialization
still uses `context.Background()` and cancellation remains advisory; T2.9.3
does not claim context cancellation or partial-initialization cleanup.

T2.9.4 retains superseded TTS synthesizers under unique loader ids until
pipeline shutdown. The retained memory and file resources are a deliberate
cost, not an optimization. Shutdown releases each tracked provider exactly
once.

T2.9.3 names every currently registered TTS adapter: Piper, Kokoro, and
MeloTTS. If inspection finds another registered adapter, amend the file scope
before implementation rather than silently omitting or expanding it.

Qwen serialization is provider-contract hardening of the existing registered
client only. It does not authorize offering remote ASR, wiring selection, or
bypassing the disclosure and consent boundary in §0.2.

T2.9.2 intentionally has no App file changes. T2.9.5A owns App manager
construction, removal of the direct pipeline owner, migration of every
existing App audio access, and the temporary bridge. T2.9.5B owns D1 request
access and the imported-file policy. T2.9.6 owns D5 settings/apply/status,
startup reconciliation, final enable/disable, and replacement of the 5A
bridge. The sequence is strict: 5A → 5B → 6. T2.9.5B must not begin until 5A
closes, and T2.9.6 must not begin until 5B closes.

Dependency graph:

```
1 → 2 → 3 → 4 → 5A → 5B → 6 → 7
5B, 6, 7 → 8
```

**Deferred ownership test placement.** The two ownership assertions live in
`internal/audio/asr/fallback/chain_test.go` as external test package
`fallback_test`. It imports `fallback`, `runtime`, and `types`. The production
graph is `fallback → asr, hw, types` and `runtime → asr, types, tts`; neither
imports the other, so the external arrangement has no cycle and also cannot
reach unexported hooks. The fake recognizer is registered through the existing
exported `ModelLoader.RegisterASR` and loaded through `LoadASR`; no production
test hook is added. The test pre-initializes the fake through the loader, calls
`chain.Initialize`, and asserts the init count remains one; it then calls
`chain.Release` and asserts the release count remains zero; finally it unloads
through the loader and asserts exactly one release. Both assertions are skipped
in T2.9.1 with the exact prerequisite `T2.9.3 (P5)` and must be unskipped at
that step's exit gate.

T2.9.3 extends the fixture with repeated `chain.Release`, repeated
`loader.Unload`, reverse-order release in both orders, and `chain.Transcribe`
assertions. T2.9.4 changes only the fixture's `LoadASR` call to pass context;
the assertions remain unchanged.

No real transcription attempt before T2.9.2 and T2.9.4 pass, and no recorded
clip accepted as evidence before T2.9.5B passes.

---

## 6. Test plan

### 6.1 Manager (fake pipelines)

| ID | Test | Pass condition |
| --- | --- | --- |
| M1 | Three generations, out-of-order releases | Each slot shut down once, after its own refs reach zero; new acquires never obtain a retired slot |
| M2 | Strict replacement: old ASR+TTS ready, candidate TTS ok, ASR fails | Old pipeline stays active and ready; factory returns nil top-level error with an ASR capability error; candidate TTS released; error names ASR; requested/ready diverge |
| M3 | Degraded startup: empty, both requested, ASR fails, TTS succeeds | Factory returns nil top-level error with ASR not ready and TTS ready; TTS active and ready; ASR error; phase degraded |
| M3b | Zero-ready startup: empty, both requested, both fail, candidate handle present | Factory returns nil top-level error with both capabilities not ready and per-capability errors; no active slot is installed and the manager remains empty; candidate `Shutdown` is called exactly once; status phase is `failed` with both capabilities requested/not-ready and errors preserved; `Acquire` fails before and after cleanup; `builds` remains nonzero while the fake `Shutdown` is blocked and returns to zero after cleanup; the rejected candidate is never observable as active |
| M4 | Explicit disable with an active manager | Disabled capability retired and shut down after refs drain; no failure reported |
| M4b | Explicit disable from an empty manager | Manager is empty and neither ASR nor TTS is requested; applying all-disabled settings leaves the manager empty, performs no replacement factory build, reports no failure, records both capabilities as not requested/not ready with no error, and leaves build accounting unchanged |
| M5 | Zero-reference retirement via Apply | Old slot starts shutdown immediately; stays in retired until Shutdown returns; new active serves |
| M6 | Zero-reference retirement via Close | Shutdown starts from Close; success only after it returns |
| M7 | Close overlaps final release with blocked Shutdown | No success until Shutdown returns and slot leaves retired |
| M8 | Timeout, then later drain | First call returns `ErrDrainTimeout` unchanged; second returns nil; state closed; no re-shutdown |
| M9 | Timeout, then drain with shutdown error | Second call returns aggregate; state closed_with_errors |
| M10 | Drain signal | Closed exactly once; timeout does not close it; concurrent Close calls do not double-close |
| M11 | Coalescing and stale builds | Final active generation equals latest desire; obsolete candidate disposed, never swapped |
| M12 | Candidate accounting | One increment before the factory; one decrement after activation or cleanup; never negative; outstanding count zero after every outcome |
| M13 | Panic boundary | Construction, initialization (including after a candidate handle is set), and shutdown panics produce `BuildFailure`/`ShutdownPanic`; manager cleans the disposed handle exactly once; no panic escapes; `closeErrs` retains failures after slots/builds disappear |
| M14 | Repeated Close | Terminal result is identical; no recomputation, traversal, or provider calls |

### 6.2 Provider contract (real temp dirs, fake and concrete providers)

| ID | Test | Pass condition |
| --- | --- | --- |
| P1 | Two recognizers initialize and release | Distinct directories; releasing A leaves B's files intact; no shared `howl-whisper` directory remains |
| P2 | Concurrent transcriptions on one recognizer | Maximum concurrent provider calls = 1 at the concrete provider boundary; base fake plus Whisper or Qwen; Whisper package limiter across two instances separately tested |
| P3 | ASR and TTS concurrently | Both run; no shared-state race under `-race`; independent resources do not serialize each other |
| P4 | Voice enumeration during synthesis | Serialized on one concrete synthesizer; `GetVoices` does not race with synthesis or `Initialize`; at least one HTTP and one process-backed adapter tested without external network |
| P5 | Chain does not initialize or release | Both T2.9.1 tests unskipped; `chain.Initialize`, repeated `chain.Release`, and `chain.Transcribe` do not change init/release counts; loader unload releases exactly once; repeated unload and reverse-order release in both orders do not release again |

Concrete adapter tests use injected command runners or `http.RoundTripper`
instances. They do not require real model execution, external network access,
or a system-installed Whisper/Piper/Kokoro binary.

### 6.3 Real pipeline (T2.9.4 gate)

| ID | Test | Pass condition |
| --- | --- | --- |
| R1 | ASR only / TTS only / both | Requested capabilities ready; independent failures recorded |
| R2 | Strict replacement preservation | Old pipeline remains ready after a failed candidate |
| R3 | Degraded initial activation | Ready capability active; failed one reported |
| R4 | Ordinary initialization failure with partial load | Exactly-once cleanup; `errors.Join` preserves both errors; builds drop after cleanup |
| P6 | Initialization failure after one provider | The provider whose initialization fails is released exactly once, as is every earlier successfully initialized provider; the original initialization error remains discoverable after cleanup errors are joined |
| R5 | Initialization panic with partial load | Same cleanup; `BuildFailure` carries panic and cleanup outcome |
| R6 | No implicit TTS | Successful transcription with TTS ready produces zero synthesizer calls |
| R7 | Chain behavior | Single initialization; no chain release; provider config preserved |
| R8 | Lazy TTS voice replacement ownership | A replacement voice is loaded under a unique loader id; the superseded synthesizer remains tracked for the pipeline lifetime; repeated replacement leaves every provider tracked, releases each exactly once at shutdown, and never creates an untracked provider |

### 6.4 Resolver, transport, temp, WAV

| ID | Test | Pass condition |
| --- | --- | --- |
| V1 | Precedence and override identity: env absolute/relative, nested, flat, missing | Explicit engine wins; env replaces the file but `Type` is the explicit engine; relative env resolves under the injected root, not CWD; missing env fails closed with `ErrModelNotFound` and its absolute path |
| V2 | `auto` policy states and stat observability | Injected `Stat` recorder: nil policy → `ErrPolicyNotConfigured`, zero `ModelFile`, and `calls == 0`; empty result → `ErrNoModelSelected`; foreign path, retyped canonical candidate, duplicated conflicting identity, or unlabeled env override → `ErrPolicyInvalidSelection`; valid selection returns the declared `Type` and `calls > 0` |
| V3 | Dangling symlink and relative CWD | Absolute path reported; no write outside the injected root; no download or CLI lookup |
| V4 | WebM vs WAV validation | WebM rejected; valid WAV accepted |
| T1 | Request round trip | Text, duration, and correlation ID returned |
| T2 | Limit ordering | Oversize blob rejected before decode; decoded, WAV, and base64 caps each reject at their stage |
| T3 | Malformed WAV table | Wrong format, channels, rate, depth, chunks, RIFF size mismatch, trailing bytes, duration mismatch all rejected specifically |
| T4 | Temp lifecycle | Binding directory absent after success, error, cancellation, and panic-safe defer |
| T5 | Legacy path guard | `blob:` and relative non-file inputs rejected |
| T5b | Imported-file limit | Canonical 600 s PCM WAV accepted at the guard; >600 s, >19,200,044 bytes, non-canonical WAV, and non-WAV inputs rejected specifically before any provider call; no fallback to provider max duration |
| T6 | Payload ceiling | Measured on the packaged runtime with documented margin, or gate blocked |

**V2 procedure.** Use `calls := 0` and a recorder
`cfg.Stat = func(string) (os.FileInfo, error) { calls++; t.Errorf("stat during nil policy"); return nil, fs.ErrNotExist }`.
Call `Resolve("auto", nil)`; assert `errors.Is(err, ErrPolicyNotConfigured)`,
a zero-value `ModelFile`, and `calls == 0`. The recorder itself fails the test
if the resolver calls it, so a seam call is both counted and reported. It
cannot observe a direct `os.Stat` bypass; that bypass is excluded structurally
by the no-direct-call rule above, and V2 does not claim otherwise. For the
positive control, inject a second
recorder that increments `calls` and delegates to `os.Stat` over a temp root
containing one canonical model; with a policy selecting that model, assert
success and `calls > 0`, proving the seam is wired rather than bypassed. A
nonexistent `Root` is never accepted as evidence of zero calls. Because `Stat`
is the only filesystem primitive, `calls == 0` is the observable test of the
nil-policy no-read claim; if implementation later adds `ReadDir` or another
primitive, V2 must be narrowed to the observed calls until that primitive is
injectable too.

### 6.4a App ownership, leases, and shutdown (5A/5B)

| ID | Test | Pass condition |
| --- | --- | --- |
| A1 | One manager and no direct pipeline owner | App has one `audioManager`; no `audioPipeline` field; no production path holds `*Pipeline` |
| A2 | Production factory order | Candidate is assigned before `Initialize`; ASR/TTS outcomes map without manager semantic changes |
| A3 | Temporary enable bridge | `Apply` is used; failed apply preserves prior readiness; no D5 fields or persistence |
| A4 | Status and voices | Status uses manager state; voices use a lease and release exactly once |
| A5 | TTS and old ASR leases | TTS and the temporary ASR bridge use operation + lease; their contexts have no artificial deadline; no 15-minute timeout |
| A6 | Close wiring | `OnShutdown` calls `Manager.Close(ctx)`; no direct pipeline shutdown; the returned error is reported, never treated as silent success |
| A7 | Five-second close and timeout | With no caller deadline, a blocked drain returns `ErrDrainTimeout` after the manager's five-second default; manager remains `closing`; no forced release; a later close finalizes after drain |
| A8 | Scope guard | No D1 transport, WAV encoder, imported-file path, limits, or T2.9.7 behavior in 5A |
| L1 | Recording lease | One operation and one lease; both released exactly once |
| L2 | Replacement | Replacement preserves a leased old generation until final release |
| L3 | Shutdown drain | Close waits for leases; timeout leaves `closing`; later close finalizes |
| L4 | Acquisition failure | Temp directory is removed and operation ends; no pipeline call |
| L5 | Request panic | Panic becomes an error after temp, lease, and operation cleanup |
| L6 | TTS/voice regression | 5B re-runs the 5A TTS and voice lease tests; this is regression verification, not new implementation |
| L7 | Timeout policy | `TranscribeAudio` and `TranscribeAudioFile` have the 15-minute deadline; TTS and voice enumeration have no such deadline |

L6 exists only to detect a 5B regression in the 5A migration. A failure is a
5A regression, not authorization to duplicate or redesign TTS/voice ownership
inside 5B.

### 6.5 Settings and status

| ID | Test | Pass condition |
| --- | --- | --- |
| S1 | Requested vs actual divergence | Enabling with a missing model shows requested true, ready false, and the specific error |
| S2 | Failed apply preserves prior readiness | Old ASR/TTS remain ready; candidate error is surfaced |
| S3 | Idempotent apply | An unchanged desired set does not rebuild or bump the applied generation |
| S4 | Startup reconciliation | A persisted explicit engine reaches ready without user action; persisted `auto` surfaces the policy error |
| S5 | Superseded apply | Rapid changes return the latest status; no stale readiness is reported |

### 6.6 UI

| ID | Test | Pass condition |
| --- | --- | --- |
| U1 | Present/absent composer | Boolean result; success only after insertion |
| U2 | Failure notification | No success toast when insertion fails |
| U3 | Mic gating | Disabled with reason when ASR not ready; cannot start |
| U4 | Push-to-talk | Pointer-down starts; pointer-up/cancel stops; blur stops; short release reports "too short"; second start ignored |
| U5 | `auto` unconfigured | Settings shows the policy error and directs engine selection |

### 6.7 JS harness

The repository has no JS test runner and no root `package.json`. Do not invent
one. Verify the encoder with a small harness page under the existing build
harness, or a temporary Node script importing a DOM-free encoder module. Record
which was used.

### 6.8 Manual matrix (requires the user, a model, and later headphones)

| # | Case | Expected |
| --- | --- | --- |
| 1 | Enable ASR with an explicit installed model, offline | ready; model path shown; no network |
| 2 | Enable ASR with `auto` | policy-not-configured error; no model chosen |
| 3 | Enable ASR without its model | path-specific error; TTS unaffected; no network |
| 4 | Enable ASR, then TTS | both ready; ASR still ready |
| 5 | Rapid on/off/on | final state matches last request; no leaked slots |
| 6 | Record 5 s, transcribe with TTS ready | text in composer; no speech; success only after insertion |
| 7 | Record beyond cap | clear rejection before handoff |
| 8 | Transcribe, then disable ASR mid-run | completes or cancels cleanly; no crash |
| 9 | Quit during transcription | process exits; no orphan `whisper-cli` |
| 10 | Hold-to-record; release early | short-release message, no transcription; blur cannot leave recording active |
| 11 | Restart with ASR enabled and an explicit engine | startup reconciliation reaches ready |

Microphone and device checks are deferred until the user has the hardware.

---

## 7. Risks and unresolved engineering questions

The imported-file duration/format policy is decided in D1. The following are
engineering unknowns.

1. **Packaged runtime limits.** Decode support for the recorded codec, the
   base64 payload ceiling, and the safe recording cap are unmeasured. T2.9.5B
   cannot pass its gate without them.
2. **Whisper CLI flags.** `--output-json`, `--suppress-non-speech-tokens`, and
   language behavior must be confirmed against the installed build's `--help`.
3. **In-process providers.** Providers that cannot be preempted rely on the
   abandon path; classification is required before the drain contract is
   considered complete.
4. **Build cancellation.** The loader uses `context.Background()` until
   T2.9.3 changes it; cancellation is a request, not a guarantee.
5. **Shutdown interruption.** `Pipeline.Shutdown` ignores its context; a hung
   release can only be timed out, not interrupted.
6. **No retry after failed cleanup.** Resources may leak; the design reports
   this rather than reclaiming it.
7. **Chain scope.** Whether a pipeline holds multiple recognizers (and thus
   needs chain-level serialization in practice) remains open.
8. **`BaseRecognizer` immutability.** Making configuration immutable after
   initialization may change the shared base type; scope unmeasured.
9. **`BuildFailure` exposure.** Whether the recovered value reaches the UI or
   only logs is a presentation choice for T2.9.6.
10. **Imported-file format narrowing.** 5B accepts only PCM WAV imports because
    no bundled duration probe is verified. Multi-format import needs a separate
    probe and is outside 5B; non-WAV files must fail closed rather than bypass
    the duration ceiling.
11. **Shutdown timeout visibility.** `Manager.Close` keeps its five-second
    default when the caller supplies no deadline. A blocked drain returns
    `ErrDrainTimeout`, leaves the manager `closing`, and never forcibly releases
a lease. The App must report that error. This is preserved and source-verified,
    not changed to make shutdown appear successful. Execution testing remains
    assigned to T2.9.5A (A7).

---

## 8. Acceptance criteria

With an explicit model installed and the machine offline:

1. ASR enables and reports ready; `auto` reports unconfigured until an engine
   is chosen.
2. Enabling ASR leaves TTS ready, and vice versa.
3. A failed enable leaves the previous capabilities working and reports why.
4. A recorded clip produces text that lands in the composer, and the success
   message appears only after it does.
5. Transcription produces no speech, even with TTS ready.
6. Hold-to-record, short-release, blur, and second-start behavior match §0.4.
7. No `blob:` URL reaches the backend; no temp audio or `.json` sidecar
   remains after success or failure.
8. No network request occurs during enable, apply, status, or transcription.
9. Manager tests pass under `-race`, including zero-reference retirement,
   timeout and eventual completion, one-shot drain signaling, panic boundaries,
   and retained disposal errors; real-pipeline gates pass separately.
10. The manual matrix passes or is recorded blocked with cause.
11. Imported transcription accepts only PCM WAV at or below 600 seconds and
    19,200,044 bytes; longer, larger, and non-WAV files fail before a provider
    call.
12. A shutdown without a caller deadline that cannot drain within five seconds
    reports `ErrDrainTimeout`; the manager remains `closing`, no lease is
    forcibly released, and the App does not report shutdown as successful.

Passing these does not close T2.2, change a stage, or declare ASR generally
functional.

---

## 9. Expected files

| File | Change |
| --- | --- |
| `internal/audio/manager.go` (T2.9.2, new; T2.9.6 type extension) | manager-owned `ManagedPipeline`/factory, slots, leases, generations, builds, drain, panic handling, minimal D3 `Settings`/`Status`; T2.9.6 extends the D5 fields only |
| `internal/audio/model_resolver.go` (new) | path precedence, model identity, policy injection and errors, preflight |
| `internal/audio/pipeline.go` | per-capability init, cleanup, no I/O under lock, no implicit TTS |
| `internal/audio/asr/fallback/chain.go` | no provider init/release; orchestration only |
| `internal/audio/asr/whisper/recognizer.go` | resolved paths, per-instance temp, no download on initialize |
| `internal/audio/asr/asr.go` | guarded initialization state, provider serialization |
| `internal/audio/runtime/loader.go` | caller context, single release ownership |
| `internal/audio/tts/**` | per-instance serialization, voice-change locking |
| `internal/audio/manager_factory.go` (T2.9.5A, new) | production manager/factory; candidate assigned before initialization |
| `app-audio.go` | T2.9.5A manager access and temporary bridge; T2.9.5B request/result bindings and imported-file policy; T2.9.6 apply/status |
| `app.go` | T2.9.5A manager construction and shutdown close; T2.9.6 startup reconciliation |
| `app_audio_ownership_test.go` (T2.9.5A, new), `app_audio_transport_test.go` (T2.9.5B, new) | ownership, close timeout, transport, lease, and timeout tests |
| `internal/audio/manager_factory_test.go` (T2.9.5A, new) | factory ordering and outcome mapping |
| `frontend/dist/js/audio-service.js` | WAV encode, base64, insertion, mic, push-to-talk |
| `frontend/dist/js/wav-encode.js` (new) | PCM WAV encoder |
| `frontend/dist/js/audio-settings.js` | apply settings, requested vs actual, policy error |
| `frontend/dist/js/settings-state.js` | defaults/migration for new fields |
| `frontend/dist/audio-settings.html`, `index.html` | control states, mic state |
| `frontend/wailsjs/**` | regenerated, not hand-edited |

---

## 10. What this design does not claim

- That the installed Whisper build accepts the WAV produced this way; that is
  a manual-matrix item.
- That `decodeAudioData` and the payload ceiling work in the packaged WebView;
  those are T2.9.5B gate items.
- That non-WAV imported files are supported in 5B. They are rejected until a
  bundled or verified duration probe exists.
- That model acquisition, integrity, or licensing is solved; that remains in
  `docs/deferred-model-ecosystem.md`.
- That remote ASR is available; it must not be offered until the disclosure and
  consent requirements are implemented.
- That any of this is implemented. This document is design-only.
