# MissingImageRequirementsREAD

Prioritized implementation guide for HOWL Chat image generation capability gaps.

Audit date: November 4, 2026. Sources: `app-*.go`, `internal/image/...`, `internal/backend/...`, `go.mod`, `go.sum`, `*.txt` help/version dumps, `runtime/` binaries, `.sd-upgrade/`.

Method: severity = how much the gap hurts usability or correctness **today**. Stages group work by dependency and user-facing value. Useability is the primary axis — a feature users can't reach or that silently produces wrong results rates higher than one that is merely incomplete.

---

## Severity legend

| Rating | Meaning |
|--------|---------|
| **P0 — Blocked** | Core flow doesn't work or silently produces wrong results; users hit this on a normal first session |
| **P1 — Painful** | Real workflow is possible but a common action fails, is confusing, or costs far more time/VRAM than it should |
| **P2 — Annoying** | Works but missing a companion UI, a polish item, or a safety check that real users will hit |
| **P3 — Nice to have** | Existing behavior is acceptable; improvement is polish or power-user escape hatch |

---

## Recommended stages

| Stage | Theme | Goal |
|-------|-------|------|
| **Stage 0 — Unblock** | Fix the things that silently break a normal session | A new user with a fresh Flux/SDXL/SD3.5 download can generate, img2img, inpaint, and outpaint without a hidden failure |
| **Stage 1 — Usable core** | Make the existing txt2img/img2img/inpaint/outpaint/video flows feel complete | Live previews where expected, sensible defaults, no silent traps, mask/control-image inputs users can actually produce |
| **Stage 2 — ControlNet & IP usable** | Make conditioning actually usable from downloaded assets | Preprocessed control images, IP-Adapter/photo-reference paths, per-asset type visibility, multiple-condition support |
| **Stage 3 — Model library & deps** | Make downloading/installing companion assets painless | HuggingFace + Civitai browsing, automatic companion resolution, trust that the right VAE/CLIP/T5/LLM was attached |
| **Stage 4 — Polish & power features** | Expose the advanced knobs already in the code, plus durability items | LR/write, sampling, caching, SLG, SCM, tensor-type-rules, per-component offload, upscaling, video quality |

---

## Stage 0 — Unblock (P0/P1)

Fix first. These are the items that make a normal session fail or produce wrong output without the user understanding why.

### P0 — SD3.5 blocked by conservative version check
- **File:** `internal/image/engine/sdcpp` version check / `isFamilySupportedByBackend`
- **Symptom:** SD3.5 models refuse to load even though the running sd-server (commit `c6dfa9d`) very likely supports them. The code comment says "For now, be conservative and require user to update."
- **Impact:** A user who downloaded SD3.5-medium sees a hard refusal. This is a real family people want *now*.
- **Fix:** Reconcile the version check with the actual `c6dfa9d` flag set (sd-help-current.txt shows `--prediction sd3_flow`, `--slg-scale`, joint-blocks support, etc. all present). Either raise the old-commit gate or document exactly which commit is the real floor for SD3.5. Do **not** leave a hard block that the backend does not actually have.
- **Dependency note:** Before flipping this on, confirm the sd-server in use actually loads an SD3.5 GGUF. The upgrade in `.sd-upgrade/extracted/` may be the better binary to ship if it is newer than `c6dfa9d`.

### P0 — img2img and inpaint have no live preview
- **Files:** `internal/image/engine/sdcpp/client.go` (`txt2img` uses sync `/sdapi/v1/txt2img`; `img2img`/`inpaint` use async `/sdcpp/v1/img_gen`)
- **Symptom:** txt2img gets a live preview because it uses the synchronous endpoint; img2img and inpaint use the async job endpoint, which the code itself says does not return previews. So img2img and inpaint show nothing until the job finishes.
- **Impact:** img2img/inpaint feel broken relative to txt2img; users assume the engine is hung. This is the single most visible usability gap in the image flows.
- **Fix options (pick one, in order of value):**
  1. If the async job endpoint can be made to emit `preview_image` in this sd.cpp build, consume it in `pollJob` the same way txt2img consumes its preview. The `sdcppJob` struct already has `PreviewImage`; the polling path already reads it. The gap is that img2img/inpaint never get a preview through that path in practice.
  2. If previews on the async path are genuinely unavailable in `c6dfa9d`, document that img2img/inpaint are fire-and-wait and decide whether to (a) switch them to the sync endpoint too, or (b) add an explicit "no live preview for transform modes" note in the UI.
- **Why this is P0:** It is the first thing a user notices when they try anything other than txt2img.

### P0 — No mask creation tools (inpaint is unusable without one)
- **Files:** inpaint path in `client.go` + `app-image-generate.go`; no mask tooling anywhere in the image packages
- **Symptom:** Inpaint requires a pre-made mask PNG. There is no brush, no SAM-style segmentation, no auto-mask from a selection. If a user does not already have a mask file, inpaint is not usable.
- **Impact:** Inpaint exists in code but is effectively a power-user-only feature until there is any way to produce a mask inside the app.
- **Fix:** Add at least one way to produce a mask. Minimum viable: a brush-based mask editor on the source image. Better: any SAM/segmentation integration that can propose a mask from a click or box. This is a frontend+image-prep concern, not an sd.cpp concern — sd.cpp already accepts `MaskPath`.

### P0 — ControlNet is unusable without a preprocessor
- **Files:** `client.go` (ControlNet fields exist and are sent), `scanner/` (ControlNets classified and listed), no preprocessor package
- **Symptom:** A downloaded ControlNet is listed and can be attached, but the app never produces the Canny/Depth/Scribble/HED/Occlusion/etc. image that the ControlNet actually needs. The user must supply a pre-processed control image from another tool.
- **Impact:** ControlNet feels non-functional to anyone who doesn't already have a preprocessor pipeline elsewhere. This is the same class of gap as the mask tooling.
- **Fix:** Add a preprocessor step that turns a raw source image into a control image for the selected ControlNet type. This can be a Go+C leader (OpenCV-style edge/depth heuristics) or a small bundled preprocessor binary, but there must be a path from "I have a photo" to "here is the control image." Until then, ControlNet is listable but not practically usable.

### P1 — No ffmpeg; MP4/GIF video output is rejected outright
- **Files:** `internal/image/engine/sdcpp/video.go` `videoContainer()`
- **Symptom:** `videoContainer()` returns an error for `mp4`/`gif` saying ffmpeg is not bundled. The only exportable containers are WebM, AVI, animated WebP.
- **Impact:** Users who want to share a video cannot get the most portable format from the app. WebM is fine for some audiences, not for many.
- **Fix:** Decide policy. Two realistic paths:
  1. Bundle ffmpeg (or a minimal WebM/MP4 writer) so MP4 is available.
  2. If bundling ffmpeg is off the table for now, at least stop presenting MP4 as if it were a first-class option; make the supported set obvious and ensure WebM/AVI/WebP are easy to play/download.
- **Dependency:** This is a runtime bundling decision, not a code logic gap. Treat it as a packaging task with a clear supported-formats surface.

### P1 — SD3.5 version check interacts with the unapplied sd.cpp upgrade
- **Files:** `.sd-upgrade/` + the version-check code
- **Symptom:** There is a newer sd.cpp build sitting in `.sd-upgrade/extracted/` (its DLLs are different sizes from the backup) that was downloaded and extracted but never swapped in. Meanwhile SD3.5 is blocked on the assumption that the current binary is too old.
- **Impact:** The two issues compound: the app both blocks a family and has a promising newer binary it never tried.
- **Fix:** Treat the upgrade as a real candidate, not a stale artifact. Verify the extracted `sd-server.exe` commit (run it with `--help` and read the commit line), then decide whether to ship it. If it is newer and supports the families you want, make it the default and retire the old backup path. If it is not usable, update the version gate and documentation to reflect the real floor instead of the current conservative refusal.

---

## Stage 1 — Usable core (mostly P1/P2)

Once Stage 0 is cleared, these make the everyday flows feel complete rather than partial.

### P1 — Masked-content and edge behavior are present but under-promoted
- **Files:** `types/types.go` (`MaskedContent`, `EdgeBehavior`, `MaskBlur`, `MaskInvert`, `MaskFeather`), `client.go` inpaint body
- **Symptom:** The inpaint mask semantics exist in the API (`fill`/`preserve`/`replace`/`only_modify`, edge `mirror`/`blur`/`noise`/`solid`), but a user has no guidance on which to pick.
- **Fix:** Add short in-app labels/examples for each masked-content mode and edge behavior. This is documentation+UI, low risk, high clarity payoff.

### P1 — Outpainting is img2img in disguise, with no seam handling
- **Files:** `app-image-generate.go` routing (`outpaint` → img2img), `client.go`, `video.go`-adjacent outpaint options
- **Symptom:** Outpaint is implemented as img2img on an expanded canvas. There is no seam blending or content-aware fill between the original region and the expanded region; `EdgeBehavior` is the only blending lever.
- **Impact:** Outpaint works but seams are visible and there is no iterative four-direction expansion — only one shot per direction or all-directions at once.
- **Fix:** After Stage 0, consider one or both of:
  1. A simple seam-blend post-pass (feather/mirror/blur already exist as options; make them do visible work for the seam).
  2. Iterative outpaint: expand one side, gen, repeat, so users can grow an image in steps. This is more valuable than a second outpaint mode flag.

### P1 — Highres fix params exist but are not surfaced
- **Files:** `types/types.go` (`EnableHR`, `HRScale`, `HRUpscaler`, `HiresWidth`, `HiresHeight`, `HiresSteps`, `HiresUpscaleTileSize`), request structs
- **Symptom:** All highres-fix fields are in the options and the sd.cpp flag set supports them, but there is no UI flow for highres fix. A user cannot request it through the normal interface.
- **Fix:** Add a highres-fix toggle + scale + upscaler picker to the generate UI. This is one of the highest-value "already works, just hidden" items.

### P2 — Sampler/RNG per-generation control is buried
- **Files:** `types/types.go` (`RNG`, `SamplerRNG`), request structs
- **Symptom:** `--rng` and `--sampler-rng` exist and are forwarded, but there is no obvious per-generation control for them.
- **Fix:** Expose sampler + RNG in the generate UI as an advanced section. Low effort, helps power users, no downside.

### P2 — TAC/SD/decode-only preview path is correct but worth calling out
- **Files:** `engine.go` buildServerArgs (`--taesd`, `--taesd-preview-only`), `client.go` `Preview: "tae"`
- **Symptom:** TAESD is correctly scoped to preview only, so the final image is still decoded by the real VAE. This is good. It just isn't surfaced as a user-visible choice.
- **Fix:** If you ever expose preview quality, make clear that TAESD affects preview only. No correctness change needed — this is a documentation/UI note, not a fix.

---

## Stage 2 — ControlNet & IP usable (P1/P2)

After Stage 1, these make conditioning pleasant rather than theoretical.

### P1 — Separate ControlNets from IP-Adapters in the UI and metadata
- **Files:** `scanner/classifier.go` (`ip-adapter`/`ipadapter` → `AssetControlNet`), `types/types.go`, scanner index buckets
- **Symptom:** IP-Adapters are classified as ControlNets by name, and the rest of the code treats them as ControlNets. There is no separate IP-Adapter concept, no IP-Adapter-specific behavior, and no way to tell a control-signal ControlNet from an image-prompt IP-Adapter in the list.
- **Impact:** Users who download an IP-Adapter see it among ControlNets and have no idea how to use it differently.
- **Fix:** Introduce an IP-Adapter asset type (or at least a sub-classification) so the UI can say "this is an image-prompt adapter, attach it like a control net but feed it a reference image, not an edge map." The scanner already sees the name; the gap is downstream representation and usage.

### P1 — Add `--ref-image` (Flux Kontext) and photo-reference paths
- **Files:** `client.go` `txt2imgRequest` / `img2imgRequest` (no `ref_image` field), sd.cpp help has `--ref-image` as a repeatable flag
- **Symptom:** Flux Kontext-style multi-reference prompting exists in sd.cpp but is not wired into the Go request structs or the engine args.
- **Impact:** A modern Flux feature is unavailable even though the backend supports it.
- **Fix:** Add a `RefImages` (or `ReferenceImagePaths`) field to `ImageOptions`, forward it as `--ref-image` in `buildServerArgs` for the relevant families, and expose it in the UI as an image-reference/prompt 보강 flow. This is the start of real IP-style prompting without needing a separate "IP frontend" to exist first.

### P2 — PhotoMaker ID path exists in sd.cpp but is not wired
- **Files:** sd.cpp help (`--pm-id-images-dir`, `--pm-id-embed-path`), not in Go structs or args
- **Symptom:** PhotoMaker-style ID image input is available in sd.cpp but not exposed.
- **Fix:** Wire `--pm-id-images-dir` and `--pm-id-embed-path` if the app wants ID/face-reference generation. Treat this as part of the same IP/reference cluster as ref-image.

### P2 — No multi-controlnet weighting
- **Files:** `client.go` (single `ControlNetPath` + `ControlStrength`)
- **Symptom:** Only one ControlNet path per request. No per-condition weighting beyond a single strength.
- **Impact:** Multi-condition workflows (e.g. pose + depth, or IP-Adapter + canny) are not expressible.
- **Fix:** Extend the request to carry multiple control conditions with individual paths, images, and strengths. This is more involved than a single field, so park it after ref-image and IP-Adapter type separation are in place.

---

## Stage 3 — Model library & deps (P1/P2)

These make downloading and assembling a working model+companions set painless.

### P1 — HuggingFace / Civitai browsing is absent; only hardcoded companion URLs exist
- **Files:** `internal/image/scanner/civitai.go` (by-hash lookup only), `internal/image/deps/deps.go` (hardcoded companion catalog URLs), TODO in `internal/image/engine/sdcpp/discovery.go`
- **Symptom:** The app can name a UUID model from Civitai by hash, and it can offer a few hardcoded companion download URLs, but there is no model/library browser, no search, no "get this model" flow.
- **Impact:** A user who lands on a HuggingFace or Civitai page for a model must download everything by hand and hope the companions match. The deps catalog helps for *known* families, but there is no general discovery.
- **Fix, in order:**
  1. Keep the hardcoded companion catalog (it already works for Flux/Flux.2/Z-Image/Qwen-Image/Anima and prevents wrong-encoder mistakes) — this is good and should stay.
  2. Add HuggingFace repo metadata fetching for discovered models: read `pipeline_tag`, `library_name`, and the repo file list so the app can say what a model is and what companions it likely needs. This is the "TODO" in discovery.go.
  3. Add Civitai model/version browsing (not just by-hash) so a user can search, see triggers, see compatible families, and kick off downloads. The by-hash lookup is useful but narrow.
  4. Add automatic companion matching from model metadata: given a downloaded checkpoint, determine needed companions from modelcard/config/tensor layout and offer to fetch the right ones. This is the "make it just work" step that turns a download into a loadable model.

### P2 — LoRA type/family classification is too coarse
- **Files:** `scanner/classifier.go`, `scanner/gguf_probe.go`, `CompatibleLoRA()`
- **Symptom:** LoRAs are classified as "lora" with a guessed family; `CompatibleLoRA()` mostly checks exact family match or unknown. There is no reliable LoRA-type metadata (e.g. what a LoRA actually affects) surfaced to the user.
- **Impact:** Users can't easily tell whether a LoRA is relevant to their model family beyond a filename guess.
- **Fix:** Improve LoRA family detection where the metadata supports it, and expose whatever type signal exists. This is lower urgency than the companion-matching work becauseLoRA application already works technically.

### P2 — Text encoder / companion auto-fill is good but not discoverable
- **Files:** `app-image-model.go` (`autofillCompanionModels`, `missingCompanions`, `validateCompanionFamilies`), `internal/image/deps/`
- **Symptom:** Companion autofill and mismatch detection exist and are well done (the deps package explicitly prevents the Qwen3-0.6B-vs-4B and Flux-AE-vs-Qwen-VAE mistakes). The problem is visibility: a user may not know why a companion was chosen or why a load is blocked.
- **Fix:** Make the dependency report (`deps.Resolve`) visible in the UI at load time, with the mismatch reason and the download/source link when relevant. This turns a previously internal check into a user-facing explanation, which is a high-value usability gain for little code.

---

## Stage 4 — Polish & power features (P2/P3)

These expose advanced knobs already present in the code/sd.cpp and harden durability.

### P2 — Caching, SLG, SCM, tensor-type-rules, and per-component offload are all present in sd.cpp but not exposed
- **Files:** sd.cpp help flag set vs. `types/types.go` and `ModelLoadOptions` / request structs
- **Specifics:**
  - Cache: `CacheMode` + `CacheOption` exist; the per-mode knobs (threshold, start, end, decay, etc.) are not broken out in the UI.
  - SLG: `SlgScale`, `SkipLayerStart`, `SkipLayerEnd`, plus the high-noise variants, exist in the request structs but are not surfaced.
  - SCM: `--scm-mask`, `--scm-policy` exist in sd.cpp but are not wired.
  - Tensor-type-rules: `--tensor-type-rules` exists in sd.cpp but is not in `ModelLoadOptions`.
  - Per-component offload: `--control-net-cpu`, `--clip-on-cpu`, `--vae-on-cpu` exist in sd.cpp but are not in the load options.
- **Fix:** Decide which of these deserve a UI surface. A reasonable line:
  - **Surface:** cache mode picker + a text field for cache options (already partly there), SLG scale toggle for DiT models.
  - **Advanced-only / text input:** SCM mask, tensor-type-rules, per-component CPU offload — expose as advanced text/config fields rather than polished UI, because they are expert knobs.
- **Why this is Stage 4:** None of these block core usability; they are power-user and debugging surfaces.

### P2 — Video missing knobs that sd.cpp already supports
- **Files:** `video.go` `buildVideoArgs`
- **Specifics not wired:**
  - LCM/EDA extra sample args for video (lcm: `noise_clip_std`, `noise_scale_start`; ltx2: `max_shift`, `base_shift`, `stretch`, `terminal`; euler_ge: `gamma`)
  - Wan VACE `--vace-strength`
  - LTXAV `--embeddings-connectors`
  - Video denoising schedule (currently a single `VideoDenoise` value)
- **Fix:** Add the most broadly useful first: VACE strength and LCM acceleration args for video. Audio and multi-shot are lower priority given the larger ffmpeg/audio gap.

### P3 — LoRA preview, VAE verification, and Civitai rename coverage
- **LoRA preview/thumbnail:** low impact, nice polish.
- **VAE decode verification:** useful for debugging, not a normal-user blocker; the existing TAESD-vs-real-VAE split already prevents the worst outcome.
- **Civitai rename coverage:** the scanner only attempts Civitai lookup for generic-looking filenames; renamed models won't get a Civitai name. Improving this is polish, not a core gap.

---

## Quick severity map (for triage)

| Item | Severity | Stage |
|------|----------|-------|
| SD3.5 blocked by conservative version check | P0 | 0 |
| img2img/inpaint no live preview | P0 | 0/1 |
| No mask creation tools | P0 | 0/1 |
| No ControlNet preprocessor | P0 | 0/2 |
| No ffmpeg; MP4/GIF rejected | P1 | 0 |
| Unapplied sd.cpp upgrade vs SD3.5 block | P1 | 0 |
| Outpainting seam/iteration gaps | P1 | 1 |
| Highres fix hidden | P1 | 1 |
| IP-Adapter mixed with ControlNets, no type separation | P1 | 2 |
| `--ref-image` / photo-reference not wired | P1 | 2 |
| No HuggingFace/Civitai browsing or auto companion matching | P1 | 3 |
| Sampler/RNG, masking semantics under-promoted | P2 | 1 |
| Cache/SLG/SCM/tensor-type/offload knobs not exposed | P2 | 4 |
| Video LCM/VACE/connectors/denoise-schedule not wired | P2 | 4 |
| LoRA type classification coarse | P2 | 3 |
| Dependency report not visible to user | P2 | 3 |
| LoRA preview, VAE verify, Civitai rename coverage | P3 | 4 |

---

## What was not changed by this rewrite

The factual content from the original audit is preserved: supported families, LoRA/ControlNet/VAE/encoder support, video capabilities, dependency catalog, version state, and the missing-code list. This document reorders that material by severity and stage and adds an explicit usability-first justification for each grouping. It does not add new capabilities beyond what the existing code and sd.cpp/llama.cpp versions can support.

---

*End of prioritized guide.*
