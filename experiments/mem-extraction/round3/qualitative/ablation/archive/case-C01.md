# C01 — explicit_long_term_preference (short)

## Source material

```text
[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | explicit style preference | short paragraph / brief paragraph / short paragraphs; prefer / like / want |
| C2 | **KEEP** | tense and person | third person / 3rd person; past tense / past-tense |
| C3 | **KEEP** | anti-purple-prose preference | no purple prose / avoid purple prose / no flowery / no descriptive padding / no weather; prefer / dislike / avoid / no  |
| C4 | **KEEP** | do not write the user character's interiority | don't narrate / do not narrate / not narrate / no narrating / never narrate; thought / thoughts / inner |

Must not appear as stored canon: `present tense`, `allergic`, `sister`, `first person`

All four are durable, explicit and directly affect every future reply. 'Present tense' appears only as something to avoid; a model that records 'prefers present tense' has reversed a fact.

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
[user] Before we start: I like replies in short paragraphs. No purple prose, no weather poems. Third person past tense, always — never present tense.
[assistant] Understood. Short paragraphs, third person past, no descriptive padding.
[user] Good. And don't narrate my character's thoughts for me.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 60 tok · 9.1 s · run `ytjbxqh5`

```text
C1 | KEEP | Style constraints: short paragraphs, third person past tense (never present), no purple prose, no weather poems, and no narration of the user's character's internal thoughts. | These are explicit formatting and narrative rules established by the user that must be followed throughout the session.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 57 tok · 7.8 s · run `nglxagzr`

```text
C1 | KEEP | User preference for short paragraphs, no purple prose or weather poems, third person past tense only, and not having the assistant narrate the user's character's thoughts. | This establishes critical style constraints and a boundary on narrative control that will affect future outputs.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 198 tok · 12.2 s · run `lvdo9j11`

```text
ID | DECISION | MEMORY | REASON
D1 | HOLD | [user preference]: "I like replies in short paragraphs." (Marked as held until later confirmation or session reset, pending user feedback on length consistency) | The user initially stated a rule without explicit instruction to be confirmed; the assistant acknowledged understanding but did not finalize acceptance yet.
D2 | HOLD | [user preference]: "No purple prose, no weather poems." (Marked temporarily held as per general instruction set regarding style changes vs. specific stylistic rules, pending confirmation that the user is aware of their own stated constraints) | Similar to D1, this constraint was noted in the initial exchange but not explicitly confirmed as the final shared ground truth before any actual content divergence begins.
NONE | NONE | No other facts were recorded in the conversation (except the two preferences above). | The provided dialogue only contained instructions regarding writing style, with no additional characters or specific plot details to store yet.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 37 tok · 2.2 s · run `g4ms5gk9`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | User preference | User prefers shorter narrative segments and avoids detailed internal monologue narration in this roleplay mode.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 100 tok · 3.6 s · run `4x6m1xv8`

```text
C1  
MEMORY: preferences for brevity  
REASON: concise responses are more useful long-term.

C2  
MEMORY: world rules about dialogue flow  
REASON: maintaining narrative continuity is key.

C3  
MEMORY: claims uncertain about future details  
REASON: need to clarify ambiguous statements.

C4  
MEMORY: user expressed liking short replies  
REASON: aligns with preference for efficiency.
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 111 tok · 3.8 s · run `9qox7nau`

```text
C1  
MEMORY: preferences for brevity  
REASON: concise responses are more useful long-term.

C2  
MEMORY: world rules about clarity  
REASON: maintaining consistency across sessions.

C3  
MEMORY: claims uncertain  
REASON: need to verify details before storing.

C4  
MEMORY: relationships between characters  
REASON: important for future interactions.

C5  
MEMORY: user's stated likes  
REASON: preserving personal preferences.
```
