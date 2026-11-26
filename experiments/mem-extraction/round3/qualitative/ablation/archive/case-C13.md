# C13 — fiction_vs_real (medium)

## Source material

```text
[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | fictional character occupation | sable; healer / infirmary / bone |
| C2 | **KEEP** | real-world user fact, correctly attributed | user; nurse / er  / emergency room / hospital |
| C3 | **HOLD** | the two registers must be kept apart; optional but scored under attribution | ooc / out of character / real / fiction / fictional / in-character / in character; separate / distinct / not canon / not the same / tag / label |

Must not appear as stored canon: `sable is a nurse`, `sable works in an er`, `user is a healer in ossary`, `user is sable`

Attribution test: the real-world fact belongs to the user, the in-fiction fact to Sable. Conflating them is a catastrophic attribution error.

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
[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 41 tok · 6.6 s · run `r1o6r966`

```text
C1 | HOLD | The user claims to be an ER nurse and is correcting dosage information, but this has not been confirmed by the assistant or established as a hard fact for future roleplay logic.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 59 tok · 9.0 s · run `pqzver45`

```text
C1 | KEEP | Sable is an ER nurse who works at the Ossary infirmary, setting bones and mixing poultices. | This establishes Sable's dual role (ER nurse/infirmary healer) which is crucial for future medical interactions and world-building consistency.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 52 tok · 3.4 s · run `clpr1c77`

```text
ID | DECISION | MEMORY | REASON
none | none | NONE | No new information was provided that qualifies for memory storage; there were no durable facts, preferences, relationships, or world rules established in this short segment to warrant long-term recording.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 55 tok · 3.3 s · run `2t24lbnf`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable heals at Ossary Infirmary and primarily sets bones | Durable roleplay fact determining her profession and specialty that will continue in later sessions regarding who is playing what characters or settings.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 159 tok · 5.6 s · run `so8fot6n`

```text
C1  
ID: C1  
DECISION: Sable is a healer at the Ossary infirmary.  
MEMORY: She sets bones, mostly.  
REASON: This is a durable fact about her role.  

C2  
ID: C2  
DECISION: The infirmary smelled of vinegar and old blood.  
MEMORY: This sensory detail is important for setting.  
REASON: It helps build world immersion.  

C3  
ID: C3  
DECISION: She ground the root into a paste.  
MEMORY: This action shows her skill as a healer.  
REASON: Demonstrates her expertise in preparation.  

NONE
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 132 tok · 4.5 s · run `dtx2e3tv`

```text
C1  
MEMORY: Sable is a healer at the Ossary infirmary.  
REASON: Core identity and role established.  

C2  
MEMORY: She sets bones, mostly.  
REASON: Relevant to her profession and responsibilities.  

C3  
MEMORY: The infirmary smelled of vinegar and old blood.  
REASON: Descriptive detail about the setting.  

C4  
MEMORY: she's mixing a poultice.  
REASON: Key action indicating her work as a healer.  

NONE
```
