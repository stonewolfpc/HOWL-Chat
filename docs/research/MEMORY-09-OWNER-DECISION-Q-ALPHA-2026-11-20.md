# MEMORY-09 — Owner Decision Q-α: D-3b Reaches Both Streams

**Status:** owner-decision record. **Q-α settled: YES — D-3b applies to acquisition and lifecycle events.** Documentation
and read-only analysis only. **No schema change, event or relation design, implementation, test, probe, benchmark,
model call, canonical edit, registry change or ledger entry is authorised.** Q-β remains OPEN. ND-1 remains OPEN.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Decision | **Q-α = YES.** D-3b's class spans both streams; a statement-less lifecycle event carries the same three obligations as a statement-less acquisition event |
| Decision scope | **reach of an existing obligation.** The obligation's content is unchanged; what was open was whether it binds lifecycle events |
| Authority | Owner, stated in session 2026-11-20 and recorded here; as with every prior confirmation, **no persisted signature exists** |
| Characterised by the owner as | **not** a new requirement — the existing decision's reach, determined |
| Was open because | no provision stated inclusion or exclusion (`...Q-ALPHA-D3B-SCOPE...` QAX-24) |
| Canonical edits | **None.** r33, v9, ledger, conformance map and D-3b's text are unchanged |
| Still open | **Q-β** (ODC-18) · **ND-1** · the validity predicate · the admissible-transition set · Q-γ (**QAD-19**) · P-G0/P-G1/P-G2 unselected · ODL-39 |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The decision as stated

**QAD-01 (Source-derived — owner).**

> "Q-α settled — YES
>
> D-3b applies to both acquisition and lifecycle events.
>
> That means a statement-less lifecycle event must:
>
> * Retain its actual cause.
> * Have an independently identifiable occurrence.
> * Never fabricate evidence or infer a statement position.
>
> The important boundaries remain:
>
> * Recording a cause does not authorize the event.
> * Recording a cause does not prove the event was justified.
> * ND-1 remains open.
> * No schema, implementation, or new requirement is authorized.
>
> …
>
> P-G1's cause-retention floor is now established across both streams. It is no longer merely a proposed extension."

**QAD-02 (Observed — provenance).** The owner's statement as supplied on 2026-11-20, transcribed without alteration.
No persisted signature exists and no canonical file was edited to record it (QAD-27).

**QAD-03 (Observed — what was actually open).** Q-α was not a question about what D-3b requires; it was a question
about **which events it binds**. The examination found the decision's **wording** stream-neutral ("events without a
source statement") while its **decision scope** was acquisition-side, with the corpus stating neither inclusion nor
exclusion (QAX-07, QAX-17, QAX-21, QAX-24). The decision supplies that missing determination.

---

## 2. What is decided

**QAD-04 (Source-derived).** **D-3b's class spans both streams.** Acquisition events and lifecycle events are both
within its reach.

**QAD-05 (Source-derived).** **A statement-less lifecycle event must retain its actual cause** — M-1 (SD-02; GRP-23)
applies to it as written.

**QAD-06 (Source-derived).** **A statement-less lifecycle event must have an independently identifiable occurrence** —
M-2 (SD-03) applies to it as written.

**QAD-07 (Source-derived).** **Neither fabrication nor an inferred statement position is permitted** for such an
event — M-3 (SD-04) applies to it as written.

**QAD-08 (Calculated — the obligation's content is unchanged).** The three parts are D-3b's existing M-1, M-2 and
M-3, applied to a class the decision now includes. No clause was added, no threshold changed, no new duty created.
This is consistent with the owner's characterisation: a reach determined, not a requirement added (QAD-09).

**QAD-09 (Observed — the owner's own framing of its character).** The decision states that no **new requirement** is
authorised, and describes P-G1's floor as "established across both streams… no longer merely a proposed extension".
The determination is therefore recorded as **scope**, not as obligation (QAD-08).

**QAD-10 (Calculated — but the determination is not derivable from the corpus).** QAX-24 found that no text stated
inclusion or exclusion. So the reach is fixed by this decision, not read off an existing provision. Both statements
hold without conflict: the **content** was already decided (D-3b), and the **reach** was not, and is now. Recording
both so neither is mistaken for the other.

**QAD-11 (Calculated — which earlier reading this resolves).** QAX-19's R-prop (the class is property-based, so
occurrence is sufficient) is the reading the decision adopts; QAX-20's R-stream (the class is the events the decision
was addressing) is not. The examination declined to choose between them because choosing required an owner
determination (QAX-25); that determination has now been made. Note that the decision is **not** an inference from
statementlessness (QAX-02's exclusion): it is an owner statement about reach, which is the only thing that could have
settled it.

---

## 3. Boundaries that attach

**QAD-12 (Source-derived — recording is not authorising).** Recording a cause does not authorise the event. This
restates D-3b's own separation (SD-05, SD-23) and I17's asymmetry: authority comes from the conditions of
establishment, not from the record (r33 `:305`). **ND-1** — who may cause a lifecycle change — remains **OPEN**, and
this decision says nothing about it.

**QAD-13 (Source-derived — recording is not justification).** Recording a cause does not prove the event was
justified. The cause identifies a producer; admissibility is a separate question with no predicate (CEV-24) and no
admissible-transition set (LCM-24). Axis H is unchanged (GRP-31).

**QAD-14 (Source-derived — no schema, implementation or new requirement).** Nothing in this decision authorises a
field, a carrier, an event type, a marker, a mechanism or an implementation, and nothing adds a requirement. The
cause's **shape** remains unsettled (SD-18) and is not touched.

**QAD-15 (Observed — Q-β is untouched).** Whether an **actor-reported rationale** also qualifies as a **ground**
remains **OPEN** (ODC-18). This decision is about reach, not about classification, and Q-β is neither narrowed nor
answered by it. The two must not be collapsed (PGR-25/PGR-26).

**QAD-16 (Observed — the statement-bearing region is untouched).** GRP-26's region — transitions that carry a source
statement — is unaffected: D-3b does not reach it, and no provision requires a basis for it. Its size remains
unestablished.

---

## 4. What this resolves, and what it does not

**QAD-17 (Calculated — resolved).**

| Was open | Now |
|---|---|
| Q-α — does D-3b reach lifecycle transitions? | **YES** (QAD-04) |
| P-G1's status — proposal or description? | **description** (QAX-29's conditional resolves; QAD-20) |
| whether P-G0 could decline the cause on statement-less lifecycle events | **no** — that would now be a violation of an existing obligation (QAD-18) |

**QAD-18 (Calculated — the first consequence, stated narrowly).** Because M-1 now binds statement-less lifecycle
events, a position that retained **no cause** for them would not be a neutral silence; it would be a non-compliance.
This is A-2 applied to a class whose reach is no longer in doubt. What that does to P-G0's **stated content** is a
separate question, taken up in the reassessment record — not decided here.

**QAD-19 (Unestablished — a new question this decision exposes, recorded not raised as work).** D-3b's class is
"events without a source statement". With both streams in reach, **which lifecycle transitions are statement-less**
becomes operative rather than academic: entry, supersession, correction, retraction and forgetting are candidates
(TGP-10/TGP-11; CLV-41), and the boundary is **not drawn** anywhere. This record does not draw it, does not name a
criterion, and does not treat it as authorised work. It is recorded so that "D-3b applies" is not read as "every
lifecycle transition carries a cause": the obligation is per-event, not per-stream (TGP-11).

**QAD-20 (Observed — P-G1's status, in the owner's terms).** "P-G1's cause-retention floor is now established across
both streams. It is no longer merely a proposed extension." Recorded as stated: P-G1 describes an obligation that
binds, and does not itself request authority.

**QAD-21 (Observed — nothing else moves).** OD-b = L2 (ODL-01); OD-a; D-3b's text; A2; DN-4; Q-B; entry ownership;
the §5/§11 provenance findings; closure (ODL-26); determinism (ODL-28); `active`/`established`; §8.1; P-EVAL; C′;
OD-g; DN-6a/DN-6b; **ODL-39**; the 118 requirements. All unchanged. No position is selected, ranked or dropped.

---

## 5. Traceability

**QAD-22 (Observed).**

| What | Where |
|---|---|
| The decision | this record, §1 |
| The reach question as it stood, and why it was open | `...Q-ALPHA-D3B-SCOPE...` QAX-07, QAX-17, QAX-21, QAX-24, QAX-25 |
| The two readings, R-prop and R-stream | QAX-19, QAX-20, QAX-23 |
| D-3b's three parts and their kinds | SD-01…SD-04; GRP-23, GRP-24 |
| Recording a cause ≠ authorising | SD-05, SD-23; I17 `:305`; ODC-25 |
| Q-β as the residual classification | ODC-18, ODC-19; CGE-43 |
| The statement-bearing region | GRP-26 |
| Axis H unchanged | GRP-31 |
| P-G1's description-or-proposal question | QAX-29; PGR-25 |
| Per-event, not per-stream | TGP-11 |

---

## 6. Limits and record status

**QAD-23 (Observed — method).** This record transcribes a determination and restates boundaries already documented.
No canon was re-read beyond the items cited; no implementation was inspected; no test, probe or experiment was run.

**QAD-24 (Observed — no experiment available).** 0/144 sessions, no corpus, no `plan_sha256`, G07/G12 unbound,
`ExecutionNotAuthorised` at `conversation.py:177`. The decision is a scope determination; it needs no execution.

**QAD-25 (Observed — the labels' limits).** QAD-08, QAD-10, QAD-11, QAD-17 and QAD-18 are **Calculated** — read from
the decision against QAX/SD/GRP, not from new text. QAD-19 is **Unestablished** on purpose. Nothing here is offered
as a textual finding in r33 or v9.

**QAD-26 (Observed — no conflict found).** The decision is consistent with D-3b's text, with the examination's
findings, and with the boundaries it restates. The one seam is QAD-10: **content was already decided; reach was
not** — and the decision is what fixes it.

**QAD-27 (Observed — record status).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Obligations created or widened | **0** — an existing obligation's reach determined |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; no canonical capture, no ledger entry |

**End of record.** **Q-α is SETTLED — YES.** D-3b's M-1, M-2 and M-3 bind statement-less lifecycle events as they
bind statement-less acquisition events; the obligation's content is unchanged and its reach is now determined
(QAD-04…QAD-08). **Recording authorises nothing and justifies nothing**; **ND-1**, the validity predicate, the
statement-less boundary for lifecycle events (QAD-19) and **Q-β** remain open; **no schema, implementation or new
requirement is authorised.**
