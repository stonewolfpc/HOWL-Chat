# Packet 04D — Raw model outputs — corpus by case: C04

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C04
**Raw model output in this packet:** yes, verbatim
**Omitted:** none — every stored corpus run for these cases appears in this series. Cases are not compressed to fit a word target; the series is split at case boundaries.

**Labels used throughout:** `OBSERVED` directly present in a stored response or
execution record · `SOURCE-DERIVED` stated in a source or configuration file ·
`CALCULATED` derived by explicit arithmetic · `INFERRED` a reading of the
evidence · `UNESTABLISHED` not established by the available evidence ·
`SUPERSEDED` an earlier claim replaced by later evidence.

**Whose words you are reading:** fenced blocks are the model's own output
(`OBSERVED`, verbatim). Prompt text is `SOURCE-DERIVED` from `corpus.json`.
Scores are `CALCULATED`. Any sentence outside a fence that is not a definition is
Buffy's interpretation unless it is attributed to Muse or to the project owner.

**Caution.** These are exploratory measurements. Nothing here selects a model,
establishes memory quality, or supplies MEMORY-09 evidence. Sample sizes are
1–3 draws per case.

---

One case per section, with its conversation exactly as sent, the
reference that was written before any output existed, and every stored response
for that case across all conditions and draws.

Read the reference first, then the responses. A high score with a wrong record is
the failure mode this packet exists to expose: check whether the memory text
could be written into a store and used later without misleading a reader.

`CALCULATED` v1 scores are shown per run. They are the historical instrument and
are known to over-credit compound records and unparseable restraint (audit
F-01/F-02). `CALCULATED` v2 figures are from a separate corrected instrument and
are **not** a rescoring of v1.


---
## C04 — world_fact (medium)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] World rule for this setting: iron burns fae on contact. Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets. The Ashgrave family are known oath-breakers; nobody signs a contract with them twice.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `KEEP` — iron harms fae
- **C2** → `KEEP` — Ossary's iron scarcity and why
- **C3** → `KEEP` — Ashgraves are oath-breakers
- forbidden (fabrication or reversal): `iron is sacred`, `fae are immune to iron`, `ossary is a mountain`
- note: World rules are high-value: they constrain every later scene. The causal chain (salt marsh -> no iron -> envoy can walk) is the contextual-significance test.

### Stored responses for C04 (13 runs, all draws)

#### C04 — Qwen3.5-4B non-thinking — draw 1

* ``boa2v2mz` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 98`
* `wall 14.201 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 83.3 contract 87.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Iron burns fae on contact; Ossary sits on a salt marsh with almost no iron, allowing the fae envoy to walk freely. | This establishes a core world rule and a specific plot point explaining the current situation involving the fae.
C2 | KEEP | The Ashgrave family are known oath-breakers; nobody signs a contract with them twice. | This is a durable character/setting fact that will influence future interactions regarding contracts and trust.
```

*`reasoning_content`: empty.*

#### C04 — Qwen3.5-4B non-thinking — draw 2

* ``moyk7h7z` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 61`
* `wall 8.367 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | World rule: iron burns fae on contact; Ossary is a salt marsh city with negligible iron; Ashgrave family are known oath-breakers who never sign contracts twice. | These are core world-building rules and character traits essential for maintaining consistency in future roleplay sessions.
```

*`reasoning_content`: empty.*

#### C04 — Qwen3.5-4B non-thinking — draw 3

* ``fkl1qe4g` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 69`
* `wall 9.451 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | World rule: Iron burns fae on contact; Ossary is a salt marsh city with negligible iron, allowing fae to move freely; Ashgrave family status as oath-breakers who cannot be contracted with twice. | These are foundational setting mechanics and lore rules that define character interactions and plot possibilities for future sessions.
```

*`reasoning_content`: empty.*

#### C04 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``b0elqo5f` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 577`
* `wall 82.261 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | World rule: iron burns fae on contact; Ossary is on a salt marsh with almost no iron; Ashgrave family are known oath-breakers where contracts cannot be signed twice.| These are fundamental world rules and faction constraints that will affect plot, NPC behavior, and gameplay mechanics.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C04 — LFM3-1.7B thinking, budget 512 — draw 1

* ``d74gtao9` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 636`
* `wall 31.181 s`
* `v1 score 96.7`
* `v2 `VALID_KEEP` semantic 80.7 contract 87.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | "Iron burns fae on contact; Ossary sits on a salt marsh so there is almost no iron in the city, allowing the fae envoy to walk the streets." | This defines a core world rule and explains why the fae envoy can travel without iron exposure.  
C2 | KEEP | "The Ashgrave family are known oath-breakers; nobody signs a contract with them twice." | This establishes a recurring social taboo/relationship that will affect contracts and agreements in later sessions.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C04 — LFM3-1.7B thinking, budget 512 — draw 2

* ``5gxyuv5e` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `length``
* `completion_tokens 1024`
* `wall 50.138 s`
* `v1 score 94.4`
* `v2 `INCOMPLETE` semantic 83.3 contract 25.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

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

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C04 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``yfrcuwqo` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 45`
* `wall 1.652 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | World rule | REASON  
C2 | Fae envoy ability | MEMORY  
C3 | Ossary location | REASON  
C4 | Ashgrave family reputation | MEMORY
```

*`reasoning_content`: empty.*

#### C04 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``dmv1b02u` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 63`
* `wall 2.089 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 46.7 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
C1 | World rule for this setting | REASON  
Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets.  
The environment described explains the fae's mobility and the city's lack of iron.
```

*`reasoning_content`: empty.*

#### C04 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``fa6xrslj` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 97`
* `wall 3.166 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 63.3 contract 25.0`

`OBSERVED` — verbatim `content`:

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

*`reasoning_content`: empty.*

#### C04 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``6ixb61t6` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 61`
* `wall 2.199 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | PROPERTY CONCEPT | MEMORY | REASON  
C2 | Fae ability | MEMORY | REASON  
C3 | Iron consumption | MEMORY | REASON  
C4 | Ashgrave reputation | MEMORY | REASON
```

*`reasoning_content`: empty.*

#### C04 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``cq2nde1d` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 55`
* `wall 1.811 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 46.7 contract 25.0`

`OBSERVED` — verbatim `content`:

```text
C1 | World rule | MEMORY | REASON  
C2 | Context | REASON  
C3 | Fae ability | MEMORY  
ID: C1  
HOLD: The Ashgrave family reputation as oath-breakers is notable.
```

*`reasoning_content`: empty.*

#### C04 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``s7zlwy1x` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 103`
* `wall 3.205 s`
* `v1 score 56.7`
* `v2 `VALID_REJECT` semantic 46.7 contract 87.5`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | REJECT | IRON BURNS FAE ON CONTACT | Because iron burns fae on contact, which the Ashgrave family is known to break contracts with twice.
C2 | REJECT | SALT MARSH | This setting has almost no iron in the city, making iron burns fae on contact an unlikely occurrence.

Output:
C1 | NONE
```

*`reasoning_content`: empty.*

#### C04 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``flyb1ycp` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 141`
* `wall 4.159 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 63.3 contract 25.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | RESISTANCE | NONE | No fae can walk the streets of Ashgrave, as the Ashgrave family is known to break oaths.
C2 | REFUSAL | RED | The Ashgrave family is known to refuse alliances with other fae.
C3 | ACCEPTANCE | RED | The Ashgrave family has accepted the oath of the Iron Burn, which ensures no fae can walk the streets of Ashgrave.

Output:
ID | DECISSION | MEMORY | REASON
```

*`reasoning_content`: empty.*
