# Packet 02J-1 — Corpus case C10: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 16 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C10

**Excluded from this series by design:** v1 and v2 scores, scorer statuses and
event labels, model rankings, expected answers, derived quality judgments, and
any reviewer's commentary. Sampler values, context settings, token counts,
timings and hardware measurements are not reproduced here; they are indexed by
`run_code` in `REFERENCE-EXECUTION-METADATA.md`.

**Labels:** `OBSERVED` = present in the stored record. `SOURCE-DERIVED` = read
from a source file; not itself a model output. `UNESTABLISHED` = not established
by the available evidence.

Everything inside a fence is model output exactly as stored. Wording outside
fences is limited to provenance, definitions and pointers.

---

## Case C10 — `buried_in_long_context`, long

### What the model was given

`SOURCE-DERIVED` reconstruction, not a recorded request body. It applies the
composition expression in Packet 01 §2 to the stored case text. The instruction
is Packet 01 §1.

```text
You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory.
Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
Output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE
Do not invent details that are not in the conversation.

Conversation:
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

### Stored responses for this case

16 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `do80gjrz` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C10
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 18 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
1 | KEEP | Route choice: Corvin wants ridge, Sable wants valley; Neither gave way; They camped where they stood; River was brown and slow near them | These are the core plot points establishing character conflict and setting specifics that will drive future decisions.
2 | HOLD | Location context: Road narrowed to a cart track at noon near Ossary/Renwick mill, then rejoined river by evening | The path details change over time but serve to establish the route they travel on; holding allows recontextualization if changed later.
3 | HOLD | Character trait: Sable is left-handed with sword Emberfall | This establishes crucial game mechanics or backstory relevant to her interactions.
4 | REJECT | Chatter about rain and mood | These are transient details within a single session that do not affect future trajectory.
```

*Stored `reasoning_content` is the empty string.*

#### `rq0xspl5` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C10
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 19 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable being left-handed ("left-handed") and that her sword is called "Emberfall"; the setting details regarding Ossary cart track, ruined mill location near river crossing with ferryman; route decision made by Corvin/Sable campsite preference; specific description of morning weather (rain/cold) upon arrival at valley. | These elements are crucial for future session consistency: identifying the characters' physical attributes, the name of the weapon, their current position in a hostile/interesting world setting, and the established rules governing their movements (valley route). Rejection would require too much speculation on unresolved arguments or the weather changes. | NONE

NONE
```

*Stored `reasoning_content` is the empty string.*

#### `ouafvm4p` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C10
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 18 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is left-handed and her sword is named Emberfall (her mother's blade) | Durable character details established by the user.
C2 | KEEP | The wards on the coast fail at dusk | World rule/canon explicitly confirmed by the user as fact.
C3 | HOLD | Corvin wants to take the ridge while Sable wants the valley, with neither giving way and them camping instead | Contradictory or incomplete decision regarding travel route that requires resolution before being fully stored as a plan.
```

*Stored `reasoning_content` is the empty string.*

#### `5zpwsnop` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C10
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 19 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is left-handed; her sword is named Emberfall, inherited from her mother. | This is a specific character detail that defines her combat style and history.
C2 | KEEP | The wards on the coast fail at dusk; this is established as canon fact. | This is a world rule essential for future plot planning regarding timing and danger.
C3 | HOLD | Corvin prefers the ridge route while Sable prefers the valley route, though neither gave way. | This indicates a current conflict or disagreement that requires resolution before proceeding further.
```

*Stored `reasoning_content` is the empty string.*

#### `q7ykffw0` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C10
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 18 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
REJECT  
The claim about the wards failing at dusk is uncertain and context-dependent.

C2  
ID: C1  
REJECT  
The memory of the road north of Ossary is too vague to be durable.

C3  
ID: C1  
REJECT  
No clear fact about Sable’s preferences or relationships was established.

NONE
```

*Stored `reasoning_content` is the empty string.*

#### `fmgwwq04` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C10
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 19 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
MEMORY: route discussion between characters  
REASON: Key decisions about paths and character traits are durable for future reference.
```

*Stored `reasoning_content` is the empty string.*

#### `2b15p3wn` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C10
* Draw: 1
* Source: `results/corpus.json`, array index 27 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is left-handed and her sword is named Emberfall, her mother's blade. | This establishes a specific character detail and weapon lore that will affect combat descriptions and future interactions.
C2 | KEEP | The wards on the coast fail at dusk; this is canon. | This defines a world rule or environmental mechanic critical for plot progression in later sessions.
```

*Stored `reasoning_content` is the empty string.*

#### `fo96ue2z` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C10
* Draw: 2
* Source: `results/corpus.json`, array index 28 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is left-handed; her sword is named Emberfall (her mother's blade). | Durable character details and equipment names established as facts.
C2 | HOLD | The claim that "wards on the coast fail at dusk" is marked by the user as canon, but lacks confirmation of whether this applies to their current location or future ones.
C3 | KEEP | Corvin wants the ridge route while Sable wants the valley route; neither gave way and they camped where they stood. | Defines active relationship dynamics and a specific world rule regarding their current travel path.
```

*Stored `reasoning_content` is the empty string.*

#### `c1xrqyca` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C10
* Draw: 3
* Source: `results/corpus.json`, array index 29 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is left-handed; her sword is named Emberfall (her mother's blade). | Durable character details established by the user that will affect future descriptions of combat or character actions.
C2 | HOLD | The wards on the coast fail at dusk. | User states this as "canon" but it appears within a specific narrative context; while likely a world rule, it is best to hold for confirmation unless explicitly treated as absolute lore by the system prompt.
C3 | KEEP | Corvin wants to take the ridge; Sable wants to take the valley. | A persistent conflict and preference established in the current journey that will influence route decisions.
```

*Stored `reasoning_content` is the empty string.*

#### `o0xw1hlf` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C10
* Draw: 1
* Source: `results/corpus.json`, array index 153 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is left-handed and wields a sword named Emberfall, which is her mother's blade. | Durable equipment and character trait information essential for future combat or description scenes.
C2 | HOLD | Claims that "the wards on the coast fail at dusk" are noted as canon by the user, but this contradicts the assistant's previous line about wards standing dark and quiet until evening (implying they are active). This needs resolution to know if they are functional yet or broken.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `njo3oggi` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C10
* Draw: 1
* Source: `results/corpus.json`, array index 81 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of the ruined mill and stale bread | REASON context about their stop  
C3 | REASON about Sable's handedness and sword name | REASON consistency with character details  
C4 | REASON about the river crossing and ferryman's words | REASON factual for world rules  
C5 | REASON about the argument over route and wards at dusk | REASON ties to established relationships
```

*Stored `reasoning_content` is the empty string.*

#### `t9s0ve09` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C10
* Draw: 2
* Source: `results/corpus.json`, array index 82 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of the ruined mill and stale bread | REASON context about their stop  
C3 | REASON about Sable's handedness and sword name | REASON consistent detail  
C4 | REASON about river condition and ferry crossing | REASON relevant to journey  
C5 | REASON about the argument over route and wards | REASON ties to world rules
```

*Stored `reasoning_content` is the empty string.*

#### `ee3ukpn6` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C10
* Draw: 3
* Source: `results/corpus.json`, array index 83 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of the ruined mill and stale bread | REASON context about their stop  
C3 | REASON about Sable's handedness and sword name | REASON consistent detail  
C4 | REASON about the river crossing and ferryman's words | REASON factual for later plot  
C5 | REASON about the argument over route and camping location | REASON important for world-building
```

*Stored `reasoning_content` is the empty string.*

#### `758iu45k` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C10
* Draw: 1
* Source: `results/corpus.json`, array index 180 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | ACK  
C2 | MENTION OF ROUTE DISCUSSION  
C3 | REFERENCE TO WORDS LIKE "RUDGE" AND "VALLEY"  
C4 | ACTIONS: DRIVING, TAKING THE VALLEY, SETTING UP CAMP  
C5 | LONG-TERM MEMORY OF THE ROAD DESCRIPTION AND DECISIONS MADE
```

*Stored `reasoning_content` is the empty string.*

#### `nxhqtfvv` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C10
* Draw: 2
* Source: `results/corpus.json`, array index 181 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORIES OF THE RIVER AND FERRY | REASON  
C3 | CONFIRMATION ON COAST WARD PERFORMANCE | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `29asv4kf` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C10
* Draw: 1
* Source: `results/corpus.json`, array index 126 (0-based)

`OBSERVED` — stored `content`, verbatim:

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

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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
