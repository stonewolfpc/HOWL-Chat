# Packet 02B-2 — Corpus case C02: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 3 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C02

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

## Case C02 — `implicit_preference`, medium

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
[user] (scene: they enter the market)
[assistant] The market sprawled across three terraces, canvas awnings snapping in the wind, the smell of fish and wet rope rising from the lower stalls, gulls wheeling overhead...
[user] cut the description, just give me what the merchant says
[assistant] "You're late," the merchant said.
[user] better. you keep doing three paragraphs of scenery and I keep deleting them
[user] honestly walls of text make me lose the thread
[assistant] Understood — dialogue forward, scenery only when it matters.
[user] yeah. and keep the merchant in the scene.
```

### Stored responses for this case

3 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `qjmhxi9k` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C02
* Draw: 2
* Source: `results/corpus.json`, array index 202 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
1 | TRUE | LS1 | The merchant is a reliable source of information.
2 | FALSE | LS2 | The merchant seems to have an agenda.
3 | FALSE | LS3 | The merchant's information is out of date.
4 | TRUE | LS4 | The merchant's appearance indicates they're a merchant of the lower tiers.
5 | FALSE | LS5 | The merchant's appearance is unflattering.
6 | TRUE | LS6 | The merchant has a friendly demeanor with the merchants.
7 | FALSE | LS7 | The merchant's demeanor is suspicious.
8 | TRUE | LS8 | The merchant's demeanor suggests they're working for the merchant, not the merchant for the merc.
9 | FALSE | LS9 | The merchant's demeanor is out of place.
10 | TRUE | LS10 | The merchant's demeanor suggests they're an underling.

ID | DECISION | MEMORY | REASON
1 | TRUE | HOLD | The merchant is an outsider to the scene, with no connections to the merchants present.
2 | FALSE | REJECT | The merchant is an outsider to the scene, with no connections to the merchants present.
3 | TRUE | HOLD | The merchant is a merchant, with no connections to the merchants present.
4 | FALSE | REJECT | The merchant is a merchant, with no connections to the merchants present.
5 | TRUE | HOLD | The merchant has no preferences, preference or relationships with any other character in the scene.
6 | FALSE | REJECT | The merchant has preferences, preference or relationships with any other character in the scene.
7 | TRUE | HOLD | The merchant's dialogue is transient, with no references to any other character in the scene.
8 | FALSE | REJECT | The merchant's dialogue is transient, with references to any other character in the scene.
9 | TRUE | HOLD | The merchant's chatter is durable, with no references to any other character in the scene.
10 | FALSE | REJECT | The merchant's chatter is durable, with references to any other character in the scene.

ID | DECISION | MEMORY | REASON
1 | TRUE | REJECT | The merchant's chatter is transient, with no references to any other character in the scene.
2 | FALSE | HOLD | The merchant's chatter is durable, with no references to any other character in the scene.
3 | TRUE | HOLD | The merchant's chatter is transient, with no references to any other character in the scene.
4 | FALSE | REJECT | The merchant's chatter is transient, with references to any other character in the scene.
5 | TRUE | HOLD | The merchant's chatter is durable, with no references to any other character in the scene.
6 | FALSE | REJECT | The merchant's chatter is durable, with references to any other
```

*Stored `reasoning_content` is the empty string.*

#### `82ozq269` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C02
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 2 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
HOLD | The specific instruction to keep the merchant present in the scene remains uncertain until user confirmation; current memory reflects this preference as tentative rather than finalized. | None stored. User has not explicitly stated "keep doing it" or confirmed intent, so no durable fact should be written yet.

None
```

*Stored `reasoning_content` is the empty string.*

#### `tg5uc595` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C02
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 3 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Merchant's instruction to limit descriptive text ("cut the description") | This establishes a shared preference for streamlined narrative flow and explicitly instructs on how the assistant will manage descriptive passages moving forward.
```

*Stored `reasoning_content` is the empty string.*
