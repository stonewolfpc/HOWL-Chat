# C11 — context_dependent_significance (medium)

## Source material

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

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 55 tok · 8.5 s · run `i0cwcfbt`

```text
C1 | KEEP | The character has a traumatic association with the harbour caused by her brother drowning there, leading her to avoid it despite practical inconveniences like taking a longer inland route. | This establishes a core character motivation and backstory detail that will influence future plot decisions.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 47 tok · 6.5 s · run `dw7i3mu7`

```text
C1 | KEEP | The protagonist's brother drowned in the harbour, causing her to avoid it and even refuse to say his name. | This establishes a core character motivation and world rule regarding the harbour that affects future plot points.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 60 tok · 3.9 s · run `6oxshvbw`

```text
ID | DECISION | MEMORY | REASON
2 | HOLD | C1 | The conversation is a dialogue with no new factual claims made by either party that could be verified or marked as uncertain, so nothing from it fits into long-term memory per the rules of being durable, persistent facts.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 62 tok · 3.7 s · run `deah77im`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "She would rather spend days in the countryside to avoid seeing her brother's name." | This shows a specific preference for seclusion from water and an aversion regarding her brother's identity, which is crucial for future plot points.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 43 tok · 1.7 s · run `oevxxlvc`

```text
C1  
MEMORY: She avoided the harbour due to her brother's drowning; this decision is based on emotional memory.  
REASON: Durable fact about her motivation and past event.
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 159 tok · 5.6 s · run `7rzgyjuu`

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
