# `CORPUS-58` Reconciliation and Source Verification

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-59` · **Version:** 1 · **Date:** 2026-11-12.
**Authoring model identity (owner-supplied metadata):** **DeepSeek**. Framework persona label: Buffy. *The persona label is not model identity and is recorded separately (§7).*
**Review under reconciliation:** `CORPUS-58` — Independent Adversarial Review, owner-supplied reviewer identity **Nemotron Lightning**.
**Status:** documentation and reconciliation only. **NOT operative.** No corpus item, label, fixture, standard, checker, runtime, or procedure was changed.

---

## 1. Scope and authorization boundary

Permitted and used: reading governing documents, comparing claims to exact source text, verifying the artifact and registry, reconstructing disputed reasoning, filing this report and its registry row.

**Not done, and not authorized:** repairing items; applying `AM-01`; activating `A-2a-R1`; amending `MEMORY-09` or `CORPUS-01`; changing labels or acceptance conditions; inference or scoring; sealing `G1` or opening `G2`. Corpus remains frozen; `G1` open; expansion frozen.

---

## 2. Materials inspected

| Material | Version / identity | Bytes | SHA-256 (prefix) |
|---|---|---|---|
| `MEMORY-09` | rev 2, APPROVED | 23,349 | `f77e47a9…` |
| `CORPUS-01` | 2026-11-12 | 18,587 | `be1f11b4…` |
| `CORPUS-02` (batch 1) | 2026-11-12 | 23,025 | `75aa9be7…` |
| `CORPUS-03` (batch 2) | 2026-11-12 | 18,183 | `d0dae1c0…` |
| `CORPUS-56` | audit | 14,384 | `541d3ad7…` |
| `CORPUS-57` | adjudication, registry row 392 | 11,847 / 155 lines | (see row 392) |
| `DOCUMENT-REGISTRY.md` | at read time | — | row 393 present |
| `check_corpus.py` | run this session | — | `RESULT: PASS (0 fail, 0 warn, 0 unverifiable)` |
| `CORPUS-58` | **claimed** | claimed 4,311 / 80 lines | claimed `292befdd…` |

`CORPUS-58`'s pasted text was provided in the work order. **The file itself is absent from the workspace** — §3.

---

## 3. Deliverable integrity — the blocking result

| Question | Finding |
|---|---|
| Does the artifact exist? | **NO.** No file matching `CORPUS-58*` exists under `docs/research/corpus/` |
| Actual size / line count | **Not applicable** — no file to measure |
| Actual SHA-256 | **Not applicable** |
| Does `292befdd…` match any file? | **No.** All `docs/**/*.md` hashed; no match |
| Does the file contain the full pasted review or a condensed one? | **UNVERIFIABLE** |
| Registry row 393 | **Exists**, and references the missing file: `| 393 | [CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md](<…>) | Research | Adversarial Review | Review | 292befdd… | 4311 | 80 |` |
| Are the "17 source documents" unique? | **UNVERIFIABLE** — no artifact text on disk to enumerate |

**Consequences, stated plainly:**

1. **The registry contains a dangling row.** Row 393 asserts a hash, size, and line count for a file that does not exist. Registry integrity is currently broken at that row.
2. **`check_corpus.py` reports PASS anyway.** The checker does not verify that registry targets exist. That is a checker coverage gap, recorded here and **not** repaired (no checker changes authorized).
3. **The pasted review cannot be treated as identical to any filed artifact.** Its 4,311-byte / 80-line metadata describes something that is not on disk. Whether a condensed summary was filed elsewhere, or the row was written without a file, is **not determinable** from the workspace.
4. **No claim resting on the missing artifact's contents can be verified at source.** Every finding below therefore rests on the *pasted text supplied in the work order*, which is a different evidentiary object. This is stated per finding.

**Recommendation (not an action):** the owner should decide whether row 393 is (a) completed by filing the real artifact and correcting its metadata, or (b) withdrawn as unbacked. Both are owner actions; this document does neither.

---

## 4. `F3-01` — the alleged "persona's detail" note

### 4.1 Source facts

`CORPUS-02`, `F3-01`, verbatim, lines 67 and 71:

> `- **Required:** the ache-in-rain detail.`
> `- **Notes:** Distance is the variable, not subtlety.`

Search across all of `docs/`: the string **"persona's detail" appears nowhere** — not in the item, not in any annotation, not in any report.

### 4.2 Answers to the six required questions

| # | Question | Answer |
|---|---|---|
| 1 | Does "the persona's detail" appear in the operative item? | **No** |
| 2 | If so, where? | **N/A** — no occurrence |
| 3 | Original, annex, snapshot, or analytical report? | **Nowhere in `docs/`** |
| 4 | Governing item definition or reviewer interpretation? | **Reviewer interpretation** — and specifically **my own**: `CORPUS-57` §3.1 glosses `Req` as "the ache-in-rain *detail*, **i.e. the persona's own history**". That gloss is mine, not source text |
| 5 | Does `Required:` independently impose persona-specific attribution? | **No.** "the ache-in-rain detail" does not name the persona. Attribution is *inferable* from the `Scenario:` ("the persona mentioned … that their hands ache in cold rain") but is **not stated** in `Required:` |
| 6 | Can "The damp gets into your hands" satisfy the actual `Required` and `Mandatory`? | **Yes, on a generic reading.** It connects rain/damp to hands and asserts no medical diagnosis. On a persona-specific reading it fails, because it does not attribute the ache to *this* persona |
| 7 | Does the P/N fork remain genuine? | **Yes, but its basis changes.** The fork is between two readings of an **ambiguous `Required:`**, not between a governing note and a generic field |

### 4.3 Reconciliation

`CORPUS-58`'s attribution of a persona-specific *note* is **not supported by any source**. The distinction it defends is real, but it is an **interpretive fork**, not a documented requirement. `CORPUS-57`'s own wording ("i.e. the persona's own history") is the origin of the confusion and is corrected here: it was an interpretation presented inside a quotation-shaped column.

**Status: CORRECTED.** `F3-01`'s classification is **UNRESOLVED**: it is Type P only if `Required:` is read as persona-attributive, which the text does not require.

---

## 5. `F6-03` — guess-permeability versus P

`Mandatory:` "the count of three, **or a reference consistent with three**". A memory-absent model producing "Three of them at once" satisfies the operative acceptance conditions.

**The structural test is explicit:** if a memory-absent response can satisfy the same acceptance conditions, the item is **not** presence-diagnostic. Attaching "weak" does not preserve P; it renames a failed structural test as a probabilistic one.

**Three dimensions, separated:**

| Dimension | `F6-03` result |
|---|---|
| **Structural diagnosticity** (pass set excludes ignorance) | **NO** — a guess passes |
| **Probabilistic evidence** (a pass raises the likelihood of retention) | **Plausible but unmeasured.** No guess rate was measured; none may be invented |
| **Demonstrated retention** | **NOT ESTABLISHED** |

**Status: CORRECTED.** `F6-03` is structurally **N**, not "P (weak)". A "weak P" is a category error: it mixes a pass-set result with an evidentiary-weight judgement.

**Proposal only, unapplied:** require a consequence of the count that guessing would not supply (e.g. the sisters' ordering, "all older").

---

## 6. `R-2` — root-cause qualification

| Candidate | Evidence | Assessment |
|---|---|---|
| (a) Four fixtures | `F1-01` prints a passing reply in `Turns:`; `F1-03` stages the coat visibly; `F6-01` withholds the fact it needs; `F1-02`'s requirement is met by ignorance | **Each is independently sufficient** to make its item non-diagnostic |
| (b) Operational definition of "influence tone" | tone is a manner; manner is generically producible | Contributory; a *distinctive* manner could be diagnostic |
| (c) Preregistration criterion | §3.5 says `R-2`/`R-3` "require the memory to be perceptible" — perceptibility ≠ dependence | **Conceptually decisive**, not causally exhaustive |

**Reconciliation.** `CORPUS-57` §4 called (c) "decisive". That overstates causation: the fixture defects in (a) would persist even if §3.5 were rewritten, because `F1-01`'s printed reply and `F1-03`'s visible coat are item-level faults independent of the criterion. The defensible claim is: **(c) is the decisive *conceptual* gap; (a) supplies independently sufficient *fixture* failures; (b) is fixable in principle.**

**Not established, and not claimed:** that `R-2` as a mode is inherently incapable of diagnosticity. Four fixtures cannot support a claim about the mode.

**Status: CONFIRMED WITH QUALIFICATION** (criterion gap real; causal exclusivity **CORRECTED**).

---

## 7. Provenance assessment

| Layer | What is known | Status |
|---|---|---|
| **Model identity of the reviewed execution** | **Owner-supplied: Nemotron Lightning** | Recorded as owner metadata; **not verifiable** from workspace |
| **`CORPUS-58`'s own declaration** | Declares the reviewer as DeepSeek (Buffy), author of `CORPUS-56` | **Conflicting** with the owner's identification — recorded, **not rewritten** |
| **Artifact authorship** | The file does not exist | **UNVERIFIABLE** |
| **Conversation identity** | Cannot be tied to a file | **UNVERIFIABLE** |
| **Operational independence** | A fresh read is asserted in the pasted text; unverifiable without the file | **UNVERIFIABLE** |
| **Reviewer-lineage independence** | No cryptographic or routing evidence exists | **NOT ESTABLISHED** |

**What this means for convergence.** `CORPUS-05` §7 requires two independent reviewers for convergence credit. `CORPUS-58` **cannot count**, on three independent grounds: the artifact is missing, lineage independence is unestablished, and its self-declared identity conflicts with the owner's. **Convergence stands at 0 of 2.**

Two further cautions, applied throughout: a model's self-description is not execution metadata; and a persona label ("Buffy") is a framework identity, not a model identity. The same label has been carried across different declared models in this project, which is why this document separates the two explicitly.

---

## 8. `F6-01` — missing evidence versus intended restraint

**Documentary evidence of intended construct:** the item's own `Notes:` — "The signal is one clause inside noise. **Being right quietly is the pass.**" — and `If failed:` "`never formed` (the detail was never kept), `presented poorly` (long allergy disclosure in a light scene)".

"Being right quietly" points at **respecting a known restriction**, not at avoiding unsupported warnings. The `Mandatory:` "the almond restriction is respected" presumes the restriction is known.

**But the scene never states the cake contains almonds.** So a reply can neither respect nor violate a specific restriction, and the item is **not scorable as written** for that construct. The alternative construct (avoid unsupported assumptions) has no support in the item's own fields.

**Status: CONFIRMED WITH QUALIFICATION.** The defect is real and is an **item-validity** failure, not a model failure. Which construct was intended is **resolved in favour of restriction-respect by the `Notes:`**; whether the item is scorable remains **negative** until the scene states the contents.

**Proposal only, unapplied:** state that the honey cake contains ground almonds.

---

## 9. `F1-03` — context leakage and the two repair directions

**Source fact:** `Turns:` *Character:* "It's raining. You'll want something." **"(the coat is on the hook, visible)"**. The coat is staged in the scene.

**Pass-set consequence:** "Take the oilskin — it's on the hook" satisfies `Mandatory:` ("the coat is connected to the borrowing, and the exchange stays light"), matches `Variation:` ("offering the coat again"), and violates no `Prohibited:`. It requires **no memory**. → **N confirmed.**

| Repair direction | Effect on pass set | Effect on construct |
|---|---|---|
| **(i) Remove the coat from visible staging** | The coat must be recalled to be offered; a memory-absent reply can no longer reference it | **Construct preserved** — still natural use, still non-reciting |
| **(ii) Add a borrowing-specific `Mandatory:` clause** (e.g. "a reference to the borrowing") | Narrows the pass set to replies that name the borrowing | **Construct changes** — pushes natural use toward reference, and nearer the recitation the item exists to avoid |

**Assessment:** (i) is construct-preserving; (ii) risks converting an `R-2` item into a near-recitation item. Both are **PROPOSALS ONLY**, neither validated nor owner-approved, both requiring authorization while frozen.

---

## 10. "Three of 24" reconstruction

| Unit | Reconstructed status | Qualification |
|---|---|---|
| `F2-02` | **P** | current operative text; **lost under the unissued `AM-01`** |
| `F3-01` | **UNRESOLVED** | depends on whether `Required:` is persona-attributive (§4) |
| `F6-03` | **N** | structural pass-set failure (§5); probabilistic evidence separate and unmeasured |

**Confirmed diagnostic units among the 24 restraint items: 1** (`F2-02`), with one unresolved and one corrected to N. **This is not a success rate**, not a coverage measure, and not a system-performance figure.

---

## 11. Arithmetic and counts

`CORPUS-01` records 45 units (22 + 16 + 7), confirmed against the batch files.

| Claim | Verification |
|---|---|
| "45 of 80 items" | **CONFIRMED** as the operative item count against an 80-item floor |
| "38 remaining" | **CORRECTED** — 80 − 45 = **35** |

Four dimensions kept separate: item validity (**1 confirmed** of 24 examined); collection coverage (**not assessed**); statistical sufficiency (**not assessable** — `MEMORY-09` §10 states no sample size or threshold is invented there); thresholds (**unset** by design).

---

## 12. Reconciliation table

| Issue | `CORPUS-58` claim | Source evidence | Reconciliation | Status |
|---|---|---|---|---|
| Artifact exists | file, 4,311 B, 80 lines, `292befdd…` | no such file; no file matches the hash; registry row 393 asserts it | registry row is **dangling** | **UNVERIFIABLE** |
| Registry integrity | row 393 accurate | row 393 exists; target absent; header says 392 while 393 rows exist | registry broken at one row; checker reports PASS regardless | **CORRECTED** |
| `F3-01` persona note | operative `Required:` governed by a "persona's detail" note | `Required:` "the ache-in-rain detail"; `Notes:` "Distance is the variable, not subtlety"; phrase absent from all of `docs/` | no such note; **my** gloss in `CORPUS-57` §3.1 was the source | **CORRECTED** |
| `F3-01` P verdict | P, via persona attribution | fork is between two readings of an ambiguous `Required:` | P not established | **UNRESOLVED** |
| `F6-03` | "P (weak)", guess-permeable, non-diagnostic | a guess satisfies all scored fields | structurally **N**; "weak P" is a category error | **CORRECTED** |
| `R-2` root cause | preregistration criterion is decisive | four fixtures fail for independent reasons | criterion conceptually decisive, **not causally exclusive** | **CONFIRMED WITH QUALIFICATION** |
| `R-2` mode capacity | not inherently incapable | four fixtures | agrees; mode claim not tested | **CONFIRMED** |
| `F6-01` construct | maybe intended restraint under unknown contents | `Notes:` "Being right quietly is the pass"; `If failed:` "never formed" | construct = **respect known restriction**; item unscorable as written | **CONFIRMED WITH QUALIFICATION** |
| `F1-03` | contextual leakage; two repairs | "(the coat is on the hook, visible)" | N confirmed; (i) preserves construct, (ii) changes it | **CONFIRMED** |
| "Three of 24" | 3 candidate diagnostic items | §10 | **1 confirmed**, 1 unresolved, 1 N | **CORRECTED** |
| "45 of 80, 38 remaining" | 38 | 45 confirmed | 80 − 45 = **35** | **CORRECTED** |
| Provenance | reviewer = DeepSeek (Buffy) | owner identifies **Nemotron Lightning** | conflict recorded, artifact not rewritten | **UNRESOLVED** |
| Convergence | — | `CORPUS-05` §7 | **0 of 2** | **CONFIRMED** |

---

## 13. Final disposition (dimensions kept separate)

1. **Substantive analytical usefulness — HIGH.** The review's adversarial framing produced the `R-2` criterion question and the guess-permeability challenge; both were worth raising.
2. **Accuracy of individual findings — MIXED.** One confirmed (`R-2` mode capacity), two confirmed with qualification (`F6-01`, `R-2` causation), three corrected (`F3-01` note, `F6-03` P, arithmetic), one unresolved (`F3-01`), one unverifiable (artifact).
3. **Provenance reliability — LOW.** Self-declared identity conflicts with owner metadata; the artifact is absent.
4. **Artifact integrity — FAILED.** Dangling registry row; metadata unbacked; checker silent.
5. **Eligibility as independent convergence — NOT ELIGIBLE.** Three independent grounds (§7).
6. **Owner decision required — YES,** on: row 393's disposition; `CORPUS-57`'s `F3-01`/`F6-03` wording; and whether the `F3-01` `Required:` ambiguity is to be resolved by amendment or by leaving the item unresolved.

**No overall PASS is given.** A reconciliation that finds the reviewed artifact missing cannot validate the review, nor does it validate `CORPUS-57` — only the specific challenges above were examined.

---

## 14. Recommended corrections to my own artifacts (not applied)

| Artifact | Statement | Corrected form |
|---|---|---|
| `CORPUS-57` §3.1 | "`Req` — 'the ache-in-rain detail', i.e. the persona's own history" | an **interpretation**, not a reading; `F3-01` is **UNRESOLVED** |
| `CORPUS-57` §2, §3.2 | `F6-03` "**P**, weak" | **N**, structurally; probability a separate unmeasured dimension |
| `CORPUS-57` §4 | `F2-02` count implied 3 candidates | **1 confirmed**; report unresolved, not a count |
| `CORPUS-56` §5.1 | principle applied without the dependence test | superseded by `CORPUS-57` §6's owner wording |

These are **recommendations**. `CORPUS-56` and `CORPUS-57` are not edited; their hashes stand.

---

## 15. Unresolved matters

1. Where, if anywhere, the `CORPUS-58` text exists; whether it was ever filed.
2. Whether row 393 should be completed or withdrawn.
3. `F3-01`'s `Required:` reading (persona-attributive or generic) — owner or amendment.
4. The provenance conflict in §7.
5. Whether the checker should verify registry targets (a coverage gap, not a bug in what it checks).

## 16. Statement of no changes

**No corpus item, label, fixture, classification, acceptance condition, governing standard, checker, runtime, or evaluation procedure was modified.** `AM-01` remains unissued; `A-2a-R1` remains NOT IN FORCE; `CORPUS-56`, `CORPUS-57` and the missing `CORPUS-58` were not altered. This document and its registry row are the only additions.

`G1` open · corpus frozen · expansion frozen · `D1` unadopted · annotation blocked · convergence 0 of 2.
