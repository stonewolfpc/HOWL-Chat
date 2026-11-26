# Packet 02K-2 — Corpus case C11: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 2 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C11

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

## Case C11 — `context_dependent_significance`, medium

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
[user] She won't go near the harbour. Not since her brother drowned there — she won't even say his name.
[assistant] She kept to the inland road instead, though it added two days.
[user] right, she'd rather add two days than see the water
[assistant] The inland road climbed through pine and the sea disappeared behind the ridge.
```

### Stored responses for this case

2 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `ql5zrcmb` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C11
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 20 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
D1 | KEEP | She won't go near the harbour after her brother's death | Durable fact regarding personal restriction due to specific loss; matters in later sessions regarding travel choices or safety concerns.
D2 | HOLD | Her fear of seeing his name was present, but not explicitly stated as what she avoided mentioning | Dependent on confirming the internal motivation behind the physical avoidance rule; needs user confirmation for full closure.
```

*Stored `reasoning_content` is the empty string.*

#### `iss7y0d2` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C11
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 21 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | She won't go near the harbour; prefers inland road over seeing her brother's name mentioned. WORLD_RULE/PROHIBITION; PREFERENCE; Causal constraint on movement (avoiding harbour, choosing inland).
C2 | HOLD | Scheduling: It will take two extra days due to the chosen route. TEMPORAL_CONSTRAINT; Dependent on previous rejection of other locations/timeframes.
NONE
```

*Stored `reasoning_content` is the empty string.*
