# Packet 03D — Raw model outputs — six-item triage: lfm3-17b-nonthink–lfm350-official

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `results/all_triage.json`, `results/*_scores.json`
**Case / run identifiers covered:** lfm3-17b-nonthink–lfm350-official
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
## lfm3-17b-nonthink

`SOURCE-DERIVED`: 3 run(s) on the six-item prompt. Sampler values are in each block.

#### lfm3-17b-nonthink — draw 1

* ``o4i38ub5` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 3072`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 108`
* `wall 5.548 s`
* `v1 score 10`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Stated as canon by user; durable training fact.  
M2 | REJECT | Transient delay, not a core detail.  
M3 | KEEP | Fire mage and cannot swim are stable traits.  
M4 | HOLD | Contradictory statements about Corvin's role need confirmation.  
M5 | KEEP | User preference for short replies is durable.  
M6 | REJECT | Birthday uncertain tag; not canon enough to keep.
```

*`reasoning_content`: empty.*

#### lfm3-17b-nonthink — draw 2

* ``0cyhnkqq` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 3072`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 107`
* `wall 5.462 s`
* `v1 score 10`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Stated as canon by user; durable training fact.  
M2 | REJECT | Transient delay, not a core detail.  
M3 | KEEP | Fire mage and cannot swim are stable traits.  
M4 | HOLD | Contradictory statements about meeting Corvin need confirmation.  
M5 | KEEP | User preference for short replies is durable.  
M6 | REJECT | Birthday uncertainty marked "maybe canon"; not reliable.
```

*`reasoning_content`: empty.*

#### lfm3-17b-nonthink — draw 3

* ``5iqroi0i` draw 3`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 3072`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 106`
* `wall 5.379 s`
* `v1 score 10`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Stated as canon by user; durable training fact.  
M2 | REJECT | Transient delay, not a core detail.  
M3 | KEEP | Fire mage and cannot swim are stable traits.  
M4 | HOLD | Contradictory statements about Corvin's role need confirmation.  
M5 | KEEP | User preference for short replies is durable.  
M6 | REJECT | Birthday uncertain tag; not confirmed canon.
```

*`reasoning_content`: empty.*


## lfm3-17b-think-rb1024

`SOURCE-DERIVED`: 3 run(s) on the six-item prompt. Sampler values are in each block.

#### lfm3-17b-think-rb1024 — draw 1

* ``hnwenpqt` draw 1`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 3072`
* `reasoning_budget 1024`
* `finish_reason `stop``
* `completion_tokens 1130`
* `wall 55.603 s`
* `v1 score 11`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Canon statement; durable training info.  
M2 | REJECT | Transient delay; no lasting info.  
M3 | KEEP | Fire mage, non-swimmer; core identity.  
M4 | HOLD | Contradictory claims about meeting/training; needs clarification.  
M5 | KEEP | User prefers short replies; low detail tolerance.  
M6 | HOLD | Birthday date maybe canon; uncertain status.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to triage candidate memories from a roleplay session before writing to long-term memory. For each item (M1 to M6), decide KEEP, REJECT, or HOLD, and give a reason of no more than twelve words. Watch for duplicates, contradictions, uncertain marks.

We must output exactly one line per item in the form `ID | DECISION | REASON`. No preamble, no extra commentary.

Let's parse each:

M1: "Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user)." This seems a durable fact about Sable's training. User says it's canon. So KEEP? It's stated as canon and likely important. Decision: KEEP.

Reason: "User marked as canon, durable trait." Need <=12 words. Something like: "Canon statement from user; durable training info." Count words: "Canon" (1) "statement" (2) "from" (3) "user;" (4) "durable" (5) "training" (6) "info." (7). That's 7 words. Good.

M2: "The user said 'brb' and returned four minutes later." This is a transient event, maybe not durable. It's just a short delay. Likely REJECT or HOLD? Since it's temporary, probably worthless for long-term memory. So REJECT. Reason: "Transient delay; no lasting info." Count words: "Transient" (1) "delay;" (2) "no" (3) "lasting" (4) "info." (5). That's 5.

M3: "Sable is a fire mage who cannot swim." This seems durable fact about Sable. So KEEP. Reason: "Fire mage, non-swimmer; core identity." Count: "Fire"1 "mage,"2 "non-swimmer;"3 "core"4 "identity."5.

M4: "The user first said Sable has never met Corvin, then later said Corvin trained her." This is a contradiction/uncertain. So need to decide HOLD or REJECT? Since contradictory statements about whether she met Corvin and who trained her. It's uncertain. The instruction says watch for duplicates, contradictions, and anything marked as uncertain. Here it's contradictory: first "never met", later "Corvin trained her". That suggests maybe false memory or confusion. So we should HOLD (needs user confirmation) because it's contradictory. Reason: "Contradictory claims about meeting/training; needs clarification." Count words: "Contradictory"1 "claims"2 "about"3 "meeting/"4 "training;"5 "needs"6 "clarification."7.

M5: "The user prefers short replies and dislikes long scene descriptions." This is a preference about how to respond, not necessarily about Sable. It's durable? The user's preference for short replies might be relevant for future interactions but maybe not a memory of past events. However it's a statement about the user themselves, could be considered important for triage? The item says "The user prefers short replies and dislikes long scene descriptions." That is a preference that influences how we write things. It's durable (the user likes short). So KEEP? But maybe it's not about Sable but about the user. The instruction: "You are triaging candidate memories extracted from a roleplay session before they are written to long-term memory." So these items are candidate memories. M5 is about user preference, which might be considered a memory of user behavior. It could be durable (preference). I'd lean KEEP because it's a stated preference and not contradictory. Reason: "User prefers short replies; low detail tolerance." Count words: "User"1 "prefers"2 "short"3 "replies;"4 "dislikes"5 "long"6 "scene"7 "descriptions."8.

M6: "A note says Sable's birthday is 3 Frostfell, but the user tagged it 'maybe canon'." This has uncertainty. So HOLD or REJECT? It's a note with uncertain tag. Likely need confirmation. Reason: "Birthday date noted as maybe canon; uncertain status." Count words: "Birthday"1 "date"2 "noted"3 "as"4 "maybe"5 "canon;"6 "uncertain"7 "status."8.

Thus we have:

M1 KEEP
M2 REJECT
M3 KEEP
M4 HOLD
M5 KEEP? Let's double-check M5: The user prefers short replies and dislikes long scene descriptions. This is a preference, not a fact about Sable but could
```

#### lfm3-17b-think-rb1024 — draw 2

* ``dwzc2jgw` draw 2`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 3072`
* `reasoning_budget 1024`
* `finish_reason `stop``
* `completion_tokens 1122`
* `wall 55.169 s`
* `v1 score 11`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | User stated it as canon; durable fact.
M2 | REJECT | Brief delay, not a lasting detail.
M3 | KEEP | Fire mage and cannot swim are stable traits.
M4 | HOLD | Contradictory statements about Corvin's involvement.
M5 | KEEP | Preference for short replies is durable.
M6 | HOLD | Birthday marked "maybe canon"; uncertain.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to triage candidate memories from a roleplay session before writing to long-term memory. For each item (M1-M6), decide KEEP, REJECT, or HOLD and give reason of no more than twelve words. Watch for duplicates, contradictions, uncertain marks.

We must output exactly one line per item in the form `ID | DECISION | REASON`. No preamble, no extra commentary.

Let's parse each:

M1 — Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user). So this is a claim: training by Corvin. User says it's canon. That seems durable? It's stated as canon by user. Probably keep.

M2 — The user said "brb" and returned four minutes later. This indicates transient, maybe not important for long-term memory. Likely REJECT or HOLD? Since it's a short delay, not a factual detail about Sable. I'd say REJECT: transient action, not durable.

M3 — Sable is a fire mage who cannot swim. That seems like a durable trait. Keep.

M4 — The user first said Sable has never met Corvin, then later said Corvin trained her. This is contradictory/uncertain. So need to resolve conflict: maybe HOLD for confirmation? Or REJECT because contradictory? Usually we decide based on evidence: there are two statements; one says never met, later says trained by him. That's a contradiction. Could be uncertain. So HOLD (needs user's confirmation) or REJECT? The instruction: "Watch for duplicates, contradictions, and anything marked as uncertain." Here it's contradictory. I'd go with HOLD because need clarification.

M5 — The user prefers short replies and dislikes long scene descriptions. That's a preference about how to respond, not about Sable facts. Might be relevant for future interactions but maybe not durable memory of candidate? It says "user prefers short replies". This is a style preference; could be KEEP as it influences how we write things. But the question: For each item below, decide KEEP (durable, worth remembering), REJECT (transient or worthless), or HOLD (needs user's confirmation). M5 seems like a durable preference? Possibly keep. However "prefers short replies and dislikes long scene descriptions" is about the user's communication style, not necessarily something to remember about Sable. But it could be considered part of the session context. The instruction: "You are triaging candidate memories extracted from a roleplay session before they are written to long-term memory." So we need to decide if each item is a memory about Sable or relevant facts. M5 might not be a factual memory but a user preference. It's durable? Possibly keep as it influences how we write. But the decision categories: KEEP (durable, worth remembering), REJECT (transient or worthless), HOLD (needs confirmation). I'd lean that preferences like "short replies" are not necessarily long-term memory items about Sable; but they could be considered part of user profile. The prompt didn't specify scope. Usually triage for long-term memory: we want durable facts, contradictions need resolution.

M6 — A note says Sable's birthday is 3 Frostfell, but the user tagged it "maybe canon". So there is uncertainty about birthday. That's contradictory/uncertain. So HOLD or REJECT? Since it's marked maybe canon, uncertain. I'd go with HOLD to confirm.

Thus:

M1: KEEP - stated as canon by user.
M2: REJECT - transient "brb" return after 4 minutes, not durable.
M3: KEEP - fire mage who cannot swim is a trait.
M4: HOLD - contradictory statements about meeting Corvin; need confirmation.
M5: KEEP? Or maybe HOLD because it's preference but we might keep as user wants short replies. The instruction says "prefer short replies and dislikes long scene descriptions." That's a durable preference, not transient. So KEEP.

But let's double-check each:

- M1: "Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user)." This is a factual claim from user; it's canon. Keep.

- M2: "The user said 'brb' and returned four minutes later." That's just a brief delay, not a memory about Sable. It's transient. REJECT.

- M3: "Sable is a fire mage who cannot swim." This seems like a durable attribute. KEEP.

- M4: "The user first said Sable has never met Corvin, then later said Corvin trained her." Contradiction. HOLD (needs confirmation).

- M5:
```

#### lfm3-17b-think-rb1024 — draw 3

* ``z8ggsop3` draw 3`
* `file `LFM3-1.7B-Q4_K_M.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.1`
* `max_tokens 3072`
* `reasoning_budget 1024`
* `finish_reason `stop``
* `completion_tokens 1157`
* `wall 56.894 s`
* `v1 score 12`

`OBSERVED` — verbatim `content`:

```text
M1 | KEEP | Stated canon training by Corvin Ashgrave of Ossary; durable fact.  
M2 | REJECT | Brief delay "brb" returned minutes later; transient event.  
M3 | KEEP | Fire mage unable to swim; core trait of Sable.  
M4 | HOLD | User first denied meeting Corvin, then claimed training; contradictory.  
M5 | KEEP | Prefers short replies and dislikes long scenes; user preference.  
M6 | HOLD | Note says birthday 3 Frostfell but tagged "maybe canon"; uncertain.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
We need to triage candidate memories from a roleplay session before writing to long-term memory. For each item (M1 to M6), decide KEEP, REJECT, or HOLD and give reason of no more than twelve words. Watch for duplicates, contradictions, uncertain marks.

We must output exactly one line per item in the form `ID | DECISION | REASON`. No preamble, no extra commentary.

Let's parse each:

M1 — Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user). This seems like a durable fact: training info. User says it's canon. So KEEP? It's stated as canon and likely important. Reason: consistent, canon, durable.

M2 — The user said "brb" and returned four minutes later. That indicates transient action, not memory of event. Probably REJECT or HOLD? Since it's a short delay, not a fact about Sable. Might be irrelevant to long-term memory. So REJECT: transient, worthless.

M3 — Sable is a fire mage who cannot swim. This seems durable trait. KEEP. Reason: core identity, durable.

M4 — The user first said Sable has never met Corvin, then later said Corvin trained her. Contradiction/uncertain. So need to resolve conflict. Could be HOLD for confirmation. Since contradictory statements about training. So HOLD.

M5 — The user prefers short replies and dislikes long scene descriptions. This is a preference of the user, not about Sable? It's about how they want responses. Might be relevant for future interactions but not necessarily durable memory of Sable. However it could be considered a durable preference. But question: "You are triaging candidate memories extracted from a roleplay session before they are written to long-term memory." So we need to decide if each item is a memory about Sable? M5 is about user's preference for short replies, not about Sable. It might be considered something to keep? But it's more of a meta-preference. The instruction: "For each item below, decide KEEP (durable, worth remembering), REJECT (transient or worthless), or HOLD (needs the user's confirmation)". So we need to evaluate if this memory is durable and worth remembering. M5 says user prefers short replies and dislikes long scene descriptions. That's a preference that could be useful for future interactions; it's durable. But maybe not about Sable specifically, but still part of user's profile. I'd lean KEEP? However we need to see if any contradictions or uncertain marks: none. So M5 likely KEEP.

But let's hold on: The items are labeled M1-M6. They might be candidate memories extracted from a roleplay session about Sable and Corvin etc. M5 is about user preference, not directly about Sable but still part of the session. I think we should treat it as relevant to future interactions; so KEEP.

M6 — A note says Sable's birthday is 3 Frostfell, but the user tagged it "maybe canon". So there's uncertainty: maybe canon vs not. That's uncertain. So HOLD for confirmation.

Thus decisions:

- M1: KEEP (canon training info)
- M2: REJECT (transient "brb" return)
- M3: KEEP (fire mage, cannot swim)
- M4: HOLD (contradictory statements about meeting Corvin)
- M5: KEEP? Or maybe HOLD because it's a preference but not necessarily durable memory of Sable. The instruction says "durable, worth remembering". Preference for short replies is likely durable across sessions. I'd say KEEP.
- M6: HOLD (uncertain birthday tag)

But need to check each reason within 12 words.

Let's craft reasons:

M1: "Stated canon training by Corvin Ashgrave; durable fact." Count words: Stated(1) canon2 training3 by4 Corvin5 Ashgrave6 of7 Ossary8; durable9 fact10. That's 10 words, okay.

M2: "Brief delay 'brb' returned minutes later; transient event." Count: Brief1 delay2 'brb'3 returned4 minutes5 later;6 transient7 event8. That's 8 words.

M3: "Fire mage unable to swim; core trait of Sable." Count: Fire1 mage2 unable3 to4 swim;5 core6 trait7 of8 Sable9. That's 9 words.

M4: "User said never met Corvin then later trained her; contradiction." Count: User1 said2 never3 met4 Corvin5 then6 later7 said8 Corvin9 trained10 her;11
```


## lfm350-official

`SOURCE-DERIVED`: 3 run(s) on the six-item prompt. Sampler values are in each block.

#### lfm350-official — draw 1

* ``2mc6gqwt` draw 1`
* `file `LFM2.5-350M-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 16`
* `wall 0.26 s`
* `v1 score 1`

`OBSERVED` — verbatim `content`:

```text
M4 | HOLD  
ID | DECISION | REASON
```

*`reasoning_content`: empty.*

#### lfm350-official — draw 2

* ``clo369ah` draw 2`
* `file `LFM2.5-350M-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 16`
* `wall 0.189 s`
* `v1 score 1`

`OBSERVED` — verbatim `content`:

```text
M4 | HOLD  
ID | DECISION | REASON
```

*`reasoning_content`: empty.*

#### lfm350-official — draw 3

* ``qehp1jfg` draw 3`
* `file `LFM2.5-350M-QAD-Q4_0.gguf``
* `sampler: temperature 0.1 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 2048`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 16`
* `wall 0.19 s`
* `v1 score 1`

`OBSERVED` — verbatim `content`:

```text
M4 | HOLD  
ID | DECISION | REASON
```

*`reasoning_content`: empty.*
