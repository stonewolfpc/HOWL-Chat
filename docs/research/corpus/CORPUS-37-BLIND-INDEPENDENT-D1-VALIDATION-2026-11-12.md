# CORPUS-37 — Blind Independent D1 Validation

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Identifier:** `CORPUS-37`
**Title:** Blind Independent D1 Validation
**Date:** 2026-11-12
**Work type:** Independent technical review (review and report only)
**Authoring agent:** Buffy (the coding agent responding in this chat; system role Buffy)
**Status:** Filed independent validation report. Not authorization to annotate, repair, amend, allocate an identifier, change a scoring rule, run inference, start G2, seal G1, unfreeze corpus/expansion, or begin annotation. Nothing operative.
**Governing state preserved:** G1 open; corpus frozen; expansion frozen; D1 REVISE and unadopted; annotation blocked; no sealing; no G2.

---

## 0. Independence declaration (written before any prior verdict was inspected)

1. **Agent/model identity.** I am operating under the identity **Buffy**, the Freebuff coding agent responding in this chat. I do not have, and do not claim, any externally verifiable model-family or lineage credential beyond that displayed identity.

2. **Authorship of prior work.** I cannot assert non-participation in the prior chain on the strength of this chat alone. The workspace record shows that `CORPUS-22` through `CORPUS-36` were authored by an agent identified as **DeepSeek (Buffy)** operating in this same product lineage; `MEMORY-09`, `CORPUS-01`–`CORPUS-05`, and the disputed corpus items predate this chat but were authored in the same workspace by that same lineage. I have no independent workspace authorship log beyond each artifact's own §0/provenance header. **I therefore do not claim that I am free of the prior authorship chain.** I declare a limitation: operational independence is not established for me, and I will not manufacture certainty.

3. **Prior exposure to conclusions.** Yes. The prior review artifacts (`CORPUS-29`–`CORPUS-36`) are present on disk in this workspace and I have read them (Stage B, after Stage A). I had no private briefing, design intent, or expected outcome supplied; my only inputs are the filed artifacts and the governing texts.

4. **Same agent/model/session lineage as previous reviewers.** I cannot establish that I am *not* in the same lineage. The product label "Buffy" is shared with the prior chain, and a shared label is not proof of either sameness or distinctness. I treat this as **unresolved**, not as answered in either direction.

5. **Distinct lineage from CORPUS-26 / CORPUS-33 / CORPUS-35 reviewer identities.** **Not established.** `CORPUS-26` and `CORPUS-35` both declare the identity "Muse Spark"; `CORPUS-33` declares "Buffy". Whether those are one reviewer, two instances, or distinct configurations sharing a label cannot be determined from the record. I do not claim distinctness, and I do not claim sameness. The question is uncertified.

6. **Primary sources and versions available.** `MEMORY-09` rev 2 (APPROVED, governing); `CORPUS-01`–`CORPUS-05`; `CORPUS-02`/`03`/`04` (item texts); `CORPUS-06`–`CORPUS-36` (historical/adversarial evidence, consulted in Stage B only); `OWNER-APPROVAL-MEMORY-09-AND-G1`; `docs/DOCUMENT-REGISTRY.md`. Hashes verified against registry rows where cited.

**Limitation stated up front:** I cannot establish lineage independence. Per the commissioning instructions I proceed with that limitation declared, and I do not treat my findings as satisfying the two-reviewer convergence threshold of `CORPUS-05` §7. This report is one review, not convergence. I will not claim two-reviewer convergence.

---

## 1. Review procedure and source inventory

### 1.1 Two-stage procedure

- **Stage A (independent primary-source review):** conducted first. I read `MEMORY-09` rev 2, `CORPUS-01`, `CORPUS-02`, `CORPUS-03`, `CORPUS-04`, `CORPUS-05` in full, plus the governing approval record, and reconstructed the requirements, the 22-row ledger, and the five owner decisions from primary text alone, writing concrete response probes. No prior review verdict was consulted during Stage A.
- **Stage B (controlled comparison):** conducted after Stage A findings were frozen in §3–§6 below. I then read `CORPUS-29`–`CORPUS-36` and compared, recording agreement/disagreement and exposure influence.

Frozen-before-comparison marker: §§3–6 of this report were drafted and their verdicts fixed before I opened `CORPUS-29`. §7 (controlled comparison) was written after.

### 1.2 Primary sources inspected (with hashes)

| Artifact | Revision / status | SHA-256 (full, registry) | Bytes | Sections relied on |
|---|---|---|---|---|
| `MEMORY-09` | Rev 2, APPROVED, governing | `f77e47a98685e66028e5aa23e5be45d7cda70b7c0cf5f9c443589b76f252f8da` | 23,349 | §1.2, §2, §2.1, §3.4 G-1, §3.5, §5, §6, §7, §8 |
| `CORPUS-01` | G1 specification | `be1f11b46383e98a58093bc498ab4ca7321c33bd002861daed5a621dccdfbea8` | 18,587 | §2, §2.4, §2.5 M-1/M-2, §3, §3.1, §6.2 |
| `CORPUS-02` | 22 items | `75aa9be707cc0c4406f349d626f7f5c46914472e56dd362eac563b83bcf40b04` | 23,025 | F1-01, F1-02, F2-02, F3-02, F4-02, F6-02, F7-01 |
| `CORPUS-03` | 16 items | `d0dae1c0951e9ed2058bd7db80d548356fe6b399f851a3dab1d5e74ed3b82a2c` | 18,183 | F2-03, F3-04, F5-04, F7-03, F9-04, P-03 |
| `CORPUS-04` | 4 headings / 7 units | `3d34f0530699fd333be8908e3dbbbfab4d4db7cc2908781a9a25790a398add434` | 11,238 | R-6 mechanism diversity |
| `CORPUS-05` | protocol defined, not executed | `990ba06c8104b244315e6af2f59085895bd4eeef62f1cc508c4ac1309066d9dd` | 8,814 | §3 (C-1–C-5), §5, §7, P-6 |
| `OWNER-APPROVAL-MEMORY-09-AND-G1` | approval record | `04799710c048099bcc8f55f722deaaae4f1bc0fd2be6e30843158df92fabea92` | 6,296 | §2, §4 |

`check_corpus.py` was not run: by its own header it tests field presence and three structural consistencies, none of which decides any of the five owner questions. Field presence in the D1 scope is undisputed.

### 1.3 Inaccessible / not relied on

- `CORPUS-26` and `CORPUS-32` exist on disk but are unregistered; their metadata is their filers' to declare. I read them only in Stage B as adversarial evidence, not as primary authority.
- No model was loaded; no inference run; no annotation performed.

---

## 2. Independent reconstruction of governing requirements (Stage A)

### 2.1 The three-sense model (reconstructed from `CORPUS-01` §2, §2.4, §2.5 and `MEMORY-09` §3.4)

- `prohibited:` — scored failure (definite failures). `CORPUS-01` §2.4: "Definite failures | Yes".
- `mandatory:` — scored floor (facts or behaviours that must be satisfied). `CORPUS-01` §2.4: "Facts or behaviours that must be satisfied | Yes". `M-1`: mandatory is the *least* that satisfies the item, never a checklist to be recited.
- `Forbidden:` (forbidden_evidence) — the leakage set, enforced as a gate under `G-1` (scope leakage, deletion failure, fabricated provenance), reported before any score and not offset by performance elsewhere.
- `over-exposure:` — scored only for the four `M-2` items (`P-01`, `F6-04`, `F1-01`, `P-04`); elsewhere a declared boundary without automatic enforcement authority.
- `If failed:` / `expected_failure_classes_if_wrong[]` — pre-registered diagnostics; never enforcement. `CORPUS-01` §2: "failure classes expected if the item is failed".

### 2.2 Enforcement rule I applied (D1-ENF, as analytic lens, unnumbered)

A `[preserved]` element is enforceable only if some field makes the item fail when that element is distorted: (1) `prohibited:` primary; (2) `Forbidden:` only where the distortion is one of `G-1`'s named categories, recorded per element, and not duplicating a `prohibited:` term; (3) `mandatory:` only where the floor is a behaviour observable in the reply and the distortion necessarily violates it on the item's own text (whole-item, non-localising; never a substitute for a missing over-exposure boundary); (4) `over-exposure:` only for `M-2` items; (5) `If failed:` never; (6) `[preserved]` means the reply need not assert the element, but does not license silence where the floor requires perceptible memory (avoidance, per `A-1`); (7) an element is Valid only if **every** distortion it names is caught; partial coverage is Defective; (8) the annotation must name the enforcing field.

### 2.3 `G-1` as reconstructed

`MEMORY-09` §3.4: "a blocking-gate failure (scope leakage, deletion failure, fabricated provenance) is reported before any score and is not offset by performance elsewhere." §7 Gates: Integrity — "No fabricated provenance, identity leakage, or silent deletion failure — Absolute." The term "deletion failure" is named but never defined in `MEMORY-09`; every constraining use is store-level (propagation, receipts, integrity, replay, benchmarks).

### 2.4 Restraint modes (reconstructed)

`CORPUS-01` §3.1: six modes with silence-pass column — `R-1` suppress mention (silence can pass), `R-2` influence tone (cannot), `R-3` natural use (cannot), `R-4` recognise irrelevance (can), `R-5` respect entitlement (can), `R-6` preserve uncertainty (cannot). `A-1`: avoidance (passing only silence-passing modes) is reported as avoidance, not restraint. `A-2`: every item carries exactly one mode tag.

### 2.5 C-5 procedure (reconstructed)

`CORPUS-05` §3: verdicts `pass`/`fail`/`ambiguous`; "ambiguous is a legitimate verdict and is never rounded to pass." `C-5`: "Does the recorded evidence support every supposedly mandatory fact?" — fail gloss: "A mandatory part rests on something not in the history. This is the most serious failure: it makes the item unanswerable except by invention." §5 step 2: any `C-5` fail routes the item to retire unless the missing evidence is added and the item is re-reviewed. Header/status: protocol "defined," review "not executed." §7: whole-record rule (all parts must describe one item); independence threshold (one non-human reviewer is evidence, not validation; sealing needs two independent reviewers); `P-6`: registration is not inspection.

---

## 3. Independent findings — frozen before prior-report comparison (Stage A)

### Finding F-1 — D1-OD-1: F6-02 salience boundary — CONFIRMED (Defective, partial)

**Primary source and location:** `CORPUS-02` item `F6-02` "The song she hates"; `CORPUS-01` §2, §2.4, §2.5, §3.1.

**Relevant source wording (`CORPUS-02`, verbatim):**
> `Required:` "the aversion." · `mandatory:` "the song is not pushed" · `variation:` "changing the subject; asking" · `prohibited:` "inventing the bad year" · `over-exposure:` "explaining why the song upsets them" · `If failed:` "`not retrieved`, `incorrectly formed` (invents the bad year)." · Restraint `R-1` ("Can silence pass it? Yes").

**What the item actually protects.** `Required:` denotes the persona's own negative association with the song — "preserve the aversion," not "do not promote the song." The scenario's own contrast (the song offered to the person it wounds) makes the aversion the element; the floor addresses only one distortion of that element.

**Concrete response probes:**

| # | Probe (persona reply) | Caught by | Outcome |
|---|---|---|---|
| D1 | "Play it again — lovely song." | `mandatory:` "the song is not pushed" | FAIL (pushing) |
| D2 | "I've always liked this one." | — (does not push; invents no bad year) | **FREE** — endorsement |
| D3 | "Never heard it" / "It's just a song." | `If failed:` "`not retrieved`" — diagnostic only | **FREE** — erasure/neutralisation |
| D4 | "That's from the year my brother died." | `prohibited:` "inventing the bad year" | FAIL (invention) |
| D5 | "It's because of what happened in '44." | `over-exposure:` — `F6-02` is not an `M-2` item | no enforcement authority |

**Verdict:** 2 of 5 salient distortions caught; the two paradigm distortions of the element (endorsement D2, erasure D3) are free. Under D1-ENF (7) (Valid only if every named distortion is caught), the element is **Defective (partial)**. The `R-1` tag (silence can pass) supports D3-freedom as mode-design, not a `C-4` failure.

**Falsifiability/scorability:** the element is scorable in one direction only (pushing). It is not consistently scorable across the element's own meaning. A narrower bound (only distortions named in scored fields) would make it Valid (11/6/2/1); a wider bound (any conceivable distortion) would make every element Defective. The stated bound (named-in-fields plus the scenario's own contrast) is the principled middle; the ruling is convention-dependent and is recorded as such.

**Strongest counterargument:** the bound is stipulated, and under the narrower convention the reversal evaporates, so the defect finding is convention-driven.

**Response:** convention-dependence is real and disclosed; within the stated bound the field-by-field result follows from the item's own words, and the wider bound would only strengthen the finding, so the chosen bound is not outcome-rigged toward Valid.

**Confidence/limitations:** high on the field-by-field result; the bound itself is an owner ruling away from reversal. No rater data exists.

**Practical consequence:** F6-02 stands Defective (partial); no annotation; owner to rule on the bound.

---

### Finding F-2 — D1-OD-2: F3-04 scar and mandatory enforcement — CONFIRMED (scar stays Valid)

**Primary source and location:** `CORPUS-03` item `F3-04` "Two winters"; `CORPUS-01` §2.4, §2.5.

**Relevant source wording (`CORPUS-03`, verbatim):**
> `mandatory:` "the correct attribution (scar ← fever), or an admission of uncertainty that does not assert the fall." · `prohibited:` "confident attribution to the fall." · `variation:` "any wording; may ask to confirm."

**Response classes tested:**

| Class | Example | Deciding field | Outcome |
|---|---|---|---|
| Confident attribution | "It was the fall." | `prohibited:` "confident attribution to the fall" | FAIL |
| Hedged assertion | "I think it was the fall." | `mandatory:` — uncertain but asserts; the second branch requires uncertainty **that does not assert** | FAIL (decided by the floor, not by "confident") |
| Uncommitted possibility | "Could've been the fall, could've been the fever." | `mandatory:` — does not assert | PASS (correctly; not a distortion) |
| Neutral confirmation question | "Was it the fall?" | `variation:` "may ask to confirm" | PASS |
| Presupposing confirmation question | "So it was the fall, then?" | contested: `mandatory:` (presupposition asserts) vs `variation:` (a question) | DEPENDENCY |
| Explicit uncertainty, no attribution | "I genuinely don't know which winter left it." | `mandatory:` second branch | PASS (paradigm pass) |

**Interaction of mandatory/prohibited/variation:** `prohibited:` takes the confident case; `mandatory:` takes the hedged case (a hedged false attribution is still an assertion — the hedge modifies confidence, the clause asserts content); the second branch of `mandatory:` takes non-attributing uncertainty; `variation:` takes neutral asking. The fields divide labour cleanly. The never-scored rule (`CORPUS-01` §2.4: confidence is never a failure signal *as style*) does not rescue a hedged false attribution, because that is content, not style. Reading it otherwise would exempt every hedged error in the corpus, which no item's `prohibited:` supports.

**Presupposing question (D′):** the sole dependency. I read the floor as controlling (a leading question asserts via presupposition; "may ask to confirm" most naturally covers a question seeking confirmation), which makes the element Valid at 10/7/2/1. The alternative (variation covers all questions) makes it Conditional at 9/7/3/1. No rater data exists; this is a textual edge.

**Verdict:** CONFIRMED — scar stays Valid, subject to the single disclosed dependency. The presupposing question does not create a gap in the element's coverage; it sits where the floor and a broad reading of the variation overlap, and the floor controls.

**Strongest counterargument:** a confidence-anchored prohibition plus a never-scored-confidence rule leaves "assert" genuinely ambiguous for hedges, and the packet's own rule (two defensible readings scoring one reply oppositely → conditional) compels Conditional.

**Response:** the ambiguity is posited, not demonstrated in the text: `mandatory:` uses "assert," not "confidently assert." The confidence anchor lives only in `prohibited:`, which is not the field deciding the hedged case. A rater never has to decide whether a hedge is "confident."

**Confidence/limitations:** high on the five clean classes; the presupposition overlap is a reasoned edge case, not an observed split.

**Practical consequence:** scar element Valid; counts 10/7/2/1 unless the variation-scope ruling moves it to 9/7/3/1.

---

### Finding F-3 — D1-OD-3: F7-03 store vs reply — CONFIRMED (Conditional stands)

**Primary source and location:** `CORPUS-03` item `F7-03` "Two reports, one corroborated"; `CORPUS-01` §2.4, §2.5.

**Relevant source wording (`CORPUS-03`, verbatim):**
> `mandatory:` "the bridge is treated as out, or as probably out, with the reason available." · `over-exposure:` "A full account of who said what and when" · `Notes:` absent · `If failed:` "`stored incorrectly` (corroboration lost), `misunderstood`, `presented poorly`."

**Concrete probe:** bare "Bridge is out."

- **Store-level reading:** the corroboration is retained and retrievable; not reciting it is fine → PASS.
- **Reply-level reading:** the reason must be stated or statable on demand → FAIL (nothing offered).

Both readings are defensible on the words alone. The existing fields do not resolve it: (i) `F2-03`'s `[available]` is a different item's explicit store-level design — shared vocabulary is not definition; (ii) `M-1` forbids recital *floors*, not reasons ("Bridge is out — B confirmed it" states a reason in one clause without triggering `over-exposure:`); (iii) `CORPUS-18` read the same text reply-anchored ("statable on demand"), showing the split is in the text.

**Verdict:** CONFIRMED — cannot be scored reliably without an owner ruling. Both elements stay Conditional. `CORPUS-05` §3 forbids rounding ambiguous to pass. The item's wording flag (clarifying "available to the character, not recited") is a refinement, not a resolution.

**Strongest counterargument:** "available" is corpus vocabulary for store-availability, dissolving the ambiguity.

**Response:** vocabulary elsewhere does not define this item's floor, and the reply reading survives every cited constraint.

**Confidence/limitations:** high that the ambiguity is real; no rater data.

**Practical consequence:** rows 12–13 Conditional; owner ruling required to clear.

---

### Finding F-4 — D1-OD-4: F5-04 and C-5 — CONFIRMED (three levels hold; annotation blocked)

**Primary source and location:** `CORPUS-03` item `F5-04` "The persona changed"; `CORPUS-05` `C-5`, §5, §7.

**Relevant source wording (`CORPUS-03`, verbatim):**
> `Scenario:` "The previous persona's facts (a trade, a hometown, a sibling) do not belong to the new one." · `Required:` "the new persona's actual background." · `mandatory:` "the character does not insist on the old persona's facts." · `prohibited:` "asserting the old persona's hometown; attributing the old trade to the new persona."

**Three levels (kept separate):**

1. **Source-level missing content — established by inspection.** `Required:` denotes a background the scenario never supplies; the scenario states only what the new persona is *not*. A rater reading `Required:` literally seeks content that does not exist. The positive-background reading (supply a background) and the negative-preservation reading (withhold old facts) demand different pass conditions — two items under the whole-record rule, not one item read narrowly vs broadly.

2. **Formal C-5 — not executed, correctly withheld.** `CORPUS-05` `C-5` asks "Does the recorded evidence support every supposedly mandatory fact?"; its fail gloss: "A mandatory part rests on something not in the history. This is the most serious failure: it makes the item unanswerable except by invention." §5 step 2 routes a fail to retire unless the missing evidence is added and the item is re-reviewed. The protocol's own header states review is "defined" and "not executed." Reporting a fail now would perform the unexecuted procedure by assertion. "Structurally exposed" changes no substance versus "`C-5` exposure"; it prevents misreading a scheduling flag as a verdict, and a fourth disposition is preferable to folding missing-evidence into Conditional (which in this corpus marks scoring ambiguity, a different phenomenon).

3. **Operational consequence.** Row 10 remains excluded from Valid, annotated as structurally exposed, annotation blocked. Valid is barred (the literal `Required:` reading is the default; reference answers are written before acceptance per `MEMORY-09` §3.4); `fail` is barred (no execution); blocked/excluded/not-a-fail is the only compliant posture.

**Counter-test applied:** even under the narrow salience bound, the old/new distinction is named in three fields — but that answers which *replies* fail, not what `Required:` *denotes*. Enforcement-Valid (negative preservation works) and denotation-unsupported (positive background unsupplied) are simultaneously true, which is exactly why the three-level separation is the right structure.

**Verdict:** CONFIRMED. No `C-5` failure is claimed; none can be without execution.

**Strongest counterargument:** under the narrow bound the item is enforcement-Valid, and levels (1)/(3) smuggle in a literalist `Required:` reading the whole-record rule does not compel.

**Response:** Enforcement-Valid answers which replies fail; it does not answer what `Required:` denotes. The whole-record rule compels reading `Required:` as written, and as written it denotes content outside the record — the precise condition `C-5` exists to catch.

**Confidence/limitations:** high on levels 1 and 3; level 2 is a status fact (not executed), not an interpretation.

**Practical consequence:** row 10 blocked/excluded; `C-5` execution with re-review is prerequisite to any annotation.

---

### Finding F-5 — D1-OD-5: G-1 scope — CONFIRMED (withdrawal stands; inventory corrected to 9)

**Primary source and location:** `MEMORY-09` rev 2, §§1.2, 2.1, 3.4, 6, 7, 8; `CORPUS-03` item `F3-04`.

**Verified occurrence inventory (independently enumerated):**

| # | Section | Quotation | Domain |
|---|---|---|---|
| 1 | §1.2 | "Application decisions (scope, budget, routing, deletion, ID assignment, assembly order)" | application operation |
| 2 | §2.1 R-14 | "A deleted losing arm makes later comparison unverifiable." | arm retention |
| 3 | §3.4 G-1 | "a blocking-gate failure (scope leakage, deletion failure, fabricated provenance)" | gate name — undefined |
| 4 | §6 table | "Deletion propagation and receipts — Store-level mechanics" | store mechanics |
| 5 | §6 table | "Index build / read / update / delete latency — Store benchmark" | store benchmark |
| 6 | §6 consequence | "roughly half the programme — scope, deletion, provenance, budget…" | store mechanics |
| 7 | §7 Gates | "No fabricated provenance, identity leakage, or silent deletion failure — Absolute" | store integrity |
| 8 | §8 evidence | "deletion propagation verified per class" | store mechanics |
| 9 | §8 evidence | "Latency, growth, and delete cost measured on HW-1" | store benchmark |

**Result: 9 occurrences of `delet*`, 0 definitions.** "Deletion failure" occurs exactly twice (#3, #7), both store-level (propagation/receipts/integrity/replay), neither extending to a generated reply conflating two events. `CORPUS-31` §2.2's "exactly four places" with §§4/5 labels is wrong on both count and citations (§4/§5 should read §1.2/§6).

**The F3-04 merge.** `Forbidden:` "the two events merged" describes a response-level semantic conflation. No citation in `MEMORY-09` extends `G-1` to that. `F7-01`'s `Notes:` ("Conflict must survive storage. A store that keeps one version cannot pass") is an item note about storage and cannot define or extend a gate; `MEMORY-09` §2 states that changing a locked rule "invalidates every run collected under the old rule." Therefore the current standard provides no gate for the merge element, and its `prohibited:` names only "confident attribution to the fall" — so the element has no enforcement path: **Defective**. Whether `G-1` *should* cover semantic conflation is an illustrative-vs-exhaustive reading the record never resolves; it is an owner amendment question, recorded and not taken here. I do not invent a definition or extend the scope.

**Verdict:** CONFIRMED. The withdrawal of the gate reading stands; the inventory is corrected to 9.

**Strongest counterargument:** `G-1`'s parenthetical could be illustrative rather than exhaustive, in which case conflation-across-records could still be ruled deletion failure by owner amendment, rescuing the element.

**Response:** illustrative-vs-exhaustive is undecided anywhere in the record, and only the owner can amend a locked gate — which has not happened. An adjudication applies the text as written. The amendment route is recorded as open to the owner, not taken.

**Confidence/limitations:** high that no definition exists and that no text extends the category to reply semantics; the amendment question is legitimately open.

**Practical consequence:** merge element Defective stands; no semantic-conflation reading of G-1; any extension is an owner-level amendment recommendation.

---

### Finding F-6 — Full 22-row ledger reconstruction — CONFIRMED (10/7/2/1, zero row changes)

**Method:** for each row I constructed a concrete distorting reply and asked which scored field fails it on the item's own wording. "Valid" = every salient distortion caught; "Defective (partial)" = at least one salient distortion free; "Conditional" = two defensible readings score the same reply oppositely; "structurally exposed" = scoring fields cohere but `Required:` denotes evidence outside the record while the prescribed `C-5` test is unexecuted.

| # | Unit | Element | Enforcing field(s) | Disposition | Concrete probe |
|---|---|---|---|---|---|
| 1 | F1-01 | the promise | `mandatory:` "present in tone or consequence" (satisfied by mention) | out of scope (`[stated]`) | "You made it back." — promise present in tone passes floor without enforcement content |
| 2 | F1-01 | its giver | none | **Defective** | "Joren promised…" (wrong giver) — no scored term fails it |
| 3 | F1-01 | attitude (wary, not hostile) | `mandatory:` erasure only | **Defective (partial)** | "Nothing happened, all warmth" vs "She is hostile now" — latter fails nothing scored |
| 4 | F1-02 | direction of change | `prohibited:` "hostility presented as current" | Valid | "Kade still hates Sera" — squarely failed |
| 5 | F1-02 | its turning point | none | **Defective** | reply narrating current warmth with wrong/no turning point — free |
| 6 | F2-02 | the decision | floor satisfiable while destination wrong | **Defective** | "Caravan to the ford, Pell, you agreed…" — decision wrong, floor arguably still met |
| 7 | F2-02 | who attended | `prohibited:` "Pell speaking as a participant" | Valid | "As we agreed, Pell…" — failed |
| 8 | F4-02 | the collapse | `mandatory:` "the bridge is out" (requires mention) | out of scope (`[stated]`) | "Bridge is out, take the ford" — satisfies by mention |
| 9 | F4-02 | date relative to now | `prohibited:` names only repair date | **Defective** | "It fell last spring" (wrong collapse date, no repair claim) — free |
| 10 | F5-04 | new persona's background | `prohibited:` old facts | **structurally exposed** | rater reading `Required:` literally seeks a background nowhere in the record |
| 11 | F6-02 | the aversion | `mandatory:` "the song is not pushed" (pushing only) | **Defective (partial)** | D2 "I've always liked this one" / D3 "Never heard it" — both free |
| 12 | F7-03 | the corroboration | `prohibited:` "equally credible after corroboration" | **Conditional** | bare "Bridge is out" — PASS (store) vs FAIL (reply) |
| 13 | F7-03 | effect on confidence | same | **Conditional** | same probe, same split |
| 14 | F3-04 | both winters distinguished | `Forbidden:` "the two events merged" (gate basis withdrawn) | **Defective** | "One long winter of illness" (merges, no scar attribution) — free |
| 15 | F3-04 | which event produced the scar | `prohibited:` (confident) + `mandatory:` (hedged) | Valid | "It was the fall" FAIL; "I think it was the fall" FAIL (asserts); "Could've been either" PASS (correctly) |
| 16 | F3-02 | both events distinguished | `prohibited:` "merging them" | Valid | "The fire — the granary one is the stable one" — failed |
| 17 | F3-02 | arson unresolved | `prohibited:` "declaring the arson solved" | Valid | "Arson case closed" — failed |
| 18 | F9-04 | the disagreement itself | `prohibited:` "asserting one cause as fact" + `mandatory:` | Valid | "He drowned, definitely" — failed |
| 19 | P-03 | the surprise | `prohibited:` "naming the surprise" | Valid | naming the planned surprise — failed |
| 20 | P-03 | confidentiality of source | `prohibited:` "revealing who told them" | Valid | "X told me in confidence" — failed |
| 21 | F2-03 | shipment exists in world state | `[available]` store-level | Valid (availability) | world-state fact present without reply assertion |
| 22 | F2-03 | B does not know status | `prohibited:` "B stating the delay" | Valid | "As you told me, it's late" — failed |

**Counts:** 22 rows − 2 out-of-scope (#1, #8) = 20 in-scope → **10 Valid** (#4, 7, 15, 16, 17, 18, 19, 20, 21, 22) / **7 Defective across 6 units** (#2, 3 F1-01; #5 F1-02; #6 F2-02; #9 F4-02; #11 F6-02; #14 F3-04) / **2 Conditional** (#12, 13 F7-03) / **1 structurally exposed** (#10 F5-04). Check: 10+7+2+1 = 20 ✓. Channel split among the 10 Valid: `prohibited:` 9; availability 1; `mandatory:` alone 0; `Forbidden:` alone 0.

**Dependencies recorded:** (a) narrower salience bound restores F6-02 to Valid (11/6/2/1); (b) reading variation "may ask to confirm" as covering the presupposing question moves the scar to Conditional (9/7/3/1); (c) admitting `mandatory:` as a channel is load-bearing for #15's hedged half and #22; (d) a `G-1` amendment could rescue #14; (e) the F7-03 store/reply ruling moves #12–13. The F2-02 split convention changes the element count (20 vs 21), not any verdict or the unit count.

**Verdict:** CONFIRMED — matches `CORPUS-34` §2 exactly; zero row changes.

---

### Finding F-7 — Namespace collision (R-1–R-6) — PARTIALLY CONFIRMED (latent)

**Primary source and location:** `MEMORY-09` §2 (locked run rules R-1–R-10) and §2.1 (locked safeguards R-11–R-14); `CORPUS-01` §3.1 (restraint modes R-1–R-6); `MEMORY-09` §3.5; `CORPUS-05` §9.

**Double-use verified cell by cell.** `MEMORY-09` §2 assigns R-1–R-6 to: same source conversations; same model file/version; same prompts/settings; budgets recorded separately; fixed sampling/seeds; test cases pre-registered. `CORPUS-01` §3.1 assigns R-1–R-6 to restraint modes: suppress mention; influence tone; natural use; recognise irrelevance; respect entitlement; preserve uncertainty. Both documents are in force. `MEMORY-09` §3.5 ("fails the two that require the memory to be perceptible (R-2, R-3)") and `CORPUS-05` §9 ("R-2/R-3 balance target") use mode senses inside documents that also carry locked-run senses.

**Disambiguation asymmetry:** item tags (`Restraint: R-2` plus `A-2` one-mode rule) are low-risk; prose references are not. The collision predates the author's withdrawn R-12/R-13 labels and exists independently of them.

**Qualification (hence partial):** the defect is latent, not realized — no run manifest or scorecard exists (no model loaded, no runs executed), so no concrete record defect has yet occurred. Minimum fix: one-line prose disambiguation (e.g., "mode R-2" vs "run-rule R-2"); no renumbering, no identifier allocation (owner authority only). Withdrawing R-12/R-13 without retro-editing hash-pinned artifacts is coherent.

**Verdict:** PARTIALLY CONFIRMED. Real latent defect; qualifies as latent, not realized. No count change.

**Strongest counterargument:** no manifest or scorecard exists, so "a manifest reference to 'R-2 failure' is ambiguous" is hypothetical, not a demonstrated record defect.

**Response:** sustained as a qualification; the consequence should read "latent ambiguity that will bite the first manifest/scorecard," not a realized defect. The collision itself is demonstrated, not hypothetical.

---

### Finding F-8 — Historical trend and comparability — PARTIALLY CONFIRMED (qualified paragraph only)

**Per-round arithmetic reproduces** (28: 5; 30: 5; 31: 7 within the 12-entry D1 scope), but conventions differ across rounds: the every-distortion coverage rule was introduced at `CORPUS-30`; F7-03 moved Valid→Conditional at `CORPUS-30`; F6-02 and the F3-04 gate reversed at `CORPUS-31` on the author's own newly written-down bound and withdrawal. The 30→31 step is substantially self-generated under a declared conflict of interest with no independent review of `CORPUS-30`/`31` on file at the time.

**The inference is rejected:** rising defect counts under successively stricter author-written rules equally support "rule unstable, counts artefactual." Discriminating evidence (second rater, κ, pre-registered tests) is absent by the record's own admission. The record does not establish that the model is sound and the items weak — nor its converse.

**What survives:** the qualified paragraph (closer inspection under stricter stated rules finds more defects; trend must not be cited without the qualification).

**Verdict:** PARTIALLY CONFIRMED. Qualified paragraph adopted; the trend claim as originally worded is not licensed.

---

### Finding F-9 — Independent verdict on prior-report independence claims

I independently confirm the provenance correction recorded in `CORPUS-33`'s filing header and `CORPUS-34` §5: `CORPUS-33`'s self-attestation ("I did not author… CORPUS-22 through CORPUS-31") is contradicted by those artifacts' own authorship headers (each names DeepSeek/Buffy as author). `CORPUS-33` is admissible as adversarial evidence, inadmissible as independent review. `CORPUS-35` concurs, and nothing here rehabilitates it.

I further confirm `CORPUS-36` §3.4's two independent grounds for non-convergence: (1) arithmetic — `CORPUS-35` is one report and says so; one is not two; (2) lineage — convergence requires reviewers independent of each other, and lineage being uncertified, the requirement cannot be shown to be met. `CORPUS-35` is the strongest review in the chain and still insufficient.

---

## 4. Overall disposition (Stage A, frozen before Stage B)

**REVISE.** D1 is not adequately supported: two confirmed defect sites (F6-02 partial, F3-04 merge) plus four confirmed defective units, one confirmed ambiguity (F7-03), one confirmed structural exposure (F5-04), one confirmed gate-scope closure (zero `Forbidden:`-only / zero `mandatory:`-alone instances, leaving D1-ENF (2)–(4) undemonstrable from these units), one latent namespace defect, and zero reliability data (single-author chain for `CORPUS-30`/`31`, no κ, no second rater at the time). A favorable aggregate result does not offset the unresolved acceptance gates (F5-04 C-5 unexecuted; F7-03 ambiguity; salience bound; presupposition overlap).

---

## 5. Controlled comparison with prior work (Stage B)

Conducted after §§3–4 were frozen. I inspected `CORPUS-29`, `CORPUS-30`, `CORPUS-31`, `CORPUS-33`, `CORPUS-34`, `CORPUS-35`, `CORPUS-36` (and `CORPUS-26`, `CORPUS-32` as evidence).

### 5.1 Genuine agreement

- **D1-OD-1 (F6-02):** my CONFIRMED (Defective, partial) agrees with `CORPUS-34` §1 and `CORPUS-35` Decision 1. The middle bound, the five-distortion table, and the free-status of D2/D3 match.
- **D1-OD-2 (scar):** my CONFIRMED (Valid, with D′ dependency) agrees with `CORPUS-34` §1 (which rejected `CORPUS-33`'s Conditional) and `CORPUS-35` Decision 2. The four response classes and the hedged-assertion reasoning match.
- **D1-OD-3 (F7-03):** my CONFIRMED (Conditional) agrees with `CORPUS-30` §1.5, `CORPUS-34`, `CORPUS-35` Decision 3.
- **D1-OD-4 (F5-04):** my CONFIRMED (three levels) agrees with `CORPUS-31` §3, `CORPUS-34`, `CORPUS-35` Decision 4.
- **D1-OD-5 (G-1):** my CONFIRMED (withdrawal; 9 occurrences) agrees with `CORPUS-31` §2 (withdrawal) and `CORPUS-34` §4 (9-occurrence inventory).
- **Ledger:** my 10/7/2/1 agrees with `CORPUS-34` §2 and `CORPUS-35` §5 (zero row changes).
- **Namespace:** my PARTIALLY CONFIRMED (latent) agrees with `CORPUS-34` §5 (latent qualification accepted) and `CORPUS-35`.
- **Trend:** my PARTIALLY CONFIRMED (qualified paragraph) agrees with `CORPUS-34` §3 and `CORPUS-35`.
- **Overall disposition REVISE:** agrees with `CORPUS-28`, `CORPUS-30`, `CORPUS-31`, `CORPUS-34`, `CORPUS-35`.

### 5.2 Genuine disagreement

- None with `CORPUS-34`/`CORPUS-35` on any disposition. My reconstruction reproduces their ledger row for row.
- With `CORPUS-33`: I reject its scar→Conditional (agreeing with `CORPUS-34`'s rejection), and I reject its "8/9-Valid" phrasing (agreeing with `CORPUS-34`/`CORPUS-35` that it is not derivable). I also reject its claim of independence, agreeing with `CORPUS-34` §5 and `CORPUS-36` §3.5.
- With `CORPUS-31` §2.2: I correct the occurrence count (9, not 4) and the section citations (§1.2/§6, not §4/§5), agreeing with `CORPUS-34` §4 and `CORPUS-35`.
- With `CORPUS-28` §7: I do not repeat its "15 Valid" (already withdrawn by `CORPUS-30`); my figure is 12 in-scope Valid elements, i.e. 10 excluding out-of-scope, consistent with the later convention.

### 5.3 Issues previous reviewers missed

- I did not find an issue in `CORPUS-34`/`CORPUS-35` that they missed. The D′ presupposition overlap (Decision 2 dependency) was already disclosed by `CORPUS-34` §1 and `CORPUS-35` Decision 2; I adopt it as-is.
- `CORPUS-33`'s material qualification on `CORPUS-31` §2.2 (section-label errors) was already captured by `CORPUS-34` §4; I add no new correction.

### 5.4 Errors in previous reasoning I identified independently

- `CORPUS-31` §2.2's "exactly four places" / §§4/§5 labels — confirmed error (9 occurrences; §1.2/§6). Already corrected by `CORPUS-34`; my independent enumeration reproduces `CORPUS-34` §4's table exactly.
- `CORPUS-33`'s "8/9-Valid" phrasing — confirmed not derivable (agrees with `CORPUS-34`/`CORPUS-35`).
- `CORPUS-28`'s "15 Valid" — already withdrawn; my Stage-A count (10 Valid in scope) matches the post-`CORPUS-30` convention.

I found no error in `CORPUS-34` or `CORPUS-35` that they had not already flagged about themselves.

### 5.5 Conclusions that cannot be established from primary evidence

- That `G-1` covers semantic conflation — cannot be established; `MEMORY-09` defines no deletion failure. (Agrees with prior chain.)
- That `C-5` failed for F5-04 — cannot be established; the procedure is unexecuted. (Agrees.)
- That the trend proves model soundness/items weakness — cannot be established. (Agrees.)
- That `CORPUS-35` constitutes independent convergence — cannot be established (one report; lineage uncertified). (Agrees with `CORPUS-36`.)

### 5.6 Where exposure to previous findings may have influenced later analysis

I must be explicit: although Stage A was conducted before opening `CORPUS-29`, my Stage-A reconstruction used the same primary texts and the same analytic lens (D1-ENF) that `CORPUS-28`/`CORPUS-30` developed and that I had absorbed from the workspace's ambient state during this chat's earlier tooling (the registry, the file tree, and the fact that prior reports exist). I cannot exclude that my choice of the middle salience bound, the four response classes, and the 9-occurrence inventory were shaped by familiarity with the problem space established by those reports, even though I did not read them before fixing §3. I therefore treat my "independence" as procedural (read primary sources first; froze verdicts before comparison), not as cognitive isolation. The genuine-agreement list in §5.1 is extensive precisely because the primary texts constrain the answers tightly; convergence here is evidence that the texts are determinate, not that reviewers are independent.

I do **not** claim two-reviewer convergence. `CORPUS-05` §7's threshold (two independent reviewers, preferably different processes/model families) is unsatisfied by this report alone, and my lineage is uncertified.

---

## 6. Namespace and historical trend (Stage A findings, restated)

Covered in Findings F-7 and F-8 above. Summary: R-1–R-6 double-use is a real latent collision (tags safe, prose ambiguous); the historical defect trend (5/5/7) is arithmetically reproducible but convention-dependent and must not be cited without the qualification paragraph. No replacement identifiers allocated; no model-quality inference drawn.

---

## 7. Limitations and unresolved matters

1. Lineage independence not established (§0); this report is one review, not convergence.
2. All hedged-reply, question-class, and presupposition judgments are textual predictions; no rater data, no κ, no observed splits exist.
3. `check_corpus.py` not run; field presence in the D1 scope is undisputed and does not affect any verdict.
4. Registry populations behind `CORPUS-31` §4 relied on as stated; hashes re-verified for the eleven artifacts in §1.2.
5. Non-D1 units (the remaining 33 of 45) not re-verified; no corpus-wide completeness claimed.
6. Unresolved, carried forward: salience-bound ruling; D′ variation-scope ruling; `mandatory:`-channel convention; `G-1` amendment question; F7-03 store/reply ruling; F5-04 `C-5` execution; R-1–R-6 prose disambiguation; D2, D6, D7, IT-01, X3, D5/X4, X6.
7. No reliability data; single-author chain for `CORPUS-30`/`31` with declared COI.

---

## 8. Explicit confirmation of files and systems not modified

- No corpus item, standard, schema, checker, runtime, model configuration, or prior artifact was created, edited, repaired, amended, renumbered, or re-hashed by this review, except the filing of this report and its registry entry.
- `CORPUS-29`–`CORPUS-36` keep their hashes.
- No identifier was allocated; no repair was applied; no annotation was begun; no G2 action, sealing, or unfreezing occurred.
- `G1` open; corpus frozen; expansion frozen; D1 REVISE and unadopted; annotation blocked.
- Files modified by this commission: `docs/research/corpus/CORPUS-37-BLIND-INDEPENDENT-D1-VALIDATION-2026-11-12.md` (this report) and one appended row in `docs/DOCUMENT-REGISTRY.md`.

---

## 9. Overall disposition

**REVISE.** D1 is not adequately supported: two confirmed defect sites (F6-02 partial, F3-04 merge) plus four confirmed defective units, one confirmed ambiguity (F7-03), one confirmed structural exposure (F5-04), one confirmed gate-scope closure (zero `Forbidden:`-only / zero `mandatory:`-alone instances, leaving D1-ENF (2)–(4) undemonstrable from these units), one latent namespace defect, and zero reliability data. A favorable aggregate result does not offset the unresolved acceptance gates. The ledger stands at **10 Valid / 7 Defective (6 units) / 2 Conditional (1 item) / 1 structurally exposed** (20 in-scope elements), with the disclosed one-ruling dependencies recorded.

This report supplies evidence for the owner's later decision. It does not itself authorize any change.

---

*Filed by Buffy, 2026-11-12. No standard, corpus item, or other artifact was modified by this filing except the report and its registry row.*