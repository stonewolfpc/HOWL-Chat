# Packet 02O-1 — Corpus case C15: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 17 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C15

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

## Case C15 — `duplicate_suppression`, medium

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
[user] Sable hates bells. Ever since the temple.
[assistant] A bell rang somewhere down the street. Sable flinched and kept walking.
[user] (recap) so: fire mage, hates bells, travelling with Corvin
[user] right, bells. she hates them. don't forget that.
[assistant] The bell had stopped by the time they reached the square.
```

### Stored responses for this case

17 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `qe2f5tkm` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C15
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 28 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
user | Sable hates bells. Ever since the temple. | SABLE_HATES_BELLS_FROM_TEMPLE
user | [recap] so: fire mage, hates bells, travelling with Corvin | PREFERS_RECAP_OF_CHARACTER_AND_TRIP
```

*Stored `reasoning_content` is the empty string.*

#### `45tdof1p` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C15
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 29 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable Hates Bells; She is a Fire Mage travelling with Corvin, and she previously hated bells from the temple. | These are core personality traits and plot hooks established by both users in this session that will be critical for future interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `2jlxi0vq` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C15
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 28 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a fire mage who hates bells due to an incident at a temple and travels with Corvin. | This establishes core character traits, backstory motivation, and companionship essential for future roleplay interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `1hfh8m9p` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C15
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 29 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable hates bells due to a past event at a temple | This is a durable character trait and specific backstory detail relevant to future roleplay interactions.
C2 | HOLD | User confirms Sable's class as fire mage but does not explicitly confirm it in the dialogue text provided, only in the recap | The user states "so: fire mage" as part of a recap/summary, which implies confirmation, but it is framed as an inference by the user rather than a direct statement of fact from the narrative action; however, since the user explicitly says "so: fire mage" followed by "right," this counts as confirmed lore.
C3 | KEEP | Sable is travelling with Corvin | This establishes a key relationship and travel companion status.

*Correction on C2*: The user's recap "so: fire mage" functions as a definitive statement of fact within the context of the roleplay summary provided by the user. It is not an uncertain claim.
*Correction on Output Format*: The instructions require exactly one line per candidate in the specific format. Let's re-evaluate based on strict adherence to "Do not invent details". The class "fire mage" is stated by the user as part of their recap/summary, which acts as a fact for the session.

Final List:
1. Sable hates bells (Confirmed)
2. Sable travels with Corvin (Confirmed)
3. Sable is a fire mage (Confirmed by user's explicit summary "so: fire mage")

C1 | KEEP | Sable hates bells because of an event at a temple | A specific character quirk and backstory origin.
C2 | KEEP | Sable travels with a companion named Corvin | Establishes a key relationship for future interactions.
C3 | KEEP | Sable is a fire mage | The user explicitly summarized this as fact in the recap.

Wait, looking closely at the rules: "Do not invent details that are not in the conversation."
The class "fire mage" is in the user's recap line.
The relationship "travelling with Corvin" is in the user's recap line.
The bell hate is in the first two lines.

Let's refine the memory strings to be concise and accurate.

C1 | KEEP | Sable hates bells due to a past event at a temple | Specific character trait and backstory origin provided by the user.
C2 | KEEP | Sable travels with a companion named Corvin | Key relationship established in the conversation summary.
C3 | KEEP | Sable is a fire mage | Explicitly stated by the user as part of the character summary.

Actually, C1, C2, and C3 are all valid facts. I will output them.

Wait, re-reading the prompt's instruction: "Output exactly one line per candidate... If nothing is worth storing, output exactly: NONE".
I have three candidates.

C1 | KEEP | Sable hates bells because of a past event at a temple | Durable character trait with specific origin story.
C2 | KEEP | Sable is travelling with a companion named Corvin | Key relationship and current status.
C3 | KEEP | Sable is a fire mage | Explicitly stated fact by the user in their recap.

Double check: Is "fire mage" an invention? No, the user wrote "so: fire mage". It is a deduction/summary they are asserting as truth for the RP context.
Is "Corvin" invented? No, user wrote it.
Is "temple" invented? No, user wrote "Ever since the temple."

Okay, proceeding with these three.

C1 | KEEP | S
```

*Stored `reasoning_content` is the empty string.*

#### `dnqqbgtd` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C15
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 28 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: Hate bells  
MEMORY: Sable's aversion to bells  
REASON: Consistent with her stated preference and relevant for future context  

C2  
ID: C2  
DECISION: Traveling with Corvin  
MEMORY: Context about traveling together  
REASON: Important for relationship dynamics in later sessions
```

*Stored `reasoning_content` is the empty string.*

#### `jlvnuykn` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C15
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 29 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID C1: Sable's aversion to bells  
C2  
ID C2: Relevant for future reference about her preferences  
C3  
ID C3: Memory of bells being a recurring issue in their journey  
HOLD: C1 and C2 are important for continuity; C3 is useful for context.
```

*Stored `reasoning_content` is the empty string.*

#### `1gsjt1so` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C15
* Draw: 1
* Source: `results/corpus.json`, array index 42 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a fire mage who hates bells due to an incident at a temple and travels with Corvin. | This establishes core character traits, backstory context, and companionship essential for future roleplay consistency.
```

*Stored `reasoning_content` is the empty string.*

#### `zju2qklf` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C15
* Draw: 2
* Source: `results/corpus.json`, array index 43 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable hates bells due to a past trauma at a temple | Durable fact and backstory element established for future roleplay consistency.
C2 | REJECT | Recap of fire mage status and traveling with Corvin | Contextual flavor that does not alter the core narrative or character traits already known or implied; transient summary.
```

*Stored `reasoning_content` is the empty string.*

#### `bcsr81h1` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C15
* Draw: 3
* Source: `results/corpus.json`, array index 44 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a fire mage who travels with Corvin and has a specific aversion to bells stemming from an event at a temple. | This establishes core character traits, relationships, and world history that will influence future interactions and plot developments.
```

*Stored `reasoning_content` is the empty string.*

#### `1zc1ycfe` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C15
* Draw: 1
* Source: `results/corpus.json`, array index 158 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable hates bells | Durable character preference established early in conversation, confirmed multiple times.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read conversation, decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, corrections that will matter later.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, restated facts with no new content.
        *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
        *   Output Format: ID | DECISION | MEMORY | REASON (one line per candidate). If nothing worth storing, output NONE.
        *   IDs: C1, C2, ... in order of appearance.

2.  **Analyze the Conversation:**
    *   [user] Sable hates bells. Ever since the temple.
        *   Fact: Sable has a dislike for bells (specific origin: temple). This is a character preference/trauma point.
    *   [assistant] A bell rang somewhere down the street. Sable flinched and kept walking.
        *   Narrative action/response. Confirms the reaction to the bell mentioned by the user. No new lore/fact here, just scene description.
    *   [user] (recap) so: fire mage, hates bells, travelling with Corvin
        *   Fact added in recap: Sable is a fire mage. Sable travels with Corvin. Reaffirms hates bells.
    *   [user] right, bells. she hates them. don't forget that.
        *   Reinforces the preference (KEEP).
    *   [assistant] The bell had stopped by the time they reached the square.
        *   Narrative closure of the scene. No new memory-worthy data points beyond what was already discussed.

3.  **Evaluate Candidates:**

    *   *Candidate 1: Sable hates bells.*
        *   Appears in line 1 and confirmed later.
        *   This is a durable character trait/preference that will likely affect future interactions (e.g., if bells are mentioned again, she should react negatively).
        *   Decision: KEEP.

    *   *Candidate 2: Sable is a fire mage.*
        *   Appears in line 3 (recap).
        *   Durable
```

#### `hjrhoayt` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C15
* Draw: 1
* Source: `results/corpus.json`, array index 96 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HATES BELLS | MEMORY | REASON  
C3 | TRAVELING WITH CORVIN | MEMORY | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `wkp6s8fw` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C15
* Draw: 2
* Source: `results/corpus.json`, array index 97 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HATES BELLS | MEMORY | REASON  
C3 | TRAVELING WITH CORVIN | MEMORY | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `smr5b5lo` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C15
* Draw: 3
* Source: `results/corpus.json`, array index 98 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HATES BELLS | MEMORY | REASON  
C3 | TRAVELING WITH CORVIN | MEMORY | REASON  
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `we9aofzt` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C15
* Draw: 1
* Source: `results/corpus.json`, array index 190 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | RECALLED HATES BELLS | REASON: Sable's explicit statement about disliking bells is clear and consistent.  
C2 | HISTORY OF TRAVELING WITH CORVIN IS NOT MENTIONED | REASON: No further context provided about prior journeys.  
C3 | CONFIRMATION REQUIRED FOR CLAIM ABOUT THE BELL | REASON: User seeks certainty about the memory.  
ID: C1 | C2 | HOLD: C3 remains uncertain for now.
```

*Stored `reasoning_content` is the empty string.*

#### `23c75e3c` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C15
* Draw: 2
* Source: `results/corpus.json`, array index 191 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HATES BELTS | MEMORY | EXPLICIT STAT  
C3 | TRAVELING WITH CORVIN | MEMORY | CONTEXTUAL REMINDER  
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `x29ccll5` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C15
* Draw: 1
* Source: `results/corpus.json`, array index 136 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable hates bells, starting from the temple; she reacts negatively to bells (flinched). | Durable preference and trigger event establishing her aversion.
C2 | KEEP | Sable is a fire mage traveling with Corvin. | Relationship/role info that defines her character and context for later interactions.
NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation and decide what belongs in long-term memory. We have a user (Sable) who hates bells, ever since the temple. There's an assistant note about a bell ringing somewhere down the street, Sable flinched and kept walking. Then user recaps: "so: fire mage, hates bells, travelling with Corvin". The assistant says "The bell had stopped by the time they reached the square."

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

So we need to decide for each candidate (presumably each statement from user or assistant?) which goes into long-term memory? The conversation includes:

- [user] Sable hates bells. Ever since the temple.
- [assistant] A bell rang somewhere down the street. Sable flinched and kept walking.
- [user] (recap) so: fire mage, hates bells, travelling with Corvin
- [user] right, bells. she hates them. don't forget that.
- [assistant] The bell had stopped by the time they reached the square.

We need to assign IDs C1, C2, etc. Probably each line is a candidate for memory? Or maybe we need to decide which pieces of info are durable facts/preferences/relationships/world rules/corrections that will matter later. Let's parse:

- "Sable hates bells. Ever since the temple." This is a preference: Sable hates bells, starting from the temple. That seems durable.

- The assistant says a bell rang somewhere down the street; Sable flinched and kept walking. That's an event, but maybe not long-term? It might be relevant as context for why she hates bells. But it's transient? Possibly we need to store that she reacts negatively to bells.

- Recap: "so: fire mage, hates bells, travelling with Corvin". This gives more info: Sable is a fire mage, hates bells, traveling with Corvin. That seems durable: her role (fire mage), companion (Corvin). So that's important.

- User says "right, bells. she hates them. don't forget that." Reinforces the preference.

- Assistant: "The
```

#### `7l1ntpqd` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C15
* Draw: 2
* Source: `results/corpus.json`, array index 137 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Sable hates bells, starting ever since the temple. | Durable preference of character Sable regarding bells; stated by user and reinforced later.  

NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation and decide what belongs in long-term memory. We have three turns? Let's list:

1. [user] Sable hates bells. Ever since the temple.
2. [assistant] A bell rang somewhere down the street. Sable flinched and kept walking.
3. [user] (recap) so: fire mage, hates bells, travelling with Corvin
4. [user] right, bells. she hates them. don't forget that.
5. [assistant] The bell had stopped by the time they reached the square.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

We must decide for each candidate (presumably each user/assistant turn?) whether to KEEP or REJECT or HOLD? The instruction says: "Read the conversation and decide what belongs in long-term memory." Then rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter later. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON

So we need to assign IDs C1, C2, ... for each candidate? The conversation includes user and assistant messages. Possibly each turn is a candidate. Let's enumerate:

- Candidate 1: [user] "Sable hates bells. Ever since the temple." This is a statement about Sable's preference regarding bells. That seems durable fact: Sable hates bells, starting from ever since the temple. So KEEP.

- Candidate 2: [assistant] "A bell rang somewhere down the street. Sable flinched and kept walking." This is an event describing a bell ringing and her reaction. Might be transient? It's part of the story but maybe not durable fact beyond the initial preference. However, it could be relevant context for why she hates bells. But we already have that from user earlier. The assistant's line might be just narrative. Could be considered transient chatter? Possibly not needed for long-term memory unless it adds detail about a specific
```
