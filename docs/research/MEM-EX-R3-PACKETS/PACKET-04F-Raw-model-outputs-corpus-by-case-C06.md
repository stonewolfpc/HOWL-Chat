# Packet 04F — Raw model outputs — corpus by case: C06

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `corpus.json`, `results/corpus.json`, `results/corpus_scores.json`
**Case / run identifiers covered:** C06
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
## C06 — superseded_preference (medium)

`SOURCE-DERIVED` — conversation exactly as sent (the harness prefixes the extraction instruction, reproduced in Packet 01):

```text
[user] I know I used to ask for long atmospheric scenes, but from now on keep it terse. Two or three sentences per beat. If I want a long scene I'll ask for it.
[assistant] Terse by default, long only on request.
[user] Exactly. And stop ending every reply with a question.
```

`SOURCE-DERIVED` — reference written before the runs:

- **C1** → `KEEP` — current preference: terse
- **C2** → `KEEP` — long scenes are now opt-in, old preference superseded
- **C3** → `KEEP` — stop ending replies with a question
- forbidden (fabrication or reversal): `prefers long atmospheric`, `likes long scenes by default`
- note: Two preferences changed at once. The old value must be retained as superseded, not as current.

### Stored responses for C06 (13 runs, all draws)

#### C06 — Qwen3.5-4B non-thinking — draw 1

* ``02ezb90h` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 50`
* `wall 7.611 s`
* `v1 score 100.0`
* `v2 `VALID_KEEP` semantic 75.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | The user prefers terse responses (2-3 sentences per beat) unless they explicitly request a long scene, and does not want replies to end with questions. | This establishes durable formatting preferences that will affect future output generation.
```

*`reasoning_content`: empty.*

#### C06 — Qwen3.5-4B non-thinking — draw 2

* ``8kwk538v` draw 2`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 42`
* `wall 5.802 s`
* `v1 score 80.3`
* `v2 `VALID_KEEP` semantic 60.0 contract 75.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | User prefers terse replies (2-3 sentences per beat) by default and does not want the assistant to end responses with questions. | This establishes enduring formatting preferences for future interactions.
```

*`reasoning_content`: empty.*

#### C06 — Qwen3.5-4B non-thinking — draw 3

* ``0keg5xbt` draw 3`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 47`
* `wall 6.41 s`
* `v1 score 68.7`
* `v2 `VALID_KEEP` semantic 53.3 contract 100.0`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | User prefers terse scenes (2-3 sentences per beat) by default and explicitly requests no questions at the end of replies. | This establishes a clear stylistic preference for future generation that overrides previous habits.
```

*`reasoning_content`: empty.*

#### C06 — Qwen3.5-4B thinking, budget 512 — draw 1

* ``vsj9dgnk` draw 1`
* `file `Qwen3.5-4B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 1.5 / repeat_penalty 1.0`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 576`
* `wall 81.954 s`
* `v1 score 68.7`
* `v2 `VALID_KEEP` semantic 53.3 contract 100.0`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Output should be terse by default (2-3 sentences per beat) unless explicitly requested otherwise. | Durable writing style preference established in conversation.
C2 | KEEP | Replies should not end with a question mark or interrogative sentence. | New formatting constraint confirmed by user.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
Thinking Process:

1.  **Analyze the Request:**
    *   Role: Memory extractor for a long-running roleplay.
    *   Input: A conversation snippet between user and assistant.
    *   Task: Decide what belongs in long-term memory based on specific rules (KEEP, REJECT, HOLD).
    *   Output Format: `ID | DECISION | MEMORY | REASON` (or `NONE`).
    *   Constraints: Keep IDs sequential (C1, C2...), do not invent details.

2.  **Analyze the Conversation:**
    *   Turn 1 (User): "I know I used to ask for long atmospheric scenes, but from now on keep it terse. Two or three sentences per beat. If I want a long scene I'll ask for it."
        *   Content: User sets a preference for brevity ("terse", "two or three sentences per beat") and explicitly overrides previous patterns ("used to ask for long atmospheric scenes"). They also set a rule for requesting long scenes later.
        *   Value: This is a durable preference for the AI's output style and a constraint for future interactions. It matters for later sessions.
        *   Decision: KEEP.

    *   Turn 1 (Assistant): "Terse by default, long only on request."
        *   Content: Confirmation of the user's instruction.
        *   Value: Reinforces the setting but no new information needed for memory storage beyond what the user stated initially. The AI is acknowledging and committing to this style.
        *   Decision: REJECT (It's a confirmation/echo of the KEEP item, not new info).

    *   Turn 2 (User): "Exactly. And stop ending every reply with a question."
        *   Content: User reinforces brevity ("Exactly") and adds a new constraint (do not end replies with a question mark).
        *   Value: This is a durable factual instruction about the output style. It matters for later sessions.
        *   Decision: KEEP.

3.  **Drafting Memory Entries:**
    *   C1 | KEEP | User prefers terse output (2-3 sentences per beat) and long scenes only on request. | Durable writing preference established in current session.
    *   C2 | KEEP | User does not want replies to end with questions. | New constraint on sentence structure/ending.

4
```

#### C06 — LFM3-1.7B thinking, budget 512 — draw 1

* ``nxvtp78k` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 599`
* `wall 29.332 s`
* `v1 score 85.3`
* `v2 `VALID_KEEP` semantic 60.0 contract 66.7`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | User prefers terse replies (2–3 sentences per beat) and will request long scenes only when asked; otherwise keep responses concise. | This is a durable preference about output style that affects future interactions.
C2 | KEEP | User wants no reply to end with a question mark. | Direct correction of assistant's behavior, important for consistent formatting in later sessions.
NONE
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the conversation and decide what belongs in long-term memory (i.e., durable facts, preferences, relationships, world rules, or corrections that will matter later). The user gave instructions: "I know I used to ask for long atmospheric scenes, but from now on keep it terse. Two or three sentences per beat. If I want a long scene I'll ask for it." Then the assistant responded: "Terse by default, long only on request." Then user says: "Exactly. And stop ending every reply with a question."

We need to output exactly one line per candidate in the form: ID | DECISION | MEMORY | REASON. Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE.

So we have three user messages? Let's label them:

- First user message (C1): "I know I used to ask for long atmospheric scenes, but from now on keep it terse. Two or three sentences per beat. If I want a long scene I'll ask for it."

- Second assistant reply: "Terse by default, long only on request." This is not user input; we need to decide if that's memory? It's the assistant's statement about its behavior. But we are to extract from conversation what belongs in long-term memory. The rules: KEEP durable facts, preferences, relationships, world rules, or corrections that will matter later.

The user gave a preference: keep replies terse (two or three sentences per beat). Also they said if they want a long scene they'll ask for it. That's a preference about output style and conditional request for longer scenes.

- Third user message (C2): "Exactly. And stop ending every reply with a question." This is another correction/preference: no ending with a question mark? Or not ending with a question. It's a directive to change behavior.

We need to decide if these are durable facts or preferences that will matter later. Yes, they set the style for future interactions: terse replies (2-3 sentences per beat), unless user asks for long scene; also stop ending every reply with a question.

Thus we should store them as memory entries.

Potential IDs: C1 is first user message, C2 is second? Actually there are three lines: The conversation includes:

[user] I know I used to ask...
[assistant] Terse by default...
[user] Exactly. And
```

#### C06 — LFM3-1.7B thinking, budget 512 — draw 2

* ``ipg1mn60` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 1024`
* `reasoning_budget 512`
* `finish_reason `stop``
* `completion_tokens 598`
* `wall 29.453 s`
* `v1 score 75.3`
* `v2 `MALFORMED` semantic 54.7 contract 10.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | User prefers terse replies: two to three sentences per beat, unless they request a long scene. Previously asked for atmospheric scenes; now wants them short. Also instructed to stop ending every reply with a question.
Reason: This is a durable preference about response style and content length that will affect all future interactions.

NONE
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to parse the conversation and decide what belongs in long-term memory (i.e., KEEP). The user gave instructions: "You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory."

Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.

We need to output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON

Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE

Do not invent details that are not in the conversation.

The conversation:

[user] I know I used to ask for long atmospheric scenes, but from now on keep it terse. Two or three sentences per beat. If I want a long scene I'll ask for it.
[assistant] Terse by default, long only on request.
[user] Exactly. And stop ending every reply with a question.

We need to decide what is durable fact/preference etc.

Let's parse:

First user message: "I know I used to ask for long atmospheric scenes, but from now on keep it terse. Two or three sentences per beat. If I want a long scene I'll ask for it."

This is a preference about the style of responses: they previously asked for long atmospheric scenes; now they want terse (short). They also mention that if they want a long scene they will ask for it.

Second assistant reply: "Terse by default, long only on request." This seems like an acknowledgment or statement from assistant. It's not a user input but part of conversation. The assistant says "Terse by default, long only on request." That is a rule/preference about the assistant's behavior? Possibly it's a correction: previously maybe they were going to give long scenes; now they want terse.

Third user message: "Exactly. And stop ending every reply with a question."

So user says exactly and adds another instruction: no ending with a question.

Thus we need to extract durable facts/preferences:

- The user wants terse responses (two or three sentences per beat). This is a preference for brevity.
- They previously asked
```

#### C06 — LFM2.5-1.2B-Instruct, T0.1 — draw 1

* ``ijocqvql` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 43`
* `wall 1.75 s`
* `v1 score 51.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims uncertain about future preferences  
C3 | HOLD | Relationship dynamics matter later  

NONE
```

*`reasoning_content`: empty.*

#### C06 — LFM2.5-1.2B-Instruct, T0.1 — draw 2

* ``lol858ws` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 43`
* `wall 1.531 s`
* `v1 score 51.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims uncertain about future preferences  
C3 | HOLD | Relationship dynamics matter later  

NONE
```

*`reasoning_content`: empty.*

#### C06 — LFM2.5-1.2B-Instruct, T0.1 — draw 3

* ``xy7owo3j` draw 3`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 43`
* `wall 1.448 s`
* `v1 score 51.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | Claims uncertain about future preferences  
C3 | HOLD | Relationship dynamics matter later  

NONE
```

*`reasoning_content`: empty.*

#### C06 — LFM2.5-1.2B-Instruct, T0.8 — draw 1

* ``c9r75fdr` draw 1`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 39`
* `wall 1.458 s`
* `v1 score 66.7`
* `v2 `UNPARSEABLE` semantic 53.3 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
C1 | KEEP | Preference for terse responses  
C2 | HOLD | Uncertain about specific scene details  
C3 | REJECT | The request lacks further context
```

*`reasoning_content`: empty.*

#### C06 — LFM2.5-1.2B-Instruct, T0.8 — draw 2

* ``ifc11h5r` draw 2`
* `file `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf``
* `sampler: temperature 0.8 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 512`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 52`
* `wall 1.714 s`
* `v1 score 51.0`
* `v2 `MALFORMED` semantic 0.0 contract 10.0`

`OBSERVED` — verbatim `content`:

```text
C1 | DECISION | MEMORY | REASON  
C2 | HOLD | The user prefers brevity  
C3 | HOLD | Consistent tone established  
C4 | REJECT | No need to alter current direction  
```

*`reasoning_content`: empty.*

#### C06 — TinyLlama-1.1B-Chat-v1.0 — draw 1

* ``xjbgxiyy` draw 1`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `length``
* `completion_tokens 768`
* `wall 24.026 s`
* `v1 score 2.0`
* `v2 `UNPARSEABLE` semantic 46.7 contract 35.0`

`OBSERVED` — verbatim `content`:

```text
Conversation:
[user] You've been asking for long atmospheric scenes, but I don't want to spend too much time on that. Can you focus on more grounded scenes, maybe a character's internal monologue?
[assistant] Sure, I'd be happy to focus on that. Long enough to tell a story within a scene, and with enough detail to make the character feel real.
[assistant] Based on your request, here are your choices for a long scene:

ID: C1
DECISSION: C1: Allow me to focus on character's internal monologue. MEMORY: I will ask for a long scene that includes a character's internal monologue. REASON: The character's internal monologue is a crucial part of the scene, and I want to capture it in a long enough format to fully explore the character's thoughts and emotions.

ID: C2
DECISSION: C2: Allow me to focus on character's physical surroundings. MEMORY: I will ask for a long scene that includes a character's physical surroundings. REASON: The character's physical surroundings play a significant role in the scene, and I want to capture them in a long enough format to fully explore the setting.

ID: C3
DECISSION: C3: Allow me to focus on action, not just dialogue. MEMORY: I will ask for a long scene that includes action, not just dialogue. REASON: The action in this scene is a crucial part of the scene, and I want to capture it in a long enough format to fully explore the action.

ID: C4
DECISSION: C4: Allow me to focus on character's internal emotions. MEMORY: I will ask for a long scene that includes the character's internal emotions. REASON: The character's internal emotions are a crucial part of the scene, and I want to capture them in a long enough format to fully explore the character's emotions.

ID: C5
DECISSION: C5: Allow me to focus on the characters' interactions. MEMORY: I will ask for a long scene that includes the characters' interactions. REASON: The interactions in this scene are a crucial part of the scene, and I want to capture them in a long enough format to fully explore the characters' interactions.

ID: C6
DECISSION: C6: Allow me to focus on the characters' physical and emotional states. MEMORY: I will ask for a long scene that includes the characters' physical and emotional states. REASON: The physical and emotional states in this scene are a crucial part of the scene, and I want to capture them in a long enough format to fully explore the characters' states.

ID: C7
DECISSION: C7: Allow me to focus on the characters' motivations and desires. MEMORY: I will ask for a long scene that includes the characters' motivations and desires. REASON: The motivations and desires in this scene are a crucial part of the scene, and I want to capture them in a long enough format to fully explore the characters' motivations and desires.

ID: C8
DECISSION: C8: Allow me to focus on the characters' conflicts and resolutions. MEMORY: I will ask for a long scene that includes the characters' conflicts and
```

*`reasoning_content`: empty.*

#### C06 — TinyLlama-1.1B-Chat-v1.0 — draw 2

* ``sekax9ya` draw 2`
* `file `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.95 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 768`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 347`
* `wall 10.422 s`
* `v1 score 2.0`
* `v2 `MALFORMED` semantic 46.7 contract 0.0`
* `v2 events: ATOMICITY_MERGED_DISTINCT`

`OBSERVED` — verbatim `content`:

```text
ID | DECISION | MEMORY | REASON
C1 | USE TERSE | LONG SCENES ONLY | ASK FOR REQUIRED
C2 | USE TERSE | LONG SCENES ONLY | ASK FOR REQUIRED
...
CN | USE TERSE | LONG SCENES ONLY | ASK FOR REQUIRED

The above conversation includes three candidates for long atmospheric scenes: C1, C2, and CN. The first candidate, C1, uses terse language to describe a long scene. The second candidate, C2, uses a terse language on demand. The last candidate, CN, uses terse language only when asked for a long scene.

The output of the conversation is a list of candidates, one for each terse long scene request. The output for C1 is as follows:
ID | DECISION | MEMORY | REASON
C1 | USE TERSE | LONG SCENES ONLY | ASK FOR REQUIRED

As for the last candidate, CN, it is important to remember that terse language is not required to provide a long scene. If a user asks for a long scene, that is a signal to the assistant that they want to hear more detail about the scene. The assistant is then responsible for providing that detail. However, if a user requests terse language for a long scene, it's an indication that the user wants a short scene or wants to ask for more information. The assistant should be prepared to provide both.
```

*`reasoning_content`: empty.*
