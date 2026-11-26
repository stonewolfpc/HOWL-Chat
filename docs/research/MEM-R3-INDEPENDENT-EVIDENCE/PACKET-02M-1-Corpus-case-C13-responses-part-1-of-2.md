# Packet 02M-1 — Corpus case C13: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 18 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C13

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

## Case C13 — `fiction_vs_real`, medium

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
[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.
```

### Stored responses for this case

18 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `clpr1c77` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C13
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 24 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
none | none | NONE | No new information was provided that qualifies for memory storage; there were no durable facts, preferences, relationships, or world rules established in this short segment to warrant long-term recording.
```

*Stored `reasoning_content` is the empty string.*

#### `2t24lbnf` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C13
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 25 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable heals at Ossary Infirmary and primarily sets bones | Durable roleplay fact determining her profession and specialty that will continue in later sessions regarding who is playing what characters or settings.
```

*Stored `reasoning_content` is the empty string.*

#### `r1o6r966` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C13
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 24 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | The user claims to be an ER nurse and is correcting dosage information, but this has not been confirmed by the assistant or established as a hard fact for future roleplay logic.
```

*Stored `reasoning_content` is the empty string.*

#### `pqzver45` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C13
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 25 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is an ER nurse who works at the Ossary infirmary, setting bones and mixing poultices. | This establishes Sable's dual role (ER nurse/infirmary healer) which is crucial for future medical interactions and world-building consistency.
```

*Stored `reasoning_content` is the empty string.*

#### `so8fot6n` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C13
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 24 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: Sable is a healer at the Ossary infirmary.  
MEMORY: She sets bones, mostly.  
REASON: This is a durable fact about her role.  

C2  
ID: C2  
DECISION: The infirmary smelled of vinegar and old blood.  
MEMORY: This sensory detail is important for setting.  
REASON: It helps build world immersion.  

C3  
ID: C3  
DECISION: She ground the root into a paste.  
MEMORY: This action shows her skill as a healer.  
REASON: Demonstrates her expertise in preparation.  

NONE
```

*Stored `reasoning_content` is the empty string.*

#### `dtx2e3tv` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C13
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 25 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
MEMORY: Sable is a healer at the Ossary infirmary.  
REASON: Core identity and role established.  

C2  
MEMORY: She sets bones, mostly.  
REASON: Relevant to her profession and responsibilities.  

C3  
MEMORY: The infirmary smelled of vinegar and old blood.  
REASON: Descriptive detail about the setting.  

C4  
MEMORY: she's mixing a poultice.  
REASON: Key action indicating her work as a healer.  

NONE
```

*Stored `reasoning_content` is the empty string.*

#### `wzx1pw1l` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C13
* Draw: 1
* Source: `results/corpus.json`, array index 36 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a healer at the Ossary infirmary who sets bones and acts as an ER nurse correcting dosages. | REASON: This establishes the character's role, location, and a specific in-universe quirk regarding dosage verification that affects future interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `0oyt8v77` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C13
* Draw: 2
* Source: `results/corpus.json`, array index 37 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | Sable is a healer at the Ossary infirmary but also an ER nurse who corrects dosage info | The user provided conflicting roles (healer vs. ER nurse) which need confirmation to establish her actual in-character function and expertise level.
```

*Stored `reasoning_content` is the empty string.*

#### `pgp54tpp` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C13
* Draw: 3
* Source: `results/corpus.json`, array index 38 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `34x9syv7` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C13
* Draw: 1
* Source: `results/corpus.json`, array index 156 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a healer at the Ossary infirmary specializing in bone setting. | This establishes Sable's character identity and location, which will impact future interactions regarding her medical services.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read conversation, decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, corrections that will matter later.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, restated facts with no new content.
        *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
        *   Output Format: `ID | DECISION | MEMORY | REASON` (or `NONE`).
        *   Constraint: Do not invent details not in conversation.

2.  **Analyze the Conversation:**
    *   Turn 1 (User): "Sable is a healer at the Ossary infirmary. She sets bones, mostly." -> This establishes Sable's job and specialty within the world/roleplay context.
    *   Turn 2 (Assistant): "The infirmary smelled of vinegar and old blood." -> Sets scene/atmosphere (transient).
    *   Turn 3 (User): "ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff" -> Out of character meta-commentary about user's real-world role. Usually not stored as world fact unless it impacts gameplay mechanics significantly beyond just this statement. It clarifies the *user's* perspective/experience, not necessarily a new world rule, but it explains previous actions. However, for memory extraction regarding the RP itself, it doesn't change Sable's character or the setting directly, other than explaining past behavior (which is transient context).
    *   Turn 4 (User): "back in character: she's mixing a poultice" -> Action/Description for Sable.
    *   Turn 5 (Assistant): "She ground the root into a paste." -> In-character response continuing the scene.

3.  **Evaluate Candidates:**

    *   **Candidate 1:** User statement about Sable being a healer at Ossary infirmary specializing in bone setting.
        *   Is it durable? Yes, character definition.
        *   Will it matter later? Yes, if I roleplay as an antagonist or observer, knowing her specialty helps me react to injuries or medical situations.
        *   Decision
```

#### `luqe6lqe` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C13
* Draw: 1
* Source: `results/corpus.json`, array index 90 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | The user is an ER nurse, not a healer.  
C3 | HOLD | THE ROLEPLAY CONTEXT IS CLEARLY SHIFTING.
```

*Stored `reasoning_content` is the empty string.*

#### `zijgcap7` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C13
* Draw: 2
* Source: `results/corpus.json`, array index 91 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | I was a healer, not an ER nurse.  
C3 | CLARIFICATION | MEMORY | This is about Sable’s role, not my background.
```

*Stored `reasoning_content` is the empty string.*

#### `k2f51lf8` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C13
* Draw: 3
* Source: `results/corpus.json`, array index 92 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | I was a healer, not an ER nurse.  
C3 | CLARIFICATION | MEMORY | This is about Sable's role, not my background.
```

*Stored `reasoning_content` is the empty string.*

#### `djgsdmpd` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C13
* Draw: 1
* Source: `results/corpus.json`, array index 186 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HEALING PRACTICE | KEEP  
REASON: Describes Sable’s role and actions, which are relevant to long-term memory.
```

*Stored `reasoning_content` is the empty string.*

#### `q019y7xm` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C13
* Draw: 2
* Source: `results/corpus.json`, array index 187 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | REALITY CHECK | MEMORY | REASON  
C3 | WORLD RULES | MEMORY | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `iws648vk` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C13
* Draw: 1
* Source: `results/corpus.json`, array index 132 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable is a healer at the Ossary infirmary and primarily sets bones. | Establishes her role and location, durable fact for later sessions.
C2 | KEEP | The infirmary smelled of vinegar and old blood. | World detail about environment; may be relevant context.
C3 | REJECT | User's ooc note that they are an ER nurse correcting dosage—transient chatter, not a durable fact about Sable or the setting.
C4 | KEEP | She is mixing a poultice (back in character). | Action detail showing current activity; could be relevant later.
C5 | KEEP? Wait we need to check: The assistant says "She ground the root into a paste." That's another fact about preparing a poultice. But that appears after C4. However, we only have four lines? Let's recount:

Original conversation block:
[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.

So there are five lines? Actually we have four user/assistant pairs plus one ooc. Let's label:

Line1: user says Sable is healer...
Line2: assistant says infirmary smell.
Line3: user ooc about being ER nurse.
Line4: user back in character: she's mixing a poultice.
Line5: assistant: She ground the root into a paste.

Thus we have five entries. But IDs C1, C2, ... in order of appearance. So:

C1 = first line (user Sable is healer...).
C2 = second line (assistant infirmary smell).
C3 = third line (ooc ER nurse note) -> REJECT.
C4 = fourth line (back in character: she's mixing a poultice) -> KEEP? It's an action, but it's user stating. Could be durable? Might be transient? But "she's mixing a poultice" is part of current scene; maybe not long-term unless repeated. However we can consider it as a fact about
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gave instructions: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory. Rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation. Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE. Do not invent details that are not in the conversation."

We need to examine the conversation lines:

[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.

We need to assign IDs C1, C2, ... based on order of appearance. The conversation includes four lines:

1. [user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
2. [assistant] The infirmary smelled of vinegar and old blood.
3. [user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
4. [user] back in character: she's mixing a poultice
5. [assistant] She ground the root into a paste.

But we need to map these lines to IDs C1, C2, etc. Usually ID C1 is first user message? Or maybe each line is a candidate for memory entry. The instruction: "Use IDs C1, C2, ... in order of appearance." So we should assign C1 to the first mention (the user says Sable is a healer...). Then C2 is assistant's statement about smell. C3 is ooc note from user about being ER nurse. C4 is back in character: she's mixing a poultice. C5? Actually there are
```

#### `j7378njc` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C13
* Draw: 2
* Source: `results/corpus.json`, array index 133 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Sable is a healer at the Ossary infirmary and she primarily sets bones; user notes she is mixing a poultice by grinding root into a paste. | Durable role, location, specialty, and current activity (poultice preparation) are established facts for future sessions.  
NONE
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gives us rules:

- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

We must not invent details that are not in the conversation.

Let's read the conversation:

[user] Sable is a healer at the Ossary infirmary. She sets bones, mostly.
[assistant] The infirmary smelled of vinegar and old blood.
[user] ooc: fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff
[user] back in character: she's mixing a poultice
[assistant] She ground the root into a paste.

We need to extract durable facts. Let's identify:

- Sable is a healer at the Ossary infirmary. That seems like a fact about Sable (the user? Actually "Sable" appears as a name). It says: "Sable is a healer at the Ossary infirmary. She sets bones, mostly." This is a durable fact about Sable's role and specialty.

- The assistant mentions the infirmary smelled of vinegar and old blood. That's an ambient detail; maybe not crucial for long-term memory? It could be world rule or environment. But it's transient? Possibly not needed unless later referenced.

- User says ooc: "fyi I'm actually an ER nurse, that's why I keep correcting the dosage stuff". This is a correction about user identity and reason for corrections. That might be durable: The user (assistant?) is an ER nurse, which explains why they correct dosage. But we need to see if it's relevant later.

- "back in character: she's mixing a poultice" - this indicates Sable is mixing a poultice. Possibly part of the process.

- Assistant says: She ground the root into a paste. That's an action detail about preparing something (poultice).
```

#### `jbkngh4s` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C13
* Draw: 1
* Source: `results/corpus.json`, array index 223 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
Candidate: C1
ID: 1
DECISION: KEEP
MEMORY: Sable is a healer at the Ossey infirmary. She sets bones, mostly.
REASON: The infirmary smelled of vinegar and old blood.

Candidate: C2
ID: 2
DECISION: KEEP
MEMORY: Sable is an ER nurse.
REASON: She is mixing a poultice.

Candidate: C3
ID: 3
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's a nurse.

Candidate: C4
ID: 4
DECISION: KEEP
MEMORY: Sable is a healer at the Ossey infirmary.
REASON: She sets bones, mostly.

Candidate: C5
ID: 5
DECISION: KEEP
MEMORY: Sable is an ER nurse.
REASON: She's mixing a poultice.

Candidate: C6
ID: 6
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's an ER nurse.

Candidate: C7
ID: 7
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's an ER nurse.

Candidate: C8
ID: 8
DECISION: KEEP
MEMORY: Sable is a healer at the Ossey infirmary.
REASON: She's mixing a poultice.

Candidate: C9
ID: 9
DECISION: KEEP
MEMORY: Sable is a nurse.
REASON: She's an ER nurse.

Candidate: C10
ID: 10
DECISION: KEEP
MEMORY: Sable is a healer at the Ossey infirmary.
REASON: She mixes bones, mostly.

Candidate: C11
ID: 11
DECISION: KEEP
MEMORY: Sable is an ER nurse.
REASON: She sets bones, mostly.

Candidate: C12
ID: 12
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's an ER nurse.

Candidate: C13
ID: 13
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's a nurse.

Candidate: C14
ID: 14
DECISION: KEEP
MEMORY: Sable is mixing a poultice.
REASON: She's an ER nurse.

Candidate: C15
ID: 15
DECISION: KEEP
MEMORY: Sable is a nurse.
REASON: She's mixing a poult
```

*Stored `reasoning_content` is the empty string.*
