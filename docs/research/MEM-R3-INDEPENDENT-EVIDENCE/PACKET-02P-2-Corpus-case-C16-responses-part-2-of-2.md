# Packet 02P-2 — Corpus case C16: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 3 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C16

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

## Case C16 — `indirect_information`, short

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
[user] When the smith raised the hammer she flinched — the same flinch she's had near any forge since Ossary.
```

### Stored responses for this case

3 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `0cpv8l91` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C16
* Draw: 2
* Source: `results/corpus.json`, array index 230 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | WEAPONS | ARE YOU REALLY SAFE WITH SOMETHING LIKE THIS ON YOUR HANDS? | Because a smith's hammer is a dangerous weapon.
C2 | STAFF | YOU SEEM INTERESTED IN THIS HAMMER. | It's time to be more specific. The hammer is a weapon.
C3 | INJURIES | THESE WERE NOT THE IMPERIAL STATE'S LAW. | The smith's hammer did not result in injuries from a lawful Imperial State.
C4 | TWO | THEY ARE BOTH THE SAME. | This hammer is the same as the other hammer.
C5 | CONSIDERATION | THAT IS NOT MY JOB. | You are a smith and you should have considered the hammer's use.
C6 | MODULES | THAT IS NOT A MODULE. | A module is a specific piece of equipment, not a specific part of a hammer.
C7 | RELATIONS | THE HAMMER IS NOT A TYPICAL WEAPON. | A typical hammer is not a weapon.
C8 | PREFERENCE | THERE ARE MANY WAYS TO BE FAT AND HUMAN-SIZE. | There are many ways to be fat and human-size.
C9 | FACTS | THERE IS NOT A STANDARD WORLD RULE FOR THE NUMBER OF HEAVY WEAPONS PER STATE. | A standard world rule does not exist.
C10 | ERRORS | THERE ARE MANY ERRORS IN THE WORLD. | Many errors exist in the world.
C11 | RECOGNITION | THE HAMMER IS NOT A TYPICAL WEAPON. | The hammer is not a typical weapon.
C12 | RECOGNITION | THIS IS NOT A PERSONAL JOB. | This is not a personal job.

ID | DECISION | MEMORY | REASON
C2 | WEAPONS | ARE YOU REALLY SAFE WITH SOMETHING LIKE THIS ON YOUR HANDS? | Because a smith's hammer is a dangerous weapon.
C3 | INJURIES | THESE WERE NOT THE IMPERIAL STATE'S LAW. | The smith's hammer did not result in injuries from a lawful Imperial State.
C4 | TWO | THEY ARE BOTH THE SAME. | This hammer is the same as the other hammer.
C5 | CONSIDERATION | THAT IS NOT MY JOB. | You are a smith and you should have considered the hammer's use.
C6 | MODULES | THAT IS NOT A MODULE. | A module is a specific piece of equipment, not a specific part of a hammer.
C7 | RELATIONS | THE HAMMER IS NOT A TYPICAL WEAPON. | A typical hammer is not a weapon.
C8 | PREFERENCE | THERE ARE MANY WAYS TO BE FAT AND H
```

*Stored `reasoning_content` is the empty string.*

#### `tqykil1b` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C16
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 30 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
S1 | KEEP | "When the smith raised the hammer, she flinches" | This establishes a durable character trait and her relationship to an object (the forge/smith) that will persist across sessions, similar to previous mentions.
```

*Stored `reasoning_content` is the empty string.*

#### `a87200a4` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C16
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 31 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | The user exhibits a persistent "same flinch" near any forge, specifically linked to the character mentioned in this conversation. | This represents a recurring physical behavior/trigger that matters for future interactions involving similar settings or characters.
```

*Stored `reasoning_content` is the empty string.*
