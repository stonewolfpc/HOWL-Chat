# MEMORY-09 S3 — Cross-Encoder Reranker: Scope Only

**Identity:** `MEMORY-09-S3-RERANKER-SCOPE-1` · **Date:** 2026-11-15
**Status: ACCEPTED as a research record (project owner, 2026-11-15).** Scoped only — **no model authorized, nothing fetched.**
The R1–R5 evaluation harness now exists: `eval/s3/rerank_eval.py` over `eval/s3/rerank_set.py`.

This document defines the job. It does not approve a download, a size, a vendor, or an execution path.

---

## 1. The job, stated precisely

P5 fails because a bi-encoder compares *topic*, not *arguments*. "Who trained you?" and "Who trained the horse?" occupy nearly the same vector neighbourhood, and the trainer record wins both. The reranker's job is to compare a query and a candidate **jointly**, so the relation and its arguments are visible at scoring time.

It is **not** tasked with: recovering facts that formation never recorded; substituting for lexical or semantic branches; or rescuing probes by threshold adjustment.

## 2. Acceptance criteria, fixed before any download

| # | Criterion | Falsified by |
|---|---|---|
| R1 | **Agent vs patient** — distinguishes who performed an action from who received it | P5 still returns the trainer for "who trained the horse?" |
| R2 | **Subject ownership** — a character-specific fact scores low for another character | P3b gains Alaric's answer when Sable asks |
| R3 | **Historical vs current** — superseded records do not outrank the current one | P4a returns the blue cabinet |
| R4 | **No P1 regression** — paraphrase matching survives reranking | P1 drops below the current passing state |
| R5 | **Abstention preserved** — irrelevant queries still select nothing | P2 or P3b stop abstaining |

R4 and R5 are **vetoes**: a reranker that gains P5 by losing P1 or abstention is rejected, not averaged.

## 3. Architecture constraints (non-negotiable)

- **Optional.** Absent reranker ⇒ current 6/7 behaviour, reported as degraded, never silently substituted.
- **No permanent residency.** Loaded per query batch or per block, then released; residual VRAM must return to baseline within budget.
- **No dependency.** Formation, lexical retrieval, ownership filtering and abstention must not import it.
- **Separate envelope.** Same budgets as the embedder (≤384 MiB during, ≤64 MiB residual, p50 ≤40 ms warm) measured with the same harness before adoption.
- **One model at a time.** If the reranker and the embedder cannot co-reside under budget, the reranker runs on demand and reports the swap.

## 4. Evaluation plan — now frozen and mechanised

| Requirement | Where it is enforced |
|---|---|
| Cases, expected outcomes, selection rule, threshold, veto logic frozen and hashed | `eval/s3/rerank-spec.json` v1.1, `cases_sha256 = ca6dd5b3541f08fb…`, threshold 0.30 fixed **before** any candidate existed |
| Harness refuses to run on a changed case set or changed scoring code | `load_spec()` exits 2 with the frozen/live hashes (demonstrated) |
| Candidate scoring isolated from mock assumptions | scorer receives query + candidate fields only; label-leak assertion; **tie-break on record_id, not input order** (order-invariance check found this defect) |
| Per-case evidence, not just totals | selected id, top score, full score map, abstended flag, failure reason |
| Existing system as comparison baseline | `baseline_system_as_is` maps each case to its frozen Gate B probe outcome (P5 FAIL is the reference point) |
| No threshold tuning on the evaluation set | threshold is read from the frozen spec and never recomputed |
| Resource measurement for a real candidate | `eval/s3/measure_embedder.py` protocol: CPU, RAM, cold load, warm p50/p95, unload residual, VRAM if GPU |

A candidate must be scored against **these exact frozen cases**, with its own resource numbers, or it is not evaluated.

## 5. What would make me recommend abandoning it

If R1 cannot be met without failing R4 or R5, the correct conclusion is that a cross-encoder is the wrong tool for this failure and the limitation gets **recorded**, not engineered around with more models.

## 6. Claim ceiling

No reranker exists here. Nothing in this document is evidence about any specific model's behaviour.
