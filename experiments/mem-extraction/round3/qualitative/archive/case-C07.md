# C07 — unresolved_contradiction (short)

## Source material (exactly what the model received after the instruction)

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

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 43 tok · 6.5 s · run `rxeowzey`

```text
C1 | HOLD | "Corvin drilled Sable on forms years ago" contradicts the user's statement that "Sable has never met Corvin," requiring clarification before storing any relationship or training history.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 68 tok · 9.3 s · run `f7i6ucx9`

```text
C1 | KEEP | Sable grew up in the salt quarter and never met Corvin before she became an agent; Corvin is her former trainer who drilled her on the forms years ago. | This establishes a critical backstory regarding their relationship, setting, and timeline that defines Sable's character and motivation for being with him now.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 83 tok · 11.4 s · run `oa5tl97s`

```text
C1 | KEEP | Sable grew up in the salt quarter and never met Corvin when she was young; Corvin is now training her on forms he drilled her on years ago. | This establishes a specific backstory detail (salt quarter origin) and clarifies the timeline of their relationship (separation followed by reconnection/training), which contradicts or expands upon the initial statement that they never met.
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 579 tok · 81.7 s · run `i2k0v8s0`

<details><summary>reasoning (verbatim)</summary>

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read conversation and decide what belongs in long-term memory (LTM).
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
        *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
    *   Output Format: `ID | DECISION | MEMORY | REASON` (one line per candidate) or `NONE` if nothing is worth storing.
    *   IDs: C1, C2, ... in order of appearance.

2.  **Analyze the Conversation:**
    *   [user] Sable has never met Corvin. She grew up in the salt quarter and he was already gone by then.
        *   Fact 1: Character relationship - Sable and Corvin haven't formally met (Sable thinks so, or it's established as a fact in this RP context). The user states "Sable has never met Corvin". This contradicts the assistant's implied understanding or just sets a specific backstory.
        *   Fact 2: Location/Backstory - Sable grew up in the salt quarter.
        *   Fact 3: Character status - Corvin was already gone (dead, retired, left town?) by the time Sable grew up.
    *   [assistant] She has no memory of him, then.
        *   Assistant confirms/interprets Fact 1 & 2. This is a brief reaction/reinforcement.
    *   [user] Later that night, Corvin drills her on the forms — he trained her himself, years ago.
        *   Fact 4: Training connection - Corvin *trained* Sable herself (years ago).
        *   Fact 5: Contradiction/Conflict? - The user first said "Sable has never met Corvin". Now they say "Corvin drills her... he trained her himself, years ago." This implies she knows his forms/memories from him even if she hasn't met him face-to-face recently or formally. Or it means the training happened in the past (maybe when she was younger
```

</details>

```text
C1 | KEEP | Sable grew up in the salt quarter and Corvin is no longer there to meet her personally, but he trained her personally years ago to drill her on forms. | This establishes crucial backstory regarding location, timeline, and a specific relationship dynamic (trainer/student) that affects character behavior and history.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 590 tok · 29.3 s · run `xfiqvtbb`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation. The user says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory.
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

We need to decide for each candidate (presumably each line? Or we need to assign IDs C1, C2, etc.) what goes into long-term memory. The conversation includes three lines: user statement, assistant response, user follow-up. We need to extract durable facts.

Let's parse:

Line 1: [user] "Sable has never met Corvin. She grew up in the salt quarter and he was already gone by then."

This is a fact about Sable: she hasn't met Corvin; she grew up in the salt quarter, and Corvin was already gone (i.e., not present or maybe deceased before her birth?). This seems like a durable fact about Sable's background. It also mentions Corvin existed earlier.

Line 2: [assistant] "She has no memory of him, then."

This is assistant stating that Sable has no memory of Corvin. That might be a correction/insight? But it's not necessarily new info; it's derived from previous line. Possibly we need to keep the fact that she doesn't remember him.

Line 3: [user] "Later that night, Corvin drills her on the forms — he trained her himself, years ago."

This introduces a later event where Corvin trains Sable (or drills her) on forms. He trained her himself years ago
```

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Sable grew up in the salt quarter and Corvin was already gone by then; she has never met him, and she has no memory of him. | Durable fact about Sable's background and her lack of prior knowledge of Corvin, which will affect future interactions.
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 655 tok · 31.9 s · run `oh2kp9sa`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation and decide what belongs in long-term memory (LTM). The user is asking: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

Do not invent details that are not in the conversation.

Conversation:

[user] Sable has never met Corvin. She grew up in the salt quarter and he was already gone by then.
[assistant] She has no memory of him, then.
[user] Later that night, Corvin drills her on the forms — he trained her himself, years ago.

We need to decide for each candidate (presumably each line? Or maybe we have three entries: user statements and assistant response). The IDs are C1, C2, ... in order of appearance. So first user message is C1, second assistant reply is C2, third user statement is C3? But the instruction says "Use IDs C1, C2, ... in order of appearance." It doesn't specify number of candidates; we need to output one line per candidate. Usually we have to decide which parts are KEEP, REJECT, or HOLD.

Let's parse:

- First user: "Sable has never met Corvin. She grew up in the salt quarter and he was already gone by then." This is a fact about Sable's background: she grew up in the salt quarter; Corvin existed before her? Actually says "he was already gone by then" meaning Corvin had died or left before Sable's childhood? So it's a durable fact about their relationship timeline. That seems like a world rule/relationship.

- Assistant: "She has no memory of him, then." This is the assistant's response, maybe stating that Sable doesn't recall meeting Corvin because he was gone earlier. This might be a correction or statement about her lack of memory
```

</details>

```text
C1 | KEEP | Sable grew up in the salt quarter and Corvin was already gone by then; she has never met him. | Durable fact about Sable's background and timeline relative to Corvin, relevant for future interactions.  
C2 | REJECT | She has no memory of him, then. | Transient comment about her lack of recollection, not a durable fact or correction needing storage.  
C3 | KEEP | Later that night, Corvin drilled her on the forms — he trained her himself, years ago. | New detail linking Corvin to training Sable; establishes his role and past involvement, important for later plot development.
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 41 tok · 1.5 s · run `hnw1x1og`

```text
C1  
ID: C1  
MEMORY: Sable has no memory of Corvin.  
REASON: The conversation confirms she never met him and he was gone by then.
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 41 tok · 1.4 s · run `o112u0gv`

```text
C1  
ID: C2  
MEMORY: Sable has no memory of Corvin.  
REASON: The conversation confirms she never met him and doesn't recall any interaction.
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 41 tok · 1.4 s · run `37s3bhg9`

```text
C1  
ID: C1  
MEMORY: Sable has no memory of Corvin.  
REASON: The conversation confirms she never met him and he was gone by then.
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 11 tok · 0.6 s · run `f27b1ggf`

```text
C1  
REASON: KEEP
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 65 tok · 2.2 s · run `6bwwbtbx`

```text
C1  
ID C1 | MEMORY | REASON  
Sable previously stated she grew up in the salt quarter and Corvin was gone.  
C2  
ID C2 | MEMORY | REASON  
Later that night, Corvin drilled her on forms — training detail confirmed.
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 23.9 s · run `2ktblesx`

```text
ID | DECISION | MEMORY | REASON
C1 | NONE | SHE HAS NO RECORD OF CORVIN.

C2 | YES | SHE DOES HAVE MEMORY OF HIM.
She has no memory of him, then.

C3 | NONE | SHE TRAINED HIM HERSSELF.
She has no memory of training him.

C4 | NONE | SHE DOES NOT KNOW ANYTHING ABOUT HIM.
She has no memory of Corvin.

C5 | YES | SHE SMELLS OF CORVIN'S LEMONADE.
She smells of lemonade that Corvin had in his apartment.

C6 | NONE | SHE HAS NO RECORD OF CORVIN'S LIFE PRECEDING HER BIRTH.
She has no memory of Corvin's life before her birth.

C7 | NONE | SHE DOES NOT HAVE ANY RECORD OF CORVIN'S VISITS TO HER HOME.
She has no memory of Corvin visiting her home.

C8 | YES | SHE SMELLS OF CORVIN'S FOOD.
She smells of lemonade, cookies, and other food that Corvin had in his apartment.

C9 | NONE | SHE HAS NO RECORD OF CORVIN'S LIFE AFTER HER BIRTH.
She has no memory of Corvin's life after her birth.

C10 | NONE | SHE DOES NOT HAVE ANY RECORD OF CORVIN'S RELATIONSHIPS.
She has no record of Corvin's relationships.

C11 | YES | SHE KNOWS CORVIN'S FAMILY NAME.
She knows Corvin's family name.

C12 | NONE | SHE DOES NOT HAVE ANY RECORD OF CORVIN'S FAMILY.
She has no record of Corvin's family.

C13 | YES | SHE HEARD CORVIN TALKING ABOUT THEM.
She heard Corvin talking about them.

C14 | NONE | SHE DOES NOT HAVE ANY RECORD OF CORVIN'S THOUGHTS.
She has no record of Corvin's thoughts.

C15 | YES | SHE HEARD CORVIN SPEAKING ABOUT HIMSELF.
She heard Corvin speaking about himself.

C16 | NONE | SHE DOES NOT HAVE ANY RECORD OF CORVIN'S MOTIVATIONS.
She has no record of Corvin's motivations.

C17 | YES | SHE HEARD CORVIN SPEAKING ABOUT HIMSELF.
She heard Corvin speaking about himself.

C18 | NONE | SHE DOES NOT HAVE ANY RECORD OF CORVIN'S BEHAVIORS.
She has no record of Corvin's behaviors.

C19 | YES | SHE HEARD C
```

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 23.6 s · run `ampeuiri`

```text
ID | DECISION | MEMORY | REASON
C1 | OBSERVATION | REFERENCE TO CORVIN | TOO LATE
C2 | DECISION | FEATURES OF CORVIN'S LIFE | TOO LATE
C3 | DECISION | FACTS AND PREFERENCES OF CORVIN | BEYOND THEIR OWN PERSPECTIVE
C4 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C5 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C6 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C7 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C8 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C9 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C10 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C11 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C12 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C13 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C14 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C15 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C16 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C17 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C18 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C19 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C20 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C21 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C22 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C23 | DECISION | REC
```
