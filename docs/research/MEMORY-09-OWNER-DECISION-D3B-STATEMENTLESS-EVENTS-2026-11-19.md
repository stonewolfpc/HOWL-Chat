# MEMORY-09 — Owner Decision D-3b: Events Without a Source Statement

**Status:** owner-decision record. **D-3b answered.** Documentation and read-only analysis only. No schema
change, no implementation, no test, no probe, no model call, no corpus, no Git operation.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Decision | **D-3b** — MEMORY-09 may create events without a source statement |
| Authority | Owner, stated in session and recorded here; **no persisted signature exists**, as with every prior confirmation |
| Canonical edits | **None.** v9, r33, ledger, conformance map unchanged; proposed ledger entry in §9, **not made** |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. The decision, verbatim

> MEMORY-09 may create events without a source statement. Such events must retain their actual cause and
> independently identifiable occurrence. Absence of conversational evidence must never be replaced with
> fabricated evidence or an inferred statement position.

### 1.1 What it settles

**SD-01 (Source-derived).** **Events may exist with no source statement.**

**SD-02 (Source-derived).** Such an event must retain **its actual cause**.

**SD-03 (Source-derived).** Such an event must have an **independently identifiable occurrence** — it is a
thing that happened, identifiable as itself, not as a position in something else.

**SD-04 (Source-derived).** Absence of conversational evidence must **never** be replaced by fabricated
evidence **or by an inferred statement position**.

### 1.2 What it does not settle

**SD-05 (Observed).** It does not decide:
who may cause a statement-less event (**DN-6a remains open**); how the occurrence is identified in the
schema; what the ordering key is; whether a cause is a field, a descriptor or a relation; any
implementation. Recording a cause is **not** an authority decision, and this record does not slide from
one into the other.

---

## 2. Why this decision is consistent with the existing design

**SD-06 (Source-derived).** It is not a novelty. Two existing provisions already contemplate knowledge
that no statement produced:

- **`via` includes `narrator_grant`** (§6.2's closed set) — a grant is an act, not a character's utterance;
- **S32** tests "ambient availability: a newcomer receives 'the city has a north gate' **without
  acquisition**" — knowledge present with no acquisition event at all.

**SD-07 (Inferred).** So the decision closes a gap between what T2 already permits (`narrator_grant`,
ambient knowledge) and what the event schema assumed (that every event points at a span). It makes the
schema's assumption explicit rather than introducing a new capability.

---

## 3. What it eliminates from the PH1A-05 option space

This is the decision's largest effect, and it is a narrowing, not a widening.

### 3.1 Span position cannot be the universal ordering key

**SD-08 (Inferred).** If an event may have no source statement, then `(turn, span position)` is
**undefined** for those events. A key that some events cannot supply is not a total order. **O-G in its
current form is eliminated as a universal remedy** — it survives only as a partial scheme for
statement-bearing events.

**SD-09 (Superseded).** The assessment's conclusion that O-G "may need no new field at all" is superseded.
It may still need no new field *for statement-bearing events*; the class without statements forces the
question regardless.

### 3.2 "Sort it at the end of the turn" is forbidden

**SD-10 (Source-derived).** The obvious repair — give a statement-less event a position at the end of its
turn so the ordering stays total — is **exactly what SD-04 prohibits**: an "inferred statement position".

**SD-11 (Inferred).** That prohibition is well aimed. Such a synthetic position would be
indistinguishable in the log from a real one: the fold would place a model-initiated loss after a
statement it never followed, and the history would be false in the same way the `revoke()` defect makes it
false (§5). The decision forbids papering over the absence with a plausible-looking value.

### 3.3 Therefore the ordering key must be carried by every event

**SD-12 (Inferred).** A total order over the pair's log requires a key that **every** event possesses and
that was **assigned to that event at write**. That is a materially different requirement from "find the
existing carrier" — and it is the first time this investigation has produced a conclusion that points
*toward* a new field rather than away from one.

**SD-13 (Unestablished).** Whether that key is a counter, an ordered occurrence identity, or a
`(turn, position)` pair whose position is assigned for every event. The options are enumerated in §4;
**none is selected here**, and the choice is the owner's.

---

## 4. The constraint space for the ordering key

**SD-14 (Inferred).** Two families remain, and they differ in what they can express:

| Family | Shape | Consequence |
|---|---|---|
| **One universal key** | every event carries an occurrence order (counter or ordered identity), statement or not | one source of order; statement position becomes *evidence*, not ordering; a single comparison rule |
| **Composite with assigned sub-position** | `(turn, sub)` where `sub` is statement-derived when available and self-assigned otherwise | keeps turn-first semantics, but one field now has **two provenances** — the "competing sources of truth" shape r32 rejects elsewhere |

**SD-15 (Inferred).** The composite family's flaw is not aesthetic: if `sub` is sometimes derived from a
statement and sometimes assigned, then two events can be *incomparable in kind* while sharing an encoding —
precisely the ambiguity SD-04 exists to prevent. The universal family avoids that by never pretending the
order came from the text.

**SD-16 (Observed).** No option is chosen. This section records the consequence of the decision, not a
design.

---

## 5. Evidence and cause are now two different things

**SD-17 (Inferred — the conceptual shift).** Until now `evidence_span` was doing double duty: what
supports the record, and implicitly where the event sits. Once an event may have **no** evidence, the two
separate:

| | What it answers | Present when there is no statement? |
|---|---|---|
| **Evidence** (`evidence_span`) | "what in the conversation supports this?" | **no** — and must not be invented (SD-04) |
| **Cause** | "what produced this event?" | **yes, always** (SD-02) |

**SD-18 (Inferred).** r32 has no field for the second. `via` names the *kind* of acquisition; it does not
name the *actor or mechanism* that caused the event, and it cannot, because the same `via` can arise from
different causes. **A cause descriptor is therefore implied by SD-02** — its shape is not settled, and
whether it is a field, a relation or part of an event envelope is an implementation-shape question this
record does not answer.

**SD-19 (Source-derived).** `evidence_span` must consequently be **nullable**, and its absence must be
**meaningful** rather than a default: "no conversational evidence exists" is a fact about the event, and
must be distinguishable from "evidence was not recorded". Those are different states, and neither may be
encoded as an empty string.

---

## 6. Effects on the open items

| Item | Effect |
|---|---|
| **D-1** (`evidence_span` definition) | unchanged in substance, plus: the definition must state **nullability** and the difference between "none exists" and "not recorded" (SD-19) |
| **D-2** (K-R15 wording) | **constrained**: any rewording must cover statement-less events. "Ordered by a persisted position" is still right in purpose; "position within the turn" is no longer sufficient on its own |
| **D-3a** (provenance defect) | **strengthened** — see §7; the defect now violates SD-04 directly |
| **D-3b** | **settled by this record** |
| **D-4** (reconstruction) | unchanged: statement-less events must be reproducible too, which a fold over a universal key supports and a span-position fold does not |
| **ND-1** (status fold) | unchanged, still separate |
| **DN-6a / DN-6b** | **still open**. The decision says events must record their cause; it does **not** say who may cause them |

---

## 7. The implementation defects, seen again

**SD-20 (Observed — `model.py:303-312`).** `revoke()` stamps its event with the **acquisition's**
`source_span`. Under SD-04 this is no longer merely a provenance error: it **manufactures conversational
evidence for an event that had its own cause** — the exact substitution the decision forbids.

**SD-21 (Observed — `model.py:276`, `:310`).** `event_id = f"ev-{len(self._events) + 1:05d}"` is derived
from append position. Such an event is **not independently identifiable**: its identity changes when the
log is rebuilt, and two runs of the same history produce different ids. That fails SD-03 as directly as it
fails K-R2.

**SD-22 (Inferred).** Both defects are now backed by an owner decision rather than by an argument from
internal coherence, which makes them unambiguous — and still **unauthorised to fix**. Repairing them is
separate work.

---

## 8. Why this decision does not settle DN-6a

**SD-23 (Inferred).** There is a tempting slide: "events may happen without a statement, so the system may
forget things on its own, so the system may decide to". Each step is invalid. SD-02 requires the event to
**record** what caused it; recording a cause presupposes that *something* caused it, not that the
something is authorised. **Who may cause forgetting remains undecided**, and O-5 (belief, model-owned) is
still not a licence to extend into knowledge removal.

---

## 9. Proposed canonical changes — identified, NOT made

| Document | Proposed change | Reason |
|---|---|---|
| `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` | add **O-16** recording D-3b with the verbatim wording and §1.1's four items | owner decisions belong in the ledger; **not authorised**, not made |
| r33 §6.2 | none proposed — `evidence_span` nullability and the cause descriptor are consequences to be shaped, not edits to make | no schema change is authorised |
| `MEMORY-09-CONFORMANCE-MAP` | **none** | the 118 are unchanged; SD-02's cause requirement is an obligation without a row, the same coverage situation as NEW-1 |

**SD-24 (Inferred).** SD-02 creates a second in-scope obligation with no conformance row — the same
species as NEW-1, not the D-4 species. Recorded, not acted on.

---

## 10. Verification

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema or code edited | **0** |
| Tests / probes run | **0** |
| v9 / r33 / ledger / map | `698254ca…` · `dfb28b60…` · `e5008020…` · `b4ba33cc…` — unchanged |

**End of record.** D-3b is settled and is not to be re-decided. DN-6a and DN-6b remain open.
