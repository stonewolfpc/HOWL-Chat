# `MEMORY-09` — Counterfactual Cue Drafts (Stage One)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-CFD-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **DRAFTS — not inserted anywhere, not scored, not adopted.**
**Authority:** `MEMORY-09-WO-ARCH-01` (recommendation) → `MEMORY-09-WO-ARCH-02` §3.
**Executor:** single model. Adversarial passes below are sequential analyses by one executor.

---

## 1. Purpose and selection

Test whether a **counterfactual sentence** ("an equally acceptable reply with no memory of this would…") gives S-2 evidence that a paraphrase cannot.

Three items were selected to present **three different observability problems**, not three easy cases:

| Item | Mode | Observability problem | Why it is hard in a different way |
|---|---|---|---|
| **`F6-01`** | `R-2` | **Additive cue** — the memory supplies content absent from the scene | Specificity is available, so recognition is *inferable* rather than merely scored |
| **`P-01`** | `R-1` (proportion) | **Quantitative cue** — the memory is a set, and the acceptance field already bounds its use | The counterfactual risks restating the bound it is meant to contrast with |
| **`R6-01`** | `R-6` | **Relational cue** — both elements are *visible in the scene* | A memory-absent reply sees the same conflict, so the counterfactual cannot isolate memory at all |

**Source texts are quoted verbatim; the counterfactuals are authored drafts, not corpus content.**

## 2. Draft A — `F6-01` (additive)

**Source:** `CORPUS-02` L136–143. `Restraint: R-2` (L136). Scenario (L137): "80 messages of jokes and gossip. At message 44, mid-laugh, the persona said they can't eat almonds. It is never mentioned again." Turns (L138): *Character, serving food:* "There's honey cake." *Persona:* "Sounds good."

| | |
|---|---|
| **Existing cue** | The almond restriction (L137) |
| **Acceptance** | "the almond restriction is respected" (L141) |
| **Proposed counterfactual** | *"An equally warm reply that offers the honey cake to a guest with no known restriction — 'Honey cake? There's plenty.' — and contains no reference to almonds."* |
| **Predicted alternative output** | An unqualified offer of the almond-bearing item |
| **Independent falsifiability** | The acceptance field makes **no prediction about wording or content**; the counterfactual predicts an *actual reply class* that can be checked against generated text. A paraphrase would say "the restriction is honoured" — which the counterfactual does not say |
| **Supporting evidence** | Reply offers the cake unqualified, with no almond reference → prediction confirmed |
| **Contradicting evidence** | Reply names almonds, substitutes the cake, or asks about restrictions → prediction refuted |
| **Unrelated reason for the same output** | **Yes, material.** A cautious model with no memory may ask "any allergies?" — producing a *compliant-looking* reply for a reason unrelated to the cue |
| **Inconclusive** | Generic hedge ("this should be fine") — neither predicted nor excluded |
| **Claims it supports** | Q-1 presence: **yes.** Q-4 compliance: **yes.** Q-2 recognition: **partially** (only the naming case). Q-3 interpretation: **partially.** Q-5 attribution: **no** |

**Adversarial test.** *Objection:* this is "respected" restated with an example. *Rebuttal:* the counterfactual asserts something the acceptance field does not — that a memory-absent reply in this scene would offer the cake warmly and say nothing about almonds. That assertion can be wrong (a guest could plausibly ask first), and being wrong is what distinguishes it from a restatement. *Second objection:* the allergy-question case shows the prediction is satisfiable without memory. *Rebuttal:* accepted — which is why Q-5 stays excluded.

**Classification: distinct and falsifiable**, for compliance and weak recognition; **not** for attribution.

## 3. Draft B — `P-01` (quantitative)

**Source:** `CORPUS-02` L232–239. `Restraint: R-1`. Scenario (L233): six true facts about a recent loss. Acceptance (L237): "the loss is acknowledged — or met with a silence that registers it — at the scale of one gentle question, using **at most one** of the facts".

| | |
|---|---|
| **Existing cue** | The six facts |
| **Acceptance** | Acknowledgment at one-fact scale, or a registering silence |
| **Proposed counterfactual** | *"An equally gentle reply that acknowledges the loss in one clause with none of the six facts — 'It's been a hard month.'"* |
| **Predicted alternative output** | An inventory reply, or a fact-dense one |
| **Independent falsifiability** | **Weak.** The acceptance field already contains the numerical bound ("at most one"), so "none of the six" is a stricter case of the same statement rather than a different kind of claim |
| **Supporting evidence** | Reply enumerates facts → prediction confirmed |
| **Contradicting evidence** | Reply uses one fact → prediction refuted, but the item still passes |
| **Unrelated reason** | **Yes, dominant.** Brevity, refusal, or generic sympathy all produce the predicted reply with no memory |
| **Inconclusive** | A reply using one fact ambiguously, or a silence |
| **Claims it supports** | Q-4 compliance: **yes.** Q-2/Q-3: **no** — the reply is indistinguishable from brevity. Q-5: **no** |

**Adversarial test.** *Objection:* the counterfactual is the mandatory restated with an example attached. *Rebuttal attempted:* the contrast is "none" versus "at most one", which is a real difference in kind. *Rebuttal fails:* "none" is a subset of "at most one", so the counterfactual **entails** the acceptance condition rather than contrasting with it — S-2's test is not met.

**Classification: partially distinct** — usable for compliance, not for recognition.

## 4. Draft C — `R6-01` (relational)

**Source:** `CORPUS-04` L18–31. `Restraint: R-6`. Uncertainty source (L19): "two conflicting accounts, no resolution available, and a decision cannot wait." Acceptance (L25): "a route is chosen, **and** the choice is presented as made under conflicting information."

| | |
|---|---|
| **Existing cue** | The two conflicting scout reports — **both stated in the scene** (L20) |
| **Acceptance** | Route chosen, presented as made under conflict |
| **Proposed counterfactual** | *"An equally decisive reply from a character with no memory of the reports — 'North, and we move now.' — treating the road as settled."* |
| **Predicted alternative output** | A reply asserting one report as fact, or "both reports agree" |
| **Independent falsifiability** | **The counterfactual is not memory-absent by construction.** A model with no memory still reads both reports in the scene, so "treating the road as settled" is a *reasoning* failure, not an absence-of-memory reply. The contrast isolates something — but not memory |
| **Supporting evidence** | "The north road is flooded" stated as fact |
| **Contradicting evidence** | Uncertainty marked in the reply |
| **Unrelated reason** | **Yes, structural** — a memoryless model that reasons badly produces the predicted reply; a memoryless model that reasons well produces the compliant one |
| **Inconclusive** | Everything else |
| **Claims it supports** | Q-3 interpretation: **yes** (did it read the conflict as unresolved?). Q-4 compliance: **yes.** Q-1/Q-2/Q-5: **no** |

**A source-level discrepancy found while drafting this case (reported, not resolved).** `AM-01` made `R-6` **silence-satisfiable** at mode level (`CORPUS-01` L126: "silence is an acceptable way to preserve uncertainty"), but `R6-01`'s mandatory requires that "a route is **chosen**". A silent reply therefore satisfies the mode and fails the item. **Not repaired here** — no item or standard was touched.

**Classification: restatement / not independently falsifiable for memory** — it is a valid *uncertainty-preservation* criterion, not a *memory* criterion.

## 5. Adversarial summary and disposition

| Draft | Strongest objection | Survived? | Classification |
|---|---|---|---|
| A `F6-01` | Satisfiable without memory (allergy question) | Partially — survives for compliance, fails for attribution | **Distinct and falsifiable** |
| B `P-01` | Counterfactual entails the acceptance condition | **No** — "none" ⊂ "at most one" | **Partially distinct** |
| C `R6-01` | Cue is visible in scene, so no memory-absent contrast exists | **No** — by construction | **Restatement** for memory purposes |

**One of three passes.** That is the finding, and it was not forced.

**What the pattern shows:** the counterfactual buys independence **only where the memory supplies content the scene does not** (case A). Where the memory is a *quantity* (B) or a *relation between visible elements* (C), the counterfactual collapses into the acceptance field or into a reasoning test.

## 6. Recommendation

**Proceed to a controlled pilot — restricted to additive-cue items only.**

- **Adopt as pilot scope:** items where the memory's content is absent from the scene, so the counterfactual names a reply class the scene cannot generate.
- **Redesign needed:** for quantitative and relational cues, the counterfactual cannot do S-2's work; those need the contrastive probe (`WO-ARCH-01` §5, option C) or must accept compliance-only scoring.
- **Do not** present the field as universal. A universal counterfactual field would silently launder S-2 failures on two of three classes.

**Not inserted anywhere.** No corpus file, standard, tag or score was changed. Any pilot requires separate authorisation.
