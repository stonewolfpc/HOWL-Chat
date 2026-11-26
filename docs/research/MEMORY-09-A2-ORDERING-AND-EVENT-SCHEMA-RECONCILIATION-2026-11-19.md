# MEMORY-09 — A2 Ordering Decision and Event-Schema Reconciliation

**Status:** owner-decision record + read-only reconciliation. **A2 selected by the owner.** No schema change,
no implementation, no test, no probe, no model call, no corpus, no Git operation, **no ledger entry**.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Scope | A2's exact scope · memory-event schema search · envelope vs separate · minimum event contract · A2 feasibility per stream · cross-stream dependencies · `evidence_span` · cause/evidence gaps · identity · D-3a/SD-02/DN-6a-b/ND-1 |
| Canon read | r33 `dfb28b6003ca356f` · v9 `698254ca262c0874` · ledger `e5008020ff68bd21` · map `b4ba33ccbfe8f8c1` |
| Canon edited | **none** |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. The A2 decision and its exact scope

**A2-01 (Source-derived — owner).** The owner selects **Family A2: independent, universal occurrence
ordering within each authoritative event stream** — acquisition/knowledge, and memory lifecycle/status.

**A2-02 (Source-derived).** A single global sequence shared between streams is **not required**, and
cross-stream allocation or a shared high-water mark must not be introduced merely to impose global order.

**A2-03 (Source-derived).** A2 does **not** establish that cross-stream relationships are never relevant.
Explicit causal relationships must be investigated where an operation needs one, **without treating them as
a substitute for chronological ordering**.

**A2-04 (Source-derived).** A1 is not to be reopened unless authoritative evidence shows A2 insufficient.

**A2-05 (Inferred — what A2 settles).** It settles **one** question per stream: *which event precedes which
other event of the same stream*. It says nothing about what a later event does to state (ND-1), nothing
about who may cause an event (DN-6a/b), and nothing about cross-stream comparison.

---

## 2. Search for an existing memory-event schema

**A2-06 (Observed — where memory events are asserted).** Exactly three authoritative passages refer to them,
and all three are the *same sentence*:

| Artifact | Line | Content |
|---|---|---|
| r33 `ARCH-STAGE-CONTRACT` | 321 | "Memory lifecycle history \| append-only **memory events** \| authoritative for status/scope/class" |
| `CONTRACT-NORMATIVE-2026-11-15` | 45 | identical row |
| r33 `ARCH-STAGE-CONTRACT` | 111, 367–368 | K-R12; two-streams rule; "what is currently true only from memory events" |

**A2-07 (Observed — the schema inventory).** r33 §6 contains **§6.1 MemoryRecord, §6.2 AcquisitionEvent,
§6.3 KnowledgeEdge, §6.4 consistency rules, §6.5 derived views**. There is no §6.x for memory events. The
normative extract has **no schema sections at all** (decisions, invariants, rules, open questions,
baseline).

**A2-08 (Observed — the historical definition, and why it does not settle this).** The v1 design of
**2026-11-11** (`DESIGN-MEMORY-CONTINUITY`, §4.2) does define an event structure:

```
Event: event_id, continuity_id, chat_id, turn_id, source_message_id, event_type,
       participants[], summary, evidence_span, observed_at, effective_at,
       origin: user | model_proposal | import | system, confidence,
       status: active | corrected | superseded | retracted | forgotten,
       supersedes_event_id?
```

**A2-09 (Inferred).** So the answer to the work order's three options is **a qualified (2) and a firm (3)**:

- the *envelope* was **defined historically** (v1, above), and omitted from r32/r33 — option (2);
- but the **lifecycle stream** was never specified as a stream in any version — v1's `Event` is a
  **narrative-event** structure carrying `participants`, `summary`, `effective_at` and a `status` field;
  lifecycle operations appear there as *operations on events* ("Correction … a new event identifies the
  earlier event"), not as a separate authoritative log — option (3).

**A2-10 (Observed — the decisive reconciliation).** The lifecycle **data** did not vanish with the schema.
r33 §6.1's `MemoryRecord` carries:

> `attribution{kind, by}`, **`status`**, `established_turn`, **`supersedes`/`superseded_by`**,
> **`correction_of`**, `conflict_group`, …

**A2-11 (Inferred — a real contradiction, previously unrecorded).** K-R12 says status, class and scope are
"a **projection** of that stream, **not mutable authoritative fields**". §6.1 nevertheless lists `status` as
a record field, and `supersedes`/`superseded_by`/`correction_of` as record **relations** for which no event
exists to project from. §5 compounds it by calling status a "**field**". The charitable reading is that §6.1
stores a *cached projection* (§6.4 says the projection check "compares to stored projection"), which is
admissible — but the relations have no source at all.

**A2-12 (Inferred — the specification gap, stated precisely).** **The lifecycle stream is authoritative by
declaration and unspecified by schema, while its data survives as mutable record fields and relations.**
That is a stronger statement than "the schema is missing": the *authority claim* (K-R12) is unfulfilled by
*any* structure in r33, and the record currently holds what K-R12 forbids it to own.

---

## 3. Shared envelope, or separate structures?

**A2-13 (Observed — evidence for each).**

| Evidence | Points to |
|---|---|
| v1 defined **one** `Event` with a discriminating `event_type` and a `status` field covering active/corrected/superseded/retracted/forgotten | **shared envelope** |
| r32/r33 define `AcquisitionEvent` **as its own schema** with acquisition-only fields (`via` closed set, `confidence_in_source`, `from_character_id`) | **separate structures** |
| K-R14: "two streams, two checks … neither projection may consult the other" | **separate** |
| K-R15 speaks of "the ordered log" for a pair; K-R12 speaks of a separate stream for memory | **separate** |
| v1's `supersedes_event_id` presupposes one identifier space across event kinds | **shared envelope** |

**A2-14 (Inferred).** r32 already **implicitly chose separate structures** — by specifying one stream's
schema and not the other. This was not recorded as a decision; it is an artifact of which stream got
written down.

**A2-15 (Inferred — consequences, no selection).**

| | Shared envelope (v1's shape) | Separate structures (r32's shape) |
|---|---|---|
| A2 fit | needs a **stream discriminator** inside one ordered space, or two sequences inside one structure — the ambiguity Family B suffered | natural: each stream owns its own sequence |
| Identity | one identifier space, so cross-kind id collisions must be excluded by construction | two spaces; cross-stream causal refs must name the stream |
| Causal refs | `supersedes_event_id` works directly (v1) | needs a typed reference (`memory_event:<id>`) |
| Precedent | defined, then dropped without a recorded reason | current, but half-specified |
| Risk | re-creates a second ordering authority inside one structure | leaves the status stream free to be specified independently |

**A2-16 (Observed).** No selection is made here. The work order forbids selecting without presenting
evidence and consequences, and both are above.

---

## 4. Minimum event-contract properties

Properties only — **no schema is proposed as a contract**.

**A2-17 (Inferred — property matrix).**

| Property | Acquisition stream | Status stream |
|---|---|---|
| **A. Identity** | `id` exists (`AcquisitionEvent.id`); uniqueness stated (K-R16); **stability unstated** and violated in the harness | **no carrier** — no event type, so no id |
| **B. Ordering** | needs a persisted position; `turn` exists, position within turn does not (r33 §6.2 has no `seq`) | **no carrier**; no schema to attach a position to |
| **C. Cause** | **no carrier** — `via` is acquisition method, not cause | **no carrier** — `status` is a value, not an event |
| **D. Evidence** | `evidence_span` exists in §6.2; **nullable/absent semantics undefined** (D-1) | **no carrier**; per D-3b a lifecycle event may have none |
| **E. Occurrence** | implied by the appended record; **not independently identifiable** (positional ids) | **no carrier**; §6.1 keeps no `status_changed_turn` |

**A2-18 (Inferred).** For the status stream, **four of five properties have no carrier at all** — not a
missing definition, a missing structure. That is the practical meaning of A2-11.

**A2-19 (Observed — the record cannot answer "when").** §6.1 has `established_turn` but **no
last-status-change turn**. So a retraction's occurrence is not recorded on the record, and no event exists
to supply it. As-of questions ("was this retracted at turn 40?") are unanswerable from either structure.

---

## 5. A2 feasibility per stream

**A2-20 (Inferred — acquisition: feasible now).** What is ordered: the pair's learn/revoke log. Assignment:
requires a persisted position (A2's own requirement, still unimplemented). Identity: needs stored ids.
Insertion-order independence: satisfiable once the fold sorts by the persisted position. Multiple events in
one turn: A2 orders them by distinct positions; the extractor's one-per-span limit is a non-issue under A2.
No turn or span: A2 does not require either.

**A2-21 (Inferred — status: not yet assessable).** A2's mechanism is expressible, but **there is nothing to
express it on**. Feasibility is blocked on the schema, not on A2.

**A2-22 (Observed — what A2 does not need, confirmed).** No wall-clock field is required for ordering:
`observed_at` exists only in v1 and is exactly the carrier K-R15 forbids relying on.

---

## 6. Cross-stream dependencies

**A2-23 (Observed — three authoritative rows whose *scenario setup* is temporal across streams).**
Corrected: the temporal word is in the setup, not in the assertion (XT-11).

| Row / rule | Operation (setup) | Assertion |
|---|---|---|
| **H10** | "Source retracted after disclosure" | the recipient's account **persists** — non-cascade |
| **K16** | "Basis retracted after promotion" | promotion **stands** — non-cascade |
| **A2a** | memory later retracted | the character still holds what it acquired — non-cascade |

**A2-24 (Superseded — see `MEMORY-09-CROSS-STREAM-TEMPORAL-CLARIFICATION-2026-11-19.md`).** The original
claim — that H10/K16 prove A2 sufficient for cross-stream temporal questions — **does not hold**, for two
reasons: (i) those rows require no computed comparison, only non-cascade behaviour (XT-11/XT-12); (ii) `turn`
is the comparandum, not A2, and it is unavailable for the turn-less events D-3b permits (XT-04–XT-08).
A2's cross-stream sufficiency is **Unestablished**. A2 itself is unaffected: it answers **within-stream**
ordering only, and its ordinals are deliberately not cross-stream comparable (XT-03–XT-05).

**A2-25 (Inferred — causal is not chronological).** Where a dependency is *causal* — v1's
`supersedes_event_id`, `correction_of`, §6.2's `source_memory_id` — the relation names **which** event
produced which, and says nothing about order. A causal link can point backwards in the log; it therefore
cannot substitute for a position, and a position cannot substitute for a cause. The owner's instruction is
structurally correct, not merely cautious.

**A2-26 (Unestablished).** Whether any operation needs cross-stream order **for the same memory within one
turn** (e.g. retracted and re-established in one turn). No authoritative operation states it; H10/K16 are
turn-separated by construction.

---

## 7. `evidence_span` — findings reconciled (D-1 remains open)

**A2-27 (Observed).** `evidence_span` appears in **four schemas** (r33 §§6.1, 6.2; v1 `Event`, v1
`StateProposal.changes[]`) and is **defined in none**. r33 mentions it once, in §6.2's field list.

**A2-28 (Unestablished — the three candidate meanings).** Identifier, text range, or both. Nothing states:
its identity (is it a stable id or an offset pair?), its position (offset into what?), its stability (do
offsets survive edit?), its resolvability (what resolves it to text?), or its **comparison semantics** (can
two spans be ordered?).

**A2-29 (Inferred — comparison semantics are load-bearing only under B, and B is superseded).** Under A2, a
statement's position in its turn is **provenance**, not the ordering authority. So A2 removes the pressure
to define span comparison — a genuine simplification produced by the decision.

**A2-30 (Source-derived).** D-3b: a span must never be fabricated and an inferred statement position must
never be substituted. D-1's remaining work is therefore definitional: nullability, the "absent" state, and
the identifier-versus-range question.

---

## 8. Cause and evidence representation gaps

**A2-31 (Inferred — cause).** CO-1 has no carrier in either stream. `via` cannot serve: two forgetting events
may both be `via=system` while one was a lifecycle operation and the other a different system process —
identical `via`, different cause. Cause must also be distinguishable from **who is authorised** (DN-6a/b
unchanged) and from **evidence** (no statement is implied by a cause).

**A2-32 (Inferred — evidence needs three states, not two).** The work order requires distinguishing "no
evidence exists" from "evidence exists but was not recorded". Both are distinct from "evidence exists and is
recorded". An empty string cannot carry the difference, and a fabricated span is forbidden (D-3b). A
two-state nullable field is **insufficient**; the contract must name the three cases.

**A2-33 (Inferred).** The harness's `revoke()` (D-3a) is precisely the third→first substitution: evidence
belonging to a different event reused to stand in for absence.

---

## 9. Identity and reconstruction

**A2-34 (Observed — two distinct defects, same class).**

| Location | Construction | Fails |
|---|---|---|
| `eval/mem1/model.py:276`, `:310` | `f"ev-{len(self._events) + 1:05d}"` | stability across reconstruction; independent identity |
| `internal/backend/types/memory.go` `NewMemoryID()` | `"mem_" + time.Now().Format("20060102150405")` | **uniqueness** (second granularity ⇒ collisions within one second) and stability |

**A2-35 (Inferred).** The second is more severe than the ordering question it sits beside: two memories
formed in the same second share an id, which breaks K-R16 outright and makes any per-memory history
ambiguous. Neither is authorised to fix.

**A2-36 (Source-derived — owner).** Identity and ordering stay separate. A stable id alone does not
establish order; a position alone does not identify an event.

---

## 10. Status of the carried items

### 10.1 D-3a — forgetting provenance

**A2-37 (Observed).** `revoke()` stamps its event with the acquisition's `source_span` (`model.py:303–312`).
Required correction properties, not implementation: (i) the event carries its own cause; (ii) absence of
conversational evidence is recorded as absence; (iii) ordering does not depend on the borrowed span; (iv)
the borrowed span is never presented as evidence for the forgetting. **Not implemented, not scheduled.**

### 10.2 SD-02 obligations against the 118

**A2-38 (Calculated — re-verified, 118 rows).**

| Clause | Nearest rows | Coverage |
|---|---|---|
| CO-1 cause recorded | K-R90 (learning only), F2 (prohibits truth on knowledge, not cause) | **uncovered** |
| CO-2 absence explicit | I15 (absence is not evidence about absence — different claim) | **uncovered** |
| CO-3 occurrence identifiable | K-R16 (uniqueness), F7 (replay determinism — no stable ids required) | **partial** |

**A2-39 (Inferred).** Documented as **in-scope coverage gaps**. No row added, altered or closed; the set
remains **118**; the map is untouched.

### 10.3 DN-6a / DN-6b — OPEN

**A2-40 (Source-derived).** "Who may cause forgetting" and "who may retract content" remain open. **A2
authorises nobody.** Neither origin, cause, nor the ability to record an event implies authority.

### 10.4 ND-1 — separate, and now larger

**A2-41 (Inferred — two distinct questions, kept apart).**

| Question | Status |
|---|---|
| **Ordering**: which status event precedes which | A2 answers *how* (once a stream exists) |
| **Fold semantics**: what a later status event does to state | **ND-1, unspecified** |

**A2-42 (Observed — ND-1 is wider than fold semantics).** The status **vocabularies do not reconcile**:

| Value | r33 K-R12 (event kinds) | r33 §5 (status field) | v1 (2026-11-11) | Produced by |
|---|---|---|---|---|
| `established` | ✔ | — | — | an event |
| `active` | — | ✔ | ✔ | **no named event** |
| `provisional` | — | ✔ | — | **no event, no operation** |
| `superseded` | ✔ | ✔ | ✔ | an event |
| `corrected` | ✔ | ✔ | ✔ | an event |
| `retracted` | ✔ | ✔ | ✔ | an event |
| `forgotten` | — | — | ✔ | dropped |

**A2-43 (Inferred, two findings — part (i) SUPERSEDED 2026-11-19: `provisional` does have a producer, formation at r33 §11 `:604`; what is missing is a typed event representation. See `MEMORY-09-OD-B-PROPOSED-ANSWER-CHECK-2026-11-19.md` ODBC-16/17).** (i) **`provisional` has no producer anywhere** — a status value nothing
can create; the implementation's candidate review uses `pending | approved | rejected`, so `provisional` may
be a *formation* state, not a truth state — **Inferred, not established**, and if true it belongs to a third
workflow ND-1 does not yet name. (ii) `established` (event) and `active` (status) are the same fact under two
names with no stated mapping.

**A2-44 (Observed — forgetting's home, resolved).** Knowledge-forgetting is `via=revoked` in the
**acquisition** stream (K-R19/R1c: "a forgotten pair has no edge"). Memory truth status is the **memory**
stream. These are two different operations with two different homes, and only the first has a schema. DN-4's
forgetting is the first; DN-6a asks who may cause it.

---

## 11. Remaining owner decisions

| # | Decision | Example | Consequence of each way |
|---|---|---|---|
| **OD-a** — **SETTLED 2026-11-19: Option B, separate structures** (`MEMORY-09-OWNER-DECISION-OD-A-2026-11-19.md`) | **Envelope or separate** (§3) | A retraction and a retraction of a *promotion* | shared: one id space, causal refs simple, risk of two sequences inside one structure; separate: clean per-stream A2, typed cross-refs required |
| **OD-b** | **Status-stream scope** | Is `provisional` a truth state or a review state? | if review: ND-1 gains a third workflow; if truth: it needs a producing operation, which none exists |
| **OD-c** | **Ordering allocation scope** (A2 follow-up) | two memories in one chat, two chats | per chat (matches `owner_chat_id`), per stream, or per memory |
| **OD-d** | **Cause descriptor shape** (CO-1) | model-initiated forgetting vs lifecycle sweep, both `via=system` | a field, a relation, or an envelope element — shape undecided, authority undecided |
| **OD-e** | **Evidence three-state contract** (CO-2, D-1) | evidence not recorded vs none exists | needs an explicit representation; nullable alone cannot express it |
| **OD-f** | **DN-6a / DN-6b** | may a model forget unprompted? | untouched by A2; blocks conformance, not scripted experiments |

**A2-45 (Inferred).** **OD-a and OD-b are the gating pair**: the envelope choice determines whether there is
one identifier space, and the status scope determines what the stream must contain. OD-c/d/e are shaped by
those answers.

---

## 12. Not inspected / not established

**A2-46 (Observed — not inspected).** `frontend/dist/js/memory-candidates.js` (built artifact; only its
comment matched); the frozen preregistration and probe taxonomy beyond what earlier records cite; the 711
`docs/**/*.md` in full — searches were term-based, so a memory-event schema under an unrelated name would
have been missed. `internal/backend/memory/*.go` was **not** read beyond the types file.

**A2-47 (Unestablished).** Whether a memory-event schema exists in a frozen or non-`docs` artifact; whether
`provisional` is a review state; whether two events for one memory in one turn is required; whether the
`established`/`active` mapping was ever stated.

**A2-48 (Observed — search errors).** Two searches reported `rg: .\models: The system cannot find the path
specified (os error 3)`. The path is absent from this checkout; **results from those searches exclude
`.models/` only**, and no other directory.

---

## 13. Verification

| Check | Result |
|---|---|
| Canonical files edited | **0** · r33 `dfb28b60…` · v9 `698254ca…` · ledger `e5008020…` · map `b4ba33cc…` |
| Code / schema edited | **0** (`model.py` mtime 2026-11-16; `engine.py` 2026-11-17) |
| Ledger entry written | **0** — no O-16 |
| Conformance rows | **118**, unchanged |
| Tests / probes / model calls | **0** |

**End of reconciliation.** A2 is recorded and scoped. The status stream is unspecified, its data lives in
record fields K-R12 forbids, and `provisional` has no producer.
