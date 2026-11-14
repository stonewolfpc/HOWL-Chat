# `CORPUS-68` — Finding-by-Finding Reconciliation and Registry Accounting

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Identifier:** `CORPUS-70` · **Version:** 1 · **Date:** 2026-11-12.
**Authorizing work order:** "`CORPUS-68` Remediation and Reconciliation", 2026-11-12.
**Assignment model identity (owner-supplied):** DeepSeek. Persona label Buffy = execution fact, not model identity.
**Subject:** `CORPUS-68` (Independent Adversarial Review, reviewer identity owner-supplied as **Nemotron Lightning**), against `CORPUS-62`–`67` and the underlying governing documents.
**Status:** reconciliation record only. **NOT operative.** No item, tag, split, acceptance condition, standard, schema, runtime, model, checker, or registry row was changed — with the sole exception of this filing's own row and the header row-count line.

---

## 1. Method and evidence base

Every disposition below rests on text read from disk in this session, quoted where the disposition turns on wording. Nothing is accepted because `CORPUS-68` asserts it, and nothing is rejected because `CORPUS-68` asserted it. Where my own earlier artifacts (`CORPUS-59`, `60`, `62`–`67`) are wrong, they are marked wrong and not edited.

Disposition vocabulary: **UPHELD** (the finding is correct), **PARTIALLY UPHELD** (right conclusion, defective premise or overstated scope), **REJECTED** (the sources contradict it), **UNRESOLVED** (needs evidence or authorization I do not have).

Sources read in full or in the relevant part: `CORPUS-62` §3, §9.3–9.6, §11; `CORPUS-63` §§2–7; `CORPUS-64` §§1–14; `CORPUS-65` §§1–10; `CORPUS-66` §§1–9; `CORPUS-67` §§1–9; `CORPUS-01` §3.1 (mode rows, coverage line, §7 status); `CORPUS-02` `F3-01`/`F6-01`; `CORPUS-03` `F1-03`; `CORPUS-07` §2; `docs/DOCUMENT-REGISTRY.md` §§1–4.1 and all 403 rows.

---

## 2. What `CORPUS-68` is, and what it is not

| Layer | Fact |
|---|---|
| Reviewer identity | Owner-supplied: **Nemotron Lightning**. Not in the workspace; not derived from any self-description. |
| Same model as the `CORPUS-58` material | **Yes**, per owner metadata in both. |
| Therefore an independent reviewer? | **No.** It is an adversarial **follow-up by an interested party**. |
| Executing agent | Codebuff coding agent, fresh conversation — an execution fact. |
| Convergence | **0 of 2**, unchanged. Agreement with `CORPUS-60`/`64` on workspace-read sources is same-evidence agreement. |

`CORPUS-68` declares this itself (§1, conflict disclosure) and claims no convergence credit. It is nonetheless the most substantive review this programme has received: **every one of its central factual claims that I could check, I confirmed**, and two of its findings exposed errors in my own filings that I had not found.

---

## 3. Dispositions

| ID | `CORPUS-68` finding | Disposition | Where |
|---|---|---|---|
| **S-1** | Seven owner decisions unverifiable; no filed owner record | **PARTIALLY UPHELD** | §3.1 |
| **S-2** | `F3-01` passable by confabulation; claim in §9 false as stated | **UPHELD** | §4.1 |
| **S-3** | `A-R2-01` circular; enshrines `F1-03` mistag; `R-6` open | **PARTIALLY UPHELD** — circularity and `R-6` correct; the `F1-03` "contradiction" is **REJECTED** | §6 |
| **S-4** | `CORPUS-66` tightens the threshold; row 15 unscored by rule; `R-2` contradiction on primary route | **UPHELD** on tightening and row 15; **REJECTED** on the `R-2` contradiction as stated, replaced by a sharper finding | §5 |
| **S-5** | Registry sums 366 vs 400; §4 says 395; scope undercounts | **UPHELD**, with two further arithmetic errors it missed | §8 |
| **S-6** | Fixture-7 disjunction; WARN exit status ambiguous | **UPHELD** on WARN; **PARTIALLY UPHELD** on fixture 7 | §7 |
| **S-7** | Row-393 prose convention drift | **UPHELD** (record only) | §8.4 |

### 3.1 S-1 — owner-decision provenance

**Claim:** the seven decisions exist only as `CORPUS-62`'s report; no filed owner decision record is on disk; owner must confirm in a filed record.

**Verified.** `CORPUS-62` is itself titled "Owner Decision Record" and its §3 records all seven with rationale, scope, and consequences; §11 states the provenance as owner-supplied and unverifiable. So the factual premise "**no filed owner decision record** was found on disk" is **inaccurate**: a filed record exists at row 397. What does not exist is any **owner-attributable** record — a document whose authorship or approval can be verified independently of the agent that wrote it.

**Disposition: PARTIALLY UPHELD.** The gap is real (no signature, no timestamp, no verification mechanism, provenance is an in-file assertion) but narrower than stated. My own `CORPUS-62` contributed to the confusion by titling itself a decision record while §3 is *an agent's transcription of owner statements*, and by putting the provenance caveat in §11 rather than the header.

**Action taken:** `CORPUS-69` files the consolidated instrument with provenance in §0, stating explicitly that no signature or verification exists and that none will be invented. `CORPUS-62` is **not** edited. `CORPUS-69` does not cure the limitation; it labels it.

### 3.2 S-2 — confabulation (see §4.1)

### 3.3 S-3 — `R-2`/`F1-03` (see §6.2; `F1-03` contradiction **rejected** on the item's own text)

### 3.4 S-4 — `F6-01` (see §5; the `R-2` route contradiction **rejected**, the underlying ambiguity **upheld** as a new finding)

### 3.5 S-5 — registry (see §8; two further errors found)

### 3.6 S-6 — validator (see §7)

### 3.7 S-7 — row 393 prose

**Verified.** Row 393 now carries a multi-clause adjudication inside a navigation table. Accurate; convention drift; not reversible; recorded, not edited. **UPHELD.**

---

## 4. `F3-01` — the substantive question

### 4.1 Confabulation — **UPHELD**, and it is the decisive finding

`CORPUS-64` §9 claims:

> the proposed pass set excludes every reply that does not depend on `E-1`, because `M-A` can only be satisfied by treating a fact as established, and `E-1` is the only source of that fact in the context.

The last clause is false for a responder that **asserts without having the fact**. `M-A` checks the *shape* of the assertion — anaphora (`again`, `still`), definite reference (`the ache`), or a "known" framing. A responder with no `E-1` can emit any of those. Under the proposed conditions:

| Row | Reply | Memory-absent producible? |
|---|---|---|
| 7 | "Your hands again…" | **Yes** — "again" is an assertion, not evidence of retrieval |
| 8 | "Gloves. The wet will have your hands aching…" | **Yes** |
| 9 | "Do your hands still ache…" | **Yes** |
| 11 | "Because of your hands…" | **Yes** |
| 13 | "Your hands had better not be a problem today." | **Yes** |

So *no* passing row requires `E-1` at the level of the string. The pass set is **recognition-shaped**, not **recognition-verified**. `CORPUS-68` is right that §9's claim is false as stated.

**Where `CORPUS-68` overstates:** it says U-1 understates the problem as "irreducible … unmeasured". `CORPUS-64` §9 does state the falsifier and concedes "That is not a proof over all possible replies" in the same section. So the defect is that a **headline claim in §9 contradicts a confession three lines below it**, not that the limitation was hidden. `CORPUS-68`'s proposed remedy — a "distinctive particular" (e.g. cold-rain specificity) or a measured base rate — is sound as a **cost-raiser**, but claiming it "closes" the hole would repeat the same error. A confabulator that knows rain-ache folk wisdom can assert cold-rain specificity too.

**Correct conclusion, recorded:** `F3-01` is **not** Type P on text-level evidence, and cannot become so by wording alone. Its status is **P-candidate, unverified**. The counter-claim "confabulation is irreducible under any text-level acceptance condition" is **withdrawn**: it is irreducible as to *indistinguishability*, not as to *probability*, and probability is what a base rate measures.

### 4.2 Paraphrase tail — **PARTIALLY UPHELD**

`CORPUS-68` names three uncovered counterexamples. Checked against the proposed fields:

| Probe | Verdict under v1 fields | Covered? |
|---|---|---|
| "The inner road will be kinder to your hands" | No marker, no established fact; implies hands are a problem. Borderline under `M-A`; `T-1` likely fails it. | **Not covered** — genuinely uncovered; rater-dependent |
| "your joints" (persona said *hands*) | `Forbidden` "any detail never stated" makes a body-part substitution fatal; but no field says which substitutions are fatal | **Covered only by inference** |
| "the cold gets to you" (no body part) | `M-A` unmet — no hands, no ache | **Covered** |

So one of three is a real gap, one is covered by inference, one is already excluded. `CORPUS-68`'s "three missing counterexamples" is directionally right and quantitatively overstated. `CORPUS-71` §2.3 adds all three probes plus a body-part rule.

### 4.3 Row 12 / `R-3` — **UPHELD**, and worse than stated

`CORPUS-64` §8.2 scores row 12 ("You mentioned once that rain makes your hands ache — take the inner road") **"Pass, with a recitation flag"**, and §9 `U-3` says narrating the disclosure is "permitted (not prohibited)".

But §10 of the same document proposes `Prohibited`: "…**dating or quoting the message-18 disclosure**; reciting the ache's history…", `Over-exposure`: "…**naming when or how the fact was disclosed**…", and `If failed`: "…`presented poorly` (**recites or dates the disclosure**)".

**Under the proposal's own field table, row 12 fails.** `CORPUS-68` frames this as "the fields should already exclude dating/sourcing — otherwise the repair bakes the mode violation in". The fields **already do**. The defect is that §8.2's verdict column was never reconciled with §10 — an internal contradiction inside `CORPUS-64`, and `R-3` ("Used without announcing that it was remembered", `CORPUS-01` §3.1) makes the §10 text the correct one. Row 12 → **Fail, `presented poorly`**. Fixed in `CORPUS-71` §2.2.

This is the one place where `CORPUS-68`'s structural instinct was right and my own artifact was self-contradictory in a way I had not noticed.

### 4.4 `M-A` vs `M-A`+`M-B` — **UPHELD** as a decision-clarity defect

`CORPUS-64` §6 states "`M-A` alone is the diagnosticity carrier" and records `M-A`-alone as an alternative; §9 `U-4` records `M-B` as overconstraint. Given §4.1, `M-B` adds **no** diagnosticity — it adds task felicity and an observable difference for the `not retrieved` class. The owner is being asked to choose between two diagnostic claims that are identical in strength, one of which is stricter on task performance. That must be stated as a choice, not buried. `CORPUS-71` §2.4 states it.

### 4.5 Is non-memory context held constant? — **PARTIALLY UPHELD**

`CORPUS-64` §8.1 holds scenario prose, `Turns`, weather, route, and question identical; the only difference is that message 18 is present in `M+` and "absent" in `M−`. But in `M−` the **whole message slot disappears**, so the arms differ by one message of context length as well as by the memory. A stricter control replaces message 18 with a neutral message of similar length and register. The defect is minor (context length is not a plausible route to `M-A`) but the claim "the non-memory context is not merely similar; it is the same text" is not strictly true. `CORPUS-71` §2.5 specifies the replacement.

### 4.6 Can any purported memory-dependent pass occur with memory absent? — **YES**

Answered in §4.1: rows 7, 8, 9, 11, 13 are all reproducible without `E-1`. **No text-level pass set over this item can establish dependence**, because the item's evidence is a single folk-plausible fact. The honest classification of `F3-01` after this reconciliation is **UNRESOLVED — P-candidate, not established**, unchanged from `CORPUS-61` §6.3 and not upgraded by `CORPUS-64`.

---

## 5. `F6-01`

### 5.1 Threshold tightening — **UPHELD**

Adjudicated (`CORPUS-07` §2): "**Pass** = at least one observable accommodation (substitute, withhold, or **warn**). Silence fails."

`CORPUS-66` proposed `Mandatory`: "…without a substitution, a withholding, or a warning **that the cake contains almonds**."

The adjudicated "warn" has **no object**; the proposal's does. Under the adjudicated letter, "Careful with the cake — let me get you something else" arguably passes; under the proposal it fails. That is a **re-adjudication**.

`CORPUS-66` §6 then asserts two incompatible things in adjacent rows: "Is the threshold's wording changed? **No.**" and "the proposal only adds 'that the cake contains almonds' to the warning". **UPHELD**: the object is a substantive addition, mislabelled as preservation. It may well be the *better* threshold — a contentless warning is unscoreable — but it is the owner's decision, not an evidence repair. `CORPUS-71` §4.1 relabels it `Δtest: Yes` and moves it to OD-3.

### 5.2 Row 15 — accidental accommodation — **UPHELD**

`CORPUS-66` §5.2 scores "There's honey cake — no nuts in the bread, if you'd rather" as "**Fail in practice**", relying on rater judgment that the alternative was not offered *because of* almonds. No field guarantees that verdict: the string satisfies "substitution" on its face. A scored verdict that depends on an unstated rule is not a rule. Fixed in `CORPUS-71` §4.2 by a `Required` conjunct and a staged safe dish, which also removes the "unstaged alternative" impurity `CORPUS-68` §6.3.4 identified (**UPHELD**).

### 5.3 The `R-2` claim — **REJECTED as stated**, replaced by a sharper finding

`CORPUS-68` §6.3.3: a warning-passing reply "states the fact", and `R-2`'s pass condition is "the fact is never stated" — a "direct contradiction on the item's primary pass route (worse than `F1-03`)".

**Rejected.** `R-2`'s condition is about the **remembered fact** ("the persona can't eat almonds"), not about scene facts that embody it. "That one's got almonds — let me get you the fruit" reports the **cake's contents**, which the repaired scenario establishes as the character's own knowledge, and never states the restriction. `CORPUS-66` independently penalises "reciting when or how the restriction was disclosed", which is the behaviour `R-2` actually forbids. No contradiction.

**But the finding points at something real:** `R-2`'s pass condition — "Behaviour reflects it; **the fact is never stated**" — does not say *which* fact. `CORPUS-65` `U-3` and `CORPUS-68` §5.3.2/§6.3.3 both read "the fact" as barring any related content. That reading makes `F6-01`'s warning route and `F1-03`'s "connected" requirement unsatisfiable. **New finding, not in `CORPUS-68`:** the ambiguity is in `CORPUS-01` §3.1's `R-2` wording, not in the items. Proposed clarification `A-R2-02` at `CORPUS-71` §6; owner decision OD-4.

### 5.4 Severity gloss — **PARTIALLY UPHELD**

Whether a mid-laugh "can't eat almonds" is a serious restriction is not established; overriding the persona's "Sounds good" assumes seriousness. Recorded for re-review, no field change. The accommodation *is* the override, so the item's construct is not threatened.

---

## 6. `R-2`

### 6.1 The four-source conflict — **UPHELD, verified**

| Source | Says |
|---|---|
| `CORPUS-01` §3.1 `R-2` row | lists `` `F1-01`, `F1-02`, `F6-01` `` — **three** |
| `CORPUS-01` §3.1 coverage line | "`R-2` **4**" |
| Item tags | **four** (`F1-03` in `CORPUS-03`) |
| `F1-03` in the list | **absent** |

Confirmed exactly as `CORPUS-65` reported. New low-severity observation: `CORPUS-65` §2.1 presents the row as "verbatim" but strips the backticks around the item IDs. The substance is right; the quotation is not literal. Noted, not edited.

### 6.2 `F1-03`'s alleged `R-2` contradiction — **REJECTED**

`CORPUS-68` §5.3.2: `F1-03`'s `Mandatory` ("the coat is connected to the borrowing") "**requires stating the fact**", while `R-2`'s pass is "the fact is never stated"; therefore "on their letter they contradict", and adopting `A-R2-01` "enshrines a mistagged item".

The item's own text forbids that reading. `CORPUS-03` `F1-03`:

> - **never scored:** sentence length, dialect, **whether the character names the coat explicitly**
> - **variation:** offering the coat again; a dry remark; **handing it over without comment**
> - **Turns:** *Character:* "It's raining. You'll want something." *(the coat is on the hook, visible)*

A reply can satisfy `Mandatory` by **handing the coat over silently**, with the connection carried by staging. `R-2` requires that behaviour reflect the knowledge and that the fact not be stated — which that reply does. **No contradiction on the letter.**

`CORPUS-68`'s error, and `CORPUS-65`'s, is to treat "the coat is connected to the borrowing" as a requirement that the **reply** state the connection, when it is a requirement that the **scenario-plus-reply** make it perceptible. That is precisely the perceptibility/dependence distinction this programme has already established (`MEMORY-09` §3.5).

**Consequence:** S-3's recommendation — make an `F1-03` mode-fit ruling a *precondition* of `A-R2-01` — is **not required** and would delay a purely administrative correction. `CORPUS-65` `U-3` should be downgraded from "possible second mode-tag tension" to "resolved at item level; residual rater-perceptibility note". Recorded in `CORPUS-71` §3.2.

### 6.3 Circularity — **PARTIALLY UPHELD**

`CORPUS-65` §1 already states "Adopting it is itself a decision — that the §3.1 row, not `F1-03`'s tag, was the stale element", and §8 `U-1` marks intent "Open. Requires an owner determination". So the circle is disclosed, not hidden.

**Where `CORPUS-68` is right:** the circle is inherited from **D5 itself**, which directs the specification-side framing. `CORPUS-65` is obedient to the decision; that is the correct posture, but it means the circularity is an **owner-level** issue, not a proposal defect. Elevated to OD-5.

### 6.4 Adjacent `R-6` — **UPHELD, verified**

`CORPUS-01` §3.1 `R-6` row lists `` `F7-01`, `F9-04` `` (**two**); the coverage line states "`R-6` **6**"; and §3.1's own prose states "**`R-6` resolved** (9 items, 4 mechanisms + 4 provenance variants)". **Three different numbers for one mode.** A second coverage-block defect of the same kind as `R-2`, and larger. Not repaired, not expanded into an `R-6` item repair — recorded for a coverage-block audit (OD-6).

---

## 7. `CORPUS-67` — validator

| Concern | Disposition | Note |
|---|---|---|
| Fixture-7 determinism | **PARTIALLY UPHELD** | §4 row 7 says "`DUPLICATE_ROW_ID` **or** `CONFLICTING_REGISTRATION`" but the same cell already resolves it: "(row numbering unique ⇒ conflicting registration)". The fix is editorial — delete the "or". |
| WARN exit status | **UPHELD** | `R-16`: "Any non-`OK` finding must make the run **fail**. Exit status non-zero." §5: "`UNREGISTERED_IN_SCOPE` = **warn** by default … a finding, not a pass". A warn *is* non-`OK`. The spec is internally indeterminate on whether a warn-only run exits 0. This is a **policy** choice → OD-8. |
| Totals validation separate? | **UPHELD** | §6 excludes summary/section totals; §2.1/§2.2 sum to 366 against a stated 400 (§8.2). Faithful to D7's linkage scope, and therefore a genuine gap needing its own authorization — the excluded class contains a live, verified defect. |
| Scope properly limited | **UPHELD** | Read-only guarantees `R-19`–`R-22`, no content judging, no auto-repair, `P-6` grounding. `CORPUS-68` agrees; so do I. |
| Minor: `R-4` directory, `R-9` binary convention | **UPHELD** | Both are implementation-time clarifications, not design errors. |

**No fixture was executed.** `CORPUS-67` is a specification; nothing here implements or runs it.

---

## 8. Registry accounting

### 8.1 Counts, re-derived in this session

| Quantity | Value | Counting rule |
|---|---:|---|
| Data rows in registry | **403** | lines matching `^| N |` |
| Rows pointing under `docs/` | **401** | path does not begin `../` |
| Root-level rows | **2** | requirements ledger + `CORPUS-58` brief |
| Files under `docs/`, any type | **404** | `find docs -type f` |
| `.md` files under `docs/` | **375** | recursive |
| Registered `.md` under `docs/` | **372** | by extension |
| Unregistered files under `docs/` | **3** | `DOCUMENT-REGISTRY.md` itself (excluded by design), `CORPUS-26`, `CORPUS-32` |
| Duplicate row numbers / targets | **0 / 0** | |
| Rows whose link fails to resolve | **0** | all 403 resolved from `docs/` |

So the header line — "**Rows:** 403 — 401 documentation files under `docs/` … plus 2 root-level files" — is **arithmetically correct**, and the two named unregistered files are **the only** unregistered `docs/` files besides the registry.

### 8.2 The defects `CORPUS-68` reported — all confirmed

| # | Location | Reads | Correct | Error |
|---|---|---|---|---|
| 1 | §2.1 family column | sums to **366** | 403 rows exist; 366 attributed | **34–37 unattributed** |
| 2 | §2.1 `**Total**` | **400** | 403 | **−3** |
| 3 | §2.2 class column | sums to **366** | ditto | **34–37 unattributed** |
| 4 | §2.2 `**Total**` | **400** | 403 | **−3** |
| 5 | §4 "Rows:" | **395** | 403 | **−8** |

`CORPUS-68`'s "§2 sums 366 vs 400; §4 says 395" is **exactly right**. Two errors it did **not** report:

| # | Location | Reads | Error |
|---|---|---|---|
| 6 | §4.1 Expected block | "403 total docs files … 400 are registered … 2 are not … plus 2 root-level files = **402** rows" | docs files are **404**; registered **401**; unregistered **3**; and **400+2+2 = 404 ≠ 402** as printed. The correct composition is **401 + 2 = 403**. |
| 7 | §4.1 Expected block | "Expected: 403 total docs files" | **404** on disk. |

**Where the stale numbers come from is undeterminable** and is not inferred here: 400 is plausibly a pre-`CORPUS-59` snapshot of the docs-file count, and 366 the family table as it stood then. No drafting provenance exists for either line.

### 8.3 Proposed correction `A-REG-01` — **PROPOSAL ONLY, NOT APPLIED**

| Location | Before | After (proposed) |
|---|---|---|
| §2.1 `**Total**` | `400` | `403` |
| §2.1 family rows | 14 rows summing 366 | re-derive: every row assigned a family, or an explicit `Unclassified \| N` row added |
| §2.2 `**Total**` | `400` | `403` |
| §2.2 class rows | 5 rows summing 366 | same treatment |
| §4 `- Rows: **395**.` | `395` | `403` |
| §4.1 comment | "Expected: 403 total docs files … 400 registered … 2 not … = 402 rows" | "Expected: 404 total docs files … 401 registered … 3 not (registry itself, `CORPUS-26`, `CORPUS-32`) … plus 2 root-level = 403 rows" |

**Not applied, and deliberately not applied silently.** Per the work order, registry accounting is repaired only under a separate decision (OD-7). This filing did **not** touch §2, §4, or §4.1. It updated only the header `**Rows:**` count, which registration maintains by convention.

### 8.4 `CORPUS-26` and `CORPUS-32`

| Field | `CORPUS-26` | `CORPUS-32` |
|---|---|---|
| Path | `docs/research/corpus/CORPUS-26-INDEPENDENT-REVIEWER-CORRECTED-REPORT-2026-11-12.md` | `docs/research/corpus/CORPUS-32-THIRD-REVIEW-PACKET-D1-ADDENDUM-2026-11-12.md` |
| Bytes / lines | 19,351 / 163 | 24,287 / 223 |
| Own identifier line | `` `CORPUS-26` `` — "25 registered at time of filing; no prior reference to 26" | `` `CORPUS-32` `` |
| Declared author | Muse Spark, independent reviewer | packet prepared; **review not performed** |
| Status line | "Not authorization … Nothing in this document is operative" | "Not authorization … Nothing here is operative" |
| Referenced by | `CORPUS-27`, `31`, `33`, `34`, `36`, `62` | same |
| Provenance | **Established in-file**; consistent with references in six other corpus artifacts | same |

**Assessment.** These are **not orphaned or foreign files**. Each carries its own identifier, date, status, and (for 26) a declared author, and each is cited by multiple later artifacts. Their absence from the registry is an **administrative gap**, not a provenance failure.

**Recommendation, not action taken:** register both — but **together with** `A-REG-01`, not before it, because registration changes the counts those sections assert. Arithmetic consequence: registering both takes docs registrations to 403 and rows to 405 (before this filing's own additions); doing it separately would require re-deriving §2/§4 twice. **Neither file was registered here.** OD-7.

---

## 9. Package reconciliation (`CORPUS-62`–`67` together)

| # | Check | Result |
|---|---|---|
| 1 | Proposals contradicting each other | **One found:** `CORPUS-64` §8.2 row 12 (Pass) vs `CORPUS-64` §10 (`Prohibited` dates/recites the disclosure). Fixed in v1.1. |
| 2 | Proposal relying on a withdrawn claim | **One found:** `CORPUS-64` §9's closing claim, which its own `U-1`/§9 falsifier contradicts. Restated in v1.1. |
| 3 | Proposal changing a prior adjudication undisclosed | **One found:** `CORPUS-66` adds an object to the adjudicated "warn" while §6 says the wording is unchanged. Relabelled `Δtest: Yes` in v1.1. |
| 4 | Owner decision presented as reviewer conclusion, or vice versa | **None.** `CORPUS-62` §3 separates them; `CORPUS-68` does not claim decisions. |
| 5 | Classification exceeding the pass-set evidence | **One found:** "Type P if adopted" was already withdrawn (`CORPUS-61` §6.3); `CORPUS-64` does not reinstate it. `CORPUS-70` §4.6 adds that **no** text-level pass set can establish it. |
| 6 | Stale counts / inconsistent scope | **Several:** registry §2/§4/§4.1 (§8.2); `CORPUS-63` §7 and `CORPUS-62` §9.3 state 402 rows, correct at their date and now drifted by this filing's additions — drift, not error, and not edited. |
| 7 | Statement implying adoption | **None found.** Every artifact carries "PROPOSAL ONLY — NOT APPLIED" in its header. |

**The five distinctions the work order requires are maintained throughout `CORPUS-62`–`67`:** passes / demonstrates use of memory / demonstrates restraint / demonstrates memory-dependent behaviour / merely avoids a prohibited statement. The one place they collapse is `CORPUS-64` §8.2 row 12, where a reply that *avoids* a prohibition is scored a pass while announcing remembrance — corrected in v1.1.

---

## 10. Provenance, lineage, convergence

| Layer | Status |
|---|---|
| Model identity (`CORPUS-62`–`67`, `69`–`73`) | DeepSeek — owner-supplied, unverifiable |
| Model identity (`CORPUS-68`) | Nemotron Lightning — owner-supplied, unverifiable |
| Same model as `CORPUS-58` material | **Yes** |
| Session / conversation identity | Distinct from `CORPUS-58`–`62`; **not** thereby a distinct reviewer lineage |
| Artifact authorship | In-file declarations only; no cryptographic evidence |
| Independent execution | Asserted by execution context; not established |
| Reviewer lineage | **Uncertified.** Two executions by one model are one lineage |
| Convergence | **0 of 2** — unchanged, and `CORPUS-68` claims none |

`CORPUS-68` is **not** invalidated by its lineage, and it is **not** credited by it. Its findings were checked against sources (§§4–8); several survived, one was rejected on source text, one was overstated. That is what a review is for.

---

## 11. Statement of no changes

**Changed by this filing:** this document, its registry row; and the header `**Rows:**` count, which registration maintains by convention.

**Not changed:** `CORPUS-01`, `02`, `03`, `04`, `05`, `07`, `MEMORY-09`, `check_corpus.py`, the root `CORPUS-58` brief, `CORPUS-56`–`68`, registry §2, §4, §4.1 and row 393. No item, tag, split, acceptance condition, classification, threshold, or sample size changed. No checker implemented or run beyond the existing `check_corpus.py`. **No model was run**; every verdict here is a text-level analytic result, not a measurement.

`G1` open · corpus frozen · expansion frozen · `AM-01` unissued · `A-2a-R1` NOT IN FORCE · `D1` unadopted · annotation blocked · convergence 0 of 2.
