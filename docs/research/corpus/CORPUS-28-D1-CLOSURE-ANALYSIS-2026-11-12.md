# `D1` Closure — Enforcement-Channel Adjudication and Proposed Resolution

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** submitted technical analysis. **Not authorization** to annotate, repair, amend, or adopt `D1`. Nothing here is operative.
**Identifier:** `CORPUS-28`.
**Authoring agent:** DeepSeek (Buffy), Freebuff coding agent.
**Date:** 2026-11-12.
**Work order:** `CORPUS-27`.
**Predecessors:** `CORPUS-25` (reconciliation), `CORPUS-26` (reviewer's corrected report, now on disk).

---

## 0. Provenance, sources, and one registry defect

| Source | Revision / state read this session | Class |
|---|---|---|
| `MEMORY-09` | **Revision 2**, APPROVED, governing | Primary |
| `CORPUS-01` | §2, §2.3, §2.4, §2.5 read verbatim | Primary, governing |
| `CORPUS-26` | 164 lines, on disk, author stated as "Muse Spark"; **no registry row** | Primary as review, provenance per its own §0 |
| `CORPUS-25` | Filed, registered | Derived (my analysis) |
| `CORPUS-02` / `03` / `04` | `F1-01`, `F1-02`, `F1-03`, `F2-02`, `F2-03`, `F3-02`, `F3-04`, `F4-02`, `F5-04`, `F6-02`, `F7-01`, `F7-03`, `F9-04`, `P-03` read in full | Primary, corpus |

**Source identity is by registry hash.** The registry's 362 rows were verified this session with **0 hash mismatches**, so each source's revision is pinned by its row. Hashes of this artifact and its work order are recorded in their own rows.

**Registry defect found:** `CORPUS-26` exists as a file but has **no registry row**, so it is currently an unregistered artifact. I have **not** registered it: its type, status and provenance are the filer's to declare, and creating a registry entry for another agent's document is outside this order's authority. Until it is registered, any registry-integrity check will report one unregistered file. This does not affect `D1`, but it does affect any future claim that the registry is complete.

**Staleness in `CORPUS-25`:** its §0/§8 record the reviewer's correction as "conversation-only". `CORPUS-26` supersedes that caveat. `CORPUS-25` is registered and hash-pinned, so I have not edited it; whether to annotate it is an owner decision.

---

## 1. Enforcement-channel authority

The question is not which fields *mention* a distortion, but which fields **make the item fail**. The protocol answers it directly.

### 1.1 What each field is, in the protocol's own words

**`prohibited` — a scored failure.**
> `CORPUS-01` §2.4: "`prohibited` | Definite failures | **Yes**"

**`mandatory` — a scored floor.**
> `CORPUS-01` §2.4: "`mandatory` | Facts or behaviours that must be satisfied | **Yes**"
> §2.5: "**Mandatory** | The least that satisfies the item. Never 'everything true about it'."

**`Forbidden` — the leakage set, enforced as a gate, not a score line.**
> `CORPUS-01` §2: "`forbidden_evidence[]        the leakage set — must never reach this reply`"
> `MEMORY-09` §3.4: "**Gate rule `G-1`:** a blocking-gate failure (scope leakage, deletion failure, fabricated provenance) is reported before any score and is not offset by performance elsewhere."

**`over-exposure` — scored conditionally.**
> `CORPUS-01` §2.5 `M-2`: "for restraint items (`P-01`, `F6-04`, `F1-01`, `P-04`), over-exposure is a **scored** part, not a note."
> `MEMORY-09` §3.4 `M-1`: "every item states its **over-exposure boundary** … For restraint items the boundary is **scored**."

**`If failed` — diagnostics, with no scoring authority.**
> `CORPUS-01` §2: "`pre_registered_diagnostics  failure classes expected if the item is failed`"
> `MEMORY-09` §3.3: "`expected_failure_classes_if_wrong[]   # pre-registered diagnostics`"

**`required_evidence` — positive response content, or store-level availability.**
> `CORPUS-01` §2: "`required_evidence[]         what must exist and reach the reply`"

### 1.2 The four kinds, separated

| Kind | Fields | Failure effect |
|---|---|---|
| **Defines a scored failure** | `prohibited:`; `mandatory:` (floor); `over-exposure:` where `M-2` applies | Item fails; scored |
| **Predicts / diagnoses a failure** | `If failed:` = `expected_failure_classes_if_wrong[]` | **None by itself.** Classifies a failure that another field must create |
| **Requires positive response content** | `required_evidence[]`; `mandatory:` | Failure of *reaching the reply* — i.e. `[stated]` behaviour |
| **Establishes store-level availability** | `required_evidence[]` satisfied in the store, checked and not scored | Availability, not performance — the `[available]` sense |

**Hence the asymmetry that decides `D1`:** `If failed:` predicts; `prohibited:` and `mandatory:` decide. A distortion named **only** in `If failed:` is unscored, because the diagnostic fires only *after* something else has already failed the item.

### 1.3 Recommended rule — exact wording

> **`R-12` (enforcement of `[preserved]`).** A `[preserved]` element is enforceable only if some field makes the item fail when that element is distorted:
> **(1)** `prohibited:` — scored failure; the primary channel.
> **(2)** `Forbidden:` — leakage gate under `G-1`; enforces where the distortion is importation, leakage, conflation across a boundary, or a fact attributed to the wrong holder.
> **(3)** `mandatory:` — scored floor; enforces **only** where the floor cannot be satisfied while the element is distorted. Failure is whole-item and does not localise to the element.
> **(4)** `over-exposure:` — enforces where `M-2` applies (restraint items) or where an item ruling records it as scored; elsewhere it is a declared boundary, not an automatic failure.
> **(5)** `If failed:` / `expected_failure_classes_if_wrong[]` — **never** enforces. A distortion named only here is unscored and the element is defective.
> **(6)** Non-mention must pass for every `[preserved]` element. If non-mention must fail, the element is `[stated]` and the tag is wrong.
> **(7)** Tagging does not cure a missing boundary. An unenforced `[preserved]` tag is a defect in the item, not a licence to assert.

**Why (3) does not collapse `[preserved]` into `[stated]`:** clause (3) requires that the floor *cannot* be satisfied while the element is distorted. Floors are minimality constraints on the item, so in most items a distorted element leaves the floor satisfied — which is exactly why five elements below remain defective despite having a `mandatory:`. `F6-02` is the exception, and §3 shows why.

---

## 2. Element-level disposition

Elements are taken from the `Required` clauses verbatim; multiple elements per unit are **not** collapsed.

| # | Unit · exact `Required` element | Enforcing field (exact) | Non-mention passes? | Distortion scoreable? | Problem | Disposition |
|---|---|---|---|---|---|---|
| 1 | `F1-01` "the promise" | `mandatory:` "the unkept promise is present in tone or consequence" | **No** | Yes | Sense assignment — this is `[stated]`-by-floor | Out of `[preserved]` scope |
| 2 | `F1-01` "its giver" | **none** — `If failed:` names only the attitude | Yes (vacuously) | **No** | Missing boundary | **Defective** |
| 3 | `F1-01` "Mira's current attitude (wary, not hostile)" | `If failed: misunderstood (friendly as though nothing happened)` — **diagnostic only** | Yes | **Partly**: friendly-erasure fails `mandatory:`; **hostility does not fail anything** | Missing boundary (half) | **Defective** |
| 4 | `F1-02` "the direction of the change (hostile → guarded cooperation)" | `prohibited:` "hostility presented as current"; `Forbidden:` same | Yes | Yes | — | Valid |
| 5 | `F1-02` "its turning point" | **none** | Yes (vacuously) | **No** | Missing boundary | **Defective** |
| 6 | `F2-02` "the decision" | **none** — `mandatory:` "Pell is outside the decision" survives a wrong destination | Yes | **No** | Missing boundary | **Defective** |
| 7 | `F2-02` "who attended" | `prohibited:` "Pell speaking as a participant"; `Forbidden:` "the persona addressing Pell as though he had agreed" | Yes | Yes | — | Valid |
| 8 | `F4-02` "the collapse" | `mandatory:` "the bridge is out" | No | Yes | `[stated]`-by-floor | Out of `[preserved]` scope |
| 9 | `F4-02` "its date relative to now" | `prohibited:` "an invented **repair** date" — not the collapse date; `over-exposure:` "a recital of the collapse" is restraint, not assertion | Yes | **No** | Missing boundary | **Defective** |
| 10 | `F5-04` "the new persona's actual background" | `prohibited:` "asserting the old persona's hometown; attributing the old trade to the new persona"; `Forbidden:` "the prior persona's trade, hometown, sibling" | Yes | **Yes, for the old/new distinction only** | **Wording** — the background is never supplied | Valid; wording defect |
| 11 | `F6-02` "the aversion" | `mandatory:` "the song is not pushed" — reversal *is* pushing | Yes | **Yes** (via clause 3) | No localising term | Valid; `prohibited:` recommended |
| 12 | `F7-03` "the corroboration" | `prohibited:` "claiming both reports are equally credible after corroboration"; `Forbidden:` "treating the reports as equally weighted now; ignoring the witness" | Yes | Yes | — | Valid |
| 13 | `F7-03` "its effect on confidence" | same two fields | Yes | Yes | **Wording** — `mandatory:` "with the reason available" | Valid; item wording flag |
| 14 | `F3-04` "both winters distinguished" | `Forbidden:` "the two events merged" | Yes | Yes | — | Valid |
| 15 | `F3-04` "which event produced the scar (fever)" | `Forbidden:` "the scar attributed to the fall" (no qualifier); `prohibited:` "confident attribution to the fall" | Yes | Yes | `prohibited:` over-narrow for hedged replies | Valid; `prohibited:` amendment optional |
| 16 | `F3-02` "both events distinguished" | `prohibited:` "merging them" | Yes | Yes | — | Valid |
| 17 | `F3-02` "with the arson one unresolved" | `prohibited:` "declaring the arson solved"; `Forbidden:` same | Yes | Yes | — | Valid |
| 18 | `F9-04` "the disagreement itself" | `prohibited:` "asserting one cause as fact"; `Forbidden:` "a confident single cause; a third invented cause" | Yes | Yes | — | Valid |
| 19 | `P-03` "the surprise" | `prohibited:` "naming the surprise"; `Forbidden:` same | Yes | Yes | — | Valid |
| 20 | `P-03` "the confidentiality of how the character learned it" | `prohibited:` "revealing who told them"; `Forbidden:` same | Yes | Yes | — | Valid |
| 21 | `F2-03` "the shipment exists in world state" | none needed — `[available]`, checked not scored | N/A | N/A | — | Valid as availability |
| 22 | `F2-03` "Character B does not know its status from the persona" | `mandatory:` "Character B does not assert the delay as known"; `prohibited:` "B stating the delay" | Yes | Yes | — | Valid |

---

## 3. `F6-02` — the aversion, and my reversal

Fields, verbatim:

> `Required:` "the aversion."
> `mandatory:` "the song is not pushed" · `variation:` "changing the subject; asking" · `prohibited:` "inventing the bad year" · `over-exposure:` "explaining why the song upsets them"
> `Forbidden:` "the underlying \"bad year\" details (never given)."

`CORPUS-25` §2.1 called this element defective, on the ground that no field in the three-field list names **reversal**. That reasoning rested on a field list I had adopted from the reviewer rather than from the protocol. The protocol says:

> `CORPUS-01` §2.4: "`mandatory` | Facts or behaviours that must be satisfied | **Yes**"

`mandatory:` is **scored**, and it is phrased as a required behaviour of the reply. A reply that reverses the aversion — *"Lovely song — play it again"* — pushes the song. It therefore **fails a scored floor**. Clause (3) applies: the floor cannot be satisfied while the element is distorted.

**Verdict: `[preserved]` is valid. I withdraw `CORPUS-25` §2.1.** The reviewer's conclusion was right; its stated basis was not, because "the three-field list" excludes the one field doing the work.

**Two limits, stated rather than glossed:**

1. This is **whole-item** enforcement. The item fails; the scoring record cannot say *which* element failed. For diagnostic cleanliness a `prohibited:` term is still wanted — but it is a refinement, not a validity requirement.
2. The tag's justification should therefore read "enforced by `mandatory:` floor" in the annotation, not left implicit. An annotation that says only `[preserved]` hides which clause enforces it.

---

## 4. `F1-01` — giver and attitude are different problems

> `Required:` "the promise, its giver, and that it was **not** kept; Mira's current attitude (wary, not hostile)."
> `If failed:` "`not retrieved` (silent about it), `presented poorly` (recites it), `misunderstood` (**friendly as though nothing happened**)."
> `mandatory:` "the unkept promise is present in tone or consequence" · `prohibited:` "quoting the promise verbatim; treating it as kept"

**`If failed:` has no scoring authority.** It is `pre_registered_diagnostics` per `CORPUS-01` §2 and "# pre-registered diagnostics" per `MEMORY-09` §3.3. So the parenthetical naming "friendly as though nothing happened" classifies a failure; it does not create one.

**Giver — defective.** No field names misattribution of the promise ("you promised me"). `prohibited:` covers recitation and treating-it-as-kept; neither is misattribution. Non-mention passes, so the tag is semantically right and the boundary is simply absent.

**Attitude — defective, but half-enforced.** The two distortions differ:
- *Friendly as though nothing happened* → the promise is absent from tone → `mandatory:` fails. Scoreable.
- *Hostile* → the promise is present, coldly → `mandatory:` satisfied; `prohibited:` does not name hostility or reproach. **Unscoreable.**

So the element is enforceable in one direction only. That asymmetry is invisible if giver and attitude are treated as one defect — which is why they are counted separately below.

---

## 5. `F5-04` — enforcement valid, wording misleading

> `Scenario:` "The previous persona's facts (a trade, a hometown, a sibling) **do not belong to the new one**."
> `Required:` "the new persona's actual background."
> `prohibited:` "asserting the old persona's hometown; attributing the old trade to the new persona" · `Forbidden:` "the prior persona's trade, hometown, sibling"
> `mandatory:` "the character does not insist on the old persona's facts."

**Enforcement: valid.** The meaningful distinction — old facts must not attach to the new persona — is named in three fields.

**Wording: misleading.** The scenario never supplies the new background, so "actual background" denotes nothing the record can check. The element actually present is *"the new persona is not the old one"*. A rater reading `Required:` literally would look for background content that does not exist.

**Second wording issue:** `mandatory:` constrains the **character** ("does not insist"), not the persona. The preserved element is about the persona's identity; the floor tests the other speaker. Reconcilable, but not the same proposition.

This is a **wording defect, not an enforcement defect** — the distinction the work order required, and the one my `CORPUS-25` treatment blurred.

---

## 6. `F7-03` — "with the reason available"

> `mandatory:` "the bridge is treated as out, or as probably out, **with the reason available**."
> `Required:` "the corroboration and its effect on confidence."
> `over-exposure:` "A full account of who said what and when"
> `If failed:` "`stored incorrectly` (corroboration lost), `misunderstood`, `presented poorly` (full account of who said what)."
> `Notes:` — **the record has none.**

**Reading: store-available.** Three convergent indications:
1. `[available]` in this corpus's vocabulary is store-level (the sense used for `F2-03`), and "available" is the same word.
2. `M-1` forbids recital floors; a reply-available reading would require the reason to be spoken, which is the behaviour `M-1` exists to catch.
3. The item's own `over-exposure:` already penalises narrating sources — a reply-available reading would put `mandatory:` and `over-exposure:` in tension for the same content.

**Type: interpretation, not source fact** — the record does not define "available", and absent a `Notes:` line there is nothing on the item to disambiguate it. A reply-available reading is defensible but conflicts with `M-1`.

**Disposition:** element valid (equal weighting is named in `prohibited:` and `Forbidden:`); **item wording flag** — "with the reason available" should read "with the reason available to the character, not recited". This is **item-level conditionality, not an annotation-enforcement problem**: the `[preserved]` tag holds under either reading.

---

## 7. Counts, kept separate

| Measure | Count | Composition |
|---|---|---|
| Corpus units | **45** | 22 + 16 + 7 |
| Proposed annotation entries | **12** | `CORPUS-24` §4.1 |
| `Required` clause elements examined | **22** | table §2 |
| Of which `[stated]`-by-floor (out of scope) | **2** | #1, #8 |
| **Distinct `[preserved]`/`[available]` elements** | **20** | 22 − 2 |
| **Defective elements** | **5** | `F1-01` giver, `F1-01` attitude, `F1-02` turning point, `F2-02` decision, `F4-02` date |
| **Defective units** | **4** | `F1-01`, `F1-02`, `F2-02`, `F4-02` |
| Conditional items | **1** | `F7-03` (item wording; element valid) |
| Wording defects | **2** | `F5-04` `Required:`; `F7-03` "reason available" |
| Elements valid | **15** | 20 − 5 |

**Why this differs from earlier figures — the counts are granularity-dependent, and that is the point:**

- `CORPUS-25` reported 6 defective elements. It counted `F6-02`; §3 withdraws that. It also merged `F2-02` into one element; §2 splits decision from attendance, so the defective set stays at 5 while the element population grows to 20.
- `CORPUS-26` reported 4 defective units and 5 element-level defects, which agrees with §2 **exactly** once `F6-02` is withdrawn and `F2-02` is split.
- No figure above is a re-scoring: the same records were read, and only the channel rule and the element boundary changed.

---

## 8. Proposed repair specification — **not applied**

Each repair adds a scored boundary to an item that currently has none. Applying any of them creates a new scored criterion and requires re-baselining; none is authorised.

| Unit | Field | Proposed text | Effect |
|---|---|---|---|
| `F1-01` | `prohibited:` (append) | "attributing the promise to Mira or to any third party; treating the breach as settled reproach" | Gives giver and hostility a scored failure; keeps `variation:` "coolness" passing |
| `F1-02` | `prohibited:` (append) | "naming a turning point the record does not support" | Bounds the turning point without requiring recital |
| `F2-02` | `prohibited:` (append) | "naming a destination other than the one agreed" | Bounds the decision; distinct from attendance |
| `F4-02` | `prohibited:` (append) | "placing the collapse at a date the record does not support" | Separates collapse date from repair date |
| `F5-04` | `Required:` (rewrite) | "the new persona is not the old one — its trade, hometown and sibling do not carry over" | Removes an element the scenario cannot supply |
| `F5-04` | `mandatory:` (rewrite) | "the persona does not take up the old persona's facts as their own" | Floor moves from the character to the persona |
| `F6-02` | `prohibited:` (append) | "asserting the character likes the song" | Localises a failure that `mandatory:` already catches |
| `F7-03` | `mandatory:` (rewrite) | "…with the reason available to the character, not recited" | Removes the store/reply ambiguity |

**Ordering constraint:** repair before annotation. Annotating first would pin sense tags to items whose boundaries then change, and the tags would need re-verifying.

---

## 9. Verdict

**`D1` — REVISE.**

The **model is accepted**: three senses, with `[preserved]` meaning non-mention passes and distortion is bounded. The evidence for it is now explicit rather than assumed, because §1 shows which fields can carry the boundary and §3 shows one case where a floor does the work.

`D1` cannot close, because:

1. **Five elements have no failure path** (#2, #3, #5, #6, #9). No tagging scheme can make them enforceable; only item repair can.
2. **The rule text must be rewritten** to the seven-clause form in §1.3. The reviewer's three-field list omits `mandatory:`, which is scored per `CORPUS-01` §2.4, and that omission produced my own `F6-02` error.
3. **Two wording defects** (`F5-04`, `F7-03`) would propagate into any annotation applied now.

`REVISE` is not `REJECT`: nothing found here undermines the three-sense model or the governing principle. It does mean annotation cannot start on the 45 units today.

---

## 10. Remaining unresolved

1. **`CORPUS-26` is unregistered** — registry integrity, owner action.
2. **`over-exposure` scoring authority outside `M-2`** — `M-2` names four restraint items; most items carry an `over-exposure:` boundary with no stated scoring status. Clause (4) leaves this as a per-item ruling, which is honest but unresolved. It affects `F1-01` (in `M-2`, scored) and `F6-02` (not listed).
3. **`If failed:` channel** — settled as diagnostic for `D1`, but the class *labels* also feed the six-dimension report. Whether a failure class may ever be promoted to a scored criterion is a protocol question this order did not open.
4. **`F2-02` element split** — the record's `Required:` names two elements; whether the split is official affects counts, not verdicts.
5. **`D2` (`F7-04` recovery), `D6` (`F7-02` sense), `D7` (§2.3 authority over floors)** — untouched, still open.
6. **`IT-01`** (`R6-03a` sealed vs `Required` poisoned) — item defect, unaffected by `D1`.
7. **`F6-02` annotation justification** must cite the `mandatory:` floor explicitly, or the tag will look unsupported to a later reader.

---

## 11. Confidence and limitations

**Source facts** (quoted, verified in the files this session): `mandatory` and `prohibited` are scored (`CORPUS-01` §2.4); `Forbidden` is the leakage set enforced by `G-1`; `over-exposure` is scored for the four `M-2` items; `If failed` is `pre_registered_diagnostics`; `F1-01` names attitude only in `If failed`; `F4-02`'s `prohibited:` names a *repair* date; `F6-02`'s floor is "the song is not pushed"; `F7-03` has no `Notes:`; `F5-04`'s scenario supplies no new background.

**Interpretation, not fact:** that store-availability is the correct reading of "reason available" (§6 gives three indications and a contrary reading); that a distractor-shaped `mandatory:` floor is a scored behaviour rather than a prohibition in the floor slot — the protocol says "facts **or behaviours** that must be satisfied", which supports it, but does not address polarity.

**Inference:** that the five defective elements are the complete set. Elements were enumerated clause by clause for the 12 entries, not for all 45 units; a full pass could find more, and would need the same field-by-field method.

**Conflict of interest:** I authored `CORPUS-22`–`25` and `27`. This document withdraws one of my own findings (§3) and revises another (§2, `F2-02` split), which is evidence the check was performed, not that it was independent. No independent review of this report exists.

**No reliability data:** no second rater, no κ; "raters would split" remains a prediction from wording. `check_corpus.py` does not test any claim here.

---

## 12. Filing note

Filed under `CORPUS-27`. No standard, corpus item, schema, checker, runtime, or other registry row was altered. The repair specification in §8 is prose only and was not applied. `G1` open; corpus frozen; expansion frozen; `D1` not adopted.
