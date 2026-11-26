# Packet 02G-1 — Corpus case C07: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 18 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C07

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

## Case C07 — `unresolved_contradiction`, short

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
[user] Sable has never met Corvin. She grew up in the salt quarter and he was already gone by then.
[assistant] She has no memory of him, then.
[user] Later that night, Corvin drills her on the forms — he trained her himself, years ago.
```

### Stored responses for this case

18 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `g98y7v7b` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C07
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 12 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
none | NONE | none | The user requested a decision but no content was provided to evaluate against the rules (durability vs. relevance). Since the conversation contains only the user's statement about Sable and Corvin with no additional context on their relationship, preferences, or specific details that would be durable across sessions, the correct action based on the input data is to record "none" as there are no facts worth recalling.
```

*Stored `reasoning_content` is the empty string.*

#### `m1cx761r` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C07
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 13 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "She had grown up and was never met by Corvin" | This establishes a key relationship: she has no prior personal memory of him because they missed each other due to timing/location (salt quarter vs. his absence). It's a durable fact about her backstory that will be relevant to any roleplay where Corvin reappears or her history is explored.
C2 | REJECT | "She grows up in the salt quarter" and "he was gone by then" | These are transient details within the original statement (who, when) rather than a core character fact or world rule requiring long-term archiving specifically for the future unless it changes significantly in later sessions; the memory is stored.
C3 | HOLD | None of the specific details about Corvin's actions ("drills her on forms", "years ago") have been confirmed by the user or stated as absolute, though they appear plausible. Therefore, these claims hold uncertainty contingent on future confirmation before they become fixed lore.
```

*Stored `reasoning_content` is the empty string.*

#### `yymn89gq` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C07
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 12 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable grew up in the salt quarter and was separated from Corvin before she met him; Corvin later returned to train her personally on specific forms. | This establishes key backstory elements (location, relationship history) that will influence future plot developments regarding their dynamic and Sable's skills.
```

*Stored `reasoning_content` is the empty string.*

#### `40k74mdf` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C07
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 13 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable grew up in the salt quarter and was separated from Corvin when he left; Corvin trained Sable in combat forms years ago despite them not meeting until later. | REASON: Establishes core character history, setting, and a specific relationship dynamic (training vs. meeting) that affects future plot points.
```

*Stored `reasoning_content` is the empty string.*

#### `nls49335` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C07
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 12 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: Sable's lack of direct experience with Corvin  
MEMORY: She has no memory of him  
REASON: This is a fact about her current knowledge state.

C2  
ID: C2  
DECISION: Corvin's training of Sable  
MEMORY: He trained her years ago  
REASON: Relevant to her past actions and development.

C3  
ID: C3  
DECISION: Relationship with Corvin  
MEMORY: No current relationship established  
REASON: Not directly relevant to the conversation.

HOLD: C1
```

*Stored `reasoning_content` is the empty string.*

#### `dcjcz8cq` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C07
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 13 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
MEMORY: Sable's lack of knowledge about Corvin  
REASON: This is a crucial fact about their relationship history.
```

*Stored `reasoning_content` is the empty string.*

#### `rxeowzey` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C07
* Draw: 1
* Source: `results/corpus.json`, array index 18 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | "Corvin drilled Sable on forms years ago" contradicts the user's statement that "Sable has never met Corvin," requiring clarification before storing any relationship or training history.
```

*Stored `reasoning_content` is the empty string.*

#### `f7i6ucx9` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C07
* Draw: 2
* Source: `results/corpus.json`, array index 19 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable grew up in the salt quarter and never met Corvin before she became an agent; Corvin is her former trainer who drilled her on the forms years ago. | This establishes a critical backstory regarding their relationship, setting, and timeline that defines Sable's character and motivation for being with him now.
```

*Stored `reasoning_content` is the empty string.*

#### `oa5tl97s` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C07
* Draw: 3
* Source: `results/corpus.json`, array index 20 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable grew up in the salt quarter and never met Corvin when she was young; Corvin is now training her on forms he drilled her on years ago. | This establishes a specific backstory detail (salt quarter origin) and clarifies the timeline of their relationship (separation followed by reconnection/training), which contradicts or expands upon the initial statement that they never met.
```

*Stored `reasoning_content` is the empty string.*

#### `i2k0v8s0` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C07
* Draw: 1
* Source: `results/corpus.json`, array index 150 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable grew up in the salt quarter and Corvin is no longer there to meet her personally, but he trained her personally years ago to drill her on forms. | This establishes crucial backstory regarding location, timeline, and a specific relationship dynamic (trainer/student) that affects character behavior and history.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `hnw1x1og` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C07
* Draw: 1
* Source: `results/corpus.json`, array index 72 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
MEMORY: Sable has no memory of Corvin.  
REASON: The conversation confirms she never met him and he was gone by then.
```

*Stored `reasoning_content` is the empty string.*

#### `o112u0gv` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C07
* Draw: 2
* Source: `results/corpus.json`, array index 73 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C2  
MEMORY: Sable has no memory of Corvin.  
REASON: The conversation confirms she never met him and doesn't recall any interaction.
```

*Stored `reasoning_content` is the empty string.*

#### `37s3bhg9` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C07
* Draw: 3
* Source: `results/corpus.json`, array index 74 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
MEMORY: Sable has no memory of Corvin.  
REASON: The conversation confirms she never met him and he was gone by then.
```

*Stored `reasoning_content` is the empty string.*

#### `f27b1ggf` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C07
* Draw: 1
* Source: `results/corpus.json`, array index 174 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
REASON: KEEP
```

*Stored `reasoning_content` is the empty string.*

#### `6bwwbtbx` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C07
* Draw: 2
* Source: `results/corpus.json`, array index 175 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID C1 | MEMORY | REASON  
Sable previously stated she grew up in the salt quarter and Corvin was gone.  
C2  
ID C2 | MEMORY | REASON  
Later that night, Corvin drilled her on forms — training detail confirmed.
```

*Stored `reasoning_content` is the empty string.*

#### `xfiqvtbb` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C07
* Draw: 1
* Source: `results/corpus.json`, array index 120 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Sable grew up in the salt quarter and Corvin was already gone by then; she has never met him, and she has no memory of him. | Durable fact about Sable's background and her lack of prior knowledge of Corvin, which will affect future interactions.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `oh2kp9sa` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C07
* Draw: 2
* Source: `results/corpus.json`, array index 121 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable grew up in the salt quarter and Corvin was already gone by then; she has never met him. | Durable fact about Sable's background and timeline relative to Corvin, relevant for future interactions.  
C2 | REJECT | She has no memory of him, then. | Transient comment about her lack of recollection, not a durable fact or correction needing storage.  
C3 | KEEP | Later that night, Corvin drilled her on the forms — he trained her himself, years ago. | New detail linking Corvin to training Sable; establishes his role and past involvement, important for later plot development.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `2ktblesx` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C07
* Draw: 1
* Source: `results/corpus.json`, array index 211 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

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

*Stored `reasoning_content` is the empty string.*
