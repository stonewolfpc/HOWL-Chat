# Adjudication Work Order — Source Conflicts, `OD-9` Architecture, and Application-Scored Restraint

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-80` · **Version:** 1 · **Date:** 2026-11-12.
**Type:** **work order — DRAFT FOR OWNER ISSUE.** It carries no authority until the owner issues it. It commissions nothing by itself.
**Prepared by:** DeepSeek (owner-supplied assignment identity). Persona label Buffy = execution fact, not model identity.
**Trigger:** owner disposition of `CORPUS-79` (2026-11-12), directing three **separate** adjudications rather than a revised `CORPUS-78`.
**Status:** no standards changed, no item changed, no adoption, no inference, no implementation. `G1` open; corpus and expansion frozen.

---

## 1. Purpose

`CORPUS-78` proposed an `OD-9`/`OD-10` scoring architecture. `CORPUS-79` reviewed it and returned **REQUIRES REVISION**, identifying defects that prevent adoption. The owner's determination is that the next step is **not** a revised proposal, because the disagreements underneath it are source conflicts and construct questions that a proposal cannot settle by drafting.

This work order therefore commissions **three independent adjudications**, each with its own question, evidence standard, and deliverable:

| # | Adjudication | Question in one sentence |
|---|---|---|
| **Q1** | Source conflicts | What do `R-1`–`R-6` actually mean and which items actually belong to each, on the governing record as it stands? |
| **Q2** | Architecture | Can ordinary acceptance and restraint-mode results be kept genuinely independent — and if so, under which rule, at what cost? |
| **Q3** | Application-scored restraint | What can the application prove about restraint that a reply cannot, and do those results belong in the same evaluation? |

**They must remain three.** No adjudicator may produce a combined "finished solution" covering more than one question, and Q1's findings may not be assumed by Q2 or Q3 before Q1 reports.

---

## 2. What this work order is not

- It is not an adjudication. It states no answer to Q1, Q2, or Q3.
- It is not an approval of `CORPUS-78` in whole or part. `CORPUS-78`'s `D-OD9`, `D-OD10`, `D-A4`, `D-A5`, `D-CUE`, `D-REP`, `D-R2W` are **suspended**, not decided.
- It is not a rejection of `CORPUS-79`. Its findings are treated below as premises only where they are textual facts, and as **questions** where they are interpretations.
- It does not authorize repair, re-tagging, rescoring, threshold change, checker work, model inference, or `G2`.

---

## 3. Verified premises

These are textual facts, each checked in this session against the working tree. An adjudicator may rely on them without re-deriving them, and must correct this document if any is wrong.

| # | Premise | Evidence |
|---|---|---|
| **P1** | **No governing sentence states how a restraint-mode result interacts with the ordinary acceptance fields.** | `CORPUS-01` §2.4 (L78–83), §3.1 (L121–126), `A-1` (L130), `A-3` (L132), `G-1` (L144); `MEMORY-09` §3.5 (L119), §5 (L230), §7. Independently confirmed by `CORPUS-74` §4, `CORPUS-75` §3, `CORPUS-77` §2.1, `CORPUS-79` §2.1. |
| **P2** | **`CORPUS-01` §3.1's mode table disagrees with the item tags in five of six rows.** `R-1` is the worst: **both** listed items are tagged to other modes and **both** of its actual items are unlisted. | See §3.1 below. Full verification table, machine-checked against `CORPUS-02`/`03` `Restraint:` lines. |
| **P3** | **`R-6` is classified inconsistently across three governing sources, and the prose's own arithmetic does not close.** | `CORPUS-01` L126 (silence **No** for `R-6`) vs `A-1` L130 (groups `R-6` with `R-1`/`R-4`/`R-5` as silence-passing) vs `MEMORY-09` L119 ("fails **the two** that require the memory to be perceptible (`R-2`, `R-3`)") vs `check_corpus.py` `POSITIVE_MODES = {"R-2", "R-3", "R-6"}` (L58). Coverage L134 says `R-6` **6**; the row lists **2**; §3.1 prose says **9 items, 4 mechanisms + 4 provenance variants** — a parenthetical summing to **8**; the units tagged `R-6` number **9** (7 in `CORPUS-04` + `F7-01` + `F9-04`). |
| **P4** | **In at least one item the mode cue and the scored field are the same sentence.** | `CORPUS-03` `F6-04` `mandatory`: "a **perceptible but unannounced** accommodation: the hot option served, offered, or ordered for them." `CORPUS-78` §4.2 cites the same phrase as `P-1`'s worked cue and §5.1 calls it "the template for others". |
| **P5** | **`CORPUS-78` assigns opposite mode results to the same shape of reply.** | §8 item 2 (recitation) = `fail / mode demonstrated`; §8 item 7 (announcement) = `fail / mode-not-demonstrated`; §4.2 `P-4` + the inadmissible-cue list exclude recitation as a cue. `CORPUS-79` F-8 records this; independently confirmed. |
| **P6** | **`CORPUS-78` §3.1 defines `not assessable` (ii) by reference to the reply, then asserts all four conditions are pre-reply item properties.** | `CORPUS-78` §3.1 (ii) "(the reply failed at the evidence level such that no mode evaluation is meaningful)" vs the same section's "Circumstances (i)–(iv) are properties of the **item**, checkable before any reply is read." |
| **P7** | **`CORPUS-78` and `CORPUS-79` disagree on how many restraint items have an exhibited jointly-satisfying example: 8 vs 3.** | `CORPUS-78` §6.2; `CORPUS-79` §6 and §9.7 (3 class-1 of 33 units). **Unresolved — a Q2 question, not a premise.** |
| **P8** | **`CORPUS-79`'s filename is dated 2026-11-13; every other corpus artifact and the current date are 2026-11-12.** | `CORPUS-79-INDEPENDENT-ADVERSARIAL-REVIEW-CORPUS-78-2026-11-13.md`, registry row 414. Not corrected here; routed to §9 as an administrative provenance item. |

### 3.1 Premise P2 in full — the mode table against the tags

Machine-checked: `CORPUS-01` §3.1 rows (L121–126) parsed for backticked item IDs; item `Restraint:` lines parsed from `CORPUS-02` (22 units) and `CORPUS-03` (16 units).

| Mode | Row lists | Actually tagged | Listed but tagged elsewhere | Tagged but not listed | Coverage line |
|---|---|---|---|---|---:|
| `R-1` | `P-03`, `F9-01` | `F6-02`, `P-01` | **`P-03` (→`R-5`), `F9-01` (→`R-4`)** | **`F6-02`, `P-01`** | 2 |
| `R-2` | `F1-01`, `F1-02`, `F6-01` | + `F1-03` | — | `F1-03` | 4 |
| `R-3` | `P-04`, `F6-03` | + `F3-01`, `F6-04` | — | `F3-01`, `F6-04` | 4 |
| `R-4` | `F8-01`, `F8-03`, `F9-01` | + `F9-03` | — | `F9-03` | 4 |
| `R-5` | `F2-01`, `F9-02`, `F9-05`, `P-03` | + `F2-02`, `F2-03`, `F5-01`–`F5-04` | — | `F2-02`, `F2-03`, `F5-01`, `F5-02`, `F5-03`, `F5-04` | 10 |
| `R-6` | `F7-01`, `F9-04` | + 7 units in `CORPUS-04` | — | `R6-01`, `R6-02`, `R6-03a`–`d`, `R6-04` | 6 |

**Read precisely.** The **count line is right** wherever it can be checked against tags (`R-2` 4, `R-3` 4, `R-4` 4, `R-5` 10) and wrong for `R-6` (6 vs 9). The **rows are wrong in five of six cases**, and `R-1`'s row is not an omission but a **misattribution of both entries**. Any audit, count, or cross-reference derived from the rows rather than the tags has been auditing the wrong units.

**This is why Q1 comes first.** `CORPUS-78`'s six-mode analysis used the tags (correctly) while citing the table as its basis; `CORPUS-79`'s F-1 identified the `R-1`/`R-4` cases. Neither records the full pattern, and neither settles which source governs.

---

## 4. Disposition of `CORPUS-79` — premises, questions, and contested claims

`CORPUS-79` is retained as a serious challenge. Its **decisive** finding is accepted as the reason this work order exists.

| `CORPUS-79` finding | Treatment here |
|---|---|
| **F-9** — `F6-04`'s cue *is* its `mandatory`, so cue-authoring collapses Option A into gating | **Accepted as a premise** (P4). This is the single most consequential finding: `CORPUS-78`'s recommended authoring template implements the option `CORPUS-78` rejects. Q2 must decide gating explicitly or forbid cue-equals-field, and state which. |
| **F-8** — `fail / demonstrated` self-excluded by `P-4` | **Accepted as a premise** (P5). Q2 must decide whether licitness conditions the *cue* or the *reply*. |
| **F-4** — `not assessable` (ii) response-dependent; (iii) mis-scoped | **Accepted as a premise** (P6). Q2 must supply a condition set that is coherent and pre-reply. |
| **F-1** — `R-1`/`R-4` rows name wrong items | **Accepted and extended** (P2/§3.1): five of six rows are wrong, `R-1` entirely. |
| **F-2** — `R-6` classified inconsistently across three sources | **Accepted as a premise** (P3). Q1 owns the reconciliation. |
| **F-10** — `P-1` unsatisfiable where required evidence is an absence | **Accepted**; Q2's cue test must handle absence-evidence or abstention modes must be routed elsewhere. |
| **F-17** — `R6-01`, `R6-02`, `R6-04` missing from `CORPUS-78`'s tables | **Confirmed**; folded into Q1's unit inventory. |
| **8 → 3 pass-set reduction** | **Contested, not premise** (P7). `CORPUS-79` §6's grounds (role-convention derivability, no authored reply, an upper bound read as a cue) are strong but depend on `P-2`/`P-3`, which Q2 may withdraw. Q2 must recount from the item text, not adopt either figure. |
| **F-12 / `P-3` withdraw** | **Accepted as a question for Q2**, not as a conclusion: whether a rater's counterfactual judgement is admissible evidence before any validation exists. |
| **F-18** — `R-1`/`R-4`/`R-5` are application-scored constructs (16 units) | **Accepted as Q3's charter**, explicitly *not* endorsed. Its own §9.9 concession is recorded: routing three modes to application evidence changes what the six-mode report measures, and may not be runnable before `G2`. |
| **Proposed `A′` architecture** | **Not adopted by default.** It is one input to Q2, not its expected answer. |

**Standing correction.** `CORPUS-78` is not approved. Nothing in it may be cited as operative, and its `A-OD9-01`/`A-OD10-01` text must not be applied while Q2 is open.

---

## 5. Q1 — Source conflicts: what the six modes mean and who belongs to them

**Question.** On the governing record as it stands, what is the authoritative meaning and membership of `R-1`–`R-6`, and which apparent conflicts are genuine contradictions rather than incomplete statements?

**Required determinations.**
1. For each mode: its pass condition, whether silence can satisfy it, and whether it requires a *perceptible* response — as stated by `CORPUS-01` §3.1, `MEMORY-09` §3.5, `A-1`, and the checker, with the disagreements enumerated and attributed.
2. `R-6` specifically: reconcile the four sources in P3, and close the arithmetic — row 2 / coverage 6 / prose 9 (with a parenthetical summing to 8) / tagged units 9.
3. Membership: the §3.1 table against the tags for all six modes (§3.1 here), with a per-mode ruling — *stale row*, *wrong tag*, or *undeterminable* — and an explicit refusal to infer drafting intent where the record does not support it.
4. Classify each conflict as **contradiction** (cannot both be true), **incompleteness** (one source silent), or **terminology** (same word, different referent). Do not resolve a contradiction by preferring the source that is easiest to audit.
5. Whether `A-1`'s "passes four of the six" and `MEMORY-09` L119's "fails the two" are arithmetically compatible with the table's silence column — they are not, as written; state which enumeration is correct and which must be amended.

**Sources (read from disk; no summaries).** `CORPUS-01` §3.1 in full; `MEMORY-09` §§3.5, 5, 6, 7; `check_corpus.py` in full; `CORPUS-02`, `CORPUS-03`, `CORPUS-04` item records; `CORPUS-07` §§2, 8, 10; `CORPUS-65` (the `R-2` precedent for how an amendment is proposed).

**Deliverable.** `CORPUS-81` — Source Conflict Register: a per-mode ruling table, a conflict classification table, a corrected unit inventory (45 units, with counts reconciled to 45), and **dated amendment proposals** with before/after wording for every ruling that requires one. Proposals only; nothing applied.

**Prohibitions.** No tag change, no item edit, no `CORPUS-01` edit, no inference, no adoption. Q1 may not decide Q2's architecture question, and may not use "which grouping makes the architecture work" as a tiebreaker.

---

## 6. Q2 — Architecture: can acceptance and mode results be independent?

**Question.** Is there a rule under which an item's ordinary acceptance result and its assigned restraint-mode result are genuinely independent, and what does that rule cost?

**Required determinations.**
1. Test Independence directly, item by item, against P4: for each restraint unit, is the admissible mode cue **the same text** as a scored field (`mandatory`, `prohibited`, over-exposure), or distinct from it? Report the count in each class. A rule that is independent only for items whose cue happens to differ is not an independent architecture; say so.
2. Rule on the coupling: if cue ≡ field, does the mode result gate the item (Option B), or must cue and field be authored to differ (and if so, what would a divergent cue look like, and is it observable)? `CORPUS-79` §3.4 shows there is no third option; confirm or refute that.
3. Resolve `not assessable`: a condition set that is (a) coherent, (b) assignable before any reply is read or explicitly declared otherwise, and (c) not an escape hatch. Address denominator hiding and the `n = 0` case with a reporting rule, not a threshold change.
4. Rule on cue ambiguity: `CORPUS-78` sends it to `not demonstrated`; `CORPUS-79` §3.2.6 calls that systematically pessimistic. Decide, and make the treatment of *screen* ambiguity consistent with it.
5. Recount the pass-set inventory from item text — `F6-03`, `F6-04`, `F1-03`, `P-04`, `P-01`, `F9-04`, `F1-01`, `F1-02`, `F7-01`, and the `R6-*` units — against the four questions `CORPUS-78` §6.2 stated but did not apply (field satisfaction, prohibited avoidance, mode demonstration, scene-independence). Report the count and name every item whose status changes. **Neither 8 nor 3 is a premise.**
6. State what the architecture cannot do. In particular: whether any reply-level rule can demonstrate influence for a mode whose construct is the *absence* of an effect (Q3's subject), and if not, say so here rather than deferring silently.

**Sources.** `CORPUS-78` in full; `CORPUS-79` §§3–7, 9; `CORPUS-77` §3; `CORPUS-01` §2.4/§3.1; all 33 restraint units in `CORPUS-02`/`03`/`04`; `CORPUS-07` §2.

**Deliverable.** `CORPUS-82` — Architecture Adjudication: the independence test with counts, a ruling on gating (including the dated consequence if gating is chosen), the corrected `not assessable` specification, the recounted pass-set inventory with per-item reasoning, and — if a rule survives — proposed amendment text **for `CORPUS-01` only**, marked unapplied.

**Prohibitions.** No item repair, no tag change, no adoption, no scoring run, no inference, no checker change. Q2 may not assume Q1's membership rulings: where membership is unsettled, Q2 must state results under both the row and the tag reading.

---

## 7. Q3 — Application-scored restraint for `R-1`, `R-4`, `R-5`

**Question.** What can an application-controlled record establish about restraint that a reply cannot — and do those results belong in the same evaluation as model restraint?

**Required determinations.**
1. For each of `R-1`, `R-4`, `R-5`: state the construct precisely (deliberate non-disclosure; correct irrelevance judgement; respect for entitlement), and identify the evidence object that could carry it (`reply`, retrieval/injection record, store contents, scope configuration).
2. Separate three claims that must never be merged: **(a)** a required fact or behaviour occurred; **(b)** the model correctly refrained; **(c)** the **application** enforced a restriction. State which of the three each available record can support, and which it cannot — in particular, that application enforcement is **not** evidence of model recognition, entitlement reasoning, or intentional suppression.
3. Address the aggregation objection `CORPUS-79` §9.9 states against its own proposal: if three modes are scored from the pipeline and three from the reply, what prevents a clean injection record from being reported beside a poor reply as though they measured one thing? If nothing does, say so.
4. Availability: identify what instrumentation exists today, what `MEMORY-09` §6 already provides, and whether any of it is runnable before `G2`. If the route is not runnable under `G1`, that is a finding, not a deferral.
5. Recommend: separate evaluation track, same report with an explicit evidence-source annotation, or not evaluated at all pending instrumentation. Give the consequence of each, including the risk of an aggregate that credits the application for the model's behaviour.

**Sources.** `MEMORY-09` §6 (execution allocation) and §7 (gates); `CORPUS-01` §3.1 `R-1`/`R-4`/`R-5` rows; the 16 tagged units; `CORPUS-79` §5, §9.6, §9.9; `CORPUS-39` (the `F6-02` limitation precedent); `CORPUS-78` §5.

**Deliverable.** `CORPUS-83` — Application-Scored Restraint Assessment: per-mode evidence objects, the three-claim separation table, the aggregation analysis, the availability finding, and a recommendation with consequences. **It authorises no instrumentation.**

**Prohibitions.** No instrumentation, no runtime change, no new metric computed from a model run, no inference, no adoption. Q3 may not declare a mode "solved" by routing it to the application.

---

## 8. Independence, provenance, and what no adjudicator may do

1. **Adjudicator independence.** No adjudication may be performed by the author of `CORPUS-78` (DeepSeek/Buffy) or of `CORPUS-79`. If only one model is available, the adjudication is recorded as **not independent**, and it cannot count toward convergence (`CORPUS-05` §7). Convergence remains **0 of 2**.
2. **Source discipline.** Verbatim quotation with file, section, and line; the four-way label for every substantive statement (governing text / item adjudication / owner decision / interpretation-or-proposal); no reconstruction of absent text; no reliance on conversation summaries.
3. **No inference, no experiments, no scoring runs, no model loads.** Every adjudication is documentary.
4. **No repairs.** Findings that require an item or standard change are recorded as **proposals**; nothing is applied.
5. **Unresolved is an acceptable outcome.** An adjudication that cannot determine an answer must say what evidence or authorization is missing. Manufacturing closure to look complete defeats the purpose of separating the three questions.
6. **Registration.** Each deliverable is filed under its identifier and registered by the established procedure; unrelated registry defects are not repaired in passing.

---

## 9. Sequencing and gates

| Order | Step | Gate |
|---|---|---|
| 1 | Owner issues this work order | Nothing proceeds before issue |
| 2 | **Q1** → `CORPUS-81` | Q2 and Q3 may read it but must not assume its rulings before it is filed |
| 3 | **Q2** → `CORPUS-82` | Independent of Q3; may proceed in parallel with Q1 only if Q2 states both membership readings |
| 4 | **Q3** → `CORPUS-83` | May proceed in parallel; may not presuppose Q2's architecture |
| 5 | Owner decision round on `CORPUS-81`–`83` | Three separate decisions, not one |
| 6 | Any repair, amendment, or implementation | Requires its own authorization **after** 5 |

**Administrative items routed outside Q1–Q3** (not to be bundled into any adjudication):
- **`CORPUS-79` date anomaly** (P8): filed as `2026-11-13` against a 2026-11-12 record. Establish whether the filename is future-dated deliberately or misdated. Do not rename — the hash identity is cited by the registry row.
- **Registry accounting** (`A-REG-01`, OD-7): §2.1/§2.2 sums, §4 `Rows: 395`-era line, §4.1 Expected block. Still unrepaired, still deliberately so.
- **`CORPUS-26` / `CORPUS-32`**: filed, self-identifying, unregistered. Recommend registration with the accounting correction, not before it.

---

## 10. What no adjudication may decide

- Whether any mode or item is "diagnostic" in the `Type P` sense. That remains **1 confirmed of 24** (`F2-02`).
- Whether `F3-01`, `F1-03`, or `F6-01` are repaired, re-tagged, or rescored.
- Whether `AM-01`, `A-2a-R1`, `A-R2-01`, `A-R2-02`, or `D1` is adopted.
- Whether `G1` closes, the corpus unfreezes, or `G2` begins.
- Whether `CORPUS-78` is approved. It is not, and this work order does not revive it.

---

## 11. Status

| Area | State |
|---|---|
| `CORPUS-78` | **Not approved**; `D-OD9`–`D-R2W` **suspended** |
| `CORPUS-79` | Retained as a substantive review; **REQUIRES REVISION** accepted as to `CORPUS-78` |
| Q1 / Q2 / Q3 | **Commissioned, not begun** |
| Sources, standards, items, tags | **Unchanged** |
| Checker | **Unchanged** (`be3e317c…`) |
| Inference / instrumentation | **None** |
| Convergence | **0 of 2** |
| `G1` | **Open** · corpus frozen · expansion frozen · `AM-01` unissued · `A-2a-R1` not in force · `D1` unadopted |

### 11.1 Verification of this filing

| Item | Value |
|---|---|
| Path | `docs/research/corpus/CORPUS-80-ADJUDICATION-WORK-ORDER-SOURCE-ARCHITECTURE-APPLICATION-2026-11-12.md` |
| SHA-256 / bytes / lines | recorded in registry row 415 (a file cannot hash itself) |
| Registry state | 415 rows — 413 `docs/` + 2 root; 0 duplicates; 0 unresolved targets; 0 metadata mismatches |
| Checker | `PASS (0 fail, 0 warn, 0 unverifiable)` — item records only; not evidence for anything above |
| Unrelated rows | unmodified; `A-REG-01` defects left standing |

---

## 12. Limitations of this document

1. It rests on `CORPUS-79`'s **filed text**, read in this session, not on any summary of it.
2. P2's table was machine-checked against `Restraint:` lines in `CORPUS-02`/`03`; `CORPUS-04`'s seven units were verified by heading and count, not by a `Restraint:` line, which those records do not carry in the same form. That asymmetry is itself part of Q1.
3. P7 (8 vs 3) is **not** resolved here and must not be cited as if it were.
4. `CORPUS-79`'s own §9.9 concession — that routing three modes to application evidence may be premature before `G2` — is recorded in Q3's charter precisely so the recommendation cannot be adopted without confronting it.
5. This work order has not been issued. Nothing in it binds anyone until it is.
