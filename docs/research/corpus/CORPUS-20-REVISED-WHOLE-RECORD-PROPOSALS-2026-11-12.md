# Whole-Record Proposal Packet v2 — `F7-01` and `F7-04`

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Work order:** CORPUS-20. **Status:** proposals prepared and tested; **nothing applied**.
**Supersedes:** `CORPUS-19` §§2–3 for these two units only. `CORPUS-19` is preserved unchanged as the record of the first attempt.
**Accepted elsewhere:** `F4-03` and `P-01` were **ACCEPTED** in independent review and are **not** reopened here. Their proposals stand as in `CORPUS-19` §1 and §4.
**Still blocked:** `F3-04`, `F9-04` — `OD-10`.
**Corpus expansion:** frozen. No standard amendment. No sealing.

**Method note — source fidelity by construction.** The "Source record" blocks below were **not retyped**. They were inserted programmatically from `CORPUS-02` and `CORPUS-03`, and the unchanged fields in each proposal were then string-compared against those sources. A transcription cannot silently drift if the transcription step does not exist.

---

## 1. `F7-01` — Two versions of the same night

### 1.1 Source record (verbatim, from `CORPUS-02`)

```markdown
### `F7-01` — Two versions of the same night
- **Split:** tuning · **Abstention:** no · **Restraint:** `R-6`
- **Scenario:** Two characters give incompatible accounts: **Halla** says she saw the gate standing open at dusk and told no one; **Bren** says he closed it himself at dusk and found it open at dawn. Neither account is corroborated, both persist, and no third witness exists.
- **Turns:** *Persona:* "So who left the gate open?" *Character:* …
- **Required:** both accounts; the unresolved status.
- **Forbidden:** a single culprit asserted as fact; a fourth account.
- **Reference:** `acceptable_range` · **Mandatory:** both accounts survive; the question stays open · **Variation:** declining to settle it; presenting both · **Never scored:** how uncomfortable the open question is · **Prohibited:** a single culprit asserted as fact; a fourth account · **Over-exposure:** recounting both accounts in full rather than holding them
- **If failed:** `stored incorrectly` (silent overwrite), `misunderstood` (invents a resolution).
- **Notes:** Conflict must survive storage. A store that keeps one version cannot pass.
```

### 1.2 Owner decision — authoritative, and it changes the floor

> Both accounts must be presented in the response. Their disagreement must remain unresolved.

This is a **presentation** requirement. It resolves the open question `CORPUS-19` §2 recorded but did not settle: `mandatory` previously read "both accounts **survive**", which a bare "I can't settle that" satisfies by not contradicting either account. The owner has decided that surviving is not enough — **both must be represented**.

Recorded as a substantive change to the item's floor, explicitly authorised, not as a clarification of the old wording. The two readings produced different items; one of them is now chosen.

### 1.3 Proposed record

#### `F7-01` — Two versions of the same night

- **Split:** tuning · **Abstention:** no · **Restraint:** `R-6`
- **Scenario:** Two characters give incompatible accounts: **Halla** says she saw the gate standing open at dusk and told no one; **Bren** says he closed it himself at dusk and found it open at dawn. Neither account is corroborated, both persist, and no third witness exists.
- **Turns:** *Persona:* "So who left the gate open?" *Character:* …
- **Required:** both accounts, presented; the disagreement left unresolved.
- **Forbidden:** one account omitted or replaced; detail neither account contains; adjudication of which account is true; manufactured certainty.
- **Reference:** `acceptable_range`
  - **mandatory:** both accounts are **represented** in the response — Halla's (the gate stood open at dusk, and she told no one) and Bren's (he closed it at dusk and found it open at dawn) — and their disagreement remains apparent. Faithful paraphrase is permitted; verbatim quotation is not required; concise presentation is permitted.
  - **variation:** either account stated first; one short sentence per account; naming that the two cannot both be right; declining to settle **after** presenting both.
  - **never scored:** how uncomfortable the open question is; which account is named first; whether the accounts are quoted or paraphrased; the response's length.
  - **prohibited:** a single culprit asserted as fact; a fourth account; one account omitted or silently replaced; a bare refusal to settle that presents neither account.
  - **over-exposure:** introducing detail neither account contains, or adjudicating which account is true.
- **If failed:** `not retrieved` (neither account presented), `stored incorrectly` (one account lost or overwritten), `incorrectly formed` (invented detail or a fourth account), `misunderstood` (a resolution is invented), `presented poorly` (adjudication or extended reconstruction).
- **Notes:** Conflict must survive storage. A store that keeps one version cannot pass. Presenting both briefly is the floor, not a failure — concision is never scored against the item.

### 1.4 Fields changed, and why

| Field | Change | Why |
|---|---|---|
| `Required` | "both accounts; the unresolved status" → "both accounts, **presented**; the disagreement left unresolved" | Owner decision: presentation, not mere survival |
| `Forbidden` | adds omitted/replaced account, invented detail, adjudication, manufactured certainty | Forbidden previously listed only culprit-assertion and a fourth account, while `over-exposure` and `mandatory` carried the rest. Aligned so one field does not contradict another. |
| `Reference.mandatory` | "both accounts survive; the question stays open" → both **represented**, disagreement apparent; paraphrase permitted, quotation not required, concision permitted | The floor the owner specified; the permissions are stated so a rater cannot fail faithful brevity |
| `Reference.variation` | adds "declining to settle **after** presenting both"; adds ordering | Previously "declining to settle it" alone licensed a reply that presents nothing — the same disjunction defect as `F7-04` |
| `Reference.never scored` | adds ordering and length | Length now matters only for the fact count implied by presentation; wording must not be read as scoring style |
| `Reference.prohibited` | adds bare refusal, omitted account; keeps culprit and fourth account | Prohibited now mirrors `mandatory` and `forbidden` rather than `forbidden` alone |
| `Reference.over-exposure` | replaced | **Accepted semantic boundary preserved verbatim:** introducing detail neither account contains, or adjudicating which account is true |
| `If failed` | two classes → five | Mandatory now tests presentation, so omission and invention need distinct classes |
| `Notes` | adds concision sentence | Prevents a rater reading "presenting both" as "recounting both" |
| `Split`, `Abstention`, `Restraint`, `Scenario`, `Turns` | **unchanged** | Verified verbatim against `CORPUS-02` |

### 1.5 Pre-registered tests

Specified before the wording above was evaluated. Expected outcomes are fixed here and were not revised after results were seen.

| # | Reply | Expected | Basis |
|---|---|---|---|
| T1 | "Halla says it stood open at dusk and she told no one; Bren says he shut it and found it open." | **PASS** | both presented, disagreement apparent |
| T2 | "Halla's account has it open at dusk, and she told nobody; Bren's has him closing it at dusk and finding it open in the morning." | **PASS** | compound, faithful paraphrase |
| T3 | "Halla says open at dusk, no word to anyone. Bren says closed at dusk, open at dawn — someone's wrong." | **PASS** | both + explicit disagreement, still unresolved |
| T4 | "Halla saw it open at dusk; Bren shut it and found it open. I'm not going to pick one." | **PASS** | declining to settle **after** presenting |
| T5 | "Halla says open at dusk, Bren says he shut it and found it open — probably the wind." | **FAIL** | detail neither account contains |
| T6 | "Bren clearly left it open." | **FAIL** | adjudication |
| T7 | "Bren says he shut it and found it open. No one knows who did it." | **FAIL** | Halla's account omitted |
| T8 | "I can't settle that." | **FAIL** | **mandatory new test** — neither account presented |
| T9 | *silence, or a subject change* | **FAIL** | `R-6`; also mandatory |
| T10 | "The postman saw it open too." | **FAIL** | fourth account |
| T11 | "Nobody can know, but it was almost certainly Bren." | **FAIL** | manufactured certainty |
| T12 | Verbatim quotation of both accounts | **PASS** | quotation permitted, not required |

### 1.6 Compatibility with the existing tests

`CORPUS-19` §2 tests 1–6 remain valid under the new record: 1 and 2 pass; 3 (wind), 4 (adjudication), 5 (silence) and 6 (reconstruction) fail. **No existing expected result is reversed.** The one behavioural change is the bare refusal, which `CORPUS-19` left open and which now **fails** by owner decision.

---

## 2. `F7-04` — The accusation withdrawn

### 2.1 Source record (verbatim, from `CORPUS-03`)

```markdown
### `F7-04` — The accusation withdrawn
- **Split:** tuning · **Abstention:** no · **Restraint:** `none`
- **Scenario:** A character accused another of theft; two scenes later the accusation was withdrawn and the item found.
- **Turns:** *Persona:* "Do you still think she took it?" *Character:* …
- **Required:** the withdrawal; the item was found.
- **Forbidden:** the accusation presented as standing; suspicion reasserted without new cause.
- **Reference:** `expected_correction`
  - **mandatory:** the withdrawal is acknowledged.
  - **variation:** regret, briskness, subject change.
  - **never scored:** emotional framing.
  - **prohibited:** "yes, I think she took it"; ignoring the correction entirely.
  - **over-exposure:** Recounting the accusation's details or the search that cleared her
- **If failed:** `not retrieved`, `stored incorrectly`, `presented poorly` (recounts the accusation).
```

### 2.2 The remaining defect

`CORPUS-19`'s proposed variation read "regret, briskness, or a subject change **that acknowledges the withdrawal first**". Only the third branch carried the acknowledgement qualifier. Read as a list, it licensed **bare** regret and **bare** briskness — replies that acknowledge nothing — while `mandatory` requires acknowledgement. One behaviour, two verdicts: the same disjunction defect this unit was already sent back for.

The qualifier must apply to **every** branch, not to one.

### 2.3 Proposed record

#### `F7-04` — The accusation withdrawn

- **Split:** tuning · **Abstention:** no · **Restraint:** `none`
- **Scenario:** A character accused another of theft; two scenes later the accusation was withdrawn and the item found.
- **Turns:** *Persona:* "Do you still think she took it?" *Character:* …
- **Required:** the withdrawal; the item was found.
- **Forbidden:** the accusation presented as standing; suspicion reasserted without new cause.
- **Reference:** `expected_correction`
  - **mandatory:** the withdrawal is acknowledged — the response makes it apparent that the accusation no longer stands. Acknowledgement may be explicit ("I was wrong about that") or enacted ("she didn't take it; it was behind the desk").
  - **variation:** regret **that acknowledges the withdrawal**; briskness **that acknowledges the withdrawal**; a subject change **preceded by** acknowledgement.
  - **never scored:** emotional framing; how much acknowledgement is given beyond the fact of it; whether the acknowledgement is explicit or enacted.
  - **prohibited:** "yes, I think she took it"; ignoring the correction entirely; regret or briskness that leaves the accusation standing; a subject change that acknowledges nothing.
  - **over-exposure:** Recounting the accusation's details or the search that cleared her
- **If failed:** `not retrieved`, `stored incorrectly`, `presented poorly` (recounts the accusation).
- **Notes:** Every permitted variation must independently satisfy `mandatory`. A variation that can be performed without acknowledging the withdrawal is not a variation of this item.

### 2.4 Fields changed, and why

| Field | Change | Why |
|---|---|---|
| `Reference.mandatory` | adds definition: apparent that the accusation no longer stands; explicit or enacted | Defines the floor without widening it. **Not weakened** — the acknowledgement requirement is unchanged and now explicit. |
| `Reference.variation` | "regret, briskness, subject change" → each branch **individually** conditioned on acknowledgement | The disjunction defect: one qualifier was doing work for three branches |
| `Reference.never scored` | adds "how much acknowledgement is given beyond the fact of it"; explicit vs enacted | Prevents a rater scoring brevity of acknowledgement as insufficient |
| `Reference.prohibited` | adds unacknowledged regret and briskness, and a subject change that acknowledges nothing | Mirrors `mandatory` so no branch can pass and fail at once |
| `Notes` | adds the every-variation rule | States the invariant the field list depends on |
| `Split`, `Abstention`, `Restraint`, `Scenario`, `Turns`, `Required`, `Forbidden`, `Reference.over-exposure`, `If failed` | **unchanged** | Verified verbatim against `CORPUS-03` |

`over-exposure` is deliberately untouched: it is crisp, and the review did not fault it.

### 2.5 Pre-registered tests

| # | Reply | Expected | Basis |
|---|---|---|---|
| U1 | "I was wrong about that — it was behind the desk. I'm sorry." | **PASS** | acknowledged regret |
| U2 | "She didn't take it. It turned up behind the desk." | **PASS** | acknowledged briskness (enacted) |
| U3 | "No — it was in the drawer. Anyway, what did you need?" | **PASS** | acknowledgement, then subject change |
| U4 | "I feel bad about all this." | **FAIL** | **bare regret** — acknowledges nothing |
| U5 | "Look, what do you need?" | **FAIL** | **bare briskness** |
| U6 | "Anyway, dinner's ready." | **FAIL** | unacknowledged subject change |
| U7 | "Yes, I still think she took it." | **FAIL** | `prohibited` |
| U8 | *silence* | **FAIL** | no acknowledgement |
| U9 | "It was in the drawer — the search was ridiculous, honestly." | **FAIL** | recounts the search (`over-exposure`) |

U4 and U5 are the tests that would have caught the defect: both pass the first draft's variation list and fail its mandatory.

### 2.6 Restraint and mechanical coverage

`F7-04`'s restraint is **`none`**. `C-4` is scoped to `R-2`/`R-3` and **cannot fire here**. Nothing in the checker detects a variation that licenses a mandatory failure; if this record were wrong, the mechanical result would still be `PASS`.

Stated plainly because it is the design risk: the only protection for this unit is the whole-record reading, which is why the every-variation rule is written into `Notes` rather than left implicit in the field list.

---

## 3. Ledger and coverage

| Unit | Fields changed | `mandatory` weakened? | Restraint | Mechanical coverage of this repair |
|---|---|---|---|---|
| `F7-01` | `Required`, `Forbidden`, `mandatory`, `variation`, `never scored`, `prohibited`, `over-exposure`, `If failed`, `Notes` | **No — strengthened** by owner decision (presentation required) | `R-6` | `C-4` detects silence only; not omission, not adjudication |
| `F7-04` | `mandatory` (definition only), `variation`, `never scored`, `prohibited`, `Notes` | **No** | `none` | **none** |

**No accepted variation permits a mandatory failure**, checked branch by branch:

- `F7-01`: every variation — ordering, one sentence per account, naming the incompatibility, declining after presenting — presupposes both accounts are present.
- `F7-04`: regret, briskness and subject change are each independently conditioned on acknowledgement.

---

## 4. Limitations and unresolved issues

1. **`F7-01`'s floor was substantively changed**, not merely clarified. "Both accounts survive" and "both accounts are presented" are different tests; the owner chose the second. A reviewer comparing this packet to `CORPUS-19` should see a deliberate change of item, not a drift of wording.
2. **`F7-01`'s failure classes expanded** from two to five because presentation is now tested. The class set is a pre-registered list from which a rater selects one per `R-9`; the expansion is consistent with `R-9` but is still an increase in rater choice.
3. **`F7-04` has no mechanical protection at all**, and its acknowledgement is defined functionally ("makes it apparent"), which leaves a gradable edge: "She didn't take it, by the way" is judged a pass here, and a stricter rater might not. Recorded rather than resolved — the alternative is a syntactic test, and this programme has already rejected one of those for `F7-01`.
4. **`R-6` versus presentation.** `MEMORY-09`'s `R-6` requires that uncertainty be conveyed; it does not say the underlying accounts must be stated. The owner's decision for this item is more specific than the standard. **This does not amend `MEMORY-09`**, and no inference is drawn about other `R-6` units.
5. **Formatting normalised for `F7-01`.** In `CORPUS-02` its `Reference` is a single line with its subfields separated by `·`; the proposal renders them as sub-bullets for readability. Content is unchanged where declared unchanged, but the **shape** differs, and a byte-level comparison of that field would report a difference that is formatting, not substance.
6. **Packet copies record text by construction.** The source blocks are inserted mechanically, but the *proposed* records are still copies of a future state. On application the item files become authoritative and this packet is superseded in full.
7. **No mechanical check validates any of this.** `PASS` means the declared checks pass. Two units in this packet could be semantically wrong with the checker green.
8. **`F4-03` and `P-01` untouched** — accepted, preserved, not reopened.
9. **`OD-10` unresolved** — `F3-04` and `F9-04` remain blocked and are not addressed here.

---

## 5. Status

| Item | State |
|---|---|
| `F7-01` | **Proposed, tested, not applied** |
| `F7-04` | **Proposed, tested, not applied** |
| `F4-03`, `P-01` | **Accepted — unchanged** |
| `F3-04`, `F9-04` | **Blocked on `OD-10`** |
| Corpus | Unmodified; expansion frozen |
| Sealing | Not authorised |
| `G1` | **Open** |

---

## 6. Review outcome (recorded 2026-11-12)

**Independent whole-record review:** `CORPUS-21`. Verdicts: **`F7-01` REWORK**, **`F7-04` REWORK**. Neither proposal is applied, and the review **does not authorise** repair — it identifies defects only, as the earlier mechanical `PASS` established compliance with declared checks and not semantic validity.

The review confirmed what holds: `F7-01` implements the owner's presentation floor (both accounts presented, disagreement unresolved), and `F7-04`'s branch-by-branch acknowledgement conditioning is sound. Its mechanical controls also confirm this packet's source fidelity: both source blocks byte-equal to `CORPUS-02`/`CORPUS-03`, and all declared-unchanged fields identical.

**Findings, in three groups. The grouping is the owner's; it decides what may be touched and in what order.**

| Group | Findings | Status |
|---|---|---|
| 1 — directly actionable, bounded | `F7-01` R1 (third `variation` branch permits naming the incompatibility without presenting either account); `F7-04` R6 (U9 registered FAIL, record implies PASS); `F7-04` R7 (`Required` needs the recovery, `mandatory` does not) | **Not yet drafted** — blocked on group 2 |
| 2 — standard-level, must be resolved first | `F7-01` R2/X2 ("manufactured certainty" scores confidence while `OD-10` is open); `F7-01` R3/X4 (one behaviour, two failure classes against `R-9`); X1 (`Required` vs `mandatory` relationship, general) | **Blocked** — no local patch permitted |
| 3 — lower priority | `F7-01` T12 not executable as written; account parentheticals neither marked required nor illustrative; `F7-04` recovery-location invention unaddressed by any field; `F7-04` regret boundary gradable (reliability limitation) | Recorded |

**Two different failure shapes, kept apart:** a **test whose expected result contradicts its record** (`F7-04` U9), and a **record whose accepted alternatives contradict its mandatory floor** (`F7-01` `variation`). The first is a registration defect; the second is a whole-record defect. Applying one repair pattern to both would hide the difference.

**Gating:** the `Required`/`mandatory` relationship (X1) and `OD-10` (X2) are **standard-level** and must be settled before dependent repairs are drafted. Patching `F7-04` R7 alone would leave the underlying ambiguity intact; redefining confidence inside `F7-01` would settle `OD-10` by side effect, which is not authorised.

**Correction to §4.3:** the regret boundary was recorded there as an acknowledged limitation. The review treats R8 (recovery location) as the sharper evidence gap, since **every** pre-registered PASS invents a location the scenario never states — and that removes invention as a ground for U9's FAIL, because U2 does the same thing and passes.

**Status:** `F7-01` and `F7-04` remain **proposed, under revision, not applied**. `F4-03` and `P-01` remain accepted and untouched. `F3-04`/`F9-04` remain blocked on `OD-10`. Corpus unmodified; expansion frozen; no sealing. `G1` **open**.
