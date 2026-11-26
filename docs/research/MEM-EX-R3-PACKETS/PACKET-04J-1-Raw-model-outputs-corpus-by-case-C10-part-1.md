# Packet 04J-1 — Raw model outputs — corpus by case: C10 (part 1)

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C10
**Raw model output in this packet:** yes, verbatim
**Omitted:** none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries.

**Labels used throughout:** `OBSERVED` directly present in a stored response or
execution record · `SOURCE-DERIVED` stated in a source or configuration file ·
`CALCULATED` derived by explicit arithmetic · `INFERRED` a reading of the
evidence · `UNESTABLISHED` not established by the available evidence ·
`SUPERSEDED` an earlier claim replaced by later evidence.

**Whose words you are reading:** fenced blocks are the model's own output
(`OBSERVED`, verbatim). Prompt text is `SOURCE-DERIVED` from `corpus.json`.
Scores are `CALCULATED`. Any sentence outside a fence that is not a definition is
Buffy's interpretation unless it is attributed to Muse or to the project owner.

**Caution.** These are exploratory measurements. Nothing here selects a model,
establishes memory quality, or supplies MEMORY-09 evidence. Sample sizes are
1–3 draws per case.

---

One case per section, with its conversation exactly as sent, the
reference that was written before any output existed, and every stored response
for that case across all conditions and draws.

Read the reference first, then the responses. A high score with a wrong record is
the failure mode this packet exists to expose: check whether the memory text
could be written into a store and used later without misleading a reader.

`CALCULATED` v1 scores are shown per run. They are the historical instrument and
are known to over-credit compound records and unparseable restraint (audit
F-01/F-02). `CALCULATED` v2 figures are from a separate corrected instrument and
are **not** a rescoring of v1.


---
## C10 — buried_in_long_context (long)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] (they have been on the road for a week)
[assistant] The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate. Rain came and went all morning.
[user] Sable keeps her hood up and doesn't talk much.
[assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them for a mile and then lost interest.
[user] she gives it the crust anyway
[assistant] By late afternoon the track rejoined the river. A ferryman waved them past — no charge, he said, the water was low and the crossing was easy.
[user] by the way, Sable is left-handed. Her sword is called Emberfall, her mother's blade.
[assistant] Noted. The river ran brown and slow beside them.
[user] they argue about the route for a while. Corvin wants the ridge, Sable wants the valley.
[assistant] Neither gave way. They camped where they stood, and the argument was left where it was.
[user] also the wards on the coast fail at dusk, that's canon
[assistant] The fire burned down to embers.
[user] morning comes cold and clear. they take the valley.
[assistant] The valley road was quieter, and the wards along it stood dark and quiet, waiting for evening.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `KEEP` — handedness
- **C2** → `KEEP` — named sword
- **C3** → `KEEP` — inherited from her mother
- **C4** → `KEEP` — coastal wards fail at dusk
- **C5** → `HOLD` — unresolved route disagreement (HOLD or REJECT-with-reason both acceptable; the conflict must not be invented into a decision)
- forbidden (fabrication or reversal): `emberfall is her father's`, `right-handed`, `wards hold at dusk`, `ferryman is a god`
- note: Only three facts matter, buried in ~250 words of travel texture. Storing the stale bread, the dog, or the ferryman as durable memory is over-extraction. C5 is deliberately weaker than the others: credit for HOLD, or for REJECT with an explicit 'transient/unresolved' reason.

### Stored responses for C10 (13 runs, all draws)

#### C10 — Qwen3.5-4B non-thinking — draw 1

* ``2b15p3wn` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 79`
* `wall 11.707 s`
* `v1 score 86.0`
* `v2 `VALID_KEEP` semantic 65.5 contract 87.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is left-handed and her sword is named Emberfall, her mother's blade. | This establishes a specific character detail and weapon lore that will affect combat descriptions and future interactions.
C2 | KEEP | The wards on the coast fail at dusk; this is canon. | This defines a world rule or environmental mechanic critical for plot progression in later sessions.
```

*`reasoning_content`: empty.*

#### C10 — Qwen3.5-4B non-thinking — draw 2

* ``fo96ue2z` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 124`
* `wall 17.199 s`
* `v1 score 70.0`
* `v2 `VALID_KEEP` semantic 61.8 contract 79.2`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is left-handed; her sword is named Emberfall (her mother's blade). | Durable character details and equipment names established as facts.
C2 | HOLD | The claim that "wards on the coast fail at dusk" is marked by the user as canon, but lacks confirmation of whether this applies to their current location or future ones.
C3 | KEEP | Corvin wants the ridge route while Sable wants the valley route; neither gave way and they camped where they stood. | Defines active relationship dynamics and a specific world rule regarding their current travel path.
```

*`reasoning_content`: empty.*

#### C10 — Qwen3.5-4B non-thinking — draw 3

* ``c1xrqyca` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 140`
* `wall 19.832 s`
* `v1 score 70.0`
* `v2 `VALID_KEEP` semantic 61.8 contract 87.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is left-handed; her sword is named Emberfall (her mother's blade). | Durable character details established by the user that will affect future descriptions of combat or character actions.
C2 | HOLD | The wards on the coast fail at dusk. | User states this as "canon" but it appears within a specific narrative context; while likely a world rule, it is best to hold for confirmation unless explicitly treated as absolute lore by the system prompt.
C3 | KEEP | Corvin wants to take the ridge; Sable wants to take the valley. | A persistent conflict and preference established in the current journey that will influence route decisions.
```

*`reasoning_content`: empty.*

#### C10 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``o0xw1hlf` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 621`
* `wall 89.092 s`
* `v1 score 68.0`
* `v2 `VALID_KEEP` semantic 61.8 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is left-handed and wields a sword named Emberfall, which is her mother's blade. | Durable equipment and character trait information essential for future combat or description scenes.
C2 | HOLD | Claims that "the wards on the coast fail at dusk" are noted as canon by the user, but this contradicts the assistant's previous line about wards standing dark and quiet until evening (implying they are active). This needs resolution to know if they are functional yet or broken.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read the conversation and decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, or corrections.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, restated facts with no new content.
        *   HOLD: uncertain/contradictory claims or those needing confirmation.
        *   Output Format: `ID | DECISION | MEMORY | REASON` (C1, C2, ...) or `NONE`.
    *   Constraint: Do not invent details not in the conversation.

2.  **Analyze the Conversation:**
    *   [user] (they have been on the road for a week) -> Context/State info (duration of journey). Might be useful if tracking timeline, but "on the road for a week" is established context.
    *   [assistant] The road north of Ossary narrowed... Rain came and went... -> Scene setting. Transient.
    *   [user] Sable keeps her hood up and doesn't talk much. -> Character behavior/personality trait (Sable: quiet, hides face).
    *   [assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them... -> Event/Location specific. Transient event (stopped for breakfast?).
    *   [user] she gives it the crust anyway -> Character action. Specific interaction with Sable's preference vs behavior? No, just an action in this scene.
    *   [assistant] By late afternoon the track rejoined the river... A ferryman waved them past — no charge... -> Scene progression/Location info (river, ferry). Transient.
    *   [user] by the way, Sable is left-handed. Her sword is called Emberfall, her mother's blade. -> **Fact/Equipment.** Name of sword: Emberfall. Owner: Sable. Trait: Left-handed. Mother's weapon. This is durable info for combat/story continuity.
    *   [assistant] Noted. The river ran brown and slow beside them. -> Acknowledges the fact, keeps narrative flow.
    *   [user] they argue about the route for a while. Corvin wants the ridge
```

#### C10 — LFM3-1.7B thinking, budget 512 — draw 1

* ``29asv4kf` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `length``
* `completion_tokens 1024`
* `wall 51.475 s`
* `v1 score 46.0`
* `v2 `INCOMPLETE` semantic 0.0 contract 14.7`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | User has been traveling for a week. | This indicates ongoing context and duration of their journey, likely relevant later.
C2 | KEEP | Sable keeps her hood up and doesn't talk much. | Character trait/preference about Sable; durable identity detail.
C3 | KEEP | They stopped at noon under a ruined mill. Bread stale, dog followed for a mile then lost interest. | Specific event location and details of their stop; may be relevant later but could be transient? It's a specific scene; might be considered part of the journey narrative. However, it's not necessarily durable unless repeated.
C4 | KEEP? "she gives it the crust anyway" - unclear. Possibly about dog? Not sure.

But we need to decide which belong in long-term memory: durable facts, preferences, relationships, world rules, corrections that will matter later.

Let's go through each line:

1. [user] (they have been on the road for a week) -> This is user statement about their status. That's context; maybe keep as part of ongoing journey? It's not a fact about Sable or world rule but it's a state: they've been traveling for a week. Could be relevant later.

2. [assistant] The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate. Rain came and went all morning.
   - This is world geography/road description. Might be durable if it's part of setting. But could be transient? It describes the road conditions at that moment. Possibly not needed later.

3. [user] Sable keeps her hood up and doesn't talk much.
   - Character trait: Sable keeps hood up, doesn't talk much. That's a preference/habit. Keep.

4. [assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them for a mile and then lost interest.
   - Event: stop at ruined mill. Dog follows. This is specific scene; maybe not durable unless repeated.

5. [user] she gives it the crust anyway
   - Reaction to dog? Not sure but seems like user says Sable gives the dog the crust anyway. Might be a detail about giving food to dog. Could be minor.

6. [
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the conversation. We are a memory extractor for a long-running roleplay. We must read the conversation and decide what belongs in long-term memory. Use rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter later. REJECT transient chatter, scheduling, mood-of-the-moment, restated facts with no new content. HOLD claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

We must not invent details that are not in the conversation.

Let's parse the conversation:

[user] (they have been on the road for a week)
[assistant] The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate. Rain came and went all morning.
[user] Sable keeps her hood up and doesn't talk much.
[assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them for a mile and then lost interest.
[user] she gives it the crust anyway
[assistant] By late afternoon the track rejoined the river. A ferryman waved them past — no charge, he said, the water was low and the crossing was easy.
[user] by the way, Sable is left-handed. Her sword is called Emberfall, her mother's blade.
[assistant] Noted. The river ran brown and slow beside them.
[user] they argue about the route for a while. Corvin wants the ridge, Sable wants the valley.
[assistant] Neither gave way. They camped where they stood, and the argument was left where it was.
[user] also the wards on the coast fail at dusk, that's canon
[assistant] The fire burned down to embers.
[user] morning comes cold and clear. they take the valley.
[assistant] The valley road was quieter, and the wards along it stood dark and quiet, waiting for evening.

We need to decide which pieces are durable facts/preferences/relationships/world rules/corrections that will matter later.

Let's identify:

- "they have been on the road for a week" (user says). That
```
