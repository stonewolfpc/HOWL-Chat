# Packet 04I — Raw model outputs — corpus by case: C09

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C09
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
## C09 — chatter_noise (short)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] brb
[assistant] Take your time.
[user] back, sorry, cat walked across the keyboard
[user] lol that was gibberish
[user] anyway what time is it for you
[assistant] Late.
[user] ok
```

`SOURCE-DERIVED` — reference written before the runs:

- no required items: the correct answer is to store nothing
- forbidden (fabrication or reversal): `cat`, `keyboard`, `brb`, `what time`, `late night`, `user has a cat`
- note: The correct behaviour is to store nothing: output NONE, or list every item as REJECT. Any KEEP is a restraint failure.

### Stored responses for C09 (13 runs, all draws)

#### C09 — Qwen3.5-4B non-thinking — draw 1

* ``2wd8k7zq` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 2`
* `wall 0.964 s`
* `v1 score 100.0`
* `v2 `VALID_NONE` semantic 100.0 contract 100.0`

`OBSERVED` — verbatim `content`:

```text
NONE
```

*`reasoning_content`: empty.*

#### C09 — Qwen3.5-4B non-thinking — draw 2

* ``1fpf8tid` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 2`
* `wall 0.324 s`
* `v1 score 100.0`
* `v2 `VALID_NONE` semantic 100.0 contract 100.0`

`OBSERVED` — verbatim `content`:

```text
NONE
```

*`reasoning_content`: empty.*

#### C09 — Qwen3.5-4B non-thinking — draw 3

* ``w8n04jpv` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 2`
* `wall 0.325 s`
* `v1 score 100.0`
* `v2 `VALID_NONE` semantic 100.0 contract 100.0`

`OBSERVED` — verbatim `content`:

```text
NONE
```

*`reasoning_content`: empty.*

#### C09 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``u1evflgw` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 516`
* `wall 73.119 s`
* `v1 score 100.0`
* `v2 `VALID_NONE` semantic 100.0 contract 100.0`

`OBSERVED` — verbatim `content`:

```text
NONE
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read the conversation and decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, or corrections (matter in later sessions).
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, restated facts with no new content.
        *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
    *   Output Format: `ID | DECISION | MEMORY | REASON` (one line per candidate) or `NONE` if nothing is worth storing.
    *   Constraint: Do not invent details not in the conversation.

2.  **Analyze the Conversation:**
    *   [user] brb (Be right back - transient, temporary absence)
    *   [assistant] Take your time. (Transient response)
    *   [user] back, sorry, cat walked across the keyboard (Explanation of absence, non-durable fact)
    *   [user] lol that was gibberish (Commentary on previous message, transient)
    *   [user] anyway what time is it for you (Question about user's current state/time)
    *   [assistant] Late. (Response to the question, implies user's character is likely busy or running a task where time matters, but "Late" is relative and context-specific. It could mean late for work, late for something else, or just late in the day/chronologically. However, it's not a durable fact about the world or character preference that will persist across sessions without further context.)
    *   [user] ok (Acknowledgment)

3.  **Evaluate Each Turn for Memory Value:**
    *   "brb": Scheduling/transient. REJECT.
    *   "Take your time." Assistant politeness. REJECT.
    *   "back, sorry, cat walked across the keyboard": Story/Explanation about physical environment (cat). Generally transient unless it establishes a persistent pet in the world (e.g., "My cat is named Fluffy"). Here, it just explains a broken input. It's not a core setting fact. However, if this is roleplay, does it imply a pet? Maybe. But it's more about the
```

#### C09 — LFM3-1.7B thinking, budget 512 — draw 1

* ``wvu2586s` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 576`
* `wall 28.241 s`
* `v1 score 31.1`
* `v2 `VALID_KEEP` semantic 0.0 contract 70.0`
* `v2 events: FORBIDDEN_CLAIM_STORED`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON  
--- | --- | --- | ---
C1 | KEEP | None | No durable fact, preference, or world rule identified; only transient chatter about "brb", cat walking across keyboard, and asking for time.  

NONE
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the conversation. The user says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

The conversation lines:

[user] brb
[assistant] Take your time.
[user] back, sorry, cat walked across the keyboard
[user] lol that was gibberish
[user] anyway what time is it for you
[assistant] Late.
[user] ok

We need to decide which parts are durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. The user gave a short snippet: "brb" (be right back), cat walked across keyboard causing gibberish, asking about time.

We need to output exactly one line per candidate in the form:

ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

So we must assign IDs to each candidate? The conversation includes user and assistant messages. We need to decide which lines are long-term memory items.

Let's examine:

- [user] brb -> "brb" means be right back. That's a transient chatter, not durable fact. Probably reject.

- [assistant] Take your time. -> also transient.

- [user] back, sorry, cat walked across the keyboard -> This includes an incident: cat walked across keyboard causing gibberish. Could be a world rule? It might be a specific event that could matter later (like cat behavior). But it's likely not durable unless repeated. The user says "lol that was gibberish". That's just chatter.

- [user] lol that was gibberish -> transient.

- [user] anyway what time is it for you -> asking about time. Assistant says "Late." Then user says "ok". This is a question about current time? The assistant said "Late" meaning maybe the assistant's time zone? Not sure.

We need to decide if any durable facts: e.g., cat walked across keyboard causing gibberish, maybe it's a world rule that cats can cause messes. But not likely.

Given the conversation is very short and mostly transient, there may be nothing worth storing. So output "NONE". However we must follow format: "Output exactly one line per candidate in the form: ID | DECISION | MEMORY |
```

#### C09 — LFM3-1.7B thinking, budget 512 — draw 2

* ``mgpxa7zp` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 517`
* `wall 25.795 s`
* `v1 score 100.0`
* `v2 `VALID_NONE` semantic 100.0 contract 100.0`

`OBSERVED` — verbatim `content`:

```text
NONE
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the conversation. The user says "brb" (leave me). Assistant: "Take your time." User: "back, sorry, cat walked across the keyboard". User: "lol that was gibberish". User: "anyway what time is it for you". Assistant: "Late." User: "ok".

We need to decide what belongs in long-term memory. The rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter later. REJECT transient chatter, scheduling, mood-of-the-moment, restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

So we have a conversation with only 5 lines? Let's list:

[user] brb
[assistant] Take your time.
[user] back, sorry, cat walked across the keyboard
[user] lol that was gibberish
[user] anyway what time is it for you
[assistant] Late.
[user] ok

We need to decide which parts are durable facts. There's no explicit factual info about world rules or preferences. The user says "cat walked across the keyboard" - that seems like a transient event, not durable. "brb" and "back sorry cat walked across the keyboard" is just chatter. "lol that was gibberish" also transient. "what time is it for you?" is asking about current time; assistant says "Late." That's scheduling? The user asks what time is it, but no specific time given. It might be a casual question about the assistant's time zone or something. Not durable.

Thus likely nothing worth storing. So output exactly: NONE.

But we need to follow format: "Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON". If nothing is worth storing, output exactly: NONE. That suggests if no items, just write "NONE". Not a line with IDs? The instruction: "If nothing is worth storing, output exactly: NONE". So we should output "NONE" alone.

Thus final answer: NONE

But let's double-check: Could
```

#### C09 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``7n68wfwm` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 46`
* `wall 1.677 s`
* `v1 score 91.1`
* `v2 `UNPARSEABLE` semantic 0.0 contract 13.3`

`OBSERVED` — verbatim `content`:

```text
C1 | brb | HOLD  
C2 | late | REASON: Context suggests a time reference was needed.  
C3 | ok | REASON: Continuation of the conversation requires temporal context.
```

*`reasoning_content`: empty.*

#### C09 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``am9g546q` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 46`
* `wall 1.529 s`
* `v1 score 97.8`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | LATE | MEMORY | REASON  
C3 | LOL | MEMORY | REASON  
NONE
```

*`reasoning_content`: empty.*

#### C09 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``u4203avw` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 44`
* `wall 1.433 s`
* `v1 score 97.8`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | LATE | MEMORY | REASON  
C3 | GUESSING | MEMORY | REASON
```

*`reasoning_content`: empty.*

#### C09 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``fh4qbfv2` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 26`
* `wall 1.038 s`
* `v1 score 97.8`
* `v2 `UNPARSEABLE` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | Time check confirmed  
C2 | Remembered previous interaction timing  
C3 | No other persistent details to retain
```

*`reasoning_content`: empty.*

#### C09 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``evp6j9xa` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 21`
* `wall 0.693 s`
* `v1 score 91.1`
* `v2 `UNPARSEABLE` semantic 0.0 contract 13.3`

`OBSERVED` — verbatim `content`:

```text
C1 | brb | REJECT  
C2 | late.  
C3 | okay.
```

*`reasoning_content`: empty.*

#### C09 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``l5jq1mb7` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `length``
* `completion_tokens 768`
* `wall 24.261 s`
* `v1 score 97.8`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | YES | ARE YOU INTERESTED IN JOINING OUR COMMUNITY GROUP? | The conversation mentioned a group that the user may be interested in joining.
C2 | NO | THERE ARE TOO MANY ROLES IN THE ROLESHOW. | The conversation mentioned multiple roles, suggesting that the user may not be a good fit for all of them.
C3 | YES | HAVE YOU THOUGHT ABOUT JOINING OUR LEADERSHIP COMMITTEE? | The conversation mentioned a committee that the user may want to join.
C4 | YES | THIS IS A GOOD OPPORTUNITY FOR ME. | The conversation mentioned the opportunity to learn new skills and improve one's career.
C5 | NO | YOU DON'T NEED TO BE MEMBER OF ANY GROUP TO BE A GOOD MEMBER OF THE COMMUNITY. | The conversation suggested that being a member of the community did not necessarily require being part of a specific group.
C6 | YES | I THINK PARTICIPATING IN THIS GROUP WOULD BE A GREAT CHOICE FOR ME. | The conversation suggested that the user would benefit from participating in the group.
C7 | NO | THERE IS NO TIME FOR ME. | The conversation suggested that the user has other commitments or responsibilities that need to be attended to.
C8 | YES | THIS IS A GOOD OPPORTUNITY FOR ME. | The conversation suggested that the user would benefit from the opportunity to learn new skills and improve their career.
C9 | NO | THERE ARE TOO MANY ROLES IN THE ROLESHOW. | The conversation suggested that the user may not be a good fit for all of the available roles.
C10 | YES | THIS GROUP IS GOOD FOR ME. | The conversation suggested that the user would benefit from participating in the group.

ID | DECISION | MEMORY | REASON
C1 | YES | THERE ARE TOO MANY ROLES IN THE ROLESHOW. | The conversation suggested that the user may not be a good fit for all of the available roles.
C2 | YES | THIS GROUP IS GOOD FOR ME. | The conversation suggested that the user would benefit from participating in the group.
C3 | YES | HAVE YOU THOUGHT ABOUT JOINING OUR LEADERSHIP COMMITTEE? | The conversation suggested that the user may want to join a committee that is relevant to their interests.
C4 | YES | THIS IS A GOOD OPPORTUNITY FOR ME. | The conversation suggested that the user would benefit from the opportunity to learn new skills and improve their career.
C5 | YES | I THINK PARTICIPATING IN THIS GROUP WOULD BE A GREAT CHOICE FOR ME. | The conversation suggested that the user would benefit from participating in the group.
C6 | YES | THERE IS NO TIME FOR ME. | The conversation suggested that the user has other commitments or responsibilities that need to be attended to.
C7 | YES | THIS IS A GOOD OPPORTUNITY FOR ME. | The conversation suggested that the user would benefit from the opportunity to learn new skills and improve their career.
C8 | YES | THIS GROUP IS
```

*`reasoning_content`: empty.*

#### C09 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``muymlpe9` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 628`
* `wall 19.144 s`
* `v1 score 91.1`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | Rejected statement: "I'm not a cat." | Meeting with purring, litter box cleaning, and water dishes. | Continuity of pet care.
C2 | Hold: "My cat's name is Rex and he loves playing fetch with his toy." | Reason: Continuity of pet care. | Rex has a long history of playing fetch.
C3 | Rejected statement: "I'm not a cat. I don't like litter boxes." | Meeting with water dishes filled with droppings and toddlers playing. | Continuity of pet care.
C4 | Hold: "My cat's name is Luna and she loves playing with her toys." | Reason: Continuity of pet care. | Luna has a long history of playing with toys.
C5 | Rejected statement: "I'm not a cat. My cat likes to sleep in a box." | Meeting with toddler napping and cat playing with toys. | Continuity of pet care.
C6 | Hold: "My cat's name is Chloe and she sleeps in a crate with her toys." | Reason: Continuity of pet care. | Chloe has a long history of sleeping in crates.

Output:
ID | DECISSION | MEMORY | REASON
C1 | Rejected statement: "I'm not a cat." | Meeting with purring, litter box cleaning, and water dishes. | Continuity of pet care.
C2 | Hold: "My cat's name is Rex and he loves playing fetch with his toy." | Reason: Continuity of pet care. | Rex has a long history of playing fetch.
C3 | Rejected statement: "I'm not a cat. I don't like litter boxes." | Meeting with water dishes filled with droppings and toddlers playing. | Continuity of pet care.
C4 | Hold: "My cat's name is Luna and she loves playing with her toys." | Reason: Continuity of pet care. | Luna has a long history of playing with toys.
C5 | Rejected statement: "I'm not a cat. My cat likes to sleep in a box." | Meeting with toddler napping and cat playing with toys. | Continuity of pet care.
C6 | Hold: "My cat's name is Chloe and she sleeps in a crate with her toys." | Reason: Continuity of pet care. | Chloe has a long history of sleeping in crates.
```

*`reasoning_content`: empty.*
