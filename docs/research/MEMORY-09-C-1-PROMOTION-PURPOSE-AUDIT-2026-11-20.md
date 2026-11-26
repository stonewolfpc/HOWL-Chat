# MEMORY-09 — C-1 Promotion-Purpose Audit

## What job §11's promotion condition was doing, and which "promotion" O-8a governs

**Status:** audit record. **Research evidence only — no canonical edit, no schema, no code, no test, no probe, no
ledger entry, no implementation, no owner decision request.** It answers one question — what job the §11 promotion
condition performs — and reports one correction: **the corpus already distinguishes two promotions**, and the
authority framework governs the other one. No fold shape, entry architecture, reading of C-1 or owner of
record-promotion is selected. **C-2 and C-3 remain unopened.**

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | what job was §11's per-kind promotion condition intended to perform in the lifecycle? |
| Finding | it states **record-promotion criteria** — when a record of a kind may leave `provisional` — not **who may perform** the promotion (§3–§5) |
| Correction | **O-8a/O-9/O-10b/O-11 are confirmed and in force, but they govern v9 *authority-promotion*, a different operation**; they do not own §11's record-promotion (§6–§7) |
| Evidence already in corpus | **ODBC-19/20** and **ODB-50** name the two senses and require that they not be merged |
| Consequence for C-1 | R-E is strongly supported, but **not because O-8a supplies the permission**; the support is that §11's exits are evidence kinds with no actor, grant or check |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. The question, and the check it required

**CPA-01 (Observed — the task).** Trace the §11 promotion paragraph, §8.1's entry, the per-kind examples and the
earlier authority decisions; determine what §11 was for; and verify O-8a, O-9, O-10b and O-11 against the exact
current r33 wording, including whether any later decision superseded or narrowed them.

**CPA-02 (Observed — three things kept apart).** The record distinguishes **what the original author intended**,
**what the text establishes**, and **what later records inferred** (§8). None is treated as a substitute for another.

**CPA-03 (Observed — fresh verification).** The current r33 was read directly for this audit (`:600`–`:616`,
`:392`–`:413`, `:72`–`:75`, `:275`–`:290`, `:315`–`:345`, `:440`–`:465`), together with the owner ledger, v9, the
OD-b check and brief, the entry-uniformity record and the conceptual-object investigation. The earlier concern that
current r33 could not be retrieved does not apply to this reading.

**CPA-04 (Observed — no recommendation).** The audit reports what the evidence supports and where it stops. It does
not ask the owner to choose a reading, select a fold, or adopt a rule.

---

## 2. The decisive distinction: two promotions

**CPA-05 (Observed — the corpus already names it).** `MEMORY-09-OD-B-PROPOSED-ANSWER-CHECK-2026-11-19.md`
**ODBC-19**:

| Sense | Meaning | Source |
|---|---|---|
| **§11 record promotion** | a record leaves `provisional` by an evidence kind | r33 `:610`–`:612` |
| **v9 authority promotion** | an actor confers world/scenario authority, under a grant | v9 rules 24, 27–38 |

**ODBC-20** adds: "These are different operations with different authorities and different effects. The same word is
used for both, and §11's 'world fact — authority or corroboration' sits at the seam."

**CPA-06 (Observed — and it is restated as a boundary).** `MEMORY-09-OD-B-OWNER-DECISION-BRIEF-2026-11-19.md`
**ODB-50**: "'Record-promotion' (§11: a record leaves `provisional` by evidence kind) and 'authority-promotion' (v9
rules 24–38: an actor confers world/scenario authority) are **different operations** with different authorities. The
seam is §11's 'world fact — authority or corroboration'. **OD-b must not merge them.**"

**CPA-07 (Calculated — the two differ in object, cause, effect and owner).**

| | **Record promotion** (§11) | **Authority promotion** (O-8a / v9) |
|---|---|---|
| Object | the record's lifecycle **status** | which record **governs authoritative truth** |
| Cause | a per-kind **evidence condition** | an **authorised act** with a typed basis |
| Effect | `provisional → active` in the status projection | a separate world/scenario record is created; the source is unchanged |
| Actor named | **none** | the domain authority holder; delegation by grant |
| Basis vocabulary | witnessing, immediacy, authority/corroboration, attribution | `record_cited` / `actor_declaration` |
| Governing text | K-R12, §11 | O-8a-1…6, O-9, O-10b, O-11; v9 rules 17–21, 24, 28–38, 43–44 |

**CPA-08 (Inferred — they may compose, but composition is not identity).** An authority-promotion could be the event
that satisfies §11's "authority" condition for a world fact — the two could sit in one pipeline. That is a possible
**composition**, not a finding that the operations are one, and nothing in the corpus states the pipeline. The audit
records it as an open seam (§10), not as a bridge.

**CPA-09 (Narrowed — correction to the boundary record's wording).** `NAB-33` said "no current item owns promotion
permission." That is **Narrowed**: authority-promotion is owned by O-8a; **record-promotion's permission** is not
addressed by O-8a and remains unowned or inapplicable. The earlier record's conclusion about the §11 clause survives;
its phrase "promotion permission" was too wide.

---

## 3. The §11 paragraph and its surrounding model

**CPA-10 (Source-derived — the text).** §11 (`:600`–`:616`):

```
evidence turn
   └─> provisional (status, class, attribution set)
          ├─> active        (episodic, preference, validated fact, belief)
          ├─> superseded    (replaced; history kept)
          ├─> corrected     (earlier statement wrong; correction chain)
          └─> retracted     (asserted, now known false; acquisitions retained)
```

> "Promotion per kind: episodic — on witnessing; preference — immediate for holder; belief — never promotable to
> fact; rumour — stays provisional, repetition is not corroboration; world fact — authority or corroboration;
> absence — attributed. Rejected: blanket three-source rule, promotion by retrieval, decay altering truth, deleting
> history."

**CPA-11 (Observed — its place).** The paragraph sits under "Lifecycle and retraction", immediately after the state
diagram and immediately before "Retraction applies to records, not to minds (R4)". It is **untagged**; §11.1 and
§11.3 carry `[P]` (`:619`, `:687`). It is prose about transitions in the lifecycle, not about domains, grants or
actors.

**CPA-12 (Source-derived — the rejected list is about causes, not persons).** Three of the four rejected items —
blanket three-source rule, promotion by retrieval, decay altering truth — are **purported warrants or causes** that
the author refuses. None names an actor, a permission or a check. The fourth (deleting history) is a preservation
rule. The list's shape is *what may not count as a ground*, which is the shape of a promotion-criteria list.

**CPA-13 (Source-derived — the same refusal appears elsewhere).** D4: provisional records are "never promoted by
retrieval" (`:72`). Guard G2: "no C4→C2 promotion by retrieval, reranking or repetition" (`:449`). Both are evidence
conditions, not permissions.

**CPA-14 (Source-derived — what a record promotion changes).** The arrow's target is `active`, a **status value**
(`:322`), and K-R12 makes status a projection of the lifecycle stream (`:111`). The paragraph does not create a
record, transfer ownership or confer authority in a domain. It changes what the record's status is.

**CPA-15 (Source-derived — the presentation consequence).** Assembly may not "present provisional as verified"
(`:283`), and C4 carries "provisional + `via=told`" as "unverified reports" (`:446`). So the status distinction is
read at presentation: `provisional` is the unverified side. `active` is the side that may be presented as verified
for its kind. That is a **verification/standing** effect, not a permission effect.

**CPA-16 (Calculated — why "record promotion" is the right name).** Every element of the paragraph — the source
state, the target state, the per-kind conditions, the rejected grounds — concerns the record's movement between
statuses. None concerns which record governs a domain or who holds a grant. The paragraph is a **transition policy
for the status stream**.

---

## 4. The per-kind trace

**CPA-17 (Calculated — each condition, what it ranges over, and whether any permission appears).**

| Kind | §11 condition | What the condition ranges over | Actor / grant / domain named? | What it changes |
|---|---|---|---|---|
| **episodic** | "on witnessing" | an acquisition route (`via=witnessed`, `:345`); an event | **no** | the record may leave `provisional` |
| **preference** | "immediate for holder" | timing and the holder relation | **no** | leaves at once for its holder |
| **belief** | "never promotable to fact" | the **target** — belief cannot become fact | **no** | a prohibition on the target, not on an actor |
| **rumour** | "stays provisional; repetition is not corroboration" | terminality for the kind; a false warrant refused | **no** | no exit; repetition does not count |
| **world fact** | "authority or corroboration" | the claim's support — source standing or independent support | **no** | the record may leave `provisional` |
| **absence** | "attributed" | the attribution set written at formation (`:604`) | **no** | the record carries its attribution |

**CPA-18 (Observed — the list is heterogeneous, and that is the point).** The six entries are not one logical type:
a cause (witnessing), a timing rule (immediate), a target prohibition (belief), a terminality rule (rumour), a
warrant pair (world fact), and an attribution rule (absence). What unites them is not "permission" but **the
condition under which this kind may leave `provisional`**.

**CPA-19 (Calculated — no permission vocabulary appears anywhere in the list).** Not one entry names who may act, a
grant, a domain, a check, or a reviewer. A permission list with no permission-holder in any of six entries would be
an odd construction; a promotion-criteria list with no actor is ordinary.

**CPA-20 (Inferred — the author's discriminative purpose).** The list separates kinds that **can** be promoted
(episodic, preference, world fact, absence) from one that **must not** be promoted to fact (belief) and one that
**never** leaves (rumour), and it rejects false grounds (repetition, retrieval, a blanket source count). That is a
purpose about **eligibility and warrant**, not about **authorisation**.

**CPA-21 (Calculated — world fact against the others).** "Authority or corroboration" is the only entry with two
alternatives. Both alternatives describe **support for the claim**: an authoritative source, or independent
corroboration. The neighbouring entries are also support/condition descriptions. Nothing in the row introduces an
actor who was previously absent from the list.

**CPA-22 (Unestablished — what qualifies under either alternative).** What counts as an authoritative source, and
what counts as sufficient corroboration, are **not stated**. That is the unresolved threshold the owner identified,
and the audit does not supply it (CPA-38).

---

## 5. §8.1, entry, and the difference between becoming `active` and being permitted to act

**CPA-23 (Source-derived — the collision).** §8.1 step 1: "Formation writes `scope=world`,
`knowledge_class=experiential`, **`status=active`**, `established_turn=40`" (`:404`). §11 says formation reaches
`provisional` and promotion reaches `active` (`:604`). Both are untagged. ECI-05/ECI-06 record the collision.

**CPA-24 (Observed — the reconciliation is an inference, not a rule).** ECI-07 reads the two as compatible if the
exit condition is satisfied **at formation**, so §8.1 records the post-advance value. ECI-08 reads authored world
facts as born `active`. ECI-09: canon does not choose. Neither reading is adopted here.

**CPA-25 (Observed — the missing premise already named).** EKR-11 calls the premise **P-EVAL**: exit conditions are
evaluated at the formation act. EKR-11: P-EVAL "is **not stated anywhere in r33** — it is inferred from the per-kind
wording only." EKR-12 adds that P-EVAL would make formation an act that both creates and advances, so the authority
to form and the authority to advance would be exercised together.

**CPA-26 (Source-derived — and the existing record already reads authority as attribution at formation).** EKR-13's
world-fact row: "authority is present at formation **if the attribution set records it**; corroboration needs later
events." That is an epistemic reading of authority — a property recorded in the attribution set, not a permission
held by an actor. It is a **later record's inference**, but it is the corpus's own reading of this exact clause.

**CPA-27 (Calculated — becoming `active` and being permitted to act are different questions).** "This record is
`active`" answers *what status the projection yields*. "This actor may promote it" answers *whether an act is
permitted*. §11's conditions are stated as properties of the record and its support; O-8a's conditions are stated as
properties of an actor and a domain. The two questions can be asked of the same event without being the same
question.

**CPA-28 (Calculated — a warrant does not need an actor, but a write still needs a writer).** If record promotion is
evidence-driven, the semantic condition can be satisfied without naming who caused it. That does not answer whether
the **write** that records the transition requires a permission — a separate, implementation-facing question. The
audit keeps semantic warrant and write permission apart (§10).

**CPA-29 (Inferred — §8.1 does not discriminate between the readings).** EKR-16: the §8.1 record is `scope=world`
**and** `knowledge_class=experiential`, so `status=active` is equally consistent with "episodic, promoted on
witnessing" and "world fact, promoted on authority". §8.1 corroborates immediate exit without identifying which
kind's rule did the work.

---

## 6. The authority framework, verified against current text

**CPA-30 (Observed — O-8a is real, confirmed and closed).** The ledger records **O-8a — promotion authority (ALL SIX
CONFIRMED — series closed)**. Its six terms: 8a-1 holder (world/scenario authoring authority, human operator by
default, per domain); 8a-2 basis (typed set, `record_cited`/`actor_declaration`, no empty set, no relevance or
coverage test); 8a-3 per-domain scope; 8a-4 explicit revocation/supersession, human-controlled in execution; 8a-5
shared regime with Q7 ambient authorship, operations distinct; 8a-6 standing model grants, explicit, scoped,
recorded, revocable.

**CPA-31 (Observed — the other three).** **O-9** (narrowly answered): disclosure may be cited as a basis, never
promotes by itself; "Evidence and authority are different things; sufficiency of a basis folds into O-8a."
**O-10b** (confirmed): no actor may unilaterally revoke or supersede another's promotion, or their own prior one;
a human edit is required. **O-11** (confirmed): promotion creates a **separate world/scenario record**; the source
character record is unchanged; promotion changes which record governs authoritative truth, not whether the
disagreement occurred.

**CPA-32 (Observed — what the framework says promotion changes).** O-8b: promotion changes **authority**, not
ownership, knowledge or belief. v9 rules 17–21: promotion changes authority; creates a world/scenario record;
confers no knowledge; compels no belief; preserves history; "authority, not history." v9 rule 24: promotion authority
is separate; generation never confers it. Rules 28–38: typed basis, citations, no relevance test, per-domain,
explicit revocation/supersession, authority changes never rewrite character state.

**CPA-33 (Observed — no later supersession or narrowing of the framework was found).** v9's confirmed inputs list
includes **O-11, O-8a (all six), O-10, O-10b** (`:15`). The ledger's consolidated O-8a record closes the series.
The only qualifications are recorded, not silent: **8a-4's execution** is human-controlled by O-10b (correction 22;
the wording is retained and the qualification is recorded alongside), and **RAT-1** gives O-10b tier-1 precedence
over K10's "latest governs." No decision supersedes O-8a, O-9, O-10b or O-11, and none narrows their **scope**.

**CPA-34 (Observed — but the scope is authority-promotion, and the records say so).** O-11's object is a separate
world/scenario record; O-8b's effect is authority; v9 rule 17 is "promotion changes authority... it creates a
world/scenario record"; O-8a-2's basis kinds are record citations and actor declarations — not witnessing,
immediacy, corroboration or attribution. **The framework's vocabulary does not match §11's per-kind conditions.**
And `MEMORY-09-OWNER-DECISION-DN1-DN4-2026-11-19.md` **OD-20** states the limit directly: for retraction of
claims/accounts, existing authority is "**none** for claims/accounts; **O-10b governs *promotions* only**."

**CPA-35 (Calculated — the framework is confirmed but not applicable to §11 by its own terms).** O-8a answers "who
may confer world/scenario authority, on what typed basis, in which domain." §11 answers "when does a record of this
kind leave `provisional`." The first is an actor-and-domain question; the second is a record-and-support question.
Reading O-8a into §11 would **expand** O-8a's scope, not apply it — the expansion the owner warned against, in the
opposite direction.

**CPA-36 (Observed — the earlier boundary record's error, stated precisely).** The boundary record treated
"promotion" as one thing and therefore reported promotion permission as unowned (NAB-33). The corpus had already
distinguished the senses (ODBC-19/20; ODB-50), and the boundary record did not cite them. **The conclusion about
§11 survives; the unqualified phrase does not** (CPA-09).

---

## 7. The check the owner asked for

**CPA-37 (Calculated — the four-way trace, stated as results).**

| Check | Result |
|---|---|
| Does **O-8a** govern who may promote? | **Yes — authority-promotion**; per-domain, basis-typed, delegable, human-controlled changes |
| Does it govern §11's `provisional → active`? | **No** on its own terms: different object, cause, effect and vocabulary (CPA-07, CPA-34) |
| Do **O-9 / O-10b / O-11** reach §11? | **No**: O-9 separates evidence from authority *within authority-promotion*; O-10b governs authority changes; O-11 describes the separate authoritative record |
| Superseded or narrowed later? | **No supersession or scope-narrowing found**; 8a-4 is qualified in execution by O-10b; RAT-1 is a precedence application |
| Does this make R-E the better reading? | **Yes, but on §11's own evidence** — exits are evidence kinds, no actor/grant/domain appears, and the list's function is per-kind eligibility (CPA-17…CPA-22) |

**CPA-38 (Calculated — the corrected simple conclusion).** The conclusion the owner proposed needs one qualifier:

> §11's clause concerns **epistemic support for record promotion** — the standing of a world fact's source, or
> independent corroboration — while **actor permission for *authority-promotion*** is governed separately by O-8a.
> **Actor permission for §11 record promotion is not governed by O-8a**, and whether it exists at all is
> unestablished.

The first clause is supported; the second is true of a different operation; the third is the part the cross-record
finding would have erased.

**CPA-39 (Calculated — why the distinction matters for C-1).** If O-8a were read as owning §11's permission, C-1
would collapse into "which framework applies," and R-D would look already answered. Because the operations differ,
C-1 remains a genuine question about the **clause's referent**, and the evidence now points to R-E for reasons
internal to §11, not by elimination.

**CPA-40 (Unestablished — the residual seam).** Whether an authority-promotion event can **satisfy** §11's
"authority" condition for a world fact — a composition of the two operations — is not stated. It is the most
plausible bridge between the layers, and it is exactly what ODBC-20 called "the seam." The audit does not adopt it.

---

## 8. Intended, established, inferred

**CPA-41 (Observed — the three registers, kept apart).**

| Register | What can be said | Status |
|---|---|---|
| **Original author's intent** | the paragraph is a per-kind promotion policy, placed in the lifecycle section, refusing false grounds (repetition, retrieval, blanket source counts) and separating kinds that may advance from those that may not | **Inferred** from placement, wording and the rejected list; no author statement exists (CLV-09: the promotion prose has one witness) |
| **What the text establishes** | per-kind conditions exist; no actor, grant, domain or check is named; `active` is a status; status is a projection; the paragraph is untagged; no event kind exists | **Source-derived / Observed** |
| **What later records inferred** | that the exits are "evidence kinds" (ODBC-10/11); that there are two promotions (ODBC-19/20; ODB-50); that authority is present at formation via attribution (EKR-13); that the clause touched excluded territory (NAE-26); that promotion permission was unowned (NAB-33); that O-8a covers promotion permission (the present correction) | **Inferred**, each with its own record |

**CPA-42 (Calculated — the registers do not agree by construction).** The author's purpose is an inference; the
text's content is narrower; later records add readings the text does not state. Treating "the corpus says" as one
voice would hide the difference — which is how the boundary record came to use "promotion" unqualified.

---

## 9. What this audit does not do

**CPA-43 (Observed).**

| Untouched | State |
|---|---|
| **C-1** | **not answered**; the audit strengthens one reading's evidential basis without selecting it |
| **C-2 / C-3** | **unopened** |
| **N-a / F-1…F-6** | unanswered; no shape selected, ranked or excluded; **F-6 is not a default** |
| **I-1 / I-2 / I-3** | I-3 excluded by K-R12; the others open |
| **N-b / N-c / N-d** | untouched |
| **O-8a / O-9 / O-10b / O-11** | **not reopened, not narrowed, not extended**; confirmed status recorded |
| **DN-6a / DN-6b** | untouched; still not promotion owners |
| **DN-1; DQ-4; validity predicate** | settled / open / absent |
| **§11's adoption status; §8.1's collision** | not decided |
| **map / requirements** | **118**; DQ-3's requirement remains the acknowledged capture gap |

**CPA-44 (Observed — no rule, predicate, schema or capture).** The audit proposes no authority rule, validity
predicate, event kind, field, key, fold or entry architecture, and creates no canonical record or ledger entry.

---

## 10. What remains unresolved

**CPA-45 (Observed — the open list, stated rather than implied).**

| # | Open question | Why it survives |
|---|---|---|
| 1 | What qualifies as an **authoritative source** | §11 names the alternative, not its threshold |
| 2 | What constitutes **sufficient corroboration** | same |
| 3 | Whether **§11 is operative** or merely descriptive | the paragraph is untagged |
| 4 | Whether **record promotion needs a permission** at all, and at which layer (write vs semantics) | O-8a does not address it; CPA-28 separates the layers |
| 5 | Whether an **authority-promotion event satisfies** §11's authority condition | the seam ODBC-20 names; CPA-40 |
| 6 | Whether **P-EVAL** holds (exits evaluated at formation) | EKR-11: inferred, not stated |
| 7 | Whether **§8.1** or **§11** governs the entry value | ECI-07/08/09 |

---

## 11. Traceability

**CPA-46 (Observed).**

| What | Where |
|---|---|
| The two promotions, named | ODBC-19, ODBC-20; ODB-50 |
| §11's text, diagram, rejected list | r33 `:600`–`:616`; CLV-09 |
| §11's exits as evidence kinds | ODBC-10, ODBC-11 |
| Provisional as entry, promotion as exit | ODBC-04, ODBC-10, ODBC-11 |
| Presentation consequence of `provisional` | r33 `:283`; C4 `:446` |
| D4 and G2 refuse promotion by retrieval | r33 `:72`, `:449` |
| §8.1 writes `status=active` | r33 `:404`; ECI-05 |
| The collision and its readings | ECI-06…ECI-09 |
| P-EVAL | EKR-11, EKR-12 |
| Authority present at formation via attribution | EKR-13 |
| O-8a (six terms) | ledger `:85`–`:180` |
| O-9 | ledger `:76`–`:80` |
| O-11 | ledger `:60`–`:73` |
| O-10b | ledger `:303`–`:327` |
| O-8b | ledger `:72`–`:75` |
| Authority-promotion's effects | v9 rules 17–21, 24, 28–38, 43–44 |
| O-10b governs promotions only | OD-20 (DN-4 decision record) |
| 8a-4 qualification; RAT-1 precedence | correction 22; v9 confirmation table; ARCH-CONTRACTS RAT-1 |
| Confirmed inputs include O-8a/O-10b/O-11 | v9 `:15` |
| The framework's map rows | K21–K26, L7 |
| The earlier unqualified phrase | NAB-33; CPA-09 |

---

## 12. Limits and record status

**CPA-47 (Observed — method).** **Newly read:** r33 `:600`–`:616`, `:392`–`:413`, `:72`–`:75`, `:275`–`:290`,
`:315`–`:345`, `:440`–`:465`; the owner ledger's O-11/O-9/O-8a/O-10b/O-8b sections; v9 rules 17–21, 24, 28–38,
43–44 and the confirmed-inputs table; ODBC-19/20; ODB-50; EKR-11/12/13/16; ECI-05…ECI-09; the DN-4 decision's OD-20;
the conceptual-object investigation's RX-28. **Not inspected:** any implementation's promotion path, r32's bytes,
the frozen preregistration, v9 beyond §7, the K13–K17 rows.

**CPA-48 (Observed — no experiment).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. The audit is textual and historical.

**CPA-49 (Observed — the labels' limits).** The author's intent (CPA-20, CPA-41) is **Inferred**; the two-promotion
distinction (CPA-05/06) is **Observed** in the corpus; the applicability conclusion (CPA-35) is **Calculated** from
the operations' stated objects; the residual seam (CPA-40) is **Unestablished**. CPA-09 is **Narrowed** and does not
rewrite NAB-33, which stands as the state at its date (OD-06).

**CPA-50 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Reading selected, rule or predicate created | **0** |
| O-8a / O-9 / O-10b / O-11 reopened or narrowed | **0** |
| C-1 answered; C-2/C-3 opened | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence** — no canonical capture, no ledger entry |

**End of record.** **The §11 promotion condition's job is to state, per memory kind, when a record may leave
`provisional`** — a **record-promotion criterion** whose conditions are causes, timing rules, prohibitions, a
terminality rule, a warrant pair and an attribution rule, and which names **no actor, grant, domain or check**
(CPA-17…CPA-22). **O-8a, O-9, O-10b and O-11 are confirmed and unsuperseded**, but they govern **v9
authority-promotion** — a different operation with a different object, cause, effect and vocabulary (CPA-05…CPA-07,
CPA-30…CPA-34). **The corpus already says so** (ODBC-19/20; ODB-50), and the earlier boundary record used
"promotion" without that qualifier (CPA-09, **Narrowed**). **R-E is strongly supported — by §11's own evidence, not
by O-8a's existence** (CPA-37/38). The corrected conclusion is: *§11 concerns epistemic support for record
promotion; actor permission for authority-promotion is governed separately by O-8a; actor permission for record
promotion is not governed by O-8a and is unestablished.* **The seam — whether an authority-promotion can satisfy
§11's condition — remains open** (CPA-40). **Nothing is selected, requested, designed or authorised.**
