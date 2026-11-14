# ENGINEERING — Isolated Native Runtime Enablement

**Date:** 2026-11-10
**Scope:** opt-in isolation of the Wails/WebView2 profile and all application-owned state so the
real desktop application can be tested natively without touching user data
**Not in scope:** F3-C, lorebook ownership, storage-default changes, Wails/Go versions,
dependency changes

## 1. Result

| Area | Status |
|---|---|
| Single opt-in isolation switch, production unchanged when unset | **Verified** |
| WebView2 profile isolated (`options.App.Windows.WebviewUserDataPath`) | **Verified** |
| Application state, logs, chats, caches, scratch files isolated | **Verified** |
| Content-library default destinations isolated | **Verified** |
| Invalid configuration fails loudly, never falls back | **Verified** |
| Unwritable isolation root rejected at runtime (exit 2) | **Verified** |
| Native window starts, frontend↔Go bindings work, Lorebooks page loads | **Verified** |
| Entry created through the real Lorebooks UI, persisted, reconstructed after restart | **Verified** |
| Canonical library file produced by "Import to Library" | **Verified after fixing W9 (§6)** |
| Undo Import rollback, and failure leaving the browser store intact | **Verified** |
| Graceful shutdown and return to ordinary launch | **Verified** |
| Real WebView2 profile and real config untouched | **Verified** (hash and mtime) |
| Full F3-B acceptance (prompt sync, activation/hot-swap, private lore, disconnected roots) | **W27 native gate satisfied (§27)**; activation/hot-swap, disconnected-root display and symlink listing remain open |

## 2. Contract

```
HOWL_WEBVIEW_DATA_DIR   absolute path, optional
```

- **Unset or empty** — ordinary behaviour, byte for byte. WebView2 uses its default profile,
  state resolves to the installation root, `os.TempDir()` is used for scratch.
- **Set** — the value names one isolation root that contains the entire test run:

| Location | Contents |
|---|---|
| `<root>\webview2` | WebView2 user-data folder (`WebviewUserDataPath`), including localStorage and service-worker state |
| `<root>\data` | state root replacing the installation root for everything the app writes |
| `<root>\data\config` | `settings.json`, image/model registries, sampler/model settings, model-intelligence cache, content-library registry and inventory, lorebook store marker/journal/index, error log and startup log, error-reporting preference |
| `<root>\data\Chats`, `Worlds`, `Characters`, `Scenarios`, `Lorebooks` | chat and creative content roots |
| `<root>\data\content\Chats|Characters|Worlds|Scenarios|Lorebooks|Howl` | default content-library destinations (production resolves these under `Documents`/`Pictures`/`Videos`) |
| `<root>\data\images\gallery`, `images\entities` | generated images and entity portraits |
| `<root>\data\workflows` | saved image workflows (production uses `%APPDATA%\HOWL_Chat\workflows`) |
| `<root>\data\persisted-chats` | legacy chat store (production uses `%APPDATA%\Roaming\HOWLChat`) |
| `<root>\data\logs`, `data\cache\tts` | loader diagnostics and TTS cache |
| `<root>\temp` | ASR/image/TTS scratch (`howl-asr`, `howl-image-temp`, `howl-tts`, `howl-kokoro-tts`, `howl-melotts-out`) |

Installation assets are deliberately **not** redirected: models (`models/`), runtime binaries
(`runtime\llama\llama-server.exe`, `sd-server.exe`, `sd-cli.exe`) and the embedded frontend
still resolve from the real checkout, because they are read-only inputs, not state.

### Validation (all fail loudly, exit code 2, no fallback)

The root must be absolute; must not be a drive root; must exist; must be a directory; must sit
outside the project tree (neither containing nor contained by it); must contain
`.howl-isolated-root` whose trimmed content is exactly `howl-isolated-test`; must be writable
(probe file created and removed); and every derived path must be contained by it. Containment is
component-wise via `filepath.Rel`, never a string prefix, so `C:\root-evil` is not inside
`C:\root`.

The marker prevents a mistyped path from adopting a directory that holds real data.
`internal/paths` panics if the variable is set but invalid, so no code path can silently
continue against production; `main.go` checks first and prints the reason before exiting.

## 3. Launch procedure

```bat
:: 1. create a marked isolation root outside the project (example)
set "HOWL_ISOLATION=C:\Users\%USERNAME%\AppData\Local\Temp\howl-isolated-run"
mkdir "%HOWL_ISOLATION%"
echo howl-isolated-test> "%HOWL_ISOLATION%\.howl-isolated-root"

:: 2. build with the project-local CLI and pinned toolchain
build\wails.cmd build

:: 3. launch isolated (same shell, variable scoped to the process)
set "HOWL_WEBVIEW_DATA_DIR=%HOWL_ISOLATION%"
start "" build\bin\howl-chat.exe

:: 4. return to ordinary operation: close the app and drop the variable
set "HOWL_WEBVIEW_DATA_DIR="
:: production launch is unchanged: build\bin\howl-chat.exe
```

`verify-isolation.sh`, `run-isolated.sh`, `run-isolated-lorebook.sh`, `seed-registry.ps1`,
`acl-test.ps1` under `build/tmp/isolation-probe/` are the exact scratch scripts used below; they
are not part of the product.

## 4. Pre-launch verification evidence

The Go probe (`build/tmp/isolation-probe`) asserts the resolved layout instead of trusting it.

```
=== A. isolation ENABLED ===
PASS isolated layout: 15 storage paths inside C:\Users\...\Temp\tmp.XJLNXBaHtN\run1
RESULT PASS

=== B. isolation DISABLED (production unchanged) ===
PASS production layout: unchanged, temp=C:\Users\STONEW~1\AppData\Local\Temp
RESULT PASS
```

Invalid configurations, each exit 2:

| Case | Message |
|---|---|
| `relative/root` | `must be an absolute path` |
| missing directory | `is not usable ... create the directory and add a .howl-isolated-root marker` |
| existing dir, no marker | `has no .howl-isolated-root marker ... refusing to isolate into a directory that may hold real data` |
| the real project root | `must be outside the project tree` |
| a directory inside the project | `must be outside the project tree` |
| unwritable root (deny ACE, temp-only) | `ISOLATION-ERROR "<root>" is not writable: open <root>\.howl-isolation-probe-...: Access is denied.` |

The unwritable case was reproduced end to end by denying `WD,AD` to the current account on a
throwaway temp directory that contained a valid marker; the probe refused it (exit 2) and the
ACE plus directory were removed afterwards. No real user folder's permissions were touched.

Unit tests in `internal/paths/isolation_test.go` cover the same rules plus sibling-prefix
containment, wrong marker content, a file as the root, and the panic-on-invalid contract.

## 5. Native smoke test

Run against `C:\Users\...\Temp\tmp.XJLNXBaHtN\run1` (marker present), binary
`build\bin\howl-chat.exe`, process started with `HOWL_WEBVIEW_DATA_DIR` set.

| Check | Result |
|---|---|
| Process and window | `howl-chat.exe`, window title `HOWL Chat` |
| Isolation banner | `[isolation] root=...\run1 webview=...\run1\webview2 state=...\run1\data` |
| Backend reachable from the webview | `OnStartup: webview reached backend` at +386.5 ms |
| WebView2 profile created in the test root | `webview2\EBWebView` incl. `Local State`, `Default`, `lockfile` — 246–279 files |
| Application state created in the test root | `data\config\logs\{startup.log,howl-errors-20261110.jsonl}`, `data\config\model-scan-dirs.json`, `data\persisted-chats\config.json` |
| Frontend rendering with Go-backed state | nav, Characters, Chat Context, `NO MODEL LOADED`, `No lorebooks selected` |
| Lorebooks page | title `Lorebooks - HOWL Chat`, store bar, `Import to Library`, `Undo Import`, `No lorebook entries yet` — zero bridge errors |
| Graceful shutdown | `CloseMainWindow` → `wails.Run returned`; zero remaining processes |
| Shutdown log | one `SYS-SHUTDOWN ... final save: no active chat to save` record per session — pre-existing and expected with no chat open |

### 5.1 UI write path and restart

1. `data\config\content-libraries.json` was seeded inside the run (test fixture) with one
   lorebook save destination pointing at `<root>\data\content\Lorebooks`, which the app loaded.
2. **Create Entry** was opened in the real Lorebooks page, `Title` and `Content` filled through
   the WebView2 accessibility tree, and saved. The page then listed the card
   `Isolated Smoke Entry` with its content and the status `Stored in browser storage.`
3. The app was closed and relaunched against the same root. The Lorebooks page reconstructed
   `Isolated Smoke Entry` and `Created inside the isolated native run; safe to delete.` from the
   isolated WebView2 profile, with the flag values (`World`, `Priority 3`, `Summary`, `Medium`)
   intact. The entry text is present in the isolated profile's
   `webview2\EBWebView\Default\Local Storage\leveldb`.
4. **`Import to Library`** initially failed with
   `Import failed. The lorebook library is not available in this window.` That was defect W9;
   after the fix in §6 the same click produced
   `Legacy lorebooks are now stored as canonical files. The browser copy was left untouched.`

Production integrity, measured before the first launch and again after the last shutdown:

```
config/settings.json: OK
config/model-registry.json: OK
config/image-settings.json: OK
profile_files=886  (before: 886)
profile_mtime=2026-09-01 23:46:00.957 (before: identical)
```

Real content roots were inspected for post-launch modification times: `Chats` (2026-06-01),
`Worlds`, `Characters`, `Scenarios` (2026-05-11), `images\gallery` (2026-11-03) and `config`
(2026-11-09 22:45:58) all predate the 22:49:41 launch. `Documents\Howl\Lorebooks` still holds no
entries.

## 6. W9 — bridge defect found by the isolated run, and its fix

**Cause.** Wails' AssetServer injects its runtime only into what it considers the main page:
`isRuntimeInjectionMatch` accepts paths ending in `/` or `/index.html`. Every other page that is
navigated to directly is served as a plain file, so `lorebooks.html` ran with no
`/wails/runtime.js` and no `/wails/ipc.js` — no `window.go`. `lorebook-repository.js` then fell
back to `localStorage` for reads and threw
`The lorebook library is not available in this window.` for `importLegacy`/`rollback`. The page
advertised an ownership transition it could never perform, in production as well as in the test
run.

**Fix — one shared seam, not a second bridge.** `webassets.go` adds an AssetServer middleware that
applies exactly the tags Wails itself injects (`/wails/runtime.js`, `/wails/ipc.js`) to every
successful `text/html` response that does not already carry them, wired once in `main.go` via
`assetserver.Options.Middleware`. Non-HTML responses, error statuses and non-GET requests pass
through untouched, and the main page is never double-injected.

**Audit of the other pages.** `webassets_test.go` walks every embedded `.html` file and asserts
exactly one runtime inject per page, so a newly added page cannot silently run without the
bridge. `TestNonHTMLAndErrorsPassThrough`, `TestMainPageIsNotInjectedTwice` and
`TestStandalonePageGetsWailsRuntime` cover the boundaries.

**Harness parity.** `harness_main.go` already proxies `window.go.main.App.*` by reflection and
serves `frontend/dist` over HTTP, so both surfaces present the same bindings object to the same
repository code; page behaviour is the same code path on both. `internal/lorebookstore` and
`internal/library` tests (which the harness drives) pass unchanged.

**Evidence of the fix in the running app.** The import error changed from the page-level
`... is not available in this window.` to a Go-side message, and after correcting a BOM in the
test fixture the import completed and produced the canonical file (§5.1).

## 7. Limitations

- **Import failure path** was observed with an unreadable registry: the page reported
  `content library registry is not initialised.`, wrote no canonical file and kept the browser
  entry. A dedicated "no destination registered" run could not be completed because the launch
  for that step did not start (no window appeared), so that exact wording is unverified; the
  registry file was left restored.
- **No model inference** was exercised (no model loaded), so model-path resolution was observed
  only as an install-root input, which is the intended behaviour.
- **F3-B items**: chat prompt synchronization was addressed at the `SetLorebooks` seam (§10) but
  has not been verified natively with a loaded model; activation/hot-swapping, private-lore
  injection, disconnected-root display and the Windows symlink listing test remain open as of this
  section. Native prompt assembly and private-lore routing were later exercised by W27 (§27);
  model-backed generation, activation/hot-swapping, disconnected-root display and the Windows
  symlink listing test remain open.

## 8. Files changed

| File | Change |
|---|---|
| `internal/paths/isolation.go` | New: contract, validation, containment, `StateRoot`/`StatePath`/`TempDir`/`WebViewDataDir` |
| `internal/paths/isolation_test.go` | New: validation, containment, enabled/disabled and panic tests |
| `internal/paths/content.go` | `InternalConfigDir` follows `StateRoot`; `DefaultContentDir` returns `<root>\data\content\...` when isolated |
| `app-isolation.go` | New: main-package helpers returning `""` when isolation is off |
| `main.go` | Validate isolation, fail with exit 2, pass `WebviewUserDataPath` |
| `app-model.go`, `app-image.go`, `app-image-generate.go`, `app-image-workflows.go`, `app-characters.go`, `app-audio.go`, `app.go`, `startup_log.go` | State, gallery, entity image, workflows, character/world/scenario roots, scratch dirs, RAG root and startup log follow the state root |
| `internal/backend/chat/service.go`, `internal/backend/chat/persistence.go`, `internal/backend/memory/service.go` | Chat content, legacy chat store and memory follow the state root |
| `internal/audio/pipeline.go`, `internal/audio/asr/asr.go`, `internal/audio/tts/{kokoro,melotts,piper}` | TTS cache and scratch output follow the state/temp roots |
| `internal/image/engine/sdcpp/engine.go` | Image temp dir follows `paths.TempDir()` |
| `webassets.go` | New: shared AssetServer middleware injecting Wails' runtime into every HTML page |
| `webassets_test.go` | New: injection, head position, encoding guard, header preservation, every-page audit |
| `app-lorebooks-store.go` | `mergeCanonicalLore`: canonical file text overlays the chat's lorebook selection after cutover |
| `app-prompt.go` | `SetLorebooks` resolves the pool through the canonical overlay |
| `app-lorebooks-store_test.go` | Canonical content wins over a chat snapshot; pre-cutover behaviour unchanged |

`go.mod` and `go.sum` are unchanged; no global environment, registry, installation path or
permission on a real user folder was modified; F3-C was not started.

## 9. Middleware targeted review

Wails' asset server sets `Content-Type`, `Content-Length` and `Cache-Control` on a file response
and performs no compression; it does not emit `ETag` or `Last-Modified`, and `bodyRecorder`
implements only `Write`/`WriteHeader`. The review found three fidelity gaps and one risk, all
closed in `webassets.go`:

- **Injection position.** Wails inserts its runtime as the *first child of `<head>`*. The first
  version inserted before `</head>`, which let a page's own head script run before the bridge
  existed. Injection now lands immediately after the opening `<head ...>` tag
  (`TestInjectionLandsAtStartOfHead`).
- **Compressed bodies.** A body carrying a non-identity `Content-Encoding` is no longer rewritten;
  Wails does not compress today, and this keeps a future encoder from corrupting pages
  (`TestCompressedHTMLIsNotRewritten`).
- **Writer interface.** The wrapper exposes `Unwrap()` for `http.ResponseController` and
  implements `Flush()` (`TestUnwrapExposesUnderlyingWriter`). Writing through the writer returned
  by `Unwrap()` would bypass injection, so that path is documented as read-only for controllers;
  `ResponseController.Flush` prefers the wrapper's `Flush`.
- **Flush semantics.** A flush now commits the buffered HTML through the injection *before*
  flushing the underlying writer (`TestFlushCommitsInjectedBody`). Flushing first would have sent
  the raw prefix and made later injection impossible.
- **Content-Length.** An exact length is set only when the whole document is delivered in one
  finish (`TestContentLengthIsRecalculatedAfterInjection`); the original file length is never
  preserved. If a handler flushes early, the header is removed *before* `WriteHeader` and never
  set again — headers already committed cannot be revised, so the server frames the remainder
  (`TestEarlyFlushRemovesContentLength`). A flush with no body yet does not commit at all, since
  that would also make the later injection impossible (`TestFlushBeforeBodyKeepsInjection`).
- **Preservation.** `Content-Type`, `Cache-Control`, `ETag`, status codes, non-HTML bodies and the
  `Content-Length` recomputation are covered by tests. The every-embedded-page walk remains in
  place, so a new page cannot ship without the bridge.

`gofmt`, `go vet .`, the bridge tests and `go build -tags "desktop,production" ./...` all pass.
No further middleware change is planned without a concrete defect.

## 10. F3-B item 1 — chat prompt synchronization

**Path.** After cutover the page hydrates canonical entries through `LorebookRepo.hydrate()`, the
repository mirrors them into the browser cache, and the chat calls `SetLorebooks` with that list.
Prompt assembly (`buildPrompt`/`buildPromptMessages` → `resolveLorePromptBlock`) resolves against
`a.lorebooks`.

**Defect found.** A chat persists its own lorebook snapshot and `data-sync.js` re-sends it when the
chat is loaded. That snapshot is a copy: a lorebook edited through the Lorebooks page after the
chat was saved never reached the prompt, because the stale copy replaced the canonical entry in
`a.lorebooks`.

**Fix.** `SetLorebooks` — the one backend entry point for the prompt pool — now passes the parsed
entries through `mergeCanonicalLore`. When the store reports cutover, entries whose ID exists in
the canonical roots take their title, content, triggers, scope and priority from the canonical
file, while the chat keeps its own `Enabled` flag (the chat owns the selection, the files own the
text). Before cutover, or with no store, the pool is returned untouched, so the browser store
stays the source. No second backend path was added.

**Tests.** `TestSetLorebooksResolvesCanonicalContent` imports an entry, edits it canonically, then
syncs a stale snapshot: the disabled case injects nothing, the enabled case injects the canonical
text and never the snapshot text. `TestSetLorebooksKeepsSnapshotBeforeCutover` pins the
pre-cutover behaviour.

### Ownership boundaries

| Scenario | Behaviour | Test |
|---|---|---|
| Canonical entry edited after chat save | Canonical text reaches the prompt | `TestSetLorebooksResolvesCanonicalContent` |
| Chat disables an entry | Nothing injected | same test (disabled case) |
| Chat enables an entry | Canonical text injected | same test (enabled case) |
| Canonical entry created after chat save | `UpdateLorebookContext` attaches the canonical entry and enables it | `TestEnablingLorebookCreatedAfterChatSave` |
| Canonical entry deleted after chat save | Dropped from a library-owned pool when the canonical view is complete | `TestDeletedCanonicalEntryDoesNotResurrect` |
| Root disconnected / unreadable | Kept as last-known state, marked stale and excluded from prompt assembly | `TestDisconnectedRootKeepsLastKnownEntry` |
| Conflicting IDs across roots | `mergeCanonicalLore` consumes `Store.List()` items, so the store's conflict rules remain the only resolver | `internal/lorebookstore` tests |
| Legacy world/scenario/character lore | Pool is not library-owned (`libraryOwned=false`), so the overlay never drops it; it keeps its own lifecycle until migration | `syncAppStateFromChatService` |

**Disconnected roots** are deliberately conservative: keeping the last-known entry avoids
treating an unplugged folder as a mass deletion, but the prompt does not yet signal that the
content is stale. That signalling and the last-known-state presentation are the remaining
F3-B work for disconnected roots; this section does not claim them.

## 11. F3-B item 3 — private character lore boundary

`Visibility` was stored and edited by the Lorebooks page but never read by prompt assembly.
Routing for character-scoped entries also treated an empty `OwnerID` as "deliver to every
character", and the Lorebooks form does not write `ownerId` at all, so character lore created
there reached any character's prompt.

- `canConsider` rejects `hidden` entries for every scope: hidden means never injected.
- `RouteResolvedEntries` still delivers an owned character entry only to its owner. An unowned
  entry keeps working only when it is public or carries no visibility label (legacy entries), so
  private unowned character lore can no longer leak.
- World and scenario routing is unchanged; those scopes have no owner to compare, so this
  section makes no claim about their visibility.
- `visibleToUser`, `visibleToOthers` and `revealCondition` remain unenforced on purpose: the UI
  writes `visibility`, and reveal-condition semantics are not defined anywhere in the
  repository, so enforcing them would be guesswork.

Tests: `TestResolveSkipsHiddenVisibility`, `TestRouteCharacterLoreVisibilityBoundary`.

### Assembled prompt and character switching

Routing tests alone would not catch a leak in prompt assembly, so the assembler is tested
directly. `TestCharacterLorePrivacyInAssembledPrompt` builds the character-specific block for
Aria and then Bram, plus the general assembled system prompt: Aria's private entry appears only
for Aria, Bram's private entry never appears for Aria, the unowned private entry appears for
neither, the hidden world entry is absent everywhere, and public and legacy entries keep working.
Rebuilding for Bram after Aria contains none of Aria's private text, and the stored chat history is
not mutated, so no cached block survives a character switch. `TestStaleLoreIsExcludedFromAssembledPrompt`
pins the same boundary for stale entries.

## 12. F3-B item 5 — disconnected-root prompt policy

Deletion and disconnection are now distinct states:

| State | Meaning | Prompt policy |
|---|---|---|
| Current | Canonical root readable and entry present | Eligible under normal rules |
| Unavailable | Root unreadable; entry retained as last-known | Excluded from prompts; marked `Stale` |
| Deleted | Complete canonical view confirms absence | Removed from a library-owned pool |

`Store.List()` returns cached last-known items for an unavailable root as well as live ones, so
`mergeCanonicalLore` splits them: only `Available` items form the current canonical map, and the
rest are kept aside for preservation. A retained entry is marked `Stale` and `canConsider`
refuses to resolve it. The entry stays in `a.lorebooks` and in `LoadChatData`
(`stale: true`) so it can still be shown and its data preserved; reconnecting the root clears the
flag on the next sync. The marker is derived state only (`json:"-"`), never written to chat data
or canonical files. Exclusion is the default: stale content is not injected and there is no
override yet.

Tests: `TestDisconnectedRootKeepsLastKnownEntry`, `TestReconnectedRootRestoresCanonicalContent`,
`TestStaleLoreIsExcludedFromAssembledPrompt`, `TestInferenceBoundaryExcludesStaleLore`.

## 13. Inference-boundary capture (W13–W17)

`App.httpClient` was a concrete `*llama.HTTPClient`, so there was no substitute seam for a test.
A narrow `inferenceClient` interface now names the exact surface the app uses (chat generation and
streaming, context size, loaded/health state); `*llama.HTTPClient` satisfies it unchanged, so
production behaviour and the request path are untouched.

`app_prompt_boundary_test.go` installs a recording client at that call point. It records the exact
messages in memory (no prompt file, no network) and returns a fixed reply so the real
`SendMessageStreamWithImage` path completes through history append and events.

`TestInferenceBoundaryCarriesOnlyCurrentCanonicalLore` proves at the boundary: a canonical edit
beats the chat snapshot, an entry created after the chat save is activated into the request, and
deleted, hidden, stale and private character entries are absent. It also asserts the request shape
(system first, user turn last) and that the reply reaches history.
`TestInferenceBoundaryExcludesStaleLore` proves an unavailable root contributes nothing even though
its last-known entry remains in the pool.

**Finding — character prompts are frontend-composed.** `GenerateCharacterResponse` sends exactly
one user message supplied by the frontend (`TestCharacterResponseBoundarySendsFrontendPromptOnly`),
and `resolveCharacterLorePromptBlock` has no production call site. The Go boundary therefore adds
no character lore at all. The W15 privacy rules protect the shared resolver and the general prompt
path, but character-turn privacy depends on the frontend composer. That is the remaining private
lore gap and must be fixed where that prompt is built, not by adding another backend route.

**Status.** This is prompt-boundary evidence in the real App path, not native GUI evidence and not
model inference. A native Wails capture with this recorder is still required before acceptance.

## 14. W18 — one lore policy for character turns

The character composer in `frontend/dist/js/chat-stream.js` rebuilt lore from `localStorage`
with `enabled` as its only filter: no ownership, no visibility, no stale state, no canonical
availability. The Go resolver was not on that path at all, so W15's rules did not protect
character turns.

The correction does not duplicate the rules in JavaScript. `BuildCharacterLoreBlock(characterName,
message)` is a bound method over the same `resolveLoreBlock` the general path uses: canonical
overlay happens at `SetLorebooks`, then visibility/ownership routing, stale exclusion, trigger
matching and budget apply. The composer calls it with the turn it is about to send and inserts the
returned block verbatim; the localStorage lore section is gone. Both prompt paths now consume one
policy.

Removing the page-side filter also removes a duplication risk: the generic path prepends character
context, and its lore comes from the same resolver rather than a browser copy.

Tests: `TestBuildCharacterLoreBlockSharesPolicy` covers Aria/Bram ownership, unowned private and
hidden exclusion, public/legacy retention and trigger requirement;
`TestCharacterResponseBoundaryUsesSharedLorePolicy` proves the composed prompt reaches the
inference boundary unchanged and contains only the character's own lore. `node --check` passes on
the modified page script.

**Limit.** The binding is unit-tested; the native character-turn UI path has not been exercised.
That, plus the native recorder run, remains the acceptance gap.

## 15. Test-only inference recorder (W20)

`inference_capture.go` (`//go:build promptcapture`) wraps the inference client at the same
boundary the app already uses. `inference_capture_off.go` (`//go:build !promptcapture`) is the
production no-op. Messages are held in memory and returned by `GetPromptCapture`, a binding that
exists only in capture builds. `BuildCharacterLoreBlock` is regenerated in
`frontend/wailsjs/go/main/App.d.ts` and `App.js`, so the composer calls a real binding.

Verified:

- Both tag sets compile. The production binary contains zero occurrences of `prompt-capture` and
  `GetPromptCapture`; the capture binary contains both, so the recorder cannot install into a
  shipped build.
- `build/wails.cmd build -tags promptcapture -o howl-promptcapture.exe` produced a separate
  binary; the existing production `howl-chat.exe` was not replaced.
- Isolated native launch of the capture build: window `HOWL Chat`, log
  `[prompt-capture] recorder armed at inference boundary`, `OnStartup: webview reached backend`,
  graceful exit, 375 files under the isolation root, and the production profile (886 files, mtime
  `2026-09-01 23:46:00.957839400`) plus config hashes unchanged.

**Not verified.** No chat or character turn has been sent through the native GUI with the
recorder armed, so character switching, captured messages, response, history and autosave remain
unexercised at the native boundary. W13–W19 therefore stay code-level or in-process evidence, not
native GUI acceptance.

## 16. First native GUI capture (ordinary chat)

The native window renders the full chat UI, but WebView2 exposes no form controls to UI
Automation: the tree contains only panes. The drive therefore uses window-relative coordinates
from a captured screenshot of the 1400x900 window, types into the message field and presses
Enter. Automation is tied to that layout.

With the capture build and the isolation root:

- The recorder log shows `[prompt-capture] recorder armed require=1 forbid=1` and, after the
  send, `[prompt-capture] call=1 messages=2 require_ok=true forbid_clear=true`.
- Two messages crossed the boundary: the assembled system prompt and the user turn. The probe
  marker was present and the forbidden marker absent. Prompt text itself never left memory.
- The turn persisted under `data/persisted-chats/chats/2026-11-10T06-18-31Z/`: `history.json`
  contains the user message and the recorder reply in order, and the chat directory contains
  meta, characters, lorebooks, memories and the NDJSON history.
- The production profile mtime and config hashes were unchanged again, and no app process was
  left running.

**Scope.** This is native capture of the generic chat path with an empty character list. It does
not yet cover character responses, character switching, canonical lorebook creation, editing an
entry after the chat was saved, enabling a canonical entry after the save, or hidden, deleted and
stale entries. The reply was confirmed in persisted history; the rendered conversation was not
re-inspected visually after the turn.

## 17. Blocker found while staging the two-character scenario

The staged scenario requires character-owned canonical lore. The current UI cannot author it.

- `frontend/dist/js/lorebooks.js` builds the save payload from the editor fields. It includes
  scope, visibility, triggers, priority and `linkedCharacters`, but never `ownerId`, and no field
  in `lorebooks.html` sets an owner.
- `internal/lorebookstore/store.go` `Save` encodes the incoming map directly. A replace does not
  merge the existing document, so a field absent from the payload is removed from the file.
- `internal/backend/lorebook/lorebook.go` routes `ScopeCharacter` entries only when
  `OwnerID == characterName`; an unowned entry is admitted only when public or legacy-unlabelled.

Consequences:

1. A character-scope entry created in the Lorebooks page has an empty owner. If it is private (the
   acceptance case), no character ever receives it: it is dead lore in the canonical library.
2. An owned entry that arrives by import is stripped of `ownerId` the first time it is edited and
   saved in that same page, because the form payload has no owner field.

This blocks the Aria/Bram native scenario at setup, not at the prompt boundary. Hand-writing
canonical JSON files would create the data but bypass the authoring workflow the scenario is meant
to verify, so that was not done. No test data was fabricated, and no product behavior was changed.

Recommended minimal correction: expose an owner for character-scope entries (bound to the linked
character or chosen explicitly), include it in the save payload, preserve an existing `ownerId` on
replace when the payload omits it, and add a save/edit round-trip test for ownership and
visibility. Until then, W13–W19 native character acceptance remains blocked on this gap.

## 18. Ownership authoring correction (W23)

The canonical editor now distinguishes ownership from association:

- The Visibility section of `lorebooks.html` has an Owner Character control, shown only for
  character scope. Its options come from the character library through the existing `GetCharacters`
  binding, with an explicit `Unowned` choice. `linkedCharacters` is untouched and never grants
  ownership.
- `buildLorebookFromForm` always supplies `ownerId`: the selected owner for character scope, an
  empty string otherwise. A private character entry saved with no owner asks for confirmation,
  because no character will receive it; the unowned-private case remains deliberately testable.
- `populateForm` loads the stored owner, retaining a value whose character no longer exists so an
  edit cannot silently drop it.

The save contract in `internal/lorebookstore/store.go` is now explicit:

- field absent from the payload → preserve the owner already on disk (legacy payload);
- field present with a value → set or change the owner;
- field present and empty → clear the owner, producing the intentional unowned state.

Tests: `TestSaveReplacePreservesOmittedOwner` imports an Aria-owned private entry and replaces it
with a payload that has no `ownerId`, requiring `Aria` to survive;
`TestSaveReplaceHonorsExplicitOwnerChangeAndClear` changes the owner to Bram and then clears it
explicitly. Both pass, as do the existing replace, identity, cutover and import tests, the app and
package suites, both build tag sets, and `node --check` on the page script.

**Not yet verified.** The editor round trip through the actual UI has not run, and neither has the
Aria/Bram native scenario. This section records an implemented authoring contract, not acceptance
of the workflow.

## 19. Ownership identity is a display name, and its source now matches chat

`RouteResolvedEntries` compares `OwnerID` with the character name, so ownership identity is the
display name, not the character record's stable `id`.

The first version of the owner selector read `GetCharacters()`, which scans `Characters/` on disk.
Chat adds characters from `CharacterStorage` (localStorage, `howl-chat-characters`), and chat
persistence stores those records. Those are different libraries: characters created in the normal
character page live in localStorage, while the isolated profile's disk `Characters/` directory can
be empty. The selector could therefore offer names the chat would never send to the resolver, or
omit names it would.

Correction: `lorebooks.html` now loads `character-storage.js`, and the selector populates from
`CharacterStorage.loadAll()` first, falling back to `GetCharacters()` only if that storage module is
unavailable. The value written is still the display name, exactly what the resolver compares.

Remaining identity limits, not changed here:

- Two characters with the same display name are indistinguishable to ownership. Private lore
  assigned to one is routed to either, because the comparison is by name.
- Renaming a character does not migrate existing canonical `ownerId` values, so previously owned
  lore stops matching the renamed character. The stored value is preserved on edit, but it is now
  stale.

Resolving either case requires ownership by stable character identity and a defined migration for
existing entries; that is a policy change beyond this correction. It should be decided before
claiming the ownership model is complete for renames or duplicate names.

## 20. Native editor round-trip attempt (incomplete)

With the rebuilt capture binary and the isolation root, the native character flow was driven by
window-relative coordinates:

- The Add Character selector opened with no characters, then created `Aria`. A second creation
  retained the previous field text and produced `AriaBram` instead of `Bram`; a rename attempt
  through the edit button did not apply.
- Reopening the selector showed both records, and they came from `CharacterStorage`, confirming in
  the native window that the corrected owner source reads the same library chat uses. This is the
  one positive result of the attempt.
- Attempts to reach the Lorebooks page through the left rail did not navigate; the window stayed
  on the chat page. WebView2 exposes no form controls to UI Automation, and the book rail icon is
  not a reliable coordinate target.

No canonical entry was created, so no `ownerId` value was written, read back, changed, cleared or
reconstructed after restart. The staged editor checks in W23/W24 remain code-level only. The
mis-created `AriaBram` record exists solely inside the isolation root; the production profile was
not touched.

Automating this editor through raw coordinates is not a sound basis for acceptance. The next
attempt should use a drive path that can address elements directly (the existing UI harness with
its RPC bindings), and reserve the native window for confirming the same flow once the editor
state transitions are known to work.

## 21. Character identity flow and proposed stable-ownership contract

Read-only investigation. No code was changed in this phase.

### 21.1 Actual identity at each boundary

| Boundary | Value used | Source |
| --- | --- | --- |
| Character authoring | display `name` + stable `id` | `CharacterStorage` (`howl-chat-characters`) |
| Lorebook owner selector | display name (current implementation) | `CharacterStorage.loadAll()` |
| Add/remove character in chat | stable `id` | `chatBridge`, `AddCharacterToChat` |
| `GenerateCharacterResponse` | stable `id` | `character.id` |
| Character prompt composition | display `name` | `chat-stream.js` |
| `BuildCharacterLoreBlock` | display `name` | `chat-stream.js` → Go binding |
| Scope routing | `OwnerID == characterName` | `lorebook.RouteResolvedEntries` |
| Canonical entry | free-form `ownerId` string | `lorebookstore` document |
| Chat persistence `characters.json` | `id` + `name` | `CharacterContext` |
| Chat meta `activeCharacters` | IDs; legacy chats use names as pseudo-IDs | `migrateLegacyMeta` |
| Disk character library | directory name; `character.json` may carry `id` | `GetCharacters`, `chat.Service.LoadChatData` |
| Memory continuity | character display name | `memory/continuity.go` |

### 21.2 Lifecycle findings

- **Stable ID exists.** `collectCharacterData` generates `char_<timestamp>_<random>` and
  `CharacterStorage.save` upserts by `id`.
- **Edit preserves it.** `saveCharacter` restores `existing.id` before saving.
- **Rename preserves it.** The name is only a field; the ID and storage key do not change.
- **Delete removes it.** `CharacterStorage.remove(id)` drops metadata, images and chat references.
  Recreating a character generates a new ID.
- **Duplicates are allowed.** No uniqueness check exists for names or IDs beyond the generated
  random suffix.
- **No duplicate/clone action and no Howl import/export path were found.** External SillyTavern
  cards are scanned read-only from content libraries; `CreateCharacter` (the disk writer) has no
  frontend callsite.
- **Legacy migration exists for localStorage** (`migrateCharactersWithIds`) and for chat meta
  (`migrateLegacyMeta`), but the latter falls back to storing display names as pseudo-IDs.

### 21.3 Why display-name ownership is unsafe

1. Duplicate display names collapse to one owner string; private lore for one reaches both.
2. Renaming changes the comparison value; existing `ownerId` values become stale and stop matching.
3. A deleted character's name can be reused by a new record, inheriting access that was not granted
   to it.
4. Legacy chat meta can contain names where IDs are expected, so the two identity spaces mix.
5. Authoring characters live in localStorage while the disk `Characters/` library is a separate set;
   a name alone does not identify which record was meant.
6. The memory subsystem also keys character scope by name, so the same ambiguity exists outside
   lorebooks.

### 21.4 Proposed contract

Reuse the existing `character.id`; do not introduce a new identifier.

- **Ownership is the stable ID.** Canonical `ownerId` stores `character.id`. Character-scope
  routing compares `OwnerID` to the character ID supplied by the prompt path.
- **The editor displays names but persists IDs.** The owner select's option value is the ID; its
  label is the current name. Duplicate names are shown with a short ID suffix so they are
  distinguishable. A missing record is shown as `Deleted character (<id>)`, never silently
  converted to unowned.
- **One identity per character turn.** `BuildCharacterLoreBlock` takes the character ID (and the
  message); the frontend passes `character.id`, while the display name remains only prompt text.
  `GenerateCharacterResponse` already uses the ID and stays unchanged.
- **General chat path** supplies the active character's ID when it has one. Without an ID, private
  character lore is not delivered (fail closed), and owned entries do not match by name.
- **Rename** changes only the label; ownership follows the unchanged ID.
- **Delete/recreate** produces a new ID, so old private lore cannot be inherited.
- **Legacy migration is versioned and explicit.**
  - The frontend supplies a `name -> [ids]` map built from the same `CharacterStorage` the chat
    uses, plus IDs found in persisted chat character records.
  - A new backend binding applies the map only where a name resolves to exactly one ID. It rewrites
    `ownerId`, records `ownerMigratedFrom` and a migration version, and journals the previous value
    for rollback. The operation is idempotent and re-runnable.
  - Ambiguous names (duplicates) and unknown names are **not** guessed. They are reported with their
    entry IDs and left unchanged. Because routing now compares IDs, private entries with unresolved
    name values never match and remain excluded until resolved manually.
  - Legacy name strings are never reinterpreted as IDs merely because they look similar. Migration
    matches only an exact, unique known name that is not also an existing character ID.
  - A migration-version marker and journal live beside the lorebook cutover data; a rollback
    binding restores prior `ownerId` values from the journal.
- **Import/export** keeps `ownerId` as an opaque ID. Imported legacy name values enter the same
  unresolved path; they are not auto-assigned.

### 21.5 Product decisions required before implementation

1. **Migration authority.** Character authoring data is frontend localStorage; the Go store cannot
   resolve names to IDs by itself. Approve the frontend-supplied, backend-journaled map, or require
   a Go-side character registry as the authority first.
2. **Legacy chats without character IDs.** Confirm that private owned lore for such chats should be
   withheld until the user re-adds the character, rather than matched by name.
3. **Snapshot label.** Decide whether canonical entries should also store `ownerName` for display
   and audit, or always resolve the label live from the character library.
4. **Memory subsystem.** Its name-based character scope has the same defect; confirm it is out of
   scope for this correction or schedule it with the same identity contract.

### 21.6 Required tests (after approval)

Stable ID across save/reload/restart/rename; identical display names with distinct ownership; 
private lore never crossing between duplicates; rename preserving ownership; delete/recreate not 
inheriting; unambiguous legacy migration; ambiguous legacy exclusion and reporting; import/edit/
save/restart; explicit owner change and clear; both prompt paths using the same identity; and the 
existing W13–W19 privacy and freshness suites.

**Status.** Phase 1 and Phase 2 were approved. Phase 3 implements the stable-identity
contract; native Wails lifecycle acceptance is recorded in §27.

### 21.7 Phase 3 implementation

- **One identity.** `character.id` (`char_...`) is the only ownership identity. A display name is
  a label and never authorizes delivery.
- **Resolver boundary.** `RouteResolvedEntries` now takes a character ID. `BuildCharacterLoreBlock`
  takes the ID (the frontend passes `character.id`), and the legacy chat service derives the active
  character ID from its character record. If no ID is available, private character lore is withheld.
- **Editor.** Owner options use the stable ID as their value and the current name as their label.
  Duplicate names carry a short ID suffix. A legacy, deleted, or otherwise unknown owner is shown
  as `Unresolved owner: <value>` and is preserved, not silently converted to unowned. The form sends
  `ownerId` and an optional `ownerName` label.
- **Persistence.** `ownerName` is descriptive only. `Save` preserves an omitted `ownerName`, applies
  an explicit label, and clears the label when the owner is explicitly cleared. The existing
  absent/preserve/set/clear rule for `ownerId` is unchanged.
- **Migration.** `MigrateLorebookOwners(identities, version)` accepts the frontend character library,
  rejects malformed or conflicting maps, and rewrites only character-scope entries whose owner is an
  exact unique name. Known stable IDs are never overwritten. Ambiguous, unknown, read-only, and
  duplicate-library entries are reported and left untouched. Writes record `ownerMigratedFrom` and
  the migration version, are idempotent, and are journaled. `RollbackLorebookOwnerMigration` restores
  the prior owner only while the entry still has the migrated value, preserving unrelated fields.
- **Tests.** Store coverage: unique migration, idempotency, ambiguous/unknown exclusion, stable-ID
  non-overwrite, malformed maps, cutover requirement, rollback, post-migration edit safety, and
  `ownerName` preservation/clear. Resolver coverage: unowned visibility boundary, duplicate names,
  rename-label independence, and `ownerName` never authorizing. App coverage: both the assembled
  prompt and the character-response binding route by ID. Frontend coverage: ID survival through
  edit/rename/restart, delete/recreate receiving a new ID, duplicate-name labels, unresolved-owner
  visibility, and migration reselection.
- **Evidence.** `go test ./internal/lorebookstore/ ./internal/backend/lorebook/ ./internal/backend/chat/ .`,
  the focused frontend `node --test` file, `node --check` for the edited scripts, and both tagged
  builds pass. The focused identity harness passes 4/4. During final verification the pre-existing
  `tests/frontend/character-configuration-integrity.test.js` voice-model-clear test failed both alone
  and in the serial suite; it exercises untouched `character-form.js`/`character-storage.js` code, so
  it is recorded as an existing failure, not attributed to this change. No native workflow was run and
  no UI screenshot was captured; the editor evidence is the VM DOM harness only, not native acceptance.

## 22. Memory and legacy identity audit

### 22.1 Memory uses an independent name-based path

Automated memory capture is not on the stable character identity yet: `memory.inferOwner` writes a
character display name into `MemoryRecord.OwnerID`; `CharacterContinuity.ExportCharacterMemories`
and `ImportCharacterMemories` compare that value to a character name; and `generateSecondaryTriggers`
derives trigger words from the owner label. The frontend `MemoryStore.getCharacterMemories` already
filters by `character.id`, so it does not match those name-owned records.

Containment applied at the prompt boundary: `chat.loadMemoriesFromService` now receives stable
character IDs and admits only character-scope memories whose `OwnerID` matches one. Name-owned
memory therefore fails closed instead of reaching a prompt. The remaining risk is continuity: a new
character that reuses an old display name could export or import memories owned by the earlier
record. Fixing that requires writing and reading memory ownership by `character.id` (or disabling
continuity for name-owned records). That is a separate memory scope and was not changed under F3-B;
it is reported here as a containment requirement.

### 22.2 Legacy disk character library

`GetCharacters` fabricates `id` from the character directory name when a JSON file has no `id`, and
`loadCharacterFromDiskByID` matches that raw value. Those directory-name IDs are not the `char_...`
identity generated by `CharacterStorage`; they cannot distinguish duplicate names and can change
with a directory rename. The lorebook editor prefers `CharacterStorage` and shows unknown values as
unresolved, so this path is not silently promoted to ownership. Treating disk IDs as an
authoritative identity would require a separate product decision and migration; none was invented
here.

### 22.3 Alternate prompt paths

The general prompt path passes an empty character ID, so only unowned/public character lore and
world/scenario lore can enter it. The legacy `chat.Service` path now resolves the active character
ID and withholds character lore when the ID is missing. `characterNames` and `characterData` remain
in use for prompt text and selection, not for lorebook authorization. A search for remaining
name-based lorebook comparisons found none.

### 22.4 Evidence classes

Go unit/integration tests, the frontend VM DOM harness, compile/build checks, and native UI evidence
are distinct. This phase has the first three. It has no native run, no element-addressable native
harness result, and no screenshot; native acceptance is not claimed.

## 23. Memory identity containment (W28)

### 23.1 Trace and reachability

Capture writes `MemoryRecord.OwnerID` through the frontend memory store and `AddMemory`, and through
`memory.ProcessMessageForEvents` → `BatchClassify`. Persistence is `memories.json` plus
`candidates.json`; retrieval flows through `GetMemories`, the character viewer, and
`chat.loadMemoriesFromService`; prompt insertion uses the resulting lorebook entries and
`memories_snapshot.json`. Generic JSON export is in Story Tools; continuity export/import lives in
`memory.CharacterContinuity`.

`ProcessMessageForEvents` and the continuity methods are not Wails bindings in this build (no
`App` wrapper or generated binding), so the current UI cannot invoke the backend continuity path.
The UI's automated-capture fallback stores chat-owned records, and the death/swap paths store a
character ID or chat session. The API boundary was nevertheless unsafe because it accepted display
names, and legacy records may already contain them.

### 23.2 Exposure finding

Before containment, continuity export/import authorized by display-name equality. A deleted and
recreated character with the same name could therefore receive the earlier character's legacy
name-owned memories if a caller passed that name. Duplicate names collapsed to one owner. Prompt
retrieval compared IDs, which normally excluded names, but a legacy owner string that happened to
equal a character ID would have been accepted without proof that it was an ID.

### 23.3 Containment implemented

- New character memories carry `owner_kind: character_id` and the stable `OwnerID`. Capture resolves
a display name to an ID only when exactly one identity matches; duplicate/ambiguous or missing
identities remain chat-owned and are never guessed.
- Continuity export requires `owner_kind == character_id` and the exact target ID. Import accepts
only records with the same stable ID and kind; legacy name-owned and foreign-owned records are
skipped, retained in storage, and never reassigned.
- Prompt insertion filters memory records before conversion (`IsStableCharacterOwnershipForAny`) and
filters pre-existing `memories_snapshot.json` character entries that lack the marker. `ToLorebookEntry`
carries the marker so automated-lore filtering can require it.
- Frontend character retrieval (`MemoryStore.getCharacterMemories` and the character viewer) requires
the marker as well. The memory list shows a visible unresolved-ownership warning, and
`GetMemoryIdentityDiagnostics` reports total/stable/unresolved counts and unresolved owner labels
without changing data.
- Legacy records are neither deleted nor rewritten. There is no invented ID and no automatic
migration.

### 23.4 Residual risks and decisions

Legacy character memories remain unresolved and are excluded from prompts and transfer until a
future explicit ownership repair. Existing `automated_character_lore.json` entries without an owner
kind are likewise withheld from character prompts (fail closed). Generic JSON export still includes
legacy records but is not an authorization transfer. No new Wails binding currently calls capture
or continuity; any future binding must pass stable identities and target IDs, not names.

### 23.5 Tests and evidence

Memory tests cover stable-ID capture, duplicate-name and missing-identity ambiguity, legacy versus
stable export, same-name recreation, foreign-identity import rejection, unchanged legacy data,
prompt-filter ownership, diagnostics reporting, and converter metadata. A chat test covers legacy
snapshot exclusion from the prompt pool. The frontend memory test covers stable-only retrieval and
explicit marker persistence. Go tests, the focused frontend memory test, both tagged builds, and
JavaScript syntax checks pass. Native memory UI evidence was not captured; this is code-level
containment only.

The W27 native gate for F3-B was satisfied later (§27). F3-C remains blocked. Mobile Remote API
remains deferred.

## 24. W28 follow-up — import classification and prompt-route audit

### 24.1 Export and import classification

There is no generic JSON import UI in this build; the only receiving memory import path is
`CharacterContinuity.ImportCharacterMemories`, which receives structured records rather than parsing
a downloaded file. Generic Story Tools export uses the frontend `MemoryRecord` DTO, which now carries
`OwnerKind`, and the internal `types.MemoryRecord` JSON carries `owner_kind`.

`TestExportRoundTripDoesNotReclassifyLegacyMemory` marshals and unmarshals a legacy record whose
`OwnerID` is exactly the destination `characterID`, alongside a stable-ID record. After the round
trip the legacy record still has an empty kind; export returns only the stable record, and import
accepts only the stable record with its original owner and kind. A legacy ID-shaped owner string is
therefore not enough to authorize transfer. `TestTypesMemoryToRecordPreservesOwnerKind` checks the
frontend export DTO keeps the same distinction.

No generic importer was added. If one is introduced later, it must preserve `owner_kind` and use the
same fail-closed rule; parsing a bare owner string as character-owned is forbidden.

### 24.2 Prompt-insertion route audit

| Route | Ownership filter | Result |
| --- | --- | --- |
| `chat.loadMemoriesFromService` → `memory.ToLorebookEntries` | `IsStableCharacterOwnershipForAny` before conversion | The only converter call site is filtered. |
| `chat.LoadChatData` → `memories_snapshot.json` | `filterStableCharacterMemoryEntries` on load | Both `SendMessage` and `SendMessageStream` consume the filtered `memorySnapshot`. |
| Character composer in `frontend/dist/js/chat-stream.js` | `MemoryStore.getCharacterMemories(character.id)` requires `ownerKind === 'character_id'` | Character prompts receive only stable-ID records. |
| `app-prompt.go` `SetLorebooks` / `a.lorebooks` | No memory service or `MemoryRecord` path | This is the lorebook pool, not character memory. |
| Candidate/fallback capture | Backend capture method is not Wails-bound; fallback writes `scope: chat` | No character-private prompt insertion. |
| Legacy `manualLore` files | Lorebook ownership, not `MemoryRecord` | Not a memory-record bypass. |

`memory/triggers.go` uses `OwnerID` to derive trigger words, not to authorize prompt delivery; legacy
records remain excluded by the prompt filters. No other `ToLorebookEntries` call exists.

### 24.3 Native verification

The first native attempt was blocked at the automation boundary: WebView2 did not expose a
remote-debugging endpoint and the accessibility tree has no form controls. That attempt produced no
memory UI evidence. It was superseded by the test-tagged seam in §25; the production binary was not
altered.

### 24.4 Verification

`go test ./internal/backend/memory/ ./internal/backend/chat/ ./internal/backend/lorebook/ .` passes,
including the JSON import round trip and snapshot filter. The focused frontend memory tests pass
2/2. The native sequence is recorded separately in §25 and is not used to claim a production-binary
UI run.

## 25. W28 native memory verification (test-tagged seam)

### 25.1 Seam contract

`app_native_test_seam.go` is built only with `native_test_seam`; the production path compiles
`app_native_test_seam_off.go`, whose start/stop hooks are no-ops. The seam starts only when
`HOWL_NATIVE_TEST_SEAM=1` and `paths.IsolationActive()` is true. It binds to
`127.0.0.1:0` with `net.Listen`, never a network interface, and logs only the loopback address. It
records no expression or result in the application log. `/__howl_native_test_seam/eval` evaluates
an expression in the real Wails page; `/__howl_native_test_seam/snapshot` loads a chat fixture and
returns only snapshot entry IDs through the test-tagged `chat.Service.NativeTestSnapshotIDs`.

The production binary was checked for `HOWL_NATIVE_TEST_SEAM`, `__howl_native_test_seam`,
`native-test-seam`, and `NativeTestSnapshotIDs`; all were absent. All four were present in
`build/bin/howl-w28-seam.exe`. The tagged seam is therefore separate from the production artifact.

### 25.2 Native evidence

The tagged binary was run against the marked isolation root
`%LOCALAPPDATA%\Temp\howl-w28-native-run1`. The startup log confirmed the isolated state and
WebView2 roots; the app reached the backend at `http://wails.localhost/`. Six memory records were
seeded into the isolated WebView2 profile: two legacy character records (one with a display-name
owner and one whose owner string equals `char_aria`), stable records for `char_aria` and
`char_bram`, a record for the recreated identity `char_aria_new`, and one chat-owned record.

| Check | Native result |
| --- | --- |
| Legacy-ownership warning | Flyout opened; exact text: `2 character memories have unresolved legacy ownership and are excluded from character prompts and transfer.` |
| Stable retrieval | `getCharacterMemories('char_aria')` and the character viewer returned only `mem_aria`. |
| Character switching | `char_aria` returned `mem_aria`; `char_bram` returned `mem_bram`; neither contained the other's record. |
| Same-name recreation | `char_aria_new` returned only `mem_aria_new`; neither Aria legacy record appeared. |
| Snapshot bypass | The fixture snapshot contained `legacy_name`, `legacy_id_shape`, and `stable`; the loaded snapshot returned only `stable`. |
| Restart persistence | After a full app restart, all six records remained; retrieval and the warning were unchanged. |
| Legacy preservation | The persisted list still contained `mem_legacy_name` and `mem_legacy_idshape`; the warning counted both. |

The snapshot result proves the filter at snapshot load. It does not exercise a live model prompt;
the two send paths consume this same filtered `memorySnapshot`, which is covered by the backend
suite. The test app was stopped afterward. The production profile and config were not used or
modified by this run.

### 25.3 Limits

This is native evidence from the explicit test-tagged binary, not a production-binary UI run. The
seam is compile-time absent from production, and the app behavior outside the seam is unchanged. No
model inference or native screenshot was captured; the evidence is the real Wails page DOM and the
isolated test endpoint.

## 26. Native test seam security audit (W28c)

The W28b seam was reviewed before reuse. The first version bound to loopback and was compile-time
absent from production, but it did not authenticate callers or validate browser origin. A hostile
page could therefore attempt a simple cross-origin POST to the evaluator, even though CORS would
hide the response. The following controls now close that gap:

- **Per-run token.** Startup generates 32 random bytes, hex-encodes them, and writes the token to
  `<isolated state>/config/native-test-seam.token`. Every endpoint requires the token in
  `X-Native-Test-Seam-Token`; health included. The token is never logged.
- **Windows ACL.** The token file is written with mode `0600`, then `icacls` removes inherited ACEs
  and grants only the current user. Verified output was only
  `STONEWOLFPC\StoneWolfPC:(F)`; the inherited `CodexSandboxUsers` ACE was gone. POSIX mode display
  is not used as the Windows authorization check.
- **Origin gate.** A request carrying `Origin` must use `http://wails.localhost` or
  `https://wails.localhost`; other origins receive 403 before token processing. A missing Origin is
  accepted only for the local test driver, which still needs the token. CORS is limited to the same
  origins and the token/content-type headers; `OPTIONS` returns 204 only after the origin check.
- **Constrained surface.** Fixed paths; `GET` health and `POST` eval/snapshot/result only; 1 MiB
  limits on eval and snapshot, 4 MiB on result; 5 s header timeout and 10 s eval timeout. Unknown
  paths return 404 and wrong methods return 405.
- **Shutdown and rotation.** Shutdown calls `http.Server.Close`. A graceful window close
  (`CloseMainWindow`/WM_CLOSE) exited the process, and a request to the former port was refused. A
  restart generated a new token and the prior token returned 401.

Observed checks: `health` without token 401, with token 200; evil origin 403; allowed-origin
preflight 204; eval without token 401, with token `{"ok":true,"value":2}`; snapshot with token
returned `{"ids":["stable"],"ok":true}` and without token 401; wrong method 405; unknown path 404.
The production binary contains none of the seam identifiers, including the token header.

### 26.1 Residual trust boundary

This is test infrastructure, not a production security boundary. A local process running as the
same Windows user can read the token file; the ACL prevents other users and the inherited sandbox
group, not same-user code. The evaluator executes in the real Wails page and can invoke every
`window.go` binding the UI can invoke, including path-taking operations. Isolation redirects
app-owned state, but it does not turn those bindings into a filesystem sandbox. The seam must only
be run against a marked isolated root by a trusted test operator; it must never be exposed on a
network interface or used with the production profile. If future work needs a boundary against
same-user local processes, use an OS-authenticated pipe or a narrow allowlisted bridge instead of
arbitrary page evaluation.

## 27. W27 native Lorebook lifecycle acceptance

### 27.1 Build and isolation

The lifecycle was exercised in the real Wails WebView2 page against the marked isolation root
`%LOCALAPPDATA%\Temp\howl-w27-native-CGJlKh`, with the content library seeded at
`data/content/Lorebooks`. The test binary `build/bin/howl-w27-acceptance.exe` was built with
`native_test_seam,promptcapture`; the untagged production binary `build/bin/howl-w27-native.exe`
contains no `native_test_seam`, `prompt-capture`, `GetPromptCapture`, or `HOWL_NATIVE_TEST_SEAM`
symbols. The test run used the real lorebook editor functions, the real repository, the product
Wails bindings, and the existing prompt-capture recorder. It did not use the production profile.

### 27.2 Native evidence

| Check | Native result |
| --- | --- |
| Character fixtures | `CharacterStorage` contained two distinct stable IDs: `char_aria_w27` (Aria) and `char_bram_w27` (Bram). |
| Native creation and ownership | Two private entries were created through the Lorebooks page editor. The owner selector offered the stable IDs; `Aria W27 Private` stored `char_aria_w27` / `Aria`, and `Bram W27 Private` stored `char_bram_w27` / `Bram`. |
| Explicit import | `importLorebooksToLibrary` produced `cutover: true`, two canonical files in `data/content/Lorebooks`, and preserved both `ownerId` values. |
| Stable-ID read-back | Reopening each entry showed the original owner ID, not a display name. |
| Rename | Renaming Aria to `Aria Renamed` kept `ownerId: char_aria_w27`; the editor reselected that ID and updated only the descriptive `ownerName` to `Aria Renamed`. |
| Prompt inclusion and exclusion | `BuildCharacterLoreBlock('char_aria_w27', ...)` contained only `ARIA-W27-PRIVATE`; the Bram call contained only `BRAM-W27-PRIVATE`; rebuilding Aria after Bram still contained only Aria's entry. |
| Inference boundary | The real `chat-stream.js` composer sent one user message through `GenerateCharacterResponse`. The `promptcapture` recorder saw Aria's composed prompt with `ARIA-W27-PRIVATE` and no Bram entry, and Bram's with `BRAM-W27-PRIVATE` and no Aria entry. |
| Deleted owner | After removing Bram from `CharacterStorage`, the editor showed `Unresolved owner: char_bram_w27` and kept that owner ID. Saving did not reassign it; `ownerName` became empty. Aria's block had no Bram entry, and the recreated `char_bram_new` had no Bram entry. |
| Restart | A graceful window close and full restart preserved `cutover: true`, both canonical files, Aria's renamed label, the unresolved Bram owner, and the same Aria/Bram prompt behavior. |
| Production surface | The untagged production binary built successfully and contains no test-seam or prompt-capture symbols. |

### 27.3 Defect found during acceptance

The first native prompt pass exposed a real policy defect. `ResolveConflicts` grouped entries
only by matched trigger phrase, so two character-scoped entries with the same trigger and priority
were collapsed before `RouteResolvedEntries` ran. Bram's entry suppressed Aria's entry even though
the entries belonged to different stable character IDs. Each entry worked alone; together, Aria's
block was empty.

The fix adds `conflictGroupKey`: character-scope entries are grouped with their `OwnerID`, so
different owners do not compete, while same-owner conflicts still resolve by the existing rules.
`FindConflicts` uses the same grouping. Regression coverage is
`TestResolveConflictsKeepsDifferentCharacterOwners`, `TestResolveConflictsStillCollapsesSameCharacterOwner`,
and `TestConflictGroupingKeepsNonCharacterScopesShared`, which pins the boundary: world and scenario
entries with the same trigger stay one conflict domain even when their `OwnerID` fields differ, so
the ownership split applies only to private character entries. Focused lorebook, chat, and app
tests pass, the frontend identity test passes 4/4, `gofmt` is clean, and the W27 production build
succeeds.

### 27.4 Evidence levels and limits

- **Native UI behavior:** verified in the real Wails page against the isolated root; no screenshot
  was captured.
- **Assembled prompt:** verified through the test-only `promptcapture` recorder at the same
  `inferenceClient` boundary the app uses. No network request was made and no prompt was written to
  disk.
- **Successful model inference:** not tested. No model was loaded or invoked, so this run proves
  prompt assembly and delivery to the inference client, not generation quality or a real model
  response.

**Result.** W27 native Lorebook lifecycle acceptance passes at these evidence levels. The W27/F3-B
native gate is satisfied; F3-C remains not started. The separate F3-B limitations recorded elsewhere
(model-backed prompt synchronization, symlink listing, and other attachment forms) are unchanged.

## 28. F3-B closure review (2026-11-10)

The eight checks below were run against the final working tree after W27/W27a. No regression or
contradictory evidence was found; the earlier ACCEPTED WITH DOCUMENTED LIMITATIONS verdict stands.

| Check | Result | Evidence |
| --- | --- | --- |
| Every W27/W28 change is represented in the report and ledger | Pass | Report §§23–28; `PROJECT_ISSUES_AND_EXPANSION.md` rows W27, W27a, W28, W28a–W28c and the F3-B closure section |
| No stale status describes completed work as pending | Pass after correction | §1 and §23.5 corrected here; §7 annotated; the F3-B and Wails-CLI records carry dated addenda for the cleared native blocker |
| No production build contains test-only seams or capture instrumentation | Pass | Rebuilt `howl-w27-native.exe`: 0 occurrences of `native_test_seam`, `HOWL_NATIVE_TEST_SEAM`, `promptcapture`, `prompt-capture`, `GetPromptCapture`, `X-Native-Test-Seam-Token`; `BuildCharacterLoreBlock` still present (5) |
| The canonical Lorebook path is the only authoritative writer | Pass (source) | Only `internal/lorebookstore` writes `howl-lorebook-v1` entry files; the direct cache writes in `lorebook-manager.js` have no callers outside that file |
| Character identity and memory privacy invariants hold across consumers | Pass | `go test ./internal/backend/memory ./internal/backend/chat .` exit 0; frontend identity 4/4 and memory identity 4/4 |
| Conflict grouping does not change world, scenario, or global behavior | Pass | `TestConflictGroupingKeepsNonCharacterScopesShared` over `ScopeWorld` and `ScopeScenario`: same trigger resolves as one domain, keyed by trigger alone |
| Affected tests and the production build pass from the final tree | Pass | `go test ./internal/backend/lorebook ./internal/backend/memory ./internal/backend/chat .` exit 0; `build/wails.cmd build -o howl-w27-native.exe` exit 0 (8.16 s) |
| No test configuration, fixture, process, or isolation setting can affect a normal launch | Pass, one observation | No `.env`, no `HOWL_*` variable in this environment, no isolation marker in the project tree; the isolated run wrote only `<root>\data\config`. See below. |

**Observation.** The real `config/logs/startup.log` contains a non-isolated process start at
08:12:27 that exited after 52 ms (`wails.Run returned`), plus three discovery scans from Go tests
between 08:13 and 08:21. They wrote the real startup log, and `config/model-scan-dirs.json` has
the same 08:12 mtime with a path normalization; its content remains the operator's own scan
directories. The isolated acceptance run wrote its own `data\config\model-scan-dirs.json`, proving
isolation itself did not leak into the project config. Nothing here changes launch behaviour:
the switch is unset, the marker is absent from the tree, and the entry paths are log-only.

**Not verified.** Successful model inference, activation/hot-swapping in a live creative workflow,
disconnected-root display in the native window, and the Windows symlink listing test remain open as
documented; none is an F3-B ownership failure.

**Status: F3-B CLOSED.** F3-C is unblocked and not started.
