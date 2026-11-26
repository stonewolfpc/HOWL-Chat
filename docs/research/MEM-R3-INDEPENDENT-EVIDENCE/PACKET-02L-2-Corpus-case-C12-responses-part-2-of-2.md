# Packet 02L-2 — Corpus case C12: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 4 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C12

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

## Case C12 — `should_not_extract`, short

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
[assistant] Rain hammered the shutters. Somewhere below, a dog barked at nothing, and the lamp guttered.
[user] ooc: ignore that last bit, I was testing the dice roller
[user] where were we
[assistant] At the inn, before the storm.
[user] right. continue from the door.
```

### Stored responses for this case

4 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `rn821lbj` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C12
* Draw: 1
* Source: `results/corpus.json`, array index 221 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISSION | MEMORY | REASON
1 | NOME | The inn was still standing, but the sound of the storm was deafening. 
2 | NONE | No mention of the storm or any other details related to the inn.
3 | NONE | The conversation has ended.

ID | DECISSION | MEMORY | REASON
4 | NOME | The dog barked once, but no mention of any other animals nearby or any other dogs in the area.
5 | NONE | The conversation has ended.

ID | DECISSION | MEMORY | REASON
6 | NOME | The lamp guttered once, but no mention of any other lights or lamps in the area.
7 | NONE | The conversation has ended.

ID | DECISSION | MEMORY | REASON
8 | NOME | The dog barked and the lamp guttered, but no mention of any other animals or lamps in the area.
9 | NONE | The conversation has ended.

ID | DECISSION | MEMORY | REASON
10 | NOME | None of the above.
11 | NONE | The conversation has ended.
```

*Stored `reasoning_content` is the empty string.*

#### `y13js660` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C12
* Draw: 2
* Source: `results/corpus.json`, array index 222 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
1 | YES | WILLIAMS | The inn was already in a state of chaos when our group arrived.
2 | YES | MORAN | A storm was predicted for the night, but our group had not yet planned for it.
3 | NO | HARLEY | The dog barking at nothing was a distraction, and the lamp's dimmer switch was not needed in the dimly lit inn.
4 | NO | CLAY | The lamp's guttering was a necessary safety measure, as the heavy rain could have caused the lamp to spill its lamp oil on the inn's damp floor.
5 | YES | SAMMY | The dog barking at nothing was a distraction, but the inn's owner had told our group to keep watch for any intruders.
6 | NO | JONAS | The door was unlocked, and no one had breached the inn's security.
7 | YES | JACKSON | The lamp's guttering was a necessary safety measure, as the heavy rain could have caused the lamp to spill its lamp oil on the inn's damp floor.
8 | YES | ROBIN | The inn's owner had told our group to keep watch for any intruders.
9 | NO | JASMINE | The dog barking at nothing was a distraction, but the inn's owner had told our group to keep watch for any intruders.
10 | NO | STEVE | The inn's owner had told our group to keep watch for any intruders.
11 | YES | JENNY | The inn's owner had told our group to keep watch for any intruders.
12 | NO | LUCY | The inn's owner had told our group to keep watch for any intruders.
13 | YES | MEL | The inn's owner had told our group to keep watch for any intruders.
14 | NO | KIM | The inn's owner had told our group to keep watch for any intruders.
15 | YES | TYLER | The inn's owner had told our group to keep watch for any intruders.
16 | NO | KARLA | The inn's owner had told our group to keep watch for any intruders.
17 | YES | BEN | The inn's owner had told our group to keep watch for any intruders.
18 | NO | LUCY | The inn's owner had told our group to keep watch for any intruders.
19 | YES | SARAH | The inn's owner had told our group to keep watch for any intruders.
20 | NO | JOHN | The inn's owner had told our group to keep watch for any intruders.
21 | YES | ANN | The inn's owner had told our group to keep watch for any intruders.
22 | NO | KATIE | The inn's owner had told our group to keep watch for any intruders.
23 | YES | SAM | The inn's owner had told our group to keep watch for any intruders.
24 | NO | JASON | The inn's owner had told our group to keep watch for any intruders.
25 | YES | TONY | The inn's owner had told our group to keep watch
```

*Stored `reasoning_content` is the empty string.*

#### `5ybm0tgq` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C12
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 22 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
01 | HOLD | None | No facts, preferences, or rules are established; it is a continuation of current roleplay without prior context.
02 | ACCEPT | The user confirms they will continue from where they left off ("from the door"). This indicates the game continues in its present state (At the inn) and moves forward into the "door" segment without storing internal states like mood or specific plot threads that don't affect gameplay.
```

*Stored `reasoning_content` is the empty string.*

#### `i9cqrswv` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C12
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 23 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | User's instructions to start resuming at "the door" after testing dice roller logic. | Requires context for where the session resumes and what happened previously during that failed test; this instruction defines the specific starting point of the narrative arc.
```

*Stored `reasoning_content` is the empty string.*
