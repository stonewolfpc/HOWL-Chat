# Audit A1 — exhaustive re-verification of Correction Addendum A1

**Audited artifact:** `docs/review/t2.9.5b-audio-research-correction-addendum-a1.md`
**SHA-256 at audit time:** `2eb69c778b6a2737a4d5415e5ab9645b0f4566c365edb75e1e2e4fb9dbfa2bbe`
**Audit date:** 2026-11-04
**Audit type:** read-only, claim-by-claim, against primary sources.

A1 was **not modified** by this audit. Findings are recorded here as errata; the artifact stands as
captured so that reviewers can see both the claim and its correction.

## Method

Every citation in A1 was re-derived from source, not from memory:

1. **Repository claims** — line-exact reads and `grep -n` against this working tree at
   `HEAD 4f5defedb1df037052dd75ada9c5d5dea21f20de`.
2. **Upstream claims** — files fetched at the exact commit `a34de8618b25dfc32ee43749db47d767a9f0a25d`
   (`raw.githubusercontent.com/ggml-org/llama.cpp/<sha>/…`) and read with `sed -n` around the cited
   lines. No llama.cpp source is vendored in this repository, so this is the only route to the
   cited lines; it verifies the *commit's* text, not that the shipped DLL was built from it.
3. **Model claims** — Hugging Face API (`/api/models/<id>`, `/tree/main`) for repository existence,
   file lists, and license metadata.
4. **Version claims** — CMake build-number plumbing, `git ls-remote` tag dereference, v0.9.0 release
   notes.

## Verdict

| | Count |
|---|---|
| Claims confirmed exactly | 31 |
| Confirmed with a citation-range imprecision | 5 |
| **Wording that overstates or misstates source behaviour** | **3** |
| Claims that remain unverified | 3 |
| Uncertainties A1 declared that are now resolved | 5 |

**No central conclusion is overturned.** The registered routes, their handlers, their capability
gates, the absence of any app-side audio call path, the model families, and the version gap all hold.
Three statements require re-wording before the artifact should be treated as review-clean, and five
citations need tightening. There are **no critical findings** and no finding that changes the
recommended Tranche A sequence.

---

## 1. Findings requiring action

### F1 — MAJOR: the `/v1/realtime` branch condition is misstated

**A1 §1.1:** "conversational session only if `mtmd_cap.audio_stream && mtmd_cap.chat` and
`?intent=transcription` is absent".

**Source** (`server-context.cpp`, `ws_realtime` handler beginning at 6115):

```cpp
const bool live_asr = !ctx_server.asr_prompt.empty() && req.get_param("intent") == "transcription";
if (ctx_server.mtmd_cap.audio_stream && ctx_server.mtmd_cap.chat && !live_asr) {
    server_realtime_conversation(...); return;
}
if (!meta->has_audio_stream) { /* error: does not support realtime transcription */ }
```

**Why it matters.** `intent=transcription` alone does **not** select transcription. Transcription is
forced only when the model also exposes a non-empty ASR prompt. A chat+stream model with
`intent=transcription` but no ASR prompt still enters the conversational branch; and a non-chat
streaming model transcribes regardless of `intent`. A1's phrasing implies the query parameter is the
switch.

**Correction.** "Conversational branch when `audio_stream && chat` **and not**
(`asr_prompt` non-empty **and** `intent=transcription`); otherwise transcription, which additionally
requires `has_audio_stream`."

### F2 — MAJOR: realtime output modality is stated more strongly than the source

**A1 §4 matrix:** realtime conversational audio, output = "audio stream + text".

**Source** (`server-realtime.h:33-34, 83-85`): "a realtime conversation over websocket … responses
follow `output_modalities` or the model's own switching"; `mtmd_helper::gen_audio gen; // set if the
model generates audio`.

**Why it matters.** A1 correctly warns elsewhere that audio input is not speech-to-speech, then
encodes exactly that implication in its own matrix. A model with `audio_stream && chat` but no
`audio_gen` produces text only.

**Correction.** Output = "text, and/or audio per `output_modalities` / model switching; audio out
requires `audio_gen`". This makes the row consistent with A1's own §4 negatives.

### F3 — MAJOR: `llama-tts` gating is cited to the wrong lines

**A1 §2:** "requires `--mmproj`; checks `mtmd_gen_audio_get_info(...) != NONE`
(`tools/tts/tts.cpp:106-130`)".

**Actual locations in the same file:**

| Claim | Line |
|---|---|
| NOT production-ready | 15 |
| `if (params.mmproj.path.empty())` → "no mmproj provided, use --mmproj" | 61-62 |
| `mtmd_gen_audio_get_info(mctx.get()).type == MTMD_GEN_AUDIO_TYPE_NONE` → "mmproj does not support audio generation" | 96-97 |
| speaker reference / TTS prompt / `gen.init` | 106-130 |

**Why it matters.** The cited range contains none of the gating logic. The claims are true; the
reference is wrong — precisely the defect a review artifact cannot afford.

**Correction.** Cite 61-62, 96-97, and 106-130 separately.

### F4 — MAJOR: licensing table groups families under the wrong terms

**A1 §8 table** places Qwen3-ASR, Qwen3-TTS, Qwen2.5-Omni and Qwen3-Omni together as "per-repo
terms; app registry records `"qwen"`", and marks Voxtral/Ultravox "unestablished".

**Verified metadata:**

| Model | HF license | Note |
|---|---|---|
| `Qwen/Qwen3-ASR-1.7B`, `Qwen/Qwen3-ASR-0.6B` | `apache-2.0` | not custom |
| `Qwen/Qwen3-TTS-12Hz-1.7B-Base` | `apache-2.0` | not custom |
| `Qwen/Qwen2.5-Omni-7B` | `other` | custom Qwen terms |
| `Qwen/Qwen3-Omni-30B-A3B-Instruct` | `other` | custom Qwen terms |
| `mistralai/Voxtral-Mini-3B-2507` | `apache-2.0` | A1 said unestablished |
| `mistralai/Voxtral-Mini-4B-Realtime-2602` | `apache-2.0` | A1 said unestablished |
| `google/gemma-4-E2B-it` | `apache-2.0`, `license_link` → `ai.google.dev/gemma/docs/gemma_4_license` | metadata says Apache-2.0; link points at a Gemma-specific licence page |
| `LiquidAI/LFM3-Audio-2B-GGUF` | `other`, `license_name: lfm1.0`, link `LICENSE` | confirms LFM 1.0 |
| `ggml-org/*-GGUF` quantised repos | no license field returned | derivative terms unestablished |

**Why it matters.** The grouping changes the obligation analysis: three of the four Qwen families A1
lumped as custom are Apache-2.0, and Voxtral is no longer unknown. Gemma 4's metadata/link mismatch is
a legal reading question, not an SPDX lookup — A1's hedge was right, but for a sharper reason.

### F5 — MINOR: `llama-mtmd-cli` realtime rule cited four lines late

A1: `mtmd-cli.cpp:576-579`. The rule is decided at 572-575 (`transcribes = audio_stream && !chat`;
`if (params.image.size() != 1)`), with the message at 576 and the call at 579. Cite **572-579**.

### F6 — MINOR: "streaming true for Whisper and Qwen" overgeneralises

A1 §7 cites `asr_models.go:50,79,137,166` as streaming-true. Verified: 50 (Tiny) true, 79 (Base) true,
108 (**Small**) false, 137 (Qwen 0.6B) true, 166 (Qwen 1.7B) true. Say "Whisper Tiny/Base and both
Qwen entries"; Small is already correct and should not be included in the defect.

### F7 — MINOR: `app-model.go` range

A1 cites 117-126; actual mmproj block is 118-125 (`mmprojPath :=` at 118, `HasVision` at 120/124).

### F8 — MINOR: `app-audio.go` range

A1 cites 107-128 for hardcoded flags. Function starts at 105; the nil-branch literals are 110-111 and
the normal-branch literals 127-128. Cite 105-135 or 110-111/127-128.

### F9 — UNVERIFIED: "95 build-number increments" / build `b12110`

`b12015` is confirmed by the project's own runtime fingerprint. `12110` is **not** independently
observed here: `CMakeLists.txt` at both refs only forwards `${BUILD_NUMBER}` (lines 155-156); the
numeric value is injected by CI. The subtraction is correct *if* both inputs are. Mark the v0.9.0
build number as reported-not-verified.

### F10 — UNVERIFIED: per-family `mmproj` quant variants

A1's E1 rows assert specific projector quantisations for Qwen2.5-Omni (`Q8_0`/`f16`), Qwen3-Omni
(`Q8_0`/`bf16`), Gemma 4 (`BF16`/`Q8_0` + MTP) and Qwen3-TTS (`Q8_0`/`bf16`). Only the Qwen3-ASR and
LFM3 trees were enumerated in this audit (both confirmed exactly). The other four trees were not
re-fetched; those cells remain E1-unreverified.

### F11 — UNVERIFIED: `mtmd-cli` "not production-ready" label

`tools/tts/tts.cpp:15` states it for `llama-tts`. The equivalent label for `llama-mtmd-cli` was not
located in this pass. Restrict the claim to `llama-tts` until the mtmd-cli statement is cited.

---

## 2. Confirmed exactly

### Upstream source (commit `a34de8618`)

| A1 claim | Verified |
|---|---|
| `server-context.h:124-135` declares the three handlers | 124 `post_transcriptions_oai`, 125 `ws_realtime`, 135 `post_tts` |
| Route registration in `server.cpp` | **186** `/v1/audio/transcriptions`, **187** `/audio/transcriptions`, **188** `/tts`, **189** WS `/v1/realtime` |
| `post_transcriptions_oai` at 6080 | exact |
| `ws_realtime` at 6115 | exact |
| `post_tts` at 6435 | exact |
| transcriptions gate `has_mtmd && allow_audio` | `if (!meta->has_mtmd || !meta->chat_params.allow_audio)` → "does not support audio input" |
| multipart `file` mandatory | `server-chat.cpp:681-686`, throws "No input file found for transcription" |
| `response_format` json only | `server-chat.cpp:691-694`, throws otherwise |
| conversion at `server-chat.cpp:673` | exact |
| ASR prompt preset at `server-chat.cpp:695` | exact (`common_chat_get_asr_prompt`) |
| `input_audio` parsing at `server-common.cpp:1290` | exact |
| `run_audio_stream` at `mtmd-cli.cpp:262` | exact |
| `/tts` gate is `mtmd_cap.audio_gen` | `has_cap_tts()` at `server-context.cpp:1034-1036`, `return mtmd_cap.audio_gen` |
| `llama-server` audio features capability-gated per model | consistent with all gates above |

### Repository

| A1 claim | Verified |
|---|---|
| `pipeline.go:438-446` builtin registrations | 438 `registerBuiltins`, 439-446 Whisper×3, Qwen×2, Piper/MeloTTS/Kokoro |
| no multimodal adapter | `pipeline.go:393-401` raises "no MultimodalBridge adapters are registered" |
| `model_resolver.go:69-76` only tiny/base/small | 72-76 exactly three entries |
| `recognizer.go:196` streaming unsupported | exact |
| Qwen is remote DashScope | `client.go:29` endpoint, `:33` `qwen3-asr-flash-filetrans`, `:96` no streaming |
| `app-prompt.go:38` vision-only gate | exact |
| `audio-settings.html:96` / `:153` | exact option lists |
| `STAGE_STATUS.md:194` build 12015 / `a34de8618` | exact |
| fingerprints in `harness-t13.log:131`, `harness-think3.log:110`, `harness-events.log:121` | all three contain `b12015-a34de8618` |
| register r14 lines 8, 9, 19, 20, 30 | CL1, CL2, MEAS1, T6, IR1 respectively |
| only `LICENSE-LLVM-OpenMP` in `runtime/llama` | confirmed by listing |
| `ggml-cuda.dll` 594,123,776 B @ 2026-10-30 15:48:38 | exact |
| `LICENSE-LLVM-OpenMP` 19,741 B | exact |
| `llama-server.exe` 9,728 B launcher | exact |
| no app-side audio route calls | `internal/backend` zero matches; `frontend` only the words "candidate transcriptions" in a tooltip and tutorial prose |

### Models and version

| A1 claim | Verified |
|---|---|
| Qwen3-ASR 1.7B main + `Q8_0`/`bf16` + both mmproj | repo tree lists both GGUFs and both `mmproj-…` files |
| Qwen3-ASR 0.6B unverified | `ggml-org/Qwen3-ASR-0.6B-GGUF` exists (tree not enumerated) |
| Voxtral Mini 3B | `ggml-org/Voxtral-Mini-3B-2507-GGUF` |
| Voxtral Realtime naming unverified | `ggml-org/Voxtral-Mini-4B-Realtime-2602-GGUF` — **4B**, not 3B |
| Ultravox 0.5 1B | `ggml-org/ultravox-v0_5-llama-3_2-1b-GGUF` |
| Qwen2.5-Omni 7B / Qwen3-Omni 30B-A3B | GGUF repos for both (3B & 7B; Instruct & Thinking) |
| Gemma 4 E2B exists and is audio-capable | `ggml-org/gemma-4-E2B-it-GGUF`; `docs/multimodal.md` lists it under audio with "audio input, vision input" |
| LFM3 has a separate GGUF repo with per-quant mmproj | 5 GGUFs + 5 matching `mmproj-` files + `LICENSE` |
| Qwen3-TTS 12Hz 1.7B | `ggml-org/Qwen3-TTS-12Hz-1.7B-Base-GGUF` |
| Pocket TTS / Soprano documented at commit | TTS README: `kyutai/pocket-tts` (per-language dirs), `ekwek/Soprano-80M`, `Soprano-1.1-80M` |
| LFM3-Audio in TTS docs | TTS README §"LFM3-Audio" |
| bundled docs cover these families | `docs/multimodal.md:98-140` |
| v0.9.0 commit `cd959bae…` | annotated tag `c4261a74…` dereferences to `cd959bae99ff4bed22bef2604f45241c09bbb8c8` |
| v0.9.0 adds Chatterbox, Nemotron 3.5, Pocket predefined voices | release notes contain CHATTERBOX, "Nemotron 3.5 Transcribe", "Pocket TTS: support predefined voice…" |
| llama.cpp MIT / whisper.cpp MIT | both LICENSE files: "MIT License, Copyright (c) 2023-2026 The ggml authors" |

---

## 3. Uncertainties A1 declared that this audit resolves

1. **Whisper redistribution (A1 §8, §11.8).** No whisper binary is bundled (`**/whisper*` → 0 files).
   The repo contains acquisition code (`downloader.go`, `downloader_acquire.go`) while
   `recognizer.go:226` states `findWhisperCLI` "never downloads or builds", and `STAGE_STATUS.md:1645`
   records that initialisation no longer downloads. Net: the app does not redistribute whisper.cpp;
   whether any caller invokes the downloader is still unestablished.
2. **Route registration lines (A1 §11.3).** Now exact: `server.cpp:186-189`.
3. **Qwen3-ASR 0.6B (A1 §11.4).** Repository exists; file list still unenumerated.
4. **Voxtral Realtime naming (A1 §11.5).** `Voxtral-Mini-4B-Realtime-2602`.
5. **Pocket TTS / Soprano (A1 §11.6).** Identified above.

---

## 4. Effect on A1's conclusions

- **§1 route/handler establishment:** unchanged and now stronger (exact registration lines).
- **§3 model compatibility:** unchanged; three naming gaps closed, one size correction (Realtime = 4B).
- **§4 matrix:** unchanged rows; realtime output wording needs F2.
- **§5 upgrade position:** unchanged; build-number input F9 flagged.
- **§6 Tranche A:** unchanged. No finding touches the contract-first sequence or the authorization
  boundary.
- **§7 selection/resource rules:** unchanged; F6 narrows one supporting citation.
- **§8 licensing:** substance unchanged (notices still missing; weights still separate), but the
  per-family table must be re-cut per F4.

## 5. Limits of this audit

- Upstream verification reads the commit's files over the network; it does not prove the shipped DLL
  was built from that commit. A1's E3-I grade remains correct and uncorrected.
- No model was loaded, no endpoint called, no server started. Every capability remains below E4.
- `ggml-cuda.dll` compatibility (A1 §1.2) is untouched by this audit and still open.
- Four model trees (Qwen2.5-Omni, Qwen3-Omni, Gemma 4, Qwen3-TTS) were not enumerated; F10 applies.

## 6. Integrity statement

- `docs/review/t2.9.5b-audio-research-correction-addendum-a1.md` was **not edited**: hash unchanged at
  `2eb69c778b6a2737a4d5415e5ab9645b0f4566c365edb75e1e2e4fb9dbfa2bbe`.
- No application source, audio architecture file, register revision, freeze record, or authorization
  state was modified. No build, server, model load, measurement, or Git history change occurred.
- This audit adds one file. `HEAD` remains `4f5defedb1df037052dd75ada9c5d5dea21f20de`.
