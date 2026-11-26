# MEMORY-09 S3 — Embedder Evaluation (selection only, nothing downloaded)

**Identity:** `MEMORY-09-S3-EMBEDDER-1` · **Date:** 2026-11-15 · **Status:** evaluation for owner approval. **No model obtained, no download, no network use.**
**Purpose:** choose the semantic branch's embedder against this experiment's requirements, not against a leaderboard.

---

## 1. What the app already has (and why it does not count)

| Code | What it actually is | Usable as the semantic branch? |
|---|---|---|
| `internal/backend/rag/embedder.go` → `LlamaEmbedder` | `/v1/embeddings` from a loaded llama-server model | Only if an **embedding-capable** model is loaded; a chat model's hidden states are not an embedding model |
| `internal/backend/rag/embedder.go` → `KeywordEmbedder` | keyword frequency vector | **No.** This is the fake-semantic pattern: lexical, mislabelled |
| `internal/backend/lorebook/embeddings.go` → `SimpleEmbedding` | TF-IDF term weighting | No — lexical |
| `internal/image/metadata/embedder.go` | image-side embedding | Not text |
| `config/model-registry.json` | chat/vision/image models; one `taesdxl` diffusion "embedding" | **No text embedding model installed** |

So: the app has no dedicated text embedder, and its fallback is precisely the thing §3 forbids. A new component is required either way.

## 2. Requirements (from the work order)

1. Local, offline inference.
2. Suited to **short factual records** and **conversational queries** (our records are one sentence; queries are one line).
3. Practical on a GTX 1080 8 GB, with CPU/RAM fallback. Note: the 1080 has no bf16/tensor cores; fp16 and quantised CPU paths matter more than raw benchmark rank.
4. Reproducible identity: model, revision, file hash, dimensions, prefixes, metric.
5. **No silent substitution** if the model fails to load — degrade loudly.

## 3. Candidates (figures verified against current sources, not recalled)

| Model | Params | Dim | Max tokens | Licence | Weight size | Notes |
|---|---:|---:|---:|---|---:|---|
| **nomic-embed-text-v1.5** | 137M | 768 | 8,192 | Apache-2.0 | 81–140 MB (Q4_K_M–Q8_0 GGUF) | GGUF for llama.cpp; **requires `search_document:` / `search_query:` prefixes**; quantisation MSE published |
| **Qwen3-Embedding-0.6B** | 596M | 1,024 | 32,768 | Apache-2.0 | ~610 MB (Q8_0 GGUF) | highest quality in its class; query needs `Instruct: <task>\nQuery:` prefix |
| bge-small-en-v1.5 | 33.4M | 384 | 512 | MIT | ~67–133 MB | cheapest CPU path; 512-token limit is fine for our records |
| bge-base-en-v1.5 | 109M | 768 | 512 | MIT | ~438 MB | stronger than small, same limits |
| granite-embedding-small-english-r2 | 47.7M | 384 | 8,192 | Apache-2.0 | ~284 MB (BF16) | long context, small |
| all-MiniLM-L6-v2 | 22.7M | 384 | 256 | Apache-2.0 | 46–90 MB | fastest; weakest retrieval quality |
| embeddinggemma-300m | ~300M | 768 | 2,048 | Gemma terms (gated) | ~318 MB | **excluded: gated licence** for an offline reproducible record |

## 4. Recommendation

**Primary: `nomic-embed-text-v1.5`, Q8_0 GGUF (~140 MB).**

Why it wins on *our* criteria rather than on MTEB: 8k context with no truncation risk for record+query pairs; Apache-2.0 with no gate; llama.cpp GGUF path already exists in this machine's toolchain (no new runtime); 768 dimensions is cheap to hold for 7 records and for hundreds later; the quantisation error is published per file, so identity can be pinned honestly. The prefix requirement is a documented hazard, not a surprise — and the adapter records prefixes in the identity so a missing prefix is a visible misconfiguration, not a silent quality loss.

**Challenger: `Qwen3-Embedding-0.6B`, Q8_0 (~610 MB).** Higher retrieval quality for a 4× size increase, still comfortable in 8 GB alongside an idle system. Worth running **if P5 remains unsolved** with nomic — that is a measured question, not a promise.

**Rejected as primary:** all-MiniLM (weakest quality, 256-token limit); bge-small (fine, but 512-token cap and no advantage here); embeddinggemma (gated licence undermines reproducible offline records).

**Important caveat, stated before measurement:** a **bi-encoder may not solve P5 at all.** "Who trained you?" and "Who trained the horse?" are near-identical in topic and vocabulary; a bi-encoder encodes each side independently and may place them close. Separating subject and relation is what a **cross-encoder reranker** is for. So the plan is: measure the bi-encoder on all probes; if P5 still fails, that is evidence for the reranker, not a reason to swap embedders again.

## 5. P5 acceptance criteria — written before the embedder is tried

So the target cannot move after the result:

1. P5 must **abstain**: no record selected, `reason` recorded.
2. P5 must abstain while **P1 still retrieves** `trained_by`. A system that abstains everywhere passes P5 and fails the point.
3. The margin between the top score and the floor is recorded, so "abstains by a hair" is distinguishable from "abstains confidently".
4. If P5 fails, the failure is reported as a **bi-encoder limitation**, and the reranker becomes the next component — not a silent threshold increase until the probe passes. Raising a threshold to make P5 pass while P1 survives is acceptable **only** if reported as a threshold change with both probes' margins shown.

## 6. Adapter contract (`eval/s3/embedder.py`)

| Requirement | Enforced by |
|---|---|
| Identity pinned from the artefact | `identity_for()` hashes the model file; identity is not read from config |
| Missing model is loud | `EmbedderUnavailable` raised; `require()` refuses to let fused scores be attributed to semantic |
| No fallback | no keyword path exists in the adapter; a degraded index stays degraded and visible |
| Prefix policy recorded | `query_prefix` / `document_prefix` are part of the identity record |
| Metric recorded | `metric` field, cosine for both candidates |

## 7. What I need from you

One decision, in one sentence: **authorise downloading `nomic-embed-text-v1.5` Q8_0 GGUF (~140 MB)** for the semantic branch — or name the challenger instead.

On approval I would: fetch exactly that file, record its SHA-256 and revision, wire it as the third branch, rerun all seven probes with branch attribution, and issue a Gate B receipt whose `semantic.available` is `true` — or a receipt that says plainly why the embedder did not close the gate.

**Nothing is downloaded, installed, or configured until you say so.**

## 8. Claim ceiling

This is a selection argument, not a measurement. No embedding model has been obtained, no vector has been computed, and no claim is made that any candidate will pass P5.
