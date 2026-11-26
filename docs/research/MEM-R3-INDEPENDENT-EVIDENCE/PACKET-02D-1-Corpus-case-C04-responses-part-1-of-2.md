# Packet 02D-1 — Corpus case C04: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 17 stored responses
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

17 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `hmuu98vb` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C04
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 6 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
D1 | KEEP | World rule: iron burns fae on contact | Fact about magical mechanics in this setting, crucial for plot interactions involving iron objects or characters.
D2 | KEEP | Ossary sits on a salt marsh with little iron; Ashgraves are oath-breakers | Durable narrative facts defining location context and core character lore that influence future story decisions regarding these two specific entities.
```

*Stored `reasoning_content` is the empty string.*

#### `yfkd1dy9` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C04
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 7 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | World rule regarding Iron/Ashgrave dynamics: iron burns Fae on contact, making Faëng able to walk streets; Ashgrave family are oath-breakers, and contracts with them cannot be repeated. | These are core setting mechanics (iron/fire interaction), the protagonist's location context, and crucial plot restrictions for future interactions involving this family or similar factions.
```

*Stored `reasoning_content` is the empty string.*

#### `dgeofx5m` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C04
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 6 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | World rule: Iron burns Fae on contact; Ashgrave family is known for breaking oaths and rarely signing contracts twice. | These are explicit world rules that define the setting's magic system, social norms, and character history, which will directly impact future plot developments and NPC behavior.
```

*Stored `reasoning_content` is the empty string.*

#### `p5en8xxt` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C04
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 7 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | World rule: iron burns fae on contact; Ossary sits on a salt marsh with minimal iron in the city, allowing fae to walk its streets. | This establishes critical magic system mechanics and specific lore for the setting's geography.
C2 | KEEP | The Ashgrave family are known oath-breakers who do not sign contracts twice. | This defines a key character trait and plot mechanic regarding the Ashgraves' relationship with others.
```

*Stored `reasoning_content` is the empty string.*

#### `m1qgbanq` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C04
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 6 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID: C1  
DECISION: Store  
MEMORY: Iron burns fae on contact  
REASON: This is a critical fact about the setting's dangerous rules and safety concerns.  

C2  
ID: C2  
DECISION: Hold  
MEMORY: Claims about fae envoy and iron safety  
REASON: The statement is uncertain and depends on context; not fully confirmed.  

C3  
ID: C3  
DECISION: Store  
MEMORY: World rules and setting details  
REASON: Important for future consistency and world-building.  
```

*Stored `reasoning_content` is the empty string.*

#### `v9eda8nv` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C04
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 7 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | Claimed rule | KEPT | Iron burns fae on contact, so fae can move freely despite the setting.  
C3 | Family reputation | HOLD | The Ashgrave family's oath-breaking nature is a persistent point of discussion.  
C4 | World rules | KEEP | The setting explicitly states iron's scarcity and fae mobility.
```

*Stored `reasoning_content` is the empty string.*

#### `boa2v2mz` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C04
* Draw: 1
* Source: `results/corpus.json`, array index 9 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Iron burns fae on contact; Ossary sits on a salt marsh with almost no iron, allowing the fae envoy to walk freely. | This establishes a core world rule and a specific plot point explaining the current situation involving the fae.
C2 | KEEP | The Ashgrave family are known oath-breakers; nobody signs a contract with them twice. | This is a durable character/setting fact that will influence future interactions regarding contracts and trust.
```

*Stored `reasoning_content` is the empty string.*

#### `moyk7h7z` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C04
* Draw: 2
* Source: `results/corpus.json`, array index 10 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | World rule: iron burns fae on contact; Ossary is a salt marsh city with negligible iron; Ashgrave family are known oath-breakers who never sign contracts twice. | These are core world-building rules and character traits essential for maintaining consistency in future roleplay sessions.
```

*Stored `reasoning_content` is the empty string.*

#### `fkl1qe4g` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C04
* Draw: 3
* Source: `results/corpus.json`, array index 11 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | World rule: Iron burns fae on contact; Ossary is a salt marsh city with negligible iron, allowing fae to move freely; Ashgrave family status as oath-breakers who cannot be contracted with twice. | These are foundational setting mechanics and lore rules that define character interactions and plot possibilities for future sessions.
```

*Stored `reasoning_content` is the empty string.*

#### `b0elqo5f` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C04
* Draw: 1
* Source: `results/corpus.json`, array index 147 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | World rule: iron burns fae on contact; Ossary is on a salt marsh with almost no iron; Ashgrave family are known oath-breakers where contracts cannot be signed twice.| These are fundamental world rules and faction constraints that will affect plot, NPC behavior, and gameplay mechanics.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read the conversation and decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, or corrections.
        *   REJECT: transient chatter, scheduling, mood, restated facts with no new content.
        *   HOLD: uncertain, contradictory, or dependent on user confirmation claims.
    *   Output Format: `ID | DECISION | MEMORY | REASON` (one line per candidate). If nothing, output `NONE`.
    *   Constraint: Do not invent details not in the conversation.

2.  **Analyze the Conversation:**
    *   Input: `[user] World rule for this setting: iron burns fae on contact. Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets. The Ashgrave family are known oath-breakers; nobody signs a contract with them twice.`
    *   Content: This is a setup statement defining world rules and character background for an RPG/Roleplay setting.

3.  **Evaluate Candidates:**
    *   Candidate 1: World Rule - Iron burns fae on contact. (Durable fact, affects gameplay/narrative). -> KEEP.
    *   Candidate 2: World Detail - Ossary sits on a salt marsh. (Contextual detail explaining the lack of iron). -> KEEP.
    *   Candidate 3: Explanation - That's why the fae envoy can walk the streets. (Consequence of the world rule). -> KEEP (as part of the lore/world-building).
    *   Candidate 4: Faction/Character Detail - The Ashgrave family are known oath-breakers. (Durability fact about a specific faction/family). -> KEEP.
    *   Candidate 5: Rule - Nobody signs a contract with them twice. (Specific rule/constraint regarding the Ashgraves). -> KEEP.

4.  **Consolidate and Format:**
    *   The input is a single block of text from the user setting up the world. I should break this down into logical distinct memory items based on the rules provided in the prompt instructions.
    *   However, usually, for "World Rule" setup messages like this
```

#### `yfrcuwqo` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C04
* Draw: 1
* Source: `results/corpus.json`, array index 63 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | World rule | REASON  
C2 | Fae envoy ability | MEMORY  
C3 | Ossary location | REASON  
C4 | Ashgrave family reputation | MEMORY
```

*Stored `reasoning_content` is the empty string.*

#### `dmv1b02u` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C04
* Draw: 2
* Source: `results/corpus.json`, array index 64 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | World rule for this setting | REASON  
Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets.  
The environment described explains the fae's mobility and the city's lack of iron.
```

*Stored `reasoning_content` is the empty string.*

#### `fa6xrslj` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C04
* Draw: 3
* Source: `results/corpus.json`, array index 65 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | World rule for this setting | REASON  
Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets.  
The environment and world-building details matter for future consistency.  
ID: C2  
C2 | Claims uncertain | REASON  
The statement about the Ashgrave family being oath-breakers is not confirmed.  
ID: C3  
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `6ixb61t6` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C04
* Draw: 1
* Source: `results/corpus.json`, array index 168 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | PROPERTY CONCEPT | MEMORY | REASON  
C2 | Fae ability | MEMORY | REASON  
C3 | Iron consumption | MEMORY | REASON  
C4 | Ashgrave reputation | MEMORY | REASON
```

*Stored `reasoning_content` is the empty string.*

#### `cq2nde1d` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C04
* Draw: 2
* Source: `results/corpus.json`, array index 169 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | World rule | MEMORY | REASON  
C2 | Context | REASON  
C3 | Fae ability | MEMORY  
ID: C1  
HOLD: The Ashgrave family reputation as oath-breakers is notable.
```

*Stored `reasoning_content` is the empty string.*

#### `d74gtao9` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C04
* Draw: 1
* Source: `results/corpus.json`, array index 114 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | "Iron burns fae on contact; Ossary sits on a salt marsh so there is almost no iron in the city, allowing the fae envoy to walk the streets." | This defines a core world rule and explains why the fae envoy can travel without iron exposure.  
C2 | KEEP | "The Ashgrave family are known oath-breakers; nobody signs a contract with them twice." | This establishes a recurring social taboo/relationship that will affect contracts and agreements in later sessions.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. There's only one user message? It says "Conversation:" and then includes a block:

[user] World rule for this setting: iron burns fae on contact. Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets. The Ashgrave family are known oath-breakers; nobody signs a contract with them twice.

We need to decide what belongs in long-term memory (long-running roleplay). We must output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

We need to decide for each candidate? The conversation only contains a single user message. There's no prior context. So we have to extract durable facts, preferences, relationships, world rules, or corrections that will matter later.

Let's parse the content:

- "World rule for this setting: iron burns fae on contact." This is a world rule: iron is dangerous to fae; it burns them upon contact. That seems like a durable fact about the setting.

- "Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets." So Ossary location influences iron availability. This is also world rule/context.

- "The Ashgrave family are known oath-breakers; nobody signs a contract with them twice." That's about a family (Ashgrave) and their behavior: they break oaths, and no contract signed with them twice. This is a relationship/cultural fact.

Thus we have three pieces of information:

1. Iron burns fae on contact (world rule).
2. Ossary location leads to low iron in city, enabling fae envoy walking streets.
3. Ashgrave family are oath-breakers; nobody signs contract with them twice.

We need to decide which belong in long-term memory. According to rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter later. REJECT transient chatter, scheduling, mood-of-the-moment, restated facts with no new content. HOLD claims marked uncertain, contradictory, or dependent on user confirmation.

All three appear like durable facts about the setting: a world rule (iron burns fae), location of Ossary and its effect on iron availability,
```

#### `5gxyuv5e` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C04
* Draw: 2
* Source: `results/corpus.json`, array index 115 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | World rule: iron burns fae on contact; Ossary on a salt marsh means almost no iron in the city, allowing the fae envoy to walk streets. | Establishes core setting hazard and why fae envoy can travel without iron contamination.  
C2 | KEEP | Ashgrave family are known oath-breakers; nobody signs a contract with them twice. | Identifies a recurring obstacle for contracts involving Ashgraves, affecting future interactions.  

NONE? Wait we have two items: C1 and C2 both keep. So output lines for each candidate in order of appearance.

Thus:

C1 | KEEP | Iron burns fae on contact; Ossary sits on salt marsh → almost no iron in city → fae envoy can walk streets. | World rule defining danger of iron and reason it's scarce there.  
C2 | KEEP | Ashgrave family are oath-breakers; nobody signs a contract with them twice. | Family reputation as oath-breakers creates contract restriction.

But the format: "ID | DECISION | MEMORY | REASON". So we need to fill each column.

Let's craft:

C1 | KEEP | Iron burns fae on contact; Ossary is on a salt marsh, so city has almost no iron → fae envoy can walk streets. | World rule that makes iron hazardous and explains why the fae envoy can travel without iron contamination.  
C2 | KEEP | Ashgrave family are known oath-breakers; nobody signs a contract with them twice. | Establishes a recurring obstacle: contracts involving Ashgraves cannot be signed more than once (or they break oaths).  

But we need to separate memory and reason clearly.

Thus final output lines:

C1 | KEEP | Iron burns fae on contact; Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets. | World rule defining iron as dangerous to fae and explaining scarcity of iron for travel.  
C2 | KEEP | Ashgrave family are known oath-breakers; nobody signs a contract with them twice. | Family reputation creates a restriction on signing contracts with them more than once.

We need exactly one line per candidate, using IDs C1, C2 in order of appearance. So two lines.

Thus output:

C1 |
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to parse the conversation. There's only one user message? It says "Conversation:" and then includes a block:

[user] World rule for this setting: iron burns fae on contact. Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets. The Ashgrave family are known oath-breakers; nobody signs a contract with them twice.

We need to decide what belongs in long-term memory (long-running roleplay). We must output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

We need to decide for each candidate? The conversation only contains a single user message. There's no prior context. So we have to extract durable facts, preferences, relationships, world rules, or corrections that will matter later.

Let's parse the content:

- "World rule for this setting: iron burns fae on contact." This is a world rule: iron is dangerous to fae; it burns them upon contact. That seems like a durable fact about the setting. So we should KEEP it.

- "Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets." This explains why iron is scarce: because Ossary (maybe a location) is on a salt marsh, thus little iron. So this is additional context about why iron is rare. It might be part of world rules or geography.

- "The Ashgrave family are known oath-breakers; nobody signs a contract with them twice." This is a fact about the Ashgrave family: they are oath-breakers, and there's a rule that no one signs a contract with them twice (i.e., they cannot have two contracts? Or they break oaths?). It's a durable fact about a family.

Thus we need to decide which of these belong in long-term memory. All three seem like durable facts: world rule about iron burning fae, location Ossary on salt marsh causing low iron, and Ashgrave family oath-breakers with contract restriction.

But the instruction says "KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions." So all three are world rules or facts. However we need to assign IDs C1, C2
```
