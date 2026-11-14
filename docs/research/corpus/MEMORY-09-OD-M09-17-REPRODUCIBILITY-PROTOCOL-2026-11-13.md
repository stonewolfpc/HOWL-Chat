# `MEMORY-09` — `OD-M09-17` Reproducibility Protocol Design (Stage Three)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-PROTO-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **PROTOCOL DESIGN ONLY — the study has not been executed.**
**Authority:** `MEMORY-09-WO-RESEARCH-01` §5, implementing `OD-M09-17`.
**Effect:** no rater recruited, no rating performed, no model run, no score, no threshold adopted.

---

## 1. The eight required elements, verbatim

Read from the governing record, **not** reconstructed: `MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md` **L104–111**.

| # | Requirement (verbatim) |
|---|---|
| **E-1** | "Rater qualification and identity/provenance requirements." |
| **E-2** | "Independent assessment procedures." |
| **E-3** | "Blinding and the information available to each rater." |
| **E-4** | "The exact judgement being compared." |
| **E-5** | "How disagreements are recorded and handled." |
| **E-6** | "The reliability measure and its justification." |
| **E-7** | "Any acceptance threshold, specified **before** results are observed." |
| **E-8** | "How insufficient sample size, missing judgements and unresolved disagreements are reported." |

The decision itself adds (L113–115) that it "does **not** authorise executing the study, collecting new evidence, scoring, or any benchmark activity."

No element was added, merged, or omitted. Topic sections in §4 that the eight do not explicitly require are marked **`[protocol detail]`**.

---

## 2. Element-by-element design

### E-1 — Rater qualification and provenance

| | |
|---|---|
| **Requirement** | Raters must be qualified and identifiable. |
| **Proposed protocol** | Two raters minimum. Each must (a) have read `CORPUS-01`'s mode definitions and the operative tag table, (b) pass a calibration set of 3 items whose ratings are pre-agreed by the owner, (c) record declared model/assistant lineage and version. |
| **Compliance evidence** | Signed rater record per rater: identifier, lineage, date, calibration outcome, version of every document read. |
| **Failure modes** | Rater shares a lineage with the executor; calibration passed by guessing; lineage declared but unverifiable. |
| **Open decisions** | Whether lineage may be self-declared, or must be artifact-attested. |

**Constraint carried forward:** declared lineage is **provenance data only**. Two raters from one lineage are *not* independent reviewers in the `CORPUS-26` sense, and no result may be described as independent review.

### E-2 — Independent assessment

| | |
|---|---|
| **Requirement** | Assessments must be independent. |
| **Proposed protocol** | Each rater rates the full unit set alone, from the same frozen materials, with no access to the other's ratings, no discussion before both are filed, and no visibility into expected results. |
| **Compliance evidence** | Time-stamped submission artifacts, one per rater, hashed on receipt; a declaration that no cross-communication occurred. |
| **Failure modes** | Informal comparison; a shared scratch file; one rater seeing the other's totals. |
| **Open decisions** | Whether the two raters may be run by the same operator. **Proposal:** permitted, disclosed, and reported as a limitation. |

### E-3 — Blinding and information available

| | |
|---|---|
| **Requirement** | Blinding must be specified. |
| **Proposed protocol** | Raters receive: item text, scenario, reply text, the operative mode definitions, and the rating scale. Raters do **not** receive: the executor's provisional verdicts, prior review conclusions, aggregate statistics, or which items are expected to pass. Unit identifiers are randomised per rater. |
| **Compliance evidence** | The exact packet supplied, by hash, for each rater. |
| **Failure modes** | Packet drift between raters; identifiers leaking expected outcomes. |
| **Open decisions** | Whether raters see the mandatory acceptance field. **Proposal:** yes — it is the published criterion, not a hint. |

### E-4 — The exact judgement compared

| | |
|---|---|
| **Requirement** | One exact judgement, stated. |
| **Proposed protocol** | Per unit, each rater assigns exactly one label from a fixed four-place scale: **`demonstrated` / `not demonstrated` / `not assessable` / `insufficient evidence`**. No free text is scored; free text is recorded but not compared. |
| **Compliance evidence** | Rater sheets containing one label per unit, no ties, no ranges. |
| **Failure modes** | Raters collapsing `not assessable` into `not demonstrated`; partial labels. |
| **Open decisions** | Whether `insufficient evidence` and `not assessable` are distinct. **Proposal:** keep distinct — they imply different remedies. |

### E-5 — Disagreements

| | |
|---|---|
| **Requirement** | Recording and handling of disagreement. |
| **Proposed protocol** | Disagreements are recorded as-is and **never resolved by discussion into a single value for reliability purposes**. A third adjudicator may produce a *separate* adjudicated record, reported alongside — never substituted for the two independent labels. |
| **Compliance evidence** | Raw two-rater matrix retained unchanged; adjudication logged separately with its own provenance. |
| **Failure modes** | Post-hoc consensus that inflates agreement; adjudication overwriting originals. |
| **Open decisions** | Whether an adjudicator is appointed at all. **Proposal:** not in the first pass. |

### E-6 — Reliability measure and justification

| | |
|---|---|
| **Requirement** | A named measure, justified. |
| **Proposed measure** | **Krippendorff's α (nominal)**, reported with raw percent agreement as a descriptive companion. |
| **Justification** | The data are nominal (four unordered labels); missing cells are likely (E-8); α handles missing data and any number of raters without imputation. Cohen's κ assumes two raters, no missing data, and is sensitive to prevalence skew — with four labels and likely `not assessable` clustering, κ would understate agreement. Percent agreement alone ignores chance and is reported only as context. |
| **Compliance evidence** | The coincidence matrix, α, and percent agreement, all reproducible from the retained sheets. |
| **Failure modes** | α computed after dropping units; measure swapped after seeing results. |
| **Open decisions** | None — the measure is a design choice this protocol is authorised to make. |

**This selects a measure. It does not select a threshold.**

### E-7 — Acceptance threshold

| | |
|---|---|
| **Requirement** | Any threshold, fixed **before** results. |
| **Status** | **PROPOSAL — requires owner approval. Not adopted. Not authorised by this document.** |
| **Proposed** | α ≥ 0.667 as "tentative", α ≥ 0.800 as "acceptable", reported against both. |
| **Justification** | 0.667/0.800 are Krippendorff's own published guidance for tentative and acceptable reliability on exploratory data; they are **not** corpus-specific and carry no evidential weight of their own. |
| **Failure modes** | The proposal silently becoming adopted because it is written here; thresholds being revised after α is seen. |
| **Open decisions** | Whether any threshold is adopted at all; whether one threshold or two bands. |

### E-8 — Insufficient sample, missing judgements, unresolved disagreement

| | |
|---|---|
| **Requirement** | Reporting of inadequacy. |
| **Proposed protocol** | Report: units offered, units rated by each rater, units rated by both, missing cells and their cause, unresolved disagreements with their labels, and a plain statement of whether the sample supports any reliability claim. **Rule:** if fewer than 20 units are double-rated, no reliability claim is made, whatever α equals. |
| **Compliance evidence** | A coverage table reconciling offered = rated + missing, per rater. |
| **Failure modes** | α computed on the easy subset; missingness reported as agreement. |
| **Open decisions** | The 20-unit floor is **proposed**, not authorised. |

---

## 3. Proposed protocol details not explicitly required by `OD-M09-17`

Marked so they are not mistaken for owner requirements:

- **`[protocol detail]` Materials freeze** — one hashed packet version, no substitutions mid-study.
- **`[protocol detail]` Order randomisation** — per-rater unit order to blunt drift.
- **`[protocol detail]` Time-boxing** — no rating session limit proposed; noted as a possible confound.
- **`[protocol detail]` Retention** — rater sheets, packets, and logs retained as artifacts with registry entries.

## 4. Explicit non-execution statement

**The study has not been executed.** No rater was recruited or appointed; no rating session occurred; no model was run; no reply was generated or scored; no α was computed; no threshold was adopted. Nothing in this document is evidence about reproducibility, model behaviour, or memory.

## 5. Known limitations of the design itself

1. **Two raters cannot separate rater error from construct ambiguity.** A low α is consistent with an unclear item as much as with an unreliable judgement.
2. **Lineage similarity is uncontrolled.** If both raters share a lineage, agreement may reflect shared priors — which is precisely the `CORPUS-26` independence problem, unresolved here.
3. **The scale's `not assessable` label is a magnet.** Items like `F6-04` invite it, and clustering there will depress α for reasons that are not rater unreliability.
4. **A good protocol proves nothing.** It establishes that a study *could* be run reproducibly, not that the corpus *is* reproducibly ratable.

## 6. Unresolved owner decisions

| # | Decision | Blocks |
|---|---|---|
| 1 | Adopt, reject, or defer the α threshold proposal | E-7 |
| 2 | Whether an adjudicator exists in pass one | E-5 |
| 3 | Whether lineage may be self-declared | E-1 |
| 4 | Whether the 20-unit floor is adopted | E-8 |
| 5 | Whether execution is authorised at all (separate work order) | everything |
