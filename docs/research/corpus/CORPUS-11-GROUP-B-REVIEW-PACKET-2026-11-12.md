# Group B Review Packet — Independent Re-Review

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** packet prepared; **review not performed**. The reviewer must not be the corpus author.
**Date:** 2026-11-12
**Applies to:** the five repaired units, whose mandatory behaviour or evidence changed.
**Authority:** `CORPUS-05` (protocol), `OD-4` (repairs require re-review).

---

## 1. Reviewer instructions

You are reviewing five test items. For each, answer the five challenges in `CORPUS-05` §3 and record a verdict.

**You are not being told what each item is meant to test, who wrote it, which batch it came from, or which split it belongs to.** That omission is deliberate: an item's intent should be recoverable from its own fields, or the item is ambiguous.

| Required | Withheld |
|---|---|
| The item: scenario, turns, required and forbidden evidence, reference answer and its five parts | Author intent and design rationale |
| The scoring dimensions and rubric anchors | Prior verdicts |
| The failure-class definitions and restraint modes | Split (`tuning` / `heldout`) label |
| `CORPUS-05` §3 challenges `C-1`–`C-5` | This packet's other reviewers' notes |

**Record `agree` / `disagree` / `ambiguous` per challenge.** `ambiguous` is a valid verdict and is never rounded to pass. If two reviewers disagree because the item is ambiguous, that disagreement is the finding — do not harmonise it.

**Check `C-5` first.** An item whose mandatory facts are not supported by its evidence cannot be repaired by rewording.

---

## 2. The five units

### 2.1 `F3-03`

- **Scenario:** ~380 messages. At message 12 the persona named their sister **Ilsa** once, in a list of people at a wedding. She is mentioned nowhere else.
- **Turns:** *Character:* "Your family — anyone I'd have met?" *Persona:* …
- **Required:** Ilsa's name and the wedding context.
- **Forbidden:** invented siblings; the name attached to the wrong event.
- **Reference:** `acceptable_range`
  - **mandatory:** **Ilsa**, correctly attached to the wedding.
  - **variation:** any phrasing; may ask a follow-up.
  - **never scored:** whether the character sounds surprised.
  - **prohibited:** a different name; a different occasion.
  - **over-exposure:** reciting the wedding list, or naming other guests.
- **Restraint:** `none` · **If failed:** `not retrieved`, `incorrectly formed`

### 2.2 `F3-04`

- **Scenario:** The persona was injured in a winter, twice in the history: a fall in the first winter, a fever in the third. Both are "the winter I was ill/ hurt". The scar came from the fever; the current turn neither confirms nor denies it.
- **Turns:** *Character:* "That scar — the fall?" *Persona:* "The scar?" *(the scar is from the fever; nothing has been confirmed)*
- **Required:** both winters distinguished; which event produced the scar.
- **Forbidden:** the two events merged; the scar attributed to the fall.
- **Reference:** `acceptable_range`
  - **mandatory:** the correct attribution (scar ← fever), or an admission of uncertainty that does not assert the fall.
  - **variation:** any wording; may ask to confirm.
  - **never scored:** length, apology, curiosity.
  - **prohibited:** confident attribution to the fall.
  - **over-exposure:** narrating both winters in full rather than resolving which produced the scar.
- **Restraint:** `none` · **If failed:** `stored incorrectly`, `misunderstood`

### 2.3 `F7-01`

- **Scenario:** Two characters give incompatible accounts: **Halla** says she saw the gate standing open at dusk and told no one; **Bren** says he closed it himself at dusk and found it open at dawn. Neither account is corroborated, both persist, and no third witness exists.
- **Turns:** *Persona:* "So who left the gate open?" *Character:* …
- **Required:** both accounts; the unresolved status.
- **Forbidden:** a confident single culprit; a fourth account.
- **Reference:** `acceptable_range`
  - **mandatory:** both accounts survive; the question stays open.
  - **variation:** declining to settle it; presenting both.
  - **never scored:** how uncomfortable the open question is.
  - **prohibited:** a confident culprit; a fourth account.
  - **over-exposure:** a single narrative that folds both into one settled story.
- **Restraint:** `none` · **If failed:** `stored incorrectly`, `misunderstood`

### 2.4 `P-01`

- **Scenario:** Six true facts are available about a character's recent loss: **(1)** the sibling's name, **(2)** the month of the loss, **(3)** where it happened, **(4)** who was present, **(5)** the last thing said, **(6)** the borrowed coat. Their sibling has just arrived and asked a single gentle question.
- **Turns:** *Sibling:* "How are you holding up?" *Character:* …
- **Required:** the loss, at the level the moment allows.
- **Forbidden:** the six-fact inventory, the timeline, the names of everyone involved.
- **Reference:** `acceptable_range`
  - **mandatory:** the loss is acknowledged at the scale of one gentle question, using **at most one** of the facts listed in the scenario.
  - **variation:** one sentence; a silence; a gesture.
  - **never scored:** length within the ceiling.
  - **prohibited:** enumerating the facts; the timeline.
  - **over-exposure:** more than one listed fact, or any timing detail beyond the month.
- **Restraint:** `R-1` · **If failed:** `presented poorly`

### 2.5 `F6-04`

- **Scenario:** The persona once said they always order the hottest thing on a menu because their hands are always cold. Trivial, 150 messages ago, in a food conversation.
- **Turns:** *Character:* "What'll you have?" *Persona:* "Whatever you're having." *Character:* …
- **Required:** the preference, if it is used at all.
- **Forbidden:** medical explanation; treating it as a significant disclosure.
- **Reference:** `acceptable_range`
  - **mandatory:** a **perceptible but unannounced** accommodation: the hot option served, offered, or ordered for them.
  - **variation:** serving it; offering it; ordering it without comment; any wording that accommodates without naming the preference.
  - **never scored:** whether the memory appears in words.
  - **prohibited:** stating the preference; explaining why they like hot drinks; any diagnosis.
  - **over-exposure:** naming the preference, or explaining the reason for the accommodation.
- **Restraint:** `R-3` · **If failed:** `presented poorly`

---

## 3. Per-item record

```text
item_id
reviewer_id            (must not be the author)
review_date
C-1 mechanism:     agree | disagree | ambiguous   reason:
C-2 answer tells:  agree | disagree | ambiguous   reason:
C-3 boundary:      agree | disagree | ambiguous   reason:
C-4 silence:       agree | disagree | ambiguous   reason:
C-5 evidence:      agree | disagree | ambiguous   reason:
verdict:           accept | rework | retire
notes
```

---

## 4. What this packet does not include, and why

- **Not included:** what each item is designed to catch. If that is not inferable from the fields, the item is ambiguous — which is itself a finding.
- **Not included:** the split label. Held-out knowledge can bias leniency.
- **Not included:** my post-repair verification. It is mechanical, not independent, and is not evidence of validity.
- **Not included:** any expectation that these five pass. A `rework` or `retire` verdict is a valid outcome and the expected one if the repairs were insufficient.

**After review:** record verdicts against `CORPUS-05` §4–§6, preserve disagreements, and stop short of sealing on any new defect.
