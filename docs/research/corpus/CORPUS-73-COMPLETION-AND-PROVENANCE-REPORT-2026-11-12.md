# Completion and Provenance Report — `CORPUS-68` Remediation Package

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Identifier:** `CORPUS-73` · **Version:** 1 · **Date:** 2026-11-12.
**Authorizing work order:** "`CORPUS-68` Remediation and Reconciliation", §§6–8.
**Assignment model identity (owner-supplied):** DeepSeek. Persona label Buffy = execution fact, not model identity.
**Execution:** fresh conversation; read-only against every source; writes limited to the five artifacts in §1 and their registry rows.
**Status:** remediation package **complete for independent review**. Nothing adopted, applied, or implemented.

---

## 1. Artifacts created

| # | Artifact | Identifier | Bytes | Lines | SHA-256 | Registry row |
|---|---|---|---|---:|---:|---|---:|
| 1 | `CORPUS-69-OWNER-DECISION-RECORD-D1-D7-2026-11-12.md` | `CORPUS-69` | 8,940 | 125 | `52dba0dd0d50c1ecb53427f920864756bd2bde980012bbe44f562758e49a750a` | 404 |
| 2 | `CORPUS-70-CORPUS-68-FINDING-RECONCILIATION-AND-REGISTRY-ACCOUNTING-2026-11-12.md` | `CORPUS-70` | 28,236 | 327 | `7c79bba299ffd232f201592a55270d0faf64dfd0aec03d0e287772ae135094fc` | 405 |
| 3 | `CORPUS-71-PROPOSAL-REVISION-REGISTER-V1-1-2026-11-12.md` | `CORPUS-71` | 16,772 | 198 | `c00ff86d99d3bfc9aea74f5322431ab40f809a091ebe0e4f3d30dfd7022023ac` | 406 |
| 4 | `CORPUS-72-OWNER-DECISION-BRIEF-OD1-OD8-2026-11-12.md` | `CORPUS-72` | 12,179 | 159 | `f24b010a04f76621e095478707f3e2d90e0e0f802ba9f76ca67283f09cfa9d81` | 407 |
| 5 | `CORPUS-73-COMPLETION-AND-PROVENANCE-REPORT-2026-11-12.md` | `CORPUS-73` | see row 408 | see row 408 | see registry row 408 | 408 |

All five are in `docs/research/corpus/`, matching the existing `CORPUS-*` convention. **This document cannot carry its own hash or size** — a file cannot hash itself — so its row in `docs/DOCUMENT-REGISTRY.md` is authoritative for its identity.

**No prior artifact was edited.** `CORPUS-62`–`68` are byte-identical to their pre-filing state; their hashes are unchanged and verified in §2.3.

---

## 2. Validation actually performed

### 2.1 Registry registration sweep

Procedure: the established registration pass — insert the five rows, renumber `1…N` preserving order, recompute every row's SHA-256, byte count, and LF line count from the working tree, update the header `**Rows:**` count, then verify.

| Check | Result |
|---|---|
| Data rows after filing | **408** |
| Rows under `docs/` / root-level rows | **406 / 2** |
| Duplicate row numbers | **0** |
| Duplicate targets (one file, two rows) | **0** |
| Rows whose linked target does not resolve | **0** |
| Hash / byte / line mismatches among all rows | **0** |
| Header `**Rows:**` line | `408 — 406 documentation files under docs/ … plus 2 root-level files` |

Row 393 (`CORPUS-58` brief) now resolves to `../CORPUS-58-…`, hash `292befdd…`, 4,311 bytes, 80 lines — the state `CORPUS-63` established, unchanged by this filing.

### 2.2 Existing mechanical checker

Command: `python docs/research/corpus/check_corpus.py`

```text
scored units: 45
split: tuning=19 heldout=26
declared in CORPUS-01: {'units': 45, 'tuning': 19, 'heldout': 26}
TOTALS MATCH
RESULT: PASS (0 fail, 0 warn, 0 unverifiable)
```

Exit status **0**. **This is a pass on item records only.** The checker never opens `docs/DOCUMENT-REGISTRY.md` — verified in source by `CORPUS-67` §2 and re-verified here — so this result says **nothing** about registry integrity, and it did not detect any of the eight accounting defects in `CORPUS-70` §8.2. `0 warn` is not a clean bill: `C-2` is warn-only by design and semantic `Required`/`Variation` conflicts are declared undecidable, which is why the `F3-01` conflict and the `F6-01` evidence gap were never mechanically detectable.

### 2.3 Governing sources unchanged

| Artifact | SHA-256 (prefix) | Unchanged |
|---|---|---|
| `CORPUS-01` | `be1f11b4…` | **Yes** |
| `CORPUS-02` | `75aa9be7…` | **Yes** |
| `CORPUS-03` | `d0dae1c0…` | **Yes** |
| `CORPUS-04` | unchanged | **Yes** |
| `CORPUS-07` | `72f6e2a2…` | **Yes** |
| `MEMORY-09` | `f77e47a9…` | **Yes** |
| `check_corpus.py` | `be3e317c…` | **Yes** |
| Root `CORPUS-58` brief | `292befdd…` | **Yes** |
| `CORPUS-62`–`67` | unchanged | **Yes** (not edited by this package) |

### 2.4 What the validation does **not** establish

- It does not establish that any proposal is correct, adoptable, or diagnostically sufficient.
- It does not establish that the registry's **summary tables** are correct — they are not (`CORPUS-70` §8.2), and this filing deliberately left them unrepaired.
- It does not establish reviewer independence, model identity, or convergence.
- It did not execute any validator fixture: `CORPUS-67` is a specification and no implementation exists.

---

## 3. Dispositions at a glance

| Area | Outcome |
|---|---|
| `CORPUS-68` findings S-1 – S-7 | **7 reconciled:** 4 upheld, 2 partially upheld, 1 rejected in part (S-3's `F1-03` contradiction) — full detail `CORPUS-70` §3 |
| Additional defects found by this reconciliation | **4:** `CORPUS-64` row 12 vs §10; `CORPUS-66` §6 self-contradiction; `CORPUS-65` stripped backticks; registry §4.1 double error |
| Proposals revised | **4** to v1.1 (`CORPUS-71` §§2–5); v1 texts preserved unedited |
| New proposals | **2:** `A-R2-02` (`CORPUS-71` §6), `A-REG-01` (`CORPUS-70` §8.3) |
| New owner decisions identified | **8** (`CORPUS-72`) |
| Prior owner decisions | **D1–D7 filed with conversation provenance** (`CORPUS-69`); **not reopened** |
| Diagnostic count | **1 of 24 confirmed** (`F2-02`) — **unchanged**; `F3-01` **not** promoted |

## 4. Limitations, stated rather than resolved

1. **No model inference.** No configuration was run, no rate measured, no experiment performed. Every diagnosticity statement is an analytic pass-set result on text.
2. **No rater study.** `F3-01`'s and `F6-01`'s proposals both depend on raters applying judgement rather than keywords; that agreement is unmeasured.
3. **The substantive `CORPUS-58` review remains unfiled.** D2 stands; nothing here relies on its content.
4. **Owner decisions are not independently verifiable.** No signature, timestamp, or attestation exists; `CORPUS-69` §0 says so and invents none.
5. **`CORPUS-68` is not an independent reviewer.** Same model as the `CORPUS-58` material; convergence **0 of 2**.
6. **Registry §2.1/§2.2/§4/§4.1 remain wrong.** Corrected values are stated for the proposer; the repair is unauthorised (`OD-7`).
7. **`CORPUS-26`/`CORPUS-32` remain unregistered.** Identity established in-file; registration recommended to accompany `A-REG-01`.
8. **Validator unimplemented; fixtures unexecuted.**
9. **`F3-01` cannot be certified Type P on text.** No wording change closes confabulation; only a base-rate measurement would, and it is not authorised while `G1` is open.

---

## 5. Final status

| Area | Status |
|---|---|
| Seven prior owner decisions | **Filed with conversation provenance** (`CORPUS-69`); approval asserted, not verifiable; D1 applied, D2 disposition recorded, D3–D7 carried into proposals |
| `CORPUS-68` findings | **Individually reconciled** — 4 upheld, 2 partially upheld, 1 rejected in part; nothing accepted on assertion |
| `F3-01` | Proposal at **v1.1**, **unapplied**; classification **UNRESOLVED — P-candidate, not established**; **not** Type P |
| `F6-01` | Proposal at **v1.1**, **unapplied**; adjudicated threshold **preserved** in v1 and **explicitly tightened** in v1.1 with `Δtest: Yes` routed to OD-3; item unscorable as filed |
| `R-2` | Specification conflict **verified** (3 listed / 4 counted / 4 tagged / `F1-03` omitted); resolution proposed as `A-R2-01`, plus `A-R2-02` for the "which fact" ambiguity; **unapplied** |
| `CORPUS-67` | Specification at **v1.1**; fixture-7 determinism fixed, WARN exit status **left as an owner decision (OD-8)**, totals explicitly out of phase-1 scope; **unimplemented** |
| Registry accounting | **Reconciled arithmetically and explicitly unresolved in the artifact**: 366 vs 400, §4 395 vs 403, §4.1 three counts and one sum wrong; corrections proposed as `A-REG-01`, **not applied** |
| `CORPUS-26` / `CORPUS-32` | Provenance **established in-file** (`CORPUS-26`: Muse Spark, independent reviewer; `CORPUS-32`: packet, review not performed); **both unregistered**; registration recommended with `A-REG-01`, **not performed** |
| Governing standards | **Unchanged** (`CORPUS-01`, `02`, `03`, `04`, `05`, `07`, `MEMORY-09` hash-verified) |
| Scored corpus and tags | **Unchanged** — no item, tag, split, acceptance condition, or classification altered |
| Checker implementation | **Unchanged** (`be3e317c…`); no registry logic added |
| Model inference | **Not performed** |
| `G1` | **Open** |
| Corpus | **Frozen** |
| `AM-01` | **Not issued** |
| `A-2a-R1` | **Not in force** |
| `D1` | **Not adopted** |
| Reviewer convergence | **0 of 2** — unchanged; `CORPUS-68` claims and receives no credit |

**Stop state.** The package is ready for independent review and owner adjudication. Nothing further was implemented, and `CORPUS-72` is the only document that asks the owner to decide anything.

`G1` open · corpus frozen · expansion frozen · annotation blocked.
