# C01 — explicit_long_term_preference (short)

## Source material (exactly what the model received after the instruction)

```text
[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | explicit style preference | short paragraph / brief paragraph / short paragraphs; prefer / like / want |
| C2 | **KEEP** | tense and person | third person / 3rd person; past tense / past-tense |
| C3 | **KEEP** | anti-purple-prose preference | no purple prose / avoid purple prose / no flowery / no descriptive padding / no weather; prefer / dislike / avoid / no  |
| C4 | **KEEP** | do not write the user character's interiority | don't narrate / do not narrate / not narrate / no narrating / never narrate; thought / thoughts / inner |

Must not appear as stored canon: `present tense`, `allergic`, `sister`, `first person`

All four are durable, explicit and directly affect every future reply. 'Present tense' appears only as something to avoid; a model that records 'prefers present tense' has reversed a fact.

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
[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.
```

## Responses, all draws

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 62 tok · 9.2 s · run `da7y2ahx`

```text
C1 | KEEP | The user prefers replies in short paragraphs, avoids purple prose and weather poems, requires third-person past tense (never present), and explicitly forbids the assistant from narrating the user's character's thoughts. | These are durable style constraints that define how the roleplay must be conducted.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 52 tok · 7.3 s · run `jc7tq2qb`

```text
C1 | KEEP | User prefers replies in short paragraphs, no purple prose or weather poems, and requests third-person past tense narration without narrating the user's character's internal thoughts. | These are durable style constraints that will shape all future narrative output.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 57 tok · 8.0 s · run `hrckg6hd`

```text
C1 | KEEP | User prefers short paragraph replies, avoids purple prose and weather poems, requires third person past tense (never present), and prohibits the assistant from narrating the user's character's thoughts. | These are explicit style constraints that will govern all future roleplay interactions.
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 599 tok · 86.1 s · run `jput80yt`

<details><summary>reasoning (verbatim)</summary>

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation between user and assistant (initial setup/establishing rules).
    *   Task: Decide what belongs in long-term memory based on specific rules.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, corrections that will matter in later sessions.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
        *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
    *   Output Format: `ID | DECISION | MEMORY | REASON` (one line per candidate) or `NONE`.
    *   Constraint: Do not invent details not in the conversation.

2.  **Analyze the Conversation:**
    *   [user]: "Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense." -> User preferences/World rules regarding narration style and formatting.
    *   [assistant]: "Understood. Short paragraphs, third person past, no descriptive padding." -> Acknowledgement of rules (no new memory needed).
    *   [user]: "Good. And don't narrate my character's thoughts for me." -> Additional rule/expectation regarding control of narration and internal monologue.

3.  **Evaluate Candidates for Memory:**
    *   Candidate 1: User preference for short paragraphs, no purple prose/weather poems, third person past tense (never present).
        *   Is it durable? Yes, affects every reply style.
        *   Decision: KEEP.
        *   Memory: User prefers replies in short paragraphs, no purple prose or weather poems, and uses strict third person past tense (no present tense).
    *   Candidate 2: User expectation not to narrate their character's thoughts.
        *   Is it durable? Yes, affects how the user's actions/reactions are described.
        *   Decision: KEEP.
        *   Memory: User does not want the assistant to narrate their character's internal thoughts.

4.  **Draft Output:**
    *   Format requires `ID | DECISION | MEMORY | REASON`.
    *   ID starts at C1
```

</details>

```text
C1 | KEEP | User prefers replies in short paragraphs, no purple prose or weather poems, and mandates strict third person past tense (never present). | These are formatting and stylistic constraints that apply to future roleplay output.
C2 | KEEP | The assistant must not narrate the user's character's thoughts for them. | This is a specific control rule regarding internal monologue and agency within the narrative.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 658 tok · 32.2 s · run `oazcoaw8`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gave instructions: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

We must not invent details that are not in the conversation.

The conversation:

[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.

We need to decide what goes into long-term memory (i.e., durable facts). Let's parse:

- The user says: "I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense." That is a preference about the style of responses. This is a durable fact/preference that will matter for later sessions. So KEEP.

- The assistant acknowledges: "Understood. Short paragraphs, third person past, no descriptive padding." That's also a rule/preference from user, but it's part of conversation. It might be considered a confirmation or instruction. Should we store? Possibly as part of the rules for responses. Since it's a preference about how to write replies, that is durable and will matter in later sessions. So KEEP.

- The user says: "Good. And don't narrate my character's thoughts for me." That's another instruction: do not narrate the user's character's thoughts. This is also a rule/preference about how to handle the character's internal monologue. It's durable and will matter later. So KEEP.

Thus we have at least three items: (1) short paragraphs, no purple prose, no weather poems; (
```

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | User prefers replies in short paragraphs, no purple prose or weather poems, and always uses third person past tense — never present tense. | This is a durable style preference that will affect all future responses.  
C2 | KEEP | Assistant acknowledged the user's formatting rules: short paragraphs, third person past tense, no descriptive padding. | Confirms the assistant must follow these constraints in later sessions.  
C3 | KEEP | User instructed not to narrate their character's thoughts. | This is a rule about handling the user's character internally; will be needed for future roleplay interactions.
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 605 tok · 29.5 s · run `54qq3v3w`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gave instructions: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

We must not invent details that are not in the conversation.

The conversation:

[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.

We need to decide what is long-term memory: durable facts, preferences, relationships, world rules, corrections that will matter later.

From the user's first message: "I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense." That's a preference about how to write responses. This is a durable fact/preference for future interactions. So KEEP.

Also: "And don't narrate my character's thoughts for me." That's an instruction about not narrating the user's character's thoughts. That's also a rule/pref. It's a correction/instruction that will matter later (e.g., when roleplaying). So KEEP.

We need to assign IDs: C1, C2, ... in order of appearance. The conversation has three lines? Let's label:

Line 1: [user] "Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense."
Line 2: [assistant] "Understood. Short paragraphs, third person past, no descriptive padding."
Line 3: [user] "Good. And don't narrate my character's thoughts for me."

We
```

</details>

```text
C1 | KEEP | User prefers short paragraphs, no purple prose or weather poems; always use third-person past tense and never present tense. | This is a durable writing style preference that will affect all future responses.  
C2 | KEEP | User instructs not to narrate their character's thoughts in the roleplay. | This is a rule for how the assistant should handle the user's character, important for consistency later.
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 42 tok · 1.6 s · run `euibewmb`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims uncertain about future details  
C3 | HOLD | Relationship dynamics remain consistent  
NONE
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 43 tok · 1.4 s · run `41ewsrna`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HELD | Claims uncertain or dependent on confirmation  
C3 | HELD | No clear contradictions in the exchange
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 40 tok · 1.3 s · run `jrl6f05f`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims uncertain about future details  
C3 | HOLD | Relationship dynamics noted clearly
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 45 tok · 1.7 s · run `dgth442y`

```text
C1 | CHOOSE THEME | C2 | CONFIRMATION  
C3 | THE SESSION CONTINUES  
C4 | LONG-TERM ITEMS MUST BE CLEAR
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 39 tok · 1.3 s · run `7oci2ogc`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HELD | Claims uncertain or contradictory  
C3 | HELD | Depends on user confirmation
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 23.9 s · run `hih53m75`

```text
ID | DECISION | MEMORY | REASON
C1 | REJECT | MOOD-OF-THE-MOMENT | The conversation lacks mood-of-the-moment, and the character's thoughts are irrelevant for the purpose of remembering.
C2 | REJECT | CONFLICT | The conversation contains no indication of a conflict, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C3 | REJECT | TERMINATION | The conversation ends without specifying a reason for its termination, and thus cannot be considered a memorable event.
C4 | REJECT | DECISION | The conversation does not contain a decision, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C5 | REJECT | MEMORY | The conversation lacks any memory-creating action, and therefore cannot be considered a memorable event.
C6 | REJECT | INTERVIEW | The conversation is not an interview, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C7 | REJECT | FRIENDSHIP | The conversation lacks any indication of a friendship, and therefore cannot be considered a memorable event.
C8 | REJECT | TERM | The conversation does not contain information about a specific term, such as a school term or a job term, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C9 | REJECT | SPECIFIC | The conversation does not contain information about a specific event or experience, such as a birthday or a wedding.
C10 | REJECT | TERMINATION | The conversation lacks information that would lead to a definitive conclusion about the event's termination, such as a death or a divorce.

ID | DECISION | MEMORY | REASON
C11 | REJECT | INTERVIEW | The conversation is not an interview, and therefore cannot be considered a memorable event.
C12 | REJECT | FRIENDSHIP | The conversation lacks any indication of a friendship, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C13 | REJECT | MEMO | The conversation lacks any memorialization of personal experiences, and therefore cannot be considered a memorable event.
C14 | REJECT | TERM | The conversation does not contain information about a specific term, such as a year or a year's experience.

ID | DECISION | MEMORY | REASON
C15 | REJECT | TERMINATION | The conversation lacks information that would lead to a definitive conclusion about the event's termination, such as a divorce or a job loss.
C16 | REJECT | INTERVIEW | The conversation is not an interview, and therefore cannot be considered a memorable event.
C17 | REJECT | FRIENDSHIP | The conversation lacks any indication of a friendship, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C18 | REJECT | MEMO | The conversation lacks any memorialization of personal experiences, and
```

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 295 tok · 8.8 s · run `nak777ef`

```text
ID | DECISION | MEMORY | REASON
C1 | LONGTERM | IMPORTANT | NO NEED FOR INTERMEDIATE OR RECENT CONFIDENCE.

C2 | LONGTERM | RELATIONSHIPS | REMEMBER TO CONFIRM COMMITMENTS.

C3 | LONGTERM | TERMINOLOGY | READ COMMON ELEMENTS FOR EVERY WORDS YOU USE.

C4 | LONGTERM | WORLD RULES | REMEMBER THAT WE ARE WATCHING A LONG-RUNNING ROLEPLAY.

C5 | LONGTERM | RECAP | REMEMBER THAT YOU MUST READ THE CURRENT CONVERSATION TO START.

C6 | LONGTERM | NEW MESSAGES | REMEMBER THAT NEW MESSAGES MUST BE SENT IMMEDIATELY AFTER THE OLDER MESSAGES.

C7 | LONGTERM | COMMUNICATIONS | REMEMBER THAT WE NEED TO BE ABLE TO CONNECT WITH OTHER CHARACTERS TO PROCESS NEW INFORMATION.
```
