# ENGINEERING — Wails CLI Integration and Build Reliability

**Date:** 2026-11-10
**Scope:** project-local Wails CLI invocation, generated TypeScript bindings, native WebView2
test isolation
**Not in scope:** F3-C, lorebook ownership, content persistence, dependency versions

## 1. Result

| Area | Status |
|---|---|
| Project-local Wails CLI wired into repeatable commands | **Verified** |
| Go toolchain pinned per invocation, no global changes | **Verified** |
| Generated bindings regenerated, deterministic, API-complete | **Verified** |
| Native executable produced by the wrapper | **Built but not runtime-tested** |
| Native WebView2 runtime acceptance | **Blocked** (isolation requires a product change) |
| TypeScript typecheck of generated bindings | **Unverified** (no TypeScript toolchain installed) |

## 2. Final invocation

```bat
build\wails.cmd build             :: native build -> build\bin\howl-chat.exe
build\wails.cmd dev               :: live-reload development session
build\wails.cmd doctor            :: environment report
build\wails.cmd generate module   :: regenerate frontend\wailsjs bindings
```

`build\wails.cmd` is the single supported entry point. It is a thin wrapper over the
project-local CLI at `build\tools\wails.exe` (v2.12.0, `build/tools/` is gitignored) and adds
only what makes the invocation reproducible:

- switches to the project root, so it works from any working directory;
- refuses to run without `wails.json`, `go.mod`, or the CLI, with the install command in the
  message;
- reads the required Go version from `go.mod` (the `toolchain` directive, else the `go`
  directive), sets `GOTOOLCHAIN` for this invocation only, and then re-checks `go version`.
  If the resolved version is not the required one it aborts rather than building with the
  machine default;
- forwards all arguments unchanged and exits with the CLI's exit code;
- changes no global PATH, GOBIN, Go environment, registry, or machine-wide installation.

`start-howl-chat.bat` dev mode previously looked for a global `wails.exe` on `GOPATH\bin`.
It now calls `build\wails.cmd dev`, so dev and build resolve the same project-local CLI and
the same toolchain.

## 3. Toolchain evidence

```
default `go`:                 go version go1.27.3 windows/amd64
GOTOOLCHAIN=go1.23.8 go:      go version go1.23.8 windows/amd64
GOVERSION under pin:          go1.23.8
GOROOT under pin:             C:\Users\StoneWolfPC\go\pkg\mod\golang.org\toolchain@v0.0.1-go1.23.8.windows-amd64
```

The module pins `go 1.23.8`; the machine default is 1.27.3, under which `wails build` fails
with `internal error: package "time" without types was imported from ...`. Nothing was
upgraded or downgraded to work around this: the module already requested 1.23.8, and the
wrapper simply makes the invocation honour it. `go.mod` and `go.sum` are unchanged by the
build.

`wails doctor` reports **Go Version go1.27.3**. That row is the Go runtime the CLI binary was
compiled with, not the compiler used for this module; the pin is confirmed by the
`GOVERSION`/`GOROOT` values above and by the successful build.

## 4. Generated bindings — retained

Decision: **retain the regenerated files, unedited.**

The regenerated diff is `App.d.ts` +109, `App.js` +216, `models.ts` +691 (980 insertions,
36 deletions). Evidence for correctness:

| Check | Result |
|---|---|
| Exported `App` methods in the live tree | 176 |
| Methods present in `App.d.ts` / `App.js` | 175 / 175 (identical sets) |
| In the Go API but not in the bindings | `OnShutdown` only |
| In the bindings but not in the Go API | none |
| Imported type namespaces declared in `models.ts` | 10 of 10 (`library`, `main`, `audio`, `lorebookstore`, `metadata`, `types`, `chat`, `settings`, `modelintel`, `rag`) |
| Regeneration is deterministic | Yes — byte-identical after `generate module` and after a full `build` |

`OnShutdown(ctx context.Context)` is the Wails lifecycle hook passed as `options.App.OnShutdown`.
The binding generator excludes methods taking `context.Context`; it is invoked by the runtime,
never from JavaScript. The frontend does not import `frontend/wailsjs` at all (0 references in
`frontend/dist`); it calls the runtime namespace `window.go.main.App`, which Wails populates
from the same `Bind` list. The bindings are therefore documentation and editor tooling, not a
runtime dependency — and the stale versions were simply out of date. No generated file was
hand-edited.

Determinism hashes (SHA-256):

```
04f70b772d8c0f2659917231a8f8811eea8109dd910d33ffec87ccc66816fd4d  App.d.ts
015e9838d1dd35719e32935f3fdf30336c735211211a5b993cca4febf0010957  App.js
a5de3cc6d6f93726b0d2ad2a6377294c39de5c01d04036956ca88d6a75696953  models.ts
```

The regenerated `models.ts` contains tab-only blank lines emitted by the generator, so
`git diff --check` reports trailing whitespace in that file. Normalising it by hand would make
the next regeneration dirty again and would be an edit to generated output, so it is recorded
here as a known, harmless property rather than fixed.

## 5. Verification performed

| Check | Result |
|---|---|
| `build\wails.cmd version` | `v2.12.0`, exit 0 |
| `build\wails.cmd doctor` | `SUCCESS Your system is ready for Wails development!` |
| `build\wails.cmd build` | exit 0, `Built '...\build\bin\howl-chat.exe' in 27.958s`, 21,803,520 bytes |
| Argument forwarding | `build\wails.cmd generate module` exit 0; bindings unchanged |
| Failure propagation | `build\wails.cmd build --definitely-not-a-flag` → CLI error, exit 1 |
| Invocation from outside the project root | `cmd /c C:\...\HOWL_Chat\build\wails.cmd version` → exit 0 from `%TEMP%` |
| Missing project files | wrapper aborts with a clear message, exit 1 |
| No `go` on PATH | wrapper aborts with a clear message, exit 1 |
| Plain Go build | `go build -tags "desktop,production"` exit 0 (29,911,040 bytes) |
| `node --check frontend/wailsjs/go/main/App.js` | exit 0 |
| `git check-ignore -v build/tools/wails.exe` | ignored by `.gitignore:97` |
| Tracked build artifacts under `build/tools`, `build/bin` | none |
| `git diff --check` | only pre-existing and generated-output whitespace (see §4) |
| `go.mod` / `go.sum` after builds | unchanged |

## 6. Native WebView2 isolation — was blocked, now resolved

> **Update (2026-11-10):** resolved by an opt-in `HOWL_WEBVIEW_DATA_DIR` switch that routes the
> WebView2 profile and every application-owned storage root through one validated test root.
> The native smoke test was then run successfully.
> See `docs/review/ENGINEERING-ISOLATED-NATIVE-RUNTIME-2026-11-10.md`. The analysis below is the
> state that motivated that change and is kept for the record.

Wails v2.12 exposes exactly one control for the WebView2 profile:
`options.App.Windows.WebviewUserDataPath` (string), copied into the Chromium loader as
`chromium.DataPath` when the Windows options are set
(`internal/frontend/desktop/windows/frontend.go:475`,
`pkg/options/windows/windows.go:101`). When it is empty the WebView2 default user-data folder
(`%APPDATA%\<binary name>`) is used, which is the user's real profile. There is no environment
variable, command-line flag, or build tag in this Wails version that overrides it, and
`main.go` does not set it.

Consequences, stated plainly:

- A separate executable directory does **not** isolate the profile; neither does a second
  window or a different working directory.
- Launching `build\bin\howl-chat.exe` today would open the user's real HOWL Chat WebView2
  profile and could read or alter stored state.
- Application storage roots are separate from the WebView2 profile and also not configurable
  (`os.UserConfigDir()` in `app-image-workflows.go`, `os.UserHomeDir()` in chat persistence and
  `internal/paths`), so a profile-only override would still leave app data pointed at the
  real folders unless those roots are redirected too.

Therefore the native app was **not launched**, and no native runtime acceptance is claimed.
This is a product-boundary question, not a test-harness detail: it needs an explicit,
reviewable configuration seam.

### Proposed smallest next task (not implemented)

**"Isolated native test mode" — env-gated storage/profile roots.**

1. In `main.go`, when `HOWL_WEBVIEW_DATA_DIR` is set, pass
   `Windows: &windows.Options{WebviewUserDataPath: <that path>}` to `wails.Run`. Unset means
   current behaviour, unchanged.
2. Apply the same gate to the app's own storage root (settings/config, chats, logs) so an
   isolated run cannot read or write real user data. `internal/paths` is the natural single
   place for that seam; the exact owner should be chosen during design, not here.
3. Acceptance for that task: two consecutive launches with a temporary root; the WebView2
   profile and app data appear only under the temporary root; the real
   `%APPDATA%\howl-chat.exe` profile and config are byte-identical before and after
   (hash comparison, no deletion); a normal unset run behaves exactly as today.

Until that task lands, native desktop acceptance for F3-B (and any later stage) stays
**Blocked**, and harness-based evidence remains the strongest runtime evidence available.

*Later note (2026-11-10):* that task landed as `HOWL_WEBVIEW_DATA_DIR` with a required marker root,
and the W27 native run is recorded in `docs/review/ENGINEERING-ISOLATED-NATIVE-RUNTIME-2026-11-10.md`
§27. This section is kept as the state of the record when it was written.

## 7. Limitations

- No TypeScript compiler is installed in the workspace (`typescript` is not resolvable, no
  global `tsc`), so the generated `.d.ts`/`models.ts` pair was checked structurally
  (method sets, namespaces) rather than typechecked. `App.js` was syntax-checked with Node.
- `build\wails.cmd` is new and **untracked**; it must be `git add`ed with the other changes or
  the wrapper is lost on a clean checkout. No commit was made.
- `wails doctor`'s Go row reflects the CLI's own build toolchain, as noted in §3.
- The native WebView2 blocker above is the only remaining obstacle to full desktop acceptance.
  *(Cleared by the isolation task; the native W27 run is recorded in the isolated-runtime report
  §27.)*

## 8. Recommended next step

Land the isolated native test mode described in §6 as its own reviewed change, then run the
native smoke test that F3-B lists as outstanding (import lorebooks, close, reopen, verify
canonical files and the derived cache). Do not start F3-C before that acceptance evidence
exists.
