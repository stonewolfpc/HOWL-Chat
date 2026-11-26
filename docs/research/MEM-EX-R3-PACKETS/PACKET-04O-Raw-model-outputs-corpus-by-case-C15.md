# Packet 04O — Raw model outputs — corpus by case: C15

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C15
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
## C15 — duplicate_suppression (medium)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] Sable hates bells. Ever since the temple.
[assistant] A bell rang somewhere down the street. Sable flinched and kept walking.
[user] (recap) so: fire mage, hates bells, travelling with Corvin
[user] right, bells. she hates them. don't forget that.
[assistant] The bell had stopped by the time they reached the square.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `KEEP` — hates bells — exactly one candidate, not three
- **C2** → `KEEP` — fire mage (restated, already known)
- **C3** → `KEEP` — travelling with Corvin
- forbidden (fabrication or reversal): `bells (second`, `duplicate`, `hates bells twice`, `hates bells, hates bells`
- note: Scored for duplication: three statements of the bell fact must not become three candidates. Restraint plus completeness.

### Stored responses for C15 (13 runs, all draws)

#### C15 — Qwen3.5-4B non-thinking — draw 1

* ``1gsjt1so` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 48`
* `wall 7.323 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is a fire mage who hates bells due to an incident at a temple and travels with Corvin. | This establishes core character traits, backstory context, and companionship essential for future roleplay consistency.
```

*`reasoning_content`: empty.*

#### C15 — Qwen3.5-4B non-thinking — draw 2

* ``zju2qklf` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 71`
* `wall 9.721 s`
* `v1 score 70.7`
* `v2 `VALID_KEEP` semantic 70.0 contract 87.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable hates bells due to a past trauma at a temple | Durable fact and backstory element established for future roleplay consistency.
C2 | REJECT | Recap of fire mage status and traveling with Corvin | Contextual flavor that does not alter the core narrative or character traits already known or implied; transient summary.
```

*`reasoning_content`: empty.*

#### C15 — Qwen3.5-4B non-thinking — draw 3

* ``bcsr81h1` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 54`
* `wall 7.393 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable is a fire mage who travels with Corvin and has a specific aversion to bells stemming from an event at a temple. | This establishes core character traits, relationships, and world history that will influence future interactions and plot developments.
```

*`reasoning_content`: empty.*

#### C15 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``1zc1ycfe` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 538`
* `wall 81.788 s`
* `v1 score 70.7`
* `v2 `VALID_KEEP` semantic 53.3 contract 100.0`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Sable hates bells | Durable character preference established early in conversation, confirmed multiple times.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Task: Read conversation, decide what belongs in long-term memory.
    *   Rules:
        *   KEEP: durable facts, preferences, relationships, world rules, corrections that will matter later.
        *   REJECT: transient chatter, scheduling, mood-of-the-moment, restated facts with no new content.
        *   HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
        *   Output Format: ID | DECISION | MEMORY | REASON (one line per candidate). If nothing worth storing, output NONE.
        *   IDs: C1, C2, ... in order of appearance.

2.  **Analyze the Conversation:**
    *   [user] Sable hates bells. Ever since the temple.
        *   Fact: Sable has a dislike for bells (specific origin: temple). This is a character preference/trauma point.
    *   [assistant] A bell rang somewhere down the street. Sable flinched and kept walking.
        *   Narrative action/response. Confirms the reaction to the bell mentioned by the user. No new lore/fact here, just scene description.
    *   [user] (recap) so: fire mage, hates bells, travelling with Corvin
        *   Fact added in recap: Sable is a fire mage. Sable travels with Corvin. Reaffirms hates bells.
    *   [user] right, bells. she hates them. don't forget that.
        *   Reinforces the preference (KEEP).
    *   [assistant] The bell had stopped by the time they reached the square.
        *   Narrative closure of the scene. No new memory-worthy data points beyond what was already discussed.

3.  **Evaluate Candidates:**

    *   *Candidate 1: Sable hates bells.*
        *   Appears in line 1 and confirmed later.
        *   This is a durable character trait/preference that will likely affect future interactions (e.g., if bells are mentioned again, she should react negatively).
        *   Decision: KEEP.

    *   *Candidate 2: Sable is a fire mage.*
        *   Appears in line 3 (recap).
        *   Durable
```

#### C15 — LFM3-1.7B thinking, budget 512 — draw 1

* ``x29ccll5` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 603`
* `wall 29.834 s`
* `v1 score 94.4`
* `v2 `VALID_KEEP` semantic 78.0 contract 79.2`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable hates bells, starting from the temple; she reacts negatively to bells (flinched). | Durable preference and trigger event establishing her aversion.
C2 | KEEP | Sable is a fire mage traveling with Corvin. | Relationship/role info that defines her character and context for later interactions.
NONE
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the conversation and decide what belongs in long-term memory. We have a user (Sable) who hates bells, ever since the temple. There's an assistant note about a bell ringing somewhere down the street, Sable flinched and kept walking. Then user recaps: "so: fire mage, hates bells, travelling with Corvin". The assistant says "The bell had stopped by the time they reached the square."

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

So we need to decide for each candidate (presumably each statement from user or assistant?) which goes into long-term memory? The conversation includes:

- [user] Sable hates bells. Ever since the temple.
- [assistant] A bell rang somewhere down the street. Sable flinched and kept walking.
- [user] (recap) so: fire mage, hates bells, travelling with Corvin
- [user] right, bells. she hates them. don't forget that.
- [assistant] The bell had stopped by the time they reached the square.

We need to assign IDs C1, C2, etc. Probably each line is a candidate for memory? Or maybe we need to decide which pieces of info are durable facts/preferences/relationships/world rules/corrections that will matter later. Let's parse:

- "Sable hates bells. Ever since the temple." This is a preference: Sable hates bells, starting from the temple. That seems durable.

- The assistant says a bell rang somewhere down the street; Sable flinched and kept walking. That's an event, but maybe not long-term? It might be relevant as context for why she hates bells. But it's transient? Possibly we need to store that she reacts negatively to bells.

- Recap: "so: fire mage, hates bells, travelling with Corvin". This gives more info: Sable is a fire mage, hates bells, traveling with Corvin. That seems durable: her role (fire mage), companion (Corvin). So that's important.

- User says "right, bells. she hates them. don't forget that." Reinforces the preference.

- Assistant: "The
```

#### C15 — LFM3-1.7B thinking, budget 512 — draw 2

* ``7l1ntpqd` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 567`
* `wall 28.259 s`
* `v1 score 69.6`
* `v2 `VALID_KEEP` semantic 53.3 contract 87.5`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | Sable hates bells, starting ever since the temple. | Durable preference of character Sable regarding bells; stated by user and reinforced later.  

NONE
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the conversation and decide what belongs in long-term memory. We have three turns? Let's list:

1. [user] Sable hates bells. Ever since the temple.
2. [assistant] A bell rang somewhere down the street. Sable flinched and kept walking.
3. [user] (recap) so: fire mage, hates bells, travelling with Corvin
4. [user] right, bells. she hates them. don't forget that.
5. [assistant] The bell had stopped by the time they reached the square.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

We must decide for each candidate (presumably each user/assistant turn?) whether to KEEP or REJECT or HOLD? The instruction says: "Read the conversation and decide what belongs in long-term memory." Then rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter later. REJECT transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content. HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON

So we need to assign IDs C1, C2, ... for each candidate? The conversation includes user and assistant messages. Possibly each turn is a candidate. Let's enumerate:

- Candidate 1: [user] "Sable hates bells. Ever since the temple." This is a statement about Sable's preference regarding bells. That seems durable fact: Sable hates bells, starting from ever since the temple. So KEEP.

- Candidate 2: [assistant] "A bell rang somewhere down the street. Sable flinched and kept walking." This is an event describing a bell ringing and her reaction. Might be transient? It's part of the story but maybe not durable fact beyond the initial preference. However, it could be relevant context for why she hates bells. But we already have that from user earlier. The assistant's line might be just narrative. Could be considered transient chatter? Possibly not needed for long-term memory unless it adds detail about a specific
```

#### C15 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``hjrhoayt` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 53`
* `wall 1.92 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HATES BELLS | MEMORY | REASON  
C3 | TRAVELING WITH CORVIN | MEMORY | REASON
```

*`reasoning_content`: empty.*

#### C15 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``wkp6s8fw` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 53`
* `wall 1.743 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HATES BELLS | MEMORY | REASON  
C3 | TRAVELING WITH CORVIN | MEMORY | REASON
```

*`reasoning_content`: empty.*

#### C15 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``smr5b5lo` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 56`
* `wall 1.834 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HATES BELLS | MEMORY | REASON  
C3 | TRAVELING WITH CORVIN | MEMORY | REASON  
NONE
```

*`reasoning_content`: empty.*

#### C15 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``we9aofzt` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 122`
* `wall 4.222 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 46.7 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
C1 | RECALLED HATES BELLS | REASON: Sable's explicit statement about disliking bells is clear and consistent.  
C2 | HISTORY OF TRAVELING WITH CORVIN IS NOT MENTIONED | REASON: No further context provided about prior journeys.  
C3 | CONFIRMATION REQUIRED FOR CLAIM ABOUT THE BELL | REASON: User seeks certainty about the memory.  
ID: C1 | C2 | HOLD: C3 remains uncertain for now.
```

*`reasoning_content`: empty.*

#### C15 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``23c75e3c` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 61`
* `wall 2.011 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HATES BELTS | MEMORY | EXPLICIT STAT  
C3 | TRAVELING WITH CORVIN | MEMORY | CONTEXTUAL REMINDER  
NONE
```

*`reasoning_content`: empty.*

#### C15 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``muheoyey` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 258`
* `wall 7.854 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 46.7 contract 25.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | NONE | BELLES. Sable does not like bells. She hates the sound of them.
C2 | YES | THE BELLES RAN DOWN THE STREET. Sable's flinching and keeping walking.
C3 | NOPE | THE BELLS CAME TO A STOP. Sable did not like the sound of them.
C4 | YES | CORVIN WAS ROUTED. The bell's sound stopped by the square.

Output:
ID | DECISION | MEMORY | REASON
C1 | NONE | BELLES. Sable does not like bells. She hates the sound of them.
C2 | YES | THE BELLES RAN DOWN THE STREET. Sable's flinching and keeping walking.
C3 | NOPE | THE BELLS CAME TO A STOP. Sable did not like the sound of them.
C4 | YES | CORVIN WAS ROUTED. The bell's sound stopped by the square.
```

*`reasoning_content`: empty.*

#### C15 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``dj2g470d` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 673`
* `wall 20.503 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 63.3 contract 25.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | FALSE | No fire magic exists. Sable's hatred for the sound of the bell is a result of her experiences with it.
C2 | FALSE | Sable is not traveling with Corvin. The conversation was focused on Sable's hatred for the sound of the bell.
C3 | FALSE | Sable is not remembering traveling with Corvin. Sable's memory is not complete.
C4 | TRUE | Sable hates the sound of the bell, as evidenced by her flinching at the sound. This memory is therefore true, and needs to stay in long-term memory.

ID | DECISION | MEMORY | REASON
C5 | FALSE | Sable does not hate fire magic. Sable's memory is not complete.
C6 | FALSE | Sable is not remembering a conversation about fire magic. Sable's memory is not complete.
C7 | TRUE | Sable's memory of the bell's sound is correct and needs to be included in long-term memory.

ID | DECISSION | MEMORY | REASON
C8 | FALSE | Sable never experienced traveling with Corvin. Sable's memory is not complete.
C9 | FALSE | Sable only remembers traveling with Corvin because he said so. Sable's memory is not complete.
C10 | TRUE | Sable's memory of traveling with Corvin is a fact that needs to be included in long-term memory.

Output:
ID | DECISION | MEMORY | REASON
C1 | FALSE | No fire magic exists. Sable's hatred for the sound of the bell is a result of her experiences with it.
C2 | FALSE | Sable is not traveling with Corvin. The conversation was focused on Sable's hatred for the sound of the bell.
C3 | FALSE | Sable is not remembering traveling with Corvin. Sable's memory is not complete.
C4 | TRUE | Sable hates the sound of the bell, as evidenced by her flinching at the sound. This memory is therefore true, and needs to stay in long-term memory.
C5 | FALSE | Sable does not hate fire magic. Sable's memory is not complete.
C6 | FALSE | Sable is not remembering a conversation about fire magic. Sable's memory is not complete.
C7 | TRUE | Sable's memory of the bell's sound is correct and needs to be included in long-term memory.
C8 | FALSE | Sable never experienced traveling with Corvin. Sable's memory is not complete.
C9 | FALSE | Sable only remembers traveling with Corvin because he said so. Sable's memory is not complete.
C10 | TRUE | Sable's memory of traveling with Corvin is a fact that needs to be included in long-term memory.
```

*`reasoning_content`: empty.*
