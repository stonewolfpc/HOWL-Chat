# MEMORY-09 S3 — Semantic Branch: Measured Execution Behaviour (corrected)

**Identity:** `MEMORY-09-S3-EMBEDDER-3` (rev 2) · **Date:** 2026-11-15
**Status: ACCEPTED as a research record (project owner, 2026-11-15).**
**Gate B: still PARTIAL, 6/7, and open.** Nothing here closes it.
**Receipts:** `eval/s3/measure-cpu.json`, `eval/s3/measure-gpu.json`, `eval/s3/embedder-receipt.json`
**Harness:** `eval/s3/measure_embedder.py` — reproducible, four phases kept separate.

---

## 0. Three corrections to my own earlier report

I published numbers that were wrong. They are corrected here, with the cause.

| Earlier claim | Corrected finding | Cause of the error |
|---|---|---|
| "GPU execution not established; `-ngl 99` shows no offload" | **Offload happens: 13/13 layers, `CUDA0 model buffer = 114.91 MiB`** | I grepped for placement lines at default verbosity. They only print at `-lv 4`. Absence of a log line is not absence of the behaviour. |
| "~251 MiB VRAM overhead even on CPU — target 0 missed" | That number was **CUDA context + weights on GPU, not CPU**. True CPU is **+94 MiB**, and that too is the *CUDA build's* context. | My "CPU" runs never disabled the GPU: llama-server **ignores `CUDA_VISIBLE_DEVICES=""`** and auto-offloads. Only `--device none -ngl 0` is a real CPU run. |
| "Cold load 16 s — budget FAIL" | **Spawn→ready 0.56 s; model load 0.14 s (CPU) / 0.24 s (GPU).** | 16 s was the first-ever read of a 146 MB file from cold disk. Subsequent starts are sub-second. The budget **passes**. |

## 1. Execution matrix, measured

Both runs: same model, same harness, idle VRAM ~1,585–1,605 MiB (median of 5 samples, 1 s apart).

| Phase | CPU (`--device none -ngl 0`) | GPU (`-ngl 99`) | Budget | Verdict |
|---|---:|---:|---|---|
| Spawn → ready | **0.56 s** | **0.56 s** | ≤ 3.0 s | **PASS** both |
| Model load | 0.14 s | 0.24 s | — | — |
| First query after load | 34.7 ms | 54.3 ms | — | cold-call cost, once |
| Warm p50 | **13.4 ms** | **10.0 ms** | ≤ 40 ms | **PASS** both |
| Warm p95 | **34.2 ms** | **29.4 ms** | ≤ 120 ms | **PASS** both |
| VRAM delta while loaded | **+94 MiB** | **+257 MiB** | ≤ 384 MiB | PASS both; CPU target 0 **missed** |
| Residual after shutdown | **−13 MiB** | **0 MiB** | ≤ 64 MiB | **PASS** both |
| Port closed after shutdown | true | true | — | **PASS** |

Weight placement verified from the logs, not inferred: CPU run `CPU_Mapped 138.65 MiB`, `offloaded 0/13`; GPU run `CUDA0 114.91 MiB` + `CPU_Mapped 23.76 MiB`, `offloaded 13/13`.

## 2. What this means for the three-choice setting

| Choice | Measured reality |
|---|---|
| **Off** | Works: degraded mode 4/4 with no server. |
| **Prefer GPU** | **Real and within budget**: +257 MiB, p50 10.0 ms. Safe when free VRAM ≥ ~400 MiB. Must fall back and *report* if not. |
| **CPU / Save VRAM** | +94 MiB, p50 13.4 ms. **Not zero** — because this binary has CUDA compiled in and initialises a context even with `--device none`. |

**Honest advertising:** "CPU / Save VRAM" should claim *"no model in VRAM, ~100 MiB fixed context"*. A CPU-only build of llama.cpp (same source, no CUDA) should take that to ~0, and that is the remaining test before the label says zero. GPU is only 3–4 ms faster at p50 on this Pascal card — the cost difference matters more than the speed difference.

## 3. Gate B, unchanged

6/7. P5 "Who trained the horse?" returns the trainer record from all three branches at once (semantic score 0.722 on the wrong record vs 0.750 on the right one). Threshold 0.60 sits in the measured gap for the other probes and deliberately does not rescue P5.

**Gate B stays open.**

## 4. Claim ceiling

Seven records. These are per-probe and per-process measurements on one machine, one runtime build, one model — not large-corpus retrieval quality, and not evidence about a conversation model's behaviour. Warm latency varied 10–26 ms p50 across runs; the tables above record the last run of each config from a fixed harness.
