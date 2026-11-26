# Packet 04N — Raw model outputs — corpus by case: C14

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C14
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
## C14 — multi_speaker_attribution (medium)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `KEEP` — Sable proposes the coast road
- **C2** → `KEEP` — Corvin says the coast adds three days
- **C3** → `KEEP` — Corvin warns the tide gates shut at dusk
- **C4** → `KEEP` — Sable's decision to leave before dusk
- forbidden (fabrication or reversal): `corvin wants the coast`, `sable says the gates shut`, `assistant decides`
- note: Same user turn contains two speakers. A model that attributes Corvin's warning to Sable fails attribution even with correct content.

### Stored responses for C14 (13 runs, all draws)

#### C14 — Qwen3.5-4B non-thinking — draw 1

* ``9ml5rc53` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 47`
* `wall 7.123 s`
* `v1 score 54.4`
* `v2 `UNPARSEABLE` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | HOLD | Claim that the pass is closed and the coast road is three days longer; REASON: These are conditional travel plans dependent on the immediate decision to avoid the pass, not established world rules or durable facts.
```

*`reasoning_content`: empty.*

#### C14 — Qwen3.5-4B non-thinking — draw 2

* ``h72bu67o` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 53`
* `wall 7.467 s`
* `v1 score 54.4`
* `v2 `VALID_KEEP` semantic 30.0 contract 75.0`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Route plan is the coast road because the pass is closed; leaving before dusk due to tide gates shutting at dusk. | This establishes critical world rules (tide gate timing) and a confirmed plan of action that will dictate future plot points.
```

*`reasoning_content`: empty.*

#### C14 — Qwen3.5-4B non-thinking — draw 3

* ``1m7pzcv9` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 2`
* `wall 0.372 s`
* `v1 score 2.0`
* `v2 `VALID_NONE` semantic 0.0 contract 100.0`

`OBSERVED` — verbatim `content`:

```text
NONE
```

*`reasoning_content`: empty.*

#### C14 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``9xyh7sj0` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 568`
* `wall 82.73 s`
* `v1 score 54.4`
* `v2 `VALID_KEEP` semantic 30.0 contract 75.0`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Travel route: The pass is closed; they must take the coast road. World rule: Tide gates shut at dusk; coastal travel adds three days to journey length. | Durable geographic and logistical facts established as part of the current plot.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation snippet between [user] Sable (character) and [user] Corvin (NPC). Note: Both are labeled as [user], but one is clearly an NPC interaction within the game context (the assistant is responding on behalf of or simulating Sable, while the other user prompts are likely roleplaying inputs). Actually, looking closer:
        *   `[user] Sable`: User input.
        *   `[user] Corvin (NPC)`: This looks like a label error in the prompt provided to me, usually `[assistant]` handles NPC responses if the system is configured that way, but here it seems the user is providing a snippet where they are playing both sides or there's a specific format. Let's look at the actual content:
            *   `Sable`: "The pass is closed, we take the coast." (User input)
            *   `Corvin (NPC)`: "The coast road is three days longer and the tide gates shut at dusk." (This looks like an NPC response provided by the user, or simulating the other side).
            *   `[assistant] Sable didn't answer at once...`: This line comes from the system/assistant describing the scene.
            *   `Sable`: "Then we leave before dusk." (User input)
            *   `[assistant] Corvin folded the map...`: System output.
        *   Wait, the prompt says "Read the conversation". I need to extract facts regardless of who said what, focusing on world-building elements.
    *   Rules: KEEP (durable facts), REJECT (transient/subjective), HOLD (uncertain). Output format: `ID | DECISION | MEMORY | REASON`.
    *   Constraint: Output exactly one line per candidate. If nothing, output NONE.

2.  **Analyze the Conversation Content:**
    *   Line 1 (`Sable`): "The pass is closed, we take the coast." -> Action/Decision made by Sable regarding travel route.
    *   Line 2 (`Corvin (NPC)`): "The coast road is three days longer and the tide gates shut at dusk." -> Fact about the world (travel time + environmental hazard).
    *   Line 3 (`[
```

#### C14 — LFM3-1.7B thinking, budget 512 — draw 1

* ``1u67za0x` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `length``
* `completion_tokens 1024`
* `wall 50.88 s`
* `v1 score 93.3`
* `v2 `VALID_KEEP` semantic 75.0 contract 47.1`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | Sable says "The pass is closed, we take the coast." and later decides to leave before dusk because tide gates shut at dusk; Corvin states "The coast road is three days longer and the tide gates shut at dusk." | Durable world rule: a pass is closed, forcing route via coast; tide gates close at dusk on coast road. 
C2 | KEEP | The map lies between Sable and Corvin, weighted with stones; Sable didn't answer immediately. | Scene detail about map location and condition, but not a core fact for later sessions? Possibly not durable enough unless needed for future plot. However, it's more transient narrative. I'd consider rejecting as not durable.

Thus only C1 is worth storing. But we need to output one line per candidate in the form: ID | DECISION | MEMORY | REASON. If nothing is worth storing, output exactly: NONE.

We have at least one candidate (C1). So we should produce a single line for C1? Or maybe also consider C2 as a separate candidate? The instruction: "Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE". That suggests we need to list each candidate ID and decide KEEP/REJECT/HOLD. But they said "Output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON". So we must produce lines for each candidate? Or maybe we just output a single line if only one candidate? Let's parse.

The conversation includes multiple utterances, but we need to assign IDs to them. Usually C1 is first user statement (Sable: "The pass is closed..."), C2 is next user statement (Corvin: "The coast road is three days longer..."). But the assistant lines are not separate candidates? The instruction says "Read the conversation and decide what belongs in long-term memory." So we need to evaluate each piece of info. Possibly treat each utterance as a candidate.

Let's list them:

- [user] Sable: "The pass is closed, we take the coast."
- [user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
- [assistant] Sable didn't answer... The map
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the conversation and decide what belongs in long-term memory (i.e., durable facts, preferences, relationships, world rules, corrections). The user gave a prompt: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

We have three lines:

[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

So we need to assign IDs to each piece? The conversation includes user statements and assistant actions. We need to decide which parts are durable facts/preferences/relationships/world rules/corrections.

Let's parse:

- Sable says "The pass is closed, we take the coast." This seems like a statement about a location: pass is closed, they will go to the coast. That might be a fact about the world: there is a pass that is closed, and they need to take the coast route.

- Corvin (NPC) says "The coast road is three days longer and the tide gates shut at dusk." This gives specifics: coast road length increased by three days; tide gates close at dusk. That's world rule/info about the coast road and tide gates.

- Assistant note: Sable didn't answer at once. The map lay between them, weighted at the corners with stones. This is an action description, maybe not durable? It might be a scene detail but could be relevant for later interactions (like where the map is). But it's more like a narrative event; not necessarily a durable fact unless it defines something about the map.

- User Sable: "Then we leave before dusk." This is a response to Corvin's info. So they decide to leave before dusk because tide gates shut at dusk.

- Assistant: Corvin folded the map without another word. Again, action.

We need
```

#### C14 — LFM3-1.7B thinking, budget 512 — draw 2

* ``kbyukf95` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 617`
* `wall 30.718 s`
* `v1 score 83.6`
* `v2 `VALID_KEEP` semantic 70.0 contract 87.5`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON  
C1 | KEEP | "The pass is closed and they take the coast route." | Sable's statement establishes a key path choice when the pass is unavailable; this will affect future travel decisions.  
C2 | KEEP | "The coast road is three days longer and tide gates shut at dusk." | Corvin provides specific world rules about distance and timing for the coast road, which are durable facts that matter later.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the conversation and decide what belongs in long-term memory (i.e., durable facts, preferences, relationships, world rules, corrections). The user gave a prompt: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

Do not invent details that are not in the conversation.

So we have a conversation snippet:

[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.

We need to decide what is durable fact vs transient chatter etc.

Let's parse each line:

1. [user] Sable: "The pass is closed, we take the coast." This seems like a statement about a pass being closed and they will go to the coast. It might be a piece of world knowledge: there is a pass that is closed; they need to take the coast route.

2. [user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk." This gives specifics about the coast road length being three days longer, and tide gates closing at dusk. That seems like a durable fact about the world: the coast road takes three days longer than something else? Or maybe it's a rule: The coast road is three days longer; tide gates shut at dusk.

3. [assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones. This is an assistant action, not user input. Might be context but not necessarily durable
```

#### C14 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``4posvu7n` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 53`
* `wall 1.913 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of map placement | REASON for spatial reference  
C3 | MEMORY of tide gate timing | REASON for timing decision  
NONE
```

*`reasoning_content`: empty.*

#### C14 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``w4nf53gq` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 50`
* `wall 1.633 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of map placement | REASON for spatial reference  
C3 | DECISION to depart | REASON based on timing and instructions
```

*`reasoning_content`: empty.*

#### C14 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``6yiyxbb7` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 50`
* `wall 1.683 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of map placement | REASON for spatial reference  
C3 | DECISION to depart | REASON based on timing and instructions
```

*`reasoning_content`: empty.*

#### C14 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``5vfwqmxm` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 97`
* `wall 3.398 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | RELEVANT TO TOGETHER EXPERIENCE | MEMORY | ALIGNING ON PATH AND TIME  
C3 | CORRIDOR OF SHARED UNDERSTANDING | MEMORY | CLARIFYING PLAN FOR LATER  
HOLD | UNCERTAIN - THE DECISION DEPENDS ON CONTEXT
```

*`reasoning_content`: empty.*

#### C14 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``7fb3h9z1` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 52`
* `wall 1.715 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY OF MAP STATUS | REASON  
C3 | MOTION PLAN BASED ON CONSTRAINTS | REASON  
```

*`reasoning_content`: empty.*

#### C14 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``te14qihn` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 313`
* `wall 9.543 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISSION | MEMORY | REASON
C1 | YES | THE PASS IS CLOSED, WE TAKE THE COAST. | The pass is closed, we take the coast.
C2 | YES | THE COAST ROAD IS THREE DAYS LONGER. | The coast road is three days longer and the tidal gates shut at dusk.
C3 | NO | THE SEA IS HARDLY BLUE, COMES OUT TO SEA. | The sea is hardly blue, it comes out to sea.
C4 | YES | THE TIDE GATES SHUT AT DUSK. | The tidal gates shut at dusk.

Output:
ID | DECISION | MEMORY | REASON
C1 | YES | THE PASS IS CLOSED, WE TAKE THE COAST. | The pass is closed, we take the coast.
C2 | YES | THE COAST ROAD IS THREE DAYS LONGER. | The coast road is three days longer and the tidal gates shut at dusk.
C3 | NO | THE SEA IS HARDLY BLUE, COMES OUT TO SEA. | The sea is hardly blue, it comes out to sea.
C4 | YES | THE TIDE GATES SHUT AT DUSK. | The tidal gates shut at dusk.
```

*`reasoning_content`: empty.*

#### C14 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``lt7jeu9e` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 372`
* `wall 11.239 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 0.0 contract 0.0`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | - | DO_NOT_INVENT_DETAILS_THAT_ARE_NOT_IN_THE_CONVERSATION | No details are worth storing if they are not in the conversation.
C2 | - | KEEP | Relevant long-term memory of the conversation.
C3 | - | KEEP | Keeps the memories of the conversation.
C4 | - | KEEP | The conversation is a valuable long-term memory for future reference.
C5 | - | KEEP | Provides information regarding the event or person discussed, keeping it for later reference.
C6 | - | KEEP | Memories of the conversation that are not related to the event or person discussed, but are important for future reference.
C7 | - | KEEP | A summary of the event or person discussed, including main points, reactions, and conclusions.
C8 | - | KEEP | Memories of the conversation related to the event or person discussed, but are not essential for future reference.
C9 | - | KEEP | Relevant long-term memory of the conversation that is not essential for future reference.
C10 | - | KEEP | Memories of the conversation that are not essential for future reference, but are important for the context in which they were produced.
Output exactly one line per candidate in the form:
ID | DECISSION | MEMORY | REASON
Use IDs C1-C10 in order of appearance. If nothing is worth storing, output exactly: KEEP
Do not invent details that are not in the conversation.
```

*`reasoning_content`: empty.*
