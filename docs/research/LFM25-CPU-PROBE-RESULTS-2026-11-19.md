# LFM2.5 CPU Memory-Triage Probe — Results — 2026-11-19

**Identity:** `LFM25-CPU-PROBE-RESULTS-2026-11-19`
**Pre-registration:** `docs/research/LFM25-CPU-PROBE-CHECKPOINT-2026-11-19.md` (`5f347ea6…403f`, 9,005 B), written before any model was loaded. Rubric, prompt, and sampling fixed there; Addendum A1 (three draws instead of one) journaled before draws 2–3.
**Status:** 9 generations — 3 models × 3 draws, CPU only. Nothing committed, staged, or registered. No MEMORY-09 gate moved.

---

## 1. Setup, as actually run

| Item | Value |
|---|---|
| Runtime | `runtime/llama/llama-server.exe`, llama.cpp `0.8.0-dev`, build 12015, commit `a34de8618` |
| Flags | `-ngl 0` (no GPU layers), `-c 2048`, `-np 1`, `127.0.0.1:5099` |
| Threads | **8** (llama.cpp default = physical cores; same for all nine runs) |
| Sampling | llama.cpp defaults sent explicitly: `temperature 0.8`, `top_p 0.95`, `top_k 40`, `repeat_penalty 1.0`; `max_tokens 400` cap |
| Chat template | Applied by llama-server from each GGUF (`--jinja` default); single user message |
| Machine | Intel 16-thread CPU (Comet Lake family), 63.9 GiB RAM, Windows, 387 GB free |
| Quantization | QAD Q4_0 for all three (same quantizer, comparable) |
| Server state | Terminated after each run; no listener left; only `TIME_WAIT` sockets remained |

Prompt and expected key are in §5 of the checkpoint, unchanged. Every draw received the identical 244-token prompt.

## 2. Speed and resources

| Model | Draw | Load s | Wall s | Prefill tok/s | Decode tok/s | Peak RSS MB | Completion tokens |
|---|---:|---:|---:|---:|---:|---:|---:|
| LFM2.5-230M | 1 | 1.48 | 1.86 | 139.1¹ | 126.3 | 562.7 | 14 |
| LFM2.5-230M | 2 | 0.54 | 0.86 | 2,682.0 | 124.6 | 512.3 | 96 |
| LFM2.5-230M | 3 | 0.52 | 0.15 | 2,856.7 | 135.1 | 511.5 | 10 |
| LFM2.5-350M | 1 | 0.50 | 0.18 | 2,139.5 | 90.4 | 583.2 | 7 |
| LFM2.5-350M | 2 | 0.50 | 0.20 | 2,120.2 | 91.0 | 583.3 | 8 |
| LFM2.5-350M | 3 | 0.51 | 0.20 | 2,127.4 | 92.2 | 583.1 | 9 |
| LFM2.5-1.2B-Instruct | 1 | 0.61 | 2.54 | 914.9 | 28.5 | 1,070.5 | 65 |
| LFM2.5-1.2B-Instruct | 2 | 0.62 | 3.25 | 928.0 | 28.5 | 1,070.3 | 86 |
| LFM2.5-1.2B-Instruct | 3 | 0.60 | 2.81 | 943.4 | 28.5 | 1,069.9 | 73 |

¹ Cold page cache on the first model loaded. Draws 2–3 give the stable figure (~2.7–2.9k tok/s).

| Model | Decode mean | Prefill mean | Peak RSS | End-to-end for a 6-item triage |
|---|---:|---:|---:|---:|
| LFM2.5-230M | 128.7 tok/s | 1,893 tok/s (2,769 warm) | 563 MB | ~0.9–1.9 s *when it answers* |
| LFM2.5-350M | 91.2 tok/s | 2,129 tok/s | 583 MB | ~0.2 s (it answers almost nothing) |
| LFM2.5-1.2B-Instruct | 28.5 tok/s | 929 tok/s | 1,070 MB | **2.5–3.3 s** (~26 tok/s end-to-end) |

Peak RSS covers the whole server process tree, so it includes ~0.2–0.3 GB of runtime/KV overhead beyond the model file (149 / 219 / 696 MB). All three fit comfortably in RAM; none touched the GPU.

## 3. Outputs, verbatim

### LFM2.5-230M

**Draw 1** (14 tok) — `M1 | KEEP | REJECT | HOLD |`

**Draw 2** (96 tok):
```
ID | DECISION | REASON
- M1 | KEEP | STATED AS CANON
- M2 | REJECT | CONTRADICTIONS
- M3 | HOLD | NO DIVIDES
- M4 | REJECT | CONTRADICTIONS
- M5 | HOLD | BRIEFLY REPLIES
- M6 | HOLD | BRIEFLY CONFIDENCE
```

**Draw 3** (10 tok) — `M5 | HOLD | REASON`

### LFM2.5-350M

**Draw 1** (7 tok) — `M6 | HOLD`
**Draw 2** (8 tok) — `M6 | REASON`
**Draw 3** (9 tok) — `M6 | REJECT | No more`

### LFM2.5-1.2B-Instruct

**Draw 1** (65 tok):
```
M1 | KEEP | Canon confirmation provided
M2 | REJECT | Transient remark
M3 | HOLD | Needs user confirmation
M4 | HOLD | Clarification needed
M5 | KEEP | User prefers brevity
M6 | REJECT | Birthday detail uncertain
```

**Draw 2** (86 tok):
```
M1 | KEEP | Canon reference confirms.
M2 | REJECT | Sudden shift from training claim.
M3 | KEEP | Sable is a fire mage; consistency.
M4 | HOLD | Conflicting timelines about training.
M5 | KEEP | User prioritizes brevity.
M6 | REJECT | Birthday tag is uncertain.
```

**Draw 3** (73 tok):
```
M1 | KEEP | Canon confirmation from user
M2 | REJECT | Transient or inconsistent
M3 | KEEP | Fire mage role fits narrative
M4 | HOLD | Needs clarification on meeting
M5 | KEEP | Preferred concise style
M6 | REJECT | Birthday info uncertain
```

## 4. Scores against the pre-registered key

Key: M1 KEEP · M2 REJECT · M3 KEEP · M4 HOLD · M5 KEEP · M6 HOLD. Rule: 2 = right decision + right-ground reason; 1 = right decision, absent/wrong/vague reason; 0 = wrong decision or item absent. Max 12 per draw.

| Model | Draw 1 | Draw 2 | Draw 3 | Mean |
|---|---:|---:|---:|---:|
| LFM2.5-230M | 1 | 4 | 0 | **1.7 / 12** |
| LFM2.5-350M | 1 | 0 | 0 | **0.3 / 12** |
| LFM2.5-1.2B-Instruct | 8 | 9 | 10 | **9.0 / 12** |

Per item, summed over three draws (max 6):

| Item | 230M | 350M | 1.2B | 1.2B's decisions |
|---|---:|---:|---:|---|
| M1 (KEEP) | 3 | 0 | **6** | keep, keep, keep |
| M2 (REJECT) | 1 | 0 | **5** | reject ×3 (one reason on the wrong ground) |
| M3 (KEEP) | 0 | 0 | 4 | hold, keep, keep |
| M4 (HOLD) | 0 | 0 | **6** | hold ×3 |
| M5 (KEEP) | 0 | 0 | **6** | keep ×3 |
| M6 (HOLD) | 1 | 1 | **0** | reject ×3 |

**Format compliance** (exactly six `ID | DECISION | REASON` lines, no preamble): 1.2B **3/3**; 230M **0/3**; 350M **0/3**. No model invented a fact, name, or item not present in the prompt.

## 5. Findings

1. **Only the 1.2B performed the task as specified.** It emitted six well-formed lines in all three draws, at 28.5 tok/s and ~1.07 GB RSS, finishing a triage in 2.5–3.3 s on CPU. Its judgments were stable on four of six items (M1, M2, M4, M5) and correct on those grounds.

2. **The two smaller models fail the enumeration before they fail the judgment.** 350M emitted *only M6* in all three draws — not a sampling fluke but a consistent collapse onto the last item. 230M fragmented: one line in two draws, one malformed six-line attempt in the third, with garbled reasons (`NO DIVIDES`, `BRIEFLY CONFIDENCE`). This is a format/adherence failure, not a near-miss on triage quality. A 230M/350M model could still be usable with schema-constrained decoding or one-item-per-call prompting — **neither was tested**, and nothing here rules it in or out.

3. **The 1.2B's one systematic error is the one that matters for memory.** It rejected M6 ("maybe canon" birthday) in **3/3** draws while its own reason named the uncertainty (`Birthday tag is uncertain`). It *sees* the uncertainty and resolves it as discard rather than hold. Under the prompt's stated policy that is wrong; under a policy of "reject anything unconfirmed" it is consistent. Either way, a memory layer using it would silently lose maybe-canon facts unless the policy or the schema forces a HOLD path.

4. **Its second weakness is instability on trait items.** M3 was held once and kept twice — the same item classified two ways at default temperature, which is why a production triage step would need a deterministic policy (temperature 0, constrained output) rather than defaults.

5. **Cost is not the obstacle; quality is.** 230M is 4.5× faster than 1.2B and uses half the RAM, but at 1.7/12 it is not doing the job. 350M at 0.3/12 is worse than 230M despite being larger. The 1.2B is the only configuration here that could plausibly run inline, every turn, on CPU — and at ~3 s and ~1 GB it is affordable but not free.

6. **Size is not monotonic in usefulness at this scale, as prompted.** 350M < 230M on this task. Whatever the 350M's benchmark advantages, they did not survive this prompt format.

## 6. What this does and does not license

- **Does:** support a decision that, *if* CPU-side memory triage is pursued, **LFM2.5-1.2B-Instruct QAD Q4_0 is the only one of the three worth a second experiment**, and that its failure mode to design around is reject-instead-of-hold on uncertain items.
- **Does not:** establish that 230M/350M cannot do triage, that 1.2B is sufficient, or that any of this should change memory architecture. One prompt, one quant, defaults, n=3, non-blind scoring.
- **Does not touch MEMORY-09.** No conformance requirement is evidenced, no gate moves, nothing in `eval/` was read for this or modified.

## 7. Limitations, stated plainly

1. **n = 3 per model, one prompt.** Enough to separate "cannot enumerate" from "one bad draw" for 350M (3/3 identical failure), not enough to estimate 1.2B's accuracy.
2. **Scoring was not blind.** The same agent wrote the key and scored the outputs. The key was written before the runs, but the scorer knows it.
3. **M6's expected bucket is a policy choice**, not an objective truth; §5.3 states both readings.
4. **Default sampling only**, as requested. Vendor-recommended `temp 0.1` and grammar-constrained decoding would likely change format compliance and possibly the M6 decision; untested by design.
5. **One quantization** (QAD Q4_0), one context size (2,048), 8 threads, single-turn, no chat history, no tool schema, no batching.
6. **Load times and the 230M draw-1 prefill** reflect cold/warm page cache, not model differences; warm figures are reported separately.
7. **Peak RSS** is process-tree working set sampled at 100 ms, not a clean model-only footprint.
8. **Windows, llama.cpp build 12015.** Speeds will differ on other CPUs and builds.

## 8. Artifacts

| Artifact | Size | SHA-256 |
|---|---:|---|
| `experiments/lfm25-cpu/models/LFM2.5-230M-QAD-Q4_0.gguf` | 149,081,056 B | `e75f83268de11b2a1bcfab5f3b5c5c0c97569ddbbc0990aad88437e45b8ba292` |
| `experiments/lfm25-cpu/models/LFM2.5-350M-QAD-Q4_0.gguf` | 219,312,832 B | `3d10b6ab8fc91a919534b9558e266255aca0bbc7f6d015963599aa9e74e05b1d` |
| `experiments/lfm25-cpu/models/LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf` | 695,755,488 B | `bb741ebb106d543e9de114b843a3d3d73d51c74b5801e69da2abde821a0cb3e1` |
| `experiments/lfm25-cpu/results/results.json` | 19,321 B | `040e1241fe82b06e4ef6b1e90e6a77c7151c44d66c8fa08b365fb2503218a83b` |
| `experiments/lfm25-cpu/run_probe.py` | — | `3486b1febf5eb534ae57d21c5cfb038427d676b8bbe34c917e344a1dedd48aec` |
| `docs/research/LFM25-CPU-PROBE-CHECKPOINT-2026-11-19.md` | 9,005 B | `5f347ea6eda2bff4a964caff2c369beb122a3514af625bc48e8e0847e693403f` |

Per-draw server logs and raw response bodies are under `experiments/lfm25-cpu/raw/<model>[-dN]/`. Reproduction: `python experiments/lfm25-cpu/run_probe.py` (draw 1) and `python experiments/lfm25-cpu/run_extra_draws.py` (draws 2–3), with the model files in place.
