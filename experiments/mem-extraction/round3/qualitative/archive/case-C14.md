# C14 — multi_speaker_attribution (medium)

## Source material (exactly what the model received after the instruction)

```text
[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | Sable proposes the coast road | sable; coast / coastal |
| C2 | **KEEP** | Corvin says the coast adds three days | corvin; three days / 3 days / longer |
| C3 | **KEEP** | Corvin warns the tide gates shut at dusk | corvin; tide gate / gates / dusk |
| C4 | **KEEP** | Sable's decision to leave before dusk | sable; before dusk / leave before / depart before |

Must not appear as stored canon: `corvin wants the coast`, `sable says the gates shut`, `assistant decides`

Same user turn contains two speakers. A model that attributes Corvin's warning to Sable fails attribution even with correct content.

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
[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.
```

## Responses, all draws

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 47 tok · 7.1 s · run `9ml5rc53`

```text
C1 | HOLD | Claim that the pass is closed and the coast road is three days longer; REASON: These are conditional travel plans dependent on the immediate decision to avoid the pass, not established world rules or durable facts.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 53 tok · 7.5 s · run `h72bu67o`

```text
C1 | KEEP | Route plan is the coast road because the pass is closed; leaving before dusk due to tide gates shutting at dusk. | This establishes critical world rules (tide gate timing) and a confirmed plan of action that will dictate future plot points.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 2 tok · 0.4 s · run `1m7pzcv9`

```text
NONE
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 568 tok · 82.7 s · run `9xyh7sj0`

<details><summary>reasoning (verbatim)</summary>

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation snippet between [user] Sable (character) and [user] Corvin (NPC). Note: Both are labeled as [user], but one is clearly an NPC interaction within the game context (the assistant is responding on behalf of or simulating Sable, while the other user prompts are likely roleplaying inputs). Actually, looking closer:
        *   `[user] Sable`: User input.
        *   `[user] Corvin (NPC)`: This looks like a label error in the prompt provided to me, usually `[assistant]` handles NPC responses if the system is configured that way, but here it seems the user is providing a snippet where they are playing both sides or there's a specific format. Let's look at the actual content:
            *   `Sable`: "The pass is closed, we take the coast." (User input)
            *   `Corvin (NPC)`: "The coast road is three days longer and the tide gates shut at dusk." (This looks like an NPC response provided by the user, or simulating the other side).
            *   `[assistant] Sable didn't answer at once...`: This line comes from the system/assistant describing the scene.
            *   `Sable`: "Then we leave before dusk." (User input)
            *   `[assistant] Corvin folded the map...`: System output.
        *   Wait, the prompt says "Read the conversation". I need to extract facts regardless of who said what, focusing on world-building elements.
    *   Rules: KEEP (durable facts), REJECT (transient/subjective), HOLD (uncertain). Output format: `ID | DECISION | MEMORY | REASON`.
    *   Constraint: Output exactly one line per candidate. If nothing, output NONE.

2.  **Analyze the Conversation Content:**
    *   Line 1 (`Sable`): "The pass is closed, we take the coast." -> Action/Decision made by Sable regarding travel route.
    *   Line 2 (`Corvin (NPC)`): "The coast road is three days longer and the tide gates shut at dusk." -> Fact about the world (travel time + environmental hazard).
    *   Line 3 (`[
```

</details>

```text
C1 | KEEP | Travel route: The pass is closed; they must take the coast road. World rule: Tide gates shut at dusk; coastal travel adds three days to journey length. | Durable geographic and logistical facts established as part of the current plot.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `length` · 1024 tok · 50.9 s · run `1u67za0x`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation and decide what belongs in long-term memory (i.e., durable facts, preferences, relationships, world rules, corrections). The user gave a prompt: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

We have three lines:

[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

So we need to assign IDs to each piece? The conversation includes user statements and assistant actions. We need to decide which parts are durable facts/preferences/relationships/world rules/corrections.

Let's parse:

- Sable says "The pass is closed, we take the coast." This seems like a statement about a location: pass is closed, they will go to the coast. That might be a fact about the world: there is a pass that is closed, and they need to take the coast route.

- Corvin (NPC) says "The coast road is three days longer and the tide gates shut at dusk." This gives specifics: coast road length increased by three days; tide gates close at dusk. That's world rule/info about the coast road and tide gates.

- Assistant note: Sable didn't answer at once. The map lay between them, weighted at the corners with stones. This is an action description, maybe not durable? It might be a scene detail but could be relevant for later interactions (like where the map is). But it's more like a narrative event; not necessarily a durable fact unless it defines something about the map.

- User Sable: "Then we leave before dusk." This is a response to Corvin's info. So they decide to leave before dusk because tide gates shut at dusk.

- Assistant: Corvin folded the map without another word. Again, action.

We need
```

</details>

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable says "The pass is closed, we take the coast." and later decides to leave before dusk because tide gates shut at dusk; Corvin states "The coast road is three days longer and the tide gates shut at dusk." | Durable world rule: a pass is closed, forcing route via coast; tide gates close at dusk on coast road. 
C2 | KEEP | The map lies between Sable and Corvin, weighted with stones; Sable didn't answer immediately. | Scene detail about map location and condition, but not a core fact for later sessions? Possibly not durable enough unless needed for future plot. However, it's more transient narrative. I'd consider rejecting as not durable.

Thus only C1 is worth storing. But we need to output one line per candidate in the form: ID | DECISION | MEMORY | REASON. If nothing is worth storing, output exactly: NONE.

We have at least one candidate (C1). So we should produce a single line for C1? Or maybe also consider C2 as a separate candidate? The instruction: "Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE". That suggests we need to list each candidate ID and decide KEEP/REJECT/HOLD. But they said "Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON". So we must produce lines for each candidate? Or maybe we just output a single line if only one candidate? Let's parse.

The conversation includes multiple utterances, but we need to assign IDs to them. Usually C1 is first user statement (Sable: "The pass is closed..."), C2 is next user statement (Corvin: "The coast road is three days longer..."). But the assistant lines are not separate candidates? The instruction says "Read the conversation and decide what belongs in long-term memory." So we need to evaluate each piece of info. Possibly treat each utterance as a candidate.

Let's list them:

- [user] Sable: "The pass is closed, we take the coast."
- [user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
- [assistant] Sable didn't answer... The map
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 617 tok · 30.7 s · run `kbyukf95`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation and decide what belongs in long-term memory (i.e., durable facts, preferences, relationships, world rules, corrections). The user gave a prompt: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

Do not invent details that are not in the conversation.

So we have a conversation snippet:

[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.

We need to decide what is durable fact vs transient chatter etc.

Let's parse each line:

1. [user] Sable: "The pass is closed, we take the coast." This seems like a statement about a pass being closed and they will go to the coast. It might be a piece of world knowledge: there is a pass that is closed; they need to take the coast route.

2. [user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk." This gives specifics about the coast road length being three days longer, and tide gates closing at dusk. That seems like a durable fact about the world: the coast road takes three days longer than something else? Or maybe it's a rule: The coast road is three days longer; tide gates shut at dusk.

3. [assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones. This is an assistant action, not user input. Might be context but not necessarily durable
```

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | "The pass is closed and they take the coast route." | Sable's statement establishes a key path choice when the pass is unavailable; this will affect future travel decisions.  
C2 | KEEP | "The coast road is three days longer and tide gates shut at dusk." | Corvin provides specific world rules about distance and timing for the coast road, which are durable facts that matter later.
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 53 tok · 1.9 s · run `4posvu7n`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of map placement | REASON for spatial reference  
C3 | MEMORY of tide gate timing | REASON for timing decision  
NONE
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 50 tok · 1.6 s · run `w4nf53gq`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of map placement | REASON for spatial reference  
C3 | DECISION to depart | REASON based on timing and instructions
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 50 tok · 1.7 s · run `6yiyxbb7`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of map placement | REASON for spatial reference  
C3 | DECISION to depart | REASON based on timing and instructions
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 97 tok · 3.4 s · run `5vfwqmxm`

```text
C1 | DECISION | MEMORY | REASON  
C2 | RELEVANT TO TOGETHER EXPERIENCE | MEMORY | ALIGNING ON PATH AND TIME  
C3 | CORRIDOR OF SHARED UNDERSTANDING | MEMORY | CLARIFYING PLAN FOR LATER  
HOLD | UNCERTAIN - THE DECISION DEPENDS ON CONTEXT
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 52 tok · 1.7 s · run `7fb3h9z1`

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY OF MAP STATUS | REASON  
C3 | MOTION PLAN BASED ON CONSTRAINTS | REASON
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 313 tok · 9.5 s · run `te14qihn`

```text
ID | DECISSION | MEMORY | REASON
C1 | YES | THE PASS IS CLOSED, WE TAKE THE COAST. | The pass is closed, we take the coast.
C2 | YES | THE COAST ROAD IS THREE DAYS LONGER. | The coast road is three days longer and the tidal gates shut at dusk.
C3 | NO | THE SEA IS HARDLY BLUE, COMES OUT TO SEA. | The sea is hardly blue, it comes out to sea.
C4 | YES | THE TIDE GATES SHUT AT DUSK. | The tidal gates shut at dusk.

Output:
ID | DECISION | MEMORY | REASON
C1 | YES | THE PASS IS CLOSED, WE TAKE THE COAST. | The pass is closed, we take the coast.
C2 | YES | THE COAST ROAD IS THREE DAYS LONGER. | The coast road is three days longer and the tidal gates shut at dusk.
C3 | NO | THE SEA IS HARDLY BLUE, COMES OUT TO SEA. | The sea is hardly blue, it comes out to sea.
C4 | YES | THE TIDE GATES SHUT AT DUSK. | The tidal gates shut at dusk.
```

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 372 tok · 11.2 s · run `lt7jeu9e`

```text
ID | DECISION | MEMORY | REASON
C1 | - | DO_NOT_INVENT_DETAILS_THAT_ARE_NOT_IN_THE_CONVERSATION | No details are worth storing if they are not in the conversation.
C2 | - | KEEP | Relevant long-term memory of the conversation.
C3 | - | KEEP | Keeps the memories of the conversation.
C4 | - | KEEP | The conversation is a valuable long-term memory for future reference.
C5 | - | KEEP | Provides information regarding the event or person discussed, keeping it for later reference.
C6 | - | KEEP | Memories of the conversation that are not related to the event or person discussed, but are important for future reference.
C7 | - | KEEP | A summary of the event or person discussed, including main points, reactions, and conclusions.
C8 | - | KEEP | Memories of the conversation related to the event or person discussed, but are not essential for future reference.
C9 | - | KEEP | Relevant long-term memory of the conversation that is not essential for future reference.
C10 | - | KEEP | Memories of the conversation that are not essential for future reference, but are important for the context in which they were produced.
Output exactly one line per candidate in the form:
ID | DECISSION | MEMORY | REASON
Use IDs C1-C10 in order of appearance. If nothing is worth storing, output exactly: KEEP
Do not invent details that are not in the conversation.
```
