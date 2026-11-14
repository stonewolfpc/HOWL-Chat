# `MEMORY-09` — Balanced Chat Baseline (Proposal)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-BCB-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **PROPOSAL — not adopted; no setting has been applied to any run.**
**Authority:** `MEMORY-09-WO-ARCH-06` §6.

---

## 1. What the baseline is

**A profile, not a number set.** "Balanced Chat Baseline" (BCB) names the *intent* — ordinary conversational operation, not tuned for coding, reasoning or creative flourish — and is realised per model family, because the publishers' recommended values genuinely differ.

**The precedent already exists in the codebase.** `app-family-defaults.go` applies *published* per-family sampling values and, for families with no publication, keeps generic defaults. Its own comment states the governing principle:

> *"Llama 3.x, Mistral-Nemo and most others publish no numbers at all, so the generic defaults are the honest choice rather than invented ones."*

**BCB adopts that principle wholesale.** It is not a new sampler policy; it is the existing one, declared as the evaluation baseline and frozen per model.

## 2. Two layers

| Layer | Applies to | Contents |
|---|---|---|
| **L-A · Common objective** | Every model | Task framing, context budget policy, stop conditions, the instruction to behave as an ordinary conversational partner. **No numeric sampler values.** |
| **L-B · Per-family published values** | Each family with a published recommendation | Temperature, top-p, top-k, min-p, penalties — **only where a publisher states them** |

**The common layer carries the *objective*; the family layer carries the *numbers*.** That is what stops BCB from being a creative preset or a reasoning preset wearing a neutral name.

## 3. Values already established in the app (evidence, not proposal)

From `app-family-defaults.go` v1, each with its cited source:

| Family | Temperature | top-p | top-k | min-p | presence | Source in code |
|---|---|---|---|---|---|---|
| qwen3 | 0.6 | 0.95 | 20 | 0 (explicit) | 0 | Qwen3 model card, thinking mode |
| gemma | 1.0 | 0.95 | 64 | 0 | — | Gemma model card; Gemma 3 ships same top-k/top-p |
| deepseek | 0.6 | 0.95 | — | 0 | — | DeepSeek-R1 card, "0.5–0.7, 0.6 recommended" |

**Families with no published numbers** (the code names Llama 3.x and Mistral-Nemo) keep the app's generic defaults. **No number may be invented to fill a gap** — a gap is recorded as a gap.

## 4. Parameter policy

| Parameter | Role | Starting point | Adjustment policy |
|---|---|---|---|
| Temperature | Prose variety vs consistency | Family value if published; else generic | Adjust only with a recorded reason and a re-frozen baseline |
| top-p | Tail truncation | 0.95 where published | Do not stack with aggressive top-k "to be safe" |
| top-k | Hard candidate cap | Family value if published | Leave unset where unpublished |
| min-p | Relative floor | 0 where published as 0 | `0` is a **value**, not "unset" — the code's `hasMinP` flag exists for this |
| repeat / presence / frequency penalty | Loop control | Published value, else 0 | Raising penalties to fix repetition changes style; record it as a condition |
| context budget | Memory pressure | Model's usable context, measured | The **test variable** for T-05, never silently changed |
| n_predict | Response length | Short, fixed per test | Must be identical across compared conditions |

## 5. Procedure

1. **Read family** from GGUF metadata; look up published values.
2. **Apply L-B** where published; note the source line.
3. **Run ordinary conversation** — a few turns of realistic chat, no task framing.
4. **Identify degradations**: loops, flatness, refusal, persona collapse, over-verbosity.
5. **Record the chosen baseline** and why, per model, before any substantive test.
6. **Freeze** it for that model for the duration of a test series.

## 6. Secondary profiles

| Profile | Purpose | Constraint |
|---|---|---|
| Creative-writing | Compare prose quality against BCB | **Secondary only** — never replaces BCB as reference |
| Reasoning | Compare deliberation-heavy operation | Secondary only |

A secondary profile may be used to *explain* a BCB result. It may not be used to *replace* a disappointing BCB result.

## 7. Guardrails

1. **The baseline is frozen before the test, not after the result.** Changing a sampler mid-series makes every earlier run a different condition.
2. **Per-model values are legitimate; per-test values are not.** Adjusting temperature because a model loops is calibration; adjusting it because a test failed is contamination.
3. **Every run records the literal values used**, not the profile name.
4. **No claim of universal optimality.** BCB is "sensible ordinary chat on this machine", nothing more.

## 8. Claim ceiling

This proposal establishes a *method* for choosing settings without inventing them. It does not establish that any family value is optimal for companion chat, that the generic defaults are good, or that BCB produces comparable behaviour across families — only that the comparison is documented well enough for the owner to judge.
