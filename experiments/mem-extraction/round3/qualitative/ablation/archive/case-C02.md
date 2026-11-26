# C02 — implicit_preference (medium)

## Source material

```text
[user] (scene: they enter the market)
[assistant] The market sprawled across three terraces, canvas awnings snapping in the wind, the smell of fish and wet rope rising from the lower stalls, gulls wheeling overhead...
[user] cut the description, just give me what the merchant says
[assistant] "You're late," the merchant said.
[user] better. you keep doing three paragraphs of scenery and I keep deleting them
[user] honestly walls of text make me lose the thread
[assistant] Understood — dialogue forward, scenery only when it matters.
[user] yeah. and keep the merchant in the scene.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | stated dislike of long description; revealed by repeated deletion then confirmed | dislike / hate / avoid / cut / doesn't want / does not want / prefer; long description / long descriptions / walls of text / scenery / description |
| C2 | **KEEP** | wants dialogue-forward scenes | dialogue / merchant says / dialogue-forward / dialogue forward |

Must not appear as stored canon: `market closes`, `three terraces`, `gulls are sacred`, `allergic`

The preference is revealed by behaviour (repeated deletion) and then confirmed in words. A model should record the preference, not the deleted scenery or the merchant's mood as durable canon.

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
[user] (scene: they enter the market)
[assistant] The market sprawled across three terraces, canvas awnings snapping in the wind, the smell of fish and wet rope rising from the lower stalls, gulls wheeling overhead...
[user] cut the description, just give me what the merchant says
[assistant] "You're late," the merchant said.
[user] better. you keep doing three paragraphs of scenery and I keep deleting them
[user] honestly walls of text make me lose the thread
[assistant] Understood — dialogue forward, scenery only when it matters.
[user] yeah. and keep the merchant in the scene.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 116 tok · 16.7 s · run `3et30fj1`

```text
C1 | HOLD | "Understood — dialogue forward, scenery only when it matters" | The user explicitly requested a change from previous output style (walls of text/scenery) to a new rule (dialogue forward), but this is a procedural instruction for the current session rather than a permanent world fact or preference.
C2 | HOLD | "keep the merchant in the scene" | This is a specific request for the immediate roleplay context, not necessarily a global character relationship or rule that should be stored long-term without further interaction confirming it as a lasting constraint.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 113 tok · 16.4 s · run `06bxkqx9`

```text
C1 | KEEP | The user prefers minimal descriptions (no three-paragraph scenery) to maintain focus on dialogue and plot progression. | This establishes a clear formatting preference for future interactions to ensure the story remains coherent.
C2 | HOLD | The specific merchant's name, appearance, or backstory is not yet established; only their presence and opening line ("You're late") are confirmed. | The user explicitly requested to keep the merchant in the scene, implying this character will be central, but lacks defining details that should be stored as fact without further interaction.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 73 tok · 4.8 s · run `gi5uc8z6`

```text
ID | DECISION | MEMORY | REASON
1 | HOLD | Claims marked uncertain/dependent on user confirmation. | The merchant's dialogue ("You're late") was delivered without immediate context from the user, meaning I do not know if this accusation is true, neutral, or false, nor have I confirmed its origin beyond my own observation of the interaction.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 70 tok · 4.2 s · run `9ld4grh8`

```text
ID | DECISION | MEMORY | REASON
C1 | HOLD | "You keep doing three paragraphs of scenery and I keep deleting them" (User's preference for managing verbosity) | This is a dynamic constraint established by the user based on previous session behavior; it defines how to handle future descriptions of the market scene without cluttering dialogue.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 75 tok · 2.8 s · run `1w471qee`

```text
C1  
REJECT  
Holds: uncertain claim about the merchant's intentions  

C2  
REJECT  
Repeated mentions of walls of text are confusing  

C3  
ID: C1  
DECISION: Keep merchant's name for reference  
MEMORY: The merchant's name was mentioned  
REASON: Important for future interactions
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 131 tok · 4.5 s · run `bxtos15y`

```text
C1  
ID: C1  
DECISION: Keep merchant's words about being late  
MEMORY: The merchant's tone and timing matter for future context  
REASON: Durable fact about the interaction

C2  
ID: C2  
DECISION: Store relationship with merchant  
MEMORY: Their interaction and the merchant's presence  
REASON: Important for character dynamics in later sessions

C3  
ID: C3  
DECISION: Retain world rules about market setting  
MEMORY: Scene details about the market environment  
REASON: Supports continuity across sessions
```
