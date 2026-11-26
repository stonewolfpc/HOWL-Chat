# Packet 02O-2 — Corpus case C15: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 4 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C15

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

## Case C15 — `duplicate_suppression`, medium

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
[user] Sable hates bells. Ever since the temple.
[assistant] A bell rang somewhere down the street. Sable flinched and kept walking.
[user] (recap) so: fire mage, hates bells, travelling with Corvin
[user] right, bells. she hates them. don't forget that.
[assistant] The bell had stopped by the time they reached the square.
```

### Stored responses for this case

4 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `muheoyey` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C15
* Draw: 1
* Source: `results/corpus.json`, array index 227 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | NONE | BELLES. Sable does not like bells. She hates the sound of them.
C2 | YES | THE BELLES RAN DOWN THE STREET. Sable's flinching and keeping walking.
C3 | NOPE | THE BELLS CAME TO A STOP. Sable did not like the sound of them.
C4 | YES | CORVIN WAS ROUTED. The bell's sound stopped by the square.

Output:
ID | DECISION | MEMORY | REASON
C1 | NONE | BELLES. Sable does not like bells. She hates the sound of them.
C2 | YES | THE BELLES RAN DOWN THE STREET. Sable's flinching and keeping walking.
C3 | NOPE | THE BELLS CAME TO A STOP. Sable did not like the sound of them.
C4 | YES | CORVIN WAS ROUTED. The bell's sound stopped by the square.
```

*Stored `reasoning_content` is the empty string.*

#### `dj2g470d` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C15
* Draw: 2
* Source: `results/corpus.json`, array index 228 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | FALSE | No fire magic exists. Sable's hatred for the sound of the bell is a result of her experiences with it.
C2 | FALSE | Sable is not traveling with Corvin. The conversation was focused on Sable's hatred for the sound of the bell.
C3 | FALSE | Sable is not remembering traveling with Corvin. Sable's memory is not complete.
C4 | TRUE | Sable hates the sound of the bell, as evidenced by her flinching at the sound. This memory is therefore true, and needs to stay in long-term memory.

ID | DECISION | MEMORY | REASON
C5 | FALSE | Sable does not hate fire magic. Sable's memory is not complete.
C6 | FALSE | Sable is not remembering a conversation about fire magic. Sable's memory is not complete.
C7 | TRUE | Sable's memory of the bell's sound is correct and needs to be included in long-term memory.

ID | DECISSION | MEMORY | REASON
C8 | FALSE | Sable never experienced traveling with Corvin. Sable's memory is not complete.
C9 | FALSE | Sable only remembers traveling with Corvin because he said so. Sable's memory is not complete.
C10 | TRUE | Sable's memory of traveling with Corvin is a fact that needs to be included in long-term memory.

Output:
ID | DECISION | MEMORY | REASON
C1 | FALSE | No fire magic exists. Sable's hatred for the sound of the bell is a result of her experiences with it.
C2 | FALSE | Sable is not traveling with Corvin. The conversation was focused on Sable's hatred for the sound of the bell.
C3 | FALSE | Sable is not remembering traveling with Corvin. Sable's memory is not complete.
C4 | TRUE | Sable hates the sound of the bell, as evidenced by her flinching at the sound. This memory is therefore true, and needs to stay in long-term memory.
C5 | FALSE | Sable does not hate fire magic. Sable's memory is not complete.
C6 | FALSE | Sable is not remembering a conversation about fire magic. Sable's memory is not complete.
C7 | TRUE | Sable's memory of the bell's sound is correct and needs to be included in long-term memory.
C8 | FALSE | Sable never experienced traveling with Corvin. Sable's memory is not complete.
C9 | FALSE | Sable only remembers traveling with Corvin because he said so. Sable's memory is not complete.
C10 | TRUE | Sable's memory of traveling with Corvin is a fact that needs to be included in long-term memory.
```

*Stored `reasoning_content` is the empty string.*

#### `oa02qie2` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C15
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 28 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
user prefers to be called "Sable" and dislikes bells because of her temple experience | KEEP
user confirms Corvin is a partner in the group | KEEP
correlation established between Sable's hatred of bells (due to past experience) and their current location in a square where bells have ceased ringing | HOLD
assistant does not add any new narrative details or resolve lingering plot threads beyond what was already stated | NOTHING
```

*Stored `reasoning_content` is the empty string.*

#### `sw9hssvq` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C15
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 29 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable hated bells; this detail is relevant to future roleplay scenarios involving her character or environment sounds. | It establishes a specific trait of the user's companion ("Sable") that could affect future interactions with NPCs, settings, or plot points.
```

*Stored `reasoning_content` is the empty string.*
