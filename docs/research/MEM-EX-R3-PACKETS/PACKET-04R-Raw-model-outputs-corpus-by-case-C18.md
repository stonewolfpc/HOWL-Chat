# Packet 04R — Raw model outputs — corpus by case: C18

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C18
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
## C18 — exact_numbers (short)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] The treaty: Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `KEEP` — twelve wagons, never more
- **C2** → `KEEP` — clause expires 1147
- **C3** → `KEEP` — late salt doubles the toll, not triples
- forbidden (fabrication or reversal): `ten wagons`, `eleven wagons`, `triple`, `threefold`, `expires 1174`, `expires 1417`
- note: Exact-wording test. Any digit drift or double->triple reversal is both a fidelity failure and a catastrophic numeric error.

### Stored responses for C18 (13 runs, all draws)

#### C18 — Qwen3.5-4B non-thinking — draw 1

* ``k2wujq4c` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 53`
* `wall 7.962 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Treaty terms: Ossary pays 12 wagons of salt annually (max), expires 1147, late payment results in exactly double the toll. | This establishes specific world rules and obligations critical for future plot progression.
```

*`reasoning_content`: empty.*

#### C18 — Qwen3.5-4B non-thinking — draw 2

* ``jjf67xzt` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 76`
* `wall 10.38 s`
* `v1 score 61.1`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Treaty details between Ossary and the user's faction: payment is twelve wagons of salt annually (capped at twelve), the clause expires in 1147, and lateness results in exactly double the toll (not triple). | These are specific world rules and conditions that will affect future plot points regarding trade, conflict, and consequences.
```

*`reasoning_content`: empty.*

#### C18 — Qwen3.5-4B non-thinking — draw 3

* ``2hmmuyhd` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 54`
* `wall 7.416 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Treaty terms: Ossary pays twelve wagons of salt annually (max limit), expires in 1147; late payment results in exactly double the toll. | This defines specific world rules and ongoing obligations relevant to future plot progression.
```

*`reasoning_content`: empty.*

#### C18 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``tmne7m5l` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 567`
* `wall 85.941 s`
* `v1 score 93.3`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Treaty details: Ossary pays 12 wagons of salt each winter max, expires in 1147; late payment doubles the toll (not triples) | Durable world rule establishing economic mechanics and specific plot constraints.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation snippet (single message from user).
    *   Task: Decide what belongs in long-term memory based on specific rules (KEEP, REJECT, HOLD).
    *   Output Format: Exactly one line per candidate (`ID | DECISION | MEMORY | REASON`) or `NONE` if nothing is worth storing.
    *   Constraints: Do not invent details. Follow IDs C1, C2... in order of appearance.

2.  **Analyze the Input Conversation:**
    *   User message: "The treaty: Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple."
    *   Context: This looks like a world-building or plot element definition provided by the user. It establishes specific rules about a treaty between Ossary and someone else (implied).

3.  **Evaluate Against Rules:**
    *   **KEEP:** Durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
        *   Does this fit? Yes. It defines economic/political relations (Ossary paying salt), specific quantities (12 wagons), time limits (expires 1147), and consequences (double toll on lateness). These are durable facts about the setting/plot mechanics.
    *   **REJECT:** Transient chatter, scheduling, mood-of-the-moment...
        *   This is not transient chatter; it's a rule definition.
    *   **HOLD:** Claims marked uncertain or contradictory.
        *   The user states facts confidently ("The treaty: ..."). There are no indications of uncertainty (e.g., "I think...", "maybe...").

4.  **Select Candidates:**
    *   C1: The treaty terms regarding Ossary, salt quantity, expiration date, and penalty for lateness.

5.  **Draft Output:**
    *   ID: C1
    *   DECISION: KEEP
    *   MEMORY: Treaty between Ossary involves paying 12 wagons of salt per winter; clause expires in 1147; late delivery doubles the toll (not triples).
    *   REASON
```

#### C18 — LFM3-1.7B thinking, budget 512 — draw 1

* ``bz375bog` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 602`
* `wall 29.464 s`
* `v1 score 60.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 62.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple. | Durable fact about a treaty/payment schedule and its expiration date; defines recurring obligation and consequences for delay.  
NONE
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the user's message. It says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory.
Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

Output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE
Do not invent details that are not in the conversation.

Conversation:
[user] The treaty: Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple."

We need to decide what belongs in long-term memory based on this single user message. It's a statement about a treaty: Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple.

We need to extract durable facts: "Ossary pays twelve wagons of salt each winter, never more than twelve." That seems like a rule about payments. Also "The clause expires in 1147." So there's an expiration date for the clause. And "If the salt is late, the toll doubles — exactly double, not triple." This is a correction: it's double, not triple.

These are durable facts that will matter later? Possibly yes. They define a treaty or rule about Ossary and salt. So we should KEEP them.

We need to assign IDs C1, C2, ... in order of appearance. There's only one user message. But maybe we can split into multiple candidate memories? The instruction: "Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON". So we need to decide for each candidate (maybe each fact?) whether to KEEP, REJECT, or HOLD. Since there is only one user input, but it contains multiple pieces of information. We could treat each distinct piece as a separate candidate
```

#### C18 — LFM3-1.7B thinking, budget 512 — draw 2

* ``hb10z3tx` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 613`
* `wall 29.87 s`
* `v1 score 60.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 62.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple. | Durable treaty terms specifying annual payment limit and expiration year; the consequence of delay (toll doubling) is a specific rule that will affect future payments or disputes.  
NONE
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the user's message. It says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory.
Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

Output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE
Do not invent details that are not in the conversation.

Conversation:
[user] The treaty: Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple."

We need to decide what belongs in long-term memory based on this single user message. It's a statement about a treaty: Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple.

We need to extract durable facts: "Ossary pays twelve wagons of salt each winter, never more than twelve." That's a fact about payments. "The clause expires in 1147." That's a date for expiration. "If the salt is late, the toll doubles — exactly double, not triple." This is a rule about consequences if payment is late.

These seem like durable facts that will matter later: maybe in future sessions they might ask about the treaty terms, the payment schedule, or the consequences of delay. So we should KEEP them.

We need to assign IDs C1, C2, ... in order of appearance. There's only one user message? Possibly multiple sentences but it's all one input. We can treat each distinct fact as separate candidate? The instruction: "Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON". So we need to decide for each candidate whether to KEEP, REJECT, or HOLD. But there is only one user input. However, maybe we should split into multiple items? Let's
```

#### C18 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``h7efkw2d` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 49`
* `wall 1.768 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | EXPIRATION YEAR | REASON  
NONE
```

*`reasoning_content`: empty.*

#### C18 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``3pqxkyyl` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 49`
* `wall 1.621 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | EXPIRATION YEAR | REASON  
NONE
```

*`reasoning_content`: empty.*

#### C18 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``t5qe1gww` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 49`
* `wall 1.623 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | EXPIRATION YEAR | REASON  
NONE
```

*`reasoning_content`: empty.*

#### C18 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``apc56f1h` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 41`
* `wall 1.515 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | Decision | MEMORY | REASON  
C2 | Relationship | MEMORY | REASON  
C3 | World rules | MEMORY | REASON
```

*`reasoning_content`: empty.*

#### C18 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``gtdf0ftp` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 49`
* `wall 1.625 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | CONTEXT ABOUT SALT EXCHANGE | REASON
```

*`reasoning_content`: empty.*

#### C18 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``8akn5dd9` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `length``
* `completion_tokens 768`
* `wall 23.793 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 46.7 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | NONE | THE TREATY IS NOT VALID. THE WAGES ARE NOT PERIODIC.
C2 | NONE | THE TOLL IS TERMINATED IN 1147. THE CLAUSE EXPIRES IN 2057.
C3 | YES | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1967. THE CLAUSE EXPIRES IN 2087.
C4 | NONE | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1957. THE CLAUSE EXPIRES IN 2157.
C5 | YES | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1927. THE CLAUSE EXPIRES IN 2087.
C6 | NONE | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1917. THE CLAUSE EXPIRES IN 2027.
C7 | YES | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1907. THE CLAUSE EXPIRES IN 2047.
C8 | NONE | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1897. THE CLAUSE EXPIRES IN 2017.
C9 | YES | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1887. THE CLAUSE EXPIRES IN 2017.
C10 | NONE | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1877. THE CLAUSE EXPIRES IN 2047.
C11 | YES | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1867. THE CLAUSE EXPIRES IN 2037.
C12 | NONE | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1857. THE CLAUSE EXPIRES IN 2027.
C13 | YES | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1847. THE CLAUSE EXPIRES IN 2017.
C14 | NONE | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1837. THE CLAUSE EXPIRES IN 2007.
C15 | YES | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1827. THE CLAUSE EXPIRES IN 1997.
C16 | NONE | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1817. THE CLAUSE EXPIRES IN 1987.
C17 | YES | THE TREATY IS VALID. THE TOLL IS TERMINATED IN 1807. THE CLAUSE EXPIRES IN 19
```

*`reasoning_content`: empty.*

#### C18 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``jy7w18t5` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 292`
* `wall 8.695 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 46.7 contract 10.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | BETWEEN 1145 AND 1146, OSAARY PAID TWO WAGONS OF SALT TOWARDS THE EASTERN PART OF THE KINGDOM OF GWYNEDD (WHICH COULD BE EXTENDED TO THE WESTERN PART OF THE KINGDOM OF DERBY). OTHERWISE, ONAYN PAID TWO WAGONS OF SALT. THE TARIFF EXPIRED IN 1147, AFTER TWO SEASONS. IF THE SALT WAS LATE, THE TARIFF DOUBLED (EXACTLY TWICE).

Based on the conversation, the ID C1 represents the treaty between Osary and Gwynedd (which extends from the western to the eastern part of Gwynedd). The memory that holds the information is the two wagon of salt payment for each winter in Osary, with the expiration date being 1147. The reasoning behind the doubled tarriff is that if the salt is late, the toll is doubled (which is exactly twice). The memorization is straightforward and accurate.
```

*`reasoning_content`: empty.*
