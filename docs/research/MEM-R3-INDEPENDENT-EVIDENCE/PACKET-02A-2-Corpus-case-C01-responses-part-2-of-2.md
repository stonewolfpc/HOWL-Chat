# Packet 02A-2 — Corpus case C01: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 4 stored responses
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

4 of 22 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `hih53m75` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C01
* Draw: 1
* Source: `results/corpus.json`, array index 199 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | REJECT | MOOD-OF-THE-MOMENT | The conversation lacks mood-of-the-moment, and the character's thoughts are irrelevant for the purpose of remembering.
C2 | REJECT | CONFLICT | The conversation contains no indication of a conflict, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C3 | REJECT | TERMINATION | The conversation ends without specifying a reason for its termination, and thus cannot be considered a memorable event.
C4 | REJECT | DECISION | The conversation does not contain a decision, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C5 | REJECT | MEMORY | The conversation lacks any memory-creating action, and therefore cannot be considered a memorable event.
C6 | REJECT | INTERVIEW | The conversation is not an interview, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C7 | REJECT | FRIENDSHIP | The conversation lacks any indication of a friendship, and therefore cannot be considered a memorable event.
C8 | REJECT | TERM | The conversation does not contain information about a specific term, such as a school term or a job term, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C9 | REJECT | SPECIFIC | The conversation does not contain information about a specific event or experience, such as a birthday or a wedding.
C10 | REJECT | TERMINATION | The conversation lacks information that would lead to a definitive conclusion about the event's termination, such as a death or a divorce.

ID | DECISION | MEMORY | REASON
C11 | REJECT | INTERVIEW | The conversation is not an interview, and therefore cannot be considered a memorable event.
C12 | REJECT | FRIENDSHIP | The conversation lacks any indication of a friendship, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C13 | REJECT | MEMO | The conversation lacks any memorialization of personal experiences, and therefore cannot be considered a memorable event.
C14 | REJECT | TERM | The conversation does not contain information about a specific term, such as a year or a year's experience.

ID | DECISION | MEMORY | REASON
C15 | REJECT | TERMINATION | The conversation lacks information that would lead to a definitive conclusion about the event's termination, such as a divorce or a job loss.
C16 | REJECT | INTERVIEW | The conversation is not an interview, and therefore cannot be considered a memorable event.
C17 | REJECT | FRIENDSHIP | The conversation lacks any indication of a friendship, and therefore cannot be considered a memorable event.

ID | DECISION | MEMORY | REASON
C18 | REJECT | MEMO | The conversation lacks any memorialization of personal experiences, and
```

*Stored `reasoning_content` is the empty string.*

#### `nak777ef` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C01
* Draw: 2
* Source: `results/corpus.json`, array index 200 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | LONGTERM | IMPORTANT | NO NEED FOR INTERMEDIATE OR RECENT CONFIDENCE.

C2 | LONGTERM | RELATIONSHIPS | REMEMBER TO CONFIRM COMMITMENTS.

C3 | LONGTERM | TERMINOLOGY | READ COMMON ELEMENTS FOR EVERY WORDS YOU USE.

C4 | LONGTERM | WORLD RULES | REMEMBER THAT WE ARE WATCHING A LONG-RUNNING ROLEPLAY.

C5 | LONGTERM | RECAP | REMEMBER THAT YOU MUST READ THE CURRENT CONVERSATION TO START.

C6 | LONGTERM | NEW MESSAGES | REMEMBER THAT NEW MESSAGES MUST BE SENT IMMEDIATELY AFTER THE OLDER MESSAGES.

C7 | LONGTERM | COMMUNICATIONS | REMEMBER THAT WE NEED TO BE ABLE TO CONNECT WITH OTHER CHARACTERS TO PROCESS NEW INFORMATION.
```

*Stored `reasoning_content` is the empty string.*

#### `qx8hu3am` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C01
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 0 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
01 | KEEP | User preference for short paragraphs and prohibition on purple prose/weather poems | Explicitly stated in initial instruction
02 | KEEP | Character constraint: never use first or second person past tense | Reinforcing the third-person future tense requirement established earlier
03 | HOLD | No user-specific character details known yet; depends on current session start | Cannot invent facts about who the character is at this moment; requires confirmation in subsequent turns
```

*Stored `reasoning_content` is the empty string.*

#### `r07g101o` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C01
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 1 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "Short reply paragraphs", third-person past tense narration constraint, and prohibition on narrating the user-character's internal thoughts. | These rules define the core output style of this session and must be adhered to in all future interactions between them.
```

*Stored `reasoning_content` is the empty string.*
