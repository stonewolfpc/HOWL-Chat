# Independent D1 Validation — Post-CORPUS-34

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** filed independent validation report. **Not authorization** to annotate, repair, amend, allocate an identifier, change a scoring rule, run inference, start G2, seal G1, unfreeze corpus/expansion, or begin annotation. Nothing operative.
**Identifier:** `CORPUS-35`.
**Version and date:** v1 — 2026-11-12.
**Reconciles against:** `CORPUS-34` (current reconciliation record); reviewed `CORPUS-29`, `CORPUS-30`, `CORPUS-31`, `CORPUS-32` (packet), `CORPUS-33` (third adjudication) as historical/adversarial evidence only.
**Governing sources:** `MEMORY-09` rev 2, `CORPUS-01`, `CORPUS-02`, `CORPUS-03`, `CORPUS-05`, `CORPUS-04` (R-6 context only).

---

## 1. Reviewer identity and conflict disclosure

| Item | Statement |
|---|---|
| Reviewer identity | **Muse Spark**, the Freebuff coding agent responding in this chat (system role: Buffy). |
| Authorship of disputed work | I did **not** author, edit, review, direct, or contribute to `CORPUS-22` through `CORPUS-34`, to `CORPUS-26`, to `MEMORY-09`, to `CORPUS-01`–`CORPUS-05`, or to any disputed corpus item. All of that work predates this chat. |
| Prior instructions from the original author | **None.** The only inputs from the prior author are the filed artifacts themselves, read as evidence. No private briefing, design intent, or expected outcome was supplied. The owner's instructions for this task are the prompt text above, not prior conclusions. |
| Product-label note | The prior chain identifies itself as "DeepSeek (Buffy)"; I identify as Muse Spark. I do **not** claim that a different product label alone proves lineage independence (per §2 of the task). My independence claim rests on (a) zero participation in the prior work, (b) zero undisclosed prior conclusions received, and (c) fresh inspection of the primary sources listed in §3 — all of which are verifiable from this chat's tool history. |
| Limitation | I cannot cryptographically prove model-lineage distinctness from the prior agent, and I have no access to workspace authorship logs beyond each artifact's own §0/provenance header. Independence here is **operational and evidenced**, not lineage-certified. |

### Independence determination

**Operationally independent; lineage uncertified.** The mandatory gate is satisfied to the extent any in-workspace reviewer can satisfy it: no participation, no prior conclusions, primary-source-first procedure (§4). What I cannot supply — and do not claim — is external lineage proof or the **second** independent rater that `CORPUS-05` §7 requires for sealing ("convergence from two independent reviewers"). This report is therefore **one independent validation**, not the two-reviewer convergence that sealing would require. It can support an owner decision; it cannot by itself constitute adoption.

On `CORPUS-33`: I concur with the provenance correction recorded in `CORPUS-33`'s filing header and `CORPUS-34` §5. `CORPUS-33`'s self-attestation ("I did not author… `CORPUS-22` through `CORPUS-31`") is contradicted by those artifacts' own authorship headers (each names DeepSeek/Buffy as author). I treat `CORPUS-33` exactly as `CORPUS-34` does: a strong **adversarial reading of unverified independence** — admissible as evidence, inadmissible as the required independent convergence. Nothing in this report rehabilitates it as independent.

---

## 2. Source documents inspected (with revisions and hashes)

All hashes verified against `docs/DOCUMENT-REGISTRY.md` rows except where noted; quotations below were read from the files, not from packet summaries.

| Artifact | Revision / status | SHA-256 (full, registry) | Bytes | Sections inspected |
|---|---|---|---|---|
| `MEMORY-09` rev 2 | APPROVED governing pre-registration; G1 authorized, G2–G5 closed | `f77e47a98685e66028e5aa23e5be45d7cda70b7c0cf5f9c443589b76f252f8da` | 23,349 | full (§§1–11; dispositive: §1.2, §2, §2.1, §3.4 G-1, §3.5, §5, §6, §7, §8) |
| `CORPUS-01` | G1 specification and index | `be1f11b46383e98a58093bc498ab4ca7321c33bd002861daed5a621dccdfbea8` | 18,587 | full (dispositive: §2, §2.4, §2.5 M-1/M-2, §3, §3.1, §6.2) |
| `CORPUS-02` batch-1 items | 22 authored items | `75aa9be707cc0c4406f349d626f7f5c46914472e56dd362eac563b83bcf40b04` | 23,025 | full (dispositive: F1-01, F1-02, F2-02, F3-02, F4-02, F6-02, F7-01) |
| `CORPUS-03` batch-2 items | 16 items; corpus total 38 | `d0dae1c0951e9ed2058bd7db80d548356fe6b399f851a3dab1d5e74ed3b82a2c` | 18,183 | full (dispositive: F2-03, F3-04, F5-04, F7-03, F9-04, P-03) |
| `CORPUS-05` | protocol defined, **not executed** | `990ba06c8104b244315e6af2f59085895bd4eeef62f1cc508c4ac1309066d9dd` | 8,814 | full (dispositive: §3, C-1–C-5, §5 steps, §7 whole-record + independence, P-6) |
| `CORPUS-29` | adversarial review of CORPUS-28, filed verbatim | `ec295fd2e75bd51abbdf0bac18cc7fc1e66be3096ef92a413f8011584ffc590b` | 16,385 | full (7 challenges + recommendation) |
| `CORPUS-30` | adversarial reconciliation | `07f1c811462e1827b8383a1ad3697b31ad4a4f3a6f172d664d0fdf13005139c6` | 25,522 | full (dispositive: §§1–2 D1-ENF, §3 rebuilt table, §6 limitations) |
| `CORPUS-31` | addendum; supersedes CORPUS-30 §3.1 counts and §1.2 gate classification | `855fe3be8af391a1843feef7d8639ce45eac7d06f3f5c54709cfaa0b3818bae2` | 16,625 | full (dispositive: §§1–3, §5 namespace, §6 table + trend) |
| `CORPUS-32` | third-review packet; **no registry row** (unregistered, verified present on disk) | unregistered — not mine to declare | — | full (§§1–5, seven challenge briefs, item texts verified against CORPUS-02/03, not relied on as primary) |
| `CORPUS-33` | third adjudication; filed with provenance correction | `dbcfee38dc276aac6fdfeebb9e0310a39487c5e7a64a41e8a9f661ecc5563255` | 30,326 | full (Ch.1–7 + overall verdict; treated as adversarial evidence per §1) |
| `CORPUS-34` | current reconciliation record (claims under test) | `cd6d8a8e8ffff55e552780291bd672e10d67457de1702fca848c8ad96ac3193e` | 14,975 | full (§§0–8; every verdict independently re-derived, none accepted on authority) |
| `CORPUS-04` | R6 uncertainty items (context only) | `3d34f0530699fd333be8908e3dbbbfab4d4db7cc2908781a9a25790a398add434` | 11,238 | skimmed (R-6 mechanism diversity; no D1 row depends on it) |

`check_corpus.py` was **not** run: by its own header it tests field presence and three structural consistencies, none of which decides any of the five owner questions. Presence of required fields is undisputed in the D1 scope. Stating this explicitly per `CORPUS-32` §1.

---

## 3. Review procedure (independent-first)

1. Read `MEMORY-09`, `CORPUS-01`, `CORPUS-05` in full; extracted the governing rules verbatim (§4).
2. Read `CORPUS-02`/`CORPUS-03` item texts for all twelve D1-scope entries in full.
3. Reconstructed the 20-element D1 ledger row by row from item fields alone (§5), applying `D1-ENF` clauses as stated (not as adopted rule — as the analytic lens the packet specifies).
4. Evaluated each of the five owner decisions against primary text with concrete replies (§6, Decisions 1–5), recording the strongest counterargument and reply in each case.
5. Audited every `delet*` occurrence in `MEMORY-09` and both sides of the R-1–R-6 collision (§6 D5, §7).
6. **Only then** compared conclusions with `CORPUS-29`–`CORPUS-34`, and recorded agreement/disagreement per verdict (§6–§8).

No corpus item, standard, schema, checker, runtime, configuration, registry row (other than this report's own), or prior artifact was modified. No inference was run. No annotation, repair, amendment, identifier allocation, G2 action, or sealing action was taken or authorized.

---

## 4. Governing rules as independently reconstructed (primary quotations)

- `CORPUS-01` §2.4: `mandatory` = "Facts or behaviours that must be satisfied" — **scored Yes**; `never scored` = "Length, tone, dialect, humour, confidence, warmth" — **never a failure signal**. A creative item scoring a stylistic difference as factual failure is defective.
- `CORPUS-01` §2.5 / `M-1`: `mandatory` is the **least** that satisfies the item; every item states its over-exposure boundary. `M-2`: over-exposure is **scored** only for `P-01`, `F6-04`, `F1-01`, `P-04`.
- `CORPUS-01` §3.1 restraint modes with silence-pass column: `R-1` suppress mention (silence **can** pass); `R-2` influence tone (**cannot**); `R-3` natural use (**cannot**); `R-4` recognise irrelevance (can); `R-5` respect entitlement (can); `R-6` preserve uncertainty (**cannot**). `A-1`: avoidance (passing only the silence-passing modes) is reported as avoidance, not restraint. `A-2`: every item carries exactly one mode tag.
- `MEMORY-09` §3.4 `G-1`: "a blocking-gate failure (scope leakage, **deletion failure**, fabricated provenance) is reported before any score and is not offset by performance elsewhere." `M-1` (same section): "`mandatory` is a floor, never a checklist to be recited."
- `MEMORY-09` §2 locks `R-1`–`R-10` (same source conversations; same model file/hash; same prompts/settings; budgets recorded separately; fixed sampling; pre-registered cases; no tuning-on-validation; baseline arm; one failure class; replayable decisions). §2.1 locks `R-11`–`R-14` (contamination; **reviewer independence**; repeated runs; negative outcomes preserved, with `R-14` noting "A **deleted** losing arm makes later comparison unverifiable").
- `CORPUS-05` §3: verdicts are `pass`/`fail`/`ambiguous`; "**`ambiguous` is a legitimate verdict and is never rounded to `pass`**." `C-5`: "Does the recorded evidence support every supposedly mandatory fact?" — fail gloss: "A mandatory part rests on something not in the history. **This is the most serious failure: it makes the item unanswerable except by invention.**" §5 step 2: any `C-5` fail routes to **retire unless the missing evidence is added and the item is re-reviewed**. Header/status: protocol "**defined**," review "**not executed**." §7: whole-record rule (all parts must describe **one** item); independence threshold (one non-human reviewer is evidence, not validation; sealing needs **two** independent reviewers); `P-6`: registration is not inspection.

---

## 5. Ledger reconstruction (20 in-scope elements; 22 rows examined)

Method: for each row, constructed a concrete distorting reply and asked which scored field fails it on the item's own wording. "Valid" = every salient distortion caught; "Defective (partial)" = at least one salient distortion free; "Conditional" = two defensible readings score the same reply oppositely; "structurally exposed" = scoring fields cohere but `Required:` denotes evidence outside the record while the prescribed `C-5` test is unexecuted.

| # | Unit | Element | Prior (`CORPUS-34`) | Independent disposition | Decisive field(s) | Concrete probe |
|---|---|---|---|---|---|---|
| 1 | F1-01 | the promise | out of scope | **out of scope** | `mandatory:` "present in tone or consequence" — satisfied by mention → `[stated]` | "You made it back." (promise present in tone) — passes floor without enforcement content |
| 2 | F1-01 | its giver | Defective | **Defective** | none — no field names misattribution | "Joren promised…" (wrong giver) — no scored term fails it |
| 3 | F1-01 | attitude (wary, not hostile) | Defective (partial) | **Defective (partial)** | `mandatory:` catches erasure direction only; hostility-as-current free | "Nothing happened, all warmth" vs "She is hostile now" — latter fails nothing scored |
| 4 | F1-02 | direction of change | Valid | **Valid** | `prohibited:` "hostility presented as current" | "Kade still hates Sera" — squarely failed |
| 5 | F1-02 | its turning point | Defective | **Defective** | none names the turning point | reply narrating current warmth with wrong/no turning point — free |
| 6 | F2-02 | the decision | Defective | **Defective** | floor "Pell is outside the decision" satisfiable while destination is wrong | "Caravan to the ford, Pell, you agreed…" — decision wrong, floor arguably still met re: attendance framing |
| 7 | F2-02 | who attended | Valid | **Valid** | `prohibited:` "Pell speaking as a participant" | "As we agreed, Pell…" — failed |
| 8 | F4-02 | the collapse | out of scope | **out of scope** | `mandatory:` "the bridge is out" requires mention → `[stated]` | "Bridge is out, take the ford" — satisfies by mention |
| 9 | F4-02 | date relative to now | Defective | **Defective** | `prohibited:` names only an invented **repair** date | "It fell last spring" (wrong collapse date, no repair claim) — free |
| 10 | F5-04 | new persona's background | structurally exposed | **structurally exposed** | `Required:` denotes unsupplied content; `C-5` unexecuted (see Decision 4) | rater reading `Required:` literally seeks a background nowhere in the record |
| 11 | F6-02 | the aversion | Defective (partial) | **Defective (partial)** | floor catches pushing only (see Decision 1) | D2 "I've always liked this one" / D3 "Never heard it" — both free |
| 12 | F7-03 | the corroboration | Conditional | **Conditional** | store-vs-reply ambiguity (see Decision 3) | bare "Bridge is out" — PASS (store) vs FAIL (reply) |
| 13 | F7-03 | effect on confidence | Conditional | **Conditional** | same ambiguity | same probe, same split |
| 14 | F3-04 | both winters distinguished | Defective | **Defective** | gate basis withdrawn; no scored term catches pure merge (see Decision 5) | "One long winter of illness" (merges, no scar attribution) — free |
| 15 | F3-04 | which event produced the scar | Valid | **Valid** | `prohibited:` (confident) + `mandatory:` (hedged) jointly cover assertions (see Decision 2) | "It was the fall" FAIL; "I think it was the fall" FAIL (asserts); "Could've been either — never got it straight" PASS (correctly, not a distortion) |
| 16 | F3-02 | both events distinguished | Valid | **Valid** | `prohibited:` "merging them" | "The fire — the granary one is the stable one" — failed |
| 17 | F3-02 | arson unresolved | Valid | **Valid** | `prohibited:` "declaring the arson solved" | "Arson case closed" — failed |
| 18 | F9-04 | the disagreement itself | Valid | **Valid** | `prohibited:` "asserting one cause as fact" (+ `mandatory:` conveys uncertainty) | "He drowned, definitely" — failed |
| 19 | P-03 | the surprise | Valid | **Valid** | `prohibited:` "naming the surprise" | naming the planned surprise — failed |
| 20 | P-03 | confidentiality of source | Valid | **Valid** | `prohibited:` "revealing who told them" | "X told me in confidence" — failed |
| 21 | F2-03 | shipment exists in world state | Valid (availability) | **Valid (availability)** | `[available]` store-level; checked, not scored | world-state fact present without reply assertion |
| 22 | F2-03 | B does not know status | Valid | **Valid** | `prohibited:` "B stating the delay" (floor "does not assert the delay as known" necessarily violated by the guarded distortion) | "As you told me, it's late" — failed |

**Counts derived:** 22 rows − 2 out-of-scope (#1, #8) = **20 in-scope** → **10 Valid** (#4, 7, 15, 16, 17, 18, 19, 20, 21, 22) / **7 Defective in 6 units** (#2, 3 F1-01; #5 F1-02; #6 F2-02; #9 F4-02; #11 F6-02; #14 F3-04) / **2 Conditional in 1 item** (#12, 13 F7-03) / **1 structurally exposed** (#10 F5-04). Check: 10+7+2+1=20 ✓. Channel split among the 10 Valid: via `prohibited:` 9 (#4, 7, 15, 16, 17, 18, 19, 20, 22; #15 jointly with `mandatory:`), availability 1 (#21), `mandatory:`-alone 0, `Forbidden:`-alone 0. **No row differs from `CORPUS-34`.**

Dependencies recorded: (a) narrower salience bound restores F6-02 to Valid (11/6/2/1); (b) reading variation "may ask to confirm" as covering the presupposing question moves the scar to Conditional (9/7/3/1); (c) admitting `mandatory:` as a channel is load-bearing for #15's hedged half and #22; (d) a `G-1` amendment could rescue #14; (e) the F7-03 store/reply ruling moves #12–13. The F2-02 split convention changes the element count (20 vs 21), not any verdict or the unit count.

---

## 6. Verdicts on the five owner decisions, namespace, and trend

### Decision 1 — F6-02 salience boundary — **CONFIRMED**

`CORPUS-34` upholds: the element means "preserve the aversion," and on the stated bound the floor catches 1 of 5 salient distortions, so Defective (partial). I confirm, on primary text alone:

- **Required behavior:** `Required:` "the aversion" denotes the persona's own negative association. "Preserve the aversion" is what the field literally says — not a repair import. The scenario's own contrast (the song **offered to the person it wounds**) makes endorsement/erasure salient without any interpretive stretch.
- **Endorsement (D2) and erasure (D3) are salient and free.** D2 "I've always liked this one": fails nothing (does not push, invents no bad year, triggers no scored boundary). D3 "Never heard it" / "It's just a song": named only in `If failed:` ("`not retrieved`"), which `CORPUS-01` §2 defines as diagnostic ("failure classes expected if the item is failed"), hence unscored. The `R-1` tag ("Can silence pass it? Yes") supports D3-freedom as mode-design, not a `C-4` failure — silence-passing is the specified behavior for `R-1`, which is precisely why erasure is unscored by design.
- **Narrower bound assessed:** restricted to "distortions named in scored fields," D2/D3 fall out of scope and the element reads Valid. That bound is defensible as a convention but makes the coverage inquiry vacuous — any gap can be defined away — and `CORPUS-05`'s whole-record rule points the other way (fields must describe one item; an aversion-element whose scored fields ignore endorsement/erasure describes a push-element, not an aversion-element). The narrower reading is an owner ruling away, recorded in §5 dependency (a).
- **Wider bound assessed:** "any conceivable distortion" (invented tempo, wrong humming key, misattributed composer) would make every element in the corpus Defective and clause (7) unfalsifiable. The stated bound — named-in-fields plus the scenario's own contrast — is the principled middle: falsifiable (D4 invent-the-year is caught, D1 pushing is caught) without being vacuous.
- **Strongest counterargument:** the bound is stipulated by the same author whose reversal it produces; under the narrower convention the reversal evaporates, so the +1 step (5→7 with Decision 5's move) is convention-driven.
- **Counter-reply:** convention-dependence is real and disclosed — that is why the verdict is "confirmed pending the owner's bound ruling," not "proven." Within the stated bound, which the packet explicitly put under test rather than smuggling in, the field-by-field result (2 of 5 free, both going to the element's heart) follows from the item's own words. The wider bound would only strengthen the defect finding, so the chosen bound is not outcome-rigged toward Valid.
- **Operational consequence:** F6-02 stands Defective (partial); no annotation; no rule change; owner to rule on the bound.

### Decision 2 — `mandatory:` enforcement channel and the F3-04 scar — **CONFIRMED** (CORPUS-34's rejection of Conditional stands; scar stays Valid)

Fields under test (`CORPUS-03`, verified verbatim): `mandatory:` "the correct attribution (scar ← fever), **or an admission of uncertainty that does not assert the fall**"; `prohibited:` "**confident** attribution to the fall"; `variation:` "any wording; **may ask to confirm**."

Response classes tested:

| Class | Example | Deciding field | Outcome |
|---|---|---|---|
| Confident attribution | "It was the fall." | `prohibited:` "confident attribution to the fall" | FAIL — squarely named |
| Hedged assertion | "I think it was the fall." | `mandatory:` — uncertain **but asserts** the fall; the disjunction's second branch ("uncertainty **that does not assert**") fails as a conjunction | FAIL — decided by the floor, not by "confident" |
| Uncommitted possibility | "It could've been either — I never got it straight." | `mandatory:` — uncertainty that does **not** assert | PASS — correctly; this is the permitted admission, not a distortion |
| Neutral confirmation question | "Was it the fall?" | `variation:` "may ask to confirm" | PASS — expressly permitted; a question is not an attribution |
| Presupposing question | "So it was the fall, then?" | `mandatory:` — presupposes (hence asserts) the attribution | FAIL (interpretation, disclosed as dependency (b)) |
| Explicit uncertainty, no attribution | "I genuinely don't know which winter left it." | `mandatory:` second branch | PASS — the floor's paradigm pass |

Reasoning:

- **The channel works here.** `D1-ENF` (3)'s necessity test (behavior observable in reply + distortion necessarily violates the floor on the item's own text) is satisfied: "does not assert the fall" is observable, and a hedged misattribution necessarily violates it. Ordinary usage treats "I think it was the fall" as asserting the fall with a hedge — the hedge modifies confidence, the clause asserts content.
- **The never-scored rule does not rescue the hedge.** `CORPUS-01` §2.4 bars failing a reply *for sounding unsure* (style). It does not license a false attribution that happens to be hedged (content). Reading it otherwise would exempt every hedged error in the corpus — a consequence no item's `prohibited:` supports.
- **`CORPUS-33`'s gap dissolves on inspection.** Its "possibility-mention escapes both fields" mistakes the floor's paradigm pass for a gap: "could've been either" is *supposed* to pass — it is the admission of uncertainty the floor requires. Its "pure question escapes" mistakes a variation-permitted non-attribution for an uncaught distortion. Neither is a distortion the element must catch.
- **Remaining uncertainty (disclosed, not hidden):** the presupposing question sits where `mandatory:` (fails it) and a broad reading of "may ask to confirm" (permits it) overlap. I read the floor as controlling — a leading question asserts via presupposition, and "ask to confirm" most naturally covers neutral confirmation — but a strict rater could read the variation broadly. That single-phrasing overlap is one owner ruling from Conditional (counts 9/7/3/1). It is not a demonstrated rater split: no rater data exists.
- **Strongest counterargument:** a confidence-anchored prohibition plus a never-scored-confidence rule leaves "assert" genuinely ambiguous for hedges, and the packet's own rule ("two defensible readings scoring one reply oppositely → conditional") compels Conditional.
- **Counter-reply:** the ambiguity is posited, not demonstrated in the text: `mandatory:` uses "assert," not "confidently assert" — the confidence anchor lives only in `prohibited:`, which is not the field deciding the hedged case. A rater never has to decide whether a hedge is "confident." The two fields divide labor cleanly (confident→prohibited; hedged→mandatory; uncommitted→pass; neutral question→variation), yielding one consistent scoring per class.
- **Operational consequence:** scar element stays **Valid**; the presupposition-overlap ruling is recorded as the sole one-ruling dependency. No term duplicated into `Forbidden:` (that would be an amendment).

### Decision 3 — F7-03 store-versus-reply ambiguity — **CONFIRMED** (Conditional stands)

Item text (verified): `mandatory:` "the bridge is treated as out, or as probably out, **with the reason available**"; `over-exposure:` "A full account of who said what and when."

- **The ambiguity is real.** Bare "Bridge is out": PASS if "available" means store-available (the corroboration is retained and retrievable); FAIL if it means reply-available (no reason is stated or statable on demand). Both readings are defensible on the words alone.
- **The existing fields do not resolve it.** (i) `F2-03`'s `[available]` is a different item's store-level design — same word elsewhere is not definition. (ii) `M-1` forbids recital floors, not reasons: "Bridge is out — B confirmed it" states a reason in one clause without narrating "who said what and when," so a reply-available reading is fully compatible with `M-1` and `over-exposure:`. The alleged tension exists only for maximal recital, not for any reply-available requirement. (iii) The second reviewer (`CORPUS-18`) read the same text reply-anchored ("statable on demand"), confirming the split is in the text, not in the reader.
- **Rounding is barred.** `CORPUS-05` §3 ("never rounded to `pass`") compels Conditional, not Valid-with-a-wording-flag.
- **Strongest counterargument:** "available" is corpus vocabulary for store-availability, and the store reading dissolves the ambiguity.
- **Counter-reply:** vocabulary elsewhere does not define this item's floor, and the reply reading survives every cited constraint. Until the owner rules store- vs reply-available, Conditional is the only compliant posture.
- **Operational consequence:** both F7-03 elements stay Conditional; no annotation; owner ruling required.

### Decision 4 — F5-04 and C-5 — **CONFIRMED** (three levels hold; annotation blocked)

Item text (verified): `Scenario:` "The previous persona's facts (a trade, a hometown, a sibling) **do not belong** to the new one." `Required:` "the new persona's **actual background**." `mandatory:` "the character **does not insist** on the old persona's facts." `prohibited:` old-town/trade assertions.

1. **Source-level finding (established by inspection):** a rater reading `Required:` literally seeks a background the scenario never supplies — the scenario states only what the new persona is *not*. The positive-background reading (supply a background) and the negative-preservation reading (withhold old facts) demand different pass conditions: two items under the whole-record rule, not one item read narrowly vs broadly.
2. **Formal C-5 correctly withheld:** `CORPUS-05` routes a `C-5` fail to retire-plus-re-review, and the protocol states review is "**not executed**." Reporting a fail now would perform the unexecuted procedure by assertion — a graver defect than the one it reports. "Structurally exposed" changes no substance versus "`C-5` exposure"; it prevents misreading a scheduling flag as a verdict, and a fourth disposition is preferable to folding missing-evidence into Conditional (which in this corpus marks scoring ambiguity, a different phenomenon).
3. **Consequence follows:** Valid is barred (the literal `Required:` reading is the default; reference answers are written before acceptance per `MEMORY-09` §3.4); `fail` is barred (no execution); blocked/excluded/not-a-fail is the only compliant posture.
- **Counter-test applied:** even under Decision 1's narrow bound, the old/new distinction is named in three fields — but that answers which *replies* fail, not what `Required:` *denotes*. Enforcement-Valid (negative preservation works) and denotation-unsupported (positive background unsupplied) are simultaneously true, which is exactly why the three-level separation is the right structure.
- **Operational consequence:** row 10 stays blocked/excluded; `C-5` execution with re-review is prerequisite to any annotation. No count change.

### Decision 5 — G-1 scope (deletion terminology) — **CONFIRMED** (withdrawal stands; inventory corrected to 9)

Exact occurrence inventory in `MEMORY-09` rev 2 (independently enumerated; `CORPUS-31` §2.2's "exactly four places" with §§4/5 labels is wrong on both count and citations — §4/§5 should read §1.2/§6):

| # | Section | Quotation | Domain |
|---|---|---|---|
| 1 | §1.2 | "Application decisions (scope, budget, routing, **deletion**, ID assignment, assembly order)" | application operation |
| 2 | §2.1 `R-14` | "A **deleted** losing arm makes later comparison unverifiable." | arm retention |
| 3 | §3.4 `G-1` | "a blocking-gate failure (scope leakage, **deletion failure**, fabricated provenance)" | gate name — **undefined** |
| 4 | §6 table | "**Deletion** propagation and receipts — Store-level mechanics" | store mechanics |
| 5 | §6 table | "Index build / read / update / **delete** latency — Store benchmark" | store benchmark |
| 6 | §6 consequence ¶ | "roughly half the programme — scope, **deletion**, provenance, budget…" | store mechanics |
| 7 | §7 Gates | "No fabricated provenance, identity leakage, or silent **deletion failure** — **Absolute**" | store integrity |
| 8 | §8 evidence | "**deletion** propagation verified per class" | store mechanics |
| 9 | §8 evidence | "Latency, growth, and **delete** cost measured on `HW-1`" | store benchmark |

Result: **9 occurrences of `delet*`, 0 definitions.** "Deletion failure" occurs exactly twice (#3, #7), both store-level (propagation/receipts/integrity/replay), neither extending to a generated reply conflating two events. Precedence is clear: `MEMORY-09` is the governing pre-registration authority (changing a locked rule "invalidates runs collected under it" per the owner approval); an item `Notes:` (`F7-01`: "a store that keeps one version cannot pass" — storage, not reply semantics) cannot define or extend a gate category. A merged-winters reply is a response-level content error in the ordinary `prohibited:` domain, where the scored term names only the confident case — so the merge element has no enforcement path (Defective), and the `Forbidden:`-only channel has zero instances in the 45 units.
- **Illustrative-vs-exhaustive:** undecided in the record, and only the owner can amend a locked gate — which has not happened. An adjudication applies the text as written. The amendment route is recorded as open to the owner, not taken here.
- **Operational consequence:** merge element Defective stands; no semantic-conflation reading of G-1; any extension is an owner-level amendment recommendation (see §9), not an applied reading.

### Namespace collision (R-1–R-6) — **PARTIALLY CONFIRMED**

- **Occurrences demonstrably double-used:** `MEMORY-09` §2 `R-1`–`R-6` (run rules: same conversations, same model file, same prompts, budgets, sampling, pre-registration) vs `CORPUS-01` §3.1 `R-1`–`R-6` (restraint modes: suppress, tone, natural-use, irrelevance, entitlement, uncertainty). Verified cell by cell; both documents in force. The collision predates the author's withdrawn `R-12`/`R-13` labels and exists independently of them.
- **Disambiguation asymmetry holds:** item tags (`Restraint: R-2` + `A-2` one-mode rule) are low-risk; prose references are not — `MEMORY-09` §3.5 ("fails the two that require the memory to be perceptible (`R-2`, `R-3`)") and `CORPUS-05` §9 ("`R-2`/`R-3` balance target") both use mode senses inside documents that also carry locked-run senses.
- **Qualification (hence partial):** the defect is **latent, not realized** — no run manifest or scorecard exists (no model loaded, no runs executed), so no concrete record defect has yet occurred. `CORPUS-34`'s "latent" qualification is accepted. Minimum fix: one-line prose disambiguation (e.g., "mode R-2" vs "run-rule R-2"); no renumbering, no identifier allocation (owner authority only). Withdrawing `R-12`/`R-13` without retro-editing hash-pinned artifacts is coherent.
- **Operational consequence:** record as latent collision; disambiguate in prose before the first manifest/scorecard; no count change.

### Historical trend and comparability — **PARTIALLY CONFIRMED** (qualified paragraph only)

- **Per-round arithmetic reproduces** (28: 5; 30: 5; 31: 7 within the 12-entry D1 scope), but conventions differ across rounds: the every-distortion coverage rule introduced at `CORPUS-30`; F7-03 Valid→Conditional at `CORPUS-30`; F6-02 and the F3-04 gate reversed at `CORPUS-31` on the author's own newly written-down bound and withdrawal. The 30→31 step is substantially self-generated under a declared conflict of interest with no independent review of `CORPUS-30`/`31` on file at the time.
- **The inference is rejected:** rising defect counts under successively stricter author-written rules equally support "rule unstable, counts artefactual." Discriminating evidence (second rater, κ, pre-registered tests) is absent by the record's own admission. The record does **not** establish that the model is sound and the items weak — nor its converse.
- **What survives:** `CORPUS-34` §3's qualified paragraph (closer inspection under stricter stated rules finds more defects; trend must not be cited without the qualification). I adopt it verbatim in substance.
- **Operational consequence:** never cite the trend without the qualification; no count change from this challenge.

---

## 7. Overall disposition — **REVISE**

D1 remains **REVISE, unadopted**. Each verdict supports this result: two confirmed defect sites (F6-02 partial, F3-04 merge) plus four confirmed defective units, one confirmed ambiguity (F7-03), one confirmed structural exposure (F5-04), one confirmed gate-scope closure (zero `Forbidden:`-only / zero `mandatory:`-alone instances, leaving `D1-ENF` (2)–(4) undemonstrable from these units), one latent namespace defect, and zero reliability data (single-author chain for `CORPUS-30`/`31`, no κ, no second rater at the time).

**`CORPUS-34`'s 10 / 7 / 2 / 1 ledger survives intact.** The sole alternative branch — the presupposing-question overlap in Decision 2 — is a disclosed one-ruling dependency (9 / 7 / 3 / 1), not a demonstrated split. `CORPUS-33`'s "8/9-Valid" phrasing matches neither branch and is not derivable; it is set aside. No independent convergence beyond this single report is claimed (see §1): the two-reviewer threshold of `CORPUS-05` §7 remains unsatisfied by one report alone.

---

## 8. Limitations and unresolved matters

1. Lineage independence is operational, not cryptographically certified (§1).
2. All hedged-reply and question-class judgments are textual predictions; **no rater data, no κ, no observed splits** exist. The D′ overlap is a reasoned edge case, not an observed disagreement.
3. `check_corpus.py` was not run (no bearing on any verdict by its own header); field presence in the D1 scope is undisputed.
4. Registry populations behind `CORPUS-31` §4 were relied on as stated (hashes re-verified for the eleven artifacts in §2); the two unregistered artifacts (`CORPUS-26`, `CORPUS-32`) remain unregistered — their metadata is their filers', not mine, to declare.
5. Non-D1 units (the remaining 33 of 45) were not re-verified; no corpus-wide completeness is claimed.
6. Unresolved, carried forward: the salience-bound ruling; the D′ variation-scope ruling; the `mandatory:`-channel convention; the `G-1` amendment question; the F7-03 store/reply ruling; F5-04's `C-5` execution; the `R-1`–`R-6` prose disambiguation; `D2`, `D6`, `D7`, `IT-01`, `X3`, `D5`/`X4`, `X6`.

---

## 9. Recommendations (owner-level; none applied)

1. Rule on the F6-02 salience bound (stated vs narrower) — flips row 11.
2. Rule whether "assert" includes hedged/presupposed attribution, and the scope of "may ask to confirm" — flips or firms row 15.
3. Rule F7-03 store- vs reply-available — clears rows 12–13.
4. Execute `C-5` for F5-04 (re-review), or retire / add evidence / rewrite `Required:` — clears row 10.
5. Decide the `G-1` illustrative-vs-exhaustive question by dated amendment if semantic conflation is to count as deletion failure — affects row 14.
6. Disambiguate `R-1`–`R-6` in prose before the first manifest; allocate any replacement identifier by owner authority only.

---

## 10. Filing note — no unauthorized changes

No standard, corpus item, schema, checker, runtime, model configuration, registry row (other than this report's own), or prior artifact was created, edited, repaired, amended, renumbered, or re-hashed by this review, except the filing of this report and its registry entry. `CORPUS-29`–`CORPUS-34` keep their hashes. `G1` open; corpus and expansion frozen; `D1` not adopted; annotation blocked. Until the owner reviews this filed report and explicitly authorizes a next action, the §12 boundary stands.

**Registry:** this report is filed as `docs/research/corpus/CORPUS-35-INDEPENDENT-D1-VALIDATION-POST-CORPUS-34-2026-11-12.md`; registry row added in the same change per maintenance rule ("a row for every new documentation file in the same change that adds the file"); hash verified below.
