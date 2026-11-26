# MEMORY-09 — ND-1 Decision Readiness

## Four candidate contents, and why the question must be scoped before it is asked

**Status:** decision-readiness brief. **Research evidence only — not canonical policy, not a requirement, not a
selection.** Read-only: no canonical edit, no schema, no code, no test, no probe, no ledger entry, no implementation.
**ND-1 is not answered, narrowed or opened as work**; this brief establishes what it has been taken to ask, what each
reading would need, and what is missing. The settled set — Q-α, the cause classification, Q-β = A, DQ-2, DQ-3 — is taken
as given and not reopened. **DN-1** (settled), **DN-6a/DN-6b** (open) and **OD-a**/**OD-b = L2** are kept distinct.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | what does **ND-1** ask, and what evidence or definitions are missing before an owner decision could be put? |
| Finding | the records use **ND-1** for **two different contents** — status fold/vocabulary semantics and transition authority — and sometimes for both (NDR-24) |
| Consequence | a decision cannot be framed until **which content** ND-1 names is fixed; the ND-1 → DQ-4 dependency holds only under the authority reading (NDR-26) |
| Does not do | scope it, answer it, design a fold, order key, event type, mapping or authority rule |
| Kept distinct | **DN-1** = v9 §8 textual clarification, confirmed 2026-11-19 · **DN-6a/DN-6b** = who may forget / retract, open · **`active`/`established`** (ODL-29), preserved |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded** ·
**Narrowed**.

---

## 1. What this brief is for, and what it refuses to do

**NDR-01 (Observed — the task).** To establish exactly what ND-1 asks, and what evidence or definitions are missing,
**without settling its implementation**.

**NDR-02 (Calculated — the order matters).** A decision request needs a subject before it needs options. ND-1's subject is
currently ambiguous in the records (§2), so the first product of this brief is a **scope question**, not an answer space.
Putting an answer space first would let the ambiguity be resolved by whichever options happened to be drafted.

**NDR-03 (Observed — what is not proposed).** No fold rule, no ordering key, no status-event type, no vocabulary mapping,
no actor parameter, no validity predicate, no schema and no amendment. The brief names missing artefacts; naming is not
designing (ODC-31's discipline).

**NDR-04 (Observed — the near-name hazards, stated once).** Three items are easily confused and are not the same:

| Item | What it is | State |
|---|---|---|
| **DN-1** | v9 §8's textual question — whether "never a view, alias or projection" governs B's record as a view of A's record | **CONFIRMED 2026-11-19** (settled; narrow reading) |
| **ND-1** | the question this brief examines | **OPEN**, and **scope-ambiguous** (§2) |
| **DN-6a / DN-6b** | who may cause forgetting / who may retract content | **OPEN**, separate (A2-40; DS-06) |

**NDR-05 (Calculated — and the difference is substantive, not typographic).** DN-1 concerns whether one recipient's record
may be a projection of another's; DN-1's confirmation removed a divergence for PH1A-03 (TC-01; DS-01). ND-1 concerns the
**memory/status stream's semantics**. Conflating them would import a settled textual finding into an open semantic one.

---

## 2. What the records say ND-1 is

**NDR-06 (Observed — the originating characterisation: fold semantics).** A2-41 keeps two questions apart: **ordering** —
"which status event precedes which", which A2 answers *how* once a stream exists — and **fold semantics** — "what a later
status event does to state", marked "**ND-1, unspecified**".

**NDR-07 (Observed — and wider than fold semantics).** A2-42: "The status **vocabularies do not reconcile**", tabulating
`established`, `active`, `provisional`, `superseded`, `corrected`, `retracted`, `forgotten` against K-R12's event kinds,
§5's status field, and v1. Two values have no named producing event in that table — `active` and `provisional` — and
`established` (an event) and `active` (a status) are "the same fact under two names with no stated mapping" (A2-43(ii)).

**NDR-08 (Observed — the `provisional` half is superseded, and must not be repeated as found).** A2-43(i) claimed
`provisional` had no producer; that finding is **SUPERSEDED 2026-11-19** — formation at r33 §11 `:604` does produce it,
and what is missing is a **typed event representation** (ODBC-16/ODBC-17). The surviving half is (ii): the
`established`↔`active` mapping.

**NDR-09 (Observed — the ordering rule, recorded as a dependency rather than a question).** DS-09: "**ND-1 — the
memory/status stream has no ordering rule.** K-R12/K-R14 require it to be independently rebuildable, but r32 specifies no
order key and no fold rule for it. Two status events in one turn are indistinguishable." Recorded "only; no requirement
added, nothing decided".

**NDR-10 (Observed — the authority characterisation, in a different record).** CEV-28 treats ND-1 as silent "on who may
cause a lifecycle change", and CEV-34 goes further: "a validity predicate **is** an answer to 'who may cause a lifecycle
change.' So L3's falsifier is not merely open; it is **entangled with an owner decision**."

**NDR-11 (Observed — and the authority characterisation is what later records inherited).** QAD-12 lists "**ND-1** — who
may cause a lifecycle change" as OPEN; QBD-10, DQD-10 and DPD-07 repeat that gloss. Meanwhile the records that examine
causation assign **authority** to DN-6a/DN-6b: A2-05 ("nothing about who may cause an event (DN-6a/b)"), A2-40, DS-06,
SD-05 ("who may cause a statement-less event (DN-6a remains open)"), SD-23 ("it does **not** say who may cause them").

**NDR-12 (Calculated — so the corpus carries two characterisations, both dated 2026-11-19).**

| Characterisation | Where | What it names |
|---|---|---|
| **C-1 — semantics** | A2-41, A2-42, A2-43(ii); DS-09 | what a later status event does to state; how the vocabulary maps; whether the stream needs an order key |
| **C-2 — authority** | CEV-28, CEV-34 (and inherited by QAD-12, QBD-10, DQD-10, DPD-07) | who may cause a lifecycle change |
| **C-1 + C-2** | the owner's own framing of the target — "lifecycle status vocabulary, fold semantics, and transition authority" | all of it, as one target |

**Neither characterisation is wrong about its own evidence**, and neither is a misquotation. What is unresolved is
**which question ND-1 is** — and, if both, whether it is one question or a bundle.

---

## 3. The four candidate contents

**NDR-13 (Calculated — stated separately so none is chosen by drafting).**

| # | Content | The question it would ask | Evidence that exists | Evidence that is missing |
|---|---|---|---|---|
| **N-a** | **fold semantics** | what a later lifecycle event does to the current state — the transition function over K-R12's kinds | K-R12 (stream authority), K-R13 (read-time join), K-R14 (independent rebuild) | **no fold rule** for the memory/status stream (A2-41; DS-09) |
| **N-b** | **vocabulary** | how `established` (event) and `active` (status) relate; where `provisional` lives; whether every status value has a producer | K-R12's kind list; §5's status field; §11's formation `:604` (ODBC-16/17) | **no stated mapping**; no typed representation for formation (A2-43) |
| **N-c** | **ordering** | whether status events need an order key, and which | K-R15 fixes `(turn, seq)` for **acquisition** only; K-R14 requires independent rebuild | **no order key and no fold rule** for the status stream (DS-09) |
| **N-d** | **authority** | who may cause a lifecycle change | I17 (`:305`); I7; K-R12's stream split; the absence of any actor on `retract` in either implementation (DS-07) | **no authority rule**; DN-6a/DN-6b open (A2-40; DS-06) |

**NDR-14 (Calculated — N-a and N-b are not separable in practice, and are separable in principle).** A fold needs values
to fold; the vocabulary question is what those values mean. One could decide N-b alone (a mapping) and leave N-a open, or
decide N-a over a mapping adopted informally — but neither *entails* the other, and a decision that bundled them would
need to say which part it settled.

**NDR-15 (Calculated — N-c is conditional on N-a).** An ordering key is needed only if the fold is order-sensitive. DS-09
records the dependency and explicitly declines to decide whether it matters ("whether this matters depends on whether the
experiment can produce two status events in one turn — not established"). So N-c is a consequence of N-a, not a
peer of it.

**NDR-16 (Calculated — N-d is a different kind of question from N-a…N-c).** N-a…N-c ask what the stream **means**; N-d asks
**who may write to it**. They share the word "lifecycle" and nothing else: no fold rule supplies an authority condition
(GRP-31's axis-H finding, applied here), and no authority rule supplies a fold.

**NDR-17 (Observed — N-d is the one already occupied by another name).** DN-6a ("who may cause a character to forget") and
DN-6b ("who may retract content or another character's account") are open and explicitly not implied by DN-4 (DS-06;
OD-22). If ND-1's content is N-d, then ND-1, DN-6a and DN-6b are three questions over the same ground, and a decision
would have to demarcate them before answering either.

**NDR-18 (Unestablished — whether N-d even reaches all lifecycle kinds).** DN-6a/6b are framed around forgetting and
retraction. Whether "who may cause a lifecycle change" is meant to cover supersession, correction, formation and scope
changes as well is **not stated** — and the answer changes N-d's size, not just its name.

---

## 4. What is missing, stated once

**NDR-19 (Observed — the absent artefacts).**

| Missing | Which content needs it | Its absence is recorded by |
|---|---|---|
| a **fold rule** over lifecycle events | N-a | A2-41; DS-09; LCM-24 (per-kind admissibility absent) |
| a **status order key** (or a finding that none is needed) | N-c | DS-09 |
| an **`established`↔`active` mapping** | N-b | A2-43(ii); ODL-29 (preserved) |
| a **typed representation of formation** | N-b | ODBC-16/ODBC-17 |
| an **authority rule** for lifecycle transitions | N-d | A2-40; DS-06; SD-05; SD-23 |
| a **validity predicate** and an evaluation point | N-d's validity side | CEV-24; CEV-30; LCM-24 |

**NDR-20 (Observed — what canon does supply, so the gap is specific).** K-R12 makes the lifecycle its own authoritative
stream with a projection; K-R13 makes cross-stream use a read-time join; K-R14 requires independent rebuild and calls a
mismatch an integrity failure; K-R15 fixes ordering **for acquisition**; §11 states **promotion** conditions (`:611`–
`:613`) and names formation (`:604`). None of these is a fold rule for the lifecycle stream, and §11's scope is promotion.

**NDR-21 (Calculated — the gap is not an omission in one clause but an absent component).** Adding an order key, a fold
rule or an actor would each be a new element, not a clarification of existing text. That is why the readiness question is
"which content", not "which wording".

**NDR-22 (Observed — no implementation limitation is being converted).** Nothing here says the system must gain a fold
rule, an actor parameter or a status-event type. The absent artefacts are recorded as absences; whether any should exist
is the decision's business, not this brief's (ODC-27; QBD-26).

---

## 5. The dependency consequence — the part that changes a later step

**NDR-23 (Observed — the recorded edge).** DPD-10 records **ND-1 → DQ-4**: "Is a validity basis a separate obligation, and
is ND-1 its prerequisite?" depends on ND-1, and the promotion path's authority clause "presumes ND-1" (GRP-53).

**NDR-24 (Calculated — that edge holds only under C-2).** If ND-1 is N-d (authority), the edge is sound. If ND-1 is N-a,
N-b or N-c (semantics), the edge is **misattributed**: DQ-4 is a validity-side question whose gate would then be
**DN-6a/DN-6b** and the absent predicate, not the stream's semantics. **Narrowed:** DPD-10 stands as recorded for the
authority reading and is conditional on the scope determination; it is not withdrawn, because under C-2 it is correct.

**NDR-25 (Calculated — so the scope determination precedes the dependency's use, not the reverse).** A DQ-4 request
prepared on "ND-1 first" would be prepared on a scope that may not be DQ-4's gate. This is a reason to scope ND-1 before
using the edge, and not a recommendation about which question to decide first.

**NDR-26 (Observed — no other edge is affected).** The predicate → set edge (DPD-11) is untouched: it is analytic and
independent of ND-1. DQ-3 is settled. DQ-5 has no live decision (D5A-23). The prohibition is independent (DSC-20).

---

## 6. What a scope determination would have to fix

**NDR-27 (Calculated — unfilled, and prior to any option space).**

| # | Term | Why it must come first |
|---|---|---|
| **S-1** | which of **N-a…N-d** ND-1 names — one, or several | every later term depends on it (NDR-13) |
| **S-2** | if several, whether they are **one question or parts** to be decided separately | DQ-2's precedent: three parts, kept apart (DQD-02; DRD-02) |
| **S-3** | the **demarcation** from DN-1 (settled), DN-6a/DN-6b (open), and `active`/`established` (ODL-29, preserved) | three near-neighbours occupy adjacent ground (NDR-04, NDR-17) |
| **S-4** | that the determination **creates no fold rule, key, mapping, event, actor, predicate or schema** | a scope answer is a boundary, not an artefact |
| **S-5** | which questions it leaves open, stated rather than implied | the standing discipline (DQD-23; D3D-23) |

**NDR-28 (Observed — S-3 is the one most likely to be lost).** `active`/`established` is already a preserved item in the
ledger of open questions, and it is *inside* N-b. A scope determination that named N-b without saying whether it absorbs
or defers `active`/`established` would either duplicate a preserved item or silently reopen it.

**NDR-29 (Observed — nothing is recommended).** This brief does not say which content ND-1 should name, does not treat
any reading as more faithful, and does not propose that the scope determination be the next decision.

---

## 7. What this brief does not do

**NDR-30 (Observed).**

| Untouched | State |
|---|---|
| **ND-1** itself | **OPEN** and **not narrowed** by this brief |
| **DN-1** | **settled** — confirmation only, not reopened |
| **DN-6a / DN-6b** | **OPEN**; not answered or narrowed |
| **`active`/`established`** (ODL-29) | preserved; not reassessed |
| **OD-a**; **OD-b = L2**; **A2**; **D-3b** | settled; unchanged |
| **DQ-4**; the **validity predicate**; the admissible-transition set | **OPEN** / absent |
| **axis H**; **QAD-19**; the **report-only** boundary; the prohibition | unchanged |
| ODL-39; closure; determinism; §8.1; P-EVAL; C′; OD-g; Q-B | unchanged |

**NDR-31 (Observed — no schema is implied, including by example).** The brief names K-R12's kinds and §5's values because
they are the evidence of the mismatch. It proposes no event type, no field, no key format, no actor model and no
migration; the absences in NDR-19 are described, not scheduled.

---

## 8. Traceability

**NDR-32 (Observed).**

| What | Where |
|---|---|
| Ordering vs fold semantics; ND-1 unspecified | A2-41 |
| The vocabulary mismatch table | A2-42 |
| `established`↔`active`; `provisional`'s producer **superseded** | A2-43; ODBC-16/ODBC-17 |
| No order key, no fold rule; two status events indistinguishable | DS-09 |
| Authority assigned to DN-6a/6b | A2-05, A2-40; DS-06; SD-05; SD-23; OD-22 |
| ND-1 as "who may cause a lifecycle change" | CEV-28, CEV-34; QAD-12; QBD-10; DQD-10; DPD-07 |
| K-R12/K-R13/K-R14; K-R15's acquisition-only ordering | r33 `:109`–`:112`, `:117`–`:120`, `:125`–`:128` |
| §11 formation `:604`; promotion conditions `:611`–`:613` | r33 |
| I17 and I7 | r33 `:305`, `:293` |
| The ND-1 → DQ-4 edge | DPD-10; GRP-53 |
| Predicate → set, untouched | DPD-11; CEV-24; LCM-24 |
| DQ-5 has no live decision | D5A-23 |

---

## 9. Limits and record status

**NDR-33 (Observed — method).** **Newly read for this record:** A2 §10.3–§11 (A2-38…A2-45), DS-05…DS-09, CEV-22…CEV-36,
SD-05/SD-23, r33 `:95`–`:139`, A2-05/A2-40, and the records that gloss ND-1 (QAD-12; QBD-10; DQD-10; DPD-07).
**Not inspected:** any implementation's status handling, r32's bytes, the frozen preregistration, v9 beyond §7, the
`MEMORY-09-ORDERING-KEY-COMPARISON` record's tables in full. No test, probe or model call was run (NDR-34).

**NDR-34 (Observed — no experiment).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. The brief is textual and dependency-relative; it needs no execution.

**NDR-35 (Observed — the labels' limits).** NDR-02, NDR-05, NDR-12…NDR-18, NDR-21, NDR-24…NDR-27 are **Calculated**.
NDR-18 and the "whether it absorbs or defers" clause of NDR-28 are **Unestablished** on purpose. NDR-24 is **Narrowed**.
A2-43(i) is **Superseded**, and is labelled so rather than repeated.

**NDR-36 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Questions answered, scoped or opened as work | **0** — one scope question identified, no content chosen |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **research evidence** — no canonical capture, no ledger entry |

**End of record.** **ND-1 is scope-ambiguous in the records.** Two characterisations are in use, both dated 2026-11-19:
**C-1** — fold semantics, vocabulary and the status stream's ordering (A2-41…A2-43; DS-09) — and **C-2** — "who may cause
a lifecycle change" (CEV-28, CEV-34), inherited by QAD/QBD/DQD/DPD, while the causation records assign authority to
**DN-6a/DN-6b**. Four candidate contents follow (**N-a** fold, **N-b** vocabulary, **N-c** ordering, **N-d** authority);
six artefacts are missing, none of them a wording problem (NDR-19). **The ND-1 → DQ-4 edge holds only under C-2**
(NDR-24). A scope determination must name the content and demarcate it from **DN-1**, **DN-6a/DN-6b** and
**`active`/`established`** before any option space is drafted (NDR-27). **Nothing is answered, narrowed, designed or
recommended; no requirement, schema or canonical capture is created.**
