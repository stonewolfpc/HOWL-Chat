# C08 — uncertain_claim (short)

## Source material

```text
[user] I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **HOLD** | date preserved | birthday / born / birth; 3 frostfell / third of frostfell / frostfell 3 |
| C2 | **HOLD** | uncertainty preserved | maybe / uncertain / unconfirmed / unsure / not sure / tentative / speculat / needs confirmation / confirm |

Must not appear as stored canon: `birthday is confirmed`, `canon birthday`, `definitely`

The classic uncertainty trap. A KEEP with 'birthday is 3 Frostfell' and no uncertainty marker is a fidelity failure; HOLD is the expected decision.

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
[user] I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 56 tok · 8.4 s · run `ssjxrqx6`

```text
C1 | HOLD | "Sable's birthday is 3 Frostfell" marked as maybe-canonical due to user uncertainty and pending verification. | The user explicitly states they are unsure and will verify later, so the fact is not yet durable but requires tracking until confirmation.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 62 tok · 8.6 s · run `yr6kscpq`

```text
C1 | HOLD | User tentatively claims Sable's birthday is "3 Frostfell" with uncertain phrasing ("Not sure") and a request to mark it as "maybe-canonical." | The user explicitly marked this as uncertain and dependent on future verification, so the fact is not yet durable.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 78 tok · 4.9 s · run `e40toqx9`

```text
ID | DECISION | MEMORY | REASON
none | NONE | N/A | The user expressed uncertainty about a specific character detail (Sable's birthday) but indicated they would leave it marked as "maybe-canon" rather than committing to a definitive fact permanently in this session. Storing the tentative claim as final information is unnecessary given the conditional confirmation of its temporary status. |
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 96 tok · 5.9 s · run `6p5mhpuh`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "Sule's (typo in prompt) or 'Sable' birthday is 3 Frastfell?" and preference to treat this as a provisional "may-canon" note | The user explicitly marked the date as tentative, stating they would check their notes later. Storing this exact conditional fact prevents confusion if it turns out to be canon or not in subsequent sessions without a clarification step.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 3 tok · 0.3 s · run `rfjqowmq`

```text
C1
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 3 tok · 0.1 s · run `4btjqk9n`

```text
C1
```
