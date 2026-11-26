# MEMORY-09 — N-a Evidence and Decision Readiness

## Fold semantics: what canon supplies, what it does not, and the alternatives that remain

**Status:** evidence and decision-readiness brief. **Research evidence only — not canonical policy, not a requirement, not
a selection.** Read-only: no canonical edit, no schema, no code, no test, no probe, no ledger entry, no implementation.
**N-a is not answered, opened as work, or formulated as an owner request.** No fold rule is proposed. **N-b** and **N-c**
are untouched except for N-c's dependence on N-a's answer. **N-d** is excluded by the scope determination and stays with
DN-6a/DN-6b. **DN-1** remains separate and settled.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Part | **N-a** — fold semantics: what a later lifecycle event does to state (NSD-03) |
| Answers | what canon states (§2), what it leaves open (§3), the genuine alternatives (§4), and the pressure points an examination would have to face (§5) |
| Does not do | choose an alternative, draft a decision, propose a fold rule, or enable N-c |
| Boundary | post-formation lifecycle events only — **entry/promotion is a pressure point, not in scope** (§5.2) |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The part, and its boundary

**NAE-01 (Observed — the part).** N-a asks **what a later lifecycle event does to state** — the transition function over
K-R12's kinds (A2-41; NDR-13). Its subject exists in canon as a *requirement that state be a projection*, and not as a
rule for how the projection is computed.

**NAE-02 (Calculated — what N-a is not).** Not N-b (how `established` relates to `active` — though the two touch at the
entry point), not N-c (whether an order key exists), not N-d (who may cause a transition), not validity (whether a
transition was admissible), not DQ-4, and not any carrier.

**NAE-03 (Observed — the dependency the scope determination preserved).** N-c's **relevance** depends on N-a's answer: if
the fold proves order-insensitive, an ordering key has no subject; if order-sensitive, it has one (NSD-12; DS-09). This
brief therefore touches order-sensitivity as a **property of each alternative**, never as a decision.

**NAE-04 (Observed — no recommendation).** The alternatives in §4 are stated with what each would require and imply. None
is ranked, none is described as more faithful, and no owner request is drafted.

---

## 2. What canon states

**NAE-05 (Source-derived — the projection requirement).** **K-R12** (`:111`): memory lifecycle is its **own authoritative
stream** — `established | superseded | corrected | retracted`, plus `scope` and `knowledge_class` changes — and "a memory's
current status, class and scope are a **projection of that stream**, not mutable authoritative fields."

**NAE-06 (Source-derived — the domain is three dimensions, not one).** K-R12 names **status, class and scope** as
projections of the same stream. So "fold semantics" has a domain wider than the `status` value: scope and knowledge-class
changes are inputs to the same projection. Whether N-a covers all three or only `status` is a boundary question this brief
raises and does not answer (NAE-24).

**NAE-07 (Source-derived — isolation).** `:367`: "Two streams, two checks. Acquisition events project to edges; memory
events project to the memory's status/scope/class. **Neither projection may consult the other.**" So the lifecycle fold may
not read acquisition events — the same purity discipline K-R11 imposes on the edge fold, in the other direction.

**NAE-08 (Source-derived — independent rebuild and integrity).** **K-R14**: each stream is rebuilt and checked
independently; a mismatch in either is an **integrity failure**. `:366`: "fold events → projection; compare to stored
projection; mismatch is an integrity error surfaced to `inspection`, never silently resolved in favour of either side."
So the fold must be **deterministic and checkable** — the requirement is strong even though the rule is absent.

**NAE-09 (Source-derived — the determinism pattern, stated for acquisition).** `:371`: "replaying the same ordered log
produces the same edge, on any machine, in any insertion order; and a missing edge is not 'unknown state' to be guessed —
it is a determinate consequence of the **last event**." This is the house pattern: ordered log, deterministic replay,
last-event determination. It is written for the edge fold; its application to the lifecycle fold is **Inferred**, not
stated.

**NAE-10 (Source-derived — the lifecycle picture in §11).** §11 (`:600`–`:616`) gives a diagram, not a rule:

```
evidence turn
   └─> provisional (status, class, attribution set)
          ├─> active        (episodic, preference, validated fact, belief)
          ├─> superseded    (replaced; history kept)
          ├─> corrected     (earlier statement wrong; correction chain)
          └─> retracted     (asserted, now known false; acquisitions retained)
```

with "Promotion per kind: episodic — on witnessing; preference — immediate for holder; belief — never promotable to fact;
rumour — stays provisional, repetition is not corroboration; world fact — authority or corroboration; absence —
attributed." Formation is at `:604` (ODBC-16/17).

**NAE-11 (Source-derived — the record's fields).** `:341`: the record carries `attribution{kind, by}`, `status`,
`established_turn`, `supersedes` / `superseded_by`. So supersession may be **relational** — a chain — and not only a
status value, which a fold over `status` alone would not produce.

**NAE-12 (Source-derived — the two vocabularies, side by side).** `:322` lists the `status` field's values as
`provisional | active | superseded | retracted | corrected`; K-R12's event kinds are `established | superseded | corrected
| retracted`. The two lists are **not** the same list: `established` is an event with no status of that name, `active` is a
status with no event of that name, and `provisional` is a status whose producing event is `established` only by
inference.

**NAE-13 (Observed — K-R17 as the only written transition-validity pattern).** K-R17 (`:117`) makes `revoked` valid only
when the pair's ordered log "currently ends in a learn event", and requires rejection at write rather than a silent no-op.
It is the house's only explicit statement of the form *an event is valid only from a state* — and it is written for the
**acquisition** stream.

---

## 3. What canon does not state

**NAE-14 (Observed — the absences, each already recorded).**

| Absent | Recorded by |
|---|---|
| a **fold rule** over lifecycle events | A2-41 ("ND-1, unspecified"); DS-09 ("no order key and no fold rule") |
| an **admissible-transition set** per kind | LCM-24; LCM-40 |
| a **producer for `active`** in K-R12's kind list | A2-42; NAE-12 |
| any **conflict resolution** — two lifecycle events whose kinds disagree | not recorded anywhere; **observed** as an absence for this brief |
| any **precedence or terminality** among kinds | not stated |
| a **`status` fold vs a chain computation** (supersession links) | not stated; `:341` implies the links exist |
| the **domain** — status only, or status+class+scope | K-R12 names all three; no rule says the fold treats them alike |

**NAE-15 (Observed — the requirement is stated, the computation is not).** NAE-05…NAE-09 give a projection requirement, a
determinism requirement and an integrity check. None of them says what the projection **is**. That is a stronger form of
absence than a missing detail: the system is required to compute something no text defines.

**NAE-16 (Calculated — and the check cannot substitute for the rule).** The projection check at `:366` compares a
computed projection to a stored one. With no rule, the comparison has no independent content: two implementations could
disagree and each be internally consistent. The check detects drift, not disagreement about meaning.

---

## 4. The genuine alternatives

**NAE-17 (Calculated — six shapes a fold rule could take, stated as properties, none recommended).**

| # | Shape | What it would make the current status | Order-sensitive? | What it would require |
|---|---|---|---|---|
| **F-1** | **last-event-wins** | the status named by the most recent lifecycle event, in stream order | **yes** | an order key — **N-c becomes operative** |
| **F-2** | **kind precedence** | the status of the highest-precedence kind present, regardless of arrival order | **no** | a precedence order among kinds — a new normative item no text states |
| **F-3** | **state-transition function** | the result of applying each event only if valid from the current state (K-R17's pattern) | **yes** | an admissible-transition set — absent (LCM-24) |
| **F-4** | **terminal states** | a terminal kind (e.g. `retracted`) absorbs all later events | partly — terminality is a precedence special case | a terminality list; interacts with F-2 |
| **F-5** | **set-valued projection** | the set of lifecycle facts, with readers interpreting `status` | **no** in principle | a departure from the single-valued `status` field (`:322`) — a schema consequence, not authorised |
| **F-6** | **no rule** | whatever the implementation computes; drift is the only check (NAE-16) | unspecified | nothing — the status quo |

**NAE-18 (Calculated — F-1 and F-2 are the two that fit the existing text without new kinds).** F-1 matches the
determinism pattern at `:371` ("consequence of the last event") and needs no new event kinds; its cost is a dependency on
ordering. F-2 matches the intuition that a retraction should not be undone by a later correction, but it requires a
precedence order that nothing states.

**NAE-19 (Calculated — F-3 imports a validity condition into the fold).** K-R17's pattern makes validity part of the
transition rule. Applied to lifecycle, F-3 would make the fold reject or ignore invalid transitions — which requires the
admissible set (absent) and, at the margins, the authority question N-d holds. So F-3 is the alternative most likely to
cross N-a's boundary.

**NAE-20 (Calculated — F-4 and F-2 are not independent).** Terminality is precedence with a specific justification. A
decision that chose either would have to say whether the other is a special case or a rival, and the brief does not
resolve that.

**NAE-21 (Calculated — F-5 has a schema consequence, which is why it is listed and not developed).** A set-valued
projection would change what `status` means as a field (`:322`, `:341`). Naming it is necessary for completeness; pursuing
it would be a schema proposal, which no part of ND-1 authorises (NSD-09).

**NAE-22 (Observed — F-6 is not a recommendation and not a null case).** "No rule" describes the present arrangement and
its cost: the projection is required, checkable, and undefined (NAE-15, NAE-16).

**NAE-23 (Calculated — what every alternative shares).** All six accept K-R12's projection requirement and the isolation
at `:367`; none changes the stream's authority, the record's fields or the two-stream split. They differ only in how the
projection is computed.

---

## 5. Pressure points an examination would have to face

**NAE-24 (Calculated — the domain).** K-R12 projects **status, class and scope** from one stream (NAE-06). A fold rule for
`status` alone would leave class and scope changes unmodelled, while a rule for all three would have to say whether the
three are computed independently or jointly. The scope determination named N-a as "fold semantics" without fixing this;
the brief records it as the first boundary question rather than resolving it.

**NAE-25 (Calculated — the entry point).** §11's diagram starts at `provisional` and promotes to `active` per kind, and
`:322` gives `provisional` and `active` as status values. K-R12's kinds contain **no `active`-producing event** (NAE-12).
So a fold over lifecycle events either (a) must derive `active` from `established` plus the per-kind promotion conditions,
(b) needs an event kind that does not exist, or (c) treats `active` as outside the fold. The brief records the three and
chooses none.

**NAE-26 (Calculated — and the promotion conditions include an authority component).** "world fact — authority or
corroboration" (§11) is a condition about **authority**, which the scope determination excluded from ND-1 and left with
DN-6a/DN-6b. So N-a may not be fully separable from the excluded territory **at the entry transition**, even though it is
separable for post-formation events. This is a genuine pressure on the clean boundary, recorded rather than smoothed.

**NAE-27 (Unestablished — whether N-a can be scoped to post-formation events only).** If the fold begins **after** a record
has a status, the entry question drops out and NAE-26 does not arise. Whether that restriction is coherent — i.e. whether
the first event is not itself a lifecycle event — is **not established**, and this brief does not assert it.

**NAE-28 (Calculated — the supersession chain).** `:341`'s `supersedes`/`superseded_by` are relations, not statuses. A
status-only fold would leave them uncomputed, and a chain-computing fold would have to say whether the chain is derived
from events or declared on them. Neither is stated.

**NAE-29 (Calculated — order-sensitivity is the pivot for N-c, not a decision here).** F-1, F-3 and F-4 are
order-sensitive; F-2 and F-5 are not in principle. So N-a's answer determines whether N-c has a subject (NAE-03), and the
brief deliberately does not predict which way that falls.

**NAE-30 (Observed — no validity predicate is required by any alternative).** None of F-1…F-6 supplies a condition under
which a transition was *warranted*; the validity gap stays where it is (GRP-31; CEV-24). F-3's admissibility is
state-relative validity of an event kind, not the validity of the act that produced it — a distinction that must survive
into any later work.

---

## 6. What a later decision would have to fix

**NAE-31 (Calculated — stated as terms, not as a request).**

| # | Term | Why explicit |
|---|---|---|
| **A-1** | the **shape** (F-1…F-6), or an explicit refusal to fix one | the alternatives are not variants of each other |
| **A-2** | the **domain** — status only, or status+class+scope (NAE-24) | K-R12 projects three dimensions |
| **A-3** | the **entry point** — post-formation only, or including `active`'s production (NAE-25) | the fold's first input decides whether authority enters |
| **A-4** | how **supersession links** are obtained (NAE-28) | status alone does not produce them |
| **A-5** | whether the rule is **order-sensitive**, and therefore whether N-c is operative (NAE-29) | N-c's relevance is N-a's consequence |
| **A-6** | that the decision **creates no event kind, field, key or schema**, and does not settle authority or validity | the standing boundaries (NSD-09; GRP-31) |

**NAE-32 (Observed — A-6 is the term that keeps a fold decision from becoming a schema decision).** A-1 through A-5 can
all be stated as rules over existing kinds and existing fields. A-6 says so explicitly, because "the fold produces `active`"
can be read as licensing a new event type (NAE-25(b)).

---

## 7. What this brief does not do

**NAE-33 (Observed).**

| Untouched | State |
|---|---|
| **N-a** | **not answered**; not opened as work; no owner request drafted |
| **N-b** | untouched — except that `active`'s entry point borders it (NAE-25) |
| **N-c** | untouched; its relevance is a consequence, not a decision (NAE-29) |
| **N-d**; **DN-6a/DN-6b** | excluded / open; unchanged |
| **DN-1** | settled; not reopened |
| **DQ-4**; the validity predicate; the admissible set | open / absent |
| **axis H** | unchanged (GRP-31) |
| **QAD-19**; the report-only boundary; the prohibition; DQ-3's A; D-3b; ODC-07 | unchanged |
| ODL-39; closure; determinism; §8.1; P-EVAL; C′; OD-g; Q-B | unchanged |

**NAE-34 (Observed — no schema, key, event or field is implied).** The brief names `status`, `class`, `scope`,
`supersedes`/`superseded_by` and K-R12's kinds because they are the evidence. It proposes no event type, no key, no field
addition and no migration.

---

## 8. Traceability

**NAE-35 (Observed).**

| What | Where |
|---|---|
| The part and its boundary | NSD-03; NDR-13 |
| N-c's dependence on N-a's **answer** | NSD-12; DS-09; NAE-03 |
| The projection requirement; three dimensions | K-R12 (`:111`) |
| Isolation between the two projections | `:367`; K-R11 |
| Independent rebuild; integrity failure; projection check | K-R14; `:366` |
| Determinism and "the last event" | `:371`; K-R15; K-R20 |
| The lifecycle picture and promotion per kind | §11 `:600`–`:616`; formation `:604` |
| Status values vs event kinds | `:322`; NAE-12 |
| Record fields, including supersession links | `:341` |
| The only written transition-validity pattern | K-R17 (`:117`) |
| No fold rule; no order key | A2-41; DS-09 |
| No admissible-transition set | LCM-24; LCM-40 |
| Axis H is unchanged by any fold | GRP-31; NAE-30 |

---

## 9. Limits and record status

**NAE-36 (Observed — method).** **Newly read for this brief:** r33 `:90`–`:145` (K-R11…K-R20), `:300`–`:375`
(the two-stream tables, record fields, projection check, fold determinism), and §11 `:600`–`:616`; plus A2-41…A2-45,
DS-09, LCM-24/LCM-40, and v9's only status reference (`:212`). **Not inspected:** any implementation's fold, r32's bytes,
the frozen preregistration, v9 beyond §7, and the K13–K17 rows. No test, probe or model call was run (NAE-37).

**NAE-37 (Observed — no experiment).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. The brief concerns what text says and what it omits.

**NAE-38 (Observed — the labels' limits).** NAE-02, NAE-16, NAE-17…NAE-23, NAE-24…NAE-26, NAE-28, NAE-29, NAE-31 and
NAE-32 are **Calculated**. NAE-27 and the "whether the chain is derived or declared" clause of NAE-28 are
**Unestablished** on purpose. Nothing here is offered as a provision of r33 or v9 beyond the quotations given.

**NAE-39 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Alternatives chosen or ranked | **0** |
| Fold rule, key, event, field or schema proposed | **0** |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence** — no canonical capture, no ledger entry |

**End of record.** **N-a's evidence is thin and specific:** K-R12 requires the projection without defining it (NAE-05),
`:367` isolates it, K-R14 and `:366` make it checkable, `:371` supplies a determinism pattern for the *other* stream, §11
supplies a picture and per-kind promotion rather than a rule, and no text states conflict resolution, precedence,
terminality, or a producer for `active` (NAE-10…NAE-15). **Six shapes remain available** (NAE-17), differing in
order-sensitivity — which is exactly what decides whether **N-c** has a subject (NAE-29). **Three pressure points** are
recorded rather than smoothed: the **domain** (status vs status+class+scope), the **entry point** where `active` must be
produced, and the **authority component in the promotion conditions** that touches the excluded territory (NAE-24…NAE-26).
**Nothing is chosen, requested, designed or authorised.**
