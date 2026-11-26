# MEMORY-09 — Lifecycle Completeness: What "Complete" Can Legitimately Claim

## Bounded evidence on closure terms, the K-R17 and K-R19 precedents, and what transfers

**Status:** evidence record. **No mechanism, marker, closure rule, schema, field, event type, implementation or new
requirement.** Read-only: no canonical edit, no test, no probe, no ledger entry. OD-b remains **settled as L2**
(ODL-01); ODL-39's K-R17 tension remains **OPEN** and is **not** treated here as inherited policy.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Scope | Q-1…Q-6 as bounded by the owner: meaning of *complete*; admissibility of a "no further entries expected" declaration (OBE-36); what K-R17 and K-R19 each guarantee; what transfers to the lifecycle stream; what stays unestablished |
| Excluded | mechanism, marker, closure rule, schema, implementation, new requirement, selection among closure approaches |
| Basis | r33: §2 (D1–D7), §3 (K-R1, K-R2, K-R11–K-R20, K-R32–K-R38), §4 (I15, I18), §5, §10.2, §11 |
| Prior records | `...CLOSURE-AND-VALIDITY-EVIDENCE...` (CEV), `...OWNER-DECISION-OD-B-L2...` (ODL), `...ENTRY-CONDITION...` (ECI) |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. Scope, and why the vocabulary has to be fixed first

**LCM-01 (Observed — what this record is for).** The owner asked for the closure/completeness record before any
evaluation of how an L2 projection behaves. The reason is structural: a projection always emits *some* state, so the
question is never whether the fold produces an answer, but **what the history it folded was entitled to claim**.

**LCM-02 (Observed — ODL-39 stays open).** This record does **not** decide that the lifecycle stream must adopt
K-R17's rejection-at-write discipline. It answers what that discipline *guarantees in its own stream*, so the
question can later be asked honestly. Section 6 records the transfer as a **shape and a prerequisite list**, not as
an inheritance.

**LCM-03 (Inferred — four claims that have been used interchangeably).** The L2 discussion uses "complete",
"closed", "terminal" and "the end of the log" as if they were one property. They are four:

| Claim | What it asserts | Refuted by |
|---|---|---|
| **Well-formed** | every entry in the stream is admissible under that stream's write rules | a write that should have been rejected |
| **Complete** | the stream contains every entry that occurred in its domain | evidence of an entry that occurred and is absent |
| **Terminally closed** | the domain has ended, and that ending is recorded | an entry that occurs after the recorded end |
| **Currently visible end** | the last entry as stored at read time | the next write |

**LCM-04 (Inferred — these are independent, in one direction each).** Well-formedness does **not** imply
completeness (§4). Completeness does **not** imply closure: a record may be completely recorded and still open to
further transitions. Visibility implies none of the other three — it is a fact about the artifact, not about the
world. Only the reverse implications are safe: a closed stream is complete *by definition of its domain ending*, and
a complete and open stream still has a visible end.

**LCM-05 (Inferred — Q-1 is a definitional question, not a design one).** "Complete" is a correspondence claim
between the stream and its domain; "no further entries currently visible" is a claim about the stored artifact. They
differ in what could refute them (LCM-03), and therefore in what evidence they require. That difference is prior to
any choice of representation, which is why it can be settled while L2's mechanism stays unauthorised.

**LCM-06 (Observed — canon does not define stream completeness).** A term search of r33 for definitions of
completeness, truncation or closure of a *stream* returns none. I15 defines what an absence **means**; I18 defines
what a subsystem owes when it might have omitted something; neither defines when a history is complete. So Q-1 can be
answered as a **distinction**, not as a quotation.

---

## 2. Q-1 — what *complete* means, distinct from a visible end

**LCM-07 (Source-derived — I15, `[X]`, r33 `:301`).**

> "**Absence of information is never treated as information about absence.** An unpopulated field, an empty list, a
> missing record and an unwritten schema each mean 'not established', never 'established as none'. A system that
> cannot safely determine an answer says so instead of manufacturing certainty."

> "**Applies uniformly to:** missing acquisition history, empty knowledge sets, unspecified permissions, undeclared
> action requirements, missing asset versions, and incomplete retrieval evidence."

**LCM-08 (Source-derived — I18, `[X]`, r33 `:306`).**

> "**The system may not use absence of evidence as evidence that nothing happened.** Every subsystem must name the
> **independent evidence** that would reveal its own silent omission — evidence it does not produce itself. A
> subsystem that cannot answer 'what would show I failed to act?' is incomplete, not merely untested."

**LCM-09 (Inferred — the operative content of "complete").** Taken together, the two invariants make completeness a
claim that **the stream cannot attest about itself**. I18 requires evidence *it does not produce itself*; a stream's
own last entry is precisely evidence it produces itself. So the only reading of "complete" that survives I18 is:

> **complete ≙ no silent omission remains undetected in this history**

— not "the log ends here", and not "no further entries are visible". Completeness is therefore **evidence-relative**
and external, while visibility is **artifact-relative** and internal.

**LCM-10 (Calculated — what each claim has behind it in canon).**

| Claim | Guaranteed by | Status |
|---|---|---|
| Well-formed | K-R16 (duplicate ids rejected at write), K-R17 (unsupportable revoke rejected at write) | **`[X]`, present** |
| Complete | nothing in r33 | **absent** — no rule, no check, no vocabulary |
| Terminally closed | deletion is a recorded, authorised end (D1; §11.1 `[P]`) and `via=revoked` closes an episode | **partial** — per-episode closure exists; record-level closure does not (CEV-16/CEV-17) |
| Currently visible end | the append-only log itself | **present**, and says nothing about the other three |

**LCM-11 (Calculated — four situations, one observation).** "No further entries currently visible" is fully
compatible with each of:

1. a rumour that will never advance (`:612`, "rumour — stays provisional");
2. a record whose successor occurred but was never written;
3. a stream whose tail was lost or truncated;
4. a record still awaiting its next transition.

**Nothing in the stored form distinguishes them.** That is the closure problem stated as a definitional fact, and it
is independent of L2: any representation that reads a state off the end of a history faces it.

**LCM-12 (Observed — two adjacent uses of "completeness" that do not help).** K-R32 requires replay to "fail loudly"
rather than "emit a verdict from partial agreement", and §10.2.2 distinguishes `reproducible` / `auditable_only` /
`expired`. Both are structurally about not reporting a partial artifact as a faithful one — the right shape — but
**§10.2 is explicitly "proposed, not adopted"**, and both concern evaluation replay, not lifecycle. They are
precedents for the idea (CEV-14/CEV-15), not constructs available to the lifecycle stream.

---

## 3. Q-2 — is an independently meaningful "no further entries expected" declaration admissible?

**LCM-13 (Source-derived — K-R12, `[X]`, r33 `:111`).**

> "Memory lifecycle is its **own authoritative stream** — `established | superseded | corrected | retracted`, plus
> `scope` and `knowledge_class` changes — and a memory's current status, class and scope are a **projection of that
> stream**, not mutable authoritative fields."

**LCM-14 (Calculated — the enumeration cannot settle Q-2).** K-R12's list admits two readings: an **exhaustive**
vocabulary of lifecycle entries, or an **exemplary** list of the changes it cares about. This brief's ODB-05 already
adopts the second reading to explain the *absence of `provisional`* — a gap the list cannot have meant to draw,
since §11 makes `provisional` the entry state. Reading the same list as exhaustive now, in order to exclude a
completion entry, would contradict ODB-05. **K-R12 therefore neither admits nor forbids a non-transition entry; on
this point it is silent** (CEV-18's fourth row).

**LCM-15 (Source-derived — the one prohibition that does bite).** K-R12's second clause excludes a **mutable
authoritative field** for status, class and scope. Whatever else follows, a stored completion flag *on the record*
is excluded by this clause, because it would be a mutable authoritative field about lifecycle state. That is a
constraint on **where** such a declaration could not live — not permission for anywhere else.

**LCM-16 (Inferred — the decisive test is independence, not expressibility).** A declaration is "independently
meaningful" only if it is meaningful **without being derived from the absence it describes**. Otherwise it restates
the problem: a stream that says "nothing follows" because nothing followed has added no information. I18 states the
test exactly — evidence *"it does not produce itself"*. Applied to OBE-36:

| Candidate source of the declaration | Independently meaningful? |
|---|---|
| the same process that writes the transitions | **no** — it is the omission reporting itself |
| a reader deriving "no successor" from the empty tail | **no** — this is the inference I15 forbids (CEV-08) |
| a recorded act by a party outside that write path | **Unestablished** — no such act exists in r33, and nothing says one is admissible |
| a kind-level policy ("rumours do not advance") | **Unestablished** — CEV-13's point: this is a statement about the kind, not a fact the stream attests |

**LCM-17 (Calculated — Q-2 verdict, four parts).**

| Question | Answer |
|---|---|
| **Required in effect?** | **Yes** — without it, an L2 fold reads absence as a terminal state, which CEV-08/CEV-10 find to be an I15/I18 breach |
| **Forbidden by canon?** | **No** — nothing in r33 restricts the stream to state-like entries; K-R12 is silent (LCM-14), and only the mutable-field clause constrains placement (LCM-15) |
| **Expressible today?** | **As an entry, plausibly; as independent evidence, not yet** — LCM-16's test has no candidate that passes |
| **Does a construct exist?** | **No** — no closure marker, no completion declaration in either stream (CEV-17) |

**LCM-18 (Inferred — what OBE-36 now amounts to).** OBE-36 asked whether a stream *can express* "no further entries
by design". The evidence relocates the difficulty: **expression is not the obstacle; independence is.** A recorded
declaration proves that a decision was made, not that the history is complete — and those are different claims
(LCM-03). OBE-36 stays open, but it is now a question about **what could count as independent evidence**, which is
also the question I18 asks of every subsystem.

---

## 4. Q-3 — what K-R17's rejection-at-write discipline actually guarantees

**LCM-19 (Source-derived — K-R17, `[X]`, r33 `:137`).**

> "`revoked` is valid only when the pair's ordered log **currently ends** in a learn event. A revoke with no open
> episode is **rejected at write**, not silently ignored, because a silent no-op hides formation bugs."

**LCM-20 (Source-derived — the surrounding write rules).** K-R15 orders a pair's events totally by `(turn, seq)` with
"no wall clock and no insertion order"; K-R16 rejects duplicate ids at write; K-R18 permits a learn at any time, and a
learn while an episode is open is a *continuation*, not a new episode.

**LCM-21 (Calculated — four guarantees, stated precisely).**

1. **No unsupportable revoke enters the log.** Every `revoked` in a stored log was preceded by at least one learn
   since the previous revoke. (K-R18 means the log does **not** strictly alternate learn/revoke; it means each revoke
   has an open episode behind it.)
2. **The fold is total.** Because the episode model is learn-opens / revoke-closes (K-R1, R1c), the last event always
   determines `known(pair)` unambiguously. There is no log shape in which the fold is undefined.
3. **Failure is loud and early.** The rule's stated purpose is that a malformed write is *not* absorbed silently —
   the failure surfaces at the writer, where a formation bug is still attributable.
4. **The reason is diagnostic, not epistemic.** K-R17 exists because "a silent no-op hides formation bugs" — it
   protects against **writes the log cannot interpret**, not against writes the log never received.

**LCM-22 (Inferred — what it does not guarantee).** **Nothing about completeness.** A log that lost its entire tail
is still perfectly well-formed; a log whose final learn was written and whose subsequent revoke was never written is
well-formed and *states the wrong current knowledge*. K-R17 constrains **admissibility**, not **coverage**. This is
the distinction LCM-03 draws, and it is why K-R17 cannot be cited as a closure precedent without qualification.

**LCM-23 (Inferred — the guarantee's real scope).** K-R17 is a rule about **transitions the current log cannot
support**. It presupposes the ability to compute the current state (the fold) and an unambiguous validity condition
("the ordered log currently ends in a learn event"). It is an **intra-stream consistency** discipline. It says
nothing about whether the log is the whole log, and it is structurally incapable of doing so: a missing write leaves
nothing for a write-time gate to reject (I18's point).

**LCM-24 (Calculated — the lifecycle analogue, and its missing prerequisite).** The nearest lifecycle analogue is
"reject a transition the projected state cannot support". Two prerequisites:

| Prerequisite | State |
|---|---|
| a total fold over lifecycle events | available in principle under L2, and required to be deterministic (ODL-28) |
| an admissible-transition set per kind | **not settled**: `active`/`established` (ODB-04), P-EVAL (EKR-11), DN-6a/DN-6b |

So even the **transferable half** of K-R17 cannot be stated today without a decision this record is not authorised
to make. Recorded as a prerequisite finding, not a proposal.

---

## 5. Q-4 — what K-R19's absence semantics actually guarantee

**LCM-25 (Source-derived — K-R19 and K-R20, `[X]`, r33 `:147`).**

> "A revoked pair has **no edge**, not an edge with a "revoked" flag. "Forgot" is the absence of the current
> relation; the fact that it was once known lives in the log, not in a half-state row."

> "`revoked_turn` is therefore **not a field**. It is derived on demand — the timestamp of the last revoke event —
> and appears only in the `knowledge_history` inspection view."

**LCM-26 (Calculated — the mechanism, stated exactly).** "Forgot" is **not** read off the absence of events. It is
the computed consequence of a **recorded** event: the last event is a revoke. The projection's negative case has a
**positive witness in the log**. That is what makes the absence meaningful: it is produced by a fact, not by a gap.

**LCM-27 (Calculated — and K-R19 covers less than "all absence").** Two distinct situations yield no edge:

| Situation | Log | What the projection can say |
|---|---|---|
| never acquired | empty (K-R1: "No edge exists without at least one event") | nothing — this is I15's "not established", not "forgot" |
| forgot | ends in a revoke | "forgot", because a learn preceded it |

So K-R19's sentence fixes the meaning of absence **for a pair that had a history**; the never-acquired case is
governed by K-R1 and I15 instead. The projection alone cannot tell the two apart — **the log can**, and only by being
read. Absence in the acquisition stream is therefore *overloaded at the projection and disambiguated at the log*.
This refines CEV-04/CEV-05: the precedent is "a negative conclusion is grounded in a recorded fact", not "absence
has one meaning everywhere".

**LCM-28 (Inferred — what K-R19 does not guarantee).** Not completeness: both situations may also arise from a lost
tail, and the log cannot say so either. Not that the *record* exists. Not anything about awareness of the
forgetting — K-R3 and R4 put that in a further acquisition event. Its guarantee is narrower and exact: **no
half-state row; history preserved elsewhere; and every "forgot" conclusion traceable to a revoke that was written.**

**LCM-29 (Inferred — the disanalogy with the lifecycle, which is the central finding).** The acquisition stream has a
**closing act** — `revoked` — that is itself a recorded, meaningful terminal event. "Rumour — stays provisional"
(`:612`) has no such act: terminality there is **the absence of any further event**, not the presence of a closing
one. So K-R19's pattern — *absence is computed from a positive recorded fact* — **does not transfer**, because the
lifecycle case has no positive fact to compute it from. To get one, the architecture would have to admit a closing
act or a completion declaration for lifecycle records; that is the decision OBE-36 poses and this record does not
make.

---

## 6. Q-5 — what transfers to the lifecycle stream, what does not, and why

**LCM-30 (Calculated — transfer table).**

| Precedent | Guarantees in its own stream | Transfers to lifecycle as | Why / why not |
|---|---|---|---|
| **K-R17** rejection at write | no unsupportable revoke; total fold; loud failure | **posture, conditionally** — reject transitions the projection cannot support | the rejectable case is an *impossible* transition; prerequisites are unset (LCM-24) |
| **K-R17** stated purpose | formation bugs surface at write | **as far as it goes** — the same argument applies to lifecycle writes | but the failure it prevents is a bad write, not an absent one (LCM-22) |
| **K-R19** absence semantics | "forgot" grounded in a recorded revoke | **as a shape only** — a negative conclusion must rest on a recorded fact | lifecycle has no closing act to rest it on (LCM-29) |
| **K-R19** no half-state | history lives in the log; no `revoked_turn` field | **yes, by analogy** — a completion flag on the record is excluded by K-R12 (LCM-15) | same principle, already binding |
| **K-R2** integrity failure | an unreconstructable projection is a failure, not a policy question | **yes, and already in force** (K-R14 rebuilds each stream) | it governs reconstruction **from present events**, not the completeness of events (LCM-31) |
| **K-R32** no partial agreement | fail loudly rather than report partial as faithful | **shape only** | §10.2 is *proposed, not adopted*; scope is replay (LCM-12) |
| **K-R38** no silent gap | "artifact absent and no tombstone" unreachable | **shape only** | proposed, not adopted; scope is retention (LCM-12) |
| **I15 / I18** | absence is never evidence; name independent evidence | **already binding on lifecycle** | this is why closure is an obligation, not why a mechanism exists (LCM-36) |

**LCM-31 (Inferred — the asymmetry that decides the framing).** K-R2's vocabulary detects **corruption**, not
**truncation**: a truncated lifecycle stream reconstructs *perfectly*, because the projection it yields is internally
consistent with the events that remain. So under L2, a lost tail is **not an integrity failure by K-R2's
definition**. Whatever addresses it must make a different claim than reconstructability — which is exactly why
closure is its own obligation (ODL-26) rather than a corollary of an existing invariant.

**LCM-32 (Inferred — why rejection-at-write cannot close the gap).** A write-path gate can only reject what arrives.
The lifecycle's characteristic failure is a transition that **never arrives**. This is not a defect of K-R17; it is
the boundary of any write-time discipline, and it is the same boundary I18 names when it demands evidence "it does
not produce itself".

**LCM-33 (Inferred — what genuinely transfers, stated as narrowly as the evidence allows).** Three things, and no
more:

1. **A posture**: prefer making an unsupportable entry unrepresentable over interpreting silence (K-R17's shape).
2. **A requirement shape**: a negative conclusion must rest on a recorded fact (K-R19's shape) — which the lifecycle
   stream currently cannot satisfy (LCM-29).
3. **An existing obligation**: I15 and I18 are already binding, so the closure obligation is not new (CEV-18) even
   though no construct exists.

None of the three is a rule for the lifecycle stream, and this record does not convert any of them into one.

**LCM-34 (Inferred — what does *not* transfer, named explicitly).** Rejection of the impossible (LCM-32); "one
meaning for absence" (the lifecycle has four possible meanings — LCM-11); and the assumption that the acquisition
stream's precedent settles lifecycle closure. That last one is ODL-39's tension, left **open**.

**LCM-35 (Observed — the nearest true precedent is unavailable).** §10.2.2's tombstone discipline is structurally the
closest thing in r33 to a closure property: "for any observer at any instant, the state 'artifact absent and no
tombstone' must be unreachable". It is **proposed, not adopted**, and scoped to retained replay artifacts. It shows
the architecture can conceive of the property; it supplies no construct (CEV-15).

**LCM-36 (Observed — the invariants, not the mechanism).** I15's semantics for absence and I18's independence test
are the reason the distinction is obligatory. They are quoted in §2; they are **not** new requirements and this
record adds none (CEV-38).

---

## 7. Q-6 — what remains unestablished

**LCM-37 (Unestablished — K-R12's enumeration).** Whether the lifecycle entry list is exhaustive. K-R12 is silent
(LCM-14), ODB-05 leans non-exhaustive for `provisional`, and no text settles the question either way.

**LCM-38 (Unestablished — the sharpest point).** Whether **any** declaration can satisfy I18's independence test for
completeness, given that the party able to declare the history complete is usually the party whose omission is in
question. If it cannot, closure may be unachievable *in-stream* by construction rather than merely unimplemented —
and that would change what L2 can promise, not just what it lacks.

**LCM-39 (Unestablished — completeness without a domain end).** Whether "complete" is definable for an open domain at
all, or only "complete so far". CEV-13's observation bears on this: no option currently expresses terminality as
established fact; L1 reports a last asserted state, L2 and L3 infer one.

**LCM-40 (Unestablished — the rejection predicate's prerequisites).** Whether rejection-at-write can be specified
before an admissible-transition set exists. It cannot today: `active`/`established` (ODB-04), P-EVAL (EKR-11) and
DN-6a/DN-6b are open (LCM-24).

**LCM-41 (Unestablished — lawful truncation).** Whether retention/archival (§10.2, Q12) creates a truncation the
lifecycle closure question must account for. Deletion is a recorded, authorised end (CEV-16), but §10.2 is proposed
and Q12 is open, so the interaction is unexamined.

**LCM-42 (Unestablished — per-kind terminality as a kind property).** Whether "rumours do not advance" can be stated
as a property of the kind rather than as a fact about a stream (CEV-13; LCM-16's last row). It is the option that
would avoid inventing a closing act, and nothing in r33 decides it.

**LCM-43 (Observed — unchanged from the entry record).** Whether the entry relation already carries completeness
meaning (ECI-25/ECI-27) remains as recorded; this record adds no entry-side finding.

---

## 8. What this record does not do, and its limits

**LCM-44 (Observed).** No mechanism, marker, closure rule, completion field, event type, predicate, retention
policy or fail-closed rule is proposed. No option is selected, ranked or preferred, and nothing here becomes a
requirement.

**LCM-45 (Observed — ODL-39 preserved).** This record does **not** conclude that the lifecycle stream must adopt
K-R17's discipline. It concludes something narrower and less convenient: the transferable part is a **posture**, its
prerequisites do not exist, and the part that would actually close the gap (a closing act) has no lifecycle analogue.
The tension stays open for the owner.

**LCM-46 (Observed — basis and method limits).** Evidence is from r33's §2, §3, §4, §5, §10.2 and §11 only. Not
inspected: v9 for a lifecycle-completeness statement; r32's bytes; any implementation's stream handling; the frozen
preregistration's treatment. The r33 sections relied on are `[X]` except §10.2, which is `[P]` and labelled as such
wherever cited. The term search for a definition of completeness (LCM-06) is term-based and could miss a paraphrase.

**LCM-47 (Observed).** No canonical file was edited; r33, v9, the owner ledger and the conformance map are
unchanged; the 118 requirements are unchanged; 0 tests were run.

**End of record.**
