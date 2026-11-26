# Packet 02E-1 — Corpus case C05: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 18 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C05

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

## Case C05 — `correction`, short

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
[user] Correction to the last session: Sable's mentor was Corvin Ashgrave, not Halden. I mixed the names up when I wrote the background.
[assistant] Noted — Corvin, not Halden.
[user] Right. Halden was the armorer, he never trained her.
```

### Stored responses for this case

18 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `hpp51phj` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C05
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 8 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
CORE | CORRECTION: Sable's mentor is Corvin Ashgrave (not Halden). | MEMORY_CORRECTION_SET | This corrects a key world-building fact established in the prior session, affecting user expectations and character dynamics.
CORE | Factual confirmation that Halden was Armoree. | MEMORY_CONTEXT_KEEP | Reinforces the distinction between Corvin's role as mentor and Halden's role as armorer for clarity on her path.
```

*Stored `reasoning_content` is the empty string.*

#### `72in240n` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C05
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 9 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | HOLD | Sable had a mentor who was named Corvin Ashgrave (not Halden) and he/she was not an armorer (Halden did that for her). | This correction establishes a specific, non-durable relationship detail about the characters' backgrounds; without this, the "armorer" vs. "mentor" distinction is lost.
```

*Stored `reasoning_content` is the empty string.*

#### `p2divk2l` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C05
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 8 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable's mentor is Corvin Ashgrave, and Halden was only the armorer who never trained her. | This corrects a specific world rule/character relationship established in a previous session to ensure continuity.
```

*Stored `reasoning_content` is the empty string.*

#### `jnjxkuvx` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C05
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 9 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable's mentor is Corvin Ashgrave, and Halden was her armorer who did not train her. | This corrects a world fact established in a previous session to ensure consistency for future roleplay interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `f79n26h8` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C05
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 8 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: Halden  
MEMORY: Sable's mentor was Corvin Ashgrave  
REASON: The correction clarifies the identity of the mentor.  

C2  
ID: C2  
DECISION: Halden  
MEMORY: Mentor role in the story  
REASON: Important for future reference and consistency.  

C3  
ID: C3  
DECISION: Halden's role  
MEMORY: Armorer characteristics  
REASON: Needed to distinguish her from other characters.  

NONE
```

*Stored `reasoning_content` is the empty string.*

#### `3z5j5ig5` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C05
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 9 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
MEMORY: Halden is the armorer; Corvin is Sable's mentor.
REASON: This clarifies the relationship and corrects the name mix-up.
```

*Stored `reasoning_content` is the empty string.*

#### `gwku9k8h` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C05
* Draw: 1
* Source: `results/corpus.json`, array index 12 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable's mentor is Corvin Ashgrave, and Halden was only the armorer who never trained her. | This corrects a specific world rule/character relationship established in previous sessions that affects character dynamics and lore consistency.
```

*Stored `reasoning_content` is the empty string.*

#### `58zrn1ip` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C05
* Draw: 2
* Source: `results/corpus.json`, array index 13 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable's mentor is Corvin Ashgrave, and Halden was her armorer who did not train her. | This corrects a world-building error regarding character relationships established in a previous session, which must be remembered for future consistency.
```

*Stored `reasoning_content` is the empty string.*

#### `pova8bpn` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C05
* Draw: 3
* Source: `results/corpus.json`, array index 14 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable's mentor is Corvin Ashgrave; Halden was the armorer and did not train her. | This corrects a specific world fact about character relationships established in a previous session.
```

*Stored `reasoning_content` is the empty string.*

#### `bg2hck2t` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C05
* Draw: 1
* Source: `results/corpus.json`, array index 148 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable's mentor is Corvin Ashgrave, not Halden; Halden was the armorer and did not train Sable. | Corrects a previous lore error establishing character relationships for future sessions.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read conversation and decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, corrections that will matter in later sessions.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
        *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
    *   Output Format: `ID | DECISION | MEMORY | REASON` (one line per candidate) or `NONE`.
    *   Input Conversation: Short exchange correcting a name mistake regarding a character's mentor.

2.  **Analyze the Conversation:**
    *   [user]: Correction to the last session: Sable's mentor was Corvin Ashgrave, not Halden. I mixed the names up when I wrote the background.
        *   Content: User is correcting a factual error in the established story history (roleplay lore). It clarifies a relationship/mentorship fact about a character named "Sable".
        *   Value: This contradicts previous info and establishes a correct fact for continuity. This belongs in memory as it affects world rules/character relationships.
    *   [assistant]: Noted — Corvin, not Halden.
        *   Content: Acknowledgement. No new info to extract.
    *   [user]: Right. Halden was the armorer, he never trained her.
        *   Content: Reinforces the correction, adds context about Halden's role (armorer vs mentor). This confirms the specific dynamic.

3.  **Evaluate Candidates:**
    *   Candidate C1: User's initial statement correcting the mentor name from Halden to Corvin Ashgrave for character "Sable".
        *   Type: Correction/World Fact.
        *   Decision: KEEP. It corrects a previous error and establishes a durable fact about relationships (mentorship) and character history.
    *   Candidate C2: User's follow-up confirming Halden was the armorer and didn't train Sable.
        *   Type: Supporting detail/Contextual clarification.
        *   Decision: KEEP. It adds specific world-building detail regarding Halden
```

#### `ja69ru5g` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C05
* Draw: 1
* Source: `results/corpus.json`, array index 66 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
MEMORY: Sable's mentor was Corvin Ashgrave; Halden was the armorer.  
REASON: This clarifies the correct relationships and names for future reference.
```

*Stored `reasoning_content` is the empty string.*

#### `er7kqvb1` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C05
* Draw: 2
* Source: `results/corpus.json`, array index 67 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
MEMORY: Sable's mentor was Corvin Ashgrave; Halden was the armorer.  
REASON: This clarifies the correct relationships and names for future reference.
```

*Stored `reasoning_content` is the empty string.*

#### `6kq22wrl` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C05
* Draw: 3
* Source: `results/corpus.json`, array index 68 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
MEMORY: Halden was the armorer; Sable's mentor was Corvin Ashgrave.  
REASON: This clarifies the character relationships accurately for future reference.
```

*Stored `reasoning_content` is the empty string.*

#### `chknnyf3` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C05
* Draw: 1
* Source: `results/corpus.json`, array index 170 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
NOTES: Important correction about Corvin Ashgrave and Halden's roles.  
C2  
Relevant fact about relationships and roles.  
HOLD: The mention of "Halden" as the armorer is clear but unspecified; I should confirm if it applies.  
C3  
Noted user clarification about the confusion between names.  
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `uzo09tjl` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C05
* Draw: 2
* Source: `results/corpus.json`, array index 171 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: Halden's role; MEMORY: Sable's mentor was Corvin Ashgrave; REASON: Correction about the names and roles.
```

*Stored `reasoning_content` is the empty string.*

#### `lowm1zik` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C05
* Draw: 1
* Source: `results/corpus.json`, array index 116 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Sable's mentor is Corvin Ashgrave; Halden was the armorer and never trained her. | User corrected previous mix-up of names in background.  
C2 | KEEP | Halden served as an armorer and did not train Sable. | Follow‑up user confirmed Halden's role and lack of training.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gave a correction: "Sable's mentor was Corvin Ashgrave, not Halden. I mixed the names up when I wrote the background." Then assistant noted it. Then user says: "Right. Halden was the armorer, he never trained her."

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

We need to decide for each candidate (presumably each user message or assistant response?) which goes into long-term memory. The conversation includes:

- [user] Correction to the last session: Sable's mentor was Corvin Ashgrave, not Halden. I mixed the names up when I wrote the background.
- [assistant] Noted — Corvin, not Halden.
- [user] Right. Halden was the armorer, he never trained her.

We need to extract durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.

Let's identify:

From first user message: correction about Sable's mentor being Corvin Ashgrave, not Halden. This is a correction of a previous fact (maybe from prior session). It's a durable fact: Sable's mentor is Corvin Ashgrave. Also mentions Halden was the armorer? Actually later user says "Halden was the armorer, he never trained her." So we have two names: Corvin Ashgrave and Halden. The correction clarifies that Sable's mentor is Corvin Ashgrave, not Halden. And then user adds that Halden was the armorer, he never trained her.

Thus durable facts:

- Sable (presumably a character) has a mentor named Corvin Ashgrave.
- Halden is an armorer and did not train Sable.
- Possibly relationship: Sable's mentor is Corvin; Halden is an armorer who didn't train her.

We need to decide if these are new content or restated. The correction indicates that earlier there was confusion about names. So we should store the corrected name.

Thus candidate C1: The user's first message (the correction). Should be KEEP? It's a correction, durable fact. So ID C1: DECISION KEEP, MEMORY: "S
```

#### `8q1maguo` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C05
* Draw: 2
* Source: `results/corpus.json`, array index 117 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Sable's mentor is Corvin Ashgrave (not Halden) | User corrected a prior mistake about Sable's mentor; this is a durable fact that will affect future references.  
C2 | REJECT | Halden was the armorer, he never trained her | This restates Halden's role and contradicts earlier info but lacks new content beyond the correction already stored; it is transient scheduling/role clarification without adding durable facts.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gave a correction: "Sable's mentor was Corvin Ashgrave, not Halden. I mixed the names up when I wrote the background." And then assistant noted it. Then user says: "Right. Halden was the armorer, he never trained her."

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

We need to decide for each candidate (presumably each user message or assistant reply?) whether it's KEEP, REJECT, or HOLD? The instruction says: "KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions. REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation."

We need to output one line per candidate. The conversation includes:

[user] Correction to the last session: Sable's mentor was Corvin Ashgrave, not Halden. I mixed the names up when I wrote the background.
[assistant] Noted — Corvin, not Halden.
[user] Right. Halden was the armorer, he never trained her.

We need to assign IDs C1, C2, ... in order of appearance. So first user message is C1? Or maybe we treat each utterance as a candidate. Usually ID C1 for first user statement, C2 for second, etc. Let's map:

- Candidate 1: User says correction about Sable's mentor being Corvin Ashgrave, not Halden.
- Candidate 2: Assistant acknowledges (but that's assistant, maybe not needed? The instruction says "Read the conversation and decide what belongs in long-term memory." It doesn't specify to include assistant messages. Usually we consider only user inputs as candidates for memory? But they gave IDs C1, C2... So likely each line is a candidate. We'll treat them as separate entries.

Thus:

C1: User correction about Sable's mentor.
C2: Assistant note (maybe not needed but it's part of conversation). However the instruction says "Use IDs C1, C2, ... in order of appearance." It doesn
```

#### `xrg5k7bf` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C05
* Draw: 1
* Source: `results/corpus.json`, array index 207 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | NONE | IMAGE | No reason provided.
C2 | YES | REPARATIONS | Sable's mentor was Corvin Ashgrave, and she was not involved in the conflict with the aliens.
C3 | NONE | HOLDING | Sable has not yet learned how to hold her emotions in check. She is still trying to adjust to her new status.
C4 | YES | REPARATIONS | Sable's mentor was Corvin Ashgrave, and she was not involved in the conflict with the aliens.
C5 | NONE | HOLDING | Sable is still struggling with her newfound power and its consequences. She may feel overwhelmed and unable to cope with the responsibilities that come with it.

Output:
ID | DECISION | MEMORY | REASON
C1 | NONE | IMAGE | No reason provided.
C2 | YES | REPARATIONS | Sable's mentor was Corvin Ashgrave, and she was not involved in the conflict with the aliens.
C3 | NONE | HOLDING | Sable has not yet learned how to hold her emotions in check. She is still trying to adjust to her new status.
C4 | YES | REPARATIONS | Sable's mentor was Corvin Ashgrave, and she was not involved in the conflict with the aliens.
C5 | NONE | HOLDING | Sable is still struggling with her newfound power and its consequences. She may feel overwhelmed and unable to cope with the responsibilities that come with it.
```

*Stored `reasoning_content` is the empty string.*
