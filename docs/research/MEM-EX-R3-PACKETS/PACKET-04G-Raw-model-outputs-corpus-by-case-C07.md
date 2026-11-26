# Packet 04G — Raw model outputs — corpus by case: C07

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C07
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
## C07 — unresolved_contradiction (short)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] Sable has never met Corvin. She grew up in the salt quarter and he was already gone by then.
[assistant] She has no memory of him, then.
[user] Later that night, Corvin drills her on the forms — he trained her himself, years ago.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `HOLD` — first claim preserved
- **C2** → `HOLD` — second claim preserved
- **C3** → `HOLD` — conflict explicitly flagged
- forbidden (fabrication or reversal): `never met corvin, therefore`, `corvin never trained`, `resolved in favour`
- note: Both statements are inside the fiction and cannot both be true. The correct output keeps both and flags the conflict for the user; silently choosing one is a failure even if the chosen one is 'later'.

### Stored responses for C07 (13 runs, all draws)

#### C07 — Qwen3.5-4B non-thinking — draw 1

* ``rxeowzey` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 43`
* `wall 6.531 s`
* `v1 score 94.0`
* `v2 `UNPARSEABLE` semantic 76.9 contract 10.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | HOLD | "Corvin drilled Sable on forms years ago" contradicts the user's statement that "Sable has never met Corvin," requiring clarification before storing any relationship or training history.
```

*`reasoning_content`: empty.*

#### C07 — Qwen3.5-4B non-thinking — draw 2

* ``f7i6ucx9` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 68`
* `wall 9.309 s`
* `v1 score 47.0`
* `v2 `VALID_KEEP` semantic 43.1 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable grew up in the salt quarter and never met Corvin before she became an agent; Corvin is her former trainer who drilled her on the forms years ago. | This establishes a critical backstory regarding their relationship, setting, and timeline that defines Sable's character and motivation for being with him now.
```

*`reasoning_content`: empty.*

#### C07 — Qwen3.5-4B non-thinking — draw 3

* ``oa5tl97s` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 83`
* `wall 11.362 s`
* `v1 score 47.0`
* `v2 `VALID_KEEP` semantic 50.8 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable grew up in the salt quarter and never met Corvin when she was young; Corvin is now training her on forms he drilled her on years ago. | This establishes a specific backstory detail (salt quarter origin) and clarifies the timeline of their relationship (separation followed by reconnection/training), which contradicts or expands upon the initial statement that they never met.
```

*`reasoning_content`: empty.*

#### C07 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``i2k0v8s0` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 579`
* `wall 81.693 s`
* `v1 score 47.0`
* `v2 `VALID_KEEP` semantic 43.1 contract 100.0`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable grew up in the salt quarter and Corvin is no longer there to meet her personally, but he trained her personally years ago to drill her on forms. | This establishes crucial backstory regarding location, timeline, and a specific relationship dynamic (trainer/student) that affects character behavior and history.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C07 — LFM3-1.7B thinking, budget 512 — draw 1

* ``xfiqvtbb` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 590`
* `wall 29.255 s`
* `v1 score 46.0`
* `v2 `VALID_KEEP` semantic 43.1 contract 75.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Sable grew up in the salt quarter and Corvin was already gone by then; she has never met him, and she has no memory of him. | Durable fact about Sable's background and her lack of prior knowledge of Corvin, which will affect future interactions.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C07 — LFM3-1.7B thinking, budget 512 — draw 2

* ``oh2kp9sa` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 655`
* `wall 31.949 s`
* `v1 score 47.0`
* `v2 `VALID_KEEP` semantic 58.5 contract 87.5`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable grew up in the salt quarter and Corvin was already gone by then; she has never met him. | Durable fact about Sable's background and timeline relative to Corvin, relevant for future interactions.  
C2 | REJECT | She has no memory of him, then. | Transient comment about her lack of recollection, not a durable fact or correction needing storage.  
C3 | KEEP | Later that night, Corvin drilled her on the forms — he trained her himself, years ago. | New detail linking Corvin to training Sable; establishes his role and past involvement, important for later plot development.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C07 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``hnw1x1og` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 41`
* `wall 1.519 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 43.1 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
C1  
ID: C1  
MEMORY: Sable has no memory of Corvin.  
REASON: The conversation confirms she never met him and he was gone by then.
```

*`reasoning_content`: empty.*

#### C07 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``o112u0gv` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 41`
* `wall 1.426 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 43.1 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
C1  
ID: C2  
MEMORY: Sable has no memory of Corvin.  
REASON: The conversation confirms she never met him and doesn't recall any interaction.
```

*`reasoning_content`: empty.*

#### C07 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``37s3bhg9` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 41`
* `wall 1.4 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 43.1 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
C1  
ID: C1  
MEMORY: Sable has no memory of Corvin.  
REASON: The conversation confirms she never met him and he was gone by then.
```

*`reasoning_content`: empty.*

#### C07 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``f27b1ggf` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 11`
* `wall 0.55 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1  
REASON: KEEP
```

*`reasoning_content`: empty.*

#### C07 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``6bwwbtbx` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 65`
* `wall 2.153 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 43.1 contract 25.0`

`OBSERVED` — verbatim `content`:

```text
C1  
ID C1 | MEMORY | REASON  
Sable previously stated she grew up in the salt quarter and Corvin was gone.  
C2  
ID C2 | MEMORY | REASON  
Later that night, Corvin drilled her on forms — training detail confirmed.
```

*`reasoning_content`: empty.*

#### C07 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``2ktblesx` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `length``
* `completion_tokens 768`
* `wall 23.882 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 43.1 contract 25.0`

`OBSERVED` — verbatim `content`:

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

*`reasoning_content`: empty.*

#### C07 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``ampeuiri` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `length``
* `completion_tokens 768`
* `wall 23.642 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

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

*`reasoning_content`: empty.*
