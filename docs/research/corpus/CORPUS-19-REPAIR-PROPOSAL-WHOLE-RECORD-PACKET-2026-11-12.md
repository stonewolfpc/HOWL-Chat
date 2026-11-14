# Repair Proposal Review Packet — whole records

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** review packet. **No proposal in this document is applied.** The corpus stands as it is; `CORPUS-02` and `CORPUS-03` are unmodified.
**Review unit:** the **complete record** for each unit, with the proposal marked inline. A field-level diff can read as correct while the record around it describes a different item — the whole-record failure that produced `F4-03` and `F7-04` in the first place.
**Excluded:** `F3-04` and `F9-04`. Their repairs depend on `OD-10` (confidence vs evidential warrant) and cannot be reviewed as items until the standard settles.
**Corpus expansion:** frozen. 34 units outside this packet untouched.

**Reviewer instructions — four questions per unit, in this order:**

1. **Whole-record:** does the record, with the proposal read as part of it, describe **one** scorable item? (Any field that licenses what another forbids is a `rework`.)
2. **Test preserved:** does the repair leave the item's **intended test** unchanged, or has it quietly become a different item?
3. **Restraint mode:** does the repair let **silence or avoidance** pass an item whose mode forbids it? `C-4` covers only `R-2`/`R-3` — on any other mode you are the only check.
4. **Checker blind spot:** if this repair were wrong, what would the mechanical checks report? "Nothing" is an acceptable and expected answer; record it.

Verdicts per `CORPUS-05` §5: **accept / rework / retire**. `retire` only for irreconcilable disagreement; a single reviewer's doubt is `rework`.

**Pre-registered examples** are numbered so a later reviewer can be scored against them. They are tests of the *boundary*, not model outputs.

---

## 1. `F4-03` — The rank

- **Split:** held-out · **Abstention:** no · **Restraint:** `none`

**Scenario.** Character was a sergeant for 100 messages, then commissioned as a lieutenant 60 messages ago. Older history still contains "Sergeant".
**Turns.** *Persona:* "Sergeant, a word?" *Character:* …
**Required.** the current rank.
**Forbidden.** the old rank accepted as current; a promotion re-narrated as news if it was already witnessed.
**Reference.** `expected_correction` — **mandatory:** the change is conveyed (currently lieutenant), without pretending it never happened. **variation:** a correction in voice; a wry "not any more"; ignoring the title and answering. **never scored:** whether the correction is gentle or blunt. **prohibited:** "yes, sergeant"; re-announcing the promotion as new information to someone who was present. **over-exposure:** Recounting the commissioning, or the circumstances of the promotion
**If failed.** `not retrieved`, `stored incorrectly` (no validity interval), `presented poorly` (recounts the commissioning).

### Proposal — `variation`

| | |
|---|---|
| **Current** | "a correction in voice; a wry 'not any more'; **ignoring the title and answering**" |
| **Proposed** | "a correction in voice; a wry 'not any more'; **an answer that conveys the current rank in passing**" |
| **Mandatory** | **unchanged** |

**Why.** "Ignoring the title and answering" satisfies neither half of mandatory: a reply that ignores "Sergeant" and just answers conveys no change. One behaviour simultaneously passes `variation` and fails `mandatory`, and matches `prohibited`'s "ignoring the correction entirely".

**Restraint mode.** `none`. `C-4` does not apply and cannot flag a regression here. The proposal keeps the obligation in the floor, so avoidance still fails mandatory.

**Checker coverage.** None. A wrong variation is invisible to every mechanical check.

**Pre-registered boundary tests**

| # | Reply | Must |
|---|---|---|
| 1 | "Not a sergeant any more — lieutenant, as of last month. What is it?" | **pass** |
| 2 | "It's lieutenant now, but go on." | **pass** |
| 3 | "Yes, sergeant?" | **fail** (`prohibited`) |
| 4 | *answers the question, never mentions rank* | **fail** (`mandatory`) |
| 5 | "That commissioning dinner was quite the evening, wasn't it?" | **fail** (`over-exposure`) |

Test 4 is the one the current variation lets through.

---

## 2. `F7-01` — Two versions of the same night

- **Split:** tuning · **Abstention:** no · **Restraint:** `R-6`

**Scenario.** Two characters give incompatible accounts: **Halla** says she saw the gate standing open at dusk and told no one; **Bren** says he closed it himself at dusk and found it open at dawn. Neither account is corroborated, both persist, and no third witness exists.
**Turns.** *Persona:* "So who left the gate open?" *Character:* …
**Required.** both accounts; the unresolved status.
**Forbidden.** a single culprit asserted as fact; a fourth account.
**Reference.** `acceptable_range` — **mandatory:** both accounts survive; the question stays open. **variation:** declining to settle it; presenting both. **never scored:** how uncomfortable the open question is. **prohibited:** a single culprit asserted as fact; a fourth account. **over-exposure:** recounting both accounts in full rather than holding them
**If failed.** `stored incorrectly` (silent overwrite), `misunderstood` (invents a resolution).
**Notes.** Conflict must survive storage. A store that keeps one version cannot pass.

### Proposal — `over-exposure`

| | |
|---|---|
| **Current** | "recounting both accounts in full rather than holding them" |
| **Proposed** | "**introducing detail neither account contains, or adjudicating which account is true**" |
| **Mandatory** | **unchanged** |

**Why.** The current boundary collides with `required`, whose content *is* both accounts. The same words that present both accounts faithfully also "recount" them. A clause-count alternative was considered and **rejected**: it fails test 2 below, accepts tests 3 and 4, and penalises ordinary paraphrase.

**Restraint mode.** `R-6` — silence cannot pass, and `C-4` **does** apply. The proposal must not create a route where saying nothing satisfies the item. It does not: silence presents neither account.

**Open question for the reviewer — not resolved here.** `mandatory` says "both accounts survive", which a bare "I can't settle that" does not do. The second reviewer treated that reply as a pass under `R-6`. Either `mandatory` requires the accounts be *presented*, or it requires only that neither be *lost* — those are different items. **Flagged, because resolving it silently would change what the item tests.**

**Checker coverage.** `C-4` catches silence-passes on `R-6`. It cannot catch the over-exposure ambiguity.

**Pre-registered boundary tests**

| # | Reply | Must |
|---|---|---|
| 1 | "Halla says it stood open at dusk; Bren says he shut it and found it open. I can't say which." | **pass** |
| 2 | "Halla's account has it open at dusk, and she told nobody; Bren's has him closing it at dusk and finding it open in the morning." | **pass** (compound, faithful) |
| 3 | "Probably the wind — Halla was half asleep anyway." | **fail** (detail neither account contains) |
| 4 | "Bren clearly did it." | **fail** (adjudication) |
| 5 | *silence, or a subject change* | **fail** (`R-6`) |
| 6 | "Picture it: dusk, the gate swinging, nobody watching." | **fail** (reconstruction) |

---

## 3. `F7-04` — The accusation withdrawn

- **Split:** tuning · **Abstention:** no · **Restraint:** `none`

**Scenario.** A character accused another of theft; two scenes later the accusation was withdrawn and the item found.
**Turns.** *Persona:* "Do you still think she took it?" *Character:* …
**Required.** the withdrawal; the item was found.
**Forbidden.** the accusation presented as standing; suspicion reasserted without new cause.
**Reference.** `expected_correction` — **mandatory:** the withdrawal is acknowledged. **variation:** regret, briskness, **subject change**. **never scored:** emotional framing. **prohibited:** "yes, I think she took it"; ignoring the correction entirely. **over-exposure:** Recounting the accusation's details or the search that cleared her
**If failed.** `not retrieved`, `stored incorrectly`, `presented poorly` (recounts the accusation).

### Proposal — `variation`

| | |
|---|---|
| **Current** | "regret, briskness, subject change" |
| **Proposed** | "regret, briskness, or a subject change **that acknowledges the withdrawal first**" |
| **Mandatory** | **unchanged** |

**Why.** Bare "subject change" permits "Anyway, dinner's ready" — no acknowledgement, which fails `mandatory` and matches `prohibited`'s "ignoring the correction entirely".

**A wider floor was proposed and withdrawn.** An earlier version widened `mandatory` to "not contradicted, and suspicion not reasserted". Because this item's restraint is **`none`**, `C-4` is silent by design, so a reply that simply avoided the topic — or said nothing — would have passed with **no mechanical finding possible**. That is the `F6-04` defect class re-created outside the checker's coverage. The floor keeps acknowledgement.

**Restraint mode.** `none` — no mechanical protection; the reviewer is the only check.

**Pre-registered boundary tests**

| # | Reply | Must |
|---|---|---|
| 1 | "No — it turned up behind the desk. Anyway, what did you need?" | **pass** (enacted acknowledgement, then subject change) |
| 2 | "That's settled. It was in the drawer." | **pass** (brisk) |
| 3 | "Anyway, dinner's ready." | **fail** (unacknowledged subject change) |
| 4 | "Yes, I still think she took it." | **fail** (`prohibited`) |
| 5 | *silence* | **fail** (`mandatory`; no acknowledgement) |

Test 3 is the one the current variation lets through. Test 1 confirms the acknowledgement need not be explicit — an enacted correction counts.

---

## 4. `P-01` — Less is better

- **Split:** tuning · **Abstention:** no · **Principle:** less-is-better · **Restraint:** `R-1`

**Scenario.** Six true facts are available about a character's recent loss: **(1)** the sibling's name, **(2)** the month of the loss, **(3)** where it happened, **(4)** who was present, **(5)** the last thing said, **(6)** the borrowed coat. Their sibling has just arrived and asked a single gentle question.
**Turns.** *Sibling:* "How are you holding up?" *Character:* …
**Required.** the loss, at the level the moment allows.
**Forbidden.** the six-fact inventory, the timeline, the names of everyone involved.
**Reference.** `acceptable_range` — **mandatory:** the loss is acknowledged — or met with a silence that registers it — at the scale of one gentle question, using **at most one** of the facts listed in the scenario. **variation:** one sentence; a silence; a gesture. **never scored:** length beyond the ceiling. **prohibited:** enumerating the facts; the timeline. **over-exposure:** more than one listed fact, or any timing detail beyond the month
**If failed.** `presented poorly` (the item is designed so that *more accurate information makes the reply worse*).
**Notes.** This is the item that a "complete memory" configuration is expected to **lose**. If no configuration can lose it, the item is badly built.

### Proposal A — `never scored`

| | |
|---|---|
| **Current** | "length beyond the ceiling" |
| **Proposed** | "**word count or verbosity, as distinct from the count of listed facts**" |

**Why.** As written it puts a failure inside the set that must never carry one: exceeding the ceiling *is* the scored failure. The charitable reading — verbosity as distinct from fact count — cannot be recovered from the words. **This defect was found by the second reviewer and missed by my own audit**, which checked `over-exposure` and not the inert field.

### Proposal B — `over-exposure`

| | |
|---|---|
| **Current** | "more than one listed fact, or any timing detail **beyond** the month" |
| **Proposed** | "more than one listed fact, or any timing detail **more precise than** the month" |

**Why.** "Beyond the month" splits between "more precise than" and "other than". Under the second reading, using the month satisfies `mandatory`'s one-fact allowance and fails `over-exposure` on the same words.

**Restraint mode.** `R-1` — silence is **permitted** and cannot be failed by `C-4` (which is scoped to `R-2`/`R-3`). This is deliberate: `mandatory` explicitly accepts "a silence that registers it", and A-1 governs at scorecard level. A future check must not "fix" this into a silence failure.

**Checker coverage.** `U-1` and presence pass regardless; neither defect is detectable mechanically.

**Pre-registered boundary tests**

| # | Reply | Must |
|---|---|---|
| 1 | "Some days are worse than others." | **pass** (no listed fact) |
| 2 | "March was hard. I'm getting through it." | **pass** (month = one fact; wordy) |
| 3 | "I don't really know how to answer that — March was hard, and I'm still working out the rest of it." | **pass** (wordy, **one** fact — tests Proposal A) |
| 4 | "March was hard, and the hospital was worse." | **fail** (two facts) |
| 5 | "Since the fourteenth." | **fail** (more precise than the month — tests Proposal B) |
| 6 | "Sit down. Tea?" | **pass or fail by A-1 only** — a gesture/silence; must not be failed by this item's `over-exposure` |

Test 6 is the restraint trap: `R-1` permits silence, so the item must not score it as under-exposure. That judgement belongs to `A-1`, one level up.

---

## Proposal ledger

| Unit | Fields changed | Mandatory | Restraint | Mechanical coverage |
|---|---|---|---|---|
| `F4-03` | `variation` | unchanged | `none` | none |
| `F7-01` | `over-exposure` | unchanged | `R-6` | `C-4` (silence only) |
| `F7-04` | `variation` | unchanged | `none` | none |
| `P-01` | `never scored`, `over-exposure` | unchanged | `R-1` | none |

**No proposal alters a `mandatory` floor.** Two earlier drafts did — `F7-04` and `F4-03` — and both were withdrawn: one because it was silence-passable outside the checker's coverage, the other because narrowing the variation preserved the item's test while widening the floor changed it.

**Nothing here is applied.** Application requires the whole-record review above, and `F3-04`/`F9-04` remain blocked on `OD-10`.

---

## Known risk in this packet

This packet **copies record text** so proposals can be read whole. A copy can drift from its source, and a drifted copy is indistinguishable from the record to a reviewer who trusts it.

Mitigations: the packet is hash-pinned in `docs/DOCUMENT-REGISTRY.md`; `CORPUS-02`/`CORPUS-03` remain authoritative; and on application the item files become the record and this packet is **superseded in full** — it must never be used to reconstruct an item. If a later reviewer finds any difference between a record here and its source, the source wins and the difference is a defect in this file, not in the corpus.
