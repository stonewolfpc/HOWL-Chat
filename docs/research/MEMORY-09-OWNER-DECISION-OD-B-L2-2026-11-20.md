# MEMORY-09 — Owner Decision OD-b: L2 Transition-Only Lifecycle Representation

**Status:** owner-decision record. **OD-b settled: L2 — transition-only.** Documentation and read-only analysis
only. **No schema change, event or relation design, closure-mechanism design, implementation, test, probe,
benchmark, model call, canonical edit, registry change or ledger entry is authorised by this record.**

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Decision | **OD-b = L2 — transition-only.** Lifecycle history is represented through transitions; current state is a **derived projection** of that history |
| Decision scope | **Lifecycle representation only** — not a schema, not a field list, not a mechanism |
| Authority | Owner, stated in session 2026-11-20 and recorded here; as with every prior confirmation, **no persisted signature exists** |
| Options not selected | **L1** (state-naming) and **L3** (act-only) — not selected, and **not** shown unsound |
| Canonical edits | **None.** r33, v9, the owner ledger and the conformance map are unchanged |
| Brief revision | OD-b owner decision brief **revision 4** (§5.4) |
| Evaluation | `...OD-B-OPTION-EVALUATION...` §11 (OBE-40…OBE-47) |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. The decision as stated

**ODL-01 (Source-derived — owner).**

> "L2 is selected: lifecycle history is represented through transitions, and current lifecycle state is derived
> as a projection of that authoritative transition history."

**ODL-02 (Source-derived — owner).**

> "The underlying lifecycle history is authoritative. The derived current state is a projection, not a competing
> authoritative status field."

**ODL-03 (Observed — provenance of the text above).** These are the owner's statements as supplied on
2026-11-20 and transcribed here. No persisted signature exists, and no canonical file was edited to record them
(ODL-34). They are quoted because they are the decision; the analysis that follows is labelled separately.

---

## 2. What is decided

**ODL-04 (Source-derived).** **The lifecycle stream records transitions.** The authoritative content of the
lifecycle structure is the history of changes, not a set of stored states.

**ODL-05 (Source-derived).** **Current lifecycle state is derived.** It is a projection computed from that
history — the stance K-R12 states ("status is projected, never a mutable field", r33 `:321`).

**ODL-06 (Source-derived).** **The derived state is not a competing authority.** A projection that disagreed with
the history would be the defect, not an alternative reading. There is one authoritative source for status.

**ODL-07 (Observed — what the selection does *not* add at entry).** Entry was already constrained and identical
across all three options (ODB-63; ECI-31/ECI-32): a lifecycle-side fact carrying `provisional` with formation's
values, because acquisition-derived status is forbidden (ECI-16). Selecting L2 changes nothing at entry; it
decides **how transitions after entry are represented**.

**ODL-08 (Inferred — status of L1 and L3).** L2 was selected; L1 and L3 were **not**. Non-selection is not a
finding against them: the evaluation reported no dominance (OBE-23), and L1's viability was the one not gated on
an open question (OBE-27). Their analysis stands as recorded, and remains available if L2's obligations (ODL-26…
ODL-31) prove unsatisfiable.

**ODL-09 (Observed).** This record fixes a **representational stance**. It names no field, no event type, no
relation, no marker and no algorithm (ODB-33; ODA-13).

---

## 3. The owner's rationale — preserved

Each clause is recorded in the owner's terms. It is not a derivation from the comparative findings (§4.21).

**ODL-10 (Source-derived — owner).** L2 provides **better control over the memory lifecycle**.

**ODL-11 (Source-derived — owner).** Maintaining **one authoritative history and deriving current state from it
reduces the risk of competing representations drifting apart**.

**ODL-12 (Source-derived — owner).** Reconstructing state from history gives the system **better opportunities to
reason about available information and make decisions along the way**.

**ODL-13 (Source-derived — owner).** The architecture should **accommodate unavoidable uncertainty without
unnecessarily halting unrelated system operations**.

**ODL-14 (Source-derived — owner).** **Uncertainty should constrain the affected decision** rather than
automatically disabling the broader memory system.

**ODL-15 (Inferred — how the rationale is to be read).** It is recorded as **owner-stated rationale** for a
representation choice, in the same sense as ODA-15…ODA-19: it explains why the choice was made, and it is **not**
a new normative requirement on MEMORY-09 (CEV-38; ODB-29).

**ODL-16 (Inferred — the boundary on ODL-13/ODL-14).** "Uncertainty constrains the affected decision" is **not**
"uncertainty may be resolved by assumption." Neither clause licenses inferring information from absence,
fabricating evidence, or treating an incomplete history as complete. **I15, D-3b and I17 remain exactly as
ratified**; the rationale addresses **availability** of the system, not **permissiveness** of the representation.
This is the same boundary the owner stated explicitly (ODL-20).

**ODL-17 (Inferred — why this is not an elegance preference).** The rationale names control, drift,
reconstructability and availability. Implementation elegance is not among its grounds, and the record does not
substitute one: reducing the decision to "transition-only is cleaner" would misstate the owner's reasons and would
also overstate the case, since L2's costs are real and unresolved (ODL-26…ODL-28).

---

## 4. The owner's failure-risk preference — recorded ordering

**ODL-18 (Source-derived — owner's risk judgment).**

| Rank | Failure mode | Owner's assessment | Comparative source |
|---|---|---|---|
| 1 | **Unwarranted state discovered during audit** — a state the record's evidence cannot support (L1's mode) | **Least problematic** of the three | OBE-11 |
| 2 | **Plausible but incorrect reconstruction from incomplete history** — truncation folding to a confident `provisional` (L2's mode) | **More problematic** than an audit-discovered error, **less** problematic than authority inversion | OBE-15 |
| 3 | **Recording an act that effectively confers status without independently established authority** (L3's mode) | **Least acceptable** and most problematic | OBE-19 |

**ODL-19 (Source-derived — owner's accompanying statement).** The owner considers some plausible-but-wrong
outcomes **nearly unavoidable** in a system interpreting imperfect conversational information. This is the reason
the middle mode is placed second rather than last.

**ODL-20 (Source-derived — owner's explicit boundary).** The owner is **not** accepting careless reconstruction
and is **not** relaxing **I15**. The distinction is that unavoidable uncertainty should be handled without making
the entire system unnecessarily unavailable.

**ODL-21 (Inferred — what kind of statement the ordering is).** It is a **risk judgment about tolerability**: which
residual failure this system can live with, given its likelihood and the difficulty of detecting and recovering
from it — the form the owner refined at OBE-37/CEV-39. It is **consistent with**, and does not revise, OBE-23's
finding that no option dominates: the owner chose which failure to accept, not which option is superior on the
sub-criteria.

**ODL-22 (Observed — how the ordering relates to the evidence).** The three modes are the ones the comparative
record identified, surfacing at **audit**, **read** and **write** respectively (OBE-22), and CEV-40 observed that
their likelihood and recoverability are **not equal**. The ordering ranks **tolerability**, not frequency or
repairability; it therefore does not assert that a truncated history can be reconstructed from what is missing.
That asymmetry is precisely why closure remains an outstanding obligation (ODL-26).

**ODL-23 (Observed — what the ranking does not do).** It does **not**:

- permit inferring information from absence, fabricating evidence, or treating incomplete history as complete;
- exempt any element from **I15** or **I17**/I7;
- assert that L2 passes every sub-criterion, or that its costs are resolved by having been accepted;
- express a preference about `active`/`established`, §8.1, C′, OD-g, ND-1 or DN-6a/DN-6b.

---

## 5. Constraints retained — unchanged by the selection

**ODL-24 (Source-derived — owner's list, verified against canon).** L2 selection does not remove or weaken any
existing requirement:

| Authority | What the selection relies on | What it does not alter |
|---|---|---|
| **I15** | Absence of information must never be treated as information about absence | unchanged; remains binding on the fold, on closure and on validity |
| **I17 / I7** | Recording or retrieving an act does not itself confer authority | unchanged; the reason a transition carries is a ground, never a warrant |
| **K-R11** | Acquisition-event folding stays confined to acquisition events | unchanged; the lifecycle fold is a separate fold |
| **K-R12** | Lifecycle status is a projection of the lifecycle stream, not a mutable authoritative field | unchanged — the selection is *compatible* with it, which is not proof of superiority (owner's earlier caution) |
| **K-R14** | Reconstruction respects the applicable stream's own event history | unchanged; no cross-stream reconstruction is implied |
| **D-3b** | Statement-less events retain their actual cause and an independently identifiable occurrence | unchanged; no fabricated evidence, no inferred statement position |
| **A2** | Each authoritative stream keeps its own deterministic occurrence ordering | unchanged; no shared global chronology is implied |
| **OD-a** | Acquisition events and lifecycle/status events remain separate structures | unchanged; the selection is downstream of it |
| **DN-4** | Forgetting stays distinct from record lifecycle status and retraction | unchanged; L2 has no vocabulary that would absorb it, and none is added |
| **Q-B** | Application status changes do not independently create character awareness | unchanged; the projection creates no awareness |

**ODL-25 (Observed).** Nothing above is reinterpreted or amended by this record. Listing an authority here
records **compatibility**, and does not convert any of them into a new requirement (CEV-38).

---

## 6. Obligations carried forward — not discharged by selecting L2

These are the costs of the selected option. They are real, unowned, and unauthorised to design (ODB-65; OBE-14).

### A. Completeness and closure

**ODL-26 (Inferred).** The system must not mistake a **truncated** lifecycle history for a **complete** one. A
missing successor cannot establish that a record is terminally `provisional`. **I15 remains binding.** The need for
completeness/closure semantics is an **architectural obligation**, but **no closure mechanism has been selected or
authorised** — and CEV-18 found such a declaration *required in effect by I15/I18 and absent from the text*
(CEV-02…CEV-06 set out the acquisition stream's two-part precedent). Whether a stream can express "no further
entries by design" remains **Unestablished** (OBE-36).

### B. Transition-ground preservation

**ODL-27 (Inferred).** A transition must not lose the **reason or grounding** for the transition. Preserving a
ground must not turn that ground into an **authority source**, and must not allow the act of recording it to
confer status (CLV-27/CLV-28; I17; OBE-14.1). L2 may decline to carry the ground, at the cost that "why" becomes
unrecoverable — that is a decision still to be made, not one this selection makes.

### C. Deterministic reconstruction

**ODL-28 (Inferred).** The **same qualifying lifecycle event history must yield the same derived projection**.
Where the available history cannot establish a trustworthy result, the system must **preserve that uncertainty**
rather than silently manufacturing a definitive state (K-R2's integrity framing; OBE-15's silent-failure mode).

### D. `active` / `established`

**ODL-29 (Observed).** The vocabulary reconciliation remains **open**. L2 **defers** the naming to the fold; it
does not eliminate it — the projection still emits one name, so K-R12's vocabulary question must be answered
either way (OBE-03; ODB-04).

### E. The §8.1 discrepancy

**ODL-30 (Observed).** The entry-uniformity and §8.1 record remains **independent**. L2 selection must not be used
to resolve **P-EVAL** (EKR-11), the `active` versus `provisional` discrepancy, or the mutable-field conflict
(EKR-19). L2 **inherits** the §8.1 interaction as a downstream obligation (EKR-22/EKR-24) without being entitled
to settle it.

**ODL-31 (Inferred).** A–E are **obligations, not authorisations**. Each requires its own recorded decision before
any design work, and none of them is satisfied by the existence of this record.

---

## 7. What remains unresolved after this decision

**ODL-32 (Observed).** None of the following is answered, narrowed or assumed by the L2 selection:

| Item | State after this decision |
|---|---|
| **OD-g** — cross-stream comparability and ordering | **OPEN**, unchanged (ODA-06; ODA-28) |
| **ND-1** — lifecycle status and authority semantics | **OPEN**, widened |
| **DN-6a / DN-6b** — authority to cause lifecycle changes | **OPEN** |
| **C′** — verification attribute placement | **OPEN**, independent of OD-b (ODB-59/ODB-60) |
| **P-EVAL** — timing of lifecycle exit-condition evaluation | **OPEN** (EKR-11) |
| **Domain identity** — world/scenario/character identity in the record model | **OPEN**, missing-structure finding (ODA-20/ODA-21), outside the state-machine scope |
| **Cause, evidence, occurrence, stable event identity** | **unchanged**; governed by D-3b, A2, OD-g and the separate identity defects (ODA-07; ODB-57) |
| **Lifecycle stream schema and closure representation** | **not drafted, not authorised** — blocked on the obligations above |
| **Exact transition-ground representation** | **OPEN** (ODL-27) |

**ODL-33 (Observed).** `active`/`established` (ODL-29) and the §8.1 question (ODL-30) belong in this table too;
they are listed under §6 because L2 inherits them rather than merely coexisting with them.

---

## 8. What this decision does not authorise

**ODL-34 (Source-derived — owner's boundary, and ODB-29).** Authorised: **a recorded decision, and documentation
of it**. Not authorised:

- schema changes, field lists, or any normative statement of structure;
- event or relation design;
- closure, completeness or fail-closed mechanism design;
- implementation or code changes of any kind;
- tests, experiments, probes, benchmarks or model calls;
- registry or requirement-ledger changes;
- changes to frozen canonical artifacts, r33, v9, or any other normative source;
- reopening or revising settled owner decisions;
- designing Howl's world/scenario/character **state machine**.

**ODL-35 (Observed).** The Howl state machine remains **outside** MEMORY-09 experimental scope (ODB-61;
ODA-18/ODA-23). No transition, activation or lookup mechanism is specified here (ODA-31).

**ODL-36 (Observed — capture status).** As with OD-a, the decision lives in this dated, untracked record until the
owner authorises canonical capture. No ledger entry exists, and **O-16 is not authorised** by this record.

---

## 9. Traceability

**ODL-37 (Observed).**

| What | Where |
|---|---|
| The decision, rationale and risk ordering | this record, §1–§4 |
| Brief now showing L2 selected, L1/L3 not selected | `...OD-B-OWNER-DECISION-BRIEF...` **revision 4**, §5.4 (ODB-72…ODB-78) |
| Comparative findings, and the selected risk preference distinguished from them | `...OD-B-OPTION-EVALUATION...` §11 (OBE-40…OBE-47) |
| Silent-failure modes and their surfacing points | OBE-11, OBE-15, OBE-19, OBE-21, OBE-22 |
| Closure evidence (Q-1) and validity evidence (Q-2) | `...CLOSURE-AND-VALIDITY-EVIDENCE...` CEV-02…CEV-18, CEV-21…CEV-34 |
| Entry constrained and identical across options | ODB-63; ECI-25/ECI-27/ECI-31/ECI-32 |
| Reason and cause separation | CLV-26/CLV-27/CLV-28 |
| §8.1 and P-EVAL kept separate | `...ENTRY-UNIFORMITY...` EKR-11, EKR-19, EKR-22, EKR-24 |
| Separate structures | `...OWNER-DECISION-OD-A...` ODA-01…ODA-14 |

---

## 10. Limits, tensions and record status

**ODL-38 (Observed — method limits).** This record transcribes a decision and restates constraints and obligations
already documented elsewhere. No canon was re-read beyond the items cited; no implementation was inspected; no
test, probe or experiment was run. No new textual evidence is offered for any open question.

**ODL-39 (Observed — tension 1, recorded not resolved).** L2's accepted residual failure — a truncated tail folding
to a confident `provisional` — is a **silent** failure of exactly the kind K-R17 rejects at write in the acquisition
stream: *"a silent no-op hides formation bugs"* (CEV-02/CEV-03). The two streams differ in semantics, and CEV-05
notes that the acquisition rule needed **two** parts (reject the impossible transition; fix what absence means).
Whether the lifecycle stream should adopt a comparable rejection-at-write discipline is **not established** — and
answering it would be closure-mechanism design, which is not authorised (ODL-34, ODL-26). Recorded so the
seam is not mistaken for a settled precedent.

**ODL-40 (Observed — tension 2, recorded not resolved).** CEV-40 concluded that a truncated history **cannot be
reconstructed from what is missing**, while the owner's ordering places reconstruction failure second rather than
last. Both statements can hold: the ordering is about which residual risk is **tolerable** given that closure is an
outstanding obligation (ODL-26), not about the loss being repairable. Recorded so the ordering is not read as a
recovery claim.

**ODL-41 (Observed — no other conflict found).** The decision statement, rationale and ranking are internally
consistent and consistent with the retained authorities in §5. No ambiguity in the decision required resolution by
inference, and none was resolved that way.

**ODL-42 (Unestablished).** Whether a closure declaration is admissible in the existing architecture (OBE-36);
whether the ordering would be revisited if closure proves inexpressible. The latter is the owner's call and is not
inferred here.

**ODL-43 (Observed — record status).** **OD-b is SETTLED — L2, representation only.** OD-a settled (Option B).
A2, D-3b and DN-4 unchanged. DN-1 open and widened. OD-g open. C′ open. The 118 ratified requirements unchanged.
0 tests. No canonical, registry, schema or code file was edited.

**End of record.**
