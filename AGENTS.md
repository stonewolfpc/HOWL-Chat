# HOWL Chat — engineering gates

These are project rules, not suggestions. They apply to every change in this repository.

## 1. Acceptance gate: the user-visible workflow decides

**No feature is complete until its actual user-facing workflow has been exercised in every supported runtime, with the expected result visibly demonstrated.**

- A successful build proves nothing about completion.
- A passing backend suite is supporting evidence, never a substitute for driving the workflow.
- "Implemented" and "tested" are different words; do not use one for the other.
- If the workflow has not been clicked through, say so plainly rather than implying it works.

Ask of every finished slice: what did the user click, in which host, and what was visible afterwards? If that sentence cannot be written with evidence, the slice is not finished.

## 2. Supported hosts

There are two, and a change to a user-facing workflow is verified in both when both are relevant:

| Host | Build | Run |
|---|---|---|
| Packaged desktop app (WebView2) | `go build -tags "desktop,production" -o build/bin/howl-chat-s1.exe .` | `HOWL_WEBVIEW_DATA_DIR=<isolated root> build/bin/howl-chat-s1.exe` |
| Browser dev host (real backend over HTTP) | `go build -tags uiharness -o build/bin/howl-ui-harness.exe .` | `HOWL_WEBVIEW_DATA_DIR=<isolated root> build/bin/howl-ui-harness.exe` — serves `frontend/dist` on `http://127.0.0.1:8799`, every `window.go.main.App.*` call proxied by reflection to the same method the desktop build calls |

The two hosts share `frontend/dist` and the same `App` bindings. Where a workflow needs a native surface the browser host does not have (file pickers, message boxes, the Wails window itself), record the difference explicitly instead of treating the hosts as equivalent. Any other runtime — a plain static file server, a mocked bridge, a harness without the real backend — is not a supported host and cannot be used as evidence.

## 3. Driving and evidence

- Drive with real input in the host: mouse events and screenshots for the packaged app, the browser tooling for the dev host.
- Read the result back from disk as well as from the screen: a UI that claims success while the store says otherwise is a defect, not a pass.
- Use an isolated state root (`HOWL_WEBVIEW_DATA_DIR`, with `.howl-isolated-root` present) for every manual run. Never drive the workflow against real user data.
- Keep destructive testing on disposable fixtures. Never touch drives or folders outside this project's isolated roots.

## 4. Ownership

- One authority per piece of state; one mutation route per authority.
- One mutation route per authority means no silent competing writers, no background overwrite of newer intent, and no automatic restoration of a superseded user decision.
- Prefer unifying or removing an owner over adding a layer. A second index, cache or mirror of the same truth is a defect unless it is explicitly the single authority.
- Every refusal is surfaced. A discarded refusal is non-conforming.

## 5. Reporting

Every report states: current owners; competing owners found; what stays and why; what is removed and what replaces it; what was implemented; what was verified and in which host; verified limitations; and whether user-visible behaviour was demonstrated.
