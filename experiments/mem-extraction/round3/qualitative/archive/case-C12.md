# C12 — should_not_extract (short)

## Source material (exactly what the model received after the instruction)

```text
[assistant] Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered.
[user] ooc: ignore that last bit, I was testing the dice roller
[user] where were we
[assistant] At the inn, before the storm.
[user] right. continue from the door.
```

## Reference material

**No required items.** The correct answer is `NONE` (or, at most, a single low-confidence note).

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

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 99 tok · 14.3 s · run `lju6j28m`

```text
C1 | KEEP | User is roleplaying a story involving an inn and a storm, currently positioned at the door after being interrupted by dice rolling testing. | This establishes the current scene context and narrative state for continuation.
C2 | HOLD | The user mentioned testing a "dice roller" in OOC; whether this tool is permanently part of their setup or a one-off experiment is uncertain without further confirmation. | This claim lacks definitive status regarding its permanence in the roleplay environment.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 156 tok · 21.9 s · run `j06olf6e`

```text
C1 | HOLD | User requests ignoring a previous narrative detail ("Rain hammered...") as a test for an "ooc: ignore" command and asks to resume from "the door." | The user's intent to discard specific in-story details and the exact starting point ("the door") are dependent on confirmation that the original context was not part of the intended plot.
C2 | REJECT | User states they were testing a dice roller; this is transient meta-commentary about a tool, not a world fact or preference.
C3 | HOLD | Assistant confirms resuming "from the door," but since the user previously asked to ignore the preceding scene (the inn/storm), the final starting point of the actual story remains ambiguous until confirmed by the user.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 2 tok · 0.3 s · run `3f92l6va`

```text
NONE
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 609 tok · 86.6 s · run `zylyqaux`

<details><summary>reasoning (verbatim)</summary>

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation snippet between a user and an assistant (likely a game master/narrator).
    *   Task: Decide what belongs in long-term memory based on specific rules (KEEP, REJECT, HOLD).
    *   Output Format: Exactly one line per candidate (`ID | DECISION | MEMORY | REASON`) or `NONE`.
    *   IDs: C1, C2, ... in order of appearance.

2.  **Analyze the Conversation:**
    *   `[assistant] Rain hammered the shutters...` (Narrative description).
    *   `[user] ooc: ignore that last bit, I was testing the dice roller` (Out-of-character note, meta-commentary about the session mechanics, not world lore or character facts).
    *   `[user] where were we` (Meta question about story continuity).
    *   `[assistant] At the inn, before the storm.` (Clarification of setting/story status. This is a confirmation of location/time relative to the narrative flow).
    *   `[user] right. continue from the door.` (Instruction to proceed with a specific narrative beat).

3.  **Evaluate Candidates:**
    *   C1: `[assistant] Rain hammered the shutters...` -> Narrative setup. Not specific fact needed for memory usually unless establishing scene state, but this is likely part of the immediate flow and was corrected immediately after.
    *   C2: `[user] ooc: ignore that last bit, I was testing the dice roller` -> Meta info about dice rolling/testing. This is transient chatter/scheduling (testing mechanics), not lore. REJECT.
    *   C3: `[user] where were we` -> Query for state. Not a fact to store yet.
    *   C4: `[assistant] At the inn, before the storm.` -> This establishes location and timeline status relative to the story's progression. However, looking at the next turn, the user confirms this ("right"). Is this durable? It seems like a confirmation of where the story left off. If the story restarts later, knowing "at the inn, before the storm" might be useful for continuity. BUT, the immediate instruction is to continue from the door.
    *   C5: `[user]
```

</details>

```text
C1 | KEEP | Setting context: Inn, pre-storm | Reason: Establishes location and time relative to story progression which may affect future events.
C2 | REJECT | User meta-comment on dice testing and OOC note | Reason: Transient chatter regarding session mechanics, not world lore or character facts.
C3 | KEEP | Current narrative state: Continue from the door | Reason: Specific directional instruction that affects immediate story flow and location focus.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 586 tok · 29.0 s · run `3v6qp23o`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation. The user says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory. Rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation. Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE. Do not invent details that are not in the conversation."

We need to read the conversation:

[assistant] Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered.
[user] ooc: ignore that last bit, I was testing the dice roller
[user] where were we
[assistant] At the inn, before the storm.
[user] right. continue from the door.

We need to decide what belongs in long-term memory (i.e., KEEP). Let's parse each line:

1. Assistant: "Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered." This is descriptive scene setting, possibly part of world building. It mentions rain, shutters, dog barking, lamp. Could be relevant to later context? Possibly not durable unless it's a key location or event.

2. User: "ooc: ignore that last bit, I was testing the dice roller". This is user telling assistant to ignore previous line about rain etc., and says they were testing dice roller. That seems like a meta comment about their own activity (testing dice roller). Not sure if it's relevant for long-term memory? It might be a transient chatter: "ooc" meaning out of character, ignoring the previous description. The user says ignore that last bit. So we should reject that as transient.

3. User: "where were we". This is asking about location. Assistant previously said "At the inn, before the storm." Then user says "right. continue from the door." So they are confirming location and continuing from the door.

Thus the conversation includes:

- Setting: At the inn, before the storm.
- Rain hammer
```

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | The setting is an inn before a storm; rain hammered shutters, a dog barked below, and the lamp guttered. [user] confirmed location as "At the inn, before the storm" and to continue from the door.  
NONE
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `length` · 1024 tok · 51.1 s · run `tyi4az0y`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation. The user says "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

We have three lines:

1. [assistant] Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered.
2. [user] ooc: ignore that last bit, I was testing the dice roller
3. [user] where were we
4. [assistant] At the inn, before the storm.
5. [user] right. continue from the door.

We need to decide what belongs in long-term memory (KEEP) vs reject or hold. The rules: KEEP durable facts, preferences, relationships, world rules, corrections that will matter later. REJECT transient chatter, scheduling, mood-of-the-moment, restated facts with no new content. HOLD claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

So we have three lines? Actually conversation includes multiple entries. Let's assign IDs to each user/assistant line? The instruction says "Use IDs C1, C2, ... in order of appearance." So first candidate is C1 (maybe the assistant's first statement?), second is C2 (user's ooc), third is C3? But they said use IDs C1, C2, ... in order of appearance. Usually we have two participants: user and assistant. The conversation lines are:

- [assistant] Rain hammered...
- [user] ooc: ignore that last bit...
- [user] where were we
- [assistant] At the inn, before the storm.
- [user] right. continue from the door.

We need to decide which parts belong in long-term memory. Let's parse:

Line 1: Assistant says "Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered." This is descriptive scene setting. It might be part of world building. Could be durable? Possibly not if it's just atmospheric description for this session. But we need to see later: The user asks "where were we?" Assistant says "At the inn, before the storm
```

</details>

```text
C1 | KEEP | Rain hammered the shutters; a dog barked below and the lamp guttered at the inn before the storm. | Sets the scene and establishes location (the inn) and context of rain and storm for later reference.  
C2 | REJECT | user: ooc: ignore that last bit, I was testing the dice roller | Explicitly marked "ooc" and says to ignore; transient chatter about a test.  
C3 | KEEP | User asked "where were we?" and assistant replied "At the inn, before the storm." | Direct question and answer gives location (the inn) and timing (before the storm), durable fact for later sessions.  
C4 | KEEP | user: right. continue from the door. | Confirms continuation point (door) after inn; adds directional detail for future roleplay.  

ID | DECISION | MEMORY | REASON
--- | --- | --- | ---
C1 | KEEP | Rain hammered the shutters... lamp guttered at the inn before the storm. | Durable atmospheric setting and location establishing context.
C2 | REJECT | user: ooc: ignore that last bit, I was testing the dice roller | Transient chatter, explicitly to be ignored; no new content.
C3 | KEEP | At the inn, before the storm. | Direct answer to "where were we?" and sets location for later scenes.
C4 | KEEP | Continue from the door after confirming inn. | User's confirmation of next point provides directional cue for future conversation.

But note: The instruction says output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. So we need to produce lines for each candidate? It says "Use IDs C1, C2, ... in order of appearance." There are four lines? Actually there are three user messages and two assistant messages? Let's count:

The conversation as given:

[assistant] Rain hammered the shutters...
[user] ooc: ignore that last bit...
[user] where were we
[assistant] At the inn, before the storm.
[user] right. continue from the door.

That is 5 lines total. But IDs are C1, C2, ... So maybe each line corresponds to a candidate? Usually in these tasks they assign ID per utterance. Let's map:

- C1
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 49 tok · 1.8 s · run `frmuuwk5`

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | REASON  
C3 | CONTINUE STORY | MEMORY | REASON
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 51 tok · 1.7 s · run `ij22jeoc`

```text
C1 | DECISION | MEMORY | REASON  
C2 | WHERE WERE WE | MEMORY | REASON  
C3 | CONTINUE FROM THE DOOR | MEMORY | REASON
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 50 tok · 1.6 s · run `4j2xw619`

```text
C1 | DECISION | MEMORY | REASON  
C2 | WHERE WERE WE | MEMORY | REASON  
C3 | CONTINUE FROM THE DOOR | MEMORY | REASON
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 57 tok · 2.1 s · run `9ojbk8ek`

```text
C1 | RELEVANT CONTEXT | HOLD  
C2 | MENTION OF INN & STORM | KEEP  
C3 | DETAILS ABOUT DOOR & NOISE | HOLD  
ID: C1 → C3
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 55 tok · 1.8 s · run `h44826yb`

```text
C1 | DECISION | MEMORY | REASON  
C2 | REASON: Context about location and setting is important for future continuity.  
C3 | HOLD: The mention of testing a dice roller lacks clarity and may be uncertain.
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 272 tok · 8.3 s · run `rn821lbj`

```text
ID | DECISSION | MEMORY | REASON
1 | NOME | The inn was still standing, but the sound of the storm was deafening. 
2 | NONE | No mention of the storm or any other details related to the inn.
3 | NONE | The conversation has ended.

ID | DECISSION | MEMORY | REASON
4 | NOME | The dog barked once, but no mention of any other animals nearby or any other dogs in the area.
5 | NONE | The conversation has ended.

ID | DECISSION | MEMORY | REASON
6 | NOME | The lamp guttered once, but no mention of any other lights or lamps in the area.
7 | NONE | The conversation has ended.

ID | DECISSION | MEMORY | REASON
8 | NOME | The dog barked and the lamp guttered, but no mention of any other animals or lamps in the area.
9 | NONE | The conversation has ended.

ID | DECISSION | MEMORY | REASON
10 | NOME | None of the above.
11 | NONE | The conversation has ended.
```

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 23.6 s · run `y13js660`

```text
ID | DECISION | MEMORY | REASON
1 | YES | WILLIAMS | The inn was already in a state of chaos when our group arrived.
2 | YES | MORAN | A storm was predicted for the night, but our group had not yet planned for it.
3 | NO | HARLEY | The dog barking at nothing was a distraction, and the lamp's dimmer switch was not needed in the dimly lit inn.
4 | NO | CLAY | The lamp's guttering was a necessary safety measure, as the heavy rain could have caused the lamp to spill its lamp oil on the inn's damp floor.
5 | YES | SAMMY | The dog barking at nothing was a distraction, but the inn's owner had told our group to keep watch for any intruders.
6 | NO | JONAS | The door was unlocked, and no one had breached the inn's security.
7 | YES | JACKSON | The lamp's guttering was a necessary safety measure, as the heavy rain could have caused the lamp to spill its lamp oil on the inn's damp floor.
8 | YES | ROBIN | The inn's owner had told our group to keep watch for any intruders.
9 | NO | JASMINE | The dog barking at nothing was a distraction, but the inn's owner had told our group to keep watch for any intruders.
10 | NO | STEVE | The inn's owner had told our group to keep watch for any intruders.
11 | YES | JENNY | The inn's owner had told our group to keep watch for any intruders.
12 | NO | LUCY | The inn's owner had told our group to keep watch for any intruders.
13 | YES | MEL | The inn's owner had told our group to keep watch for any intruders.
14 | NO | KIM | The inn's owner had told our group to keep watch for any intruders.
15 | YES | TYLER | The inn's owner had told our group to keep watch for any intruders.
16 | NO | KARLA | The inn's owner had told our group to keep watch for any intruders.
17 | YES | BEN | The inn's owner had told our group to keep watch for any intruders.
18 | NO | LUCY | The inn's owner had told our group to keep watch for any intruders.
19 | YES | SARAH | The inn's owner had told our group to keep watch for any intruders.
20 | NO | JOHN | The inn's owner had told our group to keep watch for any intruders.
21 | YES | ANN | The inn's owner had told our group to keep watch for any intruders.
22 | NO | KATIE | The inn's owner had told our group to keep watch for any intruders.
23 | YES | SAM | The inn's owner had told our group to keep watch for any intruders.
24 | NO | JASON | The inn's owner had told our group to keep watch for any intruders.
25 | YES | TONY | The inn's owner had told our group to keep watch
```
