# MEMORY-09 S3 — Gate B Branch Evaluation

**Identity:** `MEMORY-09-S3-GATEB-1` · **Date:** 2026-11-15 · **Result:** **PARTIAL — 6 of 7 probes. Gate B is NOT closed.**
**Reason it is not closed:** there is no semantic branch. No embedding model was authorized or downloaded, and a hashed n-gram space was **not** relabelled "semantic" to pretend otherwise.
**Receipt:** `eval/s3/gate-b-receipt.json` (7,221 B), written by `eval/s3/evaluate.py`
**Formation:** frozen and verified byte-identical to the Gate A baseline.

---

## 1. Branches, evaluated alone before fusion

Each probe was run against each branch in isolation; fusion and reranking are reported separately so a working branch cannot cover for a dead one.

| Probe | Question | Lexical alone | Subword alone | Fused | Reranked | Verdict |
|---|---|---|---|---|---|---|
| **P1** | "Who was your trainer, again?" | **none** | `trained_by` | `trained_by` | `trained_by` | **PASS** |
| **P2** | "What did you think of the weather?" | none | none | none | none | **PASS** (abstains) |
| **P3a** | "What is the vault code?" (no owner filter) | Alaric's record | Alaric's record | Alaric's record | Alaric's record | **PASS** (pre-filter contamination measured) |
| **P3b** | "What is the vault code?" (Sable) | none | none | none | none | **PASS** (filter removes it) |
| **P4a** | "Where is the ledger kept now?" | green | green | green | green | **PASS** |
| **P4b** | "Where was the ledger kept before?" | green, blue | none | green | **blue** | **PASS** (history reachable when asked) |
| **P5** | "Who trained the horse?" | `trained_by` | `trained_by` | `trained_by` | `trained_by` | **FAIL** (should abstain) |

**The headline:** P1 — the paraphrase that failed in Gate A — now succeeds, through the subword branch with lexical returning nothing. That is a genuine capability gain, and it is attributable, because the branches are scored separately.

## 2. Three defects found by evaluating branches alone

All three were invisible under a fused score. This is the argument for the requirement, not an aside.

| # | Symptom | Cause | Fix |
|---|---|---|---|
| 1 | Subword returned **zero hits everywhere** — a dead branch | symmetric cosine penalised short queries against long records | IDF-weighted **query coverage** instead of symmetric cosine |
| 2 | Abstention collapsed: a weather question matched "map that was drawn" | function-word n-grams (`the`, `of`, `a`) dominated | content-word n-grams only |
| 3 | "vault code" matched the trainer record at score 1.0 | out-of-vocabulary grams dropped from the denominator, so one accidental gram scored perfectly on a tiny corpus | unseen grams score at **maximum IDF** and penalise coverage |

A fused-only evaluation would have shown P1 fixed and looked like success. Defect 1 was a branch contributing nothing; defect 3 was a branch contributing noise. Both needed branch-level evidence to see.

## 3. P5 — the failure that matters

`"Who trained the horse?"` retrieves `trained_by` on **both** branches. It shares the token *trained*; it asks about a different subject entirely. No lexical or subword method can separate these: the distinction is not in the characters, it is in what the sentence is about.

This is the honest boundary of the implemented system, and it is the concrete requirement for the missing branch:

- a real semantic representation can separate *trainer-of-a-person* from *trainer-of-a-horse*;
- a reranker with cross-encoding could too;
- neither exists here, and a heuristic term-coverage reranker demonstrably cannot fix it.

**Therefore:** Gate B remains open on two counts — semantic absent, and P5 unmet. Reporting 6/7 as a pass would be exactly the "success because the example worked" failure the acceptance criteria were written to prevent.

## 4. Formation froze, and was proven to have frozen

| Module | Gate A hash | Now | |
|---|---|---|---|
| `formation.py` | `d8aac5003c66229d…` | `d8aac5003c66229d…` | identical |
| `records.py` | `1432fb63da5ae561…` | `1432fb63da5ae561…` | identical |
| `evidence.py` | `19db697a0f82d29f…` | `19db697a0f82d29f…` | identical |

The evaluator refuses to run if they differ. Records: 6 active, 1 superseded — same as Gate A.

## 5. Reproducibility record

| Component | Configuration |
|---|---|
| Lexical | BM25, k1 = 1.5, b = 0.75 |
| Subword | hashed char n-grams, n ∈ {3,4,5}, content words only, dim 2048, seed `s3-subword-v1`, metric = IDF-weighted query coverage |
| **Semantic** | **absent** — `available: false`, reason recorded in the receipt |
| Fusion | reciprocal rank fusion, k = 60, floor 0.012 |
| Rerank | heuristic: fused + 0.10 × term coverage + kind prior − 0.02 stale penalty. **Not a cross-encoder.** |
| Embedder identity | `null` |
| Determinism | total; no model, no server, no network |

## 6. What would close Gate B

1. An **embedder decision** (owner): which local model, and authorization to obtain it. `DA-R2` stays open either way.
2. A semantic branch that is real, reports degradation rather than falling back silently, and contributes visibly to `branches_used`.
3. P5 abstaining on subject mismatch — the test that separates semantics from string overlap.
4. No regression on A01–A14 and on P1–P4b.
5. A Gate B receipt whose `result` is PASS with `semantic.available = true`.

Until then the standing description is: **a two-branch system with a measured paraphrase win and a measured subject-confusion failure, with the third branch missing.**

## 7. Claim ceiling

Six of seven probes pass with formation frozen. No model was loaded, no embedding model exists in the system, nothing was downloaded, and no claim is made about how any language model would behave with these memories. The subword branch is **not** semantics and is never described as such.
