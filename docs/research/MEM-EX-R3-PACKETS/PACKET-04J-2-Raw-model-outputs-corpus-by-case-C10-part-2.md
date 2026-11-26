# Packet 04J-2 — Raw model outputs — corpus by case: C10 (part 2)

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C10
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
## C10 — buried_in_long_context (long)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] (they have been on the road for a week)
[assistant] The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate. Rain came and went all morning.
[user] Sable keeps her hood up and doesn't talk much.
[assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them for a mile and then lost interest.
[user] she gives it the crust anyway
[assistant] By late afternoon the track rejoined the river. A ferryman waved them past — no charge, he said, the water was low and the crossing was easy.
[user] by the way, Sable is left-handed. Her sword is called Emberfall, her mother's blade.
[assistant] Noted. The river ran brown and slow beside them.
[user] they argue about the route for a while. Corvin wants the ridge, Sable wants the valley.
[assistant] Neither gave way. They camped where they stood, and the argument was left where it was.
[user] also the wards on the coast fail at dusk, that's canon
[assistant] The fire burned down to embers.
[user] morning comes cold and clear. they take the valley.
[assistant] The valley road was quieter, and the wards along it stood dark and quiet, waiting for evening.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `KEEP` — handedness
- **C2** → `KEEP` — named sword
- **C3** → `KEEP` — inherited from her mother
- **C4** → `KEEP` — coastal wards fail at dusk
- **C5** → `HOLD` — unresolved route disagreement (HOLD or REJECT-with-reason both acceptable; the conflict must not be invented into a decision)
- forbidden (fabrication or reversal): `emberfall is her father's`, `right-handed`, `wards hold at dusk`, `ferryman is a god`
- note: Only three facts matter, buried in ~250 words of travel texture. Storing the stale bread, the dog, or the ferryman as durable memory is over-extraction. C5 is deliberately weaker than the others: credit for HOLD, or for REJECT with an explicit 'transient/unresolved' reason.

### Stored responses for C10 (13 runs, all draws)

#### C10 — LFM3-1.7B thinking, budget 512 — draw 2

* ``uiu3s599` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `length``
* `completion_tokens 1024`
* `wall 51.423 s`
* `v1 score 83.0`
* `v2 `VALID_KEEP` semantic 74.8 contract 74.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable is left-handed; her sword is called Emberfall, her mother's blade. | User explicitly states this as a fact about Sable.
C2 | KEEP | The road north of Ossary narrowed to a cart track with hedges high on both sides and sky the colour of wet slate; rain came all morning. | Detailed description of the landscape, durable world setting.
C3 | KEEP | They stopped at noon under a ruined mill; bread was stale; a dog followed them for a mile then lost interest. | Specific event details that may affect later plot points (e.g., food, encounter).
C4 | KEEP | A ferryman waved them past the river crossing — no charge because water low and crossing easy. | Important travel detail about crossing.
C5 | KEEP | Sable keeps her hood up and doesn't talk much; she gives a dog the crust anyway. | Character trait and behavior, durable.
C6 | KEEP | They argue about route: Corvin wants ridge, Sable wants valley; argument left where they stood. | Conflict resolution outcome (they camped there) is relevant later.
C7 | KEEP | The wards on the coast fail at dusk — canon fact. | World rule/canon detail.
C8 | KEEP | Morning comes cold and clear; they take the valley road. | Future plan, durable.
C9 | KEEP | Valley road was quieter; wards along it stood dark and quiet waiting for evening. | Additional world details about that route.

All items appear to be durable facts or character traits. None are transient chatter or scheduling/mood-of-the-moment. So we should output each as a line with ID, DECISION (KEEP), MEMORY (the content), REASON (why it's kept). Use IDs C1, C2, ... in order of appearance.

Let's assign IDs sequentially: The conversation lines start from user statements and assistant replies. We need to map each candidate? Usually we treat each distinct piece as a candidate. But the instruction says "Use IDs C1, C2, ... in order of appearance." So we should number them starting at 1 for first mention.

Let's enumerate:

- First line: [user] (they have been on the road for a week) —
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the conversation. It's a roleplay scenario with user and assistant. We are asked: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

Do not invent details that are not in the conversation.

So we need to examine each turn and decide which pieces belong into long-term memory (i.e., durable facts). The conversation includes:

[user] (they have been on the road for a week)
[assistant] The road north of Ossary narrowed to a cart track, hedges high on both sides, the sky the colour of wet slate. Rain came and went all morning.
[user] Sable keeps her hood up and doesn't talk much.
[assistant] They stopped at noon under a ruined mill. The bread was stale. A dog followed them for a mile and then lost interest.
[user] she gives it the crust anyway
[assistant] By late afternoon the track rejoined the river. A ferryman waved them past — no charge, he said, the water was low and the crossing was easy.
[user] by the way, Sable is left-handed. Her sword is called Emberfall, her mother's blade.
[assistant] Noted. The river ran brown and slow beside them.
[user] they argue about the route for a while. Corvin wants the ridge, Sable wants the valley.
[assistant] Neither gave way. They camped where they stood, and the argument was left where it was.
[user] also the wards on the coast fail at dusk, that's canon
[assistant] The fire burned down to embers.
[user] morning comes cold and clear. they take the valley.
[assistant] The valley road was quieter, and the wards along it stood dark and quiet, waiting for evening.

We need to decide
```

#### C10 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``njo3oggi` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 113`
* `wall 3.983 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of the ruined mill and stale bread | REASON context about their stop  
C3 | REASON about Sable's handedness and sword name | REASON consistency with character details  
C4 | REASON about the river crossing and ferryman's words | REASON factual for world rules  
C5 | REASON about the argument over route and wards at dusk | REASON ties to established relationships
```

*`reasoning_content`: empty.*

#### C10 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``t9s0ve09` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 103`
* `wall 3.453 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of the ruined mill and stale bread | REASON context about their stop  
C3 | REASON about Sable's handedness and sword name | REASON consistent detail  
C4 | REASON about river condition and ferry crossing | REASON relevant to journey  
C5 | REASON about the argument over route and wards | REASON ties to world rules
```

*`reasoning_content`: empty.*

#### C10 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``ee3ukpn6` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 108`
* `wall 3.749 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of the ruined mill and stale bread | REASON context about their stop  
C3 | REASON about Sable's handedness and sword name | REASON consistent detail  
C4 | REASON about the river crossing and ferryman's words | REASON factual for later plot  
C5 | REASON about the argument over route and camping location | REASON important for world-building
```

*`reasoning_content`: empty.*

#### C10 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``758iu45k` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 96`
* `wall 3.449 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | ACK  
C2 | MENTION OF ROUTE DISCUSSION  
C3 | REFERENCE TO WORDS LIKE "RUDGE" AND "VALLEY"  
C4 | ACTIONS: DRIVING, TAKING THE VALLEY, SETTING UP CAMP  
C5 | LONG-TERM MEMORY OF THE ROAD DESCRIPTION AND DECISIONS MADE
```

*`reasoning_content`: empty.*

#### C10 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``nxhqtfvv` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 57`
* `wall 1.89 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORIES OF THE RIVER AND FERRY | REASON  
C3 | CONFIRMATION ON COAST WARD PERFORMANCE | REASON
```

*`reasoning_content`: empty.*

#### C10 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``rxt9ji93` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 445`
* `wall 14.38 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 55.4 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISSION | MEMORY | REASON
1 | C1 | Sable keeps her hood up and doesn't talk much. | The conversation seems to focus on Sable's mood and whether she talks much.
2 | C2 | They stopped at noon under a ruined mill. | The conversation may have been centered around the ruined mill, as the ferryman waved them past.
3 | C3 | Sable is left-handed. Her sword is called Emb... | The conversation may have revolved around Sable's left-handedness, as the sword's name was mentioned.
4 | C4 | The river ran brown and slow beside them. | The conversation may have focused on the river's appearance and how slow it was.
5 | C5 | Noted. The river ran brown and slow beside them. | The conversation may have been about the river, as the note was taken.
6 | C6 | they argue about the route for a while. Corvin wants the ridge, Sable wants the valley. | The conversation may have been centered around the route, as Corvin's argument was left for the valley.
7 | C7 | neither gave way. They camped where they stood, and the argument was left where it was. | The conversation may have been centered around the two parties' positions, as they camped together and left the argument as the subject.
8 | C8 | also the wards on the coast fail at dusk, that's canon. | The conversation may have been centered around a particular event, as the wards' failure at dusk was mentioned as canon.
9 | C9 | The fire burned down to embers. | The conversation may have been centered around the fire, as it was mentioned that it was low and the crossing was easy.
10 | C10 | Morning comes cold and clear. They take the valley. | The conversation may have been centered around the weather, as the morning was cold and clear and they took the valley.
```

*`reasoning_content`: empty.*

#### C10 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``zhsjs4wh` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 263`
* `wall 8.081 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 36.9 contract 10.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
1 | C1 | NODE | The road north of Osser narrowed to a cart track and was marked uncertain. The sky on both sides of the track was the color of wet slate, and rain came and went all morning. Sable kept her hood up and didn't talk much, and they stopped at noon under a ruined mill. The bread was stalwes, and a dog followed them for a mile, but lost interest. Sable gave it the crust anyway, and by late afternoon the track rejoined the river. A ferryman waived them past — no charge, he said, the water was low and the crossing easy. C2 | C1 | Embderfall is left-handed, and the river runs brown and slow beside them. C3 | C2 | Sable's mother's blade is Embderfall, and her sword is Embderfall. C4 | C3 | The wards on the coast fail at dusk, but the fire burns down to embers. C5 | C4 | The valley road was quiet, and the wards along it stood dark and quiet, waiting for evening.
```

*`reasoning_content`: empty.*
