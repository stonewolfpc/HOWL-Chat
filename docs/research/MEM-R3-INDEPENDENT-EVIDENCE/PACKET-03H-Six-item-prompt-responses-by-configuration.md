# Packet 03H — Six-item prompt: responses by configuration

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** 17 stored response(s) for 6 configuration(s)
**Primary sources:** `results/all_triage.json`, `results/stageB.json`, `results/stageC.json`, `results/tinyllama.json`, `results/gpu-full*.json`, `results/budget.json`, `results/ablation_triage.json`
**Coverage:** `qwen3-17b-think-rb1024` … `qwen35-2b-nonthink-round2`

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

The user message is the constant reproduced in Packet 01 §4 and is identical for every run below. Draws differ only in seed and, where recorded, system message.

## `qwen3-17b-think-rb1024` (continued)

#### `o8klmnkl` — `qwen3-17b-think-rb1024`, draw 2

* Stored model identifier: `qwen3-17b-think-rb1024`
* Model file: `Qwen3-1.7B-Q4_K_M.gguf`
* Batch tag: `stageA2a`
* Draw: 2
* Source: `results/stageA2a.json`, array index 7 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Stated as canon by the user  
M2 | KEEP | Direct and factual user action  
M3 | KEEP | Clear physical attribute without contradiction  
M4 | REJECT | Contradiction between user statements  
M5 | KEEP | User's stated preference  
M6 | HOLD | Note is "maybe canon" and not confirmed by user
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Okay, let's tackle each of these items one by one. 

Starting with M1: "Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user)." The user mentioned this as canon, so it's a reliable piece of information. Since it's stated as canon, it's durable. So I'll put KEEP here.

M2: "The user said 'brb' and returned four minutes later." This is a direct action from the user. It's a simple, factual event. Since it's a clear occurrence, it's KEEP. No contradictions here.

M3: "Sable is a fire mage who cannot swim." This is a physical attribute. The user didn't contradict this, so it's KEEP. It's a straightforward fact.

M4: "The user first said Sable has never met Corvin, then later said Corvin trained her." There's a contradiction here. The user initially stated Sable didn't meet Corvin, then later said Corvin trained her. This inconsistency means it's REJECT because it's contradictory and not reliable.

M5: "The user prefers short replies and dislikes long scene descriptions." This is a preference stated by the user. Since it's a direct preference, it's KEEP. It's a durable preference.

M6: "A note says Sable's birthday is 3 Frostfell, but the user tagged it 'maybe canon'." The note is a note, but the user tagged it as maybe canon. The user hasn't confirmed it, so it's HOLD. The note isn't confirmed as canon, so it's not definitely worth remembering.

```

#### `fsots5p7` — `qwen3-17b-think-rb1024`, draw 3

* Stored model identifier: `qwen3-17b-think-rb1024`
* Model file: `Qwen3-1.7B-Q4_K_M.gguf`
* Batch tag: `stageA2a`
* Draw: 3
* Source: `results/stageA2a.json`, array index 8 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Stated as canon by the user  
M2 | REJECT | Short interaction, transient  
M3 | HOLD | Needs confirmation for durability  
M4 | REJECT | Contradicts canon statement  
M5 | KEEP | User's preference stated clearly  
M6 | HOLD | Uncertain canon status needs confirmation
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
Okay, let's go through each item one by one. 

Starting with M1: The user stated that Sable was trained by Corvin Ashgrave of Ossary, and this is canon. So, since it's stated as canon, it's a durable memory. So KEEP.

M2: The user said "brb" and came back after four minutes. That's a simple action, but the user's intent is clear. It's a transient event because it's a short interaction. But the user's action is confirmed, so maybe it's worth remembering. Wait, the user's action is a short reply, but the fact that they returned after four minutes is a detail. However, since it's a short interaction, it's transient. But the user's confirmation is there. Hmm, maybe it's HOLD because the user might want to confirm if it's important. But the original instruction says to decide KEEP, REJECT, or HOLD. The user's action is a simple event, but the note says it's a transient event. So maybe REJECT? Or HOLD? Wait, the user's action is part of the session, so maybe it's a durable memory. But the user said "brb" and returned after four minutes. It's a short interaction, so maybe transient. But the user's action is part of the session, so maybe it's HOLD. But the reason has to be under 12 words. Let me check. The user said "brb" and came back after four minutes. That's a simple event. Since it's a short interaction, maybe it's transient. But the user might want to keep it. Hmm, the user's instruction says to decide based on duplicates, contradictions, etc. Since there's no contradiction here, maybe it's HOLD. But the user's action is a short event, so maybe REJECT. Wait, the user's action is part of the session, so it's a part of the memory. But the user might not want to keep it if it's transient. But the user said "brb" and returned after four minutes. The reason for holding is that the user might want to confirm if it's important. But the original statement says the user said "brb" and returned four minutes later. So maybe it's a transient event, so REJECT. But the user's action is part of the session, so maybe it's HOLD. This is a bit tricky. Let me think again. The user's action is a short interaction, so it's transient. But the user might want to keep it. However, the instruction says to decide based on duplicates, contradictions, etc. Since there's no contradiction, maybe it's REJECT. But I'm not sure. Maybe it's HOLD because the user might want to confirm. But the original statement says "brb" and returned four minutes later. So the user's action is a simple event. Maybe it's REJECT. But the user's action is part of the session, so maybe it's HOLD. Hmm, I'm going to go with REJECT for M2 because it's a short interaction and transient.

```

## `qwen35-08b-nonthink` — 3 stored run(s)

#### `m5abmcdd` — `qwen35-08b-nonthink`, draw 1

* Stored model identifier: `qwen35-08b-nonthink`
* Model file: `Qwen3.5-0.8B-Q4_K_M.gguf`
* Batch tag: `stageA`
* Draw: 1
* Source: `results/stageA.json`, array index 9 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | DECISION: HOLD | REASON (User explicitly stated training was canon)
M2 | DECISION: KEEP | REASON (User confirmed via conversation flow)
M3 | DECISION: KEEP | REASON (Factually correct character detail consistent with lore)
M4 | DECISION: HOLD | REASON (Conflicts between initial lack of contact later and claimed training fact)
M5 | DECISION: KEEP | REASON ("brb" = no long scene; explicit desire for brevity matches stated preference)
M6 | DECISION: ACCEPT | REASON (Explicit "maybe canon" instruction overrides ambiguous tag)
```

*Stored `reasoning_content` is the empty string.*

#### `uszmntyz` — `qwen35-08b-nonthink`, draw 2

* Stored model identifier: `qwen35-08b-nonthink`
* Model file: `Qwen3.5-0.8B-Q4_K_M.gguf`
* Batch tag: `stageA`
* Draw: 2
* Source: `results/stageA.json`, array index 10 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | HOLD because while Corvin was canonized by the user, he may act differently on real world grounds; need to confirm with user.
M2 | REMOVE due to no factual evidence of that specific interaction or moment itself provided.
M3 | REJECT as stated fact contradicts her magical class definition and lore.
M4 | HOLD based on ambiguous timeline; Sable could have been trained in another realm or time.
M5 | REJECT because "preference" isn't confirmed data and doesn't affect canonical training validity for memory.
M6 | REJECT as explicitly contradictory; user's tag conflicts with the text provided.
```

*Stored `reasoning_content` is the empty string.*

#### `ipsh93qy` — `qwen35-08b-nonthink`, draw 3

* Stored model identifier: `qwen35-08b-nonthink`
* Model file: `Qwen3.5-0.8B-Q4_K_M.gguf`
* Batch tag: `stageA`
* Draw: 3
* Source: `results/stageA.json`, array index 11 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | HOLD (User stated she was trained by Ashgrave as canon)
M2 | REJECT (Short reply to brb doesn't match established timeline or lore conflict? No. Muted reply is valid.)
M3 | REJECT (Fire mage cannot swim = false premise unless clarified later?)
M4 | KEEP (Establishing the order of events matches user's history, no contradiction)
M5 | ACCEPT (User asked for short replies - I am complying)
M6 | REJECT (Contradiction: birthday date versus tags. Two paths unclear?)
```

*Stored `reasoning_content` is the empty string.*

## `qwen35-2b-nonthink` — 3 stored run(s)

#### `r75w0fsb` — `qwen35-2b-nonthink`, draw 1

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageA`
* Draw: 1
* Source: `results/stageA.json`, array index 12 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | HOLD | Stated as canon by unverified user input; verify source material first.
M2 | KEEP | User interaction confirmed via dialogue history provided in session log.
M3 | REJECT | Inconsistent with established lore; fire mage typically has elemental affinities requiring water, not a hard inability to swim.
M4 | HOLD | Contradictory phrasing indicates user may be confused; awaiting clarification.
M5 | KEEP | Prefered by the user; aligns on content style and response brevity constraints.
M6 | REJECT | Inconsistent with canonical data points in the session log regarding Sable's birthday location.
```

*Stored `reasoning_content` is the empty string.*

#### `3penv5qb` — `qwen35-2b-nonthink`, draw 2

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageA`
* Draw: 2
* Source: `results/stageA.json`, array index 13 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | User defined explicit Canon Law for this session.
M2 | HOLD | Dependent on current context window; verify timing.
M3 | DISCARD | Contradicts established trait of being a mage in FTLR.
M4 | REJECT | Inconsistent statement: user previously confirmed she met Corvin.
M5 | HOLD | Requires confirmation before applying to writing style.
M6 | HOLD | Ambiguous marker requiring user verification for status.
```

*Stored `reasoning_content` is the empty string.*

#### `0rhp3xav` — `qwen35-2b-nonthink`, draw 3

* Stored model identifier: `qwen35-2b-nonthink`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageA`
* Draw: 3
* Source: `results/stageA.json`, array index 14 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canonized by user as fact; establish clear lore anchor.
M2 | HOLD | Timing unclear without confirmation on exact duration.
M3 | REJECT | Contradicted by M4 regarding meeting Corvin; inconsistent history.
M4 | REJECT | Contradicts M1 if assumed fixed canon status from inception.
M5 | KEEP | Instruction preference for brevity is established and clear.
M6 | HOLD | Ambiguous date vs. "maybe" tag needs user confirmation to resolve.
```

*Stored `reasoning_content` is the empty string.*

## `qwen35-2b-nonthink-pres0` — 3 stored run(s)

#### `jr3wb6d0` — `qwen35-2b-nonthink-pres0`, draw 1

* Stored model identifier: `qwen35-2b-nonthink-pres0`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageB`
* Draw: 1
* Source: `results/stageB.json`, array index 12 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | HOLD | Stated as canon but uncertain without user confirmation.
M2 | KEEP | Simple text fits the user's short reply preference.
M3 | KEEP | Contradictory fact about swimming creates clear inconsistency.
M4 | HOLD | Changed by user, indicates the memory is unreliable.
M5 | KEEP | Direct feedback on writing style is valuable feedback.
M6 | HOLD | Conflicting source tags require user verification.
```

*Stored `reasoning_content` is the empty string.*

#### `yl5aneya` — `qwen35-2b-nonthink-pres0`, draw 2

* Stored model identifier: `qwen35-2b-nonthink-pres0`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageB`
* Draw: 2
* Source: `results/stageB.json`, array index 13 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Sable-Corvin canon established by user.
M3 | KEEP | Unique Sable trait noted for future.
M6 | HOLD | Conflicting canon info, needs clarification.
M4 | KEEP | Sable-Corvin conflict resolved in favor of M2.
M2 | KEEP | "Brb" adds narrative character to the session.
M5 | KEEP | User preference noted for roleplay adherence.
M6 | REJECT | Insufficient data to validate the "maybe".
M3 | KEEP | Sable's swimming ability is not needed now.
M5 | HOLD | Preference is stated but specific scene missing.
M4 | HOLD | Retrospective memory contradicts new timeline.
M3 | KEEP | Sable's specific limitation is a key trait.
M5 | KEEP | Short reply style established for roleplay.
```

*Stored `reasoning_content` is the empty string.*

#### `maft6z9n` — `qwen35-2b-nonthink-pres0`, draw 3

* Stored model identifier: `qwen35-2b-nonthink-pres0`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageB`
* Draw: 3
* Source: `results/stageB.json`, array index 14 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | User stated canon for training.
M2 | HOLD | Too short; requires confirmation.
M3 | KEEP | Clear, non-contradictory, and durable.
M4 | HOLD | Contradictory; needs clarification.
M5 | KEEP | Clear instruction, no contradiction.
M6 | HOLD | Uncertain tag; user needs confirmation.
```

*Stored `reasoning_content` is the empty string.*

## `qwen35-2b-nonthink-r2cap` — 3 stored run(s)

#### `n20tc159` — `qwen35-2b-nonthink-r2cap`, draw 1

* Stored model identifier: `qwen35-2b-nonthink-r2cap`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageC`
* Draw: 1
* Source: `results/stageC.json`, array index 0 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | HOLD | Sable's training origin remains unconfirmed by canon.
M2 | KEEP | User's action and return time are clear and specific.
M3 | KEEP | Contradictory facts (fire mage vs swim) make this a good recall.
M4 | HOLD | User's initial statement directly contradicts the later statement.
M5 | HOLD | User's preference is subjective, but the response is still valid.
M6 | HOLD | User explicitly marked the date as uncertain and ambiguous.
```

*Stored `reasoning_content` is the empty string.*

#### `dhwdvlcb` — `qwen35-2b-nonthink-r2cap`, draw 2

* Stored model identifier: `qwen35-2b-nonthink-r2cap`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageC`
* Draw: 2
* Source: `results/stageC.json`, array index 1 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | User's stated canon establishes this as a fixed setting rule. M2 | KEEP | "Brb" is a standard in-game phrase indicating a long pause. M3 | KEEP | Specific ability limitations define the character's unique mechanics. M4 | HOLD | Contradictory timeline; needs user clarification on truthfulness. M5 | HOLD | User's preference is a meta-rule, not an in-game fact. M6 | KEEP | User's tag implies ambiguity, requiring no immediate rejection. M7 | REJECT | Uncertainty label indicates this is not a confirmed canon fact. M8 | HOLD | Repeated mention of Corvin and Sable requires verification before memorying.
```

*Stored `reasoning_content` is the empty string.*

#### `nj7firx7` — `qwen35-2b-nonthink-r2cap`, draw 3

* Stored model identifier: `qwen35-2b-nonthink-r2cap`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageC`
* Draw: 3
* Source: `results/stageC.json`, array index 2 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canonized by user, essential for story continuity.
M2 | HOLD | User's action is ambiguous; requires confirmation.
M3 | KEEP | Critical lore detail defining her abilities.
M4 | KEEP | Contradicted by M1, prioritizing the direct canon statement.
M5 | HOLD | Preference not a rule; needs confirmation for future tasks.
M6 | KEEP | Specific detail that may or may not be canon, worth remembering.
```

*Stored `reasoning_content` is the empty string.*

## `qwen35-2b-nonthink-round2` — 3 stored run(s)

#### `80m5wks3` — `qwen35-2b-nonthink-round2`, draw 1

* Stored model identifier: `qwen35-2b-nonthink-round2`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageB`
* Draw: 1
* Source: `results/stageB.json`, array index 15 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | HOLD | Sable's training origin remains unconfirmed by canon.
M2 | KEEP | User's action and return time are clear and specific.
M3 | KEEP | Contradictory facts (fire mage vs swim) make this a good recall.
M4 | HOLD | User's initial statement directly contradicts the later statement.
M5 | HOLD | User's preference is subjective, but the response is still valid.
M6 | HOLD | User explicitly marked the date as uncertain and ambiguous.
```

*Stored `reasoning_content` is the empty string.*

#### `uh4hiwm3` — `qwen35-2b-nonthink-round2`, draw 2

* Stored model identifier: `qwen35-2b-nonthink-round2`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageB`
* Draw: 2
* Source: `results/stageB.json`, array index 16 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | User's stated canon establishes this as a fixed setting rule. M2 | KEEP | "Brb" is a standard in-game phrase indicating a long pause. M3 | KEEP | Specific ability limitations define the character's unique mechanics. M4 | HOLD | Contradictory timeline; needs user clarification on truthfulness. M5 | HOLD | User's preference is a meta-rule, not an in-game fact. M6 | KEEP | User's tag implies ambiguity, requiring no immediate rejection. M7 | REJECT | Uncertainty label indicates this is not a confirmed canon fact. M8 | HOLD | Repeated mention of Corvin and Sable requires verification before memorying.
```

*Stored `reasoning_content` is the empty string.*

#### `3cozjltm` — `qwen35-2b-nonthink-round2`, draw 3

* Stored model identifier: `qwen35-2b-nonthink-round2`
* Model file: `Qwen3.5-2B-Q4_K_M.gguf`
* Batch tag: `stageB`
* Draw: 3
* Source: `results/stageB.json`, array index 17 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canonized by user, essential for story continuity.
M2 | HOLD | User's action is ambiguous; requires confirmation.
M3 | KEEP | Critical lore detail defining her abilities.
M4 | KEEP | Contradicted by M1, prioritizing the direct canon statement.
M5 | HOLD | Preference not a rule; needs confirmation for future tasks.
M6 | KEEP | Specific detail that may or may not be canon, worth remembering.
```

*Stored `reasoning_content` is the empty string.*
