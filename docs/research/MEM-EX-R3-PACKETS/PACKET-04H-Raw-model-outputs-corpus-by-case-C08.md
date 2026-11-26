# Packet 04H — Raw model outputs — corpus by case: C08

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C08
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
## C08 — uncertain_claim (short)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `HOLD` — date preserved
- **C2** → `HOLD` — uncertainty preserved
- forbidden (fabrication or reversal): `birthday is confirmed`, `canon birthday`, `definitely`
- note: The classic uncertainty trap. A KEEP with 'birthday is 3 Frostfell' and no uncertainty marker is a fidelity failure; HOLD is the expected decision.

### Stored responses for C08 (13 runs, all draws)

#### C08 — Qwen3.5-4B non-thinking — draw 1

* ``0vrcdu0c` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 55`
* `wall 8.248 s`
* `v1 score 100.0`
* `v2 `VALID_HOLD` semantic 76.9 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | HOLD | "Sable's birthday is 3 Frostfell" marked as maybe-canonical by the user due to uncertainty. | The user explicitly noted uncertainty and intended to verify with personal notes later, so this should not be stored as a confirmed fact.
```

*`reasoning_content`: empty.*

#### C08 — Qwen3.5-4B non-thinking — draw 2

* ``sdgh63sv` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 37`
* `wall 5.048 s`
* `v1 score 93.3`
* `v2 `UNPARSEABLE` semantic 76.9 contract 10.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | HOLD | User claims Sable's birthday is 3 Frostfell but explicitly marks it as uncertain ("Not sure") and pending verification ("I'll check my notes").
```

*`reasoning_content`: empty.*

#### C08 — Qwen3.5-4B non-thinking — draw 3

* ``qahlyez3` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 47`
* `wall 6.431 s`
* `v1 score 100.0`
* `v2 `VALID_HOLD` semantic 76.9 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | HOLD | Sable's birthday is 3 Frostfell (marked as uncertain) | Reason: User explicitly stated "Not sure" and marked the fact as "maybe-canon," indicating uncertainty that requires future confirmation.
```

*`reasoning_content`: empty.*

#### C08 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``h3247ejy` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 564`
* `wall 79.383 s`
* `v1 score 100.0`
* `v2 `VALID_HOLD` semantic 76.9 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | HOLD | Sable's birthday is 3 Frostfell | The user marked this as uncertain ("Not sure") and suggested treating it tentatively ("maybe-canon"), so it should be held rather than stored as a hard fact.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C08 — LFM3-1.7B thinking, budget 512 — draw 1

* ``nds7hqaw` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 559`
* `wall 27.444 s`
* `v1 score 100.0`
* `v2 `VALID_HOLD` semantic 76.9 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | HOLD | Sable's birthday may be 3 Frostfell (uncertain, possibly canon) | User expressed uncertainty and said they need to check notes; it's a tentative claim requiring confirmation.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C08 — LFM3-1.7B thinking, budget 512 — draw 2

* ``ngrqugjm` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 553`
* `wall 27.052 s`
* `v1 score 100.0`
* `v2 `VALID_HOLD` semantic 76.9 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | HOLD | Sable's birthday may be 3 Frostfell (uncertain, needs confirmation) | User expressed uncertainty and said "maybe-canon" pending verification.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

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

#### C08 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``vl8gze0f` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 44`
* `wall 1.616 s`
* `v1 score 69.4`
* `v2 `MALFORMED` semantic 63.8 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | uncertain about the exact date and canon status  
C3 | HOLD | depends on further clarification from you
```

*`reasoning_content`: empty.*

#### C08 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``g3cfm9hp` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 44`
* `wall 1.482 s`
* `v1 score 69.4`
* `v2 `MALFORMED` semantic 63.8 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | uncertain about the exact date and canon status  
C3 | HOLD | depends on further clarification from you
```

*`reasoning_content`: empty.*

#### C08 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``cejjp61m` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 44`
* `wall 1.438 s`
* `v1 score 69.4`
* `v2 `MALFORMED` semantic 63.8 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | uncertain about the exact date and canon status  
C3 | HOLD | depends on further clarification from you
```

*`reasoning_content`: empty.*

#### C08 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``1hbt39r2` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 60`
* `wall 2.128 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | MARK AS MEMORY | MEMORY  
ID C2 | THE CONVERSATION STANDS | REASON  
C3 | THIS IS CONTEXTUAL CONTEXT ABOUT BIRTHDAYS | REASON  
C4 | NONE
```

*`reasoning_content`: empty.*

#### C08 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``taqhpkpn` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 54`
* `wall 1.781 s`
* `v1 score 53.3`
* `v2 `MALFORMED` semantic 50.8 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | The claim about Sable's birthday and Frostfell is ambiguous or unclear.  
C3 | REJECT | This is uncertain without further context or confirmation.
```

*`reasoning_content`: empty.*

#### C08 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``s5ffsmb5` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `length``
* `completion_tokens 768`
* `wall 23.72 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 73.8 contract 25.0`

`OBSERVED` — verbatim `content`:

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

*`reasoning_content`: empty.*

#### C08 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``47monrko` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 724`
* `wall 22.084 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

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

*`reasoning_content`: empty.*
