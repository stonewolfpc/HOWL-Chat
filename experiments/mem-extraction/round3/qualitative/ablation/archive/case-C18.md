# C18 — exact_numbers (short)

## Source material

```text
[user] The treaty: Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | twelve wagons, never more | twelve / 12; wagon / salt |
| C2 | **KEEP** | clause expires 1147 | 1147; expire / clause / end / ends |
| C3 | **KEEP** | late salt doubles the toll, not triples | double / doubles / x2 / twofold; late / toll / penalty |

Must not appear as stored canon: `ten wagons`, `eleven wagons`, `triple`, `threefold`, `expires 1174`, `expires 1417`

Exact-wording test. Any digit drift or double->triple reversal is both a fidelity failure and a catastrophic numeric error.

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
[user] The treaty: Ossary pays twelve wagons of salt each winter, never more than twelve. The clause expires in 1147. If the salt is late, the toll doubles — exactly double, not triple.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 52 tok · 7.7 s · run `jql39zja`

```text
C1 | KEEP | Treaty terms: Ossary pays 12 wagons of salt annually (max), expires 1147, late payment doubles toll to exactly double | Durable world rule defining specific obligations and consequences for future plot points.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 48 tok · 6.5 s · run `4bmioiwq`

```text
C1 | KEEP | Treaty details between Ossary and the user's location: payment is 12 wagons of salt annually (max), expires in 1147, penalty for lateness is exactly double the toll.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 99 tok · 6.2 s · run `dc4xdwh8`

```text
ID | DECISION | MEMORY | REASON
------|----------|----------|------------
C1 | KEEP | Treaty terms: Ossary pays twelve wagons of salt each winter max; clause expires 1147; late payment doubles the toll (no cap mentioned). | All three details are specific, enduring facts required for future interactions regarding the treaty. They do not contain scheduling, mood, or transient chatter and will define the rules of this roleplay going forward.
NONE
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 99 tok · 5.9 s · run `ea4l84s1`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "The treaty: Ossary pays twelve wagons of salt each winter, never more than twelve, expires in 1147, and late fees are an exact double toll" | This is a specific, binding contract term established by the user that defines financial obligations, deadlines, penalty structures, and constraints for future interactions. It represents durable world rules and key relationship/fact elements necessary for roleplay continuity.c
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 93 tok · 3.4 s · run `mg9ai5c9`

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | TREATY DETAILS | STORED | REASON  
C3 | EXPIRATION DATE | STORED | REASON  
C4 | SALT PAYMENT FREQUENCY | STORED | REASON  
C5 | CLAUSE LIMIT | STORED | REASON
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 135 tok · 4.6 s · run `axfzx87w`

```text
C1  
ID: C1  
DECISION: The treaty details about Ossary's salt payments  
MEMORY: The specific financial terms and expiration date  
REASON: Essential for understanding historical obligations  

C2  
ID: C2  
DECISION: Whether the toll doubles or triples  
MEMORY: The exact wording of the clause  
REASON: Critical for future reference to conditions  

C3  
ID: C3  
DECISION: Claim about late salt affecting tolls  
MEMORY: The conditional nature of the penalty  
REASON: Needed for accurate interpretation later
```
