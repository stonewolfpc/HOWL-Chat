# `MEMORY-09` — Authorisation Package Completeness Check

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-BEAC-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **COMPLETENESS CHECK — nothing run; `BEA-1` scope unchanged.**
**Target:** `MEMORY-09-BEA-1` (row 481).
**Rule applied:** this document adds *specifics* the instrument lacks. It does not widen scope. **Signature on this schedule incorporates it into `BEA-1`; neither document alone is effective.**

---

## 1. Verdict

**The proposal is not yet signable.** Of six required elements, **two are adequately addressed and four are missing or blank.**

| # | Required element | State in `BEA-1` |
|---|---|---|
| 1 | Exact test design and question | **MISSING** — "one controlled test" is not a design |
| 2 | Named models and providers | **MISSING** — "at most 3 identities" names none |
| 3 | Maximum calls including retries | **BLANK** — owner-set, no arithmetic |
| 4 | Hard cost ceiling | **BLANK** |
| 5 | Failure / inconclusive handling | **MISSING** |
| 6 | Permitted artifacts and filing location | **MISSING** |
| + | Harness build vs execution separable | **PARTIAL** — implied, not split |
| + | One concept, one contrast, model-alone comparison | **ADEQUATE** in intent |

## 2. Element 1 — the exact test

**Question:** with a canonical fact supplied by state, does a reply honour it where the same model without that state does not?

**Design — one fact, one probe, two conditions, three models.**

| Item | Value |
|---|---|
| Fixture | `F-ST-01-S1` (derived from row 478 §4 F-ST-01) |
| **Canonical fact** | An **invented**, unguessable fact: the character's mentor is *"Corvin Ashgrave of Ossary"* |
| Probe | One user turn: *"Remind me — who trained you?"* |
| Condition **A** (assisted) | System context carries the fact verbatim |
| Condition **B** (unassisted) | Identical prompt, fact absent |
| Reps | 3 per condition per model |
| Calls | 3 models × 2 conditions × 3 reps = **18 planned** |

**Why an invented fact, and this is the design's whole point:** a model cannot know *"Corvin Ashgrave of Ossary"* from priors. **`C3` is excluded by construction**, not by judgement. Any correct mention in A is state-supplied; any correct mention in B is an alarm (§5).

This is narrower than T-11 — deliberately. One concept, one contrast, exact capture.

## 3. Element 2 — models and providers

**Slots must be named at signature.** Until then they are labelled, not chosen:

| Slot | Intent | Model id | Provider |
|---|---|---|---|
| **S1-A** | Creative-writing-oriented | `___________` | `___________` |
| **S1-B** | General instruct, different family | `___________` | `___________` |
| **S1-C** | Optional third | `___________` | `___________` |

**Rule:** one provider pinned per model, recorded with revision. **A provider change voids that model's series** (row 477 P-1). No `provider="auto"`, no routing suffixes.

## 4. Element 3 — call arithmetic

| Stage | Calls |
|---|---|
| Planned | **18** |
| Retry allowance (≤1 per failed call, max 3 attempts total per call) | +12 |
| **Hard ceiling** | **30** — the phase stops at 30 calls regardless of remaining plan |

A call that fails **counts against the ceiling**. Exhausted retries are recorded, not re-planned.

## 5. Element 4 — cost

- **Ceiling:** `USD ______` — set by the owner at signature. No default is proposed, because a guessed ceiling is worse than a blank one.
- **Pre-flight:** if provider pricing for a named model cannot be retrieved **before** the run, that model is **not run** and the gap is reported. No pricing, no spend.
- **Stop rule:** projected overrun stops the phase and reports partial results.

## 6. Element 5 — failure and inconclusive handling

| Event | Action |
|---|---|
| Pinned provider 503 / retries exhausted | That model's series **void**; other models proceed |
| Fewer than 2 models complete | **Inconclusive** — reported as such, no concept claim |
| **B produces the canonical fact** | **Alarm**: the fact was not unguessable → run void, design defect reported |
| A ≈ B (no difference) | **Null result** — reported at equal prominence; not a failure of the harness |
| A partially honours (mentions name, wrong place) | Classified `C4`-adjacent; recorded, not merged with success |
| Cost ceiling reached | Stop; report partial + what was not run |
| Any unrecorded parameter | Run inadmissible (row 477 §8) |

## 7. Element 6 — artifacts and filing

**Two separate permissions.** This is the split the owner asked for:

| Permission | Allows | Network | Location |
|---|---|---|---|
| **P1 · Build** | Adapter, state layer, fixture loader, runner skeleton, record writer | **No** | `eval/harness/`, `eval/fixtures/` |
| **P2 · Execute** | 18 planned calls under caps | **Yes** | `eval/runs/S1-<date>/` |

**No app source is touched.** `eval/` is new and outside the Howl build. Every run writes: `record.json` (provenance per row 477 §8), `raw/` (bodies), `state/` (hashes before/after), `README.md` (one-sentence verdict, three observations, evidence — row 474 §2).

**P1 alone produces no evidence.** It produces a harness that cannot contact a provider until P2 is granted.

## 8. Two residual risks the check found

1. **Three models may be one control too many.** If all three behave identically, the extra calls bought nothing. Acceptable for S1 (family spread is the point), but S1's result is *"the concept held or failed on these models"*, never a rate.
2. **A single invented fact is a weak memory claim.** It tests *retention of one arbitrary datum*, not the architecture. **That is the correct scope for S1** — but it must not be reported as validating `MEMORY-09`.

## 9. Signature

> **P1 · Build** ☐ authorise ☐ decline
> **P2 · Execute** ☐ authorise with the limits below ☐ decline
>
> Models: `______` / `______` / `______` · Providers: `______` / `______` / `______`
> Call ceiling: **30** ☐ confirmed · Cost ceiling: `USD ______`
> Date: `______`
>
> *Unsigned, this document authorises nothing, and the four gaps in §1 remain open in `BEA-1`.*
