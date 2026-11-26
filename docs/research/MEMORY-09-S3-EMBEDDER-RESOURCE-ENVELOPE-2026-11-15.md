# MEMORY-09 S3 — Embedder Resource Envelope (measured baseline, pre-registered budgets)

**Identity:** `MEMORY-09-S3-EMBEDDER-2` · **Date:** 2026-11-15 · **Status:** envelope fixed. **No model downloaded, no model loaded.**
**Purpose:** decide whether fetching an embedder is worth it, using budgets written *before* measurement.
**Baseline artefact:** `eval/s3/resource-baseline.json` (2,426 B)

---

## 1. This machine, measured

| Resource | Value |
|---|---|
| GPU | NVIDIA GTX 1080, **8,192 MiB**, driver 582.78 |
| Compute capability | **6.1 (Pascal)** — no tensor cores, poor fp16 throughput |
| VRAM idle | **1,656 MiB used**, 6,404 MiB free |
| VRAM with the S1 chat model resident | ~5,806 MiB used → **~2.25 GiB free** |
| RAM | **63.9 GiB**, 49.9 GiB available |
| Cores | 16 |
| Embedding runtime | `llama-server` **b10076** already present, `--embedding` and `--pooling` confirmed |

Two consequences follow immediately:

1. **Pascal has no useful fp16 path.** A GPU embedder will run fp32/fp16-emulated. For a 140 MB model that is fine; it also means GPU is not automatically faster than 16 CPU cores.
2. **During roleplay only ~2.25 GiB VRAM is free.** Any resident embedder competes with the model the user is actually running, which is exactly what the work order says to avoid.

## 2. The three guarantees, restated as testable properties

| # | Guarantee | Test | Status |
|---|---|---|---|
| G1 | Formation works with no local AI | Gate A deterministic pipeline, 14 checks | **PASS** |
| G2 | Semantic retrieval is optional; absence is declared, not faked | `resources.degraded_mode()` — 4/4 probes pass with `semantic.available=false` | **PASS (measured today)** |
| G3 | No permanent VRAM occupancy required | CPU execution is the default; GPU requires measured headroom | **PENDING** (needs a model) |

G2 is now a measurement, not a claim: with the semantic branch absent, P1, P2, P3b and P4a all still pass and the system reports itself degraded.

## 3. Pre-registered budgets

Written now so a disappointing result cannot be reframed as a passing one.

| Budget | Threshold | Why this number |
|---|---|---|
| VRAM overhead while chat model resident | **≤ 384 MiB**; **target 0** | Leaves ~1.87 GiB of the 2.25 GiB free; 0 is achievable with CPU |
| Cold load time | ≤ 3.0 s | Beyond this, on-demand embedding is unusable per query |
| Residual VRAM after unload | ≤ 64 MiB | Proves unload actually releases |
| Query latency p50 / p95 (short query) | ≤ 40 ms / ≤ 120 ms | A memory lookup must not be felt in chat |
| Index build, 100 records | ≤ 2.0 s | Bounded one-off cost |
| CPU RSS for embedding path | ≤ 1.0 GiB | Leaves the 64 GiB box untouched |

**Preferred execution: CPU/RAM.** GPU is a fallback, not a default — it is the only option that can violate G3.

## 4. Measurement matrix (eight cells)

| Axis | Values |
|---|---|
| Chat model | **absent** / **present** (the S1 Qwen3.5-9B profile) |
| Device | **CPU** / **GPU** |
| State | **cold** (just loaded) / **warm** (after 20 embeddings) |

Per cell, record: VRAM before/after load, VRAM after unload, cold load seconds, embedding latency for 1 / 10 / 100 short records, and process RSS.

## 5. Decision rule, fixed in advance

1. If **CPU** meets the latency budget → **adopt CPU**, zero VRAM, G3 satisfied by construction. Fetch is justified.
2. If CPU misses and **GPU delta ≤ 384 MiB** with cold load ≤ 3 s → GPU permitted **only** with explicit load/unload, never resident during chat.
3. If neither meets budget → **reject the candidate**; Gate B stays open and the limitation is recorded. No threshold is adjusted to rescue it.
4. Nomic and Qwen3-0.6B are run through the same matrix; the challenger is only fetched if nomic fails a budget **or** P5 remains unsolved, not because a leaderboard says so.

## 6. What this does not settle

- Whether an embedder closes P5. A bi-encoder may not; a cross-encoder reranker is the likely answer.
- Actual latency on this machine — that requires the model, which is not authorized.
- Whether the semantic branch improves retrieval *enough* to justify its complexity at 7 records. At this corpus size the honest expectation is **no measurable benefit**; the value appears at hundreds of records, which is why the budget test is about cost, not quality.

## 7. What I need from you

Nothing yet. The envelope is set; the decision on fetching waits on your call. If you authorize, the sequence is: fetch → pin SHA-256 and revision → run the eight-cell matrix → publish measured numbers against these budgets → only then judge whether nomic earns its place.

## 8. Claim ceiling

Hardware baseline is measured and reproducible. Every latency and VRAM figure for an embedder is **unmeasured**, because no embedding model exists on this machine. Nothing here is evidence about Nomic, Qwen3, or any embedder's quality.
