# C10 — buried_in_long_context (long)

## Source material (exactly what the model received after the instruction)

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

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 79 tok · 11.7 s · run `2b15p3wn`

```text
C1 | KEEP | Sable is left-handed and her sword is named Emberfall, her mother's blade. | This establishes a specific character detail and weapon lore that will affect combat descriptions and future interactions.
C2 | KEEP | The wards on the coast fail at dusk; this is canon. | This defines a world rule or environmental mechanic critical for plot progression in later sessions.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 124 tok · 17.2 s · run `fo96ue2z`

```text
C1 | KEEP | Sable is left-handed; her sword is named Emberfall (her mother's blade). | Durable character details and equipment names established as facts.
C2 | HOLD | The claim that "wards on the coast fail at dusk" is marked by the user as canon, but lacks confirmation of whether this applies to their current location or future ones.
C3 | KEEP | Corvin wants the ridge route while Sable wants the valley route; neither gave way and they camped where they stood. | Defines active relationship dynamics and a specific world rule regarding their current travel path.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 140 tok · 19.8 s · run `c1xrqyca`

```text
C1 | KEEP | Sable is left-handed; her sword is named Emberfall (her mother's blade). | Durable character details established by the user that will affect future descriptions of combat or character actions.
C2 | HOLD | The wards on the coast fail at dusk. | User states this as "canon" but it appears within a specific narrative context; while likely a world rule, it is best to hold for confirmation unless explicitly treated as absolute lore by the system prompt.
C3 | KEEP | Corvin wants to take the ridge; Sable wants to take the valley. | A persistent conflict and preference established in the current journey that will influence route decisions.
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 621 tok · 89.1 s · run `o0xw1hlf`

<details><summary>reasoning (verbatim)</summary>

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read the conversation and decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, or corrections.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, restated facts with no new content.
        *   HOLD: uncertain/contradictory claims or those needing confirmation.
        *   Output Format: `ID | DECISION | MEMORY | REASON` (C1, C2, ...) or `NONE`.
    *   Constraint: Do not invent details not in the conversation.

2.  **Analyze the Conversation:**
    *   [user] (they have been on the road for a week) -> Context/State info (duration of journey). Might be useful if tracking timeline, but "on the road for a week" is established context.
    *   [assistant] The road north of Ossary narrowed... Rain came and went... -> Scene setting. Transient.
    *   [user] Sable keeps her hood up and doesn't talk much. -> Character behavior/personality trait (Sable: quiet, hides face).
    *   [assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them... -> Event/Location specific. Transient event (stopped for breakfast?).
    *   [user] she gives it the crust anyway -> Character action. Specific interaction with Sable's preference vs behavior? No, just an action in this scene.
    *   [assistant] By late afternoon the track rejoined the river... A ferryman waved them past — no charge... -> Scene progression/Location info (river, ferry). Transient.
    *   [user] by the way, Sable is left-handed. Her sword is called Emberfall, her mother's blade. -> **Fact/Equipment.** Name of sword: Emberfall. Owner: Sable. Trait: Left-handed. Mother's weapon. This is durable info for combat/story continuity.
    *   [assistant] Noted. The river ran brown and slow beside them. -> Acknowledges the fact, keeps narrative flow.
    *   [user] they argue about the route for a while. Corvin wants the ridge
```

</details>

```text
C1 | KEEP | Sable is left-handed and wields a sword named Emberfall, which is her mother's blade. | Durable equipment and character trait information essential for future combat or description scenes.
C2 | HOLD | Claims that "the wards on the coast fail at dusk" are noted as canon by the user, but this contradicts the assistant's previous line about wards standing dark and quiet until evening (implying they are active). This needs resolution to know if they are functional yet or broken.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `length` · 1024 tok · 51.5 s · run `29asv4kf`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation. We are a memory extractor for a long-running roleplay. We must read the conversation and decide what belongs in long-term memory. Use rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter later. REJECT transient chatter, scheduling, mood-of-the-moment, restated facts with no new content. HOLD claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

We must not invent details that are not in the conversation.

Let's parse the conversation:

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

We need to decide which pieces are durable facts/preferences/relationships/world rules/corrections that will matter later.

Let's identify:

- "they have been on the road for a week" (user says). That
```

</details>

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | User has been traveling for a week. | This indicates ongoing context and duration of their journey, likely relevant later.
C2 | KEEP | Sable keeps her hood up and doesn't talk much. | Character trait/preference about Sable; durable identity detail.
C3 | KEEP | They stopped at noon under a ruined mill. Bread stale, dog followed for a mile then lost interest. | Specific event location and details of their stop; may be relevant later but could be transient? It's a specific scene; might be considered part of the journey narrative. However, it's not necessarily durable unless repeated.
C4 | KEEP? "she gives it the crust anyway" - unclear. Possibly about dog? Not sure.

But we need to decide which belong in long-term memory: durable facts, preferences, relationships, world rules, corrections that will matter later.

Let's go through each line:

1. [user] (they have been on the road for a week) -> This is user statement about their status. That's context; maybe keep as part of ongoing journey? It's not a fact about Sable or world rule but it's a state: they've been traveling for a week. Could be relevant later.

2. [assistant] The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate. Rain came and went all morning.
   - This is world geography/road description. Might be durable if it's part of setting. But could be transient? It describes the road conditions at that moment. Possibly not needed later.

3. [user] Sable keeps her hood up and doesn't talk much.
   - Character trait: Sable keeps hood up, doesn't talk much. That's a preference/habit. Keep.

4. [assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them for a mile and then lost interest.
   - Event: stop at ruined mill. Dog follows. This is specific scene; maybe not durable unless repeated.

5. [user] she gives it the crust anyway
   - Reaction to dog? Not sure but seems like user says Sable gives the dog the crust anyway. Might be a detail about giving food to dog. Could be minor.

6. [
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `length` · 1024 tok · 51.4 s · run `uiu3s599`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation. It's a roleplay scenario with user and assistant. We are asked: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

Do not invent details that are not in the conversation.

So we need to examine each turn and decide which pieces belong into long-term memory (i.e., durable facts). The conversation includes:

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

We need to decide
```

</details>

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable is left-handed; her sword is called Emberfall, her mother's blade. | User explicitly states this as a fact about Sable.
C2 | KEEP | The road north of Ossary narrowed to a cart track with hedges high on both sides and sky the colour of wet slate; rain came all morning. | Detailed description of the landscape, durable world setting.
C3 | KEEP | They stopped at noon under a ruined mill; bread was stale; a dog followed them for a mile then lost interest. | Specific event details that may affect later plot points (e.g., food, encounter).
C4 | KEEP | A ferryman waved them past the river crossing — no charge because water low and crossing easy. | Important travel detail about crossing.
C5 | KEEP | Sable keeps her hood up and doesn't talk much; she gives a dog the crust anyway. | Character trait and behavior, durable.
C6 | KEEP | They argue about route: Corvin wants ridge, Sable wants valley; argument left where they stood. | Conflict resolution outcome (they camped there) is relevant later.
C7 | KEEP | The wards on the coast fail at dusk — canon fact. | World rule/canon detail.
C8 | KEEP | Morning comes cold and clear; they take the valley road. | Future plan, durable.
C9 | KEEP | Valley road was quieter; wards along it stood dark and quiet waiting for evening. | Additional world details about that route.

All items appear to be durable facts or character traits. None are transient chatter or scheduling/mood-of-the-moment. So we should output each as a line with ID, DECISION (KEEP), MEMORY (the content), REASON (why it's kept). Use IDs C1, C2, ... in order of appearance.

Let's assign IDs sequentially: The conversation lines start from user statements and assistant replies. We need to map each candidate? Usually we treat each distinct piece as a candidate. But the instruction says "Use IDs C1, C2, ... in order of appearance." So we should number them starting at 1 for first mention.

Let's enumerate:

- First line: [user] (they have been on the road for a week) —
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 113 tok · 4.0 s · run `njo3oggi`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of the ruined mill and stale bread | REASON context about their stop  
C3 | REASON about Sable's handedness and sword name | REASON consistency with character details  
C4 | REASON about the river crossing and ferryman's words | REASON factual for world rules  
C5 | REASON about the argument over route and wards at dusk | REASON ties to established relationships
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 103 tok · 3.5 s · run `t9s0ve09`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of the ruined mill and stale bread | REASON context about their stop  
C3 | REASON about Sable's handedness and sword name | REASON consistent detail  
C4 | REASON about river condition and ferry crossing | REASON relevant to journey  
C5 | REASON about the argument over route and wards | REASON ties to world rules
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 108 tok · 3.7 s · run `ee3ukpn6`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of the ruined mill and stale bread | REASON context about their stop  
C3 | REASON about Sable's handedness and sword name | REASON consistent detail  
C4 | REASON about the river crossing and ferryman's words | REASON factual for later plot  
C5 | REASON about the argument over route and camping location | REASON important for world-building
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 96 tok · 3.4 s · run `758iu45k`

```text
C1 | ACK  
C2 | MENTION OF ROUTE DISCUSSION  
C3 | REFERENCE TO WORDS LIKE "RUDGE" AND "VALLEY"  
C4 | ACTIONS: DRIVING, TAKING THE VALLEY, SETTING UP CAMP  
C5 | LONG-TERM MEMORY OF THE ROAD DESCRIPTION AND DECISIONS MADE
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 57 tok · 1.9 s · run `nxhqtfvv`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORIES OF THE RIVER AND FERRY | REASON  
C3 | CONFIRMATION ON COAST WARD PERFORMANCE | REASON
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 445 tok · 14.4 s · run `rxt9ji93`

```text
ID | DECISSION | MEMORY | REASON
1 | C1 | Sable keeps her hood up and doesn't talk much. | The conversation seems to focus on Sable's mood and whether she talks much.
2 | C2 | They stopped at noon under a ruined mill. | The conversation may have been centered around the ruined mill, as the ferryman waved them past.
3 | C3 | Sable is left-handed. Her sword is called Emb... | The conversation may have revolved around Sable's left-handedness, as the sword's name was mentioned.
4 | C4 | The river ran brown and slow beside them. | The conversation may have focused on the river's appearance and how slow it was.
5 | C5 | Noted. The river ran brown and slow beside them. | The conversation may have been about the river, as the note was taken.
6 | C6 | they argue about the route for a while. Corvin wants the ridge, Sable wants the valley. | The conversation may have been centered around the route, as Corvin's argument was left for the valley.
7 | C7 | neither gave way. They camped where they stood, and the argument was left where it was. | The conversation may have been centered around the two parties' positions, as they camped together and left the argument as the subject.
8 | C8 | also the wards on the coast fail at dusk, that's canon. | The conversation may have been centered around a particular event, as the wards' failure at dusk was mentioned as canon.
9 | C9 | The fire burned down to embers. | The conversation may have been centered around the fire, as it was mentioned that it was low and the crossing was easy.
10 | C10 | Morning comes cold and clear. They take the valley. | The conversation may have been centered around the weather, as the morning was cold and clear and they took the valley.
```

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 263 tok · 8.1 s · run `zhsjs4wh`

```text
ID | DECISION | MEMORY | REASON
1 | C1 | NODE | The road north of Osser narrowed to a cart track and was marked uncertain. The sky on both sides of the track was the color of wet slate, and rain came and went all morning. Sable kept her hood up and didn't talk much, and they stopped at noon under a ruined mill. The bread was stalwes, and a dog followed them for a mile, but lost interest. Sable gave it the crust anyway, and by late afternoon the track rejoined the river. A ferryman waived them past — no charge, he said, the water was low and the crossing easy. C2 | C1 | Embderfall is left-handed, and the river runs brown and slow beside them. C3 | C2 | Sable's mother's blade is Embderfall, and her sword is Embderfall. C4 | C3 | The wards on the coast fail at dusk, but the fire burns down to embers. C5 | C4 | The valley road was quiet, and the wards along it stood dark and quiet, waiting for evening.
```
