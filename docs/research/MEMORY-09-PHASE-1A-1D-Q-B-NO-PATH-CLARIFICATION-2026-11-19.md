# MEMORY-09 — Phase 1A-1D: Q-B "No" Path — Clarification of Five Statements

**Status:** research record. Read-only. **No implementation, no test, no probe, no model call, no
corpus, no canonical edit, no decision.** This pass corrects five statements in 1A-1C and does not
answer Q-B.

| Field | Value |
|---|---|
| Created | **2026-11-19** (local, CST) |
| Authority | None beyond the owner's review of 1A-1C. Creates no decision. |
| Scope | Five corrections to `…-1A-1C-Q-B-NO-PATH-PRESERVATION-VERIFICATION-2026-11-19.md` |
| Predecessors | 1A-1B brief; 1A-1C preservation verification |
| Tests executed | **0.** Static reading only. |
| Files changed | This file only. 1A-1C is **not** edited; it is superseded in the five places named below. |

**Why this record is short.** 1A-1C's seven checks stand. Five sentences in it claimed more than the
traces support, and one equivalence claim is withdrawn. Each correction is a narrowing, not a
reversal: no check changes from *preserved* to *broken*, and no check changes from *unbuilt* to *built*.

---

## 1. Correction register

| ID | 1A-1C statement now superseded | Replaced by | Basis |
|---|---|---|---|
| **C-1** | V-4: the retracted record is "invisible in **every** context, including to its owner and to any human view" | V-4 narrowed to the examined retrieval paths, with recipient accounts separately evaluated (§2) | Observed |
| **C-2** | §6 item 4: "The eligibility filter continues to exclude by status, in every context, structurally." | Character-usability excludes; inspection may expose and must never feed (§3) | Source-derived |
| **C-3** | §5: "P-5 and P-6 are the same boundary seen from two sides" | Adjacent, not equivalent; two independent invariants (§4) | Inferred + Source-derived |
| **C-4** | §3 P-7: disclosure is "the *only* route by which a character ever learns of a withdrawal" | One valid route among the acquisition routes D2 and K-R11 already name (§5) | Source-derived |
| **C-5** | §5 V-5: the harness retraction record is described as carrying "the withdrawal" without qualification | The divergence stands; what the record *means* is Unestablished (§6) | Observed + Unestablished |
| **C-6** | §7: "Only the second is consistent with the ratified set as it stands." | Only the second fits the T1 acquisition constraint; the specification is unreconciled while §8.4 stands (§7) | Source-derived |

**Nothing else in 1A-1C is withdrawn.** The four-surface partition (§1.1), RA-3's unproven status
(§1.2), the handling/awareness split (§1.3), and checks P-1, P-2, P-3, P-5, P-6 all stand.

---

## 2. C-1 — V-4 narrowed: three access paths, not one

The 1A-1C sentence treated "status is filtered" as though it implied "the record is gone". Three
different access paths exist and only one of them was examined.

| Path | Behaviour | Evidence |
|---|---|---|
| **Character retrieval** | revoked records excluded from the pool | `m09/retrieval.py:50`; `retrieve.py:126`, `:212`; `model.py:323–326` |
| **Object-level read** | revoked records returned, with status | `engine.py:1190–1195`, `:1198–1206`; `model.py:327`, `:333` |
| **Recipient-owned account** | a separate record, active, owned by the recipient; untouched by retraction of the source | `engine.py:542–559` (status change confined to the named record), `:1181` |

**V-4 — as corrected:**

> Revoked source records are excluded from the examined retrieval paths. This prevents them from being
> returned as current truth, but does not establish a human inspection surface. Object-level history and
> snapshots remain readable; recipient-owned accounts are separate records and must be evaluated
> independently.

**Consistency check (Observed).** The corrected wording removes a contradiction 1A-1C carried: P-1 says
a recipient's account survives source retraction, while the old V-4 said retracted material is invisible
"to its owner". Both cannot hold if "the retracted record" is read as any record about the retracted
claim. They hold together once the **source record** and the **recipient's account** are kept apart —
which is what O-1 and v9 rule 12 require in the first place.

**What this does not restore.** No inspection surface exists; P-5 remains *unbuilt*.

---

## 3. C-2 — The eligibility obligation, restated

1A-1C §6 item 4 was written as a universal rule and cannot be: an inspection context exists precisely
to surface what character-facing retrieval excludes. The obligation splits.

**As corrected:**

> Character-usable retrieval continues to exclude revoked source claims as current truth. The inspection
> surface may expose their status and history to the operator, but must never feed that information into
> character-usable context.

**Source basis (Source-derived).** r32 R2 gives `inspection` the widest read ("world state + knowledge
map + conflict map") and one prohibition ("must never feed the prompt"). The two clauses are a single
rule with two halves; stating only the exclusion half loses the reason the context is safe.

**Why the split is architectural, not verbal.** A single filter cannot serve both contexts: one must
exclude the record, the other must return it *with its status*. Treating them as one rule would either
hide history from the operator or expose it to the character. The separation is the mechanism.

---

## 4. C-3 — P-5 and P-6 are adjacent, not the same boundary

The 1A-1C claim is withdrawn. The two can fail independently, which is the test for non-equivalence:

- An inspection context can exist, be perfectly isolated, and character retrieval can still inject a
  status-derived notice at assembly (P-5 passes, P-6 fails).
- Character retrieval can be clean while no operator view exists at all (P-6 passes, P-5 fails).

That is today's actual state in both implementations: P-6 is preserved, P-5 is unbuilt. They were never
the same boundary; they only share a *nearby* boundary.

### 4.1 Two invariants, separately enforced

| | **I-A — Inspection isolation** | **I-B — Acquisition eligibility** |
|---|---|---|
| Protects | the operator's view from becoming a character's view | the character's context from unacquired knowledge |
| Mechanism | `inspection` is a **retrieval context with no assembly channel**: r32 §9 defines C1–C5 and no C6, so isolation is structural before it is guarded | nothing enters C2 without an `AcquisitionEvent`: K-R7, D2, K-R11 |
| Enforced at | the context boundary (K-R5, K-R6; R2's "must never feed the prompt") and assembly (G1/G2 plus the explicit prohibition) | formation/write (K-R1, K-R7) and assembly (C2's membership rule) |
| Fails as | operator data appearing in a prompt | a notice appearing in C2 with no event behind it |
| Depends on the other? | no | no |

**V-8 — as corrected (Source-derived).** The channel guard is one of **three** protections, not the sole
boundary: (i) `inspection` has no assembly channel to begin with; (ii) R2's explicit prohibition and
G1/G2 govern assembly; (iii) K-R7/D2 govern admission at write. Removing any one leaves the others, and
none of the three substitutes for P-6's own check.

### 4.2 What the disclosure responsibility adds

The owner's division is recorded as the frame and is consistent with v9 §1: the conversation layer
establishes that a disclosure occurred, its content and its recipients; MEMORY-09 records the resulting
accounts and never invents exposure. That division is **not** an invariant — it is an ownership
boundary, and it is what makes P-7's route legitimate while P-6's notice is not.

---

## 5. C-4 — "Only route" is too narrow; the routes are already enumerated

1A-1C P-7 said that under Q-B = no, disclosure is the only way a character ever learns of a withdrawal.
That narrows the architecture to one mechanism when the source names several.

**As corrected:**

> An actual disclosure is one valid route by which a character can later acquire knowledge of a
> withdrawal. MEMORY-09 must not manufacture that acquisition from status alone. Other experience,
> exposure, or authorised-transfer routes remain governed by the existing acquisition rules.

**The routes, from the sources (Source-derived).** D2's three categories — experience, exposure,
authorised transfer — and K-R11's event kinds are already the enumeration:

| Route | K-R11 kind | What it looks like for a withdrawal |
|---|---|---|
| Witnessing the withdrawal | `witnessed` | the character is present when the claim is struck or corrected |
| Being told | `told` | the retractor or a third party says so; creates an account under O-3 |
| Reading it | `read` | a document, notice or record the character actually reads |
| Inferring it | `inferred` | the character's model concludes it and the event is recorded as such |
| Verification | `verified` | the character checks and establishes it |
| Authorised transfer | `narrator_grant` | the declared, attributable arrangement |

**Caution on `inferred` (Inferred).** Inference is a **model act recorded as an event**, never a system
derivation from status. A system that computed "this character should now infer the retraction" would be
performing P-6's prohibited move under another name.

**What survives of P-7.** The mechanism is still unexercised and still matters — disclosure is the route
the ledger spends the most text on (O-3, v9 rules 46–49) and the one an owner-facing answer is most
likely to have in mind. It is simply not the only one.

**Consequence for the owner question.** "Someone tells them" is shorthand. The question must say
*acquire through a legitimate route*, or it decides more than Q-B.

---

## 6. C-5 — V-5: the divergence stands; the semantics do not

The trace establishes what the code does, not what the record is *for*. Three readings remain open, and
the code does not choose among them.

| Reading | What the record would be | Settled by the trace? |
|---|---|---|
| **S-a** The actor's declaration | the actor's own utterance that they withdraw the claim, recorded like any claim | **No** — `via="revoked"` and `kind="event"` are not a contribution declaration; the record is generated, not spoken |
| **S-b** The authoritative occurrence | application state that the retraction happened | **No** — the same object is also returned by `edges()` as an acquisition event (`model.py:347–353`), which is an application-state/event conflation, not a distinction |
| **S-c** The actor's knowledge of having acted | a memory record representing the actor's knowledge with `retracts=<id>` | **No** — ownership is copied from the target, and the target may not be the actor's in every call path |

**V-5 — as corrected (Observed).** The v9 harness generates an active `event` record
(`predicate=f"retraction:{target.predicate}"`, `text=f"{target.subject} retracted: {target.value}"`,
owning and scoping copied from the target, `model.py:282–299`, `:296–297`), and confinement to the
actor's pool is a property of `form._retract`'s own-record filter (`form.py:187–204`), not of the record
type. m09 has no retraction record at all. The divergence is real and must stay visible.

**Semantic interpretation: Unestablished.** The record must not be cited as proof that the conversation
layer established a disclosure in the form O-3 requires. It is evidence about representation, not about
event establishment.

---

## 7. C-6 — What "no" does and does not settle

1A-1C ended by saying only the "no" branch is consistent with the ratified set. That is true of the
acquisition constraint and false as a statement about the whole specification.

**As corrected:**

> A "no" answer fits the T1 knowledge-acquisition constraint (D2, O-5, A2b, E1c) without amendment. It
> does **not** reconcile r32 §8.4, which still instructs that the retraction be displayed in the
> character's channel. Resolving that — by classifying §8.4 as illustrative or by amending the rows — is
> an explicit owner-authorised act, not a consequence of answering "no".

**Why this matters (Source-derived).** 1A-1B RB-02 located the conflict *inside* r32: §8.4 against K-R7
and the K-R13 consequence, with no stated priority. A "no" answer chooses the K-R7 side of that conflict
by implication. An implication is not a classification; the owner must record it (1A-1C §6 item 5).

**So the branch is "consistent", not "reconciled".** Those are different claims and only the weaker one
is supported.

---

## 8. Revised verdicts

| # | Behaviour | Verdict after 1A-1D | Changed? |
|---|---|---|---|
| **P-1** | Acquired account survives source retraction | preserved (m09 at rest; v9 vacuous) | no |
| **P-2** | Belief not silently corrected | preserved (m09 by refusal; v9 vacuous) | no |
| **P-3** | Application recognises retracted status | preserved at the object layer | no |
| **P-4** | Retrieval does not present it as current truth | preserved **for the examined retrieval paths** | narrowed |
| **P-5** | Operator can inspect the retraction | **unbuilt** | no |
| **P-6** | No manufactured character awareness | preserved (m09 by construction; v9 by confinement) | no |
| **P-7** | Later learning through a legitimate acquisition route | mechanism unexercised; one route among several | widened |

**Unchanged conclusion.** No behaviour is lost *because* Q-B would be answered "no". The one unbuilt
behaviour (P-5) and the one unexercised route (P-7) are both in that state today, before any answer.

---

## 9. The owner question, in plain language

> **When a claim a character was told is later withdrawn, should that character's own context be
> allowed to say so — even though nobody told them and they did not witness it?**
>
> **Or should the character keep what they were told, and learn of the withdrawal only if they actually
> acquire it — by being told, seeing it, reading it, or another real route — while the application
> records the withdrawal and stops treating the claim as current?**

The second branch is not "the character never finds out". It is "the character finds out the way anyone
finds anything out". The first branch is the only one that needs a character to know something they
never encountered.

---

## 10. What this record does not do

No decision on Q-B, Q-A or Q-C. No edit to 1A-1C, the 1A-1B brief, r32, v9, the ledger or the
conformance map. No implementation change. No claim that any behaviour was executed. No statement that
either branch is correct — only that the two are now described without the five overstatements.

**End of record.**
