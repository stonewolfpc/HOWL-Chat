# HOWL Chat — Working Issue & Expansion Document

**Design law:** the app owns every fact it can verify; the model narrates.
The app is meant to be smarter than the model it hosts — that is the whole product.

**North star:** a world generator and inspiration tool for creative people.
Double-click, load a model, start making something. Zero setup knowledge required.
Fun matters. Control matters. Beginners get defaults; power users get depth, never code.

**Hardware floor:** assume 8GB VRAM, potato PCs, tiny context windows, models that
are dumber than the app. The app must be smarter than the model wherever it can be.
Anything the app can decide automatically, it decides automatically.

---

## Ordering law

1. **Broken** — something that does not work at all, or corrupts/loses work.
2. **UI** — works but is wrong, inconsistent, or unreadable.
3. **Additions** — new capability. Only after 1 and 2 are closed and the app is
   operational end to end.
4. **Ideas** — parked. Not scheduled. Only revisited if they fit the app and the
   foundations are solid.

**Stop rule:** one stage per session. A stage closes only with evidence — a
screenshot, a log line, a saved file, a count. No stage closes because the code
"looks right." At the end of every stage: answer the checkpoint questions before
opening the next one.

---

## Documentation age

This document merges two earlier working notes. Because their contents should not be
assumed correct, every claim below is tagged with its current investigation state:

- **UNVERIFIED** — written down, not yet confirmed against the running app.
- **INVESTIGATING** — currently being checked this session.
- **EVIDENCE** — checked this session with the listed proof.
- **CLOSED** — fixed and the checkpoint was answered.
- **OPEN QUESTION** — still needs an answer from you or from a fresh test.

Nothing in this document is assumed finished. A line that looks like a fix is still
only a fix once it carries **EVIDENCE** or **CLOSED**.

---

## Tier 0 — Broken (release blockers)

### T0.1 · My messages don't appear until restart
**Symptom:** I send a message, it never shows in the UI. Restart the app and it's
there. The reply renders fine.
**What we know:** the data is saved correctly — restart proves it. So this is a
render/state desync on the live path, not a storage bug. Do **not** clear caches
on load to "fix" it; that risks deleting real chats.
**Test first, before touching code:** send one message, reload, count how many
copies of that text exist in the saved chat JSON.
· 2 copies → duplicate send path.
· 1 copy → pure render bug.
Then: switch to another chat and back without restarting. If the bubble appears,
the send handler simply never re-renders.
**Evidence:** the JSON count, plus a screenshot of the bubble appearing same-session.
**State:** UNVERIFIED.

### T0.2 · The model receives my message twice
**Suspected separate bug** from T0.1 — likely the message is appended to the
prompt twice even though it's saved once. Fixing T0.1 may not fix this. Test both.
**Evidence:** a prompt dump or server log showing the message once per turn.
**State:** UNVERIFIED.

### T0.3 · Image model will not load
**Symptom:** black textless popup flashes (~0.2s), vanishes. No model load is
attempted. Loading animation runs, but no VRAM or CPU movement. Also: there is no
way to choose the image model directory at all.
**Likely cause:** stored absolute model paths still point at the old drive after
the E: → C: move, and the failure has no error surface so it dies silently.
**Fix in three parts:** (a) validate stored paths at startup and fail loudly;
(b) add a model directory chooser; (c) give that dialog a real error state instead
of a black box.
**Evidence:** invalid path → readable message. Valid path → VRAM climbs within 5s.
**State:** UNVERIFIED.

### T0.4 · Image preview is broken
Long-standing. The generated image does not preview. Find the cause and repair it.
This is the last blocker before the image page can be judged on layout (T2.3).
**Evidence:** prompt → image appears in the preview panel, correct aspect, no crop.
**State:** UNVERIFIED.

**CHECKPOINT 0** — answer before Tier 1:
- Did the JSON show 1 copy or 2?
- Does the bubble render on chat-switch without restart?
- Which drive/folder holds your image models now?

---

## Tier 1 — Core mechanics (works, but wrong or forgotten)

### T1.1 · State is not owned properly
Changing tabs must **not** unload the model or drop the connection. Coming back
should find everything where you left it. Settings — samplers especially — must
persist visually across exit and reload, not just internally.
**Evidence:** load model → switch tabs → return → still loaded, still answering.
Change samplers → restart → sliders sit where you left them.
**State:** UNVERIFIED.

### T1.2 · VRAM arbitration between chat and image models
Agreed behaviour:
- Image model requested while a chat model is loaded → check whether both fit.
  If not, a confirm dialog: this will unload your chat model to make room.
- On confirm: unload chat, re-check freed VRAM, then load image.
- Leaving the image page does **not** unload or disconnect.
- Going back to chat and pressing Load mirrors it: check image model, confirm,
  unload, load. Same dialog, wording swapped.
**Evidence:** both directions exercised, with VRAM readings before and after.
**State:** UNVERIFIED.

### T1.3 · Think tags / reasoning handling
Match how SillyTavern and Kobold handle think modes: on, or off, or hidden, and
selectable intensity where the model supports it. Hidden means it must not count
against the context budget. A thought must never be fed back to the model on the
next turn — anything inside a think tag is discarded from history.
**Evidence:** hidden mode shows no tokens in the UI and shortens the measured
prompt; a two-turn test proves no think block reappears in turn two.
**State:** UNVERIFIED.

### T1.4 · Sampler defaults the user never has to touch
Defaults stay "suggested general," but vary by model family — they differ too much
to force one set. Supported models must get correct prompt formatting
automatically. No card in the UI that the backend ignores.
**Evidence:** load two different families, confirm each gets its own input/output
format with no user action.
**State:** UNVERIFIED.

### T1.5 · Header says the wrong name
"Howl Chat - Tron Edition" → "Howl Chat". One-line fix, bundle it with T1.1.
**State:** UNVERIFIED.

**CHECKPOINT 1:**
- Which tab switch currently breaks a loaded model, if any?
- Does the image model survive leaving its page today?
- Do sampler sliders survive a restart right now, or reset?

---

## Tier 2 — UI (works, looks wrong or reads badly)

### T2.1 · Theme unification
Unify dark and light palettes across every page to match the chat page. Image page
is already consistent — don't disturb it. Light theme is the real problem: header
text stays white in many places instead of adapting, which suggests an override
left behind by the earlier coder. Hunt the overrides rather than patching each
symptom. Also verify what a system-light-theme default should look like: readable,
low strain.
**Evidence:** side-by-side screenshots of each page in both themes, no white-on-white
or white-on-light text anywhere.
**State:** UNVERIFIED.

### T2.2 · Ghost elements and stale markup
Class of bug: markup left behind by deleted features, and duplicate element IDs
where one handler writes to an ID that no longer exists while the real element sits
there working. We've already hit one of these (`ig-negative-prompt` vs `ig-negative`).
Sweep every element lookup target against the actual DOM.
**Evidence:** zero write-to-missing-ID, zero duplicate IDs, no orphaned markup.
**State:** UNVERIFIED.

### T2.3 · Image page layout — target Invoke AI organisation
Rework toward Invoke AI's structure: preview prominent, controls grouped, history
organised. Reference it (look it up) but don't chase parity — take the parts that
make it powerful and organised, drop the rest. Depends on T0.4.
**State:** UNVERIFIED.

### T2.4 · Speaker identity in chat (thumbnail or initial)
**Gap:** names alone are too slow to parse in a hurry — the eye locks onto an
image instantly but has to read a name. This is a legibility item, not decoration.

Spec:
- Thumbnail of the character's card image beside the speaker name on every message.
- No image → enlarged first initial of the character's name, light text on a dark
  blue field.
- Fixed size (small enough not to crowd the bubble, large enough to recognise at a
  glance), consistent across themes — the dark blue must not fight light theme or
  vanish into the bubble.
- Serve a resized thumbnail, not the full card image, so a large portrait isn't
  loaded once per message.

Why it's early: this is a **prerequisite for multi-character legibility** (T3.5).
Two speakers distinguished only by name is exactly the failure mode it fixes.
Bundle with T2.1 so the fallback colour is decided once for both themes.

**Evidence:** screenshot of a two-character exchange, one with an image and one
without, both instantly distinguishable; thumbnail file size confirmed small.
**State:** UNVERIFIED.

### T2.5 · Character / world / scenario image sets
**Existing capability (regressed):** a card can hold several images (up to ~10)
with one marked as default. Worlds and scenarios have the same capability, with the
image shown in the right-hand panel. This worked at some point. Find out whether it
was scrapped, overwritten, or simply broke.

**Investigation order — cheapest explanation first:**
1. **Data layer.** Can a card still store a list of images plus a default, or only
   a single path? If storage lost the list, the feature went at the data layer and
   the UI was following it down. That's "scrapped."
2. **Leftover symbols.** Search for the old names: `cycle`, `slideshow`, `gallery`,
   `interval`, `images[`, `setInterval` in the character/world/scenario code. A
   surviving handler with no DOM target is the "overwritten" signature — the same
   pattern as the `ig-negative-prompt` ghost.
3. **Lifecycle.** If it's timer-driven, check whether anything clears the timer on
   tab switch.

**Evidence:** whichever layer broke, located and named; the current storage shape
shown; any orphaned handler flagged.
**State:** UNVERIFIED.

**CHECKPOINT 2:**
- Which pages look worst in light theme to you?
- Any specific Invoke AI screen you want to use as the reference?

---

## Tier 3 — Additions (only once the app is operational)

Nothing here starts until Tiers 0–2 are closed.

### T3.1 · Context handling — verify before optimising
**First question is not "how do we optimise," it's "was the context flush ever
implemented?"** Verify that before assuming it exists. If it exists, measure the
1–2s prefill cost. Then decide whether restructuring is worth it at all.
Candidate cheap win to test: `cache_prompt: true` plus `--cache-reuse N`, which
reuses shifted KV chunks instead of re-prefilling. Full model unloads are **not**
acceptable — load times are too long.
Research before implementing, so we don't optimise into something the models
tolerate badly.
**State:** UNVERIFIED.

### T3.2 · Load speed — learn from LM Studio
LM Studio loads a model in under a quarter of our time. Find out how and adapt
whatever is portable. This may be the single biggest quality-of-life win in the app.
**State:** UNVERIFIED.

### T3.3 · Scoped lorebook + memory system
Lore must be separable: world, scenario, character. Created as a lorebook, assigned
to an entity, and it appears in that entity's list at the top of its page. Character
lore fires **only** for that character.

Mechanics to implement: primary + secondary keys, word-boundary matching, scan
depth 2–3, recursion capped at 2 with a per-entry no-recursion flag, a fixed token
budget with priority ordering, and insertion order.
Precedence on a turn: character → scenario → world; character wins on conflict.
Dedupe identical bodies across scopes, prefer the narrower scope. Give lore its own
budget ceiling (~15–20% of context) so it can't starve the character card.
Memory: per character, tied to who said it, natural-language triggers, and a hard
guard that it can never blow the context window.
**AI-authored lore:** generate entries from a world/scenario/character description.
Keys get a specificity score (reject common single nouns, prefer proper nouns and
multi-word phrases), plus a frequency test against the existing chat log — a key
that already appears everywhere is flagged as "fires too often," with the count.

**Priority note within Tier 3:** this is the first addition worth doing, because it
is the only one that changes how the app holds onto what the user built. Lore scoping
comes before the rest of Tier 3. Memory inside it is part of the same system, not a
separate feature. This is also the first addition that should be checked for overlap
with T3.1 — if context handling is broken, lore scoping is only a nicer way to overflow
the same broken pipe.

**State:** UNVERIFIED.

### T3.4 · Per-field AI character sheets
Every field gets its own generate button. The model sees the fields already filled
and must make the new one line up: name, kind, the target field's purpose, plus
truncated versions of siblings as reference, and returns only that field.
Draft + accept / regenerate / lock per field. Grammar-constrained JSON, one
pipeline with per-kind schemas.
No preset voices. Exemplar dialogue is authored per character — exemplar lines beat
trait lists for keeping a voice from bleeding.
**State:** UNVERIFIED.

### T3.5 · Multi-character chat
2+ characters kept fully separate, each receiving the appropriate system prompt,
with context flushed properly between speakers. Root cause of voice bleed is flat
context, not sampling — needs structured transcript + per-character state.

**Priority note within Tier 3:** this is the addition that depends on T2.4 and on a
sane context/flush story from T3.1. It is not the first addition; it is the one that
stays blocked until the others in this tier are solid enough to expose cross-character
context correctly.

**State:** UNVERIFIED.

### T3.6 · Ternary models as first-class
Load ternary formats without new dependencies. Must just work when the user points
at the file.
**State:** UNVERIFIED.

### T3.7 · Audio and TTS end to end
Test the audio model capability and the TTS service through the UI. Needs headphones
present for verification.
**State:** UNVERIFIED.

### T3.8 · Audio generation
Newer small audio model, reportedly fits in 8GB for short clips. Short clips are
fine. If it works, play the generated audio inside the app.
**State:** UNVERIFIED.

### T3.9 · Diffusion version check
Confirm the current Stable Diffusion release; update to the latest patched version
if one exists.
**State:** UNVERIFIED.

### T3.10 · LLM writes the image prompt — research gated
Can a local LLM prompt the image generator well? Would it be worth the cost?
Proposed flow if yes: compose the prompt before swapping models; optionally keep
enough of the LLM resident to inspect the result (only if vision-capable) and refine.
Every part of this is a big if — cost, quality, speed. Research before building.
**State:** UNVERIFIED.

### T3.11 · Compatibility for everyone else
Auto context and layering currently assume this machine. Generalise as far as
possible and ship as beta, then fix from user reports. Goal: if it runs AI at all,
it should be usable.
**State:** UNVERIFIED.

### T3.12 · Whole-app optimisation pass
Separate pass, its own stage, after everything above. Look for waste everywhere.
**State:** UNVERIFIED.

**CHECKPOINT 3:**
- Confirm the flush was/wasn't implemented before I start T3.1.
- Which addition matters most to you after the app is solid?

---

## Tier 4 — Parked ideas (no stage)

- **Reactive 3D character.** Generated avatar that animates emotionally and moves
  its mouth while speaking. Enormous VRAM cost; likely incompatible with multi-character
  chat (too many bodies). Fallback would be a 2D still, but only if animation can be
  fully automated and prompt-driven — no Live2D hand-animation. Attractive and
  probably out of scope. Revisit only if the foundations are strong and resources allow.
  **State:** UNVERIFIED / parked.
- **Multi-character with avatars.** Sub-case of the above; same resource problem.
  **State:** UNVERIFIED / parked.

---

## Documented expansion concepts — parked until foundations are solid

These came from the expansion notes. They are listed here because the notes asked for
them to be visible, but they stay at the bottom of the priority order. Nothing here is
assumed working, planned, or ready to build.

### E.1 · Model-aware image anchoring for vision-language context
Proposed idea: feed the model a composite context window containing the world image,
the scenario image, and the character's default portrait alongside the user's prompt,
so the model "sees" setting, situation, and speaker instead of guessing from text.

**Constraints recorded in the notes, to preserve verbatim:**
- High-res images can consume roughly 1,000–1,500 visual tokens each after patching.
  Three high-res reference images per turn can inject ~3,000–4,500 visual tokens
  before any chat text. On small local models, that is enough to crowd the context
  budget and slow generation.
- Passing three images every turn can dilute attention when the user is only making
  a simple move.
- A plausible compromise in the notes: character portrait and world image are passed
  mainly at context initialisation or scene change; scenario/action images are passed
  only when a new environment or active event image appears; and reference images should
  be downscaled before they reach the multimodal encoder (the notes mention 512×512 or
  768×768 as working resolutions).

**Investigation state:** UNVERIFIED.
**Current verdict:** parked behind T3.3/T3.5 and the vision model situation. This is
not assumed to be free, correct, or even desirable yet.

### E.2 · Multi-image vision support as a capability question
The notes observe that many modern open-weights vision-language models in the ~3B–12B
tier (names mentioned include Qwen2.5-VL, Pixtral 12B, and LLaVA-NeXT variants) can
accept more than one image in a single turn, because their vision encoders convert
multiple images into visual patch tokens that the transformer treats like other tokens.

**Investigation state:** UNVERIFIED.
**Current verdict:** not a standalone feature yet. If this ever becomes relevant, it
belongs with the image anchoring idea above, not as a separate promise.

### E.3 · Standout local-chat references
The notes name a few desktop/local chat tools as reference points rather than
blueprints:

- **LM Studio** — plug-and-play desktop local AI. Highlighted for Hugging Face
  integration and a VRAM calculator that colour-codes fit/split/crash before download.
- **Open WebUI** — self-hosted web UI seen as a local ChatGPT-like interface. Highlighted
  for feature density: multi-user auth, document RAG, audio I/O, web search hooks, tool
  calling.
- **Jan (Jan.ai)** — offline-first, Apache 2.0, no telemetry, clean local folders.
- **AnythingLLM** — local document knowledge bases / RAG, with chunking, embedding, and
  retrieval treated as first-class.

**Investigation state:** UNVERIFIED.
**Current verdict:** reference only. None of these is assumed to be a correct model for
HOWL Chat. Their ideas may inform T3.x priorities; they do not become scope.

### E.4 · Common local-AI friction points to avoid
The notes list recurring complaints about local AI apps and a fix direction for each.
These are recorded as things to actively avoid, not as features to build.

1. **Silent failures and frozen UIs.**
   Symptom: backend binary crashes or hits an unhandled exception, but the UI stays
   stuck on "Generating..." or a loading dot with no feedback.
   Recorded fix direction: hard socket watchdogs, process-exit listeners, and an inline
   error drawer that surfaces backend `stderr` when a process dies.

2. **Hidden file paths and locked-in data.**
   Symptom: models, chat histories, or prompt presets buried in obscure nested
   AppData/Application Support folders under scrambled IDs.
   Recorded fix direction: root-relative or user-configurable workspace directories with
   clear file structures such as `/models`, `/prompts`, `/characters`.

3. **Context window silent truncation.**
   Symptom: chat history is silently truncated when the context fills, so the model
   loses identity or instructions without telling the user.
   Recorded fix direction: a visible, real-time context meter such as
   "4,102 / 8,192 Tokens Used," so the user sees how much space remains before sliding
   memory or flushing.

4. **Rigid / opaque system prompts.**
   Symptom: system prompts are hardcoded or a single format is forced onto different
   model families.
   Recorded fix direction: transparent, user-editable system prompt templates, accessible
   per model family in plain JSON/YAML.

5. **Fragile port and dependency dependencies.**
   Symptom: users must manage external Python environments, manual port bindings, or
   sidecar node dependencies.
   Recorded fix direction: native executable compilation, dynamic internal port selection,
   and zero-dependency process orchestration.

**Investigation state:** UNVERIFIED.
**Current verdict:** these should be checked against HOWL Chat as it stands. Some may
already be handled. Some may still be real. None is assumed closed.

### E.5 · "Fake strengths" to avoid when designing the app
The notes list patterns that look good on paper but can become weaknesses:

1. **Massive extension/plugin ecosystems.**
   Risk: fragile dependency chains, bloat, conflicting background scripts after updates.
   Recorded defense: build the core (memory, image hooks, character routing, context
   flushing) into the core product. Treat third-party additions as optional power-user
   extras, not required crutches.

2. **One-size-fits-all "universal" system prompts.**
   Risk: a prompt tuned for one model family can produce bad output or choke on tags for
   another.
   Recorded defense: model-family routing. Map presets dynamically to the model family and
   let the user customize and save per-family configurations.

3. **Endless sampler sliders on the front page.**
   Risk: decision paralysis; users tweak blindly, break coherence, and blame the engine.
   Recorded defense: preset profiles with advanced drawers. Hide hyper-parameters such as
   XTC, quadratic penalization, and min-P inside collapsible "Pro" drawers, while exposing
   intuitive macro-presets such as "Strict Rule-Following," "Atmospheric Prose," and
   "Unconstrained Storytelling."

4. **Unrestricted full-context re-injection.**
   Risk: passing everything on every turn — all images, full lore, entire chat history —
   saturates context, bloats token use, raises VRAM use, and dilutes attention.
   Recorded defense: smart state partitioning. Drop turn-zero visual payloads after initial
   scene priming, rely on summarized text anchors in the system context, and flush
   intermediate character buffers to prevent cross-contamination.

**Investigation state:** UNVERIFIED.
**Current verdict:** these are design guardrails. They are especially relevant if any
Tier 3 addition touches memory, context, or image injection.

### E.6 · Custom workspace layout for files
The notes propose an explicit, obvious default layout rather than burying data inside
obscure app folders:

- Chat logs and all text of any kind, including lore entries, live under the system
  **Documents** directory by default.
- Images live under **Pictures**.
- The proposed default image tree is something like `Pictures\Howl Chat\Uncategorized`,
  with per-board/owned folders under `Pictures\Howl Chat\<whatever the user named it>`.

**Investigation state:** UNVERIFIED.
**Current verdict:** this is an app-level directory/lifecycle decision. It may interact
with T0.3 (image model directory chooser), with export/save paths, and with any future
vaults or image boards. It is recorded here because it is a concrete, user-facing choice,
not because it is assumed to be already implemented.

### E.7 · A simple "curtain" vault for private boards
The notes describe a light-weight private space for select images/boards. The stated goal
is not strong security. It is friction: keep casual browsing from exposing a private board
instantly, especially to a curious child, while still letting a capable adult recover the
files if needed.

Recorded shape of the idea:
- The vault is a hidden workspace folder. The notes mention marking the folder hidden
  (`+H`) or system-hidden (`+S +H`) on Windows so it does not show up in normal File
  Explorer views. The notes also mention a custom extension, used so that OS media
  pickers and photo apps do not pick up the files automatically.
- Access inside HOWL Chat is gated by a PIN. The notes are explicit that this is not
  intended to be cryptographic armor. The app should make the limitation obvious in any
  tutorial or onboarding, and should provide a recovery path if the PIN is forgotten.
- The notes suggest recovery ideas rather than a single mandated mechanism: for example,
  a stored local recovery code/file, or an unusual launch-time shortcut, as ways to reset
  the PIN without trapping the user. The point is recovery without cloud services and
  without bricking the user's own files.
- The notes also mention a header/byte swap idea for the stored files so that opening a
  vault file directly in a normal image viewer looks broken, while HOWL Chat can still
  render it. This is presented in the notes as optional layered obfuscation, not encryption.

**Investigation state:** UNVERIFIED.
**Current verdict:** parked behind the app being solid and behind any real need for a
private board feature. This is not assumed wanted, wanted now, or wanted at all. If it
ever moves forward, it should be treated as the "curtain" the notes describe: a visual
and casual-browsing barrier, with obvious recovery and obvious limitations.

### E.8 · Image staging board structure
The notes mention an existing staging area for the image board, with multiple boards and
the ability to slide things through them. The notes also express a desire for saved images
to mirror that board structure inside the chosen image directory.

**Investigation state:** UNVERIFIED.
**Current verdict:** this is likely related to E.6 and possibly to image-page work, but it
is not assumed implemented. Before treating it as a feature, find out what the current
staging area actually is and where it currently saves things.

---

## Research queue (do before committing to any of it)

| Topic | Question | State |
|---|---|---|
| LM Studio load speed | How does it load so fast, and what's portable to us? | UNVERIFIED |
| Context flush | Was it ever implemented? What does it actually cost per turn? | UNVERIFIED |
| KV cache reuse | Does `--cache-reuse` help on this hardware, measurably? | UNVERIFIED |
| Think modes | How do SillyTavern/Kobold expose levels, and what do models accept? | UNVERIFIED |
| Ternary formats | Which quantiser/runtime path, with no new dependencies? | UNVERIFIED |
| Invoke AI | Which layout decisions are worth copying? | UNVERIFIED |
| Stable Diffusion | Current release, and what updating costs us. | UNVERIFIED |
| LLM → image prompt | Is the quality worth the swap cost? | UNVERIFIED |
| Audio model | Exact model, VRAM, clip length on 8GB. | UNVERIFIED |
| Vision model anchoring | Is multi-image context worth the token cost, and under what rules? | UNVERIFIED |
| Local-chat references | Which ideas from LM Studio / Open WebUI / Jan / AnythingLLM transfer? | UNVERIFIED |
| Directory layout | What should the default Documents/Pictures shape be, and where do image models live now? | UNVERIFIED |

---

## Open questions for you

1. Which chat-switch or tab-switch currently loses a model? I need the actual
   reproduction, not the general rule.
2. Where do your image models live now — the new drive path?
3. Light theme: is there a specific page that's the worst offender?
4. Beta scope: ship beta right after Tier 0–2, or after T3.5 (multi-character)?
5. Is the "curtain" vault something you actually want in the app, or is it still just
   an idea you are keeping in case it becomes useful?

---

## Decision log

Append here as stages close. One line each: date, stage, evidence, result.

## F2 Requirements Ledger — 2026-11-09

This is the working capture for the F2 storage, content-library, and model-intelligence
requirements. It lives here because this is the project's existing requirements and
deferred-work document; no second requirements ledger is being maintained. The F2
engineering report is `docs/review/ENGINEERING-F2-STORAGE-MODEL-INTELLIGENCE-2026-11-09.md`.

| ID | Requirement (user intent) | Status | Dependencies / next gate |
|---|---|---|---|
| F2-A1 | Multiple user-controlled directories per content category; do not force users to move existing collections | APPROVED, backend and settings UI implemented | Content browsers still need to consume all registered roots for every category |
| F2-A2 | Explicit read-only, read/write, and save-destination roles; never modify a read-only library | APPROVED, registry roles and settings UI implemented | Existing save/import paths still call their old stores; route them through the library writer before calling this complete |
| F2-A3 | External library visibility without forced migration; show source, editability, and save destination | APPROVED, character-card discovery and source metadata implemented | Character grid/edit/delete need explicit external-item behavior; read-only cards must not be edited as local copies |
| F2-A4 | SillyTavern character-library compatibility | APPROVED, JSON and PNG `chara` discovery implemented read-only by default | Import/export adapters and editable SillyTavern operations remain future work |
| F2-A5 | Documents/Pictures/Videos defaults through Windows known folders | APPROVED, implemented and unit-tested on this host | Runtime confirmation on redirected/OneDrive known folders still required |
| F2-A6 | Separate future comprehensive cleanup operation; Reset Settings is settings-only | APPROVED and recorded; settings reset now covers unified and image stores | Cleanup UI must be separately named, previewed, and explicitly confirmed; not part of Reset App to Defaults |
| F2-B1 | Unified model directory registry across multiple roots | APPROVED, existing image/model directory registry retained and made availability-aware | A single model-root manager UI and cross-kind directory listing remain to be wired |
| F2-B2 | Unified model identification across chat, image, audio, and future TTS | APPROVED, shared `internal/modelintel` profile/cache layer added | Audio/model execution adapters, tokenizer/template inspection, and full workspace integration remain future work |
| F2-B3 | Cached identification with lightweight revalidation; no huge-file hashing in ordinary discovery | APPROVED, size + modtime cache and reuse implemented | Runtime/packaged-app behavior and cache invalidation for ancillary files remain unverified |
| F2-B4 | Evidence-based provenance and template detection; do not turn marketing names into facts | APPROVED, Verified/Source-derived/Inferred/Unknown/Conflicting model added | Tokenizer, chat-template, config, and trusted online source inspection are not implemented |
| F2-B5 | Separate modality-specific creative workspaces; shared catalog does not merge chat and image interfaces | APPROVED, capability classification keeps chat/image/TTS distinct | UI selectors need explicit capability filtering against the shared catalog |
| F2-B6 | Shared backend model intelligence and resource coordination | APPROVED as architectural direction; shared metadata cache exists | Do not replace existing modality engines; resource-manager consolidation remains a later phase |
| F2-B7 | Music/audio expansion deferred until those subsystems are ready; TTS contract frozen | DEFERRED / FROZEN | No music feature, no new TTS feature, no TTS contract change in F2 |

Additional decisions recorded during F2:

- `InternalConfigDir()` remains under the project root's `config` directory for compatibility;
  changing it is a migration decision, not a silent F2 side effect.
- Registry entries for missing model files remain persisted and are marked unavailable;
  unavailable entries are not offered to pickers.
- Removing a content-library registration removes only the association; the directory and
  all files are retained.
- No destructive migration, content move, model move, telemetry, cloud inference, or
  background network activity was introduced.

## F2-B Deferred Future Requirements — 2026-11-09

These are recorded for later planning, not implemented in F2-B. Each entry names the
subsystem/roadmap location that should own it; do not create a second competing ledger.

| ID | Future requirement | Eventual owner / roadmap slot | Dependencies and gate |
|---|---|---|---|
| FUT-C1 | Up to five chat message cues shown in the message composer and persisted with the chat | Chat UI / prompt assembly (post-F2 chat behavior phase) | Chat schema versioning, cue-to-prompt mapping, and user control; do not alter TTS contract |
| FUT-C2 | Image cue support in chat, with final supported count established by research | Chat UI / multimodal prompt assembly | Image attach/read pipeline, model capability evidence, bounded image count research |
| FUT-C3 | llama.cpp feature exposure: Flash Attention, GPU offload, basic multi-GPU, hardware-specific optimizations | Chat model runtime settings / llama.cpp adapter | Verified backend flags per model/runtime; hardware probe; safe fallback when unsupported |
| FUT-C4 | Architecture-, family-, and model-era-specific options, capability-gated | Model intelligence + chat model settings | Reliable architecture/family evidence; per-capability allowlist; never filename-only certainty |
| FUT-C5 | Comprehensive sampling controls, including XTC and rediscovery research for the three unidentified samplers | Chat sampling settings | Runtime sampler support matrix; named/versioned controls; research provenance recorded |
| FUT-ST1 | SillyTavern character-card import from PNG metadata and JSON formats | Character import adapter / content library | Read-only discovery already exists; add explicit preview and copy-on-import, never in-place rewrite |
| FUT-ST2 | Intelligent mapping of supplied card fields into Howl fields | Character import adapter | Field-by-field mapping table, conflict policy, and user-visible provenance |
| FUT-ST3 | Preserve and intelligently parse information dumped into Description | Character import adapter / character form | Conservative parser with original text retained; explicit user review before applying |
| FUT-ST4 | Character LoRA support | Character model configuration / llama.cpp adapter | LoRA path discovery, compatibility validation, runtime loading, and resource accounting |
| FUT-ST5 | Lorebook import and compatible configuration adoption | Lorebook import adapter / content library | Import preview, entry-level validation, conflict handling, and no silent overwrite |
| FUT-CH1 | Chat-scoped Author's Note with popup, floating, hide, and side-panel behavior | Chat composer / chat settings | Per-chat persistence, prompt insertion order, and visible enable/disable state |
| FUT-CH2 | Single-character chats retain configured context capacity before rollover/memory strategies engage | Context manager / memory continuity | Token accounting per model, reserved response budget, and observable rollover diagnostics |
| FUT-CH3 | Multi-character context isolation and continuity, with flushing between turns or a demonstrably stronger alternative | Multi-character chat engine / memory continuity | Per-character transcript boundaries, shared-scene state, testable isolation, and no cross-character leakage |
| FUT-IMG1 | Hardware-aware RAM/VRAM placement for image generation | Image engine resource manager | Hardware probe, backend capability matrix, safe fallback, and measured memory behavior |
| FUT-IMG2 | Safe CPU/system-memory offloading where supported | Image engine resource manager | Backend-supported offload only; guardrails for unstable or insufficient memory |
| FUT-IMG3 | Maximum practical OOM prevention | Image engine + shared resource coordination | Preflight estimates, bounded retries, cancellation, and error reporting that preserves existing output |

These requirements remain deferred. They must not be claimed as implemented by F2-B
merely because the shared model/content foundations now expose extension points.

## F2-C Content Ownership Decisions — 2026-11-09

This section records architecture decisions only. No cutover, migration, or ownership
change was implemented by F2-C. The authoritative evidence and per-category design are
in `docs/review/ARCHITECTURE-F2C-CONTENT-OWNERSHIP-CUTOVER-2026-11-09.md`.

| ID | Decision | Status | Next gate |
|---|---|---|---|
| F2-C1 | One authoritative owner per content item; content libraries are address/permission registrations, not a second copy. | DESIGN RECORDED | Category adapters and read-back paths |
| F2-C2 | F2-B automatic JSON copies are untrusted export/compatibility snapshots, not canonical content, and must not be read back automatically. | DESIGN RECORDED | Explicit export semantics before any migration |
| F2-C3 | Browser localStorage/IndexedDB becomes a rebuildable cache after cutover, never the canonical owner. | DESIGN RECORDED | Per-category migration and rollback evidence |
| F2-C4 | Chats retain their multi-file folder format; register the existing AppData chat root in place before considering a move. | DESIGN RECORDED | Root-aware chat persistence and active-pointer migration |
| F2-C5 | Images and videos are binary-plus-sidecar generated media with stable identity and atomic commit. | DESIGN RECORDED | Media catalog and video sidecar implementation |
| F2-C6 | `Save` is create-only; replacing an existing item requires `Write` or an explicit collision choice. | DESIGN RECORDED | Registry role semantics and tests |
| F2-C7 | No automatic migration or background synchronization; migration is explicit, verified, journaled, and reversible. | DESIGN RECORDED | Preflight, journal, rollback, and user reporting |
| F2-C8 | Disconnected roots retain registrations and last-known index entries marked unavailable. | DESIGN RECORDED | Availability-aware index cache |

F2-B future requirements remain deferred and assigned to their existing roadmap slots.
No F2-C decision promotes any of them to implemented.


## F3-A Ownership Safety Decisions — 2026-11-09

F3-A established truthful save/export behavior, collision-safe library writes, bounded
gallery access, unavailable-root inventory, and scoped reset/wipe behavior. No content
category was migrated. Evidence and test results are in
`docs/review/ENGINEERING-F3A-OWNERSHIP-SAFETY-MIRROR-SEMANTICS-2026-11-09.md`.

| ID | Decision | Status | Next gate |
|---|---|---|---|
| F3-A1 | Automatic F2-B library mirror writes are removed from character, world, scenario, and lorebook saves. Legacy JSON files are left untouched and are not read back. | IMPLEMENTED | Optional explicit export UI |
| F3-A2 | `ExportContentItem` is the only content-library write binding; it returns a structured `OperationResult` and writes a non-authoritative export envelope. | IMPLEMENTED | Expose an explicit export action if desired |
| F3-A3 | Export identities use a stable hash, safe segments, and collision detection. Existing files are never replaced by a normal export. | IMPLEMENTED | Per-category migration adapters |
| F3-A4 | Replacement requires `RoleWrite` plus an explicit `replace=true`; `RoleSave` alone is create-only. | IMPLEMENTED | Review role UX before migration |
| F3-A5 | Registered roots retain a persisted last-known character inventory; unavailable roots show stale items as unavailable, never as deleted or live. | IMPLEMENTED | Add live scanners for other categories if needed |
| F3-A6 | Gallery read/list/delete operations are contained by the gallery root, including symlink/reparse-point checks; atomic image+sidecar commit is used. | IMPLEMENTED | Media catalog and video sidecar design |
| F3-A7 | Wipe uses the configured chat root and gallery root, preflights both, preserves external roots, reports failures, and is not a settings reset. | IMPLEMENTED | Fully previewed comprehensive cleanup contract |
| F3-A8 | Deletion/rename propagation and full lifecycle synchronization are deferred until authoritative ownership moves. | DEFERRED | Lorebook cutover (F3-B) with rollback evidence |

F2-B future requirements and F2-C ownership decisions remain in force. F3-A does not
promote any deferred feature to implemented.

## F3-A-R1 Safety Repairs — 2026-11-09

F3-A-R1 closes the concrete F3-A review defects without starting F3-B. It
rejects unsafe replacement targets, makes wipe staging/rollback and state reset
explicit, and applies gallery containment to listed image candidates. Evidence is
in `docs/review/ENGINEERING-F3A-R1-SAFETY-REPAIRS-2026-11-09.md`.

| ID | Decision | Status | Next gate |
|---|---|---|---|
| F3-A-R1-1 | `WriteFile` uses `Lstat` and permits replacement only for an existing regular file. Directories, symlinks, and other non-regular targets return `ErrUnsafeTarget` and remain untouched. | IMPLEMENTED | F3-A-R1 acceptance |
| F3-A-R1-2 | Replacement backups use unique sibling paths; commit failure restores the original when possible and reports the recoverable backup path when restore fails. | IMPLEMENTED | F3-A-R1 acceptance |
| F3-A-R1-3 | Export directory identity is the content-library identity hash, not the display title; display names live in the export envelope. Empty identity with nil data is rejected. | IMPLEMENTED | F3-A-R1 acceptance |
| F3-A-R1-4 | Wipe stages both application-owned roots before resetting state, restores the first if the second stage fails, resets `LastChatID`, and clears the inventory cache. | IMPLEMENTED | F3-A-R1 acceptance |
| F3-A-R1-5 | Wipe reports recoverable staged data if cleanup fails instead of reporting a clean wipe. External roots remain untouched. | IMPLEMENTED | F3-A-R1 acceptance |
| F3-A-R1-6 | Gallery image candidates use the same canonical containment as reads; outside-resolving image or sidecar symlinks are not listed. | IMPLEMENTED | F3-A-R1 acceptance |
| F3-A-R1-7 | Export remains backend-only, character inventory remains the only live scanner, video sidecars remain incomplete, and delete/rename propagation remains deferred. | LIMITATION RECORDED | F3-B readiness review |
| F3-A-R1-8 | The final full test suite passed; the earlier audio failure did not reproduce in 10 isolated runs or the final suite. No audio code was changed and the flake remains unresolved. | OBSERVED / UNRESOLVED FLAKE | Monitor only; no audio change authorized |

F3-B remains unauthorized pending acceptance of this repair.

## Deferred Feature Requirements — 2026-11-09

**Governing product principle:** **Immersion is the goal. Continuity is the capability that supports it.**

Howl should make long sessions across characters manageable, preserving as much established
information and state as practicable while avoiding interruptions that detract from the
experience.

These are deferred requirements. They are documented for future design and sequencing only.
No item in this section is implemented, authorized, or complete. F3-A-R1 acceptance remains
the immediate gate, and F3-B ownership remains the next authorized sequence. These additions
must not reorder that sequence.

### Protected existing behavior — REG-1

The existing World, Scenario, Character, and Lorebook relationship system currently works and
is a regression constraint for every future ownership, storage, context, memory, or chat-UI
change. It must not be redesigned, simplified, removed, or silently replaced.

- Worlds are independently usable and hot-swappable.
- Worlds can have multiple attached scenarios, characters, and lorebooks.
- Scenarios and characters can be attached to multiple parents and remain independently usable.
- Scenarios can operate without worlds.
- Characters can operate without worlds or scenarios.
- Lorebooks can be saved with attachments or separately and reused.
- Attachments are reusable relationships, not exclusive ownership or physical nesting.
- Attached assets remain visible in their designated UI sections with clear `Attached to —`
  relationship labels, even when not active.
- Selecting a world activates its applicable attached content.
- Character lore supports private scope and public/shared scope.
- World and scenario lore are public/shared.
- Character-private lore is available only to that character; character-public lore can be
  triggered and seen by the active group.
- Context assembly must respect active relationships, scope, and context limits.

Status: **PROTECTED — REGRESSION CONSTRAINT, NOT A FEATURE REQUEST**

### Deferred feature register

| ID | Requirement | Roadmap placement | Dependencies | Status |
|---|---|---|---|---|
| FUT-CHUI1 | Chat bubble halo, spacing, and visual polish | Post-F3-B — Chat UI and appearance | Existing message rendering; preserved neon/TRON theme; readable contrast | DEFERRED — DOCUMENTED ONLY |
| FUT-CHUI2 | User-section appearance controls with coordinated colors and per-character independence | Post-F3-B — Chat UI and appearance | Per-character presentation state; theme defaults; multi-character identity model | DEFERRED — DOCUMENTED ONLY |
| FUT-CHTOK1 | Independent context-usage and message-token visibility controls | Post-F3-B — Chat UI and context visibility | Active model/runtime tokenizer evidence; context accounting; no false precision | DEFERRED — DOCUMENTED ONLY |
| FUT-CHID1 | Character thumbnail and fallback identity presentation in chat | Post-F3-B — Chat UI, character presentation, and multi-character identity | Character image storage; per-character color selection; group distinguishability | DEFERRED — DOCUMENTED ONLY |
| FUT-CHAIU1 | Optional AI-as-user role/persona capability | Post-F3-B — Advanced AI-as-user design | Separate feasibility study; explicit controls; multi-character isolation and context flushing | DEFERRED — DOCUMENTED ONLY |
| FUT-MEM1 | Application-side structured memory, inventory, and continuity | Post-F3-B — Application-side memory and inventory | Structured state schema; scope-aware retrieval; per-character and per-session isolation; context limits | DEFERRED — DOCUMENTED ONLY |
| FUT-LORE1 | Automated lore and memory intelligence overhaul | Post-F3-B — Lorebook and natural-language trigger intelligence | Durable-fact extraction; manual-lore conflict checks; review/correction UX; long-session reliability | DEFERRED — DOCUMENTED ONLY |
| FUT-PRINCIPLE1 | Governing product principle for immersion and continuity | No implementation slot — design principle | Applies across chat, memory, lore, and multi-character work | RECORDED — NOT A FEATURE |

### FUT-CHUI1 — Chat bubble appearance and polish

- Restore the missing cyan halo around the user message bubble in the default UI.
- Add a small, consistent gap between adjacent message bubbles so their borders no longer
  touch or visually overlap.
- Preserve the existing softened neon/TRON design language and its readability.

### FUT-CHUI2 — User appearance controls and per-character presentation

- Provide a user-section appearance control with approximately 20 coordinated color choices.
- A selected color should affect the message border, halo, and font color.
- Support independent per-character appearance settings in multi-character chat.
- Preserve sensible defaults and existing theme behavior.

### FUT-CHTOK1 — Token visibility controls

Provide independently configurable display options for:

- Current context usage as a percentage of the available context.
- Token count for the user's newly entered message.
- Token count for a specific generated response.

Ensure these values are meaningful for the active model/runtime and do not imply greater
precision than the tokenizer information supports.

### FUT-CHID1 — Character identity in chat

- Display the character's image thumbnail beside their name/messages.
- When no image exists, display a circular fallback avatar.
- The fallback uses the first initial of the character's name in white.
- The circle uses that character's selected border color as its background.
- Ensure multiple characters remain visually distinguishable in group conversations.

### FUT-CHAIU1 — Potential AI-as-user feature

Record as a major optional capability requiring separate design and feasibility review:

- Allow AI to take over the user role/persona in chat.
- Define explicit activation, deactivation, and user-control boundaries.
- Investigate interaction with multi-character context isolation and context flushing.
- Prevent personality, knowledge, and voice bleed between participants.
- Do not assume implementation details before architecture review.

### FUT-MEM1 — Application-side memory and continuity system

Record as a major Howl capability, not merely a minor chat enhancement.

The application must maintain useful continuity across long-running, multi-character sessions
without relying on the model to remember reliably. Examples include:

- Inventory and item ownership.
- Item acquisition, transfer, consumption, and use.
- State changes such as a health potion being used several turns earlier.
- Clothing and outfit continuity.
- Character state changes when clothing or equipment changes.
- Broader persistent narrative and session facts.

The system should maintain structured state and supply relevant information to the active
conversation when appropriate. Not every stored fact belongs in every prompt; context selection
must remain scope-aware and efficient.

### FUT-LORE1 — Automated lore and memory intelligence

Record a substantial overhaul and enhancement of the automated lore-entry/memory system:

- Improve natural-language trigger detection and retrieval.
- Improve extraction of durable facts and state changes from conversations.
- Support manual review, editing, correction, and override.
- Check automatically generated memories against relevant manual lorebooks.
- Detect likely matches, conflicts, and duplicate information.
- Flag uncertain or conflicting memory candidates for user review.
- Do not silently overwrite established manual lore.
- Avoid intrusive review prompts during roleplay.
- Keep automated behavior helpful and unobtrusive.
- Improve reliability across long conversations and multiple characters.

### Roadmap classification

These are deferred requirements, not permission to implement.

- **Chat UI and appearance:** FUT-CHUI1, FUT-CHUI2, FUT-CHTOK1, FUT-CHID1.
- **Multi-character chat and identity:** FUT-CHID1, FUT-CHAIU1, FUT-MEM1.
- **Application-side memory and inventory:** FUT-MEM1.
- **Lorebook and natural-language trigger intelligence:** FUT-LORE1.
- **Advanced AI-as-user design:** FUT-CHAIU1, after an explicit feasibility review.
- **Protected current relationship behavior:** REG-1, applying to every future stage.
- **Shared hardware/resource foundation:** FUT-HW1 in *Deferred Hardware Intelligence and Resource Orchestration — 2026-11-09*; staged after F3-B, with no change to the existing order.
- **Prompt and model-format intelligence:** FUT-PROMPT1 in *Deferred Prompt and Model-Format Intelligence — 2026-11-09*; alongside the unified model-intelligence and runtime-capability work, with no change to the existing order.
- **Hardware monitor (presentation only):** FUT-HW2, dependent on FUT-HW1 telemetry; placement, appearance, and metrics undecided.

Dependencies include the active World/Scenario/Character/Lorebook relationship model, context
assembly and limits, per-character scope isolation, tokenizer/runtime capability evidence, and
user-review UX for uncertain memory candidates. No future ownership or storage change may
weaken REG-1.

## Deferred Hardware Intelligence and Resource Orchestration — 2026-11-09

**Requirement ID: `FUT-HW1`** — unified, application-wide, hardware-aware model resource
orchestration.

**Status: DEFERRED — DOCUMENTED ONLY.** Nothing in this section is implemented, investigated
against source, authorized, or complete. It records architecture intent, dependencies, and
acceptance criteria so the requirement can be surfaced when the roadmap reaches the shared
model/resource-management foundation. It does not authorize implementation, does not change
model loading, and does not alter the existing engineering sequence: **F3-A-R1 acceptance
remains the immediate gate, F3-B storage ownership remains next, and `FUT-HW1` must not be
started during either.**

### Why this requirement exists

Howl must maximise usable performance while keeping the operating system, other applications,
and Howl itself stable across a broad range of hardware configurations.

This is not a VRAM monitor and it is not a model-loading convenience feature. A monitor reads
numbers; a convenience loader opens one file. `FUT-HW1` must coordinate all of: model
identification, resource estimation, admission, GPU/RAM/CPU placement, live monitoring,
workload scheduling, unloading, process cleanup, and verified transitions.

**Recorded user concern (unresolved):** loading a large model can currently fail when the
model does not fit entirely in VRAM, and the existing resource-aware layering system is **not
demonstrated** to work reliably. The presence of layering controls or settings is not evidence
that layering works. No entry in this ledger may describe current GPU/CPU layering or
offloading as functional without measured evidence.

### Backend priority and preservation of existing work

- Hardware-aware resource management and reliable model layering/offloading are **foundational
  backend requirements** for Howl.
- Some related systems and capabilities already exist in the application. Their presence must
  not be treated as proof of functional or reliable resource orchestration.
- Future engineering must inspect and preserve useful existing work, identify actual defects and
  missing links, and improve it rather than automatically replacing it.
- Large models can currently fail to load when they cannot fit entirely in VRAM. This is a known
  product concern that requires resolution, not documentation.
- The objective is reliable, hardware-aware operation across supported configurations, not
  merely exposing more settings.
- None of this authorises rewriting existing loading or process-management code; see §9 and
  AC-HW1-9.

### 1 · Intelligent model admission and layering

When a user selects a model and initiates work, Howl should:

- Identify the model architecture, format, quantization, runtime, and applicable loading
  capabilities.
- Estimate model weights, runtime overhead, context/KV-cache requirements where applicable,
  and expected working memory.
- Inspect available GPU VRAM, system RAM, CPU capabilities, and other relevant resources.
- Account for current consumption by other applications and by Howl's own processes.
- Calculate a safe resource plan **before** loading.
- Preserve a configurable minimum free-VRAM safety reserve.
- Prefer retaining as much useful work on the GPU as is safely possible, for performance.
- Automatically layer or offload supported components into system RAM/CPU when the complete
  model cannot safely fit in VRAM.
- Avoid treating an unsupported offload configuration as though it were supported.
- Fail clearly and safely when no viable loading plan exists.

The safety reserve must be **dynamically calculated** from actual available resources and
current conditions. A fixed percentage alone is insufficient.

Default policy favours stability and successful operation over aggressive VRAM utilisation.
Power users must be able to override caps and resource preferences, with clear warnings and
explicit control (see §7).

### 2 · Process-aware live protection

Resource safety cannot end when loading succeeds. Howl must monitor actual resource use during
loading, inference, generation, and other model-backed work.

- Detect VRAM and RAM pressure as it develops.
- Distinguish expected usage from actual process/system usage where possible.
- Account for external applications that consume resources after a model has loaded.
- Detect when continuing, or admitting additional work, risks OOM or system instability.
- Stop admitting new work when safety thresholds are reached.
- Notify the user clearly.
- Coordinate active operations and pending work rather than abruptly terminating processes
  while they are executing.
- Safely unload or offload resources when necessary.
- **Verify that resource release actually occurred.**

A request to unload a model is not proof that its weights have left VRAM. Backend processes,
child processes, memory mappings, and runtime ownership must all be considered.

### 3 · Lazy loading and seamless user workflow

The intended default workflow is:

1. User selects a Chat model.
2. User sends a message.
3. Howl performs admission and loads the selected model automatically if necessary.
4. Generation begins only after loading and state verification succeed.

There must not be a mandatory separate Load action in the normal workflow. The same principle
applies to Image, Audio, TTS, Vision, and future model-backed subsystems.

When switching workloads:

1. Determine whether the currently loaded models can safely coexist.
2. Retain multiple models when resources permit.
3. If resources are insufficient, coordinate the necessary unload.
4. If an operation is active, notify the user and queue the new work order.
5. Reach a safe stopping point for the active operation according to its cancellation/completion
   policy.
6. Unload the previous model.
7. Release and clean up its owned child processes and backend resources.
8. Verify that the required resources have actually been released.
9. Calculate/admit the next model's loading plan.
10. Load the requested model.
11. Verify the new model and resource state.
12. Execute the queued work.

The governing transition sequence is:

**Unload → Verify → Load → Verify → Work**

The next workload must not start merely because an unload or load request returned without an
error.

### 4 · Cross-subsystem coordination

`FUT-HW1` is an application-wide shared resource-management capability, not a per-subsystem
feature. At minimum it must plan for:

- Chat/LLM inference.
- Image generation.
- Audio generation and processing.
- TTS.
- Vision and multimodal models.
- Future model-backed capabilities.

Subsystems may retain separate engines and adapters, but they must participate in common
resource admission, ownership, scheduling, safety, and lifecycle coordination. No subsystem may
independently consume GPU resources without respecting the shared resource manager.

This requirement must be **reconciled with the existing process-ownership and
shutdown-lifecycle foundation** rather than creating a second competing manager. Two managers
with overlapping authority is a regression, not an implementation.

### 5 · Hardware compatibility and research

Aim for the broadest practical hardware support. Investigate and capability-gate:

- NVIDIA GPUs.
- AMD GPUs.
- Intel discrete GPUs.
- Integrated GPUs.
- Multiple GPUs.
- CPU-only execution.
- Different system RAM capacities and memory configurations.
- Different GPU memory capacities.
- Runtime-specific GPU offloading and CPU/RAM layering capabilities.

Each relevant backend and hardware family needs its own research. Some configurations may
require substantial research, implementation, or beta testing.

Maintain a compatibility/evidence matrix that keeps these states distinct:

| State | Meaning | Permitted wording |
|---|---|---|
| Verified | Measured on real hardware with recorded evidence | "supported" |
| Documentation-supported | Backend docs or source confirm the capability; not measured here | "supported by backend documentation, unverified here" |
| Implemented, unverified | Code path exists; no measurement | "implemented, unverified" |
| Experimental | Partially works with known instability or gaps | "experimental" |
| Unsupported | Backend cannot do it | "unsupported" |
| Unknown | Not investigated | "unknown" |

Universal compatibility must not be claimed without evidence.

### 6 · Performance and safety objectives

The system must balance:

- Maximum practical performance.
- Low latency.
- Safe VRAM reserves.
- Safe system RAM usage.
- CPU availability.
- Stable model loading.
- Reliable process cleanup.
- Prevention of OOM and application failure.
- Responsiveness while operations are active.

Stability must not be sacrificed simply to maximise GPU utilisation.

### 7 · Power-user controls

Transparent overrides for users who understand their hardware and accept additional risk.
Candidate controls:

- VRAM reserve.
- Resource caps.
- GPU/RAM placement preferences.
- Offloading preferences.
- Co-residency policy.
- Workload priority.
- Automatic unloading behaviour.

Defaults must remain safe. Overrides must not silently disable unrelated lifecycle or
process-safety protections; every override must state what it can break.

### 8 · Roadmap placement and dependencies

`FUT-HW1` is a **major deferred shared-foundation requirement**. It belongs in the
hardware/runtime/model-management architecture stage, before finalising the model-loading and
cross-subsystem integration work.

| Dependency | What `FUT-HW1` needs from it |
|---|---|
| Existing process ownership and shutdown lifecycle (F1) | Single authority for owned processes and cleanup; no second manager |
| Unified model intelligence (F2) | Reliable architecture, family, quantization, and runtime identification |
| Model-family and runtime capability detection | Per-capability evidence used for admission and capability gating |
| Chat and Image model managers | Placement, load/unload, and state reporting through one coordinator |
| Audio/TTS/Vision runtime integration | Participation in shared admission and scheduling |
| Application-wide concurrency and resource scheduling | One queue and one admission order across subsystems |
| Hardware monitoring and user-facing diagnostics | Live VRAM/RAM/CPU/pressure readings the user can inspect |

Related existing ledger entries that `FUT-HW1` must absorb rather than duplicate or contradict:
T1.2 (VRAM arbitration between chat and image models), T3.2 (load speed), T3.11 (compatibility
for everyone else), `FUT-C3` (llama.cpp feature exposure), `FUT-IMG1`, `FUT-IMG2`, `FUT-IMG3`
(image RAM/VRAM placement and OOM prevention), and `FUT-ST4` (LoRA resource accounting).

`FUT-HW1` must not be implemented during F3-A-R1 or F3-B storage-ownership work. The existing
engineering sequence is unchanged; this requirement is surfaced when the roadmap reaches the
shared model/resource-management foundation.

### 9 · Preservation and restrictions

- Do not replace existing working model functionality without evidence and authorization.
- Do not make TTS-specific implementation changes under this documentation task.
- Do not change the current roadmap order.
- Do not create a second authoritative requirements document.
- Do not claim existing layering or offloading works merely because controls or settings exist.

### 10 · Acceptance criteria

None of these are satisfied today. They define the evidence a future `FUT-HW1` stage must
produce; a compilation, a settings screen, or a successful unload request satisfies none of
them.

| ID | Criterion |
|---|---|
| AC-HW1-1 | Admission produces a written plan (identification, estimates, available resources, reserve, chosen placement) before any model load is attempted. |
| AC-HW1-2 | The free-VRAM reserve is computed from live measurements and current conditions, not a fixed percentage, and is configurable. |
| AC-HW1-3 | When a model cannot fit fully in VRAM, Howl either layers only backend-supported components or refuses with a clear, safe error; unsupported configurations are never attempted as if supported. |
| AC-HW1-4 | Every load, unload, and offload transition is followed by a verification step, and the verified state is what the UI reports. |
| AC-HW1-5 | Unload evidence shows the weights actually left VRAM and that owned child processes and backend resources are gone. |
| AC-HW1-6 | Pressure during inference is detected and new work is refused or queued before OOM, with a user-visible explanation. |
| AC-HW1-7 | Workload switching follows Unload → Verify → Load → Verify → Work, with active work stopped at a defined safe point rather than killed mid-operation. |
| AC-HW1-8 | Chat, Image, Audio, TTS, and Vision all obtain GPU resources through the shared coordinator; no subsystem bypasses admission. |
| AC-HW1-9 | The shared coordinator is reconciled with the existing process-ownership foundation, with one documented authority and no duplicate cleanup path. |
| AC-HW1-10 | A compatibility/evidence matrix exists in which every hardware/backend combination carries one of the six defined states and its evidence source. |
| AC-HW1-11 | Each power-user override is documented with its risk, and no override silently disables lifecycle or process-safety protection. |
| AC-HW1-12 | Long-run evidence (repeated load/unload/switch cycles) shows no resource leak, no stale process, and no unverified state reported as loaded. |

### 11 · Possible presentation layer — hardware monitor (`FUT-HW2`, UNDECIDED)

A possible future user-facing hardware/resource monitor is recorded here. It is a
**presentation layer over `FUT-HW1` telemetry**, not a separate resource manager and not a
second source of truth.

- It may eventually be persistently visible across application pages.
- It may display Howl's application/process resource consumption and relevant system-resource
  state.
- It must be backed by the shared resource-management and monitoring foundation in `FUT-HW1`.
- Its exact placement, appearance, metrics, interaction, and presentation are explicitly
  **UNDECIDED**.
- It is not a prerequisite for implementing the `FUT-HW1` backend; the backend must be able to
  proceed with no monitor UI at all.
- No UI design is frozen or specified here. Nothing in this entry authorises inventing a
  layout, metric set, or interaction model on the user's behalf.

**`FUT-HW2` — POSSIBLE, UNDECIDED, DOCUMENTED ONLY. Dependent on `FUT-HW1`. Not a competing
resource manager.**

**Status: `FUT-HW1` — DEFERRED, DOCUMENTED ONLY, NOT IMPLEMENTED, NOT AUTHORIZED.**

## Deferred Prompt and Model-Format Intelligence — 2026-11-09

**Requirement ID: `FUT-PROMPT1`** — expanded model-family support with automatic,
evidence-based prompt and chat-format selection.

**Status: DEFERRED — DOCUMENTED ONLY.** Not implemented, not investigated against source by
this entry, not authorised for implementation. Recorded so the requirement can be surfaced
alongside the existing unified model-intelligence and runtime-capability work. It does not
change the current engineering sequence: F3-A-R1 acceptance remains the immediate gate,
followed by F3-B.

### Requirement

Howl must support substantially more model families and automatically determine the correct
prompt/chat formatting for the model actually loaded.

- Expand supported model-family and architecture recognition.
- Detect applicable prompt/chat templates and formatting from reliable model metadata and
  runtime information.
- Automatically populate the appropriate prompt configuration when a model loads.
- Support family-specific prompt tags and formatting conventions.
- Support think/reasoning tags and model-specific reasoning conventions.
- Support Jinja-based chat templates and associated configuration where applicable.
- Support ChatML and other appropriate formats.
- Use ChatML as a sensible generic/default fallback where appropriate, not as a universal
  assumption.
- Avoid guessing from marketing names when actual model configuration, tokenizer metadata,
  templates, or runtime evidence are available.
- Distinguish verified configuration from inferred, conflicting, or unknown configuration.
- Preserve user overrides and make the selected configuration inspectable.
- Ensure incompatible formatting options are not silently applied to a model.

The intended behaviour is similar in spirit to SillyTavern's automatic model/template-aware
defaults, while preserving Howl's offline-first architecture and model-intelligence principles.

### Evidence states and required wording

A loaded model's formatting state must be reported as one of:

| State | Definition | Permitted wording |
|---|---|---|
| Verified | Confirmed from metadata or runtime and exercised against a real model | "verified" |
| Inferred | Derived from consistent metadata/family evidence but not exercised | "inferred" |
| Conflicting | Multiple sources disagree | "conflict — review required" |
| Unknown | No reliable evidence; default applied | "unknown — default applied" |

An inferred or defaulted configuration must never be presented as verified.

### Acceptance criteria

None of these are satisfied by this entry. They define the evidence a future `FUT-PROMPT1`
stage must produce.

| ID | Criterion |
|---|---|
| AC-PROMPT1-1 | Each supported family carries its recognition rule, evidence source, and formatting coverage in one inspectable matrix. |
| AC-PROMPT1-2 | Loading a model populates prompt/format configuration automatically, with no manual template selection required for the verified families. |
| AC-PROMPT1-3 | Family-specific tags, reasoning tags, ChatML, and Jinja templates are each exercised on at least one real model, with the evidence recorded. |
| AC-PROMPT1-4 | The active configuration is visible to the user, labelled with its evidence state, and shows where it came from. |
| AC-PROMPT1-5 | A user override survives reload, and is never silently overwritten by detection. |
| AC-PROMPT1-6 | A detected conflict is surfaced rather than resolved silently, and an incompatible option is refused or disabled rather than applied. |
| AC-PROMPT1-7 | Marketing-name-only evidence is never sufficient to select a template when metadata, tokenizer, template, or runtime evidence exists. |
| AC-PROMPT1-8 | Detection, template parsing, and configuration persist as one reusable subsystem — no parallel duplicate implementations. |

### Roadmap placement and dependencies

Belongs alongside the existing unified model-intelligence and runtime-capability work, before
finalising model-loading and cross-subsystem integration. It must coordinate with:

- Model identification and architecture detection.
- Tokenizer and metadata inspection.
- Chat-template/Jinja parsing.
- Family-specific formatting and reasoning tags.
- Runtime compatibility.
- Model loading and post-load verification.
- User overrides and configuration persistence.

Related existing ledger entries: T1.3 (think tags), T1.4 (correct prompt formatting for
supported families without user action), T3.6 (ternary models), `FUT-C3` (llama.cpp feature
exposure), `FUT-C4` (architecture/family/model-era options), and `FUT-HW1` (post-load state
verification).

Hardware-monitor presentation (`FUT-HW2`) remains dependent on `FUT-HW1` backend telemetry and
has no dependency on `FUT-PROMPT1`.

Neither addition changes the current engineering sequence. F3-A-R1 acceptance remains the
immediate gate, followed by F3-B.

### Preservation of existing work

- Significant work on model and prompt intelligence already exists.
- This requirement is **not** authorisation to replace that work.
- When the roadmap reaches this subsystem, inspect the existing implementation and
  documentation first.
- Identify what already works, what is incomplete, what is unreliable, and what is missing.
- Preserve verified functionality and extend it.
- Avoid duplicating model-family detection, template parsing, or configuration systems
  unnecessarily.

**Status: `FUT-PROMPT1` — DEFERRED, DOCUMENTED ONLY, NOT IMPLEMENTED, NOT AUTHORIZED.**

## F3-B Scope Gate — 2026-11-10

F3-B is the next authorized slice: **lorebook ownership cutover**. Before any implementation,
the documented scope was checked against the accepted F2-C architecture and the F3-A/F3-A-R1
safety baseline. The **intent** is unambiguous. The **contract** required to implement it
safely is not. No F3-B code, migration, or content change was made.

### Documented definition, as it stands

- `docs/review/ARCHITECTURE-F2C-CONTENT-OWNERSHIP-CUTOVER-2026-11-09.md` §10:
  "F3-B — Lorebooks: canonical JSON file store, explicit import from localStorage, read-back
  through the library adapter, legacy store retained."
- `docs/review/ENGINEERING-F3A-OWNERSHIP-SAFETY-MIRROR-SEMANTICS-2026-11-09.md` §10: stable
  IDs, writes through the authoritative lorebook store, explicit exports, and unavailable-root
  and rollback evidence, "before touching worlds, scenarios, characters, chats, images, or
  videos".
- `F3-A8`: deletion/rename propagation is deferred until "lorebook cutover (F3-B) with rollback
  evidence".

There is no F3-B requirements table in this ledger and no per-behavior acceptance list beyond
those three references.

### Blocking decisions

| ID | Decision required | Evidence it blocks implementation | Options |
|---|---|---|---|
| F3-B-Q1 | Canonical file granularity | F2-C §5.4 leaves it open: "one JSON file per lorebook entry, or one versioned lorebook document per collection if a collection concept is added". The current store is a flat entry array and no collection wrapper exists. | One file per entry; one versioned document per collection |
| F3-B-Q2 | Which readers become canonical in F3-B | Sixteen frontend modules read `howl-chat-lorebooks` directly, including chat prompt assembly (`chat-stream.js`), context building (`context-data.js`, `context-selectors.js`), and the world/scenario/character attachment dropdowns. F2-C §5.4 says the lorebook page reads the root and localStorage "becomes a cache", but defines no cache population, refresh, failure, or rebuild contract. | Page-only cutover with a defined cache contract; or an app-wide canonical read path, which touches chat and content forms |
| F3-B-Q3 | Stable identity policy for existing entries | Entry IDs are `Date.now()` strings (`lorebooks.js`) and are referenced by worlds, scenarios, and characters via `howl.lorebooks`. F2-C requires stable IDs. Preserving current IDs keeps REG-1 attachments but is not stable; minting new IDs breaks those attachments unless references are rewritten. | Preserve and remap; mint and rewrite references; dual key |
| F3-B-Q4 | Explicit-import contract | Import is required, but its unit, idempotency, origin marking, failure rollback, and whether the legacy array is modified are unspecified. F2-C's generic collision rules do not define the lorebook mapping. | Define import scope, repeat-import behavior, journal, and rollback evidence |
| F3-B-Q5 | Whether F3-A8 propagation is in F3-B | `F3-A8` names the lorebook cutover as the gate for delete/rename propagation but does not require implementing it there. | Include propagation; or defer it with a new explicit gate |

### Conflict with the accepted architecture

F3-B is defined as lorebook-only and sequenced before chats, worlds, scenarios, and characters.
Chats embed full `lorebooks.json` snapshots, and chat/context/attachment code reads the browser
store directly. A lorebook-only cutover therefore cannot produce one authoritative owner
without either (a) keeping the browser store authoritative for every non-page reader, which is
the unlabelled dual-source state F3-A removed, or (b) changing chat and content-form readers,
which the F3-B sequencing excludes. F2-C §5.4's cache sentence does not resolve this; it only
names the browser store's eventual role.

### Status

F3-B: **AUTHORIZED, SCOPE GATE OPEN, NOT STARTED.** The five decisions above are the only
remaining blockers. No implementation, migration, or user-content change occurred. `REG-1`
remains protected and unchanged, and the roadmap order is unchanged.

## F3-B Lorebook Ownership Cutover — 2026-11-10

F3-B establishes one canonical owner for lorebooks. Entries are JSON documents in
registered lorebooks roots, written only by the shared store. The browser store is now a
derived cache and the legacy seed for the explicit import. No user content was migrated,
deleted, moved, renamed, or rewritten by this stage; the ownership transition is an explicit,
journaled, reversible action the user performs.

Evidence and test results are in
`docs/review/ENGINEERING-F3B-LOREBOOK-OWNERSHIP-2026-11-10.md`.

### Resolved decisions

| ID | Question | Decision | Why |
|---|---|---|---|
| F3-B-Q1 | Canonical file granularity | One JSON document per entry, named by stable identity: `<root>/<id>.json`, envelope `howl-lorebook-v1` | Entries are independently attachable items (REG-1), so per-entry files give atomic single-item writes, no whole-collection rewrite, no collection concept the data does not have |
| F3-B-Q2 | Authoritative read path | One store behind the registry; all lorebook reads go through it. A repository hydrates a derived browser cache so existing synchronous readers keep working; chat prompt sync refreshes before sending | Avoids the rejected dual-source state without rewriting sixteen render paths; canonical is authoritative, the cache is rebuildable |
| F3-B-Q3 | Stable identity | Existing identities are preserved verbatim on import; new entries receive random `lb_` identities; duplicate legacy identities are reassigned and reported, never silently re-linked | Preserves every world/scenario/character attachment without rewriting user content; stability holds from import onward |
| F3-B-Q4 | Import and recovery | Explicit, page-level import of all browser entries; create-only; idempotent by identity plus content hash; conflicts reported without overwrite; journal records created files; rollback removes only those files and clears the cutover marker | Repeatable and reversible with the legacy store untouched as the recovery source |
| F3-B-Q5 | Delete/rename propagation | No separate propagation subsystem. With file ownership, delete removes the canonical file through an explicit action; rename keeps the identity and therefore the file; cross-reference cleanup stays where it already lives | Identity-based paths make propagation unnecessary; adding a second mechanism would duplicate ownership |

### Implemented decisions

| ID | Decision | Status |
|---|---|---|
| F3-B1 | `internal/lorebookstore` is the only writer of lorebook content. Files use the `howl-lorebook-v1` envelope and are reached only through registered roots and library roles. | IMPLEMENTED |
| F3-B2 | Listing is bounded and sequential, one level per root. Foreign formats (including old F2-B export envelopes) are ignored and counted, never read back. Duplicate identities across roots are reported as conflicts. | IMPLEMENTED |
| F3-B3 | A disconnected root keeps its registration and last-known entries, returned as unavailable and non-editable from a bounded index cache. An empty but reachable root is reported empty rather than stale. | IMPLEMENTED |
| F3-B4 | Create requires a save root; replace requires the write role on the root that already holds the entry and re-validates the target through the F3-A-R1 registry protections. Delete re-derives the path from identity and requires a writable root. | IMPLEMENTED |
| F3-B5 | Explicit import and rollback are journaled under the settings directory; rollback removes only files the journal records as created and never touches the browser store or unrelated files. | IMPLEMENTED |
| F3-B6 | `frontend/dist/js/lorebook-repository.js` is the single frontend data path. Before import the browser store remains the source; after import the cache is rebuilt from canonical files. | IMPLEMENTED |
| F3-B7 | The lorebooks page reads through the repository, writes through it, and exposes the import, undo, and status controls. Chat lorebook sync refreshes from the canonical store before sending. | IMPLEMENTED |
| F3-B8 | `REG-1` behavior is unchanged: no attachment, scope, activation, or hot-swap code was modified. | PRESERVED |
| F3-B9 | An import interrupted after file creation is recovered on retry: verified identical files are adopted, journaled as adopted, and the retry completes the marker instead of reporting nothing to do. Rollback leaves adopted files in place because this store did not create them. | IMPLEMENTED |
| F3-B10 | Duplicate identities across roots are reported as conflicts. The first registered root wins deterministically and the same entry is returned on every read, so an existing attachment reference resolves consistently to that winner until the duplication is resolved. | IMPLEMENTED |
| F3-B11 | Consumer audit: every live catalog read goes through the single `howl-chat-lorebooks` cache that the repository hydrates, and every live catalog write goes through the repository. Only the store opens lorebook entry files; no consumer scans roots independently. The chat page still loads `lorebook-manager.js`, whose `saveLorebook`, `deleteLorebook`, and `openLorebookCreationPopup` write the cache directly but have no callers; recorded as a landmine for acceptance. Per-chat `lorebooks.json` snapshots remain chat-scoped copies, not a root scan. | VERIFIED (source) |

F3-B remains scoped to lorebooks. Worlds, scenarios, characters, chats, images, and videos
keep their current ownership until their own stages. Roadmap order is unchanged; the next
ownership slice remains worlds and scenarios (F3-C).

### Acceptance status — 2026-11-10

F3-B is **implementation complete**. It is not yet accepted. Evidence levels are kept
separate:

| Evidence level | State |
|---|---|
| Implementation complete | Yes |
| Automated tests passed | Yes — store, bindings, race, full suite |
| Integration behavior verified | Yes — real bindings and real pages through the project UI harness |
| Runtime acceptance | Partial — harness-verified; the native WebView2 window was not exercised |

Open acceptance checks:

1. Duplicate identity resolution in use: confirm an attachment pointing at a duplicated identity opens the same winner the store returns, and that the conflict is visible.
2. Interruption recovery in a real process: interrupt an import at the file-write, journal, and marker boundaries, restart, retry, and confirm adoption completes ownership. Then confirm rollback returns to browser storage and leaves adopted files.
3. Lorebooks page behavior in the application: create, edit, rename, delete, import, undo, disconnected-root display, and chat prompt refresh.
4. `REG-1` in use: attachments, public and private character lore, world activation, and hot-swapping under the repository and cache arrangement.
5. Persistence across restart: import successfully, close Howl, reopen, and confirm canonical entries, the derived cache, and every world/scenario/character attachment reconstruct without loss or duplication.
6. Legacy write-path check: confirm no reachable UI calls `lorebookManager.saveLorebook`, `lorebookManager.deleteLorebook`, or `lorebookManager.openLorebookCreationPopup` on the chat page, since those write the cache directly and would be discarded by the next hydration.

Runtime pass — 2026-11-10 (UI harness, isolated folder, real Go bindings):

1. Duplicate identity: a second root holding the same identity produced `1 identity conflict need
   review`, one entry resolved deterministically, and a world attachment pointed at that same
   identity.
2. Interruption: deleting the marker while files existed (a crash after file write), reloading,
   and importing again produced journal actions `adopted`, unchanged file mtimes, and a restored
   marker with the message "An interrupted import was verified and ownership completed."
3. Lorebooks page: create (world, scenario, character scopes), rename with unchanged file name,
   delete removing one file, import, undo restoring browser storage, and a surfaced
   `no lorebook save destination is configured` failure all behaved as designed.
4. `REG-1`: a world was created through its form with `Moonlit Archive II` attached; the saved
   reference is the canonical identity and survived a reload. Character and scenario attachment
   forms, activation, and hot-swapping were not exercised.
5. Restart: the harness was stopped, restarted, the browser cache wiped, and the page reloaded;
   both entries rebuilt from files with cutover still true.
6. Legacy write path: source-verified only.

### Final assessment — 2026-11-10

Harness integration acceptance passed; runtime acceptance remains partial. The remaining work is
bounded and does not reopen the architecture.

- **Stale Import button (low severity).** Observed only when the marker was deleted while the
  page stayed open, which is an out-of-band simulation. After a real restart the page hydrates
  the marker-less state and Import is enabled, so crash recovery does not normally require a
  reload. The limitation is that the store bar does not self-refresh when the marker changes
  underneath a live session; a manual refresh is the remedy. Not a recovery defect.
- **`REG-1` partially verified.** World attachment saved a canonical reference and survived a
  reload. Scenario and character attachment forms, activation, private/public lore, and
  hot-swapping are untested.
- **Native runtime outstanding.** The harness reflects bindings over HTTP; a packaged WebView2
  session is still required for final acceptance.

Remaining bounded acceptance work:

1. Native WebView2 smoke test: import, close, reopen.
2. Character and scenario attachment forms, and attachment resolution after import.
3. Chat prompt synchronization with a running model.
4. Activation and hot-swapping in a creative workflow, including private/public lore.
5. Windows symlink listing test with suitable privileges.

Status remains: implementation complete, automated tests passed, integration verified, runtime
partial. F3-B stays open for final acceptance; no architecture change is authorised.

### Final acceptance pass — 2026-11-10

Verdict: **ACCEPTED WITH DOCUMENTED LIMITATIONS**. This supersedes the open status above. The W27
native Lorebook lifecycle gate is now satisfied by the evidence in
`docs/review/ENGINEERING-ISOLATED-NATIVE-RUNTIME-2026-11-10.md` §27: two stable character IDs,
native editor ownership, canonical import, rename, character-specific prompt inclusion/exclusion,
inference-boundary capture, unresolved deleted owner, and restart persistence all passed. No model
inference was performed. F3-C is not started; the other documented limitations remain unchanged.

Checks exercised in this pass (isolated folder, project UI harness, real bindings):

| Check | Status | Evidence |
|---|---|---|
| Native desktop build | Passed | `go build -tags desktop` exit 0; binary produced |
| Native WebView2 runtime | Blocked | WebView2 user-data folder cannot be isolated from the real app profile with available options; launching would expose/alter the user's stored data, so it was not run |
| Scenario attachment | Passed | `Black River Vigil` saved `lorebooks: [{id: 1794283090983, title: Tide Oath}]`, survived reload |
| Character attachment | Passed | `Rill Final` saved `howl.lorebooks: [{id: 1794283090983, title: Tide Oath}]` |
| Reusable attachment | Passed | The same identity is attached to both the scenario and the character |
| Legacy catalog writes | Passed | No caller of `lorebookManager.saveLorebook`, `deleteLorebook`, or `openLorebookCreationPopup` outside `lorebook-manager.js` |
| Symlink listing | Blocked | `mklink` returns "You do not have sufficient privilege to perform this operation"; test remains skipped |
| Chat prompt synchronization | Unverified | No prompt-preview binding exists and no inference server is configured, so `SetLorebooks` delivery to a prompt cannot be observed |
| Activation / hot-swapping | Unverified | Requires an active chat workspace; not exercised |
| Private character lore injection | Unverified | Private visibility field is stored canonically, but injection scope was not exercised |
| Disconnected-root display | Unverified | Not exercised in this pass |

Additional observation (pre-existing, low severity, not F3-B): reopening the character creation
form can display stale selected-lorebook chips because `closeCharacterCreationPopup` clears the
state without re-rendering the chip list, and `openCharacterCreationPopup` does not render it.
Saving in that state correctly persists the empty selection. This belongs to the character
surface, not to lorebook ownership.

Limitations are bounded and do not undermine F3-B ownership, persistence, recovery, or
integration guarantees: import, identity preservation, rollback, restart reconstruction,
interruption adoption, duplicate conflict reporting, and world/scenario/character attachment
resolution were all exercised successfully.

### F3-B closure review — 2026-11-10

Status: **F3-B CLOSED**. The dependency chain was re-checked after W27/W27a and no regression or
contradictory evidence was found; the ACCEPTED WITH DOCUMENTED LIMITATIONS verdict stands. F3-C is
unblocked and not started.

| Check | Result | Evidence |
|---|---|---|
| Every W27/W28 change is represented in the report and ledger | Pass | Report §§23–28; rows W27, W27a, W28, W28a–W28c |
| No stale status describes completed work as pending | Pass after correction | Report §1/§23.5 corrected, §7 annotated; F3-B §13 and Wails-CLI addenda record the cleared blocker |
| No production build carries seams or capture instrumentation | Pass | Rebuilt `howl-w27-native.exe`: 0 hits for seam, token-header and `promptcapture` symbols; `BuildCharacterLoreBlock` present |
| Canonical lorebook path is the only authoritative writer | Pass (source) | Only `internal/lorebookstore` writes entry files; `lorebook-manager.js` cache writes have no callers |
| Identity and memory privacy invariants hold across consumers | Pass | memory, chat and root Go suites exit 0; frontend identity and memory tests 4/4 each |
| Conflict grouping leaves world/scenario/global behavior unchanged | Pass | `TestConflictGroupingKeepsNonCharacterScopesShared` |
| Affected tests and production build pass from the final tree | Pass | `go test` lorebook/memory/chat/root exit 0; `build\wails.cmd build` exit 0 |
| No test configuration or isolation setting can affect a normal launch | Pass, one observation | No `.env`, no `HOWL_*` variable, no marker in the tree; isolated run wrote only its own root. A non-isolated 52 ms process start and Go-test discovery scans remain in the real startup log (log-only). |

Not verified, unchanged: successful model inference, activation/hot-swapping in a live creative
workflow, disconnected-root display in the native window, and the Windows symlink listing test.
None is an F3-B ownership failure.

### F3-C dependency and acceptance review — 2026-11-10

Status: **DECISIONS RECORDED — IMPLEMENTATION NOT STARTED.** No code was changed. The decision and
compatibility review is `docs/review/REVIEW-F3C-DECISIONS-AND-COMPATIBILITY-2026-11-10.md`;
the earlier dependency review is `docs/review/REVIEW-F3C-DEPENDENCY-AND-ACCEPTANCE-2026-11-10.md`.

Inherited from F3-B: one canonical writer per category, identity-named paths, create/replace role
rules, explicit journaled import with rollback, cutover markers, bounded scan with conflict
reporting, unavailable-root handling, atomic writes with structured failures, and the
isolation-plus-prompt-boundary native evidence method. World/scenario registry categories,
roles, default destinations, and isolation already exist.

Not inherited, and therefore F3-C's actual work: there is no canonical world or scenario store.
Ownership is still `howl-worlds` / `howl-chat-scenarios` localStorage; native world/scenario
directories are cwd/exe-relative; `ExportContentItem` writes non-authoritative copies and has no
frontend caller; chat reads resolve worlds/scenarios by name; chat snapshots embed full documents
with no revision field; world/scenario lore uses the name as `OwnerID`.

| ID | Decision needed | Status |
|---|---|---|
| F3-C-Q1 | Folder per item with a primary document; assets are referenced through the media manifest | DECIDED |
| F3-C-Q2 | Stable IDs; preserve unique legacy IDs verbatim, mint `wld_`/`scn_` for new items, journal every mapping | DECIDED |
| F3-C-Q3 | Explicit, create-only, idempotent, journaled, reversible folder import; browser data untouched | DECIDED |
| F3-C-Q4 | Pinned chat snapshots with recorded revisions, visible staleness, and an explicit refresh operation | DECIDED |
| F3-C-Q5 | Preserve dangling references, display them unavailable, exclude them from prompts, never substitute | DECIDED |
| F3-C-Q6 | World/scenario lore resolves by stable IDs in this slice; renames and recreation never transfer ownership | DECIDED |

Compatibility pass — nine findings (seven conflicts, two gaps), none overturning a decision:

| ID | Conflict | Disposition |
|---|---|---|
| C1 | World IDs are bare timestamps; scenarios use two shapes (`scenario-<ms>-<rand>` and bare `<ms>`) | Identity stays opaque; preserve unique IDs, reassign duplicates through a journaled map |
| C2 | `chat/persistence.go` backfills `WorldID`/`ScenarioID` from names | Migrate by explicit mapping; unresolved names stay unresolved |
| C3 | `chat/service.go` reads and writes world/scenario lore by name path, outside the store | PREREQUISITE — resolve by store ID and relocate generated lore writes |
| C4 | Reference resolvers fall back to display-name matching | Typed refs; ID-only resolution in the canonical path |
| C5 | Snapshots are already pinned but carry no revision and refresh touches character sets | Add revision fields, staleness indicator, explicit refresh |
| C6 | World/scenario editors rebuild documents from form fields and drop unknown keys | PREREQUISITE — repository save must merge unknown fields |
| C7 | Folder-per-item assets conflict with the F3-F paired media root | Item folder holds `media.json` references; binaries stay F3-F |
| C8 | No native world/scenario writer exists; legacy `Worlds/<Name>` copies are non-authoritative | Import from browser records; keep legacy folders read-only |

Added acceptance requirements: typed cross-category references (replacing `[]interface{}` and the
duplicated `lorebooks`/`howl.lorebooks` arrays), distinguishable failure states for root/identity/
import/reference/write failures, and lossless preservation of unknown legacy fields.

Still open for decision, not blocking the store: refresh semantics (document only vs document plus
character set), whether generated world/scenario lore is canonical or chat-scoped, and whether any
binary asset may live inside an item folder before F3-F.

Acceptance for the slice requires: stable IDs and rename stability, a single writer with surfaced
failures, reversible explicit import, a canonical read path with a rebuildable cache, reference
resolution by identity, a visible snapshot revision, unavailable-root exclusion from prompts,
native isolated evidence with prompt-boundary capture, and unchanged lorebook/memory regressions.
Model inference remains separate evidence. F3-C is unblocked; F3-D through F3-H remain in order.

### F3-C decisions resolved and contracts frozen — 2026-11-10

Status: **CONTRACTS FROZEN — IMPLEMENTATION NOT STARTED.** No code was changed. The frozen record
is `docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md`.

| ID | Resolution |
|---|---|
| R1 | Document refresh and scenario cast sync are separate explicit operations; loading or editing never mutates a chat; a missing source keeps the snapshot and refuses refresh |
| R2 | Manual item lore is canonical; dialogue-derived lore is chat-owned in `memories_snapshot.json` until explicitly promoted; legacy `automated_*_lore.json` imports as `legacy_generated_unattributed`, flagging it without inventing provenance; chat stops writing world/scenario lore |
| R3 | Item folders hold the document and `media.json` references only; binaries stay with F3-F; legacy IndexedDB refs preserved as references |

Frozen contracts: `howl-world-v1` / `howl-scenario-v1` envelopes with store-computed `revision`;
opaque IDs preserving unique legacy values (`wld_`/`scn_` for new); `<root>/<id>/` paths; single
writer with create/replace roles; lossless update merging unknown fields with `baseRevision` checks;
typed `{kind, id}` references resolved by ID only; explicit preflighted journaled import and rollback;
chat snapshot revisions with visible staleness; and distinct failure states. Acceptance A1–A13
includes regressions for the second writer, lossy saves, and title fallback.

Consequences kept visible: new generated lore no longer becomes shared world knowledge automatically;
legacy generated lore is canonical but visibly unattributed; character generated lore keeps its
current behavior until F3-D and must not be extended here.

### F3-C store foundation implemented — 2026-11-10

Status: **STORE FOUNDATION IMPLEMENTED AND TESTED — CONSUMERS NOT CONNECTED.** `internal/worldstore`
provides canonical envelopes (`howl-world-v1` / `howl-scenario-v1`), identity validation with
legacy-preserving IDs, create-only save, role-enforced update, deterministic `rv1-` revisions with
mismatch detection, lossless merge of unknown fields with `baseRevision` conflicts, and typed
`{kind, id}` references resolved by ID only. Record:
`docs/review/IMPLEMENTATION-F3C-STORE-FOUNDATION-2026-11-10.md`.

| Check | Result |
|---|---|
| `gofmt` and `go vet` on the new package | clean, exit 0 |
| `go test ./internal/worldstore ./internal/library` | both `ok`; 33 store/revision/reference/concurrency/availability/authority tests pass |
| A13 regression (`lorebook`, `chat`, `memory`, root) | all `ok`, exit 0 |
| Frontend identity and memory tests | exit 0, 4/4 each |

Deliberately not implemented at foundation time: import and rollback, chat-service conversion,
snapshot revisions and refresh, media manifests, editor/repository wiring, and the unavailable-root
index. The chat-service conversion was completed separately on 2026-11-10; see below.

Correctness gates closed later on 2026-11-10: read-check-write is serialized process-wide, so
concurrent updates sharing a base revision produce exactly one winner and the rest
`revision_conflict`, including across two `Store` instances; concurrent creates produce one `saved`
and the rest `collision`; a blocked atomic temp write produces `write_failed` with the canonical
document byte-identical; `go test -race ./internal/worldstore` passes. `Lookup` now separates
proven absence (`ErrNotFound`) from unavailable or incompletely searched storage
(`ErrUnavailable`), and `Update` reports `destination_unavailable` rather than `not_found`. Consumers
must never recreate, import, overwrite, or reassign references on an indeterminate search. The write
lock is per process; cross-process writers are outside the single-process desktop contract.

Root authority resolved 2026-11-10 in `docs/review/DECISION-F3C-ROOT-AUTHORITY-2026-11-10.md`:
roots are equivalent search locations, never ranked; path order is retrieval order only. A found
copy is canonical only when every root was searched cleanly and exactly one root holds the
identity. Otherwise the result is found-but-authority-indeterminate (`AuthorityIndeterminate`,
`Copies`), the item stays readable with `Editable=false`, and `Update` refuses with
`authority_indeterminate` without writing. This closes the stale-copy question without adding an
availability index or a duplicate resolver.

### F3-C chat consumer converted — 2026-11-10

Status: **CHAT SERVICE IS A CANONICAL CONSUMER; WORLD/SCENARIO GENERATED-LORE WRITES REMOVED.**
Record: `docs/review/IMPLEMENTATION-F3C-CHAT-CONSUMER-2026-11-10.md`.

`chat.Service` resolves world/scenario documents and lore by stable ID through `worldstore`; display
names are labels only. Manual lore and legacy automated lore **already inside a canonical item
folder** are classified `legacy_generated_unattributed` and stay prompt-eligible with no invented
attribution. Files still in old name-keyed folders are not read at runtime; LF-1 authorizes them only
as explicit, journaled migration input, not as a fallback. New conversation discoveries stay
chat-owned. Name-keyed generated-lore and memory
filters are ID-only, with empty or foreign owners excluded. `SaveChatData` no longer writes
`automated_world_lore.json` or `automated_scenario_lore.json`; character generated lore is retained
behind a single allowlist until F3-D. Chat persistence no longer copies names into
`WorldID`/`ScenarioID`.

Classification detail: `docs/review/F3C-LEGACY-LORE-CLASSIFICATION-2026-11-10.md`.

| Check | Result |
|---|---|
| `gofmt` and `go vet` on touched packages | clean, exit 0 |
| `go test ./internal/worldstore` | `ok`; 38 tests |
| `go test ./internal/backend/chat` | `ok`; 23 tests |
| `go test -race ./internal/worldstore` | pass |
| A13 regression (`lorebook`, `chat`, `memory`, root) | all `ok`, exit 0 |
| Frontend identity and memory tests | exit 0, 4/4 each |
| Source audit | no direct `Worlds`/`Scenarios` paths in `internal/backend/chat` |

Not started: editor wiring, import/rollback, promotion, snapshot revisions/refresh, cast sync, media
manifests, availability index, and F3-D character generated lore. `chat.Service` is not the desktop
inference path; the app's prompt builder consumes the chat-owned snapshot and is unchanged in this
slice. Without import, a legacy chat whose world exists only by name resolves to `not_set` rather
than falling back to a name folder.

### F3-C repository-wide single-writer audit — 2026-11-10

Status: **SOURCE-LEVEL SINGLE WRITER; A10 NOT YET DECLARED.** Record:
`docs/review/AUDIT-F3C-SINGLE-WRITER-2026-11-10.md`.

Two real legacy surfaces were found and closed:

- Dead name-keyed disk readers in `app-characters.go` (`worldExistsOnDisk`,
  `scenarioExistsOnDisk`, and their directory helpers) were removed. They had no callers, but they
  were still legacy readers in the repository.
- `ExportContentItem` could write a non-authoritative `howl-content-export-v1` envelope to
  `<hash>/world.json` or `<hash>/scenario.json` under a registered root. It now refuses worlds and
  scenarios before any destination lookup; generic export tests use characters instead. This is a
  deliberate capability gap until store-backed world/scenario export exists.

After those changes: `worldstore.Save`/`Update` are the only canonical document writers; chat cannot
write under item roots; no production deletion path targets world/scenario roots; the desktop prompt
builder still consumes the chat-owned snapshot; the desktop editor still writes browser
`localStorage` records. The browser records are the intended explicit import source, not a canonical
writer.

A10 source audit: PASS. A10 acceptance: NOT CLOSED - frozen evidence also requires a native runtime
check that chat cannot write item folders. No native run was performed.

Resolved before import: LF-1 amends frozen §6 to authorize one explicit, journaled, reversible
migration of old `Worlds/<name>` / `Scenarios/<name>` folders. Mapping requires exact stable-ID
evidence or user resolution, never a title alone; unresolved content is quarantined, source folders
are never modified, and runtime name reads/writes remain prohibited. Record:
`docs/review/DECISION-F3C-LEGACY-FOLDER-MIGRATION-2026-11-10.md`. Migration and import code are not
started; A9 was still open at that point (later passed).

### F3-C import design recorded — 2026-11-10

Status: **DESIGN RECORDED; IMPLEMENTATION NOT STARTED.** Record:
`docs/review/DESIGN-F3C-IMPORT-2026-11-10.md`.

The design covers two explicit sources: browser records supplied as a deep copy by the renderer, and
LF-1 legacy name folders selected by the user. Preflight is observational and writes nothing; it
returns a plan with per-item dispositions and a plan hash. Execution is refused when the plan is
stale. Canonical writes go through new store import primitives, never through the registry directly.
Staging lives under a hidden non-item path; commit is an atomic no-overwrite rename; journal and
marker live under config. Rollback removes only journal-created, unchanged items and never touches
adopted items or source folders. Identity requires exact stable-ID evidence or explicit user
resolution; duplicates, missing IDs, conflicts, unresolved owners, and inline media are quarantined
with distinct reasons. No silent document or lore merge is allowed.

Acceptance is LF1-LF12 plus A9, A10 native runtime evidence, and A13. No import code, store import
primitive, Wails binding, or UI was started.

Design gate tightened 2026-11-10 against the actual store and filesystem behavior. Resolved in the
record: per-item (not transaction-wide) crash consistency with a `planned` journal state and a
rename-before-`committed` recovery path; commit-time revalidation under the existing process-wide
`writeMu`; staging isolation under `.howl-import-staging` with an ownership marker and same-root
rename; concurrent staging but serialized commit/save/update/removal; and two-phase rollback that
verifies every document/lore/media hash, rejects unexpected files, and stops on ambiguous evidence.
The crash-state and acceptance matrices now include deterministic fault injection at every boundary.

### F3-C explicit import and LF-1 migration implemented — 2026-11-10

Status: **IMPLEMENTED AND FOCUSED GATES PASS; NATIVE A9/A10 EVIDENCE OUTSTANDING.** Record:
`docs/review/IMPLEMENTATION-F3C-IMPORT-2026-11-10.md`.

`internal/worldstore/import.go` owns hidden operation-scoped staging, no-overwrite commit under the
existing process-wide `writeMu`, owned-staging re-read, commit verification, and hash-guarded
removal. `internal/contentimport` implements observational browser/LF-1 preflight, a re-derived
plan hash, journaled per-item execution, explicit resume, and two-phase rollback. Wails bindings in
`app-imports.go` expose preview, execute, resume, rollback, and read-only journal listing; the
services are initialized in `app.go`. No editor, snapshot, media-byte, lore-promotion, or F3-D work
was added.

| Check | Result |
|---|---|
| A13 regression (`lorebook`, `chat`, `memory`, root) | all `ok`, exit 0 |
| `go test ./internal/worldstore ./internal/backend/chat ./internal/contentimport .` | all `ok` |
| `go test -race ./internal/worldstore ./internal/contentimport` | pass |
| `go vet ./internal/worldstore ./internal/contentimport .` | exit 0 |
| Frontend identity tests | 4/4 and 2/2, exit 0 |

Crash tests cover after journal, complete staging, incomplete staging, after prepared, after the
canonical rename but before `committed`, after `committed` before the marker, and terminal marker
no-op. Concurrent imports and ordinary creates have one winner; commit revalidates registration,
`Lookup`, final absence, and hashes. Rollback verifies every document/lore/media hash first and
refuses changed lore, adopted items, and unexpected files; source trees stay byte-identical.

At that stage A9 had no isolated native import run and A10 had no runtime evidence: the source
audit passed, but no native chat load/save and generated-lore synchronization run had been performed
to prove canonical item trees were unchanged. Both later passed; see the native acceptance section
below.

The post-audit correction pass now normalizes references losslessly into the frozen typed shape,
preserves unknown fields, and keeps ID-only resolution.

### F3-C import audit correction pass — 2026-11-10

Status: **FOCUSED CORRECTIONS COMPLETE (A9/A10 LATER PASSED); A13 GREEN.** The accepted audit findings
F-01 through F-09 were addressed without UI wiring, F3-D, media-byte migration, or native-acceptance
claims. Records: `docs/review/AUDIT-F3C-IMPORT-BOUNDARIES-2026-11-10.md`,
`docs/review/DESIGN-F3C-IMPORT-2026-11-10.md`, and
`docs/review/IMPLEMENTATION-F3C-IMPORT-2026-11-10.md`.

| Finding | Disposition |
|---|---|
| F-01 typed references | Corrected losslessly; legacy-name/malformed/missing-ID/unknown-kind are classified; runtime resolution is ID-only |
| F-02 registered legacy roots | Corrected; every registered root and canonical parent/child path is rejected with `source_registered` or destination overlap |
| F-03 media adoption | Corrected; normalized media must match exactly or the item is `media_conflict` |
| F-04 empty owners | Corrected; absent/empty/whitespace owners are stamped; foreign owners are not rewritten |
| F-05 staging cleanup | Corrected; abort is serialized and exact-content checked; cleanup failures are visible; rollback refuses artifacts |
| F-06 Wails bindings | Regenerated; five content-import methods present; native build passes; native UI acceptance remains A9 |
| F-07 result aggregation | Corrected; only identity/lore/media conflicts count as conflicts |
| F-08 marker record | Corrected in design: `owner.json` at `<itemID>.owner.json` beside staging; no journal invalidation |
| F-09 media schema | Corrected; canonical/legacy reference-only envelopes validated; unknown fields preserved; inline payloads detected |

| Check | Result |
|---|---|
| `go test ./internal/contentimport ./internal/worldstore -count=1` | both `ok` |
| A13 regression (`lorebook`, `chat`, `memory`, root) | all `ok`, exit 0 |
| `go test -race ./internal/worldstore ./internal/contentimport -count=1` | both `ok` |
| `go vet ./...` | exit 0 |
| Frontend identity tests | 4/4 and 2/2, exit 0 |
| `build\wails.cmd generate module` | exit 0; `ContentImport` methods present in generated bindings |
| `build\wails.cmd build` | exit 0; `build\bin\howl-chat.exe` produced |

A recorded red-before/green-after run is not available for every finding because the focused tests
were authored after the first implementation edit. The audit and baseline records retain the
pre-correction source evidence. No frozen-contract meaning was changed: F-01 implements the existing
typed-reference rule, F-08 is a documentation refinement, and F-03/F-09 implement the existing
media-content rule. **At that stage A9 and A10 were still open; they later passed. A13 is regression evidence only.**

### F3-C final marker/cleanup/media correction — 2026-11-10

Status: **CORRECTION COMPLETE (A9/A10 LATER PASSED); A13 GREEN.** The final targeted verification
found marker-integrity, cleanup-aggregation, and media-classification defects. They were corrected
without UI/editor wiring, media-byte migration, or F3-D work. Records:
`docs/review/FINAL-VERIFICATION-F3C-IMPORT-2026-11-10.md`,
`docs/review/IMPLEMENTATION-F3C-IMPORT-2026-11-10.md`,
`docs/review/AUDIT-F3C-IMPORT-BOUNDARIES-2026-11-10.md`, and
`docs/review/DESIGN-F3C-IMPORT-2026-11-10.md`.

| Correction | Result |
|---|---|
| Operation-scoped completion markers | Marker path and contents bind category, operation, plan, destination, item set, and terminal status; exclusive creation and ownership checks prevent overwrite |
| Marker disagreement | Missing, malformed, stale, contradictory, and foreign markers return `integrity_failure` before recovery; canonical content is preserved |
| Partial rollback | `rollback_partial` removes the completion marker after ownership verification |
| Cleanup aggregation | `cleanup_failed` and mixed outcomes count as quarantined, with operation `ambiguous`; prepared rollback is `rolled_back`, not `cleanup_complete`, and is not counted; in-progress items use `pending`, not `quarantine` |
| Media classification | Malformed and unsupported manifests are distinct; missing format is malformed; legacy/canonical unknown fields and legacy adoption are tested |
| Reference evidence | Full edge matrix and lossless update/reload retention are tested; no production runtime consumer exists in this slice, recorded as an integration boundary |

| Check | Result |
|---|---|
| `go test ./internal/contentimport ./internal/worldstore -count=1` | both `ok` |
| `go test ./internal/backend/lorebook ./internal/backend/chat ./internal/backend/memory . -count=1` | all `ok`; A13 exit 0 |
| `go test -race ./internal/worldstore ./internal/contentimport -count=1` | both `ok`; exit 0 |
| `go vet ./...` | exit 0 |
| `gofmt -l internal/contentimport internal/worldstore` | no output |
| Frontend identity tests | 4/4 and 2/2, exit 0 |
| `build\wails.cmd generate module` | exit 0 |
| `build\wails.cmd build` | exit 0; native binary produced |

A9 and A10 were not run by this correction. Their later native acceptance is recorded below. A13
remains a separate regression gate and passed.

### F3-C native A9/A10/A12 acceptance passed — 2026-11-10

Status: **A9 PASS; A10 PASS; A12 PASS.** Records:
`docs/review/VERIFICATION-F3C-NATIVE-A9-A10-2026-11-10.md` and
`docs/review/VERIFICATION-F3C-NATIVE-A12-2026-11-10.md`.

A test-tagged native WebView2 build ran against a marked isolated root outside the project. A9 used
the real content-import preview/execute/resume/rollback/journal bindings: two importable items and
three quarantines, canonical creation, no source mutation, rollback of only journaled creations,
idempotent adoption, and a zero-write identity conflict. A10 created canonical world and scenario
items, loaded chat data by stable ID, called `LoadChat`, `SaveChat`, `SaveChatData`, and the
generated-lore save path, and left both canonical item trees byte-identical; chat files were written
only under `Chats/A10Chat`. The production build contained none of the test seam symbols.

| Check | Result |
|---|---|
| A9 native import | **PASS** |
| A10 native single writer | **PASS** |
| `go test ./...` | all `ok`, exit 0 |
| `go vet ./...` | exit 0 |
| Production `howl-prod-a10.exe` seam scan | 0 matches for seam/token/test symbols |
| A12 prompt-boundary native evidence | **PASS** |

Limits: no import UI, no model inference, and the post-rename cleanup fault injection remains
source-inspected. The A10 generated-lore path used a build-tagged seam because no production binding
exposes it; that seam is absent from production. A12 captured the production prompt with synthetic
world, character, scenario, lore, and user markers in production order; disabled lore was absent.

### Hardware compatibility targets — 2026-11-10

Recorded as future validation targets, not implementation. Two additional laptops are available: a
gaming laptop for GPU/VRAM/driver/resource coverage, and a low-end laptop for memory pressure, CPU
fallback, small-model loading, and graceful degradation (roughly 1B as the first candidate, 3B only
if that hardware supports it).

| ID | Target | Status |
|---|---|---|
| HW-C1 | Hardware compatibility matrix covering model families, VRAM/RAM, offload, generation performance, and failure behavior | DEFERRED — NOT STARTED |
| HW-C2 | Low-end profile must degrade gracefully rather than fail opaquely | DEFERRED — NOT STARTED |
| HW-C3 | Compatibility claims require measurements on the named machine, not inference from the desktop | DEFERRED — NOT STARTED |

### Desktop-complete architecture requirement — 2026-11-10

Recorded as a binding product constraint, not an implementation item. Howl is a complete desktop
application on its own; a mobile companion is optional and never a dependency.

| ID | Requirement | Status |
|---|---|---|
| MOB-1 | Desktop functionality stays fully usable with no phone, no network service, and no companion app | CONSTRAINT RECORDED |
| MOB-2 | Mobile connectivity is never required for local inference, creative workflows, or normal application operation | CONSTRAINT RECORDED |
| MOB-3 | Connecting a phone extends an existing capability; it never unlocks functionality the desktop needs to work | CONSTRAINT RECORDED |
| MOB-4 | Android is the mobile test platform when that phase is authorized | CONSTRAINT RECORDED |

No mobile work is authorized or started. Any future remote API or pairing design must satisfy
MOB-1 through MOB-3 before it is accepted, and must not become a hidden dependency of desktop
startup, storage, or inference.

## Wails CLI Integration and Build Reliability — 2026-11-10

Native desktop builds run through one project-local entry point, `build\wails.cmd`, which drives
`build\tools\wails.exe` (v2.12.0) with the Go toolchain this module pins (go1.23.8) and changes
nothing global: no PATH, GOBIN, Go environment, registry, or machine-wide install.

| Command | Purpose |
|---|---|
| `build\wails.cmd build` | native build → `build\bin\howl-chat.exe` |
| `build\wails.cmd dev` | live-reload development session |
| `build\wails.cmd doctor` | environment report |
| `build\wails.cmd generate module` | regenerate `frontend\wailsjs` bindings |

Decisions:

| ID | Decision | Status |
|---|---|---|
| W1 | The wrapper is the supported command surface. It switches to the project root, preflights the CLI and project files, pins the toolchain from `go.mod` for that invocation only, verifies `go version` actually resolved to it, forwards arguments unchanged, and propagates the CLI exit code. It aborts rather than falling back to the machine default toolchain. | IMPLEMENTED |
| W2 | `start-howl-chat.bat` dev mode calls the wrapper instead of looking for a global `wails.exe` on `GOPATH\bin`. | IMPLEMENTED |
| W3 | The regenerated bindings are retained, unedited. 176 exported `App` methods, 175 bound; the only omission is `OnShutdown(context.Context)`, a lifecycle hook the generator never binds and the frontend never calls. Regeneration is byte-identical across runs. | RETAINED |
| W4 | Native WebView2 runtime acceptance stays blocked: Wails v2.12 exposes the profile folder only through `options.App.Windows.WebviewUserDataPath`, which `main.go` does not set, so a launch would use the real user profile. | BLOCKED |

No dependency, Go, or Wails version was changed; `go.mod` and `go.sum` are untouched by the
build. F3-C was not started.

Evidence and verification results are in
`docs/review/ENGINEERING-WAILS-CLI-INTEGRATION-2026-11-10.md`.

The smallest enabling task for native acceptance, recorded there and not implemented, is an
env-gated isolated WebView2 profile plus app storage root, so an isolated run cannot read or
write real user data while an unset run behaves exactly as today.

## Isolated Native Runtime Enablement — 2026-11-10

The real desktop application can now be tested natively against a dedicated test root without
reading or modifying user data. The single opt-in switch is `HOWL_WEBVIEW_DATA_DIR`; it must name
an absolute, writable directory outside the project tree containing a `.howl-isolated-root` marker
whose content is exactly `howl-isolated-test`. With it set, the WebView2 profile is
`<root>\webview2`, all application state is `<root>\data`, and scratch files are `<root>\temp`.
Unset means production behaviour, unchanged; set-but-invalid aborts startup with exit 2 and never
falls back to production directories. Install assets (models, runtime binaries) still resolve from
the real checkout because they are inputs, not state.

| ID | Decision | Status |
|---|---|---|
| W5 | One opt-in root holds the WebView2 profile, all application state and scratch space; production resolution is untouched when the variable is unset | IMPLEMENTED |
| W6 | A marker file is required, so a mistyped root can never adopt a directory that holds real data; every derive is checked with component-wise containment | IMPLEMENTED |
| W7 | Native smoke test against a verified root: window, frontend↔Go bindings, Lorebooks page, graceful shutdown; real profile and config byte-identical before and after | VERIFIED |
| W8 | Unwritable-root rejection and creating content through the Lorebooks UI were not exercised end to end | OPEN |

No storage default, dependency, Go or Wails version was changed. F3-C was not started. F3-B
acceptance items that this task did not touch (prompt synchronization, activation/hot-swapping,
private-lore injection, disconnected roots, symlink listing) remain open.

Evidence: `docs/review/ENGINEERING-ISOLATED-NATIVE-RUNTIME-2026-11-10.md`.

### Isolated native run — UI round (follow-up, 2026-11-10)

- Content-library defaults now resolve inside the isolation root (`data\content\<Category>`), so an
  import can never target the real Documents/Pictures/Videos folders.
- Unwritable isolation root exercised end to end: a temp-only deny ACE produced
  `ISOLATION-ERROR ... is not writable: ... Access is denied.` and exit code 2; no real folder's
  permissions were changed.
- A lorebook entry was created through the real Create Entry dialog, persisted in the isolated
  WebView2 profile, and was reconstructed with its content and flags after a restart.
- `Import to Library` cannot succeed: `lorebooks.html` runs without the Wails bindings
  (`window.go` absent), so `lorebook-repository.js` throws "The lorebook library is not available
  in this window." and the page silently falls back to localStorage. Pre-existing defect,
  unrelated to isolation; no canonical file is produced.

| ID | Decision | Status |
|---|---|---|
| W8 | Unwritable-root rejection proven at runtime on a temp-only root | VERIFIED |
| W9 | Lorebooks page could not perform the library import: Wails injects its runtime only into paths ending in `/` or `/index.html`, so standalone pages ran without `window.go` | FIXED — AssetServer middleware applies the same inject to every HTML page |
| W11 | One shared injection seam (`webassets.go` + `assetserver.Options.Middleware`) covers all current and future pages; a test walks every embedded `.html` file | IMPLEMENTED |
| W12 | Import creates the canonical file in the registered isolated root with identity and metadata; Undo Import removes only imported copies; restart reads the canonical state | VERIFIED |
| W10 | An entry created in the isolated run survives restart and is reconstructed from the isolated profile | VERIFIED |

Evidence: `docs/review/ENGINEERING-ISOLATED-NATIVE-RUNTIME-2026-11-10.md` (§5.1, §6).

### W13 — F3-B prompt synchronization with canonical lorebook state

| ID | Decision | Status |
|---|---|---|
| W13 | After cutover the chat persisted a lorebook snapshot and re-sent it on load, so a canonically edited entry never reached the prompt; `SetLorebooks` now overlays canonical file content on the chat selection (chat keeps `Enabled`, files own the text) | FIXED IN CODE — unit-tested, not yet natively verified (no model loaded) |

Evidence: `docs/review/ENGINEERING-ISOLATED-NATIVE-RUNTIME-2026-11-10.md` §10; `TestSetLorebooksResolvesCanonicalContent`.


| W14 | Ownership boundaries after cutover: deleted canonical entries are dropped from a library-owned pool on a complete view, disconnected roots keep last-known state, and a lorebook created after a chat save can be enabled into that chat | FIXED IN CODE — unit-tested, not natively accepted |


| W15 | Private character lore leaked: `visibility` was never read and unowned character-scope entries (the Lorebooks form writes no `ownerId`) were routed to every character; hidden entries are now excluded and unowned character lore is shared only when public | FIXED IN CODE — unit-tested, not natively accepted |


| W16 | Disconnected roots and deletion were one state: retained last-known entries were still injected as current canon; entries are now marked `Stale` (in memory only, never persisted), excluded from prompt assembly, and cleared when the root reconnects | FIXED IN CODE — unit-tested, not natively accepted |


| W17 | No inference seam existed for prompt-boundary verification (`App.httpClient` was a concrete `*llama.HTTPClient`); added a narrow `inferenceClient` interface and a recording client that captures the exact messages at the real send boundary | IMPLEMENTED + PROMPT-BOUNDARY TESTED IN-PROCESS — not native GUI, not model inference |

| W18 | Character turns send a frontend-composed prompt and `resolveCharacterLorePromptBlock` has no production call site, so character-lore privacy depends on the frontend composer | OPEN |


| W19 | W18 corrected: the frontend character composer no longer rebuilds lore from localStorage; it calls the bound `BuildCharacterLoreBlock`, which uses the same resolver for ownership, visibility, activation, deletion and stale roots | FIXED IN CODE — binding unit-tested; native UI path not yet exercised |


| W20 | Test-only inference recorder: `promptcapture` build tag wraps the real boundary, production is a compile-time no-op with zero capture symbols; `BuildCharacterLoreBlock` is in generated bindings and the composer fails loudly when the resolver is unavailable | IMPLEMENTED + NATIVE ARMING VERIFIED — GUI chat/character prompt exercise pending |


| W21 | Native GUI ordinary-chat capture through the armed recorder: 2 messages at the real boundary, probe present, forbidden marker absent, user turn + reply persisted to `history.json`; WebView2 exposes no UIA form controls, so the drive is coordinate-based at 1400x900 | NATIVE ORDINARY-CHAT CAPTURE VERIFIED — character/lorebook scenario pending |


| W22 | Staging the two-character native scenario exposed an authoring gap: the canonical Lorebooks editor never writes `ownerId` (and `Store.Save` replaces the document verbatim), so character-scope private lore created there is unowned and never routed; editing an imported owned entry strips its owner | OPEN — blocks native character acceptance; needs an owner control plus preserve-on-replace |


| W23 | W22 corrected: character-scope entries have an explicit Owner control distinct from Linked Characters; the editor sends `ownerId`; replace preserves it only when the payload omits the field and honors explicit set/clear; private unowned saves require confirmation | FIXED IN CODE — store contract unit-tested; UI round-trip and native scenario pending |


| W24 | Ownership identity is the display name and the selector was reading the on-disk character set while chat reads `CharacterStorage`; `lorebooks.html` now loads character storage and the selector prefers it, falling back to `GetCharacters`. Duplicate names and renames remain unresolved because the resolver compares names | FIXED IN CODE — identity-source mismatch; native round trip still pending |


| W25 | Native editor round-trip attempted: the corrected owner selector lists characters from `CharacterStorage` inside the native window, but coordinate driving could not reliably reach the Lorebooks page or complete the second character rename; no canonical owner value was exercised | INCOMPLETE — needs element-addressable drive (UI harness) before native confirmation |


| W26 | Identity-flow investigation and stable-ownership proposal: reuse `character.id`; editor persists IDs; resolver routes by character ID; frontend-supplied, backend-journaled, versioned name→ID migration; ambiguous legacy values excluded and reported | PHASE 1/2 COMPLETE — awaiting product decisions on migration authority, legacy chats without IDs, ownerName snapshot, memory scope |


| W27 | Stable character ownership implemented: `character.id` is the sole authorization identity; resolver/editor/prompt paths use IDs; optional `ownerName` is descriptive; versioned journaled name→ID migration with rollback; memory name-owned records fail closed at prompt; disk directory-name IDs documented as legacy | NATIVE WAILS ACCEPTANCE PASS — native editor/import/rename/prompt-boundary/unresolved-owner/restart scenario complete; no model inference performed |

| W27a | W27 acceptance exposed cross-owner conflict collapse: same-trigger character entries with different stable owners were resolved before routing. `conflictGroupKey` now separates character ownership domains; same-owner conflicts still resolve normally | FIXED + REGRESSION TESTS PASS — native Aria/Bram prompt blocks and inference-boundary prompts isolate correctly |


| W28 | Memory identity containment: explicit `owner_kind`; stable-ID capture; continuity export/import requires same character ID; prompt and snapshot filters fail closed for legacy/name-owned records; diagnostics + visible warning; legacy data preserved without reassignment | GO TESTS + FRONTEND MEMORY TEST + BUILDS PASS — native memory checks completed via W28b test seam; legacy ownership repair remains a product decision |

| W28a | W28 follow-up: generic export/import classification test proves a legacy record whose owner string equals the target ID is rejected; prompt-route audit finds no second memory composer; initial native UI attempt blocked by WebView2 automation surface | CODE + ROUND-TRIP TESTS PASS — native checks completed in W28b |

| W28b | Test-tagged loopback native seam: explicit opt-in plus marked isolation, loopback-only, production no-op; real Wails page verified warning, stable retrieval, character switching, same-name recreation, snapshot filter, restart persistence and legacy preservation; production binary contains no seam identifiers | 7/7 NATIVE TEST-SEAM CHECKS PASS — snapshot is load-level, not live-model prompt evidence |

| W28c | Security audit of W28b test seam: per-run token, Origin enforcement, fixed methods/paths/body caps, user-only Windows ACL, token rotation and graceful-shutdown closure; residual same-user process and page-binding reach documented | SECURITY TESTS PASS — test-only loopback boundary, not production security |

