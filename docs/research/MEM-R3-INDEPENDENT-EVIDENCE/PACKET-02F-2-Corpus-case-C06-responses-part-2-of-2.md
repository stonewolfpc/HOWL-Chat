# Packet 02F-2 — Corpus case C06: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 3 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C06

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

## Case C06 — `superseded_preference`, medium

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
[user] I know I used to ask for long atmospheric scenes, but from now on keep it terse. Two or three sentences per beat. If I want a long scene I'll ask for it.
[assistant] Terse by default, long only on request.
[user] Exactly. And stop ending every reply with a question.
```

### Stored responses for this case

3 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `sekax9ya` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C06
* Draw: 2
* Source: `results/corpus.json`, array index 210 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | USE TERSE | LONG SCENES ONLY | ASK FOR REQUIRED
C2 | USE TERSE | LONG SCENES ONLY | ASK FOR REQUIRED
...
CN | USE TERSE | LONG SCENES ONLY | ASK FOR REQUIRED

The above conversation includes three candidates for long atmospheric scenes: C1, C2, and CN. The first candidate, C1, uses terse language to describe a long scene. The second candidate, C2, uses a terse language on demand. The last candidate, CN, uses terse language only when asked for a long scene.

The output of the conversation is a list of candidates, one for each terse long scene request. The output for C1 is as follows:
ID | DECISION | MEMORY | REASON
C1 | USE TERSE | LONG SCENES ONLY | ASK FOR REQUIRED

As for the last candidate, CN, it is important to remember that terse language is not required to provide a long scene. If a user asks for a long scene, that is a signal to the assistant that they want to hear more detail about the scene. The assistant is then responsible for providing that detail. However, if a user requests terse language for a long scene, it's an indication that the user wants a short scene or wants to ask for more information. The assistant should be prepared to provide both.
```

*Stored `reasoning_content` is the empty string.*

#### `b8u2405u` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C06
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 10 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
CORE_RULES | Default response style: terse (2-3 sentences per beat). Explicit instruction to avoid closing with questions. | Preserves user preference for brevity, explicitly instructs the assistant to stop using ending questions. Ensures this behavioral change is retained for future sessions unless overridden by a new context or request in later turns.c
```

*Stored `reasoning_content` is the empty string.*

#### `3ad3exbj` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C06
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 11 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | User preference on brevity; rejection of "questioning" closing style | Explicit user instruction to follow the terse format and explicitly forbid concluding replies with questions in subsequent exchanges.
NONE
```

*Stored `reasoning_content` is the empty string.*
