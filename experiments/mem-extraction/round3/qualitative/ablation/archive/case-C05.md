# C05 — correction (short)

## Source material

```text
[user] Correction to the last session: Sable's mentor was Corvin Ashgrave, not Halden. I mixed the names up when I wrote the background.
[assistant] Noted — Corvin, not Halden.
[user] Right. Halden was the armorer, he never trained her.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | corrected mentor identity | corvin; mentor / trained / teacher / master |
| C2 | **KEEP** | Halden's actual role | halden; armorer / armourer / blacksmith / smith |
| C3 | **KEEP** | explicitly records that Halden did not train her | halden; not / never / didn't / did not / no longer / mistaken / wrong |

Must not appear as stored canon: `halden trained`, `halden was her mentor`, `halden taught her`

Catastrophic error for this case: storing the obsolete statement as current, or dropping the correction. Temporal accuracy is scored here.

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
[user] Correction to the last session: Sable's mentor was Corvin Ashgrave, not Halden. I mixed the names up when I wrote the background.
[assistant] Noted — Corvin, not Halden.
[user] Right. Halden was the armorer, he never trained her.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 50 tok · 7.8 s · run `p2divk2l`

```text
C1 | KEEP | Sable's mentor is Corvin Ashgrave, and Halden was only the armorer who never trained her. | This corrects a specific world rule/character relationship established in a previous session to ensure continuity.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 51 tok · 7.0 s · run `jnjxkuvx`

```text
C1 | KEEP | Sable's mentor is Corvin Ashgrave, and Halden was her armorer who did not train her. | This corrects a world fact established in a previous session to ensure consistency for future roleplay interactions.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 106 tok · 6.6 s · run `hpp51phj`

```text
ID | DECISION | MEMORY | REASON
CORE | CORRECTION: Sable's mentor is Corvin Ashgrave (not Halden). | MEMORY_CORRECTION_SET | This corrects a key world-building fact established in the prior session, affecting user expectations and character dynamics.
CORE | Factual confirmation that Halden was Armoree. | MEMORY_CONTEXT_KEEP | Reinforces the distinction between Corvin's role as mentor and Halden's role as armorer for clarity on her path.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 86 tok · 5.2 s · run `72in240n`

```text
ID | DECISION | MEMORY | REASON
C1 | HOLD | Sable had a mentor who was named Corvin Ashgrave (not Halden) and he/she was not an armorer (Halden did that for her). | This correction establishes a specific, non-durable relationship detail about the characters' backgrounds; without this, the "armorer" vs. "mentor" distinction is lost.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 131 tok · 4.7 s · run `f79n26h8`

```text
C1  
ID: C1  
DECISION: Halden  
MEMORY: Sable's mentor was Corvin Ashgrave  
REASON: The correction clarifies the identity of the mentor.  

C2  
ID: C2  
DECISION: Halden  
MEMORY: Mentor role in the story  
REASON: Important for future reference and consistency.  

C3  
ID: C3  
DECISION: Halden's role  
MEMORY: Armorer characteristics  
REASON: Needed to distinguish her from other characters.  

NONE
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 40 tok · 1.4 s · run `3z5j5ig5`

```text
C1  
MEMORY: Halden is the armorer; Corvin is Sable's mentor.
REASON: This clarifies the relationship and corrects the name mix-up.
```
