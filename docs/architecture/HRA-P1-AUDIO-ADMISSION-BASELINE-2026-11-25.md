# HRA-P1 — Audio Admission Baseline

## Record `HRA-P1-AUDIO-2026-11-25-01`

| Field | Value |
|---|---|
| Date | **2026-11-25** |
| Type | Baseline/classification record |
| Test | `internal/audio/TestASRStartConsultsAdmissionWhileChatLoaded` |
| Gate impact | **Gate 1 did not cause the failure** |
| Classification | **Pre-existing policy/test mismatch; audio is a required HRA-P1 completion gate** |
| Production changes made for this record | **None** |

---

## 1. Question

Before Gate 1 is declared closed, determine:

1. Did `TestASRStartConsultsAdmissionWhileChatLoaded` fail before Gate 1?
2. Does the test's expected refusal still match the documented admission policy?
3. If the failure is a pre-existing mismatch, record it separately rather than carrying it as an unexplained Gate 1 regression.

This record does not change Gate 1, the test, or the admission policy.

## 2. Baseline method

`git show HEAD:internal/resource/admission.go` is unavailable: the source exists on disk but not in the `HEAD` tree. To establish behavior without altering the worktree, I used a Go overlay:

- copied the current `internal/resource/admission.go`;
- reversed only the Gate 1 changes in the copy:
  - `grant` records `RequiredVRAMMB` again;
  - `blockingModels` excludes `GPUMemoryMB <= 0` again;
  - resident-stage projection does not add new adapter bytes;
  - no `committedVRAMMB` normalization;
- compiled the unchanged `internal/audio` test against that overlay.

The temporary overlay directory is not retained after this record; the commands below show the form that was executed. The overlay is validated by the graduated D-1/D-9 tests:

```text
go test -count=1 -overlay=.hra-baseline/overlay.json -run 'TestGrantReservesProjectedPeak|TestRAMOnlyResidentBlocksRAMAdmission' ./internal/resource
```

Both tests fail under the overlay with their original pre-Gate-1 symptoms:

```text
pending reservation = 0 MB, want the projected peak 512 MB
RAM-only resident was not reported as the blocker: []
```

The overlay therefore reproduces the pre-Gate-1 admission behavior.

## 3. Baseline result

Baseline audio test:

```text
go test -count=1 -overlay=.hra-baseline/overlay.json -run TestASRStartConsultsAdmissionWhileChatLoaded -v ./internal/audio
--- FAIL: TestASRStartConsultsAdmissionWhileChatLoaded
    admission_test.go:40: ASR must not become ready when admission refuses the start
```

Current Gate 1 audio test:

```text
go test -count=1 -run TestASRStartConsultsAdmissionWhileChatLoaded -v ./internal/audio
--- FAIL: TestASRStartConsultsAdmissionWhileChatLoaded
    admission_test.go:40: ASR must not become ready when admission refuses the start
```

The failure is identical before and after Gate 1. Gate 1 did not change the ASR admission decision, the soft-reserve branch, or the ASR profile.

## 4. Why the test's expectation is stale

The test supplies:

- free VRAM: **4500 MB**;
- loaded chat model: **3000 MB VRAM / 4000 MB RAM**;
- current local Qwen ASR profile: **4096 MB VRAM / 3840 MB RAM**.

The authority returns:

```text
admitted=true
required=4096 free=4500 reserve=450 available=4050 soft_allowance=322 soft_used=46
reason=admitted: 4096 MB VRAM required, 4050 MB available after 450 MB reserve;
soft reserve: used 46 MB of the 322 MB allowance,
keeping 404 MB free above the 128 MB hard floor
```

The documented policy is:

- `DefaultReservePercent` is 10%, but it is a target, not an absolute floor;
- `softReserveAllowanceMB` is 500 MB and may be used when it is the difference between running and refusing;
- the hard reserve floor is 128 MB;
- explicit reserves are never softened.

With a 450 MB target reserve, the soft floor is `max(128, 450-500) = 128 MB`, leaving a 322 MB allowance. The request uses 46 MB of that allowance and leaves 404 MB free, above the 128 MB hard floor.

Therefore the current decision **admits** the ASR start. The test expects refusal, but no documented rule says:

- ASR may not use the soft reserve; or
- a resident chat model disables the soft allowance; or
- the 10% reserve is a hard floor for ASR.

The test's expected refusal is not supported by the documented policy at its current fixture values.

## 5. Classification

**PRE-EXISTING POLICY/TEST MISMATCH.**

The test predates or has drifted from the current ASR profile/policy combination. It fails on the pre-Gate-1 overlay exactly as it fails after Gate 1. Gate 1 requires no correction for this failure.

The test's stated intent — the pipeline must consult admission before starting the ASR child — is still valid. Its fixture no longer creates the refusal it assumes.

## 6. Disposition

Gate 1 may be treated as closed for D-1/D-9 because this failure is classified as pre-existing and unrelated. Do not use the failing audio test as evidence against Gate 1.

The test needs a separate, explicit disposition:

1. **Preserve the refusal test intent:** change the fixture to a request the policy actually refuses (for example, a free-memory value that leaves less than the 128 MB hard floor), then assert no child starts. This changes the test fixture, not the policy.
2. **Test the admitted case:** assert that this fixture is admitted and that the pipeline still does not start a child before admission returns. This requires changing the test's expected outcome.
3. **Change the policy:** only with a documented decision that ASR cannot use the soft allowance. That would be a production policy change and is not authorized by this record.

No disposition is applied here. Audio integration is now a required HRA-P1 completion gate (`HRA-P1-IMPLEMENTATION-PLAN-2026-11-25.md` §7.6 and `HRA-P1-RESOURCE-LIFECYCLE-CONTRACT-2026-11-25.md` §12.3). The test remains red and visible; it must not be waived, skipped, renamed as a known failure, or weakened to make the suite green. The three options above remain unauthorized until the intended reserve contract is decided and the audio lifecycle is integrated as a first-class workload.

Until the audio gate closes, HRA-P1 is not complete even if the remaining image/workflow tests pass. The current failure is evidence of an unresolved policy/test mismatch, not proof that ASR bypasses admission.

## 7. TTS baseline and offline-only scope (2026-11-25)

No disposition is applied to the ASR fixture by this section. It records the TTS side of the same gate, because "ASR and TTS participate in the same lifecycle and resource guarantees" cannot be demonstrated by the ASR path alone.

**TTS has never been exercised on this machine.** There is no installed Piper voice (`models/en_US-lessac-medium.onnx` and `bin/piper.exe` are absent, and `models/` contains no `.onnx` file), no Kokoro asset set (`models/kokoro` and `runtime/tts` are absent), and no MeloTTS service on the configured loopback URL. `config/settings.json` records `ttsEnabled: false`, `ttsEngine: "piper"`, `ttsVoice: "en_US-lessac-medium"`. The correct statement is therefore "TTS is unverified here", not "TTS works" and not "TTS is broken".

**Every supported TTS engine is local.** Piper runs as a local executable spawned per request. Kokoro runs local ONNX weights through a local Python bridge. MeloTTS is an HTTP service addressed at `http://127.0.0.1:8888` by default — loopback, not a cloud endpoint. Model acquisition is a separate, manual download step; no synthesis call leaves this machine. The only network-capable audio adapters in the build are the key-gated DashScope Qwen ASR clients, and automatic ASR selection filters to `Local && ApprovedLocal`, so no cloud engine is selected automatically.

**Admission coverage is not the same as lifecycle coverage.** TTS does call the same `resource.Authority` before a provider is constructed (`initTTS`, `loadTTSForRequest`, `ensureTTSVoice`), and `internal/audio/tts_lifecycle_characterization_test.go` now pins that. It does not yet prove that a TTS child is owned, killed and reaped the way the ASR server and the image server are: the Piper and Kokoro providers spawn with `exec.CommandContext` outside job ownership. The ownership table, the named gaps, and the host acceptance protocol that must be run with the user present are recorded in `HRA-P1-IMPLEMENTATION-PLAN-2026-11-25.md` §7.6.1.

## 8. Limitations

- No supported host was run.
- No user-visible workflow was demonstrated.
- The baseline is a source overlay, not a historical checkout, because the source is absent from `HEAD`.
- Gate 2 has not started.
- TTS has no installed assets on this machine, so no TTS synthesis has been characterized end to end.
