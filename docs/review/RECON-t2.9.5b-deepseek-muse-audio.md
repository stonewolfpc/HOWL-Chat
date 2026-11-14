# T2.9.5B Audio — DeepSeek / Muse Audit Reconciliation (RECON-1)

**Type:** research and documentation only. No source, register, freeze, dependency, toolchain,
runtime, model, measurement, commit, or Git-history change was made by this work.
**Date:** 2026-11-04. **HEAD:** `4f5defedb1df037052dd75ada9c5d5dea21f20de` (verified before and after).
**Upstream pin:** `a34de8618b25dfc32ee43749db47d767a9f0a25d` (fetched live this pass via
`raw.githubusercontent.com`; line numbers below are from those blobs).

---

## 1. Executive reconciliation

I re-read A1, the DeepSeek audit, and Muse's re-verification in full, re-fetched the cited
upstream files at the pinned commit, and re-grepped the current working tree. Hashes of all three
inputs match their recorded values (§2).

**Bottom line:** Muse's challenge succeeds where it matters and fails where it matters.

- **Muse is right on F11.** The `llama-mtmd-cli` NOT-production-ready label exists (block at
  `tools/mtmd/mtmd-cli.cpp:34-38`, operative sentence line 35 — verified this pass). The DeepSeek
  audit's F11 ("not located… restrict the claim to llama-tts") is **wrong and is hereby withdrawn
  by its own author**. A1 §2's table stands as written.
- **Muse confirms F1–F10.** No DeepSeek finding other than F11 is overturned. Two Muse refinements
  are accepted: the F1 imprecision is inherited from upstream's own code comment, and the F5 rule
  spans decision 572 through dispatch 579 with the image-count check at 575 (a line Muse's own
  table omits — corrected in §4).
- **Muse's D2 needs one correction (mine, §4.D2).** `HTTPClient.LoadModel` is not quite
  caller-less: `model.Manager` is constructed in production (`app.go:102`) and its `LoadModel`
  (`model/manager.go:75`) delegates to it — but nothing ever calls `Manager.LoadModel`, so the
  chain is dead in practice and would 404 against the app's own single-model server if invoked.
  The defect stands; the call-graph description is corrected.
- **One Muse claim is stale due to tree drift (§4.D-status).** `GetAudioStatus` no longer
  hardcodes all four flags: current `app-audio.go:105-135` reports live `asr_ready`/`tts_ready`
  from the audio manager while `multimodal_ready`/`realtime_ready` remain hardcoded `false`.
  The F8/A1-§7 substance (multimodal + realtime misreport) survives; the "all hardcoded" wording
  does not.
- **D3, D4, D5, D6, D7 verified as stated**, with severity challenges recorded in §4 (D4 is
  evidence hygiene, not a defect; D5's real content is narrower than its headline; D7's Bark row
  is stale-doc, not a product defect).
- **No finding changes the authorized position:** A1 §6 Tranche A (contract-first, owner
  authorization required), no upgrade, no implementation. The consolidated debt register (§11)
  has 11 items with shared root causes, not 11 independent patches.

The original research report remains **unrecovered** (Muse's search account accepted; I did not
re-run a filesystem hunt because no new evidence for its location exists — its absence is itself
the finding). Every "prior report" disposition therefore rests solely on A1 §10. Nothing in this
artifact reconstructs the original report.

---

## 2. Audit baseline and artifact hashes

| Artifact | SHA-256 (recomputed this pass) | Expected | Result |
|---|---|---|---|
| A1 `docs/review/t2.9.5b-audio-research-correction-addendum-a1.md` | `2eb69c778b6a2737a4d5415e5ab9645b0f4566c365edb75e1e2e4fb9dbfa2bbe` | `2eb69c77…9dbfa2bbe` | MATCH |
| DeepSeek audit `docs/review/AUDIT-t2.9.5b-audio-research-a1.md` | `502e98e36eccb2e08be929d172e80b22e8c1ccd0c5751e2c44ed381ce86cc8cb` | `502e98e3…381ce86cc8cb` | MATCH |
| Muse re-verification `docs/review/AUDIT-t2.9.5b-full-reverification-muse.md` | `7de11f73ef42d4e825998c5f8d5ed7c0257ecf32d4ef6967c437cf03fe9690e2` | `7de11f73…fe9690e2` | MATCH |
| Review index `docs/review/README.md` | `68b85f1a60913db40b27db834092b696b6fc9eff2787b6cf438d74cc0d71dced` | `68b85f1a…0d71dced` | MATCH |

Working tree at reconciliation time: HEAD unchanged `4f5defe`; `git status` shows pre-existing
modifications (`app-audio.go`, `app.go`, `config/*`, `internal/audio/*`, `IMAGE_PAGE_*`,
`docs/audio-architecture/README.md`, plus `build/t21`/`build/logs` untracked logs). This drift —
notably in `app-audio.go` and the audio subsystem — postdates Muse's line numbers and is the
reason §4 re-verifies every cited line rather than trusting either audit's numbers. This artifact
adds exactly one file and modifies nothing else.

---

## 3. Claim-by-claim DeepSeek × Muse comparison

Grades: **Confirmed** (source supports as written), **Corrected** (true claim, wrong
citation/wording — fix recorded), **Disputed** (audits disagree — resolved here by source),
**Unestablished** (no sufficient source), **Superseded** (replaced by later verified state).

### 3.1 Upstream source claims

| # | Claim | DeepSeek | Muse | Reconciled (this pass) |
|---|---|---|---|---|
| F1 | `/v1/realtime` branch condition | MAJOR, corrected | Confirmed + comment-provenance nuance | **Confirmed.** Fetched `server-context.cpp:6130-6135`: comment (6131-6132) uses the `?intent=transcription` shorthand; boolean (6133) requires `!asr_prompt.empty() && intent==transcription`; conversational branch (6134) needs `audio_stream && chat && !live_asr`; fallback gate `!has_audio_stream` (6138). Muse's provenance nuance accepted: A1 paraphrased the comment, not the boolean |
| F2 | realtime output modality | MAJOR, corrected | Confirmed | **Confirmed.** Muse cites `server-realtime.h:33-34,83-85`; consistent with the fetched branch code (conversation path exists independent of `audio_gen`). A1 matrix cell "audio stream + text" overstates — corrected output: "text, and/or audio per `output_modalities`/model switching; audio out requires `audio_gen`" |
| F3 | `llama-tts` gates cited to 106-130 | MAJOR, corrected | Confirmed with split table | **Confirmed.** Fetched `tts.cpp`: label block 14-18 (operative line 15), `--mmproj` gate 61-62, `MTMD_GEN_AUDIO_TYPE_NONE` reject 96-97 (both derived from my sed windows and matching both audits), speaker/prompt/init 106-130. Correction stands |
| F5 | mtmd-cli realtime rule | MINOR, cite 572-579 | Confirmed; decision 572, branch 574, error 576, dispatch 579 | **Corrected (both slightly imprecise).** Fetched `mtmd-cli.cpp:571-579`: decision `transcribes` at 572, branch at 574, image-count check at **575**, message at 576, dispatch at 579. Muse's table drops line 575. Span 572-579 covers all; exact set is 572/574/575/576/579 |
| F11 | mtmd-cli production label | UNVERIFIED (not located) | **Refuted**: label at 35-36 | **Disputed → resolved for Muse.** Fetched block: `/**` at 34, "NOT a production-ready binary" at **35**, "playground…" at 36, contributor note 37, `*/` 38. Label exists. **F11 is withdrawn.** A1 §2 row stands unchanged |
| Routes | `server.cpp:186-189` | confirmed | confirmed | **Confirmed.** Fetched: 186 transcriptions, 187 alias, 188 `/tts`, 189 WS realtime. A1's file-level citation is now closed to exact lines |
| Router `/models/load` | (not examined) | Muse D2 cites `server.cpp:529-533` | **Confirmed.** Fetched: `post /models/load` at 530-531 inside the "custom routes for router" block, registered only under `server_run_router`. Single-model spawn path never registers it |
| Gates | transcriptions `has_mtmd && allow_audio` | confirmed | confirmed | **Confirmed.** Fetched `server-context.cpp:6082-6084` exact |
| F9 | `b12110` / "95 increments" | UNVERIFIED | Confirmed as Unestablished | **Confirmed Unestablished.** Neither audit observed `b12110`; no new source this pass. The upgrade recommendation never depended on it |

### 3.2 Repository claims

| # | Claim | DeepSeek | Muse | Reconciled (this pass, current tree) |
|---|---|---|---|---|
| F6 | streaming flags | MINOR (Small=false) | Confirmed | **Confirmed.** Current `asr_models.go:50,79` true; `:108` false; `:137,166` true. A1 §7 must read "Tiny/Base + both Qwen entries" |
| F7 | `app-model.go` mmproj block | range 118-125 | Confirmed | **Confirmed.** Current lines 118-125 exact (`mmprojPath :=` 118, `HasVision` 120/124). Note: line 121 assigns `profile.MMPProjPath` (double-P) — a possible field-name typo, recorded as debt T12, not investigated further here |
| F8 | `app-audio.go` hardcoded flags | 105-135 or 110-111/127-128 | Confirmed | **Corrected for drift.** Current `app-audio.go:105-135`: nil branch hardcodes all four (110-113); live branch reports real `asr_ready`/`tts_ready` (128-129) but still hardcodes `multimodal_ready:false, realtime_ready:false` (130-131). Muse's audit predates the live-status change |
| A1 §1.3 | zero app-side audio route calls | confirmed | Confirmed with qualification | **Confirmed.** No `input_audio`/`/v1/realtime`/`/tts`/transcriptions construction in `internal/backend`. Muse's qualification (mmproj auto-detect is vision plumbing an audio projector could satisfy by filename) is accepted and folded into debt T1 |
| Whisper binary | 0 files bundled | confirmed | Extended (downloader dormant) | **Confirmed + extended.** `NewBinaryDownloader`/`DownloadLatest`/`Acquire` defined (`downloader.go:22-36`, `downloader_acquire.go:29,45,132`) with **zero production callers** (only self-delegation at `downloader.go:36`). `findWhisperCLI` "never downloads or builds" (`recognizer.go:226-228`, current tree). Distributor question narrowed to dormant code — kept alive, not closed |
| A1 §11.2 | `ggml-cuda.dll` outlier | confirmed | Extended + D4 | **Confirmed.** Current listing: 594,123,776 B @ 2026-10-30 15:48 vs set 10:46-10:49; only notice file remains `LICENSE-LLVM-OpenMP` (19,741 B). Provenance still Unestablished; closure criteria in §5/§11 |
| F10 | four families' mmproj quants | UNVERIFIED | Confirmed as Unestablished | **Confirmed Unestablished.** No per-tree re-fetch was possible in this environment (no Hub access configured here); Muse's secondary data is the best available and stays secondary |

### 3.3 Resolved tallies

- DeepSeek findings surviving: F1, F2, F3, F4, F5, F6, F7, F8, F9, F10 (10 of 11).
- **F11 withdrawn** (author's own correction — the label was one fetch away and I failed to fetch that file).
- Muse D-findings: D1 accepted (it is the F11 refutation); D2 accepted with call-graph correction; D3, D5, D6, D7 accepted with severity challenges noted; D4 accepted as hygiene, not defect.

---

## 4. Muse D1–D7: independent investigation

### D1 — mtmd-cli label (Muse correct; DeepSeek F11 withdrawn)

Verified above (§3.1). Root cause of my miss: I fetched `tts.cpp` but never fetched
`mtmd-cli.cpp` in the audit pass. No systemic cause; a plain coverage gap. Regression: any
future "label not located" claim must name the files actually fetched. No user-visible impact
(the label governs upstream CLI use, not the app). **Status: closed, A1 unchanged.**

### D2 — router-only `POST /models/load` + `LoadModelWithMMProj` ignores its argument

**Verified real defect; call-graph description corrected.** Current sources:
`http_client.go:276-313` (`LoadModel` POSTs `{base}/models/load`, auto-attaches found mmproj at
289-292); `:393-397` (`LoadModelWithMMProj` ignores `mmprojPath`, delegates with `nil`
options); `model/manager.go:75` delegates to `HTTPClient.LoadModel`. Production wiring:
`model.NewManager(client)` is constructed (`app.go:102`) and its `IsModelLoaded` is consulted
(`chat/service.go:449,526`) — but **`Manager.LoadModel` itself has zero callers**
(repo-wide grep: only its definition and the internal `http_client.go:396` self-call). So the
live model path is `App.LoadModel` → kill + respawn (`app-model.go:27+`), while the entire
`Manager.LoadModel → HTTP /models/load` chain is reachable-but-never-invoked dead code that
would 404 against the app's own server. Root cause: most likely a Gollama-era client shape (or
aspirational router-mode support) carried across the migration — consistent with §8, but the
pre-migration tree is unavailable so this stays **Inferred**, not established. User-visible
impact today: none (never invoked); latent risk: any future caller gets a silent 404-shaped
failure plus a dropped mmproj argument. Regression coverage required: a no-caller assertion for
`/models/load` against single-model spawn, and an options-passthrough test if the method is
ever kept. Obsolete on cleanup: both methods (or re-targeting to the spawn path) — owner
decision, not this artifact.

### D3 — single-image + text-only history ceiling

**Verified real constraint; severity accepted, root cause refined.** Current sources:
`app-prompt.go:25` (single `imageData: string` Wails parameter); `http_client.go:576-583`
(single `image_url` part, one of exactly two construction sites — `:582`, `:712`);
`buildPromptMessages` (`app-prompt.go:620-745`) emits system + roleplay context + text history
with **zero image parts** (current-tree grep confirms). So: at most one image per request, and
turn N+1 cannot see turn N's image — the saved conversation and the model's context diverge
silently. This is **independent of runtime capability** (upstream loops over content parts; the
server side is not the limiter) and independent of the Gollama migration (no pre-migration code
exists to compare — "lost vs never-had" is Unestablished). It is also independent of audio work
but shares the pipeline the audio tranches must extend (Tranche B touches the same
`buildRequest*` functions). User-visible impact: multi-image prompts impossible; second-turn
references to an earlier image ("as shown above") silently fail. Regression coverage: one-image
cap assertion, history-image-carriage assertion (currently asserting absence — the contract to
beat), MIME fidelity (see D7-analogue below). Fix is additive (multi-part construction +
history carriage), never a rebuild — consistent with §8.

### D4 — `llama-connect.exe` Sep-2 outlier

**Verified fact; severity challenged: hygiene note, not a defect.** Current listing:
`llama-connect.exe` 11,473,920 B @ 2026-09-02 17:14 amid an Oct-30 set, plus the known
`ggml-cuda.dll` 594 MB @ Oct-30 15:48. Neither file's provenance is established, and no
`.src_sha`, receipt, or SBOM exists for the runtime set. But a timestamp alone proves nothing
about ABI or function, and `llama-connect` is a connectivity helper, not the inference path.
Correct disposition: fold into the single provenance record (§11-T5), do not treat as a second
independent defect. I accept Muse's conflation warning about the root-level SD-stack DLLs
(different stack, different files) — no evidence of cross-stack mixing was found.

### D5 — Go directive 1.23.8 vs installed 1.27.3

**Verified fact; headline challenged as misleading.** Current `go.mod:3` declares `go 1.23.8`;
installed toolchain reports `go1.27.3`. No `.github/workflows/` exists (empty — no CI matrix to
consult), so "build scripts use a different toolchain" is Unestablished beyond the local
install. The directive still governs language semantics and `go vet` stdversion checks for the
module regardless of toolchain, so the premise "Go never updated" is only half-stale: the
*toolchain* moved, the *module contract* did not. No dependency incompatibility was observed
(this pass ran no builds by authorization). Disposition: informational; any directive bump is a
staged, separately-tested change (§9). Not debt in the "carried-forward defect" sense.

### D6 — kill-respawn + exclusive port 8080

**Verified intentional-but-fragile lifecycle; "defect" downgraded to known fragility.**
Current `app-model.go:27-90`: `LoadModel` requires port 8080 free (`waitForPortFree`, hard fail
at 79-82), kills the existing server, respawns with `-m` (+ `--mmproj` at 196-201 per both
audits; spawn at `app-llama.go:176-257`). Crash-restart path re-checks the port (254-256).
This is the app's *actual* model lifecycle — simple, legible, and consistent with single-model
`llama-server`. Its costs are real: every switch drops chat availability briefly, and any port
occupant hard-fails the load. But no queue, no hot-swap, and no port agility exist upstream in
single-model mode to inherit — this is app design, not migration residue. It interacts with A1
§7's resource rule (no mechanism yet for "refuse-or-queue instead of evicting chat") and with
any future audio-model sidecar. Regression coverage: port-occupied refusal message, respawn
restores prior model, no second server process leaks. No obsolete code identified.

### D7 — stale `docs/audio-architecture/README.md` (Bark, PENDING stages)

**Verified stale doc; severity challenged: doc defect, not product defect.** Current lines:
stages 2-5 marked PENDING (`:64-80`) while ASR/TTS/transcribe-UI/single-image-vision code
exists; TTS matrix lists Bark (`:207`) with zero registered Bark loader
(`pipeline.go:438-446` registers Piper/MeloTTS/Kokoro in both audits; unchanged in current
tree per A1-§7 evidence); FAQ advises "upgrade to MeloTTS or Bark" (`:387`). Deeper wrinkle
Muse underplays: the doc's ASR sizes (39/74/244 MB at `:195-197`) describe upstream *bin*
weights while `model_resolver.go:69-76` demands `ggml-*.bin` files — the doc and the resolver
do not even agree on what "tiny/base/small" means as artifacts. Root cause: May-2026 plan
frozen in place while code moved on. User-visible impact: indirect (a reader, not the app,
is misled) — except `frontend/dist` tutorial prose (§8), which ships to users. Disposition:
one-line status header when docs are next touched; no rewrite authorized here. Related
components also affected: `asr-repair-design.md`, `STAGE_STATUS.md`, and `docs/review/*`
overlap with no stated precedence — the register (r14) is the only controlling artifact.

---

## 5. A1 §3-§4 items re-checked (requirement 3A)

- **Realtime routing + `live_asr`:** reconciled wording — "Conversational branch when
  `audio_stream && chat` and NOT (`asr_prompt` non-empty AND `intent=transcription`); otherwise
  the transcription path, which additionally requires `has_audio_stream`." (sources §3.1).
- **Realtime output:** "text, and/or audio per `output_modalities`/model switching; audio out
  requires `audio_gen`." A1 §4 cell corrected; A1 §4 negatives already say this.
- **TTS citation:** split to 15 / 61-62 / 96-97 / 106-130. Claims true; references fixed.
- **Model families:** no change to A1 §3 rows beyond the audit's accepted corrections
  (Realtime = 4B naming; Qwen3-Omni documented incl. Thinking variant per Muse's multimodal.md
  read). Per-family quant cells stay E1-unverified (F10).
- **Projector rule:** standing — `mmproj` presence never identifies projector kind; `audio_gen` /
  `audio_stream` / `chat` are negotiated flags. Current `app-model.go:118-125` still violates it
  (`HasVision` from filename).
- **Licenses:** F4 re-cut accepted (Qwen3-ASR/TTS + Voxtral Apache-2.0; Qwen2.5-Omni + Qwen3-Omni
  custom; Gemma 4 SPDX-vs-licence-page mismatch sharper than A1's hedge; LFM 1.0; ggml-org
  converts license-less). Runtime notice gap stands and is understated in A1 (vendored
  miniaudio/cpp-httplib second-order notices also missing).
- **Pocket/Soprano:** identified per DeepSeek audit (upstream `kyutai/pocket-tts`,
  `ekwek/Soprano-80M` + `Soprano-1.1-80M`), still D+S at best, naming resolved, verification open.
- **Voxtral Realtime identity:** `Voxtral-Mini-4B-Realtime-2602` (4B, not 3B) — A1-era naming
  corrected.
- **Build numbers:** `b12015-a34de8618` observed (fingerprints in project logs); `b12110`
  unobserved. Distance arithmetic conditional only.
- **mtmd-cli wording:** F11 withdrawn; A1 stands.
- **Whisper/downloader/redistribution:** no bundled binary; acquisition code dormant with zero
  production callers; `findWhisperCLI` never downloads. Distributor question narrowed, not closed.

---

## 6. Reconciled audio capability matrix

Grades as in Muse §5. Nothing below E4 (no execution was authorized or performed).

| Capability | Grade | Limiting layer + evidence |
|---|---|---|
| Ordinary text chat | Working (not E2E here) | single-model spawn `-m` (`app-llama.go:182-191` per audits); not re-probed |
| Vision, single image, current turn | Constrained | `app-prompt.go:25`; `http_client.go:576-583,703-713`; `HasVision` gate `app-model.go:118-125`, `app-prompt.go:38`; `--mmproj` at spawn `app-llama.go:196-201` |
| Multiple images / request | Not implemented | structural: one `[]byte` → one `image_url` part; exactly two construction sites (`:582`, `:712`) |
| Image in history | Not implemented | `buildPromptMessages` 620-745 emits zero image parts (current-tree verified) — history and model context diverge from turn 2 |
| Image MIME fidelity | Latent defect | both sites hardcode `data:image/jpeg;base64,` regardless of bytes (`:583`, `:713`) |
| Local file ASR (Whisper) | Partial | `Transcribe` path exists; registry tiny/base/small only (`model_resolver.go:69-76`) vs UI offering medium/large + Qwen-API (`audio-settings.html:96`) — UI-exposed-unsupported for the rest |
| Live/streaming ASR | UI-exposed-unsupported | both adapters explicitly reject (`recognizer.go:196-202`, `qwen/client.go:95-101`) while registry declares `SupportsStreaming:true` for Tiny/Base/Qwen (`asr_models.go:50,79,137,166`) |
| Direct audio understanding | Not implemented | zero `input_audio` construction in `internal/backend` |
| Realtime speech-to-speech | Not implemented | no WS client in `internal/backend/llama`; `realtime_ready:false` hardcoded even in the new live-status shape (`app-audio.go:130-131`) |
| Native TTS (`audio_gen`) | Not implemented | no `/tts` caller; Piper/Kokoro/MeloTTS/browser independent of llama.cpp |
| Status honesty | Partial, improving | `asr_ready`/`tts_ready` now live from manager; `multimodal_ready`/`realtime_ready` still hardcoded false |
| Capability detection | Partial + misleading | filename inference (T1); no audio capability field anywhere |
| Model lifecycle | Working, fragile | kill + respawn; port-8080 exclusive with hard fail (`app-model.go:78-82`); crash-restart re-check (`:254-256`) |
| Whisper acquisition | Dormant code | zero production callers; never-downloads finder |

---

## 7. Application capability and multimodal-history findings

The history investigation (requirement: does the record preserve what was actually sent?)
returns a clean, bounded answer:

1. `buildPromptMessages` builds system + roleplay context + **text** history only. The current
   turn's image travels a separate path (`SendMessageStreamWithImage` → single-part request),
   and the line-704-745 tail appends history as text with an explicit "untouched" comment for
   saved history. **The model's context from turn 2 onward does not contain turn 1's image.**
   Saved chat history presumably records that an image was sent (not traced this pass —
   recorded as a gap, §14.5 of Muse's artifact stands), but the *inference* record diverges
   regardless.
2. This is independent of runtime capability (upstream accepts multi-part content), independent
   of Gollama (no pre-migration comparator exists), and independent of audio — but it is the
   same `buildRequest*` surface Tranche B must extend, so audio planning inherits the
   constraint: any `input_audio` history carriage must be designed, not assumed.
3. Vision is therefore **working but artificially constrained**: single-image, current-turn,
   JPEG-labelled. "Artificially" is precise — the ceiling is three app functions, not the
   server, the model, or the migration.

---

## 8. Backend migration findings (Gollama → llama.cpp)

Residual references re-verified in the current tree: `prompt/builder.go:6`,
`types/model_profile.go:97,144` (comment-only), `frontend/dist/js/tutorials-data-model.js`
(+ `frontend/dist/index.html` bundle strings) — user-visible stale prose in the shipped bundle.
No Gollama plugin import, cgo binding, or Gollama-API call shape exists in Go code; the backend
is `exec` + HTTP against `llama-server` (`app-llama.go:171-257`).

On the migration-loss question (requirement 4): **no evidence of lost capability retrievable
from this tree.** The pre-migration code is absent, so "vision worked under Gollama" is
testimonial with unrecoverable mechanism. What *is* established: the current ceilings (D3, MIME,
port lifecycle) are all explicable as native app-layer design around single-model
`llama-server`, and D2 is the only item with a plausible migration-leftover shape — graded
Inferred, not established. The common-root-cause hunt therefore ends here: **T1 (filename
inference) is the one shared root** behind the vision gate, the audio-projector blind spot,
and the status/UI/resolver disagreement. Everything else has its own local cause (§11).

The migration's genuine gains (reasoning streaming, `/props`, `/v1/models`, multi-arch GGUF)
stand unchallenged; nothing here recommends against llama.cpp.

---

## 9. Go and llama.cpp upgrade assessment

**Go.** Declared `go 1.23.8` (`go.mod:3`); installed `go1.27.3`; no `.github/workflows/`
directory exists, so no CI toolchain pin was found to conflict. Muse's D5 is accepted as fact
with the headline corrected: toolchain moved, module contract did not. No incompatibility
observed (no build run — unauthorized). Supported-target and migration-work questions
**cannot be answered from this pass** (no CI matrix, no `wails.json` build settings enumerated
here, no vet/build output). Position: leave the directive alone; any bump is a staged change
with its own build + vet + smoke-test plan, never bundled with a runtime move. No upgrade
performed, none authorized.

**llama.cpp runtime.** Current: `b12015-a34de8618` (2026-10-30). Candidate: v0.9.0 (`cd959bae`,
2026-11-04) adding Chatterbox TTS, Nemotron 3.5 Transcribe, predefined Pocket voices, voice APIs
in `mtmd.h`, `MTMD_VIDEO_TEMPLATE_QWEN3VL`, ggml v0.29.0 (batched copy API superseding `_2d`),
miniaudio 0.11.26 / cpp-httplib 0.64.0 (release-page-derived, secondary here). Assessment:
additive MINOR with no stated `llama.h`/GGUF/mmproj break, but it touches `/props` shape and
deprecates a ggml copy API — adoption needs a named-gap justification plus a verification plan
(props-shape check, `_2d` call-site sweep, single-build re-hashing). **No upgrade performed,
none recommended now**: nothing in the Whisper + native-audio-input scope (A1 §5) requires
v0.9.0, and the provenance record (§11-T5) must be clean before any runtime changes anyway.

---

## 10. Licensing and redistribution position

- **Runtime notices: packaging defect, confirmed.** `runtime/llama/` ships ~40 binaries with
  only `LICENSE-LLVM-OpenMP`. Upstream llama.cpp is MIT ("Copyright (c) 2023-2026 The ggml
  authors" — secondary basis, consistent across both audits); the MIT copyright + permission
  text must accompany the binaries. A1 understates by omitting vendored miniaudio/cpp-httplib
  second-order notices. Cure is packaging, not research. No legal conclusion beyond the
  artifacts is rendered.
- **Whisper binaries: not redistributed.** No `whisper*` binary outside `build/` (both audits +
  current tree). Acquisition code (`downloader*.go`) is dormant with zero production callers.
  If ever wired to a UI affordance, the app becomes a download-on-behalf distributor of
  MIT-licensed binaries (notice duty follows the bits) — decision point recorded, not resolved.
- **Weight terms (per-family, SPDX ≠ terms):** Qwen3-ASR/TTS + Voxtral Apache-2.0; Qwen2.5-Omni
  + Qwen3-Omni custom; Gemma 4 Apache-2.0 tag with a Gemma-licence link (do not rely on the
  tag); LFM 1.0 (attribution/acceptable-use may reach UI); gated `kyutai/pocket-tts`
  (cc-by-4.0); several ggml-org convert repos with **no license field** (derivative-terms gap).
  Per-family review required before any bundling or download-on-behalf feature. Weights are not
  bundled today (E1).

---

## 11. Consolidated technical-debt register

Standing rule preserved: no knowingly carried-forward digital debt — every item below has a
bounded fix surface or an explicit owner decision recorded. Any future modularization follows
COPY → PA decided → PASTE → EDIT → VERIFY → DELETE THE OLD; untested replacement must never
displace working code.

| # | Item | Root cause | Evidence | User-visible impact | Corrective direction → dependencies → regression cover |
|---|---|---|---|---|---|
| T1 | Capability inferred from projector filename | no negotiated-capability record | `app-model.go:118-125`; candidate list + wildcard `http_client.go:318-391` (per audits) | audio-mmproj would falsely set `HasVision`; vision-gated UI misreports | negotiated record (`/props` + family knowledge), never filename → before any native audio admission → audio-mmproj-present ⇒ vision-false test |
| T2 | Declared-but-unsupported streaming | registry rows added without adapter check | `asr_models.go:50,79,137,166` vs `recognizer.go:196-202`, `qwen/client.go:95-101` | user selects streaming, gets an error (or silent non-stream) | one-line-per-row correction → none → registry-vs-adapter contract test |
| T3 | Status/registry/UI triple disagreement | three independent sources of truth | UI `:96` vs `model_resolver.go:69-76` vs `app-audio.go:105-135` (now partially live) | UI offers medium/large/Qwen-API the resolver rejects | single capability source feeding all three (= A1 Tranche-A item 1) → after T1 → triple-agreement test. Largest sequencing risk (touches UI contract) |
| T4 | Dead `/models/load` chain + ignored mmproj arg | carried-over client shape (Inferred Gollama-era/aspirational) | `http_client.go:276-313,393-397`; `manager.go:75`; zero callers of `Manager.LoadModel` | none today; 404 + dropped arg if ever called | delete or re-target to spawn path → owner cleanup decision → no-caller proof re-check + options-passthrough test if kept |
| T5 | Runtime provenance record | no receipt/SBOM at bundle time | `ggml-cuda.dll` 15:48 outlier; `llama-connect.exe` Sep-2 outlier; notice absence (§10) | none directly; every load claim inherits E3-I | record source URL + hash or replace from pinned asset; ship MIT notices → before any runtime-dependent audio work → hash-recorded single-build assertion |
| T6 | Single-image / text-only history | app-layer single-part construction | `app-prompt.go:25,620-745`; `http_client.go:576-583,703-713` | multi-image impossible; turn-2 image references silently fail | additive multi-part + history carriage → design + authorization → cap + carriage + MIME-fidelity tests |
| T7 | JPEG MIME hardcode | copy-paste default at both sites | `http_client.go:583,713` | benign today (server sniffs content), latent corruption risk | sniff-or-label at two call sites → with T6 → byte-fidelity test |
| T8 | Port-8080 exclusivity, kill-respawn | single-model server design, no queue | `app-model.go:27-90,254-256` | switch drops chat briefly; occupant hard-fails load | refuse-or-queue policy per A1 §7 → resource-arbitration design → refusal-message + respawn-restores + no-leak tests |
| T9 | Repeated Whisper `Initialize` orphans temp dir | missing one-shot guard | STAGE_STATUS F1 (self-recorded); `recognizer.go` Release/tempDir (current tree confirms handling exists) | disk litter on repeated init | one-shot init guard (as STAGE_STATUS proposes) → small, isolated → repeat-init no-orphan test |
| T10 | Stale docs + gollama prose in bundle | plan docs frozen while code moved | `docs/audio-architecture/README.md:64-80,195-207,387`; `tutorials-data-model.js`; dist bundle | misled readers; user-visible tutorial staleness | one-line status header; fix tutorial strings on next touch → owner docs decision → doc-vs-registry consistency check |
| T12 | `profile.MMPProjPath` (double-P) | probable field-name typo | `app-model.go:121` (current tree) | none if consistent; confusion + serialization mismatch if not | confirm-or-rename with full-repo reference sweep → trivial, isolated → compile + grep-clean |

Shared-root summary: T1 is the root behind the detection family (T1→T3); T6/T7 share the
`buildRequest*` surface; T4/T5 share the runtime-client boundary; T8 stands with the resource
policy; T9/T10/T12 are isolated. **No unavoidable debt identified.**

---

## 12. Dependency-ordered remediation recommendations

Recommendations only — no authorization is granted or implied by this artifact.

1. **Errata note (no edits to A1):** F1 branch wording, F2 output cell, F3 split citation, F4
   table re-cut, F5 span 572-579, F6 "Tiny/Base + Qwen", F7/F8 current ranges, F11 withdrawal.
   A1 stays byte-identical (review-chain anchor).
2. **Packaging cure (T5-license half):** ship llama.cpp/ggml MIT + vendored-lib notices. No code
   change, no upgrade.
3. **Provenance closure (T5):** source URL + hash for the 594 MB `ggml-cuda.dll` (or replace
   from a pinned official asset); disposition `llama-connect.exe` in the same record; re-hash
   the single-build set. Gates all runtime-dependent audio work.
4. **Contract-first Tranche A** per A1 §6, with T1→T3→T2 as the ordered defect backlog and
   T4/T6/T7/T8 scoped as its ceiling. Tranche B/C/D unchanged and additive.
5. **Docs hygiene (T10):** status header on `docs/audio-architecture/README.md`; tutorial
   strings on next touch; register r14 remains the sole controlling artifact.
6. **Lifecycle + history (T8, T6):** refuse-or-queue policy; additive multi-part/history
   carriage — designed, not assumed, and only under explicit authorization.
7. **Go directive:** leave `go 1.23.8` until a staged toolchain pass (CI enumeration + build +
   vet + smoke), never bundled with a runtime move.

---

## 13. Required acceptance and regression coverage

For any future authorized implementation tranche (not this artifact): no-caller proof for
`/models/load` against single-model spawn; options/mmproj passthrough test if kept;
audio-mmproj ⇒ vision-false; registry-vs-adapter streaming agreement; status/UI/resolver
triple agreement; one-image cap + history-carriage + MIME-fidelity assertions; port-occupied
refusal message; respawn-restores-prior-model; no leaked second server process; repeat-Whisper-init
no-orphan; hash-recorded single-build assertion; doc-vs-registry consistency check. E4 (load +
negotiated capability on this host) is the minimum bar for admitting any native provider; E5
(correct output on a user request) for user-facing claims.

---

## 14. Explicit owner decisions still required

1. Authorize or decline the A1-errata note (and its exact text).
2. Authorize or decline the T4 cleanup (delete vs re-target).
3. Authorize or decline wiring the dormant whisper downloader — with the distributor
   determination per family made first.
4. Authorize or decline any native ASR admission (A1 Tranche A entry criteria stand:
   single-build hash record + Qwen3-ASR main + audio mmproj + passing transcription probe).
5. Authorize or decline the T8 refuse-or-queue policy and the T6 history-carriage design.
6. Authorize or decline the Go directive bump and its test plan.
7. Authorize or decline any llama.cpp upgrade (named-gap justification required).
8. Decide the docs-precedence question (register vs STAGE_STATUS vs audio-architecture vs
   review artifacts) or explicitly leave the register as sole control.
9. Dispose IR1 / MEAS1-execution / T6 / CL3 per the controlling register process — none is
   touched by this artifact.

No implementation, measurement, upgrade, or authorization-state change is effected here.

---

## 15. Integrity close

- Inputs verified by hash before reliance (§2, all MATCH). A1, the DeepSeek audit, the review
  index, all register revisions, freeze records, and all source files are byte-identical to
  their pre-reconciliation state, except for the working-tree modifications that predate this
  task (listed in §2) and the single file below.
- **Executed:** file reads, `grep`/`sed` line reads, `sha256sum`, `ls`, `git rev-parse/status`
  (read-only), `go version` (no build), live HTTPS fetches of upstream blobs at `a34de8618`
  into `/tmp` (removed at close). **Never:** builds, tests, servers, model loads, inference,
  measurements, register/freeze edits, commits, deletions, or source modifications.
- **Created:** exactly one file — `docs/review/RECON-t2.9.5b-deepseek-muse-audio.md` (this
  artifact). Its path, SHA-256, size, HEAD, and working-tree state are reported in the
  completion message as part of the audit record.
