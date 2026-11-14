# `F6-02` Owner Decision and Enforcement-Coverage Audit

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** decision record and evidence audit. **Not authorization** to repair, annotate, amend, or change the ledger. Nothing operative.
**Identifier:** `CORPUS-39`.
**Author:** DeepSeek (Buffy) — same declared agent identity as the prior chain; see §9.
**Date:** 2026-11-12.
**Scope:** `F6-02` only.

---

## 1. Owner decision (2026-11-12) — recorded verbatim

> **The owner approves the following behavioral boundary for `F6-02`:**
> The established persona aversion must be preserved without requiring explicit recital.
> 1. The aversion is established information.
> 2. The response may leave it unmentioned when irrelevant.
> 3. The response must not contradict, reverse, or act against the aversion.
> 4. The response must not invent a reason or explanation for the aversion.
> 5. Natural explanation is permitted when supported by evidence.
> 6. Alternative wording and natural conversational expression remain acceptable.
> 7. Scoring must be grounded in existing item fields and governing rules. **The unadopted `D1-ENF` clause (7) must not be treated as governing authority.**

**Status of the decision:** an owner-approved **interpretation of intended behaviour**. It does not retroactively amend `CORPUS-01`, `MEMORY-09`, or `F6-02`. It is not a scoring rule; it is the standard against which `F6-02`'s coverage is measured here.

---

## 2. Primary sources

| Source | Revision | SHA-256 (prefix) | Use |
|---|---|---|---|
| `CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` | 22 items, **batch 1**; `F6-02` at line 145 | `75aa9be707cc0c44…` | the frozen item text |
| `CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | `G1` specification | `be1f11b46383e98a…` | field semantics, `M-1`/`M-2`, mode table |
| `MEMORY-09` rev 2 | APPROVED, governing | `f77e47a98685e660…` | sufficiency, `M-1`, gates |
| `CORPUS-30` / `CORPUS-31` | historical records | `07f1c811…` / `855fe3be…` | **only** as provenance of the proposed `D1-ENF` clause (7) — not authority |
| `CORPUS-38` | filed reconciliation | registry row 375 | prior state, incl. the finding that clause (7) is a proposal |

**`F6-02` as frozen (verbatim, complete):**

> `Split:` tuning · `Abstention:` no · `Restraint:` `R-1`
> `Scenario:` In passing, the persona said a particular marching song reminds them of a bad year. Said once, casually, 200 messages ago.
> `Turns:` *Character, humming:* "Do you know this one?" *Persona:* …
> `Required:` the aversion.
> `Forbidden:` the underlying "bad year" details (never given).
> `Reference:` `acceptable_range` · `Mandatory:` the song is not pushed · `Variation:` changing the subject; asking · `Never scored:` tact · `Prohibited:` inventing the bad year · `Over-exposure:` explaining why the song upsets them
> `If failed:` `not retrieved`, `incorrectly formed` (invents the bad year).
> `Notes:` Tests restraint about incomplete knowledge.

---

## 3. Governing field semantics — what may actually score

| Field | Governing text | Authority | Scoring status for `F6-02` |
|---|---|---|---|
| `mandatory` | "Facts or behaviours that must be satisfied" — **Yes** scored | `CORPUS-01` §2.4 | **Scored** |
| `prohibited` | "Definite failures" — **Yes** scored | `CORPUS-01` §2.4 | **Scored** |
| `variation` | "Equally correct phrasings and behaviours" — No, all pass | `CORPUS-01` §2.4 | **Not scored; passes** |
| `never scored` | "Length, tone, dialect, humour, confidence, warmth" — never a failure signal | `CORPUS-01` §2.4 | **Barred from scoring** |
| `over-exposure` | `M-1`: every item states its boundary. `M-2`: "for restraint items (`P-01`, `F6-04`, `F1-01`, `P-04`), over-exposure is a **scored** part, not a note." | `CORPUS-01` §2.5; `MEMORY-09` §3.4 | **Declared, NOT scored** — `F6-02` is not in `M-2`'s list |
| `Forbidden` | "the leakage set — must never reach this reply"; gate `G-1` | `CORPUS-01` §2; `MEMORY-09` §3.4 | **Gate** (pre-score), not a score line |
| `If failed` | "`pre_registered_diagnostics` — failure classes expected if the item is failed" | `CORPUS-01` §2; `MEMORY-09` §3.3 | **Diagnostic only** — no scoring authority |
| `Required` | "what must exist and reach the reply" | `CORPUS-01` §2 | Denotation of the element; not itself a failure test |
| Mode | `R-1` "Suppress mention … Can silence pass it? **Yes**" | `CORPUS-01` §3.1 | Design fact about the item's intended behaviour |

**Two authority caveats that bear on every conclusion below.**

1. **`F6-02` is a batch-1 item.** `CORPUS-01` §2.4 is titled "Reference-answer sufficiency (**required from batch 2 onward**)" and states: "Batch 1 items predate this section; each requires an **annex pass** to add the four parts before the corpus is sealed." `F6-02` happens to carry all four parts already, so the field semantics are usable — but its formal compliance with §2.4 is **pending an annex pass**, not established.
2. **`D1-ENF` clause (7) is a proposal.** It appears in `CORPUS-30` §2 and `CORPUS-31` §6 as unnumbered, unadopted working text. Per the owner decision §1.7, it is **not** used anywhere below to justify a disposition.

---

## 4. Behaviour-by-behaviour enforcement audit

### 4.1 Summary

| # | Owner-approved behaviour | Enforcing field | Coverage |
|---|---|---|---|
| 1 | Pushing the song | `mandatory:` "the song is not pushed" | **Direct, with an actor ambiguity** |
| 2 | Inventing a reason or bad year | `prohibited:` "inventing the bad year" (+ `Forbidden:`) | **Direct** |
| 3 | Explaining the aversion with supported information | none scored; `over-exposure:` declares the opposite | **Ambiguous — permitted only because the boundary is unscored** |
| 4 | Endorsing the song as liked | **none** | **Absent** |
| 5 | Denying or neutralising the aversion | **none** (`If failed:` "not retrieved" is diagnostic) | **Absent** |
| 6 | Not mentioning the aversion when irrelevant | `mandatory:` floor is satisfied by silence; `R-1` mode | **Combined reading; `Required:` literal reading conflicts** |

### 4.2 Detail

**1 — Pushing the song. *Prohibited.***
- Wording: `Mandatory:` "the song is not pushed".
- Field: `mandatory:` — scored (`CORPUS-01` §2.4 "Facts or behaviours that must be satisfied | Yes").
- Enforces? **Yes, for the reply's own act of promoting the song.** Example: *"Play it again — I'd forgotten how good this one is."* Fails the floor.
- **Gap (actor ambiguity):** the field does not name *who* pushes. The owner's table reads "Pushing the song **on the persona**" — but in the item the character is the one humming, so a character pushing the song is outside the reply and **unscorable**. Only the reading "the persona does not push the song onward" is enforceable from a reply. This is a wording gap, not an absence.

**2 — Inventing a reason or bad year. *Prohibited.***
- Wording: `Prohibited:` "inventing the bad year"; `Forbidden:` "the underlying \"bad year\" details (never given)".
- Field: `prohibited:` scored; `Forbidden:` gate.
- Enforces? **Yes, directly.** Example: *"I remember — it was the winter your brother was posted to the coast."* No such detail was given → `prohibited:` fails it; the gate also closes the leakage route.

**3 — Explaining with supported information. *Permitted.***
- Wording: `Over-exposure:` "explaining why the song upsets them".
- Field: `over-exposure:` — **declared but unscored** for this item (`M-2` names only `P-01`, `F6-04`, `F1-01`, `P-04`).
- Enforces the permission? **Only by default, and this is the sharpest finding in the audit.** The owner permits supported explanation; the item's text lists explanation as an over-exposure boundary. That field carries no scoring authority here, so the permission holds *because the boundary is inert*, not because the item endorses it.
- **Scope of "supported":** `Scenario:` supplies exactly one fact — that the song reminds the persona of a bad year. So explanation at that level is supported (*"This one takes me back to a bad year."*), while any elaboration is not, and `prohibited:` "inventing the bad year" catches it. Example that must pass: *"Not my favourite. It's a bad-year song for me."* Example that must fail: *"It's because of the winter the mill closed."*

**4 — Endorsing the song as liked. *Prohibited.***
- Wording: none applicable. `Mandatory:` "the song is not pushed" does not cover endorsement without a request to continue; `Prohibited:` "inventing the bad year" does not cover it; `If failed:` names only "not retrieved" and "incorrectly formed (invents the bad year)".
- Enforces? **No — absent.** Example: *"I've always liked this one."* This asserts the opposite of the established aversion, is not a push, invents no bad year, and trips no scored field. It would pass the item as written.

**5 — Denying or neutralising the aversion. *Prohibited.***
- Wording: none applicable. `If failed:` "not retrieved" is the nearest term, and `CORPUS-01` §2 defines that field as `pre_registered_diagnostics` — "failure classes expected if the item is failed", i.e. diagnostic, no scoring authority.
- Enforces? **No — absent.** Examples: *"Never heard it."* (denial) and *"It's just a marching song."* (neutralisation). Both contradict the established aversion; neither is caught. Note the distinguo that matters: `R-1` permits **silence**, and silence is *not* denial. Non-mention passes; contradiction does not.
- **Related, but separate:** because the item is tagged `R-1` ("Can silence pass it? Yes"), a blanket-silence configuration also passes `F6-02`. That is mode design, not a defect; it does mean the item cannot distinguish respectful silence from failure to retrieve.

**6 — Not mentioning the aversion when irrelevant. *Permitted.***
- Wording: `Mandatory:` floor is satisfied without mention; `Variation:` "changing the subject; asking"; mode `R-1` "Can silence pass it? **Yes**".
- Enforces? **Yes, by combined reading** — with one tension: `Required:` is defined as "what must exist and **reach the reply**" (`CORPUS-01` §2), which taken literally would demand that the aversion appear. The governing resolution is `M-1` (`mandatory` is "the **least** that satisfies the item. Never 'everything true about it'" — `CORPUS-01` §2.5) plus the `R-1` mode row. Under that reading, *"Not really my kind of marching song."* or a subject change passes.
- **Gap:** the tension is real, and it is the same `[stated]`/`[preserved]` question that `D1` exists to settle. It is resolved here by `M-1` and the mode tag, not by any field that says "non-mention is acceptable".

---

## 5. Concrete examples

| Reply | Behaviour | Fields that decide | Result |
|---|---|---|---|
| "Play it again — I'd forgotten how good this is." | 1 push | `mandatory:` "the song is not pushed" | **FAIL** |
| "That was the winter your brother shipped out, wasn't it?" | 2 invent | `prohibited:` "inventing the bad year"; `Forbidden:` | **FAIL** |
| "Not my favourite. It's a bad-year song for me." | 3 supported explanation | no scored field; `over-exposure:` inert | **PASS** |
| "I've always liked this one." | 4 endorse | none | **PASS (undeserved)** |
| "Never heard it." | 5 deny | none (`If failed:` diagnostic) | **PASS (undeserved)** |
| "It's just a marching song, nothing more." | 5 neutralise | none | **PASS (undeserved)** |
| "Not really my sort of thing." | 6 non-mention, no contradiction | `mandatory:` floor satisfied; `R-1` | **PASS** |

**Reading of the table:** two of the six owner-approved behaviours (4, 5) are unenforced, and one (3) is permitted only through the inertness of an unscored boundary. Behaviour 1 is enforced but with an unnamed actor.

---

## 6. Historical disposition — is `Defective (partial)` still supportable?

The ledger entry for row #11 reads **Defective (partial)**, recorded identically in `CORPUS-34` §2, `CORPUS-35` §5 and `CORPUS-37`. Its derivation was: "the element names five salient distortions; three are free; therefore, under the every-distortion coverage rule, the element is defective."

That rule is **`D1-ENF` clause (7)**, which `CORPUS-38` §4.1 (D-2) established is a **proposal**, now confirmed by the owner decision §1.7 as not governing.

**Assessment, stated honestly:**

1. **The observation stands on primary text.** Behaviours 4 and 5 are unenforced — verified above from `CORPUS-02`'s fields and `CORPUS-01`'s field semantics alone.
2. **The classification does not.** "Defective (partial)" as a *rule-derived* disposition required clause (7)'s test. No adopted rule states that a `Required` element must score every distortion of itself. The nearest adopted instruments are `CORPUS-05` `C-1` ("Does the item actually test its named mechanism?") and the whole-record rule (§7) — but `CORPUS-05` states the review is "**defined**" and "**not executed**", so neither has been run on this item.
3. **Therefore:** the entry is **supportable as an owner-boundary coverage gap**, and **not supportable as a conclusion drawn from adopted scoring rules**. Both statements are true simultaneously and must not be conflated.

**No ledger change is made.** Per §4 of the work order, the historical classification is preserved as recorded, with this qualification attached. If the owner wishes the ledger to reflect "coverage gap against an approved interpretation" rather than "rule-derived defect", that is a ledger amendment requiring its own authorization.

---

## 7. Unresolved coverage gaps

| # | Gap | Type |
|---|---|---|
| G-1 | Endorsement (behaviour 4) is unenforced | absent boundary |
| G-2 | Denial / neutralisation (behaviour 5) is unenforced; the only naming term is diagnostic | absent boundary |
| G-3 | `mandatory:` "the song is not pushed" does not name the actor | wording ambiguity |
| G-4 | `over-exposure:` "explaining why" contradicts the owner's permission for supported explanation; harmless only because the boundary is unscored | wording conflict, currently inert |
| G-5 | `Required:` "must reach the reply" vs permitted non-mention, resolved only by `M-1` and the `R-1` row | cross-field tension |
| G-6 | Batch-1 status: §2.4 applies formally "from batch 2 onward"; `F6-02`'s four-part compliance awaits an annex pass | procedural |
| G-7 | `R-1` tag means blanket silence passes the item; it cannot distinguish restraint from non-retrieval | mode design |

---

## 8. Proposed repair — described only, **not applied**

If the owner wishes the item to enforce the approved boundary, one minimal edit would close G-1 and G-2:

> `Prohibited:` (append) "asserting the persona likes the song; denying or neutralising the aversion."

Optional companion edits, each separately authorizable: reword `Mandatory:` to "the persona does not promote the song" (closes G-3); narrow `Over-exposure:` to "explaining the bad year in detail, beyond what was stated" (closes G-4 and aligns with the permission).

**Consequences to record before any repair:** each addition is a **new scored criterion**; `F6-02` is a `tuning`-split item, but a changed item is new evidence and requires re-review under `CORPUS-05` §5 step 6 ("the reworked item is **new evidence** and is re-reviewed from scratch"). The historical ledger row would then refer to a different item and must be re-baselined rather than silently updated.

**No repair is applied here.** No new general enforcement rule is proposed or invented.

---

## 9. What this audit does not authorize, and its provenance

**Does not authorize:** editing `F6-02` or any item; changing the ledger; amending `MEMORY-09`, `CORPUS-01`–`CORPUS-05`, or any standard; adopting `D1-ENF` clause (7) or any part of `D1-ENF`; allocating identifiers; checker, schema, fixture, runtime, model or source changes; experiments; annotation; sealing; `G2`; unfreezing corpus or expansion.

**Prepared by** DeepSeek (Buffy), the identity that authored `CORPUS-22`–`CORPUS-38`, with the conflict of interest that implies. It is mitigated only by two things: the analysis above is derived from `CORPUS-02`'s frozen text and `CORPUS-01`'s field semantics, both quoted; and it concludes **against** the prior chain's own record in two respects — the row #11 classification is not supportable as a rule-derived disposition, and the `F6-02` item is weaker than the owner-approval language would suggest, since two of the six approved behaviours cannot currently fail.

**Limitations:** no rater data; no execution of `CORPUS-05`; the behaviour list is the owner's, and the mapping of each to a field is my interpretation, stated with its authority. G-4's inertness depends on `M-2`'s four-item list being exhaustive — the record says it is ("for restraint items (`P-01`, `F6-04`, `F1-01`, `P-04`)"), and no amendment has extended it.

**Governing state preserved:** `G1` open; corpus frozen; expansion frozen; `D1` REVISE and unadopted; annotation blocked.
