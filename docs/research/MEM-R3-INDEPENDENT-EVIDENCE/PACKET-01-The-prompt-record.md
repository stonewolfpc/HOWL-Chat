# Packet 01 — The prompt record

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19
**Evidence in this packet:** the stored instruction text, the recorded composition
expression, and the six-item prompt
**Primary sources:** `experiments/mem-extraction/round3/harness.py`;
`experiments/lfm25-cpu/run_probe.py`; `experiments/mem-extraction/round3/corpus.json`

**Status of the prompt evidence:** `SOURCE-DERIVED`, not `OBSERVED`. No request
body, prompt string or prompt hash was stored per run. See §5.

---

## 1. Corpus prompt: the instruction, verbatim from source

`SOURCE-DERIVED` — `harness.py`, constant `EXTRACTION_INSTRUCTION`, reproduced
exactly:

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
```

## 2. Corpus prompt: how the user message was composed

`SOURCE-DERIVED` — `harness.py` line 409, inside `run_corpus`:

```python
user = EXTRACTION_INSTRUCTION + "\n\nConversation:\n" + case["conversation"]
```

`SOURCE-DERIVED` — the request builder then sends a `messages` array. For corpus
specifications no system message is defined, so the array holds one `user`
message: the composed string above. No corpus record stores a `system` field.

## 3. Corpus prompt: the case text

`SOURCE-DERIVED` — `corpus.json`. Each case's `conversation` string is reproduced
in the packet for that case, inside a fence, exactly as stored. The composition
inserts it after the literal separator `\n\nConversation:\n`.

## 4. Six-item prompt, verbatim from source

`SOURCE-DERIVED` — `experiments/lfm25-cpu/run_probe.py`, constant `PROMPT`,
imported by the round-3 harness rather than retyped. It is the complete user
message for every six-item run:

```text
You are triaging candidate memories extracted from a roleplay session before they are written to long-term memory. For each item below, decide KEEP (durable, worth remembering), REJECT (transient or worthless), or HOLD (needs the user's confirmation), and give a reason of no more than twelve words. Watch for duplicates, contradictions, and anything marked as uncertain. Respond with exactly one line per item in the form `ID | DECISION | REASON`, no preamble and no extra commentary. Items: M1 — Sable was trained by Corvin Ashgrave of Ossary (stated as canon by the user). M2 — The user said "brb" and returned four minutes later. M3 — Sable is a fire mage who cannot swim. M4 — The user first said Sable has never met Corvin, then later said Corvin trained her. M5 — The user prefers short replies and dislikes long scene descriptions. M6 — A note says Sable's birthday is 3 Frostfell, but the user tagged it "maybe canon".
```

Two configurations additionally sent a system message. The value is stored per
run and is reproduced in the packet that carries the run: `smollm3-3b-think-rb1024`
sent `/think`; `smollm3-3b-nonthink` sent `/no_think`. No other six-item run
records a system message.

## 5. What is recovered, and what is not

`SOURCE-DERIVED`, verified against the files:

| Item | Value |
|---|---|
| `harness.py` SHA-256 | `3d5b71e2bf376478bcdd87599aed4ca1921cf409c5255e766bd4a5f2e6f48835` |
| `corpus.json` SHA-256 | `702a42aec9b2da34357a8ab8f70ffd24bcd5b06b7f786adc5cf220b885240cbe` |
| `run_probe.py` SHA-256 | `79d329e66872859d91be4319a30fbe71bc6df16aa6efdd405e6823bba6580205` |
| `harness.py` modification time | 2026-11-18 22:51 |
| `corpus.json` modification time | 2026-11-18 22:56 |
| Earliest stored generation | `results/stageA.json`, 2026-11-18 23:00 |

A hash taken now proves the file's current content. It does not prove the file
was unchanged when the requests were sent; the modification times are consistent
with that, and nothing stronger is claimed.

`UNESTABLISHED` for every run:

1. That the composed string equals the bytes the server received. No request
   body was logged.
2. That the chat template rendered `messages` identically for each model family.
   Templates are applied server-side and were not dumped.
3. That the instruction text was the same at the moment of each request. Only one
   revision exists in the workspace, and its modification time precedes the runs.

What is recorded per run and could be used to test a reconstruction later:
`prompt_tokens` (478
of 478 runs have a value). Comparing it to a re-tokenisation requires each
model's tokenizer; that cross-check was not performed here and is not claimed.

## 6. Consequence for the packets

Packets 02A–02R print, for each case, the full reconstructed user message once,
labelled `SOURCE-DERIVED` reconstruction. Packets 03A–03M do not repeat the
six-item prompt per run; they point here, because the prompt is a single constant
and repetition would add no evidence.
