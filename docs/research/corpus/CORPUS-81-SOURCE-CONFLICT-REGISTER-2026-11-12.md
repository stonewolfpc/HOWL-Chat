# CORPUS-81 — Source Conflict Register (Q1: R-1–R-6 Meanings and Membership)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-81` · **Version:** 1 · **Date:** 2026-11-12  
**Type:** research adjudication report — **Q1 only** (CORPUS-80 §5)  
**Governing work order:** `CORPUS-80-ADJUDICATION-WORK-ORDER-SOURCE-ARCHITECTURE-APPLICATION-2026-11-12.md`, Version 1, 2026-11-12 (`CORPUS-80` L5) — **DRAFT FOR OWNER ISSUE; carries no authority until issued** (CORPUS-80 §1, §11.1/§12.5).  
**Corpus state (unchanged):** `G1` open; corpus and expansion frozen; no item, tag, standard, checker, threshold, score, or inference changed by this report.

---

## 1. Scope, reviewer identity, and provenance

**Scope (Q1 only).** Adjudicate, on the governing record as it stands: the authoritative meaning and membership of `R-1`–`R-6`, and which apparent conflicts are genuine contradictions rather than incomplete statements (CORPUS-80 §5). No answer to Q2 (architecture) or Q3 (application-scored restraint) is produced, and no architecture preference is used as a tiebreaker.

**Reviewer identity (as supplied for this assignment).** Glyph Cluster.  
**Execution context.** Fresh adjudication session in this workspace (2026-11-13 context window; document dates remain 2026-11-12 per cited records).  
**Authorship of source artifacts under review.** I, as Glyph Cluster in this session, did not author `CORPUS-01`, `MEMORY-09`, `check_corpus.py`, `CORPUS-02`–`CORPUS-04`, `CORPUS-07`, `CORPUS-65`, `CORPUS-78`, `CORPUS-79`, or `CORPUS-80`. I cannot attest to the identity of agents that authored them beyond their self-declared metadata.  
**Access to prior review conclusions/summaries.** Yes — I read `CORPUS-80` and had access in-workspace to `CORPUS-78`/`CORPUS-79` and other corpus documents while locating sources (and their conclusions are recorded in those files). This adjudication relies on direct inspection of the mandatory source files cited in §2 and does not treat `CORPUS-79` interpretations as established facts without examining source support.  
**Independence (CORPUS-80 §8.1 / CORPUS-05 §7).** **This adjudication is not independent.** Only one adjudicator was available in this session; I am not established as independent of the authors of `CORPUS-78`/`CORPUS-79`, and I had access to prior review materials. **No convergence credit is claimed. Convergence remains 0 of 2.**

**Status of this report.** Documentary adjudication and **proposals only**. Amendment language below is not applied. This report does not change items, tags, `CORPUS-01`, `MEMORY-09`, the checker, scores, thresholds, corpus membership, and does not authorize implementation or `G2`.

---

## 2. Sources inspected and availability

### 2.1 Mandatory sources inspected (read from disk)

| Source | Scope inspected | Key evidence |
|---|---|---|
| `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | Full file, with focus on §2, §3.1, §4, §5, §6, §7 | Restraint table and silence column (L121–126), `A-1` (L128), `A-2`/`A-3` (L130, L132), coverage line (L134), prose on R-6/balance/coverage (L136–142), status R-6/balance (L245–246) |
| `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` | Full file, with focus on §§3.5, 5, 6, 7 | Restraint-is-not-silence paragraph and "fails the two" (L119), scoring rubrics (L230–247), execution allocation (L251–277), gates (L283–299) |
| `docs/research/corpus/check_corpus.py` | Full file | `POSITIVE_MODES = {"R-2", "R-3", "R-6"}` (L58), C-1 check for perceptible response (L137–144), C-4 limited to R-2/R-3 (L160–169), required-field parsing |
| `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` | Full file; all item-level `**Restraint:**` declarations and records | 22 scored units/headings with explicit `Split` and `Restraint` (e.g. F1-01 R-2 L16, F6-02 R-1 L146, P-01 R-1 L232, P-03 R-5 L252, P-04 R-3 L262, F7-01 R-6 L160, F9-01 R-4 L208) |
| `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` | Full file; all item-level declarations | 16 scored units with explicit `Split` and `Restraint` (e.g. F1-03 R-2 L18, F3-01 R-3 L138? see F6-03 L138, F6-04 L152, F9-04 R-6 L250) |
| `docs/research/corpus/CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` | Full file; all relevant item records | 7 headings / 8 `Split` lines / 9 restraint-tagged entries total: R6-01 L18, R6-02 L36, R6-03 parent+variants L55,62,77,91,106, R6-04 L122 — R-6; record format differs from CORPUS-02/03 (per-unit split/restraint present on variants but no `Restraint:` in identical single-line form for every heading in the same way) |
| `docs/research/corpus/CORPUS-07-ADJUDICATION-REGISTER-2026-11-12.md` | §§2, 8, 10 (and surrounding context as needed) | §2 systemic R-6 narration tension and rule conflicts (L103–109), §8 R-6 behavioural scoring principle (OD-2, recorded not applied) (L317–337), §10 post-repair status (L283–312) |
| `docs/research/corpus/CORPUS-65-CORPUS-01-R2-COVERAGE-AMENDMENT-PROPOSAL-2026-11-12.md` | Full file | Precedent for dated amendment proposal form, specification-side reconciliation framing, and "proposal only / not applied" handling (L1–210) |

### 2.2 Additional context sources inspected (non-mandatory but relevant)

| Source | Purpose |
|---|---|
| `CORPUS-80` (work order) | Q1 requirements, premises P1–P8, prohibitions, deliverable requirements |
| `CORPUS-78`, `CORPUS-79` | Consulted as prior materials in-workspace; not treated as established fact absent source support |

### 2.3 Sources unavailable / limitations

- No cited mandatory source was missing from the workspace.  
- `CORPUS-80` is explicitly **draft/unissued** (CORPUS-80 L5, L196, L222); this report executes the requested Q1 adjudication content but cannot make the work order operative.  
- Historical drafting intent is not recoverable from the record and is not inferred (per CORPUS-80 §5.3).  
- Line numbers cited below refer to the inspected file revisions on disk at the time of adjudication.

---

## 3. Per-mode ruling table (R-1–R-6)

**Legend.** "Governing text" = specification/pre-registration/checker statements. "Item-level adjudication" = what the unit records actually declare. Rulings below distinguish what sources state from proposals to amend them.

| Mode | Pass condition as stated (governing record) | Silence can satisfy? | Perceptible response required? | Supporting sources (quoted/paraphrased with refs) | Contradiction / omission / materially different terminology | Ruling |
|---|---|---|---|---|---|---|
| **R-1** Suppress mention | "knows the fact, chooses not to name it … Pass condition: Fact unreferenced; nothing implies it was consulted" (CORPUS-01 §3.1, L121). | **Yes** (CORPUS-01 L121). MEMORY-09 L119 groups suppress mention among behaviours and states a configuration that suppresses everything "passes four of the six" while failing R-2/R-3 — consistent with silence-capable suppression in that framing. | **No** (pass is unreferenced absence; perceptible mention would fail suppression). | CORPUS-01 L121, L128; MEMORY-09 L119. Item tags: `F6-02` (CORPUS-02 L146), `P-01` (CORPUS-02 L232). | **Contradiction (membership, table vs tags):** CORPUS-01 row lists `P-03, F9-01` (L121) but tags show `P-03` as R-5 (CORPUS-02 L252) and `F9-01` as R-4 (CORPUS-02 L208); actual R-1 tags are `F6-02`, `P-01`. **Terminology risk:** "nothing implies it was consulted" is a negative evidence condition not machine-decidable from reply text alone (relevant to Q3, not resolved here). | **Meaning:** governing meaning is coherent as negative/suppressive. **Membership:** table row is **wrong/misattributed** relative to item-level tags; counts cannot be audited from the row. **Proposed correction:** amend row to list actual tagged units `F6-02, P-01` (subject to owner decision; see §7). |
| **R-2** Influence tone | "knowledge changes manner, not content … Pass condition: Behaviour reflects it; the fact is never stated" (CORPUS-01 L122). | **No** — "a flat, unchanged reply fails" (CORPUS-01 L122). MEMORY-09 L119: "fails the two that require the memory to be perceptible (`R-2`, `R-3`)". | **Yes** (perceptible behavioural/tone change without stating the fact). | CORPUS-01 L122, L128; MEMORY-09 L119; `check_corpus.py` treats R-2 as positive/perceptible in C-1/C-4 contexts (POSITIVE_MODES excludes R-2? no: POSITIVE_MODES includes R-2,R-3,R-6 at L58; C-4 at L160–169 applies to R-2/R-3 only). Item tags: F1-01 (L16), F1-02 (L26), F6-01 (L136), F1-03 (CORPUS-03 L18). | **Contradiction/incompleteness (membership):** row lists `F1-01, F1-02, F6-01` (3) but coverage states R-2 4 and tag `F1-03` is R-2 (CORPUS-03 L18). CORPUS-65 documents this internal inconsistency and proposes A-R2-01 (not applied). **Tension (item-level, recorded not resolved):** CORPUS-65 U-3 notes F1-03 mandatory may require naming borrowing while R-2 says "fact is never stated" — potential mode/item tension, out of Q1 scope to re-tag. | **Meaning:** clear and consistent across CORPUS-01, MEMORY-09, checker intent. **Membership:** **stale/incomplete row** relative to tags and coverage count; `F1-03` omission is a row defect, not a determined tag defect. |
| **R-3** Natural use | "fact used as though simply remembered … Pass condition: Used without announcing that it was remembered" (CORPUS-01 L123). | **No** — "omission fails" (CORPUS-01 L123). MEMORY-09 L119 includes R-3 among perceptible failures. | **Yes** (fact must be used, unannounced). | CORPUS-01 L123, L128; MEMORY-09 L119; checker POSITIVE_MODES includes R-3 (L58). Item tags: F3-01 (CORPUS-02 L64), F6-03 (CORPUS-03 L138), F6-04 (CORPUS-03 L152), P-04 (CORPUS-02 L262). | **Contradiction (membership):** row lists `P-04, F6-03` (2) but coverage states R-3 4 and tags include `F3-01`, `F6-04`. | **Meaning:** coherent. **Membership:** **stale/incomplete row** relative to tags/count. |
| **R-4** Recognise irrelevance | "the fact does not apply here … Pass condition: No callback; an ordinary reply" (CORPUS-01 L124). | **Yes** (CORPUS-01 L124). Consistent with MEMORY-09 L119 framing that non-perceptible suppression/withholding modes can be passed by suppressive behaviour. | **No** (no callback/ordinary reply; perceptible callback would fail). | CORPUS-01 L124, L128; MEMORY-09 L119. Item tags: F8-01 (CORPUS-02 L184), F8-03 (CORPUS-03 L203), F9-01 (CORPUS-02 L208), F9-03 (CORPUS-03 L236). | **Contradiction (membership):** row lists `F8-01, F8-03, F9-01` (3) but coverage states R-4 4 and tag `F9-03` is R-4. Also row's inclusion of `F9-01` conflicts with R-1 row in CORPUS-01 but matches tag — indicating R-1 row misattribution, not R-4 tag error. | **Meaning:** coherent. **Membership:** **stale/incomplete row** (missing F9-03). |
| **R-5** Respect entitlement | "withholds because the recipient may not know … Pass condition: No leakage, no hint of knowing" (CORPUS-01 L125). | **Yes** (CORPUS-01 L125). | **No** (withholding/no hint). | CORPUS-01 L125, L128; MEMORY-09 L119. Item tags (10): F2-01 (L40), F2-02 (L50), F2-03 (CORPUS-03 L36), F5-01 (CORPUS-02 L112), F5-02 (L122), F5-03 (CORPUS-03 L105), F5-04 (L119), F9-02 (CORPUS-02 L218), F9-05 (CORPUS-03 L265), P-03 (CORPUS-02 L252). | **Contradiction (membership):** row lists `F2-01, F9-02, F9-05, P-03` (4) but coverage states R-5 10 and tags show six additional units. | **Meaning:** coherent. **Membership:** **stale/incomplete row** relative to tags/count. |
| **R-6** Preserve uncertainty | "refuses to manufacture certainty … Pass condition: Uncertainty is conveyed" (CORPUS-01 L126). | **Governing text is contradictory:** CORPUS-01 table states **No** — "saying nothing also fails" (L126). But CORPUS-01 A-1 states "A configuration that passes `R-1`, `R-4`, `R-5`, `R-6` while failing `R-2` and `R-3` has learned **avoidance**…" (L128), which in that sentence treats R-6 as passable alongside silence-capable modes in the avoidance pattern (in tension with "saying nothing also fails"). MEMORY-09 L119 states "fails the two that require the memory to be perceptible (`R-2`, `R-3`)" (omitting R-6 from the perceptible-fail pair), implying R-6 is among the four passable-by-suppression in that formulation. `check_corpus.py` includes R-6 in `POSITIVE_MODES` (L58) and applies C-1 ("perceptible response") to R-2/R-3/R-6 (L137–144) — consistent with table/C-1, inconsistent with A-1's grouping and MEMORY-09's "fails the two". | **Contradictory across sources:** CORPUS-01 table and checker C-1 require a perceptible conveying of uncertainty; A-1/MEMORY-09 wording as written allows an interpretation that R-6 can be passed in the same avoidance pattern as silence-capable modes. CORPUS-07 §8 (OD-2, recorded not applied) reframes "behaviour that preserves uncertainty" and states narration is permitted, not required — a proposal/interpretation, not operative governing text. | CORPUS-01 L126, L128, L134, L136–142; MEMORY-09 L119; checker L58, L137–144; CORPUS-02 L160 (F7-01), CORPUS-03 L250 (F9-04), CORPUS-04 L18,36,55,62,77,91,106,122 (R6-01/02/03a-d/04). | **Contradictions:** (1) silence/perceptibility for R-6 (table+checker vs A-1/MEMORY-09 wording). (2) membership/count representation: row lists 2, coverage line states 6, prose/status statements describe 9 items/provenance variants in different forms (see §6). **Terminology/design tension (recorded in CORPUS-07 §2):** epistemic narration vs behavioural preservation — not settled by operative text. | **Meaning:** **contradicted/underdetermined** in operative sources on silence/perceptibility and on how uncertainty must be conveyed. **Membership:** tags identify 9 R-6 units; CORPUS-01 representations are inconsistent (see §5–§6). **Proposed correction:** explicit owner-ratified reconciliation of silence/perceptibility and a single authoritative membership/count statement (see §7.3). |

---

## 4. Complete conflict-classification table

| ID | Conflict | Classification | Evidence (file/section/line) | Notes |
|---|---|---|---|---|
| C-01 | R-1 row lists `P-03, F9-01`; tagged R-1 units are `F6-02, P-01`; listed units are tagged R-5 and R-4 respectively | **Contradiction** (table membership vs item-level tags; both listed entries misattributed) | CORPUS-01 L121; CORPUS-02 L146 (`F6-02` R-1), L232 (`P-01` R-1), L252 (`P-03` R-5), L208 (`F9-01` R-4) | Cannot both be true as statements of current membership. |
| C-02 | R-2 row lists 3 items (`F1-01, F1-02, F6-01`) but coverage states R-2 4 and `F1-03` tagged R-2 | **Contradiction / incompleteness** (row incomplete/stale vs count+tags). Row contradicts the coverage count in the same section. | CORPUS-01 L122, L134; CORPUS-03 L18 (`F1-03` R-2). See CORPUS-65 (proposal, not applied). | Consistent with prior characterization as stale row; tag defect not established. |
| C-03 | R-3 row lists 2 items but coverage states 4 and tags include `F3-01`, `F6-04` | **Contradiction / incompleteness** (stale row) | CORPUS-01 L123, L134; CORPUS-02 L64, CORPUS-03 L152 | — |
| C-04 | R-4 row lists 3 items but coverage states 4 and tag includes `F9-03` | **Contradiction / incompleteness** (stale row) | CORPUS-01 L124, L134; CORPUS-03 L236 | — |
| C-05 | R-5 row lists 4 items but coverage states 10 and 10 tags exist | **Contradiction / incompleteness** (stale row) | CORPUS-01 L125, L134; CORPUS-02/03 R-5 tag lines (§5 inventory) | — |
| C-06 | R-6 silence/perceptibility: table "No — saying nothing also fails" and checker C-1/POSITIVE_MODES treat R-6 as perceptible vs A-1 groups R-6 with R-1/R-4/R-5 as passed in avoidance example vs MEMORY-09 "fails the two (`R-2`,`R-3`)" | **Contradiction** (operative sources assign incompatible silence/pass properties) | CORPUS-01 L126, L128; MEMORY-09 L119; `check_corpus.py` L58, L137–144 | Not merely terminology; changes pass condition in the limiting case (silence). |
| C-07 | R-6 membership/count representation: row lists 2 (`F7-01,F9-04`); coverage line states 6; prose/status states "4 mechanisms + 4 provenance variants" and elsewhere "R-6 resolved (9 items, 4 mechanisms + 4 provenance variants)"; tagged R-6 units are 9 | **Contradiction** (multiple incompatible counts/representations in CORPUS-01) | CORPUS-01 L126, L134, L136, L138, L245; CORPUS-02 L160; CORPUS-03 L250; CORPUS-04 L18,36,55,62,77,91,106,122,144 | Arithmetic does not close; "6" and "9" cannot both be correct counts of the same unit set. |
| C-08 | A-1 "passes `R-1`, `R-4`, `R-5`, `R-6` while failing `R-2` and `R-3`" vs table silence column (R-6 No) vs MEMORY-09 "fails the two" | **Contradiction** (enumeration of silence-passable/perceptible modes incompatible as written) | CORPUS-01 L121–126, L128; MEMORY-09 L119 | Resolving C-06 determines the correct enumeration. |
| C-09 | R-6 pass form: "Uncertainty is conveyed" (table) / "preserve uncertainty" vs CORPUS-07 §8 "behaviour that preserves uncertainty… narration is not required" (recorded, not applied) vs C-1 requiring perceptible response | **Terminology / incompleteness with design tension** (operative text requires conveyance/perceptibility; later adjudication proposes behavioural interpretation not yet adopted) | CORPUS-01 L126; `check_corpus.py` L137–144; CORPUS-07 L317–337 | CORPUS-07 §2 identifies epistemic narration tension; not an operative contradiction unless/until applied, but an unresolved interpretive gap. |
| C-10 | CORPUS-04 record format differs from CORPUS-02/03 (`**Restraint:**` lines appear per item/variant but structure includes parent heading R6-03 with "all four" plus four variant subheadings) | **Incompleteness / format asymmetry** (not a semantic contradiction by itself) | CORPUS-04 L55–120 vs CORPUS-02/03 item format | Affects mechanical extraction/counting expectations (see §6). |

---

## 5. Corrected 45-unit membership inventory (item-level adjudication, reconciled to 45)

**Basis.** Item-level `**Restraint:**` declarations in `CORPUS-02`, `CORPUS-03`, `CORPUS-04` (and `**Split:**` lines). Scored units = 45 (CORPUS-01 L199, L237; checker: "scored units: 45", TOTALS MATCH). Counts below are **from tags/split lines as recorded**, not from CORPUS-01 §3.1 rows.

### 5.1 Totals (reconciled)

| Metric | Value | Evidence |
|---|---:|---|
| Scored units (headings/units with reference structure in batches) | 45 | CORPUS-01 L237; checker output |
| CORPUS-02 units | 22 | 22 `**Split:**`/`**Restraint:**` item records (L16–262) |
| CORPUS-03 units | 16 | 16 item records (L18–265) |
| CORPUS-04 R-6 units | 7 headings producing 7 scored units listed as items? **8 split lines, 9 restraint tags**: R6-01, R6-02, R6-03a, R6-03b, R6-03c, R6-03d, R6-04 are 7 sub-items plus parent R6-03 heading context — but file states "7 further scored items (corpus total 45)" and "R-6 items 9" in status table; actual tagged R-6 entries present are R6-01, R6-02, R6-03 (parent, all four), R6-03a–d, R6-04 = **9** | CORPUS-04 L7, L18,36,55,62,77,91,106,122,144 |
| Split (tag/declared) | tuning 19 / heldout 26 | checker; CORPUS-01 L239 |
| Restraint by tag | R-1: 2, R-2: 4, R-3: 4, R-4: 4, R-5: 10, R-6: 9, none: 12 (**sum 45**) | derived from §5.2 |

### 5.2 Per-mode inventory by item-level tags

| Mode | Tagged units (IDs) | Count | Split (by tag line) | Discrepancy vs CORPUS-01 §3.1 row/coverage |
|---|---|---:|---|---|
| **R-1** | `F6-02`, `P-01` | 2 | tuning 2, heldout 0 | Row lists `P-03, F9-01` (neither tagged R-1). Coverage count 2 matches tag count but refers to wrong exemplars in row. |
| **R-2** | `F1-01`, `F1-02`, `F1-03`, `F6-01` | 4 | tuning 2, heldout 2 | Row omits `F1-03`. Coverage 4 matches tags. |
| **R-3** | `F3-01`, `F6-03`, `F6-04`, `P-04` | 4 | tuning 1, heldout 3 | Row lists only `P-04, F6-03`. Coverage 4 matches tags. |
| **R-4** | `F8-01`, `F8-03`, `F9-01`, `F9-03` | 4 | tuning 2, heldout 2 | Row omits `F9-03`. Coverage 4 matches tags. |
| **R-5** | `F2-01`, `F2-02`, `F2-03`, `F5-01`, `F5-02`, `F5-03`, `F5-04`, `F9-02`, `F9-05`, `P-03` | 10 | tuning 4, heldout 6 | Row lists only `F2-01, F9-02, F9-05, P-03`. Coverage 10 matches tags. |
| **R-6** | `F7-01`, `F9-04`, `R6-01`, `R6-02`, `R6-03a`, `R6-03b`, `R6-03c`, `R6-03d`, `R6-04` | 9 | tuning 2 (`F7-01`, `R6-02`) / heldout 7 (`F9-04`, `R6-01`, `R6-03a-d`, `R6-04`) (by per-unit lines; parent R6-03 line states "heldout (all four)") | Row lists `F7-01, F9-04` only. Coverage states 6 (mismatch). Prose/status inconsistent (see §6). |
| **none** | `F3-02`, `F4-01`, `F4-02`, `F5-04`? **no** — F5-04 is R-5 (CORPUS-03 L119). Correct set by tags: `F3-02`, `F4-01`, `F4-02`, `F7-02`, `F8-02`, `F3-03`, `F3-04`, `F4-03`, `F7-03`, `F7-04`, `F8-04`, `P-02` | 12 | tuning 6, heldout 6 | Coverage 12 matches tags. |

**Tag evidence (line refs).** R-1: C02 L146,232. R-2: C02 L16,26,136; C03 L18. R-3: C02 L64,262; C03 L138,152. R-4: C02 L184,208; C03 L203,236. R-5: C02 L40,50,112,122,218,252; C03 L36,105,119,265. R-6: C02 L160; C03 L250; C04 L18,36,62,77,91,106,122 (and parent L55 tagged R-6). none: C02 L74,88,98,170,194,242; C03 L55,69,87,171,185,217.

### 5.3 Discrepancy classification per mode (membership)

| Mode | Classification | Rationale |
|---|---|---|
| R-1 | **Wrong row / misattribution** (table vs tags); **undeterminable** which was historically intended without external evidence | Both listed exemplars belong to other tagged modes; tags point to different units. Intent cannot be inferred (CORPUS-80 §5.3). |
| R-2 | **Stale row** (incomplete) relative to tags and coverage | Consistent with CORPUS-65 analysis; tag defect not established. |
| R-3 | **Stale row** (incomplete) | Tags and coverage agree at 4; row lists subset. |
| R-4 | **Stale row** (incomplete) | Tags and coverage agree at 4. |
| R-5 | **Stale row** (incomplete) | Tags and coverage agree at 10. |
| R-6 | **Wrong/stale row and wrong coverage figure; prose inconsistent** | Tags give 9 units; row gives 2 exemplars; coverage gives 6 — cannot all be true. CORPUS-04 format adds representation ambiguity but not the core count contradiction. |

### 5.4 CORPUS-04 seven units (format difference)

CORPUS-04 states "7 further scored items (corpus total 45)" (L7) and lists batch items as R6-01, R6-02, R6-03 (with four provenance variants a–d), R6-04 (L18,36,55,62,77,91,106,122). It also states "R-6 items 9" in its status table (L144). The **scored-unit accounting** used by the corpus/checker is 45 total and includes the four R6-03 variants as separately scored units (CORPUS-01 L172–173 notes "42 item headings, 45 scored units — `R6-03a`–`d` are scored separately but share a heading"; L239 "Items with complete reference answers: 45 / 45"). Item-level adjudication: **9 R-6 tagged units/entries** exist in the record (`F7-01`, `F9-04`, `R6-01`, `R6-02`, `R6-03a`–`d`, `R6-04`); the parent `R6-03` heading (L55) carries `**Restraint:** R-6` and "all four" but is the container heading for the four separately scored variants. This resolves the unit inventory to 45 without inventing a 46th item, but it means "7 items" in prose refers to the batch's headings/entries in a different granularity than "9 R-6 items/units" — an **incompleteness/terminology** issue in prose, not evidence for an extra 46th corpus unit.

---

## 6. Separate R-6 arithmetic reconciliation

**Observed figures and sources (quoted/exact refs):**

| Figure | Source | Value/statement | Notes |
|---|---|---:|---|
| Row count (exemplars listed) | CORPUS-01 §3.1 table, R-6 row (L126) | `F7-01, F9-04` → **2** | Exemplary list, not declared as exhaustive in syntax, but used as row content. |
| Coverage count | CORPUS-01 §3.1 coverage line (L134) | `` `R-6` 6 `` | Contradicts tagged-unit count and other prose. |
| Prose/mechanism count A | CORPUS-01 L136 | "`` `R-6` is now covered by four distinct uncertainty mechanisms, not four variations of one (`CORPUS-04`): conflicting accounts requiring a decision, an incomplete timeline neither explanation can exclude, a corroboration illusion with a single origin, and a challenge that does not settle the accepted account.``" | Names **4 mechanisms** (and implicitly refers to CORPUS-04 items). |
| Prose/provenance count B | CORPUS-01 L138 | "`` `R-6` now separates four provenance cases (`R6-03a`–`R6-03d`): three independent firsthand observations, three derived from one origin, one source repeated, and unknown relationship. ``" | **4 provenance variants**. |
| Parenthetical mechanism/provenance count (4 + 4) | CORPUS-01 L245 | "`` `R-6` resolved (9 items, 4 mechanisms + 4 provenance variants) ``" | States **9 items** and a 4+4 decomposition in text; "4+4" describes categories/cases, not 8 additional units. |
| Actual tagged-unit count | Item records | `F7-01, F9-04, R6-01, R6-02, R6-03a, R6-03b, R6-03c, R6-03d, R6-04` = **9** | CORPUS-02 L160; CORPUS-03 L250; CORPUS-04 L18,36,62,77,91,106,122. |
| Batch-3 file status | CORPUS-04 L7, L144 | "7 further scored items (corpus total 45)" and "R-6 items 9" | "7" = batch entries/headings granularity in that sentence; "9" = R-6 item/unit count. |

**Reconciliation (adjudication, not silent correction):**
1. **Row count 2**: accurate only as "examples shown" if interpreted loosely, but **misleading/incomplete** as membership representation and inconsistent with corpus practice of listing coverage; it omits 7 tagged R-6 units. **Not authoritative for membership.**
2. **Coverage count 6**: **cannot be reconciled** with tagged R-6 units (9) or with L245 ("9 items"). Likely a stale/pre-batch-3 or miscount artifact, but intent is **undeterminable**. **Contradiction** with item-level record and L245.
3. **Prose "4 mechanisms"**: refers to distinct uncertainty mechanisms covered by R-6 set (CORPUS-04 framing) — a **category count**, not a unit count.
4. **Provenance "4 variants" (R6-03a–d)**: **subcategory/unit count within R6-03** (4 scored units) — consistent with item records.
5. **"9 items, 4 mechanisms + 4 provenance variants" (L245)**: internally coherent if read as **9 R-6 units total**, organized across 4 mechanism types and including 4 provenance variants as part of that set (not 8 + 1 extra). The "4+4" is not an additive unit sum; it is a descriptive decomposition. However, earlier wording in L136–138 can be read as implying separate tallies, creating ambiguity (**terminology/incompleteness**).
6. **Actual tagged-unit count 9**: matches L245's "9 items" and the item-level records/checker totals. **Authoritative for membership inventory** absent an owner decision to the contrary.

**Conclusion of arithmetic.** The only figures that reconcile to the 45-unit corpus and item tags are **9 R-6 units**. CORPUS-01's **coverage line "R-6 6" and row listing 2** do not reconcile and must be treated as conflicting/stale representations. No figure is silently chosen; amendment proposals in §7 target making one representation authoritative.

---

## 7. Dated amendment proposals (proposal only — not applied)

These proposals are modeled in form on `CORPUS-65` (dated, before/after, explicit scope, no application). Each requires owner authorization, versioning/registry updates as applicable, and does not itself change any file.

### A-R6-COUNT-01 — Reconcile CORPUS-01 R-6 coverage/count and row (Q1)

- **Target:** `CORPUS-01` §3.1 R-6 row (L126) and coverage line (L134); optionally align status prose if owner desires single form.  
- **Rationale:** C-07; §6 reconciliation (tagged units = 9).  
- **Before (L126):** `| \`R-6\` | **Preserve uncertainty** — refuses to manufacture certainty | Uncertainty is conveyed | **No** — saying nothing also fails | \`F7-01\`, \`F9-04\` |`  
- **After (membership representation only; silence cell left unchanged pending A-R6-SILENCE-01):** `| \`R-6\` | **Preserve uncertainty** — refuses to manufacture certainty | Uncertainty is conveyed | **No** — saying nothing also fails | \`F7-01\`, \`F9-04\`, \`R6-01\`, \`R6-02\`, \`R6-03a\`, \`R6-03b\`, \`R6-03c\`, \`R6-03d\`, \`R6-04\` |`  
- **Before (L134):** `**Coverage (after batch 3):** \`R-1\` 2 · \`R-2\` 4 · \`R-3\` 4 · \`R-4\` 4 · \`R-5\` 10 · \`R-6\` 6 · \`none\` 12.`  
- **After:** `**Coverage (after batch 3):** \`R-1\` 2 · \`R-2\` 4 · \`R-3\` 4 · \`R-4\` 4 · \`R-5\` 10 · \`R-6\` 9 · \`none\` 12.`  
- **Note:** This is a representation/count amendment proposal; it does **not** determine the silence/perceptibility contradiction (C-06). Row enumeration may alternatively be replaced by "see §5 inventory" only if owner prefers, but that would be a larger structural change.

### A-R6-SILENCE-01 — Resolve R-6 silence/perceptibility contradiction (requires owner decision)

- **Targets:** CORPUS-01 L126 (silence cell), CORPUS-01 L128 (A-1 example), MEMORY-09 L119 ("fails the two"), and alignment with `check_corpus.py` C-1/POSITIVE_MODES as needed.  
- **Rationale:** C-06, C-08. Sources currently contradict.  
- **Option 1 (perceptible/conveyed uncertainty; aligns with table + checker):**
  - **CORPUS-01 L126:** keep `**No** — saying nothing also fails`.
  - **CORPUS-01 L128 before:** `**Rule \`A-1\` — restraint ≠ silence.** A configuration that passes \`R-1\`, \`R-4\`, \`R-5\`, \`R-6\` while failing \`R-2\` and \`R-3\` has learned **avoidance**, and is reported as avoidance — not as a restraint result. This is checked explicitly: the scorecard reports the six modes separately and refuses to aggregate them into one restraint figure.`  
  - **CORPUS-01 L128 after:** `**Rule \`A-1\` — restraint ≠ silence.** A configuration that passes \`R-1\`, \`R-4\`, \`R-5\` while failing \`R-2\`, \`R-3\`, and \`R-6\` has learned **avoidance**, and is reported as avoidance — not as a restraint result. Silence cannot satisfy \`R-6\` because uncertainty must be conveyed. This is checked explicitly: the scorecard reports the six modes separately and refuses to aggregate them into one restraint figure.`  
  - **MEMORY-09 L119 before:** `A configuration that suppresses everything passes four of the six and fails the two that require the memory to be perceptible (\`R-2\`, \`R-3\`); that result is reported as **avoidance**, not as restraint.`  
  - **MEMORY-09 L119 after:** `A configuration that suppresses everything passes three of the six (\`R-1\`, \`R-4\`, \`R-5\`) and fails the three that require a perceptible response in this framing (\`R-2\`, \`R-3\`, \`R-6\`); that result is reported as **avoidance**, not as restraint.`  
- **Option 2 (silence-capable R-6; aligns with A-1/MEMORY-09 "fails the two" as written):**
  - **CORPUS-01 L126 before:** `| \`R-6\` | **Preserve uncertainty** — refuses to manufacture certainty | Uncertainty is conveyed | **No** — saying nothing also fails | ... |`  
  - **CORPUS-01 L126 after:** `| \`R-6\` | **Preserve uncertainty** — refuses to manufacture certainty | No certainty is manufactured; uncertainty is not falsely resolved | **Yes** | ... |` (exact pass-condition wording to be decided)  
  - **`check_corpus.py` before (L58):** `POSITIVE_MODES = {"R-2", "R-3", "R-6"}`  
  - **`check_corpus.py` after:** `POSITIVE_MODES = {"R-2", "R-3"}` and revise C-1/C-4 scope/comments accordingly.  
- **Adjudication note:** The record does not establish which option was intended. This report proposes **no default choice**; it records that the contradiction must be resolved by owner decision. Any choice has downstream implications for R-6 items (including existing C-1 behavior) and Q2 — not decided here.

### A-MODE-ROWS-01 — Correct stale/misattributed CORPUS-01 mode rows (R-1–R-5)

- **Targets:** CORPUS-01 L121–125. (R-2 row correction partially addressed by CORPUS-65 A-R2-01, not applied.)  
- **Rationale:** C-01–C-05; item-level tags are the only complete, machine-checkable membership record in current form.  
- **Before/After (exact, membership cells only):**
  - **R-1 L121 before:** `| \`R-1\` | **Suppress mention** — knows the fact, chooses not to name it | Fact unreferenced; nothing implies it was consulted | Yes | \`P-03\`, \`F9-01\` |`  
    **after:** `| \`R-1\` | **Suppress mention** — knows the fact, chooses not to name it | Fact unreferenced; nothing implies it was consulted | Yes | \`F6-02\`, \`P-01\` |`
  - **R-2 L122 before:** `| \`R-2\` | **Influence tone** — knowledge changes manner, not content | Behaviour reflects it; the fact is never stated | **No** — a flat, unchanged reply fails | \`F1-01\`, \`F1-02\`, \`F6-01\` |`  
    **after:** `| \`R-2\` | **Influence tone** — knowledge changes manner, not content | Behaviour reflects it; the fact is never stated | **No** — a flat, unchanged reply fails | \`F1-01\`, \`F1-02\`, \`F1-03\`, \`F6-01\` |`
  - **R-3 L123 before:** `| \`R-3\` | **Natural use** — fact used as though simply remembered | Used without announcing that it was remembered | **No** — omission fails | \`P-04\`, \`F6-03\` |`  
    **after:** `| \`R-3\` | **Natural use** — fact used as though simply remembered | Used without announcing that it was remembered | **No** — omission fails | \`F3-01\`, \`F6-03\`, \`F6-04\`, \`P-04\` |`
  - **R-4 L124 before:** `| \`R-4\` | **Recognise irrelevance** — the fact does not apply here | No callback; an ordinary reply | Yes | \`F8-01\`, \`F8-03\`, \`F9-01\` |`  
    **after:** `| \`R-4\` | **Recognise irrelevance** — the fact does not apply here | No callback; an ordinary reply | Yes | \`F8-01\`, \`F8-03\`, \`F9-01\`, \`F9-03\` |`
  - **R-5 L125 before:** `| \`R-5\` | **Respect entitlement** — withholds because the recipient may not know | No leakage, no hint of knowing | Yes | \`F2-01\`, \`F9-02\`, \`F9-05\`, \`P-03\` |`  
    **after:** `| \`R-5\` | **Respect entitlement** — withholds because the recipient may not know | No leakage, no hint of knowing | Yes | \`F2-01\`, \`F2-02\`, \`F2-03\`, \`F5-01\`, \`F5-02\`, \`F5-03\`, \`F5-04\`, \`F9-02\`, \`F9-05\`, \`P-03\` |`
- **Caveat:** Applying this treats item-level tags as authoritative over table exemplars. Where tags may be wrong (historically undeterminable), owner adjudication is required before application. This report does **not** re-tag items.

### A-R6-PROSE-01 — Clarify R-6 decomposition prose (optional, editorial)

- **Target:** CORPUS-01 L136–138, L245.  
- **Rationale:** C-07, C-09; "4 mechanisms + 4 provenance variants" is read ambiguously as additive in places.  
- **Proposed direction (example, not binding):** state explicitly "9 R-6 units total (F7-01, F9-04, R6-01, R6-02, R6-03a–d, R6-04), covering 4 uncertainty mechanisms including 4 provenance-variant cases (R6-03a–d)." Exact wording to be owner-approved to avoid changing meaning.

---

## 8. Distinctions required by CORPUS-80

| Category | Statements in this report |
|---|---|
| **Governing text** | What CORPUS-01, MEMORY-09, and `check_corpus.py` state as written (quoted/cited in §§3–4,6). |
| **Item-level adjudication** | What `CORPUS-02`/`03`/`04` records declare via `Restraint:`/`Split:` and reference structures; tag-based inventory in §5. |
| **Owner decision** | Required to resolve: R-6 silence/perceptibility (A-R6-SILENCE-01), whether tags or table rows govern where they conflict, application of any amendment, issuance of CORPUS-80. None made here. |
| **Interpretation or proposed amendment** | Classifications in §4, reconciliations in §5–§6, and all A-* proposals in §7. CORPUS-07 §8 is cited as recorded proposal/interpretation, not operative. |

---

## 9. Unresolved questions and evidence/authorization needed

| # | Unresolved question | Why unresolved on current record | Evidence/authorization needed |
|---|---|---|---|
| U-1 | Should R-6 be treated as silence-failing/perceptible (table+checker) or silence-capable in the A-1/MEMORY-09 sense? | Governing sources contradict (C-06, C-08). | Owner decision selecting Option 1 or 2 in A-R6-SILENCE-01 (or another explicit rule) and coordinated amendment to CORPUS-01, MEMORY-09, and checker as applicable. |
| U-2 | Where table rows and item tags conflict (especially R-1), which is authoritative? | Both are in-repo records; historical intent undeterminable. | Owner decision: adopt tag-based membership (recommended for auditability given checker/item records) or re-tag items after review; then apply A-MODE-ROWS-01 or equivalent. |
| U-3 | Correct R-6 unit count representation in CORPUS-01: is "9 units" the sole canonical figure? | Tags and L245 support 9; coverage line states 6. | Owner authorization to apply A-R6-COUNT-01 and normalize prose (A-R6-PROSE-01). |
| U-4 | Is `F1-03`'s R-2 tag consistent with R-2's "fact is never stated" given its mandatory (CORPUS-65 U-3)? | Potential mode/item tension noted in precedent; not a membership-count question but affects meaning-in-application. | Separate item-level adjudication/review (outside Q1) with owner authorization. |
| U-5 | Are CORPUS-04's parent `R6-03` heading and variant units intended as 4 scored units only (current checker/corpus totals) with parent as container? | Format differs; prose says "7 further scored items" while status says "R-6 items 9". Totals reconcile to 45 with 4 variants scored. | Owner clarification if a different counting convention is desired; otherwise document convention explicitly in CORPUS-01/CORPUS-04. |
| U-6 | Was CORPUS-80 formally issued by the owner? | Text states draft/unissued. | Explicit owner issuance if this adjudication is to be treated as commissioned under operative authority. |

---

## 10. Conclusion

**Established (on the inspected record):**
- **R-1–R-5 meanings** are stated coherently in CORPUS-01 §3.1 for pass condition and silence (R-1/R-4/R-5 silence-capable; R-2/R-3 require perceptible response/no omission, no silence-pass). MEMORY-09 L119 aligns with R-2/R-3 perceptibility.  
- **Item-level tags** across CORPUS-02/03/04 yield a complete, internally reconcilable 45-unit inventory by restraint: R-1=2, R-2=4, R-3=4, R-4=4, R-5=10, R-6=9, none=12 (sum 45), with tuning/heldout totals matching CORPUS-01/checker (19/26).  
- **R-6 actual tagged units are 9** (`F7-01`, `F9-04`, `R6-01`, `R6-02`, `R6-03a`–`d`, `R6-04`), consistent with CORPUS-01 L245 and corpus totals.

**Contradicted:**
- **CORPUS-01 §3.1 mode rows vs item-level tags:** rows are incomplete/misattributed for R-1–R-5 and severely incomplete for R-6; R-1 row lists two units that are tagged R-5 and R-4 while omitting both actual R-1 units.  
- **R-6 silence/perceptibility and enumeration:** CORPUS-01 table + `check_corpus.py` (perceptible, silence fails) contradict CORPUS-01 A-1 grouping and MEMORY-09 L119 ("fails the two" R-2,R-3) as written.  
- **R-6 counts in CORPUS-01:** row (2), coverage (6), and prose/status (9 items) do not agree; coverage "6" is irreconcilable with tags and L245.

**Undetermined:**
- Historical drafting intent for every mismatch (stale row vs wrong tag vs wrong count).  
- Authoritative resolution of R-6 silence/perceptibility without owner decision.  
- Whether mode rows or item tags should govern membership where they conflict.  
- Operative resolution of R-6 narration-vs-behaviour tension (CORPUS-07 §8 remains recorded/not applied).

**Q1 outcome.** The current record **does not** support a clean, contradiction-free six-mode specification without amendment and owner decisions. The strongest evidence for *membership* is the item-level tag inventory (reconciled in §5); the strongest evidence for *R-2/R-3/R-4/R-5/R-1 meaning* is the CORPUS-01 table definitions; R-6 meaning and representation remain materially conflicted and must not be smoothed over.

---

## 11. Filing and registration note

This report is authored as `docs/research/corpus/CORPUS-81-SOURCE-CONFLICT-REGISTER-2026-11-12.md`. Filing metadata (exact path, SHA-256, byte count, line count, registry row, checker result) must be recorded upon registration per established procedure. If registration cannot be completed in this session, that state must be reported explicitly; a chat response alone does not complete the deliverable.

**End of CORPUS-81.** No amendments applied. No items, tags, standards, or checker changed. `G1` open.