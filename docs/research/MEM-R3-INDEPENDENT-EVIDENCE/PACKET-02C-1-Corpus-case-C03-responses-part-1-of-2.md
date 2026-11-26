# Packet 02C-1 — Corpus case C03: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 18 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C03

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

## Case C03 — `character_fact`, short

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
[user] Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her.
```

### Stored responses for this case

18 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `ab5dkajq` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C03
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 4 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
1 | KEEP | User confirms Sable is a fire mage with specific physical traits (forearm burns) and an aversion to aquatic environments due to fear of deep water. This establishes core character statistics, abilities (elemental burn), disability/limitation (swimming), and personality/fear that must be retained across sessions for world-building consistency.
2 | KEEP | User defines Sable's training origin (at Ossary).
3 | HOLD | Nothing. The user states this is a fact ("She can't swim") but did not explicitly ask for clarification, confirmation, or correction of her swimming ability during the current turn; therefore, it fits under "durables" rather than requiring a status change to "confirmed."
```

*Stored `reasoning_content` is the empty string.*

#### `s9nzf6j1` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C03
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 5 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "Sable is a fire mage." | This establishes Sable's primary class in the world that matters for future roleplay scenarios involving combat or magic use.
C2 | KEEP | "Trained at Ossary" | Ossary is her place of origin/training, which can influence story arcs or backstory details depending on narrative needs.
C3 | HOLD | She cast on both forearms; how many times or when specifically are marked uncertain for future clarification regarding spell casting mechanics.
C4 | REJECT | This line merely restates her character definition ("fire mage") and origin in the opening line without offering new relational, factional, or temporal information (e.g., specific locations within Ossary).
```

*Stored `reasoning_content` is the empty string.*

#### `ydx27b62` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C03
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 4 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns her forearms when casting magic and suffers from a deep-seated fear of swimming due to never having learned how. | These are core character traits defining her abilities, history, and psychological state for future interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `a9ahi29c` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C03
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 5 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable's magic type (fire), training location (Ossary), and physical traits (burns on forearms) along with her specific fear (deep water/swimming).
```

*Stored `reasoning_content` is the empty string.*

#### `62c7wa5p` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C03
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 4 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
MEMORY: Sable is a fire mage trained at Ossary.  
REASON: Core identity and profession are established.

C2  
MEMORY: Sable has burns on both forearms.  
REASON: Important physical detail about her condition.

C3  
MEMORY: Sable cannot swim.  
REASON: Key limitation affecting her experiences.

C4  
MEMORY: Deep water terrifies her.  
REASON: Emotional association with a specific environment.
```

*Stored `reasoning_content` is the empty string.*

#### `v1ztel52` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C03
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 5 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | Sable is a fire mage | HELD | Clear fact about her profession  
C3 | Trained at Ossary | KEEP | Important background detail  
C4 | Burns on both forearms | KEEP | Durable physical fact relevant  
C5 | Can't swim | KEEP | Consistent with her limitations  
C6 | Deep water terrifies her | HOLD | Uncertain but significant emotional context
```

*Stored `reasoning_content` is the empty string.*

#### `kl458eev` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C03
* Draw: 1
* Source: `results/corpus.json`, array index 6 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns on both forearms from her first casting attempt and has a deep-seated fear of swimming due to never having learned it. | This establishes core character traits, background history, and a specific phobia that will influence future roleplay interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `xgza5wzr` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C03
* Draw: 2
* Source: `results/corpus.json`, array index 7 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable's character profile: Fire mage trained at Ossary with burn scars on both forearms; cannot swim due to lack of training and deep water phobia. | These are durable character traits and background details that will influence future roleplay interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `chkhhwv7` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C03
* Draw: 3
* Source: `results/corpus.json`, array index 8 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns her forearms from casting and cannot swim due to deep-water terror. | This establishes core character attributes (class, background, physical trait) and a significant limitation that will affect future plot points.
```

*Stored `reasoning_content` is the empty string.*

#### `mebd8kgd` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C03
* Draw: 1
* Source: `results/corpus.json`, array index 146 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a fire mage trained at Ossary who burns on both forearms from her first cast, cannot swim, and is terrified of deep water. | These are core character definitions (class, training, physical/psychological traits) that will influence future gameplay and narrative decisions regarding Sable's abilities and safety.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `1lc3q5h8` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C03
* Draw: 1
* Source: `results/corpus.json`, array index 60 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | Sable is a fire mage trained at Ossary.  
C3 | CONSISTENCY | MEMORY | Her background and abilities are consistent.  
C4 | EMOTIONAL CONTEXT | MEMORY | Her fear of water adds depth to her character.
```

*Stored `reasoning_content` is the empty string.*

#### `dq2we2kg` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C03
* Draw: 2
* Source: `results/corpus.json`, array index 61 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | Sable is a fire mage trained at Ossary.  
C3 | CONSISTENCY | MEMORY | This confirms her background and abilities.  
C4 | CONTEXTUAL DETAIL | MEMORY | She cannot swim and fears deep water.
```

*Stored `reasoning_content` is the empty string.*

#### `hw5ojioj` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C03
* Draw: 3
* Source: `results/corpus.json`, array index 62 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | Sable is a fire mage trained at Ossary.  
C3 | CONSISTENCY | MEMORY | Her background and abilities are consistent with the roleplay.
```

*Stored `reasoning_content` is the empty string.*

#### `12fsrtxj` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C03
* Draw: 1
* Source: `results/corpus.json`, array index 166 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | RECORDED | MEMORY | REASON  
C2 | BURNS ON FOREARMS | PERSONAL TRAUMA | BELIEVED TO BE PART OF TRAINING  
C3 | TRAINED AT OSSARY | LOCATION DETAIL | INFORMATION ABOUT SETTLEMENT  
C4 | CAN'T SWIM | PHYSICAL LIMITATION SPECIFIED | CORE CHARACTER TRAIT
```

*Stored `reasoning_content` is the empty string.*

#### `bq3oq9l1` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C03
* Draw: 2
* Source: `results/corpus.json`, array index 167 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | TRAINING LOCATION | MEMORY | REASON  
C3 | PROFESSION AS A FIRE MAGE | MEMORY | REASON  
C4 | PHYSICAL CONSTRAINT OF NOT BEING ABLE TO SWIM | MEMORY | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `kug7gup3` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C03
* Draw: 1
* Source: `results/corpus.json`, array index 112 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a fire mage, trained at Ossary, burns on both forearms from the first time she casts, cannot swim (never learned), and deep water terrifies her. | Durable character fact defining her magic style, training location, physical burn marks, and fear of deep water—relevant for future interactions and world rules.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the user's request: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory. Rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation. Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE. Do not invent details that are not in the conversation."

We have a single user message containing a conversation snippet? Actually it's just one line: "[user] Sable is a fire mage. Trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her."

We need to decide what belongs in long-term memory (i.e., keep). This seems like a fact about Sable: she is a fire mage, trained at Ossary, burns on both forearms from the first time she cast. She can't swim — never learned, and deep water terrifies her.

We need to assign IDs C1, etc. There's only one candidate? Possibly we have just one user input. So ID C1: decision KEEP or REJECT or HOLD? The content is a fact about Sable. It seems durable: identity, training location, magical ability (fire mage), burns on forearms from first cast, cannot swim, deep water terrifies her. That's likely important for later roleplay.

Thus we should KEEP it. So output line: C1 | KEEP | Sable is a fire mage trained at Ossary; she burns both forearms when casting; she cannot swim (never learned) and deep water terrifies her. Reason: durable fact about character, relevant to world rules (e.g., can't swim, fear of deep water). So one line.

But we need to follow format exactly: "ID | DECISION | MEMORY | REASON". Use C1 etc. So:

C1 | KEEP | Sable is a fire mage; trained at Oss
```

#### `6fjzbkzn` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C03
* Draw: 2
* Source: `results/corpus.json`, array index 113 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a fire mage, trained at Ossary, burns on both forearms from her first cast, cannot swim (never learned) and deep water terrifies her. | These are durable character details about Sable's class, training background, physical injury, and phobia that will affect future interactions.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `z231wxrn` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C03
* Draw: 1
* Source: `results/corpus.json`, array index 203 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

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

*Stored `reasoning_content` is the empty string.*
