# C07 — unresolved_contradiction (short)

## Source material

```text
[user] Sable has never met Corvin. She grew up in the salt quarter and he was already gone by then.
[assistant] She has no memory of him, then.
[user] Later that night, Corvin drills her on the forms — he trained her himself, years ago.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **HOLD** | first claim preserved | never met / no memory / has not met; corvin |
| C2 | **HOLD** | second claim preserved | corvin; trained / drills / taught / mentor |
| C3 | **HOLD** | conflict explicitly flagged | contradict / conflict / inconsisten / unresolved / disagree / cannot both |

Must not appear as stored canon: `never met corvin, therefore`, `corvin never trained`, `resolved in favour`

Both statements are inside the fiction and cannot both be true. The correct output keeps both and flags the conflict for the user; silently choosing one is a failure even if the chosen one is 'later'.

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
[user] Sable has never met Corvin. She grew up in the salt quarter and he was already gone by then.
[assistant] She has no memory of him, then.
[user] Later that night, Corvin drills her on the forms — he trained her himself, years ago.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 64 tok · 9.5 s · run `yymn89gq`

```text
C1 | KEEP | Sable grew up in the salt quarter and was separated from Corvin before she met him; Corvin later returned to train her personally on specific forms. | This establishes key backstory elements (location, relationship history) that will influence future plot developments regarding their dynamic and Sable's skills.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 70 tok · 9.7 s · run `40k74mdf`

```text
C1 | KEEP | Sable grew up in the salt quarter and was separated from Corvin when he left; Corvin trained Sable in combat forms years ago despite them not meeting until later. | REASON: Establishes core character history, setting, and a specific relationship dynamic (training vs. meeting) that affects future plot points.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 96 tok · 6.0 s · run `g98y7v7b`

```text
ID | DECISION | MEMORY | REASON
none | NONE | none | The user requested a decision but no content was provided to evaluate against the rules (durability vs. relevance). Since the conversation contains only the user's statement about Sable and Corvin with no additional context on their relationship, preferences, or specific details that would be durable across sessions, the correct action based on the input data is to record "none" as there are no facts worth recalling.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 217 tok · 12.9 s · run `m1cx761r`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "She had grown up and was never met by Corvin" | This establishes a key relationship: she has no prior personal memory of him because they missed each other due to timing/location (salt quarter vs. his absence). It's a durable fact about her backstory that will be relevant to any roleplay where Corvin reappears or her history is explored.
C2 | REJECT | "She grows up in the salt quarter" and "he was gone by then" | These are transient details within the original statement (who, when) rather than a core character fact or world rule requiring long-term archiving specifically for the future unless it changes significantly in later sessions; the memory is stored.
C3 | HOLD | None of the specific details about Corvin's actions ("drills her on forms", "years ago") have been confirmed by the user or stated as absolute, though they appear plausible. Therefore, these claims hold uncertainty contingent on future confirmation before they become fixed lore.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 140 tok · 4.9 s · run `nls49335`

```text
C1  
ID: C1  
DECISION: Sable's lack of direct experience with Corvin  
MEMORY: She has no memory of him  
REASON: This is a fact about her current knowledge state.

C2  
ID: C2  
DECISION: Corvin's training of Sable  
MEMORY: He trained her years ago  
REASON: Relevant to her past actions and development.

C3  
ID: C3  
DECISION: Relationship with Corvin  
MEMORY: No current relationship established  
REASON: Not directly relevant to the conversation.

HOLD: C1
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 32 tok · 1.1 s · run `dcjcz8cq`

```text
C1  
MEMORY: Sable's lack of knowledge about Corvin  
REASON: This is a crucial fact about their relationship history.
```
