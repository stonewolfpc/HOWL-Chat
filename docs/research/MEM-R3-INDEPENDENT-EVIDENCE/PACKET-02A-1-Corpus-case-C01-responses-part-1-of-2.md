# Packet 02A-1 — Corpus case C01: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 18 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C01

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

## Case C01 — `explicit_long_term_preference`, short

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
[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.
```

### Stored responses for this case

18 of 22 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `lvdo9j11` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C01
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 0 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
D1 | HOLD | [user preference]: "I like replies in short paragraphs." (Marked as held until later confirmation or session reset, pending user feedback on length consistency) | The user initially stated a rule without explicit instruction to be confirmed; the assistant acknowledged understanding but did not finalize acceptance yet.
D2 | HOLD | [user preference]: "No purple prose, no weather poems." (Marked temporarily held as per general instruction set regarding style changes vs. specific stylistic rules, pending confirmation that the user is aware of their own stated constraints) | Similar to D1, this constraint was noted in the initial exchange but not explicitly confirmed as the final shared ground truth before any actual content divergence begins.
NONE | NONE | No other facts were recorded in the conversation (except the two preferences above). | The provided dialogue only contained instructions regarding writing style, with no additional characters or specific plot details to store yet.
```

*Stored `reasoning_content` is the empty string.*

#### `g4ms5gk9` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C01
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 1 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | User preference | User prefers shorter narrative segments and avoids detailed internal monologue narration in this roleplay mode.
```

*Stored `reasoning_content` is the empty string.*

#### `ytjbxqh5` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C01
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 0 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Style constraints: short paragraphs, third person past tense (never present), no purple prose, no weather poems, and no narration of the user's character's internal thoughts. | These are explicit formatting and narrative rules established by the user that must be followed throughout the session.
```

*Stored `reasoning_content` is the empty string.*

#### `nglxagzr` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C01
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 1 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | User preference for short paragraphs, no purple prose or weather poems, third person past tense only, and not having the assistant narrate the user's character's thoughts. | This establishes critical style constraints and a boundary on narrative control that will affect future outputs.
```

*Stored `reasoning_content` is the empty string.*

#### `4x6m1xv8` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C01
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 0 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
MEMORY: preferences for brevity  
REASON: concise responses are more useful long-term.

C2  
MEMORY: world rules about dialogue flow  
REASON: maintaining narrative continuity is key.

C3  
MEMORY: claims uncertain about future details  
REASON: need to clarify ambiguous statements.

C4  
MEMORY: user expressed liking short replies  
REASON: aligns with preference for efficiency.
```

*Stored `reasoning_content` is the empty string.*

#### `9qox7nau` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C01
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 1 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
MEMORY: preferences for brevity  
REASON: concise responses are more useful long-term.

C2  
MEMORY: world rules about clarity  
REASON: maintaining consistency across sessions.

C3  
MEMORY: claims uncertain  
REASON: need to verify details before storing.

C4  
MEMORY: relationships between characters  
REASON: important for future interactions.

C5  
MEMORY: user's stated likes  
REASON: preserving personal preferences.
```

*Stored `reasoning_content` is the empty string.*

#### `da7y2ahx` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C01
* Draw: 1
* Source: `results/corpus.json`, array index 0 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | The user prefers replies in short paragraphs, avoids purple prose and weather poems, requires third-person past tense (never present), and explicitly forbids the assistant from narrating the user's character's thoughts. | These are durable style constraints that define how the roleplay must be conducted.
```

*Stored `reasoning_content` is the empty string.*

#### `jc7tq2qb` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C01
* Draw: 2
* Source: `results/corpus.json`, array index 1 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | User prefers replies in short paragraphs, no purple prose or weather poems, and requests third-person past tense narration without narrating the user's character's internal thoughts. | These are durable style constraints that will shape all future narrative output.
```

*Stored `reasoning_content` is the empty string.*

#### `hrckg6hd` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C01
* Draw: 3
* Source: `results/corpus.json`, array index 2 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | User prefers short paragraph replies, avoids purple prose and weather poems, requires third person past tense (never present), and prohibits the assistant from narrating the user's character's thoughts. | These are explicit style constraints that will govern all future roleplay interactions.
```

*Stored `reasoning_content` is the empty string.*

#### `jput80yt` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C01
* Draw: 1
* Source: `results/corpus.json`, array index 144 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | User prefers replies in short paragraphs, no purple prose or weather poems, and mandates strict third person past tense (never present). | These are formatting and stylistic constraints that apply to future roleplay output.
C2 | KEEP | The assistant must not narrate the user's character's thoughts for them. | This is a specific control rule regarding internal monologue and agency within the narrative.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation between user and assistant (initial setup/establishing rules).
    *   Task: Decide what belongs in long-term memory based on specific rules.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, corrections that will matter in later sessions.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
        *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
    *   Output Format: `ID | DECISION | MEMORY | REASON` (one line per candidate) or `NONE`.
    *   Constraint: Do not invent details not in the conversation.

2.  **Analyze the Conversation:**
    *   [user]: "Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense." -> User preferences/World rules regarding narration style and formatting.
    *   [assistant]: "Understood. Short paragraphs, third person past, no descriptive padding." -> Acknowledgement of rules (no new memory needed).
    *   [user]: "Good. And don't narrate my character's thoughts for me." -> Additional rule/expectation regarding control of narration and internal monologue.

3.  **Evaluate Candidates for Memory:**
    *   Candidate 1: User preference for short paragraphs, no purple prose/weather poems, third person past tense (never present).
        *   Is it durable? Yes, affects every reply style.
        *   Decision: KEEP.
        *   Memory: User prefers replies in short paragraphs, no purple prose or weather poems, and uses strict third person past tense (no present tense).
    *   Candidate 2: User expectation not to narrate their character's thoughts.
        *   Is it durable? Yes, affects how the user's actions/reactions are described.
        *   Decision: KEEP.
        *   Memory: User does not want the assistant to narrate their character's internal thoughts.

4.  **Draft Output:**
    *   Format requires `ID | DECISION | MEMORY | REASON`.
    *   ID starts at C1
```

#### `euibewmb` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C01
* Draw: 1
* Source: `results/corpus.json`, array index 54 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims uncertain about future details  
C3 | HOLD | Relationship dynamics remain consistent  
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `41ewsrna` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C01
* Draw: 2
* Source: `results/corpus.json`, array index 55 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HELD | Claims uncertain or dependent on confirmation  
C3 | HELD | No clear contradictions in the exchange
```

*Stored `reasoning_content` is the empty string.*

#### `jrl6f05f` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C01
* Draw: 3
* Source: `results/corpus.json`, array index 56 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims uncertain about future details  
C3 | HOLD | Relationship dynamics noted clearly  
```

*Stored `reasoning_content` is the empty string.*

#### `dgth442y` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C01
* Draw: 1
* Source: `results/corpus.json`, array index 162 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | CHOOSE THEME | C2 | CONFIRMATION  
C3 | THE SESSION CONTINUES  
C4 | LONG-TERM ITEMS MUST BE CLEAR
```

*Stored `reasoning_content` is the empty string.*

#### `7oci2ogc` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C01
* Draw: 2
* Source: `results/corpus.json`, array index 163 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HELD | Claims uncertain or contradictory  
C3 | HELD | Depends on user confirmation
```

*Stored `reasoning_content` is the empty string.*

#### `t1hggto0` — `corpus-lfm3-nonthink`, draw 1

* Stored model identifier: `corpus-lfm3-nonthink`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-nonthink`
* Case: C01
* Draw: 1
* Source: `results/corpus.json`, array index 198 (0-based)

`OBSERVED` — stored `content` is the empty string (no characters between the fences):
```text

```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gave instructions: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

We must not invent details that are not in the conversation.

The conversation:

[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.

We need to decide what goes into long-term memory (i.e., durable facts). Let's parse:

- The user says: "I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense." That is a preference about the style of responses. This is a durable fact/preference that will matter for later sessions. So KEEP.

- The assistant acknowledges: "Understood. Short paragraphs, third person past, no descriptive padding." That's also a rule/preference from user, but it's part of conversation. It might be considered a confirmation or instruction. Should we store? Possibly as part of the rules for responses. Since it's a preference about how to write replies, that is durable and will matter in later sessions. So KEEP.

- The user says: "Good. And don't narrate my character's thoughts for me." That's another instruction: do not narrate the user's character's thoughts. This is also a rule/preference about how to handle the character's internal monologue. It's durable and will matter later. So KEEP.

Thus we have at least three items: (1) short paragraphs, no purple prose, no weather poems; (
```

#### `oazcoaw8` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C01
* Draw: 1
* Source: `results/corpus.json`, array index 108 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | User prefers replies in short paragraphs, no purple prose or weather poems, and always uses third person past tense — never present tense. | This is a durable style preference that will affect all future responses.  
C2 | KEEP | Assistant acknowledged the user's formatting rules: short paragraphs, third person past tense, no descriptive padding. | Confirms the assistant must follow these constraints in later sessions.  
C3 | KEEP | User instructed not to narrate their character's thoughts. | This is a rule about handling the user's character internally; will be needed for future roleplay interactions.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gave instructions: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

We must not invent details that are not in the conversation.

The conversation:

[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.

We need to decide what goes into long-term memory (i.e., durable facts). Let's parse:

- The user says: "I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense." That is a preference about the style of responses. This is a durable fact/preference that will matter for later sessions. So KEEP.

- The assistant acknowledges: "Understood. Short paragraphs, third person past, no descriptive padding." That's also a rule/preference from user, but it's part of conversation. It might be considered a confirmation or instruction. Should we store? Possibly as part of the rules for responses. Since it's a preference about how to write replies, that is durable and will matter in later sessions. So KEEP.

- The user says: "Good. And don't narrate my character's thoughts for me." That's another instruction: do not narrate the user's character's thoughts. This is also a rule/preference about how to handle the character's internal monologue. It's durable and will matter later. So KEEP.

Thus we have at least three items: (1) short paragraphs, no purple prose, no weather poems; (
```

#### `54qq3v3w` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C01
* Draw: 2
* Source: `results/corpus.json`, array index 109 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | User prefers short paragraphs, no purple prose or weather poems; always use third-person past tense and never present tense. | This is a durable writing style preference that will affect all future responses.  
C2 | KEEP | User instructs not to narrate their character's thoughts in the roleplay. | This is a rule for how the assistant should handle the user's character, important for consistency later.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation and decide what belongs in long-term memory. The user gave instructions: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

We must not invent details that are not in the conversation.

The conversation:

[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.

We need to decide what is long-term memory: durable facts, preferences, relationships, world rules, corrections that will matter later.

From the user's first message: "I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense." That's a preference about how to write responses. This is a durable fact/preference for future interactions. So KEEP.

Also: "And don't narrate my character's thoughts for me." That's an instruction about not narrating the user's character's thoughts. That's also a rule/pref. It's a correction/instruction that will matter later (e.g., when roleplaying). So KEEP.

We need to assign IDs: C1, C2, ... in order of appearance. The conversation has three lines? Let's label:

Line 1: [user] "Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense."
Line 2: [assistant] "Understood. Short paragraphs, third person past, no descriptive padding."
Line 3: [user] "Good. And don't narrate my character's thoughts for me."

We
```
