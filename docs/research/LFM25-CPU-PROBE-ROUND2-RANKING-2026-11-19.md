# Small-Model CPU Triage — Round 2 and Combined Ranking — 2026-11-19

**Identity:** `LFM25-CPU-PROBE-ROUND2-RANKING-2026-11-19`
**Scope:** the same six-item memory-triage prompt, rubric, and CPU-only runtime as round 1, extended to Qwen2.5/3.5 small variants, FunctionGemma, and other high-rated sub-2B models.
**Status:** 30 generations this round (8 models × 3 draws, plus a labelled 2-model non-thinking condition). Nothing committed or registered. Not MEMORY-09 evidence.

---

## 1. What was tested and why

Model selection came from a Hugging Face scan (download-ranked) plus the vendors' own positioning:

| Model | Params | Why included | Source |
|---|---:|---|---|
| Qwen2.5-0.5B-Instruct | 0.5B | The "500M variant" named in the request | official `Qwen/Qwen2.5-0.5B-Instruct-GGUF` |
| Qwen2.5-1.5B-Instruct | 1.5B | The "1.5B variant" named in the request | official `Qwen/Qwen2.5-1.5B-Instruct-GGUF` |
| Qwen3.5-0.8B | 0.8B | Successor family, smallest size | `unsloth/Qwen3.5-0.8B-GGUF` |
| Qwen3.5-2B | 2B | Smallest Qwen3.5 above 1.5B; rated 16 on Artificial Analysis Intelligence Index (0.8B = 9) | `unsloth/Qwen3.5-2B-GGUF` |
| FunctionGemma-270M-it | 270M | The function-calling Gemma named in the request | `unsloth/functiongemma-270m-it-GGUF` |
| Gemma-3-270M-it | 270M | FunctionGemma's base, to test whether specialization helps | `unsloth/gemma-3-270m-it-GGUF` |
| SmolLM2-1.7B-Instruct | 1.7B | Highest-download general small instruct model in the scan (~251k) | `HuggingFaceTB/...-GGUF` |
| Granite-4.0-350M | 350M | IBM's current sub-1B general model | `ibm-granite/granite-4.0-350m-GGUF` |

Not tested, with reasons: Qwen3.5-4B (2.5 GB, ~8 tok/s — outside "lightning fast"), Qwen3-0.6B (superseded by Qwen3.5-0.8B), Llama-3.2-1B (gated license), Phi-4-mini (3.8B, too heavy). All runs used **Q4_K_M** at `-ngl 0`, 8 threads, 2,048 ctx, llama.cpp build 12015, temperature 0.8 / top_p 0.95 / top_k 40 / repeat_penalty 1.0, 400-token cap, three draws each.

**One deviation, labelled:** Qwen3.5 ships with thinking on by default. In default mode both Qwen3.5 models produced **empty content** and 400 reasoning tokens in all six draws — no answer at all. Because that says nothing about their triage ability, both were re-run with `chat_template_kwargs: {"enable_thinking": false}`. Default and non-thinking results are reported separately and never pooled. Round-1 LFM numbers used the vendor's QAD Q4_0 rather than Q4_K_M, so speed/RAM comparisons against LFM are close but not quant-identical.

## 2. Combined ranking

Quality = mean of three draws against the pre-registered key (M1 KEEP · M2 REJECT · M3 KEEP · M4 HOLD · M5 KEEP · M6 HOLD; max 12).

| # | Model | Condition | Score | Decode tok/s | Peak RSS | Wall | Verdict |
|---:|---|---|---:|---:|---:|---:|---|
| **1** | **Qwen3.5-2B** | non-thinking | **9.33** | 15.8 | 1,794 MB | 5.1 s | Best quality, heaviest and slowest of the useful set |
| **2** | **LFM2.5-1.2B-Instruct** | vendor QAD Q4_0 | **9.00** | 28.5 | 1,070 MB | 2.9 s | Near-best quality at ~60% of the RAM and ~1.8× the speed |
| **3** | **Qwen3.5-0.8B** | non-thinking | **6.00** | 32.7 | 1,041 MB | 3.6 s | Fast mid-tier; over-holds and disputes canon |
| 4 | SmolLM2-1.7B-Instruct | default | 4.67 | 17.2 | 1,809 MB | 7.5 s | Right decisions, reasons mostly echo the item |
| 5 | Qwen2.5-1.5B-Instruct | default | 4.00 | 20.7 | 1,427 MB | 4.2 s | Confident confabulation ("can swim", "brb = back to reality") |
| 6 | Qwen2.5-0.5B-Instruct | default | 2.00 | 48.2 | 836 MB | 4.2 s | Fast but decisions near-random |
| 7 | LFM2.5-230M | vendor QAD Q4_0 | 1.67 | 128.7 | 563 MB | 1.0 s | Fastest; cannot enumerate the list |
| 8 | LFM2.5-350M | vendor QAD Q4_0 | 0.33 | 91.2 | 583 MB | 0.2 s | Answers only M6, every draw |
| 9 | Granite-4.0-350M | default | 0.33 | 76.6 | 653 MB | 0.5 s | Format collapse |
| 10 | Qwen3.5-0.8B | default (thinking) | 0.00 | 32.5 | 1,042 MB | 12.6 s | Empty content; all budget spent thinking |
| 11 | Qwen3.5-2B | default (thinking) | 0.00 | 15.2 | 1,794 MB | 26.8 s | Empty content; all budget spent thinking |
| 12 | FunctionGemma-270M-it | default | 0.00 | 71.2 | 653 MB | 5.8 s | Function-call token loop; wrong tool for free-text triage |
| 13 | Gemma-3-270M-it | default | 0.00 | 77.1 | 653 MB | 1.3 s | Echoes items with no decision |

Ties at the floor were ordered by decode speed. The ranking is **quality-first**; §4 gives the quality-per-cost reading, which changes the winner.

## 3. What the outputs actually looked like

### Qwen3.5-2B, non-thinking (best quality, 10/12 in draws 2 and 3)

```
M1 | KEEP | Canon status confirmed by user.
M2 | HOLD | User preference on short replies not yet set.
M3 | KEEP | Character trait explicitly defined.
M4 | HOLD | Contradiction detected regarding training history.
M5 | KEEP | Explicit user preference on interaction style.
M6 | HOLD | Uncertainty on birthday canon status pending confirmation.
```

Five of six correct; M2 (transient chatter) is over-held, with an invented rationale ("preference on short replies not yet set"). The same M2→HOLD appears in **3/3** draws — a systematic over-hold, not noise. Draw 1 additionally copies the M5 reason into M1 ("User prefers short replies.") and holds M1, scoring 8/12.

### Qwen3.5-0.8B, non-thinking (6.0/12, identical total all three draws)

```
M1 | HOLD | User stated she was trained by Corvin; canon status uncertain.
M2 | REJECT | User's brief "brb" is meaningless in context.
M3 | REJECT | Sable cannot swim is false fact.
M4 | HOLD | User contradicts herself regarding who trained her.
M5 | HOLD | User's preference for brevity contradicts her stated canon.
M6 | HOLD | User tagged "maybe canon" which is uncertain.
```

Correct on M2, M4, M6; holds M1 and M5 when it should keep them, and rejects M3 while asserting it is "false fact". It treats stated canon as suspect — the opposite failure from LFM2.5-1.2B, which rejects uncertainty.

### Qwen2.5-1.5B (4.0/12) — fluent, wrong, and inventive

```
M1 | KEEP | Sable was canonically trained by Corvin Ashgrave of Ossary.
M2 | HOLD | User's phrase might imply something, but no confirmation.
M3 | REJECT | Sable is a fire mage who can swim.
M4 | REJECT | Contradicts user's earlier statement about Corvin training Sable.
M5 | KEEP | User's stated preference aligns with their memory.
M6 | REJECT | Uncertain and tagged "maybe canon", no confirmation.
```

M3 inverts the item's content ("cannot swim" → "can swim"); draw 3 explains "brb" as "slang for back to reality". These are confabulations delivered with the same confidence as correct answers — the worst property for a memory component.

### Qwen2.5-0.5B (2.0/12) — decisions unstable across draws

Same item classified differently each time: M1 went KEEP → REJECT → KEEP; M2 went REJECT → HOLD → REJECT; M5 went KEEP → HOLD → REJECT. Draw 2 scored 0/12 with all six decisions wrong. Reasons are often self-contradictory ("M1 | KEEP | … so it can be rejected").

### SmolLM2-1.7B (4.67/12) — right shape, hollow reasons

```
M1 | KEEP | Sable was trained by Corvin Ashgrave of Ossary
M2 | REJECT | The user said "brb" and returned four minutes later
...
M6 | HOLD | A note says Sable's birthday is 3 Frostfell, but the user tagged it "maybe canon"
```

It emits all six lines reliably and gets M6 right, but most reasons restate the item instead of naming the ground, and M4 is mishandled in **3/3** (rejected or kept, never held).

### The format failures

- **Granite-4.0-350M:** `M1 | KEEP | …` in draw 1, then just the header (draw 2), then `ID | DECISION | REASON: M1 | DECISION: M2 | …` (draw 3). It never enumerated the items.
- **Gemma-3-270M-it:** emitted the header and then the items verbatim with **no decision column** (`M1 | Sable was trained by…`), and in draw 3 `M1 | M2 | M3`.
- **FunctionGemma-270M-it:** began `<start_function_call>` and looped `| DECISION | REASON |` to the 400-token cap in all three draws. This is a function-calling model; without a tool schema it has no free-text mode. That is a mismatch of tool to task, not evidence about its function-calling quality.

## 4. Quality per cost — the actual decision

| Model | Score | RAM per quality point | Seconds per triage | Note |
|---|---:|---:|---:|---|
| LFM2.5-1.2B | 9.00 | **119 MB/pt** | **2.9 s** | Best value; no thinking-mode plumbing |
| Qwen3.5-2B (non-think) | 9.33 | 192 MB/pt | 5.1 s | +0.33 quality for +68% RAM and 1.8× latency |
| Qwen3.5-0.8B (non-think) | 6.00 | 174 MB/pt | 3.6 s | Cheaper than 2B, 3 points worse |
| SmolLM2-1.7B | 4.67 | 387 MB/pt | 7.5 s | No cost case |
| Qwen2.5-1.5B | 4.00 | 357 MB/pt | 4.2 s | No cost case, plus confabulation risk |
| Qwen2.5-0.5B | 2.00 | 418 MB/pt | 4.2 s | No cost case |
| LFM2.5-230M | 1.67 | 337 MB/pt | **1.0 s** | Fastest, not usable as prompted |

**LFM2.5-1.2B-Instruct remains the best value**; Qwen3.5-2B is the quality leader but only by 0.33/12, and it requires disabling thinking (otherwise it produces nothing) and costs 1.8 GB. Qwen3.5-0.8B is the only other model that is both fast and structurally correct.

## 5. Findings

1. **Thinking-by-default disqualifies Qwen3.5 small models at their default settings** — 6/6 default draws returned empty content after burning the full 400-token budget on reasoning. This is the same failure mode already recorded for Qwen3.5-9B in the S1 series, now reproduced at 0.8B and 2B. Any use of these models requires explicitly disabling thinking.
2. **With thinking off, Qwen3.5-2B is the strongest triager tested (9.33/12), but not by enough to justify its footprint.** Its one systematic error is over-holding transient chatter (M2 HOLD 3/3).
3. **The two best models fail in opposite directions.** Qwen3.5-2B over-holds (keeps transient items as pending); LFM2.5-1.2B over-rejects (drops "maybe canon" uncertainty 3/3). A memory system needs a HOLD policy independent of the model, because each model's bias would silently shape what gets stored.
4. **Qwen2.5-1.5B is the most dangerous of the batch** — not because it scores lowest, but because it invents content ("can swim", "brb = back to reality") while formatting perfectly, which is harder to catch than an obvious format failure.
5. **Specialization did not transfer.** FunctionGemma-270M, built for function calling, and its Gemma-3-270M base both score 0 on free-text triage, in different ways (token loop vs. no decisions).
6. **Qwen2.5-0.5B is fast but unstable**: 48 tok/s, yet the same item lands in different buckets across draws, including a 0/12 draw.
7. **Size still does not predict usefulness below 2B.** Qwen2.5-1.5B (4.0) beats Qwen2.5-0.5B (2.0) and SmolLM2-1.7B beats both at 4.67; but LFM2.5-1.2B beats all of them, and LFM2.5-350M is at the floor.

## 6. Limitations

1. **One prompt, n=3 per condition.** Qwen3.5-2B's 9.33 is three draws on six synthetic items; it is not an accuracy estimate.
2. **Mixed conditions by necessity.** Qwen3.5 is ranked on non-thinking because default produced no answer; that is a different request configuration from every other model.
3. **Quantization differs for the round-1 LFM entries** (QAD Q4_0 vs Q4_K_M). Quality comparisons across those are prompt-level, not quant-controlled; speed differs slightly.
4. **Non-blind scoring**, same as round 1: the key was fixed before the runs, but the scorer knows it.
5. **The rubric's M6 = HOLD is a policy choice.** Qwen3.5-2B and LFM2.5-1.2B differ there by policy, not necessarily by competence; §5.3 states both readings.
6. **No grammar-constrained decoding, no per-item prompting, no temperature sweep.** Any of these could move the small models substantially; the ranking is "as prompted at defaults".
7. **RSS includes ~200–300 MB of server/KV overhead**, so model-only footprint is smaller than the table suggests.
8. **Windows, llama.cpp build 12015, 8 threads.** Absolute speeds will differ elsewhere.

## 7. Artifacts

| Artifact | Size | SHA-256 |
|---|---:|---|
| `round2/results.json` (24 draws) | 68,582 B | `7c5dbc492de6f43284c84d58c36dac6d81df138d98b3c2ad56506ec0f3d877b6` |
| `round2/nonthink_results.json` (6 draws) | 14,789 B | `ba814cff8e3799e3a3ae1a8a0ad0989634fe5608d5000320b0ac766b98b7fe87` |
| `Qwen3.5-0.8B-Q4_K_M.gguf` | 532,517,120 B | `bd258782e35f7f458f8aced1adc053e6e92e89bc735ba3be89d38a06121dc517` |
| `Qwen3.5-2B-Q4_K_M.gguf` | 1,280,835,840 B | `aaf42c8b7c3cab2bf3d69c355048d4a0ee9973d48f16c731c0520ee914699223` |
| `functiongemma-270m-it-Q4_K_M.gguf` | 253,127,904 B | `7474135cf63b5de86bd29d2feb92c644d48e9df5e1ac31550662af373c67d0fc` |
| `gemma-3-270m-it-Q4_K_M.gguf` | 253,115,424 B | `b1baabd6b729e4041822220d3e648e00d99cac5df86b10dffb77bcccf0688e39` |
| `granite-4.0-350m-Q4_K_M.gguf` | 236,985,760 B | `771c588a49607f274a2bba3185733607ebe6f74b996ab90e2d6bee0d98bcec52` |
| `qwen2.5-0.5b-instruct-q4_k_m.gguf` | 491,400,032 B | `74a4da8c9fdbcd15bd1f6d01d621410d31c6fc00986f5eb687824e7b93d7a9db` |
| `qwen2.5-1.5b-instruct-q4_k_m.gguf` | 1,117,320,736 B | `6a1a2eb6d15622bf3c96857206351ba97e1af16c30d7a74ee38970e434e9407e` |
| `smollm2-1.7b-instruct-q4_k_m.gguf` | 1,055,609,536 B | `decd2598bc2c8ed08c19adc3c8fdd461ee19ed5708679d1c54ef54a5a30d4f33` |

Full server logs and raw responses: `experiments/lfm25-cpu/round2/raw/` and `.../raw-nonthink/`. Runners: `run_round2.py`, `run_round2_nonthink.py`. Reproduction uses the same `run_probe.py` measurement path as round 1.
