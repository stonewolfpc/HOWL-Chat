# MEMORY-09 — Preregistration Freeze Record

**Status: FROZEN.** Recorded 2026-11-15 with the owner's approval. This record is the freeze
artifact; the four frozen documents are **not edited** by it, so their hashes remain the
identities that were approved.

**What freezing means here:** the experimental design has fixed comparisons, fixed error control,
fixed magnitude criteria and stated limitations. No further negotiation of success criteria is
permitted after this point. Corrections are append-only amendments registered in §4.

---

## 1. Frozen preregistration set

| # | Document | Bytes | SHA-256 (full) |
|---|---|---|---|
| 1 | `docs/research/MEMORY-09-PREREGISTRATION-2026-11-15.md` | 9,235 | `8c788a96ea6a866d4fd06e13496ddbfa46b649da5a5e0c8ef9bace960844ff12` |
| 2 | `docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-1-2026-11-15.md` | 9,800 | `975033401ae1beaa9caaa5d5ae55aab0ca3ddb4f5159eb41e9feb8e01a4c9b1a` |
| 3 | `docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-2-2026-11-15.md` | 8,717 | `c3ffd7ed735ed4cf68f8b455555af2a1d10eb38273b07a756f269076cd917f8f` |
| 4 | `docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-3-2026-11-15.md` | 6,803 | `c5af036265b03662b26feaf40d358ad08c3063fa0055c72d2fb0c2f41c236281` |

Machine-checkable manifest (extract and run `sha256sum -c` from the repository root):

```text
8c788a96ea6a866d4fd06e13496ddbfa46b649da5a5e0c8ef9bace960844ff12  docs/research/MEMORY-09-PREREGISTRATION-2026-11-15.md
975033401ae1beaa9caaa5d5ae55aab0ca3ddb4f5159eb41e9feb8e01a4c9b1a  docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-1-2026-11-15.md
c3ffd7ed735ed4cf68f8b455555af2a1d10eb38273b07a756f269076cd917f8f  docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-2-2026-11-15.md
c5af036265b03662b26feaf40d358ad08c3063fa0055c72d2fb0c2f41c236281  docs/research/MEMORY-09-PREREGISTRATION-AMENDMENT-3-2026-11-15.md
```

---

## 2. What is frozen, in substance

| Element | Frozen value | Where |
|---|---|---|
| Independent conversations | `N` = 16, paired | Prereg §2 |
| Repetitions | 3 matched generation seeds per conversation × arm | Prereg §4, Amend. 1 §5 |
| Probes | 13 per conversation, quota fixed | Prereg §5 |
| Coverage | ≥3 rollovers; every probe's span absent from the window | Prereg §5 |
| Common outcome | **M1** answer quality, identical scoring in every arm | Amend. 1 §2 |
| Memory-specific outcome | **M2**, A0 = `N/A`, D = manipulation check | Amend. 1 §3 |
| Claims | **C1** A1 > A0 on M1; **C2** A1 > D on M1; full claim = C1 ∧ C2 | Amend. 2 §1 |
| Error budget | family-wise α = .05 → **α = .025 per claim** | Amend. 2 §3.1 |
| Test | exact two-sided sign test on `n_eff`; thresholds by `n_eff` | Amend. 2 §3.2 |
| Ties | excluded from the sign count, **included** in the median | Amend. 2 §2, Amend. 3 §2 |
| Interpretability floor | `n_eff < 12` → indeterminate | Amend. 2 §2 |
| Magnitude | median `d_i` ≥ 0.10, **operational** guard only | Amend. 2 §3.3, Amend. 3 §1 |
| Attainable threshold | smallest qualifying median = 4/39 ≈ 0.102564 | Amend. 3 §1.2 |
| Invalid pairs | outcome-blind, ≤2 of 16, no replacement; veto fails the run | Amend. 2 §4 |
| Stopping | fixed sample, outcome-blind halts only, no post-hoc extension | Prereg §7 |

**Carried forward into results (non-negotiable wording):** 0.10 is a preregistered guard against
trivial differences, **not** an established practical-significance threshold. No result may be
described as "practically meaningful" on the strength of it.

---

## 3. What is preserved unchanged

| Document | SHA-256 prefix | Role |
|---|---|---|
| `MEMORY-09-EXPERIMENT-PLAN-2026-11-15.md` | `6c4078bd…` | approved plan v3 |
| `MEMORY-09-G0-PREFLIGHT-2026-11-15.md` | `5c245283…` | G0 receipt |
| `MEMORY-09-G0-RECONCILIATION-2026-11-15.md` | `2b9c55b6…` | D1/D2/O1 record |
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` | `e88ef2d1…` | design contract r32 |
| `MEMORY-09-CONTRACT-NORMATIVE-2026-11-15.md` | `87154a9b…` | generated extract |
| frozen S3 artifacts (8) | per plan §1.1 | baseline apparatus |
| all evaluation receipts | per G0 receipt | immutable results |

---

## 4. Amendment register

| # | Document | Status | Hash |
|---|---|---|---|
| 1 | Metric symmetry | accepted | `97503340…` |
| 2 | Capability acceptance rule | accepted | `c3ffd7ed…` |
| 3 | MMD interpretation and calculation | accepted | `c5af0362…` |

Any later amendment is appended with the next number, dated, and **must itself be registered in
the G4 receipt**. An unregistered amendment does not bind the analysis.

---

## 5. G4 preconditions (binding on the future run)

Before the **first model call**, the G4 run must:

1. recompute the four hashes in §1 and abort on any mismatch;
2. write all four full hashes into the G4 receipt;
3. record every registered amendment's hash;
4. record the seed table's hash (§1 of Amend. 1 §5 derivation) before the first call.

A run that starts without these recorded is **invalid** and cannot support C1 or C2.

---

## 6. Authorization state

| Item | State |
|---|---|
| Preregistration freeze | **authorized, completed by this record** |
| G1 implementation (`eval/mem1`) | **not authorized** |
| G4 execution / subject runs / model calls | **not authorized** |
| Frozen baseline artifacts, Howl Chat | **untouched** |
| Commits or pushes | **none** |

`eval/mem1` does not exist. No model call has been made under this design.

The design is fixed. What remains is to build the instrument and run it — separately authorized.
