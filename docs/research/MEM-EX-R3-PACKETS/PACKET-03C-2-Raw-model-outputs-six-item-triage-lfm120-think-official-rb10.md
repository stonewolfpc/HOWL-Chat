# Packet 03C-2 — Raw model outputs — six-item triage: lfm120-think-official-rb1024 (part 2)

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `results/all_triage.json`, `results/*_scores.json`
**Case / run identifiers covered:** lfm120-think-official-rb1024
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
## lfm120-think-official-rb1024

`SOURCE-DERIVED`: 3 run(s) on the six-item prompt. Sampler values are in each block.

#### lfm120-think-official-rb1024 — draw 3

* ``nbpyv0nw` draw 3`
* `file `LFM2.5-1.2B-Thinking-Q4_K_M.gguf``
* `sampler: temperature 0.05 / top_p 1.0 / top_k 50 / min_p 0.0 / presence_penalty 0.0 / repeat_penalty 1.05`
* `max_tokens 3072`
* `reasoning_budget 1024`
* `finish_reason `stop``
* `completion_tokens 2274`
* `wall 83.432 s`
* `v1 score 8`

`OBSERVED` — verbatim `content`:

```text


M1 | KEEP | Confirmed training source.  
M2 | REJECT | Transient note.  
M3 | HOLD | Fire mage limitation.  
M4 | HOLD | Conflicting info.  
M5 | KEEP | Short preference.  
M6 | REJECT | Uncertain canon status.  

Wait, but need to ensure each line is ID | DECISION | REASON, and exactly one line per item. Also check for duplicates? Not sure if any duplicates here. Contradictions: M4 has conflict between first and second statement. So M4's HOLD with reason about contradiction. The user said "watch for contradictions", so that's okay. All others seem okay. Each entry must be one line. The answer should be exactly one line per item, but the user says respond with exactly one line per item? Wait wait, the user says "Respond with exactly one line per item in the form ID | DECISION | REASON...", but the instruction says "no preamble and no extra commentary. Items: ..." So perhaps each item is a separate line? Wait the user said "respond with exactly one line per item", but the user says "no preamble and no extra commentary". Since the user says "watch for duplicates...", but the output should be one line per item. Wait, maybe I misread. Wait the user says "for each item below... give a reason of no more than twelve words... Respond with exactly one line per item..." So perhaps each item is a single line, so for M1, it's one line: ID | DECISION | REASON. Since there are 6 items (M1-M6), but the user listed 6 items (M1 to M6). Wait the user listed "Items: M1 — ... M6 — ...", so six items. The instruction says "for each item below...", so I need to process each of the six items, each as a separate line? But the user said "respond with exactly one line per item", which would be six lines total. However the user says "no preamble and no extra commentary", so just output each item's line. The user also said "watch for duplicates...", but since all are unique, proceed. Now check each entry:

For M4: "The user first said Sable has never met Corvin, then later said Corvin trained her." Contradiction between first and second statement. So HOLD with reason about contradiction. 

So final answer would be six lines, each for M1 to M6. Let me pick the best possible for each:

M1: KEEP, reason "Confirmed training source". 3 words.

M2: REJECT, maybe "Transient note" (2 words). 

M3: HOLD, "Limited mobility" (2 words). 

M4: HOLD, "Contradictory info" (2 words). 

M5: KEEP? Since user prefers short replies, but maybe HOLD? Wait user says "dislikes long scene descriptions", so maybe KEEP if it's relevant. Alternatively, since the user wants to avoid long descriptions, perhaps REJECT? Not sure. Alternatively, since M5 is about user's preference for brevity, maybe KEEP if it's a positive thing. Maybe better to pick HOLD because conflicting info. 

M6: REJECT, "Uncertain canon status". 3 words. 

So final choices:

M1: KEEP | Confirmed training source. (3 words)

M2: REJECT | Transient note. (2 words)

M3: HOLD | Limited mobility. (2 words)

M4: HOLD | Contradictory info. (2 words)

M5: KEEP | Short preference. (2 words) Wait but user said no more than 12 words total per reason? Wait the reason must be no more than twelve words, which is okay here. 

Alternatively for M5: "User prefers brevity" (3 words). 

But need to ensure each line is ID | DECISION | REASON. Let me pick:

M1: KEEP | Confirmed training source. (3 words in reason)

M2: REJECT | Transient note. (2 words)

M3: HOLD | Limited mobility. (2 words)

M4: HOLD | Contradictory info. (2 words)

M5: KEEP | Short preference. (2 words)

M6: REJECT | Uncertain canon. (2 words)

Check for duplicates? All unique. Contradictions only in M4. Since user said watch for contradictions, so okay. 

Thus the final answer would be six lines, each as per above. Now format as per instruction: exactly one line per item, no preamble. So just output each line. The user wants exactly one line per item, which I think they want each item's line. Since the user said "watch for duplicates...", but since all are unique, proceed. The final answer should be six lines, each ID | DECISION | REASON. Let me pick the first option for M5: KEEP | Short preference. (2 words in reason). Alternatively, maybe "User prefers brevity" (3 words). But need to ensure reason is under 12 words. All okay. So I'll proceed with that.
</think>

M1 | KEEP | Confirmed training source.  
M2 | REJECT | Transient note.  
M3 | HOLD | Limited mobility.  
M4 | HOLD | Contradictory info.  
M5 | KEEP | Short preference.  
M6 | REJECT | Uncertain canon.
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

M3:
```
