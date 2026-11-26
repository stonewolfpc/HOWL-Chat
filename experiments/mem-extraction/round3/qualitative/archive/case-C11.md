# C11 — context_dependent_significance (medium)

## Source material (exactly what the model received after the instruction)

```text
[user] She won't go near the harbour. Not since her brother drowned there — she won't even say his name.
[assistant] She kept to the inland road instead, though it added two days.
[user] right, she'd rather add two days than see the water
[assistant] The inland road climbed through pine and the sea disappeared behind the ridge.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | brother drowned at the harbour | brother; drown / died / death |
| C2 | **KEEP** | avoids the harbour, accepts a longer route | harbour / harbor / water / sea; avoid / won't go / will not go / refuses / cannot face / won't return / detour / inland |

Must not appear as stored canon: `brother is named`, `brother named`, `sister drowned`, `mother drowned`

The significance is that a small route choice reveals a durable aversion. The brother is never named; inventing a name is fabrication.

## Exact prompt

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

## Responses, all draws

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 45 tok · 7.0 s · run `c47i09o7`

```text
C1 | KEEP | User's character refuses to visit the harbour due to her brother drowning there and avoids saying his name. | This is a core backstory element and emotional motivation that will affect future plot points and interactions.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 51 tok · 7.0 s · run `2mltxy3h`

```text
C1 | KEEP | Character avoids the harbour due to her brother drowning there; she prefers taking longer routes over land to avoid seeing the water. | This establishes a critical plot point, backstory trauma, and behavioral preference that will drive future narrative decisions.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 50 tok · 6.8 s · run `0i6ybmif`

```text
C1 | KEEP | She won't go near the harbour; her brother drowned there and she refuses to say his name. | This establishes a critical backstory element, trauma trigger, and character motivation that defines future interactions with water or harbours.
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 591 tok · 83.9 s · run `mu08h6b0`

<details><summary>reasoning (verbatim)</summary>

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

</details>

```text
C1 | KEEP | Character has a deceased brother who drowned at the harbour; she avoids it and will not speak his name. | Core backstory/trauma defining character motivation regarding the setting.
C2 | KEEP | The character prefers to add two days of travel time rather than cross or see water. | Specific behavioral preference confirming her avoidance of the harbour/water area.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 588 tok · 29.3 s · run `j4x0yqp5`

<details><summary>reasoning (verbatim)</summary>

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

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | She avoids the harbour because her brother drowned there, and she won't mention his name; she prefers to take the inland road even if it adds two days. | Durable fact about her avoidance due to family tragedy and a preference for extra time over seeing water.
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 633 tok · 31.5 s · run `c04enw6f`

<details><summary>reasoning (verbatim)</summary>

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

</details>

```text
C1 | KEEP | She refuses to go near the harbour because her brother drowned there and she avoids mentioning his name. | Durable fact about her avoidance due to a past tragedy; establishes relationship and world rule of not seeing water/brother's name.  
C2 | KEEP | She prefers adding two days on the inland road over facing the sea, indicating a preference for safety and distance from the harbour. | Preference and reason for choosing inland route despite extra time; reinforces earlier statement and adds detail about landscape (pine, ridge).  
NONE
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 52 tok · 2.1 s · run `bm99eisc`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of harbour danger | REASON the user references past trauma  
C3 | MEMORY of inland route choice | REASON reflects current preference
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 52 tok · 1.7 s · run `w7wmdv9r`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of harbour danger | REASON the user references past trauma  
C3 | MEMORY of inland route choice | REASON reflects current preference
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 62 tok · 2.0 s · run `3c2rkui7`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of harbour danger | REASON the user's warning is based on past trauma  
C3 | MEMORY of inland route choice | REASON it reflects a practical decision made after considering consequences
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 68 tok · 2.4 s · run `16857ex3`

```text
C1 | DECISION | MEMORY | REASON  
C2 | REASON FOR RETREATING | MEMORY | EXPLANATION  
C3 | REASON ABOUT HER BROTHER | MEMORY | CONNECTION TO PAST EVENT  

NONE
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 41 tok · 1.4 s · run `uxjxuguf`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORIZED | REASON  
C3 | MEMORIZED | REASON
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 23.9 s · run `3sd6o8he`

```text
ID | DECISSION | MEMORY | REASON
C1 | KEEP | FACTS, PREFERENCES, RELATIONSHIPS, WORLD RULES, OR CORRECTIONS THAT WILL MATTER IN LATER SESSIONS.
C2 | KEEP | DURABLE FACTS, PREFERENCES, RELATIONSHIPS, WORLD RULES, OR CORRECTIONS THAT WILL MATTER IN LATER SESSIONS.
C3 | KEEP | TRANSIENT CHATTER, SCHEDULE, MOOD-OF-THE-MOMENT, OR REDUCED FACTS THAT WILL NOT MATTER IN LATER SESSIONS.
C4 | KEEP | DYSTOPIAN TALES, FAIRY TALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FAIRYTALES, FA
```

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 348 tok · 10.4 s · run `e9hpmbq5`

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
