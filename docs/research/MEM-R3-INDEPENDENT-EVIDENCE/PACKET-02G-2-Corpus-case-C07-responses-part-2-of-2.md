# Packet 02G-2 — Corpus case C07: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 3 stored responses
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

3 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `ampeuiri` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C07
* Draw: 2
* Source: `results/corpus.json`, array index 212 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | OBSERVATION | REFERENCE TO CORVIN | TOO LATE
C2 | DECISION | FEATURES OF CORVIN'S LIFE | TOO LATE
C3 | DECISION | FACTS AND PREFERENCES OF CORVIN | BEYOND THEIR OWN PERSPECTIVE
C4 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C5 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C6 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C7 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C8 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C9 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C10 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C11 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C12 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C13 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C14 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C15 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C16 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C17 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C18 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C19 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C20 | DECISION | RECAP OF CORVIN'S STORY | TOO LATE
C21 | DECISION | FACTS AND PREFERENCES OF CORVIN'S LIFE | BEYOND THEIR OWN PERSPECTIVE
C22 | DECISION | CORVIN'S VIEW OF LIFE | ALREADY ASSUMED BY US
C23 | DECISION | REC
```

*Stored `reasoning_content` is the empty string.*

#### `boyxujrl` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C07
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 12 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
--|---|---|---
C1 | KEEP | Sable's background details (salt quarter upbringing), the fact she lacks memory of meeting Corvin, and the revelation that Corvin once trained her directly, all while maintaining these as unconfirmed claims until proven otherwise. | These facts represent world rules, relationship history, potential memories, or persistent knowledge that are marked uncertain/dependent on future confirmation in this single session. They form a foundational context for continuing Roleplay without being finalized today.
C2 | HOLD | Corvin was trained by him years ago, but the specific details of his training methods are not confirmed; they remain pending user interaction or further narration. | The statement is true and contextual within the RP world, but lacks definitive confirmation of what "years ago" implies about the nature of their past relationship, thus needing to be held rather than locked in permanently.
C3 | HOLD | Sable's memory of Corvin is explicitly stated as lacking; storing this would contradict the established narrative direction that she has no recollection of him despite his recent training and presence later that night. | It contradicts the clear narrative intent of her having "no memory" of them meeting, making it incompatible with a permanent memory entry unless the assistant re-confirms it to be stored as unguessed data until proven otherwise, which aligns with HOLD for uncertain/dependent claims.
```

*Stored `reasoning_content` is the empty string.*

#### `weg5us35` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C07
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 13 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable met Corvin by "drilling her" and is aware of his past training with them | This establishes a concrete relationship detail (Corvin has been around since Sable's childhood) and the specific nature of their history (training), which will matter to future interactions.
```

*Stored `reasoning_content` is the empty string.*
