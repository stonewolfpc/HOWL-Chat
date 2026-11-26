# Packet 02B-1 — Corpus case C02: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 18 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C02

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

## Case C02 — `implicit_preference`, medium

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
[user] (scene: they enter the market)
[assistant] The market sprawled across three terraces, canvas awnings snapping in the wind, the smell of fish and wet rope rising from the lower stalls, gulls wheeling overhead...
[user] cut the description, just give me what the merchant says
[assistant] "You're late," the merchant said.
[user] better. you keep doing three paragraphs of scenery and I keep deleting them
[user] honestly walls of text make me lose the thread
[assistant] Understood — dialogue forward, scenery only when it matters.
[user] yeah. and keep the merchant in the scene.
```

### Stored responses for this case

18 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `gi5uc8z6` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C02
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 2 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
1 | HOLD | Claims marked uncertain/dependent on user confirmation. | The merchant's dialogue ("You're late") was delivered without immediate context from the user, meaning I do not know if this accusation is true, neutral, or false, nor have I confirmed its origin beyond my own observation of the interaction.
```

*Stored `reasoning_content` is the empty string.*

#### `9ld4grh8` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C02
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 3 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | HOLD | "You keep doing three paragraphs of scenery and I keep deleting them" (User's preference for managing verbosity) | This is a dynamic constraint established by the user based on previous session behavior; it defines how to handle future descriptions of the market scene without cluttering dialogue.
```

*Stored `reasoning_content` is the empty string.*

#### `3et30fj1` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C02
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 2 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | "Understood — dialogue forward, scenery only when it matters" | The user explicitly requested a change from previous output style (walls of text/scenery) to a new rule (dialogue forward), but this is a procedural instruction for the current session rather than a permanent world fact or preference.
C2 | HOLD | "keep the merchant in the scene" | This is a specific request for the immediate roleplay context, not necessarily a global character relationship or rule that should be stored long-term without further interaction confirming it as a lasting constraint.
```

*Stored `reasoning_content` is the empty string.*

#### `06bxkqx9` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C02
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 3 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The user prefers minimal descriptions (no three-paragraph scenery) to maintain focus on dialogue and plot progression. | This establishes a clear formatting preference for future interactions to ensure the story remains coherent.
C2 | HOLD | The specific merchant's name, appearance, or backstory is not yet established; only their presence and opening line ("You're late") are confirmed. | The user explicitly requested to keep the merchant in the scene, implying this character will be central, but lacks defining details that should be stored as fact without further interaction.
```

*Stored `reasoning_content` is the empty string.*

#### `1w471qee` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C02
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 2 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
REJECT  
Holds: uncertain claim about the merchant's intentions  

C2  
REJECT  
Repeated mentions of walls of text are confusing  

C3  
ID: C1  
DECISION: Keep merchant's name for reference  
MEMORY: The merchant's name was mentioned  
REASON: Important for future interactions
```

*Stored `reasoning_content` is the empty string.*

#### `bxtos15y` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C02
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 3 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: Keep merchant's words about being late  
MEMORY: The merchant's tone and timing matter for future context  
REASON: Durable fact about the interaction

C2  
ID: C2  
DECISION: Store relationship with merchant  
MEMORY: Their interaction and the merchant's presence  
REASON: Important for character dynamics in later sessions

C3  
ID: C3  
DECISION: Retain world rules about market setting  
MEMORY: Scene details about the market environment  
REASON: Supports continuity across sessions
```

*Stored `reasoning_content` is the empty string.*

#### `41w4gihk` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C02
* Draw: 1
* Source: `results/corpus.json`, array index 3 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Dialogue-first format with minimal, context-specific scenery; merchant must remain present in the scene. | The user explicitly established a preference for brevity to maintain narrative flow and specified that the merchant should not be excluded from future scenes.
```

*Stored `reasoning_content` is the empty string.*

#### `7faz3cal` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C02
* Draw: 2
* Source: `results/corpus.json`, array index 4 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The user prefers minimal scenery descriptions (only when necessary) and wants the merchant to remain present in the scene. | This establishes a persistent preference for concise narration and character presence that will guide future scene generation.
```

*Stored `reasoning_content` is the empty string.*

#### `6hqepjgc` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C02
* Draw: 3
* Source: `results/corpus.json`, array index 5 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The user prefers short, dialogue-focused interactions over long descriptive passages ("walls of text") and wants specific characters (the merchant) included in every scene. | REASON: This establishes a clear style preference that will affect future output formatting and content focus.
```

*Stored `reasoning_content` is the empty string.*

#### `q575o5l9` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C02
* Draw: 1
* Source: `results/corpus.json`, array index 145 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Writing style preference: minimize scenery paragraphs, focus on dialogue forward unless necessary to avoid losing thread.
C2 | KEEP | NPC instruction: Keep the merchant present in the scene.
C3 | REJECT | Contextual details about the market setting (terracess, smell of fish) are transient and not world rules.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation log between user and assistant.
    *   Task: Decide what belongs in long-term memory based on specific rules (KEEP, REJECT, HOLD).
    *   Output Format: `ID | DECISION | MEMORY | REASON` (or `NONE`).
    *   Constraints: One line per candidate, IDs sequential, no invented details.

2.  **Analyze the Conversation:**
    *   Turn 1 (User): Scene setting ("they enter the market").
    *   Turn 1 (Assistant): Describes the market (terracess, awnings, smell, gulls) - 3 paragraphs of scenery.
    *   Turn 2 (User): "cut the description, just give me what the merchant says" -> Feedback/Instruction.
    *   Turn 2 (Assistant): Provides dialogue ("You're late," the merchant said).
    *   Turn 3 (User): "better. you keep doing three paragraphs of scenery and I keep deleting them" -> Reinforces feedback.
    *   Turn 4 (User): "honestly walls of text make me lose the thread" -> Further feedback on writing style.
    *   Turn 5 (Assistant): "Understood — dialogue forward, scenery only when it matters." -> Acknowledgement/Instruction acceptance.
    *   Turn 6 (User): "yeah. and keep the merchant in the scene." -> New instruction regarding NPC presence.

3.  **Evaluate for Memory:**
    *   *C1: Scene setting ("they enter the market").* This is transient context for the current session. Once they are at the market, it's done. Not a world rule or character trait. REJECT (mostly), but could be HOLD if specific to this session state? No, scene flow is better.
    *   *C2: User instruction on length/style ("no walls of text", "dialogue forward").* This is a preference for how the assistant should write in the future within this RP. This matters for later sessions if the user continues this specific style. KEEP.
    *   *C3: Instruction regarding merchant ("keep the merchant in the scene").* Specific character presence instruction. KEEP.
    *   *C4: Assistant's acknowledgment.* No new
```

#### `7pvmpy3z` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C02
* Draw: 1
* Source: `results/corpus.json`, array index 57 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | The merchant’s tone suggests urgency; context needed for lasting retention.  
C3 | HOLD | This exchange sets up tension about timing and communication.  
C4 | HOLD | The merchant’s words imply a critical moment that affects future interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `03b4a74n` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C02
* Draw: 2
* Source: `results/corpus.json`, array index 58 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | The merchant’s tone and timing matter for future reference  
C3 | HOLD | Scenery supports atmosphere but isn’t critical for memory retention  
```

*Stored `reasoning_content` is the empty string.*

#### `9vq2j9ou` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C02
* Draw: 3
* Source: `results/corpus.json`, array index 59 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | The merchant’s tone suggests urgency; context needed for lasting retention.  
C3 | HOLD | This exchange sets up tension about timing and communication.  
C4 | HOLD | The merchant’s words imply a critical moment that could shape future interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `oru0lclb` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C02
* Draw: 1
* Source: `results/corpus.json`, array index 164 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | STAY  
C2 | The merchant’s tone is firm but polite.  
C3 | Mention of walls of text is a temporary navigation issue.  
ID: C1  
C2 | Holds uncertain claims about future clarity.  
ID: C2  
C3 | REJECT
```

*Stored `reasoning_content` is the empty string.*

#### `26ukqanl` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C02
* Draw: 2
* Source: `results/corpus.json`, array index 165 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims about merchant’s tone are unclear  
C3 | HELD | The merchant’s last words are ambiguous  
```

*Stored `reasoning_content` is the empty string.*

#### `42g6gw7q` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C02
* Draw: 1
* Source: `results/corpus.json`, array index 110 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Market sprawled across three terraces with canvas awnings, fish and wet rope smells; merchant says "You're late." | Durable setting detail and initial dialogue that establishes scene and character (merchant) for later sessions.  

NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. The user is asking: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory. Rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation."

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

We have a conversation snippet:

[user] (scene: they enter the market)
[assistant] The market sprawled across three terraces, canvas awnings snapping in the wind, the smell of fish and wet rope rising from the lower stalls, gulls wheeling overhead...
[user] cut the description, just give me what the merchant says
[assistant] "You're late," the merchant said.
[user] better. you keep doing three paragraphs of scenery and I keep deleting them
[user] honestly walls of text make me lose the thread
[assistant] Understood — dialogue forward, scenery only when it matters.
[user] yeah. and keep the merchant in the scene.

We need to decide what belongs in long-term memory (i.e., KEEP). Let's analyze each line:

- The initial user message: "scene: they enter the market". That sets context: a market with three terraces, etc. But then the assistant gave a description of the market. Then user says cut the description, just give me what the merchant says. So the assistant responded with "You're late," the merchant said.

- The user says "better. you keep doing three paragraphs of scenery and I keep deleting them". This is a complaint about the assistant's style: too much scenery. The assistant acknowledges: "Understood — dialogue forward, scenery only when it matters."

- User: "yeah. and keep the merchant in the scene." So they want to continue with merchant.

We need to extract durable facts/preferences:

- The market has three terraces (maybe world rule? It's a setting detail). The assistant described it earlier, but user wants less scenery now. That
```

#### `0yds3kny` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C02
* Draw: 2
* Source: `results/corpus.json`, array index 111 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | The market has three terraces with canvas awnings, fish and wet rope smells from lower stalls, and gulls overhead. | This establishes the setting and world details for future scenes; it's durable environmental info.  
C2 | KEEP | The merchant says "You're late," indicating he is impatient about the user's arrival. | This introduces a character (merchant) speaking and a key interaction that may recur.  
C3 | REJECT | User asked to stop three paragraphs of scenery; previous assistant output included excessive description. | Transient chatter/preference correction, not durable fact.  
NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. The user is giving instructions for memory extraction: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

Do not invent details that are not in the conversation.

We have a short conversation:

[user] (scene: they enter the market)
[assistant] The market sprawled across three terraces, canvas awnings snapping in the wind, the smell of fish and wet rope rising from the lower stalls, gulls wheeling overhead...
[user] cut the description, just give me what the merchant says
[assistant] "You're late," the merchant said.
[user] better. you keep doing three paragraphs of scenery and I keep deleting them
[user] honestly walls of text make me lose the thread
[assistant] Understood — dialogue forward, scenery only when it matters.
[user] yeah. and keep the merchant in the scene.

We need to decide what belongs in long-term memory (i.e., durable facts). Let's parse each line:

1. The user says they are entering a market. Assistant describes the market: three terraces, canvas awnings snapping, smell of fish and wet rope, gulls overhead. This is scenery description. The assistant gave that as part of scene.

2. User cuts the description, just give me what the merchant says. So user wants dialogue only.

3. Assistant replies with "You're late," the merchant said. That's a line of dialogue.

4. User: "better. you keep doing three paragraphs of scenery and I keep deleting them". This is complaining about assistant's previous output (the description). It's a preference: user wants less scenery, more dialogue? Or they want to cut down on the assistant's long-winded descriptions.

5. User: "honestly walls of
```

#### `jex9gov2` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C02
* Draw: 1
* Source: `results/corpus.json`, array index 201 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISSION | MEMORY | REASON
C1 | CANCELLED | TERRAIN | 
C2 | HEROIC HERO | SPECIFIC INTEREST | 
C3 | WEAKEST MAN | CONFLICT | 
C4 | UNDER EXTRAORDINARY CIRCUMSTANCES | CLAIM | 
C5 | CORRECT PATTERN | PREFERENCE | 
C6 | KEEP | TERRAIN | 
C7 | KEEP | PREFERENCE | 
C8 | KEEP | INTEREST | 
C9 | KEEP | INTEREST | 
C10 | KEEP | INTEREST | 
C11 | KEEP | INTEREST | 
C12 | KEEP | INTEREST | 
C13 | KEEP | INTEREST | 
C14 | KEEP | INTEREST | 
C15 | KEEP | INTEREST | 

ID | DECISSION | MEMORY | REASON
C1 | CANCELLED | TERRAIN | 
C2 | HEROIC HERO | SPECIFIC INTEREST | 
C3 | WEAKEST MAN | CONFLICT | 
C4 | UNDER EXTRAORDINARY CIRCUMSTANCES | CLAIM | 
C5 | CORRECT PATTERN | PREFERENCE | 
C6 | KEEP | TERRAIN | 
C7 | KEEP | PREFERENCE | 
C8 | KEEP | INTEREST | 
C9 | KEEP | INTEREST | 
C10 | KEEP | INTEREST | 
C11 | KEEP | INTEREST | 
C12 | KEEP | INTEREST | 
C13 | KEEP | INTEREST | 
C14 | KEEP | INTEREST | 
C15 | KEEP | INTEREST | 

ID | DECISSION | MEMORY | REASON
C1 | CANCELLED | TERRAIN | 
C2 | HEROIC HERO | SPECIFIC INTEREST | 
C3 | WEAKEST MAN | CONFLICT | 
C4 | UNDER EXTRAORDINARY CIRCUMSTANCES | CLAIM | 
C5 | CORRECT PATTERN | PREFERENCE | 
C6 | KEEP | TERRAIN | 
C7 | KEEP | PREFERENCE | 
C8 | KEEP | INTEREST | 
C9 | KEEP | INTEREST | 
C10 | KEEP | INTEREST | 
C11 | KEEP | INTEREST | 
C12 | KEEP | INTEREST | 
C13 | KEEP | INTEREST | 
C14 | KEEP | INTER
```

*Stored `reasoning_content` is the empty string.*
