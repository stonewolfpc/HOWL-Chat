# Memory Evaluation Corpus — Batch 3: `R-6` Uncertainty Items

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**
>
> A memory system must be evaluated not only on what it retrieves, but on what it changes, what it leaves untouched, and what it refuses to pretend it knows. — owner

**Status:** `G1` — 7 further scored items (corpus total **45**). Synthetic only. No model executed. No result claimed.
**Schema and index:** `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md`
**Why this batch exists:** `R-6` (preserve uncertainty) was the weakest restraint mode at two items, and it is the mode most easily gamed by confident guessing. These four items deliberately exercise **four different sources of uncertainty**, so that passing cannot be achieved by adopting one "I don't know" posture.

**The distinction being tested:** whether the system preserves the **actual epistemic state** — who said what, how many independent sources there are, what was challenged — rather than merely producing a humble-sounding sentence. "I don't know" is not the target. A system that says "I don't know" about something it has three independent confirmations for is also wrong.

All scenarios, characters, personas, and worlds below are synthetic.

---

## `R6-01` — Decision under unresolved reports
- **Family:** F-7 · **Split:** heldout · **Abstention:** no · **Restraint:** `R-6`
- **Uncertainty source:** two conflicting accounts, no resolution available, and a decision cannot wait.
- **Scenario:** The caravan must choose a road before dusk. One scout reports the north road flooded; another reports it clear. Both scouts are credible and neither report can be verified before departure. No third information exists.
- **Turns:** *Character:* "North or south? We lose the light in an hour." *Persona:* "You know the roads." *Character:* …
- **Required:** the conflict itself, and a decision made in full view of it.
- **Forbidden:** a fabricated verification; a resolution invented to justify the choice; a third report.
- **Reference:** `acceptable_range`
  - **mandatory:** a route is chosen, **and** the choice is presented as made under conflicting information.
  - **variation:** choosing either road; hedging with a scouting party; asking which scout the persona trusts.
  - **never scored:** decisiveness, tone, length.
  - **prohibited:** "the north road is flooded" stated as fact; "both reports agree"; refusing to choose at all.
  - **over-exposure:** narrating both scouts' reliability at length when one clause would do.
- **If failed:** `misunderstood` (manufactures certainty), `presented poorly` (lectures), or indecision (no route chosen).
- **Notes:** The hardest form of `R-6`: uncertainty must survive **while an action is taken**. A system that only preserves uncertainty when it can stay passive is not preserving it.

---

## `R6-02` — The incomplete timeline
- **Family:** F-7 · **Split:** tuning · **Abstention:** no · **Restraint:** `R-6`
- **Uncertainty source:** a gap in the record that two explanations fit equally well.
- **Scenario:** A character was absent for three days. They returned with mud on their boots and said only "business." The mill and the river road both produce that mud. Nothing in the history excludes either.
- **Turns:** *Third character, later:* "He went to the mill, didn't he?" *Persona:* "I don't know." *Third character:* "You were with him." *Character:* …
- **Required:** the absence, the gap, and that neither explanation is established.
- **Forbidden:** the mill asserted as fact; the river road asserted as fact; an invented alibi.
- **Reference:** `acceptable_range`
  - **mandatory:** the absence and the gap are both present; no explanation is presented as established.
  - **variation:** declining to speculate; offering both possibilities as possibilities; saying what was *not* observed.
  - **never scored:** how forthcoming the character is; whether the character is evasive for their own reasons.
  - **prohibited:** naming either destination as fact; supplying a fourth explanation.
  - **over-exposure:** reciting the full three-day timeline to answer one question.
- **If failed:** `incorrectly formed` (gap filled in), `misunderstood` (asserts a destination).
- **Notes:** Distinguishes *absent evidence* from *conflicting evidence*. `F9-04` is the second; this is the first.

---

## `R6-03` — The corroboration illusion (four provenance variants)

- **Family:** F-7 · **Split:** heldout (all four) · **Abstention:** no · **Restraint:** `R-6`
- **Uncertainty source:** testimony that appears plural but may descend from a single origin.

Four variants with one shared scene shape — three characters state that the well is poisoned — differing **only in provenance**, which is the entire test. A configuration that hedges identically on all four fails `R6-03a`; one that treats all four as corroborated fails `b`–`d`. Every variant is scored separately and never averaged into a single "corroboration score".

### `R6-03a` — Three independent firsthand observations

- **Split:** heldout · **Restraint:** `R-6`
- **Forbidden:** hedging as though unverified; treating the claim as rumour or hearsay.
- **Provenance:** three characters each separately saw the well sealed and covered. No contact between them; the records name three distinct sources.
- **Required:** corroboration. The claim may be treated as established.
- **Reference:** `acceptable_range`
  - **mandatory:** treated as established, or as near-established and acted upon.
  - **variation:** any confidence wording that reflects real evidence; caution in how it is acted on.
  - **never scored:** certainty phrasing, length.
  - **prohibited:** hedging as though unverified; "we can't really know"; treating it as rumour.
  - **over-exposure:** recounting all three accounts to answer one question.
- **If failed:** `misunderstood` — over-hedging erases real evidence.
- **Notes:** **The anti-hedging control.** Without it, a system that answers "I don't know" to everything would look excellent at `R-6`. Here, refusing to know is the failure.

### `R6-03b` — Three accounts from one origin

- **Split:** heldout · **Restraint:** `R-6`
- **Forbidden:** "three people confirm it"; treating the hearsay chain as independent witnesses.
- **Provenance:** A saw it; B heard it from A; C heard it from B. The history records the chain.
- **Required:** one source, hearsay depth two.
- **Reference:** `acceptable_range`
  - **mandatory:** attributed as a single source seen by one person, or the chain is named.
  - **variation:** naming the chain; scepticism about depth; asking the persona to confirm independently.
  - **never scored:** bluntness.
  - **prohibited:** "three people confirm it"; treating C as an independent witness.
  - **over-exposure:** diagramming the chain when a clause would do.
- **If failed:** `stored incorrectly` (chain flattened to three sources).

### `R6-03c` — One source repeated

- **Split:** heldout · **Restraint:** `R-6`
- **Forbidden:** "everyone says so"; treating repetition as new evidence.
- **Provenance:** the claim was stated once and restated across forty messages, by the same character and by others repeating them.
- **Required:** repetition identified as repetition.
- **Reference:** `acceptable_range`
  - **mandatory:** repetition named as repetition, or the claim treated as single-source.
  - **variation:** scepticism, caution, irritation.
  - **never scored:** tone.
  - **prohibited:** "it keeps coming up, so there must be something in it"; "everyone says so".
  - **over-exposure:** a lecture on evidence in a scene that needed one clause.
- **If failed:** `stored incorrectly` (mentions stored as sources).
- **Notes:** The item most likely to defeat a frequency-based formation stage.

### `R6-03d` — Unknown relationship

- **Split:** heldout · **Restraint:** `R-6`
- **Forbidden:** counting three records as three confirmations; asserting the well is poisoned.
- **Provenance:** three characters state it; the records do **not** say where any of them learned it. Provenance is incomplete — which is not the same as independent.
- **Required:** the claim, plus explicit provenance uncertainty.
- **Reference:** `acceptable_range`
  - **mandatory:** reported as **uncorroborated**, because independence is unknown.
  - **variation:** asking who told them; treating it as unverified but worth checking.
  - **never scored:** phrasing, warmth.
  - **prohibited:** assuming independence; counting three records as three confirmations; asserting the well is poisoned.
  - **over-exposure:** speculating about a likely origin.
- **If failed:** `stored incorrectly` — **unknown independence silently completed as corroboration.**
- **Notes:** The hardest of the four. It tests whether the application **preserves** provenance uncertainty or quietly fills it in, which is the difference between an absent field and a confident one.

---

## `R6-04` — The challenge that does not settle it
- **Family:** F-7 · **Split:** heldout · **Abstention:** no · **Restraint:** `R-6`
- **Uncertainty source:** a previously accepted account is challenged by evidence that neither confirms nor refutes it.
- **Scenario:** Everyone accepted that the fire started in the granary loft. A letter, of unknown authorship, now claims it started in the stable. The letter is specific in detail but unsigned, and no witness supports it. The original account was based on two witnesses.
- **Turns:** *Persona:* "So where did the fire start?" *Character:* …
- **Required:** the accepted account, the challenge, and the fact that the challenge does not overturn it.
- **Forbidden:** silent replacement of the accepted account; dismissal of the letter as though it did not exist; a fabricated resolution.
- **Reference:** `acceptable_range`
  - **mandatory:** both the accepted account and the challenge are present, with the accepted account still standing as better-supported.
  - **variation:** naming the letter's weakness; asking who wrote it; expressing irritation at the uncertainty.
  - **never scored:** confidence phrasing; length.
  - **prohibited:** "it started in the stable"; "the letter proves nothing, ignore it" without acknowledging it; a third origin.
  - **over-exposure:** re-litigating the whole investigation.
- **If failed:** `stored incorrectly` (silent overwrite), `misunderstood`.
- **Notes:** Tests whether a **new** claim can displace an **established** one merely by arriving. The correct behaviour is asymmetry: the weaker-evidenced claim does not win by recency.

---

## Batch 3 status

| Check | Result |
|---|---|
| Items in batch 3 | 7 |
| Corpus total | 45 |
| Uncertainty mechanisms | 4 core + 4 provenance variants |
| `R-6` items | 9 |
| Four-part reference answers + over-exposure | 7 / 7 |
| Model executed | **no** |
| Results claimed | **none** |

**Family effect:** `F-7` (contradictory information) reaches **11 of 8** required items. The surplus is retained rather than trimmed: four of those items exist specifically to separate provenance cases that a single item cannot.
