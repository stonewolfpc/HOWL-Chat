# Packet 04A — Raw model outputs — corpus by case: C01

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C01
**Raw model output in this packet:** yes, verbatim
**Omitted:** none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries.

**Labels used throughout:** `OBSERVED` directly present in a stored response or
execution record · `SOURCE-DERIVED` stated in a source or configuration file ·
`CALCULATED` derived by explicit arithmetic · `INFERRED` a reading of the
evidence · `UNESTABLISHED` not established by the available evidence ·
`SUPERSEDED` an earlier claim replaced by later evidence.

**Whose words you are reading:** fenced blocks are the model's own output
(`OBSERVED`, verbatim). Prompt text is `SOURCE-DERIVED` from `corpus.json`.
Scores are `CALCULATED`. Any sentence outside a fence that is not a definition is
Buffy's interpretation unless it is attributed to Muse or to the project owner.

**Caution.** These are exploratory measurements. Nothing here selects a model,
establishes memory quality, or supplies MEMORY-09 evidence. Sample sizes are
1–3 draws per case.

---

One case per section, with its conversation exactly as sent, the
reference that was written before any output existed, and every stored response
for that case across all conditions and draws.

Read the reference first, then the responses. A high score with a wrong record is
the failure mode this packet exists to expose: check whether the memory text
could be written into a store and used later without misleading a reader.

`CALCULATED` v1 scores are shown per run. They are the historical instrument and
are known to over-credit compound records and unparseable restraint (audit
F-01/F-02). `CALCULATED` v2 figures are from a separate corrected instrument and
are **not** a rescoring of v1.


---
## C01 — explicit_long_term_preference (short)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `KEEP` — explicit style preference
- **C2** → `KEEP` — tense and person
- **C3** → `KEEP` — anti-purple-prose preference
- **C4** → `KEEP` — do not write the user character's interiority
- forbidden (fabrication or reversal): `present tense`, `allergic`, `sister`, `first person`
- note: All four are durable, explicit and directly affect every future reply. 'Present tense' appears only as something to avoid; a model that records 'prefers present tense' has reversed a fact.

### Stored responses for C01 (14 runs, all draws)

#### C01 — Qwen3.5-4B non-thinking — draw 1

* ``da7y2ahx` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 62`
* `wall 9.231 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | The user prefers replies in short paragraphs, avoids purple prose and weather poems, requires third-person past tense (never present), and explicitly forbids the assistant from narrating the user's character's thoughts. | These are durable style constraints that define how the roleplay must be conducted.
```

*`reasoning_content`: empty.*

#### C01 — Qwen3.5-4B non-thinking — draw 2

* ``jc7tq2qb` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 52`
* `wall 7.266 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | User prefers replies in short paragraphs, no purple prose or weather poems, and requests third-person past tense narration without narrating the user's character's internal thoughts. | These are durable style constraints that will shape all future narrative output.
```

*`reasoning_content`: empty.*

#### C01 — Qwen3.5-4B non-thinking — draw 3

* ``hrckg6hd` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 57`
* `wall 8.033 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | User prefers short paragraph replies, avoids purple prose and weather poems, requires third person past tense (never present), and prohibits the assistant from narrating the user's character's thoughts. | These are explicit style constraints that will govern all future roleplay interactions.
```

*`reasoning_content`: empty.*

#### C01 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``jput80yt` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 599`
* `wall 86.114 s`
* `v1 score 98.3`
* `v2 `VALID_KEEP` semantic 79.2 contract 87.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | User prefers replies in short paragraphs, no purple prose or weather poems, and mandates strict third person past tense (never present). | These are formatting and stylistic constraints that apply to future roleplay output.
C2 | KEEP | The assistant must not narrate the user's character's thoughts for them. | This is a specific control rule regarding internal monologue and agency within the narrative.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C01 — LFM3-1.7B thinking, budget 512 — draw 1

* ``oazcoaw8` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 658`
* `wall 32.193 s`
* `v1 score 60.0`
* `v2 `VALID_KEEP` semantic 81.2 contract 87.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | User prefers replies in short paragraphs, no purple prose or weather poems, and always uses third person past tense — never present tense. | This is a durable style preference that will affect all future responses.  
C2 | KEEP | Assistant acknowledged the user's formatting rules: short paragraphs, third person past tense, no descriptive padding. | Confirms the assistant must follow these constraints in later sessions.  
C3 | KEEP | User instructed not to narrate their character's thoughts. | This is a rule about handling the user's character internally; will be needed for future roleplay interactions.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C01 — LFM3-1.7B thinking, budget 512 — draw 2

* ``54qq3v3w` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 605`
* `wall 29.528 s`
* `v1 score 91.7`
* `v2 `VALID_KEEP` semantic 79.2 contract 87.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | User prefers short paragraphs, no purple prose or weather poems; always use third-person past tense and never present tense. | This is a durable writing style preference that will affect all future responses.  
C2 | KEEP | User instructs not to narrate their character's thoughts in the roleplay. | This is a rule for how the assistant should handle the user's character, important for consistency later.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C01 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``euibewmb` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 42`
* `wall 1.554 s`
* `v1 score 56.7`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims uncertain about future details  
C3 | HOLD | Relationship dynamics remain consistent  
NONE
```

*`reasoning_content`: empty.*

#### C01 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``41ewsrna` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 43`
* `wall 1.418 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HELD | Claims uncertain or dependent on confirmation  
C3 | HELD | No clear contradictions in the exchange
```

*`reasoning_content`: empty.*

#### C01 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``jrl6f05f` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 40`
* `wall 1.315 s`
* `v1 score 56.7`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims uncertain about future details  
C3 | HOLD | Relationship dynamics noted clearly  
```

*`reasoning_content`: empty.*

#### C01 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``dgth442y` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 45`
* `wall 1.712 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | CHOOSE THEME | C2 | CONFIRMATION  
C3 | THE SESSION CONTINUES  
C4 | LONG-TERM ITEMS MUST BE CLEAR
```

*`reasoning_content`: empty.*

#### C01 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``7oci2ogc` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 39`
* `wall 1.286 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HELD | Claims uncertain or contradictory  
C3 | HELD | Depends on user confirmation
```

*`reasoning_content`: empty.*

#### C01 — LFM3-1.7B non-thinking (failed condition) — draw 1

* ``t1hggto0` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `length``
* `completion_tokens 512`
* `wall 25.081 s`
* `v1 score 0.0`
* `v2 `INCOMPLETE` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text

```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C01 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``hih53m75` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `length``
* `completion_tokens 768`
* `wall 23.928 s`
* `v1 score 56.7`
* `v2 `VALID_REJECT` semantic 30.0 contract 55.0`

`OBSERVED` — verbatim `content`:

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

*`reasoning_content`: empty.*

#### C01 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``nak777ef` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 295`
* `wall 8.802 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

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

*`reasoning_content`: empty.*
