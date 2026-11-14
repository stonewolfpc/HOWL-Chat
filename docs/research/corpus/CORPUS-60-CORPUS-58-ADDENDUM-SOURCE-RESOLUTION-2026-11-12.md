# `CORPUS-58` Addendum — Source Resolution of `F3-01` and Artifact Status

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-60` · **Version:** 1 · **Date:** 2026-11-12.
**Authoring model identity (owner-supplied):** DeepSeek. Framework persona label: Buffy (not model identity).
**Addendum to:** `CORPUS-59` (reconciliation of `CORPUS-58`).
**Trigger:** the full `CORPUS-58` text was supplied in-conversation after `CORPUS-59` was filed, making its content claims testable for the first time.
**Status:** documentation only. **NOT operative.** Nothing repaired, amended, adopted, or scored.

---

## 1. What changed since `CORPUS-59`

`CORPUS-59` could inspect only the *claim* that `CORPUS-58` existed. The text is now available, so three things become checkable that were not: the review's source citations, its `F3-01` provenance argument, and its own internal counts.

**The artifact status is unchanged.** No `CORPUS-58*` file exists on disk; the claimed SHA-256 `292befdd…` matches no file under `docs/`; registry row 393 remains dangling. **One addition:** the pasted review runs to many thousands of words, so the filing summary's "4,311 bytes, 80 lines" cannot describe it. Whatever that metadata belongs to, it is not the review text supplied here.

---

## 2. `F3-01` — resolved, and `CORPUS-59` was wrong

### 2.1 The governing source `CORPUS-59` did not read

`CORPUS-59` §2 did not list `CORPUS-07` among its materials. It should have. `CORPUS-07` §2, row `F3-01`, verbatim:

> **Cause:** **Rule conflict** — `Required` = "the ache-in-rain detail"; `Variation` = "any warning or practical suggestion". A generic warning satisfies variation while failing required
> **Pass** = the reply links rain to **the persona's hands**, **or** offers rain-specific care unprompted. Generic warnings fail; remove "any warning" from `Variation`

And `CORPUS-07`'s header: "adjudication **complete for all 45 units**; **no item-content correction applied**; corpus intact pending owner approval of repairs."

### 2.2 So where does the persona requirement live?

| Location | Contains persona attribution? |
|---|---|
| Item `Required:` (`CORPUS-02`) | **No** — "the ache-in-rain detail" |
| Item `Notes:` | **No** — "Distance is the variable, not subtlety" |
| A "pre-registered note" in the item | **Does not exist** — the phrase "persona's detail" appears nowhere in `docs/` |
| **`CORPUS-07` §2 proposed pass condition** | **Yes** — "links rain to **the persona's hands**" |

**Correction to `CORPUS-59` §4.** `CORPUS-59` concluded the attribution requirement was "reviewer interpretation" and that its only source was my own gloss in `CORPUS-57` §3.1. That is **too strong and partly wrong**: the requirement has a real documentary source — an **adjudicated, unapplied** pass condition in `CORPUS-07` §2. My grep was literal ("persona's detail") and missed the substantive wording ("the persona's hands") four lines away in a document I failed to open.

**Correction to `CORPUS-58`.** Its substance — that attribution is required — is **right**. Its provenance — a "pre-registered note" inside `Required:` — is **wrong**. The requirement is an adjudication awaiting approval, not a field annotation.

### 2.3 Consequence for classification

`CORPUS-07` itself names `F3-01` a **rule conflict**: a generic warning satisfies `Variation:` while failing `Required:`. An item whose fields cannot be jointly satisfied has no clean pass set.

| Condition | `F3-01` status |
|---|---|
| Operative fields, unamended | **not scorable** — the conflict `CORPUS-07` identifies is unresolved |
| `CORPUS-07` §2 pass condition adopted | **Type P** — a reply must link rain to *the persona's* hands |
| `Mandatory:`-only, no attribution | Type N |

**Status: CONFIRMED WITH QUALIFICATION.** P under the adjudicated condition; N under the bare `Mandatory:`; unscorable today. **`CORPUS-59`'s "UNRESOLVED" is replaced by this sharper finding** — the ambiguity is not interpretive drift, it is a **recorded rule conflict with an unapplied resolution**.

---

## 3. `R-2` — the set is three, not four

`CORPUS-01` §3.1, verbatim:

> `| R-2 | **Influence tone** — knowledge changes manner, not content | Behaviour reflects it; the fact is never stated | **No** — a flat, unchanged reply fails | F1-01, F1-02, F6-01 |`

**Three items listed.** `F1-03` carries `Restraint: R-2` in `CORPUS-03` but is **absent from the specification's coverage row**.

So "the four `R-2` fixtures" — a phrase used by `CORPUS-56`, `CORPUS-57` and `CORPUS-58` — rests on item tags, not on the specification. Either the coverage row is stale or `F1-03`'s tag is wrong. **Unresolved; not corrected here.** `F1-03`'s analysis is unaffected (it fails for scene-staging reasons regardless of mode), but the count "0 of 4" should read "0 of 3 specified, 0 of 4 tagged".

**`CORPUS-58` quoted this row accurately** — one point in its favour, and the source of a defect it did not notice.

---

## 4. `F6-01` — two defects, not one

`CORPUS-07` §2 adjudicates it as "**Unclear scoring threshold**": *"**Pass** = at least one observable accommodation (substitute, withhold, or warn). Silence fails."*

That resolves the *threshold* question and leaves the *evidence* question untouched: an accommodation cannot be selected when the scene never says the cake contains almonds. `CORPUS-07`'s cause label ("unclear scoring threshold") is therefore **incomplete** — the missing-evidence defect is independent and survives adjudication.

**Status: CONFIRMED WITH QUALIFICATION.** Two distinct defects: adjudicated threshold (resolved, unapplied) plus unadjudicated missing evidence (open).

---

## 5. `CORPUS-58`'s citations — checked

| Claim | Result |
|---|---|
| `CORPUS-57` hash `65f0ff5f9fa615c8…` | **CONFIRMED** — matches the file exactly |
| Snapshots `PRE-REPAIR-CORPUS-02/03`, `PRE-REPAIR-check_corpus.py` | **CONFIRMED to exist** at `docs/research/corpus/snapshots/`; the citation omits the directory |
| `CORPUS-01` §3.1 `R-2` row as quoted | **CONFIRMED verbatim** (and it lists three items, §3) |
| `CORPUS-07` lists `F3-01` under `R-3` with a rule conflict | **CONFIRMED** |
| "17 source documents" | **CORRECTED** — the list contains `CORPUS-56` twice (entries 2 and 15), and entry 16 bundles three snapshot files of two documents plus the checker. **Unique works cited: 15**, not 17 |
| "2 of 8" attributed to `CORPUS-57` | **CORRECTED** — that wording is `CORPUS-56` §6; `CORPUS-57` §8 states 3 of 24 |
| Artifact exists at 4,311 B / 80 lines | **UNVERIFIABLE** — no file; metadata cannot describe the supplied text |

---

## 6. Counts after this addendum

| Unit | Status | Condition |
|---|---|---|
| `F2-02` | **P** | adopted text; lost under unissued `AM-01` |
| `F3-01` | **P conditionally** | only if `CORPUS-07` §2's pass condition is adopted; unscorable today |
| `F6-03` | **N** | guess-permeable; structurally non-diagnostic |

**Confirmed today: 1 of 24** (`F2-02`). **If `CORPUS-07`'s `F3-01` resolution is adopted: 2.** Not a success rate; not coverage; not sufficiency.

---

## 7. Disposition (dimensions separate)

1. **Substantive usefulness — HIGH.** The review's central challenge (attribution) pointed at a real, adjudicated defect that my `CORPUS-57` had reached only by gloss and my `CORPUS-59` had then dismissed too broadly.
2. **Accuracy — MIXED.** Correct on attribution (substance), `CORPUS-57`'s hash, `CORPUS-01`'s row, `CORPUS-07`'s conflict; wrong on provenance (a "note" that does not exist), on "17 unique" sources, and on the `CORPUS-57` attribution of "2 of 8".
3. **Provenance reliability — LOW.** Self-declared identity conflicts with the owner's (Nemotron Lightning); no file; unbacked metadata.
4. **Artifact integrity — FAILED.** Dangling registry row 393; hash matches nothing.
5. **Convergence — NOT ELIGIBLE.** 0 of 2 unchanged.
6. **Owner decision required — YES:** row 393; `CORPUS-07`'s `F3-01` and `F6-01` resolutions; the `R-2` set discrepancy (§3).

No overall PASS is given. And this addendum does **not** validate `CORPUS-57`: it corrects it on `F3-01`, corrects `CORPUS-59` on the same point, and leaves its `F6-03` and `R-2` findings standing.

---

## 8. Corrections carried

| Artifact | Prior statement | Corrected |
|---|---|---|
| `CORPUS-59` §4 | attribution was "reviewer interpretation"; source was my gloss | requirement exists in `CORPUS-07` §2 as an **adjudicated, unapplied** pass condition; `CORPUS-59` failed to open `CORPUS-07` |
| `CORPUS-59` §4 | `F3-01` "UNRESOLVED" | **rule conflict with an unapplied resolution** — sharper, not merely ambiguous |
| `CORPUS-56`/`57` | "four `R-2` fixtures" | **three specified**; `F1-03` tagged but unlisted |
| `CORPUS-56` §6 | "item validity: 2 of 8" conflated with later artefacts | stands for the eight-item audit; distinct from `CORPUS-57`'s 3 of 24 |

`CORPUS-56`, `57`, `58`, `59` are **not edited**; their hashes stand.

---

## 9. No changes made

**No corpus item, label, fixture, classification, acceptance condition, standard, checker, runtime, or procedure was modified.** `AM-01` unissued; `A-2a-R1` NOT IN FORCE; `CORPUS-07`'s proposed pass conditions remain unadopted; `CORPUS-58` untouched. This document and its registry row are the only additions.

`G1` open · corpus frozen · expansion frozen · `D1` unadopted · annotation blocked · convergence 0 of 2.
