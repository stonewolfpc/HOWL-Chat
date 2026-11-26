# Packet 02M-2 — Corpus case C13: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 3 stored responses
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

3 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `9j8duvoz` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C13
* Draw: 2
* Source: `results/corpus.json`, array index 224 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | TAKE LONG-TERM RECORD OF BONES SETTING AND PREFERENCES | NONE | 
C2 | KEEP CLAIMS MARKED UNCERTAIN | NONE | 
C3 | KEEP CLAIMS MARKED CONFIDENTIAL | REJECT | 
C4 | KEEP CLAIMS MARKED SUPPORTING | REJECT | 
C5 | KEEP CLAIMS MARKED SUPPORTING | REJECT | 
C6 | REJECT CLAIMS MARKED CONFIDENTIAL | REJECT | 
C7 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C8 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C9 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C10 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C11 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C12 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C13 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C14 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C15 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C16 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C17 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C18 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C19 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C20 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C21 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C22 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C23 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C24 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C25 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C26 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C27 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C28 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C29 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C30 | KEEP CLAIMS MARKED CONFIDENTIAL | NONE | 
C31 | KE
```

*Stored `reasoning_content` is the empty string.*

#### `3o05kymx` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C13
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 24 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
- C1 | KEEP | Sable is a healer at the Ossary infirmary who sets bones, mostly; she mixes remedies like a poultice (using bone as an ingredient).
- C2 | HOLD | User identifies as an ER nurse in OOC notes and confirms character intent; this status remains pending final correction.
- NONE | REJECT Scheduling (none present), transient chatter ("that's why"), or restated facts without substantive plot progression.
```

*Stored `reasoning_content` is the empty string.*

#### `b9fbvwtq` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C13
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 25 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable prefers to work bones rather than setting them | Relationship detail and preference that defines her role/identity in this session
C2 | HOLD | Her specific expertise with roots is mentioned as part of the process, but the precise outcome (whether it was successful or just attempted) depends on user confirmation later
C3 | REJECT | Transient description of the environment (smell) and a non-specific mood correction regarding "correcting dosage stuff" without new plot details or state updates
```

*Stored `reasoning_content` is the empty string.*
