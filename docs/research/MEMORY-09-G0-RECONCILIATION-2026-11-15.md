# MEMORY-09 — G0 Closeout Reconciliation

**Status:** documentation-only reconciliation of the G0 receipt's three findings (D1, D2, O1).
No artifact was rewritten. No architecture rule was added. No experiment was run. No
implementation exists.

**Historical artifacts preserved unchanged** (verified before and after writing this record):

| Artifact | Size | SHA-256 |
|---|---|---|
| `docs/research/MEMORY-09-EXPERIMENT-PLAN-2026-11-15.md` (v3, approved) | 21,616 B | `6c4078bd8c1fd29f1c74691afd0e25273000ac8bbddf6cc5c7f162bcfbe6bbb9` |
| `docs/research/MEMORY-09-G0-PREFLIGHT-2026-11-15.md` | 8,445 B | `5c245283c1712ed59ade4ab51f6310b35b4ad428d6b446061c50b6f577f162d0` |
| Frozen S3 artifacts (8) | — | as recorded in plan §1.1; all MATCH |

This record **supersedes the erroneous labels** in plan §1.1 for citation purposes. It does not
edit the plan, because the plan's hash is the identity G0 authorized.

---

## 1. D1 — Temporal provenance labels

### 1.1 The erroneous reference, preserved verbatim as historical context

Plan v3, §1.1, final row:

```text
| receipts | `43d744e9…` (temporal v2), `06f8413885aeb992` (A/B), `caed2a2f…` (temporal v1) | immutable results |
```

The A/B entry is correct. Both temporal labels are wrong, and for different reasons.

### 1.2 Verified identities

| Artifact | Path | Size | SHA-256 |
|---|---|---|---|
| Temporal **spec** (v2 case set) | `eval/s3/temporal-spec.json` | 7,961 B | `43d81e07e236df8c0128ec8e2c86a99c188514a75265892fb10d5ad7049bd93a` |
| Temporal **receipt** (v2 run) | `eval/s3/temporal-receipt.json` | 5,812 B | `caed2a2f6791e44c269c3ed1bd69e864b580ab9065b76a522f23235be84d8ca2` |
| Temporal v2 **case set** (content hash) | `cases_sha256` field in both files above | — | `43d744e982a46fa7d6a6c51fc8240ab73eb0a6144789800080dec6c686419611` |

Receipt internal fields, read directly:

- `experiment` = `E-TEMPORAL-SET`
- `version` = **`v2`**
- `cases_sha256` = `43d744e982a46fa7d6a6c51fc8240ab73eb0a6144789800080dec6c686419611`
- `frozen_note` states the run used a fresh store, new owner ids and new cases, and did not
  reference or modify the frozen seven records or the R1–R5 case set.

### 1.3 Corrected mapping

| Plan label | Actually is |
|---|---|
| `` `43d744e9…` (temporal v2) `` | the **v2 case-set content hash**, not a receipt identity |
| `` `caed2a2f…` (temporal v1) `` | the **v2 run receipt** — its own `version` field says `v2`; not v1 |
| `` `06f8413885aeb992` (A/B) `` | correct: `eval/s3/ab-rerank-receipt.json` |

**Additionally:** no temporal **v1** receipt artifact exists under `eval/s3/`. The v1 run was
withdrawn for a harness defect, so `` (temporal v1) `` does not merely point at the wrong file —
it has no referent on disk at all. The only temporal receipt present is the v2 one.

### 1.4 Citation rule from here on

- A **run** is cited by `(path, sha256)` of its receipt: temporal v2 =
  `eval/s3/temporal-receipt.json` @ `caed2a2f…`.
- A **case set** is cited by the field name `cases_sha256`, never as a receipt:
  v2 cases = `43d744e9…`.
- A **withdrawn run** is cited as withdrawn, never as a live receipt.

No file was renamed or deleted; the wrong labels exist only in the plan text and in the G0
receipt, both preserved as history.

---

## 2. D2 — `Q-MF-01`

### 2.1 Where it actually lives

`eval/s3/formation.py`, line 7, module docstring:

> AI-assisted formation is a separate, optional method (`Q-MF-01` remains open).

The architecture contract contains **no** `Q-MF` identifier. Its open-question list covers
Q1–Q7, Q14a and Q14b, plus unnumbered follow-ons in §14 (group identity, ambient beliefs,
inspection exposure, schema versioning, retention, deletion). None of them asks about formation
method.

### 2.2 Classification

`Q-MF-01` is a **research question**, not a contract question:

- It asks an empirical question — *does model-assisted memory formation outperform the
  deterministic extractor, and at what cost?* — whose answer requires an experiment, not a
  design decision.
- The contract's open questions are **design decisions** awaiting owner approval. Adding a
  formation-method question there would imply the architecture has a formation rule to decide,
  which it deliberately does not: formation method is `Q-MF-01`'s subject precisely because it
  is unsettled.

### 2.3 Authoritative reference (assigned here, explicitly non-normative)

| Field | Value |
|---|---|
| Identifier | **`RQ-MF-01`** |
| Kind | research question (non-normative) |
| Statement | Does model-assisted formation produce records at least as good as the deterministic extractor on the experiment's own measures, and at what prompt/token/latency cost? |
| Original label | `Q-MF-01` (retained as an alias, as written in `formation.py`) |
| Canonical location | this record, §2.3 |
| Also referenced in | `eval/s3/formation.py:7`; plan v3 §3.1 and §5 disclosure |
| Status | **open, untested** |

**Explicit non-claims:** this record does not resolve formation, does not choose a formation
method, and does not add a normative rule to the contract. The plan's existing statement stands
unchanged: A1 measures memory *use* with deterministic formation and cannot answer `RQ-MF-01`.

**Suggested follow-up (requires owner approval, not done here):** register `RQ-MF-01` in the
contract's §14 as a research question clearly separated from the design questions, so that a
later reader does not conclude the contract simply forgot it.

---

## 3. O1 — Duplicate reranker specification

### 3.1 Identities

| Path | Size | SHA-256 |
|---|---|---|
| `eval/s3/rerank-spec.json` | 11,617 B | `bbf400c225e6c90a20a6c40da920b567e2d9137a6ede6f3399507df8bff804e2` |
| `eval/s3/rerank-spec-v1.1.json` | 11,617 B | `bbf400c225e6c90a20a6c40da920b567e2d9137a6ede6f3399507df8bff804e2` |

Byte-identical. Neither is stale, and neither may be deleted: each is loaded by different code.

### 3.2 Why both exist

| Consumer | Path loaded | Role |
|---|---|---|
| `eval/s3/rerank_eval.py` | `rerank-spec.json` (line 28, `SPEC_PATH`) | the v1.1 harness's live protocol |
| `eval/s3/rerank_eval_v12.py` | `rerank-spec-v1.1.json` (line 30, `SPEC_V11`) | the **immutability reference** the v1.2 module checks against |

So `rerank-spec-v1.1.json` is not a redundant copy: it is the v1.2 module's guard that v1.1 did
not change. Declaring one "canonical" and removing the other would break that check.

### 3.3 Citation guidance

- Cite **by purpose**, not by preference: v1.1 harness protocol →
  `eval/s3/rerank-spec.json`; v1.2's immutability reference →
  `eval/s3/rerank-spec-v1.1.json`.
- Always include the hash; since both carry `bbf400c2…`, a hash-only citation is unambiguous
  even when the path differs.
- The v1.2 protocol itself is a different artifact:
  `eval/s3/rerank-spec-v1.2.json` @ `5414d0a7…`.

**No deletion, no rename, no edit.**

---

## 4. Closeout

- D1: corrected mapping recorded, historical error preserved, no artifact rewritten.
- D2: classified as a research question with a stable non-normative identifier; formation
  remains unresolved.
- O1: duplicate explained by consumer role; both retained; citation rule stated.

**Next required step:** the A1 preregistration (`MEMORY-09-PREREGISTRATION-2026-11-15.md`),
which fixes `N` and repetitions before any run. **G1 implementation remains unauthorized.**
