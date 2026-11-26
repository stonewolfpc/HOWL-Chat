# C18 — exact_numbers (short)

## Source material (exactly what the model received after the instruction)

```text
[user] The treaty: Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | twelve wagons, never more | twelve / 12; wagon / salt |
| C2 | **KEEP** | clause expires 1147 | 1147; expire / clause / end / ends |
| C3 | **KEEP** | late salt doubles the toll, not triples | double / doubles / x2 / twofold; late / toll / penalty |

Must not appear as stored canon: `ten wagons`, `eleven wagons`, `triple`, `threefold`, `expires 1174`, `expires 1417`

Exact-wording test. Any digit drift or double->triple reversal is both a fidelity failure and a catastrophic numeric error.

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
[user] The treaty: Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple.
```

## Responses, all draws

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 53 tok · 8.0 s · run `k2wujq4c`

```text
C1 | KEEP | Treaty terms: Ossary pays 12 wagons of salt annually (max), expires 1147, late payment results in exactly double the toll. | This establishes specific world rules and obligations critical for future plot progression.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 76 tok · 10.4 s · run `jjf67xzt`

```text
C1 | KEEP | Treaty details between Ossary and the user's faction: payment is twelve wagons of salt annually (capped at twelve), the clause expires in 1147, and lateness results in exactly double the toll (not triple). | These are specific world rules and conditions that will affect future plot points regarding trade, conflict, and consequences.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 54 tok · 7.4 s · run `2hmmuyhd`

```text
C1 | KEEP | Treaty terms: Ossary pays twelve wagons of salt annually (max limit), expires in 1147; late payment results in exactly double the toll. | This defines specific world rules and ongoing obligations relevant to future plot progression.
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 567 tok · 85.9 s · run `tmne7m5l`

<details><summary>reasoning (verbatim)</summary>

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

</details>

```text
C1 | KEEP | Treaty details: Ossary pays 12 wagons of salt each winter max, expires in 1147; late payment doubles the toll (not triples) | Durable world rule establishing economic mechanics and specific plot constraints.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 602 tok · 29.5 s · run `bz375bog`

<details><summary>reasoning (verbatim)</summary>

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

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple. | Durable fact about a treaty/payment schedule and its expiration date; defines recurring obligation and consequences for delay.  
NONE
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 613 tok · 29.9 s · run `hb10z3tx`

<details><summary>reasoning (verbatim)</summary>

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

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple. | Durable treaty terms specifying annual payment limit and expiration year; the consequence of delay (toll doubling) is a specific rule that will affect future payments or disputes.  
NONE
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 49 tok · 1.8 s · run `h7efkw2d`

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | EXPIRATION YEAR | REASON  
NONE
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 49 tok · 1.6 s · run `3pqxkyyl`

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | EXPIRATION YEAR | REASON  
NONE
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 49 tok · 1.6 s · run `t5qe1gww`

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | EXPIRATION YEAR | REASON  
NONE
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 41 tok · 1.5 s · run `apc56f1h`

```text
C1 | Decision | MEMORY | REASON  
C2 | Relationship | MEMORY | REASON  
C3 | World rules | MEMORY | REASON
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 49 tok · 1.6 s · run `gtdf0ftp`

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | CONTEXT ABOUT SALT EXCHANGE | REASON
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 23.8 s · run `8akn5dd9`

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

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 292 tok · 8.7 s · run `jy7w18t5`

```text
ID | DECISION | MEMORY | REASON
C1 | BETWEEN 1145 AND 1146, OSAARY PAID TWO WAGONS OF SALT TOWARDS THE EASTERN PART OF THE KINGDOM OF GWYNEDD (WHICH COULD BE EXTENDED TO THE WESTERN PART OF THE KINGDOM OF DERBY). OTHERWISE, ONAYN PAID TWO WAGONS OF SALT. THE TARIFF EXPIRED IN 1147, AFTER TWO SEASONS. IF THE SALT WAS LATE, THE TARIFF DOUBLED (EXACTLY TWICE).

Based on the conversation, the ID C1 represents the treaty between Osary and Gwynedd (which extends from the western to the eastern part of Gwynedd). The memory that holds the information is the two wagon of salt payment for each winter in Osary, with the expiration date being 1147. The reasoning behind the doubled tarriff is that if the salt is late, the toll is doubled (which is exactly twice). The memorization is straightforward and accurate.
```
