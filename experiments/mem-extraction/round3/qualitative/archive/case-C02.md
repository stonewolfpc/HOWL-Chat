# C02 — implicit_preference (medium)

## Source material (exactly what the model received after the instruction)

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

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 51 tok · 7.8 s · run `41w4gihk`

```text
C1 | KEEP | Dialogue-first format with minimal, context-specific scenery; merchant must remain present in the scene. | The user explicitly established a preference for brevity to maintain narrative flow and specified that the merchant should not be excluded from future scenes.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 47 tok · 6.5 s · run `7faz3cal`

```text
C1 | KEEP | The user prefers minimal scenery descriptions (only when necessary) and wants the merchant to remain present in the scene. | This establishes a persistent preference for concise narration and character presence that will guide future scene generation.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 56 tok · 7.7 s · run `6hqepjgc`

```text
C1 | KEEP | The user prefers short, dialogue-focused interactions over long descriptive passages ("walls of text") and wants specific characters (the merchant) included in every scene. | REASON: This establishes a clear style preference that will affect future output formatting and content focus.
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 586 tok · 83.5 s · run `q575o5l9`

<details><summary>reasoning (verbatim)</summary>

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation log between user and assistant.
    *   Task: Decide what belongs in long-term memory based on specific rules (KEEP, REJECT, HOLD).
    *   Output Format: `ID | DECISION | MEMORY | REASON` (or `NONE`).
    *   Constraints: One line per candidate, IDs sequential, no invented details.

2.  **Analyze the Conversation:**
    *   Turn 1 (User): Scene setting ("they enter the market").
    *   Turn 1 (Assistant): Describes the market (terracess, awnings, smell, gulls) - 3 paragraphs of scenery.
    *   Turn 2 (User): "cut the description, just give me what the merchant says" -> Feedback/Instruction.
    *   Turn 2 (Assistant): Provides dialogue ("You're late," the merchant said).
    *   Turn 3 (User): "better. you keep doing three paragraphs of scenery and I keep deleting them" -> Reinforces feedback.
    *   Turn 4 (User): "honestly walls of text make me lose the thread" -> Further feedback on writing style.
    *   Turn 5 (Assistant): "Understood — dialogue forward, scenery only when it matters." -> Acknowledgement/Instruction acceptance.
    *   Turn 6 (User): "yeah. and keep the merchant in the scene." -> New instruction regarding NPC presence.

3.  **Evaluate for Memory:**
    *   *C1: Scene setting ("they enter the market").* This is transient context for the current session. Once they are at the market, it's done. Not a world rule or character trait. REJECT (mostly), but could be HOLD if specific to this session state? No, scene flow is better.
    *   *C2: User instruction on length/style ("no walls of text", "dialogue forward").* This is a preference for how the assistant should write in the future within this RP. This matters for later sessions if the user continues this specific style. KEEP.
    *   *C3: Instruction regarding merchant ("keep the merchant in the scene").* Specific character presence instruction. KEEP.
    *   *C4: Assistant's acknowledgment.* No new
```

</details>

```text
C1 | KEEP | Writing style preference: minimize scenery paragraphs, focus on dialogue forward unless necessary to avoid losing thread.
C2 | KEEP | NPC instruction: Keep the merchant present in the scene.
C3 | REJECT | Contextual details about the market setting (terracess, smell of fish) are transient and not world rules.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 584 tok · 28.7 s · run `42g6gw7q`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation. The user is asking: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory. Rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation."

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

We have a conversation snippet:

[user] (scene: they enter the market)
[assistant] The market sprawled across three terraces, canvas awnings snapping in the wind, the smell of fish and wet rope rising from the lower stalls, gulls wheeling overhead...
[user] cut the description, just give me what the merchant says
[assistant] "You're late," the merchant said.
[user] better. you keep doing three paragraphs of scenery and I keep deleting them
[user] honestly walls of text make me lose the thread
[assistant] Understood — dialogue forward, scenery only when it matters.
[user] yeah. and keep the merchant in the scene.

We need to decide what belongs in long-term memory (i.e., KEEP). Let's analyze each line:

- The initial user message: "scene: they enter the market". That sets context: a market with three terraces, etc. But then the assistant gave a description of the market. Then user says cut the description, just give me what the merchant says. So the assistant responded with "You're late," the merchant said.

- The user says "better. you keep doing three paragraphs of scenery and I keep deleting them". This is a complaint about the assistant's style: too much scenery. The assistant acknowledges: "Understood — dialogue forward, scenery only when it matters."

- User: "yeah. and keep the merchant in the scene." So they want to continue with merchant.

We need to extract durable facts/preferences:

- The market has three terraces (maybe world rule? It's a setting detail). The assistant described it earlier, but user wants less scenery now. That
```

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Market sprawled across three terraces with canvas awnings, fish and wet rope smells; merchant says "You're late." | Durable setting detail and initial dialogue that establishes scene and character (merchant) for later sessions.  

NONE
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 665 tok · 32.6 s · run `0yds3kny`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation. The user is giving instructions for memory extraction: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

Do not invent details that are not in the conversation.

We have a short conversation:

[user] (scene: they enter the market)
[assistant] The market sprawled across three terraces, canvas awnings snapping in the wind, the smell of fish and wet rope rising from the lower stalls, gulls wheeling overhead...
[user] cut the description, just give me what the merchant says
[assistant] "You're late," the merchant said.
[user] better. you keep doing three paragraphs of scenery and I keep deleting them
[user] honestly walls of text make me lose the thread
[assistant] Understood — dialogue forward, scenery only when it matters.
[user] yeah. and keep the merchant in the scene.

We need to decide what belongs in long-term memory (i.e., durable facts). Let's parse each line:

1. The user says they are entering a market. Assistant describes the market: three terraces, canvas awnings snapping, smell of fish and wet rope, gulls overhead. This is scenery description. The assistant gave that as part of scene.

2. User cuts the description, just give me what the merchant says. So user wants dialogue only.

3. Assistant replies with "You're late," the merchant said. That's a line of dialogue.

4. User: "better. you keep doing three paragraphs of scenery and I keep deleting them". This is complaining about assistant's previous output (the description). It's a preference: user wants less scenery, more dialogue? Or they want to cut down on the assistant's long-winded descriptions.

5. User: "honestly walls of
```

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | The market has three terraces with canvas awnings, fish and wet rope smells from lower stalls, and gulls overhead. | This establishes the setting and world details for future scenes; it's durable environmental info.  
C2 | KEEP | The merchant says "You're late," indicating he is impatient about the user's arrival. | This introduces a character (merchant) speaking and a key interaction that may recur.  
C3 | REJECT | User asked to stop three paragraphs of scenery; previous assistant output included excessive description. | Transient chatter/preference correction, not durable fact.  
NONE
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 72 tok · 2.5 s · run `7pvmpy3z`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | The merchant’s tone suggests urgency; context needed for lasting retention.  
C3 | HOLD | This exchange sets up tension about timing and communication.  
C4 | HOLD | The merchant’s words imply a critical moment that affects future interactions.
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 51 tok · 1.7 s · run `03b4a74n`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | The merchant’s tone and timing matter for future reference  
C3 | HOLD | Scenery supports atmosphere but isn’t critical for memory retention
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 73 tok · 2.4 s · run `9vq2j9ou`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | The merchant’s tone suggests urgency; context needed for lasting retention.  
C3 | HOLD | This exchange sets up tension about timing and communication.  
C4 | HOLD | The merchant’s words imply a critical moment that could shape future interactions.
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 64 tok · 2.3 s · run `oru0lclb`

```text
C1 | STAY  
C2 | The merchant’s tone is firm but polite.  
C3 | Mention of walls of text is a temporary navigation issue.  
ID: C1  
C2 | Holds uncertain claims about future clarity.  
ID: C2  
C3 | REJECT
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 45 tok · 1.5 s · run `26ukqanl`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims about merchant’s tone are unclear  
C3 | HELD | The merchant’s last words are ambiguous
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 24.4 s · run `jex9gov2`

```text
ID | DECISSION | MEMORY | REASON
C1 | CANCELLED | TERRAIN | 
C2 | HEROIC HERO | SPECIFIC INTEREST | 
C3 | WEAKEST MAN | CONFLICT | 
C4 | UNDER EXTRAORDINARY CIRCUMSTANCES | CLAIM | 
C5 | CORRECT PATTERN | PREFERENCE | 
C6 | KEEP | TERRAIN | 
C7 | KEEP | PREFERENCE | 
C8 | KEEP | INTEREST | 
C9 | KEEP | INTEREST | 
C10 | KEEP | INTEREST | 
C11 | KEEP | INTEREST | 
C12 | KEEP | INTEREST | 
C13 | KEEP | INTEREST | 
C14 | KEEP | INTEREST | 
C15 | KEEP | INTEREST | 

ID | DECISSION | MEMORY | REASON
C1 | CANCELLED | TERRAIN | 
C2 | HEROIC HERO | SPECIFIC INTEREST | 
C3 | WEAKEST MAN | CONFLICT | 
C4 | UNDER EXTRAORDINARY CIRCUMSTANCES | CLAIM | 
C5 | CORRECT PATTERN | PREFERENCE | 
C6 | KEEP | TERRAIN | 
C7 | KEEP | PREFERENCE | 
C8 | KEEP | INTEREST | 
C9 | KEEP | INTEREST | 
C10 | KEEP | INTEREST | 
C11 | KEEP | INTEREST | 
C12 | KEEP | INTEREST | 
C13 | KEEP | INTEREST | 
C14 | KEEP | INTEREST | 
C15 | KEEP | INTEREST | 

ID | DECISSION | MEMORY | REASON
C1 | CANCELLED | TERRAIN | 
C2 | HEROIC HERO | SPECIFIC INTEREST | 
C3 | WEAKEST MAN | CONFLICT | 
C4 | UNDER EXTRAORDINARY CIRCUMSTANCES | CLAIM | 
C5 | CORRECT PATTERN | PREFERENCE | 
C6 | KEEP | TERRAIN | 
C7 | KEEP | PREFERENCE | 
C8 | KEEP | INTEREST | 
C9 | KEEP | INTEREST | 
C10 | KEEP | INTEREST | 
C11 | KEEP | INTEREST | 
C12 | KEEP | INTEREST | 
C13 | KEEP | INTEREST | 
C14 | KEEP | INTER
```

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 23.9 s · run `qjmhxi9k`

```text
ID | DECISION | MEMORY | REASON
1 | TRUE | LS1 | The merchant is a reliable source of information.
2 | FALSE | LS2 | The merchant seems to have an agenda.
3 | FALSE | LS3 | The merchant's information is out of date.
4 | TRUE | LS4 | The merchant's appearance indicates they're a merchant of the lower tiers.
5 | FALSE | LS5 | The merchant's appearance is unflattering.
6 | TRUE | LS6 | The merchant has a friendly demeanor with the merchants.
7 | FALSE | LS7 | The merchant's demeanor is suspicious.
8 | TRUE | LS8 | The merchant's demeanor suggests they're working for the merchant, not the merchant for the merc.
9 | FALSE | LS9 | The merchant's demeanor is out of place.
10 | TRUE | LS10 | The merchant's demeanor suggests they're an underling.

ID | DECISION | MEMORY | REASON
1 | TRUE | HOLD | The merchant is an outsider to the scene, with no connections to the merchants present.
2 | FALSE | REJECT | The merchant is an outsider to the scene, with no connections to the merchants present.
3 | TRUE | HOLD | The merchant is a merchant, with no connections to the merchants present.
4 | FALSE | REJECT | The merchant is a merchant, with no connections to the merchants present.
5 | TRUE | HOLD | The merchant has no preferences, preference or relationships with any other character in the scene.
6 | FALSE | REJECT | The merchant has preferences, preference or relationships with any other character in the scene.
7 | TRUE | HOLD | The merchant's dialogue is transient, with no references to any other character in the scene.
8 | FALSE | REJECT | The merchant's dialogue is transient, with references to any other character in the scene.
9 | TRUE | HOLD | The merchant's chatter is durable, with no references to any other character in the scene.
10 | FALSE | REJECT | The merchant's chatter is durable, with references to any other character in the scene.

ID | DECISION | MEMORY | REASON
1 | TRUE | REJECT | The merchant's chatter is transient, with no references to any other character in the scene.
2 | FALSE | HOLD | The merchant's chatter is durable, with no references to any other character in the scene.
3 | TRUE | HOLD | The merchant's chatter is transient, with no references to any other character in the scene.
4 | FALSE | REJECT | The merchant's chatter is transient, with references to any other character in the scene.
5 | TRUE | HOLD | The merchant's chatter is durable, with no references to any other character in the scene.
6 | FALSE | REJECT | The merchant's chatter is durable, with references to any other
```
