# MEMORY-09 — Owner Decisions DN-1 and DN-4

**Status:** owner-decision record. **DN-1 confirmed. DN-4 = YES.** Documentation and read-only analysis
only. **No implementation, test, probe, benchmark, model call or corpus execution is authorised by this
record.**

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Decisions | **DN-1** (v9 §8 textual clarification) · **DN-4** (forgetting = YES) |
| Authority | Owner, stated in session and recorded here; as with every prior confirmation, **no persisted signature exists** |
| Canonical edits | **None.** v9, r33, the owner ledger and the conformance map are unchanged (§10 lists the proposed changes that were *not* made) |
| Files created | This record and one reconciliation report |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. DN-1 — v9 §8 textual clarification (CONFIRMED)

### 1.1 Exact owner wording — preserved verbatim

> I confirm that v9 §8's prohibition on obtaining a view, alias, or projection applies to the source record itself. It does not prohibit a recipient from maintaining an independently owned account backed by their own legitimate acquisition event.

### 1.2 What this confirms

**OD-01 (Source-derived).** The prohibition concerns obtaining a view, alias or projection **of the
source record**.

**OD-02 (Source-derived).** A recipient may maintain an **independently owned account backed by their
own legitimate acquisition event**.

**OD-03 (Source-derived).** The confirmation does **not** grant access to the source record, its
unrelated contents, or its private derivation history.

**OD-04 (Source-derived).** The recipient's account remains **independently owned** and must retain
**its own provenance**.

### 1.3 What it is not

**OD-05 (Observed).** This is **not an amendment to v9.** The owner confirmed the interpretation of
existing wording; no rule text changed. The v9 hash is unchanged at `698254ca262c0874`.

**OD-06 (Observed).** No historical artifact was rewritten to make the interpretation appear previously
recorded. The Phase 1A-1D and DN-1 preparation reports retain their "open pending confirmation" status
as records of the state at the time.

### 1.4 Evidence the reading rests on (recorded, not re-argued)

| Passage | What it shows |
|---|---|
| **Rule 12**, "Provenance without bypass" | the prohibition is "never access to **the source record**" |
| **§8** | the clause sits in a list of what must not be obtained *from the source* |
| **§7 step 4** | the account is created with "source, **event/turn**, communicated content, origin descriptor" |
| **Rules 50–51** | `origin_class` "answers whether the account is **backed by an acquisition event**"; an own account "creates **or implies no acquisition event**" |

**OD-07 (Inferred).** The confirmation is consistent with all four; the absolute reading conflicted with
all four and would have required recording an internal inconsistency in v9. That alternative is now
foreclosed.

### 1.5 Effect on PH1A-03

**OD-08 (Source-derived).** The **PH1A-03 divergence is removed**. There is no T2/T3 conflict on this
point: v9 does not prohibit an account being backed by, and reconstructable from, the recipient's own
acquisition events.

**OD-09 (Source-derived).** K-R1, K-R2 and K-R11 therefore govern, attached to the account as the
carrier of the knowledge relation. The owner's wording — "account backed by their own legitimate
acquisition event" — supplies the carrier phrasing that the earlier analysis marked `Unestablished`.

**OD-10 (Inferred).** Nothing downstream needs to change except the status of PH1A-03 and the
findings listed as superseded in the reconciliation report. No other conclusion depended on the
absolute reading.

---

## 2. DN-4 — Forgetting (YES)

### 2.1 The decision

**OD-11 (Source-derived).** **YES.** A character can lose knowledge they previously acquired, through a
**distinct forgetting operation**.

**The decision establishes** — exactly these six things:

1. Forgetting is **distinct from retraction**.
2. Forgetting is a **new event**.
3. The **account and its history remain preserved**.
4. The character **no longer retains the knowledge as usable knowledge** after the operation.
5. The **historical fact that the character once acquired the knowledge remains represented**.
6. Forgetting **must not** be implemented as an ordinary content-status change, and must not be
   silently conflated with retraction.

**The decision does not establish** — and these remain unresolved:

- the exact event schema or payload;
- event ordering or sequence representation;
- who is authorised to cause forgetting;
- whether the character can initiate forgetting;
- whether forgetting is reversible;
- the precise retained-history or projection structure;
- the relationship between forgetting and every other lifecycle operation;
- any implementation strategy or code change.

**OD-12 (Inferred).** Items 2 and 4 together are the whole behavioural content: the knowledge relation
ends, and it ends by its own recorded event rather than by mutating the object that carries the content.

### 2.2 Why this is the answer the sources point at

**OD-13 (Source-derived).** T2 already specifies the mechanism in detail — K-R3 (a new event, never a
deleted row), K-R17 (a revoke is valid only when the log ends in a learn event), K-R19 (a forgotten pair
has **no edge**), K-R20 (`revoked_turn` derived, not a field), K-R24/K-R25 (forgetting and difficulty
recalling distinguishable in every view) — and probes S37/S38/S41 test exactly this. The 118 contain no
forgetting row, which is why the operation's **existence** needed an owner answer rather than a
derivation.

---

## 3. Consequences of DN-4 = YES, requirement by requirement

### 3.1 Applicable and now required — every one individually

| Requirement | Consequence of YES |
|---|---|
| **K-R1** | the knowledge relation is the fold of that pair's acquisition events, and **no relation exists without at least one event**; a forgetting event closes the current episode → **fully applicable, now load-bearing** |
| **K-R2** | a projection that cannot be rebuilt from its events is an **integrity failure**; the projection must now rebuild with forget events included |
| **K-R3** | revocation **is** a new event and never a deleted history row — this is the rule YES invokes directly |
| **K-R11** | the fold depends only on acquisition-family events; the `revoked` member becomes **live rather than dead** |
| **K-R14** | each stream rebuilt and checked independently; the acquisition stream must now be checked with forget events present |
| **K-R15** | events ordered by `(turn, seq)`; **unchanged in text, materially harder** — a forget and a learn in the same turn must resolve deterministically (§4) |
| **K-R16** | event ids unique; unchanged, but a second event kind now depends on it |
| **K-R17** | a revoke is valid only when the log currently ends in a learn event; **invalid forgets must be rejected at write** |
| **K-R18** | a learn while an episode is open is a continuation; unchanged |
| **K-R19** | a forgotten pair has **no edge** — not an edge with a flag |
| **K-R20** | `revoked_turn` is **derived**, not a field |
| **K-R21 / K-R90** | `via` is a summary of the current episode; "how did they learn this" is answered from events — **the event log must survive forgetting** (this is establish-item 5) |
| **K-R22 / K-R23** | continuation may upgrade, never re-date; `via` never downgrades |
| **K-R24** | recall difficulty is a **surfacing** property, never eligibility; it must not hide a record from a direct question — **now must be kept distinct from K-R19's absence** |
| **K-R25** | forgetting and difficulty-recalling **remain distinguishable in every view that reports either** — required directly by YES |
| **K-R26** | `via` precedence is a fixed total order; unchanged |
| **I4 / I11** | history never destroyed; acquisition history authoritative — both now constrain the forget path explicitly |
| **I15** | absence is never evidence about absence; unchanged |
| **F1** | knowledge reconstructable from acquisition and coverage records — **now includes forget events** |
| **F1b** | reconstruction compared against any projection; mismatch an integrity failure — **now has content to check**, where before there was nothing |
| **S37** | learn → revoke → relearn multi-cycle: **required**, and currently unrunnable |
| **S38** | invalid revoke rejected at write; same-turn resolution by `seq`: **required**, and currently unrepresentable |
| **S41** | forgetting vs recall difficulty distinguishable: **required** |

### 3.2 Scope clarified — applicable, but must not be over-read

| Item | Clarification |
|---|---|
| **A2a** ("Character still holds what it acquired") | the row's context is *"Account; memory later retracted"*. It governs **source retraction**, and is **not** contradicted by YES: forgetting is a different operation with a different cause. A2a must not be read as "a character can never lose knowledge" |
| **C1** ("Sable told by Alaric; Alaric later revoked → Sable still knows") | same point at row level: a **source's** revocation does not remove a recipient's knowledge. It is evidence for the separation YES establishes, not against it |
| **K-R13** (read-time join) | unaffected in text; its *use* is now inspection-only, because O-14 removed the character-facing notice. The join must not be computed for a character's channel |
| **O-2** (chat-deletion lifecycle) | remains **open**; deletion (K-R43, removal of records with the chat) is **not** forgetting and must not be merged with it |
| **D5** (decay affects surfacing only) | unchanged, and it is what makes the separation clean: decay never needed forgetting, and forgetting is never represented by decay |

### 3.3 Remains unmet — by the implementations, not by the requirement

**OD-14 (Observed).** YES changes no requirement text, so nothing that was met becomes unmet. It does
make several existing obligations testable for the first time, and **no implementation meets them**:

| Requirement | v9 harness (`model.py`) | Engine (`m09/`) |
|---|---|---|
| K-R1/K-R2/K-R14 (fold, rebuild, two streams) | events exist; `edges()` maps 1:1 and never folds | **no events at all**; status is a field |
| K-R3 (forget is an event) | `revoke()` sets record status **and** appends a `via="revoked"` event — a hybrid | no event |
| K-R17 (reject invalid revoke) | not checked | no operation |
| K-R19 (no edge after forget) | **violated** — one edge per event, including `revoked` | no edges |
| K-R20 (`revoked_turn` derived) | `revoked_at_turn` is a field | no field, no derivation |
| K-R25 (two states distinguishable) | not distinguishable | not distinguishable |
| K-R15 (`turn, seq`) | `turn` only | no events |

---

## 4. Conformance rows and probes

**OD-15 (Calculated).** **0 of the 118 requirements change text, and none is added or removed.** YES
creates no conformance row. Three effects are worth recording:

| Row | Effect |
|---|---|
| **F1 / F1b** | unchanged in text; **stricter in practice** — reconstruction must now include forget events |
| **A2a** | unchanged; **clarified** as source-retraction-only (§3.2) |
| **E1a / E1b / H10** | unchanged — the retraction path is untouched by YES |

**OD-16 (Calculated).** **Three T2 probes become directly relevant and remain unrunnable**: S37, S38,
S41. None is among the 118; all are r32 §12 probes, and none has ever been executed.

**OD-17 (Inferred — a genuine gap).** YES establishes a T2 obligation with **zero coverage in the 118**.
Whether that is acceptable (the 118 are the experiment's requirements, and forgetting may be out of
their scope) or whether a row must be added is **not settled by existing authority**: O-6b closed the
set at 118, and adding a requirement needs an owner decision. This is raised in §9, not decided.

---

## 5. The current implementation divergence — preserved

**OD-18 (Observed).** Both implementations conflate the operations, and **neither now matches either
branch**:

- the engine removes knowledge (a record status change hides the account from its owner,
  `m09/retrieval.py:50`) with **no knowledge-stream event** — that is Yes's requirement missing and
  No's persistence violated;
- the harness revokes the target record and writes a retraction record — knowledge removal plus a
  content-status statement, again with **no distinct forgetting event**, and with an edge still emitted.

**OD-19 (Inferred).** This is a **current implementation divergence**, not an implicitly selected
policy. It must not be described as "the system has been behaving as though forgetting existed": it has
been behaving as though forgetting *and* retraction were one status mutation, which is neither answer.
Under YES it is now non-conformant in a newly specified way.

---

## 6. Downstream: DN-6

**OD-20 (Inferred).** DN-4 = YES **splits** the authority question into two, and they are not
interchangeable:

| | **DN-6a — who may cause forgetting** | **DN-6b — who may retract content** |
|---|---|---|
| What it governs | ending a character's knowledge relation | changing a record's truth status |
| Existing authority | **none** | **none** for claims/accounts; O-10b governs *promotions* only |
| Closest principle | O-5: the model owns the character's mental state; MEMORY-09 records it | O-1 ownership; O-10b human control of authority changes |
| Status | **open** | **open** |

**OD-21 (Inferred).** O-5's principle points toward the character's **model** holding forgetting
authority, but O-5 governs *belief*, not knowledge, and no text extends it. That is an inference, not a
decision. Whether a character may initiate their own forgetting, and whether the operator may impose it,
remain open.

**OD-22 (Source-derived).** Nothing in DN-4 answers DN-6a or DN-6b. The decision explicitly leaves
authorisation unresolved.

---

## 7. Downstream: PH1A-12

**OD-23 (Inferred).** The **mechanism set is now fixed at the behavioural level**: two streams (DN-3,
by tier precedence), acquisition events backing accounts (DN-1 + K-R1/K-R11), forgetting as a distinct
event (DN-4), status as a projection.

**OD-24 (Observed).** No implementation expresses that set, and the `seq` schema gap remains
(ARCH-03 F6; PH1A-05 open). So PH1A-12 — which implementation is the experimental subject, and where
the bridge lives — is **answerable in principle but not yet answerable well**, and this record does not
select an implementation or a bridge design.

**OD-25 (Inferred).** PH1A-12's predecessor list is now: **PH1A-05** (`turn, seq` structure) and the
DN-6a/DN-6b authority answers (needed for conformance, not strictly for a scripted experiment). DN-1,
DN-3 and DN-4 are discharged.

---

## 8. What these decisions authorise

**OD-26 (Observed).** Documentation only. No implementation, no code change, no test, probe, benchmark
or experiment, no model call, no corpus, no event-schema design presented as settled, no registry edit,
no change to v9, r33, the ledger or the conformance map, no modification of a frozen artifact, no Git
operation.

---

## 9. Genuinely new questions raised by the evidence

Only one, and it is not answerable from existing authority:

**NEW-1 — Does the 118-requirement set need a forgetting row?** YES creates a T2 obligation with no
coverage in the ratified set (OD-17). Options: treat forgetting as outside the experiment's requirement
scope, recorded as a coverage gap like K-R46–K-R48 under D-4; or authorise one additional requirement.
O-6b closed the set at 118, so this cannot be resolved by inference.

**NEW-2 — Is PH1A-05 now blocking?** YES makes same-turn ordering material for the first time (a forget
and a learn in one turn). The `seq` structure question was already open; it now has a consequence.

Neither is DN-1 or DN-4 being re-opened.

---

## 10. Proposed canonical changes — identified, NOT made

| Document | Proposed change | Reason |
|---|---|---|
| `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` | add an **O-15** section recording DN-1 and DN-4 with the verbatim wording and the §2/§3 scope boundaries | the ledger is the durable home of owner decisions; these are owner decisions |
| `MEMORY-09-CONFORMANCE-MAP-2026-11-17.md` | **none proposed** unless NEW-1 is answered "add a row" | the 118 are unchanged by DN-4 |
| `MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` | **none** | DN-1 confirms existing wording; it does not amend v9 |
| `…-ARCH-STAGE-CONTRACT-…` (r33) | **none** | DN-4 selects K-R3/K-R19 as governing; it changes no r32 text |

These edits are **not authorised by this record** and were not made.

---

## 11. Verification

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Implementation files edited | **0** |
| Tests / probes run | **0** |
| v9 hash | `698254ca262c0874` — unchanged |
| r33 hash | `dfb28b6003ca356f` — unchanged |
| Ledger hash | `e5008020ff68bd21` — unchanged |
| Conformance hash | `b4ba33ccbfe8f8c1` — unchanged |

**End of record.** DN-1 and DN-4 are settled and must not be re-decided.
