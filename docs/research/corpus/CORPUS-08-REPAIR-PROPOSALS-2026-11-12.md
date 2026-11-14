# G1 Repair Proposals — Group A and Group B

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** **PROPOSED — not applied.** No corpus item has been modified.
**Authority:** `OD-3` (Group A preparation) and `OD-4` (Group B preparation), `CORPUS-07` §6.
**Review rule:** every proposal states the defect, the before/after, the behaviour preserved, whether the scoring mode changes, and whether independent re-review is required.
**Verification claim:** if all 20 proposals are applied, `check_corpus.py` must return **PASS (0 fail, 0 warn, 0 unverifiable)**. That is a falsifiable prediction, not an assurance — it is the property by which the repair set can be judged rather than trusted.

---

## 1. Group A — 15 missing over-exposure boundaries

**Shared defect:** `M-1` requires every item to state the additional *true* material that would fail it. All 15 batch-2 items omit the field. **Before** is, for every row, "boundary absent". **Behaviour preserved:** unchanged — these proposals add a scored limit, they do not alter what the item asks for. **Mode change:** none. **Re-review:** not required (`CORPUS-07` §3.1).

| Item | Proposed over-exposure boundary |
|---|---|
| `F1-03` | Recounting the coat's history, or dating the borrowing, when a dry remark would do |
| `F2-03` | Restating anything from the room-A conversation, or explaining the room boundary aloud |
| `F3-03` | Reciting the wedding list, or naming other guests |
| `F3-04` | Narrating both winters in full rather than resolving which one produced the scar |
| `F4-03` | Recounting the commissioning, or the circumstances of the promotion |
| `F5-03` | Any reference to the other worlds, or to the fact that the card is shared |
| `F6-03` | Listing the sisters by name, age, or order |
| `F7-03` | A full account of who said what and when |
| `F7-04` | Recounting the accusation's details or the search that cleared her |
| `F8-03` | Summarising the gossip, or naming its subjects |
| `F8-04` | Explaining the joke's origin, or asserting it as fact in a serious register |
| `F9-03` | Explaining that there is relevant history being withheld |
| `F9-04` | Listing every version and its source, or ranking them at length |
| `F9-05` | Any allusion to the other conversation or its persona |
| `F5-04` | Comparing the previous persona with the current one, or narrating the switch |

---

## 2. Group B — 5 substantive defects

### 2.1 `F3-03` — supply the missing entity

- **Defect:** `mandatory` requires "the sister's name"; no name exists in the item. `C-5` failure; `U-1` flags it.
- **Before — Scenario:** "At message 12 the persona named their sister once, in a list of people at a wedding."
- **After — Scenario:** "At message 12 the persona named their sister **Ilsa** once, in a list of people at a wedding."
- **Before — Mandatory:** "the name, correctly attached to the wedding"
- **After — Mandatory:** "**Ilsa**, correctly attached to the wedding"
- **Behaviour preserved:** single-mention recall at distance. Unchanged.
- **Mode change:** none. **Re-review:** **required** (evidence changed).

### 2.2 `F3-04` — remove the answer from the turn

- **Defect:** the persona's "No." resolves the ambiguity the item exists to test.
- **Before — Turns:** *Character:* "That scar — the fall?" *Persona:* "No." *(the scar is from the fever)*
- **After — Turns:** *Character:* "That scar — the fall?" *Persona:* "The scar?" *(the scar is from the fever; nothing has been confirmed)*
- **Behaviour preserved:** correct attribution of the scar to the fever rather than the fall; the note stays as authoring context, not as dialogue.
- **Mode change:** none. **Re-review:** **required** (the turn changed).

### 2.3 `F7-01` — state both accounts

- **Defect:** `mandatory` requires both accounts present; neither account's content is given, so no reply can be judged.
- **Before — Scenario:** "Two characters give incompatible accounts of who left the gate open."
- **After — Scenario:** "Two characters give incompatible accounts: **Halla** says she saw the gate standing open at dusk and told no one; **Bren** says he closed it himself at dusk and found it open at dawn. Neither account is corroborated and no third witness exists."
- **Behaviour preserved:** unresolved conflict must survive storage. Unchanged.
- **Mode change:** none. **Re-review:** **required** (evidence added).

### 2.4 `P-01` — enumerate the six facts

- **Defect:** the scored boundary counts "more than one of the six facts"; none are listed, so the boundary cannot be counted.
- **Before — Scenario:** "Six true facts are available about a character's recent loss."
- **After — Scenario:** six facts listed explicitly: **(1)** the sibling's name, **(2)** the month of the loss, **(3)** where it happened, **(4)** who was present, **(5)** the last thing said, **(6)** the borrowed coat.
- **Before — Mandatory:** "the loss is acknowledged at the scale of one gentle question"
- **After — Mandatory:** "the loss is acknowledged at the scale of one gentle question, using **at most one** of the facts listed in the scenario"
- **Before — Over-exposure:** "more than one of the six facts, or any timeline detail"
- **After — Over-exposure:** "more than one listed fact, or any timing detail beyond the month"
- **Behaviour preserved:** restraint under abundance — the item's entire purpose. Strengthened, because the boundary is now countable.
- **Mode change:** none (`R-1`). **Re-review:** **required** (boundary became scored and countable).

### 2.5 `F6-04` — Option A, confirmed (`OD-5`)

- **Defect:** `mandatory` is "nothing"; a graded `R-3` item scored by a negative requirement. Provenance, non-use, and avoidance are indistinguishable. `C-1` reproduces this mechanically.
- **Before — Mandatory:** "nothing — this item may be passed by not using the memory"
- **After — Mandatory:** "a **perceptible but unannounced** accommodation: the hot option served, offered, or ordered for them"
- **Before — Variation:** "ordering something hot for them; a small remark; no reference at all"
- **After — Variation:** "serving it; offering it; ordering it without comment; any wording that accommodates without naming the preference"
- **Before — Prohibited:** "a diagnosis; framing it as a revelation"
- **After — Prohibited:** "stating the preference; explaining why they like hot drinks; any diagnosis"
- **Before — Over-exposure:** "naming the preference unprompted, or any explanation of *why* they like hot drinks"
- **After — Over-exposure:** "naming the preference, or explaining the reason for the accommodation"
- **Behaviour preserved:** natural, restrained use of trivia. **Changed in one respect:** silence now fails, which is the point — the item can no longer be passed by a system that never retrieves.
- **Mode change:** none (`R-3` retained). **Re-review:** **required** (intended test changed).

---

## 3. Verification after application

If every proposal above is applied, the checker must report:

```text
scored units: 45
split: tuning=19 heldout=26
TOTALS MATCH
RESULT: PASS (0 fail, 0 warn, 0 unverifiable)
```

Two specific effects to confirm rather than assume:

| Check | Before | Expected after |
|---|---|---|
| `missing over-exposure` | 15 | **0** |
| `C-1 F6-04` | fail | **pass** |
| `U-1 P-01` (`the six facts`) | unverifiable | **pass** — mandatory no longer uses the generic phrase |
| `U-1 F3-03` (`the name`) | unverifiable | **pass** — mandatory names Ilsa directly |

**If the checker still reports a finding after application, the repair set is wrong and is reported as wrong.** That is the purpose of stating the prediction in advance.

---

## 4. What these proposals do not do

- They do not add items, change split assignments, or alter restraint modes.
- They do not resolve the three **rule conflicts** (`F3-01`, `F7-02`, `R6-01`) or the nine other ambiguities — those are `CORPUS-07` §2 and remain separate.
- They do not apply the `R-6` behavioural model of `CORPUS-07` §8 to any item.
- They do not seal the corpus, close `G1`, or authorise `G2`.

**Application requires separate owner approval.** Until then the originals stand unchanged and the checker's `FAIL` remains the accurate description of the corpus.
