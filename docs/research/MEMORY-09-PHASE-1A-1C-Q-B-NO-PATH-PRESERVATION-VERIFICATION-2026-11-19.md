# MEMORY-09 — Phase 1A-1C: Q-B "No" Path — Preservation Verification

**Status:** research record. Read-only. **No implementation, no test, no probe, no model call, no
corpus, no canonical edit, no decision.** This pass does not answer Q-B, does not amend A2b/E1c/D2,
does not adopt RA-3, and does not change the 1A-1B brief.

| Field | Value |
|---|---|
| Created | **2026-11-19** (local, CST) |
| Authority | None beyond the owner's review direction on the 1A-1B brief. Creates no decision. |
| Scope | Whether the seven named behaviours survive a *hypothetical* Q-B = "no", verified against source and against both implementations |
| Predecessor | `MEMORY-09-PHASE-1A-1B-RETRACTION-AWARENESS-AND-RELATIONSHIP-DECISION-BRIEF-2026-11-19.md` |
| Tests executed | **0.** Every behavioural statement is a static reading at the hashes in Appendix B. |
| Files changed | This file only. Nothing staged, committed, registered or frozen. |

**Why this pass exists.** The owner's review accepted the 1A-1B findings and asked for one narrow
check before Q-B is put: that answering "no" does not silently discard a behaviour the ledger already
requires. This record performs that check. It is a preservation audit, not a decision, and it does not
make the architecture more complicated — it separates seven things that were previously compressed
into one.

---

## 1. Review disposition recorded

Three items from the owner's review are carried forward as **research dispositions**. None is a
decision, and none is a canonical edit.

### 1.1 RB-05 accepts a scope correction to its wording

The owner's qualification is right and is narrower than a rebuttal: D2 bars **one** surface, not four.
The 1A-1B brief sometimes wrote as though the internal join and the character's acquisition were one
operation. They are four distinct information positions, and only the last is governed by D2.

| Surface | Holds the retraction? | Who can see it | D2 engaged? |
|---|---|---|---|
| **S-i** Application state | yes — status projection, revocation reason | system internals | **No** |
| **S-ii** Human inspection | yes — `inspection` context (r32 R2) | the operator | **No** |
| **S-iii** Retrieval eligibility | yes — used as a filter, never surfaced as text | the eligibility filter | **No** (it excludes; it does not inform) |
| **S-iv** Character-usable context | yes — text in C2 | the character's model | **Yes**, and only here |

**V-0 (Source-derived).** RB-05 stands for S-iv and only for S-iv: a character-usable notice with no
acquisition event is knowledge without experience, exposure or authorised transfer. S-i, S-ii and
S-iii are not character knowledge and D2 does not reach them. This is a correction to the brief's
phrasing (§9.1, §9.4), not to its finding.

**V-0b (Inferred).** The distinction is load-bearing for Q-B's wording: a "no" answer must be read as
"no to S-iv", never as "no to S-i–S-iii". The owner's formulation — MEMORY-09 knows, the character
does not — is exactly this partition.

### 1.2 RA-3 is a design candidate, not an established architecture

The 1A-1B matrix shows RA-3 is the **only** relation of the three that can carry a shared content
identity without violating O-1. That is a statement about the three relations described, not a proof
about RA-3.

Three properties remain unproven, and each is a design obligation rather than a reading:

| Property | What would establish it | Current status |
|---|---|---|
| **Safe** | the identity cannot be used to enumerate holders or reconstruct an audience (F1) | **Unestablished** — no mechanism defined |
| **Necessary** | no cheaper mechanism supplies sameness for the retraction join | **Unestablished** — not compared against alternatives |
| **Sufficient** | it can compute the notice for a recipient *and* for two independent witnesses | **Unestablished** — no join specified |

**V-0c (Observed).** The only textual obstacle on the v9 side is the `proposition_id` refusal on
private records (`records.py:14–20`, `:155–158`), whose stated reason is enumeration traversal. RA-3
would have to be new architecture under O-6a case 3 and declared as such. Nothing here adopts it.

### 1.3 Two responsibilities, named apart

The owner's distinction is accepted as the organising frame for this pass, because without it a "no"
to Q-B reads as a refusal to handle retractions at all.

| | **Retraction handling** | **Retraction awareness** |
|---|---|---|
| Question | can the system recognise that a claim is withdrawn and stop presenting it as current truth? | is the character informed that something they hold was withdrawn? |
| Owner | MEMORY-09 (status, eligibility, provenance) | the conversation layer (disclosure), never MEMORY-09 |
| Depends on Q-B? | **No** — required either way | **Yes** — this is Q-B |
| Where it lands | S-i, S-ii, S-iii | S-iv |
| Failure mode if absent | stale claim presented as current | character's model silently given knowledge it never acquired |

**V-0d (Source-derived).** Handling is entailed by E1a, I4, I7, K-R12 and the eligibility stage of the
r32 stage contract (§4: eligibility excludes by status). Awareness is entailed by nothing in the
ratified set (§2.3 of 1A-1B). They are separately required, and only one of them is contested.

---

## 2. The hypothetical under test, stated precisely

> **Q-B = no.** A character's usable context never reports a retraction they were not told about.
> The character keeps what they acquired, keeps their beliefs, and learns of the withdrawal only if
> someone tells them. MEMORY-09 records the withdrawal and stops treating the claim as current.

Two things this hypothetical does **not** say, both of which the check below verifies rather than
assumes:

1. It does not say the retraction is forgotten by the system.
2. It does not say the character is permanently unable to learn.

**Test method.** For each behaviour: (a) name the source requirement that carries it; (b) read both
implementations for a surface that would satisfy or break it; (c) give a verdict in one of four
values — **preserved**, **preserved at rest only**, **unbuilt**, **vacuous** — and (d) name any residual
risk. "Preserved" means no named text is violated and a surface exists. It is not a test result.

---

## 3. The seven preservation checks

### P-1 — The original acquired account survives source retraction

**Required by:** O-3 (verbatim: "subsequent retraction of the source does not erase the recipient's
acquired account"); H10; A2a; v9 rule 49.

**Observed, m09.** `retract` changes the status of the **named record only** (`engine.py:542–559`;
`rec.status = REVOKED` is the whole effect). Accounts are separate records with their own
`record_id` and owner, and `accounts_of` returns them without consulting the source
(`engine.py:1181`). No cascade exists from a retracted source to a recipient account.

**Observed, v9 harness.** There are no recipient accounts to preserve: `_admit` writes every formed
record against the speaker with `via="witnessed"` (`form.py:206–240`), and `disclose` is never called
(`m09/engine.py:575` has no call site in `eval/mem1`). The property is therefore **not exercised**
rather than contradicted.

**Verdict: preserved** (source and m09 at rest) · **vacuous in the v9 harness**.

**Residual risk.** Survival *at rest* is not survival *at the character's surface*. If the account
itself is ever revoked, `perspective` drops it for its own owner (`m09/retrieval.py:50`), which is the
1A-1 RP-02/GB-8 defect. P-1 does not depend on that path — a source retraction is not an account
retraction — but the two must not be conflated when the notice question is implemented.

### P-2 — The character's belief is not silently corrected

**Required by:** O-5 (verbatim: "A character's mind is not a database field for the system to edit");
v9 rules 15, 38, 53–55; S34; E1a.

**Observed, m09.** `retract` refuses a belief outright:

> `BeliefError("a belief is not retractable here; belief lifecycle is model-gated (D-B) -- use revoke_belief")`

inside `engine.py:542–559`. There is no path from a source retraction to a belief record, and belief
revision has its own model-gated operation (`declare_belief` / `revoke_belief`,
`engine.py:311–361`).

**Observed, v9 harness.** Belief is not representable: `KINDS` is
`("fact", "preference", "commitment", "event", "state_change", "absence")` (`model.py:48`). The
harness cannot correct a belief because it holds none. Not a preservation success — an absent
representation.

**Verdict: preserved** (m09, by refusal) · **vacuous in the v9 harness**.

**Residual risk.** "Not corrected" is satisfied trivially while awareness is absent, because there is
nothing to correct *with*. Under Q-B = yes the same guarantee has to survive a notice in C2, and
nothing in either implementation would enforce that separation beyond the C2/C3 channel boundary.

### P-3 — The application can still recognise the retracted status

**Required by:** E1a ("Only the memory's truth status changes"); I4 (history never destroyed);
K-R12 (status is a projection of its own stream); r32 stage contract §4 (eligibility excludes by
status).

**Observed, m09.** Status is a first-class field with a closed set
(`records.py:69–72`: `ACTIVE | SUPERSEDED | REVOKED`), validated in `__post_init__` (`:113`). The
record is retained; `engine.history()` returns it regardless of status (`engine.py:1190–1195`), and
`snapshot()` serialises it with its status (`engine.py:1198–1206`).

**Observed, v9 harness.** `revoke` sets `status` and records `revoked_at_turn` plus a reason map
(`model.py:303–311`); `history()` and `versions()` return revoked records (`:327`, `:333`);
`active_value()` returns `None` when no active version exists (`:338–343`).

**Verdict: preserved** — in both, at the object layer.

**Residual risk.** Recognising is not exposing. See P-5: no retrieval context in either implementation
can be asked for it.

### P-4 — Retrieval does not treat the retracted claim as current established truth

**Required by:** I1 (eligibility precedes ranking), I6 (empty set → abstain), I7 (retrievability does
not confer authority), r32 §4 eligibility stage, v9 rule 38.

**Observed, m09.** Two independent gates: the character pool excludes `status == "revoked"`
(`m09/retrieval.py:50`), and the authoritative views filter `r.status == ACTIVE`
(`engine.py:175`, `:328`, `:357`). A retracted world fact therefore abstains rather than returning
the stale value.

**Observed, v9 harness.** `MemoryStore.current()` filters `status == ACTIVE` (`model.py:323–326`), and
`eligible` re-checks it independently (`retrieve.py:126`), so a revoked record cannot enter a
candidate pool in either declared context. `retrieve()` repeats the check at `:212`.

**Verdict: preserved** — and preserved *structurally*, by exclusion rather than by annotation.

**Residual risk (V-4).** Exclusion is not neutral. Because the check is a hard filter in both
implementations, the retracted record is invisible in **every** context, including to its owner and to
any human view — the read-time join of K-R13 has no implementation anywhere (RP-13). A "no" answer
does not create this gap; it inherits it, and P-5 is where it bites.

### P-5 — The human operator can inspect the retraction and its effects

**Required by:** r32 R2 (`inspection`: "world state + knowledge map + conflict map", "must never feed
the prompt"); §9 guard structure (the narrator may hold what the character does not, provided it is not
injected into the character channel).

**Observed, m09.** There is no inspection context. `CONTEXTS = ("perspective", "world",
"belief_facts")` (`m09/retrieval.py:29`), and `dispatch` returns `unknown_context` for anything else
(`:98`). Object-level accessors exist — `history()` (`engine.py:1190`) and `snapshot()`
(`:1198`) — but no view consumes them and no context is named for them.

**Observed, v9 harness.** `CONTEXTS = ("character_perspective", "world_state")`
(`retrieve.py:24`). The frozen four-context table is not implemented; the existing G1 review receipts
record exactly this as blocked on a contract decision, not on code
(`g1-review-receipt-r2.json:19`, `:73`).

**Verdict: unbuilt.** The requirement is clear and unimplemented in both. This is a build gap, not a
semantic loss: a "no" answer neither creates nor worsens it, but the behaviour the owner listed is
**not** preserved today in any surface an operator can use.

**Residual risk.** Without an `inspection` context, "MEMORY-09 knows the retraction" is true only of
the object layer. A reviewer cannot observe the claim's withdrawal except by reading code. If the
owner's seven behaviours are treated as acceptance criteria for the "no" path, this is the one that
fails today.

### P-6 — The system does not manufacture retraction awareness for characters

**Required by:** D2; K-R7 ("requires an `AcquisitionEvent`, never a context switch alone"); K-R13's
consequence; v9 rule 10; A2b; E1c.

**Observed, m09.** No code path writes a status, marker or notice into a character's usable context.
`perspective` *removes* non-active records and adds nothing (`m09/retrieval.py:46–62`); `disclose`
creates accounts only from an explicit event with supplied recipients (`engine.py:575–700`). m09 has
no retraction-fact object at all, so there is nothing that could be injected.

**Observed, v9 harness — the one asymmetry worth recording.** `model.retract` emits a **new active
record** carrying the withdrawal in its text: `kind="event"`,
`predicate=f"retraction:{target.predicate}"`, `text=f"{target.subject} retracted: {target.value}"`,
owned by the target's owner and copying its scope and audience (`model.py:282–299`,
`:296–297`). It is consumed as answer evidence (`score.py:86–89`) and appears in fields as an
expected record (`fixtures.py:242`).

Because `form._retract` only reaches records owned by the span's own speaker
(`form.py:187–204`, `current(owner_id=owner_id)`), the record lands only in the retractor's own pool.
A character therefore receives a retraction statement about **their own act**, not about someone
else's. That is not manufactured awareness of another's retraction — but it is a surface that would
become one the moment scope or audience were generalised.

**Verdict: preserved** — by construction in m09, by confinement in the v9 harness.

**Residual risk (V-5).** The v9 harness has a retraction record and m09 does not. Any notice design
must decide which of the two it extends; extending the v9 shape without its ownership confinement
would produce exactly the eventless notice D2 forbids.

### P-7 — A character can learn of the retraction later through an actual disclosure

**Required by:** O-3 (disclosure needs no authorisation and creates recipient accounts); v9 rules
46–49; H10; J3/J7; rule 47 (narration may report, never manufacture).

**Observed, m09.** `disclose` is implemented and creates one account per supplied recipient with a
fresh descriptor, no source id and no status check on the source (`engine.py:575–700`). It is
**never called** in `eval/mem1`: there is no call site, and `form._admit` cannot create a recipient
account (`form.py:206–240`).

**Observed, v9 harness.** No disclosure operation exists. `form` writes records for the speaker only,
always `via="witnessed"`; narration is excluded from retraction (`form.py:187–204`) and nothing
creates a second party's account.

**Verdict: mechanism present, path unexercised** (m09) · **absent** (v9 harness).

**Residual risk (V-6).** Under Q-B = no this is the *only* route by which a character ever learns of a
withdrawal, so its absence is load-bearing rather than incidental. Note also what is **not** required:
no retraction-fact object is needed for the retractor to say "I withdrew that". The content is an
ordinary claim; the conversation layer establishes the event; accounts follow. m09's lack of a
retraction record (P-6) does not block P-7, and the v9 harness's retraction record does not by itself
provide it.

---

## 4. Consolidated result

| # | Behaviour | Source basis | m09 | v9 harness | Verdict |
|---|---|---|---|---|---|
| **P-1** | Acquired account survives source retraction | O-3, H10, A2a, v9 r49 | preserved at rest | vacuous | **preserved** |
| **P-2** | Belief not silently corrected | O-5, v9 r15/r38/r53–55, S34 | preserved by refusal | vacuous | **preserved** |
| **P-3** | Application recognises retracted status | E1a, I4, K-R12 | preserved | preserved | **preserved** |
| **P-4** | Retrieval does not present it as current truth | I1, I6, I7, §4 | preserved | preserved | **preserved** |
| **P-5** | Human operator can inspect the retraction | r32 R2 `inspection` | absent | absent | **unbuilt** |
| **P-6** | No manufactured character awareness | D2, K-R7, K-R13c, v9 r10 | preserved | preserved (confined) | **preserved** |
| **P-7** | Later learning through disclosure | O-3, v9 r46–49, H10 | implemented, uncalled | absent | **path unexercised** |

**What the table shows.** Five of the seven are preserved at source level, one is preserved only as a
mechanism with no exercised path, and one — the operator's own inspection — has no surface in either
implementation. **No behaviour is lost *because* Q-B would be answered "no".** The two weak rows are
weak today, before any answer, and would be equally weak under "yes".

**What it does not show.** Preserved is not executed. Every "preserved" verdict is a static reading;
0 of the relevant rows (A2a, A2b, E1a–E1c, H10) are implemented as probes, and the corpus does not
exist (1A-1: 0/144 sessions, G07/G12 unbound).

---

## 5. Residual risks, stated once each

| ID | Risk | Consequences if unaddressed | Label |
|---|---|---|---|
| **V-4** | Status is enforced by hard exclusion in every context, so no surface can report a retracted record — not the owner's, not the operator's | P-4's guarantee is bought at P-5's cost; the "handling" half exists only as a filter | Observed |
| **V-5** | The v9 harness represents a retraction as an active record in the actor's own pool; m09 represents it as a status with no record | Two incompatible shapes for the same fact; generalising the v9 shape would create an eventless notice | Observed |
| **V-6** | Disclosure of a withdrawal has no exercised path | Under "no", later learning is the only route and it is untested | Observed |
| **V-7** | A2b's precondition ("Same character, same turn") is still unread | The "no" answer does not need it; the "yes" answer does, and it cannot be supplied by interpretation | Unestablished (1A-1B GB-2) |
| **V-8** | The `inspection` context must never feed the prompt, and neither implementation has a guard because neither has the context | When it is built, the guard is the only thing separating S-ii from S-iv | Source-derived |

**V-8 is the one that matters structurally.** P-5 and P-6 are the same boundary seen from two sides:
the moment an operator view exists, the only thing preventing it from becoming a character notice is
the channel guard in r32 §9 (`G1`–`G4`) and the absence of any code path from `inspection` to
assembly. A "no" answer is only as strong as that guard.

---

## 6. What a "no" answer would commit to — obligations, not decisions

If the owner answers Q-B = no, the following are consequences of that answer rather than new choices.
None is adopted here.

1. S-iv carries no retraction text, marker or annotation derived from status.
2. MEMORY-09 keeps status, reason and history, and keeps presenting nothing about it to a character.
3. The `inspection` context is required to exist, and is required never to feed assembly.
4. The eligibility filter continues to exclude by status, in every context, structurally.
5. §8.4's display sentence is recorded as illustrative and non-binding — **this is an amendment-shaped
   act** and would be the owner's, not this record's.
6. A2b, E1c and D2 need no amendment.
7. Q-A's sameness half is not forced; RA-2 remains conformant.
8. Q-C is moot.
9. The disclosure path must be built and exercised, because it becomes the sole route to awareness.

Item 5 is the only one that touches canonical text, and it is exactly why "no" is not free: it settles
r32's internal inconsistency by demoting one of r32's own worked examples.

---

## 7. What this record does not do

No decision on Q-B, Q-A or Q-C. No amendment to the brief, the ledger, r32, v9 or the conformance map.
No implementation change. No claim that any behaviour was executed, demonstrated or reproduced. No
statement that the "no" path is the correct answer — only that, if it were chosen, six of the seven
behaviours survive on paper and the seventh was never built.

**Answer-readiness.** The question put to the owner can now be stated without architecture in it:

> When a claim a character was told is later withdrawn, should that character's own context be allowed
> to mention the withdrawal, even though nobody told them? Or should the character keep what they were
> told and learn of the withdrawal only if someone says so — while the application records it and stops
> treating the claim as current?

Both branches are fully consistent with the system knowing. Only the second is consistent with the
ratified set as it stands.

---

## Appendix A — Implementation traces verified this session

| Trace | Location | What it shows |
|---|---|---|
| **RBT-11** | `m09/engine.py:1181` `accounts_of` | accounts returned without consulting their source's status |
| **RBT-12** | `m09/engine.py:1190–1195`, `:1198–1206` `history`, `snapshot` | status and full record set readable at the object layer |
| **RBT-13** | `m09/retrieval.py:29`, `:98` | three contexts; unknown context abstains — no `inspection` |
| **RBT-14** | `retrieve.py:24`, `:126`, `:212` | two contexts; status filtered before scoring and again in `retrieve` |
| **RBT-15** | `model.py:282–299`, `:296–297` | retraction emits an active `event` record carrying the withdrawal in text |
| **RBT-16** | `model.py:303–311` | revocation keeps the record, sets status, records turn and reason |
| **RBT-17** | `model.py:48` | `KINDS` carries no belief — belief is not representable in the v9 harness |
| **RBT-18** | `m09/engine.py:311–361` | belief created and revoked only by the character's own model |
| **RBT-19** | `form.py:187–204` | retraction confined to the speaking character's own current records |
| **RBT-20** | `form.py:206–240` | every formed record is the speaker's, `via="witnessed"`; no recipient account |
| **RBT-21** | `m09/engine.py:575–700` | `disclose` implemented; no call site anywhere in `eval/mem1` |
| **RBT-22** | `g1-review-receipt-r2.json:19`, `:73` | four-context contract recorded as blocked on a decision, not on code |

## Appendix B — Sources read this session

Static reading only; **no file was executed and no hash below is a result**.

| Source | Lines | Tier |
|---|---|---|
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` (r32) | 1028 | T2 |
| `MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` | 236 | T3 |
| `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` | 636 | T1 |
| `MEMORY-09-CONFORMANCE-MAP-2026-11-17.md` | 211 | T1 (ratified) |
| `eval/mem1/m09/engine.py` | 1209 | Observed |
| `eval/mem1/m09/records.py` | 252 | Observed |
| `eval/mem1/m09/retrieval.py` | 101 | Observed |
| `eval/mem1/model.py` | 398 | Observed |
| `eval/mem1/form.py` | 242 | Observed |
| `eval/mem1/retrieve.py` | 235 | Observed |
| `eval/mem1/conversation.py` | 414 | Observed |

## Appendix C — Claim labels used

`Observed` — read directly from the file and line range cited.
`Source-derived` — follows from a quoted requirement or rule, with the derivation stated.
`Inferred` — a judgement this record makes and offers for challenge.
`Unestablished` — not answered by any source or code read here.
`Superseded` — **0 labels applied**; the only appearance of the word is this sentence. Nothing in this
record replaces a prior finding. RB-05 is narrowed in wording (§1.1), not withdrawn.

**End of record.** No canonical artifact was modified. Q-B remains the owner's to answer.
