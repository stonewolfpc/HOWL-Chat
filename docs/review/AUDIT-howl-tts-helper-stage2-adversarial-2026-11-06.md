# Adversarial Review — Stage 2 Fake Helper / Real-Process Interoperability

**Date:** 2026-11-06
**Reviewer:** Muse (independent adversarial pass)
**Subject:** `docs/review/IMPLEMENTATION-howl-tts-helper-fake-process-stage2-2026-11-06.md`
**Subject SHA-256 (verified):** `adb96b068142f17244eab0b8b86dd0714cc1660543d9c45af8f9b62da612b101`
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` §17, incl. A1/A1-D7
**Baseline claimed:** accepted Stage 1 implementation report + `CORRECTIVE-howl-tts-helper-supervisor-stage1-2026-11-06.md`

**Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. Not an authorization. Stage 3 remains
unauthorized; no native provider, pipeline, cache, character, or UI work is endorsed here.

**Method:** hash verification of all 22 declared artifacts, full read of the new Stage 2 sources and
tests, cross-check of both codecs against the frozen DTO/capability tables, independent re-execution
of every reported command, one full-suite run, and one temporary parity probe (written, run, deleted).
No production file was modified. The probe directory was verified back to its six files.

---

## 1. Evidence verification — reproducible, with one chronology correction

| Reported claim | Independent result |
|---|---|
| Report SHA-256 | **Exact match** |
| 11 baseline hashes (§3) | **All match**, including my own Stage 1 audit (`03ed0a46…`) |
| 10 Stage 2 artifact hashes (§3) | **All match** |
| `go test` fakehelper + helper → ok | **Reproduces** (0.810 s / 4.487 s) |
| 12 `TestStage2…` → PASS | **Reproduces**, one correction (§2.5): 11 launch a process, 1 does not |
| `go test -race` → ok | **Reproduces** (1.824 s / 6.350 s), no races |
| `go vet` → exit 0 | **Reproduces** |
| `CGO_ENABLED=0 go build ./...` → exit 0 | **Reproduces** |
| Root `go test .` → ok | **Reproduces** (`ok howl-chat 12.571s` in my full run) |
| Full suite final retry → exit 0 | **Not reproducible on demand — see §2.6.** My run failed in `internal/audio`, same test, same signature |
| Disclosed flake is pre-existing and unrelated | **Corroborated** by three independent signals (§2.6) |
| No commit/push; new dirs untracked | **Confirmed** (`??` on `cmd/howl-tts-helper-fake/`; helper dir still untracked) |
| No production wiring | **Confirmed** — zero imports of `tts/helper` or `tts/fakehelper` outside their own trees |

The report describes the code on disk. Nothing was overstated about *what was tested*; two
*characterizations* need correcting (§2.5, §2.7).

---

## 2. Findings

### S2-1 (MODERATE — protocol conformance) Capability negotiation is unimplemented, and the golden fixtures encode a non-conformant `hello`

This is the parity question answered with evidence, and it cuts the way the owner suspected.

- `internal/audio/tts/helper/client.go` contains **zero** references to capabilities. The parent stores
  them for status display and never consults them. (The interop test asserts
  `len(status.Capabilities) >= 5`; nothing asserts that a *required* capability is present.)
- `validateCapabilities` (`protocol.go`) checks only the namespace regex, `version >= 1`, and duplicate
  IDs. §17.6.1 is an enumerated registry (`howl.tts.catalog|load|unload|synthesize|cancel|speed|language|ssml|streaming`);
  membership is never checked on either side.
- §17.6.1 is normative: *"A capability with `Supported:false` is false for that provider; it is not
  silently omitted. A request requiring a false capability returns `unsupported_capability`."*
  Neither implementation returns `unsupported_capability` for anything.
- `testdata/golden/hello.response.json` advertises service capability **`protocol.hello`** and provider
  capability **`synthesis.file`** — neither exists in §17.6.1 — and **omits every core service
  capability**. Under the contract's own rule that is a silently-omitted capability set, i.e. the
  normative fixture is itself non-conformant.
- The fake helper, by contrast, emits exactly the five contract service IDs and the two provider IDs.
  The two artifacts therefore **disagree with each other**, not merely with the contract:
  helper = conformant, fixture = not.

**Why this matters:** a helper that simply omits `howl.tts.cancel` or `howl.tts.synthesize` completes
the handshake cleanly, and the adapter will then call operations the helper never claimed to support.
Degradation is tolerable today (`unsupported_operation` → typed failure → kill), but the parent is
asserting a capability set it never validated, which is precisely the drift risk named in the brief.
This becomes material the moment a real provider advertises partial capabilities (Stage 3+).

**Fix direction (not this review):** decide whether capability membership is enforced (registry check
+ `unsupported_capability` gate before `catalog`/`load`/`synthesize`/`cancel`) or explicitly declared
advisory; and correct `hello.response.json` to the §17.6.1 registry either way.

### S2-2 (MINOR — inaccurate claim) "Renames to the final path without overwriting" is not what `os.Rename` does on Windows

`server.go` publishes with `os.Rename(part, final)` after an `Lstat` non-existence check. Go's Windows
`rename` passes `MOVEFILE_REPLACE_EXISTING`, so it **will** replace a file that appears between the
check and the rename. The `.part` side is genuinely exclusive (`O_EXCL`, verified) and the pre-check is
real, so the window is small — but §17.10's no-overwrite publish rule is satisfied only by the
pre-check, not by the rename. The report's §6 sentence should be corrected, and a real helper should
either drop `REPLACE_EXISTING` or verify post-publish identity.

### S2-3 (MINOR — evidence gap) Kill-before-delete *ordering* is not demonstrated against a real child

A1-D6 requires terminate → observe exit → delete. Stage 2's forced-cancel and timeout tests assert the
**end state** (dead PID, no final, no `.part`), which is consistent with that ordering but does not
distinguish it from delete-then-kill. The only ordering proof is Stage 1's
`TestRejectedCleanupTerminatesBeforeRemovingOutput`, which uses the in-memory `fakeProcess`, not an OS
child. So ordering at the real process boundary is inferred, not observed. The corrective report's
phrase "starts a real fake process" means "through the real supervisor", which is accurate but easy to
misread — worth pinning explicitly.

### S2-4 (MINOR — real-boundary coverage gaps)

Not exercised against a real child, and not claimed as such, but worth listing because each is a
contract path the fake helper *could* model with one more mode:

- late success arriving after the parent has timed out (Stage 1 covers it with the fake transport only);
- a helper that ignores `shutdown` (shutdown-grace kill path);
- a kill *during* WAV writing — the only path that can leave a `.part` behind. `hang-synth` blocks
  before `createPart`, so `.part` cleanup after a hard kill is never exercised end-to-end;
- `request_not_cancellable` (cancel during `load`) and `busy` through the real client;
- the fake's own post-publication-cancel branch (`server.go`, "cancelled after publication") is
  unreachable in the current tests, since cooperative cancel releases before the part is created.

### S2-5 (INFORMATIONAL — precision) "Twelve real-process tests" is eleven

`TestStage2MissingHelperIsNonfatal` deliberately launches nothing (fake `LookupEnv`). Eleven tests use a
real child. The matrix itself is accurate per row; only the §5 sentence overcounts.

### S2-6 (INFORMATIONAL — chronology, correctly disclosed) The suite is not deterministically green

My independent `go test ./...` **failed** in `internal/audio`:

```text
--- FAIL: TestConcurrentDifferentConfigurationsRemainIsolated (0.03s)
    pipeline_cache_test.go:401: concurrent different stats = hits 5 misses 2, want 0/2
```

Identical test, identical signature to the disclosed failure. Three signals confirm it is not Stage 2:

1. isolated `go test ./internal/audio` passes (3.411 s) and focused `-count=3` passes;
2. `internal/audio/pipeline_cache_test.go` mtime is **07:56**, well before every Stage 2 file
   (16:05–16:17), and no Stage 2 file is in `internal/audio/`;
3. it fails only under full-suite load — a load-sensitive concurrency assertion, not a helper defect.

**Chronology, preserved rather than flattened:** earlier full run failed → report disclosed it →
final retry passed → my independent full run failed again → isolated and focused runs pass. The honest
summary is *"Stage 2 packages are stable; the repository suite is not deterministically green because
of a pre-existing flake in `internal/audio`."* It must not be recorded as a Stage 2 pass.

### S2-7 (INFORMATIONAL — admitted limitations, confirmed accurate)

Goroutine-leak evidence is `baseline+2` tolerance, not a registry; there is no waiter counter (single
`Wait()` is structural, and `GetExitCodeProcess` proves exit, not reaping); the final-file symlink test
is a genuine `SKIP` for privilege. All correctly labeled as limitations.

---

## 3. What genuinely survives challenge

- **The real-process boundary is real.** `TestMain` builds the helper with `CGO_ENABLED=0` and exits 1
  on build failure (no silent skip); tests use the production `Discoverer` → `Supervisor.Start` →
  `StartProcess` → `proc.Start`. I confirmed `internal/proc.Start` creates with `CREATE_SUSPENDED`,
  assigns to the job object, resumes only afterwards, and kills + errors if either step fails — so a
  nil-error start *implies* ownership. No fake transport, TCP, pipe, or audio IPC is used.
- **Parity is not merely mutual agreement — verified by probe.** I wrote a temporary test that feeds the
  eight contract-native golden **request** fixtures (hand-written, `piper`-flavoured, not fake output)
  through the fake helper's own strict codec: **all eight accepted**, including per-operation strict
  parameter decoding. That is the direction the submission never tested, and it passes. Response
  direction is covered end-to-end by the twelve interop tests through the real client decoder.
- **Fault injection lives at the process boundary.** Modes are env-selected child behaviours
  (`hold`, `hang`, `crash-on-health`, `malformed-hello`, `incompatible`, `bad-wav`). No parent-side hook
  short-circuits the behaviour being proven. The only white-box seam is
  `c.opts.OperationTimeout = 100ms` in one test, which changes real configuration rather than bypassing
  logic — acceptable, worth noting.
- **Cancellation and publication barrier, as tested:** cooperative cancel keeps the helper alive with no
  output; forced cancel and timeout both kill the real PID and leave nothing; `bad-wav` proves the
  parent distrusts helper metadata and rejects at the acceptance/validation step with exact cleanup.
- **WAV discipline:** canonical 44-byte header, deterministic ramp, `O_EXCL` part file, self-validation
  before publish, hash re-checked by the parent. Parent-side containment (absolute + `Rel` + leaf
  lstat + exact-path cleanup from caller parameters) is unchanged from Stage 1 and still sound.
- **Failure isolation:** malformed, incompatible, crashed, and missing helpers all produce typed
  failures with no surviving PID, no output, and no reader growth; root Chat passes.
- **Scope is clean.** No provider registration, no `tts/helper` import in app/pipeline code, no CGO, no
  native runtime, no model/asset/UI/cache change. The fake helper is a command plus a test fixture.

---

## 4. Disposition

| Area | Assessment |
|---|---|
| Real Windows process boundary | **Pass** — genuinely exercised, ownership enforced by construction |
| Protocol parity (request direction) | **Pass** — plus my fixture probe |
| Protocol parity (capability layer) | **Gap — S2-1** |
| Cancellation / timeout / shutdown | **Pass as tested**; late-response and shutdown-ignore paths unexercised (S2-4) |
| WAV generation and parent validation | **Pass**; publish-overwrite claim inaccurate (S2-2) |
| Kill-before-delete ordering at OS level | **Inferred, not observed** (S2-3) |
| Failure isolation and leak checks | **Pass** within the admitted limits |
| Test independence | **Pass** — real boundary, no bypassing hooks |
| Scope | **Pass** |
| Full-suite determinism | **Fail, correctly disclosed, not Stage 2** (S2-6) |

**Recommendation: ACCEPT Stage 2 for its stated milestone — real-process interoperability is proven —
with S2-1 recorded as a required precondition for Stage 3 and S2-2 corrected in the record.**

Rationale: the milestone the owner set ("we are no longer testing only simulated process behaviour")
is met and independently reproduced, and parity survives except at the capability layer. That layer
cannot be meaningfully tested with a helper that advertises everything correctly; it becomes load-bearing
exactly when a real provider starts advertising a partial capability set. Fixing it now, while the fake
helper is the only helper, is far cheaper than discovering it during native integration.

**Not authorized by this review:** Stage 3, native Kokoro/G2P, ONNX Runtime, provider registration,
pipeline/character/cache/UI integration. The deferrals (restart budget, backoff, idle close; F-4, F-5)
remain exactly as recorded, and no §17.8 provisional value may be cited as measured.

---

## 5. Note on the Stage 1 chronology

Stage 1 is recorded by the owner as accepted and closed via `CORRECTIVE-…-stage1-2026-11-06.md`, which
I had not reviewed before this pass. On inspection, the three findings I raised are addressed in code
and backed by real regression tests: `pendingRequest.deliver` is non-blocking and first-outcome-only,
and `TestFailAllFullOutcomeDoesNotBlockOrSkipKill` asserts the full-channel case *and* `waitCount == 1`;
the failed-handshake path now terminates and observes the old reader before retry; `cleanupRejected`
gates deletion on termination. Closure is defensible. Two precision notes: that corrective report has
not had an independent line-by-line review, and its F-3 ordering evidence is at the client level over
the in-memory process, not over an OS child (S2-3).

---

## 6. Boundaries of this review

**Performed:** 22 hash verifications; full read of `fakehelper/{protocol,server,wav,pathcheck_*}.go`,
`cmd/howl-tts-helper-fake/main.go`, all `interop_*_test.go`, the carried Stage 1 corrective code paths;
contract cross-check against §§17.3–17.6, 17.10, 17.11; re-execution of every reported command; one
full-suite run; one temporary parity probe (deleted; package restored to six files).

**Not performed:** no implementation change, no fix, no permanent test, no helper or fake modification,
no native build, no download, no pipeline or UI work, no commit, push, branch, staging, or deploy.
