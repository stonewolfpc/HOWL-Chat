# MEMORY-09 S3 — Reranker Candidate Enumeration & Compatibility Review

**Identity:** `MEMORY-09-S3-RERANKER-CANDIDATES-1` · **Date:** 2026-11-15
**Status:** review only. **Nothing downloaded. No candidate approved.** Gate B open at 6/7.

---

## 0. The apparatus is frozen

| Artefact | SHA-256 (16) |
|---|---|
| `eval/s3/rerank-spec.json` v1.1 (cases, outcomes, threshold 0.30, vetoes) | `bbf400c2…` |
| case payload hash | `ca6dd5b3541f08fb…` |
| `eval/s3/rerank_eval.py` | `f59959f3…` |
| `eval/s3/rerank_set.py` | `23badd5e…` |
| last harness receipt | `5b8da5be…` |

No further change to the harness, case set, threshold or veto logic unless a **demonstrated defect** requires it. A candidate faces these exact cases.

## 1. Runtime compatibility, verified locally

`llama-server` b10076 (305ba519a) on this machine **already supports** the reranking path — no new runtime needed:

```
--rerank, --reranking        enable reranking endpoint on server
--pooling {none,mean,cls,last,rank}
```

Verified by `--help` on the installed binary, not from documentation.

**One build-age caveat that matters.** llama.cpp release **b11223 (2026-09-27)** added *"server: allow RANK pooling batch splitting for causal LLM rerankers (Qwen3, Qwen3-VL)"*. Our **b10076 predates it**. Bidirectional cross-encoders (BERT/XLM-R) require all tokens in one physical batch anyway; causal rerankers are the ones that need that splitting. So Qwen3-Reranker carries a compatibility risk on this build that bge does not.

## 2. Candidates, with measured metadata

Params and licences read from the Hugging Face API; GGUF sizes read from the repo tree. No weights fetched.

| Candidate | Architecture | Params | Licence | GGUF available | Smallest / recommended |
|---|---|---:|---|---|---|
| **bge-reranker-v2-m3** | XLM-R, bidirectional | 567.8 M | **Apache-2.0** | yes (`smarttasks`, Q4/Q5/Q8) | 438.4 MB Q4_K_M · 635.7 MB Q8_0 |
| **Qwen3-Reranker-0.6B** | Qwen3, causal | 0.6 B | **Apache-2.0** | yes, incl. maintainer build (`ggml-org`) | 639.2 MB Q8_0 |
| bge-reranker-base | XLM-R | 278.0 M | MIT | no official GGUF — conversion needed | — |
| ms-marco-MiniLM-L6-v2 | BERT, English | 22.7 M | Apache-2.0 | no official GGUF — conversion needed | — |
| jina-reranker-v2-base-multilingual | — | 278.4 M | **CC-BY-NC-4.0** | community only | **non-commercial — see §4** |
| mxbai-rerank-large-v2 | — | 1.54 B | Apache-2.0 | no official GGUF | 1.5 GB+ |

## 3. Envelope consequence, stated precisely (corrected)

My first version said file size makes GPU residency impossible. **File size is not runtime VRAM.** What the arithmetic actually supports is narrower, and the four cases must be kept apart:

| Case | Status |
|---|---|
| **Full GPU residency inside the 384 MiB budget** | **Not established, and likely impossible** — 438.4 MB of weights alone exceeds it |
| **CPU execution** | Viable in principle; **performance unmeasured** |
| **Partial GPU offload** (some layers on GPU, rest on CPU) | **Untested.** This is the case file size says nothing about |
| **GPU alongside active chat/image generation** | Governed by Howl's protected VRAM reserve, whatever the reranker needs |

So the correct conclusion is *"full GPU residency is not available under this envelope"*, **not** *"CPU-only is a technical requirement"*. Partial offload and reserve-aware co-residency stay open until measured. The user-facing CPU/VRAM preference remains appropriate and unchanged.

CPU cost is the open question: bge-reranker-v2-m3 is 24-layer XLM-R, and reranking 4–6 candidates means 4–6 forward passes per query. Whether that fits the ≤ 40 ms p50 budget is **unmeasured** in both directions.

## 4. Licence note

`jina-reranker-v2-base-multilingual` is **CC-BY-NC-4.0**. If Howl ever ships commercially, that candidate is excluded regardless of quality. Recorded now so it is not discovered after adoption.

## 5. Recommended evaluation order, if and when authorized

1. **bge-reranker-v2-m3 Q8_0** (635.7 MB, Apache-2.0) — bidirectional, `rank` pooling, documented llama.cpp path, no build-age caveat. Best first test of whether a cross-encoder can separate *agent* from *patient* at all.
2. **Qwen3-Reranker-0.6B Q8_0** (639.2 MB, maintainer GGUF) — only if (1) shows the capability exists but scores weakly. Carries the b10076 risk from §1 and needs a conversion sanity check.
3. Everything else only if the first two fail for reasons that are architectural rather than capability-related.

## 6. Pre-fetch verification checklist (each is a gate, not a formality)

1. Pin the revision; record file SHA-256 at download; re-verify before scoring.
2. Confirm the GGUF carries the reranker head (`cls.output.weight` / RANK pooling metadata). Community Qwen3 reranker GGUFs are widely reported to score ~4.5e-23 without it.
3. **Instrument sanity check before any case run — rank-based, not magnitude-based** (corrected: cross-encoder scores are not universally calibrated, so requiring ~1/~0 would reject a working model for its output scale). A candidate passes the instrument check when:
   - the relevant candidate **ranks above** the deliberately misleading one;
   - the irrelevant candidate **ranks below** the relevant one;
   - the score difference is **nonzero and reproducible** across repeated calls;
   - **reversing candidate order** does not change the ranking outcome;
   - output contains **no NaN, no infinity, and is not degenerate-constant**.

   Raw scores are recorded for the record, but no universal magnitude is imposed. A model failing this check is a **broken instrument and is rejected as such**, not scored on R1–R5.

   **Separation of powers:** the frozen selection rule (keep ≥ 0.30, else abstain) applies to *harness* scale. A candidate's score transform is reported alongside its results and **must not silently replace the frozen rule**; if a transform is needed to score the candidate at all, that is a documented change of instrument, not of evaluation.
4. Then, and only then: the frozen 9 cases, unchanged, plus the `measure_embedder.py` protocol for CPU/RAM/cold-load/warm-latency/unload.

## 7. What this document does not claim

No candidate has been run, scored, or downloaded. Every statement about which model will handle subject confusion better is a **hypothesis**, not evidence. Sizes and licences are measured metadata; behaviour is unknown.

## 8. Pinned artefact for the fetch decision — **not approved for download**

| Field | Value |
|---|---|
| Repo | `smarttasks/bge-reranker-v2-m3-GGUF` |
| Revision | `fc50d873dcd84a7be9760cfecd31f513d5ebc436` |
| File | `bge-reranker-v2-m3-Q8_0.gguf` |
| Bytes | **635,674,432** (635.7 MB) |
| Publisher SHA-256 | `d32a877b7266f5275172458616824ce6728502416672052c2da795da3389201e` |
| Licence | **Apache-2.0** (repo tag); base model `BAAI/bge-reranker-v2-m3`, 567,755,777 params, XLM-R |
| Alternatives in same repo | Q4_K_M 438,374,880 B / Q5_K_M 468,390,368 B, each with its own published hash |
| Publisher's claims | MANIFEST v1, 0 risk findings; README states clean XLM-RoBERTa→llama.cpp conversion and preserved ranking behaviour across quants — **claims, not verification** |

The publisher hash lets the download be verified byte-for-byte without trusting the transport. The tensor/head check (`cls.output.weight`, rank pooling metadata) still requires the file and happens **after** fetch, before any case is run.

### Resource expectations, labelled as expectations

| Quantity | Expectation | Basis |
|---|---|---|
| Disk | 635.7 MB | measured |
| RAM resident (CPU) | ~0.7–1.0 GB | weights + runtime overhead; **unmeasured** |
| Warm latency per candidate pair | tens to low hundreds of ms on 16 cores | 24-layer XLM-R, **unmeasured** |
| Per query (4–6 candidates) | 4–6× the pair cost | **unmeasured** |
| Full GPU residency | excluded by 384 MiB envelope | arithmetic |
| Partial offload | unknown | untested |

If latency or RAM proves binding, **Q4_K_M at 438.4 MB** is the fallback within the same pinned revision, at a fidelity cost that must then be reported.

## 9. Dispositions

| Item | Decision |
|---|---|
| Candidate enumeration | accepted as research |
| bge-reranker-v2-m3 | **first candidate; not approved for download** |
| Qwen3-Reranker-0.6B | reserve candidate |
| 384 MiB envelope | preserved; GPU interpretation clarified in §3 |
| Score sanity check | refined to rank-based, §6.3 |
| Frozen R1–R5 apparatus | untouched |
| Gate B | 6/7, open |

**The reranker is optional infrastructure, never a memory-system dependency.** Formation, lexical retrieval, semantic retrieval when available, and explicit abstention must all survive its absence — and its rejection.

## 10. Claim ceiling

Gate B is 6/7 and open. The horse/trainer failure is unresolved. This review changes what we *could* test next, not what we know.
