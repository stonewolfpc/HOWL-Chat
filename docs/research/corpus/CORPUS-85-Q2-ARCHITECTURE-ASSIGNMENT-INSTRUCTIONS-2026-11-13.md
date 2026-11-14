# Q2 Assignment Instructions — Architecture Adjudication (`CORPUS-82`, unexecuted)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-85` · **Version:** 1 · **Date:** 2026-11-13.
**Type:** assignment instruction. **It commissions work; it performs none.** No finding, ranking, or answer appears in this file.
**Issued by:** `CORPUS-89` (owner approval and issuance instrument, 2026-11-12), superseding the proposed issuance recorded in `CORPUS-84` §1. **Charter:** `CORPUS-80` §6, unamended.
**Prepared by:** DeepSeek, as secretary. **DeepSeek is not eligible to execute this assignment** (`CORPUS-80` §8.1: author of `CORPUS-78`).
**Execution status:** **NOT EXECUTED.** The deliverable identifier `CORPUS-82` is reserved and unassigned to any file.

---

## 1. Eligibility — read this before accepting

**Barred:** the author of `CORPUS-78` (DeepSeek/Buffy) — it proposed the architecture under test. The author of `CORPUS-79` (Muse) — reserved by the owner for a later independent challenge, and its `A′` is a candidate under test. The author of `CORPUS-81` (Glyph Cluster) — already the Q1 adjudicator, and the owner has ruled that it must not adjudicate its own work again.

**Provenance of the three bars.** The `CORPUS-78` and `CORPUS-79` exclusions are required by `CORPUS-80` §8.1 as written. The `CORPUS-81`-author exclusion is an **additional assignment restriction included in these prepared instructions** — it rests on the owner's instruction that Glyph Cluster not be asked to adjudicate its own work again, and its extension from that Q1 context to Q2/Q3 is the **secretary's**, not an owner decision (`CORPUS-88` §2). It is advice carried into the instructions, not a rule inherited from `CORPUS-80`. It follows that the **charter's mandatory exclusions may already leave no eligible adjudicator**, and this third restriction narrows the field further as advice the owner may adopt or drop.

**Consequence.** If those three are the only identities available in this workspace, **no eligible adjudicator exists**, and Q2 cannot be performed as an independent adjudication. In that case the correct action is:

1. **Do not** run Q2 under a borrowed or reassigned label.
2. **Do not** treat the near-miss as independent. Run it only as a **declared non-independent documentary exercise**, stated in the deliverable's first section, in the same words `CORPUS-81` used.
3. Record that the result **cannot count toward convergence**, and that `G1`'s convergence requirement remains unmet.
4. Escalate to the owner for a scheduling decision — not for a waiver.

**A non-independent Q2 is still permitted and still useful.** It is simply not evidence of independence, and it cannot be cited as such.

---

## 2. The six determinations required

Answer each **separately**, with its own evidence and its own confidence. A single overall verdict ("the architecture is sound / unsound") is **not** an acceptable deliverable.

| # | Determination | The question, stated precisely |
|---|---|---|
| **Q2.1** | Cue independence | Can an instrument exist whose acceptance decision and mode result are **genuinely** independent — such that `pass / mode-not-demonstrated` is reachable — given that a cue must be perceptible? Test this against `F6-04`, where the cue *is* the `mandatory` text (`CORPUS-80` §3 P4). |
| **Q2.2** | The two failure shapes | Distinguish **(F-i)** the item cannot produce the mode at all from **(F-ii)** the item can produce it and the reply does not demonstrate it. State whether the corpus can currently express (F-ii) for a gated field, and what would have to be true for it to do so. |
| **Q2.3** | The not-assessable conditions | Rule on the four conditions `CORPUS-78` §3.1 proposed. Two of them — (ii) response-dependent, (iii) mis-scoped — were found to be decidable only by reading the reply, contradicting the claim that all four are pre-reply (`CORPUS-79` F-4). Either repair the condition set or state that no pre-reply condition set exists. |
| **Q2.4** | Gating ruling | State whether a **non-gating** dual-field instrument is achievable at all, and if so, name the item shapes where it is. If hidden gating is unavoidable for perceptible-field modes, say that — it is a legitimate outcome, and better than a claim the record cannot support. |
| **Q2.5** | Pass-set recount | Recount from **item text**, not from declarations, how many of the 45 units have exactly one valid mode assignment. **Neither 8 (as carried in `CORPUS-78`/`CORPUS-79`) nor 3 (as urged by `CORPUS-79` F-17) may be assumed.** Publish the per-unit tally so the count can be checked. Expected divergence from both figures is a finding, not an error. |
| **Q2.6** | Application-scored modes | Determine whether `R-1`, `R-4`, `R-5` — which record what the model *did* rather than what it *said* — can be scored from a reply at all. Do **not** design the application instrument; that is Q3 (`CORPUS-86`). State only what the reply can and cannot establish. |

---

## 3. Required sources

Read from disk; cite file, section, and line for every substantive claim.

**Mandatory:** `MEMORY-09` §§3.5, 5, 6, 7 (esp. L119, L230) · `CORPUS-01` §2.4 (L78–83), L96, L121–126, L128–136, L144 · `CORPUS-07` §§2, 8, 10 · `CORPUS-78` §§3.1, 6.2, 8, 9 · `CORPUS-79` §§4–7, 9 · `CORPUS-80` §§3–4, 6, 8 · `CORPUS-05` §7 · `check_corpus.py` (L58 and C-1/C-4).

**Q1 is required as an input and is not a premise.** `CORPUS-81` is on file (`CORPUS-84` §4) and **its findings remain unadopted**. When Q2 depends on mode membership, it must state the result **under both readings** where they differ — the specification table's membership and the item-level tags' membership — and must not silently adopt either. A Q2 that presupposes Q1's conclusions is out of charter.

---

## 4. Evidence standard and label discipline

Every substantive sentence carries one label: **governing text**, **item adjudication**, **reviewer interpretation**, **proposal**, or **unresolved**. Specifically:

- A reviewer's assertion is **not** evidence that wording exists. Quote the record.
- Model identity is **never** inferred from a persona label, a filename, or a self-description.
- A count is a finding only with the per-unit list attached.
- An unreachable cell (`pass / mode-not-demonstrated`) is evidence about the **instrument**, not about the corpus, until Q2.4 is ruled on.
- Where the record is silent, write **unresolved**. Do not fill the gap and proceed.

---

## 5. Prohibitions

Do **not**: amend any item, tag, split, threshold, or score · edit `MEMORY-09`, `CORPUS-01`, or `CORPUS-07` · implement a checker, validator, or runtime change · run any model evaluation or experiment · infer model identity from output · begin `G2` · close `G1` · treat `A′` or Option A as the expected answer · cite convergence · seal or unfreeze anything.

**Do not repair the result by raising the standard after seeing it.** If the recount lands somewhere inconvenient, publish it.

---

## 6. Deliverable specification (`CORPUS-82`)

- Filename: `CORPUS-82-ARCHITECTURE-ADJUDICATION-<topic>-<date>.md`.
- §0 identity and **independence disclosure stated explicitly**, including the eligibility analysis from §1 above and whether the author is eligible at all.
- §1 scope, sources, and method.
- §2–§7: one section per determination Q2.1–Q2.5, each with question, evidence, ruling, confidence, and what would falsify the ruling.
- §8 the integrated picture — explicitly **not** a single verdict, and explicitly naming which determinations remain independent of each other.
- §9 what this adjudication cannot establish.
- §10 status: nothing adopted, nothing applied, convergence unchanged.
- Registered in `docs/DOCUMENT-REGISTRY.md` with the established row form.

---

## 7. Acceptance criteria — checkable

1. All six determinations answered, none omitted, none merged.
2. Eligibility stated; if ineligible, the label says so.
3. Recount includes a **per-unit tally** and neither 8 nor 3 is asserted without it.
4. Every source claim line-cited; zero uncited assertions about what the record says.
5. Membership-dependent findings give **both** readings (`CORPUS-81` unadopted).
6. No item, tag, threshold, or score changed; no inference performed.
7. Prohibitions honoured and self-certified.
8. Filed and registered; conversation output alone does not count as completion.

---

## 8. What acceptance does not mean

A conforming `CORPUS-82` **does not** validate the architecture, **does not** adopt `A′` or Option A, **does not** authorize a corpus repair, and **does not** advance convergence. It answers Q2. Adoption is the owner's decision, taken afterwards, on all three questions.
