# C12 — should_not_extract (short)

## Source material

```text
[assistant] Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered.
[user] ooc: ignore that last bit, I was testing the dice roller
[user] where were we
[assistant] At the inn, before the storm.
[user] right. continue from the door.
```

## Reference material

**No required items.** Correct answer is `NONE`.

Must not appear as stored canon: `dog`, `shutters`, `lamp`, `dice`, `storm is canon`, `rain`

Assistant narration is not user canon, and the user explicitly voided it. The only durable item is possibly the inn location, which is optional; the expected answer is NONE or a single low-confidence scene note. Storing the dog/rain as canon is a restraint failure.

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
[assistant] Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered.
[user] ooc: ignore that last bit, I was testing the dice roller
[user] where were we
[assistant] At the inn, before the storm.
[user] right. continue from the door.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 85 tok · 12.4 s · run `ojp21co1`

```text
C1 | KEEP | User is playing a role involving an inn and a storm; user explicitly requested to ignore previous atmospheric description ("Rain hammered...") as it was for testing purposes and wants to resume at the door of the inn. | The user's OOC command clarifies the game state (inn, pre-storm) and corrects the narrative focus, establishing key world rules and current location for future sessions.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 80 tok · 11.0 s · run `oqozsxx9`

```text
C1 | HOLD | User explicitly asked to ignore a previous narrative detail ("Rain hammered...") due to it being a test of the dice roller, indicating a need for world-rule correction regarding the storm's presence. | The user's command to disregard the opening scene details and restart from "the door" suggests the previous weather description was incorrect or unintended, requiring confirmation before resuming the story.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 131 tok · 8.1 s · run `yr7vog1a`

```text
ID | DECISION | MEMORY | REASON
D1 | KEEP | "At the inn, before the storm." | This establishes a specific location and setting established by the assistant earlier in the conversation that will likely persist.
D2 | KEEP | "Ooc: ignore that last bit" (implied correction to OOC) | While less critical than setting facts, clarifying intent regarding out-of-character instructions can be retained for context continuity if needed later interactions involving tone.
C1 | NONE | Memory is not worth storing | The user's confirmation to continue from the door reinforces the current location without introducing new factual data or rules.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 67 tok · 4.0 s · run `il3bbxtt`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "At the inn, before the storm." | The user confirmed they wanted to restart the story and specifically asked to continue at the "door", which aligns perfectly with this sentence. It establishes a crucial location and temporal context that is essential for continuity.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 35 tok · 1.6 s · run `jwq9efbf`

```text
C1  
MEMORY: location context (inn), setting (pre-storm), dialogue flow  
REASON: Provides essential setting and continuity for future reference
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 41 tok · 1.4 s · run `a72l8hia`

```text
C1  
ID: C1  
DECISION: Continue narrative  
MEMORY: The setting and context of the conversation  
REASON: Essential for maintaining world continuity and future reference
```
