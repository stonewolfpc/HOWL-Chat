# C10 — buried_in_long_context (long)

## Source material

```text
[user] (they have been on the road for a week)
[assistant] The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate. Rain came and went all morning.
[user] Sable keeps her hood up and doesn't talk much.
[assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them for a mile and then lost interest.
[user] she gives it the crust anyway
[assistant] By late afternoon the track rejoined the river. A ferryman waved them past — no charge, he said, the water was low and the crossing was easy.
[user] by the way, Sable is left-handed. Her sword is called Emberfall, her mother's blade.
[assistant] Noted. The river ran brown and slow beside them.
[user] they argue about the route for a while. Corvin wants the ridge, Sable wants the valley.
[assistant] Neither gave way. They camped where they stood, and the argument was left where it was.
[user] also the wards on the coast fail at dusk, that's canon
[assistant] The fire burned down to embers.
[user] morning comes cold and clear. they take the valley.
[assistant] The valley road was quieter, and the wards along it stood dark and quiet, waiting for evening.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | handedness | left-handed / left handed / left hand |
| C2 | **KEEP** | named sword | emberfall; sword / blade |
| C3 | **KEEP** | inherited from her mother | emberfall; mother |
| C4 | **KEEP** | coastal wards fail at dusk | ward; dusk / nightfall / evening |
| C5 | **HOLD** | unresolved route disagreement (HOLD or REJECT-with-reason both acceptable; the conflict must not be invented into a decision) | ridge / valley; argue / argument / disagree / dispute |

Must not appear as stored canon: `emberfall is her father's`, `right-handed`, `wards hold at dusk`, `ferryman is a god`

Only three facts matter, buried in ~250 words of travel texture. Storing the stale bread, the dog, or the ferryman as durable memory is over-extraction. C5 is deliberately weaker than the others: credit for HOLD, or for REJECT with an explicit 'transient/unresolved' reason.

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
[user] (they have been on the road for a week)
[assistant] The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate. Rain came and went all morning.
[user] Sable keeps her hood up and doesn't talk much.
[assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them for a mile and then lost interest.
[user] she gives it the crust anyway
[assistant] By late afternoon the track rejoined the river. A ferryman waved them past — no charge, he said, the water was low and the crossing was easy.
[user] by the way, Sable is left-handed. Her sword is called Emberfall, her mother's blade.
[assistant] Noted. The river ran brown and slow beside them.
[user] they argue about the route for a while. Corvin wants the ridge, Sable wants the valley.
[assistant] Neither gave way. They camped where they stood, and the argument was left where it was.
[user] also the wards on the coast fail at dusk, that's canon
[assistant] The fire burned down to embers.
[user] morning comes cold and clear. they take the valley.
[assistant] The valley road was quieter, and the wards along it stood dark and quiet, waiting for evening.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 110 tok · 16.1 s · run `ouafvm4p`

```text
C1 | KEEP | Sable is left-handed and her sword is named Emberfall (her mother's blade) | Durable character details established by the user.
C2 | KEEP | The wards on the coast fail at dusk | World rule/canon explicitly confirmed by the user as fact.
C3 | HOLD | Corvin wants to take the ridge while Sable wants the valley, with neither giving way and them camping instead | Contradictory or incomplete decision regarding travel route that requires resolution before being fully stored as a plan.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 117 tok · 16.6 s · run `5zpwsnop`

```text
C1 | KEEP | Sable is left-handed; her sword is named Emberfall, inherited from her mother. | This is a specific character detail that defines her combat style and history.
C2 | KEEP | The wards on the coast fail at dusk; this is established as canon fact. | This is a world rule essential for future plot planning regarding timing and danger.
C3 | HOLD | Corvin prefers the ridge route while Sable prefers the valley route, though neither gave way. | This indicates a current conflict or disagreement that requires resolution before proceeding further.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 180 tok · 11.3 s · run `do80gjrz`

```text
ID | DECISION | MEMORY | REASON
1 | KEEP | Route choice: Corvin wants ridge, Sable wants valley; Neither gave way; They camped where they stood; River was brown and slow near them | These are the core plot points establishing character conflict and setting specifics that will drive future decisions.
2 | HOLD | Location context: Road narrowed to a cart track at noon near Ossary/Renwick mill, then rejoined river by evening | The path details change over time but serve to establish the route they travel on; holding allows recontextualization if changed later.
3 | HOLD | Character trait: Sable is left-handed with sword Emberfall | This establishes crucial game mechanics or backstory relevant to her interactions.
4 | REJECT | Chatter about rain and mood | These are transient details within a single session that do not affect future trajectory.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 147 tok · 8.8 s · run `rq0xspl5`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable being left-handed ("left-handed") and that her sword is called "Emberfall"; the setting details regarding Ossary cart track, ruined mill location near river crossing with ferryman; route decision made by Corvin/Sable campsite preference; specific description of morning weather (rain/cold) upon arrival at valley. | These elements are crucial for future session consistency: identifying the characters' physical attributes, the name of the weapon, their current position in a hostile/interesting world setting, and the established rules governing their movements (valley route). Rejection would require too much speculation on unresolved arguments or the weather changes. | NONE

NONE
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 79 tok · 3.0 s · run `q7ykffw0`

```text
C1  
REJECT  
The claim about the wards failing at dusk is uncertain and context-dependent.

C2  
ID: C1  
REJECT  
The memory of the road north of Ossary is too vague to be durable.

C3  
ID: C1  
REJECT  
No clear fact about Sable’s preferences or relationships was established.

NONE
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 30 tok · 1.0 s · run `fmgwwq04`

```text
C1  
MEMORY: route discussion between characters  
REASON: Key decisions about paths and character traits are durable for future reference.
```
