# Adversarial Review — Capability Contract and No-Overwrite Publication Correction

**Date:** 2026-11-06
**Reviewer:** Muse (independent adversarial pass)
**Subject:** `docs/review/CORRECTIVE-howl-tts-helper-capabilities-no-overwrite-2026-11-06.md`
**Subject SHA-256 (verified):** `c9f498f61e0444ab9a66b347f466da7ac456e5afb0a30534621ab453de83857f`
**Companion:** `docs/review/IMPLEMENTATION-howl-tts-helper-fake-process-stage2-2026-11-06.md`
**Companion SHA-256 (verified):** `ec85b2b2c3760542699bcd406f1f3b6b3eb2cdc0165b9d6a3b21977a06741227`
**Prior review:** `docs/review/AUDIT-howl-tts-helper-stage2-adversarial-2026-11-06.md` (S2-1…S2-7)
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` §17

**Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. Scope limited to the S2-1/S2-2 correction and
the S2-3/S2-5 evidence corrections. The accepted real-process milestone was not reopened. Stage 3
remains unauthorized. No production file was modified; two temporary probes were run and deleted.
Nothing was committed or pushed.

---

## 1. Evidence verification

| Reported claim | Independent result |
|---|---|
| Correction SHA-256 | **Exact match** |
| Updated Stage 2 report SHA-256 | **Exact match** |
| 11 changed-file hashes (§7) | **All match**, including the new `publish.go` |
| Untouched files still byte-identical | **Confirmed** — `supervisor.go` (`3dd467d3…`), `output.go` (`c50fce52…`), `discovery.go` (`af785db5…`), both path-check files, `internal/proc/job_windows.go` (`1edef75a…`), `wav.go` (`87106725…`), `main.go` (`50b61853…`), both `interop_process_*` files |
| `go test` fakehelper + helper → ok | **Reproduces** (1.133 s / 5.367 s) |
| `go test -race` → ok | **Reproduces** (2.115 s / 6.996 s), no races |
| `go vet` → exit 0 | **Reproduces** |
| Full suite final run → exit 0 | **Reproduces** — `internal/audio` 4.862 s, `image/scanner` 230.659 s, exit 0 |
| "Eleven of the twelve `TestStage2…` launch a real process" | **Confirmed** — 12 functions; `TestStage2MissingHelperIsNonfatal` launches none |
| S2-3 wording corrected | **Confirmed** at report lines 217 and 318; no observed-ordering claim remains |
| Golden fixture updated | **Confirmed** — five service IDs plus provider `howl.tts.synthesize`; no `protocol.hello` / `synthesis.file` |
| No supervisor/ownership change | **Confirmed** by hash |

**Chronology, preserved rather than tidied.** The full-suite record now spans three independent runs:

1. Stage 2 era, implementer's first run: failed in `internal/audio` (`hits 5 misses 2`, want `0/2`).
2. My Stage 2 audit run: **failed, same test, same signature**.
3. Correction run: passed. **My correction-era run: passed.**

Two failures and two passes across four runs, with the failure reproducible only under full-suite load
and never in isolation or in the focused `-count=3` form. That is a load-sensitive pre-existing flake,
not a Stage 2 or correction defect; it also means the suite is **not deterministically green** and must
not be recorded as such. `pipeline_cache_test.go` mtime (07:56) still predates every Stage 2 and
correction artifact.

---

## 2. Gate audit (the S2-1 correction)

### 2.1 Does every operation pass through a gate before a frame can be written?

Yes. The parent has exactly **two** frame-emission sites — `client.go:409` (`call`/`callOn`, used for
`hello`, `health`, `catalog`, `load`, `unload`, `cancel`, `shutdown`) and `client.go:744`
(`synthesize`). Both sit behind the gates below, and no `Supervisor` method encodes or writes a frame,
so there is no second path to the wire.

| Operation | Service gate | Provider gate | Frame-not-sent assertion |
|---|---|---|---|
| `catalog` | `howl.tts.catalog` | — | yes (`TestClientRejectsMissingServiceCapability`) |
| `load` | `howl.tts.load` | — | gate read-verified |
| `unload` | `howl.tts.unload` | — | gate read-verified (see C-2) |
| `synthesize` | `howl.tts.synthesize` | `howl.tts.synthesize` | yes |
| `synthesize`, `speed != 1.0` | — | `howl.tts.speed` | yes |
| `synthesize`, non-empty `language` | — | `howl.tts.language` | yes |
| `cancel` | `howl.tts.cancel` | — | gate read-verified |
| `health`, `shutdown`, `hello` | intentionally ungated | — | correct — no registry ID exists for them |

Every optional request feature is covered: `synthesisParams` carries only `speed` and `language` as
optional capabilities, and there is no SSML field or streaming method, so `howl.tts.ssml` and
`howl.tts.streaming` have no request surface to gate.

### 2.2 Missing / unknown / duplicate / wrong version

| Case | Behaviour | Evidence |
|---|---|---|
| Unknown ID (service or provider) | `protocol_violation` | `capabilityVersions` membership check; tested with `protocol.hello` and `synthesis.file` |
| Wrong version | `protocol_violation` | tested (`Version = 2`) |
| Duplicate ID | `protocol_violation` | implemented (`seen` map); **not tested** — C-4 |
| Missing mandatory service ID | `protocol_violation` + bounded `detail["missing"]` | tested twice |
| Absent provider feature | `unsupported_capability` at request time | tested for synthesis/speed/language |

`unsupported_capability` maps to `ErrCodeSynthesisFailed`, matching the contract's registry row.

### 2.3 Fake-helper side

Complete and symmetric: `cancel` is gated explicitly (`server.go:134`), `catalog`/`load`/`unload`/
`synthesize` through `serviceCapabilityForOperation`, and synthesis additionally requires provider
synthesis, speed, and language. `partial-capabilities` keeps all five service IDs while dropping the
provider features, so an `unsupported_capability` response is reachable without a non-conformant
handshake — the right construction.

### 2.4 Golden fixture

`hello.response.json` now advertises exactly the five service IDs and the registered provider
synthesis ID. Critically, **`TestGoldenFixtures` runs `validateHello` on the fixture**
(`golden_test.go:54`), so the artifact that drifted last time is now subject to the same validator as
the wire. A future non-registry ID fails the suite rather than passing quietly.

---

## 3. Publication audit (the S2-2 correction)

`publishNoReplace`: `os.Link(part, final)` → `os.Remove(part)`; on unlink failure, remove the new final
and return the error.

- **No-overwrite is now enforced by the primitive, not by a check.** `CreateHardLink` fails with
  `ERROR_ALREADY_EXISTS`; there is no replace flag to abuse. The check-then-rename window is gone.
- **Existing file byte-identical:** asserted — `TestPublishNoReplaceRefusesExistingFinal` writes
  `existing`, attempts publication, and re-reads it unchanged; then removes it and confirms a second
  publication succeeds and leaves no part file.
- **Part-file creation is still `O_EXCL`** at write time, so name selection cannot clobber.
- **Publish-failure cleanup:** the caller removes the part on error (`server.go:574`), so a refused
  publication leaks nothing.
- **Cleanup-failure branch is sound but untested.** Forcing `os.Remove(part)` to fail is not practical
  here; the branch removes the final and reports failure, leaving no half-published name — the correct
  conservative behaviour.

One caveat on the *guidance*, not the code, in C-3.

---

## 4. Findings

### C-1 (MODERATE, non-blocking — contract text vs. implemented policy) The mandatory five-service rule is stricter than §17.6.1

The parent now **rejects the handshake** (`protocol_violation`) when any of `catalog`, `load`, `unload`,
`synthesize`, `cancel` is absent. The frozen contract does not state a mandatory minimum. §17.6.1 says
`hello` advertises service capabilities, and that a false capability "is not silently omitted"; §17.6.3
prescribes the remedy for an unavailable capability as a per-request `unsupported_capability`, and the
`Capability` DTO has no `supported` field, so absence is the only encoding of "false".

Consequence: a helper that honestly cannot `unload` (or `cancel`) is now rejected outright, although
the contract's own matrix would degrade it per request — and the client already has a kill-based
fallback for an uncancellable synthesis. The correction thus forecloses a contract-permitted
degradation path.

This is a **policy addition, not a bug**: fail-fast is defensible and safer, and it narrows S2-1's blast
radius. But it is not derivable from the frozen text, so it must not live only in code. Either state the
mandatory minimum in §17.6.1, or relax the handshake to warn-and-degrade and let the per-request gates
carry the load. **Recommendation: add the rule to the contract before Stage 3**; the implementation can
stand as-is.

### C-2 (MINOR — behaviour regression) `Unload` and `Cancel` gate before the lifecycle, so a never-started client reports a false capability diagnosis

`Catalog`, `Load`, and `Synthesize` call `c.Start(ctx)` before their gate. `Unload` and `Cancel` do
not, and `requireServiceCapability` reads the zero-value `c.hello`.

Verified by probe (temporary test, since deleted):

```text
UNLOAD_NOT_STARTED code=unsupported_capability msg="TTS helper does not advertise howl.tts.unload" framesSent=0
CANCEL_NOT_STARTED code=unsupported_capability msg="TTS helper does not advertise howl.tts.cancel" framesSent=0
```

No helper was ever contacted, yet the error asserts a fact about one. Before the correction these calls
would have returned the discovery/start failure, which was accurate. Not reachable in production today
(nothing wires the client yet), but it is a wrong diagnosis a caller can act on and it will be the first
thing a wired caller hits. Fix: call `Start` before the gate as the other operations do, or return a
"not connected" failure when `c.hello` has not been established.

### C-3 (INFORMATIONAL — guidance over-broad) Hard links are an NTFS/same-volume pattern

The mechanism is correct for the fake helper: part and final share a directory, so the link is
same-volume by construction. The report then recommends the pattern to real helpers. Hard links are
unavailable on FAT/exFAT and on many SMB shares, so a real helper adopting it would turn a legitimate
output path on such a volume into `output_failed`. If a real helper needs the guarantee, prefer a
`MoveFileExW` without `REPLACE_EXISTING` (fails when the destination exists) or `CREATE_NEW` on the
final name. Worth a sentence in §4 so the guidance is not inherited blindly.

### C-4 (INFORMATIONAL) Duplicate-capability rejection is implemented but untested

`validateCapabilities` rejects duplicates via a `seen` map; `TestValidateHelloRejectsInvalidCapabilitySets`
covers unknown service, unknown provider, wrong version, and missing — not duplicate. The code is
correct by inspection; the test matrix is one row short of the challenge list.

### C-5 (INFORMATIONAL) `hasCapability` ignores version

Gates match on ID only. This is safe because `validateHello` rejects any non-v1 descriptor before
`c.hello` is stored — but in-package tests mutate `c.hello` directly, so the invariant is upheld by
convention rather than by construction. No action needed at v1; note it for any future protocol minor.

### C-6 (INFORMATIONAL, confirmed good) Three copies of the nine IDs, with the right one now enforced

The IDs exist as constants plus `capabilityVersions` plus `mandatoryServiceCapabilities`, again in the
fake helper's constants, and again as literals in the golden fixture. Duplication between the parent and
the fake is deliberate and desirable — two independent implementations are what make parity meaningful —
and the fixture is now validated rather than merely decoded, which is what closes the drift hole that
produced S2-1.

### Pre-existing contract inconsistency (not this correction)

§17.6.2 says an unadvertised extension is rejected with `unsupported_extension`; the error table around
line 415 says `unsupported_capability`. Both implementations use `unsupported_extension`, matching the
normative section. Worth reconciling in the contract, not in code.

---

## 5. What survives challenge

- The three registries (parent, fake helper, golden fixture) now agree on the same nine IDs and v1
  versions, and the fixture cannot drift silently.
- Every operation and every optional feature is gated above the only two frame-emission sites; there is
  no bypass through the supervisor, which emits no frames at all.
- The no-overwrite guarantee is atomic rather than check-based, with byte-identical preservation and no
  leak on refusal.
- The S2-3 and S2-5 corrections are accurate and match the artifacts.
- The change surface is exactly the eleven declared files; ownership, discovery, output validation, path
  checks, and process startup are untouched and hash-identical.
- Tests, race, vet, and a full-suite run reproduce.

---

## 6. Disposition

| Item | Assessment |
|---|---|
| S2-1 capability registry / membership / gates | **Resolved** — with C-1 (contract text) outstanding |
| S2-2 no-overwrite publication | **Resolved** — atomic, tested, leak-free |
| S2-3 ordering wording | **Resolved** |
| S2-5 test count | **Resolved** |
| Evidence and chronology | **Verified**, flake preserved, suite not deterministically green |
| New defects introduced | **C-2 only**, minor and latent |

**Verdict: no blocking discrepancy. The correction gate can close.**

C-1 is a documentation gap in the frozen contract, not a code fault, and the implemented behaviour is the
safer of the two options; C-2 is a latent misdiagnosis with no production caller. Both should be settled
before Stage 3 wiring begins, and neither should hold the correction gate.

**Still not authorized:** Stage 3, native Kokoro/G2P, ONNX Runtime, provider registration, or any
pipeline/character/cache/UI integration. The Stage 2 deferrals (restart budget, backoff, idle close,
F-4, F-5) and the accepted evidence limits stand unchanged, and no §17.8 provisional value may be cited
as measured.

---

## 7. Boundaries of this review

**Performed:** verification of 13 hashes (2 reports, 11 changed files) plus 10 unchanged-file hashes;
full read of `publish.go`, the capability tables and gates in `protocol.go`/`client.go`, the fake
helper's dispatch and capability advertisement, and every new test; contract cross-check against
§§17.4.1, 17.6.1–17.6.3, 17.9; two temporary probes (gate-before-lifecycle classification, deleted;
the earlier fixture parity probe belongs to the Stage 2 audit); focused, race, vet, and full-suite runs.

**Not performed:** no fix, no permanent test, no code or fixture change, no commit, push, staging,
branch change, or deployment.
