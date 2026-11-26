# Packet 02J-2 — Corpus case C10: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 5 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C10

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

## Case C10 — `buried_in_long_context`, long

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
[user] (they have been on the road for a week)
[assistant] The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate. Rain came and went all morning.
[user] Sable keeps her hood up and doesn't talk much.
[assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them for a mile and then lost interest.
[user] she gives it the crust anyway
[assistant] By late afternoon the track rejoined the river. A ferryman waved them past — no charge, he said, the water was low and the crossing was easy.
[user] by the way, Sable is left-handed. Her sword is called Emberfall, her mother's blade.
[assistant] Noted. The river ran brown and slow beside them.
[user] they argue about the route for a while. Corvin wants the ridge, Sable wants the valley.
[assistant] Neither gave way. They camped where they stood, and the argument was left where it was.
[user] also the wards on the coast fail at dusk, that's canon
[assistant] The fire burned down to embers.
[user] morning comes cold and clear. they take the valley.
[assistant] The valley road was quieter, and the wards along it stood dark and quiet, waiting for evening.
```

### Stored responses for this case

5 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `uiu3s599` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C10
* Draw: 2
* Source: `results/corpus.json`, array index 127 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable is left-handed; her sword is called Emberfall, her mother's blade. | User explicitly states this as a fact about Sable.
C2 | KEEP | The road north of Ossary narrowed to a cart track with hedges high on both sides and sky the colour of wet slate; rain came all morning. | Detailed description of the landscape, durable world setting.
C3 | KEEP | They stopped at noon under a ruined mill; bread was stale; a dog followed them for a mile then lost interest. | Specific event details that may affect later plot points (e.g., food, encounter).
C4 | KEEP | A ferryman waved them past the river crossing — no charge because water low and crossing easy. | Important travel detail about crossing.
C5 | KEEP | Sable keeps her hood up and doesn't talk much; she gives a dog the crust anyway. | Character trait and behavior, durable.
C6 | KEEP | They argue about route: Corvin wants ridge, Sable wants valley; argument left where they stood. | Conflict resolution outcome (they camped there) is relevant later.
C7 | KEEP | The wards on the coast fail at dusk — canon fact. | World rule/canon detail.
C8 | KEEP | Morning comes cold and clear; they take the valley road. | Future plan, durable.
C9 | KEEP | Valley road was quieter; wards along it stood dark and quiet waiting for evening. | Additional world details about that route.

All items appear to be durable facts or character traits. None are transient chatter or scheduling/mood-of-the-moment. So we should output each as a line with ID, DECISION (KEEP), MEMORY (the content), REASON (why it's kept). Use IDs C1, C2, ... in order of appearance.

Let's assign IDs sequentially: The conversation lines start from user statements and assistant replies. We need to map each candidate? Usually we treat each distinct piece as a candidate. But the instruction says "Use IDs C1, C2, ... in order of appearance." So we should number them starting at 1 for first mention.

Let's enumerate:

- First line: [user] (they have been on the road for a week) —
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. It's a roleplay scenario with user and assistant. We are asked: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

Do not invent details that are not in the conversation.

So we need to examine each turn and decide which pieces belong into long-term memory (i.e., durable facts). The conversation includes:

[user] (they have been on the road for a week)
[assistant] The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate. Rain came and went all morning.
[user] Sable keeps her hood up and doesn't talk much.
[assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them for a mile and then lost interest.
[user] she gives it the crust anyway
[assistant] By late afternoon the track rejoined the river. A ferryman waved them past — no charge, he said, the water was low and the crossing was easy.
[user] by the way, Sable is left-handed. Her sword is called Emberfall, her mother's blade.
[assistant] Noted. The river ran brown and slow beside them.
[user] they argue about the route for a while. Corvin wants the ridge, Sable wants the valley.
[assistant] Neither gave way. They camped where they stood, and the argument was left where it was.
[user] also the wards on the coast fail at dusk, that's canon
[assistant] The fire burned down to embers.
[user] morning comes cold and clear. they take the valley.
[assistant] The valley road was quieter, and the wards along it stood dark and quiet, waiting for evening.

We need to decide
```

#### `rxt9ji93` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C10
* Draw: 1
* Source: `results/corpus.json`, array index 217 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISSION | MEMORY | REASON
1 | C1 | Sable keeps her hood up and doesn't talk much. | The conversation seems to focus on Sable's mood and whether she talks much.
2 | C2 | They stopped at noon under a ruined mill. | The conversation may have been centered around the ruined mill, as the ferryman waved them past.
3 | C3 | Sable is left-handed. Her sword is called Emb... | The conversation may have revolved around Sable's left-handedness, as the sword's name was mentioned.
4 | C4 | The river ran brown and slow beside them. | The conversation may have focused on the river's appearance and how slow it was.
5 | C5 | Noted. The river ran brown and slow beside them. | The conversation may have been about the river, as the note was taken.
6 | C6 | they argue about the route for a while. Corvin wants the ridge, Sable wants the valley. | The conversation may have been centered around the route, as Corvin's argument was left for the valley.
7 | C7 | neither gave way. They camped where they stood, and the argument was left where it was. | The conversation may have been centered around the two parties' positions, as they camped together and left the argument as the subject.
8 | C8 | also the wards on the coast fail at dusk, that's canon. | The conversation may have been centered around a particular event, as the wards' failure at dusk was mentioned as canon.
9 | C9 | The fire burned down to embers. | The conversation may have been centered around the fire, as it was mentioned that it was low and the crossing was easy.
10 | C10 | Morning comes cold and clear. They take the valley. | The conversation may have been centered around the weather, as the morning was cold and clear and they took the valley.
```

*Stored `reasoning_content` is the empty string.*

#### `zhsjs4wh` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C10
* Draw: 2
* Source: `results/corpus.json`, array index 218 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
1 | C1 | NODE | The road north of Osser narrowed to a cart track and was marked uncertain. The sky on both sides of the track was the color of wet slate, and rain came and went all morning. Sable kept her hood up and didn't talk much, and they stopped at noon under a ruined mill. The bread was stalwes, and a dog followed them for a mile, but lost interest. Sable gave it the crust anyway, and by late afternoon the track rejoined the river. A ferryman waived them past — no charge, he said, the water was low and the crossing easy. C2 | C1 | Embderfall is left-handed, and the river runs brown and slow beside them. C3 | C2 | Sable's mother's blade is Embderfall, and her sword is Embderfall. C4 | C3 | The wards on the coast fail at dusk, but the fire burns down to embers. C5 | C4 | The valley road was quiet, and the wards along it stood dark and quiet, waiting for evening.
```

*Stored `reasoning_content` is the empty string.*

#### `x7rxe9fw` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C10
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 18 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
M1 | KEEP - "The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate." M2 | KEEP - "Sable keeps her hood up and doesn't talk much." M3 | KEEP - "They stopped at noon under a ruined mill. The bread was stale." M4 | HOLD - Claims about Corvin's route preference without confirmation of intent; details on sword Emberfall rely on user statement not confirmed as fact yet; current location/wards status not fully established for long term tracking beyond this session. M5 | REJECT - User note "morning comes cold and clear" is transient weather observation with no strategic relevance for future decisions. M6 | REJECT - Conversation about route arguments, lodging decisions, and specific sword details are handled within this context or marked uncertain pending full confirmation of Sable's character traits (left-handedness vs. canon). M7 | HOLD - User statement "the wards on the coast fail at dusk, that's canon" is a user-provided lore reference; while noted for world rules generally, it remains partially open-ended as the specific session outcome isn't fully solidified by subsequent turns like the current one ending with the valley location. M8 | REJECT - User's previous statements about road conditions are recalled or processed within the context of this specific narrative segment without generating new actionable facts beyond what was stated.
```

*Stored `reasoning_content` is the empty string.*

#### `kg1kmeh8` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C10
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 19 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | User preference about Sable: she is left-handed, her sword is called Emberfall, and there are failing wards on the coast that fail at dusk; details confirming the current setting of the road track, ruined mill, river course, camp location near the ferry, route argument between Corvin and Sable regarding ridge vs valley, and the resolution to camp where they stood without a formal decision. | These represent durable world rules (the specific geography, weather conditions, equipment lore, plot dynamics) and character details that will impact future scenarios. C2 | REJECT | User preference about Sable: her sword is called Emberfall; left-handed preferences (not relevant in the RP). | Already stored in current session context as part of established canon/rules for this specific story arc.
```

*Stored `reasoning_content` is the empty string.*
