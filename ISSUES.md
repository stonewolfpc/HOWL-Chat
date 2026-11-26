# HOWL Chat — Working Issue Document

**Design law:** the app owns every fact it can verify; the model narrates.
The app is meant to be smarter than the model it hosts — that is the whole product.

**North star:** a world generator and inspiration tool for creative people.
Double-click, load a model, start making something. Zero setup knowledge required.
Fun matters. Control matters. Beginners get defaults; power users get depth, never code.

**Hardware floor:** assume 8GB VRAM, potato PCs, tiny context windows, models that
are dumber than the app. The app must be smarter than the model wherever it can be.
Anything the app can decide automatically, it decides automatically.

---

## Ordering law

1. **Broken** — something that does not work at all, or corrupts/loses work.
2. **UI** — works but is wrong, inconsistent, or unreadable.
3. **Additions** — new capability. Only after 1 and 2 are closed and the app is
   operational end to end.
4. **Ideas** — parked. Not scheduled. Only revisited if they fit the app and the
   foundations are solid.

**Stop rule:** one stage per session. A stage closes only with evidence — a
screenshot, a log line, a saved file, a count. No stage closes because the code
"looks right." At the end of every stage: answer the checkpoint questions before
opening the next one.

---

## Tier 0 — Broken (release blockers)

### T0.1 · My messages don't appear until restart
**Symptom:** I send a message, it never shows in the UI. Restart the app and it's
there. The reply renders fine.
**What we know:** the data is saved correctly — restart proves it. So this is a
render/state desync on the live path, not a storage bug. Do **not** clear caches
on load to "fix" it; that risks deleting real chats.
**Test first, before touching code:** send one message, reload, count how many
copies of that text exist in the saved chat JSON.
· 2 copies → duplicate send path.
· 1 copy → pure render bug.
Then: switch to another chat and back without restarting. If the bubble appears,
the send handler simply never re-renders.
**Evidence:** the JSON count, plus a screenshot of the bubble appearing same-session.

### T0.2 · The model receives my message twice
**Suspected separate bug** from T0.1 — likely the message is appended to the
prompt twice even though it's saved once. Fixing T0.1 may not fix this. Test both.
**Evidence:** a prompt dump or server log showing the message once per turn.

### T0.3 · Image model will not load
**Symptom:** black textless popup flashes (~0.2s), vanishes. No model load is
attempted. Loading animation runs, but no VRAM or CPU movement. Also: there is no
way to choose the image model directory at all.
**Likely cause:** stored absolute model paths still point at the old drive after
the E: → C: move, and the failure has no error surface so it dies silently.
**Fix in three parts:** (a) validate stored paths at startup and fail loudly;
(b) add a model directory chooser; (c) give that dialog a real error state instead
of a black box.
**Evidence:** invalid path → readable message. Valid path → VRAM climbs within 5s.

### T0.4 · Image preview is broken
Long-standing. The generated image does not preview. Find the cause and repair it.
This is the last blocker before the image page can be judged on layout (T2.3).
**Evidence:** prompt → image appears in the preview panel, correct aspect, no crop.

### T0.5 · Concurrent instances kill each other's model processes
`internal/proc/job_windows.go:110-136` (`ReapOrphans`) force-stops every
`llama-server.exe` / `sd-server.exe` whose `ExecutablePath` begins with the
project root. It consults no parent PID, no child list and no recorded owner, so a
second instance - or the development UI harness - takes the model child out from
under a running app. Observed 2026-11-01: harness startup reaped a live app's image
child (`[proc] reaped 1 orphaned model process(es)`); the victim identity (PID
24780 under PID 9252) is author-observed process inspection, not something the
saved log names. This is not Image-Studio-specific: any concurrently running Howl
instance or harness shares the same executable-path space.
**Before any repair:** a read-only audit of the ownership boundary - who may kill
what, and how ownership is proven. Repair is not in scope for the next stage.
**Evidence:** two instances running → starting the second leaves the first's child
alive and reports it as left alone, or an explicit ownership check refuses the
kill, with a test.

### T0.6 · The capture probe opens a stream on the ASR apply path
`app-media.go:16` documents the capture port as presence-level ("it must never
activate the endpoint or open a stream") but its body calls
`microphoneStatusNow()`, the initialize-level probe (`mic_status.go:68-74` →
`mic_windows.go:116`, `:165`): `Activate(IAudioClient)` +
`Initialize(SHARED, endpoint mix format)`, then release, no `Start`.
`internal/media/audio.go:171,188` runs that port twice per explicit ASR apply —
once before the model build and once after it.

Observed 2026-11-23, WH-1000XM4 (Bluetooth Classic, default output and input):
isolating that exact probe reproduced the reported interference twice (playback
louder, flat, bass gone, brief dropout, then recovery), while the presence probe
and two full application launches — one isolated, one real profile — produced no
endpoint change and no audible difference. A sampler using
`IMMNotificationClient` plus property-store polling (no `Activate`, no stream)
recorded zero transitions across both launches. Windows 11 selects HFP instead of
A2DP whenever an application opens the capture endpoint, and resamples other
applications' playback on that device to the HFP format (Microsoft, Bluetooth
Classic audio).

**Not part of this defect:** `GetAudioStatus` used the same initialize-level probe
on every page load and tab return until `67c1a37` (2026-11-23). Builds older than
that reproduce the launch symptom; HEAD does not.

**Before any repair:** decide whether ASR activation needs an initialize-level
probe at all, and make the port match its own contract. The presence probe
already answers "is there a usable device".
**Evidence:** a test that applies ASR settings with ASR enabled and asserts the
capture port is never invoked, plus a live apply on a Bluetooth headset with no
profile switch.

**Plan (awaiting approval, not implemented):** `captureProbe` returns
`microphonePresenceNow()` and its comment states presence semantics and why; the
`MicStateReady` copy stops claiming initialization that no longer happens; the
initialize-level probe stays as a diagnostics/test seam with no production
caller; a regression test replaces the initialize seam with one that fails the
test if called and asserts an ASR apply never reaches it.

**Semantics (settled 2026-11-23, before implementation).**
`IMMDeviceEnumerator::GetDefaultAudioEndpoint(eCapture, eConsole)` plus
`IMMDevice::GetState` proves exactly four things: a default console capture
endpoint is selected, the role resolved, `GetState` returned `S_OK`, and the
state is `DEVICE_STATE_ACTIVE` (present, enabled, not unplugged).
It does **not** prove that the endpoint can be opened, that its mix format is
accepted, that Windows privacy permissions allow capture, that no other client
holds it, that samples will flow, or that the renderer will pick the same device.
`DEVICE_STATE_ACTIVE` is a device-state flag, not "in use".

That claim surfaces as `MicrophoneStatus.Message`, which `newMicrophoneStatus`
sets from `microphoneMessageFor(state)` (`mic_status.go:140-148`) and which the
UI shows in `audio-settings.js:521` and `audio-service.js:122`. The precise
probe text already sits in `Detail` and is accurate (`mic_windows.go:356-362`).
The generic messages therefore have to describe the *gate*, not the device:

| State | Current message | Correction |
|---|---|---|
| `ready` | "initialized for shared-mode capture. Recording has not started." | "The Windows default microphone is present and enabled. Recording has not started or been verified." |
| `unavailable` | "...unavailable (disabled, disconnected, or in use)." | Drop "in use": presence cannot detect it. |
| `capture-init-failed` | "present but could not be initialized for recording." | "Windows could not confirm a default microphone; the probe failed before any capture was attempted." |
| `access-denied` | unchanged | Presence can return `E_ACCESSDENIED`; accurate as written. |
| `no-device` | `MicrophoneNoDeviceMessage` | **Do not reword** — frozen audit contract. |

The `MicStateReady` doc comment (`mic_status.go:8-13`) and `Usable()`
(`:53-56`) already describe both probe levels correctly; only `MicStateUnavailable`
(`:19-21`) carries the same "in use" overclaim and gets the same trim.

The regression test must fail when the initialize seam *runs*, not merely when
the status is wrong: replace `microphoneProbe` with a function that records an
invocation and calls `t.Fatalf`, drive `captureProbe`, `GetAudioStatus` and an
ASR-enabled `ApplyAudioSettings`, and assert zero invocations. Add the refusal
case (presence reports no-device → `CaptureUnavailableError`) so the gate stays
proven.

### T0.7 · Controlled Folder Access blocked a content-root write, and startup never noticed
Observed 2026-11-23 19:54:28, during the real-profile launch (`exp6`): Windows
Defender event **1123** — `C:\Users\StoneWolfPC\HOWL_Chat\build\bin\howl-chat-s1.exe`
was blocked from modifying `%userprofile%\Documents\Howl`. The block landed 5.6 s
after process start (19:54:22.72 → 19:54:28.33).

Correction to this entry's first reading: `Documents\Howl` is **not** the state
root. `paths.StateRoot()` is the project root unless isolation is active
(`internal/paths/isolation.go:185-190`), so settings, chats, logs and registries
live under `<root>/config` and were untouched. `Documents\Howl` is the
**content-library root** (`internal/paths/content.go:25-63`).

Startup still reported success — `[lifecycle] ready: 6 stages in 154ms`, webview
attached, rescan completed — so a successful startup does **not** prove the app
had normal access to its content directories, and nothing was surfaced to the user.

**Identified operation.** The blocked write was directory **creation**, not a file
write: `os.MkdirAll` at `app-worlds.go:47` inside `ensureDefaultLibrary`
(`app-worlds.go:35-55`), on the path from `paths.DefaultContentDir`
(`internal/paths/content.go:43-63`) = `<Documents known folder>\Howl\<Category>`.
Go implements that as `CreateDirectoryW` (`Documents\Howl` first, then the category
directory), so the denied leaf is `Worlds` or `Scenarios` — whichever the renderer
asked for first. Both take the identical path and no log line records the attempt.

It is reached from a **read** binding: `listDocuments` (`app-worlds.go:61-62`) ←
`GetWorlds` (`:78`) / `GetScenarios` (`:90`), driven by the renderer's startup load
(`world-repository.js:151` inside `_load`, `loadAll()` at module init `:268`;
`scenario-repository.js:151`). Also reachable from `:106` and `:193`.

**Proof of identity.** After a launch with protection off, `Documents\Howl`
contains exactly `Scenarios` and `Worlds`, both stamped `2026-11-23 20:25:50` —
3 s after that process started (20:25:47), the same window in which the protected
run was blocked. During the blocked run the directory did not exist at all.

**Required or optional:** required for content features, optional for boot. It
provides the approved save roots for Worlds and Scenarios (and Chats / Lorebooks
on first use). When it fails silently, `world-repository.js:154`
(`console.warn('[worldRepository] backend load failed')`) falls back to the stale
`localStorage` cache or an empty list, the library is never registered so a later
save has no default root, and every read retries the blocked creation. The user
loses their worlds/scenarios views and any default-location save without being
told why. Settings, chats and logs were unaffected — they live under the project
root, not Documents.

It is also a mutation on a read path: `GetWorlds` creates directories, contrary to
the rule stated at `app-libraries.go:48-50` that content directories are created by
explicit user action.

**Repro requirement:** with protection enabled, a launch must show the blocked
path in the Defender log *and* either write successfully or surface the failure.
Silent success is the defect.

**Verification plan (four cases plus 3b, injected through code seams, read back
from disk as well as the screen):**

| # | Case | Injection in a disposable isolated root | Expected result |
|---|---|---|---|
| 1 | Existing directory, readable but not writable | Pre-create `content/Worlds` holding one world document; deny write for the current user (`icacls <root>\content\Worlds /deny "%USERNAME%":(W)`), leave read access | Adopted and registered; the existing world is **listed and opens** (viewing works); library status `read-only` with a visible reason; new/save/delete/import disabled and explain why; nothing written (contents and mtimes unchanged); never shown as empty or unavailable |
| 2 | Missing directory, creation permitted | Fresh isolated root, nothing pre-created | Root and category created, probe file written/read back/removed, registered as default-save, library `ready`; Worlds panel lists and saves. Verify on disk: directories exist, no probe file left, `config/content-libraries.json` holds the registration |
| 3 | Directory creation denied | `content` exists as a **file** (ENOTDIR); repeat with `content` ACL-denied (`E_ACCESSDENIED`) | Core ready; Worlds library `unavailable`, naming the path and failed step; Scenarios untouched; other panels usable; one diagnostic; nothing created |
| 3b | One category denied, others permitted | Same as case 3, but only for Worlds; leave `content/Scenarios` absent and creatable | Worlds `unavailable` (creation denied) while Scenarios reaches `ready` and its panel loads and saves; Lorebooks and Chats unaffected; core ready. Proves no status is shared across categories |
| 4 | Optional logging unavailable | `config/logs` as a file, or ACL-denied | Core ready and usable; one "log unavailable" diagnostic; no library marked unavailable |

Unit level, behind the injected filesystem seam: each of the seven probe steps
fails independently and reports the right `Step` and classification; a probe
against an existing writable directory leaves no artifacts; and a readable but
unwritable library returns its documents with `Readable=true, Writable=false`
while every mutation binding returns the typed `read-only` refusal.

Both hosts are exercised: packaged desktop (`build/bin/howl-chat-s1.exe`) and the
browser dev host (`-tags uiharness`), each on its own isolated root. The folder
picker is a native surface, so "Choose folder" is verified in the desktop app only
and that difference is recorded rather than treated as equivalent.

**CFA verdict:** still optional. Re-enabling Controlled Folder Access is useful only
as confirmation that event 1123 reproduces this exact path; because
`Documents\Howl` now exists, such a run would have to target a fresh category
directory (e.g. Lorebooks) rather than the parent. The four deterministic isolated
root cases above are the repeatable gate.

**Readiness vs. this write — the call order**
1. Process start → `NewApp` (`app.go:119`) → `Lifecycle.Boot(app.bootStages())`.
2. Six stages (`app-boot.go:38-43`): core services, audio owner, settings store,
   registries, image domain, inference capture. **None creates or validates a
   user-content directory.**
3. `Lifecycle.Boot` sets `ready = firstErr == nil` (`internal/host/lifecycle.go:110`)
   and prints `[lifecycle] ready: 6 stages` at ≈+154 ms.
4. `startup()` prints `webview attached … ready` (`app.go:137`), then starts
   background rescan and model discovery.
5. Renderer loads → `GetWorlds` → `ensureDefaultLibrary` → `MkdirAll` → denied
   (≈+5.6 s) → the error returns to JS at `world-repository.js:151` and dies in the
   `console.warn` at `:154`.

So readiness is declared before the write is even attempted, and the write's error
never travels toward it. The lifecycle *can* fail a boot (`ready=false`,
`[lifecycle] startup stopped at "stage"`, `lifecycle.go:110-115`), but
`app.go:120-121` only prints `[Lifecycle] startup incomplete` and the window opens
anyway; `Ready()` is consumed solely by the stdout line at `app.go:137`.

**Verdict: a required-write failure is structurally unreachable by the user.** It
is swallowed twice — the boot gate does not check writability, and the content
path's error ends in a console warning. Worth stating plainly: `Lifecycle.Ready()`
has exactly one consumer today, the stdout line at `app.go:137`. Global readiness
gates nothing. A boot stage that failed on a content library would therefore be
cosmetic and still invisible, so it is **not** the repair.

#### Revised repair design (after review) — plan only, not implemented

**R1 · Feature-level readiness; the core stays ready.**
No boot stage fails on content storage. `Lifecycle.Ready()` keeps its current
meaning (core composition succeeded) and is not reused as a content gate.

The dependency is demonstrated, not assumed. The content-library root is read only
by: the world/scenario bindings (`app-worlds.go`), cutover rollback
(`app-cutover.go:144,152`), the chat service's canonical world/scenario/lore load
(`internal/backend/chat/service.go:107-129` and `loadCanonicalWorld`), which
already returns `ContextUnavailable` and contributes nothing to the prompt, and
content import (`internal/contentimport/service.go:34`). None of these runs in a
boot stage and none is needed for settings, chats, audio or image work, so core
readiness has no dependency on this root.

Per-library state is one authority in the app layer, behind one binding. **The
product guarantee it exists to enforce: a storage failure must never masquerade as
an empty library, a successful save, or a globally broken application.**

Readiness is reported **per content library**, never once for the content system.
Every entry carries its category and library identity, and each library is
evaluated on its own path and probe result, so Worlds being unavailable cannot
imply anything about Scenarios, Lorebooks, Chats or Characters — they are separate
registrations.

```
type ContentLibraryState struct {
    Category  string // worlds | scenarios | lorebooks | chats | characters
    LibraryID string // registration ID; empty only on the synthesized not-initialized entry
    Path      string
    Readable  bool   // content can be listed and opened
    Writable  bool   // a mutation can succeed (write probe passed)
    SaveRoot  bool   // designated default save destination
    Status    string // ready | read-only | not-readable | not-initialized | unavailable
    Reason    string // user-facing, e.g. "this library is read-only: access denied"
    Operation string // failed provisioning/probe step, when there was one
    Retryable bool
}
func (a *App) GetContentReadiness() []ContentLibraryState
```

Reading and modifying are separate capabilities, so there is deliberately **no
single `Ready` flag**:

| Status | Readable | Writable | Meaning, and what the UI does |
|---|---|---|---|
| `ready` | yes | yes | Normal operation |
| `read-only` | yes | **no** | Existing worlds/scenarios stay visible and openable; create, save, delete and import are disabled with the reason shown. Never hidden, never labelled unavailable |
| `not-readable` | no | no | Registered but access fails (permission, missing drive, offline). Listed with the reason; cached items shown as stale |
| `not-initialized` | no | no | No usable library yet. Offer initialization |
| `unavailable` | no | no | Initialization or access failed for another reason; name the failed step |

Category rollup is derived and never shared: a category is readable when any of
its libraries is readable, and writable when any readable library with a write
role passes its probe. A category with no registration produces exactly one
synthesized `not-initialized` entry. No status is copied between categories.

The chat service's `GetContextStatus()` (`service.go:132-136`) is the existing
precedent: state is reported, unavailability is displayed, and no data is
substituted.

**R2 · Provisioning is its own authority; reads never mutate.**
A new `app-content-roots.go` owns the only provisioning route: resolve
`paths.DefaultContentDir`, create the root, create the category directory, prove
writability, register it as default-save. `ensureDefaultLibrary` is removed from
`listDocuments` (`app-worlds.go:61-62`), `saveDocument` (`:106`) and
`importDocuments` (`:193`).

- Reads become pure. With no registered root, `GetWorlds`/`GetScenarios` return an
  empty list **and** feature state `not-initialized` — not an error, so the panel
  can offer setup instead of presenting a backend failure.
- A write into an uninitialized category is refused with a typed error naming the
  missing root; the renderer turns that into the "Set up the default location"
  action. A save never creates a directory, so a CFA block can no longer hide
  inside a save.
- `AddDefaultContentLibrary` (`app-libraries.go:46-58`) keeps its explicit
  user-action meaning and delegates to the same authority. This removes the second
  provisioner rather than adding a layer.
- Initialization runs once, **synchronously from `startup()` before the background
  scans**, so the renderer's first `GetWorlds` sees a settled state instead of
  racing provisioning. It never fails the boot; each category's outcome lands in
  `ContentReadiness` and is logged once.
- Existing but unregistered default directories are adopted by that step
  (registration is not a write). A directory that is readable but unwritable is
  still registered with `Writable=false`: its contents remain visible and
  openable, and every mutation is refused with a typed `read-only` error rather
  than the feature being called unavailable.

**R3 · The writeability probe proves what it claims.**
`ProbeWritable` reports each step separately rather than one boolean:

```
type ProvisionReport struct {
    RootCreated     bool // false when the root already existed
    CategoryCreated bool
    FileCreated     bool
    BytesWritten    int
    ReadBack        bool
    Cleanup         bool
    Step            string // the step that failed, "" on success
    Err             error
}
```

Steps: resolve path → create root → create category directory → create a randomly
named probe file → write N bytes, `Sync`, close → read back and compare → remove
the probe file → remove the category directory **only** if this call created it and
it is empty. A pre-existing directory is never removed.

Errors are classified, not stringified: `permission-denied` (`os.ErrPermission` /
`ERROR_ACCESS_DENIED`), `path-blocked` (a regular file where the directory
belongs), `read-only` volume, `not-initialized`, `unknown`. A directory that
already exists but cannot be written is a first-class case: reads succeed and
return its documents, the probe records `Writable=false`, mutations are refused
with the `read-only` reason, and nothing claims success.

The probe only ever targets a disposable isolated root, and its filesystem calls
sit behind an injectable seam so each of the seven steps fails deterministically in
a unit test without ACLs.

**R4 · Every refusal is surfaced, and viewing is never confused with modifying.**
- A banner per **library** (category and path named) driven by its status:
  `read-only` says the content can be viewed but not changed and gives the reason;
  `not-readable` and `unavailable` give the reason and failed step;
  `not-initialized` offers "Set up default location" (provisioning), and "Choose
  folder" stays available (existing picker).
- Mutating controls — new, save, delete, import — are disabled in a `read-only`
  library and explain why if used. Viewing and opening stay enabled.
- One diagnostic per initialization attempt, not one per read. Today
  `a.logErrorSilent` plus the renderer `console.warn` make every retry invisible.
- Cached or legacy renderer data is marked stale only when its library is not
  readable. A `read-only` library's documents are live and must not be marked
  stale. Today `world-repository.js:154-176` presents the cache as if it were live.

**Unchanged:** `worldstore` envelopes, the registry's shape (`library.View`
already exposes `Exists/Readable/Writable/Available`), and the existing
`Documents\Howl` directories.

**Implementation order** (each slice builds, vets and keeps the current suites
green, and its workflow is driven in-host before the next slice starts):
1. `app-content-roots.go`: provisioning authority, step-wise probe, error
   classification, injected filesystem seam. Unit tests fail each probe step.
2. Remove provisioning from reads/saves/imports; add the typed "not initialized"
   refusal and `GetContentReadiness`; root tests assert that `GetWorlds` never
   creates a directory (probe the isolated root before and after).
3. Call initialization from `startup()` before the background scans; record state;
   one diagnostic per attempt.
4. Renderer: per-library banner keyed by status, `read-only` disables mutations
   while keeping content visible and openable, route save refusals to the banner,
   and mark cached worlds/scenarios stale only when their library is not readable.
5. Four cases plus 3b in both hosts, read back from disk and screen, recorded here
   with the observed results.

#### Implementation and verification (2026-11-24)

Implemented as designed. Files changed:
- `app-content-roots.go` (new): the single provisioning authority, step-wise
  write probe (`RootCreated/CategoryCreated/FileCreated/BytesWritten/ReadBack/
  Cleanup/RootRemoved/CategoryRemoved`, failing `Step`), injected filesystem seam,
  classified errors, per-library `ContentLibraryState`, derived category rollup,
  `GetContentReadiness`, `EnsureDefaultContentLibrary`, and the create gate.
- `app-worlds.go`: `ensureDefaultLibrary` removed. `GetWorlds`/`GetScenarios` are
  pure reads; creates (`saveDocument`, `importDocuments`) pass through
  `requireCreateWritable`.
- `app-libraries.go`: `AddDefaultContentLibrary` and `AddContentLibrary` delegate to
  the one authority; a user-added read-only root is registered, not rejected.
- `app.go`, `harness_main.go`: content provisioning/probing runs once before the
  renderer can read it, in both hosts; the optional log is verified first.
- `startup_log.go`, `app-errors.go`: an unwritable startup log is reported as a
  medium `SYS-LOGGING` diagnostic; diagnostics fall back to `paths.TempDir()` so
  the limitation is reportable instead of vanishing with the logger.
- `frontend/dist/js/content-readiness.js` (new) + `world-repository.js`,
  `scenario-repository.js`, `world.html`, `scenario.html`: per-library banner with
  path/status/reason and a setup action, mutation controls disabled when not
  writable, cache marked stale only when the library is not readable.

**Observed results** (isolated roots, browser dev host, read from disk as well as
screen):
1. **Readable, not writable** (existing `Worlds` with one world; deny
   `WD,AD,WEA,WA` for the current user): library `read-only`, `readable=true`,
   `writable=false`; the world still lists; create refused with the read-only
   reason and the directory unchanged; banner shows path + read-only + reason;
   Create and Save disabled. Not shown as empty or unavailable.
2. **Missing, creation permitted**: startup created `content/{Worlds,Scenarios}`
   and registered both with `default_save`; a world created through the panel UI
   landed as `Worlds/<id>/world.json`; banner hidden, Create enabled.
3. **Creation denied**: `content/Worlds` as a file → `worlds unavailable`
   (retryable false, path named), `scenarios ready`, its directory created; world
   page showed the reason plus the setup action with Create disabled; scenario page
   had no banner and Create enabled. Removing the blocker and clicking “Set up
   default location” provisioned and registered Worlds and cleared the banner.
4. **Optional logging unavailable** (`config/logs` as a file): Worlds and
   Scenarios stayed `ready` and writable and pages served. The earlier run
   recorded **two** `SYS-LOGGING` diagnostics for one condition (diagnostics
   fallback and the startup-log open); after reconciling, the fallback is now the
   single report (`logging.initialize`) and a re-run on a fresh root recorded
   **exactly one** medium `SYS-LOGGING` entry, naming the configured path and the
   OS error, with the fallback jsonl still written under `temp/howl-diagnostics`.

**Packaged desktop host — UI driven, not just started.** The production build was
launched on isolated roots and driven with real mouse input and WebView2 windows,
reading both the rendered window and the files on disk:
- Fresh root: World page had no banner and an enabled Create; the creation form
  accepted a typed title and wrote `Worlds/<id>/world.json`, confirmed by the
  app's own save dialog and by reading the file.
- `Worlds` denied (`WD,AD,WEA,WA`): banner `Worlds: 1 of 1 libraries needs
  attention` naming the path and reason, existing world still shown, Create
  visibly disabled; clicking the disabled control opened nothing and changed
  nothing on disk.
- `Worlds` replaced by a file: banner reported `unavailable` (file, not a
  directory) with the setup action and Create disabled; the Scenarios page had no
  banner and an enabled Create — one category's failure did not touch the other.
  Deleting the blocker and clicking “Set up default location” provisioned and
  registered `Worlds`, cleared the banner and re-enabled Create; the restored
  world then listed again with no banner.
- Two libraries: banner `1 of 2 libraries needs attention` listed both paths
  (default `read-only`, sibling `ready`); Create reflected the default save root
  and the sibling's world stayed visible and editable.

**Multi-library audit (two worlds libraries).** `saveWritable` was added to the
category rollup: `writable` stays true when any library can be written (editing
survives a failing sibling), while create follows the **default save root**. The
first pass exposed a real defect: `saveDocument` gated every save on the default
root, so a read-only default refused updates to a document that lived in a
healthy sibling (`SaveWorld FAILED` with the read-only reason). The gate now
applies only to creation (no id, or proven absence). Verified both directions
with a writable sibling owning a world while the default was read-only: the edit
saved (`Sibling World Edited`) and the read-only root was untouched; with the
default writable and the sibling read-only, create stayed enabled and the banner
counted one failure without hiding it. The probe only proves that *new* files
cannot be created there (a directory-level deny leaves existing files writable),
so the banner and reason now say “new items refused” instead of claiming all
changes are refused.

**Tests:** `go build ./...`, `go vet ./...`, `go test -count=1 .` all pass, including
new coverage that fails each probe step independently, keeps a read-only library
readable, counts a failing library in the rollup of two, proves a read-only
sibling does not disable create when the default save root is writable, proves a
read-only default refuses create while a healthy sibling's update still saves,
isolates categories, proves `GetWorlds` creates nothing, and lets the probe
override a registry view it disproves.

**Limitations:** no CFA confirmation run (still optional). The folder picker is a
native surface and was not exercised (the setup action used the existing default
location). Form save success surfaces as a native JS alert in the WebView2
window, which blocks scripted reads until dismissed; outcomes were confirmed from
the backend log and the files on disk instead. During verification the machine
showed Controlled Folder Access active and blocked `python.exe` writing under
`Documents\Codex` — an unrelated application and folder. No Defender setting or
exclusion was changed, and no test touched `Documents\Howl`.

**Status:** implemented and verified in both hosts, including packaged desktop UI
interactions (banner, disabled create, category isolation, setup action, restore)
and the multi-library rollup in both directions. Available for review before any
commit.

**CHECKPOINT 0** — answer before Tier 1:
- Did the JSON show 1 copy or 2?
- Does the bubble render on chat-switch without restart?
- Which drive/folder holds your image models now?

---

## Tier 1 — Core mechanics (works, but wrong or forgotten)

### T1.1 · State is not owned properly
Changing tabs must **not** unload the model or drop the connection. Coming back
should find everything where you left it. Settings — samplers especially — must
persist visually across exit and reload, not just internally.
**Evidence:** load model → switch tabs → return → still loaded, still answering.
Change samplers → restart → sliders sit where you left them.

### T1.2 · VRAM arbitration between chat and image models
Agreed behaviour:
- Image model requested while a chat model is loaded → check whether both fit.
  If not, a confirm dialog: this will unload your chat model to make room.
- On confirm: unload chat, re-check freed VRAM, then load image.
- Leaving the image page does **not** unload or disconnect.
- Going back to chat and pressing Load mirrors it: check image model, confirm,
  unload, load. Same dialog, wording swapped.
**Evidence:** both directions exercised, with VRAM readings before and after.

### T1.3 · Think tags / reasoning handling
Match how SillyTavern and Kobold handle think modes: on, or off, or hidden, and
selectable intensity where the model supports it. Hidden means it must not count
against the context budget. A thought must never be fed back to the model on the
next turn — anything inside a think tag is discarded from history.
**Evidence:** hidden mode shows no tokens in the UI and shortens the measured
prompt; a two-turn test proves no think block reappears in turn two.

### T1.4 · Sampler defaults the user never has to touch
Defaults stay "suggested general," but vary by model family — they differ too much
to force one set. Supported models must get correct prompt formatting
automatically. No card in the UI that the backend ignores.
**Evidence:** load two different families, confirm each gets its own input/output
format with no user action.

### T1.5 · Header says the wrong name
"Howl Chat - Tron Edition" → "Howl Chat". One-line fix, bundle it with T1.1.

**CHECKPOINT 1:**
- Which tab switch currently breaks a loaded model, if any?
- Does the image model survive leaving its page today?
- Do sampler sliders survive a restart right now, or reset?

---

## Tier 2 — UI (works, looks wrong or reads badly)

### T2.1 · Theme unification
Unify dark and light palettes across every page to match the chat page. Image page
is already consistent — don't disturb it. Light theme is the real problem: header
text stays white in many places instead of adapting, which suggests an override
left behind by the earlier coder. Hunt the overrides rather than patching each
symptom. Also verify what a system-light-theme default should look like: readable,
low strain.
**Evidence:** side-by-side screenshots of each page in both themes, no white-on-white
or white-on-light text anywhere.

### T2.2 · Ghost elements and stale markup
Class of bug: markup left behind by deleted features, and duplicate element IDs
where one handler writes to an ID that no longer exists while the real element sits
there working. We've already hit one of these (`ig-negative-prompt` vs `ig-negative`).
Sweep every element lookup target against the actual DOM.
**Evidence:** zero write-to-missing-ID, zero duplicate IDs, no orphaned markup.

**Status 2026-11-02: inventory complete, stage OPEN, nothing repaired.** Read-only
classification of every lookup in `frontend/dist/js` against the nine pages.
Full record in `STAGE_STATUS.md` → "T2.2 · Ghost-ID inventory (evidence revision)".
The inventory counts *references*, not read-versus-write operations, so this stage's
own acceptance bar — zero write-to-missing-ID — **cannot yet be scored**; 30 reference
sites remain, and they are not claims that all 30 are writes or runtime failures.
Nothing here authorises a repair.

**Reconciled baseline (this working copy at `1477530`, source-derived; no page loaded):**
372 ids named by literal `getElementById`, 601 ids in static markup, 325 provided
by markup, **47 absent**, 29 of those created by script before use, **18 never
created anywhere (ghosts)** over **30 reference sites**. Distinct creation ids 50
(21 never looked up); selector-based lookups 33 raw / 31 distinct / 27 leading
tokens. Identities: 325 + 47 = 372, 29 + 18 = 47.

**Superseded figures — do not reuse.** An earlier pass reported 75 absent ids as
18 ghosts + 51 runtime-created + 4 dynamic prefixes + 10 compound fragments,
which sums to 83; the same pass reported **400 lookup ids** where the corrected
count is 372. All five figures — **75, 400, 51, 4, 10** — are withdrawn and must
not be reused as the corrected baseline: 51 exceeds the entire created-id set (50),
so it cannot be a subset count; the 4 prefixes and 10 fragments were selector
strings and template prefixes, not ids; 75 counted distinct *strings*, not distinct
*ids*; and 400 is superseded by the regenerated 372. None of the old figures is
reproducible — the original regex was not retained — so they are not reconciled
against; they are replaced by the measurements above (325 provided + 47 absent =
372; 29 created-before-use + 18 never created = 47).

**Candidate findings (the audio handoff is repaired in code; the other three remain unaddressed):**

1. **Audio microphone handoff — repaired in the microphone/history tranche
   (2026-11-04).** The renderer no longer uses `MediaRecorder` or a `blob:` URL.
   It captures mono PCM through Web Audio, encodes a 16 kHz PCM16 WAV locally
   (`audio-wav.js`), and calls `TranscribeAudioData` with one base64 payload. The
   Go binding writes a uniquely named `%TEMP%/howl-asr/recording-*.wav`, hands
   that path only to the managed pipeline, and removes it on success, error,
   cancellation, or panic. The composer target is now `#message-input`; the
   transcript is appended into the composer, not sent automatically. Go tests
   cover binding/pipeline handoff, temp cleanup on success/error/cancellation,
   concurrent unique paths, and the frontend contract; a Node assertion covers
   the WAV header and 48 kHz-to-16 kHz output. **Limits:** no live Wails
   microphone run, no persisted-history reopen test, no concurrent chat + ASR
   measurement, and no general transcription-quality claim yet.
2. **Two write-only character properties.** `first_mes` and `scenario` are collected
at `character-form.js:61-62` (always `''`, missing fields) and restored at
`character-grid.js:127` (guarded no-op). Neither has a consumer anywhere in the
frontend. Recorded, not deleted, not remapped: whether they belong to the intended
character schema is a product decision.
3. **Three manager modules have no discovered load path** (`character-manager.js`,
`scenario-manager.js`, `world-manager.js`) — no script tag, no import, no dynamic
loader, no worker, no cross-file reference, no inline handler, and no Go-side
reference, though `//go:embed all:frontend/dist` packages them into the binary.
Their element contracts do not match the live forms and they contain selector-modal
code with no live counterpart found by name. **Not redundant by proof and not safe
to delete**; behavioural equivalence is unestablished.
4. **Three parallel form/manager pairs** (`character-form`/`character-manager`,
`scenario-form`/`scenario-manager`, `world-form`/`world-manager`) share function
names while differing in structure. Candidate structural finding only.

**Carried limitations:** 16 template-literal lookups are outside the inventory;
the selector-token figure has one valid extraction method (an earlier attempt
returned 0 because its `sed` failed — do not reuse it); creation detection covers
`id="X"` and `.id = 'X'` only; all reachability is static inference, not observation.

### T2.3 · Image page layout — target Invoke AI organisation
Rework toward Invoke AI's structure: preview prominent, controls grouped, history
organised. Reference it (look it up) but don't chase parity — take the parts that
make it powerful and organised, drop the rest. Depends on T0.4.

**CHECKPOINT 2:**
- Which pages look worst in light theme to you?
- Any specific Invoke AI screen you want to use as the reference?

---

## Tier 3 — Additions (only once the app is operational)

Nothing here starts until Tiers 0–2 are closed.

### T3.1 · Context handling — verify before optimising
**First question is not "how do we optimise," it's "was the context flush ever
implemented?"** Verify that before assuming it exists. If it exists, measure the
1–2s prefill cost. Then decide whether restructuring is worth it at all.
Candidate cheap win to test: `cache_prompt: true` plus `--cache-reuse N`, which
reuses shifted KV chunks instead of re-prefilling. Full model unloads are **not**
acceptable — load times are too long.
Research before implementing, so we don't optimise into something the models
tolerate badly.

### T3.2 · Load speed — learn from LM Studio
LM Studio loads a model in under a quarter of our time. Find out how and adapt
whatever is portable. This may be the single biggest quality-of-life win in the app.

### T3.3 · Scoped lorebook + memory system
Lore must be separable: world, scenario, character. Created as a lorebook, assigned
to an entity, and it appears in that entity's list at the top of its page. Character
lore fires **only** for that character.
Mechanics to implement: primary + secondary keys, word-boundary matching, scan
depth 2–3, recursion capped at 2 with a per-entry no-recursion flag, a fixed token
budget with priority ordering, and insertion order.
Precedence on a turn: character → scenario → world; character wins on conflict.
Dedupe identical bodies across scopes, prefer the narrower scope. Give lore its own
budget ceiling (~15–20% of context) so it can't starve the character card.
Memory: per character, tied to who said it, natural-language triggers, and a hard
guard that it can never blow the context window.
**AI-authored lore:** generate entries from a world/scenario/character description.
Keys get a specificity score (reject common single nouns, prefer proper nouns and
multi-word phrases), plus a frequency test against the existing chat log — a key
that already appears everywhere is flagged as "fires too often," with the count.

### T3.4 · Per-field AI character sheets
Every field gets its own generate button. The model sees the fields already filled
and must make the new one line up: name, kind, the target field's purpose, plus
truncated versions of siblings as reference, and returns only that field.
Draft + accept / regenerate / lock per field. Grammar-constrained JSON, one
pipeline with per-kind schemas.
No preset voices. Exemplar dialogue is authored per character — exemplar lines beat
trait lists for keeping a voice from bleeding.

### T3.5 · Multi-character chat
2+ characters kept fully separate, each receiving the appropriate system prompt,
with context flushed properly between speakers. Root cause of voice bleed is flat
context, not sampling — needs structured transcript + per-character state.

### T3.6 · Ternary models as first-class
Load ternary formats without new dependencies. Must just work when the user points
at the file.

### T3.7 · Audio and TTS end to end
Test the audio model capability and the TTS service through the UI. Needs headphones
present for verification.

### T3.8 · Audio generation
Newer small audio model, reportedly fits in 8GB for short clips. Short clips are
fine. If it works, play the generated audio inside the app.

### T3.9 · Diffusion version check
Confirm the current Stable Diffusion release; update to the latest patched version
if one exists.

### T3.10 · LLM writes the image prompt — research gated
Can a local LLM prompt the image generator well? Would it be worth the cost?
Proposed flow if yes: compose the prompt before swapping models; optionally keep
enough of the LLM resident to inspect the result (only if vision-capable) and refine.
Every part of this is a big if — cost, quality, speed. Research before building.

### T3.11 · Compatibility for everyone else
Auto context and layering currently assume this machine. Generalise as far as
possible and ship as beta, then fix from user reports. Goal: if it runs AI at all,
it should be usable.

### T3.12 · Whole-app optimisation pass
Separate pass, its own stage, after everything above. Look for waste everywhere.

**CHECKPOINT 3:**
- Confirm the flush was/wasn't implemented before I start T3.1.
- Which addition matters most to you after the app is solid?

---

## Tier 4 — Parked ideas (no stage)

- **Reactive 3D character.** Generated avatar that animates emotionally and moves
  its mouth while speaking. Enormous VRAM cost; likely incompatible with multi-character
  chat (too many bodies). Fallback would be a 2D still, but only if animation can be
  fully automated and prompt-driven — no Live2D hand-animation. Attractive and
  probably out of scope. Revisit only if the foundations are strong and resources allow.
- **Multi-character with avatars.** Sub-case of the above; same resource problem.

---

## Research queue (do before committing to any of it)

| Topic | Question |
|---|---|
| LM Studio load speed | How does it load so fast, and what's portable to us? |
| Context flush | Was it ever implemented? What does it actually cost per turn? |
| KV cache reuse | Does `--cache-reuse` help on this hardware, measurably? |
| Think modes | How do SillyTavern/Kobold expose levels, and what do models accept? |
| Ternary formats | Which quantiser/runtime path, with no new dependencies? |
| Invoke AI | Which layout decisions are worth copying? |
| Stable Diffusion | Current release, and what updating costs us. |
| LLM → image prompt | Is the quality worth the swap cost? |
| Audio model | Exact model, VRAM, clip length on 8GB. |

---

## Open questions for you

1. Which chat-switch or tab-switch currently loses a model? I need the actual
   reproduction, not the general rule.
2. Where do your image models live now — the new drive path?
3. Light theme: is there a specific page that's the worst offender?
4. Beta scope: ship beta right after Tier 0–2, or after T3.5 (multi-character)?

---

## Decision log

Append here as stages close. One line each: date, stage, evidence, result.

- 2026-11-23 · Audio interference investigation (research only, no code change) — reproduced the reported startup interference with one operation: Howl's initialize-level capture probe (`Activate(IAudioClient)` + `Initialize` on the default console capture endpoint) flipped the WH-1000XM4 to HFP and was stage-reported twice. The presence probe, an isolated packaged launch and a real-profile packaged launch all produced no endpoint change (render stayed `Headphones (WH-1000XM4)` 48 kHz/2ch, zero `IMMNotificationClient` events, zero property-store changes) and no audible difference. `GetAudioStatus` called the initialize probe on every page load until `67c1a37`, so pre-23rd builds still show it at launch; HEAD does not. Residual defect recorded as T0.6. Deferred by decision: render-side `GetMixFormat` and session-volume instrumentation (each needs an `Activate`). Evidence: `build/tmp/audio-sampler/logs/exp1_startup.log`, `exp2_presence.log`, `exp3_initialize.log`, `exp4_a2dp_baseline_20261123_195152.log`, `exp5_initialize.log`, `exp6_real_profile.log`.

- 2026-11-04 · Microphone/history tranche (code close-out, live run pending) — replaced the WebM `MediaRecorder`/`blob:` handoff with renderer-side 16 kHz PCM16 WAV capture and a `TranscribeAudioData` binding; temp files are unique and removed on success, error, cancellation, and panic; composer target corrected to `#message-input`. Evidence: root tests, full `internal/audio/...` tests, `go build ./...`, JS syntax, Node WAV header/resample assertions. Not yet evidence: live Wails mic capture, transcript send + conversation reload/history, concurrent chat + ASR, real Whisper. No registry or family-list changes.
- 2026-11-02 · Image generation settlement guard — source review PASS; stage-recorded tests PASS; independent review PASS; real-app acceptance PASS on verified build `d7e8aacd…` (PID 9252): abort leaves child alive, next request refused as abandoned, Unload/Reload recovers, work admitted after reset. Closed with no code change; the earlier live run had tested a pre-guard 18:59 binary. Unresolved: why the old sd-server exited.
- 2026-11-02 · T2.2 ghost-ID inventory (evidence revision, no code change) — 18 ghost ids over 30 sites confirmed against a regenerated accounting: 372 literal `getElementById` ids, 601 markup ids, 325 provided by markup, 47 absent, 29 created before use, 18 never created. Withdrawn as superseded: 75 absent ids / 400 lookup ids / 51 runtime-created / 4 prefixes / 10 compound fragments (summed to 83; 51 exceeded the 50-id creation set; 75 counted strings, not ids). Recorded, none repaired: the audio `chat-input` vs `message-input` mismatch is conditional on an upstream `blob:` URL handoff (`saveBlobToTempFile` returns `createObjectURL`), so it is not yet a one-line fix; character `first_mes`/`scenario` are persisted write-only empty strings with no consumer; three manager modules have no discovered load path but unestablished behavioural equivalence, so they are legacy-looking, not provably redundant; the three form/manager pairs are parallel implementations, not duplicates. Limits: 16 template-literal lookups unclassified, one valid selector-extraction method, static inference only. Checkpoint `1477530` untouched.
- 2026-11-01 · T2.1 Stage C (remaining pages theme verification) — five pages measured in both themes through `setThemeMode` with transitions off: 0 offenders at the probe's `|dL| * 255 < 60` luminance heuristic (text >=16 px only; the 14 px button labels were assessed visually), 0 stuck-dark areas in light, expected dark panels only in dark, and the buttons checked render white. No code changed. Independent review: PASS WITH FINDINGS, no blocking findings; closed. Corrections recorded: `real-settings.html` no longer exists (deleted by T2.2); `--accent-secondary` is deliberately constant cyan in both modes; the "stuck light flyout" reading was retracted, its proposed 4 s-transition cause withdrawn as unsupported and contradicted (no `@property` registration, probe overrides transitions), cause Unestablished; `ReapOrphans` kills by executable path under the project root, not by recorded state or parent identity. Evidence limits: per-page results and screenshots were not persisted, and the harness wrote chat state through real bindings. Open: override-sheet deletion, top-bar collision, `audio-settings` dark body background, stale `real-settings` references, and the cross-cutting process-ownership defect.


### T2.4 · Speaker identity in chat (thumbnail or initial)
**Gap:** names alone are too slow to parse in a hurry — the eye locks onto an
image instantly but has to read a name. This is a legibility item, not decoration.

Spec:
- Thumbnail of the character's card image beside the speaker name on every message.
- No image → enlarged first initial of the character's name, light text on a dark
  blue field.
- Fixed size (small enough not to crowd the bubble, large enough to recognise at a
  glance), consistent across themes — the dark blue must not fight light theme or
  vanish into the bubble.
- Serve a resized thumbnail, not the full card image, so a large portrait isn't
  loaded once per message.

Why it's early: this is a **prerequisite for multi-character legibility** (T3.5).
Two speakers distinguished only by name is exactly the failure mode it fixes.
Bundle with T2.1 so the fallback colour is decided once for both themes.

**Evidence:** screenshot of a two-character exchange, one with an image and one
without, both instantly distinguishable; thumbnail file size confirmed small.

### Stage closure ritual (applies to every stage)
When a stage's fixes are in, don't open the next one:
1. Walk the whole app, not just the stage's area — the way you did when building
   this document, one thing at a time.
2. Anything new goes in this document, appended to the correct tier. It is **not**
   fixed on the spot, and it does not change the current stage's scope.
3. Compare the fresh pass against the document: what's closed, what's still open,
   what's new.
4. Only then choose the next stage.

This is the loop that produced the good list today. Protect it.

### T2.5 · Character / world / scenario image sets
**Existing capability (regressed):** a card can hold several images (up to ~10)
with one marked as default. Worlds and scenarios have the same capability, with the
image shown in the right-hand panel. This worked at some point. Find out whether it
was scrapped, overwritten, or simply broke.

**Investigation order — cheapest explanation first:**
1. **Data layer.** Can a card still store a list of images plus a default, or only
   a single path? If storage lost the list, the feature went at the data layer and
   the UI was following it down. That's "scrapped."
2. **Leftover symbols.** Search for the old names: `cycle`, `slideshow`, `gallery`,
   `interval`, `images[`, `setInterval` in the character/world/scenario code. A
   surviving handler with no DOM target is the "overwritten" signature — the same
   pattern as the `ig-negative-prompt` ghost.
3. **Lifecycle.** If it's timer-driven, check whether anything clears the timer on
   tab switch, chat switch, or re-render. A timer killed by a later lifecycle change
   explains "it used to work" precisely.
4. **Overrides.** If the earlier coder's UI overrides removed the container or its
   styles, this is the same hunt as T2.2. Run them in one session.

**Spec once repaired:**
- One shared image-set mechanism bound in three places (chat speaker, world panel,
  scenario panel). Do not write three timers — this is the same component three times.
- Default image is the anchor: it shows on load and when cycling is off.
- Optional cycling, off by default, with a sane interval (a few minutes, not seconds).
  Configurable duration; per-character override later if it earns it.
- Never cycle while the character isn't speaking in a group chat — wasted decoding.
- Preload the next image before swapping so there's no blank flash, and keep the
  container size fixed so the layout doesn't jump between aspect ratios.
- Clear the interval on unmount and on chat/tab change. Stacked intervals are the
  classic cause of "it speeds up the longer I leave it running."
- Respect reduced-motion and give it a visible off switch.
- Avatars (T2.4) should draw the default image from this same set. One source of truth.

**Evidence:** a card with 4 images, default shows first, cycling advances on the
timer, no flash, no layout jump, timer stops on tab switch and does not stack;
world and scenario panels behave identically.

### T3.13 · D&D / tabletop mode
**Gate:** requires T3.3 (typed state + extraction) and T3.5 (multi-character).
Not startable before those, because its whole problem is theirs.

**Why it keeps failing:** local models are bad at being the source of truth. They
drop clothing, invent inventory, lose HP, contradict the map. No prompt fixes that.
It's a memory architecture problem, not a phrasing problem.

**Core rule: the app owns state; the model narrates.**
- Typed state store: characters (outfit, condition, injuries), inventory (item,
  quantity, holder), location and time, quests/flags, combat state.
- A compact state block is injected every turn (~200–400 tokens). Facts are
  restated, never remembered. This is what stops the clothing drift after 3 messages.
- The model may only **propose** changes, as grammar-constrained JSON deltas
  (`add_item`, `remove_item`, `set_outfit`, `damage`). The app validates and applies
  them. Unparseable or invalid → ignored, state unchanged. Prose never mutates state.
- Dice, HP, initiative, and rules checks run in code. The model rolls nothing and
  calculates nothing; it narrates outcomes the app hands it.
- Rules text is retrieved from a scoped rules lorebook, not stuffed, so a rulebook
  can't eat the context window.

**Two modes:**
- *Model as DM* — AI narrates and adjudicates; the app holds state, rolls, enforces
  rules, feeds results back.
- *User as DM* — you narrate; app gives you the state panel, rolls, and NPC tools;
  AI plays the NPCs.

**The inventory panel is the point.** Visible, editable, always showing what's
carried and by whom. A state panel is what stops "we forgot the rope" drift — not a
better system prompt. This is the part that has to be right.

**Evidence:** 10 exchanges with no outfit drift; inventory changes only through
validated deltas; every roll traceable to an app-side result.

## Architectural law · App owns state, model narrates

The app holds every fact it can enumerate and verify. The model is handed those facts
each turn and returns prose and *proposed* changes. Nothing important is remembered by
the model; nothing important is inferred from vibes.

**Why:** local models are bad at being the source of truth. They hold facts in the same
soft attention that carries tone and personality, so state drifts within a few turns.
Restating a fact every turn costs a few tokens and never drifts.

### Authority ladder (resolve conflicts this way)
1. **What the user explicitly states** — highest. If you say you took off the coat, the
   app records the coat as off. Never argue with the user's own narration.
2. **A validated model delta** — structured, schema-checked, applied.
3. **App defaults** — used when nothing else says otherwise.
Rule: the app wins on facts it owns (item count, HP, seed, model loaded). The user wins
on intent (what they say just happened). The model wins on nothing that is checkable.

### Decision test
| Question | Verdict |
|---|---|
| Can the app enumerate and verify it? | App owns it. Model may propose deltas. |
| Is it prose, tone, description, improvisation? | Model owns it. Leave it alone. |
| Is it a judgement call (is the NPC convinced)? | Model narrates. If rules are on, the app adjudicates the outcome. |
| Is it relevance detection (is this lore relevant)? | App computes it, but **must show and allow override** — heuristics are not truth. |
| Does owning it risk contradicting the user's own words? | Don't own it. Ask or record what they said. |

### Where this already applies in the plan
| Area | Currently | Should be |
|---|---|---|
| Transcript identity / order | App re-renders from partial state (T0.1) | App owns the message list; UI renders it |
| Sampler values | Reset visually on restart (T1.1) | App owns and loads them at startup |
| Loaded model + VRAM | Lost on tab switch (T1.1) | App owns one load state, tabs observe it |
| Prompt format per family | Per-user guessing (T1.4) | App selects the template from the model |
| Think blocks | Risk of feeding reasoning back (T1.3) | App parses, strips, never echoes |
| Which lore fires | Model would decide by vibes | App evaluates triggers deterministically, with a visible override |
| Lore/memory budget | Unbounded | App enforces the ceiling and priority order |
| Character card fields | Model edits prose freely | App owns structure; model fills one field at a time |
| Image prompt/seed/dims/LoRA/board | Scattered | App owns; model may only propose prompt text |
| Outfit, inventory, HP, location | Model's memory (drifts in 3 turns) | Typed store injected each turn (T3.13) |
| Which image is default, cycle timer | Unclear / broken (T2.5) | App owns the default and the timer |
| Speakership / active character | Hardcoded assumption | App owns the speaker set per chat |

### Where NOT to apply it
- Do not let regex triggers silently outrank the model on meaning. Surface the call
  ("this entry fired because it saw 'Moonblade'"), let the user disable it. A hidden
  heuristic that is wrong is worse than a model guess, because nobody can see it.
- Do not freeze anything the user wants to hand-adjust mid-scene. Every app-owned
  field needs a manual edit path.
- Do not make the app author prose. Owning facts is not writing the story.
