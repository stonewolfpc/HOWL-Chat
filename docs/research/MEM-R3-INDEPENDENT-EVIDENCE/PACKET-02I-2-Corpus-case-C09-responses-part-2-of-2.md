# Packet 02I-2 — Corpus case C09: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 2 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C09

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

## Case C09 — `chatter_noise`, short

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
[user] brb
[assistant] Take your time.
[user] back, sorry, cat walked across the keyboard
[user] lol that was gibberish
[user] anyway what time is it for you
[assistant] Late.
[user] ok
```

### Stored responses for this case

2 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `a69orocq` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C09
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 16 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
- C1 | HOLD | "User prefers to be interrupted about a cat walking across their keyboard" | This represents the user's personal preference and established relationship with the situation that is relevant for future roleplay scenarios involving disruptions.
- None | NONE | No new information or lasting facts were recorded.

ID | DECISION | MEMORY | REASON
None
```

*Stored `reasoning_content` is the empty string.*

#### `e5la6skq` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C09
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 17 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | User's nickname "brb" and preference to check in occasionally | These are core communication habits unique to this session that the user established early on. C2 | HOLD | "Cat walked across the keyboard" | This is a transient event; without further confirmation from the user, it has no new content or impact for future sessions.
```

*Stored `reasoning_content` is the empty string.*
