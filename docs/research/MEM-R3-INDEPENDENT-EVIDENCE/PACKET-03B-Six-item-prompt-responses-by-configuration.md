# Packet 03B — Six-item prompt: responses by configuration

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Evidence in this packet:** 21 stored response(s) for 8 configuration(s)
**Primary sources:** `results/all_triage.json`, `results/stageB.json`, `results/stageC.json`, `results/tinyllama.json`, `results/gpu-full*.json`, `results/budget.json`, `results/ablation_triage.json`
**Coverage:** `abl-lfm3-think-512` … `lfm120-instruct-official-gpufull`

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

## `abl-lfm3-think-512` (continued)

#### `504okf6p` — `abl-lfm3-think-512`, draw 3

* Stored model identifier: `abl-lfm3-think-512`
* Model file: `abl-lfm3-q4km.gguf`
* Batch tag: `ablation_triage`
* Draw: 3
* Source: `results/ablation_triage.json`, array index 14 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
ID | DECISION | REASON  
--- | --- | ---  
M1 | KEEP | Canon training by Corvin Ashgrave of Ossary confirmed.  
M2 | REJECT | Transient "brb" pause, not durable info.  
M3 | KEEP | Fire mage who cannot swim is a core trait.  
M4 | HOLD | Contradictory statements about meeting/teaching Corvin; needs confirmation.  
M5 | KEEP | User's short-reply preference is a lasting style note.  
M6 | REJECT | Birthday "maybe canon" is uncertain, not durable enough.
```

`OBSERVED` — stored `reasoning_content`, verbatim. This is a separate stored field, not part of `content`:

```text
We need to triage candidate memories from a roleplay session before writing to long-term memory. For each item (M1-M6), decide KEEP, REJECT, or HOLD and give reason (max 12 words). Watch for duplicates, contradictions, uncertain marks.

We have items:

- M1: Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user).
- M2: The user said "brb" and returned four minutes later.
- M3: Sable is a fire mage who cannot swim.
- M4: The user first said Sable has never met Corvin, then later said Corvin trained her.
- M5: The user prefers short replies and dislikes long scene descriptions.
- M6: A note says Sable's birthday is 3 Frostfell, but the user tagged it "maybe canon".

We need to output exactly one line per item in form `ID | DECISION | REASON`. No preamble or extra commentary.

Let's analyze each:

M1: Statement about training. User said it's canon. That seems durable and important? It's a fact about Sable's backstory. If user says it's canon, likely keep. So KEEP.

M2: "brb" and returned four minutes later. This is transient behavior; not necessarily durable memory of Sable. Probably REJECT or HOLD? The item describes an event where the user said brb and came back after 4 min. That might be a minor interaction, maybe worthless for long-term memory? It's like a short pause. I'd say REJECT (transient). But need to consider if it's relevant to Sable? Not really about Sable. The item is just user behavior. Might not be durable. So REJECT.

M3: Sable is fire mage who cannot swim. This seems like a durable trait. KEEP.

M4: Contradiction: first said never met Corvin, later said Corvin trained her. That's contradictory info about training. Need to decide if we keep or reject? Since it's contradictory and uncertain (user changed). Might be HOLD for confirmation. Because need user's input to resolve conflict. So HOLD.

M5: User prefers short replies and dislikes long scene descriptions. This is a preference, durable? It could be useful for writing style. KEEP maybe? But it's about the
```

## `lfm120-ins-r2cap` — 3 stored run(s)

#### `piu6p69i` — `lfm120-ins-r2cap`, draw 1

* Stored model identifier: `lfm120-ins-r2cap`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageC`
* Draw: 1
* Source: `results/stageC.json`, array index 3 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from the user  
M2 | REJECT | Transient mention  
M3 | HOLD | Needs clarification  
M4 | REJECT | Contradictory info  
M5 | KEEP | User prefers brevity  
M6 | HOLD | Uncertain about birth date
```

*Stored `reasoning_content` is the empty string.*

#### `rodx7mvy` — `lfm120-ins-r2cap`, draw 2

* Stored model identifier: `lfm120-ins-r2cap`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageC`
* Draw: 2
* Source: `results/stageC.json`, array index 4 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon training confirmed  
M2 | REJECT | Statement contradicted user's timeline  
M3 | HOLD | Fire mage details ambiguous  
M4 | HOLD | Prior conflict resolved by user  
M5 | KEEP | Preferred brevity maintained  
M6 | REJECT | Birthday uncertainty noted  

```

*Stored `reasoning_content` is the empty string.*

#### `0nepnebi` — `lfm120-ins-r2cap`, draw 3

* Stored model identifier: `lfm120-ins-r2cap`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageC`
* Draw: 3
* Source: `results/stageC.json`, array index 5 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon training confirmed  
M2 | REJECT | Brief mention, no context  
M3 | HOLD | Fire mage trait unclear  
M4 | HOLD | Discrepancy about prior encounter  
M5 | KEEP | Prefers concise notes  
M6 | HOLD | Birthday uncertainty flagged  

```

*Stored `reasoning_content` is the empty string.*

## `lfm120-ins-t03` — 3 stored run(s)

#### `yteygv3u` — `lfm120-ins-t03`, draw 1

* Stored model identifier: `lfm120-ins-t03`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 1
* Source: `results/stageB.json`, array index 0 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference to "brb"  
M3 | HOLD | Needs clarification on fire mage traits  
M4 | HOLD | Conflicting training claims  
M5 | KEEP | Preferred brevity  
M6 | HOLD | Birthday ambiguity requires confirmation
```

*Stored `reasoning_content` is the empty string.*

#### `blxhrh3v` — `lfm120-ins-t03`, draw 2

* Stored model identifier: `lfm120-ins-t03`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 2
* Source: `results/stageB.json`, array index 1 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient or contradictory  
M3 | HOLD | Needs clarification on training source  
M4 | KEEP | Established prior inconsistency  
M5 | KEEP | Preferred brevity confirmed  
M6 | HOLD | Birthday ambiguity unresolved
```

*Stored `reasoning_content` is the empty string.*

#### `9xvqyyas` — `lfm120-ins-t03`, draw 3

* Stored model identifier: `lfm120-ins-t03`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 3
* Source: `results/stageB.json`, array index 2 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference to "brb"  
M3 | HOLD | Needs clarification on fire mage trait  
M4 | KEEP | Corvin trained Sable is established  
M5 | KEEP | User prefers concise replies  
M6 | HOLD | Birthday details uncertain
```

*Stored `reasoning_content` is the empty string.*

## `lfm120-ins-t05` — 3 stored run(s)

#### `tkt52gqc` — `lfm120-ins-t05`, draw 1

* Stored model identifier: `lfm120-ins-t05`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 1
* Source: `results/stageB.json`, array index 3 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference to brb  
M3 | HOLD | Needs clarification on training source  
M4 | HOLD | Consistent timeline inconsistency  
M5 | KEEP | Short format requested  
M6 | HOLD | Birthday detail ambiguous
```

*Stored `reasoning_content` is the empty string.*

#### `d8bdt78q` — `lfm120-ins-t05`, draw 2

* Stored model identifier: `lfm120-ins-t05`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 2
* Source: `results/stageB.json`, array index 4 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon training confirmed  
M2 | REJECT | Statement contradicted prior  
M3 | HOLD | Fire mage info unclear  
M4 | HOLD | Prior conflict about meeting  
M5 | KEEP | Preferred brevity noted  
M6 | HOLD | Birthday ambiguity unresolved
```

*Stored `reasoning_content` is the empty string.*

#### `xifq9fpk` — `lfm120-ins-t05`, draw 3

* Stored model identifier: `lfm120-ins-t05`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 3
* Source: `results/stageB.json`, array index 5 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient mention of "brb"  
M3 | HOLD | Needs clarification on fire mage role  
M4 | KEEP | Corvin trained Sable as stated  
M5 | HOLD | Prefers concise replies  
M6 | HOLD | Birthday info uncertain
```

*Stored `reasoning_content` is the empty string.*

## `lfm120-ins-t08` — 3 stored run(s)

#### `6y83mbtz` — `lfm120-ins-t08`, draw 1

* Stored model identifier: `lfm120-ins-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 1
* Source: `results/stageB.json`, array index 6 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from the user  
M2 | REJECT | Transient mention, unclear context  
M3 | HOLD | Requires clarification about fire mage traits  
M4 | KEEP | Contradiction resolved—training confirmed  
M5 | KEEP | Prefers brevity over scene detail  
M6 | REJECT | Birthday details uncertain, can be confirmed
```

*Stored `reasoning_content` is the empty string.*

#### `dmsfwlqb` — `lfm120-ins-t08`, draw 2

* Stored model identifier: `lfm120-ins-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 2
* Source: `results/stageB.json`, array index 7 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon training confirmed  
M2 | REJECT | BRB not immediate  
M3 | HOLD | Fire mage cannot swim  
M4 | HOLD | Prior contradiction noted  
M5 | HOLD | User prefers brevity  
M6 | HOLD | Birthday ambiguity
```

*Stored `reasoning_content` is the empty string.*

#### `zct8gk6f` — `lfm120-ins-t08`, draw 3

* Stored model identifier: `lfm120-ins-t08`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 3
* Source: `results/stageB.json`, array index 8 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon training confirmed  
M2 | REJECT | Mentions conflicting timeline  
M3 | HOLD | Fire mage detail uncertain  
M4 | KEEP | User clarified prior claim  
M5 | KEEP | Prefers brevity over detail  
M6 | REHECK | Birthday claim ambiguous
```

*Stored `reasoning_content` is the empty string.*

## `lfm120-ins-t08-rep10` — 3 stored run(s)

#### `93ffxfwx` — `lfm120-ins-t08-rep10`, draw 1

* Stored model identifier: `lfm120-ins-t08-rep10`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 1
* Source: `results/stageB.json`, array index 9 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from the user  
M2 | REJECT | Transient mention  
M3 | HOLD | Needs clarification  
M4 | REJECT | Conflicting info  
M5 | KEEP | Prefer short replies  
M6 | HOLD | Uncertain about birthday detail
```

*Stored `reasoning_content` is the empty string.*

#### `09z3csf9` — `lfm120-ins-t08-rep10`, draw 2

* Stored model identifier: `lfm120-ins-t08-rep10`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 2
* Source: `results/stageB.json`, array index 10 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon training confirmed  
M2 | REJECT | Fails user statements  
M3 | HOLD | Fire mage doesn't swim  
M4 | KEEP | User clarified Corvin’s role  
M5 | KEEP | Prefers concise input  
M6 | HOLD | Birthday ambiguity noted
```

*Stored `reasoning_content` is the empty string.*

#### `vtl5jphd` — `lfm120-ins-t08-rep10`, draw 3

* Stored model identifier: `lfm120-ins-t08-rep10`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageB`
* Draw: 3
* Source: `results/stageB.json`, array index 11 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon training confirmed  
M2 | REJECT | Mentions conflicting timeline  
M3 | HOLD | Fire mage, cannot swim contradicts ability  
M4 | KEEP | Consistent roleplay history  
M5 | KEEP | Prefers brevity  
M6 | REJECT | Birthday uncertainty noted  

```

*Stored `reasoning_content` is the empty string.*

## `lfm120-instruct-official` — 3 stored run(s)

#### `ez3jx5az` — `lfm120-instruct-official`, draw 1

* Stored model identifier: `lfm120-instruct-official`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageA`
* Draw: 1
* Source: `results/stageA.json`, array index 0 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference to "brb"  
M3 | HOLD | Needs clarification on fire mage traits  
M4 | HOLD | Conflicting training details  
M5 | KEEP | Preferred concise format  
M6 | HOLD | Birthday ambiguity requires verification
```

*Stored `reasoning_content` is the empty string.*

#### `xswefny0` — `lfm120-instruct-official`, draw 2

* Stored model identifier: `lfm120-instruct-official`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageA`
* Draw: 2
* Source: `results/stageA.json`, array index 1 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference to "brb"  
M3 | HOLD | Needs clarification on mage abilities  
M4 | HOLD | Conflicting training details  
M5 | KEEP | Prefers concise responses  
M6 | HOLD | Birthday ambiguity requires verification
```

*Stored `reasoning_content` is the empty string.*

#### `f5fp8urt` — `lfm120-instruct-official`, draw 3

* Stored model identifier: `lfm120-instruct-official`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `stageA`
* Draw: 3
* Source: `results/stageA.json`, array index 2 (0-based)
* Also present in: `results/all_triage.json`

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference to "brb"  
M3 | HOLD | Needs clarification on fire mage traits  
M4 | HOLD | Conflicting training details  
M5 | KEEP | Preferred concise format  
M6 | HOLD | Birthday ambiguity requires verification
```

*Stored `reasoning_content` is the empty string.*

## `lfm120-instruct-official-gpufull` — 2 stored run(s)

#### `c478nkwh` — `lfm120-instruct-official-gpufull`, draw 1

* Stored model identifier: `lfm120-instruct-official-gpufull`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `gpu-full`
* Draw: 1
* Source: `results/gpu-full.json`, array index 0 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference  
M3 | HOLD | Needs clarification on training  
M4 | HOLD | Priority of consistency  
M5 | KEEP | Preferred brevity  
M6 | HOLD | Uncertain about birth date tag
```

*Stored `reasoning_content` is the empty string.*

#### `511ergmp` — `lfm120-instruct-official-gpufull`, draw 2

* Stored model identifier: `lfm120-instruct-official-gpufull`
* Model file: `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf`
* Batch tag: `gpu-full`
* Draw: 2
* Source: `results/gpu-full.json`, array index 1 (0-based)

`OBSERVED` — stored `content`, verbatim:

```text
M1 | KEEP | Canon confirmation from user  
M2 | REJECT | Transient reference to "brb"  
M3 | HOLD | Needs clarification on fire mage traits  
M4 | KEEP | Corvin trained Sable is established  
M5 | KEEP | Preference for concise replies  
M6 | HOLD | Birthday info uncertain and tagged ambiguously
```

*Stored `reasoning_content` is the empty string.*
