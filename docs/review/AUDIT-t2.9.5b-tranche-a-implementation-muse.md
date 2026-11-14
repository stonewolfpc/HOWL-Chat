# T2.9.5B Tranche A — Independent Implementation Audit (Muse, read-only)

**Status:** independent implementation audit. Read-only; no source, test, config, or
documentation file was modified, built, executed (except `llama-server.exe --help`
as a read probe and byte-level string scans), or committed by this audit.
**Authority:** accepted RECON-2 r2.1 (`docs/review/RECON-t2.9.5b-deepseek-muse-audio-r2.md`);
the controlling authorization register r14 (`bd10519b…c5f94cc44`, sidecar verified
byte-identical this pass) is untouched and unchallenged.
**Implementation report audited:** `docs/audio-architecture/t2.9.5b-tranche-a-implementation.md`.

## 1. Baseline integrity

- HEAD `4f5defe` matches the report's stated starting commit. 34 tracked modifications
  confirmed by `git diff --name-only`. The 10 files the report names as newly modified
  are all present in the diff.
- All 12 SHA-256 hashes in report §12 re-computed independently: **all match**.
- `docs/review/` contains no tracked content and no tranche edits; register r10–r14
  sidecars intact (r14 and r13 re-hashed this pass, match).
- The 24-pre-existing / 10-new split cannot be re-proven from VCS (no start snapshot),
  but file mtimes corroborate the account: the 11-02 batch (whisper recognizer,
  asr.go, fallback chain, qwen client, loader, TTS adapters, app.go) predates the
  11-04 tranche batch. Taken on report with that corroboration.
- **Report inventory inaccuracies (documentation, low severity):**
  - §2 "Modified" lists `internal/audio/manager.go`, `manager_factory.go`,
    `model_resolver.go` — these are **untracked new files**, not tracked modifications.
    `manager.go` alone is ~1100 lines of generation/lease/operation machinery that the
    report's §8 describes as pre-existing ("the manager still owns…"). Whether authored
    11-02 or 11-04, it is new to version control and should be listed as new.
  - §2 "New" omits new test files (`asr_selection_test.go`, `asr_registry_test.go`,
    `app_audio_capability_test.go`, manager/resolver/TTS test files) and new non-test
    files (the three above).
  - `docs/audio-architecture/README.md` was edited in the tranche window (+32 lines,
    including the new capability table) but is counted in neither the 10 nor discussed.
- No test, build, or runtime verification was re-executed by this audit (read-only
  constraint); report §9's claimed `go test`/`go build` results are therefore
  **taken on report**, not independently confirmed. The claim files all exist.

## 2. Capability contract (`internal/audio/capability`)

Design is honest; enforcement reach is overstated.

**Confirmed correct.**
Tri-state semantics hold: `unknown` blocks via `Dimension.OK()`; `Selectable()`
requires every dimension. Explicit selection (`Select`) is honoured-or-refused with
full blockers and never substituted. `Fallbacks` enforces same-task/same-locality
(no local failure can silently become a remote API call). Verification grades mean
what the report says (`implementation` = adapter wired, not run; nothing claims
`end_to_end`). Remote Qwen fails closed without `HOWL_QWEN_API_KEY`. Local Qwen
(`SupportsCPU:false`, GPU-only) blocks on CPU-only hosts. `Blockers()` reasons reach
the settings UI verbatim. No races: `Catalog` is rebuilt per call (no stale snapshot
cache); GPU detection is `sync.Once`; pipeline mutex discipline around
`initialization`/`status` is intact.

**F-C1 (moderate): auto-selection bypasses the contract it advertises.** Report §3:
"automatic selection considers only selectable records, prefers verified then
implemented then local." In production, `asrCandidatesForRecords` (pipeline.go:291)
auto branch filters the planner list by `Local && ApprovedLocal` only — it never
consults `Selectable()` or any dimension, and the freshly built `records` are
discarded unused in that branch. `capability.AutoSelect` and `capability.Fallbacks`
have **zero production callers** (tests only). Consequence: auto can plan an engine
whose binary is absent, whose model file is missing, or whose hardware is ineligible;
failure surfaces at load, not at selection. The contract today governs *explicit*
selection and *reporting*, not automatic selection.

**F-C2 (moderate): explicit TTS is registration-checked, not contract-checked.**
`ttsCandidatesFor` (pipeline.go:318) verifies only registry membership, never
`capability.Select`. An explicit TTS engine with blocked dimensions (e.g. missing
voice file) is refused at load, not with the contract's blockers. Report §5
("explicit TTS is honoured or refused") overstates the mechanism; the outcome is
equivalent (refusal, no substitution) but the cited enforcement point is wrong.

**F-C3 (significant — factually wrong string in shipped code):**
`ProductionEnv` hardcodes `LlamaAudioRoutes: false`, and the contract reports
"the bundled llama-server build exposes no audio request routes" as fact. §6 shows
this is unestablished and the binary evidence points the other way. A hardcoded
`false` goes silently stale if the runtime is replaced; this must become
unknown-pending-probe (or true, after the §6 re-probe passes).

**F-C4 (low): two path policies; resolver unwired.** Pipeline loads via
`defaultASRModelPath` (pipeline.go:664); the new `Resolver` uses its own
nested/flat search and has **no production callers** (tests only). The resolver
"cannot disagree" only because nothing reads it. Wire it into the load path or
remove it; same for `AutoSelect`/`Fallbacks` (wire or downgrade report §3).

## 3. DEBT-T1 — projector classification: PASS with observation

`gguf.InspectProjector`/`ClassifyProjectorMetadata` verified line by line.
Keys are the genuine families (`clip.has_vision_encoder`, `clip.has_audio_encoder`,
`clip.vision.*`/`clip.audio.*`); explicit-false veto works (a `false` flag wins over
key-space fallback); empty/unreadable/contradictory metadata fails closed to
`unknown` with vision disabled. `HasVision`/`HasAudioInput`/`ProjectorKind` are set
atomically from one inspection in `app-model.go` — they cannot diverge. Existing
vision behaviour preserved for declaring projectors; error text names the
projector's declaration. Tests build **real minimal GGUF files** parsed through
`gguf-parser-go` (not assumption replays), covering audio-must-not-enable-vision,
vision-preserved, flags, key spaces, and explicit-false.
Observation: a projector with silent metadata loses vision gating (report debt #4
discloses this exactly); the `clip.audio.*` fallback key space is unestablished for
real audio families — safe only because failure is closed. No filename inference
anywhere in the path.

## 4. DEBT-T3 — registry/resolver/status/UI: PASS with findings

Genuine fixes confirmed: `approvedASRModels` derives from the registry (resolver
cannot admit what the registry does not); stale/non-identities removed from UI and
resolver (`whisper-medium`, `whisper-large`, `qwen-asr`); UI now offers exactly
registry identities + `auto` (verified in `audio-settings.html` diff);
`GetAudioStatus` carries the contract's records and the UI renders reasons;
`ApplyAudioSettings` is registered in Go, Wails bindings (`.d.ts`/`.js`), and the
save path; explicit engines reach `PipelineConfig` and refuse-with-reason.
`config/model-registry.json` is the image-model scan cache — **not** a duplicate
audio registry. `initMultimodal` fails honestly when multimodal is requested.
Dead constants `ASRWhisperMed/Large/Fallback` remain (one benign use in a token
heuristic) — recommend removal.

**F-T3-1 (moderate): saved engines are not re-applied at startup.** The backend
starts with `NewProductionManager()` and zero `Settings`; nothing applies persisted
choices on launch. The frontend persists `asrEngine`/`ttsEngine` (settings.json +
localStorage) and restores the dropdowns — so after a restart the UI can display an
engine the backend has not applied until the user re-saves. Report debt #9 says
persistence is "not wired" but the frontend half *is* wired; the precise gap is
backend re-application. Document or fix; do not let UI and backend disagree silently.

**F-T3-2:** auto-ASR dimension bypass, same root as F-C1.

## 5. DEBT-T2 — streaming honesty: PASS

All six registry entries declare `SupportsStreaming: false`; all three adapters
(whisper, DashScope, fallback chain) return unsupported from `TranscribeStream`;
`planASR` filters on `RequireStreaming`; the updated planner test asserts an empty
plan (exercising the real `BuildSelectionPlan`, not a mock); the registry/agreement
test fails on any future declaration-without-adapter, including count drift. No
other declared-but-unimplemented capability found: TTS `SynthesizeStream` exists on
all three providers (Piper's `SupportsStreaming:true` is backed; MeloTTS/Kokoro
`false`-but-implemented is conservative), and the contract correctly scopes itself
to file-transcription/streaming/file-adjacent tasks.

## 6. Native-ASR gate: evidence FAILS, conclusion stands on narrower grounds

**F-G1 (significant): the gate scanned the wrong binary.**
Report §7 evidence #1 ran string scans on `runtime/llama/llama-server.exe` (9,728
bytes, mtime 2026-10-30 10:49). That file is a launcher **stub**: `MZ` header,
imports `llama-server-impl.dll`, exports a `llama_server` entry import. The actual
server is `llama-server-impl.dll` (10,088,960 bytes), which **contains** the audio
request surface: `/v1/audio/transcriptions` (1×), `input_audio` (8×),
`input_audio_buffer_*` handlers, `Request converted: OpenAI Transcriptions ->
OpenAI Chat Completions`, `No input file found for transcription`,
`Only 'json' response_format is supported for transcription`, plus realtime
transcription-session strings. `mtmd.dll` (2.1 MB) contains `clip_has_audio`,
`has_audio`, and qwen2/omni-era symbols. Zero-matches on the stub prove nothing
about the server. The `--help` observation (mmproj flags only) is inconclusive:
audio HTTP endpoints in recent llama.cpp need no dedicated CLI flags.
Grade: whether the bundled build serves the audio route is now **unknown**, with
binary evidence leaning present — not "established absent."
Required: re-probe against the implementation (route table or a no-model
404-vs-405 probe of `POST /v1/audio/transcriptions`), then correct report §7/§10
debt #1/§11, and fix F-C3. Owner decision #1 ("adopt a runtime build with the
audio request surface") rests on a premise now in doubt and must not proceed on
the stub scan.

**What still holds (verified this pass).** No model assets: `models` → broken
symlink (`/e/HOWL_Chat/models` absent), no `.gguf` in tree (report's path form
`E:\HOWL_Chat\models` matches). No `whisper-cli` on PATH, no `bin/whisper/`.
Local Qwen is safely unavailable at four layers: registry `ApprovedLocal:false`,
contract blocks (provider+model+task), no adapter registered in
`registerBuiltins` (load fails closed), resolver excludes it (no `LocalFile`).
Remote Qwen cannot leak into local auto (admission filter + `PreferredBackend`
local + `RecommendASR` admitted-only). E4 and E5 are not claimed — confirmed.

## 7. Regressions and debt movement

- Tranche scope respected: chat untouched except image-error strings; Whisper
  lifecycle/limiter/temp-file changes are the 11-02 pre-existing batch, preserved
  with tests; TTS voice path (`ensureTTSVoice`/`voiceModelPath`/`chooseASR`)
  untouched by tranche (empty `git diff -G`).
- Pre-existing DEBT-T9 (re-init temp-dir orphan in `recognizer.go`) is untouched by
  the tranche and correctly not claimed fixed — still open, still provisional.
- `EnableAudioASR(true)` resetting TTS desire matches old behaviour (old code also
  rebuilt with `TTSEnabled:false`) — no regression.
- Debt accounting is accurate except report debt #1's root cause ("runtime predates
  or excludes the ASR surface") — invalidated by F-G1; rewrite as "audio surface
  presence unestablished pending impl-DLL/route re-probe; assets absent regardless."
- Debt moved, not eliminated, and correctly labelled: static (not negotiated)
  facts, MeloTTS `unknown` reachability, unverified remote Qwen, legacy status keys,
  D5 persistence — all visible in-contract as `unavailable`/`unknown`/sub-`end_to_end`.

## 8. Required corrections (dependency order)

1. **Re-probe the runtime audio surface** (impl DLL route table or no-model
   404-vs-405 probe); correct report §7/§10/§11 and owner decision #1's premise;
   fix `LlamaAudioRoutes` + contract message (F-G1, F-C3). Blocks any
   runtime-replacement decision. Nothing else is safe to decide on the stub scan.
2. **Close the selection-enforcement gap**: make auto-ASR consult `Selectable()`
   (or wire `AutoSelect`), run explicit TTS through `capability.Select`, and either
   wire the `Resolver`/`Fallbacks` into production or downgrade report §3/§5 to
   declaration-plus-reporting (F-C1, F-C2, F-C4).
3. **Startup application of saved engines** or a precise documented limitation
   (F-T3-1).
4. **Report inventory**: list manager/factory/resolver as new, enumerate all new
   tests, account the README edit (documentation accuracy).
5. Remove dead `ASRWhisperMed/Large/Fallback` constants or justify retention.

## 9. Verdict: PASS WITH CORRECTIONS

The implementation is fail-closed at every admission boundary, the contract reports
limitations honestly (including what was not run), the three debt fixes are real,
and no regression was introduced by the tranche. It is safe to build upon **after**
corrections 1–3: the native-gate binary evidence is invalid (wrong file scanned,
false absence claimed in shipped code), and the "single source of truth" governs
explicit selection and reporting but not automatic selection. E4/E5 remain
unclaimed; the bundled runtime's audio capability is now unknown-leaning-present,
with absent assets still independently blocking admission.
