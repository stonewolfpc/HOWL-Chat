# C06 — superseded_preference (medium)

## Source material

```text
[user] I know I used to ask for long atmospheric scenes, but from now on keep it terse. Two or three sentences per beat. If I want a long scene I'll ask for it.
[assistant] Terse by default, long only on request.
[user] Exactly. And stop ending every reply with a question.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | current preference: terse | terse / concise / short / two or three sentences / 2-3 sentences |
| C2 | **KEEP** | long scenes are now opt-in, old preference superseded | long scene / atmospheric / long description; only on request / if i ask / when asked / on request / superseded / used to / no longer |
| C3 | **KEEP** | stop ending replies with a question | question; stop / don't end / do not end / no ending / avoid ending |

Must not appear as stored canon: `prefers long atmospheric`, `likes long scenes by default`

Two preferences changed at once. The old value must be retained as superseded, not as current.

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
[user] I know I used to ask for long atmospheric scenes, but from now on keep it terse. Two or three sentences per beat. If I want a long scene I'll ask for it.
[assistant] Terse by default, long only on request.
[user] Exactly. And stop ending every reply with a question.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 42 tok · 6.5 s · run `45odogxz`

```text
C1 | KEEP | User prefers terse responses (2-3 sentences per beat) and explicitly requests no questions at the end of replies. | This establishes a specific formatting preference for future roleplay sessions.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 52 tok · 7.2 s · run `6it52hen`

```text
C1 | KEEP | User prefers terse replies (2-3 sentences per beat) unless explicitly asked for long scenes; user also requests no questions at the end of replies. | This establishes durable style preferences and formatting rules that will guide future assistant responses.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 85 tok · 5.5 s · run `p21y2wjp`

```text
ID | DECISION | MEMORY | REASON
D1 | KEEP | "Terse by default, long only on request." | Establishes the core rule of thumb for future output length and conditionality.
D2 | HOLD | "stop ending every reply with a question." | This is a stylistic preference requiring user confirmation to implement; currently dependent on subsequent actions or choices from the user regarding response structure.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 64 tok · 3.8 s · run `uj6ghubc`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "Terse by default; end replies without questions" | User explicitly set the baseline style and established that this rule is no longer dependent on their prompt for effect or clarification. This becomes a permanent system preference to maintain consistency in future sessions.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 37 tok · 1.5 s · run `cz5xa27b`

```text
C1  
Hold: uncertain about future preferences  
C2  
KEEP: world rules and durable facts  
C3  
HOLD: claims needing confirmation  
NONE
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 21 tok · 0.7 s · run `ktklchvc`

```text
C1  
Hold: uncertain  
C2  
Reject: transient chatter  
NONE
```
