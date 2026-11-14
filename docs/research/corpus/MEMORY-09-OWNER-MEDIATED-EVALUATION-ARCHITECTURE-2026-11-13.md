# `MEMORY-09` — Owner-Mediated Evaluation Architecture

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-OME-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **ARCHITECTURE PROPOSAL — no test executed, no standard amended.**
**Authority:** `MEMORY-09-WO-ARCH-06` §1–§4, §10.
**Executor:** single model. Nothing here claims independent review.

---

## 1. The question this architecture serves

> **Does the `MEMORY-09` architecture produce the intended behavioural improvements, independently of the model running underneath it?**

That question is *comparative* (architecture vs no architecture) and *situated* (this hardware, these models, realistic chat use). It is not "is the corpus valid" — a different question, addressed elsewhere.

## 2. Four layers, four claim ceilings

| Layer | Who acts | Produces | May claim |
|---|---|---|---|
| **L0 · Mechanical state** | Deterministic code / inspection | Inventory values, clothing list, character facts, index lookups | **Facts about stored state** |
| **L1 · Controlled run** | Runtime, one configuration | A generated reply under recorded settings | **That this reply occurred under these conditions** |
| **L2 · Interpretation** | AI evaluator or the owner | Pattern, contradiction, likely cause | **A hypothesis about what the reply shows** |
| **L3 · Adjudication** | **Owner only** | Disposition | **What the project concludes** |

**A result at L1 never becomes an L3 conclusion without passing through L2 and an owner record.** The layers exist so no layer inherits authority it did not earn.

## 3. Reconcile with `OD-M09-17` — what actually needs two raters

`OD-M09-17` (addendum L88–L118) requires a **two-rater study before promotion beyond single-source confidence** for **cue-level results** (L90). It is not a blanket requirement that every observation be double-rated. The distinction that resolves this is **what the claim asserts**:

| Activity | Two raters required? | Why |
|---|---|---|
| **Formal reproducibility study** (raters judge the same units; α reported) | **Yes** — that is exactly what L88 requires | The claim is *about agreement between judges* |
| **Controlled repeated-run engineering test** (one config, N runs, observable outcome) | **No** — but labelled single-source | The claim is *about system behaviour*, measured by observation, not by rater judgement |
| **Cross-model compatibility test** | **No** | Different models are different conditions, not raters |
| **Sampling sensitivity test** | **No** | Same, with sampling as the varied factor |
| **Owner-led qualitative review** | **No** — labelled owner judgement | L3, not a rating |
| **Mechanical state/continuity verification** | **No** | Deterministic inspection; no judgement at all |

**The line:** *judgement reproducibility* needs two raters. *Behavioural observation under recorded conditions* does not — provided it is never reported as though inter-rater agreement had been established.

**This is not a reinterpretation.** `OD-M09-17` says cue-level results may be reported "at **single-source confidence**" (L90). Engineering results sit below that ceiling by construction: they support "the system did X under Y", never "the corpus is reliably ratable".

**Owner decision requested:** confirm this reading (§3 table) as the governing split. Until confirmed, engineering results are recorded as single-source and labelled as such.

> **Revision note (2026-11-13, added after source audit).** The split above is a **proposed interpretation**, not a settled reading. The audit at `MEMORY-09-OD-M09-17-SOURCE-AUDIT-2026-11-13.md` (row 476) finds the *scope* claim supported — the decision addresses cue-level judgement — but corrects the implication that the record **permits** engineering evaluation. It does not reach it, and grants no permission: `OD-M09-17` §1 is "documentation direction only" and its effect statement authorises "no … evaluation". Engineering evaluation therefore needs a separate owner authorisation.

## 4. The evaluation record (A–D)

| Block | Fields |
|---|---|
| **A · Test identity** | Test ID/version; objective; expected behaviour; exact inputs + preceding context; model identity/version; runtime; **exact sampling values**; memory/retrieval/state configuration |
| **B · Evidence** | Full unedited response; relevant prior messages; observed application state; expected vs observed; mechanical measurements; errors/warnings |
| **C · Interpretation** | Plain-language explanation first; mechanical findings; AI interpretation (labelled); uncertainties and alternatives; supports/contradicts/does-not-address |
| **D · Owner review** | Owner assessment; agree/disagree; follow-up requested; override + rationale; final disposition |

**Rule:** D is never filled in before C, and C is never filled in before B. The record's order is the audit trail.

## 5. Presentation contract

1. **First paragraph is plain language.** No statistics, no code, no file paths.
2. **Second layer is the mechanical evidence** — counts, values, diffs.
3. **Third layer is available on request** — raw outputs, full transcripts, hashes.

A result that cannot be stated in one plain sentence has not yet been understood, only computed.

## 6. Labelling discipline

Every statement in a record carries one of four tags:

| Tag | Meaning |
|---|---|
| `[MECH]` | Reproduced by inspection or deterministic calculation |
| `[OBS]` | Observed in a recorded run |
| `[AI]` | Interpretation by a model — **never authority** |
| `[OWNER]` | Judgement by the owner — final, but recorded with its basis |

**An `[AI]` statement may never overwrite an `[OBS]` one, and an `[OWNER]` override may never delete one.**

## 7. What this architecture deliberately does not do

- It does not convert repeat runs into independent evidence (§2.7 of the order; one model run 20 times is still one model).
- It does not let a checklist substitute for a test: an unexecuted row is `NOT RUN`, never `PASS`.
- It does not let the owner's preference erase a mechanical defect — an override **records** the defect and the decision, it does not remove the finding.
- It does not use α for run-to-run variability. α measures **agreement between raters on the same units**; run-to-run spread is a different quantity (see the test matrix, T-01).

## 8. Failure modes of this architecture itself

| Risk | Mitigation |
|---|---|
| Owner becomes a rubber stamp because records are long | Presentation contract (§5): one plain sentence, then evidence |
| AI interpretation quietly becomes the finding | Tag discipline (§6); `[AI]` never appears in block D |
| Convenient configurations get chosen per test | Baseline frozen per model before tests (§472); changes recorded as conditions |
| Null results get buried | `NOT RUN` and `NO DIFFERENCE` are first-class dispositions, reported like any other |
| The architecture is judged by whether results please | The claim ceiling is fixed at L1–L2 until L3 is recorded, whatever the result |

## 9. Claim ceiling of this document

It proposes a process. **No test has been run, no model loaded, no result exists.** It does not establish that `MEMORY-09` improves behaviour, that any model is compatible, or that the owner will find the output legible — the last of which only the owner can judge once a record exists.
