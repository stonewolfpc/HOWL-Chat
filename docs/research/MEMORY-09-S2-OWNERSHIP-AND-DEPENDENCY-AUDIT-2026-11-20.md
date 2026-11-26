# MEMORY-09 — S-2 Ownership and Dependency Audit

## What the lifecycle-representation decision inherits, what it must not absorb, and whether it can be bounded

**Status:** final ownership and dependency audit. **Research evidence only — no canonical edit, no schema, no code,
no test, no probe, no ledger entry, no implementation, no owner authorization request.** **No carrier is selected;
the 118-row map is unchanged.** C-1 and C-2 are not reopened. Q2 = Yes and Q1 = S-2 are carried forward.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | can the S-2 decision be authorized without absorbing authority, validity, N-b, N-c or N-a's fold, and without implicitly authorizing implementation or canonical capture? |
| Verdict | **passes with explicit fences** — S-2 is separable, but only if the authorization request names its subject, exclusions and gates |
| Inherited constraints | **K-R12** (locus, domain, projection, non-mutation); **K-R14** (independent rebuild, no stored join); **§6.4** (memory projection check); **K-R2's integrity principle** (literal object: the edge; applied to the memory projection by §6.4/K-R14); **K-R13** (read-time join); **K-R15** (acquisition ordering only) |
| Must not absorb | authority/permission (§11, C-1 R-E); validity; causes (D5, K-R9, K-R10, I13); N-b; N-c; N-a's fold; O-8a; v9 authority `domain`; requirement capture; implementation; canonical capture |
| Next artifact | a **narrowly scoped authorization request for S-2**, not a carrier proposal — permissible only after the conditions in §8 are met; **not drafted here** |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Task and carried state

**S2A-01 (Observed — the owner's request).** Before any S-2 authorization request: confirm every constraint S-2 would
inherit (especially K-R12, K-R14, K-R2); identify unresolved questions S-2 must **not** absorb (authority, validity,
N-b, N-c, N-a's fold); and confirm whether S-2 can be bounded without implicitly authorizing implementation or
canonical capture. If the audit passes, the next artifact **can** be a narrowly scoped authorization request — not a
carrier proposal.

**S2A-02 (Observed — carried forward).** Q2 = **Yes** (Q2A-33). Q1 = **S-2**, the lifecycle representation outside
N-a (Q1S-31). **S-3** is the acceptance test, not a co-owner. **S-1** is the downstream fold consumer. The current
vocabulary has no class/scope entry; the null carrier is buildable but incomplete (C3D-11; KRO-27).

**S2A-03 (Observed — the boundary).** No carrier, schema, kind, field, requirement, map change, code, migration or
canonical capture. This record audits **whether S-2 can be opened safely**, not what it should decide.

**S2A-04 (Observed — what "bounded" means here).** An S-2 decision is bounded if its **subject**, **deliverable**,
**exclusions** and **gates** are stated so that deciding the representation does not, by itself, decide authority,
validity, the fold, a conformance requirement, an implementation, or a canonical text change.

---

## 2. S-2's ownership, restated

**S2A-05 (Calculated — the subject).** S-2 decides **how a class or scope change is represented in the lifecycle
stream**, including whether to represent it at all; the logical form of the carrier; the projection rule that derives
the current value; the check that compares computed to stored; and the entry shape needed for such a record.

**S2A-06 (Calculated — what it does not own).** S-2 does not own the fold (S-1), the status vocabulary (N-b), the
status ordering rule (N-c), authority or permission (§11/C-1), validity (D6 and the unresolved validity predicate),
the causes of class/scope change (D5/K-R9/K-R10/I13), the requirement capture (map), or implementation/canonical
state.

**S2A-07 (Calculated — why S-2 is the right owner).** K-R12's obligation is a property of the stream and its
projection, not of the fold (KRO-37). A-6 bars S-1 from creating the carrier; K-R14 bars an external join; S-3 is a
test, not an authority. S-2 is the only frame whose scope includes the carrier without importing A-6 (Q1S-28…31).

---

## 3. Constraint audit: what S-2 inherits

### 3.1 K-R12 — direct

**S2A-08 (Source-derived — the rule, `:111`).** The lifecycle is its own authoritative stream; its subject includes
the four status kinds **plus** scope and class changes; current status, class and scope are a **projection** of the
stream, not mutable authoritative fields.

**S2A-09 (Calculated — what S-2 inherits).**

| Clause | S-2's obligation |
|---|---|
| **locus** | the carrier must belong to the lifecycle stream, not to the acquisition stream or an external store |
| **domain** | the representation must cover class/scope changes, not only status |
| **projection** | current class/scope must be derivable from the stream; S-2 must supply the derivation rule |
| **non-mutation** | the authoritative value must not be a mutable field; a cache is permitted only as a projection with a rebuild rule (P-5) |

### 3.2 K-R14 — direct

**S2A-10 (Source-derived — the rule, `:113`).** "Each stream is rebuilt and checked independently; a mismatch in
either is an integrity failure. The cross-stream join has no stored form, so it cannot drift."

**S2A-11 (Calculated — what S-2 inherits).** The class/scope carrier must be rebuildable from the lifecycle stream
alone; the projection must not consult the acquisition stream; no cross-stream join may be stored. **P-4 (external
operation joined at projection) is therefore excluded by K-R14, not by A-6** (C3D-13; KRO-25). S-2 may consider an
in-stream carrier only.

### 3.3 §6.4 — direct

**S2A-12 (Source-derived — the projection check, `:366`).** "Projection check is deterministic and cheap enough to
run on write: fold events → projection; compare to stored projection; mismatch is an integrity error surfaced to
`inspection`, never silently resolved in favour of either side." And: "Acquisition events project to edges; memory
events project to the memory's status/scope/class. Neither projection may consult the other."

**S2A-13 (Calculated — what S-2 inherits).** S-2 must specify a class/scope projection rule that is **checkable**:
there must be an expected value to compare against. A null carrier over a stored cache fails this check silently —
it recomputes the formation value and cannot detect a missing change (C3D-08). S-2's representation must make the
class/scope change visible to the check.

### 3.4 K-R2 — the integrity principle, not the edge clause

**S2A-14 (Source-derived — the rule, `:83`).** "An edge carries `derived_from_event_ids`. A projection that cannot be
reconstructed from its events is an **integrity failure**, not a policy question."

**S2A-15 (Calculated — K-R2's literal object is the edge).** K-R2 names `derived_from_event_ids` on `KnowledgeEdge`,
so its **edge-specific clause** does not directly bind the memory projection. Its **integrity principle** —
unreconstructable projection = integrity failure — is applied to the memory stream by §6.4 and K-R14.

**S2A-16 (Calculated — what S-2 inherits).** S-2 inherits the **principle**, not the edge clause: a class/scope
projection that cannot be reconstructed from lifecycle events is an integrity failure, not a policy choice. The
audit does not overstate K-R2's reach; the binding text for the memory projection is §6.4/K-R14.

### 3.5 K-R13 — read-time join

**S2A-17 (Source-derived — the rule, `:112`).** "'This character's knowledge includes a retracted claim' is a
**read-time join** of edge × memory projection. It is never stored on the edge."

**S2A-18 (Calculated — what S-2 inherits).** S-2 must not store a cross-stream join; class/scope representation
belongs to the memory stream, and any use with knowledge is a read-time join. This also fences S-2 from N-a's fold
and from entry-architecture shortcuts that would put memory status/class/scope on the edge (§6.3).

### 3.6 K-R15 — acquisition ordering only

**S2A-19 (Source-derived — the rule, §6.4).** "Fold determinism (K-R15, K-R20): replaying the same ordered log
produces the same edge…" K-R15 fixes `(turn, seq)` for the **acquisition** fold; it does not supply an ordering rule
for the memory lifecycle.

**S2A-20 (Calculated — what S-2 inherits, and what it must not borrow).** S-2 may need an ordering rule for
class/scope entries; it **cannot borrow** K-R15's acquisition key without a new decision. That ordering question is
**S-2's own**, not N-c: N-c asks whether the **status** stream needs an order key (NSD-11/12; Q1S-26). S-2 must not
resolve N-c while answering its own ordering need.

### 3.7 §11 — authority and deletion

**S2A-21 (Source-derived — §11, `:608`–`:614`).** Promotion per kind: episodic on witnessing; preference immediate
for the holder; belief never promotable to fact; rumour stays provisional; world fact by authority or corroboration;
absence attributed. Retraction applies to records, not minds. §11.1 is a **chat deletion** contract (proposed, `[P]`).

**S2A-22 (Calculated — what S-2 inherits and must not absorb).** S-2 must not decide **who may promote** a record or
whether a class/scope change is permitted; C-1 = R-E confines §11's "authority" to epistemic support for
record-promotion, not permission (C1D-01…13). §11.1's deletion **by scope** is deletion, not a scope **change**; S-2
must not conflate them.

### 3.8 D1/D5, K-R8/K-R9/K-R10/I13 — scope values and class causes

**S2A-23 (Source-derived — scope values, D1/§5).** `world` / `scenario` / `character_public` / `character_private`
are **scopes**, not owners (D1); §5 lists them as the `scope` field (`:318`). S-2 must not invent values or change
ownership semantics.

**S2A-24 (Source-derived — class causes).** D5: decay never affects truth, ownership, **scope**, eligibility or
history. K-R10: an experiential event never downgrades to ambient by repetition or age; decay does not change class.
K-R9: no record becomes ambient because many characters know it; popularity is not authorship. K-R8: ambient applies
only to authored world/scenario definitions. I13: class is authored, not inferred from popularity or age.

**S2A-25 (Calculated — what S-2 inherits and must not absorb).** S-2 must represent a class/scope change **without**
deciding which causes are permitted, whether authorship is sufficient, or whether a change is warranted. Those are
the cause/warrant questions K-R9/K-R10/D5/I13 constrain; S-2 must not relax or extend them. It also must not infer a
class or scope from retrieval, popularity or decay.

### 3.9 Axis separation: r33 `scope` ≠ v9 `domain`; O-8a

**S2A-26 (Source-derived — the separation).** r33's `scope` is the memory field (`world`/`scenario`/…); v9's
authority **domain** is the per-world/scenario domain of O-8a (RX-17; C3Q-05). They are different axes. O-8a's
authority-promotion is a separate world/scenario record, not §11's record-promotion and not a memory `scope` change
(CPA-05…07, CPA-30…36).

**S2A-27 (Calculated — S-2's fence).** S-2 must use r33's `scope` values and must not merge them with v9's domain or
O-8a's authority-promotion. A class/scope representation is not an authority grant.

### 3.10 A-6, NSD-09, the map and the read-only boundary

**S2A-28 (Source-derived — A-6 is not inherited).** A-6 binds the **N-a decision** (NAE-31/32; KRO-17). S-2 is not
N-a, so A-6 does not apply by default. **The authorization request must state whether A-6 is imported.** If imported,
S-2 cannot consider a new kind or field, and the representation question is effectively reduced to P-3/P-5 or a
K-R14-barred P-4 — a scope choice the owner must make explicitly.

**S2A-29 (Source-derived — NSD-09 is not inherited).** NSD-09 binds the **ND-1 determination**, not a separately
authorized S-2 decision (KRO-19). It is the **authorisation boundary**: identifying S-2 does not authorize it. The
request must be an explicit owner authorization, not an inference from NSD-09.

**S2A-30 (Observed — the map and read-only boundary).** The map stays at **118**; no requirement is captured. This
record and any S-2 decision record are research/architecture documents, not canonical edits.

---

## 4. Dependency audit

**S2A-31 (Calculated — the dependency graph).**

```
S-2  representation decision
  ├─ inherits K-R12 / K-R14 / §6.4 / K-R2 principle / K-R13 / K-R15
  ├─ is tested by S-3 (O-KR12 projection acceptance)
  ├─ feeds S-1 (fold scope decision, separate and downstream)
  ├─ owns Q5 (carrier), but does not answer it by being authorized
  └─ narrows Q4 (deferred vs excluded), but does not decide it
```

**S2A-32 (Calculated — S-2 → S-3).** S-2 proposes a representation; S-3 tests whether it yields the required
projection. S-3 does not design the carrier and does not authorize it. If S-2 declines to represent, S-3 records the
gap as an unmet obligation, not as a semantic denial (Q2 = Yes).

**S2A-33 (Calculated — S-2 → S-1).** S-1's fold scope is downstream: S-2 decides the record; S-1 decides whether the
fold consumes it. S-2 must not decide the transition function, the fold shape, the entry point or order sensitivity.
If S-2 creates a carrier, S-1 still needs its own amendment to fold it (Q1S-38).

**S2A-34 (Calculated — S-2 ↔ Q5).** S-2 owns the carrier question. The authorization request authorizes the
**decision process**, not a preselected carrier. S-2 may choose, or decline; either is a decision, not an
implementation.

**S2A-35 (Calculated — S-2 ↔ Q4).** S-2's outcome narrows Q4: "deferred to S-2" is coherent; "excluded" is
inconsistent with Q2 = Yes. The disposition remains the owner's; S-2 does not decide it by deciding representation
(Q1S-40).

**S2A-36 (Calculated — S-2 ↔ entry architecture).** S-2 may need to decide the entry **shape** for class/scope
records. That is representation, not fold behavior. S-2 must not use entry architecture to decide how a later event
changes current state — that is S-1.

**S2A-37 (Calculated — S-2 ↔ ordering).** S-2 may need an ordering rule for its entries; that is a new S-2 question.
It must not decide N-c (status ordering) and must state whether its rule is independent of N-c (S2A-20).

---

## 5. Non-absorption audit

**S2A-38 (Calculated — the questions S-2 must not absorb).**

| Question | Owner | Why S-2 must not absorb it | Fence |
|---|---|---|---|
| **authority / permission** — who may change class/scope | §11 / C-1 R-E; separate | representation is not permission | S-2 may record a change; it may not decide who is allowed to cause one |
| **validity** — whether a change is warranted or true | D6; unresolved validity predicate | class/scope are not truth values | S-2 represents; it does not validate |
| **causes** — authorship, repetition, age, decay, popularity | D5 / K-R9 / K-R10 / I13 | these rules constrain causes; S-2 must not relax them | S-2 may carry a value; it may not decide what may produce it |
| **N-b** — status vocabulary (`active`/`established`) | S-1 | status mapping is a fold-input question | S-2 must not rename or map statuses |
| **N-c** — status-stream ordering key | S-1, conditional on A-5 | N-c is about status, not class/scope | S-2's own ordering rule must be labelled separate and must not resolve N-c |
| **N-a's fold** — shape, entry point, links, order sensitivity | S-1 | the fold consumes the representation; it does not define it | S-2 decides the record; S-1 decides the transition function |
| **§11 promotion** — record-promotion authority | §11 / C-1 R-E | promotion is a different operation from class/scope change | S-2 must not decide promotion permission or sufficiency |
| **O-8a authority-promotion** | separate world/scenario record | a different axis from r33 `scope` | S-2 must not merge them |
| **v9 authority `domain`** | v9 / O-8a | a different axis from r33 `scope` | S-2 must use r33 values only |
| **requirement capture** | map / separate step | the map is a conformance artefact, not a design record | S-2's decision does not add a row; capture is separate |
| **implementation** | separate authorization | a decision is not code, schema or migration | S-2's deliverable is a decision record |
| **canonical capture** | separate authorization | a decision record is not a canonical edit | no r33/v9/ledger/registry change |

**S2A-39 (Calculated — the most likely absorptions).** The three highest-risk expansions are: (i) **authority** —
"who may change scope" feels adjacent to representation; (ii) **entry architecture** — shaping entries can slide into
fold behavior; (iii) **ordering** — a class/scope entry may need an order rule, and N-c is nearby. The request must
name all three explicitly.

---

## 6. Boundary specification for the S-2 decision

**S2A-40 (Calculated — in scope).**

| In scope | Why |
|---|---|
| whether to represent a class/scope change at all, or decline | the decision includes the null option; declining is a decision, not a failure |
| the **logical** carrier form (in-stream kind, attribute, or another in-stream form) | this is the representation question |
| the projection rule that derives the current value | required by K-R12/§6.4 |
| the check that compares computed to stored | required by §6.4 |
| the entry shape needed for a class/scope record | representation, not fold |
| an ordering rule **for class/scope entries**, labelled as S-2's own | may be needed; must not become N-c |

**S2A-41 (Calculated — out of scope).**

| Out of scope | Gate |
|---|---|
| physical schema, DDL, indexes, storage layout | separate implementation authorization |
| code, migration, tests, probes, model calls | separate implementation authorization |
| canonical text edits (r33, v9, ledger, registry) | separate canonical-capture authorization |
| requirement capture / map change | separate capture step; map stays 118 |
| authority, permission, sufficiency, validity | separate questions |
| causes (authorship, repetition, age, decay, popularity) | D5/K-R9/K-R10/I13; not S-2's |
| N-b, N-c, N-a's fold | S-1 |
| §11 promotion, O-8a authority-promotion, v9 domain | separate axes |

**S2A-42 (Calculated — the deliverable).** A **decision record** stating: the chosen logical representation or a
reasoned decline; the projection rule; the check; the entry shape; any S-2 ordering rule; the inherited constraints;
and an explicit statement that implementation and canonical capture are **not** authorized by it.

**S2A-43 (Calculated — the success criterion).** The decision passes S-3 if it yields a class/scope projection that
is derivable from the lifecycle stream, independently rebuildable (K-R14), and checkable (§6.4). If it declines, it
must record the unmet obligation, not deny the semantic domain (Q2 = Yes).

---

## 7. Can it be bounded?

**S2A-44 (Calculated — yes, with conditions).** S-2 is separable from authority, validity, the fold, N-b, N-c,
implementation and canonical capture. The separation is not automatic: it depends on the authorization request
stating the subject, exclusions and gates. Without those, "representation" can silently expand into any of them.

**S2A-45 (Calculated — the conditions).**

| # | Condition |
|---|---|
| **C-1** | the request names **S-2** as the subject and states that it does not authorize S-1, S-3 or a carrier |
| **C-2** | the request states the **logical-only** deliverable: no code, schema, migration, tests or probes |
| **C-3** | the request states **no canonical capture**: no r33/v9/ledger/registry edit |
| **C-4** | the request states **no requirement capture**: the map stays 118 |
| **C-5** | the request states the **non-absorption** list: authority, validity, causes, N-b, N-c, N-a fold, §11, O-8a, v9 domain |
| **C-6** | the request states **whether A-6 is imported**; if it is, S-2's schema-bearing options are barred and the scope is narrower |
| **C-7** | the request states that **implementation and canonical capture require separate authorization** |
| **C-8** | the request states that the decision may **decline**, and that a decline is a decision, not a failure |

**S2A-46 (Calculated — why C-6 is the pivotal condition).** If A-6 is imported, S-2 cannot consider a new kind or
field; P-4 is already barred by K-R14; only P-3/P-5 remain. That would make S-2 a decision about accepting the gap,
not about representing the change. Whether that is the owner's intent must be explicit, not inferred.

**S2A-47 (Calculated — the audit verdict).** **Passes with conditions C-1…C-8.** The next artifact **can** be a
narrowly scoped authorization request for S-2, not a carrier proposal. That request is **not drafted here**; drafting
it is a separate step, and it would still require the owner's approval.

---

## 8. Residual risks

**S2A-48 (Calculated — what could still go wrong).**

| Risk | Why it matters | Mitigation |
|---|---|---|
| **A-6 ambiguity** | importing it silently narrows S-2 to a gap-acceptance decision | C-6: state it explicitly |
| **entry-architecture slide** | shaping entries becomes fold behavior | C-5: S-2 decides shape, S-1 decides transitions |
| **ordering slide** | a class/scope order rule becomes N-c | C-5: label S-2's rule separate; N-c stays S-1's |
| **authority slide** | representation becomes permission | C-5: §11/C-1 R-E fence |
| **implementation slide** | a decision record includes DDL or code | C-2/C-7: logical-only; separate gate |
| **capture slide** | the decision adds a map row or edits canon | C-3/C-4: separate gates |
| **decline-as-denial** | a reasoned decline is read as "class/scope cannot change" | C-8: decline records the gap; Q2 = Yes stands |

**S2A-49 (Observed — no risk is resolved here).** The audit identifies the fences; it does not build them. The
authorization request would have to carry them.

---

## 9. What this record does not do

**S2A-50 (Observed).**

| Untouched | State |
|---|---|
| **S-2 decision** | not opened; not authorized |
| **carrier / schema / kind / field** | none selected or proposed |
| **C-3 / Q4 / Q5** | C-3 decomposed; Q5 owned by S-2; Q4 narrowed, not decided |
| **S-1 / N-a / N-b / N-c** | open; S-1 is downstream |
| **§11 / O-8a / v9 domain** | fenced, not decided |
| **A-6 / NSD-09** | not revised; A-6's applicability to S-2 is a condition (C-6) |
| **map / requirements** | **118**; no capture |
| **code / schema / migration / tests** | **0** |
| **owner request** | **not drafted** |

---

## 10. Traceability

**S2A-51 (Observed).**

| What | Where |
|---|---|
| K-R12: locus, domain, projection, non-mutation | `:111`; extract `:106` |
| K-R14: independent rebuild; no stored join | `:113` |
| §6.4: projection check; two streams, two checks | `:366`–`:367` |
| K-R2: edge `derived_from_event_ids`; integrity principle | `:83` |
| K-R13: read-time join; never stored | `:112` |
| K-R15: acquisition ordering only | §6.4 |
| §11: promotion per kind; retraction; §11.1 deletion | `:608`–`:619` |
| D1/D5: scopes; decay | `:70`–`:72` |
| K-R8/K-R9/K-R10/I13: class causes | `:222`–`:224`, `:299` |
| Scope values | §5 `:318` |
| r33 `scope` ≠ v9 domain; O-8a | RX-17; C3Q-05; CPA-05…07, CPA-30…36 |
| A-6 binds N-a, not S-2 | NAE-31/32; KRO-17 |
| NSD-09 binds ND-1, not S-2 | owner scope C-6; KRO-19 |
| Q2 = Yes; Q1 = S-2 | Q2A-33; Q1S-31 |
| Carrier families; P-4 excluded by K-R14 | C3D-10/11/13; KRO-22/23/25 |

---

## 11. Limits and record status

**S2A-52 (Observed — method).** **Newly read:** K-R2 (`:83`); K-R12–K-R14 (`:111`–`:113`); §6.4 (`:364`–`:370`);
§11 (`:600`–`:619`); D1/D5 (`:70`–`:72`); scope values (`:318`); C3D-10/11/13; Q1S-28…31; KRO-06/17/19/22/23/25.
**Not inspected:** any implementation's lifecycle handling, r32's bytes (absent), the frozen preregistration, v9 beyond
§7. 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound, `ExecutionNotAuthorised` at `conversation.py:177`.

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Carrier selected; S-2 opened | **0** |
| Authorization request drafted | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence** — no canonical capture, no ledger entry |

**End of record.** The audit **passes with conditions C-1…C-8** (S2A-45/47). S-2 inherits **K-R12** directly,
**K-R14** directly, **§6.4** directly, and **K-R2's integrity principle** through §6.4/K-R14 (its literal edge clause
does not bind the memory projection); it also inherits **K-R13** (no stored join) and **K-R15's acquisition-only
ordering** (S2A-08…S2A-20). It must not absorb **authority, validity, causes, N-b, N-c, N-a's fold, §11, O-8a or v9
`domain`** (S2A-38/39). It can be bounded to a **logical decision record** with separate gates for implementation and
canonical capture, **provided** the request states the subject, exclusions, A-6 status and decline option
(S2A-40…S2A-47). **A-6's applicability is the pivotal condition: importing it would reduce S-2 to a gap-acceptance
decision** (S2A-46). **The next artifact may be a narrowly scoped authorization request for S-2 — not a carrier
proposal — but it is not drafted here** (S2A-47/50). **Nothing is selected, opened, implemented or authorized.**
