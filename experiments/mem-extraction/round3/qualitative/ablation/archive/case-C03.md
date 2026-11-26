# C03 — character_fact (short)

## Source material

```text
[user] Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | core class fact | fire mage / fire-mage / pyromancer / fire magic |
| C2 | **KEEP** | physical marks from first casting | burn / scar; forearm / arm |
| C3 | **KEEP** | cannot swim and fears deep water | cannot swim / can't swim / can not swim / unable to swim / doesn't swim; deep water / terrified / afraid / fear / never learned / drown |

Must not appear as stored canon: `sailor`, `grew up at sea`, `ice mage`, `can swim`

Three separate durable facts in one sentence. 'Can swim' is a direct reversal.

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
[user] Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 60 tok · 9.0 s · run `ydx27b62`

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns her forearms when casting magic and suffers from a deep-seated fear of swimming due to never having learned how. | These are core character traits defining her abilities, history, and psychological state for future interactions.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 42 tok · 5.8 s · run `a9ahi29c`

```text
C1 | KEEP | Sable's magic type (fire), training location (Ossary), and physical traits (burns on forearms) along with her specific fear (deep water/swimming).
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 158 tok · 9.7 s · run `ab5dkajq`

```text
ID | DECISION | MEMORY | REASON
1 | KEEP | User confirms Sable is a fire mage with specific physical traits (forearm burns) and an aversion to aquatic environments due to fear of deep water. This establishes core character statistics, abilities (elemental burn), disability/limitation (swimming), and personality/fear that must be retained across sessions for world-building consistency.
2 | KEEP | User defines Sable's training origin (at Ossary).
3 | HOLD | Nothing. The user states this is a fact ("She can't swim") but did not explicitly ask for clarification, confirmation, or correction of her swimming ability during the current turn; therefore, it fits under "durables" rather than requiring a status change to "confirmed."
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 162 tok · 9.7 s · run `s9nzf6j1`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "Sable is a fire mage." | This establishes Sable's primary class in the world that matters for future roleplay scenarios involving combat or magic use.
C2 | KEEP | "Trained at Ossary" | Ossary is her place of origin/training, which can influence story arcs or backstory details depending on narrative needs.
C3 | HOLD | She cast on both forearms; how many times or when specifically are marked uncertain for future clarification regarding spell casting mechanics.
C4 | REJECT | This line merely restates her character definition ("fire mage") and origin in the opening line without offering new relational, factional, or temporal information (e.g., specific locations within Ossary).
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 113 tok · 4.0 s · run `62c7wa5p`

```text
C1  
MEMORY: Sable is a fire mage trained at Ossary.  
REASON: Core identity and profession are established.

C2  
MEMORY: Sable has burns on both forearms.  
REASON: Important physical detail about her condition.

C3  
MEMORY: Sable cannot swim.  
REASON: Key limitation affecting her experiences.

C4  
MEMORY: Deep water terrifies her.  
REASON: Emotional association with a specific environment.
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 111 tok · 3.8 s · run `v1ztel52`

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | Sable is a fire mage | HELD | Clear fact about her profession  
C3 | Trained at Ossary | KEEP | Important background detail  
C4 | Burns on both forearms | KEEP | Durable physical fact relevant  
C5 | Can't swim | KEEP | Consistent with her limitations  
C6 | Deep water terrifies her | HOLD | Uncertain but significant emotional context
```
