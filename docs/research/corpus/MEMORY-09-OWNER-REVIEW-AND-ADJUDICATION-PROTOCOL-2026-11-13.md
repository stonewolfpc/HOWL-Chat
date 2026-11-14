# `MEMORY-09` — Owner Review and Adjudication Protocol

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-ORAP-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **PROTOCOL PROPOSAL — no record has been reviewed under it.**
**Authority:** `MEMORY-09-WO-ARCH-06` §4, §10.

---

## 1. Principle

The owner is the **final interpretive authority**. AI evaluation is a research assistant: it finds, compares, and proposes; it never concludes. **Both** directions of overreach are forbidden — an AI's unfavourable finding may not silently invalidate something the owner values, and an owner's approval may not silently erase a technical defect.

## 2. Presentation order (non-negotiable)

| Layer | Content | Length |
|---|---|---|
| **1 · Verdict sentence** | One plain sentence: what happened and whether it supports the hypothesis | 1 line |
| **2 · Three observations** | The three things that decide it, in ordinary language | ≤3 bullets |
| **3 · Evidence** | Values, counts, diffs, paths | as needed |
| **4 · Raw** | Full transcript, state dumps, hashes | on request |

**A record that starts with numbers has failed this protocol.** The owner should be able to act after layer 2; layers 3–4 exist so the act can be checked later.

## 3. Owner action menu

| Action | Meaning | Recorded as |
|---|---|---|
| **Accept** | Agrees with the stated interpretation | `[OWNER] accepted` |
| **Dispute** | Disagrees with the interpretation, not the observations | `[OWNER] disputed` + reason |
| **Request investigation** | Wants a narrower test before deciding | `[OWNER] follow-up` + question |
| **Override** | Sets a different disposition than the evidence suggests | `[OWNER] override` + rationale |
| **Suspend** | Insufficient evidence to decide | `[OWNER] suspended` |

**Dispute and override are different acts.** Disputing an *interpretation* leaves the observations standing. Overriding a *disposition* leaves both the observations and the original disposition standing, and adds the decision beside them.

## 4. Override record — required fields

| Field | Why |
|---|---|
| Original result, verbatim | The evidence must survive |
| Original interpretation | So the disagreement is legible later |
| Owner's disposition | The decision actually taken |
| Rationale in the owner's words | Why the evidence was outweighed |
| What would change the owner's mind | Keeps the override falsifiable |
| Date and scope | Which tests it applies to |

**An override never edits blocks B or C.** It appends block D. The record keeps both readings permanently.

## 5. Preserving disagreement

When AI interpretation and owner judgement differ, **both are retained with their evidence**, and the disposition names which one governs and why. A later reviewer must be able to see that a defect was known and consciously accepted — or that an owner preference was overruled by mechanical fact.

## 6. Statistics, when they appear at all

- **α is not used for run-to-run variation.** α measures agreement between *raters judging the same units*; spread across repeated runs of one model is a different quantity. Reporting α for repeated runs would be a category error and is prohibited here.
- When a rate is reported ("7 of 10 runs"), state the **numerator, denominator, and what one run is**.
- No p-values, no confidence intervals, no thresholds unless the owner has adopted them. A bare percentage with no denominator is not a measurement.

## 7. Targeted investigation

When the owner flags a concern, the default is **not** a restart. The procedure:

1. Restate the concern as a question the record can answer.
2. Identify the smallest test that would answer it (often one condition, 3 runs).
3. Run only that, under the frozen baseline.
4. Append as a linked record; the original is untouched.

**A full re-run is justified only when a frozen condition changed** (model, quantisation, sampler, architecture), not when a result was unwelcome.

## 8. Escalation triggers — when the owner should not decide alone

| Situation | Why |
|---|---|
| Two mechanical facts contradict each other | Something is wrong with the record, not with the judgement |
| State dumps disagree with the reply's claims | A possible architecture defect; needs investigation, not adjudication |
| A result depends on a fact the model could have guessed (C3 not excluded) | Attribution is unresolved; the record must say so |
| The owner's rationale would erase a `[C4]` incorrect-despite-correct-state finding | The most important defect class must not be adjudicated away |

## 9. Illustrative record shape

> **T-07 · inventory removal — supports hypothesis.**
> The character stopped offering the lantern once it was removed from the inventory. Three runs, no phantom item, state dump matched at every turn. *Two earlier runs under a higher repeat penalty showed the lantern reappearing — recorded as a separate condition, not merged.*

Then evidence, then owner block. Note that the last sentence reports a **defect**, not a success — records must be allowed to contain both.

## 10. Anti-patterns

| Anti-pattern | Why it fails |
|---|---|
| Score-only summary ("7/10, pass") | Discards what the owner needs to judge |
| AI verdict in block D | Authority laundering |
| Silent deletion of a failed run | Falsifies the record |
| Restarting after an unwelcome result | Converts evaluation into confirmation |
| Percentages without denominators | Unfalsifiable |
| α applied to repeat runs | Wrong quantity |
| Override without "what would change my mind" | Unfalsifiable decision |

## 11. Claim ceiling

This protocol governs *how* results are presented and adjudicated. It establishes no result, and it cannot make a poorly designed test legible — if the matrix row was wrong, the protocol will faithfully present a useless finding.
