# Packet 02H-2 — Corpus case C08: responses (part 2 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 3 stored responses
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

3 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `47monrko` — `corpus-tinyllama`, draw 2

* Stored model identifier: `corpus-tinyllama`
* Model file: `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf`
* Stored base name: `tinyllama-11b-chat`
* Case: C08
* Draw: 2
* Source: `results/corpus.json`, array index 214 (0-based)
* Also present in: `results/corpus_tinyllama.json`

`OBSERVED` — stored `content`, verbatim:

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

*Stored `reasoning_content` is the empty string.*

#### `89i18iwr` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C08
* Draw: 1
* Source: `results/control_corpus_2b.json`, array index 14 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
1 | HOLD | Sable's birth year is unknown/uncertain until user confirmation. | The user explicitly marked the date as suspect and asked to verify against notes; thus, no definitive fact is established for permanent retention.
```

*Stored `reasoning_content` is the empty string.*

#### `fjqmajpi` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `control_corpus_2b`
* Case: C08
* Draw: 2
* Source: `results/control_corpus_2b.json`, array index 15 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | "SABLE'S BIRTHDAY IS 3 FROSTFELL" + "MARKS IT MAYBE-CANON FOR NOW" | To ensure core canon information regarding Sable's specific birthday is preserved for future session updates, while the user's caveat to check notes allows correction if their research differs.
```

*Stored `reasoning_content` is the empty string.*
