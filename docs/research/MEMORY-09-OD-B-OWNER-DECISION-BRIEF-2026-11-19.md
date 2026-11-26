# MEMORY-09 — OD-b Owner Decision Brief (Revision 4)

## How the Lifecycle Representation Relates to the Documented Status Progression

**Status:** decision brief for the owner, **revision 4** — **OD-b is settled: the owner selected L2 on 2026-11-20**
(§5.4; `MEMORY-09-OWNER-DECISION-OD-B-L2-2026-11-20.md`). The brief records the selection and the obligations it
inherits. **Not a schema proposal** — no field list is offered (§12).
Read-only: **no schema change, implementation, test, probe, model call, canonical edit, registry change or
ledger entry is authorised**, by this revision or by the selection itself (ODB-29).

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Revision | **4** — **owner selection recorded (§5.4)**. Revision 3's narrowing is preserved unchanged; no earlier analysis rewritten |
| Decision | **OD-b = L2 — transition-only.** Owner, 2026-11-20. Rationale, failure-risk ordering and carried-forward obligations: `MEMORY-09-OWNER-DECISION-OD-B-L2-2026-11-20.md` |
| Not selected | **L1** and **L3** — not selected, and **not** shown unsound (ODB-76) |
| Retired by evidence | Reading B (review state) and Reading C (not a lifecycle state) — §4 |
| Preserved separately | **C′** — whether verification information is represented independently of status (§10) |
| Settled elsewhere | **A2**, **D-3b**, **D-3a**, **DN-4**, **OD-a** — none reopened |
| Open elsewhere | **DN-1** (widened), **OD-g**, **DN-6a**/**DN-6b**, **C′** (§10), **P-EVAL** (§5.3) |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 0. What changed in revisions 2–4

**ODB-37 (Observed — changelog).** Revision 1 asked *"is `provisional` a truth/status state or a review
state?"*. That question is now answered by canon rather than by choice, and the brief's job changes:

| # | Change | Cause |
|---|---|---|
| 1 | §11's transition diagram added as the **controlling** statement | proposed-answer check, ODBC-03 |
| 2 | **Reading B retired** — not unsupported merely, but contradicted by §11's evidence-kind exits | ODBC-11 |
| 3 | **Reading C retired** — §11 makes `provisional` the entry state | ODBC-08 |
| 4 | **Producer finding corrected** — formation produces `provisional`; what is missing is the typed representation | ODBC-05/16/17 |
| 5 | **OD-b reframed** to the representation question (§5) | owner direction |
| 6 | **C′ extracted** as a separate design question, not a third interpretation (§10) | owner direction |
| 7 | Caution added: formation ≠ cause/evidence/occurrence/authority (§7) | owner direction |
| 8 | Experimental boundary restated (§11) | owner direction |
| 9 | **Entry, cross-stream dependency and ordering removed from the discriminators** — they are constrained, not chosen; §5/§5.1 rewritten | ECI-31/ECI-32 |
| 10 | §8.1's `status=active` retained as an **independent** open question, outside OD-b | `...ENTRY-UNIFORMITY...` EKR-19 |
| 11 | **L2 selected by the owner (2026-11-20)** — §5.4 records the selection, its rationale, the owner's failure-risk ordering and the obligations L2 inherits; the L1/L3 assessment is retained as written | owner decision |

**ODB-38 (Inferred).** The reframing is not a narrowing for convenience. Revision 1's question had two answers,
one of which canon forbids; leaving it open would have invited a decision canon has already excluded.

---

## 1. What this brief is, and is not

**ODB-01 (Observed).** This document states the documented progression, the representational options, and their
consequences. It **does not draft the lifecycle schema**, does not name fields, and does not resolve OD-b.

**ODB-02 (Inferred).** The discipline matters because a field list is an answer in disguise, and would now
encode an answer to OD-b, to `active`/`established`, to `forgotten` and to OD-g at once.

---

## 2. The documented progression

**ODB-39 (Source-derived — verbatim, r33 `:600–612`, §11 "Lifecycle and retraction").**

```
evidence turn
   └─> provisional (status, class, attribution set)
          ├─> active        (episodic, preference, validated fact, belief)
          ├─> superseded    (replaced; history kept)
          ├─> corrected     (earlier statement wrong; correction chain)
          └─> retracted     (asserted, now known false; acquisitions retained)
```

"Promotion per kind: episodic — on witnessing; preference — immediate for holder; belief — never promotable to
fact; **rumour — stays provisional**, repetition is not corroboration; world fact — authority or corroboration;
absence — attributed. Rejected: blanket three-source rule, promotion by retrieval, decay altering truth,
deleting history."

**ODB-40 (Observed — the five structural facts this fixes).**

1. **Entry** — formation produces a `provisional` record from the evidence turn, writing `class` and
   `attribution set` with it.
2. **`provisional` is the entry state**, and `active` is its first child.
3. **Exit is evidence-kind dependent** — witnessing, immediacy, authority/corroboration, per kind.
4. **Exit is not caused by** retrieval, repetition, approval, decay, or a blanket corroboration count
   (all explicitly rejected at `:612`).
5. **A legitimate terminal state exists** — "rumour — stays provisional" means no outgoing transition is
   required, ever.

**ODB-41 (Inferred — the requirement this creates).** Any lifecycle representation must be able to express
**entry at formation**, **exit conditional on evidence kind**, and **a state with no successor**. Those three
are structural, not stylistic; an option that cannot express one of them is wrong.

**ODB-42 (Observed — retraction preserves acquisition history).** §11's `retracted` node reads "asserted, now
known false; **acquisitions retained**", and r33 `:616` — "Retraction applies to records, not to minds (R4)."

**ODB-43 (Observed — `provisional` must not be presented as verified).** r33 `:283`, Assembly prohibited
behaviours: "present provisional as verified".

**ODB-44 (Observed — D4).** r33 `:72` — "Provisional records are retrievable for characters holding them,
attributed, never promoted by retrieval." **Unverified is not inactive.**

---

## 3. Three enumerations, one progression

**ODB-03 (Observed — the enumerations).**

| Source | Vocabulary | Line |
|---|---|---|
| r33 §5 table, "Truth/lifecycle status" | `provisional / active / superseded / retracted / corrected` | `:322` |
| r33 **§11 diagram** | entry `provisional` → `active` / `superseded` / `corrected` / `retracted` | `:600–612` |
| r33 **K-R12** (lifecycle stream) | `established \| superseded \| corrected \| retracted`, plus `scope` and `knowledge_class` changes | `:111` |

**ODB-04 (Calculated — the residual mismatch).** §5 and §11 agree on five names. **K-R12** agrees on three and
differs in both directions: it omits `provisional` and `active`, and introduces **`established`**, which
appears in neither §5 nor §11 as a state.

**ODB-05 (Inferred — corrected explanation).** K-R12's omission of `provisional` is **incomplete enumeration**
of the progression, not a category boundary (ODBC-15). K-R12 is a list of the transitions it cares about; the
entry state is missing from the list, not from the lifecycle.

**ODB-06 (Superseded — corrected 2026-11-19).** The earlier claim that `provisional` "has no producer
anywhere" (A2-43(i)) is **withdrawn**: formation produces it at `:604`. The real gap is a **missing typed
representation** of the transitions, not a missing producer.

**ODB-07 (Observed — `forgotten` is in none of the three).** v1's `Event` carries `forgotten`; **DN-4** requires
forgetting as a distinct operation and event. No enumeration above contains its name. Unresolved independently
of OD-b.

---

## 4. What canon settles, and what is retired

**ODB-08 (Source-derived — settled by evidence, for owner confirmation).** `provisional` **is a lifecycle
state**: the entry state of the documented progression, produced by formation, exited only by evidence kind.
This is read from §11 `:604`, §5 `:322`, D4 `:72`, C4 `:446`, `:612`, S12 `:941`, S26 `:956`.

**ODB-09 (Superseded — Reading B retired).** A review-state reading is **contradicted**, not merely
unsupported: §11's exit conditions are evidence kinds, and it explicitly rejects promotion by retrieval and a
blanket corroboration rule. No reviewer, approver, review queue or approval step exists in r33 or v9. Adopting
B would require inventing an actor canon has never named — the outcome the owner's "review authority must not
appear by implication" protection forbids.

**ODB-10 (Superseded — Reading C retired as a competing reading).** §11 makes `provisional` the lifecycle's
entry state and `active`'s parent, so "not a lifecycle state" fails. C explained only the K-R12 omission, and
that is better explained as incompleteness (ODB-05).

**ODB-11 (Inferred — Reading C's residue survives as C′).** The value also carries a **verification**
dimension that assembly reads (`:283`, C4 `:446`) and that conflict materiality reads (`:459`). Whether that
dimension needs its own representation, independent of status, is a **separate design question** (§10) — not a
third OD-b interpretation.

**ODB-12 (Observed — one field, two dimensions).** §5's row is labelled "**Truth**/lifecycle status" and lists
values from both dimensions. That label is the origin of the conflation; C′ is the question that remains about
it.

---

## 5. The decision that was before the owner

**ODB-45 (Inferred — resolved 2026-11-20, ODB-72).** The progression is fixed. What was undecided is **what the
lifecycle structure records so that the progression is recoverable** — states, transitions, or grounding acts. The
three options below are **representational stances**, deliberately stated without fields. The owner selected **L2**;
the options are retained as stated because §5.1–§5.3's reasoning is what explains the choice.

**ODB-62 (Observed — revision 3 change).** The selection criteria have narrowed. **Entry location, cross-stream
dependency and cross-stream ordering are no longer discriminators**: the first is constrained by
K-R11/K-R12/K-R2/K-R14/D4 (ECI-16/17), the second follows from it (ECI-18), and the third is never needed
(ECI-34/35). What remains is §5.2's three questions.

### 5.1 What is constrained out of the decision

| Criterion | Status | Why it cannot discriminate |
|---|---|---|
| **Where the entry condition lives** | **constrained** — lifecycle-side | acquisition-derived status is forbidden by four independent rules (ECI-16) |
| **Cross-stream dependency at entry** | **constrained** — content copy only | values may cross; status may not (ECI-13/14) |
| **Cross-stream ordering** | **not required** | a copied turn is one fact recorded twice, not a comparison (ECI-35) |

**ODB-63 (Inferred).** At entry the three options are **identical in requirement**: each must record a
lifecycle-side fact carrying `provisional` with formation's values. They diverge only **after** entry
(ECI-31/32). Selecting among them is therefore a decision about **transitions**, not about entry.

### L1 — State-naming

The stream records each lifecycle event with the state it enters or leaves, from the §11 vocabulary.

- **Expresses:** the progression directly; each of ODB-40's five facts maps to an event.
- **Cost:** the event vocabulary duplicates §5's table, so the two must stay in step; the
  `active`/`established` mismatch becomes an immediate, unavoidable naming choice; rumour's terminal state needs
  an event that names a state with no successor.
- **Does not decide:** who may act (ND-1), when (OD-g), cause or evidence shape (D-3b).

### L2 — Transition-only

The stream records the **relation** (entered, superseded-by, corrected-by, retracted), and states are derived
by folding the relations.

- **Expresses:** K-R12's own principle — "status is projected, never a mutable field" (`:321`) — most directly.
  Terminal `provisional` becomes "entered, never related onward", which needs no special case.
- **Cost:** the fold must be **total**; an unclassifiable relation leaves a record's state undefined, which is a
  P-2 definitional risk. Entry also needs an **entry relation or a closure fact**, otherwise the fold cannot tell
  a rumour from a record whose successor was never recorded (`...ENTRY-CONDITION...` ECI-25, ECI-27). A
  transition-only vocabulary also has no natural place for *unverified* as a property — which is precisely where
  C′ returns with force (§10).
- **Does not decide:** the same three.

### L3 — Act-only

The stream records the **grounding act** (the formation evidence, the promotion evidence kind, the correction,
the retraction), and both state and transition are derived.

- **Expresses:** the separation the owner has insisted on — cause is recorded as cause, not as status; and the
  per-kind exit rules are the natural content.
- **Cost:** the derivation rules become the real specification, and they must be total across every kind,
  including the kind that never exits. There is also an inversion risk: act-only can make cause authoritative
  *over* status, which is the separation failure in the opposite direction.
- **Does not decide:** the same three.

**ODB-46 (Inferred).** The options differ in **where the specification's weight sits** (state names, relations,
or acts), not in what the progression is. All three can express ODB-40's five facts; they differ in what they
make easy, and in which omission becomes silent.

### 5.2 The three discriminators

**ODB-64 (Inferred — L1's question).** *What does explicit state representation cost?* Two known costs: the
`active`/`established` reconciliation becomes an immediate naming choice (ODB-04), and each §11-adjacent concept
outside the five names needs a home — **forgetting first** (CLV-38). What L1 buys: no fold, no closure argument
for the rumour, and a state that cannot be lost silently (ECI-21/ECI-37).

**ODB-65 (Inferred — L2's question).** *Can the fold and its closure guarantee be made complete **without losing
the reason for a transition**?* The reason is not carried by a relation (CLV-27), so retaining it means
referencing a ground inside the transition — which is where the cause could drift into authorising the status
(CLV-28). L2 also inherits the closure obligation (ECI-25) and the phantom-provisional mode (ECI-27), and under
§8.1 read literally, a wrong default (EKR-22).

**ODB-66 (Inferred — L3's question).** *Can derivation stay deterministic without the grounding act becoming an
unintended authority over status?* Determinism needs a total derivation across every kind, including the kind
that never exits; and the inversion risk is the separation failure in the opposite direction (CLV-26).

**ODB-67 (Inferred — none is preferred by the narrowing).** L2 fits K-R12's projection principle most naturally,
but that is a **fit**, not an argument: its closure burden is substantive (ODB-65). L1 and L3 carry their own
costs (ODB-64/66). No option is selected **in this section**, and none is excluded on evidence. **Updated
2026-11-20:** the owner selected **L2** (ODB-72); L1 and L3 were not selected, and were not shown unsound.

**ODB-47 (Resolved 2026-11-20 — originally Unestablished).** Which option the owner intends. Canon does not choose; it constrains. **Answered: L2** (ODB-72).

**ODB-70 (Observed — where the criteria are worked out).** The three questions are assessed per sub-criterion in
`...OD-B-OPTION-EVALUATION...` (OBE-08/13/17). It reports no dominance (OBE-23) and selects nothing; the residual
choice is OBE-29, refined in **form** by OBE-37. The owner's selection is recorded separately
(`MEMORY-09-OWNER-DECISION-OD-B-L2-2026-11-20.md`), with the evaluation's post-decision note at OBE-40…OBE-47.

**ODB-71 (Observed — the two conditional viability questions).** L2's closure question and L3's validity boundary
are examined in `...CLOSURE-AND-VALIDITY-EVIDENCE...`: a closure declaration is **required in effect by I15/I18 and
absent from the text** (CEV-18), and the recording-versus-conferring distinction is **stated by I17 and
unenforceable today** (CEV-32). Both remain open; neither is a mechanism.

### 5.3 A separate open question, outside OD-b

**ODB-68 (Observed — retained, not absorbed).** §8.1 `:404` has formation writing `status=active`, against §11's
entry `provisional`. This is **not** an OD-b discriminator and **not** a reason to prefer an option. It is a
per-kind semantic question with two parts that must not merge — does formation *always* enter `provisional`, and
can an exit be satisfied *at* formation — plus a storage-level problem: §8.1 writes a field K-R12 says is a
projection (`...ENTRY-UNIFORMITY...` EKR-19).

**ODB-69 (Inferred — why it stays separate).** A fold-based selection **inherits** the question as a downstream
obligation (EKR-22/24) without being entitled to settle it. Keeping it outside OD-b is what stops a worked case
from amending §11's diagram by inference.

---

### 5.4 The owner's selection, recorded

**ODB-72 (Observed — owner, 2026-11-20).** **OD-b = L2 — transition-only.** Lifecycle history is represented
through transitions; current lifecycle state is derived as a **projection** of that authoritative history, and the
derived state is not a competing authoritative status field. **L1 and L3 were not selected.** The decision, the
owner's rationale and the owner's failure-risk ordering are recorded in
`MEMORY-09-OWNER-DECISION-OD-B-L2-2026-11-20.md` (ODL-01…ODL-23).

**ODB-73 (Observed — why the rationale is not restated here).** The rationale and the risk ordering have one
recorded home, to prevent a second copy from drifting from it. This brief points to them (`ODL-10…ODL-23`) and
records only what the selection means for its own analysis.

**ODB-74 (Observed — what stands unchanged).** §5.1's constrained-out table, §5.2's three discriminators and
§5.3's separation of §8.1 are **unchanged**. They explain why the decision was a choice among three viable
stances, not which stance to take; a recorded choice does not convert them into a justification of it.

**ODB-75 (Inferred — L2's costs are now obligations).** ODB-65's closure and reason questions, and EKR-22's §8.1
interaction, are **inherited obligations of the selected option**, not objections that lapsed when it was chosen.
Concretely: closure (CEV-18 — required in effect, absent from the text), transition ground (CLV-27/CLV-28), and
deterministic reconstruction that preserves uncertainty instead of manufacturing a state. None is designed, and
none is authorised to be designed here.

**ODB-76 (Inferred — L1 and L3 are not falsified).** Non-selection is not a finding against them: OBE-23 reported
no dominance, and OBE-27 recorded L1's viability as currently unthreatened by an open question. Their rows in §6
remain valid as recorded, and remain readable if L2's obligations prove unsatisfiable.

**ODB-77 (Observed — nothing is authorised).** The selection authorises a **recorded decision only** (ODB-29). The
next legitimate step is still a separate **responsibilities-only** proposal (ODB-30) — not fields, not events, not
a closure marker.

**ODB-78 (Observed — the matrix is pre-decision).** §6's matrix is the comparison that preceded the choice and is
retained as written. Its L2 row entries "no onward relation **+ closure**" and "**not carried** — needs a ground
reference" now read as the inherited obligations of ODB-75.

---

## 6. Consequence matrix

**ODB-48 (Inferred — coverage of §11's structural facts).**

| §11 fact | L1 | L2 | L3 |
|---|---|---|---|
| **Entry at formation** | **constrained — identical** | **constrained — identical** | **constrained — identical** |
| Exit conditional on evidence kind | event names the target state | relation names the change | the act's kind supplies it |
| Terminal `provisional` (rumour) | state named, no successor | no onward relation **+ closure** | no exit act **+ closure** |
| Reason for the transition | payload choice | **not carried** — needs a ground reference | inherent, must not invert |
| Closure (complete vs lost successor) | **not required** | **required** | **required** |
| Retraction retains acquisitions (`:616`) | target state only | relation only | act only |
| `active` vs `established` naming | **must be chosen immediately** | deferred to the fold | deferred to derivation |
| Verification attribute (C′) | **no home** | **no home** | **no home** |

**ODB-49 (Inferred).** The last row is the finding that matters for sequencing: **no option represents the
verification dimension**. Whichever is chosen, C′ remains open, and choosing an option without noticing this
would bury it.

---

## 7. What formation producing `provisional` does **not** establish

**ODB-57 (Inferred — owner caution, adopted).** That formation creates a `provisional` record says nothing
about the **cause, evidence, occurrence or authority** of that creation event. Those remain governed by prior
decisions and open questions:

| Responsibility | Governing position |
|---|---|
| **Cause** | must retain the actual cause; no fabricated cause (D-3b) |
| **Evidence** | `evidence_span` declared non-nullable in four schemas, defined in none; D-3b permits statement-less events, which the nullability contradicts |
| **Occurrence** | A2 per-stream order; a turn-less event's cross-stream position is **undefined** (OD-g) |
| **Authority** | **ND-1** open and widened; who may cause a lifecycle change is not established |
| **Identity** | unstable `len()`-derived ids and second-granularity `NewMemoryID()` collisions are independent defects, not solved by OD-b |

**ODB-58 (Inferred).** A lifecycle structure that assumes formation's creation event has an authority, a
position, or an evidence span already settled would import three unowned decisions. §5's options therefore
constrain only *what is recorded about the progression*, never *what justifies it*.

---

## 8. What OD-b must not absorb

**ODB-27 (Source-derived — owner boundary).**

| Issue | Position after this revision |
|---|---|
| `provisional` | **lifecycle entry state** (ODB-08) — owner confirmation pending, not a design choice |
| `active` / `established` | vocabulary mismatch, **unresolved** (ODB-04) |
| `forgotten` | required by DN-4, **absent from all three enumerations** (ODB-07) |
| Verification attribute | **C′**, separate design question (§10) |
| Record-promotion vs authority-promotion | different operations and authorities (ODB-50; `...CHECK...` ODBC-19/20) |
| Cross-stream chronology | **OD-g**, unresolved (ODA-06/ODA-28) |
| Domain identity | missing-structure finding, **outside this experiment's state-machine scope** (ODA-20–ODA-26) |

**ODB-50 (Inferred).** "Record-promotion" (§11: a record leaves `provisional` by evidence kind) and
"authority-promotion" (v9 rules 24–38: an actor confers world/scenario authority) are **different operations**
with different authorities. The seam is §11's "world fact — authority or corroboration" (`:612`). OD-b must not
merge them; L1–L3 all touch this seam and none may settle it.

**ODB-28 (Inferred).** A decision that answered any row of this table in passing would be a decision the owner
did not make.

---

## 9. What the owner's answer authorises

**ODB-29 (Source-derived).** An OD-b answer authorises **a recorded decision and nothing else** — no schema
change, no implementation, no tests, no canonical edit, no registry change, no ledger entry.

**ODB-30 (Inferred).** After OD-b, the next legitimate step is a **separate proposal for the lifecycle
structure's responsibilities** — still not fields — which the owner must authorise before any schema exists.

---

## 10. C′ — a separate design question, not a third interpretation

**ODB-59 (Inferred).** C′ asks: **should verification condition be represented independently of lifecycle
status?** Evidence that it exists as a distinct concern: C4's retrieval channel keys on "provisional +
`via=told`" (`:446`); assembly forbids presenting provisional as verified (`:283`); conflict materiality
compares "provisional vs established" (`:459`); `confidence_in_source` already exists on the acquisition side.

**ODB-60 (Inferred — its relationship to OD-b).** C′ is **downstream of OD-b but not settled by it**, and
ODB-48's last row shows none of L1–L3 gives it a home. It should be decided after OD-b, in its own record, with
its own boundary. It must not be used to reintroduce a third OD-b reading.

---

## 11. Experimental boundary

**ODB-61 (Observed — owner direction, recorded).** **Howl's dynamic state machine remains Howl's
responsibility.** MEMORY-09 establishes the event and memory representation Howl may eventually use; it does
**not** design the machinery that switches worlds, scenarios and characters. r33/v9 hold the active world and
scenario constant, and v9 (`:228–229`) places switching outside the experiment. This matches the OD-a record's
ODA-18/ODA-23 and authorises no state-machine work.

---

## 12. Why this brief still contains no field list

**ODB-33 (Inferred).** Any concrete field set would encode an answer to OD-b, to `active`/`established`, to
`forgotten`, to C′ and to OD-g simultaneously, presenting five unowned decisions as implementation detail.

---

## 13. Limits

**ODB-34 (Observed).** Term-based search over `docs/research` markdown plus one Go types file. Not read: all
711 documents; `internal/backend/memory/*.go` beyond the types file; §11's authority level beyond r33 `:3`'s
"design record, not implementation" (r33 is untracked, so no commit-level chronology exists).**§5–§11 provenance,
vocabulary coexistence and the L1–L3 consequence tests are now examined separately** in
`MEMORY-09-SECTION-5-11-PROVENANCE-AND-L1-L3-COMPARISON-2026-11-19.md` — that record makes no selection —
and the **entry condition under OD-a's separate structures** is examined in
`MEMORY-09-ENTRY-CONDITION-AND-CROSS-STRUCTURE-INTEGRATION-2026-11-19.md` (ECI-31/ECI-32: at entry the three
options converge; they diverge only post-entry).

**ODB-35 (Resolved in part — 2026-11-19).** The r3/r4 proposed-answer check is **closed**: no distinct meaning
(ODBC-01/02). Still unestablished: whether the Go candidate pipeline is the ancestry of the §5 value; whether
any unratified requirement needs `provisional` to mean review.

**ODB-36 (Observed).** No canonical file was edited; the 118 requirements are unchanged; no canonical capture
of any decision has occurred.

**ODB-79 (Observed — 2026-11-20).** No canonical file was edited to record the L2 selection either; the 118
requirements are unchanged; the four canon hashes are unchanged; 0 tests were run.

**End of brief (revision 4).**
