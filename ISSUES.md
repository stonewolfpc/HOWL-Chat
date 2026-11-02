# HOWL Chat — Working Issue Document

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

### T0.2 · The model receives my message twice
**Suspected separate bug** from T0.1 — likely the message is appended to the
prompt twice even though it's saved once. Fixing T0.1 may not fix this. Test both.
**Evidence:** a prompt dump or server log showing the message once per turn.

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

### T0.4 · Image preview is broken
Long-standing. The generated image does not preview. Find the cause and repair it.
This is the last blocker before the image page can be judged on layout (T2.3).
**Evidence:** prompt → image appears in the preview panel, correct aspect, no crop.

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

### T1.2 · VRAM arbitration between chat and image models
Agreed behaviour:
- Image model requested while a chat model is loaded → check whether both fit.
  If not, a confirm dialog: this will unload your chat model to make room.
- On confirm: unload chat, re-check freed VRAM, then load image.
- Leaving the image page does **not** unload or disconnect.
- Going back to chat and pressing Load mirrors it: check image model, confirm,
  unload, load. Same dialog, wording swapped.
**Evidence:** both directions exercised, with VRAM readings before and after.

### T1.3 · Think tags / reasoning handling
Match how SillyTavern and Kobold handle think modes: on, or off, or hidden, and
selectable intensity where the model supports it. Hidden means it must not count
against the context budget. A thought must never be fed back to the model on the
next turn — anything inside a think tag is discarded from history.
**Evidence:** hidden mode shows no tokens in the UI and shortens the measured
prompt; a two-turn test proves no think block reappears in turn two.

### T1.4 · Sampler defaults the user never has to touch
Defaults stay "suggested general," but vary by model family — they differ too much
to force one set. Supported models must get correct prompt formatting
automatically. No card in the UI that the backend ignores.
**Evidence:** load two different families, confirm each gets its own input/output
format with no user action.

### T1.5 · Header says the wrong name
"Howl Chat - Tron Edition" → "Howl Chat". One-line fix, bundle it with T1.1.

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

### T2.2 · Ghost elements and stale markup
Class of bug: markup left behind by deleted features, and duplicate element IDs
where one handler writes to an ID that no longer exists while the real element sits
there working. We've already hit one of these (`ig-negative-prompt` vs `ig-negative`).
Sweep every element lookup target against the actual DOM.
**Evidence:** zero write-to-missing-ID, zero duplicate IDs, no orphaned markup.

### T2.3 · Image page layout — target Invoke AI organisation
Rework toward Invoke AI's structure: preview prominent, controls grouped, history
organised. Reference it (look it up) but don't chase parity — take the parts that
make it powerful and organised, drop the rest. Depends on T0.4.

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

### T3.2 · Load speed — learn from LM Studio
LM Studio loads a model in under a quarter of our time. Find out how and adapt
whatever is portable. This may be the single biggest quality-of-life win in the app.

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

### T3.4 · Per-field AI character sheets
Every field gets its own generate button. The model sees the fields already filled
and must make the new one line up: name, kind, the target field's purpose, plus
truncated versions of siblings as reference, and returns only that field.
Draft + accept / regenerate / lock per field. Grammar-constrained JSON, one
pipeline with per-kind schemas.
No preset voices. Exemplar dialogue is authored per character — exemplar lines beat
trait lists for keeping a voice from bleeding.

### T3.5 · Multi-character chat
2+ characters kept fully separate, each receiving the appropriate system prompt,
with context flushed properly between speakers. Root cause of voice bleed is flat
context, not sampling — needs structured transcript + per-character state.

### T3.6 · Ternary models as first-class
Load ternary formats without new dependencies. Must just work when the user points
at the file.

### T3.7 · Audio and TTS end to end
Test the audio model capability and the TTS service through the UI. Needs headphones
present for verification.

### T3.8 · Audio generation
Newer small audio model, reportedly fits in 8GB for short clips. Short clips are
fine. If it works, play the generated audio inside the app.

### T3.9 · Diffusion version check
Confirm the current Stable Diffusion release; update to the latest patched version
if one exists.

### T3.10 · LLM writes the image prompt — research gated
Can a local LLM prompt the image generator well? Would it be worth the cost?
Proposed flow if yes: compose the prompt before swapping models; optionally keep
enough of the LLM resident to inspect the result (only if vision-capable) and refine.
Every part of this is a big if — cost, quality, speed. Research before building.

### T3.11 · Compatibility for everyone else
Auto context and layering currently assume this machine. Generalise as far as
possible and ship as beta, then fix from user reports. Goal: if it runs AI at all,
it should be usable.

### T3.12 · Whole-app optimisation pass
Separate pass, its own stage, after everything above. Look for waste everywhere.

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
- **Multi-character with avatars.** Sub-case of the above; same resource problem.

---

## Research queue (do before committing to any of it)

| Topic | Question |
|---|---|
| LM Studio load speed | How does it load so fast, and what's portable to us? |
| Context flush | Was it ever implemented? What does it actually cost per turn? |
| KV cache reuse | Does `--cache-reuse` help on this hardware, measurably? |
| Think modes | How do SillyTavern/Kobold expose levels, and what do models accept? |
| Ternary formats | Which quantiser/runtime path, with no new dependencies? |
| Invoke AI | Which layout decisions are worth copying? |
| Stable Diffusion | Current release, and what updating costs us. |
| LLM → image prompt | Is the quality worth the swap cost? |
| Audio model | Exact model, VRAM, clip length on 8GB. |

---

## Open questions for you

1. Which chat-switch or tab-switch currently loses a model? I need the actual
   reproduction, not the general rule.
2. Where do your image models live now — the new drive path?
3. Light theme: is there a specific page that's the worst offender?
4. Beta scope: ship beta right after Tier 0–2, or after T3.5 (multi-character)?

---

## Decision log

Append here as stages close. One line each: date, stage, evidence, result.

- 2026-11-02 · Image generation settlement guard — source review PASS; stage-recorded tests PASS; independent review PASS; real-app acceptance PASS on verified build `d7e8aacd…` (PID 9252): abort leaves child alive, next request refused as abandoned, Unload/Reload recovers, work admitted after reset. Closed with no code change; the earlier live run had tested a pre-guard 18:59 binary. Unresolved: why the old sd-server exited.
- 2026-11-02 · T2.1 Stage C (remaining pages theme verification) — five pages measured in both themes through `setThemeMode` with transitions off: 0 low-contrast nodes and 0 stuck-dark areas in light, expected dark panels only in dark, filled buttons carry `--on-accent`. No code changed (nothing needed correcting). Corrections recorded: `real-settings.html` no longer exists (deleted by T2.2), and a "stuck light flyout" reading was retracted as a 600 ms-vs-4 s transition artifact. Open: override-sheet deletion, top-bar collision, `audio-settings` dark body background, and process reaping that trusts recorded state over parent identity.


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

### Stage closure ritual (applies to every stage)
When a stage's fixes are in, don't open the next one:
1. Walk the whole app, not just the stage's area — the way you did when building
   this document, one thing at a time.
2. Anything new goes in this document, appended to the correct tier. It is **not**
   fixed on the spot, and it does not change the current stage's scope.
3. Compare the fresh pass against the document: what's closed, what's still open,
   what's new.
4. Only then choose the next stage.

This is the loop that produced the good list today. Protect it.

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
   tab switch, chat switch, or re-render. A timer killed by a later lifecycle change
   explains "it used to work" precisely.
4. **Overrides.** If the earlier coder's UI overrides removed the container or its
   styles, this is the same hunt as T2.2. Run them in one session.

**Spec once repaired:**
- One shared image-set mechanism bound in three places (chat speaker, world panel,
  scenario panel). Do not write three timers — this is the same component three times.
- Default image is the anchor: it shows on load and when cycling is off.
- Optional cycling, off by default, with a sane interval (a few minutes, not seconds).
  Configurable duration; per-character override later if it earns it.
- Never cycle while the character isn't speaking in a group chat — wasted decoding.
- Preload the next image before swapping so there's no blank flash, and keep the
  container size fixed so the layout doesn't jump between aspect ratios.
- Clear the interval on unmount and on chat/tab change. Stacked intervals are the
  classic cause of "it speeds up the longer I leave it running."
- Respect reduced-motion and give it a visible off switch.
- Avatars (T2.4) should draw the default image from this same set. One source of truth.

**Evidence:** a card with 4 images, default shows first, cycling advances on the
timer, no flash, no layout jump, timer stops on tab switch and does not stack;
world and scenario panels behave identically.

### T3.13 · D&D / tabletop mode
**Gate:** requires T3.3 (typed state + extraction) and T3.5 (multi-character).
Not startable before those, because its whole problem is theirs.

**Why it keeps failing:** local models are bad at being the source of truth. They
drop clothing, invent inventory, lose HP, contradict the map. No prompt fixes that.
It's a memory architecture problem, not a phrasing problem.

**Core rule: the app owns state; the model narrates.**
- Typed state store: characters (outfit, condition, injuries), inventory (item,
  quantity, holder), location and time, quests/flags, combat state.
- A compact state block is injected every turn (~200–400 tokens). Facts are
  restated, never remembered. This is what stops the clothing drift after 3 messages.
- The model may only **propose** changes, as grammar-constrained JSON deltas
  (`add_item`, `remove_item`, `set_outfit`, `damage`). The app validates and applies
  them. Unparseable or invalid → ignored, state unchanged. Prose never mutates state.
- Dice, HP, initiative, and rules checks run in code. The model rolls nothing and
  calculates nothing; it narrates outcomes the app hands it.
- Rules text is retrieved from a scoped rules lorebook, not stuffed, so a rulebook
  can't eat the context window.

**Two modes:**
- *Model as DM* — AI narrates and adjudicates; the app holds state, rolls, enforces
  rules, feeds results back.
- *User as DM* — you narrate; app gives you the state panel, rolls, and NPC tools;
  AI plays the NPCs.

**The inventory panel is the point.** Visible, editable, always showing what's
carried and by whom. A state panel is what stops "we forgot the rope" drift — not a
better system prompt. This is the part that has to be right.

**Evidence:** 10 exchanges with no outfit drift; inventory changes only through
validated deltas; every roll traceable to an app-side result.

## Architectural law · App owns state, model narrates

The app holds every fact it can enumerate and verify. The model is handed those facts
each turn and returns prose and *proposed* changes. Nothing important is remembered by
the model; nothing important is inferred from vibes.

**Why:** local models are bad at being the source of truth. They hold facts in the same
soft attention that carries tone and personality, so state drifts within a few turns.
Restating a fact every turn costs a few tokens and never drifts.

### Authority ladder (resolve conflicts this way)
1. **What the user explicitly states** — highest. If you say you took off the coat, the
   app records the coat as off. Never argue with the user's own narration.
2. **A validated model delta** — structured, schema-checked, applied.
3. **App defaults** — used when nothing else says otherwise.
Rule: the app wins on facts it owns (item count, HP, seed, model loaded). The user wins
on intent (what they say just happened). The model wins on nothing that is checkable.

### Decision test
| Question | Verdict |
|---|---|
| Can the app enumerate and verify it? | App owns it. Model may propose deltas. |
| Is it prose, tone, description, improvisation? | Model owns it. Leave it alone. |
| Is it a judgement call (is the NPC convinced)? | Model narrates. If rules are on, the app adjudicates the outcome. |
| Is it relevance detection (is this lore relevant)? | App computes it, but **must show and allow override** — heuristics are not truth. |
| Does owning it risk contradicting the user's own words? | Don't own it. Ask or record what they said. |

### Where this already applies in the plan
| Area | Currently | Should be |
|---|---|---|
| Transcript identity / order | App re-renders from partial state (T0.1) | App owns the message list; UI renders it |
| Sampler values | Reset visually on restart (T1.1) | App owns and loads them at startup |
| Loaded model + VRAM | Lost on tab switch (T1.1) | App owns one load state, tabs observe it |
| Prompt format per family | Per-user guessing (T1.4) | App selects the template from the model |
| Think blocks | Risk of feeding reasoning back (T1.3) | App parses, strips, never echoes |
| Which lore fires | Model would decide by vibes | App evaluates triggers deterministically, with a visible override |
| Lore/memory budget | Unbounded | App enforces the ceiling and priority order |
| Character card fields | Model edits prose freely | App owns structure; model fills one field at a time |
| Image prompt/seed/dims/LoRA/board | Scattered | App owns; model may only propose prompt text |
| Outfit, inventory, HP, location | Model's memory (drifts in 3 turns) | Typed store injected each turn (T3.13) |
| Which image is default, cycle timer | Unclear / broken (T2.5) | App owns the default and the timer |
| Speakership / active character | Hardcoded assumption | App owns the speaker set per chat |

### Where NOT to apply it
- Do not let regex triggers silently outrank the model on meaning. Surface the call
  ("this entry fired because it saw 'Moonblade'"), let the user disable it. A hidden
  heuristic that is wrong is worse than a model guess, because nobody can see it.
- Do not freeze anything the user wants to hand-adjust mid-scene. Every app-owned
  field needs a manual edit path.
- Do not make the app author prose. Owning facts is not writing the story.
