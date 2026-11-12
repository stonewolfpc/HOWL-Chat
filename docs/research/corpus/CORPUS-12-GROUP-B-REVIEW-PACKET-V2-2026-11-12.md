# Group B Review Packet — Version 2 (Corrected Commission)

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** packet **v2 prepared and commissioned**; review **performed** — see `CORPUS-13`.
**Date:** 2026-11-12
**Supersedes for commissioning purposes:** `CORPUS-11-GROUP-B-REVIEW-PACKET-2026-11-12.md` — **preserved byte-identical**, not replaced. Its hash is unchanged: `b51c4b680dd1f300bdf659228ed4def462f0047ab8fc3ebe39fa419ea7d826aa`.
**Does not supersede:** `CORPUS-11`'s existence, status, or secondary-evidence standing. `CORPUS-11` remains the record of what was *proposed* to be commissioned; this file is the record of what was *actually* issued.
**Authority:** `CORPUS-05` (protocol), `CORPUS-07` `OD-4`, work order `CORPUS-11` §1 (packet versioning) and §2 (blindness).

---

## 1. Why a version 2 exists

Version 1 could not be used unamended. Three defects were found in it **before** any verdict was collected for four of the five items, and one was found after the first item had already been reviewed. All three are packet defects, not item defects. None of them changed a corpus byte.

| # | Packet defect | Evidence | Consequence |
|---|---|---|---|
| `PV-1` | `F7-01` printed `Restraint: none`; the item registers `R-6` | `CORPUS-02` line 160: `- **Split:** tuning · **Abstention:** no · **Restraint:** \`R-6\`` | Flips the silence test. `R-6` is the one mode where **saying nothing also fails**; `none` is not a restraint mode at all. A reviewer briefed `none` would have scored `C-4` wrong. |
| `PV-2` | `F7-01` over-exposure paraphrased as "a single narrative that folds both into one settled story"; the item registers "a detailed reconciliation attempt that settles it" | `CORPUS-02` line 160 | Narrower than the registered boundary. A reconciliation attempt that fails to settle would have been read as permitted. |
| `PV-3` | The briefing omitted the restraint taxonomy and the meaning of `Restraint: none` | round-1 briefing issued to the `F3-03` reviewer | The reviewer reported `Restraint: none` + an over-exposure clause as a **field contradiction** and made it the blocking defect. The registered rules resolve that combination explicitly. `F3-03` was therefore re-reviewed under the corrected briefing; both rounds are preserved (`CORPUS-13` §3). |

**Correction applied:** v2 reproduces every field from the corpus files, states the taxonomy and the `none` interpretation rule, and is the text actually issued for `F3-04`, `F7-01`, `P-01`, `F6-04`, and the `F3-03` re-review.

**Not corrected, and deliberately so:** `CORPUS-11`'s body still contains `PV-1` and `PV-2`. Rewriting it would destroy the record of what was reviewed and when. The discrepancy is registered instead (`CORPUS-14` `ND-1`).

### 1.1 Field-fidelity check — v1 print vs registered value

Performed against the working-tree corpus files, not against `CORPUS-08`'s proposals.

| Item | Field | `CORPUS-11` printed | Registered (corpus file) | Match |
|---|---|---|---|---|
| `F3-03` | Restraint | `none` | `none` | ✓ |
| `F3-03` | over-exposure | "reciting the wedding list, or naming other guests" | "Reciting the wedding list, or naming other guests" | ✓ (case only) |
| `F3-04` | Restraint | `none` | `none` | ✓ |
| `F3-04` | over-exposure | "narrating both winters in full…" | "Narrating both winters in full rather than resolving which produced the scar" | ✓ (case only) |
| `F7-01` | **Restraint** | **`none`** | **`R-6`** | ✗ **`PV-1`** |
| `F7-01` | **over-exposure** | **"a single narrative that folds both into one settled story"** | **"a detailed reconciliation attempt that settles it"** | ✗ **`PV-2`** |
| `P-01` | Restraint | `R-1` | `R-1` | ✓ |
| `P-01` | over-exposure | "more than one listed fact, or any timing detail beyond the month" | identical | ✓ |
| `F6-04` | Restraint | `R-3` | `R-3` | ✓ |
| `F6-04` | over-exposure | "naming the preference, or explaining the reason for the accommodation" | identical | ✓ |

One material mismatch (`F7-01` restraint tag), one paraphrase (`F7-01` over-exposure), eight matches.

---

## 2. What the reviewer receives — and what is withheld

| Provided | Withheld |
|---|---|
| The item: scenario, turns, required and forbidden evidence, reference answer and all five parts | Author intent and design rationale |
| The restraint taxonomy `R-1`–`R-6` and the item's tag | Batch identity |
| `M-1`, `M-2`, the floor-not-checklist rule, the never-scored rule | Tuning / held-out split |
| The seven challenges and the silence test | Prior verdicts, including `CORPUS-06`'s |
| The verdict convention (`agree` / `disagree` / `ambiguous`) | DeepSeek's post-repair verification and every expected outcome |
| The corpus-wide prose convention for evidence sufficiency | `CORPUS-08`'s before/after rationale and the repair manifest |

**No evidence was removed to achieve blindness.** The five items are reproduced complete — every field, including `If failed` and `Notes`. Blindness is achieved by withholding *interpretation*, never by withholding *material*.

### 2.1 Rules text issued verbatim to every reviewer

```text
- M-1: every item states its over-exposure boundary (additional TRUE material that would still fail it).
- M-2: for restraint items over-exposure is a scored part, not a note. A reply can satisfy every
  required fact and still fail by saying too much.
- mandatory is a FLOOR, not a checklist to recite. 'Never scored' parts (length, tone, dialect,
  humour, confidence, warmth) are NEVER failure signals.
- Restraint taxonomy, six modes, each independent:
  R-1 suppress mention (knows the fact, chooses not to name it; nothing implies it was consulted;
      silence CAN pass)
  R-2 influence tone (knowledge changes manner, not content; silence CANNOT pass)
  R-3 natural use (fact used as though simply remembered; silence CANNOT pass)
  R-4 recognise irrelevance (no callback; silence CAN pass)
  R-5 respect entitlement (withholds because recipient may not know; silence CAN pass)
  R-6 preserve uncertainty (refuses to manufacture certainty; silence CANNOT pass)
- CRITICAL INTERPRETATION RULE: 'Restraint: none' means the item is NOT a restraint-mode item at all.
  It does NOT mean over-exposure is absent or unscored. An item tagged 'Restraint: none' that also
  carries an over-exposure boundary is explicitly internally consistent and M-1 compliant.
  Do not report that combination as a contradiction.
- Corpus-wide convention: scenarios are synthetic specification prose, not verbatim transcript dumps.
  'Evidence sufficiency' means the scenario states the fact clearly enough that a reply can be judged
  against it; it does not require a quoted transcript line unless the item depends on exact wording.
- Seven challenges: (1) evidence sufficiency; (2) mandatory behaviour; (3) variation; (4) prohibited
  behaviour; (5) over-exposure; (6) scoring clarity; (7) test validity.
  Also assess: can the item be passed by blanket silence or blanket abstention?
- Verdicts must be exactly one of: agree / disagree / ambiguous.
  Ambiguity is NEVER rounded up to agreement or a pass.
  Cite the exact field, evidence passage, or registered rule supporting each judgment.
```

---

## 3. The five units, as issued

Field text below is copied from the corpus files. `CORPUS-12` is a transcript, not a summary.

### 3.1 `F3-03` — The name at message 12

- **Scenario:** ~380 messages. At message 12 the persona named their sister **Ilsa** once, in a list of people at a wedding. She is mentioned nowhere else.
- **Turns:** *Character:* "Your family — anyone I'd have met?" *Persona:* …
- **Required:** the sister's name and the wedding context.
- **Forbidden:** invented siblings; the name attached to the wrong event.
- **Reference:** `acceptable_range`
  - **mandatory:** **Ilsa**, correctly attached to the wedding.
  - **variation:** any phrasing; may ask a follow-up.
  - **never scored:** whether the character sounds surprised.
  - **prohibited:** a different name; a different occasion.
- **over-exposure:** Reciting the wedding list, or naming other guests
- **Restraint:** `none` · **If failed:** `not retrieved` (no name), `incorrectly formed` (wrong name or event).

**Item-specific scrutiny issued:** verify that Ilsa is explicitly established by the supplied evidence; check that the repair did not introduce the name as an unsupported fact or make the answer available in the current turn.

### 3.2 `F3-04` — Two winters

- **Scenario:** The persona was injured in a winter, twice in the history: a fall in the first winter, a fever in the third. Both are "the winter I was ill/ hurt". The scar came from the fever; the current turn neither confirms nor denies it.
- **Turns:** *Character:* "That scar — the fall?" *Persona:* "The scar?" *(the scar is from the fever; nothing has been confirmed)*
- **Required:** both winters distinguished; which event produced the scar.
- **Forbidden:** the two events merged; the scar attributed to the fall.
- **Reference:** `acceptable_range`
  - **mandatory:** the correct attribution (scar ← fever), or an admission of uncertainty that does not assert the fall.
  - **variation:** any wording; may ask to confirm.
  - **never scored:** length, apology, curiosity.
  - **prohibited:** confident attribution to the fall.
- **over-exposure:** Narrating both winters in full rather than resolving which produced the scar
- **Restraint:** `none` · **If failed:** `stored incorrectly` (merged), `misunderstood`.

**Item-specific scrutiny issued:** verify that the revised turn "The scar?" no longer gives away the answer; check that long-range retrieval remains necessary and that the item has not become a current-turn comprehension test.

### 3.3 `F7-01` — Two versions of the same night

- **Scenario:** Two characters give incompatible accounts: **Halla** says she saw the gate standing open at dusk and told no one; **Bren** says he closed it himself at dusk and found it open at dawn. Neither account is corroborated, both persist, and no third witness exists.
- **Turns:** *Persona:* "So who left the gate open?" *Character:* …
- **Required:** both accounts; the unresolved status.
- **Forbidden:** a confident single culprit; a fourth account.
- **Reference:** `acceptable_range`
  - **mandatory:** both accounts survive; the question stays open.
  - **variation:** declining to settle it; presenting both.
  - **never scored:** how uncomfortable the open question is.
  - **prohibited:** a confident culprit; a fourth account.
- **over-exposure:** a detailed reconciliation attempt that settles it
- **Restraint:** `R-6` · **If failed:** `stored incorrectly` (silent overwrite), `misunderstood` (invents a resolution).

**Item-specific scrutiny issued:** verify that both Halla's and Bren's accounts are supplied, distinguishable, and sufficient to evaluate whether both survive; check that the required response does not force an unsupported reconciliation or imply one account must be true.

### 3.4 `P-01` — Less is better

- **Principle:** less-is-better · **Restraint:** `R-1` · **Abstention:** no
- **Scenario:** Six true facts are available about a character's recent loss: **(1)** the sibling's name, **(2)** the month of the loss, **(3)** where it happened, **(4)** who was present, **(5)** the last thing said, **(6)** the borrowed coat. Their sibling has just arrived and asked a single gentle question.
- **Turns:** *Sibling:* "How are you holding up?" *Character:* …
- **Required:** the loss, at the level the moment allows.
- **Forbidden:** the six-fact inventory, the timeline, the names of everyone involved.
- **Reference:** `acceptable_range`
  - **mandatory:** the loss is acknowledged at the scale of one gentle question, using **at most one** of the facts listed in the scenario.
  - **variation:** one sentence; a silence; a gesture.
  - **never scored:** length beyond the ceiling.
  - **prohibited:** enumerating the facts; the timeline.
- **over-exposure:** more than one listed fact, or any timing detail beyond the month
- **If failed:** `presented poorly`.

**Item-specific scrutiny issued:** verify that all six facts are enumerated and supported; check that the countable boundary measures restraint rather than rewarding arbitrary omission; confirm that a response meeting the mandatory acknowledgment still fails if it exceeds the boundary; note any tension between `variation: a silence` and the mandatory acknowledgment.

### 3.5 `F6-04` — Cold hands, warm drink

- **Restraint:** `R-3` · **Abstention:** no
- **Scenario:** The persona once said they always order the hottest thing on a menu because their hands are always cold. Trivial, 150 messages ago, in a food conversation.
- **Turns:** *Character:* "What'll you have?" *Persona:* "Whatever you're having." *Character:* …
- **Required:** the preference, if it is used at all.
- **Forbidden:** medical explanation; treating it as a significant disclosure.
- **Reference:** `acceptable_range`
  - **mandatory:** a **perceptible but unannounced** accommodation: the hot option served, offered, or ordered for them.
  - **variation:** serving it; offering it; ordering it without comment; any wording that accommodates without naming the preference.
  - **never scored:** whether the memory appears.
  - **prohibited:** stating the preference; explaining why they like hot drinks; any diagnosis.
- **over-exposure:** naming the preference, or explaining the reason for the accommodation
- **If failed:** `presented poorly` (over-weighted trivia).
- **Notes:** Deliberately permits non-use. Trivia that is available need not be used, and using it heavily is the failure.

**Item-specific scrutiny issued:** verify the natural-use behaviour is genuine rather than coincidental — perceptible accommodation required, preference unannounced, silence or non-use fails, multiple natural expressions pass, no scripted wording required; note the tension between `Required` and the mandatory accommodation.

---

## 4. Per-item record format issued

```text
item_id
C-1 evidence sufficiency : agree | disagree | ambiguous   reason:
C-2 mandatory behaviour  : agree | disagree | ambiguous   reason:
C-3 variation            : agree | disagree | ambiguous   reason:
C-4 prohibited behaviour : agree | disagree | ambiguous   reason:
C-5 over-exposure        : agree | disagree | ambiguous   reason:
C-6 scoring clarity      : agree | disagree | ambiguous   reason:
C-7 test validity        : agree | disagree | ambiguous   reason:
silence test             : agree | disagree | ambiguous   reason:
verdict                  : accept | rework | retire
```

**Challenge numbering note.** `CORPUS-05` §3 numbers its five challenges `C-1`–`C-5` with a different meaning (`C-1` mechanism, `C-2` answer tells, `C-3` boundary, `C-4` silence, `C-5` evidence). This packet uses `C-1`–`C-7` for the seven review questions required by the work order, which is a superset. Where a verdict needs to be mapped back to `CORPUS-05`, the mapping is: packet `C-1`→`C-5`, `C-2`/`C-3`/`C-4`→`C-2`, `C-5`→`C-3`, silence test→`C-4`, packet `C-6`/`C-7`→no `CORPUS-05` equivalent. **Registered as `ND-10`** — two challenge numbering schemes now exist in the corpus records, which is a naming hazard.

---

## 5. Boundaries observed in this packet version

- No corpus item was edited. No field was corrected.
- `CORPUS-11` was not rewritten, replaced, or deleted; its hash still matches the registry.
- `MEMORY-09` was not touched.
- `check_corpus.py` was not modified and was not re-run to make a finding disappear.
- No item added, removed, split, or retagged.
- No model loaded, no experiment run.
- `G1` remains open. This packet seals nothing.

**Next:** verdicts in `CORPUS-13`; reconciliation, defect register, and escalation in `CORPUS-14`.
