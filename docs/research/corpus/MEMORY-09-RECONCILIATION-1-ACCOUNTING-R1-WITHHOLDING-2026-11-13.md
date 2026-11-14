# `MEMORY-09` — Reconciliation 1: Accounting, `R-1` Exemplars, Withholding Evidence

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-RCN-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **NON-OPERATIVE RECORD — findings only.**
**Authorization:** owner-authorized, 2026-11-13 (`DeepSeek Work Order — File and Register MEMORY-09 Three-Question Reconciliation`).
**Nature:** documentation of a read-only reconciliation. It adopts nothing, amends nothing, and repairs nothing.

---

## 0. Source protection

No source or governing artifact was modified in producing this record. Specifically **unchanged**: all corpus item files, every `Restraint:` tag and membership, `CORPUS-01` and its coverage statements, `MEMORY-09`, `OD-M09-03`, the owner decision record and its addendum, `AM-01`, `OD-M09-05`, `CORPUS-78`, scoring rules, thresholds, splits, `check_corpus.py`, and all historical review artifacts. Defects found in those sources are **documented here and left in place**, as §6 records.

**Nothing in this document is an owner decision.** Recommendations are labelled as recommendations, findings as findings, and calculations as calculations.

---

## 1. Method and sources inspected

Read-only inspection, 2026-11-13. Counts were recomputed from source text rather than carried over from any prior summary, and every quoted field was read from the item file.

| Source | Digest prefix | Role |
|---|---|---|
| `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` | `75aa9be7…` | Items incl. `F6-02`, `P-01`, `P-03`, `F9-01` |
| `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` | `d0dae1c0…` | Batch-2 items |
| `docs/research/corpus/CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` | `3d34f053…` | `R6-01`…`R6-04` incl. container `R6-03` |
| `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | `a2bfad11…` | Mode table, `A-1`…`A-3`, coverage line, §7 status |
| `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` | `f77e47a9…` | Threshold and provenance doctrine |
| `docs/research/corpus/OD-M09-05-MEMBERSHIP-DIVERGENCE-RECORD-2026-11-13.md` | `3180c0a1…` | Prior membership divergence findings |
| `docs/research/corpus/OD-M09-03-LAYERED-CUE-STANDARD-2026-11-13.md` | `f3178e1e…` | Layers, `S-1`…`S-6`, three-level distinction |
| `docs/research/MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md` | `313a0a82…` | Owner dispositions 15–18 |
| `inference_capture.go` | 4,064 B | Prompt-capture build tag and marker audit |

Heading detection: `^#{2,3} ` + backticked ID. Scored-body detection: presence of `**Split:**`, `**Required:**`, `**Forbidden:**`, `**Restraint:**`.

---

## 2. Q1 — The 45 / 46 / 42 accounting, resolved

### 2.1 Three conventions, each internally correct

| Convention | Count | Definition (as observed in source) | Source |
|---|---:|---|---|
| **Headings** | **46** | Every level-2/3 item heading across the three item files (22 + 16 + 8) | Item files |
| **Scored units** | **45** | Headings carrying a complete scored body — `Split`, `Required`, `Forbidden`, `Restraint` | `CORPUS-01` §7 ("scored units: **45**") |
| **Family-collapsed** | **42** | Variants consolidated into their parent family | `CORPUS-01` §5 coverage line (`2+4+4+4+10+6+12`) and §7 ("Item headings: **42**") |

### 2.2 The arithmetic, and its single cause

**`R6-03` is a container, not a scored unit.** Its heading appears at `CORPUS-04` L53 — "The corroboration illusion (four provenance variants)". It carries `Family: F-7`, `Split: heldout (all four)`, `Abstention: no`, `Restraint: R-6` (L55), plus an uncertainty-source field and framing prose. It carries **no `Required` and no `Forbidden`** — the four variants at L60, L75, L89, L104 do.

- **46 headings − 1 container (`R6-03`) = 45 scored units.** ✓
- **46 headings − 4 variants consolidated into `R6-03` = 42 family units.** ✓
- **46 `Restraint:` tags = 45 scored units + 1 container tag.** The container is tagged (L55), which is why tag totals exceed scored units by exactly one.

Both `45` and `42` are correct under their own definitions; neither is a correction of the other. The defect is that **no document states which convention it is using**, so the two read as competing claims about one quantity.

### 2.3 The six `R-6` counts, kept apart

| Count | Value | Definition |
|---|---:|---|
| Table membership (`CORPUS-01` §5) | **2** | `F7-01`, `F9-04` only — predates the `R6` series |
| Coverage line | **6** | 4 mechanisms + 2 legacy, variants collapsed |
| §7 status | **9** | "`R-6` resolved (9 items, 4 mechanisms + 4 provenance variants)" — its own arithmetic gives **8** |
| Tag count | **10** | 8 `R6` headings + `F7-01` + `F9-04` |
| Dedicated fixtures | **8** | `R6-01`, `R6-02`, `R6-03`, `R6-03a`…`d`, `R6-04` |
| **Scored `R6` units** | **7** | 8 fixtures − the `R6-03` container |

The last line is new to this record and is the count needed for sample-size reasoning: the container is a fixture heading but not a scorable unit.

---

## 3. Q2 — `R-1` currently has no clean exemplar

`R-1` is defined as *suppress mention — knows the fact, chooses not to name it; fact unreferenced; nothing implies it was consulted*, with silence passing (`CORPUS-01` L121).

| Item | Tag | Observed fields | Finding |
|---|---|---|---|
| **`F6-02`** (L145) | `R-1` | L146 `Abstention: no`; L149 `Required: the aversion`; L151 `Mandatory: the song is not pushed`; over-exposure: explaining why the song upsets them; forbidden: the "bad year" details, "never given" | **Closest candidate, but a hybrid.** `Required: the aversion` asks for something perceptible inside an `R-1` tag. Under `R-1`, silence passes; if the aversion must be *displayed*, silence fails. That is an `R-2`-shaped requirement (`R-2` silence column: "**No** — a flat, unchanged reply fails") |
| **`P-01`** (L231) | `R-1` | L232 `Abstention: no`; L235 `Required: the loss, at the level the moment allows`; L237 `Mandatory: … acknowledged — or met with a silence that registers it … using at most one of the facts`; over-exposure: more than one listed fact | **Proportion, not suppression.** The loss is acknowledged; what is bounded is *how much* is disclosed. Failure mode is over-telling, not leaking one specific secret. Silence-satisfiable, so compatible with `R-1` without exemplifying it |
| **`P-03`** (L251) | `R-5` | L252 `Abstention: **yes**`; L255 `Required: the surprise, and the confidentiality of how the character learned it`; L257 `Mandatory: the secret is kept` | **Tag correct.** Confidentiality/entitlement item. Its own note: "Factually correct, socially wrong… proves link 1 correctness does not excuse links 2–4". **Do not move to `R-1`** |
| **`F9-01`** (L207) | `R-4` | L208 `Abstention: **yes**`; L211 `Required: nothing from the siege`; L213 `Mandatory: no siege material reaches the reply` | **Tag correct.** Availability-versus-relevance item. **Do not move to `R-1`** |

### 3.1 Finding

**`R-1` has no established clean exemplar.** `F6-02` carries a perceptibility tension; `P-01` tests proportionate disclosure. The two items the mode table lists as `R-1` — `P-03` and `F9-01` — are correctly tagged `R-5` and `R-4`.

An incidental discriminator: the two `R-5`/`R-4` items carry `**Abstention:** yes`; both `R-1` items carry `no`. The corpus already encodes entitlement- and relevance-based withholding separately from mention suppression.

### 3.2 What this finding does not authorise

- It does **not** declare the tags wrong. Tag membership is the adopted authority (`OD-M09-05`); imperfect exemplification is not mis-tagging.
- It does **not** authorise re-tagging, rewriting, or repairing any item.
- It does **not** establish that `R-1` is unmeasurable — only that no current item cleanly instantiates it.

This is a **construct-validity** observation, distinct from the **operational membership** question already recorded in `OD-M09-05` D-1.

---

## 4. Q3 — Withholding modes: evidence and its ceiling

### 4.1 Five distinct facts, routinely conflated

| Level | Fact | What would show it | Held? |
|---|---|---|---|
| 1 | Information exists in storage | Storage/ledger inspection | Not inspected here |
| 2 | Application selected or retrieved it | Retrieval/selection record | No record cited |
| 3 | It was injected into the prompt | Prompt capture — `inference_capture.go` records assembled prompt and audits marker presence (`require_ok=%t forbid_clear=%t`) | Mechanism exists; shows injection only |
| 4 | The model used or was influenced by it | A validated causal-evidence channel | **None exists** |
| 5 | The reply exhibits or avoids disclosure | Reply text / scoring | Observable |

**Level 3 is not level 4.** Prompt capture answers "was this text in the prompt"; it cannot answer "did the model use it". **Level 5 is not evidence about levels 2–4**: absence of disclosure is compatible with successful withholding, with ignorance, with non-retrieval, and with generic reticence.

### 4.2 The evidence ceiling for `R-1`, `R-4`, `R-5`

> **Compatible with withholding; no prohibited disclosure observed; memory-guided withholding not established.**

That is the strongest claim reply text alone can support for a withholding mode. The reason is structural, not stylistic: these modes' positive signature is an **absence**, and absence has many causes.

### 4.3 Three evidence requirements

| Requirement | What it would supply | Available today? |
|---|---|---|
| **Availability / retrieval evidence** — independent of the reply | Distinguishes "correctly withheld" from "never retrieved" — the distinction `F9-01`'s own note demands | **No** reply-side channel; would need system records or a companion item forcing use |
| **Leakage / over-use detection** | Identifies prohibited disclosure or excessive use, rather than scoring wording | **Partly** — over-exposure fields exist, but they are scored text, not detector output |
| **Negative control / controlled contrast** | Separates memory-dependent behaviour from behaviour that would occur anyway | **No** — and constructing fixtures is corpus work, not authorised |

**A paired memory-present/memory-absent comparison may strengthen attribution but is not automatically sufficient.** Its design, controls, interpretation and validity would require separate specification and authorisation. It is described here as a candidate, not a solution.

### 4.4 What is not claimed

- Not claimed: that availability, injection or retrieval proves model use.
- Not claimed: that any current channel establishes causal memory influence.
- **Also not claimed: that no future validated method could ever establish influence.** The limit is a statement about present channels.

---

## 5. Findings, recommendations, and unresolved matters kept apart

| Kind | Content |
|---|---|
| **Observed source text** | `R6-03` has no `Required`/`Forbidden`; `F6-02` L149 requires "the aversion"; `P-01` L237 bounds disclosure to "at most one" fact; `P-03` and `F9-01` carry `Abstention: yes` |
| **Calculated counts** | 46 headings · 45 scored units · 42 family units · `R-6`: 2/6/9/10/8/7 under their six definitions |
| **Source-derived findings** | The 45/46/42 split is definitional with one cause (`R6-03` container + variant collapsing); tag totals exceed scored units by one for the same reason; `R-1` has no clean exemplar |
| **Interpretations** | `F6-02` is an `R-1`/`R-2` hybrid; `P-01` tests proportion rather than suppression. These are readings of the item text, not measurements |
| **Recommendations** (not decisions) | (1) State the accounting convention wherever a count appears — scored units 45, headings 46, family units 42. (2) Record `R-1`'s exemplar gap as a known design question. (3) Record the withholding evidence ceiling as the bound on `R-1`/`R-4`/`R-5` results. (4) Consider whether a container heading should carry a `Restraint:` tag at all, since it inflates tag totals |
| **Unresolved** | Whether `F6-02`'s "aversion" is respect or display; whether `P-01` belongs to a proportion construct outside the six modes; which unit any future 45-unit census excludes; whether availability evidence is obtainable at all within G1 bounds |

---

## 6. Source defects recorded, deliberately not applied

| # | Artifact / location | Defect | Impact | Correction | Authorization needed |
|---|---|---|---|---|---|
| SD-1 | `CORPUS-01` §5 coverage line vs §7 | Two conventions (42 family, 45 scored) stated without labels | Counts read as contradictory | Label both conventions | **Yes** — standards edit |
| SD-2 | `CORPUS-01` §7 `R-6` line | "9 items, 4 mechanisms + 4 provenance variants" — arithmetic gives 8 | `R-6` count unreliable | Restate as 8 fixtures / 7 scored / 10 tagged | **Yes** — standards edit |
| SD-3 | `CORPUS-04` L55 | Container `R6-03` carries `Restraint: R-6`, so tag totals include a non-scored heading | Tag count ≠ scored-unit count by 1 | State convention, or omit tag on containers | **Yes** — item/standards edit |
| SD-4 | `CORPUS-01` §5 `R-1` row | Lists `P-03`, `F9-01`, both tagged elsewhere | Table contradicts tags — already recorded in `OD-M09-05` D-1 | Reconcile table | **Yes** — standards edit |
| SD-5 | `CORPUS-01` §5 table | `P-03` and `F9-01` listed twice, incompatible with `A-2` — recorded in `OD-M09-05` D-4 | Multi-membership vs single tag | Same as SD-4 | **Yes** |

All five remain **unchanged in the source**. None is repaired here.

---

## 7. Non-operative status

This record:

- adopts no standard, threshold, statistic, count convention or vocabulary;
- amends no artifact and authorises no item, tag, table, coverage, score, split or checker change;
- asserts no empirical result, no construct validity, no reproducibility;
- does not restart `P-3`, does not merge it, does not design a study;
- does not run inference, benchmarks or `G2`.

Its findings are analytical: recomputed source facts plus labelled interpretation.

---

## 8. Provenance and verification

| Layer | Fact |
|---|---|
| Authorization | Owner, 2026-11-13 |
| Performed by | DeepSeek (per work order designation), repository worker — artifact-reported identity; labels are provenance data, not lineage proof |
| Sources inspected | §1 table, by digest prefix |
| Model executed | No |
| `P-3` | Not started, not restarted |
| Verification | Counts recomputed from source; item quotations read from file at cited lines; registry metadata verified against bytes/lines/SHA-256 |
| Administrative checker | Run only because registry inputs changed; a `PASS` is item-record consistency, **not** construct validation |

---

## 9. Revision history

| Rev | Date | Change |
|---|---|---|
| 1 | 2026-11-13 | Initial record. Compiled from a read-only reconciliation; every count and quotation re-verified against source before filing. Corrections made to the working draft before filing: added the sixth `R-6` count (**7 scored units**, excluding the container) after confirming `R6-03` has no scored body; added the tag-total explanation (46 = 45 + container tag) rather than leaving it as an unexplained discrepancy; recorded `Abstention:` values as an observed discriminator; downgraded "paired control" from remedy to candidate, with its sufficiency explicitly unresolved; added §6 covering five source defects left unapplied |
