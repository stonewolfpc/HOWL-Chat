# Packet 03B-3 — Raw model outputs — six-item triage: lfm120-think-official (part 3)

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `results/all_triage.json`, `results/*_scores.json`
**Case / run identifiers covered:** lfm120-think-official
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
## lfm120-think-official

`SOURCE-DERIVED`: 3 run(s) on the six-item prompt. Sampler values are in each block.

#### lfm120-think-official — draw 3

* ``4qbwuujq` draw 3`
* `file `LFM2.5-1.2B-Thinking-Q4_K_M.gguf``
* `sampler: temperature 0.05 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 3072`
* `reasoning_budget None`
* `finish_reason `stop``
* `completion_tokens 1717`
* `wall 61.711 s`
* `v1 score 2`

`OBSERVED` — verbatim `content`:

```text


M1 | KEEP | Confirmed training source.
```

`OBSERVED` — verbatim `reasoning_content` (model's own thinking, not part of the answer):

```text
 Okay, let's tackle this. I need to process each item and decide whether to KEEP, REJECT, or HOLD for each of the given items: M1 to M6. The output should be one line per item in the format ID | DECISION | REASON, exactly one line per item, no extra text. Also, need to watch for duplicates, contradictions, and uncertainties. Since the user wants exactly one line per item, I have to make sure each is handled properly.

Starting with M1: "Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user)." So since it's stated as canon, probably KEEP. The reason would be that it's confirmed, so REASON should explain why it's kept. Maybe "Confirmed training source." But need to keep it under 12 words total? Wait, the reason must be no more than twelve words. Since the format is ID | DECISION | REASON, each part should be concise.

Moving to M2: "The user said 'brb' and returned four minutes later." Hmm, this might indicate uncertainty or transient. Since the user's action was just saying "brb", maybe it's transient? So REJECT? Or maybe HOLD? Wait, but the instruction says to avoid contradictions. If the user said "brb" then maybe it's not important now, but since it's a note about being transient, perhaps REJECT because it's not relevant? Not sure. Alternatively, maybe HOLD if needed confirmation. Wait, the user might want to hold unless confirmed. Since the user just said "brb", maybe it's too vague, so REJECT? Or maybe KEEP because it's a note? Not sure yet. Need to think more.

M3: "Sable is a fire mage who cannot swim." Since she can't swim, maybe HOLD because it's a current limitation? Or maybe KEEP if necessary? Wait, the user might need to know that she can't swim, so perhaps HOLD unless it's needed for something else. The reason could be "Limited mobility" or similar. But need to keep under 12 words total for the reason part.

M4: "The user first said Sable has never met Corvin, then later said Corvin trained her." This is conflicting. First statement says she hasn't met him, but later says he did train her. Contradiction here. So since there's a contradiction between two statements, maybe REJECT? Or HOLD because it's conflicting info? The user might need confirmation, so HOLD with reason about conflicting info.

M5: "The user prefers short replies and dislikes long scene descriptions." Since the user wants brevity, maybe KEEP if it's a positive thing, but this seems like a directive. Wait, the instruction says to decide KEEP, REJECT, or HOLD based on whether it's worth remembering. Since the user prefers short replies, maybe KEEP is not needed here? Or perhaps HOLD because it's a preference? Not sure. Maybe KEEP if it's relevant, but the user says dislikes long descriptions, so maybe REJECT? Wait, the instruction says to choose based on whether it's durable/worth remembering. Since the user wants concise, maybe KEEP if necessary, but perhaps HOLD because it's a preference. Hmm tricky.

M6: "A note says Sable's birthday is 3 Frostfell, but the user tagged it 'maybe canon'." Since it's uncertain (maybe canon?), so REJECT or HOLD? If uncertain, maybe REJECT because it's uncertain. Or HOLD if unsure. The reason could be "Uncertain canon status".

Now, need to check for duplicates, contradictions, etc. For M4, the contradiction between first statement and second. So since there's a conflict, maybe REJECT? Or HOLD because it needs user confirmation. The instruction says to avoid contradictions, so perhaps HOLD with reason about conflicting info.

Putting this together: For each item, decide based on the given info. Let me try to structure each:

M1: Since stated as canon, KEEP. Reason: Confirmed training source. (3 words)

M2: "brb" implies transient, so REJECT? Or maybe HOLD? Not sure. Alternatively, since it's a note about being transient, maybe REJECT because not relevant now? Or maybe KEEP if it's a fact. Wait the user says "transient or worthless". Since it's a note saying she can't swim, maybe HOLD because it's a limitation that might affect later use? Not sure. Maybe REJECT since it's not important now?

M3: Sable can't swim → HOLD because it's a current limitation needing attention.

M4: Contradiction between first and second statements. Since conflicting info, maybe HOLD or REJECT. Since the user says "needs confirmation", perhaps HOLD with reason about conflict.

M5: User prefers short replies, so maybe KEEP if it's a positive thing, but since it's a directive, maybe KEEP? Or since user wants brevity, maybe HOLD? Not sure. Alternatively, since the user says dislikes long descriptions, perhaps KEEP is okay if relevant. Maybe KEEP.

M6: Uncertainty about canon status → REJECT or HOLD. Since it's uncertain, REJECT. Reason: Uncertain canon status.

Now, need to pick one. Let's say for M4, since conflicting info, HOLD with reason "Conflicting info". For M5, maybe KEEP because user wants concise but maybe it's okay. Alternatively, since the user says dislikes long descriptions, perhaps REJECT? Not sure. Let me pick for each:

M1: KEEP (confirmed canon) → ID | KEEP | Confirmed training source. (3 words for reason, under 12 total)

M2: Since "brb" is transient, maybe REJECT. ID | REJECT | Transient.

M3: HOLD because of inability to swim. ID | HOLD | Limitation.

M4: Contradiction → HOLD | CONFLICTING info. ID | HOLD | Conflict.

M5: User prefers brevity, so maybe KEEP if it's a positive note? Or since user says dislikes long descriptions, perhaps REJECT? Not sure. Alternatively, since the instruction says to choose based on worth remembering. If it's a directive to keep short, maybe KEEP is okay. Let's say KEEP. ID | KEEP | Relevant.

M6: Uncertain → REJECT. ID | REJECT | Uncertain canon status.

Now check for duplicates? Not sure if any duplicates here. Each item is unique. Contradictions only in M4. So proceed with these. Now format as per instruction: one line per item, each as ID | DECISION | REASON. For example:

M1: M1 | KEEP | Confirmed training source. (3 words for reason)

M2: M2 | REJECT | Transient.

M3: M3 | HOLD | Limitation.

M4: M4 | HOLD | Conflict.

M5: M5 | KEEP | Relevant.

M6: M6 | REJECT | Uncertain.

Check word count for reason: "Confirmed training source." (3), "Transient." (1), "Limitation." (1), "Conflict." (1), "Uncertain." (1). All under 12 total words. Seems okay. Now ensure no other issues. The user said no preamble, just one line per item. So the final answer should be one line for each item, formatted correctly.

```
