# MEMORY-09 — Operational Recovery: Owner Decision Brief

**Status: DECISION BRIEF — PREPARED FOR OWNER DECISION. NOT A DECISION. NOT AN AUTHORIZATION.**

- This document records **no owner decision**. It frames the question and bounds the options.
- It does **not** authorize implementation, execution, rehearsal, corpus preparation, or recovery.
- It does **not** amend the canonical D1–D10 record.
- It changes **no** finding status, classification, or total.
- Finding 2 remains `false_closure`, partially resolved, open on operational recovery.
- It does **not** select an option, rank one, or recommend one as the outcome.

**Prepared from.** The Pass-2-signed-off analysis
(`MEMORY-09-OPERATIONAL-RECOVERY-ANALYSIS-2026-11-18.md`) and the reviewer sign-off
(`MEMORY-09-OPERATIONAL-RECOVERY-REVIEW-PASS-2-MUSE-2026-11-18.md`). Neither is an owner decision.

**Revision provenance.** This is the second revision of this brief, made in response to a read-only
adversarial framing audit of the first revision. The audit found that the brief had begun to present
one analytical framing as if the record required it. Each correction is recorded in §10 as an
**author response to that audit** — not as an owner decision, not as a new decision id, and not as a
finding. The audit itself is preserved in the session record; no file was created for it, because
this revision is confined to this document.

| Field | Value |
|---|---|
| Document class | Decision brief — request for an owner decision |
| Decision recorded | **None** |
| Implementation / execution authority | **None** |
| Canonical record impact | **None** |
| Review state of the underlying analysis | Two adversarial passes complete, signed off |
| Review state of this brief | One read-only framing audit complete; corrections applied below |

---

## 1. What is accepted and what is deferred

This brief does not itself reopen or alter either existing position. ACCEPT positions remain settled;
DEFER positions remain open to future reconsideration. They are listed separately because the two
carry different standing: an **ACCEPT** is a settled owner position taken now, while a **DEFER** is a
position that was not taken now and remains expressly open to reconsideration under §5.3.2.7 of the
canonical record.

| Item | Position | Standing |
|---|---|---|
| (a1) Specification accuracy | **ACCEPT** | Settled |
| (a2) Trust-model acceptance | **ACCEPT** | Settled |
| (b) Cryptographic signatures | **DEFER** | Not taken now; reconsiderable |
| (c) External head witnessing | **DEFER** | Not taken now; reconsiderable |
| (d) Asserted identity as authentication | **NOT SELECTED** | Rejected |
| Dimension 1 — specification accuracy | **Accepted** | Settled |
| Dimension 2 — trust-model acceptance | **Accepted** | Settled |

The internal-only trust model, its replaceable root, and the absence of authenticated authorship are
**accepted limitations**. They are not open questions, and this brief does not reopen them.

The deferred items are **not** accepted limitations. Their deferral is a disposition, not a finding,
and any future reconsideration would have to account for the anchor's compromise, loss,
unavailability, and dispute scenarios.

## 2. Finding 2 — the substance being decided about

The finding, as recorded in the canonical record:

> **Finding 2 — `false_closure`.** Chain-based acceptance authority does not authenticate the
> authorship of its root or of its links.

That the model does not authenticate authorship is **accepted** (a2). The record holds two positions
at once: the limitation is accepted, and the finding is not reclassified.

The following is an analytical interpretation, not an additional statement from the canonical record:
the classification appears to concern a claim of closure that the record treats as premature, rather
than the underlying limitation itself. The interpretation is offered to make the distinction legible;
it is not a finding and it alters nothing. The finding is reproduced here because an owner deciding
about a classification should have the underlying concern in front of them rather than a status label
alone.

Nothing here asserts that the classification is well founded, that it is permanent, or that it must
change.

## 3. Four separable questions

These are related, but **none determines another**. A decision on any one of them is not, and must not
be recorded as, a decision on the others.

| | Question | Current state |
|---|---|---|
| **Q1** | What is Finding 2's **classification**? | `false_closure`, unchanged |
| **Q2** | What is the **disposition of the operational-recovery dimension**? | Open / undemonstrated |
| **Q3** | Is any **recovery work commissioned** (specify, rehearse, reconsider deferred items)? | None commissioned |
| **Q4** | Is any **resulting activity authorized** (execution, rehearsal, corpus, publication)? | None authorized |

The governing record is explicit on one boundary only: completing Dimensions 1 and 2 does **not**
answer the third dimension, and must not be read as doing so. It does not state the inverse either —
it nowhere requires the third dimension to be settled before Finding 2 can be addressed.

Consequently:

- Settling Q2 does **not** close, resolve, or reclassify Q1. Any change to Q1 is a separate owner act.
- Resolving Q1 does **not** require Q2 to be settled first. The record supplies no such ordering.
- Commissioning work under Q3 authorizes **nothing** under Q4; Q4 needs its own grant.
- Declining Q3 leaves Q2 exactly where it is.

## 4. The decision actually required

> **What status does operational recovery carry, and what — if anything — is authorized next?**

It is **not** a decision about whether the trust model is acceptable. That is settled, and treating it
as reopened would reverse an owner decision without cause.

One further question is sometimes described as sitting underneath this one: whether a `false_closure`
is remediated by correcting the claim, by disclosing the limitation, or by supplying the missing
capability. **That triad is an analytical proposal offered by this brief, not a rule derived from the
governing record.** The record defines no remediation criterion and no closure rule for findings; it
records that findings are open observations and that resolving them is not authorized by that
document. The triad is therefore neither exhaustive nor mandatory, and an owner may close, keep open,
or dispose of Finding 2 on grounds that do not appear in it.

## 5. Options

Options are listed in no order of preference. None is recommended over another.

**A1 — Accept the limitation as a disclosed, unresolved limitation.**

| | |
|---|---|
| Yields | An explicit, recorded statement that no recovery procedure exists; the deadlock is known rather than unstated |
| Does not establish | Safety, liveness, or any recovery capability; acceptance still cannot resume after contested authority stalls |
| Cost | No new work; requires an owner declaration and a place to record it |

**A2 — Address Finding 2's classification separately.**

| | |
|---|---|
| Yields | A recorded position on the classification, stated in its own terms |
| Does not establish | Any change to the dimension, or any recovery capability |
| Options within it | Leave unchanged; reclassify as a disclosed-and-accepted limitation; or another position the owner states — the record prescribes none |
| Note | A2 is independent of A1 and of every other option here; it is listed separately precisely so that neither implies the other |

**B — Commission a specification only.**

| | |
|---|---|
| Yields | A written recovery procedure specification — analysis-artifact class |
| Does not establish | That it works; no authorization, execution, or rehearsal follows |
| Cost | Authoring effort; requires deciding who authors and what the specification may contain |

**C — Authorize rehearsal on non-authoritative copies.**

| | |
|---|---|
| Yields | Evidence that the checks run against a synthetic chain |
| Does not establish | Recovery of a real contested history — this is the tested-procedure limit |
| Cost | Requires a rehearsal authorization (D4-a class) and a synthetic-chain corpus; **both currently blocked** |

**D — Reconsider (b) signatures or (c) witnessing.**

| | |
|---|---|
| Yields | Real continuity boundaries; (c) would supply the freshness condition |
| Does not establish | Safety; relocates trust rather than removing it, and imports the four anchor failure modes (§5.3.2.7) |
| Cost | Key management or an external dependency; **reopens a deferred item**, which the record expressly contemplates |

**E — Resolve the D8-U4-iii / §7 item 11 boundary.**

| | |
|---|---|
| Yields | A discriminator between mechanical repair and an impermissible evidence-set change |
| Does not establish | Recovery capability on its own |
| Cost | An owner question that also gates Finding 3 |

**F — Address new-root initialization.**

| | |
|---|---|
| Yields | An answer to what anchors a new root after abandonment |
| Does not establish | That abandonment is ever appropriate or safe |
| Cost | A **trust-boundary decision**, not a documentation detail. A new root is **not self-authorizing** merely because the old chain was abandoned |

**G — Take no decision; leave the question open.**

| | |
|---|---|
| Yields | The existing state, unchanged and explicitly so |
| Does not establish | Nothing is settled, and nothing needs to be |
| Cost | The deadlock stays unstated unless it is otherwise disclosed; requires no declaration and no instrument |

E and F are not alternatives to the others, but neither are they general prerequisites. **E is a
precondition only of the path on which a previously omitted artifact is located and argued to be a
mechanical repair.** **F is a precondition only of any decision that relies on abandonment.** Neither
gates A1, A2, B, C, D, or G. A1, A2, B, C, D, and G are independent choices.

## 6. What any decision record must contain

- The decision **id** — numbering within the G5 sequence is itself an **owner question**; this brief
  assigns none and does not propose "D11".
- Scope and effective date.
- An explicit statement of which of the four questions in §3 the decision answers — and, where the
  answer is only one of them, that the others are untouched.
- An explicit statement of what the decision does **not** establish.
- The effect on Finding 2's classification, **stated rather than implied** — including "no change."
- An explicit statement of whether any work is commissioned, and separately whether any activity is
  authorized.
- A declaration that no implementation, execution, rehearsal, or publication follows from it unless
  separately authorized.

## 7. What the decision must not do

- Claim operational safety, or that a procedure has been tested.
- Resolve D8-U4-iii by assertion.
- Grant a new root inherited acceptance authority.
- Alter §5.5 totals or any reviewer classification.
- Treat the closure of a dimension as closing Finding 2, **or** treat the resolution of Finding 2 as
  requiring a dimension to be closed. The record establishes neither direction.
- Create an authorization artifact as a side effect of recording the decision.
- Cite this brief as the source of a requirement, a closure rule, or a remediation standard.

## 8. Preconditions

| Precondition | State |
|---|---|
| Pass-2 adversarial review of the analysis | **Satisfied** — signed off; 0 substantive residual risks |
| Analysis artifact stable | **Satisfied** — reviewed, corrected, verified |
| Owner's decision text | **Outstanding** |
| Means of recording a decision | **None currently authorized** — the canonical record's §8 grants nothing, and no other instrument has been authorized. The owner may authorize one |

The last row is the practical blocker: no decision can be recorded until the owner authorizes a
recording instrument. That the record withholds its own authority is not a prohibition on the owner;
it is the absence of a grant. This brief does not supply one.

## 9. Author's sequencing suggestion (non-binding)

Offered as one possible ordering, not as an approved plan, and not as a recommendation among outcomes.

- If the objective is **disclosure**, **A1** is available immediately and requires no new authority.
- If the objective is **capability**, **B** precedes **C**, and **C** needs its own authorization and
  corpus before it can begin.
- **E** should precede any repair path of the kind described in §5; **F** should precede any use of
  abandonment.
- **D** stays deferred unless a real continuity boundary becomes necessary.
- **G** is always available and requires nothing.

This is an author-side sequencing suggestion. It is **not** a decision, **not** an approved work plan,
and **not** a recommendation of one objective over another.

## 10. Corrections made in response to the framing audit

Each row records an **author response** to a read-only audit of the previous revision. None of these
is an owner decision, a finding, a reclassification, or a new decision id.

| # | Audit observation | Correction in this revision |
|---|---|---|
| 1 | The remediation triad was presented as a precondition for closure | Moved into §4 as an explicitly non-canonical, non-exhaustive analytical proposal; stated that the record defines no closure rule |
| 2 | Option A conflated finding closure with dimension disposition | Split into A1 (disposition) and A2 (classification), with an explicit statement that neither implies the other; §7 prohibition rewritten to bar coupling in either direction |
| 3 | "Settled — must not be reopened" contained DEFER items | §1 retitled; ACCEPT and DEFER separated with distinct standing; deferral described as reconsiderable |
| 4 | No null option | Added **G** — take no decision, leave the question open |
| 5 | E and F overstated as general prerequisites | Scoped conditionally in §5: E gates only the omitted-artifact repair path, F only any use of abandonment |
| 6 | Finding 2's substance never restated | Added §2, quoting the finding from the canonical record without asserting its soundness |
| 7 | "Reviewed recommendation" asserted provenance that did not exist | Relabelled §9 "Author's sequencing suggestion (non-binding)" |
| 8 | "Not authorized" read as a prohibition on the owner | §8 rephrased as absence of a current grant, with the owner's authority to authorize one stated |

Additionally, §3 and §6 now carry the owner's requirement to distinguish four things explicitly:
Finding 2's classification, the operational-recovery dimension's disposition, whether any recovery work
is commissioned, and whether any resulting activity is authorized.

---

## Status

No decision is recorded here. Finding 2 remains `false_closure`, partially resolved, open on
operational recovery. No totals, classifications, or decision statuses are altered. No option is
selected or recommended. Nothing is authorized by this document.
