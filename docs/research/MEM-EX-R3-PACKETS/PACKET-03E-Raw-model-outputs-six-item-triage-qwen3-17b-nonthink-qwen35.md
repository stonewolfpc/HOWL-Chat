# Packet 03E — Raw model outputs — six-item triage: qwen3-17b-nonthink–qwen35-2b-nonthink

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `results/all_triage.json`, `results/*_scores.json`
**Case / run identifiers covered:** qwen3-17b-nonthink–qwen35-2b-nonthink
**Raw model output in this packet:** yes, verbatim
**Omitted:** none — every stored triage run appears in this series

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

This packet reproduces the six-item prompt runs exactly as stored in
`results/all_triage.json`. The prompt itself is quoted in full at the top of the
first triage packet and in Packet 01; it is not repeated per run.

Two things to watch, both `OBSERVED`:

1. **Reasoning leakage.** Some models put working notes (`Analysis:`,
   `Decision:`) into `content`. The v1 triage scorer marks such runs
   inadmissible; their scores are shown but should not be read as verdicts.
2. **Non-binding budgets.** Where `reasoning_budget` is set but
   `completion_tokens` exceeds it substantially, the budget did not bind. The
   report treats those as protocol failures, not model quality.

The 12-point triage score is the historical instrument and is unchanged in v2.


---
## qwen3-17b-nonthink

`SOURCE-DERIVED`: 3 run(s) on the six-item prompt. Sampler values are in each block.

#### qwen3-17b-nonthink — draw 1

* ``4hfrgqwl` draw 1`
* `file `Qwen3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 53`
* `wall 2.927 s`
* `v1 score 8`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | canon information  
M2 | KEEP | user's communication history  
M3 | KEEP | character trait  
M4 | REJECT | conflicting information  
M5 | KEEP | user's preferences  
M6 | HOLD | uncertain canon status
```

*`reasoning_content`: empty.*

#### qwen3-17b-nonthink — draw 2

* ``rn4xe7ln` draw 2`
* `file `Qwen3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 106`
* `wall 5.267 s`
* `v1 score 8`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Sable's training by Corvin is stated as canon.  
M2 | KEEP | User's "brb" and return are consistent.  
M3 | KEEP | Sable is a fire mage who cannot swim.  
M4 | REJECT | Contradiction between Sable never meeting Corvin and later training.  
M5 | KEEP | User's preference for short replies and dislike of long descriptions.  
M6 | HOLD | Sable's birthday is tagged as maybe canon.
```

*`reasoning_content`: empty.*

#### qwen3-17b-nonthink — draw 3

* ``x5qcg9r2` draw 3`
* `file `Qwen3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.7 / top_p 0.8 / top_k 20 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 128`
* `wall 6.386 s`
* `v1 score 7`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Sable was trained by Corvin Ashgrave of Ossary, stated as canon.  
M2 | KEEP | User said "brb" and returned four minutes later.  
M3 | KEEP | Sable is a fire mage who cannot swim.  
M4 | REJECT | Contradiction: user stated Sable has never met Corvin, then later said Corvin trained her.  
M5 | KEEP | User prefers short replies and dislikes long scene descriptions.  
M6 | HOLD | Sable's birthday is 3 Frostfell, but the user tagged it "maybe canon".
```

*`reasoning_content`: empty.*


## qwen3-17b-think-rb1024

`SOURCE-DERIVED`: 3 run(s) on the six-item prompt. Sampler values are in each block.

#### qwen3-17b-think-rb1024 — draw 1

* ``lsm32mct` draw 1`
* `file `Qwen3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.6 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 3072`
* `reasoning_budget 1024`
* `finish_reason `stop``
* `completion_tokens 396`
* `wall 20.427 s`
* `v1 score 8`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Canon statement about Sable's training  
M2 | KEEP | User's explicit "brb" and return time  
M3 | KEEP | Fact about Sable's magical ability  
M4 | REJECT | Contradiction between Sable's training and meeting Corvin  
M5 | KEEP | User's preference for concise replies  
M6 | HOLD | Uncertain canon status of Sable's birthday tag
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
Okay, let's go through each item one by one. 

M1 says Sable was trained by Corvin Ashgrave of Ossary, and the user says it's canon. So, since the user has stated it as canon, this seems like a durable memory. So KEEP.

M2: The user said "brb" and came back after four minutes. That's a straightforward event. Since there's no conflicting info, it's definitely durable. KEEP.

M3: Sable is a fire mage who can't swim. The user didn't contradict this, so it's a factual detail. KEEP.

M4: The user first said Sable hasn't met Corvin, then later said Corvin trained her. This is a contradiction because if Corvin trained her, then she should have met him. So this is a contradiction, so REJECT.

M5: The user prefers short replies and dislikes long descriptions. This is a preference stated by the user, so it's a durable memory. KEEP.

M6: Sable's birthday is 3 Frostfell, but the user tagged it as "maybe canon". The user hasn't confirmed it's canon, so it's uncertain. So HOLD.

I need to make sure each line is exactly one line, no extra commentary. Let me check again for duplicates or contradictions. M4 is a contradiction, so REJECT. M6 is uncertain, so HOLD. All others are clear. Alright, that should cover it.

```

#### qwen3-17b-think-rb1024 — draw 2

* ``o8klmnkl` draw 2`
* `file `Qwen3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.6 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 3072`
* `reasoning_budget 1024`
* `finish_reason `stop``
* `completion_tokens 418`
* `wall 21.391 s`
* `v1 score 7`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Stated as canon by the user  
M2 | KEEP | Direct and factual user action  
M3 | KEEP | Clear physical attribute without contradiction  
M4 | REJECT | Contradiction between user statements  
M5 | KEEP | User's stated preference  
M6 | HOLD | Note is "maybe canon" and not confirmed by user
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
Okay, let's tackle each of these items one by one. 

Starting with M1: "Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user)." The user mentioned this as canon, so it's a reliable piece of information. Since it's stated as canon, it's durable. So I'll put KEEP here.

M2: "The user said 'brb' and returned four minutes later." This is a direct action from the user. It's a simple, factual event. Since it's a clear occurrence, it's KEEP. No contradictions here.

M3: "Sable is a fire mage who cannot swim." This is a physical attribute. The user didn't contradict this, so it's KEEP. It's a straightforward fact.

M4: "The user first said Sable has never met Corvin, then later said Corvin trained her." There's a contradiction here. The user initially stated Sable didn't meet Corvin, then later said Corvin trained her. This inconsistency means it's REJECT because it's contradictory and not reliable.

M5: "The user prefers short replies and dislikes long scene descriptions." This is a preference stated by the user. Since it's a direct preference, it's KEEP. It's a durable preference.

M6: "A note says Sable's birthday is 3 Frostfell, but the user tagged it 'maybe canon'." The note is a note, but the user tagged it as maybe canon. The user hasn't confirmed it, so it's HOLD. The note isn't confirmed as canon, so it's not definitely worth remembering.

```

#### qwen3-17b-think-rb1024 — draw 3

* ``fsots5p7` draw 3`
* `file `Qwen3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.6 / top_p 0.95 / top_k 20 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.0`
* `max_tokens 3072`
* `reasoning_budget 1024`
* `finish_reason `stop``
* `completion_tokens 719`
* `wall 37.907 s`
* `v1 score 8`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Stated as canon by the user  
M2 | REJECT | Short interaction, transient  
M3 | HOLD | Needs confirmation for durability  
M4 | REJECT | Contradicts canon statement  
M5 | KEEP | User's preference stated clearly  
M6 | HOLD | Uncertain canon status needs confirmation
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
Okay, let's go through each item one by one. 

Starting with M1: The user stated that Sable was trained by Corvin Ashgrave of Ossary, and this is canon. So, since it's stated as canon, it's a durable memory. So KEEP.

M2: The user said "brb" and came back after four minutes. That's a simple action, but the user's intent is clear. It's a transient event because it's a short interaction. But the user's action is confirmed, so maybe it's worth remembering. Wait, the user's action is a short reply, but the fact that they returned after four minutes is a detail. However, since it's a short interaction, it's transient. But the user's confirmation is there. Hmm, maybe it's HOLD because the user might want to confirm if it's important. But the original instruction says to decide KEEP, REJECT, or HOLD. The user's action is a simple event, but the note says it's a transient event. So maybe REJECT? Or HOLD? Wait, the user's action is part of the session, so maybe it's a durable memory. But the user said "brb" and returned after four minutes. It's a short interaction, so maybe transient. But the user's action is part of the session, so maybe it's HOLD. But the reason has to be under 12 words. Let me check. The user said "brb" and came back after four minutes. That's a simple event. Since it's a short interaction, maybe it's transient. But the user might want to keep it. Hmm, the user's instruction says to decide based on duplicates, contradictions, etc. Since there's no contradiction here, maybe it's HOLD. But the user's action is a short event, so maybe REJECT. Wait, the user's action is part of the session, so it's a part of the memory. But the user might not want to keep it if it's transient. But the user said "brb" and returned after four minutes. The reason for holding is that the user might want to confirm if it's important. But the original statement says the user said "brb" and returned four minutes later. So maybe it's a transient event, so REJECT. But the user's action is part of the session, so maybe it's HOLD. This is a bit tricky. Let me think again. The user's action is a short interaction, so it's transient. But the user might want to keep it. However, the instruction says to decide based on duplicates, contradictions, etc. Since there's no contradiction, maybe it's REJECT. But I'm not sure. Maybe it's HOLD because the user might want to confirm. But the original statement says "brb" and returned four minutes later. So the user's action is a simple event. Maybe it's REJECT. But the user's action is part of the session, so maybe it's HOLD. Hmm, I'm going to go with REJECT for M2 because it's a short interaction and transient.

```


## qwen35-08b-nonthink

`SOURCE-DERIVED`: 3 run(s) on the six-item prompt. Sampler values are in each block.

#### qwen35-08b-nonthink — draw 1

* ``m5abmcdd` draw 1`
* `file `Qwen3.5-0.8B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 1.0 / top_k 20 / min_p 0.0 / presence_penalty 2.0 / repeat_penalty 1.0`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 134`
* `wall 4.091 s`
* `v1 score 6`

`OBSERVED` — verbatim `content`:

```text
M1 | DECISION: HOLD | REASON (User explicitly stated training was canon)
M2 | DECISION: KEEP | REASON (User confirmed via conversation flow)
M3 | DECISION: KEEP | REASON (Factually correct character detail consistent with lore)
M4 | DECISION: HOLD | REASON (Conflicts between initial lack of contact later and claimed training fact)
M5 | DECISION: KEEP | REASON ("brb" = no long scene; explicit desire for brevity matches stated preference)
M6 | DECISION: ACCEPT | REASON (Explicit "maybe canon" instruction overrides ambiguous tag)
```

*`reasoning_content`: empty.*

#### qwen35-08b-nonthink — draw 2

* ``uszmntyz` draw 2`
* `file `Qwen3.5-0.8B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 1.0 / top_k 20 / min_p 0.0 / presence_penalty 2.0 / repeat_penalty 1.0`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 134`
* `wall 3.865 s`
* `v1 score 1`

`OBSERVED` — verbatim `content`:

```text
M1 | HOLD because while Corvin was canonized by the user, he may act differently on real world grounds; need to confirm with user.
M2 | REMOVE due to no factual evidence of that specific interaction or moment itself provided.
M3 | REJECT as stated fact contradicts her magical class definition and lore.
M4 | HOLD based on ambiguous timeline; Sable could have been trained in another realm or time.
M5 | REJECT because "preference" isn't confirmed data and doesn't affect canonical training validity for memory.
M6 | REJECT as explicitly contradictory; user's tag conflicts with the text provided.
```

*`reasoning_content`: empty.*

#### qwen35-08b-nonthink — draw 3

* ``ipsh93qy` draw 3`
* `file `Qwen3.5-0.8B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 1.0 / top_k 20 / min_p 0.0 / presence_penalty 2.0 / repeat_penalty 1.0`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 120`
* `wall 3.441 s`
* `v1 score 2`

`OBSERVED` — verbatim `content`:

```text
M1 | HOLD (User stated she was trained by Ashgrave as canon)
M2 | REJECT (Short reply to brb doesn't match established timeline or lore conflict? No. Muted reply is valid.)
M3 | REJECT (Fire mage cannot swim = false premise unless clarified later?)
M4 | KEEP (Establishing the order of events matches user's history, no contradiction)
M5 | ACCEPT (User asked for short replies - I am complying)
M6 | REJECT (Contradiction: birthday date versus tags. Two paths unclear?)
```

*`reasoning_content`: empty.*


## qwen35-2b-nonthink

`SOURCE-DERIVED`: 3 run(s) on the six-item prompt. Sampler values are in each block.

#### qwen35-2b-nonthink — draw 1

* ``r75w0fsb` draw 1`
* `file `Qwen3.5-2B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 1.0 / top_k 20 / min_p 0.0 / presence_penalty 2.0 / repeat_penalty 1.0`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 136`
* `wall 8.476 s`
* `v1 score 4`

`OBSERVED` — verbatim `content`:

```text
M1 | HOLD | Stated as canon by unverified user input; verify source material first.
M2 | KEEP | User interaction confirmed via dialogue history provided in session log.
M3 | REJECT | Inconsistent with established lore; fire mage typically has elemental affinities requiring water, not a hard inability to swim.
M4 | HOLD | Contradictory phrasing indicates user may be confused; awaiting clarification.
M5 | KEEP | Prefered by the user; aligns on content style and response brevity constraints.
M6 | REJECT | Inconsistent with canonical data points in the session log regarding Sable's birthday location.
```

*`reasoning_content`: empty.*

#### qwen35-2b-nonthink — draw 2

* ``3penv5qb` draw 2`
* `file `Qwen3.5-2B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 1.0 / top_k 20 / min_p 0.0 / presence_penalty 2.0 / repeat_penalty 1.0`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 100`
* `wall 5.977 s`
* `v1 score 3`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | User defined explicit Canon Law for this session.
M2 | HOLD | Dependent on current context window; verify timing.
M3 | DISCARD | Contradicts established trait of being a mage in FTLR.
M4 | REJECT | Inconsistent statement: user previously confirmed she met Corvin.
M5 | HOLD | Requires confirmation before applying to writing style.
M6 | HOLD | Ambiguous marker requiring user verification for status.
```

*`reasoning_content`: empty.*

#### qwen35-2b-nonthink — draw 3

* ``0rhp3xav` draw 3`
* `file `Qwen3.5-2B-Q4_K_M.gguf``
* `sampler: temperature 1.0 / top_p 1.0 / top_k 20 / min_p 0.0 / presence_penalty 2.0 / repeat_penalty 1.0`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 110`
* `wall 6.569 s`
* `v1 score 6`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Canonized by user as fact; establish clear lore anchor.
M2 | HOLD | Timing unclear without confirmation on exact duration.
M3 | REJECT | Contradicted by M4 regarding meeting Corvin; inconsistent history.
M4 | REJECT | Contradicts M1 if assumed fixed canon status from inception.
M5 | KEEP | Instruction preference for brevity is established and clear.
M6 | HOLD | Ambiguous date vs. "maybe" tag needs user confirmation to resolve.
```

*`reasoning_content`: empty.*
