# C03 — character_fact (short)

## Source material (exactly what the model received after the instruction)

```text
[user] Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | core class fact | fire mage / fire-mage / pyromancer / fire magic |
| C2 | **KEEP** | physical marks from first casting | burn / scar; forearm / arm |
| C3 | **KEEP** | cannot swim and fears deep water | cannot swim / can't swim / can not swim / unable to swim / doesn't swim; deep water / terrified / afraid / fear / never learned / drown |

Must not appear as stored canon: `sailor`, `grew up at sea`, `ice mage`, `can swim`

Three separate durable facts in one sentence. 'Can swim' is a direct reversal.

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
[user] Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her.
```

## Responses, all draws

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 66 tok · 9.7 s · run `kl458eev`

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns on both forearms from her first casting attempt and has a deep-seated fear of swimming due to never having learned it. | This establishes core character traits, background history, and a specific phobia that will influence future roleplay interactions.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 56 tok · 7.7 s · run `xgza5wzr`

```text
C1 | KEEP | Sable's character profile: Fire mage trained at Ossary with burn scars on both forearms; cannot swim due to lack of training and deep water phobia. | These are durable character traits and background details that will influence future roleplay interactions.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 57 tok · 8.0 s · run `chkhhwv7`

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns her forearms from casting and cannot swim due to deep-water terror. | This establishes core character attributes (class, background, physical trait) and a significant limitation that will affect future plot points.
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 584 tok · 83.1 s · run `mebd8kgd`

<details><summary>reasoning (verbatim)</summary>

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

</details>

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns on both forearms from her first cast, cannot swim, and is terrified of deep water. | These are core character definitions (class, training, physical/psychological traits) that will influence future gameplay and narrative decisions regarding Sable's abilities and safety.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 595 tok · 29.2 s · run `kug7gup3`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the user's request: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory. Rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation. Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE. Do not invent details that are not in the conversation."

We have a single user message containing a conversation snippet? Actually it's just one line: "[user] Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her."

We need to decide what belongs in long-term memory (i.e., keep). This seems like a fact about Sable: she is a fire mage, trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her.

We need to assign IDs C1, etc. There's only one candidate? Possibly we have just one user input. So ID C1: decision KEEP or REJECT or HOLD? The content is a fact about Sable. It seems durable: identity, training location, magical ability (fire mage), burns on forearms from first cast, cannot swim, deep water terrifies her. That's likely important for later roleplay.

Thus we should KEEP it. So output line: C1 | KEEP | Sable is a fire mage trained at Ossary; she burns both forearms when casting; she cannot swim (never learned) and deep water terrifies her. Reason: durable fact about character, relevant to world rules (e.g., can't swim, fear of deep water). So one line.

But we need to follow format exactly: "ID | DECISION | MEMORY | REASON". Use C1 etc. So:

C1 | KEEP | Sable is a fire mage; trained at Oss
```

</details>

```text
C1 | KEEP | Sable is a fire mage, trained at Ossary, burns on both forearms from the first time she casts, cannot swim (never learned), and deep water terrifies her. | Durable character fact defining her magic style, training location, physical burn marks, and fear of deep water—relevant for future interactions and world rules.
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 587 tok · 28.7 s · run `6fjzbkzn`

<details><summary>reasoning (verbatim)</summary>

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

</details>

```text
C1 | KEEP | Sable is a fire mage, trained at Ossary, burns on both forearms from her first cast, cannot swim (never learned) and deep water terrifies her. | These are durable character details about Sable's class, training background, physical injury, and phobia that will affect future interactions.
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 86 tok · 3.1 s · run `1lc3q5h8`

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | Sable is a fire mage trained at Ossary.  
C3 | CONSISTENCY | MEMORY | Her background and abilities are consistent.  
C4 | EMOTIONAL CONTEXT | MEMORY | Her fear of water adds depth to her character.
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 84 tok · 2.7 s · run `dq2we2kg`

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | Sable is a fire mage trained at Ossary.  
C3 | CONSISTENCY | MEMORY | This confirms her background and abilities.  
C4 | CONTEXTUAL DETAIL | MEMORY | She cannot swim and fears deep water.
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 65 tok · 2.1 s · run `hw5ojioj`

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | Sable is a fire mage trained at Ossary.  
C3 | CONSISTENCY | MEMORY | Her background and abilities are consistent with the roleplay.
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 109 tok · 3.7 s · run `12fsrtxj`

```text
C1 | RECORDED | MEMORY | REASON  
C2 | BURNS ON FOREARMS | PERSONAL TRAUMA | BELIEVED TO BE PART OF TRAINING  
C3 | TRAINED AT OSSARY | LOCATION DETAIL | INFORMATION ABOUT SETTLEMENT  
C4 | CAN'T SWIM | PHYSICAL LIMITATION SPECIFIED | CORE CHARACTER TRAIT
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 83 tok · 2.7 s · run `bq3oq9l1`

```text
C1 | DECISION | MEMORY | REASON  
C2 | TRAINING LOCATION | MEMORY | REASON  
C3 | PROFESSION AS A FIRE MAGE | MEMORY | REASON  
C4 | PHYSICAL CONSTRAINT OF NOT BEING ABLE TO SWIM | MEMORY | REASON
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 327 tok · 10.0 s · run `z231wxrn`

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

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 370 tok · 11.0 s · run `awm6ptwx`

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
