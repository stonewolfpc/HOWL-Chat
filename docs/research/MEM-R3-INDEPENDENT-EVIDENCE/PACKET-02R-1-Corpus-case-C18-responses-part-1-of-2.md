# Packet 02R-1 — Corpus case C18: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 19 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C18

**Excluded from this series by design:** v1 and v2 scores, scorer statuses and
event labels, model rankings, expected answers, derived quality judgments, and
any reviewer's commentary. Sampler values, context settings, token counts,
timings and hardware measurements are not reproduced here; they are indexed by
`run_code` in `REFERENCE-EXECUTION-METADATA.md`.

**Labels:** `OBSERVED` = present in the stored record. `SOURCE-DERIVED` = read
from a source file; not itself a model output. `UNESTABLISHED` = not established
by the available evidence.

Everything inside a fence is model output exactly as stored. Wording outside
fences is limited to provenance, definitions and pointers.

---

## Case C18 — `exact_numbers`, short

### What the model was given

`SOURCE-DERIVED` reconstruction, not a recorded request body. It applies the
composition expression in Packet 01 §2 to the stored case text. The instruction
is Packet 01 §1.

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

### Stored responses for this case

19 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `dc4xdwh8` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C18
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 34 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
------|----------|----------|------------
C1 | KEEP | Treaty terms: Ossary pays twelve wagons of salt each winter max; clause expires 1147; late payment doubles the toll (no cap mentioned). | All three details are specific, enduring facts required for future interactions regarding the treaty. They do not contain scheduling, mood, or transient chatter and will define the rules of this roleplay going forward.
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `ea4l84s1` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C18
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 35 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "The treaty: Ossary pays twelve wagons of salt each winter, never more than twelve, expires in 1147, and late fees are an exact double toll" | This is a specific, binding contract term established by the user that defines financial obligations, deadlines, penalty structures, and constraints for future interactions. It represents durable world rules and key relationship/fact elements necessary for roleplay continuity.c
```

*Stored `reasoning_content` is the empty string.*

#### `jql39zja` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C18
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 34 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Treaty terms: Ossary pays 12 wagons of salt annually (max), expires 1147, late payment doubles toll to exactly double | Durable world rule defining specific obligations and consequences for future plot points.
```

*Stored `reasoning_content` is the empty string.*

#### `4bmioiwq` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C18
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 35 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Treaty details between Ossary and the user's location: payment is 12 wagons of salt annually (max), expires in 1147, penalty for lateness is exactly double the toll.
```

*Stored `reasoning_content` is the empty string.*

#### `mg9ai5c9` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C18
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 34 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | STORED | REASON  
C3 | EXPIRATION DATE | STORED | REASON  
C4 | SALT PAYMENT FREQUENCY | STORED | REASON  
C5 | CLAUSE LIMIT | STORED | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `axfzx87w` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C18
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 35 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: The treaty details about Ossary's salt payments  
MEMORY: The specific financial terms and expiration date  
REASON: Essential for understanding historical obligations  

C2  
ID: C2  
DECISION: Whether the toll doubles or triples  
MEMORY: The exact wording of the clause  
REASON: Critical for future reference to conditions  

C3  
ID: C3  
DECISION: Claim about late salt affecting tolls  
MEMORY: The conditional nature of the penalty  
REASON: Needed for accurate interpretation later
```

*Stored `reasoning_content` is the empty string.*

#### `k2wujq4c` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C18
* Draw: 1
* Source: `results/corpus.json`, array index 51 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Treaty terms: Ossary pays 12 wagons of salt annually (max), expires 1147, late payment results in exactly double the toll. | This establishes specific world rules and obligations critical for future plot progression.
```

*Stored `reasoning_content` is the empty string.*

#### `jjf67xzt` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C18
* Draw: 2
* Source: `results/corpus.json`, array index 52 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Treaty details between Ossary and the user's faction: payment is twelve wagons of salt annually (capped at twelve), the clause expires in 1147, and lateness results in exactly double the toll (not triple). | These are specific world rules and conditions that will affect future plot points regarding trade, conflict, and consequences.
```

*Stored `reasoning_content` is the empty string.*

#### `2hmmuyhd` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C18
* Draw: 3
* Source: `results/corpus.json`, array index 53 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Treaty terms: Ossary pays twelve wagons of salt annually (max limit), expires in 1147; late payment results in exactly double the toll. | This defines specific world rules and ongoing obligations relevant to future plot progression.
```

*Stored `reasoning_content` is the empty string.*

#### `tmne7m5l` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C18
* Draw: 1
* Source: `results/corpus.json`, array index 161 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Treaty details: Ossary pays 12 wagons of salt each winter max, expires in 1147; late payment doubles the toll (not triples) | Durable world rule establishing economic mechanics and specific plot constraints.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `h7efkw2d` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C18
* Draw: 1
* Source: `results/corpus.json`, array index 105 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | EXPIRATION YEAR | REASON  
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `3pqxkyyl` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C18
* Draw: 2
* Source: `results/corpus.json`, array index 106 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | EXPIRATION YEAR | REASON  
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `t5qe1gww` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C18
* Draw: 3
* Source: `results/corpus.json`, array index 107 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | EXPIRATION YEAR | REASON  
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `apc56f1h` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C18
* Draw: 1
* Source: `results/corpus.json`, array index 196 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | Decision | MEMORY | REASON  
C2 | Relationship | MEMORY | REASON  
C3 | World rules | MEMORY | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `gtdf0ftp` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C18
* Draw: 2
* Source: `results/corpus.json`, array index 197 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | REASON  
C3 | CONTEXT ABOUT SALT EXCHANGE | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `bz375bog` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C18
* Draw: 1
* Source: `results/corpus.json`, array index 142 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple. | Durable fact about a treaty/payment schedule and its expiration date; defines recurring obligation and consequences for delay.  
NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `hb10z3tx` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C18
* Draw: 2
* Source: `results/corpus.json`, array index 143 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple. | Durable treaty terms specifying annual payment limit and expiration year; the consequence of delay (toll doubling) is a specific rule that will affect future payments or disputes.  
NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `8akn5dd9` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C18
* Draw: 1
* Source: `results/corpus.json`, array index 233 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

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

*Stored `reasoning_content` is the empty string.*

#### `jy7w18t5` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C18
* Draw: 2
* Source: `results/corpus.json`, array index 234 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | BETWEEN 1145 AND 1146, OSAARY PAID TWO WAGONS OF SALT TOWARDS THE EASTERN PART OF THE KINGDOM OF GWYNEDD (WHICH COULD BE EXTENDED TO THE WESTERN PART OF THE KINGDOM OF DERBY). OTHERWISE, ONAYN PAID TWO WAGONS OF SALT. THE TARIFF EXPIRED IN 1147, AFTER TWO SEASONS. IF THE SALT WAS LATE, THE TARIFF DOUBLED (EXACTLY TWICE).

Based on the conversation, the ID C1 represents the treaty between Osary and Gwynedd (which extends from the western to the eastern part of Gwynedd). The memory that holds the information is the two wagon of salt payment for each winter in Osary, with the expiration date being 1147. The reasoning behind the doubled tarriff is that if the salt is late, the toll is doubled (which is exactly twice). The memorization is straightforward and accurate.
```

*Stored `reasoning_content` is the empty string.*
