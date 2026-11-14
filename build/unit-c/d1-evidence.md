# D1 implementation — evidence package

Date: 2026-11-04 (evening session, unattended)
Authorization: D1 only, Revision 7. D3 not touched. No commit made.
Owner acceptance: 2026-11-05 — **ACCEPTED — OWNER ACCEPTANCE WITH DOCUMENTED
VERIFICATION LIMITATIONS**; see §10.

## 1. Implementation summary

Two new files plus wiring. The iframe never holds the lock or the real backend
promise; the host document does.

New:
- `frontend/dist/js/iframe-bridge-host.js` — parent-side host: generation
  counter, per-modal listener, pending map + timers + retired ids, activation
  lock, `activation-settled` notification, origin/source validation, seven-method
  allowlist.
- `frontend/dist/js/wails-bridge-iframe.js` — iframe shim, installed only when
  `window.go` is absent. Envelope handshake, queued pre-ready calls, generation-
  qualified ids, timeouts, retirement, `pagehide` cleanup.
- `internal/audio/manager_applied_test.go` — applied-generation tests.
- `tests/e2e/specs/iframe-bridge.spec.js` — 6 bridge tests over a real HTTP
  origin.

Modified:
- `frontend/dist/js/settings-state.js` — `_writeBackendSettings` throws instead
  of returning silently when the binding is absent.
- `frontend/dist/js/audio-settings.js` — `_confirmApplied`, `_bridgeReady`,
  `_activationBusy`, `_refreshApplyInProgress`; Save gated on handshake; Retry
  gated on lock/`apply_in_progress`; `Applied` only on confirmed effective
  configuration, otherwise Indeterminate.
- `frontend/dist/js/menu-functions.js` — attach on modal open, detach on remove.
- `frontend/dist/audio-settings.html` — shim loaded before `settings-state.js`.
- `frontend/dist/{index,lorebooks,world,character,scenario,story-tools}.html` —
  host loaded immediately before `menu-functions.js`.
- `internal/audio/manager.go` — `AppliedSettings()` returns the active slot's
  settings under the manager mutex.
- `app-audio.go` — `GetAudioStatus` publishes `apply_in_progress`
  (`status.Builds > 0`) and `applied` (active generation only).
- `tests/e2e/specs/audio-settings-lifecycle.spec.js` — fake backend now
  publishes `applied` so the specification's `Applied` assertions remain valid.
- `build/unit-c/d1-plan.md` — §0 marked measured; Revision 8.

No generic forwarding, no `eval`, no wildcard origin in production, no settings
schema change, no changes to `applyAudioSettings`' return contract.

## 2. Packaged origin gate (mandatory)

Measured on the real packaged binary with a temporary `ReportOrigin` binding and
probe script, both deleted afterwards:

```
window.location.origin === "http://wails.localhost"
```

Decision: **non-opaque, so §3.2.1's production rule applies unchanged**
(`event.origin === window.location.origin`, origin non-empty and not `"null"`).
No wildcard, no opaque-origin exception, no fallback enabled in production. The
`file://` source-only path remains test-only.

Side effects of the measurement, and their disposition:
- The probe app spawned `llama-server` (the already-known D3 behaviour). Its PID
  was confirmed as a child of the probe PID and stopped by PID.
- The probe rewrote `%APPDATA%\HOWLChat\chats\2026-10-31T14-06-59Z\history.json`.
  It was restored from a pre-run backup and re-hashed to the baseline value.
- No audio playback, no device routing, no microphone use occurred.

## 3. Commands executed

| Command | Result |
|---|---|
| `gofmt -l internal/audio/manager.go internal/audio/manager_applied_test.go app-audio.go` | clean after one `gofmt -w app-audio.go` |
| `go build ./...` | exit 0 |
| `go vet ./...` | exit 0 |
| `go test ./internal/audio/` | `ok howl-chat/internal/audio 1.965s`, exit 0 |
| `go test ./internal/audio/ -run TestAppliedSettings -v` | 2 PASS |
| `node --check` on host, shim, audio-settings, settings-state, menu-functions, both specs | exit 0 each |
| `CGO_ENABLED=0 GOFLAGS=-buildvcs=false go build -tags "desktop,production" -ldflags "-w -s -H windowsgui" -o build/bin/howl-chat-d1.exe .` | exit 0 |
| `npx playwright test specs/iframe-bridge.spec.js specs/audio-settings-lifecycle.spec.js --config=<temporary headless config>` | **14 passed**, exit 0 |

Artifact: `build/bin/howl-chat-d1.exe` (22,567,936 bytes), built after the probe
code was removed.

## 4. Plan tests 1–14 — actual status

| # | Test | Status | Evidence |
|---|---|---|---|
| 1 | Bridge unit tests (source/origin/malformed/duplicate/late) | **PARTIAL** | unknown-method rejection and source scoping asserted; malformed payload, duplicate id and late-reply counters not exercised |
| 2 | Standalone regression | **PASS** | 8/8 `audio-settings-lifecycle` tests pass |
| 3 | HTTP parent-page e2e | **PASS (fake backend)** | 6/6 `iframe-bridge` tests over `http://127.0.0.1`; not `harness/shim.js` |
| 4 | 5× open/close lifecycle | **PARTIAL** | one detach/reattach asserted: generation increments, listener removed; 5× repetition not run |
| 5 | Real app, fixture intact | **NOT RUN** | requires driving the packaged panel; only the origin probe ran the app |
| 6 | Late responses | **NOT RUN** | needs budgets shortened or a real 180 s wait |
| 7 | Origin matrix | **PARTIAL** | `file://` (lifecycle spec) and non-opaque HTTP both exercised; opaque-origin rejection not explicitly asserted |
| 8 | Activation serialization | **NOT RUN** | needs two documents against the real backend |
| 9 | Settle after retire | **NOT RUN** | — |
| 10 | Real parent/iframe lifetimes | **PARTIAL** | detach/reattach + no inherited pending entry; multiple cycles not run |
| 11 | Navigation by phase | **NOT RUN** | — |
| 12 | Backend serialization proof | **NOT RUN** | `applyMu` verified by source reading only, not re-proven at runtime |
| 13 | `Builds` gauge + duplicate behaviour | **PARTIAL** | `AppliedSettings` tested; `apply_in_progress` mapping and sequential-rebuild not asserted |
| 14 | `Applied` qualification | **PASS** | matched → `Applied`; withheld `applied` block → `Indeterminate — effective configuration unconfirmed` |

## 5. Acceptance rows

Per-row status, so the summary cannot be read more strongly than the tests:

| Row | Status | Basis |
|---|---|---|
| 1 | PARTIAL | script order and handshake tested; bridge-unit refusal paths partly |
| 2 | PARTIAL | source scoping and production origin rule exercised; opaque-origin rejection not asserted |
| 3 | PASS | unlisted-method test |
| 4 | PARTIAL | correlation and timeout code inspected; no timeout exercised |
| 5 | PARTIAL | failed-write surfaced; timeout error text not exercised |
| 6 | PASS | Save disabled until handshake/hydration |
| 7 | PARTIAL | false labels covered in tested fixture paths only |
| 8 | NOT RUN | 5× listener count |
| 9 | PASS | standalone suite |
| 10 | PASS | hashes unchanged |
| 11 | PARTIAL | no opaque production trust; negative case not asserted |
| 12 | NOT RUN | late reply |
| 13 | PASS | origin measured |
| 14 | PARTIAL | host refusal now tested; cross-document overlap not |
| 15 | PARTIAL | unresolved lock keeps Apply locked; settlement path untested |
| 16 | PARTIAL | ownership in code; single detach/reattach only |
| 17 | PARTIAL | generation-qualified ids; old-result delivery not driven |
| 18 | NOT RUN | cross-navigation serialization (source review only) |
| 19 | NOT RUN | navigation by phase |
| 20 | PARTIAL | `apply_in_progress` consumed; gauge latching not asserted via status |
| 21 | PARTIAL | UI gate tested; sequential duplicate rebuild not |
| 22 | PASS | `Applied` requires equality (fixture) |
| 23 | PASS | unconfirmable state renders Indeterminate (fixture) |

## 6. Protected state — before/after

| File | Baseline | After |
|---|---|---|
| `config/settings.json` | `e786f53f…1757bf` | `e786f53f…1757bf` (unchanged) |
| `%APPDATA%\HOWLChat\config.json` | `d43ce0d6…eef132` | `d43ce0d6…eef132` (unchanged) |
| `…\chats\2026-10-31T14-06-59Z\history.json` | `fd9eee44…a28639` | `fd9eee44…a28639` (restored after probe) |

ASR settings remain `asrEnabled: true`, `asrEngine: qwen3-asr-1.7b-local`,
devices `default`. No audio device routing was performed; nothing was played.

## 7. Limitations and unresolved items

- **Bridge e2e uses a fake backend**, not the packaged app. It proves protocol,
  ownership, origin handling and label truthfulness, not real model activation.
  It is not physical-device or end-to-end proof.
- **Timeout paths are not exercised end to end.** Budgets are 20 s–300 s; the
  retirement logic is implemented but only inspected, not driven.
- **Sequential-duplicate and cross-navigation serialization** rely on the
  backend `applyMu` as established by source review; no runtime proof was added.
- **`superseded` is still discarded** by `applyAudioSettings`; a superseded call
  can still return `nil`. The panel no longer converts that into `Applied`
  unless the published effective configuration matches, but the backend
  truthfulness gap remains (recorded, not fixed — outside D1).
- **D3 ASR child lifecycle is untouched.** The probe app did orphan a
  `llama-server`, which was stopped manually by PID; this is the known D3 defect,
  not a regression.
- A headless Chromium shell was installed to `%LOCALAPPDATA%\ms-playwright` to
  run the specs; the temporary headless config was deleted after the run.
- No commit, no push, no deployment.

## 8. Correction batch (2026-11-05)

Authorized scope: host lock enforcement, tri-state status handling,
restored-lock no-release policy, additive `desired_gen`/`applied_gen`, targeted
tests, and documentation corrections. `Apply`'s return contract and the
`superseded` bool are unchanged.

Results:

| Check | Result |
|---|---|
| 10 activation-lock harness tests | **10 passed (13.3s)**, exit 0 |
| `TestStatusExposesQueuedGeneration` | PASS |
| `TestAppliedSettingsSupersededGeneration` | PASS (characterization only) |
| `go test ./internal/audio/` | ok, exit 0 |
| `go build ./...` | exit 0 |
| `node --check` (host, panel) | exit 0 each |

Defect found by the batch: the panel read `window.howlBridgeHost` inside the
iframe, where it is undefined, so the host's restored lock was invisible to it.
The panel now resolves the same-origin parent host. Tests 5–9 caught this; they
failed before the fix and pass after it.

Correction limits:
- The 10 tests are fixture-backed over HTTP. No packaged app, no physical device.
- The restart boundary is **not proven**. The fresh-browsing-context test is a
  proxy only.
- The earlier 14 tests (6 bridge + 8 lifecycle) were **not re-run** in this
  batch, per the authorized check list; their earlier result is not evidence
  about the latest panel edit.
- Plan tests 5, 6, 8, 9, 11, 12 remain NOT RUN; acceptance rows 20–21 remain
  PARTIAL; D1 remains acceptance-pending.

Regression re-run after the correction: `iframe-bridge.spec.js` +
`audio-settings-lifecycle.spec.js` = **14 passed (4.8s)**, exit 0. The stricter
`applyInProgress === false` requirement and the parent-host lookup did not break
the existing bridge or lifecycle behavior.

## 9. Dismissal correction (2026-11-05, authorized)

Production changes (both in `audio-settings.js`): `close()` dismisses via
`_finishClose()` when a restored lock is held, before any save decision; and
`saveAudioSettings()` refuses before `flushSettings()`, so a programmatic call
cannot write settings.

Result: **15 of 16 harness tests passed; 1 failed; batch stopped** under the
authorization's stop condition (§7: a check failed).

| Test | Result |
|---|---|
| 11 Escape (panel) | PASS |
| 12 close button | PASS |
| 13 overlay click | **FAIL — route cannot be exercised** |
| 14 parent-window Escape | PASS |
| 15 reopen after dismissal | PASS |
| 16 programmatic save refused | PASS |
| `node --check` / `go build ./...` | exit 0 / exit 0 |

Root cause of the failure: the spec fixture mounts the panel in a bare
`<iframe id="f">` and does not load `menu-functions.js`. There is therefore no
`#audio-settings-modal` element and no parent overlay-click handler, so the
overlay route does not exist in the harness. `locator('#audio-settings-modal')`
timed out after 30 s.

Consequence for the other assertions: "parent modal removed" is vacuously true
in tests 11, 12 and 14 because that element never existed. Those tests still
prove zero `SetSettings` calls, zero `ApplyAudioSettings` calls, an unchanged
`activation()` snapshot and a surviving session marker.

A later authorized fixture-only correction made the overlay route testable. The
earlier failed run above is retained as history.

### 9.1 Fixture correction (2026-11-05, test-only)

The harness now loads the real `menu-functions.js` and opens the modal only
through `openAudioSettings()`; nothing in the parent integration is re-written
in the fixture. Added verification:

- Stub isolation: `menuFunctions.openAudioSettings === openAudioSettings`, its
  body contains `_openAudioSettingsIframe`, and no stub name contains "audio".
- Overlay click: the chosen point satisfies both
  `elementFromPoint(x, y).id === 'audio-settings-modal'` and lies outside
  `.real-settings-container`'s bounding rect.
- Page errors are collected from before navigation and asserted empty.
- Parent Escape handler lifecycle: keydown listener count is 1 after load, 2
  while the modal is open, and back to 1 after dismissal.
- Every dismissal route asserts: modal removed, `listening === false`, no
  `audio-settings.html` frame, zero `SetSettings`/`ApplyAudioSettings` calls,
  unchanged `activation()` snapshot, and the session marker still present.

Result: **16/16 pass (16.6s)**, exit 0. `node --check` exit 0; `go build ./...`
exit 0.

First run of this fixture: 14 passed, 2 failed (`Frame was detached` in the
"authoritative idle" and "fresh context" tests). Cause was test-side, not
production: a successful activation closes the modal, detaching the iframe while
`panel.evaluate` was still returning. Both call sites now fire-and-forget and
poll the fixture counter.

No protected settings or history files were touched; no commit was made.

## 10. Owner acceptance disposition (2026-11-05)

**D1 — ACCEPTED — OWNER ACCEPTANCE WITH DOCUMENTED VERIFICATION LIMITATIONS.**

The owner accepted D1 in its then-current implemented state on 2026-11-05, on
the basis recorded in `build/unit-c/d1-plan.md` §8: the iframe bridge and
production-origin validation; parent-owned locking and generation-qualified
handling; truthful Saved / Applied / Ready distinctions including the
restored-lock behavior; fail-closed handling of unresolved activation state; the
two production dismissal/save guards; the corrected 16-test harness over the
real parent integration; and the reported 16/16 activation-lock result with
`node --check` and `go build ./...` at exit 0.

This section supersedes the acceptance-pending status statements earlier in this
document (including §4 and §9). Those statements are left unaltered as history;
their results, failed runs, corrections and chronology stand as recorded.

### 10.1 Limitations that remain explicit

Acceptance does not establish:

- the packaged-app restart boundary (test 10 is a fresh-context proxy only);
- a rerun of the 14 bridge/lifecycle regression tests;
- the remaining packaged-app verification and, in §5, the plan tests 5, 6, 8,
  9, 11, 12 marked NOT RUN and acceptance rows 20–21 marked PARTIAL;
- physical-device or real-audio behavior.

These are **unverified/deferred**. They are not PASS results, are not implied to
have been verified, and are not evidence of a known implementation failure. The
§5 table still shows every earlier PARTIAL and NOT RUN entry unchanged.

Protection state at acceptance: `settings.json` `e786f53f…1757bf`, AppData
`config.json` `d43ce0d6…eef132`, `history.json` `fd9eee44…a28639`; repository
head `4f5defe`; no commit, push, release or deployment.

