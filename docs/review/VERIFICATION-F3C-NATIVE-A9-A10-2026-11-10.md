# VERIFICATION — F3-C Native Acceptance A9 and A10

**Date:** 2026-11-10
**Status:** **A9 PASS; A10 PASS.** A12 passed separately in
`VERIFICATION-F3C-NATIVE-A12-2026-11-10.md`.
**Method:** real Wails bindings in an isolated native WebView2 process, a marked state root outside
the project, a test-only `native_test_seam` build, and byte-level canonical-tree hashing.

## 1. Isolation and harness

Both runs used a fresh temporary root outside the project containing `.howl-isolated-root` with
`howl-isolated-test`. `HOWL_WEBVIEW_DATA_DIR` directed all application state, the WebView2 profile,
chat storage, content libraries, and logs into that root. No production user profile or project
content root was used.

The test binary was built with:

```text
build\wails.cmd build -skipbindings -tags native_test_seam -o howl-a9-seam.exe
```

The native seam bound to loopback and required `X-Native-Test-Seam-Token`. The A9 driver called the
five real content-import bindings through `window.go.main.App`; the A10 driver called the real
chat bindings and used the existing native snapshot endpoint. A10's generated-lore save path is not
exposed by a production Wails binding, so one additional build-tagged seam endpoint invoked the
existing `chat.Service.SaveChatData(..., true)` path. That method and endpoint are absent from the
production build.

## 2. A9 — explicit import and LF-1 migration

Driver: `build/tmp/a9-acceptance/run-a9-native.ps1`
Trace evidence: `build/tmp/a9-acceptance/evidence-20261110-134608/`
Earlier evidence: `build/tmp/a9-acceptance/evidence-20261110-132731/`,
`build/tmp/a9-acceptance/evidence-20261110-133602/`

Seeded Worlds destination plus a legacy source with one valid item, a duplicate-identity pair, and
a folder with no document ID. Browser records were supplied as an explicit request.

| Check | Result |
|---|---|
| Native preview | `ready=2`, `writes=2`, `quarantined=3` |
| Quarantine evidence | two `duplicate_identity`; one `missing_document_id` |
| Browser request integrity | unchanged before/after preview and execute |
| Native execute | `completed_partial`; both importable items committed |
| Canonical files | `world.json`, `world_lore.json`, and media reference created |
| Source tree | byte-identical before/after execute |
| Explicit resume | no-op with the same aggregates |
| Native rollback | `rolled_back`; both journaled creations removed |
| Journal | rolled-back operation listed; no `integrity_failure` |
| Post-rollback re-preview | `ready=2`, `adoptable=0`, `writes=2`; original items are absent |
| Separate idempotency fixture | `wld_a9_idem` committed by one operation, adopted by a later operation; zero writes |
| Conflict preview | one `identity_conflict`; zero writes; canonical file unchanged |
| Final source tree | byte-identical after the full run |

The post-rollback and idempotency claims are separate. The rollback removed
`wld_a9_browser` and `wld_a9_legacy`; a re-preview of the original request then planned two
writes and zero adoptions. The later idempotency fixture created `wld_a9_idem` in operation
`imp_1d7b00d18dc90b59` (committed) and adopted it in operation `imp_196a6176841108a5` (zero
writes). It did not adopt either rolled-back item.

This is native binding and runtime evidence. It does not exercise the import UI, and the separate
post-rename cleanup fault-injection limitation remains source-inspected.

## 3. A10 — chat single-writer runtime

Driver: `build/tmp/a10-acceptance/run-a10-native.ps1`
Evidence: `build/tmp/a10-acceptance/evidence-20261110-134557/`

Canonical world and scenario items were created through the native import bindings. Manual and
generated lore files were then written into those item folders using stable owner IDs. The chat
fixture referenced both canonical IDs.

| Check | Result |
|---|---|
| Canonical items | one world and one scenario item created natively |
| Canonical pre-state | only the two expected item folders; full file hashes captured |
| `LoadChat` | chat ID, world ID, scenario ID, and one history entry restored |
| `SaveChat` | chat-owned data saved; canonical hashes unchanged |
| `SaveChatData` | frontend chat state saved; canonical hashes unchanged |
| Canonical chat load | world and scenario status `available` |
| Generated-lore snapshot | `lb_a10_world_auto`, `lb_a10_scenario_auto` present |
| Generated-lore save | `SaveChatData(..., true)` completed with both entries present |
| Canonical post-state | world and scenario trees byte-identical to the pre-state |
| Chat-owned files | `chat.json`, `history.json`, `memories.json` written under `Chats/A10Chat` |
| Sibling creation | no canonical item or lore sibling was created |

The production call chain was traced directly: `handleSnapshot` calls
`chat.Service.LoadChatData`; `handleSaveChat` calls the test-only wrapper
`Service.NativeTestSaveChatData`, which calls the production `Service.SaveChatData(..., true)`;
that method writes `Chats/A10Chat/chat.json` and `history.json`, calls the production
`memory.Service.SaveChatData`, and runs the production `automatedLoreSyncTarget` allowlist
(character-only). The A10 driver hashed the complete recursive file sets under
`data/content/Worlds` and `data/content/Scenarios`, including hidden files, before the run and
after LoadChat, SaveChat, SaveChatData, canonical LoadChatData, and generated-lore synchronization.
The after-manifests are `a10-loadchat-*`, `a10-savechat-*`, `a10-savechatdata-*`,
`a10-canonicalload-*`, and `a10-generatedlore-*` for both roots.

This is not the desktop prompt builder; A12 covers that boundary separately. It does not claim
model inference.

## 4. Verification after the native runs

| Command | Result |
|---|---|
| `go test ./internal/contentimport ./internal/worldstore ./internal/backend/chat` | all `ok` |
| `go test ./...` | all packages `ok`, exit 0 |
| `go vet ./...` | exit 0 |
| `gofmt -l` on changed Go files | no output |
| Production build `build\wails.cmd build -skipbindings -o howl-prod-a10.exe` | exit 0; binary produced |
| Production scan for `__howl_native_test_seam` | 0 |
| Production scan for `HOWL_NATIVE_TEST_SEAM` | 0 |
| Production scan for `X-Native-Test-Seam-Token` | 0 |
| Production scan for `NativeTestSaveChatData` | 0 |
| Production scan for `native_test_seam` | 0 |

## 5. Remaining limits

- No import UI was wired or exercised; A9 used the real Wails bindings directly.
- No model inference is claimed here. A12 prompt-boundary evidence is recorded separately in
  `VERIFICATION-F3C-NATIVE-A12-2026-11-10.md`.
- The A10 generated-lore path was reached through a build-tagged native seam because no production
  binding exposes `chat.Service.SaveChatData`; the production binary contains neither.
- Post-rename cleanup fault injection, fsync/power-loss guarantees, and cross-process writer
  serialization remain outside this evidence.
- A13 remains a regression gate, not a substitute for A9 or A10.
