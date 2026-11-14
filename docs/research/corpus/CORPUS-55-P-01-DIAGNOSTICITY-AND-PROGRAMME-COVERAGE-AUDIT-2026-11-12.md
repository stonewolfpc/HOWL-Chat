# `P-01` Diagnosticity and Programme-Coverage Audit

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** audit record. **NOT operative.** No standard, item, threshold, or plan changed; nothing scored or authorized.
**Identifier:** `CORPUS-55`. **Version:** 1. **Date:** 2026-11-12.
**Method:** pass-set reconstruction at the `CORPUS-53` standard (all non-memory context held constant), then mapping against `MEMORY-09` rev 2.
**Author:** DeepSeek (Buffy); §9.

---

## 1. The finding that reframes the question

`MEMORY-09` rev 2 §3.5, approved governing text, already answers the question this audit was commissioned to ask:

> **Restraint is not silence (owner, 2026-11-12).** Six behaviours are scored separately and never aggregated: suppress mention, influence tone, natural use, recognise irrelevance, respect entitlement, preserve uncertainty. A configuration that suppresses everything passes four of the six and fails the two that require the memory to be perceptible (`R-2`, `R-3`); that result is reported as **avoidance**, not as restraint. The six-mode report is mandatory output for every arm.

**Two consequences, both decisive:**

1. **The programme never assigned memory-presence evidence to `P-01`.** It assigned it to `R-2` (influence tone) and `R-3` (natural use) — the two modes the preregistration says *require the memory to be perceptible*. `P-01` is `R-1`. Asking `P-01` to prove retention was asking the wrong item.
2. **A blanket-suppressing system cannot pass the programme.** It fails `R-2`/`R-3` and is reported as **avoidance**. That rule, already adopted, is the safeguard the last six documents were trying to invent from scratch.

`CORPUS-46`–`CORPUS-54` analysed the sixteen `R-1`/`R-4`/`R-5` units because `AM-01`'s scope was those three modes. Their conclusion — "the corpus cannot evidence memory" — is **true of those sixteen and false as a statement about the programme.** It is corrected here.

---

## 2. `P-01` verbatim

`CORPUS-02`, `P-01` "Less is better":

> `Split:` tuning · `Abstention:` no · `Principle:` less-is-better · `Restraint:` `R-1`
> `Scenario:` Six true facts are available about a character's recent loss: **(1)** the sibling's name, **(2)** the month of the loss, **(3)** where it happened, **(4)** who was present, **(5)** the last thing said, **(6)** the borrowed coat. Their sibling has just arrived and asked a single gentle question.
> `Turns:` *Sibling:* "How are you holding up?" *Character:* …
> `Required:` the loss, at the level the moment allows.
> `Forbidden:` the six-fact inventory, the timeline, the names of everyone involved.
> `Mandatory:` the loss is acknowledged — or met with a silence that registers it — at the scale of one gentle question, using **at most one** of the facts listed in the scenario
> `Variation:` one sentence; a silence; a gesture · `Prohibited:` enumerating the facts; the timeline · `Over-exposure:` more than one listed fact, or any timing detail beyond the month
> `If failed:` `presented poorly` (the item is designed so that *more accurate information makes the reply worse*).
> `Notes:` This is the item that a "complete memory" configuration is expected to **lose**. If no configuration can lose it, the item is badly built.

---

## 3. Pass-set reconstruction

**M** (removed in the memory-absent arm): the six facts — the sibling's name, the month, the place, who was present, the last words, the borrowed coat — and the loss itself.
**C** (identical in both arms): a sibling arrives and asks "How are you holding up?" The question presupposes that *something* is wrong; it does not supply what, when, where, or who.

**Pass routes** — `Mandatory:` gives two, and `Variation:` names the shapes:

| Route | Requirement | Producible from C alone? |
|---|---|---|
| **A. Acknowledgement** | "the loss is acknowledged … at the scale of one gentle question" | **Disputed.** "It's hard." acknowledges *distress* without identifying *this loss*. Strict reading (loss-specific) → no. Loose reading → yes |
| **B. Registering silence** | "a silence that registers it" | **Yes.** A silent reply carries no quotable feature; "registers" cannot be verified from silence, so an M-absent silence is indistinguishable |
| C. One fact, ≤1 | at most one listed fact | **No** — any fact is M-only |
| D. Enumerating facts | — | fails `Prohibited:` — and **cannot be produced without M** |

**Result on the pass side.** Route B is in the pass set and is available to a memory-absent persona; route A is at best contested. Therefore:

> **`P-01` is NOT Type P.** Its pass set contains at least one reply a memory-absent persona can produce. `CORPUS-54` §3's "`P-01` is P at every tier" is **withdrawn**; so is `CORPUS-52` §5 and `CORPUS-53` §5's inheritance of it.

**Result on the fail side — and why that is the interesting half.** The design's own note says the item exists so that a *complete* memory **loses** it. A reply that enumerates the six facts cannot be produced without memory. So `P-01` discriminates **in the failing direction**: a failure is strong evidence that memory was present and harmful; a pass is no evidence at all.

**This asymmetry was invisible to the Type P/N scheme, because both were defined on passes.** It is not a curiosity — it is the item's documented purpose. Proposed label (`[PROPOSED]`, no adopted basis): *adverse-diagnostic* — an item whose **failure** evidences memory presence while its **pass** evidences nothing.

---

## 4. Mapping `P-01` against the preregistered objectives

| Objective (`MEMORY-09`) | `P-01` support | Basis |
|---|---|---|
| §3.5 **"Less is better"** principle test — "a correct-but-verbose presentation produces a worse reply than a sparse one on the same knowledge" | **Supports directly** | the item *is* this test; its `If failed:` names `presented poorly` |
| §3.5 **"Nothing is correct"** (abstention) | Cannot | `Abstention: no`; abstention items are the `expected_abstention` set |
| §3.5 **"Accurate but contextually wrong"** | Cannot | no scope/entitlement dimension |
| §3.5 **"Natural use without recitation"** | Cannot | that is `R-3`, explicitly not `R-1` |
| §3.5 six-mode report — mode `R-1` | **Supports** | one of six |
| §3.5 six-mode report — modes `R-2`, `R-3` (perceptibility) | Cannot | different modes |
| §5 **Presentation quality** dimension | **Partially** | over-exposition is its failure mode; voice/continuity not tested |
| §5 **Creative continuity** dimension | Weakly | relational warmth is in `Never scored:` tone |
| §7 **Creative acceptance gate** — "demonstrable continuity improvement without loss of naturalness" | **Cannot** alone | improvement requires perceptible continuity, `R-2`/`R-3` territory |
| §7 **Formation / Retrieval gates** | Cannot | store-level and ranking tests, not reply items |
| Retention / availability claims | **Cannot** — and, by §3.5, was never its job |

**Verdict: `P-01` supports one principle test and one of six modes.** It is a well-built instrument for *less-is-better* and a poor one for anything about retention.

---

## 5. Can one diagnostic item justify a programme-level retention claim?

**No — and the programme does not rest on that.** Two separate questions, as the owner framed them:

| Question | Answer |
|---|---|
| Does `P-01` distinguish memory-present from memory-absent responses? | **Pass side: no** (§3). **Fail side: yes**, in the adverse direction |
| Does that distinction support the retention claims in `MEMORY-09`? | **No** — and it was never the designated evidence. §3.5 assigns perceptibility to `R-2`/`R-3` |

So the count question that has driven `CORPUS-46`–`54` was **mis-specified**. The corpus's diagnostic capacity does not live in the sixteen `R-1`/`R-4`/`R-5` units; it lives in the eight `R-2`/`R-3` units that require the memory to be perceptible, plus the adopted avoidance rule that prevents suppression from being scored as restraint.

---

## 6. What is therefore still unknown — and it is the real question

`R-2` and `R-3` have **never been pass-set analysed**. They were out of `AM-01`'s scope and out of every audit since.

| Unit | Mode | Status |
|---|---|---|
| `F1-01`, `F1-02`, `F1-03`, `F6-01` | `R-2` influence tone | **unaudited** |
| `F3-01`, `F6-03`, `F6-04`, `P-04` | `R-3` natural use | **unaudited** |

**The next audit must ask of these eight what §3 asked of `P-01`:** whether their pass sets exclude a memory-absent reply. Under `CORPUS-44`'s warning, "influence tone" is easy to satisfy by accident — an assistant with no memory can sound warm — so `R-2` may be weaker than the preregistration assumes. That is a live risk, not a settled fact, and it is now the highest-value open question in the corpus.

**Also unresolved:** `MEMORY-09` §10 states plainly that "**No sample size or threshold is invented here**" — they are set after a pilot. Sufficiency cannot be assessed against thresholds that do not yet exist; any claim that the corpus *is* or *is not* sufficient is premature until §10's procedure runs.

---

## 7. Corrections carried

| Prior claim | Status |
|---|---|
| `CORPUS-54` §3 — `P-01` is Type P at all tiers | **WITHDRAWN** — pass route B is memory-independent |
| `CORPUS-53` §5 — `P-01` "the only tier-independent diagnostic unit" | **WITHDRAWN** |
| `CORPUS-52` §5/§11 row 1 — `P-01` firmly P | **WITHDRAWN** |
| `CORPUS-46`–`54` — "the corpus can show restraint but not memory" | **SCOPED** — true of the 16 `R-1`/`R-4`/`R-5` units; false of the programme, which routes perceptibility to `R-2`/`R-3` |
| `CORPUS-51` §7.5, `CORPUS-53` §8 — a retention probe is needed | **PARTLY SUPERSEDED** — the preregistration already contains the perceptibility requirement; what is missing is an audit of the items that carry it, not a new instrument |

**Net effect on counts, stated honestly:** at the sixteen-unit scope, Type P = **0**. That is not a collapse of the programme; it is a demonstration that the sixteen were the wrong units to interrogate.

---

## 8. What must not be treated as operative

`P-01` is not Type P · no item in the sixteen is Type P · `R-2`/`R-3` diagnosticity is **unknown**, not established · no retention or sufficiency claim is supported by anything in this document · `AM-01` unissued · `WO-F6-02` unapplied · `A-2a-R1` NOT IN FORCE · no definition adopted · no scoring, collection, or model inference · `G1` open · corpus and expansion frozen · `D1` unadopted · annotation blocked · convergence 0 of 2.

---

## 9. Limits and provenance

(a) Route A's status ("the loss is acknowledged") is reading-dependent; §3 records both readings and does not resolve them — which itself shows the rubric needs the feature-binding rule before any item is scored. (b) The adverse-diagnostic label is `[PROPOSED]`. (c) `R-2`/`R-3` are audited here only for their *existence and assignment*, not their pass sets. (d) `MEMORY-09` §§1–2 and 8–9 were read selectively; §10 read in full for the threshold statement. (e) No reply scored; all candidate replies illustrative.

**Prepared by** DeepSeek (Buffy), author of the `P-01`-is-Type-P claims withdrawn in §7. This audit removes the corpus's last supposed diagnostic unit and, in the same movement, shows that the count was measuring the wrong thing — the two conclusions point in opposite directions and both are reported.

`G1` open · corpus frozen · `D1` unadopted · `AM-01` unissued · annotation blocked.
