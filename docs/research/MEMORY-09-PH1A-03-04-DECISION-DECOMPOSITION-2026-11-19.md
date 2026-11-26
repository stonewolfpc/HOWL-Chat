# MEMORY-09 — PH1A-03/04 Decision Decomposition

**Status:** research report. **READ-ONLY.** No architecture selected, no hybrid designed, no
implementation, no test, no probe, no model call, no corpus, no canonical edit, no registry change, no
Git operation. The decision is **not** put to the owner here.

| Field | Value |
|---|---|
| Created | **2026-11-19** (local, CST) |
| Question | Is PH1A-03/04 one indivisible decision or a bundle of separable mechanisms? |
| Baseline | r33 `dfb28b6003ca356f` · ledger `e5008020ff68bd21` (O-14 `:205`, O-6a `:424`) · conformance `b4ba33ccbfe8f8c1` · Phase 1A CO `82173a6f78ac84c4` |
| Tests executed | **0** |
| Files changed | One new file. Untracked, unregistered, unfrozen |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. The two items, verbatim

**DC-01 (Observed).** The authoritative definitions are rows in the investigation's decision table
(`MEMORY-09-PHASE-1A-CONCEPTUAL-OBJECT-INVESTIGATION.md:918–919`):

| ID | Question as recorded | Type recorded | Tier recorded |
|---|---|---|---|
| **PH1A-03** | "Does acquisition event-sourcing survive v9?" | Owner decision / source reading | T2 vs T3 |
| **PH1A-04** | "Does K-R12's memory-lifecycle stream survive?" | Owner decision | T2 |

Both name **K-R** mechanisms (T2) and ask whether they survive a **v9** model (T3). Neither asks which
of two implementations to use.

**DC-02 (Observed).** The investigation's substantive treatment is §7.2 (`:535–551`), which already
separates the two:

- acquisition event-sourcing "remains a T2 obligation **wherever T1 has not displaced it**", and T1
  displaces it "only for the private-disclosure case (O-1 requires a separate recipient record)";
- the lifecycle stream (K-R12) "is **not displaced by anything**" and is unimplemented in both.

**DC-03 (Inferred).** That paragraph contains the whole decomposition. The two items are not parallel
questions of equal weight: one is a conflict whose scope T1 has already narrowed, the other is an
obligation nothing contradicts. Treating them as one decision is what made them look symmetrical.

---

## 2. What each option is intended to satisfy

**DC-04 (Source-derived).** Tier-1 decisions are silent on storage. D-01 (containment), O-1 (exclusive
ownership), O-3 (disclosure persistence), O-5 (belief independence), O-10/O-10b (authority changes),
O-13 (reach), O-14 (awareness) constrain *what must be true*, never *how it is held*. No T1 decision
names a stream, a field, a projection or an event.

**DC-05 (Source-derived).** The T2 mechanisms in play:

| Mechanism | What it requires | Kind |
|---|---|---|
| **K-R1** | an edge exists only as the fold of that pair's acquisition events | storage |
| **K-R2** | a projection that cannot be rebuilt is an integrity failure | storage |
| **K-R3** | revocation is a new event, never a deleted row | storage + behaviour |
| **K-R11** | the fold depends only on acquisition-family events | storage |
| **K-R12** | status/scope/class are projections of a lifecycle stream, not fields | storage |
| **K-R13** | the retraction join is read-time and never stored | behaviour (now moot, §4) |
| **K-R14** | each stream is rebuilt and checked independently | storage |
| **K-R15** | events ordered by `(turn, seq)` | storage + behaviour |
| **K-R17** | a revoke with no open episode is rejected | behaviour |
| **K-R19** | a forgotten pair has **no edge** | behaviour |
| **K-R20** | `revoked_turn` is derived, not a field | storage |

| Requirement | Status under a stream model | Under a record-with-field model |
|---|---|---|
| **A2a / E1a / E1b / H10** | satisfiable | satisfiable |
| **O-1 / rule 12** | satisfiable if records carry an owner | satisfied already |
| **F1 / F1b** | satisfied (rebuild from events) | **F1** satisfied by accounts; **F1b** has no projection to compare — `Unestablished` |
| **I4 / I11 / K-R3 / K-R19** | satisfied | **not satisfied**: no forget event, no absence-of-edge state |

**DC-06 (Calculated).** Of the eleven T2 mechanisms above: **six are purely storage** (K-R1, K-R2,
K-R11, K-R12, K-R14, K-R20), **three are behavioural** (K-R3, K-R17, K-R19), **one spans both**
(K-R15, whose ordering is storage and whose determinism is observable), and **one is now moot**
(K-R13, because O-14 removed the notice). 6 + 3 + 1 + 1 = 11 ✓. A decision that bundles them decides
storage and behaviour together, which is precisely the conflation this pass is asked to undo.

---

## 3. Storage representation versus behavioural semantics

**DC-07 (Inferred).** The two are separable in principle because a behaviour can be produced without the
storage shape that usually carries it, and a storage shape can exist without any behaviour depending on
it:

| | Behaviour observable to a character or operator | Storage shape only |
|---|---|---|
| **Retraction** | source status changes; accounts persist (E1a, E1b, H10, A2a) | whether status is a field or a projection (K-R12) |
| **Forgetting** | the character no longer knows (K-R3, K-R19) | whether it is an event or a deleted row |
| **Ordering** | same-turn resolution is deterministic (S38) | `seq` |
| **Recoverability** | an operator can see how knowledge changed | rebuildability (K-R2, K-R14) |

**DC-08 (Source-derived).** Only the left column is constrained by T1. The right column is constrained
by T2, which governs by tier order and does not need a T1 justification to apply.

---

## 4. The seven independencies, tested one by one

### 4.1 Append-only event history — **independently decidable**

**DC-09 (Source-derived).** Governed by K-R1/K-R11/I11. The question is whether an acquisition record
must exist at all. It does not depend on how status is represented: a store can keep acquisition events
while status remains a field.

**DC-10 (Inferred).** It is also the one place where the alleged T2/T3 conflict may not exist. v9 §8's
prohibition reads: the recipient's origin reference is "a descriptor … not a resolvable identifier. The
resolvable link exists only on the source side. Each recipient record has its own identity and
acquisition provenance; it is **never a view, alias or projection**." The antecedent of "it" is the
source record. Read with that antecedent, the sentence forbids B's record being a view *of A's record*;
it does not forbid B's record being derived from B's acquisition events. Under that reading, **rule 12
and K-R1 govern different relations and do not conflict** — so PH1A-03 is a *reading* question, not a
divergence.

**DC-10b (Unestablished).** The competing reading — "projection" with no antecedent, forbidding any
derivation — would make it a case-2 divergence. Nothing in v9 resolves which reading is intended. **This
is the single textual point on which PH1A-03 turns.**

### 4.2 Status representation — **independently decidable, and probably not an owner decision**

**DC-11 (Source-derived).** K-R12 requires status to be a projection. v9 rules 20, 37 and 42 describe the
*effects* of lifecycle operations (history preserved; revocation and supersession explicit; three
operations with three meanings) and say nothing about whether status may be a stored field.

**DC-12 (Inferred).** Silence is not contradiction — O-6a §2 states it directly. O-6a's three cases are
(v9 vs T1), (v9 vs applicable r32), (r32 silent, v9 speaks). This is a **fourth combination**: r32 speaks,
v9 is silent. No case is needed, because plain tier precedence already decides it — T2 governs where it
is in scope and does not conflict with T1. **PH1A-04 therefore needs no owner decision; it needs only a
recorded finding that v9 is silent, and would need a decision only if the owner chose to diverge from
K-R12.**

### 4.3 Account lifecycle — **independently decidable**

**DC-13 (Source-derived).** v9 rules 53–55 govern belief lifecycle; the account lifecycle is fixed by
O-1/O-3 and H10, all T1. Nothing in either stream mechanism changes when an account is created, copied or
retained. Independent of 4.1, 4.2 and 4.4.

### 4.4 Forgetting — **independently decidable, and behavioural**

**DC-14 (Observed).** Neither implementation can forget. The engine has no concept; the harness's
`revoke()` sets record status and still emits an edge, contradicting K-R19 (§7.3). Forgetting is the only
mechanism whose *observable* behaviour differs between the two models: in a stream model "forgot" is the
absence of an edge with history intact; in a record model it is indistinguishable from retraction
(CO-07).

**DC-15 (Source-derived).** No conformance row tests forgetting. It is a T2 obligation with no probe.

### 4.5 Retraction — **independently decidable, and now minimal**

**DC-16 (Source-derived).** Retraction's content is fixed by T1 (E1a, E1b, H10, A2a, O-3) and is the same
under both models. O-14 removed its only representation-sensitive behaviour: the notice join. What
remains is "status changes, accounts persist" — satisfiable by a field or a projection.

### 4.6 Ownership and authority to modify or retract — **independently decidable, and uncovered**

**DC-17 (Observed).** O-1 fixes *ownership*; nothing fixes *authority to retract an account*. `retract`
takes no actor in either implementation (`engine.py:542–559`; `model.py:282–299`), and no T1 or T2 text
assigns the operation (RQ-03). This is orthogonal to storage: it is the same gap under either model.

### 4.7 Recoverability and historical projection — **dependent on 4.1 and 4.2**

**DC-18 (Source-derived).** F1/F1b require knowledge to be reconstructable and a projection compared
against its reconstruction. K-R2/K-R14 make the comparison an integrity check. If no events exist
(4.1 = no), F1b has nothing to compare and becomes `Unestablished` rather than satisfied; if status is a
field (4.2 = field), K-R14's second half has no stream. Recoverability is thus the *consequence* of the
two storage choices, not a third choice.

---

## 5. Requirements neither model satisfies

**DC-19 (Observed + Source-derived).**

| Requirement | Why neither satisfies it |
|---|---|
| **K-R15** `(turn, seq)` | r32's own `AcquisitionEvent` schema has `turn` but no `seq` and no chat field (§7.4; ARCH-03 F6); the harness has `turn` only; the engine has no events. The obligation is specified against a schema that cannot carry it |
| **K-R19** | the harness emits one edge per event including `revoked` (§7.3), so "forgot" has an edge — the opposite of the rule |
| **F1b** | the engine has no projection to compare; its account *is* the knowledge representation, and §7.6 records that whether that makes it authoritative is `Unestablished` |
| **A2a at retrieval** | the engine's `retract()` can revoke a private account, and `perspective()` then hides it from its owner (§7.5) |
| **B3** event vs definition | no representation in either (PH1A-09) |
| **L1** human contribution | engine refuses human contributors (CO-08) |

**DC-20 (Inferred).** This is the finding that most changes the decision: **neither implementation is a
conforming subject for either model.** "Which model wins" is therefore not the operational question;
"which partial implementation is extended toward the model that governs" is.

---

## 6. What O-6a permits

**DC-21 (Source-derived).** O-6a's tier order constrains *conflicts*, not *bundles*. Nothing in it
requires both streams or neither:

- adopting r32's mechanisms is the default wherever T1 is silent and v9 does not contradict them;
- dropping an applicable r32 mechanism is **case 2** — a divergence that must be registered and decided,
  not an inference from a v9 rule;
- adding a mechanism r32 does not have is **case 3** — new architecture, declared as such.

**DC-22 (Inferred).** A mixed outcome is therefore **procedurally permitted**, provided each element is
declared under the case that governs it and no element contradicts T1. This is a statement about
permissibility, not a design: no combination is proposed, sketched or recommended here.

**DC-23 (Source-derived).** One asymmetry matters: dropping K-R12 is a divergence *from an obligation
nothing contradicts*, whereas dropping acquisition events may not even be a divergence if §4.1's reading
holds. The first would need an owner decision to depart from T2; the second may need only a textual
ruling.

---

## 7. PH1A-12, reassessed

**DC-24 (Superseded).** The post-Q-B review's finding DB-19 treated PH1A-03/04 as one root that must be
settled before PH1A-12, and DB-14 said PH1A-12 "is not answerable as posed until the mechanism set is
fixed". The dependency is real but was stated too coarsely. It is a dependency on **the set**, and two of
that set's members may not be decisions at all (§4.2, §4.1).

**DC-25 (Inferred).** With the decomposition, PH1A-12's content changes:

| Old framing | Decomposed framing |
|---|---|
| "Which model does the experiment test?" | "Which partial implementation is extended to express the T2-governed mechanisms, and where does the bridge live?" |
| blocked by PH1A-03/04 as a decision | blocked by (a) the §4.1 reading, (b) whether the owner diverges from K-R12, (c) the `seq` schema gap |
| answerable by choosing a codebase | answerable only after the required mechanism set is fixed — otherwise it answers the stream question by convenience |

**DC-26 (Observed).** The audit's premise still holds: no document names the subject (audit §15.2 item
7), and P9 (no bridge) is structural. What the decomposition adds is that the subject question cannot be
answered *well* before (a)–(c), even though it can be answered *quickly*.

---

## 8. RQ-03 reconciled with the sequence

**DC-27 (Source-derived).** RQ-03 — who may retract a character's account — is **independent of every
storage question**. It is the same gap under a stream model and a record model: `retract` has no actor,
no owner check and no authority gate in either implementation.

**DC-28 (Inferred).** It is **not** independent of §4.4. If forgetting and retraction are separated
(K-R3/K-R19), RQ-03 splits into two questions — who may cause *forgetting* (a knowledge operation) and
who may *retract content* (a truth operation) — with different authorities. If they are not separated,
one question covers both. So RQ-03 is conditionally dependent on the forgetting decision, and
unconditionally independent of the stream decisions.

**DC-29 (Inferred).** Its position in the sequence should therefore change. It is not a prerequisite of
PH1A-03/04 and must not be folded into it. It **is** a prerequisite of subject selection, because
whichever subject is chosen carries its retraction behaviour into the run, and A2a is a ratified row that
the engine currently fails at retrieval (§7.5).

---

## 9. Decision-dependency map

**Type key:** **R** = genuinely required from the owner · **C** = consequence of another decision ·
**D** = safely deferrable · **B(x)** = cannot be made until x is settled.

| Node | Question | Type | Depends on | Why |
|---|---|---|---|---|
| **DN-1** | Is v9 §8's "never a view, alias or projection" read with its antecedent (source record) or absolutely? | **R** | — | the only textual point deciding whether PH1A-03 is a divergence or a non-conflict (§4.1) |
| **DN-2** | Must acquisition events be stored, in addition to accounts? | **C** of DN-1 | DN-1 | if DN-1 = antecedent, K-R1/K-R11 govern and the answer is yes; if absolute, case-2 decision |
| **DN-3** | Must lifecycle status be a projection rather than a field (K-R12)? | **C** of tier order | — | v9 is silent; T2 governs; only a deliberate divergence would be an owner decision (§4.2) |
| **DN-4** | Must forgetting be distinguishable from retraction? | **R** | DN-2 | behavioural, T2, untested by any row; unanswerable without an acquisition record to be absent |
| **DN-5** | Must same-turn ordering be deterministic, and in which structure? | **R** | — | PH1A-05; T2 obligation no schema can carry (§5) |
| **DN-6** | Who may retract or revoke an account? | **R** | DN-4 (splits if yes) | RQ-03; uncovered by any tier (§4.6) |
| **DN-7** | Which implementation is the experimental subject, and where is the bridge? | **B(DN-1..DN-5, DN-6)** | as listed | PH1A-12; content changes under decomposition (§7) |
| **DN-8** | Do the four knowledge classes survive as a taxonomy? | **R** | — | PH1A-01, independent of all above |
| **DN-9** | Where does applicability live on MEMORY-09's side? | **R** | — | Q-D/RQS-01, independent; blocks state-aware experiments |
| **DN-10** | Incompatibility predicate; materiality scope; event-vs-definition; REV 2 registry | **D** | — | PH1A-07/08/09/13; none gates subject selection |
| **DN-11** | Human contribution as a claim (L1) | **D** | — | PH1A-10; implementation-level, one row |
| **DN-12** | Ownership of the retrieval trigger | **D** | — | PH1A-11; Howl boundary, no experiment depends on it |
| **DN-13** | Q-A (sameness) | **D** | — | no requirement needs it (DB-04) |
| **DN-14** | Q-C (notice attribution) | **D — moot** | — | O-14 removed its object (DB-06) |

### 9.1 What this map shows

**DC-30 (Inferred).** PH1A-03/04 should **not** remain one decision. It decomposes into:

1. **one reading** (DN-1) — a textual question with two answers, neither an architecture;
2. **one consequence** (DN-2) that follows from the reading;
3. **one obligation that is not a decision at all** (DN-3) — T2 governs because v9 is silent;
4. **two behavioural decisions that were hiding inside the bundle** (DN-4 forgetting, DN-5 ordering),
   of which DN-5 is already PH1A-05.

**DC-31 (Inferred).** Three further results follow, and they are the reason this pass was worth doing:

- **The bundle was doing work it should not.** Presenting DN-3 as an owner choice invites a divergence
  from K-R12 that nothing motivates, and presenting DN-1 as an architecture choice hides that it is a
  one-sentence reading.
- **The genuinely open behavioural question is forgetting** (DN-4), not event-sourcing. It is the only
  mechanism whose observable difference the two models produce, it has no probe, and it is what
  distinguishes "misinformed" from "amnesiac" — which is the same distinction O-14 just protected.
- **RQ-03 (DN-6) belongs after DN-4, not inside PH1A-03/04.** Folding it in would decide by implication
  whether forgetting and retraction are one operation.

---

## 10. What this report does not do

No architecture selected. No hybrid designed — §6 assesses permissibility only and proposes no
combination. No implementation, test, probe, model call or corpus. No canonical edit, registry change or
Git action. The PH1A-03/04 decision is **not** put to the owner.

**Limits.** All findings are static readings at the hashes in the header. §4.1's two readings of v9 §8
are offered as readings, not ranked; DN-1 exists precisely because this report cannot choose between
them. The engine's conformance claims are source readings, not executions.

**End of report.**
