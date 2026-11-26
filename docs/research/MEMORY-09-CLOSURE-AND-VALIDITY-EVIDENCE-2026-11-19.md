# MEMORY-09 — Closure and Validity: Targeted Evidence

## Whether a "no further entries" declaration is permitted, and whether recording can be separated from conferring

**Status:** evidence record. **No mechanism, no schema, no selection.** OD-b remains open. Read-only: no canonical
edit, no test, no probe, no ledger entry, no implementation authorisation.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Q-1 | Does the existing architecture permit an **independently meaningful declaration that no further lifecycle entries are expected**? (OBE-36) |
| Q-2 | Can the architecture **demonstrably enforce** a distinction between recording an act and that act conferring status? (OBE-18) |
| Method | targeted textual evidence from r33 only; no design |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. Why these two questions were worth separating

**CEV-01 (Inferred).** They looked like the same kind of gap — "the option needs something the canon lacks" — but
they are not. Q-1 asks whether a **permitted** representation exists. Q-2 asks whether a **required** distinction
is enforceable. The evidence below shows each already has a governing rule in r33, and that the rules point in
**opposite directions**: closure is *required in effect but absent in text*; validity is *stated but
unenforceable*.

---

# Part A — Closure

## 2. The architecture already reads meaning from a log's end

**CEV-02 (Source-derived — K-R17, `:137`).**

> "`revoked` is valid only when the pair's ordered log **currently ends** in a learn event. A revoke with no open
> episode is **rejected at write**, not silently ignored, because a silent no-op hides formation bugs."

**CEV-03 (Calculated — three properties that matter here).**

1. The **terminal position of a log is load-bearing** in this architecture; it is not merely storage order.
2. An impossible terminal transition is **rejected at write** rather than tolerated — the stream is kept free of
   states that cannot be interpreted.
3. The stated reason is precisely the failure mode OD-b is weighing: *"a silent no-op hides formation bugs."*

## 3. The existing solution to exactly this ambiguity, in the other stream

**CEV-04 (Source-derived — K-R19, `:147`).**

> "A revoked pair has **no edge**, not an edge with a 'revoked' flag. 'Forgot' is the absence of the current
> relation; the fact that it was once known lives in the log, not in a half-state row."

**CEV-05 (Inferred — the pattern, stated plainly).** The acquisition stream faces the same shape of question —
"does an absent successor mean something?" — and answers it by **combining two rules**: K-R17 rejects the
impossible transition at write, and K-R19 fixes what absence means so that absence has exactly **one** meaning.
Neither rule alone would suffice: rejection without a stated meaning leaves absence ambiguous; a stated meaning
without rejection leaves the stream interpretable in two ways.

**CEV-06 (Inferred — the precedent this establishes).** The architecture **does** accept the idea that a stream's
shape can carry a second-order fact (*that the history ends here in this way*), and prefers **make the wrong state
unrepresentable** over inferring intent from silence. What it has not done is apply that pattern to the memory
lifecycle stream.

## 4. The governing invariants make L2's and L3's silent failure an I15/I18 breach

**CEV-07 (Source-derived — I15, `[X]`).**

> "**Absence of information is never treated as information about absence.** An unpopulated field, an empty list, a
> missing record and an unwritten schema each mean 'not established', never 'established as none'. …
> **Applies uniformly to:** missing acquisition history, empty knowledge sets, unspecified permissions, undeclared
> action requirements, missing asset versions, and incomplete retrieval evidence."

**CEV-08 (Inferred — the direct consequence).** A fold that reads "no promoting event follows" as "this record is
terminally `provisional`" is treating absence of information (no further entries present) as information about
absence (no further entries occurred). **I15 forbids that inference.** L2's and L3's silent-failure mode is
therefore not merely a design risk: as specified without closure, it **violates an existing `[X]` invariant**.

**CEV-09 (Source-derived — I18, `[X]`).**

> "**The system may not use absence of evidence as evidence that nothing happened.** Every subsystem must name the
> **independent evidence** that would reveal its own silent omission — evidence it does not produce itself."

**CEV-10 (Calculated — I18 is a hard requirement on a lifecycle representation).** A lifecycle stream whose lost
tail produces a plausible terminal state has exactly the property I18 forbids: its own failure is invisible, and it
cannot name evidence of that failure that it does not itself produce. So under I18, a fold-based or act-based
representation must be able to point to **something outside the stream** that would reveal truncation.

**CEV-11 (Inferred — the important asymmetry this creates).** CEV-08 and CEV-10 apply to **inference from
absence**, not to **assertion**. This is the precise reason L1 is exempt from closure and L2/L3 are not:

| | Does it infer a fact from absence? | I15 posture |
|---|---|---|
| **L1** | no — reports the last asserted state | **compliant**: "the state is `provisional`; whether it advances is not established" |
| **L2** | yes — state derived from the absence of a successor | **breach unless closed** |
| **L3** | yes — state derived from the absence of an exit act | **breach unless closed** |

**CEV-12 (Inferred — a correction to OBE-08's basis).** In the option evaluation I credited L1's terminal-`provisional`
handling to "the state is stored." That is the mechanism, not the reason. The reason is **CEV-11**: L1 never
derives anything from a missing entry, so it never needs to know whether the history is complete. This strengthens
the finding and narrows it: L1's advantage is **compliance with I15 by construction**, not merely convenience.

**CEV-13 (Inferred — and one price L1 still pays).** Because L1 does not infer, it also cannot say *"this record will
never advance"* as a fact. The rumour's "stays provisional" is then a **policy statement about the kind**, not a
state the stream attests. That is honest under I15, but it means no option currently expresses terminality as
established fact — a point in favour of keeping the question open rather than resolving it inside OD-b.

## 5. Precedents for a completeness declaration, and what they are worth

**CEV-14 (Source-derived — §10.2, and its status line).** The evaluation-log section distinguishes outcomes as
`reproducible`, `auditable_only` or `expired` — **"never collapsed into 'replay worked' or 'replay broke'"** — and
states that "the `expired` state lets the system be honest instead of approximating." The same subsection carries
**"Status: proposed, not adopted."**

**CEV-15 (Calculated — what that precedent establishes).** r33 already contemplates a **third axis** describing
whether a record can be trusted to be *complete and faithful*, kept separate from the record's own subject. That is
structurally the shape Q-1 asks for. But it is **proposed, not adopted**, and it lives outside the lifecycle
stream, so it is a precedent for the **idea**, not a construct available for use.

**CEV-16 (Source-derived — D1 and §11.1 `[P]`).** Chat deletion removes a chat's experiential memories; §11.1 is
tagged "proposed, not adopted". Deletion is therefore a **recorded, authorised reason for a history to stop** —
not a completion, but an independently meaningful end.

**CEV-17 (Inferred — the nearest thing to an answer in the existing text).** The architecture already distinguishes
"the history ended because something recorded an end" (deletion, `via=revoked`) from "the history merely stops."
It has **no construct** for the second case in the lifecycle stream, and no construct declaring *completion* in
either stream.

## 6. Q-1 verdict

**CEV-18 (Calculated — four-part verdict, no design).**

| Question | Answer |
|---|---|
| Is such a declaration **required in effect**? | **Yes** — CEV-08/CEV-10: without it, L2/L3 breach I15 and cannot satisfy I18 |
| Is it **permitted in principle**? | **Yes** — CEV-05/CEV-06: the architecture already treats a stream's shape as carrying a second-order fact, and CEV-15 shows the same shape contemplated elsewhere |
| Does a **construct exist today**? | **No** — no closure marker for the lifecycle stream, and no completion declaration in either stream (CEV-17) |
| Is it **forbidden**? | **No** — nothing in r33 prohibits it; nothing in K-R12 restricts the stream to state-like entries |

**CEV-19 (Inferred — the honest summary).** Q-1's answer is not "the architecture permits it" or "it does not."
It is: **the architecture obliges the distinction and has no means to express it.** That is a stronger finding
than OBE-36 anticipated, and it applies to L2 and L3 *equally* — so it does not discriminate between them.

**CEV-20 (Inferred — a consequence for the falsifier).** OBE-26's L2 falsifier ("closure cannot be specified
without a successor relation") is now sharper: the question is not whether a successor relation exists, but whether
a declaration that is **neither state nor transition** is admissible in a stream that K-R12 describes only as
carrying status/scope/class changes.

---

# Part B — Validity versus recording

## 7. The distinction is already canonical

**CEV-21 (Source-derived — I17, `[X]`).**

> "**A result earns authority from the conditions under which it was established, not from when or where it was
> written.** … never 'it was computed most recently', always 'it was computed under these conditions'."

**CEV-22 (Source-derived — I7, `[X]`).** "Retrievability does not confer authority."

**CEV-23 (Calculated — what I17 settles for Q-2).** I17 states the separation L3 needs **as a principle, in the
canon, with an `[X]` tag**: standing comes from conditions of establishment, explicitly **not** from the act of
writing. So "recording an act confers status" is not merely undesirable — it is the exact error I17 names.

**CEV-24 (Inferred — the qualification).** I17 constrains **derivation**: authority must be traced to conditions.
It does **not** specify which conditions, who evaluates them, or what happens when they are unknown. The principle
is present; the **predicate** is not.

## 8. The default is fail-closed, which makes the gap sharper

**CEV-25 (Source-derived — I15, permissions clause).** I15's uniform list includes **"unspecified permissions"** and
requires that they mean *"not established"*, never *"established as none."*

**CEV-26 (Source-derived — K-R5, `:207`).** "Every retrieval call names its context. **Default is
`character_perspective`** (fail closed)."

**CEV-27 (Calculated — the combined effect).** For an act whose authorisation is unspecified, the house rules give
**not established**, and the house default is **refuse**. Applied to L3's derivation, an act of unknown validity
cannot be treated as valid — so the derivation cannot emit the state the act would otherwise confer.

**CEV-28 (Inferred — the consequence, which is the real finding for L3).** Under L3 as specified, derivation must
answer *"was this act valid?"* **before** it emits any status. Today it cannot answer, ND-1 is silent on who may
cause a lifecycle change, and failing closed means emitting **no status at all** for any record whose act validity
is unknown. So L3's exposure is not "some states may be wrong"; it is that **a fail-closed L3 emits nothing, and a
non-fail-closed L3 commits exactly the inversion I17 forbids.** There is no third option available without a
validity predicate.

**CEV-29 (Inferred — why this is not a criticism of L3).** The same gap sits under L1 and L2 as a *provenance*
question rather than a *derivation* question: they record what the stream says and can carry the conditions
alongside, so an unknown condition degrades traceability without destroying the state. L3 makes the act the ground,
so the unknown condition becomes a hole in the representation. That is the structural cost OBE-18 named, now with
its mechanism.

## 9. What "demonstrably enforceable" would require

**CEV-30 (Inferred — stated without designing it).** Three things, none of which exists:

1. A **validity predicate** that is not "the act is present" — i.e. conditions checked independently of the record
   whose standing they support.
2. An **evaluation point** — when the predicate runs, given that exits may be evaluated at formation (P-EVAL,
   EKR-11) or later.
3. A **failure mode for an unevaluable act** that is visible, in the shape K-R2 already demands: "A projection that
   cannot be reconstructed from its events is an **integrity failure**, not a policy question."

**CEV-31 (Inferred — the one thing that would make L3 defensible).** If (1) and (2) exist and (3) reports rather
than guesses, L3's inversion is closed: recording an act supplies a *claim*, and derivation supplies standing only
where the predicate holds. Without (1), no amount of care at (3) helps, because the derivation cannot tell a valid
act from an invalid one.

## 10. Q-2 verdict

**CEV-32 (Calculated — four-part verdict).**

| Question | Answer |
|---|---|
| Is the distinction **canonically required**? | **Yes** — I17, with I7 as its sibling |
| Is the **default** consistent with it? | **Yes** — I15's permissions clause plus fail-closed (K-R5 family) |
| Is it **enforceable today**? | **No** — no validity predicate, no evaluation point, ND-1 open |
| Does an implementation exist that could demonstrate it? | **No** — nothing in the inspected text; not an implementation question (CEV-35) |

**CEV-33 (Inferred — summary).** Q-2's answer is the mirror of Q-1's: **the requirement is stated and the mechanism
is absent.** L3 remains **conditional**, as the owner assessed, and the condition is precise: a validity predicate
independent of act presence.

**CEV-34 (Inferred — and the condition is not deferred indefinitely).** Unlike Q-1, Q-2 **cannot** be left to a
later design stage without deciding ND-1, because a validity predicate *is* an answer to "who may cause a lifecycle
change." So L3's falsifier is not merely open; it is **entangled with an owner decision**, which is the strongest
form of "conditional" available in this workstream.

---

## 11. What this record does not do

**CEV-35 (Observed).** No mechanism, no field, no event type, no marker, no predicate is proposed. Part A states
what is required and what exists; Part B states what is required and what is missing.

**CEV-36 (Observed).** ND-1 is **not answered, narrowed or assumed**. P-EVAL (EKR-11) is **not** adopted. §8.1 is
untouched. Entry remains constrained and identical across options (ODB-63).

**CEV-37 (Observed).** No option is selected, ranked or preferred. CEV-11 and CEV-13 explain a structural property
of the three; they do not weigh it.

**CEV-38 (Observed).** The 118 requirements are unchanged. Nothing here becomes a requirement: the record reports
which **existing** invariant already bears on each question.

---

## 12. The refined owner question, and one caution

**CEV-39 (Observed — owner wording, recorded verbatim).** The decision is not "which failure would you accept" but:

> *"Which failure mode is least tolerable, considering both its likelihood and the difficulty of detecting and
> recovering from it?"*

**CEV-40 (Inferred — why the refinement changes the shape of the answer).** The original form implied equal
exposure across the three options. Likelihood is **not** equal: L1's unwarranted state requires an omission *and* an
audit failure; L2's requires only a lost tail, which partial histories make ordinary; L3's requires only a write
path. Recovery differs too — an unwarranted state can be re-grounded from provenance, a truncated history **cannot
be reconstructed from what is missing**, and an inversion can be re-derived once a predicate exists. The refined
question therefore has a defensible answer, where the original did not.

**CEV-41 (Observed — carrying the owner's caution).** "Currently unthreatened by an open question" (OBE-27) must
not be read as "safest overall." Certainty about prerequisites and safety of the resulting choice are **different
claims**; this record adopts neither. L1's vocabulary and state-representation costs are untouched by CEV-11.

---

## 13. Limits

**CEV-42 (Observed).** Basis: r33 §2 (D1, D4, D6), §3 (K-R2, K-R17, K-R19), §4 (I7, I15, I17, I18), §5, §10.2,
§11, §11.1. **Not inspected:** whether any implementation already enforces a validity predicate (out of scope by
instruction); r32's bytes; whether `attribution.kind` could carry closure meaning (none of the seven values
suggests it, but this was not exhaustively tested); the frozen preregistration's treatment. **Weakest point:**
CEV-10's claim that I18 applies to representation rather than to subsystems. I18 is written about **subsystems**; applying
it to a representation stance is an **Inference**, and it is the load-bearing step in CEV-18's first row.

**End of record.**
