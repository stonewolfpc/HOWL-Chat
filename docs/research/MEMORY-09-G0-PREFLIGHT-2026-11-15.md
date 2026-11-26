# MEMORY-09 — G0 Preflight Receipt

**Gate:** G0 (plan approval / preflight) of the v3 experiment plan.
**Scope executed:** verification only. No `eval/mem1/` created, no experimental subject run, no
evaluation data collected, no frozen artifact modified, no Howl Chat change, no commit.

**Disposition: PASS**, with **two recorded deviations** and **one ambiguity deferred to G4**
(all detailed in §5). A G0 pass authorizes nothing beyond scope: implementation still requires
separate authorization.

---

## 1. Authorized plan identity

| Field | Cited | Actual | Result |
|---|---|---|---|
| Path | `docs/research/MEMORY-09-EXPERIMENT-PLAN-2026-11-15.md` | present | PASS |
| Size | 21,616 B | 21,616 B | PASS |
| SHA-256 | `6c4078bd8c1fd29f1c74691afd0e25273000ac8bbddf6cc5c7f162bcfbe6bbb9` | identical | PASS |

The plan was verified **before** any other check, so every subsequent result is tied to this
exact artifact.

## 2. Frozen prerequisites (§1.1 of the plan)

Every artifact the plan declares frozen was re-hashed from disk.

| Artifact | Cited | Actual | Result |
|---|---|---|---|
| `eval/s3/rerank-spec.json` | `bbf400c225e6c90a` | `bbf400c225e6c90a…` | MATCH |
| `eval/s3/rerank-spec-v1.2.json` | `5414d0a7` | `5414d0a7ced42386…` | MATCH |
| `eval/s3/rerank_eval.py` | `f59959f30ec336ae` | `f59959f30ec336ae…` | MATCH |
| `eval/s3/rerank_set.py` | `23badd5e7825718d` | `23badd5e7825718d…` | MATCH |
| `eval/s3/formation.py` | `d8aac500` | `d8aac5003c66229d…` | MATCH |
| `eval/s3/records.py` | `1432fb63` | `1432fb63da5ae561…` | MATCH |
| `eval/s3/evidence.py` | `19db697a` | `19db697a0f82d29f…` | MATCH |
| `eval/s3/temporal-spec.json` | `43d81e07e236df8c` | `43d81e07e236df8c…` | MATCH |

**8/8 MATCH.** No frozen artifact has moved since the contract checkpoint.

## 3. Baseline and provenance

Suites run read-only, unchanged apparatus:

| Check | Result | Consistency with plan §1.2 |
|---|---|---|
| `eval/s3/selfcheck.py` | PASS (14 checks, 14 passed) | agrees |
| `eval/s3/evaluate.py` | PARTIAL (6/7 probes, `formation_frozen=True`) | agrees — P5 still open |
| `eval/runs/verify-s1-batch.py` | PASS (5 calls, 0 problems) | agrees |
| `python -m eval.harness.validate` | 89/89 passed | agrees |
| `corpus/check_corpus.py` | PASS (0 fail, 0 warn) | agrees |

Additional provenance confirmation:

- **P5 remains open.** The Gate B receipt contains the probe query `Who trained the horse?`, and
  `evaluate.py` still reports 6/7 PARTIAL. The plan's "not a completed baseline" statement is
  accurate today.
- **BGE remains closed.** The BGE record states *"Recommendation: do not ship bge"* and *"§14
  Final disposition — not shipped"*, including "do not fetch another ranking candidate".
  The plan's §10 non-goal matches this.
- **Resource envelope figures agree.** `reranker-resources.json` declares budgets
  `vram_mib 384`, `p50_ms 40`, `p95_ms 120`, `rss_mb 1024` with configs
  `cpu / partial_6 / partial_12 / full_gpu` — consistent with the plan's "no configuration meets
  both" statement.

## 4. Internal consistency of the plan

| Property | Result | Evidence |
|---|---|---|
| Section numbering | PASS | 19 subsections, 1.1 → 5.3, strict numeric order |
| Cross-references resolve | PASS | `§2.4`, `§2.5`, `§2.6`, `§3.2` are plan sections; `§10.1`, `§11.3` resolve in the design record |
| Arms closed set | PASS | A0, A1, A2, C, D — all declared in §2.1, none used elsewhere |
| Outcome ladder | PASS | four levels present; headline is *memory-supported answer* |
| Token accounting | PASS | `tokens.system/history/memory/overhead/total/reserve` all defined in §2.4 |
| Placebo fidelity | PASS | 8/8 required rows in §2.5 |
| Veto set | PASS | cross-character leakage, erased history, fabricated provenance, temporal correctness |
| Placeholders | PASS | no `TODO`/`TBD`/`FIXME` anywhere |
| Retired terminology | PASS | "M-correct" appears **only** inside the v2 revision note, quoting the term it replaced |

## 5. Deviations, ambiguities, observations

### D1 — Receipt labels in §1.1 are misattributed (documentation defect, non-blocking for G0)

The plan's receipt row reads:

> `| receipts | `43d744e9…` (temporal v2), `06f8413885aeb992` (A/B), `caed2a2f…` (temporal v1) |`

Verified facts:

- `43d744e9…` is the **v2 case-set hash** (`cases_sha256`), not a receipt identity. It appears
  inside `temporal-spec.json` and `temporal-receipt.json`.
- `caed2a2f…` is the **v2 receipt file** (`eval/s3/temporal-receipt.json`), whose own `version`
  field reads `v2`. It is not a v1 receipt.
- The A/B entry is correct: `06f8413885aeb992` matches `ab-rerank-receipt.json`.

**Impact:** provenance itself is intact — no file changed. What is wrong is the *label*, which
would mislead a reader mapping hashes to artifacts. It does not affect what the experiment will
measure.

**Required action:** a plan erratum before any G4 result is reported, correcting the row to name
the receipt file and its hash, and marking the case-set hash as such. Not corrected during G0,
because editing the plan would change the hash this gate authorized.

### D2 — `Q-MF-01` is not registered in the contract (identifier ambiguity, non-blocking)

The plan says `Q-MF-01` remains open. That identifier exists in `eval/s3/formation.py`'s
docstring ("AI-assisted formation is a separate, optional method") but **not** in the design
record's open-question list, which numbers Q1–Q15 and has no formation-method question.

**Impact:** the plan's disclosure is substantively right — A1 cannot answer a model-assisted
formation question — but the reference points at code, not at the contract. Either the contract
should register the question, or the plan should cite the code location.

**Required action:** none for G0; reconcile before the plan is treated as a frozen experiment
specification.

### O1 — Duplicate protocol file (observation, benign)

`eval/s3/rerank-spec.json` and `eval/s3/rerank-spec-v1.1.json` are byte-identical
(`bbf400c2…`). Two paths, one content. Harmless, but any future citation should name one path
and note the identity.

### A1 — `N` and repetition count are still undefined (deferred to G4)

§4.4 specifies conversation *shape* (≥3 rollovers, planted fact classes, correction after
rollover, two characters, an unanswerable probe) but the plan does not fix the number of
conversations or repetitions. G4 cannot be scored without them, and choosing them after seeing
results would be exactly the tuning the plan forbids.

**Required action:** decide `N` and reps **before** the G4 run, and record the decision in the
run receipt with its rationale. This is a genuine gap in the plan, not a defect in this gate.

## 6. Boundaries observed

| Boundary | Status |
|---|---|
| G0 only | PASS — verification and this receipt only |
| Do not create `eval/mem1/` | PASS — absent |
| Do not modify frozen artifacts or pipeline components | PASS — 8/8 hashes unchanged after all checks |
| Do not run experimental subjects or collect results | PASS — no subject exists, no data collected |
| Do not modify Howl Chat | PASS — no changes under `internal/`, `app-*`, `frontend/` |
| Do not revisit BGE or add a reranker | PASS — read-only references only |
| Do not commit or push | PASS — nothing staged or committed |

Three pre-existing tracked modifications remain untouched and unrelated:
`docs/research/OWNER-ENVIRONMENT-HISTORY-2026-11-14.md`, `eval/harness/validate.py`, and
`eval/harness/__pycache__/validate.cpython-312.pyc`.

Ports 8091/8092 closed at the end of the gate.

## 7. Disposition

**G0: PASS.**

- Plan identity verified exactly: 21,616 B, `6c4078bd…`.
- All eight frozen prerequisites match; baseline unchanged; P5 still open; BGE still closed.
- The plan is internally consistent: numbering, cross-references, arm set, metric ladder,
  token accounting, placebo fidelity and veto set all check out.
- One documentation defect (D1) and one identifier ambiguity (D2) recorded above; neither
  changes what the experiment measures, and both have explicit required actions.
- One real gap (A1: `N` and reps) is deferred with a named deadline, because deciding it after
  results would be tuning.

**Implementation remains unauthorized.** This receipt authorizes nothing past G0: the subject
build at `eval/mem1/` requires separate authorization, and D1 should be corrected before any
result is reported.
