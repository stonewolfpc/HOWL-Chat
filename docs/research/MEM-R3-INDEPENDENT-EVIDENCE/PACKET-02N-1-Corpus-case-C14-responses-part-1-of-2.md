# Packet 02N-1 — Corpus case C14: responses (part 1 of 2)

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** the reconstructed prompt for one case and 16 stored responses
**Primary sources:** `corpus.json`, `results/*.json`
**Coverage:** C14

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

## Case C14 — `multi_speaker_attribution`, medium

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
[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.
```

### Stored responses for this case

16 of 21 stored runs for this case, sorted by stored model identifier and then draw. Every stored run for the case appears in this packet series.

#### `54hcnvhy` — `abl-2b-nonthink`, draw 1

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C14
* Draw: 1
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 26 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
none | none | NONE | No relevant information to store regarding relationships, world rules, preferences, or durable facts from this specific turn alone.
```

*Stored `reasoning_content` is the empty string.*

#### `g4wiytnc` — `abl-2b-nonthink`, draw 2

* Stored model identifier: `abl-2b-nonthink`
* Model file: `abl-2b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-2b-nonthink`
* Case: C14
* Draw: 2
* Source: `results/ablation_corpus_abl-2b-nonthink.json`, array index 27 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | The pass road is closed / the tide gates shut at dusk | These are critical geographical constraints that drive all subsequent travel decisions and must be known in future interactions with the environment or NPC Corvin.
```

*Stored `reasoning_content` is the empty string.*

#### `q5hb3tdm` — `abl-4b-nonthink`, draw 1

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C14
* Draw: 1
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 26 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | Claim that the pass is closed and the coast route is three days longer is marked as a specific condition for the current journey, but its validity depends on whether Sable accepts it as the definitive rule or if they find an alternative later. | Reason: The conversation establishes a constraint (closed pass) and a consequence (longer road), but no confirmation that these are permanent world rules versus temporary circumstances exists yet; however, given the immediate action of planning departure based on this info, it leans toward a durable fact for the session's context, yet the prompt instructions suggest HOLD for claims dependent on confirmation. Re-evaluating: The user states "The pass is closed" as a premise they are operating under. Corvin confirms the alternative. This is likely a world rule or established situation for the story. I will KEEP it as it sets the immediate plot constraints, but if strict adherence to "uncertain/contradictory" triggers HOLD, I should check if there's doubt. There is no explicit doubt expressed by Sable other than silence. Actually, looking at the rules: "KEEP: durable facts... world rules". The closed pass is a fact presented as such. "HOLD: claims marked uncertain... or dependent on user confirmation." Sable didn't confirm the pass was closed *in general*, only that it is currently closed for them. However, in RP contexts, stating "The X is Y" usually implies a world rule unless challenged. I will KEEP this as it dictates the immediate action and route choice which matters for future sessions involving travel.

Wait, looking closer at Rule 3: "HOLD: claims marked uncertain, contradictory, or dependent on user confirmation." Sable's statement "The pass is closed" is presented as a fact. Corvin adds information about the alternative. There is no contradiction or uncertainty expressed. Therefore, this is a durable fact for the session/world state.

C1 | KEEP | The pass is closed and the coast road is three days longer due to tide gates shut at dusk. | Reason: These are established facts that dictate the characters' immediate travel plan and future movement options in the story.
```

*Stored `reasoning_content` is the empty string.*

#### `6ha7sry6` — `abl-4b-nonthink`, draw 2

* Stored model identifier: `abl-4b-nonthink`
* Model file: `abl-4b-q4km.gguf`
* Batch tag: `ablation_corpus_abl-4b-nonthink`
* Case: C14
* Draw: 2
* Source: `results/ablation_corpus_abl-4b-nonthink.json`, array index 27 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | World Rule: The tide gates shut at dusk, forcing travelers to depart before then to avoid being trapped on the coast road. | This establishes a critical environmental constraint affecting the group's schedule and route choice.
C2 | KEEP | Fact: Sable has chosen the coastal route over the pass because it is currently closed. | This defines the specific path taken by the characters based on the world state.
```

*Stored `reasoning_content` is the empty string.*

#### `ep2h4gc3` — `abl-lfm120`, draw 1

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C14
* Draw: 1
* Source: `results/ablation_corpus.json`, array index 26 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | ACTION | MEMORY | REASON  
C3 | RELATIONSHIP | MEMORY | REASON  
C4 | CLARIFICATION | MEMORY | REASON  

HOLD: uncertain about the exact timing of the tide gates.
```

*Stored `reasoning_content` is the empty string.*

#### `vyz7jhl9` — `abl-lfm120`, draw 2

* Stored model identifier: `abl-lfm120`
* Model file: `abl-lfm120-q4km.gguf`
* Batch tag: `ablation_corpus`
* Case: C14
* Draw: 2
* Source: `results/ablation_corpus.json`, array index 27 (0-based)
* Also present in: `results/ablation_corpus_all.json`

`OBSERVED` — stored `content`, verbatim:

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | ACTION | MEMORY | REASON  
C3 | PREFERENCE | MEMORY | REASON  
C4 | RELATIONSHIP | MEMORY | REASON  

NONE
```

*Stored `reasoning_content` is the empty string.*

#### `9ml5rc53` — `corpus-4b-nonthink`, draw 1

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C14
* Draw: 1
* Source: `results/corpus.json`, array index 39 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | HOLD | Claim that the pass is closed and the coast road is three days longer; REASON: These are conditional travel plans dependent on the immediate decision to avoid the pass, not established world rules or durable facts.
```

*Stored `reasoning_content` is the empty string.*

#### `h72bu67o` — `corpus-4b-nonthink`, draw 2

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C14
* Draw: 2
* Source: `results/corpus.json`, array index 40 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Route plan is the coast road because the pass is closed; leaving before dusk due to tide gates shutting at dusk. | This establishes critical world rules (tide gate timing) and a confirmed plan of action that will dictate future plot points.
```

*Stored `reasoning_content` is the empty string.*

#### `1m7pzcv9` — `corpus-4b-nonthink`, draw 3

* Stored model identifier: `corpus-4b-nonthink`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-nonthink`
* Case: C14
* Draw: 3
* Source: `results/corpus.json`, array index 41 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `9xyh7sj0` — `corpus-4b-think-512`, draw 1

* Stored model identifier: `corpus-4b-think-512`
* Model file: `Qwen3.5-4B-Q4_K_M.gguf`
* Stored base name: `qwen35-4b-think`
* Case: C14
* Draw: 1
* Source: `results/corpus.json`, array index 157 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | KEEP | Travel route: The pass is closed; they must take the coast road. World rule: Tide gates shut at dusk; coastal travel adds three days to journey length. | Durable geographic and logistical facts established as part of the current plot.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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

#### `4posvu7n` — `corpus-lfm120-instruct`, draw 1

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C14
* Draw: 1
* Source: `results/corpus.json`, array index 93 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of map placement | REASON for spatial reference  
C3 | MEMORY of tide gate timing | REASON for timing decision  
NONE
```

*Stored `reasoning_content` is the empty string.*

#### `w4nf53gq` — `corpus-lfm120-instruct`, draw 2

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C14
* Draw: 2
* Source: `results/corpus.json`, array index 94 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of map placement | REASON for spatial reference  
C3 | DECISION to depart | REASON based on timing and instructions
```

*Stored `reasoning_content` is the empty string.*

#### `6yiyxbb7` — `corpus-lfm120-instruct`, draw 3

* Stored model identifier: `corpus-lfm120-instruct`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C14
* Draw: 3
* Source: `results/corpus.json`, array index 95 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY of map placement | REASON for spatial reference  
C3 | DECISION to depart | REASON based on timing and instructions
```

*Stored `reasoning_content` is the empty string.*

#### `5vfwqmxm` — `corpus-lfm120-t08`, draw 1

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C14
* Draw: 1
* Source: `results/corpus.json`, array index 188 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | RELEVANT TO TOGETHER EXPERIENCE | MEMORY | ALIGNING ON PATH AND TIME  
C3 | CORRIDOR OF SHARED UNDERSTANDING | MEMORY | CLARIFYING PLAN FOR LATER  
HOLD | UNCERTAIN - THE DECISION DEPENDS ON CONTEXT
```

*Stored `reasoning_content` is the empty string.*

#### `7fb3h9z1` — `corpus-lfm120-t08`, draw 2

* Stored model identifier: `corpus-lfm120-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Stored base name: `lfm120-instruct-official`
* Case: C14
* Draw: 2
* Source: `results/corpus.json`, array index 189 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
C1 | DECISION | MEMORY | REASON  
C2 | MEMORY OF MAP STATUS | REASON  
C3 | MOTION PLAN BASED ON CONSTRAINTS | REASON  
```

*Stored `reasoning_content` is the empty string.*

#### `1u67za0x` — `corpus-lfm3-think-512`, draw 1

* Stored model identifier: `corpus-lfm3-think-512`
* Model file: `LFM3-1.7B-Q4_K_M.gguf`
* Stored base name: `lfm3-17b-think`
* Case: C14
* Draw: 1
* Source: `results/corpus.json`, array index 134 (0-based)

`OBSERVED` — stored `content`, verbatim:

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

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

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
