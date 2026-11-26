# MEMORY-09 — DN-4 Preparation: The Forgetting Scope Decision

**Status:** research report and decision preparation. **READ-ONLY.** No architecture selected, no
implementation, no test, no probe, no model call, no corpus, no canonical edit, no registry change, no
Git operation. DN-1 is **not** confirmed here and DN-4 is **not** decided here.

| Field | Value |
|---|---|
| Created | **2026-11-19** (local, CST) |
| Purpose | Prepare DN-4 as the next owner decision; verify the No path requirement by requirement; keep DN-1 separate |
| Baseline | v9 `698254ca262c0874` · r33 `dfb28b6003ca356f` · ledger `e5008020ff68bd21` · conformance `b4ba33ccbfe8f8c1` · CO `82173a6f78ac84c4` |
| Tests executed | **0** |
| Files changed | One new file. Untracked, unregistered, unfrozen |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. Two items, deliberately not bundled

**FD-01 (Inferred).** DN-1 and DN-4 are different kinds of question and must not be approved together:

| | **DN-1** | **DN-4** |
|---|---|---|
| Kind | textual: which antecedent governs a sentence | behavioural: whether an operation exists |
| Answer form | confirmation (one sentence) | decision (yes or no) |
| If left open | an ambiguity persists in a ratified document | a non-conforming behaviour continues |
| Authority | a reading of existing text | a new owner decision |
| Consequence of the alternative | record an internal inconsistency in v9 | a scope determination with registered exclusions |

Bundling them would let a textual confirmation appear to authorise a behavioural choice, or the reverse.

---

## 2. DN-1, presented for confirmation

**FD-02 (Source-derived).** The narrow reading is supported by four passages, and the absolute reading
conflicts with all four:

1. **Rule 12**, titled "Provenance without bypass": the prohibition is "never access to **the source
   record**".
2. **§8**: the clause sits inside a list of what must not be obtained *from the source* — "the source
   record, its unrelated contents, its derivation history, or any pointer resolving to them".
3. **§7 step 4**: the account is created with "source, **event/turn**, communicated content, origin
   descriptor" — an event reference is required content.
4. **Rules 50–51**: `origin_class` "answers whether the account is **backed by an acquisition event**",
   and an own account "creates **or implies no acquisition event**" — the distinction presupposes
   accounts that are backed by one.

**FD-03 (Inferred).** The request is a **confirmation**, not an amendment: that §8's "never a view,
alias or projection" bars a view of the source record, not derivation from the recipient's own
acquisition events. Confirming it means no PH1A-03 divergence exists. Declining it means v9 contradicts
itself across rules 12 and 50–51 and must be reconciled by amendment before anything downstream is
coherent.

**FD-04 (Inferred).** DN-1's status remains **open pending the owner's word**. The evidence is strong;
the authority to settle it is not this record's.

---

## 3. Why DN-4 is a real question, not a derived one

**FD-05 (Source-derived).** The hierarchy settles how knowledge-loss must be represented **if it
exists** — K-R3 (a new event, never a deleted row), K-R19 (a forgotten pair has no edge) — and does not
settle **whether the application offers it**:

- no T1 decision requires a character to be able to lose knowledge; O-2 (lifecycle consequences) is
  **open**;
- **0 of the 118** conformance requirements test character knowledge-loss; the `revoke` rows concern
authority, belief, or a teller's revocation that leaves the recipient knowing (C1) — see FD-10;
- T2 presupposes the operation in K-R3, K-R17, K-R19, K-R24/K-R25 and probes S37/S41, but presupposing
  a mechanism is not mandating that the product expose it.

**FD-06 (Inferred).** That gap — representation settled, existence unset — is exactly the shape of a
genuine owner decision. It is not a storage preference dressed as one.

---

## 4. The question, worded to avoid two traps

> **Can a character lose knowledge — forget something they once knew?**

**Yes.** Forgetting is a **distinct operation** from retraction: it ends what the character holds while
the record and its history remain intact, and it is recorded as its own event on the knowledge stream.

*What Yes establishes:* the behavioural distinction. It does **not** settle the event shape, the
retained-history representation, the ordering scheme, or who may perform the operation — those follow
from K-R3/K-R19 and from DN-6, not from this answer.

**No.** MEMORY-09 has **no knowledge-removal operation**. Records and accounts persist; content status
continues to change independently through retraction, correction, supersession and authority changes.

*What No does not mean:* that forgetting and retraction are interchangeable. It means forgetting does
not exist as an operation, so there is nothing to conflate. Retraction keeps its own meaning (the
content is no longer current) and the knower keeps what they acquired.

**FD-07 (Inferred).** The distinction between those two "no"s matters more than it looks. A "No" that
were read as "retraction and forgetting are the same thing" would be the current defect raised to
policy; a "No" read correctly removes one operation and leaves the other's semantics untouched.

---

## 5. The No path, verified requirement by requirement

The review asked for this specifically: which requirements become **inapplicable**, which remain
**applicable unchanged**, and which remain applicable but need **a different mechanism**. The answer is
not uniform, and the list is shorter than the earlier report implied.

### 5.1 Inapplicable under No — the mechanism each rule constrains would not exist

| Rule | What it constrains | Why it is inapplicable |
|---|---|---|
| **K-R3** | revocation as a new event; "forgetting … never rewrites what happened" | no revocation of knowledge occurs, so there is no such event to constrain. The historical-integrity principle is **independently carried by I4** ("history is never destroyed"), not lost with the rule |
| **K-R17** | a revoke is valid only when the log ends in a learn event; invalid revokes rejected at write | no revoke operation to validate |
| **K-R19** | a forgotten pair has **no edge**, not an edge with a flag | no forgetting, so no absent-edge state |
| **K-R20** | `revoked_turn` is derived, not a field | no revoke to derive it from |
| **K-R25** | forgetting and difficulty-recalling distinguishable in every view | only the second state exists |
| **K-R11** (partially) | the closed event set includes `revoked` | the value becomes unused — a dead member of a closed set, not a violation |
| **S37** | learn → revoke → relearn multi-cycle | the probe cannot be run as written |
| **S41** | forgetting vs recall-difficulty distinguishable | one state |

**FD-08 (Calculated).** Seven of these are T2 rules or T2 contract probes; **none is among the 118**
ratified requirements. S37 and S41 are r32 §12 probes, never executed.

### 5.2 Still applicable, unchanged — no rule becomes a different mechanism

| Rule / requirement | Why it survives No |
|---|---|
| **K-R1, K-R2, K-R11** (fold, rebuildability, event-family fold) | acquisition history is unaffected by whether knowledge can be removed |
| **K-R14** (two streams rebuilt independently) | unaffected |
| **K-R24** (recall difficulty is surfacing only; must not hide a record) | constrains *surfacing*, which still exists; the "never by removing an edge" clause becomes vacuous, not the prohibition |
| **K-R15** (`turn, seq`) | unaffected — and unsatisfiable in either implementation regardless |
| **F1 / F1b** (reconstructability) | unaffected |
| **D5** (decay affects surfacing only) | unchanged — and it is *why* No is viable: decay never needed forgetting |
| **A2a** ("still holds what it acquired") | **strengthened** — nothing can remove an account, so the guarantee cannot be broken |
| **E1a, E1b, H10, O-1, O-3, O-5, O-14** | untouched |

**FD-09 (Source-derived).** **No rule falls into the third category** — "applicable but requiring a
different mechanism" — in any strong sense. The applicable rules keep their mechanisms; the inapplicable
ones are not replaced. This is a correction to the earlier report, which listed the exclusions without
distinguishing them.

**FD-10 (Calculated).** No ratified requirement is broken by No. **Zero** of the 118 test character
knowledge-loss: `forget` returns no row, and every `revoke` match is a different subject — authority
(B7c, K26), belief (D4, D5, governed by O-5), and **C1**.

**C1 is the closest the map comes to this question, and it points the same way**: *"Sable told by
Alaric; Alaric later revoked → Sable still knows; origin told, teller cited."* The 118 already treat a
source's revocation and a knower's knowledge as separate things. A2a is strengthened rather than
weakened.

**No is conformant with the 118** and departs only from T2 mechanisms that presuppose an operation the
owner would have declared out of scope.

### 5.3 How that departure should be recorded — two readings, unresolved

**FD-11 (Unestablished).** Whether No is a **case-2 divergence** (v9 contradicts an applicable r32
mechanism) or a **D-4-style scope determination** (MEMORY-09 is simply absent from the operation, as it
is from transfer) is not settled by O-6a:

- **Case 2** fits if K-R3/K-R19 are treated as applicable obligations that a No decision sets aside;
- **D-4** fits better if the operation is never in MEMORY-09's scope, since no ratified requirement
  needs it and O-2 is open. D-4 set the precedent: K-R46–K-R48 were recorded as a **known coverage
  gap** and a scope determination, not as divergences.

The owner's answer should say which form it takes; the two have different registry consequences.

---

## 6. The third position is a divergence, not a default policy

**FD-12 (Observed).** Both implementations conflate the two operations today:

- the engine's `retract()` revokes a private account, and `perspective()` then hides it from its owner
  (§7.5 of the CO investigation) — knowledge is removed, with no knowledge-stream event;
- the harness's `retract()` revokes the target record and writes a retraction record — the same
  removal, plus a content-status statement.

**FD-13 (Inferred — the point worth preserving).** The current behaviour satisfies **neither** branch:

- it is **not "Yes"**, because the removal produces no knowledge-stream event, no episode closure and
  no distinct record of the loss;
- it is **not "No"**, because it removes knowledge, which No forbids outright.

So leaving DN-4 unanswered is not neutral and is not an implicit choice of either answer. It continues a
**current implementation divergence** that is non-conformant under both. That is a different status from
"we have been behaving as if No were true", and the record must not blur them.

---

## 7. What DN-4 does not settle

**FD-14 (Inferred).** Even a complete answer leaves these open, and none should be read into it:

| Item | Where it comes from instead |
|---|---|
| who may cause forgetting, and who may retract content | **DN-6**, which splits under Yes and dissolves under No |
| the event shape and ordering of a forgetting event | K-R3/K-R15 |
| what a "no longer holds it" state looks like in an inspection view | P-5 (unbuilt, unauthorised) |
| whether the experiment's subject expresses the operation | **PH1A-12**, after the mechanism set is fixed |

---

## 8. Sequence

**FD-15 (Inferred).** Two separate owner responses, in this order, with nothing else attached:

1. **DN-1 — confirmation.** One sentence: the narrow reading, or an instruction to record v9's internal
   inconsistency for reconciliation.
2. **DN-4 — decision.** Yes (distinct operation) or No (no removal operation), with the §5.3 form
   (divergence or scope determination) named.

Then, and only then: **DN-6** (authority, shaped by the DN-4 answer), then **PH1A-12** (subject, once the
mechanism set is fixed), then the corpus and the bridge — none of which is authorised by this record.

**FD-16 (Superseded).** The earlier report's enumeration of No-path exclusions is superseded by §5's
three-way classification; the exclusions stand, but their character differs rule by rule.

---

## 9. Limits

**FD-17 (Observed).** All findings are static readings at the hashes in the header. Nothing was
executed. The claim that 0 of the 118 mention forgetting is a search result, not an execution.

No architecture selected, no implementation, no test, no probe, no model call, no corpus. No canonical
edit, registry change or Git action. No execution boundary moved: implementation, testing and probes
remain unauthorised.

---

**End of report.** DN-1 awaits confirmation; DN-4 awaits a decision. Neither is answered here.
