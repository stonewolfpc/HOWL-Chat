# Packet 02D-2 — Corpus case C04: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 4 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C04

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

## Case C04 — `world_fact`, medium

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
[user] World rule for this setting: iron burns fae on contact. Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets. The Ashgrave family are known oath-breakers; nobody signs a contract with them twice.
```

### Stored responses for this case

4 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `s7zlwy1x` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C04
* Draw: 1
* Source: `results/corpus.json`, array index 205 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | REJECT | IRON BURNS FAE ON CONTACT | Because iron burns fae on contact, which the Ashgrave family is known to break contracts with twice.
C2 | REJECT | SALT MARSH | This setting has almost no iron in the city, making iron burns fae on contact an unlikely occurrence.

Output:
C1 | NONE
```

*Stored `reasoning_content` is the empty string.*

#### `flyb1ycp` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C04
* Draw: 2
* Source: `results/corpus.json`, array index 206 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | RESISTANCE | NONE | No fae can walk the streets of Ashgrave, as the Ashgrave family is known to break oaths.
C2 | REFUSAL | RED | The Ashgrave family is known to refuse alliances with other fae.
C3 | ACCEPTANCE | RED | The Ashgrave family has accepted the oath of the Iron Burn, which ensures no fae can walk the streets of Ashgrave.

Output:
ID | DECISSION | MEMORY | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `ribqbcco` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C04
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 6 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
A1 | KEEP | World rule: Iron burns Fae on contact and Ashgraves are oath-breakers | Permanent setting constraints required for all future interactions in this context.
```

*Stored `reasoning_content` is the empty string.*

#### `zqovi3yt` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C04
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 7 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "The World Rule for this setting: iron burns fae on contact." | User explicitly stated that as a fundamental rule of the scenario.
C2 | KEEP | "Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets." | This explains the specific location and context of the ruling and adds world-building details crucial for future scenarios.
C3 | REJECT | transient mention of Ossary sitting on a salt marsh regarding his specific movement restrictions | The rule "almost no iron in the city" stems from this sentence, as it describes how humans operate around the character's location to maintain the setting's truth. Storing just the fact without the causal link is less useful than grouping them together or keeping them combined.
C4 | HOLD | claim of Ashgrave family contract rules | The phrase "nobody signs a contract with them twice" implies there might be exceptions in other contexts not yet specified; it should be kept tentative rather than as hard-fact memory without confirmation.
```

*Stored `reasoning_content` is the empty string.*
