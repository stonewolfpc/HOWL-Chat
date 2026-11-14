# T2.9.5B Audio — RECON-2: consolidated DeepSeek ↔ Muse reconciliation

**Type:** research and documentation only. No source, register, freeze, dependency, toolchain, runtime,
model, measurement, commit, or Git-history change was made by this work.
**Date:** 2026-11-04. **HEAD:** `4f5defedb1df037052dd75ada9c5d5dea21f20de` (before and after).
**Upstream pin:** `a34de8618b25dfc32ee43749db47d767a9f0a25d` (raw blobs fetched live this pass).
**Revision:** r2.1 — integrity and evidence-grade correction pass. Section-level changes are logged
in §16; no finding in §3 changed substance, and the three document-integrity defects are corrected
rather than papered over.

**Authority:** this artifact is the authoritative reconciliation. It **supersedes RECON-1**
(`docs/review/RECON-t2.9.5b-deepseek-muse-audio.md`), which is preserved byte-identical at
`d40b34b845a35a32fcf24a4ff42b61e14c661f23fb2e4f0fca0e6410a921b552` and remains in the chain as
evidence. It contains both auditors' positions explicitly, marked **M** (Muse) and **D** (DeepSeek),
with source adjudication underneath. Disagreements are settled by the fetched source, never by
seniority, detail, or vote.

---

## 1. Executive reconciliation

**Muse's independent re-verification is largely correct, and one of my findings was flatly wrong.**
The `llama-mtmd-cli` "NOT a production-ready binary" label exists at `tools/mtmd/mtmd-cli.cpp:34-38`
(operative line 35). My audit's F11 said the label was not located and advised restricting the claim
to `llama-tts`; that was a coverage failure, not a judgement call, and it is **withdrawn**. A1 §2
stands as written.

**RECON-1 is substantially sound but contains three defects of its own**, which this artifact
corrects:

1. **F1 line numbers are wrong.** RECON-1 reports the decisive realtime lines as 6130-6135 with the
   boolean at 6133 and the fallback at 6138. The blob shows comments at 6132-6133, boolean at 6134,
   branch at 6135, fallback gate at 6140. Muse's 6132-6136 is the accurate range; RECON-1's
   "verified this pass" numbering is off by one throughout and off by two at the fallback.
2. **The "tree drift" explanation is false.** RECON-1 claims Muse's status-shape citations were
   stale because `app-audio.go` changed after Muse's audit. `app-audio.go` was last written
   **2026-11-02 20:39**, before either audit; both audits read the same file, and both cited the same
   two still-hardcoded flags correctly. RECON-1's replacement lines (110-113 / 128-131) are off by
   two to three, and its claim that either audit said "all four flags are hardcoded" is a wording
   neither of us used.
3. **D2's "correction" targets a claim Muse did not make.** Muse said `LoadModel` and
   `LoadModelWithMMProj` have zero callers — true. RECON-1 answers that `model.Manager` *is*
   constructed and its `IsModelLoaded` *is* consulted — also true, but not in conflict. The useful
   synthesis is: the manager is live for status queries; its `LoadModel` chain is never invoked.

That is three RECON-1 defects. With my withdrawn F11 the auditor-level defect count is four; these
are accounted for separately from the 22 claim rows, because an error in how an audit argued is not
the same thing as a verdict on a claim (§1 tally below).

**Everything else converges.** Licences, Go/toolchain, the `llama-connect.exe` and second
`ggml-cuda.dll` findings, the port-8080 lifecycle, the stale audio-architecture doc, the
single-image/text-only-history ceiling, the dead router-mode surface, and the mtmd-cli rule all
re-verify. No finding changes A1's conclusions or the authorized position (nothing is authorized).

### Headline tally — 22 unique claim rows

Dispositions are mutually exclusive: every row in §3 carries exactly one primary class, and the
counts sum to 22. Secondary notes (for example a corrected citation whose substance still holds) are
recorded inside the row and are **not** counted a second time.

| Class | Count | Rows |
|---|---|---|
| Confirmed — substance and citation hold | 11 | 8, 12, 13, 14, 15, 16, 17, 19, 20, 21, 22 |
| Corrected — A1 needs re-wording or a citation fix | 7 | 1–7 |
| Unestablished | 2 | 9 (`b12110`), 10 (projector quants) |
| Withdrawn — the finding was wrong | 1 | 11 (my F11) |
| Downgraded — filed as debt, is not debt | 1 | 18 (`MMPProjPath`, now a non-debt observation, §4.10) |

Total: 22.

### Auditor-level defects — 4, counted separately

These are errors in the audits and in RECON-1's argument, not verdicts on code claims, so folding them
into the 22 would double-count. The D2 framing defect lives here rather than in the claim tally for
exactly that reason: it changes no row's verdict, only RECON-1's description of what it was correcting.

| Defect | Artifact | Disposition |
|---|---|---|
| F11 "production label not located" | my audit (D) | **Withdrawn** — label exists at `mtmd-cli.cpp:35` (§4.1) |
| `/v1/realtime` line numbers | RECON-1 | **Corrected** — actual 6134 / 6135 / 6140 (§4.2) |
| "Tree drift" status explanation | RECON-1 | **Corrected** — false; `app-audio.go` mtime predates all audits (§4.3) |
| D2 framed as correcting Muse | RECON-1 | **Corrected** — both observations true; no contradiction to correct (§4.4) |

### Original research report

Still **NOT RECOVERED**. No claim below depends on it. Every "original DeepSeek finding" that
survives is the version recorded in A1 §10 or reconstructed from A1 itself; where a conclusion rests
only on A1, this artifact says so. Nothing here reconstructs or impersonates the lost report.

---

## 2. Baseline and artifact hashes

| Artifact | SHA-256 (recomputed) | Result |
|---|---|---|
| A1 `docs/review/t2.9.5b-audio-research-correction-addendum-a1.md` | `2eb69c778b6a2737a4d5415e5ab9645b0f4566c365edb75e1e2e4fb9dbfa2bbe` | MATCH |
| My audit `docs/review/AUDIT-t2.9.5b-audio-research-a1.md` | `502e98e36eccb2e08be929d172e80b22e8c1ccd0c5751e2c44ed381ce86cc8cb` | MATCH |
| Muse `docs/review/AUDIT-t2.9.5b-full-reverification-muse.md` | `7de11f73ef42d4e825998c5f8d5ed7c0257ecf32d4ef6967c437cf03fe9690e2` | MATCH |
| RECON-1 `docs/review/RECON-t2.9.5b-deepseek-muse-audio.md` | `d40b34b845a35a32fcf24a4ff42b61e14c661f23fb2e4f0fca0e6410a921b552` | MATCH |

Working tree: `HEAD 4f5defe`; 24 pre-existing tracked modifications (`app-audio.go` and
`docs/audio-architecture/README.md` among them, both **modified 2026-11-02**, i.e. before all three
audits — the fact that disproves RECON-1's drift claim). New files are confined to `docs/review/`.
No register, freeze, or authorization artifact was touched.

---

## 3. Claim-by-claim comparison

**M** = Muse re-verification. **D** = my audit. **R1** = RECON-1. Verdict cites the fetched source.

| # | Claim | D | M | R1 | Verdict |
|---|---|---|---|---|---|
| 1 | `/v1/realtime` branch condition | MAJOR corrected | confirmed + comment provenance | confirmed, lines 6130-6135 | **Confirmed.** Comments 6132-6133; `live_asr` 6134; branch 6135; fallback 6140. M's range right, R1's off by one |
| 2 | Realtime output modality | MAJOR corrected | confirmed | confirmed | **Confirmed.** `server-realtime.h:33-34,83-85`: follows `output_modalities`; audio only if `gen` set |
| 3 | `llama-tts` gate citations | 61-62 / 96-97 | same | same | **Confirmed.** Label block 14-18 (operative 15), `--mmproj` 61-62, NONE reject 96-97, speaker/prompt/init 106-130 |
| 4 | Licence grouping | MAJOR re-cut | confirmed, two refinements | accepted | **Confirmed.** Table in §10 |
| 5 | `mtmd-cli` realtime rule | 572-579 | 572/574/576/579 | 572/574/575/576/579 | **Confirmed.** Image-count check is 575; R1's set is complete, M's missed one line |
| 6 | "streaming true for Whisper" | MINOR | confirmed | confirmed | **Confirmed.** `asr_models.go` 50/79 true, **108 Small false**, 137/166 true |
| 7 | `app-model.go` mmproj block | 118-125 | confirmed | confirmed | **Confirmed.** 118 `mmprojPath :=`, 120/124 `HasVision` |
| 8 | `app-audio.go` hardcoded flags | 110-111 / 127-128 | same | "all four", 110-113 / 128-131, drift | **D and M confirmed; R1 corrected.** Nil branch: asr 108, tts 109, multimodal 110, realtime 111. Live branch: asr 125, tts 126, multimodal 127, realtime 128 |
| 9 | `b12110` / "95 increments" | unverified | unestablished | unestablished | **Unestablished.** `b12015` observed in logs; `12110` never observed; subtraction conditional |
| 10 | Four families' mmproj quants | unverified | unestablished | unestablished | **Unestablished.** Trees not enumerated; licence fields now enumerated (§10) |
| 11 | `mtmd-cli` production label | "not located" | **refuted**, it exists | confirmed | **D wrong, M right.** Label at 34-38, line 35. F11 withdrawn |
| 12 | Router-only `/models/load` | not examined | D2, zero callers | "correction": manager live | **Confirmed with correction.** `Manager.LoadModel` zero callers; manager itself live for `IsModelLoaded` (chat/service.go:449,526) |
| 13 | Single image + text-only history | not examined | D3 | confirmed | **Confirmed.** `app-prompt.go:25`; history 620-747; `HasImage bool` persists at chat/types.go:78 |
| 14 | `llama-connect.exe` Sep-2 outlier | not examined | D4 | confirmed, hygiene | **Confirmed.** 11,473,920 B @ 2026-09-02 17:14 |
| 15 | Go 1.23.8 vs 1.27.3 | not examined | D5, premise stale | confirmed, headline narrowed | **Confirmed.** `go.mod:3`; `go1.27.3`; no CI workflows |
| 16 | Port-8080 kill/respawn | not examined | D6 | confirmed, design not residue | **Confirmed.** Check 78, hard fail 79-81; crash re-check 254-256 |
| 17 | Stale audio-architecture doc | not examined | D7 | confirmed + sizes wrinkle | **Confirmed and stronger.** See §4.7 |
| 18 | `MMPProjPath` double-P | not examined | — | "possible typo" | **Downgraded, not debt.** Consistent at 3 sites, no json tag → cosmetic; removed from the debt table (§4.10) |
| 19 | Zero app-side audio calls | confirmed | confirmed + mmproj qualification | folded into "T1" (now DEBT-T1) | **Confirmed.** No audio-endpoint construction in `internal/backend` |
| 20 | Whisper redistribution | resolved: not bundled | extended: dormant code | same | **Confirmed.** No `**/whisper*` binary; `findWhisperCLI` never downloads; `Acquire` uncalled |
| 21 | Routes, gates, formats | confirmed | confirmed | confirmed | **Confirmed.** `server.cpp:186-189`; gate 6082-6084; file 681-685; json 691-694; `input_audio` 1290 |
| 22 | Version gap / v0.9.0 content | confirmed | confirmed | confirmed | **Confirmed.** Tag `c4261a74…` → commit `cd959bae…`; notes contain Chatterbox, Nemotron 3.5, Pocket voices, ggml v0.29.0, miniaudio 0.11.26, cpp-httplib 0.64.0 |

---

## 4. Resolved disagreements

### 4.1 F11 — my error, withdrawn

`tools/mtmd/mtmd-cli.cpp:34-38`:

```cpp
/**
 * Please note that this is NOT a production-ready binary.
 * It is a playground for trying multimodal support in llama.cpp.
 * For contributors: please keep this code simple and easy to understand. ...
 */
```

A1 §2's claim was correct. **Root cause of my miss:** in the audit pass I fetched `tools/tts/tts.cpp`
but never `tools/mtmd/mtmd-cli.cpp`, then reported the label as unlocated. Mitigation recorded for
future audits: a "not located" verdict must name the files searched. No user-visible impact; the
label governs upstream CLI use, not the app.

### 4.2 F1 citations — Muse right, RECON-1 wrong

Actual blob mapping at `server-context.cpp`:

| Line | Content |
|---|---|
| 6132-6133 | two comments, the second using the `?intent=transcription` shorthand |
| 6134 | `const bool live_asr = !ctx_server.asr_prompt.empty() && req.get_param("intent") == "transcription";` |
| 6135 | `if (ctx_server.mtmd_cap.audio_stream && ctx_server.mtmd_cap.chat && !live_asr) {` |
| 6140 | `if (!meta->has_audio_stream) {` → "does not support realtime transcription" |

Muse's 6132-6136 spans comment → branch → call and is accurate. RECON-1's 6130-6135 puts the comment
at 6131-6132 and the boolean at 6133 — both off by one — and puts the fallback at 6138, which is the
closing brace. The substance of RECON-1's F1 conclusion is right; its numbers are not.

### 4.3 Status flags — no drift, and both prior audits were right

`app-audio.go` (8,223 B, mtime **2026-11-02 20:39**) predates all three audits. Current shape:

- nil branch 106-114: `asr_ready` 108, `tts_ready` 109, `multimodal_ready` 110, `realtime_ready` 111,
  `active` 112 — all literal false.
- live branch 124-132: `asr_ready` 125 and `tts_ready` 126 from `status`; `multimodal_ready` 127 and
  `realtime_ready` 128 literal false.

Muse cited 110-111 / 127-128 and said `realtime_ready` is hardcoded: correct. I cited the same pair
for the two hardcoded capability flags: correct. RECON-1's claim that Muse "predates the live-status
change" is unsupported by mtime, and its "all four hardcoded" quotation is not in either audit.
**The accurate A1 errata item is narrow:** multimodal and realtime remain hardcoded false; ASR/TTS
readiness is live.

### 4.4 D2 — both observations true, RECON-1's framing wrong

- `model.NewManager(client)` at `app.go:102`; stored on `App.modelManager` at 179; passed to
  `chat.NewService` at 116.
- `chat/service.go:449,526` call `s.modelManager.IsModelLoaded()` — the manager is production-live.
- `Manager.LoadModel` (`manager.go:40`) delegates to `llamaClient.LoadModel` at 75, which POSTs
  `/models/load` (`http_client.go:276-313`) — an endpoint registered only inside the router block
  (`server.cpp:525-533`, `/models/load` at 530). The app always spawns single-model servers.
- `LoadModelWithMMProj` (`http_client.go:393-397`) ignores its `mmprojPath` argument.
- **Zero callers** of `Manager.LoadModel` or `LoadModelWithMMProj`.

Verdict: Muse's "dead, misleading code" is right; RECON-1's addition (manager live for status) is
right; RECON-1's description of it as correcting Muse is wrong. Latent severity only — nothing
invokes the chain today.

### 4.5 `mtmd-cli` rule — RECON-1 has the complete set

572 decision, 574 branch, **575** image-count check, 576 message, 577 return, 579 dispatch. RECON-1
is the only one of the three that records 575. Muse's set omitted it; my span covered it.

### 4.6 TTS label block — RECON-1 right, Muse off by one

Block opens at 14, operative sentence 15, closes 18. Muse wrote 13-17; RECON-1 wrote 14-18. A1's
substantive claim never depended on the block boundary.

### 4.7 Audio-architecture doc — stale in a way none of the three fully recorded

`docs/audio-architecture/README.md` (13,621 B, mtime 2026-11-02 10:00; modified vs HEAD):

- Stages 2-5 marked **PENDING** at 64, 69, 75, 78 while ASR, TTS, transcription UI and single-image
  vision ship.
- TTS matrix lists **Bark** at 207; no Bark loader exists (`pipeline.go:438-446` registers
  Piper/MeloTTS/Kokoro).
- FAQ at 387 advises "upgrade to MeloTTS or Bark".
- ASR sizes 39/74/244 MB at 195-197 are not the `ggml-*.bin` artifacts the resolver approves
  (`model_resolver.go:72-76`).
- **New here:** the doc's file tree names paths that do not exist —
  `internal/audio/asr/whisper/local.go`, `internal/audio/asr/qwen-asr/client.go`,
  `internal/audio/multimodal/`, `internal/audio/realtime/` are all absent. The real tree is
  `asr/whisper/recognizer.go`, `asr/qwen/client.go`, `tts/{piper,melotts,kokoro}`, `runtime/`, `hw/`.
  The document is a plan drawing, not a description of the tree.

### 4.8 Licensing — three-way agreement, with new primary rows

Verified this pass (HF API): `ggml-org` **convert repos with no licence field at all** —
Qwen3-ASR-1.7B-GGUF, Qwen3-ASR-0.6B-GGUF, Qwen3-Omni-30B-A3B-Instruct-GGUF,
Qwen3-TTS-12Hz-1.7B-Base-GGUF, Voxtral-Mini-3B-2507-GGUF. `ggml-org/ultravox-v0_5-llama-3_2-1b-GGUF`
→ `mit`. `ggml-org/gemma-4-E2B-it-GGUF` → `apache-2.0`. `ekwek/Soprano-80M` → `apache-2.0`.
`kyutai/pocket-tts` → `cc-by-4.0`, `gated: auto`. Qwen3-Omni **Thinking** repo was not checked;
only Instruct.

### 4.9 Go and CI — gap closed

`go.mod:3` declares `go 1.23.8`; installed toolchain `go1.27.3`. `.github/` contains only
`copilot-instructions.md` — **no `workflows/` directory**, so there is no CI toolchain pin to
reconcile. RECON-1's substantive point (toolchain moved, module contract did not) is correct;
Muse's "not enumerated" gap is now resolved rather than merely repeated.

### 4.10 `MMPProjPath` — non-debt observation, removed from the register

`MMPProjPath` is spelled that way in its declaration (`types/model_profile.go:72-73`) and at both use
sites (`app-model.go:121`, `app-prompt.go:40`). There is no competing `MMProjPath`; the field carries
no JSON tag, so there is no serialization mismatch, no runtime effect, and no failure mode. It is a
naming wart only.

**Disposition:** recorded here as a non-debt observation and **removed from the technical-debt table**
(§11). No rename is recommended, not even as an optional item: the "no knowingly carried-forward
debt" rule cuts against inventing cosmetic work, and a cosmetic row retained purely to preserve
numbering is precisely that. If a future authorized change touches the struct for another reason, the
spelling may ride along at zero marginal cost — an observation, not a task.

### 4.11 Editorial defect in RECON-1

Its §7.1 contains a garbled clause ("same family, un acknowled redistribution/package gap") and its
§2 describes `.github/workflows/` as "empty" when the directory does not exist at all. Neither
affects a finding, but an authoritative artifact should not carry text no one can parse.

---

## 5. Unresolved and evidence gaps

1. **Original research report** — unrecovered; conclusions rest on A1 §10 alone.
2. **`b12110`** — never observed; every statement about v0.9.0's build number is reported, not measured.
3. **Four projector-quant cells** (Qwen2.5-Omni, Qwen3-Omni, Gemma 4, Qwen3-TTS) — unenumerated.
4. **`ggml-cuda.dll` / `llama-connect.exe` provenance** — no receipt, SBOM, or source hash; ABI
   compatibility with the server binary remains E3-I.
5. **Image bytes in saved history** — `chat/types.go:78` persists `HasImage bool` only; whether any
   other store keeps the image (and in what form) was not traced. Inference context definitively
   excludes it.
6. **Pre-migration capability set** — no Gollama-era code remains, so "vision regressed vs was never
   multi-image" is unrecoverable.
7. **Upstream newer than v0.9.0** — not enumerated; no claim made beyond that tag.
8. **CVE/advisory status** for the pinned runtime, whisper.cpp, Go and Wails — not assessed.

---

## 6. Reconciled audio capability matrix

No E4: nothing was loaded, run, or measured.

| Capability | Grade | Limiting layer |
|---|---|---|
| Local file ASR (Whisper) | Partial | `Transcribe` path exists; resolver approves tiny/base/small only (`model_resolver.go:72-76`) while UI offers medium/large (`audio-settings.html:96`) |
| Live/streaming ASR | UI-exposed-unsupported | Both adapters reject (`recognizer.go:196-202`, `qwen/client.go:95-101`); registry declares true for Tiny/Base/Qwen (`asr_models.go:50,79,137,166`) |
| Direct audio understanding | Not implemented | Zero `input_audio` construction in `internal/backend` |
| Realtime conversation | Not implemented | No WS client; `realtime_ready` hardcoded false at `app-audio.go:111,128` |
| Native TTS (`audio_gen`) | Not implemented | No `/tts` caller; Piper/Kokoro/MeloTTS are independent of llama.cpp |
| Status honesty | Partial | `asr_ready`/`tts_ready` live (125-126); multimodal/realtime hardcoded false (127-128) |
| Capability detection | Partial + misleading | `HasVision` from mmproj presence (`app-model.go:118-125`); no audio capability field |
| Vision — historical capability | Owner-reported working (pre-audit) | Testimony only; not measured in this pass |
| Vision — current implementation, single image | Implementation present; end-to-end **not reverified** | One `imageData` (`app-prompt.go:25`); one `image_url` part (`http_client.go:582-583`, `712-713`) |
| Vision — multi-image | App-layer limit established; per-model behavior **unestablished** | Upstream loops over content parts, but no model was exercised with two or more images |
| Vision in history | Not implemented | `buildPromptMessages` 620-747 emits zero image parts |
| Whisper acquisition | Dormant code | Downloader defined, zero callers; finder never downloads |

---

## 7. Application capability and multimodal history

The user-visible record and the inference record diverge, and the gap is exactly located:

1. `SendMessageStreamWithImage(message, imageData string)` accepts **one** base64 image
   (`app-prompt.go:25`). It decodes and forwards one `[]byte`.
2. `buildRequestWithImage` and `GenerateChatStreamWithImage` each construct exactly **one**
   `image_url` part and label it `data:image/jpeg;base64,` unconditionally
   (`http_client.go:582-583`, `712-713`) — PNG/WebP bytes are mislabelled. Upstream sniffs media by
   content, so this is latent, not currently fatal.
3. `buildPromptMessages` (620-747) emits system, roleplay context and text history only — **no image
   parts**. From turn 2 onward the model cannot see turn 1's image.
4. Persistence records `HasImage bool` (`chat/types.go:78`): the transcript knows an image existed,
   but the image is not replayed into context.

**Vision status — four separate evidence grades**, so that neither more nor less is claimed than the
pass supports:

1. **Historical working capability:** owner-reported. Vision is understood to have worked before this
audit; that is testimony, not a measurement recorded here.
2. **Current implementation:** the single-image request path exists and is located — one `imageData`,
one `image_url` part, one decoder.
3. **Current end-to-end operation:** **not reverified.** Nothing was built, loaded or run in this
pass, so "works today" is unsupported and is not claimed.
4. **Multi-image:** the app-layer ceiling is established at three functions; each model's own
multi-image behaviour is **unestablished** — upstream accepting N content parts proves the server can
carry them, not that a given projector handles them correctly.

What does follow: the ceiling is app-layer construction, not the server, the model, or the migration —
`app-llama.go:196-198` already passes `--mmproj`. Any fix is additive multi-part construction plus
history carriage, not a rebuild, and nothing here authorizes it. An E4 load with a real image is the
minimum evidence to claim current operation; two images plus one history turn is the minimum to claim
multi-image.

---

## 8. Backend migration (Gollama → llama.cpp)

- **No functional Gollama import, cgo binding, or API shape** in Go code; the backend is `exec` +
  HTTP against `llama-server`.
- **Residual text:** comments at `prompt/builder.go:6`, `types/model_profile.go:97,144`; user-visible
  tutorial prose at `frontend/dist/js/tutorials-data-model.js:54,79` ("auto-managed by gollama",
  "gollama backend may auto-manage this", "use the gollama CLI").
- **Shared root cause:** exactly one — **DEBT-T1, capability inferred from projector filename**. It
  produces the vision gate, the audio-projector blind spot, and the status/UI/resolver disagreement.
  D2 is the only item with a plausible migration-leftover shape (router-mode client retained under a
  single-model design), and it stays **Inferred**: the pre-migration tree is gone.
- **Not migration losses:** single-image construction, JPEG labelling, port-8080 exclusivity and the
  audio-architecture doc are native app-layer design and doc drift. No evidence supports rebuilding a
  working subsystem, and none is recommended.
- **Genuine gains, unchallenged:** `reasoning_content` streaming with a think-splitter
  (`app-prompt.go:57-89`), `--reasoning-budget` (verified against this binary's `--help`,
  `STAGE_STATUS.md:190-193`), `/props`, `/v1/models`, multi-architecture GGUF.

---

## 9. Go and llama.cpp upgrade assessment

**Go.** Declared `go 1.23.8`; toolchain `go1.27.3`; no CI. The directive governs language semantics
and `go vet` stdversion checks regardless of toolchain, so the accurate statement is: toolchain
current, module contract three minor versions behind. No incompatibility was observed — no build was
run. A bump is a staged change with its own build, vet and smoke plan, never bundled with a runtime
move. **No upgrade performed or recommended now.**

**llama.cpp.** Shipped `b12015-a34de8618` (2026-10-30). v0.9.0 (`cd959bae`, 2026-11-04) adds
Chatterbox TTS, Nemotron 3.5 Transcribe, predefined Pocket voices, voice APIs, ggml v0.29.0 (batched
copy superseding `_2d`), miniaudio 0.11.26, cpp-httplib 0.64.0. Additive minor release; it touches
`/props` shape and deprecates a ggml API. Nothing in the Whisper + native-audio-input scope requires
it, and the runtime provenance record is unclean (`ggml-cuda.dll`, `llama-connect.exe`). **No upgrade
performed or recommended now**; adopt only against a named gap, with props-shape and `_2d` call-site
checks.

---

## 10. Licensing and redistribution position

Evidence grades, not legal conclusions.

### 10.1 Runtime notices — directory-level gap confirmed; distribution compliance unestablished

`runtime/llama/` ships ~40 binaries with exactly one notice file, `LICENSE-LLVM-OpenMP` (19,741 B).
No llama.cpp MIT text, no ggml notice, no miniaudio or cpp-httplib notice. Upstream is MIT
("Copyright (c) 2023-2026 The ggml authors"), whose terms require the notice to accompany
redistributed copies.

**What this establishes:** a notice gap **inside that directory** — a real packaging concern, and
consistent with A1 §8's warning about second-order notices.

**What it does not establish:** whether the complete distributed application supplies the same notices
by another route — installer payload, release archive, a top-level `LICENSES/`, or notices generated
at build time. The obligation attaches to distribution, not to one folder read in isolation, and the
distribution layout was **not** enumerated in this pass.

**Disposition:** downgrade from "confirmed packaging defect" to **notice gap observed in
`runtime/llama/`; complete distribution-package compliance unestablished**. Closing it requires
inspecting the actual shipped package before any packaging change is scoped. No legal conclusion is
made, here or elsewhere.

### 10.2 Weights (per-family; SPDX ≠ terms)

| Model | Licence | Note |
|---|---|---|
| Qwen3-ASR 1.7B / 0.6B (upstream) | `apache-2.0` | not custom Qwen terms |
| Qwen3-TTS 12Hz 1.7B Base (upstream) | `apache-2.0` | — |
| Voxtral Mini 3B / 4B Realtime (upstream) | `apache-2.0` | Realtime is 4B |
| Qwen2.5-Omni-7B | `other` | custom terms |
| Qwen3-Omni-30B-A3B-Instruct | `other` | Thinking unchecked |
| Gemma 4 E2B/E4B | `apache-2.0` + `license_link` → `ai.google.dev/gemma/docs/gemma_4_license` | do not rely on the tag alone |
| LFM3-Audio 2B | `other` / `lfm1.0` | licence file in repo |
| Ultravox GGUF | `mit` | base model terms still apply |
| Soprano-80M / 1.1-80M | `apache-2.0` | TTS systems, not vocoders |
| Pocket TTS | `cc-by-4.0`, gated (auto) | per-family review needed |
| ggml-org converts: Qwen3-ASR ×2, Qwen3-Omni-Instruct, Qwen3-TTS, Voxtral-3B | **no licence field** | derivative-terms gap |

Weights are not bundled today. No redistribution determination is made here.

### 10.3 Whisper

No whisper binary is bundled (`**/whisper*` → 0). The downloader exists but is dormant: `Acquire` and
`DownloadLatest` have zero production callers, and `findWhisperCLI` "never downloads or builds"
(`recognizer.go:226-228`). If ever wired to a UI affordance, notice duty follows the downloaded bits.

---

## 11. Consolidated technical-debt register

Any future modularization follows COPY → PASTE → EDIT → VERIFY → DELETE THE OLD; an untested
replacement never displaces working code.

**Namespace rule.** Every technical-debt ID carries the `DEBT-` prefix. The bare token `T6` is **not**
used in this artifact for debt: `T6` belongs to the authorization register (`IR1` / `MEAS1` / `T6` /
`CL3`), which is untouched. This removes the collision where a debt item and an authorization gate
shared one identifier and could be mistaken for each other.

| # | Item | Root cause | Evidence | Impact | Direction → dependency → regression |
|---|---|---|---|---|---|
| DEBT-T1 | Capability inferred from filename | no negotiated-capability record | `app-model.go:118-125`; `http_client.go:318-391` | audio projector would set `HasVision`; UI misreports | `/props` + family knowledge → gates DEBT-T2/DEBT-T3 → audio-mmproj ⇒ vision-false |
| DEBT-T2 | Declared-but-unsupported streaming | registry rows added without adapter check | `asr_models.go:50,79,137,166` vs rejections | user selects streaming, gets error | correct rows (Small is already false) → DEBT-T1 → registry-vs-adapter contract test |
| DEBT-T3 | Status/registry/UI disagreement | three sources of truth | `audio-settings.html:96`; resolver 72-76; `app-audio.go:108-111,125-128` | UI offers models the resolver rejects | single capability source (= A1 Tranche A item 1) → DEBT-T1 → triple-agreement test |
| DEBT-T4 | Dead `/models/load` chain + ignored arg | carried-over client shape (Inferred) | `http_client.go:276-313,393-397`; zero `Manager.LoadModel` callers | none today; 404 if invoked | delete or re-target → owner decision → no-caller proof |
| DEBT-T5 | Runtime provenance, two binaries | no receipt/SBOM; purpose of `llama-connect.exe` unknown | `ggml-cuda.dll` 594,123,776 B @ Oct 30 15:48; `llama-connect.exe` 11,473,920 B @ Sep 2 17:14 | every load claim inherits E3-I | identify purpose + source/hash for both → before runtime work → hash-recorded build |
| DEBT-T6 | Single-image / text-only history | single-part app construction | `app-prompt.go:25,620-747`; `http_client.go:576-583,703-713` | multi-image impossible; turn-2 image invisible | additive multi-part + carriage → design + authorization → cap/carriage tests |
| DEBT-T7 | JPEG MIME hardcode | copy-paste default at both sites | `http_client.go:583,713` | **latent only** — no observed failure; upstream sniffs by content | sniff-or-carry type → with DEBT-T6 → byte-fidelity test |
| DEBT-T8 | Port-8080 exclusivity, kill-respawn | single-model design, no queue | `app-model.go:78-81,254-256` | switch drops chat; occupant hard-fails load | refuse-or-queue policy (A1 §7) → arbitration design → refusal/respawn/no-leak tests |
| DEBT-T9 | **Provisional** — repeated Whisper init can orphan a temp dir | no guard against re-initialization while `r.tempDir` is set | `recognizer.go:82` (`MkdirTemp`), `90` (`r.tempDir = tempDir`, unconditional) vs `209-213` (`Release` removes only the current path); originally surfaced from `STAGE_STATUS.md:1690-1698` | source-level mechanism confirmed; **production reachability unestablished** — repeated successful `Initialize` on one instance was not demonstrated | trace callers and re-init paths first; then decide cleanup-on-overwrite vs guard vs reuse. **No one-shot guard is prescribed** |
| DEBT-T10 | Stale docs + gollama prose | plans frozen, code moved | audio-architecture 64-80,195-207,387 + nonexistent paths; tutorial 54,79 | readers and users misled | status header; fix tutorial strings on touch → owner docs decision → doc-vs-tree check |
| DEBT-T11 | Dormant Whisper downloader | acquisition code retained, unwired | `downloader*.go`; zero callers | none today; distributor duty if wired | keep-or-remove with notice analysis → owner decision → no-caller proof |

`MMPProjPath` (formerly the twelfth row) is **not** in this register; it is a non-debt observation
(§4.10).

Shared roots: DEBT-T1 → DEBT-T2/DEBT-T3; DEBT-T6/DEBT-T7 share the `buildRequest*` surface; DEBT-T4/
DEBT-T5 share the runtime-client boundary. No unavoidable debt identified.

---

## 12. Dependency-ordered remediation

Recommendations only; nothing is authorized.

1. **Errata note for A1** (never editing A1): F1 wording, F2 output cell, F3 split citation, F4 table
   re-cut, F5 exact set 572-579, F6 Tiny/Base + Qwen, F7/F8 ranges, **F11 withdrawn**.
2. **Packaging:** ship llama.cpp/ggml MIT + vendored-lib notices. No code, no upgrade.
3. **Provenance:** identify what `llama-connect.exe` is, whether anything invokes it, and where it came
   from; record source URL + hash for the 594 MB `ggml-cuda.dll`, or replace from a pinned asset.
   **No deletion or disposal is recommended** — the Sep-2 timestamp is a reason to investigate, not a
   reason to remove, and removal would be a shipped-package change requiring owner approval (§14.1).
4. **Tranche A contract-first** per A1 §6. The ordering is a recommendation and authorizes nothing:
   **DEBT-T1 → DEBT-T3 → DEBT-T2**, because the single capability record DEBT-T1 introduces is what
   DEBT-T3 reads and what DEBT-T2's corrected rows must agree with. DEBT-T4, DEBT-T6, DEBT-T7, DEBT-T8 are the
   tranche's ceiling, not its backlog.
5. **Docs:** status header on the audio-architecture doc; correct the two gollama tutorial strings on
   next touch; register r14 remains the sole controlling artifact.
6. **Lifecycle/history:** DEBT-T8 policy and DEBT-T6 carriage — designed, not assumed. DEBT-T7's MIME
   fidelity rides with DEBT-T6 rather than as separate work.
7. **Go directive:** leave it until a staged pass with build + vet + smoke.

---

## 13. Required acceptance and regression coverage

For any future authorized tranche: no-caller proof for `/models/load`; options/mmproj passthrough if
kept; audio-mmproj ⇒ vision-false; registry-vs-adapter streaming agreement; status/UI/resolver triple
agreement; one-image cap + history-carriage + MIME fidelity; port-occupied refusal; respawn restores
model; no leaked server process; repeat-init no orphan; hash-recorded build; doc-tree consistency.
E4 (load + negotiated capability on this host) is the floor for admitting native audio; E5 (correct
output on a user request) for any user-facing claim.

---

## 14. Decisions, prerequisites, and register control

Three separate lists. A recommendation in this artifact is **not** an authorization, and appearing in
§12 does not move an item into any tranche.

### 14.1 Owner approval required (shipped product or policy)

1. The A1 errata note and its exact text, including the withdrawn F11.
2. Any change to shipped runtime packaging (notices) or to the distribution layout — after §14.2 item 3.
3. Retention, replacement, or removal of any binary in `runtime/` — `llama-connect.exe` and the
   duplicate-size `ggml-cuda.dll` included — and only after §14.2 item 2.
4. DEBT-T4 cleanup (delete vs re-target the `/models/load` chain).
5. Wiring the dormant Whisper downloader, with the distributor determination first.
6. Any native ASR admission (A1 Tranche A entry criteria stand; E4 is the floor).
7. DEBT-T8 refuse-or-queue policy and the DEBT-T6 history-carriage design.
8. DEBT-T10 docs precedence, or confirmation that the register remains sole control.
9. A Go directive bump and its test plan.
10. Any llama.cpp upgrade, against a named gap.

### 14.2 Research required before any decision

1. **DEBT-T9:** production reachability of repeated Whisper initialization (callers, re-init paths) and
the correct fix — cleanup-on-overwrite, guard, or reuse.
2. Runtime provenance: purpose, source and hash for `llama-connect.exe` and `ggml-cuda.dll`; ABI
compatibility with the server binary.
3. Distribution-package notice inspection (installer payload or release archive), not just
`runtime/llama/`.
4. Unresolved model/projector compatibility: the four projector-quant cells, and per-model multi-image
behaviour.
5. `b12110` — one log observation, or drop the claim.
6. DEBT-T11 keep-or-remove analysis, if wiring is ever considered.

### 14.3 Controlled exclusively by the register (untouched here)

`IR1`, `MEAS1`, register `T6`, and `CL3` remain register matters. This artifact neither disposes nor
renumbers them. The bare token `T6` in this document always refers to the register item; every
technical-debt ID is prefixed `DEBT-` (§11).

---

## 15. Integrity close

- Inputs verified by hash before reliance (§2), including RECON-1. Revision r2.1 changed **this file
  and the review index only**; the four input artifacts remain byte-identical.
- **Created:** this file. **Modified:** `docs/review/README.md` (index entry and result line only).
  **Unchanged:** A1, my audit, Muse's audit, RECON-1, all register/freeze/authorization artifacts, all
  source files.
- **Executed:** file reads, `grep`/`sed`/`awk`, `sha256sum`, `ls`, `git rev-parse/status` (read-only),
  `go version` (no build), HTTPS fetches of upstream blobs at the pinned commit, Hugging Face API
  reads, GitHub release/tag reads.
- **Never:** builds, tests, servers, model loads, inference, measurements, register/freeze edits,
  commits, deletions, or source modifications.
- Working tree: 24 pre-existing tracked modifications; `HEAD 4f5defe` unchanged.

---

## 16. Revision r2.1 — correction-pass log

Owner review identified eight defects in the first RECON-2 text. Each is corrected below from the
cited source and reasoning, not from the reviewer's authority, and none was resolved by softening a
finding.

| # | Defect | Correction applied |
|---|---|---|
| A | Debt `T6` collided with authorization-register `T6` | All debt IDs renamed `DEBT-T1`…`DEBT-T11`; namespace rule stated in §11; bare `T6` now refers only to the register (§14.3) |
| B | `MMPProjPath` both downgraded and retained as debt | Removed from the debt table; recorded as a non-debt observation (§4.10); no rename recommended, not even optional |
| C | Tally did not reconcile; D2 framing unaccounted | Re-cut as 22 mutually exclusive claim rows (11 / 7 / 2 / 1 / 1) plus a separate four-row auditor-defect table that names the D2 framing defect (§1) |
| D | DEBT-T9 over-graded, fix prescribed | Regraded **Provisional**; mechanism cited at `recognizer.go:82,90,209-213`; production reach stated unestablished; one-shot-guard prescription withdrawn (§11) |
| E | Licensing conclusion too broad | Narrowed to a directory-level gap; complete distribution compliance unestablished; inspection listed as a prerequisite (§10.1, §14.2) |
| F | `llama-connect.exe` disposal recommended | Recommendation removed; identification of purpose, use and origin required first (§12, §14.2) |
| G | Vision status over-claimed | Split into four grades — historical / current implementation / current end-to-end / multi-image (§6, §7) |
| H | Owner-decision accounting incomplete | One list replaced by three: owner approval, research prerequisites, register-controlled (§14) |

Two further items were made explicit rather than left implicit: DEBT-T7's JPEG MIME defect is graded
**latent, no observed failure**, and the `DEBT-T1 → DEBT-T3 → DEBT-T2` ordering is stated as a
recommendation that authorizes nothing (§12.4).

### 16.1 Verification of this pass

- Every §3 claim row carries exactly one class; the five counts sum to 22.
- No bare `T<n>` token remains as a debt identifier; `DEBT-` prefixed IDs are the only debt IDs.
- No historical finding was rewritten to make an earlier audit look correct: my F11 stays withdrawn,
  RECON-1's three defects stay named, RECON-1 and both audits stay byte-identical.
- No source, register, freeze, measurement, dependency, runtime or model change accompanies this pass.

### 16.2 What is not settled by this pass

These corrections were applied from the cited source and from the review's reasoning. A second auditor
has **not** independently certified revision r2.1; the disputed line locations were cross-checked
against the blob and the repository, and the re-grades are documented above, but the revision itself
has one author.

Six questions remain open on evidence rather than on wording, and are listed in §14.2: DEBT-T9
production reachability, runtime provenance for both binaries, distribution-package notices, per-model
multi-image behaviour, `b12110`, and the four projector-quant cells. Until those are answered, nothing
here is implementation authorization and no controlling register entry changes.
