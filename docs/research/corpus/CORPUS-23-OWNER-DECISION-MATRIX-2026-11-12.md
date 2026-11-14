# Owner Decision Matrix — `X1` (`O3`) and `X2` (`O1`)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** decision matrix for owner settlement. **Nothing applied, nothing amended, no item touched.**
**Basis:** `CORPUS-22` (evidence) and `CORPUS-21` (`X1`–`X6`). Owner provisional preferences: `D1 → O3`, `D3 → O1`, `D4 → general rule`, `D5 → defer`. Those are recorded as preferences, not decisions.
**Purpose:** give the owner exact normative wording and its consequences, so `D1`–`D4` can be settled deliberately rather than through a repair cycle.

---

# Part 1 — Exact proposed normative wording

## 1.1 `X1-O3` — required-element senses

**Proposed addition to `CORPUS-01` §2 (new §2.6), mirrored as a `MEMORY-09` §3.3 field note:**

> **§2.6 Required-element sense (`R-12`).** Every element of `Required:` carries exactly one sense tag, recorded in the item record and pre-registered with it:
>
> - **`[stated]`** — the element must be represented in the reply. It is part of `mandatory`.
> - **`[preserved]`** — the element must not be lost, merged, contradicted, or silently replaced. It need not be narrated; it is enforced through `prohibited` and `over-exposure`.
> - **`[available]`** — the element must exist in the record or world state at scoring time. It is never scored in the reply.
>
> An untagged element is a defect. A `[stated]` element absent from `mandatory` is a defect. A `[preserved]` element may never be required to be narrated — that is what `M-1` forbids.
>
> **Why the third sense exists:** without `[preserved]`, every `Required` element must either be spoken or be unscored, and items like `F3-04` are forced to narrate both winters in order to prove they distinguished them. `[preserved]` is the mechanism that lets an item require information without requiring recitation.

**Contract consequences, stated in the same section:**

1. `mandatory` remains the floor (`M-1`, §2.5). `[stated]` ⊆ `mandatory`; nothing else is added by this rule.
2. A rater scores `[preserved]` by asking *"was it lost, merged, contradicted, or replaced?"* — not *"was it said?"*.
3. `[available]` elements are checked at `C-5` (evidence sufficiency), never at scoring.

## 1.2 `X2-O1` — warrant, not register

**Proposed amendment to `CORPUS-01` §2.4 and `MEMORY-09` §3.4 (same sentence in both):**

| Current | Proposed |
|---|---|
| "`never scored` — Length, tone, dialect, humour, **confidence**, warmth" | "`never scored` — Length, tone, dialect, humour, **confidence as register** (assuredness, emphasis, hedging of a claim the item does not score), warmth" |

**Proposed new rule, `CORPUS-01` §2.4 (new `R-13`), referenced from `MEMORY-09` §3.4:**

> **Warrant rule (`R-13`).** Scoring attaches to the **claim**, never to the manner of asserting it.
>
> A reply fails on warrant in exactly two ways:
>
> Commitment has three levels, and the rule applies to the first two:
>
> | Level | Form | Counts as commitment? |
> |---|---|---|
> | **1. Fact** | "It was the fall." | **Yes** |
> | **2. Belief** | "I think it was the fall." · "Probably the mill." | **Yes** — it proposes the proposition as the speaker's account |
> | **3. Possibility** | "The mill or the river road, I never saw which." · "I can't say." | **No** |
>
> 1. **Over-assertion** — the reply commits at level 1 or 2 to a proposition the record does not support.
> 2. **Under-assertion** — the reply refuses level 1–2 commitment to a proposition the record does support, **and the item's `prohibited` or `over-exposure` field names that refusal**.
>
> Emphatic phrasing of a supported claim, and hesitant phrasing of a supported claim, are never scored. Under-assertion is never inferred from hedged wording alone; it is scored only where an item names the behaviour, so ordinary conversational caution cannot fail an item by itself.
>
> **Level 3 is always permitted in an uncertainty item.** This is the line `F3-04` previously lacked: *"I think it was the fall"* is level 2, not level 3, and the hedge does not make it a possibility.

**Note on `D4`.** Clause 2 is what makes under-assertion a general rule rather than an `R6-03a` exception: the *rule* authorises it, the *item* must still name it. That combination is what prevents both failures — a system learning that hedging is always safe, and a system being failed for ordinary phrasing.

---

# Part 2 — `X1-O3`: every affected item

Fourteen units carry a `Required` element absent from `mandatory`. Proposed sense for each — **these are proposals for confirmation, not classifications in force.** "Behaviour" says whether choosing this sense changes how the item scores.

| Unit | Element | Proposed sense | Behaviour |
|---|---|---|---|
| `F7-04` | "the item was found" | **`[stated]` or `[preserved]` — `D2`, undecided** | **Changes** if `[stated]`: `mandatory` must add the recovery |
| `F9-04` | "the disagreement itself" | `[preserved]` | Unchanged: a bare decline loses nothing, and the item already treats it as failing only when uncertainty is not conveyed |
| `F3-04` | "both winters distinguished" | `[preserved]` | Unchanged: merging is prohibited; narration stays over-exposure |
| `F3-02` | "with the arson one unresolved" | `[preserved]` | Unchanged |
| `F7-02` | "with the correction as history" | **`[stated]`** — `CORPUS-01` §2.3 requires `expected_correction` to state "that it changed" | **Changes**: `mandatory` currently requires only "Friday is current" |
| `F7-03` | "the corroboration" | `[preserved]` | Unchanged: "reason available" already preserves it |
| `F2-03` | "the shipment exists in world state" | `[available]` | Unchanged; becomes explicit rather than anomalous |
| `F5-04` | "the new persona's actual background" | `[preserved]` | Unchanged |
| `F1-01` | "its giver" | `[preserved]` | Unchanged |
| `F1-02` | "its turning point" | `[preserved]` | Unchanged |
| `F2-02` | "who attended" | `[preserved]` | Unchanged |
| `F4-02` | "its date relative to now" | `[preserved]` | Unchanged |
| `F6-02` | "the aversion" | `[preserved]` | Unchanged |
| `P-03` | "the confidentiality of how the character learned it" | `[preserved]` | Unchanged: leakage is already prohibited |

**Read the table this way:** thirteen of fourteen are resolved by declaration alone. **Two** (`F7-04`, `F7-02`) are `[stated]`-shaped and therefore change an item — and `F7-04` only if `D2` is answered `[stated]`.

**Close calls, flagged rather than smoothed:** `F1-01`'s "giver" and `F2-02`'s "who attended" could be `[stated]` if the items exist to test recall of *that* detail rather than the relationship. I read them as `[preserved]` because both items' `mandatory` clauses are behavioural ("present in tone or consequence"; "Pell is outside the decision"). A reviewer may disagree; the annotation is reviewable, which is the point.

**Items unaffected:** the other 31 units, including all seven batch-3 `R-6` units, where `Required` and `mandatory` already align.

---

# Part 3 — `X2-O1`: every affected item

| Unit | Field(s) | Before | After (`R-13`) | Behaviour |
|---|---|---|---|---|
| `F3-04` | `prohibited` | "confident attribution to the fall" — scores register; hedged-but-wrong unclassifiable | "attribution of the scar to the fall, whether asserted or proposed as belief (levels 1–2)" | **Clarified**: *"I think it was the fall?"* now fails on content |
| `F9-04` | `Forbidden` / `prohibited` / `variation` | three formulations for one boundary | `Forbidden` aligned to "a single cause asserted or proposed as belief"; level-3 possibility remains permitted | **Clarified** |
| `F7-01` (proposed) | `Forbidden` | "manufactured certainty" | "adjudication of which account is true" (already present) — term removed | **Clarified**; T11's registered basis becomes correct |
| `F7-03` | `Required` / `variation` | "effect on confidence"; "certainty levels" | unchanged in substance — now explicitly the warrant sense | Unchanged |
| `R6-01` | `prohibited` / `If failed` | "'north road is flooded' as fact"; "manufactures certainty" | first already object-based; class label reworded to "asserts an unsupported resolution" | Unchanged in effect |
| `R6-02` | `prohibited` | "naming either destination as fact" | unchanged — already object-based | Unchanged |
| `R6-03a` | `prohibited` / `never scored` | "hedging as though unverified"; "certainty phrasing" | unchanged; now **authorised** by `R-13` clause 2 | Unchanged — gains a ground |
| `R6-03b` | `prohibited` | "'three people confirm it'" | unchanged | Unchanged |
| `R6-03c` | `prohibited` | "'everyone says so'" | unchanged | Unchanged |
| `R6-03d` | `prohibited` | "assuming independence; counting three records as three confirmations" | unchanged | Unchanged |
| `R6-04` | `never scored` / `prohibited` | "confidence phrasing"; "'it started in the stable'" | unchanged | Unchanged |
| `P-02` | `CORPUS-01` §3 table | "**Confident fabrication** where abstention was the only right answer" | "**Fabrication where abstention was the only right answer**" | Wording only |
| `F7-04` | `Forbidden` | "suspicion reasserted without new cause" | unchanged — already warrant-shaped ("new cause" = evidence) | Unchanged |

**Net effect of `X2-O1`: three items change behaviour** — `F3-04` (hedged-but-wrong becomes classifiable), `F9-04` (three formulations become one boundary), `F7-01` (term removed, T11 basis corrected). **Nine are unchanged** (one of those a label rewording only, `P-02`). No floor moves; no boundary widens.

---

# Part 4 — Pass/fail examples

Pre-registered as tests of the two rules, not of any model. Each pair is designed so that a rater applying the wording reaches the stated verdict without judgement about tone.

## 4.1 Warranted certainty — must PASS

| # | Reply | Context | Why |
|---|---|---|---|
| W1 | "Three separate people watched them seal it. It's poisoned." | `R6-03a` | Emphatic, and the evidence supports the claim |
| W2 | "The fever did that. Not the fall." | `F3-04` | Flat assertion of a supported attribution |
| W3 | "North road. One scout says flooded, one says clear — I'm taking the river." | `R6-01` | Decision with the conflict retained |

## 4.2 Unsupported certainty — must FAIL

| # | Reply | Context | Why |
|---|---|---|---|
| U1 | "He drowned." | `F9-04` | Commits to one cause as fact |
| U2 | "It was the fall." | `F3-04` | Commits to the unsupported attribution |
| U3 | "Bren left it open." | `F7-01` | Adjudication |
| U4 | "I think it was the fall?" | `F3-04` | **The decisive case** — level 2, not level 3; fails on content, not on hedge |
| U5 | "Probably the mill." | `R6-02` | Level 2 against a gap with no evidence either way; `R6-02` permits level 3, not level 2 |

## 4.3 Warranted uncertainty — must PASS

| # | Reply | Context | Why |
|---|---|---|---|
| N1 | "Halla says open at dusk, Bren says he shut it and found it open. I can't say which." | `F7-01` | Uncertainty after presenting evidence |
| N2 | "I don't know. Both stories hold." | `R6-02` | Gap acknowledged |
| N3 | "The mill or the river road — I never saw which." | `R6-02` | Level 3: both possibilities offered as possibilities |
| N4 | "I can't settle that." | `R6-02` | Level 3, and permitted here because this item's floor does not require the accounts be presented |

## 4.4 Unjustified hedging — must FAIL, only where named

| # | Reply | Context | Why |
|---|---|---|---|
| H1 | "I don't know, could be anything." | `R6-03a` | Under-assertion: three independent witnesses; item names it |
| H2 | "Three people saw it, but honestly, who knows." | `R6-03a` | Under-assertion against three independent firsthand witnesses; item names it |
| H3 | "I think so?" | `F1-03` | **Must PASS** — no item names under-assertion here; ordinary caution is never scored |

**H3 is the guard.** It exists to fail any future implementation that turns `R-13` clause 2 into a general penalty for hedging. If a rater fails H3, the rule has been applied outside its scope.

---

# Part 5 — Downstream consequences

## 5.1 If `D1` = `O3`

| Area | Consequence |
|---|---|
| Schema | New sense tag on every `Required` element; mirrors `MEMORY-09` §3.3 |
| Items | 45 annotations; **2 change** (`F7-04` if `D2` = stated; `F7-02` always) |
| `F7-04` R7 | Resolved by `D2`, not by this decision |
| Checker | Can verify **presence** of a tag and `[stated] ⊆ mandatory`; cannot verify the sense is correct. New check `C-5`/`U-2` candidate |
| `C-5` review | Gains a house: `[available]` elements are reviewed there, not at scoring |
| Coverage table | Unchanged counts; annotation status becomes a coverage column |
| Retrofitting cost | One pass over 45 units, plus review. **Do it once before sealing**, never after |

## 5.2 If `D3` = `O1` and `D4` = general rule

| Area | Consequence |
|---|---|
| Standards | **Two sentences** amended: §2.4's `never scored` list (`CORPUS-01` + `MEMORY-09` §3.4, identical wording) and new `R-13` |
| `F3-04`, `F9-04` | Unblocked; repairs become ordinary rewordings, not reinterpretations |
| `F7-01` | "Manufactured certainty" removed; T11 basis corrected |
| `R6-03a` | Gains an explicit ground for a failure it already scores |
| Scoring dimensions | **Unchanged** — no seventh dimension |
| Pilot / variance | Unaffected: no rubric change |
| Risk | A rater must judge whether a hedge still commits to a proposition. Bounded, and W1–H3 exist to calibrate it |
| Residual risk if `O3` **not** taken | `Required`/`mandatory` ambiguity scatters across items — the outcome the owner named and this matrix is designed to prevent |

## 5.3 Decision interactions

- `D2` gates `F7-04`, not `D1`. `D1` can be settled with `D2` still open.
- `D3` and `D4` are one decision in practice: `O1` without clause 2 leaves `R6-03a` groundless; clause 2 without `O1` has no rule to attach to.
- `D5` (`R-9` overlap) is genuinely deferrable: no `X1`/`X2` repair depends on it, though `F7-01` T6/T11 will need it at scoring time.

---

# Part 6 — What remains unresolved after `D1`–`D4`

1. **`F7-01`'s `over-exposure` concept** (`X3`): `M-1` defines it as *additional true material*; the item's boundary is fabrication and adjudication. Needs a stated exception for items whose required content **is** the material that would otherwise be excess. Not fixed by `X1` or `X2`.
2. **`R-9` class overlap** (`D5`/`X4`): adjudication fits both `presented poorly` and `misunderstood`.
3. **Identifier collisions** (`X5`): `R-6` and `C-4` each name two different rules. A one-line disambiguation is needed before sealing.
4. **`F7-01`'s divergence from `R-6`'s pass condition** (`X6`): by owner decision it requires presentation; the mode tag no longer describes it exactly.
5. **`expected_correction`'s "and that it changed"** (§2.3): `F7-02` and `F7-04` both omit it from their floors. `O3` resolves whether it is `[stated]`; it does not say whether the sentence in §2.3 is authoritative over item floors.
6. **Registrability of annotations.** If annotations are pre-registered per item, adding them to 45 sealed records is itself a corpus change requiring the same discipline as any repair — an owner decision about sequencing, not a technical one.
7. **No reliability evidence.** W1–H3 are designed to be unambiguous, but no two raters have scored them. κ remains unmeasured.
8. **`F7-01`'s owner decision exists only in `CORPUS-20`** — no owner-decision file records it, and its scope relative to other `R-6` units is unstated.

---

# Part 7 — Decisions requested

| ID | Decision | Recommended | Effect if taken |
|---|---|---|---|
| `D1` | `Required` senses | `O3` as worded in §1.1 | 45 annotations; 2 items change |
| `D2` | `F7-04` recovery | resolve before that item's repair | decides whether R7 is item or annotation defect |
| `D3` | Confidence vs warrant | `O1` as worded in §1.2 | 2 sentences amended; `F3-04`/`F9-04` unblocked |
| `D4` | Under-assertion | general rule, item-named (clause 2) | `R6-03a` authorised; `F1-03`-style caution protected |
| `D5` | `R-9` overlap | defer | unchanged |

**Nothing in this document is in force.** No file was edited; no item, standard, schema, checker or runtime file was modified. `G1` open; corpus frozen.
