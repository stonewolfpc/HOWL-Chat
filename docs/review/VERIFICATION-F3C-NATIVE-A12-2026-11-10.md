# VERIFICATION — F3-C Native A12 Prompt Boundary

**Date:** 2026-11-10
**Status:** **A12 PASS.** The production prompt path was captured at the inference boundary in an
isolated native WebView2 process. This is prompt-boundary evidence, not model inference.
**Driver:** `build/tmp/a12-acceptance/run-a12-native.ps1`
**Evidence:** `build/tmp/a12-acceptance/evidence-20261110-134812/`

## 1. Build and isolation

The native binary was built with the test-only tags:

```text
build\wails.cmd build -skipbindings -tags "native_test_seam,promptcapture" -o howl-a12-seam.exe
```

The run used a fresh temporary root outside the project containing `.howl-isolated-root` with
`howl-isolated-test`. `HOWL_WEBVIEW_DATA_DIR` directed all application state into that root. The
native seam bound only to loopback and required `X-Native-Test-Seam-Token`.

The capture route is compiled only with `native_test_seam && promptcapture`. A native-only build
gets a `501` stub; a production build gets neither the route nor the recorder.

## 2. Production path exercised

The driver used real Wails bindings, not a test prompt fixture:

1. `window.go.main.App.SaveChatData(...)` applied the synthetic world, scenario, character, and
   lore entries to the production `App` state and its normal persistence manager.
2. `window.go.main.App.SendMessageStream(...)` entered the production prompt path:
   `SendMessageStream` → `SendMessageStreamWithImage` → `recordUserTurn` →
   `buildPromptMessages` → `httpClient.GenerateChatStream`.
3. The `promptcapture` wrapper sat at that same inference-client boundary and recorded the messages
   without calling a live model.
4. The native seam returned the recorder's in-memory messages; the prompt itself was not read from
   a fixture or reconstructed after the fact.

The production `buildPromptMessages` assembled the system prompt and the user message. The
captured `system` content contained the world, character, scenario, and lore sections in that
order; the captured `user` content contained the user marker and the production system-note
suffix.

## 3. Result

| Check | Result |
|---|---|
| Capture enabled | `true` |
| Inference-boundary calls | `1` |
| Captured roles | `system`, `user` |
| World marker | Present in system prompt at offset 164 |
| Character marker | Present in system prompt at offset 278 |
| Scenario marker | Present in system prompt at offset 497 |
| Lore marker | Present in system prompt at offset 586 |
| User marker | Present in user message |
| Disabled lore marker | Absent |
| Isolated native process | Passed; no live model |

The exact captured messages are in `a12-prompt-capture.json`. The synthetic input is in
`a12-input.json`; it is input to the production state and prompt builder, not an expected-prompt
fixture.

## 4. Production binary check

A separate untagged production build passed:

```text
build\wails.cmd build -skipbindings -o howl-prod-a12.exe
```

Scan of `howl-prod-a12.exe` returned `0` occurrences for all of:

- `__howl_native_test_seam`
- `HOWL_NATIVE_TEST_SEAM`
- `X-Native-Test-Seam-Token`
- `NativeTestSaveChatData`
- `native_test_seam`
- `promptcapture`
- `prompt-capture`
- `GetPromptCapture`

`gofmt -l` on the changed seam files was clean. Focused content-import, worldstore, and chat tests
passed; `go vet ./...` exited 0.

## 5. Limits

- The captured content was synthetic and explicitly marked; no real user content was used.
- No live model was called and no model output quality is claimed.
- The native import UI was not exercised; A9 remains binding-level acceptance only.
- This record closes A12, not overall F3-C acceptance.
- Post-rename cleanup fault injection, fsync/power-loss guarantees, and cross-process writer
  serialization remain outside this evidence.
