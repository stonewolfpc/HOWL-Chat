# MEMORY-09 — DN-1 Textual Clarification and the Forgetting Scope Question

**Status:** research report. **READ-ONLY.** No architecture selected, no implementation, no test, no
probe, no model call, no corpus, no canonical edit, no registry change, no Git operation.

| Field | Value |
|---|---|
| Created | **2026-11-19** (local, CST) |
| Purpose | (a) perform DN-1 as a tightly scoped textual check; (b) state DN-2/DN-3 as consequences of the hierarchy; (c) frame the one genuinely open behavioural question |
| Baseline | v9 `698254ca262c0874` · r33 `dfb28b6003ca356f` · ledger `e5008020ff68bd21` · conformance `b4ba33ccbfe8f8c1` · Phase 1A CO `82173a6f78ac84c4` |
| Tests executed | **0** |
| Files changed | One new file. Untracked, unregistered, unfrozen |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. What DN-1 asked, and what the check found

**TC-01 (Observed).** DN-1 is the question whether v9 §8's final clause — *"Each recipient record has
its own identity and acquisition provenance; it is never a view, alias or projection"* — prohibits
B's record being a view of **A's record**, or prohibits it being derived from **B's acquisition
events**. Only the second reading creates a conflict with K-R1.

**TC-02 (Source-derived).** The check was run on the whole of v9's disclosure material, not on the
sentence in isolation. The result is stronger than the decomposition's grammatical argument: **the
absolute reading contradicts two of v9's own rules.**

---

## 2. The textual evidence

### 2.1 Rule 12 is about access to the source, in its own title and content

**TC-03 (Observed).** v9 rule 12, *"Provenance without bypass"*:

> "recipient records carry a restricted origin descriptor, **never access to the source record**."

The rule's prohibition is access to the source record. It says nothing about derivation.

### 2.2 §8 is a list about the source, and the clause sits inside that list

**TC-04 (Observed).** v9 §8:

> "**Recipient must not obtain:** the source record, its unrelated contents, its derivation history, or
> any pointer resolving to them.
>
> **Mechanism:** the recipient's origin reference is a **descriptor** (speaker, record type, turn), not
> a resolvable identifier. The resolvable link exists only on the source side. Each recipient record
> has its own identity and acquisition provenance; it is never a view, alias or projection."

"View, alias or projection" is the last item in a paragraph whose subject is *what must not be obtained
from the source*. Read absolutely, the clause would leave the paragraph's own topic — pointers to the
source — and prohibit a different relation the paragraph never mentions.

### 2.3 v9 requires the account to carry an event reference

**TC-05 (Observed).** v9 §7, step 4:

> "**Memory system** creates, per exposed character, their own acquired account: source, **event/turn**,
> communicated content, origin descriptor."

An account is required to record an **event/turn**. If accounts could not be tied to acquisition
events, that field would have no referent. Rule 49 reinforces it: *"the disclosure **event** and the
recipient's acquired account persist after any retraction of the source."*

### 2.4 The decisive one: v9 already distinguishes accounts *by* their acquisition event

**TC-06 (Source-derived — the finding of this pass).** v9 rule 50 and rule 51 use the concept the
absolute reading would prohibit:

| Rule | Text | What it presupposes |
|---|---|---|
| **50** | "`via` answers how information entered the character's knowledge; `origin_class` answers whether the account is **backed by an acquisition event**" | that some accounts *are* backed by acquisition events |
| **51** | "a character's own account … creates **or implies no acquisition event**" | that the negative is meaningful — i.e. that other account kinds do |

A document that defines one class of account as "backed by an acquisition event" and another as
"creating or implying none" cannot also mean that no account may be derived from an acquisition event.
**Under the absolute reading, v9 would contradict itself**, with no stated priority — the same defect
shape O-14 had to resolve in r32, in a document where no such resolution exists.

---

## 3. The two readings, and what each costs

| | **Narrow reading** (antecedent = the source record) | **Absolute reading** (no derivation) |
|---|---|---|
| Consistent with rule 12's title and content | yes | no — the rule is about access to the source |
| Consistent with §8's list | yes | no — leaves the paragraph's subject |
| Consistent with §7 step 4 ("event/turn") | yes | no — the field becomes meaningless |
| Consistent with rules 50–51 (`origin_class`) | yes | **no — direct contradiction** |
| Consistent with K-R1/K-R11 | yes | no |
| Cost if adopted | none identified | requires reconciling rules 12, 50, 51 and §7 step 4, which is an amendment-shaped act |

**TC-07 (Source-derived).** The narrow reading is the only one under which v9 is internally coherent.
This is not a preference between two equally supported readings; the absolute reading has **no**
supporting text and conflicts with four passages — rule 12, §8, §7 step 4, and rules 50–51.

**TC-08 (Inferred).** It follows that PH1A-03's alleged T2/T3 divergence **does not exist on the
coherent reading**. There is no case-2 divergence to decide, because there is no contradiction to
resolve.

---

## 4. What this does *not* establish

**TC-09 (Inferred).** The finding is a **reading requiring confirmation**, not a ruling and not an
amendment:

- It does not amend v9. v9 is ratified (O-6b); confirming which of its sentences governs is not
  editing it.
- It does not establish that accounts *must* be event-derived. It establishes only that v9 does not
  forbid it — which is what PH1A-03 needed.
- It preserves the ambiguity in the record. DN-1 remains open until the owner confirms, because the
  consequence is architectural even though the question is textual.

**If the owner declines the narrow reading**, the consequence is not "PH1A-03 becomes a divergence". It
is that **v9 acquires an internal inconsistency** between rules 12 and 50–51, which would have to be
recorded as a defect in v9 and reconciled by amendment. That is a heavier outcome than the one-sentence
confirmation, and the owner should see it before choosing.

---

## 5. DN-2 as a consequence, not a choice

**TC-10 (Source-derived).** On the narrow reading, K-R1, K-R2 and K-R11 apply without any owner
decision:

- **K-R1** — the knowledge relation is the fold of that pair's acquisition events, and no relation
  exists without at least one event.
- **K-R2** — a projection that cannot be rebuilt from its events is an integrity failure.
- **K-R11** — the fold depends only on acquisition-family events.

**TC-11 (Observed + Inferred).** One qualification, stated rather than glossed: K-R1's text says
"the **edge** for (character, memory)". v9 has no edges, and O-1 requires B's record to be a separate
account (T1, and it displaces r32's shared-reference relation for private disclosure — CO §7.2). So the
**substance** of K-R1/K-R11 attaches to whatever carries the knowledge relation; under v9 that carrier
is the account. Whether that transposition is itself a "gap-fill" (O-6a case 3) or simply the
representation consequence of O-1 is **`Unestablished`** and is not resolved here. What is established
is that the obligation survives; what is not is whether the surviving form needs to be declared as new
architecture.

---

## 6. DN-3 as a consequence, not a choice

**TC-12 (Source-derived).** K-R12 requires status, scope and class to be projections of a lifecycle
stream. v9 rules 20, 37 and 42 describe the *effects* of lifecycle operations and say nothing about
whether status may be a stored field.

**TC-13 (Inferred).** This is the **fourth combination** O-6a's three cases do not name — r32 speaks,
v9 is silent. No case is required: plain tier precedence decides it, because T2 governs where it is in
scope and does not conflict with T1. **PH1A-04 is therefore not an owner decision.** A decision would be
required only to *depart* from K-R12, and nothing in the record motivates that departure.

---

## 7. The three categories, kept apart

**TC-14 (Inferred).** The distinction the review asked for, applied:

| Category | Items | How it must be recorded |
|---|---|---|
| **A. Settled by existing authority** | K-R12 governs (T2, v9 silent); K-R1/K-R2/K-R11 govern on the narrow reading; O-14; A2b/E1c/D2 | as findings of tier precedence, **not** as decisions |
| **B. Interpretation requiring confirmation** | DN-1 — the antecedent of "never a view, alias or projection" | as a reading with its evidence, **not** as an amendment and **not** as a decision |
| **C. Genuinely new owner decision** | DN-4 — whether knowledge-loss is in scope | as a behavioural decision, §8 below |

**TC-15 (Inferred).** Two failure modes this separation prevents: DN-1 hardening into an architectural
rule without confirmation, and DN-3 being read as an optional design preference when it is an
obligation.

---

## 8. DN-4 — the one genuinely open behavioural question

**TC-16 (Source-derived).** The hierarchy settles **how** knowledge-loss must be represented if it
exists (K-R3: a new event, never a deleted row; K-R19: a forgotten pair has no edge, not an edge with a
flag). It does not settle **whether the application must offer it**:

- No T1 decision requires a character to be able to lose knowledge. O-2, the lifecycle-consequences
  question, is **open** and explicitly out of experiment scope.
- None of the 118 conformance requirements has a forgetting row.
- T2 presupposes the operation (K-R3, K-R17, K-R19, K-R24/K-R25; probes S37, S41), but presupposing a
  mechanism is not mandating that the product expose it.

**So the question is scope, not representation.** Stated plainly:

> **Can a character lose knowledge — forget something they once knew?**
>
> - **Yes** — then K-R3 and K-R19 apply as written: forgetting is a new event on the knowledge stream,
>   the pair ends up with no relation, and the record and its history remain. It is distinct from
>   retraction, which changes truth status and leaves the knower holding what they acquired.
> - **No** — then no operation removes knowledge. Accounts persist; only content status changes.
>   K-R3, K-R17, K-R19, K-R24/K-R25 and probes S37/S41 become **scope exclusions**, recorded the way
>   K-R46–K-R48 were recorded under D-4, rather than divergences.

### 8.1 What each answer locks

| Consequence | **Yes** | **No** |
|---|---|---|
| A2a ("still holds what it acquired") | holds unless a forgetting event occurs — an operation can make it false | **unbreakable** — nothing can remove an account |
| The engine's `retract()` on a private account | non-conforming either way (it neither forgets with an event nor preserves) | non-conforming; no part of that operation could be retained as-is |
| RQ-03 / DN-6 | **splits**: who may cause forgetting vs who may retract content | **dissolves**: nobody may revoke an account |
| K-R3/K-R17/K-R19 | apply in full | vacuous, as a recorded exclusion |
| CO-07's token collision | resolved by keeping `via=revoked` and `status=retracted` apart | resolved by having only one status operation |
| Narrative capability | characters can forget | characters never forget; only content changes |

### 8.2 The status quo, stated so it is not chosen by default

**TC-17 (Observed).** Both implementations already behave like a third position — forgetting conflated
with retraction — which is neither answer. The engine revokes a private account and hides it from its
owner (§7.5); the harness revokes the record and writes a retraction record. That is a case-2 departure
from T2 with no upside and no decision behind it. **If DN-4 is not answered, this is what continues.**

---

## 9. DN-6 and PH1A-12 after the decomposition

**TC-18 (Source-derived).** DN-6 (who may retract or revoke an account) is strictly downstream of DN-4:
it has two parts under "yes", one under "no", and no stable meaning before the scope question is
answered. It must not be asked first.

**TC-19 (Inferred).** PH1A-12 now has **exactly one** owner-decision predecessor — DN-4 — plus one
confirmation (DN-1):

| Predecessor | Kind | Effect on PH1A-12 |
|---|---|---|
| DN-1 | confirmation | if confirmed, no divergence; if declined, a v9 defect must be recorded first |
| DN-3 | settled by precedence | none — removes a predecessor the earlier framing had |
| DN-4 | **owner decision** | fixes whether the subject must express knowledge-loss |
| DN-5 (`seq`) | open T2 gap (PH1A-05) | the subject's schema cannot carry K-R15 either way |

**TC-20 (Superseded).** The post-Q-B review's DB-19 ("PH1A-03/04 is the single next decision") is
superseded by this decomposition. What remains as a decision is DN-4 alone.

---

## 10. What this report does not do

No architecture selected; no hybrid designed; no implementation, test, probe, model call or corpus. No
canonical edit, registry change or Git action. DN-1 is **not** answered on the owner's behalf, and DN-4
is **not** decided. Nothing here authorises implementation, testing or execution.

**Limits.** The textual finding (§2) rests on four passages read in full; it is a reading of ratified
text, not an amendment to it, and it remains subject to the owner's confirmation. The `Unestablished`
item in TC-11 is left open rather than resolved by assumption.

---

## 11. The two things now ready for the owner

1. **Confirmation (one sentence):** that v9 §8's "never a view, alias or projection" prohibits a view of
   the **source record**, not derivation from the recipient's own acquisition events — with the
   consequence that no PH1A-03 divergence exists, and the alternative would require recording an
   internal inconsistency in v9.
2. **Decision (one question):** can a character lose knowledge — forget something they once knew? If
   yes, K-R3/K-R19 apply as written; if no, knowledge-loss is an explicit scope exclusion and no
   operation may remove an account.

**End of report.**
