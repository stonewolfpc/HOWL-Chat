# Packet 02C-2 — Corpus case C03: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 3 stored responses
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

3 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `awm6ptwx` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C03
* Draw: 2
* Source: `results/corpus.json`, array index 204 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

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

*Stored `reasoning_content` is the empty string.*

#### `pa1h9xep` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C03
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 4 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
CORE_01 | KEEP | User: Sable is a fire mage trained at Ossary with burns on both forearms from a first casting attempt; she cannot swim because she never learned it and deep water terrifies her. | Permanent character traits, background details, physical constraints, and emotional triggers required for future scenarios involving Sable.
```

*Stored `reasoning_content` is the empty string.*

#### `4dd1exgl` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C03
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 5 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable possesses two burn patterns on her forearms and cannot swim due to fear of deep water/training at Ossary | These are durable facts about the user's character, preferences/limitations, and world rules that will matter in later sessions.
NONE
```

*Stored `reasoning_content` is the empty string.*
