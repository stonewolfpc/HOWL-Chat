# MEMORY-09 — S-2 Authorization Request

## Authorizing the lifecycle-representation decision process, not its result

**Status:** draft authorization request. **Not an owner decision. No carrier is proposed or preselected. No schema,
code, migration, test, probe, requirement, map change or canonical capture is requested or authorized.** This
document asks the owner to authorize **the S-2 decision process** within stated boundaries.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Requested of | **the owner** |
| Subject | **S-2 only** — the lifecycle representation of post-formation `scope` and `knowledge_class` changes |
| Request | authorize S-2 to decide the **logical representation**, or to decline with the obligation recorded as unmet, under the boundaries below |
| Pivotal choice | **whether A-6 applies to S-2** — this defines the decision's scope, not its result |
| Not requested | a carrier design, a schema, implementation, requirement capture, canonical capture, or any decision about S-1/S-3 |
| Read with | **Howl separation addendum** (`MEMORY-09-S2-HOWL-SEPARATION-ADDENDUM-2026-11-20.md`, S2H-01…S2H-22); part of this request and binding on any A/B/C/D selection |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Why this request exists

**S2R-01 (Observed — the finding chain).** Q2 = **Yes**: post-formation `scope` and `knowledge_class` changes are
within the model's semantic domain (Q2A-33). K-R12 creates **O-KR12(d)**: every actual class/scope change must be
recorded in the lifecycle stream and projected from it (KRO-06). The current vocabulary has no class/scope entry;
the null carrier is buildable but incomplete (C3D-11; KRO-27).

**S2R-02 (Observed — the ownership finding).** Q1 = **S-2**, the lifecycle representation outside N-a (Q1S-31).
**S-3** (K-R12's projection requirement) is the acceptance test, not a co-owner. **S-1** (N-a's fold) is a
downstream consumer.

**S2R-03 (Observed — the audit finding).** The S-2 ownership and dependency audit **passes with conditions C-1…C-8**
(S2A-45/47). The decisive condition is **A-6's applicability** (S2A-46): if A-6 carries into S-2, the decision cannot
genuinely consider a new carrier and becomes a gap-acceptance decision; if it does not, S-2 may consider logical
representations, but still not physical implementation.

**S2R-04 (Observed — what this document is).** A request for **authorization to decide**, not a decision and not a
carrier proposal. It contains no candidate carrier designs. The owner is asked to fix the process, its scope and its
gates — not to preselect its result.

---

## 2. Subject: S-2 only

**S2R-05 (Calculated — the subject).** S-2 decides **whether, and in what logical form**, a post-formation class or
scope change is represented in the lifecycle stream, such that K-R12's projection obligation can be satisfied — **or**
issues a reasoned decline recording the obligation as unmet.

**S2R-06 (Calculated — what is not the subject).**

| Subject | Status in this request |
|---|---|
| **S-3** — K-R12 projection requirement | **acceptance test only**; not authorized to design or decide a carrier |
| **S-1** — N-a fold scope | **downstream**; not authorized; no fold behavior, shape, entry point or order sensitivity |
| **N-b** — status vocabulary | not authorized |
| **N-c** — status-stream ordering | not authorized; S-2's own ordering rule, if any, is labelled separately |
| **§11 / C-1** — promotion authority | not authorized |
| **O-8a / v9 `domain`** — authority-promotion axis | not authorized |

**S2R-07 (Calculated — why only S-2).** S-1 cannot own the carrier (A-6 bars it); S-3 is a requirement, not an
authority; the other subjects are different axes. Opening more than S-2 would reproduce the conflation this boundary
was built to prevent (Q1S-09…14).

---

## 3. The decision question

**S2R-08 (Calculated — the primary question).**

> Shall a post-formation `scope` or `knowledge_class` change be represented in the lifecycle stream — and if so, in
> what **logical** form, with what projection rule and check — or shall representation be declined, with K-R12's
> obligation recorded as unmet?

**S2R-09 (Calculated — the logical sub-questions).** If the answer is to represent:

| # | Sub-question | Boundary |
|---|---|---|
| 1 | what logical carrier form carries the change **in the lifecycle stream** | no physical schema, no code |
| 2 | what projection rule derives the current class/scope from the stream | must satisfy K-R12 and §6.4 |
| 3 | what check compares computed to stored, and what counts as a mismatch | must be deterministic; mismatch is an integrity error |
| 4 | what entry shape a class/scope record has | shape only; not the fold's transition function |
| 5 | whether an ordering rule is needed for those entries | **S-2's own**, labelled separately; must not settle N-c |

**S2R-10 (Calculated — what is deliberately absent).** No candidate carrier is named, described, compared or
recommended. The request does not contain an option set. Developing options is part of the authorized S-2 process, not
part of this request.

---

## 4. Scope: in and out

**S2R-11 (Calculated — in scope).**

| In scope | Why |
|---|---|
| whether to represent or decline | the null outcome is a decision, not a failure |
| the logical carrier form | the representation question |
| the projection rule | required by K-R12/§6.4 |
| the check | required by §6.4 |
| the entry shape | representation, not fold |
| an S-2-specific ordering rule, if needed | may be necessary; must not become N-c |

**S2R-12 (Calculated — out of scope).**

| Out of scope | Gate |
|---|---|
| physical schema, DDL, indexes, storage layout | separate implementation authorization |
| code, migration, tests, probes, model calls | separate implementation authorization |
| canonical edits (r33, v9, ledger, registry) | separate canonical-capture authorization |
| requirement capture / map change | separate capture step; map stays **118** |
| authority, permission, sufficiency | separate question |
| validity, warrant, truth | separate question |
| causes of class/scope change | D5 / K-R9 / K-R10 / I13; not S-2's |
| N-b, N-c, N-a's fold | S-1 |
| §11 promotion, O-8a authority-promotion, v9 `domain` | separate axes |

**S2R-13 (Calculated — the logical/physical line).** "Logical" means the record's **semantics and invariants**: what
the carrier means, what the projection derives, what the check compares, what shape the entry has. "Physical" means
schema, storage, code, migration and deployment. The S-2 decision may cross the first line; it does not cross the
second.

---

## 5. Inherited constraints

**S2R-14 (Calculated — the constraints S-2 must satisfy).**

| Constraint | What S-2 inherits |
|---|---|
| **K-R12** | carrier in the lifecycle stream; covers class/scope; current value is a projection; no mutable authoritative field |
| **K-R14** | independently rebuildable from the lifecycle stream; no cross-stream input; no stored join |
| **§6.4** | memory events project to status/scope/class; the check is deterministic; mismatch is an integrity error |
| **K-R2** | the **integrity principle** — an unreconstructable projection is an integrity failure; **not** its edge-specific `derived_from_event_ids` literal clause, whose object is `KnowledgeEdge` |
| **K-R13** | knowledge × memory is a read-time join; no stored cross-stream join |
| **K-R15** | acquisition ordering only; not an automatic ordering solution for lifecycle entries |

**S2R-15 (Source-derived — K-R2's exact reach).** K-R2 says: "An edge carries `derived_from_event_ids`. A projection
that cannot be reconstructed from its events is an integrity failure, not a policy question." Its literal object is
the **edge**. The memory projection's binding texts are §6.4 and K-R14; S-2 inherits K-R2's **principle**, not its
edge clause. The request must not be read as importing `derived_from_event_ids` onto the memory stream.

**S2R-16 (Source-derived — K-R15's exact reach).** K-R15 fixes `(turn, seq)` for the **acquisition** fold. It does
not supply an ordering rule for lifecycle entries. If S-2 needs one, that is a new S-2 question; it must be labelled
S-2's own and must not silently settle N-c.

---

## 6. Non-absorption boundaries

**S2R-17 (Calculated — questions S-2 must not absorb).**

| Question | Owner | Why not S-2 | Fence |
|---|---|---|---|
| **authority / permission** — who may change class/scope | §11 / C-1 R-E | representation is not permission | S-2 may record a change; it does not decide who may cause one |
| **validity** — whether a change is warranted or true | D6; unresolved validity predicate | class/scope are not truth values | S-2 represents; it does not validate |
| **causes** — authorship, repetition, age, decay, popularity | D5 / K-R9 / K-R10 / I13 | these rules constrain causes | S-2 may carry a value; it may not relax or extend them |
| **N-b** — `active`/`established` vocabulary | S-1 | status mapping is a fold-input question | S-2 must not rename or map statuses |
| **N-c** — status-stream ordering | S-1, conditional on A-5 | N-c concerns status | S-2's ordering rule is separate and must not resolve N-c |
| **N-a's fold** — shape, entry point, links, order sensitivity | S-1 | the fold consumes the representation | S-2 decides the record; S-1 decides the transition function |
| **§11 promotion** | §11 / C-1 R-E | promotion is a different operation | S-2 must not decide promotion permission or sufficiency |
| **O-8a** — authority-promotion | separate world/scenario record | a different axis | S-2 must not merge it with r33 `scope` |
| **v9 authority `domain`** | v9 / O-8a | a different axis | S-2 uses r33 `scope` values only |
| **state-machine orchestration / UX** | Howl state machine / application behavior | S-2 defines memory capabilities, not selection, switching, combination, activation/deactivation or restore | S-2 may provide the information; it does not decide how it is consumed (S2H-11…S2H-14) |

**S2R-18 (Calculated — the three highest-risk slides).** **Authority** (representation becoming permission),
**entry architecture** (shape becoming fold behavior), and **ordering** (an S-2 rule becoming N-c). The request names
all three so the S-2 process cannot absorb them by implication.

---

## 7. The A-6 applicability question

**S2R-19 (Observed — why this is pivotal).** A-6 is a term of the **N-a decision**: it creates no event kind, field,
key or schema, and settles neither authority nor validity (NAE-31/32; KRO-17). S-2 is not N-a, so A-6 does not apply
**by default**. But if the owner imports it, S-2's decision space changes fundamentally (S2A-46).

**S2R-20 (Calculated — the two options).**

| Option | What it means | What S-2 can genuinely decide | What it becomes |
|---|---|---|---|
| **A — A-6 applies** | S-2 inherits the no-kind/no-field/no-key/no-schema bar | only representations that use the existing vocabulary; if none satisfies K-R12, a reasoned decline | a **gap-acceptance** decision, not a representation decision |
| **B — A-6 does not apply** | S-2 may consider logical representations that require new kinds/fields/keys | the logical carrier, projection rule, check, entry shape | a **representation** decision, still without physical implementation |

**S2R-21 (Calculated — this is a scope choice, not a carrier choice).** Selecting A or B does not select a carrier.
It defines whether the decision may consider a new carrier at all. The request contains no candidate designs either
way.

**S2R-22 (Calculated — the question must be answered).** If neither option is selected, the authorization is
incomplete: S-2 cannot open with a defined decision space. The owner may also decline to authorize S-2 (S2R-30).

---

## 8. Deliverable and gates

**S2R-23 (Calculated — the deliverable).** A **decision record** stating: the chosen logical representation or a
reasoned decline; the projection rule; the check; the entry shape; any S-2-specific ordering rule; the inherited
constraints; the non-absorption boundaries; and an explicit statement that implementation, requirement capture and
canonical capture are **not** authorized by it.

**S2R-24 (Calculated — the success criterion).** The decision passes S-3 if it yields a class/scope projection that is
derivable from the lifecycle stream, independently rebuildable (K-R14) and checkable (§6.4). If it declines, it must
record which part of O-KR12 remains unmet; it must not deny the semantic domain (Q2 = Yes).

**S2R-25 (Calculated — separate gates).**

| Gate | What it covers | State after S-2 |
|---|---|---|
| **Implementation** | schema, DDL, code, migration, tests, probes, model calls | **not authorized**; separate owner authorization required |
| **Requirement capture** | adding a conformance requirement; changing the map | **not authorized**; map stays **118**; separate step |
| **Canonical capture** | editing r33, v9, ledger or registry | **not authorized**; separate step |
| **Owner decision on the S-2 record** | accepting, modifying or rejecting the S-2 result | a separate act; this request does not presume it |

**S2R-26 (Calculated — no implicit authorization).** Authorizing S-2 authorizes **analysis and a decision record**; it
does not authorize any of the four gates above. The decision record must repeat that sentence.

---

## 9. Decline as a valid outcome

**S2R-27 (Calculated — decline is a decision).** The owner may decline to represent class/scope changes. A decline is
not a failure to decide and not a semantic denial. It must record:
- that O-KR12 remains **unmet** for the declined dimension(s);
- that Q2 = Yes stands;
- that the current null carrier is **incomplete**, not a complete model;
- that no implementation, capture or canonical change follows.

**S2R-28 (Calculated — why the option must be explicit).** Without it, the authorization pressures the owner toward
an architectural answer. The decline option preserves the difference between "we choose not to represent" and "the
model cannot represent" — only the first is a decision.

**S2R-29 (Calculated — decline does not close S-2 forever).** A decline can be revisited by a later authorization.
It is a disposition, not an irreversible architectural fact.

---

## 10. What the owner is not being asked to decide

**S2R-30 (Calculated).**

| Not asked | Why |
|---|---|
| a carrier design or option | the request authorizes the process, not its result |
| physical schema or code | separate implementation gate |
| who may change class/scope | authority is §11/C-1, not S-2 |
| whether a change is valid or warranted | validity is separate |
| N-b, N-c or N-a's fold | S-1 |
| §11 promotion, O-8a, v9 `domain` | separate axes |
| state-machine design or UX behavior | Howl state machine / application behavior; see S2H |
| how memories are selected, switched, combined, activated or restored | state machine, not memory representation |
| a requirement or map change | separate capture step |
| canonical edits | separate canonical-capture step |

---

## 11. Authorization options

**S2R-31 (Calculated — the owner's choices).**

**All options are subject to the Howl separation addendum (S2H-01…S2H-22).**

| Option | Authorization |
|---|---|
| **A** | **Authorize S-2 with A-6 applying.** The decision is confined to existing-vocabulary representations; if none satisfies K-R12, a reasoned decline/gap acceptance is the expected outcome. |
| **B** | **Authorize S-2 with A-6 not applying.** The decision may consider logical representations requiring new kinds/fields/keys, subject to K-R12/K-R14/§6.4; no physical implementation. |
| **C** | **Decline to authorize S-2.** The status quo stands: the class/scope obligation remains acknowledged and unrepresented. |
| **D** | **Authorize with modifications.** The owner states the modifications; they must preserve the S-2-only subject and the separate gates. |

**S2R-32 (Observed — what a response should state).** The owner's response should state the chosen option, and for A
or B, confirm that the non-absorption boundaries in §6 and the separate gates in §8 apply. No carrier should be
named in the response; naming one would preselect the S-2 result.

---

## 12. What this request does not do

**S2R-33 (Observed).**

| Untouched | State |
|---|---|
| **carrier / schema / kind / field** | none proposed, named or preselected |
| **S-2 decision** | not opened; not authorized |
| **S-1 / S-3** | not authorized; S-3 is the test; S-1 is downstream |
| **authority / validity / causes** | fenced, not decided |
| **N-b / N-c / N-a fold** | not authorized |
| **§11 / O-8a / v9 domain** | fenced, not decided |
| **implementation / requirement capture / canonical capture** | not authorized |
| **map / requirements** | **118**; no change |
| **code / schema / migration / tests** | **0** |
| **state-machine design / UX** | not authorized; S2H |

---

## 13. Traceability

**S2R-34 (Observed).**

| What | Where |
|---|---|
| Q2 = Yes | Q2A-33 |
| O-KR12; class/scope gap | KRO-06; KRO-27 |
| Q1 = S-2; S-3 test; S-1 downstream | Q1S-31; Q1S-34/35 |
| Audit passes with C-1…C-8; A-6 pivotal | S2A-45/46/47 |
| K-R12: locus, domain, projection, non-mutation | `:111` |
| K-R14: independent rebuild; no stored join | `:113` |
| §6.4: projection check; two streams, two checks | `:366`–`:367` |
| K-R2: edge clause; integrity principle | `:83` |
| K-R13: read-time join | `:112` |
| K-R15: acquisition ordering only | §6.4 |
| A-6 binds N-a, not S-2 by default | NAE-31/32; KRO-17 |
| NSD-09 binds ND-1, not S-2 | owner scope C-6; KRO-19 |
| Non-absorption; highest-risk slides | S2A-38/39 |
| Separate gates | S2A-41; S2R-25 |
| MEMORY-09 / Howl separation | S2H-01…S2H-22; owner, 2026-11-20 |

---

## 14. Record status

**S2R-35 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Carrier proposed or preselected | **0** |
| S-2 opened or authorized | **0** — this is a request, not an authorization |
| Authorization request contains candidate designs | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **authorization request draft** — not an owner decision, no canonical capture |

**End of request.** The owner is asked to authorize **S-2 only** — the logical representation decision for
post-formation `scope` and `knowledge_class` changes — under the inherited constraints (K-R12, K-R14, §6.4, K-R2's
integrity principle, K-R13, K-R15), the non-absorption boundaries (authority, validity, causes, N-b, N-c, N-a's fold,
§11, O-8a, v9 `domain`), and the separate gates for implementation, requirement capture and canonical capture. **The
pivotal choice is A-6's applicability** (S2R-20…22): A makes S-2 a gap-acceptance decision; B makes it a
representation decision; neither selects a carrier. **Decline is a valid outcome** and must record O-KR12 as unmet
(S2R-27…29). **No candidate carrier design appears in this request, and no decision is made by it** (S2R-04, S2R-33).
