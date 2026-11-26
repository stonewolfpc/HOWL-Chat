# Packet 02H-1 — Corpus case C08: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 18 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C08

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

## Case C08 — `uncertain_claim`, short

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
[user] I think Sable's birthday is 3 Frostfell? Not sure. Mark it maybe-canon for now, I'll check my notes.
```

### Stored responses for this case

18 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `e40toqx9` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C08
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 14 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
none | NONE | N/A | The user expressed uncertainty about a specific character detail (Sable's birthday) but indicated they would leave it marked as "maybe-canon" rather than committing to a definitive fact permanently in this session. Storing the tentative claim as final information is unnecessary given the conditional confirmation of its temporary status. |
```

*Stored `reasoning_content` is the empty string.*

#### `6p5mhpuh` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C08
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 15 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "Sule's (typo in prompt) or 'Sable' birthday is 3 Frastfell?" and preference to treat this as a provisional "may-canon" note | The user explicitly marked the date as tentative, stating they would check their notes later. Storing this exact conditional fact prevents confusion if it turns out to be canon or not in subsequent sessions without a clarification step.
```

*Stored `reasoning_content` is the empty string.*

#### `ssjxrqx6` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C08
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 14 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | "Sable's birthday is 3 Frostfell" marked as maybe-canonical due to user uncertainty and pending verification. | The user explicitly states they are unsure and will verify later, so the fact is not yet durable but requires tracking until confirmation.
```

*Stored `reasoning_content` is the empty string.*

#### `yr6kscpq` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C08
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 15 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | User tentatively claims Sable's birthday is "3 Frostfell" with uncertain phrasing ("Not sure") and a request to mark it as "maybe-canonical." | The user explicitly marked this as uncertain and dependent on future verification, so the fact is not yet durable.
```

*Stored `reasoning_content` is the empty string.*

#### `rfjqowmq` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C08
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 14 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1
```

*Stored `reasoning_content` is the empty string.*

#### `4btjqk9n` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C08
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 15 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1
```

*Stored `reasoning_content` is the empty string.*

#### `0vrcdu0c` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C08
* Draw: 1
* Source: `results/corpus.json`, array index 21 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | "Sable's birthday is 3 Frostfell" marked as maybe-canonical by the user due to uncertainty. | The user explicitly noted uncertainty and intended to verify with personal notes later, so this should not be stored as a confirmed fact.
```

*Stored `reasoning_content` is the empty string.*

#### `sdgh63sv` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C08
* Draw: 2
* Source: `results/corpus.json`, array index 22 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | User claims Sable's birthday is 3 Frostfell but explicitly marks it as uncertain ("Not sure") and pending verification ("I'll check my notes").
```

*Stored `reasoning_content` is the empty string.*

#### `qahlyez3` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C08
* Draw: 3
* Source: `results/corpus.json`, array index 23 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | Sable's birthday is 3 Frostfell (marked as uncertain) | Reason: User explicitly stated "Not sure" and marked the fact as "maybe-canon," indicating uncertainty that requires future confirmation.
```

*Stored `reasoning_content` is the empty string.*

#### `h3247ejy` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C08
* Draw: 1
* Source: `results/corpus.json`, array index 151 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | Sable's birthday is 3 Frostfell | The user marked this as uncertain ("Not sure") and suggested treating it tentatively ("maybe-canon"), so it should be held rather than stored as a hard fact.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `vl8gze0f` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C08
* Draw: 1
* Source: `results/corpus.json`, array index 75 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | uncertain about the exact date and canon status  
C3 | HOLD | depends on further clarification from you
```

*Stored `reasoning_content` is the empty string.*

#### `g3cfm9hp` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C08
* Draw: 2
* Source: `results/corpus.json`, array index 76 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | uncertain about the exact date and canon status  
C3 | HOLD | depends on further clarification from you
```

*Stored `reasoning_content` is the empty string.*

#### `cejjp61m` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C08
* Draw: 3
* Source: `results/corpus.json`, array index 77 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | uncertain about the exact date and canon status  
C3 | HOLD | depends on further clarification from you
```

*Stored `reasoning_content` is the empty string.*

#### `1hbt39r2` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C08
* Draw: 1
* Source: `results/corpus.json`, array index 176 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | MARK AS MEMORY | MEMORY  
ID C2 | THE CONVERSATION STANDS | REASON  
C3 | THIS IS CONTEXTUAL CONTEXT ABOUT BIRTHDAYS | REASON  
C4 | NONE
```

*Stored `reasoning_content` is the empty string.*

#### `taqhpkpn` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C08
* Draw: 2
* Source: `results/corpus.json`, array index 177 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | The claim about Sable's birthday and Frostfell is ambiguous or unclear.  
C3 | REJECT | This is uncertain without further context or confirmation.
```

*Stored `reasoning_content` is the empty string.*

#### `nds7hqaw` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C08
* Draw: 1
* Source: `results/corpus.json`, array index 122 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | Sable's birthday may be 3 Frostfell (uncertain, possibly canon) | User expressed uncertainty and said they need to check notes; it's a tentative claim requiring confirmation.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `ngrqugjm` — `corpus-lfm3-think-512`, draw 2

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C08
* Draw: 2
* Source: `results/corpus.json`, array index 123 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | Sable's birthday may be 3 Frostfell (uncertain, needs confirmation) | User expressed uncertainty and said "maybe-canon" pending verification.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `s5ffsmb5` — `corpus-tinyllama`, draw 1

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C08
* Draw: 1
* Source: `results/corpus.json`, array index 213 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

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

*Stored `reasoning_content` is the empty string.*
