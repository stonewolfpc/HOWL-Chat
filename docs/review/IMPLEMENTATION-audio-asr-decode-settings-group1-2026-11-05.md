# Group 1 — ASR Whisper Decode Settings: Implementation and Evidence

**Date:** 2026-11-05  
**Scope:** `asrLanguage`, `asrThreads`, `asrBeamSize`, `asrBestOf`  
**Status:** Implemented; independent review requested. No commit or push was made.

## Result

The four settings now reach the final `whisper-cli` invocation when, and only when, an explicit Whisper engine is selected. Automatic engine selection and every non-Whisper engine retain their previous decode behavior and are reported as `unavailable` or `unsupported`, never as applied.

The Group 0 34-field desired-settings contract, microphone gate, startup no-auto-load behavior, retry path, ASR lifecycle, and independent TTS path were preserved. The legacy four-field `applied` projection remains a compatibility view and is not authoritative.

## Files Changed

### Operational path

- `internal/audio/settings_contract.go` — four decode keys became operational; explicit-Whisper language/profile validation; engine-specific status; invalid decode values block activation only for explicit Whisper; active decode values are retained in `Applied` before a replacement is classified.
- `internal/audio/manager.go` — `audio.Settings` carries `ASRLanguage`, `ASRThreads`, `ASRBeamSize`, and `ASRBestOf` through the active generation.
- `internal/audio/types/interfaces.go` — `PipelineConfig` carries the same four values.
- `internal/audio/manager_factory.go` — maps the operational settings into `PipelineConfig`.
- `internal/audio/pipeline.go` — forwards configured decode values to explicit Whisper and preserves the previous `auto`/`6` values for automatic and non-Whisper selections.
- `app-audio.go` — passes the desired contract to `HasBackendValidationIssue`; updated the compatibility comment. No Group 1 change was made to microphone probing or the activation gate.
- `frontend/dist/js/audio-settings.js` — aggregate requires the four decode keys only for explicit Whisper, reports `unavailable`, and includes the values in the applied-count calculation.

### Tests

- `internal/audio/asr/whisper/decode_args_test.go` (new) — captures the final CLI arguments.
- `internal/audio/pipeline_decode_test.go` (new) — proves pipeline forwarding and non-Whisper gating, including local llama legacy values.
- `internal/audio/settings_contract_test.go` — updated counts and engine-specific status coverage.
- `internal/audio/capability_seam_test.go` — adds the local llama test provider to the selectable test catalog for the gating test.
- `app_audio_contract_test.go` — updated counts/applied map and added invalid-thread, profile-language, and unsupported-engine cases.
- `app_audio_ownership_test.go` — test helper preserves the new settings when constructing desired state.
- `tests/e2e/specs/iframe-bridge.spec.js`
- `tests/e2e/specs/audio-settings-lifecycle.spec.js`
- `tests/e2e/specs/activation-lock.spec.js`

The three Playwright fixtures now produce the four decode keys only for explicit Whisper and assert 8 backend settings applied / 24 unsupported for that fixture, while the non-Whisper fixture remains at 28 unsupported.

`internal/audio/settings_contract.go` and `app_audio_contract_test.go` were already untracked Group 0 files in this worktree. They were modified, not created, by Group 1.

## Setting-by-Setting Data Flow

| Setting | Persisted key and type | Validation / default | Renderer to Go | Provider field | Final Whisper argument | Behavior controlled | Supported when |
|---|---|---|---|---|---|---|---|
| Spoken language | `asrLanguage`, string | allowed language list; profile-language check for explicit Whisper; default `auto` | `DesiredSettings` → `audio.Settings.ASRLanguage` → `PipelineConfig.ASRLanguage` → `types.ASRConfig.Language` | `ASRConfig.Language` | `--language <value>` only when not `auto`/empty | Selects the decoder language; omission retains whisper.cpp automatic detection | Explicit local Whisper CLI only |
| CPU threads | `asrThreads`, int | `1..16`; default `4` | same path → `PipelineConfig.ASRThreads` | `ASRConfig.NumThreads` | `--threads <n>` | Worker-thread count used by whisper.cpp | Explicit local Whisper CLI only |
| Beam size | `asrBeamSize`, int | `1..10`; default `5` | same path → `PipelineConfig.ASRBeamSize` | `ASRConfig.BeamSize` | `--beam-size <n>` | Beam-search width used by whisper.cpp | Explicit local Whisper CLI only |
| Best-of | `asrBestOf`, int | `1..10`; default `5` | same path → `PipelineConfig.ASRBestOf` | `ASRConfig.BestOf` | `--best-of <n>` | Candidate count considered by whisper.cpp decoding | Explicit local Whisper CLI only |

The adapter already emitted these flags from `ASRConfig`; the missing link was the pipeline hardcoding `Language: "auto"` and `NumThreads: 6` and never supplying beam/best-of. `asrDecodeConfigFor` now uses the configured values only when `isWhisperASREngine` is true. That predicate requires a known local profile whose runtime is `types.RuntimeWhisperCLI`, so `auto` is deliberately excluded because it does not identify the actual decoder.

For automatic and non-Whisper selections, `asrDecodeConfigFor` returns the prior pipeline values: `Language: "auto"`, `NumThreads: 6`, and zero beam/best-of. This prevents the new Whisper settings from silently changing llama-server or Qwen behavior. Local llama still receives its existing 6-thread value; its language field is a legacy passthrough, not an application of the new `asrLanguage` setting. Its user-facing decode settings remain unsupported.

### Intentional behavior change for explicit Whisper

An explicit Whisper selection now uses the established contract defaults (`4` threads, beam size `5`, best-of `5`, language `auto`) instead of the pipeline's old hardcoded `6` threads and absent beam/best-of arguments. This is the wiring requested by Group 1, not a new settings-contract default. Automatic and non-Whisper selections retain the old `6`-thread pipeline value.

Source references in the current worktree:

- `internal/audio/settings_contract.go:120` — `DesiredSettings.Backend()` maps the four operational values into `audio.Settings`.
- `internal/audio/settings_contract.go:136` — `DesiredSettings.Validate()` owns defaults/ranges and explicit-Whisper profile-language validation.
- `internal/audio/settings_contract.go:276` — `isWhisperASREngine` requires a local profile with `RuntimeWhisperCLI`.
- `internal/audio/settings_contract.go:478` — `setASRDecode` derives the engine-specific status.
- `internal/audio/pipeline.go:382` — `asrDecodeConfigFor` gates the values before provider initialization.
- `internal/audio/asr/whisper/recognizer.go:300` — `runWhisperCLI` emits `--language`, `--beam-size`, `--best-of`, and `--threads`.

## Status Semantics

`BuildSettingsReport` now classifies each decode key as follows:

- ASR disabled: `unavailable` with “ASR is disabled”.
- ASR enabled but engine `auto`/empty: `unavailable` with a “choose an explicit Whisper engine” reason.
- ASR enabled with an explicit non-Whisper engine: `unsupported` with an engine-specific reason; the key never enters `Applied`.
- Explicit Whisper, active generation matches, provider ready: `applied` and the value enters `Applied`.
- Explicit Whisper, no matching active generation, in-progress replacement, or apply failure: `pending` or `rejected`; the previous effective decode values remain in `Applied` when an explicit Whisper generation is active.
- Invalid requested value: `invalid`; activation is blocked before the manager/factory is called for an explicit Whisper request.
- Invalid values for an unsupported engine remain reportable but do not block the core capability fields.

The aggregate UI requires the four decode keys only when the selected engine is `whisper-tiny`, `whisper-base`, or `whisper-small`. An explicit Whisper fixture therefore reports 8 of 8 backend settings applied; the non-Whisper fixture still reports the decode settings as unsupported. The legacy `GetAudioStatus()["applied"]` map remains limited to the four capability fields.

## Verification Executed

### Go tests

Command:

```text
GOTOOLCHAIN=go1.23.8 go test -tags "desktop,production" ./internal/audio ./internal/audio/asr/whisper .
```

Result: **PASS** — all three packages reported `ok`:

- `howl-chat/internal/audio`
- `howl-chat/internal/audio/asr/whisper`
- `howl-chat` (root)

The tests cover:

- `TestWhisperDecodeSettingsReachFinalCLI` — exact `--language`, `--threads`, `--beam-size`, and `--best-of` values at the final captured invocation.
- Automatic language — no `--language` argument is emitted for `auto`.
- `TestWhisperThreadChangeChangesInvocation` — changing configured threads changes `--threads`.
- Beam size and best-of are asserted independently with different values.
- `TestPipelineForwardsWhisperDecodeConfig` — explicit Whisper receives all four values.
- `TestASRDecodeConfigIsOnlyForExplicitWhisper` — auto and non-Whisper selections receive the legacy values, not configured decode values.
- `TestPipelineKeepsLegacyASRDefaultsForLocalLlama` — local llama receives `auto`/`6`/no beam/no best-of.
- `TestSettingsReportWhisperDecodeStatusIsEngineSpecific` — Qwen gets `unsupported`; auto gets `unavailable`; neither enters `Applied`.
- `TestInvalidWhisperDecodeValueFailsBeforeActivation` — invalid threads are rejected before factory activation and the active value `8` remains in `Applied`.
- `TestWhisperLanguageOutsideProfileFailsBeforeActivation` — `whisper-tiny` + `nl` is rejected with a language-specific validation reason.
- `TestInvalidDecodeValueForUnsupportedEngineDoesNotBlockCoreApply` — an invalid decode value for Qwen does not block the core engine apply and is not claimed as applied.
- Existing failed-apply and microphone/no-auto-load/TTS lifecycle tests remain green.

### Audio frontend regression suite

Command, run from `tests/e2e`:

```text
npx --no-install playwright test --config=playwright.audio.config.js specs/iframe-bridge.spec.js specs/audio-settings-lifecycle.spec.js specs/activation-lock.spec.js
```

Result: **31 passed** in 16.8 seconds.

This includes the iframe bridge handshake, applied-state truthfulness, unsupported/loaded setting reporting, microphone gating, startup restore-without-apply, activation locks, retry behavior, and TTS/lifecycle regression coverage.

### Production build

Command:

```text
GOTOOLCHAIN=go1.23.8 CGO_ENABLED=0 "$HOME/go/bin/wails.exe" build -devtools -s -nopackage -m -nosyncgomod -tags "desktop,production"
```

Result: **PASS** — Wails CLI v2.12.0 built `C:\Users\StoneWolfPC\HOWL_Chat\build\bin\howl-chat.exe` with tags `[desktop,production]`.

### Static checks

- `node --check frontend/dist/js/audio-settings.js` — pass.
- `gofmt -l` over all changed Go files — no output.
- Scoped `git diff --check` over the tracked Group 1 paths — pass.
- Whole-worktree `git diff --check` — not clean because of pre-existing unrelated trailing whitespace in `IMAGE_PAGE_CHANGE_LEDGER.md` and generated `frontend/wailsjs/go/models.ts`. Those files were already modified in the worktree and are outside Group 1.

## Evidence Grades

- **Observed:** test exit statuses, exact captured Whisper arguments, Playwright results, Wails build receipt.
- **Source-derived:** the validation ranges, defaults, engine predicate, status transitions, and the fact that `runWhisperCLI` emits these four flags.
- **Calculated:** 8 applied / 24 unsupported / 2 renderer for the explicit-Whisper fixture with both capabilities enabled; 7 applied with TTS disabled.
- **Inferred:** the semantic effect of the four whisper.cpp flags on decoding quality or latency was not measured.
- **Unestablished:** a real Whisper model decode, packaged WebView2 runtime smoke, and restart behavior under the production binary were not exercised in this slice.
- **Superseded:** the Group 0 statement that these four settings are unsupported is superseded only for explicit local Whisper engines. It remains true for automatic selection and all non-Whisper engines.

## Limitations and Review Notes

1. **Arguments were captured through the test seam.** The tests prove the exact CLI arguments, not the acoustic output of a real Whisper model. A live `whisper-cli` decode with a model is still needed to confirm quality and latency effects.
2. **Production binary was built, not launched.** The Wails build validates compilation and binding generation; no interactive WebView2 smoke test was performed.
3. **Automatic Whisper selection is intentionally excluded.** `auto` cannot prove which provider will decode, so the values are reported `unavailable` until the user chooses an explicit Whisper engine. This is deliberate truthfulness, not an applied-but-hidden behavior.
4. **llama/Qwen retain legacy behavior.** Local llama still receives its prior 6-thread value and legacy language field; Qwen receives the prior values. The four user settings remain unsupported there and do not alter their behavior.
5. **No unrelated subsystem edits were made for Group 1.** Microphone probing, capture devices, routing, VAD, push-to-talk, gain, TTS, cache, and unrelated settings were not changed. The Group 0 worktree already contained unrelated modifications and untracked artifacts; those remain untouched. The Wails production build regenerated frontend bindings as part of verification; generated binding files were already modified in the worktree and are not Group 1 source edits.
6. **No commit or push.** All changes are in the worktree for independent review.

## Scoped Diff Review

Group 1 touched the operational path, focused tests, and the three audio frontend specs listed above. The worktree was not clean before Group 1: `app-audio.go`, `frontend/dist/js/audio-settings.js`, and the Playwright fixtures already carried accepted Group 0 changes, and `internal/audio/settings_contract.go` / `app_audio_contract_test.go` were untracked Group 0 files. A repository-wide diff therefore cannot isolate Group 1 by itself; the review above is path-scoped.

No commit, push, or production publication was performed.
