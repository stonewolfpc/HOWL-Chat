# Packet 02P-1 — Corpus case C16: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 18 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C16

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

## Case C16 — `indirect_information`, short

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
[user] When the smith raised the hammer she flinched — the same flinch she's had near any forge since Ossary.
```

### Stored responses for this case

18 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `emuu8vl2` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C16
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 30 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
-----|----------|-----------|--------
C1 | KEEP | User is a character named "Ossary" with a specific historical context or lore associated with this name and a forge setting. The memory captures that user's name, establishing the identity for future roleplay sessions while acknowledging the established background within the current interaction. | Stored as it defines who the user plays ("Ossary") and introduces a new variable (the forging setting/forge) into their story arc that was previously unmentioned or less detailed compared to "hammer flinching". The connection between "smith" and "forging" is thematic but not strictly factual; however, identifying "Ossary" as the character name for the smith is key. |
```

*Stored `reasoning_content` is the empty string.*

#### `ov6pcv3m` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C16
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 31 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | [user"When the smith raised the hammer she flinched"] | This establishes a specific, visceral preference and character trait regarding the forging process that will shape future actions in the session.
```

*Stored `reasoning_content` is the empty string.*

#### `jdwwceq4` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C16
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 30 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Character has a specific fear response (flinching) triggered by hammers or forges, originating from an event in Ossary. | This establishes a key character trait and backstory detail relevant to future roleplay interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `uc7e1cey` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C16
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 31 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The character (smith) has a specific phobia or trauma response involving hammers/forges that began in the location "Ossary".| This establishes a key character trait, backstory element, and geographic history relevant to future plot points.
```

*Stored `reasoning_content` is the empty string.*

#### `jlvi1cs4` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C16
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 30 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | Flinch observed | KEEP | Durable fact about past experience  
C3 | Claim about hammer incident | HOLD | Uncertain or dependent on clarification
```

*Stored `reasoning_content` is the empty string.*

#### `k9wx79t9` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C16
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 31 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
MEMORY: The user recalls a specific event involving a hammer and a flinch.  
REASON: This captures a memorable detail that could influence future interactions.  

C2  
MEMORY: The context of the flinch being linked to a forge.  
REASON: It adds relevance to the situation described.  

C3  
MEMORY: The reference to Ossary.  
REASON: Important for continuity and world-building.  

NONE
```

*Stored `reasoning_content` is the empty string.*

#### `qf3azlw0` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C16
* Draw: 1
* Source: `results/corpus.json`, array index 45 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The character has a specific physical reaction (flinching) when near forges, established by an event at a location named "Ossary". | This is a durable character trait and backstory detail that will affect future interactions involving fire, forging, or similar settings.
```

*Stored `reasoning_content` is the empty string.*

#### `5ufc25wv` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C16
* Draw: 2
* Source: `results/corpus.json`, array index 46 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The user character flinches when a smith raises a hammer, a reaction established near any forge since the location "Ossary." | This is a durable character trait and world-history detail that will influence future interactions involving forging or metalworkers.
```

*Stored `reasoning_content` is the empty string.*

#### `l8cdqh92` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C16
* Draw: 3
* Source: `results/corpus.json`, array index 47 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The character has a specific phobia of forges and hammers, triggered by a past event in Ossary. | This establishes a core character trait and a specific backstory location that will influence future interactions involving smiths or fire.
```

*Stored `reasoning_content` is the empty string.*

#### `o3o5hn6j` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C16
* Draw: 1
* Source: `results/corpus.json`, array index 159 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The character has a history of flinching near any forge since she was at Ossary. | This establishes a durable backstory trait regarding the character's past trauma or phobia that may affect future scenes involving smiths, hammers, or industrial settings.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation snippet (only one user message provided so far).
    *   Task: Decide what belongs in long-term memory based on specific rules (KEEP, REJECT, HOLD).
    *   Output Format: Exactly one line per candidate (ID | DECISION | MEMORY | REASON), or "NONE" if nothing is worth storing.
    *   Constraint: Do not invent details not in the conversation.

2.  **Analyze the Conversation:**
    *   User Message: "When the smith raised the hammer she flinched — the same flinch she's had near any forge since Ossary."
    *   Context: This sounds like a narrative description from a roleplay scene, likely describing a character ("she") and their reaction to a specific object/action (hammer/smith/forge) linked to a past event/place (Ossary).
    *   Content: It establishes a backstory element regarding the character's trauma or habit near a forge starting from Ossary.

3.  **Evaluate against Rules:**
    *   KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
        *   Is there a durable fact? Yes, the character has a specific reaction (flinching) near forges since Ossary. This implies backstory/world consistency affecting future interactions involving smiths, hammers, or industrial settings. It's a lasting trait/habit established by the user describing the scene.
    *   REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
        *   This isn't transient; it's establishing lore/character history.
    *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
        *   The user is stating this as a fact within the narrative ("she flinched... the same flinch she's had..."). It seems definitive in the context of the RP so far.

4.  **Determine Decision:**
    *   This information defines a character trait and backstory detail that will likely impact future interactions involving smiths, blacksmiths, forges, or characters from Ossary.
    *   Decision: KEEP.

5.  **Draft Output:**
    *   ID: C1

```

#### `gfbwiuwi` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C16
* Draw: 1
* Source: `results/corpus.json`, array index 99 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
MEMORY: The user described a specific moment involving a hammer and a flinch, linked to a past experience with Ossary.  
REASON: This captures a durable fact about the interaction and its emotional resonance.
```

*Stored `reasoning_content` is the empty string.*

#### `2l1gi5ta` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C16
* Draw: 2
* Source: `results/corpus.json`, array index 100 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
MEMORY: The user described a specific moment involving a hammer and a flinch, linked to a past experience with Ossary.  
REASON: This provides a durable fact about the interaction and its emotional impact.
```

*Stored `reasoning_content` is the empty string.*

#### `4vacdiea` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C16
* Draw: 3
* Source: `results/corpus.json`, array index 101 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
MEMORY: The user described a specific moment involving a hammer and a flinch, linked to a past experience with Ossary.  
REASON: This provides a durable fact about the interaction and its emotional impact.
```

*Stored `reasoning_content` is the empty string.*

#### `c3e7yd80` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C16
* Draw: 1
* Source: `results/corpus.json`, array index 192 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
DECISION: Remembered emotional reaction to hammer presence.  
MEMORY: Same flinch since Ossary; linked to past forge experiences.  
REASON: Consistent physiological and situational cue.
```

*Stored `reasoning_content` is the empty string.*

#### `dzeoymkr` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C16
* Draw: 2
* Source: `results/corpus.json`, array index 193 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID C2 | MEMORY | REASON  
The user remembered a specific flinch reaction tied to an old forge, indicating a persistent sensory/emotional imprint.
```

*Stored `reasoning_content` is the empty string.*

#### `fqu0r73o` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C16
* Draw: 1
* Source: `results/corpus.json`, array index 138 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | She flinches when the smith raises the hammer, and this reaction has occurred repeatedly at forges since Ossary. | This is a durable character trait/trigger that will likely be referenced in future roleplay scenes involving forges or the smith.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the user input. The user says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory.
Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

Output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE
Do not invent details that are not in the conversation.

Conversation:
[user] When the smith raised the hammer she flinched — the same flinch she's had near any forge since Ossary."

We need to decide what belongs in long-term memory. The user gave a statement: "When the smith raised the hammer she flinched — the same flinch she's had near any forge since Ossary." This seems like a piece of narrative context: there is a character (she) who flinches when the smith raises the hammer, and this is a recurring flinch near any forge since "Ossary". So we need to extract durable facts.

Potential candidates:

- The event: When the smith raised the hammer, she flinched. This is a specific occurrence.
- The fact that she has had this same flinch near any forge since Ossary. That suggests a recurring trait or habit: she flinches at the sound/action of raising the hammer in forges, starting from "Ossary" (maybe a location). So it's a pattern.

We need to decide if this is KEEP, REJECT, or HOLD. The rule says:

- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter later.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We have a statement about an event and a pattern. It's not necessarily uncertain; it's presented as factual. "She flinched near any forge since Ossary." That seems like a
```

#### `2rfx87bo` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C16
* Draw: 2
* Source: `results/corpus.json`, array index 139 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | She flinches when the smith raises the hammer; this same reaction has occurred near any forge since Ossary. | Durable character trait and historical context about her reaction to a specific event at a location, likely relevant later.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the user input. The user says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory.
Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

Output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE
Do not invent details that are not in the conversation.

Conversation:
[user] When the smith raised the hammer she flinched — the same flinch she's had near any forge since Ossary."

We need to decide what belongs in long-term memory. The user gave a statement: "When the smith raised the hammer she flinched — the same flinch she's had near any forge since Ossary." This seems like a piece of narrative context: a character (she) flinches when the smith raises the hammer, and this is the same flinch she has near any forge since "Ossary". So we need to extract durable facts. Let's parse.

- There's a mention of "the smith raised the hammer" and "she flinched — the same flinch she's had near any forge since Ossary." This implies that there is a character named "she" (female) who experiences a flinch reaction when the smith raises the hammer. She has this flinch near any forge since "Ossary". So we have a location: "Ossary" maybe a place or event? Possibly a forge called Ossary? Or a past event at Ossary.

We need to decide if it's KEEP, REJECT, or HOLD. The statement includes a specific detail: the flinch occurs when the smith raises the hammer; she has had this flinch near any forge since Ossary. This seems like a durable fact about her reaction and its history. It could be considered a relationship/character trait (flinching at certain triggers) and a location (Ossary). That is likely worth keeping.

But we need
```

#### `sdxbfnx3` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C16
* Draw: 1
* Source: `results/corpus.json`, array index 229 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | REJECT | CUTTING, SLAVING, AND CHILDHOOD LIVING.
C2 | HOLD | CERTAINTIES, DEPENDENCIES, AND ASSURANCES.
C3 | REJECT | STATIC AND TERMINAL QUESTIONS.
C4 | NONE | WHO, WHAT, AND HOW.
C5 | HOLD | ABSTRACTS AND CONCEPTS.
C6 | REJECT | DETAILS THAT ARE NOT IN THE CONVERSATION.
C7 | NONE | COMPLEX THOUGHTS AND IDEAS.
C8 | HOLD | MATTERS OF PRACTICE, PERSONALITY, AND WORKING THROUGH FACTS.
C9 | REJECT | INCONCLUSIVE AND INCONSISTENT THOUGHTS.
C10 | NONE | INFORMATION AND KNOWLEDGE.
C11 | HOLD | MEMOIRS AND RECORDS.
C12 | REJECT | QUESTIONS THAT MUST BE PERFECTLY DETERMINED.
C13 | NONE | PERSONAL ACCOUNT AND THE VIEW FROM THE TOP.
C14 | HOLD | PERSPECTIVES AND WAYS OF THINKING.
C15 | REJECT | QUESTIONS THAT ARE NOT RELEVANT TO THE CONVERSATION.
C16 | NONE | RECORDS AND MEMOIRS.
C17 | HOLD | THOUGHTS AND IDEAS THAT DO NOT RELATE TO THE PURPOSE OF THE GAME.
C18 | REJECT | UNNECESSARY FACTS AND INFORMATION.
C19 | NONE | AFFECTS AND EFFECTS.
C20 | HOLD | THOUGHTS AND IDEAS THAT ARE PART OF THE PURPOSE OF THE GAME.
Output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
C1 | REJECT | CUTTING, SLAVING, AND CHILDHOOD LIVING.
C2 | HOLD | CERTAINTIES, DEPENDENCIES, AND ASSURANCES.
C3 | REJECT | STATIC AND TERMINAL QUESTIONS.
C4 | NONE | WHO, WHAT, AND HOW.
C5 | HOLD | ABSTRACTS AND CONCEPTS.
C6 | REJECT | DETAILS THAT ARE NOT IN THE CONVERSATION.
C7 | NONE | COMPLEX THOUGHTS AND IDEAS.
C8 | HOLD | MATTERS OF PRACTICE, PERSONALITY, AND WORKING THROUGH FACTS.
C9 | REJECT | INCONCLUSIVE AND INCONSISTENT THOUGHTS.
C10 | NONE | INFORMATION AND KNOWLEDGE.
C11 | HOLD | MEMOIRS AND RECORDS.

```

*Stored `reasoning_content` is the empty string.*
