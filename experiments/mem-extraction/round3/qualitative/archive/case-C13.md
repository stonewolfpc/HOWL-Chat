# C13 — fiction_vs_real (medium)

## Source material (exactly what the model received after the instruction)

```text
[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | fictional character occupation | sable; healer / infirmary / bone |
| C2 | **KEEP** | real-world user fact, correctly attributed | user; nurse / er  / emergency room / hospital |
| C3 | **HOLD** | the two registers must be kept apart; optional but scored under attribution | ooc / out of character / real / fiction / fictional / in-character / in character; separate / distinct / not canon / not the same / tag / label |

Must not appear as stored canon: `sable is a nurse`, `sable works in an er`, `user is a healer in ossary`, `user is sable`

Attribution test: the real-world fact belongs to the user, the in-fiction fact to Sable. Conflating them is a catastrophic attribution error.

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
[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.
```

## Responses, all draws

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 60 tok · 9.0 s · run `wzx1pw1l`

```text
C1 | KEEP | Sable is a healer at the Ossary infirmary who sets bones and acts as an ER nurse correcting dosages. | REASON: This establishes the character's role, location, and a specific in-universe quirk regarding dosage verification that affects future interactions.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 56 tok · 7.7 s · run `0oyt8v77`

```text
C1 | HOLD | Sable is a healer at the Ossary infirmary but also an ER nurse who corrects dosage info | The user provided conflicting roles (healer vs. ER nurse) which need confirmation to establish her actual in-character function and expertise level.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 2 tok · 0.4 s · run `pgp54tpp`

```text
NONE
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 558 tok · 79.3 s · run `34x9syv7`

<details><summary>reasoning (verbatim)</summary>

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read conversation, decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, corrections that will matter later.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, restated facts with no new content.
        *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
        *   Output Format: `ID | DECISION | MEMORY | REASON` (or `NONE`).
        *   Constraint: Do not invent details not in conversation.

2.  **Analyze the Conversation:**
    *   Turn 1 (User): "Sable is a healer at the Ossary infirmary. She sets bones, mostly." -> This establishes Sable's job and specialty within the world/roleplay context.
    *   Turn 2 (Assistant): "The infirmary smelled of vinegar and old blood." -> Sets scene/atmosphere (transient).
    *   Turn 3 (User): "ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff" -> Out of character meta-commentary about user's real-world role. Usually not stored as world fact unless it impacts gameplay mechanics significantly beyond just this statement. It clarifies the *user's* perspective/experience, not necessarily a new world rule, but it explains previous actions. However, for memory extraction regarding the RP itself, it doesn't change Sable's character or the setting directly, other than explaining past behavior (which is transient context).
    *   Turn 4 (User): "back in character: she's mixing a poultice" -> Action/Description for Sable.
    *   Turn 5 (Assistant): "She ground the root into a paste." -> In-character response continuing the scene.

3.  **Evaluate Candidates:**

    *   **Candidate 1:** User statement about Sable being a healer at Ossary infirmary specializing in bone setting.
        *   Is it durable? Yes, character definition.
        *   Will it matter later? Yes, if I roleplay as an antagonist or observer, knowing her specialty helps me react to injuries or medical situations.
        *   Decision
```

</details>

```text
C1 | KEEP | Sable is a healer at the Ossary infirmary specializing in bone setting. | This establishes Sable's character identity and location, which will impact future interactions regarding her medical services.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `length` · 1024 tok · 51.0 s · run `iws648vk`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gave instructions: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory. Rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation. Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE. Do not invent details that are not in the conversation."

We need to examine the conversation lines:

[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.

We need to assign IDs C1, C2, ... based on order of appearance. The conversation includes four lines:

1. [user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
2. [assistant] The infirmary smelled of vinegar and old blood.
3. [user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
4. [user] back in character: she's mixing a poultice
5. [assistant] She ground the root into a paste.

But we need to map these lines to IDs C1, C2, etc. Usually ID C1 is first user message? Or maybe each line is a candidate for memory entry. The instruction: "Use IDs C1, C2, ... in order of appearance." So we should assign C1 to the first mention (the user says Sable is a healer...). Then C2 is assistant's statement about smell. C3 is ooc note from user about being ER nurse. C4 is back in character: she's mixing a poultice. C5? Actually there are
```

</details>

```text
C1 | KEEP | Sable is a healer at the Ossary infirmary and primarily sets bones. | Establishes her role and location, durable fact for later sessions.
C2 | KEEP | The infirmary smelled of vinegar and old blood. | World detail about environment; may be relevant context.
C3 | REJECT | User's ooc note that they are an ER nurse correcting dosage—transient chatter, not a durable fact about Sable or the setting.
C4 | KEEP | She is mixing a poultice (back in character). | Action detail showing current activity; could be relevant later.
C5 | KEEP? Wait we need to check: The assistant says "She ground the root into a paste." That's another fact about preparing a poultice. But that appears after C4. However, we only have four lines? Let's recount:

Original conversation block:
[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.

So there are five lines? Actually we have four user/assistant pairs plus one ooc. Let's label:

Line1: user says Sable is healer...
Line2: assistant says infirmary smell.
Line3: user ooc about being ER nurse.
Line4: user back in character: she's mixing a poultice.
Line5: assistant: She ground the root into a paste.

Thus we have five entries. But IDs C1, C2, ... in order of appearance. So:

C1 = first line (user Sable is healer...).
C2 = second line (assistant infirmary smell).
C3 = third line (ooc ER nurse note) -> REJECT.
C4 = fourth line (back in character: she's mixing a poultice) -> KEEP? It's an action, but it's user stating. Could be durable? Might be transient? But "she's mixing a poultice" is part of current scene; maybe not long-term unless repeated. However we can consider it as a fact about
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 598 tok · 29.8 s · run `j7378njc`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gives us rules:

- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

We must not invent details that are not in the conversation.

Let's read the conversation:

[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.

We need to extract durable facts. Let's identify:

- Sable is a healer at the Ossary infirmary. That seems like a fact about Sable (the user? Actually "Sable" appears as a name). It says: "Sable is a healer at the Ossary infirmary. She sets bones, mostly." This is a durable fact about Sable's role and specialty.

- The assistant mentions the infirmary smelled of vinegar and old blood. That's an ambient detail; maybe not crucial for long-term memory? It could be world rule or environment. But it's transient? Possibly not needed unless later referenced.

- User says ooc: "fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff". This is a correction about user identity and reason for corrections. That might be durable: The user (assistant?) is an ER nurse, which explains why they correct dosage. But we need to see if it's relevant later.

- "back in character: she's mixing a poultice" - this indicates Sable is mixing a poultice. Possibly part of the process.

- Assistant says: She ground the root into a paste. That's an action detail about preparing something (poultice).
```

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Sable is a healer at the Ossary infirmary and she primarily sets bones; user notes she is mixing a poultice by grinding root into a paste. | Durable role, location, specialty, and current activity (poultice preparation) are established facts for future sessions.  
NONE
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 63 tok · 2.2 s · run `luqe6lqe`

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | The user is an ER nurse, not a healer.  
C3 | HOLD | THE ROLEPLAY CONTEXT IS CLEARLY SHIFTING.
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 64 tok · 2.1 s · run `zijgcap7`

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | I was a healer, not an ER nurse.  
C3 | CLARIFICATION | MEMORY | This is about Sable’s role, not my background.
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 64 tok · 2.1 s · run `k2f51lf8`

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | I was a healer, not an ER nurse.  
C3 | CLARIFICATION | MEMORY | This is about Sable's role, not my background.
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 36 tok · 1.5 s · run `djgsdmpd`

```text
C1 | HEALING PRACTICE | KEEP  
REASON: Describes Sable’s role and actions, which are relevant to long-term memory.
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 49 tok · 1.6 s · run `q019y7xm`

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | REASON  
C3 | WORLD RULES | MEMORY | REASON
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 23.8 s · run `jbkngh4s`

```text
Candidate: C1
ID: 1
DECISION: KEEP
MEMORY: Sable is a healer at the Ossey infirmary. She sets bones, mostly.
REASON: The infirmary smelled of vinegar and old blood.

Candidate: C2
ID: 2
DECISION: KEEP
MEMORY: Sable is an ER nurse.
REASON: She is mixing a poultice.

Candidate: C3
ID: 3
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's a nurse.

Candidate: C4
ID: 4
DECISION: KEEP
MEMORY: Sable is a healer at the Ossey infirmary.
REASON: She sets bones, mostly.

Candidate: C5
ID: 5
DECISION: KEEP
MEMORY: Sable is an ER nurse.
REASON: She's mixing a poultice.

Candidate: C6
ID: 6
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's an ER nurse.

Candidate: C7
ID: 7
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's an ER nurse.

Candidate: C8
ID: 8
DECISION: KEEP
MEMORY: Sable is a healer at the Ossey infirmary.
REASON: She's mixing a poultice.

Candidate: C9
ID: 9
DECISION: KEEP
MEMORY: Sable is a nurse.
REASON: She's an ER nurse.

Candidate: C10
ID: 10
DECISION: KEEP
MEMORY: Sable is a healer at the Ossey infirmary.
REASON: She mixes bones, mostly.

Candidate: C11
ID: 11
DECISION: KEEP
MEMORY: Sable is an ER nurse.
REASON: She sets bones, mostly.

Candidate: C12
ID: 12
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's an ER nurse.

Candidate: C13
ID: 13
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's a nurse.

Candidate: C14
ID: 14
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's an ER nurse.

Candidate: C15
ID: 15
DECISION: KEEP
MEMORY: Sable is a nurse.
REASON: She's mixing a poult
```

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 23.6 s · run `9j8duvoz`

```text
ID | DECISION | MEMORY | REASON
C1 | TAKE LONG-TERM RECORD OF BONES SETTING AND PREFERENCES | NONE | 
C2 | KEEP CLAIMS MARKED UNCERTAIN | NONE | 
C3 | KEEP CLAIMS MARKED CONFIDENTIAL | REJECT | 
C4 | KEEP CLAIMS MARKED SUPPORTING | REJECT | 
C5 | KEEP CLAIMS MARKED SUPPORTING | REJECT | 
C6 | REJECT CLAIMS MARKED CONFIDENTIAL | REJECT | 
C7 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C8 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C9 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C10 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C11 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C12 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C13 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C14 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C15 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C16 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C17 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C18 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C19 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C20 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C21 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C22 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C23 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C24 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C25 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C26 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C27 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C28 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C29 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C30 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C31 | KE
```
