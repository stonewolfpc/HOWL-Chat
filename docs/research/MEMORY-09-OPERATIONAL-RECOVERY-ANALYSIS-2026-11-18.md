# MEMORY-09 — Operational Recovery Analysis (Reviewed Artifact — Two Adversarial Passes)

**Status: ANALYSIS — ADVERSARIAL REVIEW COMPLETE (PASS 2 SIGNED OFF); NOT OWNER-APPROVED; NOT A CANONICAL AMENDMENT.**

- This document is **not** a canonical D1–D10 decision amendment.
- It is **not** an owner-approved operational-recovery procedure.
- It does **not** authorize implementation, execution, rehearsal, or recovery.
- **Both adversarial passes complete.** Pass 1 corrections incorporated; Pass 2 verified all seven (0 contradictions, 0 boundary violations, 0 substantive false-closure risks).
- **Review sign-off is not acceptance.** No owner approval is conferred by either pass.
- Finding 2 remains `false_closure`, **partially resolved, open on operational recovery**.
- No finding totals, classifications, or decision statuses are changed by this document.

**Source.** The owner's corrected analysis supplied in the 2026-11-18 conversation, recorded
here as the review artifact.

| Field | Value |
|---|---|
| Document class | Analysis artifact — both adversarial passes complete; **not canonical, not owner-approved** |
| Relationship to canonical record | Separate artifact; **does not amend** `MEMORY-09-G5-SEQUENCING-D1-D10-DECISION-RECORD-2026-11-18.md` |
| Review state | **Complete** — Pass 1 (7 corrections, incorporated) and Pass 2 (signed off, 3 optional editorial nits). **Not owner-approved** |
| Approval state | Not owner-approved; not a closed specification |
| Implementation authority | **None** |
| Verification state | No execution, no rehearsal, no rehearsal harness, no source re-verification |
| Git impact | New untracked documentation file only; nothing staged, committed, or pushed |

---

# Operational recovery — corrected analysis

## Q1. Triggers

Observable events only, and bounded: a quiescent unresolved chain blocks nothing, so it triggers
nothing. Only a blocked pending acceptance makes a condition actionable.

| Observation | Character | Taxonomy |
|---|---|---|
| Referenced artifact **absent**, reference unresolvable | gap | `incomplete` — *proposed application* |
| Referenced artifact **present, digest ≠ recorded** | conflict | contradiction family — *proposed; subtype coverage not yet justified* |
| Competing successors, co-presented | conflict | existing |
| Multiple heads, no valid path | conflict | existing |
| Cycle at the head | structural conflict | existing |
| **Presentation asymmetry** — two parties report different heads | **not detectable by one validator** | social trigger |

Different operations (resolution failure vs recomputation) and different preserved objects. Repair
standing, if any, remains per §7 item 11 — **unresolved**. Presentation asymmetry remains
non-self-triggering; the trigger is necessarily social.

## Q2. Preservation

The collection sits **outside** the contested presentation — adding a manifest to the contested
directory would mutate the state it preserves and could be read as a chain artifact or alter the
apparent head set.

Records: presented directory or snapshot identity with method and time; every artifact with digest;
the observed head set; which validator observed which presentation; refusals and trigger
observations; the time and source of each; and the terminal record, if preserved, **as a claim only
— never as an adjudicator** (D8-U5).

**All observed branches, not "both"** — asymmetry means the known count cannot be assumed complete.

Preservation and repair are strictly separated:

> Locating a previously omitted artifact may alter the committed evidence set under the same
> `run_id`. The existing record supplies **no testable discriminator** between mechanical recovery
> and an impermissible evidence-set change. Accordingly this path is unresolved and cannot be
> treated as an established recovery operation.

This preserves §7 item 11 exactly: an owner question, not a new finding, not a resolved rule. The
collection itself inherits the accepted internal-only limits; its integrity is asserted, not
derived.

## Q3. Continuity

The anchor is an **exclusion filter**. It excludes non-continuations; it does not select an
authorized continuation.

| Condition | Mechanical result |
|---|---|
| Candidate makes **no claim** of continuity with H | **Excluded** by filter; **no continuity defect asserted or implied** for that candidate |
| Candidate **claims** continuity but cannot demonstrate it | **`incomplete`**; refusal and escalation |
| H **absent from all** presented candidates | **No global conclusion**; continued refusal / escalation where continuity is required |

The third row is not a fourth kind of evidence defect. It is **the limit of what the presentation
establishes** — the instrument cannot certify that the candidate set is complete, so an empty
filter result is evidence about the presentation, never proof that no continuation exists.

**Freshness.** With attested heads H₁ < … < Hₙ, testing partitions on Hₙ and leaves (Hₙ, head]
unanchored; divergence inside that interval is invisible. A fully attested sequence pins history
only **through its last attested point**. Closing it needs an attested current head or an
explicitly closed sequence. Witnessing value depends qualitatively on both coverage(what) and
freshness(when) — **illustrative, not a defined metric** — and both inherit (a2)'s limit, since an
anchor can be compromised, lost, unavailable, or disputed (§5.3.2.7).

Where |C′| = 1, the survivor is the sole *presented* continuation of H **within this presentation** —
uncontested, **not authorized**. Owner acceptance is still required.

## Q4. Abandonment

Abandonment has no mechanical trigger. It belongs exclusively to the authority layer:

- The owner may **declare** the old chain abandoned.
- The declaration records **who, when, scope, reason**.
- It does **not** claim the instrument proved no continuation existed.
- The new root starts a **new chain** and inherits **no** old acceptance authority.
- The declaration is preserved **as** a declaration under the existing internal-only trust model;
  its authorship is **asserted, not independently authenticated**, per Finding 2 / D4-a.

Citation boundary, three-part: prior records remain citeable as surviving documents or historical
claims; they establish no authenticated ordering or continuity through the abandoned chain; the new
root does not inherit their authority. Content survives; ordering and continuity do not.

**Abandonment is not D8-U4-iii succession.** A successor under U4-iii continues under the original
run's authority; abandonment ends that chain and begins another carrying none of it.

## Q5. Resumption

Three separately recorded facts, distinct provenance, **no synthesised `recovery_succeeded`
verdict**:

| Field | Provenance | Establishes |
|---|---|---|
| Traversal result, scoped to the presentation examined | mechanical | no detectable conflict within this presentation |
| Owner designation | authority declaration | the history accepted for continued use |
| Non-existence of other branches | **unestablished** | not establishable by this instrument |

The designation names the **specific head by digest** — not a directory or branch label — with
scope and time. Labels are mutable; only digests are identities.

## Three layers

1. **Mechanical** — observes, filters, classifies, refuses, preserves. Does not adjudicate.
2. **Authority** — the external owner declares successor, abandonment, or resumption, with
   explicit scope and provenance.
3. **Evidence** — preserves observations and declarations without converting either into a
   stronger claim than its provenance supports. A declaration is recorded **as** a declaration.

**Qualification:** the authority layer can designate a history for continued use, but designation
does not demonstrate that the recovery procedure is operationally **safe**.

That splits cleanly inside §5.3.2.4's third dimension:

| Sub-question of "operational recovery" | Status |
|---|---|
| Is there a valid external authority to designate? | **Supported** — the owner is the external resolver |
| Is the procedure operationally safe? | **Open / undemonstrated** |

So the open item is exactly operational safety, not a re-litigation of authority — and it stays
open regardless of how cleanly a designation is recorded.

## Unchanged

Tested-procedure limit: the instrument cannot validate the procedure that repairs it; a rehearsal on
a non-authoritative copy exercises the checks without establishing recovery of a real contested
history.

Open: D8-U2 subtype coverage (proposed, unjustified); D8-U4-iii / §7 item 11 (no testable
discriminator); (b) and (c) deferred, freshness belonging to (c); artifact governance across
specification / authorization / evidence. No procedure exists; nothing tested or rehearsed.

## Additional unresolved questions (first adversarial review)

Reviewer questions, **not findings**, and **not counted in any total**:

1. **Quiescence provenance.** "A quiescent unresolved chain blocks nothing, so it triggers nothing"
   is consistent with §5.3.2.3 liveness, but liveness states the absence of a *mechanism*, not the
   absence of a *duty*. Whether quiescence-triggerlessness restates liveness or introduces a new
   liveness rule is unresolved. It is a reviewer restatement, **not an owner decision**.
2. **New-root initialization.** If a new root inherits no prior acceptance authority, what
   establishes its initial trust anchor? D6 requires prior links for chain-based authority and a new
   root has none. Out of scope here, listed so that "new root" is not read as self-authorizing.

Both belong alongside (b) and (c) as deferred, not resolved.

---

Finding 2 remains `false_closure`, partially resolved, open on operational recovery. No totals,
classifications, or decision statuses changed. This is a review artifact, not a §5.3 subsection.

---

## Review history

| Pass | Status | Outcome |
|---|---|---|
| 1 | Returned | 0 contradictions, 0 boundary violations; 5 false-closure risks; 2 additional reviewer questions; 7 wording corrections — **all incorporated** |
| 2 | **Returned — signed off** | 0 contradictions, 0 boundary violations, **0 substantive false-closure risks**; all 7 corrections verified; 3 optional editorial nits; unresolved unchanged |

### Corrections incorporated from pass 1

| # | Location | Correction |
|---|---|---|
| 1 | Q3, singleton case | "sole *presented* continuation … **within this presentation**" — closes the completeness inference |
| 2 | Q3, exclusion row | "no continuity defect asserted or implied for that candidate" — scoped to continuity with H |
| 3 | Q3, freshness | coverage/freshness dependency marked **qualitative, not a defined metric** |
| 4 | Q4, bullets | Declaration authorship **asserted, not independently authenticated** (Finding 2 / D4-a) |
| 5 | Q1, closing paragraph | "Repair standing, if any, remains per §7 item 11 — **unresolved**" |
| 6 | Q4, end | "**Abandonment is not D8-U4-iii succession**" — new root versus successor under original authority |
| 7 | Q2, records list | Terminal record, if preserved, **as a claim only — never as an adjudicator** (D8-U5) |

### Pass-2 outcome

Reviewer: Muse (AI reviewer). Review artifact:
`docs/research/MEMORY-09-OPERATIONAL-RECOVERY-REVIEW-PASS-2-MUSE-2026-11-18.md`.

| Result | Count |
|---|---|
| Confirmed contradictions | **0** |
| Boundary violations | **0** |
| Substantive false-closure risks remaining | **0** |
| Pass-1 corrections verified | **7 of 7** |
| Unresolved questions | unchanged; not findings, not counted |

Reviewer's verdict: all seven corrections "incorporated faithfully and correctly placed"; no correction
introduced a new contradiction or boundary violation; the open item "remains exactly operational
safety, per §5.3.2.4."

**Optional editorial nits** — reviewer states no re-review required; **recorded, not applied**:

| # | Location | Nit |
|---|---|---|
| N1 | Q3 singleton | "uncontested" inherits "within this presentation" only implicitly |
| N2 | Q1 quiescence sentence | bald in Q1; cross-reference to *Additional unresolved questions* (1) suggested |
| N3 | Header / history labelling | "Unreviewed" reads as verbally tense alongside "review received"; Document class could be tightened |

**Sign-off is not acceptance.** The review certifies wording fidelity and boundary compliance of
this artifact. It confers no owner approval, no implementation or recovery authority, and no change
to Finding 2.

## Review status — both passes complete

Both adversarial passes reviewed **this saved artifact**, not a summary or an earlier conversational
draft. The scrutiny applied was:

- Q1 trigger classification and presentation asymmetry.
- Q2 preservation versus authentication, and the unresolved D8-U4-iii boundary.
- Q3 exclusion, incompleteness, continuity, and the absence-of-H limitation.
- Q4 owner-declared abandonment and the non-inheritance of prior acceptance authority.
- Q5 separation of traversal, owner designation, and global branch existence.
- The three-layer model, and the distinction between accepted external authority and
  undemonstrated operational safety.
- Compatibility with D6, D8-U2, D8-U4-iii, D8-U5, §5.3.2.4, and §7 item 11.

The review distinguished **confirmed contradictions**, **false-closure risks**, **boundary
violations**, **unresolved questions**, and **compatible conclusions**. Pass 2 applied all of these
and found no substantive residual risk.

**No further pass is requested.** The three optional editorial nits remain unapplied and are recorded
under *Pass-2 outcome*; applying them would be a new edit, not a review requirement.
