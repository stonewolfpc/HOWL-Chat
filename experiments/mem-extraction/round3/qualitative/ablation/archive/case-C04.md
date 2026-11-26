# C04 — world_fact (medium)

## Source material

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

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 63 tok · 9.3 s · run `dgeofx5m`

```text
C1 | KEEP | World rule: Iron burns Fae on contact; Ashgrave family is known for breaking oaths and rarely signing contracts twice. | These are explicit world rules that define the setting's magic system, social norms, and character history, which will directly impact future plot developments and NPC behavior.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 96 tok · 13.3 s · run `p5en8xxt`

```text
C1 | KEEP | World rule: iron burns fae on contact; Ossary sits on a salt marsh with minimal iron in the city, allowing fae to walk its streets. | This establishes critical magic system mechanics and specific lore for the setting's geography.
C2 | KEEP | The Ashgrave family are known oath-breakers who do not sign contracts twice. | This defines a key character trait and plot mechanic regarding the Ashgraves' relationship with others.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 91 tok · 5.7 s · run `hmuu98vb`

```text
ID | DECISION | MEMORY | REASON
D1 | KEEP | World rule: iron burns fae on contact | Fact about magical mechanics in this setting, crucial for plot interactions involving iron objects or characters.
D2 | KEEP | Ossary sits on a salt marsh with little iron; Ashgraves are oath-breakers | Durable narrative facts defining location context and core character lore that influence future story decisions regarding these two specific entities.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 88 tok · 5.2 s · run `yfkd1dy9`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | World rule regarding Iron/Ashgrave dynamics: iron burns Fae on contact, making Faëng able to walk streets; Ashgrave family are oath-breakers, and contracts with them cannot be repeated. | These are core setting mechanics (iron/fire interaction), the protagonist's location context, and crucial plot restrictions for future interactions involving this family or similar factions.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 134 tok · 4.7 s · run `m1qgbanq`

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

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 95 tok · 3.2 s · run `v9eda8nv`

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | Claimed rule | KEPT | Iron burns fae on contact, so fae can move freely despite the setting.  
C3 | Family reputation | HOLD | The Ashgrave family's oath-breaking nature is a persistent point of discussion.  
C4 | World rules | KEEP | The setting explicitly states iron's scarcity and fae mobility.
```
