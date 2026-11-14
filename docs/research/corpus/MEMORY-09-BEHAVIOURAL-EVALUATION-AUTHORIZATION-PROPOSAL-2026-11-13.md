# `MEMORY-09` — Behavioural Evaluation Authorisation (Proposal)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-BEA-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **PROPOSED — NOT EFFECTIVE. Requires the owner's signature.**
**Grounds:** row 476 (source audit): `OD-M09-17` is bounded to cue-level judgement, authorises no evaluation, and does not reach behavioural observation.

---

## 1. Why a separate instrument is required

The audit established a three-part finding that this instrument exists to satisfy:

1. `OD-M09-17` **does not reach** behavioural evaluation.
2. It **grants no permission** — the addendum "authorises no … evaluation, evidence collection".
3. Therefore behavioural evaluation **needs its own explicit authorisation**, which no existing record supplies.

**This is not** a reinterpretation of `OD-M09-17`, and it does not amend it. It is a separate grant, deliberately narrow.

## 2. What it would authorise

| # | Authorised | Bounds |
|---|---|---|
| 1 | Build the Python harness: adapter (row 477 v1.1), state layer and fixtures (row 478) | Code only; no corpus access |
| 2 | **One** controlled test — S1 concept check, assisted vs unassisted (row 479 §6) | Named test id; frozen configs |
| 3 | Hosted inference via `huggingface_hub.InferenceClient` | Provider pinned; no `auto` |
| 4 | Named models only | At most **3** model identities in S1 |
| 5 | Bounded spend | A hard run cap and cost cap, set by the owner at signature |

## 3. Mandatory conditions

1. **Adapter compliance** with row 477 v1.1: provider pinned, `auto` prohibited, no model-id routing suffixes, fail-closed on 503.
2. **Provenance completeness**: every call records model id + revision, provider requested/used, params sent, returned model, `system_fingerprint` where supplied, client version, hashes; **missing fields marked `unavailable`**.
3. **State evidence**: every run records state hashes before and after (row 478 §3).
4. **Records** in the A–D shape (row 470 §4) with `[MECH] [OBS] [AI] [OWNER]` tags.
5. **No-bridge rule** (row 480 §5): engineering results may not be combined with cue results into one claim.
6. **Fail-closed** (row 477 F-1): a failed pinned call is recorded, never rerouted.

## 4. What it does not authorise

| Not authorised | Because |
|---|---|
| S2, S3, S4 | Separate stages, separate decisions |
| Local inference or hardware measurement | S3 |
| Any Howl source or runtime change | Out of scope entirely |
| Corpus, tag, split, score, standard changes | Standing prohibition |
| Rating, rater recruitment, α on real data | Different instrument |
| Two-rater or reproducibility claims | `OD-M09-17` territory, unmet |
| Citing engineering results alongside cue results | No-bridge rule |
| `G2`, `P-3`, `L4` | Unauthorised in every order to date |

## 5. Bounds and lapse

- **Run cap:** set at signature; no run beyond it without a new grant.
- **Cost cap:** set at signature; a projected overrun stops the phase rather than continuing.
- **No exploratory sweeps:** every run belongs to a named test id.
- **Lapse:** the grant expires when S1 is reported or the cap is reached, whichever comes first. Progress on one does not extend the other.

## 6. Records required at completion

A single S1 record containing: the test id and frozen config, the models and providers used, every call's provenance block, the state hashes, the assisted/unassisted comparison with its disclosure of any prompt-shape difference, the `C2`/`C3` classification with its basis, and an explicit statement of what the result does **not** establish.

## 7. Non-claims

This instrument, if signed, would authorise **work**, not conclusions. It would not establish that the memory concept works, that any model is compatible, or that `MEMORY-09` improves behaviour. A null result is a valid outcome and is reportable at the same prominence as a positive one.

## 8. Signature

> **Owner decision:** ☐ Authorise as written ☐ Authorise with the following changes ☐ Decline
>
> Run cap: ______ · Cost cap: ______ · Date: ______
>
> *Until signed, nothing in this document authorises any activity.*

> **Completeness note (2026-11-13).** This instrument is **not yet signable**. The check at `MEMORY-09-AUTHORIZATION-PACKAGE-COMPLETENESS-CHECK-2026-11-13.md` (row 482) finds four of its six required particulars missing or blank — exact test design, named models and providers, call ceiling, cost ceiling, failure handling and filing locations — and supplies them. **Signature on the check schedule incorporates it here; this document alone remains incomplete.**
