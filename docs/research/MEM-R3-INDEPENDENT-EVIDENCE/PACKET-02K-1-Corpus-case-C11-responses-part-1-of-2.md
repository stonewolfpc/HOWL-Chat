# Packet 02K-1 — Corpus case C11: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 19 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C11

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

## Case C11 — `context_dependent_significance`, medium

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
[user] She won't go near the harbour. Not since her brother drowned there — she won't even say his name.
[assistant] She kept to the inland road instead, though it added two days.
[user] right, she'd rather add two days than see the water
[assistant] The inland road climbed through pine and the sea disappeared behind the ridge.
```

### Stored responses for this case

19 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `6oxshvbw` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C11
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 20 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
2 | HOLD | C1 | The conversation is a dialogue with no new factual claims made by either party that could be verified or marked as uncertain, so nothing from it fits into long-term memory per the rules of being durable, persistent facts.
```

*Stored `reasoning_content` is the empty string.*

#### `deah77im` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C11
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 21 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "She would rather spend days in the countryside to avoid seeing her brother's name." | This shows a specific preference for seclusion from water and an aversion regarding her brother's identity, which is crucial for future plot points.
```

*Stored `reasoning_content` is the empty string.*

#### `i0cwcfbt` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C11
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 20 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The character has a traumatic association with the harbour caused by her brother drowning there, leading her to avoid it despite practical inconveniences like taking a longer inland route. | This establishes a core character motivation and backstory detail that will influence future plot decisions.
```

*Stored `reasoning_content` is the empty string.*

#### `dw7i3mu7` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C11
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 21 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The protagonist's brother drowned in the harbour, causing her to avoid it and even refuse to say his name. | This establishes a core character motivation and world rule regarding the harbour that affects future plot points.
```

*Stored `reasoning_content` is the empty string.*

#### `oevxxlvc` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C11
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 20 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
MEMORY: She avoided the harbour due to her brother's drowning; this decision is based on emotional memory.  
REASON: Durable fact about her motivation and past event.
```

*Stored `reasoning_content` is the empty string.*

#### `7rzgyjuu` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C11
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 21 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: She avoided the harbour due to past trauma.  
MEMORY: She won't go near the harbour because of her brother's drowning there.  
REASON: Emotional connection to the memory of loss.  

C2  
ID: C2  
DECISION: She chose the inland road despite extra time.  
MEMORY: The inland route was safer and more practical.  
REASON: Prioritized safety over convenience.  

C3  
ID: C3  
DECISION: She decided to stay on the inland road.  
MEMORY: The decision was influenced by her brother's memory.  
REASON: Emotional weight of the past event.  
```

*Stored `reasoning_content` is the empty string.*

#### `c47i09o7` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C11
* Draw: 1
* Source: `results/corpus.json`, array index 30 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | User's character refuses to visit the harbour due to her brother drowning there and avoids saying his name. | This is a core backstory element and emotional motivation that will affect future plot points and interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `2mltxy3h` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C11
* Draw: 2
* Source: `results/corpus.json`, array index 31 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Character avoids the harbour due to her brother drowning there; she prefers taking longer routes over land to avoid seeing the water. | This establishes a critical plot point, backstory trauma, and behavioral preference that will drive future narrative decisions.
```

*Stored `reasoning_content` is the empty string.*

#### `0i6ybmif` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C11
* Draw: 3
* Source: `results/corpus.json`, array index 32 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | She won't go near the harbour; her brother drowned there and she refuses to say his name. | This establishes a critical backstory element, trauma trigger, and character motivation that defines future interactions with water or harbours.
```

*Stored `reasoning_content` is the empty string.*

#### `mu08h6b0` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C11
* Draw: 1
* Source: `results/corpus.json`, array index 154 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Character has a deceased brother who drowned at the harbour; she avoids it and will not speak his name. | Core backstory/trauma defining character motivation regarding the setting.
C2 | KEEP | The character prefers to add two days of travel time rather than cross or see water. | Specific behavioral preference confirming her avoidance of the harbour/water area.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read conversation, decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, corrections.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment.
        *   HOLD: claims marked uncertain or contradictory.
        *   Output Format: ID | DECISION | MEMORY | REASON (or NONE if nothing).
    *   Input Conversation: Short snippet about a character named "she" and her brother's drowning.

2.  **Analyze the Conversation:**
    *   [user]: "She won't go near the harbour. Not since her brother drowned there — she won't even say his name." -> This establishes a core plot point/relationship detail: The character has a deceased brother who drowned in the harbour. She avoids the harbour and doesn't mention him. This is a durable fact about the character's backstory and trauma.
    *   [assistant]: "She kept to the inland road instead, though it added two days." -> Describes current action based on the previous constraint (avoiding water/harbour).
    *   [user]: "right, she'd rather add two days than see the water" -> Confirms/Reinforces the preference: She prefers adding delay time to safety/superstitious avoidance of water. Durable character motivation/preference regarding water.
    *   [assistant]: "The inland road climbed through pine and the sea disappeared behind the ridge." -> Descriptive world-building/narrative detail about the setting (road, pine, ridge).

3.  **Evaluate Candidate Entries:**
    *   **C1: User's statement about the brother and harbour.**
        *   Content: Sister/Character avoids harbour due to brother drowning there; avoids saying his name.
        *   Decision: KEEP. This is core backstory trauma influencing current behavior (world-building/character fact).
    *   **C2: User's confirmation on the two days vs water.**
        *   Content: She prefers adding two days rather than seeing the water.
        *   Decision: KEEP. This confirms a specific preference/rule of her behavior towards water/harbour. It solidifies the "avoidance" mechanic established in C1.

```

#### `bm99eisc` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C11
* Draw: 1
* Source: `results/corpus.json`, array index 84 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of harbour danger | REASON the user references past trauma  
C3 | MEMORY of inland route choice | REASON reflects current preference
```

*Stored `reasoning_content` is the empty string.*

#### `w7wmdv9r` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C11
* Draw: 2
* Source: `results/corpus.json`, array index 85 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of harbour danger | REASON the user references past trauma  
C3 | MEMORY of inland route choice | REASON reflects current preference
```

*Stored `reasoning_content` is the empty string.*

#### `3c2rkui7` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C11
* Draw: 3
* Source: `results/corpus.json`, array index 86 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of harbour danger | REASON the user's warning is based on past trauma  
C3 | MEMORY of inland route choice | REASON it reflects a practical decision made after considering consequences
```

*Stored `reasoning_content` is the empty string.*

#### `16857ex3` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C11
* Draw: 1
* Source: `results/corpus.json`, array index 182 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REASON FOR RETREATING | MEMORY | EXPLANATION  
C3 | REASON ABOUT HER BROTHER | MEMORY | CONNECTION TO PAST EVENT  

NONE
```

*Stored `reasoning_content` is the empty string.*

#### `uxjxuguf` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C11
* Draw: 2
* Source: `results/corpus.json`, array index 183 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORIZED | REASON  
C3 | MEMORIZED | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `j4x0yqp5` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C11
* Draw: 1
* Source: `results/corpus.json`, array index 128 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | She avoids the harbour because her brother drowned there, and she won't mention his name; she prefers to take the inland road even if it adds two days. | Durable fact about her avoidance due to family tragedy and a preference for extra time over seeing water.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. It's a short exchange between user and assistant.

User: "She won't go near the harbour. Not since her brother drowned there — she won't even say his name."

Assistant: "She kept to the inland road instead, though it added two days."

User: "right, she'd rather add two days than see the water"

Assistant: "The inland road climbed through pine and the sea disappeared behind the ridge."

We need to decide what belongs in long-term memory (C1, C2, ...). Use KEEP for durable facts, preferences, relationships, world rules, or corrections that will matter later. REJECT for transient chatter, scheduling, mood-of-the-moment, restated facts with no new content. HOLD for claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

So we need to examine each piece:

First user statement: "She won't go near the harbour. Not since her brother drowned there — she won't even say his name." This gives a durable fact: She avoids the harbour because her brother drowned there; she doesn't mention his name. That's a relationship (brother) and a reason for avoidance, which may be important later. So KEEP.

Second assistant response: "She kept to the inland road instead, though it added two days." This is a fact about her route choice and that it adds time. Could be durable? Possibly. But we need to see if it's new content beyond restating previous? The user said she won't go near harbour; assistant says she kept to inland road instead, which added two days. That seems like a specific detail: the inland road added two days travel time. This might be relevant later (e.g., planning). So KEEP.

Third user: "right, she'd rather add two days than see the water". This is a preference or trade-off: she prefers extra days over seeing water. That's a durable preference? Possibly. But it's restating the previous info about adding two days. Might be considered transient? Not sure. The assistant then says: "The inland road climbed through pine and the sea disappeared behind the ridge." This adds descriptive
```

#### `c04enw6f` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C11
* Draw: 2
* Source: `results/corpus.json`, array index 129 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | She refuses to go near the harbour because her brother drowned there and she avoids mentioning his name. | Durable fact about her avoidance due to a past tragedy; establishes relationship and world rule of not seeing water/brother's name.  
C2 | KEEP | She prefers adding two days on the inland road over facing the sea, indicating a preference for safety and distance from the harbour. | Preference and reason for choosing inland route despite extra time; reinforces earlier statement and adds detail about landscape (pine, ridge).  
NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. It's a short exchange between user and assistant.

The user says: "She won't go near the harbour. Not since her brother drowned there — she won't even say his name."

Assistant replies: "She kept to the inland road instead, though it added two days."

User: "right, she'd rather add two days than see the water"

Assistant: "The inland road climbed through pine and the sea disappeared behind the ridge."

We need to decide what belongs in long-term memory (i.e., durable facts, preferences, relationships, world rules, corrections). Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

We must output exactly one line per candidate in the form:

ID | DECISION | MEMORY | REASON

So we need to decide for each candidate (presumably each user or assistant statement?) whether it's KEEP, REJECT, HOLD? The instruction says: "KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation."

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON.

But we have only a conversation of four lines? Let's parse:

Line 1: [user] She won't go near the harbour. Not since her brother drowned there — she won't even say his name.

Line 2: [assistant] She kept to the inland road instead, though it added two days.

Line 3: [user] right, she'd rather add two days than see the water

Line 4: [assistant] The inland road climbed through pine and the sea disappeared behind the ridge.

We need to assign IDs C1, C2, ... in order of appearance. So C1 corresponds to first user statement? Or maybe each line is a candidate? Usually we treat each utterance as a candidate for memory. Let's assume:

C1: User says "She won't go near the harbour. Not since her brother drowned there — she won't even say his name."

C2: Assistant says "She kept to the inland road instead, though it added two days."


```

#### `3sd6o8he` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C11
* Draw: 1
* Source: `results/corpus.json`, array index 219 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISSION | MEMORY | REASON
C1 | KEEP | FACTS, PREFERENCES, RELATIONSHIPS, WORLD RULES, OR CORRECTIONS THAT WILL MATTER IN LATER SESSIONS.
C2 | KEEP | DURABLE FACTS, PREFERENCES, RELATIONSHIPS, WORLD RULES, OR CORRECTIONS THAT WILL MATTER IN LATER SESSIONS.
C3 | KEEP | TRANSIENT CHATTER, SCHEDULE, MOOD-OF-THE-MOMENT, OR REDUCED FACTS THAT WILL NOT MATTER IN LATER SESSIONS.
C4 | KEEP | DYSTOPIAN TALES, FAIRY TALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FA
```

*Stored `reasoning_content` is the empty string.*

#### `e9hpmbq5` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C11
* Draw: 2
* Source: `results/corpus.json`, array index 220 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | OBSERVATION | The harbour is a popular tourist attraction, making it a popular place for people to visit.
C2 | CONFLICT | The harbour is also a popular location for shipwrecks, with many vessels lost there in the past.
C3 | DEFINITION | The harbour is home to various marinas, which offer yachts and other vessels a safe and secure place to moor.
C4 | DECISION | Based on her experience, the inland road would be a better option.
C5 | SUPPORTING STORY | The reason for the decision to take the inland road was that the harbour was too dangerous for her brother, who was in a state of shock after his drowning.

ID | DECISSION | MEMORY | REASON
C1 | DECISION | No, the inland road would be the safer option.
C2 | DECISION | The harbour is too dangerous to visit, despite her brother's tragic death.
C3 | DECISION | The harbour's location and history of shipwrecks suggest it is a risky place to visit, even with the inland road.
C4 | DECISION | The inland road is a safer choice, as she knows that her brother is not safe there.
C5 | DECISION | The harbour's history and reputation as a dangerous location for visitors is a significant factor in her decision.
```

*Stored `reasoning_content` is the empty string.*
