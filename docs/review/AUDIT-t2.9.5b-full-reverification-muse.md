# Full Adversarial Re-verification — T2.9.5B Audio Research, Addendum A1, and Howl Chat Application State

**Auditor:** Muse (independent adversarial re-audit).
**Role:** strictly READ-ONLY except for the creation of this artifact. No source, register, freeze,
Git, build, server, model, inference, or measurement action was taken.
**Audit date:** 2026-11-04. Upstream access date: November 2026.

---

## 1. Executive summary

I re-derived the contested claims from primary sources instead of trusting either DeepSeek document:
repository lines at `HEAD 4f5defe`, upstream llama.cpp files fetched at commit
`a34de8618b25dfc32ee43749db47d767a9f0a25d`, Hugging Face Hub API metadata (partly fetched
directly by me, partly via a web-research subagent whose source URLs are recorded), the v0.9.0
release page (read directly), and the working tree.

**Standing verdict: no central conclusion of A1 is overturned.** The routes, handlers, gates,
app-side absence of native audio, model families, and version gap all hold. The DeepSeek audit's
three MAJOR wording findings (F1–F3) are **confirmed**; its F4 licensing re-cut is **confirmed in
substance with two refinements**; F5–F10 are **confirmed** (F5 with an exact-line correction);
**F11 is refuted** — the `llama-mtmd-cli` "not production-ready" label exists at
`tools/mtmd/mtmd-cli.cpp:35-36`, so A1's claim stands and the DeepSeek audit's "UNVERIFIED" is wrong.

**New findings this audit adds** (none of them in A1 or the DeepSeek audit):

| # | Finding | Severity |
|---|---|---|
| D1 | DeepSeek F11 is factually wrong (label exists, `mtmd-cli.cpp:35-36`) | Correction of auditor |
| D2 | `HTTPClient.LoadModel` / `LoadModelWithMMProj` target router-only `POST /models/load`, which the app's own single-model server never serves; `LoadModelWithMMProj` additionally ignores its `mmprojPath` argument. Zero callers of either — dead, misleading code | Medium (latent defect) |
| D3 | Exactly one image per request at the Wails boundary (`app-prompt.go:25`); history builder (`app-prompt.go:620-745`) emits text only, never re-sends images. Multi-image and image-in-history are structurally unsupported | Medium (capability ceiling) |
| D4 | `runtime/llama/llama-connect.exe` (Sep 2 17:14) is a second timestamp outlier nobody flagged; project-root SD-stack DLLs (incl. a *different* 351 MB `ggml-cuda.dll`, Oct 31) are a separate stack from `runtime/llama` — conflation risk for the provenance question | Low / evidence hygiene |
| D5 | `go.mod` declares `go 1.23.8`, installed toolchain is `go 1.27.3` — the "Go never updated" premise is stale; the real question is language-version semantics, not toolchain age | Correction of premise |
| D6 | `LoadModel` kills and respawns the server per model with hardcoded port `8080` (`app-model.go:52-60`); any occupant fails the load. No queue, no port fallback | Low (known UX fragility) |
| D7 | `docs/audio-architecture/README.md` stages 2–5 still read PENDING while ASR/TTS/vision-shipping code exists; its TTS matrix lists Bark, which no loader registers. Stale design doc contradicts shipped state | Low (docs) |

Claim tallies for the A1 × DeepSeek-audit challenge are in §4. The capability matrix is in §5.
Nothing in this artifact authorizes implementation, measurement, or any register change (§11 of the
task; restated in §13 here).

---

## 2. Scope and audit baseline

### 2.1 Source artifacts and hashes (verified before reliance)

| Artifact | SHA-256 (recomputed this pass) | Expected | Result |
|---|---|---|---|
| `docs/review/t2.9.5b-audio-research-correction-addendum-a1.md` | `2eb69c778b6a2737a4d5415e5ab9645b0f4566c365edb75e1e2e4fb9dbfa2bbe` | `2eb69c77…9dbfa2bbe` | MATCH |
| `docs/review/AUDIT-t2.9.5b-audio-research-a1.md` | `502e98e3…381ce86cc8cb` | `502e98e3…381ce86cc8cb` | MATCH |
| `docs/review/README.md` | `68b85f1a…0d71dced` | `68b85f1a…0d71dced` | MATCH |

All three 64-hex values matched exactly, i.e. A1, the DeepSeek audit, and the index were **not edited**. (Line counts: A1 253, audit 272, README 57 — unchanged.)

### 2.2 Repository baseline

- `HEAD 4f5defedb1df037052dd75ada9c5d5dea21f20de` (unchanged throughout; verified before and after).
- Working tree: pre-existing modifications only (e.g. `app-audio.go`, `app.go`, `internal/audio/*`,
  `IMAGE_PAGE_CHANGE_LEDGER.md`, plus untracked `docs/review/*` from the capture session and
  unrelated untracked logs). This audit created exactly one file (this artifact); nothing else was
  touched. Full inventory and this file's hash are reported at close (§15).
- Original research report: still **NOT RECOVERED** (per `docs/review/README.md §3`, whose search
  account I did not need to repeat — I independently confirmed no `*research*`/`*A1*`/`*addendum*`
  file exists and no Pocket/Soprano content exists in any `.md`). All challenges to "prior report"
  claims therefore run through A1 §10's disposition table, marked accordingly.

### 2.3 Method and evidence grades used here

Evidence classes per the task: **Observed** (direct read), **Source-derived** (exact line/window),
**Calculated** (arithmetic on observed values), **Inferred** (stated as such, never as fact),
**Unestablished**, **Superseded**. Claim classes: **Confirmed / Corrected / Disputed /
Unestablished / Superseded**. Upstream files were fetched as raw blobs at the pinned commit over
the network (read path only; no execution, no vendoring into the repo — temp copies lived outside
the project and are removed at close). HF metadata: three repos re-fetched directly by me
(Qwen3-ASR-1.7B, Voxtral-Mini-4B-Realtime-2602, Qwen2.5-Omni-7B, Gemma-4-E2B-it); the remainder is
web-researcher-reported with source URLs, graded as secondary and flagged where it is the sole basis.

---

## 3. Upstream source re-verification (commit `a34de8618`)

All line numbers below are Source-derived from blobs fetched at
`a34de8618b25dfc32ee43749db47d767a9f0a25d` on 2026-11-04.

### 3.1 Route registration — Confirmed exactly

`tools/server/server.cpp:186-189` (Observed 813-line file):

```
186  ctx_http.post("/v1/audio/transcriptions",  ex_wrapper(routes.post_transcriptions_oai));
187  ctx_http.post("/audio/transcriptions",     ex_wrapper(routes.post_transcriptions_oai));
188  ctx_http.post("/tts",                      ex_wrapper(routes.post_tts));
189  ctx_http.ws  ("/v1/realtime",              routes.ws_realtime);
```

Handler declarations `tools/server/server-context.h:124-135`: `post_transcriptions_oai` at 124,
`ws_realtime` at 125, `post_tts` at 135. Handler definitions in `server-context.cpp`: 6080, 6115,
6435 — all exact. A1's retained imprecision (registration cited "at file level") is now closed:
**registration block is `server.cpp:186-189`**. Grade: Confirmed (also confirms DeepSeek §2 row).

### 3.2 Transcriptions gate and formats — Confirmed exactly

- Gate `server-context.cpp:6082`: `if (!meta->has_mtmd || !meta->chat_params.allow_audio)` →
  "The current model does not support audio input." (A1's `has_mtmd && allow_audio` formulation holds.)
- Mandatory multipart `file`: `server-chat.cpp:680-685` (`find("file")`, else throw
  `No input file found for transcription`). A1's "681-686" is a ±1 window of the same statement —
  acceptable range, noted.
- `response_format` json-only: `server-chat.cpp:691-694` (throw otherwise). A1's "691-694" exact.
- Conversion `server-chat.cpp:673` (`convert_transcriptions_to_chatcmpl` def line) exact; ASR prompt
  preset `:695` (`common_chat_get_asr_prompt`) exact.
- `input_audio` parsing `server-common.cpp:1290` (`} else if (type == "input_audio") {`) exact,
  with the `allow_audio` error at 1291 ("you may need to provide the mmproj").

### 3.3 Realtime branch — DeepSeek F1 CONFIRMED (with one nuance added)

`server-context.cpp:6132-6136`:

```cpp
// models that answer after the audio hold a conversation, see server-realtime.h
// live transcription: Voxtral Realtime (no chat), or an "asr" model in a session opened with ?intent=transcription
const bool live_asr = !ctx_server.asr_prompt.empty() && req.get_param("intent") == "transcription";
if (ctx_server.mtmd_cap.audio_stream && ctx_server.mtmd_cap.chat && !live_asr) {
```

A1's wording ("conversational session only if `audio_stream && chat` and `?intent=transcription`
is absent") is **imprecise in exactly the way DeepSeek F1 states**: transcription requires the
conjunction of non-empty `asr_prompt` AND `intent=transcription`. Two concrete consequences A1's
phrasing gets wrong: (a) a chat+stream model contacted with `intent=transcription` but no ASR
prompt still takes the conversational branch; (b) a non-chat streaming model transcribes regardless
of `intent`. **Additional nuance neither document records:** A1's phrasing mirrors the code
*comment* at 6132-6133, which itself uses the `?intent=transcription` shorthand — i.e. the
imprecision is inherited from upstream's own comment, but the boolean on the next line is
authoritative and stricter. Claim: **Corrected** (F1 stands, provenance of the error identified).

### 3.4 Realtime output modality — DeepSeek F2 CONFIRMED

`tools/server/server-realtime.h:33-34`: "responses follow `output_modalities` or the model's own
switching"; `:83-85`: `mtmd_helper::gen_audio gen; // set if the model generates audio`. A1 §4
matrix cell "audio stream + text" therefore **overstates**: a model with `audio_stream && chat`
but without `audio_gen` yields text only. A1's own §4 negatives ("Audio input ≠ full
speech-to-speech… additionally needs `audio_gen`") contradict its matrix cell. Claim:
**Corrected** — output must read "text, and/or audio per `output_modalities`/model switching;
audio out requires `audio_gen`".

### 3.5 TTS gate citation — DeepSeek F3 CONFIRMED

Same file throughout (`tools/tts/tts.cpp`, 223 lines):

| Statement | Lines | Status |
|---|---|---|
| NOT production-ready playground notice | 13-17 (operative sentence at 15) | A1 correct, citation wrong |
| `--mmproj` required | 61-62 | A1 correct, citation wrong |
| `mtmd_gen_audio_get_info(...).type == MTMD_GEN_AUDIO_TYPE_NONE` → reject | 96-97 | A1 correct, citation wrong |
| speaker ref / TTS prompt / `gen.init` | 106-130 | this is what A1 cited for the gates above |

Claim: **Corrected** — split the citation into 15 / 61-62 / 96-97 / 106-130.

### 3.6 mtmd-cli realtime rule — DeepSeek F5 CONFIRMED with tighter lines

Actual lines (`tools/mtmd/mtmd-cli.cpp`): decision `572` (`const bool transcribes = ...`), branch
`574`, single-`--audio` error `576`, dispatch `579`. DeepSeek's "cite 572-579" covers all four;
A1's "576-579" omits the decision. **Corrected** to 572/574/576/579 (span 572-579 acceptable).

### 3.7 mtmd-cli production label — DeepSeek F11 REFUTED (D1)

`tools/mtmd/mtmd-cli.cpp:35-36`:

```cpp
* Please note that this is NOT a production-ready binary.
* It is a playground for trying multimodal support in llama.cpp.
```

A1 §2's "Source-labelled not production-ready" for `llama-mtmd-cli` is **Confirmed**; the DeepSeek
audit's F11 ("label not located… restrict the claim to llama-tts") is **wrong** — the label was
one `grep` away. This is a finding against the DeepSeek audit, not against A1. (Consequence for A1:
none; its table stands as written.)

### 3.8 Remaining upstream citations — all Confirmed

- `has_cap_tts()` → `mtmd_cap.audio_gen` at `server-context.cpp:1034-1036`. Exact.
- `run_audio_stream` at `mtmd-cli.cpp:262` (`static int run_audio_stream(...)`). Exact.
- CMakeLists `${BUILD_NUMBER}` plumbing at root `CMakeLists.txt:155-156`
  (`set(LLAMA_BUILD_NUMBER ${BUILD_NUMBER})`), forwarded to ggml at `:223`. So the numeric build
  id is injected by CI — DeepSeek F9's caution stands (see §3.9).
- `docs/multimodal.md` at this commit documents: Gemma 4 audio (98-101, 137-140, "audio input,
  vision input"), Ultravox 0.5 (108-110), Voxtral (116-117), Qwen2.5-Omni (127-130, audio+vision),
  Qwen3-Omni **both** Instruct and Thinking (132-135). A1's "Qwen3-Omni documented" holds, and the
  Thinking variant resolves a further ambiguity A1 never raised.
- TTS README at this commit documents: Qwen3-TTS (22-37), Pocket TTS with per-language-dir
  conversion (38-60, upstream `kyutai/pocket-tts`), Soprano-80M **and** Soprano-1.1-80M with
  `decoder.pth` conversion (63-85), LFM3-Audio (88-111). No Chatterbox/Nemotron/predefined-voices
  strings — consistent with A1's "arrive in v0.9.0" and with the v0.9.0 notes I read directly (§8).
- Upstream LICENSE is MIT ("Copyright (c) 2023-2026 The ggml authors") — via web-researcher read
  of the repo LICENSE; secondary basis, consistent with the known ggml license text.

### 3.9 Build number `b12110` — DeepSeek F9 CONFIRMED (still Unverified)

`b12015` is Observed (server fingerprint in this project's own logs, e.g.
`build/t1/harness-t13.log` chunks carrying `"system_fingerprint":"b12015-a34de8618"`, and
`STAGE_STATUS.md:194`). `b12110`/`cd959bae…` for v0.9.0 is Reported (release page read directly
confirms the tag `cd959ba` and date 4 Nov 2026; the nightly number `b12110` is per the DeepSeek
audit's `git ls-remote` work, not re-observed by me). The "95 increments" subtraction is
Calculated-correct *conditional* on both inputs. Claim stays **Unestablished** as an observed fact;
the upgrade recommendation does not depend on it (§8).

---

## 4. A1 × DeepSeek-audit challenge — full disposition

| # | Subject | A1 claim | DeepSeek audit | This audit |
|---|---|---|---|---|
| F1 | `/v1/realtime` branch | imprecise (`intent` as switch) | MAJOR, corrected | **Confirmed** + comment-provenance nuance (§3.3) |
| F2 | realtime output "audio stream + text" | overstates | MAJOR, corrected | **Confirmed** (§3.4) |
| F3 | `llama-tts` gate lines 106-130 | wrong range | MAJOR, corrected | **Confirmed**; exact lines in §3.5 |
| F4 | license grouping (Qwen custom; Voxtral/Ultravox unestablished) | needs re-cut | MAJOR | **Confirmed in substance**, two refinements: (a) Gemma 4's `license_link` targets `ai.google.dev/gemma/docs/gemma_4_license` (researcher-verified; my direct fetch confirms `license:apache-2.0` tag) — the SPDX label alone is insufficient, which is a sharper reason than A1's hedge; (b) `Qwen2.5-Omni-7B` is `license:other` (**directly verified** by me), Qwen3-ASR-1.7B and Voxtral-Realtime-4B are `apache-2.0` (**directly verified** by me); remaining rows rest on researcher HF-API reads (§10) |
| F5 | mtmd-cli rule 576-579 | 4 lines late | MINOR | **Confirmed**; exact decision lines 572/574/576/579 |
| F6 | "streaming true for Whisper and Qwen" | overgeneralises (Small is false) | MINOR | **Confirmed** (`asr_models.go:50` T true, `:79` B true, `:108` Small **false**, `:137`/`:166` Qwen true) |
| F7 | `app-model.go:117-126` mmproj block | range off by one | MINOR | **Confirmed**; actual block 118-125 (`mmprojPath :=` at 118, `HasVision` at 120/124) |
| F8 | `app-audio.go:107-128` hardcoded flags | range imprecise | MINOR | **Confirmed**; func starts 105, nil-branch literals 110-111, normal-branch 127-128 (with `multimodal_ready:false`, `realtime_ready:false` at both sites) |
| F9 | `b12110` / "95 increments" | reported-not-verified | UNVERIFIED | **Confirmed** as Unestablished (§3.9) |
| F10 | four families' mmproj quant variants | E1-unreverified | UNVERIFIED | **Confirmed** as Unestablished — I did not re-fetch those four trees either; researcher data exists (§10) but is secondary |
| F11 | mtmd-cli "not production-ready" | A1 claims label | UNVERIFIED | **Refuted — label exists at mtmd-cli.cpp:35-36** (D1). A1 Confirmed |
| A1-§1.3 | zero app-side route calls | E2 | confirmed | **Confirmed with one qualification**: no *audio-endpoint* calls exist (`input_audio`, `/tts`, `/v1/realtime`, transcriptions — zero hits in `internal/backend/llama`), but audio-*adjacent* plumbing exists: `mmproj` auto-detect + `--mmproj` spawn (`http_client.go:288-294`, `app-llama.go:196-201`), which is vision-capability plumbing that an audio projector would also satisfy by filename. The "no native audio path" conclusion stands; the sentence "no audio … request code at all" should say "no audio-*endpoint* request code" |
| A1-§5 | bundled build `b12015-a34de8618` 2026-10-30 | E1/E3-I | confirmed | **Confirmed** (fingerprint Observed in logs; `--help` record at STAGE_STATUS:194 secondary) |
| A1-§11.1 | E3-I linkage, no SBOM | declared | accepted | **Confirmed** — nothing found that upgrades it (no `.src_sha`, build receipt, or SBOM anywhere in the repo; `.sd-upgrade/` concerns the SD stack only) |
| A1-§11.8 | whisper redistribution | unestablished | resolved: no binary bundled | **Extended**: still no `**/whisper*` binary anywhere outside `build/` (Observed); `BinaryDownloader.Acquire/DownloadLatest` (whisper GitHub-release ZIP + unix build-from-source paths, `downloader_acquire.go:17,29,45,132-145`) has **zero production callers** (Observed, §10) — the download-on-behalf question narrows to dormant code, but the code's *existence* keeps the distributor question alive if ever wired up |

No other A1 row was overturned. Tallies: A1 claims re-checked ≈ 45 (31 exact-confirms and 5
range-imprecisions per the DeepSeek audit all re-observed by me at the cited lines); wording
corrections required: F1, F2, F3, F4-table (4, of which F11-removal nets against the DeepSeek
audit); unverified remainders: F9, F10, plus §14 items.

---

## 5. Application capability matrix (UI → runtime, observed at HEAD)

Grades: Working / Constrained / Partial / Runtime-only / UI-exposed-unsupported / Documented-only /
Not-implemented / Broken / Unestablished. "E2E" means end-to-end verified — nothing below claims it
(no execution was run or authorized).

| Capability | Grade | Limiting layer + evidence |
|---|---|---|
| Ordinary text chat | Working (not E2E here) | `llama-server` single-model spawn `-m` (`app-llama.go:182-191`); chat completions path in `http_client.go`. Not re-probed |
| Vision, single image, current turn | Constrained | Wails boundary takes exactly one `imageData` (`app-prompt.go:25`); `buildRequestWithImage(model, prompt, options, stream, imageData)` (`http_client.go:576`); single `image_url` part (`:582`, `:712-713`); vision gate `HasVision` from mmproj presence (`app-model.go:118-125`, `app-prompt.go:38`); `--mmproj` passed at spawn (`app-llama.go:196-201`) |
| Multiple images per request | Not-implemented | Structural: one `[]byte` through `SendMessageStreamWithImage` → `GenerateChatStreamWithImage` → one content part. Upstream accepts multiple `image_url` parts; the app never builds them |
| Image dimensions/limits/preprocessing | Unestablished | No resize/compress/size-cap code found on the app path (`buildRequestWithImage` base64-encodes as `image/jpeg` unconditionally — note: PNG/WebP bytes are labelled `data:image/jpeg`, Observed at `:582-583`/`:712-713`); any effective limit would come from server-side handling, unprobed |
| Image in conversation history | Not-implemented | `buildPromptMessages` (`app-prompt.go:620-745`, 126 lines) emits **zero** image/vision parts (Observed) — history is text-only; a second turn cannot see the first turn's image |
| Audio transcription (local) | Partial | Whisper `Transcribe` file path exists (`recognizer.go`); UI records → temp file → `TranscribeAudio` → `insertTranscriptIntoChat` (`audio-service.js:114-129,157`); but registry approves only tiny/base/small while UI offers medium/large (`audio-settings.html:96` vs `model_resolver.go:69-76`) — UI-exposed-unsupported for medium/large |
| Live/streaming transcription | UI-exposed-unsupported | `TranscribeStream` returns "not supported" in Whisper (`recognizer.go:196-202`) and Qwen (`client.go:95-101`); registry declares `SupportsStreaming:true` for Tiny/Base/Qwen (`asr_models.go:50,79,137,166`); `GetAudioStatus` hardcodes `realtime_ready:false` (`app-audio.go:110-111,127-128`) |
| Direct audio understanding (`input_audio`) | Not-implemented | Zero `input_audio` construction in `internal/backend` (Observed); chat parts carry images only (`app-prompt.go:38` vision gate) |
| Realtime spoken conversation | Not-implemented | No WebSocket client in `internal/backend/llama` (Observed: no `gorilla/websocket` use there; the indirect dep in go.mod serves other code); status hardcodes `realtime_ready:false` |
| TTS (native `audio_gen`) | Not-implemented | No `/tts` caller; engines are Piper/Kokoro/MeloTTS/browser (`pipeline.go:444-446`, `audio-settings.html:153`) — independent of llama.cpp, as A1 states |
| Video / multimodal beyond image+text | Not-implemented | Omni video unsupported per deferred note D (see §7); no video path in `http_client.go` |
| Model discovery / capability detection | Partial + misleading | `findMMProj` + `HasVision` boolean conflates vision/audio projectors by filename (A1 §3 standing rule; Observed at `app-model.go:118-125` and the IQ/Q/K candidate list at `http_client.go:323-360`); audio has no capability field at all |
| Auto/manual provider selection | Partial + defects | Auto never filters by *task* (A1 §7); explicit selection can choose medium/large/Qwen-API engines the resolver rejects (§5 row 6); locality switch (local↔DashScope) undisclosed (`client.go:29,33`) |
| Resource mgmt / loading / cancellation | Constrained | Per-model kill+respawn (`app-model.go:52-78`); port 8080 exclusive + hard fail (`:56-60`); stream cancellation exists (`app-prompt.go:54-55` `streamCancel`); VRAM arbitration file exists (`app-vram-arbitration.go:43` CUDA-context overhead) — behaviour unprobed |
| Conversation persistence | Working (not E2E here) | `internal/backend/chat/` service + `persistence.go` `PersistenceManager`; user data dirs (`Chats/`, `Characters/`…) Observed at root |
| Whisper binary acquisition | Dormant code | `downloader.go`/`downloader_acquire.go` implement release-ZIP + build-from-source fetch; `findWhisperCLI` "never downloads or builds" (`recognizer.go:226-228`); zero callers of `Acquire`/`DownloadLatest` (Observed) |

---

## 6. Backend migration audit: Gollama → llama.cpp

### 6.1 Residual references (Observed, all read in full context)

1. `internal/backend/prompt/builder.go:6` — "application-level prompt construction that is not part
   of gollama.cpp" (comment only).
2. `internal/backend/types/model_profile.go:97` — "`RuntimeSettings` … passed to gollama.cpp"
   (comment only; the struct actually feeds `spawnLlamaServerWithModel` args).
3. `internal/backend/types/model_profile.go:144` — "not passed to gollama.cpp" (comment only).
4. `frontend/dist/js/tutorials-data-model.js:54,79` — user-facing tutorial prose: "auto-managed by
   gollama", "The gollama backend may auto-manage this", "use the gollama CLI". **User-visible.**
5. `frontend/dist/index.html` + `.freebuff/backup/index.html` — matched the case-insensitive
   `gollama` sweep; dist bundle strings (build artifact, user-visible in the shipped page).

No functional Gollama plugin import, cgo binding, or Gollama-API call shape was found in Go code:
the backend is `exec`+HTTP against `llama-server` (`app-llama.go:252-257`, `http_client.go`).
Grades: references 1-3 cosmetic-dead (docs debt); 4-5 user-visible stale prose (should be corrected
whenever tutorials are next touched — no doc rewrite authorized here).

### 6.2 Capabilities lost, constrained, or left behind (Source-derived)

1. **Single-image ceiling (D3).** Whether or not Gollama-era vision took multiple images is not
   established in this tree (no pre-migration code remains to compare), but the current ceiling is
   structural and located exactly: one `imageData: string` Wails parameter (`app-prompt.go:25`) →
   one `[]byte` → one `image_url` part (`http_client.go:576-583`, `:703-713`) → text-only history
   (`app-prompt.go:620-745`). Upstream `server-common.cpp:1282-1288` loops over *parts*, so the
   server side is not the limiter. Any "vision regression" investigation starts and ends in these
   three app functions — **no rebuild of vision is justified**; the fix is additive multi-part
   construction plus history carriage (design + authorization required, not given).
2. **Image MIME mislabel.** `buildRequestWithImage` and `GenerateChatStreamWithImage` hardcode
   `data:image/jpeg;base64,` regardless of actual bytes (`http_client.go:582-583`, `:712-713`).
   Upstream sniffs media by content (`handle_media`), so this may be benign today, but it is a
   latent correctness defect with a two-line fix surface.
3. **Dead router-mode code (D2).** `HTTPClient.LoadModel` POSTs `{base}/models/load`
   (`http_client.go:275-313`), an endpoint that exists upstream **only in router mode**
   (`server.cpp:529-533` inside `server_run_router`). The app always spawns single-model servers
   (`spawnLlamaServerWithModel`, `-m`, no router flag). `LoadModelWithMMProj` (`:393-397`) ignores
   its own `mmprojPath` argument and delegates to `LoadModel`. Zero callers of either (Observed).
   This is either a Gollama-era leftover shape or aspirational code — either way it is dead,
   misleading, and would 404 if ever called. Removal is a future cleanup decision, not this audit's.
4. **mmproj kind-blindness (confirms A1 §3 rule).** `findMMProj` matches ~30 vision-quant-style
   names (`http_client.go:323-360`) then falls back to wildcard scan; any audio-only projector file
   satisfying the pattern would set `HasVision=true`. With an audio model loaded, vision-gated UI
   would misreport. No audio capability field exists anywhere in the profile types.
5. **Port exclusivity (D6).** `LoadModel` requires port 8080 free and hard-fails otherwise
   (`app-model.go:56-60`); model switch = kill + respawn (`:72-78`). Ordinary chat is therefore
   briefly *unavailable* on every switch — the resource-protection direction A1 §7 asks for (never
   evict active chat for audio) has no mechanism yet.

### 6.3 What the migration genuinely bought (to avoid a one-sided account)

Single-model `llama-server` exposes `reasoning_content` streaming (used at `app-prompt.go:57-89`
with a think-splitter), `--reasoning-budget` controls (`app-llama.go:214-230`, verified against
this binary's `--help` per STAGE_STATUS:190-199), `/props`, `/v1/models` (used by
`GetLoadedModel`, `http_client.go:405`), and multi-architecture GGUF support the old plugin
lacked. Recommending against llama.cpp would require evidence this audit did not find.

---

## 7. Documentation adversarial review (abridged to material contradictions)

1. `docs/audio-architecture/README.md` vs shipped code: ASCII diagram marks Stages 2 (ASR), 3
   (TTS), 4 (multimodal), 5 (realtime) **PENDING** (lines 64-80) while `internal/audio/asr`,
   `/tts`, `TranscribeAudio`, and single-image vision all exist. The doc describes a May-2026 plan
   as current status. Its TTS matrix lists **Bark** (`:207`) — no loader registers Bark
   (`pipeline.go:438-446` registers Piper/MeloTTS/Kokoro only). Its ASR matrix sizes (39/74/244 MB
   at `:195-197`) describe upstream *bin* weights, while the resolver demands `ggml-tiny.bin` etc.
   (`model_resolver.go:72-76`) — same family, un acknowled redistribution/package gap (see §10).
   **Stale; do not cite for capability claims.**
2. `docs/deferred-model-ecosystem.md` remains the soundest audio-adjacent doc in the tree (role-aware
   resolution, §1.6 non-assumptions). It is explicitly deferred/unauthorizing — no conflict with the
   register.
3. `frontend/dist` tutorial prose references gollama CLI/backend behaviour (§6.1 item 4) —
   user-visible staleness inside the shipped bundle.
4. STAGE_STATUS self-records a Whisper re-init defect (F1, `STAGE_STATUS.md:1690-1700+, latent temp-dir
   orphan on repeated `Initialize`) and invented `default-tts-voice` options (`:1014`). Both are
   honestly recorded as defects — credit where due — but they corroborate §11 items T1/T2.
5. No authoritative-document map exists; `STAGE_STATUS.md`, `asr-repair-design.md` (frozen design),
   `docs/audio-architecture/*`, and `docs/review/*` overlap without a stated precedence. I do **not**
   create an index (task §5 forbids another index system); I record that the register (r14) is the
   only controlling artifact and the audio docs are subordinate/uncontrolled.

---

## 8. November 2026 upstream comparison (primary sources read directly unless noted)

- **Releases.** v0.9.0 tag page read directly (4 Nov 2026, `cd959ba`): audio-relevant additions are
  **Chatterbox TTS support, Nemotron 3.5 Transcribe (offline, 40 languages, prompt context, speaker
  attribution), Pocket TTS predefined voices**, plus server `/props` predefined-voice reporting and
  TTS-endpoint predefined-voice acceptance; `MTMD_GEN_AUDIO_TYPE_CHATTERBOX` + voice APIs in
  `mtmd.h`; `MTMD_VIDEO_TEMPLATE_QWEN3VL`; ggml **v0.29.0** with a batched tensor-copy API
  *superseding the `_2d` API* (direct-linkage deprecation risk); vendor miniaudio 0.11.26,
  cpp-httplib 0.64.0. A1's "only in v0.9.0" list (Chatterbox, Nemotron 3.5, predefined Pocket
  voices) is **Confirmed**; none is required by the Whisper + native-audio-input scope (A1 §5
  recommendation stands).
- **Distance.** b12015 (30 Oct 2026, commit `a34de86`) → v0.9.0 (`cd959ba`) is a small additive step;
  v0.9.0 is a semver MINOR (no `llama.h` break at that tag — release-notes-derived, secondary
  inference from the project's stated versioning rules). Upgrade only against a demonstrated gap;
  any upgrade must re-check `/props` JSON shape (new voices field) and any direct `ggml_*_2d`
  copy call sites. No newer tag than v0.9.0 was visible at access (researcher-reported; I did not
  re-list tags — Unestablished beyond v0.9.0).
- **GGUF/mmproj formats.** No container or projector format break was stated in the v0.9.0 notes
  (Observed absence in the notes I read; not a proof of absence elsewhere — graded Inferred-safe).
- **Security/maintenance.** The release notes are feature notes, not an advisory feed; no CVE sweep
  was performed (out of read-only scope without a dependency manifest audit — recorded Unestablished
  in §14).

---

## 9. Go toolchain assessment

- **Declared:** `go 1.23.8` (`go.mod:3`, Observed). Deps include `wails/v2 v2.12.0`, `gopsutil/v3
  v3.24.5`, `gguf-parser-go v0.24.0`, `gorilla/websocket v1.5.3` (indirect), `x/sys v0.33.0`.
- **Installed:** `go1.27.3 windows/amd64` (Observed via `go version`; no build was run).
- The task's premise ("Go never updated since inception") is **stale as a toolchain claim** (D5):
  the installed toolchain postdates the directive by four minor lines. The *real* questions are
  (a) the `go 1.23.8` directive pins pre-1.24 language semantics/GODEBUG defaults even under a 1.27
  toolchain (secondary-source semantics, standard Go behaviour — treat as guidance, not finding);
  (b) `go vet`'s `stdversion` check under newer toolchains will flag stdlib symbols newer than the
  directive (actionable if the directive is ever bumped); (c) Wails v2.12.0 currency against 1.27
  is unassessed here (no CI config was found to read — `ls` shows `.github/` exists; its workflows
  were not enumerated — Unestablished, §14).
- **No upgrade is recommended or authorized by this audit.** If ever pursued: stage separately from
  any llama.cpp move; first enumerate `.github/` workflows and `wails.json` build matrix (not done
  here).

---

## 10. Licensing, redistribution, and legal dependencies

Grades here are evidence grades, **not legal conclusions**. No compliance/non-compliance verdict is
rendered beyond what the artifacts prove.

### 10.1 Runtime code notices — Confirmed gap (Observed)

`runtime/llama/` contains exactly one notice file: `LICENSE-LLVM-OpenMP` (19,741 B, Observed by
directory listing). No llama.cpp MIT text, no ggml notice, no third-party bundle (miniaudio,
cpp-httplib — both vendored upstream per the v0.9.0 notes I read) accompanies ~40 shipped
`llama-*.exe`/`llama-*-impl.dll`/`ggml-*.dll`/`mtmd.dll`/`llama.dll` binaries. Upstream llama.cpp
is MIT ("Copyright (c) 2023-2026 The ggml authors" — secondary basis), whose terms require the
copyright + permission notice to accompany copies. **The obligation-to-ship-notice is established
as a packaging defect at E1/E2; whether any downstream/Projects-terms cure it was not examined.**
A1 §8's conclusion stands and is, if anything, understated: it omits the miniaudio/httplib
second-order notices.

### 10.2 Bundled vs downloaded binaries (Observed)

- No `whisper`-named binary exists anywhere outside `build/` (Observed). The app does not
  redistribute whisper.cpp. **Confirmed.**
- `BinaryDownloader` (`downloader.go:21-35`) + `Acquire` (`downloader_acquire.go:29-145`)
  implement: Windows — fetch latest whisper.cpp release ZIP via GitHub API (`:17,:95-117`); Unix —
  `git`+`cmake`/`gcc` build-from-source (`:132-145`). **Zero production callers** of
  `Acquire`/`DownloadLatest` (Observed). Dormant code — but if ever wired to a UI affordance, the
  app becomes a download-on-behalf distributor of MIT-licensed binaries (notice duty follows the
  bits). Record, do not resolve.
- `runtime/llama` ships version-mixed neighbours: `llama-connect.exe` Sep 2 17:14 amid an otherwise
  Oct-30 set (D4, Observed) — same class of anomaly as `ggml-cuda.dll`, previously unremarked.

### 10.3 ggml-cuda.dll provenance (D4 extended — still Unestablished, now bounded)

Observed facts: `runtime/llama/ggml-cuda.dll` 594,123,776 B, Oct 30 15:48:38 vs set Oct 30
10:46-10:49; root-level SD-stack `ggml-cuda.dll` 350,998,016 B, Oct 31 17:10 (a *different* file
for a different stack — do not conflate); `.sd-upgrade/sd-new.zip` (354 MB, Oct 31) + `extracted/`
+ `backup-37bc88b/` receipts concern the SD stack, not llama; no `.src_sha`, `convert.log`,
SBOM, or build receipt exists for the runtime set. What would establish provenance/compat (A1 §11.2
still open): (1) `llama-server --version` + `--build-info` output on this machine (fingerprint only
proves the server binary's self-report, not DLL ABI); (2) a recorded source URL + hash for the
594 MB file (official llama.cpp CUDA release asset vs third-party build); (3) a load smoke test
(CUDA path init log) — execution, not authorized here. Timestamp alone proves nothing either way;
the *size delta across stacks* explanation (SD 351 MB vs runtime 594 MB) is now documented so a
future investigator does not chase the wrong delta.

### 10.4 Model-weight licenses (direct + secondary)

Directly verified by me via Hub API (Observed, Nov 2026): `Qwen/Qwen3-ASR-1.7B` → `apache-2.0`;
`mistralai/Voxtral-Mini-4B-Realtime-2602` → `apache-2.0` (arch `VoxtralRealtimeForConditionalGeneration`,
13 languages); `Qwen/Qwen2.5-Omni-7B` → `other` (custom Qwen terms); `google/gemma-4-E2B-it` tags
include `license:apache-2.0` (cardData `license_link` per researcher targets the Gemma 4 licence
page — the SPDX tag alone must not be relied on). Secondary (researcher HF-API reads with URLs, §2.3):
Qwen3-ASR-0.6B `apache-2.0`; Qwen3-TTS-Base `apache-2.0`; Qwen3-Omni-Instruct `other`; Voxtral-Mini-3B
upstream `apache-2.0`; Ultravox GGUF repo `mit`; Gemma-4-E2B-it GGUF `apache-2.0`; LFM3-Audio-GGUF
`other`/`lfm1.0` with LICENSE file; Soprano repos `apache-2.0` (TTS systems with `decoder.pth`,
not standalone vocoders); pocket-tts `cc-by-4.0` + **gated** (auto-approval) — redistribution of
*that* weight file needs per-family review; several ggml-org convert repos carry **no license
field at all** (Qwen3-ASR ×2, Qwen3-Omni ×2, Qwen3-TTS, Voxtral-3B) — derivative-terms gap A1 did
not enumerate. The ggml-org mmproj quant matrix is inconsistent across families (researcher file
lists) — A1 §3/F10's per-family quant cells therefore cannot be promoted past E1 without per-tree
reads; F10 stands.

---

## 11. Technical debt and architectural findings

Root-cause orientation; COPY→PASTE→EDIT→VERIFY→DELETE-THE-OLD applies to any future modularization
(untested replacement must never displace working code).

- **T1. Capability inference from filenames (root cause of §5 rows 6,12).** `HasVision =
  mmproj-exists` (`app-model.go:118-125`); ~30-name candidate list + wildcard (`http_client.go:
  318-391`); no audio-capability field. Fix direction: negotiated-capability record (`/props` +
  model family knowledge), never filename. Regression coverage needed: audio-mmproj-present ⇒
  vision still false.
- **T2. Declared-but-unsupported streaming.** Registry `SupportsStreaming:true` vs both adapters'
  explicit rejections (`asr_models.go:50,79,137,166` vs `recognizer.go:196-202`,
  `qwen/client.go:95-101`). One-line-per-row correction surface; needs a contract test pinning
  registry-vs-adapter agreement.
- **T3. Status/registry/UI triple disagreement.** UI offers medium/large/Qwen-API
  (`audio-settings.html:96`); resolver approves tiny/base/small (`model_resolver.go:69-76`);
  status hardcodes `multimodal_ready:false, realtime_ready:false` (`app-audio.go:105-135`).
  Three sources of truth; the UI is the one the user trusts. Fix direction: single capability
  source feeding all three (this is A1 Tranche-A item 1 — design/authorization required).
- **T4. Dead router-mode surface (D2).** `LoadModel`/`LoadModelWithMMProj` (`http_client.go:
  275-313,393-397`): wrong endpoint family for this app's spawn model + ignored argument. Delete
  (or wire) only under a cleanup authorization with a no-caller proof re-check.
- **T5. Dormant downloader ( §10.2).** Keep-or-remove decision with a notice-duty analysis; do not
  half-wire it to a button.
- **T6. Single-image/text-only-history ceiling (D3).** Located exactly (§6.2.1); additive fix only.
- **T7. JPEG MIME hardcode (§6.2.2).** Sniff-or-label fix at two call sites.
- **T8. Kill-respawn + port-8080 exclusivity (D6).** No queue/fallback; model switch drops chat
  availability. Fix direction belongs with resource arbitration (A1 §7 resource rule).
- **T9. Repeated Whisper `Initialize` orphans temp dir** (STAGE_STATUS F1, honestly self-recorded;
  `recognizer.go` `Release`/`tempDir` handling). One-shot init guard proposed there; not implemented.
- **T10. Gollama prose in shipped bundle + stale audio-architecture docs (§§6.1, 7).** Docs debt;
  no runtime effect except tutorial text.

Unavoidable debt: none identified that must permanently remain; every item above has a bounded fix
surface. Largest sequencing risk is T3 (single capability source) because it touches UI contract.

---

## 12. Risks and dependencies

1. **Provenance risk** (runtime DLL set): one unexplained outlier (`ggml-cuda.dll`) + one
   unremarked (`llama-connect.exe`, Sep 2). Any audio/load claim built on this set inherits E3-I.
2. **License-notice risk**: MIT-notice absence is a fact, not a suspicion (§10.1). Cure is packaging,
   not research.
3. **Weight-terms risk**: gated `kyutai/pocket-tts` (cc-by-4.0), `other`-licensed Omni/LFM3 families,
   license-less ggml-org converts. Per-family review required before any bundling or
   download-on-behalf feature.
4. **Capability-misreport risk**: T1/T3 mean the app can promise (UI) what the stack cannot do
   (streaming, medium/large local ASR, multimodality) while hiding what it could do (none today for
   audio, by design absence rather than failure).
5. **Upgrade risk**: v0.9.0 is additive but touches `/props` shape, voice APIs, and ggml copy API —
   any adoption needs its own verification plan (A1 §5 agrees; I add the `_2d` deprecation).
6. **History-carriage risk**: text-only history (§5) silently drops image context on turn 2 — a
   correctness trap for any future multi-turn vision work, independent of audio.

---

## 13. Recommended sequencing (recommendations only — zero authorization)

1. **Correct the four A1 paragraphs** (F1 branch wording, F2 output cell, F3 split citation,
   F4 table re-cut) *in a new errata note that cites A1 without editing it* — A1 must stay
   byte-identical (its hash is the review chain's anchor).
2. **Packaging cure**: ship the llama.cpp/ggml MIT notice + vendored-lib notices with the installer
   (no upgrade, no code change needed to start).
3. **Provenance closure**: record source URL + hash for the 594 MB `ggml-cuda.dll` or replace it
   from a pinned official asset; note `llama-connect.exe` disposition in the same record.
4. **Contract-first Tranche A** per A1 §6 (unchanged by this audit), with T1/T2/T3 as its
   defect backlog and D2/D3 scoping its ceiling.
5. **Docs hygiene**: mark `docs/audio-architecture/README.md` superseded-or-current (one-line
   status header, not a rewrite) and fix the two gollama tutorial strings when tutorials are next
   touched.
6. **Go directive**: leave `go 1.23.8` until a staged toolchain pass with CI enumeration; do not
   bundle with any llama.cpp move.

---

## 14. Explicitly unresolved questions

1. Original research report text — NOT RECOVERED; A1 §10 dispositions are the only record.
2. `b12110` as observed fact (F9); any tag newer than v0.9.0.
3. Four families' mmproj quant cells (F10) on primary reads (researcher data is secondary).
4. `ggml-cuda.dll` + `llama-connect.exe` provenance/ABI (needs §10.3 items 1-3; item 3 is execution).
5. Whether any caller could reach `BinaryDownloader.Acquire` via UI/script (zero Go callers observed;
   frontend invocation not exhaustively traced — `audio-settings` JS was read for transcript/TTS
   paths only).
6. Wails v2.12.0 × Go 1.27 currency; `.github/` CI matrix contents; `wails.json` build settings.
7. Pre-migration (Gollama-era) capability set — no pre-migration code in tree to compare; "vision
   worked" is testimonial, mechanism unrecoverable here.
8. Effective image size/format limits end-to-end (server-side behaviour unprobed by design).
9. Any advisory feed for llama.cpp/whisper.cpp/Go CVEs relevant to pinned versions.
10. Per-family redistribution terms for gated/custom/license-less weight repos (legal reading, not
    SPDX lookup).

---

## 15. Source citations and integrity close

### 15.1 Exact citations relied on (all read this pass)

**Upstream** (`a34de8618`, raw blobs): `server.cpp:186-189,529-533`; `server-context.h:124-135`;
`server-context.cpp:1034-1036,6080-6095,6115-6165,6132-6136,6435-6450`; `server-chat.cpp:673-700`;
`server-common.cpp:1290-1300`; `server-realtime.h:33-34,83-85`; `mtmd-cli.cpp:35-36,262,572-579`;
`tts.cpp:13-17,61-62,96-97,106-130`; `docs/multimodal.md:98-140`; `tools/tts/README.md:22-111`;
root `CMakeLists.txt:155-156,223`; v0.9.0 release page (4 Nov 2026, `cd959ba`) read directly.
**Repo**: `app-prompt.go:25,38,51-55,57-89,620-745`; `app-model.go:26-78,108-135`;
`app-model-template.go:17`; `app-llama.go:168-257,259-272`; `internal/backend/llama/http_client.go:
55-59,275-313,318-391,393-453,576-583,703-713,768-769`; `internal/backend/chat/` (service,
persistence, `image_compat.go`); `internal/backend/vision/vision.go` (exists; behaviour not
audited); `internal/backend/prompt/builder.go:4-8`; `internal/backend/types/model_profile.go:
97,144`; `internal/audio/{pipeline,asr/{whisper/{recognizer,downloader,downloader_acquire},qwen/client},types/asr_models,model_resolver,hw,manager}.go` at cited lines;
`frontend/dist/{audio-settings.html:96,153,audio-service.js:114-157,js/tutorials-data-model.js:54,79,index.html}`;
`runtime/llama/` full listing + root `*.dll` listing; `go.mod:3,9`; `build_howl_cgo.bat`;
`STAGE_STATUS.md:188-200,1014,1454-1700`; `backups/BACKUP-2026-11-01.md:38-40`; `docs/review/*`
(3 files); `docs/audio-architecture/README.md:64-80,195-207`; `docs/deferred-model-ecosystem.md`
(supporting, no conflict); `docs/t2.9.5b-authorization-decision-register-r14.md:3-20,30,40,115,121`.
**Hub** (direct): `Qwen/Qwen3-ASR-1.7B`, `mistralai/Voxtral-Mini-4B-Realtime-2602`,
`Qwen/Qwen2.5-Omni-7B`, `google/gemma-4-E2B-it` API records. **Secondary** (researcher URLs on file
in the working session): remaining HF repos, Go release notes, whisper.cpp releases, v0.9.0 tag
metadata, llama.cpp LICENSE text.

### 15.2 Register preservation (per task §8 — independently verified, not altered)

r14 (DRAFT — NOT FROZEN — NOT AUTHORIZED) states: CL1/C1/CL2 closed as recorded; **CL3 OPEN**
(`:115`); **IR1 `[I]`** requiring separate disposition (`:30`); **MEAS1 OPEN**, execution not
authorized (`:19`); **T6 OPEN/BLOCKED** (`:20`); gate **NOT AUTHORIZED** (`:15`). This artifact
changes none of it, authorizes nothing, and treats the audio scope as research only. No MEAS1
execution, no CL/IR/T disposition, no implementation.

### 15.3 What was and was not executed

Executed: file reads, `grep`/`sed`/`awk` searches, `sha256sum`, `ls`, `git rev-parse/status/log`
(read-only), `go version` (no build), raw-blob HTTPS fetches of upstream sources, Hub/GitHub API
page reads. **Never**: builds, tests, servers, model loads, inference, measurements, register/freeze
edits, commits, deletions, or file modifications outside this artifact. Temp upstream blobs lived at
`/tmp/tools_*`/`/tmp/upstream_server.cpp` (outside the repo, removed at close).

### 15.4 Deliverable inventory

Created: exactly one file — `docs/review/AUDIT-t2.9.5b-full-reverification-muse.md` (this artifact).
Its SHA-256, byte/line counts, `HEAD`, and working-tree state are reported in the completion message
and are part of the audit record. A1, the DeepSeek audit, the review index, all register revisions,
all freeze records, and all source files are byte-identical to their pre-audit state.
