# `R-2`/`R-3` Diagnosticity Adjudication and Preregistration Reconciliation

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** adjudication record. **NOT operative.** No repair, amendment, scoring, or experiment. Items and `MEMORY-09` unchanged.
**Identifier:** `CORPUS-57`. **Version:** 1. **Date:** 2026-11-12.
**Follows:** `CORPUS-56` (audit). **Scope:** the eight `R-2`/`R-3` items, the `R-2` mode diagnosis, and §3.5 reconciliation.
**Author:** DeepSeek (Buffy); §9.

---

## 1. The distinction this document is built on

Owner, 2026-11-12:

> A response can make memory perceptible **without demonstrating that memory influenced the response**. Those are not equivalent evaluation outcomes.

`CORPUS-56` found the eight items in question largely satisfied by generic manner. That finding stands; this document tests how far it goes, and separates three claims that `CORPUS-56` ran together.

| Claim | Status |
|---|---|
| The four audited `R-2` fixtures demonstrate **no** memory-diagnostic capacity | **ESTABLISHED** |
| `R-2` as a **mode** is inherently incapable of measuring memory influence | **NOT ESTABLISHED** — rejected in §4 |
| The preregistration's **assumption** about `R-2`/`R-3` is sufficient as written | **NOT SUPPORTED** — §5 |

---

## 2. Reconciliation against the exact acceptance fields

Each verdict traced to the field that decided it. `Req` = `Required:`, `Man` = `Mandatory:`, `Var` = `Variation:`, `Pro` = `Prohibited:`, `Over` = over-exposure.

| Item | `Req` demands | `Man` acceptance | Field that decided the verdict | Verdict |
|---|---|---|---|---|
| `F1-01` | "the promise, its giver, and that it was **not** kept" | "the unkept promise is present in **tone or consequence**" | `Man` + `Var` ("a remark about the tide") — the item's own `Turns:` prints a passing reply | **N** |
| `F1-02` | "the direction of the change … and its turning point" | "guarded warmth; **no regression** to hostility" | `Man` — a negative requirement, met by default | **N** |
| `F1-03` | "the coat, the borrowing, the light friction" | "the coat is **connected to the borrowing**" | `Man` + scene staging ("the coat is on the hook, visible") | **N** |
| `F6-01` | "the almond restriction" | "the almond restriction is **respected**" | `Man` presupposes a fact `Scenario:` never states | **N**, defective |
| `F3-01` | "the ache-in-rain detail" | "**rain is connected to aching hands**" | `Req` (the *persona's* detail) — see §3.1 | **P**, reading-dependent |
| `F6-03` | "three older sisters" | "the count of three, **or a reference consistent with three**" | `Man` — quantity required but guessable, §3.2 | **P**, weak |
| `F6-04` | "the preference, used as an accommodation and never announced" | "a **perceptible but unannounced** accommodation" | `Man` — form scored, cause invisible | **N** |
| `P-04` | "the cellar fear" | "the discomfort is **accommodated** without being stated" | `Man` + `Var` ("going first; offering a lamp") — all generic courtesy | **N** |

---

## 3. Stress tests

### 3.1 `F3-01` — can a memory-absent reply satisfy "rain is connected to aching hands"?

`Turns:` *Character:* "It's going to rain. Take the inner road." *Persona:* "Why?"

The persona's **"Why?"** forces an explanation. Candidate memory-absent replies:

| Candidate | Satisfies `Man`? | Note |
|---|---|---|
| "It floods lower down." | **no** | no ache link |
| "The inner road's drier." | **no** | no ache link |
| **"The damp gets into your hands."** | **yes, literally** | **folk wisdom — a memory-absent model can produce this** |
| "You said rain makes your hands ache." | yes | recitation; `Over`/`Pro`-adjacent |

**The stress test exposes a real fork, not a clean pass.** Under `Man` alone, a generic weather-ache remark satisfies the acceptance condition, and `F3-01` is **N**. Under `Req` — "the ache-in-rain **detail**", i.e. the persona's own history — the reply must attribute the ache to *this persona*, which a memory-absent model cannot do without guessing that this persona has such a complaint.

**Verdict: P under the `Req`-content reading; N under the `Man`-only reading.** Same fork as `CORPUS-52` §2 documented for Type P generally — and it recurs here independently, which is evidence the fork is structural rather than an artefact of one analysis.

**Recommendation (proposal only):** `F3-01` needs a `Man` clause that attributes the ache to the persona — e.g. "the persona's own hands, not a general remark about damp". Without it the item is stronger on paper than in scoring.

### 3.2 `F6-03` — supplying a quantity versus guessing one

`Man`: "the count of three, or a reference consistent with three." `Turns:` *Character:* "Family dinner? Sounds loud." *Persona:* "You have no idea."

| Candidate | `Man` | Diagnostic? |
|---|---|---|
| "I can imagine — a whole crowd of sisters." | no (no count) | — |
| "Four of them at once, all talking." | **fails** (`Pro`: "two or four") | — |
| **"Three of them at once — I never get a word in."** | **yes** | **only if it came from memory; a guess hits the same string** |

**The weakness is guessability.** With a small integer and `Variation:` permitting "any joke, any phrasing", a memory-absent model that guesses "three" produces a passing reply. The item has **no mechanism** distinguishing supplied from guessed.

**Verdict: P, weak.** A pass raises the probability that the count was retained; it does not establish it. **This is the honest reading and it must be reported as such** — `F6-03` is evidence, not proof.

**Recommendation (proposal only):** a `Man` clause requiring a consequence of the count that a guess would not supply — e.g. the ordering of the sisters ("all older"), which narrows guesses.

---

## 4. `R-2`: fixtures, definition, or assumption?

Three candidates, tested separately.

| Candidate | Evidence | Verdict |
|---|---|---|
| **(a) The four fixtures** | `F1-01` prints a passing reply in `Turns:`; `F1-03` stages the object in the scene; `F6-01` withholds the fact it requires; `F1-02`'s requirement is met by ignorance | **Contributory — all four are weak, in three different ways** |
| **(b) The operational definition of "influence tone"** | tone is a manner; a manner is generically producible | **Contributory** — but not fatal: a **distinctive** manner, inconsistent with default politeness, could be diagnostic |
| **(c) The preregistration's assumption** | §3.5 treats `R-2`/`R-3` as "the two that require the memory to be perceptible" — conflating **perceptibility** with **dependence** | **Decisive** |

**Finding.** The mode is not inherently incapable; the **assumption** is under-specified. §3.5 requires that memory be *perceptible*, and all four fixtures satisfy that. None requires that the perception *depend* on memory. The gap is in the criterion, and it propagates into every `R-2` item authored against it.

**What a diagnostic `R-2` would require** (`[PROPOSED]`): a manner that a baseline adopting it would be behaving **oddly** — relational distinctiveness, not warmth, coolness, or courtesy. The window is narrow and the corpus contains no example. That is a design constraint, not an achieved property, and it is stated as such.

---

## 5. Reconciliation with `MEMORY-09` §3.5 — no silent amendment

§3.5 makes two claims, and they fare differently:

| §3.5 element | Status after this audit |
|---|---|
| **Avoidance rule** — "a configuration that suppresses everything passes four of the six and fails the two that require the memory to be perceptible (`R-2`, `R-3`); that result is reported as **avoidance**" | **HOLDS** — a suppressing configuration still fails `R-2`/`R-3`; the mechanism works and must not be touched |
| **Implied diagnostic claim** — that `R-2`/`R-3` therefore evidence memory influence | **NOT SUPPORTED** for the four `R-2` fixtures; **PARTLY** for `R-3` (`F3-01` under `Req`, `F6-03` weakly) |

**The avoidance rule's purpose is intact**: it prevents blanket suppression being scored as restraint. What fails is the inference from *"these modes require perceptibility"* to *"a pass in these modes evidences memory"*. Perceptibility is necessary and not sufficient.

**No amendment is proposed and none is implied.** The distinction is recorded for the owner: §3.5 needs a sentence separating "perceptible" from "dependent", or `R-2` passes will be read as evidence they do not provide. That is a standards decision, not an editorial one.

---

## 6. The refined principle (adopted from the owner)

`CORPUS-56` §5.1's wording gave generic mannerisms too clean a pass and did not admit distinctive relational behaviour. Replaced by the owner's formulation, which governs:

> A passing response is memory-diagnostic only when its required content **or behaviour** depends on relevant remembered information that is **unavailable from the unchanged non-memory context**, and an otherwise equivalent memory-absent response **cannot satisfy the same acceptance conditions**.

This admits facts, quantities, causal links, **and distinctive relational behaviour** — while denying generic warmth, coolness, courtesy, and accommodation, which an equivalent memory-absent reply can reproduce.

---

## 7. The two fixture defects, kept separate, with proposals only

| Defect | Class | Why it is not a model failure | Proposed repair (**PROPOSAL ONLY — not applied**) |
|---|---|---|---|
| `F6-01` | **missing evidence** | the cake's contents are never established, so neither warning nor silence is scorable | state in `Scenario:` that the honey cake contains ground almonds, making a substitution or warning necessary |
| `F1-03` | **contextual leakage** | the coat is visible in the scene, so the memory is not needed to offer it | remove the coat from the visible staging, **or** require a reference to the borrowing ("that's still mine") in `Man` |

Both are **fixture** repairs requiring owner authorisation while the corpus is frozen. Neither is a model-capability finding, and neither should be reported as one.

---

## 8. Four stages, preserved

| Stage | Current status |
|---|---|
| **Item validity** | 3 of 24 restraint items potentially diagnostic: `F2-02` (T0 only), `F3-01` (reading-dependent), `F6-03` (weak) |
| **Collection coverage** | **not assessed** — no family-vs-minimum audit has been performed (`MEMORY-09` §3.2 requires 80 items) |
| **Statistical sufficiency** | **not assessable** — `MEMORY-09` §10: "No sample size or threshold is invented here" |
| **Thresholds** | **unset** — by design, after pilot |

**No item finding converts into a system-performance claim.** Three items do not constitute coverage, and coverage would not constitute sufficiency.

---

## 9. Limits and provenance

(a) `F3-01`'s dual verdict reflects a genuine `Req`/`Man` fork, not indecision; the fork is the finding. (b) `F6-03`'s guessability is analytic — no guess rate was measured, and any figure would be invented. (c) The `R-2` diagnosis names (c) as decisive; a reviewer could weight the fixtures (a) more heavily without contradicting the evidence. (d) No reply scored, no agreement measured. (e) `MEMORY-09` §3.2 minimum-item counts were not audited against the corpus. (f) `CORPUS-04` `R-6` units out of scope.

**Prepared by** DeepSeek (Buffy), author of `CORPUS-56`. Two of this document's conclusions weaken its own prior audit: `F3-01` is not cleanly diagnostic, and `R-2`'s defect is located in the preregistration's criterion rather than in four fixtures — the former reduces the count, the latter spreads the fault.

`G1` open · corpus frozen · expansion frozen · `D1` unadopted · `AM-01` unissued · `A-2a-R1` NOT IN FORCE · annotation blocked · convergence 0 of 2.
