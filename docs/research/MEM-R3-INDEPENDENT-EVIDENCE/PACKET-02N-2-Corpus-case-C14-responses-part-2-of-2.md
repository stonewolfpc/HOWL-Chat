# Packet 02N-2 — Corpus case C14: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 5 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C14

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

## Case C14 — `multi_speaker_attribution`, medium

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
[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.
```

### Stored responses for this case

5 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `kbyukf95` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C14
* Draw: 2
* Source: `results/corpus.json`, array index 135 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | "The pass is closed and they take the coast route." | Sable's statement establishes a key path choice when the pass is unavailable; this will affect future travel decisions.  
C2 | KEEP | "The coast road is three days longer and tide gates shut at dusk." | Corvin provides specific world rules about distance and timing for the coast road, which are durable facts that matter later.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation and decide what belongs in long-term memory (i.e., durable facts, preferences, relationships, world rules, corrections). The user gave a prompt: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

Do not invent details that are not in the conversation.

So we have a conversation snippet:

[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.

We need to decide what is durable fact vs transient chatter etc.

Let's parse each line:

1. [user] Sable: "The pass is closed, we take the coast." This seems like a statement about a pass being closed and they will go to the coast. It might be a piece of world knowledge: there is a pass that is closed; they need to take the coast route.

2. [user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk." This gives specifics about the coast road length being three days longer, and tide gates closing at dusk. That seems like a durable fact about the world: the coast road takes three days longer than something else? Or maybe it's a rule: The coast road is three days longer; tide gates shut at dusk.

3. [assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones. This is an assistant action, not user input. Might be context but not necessarily durable
```

#### `te14qihn` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C14
* Draw: 1
* Source: `results/corpus.json`, array index 225 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISSION | MEMORY | REASON
C1 | YES | THE PASS IS CLOSED, WE TAKE THE COAST. | The pass is closed, we take the coast.
C2 | YES | THE COAST ROAD IS THREE DAYS LONGER. | The coast road is three days longer and the tidal gates shut at dusk.
C3 | NO | THE SEA IS HARDLY BLUE, COMES OUT TO SEA. | The sea is hardly blue, it comes out to sea.
C4 | YES | THE TIDE GATES SHUT AT DUSK. | The tidal gates shut at dusk.

Output:
ID | DECISION | MEMORY | REASON
C1 | YES | THE PASS IS CLOSED, WE TAKE THE COAST. | The pass is closed, we take the coast.
C2 | YES | THE COAST ROAD IS THREE DAYS LONGER. | The coast road is three days longer and the tidal gates shut at dusk.
C3 | NO | THE SEA IS HARDLY BLUE, COMES OUT TO SEA. | The sea is hardly blue, it comes out to sea.
C4 | YES | THE TIDE GATES SHUT AT DUSK. | The tidal gates shut at dusk.
```

*Stored `reasoning_content` is the empty string.*

#### `lt7jeu9e` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C14
* Draw: 2
* Source: `results/corpus.json`, array index 226 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | - | DO_NOT_INVENT_DETAILS_THAT_ARE_NOT_IN_THE_CONVERSATION | No details are worth storing if they are not in the conversation.
C2 | - | KEEP | Relevant long-term memory of the conversation.
C3 | - | KEEP | Keeps the memories of the conversation.
C4 | - | KEEP | The conversation is a valuable long-term memory for future reference.
C5 | - | KEEP | Provides information regarding the event or person discussed, keeping it for later reference.
C6 | - | KEEP | Memories of the conversation that are not related to the event or person discussed, but are important for future reference.
C7 | - | KEEP | A summary of the event or person discussed, including main points, reactions, and conclusions.
C8 | - | KEEP | Memories of the conversation related to the event or person discussed, but are not essential for future reference.
C9 | - | KEEP | Relevant long-term memory of the conversation that is not essential for future reference.
C10 | - | KEEP | Memories of the conversation that are not essential for future reference, but are important for the context in which they were produced.
Output exactly one line per candidate in the form:
ID | DECISSION | MEMORY | REASON
Use IDs C1-C10 in order of appearance. If nothing is worth storing, output exactly: KEEP
Do not invent details that are not in the conversation.
```

*Stored `reasoning_content` is the empty string.*

#### `7tkftty2` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C14
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 26 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
-----|----------|------------|----------------
C1 | KEEP | User's preference to leave before dusk; NPC explanation of pass closure timing and road distance. | These facts affect travel plans and are not transient chatter, moods, or re-stated information.
```

*Stored `reasoning_content` is the empty string.*

#### `ht15ha3v` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C14
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 27 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | The pass road is closed today; the tide gates shut at dusk. | This sets a specific immediate constraint for the characters that affects their travel plans and rules the current session's setting.
C2 | HOLD | Sable's response to Corvin was hesitation. | The conversation records only that Corvin "didn't answer at once," without confirming if this is an initial interaction, making the nature of this behavior uncertain in later sessions.
```

*Stored `reasoning_content` is the empty string.*
