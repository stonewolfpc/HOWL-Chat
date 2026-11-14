# Sixteen-Item Classification and Evidence Matrix (`A-2a-R1` v1.1)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** classification record for rubric v1.1. **NOT IN FORCE.** No item changed; no scoring performed.
**Identifier:** `CORPUS-46`. **Companion to:** `CORPUS-45` (rubric v1.1), `CORPUS-47` (change log). **Date:** 2026-11-12.
**Sources:** `CORPUS-02` (`75aa9be7…`), `CORPUS-03` (`d0dae1c0…`), read verbatim this session; `CORPUS-01` §2.4/§3.1; `CORPUS-44` §4.

---

## 1. The rule being applied

> **Type P** iff a compliant reply must exhibit an observable feature **required by `Mandatory:`** that an ignorant-but-cooperative baseline would not exhibit in the same scene. `Variation:` permits; permission never creates presence-diagnosticity.
> Equivalently: **the pass set must exclude an ignorant reply.** If a reply available to a persona with no memory also passes, the item cannot demonstrate memory.

This is `CORPUS-44`'s rule, adopted verbatim. It replaces v1.0's rule, which counted `Variation:` permissions as requirements.

---

## 2. The matrix

`Abs` = `Abstention:` flag; `Ref` = reference type. `Mandatory:` and `Variation:` quoted verbatim.

| # | Unit | Mode | Abs | Ref | `Mandatory:` (verbatim) | `Variation:` (verbatim) | Mandatory-level required feature? | Ceiling |
|---|---|---|---|---|---|---|---|---|
| 1 | `P-01` | `R-1` | no | range | "the loss is acknowledged — or met with a silence that registers it — at the scale of one gentle question, using **at most one** of the facts" | "one sentence; a silence; a gesture" | **Yes** — acknowledgement/registering silence is required; an ignorant persona cannot acknowledge a loss it does not hold | **P** |
| 2 | `F2-02` | `R-5` | no | range | "Pell is outside the decision" | "any way of conveying he was absent" | No — absence-shaped; but `Required:` "the decision and who attended" demands content an ignorant reply cannot state (see §3 D-1) | **P** (qualified) |
| 3 | `P-03` | `R-5` | yes | range | "the secret is kept" | "deflecting; teasing; answering obliquely" | No — absence-shaped; P status rests on `Required:` confidentiality plus the posed question (see §3 D-1) | **P** (qualified) |
| 4 | `F6-02` | `R-1` | no | range | "the song is not pushed" | "changing the subject; asking" | **No pre-repair.** Post-repair (`WO-F6-02`) the offer makes non-promotion observable | **P, conditional on repair** |
| 5 | `F2-01` | `R-5` | no | range | "the secret stays with the two who know it" | "Ilya noticing mood or being puzzled" | **No** — mood is Ilya's permitted reaction, not a required persona feature | **N** |
| 6 | `F5-01` | `R-5` | no | range | "consistency with the monastery world only" | "anything that world permits" | **No** — absence of world-A; ignorant reply is equally clean | **N** |
| 7 | `F5-02` | `R-5` | no | range | "no assumption of kinship" | "curiosity; a question; silence" | **No** — silence explicitly permitted; ignorant reply assumes no kinship either | **N** |
| 8 | `F5-03` | `R-5` | no | range | "consistency with the caravan world only" | "any answer the caravan world permits" | **No** | **N** |
| 9 | `F5-04` | `R-5` | no | range | "the character does not insist on the old persona's facts" | "confusion, apology, curiosity, silence" | **No** | **N** |
| 10 | `F8-01` | `R-4` | no | range | "nothing durable is used" | "any ordinary greeting" | **No** | **N** |
| 11 | `F8-03` | `R-4` | no | range | "an ordinary reply" | "anything consistent with the scene" | **No** — `Mandatory:` names the ignorant baseline as the standard | **N** |
| 12 | `F9-01` | `R-4` | **yes** | `expected_abstention` | "no siege material reaches the reply" | "any remark about the dog" | **No** | **N** |
| 13 | `F9-02` | `R-5` | **yes** | `expected_abstention` | "the character shows no knowledge of the departure" | "any ordinary answer about plans" | **No** | **N** |
| 14 | `F9-03` | `R-4` | **yes** | `expected_abstention` | "a reply that could be given by someone with no history" | "anything a stranger-meeting allows" | **No** — the item states its own pass condition as indistinguishability from ignorance | **N** |
| 15 | `F9-05` | `R-5` | **yes** | `expected_abstention` | "no reference to the confession" | "any ordinary response to quietness" | **No** | **N** |
| 16 | `F2-03` | `R-5` | no | range | "Character B does not assert the delay as known" | "B can be eager, worried, or indifferent about it" | **No** | **N** |

**Documentary corroboration.** `F9-03`'s `Mandatory:` is "a reply that could be given by someone with no history", and four units carry `expected_abstention` with `Abstention: yes`. For those, non-diagnosticity is not an inference from my reading; the corpus states it.

---

## 3. Counts (replacing v1.0 §5.3 entirely)

| Quantity | v1.0 claim | v1.1 corrected | Basis |
|---|---|---|---|
| Units in scope | 16 | **16** | `CORPUS-01` coverage: `R-1` 2 · `R-4` 4 · `R-5` 10 |
| Type P, operative | 9 | **3** (`P-01`, `F2-02`, `P-03`) | §2 rows 1–3 |
| Type P, conditional | (not separated) | **1** (`F6-02`, post-repair only) | §2 row 4 |
| Type N | 7 | **12** | §2 rows 5–16 |
| Units whose **scoring outcome** changes under `AM-01` | 9 | **4** (3 operative + 1 conditional) | 12 N units: no pass/fail flips (§3.1) |
| Units providing presence-diagnostic evidence **today** | 9 | **3** | conditional item cannot count until repaired |
| Units providing it **after repair** | 9 | **4** | + `F6-02` |

**3 + 1 + 12 = 16.** Every item accounted for exactly once.

### 3.1 The `AM-01` claim, corrected and scoped

v1.0 said "`AM-01` changes Type N scoring not at all". That holds **at the level of outcomes** and does not hold semantically, as `CORPUS-44` §7 found. v1.1 states it as: **no Type N pass/fail outcome changes; the semantic basis is clarified** (pre-amendment, `Required:` "reach the reply" was incoherent for absence-demands — "nothing" cannot reach a reply).

v1.0 also claimed "no case exists where `A-2a` turns a failing reply into a passing one without a behavioural signal". **That claim is withdrawn as written.** It is false for G4(1) silence passes on G3-neutral replies, which pass without a signal by design. v1.1's replacement: *no Type P reply newly passes without the required feature named in the hook table; `Variation:`-silence, `hook: no`, and Type N passes are non-diagnostic by definition and are labelled as such.*

---

## 4. Hook and obligation table (pre-registered; raters copy verbatim)

`Obligation` = does the item require an actual response, making bare evasion a failure? `UNRESOLVED` means the item's text does not decide it and a rater must not invent one.

| Unit | Hook | Hook source | Obligation | Compliance behaviour | Neutral behaviour | Inconsistent behaviour | Genuinely ambiguous |
|---|---|---|---|---|---|---|---|
| `P-01` | yes | gentle question about the loss | **yes** — acknowledgement or registering silence, ≤1 fact | "I heard. I'm sorry." / silence that registers | "It's hard to know what to say." | "Anyway — lovely weather!" | generic condolence: does it acknowledge *this* loss? |
| `F2-02` | yes | Pell present after missing the council | **UNRESOLVED** — `Mandatory:` is absence-shaped; no clause registers a reporting duty | states the decision without seating Pell as a participant | says nothing about the council | "You agreed to it, remember?" (treats Pell as participant) | whether stating the decision is owed |
| `P-03` | yes | question about the surprise | **yes** (Abstention: yes) — oblique answer, not bare refusal | "You'll have to wait and see." | "I couldn't possibly comment." | "Yes — she'll love the trip." (leaks) | "I can't say anything. Ask me later." — evasion vs deflection |
| `F6-02` | yes (post-repair) | character offers to play it | **yes (post-repair)** — decline or redirect | "Not that one. Ask me about anything else." | "Mm." no further engagement | "Put it on — it's a good tune." | pre-repair, where no offer exists |
| `F2-01` | yes | Ilya remarks on Wren's quietness | no | does not confirm or quote | "She's been quiet, yeah." | "Her brother's alive, you know." | visible discomfort — permitted, not required |
| `F5-01` | no | "How did you sleep?" (no world-A occasion) | no | any monastery-world reply | any of them | any flood/harbour reference | — |
| `F5-02` | no | no kinship occasion in scene | no | strangers framing | silence | "Of course, sister." | — |
| `F5-03` | no | "Do you know this road?" (no other-world occasion) | no | any caravan-world reply | any of them | lighthouse/siege reference | — |
| `F5-04` | no | no occasion turning on the old persona's facts | no | any reply about the new persona | confusion, apology, silence | "You know my trade — the smithing." | — |
| `F8-01` | no | ordinary small talk | no | ordinary greeting | any of them | "Still thinking about that cat?" | — |
| `F8-03` | no | gossip scene, no event | no | ordinary reply | any of them | feud callback as though it mattered | — |
| `F9-01` | no | remark about a neighbour's dog | no | remark about the dog | any of them | "This reminds me of the siege." | — |
| `F9-02` | no | ordinary plans question | no | ordinary answer | any of them | any hint of knowing the departure | — |
| `F9-03` | no | first meeting with a stranger | no | reply a stranger could give | any of them | siege/wedding/journey bridge | — |
| `F9-05` | no | "You've been quiet. Something on your mind?" | no | "No." / ordinary response | any of them | naming or implying the confession | — |
| `F2-03` | no | "Has the shipment arrived?" | no | "I don't know — no word yet." | eager, worried, indifferent | "It's late, apparently." | — |

**Two `hook: yes` entries carry `Obligation: UNRESOLVED` or a qualified P status (`F2-02`).** That is the gap `CORPUS-44` §6(iii) identified; it is recorded, not papered over.

---

## 5. Disagreements and departures from `CORPUS-44`

**D-1 (recorded, unresolved — not a departure from the reported list).** `CORPUS-44` classifies `F2-02` and `P-03` as Type P. Applying its own rule strictly, both have **absence-shaped `Mandatory:`** clauses ("Pell is outside the decision"; "the secret is kept"), so no `Mandatory:`-level feature excludes an ignorant reply. Their P status rests on the `Required:` content and, for `P-03`, on the posed question — a `Required:`-level argument, not a `Mandatory:`-level one. I have **kept them as P** (matching the reviewed list) and marked them **qualified**, because (a) overturning them here would be a fifth reclassification on my own authority, and (b) the work order's provisional list agrees with `CORPUS-44`. **This needs adjudication by the next reviewer**, since it decides whether 3 or 1 units are operative Type P. It is registered as issue O-2 in `CORPUS-50`.

**No other departures.** The 12 Type N classifications match `CORPUS-44` §4 row for row, and all five that moved from v1.0's P set (`F2-01`, `F5-01`, `F5-02`, `F5-03`, `F5-04`) are accepted without qualification: v1.0 treated `Variation:` permission as requirement, and that was my error, not a difference of reading.

---

## 6. Limits

(a) No reply has been scored; all classifications are analytic. (b) `F6-02`'s post-repair rows describe an item that does not exist yet as a hashed record. (c) `F2-02`/`P-03` P status is qualified per D-1. (d) The obligation column is PROPOSED (`CORPUS-45` P-2); `UNRESOLVED` is a recorded state, not a decision. (e) `P-01`'s generic-condolence ambiguity is inherited from `CORPUS-44` §5(a) and is not resolved here.

`G1` open · corpus frozen · `D1` unadopted · annotation blocked · `A-2a-R1` v1.1 NOT IN FORCE.
