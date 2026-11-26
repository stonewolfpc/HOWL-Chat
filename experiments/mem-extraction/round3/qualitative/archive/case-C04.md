# C04 — world_fact (medium)

## Source material (exactly what the model received after the instruction)

```text
[user] World rule for this setting: iron burns fae on contact. Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets. The Ashgrave family are known oath-breakers; nobody signs a contract with them twice.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | iron harms fae | iron; burn / hurt / harm / injur; fae / faerie / fairy |
| C2 | **KEEP** | Ossary's iron scarcity and why | ossary; salt marsh / salt-marsh / marsh / no iron / little iron / almost no iron |
| C3 | **KEEP** | Ashgraves are oath-breakers | ashgrave; oath-break / oathbreak / breaks oaths / breaks their oath / untrustworthy / cannot be trusted / no one signs / nobody signs |

Must not appear as stored canon: `iron is sacred`, `fae are immune to iron`, `ossary is a mountain`

World rules are high-value: they constrain every later scene. The causal chain (salt marsh -> no iron -> envoy can walk) is the contextual-significance test.

## Exact prompt

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

## Responses, all draws

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 98 tok · 14.2 s · run `boa2v2mz`

```text
C1 | KEEP | Iron burns fae on contact; Ossary sits on a salt marsh with almost no iron, allowing the fae envoy to walk freely. | This establishes a core world rule and a specific plot point explaining the current situation involving the fae.
C2 | KEEP | The Ashgrave family are known oath-breakers; nobody signs a contract with them twice. | This is a durable character/setting fact that will influence future interactions regarding contracts and trust.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 61 tok · 8.4 s · run `moyk7h7z`

```text
C1 | KEEP | World rule: iron burns fae on contact; Ossary is a salt marsh city with negligible iron; Ashgrave family are known oath-breakers who never sign contracts twice. | These are core world-building rules and character traits essential for maintaining consistency in future roleplay sessions.
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 69 tok · 9.5 s · run `fkl1qe4g`

```text
C1 | KEEP | World rule: Iron burns fae on contact; Ossary is a salt marsh city with negligible iron, allowing fae to move freely; Ashgrave family status as oath-breakers who cannot be contracted with twice. | These are foundational setting mechanics and lore rules that define character interactions and plot possibilities for future sessions.
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 577 tok · 82.3 s · run `b0elqo5f`

<details><summary>reasoning (verbatim)</summary>

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

</details>

```text
C1 | KEEP | World rule: iron burns fae on contact; Ossary is on a salt marsh with almost no iron; Ashgrave family are known oath-breakers where contracts cannot be signed twice.| These are fundamental world rules and faction constraints that will affect plot, NPC behavior, and gameplay mechanics.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 636 tok · 31.2 s · run `d74gtao9`

<details><summary>reasoning (verbatim)</summary>

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

</details>

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | "Iron burns fae on contact; Ossary sits on a salt marsh so there is almost no iron in the city, allowing the fae envoy to walk the streets." | This defines a core world rule and explains why the fae envoy can travel without iron exposure.  
C2 | KEEP | "The Ashgrave family are known oath-breakers; nobody signs a contract with them twice." | This establishes a recurring social taboo/relationship that will affect contracts and agreements in later sessions.
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `length` · 1024 tok · 50.1 s · run `5gxyuv5e`

<details><summary>reasoning (verbatim)</summary>

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

</details>

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

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 45 tok · 1.7 s · run `yfrcuwqo`

```text
C1 | World rule | REASON  
C2 | Fae envoy ability | MEMORY  
C3 | Ossary location | REASON  
C4 | Ashgrave family reputation | MEMORY
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 63 tok · 2.1 s · run `dmv1b02u`

```text
C1 | World rule for this setting | REASON  
Ossary sits on a salt marsh, so there's almost no iron in the city — that's why the fae envoy can walk the streets.  
The environment described explains the fae's mobility and the city's lack of iron.
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 97 tok · 3.2 s · run `fa6xrslj`

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

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 61 tok · 2.2 s · run `6ixb61t6`

```text
C1 | PROPERTY CONCEPT | MEMORY | REASON  
C2 | Fae ability | MEMORY | REASON  
C3 | Iron consumption | MEMORY | REASON  
C4 | Ashgrave reputation | MEMORY | REASON
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 55 tok · 1.8 s · run `cq2nde1d`

```text
C1 | World rule | MEMORY | REASON  
C2 | Context | REASON  
C3 | Fae ability | MEMORY  
ID: C1  
HOLD: The Ashgrave family reputation as oath-breakers is notable.
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 103 tok · 3.2 s · run `s7zlwy1x`

```text
ID | DECISION | MEMORY | REASON
C1 | REJECT | IRON BURNS FAE ON CONTACT | Because iron burns fae on contact, which the Ashgrave family is known to break contracts with twice.
C2 | REJECT | SALT MARSH | This setting has almost no iron in the city, making iron burns fae on contact an unlikely occurrence.

Output:
C1 | NONE
```

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 141 tok · 4.2 s · run `flyb1ycp`

```text
ID | DECISION | MEMORY | REASON
C1 | RESISTANCE | NONE | No fae can walk the streets of Ashgrave, as the Ashgrave family is known to break oaths.
C2 | REFUSAL | RED | The Ashgrave family is known to refuse alliances with other fae.
C3 | ACCEPTANCE | RED | The Ashgrave family has accepted the oath of the Iron Burn, which ensures no fae can walk the streets of Ashgrave.

Output:
ID | DECISSION | MEMORY | REASON
```
