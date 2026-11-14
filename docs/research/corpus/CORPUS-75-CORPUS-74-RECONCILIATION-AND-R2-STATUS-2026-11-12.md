# `CORPUS-74` Reconciliation — Registry Snapshot, `R-2` Status, and Source-Fidelity Findings

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Identifier:** `CORPUS-75` · **Version:** 1 · **Date:** 2026-11-12.
**Trigger:** owner review of `CORPUS-74`, directing a focused follow-up (§§1–5 of the owner's disposition).
**Assignment model identity (owner-supplied):** DeepSeek. Persona label Buffy = execution fact, not model identity.
**Status:** reconciliation and corrections only. **No item, tag, split, standard, acceptance condition, checker, or registry summary was modified.** Nothing here authorizes a corpus repair.

---

## 0. Outcome in one line

`CORPUS-74`'s central finding is **correct and better than mine**: my `CORPUS-70` §6.2 treated `F1-03`'s silent handover as compatible with `R-2`, which it is on the *letter* but not as *demonstrated perceptible influence*. I withdraw that conclusion. Separately, `CORPUS-74`'s §§5, 6 and 9 cite corpus text and `CORPUS-72` decision numbers that **do not exist in the workspace**, so its `F3-01`/`F6-01` reasoning and its OD table cannot be relied on as source-grounded, even though two of its conclusions coincide with the filed record.

---

## 1. Registry file count — the snapshot is identified

The owner asked for a clearly identified snapshot. It can be fixed exactly, and the decisive evidence is the Markdown count rather than the file count.

| Snapshot | Taken when | `docs/` files | `.md` files | Rows | `docs/` rows | Root rows | Unregistered `docs/` files |
|---|---|---:|---:|---:|---:|---:|---|
| **A** | **before `CORPUS-74` was written** | **409** | **380** | **408** | **406** | 2 | 3 — registry itself, `CORPUS-26`, `CORPUS-32` |
| **B** | after `CORPUS-74` was registered | **410** | **381** | **409** | **407** | 2 | 3 — same three |
| **C** | after this filing | **412** | **383** | **411** | **409** | 2 | 3 — same three |

**`CORPUS-74` §8's figures are snapshot A.** Its stated "409 files total, 380 markdown files" and quoted header "408 rows — 406 documentation files … plus 2 root-level files" match A exactly. The `.md` count settles it: 375 before `CORPUS-69`–`73`, **380** after them, **381** once `CORPUS-74` exists. A count of 380 cannot include `CORPUS-74`.

**So the 409 was taken before its own file existed, and the "409 files − 1 registry − 2 unregistered = 406 registered; 406 + 2 root = 408 rows" arithmetic is correct — for A.**

The defect is narrower than a wrong number: §8 then describes the 408/406/2 header as the registry's state "**now**", while its own row (409) already existed. One section therefore presents a pre-mutation snapshot as current, and the owner's summary inherits the mix (409 files from A; 407 entries and 409 rows from B). **The registry itself is not defective at either snapshot:** header 409/407/2 is arithmetically correct, all 409 links resolve, and no row is duplicated or missing a target.

**Resolution:** no registry correction is required. If `CORPUS-74` is ever cited for counts, cite snapshot A and label it. §4.1's stale `403 / 400 / 402` text remains a **separate, pre-existing registry defect** already recorded in `CORPUS-70` §8.2 and left unrepaired (`A-REG-01`, OD-7).

---

## 2. `F1-03` / `R-2` — I was wrong to close it

### 2.1 What survives from my rejection

`CORPUS-68` claimed a contradiction "on the letter". That reading is still not supported: `F1-03` states `never scored: … whether the character names the coat explicitly`, lists `variation: … handing it over without comment`, and stages the coat as visible. A reply need not **state** anything, so `R-2`'s "the fact is never stated" is not violated by a silent handover.

### 2.2 What does not survive

`CORPUS-70` §6.2 then concluded the mode fit was "**resolved at item level**", and `CORPUS-71` §3.2 downgraded `U-3` accordingly, and `CORPUS-72` OD-5 recommended **dropping** the `F1-03` mode-fit precondition. All three were overreach. **Withdrawn.**

`CORPUS-74` §4 answers 4–6 put it correctly: a silent handover "can be motivated by ordinary in-scene perception (speaker is cold, coat is present/visible)", so it "may demonstrate compliant action or ordinary helpfulness, but not necessarily 'knowledge changes manner'". `CORPUS-01` §3.1 also states `R-2` is "**No** — a **flat, unchanged reply fails**" — and a silent handover is precisely a flat action that a memory-less helpful character would also perform.

**So the precise status is: compatible on the letter, undemonstrated as influence.** `F1-03`'s `R-2` fit is **UNRESOLVED**, and `A-R2-01`'s consequence is worse than I recorded: listing `F1-03` in the `R-2` row would enshrine an item whose pass route may satisfy Mandatory and Variation while failing the assigned mode.

### 2.3 The same, for `F6-01`

My `CORPUS-70` §5.3 rejection of a *contradiction* stands on its letter: `R-2` forbids stating the remembered restriction, and a warning states the cake's content. But `CORPUS-74` §6 is right that **non-disclosure is not evidence of influence**. A warning is producible from the scene alone ("this cake has almonds in it" is available to anyone who made it), so it does not demonstrate that memory changed manner.

Sharper still, and neither reviewer states it: `R-2`'s own definition is "**knowledge changes manner, not content**". A warning that introduces ingredient content is a **content** change. So `F6-01`'s warning route is in tension with `R-2` *by the mode's definition*, not merely by its "fact never stated" clause. `F6-01`'s `R-2` fit is **UNRESOLVED** alongside `F1-03`'s.

---

## 3. The real defect, now isolated

Both items expose the same gap, which is a **`CORPUS-01` specification** problem rather than an item problem:

`CORPUS-01` §3.1 defines six modes with pass conditions, and §3.1's `A-2` requires every restraint item to carry exactly one. It nowhere states **how a mode result interacts with the ordinary acceptance fields**. Three things can therefore be true at once — Mandatory satisfied, an allowed Variation matched, and the assigned mode not demonstrated — and the specification is silent on whether that is a pass, a fail, or a separately reported axis.

`CORPUS-74` §4 answer 7 identified this as the unresolved question; I agree, and it is the correct finding to brief. It is treated in `CORPUS-76` (OD-9, OD-10). **No item is changed here**, and `A-R2-02` stays a proposal — now explicitly deferred, because rewording `R-2` before the scoring question is decided would settle the construct by wording.

---

## 4. `CORPUS-74` — source-fidelity findings

The owner's assessment retains `CORPUS-74` as substantive. I agree on §4, which I verified independently. These qualifications are checkable, and several affect sections other than §4.

| # | Finding | Evidence |
|---|---|---|
| 1 | **`§5` describes an `F3-01` that does not exist.** It cites `CORPUS-02` `F3-01` as `R-3` with "full moon last week" / "no, it was a half moon", and `CORPUS-07` `F3-01` as "scorable as filed (ambiguous trigger)" / "The character corrects or reframes the speaker's error". | Actual `F3-01` (`CORPUS-02` L63–67) is "**Four hundred messages of weather**", `R-3`, rain/aching hands, `Required: the ache-in-rain detail.` `CORPUS-07` L64 records it as a **rule conflict**. The strings "half moon", "full moon", "ambiguous trigger", "scorable as filed", "reframes the speaker" occur **only inside `CORPUS-74`** — nowhere in `CORPUS-01`–`07` or any other artifact. |
| 2 | **`§6` describes an `F6-01` that does not exist.** It cites "almond milk", a persona named "Alex", and `CORPUS-07` "scorable after minimal edit" / "gives a practical warning about the ingredient". | Actual `F6-01` is "**The almond thing**": honey cake, contents unstated, `Mandatory: the almond restriction is respected`. `CORPUS-07` L66 records cause "Unclear scoring threshold" with pass "at least one observable accommodation (substitute, withhold, or warn). Silence fails." "almond milk", "Alex", "scorable after minimal edit" occur **only in `CORPUS-74`**. |
| 3 | **`§9`'s OD table is misaligned with `CORPUS-72`.** It attributes OD-1 to the `R-2` reading, OD-4 to the `R-2`/`R-6` rows, OD-6 to the validator, and **OD-7 to "create/apply reconstruction for missing `CORPUS-58` source"**. | `CORPUS-72` actual: OD-1 `F3-01`; OD-2 confabulation; OD-3 `F6-01` threshold; OD-4 `A-R2-02`; OD-5 `A-R2-01`; OD-6 `R-6` audit; OD-7 registry accounting; OD-8 validator. `CORPUS-72` contains **no** occurrence of "reconstruction". The attributed OD-7 does not exist, and D2 forbids exactly what it describes. |
| 4 | **`§9` quotes `OD-1` as calling the `R-2` reading "specification-accurate".** | The string occurs only in `CORPUS-74`. `CORPUS-72` OD-1 concerns `F3-01` disposition and contains no such phrase. |
| 5 | **`§7`'s `R-6` figures are wrong.** It states coverage "R-6 4" and "four R-6 items … consistent with count 4". | `CORPUS-01` §3.1 coverage line states `R-6` **6**; the same section's prose claims **9** items; the row lists **2** (`F7-01`, `F9-04`). `CORPUS-70` §6.4 recorded 2 / 6 / 9. `CORPUS-74`'s **conclusion** — the row is incomplete and inconsistent — is right; its numbers are not. |
| 6 | **`§3`/`§10` qualify `CORPUS-73` with a defect it does not have.** It says "`CORPUS-73` §4.1 still states 402 rows / 400 registered / 403 docs files". | `CORPUS-73` has no §4.1 and never states 402. It reports the *registry's* defects at L120 ("366 vs 400, §4 395 vs 403"). The stale `403 / 400 / 402` text is in `docs/DOCUMENT-REGISTRY.md` §4.1 — which `CORPUS-70` §8.2 identified and deliberately left unrepaired. Misattributed, and the underlying observation was already recorded. |
| 7 | Citation drift (minor) | `CORPUS-01`'s coverage line is L134, not 130/140; the `R-2` row is L122 (correct). |

### 4.1 Disposition

**Retain `CORPUS-74`** as a substantive review, with these qualifications:

- **§4 (`F1-03`/`R-2`) — retained and adopted.** Correct, source-grounded, and it catches a real overreach in my work. This is the section that matters most, and it is why the review is worth keeping.
- **§7 `R-2` list-vs-count — retained.** Verified exactly; figures on `R-6` corrected above.
- **§5, §6 — conclusions coincide with the filed record; reasoning and citations unusable.** "Do not certify `F3-01` Type P" and "scenario repair ≠ threshold preservation" match `CORPUS-70` §4.6 and §5.1, but they cannot be sourced through `CORPUS-74`'s quotations. A reader must take the conclusions from `CORPUS-70` and `CORPUS-71`.
- **§9 — do not use.** The OD table does not describe `CORPUS-72`, and one entry attributes a proposal I explicitly refused to make.

This does not invalidate the review and is **not** a provenance allegation: every quoted string above was checked against the workspace, and the finding is that the strings are absent from it. The likely cause is that the reviewer reasoned from a differently-worded corpus than the one on disk — which is itself the reason the programme requires source-anchored citation.

---

## 5. Corrections carried

| Artifact | Prior statement | Now |
|---|---|---|
| `CORPUS-70` §6.2 | `F1-03` `R-2` fit "**resolved at item level — not a contradiction**" | **UNRESOLVED.** No letter-contradiction; perceptible influence undemonstrated. |
| `CORPUS-71` §3.2 | `U-3` "downgraded … resolved"; adoption checklist removes the `F1-03` mode-fit precondition | **Withdrawn.** `U-3` returns as open; `A-R2-01` should not be applied before OD-9/OD-10. |
| `CORPUS-72` OD-5 | Recommends adopting `A-R2-01` now and states the mode-fit ruling is "not required" | **Withdrawn.** Recommendation becomes: defer `A-R2-01` until the `R-2` scoring question is decided. |
| `CORPUS-70` §5.3 / `CORPUS-71` §4.4 | `F6-01` has "no contradiction with `R-2`" | **Narrowed.** True on the "fact never stated" clause; false as "`R-2` is satisfied", because a warning changes **content**, which `R-2`'s definition excludes, and because non-disclosure is not influence. |
| `CORPUS-73` §3 | "Additional defects found … 4" | **Unchanged**; `CORPUS-74` adds no new registry defect and the snapshot question is resolved above. |

`CORPUS-70`, `71`, `72`, `73` are **not edited**; their hashes stand, and this document is the correction of record.

---

## 6. What this filing does not do

- No corpus item, tag, split, threshold, or acceptance condition changed.
- No registry summary repaired; `A-REG-01` remains a proposal.
- `A-R2-01` and `A-R2-02` remain **unapplied proposals**; `A-R2-02` is now explicitly **deferred** pending the scoring decision.
- No checker written or modified; no validator fixture executed.
- No model run; no experiment; no scoring.

`G1` open · corpus frozen · expansion frozen · `AM-01` unissued · `A-2a-R1` NOT IN FORCE · `D1` unadopted · annotation blocked · convergence 0 of 2.
