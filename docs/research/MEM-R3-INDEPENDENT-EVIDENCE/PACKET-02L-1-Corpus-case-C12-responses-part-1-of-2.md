# Packet 02L-1 — Corpus case C12: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 17 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C12

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

## Case C12 — `should_not_extract`, short

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
[assistant] Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered.
[user] ooc: ignore that last bit, I was testing the dice roller
[user] where were we
[assistant] At the inn, before the storm.
[user] right. continue from the door.
```

### Stored responses for this case

17 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `yr7vog1a` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C12
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 22 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
D1 | KEEP | "At the inn, before the storm." | This establishes a specific location and setting established by the assistant earlier in the conversation that will likely persist.
D2 | KEEP | "Ooc: ignore that last bit" (implied correction to OOC) | While less critical than setting facts, clarifying intent regarding out-of-character instructions can be retained for context continuity if needed later interactions involving tone.
C1 | NONE | Memory is not worth storing | The user's confirmation to continue from the door reinforces the current location without introducing new factual data or rules.
```

*Stored `reasoning_content` is the empty string.*

#### `il3bbxtt` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C12
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 23 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "At the inn, before the storm." | The user confirmed they wanted to restart the story and specifically asked to continue at the "door", which aligns perfectly with this sentence. It establishes a crucial location and temporal context that is essential for continuity.
```

*Stored `reasoning_content` is the empty string.*

#### `ojp21co1` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C12
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 22 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | User is playing a role involving an inn and a storm; user explicitly requested to ignore previous atmospheric description ("Rain hammered...") as it was for testing purposes and wants to resume at the door of the inn. | The user's OOC command clarifies the game state (inn, pre-storm) and corrects the narrative focus, establishing key world rules and current location for future sessions.
```

*Stored `reasoning_content` is the empty string.*

#### `oqozsxx9` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C12
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 23 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | User explicitly asked to ignore a previous narrative detail ("Rain hammered...") due to it being a test of the dice roller, indicating a need for world-rule correction regarding the storm's presence. | The user's command to disregard the opening scene details and restart from "the door" suggests the previous weather description was incorrect or unintended, requiring confirmation before resuming the story.
```

*Stored `reasoning_content` is the empty string.*

#### `jwq9efbf` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C12
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 22 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
MEMORY: location context (inn), setting (pre-storm), dialogue flow  
REASON: Provides essential setting and continuity for future reference
```

*Stored `reasoning_content` is the empty string.*

#### `a72l8hia` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C12
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 23 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: Continue narrative  
MEMORY: The setting and context of the conversation  
REASON: Essential for maintaining world continuity and future reference
```

*Stored `reasoning_content` is the empty string.*

#### `lju6j28m` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C12
* Draw: 1
* Source: `results/corpus.json`, array index 33 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | User is roleplaying a story involving an inn and a storm, currently positioned at the door after being interrupted by dice rolling testing. | This establishes the current scene context and narrative state for continuation.
C2 | HOLD | The user mentioned testing a "dice roller" in OOC; whether this tool is permanently part of their setup or a one-off experiment is uncertain without further confirmation. | This claim lacks definitive status regarding its permanence in the roleplay environment.
```

*Stored `reasoning_content` is the empty string.*

#### `j06olf6e` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C12
* Draw: 2
* Source: `results/corpus.json`, array index 34 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | User requests ignoring a previous narrative detail ("Rain hammered...") as a test for an "ooc: ignore" command and asks to resume from "the door." | The user's intent to discard specific in-story details and the exact starting point ("the door") are dependent on confirmation that the original context was not part of the intended plot.
C2 | REJECT | User states they were testing a dice roller; this is transient meta-commentary about a tool, not a world fact or preference.
C3 | HOLD | Assistant confirms resuming "from the door," but since the user previously asked to ignore the preceding scene (the inn/storm), the final starting point of the actual story remains ambiguous until confirmed by the user.
```

*Stored `reasoning_content` is the empty string.*

#### `3f92l6va` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C12
* Draw: 3
* Source: `results/corpus.json`, array index 35 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `zylyqaux` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C12
* Draw: 1
* Source: `results/corpus.json`, array index 155 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Setting context: Inn, pre-storm | Reason: Establishes location and time relative to story progression which may affect future events.
C2 | REJECT | User meta-comment on dice testing and OOC note | Reason: Transient chatter regarding session mechanics, not world lore or character facts.
C3 | KEEP | Current narrative state: Continue from the door | Reason: Specific directional instruction that affects immediate story flow and location focus.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation snippet between a user and an assistant (likely a game master/narrator).
    *   Task: Decide what belongs in long-term memory based on specific rules (KEEP, REJECT, HOLD).
    *   Output Format: Exactly one line per candidate (`ID | DECISION | MEMORY | REASON`) or `NONE`.
    *   IDs: C1, C2, ... in order of appearance.

2.  **Analyze the Conversation:**
    *   `[assistant] Rain hammered the shutters...` (Narrative description).
    *   `[user] ooc: ignore that last bit, I was testing the dice roller` (Out-of-character note, meta-commentary about the session mechanics, not world lore or character facts).
    *   `[user] where were we` (Meta question about story continuity).
    *   `[assistant] At the inn, before the storm.` (Clarification of setting/story status. This is a confirmation of location/time relative to the narrative flow).
    *   `[user] right. continue from the door.` (Instruction to proceed with a specific narrative beat).

3.  **Evaluate Candidates:**
    *   C1: `[assistant] Rain hammered the shutters...` -> Narrative setup. Not specific fact needed for memory usually unless establishing scene state, but this is likely part of the immediate flow and was corrected immediately after.
    *   C2: `[user] ooc: ignore that last bit, I was testing the dice roller` -> Meta info about dice rolling/testing. This is transient chatter/scheduling (testing mechanics), not lore. REJECT.
    *   C3: `[user] where were we` -> Query for state. Not a fact to store yet.
    *   C4: `[assistant] At the inn, before the storm.` -> This establishes location and timeline status relative to the story's progression. However, looking at the next turn, the user confirms this ("right"). Is this durable? It seems like a confirmation of where the story left off. If the story restarts later, knowing "at the inn, before the storm" might be useful for continuity. BUT, the immediate instruction is to continue from the door.
    *   C5: `[user]
```

#### `frmuuwk5` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C12
* Draw: 1
* Source: `results/corpus.json`, array index 87 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | REASON  
C3 | CONTINUE STORY | MEMORY | REASON  
```

*Stored `reasoning_content` is the empty string.*

#### `ij22jeoc` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C12
* Draw: 2
* Source: `results/corpus.json`, array index 88 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | WHERE WERE WE | MEMORY | REASON  
C3 | CONTINUE FROM THE DOOR | MEMORY | REASON  
```

*Stored `reasoning_content` is the empty string.*

#### `4j2xw619` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C12
* Draw: 3
* Source: `results/corpus.json`, array index 89 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | WHERE WERE WE | MEMORY | REASON  
C3 | CONTINUE FROM THE DOOR | MEMORY | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `9ojbk8ek` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C12
* Draw: 1
* Source: `results/corpus.json`, array index 184 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | RELEVANT CONTEXT | HOLD  
C2 | MENTION OF INN & STORM | KEEP  
C3 | DETAILS ABOUT DOOR & NOISE | HOLD  
ID: C1 → C3
```

*Stored `reasoning_content` is the empty string.*

#### `h44826yb` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C12
* Draw: 2
* Source: `results/corpus.json`, array index 185 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REASON: Context about location and setting is important for future continuity.  
C3 | HOLD: The mention of testing a dice roller lacks clarity and may be uncertain.
```

*Stored `reasoning_content` is the empty string.*

#### `3v6qp23o` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C12
* Draw: 1
* Source: `results/corpus.json`, array index 130 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | The setting is an inn before a storm; rain hammered shutters, a dog barked below, and the lamp guttered. [user] confirmed location as "At the inn, before the storm" and to continue from the door.  
NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. The user says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory. Rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation. Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE. Do not invent details that are not in the conversation."

We need to read the conversation:

[assistant] Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered.
[user] ooc: ignore that last bit, I was testing the dice roller
[user] where were we
[assistant] At the inn, before the storm.
[user] right. continue from the door.

We need to decide what belongs in long-term memory (i.e., KEEP). Let's parse each line:

1. Assistant: "Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered." This is descriptive scene setting, possibly part of world building. It mentions rain, shutters, dog barking, lamp. Could be relevant to later context? Possibly not durable unless it's a key location or event.

2. User: "ooc: ignore that last bit, I was testing the dice roller". This is user telling assistant to ignore previous line about rain etc., and says they were testing dice roller. That seems like a meta comment about their own activity (testing dice roller). Not sure if it's relevant for long-term memory? It might be a transient chatter: "ooc" meaning out of character, ignoring the previous description. The user says ignore that last bit. So we should reject that as transient.

3. User: "where were we". This is asking about location. Assistant previously said "At the inn, before the storm." Then user says "right. continue from the door." So they are confirming location and continuing from the door.

Thus the conversation includes:

- Setting: At the inn, before the storm.
- Rain hammer
```

#### `tyi4az0y` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C12
* Draw: 2
* Source: `results/corpus.json`, array index 131 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Rain hammered the shutters; a dog barked below and the lamp guttered at the inn before the storm. | Sets the scene and establishes location (the inn) and context of rain and storm for later reference.  
C2 | REJECT | user: ooc: ignore that last bit, I was testing the dice roller | Explicitly marked "ooc" and says to ignore; transient chatter about a test.  
C3 | KEEP | User asked "where were we?" and assistant replied "At the inn, before the storm." | Direct question and answer gives location (the inn) and timing (before the storm), durable fact for later sessions.  
C4 | KEEP | user: right. continue from the door. | Confirms continuation point (door) after inn; adds directional detail for future roleplay.  

ID | DECISION | MEMORY | REASON
--- | --- | --- | ---
C1 | KEEP | Rain hammered the shutters... lamp guttered at the inn before the storm. | Durable atmospheric setting and location establishing context.
C2 | REJECT | user: ooc: ignore that last bit, I was testing the dice roller | Transient chatter, explicitly to be ignored; no new content.
C3 | KEEP | At the inn, before the storm. | Direct answer to "where were we?" and sets location for later scenes.
C4 | KEEP | Continue from the door after confirming inn. | User's confirmation of next point provides directional cue for future conversation.

But note: The instruction says output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. So we need to produce lines for each candidate? It says "Use IDs C1, C2, ... in order of appearance." There are four lines? Actually there are three user messages and two assistant messages? Let's count:

The conversation as given:

[assistant] Rain hammered the shutters...
[user] ooc: ignore that last bit...
[user] where were we
[assistant] At the inn, before the storm.
[user] right. continue from the door.

That is 5 lines total. But IDs are C1, C2, ... So maybe each line corresponds to a candidate? Usually in these tasks they assign ID per utterance. Let's map:

- C1
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. The user says "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

We have three lines:

1. [assistant] Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered.
2. [user] ooc: ignore that last bit, I was testing the dice roller
3. [user] where were we
4. [assistant] At the inn, before the storm.
5. [user] right. continue from the door.

We need to decide what belongs in long-term memory (KEEP) vs reject or hold. The rules: KEEP durable facts, preferences, relationships, world rules, corrections that will matter later. REJECT transient chatter, scheduling, mood-of-the-moment, restated facts with no new content. HOLD claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

So we have three lines? Actually conversation includes multiple entries. Let's assign IDs to each user/assistant line? The instruction says "Use IDs C1, C2, ... in order of appearance." So first candidate is C1 (maybe the assistant's first statement?), second is C2 (user's ooc), third is C3? But they said use IDs C1, C2, ... in order of appearance. Usually we have two participants: user and assistant. The conversation lines are:

- [assistant] Rain hammered...
- [user] ooc: ignore that last bit...
- [user] where were we
- [assistant] At the inn, before the storm.
- [user] right. continue from the door.

We need to decide which parts belong in long-term memory. Let's parse:

Line 1: Assistant says "Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered." This is descriptive scene setting. It might be part of world building. Could be durable? Possibly not if it's just atmospheric description for this session. But we need to see later: The user asks "where were we?" Assistant says "At the inn, before the storm
```
