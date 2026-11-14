# INVESTIGATION — Go-native Kokoro/G2P hosting: in-process vs helper process

**Date:** 2026-11-06  
**Repository:** `C:\Users\StoneWolfPC\HOWL_Chat`  
**Branch / HEAD:** `master` / `9c7d2a1`  
**Mode:** READ-ONLY source and existing-evidence investigation. No application launch, build, test, benchmark, ONNX Runtime initialization, inference, model/voice access, dependency change, configuration write, production change, commit, or push was performed.  
**Deliverable status:** Investigation only. This document does **not** select an architecture, approve a runtime DLL, authorize integration, or amend any prior report.  
**Workspace state:** The working tree was already dirty when this pass began. `git status` shows pre-existing modifications in application, audio, frontend, configuration, and build files; this investigation made none of them. Source citations refer to the working-tree versions actually read, which may differ from `HEAD`. The only file created by this pass is this report.

---

## Evidence vocabulary

| Grade | Meaning here |
|---|---|
| **Source-observed** | Directly visible in the cited file and symbol. This proves code structure, not runtime behavior. |
| **Source-derived** | A conclusion from combining source-observed facts. It remains unexecuted. |
| **Prior evidence** | A result recorded in an existing isolated artifact, with its own scope and limits. It is not re-measured here. |
| **Unestablished** | Requires a build, runtime test, measurement, external approval, or product decision not available in this read-only pass. |

“Verified” below never means “built or run in this investigation.” It means the source or artifact was read directly.

---

## 1. Executive findings

1. **Neither architecture is selected.** Both are technically plausible. They differ less in whether ONNX Runtime can run and more in **where its build, native lifetime, crash radius, and update responsibility live**.

2. **Howl’s release build is explicitly pure Go today.** `build_howl_cgo.bat:13-15` sets `CGO_ENABLED=0` and states that llama.cpp runs as an external `llama-server.exe`, so no CGO, MSVC/clang, or DLL bundling is required. Adding an in-process CGO/ORT dependency directly contradicts that release invariant.

3. **The two build paths do not enforce the same CGO setting.** `build_howl_cgo.bat` sets `CGO_ENABLED=0`; `start-howl-chat.bat` sets only `GOFLAGS=-buildvcs=false` before its auto-build and dev-mode paths. On this machine `go env CGO_ENABLED` is `1` and `CC` is `gcc`. The launcher can therefore produce a CGO-bearing binary that the release script deliberately refuses, or fail on a machine without a C compiler. This is a pre-existing inconsistency; this report does not fix it.

4. **The candidate stack raises the effective build floor to Go 1.25.** `go.mod:3` currently declares `go 1.23.8`. `go-kokoro v0.3.0` and `ortenv v0.2.0` declare `go 1.25.0`; the prior source-level G2P review records `go-g2p v0.2.0` as Go 1.25. The installed toolchain here is `go1.27.3 windows/amd64`, but the release scripts do not pin `GOTOOLCHAIN`. The floor change must be deliberate whether the code runs in-process or in a helper.

5. **Howl already uses sidecar processes successfully, but not a generic helper protocol.** Chat owns `llama-server` on port 8080; image owns `sd-server` on port 8081 and `sd-cli` for video; local Qwen ASR owns a separate `llama-server` on an ephemeral port. The reusable patterns are process-job ownership, HTTP readiness, `Wait()`-based liveness, kill-tree cleanup, and explicit status errors. There is no existing generic request/IPC protocol for a TTS helper.

6. **The audio manager is a strong fit for lazy activation.** `NewProductionManager` constructs no model, binary, or network I/O (`internal/audio/manager_factory.go:9-12`). TTS is initialized only when `Apply` builds a generation (`internal/audio/pipeline.go:723`). Both architectures can therefore avoid loading TTS at startup; a helper can be started only when TTS is explicitly requested.

7. **The same audio boundary favors file-path transfer.** The current TTS contract returns a WAV path, not audio bytes (`internal/audio/pipeline.go`, `synthesizeWithCache`; `types.SynthesisResult.AudioPath`). A helper can write the WAV and return its path on the same filesystem, avoiding an audio-byte copy over IPC. This is a structural advantage, not a measured latency claim.

8. **The main unresolved reliability difference is native failure containment.** An in-process native crash can terminate the entire Wails application. A helper can be detected, killed, and restarted without killing Chat, at the cost of a process, a protocol, and new failure modes. `go-kokoro` also states that native `Run` cannot be interrupted mid-call; a helper can terminate the process to cancel, while an in-process caller cannot safely force a native call to stop.

9. **A passing exact ONNX Runtime probe is not an architecture or packaging approval.** The existing probe demonstrated two sessions and both ownership modes for exact 1.22.0 and 1.25.0 files on this Windows x64 machine; 1.17 failed the binding’s C API request. That establishes a runtime foundation, not an application integration, a production DLL choice, a redistribution right, or a helper design.

10. **Conditional posture, not a selection:** if preserving the `CGO_ENABLED=0` main executable and native-crash isolation outweigh the cost of a new process and protocol, the helper has the lower architecture-risk profile. If avoiding IPC and accepting a deliberate release-contract change are acceptable, in-process remains viable, provided ORT ownership is centralized and the native failure policy is explicit. The next step should be a decision on those constraints, not an implementation.

---

## 2. Candidate boundary

This report concerns the newer candidate pair, not the earlier Python-backed or `go-kokoro-tts` spike:

| Component | Exact candidate | Relevant constraint | Evidence status |
|---|---|---|---|
| Kokoro frontend | `github.com/androiddrew/go-kokoro v0.3.0` | Go 1.25+, CGO, ONNX Runtime | Module `go.mod` read; README read. Windows execution not established. |
| English G2P | `github.com/androiddrew/go-g2p v0.2.0` | Go 1.25, CGO, ONNX Runtime; optional neural fallback | Prior source-level feasibility report. Module was not present in the local cache and was not independently reopened here. Windows build not established. |
| ORT ownership | `github.com/androiddrew/ortenv v0.2.0` | Go 1.25; one process-global environment retained until exit | Module read; prior runtime probe passed for exact 1.22.0 and 1.25.0. |
| ORT binding | `github.com/yalue/onnxruntime_go v1.22.0` | `go 1.19`; CGO; C API 22 | Module read; prior runtime probe passed. |
| Native runtime | Exact 1.22.0 fixture and exact 1.25.0 file | PE x86-64; Microsoft-signed; C API 22 | Prior probe passed both; 1.17 negative control failed. Production approval is open. |
| Text tokenizer | `androiddrew/go-kokoro v0.3.0` tokenizer | Preserves `ᵊ`, fails closed on unknown symbols | Prior C-01 probe passed; not re-run here. |

The earlier `chasedputnam/go-kokoro-tts v1.0.1` + `onnxruntime_go v1.30.0` spike is **not** the current candidate. Its measurements and its Python/eSpeak limitations must not be transferred to this comparison.

---

## 3. Actual Howl constraints

### 3.1 Build and launch paths

| Fact | Evidence | Consequence |
|---|---|---|
| Release build forces `CGO_ENABLED=0` | `build_howl_cgo.bat:15` | Any package that imports C cannot be part of the main executable unless this contract is deliberately changed. |
| Release rationale is “pure Go, no toolchain, no DLL bundling” | `build_howl_cgo.bat:13-15` | An in-process ORT route invalidates the stated reason for the release shape, not just one flag. |
| Release command builds only `.` into one executable | `build_howl_cgo.bat:19` | There is no helper build target or DLL-copy step to reuse. |
| Launcher auto-build sets only `GOFLAGS` | `start-howl-chat.bat:19-20` | It inherits the caller’s/default `CGO_ENABLED`; on this host that is `1` with `gcc`. |
| Launcher dev path also uses the default CGO setting | `start-howl-chat.bat` `:devmode` | `wails dev` and the auto-build are not covered by the release script’s CGO=0 invariant. |
| Module language version is 1.23.8 | `go.mod:3` | The candidate’s Go 1.25 requirement is a toolchain-floor change relative to the declared module. |
| Current host | `go version`, `go env` | `go1.27.3 windows/amd64`, `GOOS=windows`, `GOARCH=amd64`, `GOTOOLCHAIN=auto`, `CGO_ENABLED=1`, `CC=gcc`. This host is not the release contract. |
| No ORT in release bin | `build/bin/onnxruntime.dll` absent | Any in-process or helper ORT route needs a packaging addition; the current release output does not carry it. |

The discrepancy is **not** merely cosmetic. With a CGO-bearing import:

- the release path fails closed at build time because CGO is disabled;
- the launcher path may build a different binary if a compiler is present;
- a machine without a compiler may fail the launcher build while the release script’s documented rationale says no compiler is needed;
- the produced executables can differ in DLL dependencies, size, startup behavior, and crash surface;
- CI or a clean machine can validate one path while users exercise the other.

No fix is proposed here. Resolving this is a prerequisite to claiming that CGO is a controlled property of any Howl build.

### 3.2 Existing process and sidecar patterns

| Existing pattern | Source | What is reusable |
|---|---|---|
| Job-object ownership with `KILL_ON_JOB_CLOSE` | `internal/proc/job_windows.go:28-105`, `:187-197` | A helper can be owned by the app and cannot survive an app crash or force-kill if attached correctly. |
| Suspended start, assign, then resume | `internal/proc/job_windows.go:80-116` | This is the race-free pattern. `proc.Start` launches suspended, attaches, and resumes only after ownership succeeds. |
| Post-start `Attach` | `app-llama.go` spawn path; `internal/image/engine/sdcpp/engine.go` load path | Existing code uses `cmd.Start()` then `proc.Attach(cmd)`. That leaves a small ownership race and is not the best pattern for a new helper. |
| Kill tree | `internal/proc/job_windows.go:164-182` | `taskkill /T /F /PID` is the current descendant-cleanup primitive. |
| Orphan reaping | `internal/proc/job_windows.go:200-218` | Only matches `llama-server.exe` and `sd-server.exe` by path. A new helper name is not covered. |
| HTTP readiness | `internal/audio/asr/llama/llama.go:272-282`; `internal/image/engine/sdcpp/client.go:327-338` | Poll a bounded endpoint, but also watch process exit; readiness and liveness are different questions. |
| One `Wait()` owner | `internal/image/engine/sdcpp/supervisor.go` | A single reaper goroutine owns `Wait()` and publishes exit; other code reads its verdict. |
| Bounded startup | `internal/audio/asr/llama/llama.go:28`, `:120-145` | 120-second ASR startup timeout with exit detection. This is an existing example, not a TTS recommendation. |
| Explicit version check | `app-llama.go:201-217`; `app-image-diagnostics.go:108-145` | `llama-server --version` and `sd-server --help` are used to validate a binary/version. |
| HTTP services | `app-llama.go` port 8080; `internal/image/engine/sdcpp/engine.go` port 8081; `internal/audio/asr/llama/llama.go:78-118` ephemeral ASR port | Local HTTP is the established IPC-like boundary in this application. |
| Per-request subprocess | `internal/audio/tts/kokoro/local.go:254`; `internal/audio/tts/piper/local.go:208` | Existing TTS invokes a CLI/Python process per synthesis with `exec.CommandContext`; it is not a persistent TTS server. |
| Remote TTS HTTP | `internal/audio/tts/melotts/client.go:28`, `:45`, `:237`, `:282` | Existing TTS already supports an HTTP backend with context-aware requests and an 8-minute client timeout. |

The existence of these patterns does **not** prove a new helper is safe. They show what must be built and which lifecycle questions already have precedents.

### 3.3 Audio boundary and lazy activation

- `App.NewApp` constructs `audio.NewProductionManager()` without loading a provider (`app.go:206`; `internal/audio/manager_factory.go:9-12`).
- `ApplyAudioSettings` is the explicit activation path; it validates the desired contract before asking the manager to build a generation (`app-audio.go`).
- `Pipeline.Initialize` loads only requested capabilities. TTS errors are recorded per capability, not as a whole-app failure (`internal/audio/pipeline.go:101-181`, `:723-751`).
- `processManagedAudio` deliberately adds no request deadline (`app-audio.go:42-44`). The manager has leases, operation cancellation, and a bounded close, but a TTS request itself is timeout-neutral.
- `OnShutdown` saves the chat, stops `llama-server`, shuts down the image engine, closes the audio manager, then closes the process job (`app.go:283-313`). Any future ORT teardown must fit this order rather than inventing a competing shutdown path.

### 3.4 Packaging state

`build/bin` contains `howl-chat.exe` and test harnesses but no `onnxruntime.dll`. The release script does not copy a native runtime, a helper executable, model assets, voice assets, or license notices. The repository root contains several DLLs, but a root-level file is not evidence that the release process stages it beside the executable. Both architectures therefore require a packaging design beyond the current build command.

---

## 4. Option A — in-process Go/CGO and ONNX Runtime

### 4.1 Go and CGO contract

**Source-observed:** the candidate modules require Go 1.25.0 and CGO. The project declares Go 1.23.8 and its release script disables CGO. The local host has Go 1.27.3 and `gcc`, but that is a development environment, not a guaranteed release or user environment.

**Consequences:**

- The main executable must be built with CGO enabled, or the CGO-bearing package must be excluded by build tags. Excluding it means the in-process architecture is not present in that binary.
- A C compiler is required on every build machine that produces the release; the current release rationale explicitly avoids that requirement.
- End users normally do not need the compiler, but they do need the selected C runtime and the ORT DLL beside the executable or otherwise resolvable.
- The Go 1.25 floor becomes part of the project’s build contract. The current `go.mod` says 1.23.8; the release scripts do not pin a toolchain.
- A single release artifact now carries native code. Its provenance, signing, and dependency manifest need the same treatment as the existing llama/sd binaries, not only Go module checks.

This is a **release-contract change**, not a local build tweak. It should be decided before any in-process implementation.

### 4.2 DLL discovery and packaging

`onnxruntime_go` exposes `SetSharedLibraryPath(path)` and requires it before `InitializeEnvironment()`; otherwise the default is `onnxruntime.dll` on Windows (`onnxruntime_go.go:65-70`, `:85`). `ortenv.Init(library)` accepts an explicit path or an OS-loader name, normalizes explicit paths, and pins the selector (`ortenv/runtime.go`, `Init`).

For an app-local deployment:

- stage the exact `onnxruntime.dll` next to the executable, or resolve an explicit absolute path;
- ensure it is the intended architecture (x86-64 in the tested files);
- stage the CRT dependencies it imports (`MSVCP140`, `MSVCP140_1`, `VCRUNTIME140`, `VCRUNTIME140_1`, and the CRT API sets as applicable);
- keep provider DLLs such as `onnxruntime_providers_shared.dll` out unless a provider actually requires them; the prior CPU probe did not need that file;
- do not rely on `System32` containing a compatible ORT. The 1.17 System32 control failed because it did not provide the binding’s requested C API;
- do not assume a successful DLL load proves every operator, provider, model, or voice works.

The exact 1.22.0 and 1.25.0 files passed the isolated CPU probe, including signature and import checks. That is evidence for those files only. It is not a redistribution approval and does not establish that a different 1.25 build or a future runtime is compatible.

### 4.3 Runtime ownership and process lifetime

ONNX Runtime has one environment per process. `ortenv` retains that environment and its native library until process exit; it deliberately never unloads them because CUDA providers can retain stale pointers (`ortenv/README.md`, “Why retain the environment?”). `go-kokoro` documents two modes (`go-kokoro/README.md:185-207`):

- set `ORTLibrary`: the engine initializes through `ortenv` and the environment remains until process exit;
- leave it empty: initialize `onnxruntime_go` yourself before `New`, and destroy only after every engine is closed.

The prior probe demonstrated both modes with exact 1.22.0 and 1.25.0, two sessions on one environment, and clean process exit. It also demonstrated that a second direct initialization is rejected and that 1.17 cannot satisfy the binding.

The important consequences are not the happy path:

- all in-process ORT users must agree on one selector and one owner;
- a direct caller and `ortenv` cannot both own the environment;
- destroying an environment outside `ortenv` forces a process restart rather than reinitialization;
- changing the DLL requires updating every in-process ORT user together;
- a future in-process ASR, image, or embedding path using ORT creates a shared-lifetime dependency that does not exist today;
- sessions, tensors, and options must still be destroyed in the correct order even though the environment is retained.

An in-process TTS integration that initializes ORT at first `Apply` and never unloads it would match the documented model. It would also make the ORT selector and version a property of the whole application, not of one provider.

### 4.4 Startup, shutdown, and failure

**Startup:** the current manager is provider-free, so in-process TTS can remain lazy. `Apply` would build the pipeline, load the model, create the ORT session, and report a capability error if any step fails. It need not block the initial Wails window.

**Shutdown:** `Pipeline.Shutdown` unloads provider-owned models through the loader (`internal/audio/pipeline.go:523-548`), but under `ortenv` the ORT environment remains until process exit. If the caller owns the environment, `DestroyEnvironment` must run after all sessions and inference calls are finished. The existing `OnShutdown` order is chat save → chat server kill → image shutdown → audio close → job close (`app.go:283-313`); an in-process ORT teardown must be inserted deliberately, not inferred.

**Initialization/session failure:** `initTTS` loops candidates, records an error, and leaves `TTSReady` false (`internal/audio/pipeline.go:723-751`). The manager can reject a TTS-only build while retaining an existing active generation, according to its source-level state machine. That is a useful containment boundary for ordinary errors. It is not a containment boundary for a native access violation, abort, or corrupted native state: those can terminate the process, and Go cannot recover from all of them.

**Can Chat continue without TTS?** At the source level, yes: construction is provider-free, chat does not depend on the audio manager, and a failed TTS capability is recorded separately. This has not been runtime-tested for this candidate, and it does not cover a native crash.

### 4.5 Upgrade and maintenance

- Every ORT update is an application update, not a provider update.
- The binding (`onnxruntime_go v1.22.0`) and the exact native runtime must be tested together; the prior probe showed C API 22 acceptance for 1.22.0 and 1.25.0, and rejection for 1.17.
- The Go 1.25 floor and CGO toolchain become release prerequisites.
- A future in-process ORT consumer must coordinate initialization and shutdown with TTS.
- Model/voice changes remain separate, but packaging and exact-version approval are coupled to the application release.
- Debugging involves one process, which is simpler for ordinary errors; a native fault has a larger blast radius.

---

## 5. Option B — separate helper process

### 5.1 Independent build and distribution

A helper can be a separate executable that imports the CGO/ORT stack while the main Howl executable does not. This is technically possible without changing the main binary’s `CGO_ENABLED=0` setting **if** the main package does not import the helper’s packages.

There are two packaging forms:

1. **Separate Go module** for the helper: keeps its Go 1.25 directive and dependency graph out of Howl’s module. This is the cleanest way to preserve the main module’s floor, but adds a second module, versioning, build, test, and release artifact.
2. **Additional package/target in the existing module**: simpler file layout, but the module’s Go directive/dependency graph and tooling now include the helper’s requirements. The main package can still avoid importing it, but `go mod tidy`, CI, and toolchain selection see it.

Either form needs:

- its own CGO build step and compiler;
- its own app-local `onnxruntime.dll` and CRT dependencies;
- its own exact-version/version-handshake check;
- model and voice assets available under a defined directory contract;
- license/notice files for the helper and native runtime;
- a release step that stages the helper even when TTS is optional;
- a failure policy for a missing helper (TTS unavailable, not application failure).

No Kokoro/G2P helper target exists today. `build_howl_cgo.bat` builds the main application; `start-server.bat`/`.sh` locate and run `sd-server.exe`; `build/bin` contains no helper and no ORT DLL.

### 5.2 IPC options

| Option | Existing precedent in Howl | Advantages | Costs/unknowns |
|---|---|---|---|
| Loopback HTTP | `llama-server` 8080, `sd-server` 8081, ASR llama ephemeral port, MeloTTS HTTP client | Already understood; JSON; context-aware requests; `/health`; easy to test with ordinary tools | Needs a port, bind/security policy, request framing, auth/authorization, and collision handling. Local processes can reach loopback. No performance numbers are established. |
| Child stdin/stdout | `os/exec`; `StdoutPipe` already used for sd-server progress | No listening port; lifetime tied to the child; simple for one client | Needs a framed protocol, deadlock/backpressure discipline, and explicit cancel semantics. A child stdout used for both logs and protocol needs separation. |
| Windows named pipe | No named-pipe transport found in the inspected Go source; `os.Pipe` appears only in a test | Native Windows namespace; no TCP port | New platform-specific code and ACL/security design; not an existing pattern here. |
| Wails binding/event bridge | Used between the Go app and the WebView renderer | Already present for UI calls | Not a helper transport. It does not connect the Go backend to an external process. |

The existing HTTP pattern is the lowest-new-concept option, but it is not automatically the best. A local HTTP server must bind loopback only, use a per-launch token or equivalent authorization, and avoid exposing model/text data to unrelated local processes. The existing llama/sd servers are trusted local services with no such token in the inspected code; a new helper should not copy that omission without a decision.

Any IPC design must carry at least: protocol version, request ID, operation (`synthesize`/health/version), text, voice, speed, engine, output format, cancellation, result path or bytes, and typed errors. It must define maximum message size and unknown-field behavior. These are design requirements, not implementation in this report.

### 5.3 Startup, health, shutdown, crash detection, and restart

**Startup:** the helper can be started only when TTS is explicitly applied. The audio manager’s provider-free construction and `Apply` path support that. Startup should use `proc.Start` (suspended → job assignment → resume), not `cmd.Start` followed by `proc.Attach`. The ASR adapter already demonstrates the safer pattern (`internal/audio/asr/llama/llama.go:104-115`).

**Readiness:** separate process liveness from service readiness. The image engine’s supervisor exists because polling readiness alone cannot notice a dead child (`internal/image/engine/sdcpp/supervisor.go`). A helper should have exactly one `Wait()` owner and a bounded health probe; readiness failure or child exit should be reported distinctly, with the child’s stderr tail where possible.

**Shutdown:** request a graceful stop, wait with a bound, then kill the tree and close the job. `proc.KillTree` and `CloseJob` are available; the helper must not outlive the app. Do not call `Wait()` from two goroutines.

**Crash detection:** the `Wait()` result is the authoritative process-exit signal. A closed IPC pipe is a hint, not a substitute for `Wait()`.

**Restart:** the existing application does not provide a generic automatic restart supervisor. `sd-server` can be recreated by `EnsureReady` when it is not running (`internal/image/engine/sdcpp/engine.go:239-279`), and `llama-server` ASR reports a failed initialization rather than silently restarting. A TTS helper needs an explicit policy: restart on next request, bounded restart with backoff, or mark TTS unavailable until the user retries. No restart policy or timing is established here.

### 5.4 Cancellation and timeouts

Existing TTS CLI/Python paths use `exec.CommandContext` (`internal/audio/tts/kokoro/local.go:254`; `internal/audio/tts/piper/local.go:208`), and MeloTTS uses context-aware HTTP (`internal/audio/tts/melotts/client.go:237`, `:282`). The ASR HTTP path has a 15-minute client timeout and a 120-second startup timeout. The shared `processManagedAudio` path deliberately has **no request deadline** (`app-audio.go:42-44`).

A helper protocol therefore needs:

- a bounded connect/startup timeout;
- a per-request deadline or an explicit caller-controlled cancellation;
- a cancel message with request ID, followed by a bounded wait;
- a kill-and-restart fallback if the helper does not acknowledge cancellation;
- a rule for late results after cancellation (discard, not deliver);
- a health check that cannot block the Chat UI thread.

`go-kokoro` documents that native `Run` cannot be interrupted mid-call (`go-kokoro/README.md`, provider/lifetime section). In-process cancellation can stop before or after the native call but cannot safely preempt it. A helper can terminate the process, but that throws away the loaded model and forces a restart. This is a real architectural tradeoff, not a free win.

### 5.5 Audio transfer and memory

The current pipeline returns `AudioOutputPath`; the Wails binding gives the renderer a file path. A helper can write a WAV in a shared temp/cache directory and return the path, matching that contract. This avoids serializing the waveform through IPC and avoids holding a second full copy in the main process.

If streaming is required, the existing TTS implementations read the complete file and emit 32 KiB chunks (`internal/audio/tts/piper/local.go`, `SynthesizeStream`; `kokoro/local.go`, same pattern; `melotts/client.go:150-181`). A helper could stream bytes or path-backed chunks, but that requires protocol and lifetime design. No latency, throughput, or memory measurement exists for either option, and none should be inferred from the existing code.

Shared filesystem assumptions matter: both processes must use the same temp/cache root, cleanup must be owned by one side, and a path returned after helper death must not be treated as valid indefinitely.

### 5.6 Missing or incompatible helper, and independent updates

A helper design needs a deterministic discovery contract:

- an app-relative `helper/` or `runtime/tts/` path, with an override for development;
- a version/protocol handshake before the first synthesis;
- `--version` or a health payload containing helper version, protocol version, and ORT version;
- a typed “missing,” “incompatible,” and “failed to start” status;
- a clear UI/status message rather than a generic synthesis failure.

Existing examples are `ValidateLlamaServer` (`app-llama.go:201-217`) and `GetSDServerVersion` (`app-image-diagnostics.go:108-145`). Both are binary checks, not a full protocol contract.

A helper can be updated independently in principle: replace the executable and its app-local DLL, keep a compatible protocol, and let the app validate. In practice, the model/voice assets, ORT C API, and protocol version are a compatibility triple. An “independent” update is not independent of those contracts.

### 5.7 Failure isolation and its limits

A separate process prevents a helper crash from directly taking down the Wails process, provided the main process treats helper failure as a recoverable TTS capability failure and does not block on a dead pipe. That is a meaningful reliability advantage over in-process native code.

It is not total isolation:

- the helper still shares CPU, RAM, disk, and any GPU/VRAM used by the selected provider;
- the app’s existing resource manager does not govern TTS, and the audio pipeline has no shared admission gate;
- a hung helper can still block a request if the protocol has no timeout;
- a helper started under the app’s job object is killed when the app exits, so it cannot be a durable background service;
- a helper port or pipe is a local attack surface;
- a helper update can break the app if protocol/version checks are incomplete.

No runtime test in this investigation has exercised a helper crash, restart, or audio transfer.

---

## 6. Comparison matrix

| Criterion | In-process Go/CGO | Separate helper | Basis |
|---|---|---|---|
| Main release stays `CGO_ENABLED=0` | **No**, unless the CGO code is excluded entirely | **Yes**, if main does not import helper packages | `build_howl_cgo.bat:15`; package boundary |
| Build toolchain | Main build needs CGO compiler; Go 1.25 floor | Helper build needs CGO compiler; main can avoid it | `go.mod:3`; candidate module requirements |
| Packaging surface | One main exe + ORT DLL + CRT + notices | Main exe + helper exe + ORT DLL + CRT + notices + protocol | `build/bin` currently has neither |
| ORT ownership | Process-global; every in-process user must cooperate | Isolated to helper process | `ortenv/README.md`; `go-kokoro/README.md` |
| Native crash radius | Whole Howl application | Helper; main can report TTS unavailable | Source-derived; not crash-tested |
| Chat impact of TTS failure | Ordinary errors can be per-capability; native faults can kill app | Helper failure can be contained if IPC is bounded | Source-derived; no failure injection |
| Startup cost | No new process; model/ORT load still lazy | Extra process + handshake; can be lazy | `manager_factory.go:9-12` |
| Lazy activation | Yes, on `Apply` | Yes, on `Apply` | `pipeline.go:723` |
| IPC | None | Required: HTTP, stdio, or named pipe | No generic helper protocol exists |
| Audio transfer | File path or in-process buffer | File path is natural; bytes/stream add protocol work | `pipeline.go`, TTS `AudioPath` |
| Cancellation | Context can stop before/after; native `Run` is not interruptible | Can kill helper, but loses loaded model and requires restart | `go-kokoro` README; `exec.CommandContext` precedent |
| Restart after failure | Restart the whole app | Restart helper; policy not established | `sd-server` has recreation, not a generic supervisor |
| Independent update | Coupled to app release | Possible with protocol/version checks; assets still coupled | Source-derived |
| Existing pattern fit | New in-process native runtime; no current ORT user | Fits existing sidecar/HTTP/job patterns, but needs a new protocol | `app-llama.go`, `sdcpp`, ASR llama |
| Testability | Unit tests can fake TTS; native lifecycle needs process tests | Protocol and lifecycle need two-process tests | Existing test seams |
| Maintenance | One ORT owner/version, one app release | Two artifacts and a protocol, but isolation | Source-derived |
| Security | DLL/CRT trust and native code in main process | Adds local IPC endpoint/pipe and helper binary trust | Source-derived |
| Performance | No IPC; no measurement | Potential IPC/process cost; no measurement | **Unestablished** |

No row is a performance measurement. “Lower cost” is not claimed without numbers.

---

## 7. What remains unestablished

1. **Windows build of the candidate pair.** `go-kokoro` documents Linux testing; `go-g2p` Windows/CGO/ORT build remains unproven. The prior ORT probe used only the binding’s stock model, not G2P or Kokoro.
2. **Exact DLL approval.** 1.22.0 and 1.25.0 passed this machine’s isolated CPU probe; production version, redistribution terms, signing, and update policy remain open.
3. **CGO/release decision.** Whether Howl will intentionally become CGO-bearing, remain pure Go with a helper, or use build tags is a product/build decision not made here.
4. **Failure injection.** No test has established whether ORT initialization failure, session failure, or a native crash leaves the app usable. Source-level capability isolation is not crash isolation.
5. **Helper lifecycle.** No helper exists; startup, handshake, crash detection, restart, and shutdown are design requirements, not demonstrated behavior.
6. **IPC/audio measurements.** Startup latency, request latency, throughput, memory duplication, and audio transfer cost are unmeasured. Any numbers require a later controlled test.
7. **Resource contention.** TTS is outside the shared resource manager. A helper does not automatically prevent CPU/RAM/VRAM contention.
8. **Protocol/update compatibility.** No versioned helper protocol or compatibility matrix exists.
9. **G2P quality and integration.** This report does not evaluate pronunciation, tokenizer behavior, or synthesis quality; prior evidence is limited to the recorded probes.

---

## 8. Disposition

**No architecture is selected.**

- **In-process** is technically viable only as a deliberate release change: CGO enabled, Go 1.25 floor accepted, one ORT owner, exact app-local runtime, packaging/CRT/notices added, and an explicit policy for native failure. It offers the fewest moving parts and no IPC, but its failure radius is the application.
- **Helper** is technically viable and preserves the main executable’s pure-Go release shape if the helper is a separate package/module and is not imported by `main`. It offers crash containment and independent update potential, but requires a second artifact, a protocol, process lifecycle, IPC security, and measured cancellation/audio behavior. It is not automatically safer merely because existing sidecars exist.

The decision should be made on the constraints this investigation surfaced: **release purity, crash containment, and update independence versus protocol complexity and a second artifact**. The next step is a bounded decision record, not implementation. If a spike is authorized later, it should measure the chosen transport and failure behavior rather than assume them.

---

## 9. Source and evidence index

**Build and project**

- `build_howl_cgo.bat:13-19`
- `start-howl-chat.bat:19-20`, `:devmode`
- `go.mod:3`
- `app.go:206`, `:283-313`
- `build/bin/` listing (no `onnxruntime.dll`)
- `go version`; `go env GOOS GOARCH GOVERSION GOTOOLCHAIN CGO_ENABLED CC`

**Audio boundary**

- `internal/audio/manager_factory.go:9-12`
- `internal/audio/manager.go` (`Apply`, `Acquire`, `Close`)
- `internal/audio/pipeline.go:101-181`, `:523-548`, `:723-751`
- `app-audio.go:23-56`, `:175-198`
- `internal/audio/tts/kokoro/local.go:232-254`
- `internal/audio/tts/piper/local.go:183-208`
- `internal/audio/tts/melotts/client.go:28`, `:45`, `:150-181`, `:237-282`
- `scripts/kokoro-tts.py` (current Python/ONNX bridge invoked per synthesis)

**Process and sidecar patterns**

- `internal/proc/job_windows.go:28-116`, `:164-218`
- `app-llama.go:201-217`
- `internal/image/engine/sdcpp/engine.go:239-279`, `:413-430`
- `internal/image/engine/sdcpp/supervisor.go`
- `internal/image/engine/sdcpp/client.go:327-338`
- `internal/audio/asr/llama/llama.go:28`, `:78-145`, `:272-282`
- `app-image-diagnostics.go:108-145`

**ONNX Runtime and candidate**

- `build/tmp/ort-windows-runtime-probe/RESULTS.md`
- `build/tmp/ort-windows-runtime-probe/modcache/github.com/yalue/onnxruntime_go@v1.22.0/onnxruntime_go.go:65-70`, `:85-135`
- `build/tmp/ort-windows-runtime-probe/modcache/github.com/androiddrew/ortenv@v0.2.0/runtime.go`
- `build/tmp/ort-windows-runtime-probe/modcache/github.com/androiddrew/ortenv@v0.2.0/README.md`
- `build/tmp/kokoro-tokenizer-probe/modcache/github.com/androiddrew/go-kokoro@v0.3.0/go.mod`
- `build/tmp/kokoro-tokenizer-probe/modcache/github.com/androiddrew/go-kokoro@v0.3.0/README.md:185-207`
- `build/tmp/kokoro-g2p-feasibility/G2P-TECHNICAL-FEASIBILITY.md` (prior source-level G2P evidence; Windows build unestablished)
- `docs/review/AUDIT-kokoro-g2p-feasibility-challenge-2026-11-06.md` (prior audit; CGO discrepancy corroboration)

**Not evidence for this candidate**

- `build/tmp/native-kokoro-spike/RESULTS.md`: different `go-kokoro-tts v1.0.1` candidate, different binding, Python/eSpeak limitations, and a small-sample performance run. Its numbers must not be applied to `go-kokoro v0.3.0` or to either hosting architecture.

### 9.1 Source hash manifest (SHA-256)

These are SHA-256 hashes of the exact working-tree bytes read for the cited source files. They identify the inspected versions despite the pre-existing dirty worktree; they are not Git blob IDs and may differ from `HEAD`. Command-output facts (`go version`, `go env`, directory listings) are not file hashes.

| Inspected file | SHA-256 |
|---|---|
| `build_howl_cgo.bat` | `59e860d117796c2d63b7816c588e2172c999bac5cd005f09828981e35f251e65` |
| `start-howl-chat.bat` | `df04082b2f6619622e1fddcec3623a143697b48d61a972182d4b4b405181216e` |
| `go.mod` | `1c4e45eb66fe74bdee41ef3049f242f0a1fdb185f95e000938272621189644c2` |
| `app.go` | `41197643f04b798b10b803119cf4f019304627f0f5eb3f1f835099b4303ba644` |
| `app-audio.go` | `ff5041f19c7264cbc3cfae4baa91b49d1f58189c021a9c8155bb1833713cc9d9` |
| `app-llama.go` | `5ee8531edd7a54937297c53a56cf8b9ace40cbfe86e54caacbd9238108978331` |
| `app-image-diagnostics.go` | `4c7b3e97aa1ab8c9228371cbc7da680dafc50c690c4d5186d7a1e7ba4db70e7f` |
| `scripts/kokoro-tts.py` | `b75723ac47cf3a525dcda0eba05c81292eec790d14a8e5155e260e0ca58db6d1` |
| `internal/audio/manager_factory.go` | `89eb0d0d3f60ec8a07a7d79471b21fc899de3ab259b1e982ec5190183d697808` |
| `internal/audio/manager.go` | `ee7a4827f08ffe13af0bcce2cd6d0a5a7c800c79c6902d12bfda79be28b5ef3b` |
| `internal/audio/pipeline.go` | `ae386c36f8112e02d78f05ca1304cb5e9f4faf0551e64beeb33e18ef86fdfaba` |
| `internal/audio/tts/kokoro/local.go` | `4b954e2c970a7baf4305333c345927455dbb0043fe4d9287e902eb26fe7a5bf6` |
| `internal/audio/tts/piper/local.go` | `9e192bbe52599d954d1fe6d91646e8780b85e7cf7a3bce32123d87f41285892f` |
| `internal/audio/tts/melotts/client.go` | `13e7b93041482f04a953fe3c1c32aafa5f3347db5f33646143cc73c69c52758f` |
| `internal/proc/job_windows.go` | `1edef75a7fc753b949707cb1066056ef137cfc511afca9fea53f28ecee663c46` |
| `internal/image/engine/sdcpp/engine.go` | `803dae0e069526fca1341b95335466ce9b1cdd4ab5fe098bdde82d397301d60d` |
| `internal/image/engine/sdcpp/supervisor.go` | `af23915040b1bcd057d5c5cb3c58809edd97c5c86a2197646b39abf9204e2417` |
| `internal/image/engine/sdcpp/client.go` | `e690734df3ecece30416bfbd1f467e22b5093f375ae3c49e9a77982015b38e58` |
| `internal/audio/asr/llama/llama.go` | `eda8974ac87895e9c22a28485f37eda2b0f28fa303852aedaf2ef6f37351c5da` |
| `build/tmp/ort-windows-runtime-probe/RESULTS.md` | `3a549dbbef5a4c90df4a8300f26b6921d78d4e11984a9c570daa15cedf502fde` |
| `build/tmp/ort-windows-runtime-probe/modcache/github.com/yalue/onnxruntime_go@v1.22.0/onnxruntime_go.go` | `2fdc5e171b60262301794cfe063164ebeafd42a66b3e5bfae534aae143570933` |
| `build/tmp/ort-windows-runtime-probe/modcache/github.com/androiddrew/ortenv@v0.2.0/runtime.go` | `89e97e5176c607f438845d429ef49151b7780c3c10ad1635e9b334f9b7fcdf42` |
| `build/tmp/ort-windows-runtime-probe/modcache/github.com/androiddrew/ortenv@v0.2.0/README.md` | `ad16f3523a2547034fe824bbc2fcc885383cc0435f9baaa8cc585aa28fa25965` |
| `build/tmp/kokoro-tokenizer-probe/modcache/github.com/androiddrew/go-kokoro@v0.3.0/go.mod` | `174dd3e595ae1924ed9228f1fc82579e97b02de553bd53077965f89b81aef44a` |
| `build/tmp/kokoro-tokenizer-probe/modcache/github.com/androiddrew/go-kokoro@v0.3.0/README.md` | `db8cd89ad41e168d1da9cd980d6ff056c97f1b2095835e38021d5936aad3cf1d` |
| `build/tmp/kokoro-g2p-feasibility/G2P-TECHNICAL-FEASIBILITY.md` | `c7eb7667bac5be3f83ccc7bf4b048def7cff309290f61a5f2a0dd738f048a26c` |
| `docs/review/AUDIT-kokoro-g2p-feasibility-challenge-2026-11-06.md` | `6c62db923384d09ba8a79986c95ea94089d70a77ce62a66bf16bdb1f72b25935` |

---

## 10. Explicit non-claims

- This report does not approve ONNX Runtime 1.22.0, 1.25.0, or any distribution.
- This report does not approve in-process hosting or a helper process.
- This report does not establish that a helper is more secure or faster.
- This report does not establish that an in-process failure is recoverable.
- This report does not establish Windows behavior for `go-g2p` or `go-kokoro`.
- This report does not alter any production code, build script, dependency, model, voice, or prior research artifact.
