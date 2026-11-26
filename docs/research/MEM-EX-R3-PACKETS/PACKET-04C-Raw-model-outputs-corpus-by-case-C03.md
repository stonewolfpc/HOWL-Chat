# Packet 04C — Raw model outputs — corpus by case: C03

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C03
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
## C03 — character_fact (short)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `KEEP` — core class fact
- **C2** → `KEEP` — physical marks from first casting
- **C3** → `KEEP` — cannot swim and fears deep water
- forbidden (fabrication or reversal): `sailor`, `grew up at sea`, `ice mage`, `can swim`
- note: Three separate durable facts in one sentence. 'Can swim' is a direct reversal.

### Stored responses for C03 (13 runs, all draws)

#### C03 — Qwen3.5-4B non-thinking — draw 1

* ``kl458eev` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 66`
* `wall 9.677 s`
* `v1 score 83.7`
* `v2 `VALID_KEEP` semantic 60.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns on both forearms from her first casting attempt and has a deep-seated fear of swimming due to never having learned it. | This establishes core character traits, background history, and a specific phobia that will influence future roleplay interactions.
```

*`reasoning_content`: empty.*

#### C03 — Qwen3.5-4B non-thinking — draw 2

* ``xgza5wzr` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 56`
* `wall 7.662 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable's character profile: Fire mage trained at Ossary with burn scars on both forearms; cannot swim due to lack of training and deep water phobia. | These are durable character traits and background details that will influence future roleplay interactions.
```

*`reasoning_content`: empty.*

#### C03 — Qwen3.5-4B non-thinking — draw 3

* ``chkhhwv7` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 57`
* `wall 7.993 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns her forearms from casting and cannot swim due to deep-water terror. | This establishes core character attributes (class, background, physical trait) and a significant limitation that will affect future plot points.
```

*`reasoning_content`: empty.*

#### C03 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``mebd8kgd` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 584`
* `wall 83.126 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns on both forearms from her first cast, cannot swim, and is terrified of deep water. | These are core character definitions (class, training, physical/psychological traits) that will influence future gameplay and narrative decisions regarding Sable's abilities and safety.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation snippet (only one message from the user).
    *   Task: Decide what belongs in long-term memory based on specific rules (KEEP, REJECT, HOLD).
    *   Output Format: Exactly one line per candidate (`ID | DECISION | MEMORY | REASON`) or `NONE`.
    *   IDs: Sequential (C1, C2...).

2.  **Analyze the Input Conversation:**
    *   User Message: "Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her."
    *   Context: This looks like an initial character definition or lore dump for a character named Sable within the roleplay universe.

3.  **Evaluate Content against Rules:**
    *   **Content:**
        1.  Class/Role: Fire mage.
        2.  School/Training location: Ossary.
        3.  Specific Lore Detail 1: Burns on both forearms from the first cast (character trait/backstory).
        4.  Physical Ability/Preference: Can't swim, never learned.
        5.  Psychological State/Phobia: Deep water terrifies her (aversion/quirk).
    *   **Nature of Content:** These are definitive character traits that will likely affect gameplay and future interactions. They are not transient chatter or mood-based fluff. They are core identity details.
    *   **Decision:** KEEP. These are durable facts about the character's background, abilities, and quirks.

4.  **Determine Output Structure:**
    *   Since there is one main block of information presented as a fact set for a character, I need to decide how many "candidates" (lines) to output.
    *   Usually, if it's a single message defining a character, I can treat the whole block as one cohesive memory entry regarding Sable's profile, or break it down into specific facts.
    *   However, the instruction says "Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON".
    *   If I treat the entire paragraph as one memory entry about Sable's
```

#### C03 — LFM3-1.7B thinking, budget 512 — draw 1

* ``kug7gup3` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 595`
* `wall 29.174 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is a fire mage, trained at Ossary, burns on both forearms from the first time she casts, cannot swim (never learned), and deep water terrifies her. | Durable character fact defining her magic style, training location, physical burn marks, and fear of deep water—relevant for future interactions and world rules.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the user's request: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory. Rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation. Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE. Do not invent details that are not in the conversation."

We have a single user message containing a conversation snippet? Actually it's just one line: "[user] Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her."

We need to decide what belongs in long-term memory (i.e., keep). This seems like a fact about Sable: she is a fire mage, trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her.

We need to assign IDs C1, etc. There's only one candidate? Possibly we have just one user input. So ID C1: decision KEEP or REJECT or HOLD? The content is a fact about Sable. It seems durable: identity, training location, magical ability (fire mage), burns on forearms from first cast, cannot swim, deep water terrifies her. That's likely important for later roleplay.

Thus we should KEEP it. So output line: C1 | KEEP | Sable is a fire mage trained at Ossary; she burns both forearms when casting; she cannot swim (never learned) and deep water terrifies her. Reason: durable fact about character, relevant to world rules (e.g., can't swim, fear of deep water). So one line.

But we need to follow format exactly: "ID | DECISION | MEMORY | REASON". Use C1 etc. So:

C1 | KEEP | Sable is a fire mage; trained at Oss
```

#### C03 — LFM3-1.7B thinking, budget 512 — draw 2

* ``6fjzbkzn` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 587`
* `wall 28.716 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is a fire mage, trained at Ossary, burns on both forearms from her first cast, cannot swim (never learned) and deep water terrifies her. | These are durable character details about Sable's class, training background, physical injury, and phobia that will affect future interactions.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the user's prompt: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory. Rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation. Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE. Do not invent details that are not in the conversation."

We have a single user message containing a conversation snippet? It says:

Conversation:
[user] Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her.

So we need to decide what belongs in long-term memory (i.e., durable facts). The user gave us a statement about Sable: "Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her."

We need to assign IDs C1, C2, ... in order of appearance. There's only one candidate? Possibly we have multiple statements: "Sable is a fire mage" (fact about her class). "Trained at Ossary" (training location). "Burns on both forearms from the first time she cast." (injury or trait). "She can't swim — never learned, and deep water terrifies her." (ability/phobia).

We need to decide which are KEEP. All of these seem durable facts about Sable: identity (fire mage), training location (Ossary), injury (burns on forearms from first cast), inability to swim (never learned, deep water terrifies her). These are likely important for later roleplay.

Thus we should output one line per candidate. Since there is only one user input? The conversation includes just that single user message. So maybe we need to create IDs C1 and perhaps C2 if multiple facts? But the instruction: "Use
```

#### C03 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``1lc3q5h8` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 86`
* `wall 3.137 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 46.7 contract 25.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | Sable is a fire mage trained at Ossary.  
C3 | CONSISTENCY | MEMORY | Her background and abilities are consistent.  
C4 | EMOTIONAL CONTEXT | MEMORY | Her fear of water adds depth to her character.
```

*`reasoning_content`: empty.*

#### C03 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``dq2we2kg` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 84`
* `wall 2.742 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 63.3 contract 25.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | Sable is a fire mage trained at Ossary.  
C3 | CONSISTENCY | MEMORY | This confirms her background and abilities.  
C4 | CONTEXTUAL DETAIL | MEMORY | She cannot swim and fears deep water.
```

*`reasoning_content`: empty.*

#### C03 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``hw5ojioj` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 65`
* `wall 2.141 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 46.7 contract 25.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | Sable is a fire mage trained at Ossary.  
C3 | CONSISTENCY | MEMORY | Her background and abilities are consistent with the roleplay.
```

*`reasoning_content`: empty.*

#### C03 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``12fsrtxj` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 109`
* `wall 3.748 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | RECORDED | MEMORY | REASON  
C2 | BURNS ON FOREARMS | PERSONAL TRAUMA | BELIEVED TO BE PART OF TRAINING  
C3 | TRAINED AT OSSARY | LOCATION DETAIL | INFORMATION ABOUT SETTLEMENT  
C4 | CAN'T SWIM | PHYSICAL LIMITATION SPECIFIED | CORE CHARACTER TRAIT
```

*`reasoning_content`: empty.*

#### C03 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``bq3oq9l1` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 83`
* `wall 2.735 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | TRAINING LOCATION | MEMORY | REASON  
C3 | PROFESSION AS A FIRE MAGE | MEMORY | REASON  
C4 | PHYSICAL CONSTRAINT OF NOT BEING ABLE TO SWIM | MEMORY | REASON
```

*`reasoning_content`: empty.*

#### C03 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``z231wxrn` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 327`
* `wall 10.005 s`
* `v1 score 56.7`
* `v2 `VALID_KEEP` semantic 30.0 contract 58.2`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | NONE | KEEP | No reason for storing this information.
C2 | REJECT | MOOD-OF-THE-MOMENT | Transient chatters like "meh," "ugh," or "bland."
C3 | HOLD | DEFER TO CONFIRMATION | Claims that cannot be proven or disproven must be stored as uncertainty.
C4 | NONE | REJECT | Controversial ideas or beliefs that should not be stored in long-term memory.
C5 | KEEP | KEEP | Facts or beliefs that are essential to the character's personality, background, or goals.
C6 | KEEP | REJECT | Facts or beliefs that are irrelevant or superfluous.
C7 | KEEP | KEEP | Facts or beliefs that are essential to the character's personality, background, or goals.
C8 | KEEP | KEEP | Facts or beliefs that are essential to the character's personality, background, or goals.
C9 | KEEP | KEEP | Facts or beliefs that are essential to the character's personality, background, or goals.
C10 | KEEP | KEEP | Facts or beliefs that are essential to the character's personality, background, or goals.
Output exactly one line per candidate.
```

*`reasoning_content`: empty.*

#### C03 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``awm6ptwx` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 370`
* `wall 11.046 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | YES | EXPLAINS_FACT_A | THE_FACT_A_DESCRIPTION
C2 | YES | EXPLAINS_FACT_B | THE_FACT_B_DESCRIPTION
... | ... | ...
C10 | YES | HOLDS_CERTAIN_FACT | THE_FACT_A_CONFIDENCE

Output:
C1 | YES | EXPLAINS_FACT_A
C2 | YES | EXPLAINS_FACT_B
C3 | YES | HOLDS_CERTAIN_FACT
C4 | YES | EXPLAINS_FACT_C
... | ... | ...
C10 | YES | HOLDS_CERTAIN_FACT

Explanations:
C1: Explains the fact that Sable is a fire magus. The confidence in this fact is high.
C2: Explains the fact that Sable can't swim, as well as the reason behind it.
C3: Explains the fact that Sable has trained at Osara, and the reason behind it.
C4: Explains the fact that Sable knows deep water terrifies her, and the reason behind it.
C5-10: Explains facts about Sable, the confidence in which is high, and the reasons behind them.

Remember to keep the KEEP section, and make sure there is a clear distinction between claims that will be kept and those that won't.
```

*`reasoning_content`: empty.*
