# Corrective — Capability Contract and No-Overwrite Publication

**Date:** 2026-11-06  
**Scope:** S2-1, S2-2, and the S2-3/S2-5 evidence corrections only  
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` §17  
**Review basis:** `docs/review/AUDIT-howl-tts-helper-stage2-adversarial-2026-11-06.md`  
**Status:** CORRECTION IMPLEMENTED FOR REVIEW. STAGE 2 REMAINS ACCEPTED FOR ITS REAL-PROCESS MILESTONE. STAGE 3 IS NOT AUTHORIZED. NO COMMIT OR PUSH.

---

## 1. Authorization boundary

The owner accepted Stage 2 for its real-process interoperability milestone and authorized a narrow correction before any native-provider work:

- resolve S2-1 while the fake helper is the only implementation;
- include the Windows no-overwrite discrepancy (S2-2);
- correct the S2-3 ordering wording and the S2-5 test count;
- do not reopen the entire real-process implementation.

This correction does not add a native provider, model, voice, ONNX Runtime, CGO, Python, production pipeline, character, cache, UI, packaging, download, or Stage 3 work. No commit or push occurred.

---

## 2. Findings and disposition

| Finding | Correction |
|---|---|
| **S2-1** capability negotiation | **Resolved in code and tests.** Frozen nine-ID registry, mandatory service capabilities, exact version check, parent-side operation/feature gates, helper-side gates, and `unsupported_capability` on both sides. |
| **S2-2** Windows rename replacement | **Resolved.** Final publication no longer uses `os.Rename`; `publishNoReplace` creates a hard link and refuses an existing destination. |
| **S2-3** ordering evidence | **Corrected in the Stage 2 report.** Tests assert end state, not observed kill-before-delete ordering. Ordering is structural plus the Stage 1 in-memory proof, not observed at the OS-child boundary. |
| **S2-5** test count | **Corrected.** Eleven `TestStage2...` tests launch a process; `TestStage2MissingHelperIsNonfatal` launches none. |

---

## 3. Capability-contract correction

### 3.1 Frozen registry and mandatory set

The main-side codec and the fake helper now use the §17.6.1 identifiers:

| Capability | Version |
|---|---:|
| `howl.tts.catalog` | 1 |
| `howl.tts.load` | 1 |
| `howl.tts.unload` | 1 |
| `howl.tts.synthesize` | 1 |
| `howl.tts.cancel` | 1 |
| `howl.tts.speed` | 1 |
| `howl.tts.language` | 1 |
| `howl.tts.ssml` | 1 |
| `howl.tts.streaming` | 1 |

The five service operations are mandatory in `hello`: `catalog`, `load`, `unload`, `synthesize`, and `cancel`. The parent rejects a hello that omits any of them. Unknown IDs, duplicate IDs, and any version other than 1 are `protocol_violation`; the missing-capability failure carries a bounded `detail["missing"]` listing the absent mandatory IDs.

This does not invent a `supported` JSON field: the frozen `Capability` DTO has only `id` and `version`. Presence means v1 support, and optional IDs remain valid registry members.

### 3.2 Parent-side enforcement

`Client` now gates operations against the validated hello before sending a frame:

| Request | Required capability |
|---|---|
| `catalog` | service `howl.tts.catalog` |
| `load` | service `howl.tts.load` |
| `unload` | service `howl.tts.unload` |
| `synthesize` | service `howl.tts.synthesize` and provider `howl.tts.synthesize` |
| `synthesize` with `speed != 1.0` | provider `howl.tts.speed` |
| `synthesize` with a non-empty language | provider `howl.tts.language` |
| `cancel` | service `howl.tts.cancel` |

A failed gate returns the frozen `unsupported_capability` code. The tests assert that no operation frame is sent when a gate fails. `health` and `shutdown` are not gated because they are control operations and have no capability ID in the registry.

### 3.3 Helper-side enforcement

The fake helper applies the same service gate before executing `catalog`, `load`, `unload`, `synthesize`, or `cancel`. Synthesis additionally checks its provider capabilities for synthesis, speed, and language before validating the optional request feature.

The `partial-capabilities` test mode advertises only provider synthesis, keeping all five service capabilities valid. This permits a real `unsupported_capability` response for a speed request without making the handshake itself non-conformant.

### 3.4 Golden fixture

`internal/audio/tts/helper/testdata/golden/hello.response.json` now advertises the five service IDs and the registered provider synthesis ID. It no longer advertises the non-registry `protocol.hello` or `synthesis.file` names. `TestGoldenFixtures` therefore validates a contract-native hello.

---

## 4. No-overwrite publication

`internal/audio/tts/fakehelper/publish.go` replaces the Windows-replacing `os.Rename` publication:

1. `os.Link(part, final)` creates the final directory entry atomically and fails if `final` already exists.
2. `os.Remove(part)` removes the temporary name.
3. If unlinking the part fails, the newly created final entry is removed and publication reports failure.

The part file and final file are in the same directory, so the hard link is on the same volume. The previous `Lstat` pre-check remains as an early diagnostic, but the no-overwrite guarantee is now enforced by the publication primitive itself. A direct test creates an existing final, attempts publication, and confirms the final contents are unchanged; a second publication succeeds after the final is removed.

This guarantee applies to the fake helper's output path and is the pattern a real helper must follow. It does not claim cross-volume hard-link support.

---

## 5. Evidence-wording corrections

### S2-3

`TestStage2ForcedCancellationKillsAndCleans` and `TestStage2TimeoutKillsAndCleans` assert the end state: the PID is gone and no final or `.part` remains. They do not distinguish terminate → observe exit → delete from delete-then-kill by observation at the OS-child boundary. The ordering is structural in `Client.cleanupRejected` and `Supervisor.Terminate`; the Stage 1 in-memory test is the ordering proof. The Stage 2 report now says this explicitly and no longer implies real-process ordering was observed.

### S2-5

There are twelve `TestStage2...` functions. Eleven launch a real child; `TestStage2MissingHelperIsNonfatal` deliberately has no helper to launch. The Stage 2 report now says “eleven … launch a real process; … launches none” rather than “twelve real-process tests.”

### Full-suite chronology

The correct record is:

- Stage 2 package tests: pass.
- Real-process matrix: reproduced.
- Race, vet, and no-CGO build: pass.
- Repository-wide suite: not deterministically green because of the pre-existing `internal/audio` cache test.

The existing failure remains `TestConcurrentDifferentConfigurationsRemainIsolated` with `hits 5 misses 2`, want `0/2`; isolated and audio-tree runs pass. It is a load-sensitive pre-existing test, not a Stage 2 or capability-correction defect.

---

## 6. Verification

Environment: `go1.27.3 windows/amd64`; module `howl-chat`.

```text
go test ./internal/audio/tts/fakehelper ./internal/audio/tts/helper -count=1
ok  howl-chat/internal/audio/tts/fakehelper  1.378s
ok  howl-chat/internal/audio/tts/helper      4.322s
exit 0
```

```text
go test -race ./internal/audio/tts/fakehelper ./internal/audio/tts/helper -count=1
ok  howl-chat/internal/audio/tts/fakehelper  2.566s
ok  howl-chat/internal/audio/tts/helper      6.812s
exit 0
```

```text
go vet ./internal/audio/tts/fakehelper ./internal/audio/tts/helper
exit 0

CGO_ENABLED=0 go build ./...
exit 0

go test ./internal/audio/... -count=1
exit 0
```

The full audio tree passed, including `internal/audio` at 3.443s.

New focused tests:

- `TestValidateHelloRejectsInvalidCapabilitySets` (unregistered service/provider, wrong version, missing service);
- `TestValidateHelloNamesMissingCoreCapabilities`;
- `TestClientRejectsMissingServiceCapability`;
- `TestClientRejectsUnadvertisedProviderFeatures` (synthesis, speed, language);
- `TestRunAdvertisesFrozenCapabilities`;
- `TestHandleSynthesizeRejectsUnadvertisedOptionalCapability`;
- `TestPublishNoReplaceRefusesExistingFinal`.

The real-process handshake test now asserts exactly five service capabilities and their v1 versions, replacing the previous `>= 5` count.

Repository-wide correction run:

```text
go test ./... -count=1
exit 0
ok  howl-chat  12.633s
ok  howl-chat/internal/audio  5.729s
ok  howl-chat/internal/audio/tts/fakehelper  0.934s
ok  howl-chat/internal/audio/tts/helper  6.651s
ok  howl-chat/internal/image/scanner  232.315s
```

This one run passed, but it does not make the repository suite deterministically green. The earlier Stage 2 full run and the independent audit both reproduced the same `internal/audio` cache-test failure; isolated, focused, and audio-tree runs pass. The correction did not touch that test or the cache implementation.

---

## 7. Changed files and corrected hashes

| File | SHA-256 |
|---|---|
| `internal/audio/tts/helper/protocol.go` | `3a0b1cd57655b6731641b2e7f6a5f0dad42dd45504fd503c2bc1c0bb1ad30204` |
| `internal/audio/tts/helper/client.go` | `b33c9f5151ba36f4b1363d19b0c8875a6211295077f40d2e413665866577160a` |
| `internal/audio/tts/helper/protocol_test.go` | `2c3b8e1e2470f83f23f92dda3849fd07fd7d55b93e93bfe875a913fff1dedb8f` |
| `internal/audio/tts/helper/client_test.go` | `ff347711e76b78d0b1c8a3c8b3aa378acb535356b8a7d7745211737d00ce0a6d` |
| `internal/audio/tts/helper/fake_test.go` | `c5d2a857fb2d70c23e3e7d6d41c47438eb96443f1da01003cf505e55a5b3a3a1` |
| `internal/audio/tts/helper/interop_test.go` | `28ed6d11b551fc2aa86c4bec02be4e7311002bb45f52e654578d75b3d8cf203d` |
| `internal/audio/tts/helper/testdata/golden/hello.response.json` | `e0d5ded1065a8e7c5a3a870d138d7c74bb05012237ac33f0b5b72951c9414709` |
| `internal/audio/tts/fakehelper/protocol.go` | `ece9097423381653a999a126fe5e4bcbff084e4d80a63a9fef77890ea8dd79ea` |
| `internal/audio/tts/fakehelper/server.go` | `c4c613306569f1f34e516084ef0b90199e2fdece673433dbb3f62d9045f2767e` |
| `internal/audio/tts/fakehelper/server_test.go` | `1f2234bd2720231f0915a2d6f1b009ba3f093012241c0b6c242632fadf0ca846` |
| `internal/audio/tts/fakehelper/publish.go` | `f6844a5c54c3a8489d1fdeaf0deaf5325db92dedcc67d0fbc0c2cd2c4c7c986d` |

The accepted Stage 1 baseline hashes for `protocol.go` (`8a9f710b…`) and `client.go` (`e2b2e713…`) remain in the Stage 2 report's baseline table. The corrected hashes above supersede them; no supervisor, ownership, discovery, output-validation, or path-check file changed.

---

## 8. Status

- S2-1: resolved for Stage 3 entry.
- S2-2: resolved for the fake-helper publication path.
- S2-3 and S2-5: evidence wording corrected.
- Stage 2: accepted for its real-process interoperability milestone.
- Stage 3: not authorized.
- No commit, push, staging, branch change, or deployment occurred.
