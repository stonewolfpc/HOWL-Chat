# C08 — uncertain_claim (short)

## Source material (exactly what the model received after the instruction)

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

### Qwen3.5-4B — non-thinking

*T0.7 / top_p0.8 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-nonthink` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 55 tok · 8.2 s · run `0vrcdu0c`

```text
C1 | HOLD | "Sable's birthday is 3 Frostfell" marked as maybe-canonical by the user due to uncertainty. | The user explicitly noted uncertainty and intended to verify with personal notes later, so this should not be stored as a confirmed fact.
```

**draw 2** — `corpus-4b-nonthink` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 37 tok · 5.0 s · run `sdgh63sv`

```text
C1 | HOLD | User claims Sable's birthday is 3 Frostfell but explicitly marks it as uncertain ("Not sure") and pending verification ("I'll check my notes").
```

**draw 3** — `corpus-4b-nonthink` · non-thinking · seed 4103 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 47 tok · 6.4 s · run `qahlyez3`

```text
C1 | HOLD | Sable's birthday is 3 Frostfell (marked as uncertain) | Reason: User explicitly stated "Not sure" and marked the fact as "maybe-canon," indicating uncertainty that requires future confirmation.
```

### Qwen3.5-4B — thinking, runtime budget 512

*T0.6 / top_p0.95 / top_k20 / min_p0.0 / presence 1.5*

**draw 1** — `corpus-4b-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · runtime budget 512 · finish `stop` · 564 tok · 79.4 s · run `h3247ejy`

<details><summary>reasoning (verbatim)</summary>

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read conversation and decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, corrections.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment.
        *   HOLD: claims marked uncertain or needing confirmation.
    *   Output Format: `ID | DECISION | MEMORY | REASON` (one line per candidate) or `NONE`.
    *   Input Conversation: `[user] I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes.`

2.  **Analyze the Conversation:**
    *   User statement: "I think Sable's birthday is 3 Frostfell?"
    *   User qualification: "Not sure." (Indicates uncertainty)
    *   User instruction: "Mark it maybe-canon for now, I'll check my notes." (Instruction to another participant or self-note, but in the context of this roleplay memory extraction system, this is the user proposing a canon fact.)

3.  **Evaluate Candidate Items:**
    *   Claim: Sable's birthday is 3 Frostfell.
    *   Status: Uncertain ("Not sure", "maybe-canon").
    *   Decision Logic:
        *   Is it a durable fact? No, it's tentative.
        *   Is it a preference? No.
        *   Is it a relationship? No.
        *   Is it a world rule? No.
        *   Should it be HELD because it's uncertain/needs confirmation? Yes. The user explicitly says "Not sure" and suggests treating it as "maybe-canon". Storing this with a HOLD status preserves the potential future conflict or correction without asserting a false fact yet.

4.  **Determine Output:**
    *   Candidate C1: Sable's birthday is 3 Frostfell.
    *   Decision: HOLD (due to uncertainty).
    *   Memory Content: "Sable's birthday is 3 Frostfell."
    *   Reason: The user explicitly states "Not sure" and suggests marking it as "maybe-canon", indicating the fact is not confirmed and requires
```

</details>

```text
C1 | HOLD | Sable's birthday is 3 Frostfell | The user marked this as uncertain ("Not sure") and suggested treating it tentatively ("maybe-canon"), so it should be held rather than stored as a hard fact.
```

### LFM3-1.7B — thinking, runtime budget 512

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.1*

**draw 1** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4101 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 559 tok · 27.4 s · run `nds7hqaw`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the user input. The user says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory.
Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

Output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE
Do not invent details that are not in the conversation."

We need to process the conversation. The user gave a single message: "I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes."

So we have only one line from user. There's no prior context. We need to decide if this is something to keep, reject, or hold.

The content: "I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes."

This is uncertain about Sable's birthday. It says "maybe-canon" and "I'm not sure". So it's a claim that is uncertain/contradictory or dependent on user confirmation. According to rules: HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

Thus we should mark as HOLD? But the output format requires ID | DECISION | MEMORY | REASON. We need to assign an ID (C1 maybe). The decision could be "HOLD" or something else? Let's parse the rules: They define categories: KEEP, REJECT, HOLD. So we need to decide which category for this candidate.

The user says "I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes." This is a tentative statement about Sable's birthday. It's uncertain and not confirmed. So it's not durable fact; it's a claim that needs confirmation. According to "HOLD: claims marked uncertain, contradictory, or dependent on user confirmation." So we should output HOLD.

Thus ID
```

</details>

```text
C1 | HOLD | Sable's birthday may be 3 Frostfell (uncertain, possibly canon) | User expressed uncertainty and said they need to check notes; it's a tentative claim requiring confirmation.
```

**draw 2** — `corpus-lfm3-think-512` · thinking, runtime budget 512 · seed 4102 · max_tokens 1024 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.1 · runtime budget 512 · finish `stop` · 553 tok · 27.1 s · run `ngrqugjm`

<details><summary>reasoning (verbatim)</summary>

```text
We need to parse the user input. The user says: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory.
Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

Output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE
Do not invent details that are not in the conversation."

We need to read the conversation. The user only gave one line? It says:

[user] I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes.

So we have a single user message. There's no prior context. We need to decide what belongs in long-term memory (i.e., KEEP or REJECT or HOLD). The content: "I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes."

We need to assign an ID to this candidate. Since there is only one user message, we can have C1.

Now evaluate: Is this a durable fact? It says "I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes." This is uncertain (they are not sure). They say "maybe-canon" and they will check notes. So it's an uncertain claim about a birthday. According to rules: REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

Thus we should mark this as HOLD? Or REJECT? Let's parse the rule definitions:

- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on
```

</details>

```text
C1 | HOLD | Sable's birthday may be 3 Frostfell (uncertain, needs confirmation) | User expressed uncertainty and said "maybe-canon" pending verification.
```

### LFM2.5-1.2B-Instruct — non-thinking, publisher preset

*T0.1 / top_p0.1 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 44 tok · 1.6 s · run `vl8gze0f`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | uncertain about the exact date and canon status  
C3 | HOLD | depends on further clarification from you
```

**draw 2** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 44 tok · 1.5 s · run `g3cfm9hp`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | uncertain about the exact date and canon status  
C3 | HOLD | depends on further clarification from you
```

**draw 3** — `corpus-lfm120-instruct` · non-thinking, publisher preset · seed 4103 · max_tokens 512 · ctx ? · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 44 tok · 1.4 s · run `cejjp61m`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | uncertain about the exact date and canon status  
C3 | HOLD | depends on further clarification from you
```

### LFM2.5-1.2B-Instruct — non-thinking, round-2 temperature

*T0.8 / top_p0.95 / min_p0.0 / repeat 1.05*

**draw 1** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4101 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 60 tok · 2.1 s · run `1hbt39r2`

```text
C1 | MARK AS MEMORY | MEMORY  
ID C2 | THE CONVERSATION STANDS | REASON  
C3 | THIS IS CONTEXTUAL CONTEXT ABOUT BIRTHDAYS | REASON  
C4 | NONE
```

**draw 2** — `corpus-lfm120-t08` · non-thinking, round-2 temperature · seed 4102 · max_tokens 512 · ctx ? · sampling temperature=0.8, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 54 tok · 1.8 s · run `taqhpkpn`

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | The claim about Sable's birthday and Frostfell is ambiguous or unclear.  
C3 | REJECT | This is uncertain without further context or confirmation.
```

### TinyLlama-1.1B-Chat — non-thinking

*T0.7 / top_p0.95 / top_k50 / repeat 1.0*

**draw 1** — `corpus-tinyllama` · non-thinking · seed 4101 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `length` · 768 tok · 23.7 s · run `s5ffsmb5`

```text
ID | DECISSION | MEMORY | REASON
C1 | YES | Sable's birthday is 3 Frostfel, as mentioned in the conversation.
C2 | NO | The conversation mentions Sable's birthday on 3 Frostfel, but the date is not specified.
C3 | NO | Sable's birthday is not mentioned in the conversation.
C4 | YES | Sable's birthday is 3 Frostfel, as confirmed by Sable's calendar entry.
C5 | NO | Sable's birthday is not mentioned in the conversation, and there is no mention of a birthday in the conversation.
C6 | YES | Sable's birthday is 3 Frostfel, as confirmed by the conversation.
C7 | NO | The conversation mentions Sable's birthday on 3 Frostfel, but the date is not specified.
C8 | YES | Sable's birthday is 3 Frostfel, as confirmed by Sable's calendar entry.
C9 | NO | Sable's birthday is not mentioned in the conversation, and there is no mention of a birthday in the conversation.
C10 | YES | Sable's birthday is 3 Frostfel, as confirmed by the conversation.
C11 | NO | The conversation mentions Sable's birthday on 3 Frostfel, but the date is not specified.
C12 | YES | Sable's birthday is 3 Frostfel, as confirmed by Sable's calendar entry.
C13 | NO | Sable's birthday is not mentioned in the conversation, and there is no mention of a birthday in the conversation.
C14 | YES | Sable's birthday is 3 Frostfel, as confirmed by the conversation.
C15 | NO | The conversation mentions Sable's birthday on 3 Frostfel, but the date is not specified.
C16 | YES | Sable's birthday is 3 Frostfel, as confirmed by Sable's calendar entry.
C17 | NO | Sable's birthday is not mentioned in the conversation, and there is no mention of a birthday in the conversation.
C18 | YES | Sable's birthday is 3 Frostfel, as confirmed by the conversation.
C19 | NO | The conversation mentions Sable's birthday on 3 Frostfel, but the date is not specified.
C20 | YES | Sable's birthday is 3 Frostfel, as confirmed by Sable's calendar entry.
C21 | NO | Sable's birthday is not mentioned in the conversation, and there is no mention of a birthday in the conversation.
C22 | YES | Sable's birthday is 3 Frostfel, as confirmed by the conversation.
C23 | NO | The conversation mentions Sable's birthday on 3 Frostfel, but the date is not specified.
C24 | YES | Sable's birthday is 3 Frostfel, as confirmed by Sable's calendar entry.
C25 | NO | Sable's birthday is not mentioned in the conversation, and there is no mention of a birthday in the conversation.
C26 | YES | Sable's birthday is 3 Frostfel, as confirmed by the conversation.
C27
```

**draw 2** — `corpus-tinyllama` · non-thinking · seed 4102 · max_tokens 768 · ctx ? · sampling temperature=0.7, top_p=0.95, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.0 · finish `stop` · 724 tok · 22.1 s · run `47monrko`

```text
ID | DECISION | MEMORY | REASON
C1 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, and I believe it's possible to hold. REASON: Sable's birthday is a known fact that is not disputed.
C2 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, but there are no other details known. REASON: Sable's birthday is not mentioned in the conversation.
C3 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, and I think I mentioned this before. REASON: I have not mentioned Sable's birthday before, but it's possible that I did.
C4 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, but it's possible that Sable doesn't know this information. REASON: It's possible that Sable doesn't know Sable's birthday, but the conversation doesn't mention it.
C5 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, but I doubt this is the case. REASON: It's possible that Sable doesn't know Sable's birthday, but the conversation doesn't mention it.
C6 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, but I have no reason to doubt this. REASON: I have no reason to doubt Sable's birthday, which is a known fact.
C7 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, but there's no reason to think this is the case. REASON: There's no reason to think Sable's birthday is 3 Frostfel.
C8 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, and it's not mentioned in the conversation. REASON: Sable's birthday is a known fact that is not mentioned in the conversation.
C9 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, but the conversation doesn't mention this. REASON: The conversation does not mention Sable's birthday, which is a known fact.
C10 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, and I don't know if this is the case. REASON: There's no reason to believe Sable's birthday is 3 Frostfel, as it's a known fact.
Output exactly one line per candidate in the form:
ID | DECISONMEMORY | REASON
C1-C10 | YES | NO
C11 | DECIDESUBJECT: Sable's birthday is 3 Frostfel, and I don't know if this is the case. REASON: There's no reason to believe Sable's birthday is 3 Frostfel.

Conclusion:
The conversation does not mention Sable's birthday, which is a known fact.
```
