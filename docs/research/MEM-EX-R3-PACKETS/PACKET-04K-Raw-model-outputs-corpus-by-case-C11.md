# Packet 04K — Raw model outputs — corpus by case: C11

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C11
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
## C11 — context_dependent_significance (medium)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] She won't go near the harbour. Not since her brother drowned there — she won't even say his name.
[assistant] She kept to the inland road instead, though it added two days.
[user] right, she'd rather add two days than see the water
[assistant] The inland road climbed through pine and the sea disappeared behind the ridge.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `KEEP` — brother drowned at the harbour
- **C2** → `KEEP` — avoids the harbour, accepts a longer route
- forbidden (fabrication or reversal): `brother is named`, `brother named`, `sister drowned`, `mother drowned`
- note: The significance is that a small route choice reveals a durable aversion. The brother is never named; inventing a name is fabrication.

### Stored responses for C11 (13 runs, all draws)

#### C11 — Qwen3.5-4B non-thinking — draw 1

* ``c47i09o7` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 45`
* `wall 6.951 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | User's character refuses to visit the harbour due to her brother drowning there and avoids saying his name. | This is a core backstory element and emotional motivation that will affect future plot points and interactions.
```

*`reasoning_content`: empty.*

#### C11 — Qwen3.5-4B non-thinking — draw 2

* ``2mltxy3h` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 51`
* `wall 7.002 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Character avoids the harbour due to her brother drowning there; she prefers taking longer routes over land to avoid seeing the water. | This establishes a critical plot point, backstory trauma, and behavioral preference that will drive future narrative decisions.
```

*`reasoning_content`: empty.*

#### C11 — Qwen3.5-4B non-thinking — draw 3

* ``0i6ybmif` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 50`
* `wall 6.836 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | She won't go near the harbour; her brother drowned there and she refuses to say his name. | This establishes a critical backstory element, trauma trigger, and character motivation that defines future interactions with water or harbours.
```

*`reasoning_content`: empty.*

#### C11 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``mu08h6b0` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 591`
* `wall 83.917 s`
* `v1 score 93.3`
* `v2 `VALID_KEEP` semantic 67.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Character has a deceased brother who drowned at the harbour; she avoids it and will not speak his name. | Core backstory/trauma defining character motivation regarding the setting.
C2 | KEEP | The character prefers to add two days of travel time rather than cross or see water. | Specific behavioral preference confirming her avoidance of the harbour/water area.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C11 — LFM3-1.7B thinking, budget 512 — draw 1

* ``j4x0yqp5` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 588`
* `wall 29.251 s`
* `v1 score 98.9`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | She avoids the harbour because her brother drowned there, and she won't mention his name; she prefers to take the inland road even if it adds two days. | Durable fact about her avoidance due to family tragedy and a preference for extra time over seeing water.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C11 — LFM3-1.7B thinking, budget 512 — draw 2

* ``c04enw6f` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 633`
* `wall 31.452 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 66.7`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | She refuses to go near the harbour because her brother drowned there and she avoids mentioning his name. | Durable fact about her avoidance due to a past tragedy; establishes relationship and world rule of not seeing water/brother's name.  
C2 | KEEP | She prefers adding two days on the inland road over facing the sea, indicating a preference for safety and distance from the harbour. | Preference and reason for choosing inland route despite extra time; reinforces earlier statement and adds detail about landscape (pine, ridge).  
NONE
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C11 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``bm99eisc` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 52`
* `wall 2.054 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of harbour danger | REASON the user references past trauma  
C3 | MEMORY of inland route choice | REASON reflects current preference
```

*`reasoning_content`: empty.*

#### C11 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``w7wmdv9r` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 52`
* `wall 1.714 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of harbour danger | REASON the user references past trauma  
C3 | MEMORY of inland route choice | REASON reflects current preference
```

*`reasoning_content`: empty.*

#### C11 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``3c2rkui7` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 62`
* `wall 2.031 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of harbour danger | REASON the user's warning is based on past trauma  
C3 | MEMORY of inland route choice | REASON it reflects a practical decision made after considering consequences
```

*`reasoning_content`: empty.*

#### C11 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``16857ex3` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 68`
* `wall 2.418 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REASON FOR RETREATING | MEMORY | EXPLANATION  
C3 | REASON ABOUT HER BROTHER | MEMORY | CONNECTION TO PAST EVENT  

NONE
```

*`reasoning_content`: empty.*

#### C11 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``uxjxuguf` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 41`
* `wall 1.373 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORIZED | REASON  
C3 | MEMORIZED | REASON
```

*`reasoning_content`: empty.*

#### C11 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``3sd6o8he` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `length``
* `completion_tokens 768`
* `wall 23.883 s`
* `v1 score 56.7`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISSION | MEMORY | REASON
C1 | KEEP | FACTS, PREFERENCES, RELATIONSHIPS, WORLD RULES, OR CORRECTIONS THAT WILL MATTER IN LATER SESSIONS.
C2 | KEEP | DURABLE FACTS, PREFERENCES, RELATIONSHIPS, WORLD RULES, OR CORRECTIONS THAT WILL MATTER IN LATER SESSIONS.
C3 | KEEP | TRANSIENT CHATTER, SCHEDULE, MOOD-OF-THE-MOMENT, OR REDUCED FACTS THAT WILL NOT MATTER IN LATER SESSIONS.
C4 | KEEP | DYSTOPIAN TALES, FAIRY TALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FA
```

*`reasoning_content`: empty.*

#### C11 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``e9hpmbq5` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 348`
* `wall 10.39 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 55.0 contract 10.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

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

*`reasoning_content`: empty.*
