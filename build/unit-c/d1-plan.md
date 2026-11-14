# D1 — Audio Settings backend bridge: implementation plan

Status: **ACCEPTED — OWNER ACCEPTANCE WITH DOCUMENTED VERIFICATION LIMITATIONS**
(2026-11-05). Superseded status, retained for the record: IMPLEMENTED —
ACCEPTANCE PENDING (Revisions 7–9), during which the activation-lock correction
(§3.6.3) was implemented and its correction evidence was pending review.
Acceptance was granted with the limitations recorded in §8; it does not assert
that the packaged-app and lifecycle checks are complete.
Gate status: §0 origin measured as `http://wails.localhost` (non-opaque; rule
satisfied). 3.4.1 `applied` identity implemented in `GetAudioStatus`.
Implementation and test evidence: `build/unit-c/d1-evidence.md`.
Revision 2: added the opaque-origin rule (3.2.1) and the late-response
retirement rule (3.3).
Revision 3: added the activation-serialization rule (3.6).
Revision 4: made host/iframe ownership explicit (3.0) and added the
document-lifetime rule (3.6.1) and the old-generation response rule (3.7).
Revision 5: replaced the `sessionStorage`-as-lock assumption with the verified
backend `applyMu` guarantee, recorded what `GetAudioStatus` does not publish,
and reported the required additive status field (3.6.1).
Revision 6: verified `Builds` is an active gauge (dropped the `Phase` disjunct)
and added the duplicate-activation rule (3.6.2).
Revision 7: adds the `Applied` qualification (3.4.1) — a nil error is not proof
of effective configuration, and the current API cannot establish it.
Revision 8: records the measured packaged origin and marks §0 satisfied.
Revision 9 (2026-11-05): records that implementation was authorized and
performed; adds §3.6.3 (lock enforcement, restored-lock no-release policy,
`desired_gen`/`applied_gen` observability). Acceptance remains pending.

## 0. Open item that gates implementation

**MEASURED — gate satisfied.** The packaged app reports
`window.location.origin === "http://wails.localhost"`.

Measured on the real packaged binary with a temporary build-only probe
(`ReportOrigin` binding writing `window.location.origin` to a file), which was
then deleted and the app rebuilt without it. The value is non-opaque, so the
production rule in 3.2.1 (`event.origin === window.location.origin`, origin
non-empty and not `"null"`) is satisfiable exactly as written and is applied
unchanged. No fallback, wildcard, or opaque-origin exception is enabled in
production.

## 1. Scope

Restore a supported backend path for the Audio Settings panel and make its
save/apply reporting truthful.

In scope:
- A bridge that exists before `settings-state.js` / `audio-settings.js` initialize.
- A host module in the six parent documents that owns the lock and the real
  backend promise.
- Additive `GetAudioStatus` fields only: `apply_in_progress`, and an `applied`
  identity block read from the active slot. No settings schema change, no new
  API surface.
- Truthful `Saved` / `Applied` / `Ready` reporting.
- Focused regression tests.

Out of scope (explicitly excluded):
- D3 ASR child-process lifecycle.
- Settings schema changes.
- ASR/TTS engine behaviour or lifecycle.
- Chat persistence (D2 withdrawn — no defect).
- The existing `postMessage(..., '*')` callers in `character-form.js` / `world-form.js`.
- Audio device routing.

## 2. Established facts

Observed in the packaged app:
- The panel is an iframe: `menu-functions.js` builds `<iframe src="audio-settings.html#...">`.
- Inside that iframe the backend is unreachable: both cards read
  "Status: backend not available (dev mode)"; Save reports
  "Apply failed — Error: audio backend is unavailable".
- Save reported **"Saved"** while `config/settings.json` mtime stayed unchanged
  (predating the save) — a false success.

Established from source:
- `settings-state.js` `_writeBackendSettings()` does
  `if (!window.go?.main?.App?.SetSettings) return;` — it returns silently, so a
  missing backend cannot fail the save path.
- `audio-settings.js` sets `Saved` at `:582`/`:617` *before* apply, and
  `Saved`/`Applied` at `:545` from `_readReadiness`.
- The main document already installs `message` listeners (`sidebar-panels.js:95/117/131`).
- `harness/shim.js` already provides a supported `window.go` proxy over `/__rpc`
  (build tag `uiharness`) — the existing precedent for a non-Wails host.
- `tests/e2e/specs/audio-settings-lifecycle.spec.js` loads the panel standalone
  from `file://` with an in-page `window.go`; the bridge must no-op when
  `window.go` already exists.

Not established (do not assume):
- What Wails v2.12 injects into subframes. No `vendor/` directory exists and the
  module cache is outside the project, so the plan does not depend on it.
- The document origin the packaged app is served from (see section 0).

## 3. Design

Two files. The lock and the real backend promise must outlive the panel, so they
cannot live in a script that only the panel loads.

**`frontend/dist/js/iframe-bridge-host.js` — new, parent-side.** Loaded in every
document that can open the modal, immediately *before* `js/menu-functions.js`:
`index.html:85`, `lorebooks.html:355`, `world.html:262`, `character.html:327`,
`scenario.html:239`, `story-tools.html:65`. Exposes `window.howlBridgeHost` and
owns everything that must survive panel teardown.

**`frontend/dist/js/wails-bridge-iframe.js` — new, iframe-side.** Loaded as the
**first** script in `audio-settings.html`, before `settings-state.js`
(`audio-settings.html:293`). Owns only the panel's request handles and its own
presentation state.

```
// iframe-side: installs only when the native bridge is absent
if (!window.go) { window.go = installShim(); }
```

Installing only-if-absent means a future Wails version that injects into frames
wins automatically, the standalone `file://` spec is unaffected, and the panel
still has no backend of its own to fall back on.

### 3.0 Ownership

| Concern | Owner | Lifetime |
|---|---|---|
| Native `window.go` | Wails, main document | app process |
| Real `ApplyAudioSettings` promise | host | until settlement, or document destruction (3.6.1) |
| Activation lock (advisory) | host | document, mirrored to `sessionStorage`; the backend `applyMu` is authoritative (3.6.1) |
| Generation counter, pending map, timers, retired ids | host | document |
| `message` listener install/remove | host, driven by `menu-functions.js` open/remove | one per modal |
| Request `seq`, its own promise handles, `Saved`/`Applied`/`Ready` | iframe shim | iframe document |

The iframe never holds the lock, never holds the real backend promise, and never
decides whether an activation may start. Closing the iframe retires only its own
request handles (3.3) and removes the host listener for that window; it does not
release or cancel the host-owned operation.

Reopening is therefore a new generation with no inherited state: the panel
issues fresh `GetAudioStatus`/`GetSettings` calls and cannot adopt the previous
activation's request or result as its own.

`ALLOWED` — exactly seven methods, derived from real call sites:

| Method | Call site |
|---|---|
| `GetAudioStatus` | `audio-settings.js` status poll + `_readReadiness` |
| `ApplyAudioSettings` | `audio-settings.js` `_applySnapshot` |
| `GetTTSModels` | `audio-settings.js` `populateTTSDropdown` |
| `OpenDirectoryPicker` | `audio-settings.js` `scanTTSDirectory` |
| `ScanTTSDir` | `audio-settings.js` `scanTTSDirectory` |
| `GetSettings` | `settings-state.js` `loadSettings` |
| `SetSettings` | `settings-state.js` `_writeBackendSettings` |

Any unlisted property returns `undefined`, so an unlisted call fails loudly
rather than silently succeeding.

### 3.1 Protocol

iframe → parent:
```
{ __howlBridge: 1, type: 'call', id: '<generation>:<seq>', method: <string>, args: [..] }
{ __howlBridge: 1, type: 'hello' }
```
parent → iframe:
```
{ __howlBridge: 1, type: 'ready' }                       // handshake ack
{ __howlBridge: 1, type: 'result', id, ok: true, result }
{ __howlBridge: 1, type: 'result', id, ok: false, error: { message } }
{ __howlBridge: 1, type: 'activation-settled' }          // uncorrelated: no id
```

`result` replies only ever carry the id of a request issued by that same
generation. `activation-settled` is not a reply and matches nothing (3.7).

Handshake: bridge posts `hello` on load. Calls made before `ready` are queued
and dispatched after it, inside the same timeout budget. If no `ready` arrives
within `HANDSHAKE_TIMEOUT` (5 s), the bridge transitions to `unavailable`
permanently for that document.

### 3.2 Validation (parent side)

Accept a message only when **all** hold:
1. `event.source === iframe.contentWindow` — primary identity check.
2. Origin rule — see 3.2.1.
3. `data && data.__howlBridge === 1`.
4. `typeof data.method === 'string' && ALLOWED.has(data.method)`.
5. `Array.isArray(data.args)`.

Rejections are dropped silently from the caller's perspective (no reply is sent
to a non-matching source) and counted in a diagnostic counter. Unknown methods
received from the *correct* source get an explicit `ok: false` so the bug
surfaces instead of hanging.

The parent never forwards `window.go` wholesale and never evaluates strings.

### 3.2.1 Origin rule

Opaque-origin equality (`'null' === 'null'`) proves nothing, so it is **never**
a production trust rule.

- **Production path** — `window.location.protocol !== 'file:'`.
  Accept only when `event.origin === window.location.origin` **and**
  `window.location.origin` is non-empty and not `'null'`.
  If the packaged document origin turns out to be opaque, the bridge must be
  extended with an equally strong alternative boundary *and* documented here
  before it is enabled. Until then it fails closed: Save disabled, backend
  reported unavailable. It must **not** degrade to trusting `event.source`
  alone.
- **Test path** — `window.location.protocol === 'file:'`.
  Accept on `event.source === iframe.contentWindow` alone. This exists only for
  the standalone Playwright spec, which loads `audio-settings.html` directly
  from disk, and is unreachable from the packaged app.
- The production branch must be exercised under the HTTP harness as well as the
  `file://` branch, so the permissive path can never become the only tested one.
  The harness (`harness_main.go`) serves the same `frontend/dist` over HTTP at
  `127.0.0.1:8799`, giving a verified non-opaque origin
  (`http://127.0.0.1:8799`) for that branch. That is a test origin only; the
  packaged Wails origin remains the section 0 gate.

### 3.3 Correlation, timeouts, retirement

- `id` is namespaced by generation: `` `${generation}:${seq}` ``. Replies echo
  `id` verbatim and pending entries are keyed by that full string, so a reply
  from generation N cannot match a request issued in generation N+1 even when
  sequence numbers coincide.
- Each pending entry stores `{ generation, method, timer, resolve, reject }`.
  A reply is accepted only when `entry.generation === currentGeneration`.
- Budgets: `GetSettings`/`GetAudioStatus`/`GetTTSModels` 20 s;
  `SetSettings` 60 s; `ApplyAudioSettings` 180 s (local model init dominates
  this and is not bounded in practice); `OpenDirectoryPicker` 300 s.
- **On timeout:** reject with `BridgeTimeoutError`, **delete** the pending entry,
  and add its `id` to a `retired` set for that generation. A reply arriving
  later finds no entry, is dropped, increments a `lateReply` counter, and
  changes no state.
- **On close:** increment `generation`, reject every pending call with
  `BridgeClosedError`, clear all timers, clear `retired`, and remove the parent
  listener. The parent sets `active = false` before teardown, so a backend
  promise that settles after close never posts a reply at all.
- Retirement is enforced by the pending-map lookup plus the generation check —
  **not** by the original promise having already rejected. A settled promise
  must not be the only thing standing between a late reply and a reopened
  panel.

### 3.4 Fail-closed states

| State | Required evidence |
|---|---|
| `Saved` | `SetSettings` resolved |
| `Applied` | `ApplyAudioSettings` resolved **and** the published effective configuration matches the request (enabled flags + engines). A nil error alone is not proof — see 3.4.1. |
| `Ready` | fresh `GetAudioStatus` result |
| Any failure | explicit error; no success badge inferred |
| Timed out | explicit error, reported as **indeterminate**: the operation may still be running backend-side. A timeout is not evidence that nothing happened, and it must never render as a partial success. |

Required changes:
- `_writeBackendSettings` must **throw** (`new Error('settings backend unavailable')`)
  instead of `return` when the binding is absent. The main document is unaffected
  because it has `window.go`; the harness shim supplies one too.
- Save stays disabled until `_hydrated` **and** the bridge reports `ready`.
- The apply path must not set `Applied` from a resolved promise; it must compare
  the effective configuration (3.4.1).

### 3.4.1 `Applied` requires confirmed effective configuration

`Applied` must mean *this request's configuration is the effective
configuration*, not that the call returned without error.

Why a nil error is insufficient, verified in source:

- `applyAudioSettings` (`app-audio.go:347`) discards both the status and the
  `superseded` bool from `Manager.Apply`, so a request overtaken by a newer
  desired generation still returns `nil`.
- `GetAudioStatus` publishes readiness (`asr_ready`, `asr_readiness`,
  `tts_ready`) but no engine identity and no generation identifier, so readiness
  cannot be attributed to a specific request.
- `GetSettings` returns the persisted desired settings (`appSettings`), not the
  configuration of the active generation.
- `audioReadiness` reads `audioManager.Desired()`, which is the newest request
  rather than the installed generation.

The effective configuration does exist in the backend: `slot.settings` (with
`asrReady`/`ttsReady`) holds the active generation's settings. It is simply not
published.

**Rules:**

- Render `Applied` only when the request resolved without error **and** the
  effective configuration equals the requested one — enabled flags and engine
  strings.
- If the effective configuration cannot be compared — not published, or replaced
  by a newer generation — render
  `Indeterminate — activation reported success, effective configuration
  unconfirmed`. Never derive `Applied` from the absence of an error.
- `Saved` is unaffected: it reflects `SetSettings` persistence only, and must
  not be reused as evidence about activation.

**Minimal additive option, no schema change.** Publish an `applied` block in
`GetAudioStatus` from the active slot (`asr_enabled`, `tts_enabled`,
`asr_engine`, `tts_engine`). With `apply_in_progress` (3.6.1) this makes
`Applied` provable by equality with the request. Until it exists, the panel's
ceiling is Indeterminate.

### 3.5 Lifecycle safety

- One listener per modal instance, bound to that specific `contentWindow`,
  installed in `_openAudioSettingsIframe`, removed in `_removeAudioSettingsModal`.
- The generation counter plus `active` flag implement the retirement rules in 3.3.
- Reopening must not accumulate listeners, pending entries, or cross-talk.

### 3.6 Activation serialization — apply after an indeterminate timeout

A bridge timeout retires the *client* request (3.3). It does **not** cancel the
backend operation, which may still be initializing a local model and holding its
resources. Permitting Retry at that moment would start a second activation on
top of a live one.

- **One activation in flight per app session.** The lock is a field of
`window.howlBridgeHost` in the parent document (3.0) — never of the iframe shim,
which dies with the panel. An iframe generation change, including close and
reopen, must not release it.
- **The real promise is host-owned.** The host keeps the actual
`window.go.main.App.ApplyAudioSettings(...)` promise. The iframe's request is a
proxy for it; retiring the proxy (3.3) never cancels the promise.
- `retryAudioActivation()` (`audio-settings.js:608`) currently guards only on the
per-document `_saving` flag. It must consult the host lock as well, and the
`#audio-retry-activation` button must be disabled while the lock is held.
- **Retirement is client-side only.** When an `ApplyAudioSettings` call times
out, the parent retires the pending entry but keeps the real backend promise
and the lock. The client request is gone; the operation is not.
- **Release on actual settlement.** The lock releases when the parent's backend
promise settles and only then — not on timeout, not on close — within the
lifetime of the document (document destruction is handled in 3.6.1). The outcome is
delivered to whatever generation is current, or dropped without a state write
if no panel is open; dropping it may only release the lock.
- **Reconciliation before any second activation.** While the lock is held:
`Apply` and `#audio-retry-activation` are disabled, the panel shows
`Indeterminate — activation in progress`, and reopening refreshes `GetAudioStatus`/`GetSettings` to reflect live backend
state without permitting a new activation.
- **Undeterminable outcome.** If the operation cannot be determined, the panel
reports `State unknown — restart required` and keeps `Apply` disabled for the
rest of the session. It must not offer a second activation as recovery.

Failing closed is intentional: a stuck lock is recoverable by restarting the
app, whereas overlapping activations risk duplicate model loads and conflicting
resource ownership.

### 3.6.1 Document lifetime and cross-navigation serialization

The host is a per-document singleton, and `menu-functions.js` is loaded by six
documents, so navigation destroys the host and any renderer-side lock with it.
That is not the authority. The real serialization is already in the backend.

**Established guarantee (verified in source).** `audio.Manager.Apply` takes
`m.applyMu` at `internal/audio/manager.go:234` with a deferred unlock at `:235`
(field at `:181`) and holds it for the entire call, including the synchronous
factory build that runs inside the loop beginning at `:237`. The Manager is
owned by `App` in the Go process, so renderer document navigation cannot affect
it. Two activations therefore cannot execute concurrently, whichever document
issues them: a second call blocks until the first returns, and then either
rebuilds or is skipped by `canSkipLocked`.

**What the backend does not publish.** Whether an apply is in progress.
`audio.Status` carries `Phase` (set to `"applying"` during the build) and
`Builds` (`st.Builds = m.builds`, `:931`), but `App.GetAudioStatus`
(`app-audio.go:211`) exposes neither: its `active` key is only a
phase-not-closed test and is true when idle. There is also no activation id. A
fresh document therefore cannot tell from the status map whether the previous
activation is still running.

**Rules:**

- Serialization is the backend's responsibility and is not re-implemented in the
  renderer. The host lock remains only to avoid issuing redundant queued work
  and to keep the panel's own states truthful.
- The `sessionStorage` marker is a UI hint. It may drive
  `Indeterminate — activation may still be running`; it is never proof that an
  activation completed, so it must not clear a state or re-enable Apply.
- The panel must never derive completion from `asr_ready`, `asr_readiness` or
  `active`. Those describe the last published generation, not whether a build is
  still running.
- Completion is never inferred. It is either reported by the in-flight request
  in this document, or published authoritatively by the backend (below).
- A reopened panel must not schedule a duplicate activation. It waits until
  `apply_in_progress` is false before enabling `Apply`; the backend does not
  suppress a sequential duplicate (3.6.2).

**Required minimal backend addition — no schema change.** Publish
`apply_in_progress` in `GetAudioStatus` as `status.Builds > 0`. Two lines in
`app-audio.go`; no new state and no new API surface.

`Builds` is an **active gauge, not a cumulative counter** — verified: `m.builds++`
at `:279` immediately before the build; `m.builds--` on success in
`activateLocked` (`:761`) and on every cleanup path in `recordCleanupLocked`
(`:855`); and it is the outstanding-work counter consumed by
`drainConditionLocked` (`:885`) and `outstandingLocked` (`:909`). Existing tests
already assert `Builds == 0` after both success and failure
(`manager_test.go:386`, `:419`, `:460`, `:503`, `:539`) and `== 1` inside the
factory (`:957`–`:974`). A gauge used as the shutdown-drain condition cannot
latch positive without also blocking Close.

`Phase == "applying"` (`:283`) is deliberately **not** used in the expression:
`statusLocked` overrides `Phase` during closing and after terminal, so it is a
display projection rather than a reliable in-flight indicator.

Until that field exists, the panel reports `Indeterminate — activation may still
be running` after navigation. It must not claim either completion or absence of
an activation.

### 3.6.2 Duplicate activation, coalescing, and resource overlap

`applyMu` provides concurrency safety, not content dedupe. Verified: every
`Apply` entry advances `m.desiredGen` (`:229`–`:231`), and `canSkipLocked` skips
only when
`m.appliedGen >= buildGen` (`:733`). Coalescing therefore covers *overlapping*
callers: a queued caller whose arrival superseded an in-flight build lets that
build be disposed and the queued caller reuse the newest generation. A
**sequential** duplicate — the first call fully settled, then the same settings
applied again — advances `desiredGen` past `appliedGen` and rebuilds, including
a full model load.

Consequences recorded for this plan:

- The UI gate is the duplicate-suppression mechanism: a reopened panel waits for
  `apply_in_progress == false` before enabling activation.
- No backend content-dedupe is added. It would change retry semantics after a
  failed build, and no evidence requires it for the safety property.
- Related truthfulness issue, recorded not fixed: `applyAudioSettings`
  (`app-audio.go:347`) discards the `superseded` bool from `Manager.Apply`, so a
  superseded call returns `nil` and can currently be reported as applied.
  Surfacing it (sentinel error or status field) changes an existing binding's
  contract, so it is a separate decision.

**Residual limit, recorded deliberately.** `Apply` returns once the new
generation is installed, while the replaced generation is shut down in a
goroutine (`go m.runShutdown(claimed)`, `:430` and `:485`). So "no overlapping
activations" does not mean "no overlapping resource release": the previous
model may still be tearing down when the next apply begins. This is the
manager's intended retire semantics, and it is not described here as full
resource serialization.

### 3.6.3 Lock enforcement and restored-lock policy (Revision 9)

The host enforces its activation lock: an `ApplyAudioSettings` request that
arrives while the lock is held is refused with `ActivationLockedError` **before**
any pending entry, timer or backend call is created. The lock is no longer only
an indicator.

A lock restored from `sessionStorage` is never released in-process — not by
elapsed time, generation equality, repeated idle reads, the panel, or a manual
control. It is reported as `State unknown — restart required` and blocks Save
and Retry until the application is restarted. The reason is the interval before
the manager observes a call: a request that left the old renderer but has not
reached Go is invisible to every backend counter, so no renderer-side rule can
prove its absence.

Dismissal is permitted and performs no write, no activation, and no lock change;
reopening reproduces the unresolved state.

`GetAudioStatus` publishes `desired_gen` and `applied_gen` as **observability
only**. Invariant: `desired_gen == applied_gen && builds == 0` means the
manager's published state has no outstanding generation or active build. It
does not prove that no caller remains inside `Apply` after installation, and it
never authorizes releasing a lock.

`Apply`'s return contract and its discarded `superseded` bool are unchanged.

### 3.7 Old-generation responses

An activation started by generation N must never be delivered as the answer to a
request from generation N+1.

- Request ids are generation-qualified (3.3), so an N reply cannot match an N+1
pending entry, and retiring a pending entry removes it from the host map.
- When an N activation settles, the host **discards its result payload**. It
sends at most `{ type: 'activation-settled' }` — no `id`, no result, no error.
- That notification has exactly two effects: release the lock, and prompt a
fresh `GetAudioStatus` read initiated by the panel itself.
- It must not resolve or reject any pending request, must not be rendered as
`Applied`, and must not repaint `Saved`/`Ready` from the old result.
- If no panel is open when it settles, the host releases the lock and drops the
payload silently.

## 4. Test plan

1. **Bridge unit tests**: wrong source, wrong origin, unknown method, malformed
   payload, duplicate `id`, reply for an unknown `id`.
2. **Standalone regression**: `audio-settings-lifecycle.spec.js` passes
   unchanged (proves the shim no-ops when `window.go` exists).
3. **Harness e2e**: a parent page embedding the panel with `harness/shim.js`,
   giving a real reflected `window.go`. Assert `local-initialized` readiness, a
   real `SetSettings` write, and no `Saved` badge when `SetSettings` rejects.
4. **Lifecycle**: open/close the panel 5×; assert listener count 1, pending 0,
   no stale promise resolution.
5. **Real app, fixture intact**: truthful states, `config/settings.json` mtime
   changes on save, post-test hashes recorded **and compared** to the fixture.
6. **Late responses**: reply after timeout; reply after close; reply after
   close→reopen. Each must change no state, resolve nothing, and leave the
   reopened panel's pending map untouched.
7. **Origin matrix**: `file://` (source-only trust) and the HTTP harness at
   `http://127.0.0.1:8799` (origin equality required). Assert an opaque-origin
   message is rejected on the production path, and a mismatched non-opaque
   origin is rejected on both.
8. **Activation serialization**: with an activation in flight, attempt a second
   apply, a close→reopen→apply, and an apply after the client timeout. Assert
   exactly one backend `ApplyAudioSettings` call occurred (counted at the
   harness shim) and that `Apply` stays disabled until the first promise
   settles.
9. **Settle after retire**: let the client time out, then let the backend
   promise settle. Assert the lock releases, the retired generation receives no
   state write, and exactly one further activation is then permitted.
10. **Real parent/iframe lifetimes**: start an activation, destroy the iframe
    (`_removeAudioSettingsModal`), reopen a panel, then settle the original
    promise. Assert the new generation has zero pending entries and no inherited
    id; that `Apply` and `#audio-retry-activation` are disabled while the lock is
    held; that the new panel issues its own `GetAudioStatus` read; that no
    `Applied` badge appears from the old result; and that exactly one backend
    `ApplyAudioSettings` call was made. This must drive real iframe creation and
    destruction, not repeated calls against a persistent shim.
11. **Navigation during activation, by phase**: start an activation and navigate
    while it is (a) issued but not yet building, (b) inside the factory build,
    (c) after the new generation is installed, (d) after a failed build. In each
    case the new document must not produce concurrent execution, and its state
    must be truthful: no `Applied` badge without evidence, and
    `Indeterminate — activation may still be running` where the backend has not
    yet said otherwise. Repeat with `apply_in_progress` present and absent.
12. **Backend serialization proof**: issue `ApplyAudioSettings` concurrently
    from two different documents and assert they do not overlap — the second
    build does not begin until the first returns, and at most one build is ever
    active. This exercises the existing `applyMu` guarantee, not a shim.
13. **Gauge semantics and duplicate behavior**: assert `apply_in_progress` is
    false at idle, true inside the factory, and false again after success and
    after a failed build — it must never latch true. Then run two sequential
    identical applies and assert the second rebuilds, proving the UI gate (not
    backend dedupe) is what prevents duplicate activation, while two overlapping
    applies still yield at most one live build.
14. **`Applied` qualification**: assert the badge appears only when the
    published effective configuration matches the request. Then force both
    failure shapes — a superseded apply (nil error, newer generation installed)
    and an uncomparable effective state — and assert each renders
    `Indeterminate — ... effective configuration unconfirmed`, never `Applied`.
    Run this with `applied` present and absent.

## 5. Acceptance criteria

| # | Requirement | Evidence |
|---|---|---|
| 1 | Bridge exists before panel scripts run | script order + handshake test |
| 2 | Strict origin and source validation | unit test 1 + origin matrix 7 |
| 3 | Narrow method allowlist (7) | unlisted-method test |
| 4 | Request/response correlation + timeouts | timeout + late-reply tests |
| 5 | Backend errors surfaced, not swallowed | rejection test |
| 6 | Save disabled until handshake | UI-state test |
| 7 | No false `Saved`/`Applied`/`Ready` | failed-write + timeout tests |
| 8 | No duplicate listeners on reopen | 5× open/close test |
| 9 | Main document bindings unaffected | standalone regression suite |
| 10 | Fixture preserved | pre/post hashes recorded and compared |
| 11 | Opaque origin never trusted in production | origin matrix 7 |
| 12 | Late reply cannot mutate a retired/closed panel | late-response test 6 |
| 13 | Production origin verified and recorded | section 0 |
| 14 | No overlapping activations | serialization tests 8–9 |
| 15 | Indeterminate apply keeps Apply locked until settlement | test 9 + disabled-state assertion |
| 16 | Parent owns lock and real promise; iframe owns only handles and state | §3.0 table + test 10 |
| 17 | Old-generation result never answers a new generation | test 10 |
| 18 | Cross-navigation serialization rests on backend `applyMu`, not on `sessionStorage` | §3.6.1 + test 12 |
| 19 | No completion inferred from status; in-progress state is authoritative or explicitly indeterminate | test 11 |
| 20 | `apply_in_progress` reflects an active build only, never history | test 13 |
| 21 | Duplicate activation prevented by the UI gate, not by assumed backend dedupe | test 13 |
| 22 | `Applied` requires confirmed effective configuration, not a nil error | test 14 |
| 23 | Unconfirmable effective state renders Indeterminate, never `Applied` | test 14 |

## 6. Risks

- **Opaque frame origin.** Under `file://` the origin is `'null'`, where origin
  equality proves nothing. Accepted only on that test path; production requires
  a verified non-opaque origin and otherwise fails closed. See 3.2.1, test 7.
- **Unknown production origin.** Not derivable from project sources; gated by
  section 0 before the bridge is considered complete.
- **Future Wails frame injection.** The shim installs only when `window.go` is
  absent, so a native path wins automatically.
- **Existing loose IPC.** `character-form.js`/`world-form.js` post with `'*'`.
  The bridge namespaces on `__howlBridge` and validates source, so it is not
  confused with those messages. They are left unchanged (out of scope).
- **Apply latency.** 180 s budget; exceeding it reports an indeterminate
  timeout, never success.
- **Held activation lock.** If the backend promise never settles, `Apply` stays
  disabled for the session and the panel reports `State unknown — restart
  required`. Accepted deliberately: a restart is a safe recovery, overlapping
  activations are not.
- **Renderer lock is advisory.** Navigation destroys the host and the real
  promise; the `sessionStorage` marker is a UI hint only. Concurrency safety
  rests on the backend `applyMu` (3.6.1), which is why a queued second apply is
  slow rather than dangerous.
- **Previous generation may still be shutting down.** `Apply` returns before the
  replaced slot finishes `runShutdown`, so activations are serialized but
  resource release can briefly overlap (3.6.1). Not fixed here; recorded.
- **Sequential duplicates rebuild.** The backend coalesces only overlapping
  callers; an identical later apply reloads the model. Serialized but costly, so
  the panel's gate (3.6.2) is load-bearing rather than cosmetic.
- **`Applied` is unprovable until the `applied` block is published.** The label
  is withheld and reported as Indeterminate rather than inferred from a nil
  error (3.4.1). This is a deliberate visibility loss, not a UI bug.

## 7. Rollback

Remove the script tag from `audio-settings.html` and the host tag from the six
page documents, remove the listener block from `menu-functions.js`, and restore
`_writeBackendSettings`'s silent return.
No data migration, no schema change.

## 8. Owner acceptance (2026-11-05)

**D1 — ACCEPTED — OWNER ACCEPTANCE WITH DOCUMENTED VERIFICATION LIMITATIONS.**

The owner accepted D1 on 2026-11-05 in its then-current implemented state. The
accepted state comprises:

- the iframe bridge and production-origin validation (§3.2.1);
- parent-owned activation locking and generation-qualified request handling
  (§3.0, §3.6);
- truthful Saved / Applied / Ready distinctions, including the restored-lock
  behavior (§3.4.1, §3.6.3);
- fail-closed handling of unresolved activation state (§3.6.3);
- the two production dismissal/save guards in
  `frontend/dist/js/audio-settings.js`;
- the corrected 16-test activation-lock harness exercising the real parent
  integration;
- the reported passing checks: 16/16 activation-lock tests, `node --check`, and
  `go build ./...`.

Evidence for that state is in `build/unit-c/d1-evidence.md` (§8, §9.1).

### 8.1 Verification limitations that remain

Acceptance does **not** establish, and the record must not be read as asserting:

1. The packaged-app restart boundary. The fresh-context test is a proxy only.
2. A rerun of the 14 bridge/lifecycle regression tests.
3. The remaining packaged-app verification and acceptance rows identified in the
evidence (plan tests 5, 6, 8, 9, 11, 12; acceptance rows 20–21).
4. Physical-device or real-audio behavior.

These items remain **unverified/deferred**. They are not converted to PASS, not
implied verified, and not evidence of a known implementation failure.

### 8.2 Closure boundary

This acceptance closes D1 at the owner's declared boundary. It does not
authorize further D1 implementation or verification work. D3 remains deferred.
