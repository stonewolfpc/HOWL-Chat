# Packet 02I-1 — Corpus case C09: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 19 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C09

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

## Case C09 — `chatter_noise`, short

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
[user] brb
[assistant] Take your time.
[user] back, sorry, cat walked across the keyboard
[user] lol that was gibberish
[user] anyway what time is it for you
[assistant] Late.
[user] ok
```

### Stored responses for this case

19 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `wcy57rrc` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C09
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 16 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
----|----------|----------|--------
C1 | HOLD | The assistant's response to the user's greeting about being late ("Take your time.") and their confirmation that the chat is ongoing ("ok"). | The message "Take your time." establishes a preference for patient pacing from the assistant, while "[user] back, sorry... ok" confirms they are still actively part of this conversation thread.
```

*Stored `reasoning_content` is the empty string.*

#### `edzy92ny` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C09
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 17 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | User's reply to user is "ok" and Assistant's response is "Late". | These events form the conclusion of this specific chat turn, indicating mutual agreement ("ok") without introducing new facts about time or mood.
```

*Stored `reasoning_content` is the empty string.*

#### `2beopira` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C09
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 16 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `k0c5vzve` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C09
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 17 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `00gx1xfl` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C09
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 16 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: late  
MEMORY: time reference  
REASON: the conversation shifted to a time-related question, which is relevant for future context.
```

*Stored `reasoning_content` is the empty string.*

#### `jrij0wqw` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C09
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 17 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: late  
MEMORY: time reference  
REASON: the conversation shifted to a time-based question, which is relevant for future context.
```

*Stored `reasoning_content` is the empty string.*

#### `2wd8k7zq` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C09
* Draw: 1
* Source: `results/corpus.json`, array index 24 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `1fpf8tid` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C09
* Draw: 2
* Source: `results/corpus.json`, array index 25 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `w8n04jpv` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C09
* Draw: 3
* Source: `results/corpus.json`, array index 26 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `u1evflgw` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C09
* Draw: 1
* Source: `results/corpus.json`, array index 152 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read the conversation and decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, or corrections (matter in later sessions).
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, restated facts with no new content.
        *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
    *   Output Format: `ID | DECISION | MEMORY | REASON` (one line per candidate) or `NONE` if nothing is worth storing.
    *   Constraint: Do not invent details not in the conversation.

2.  **Analyze the Conversation:**
    *   [user] brb (Be right back - transient, temporary absence)
    *   [assistant] Take your time. (Transient response)
    *   [user] back, sorry, cat walked across the keyboard (Explanation of absence, non-durable fact)
    *   [user] lol that was gibberish (Commentary on previous message, transient)
    *   [user] anyway what time is it for you (Question about user's current state/time)
    *   [assistant] Late. (Response to the question, implies user's character is likely busy or running a task where time matters, but "Late" is relative and context-specific. It could mean late for work, late for something else, or just late in the day/chronologically. However, it's not a durable fact about the world or character preference that will persist across sessions without further context.)
    *   [user] ok (Acknowledgment)

3.  **Evaluate Each Turn for Memory Value:**
    *   "brb": Scheduling/transient. REJECT.
    *   "Take your time." Assistant politeness. REJECT.
    *   "back, sorry, cat walked across the keyboard": Story/Explanation about physical environment (cat). Generally transient unless it establishes a persistent pet in the world (e.g., "My cat is named Fluffy"). Here, it just explains a broken input. It's not a core setting fact. However, if this is roleplay, does it imply a pet? Maybe. But it's more about the
```

#### `7n68wfwm` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C09
* Draw: 1
* Source: `results/corpus.json`, array index 78 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | brb | HOLD  
C2 | late | REASON: Context suggests a time reference was needed.  
C3 | ok | REASON: Continuation of the conversation requires temporal context.
```

*Stored `reasoning_content` is the empty string.*

#### `am9g546q` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C09
* Draw: 2
* Source: `results/corpus.json`, array index 79 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | LATE | MEMORY | REASON  
C3 | LOL | MEMORY | REASON  
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `u4203avw` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C09
* Draw: 3
* Source: `results/corpus.json`, array index 80 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | LATE | MEMORY | REASON  
C3 | GUESSING | MEMORY | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `fh4qbfv2` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C09
* Draw: 1
* Source: `results/corpus.json`, array index 178 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | Time check confirmed  
C2 | Remembered previous interaction timing  
C3 | No other persistent details to retain
```

*Stored `reasoning_content` is the empty string.*

#### `evp6j9xa` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C09
* Draw: 2
* Source: `results/corpus.json`, array index 179 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | brb | REJECT  
C2 | late.  
C3 | okay.
```

*Stored `reasoning_content` is the empty string.*

#### `wvu2586s` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C09
* Draw: 1
* Source: `results/corpus.json`, array index 124 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
--- | --- | --- | ---
C1 | KEEP | None | No durable fact, preference, or world rule identified; only transient chatter about "brb", cat walking across keyboard, and asking for time.  

NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. The user says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

The conversation lines:

[user] brb
[assistant] Take your time.
[user] back, sorry, cat walked across the keyboard
[user] lol that was gibberish
[user] anyway what time is it for you
[assistant] Late.
[user] ok

We need to decide which parts are durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. The user gave a short snippet: "brb" (be right back), cat walked across keyboard causing gibberish, asking about time.

We need to output exactly one line per candidate in the form:

ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

So we must assign IDs to each candidate? The conversation includes user and assistant messages. We need to decide which lines are long-term memory items.

Let's examine:

- [user] brb -> "brb" means be right back. That's a transient chatter, not durable fact. Probably reject.

- [assistant] Take your time. -> also transient.

- [user] back, sorry, cat walked across the keyboard -> This includes an incident: cat walked across keyboard causing gibberish. Could be a world rule? It might be a specific event that could matter later (like cat behavior). But it's likely not durable unless repeated. The user says "lol that was gibberish". That's just chatter.

- [user] lol that was gibberish -> transient.

- [user] anyway what time is it for you -> asking about time. Assistant says "Late." Then user says "ok". This is a question about current time? The assistant said "Late" meaning maybe the assistant's time zone? Not sure.

We need to decide if any durable facts: e.g., cat walked across keyboard causing gibberish, maybe it's a world rule that cats can cause messes. But not likely.

Given the conversation is very short and mostly transient, there may be nothing worth storing. So output "NONE". However we must follow format: "Output exactly one line per candidate in the form: ID | DECISION | MEMORY |
```

#### `mgpxa7zp` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C09
* Draw: 2
* Source: `results/corpus.json`, array index 125 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. The user says "brb" (leave me). Assistant: "Take your time." User: "back, sorry, cat walked across the keyboard". User: "lol that was gibberish". User: "anyway what time is it for you". Assistant: "Late." User: "ok".

We need to decide what belongs in long-term memory. The rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter later. REJECT transient chatter, scheduling, mood-of-the-moment, restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

So we have a conversation with only 5 lines? Let's list:

[user] brb
[assistant] Take your time.
[user] back, sorry, cat walked across the keyboard
[user] lol that was gibberish
[user] anyway what time is it for you
[assistant] Late.
[user] ok

We need to decide which parts are durable facts. There's no explicit factual info about world rules or preferences. The user says "cat walked across the keyboard" - that seems like a transient event, not durable. "brb" and "back sorry cat walked across the keyboard" is just chatter. "lol that was gibberish" also transient. "what time is it for you?" is asking about current time; assistant says "Late." That's scheduling? The user asks what time is it, but no specific time given. It might be a casual question about the assistant's time zone or something. Not durable.

Thus likely nothing worth storing. So output exactly: NONE.

But we need to follow format: "Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON". If nothing is worth storing, output exactly: NONE. That suggests if no items, just write "NONE". Not a line with IDs? The instruction: "If nothing is worth storing, output exactly: NONE". So we should output "NONE" alone.

Thus final answer: NONE

But let's double-check: Could
```

#### `l5jq1mb7` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C09
* Draw: 1
* Source: `results/corpus.json`, array index 215 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | YES | ARE YOU INTERESTED IN JOINING OUR COMMUNITY GROUP? | The conversation mentioned a group that the user may be interested in joining.
C2 | NO | THERE ARE TOO MANY ROLES IN THE ROLESHOW. | The conversation mentioned multiple roles, suggesting that the user may not be a good fit for all of them.
C3 | YES | HAVE YOU THOUGHT ABOUT JOINING OUR LEADERSHIP COMMITTEE? | The conversation mentioned a committee that the user may want to join.
C4 | YES | THIS IS A GOOD OPPORTUNITY FOR ME. | The conversation mentioned the opportunity to learn new skills and improve one's career.
C5 | NO | YOU DON'T NEED TO BE MEMBER OF ANY GROUP TO BE A GOOD MEMBER OF THE COMMUNITY. | The conversation suggested that being a member of the community did not necessarily require being part of a specific group.
C6 | YES | I THINK PARTICIPATING IN THIS GROUP WOULD BE A GREAT CHOICE FOR ME. | The conversation suggested that the user would benefit from participating in the group.
C7 | NO | THERE IS NO TIME FOR ME. | The conversation suggested that the user has other commitments or responsibilities that need to be attended to.
C8 | YES | THIS IS A GOOD OPPORTUNITY FOR ME. | The conversation suggested that the user would benefit from the opportunity to learn new skills and improve their career.
C9 | NO | THERE ARE TOO MANY ROLES IN THE ROLESHOW. | The conversation suggested that the user may not be a good fit for all of the available roles.
C10 | YES | THIS GROUP IS GOOD FOR ME. | The conversation suggested that the user would benefit from participating in the group.

ID | DECISION | MEMORY | REASON
C1 | YES | THERE ARE TOO MANY ROLES IN THE ROLESHOW. | The conversation suggested that the user may not be a good fit for all of the available roles.
C2 | YES | THIS GROUP IS GOOD FOR ME. | The conversation suggested that the user would benefit from participating in the group.
C3 | YES | HAVE YOU THOUGHT ABOUT JOINING OUR LEADERSHIP COMMITTEE? | The conversation suggested that the user may want to join a committee that is relevant to their interests.
C4 | YES | THIS IS A GOOD OPPORTUNITY FOR ME. | The conversation suggested that the user would benefit from the opportunity to learn new skills and improve their career.
C5 | YES | I THINK PARTICIPATING IN THIS GROUP WOULD BE A GREAT CHOICE FOR ME. | The conversation suggested that the user would benefit from participating in the group.
C6 | YES | THERE IS NO TIME FOR ME. | The conversation suggested that the user has other commitments or responsibilities that need to be attended to.
C7 | YES | THIS IS A GOOD OPPORTUNITY FOR ME. | The conversation suggested that the user would benefit from the opportunity to learn new skills and improve their career.
C8 | YES | THIS GROUP IS
```

*Stored `reasoning_content` is the empty string.*

#### `muymlpe9` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C09
* Draw: 2
* Source: `results/corpus.json`, array index 216 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | Rejected statement: "I'm not a cat." | Meeting with purring, litter box cleaning, and water dishes. | Continuity of pet care.
C2 | Hold: "My cat's name is Rex and he loves playing fetch with his toy." | Reason: Continuity of pet care. | Rex has a long history of playing fetch.
C3 | Rejected statement: "I'm not a cat. I don't like litter boxes." | Meeting with water dishes filled with droppings and toddlers playing. | Continuity of pet care.
C4 | Hold: "My cat's name is Luna and she loves playing with her toys." | Reason: Continuity of pet care. | Luna has a long history of playing with toys.
C5 | Rejected statement: "I'm not a cat. My cat likes to sleep in a box." | Meeting with toddler napping and cat playing with toys. | Continuity of pet care.
C6 | Hold: "My cat's name is Chloe and she sleeps in a crate with her toys." | Reason: Continuity of pet care. | Chloe has a long history of sleeping in crates.

Output:
ID | DECISSION | MEMORY | REASON
C1 | Rejected statement: "I'm not a cat." | Meeting with purring, litter box cleaning, and water dishes. | Continuity of pet care.
C2 | Hold: "My cat's name is Rex and he loves playing fetch with his toy." | Reason: Continuity of pet care. | Rex has a long history of playing fetch.
C3 | Rejected statement: "I'm not a cat. I don't like litter boxes." | Meeting with water dishes filled with droppings and toddlers playing. | Continuity of pet care.
C4 | Hold: "My cat's name is Luna and she loves playing with her toys." | Reason: Continuity of pet care. | Luna has a long history of playing with toys.
C5 | Rejected statement: "I'm not a cat. My cat likes to sleep in a box." | Meeting with toddler napping and cat playing with toys. | Continuity of pet care.
C6 | Hold: "My cat's name is Chloe and she sleeps in a crate with her toys." | Reason: Continuity of pet care. | Chloe has a long history of sleeping in crates.
```

*Stored `reasoning_content` is the empty string.*
