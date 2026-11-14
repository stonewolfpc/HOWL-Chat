# Muse — Revision 7 delta review

**Date:** 2026-11-11
**Reviewer:** Muse
**Review type:** Focused documentation-delta review (read-only)
**Subject:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`, revision 7
**Request:** `docs/review/AMENDMENT-02-REVISION-7-MUSE-DELTA-REVIEW-REQUEST-2026-11-11.md`
**Status:** Complete; no document or source file modified

## 1. Verdict

**PASS WITH REQUIRED CHANGES**

Revision 7 closes `M3-01` through `M3-05` substantively. This is not a case of adding six criteria and declaring victory: the underlying sections were changed, and the criteria were written to test the changed text.

- `M3-02` is the clearest improvement. §7.1 no longer describes an unconditional transaction. It opens with "Selecting a new scenario is a proposal until the confirmation gate resolves," inserts a gate as step 2–3, and names exactly what decline leaves untouched — active world/scenario, participant state, visible history, memory, and model context — while permitting the rejected attempt to survive "only as a non-active proposal." `AC-C68` mirrors that list rather than paraphrasing it. `AC-C23` was also rewritten from "immediately records and displays" to "After any required `OD-03` confirmation," which removes the wording that made the ordering ambiguous.
- `M3-03` is handled with the right mechanism, not just a rule. §9.3 now states that Chat-derived entries are conversation-scoped by default, enter character-persistent memory only through explicit promotion, record source scope, prior scope, and promotion reason, and — the load-bearing clause — "Retrieval must enforce scope before merging the entry into any character context." Placing enforcement at retrieval rather than at write time is correct, because a write-time-only rule is defeated by any index that was built earlier.
- `M3-04` is closed at the rule level by §7.6 and at the test level by `AC-C71`, including the originating-chat history half.
- `M3-05` is genuinely closed: §7.2 expanded from "the change reason and actor" to actor, reason, affected scope, previous value, resulting value, and confirmation outcome, and `AC-C72` requires that record to survive reload and remain retrievable.
- `M3-06` is correctly declined rather than "fixed," which preserves ID stability.

Three required changes remain. They are small, documentation-only, and all sit on the new gate rather than on the M3 findings: the gate has no durability rule, promotion has no atomicity or reversal rule, and `AC-C71` is narrower than the equivalent privacy criterion. Details in §5.

This verdict is a documentation-delta result only. It is not architecture approval, not acceptance, and not research or implementation authorization.

## 2. Identity verification

| Artifact | Expected | Recomputed | Result |
|---|---|---|---|
| Amendment 02, revision 7 | `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1` | same | Match |
| Size / lines | `76,565 B / 1,241 lines` | same | Match |
| Sidecar | `sha256sum -c` | `OK` | Verified |
| Revision 6 archive | `ef2d5dd88fc30cf49424f6bcf20f50aca5c599ae6d9a8f76a33bc65673370433` | same | Match |
| Revision 6 archive size | `71,428 B` | same | Match |
| Integration plan | `1e92a54434e8e3269015c47374ad57bd6ae755e64e27a07954c083ab952510cf` | same | Unchanged |
| Amendment 01 | `f3a0f8b9c5d19129db85d559f1880f21451cc7aa70596d37a408952c85c3f4ea` | same | Unchanged |
| Revision-3 review | `0c82874b8fc5c0b82aa3a0e42b83802a8e21b2a2384b453758a4c5a429b1acf3` | same | Unchanged |
| Revision-6 OD-04 addendum | `b7be1d46abd2c82fd2aa6edda824b9ddf4c6f1f3bf90df500806d342b91f625d` | same | Unchanged |

The revision-6 archive exists as a separate file at exactly the recorded digest and byte size. It was not rewritten in place, and the current amendment was not edited and re-hashed under the same identity.

## 3. Disposition of the M3 findings

I checked each claim against the text rather than accepting the applied status.

| Finding | Claimed | My assessment |
|---|---|---|
| `M3-01` | Applied | **Confirmed.** §16.8 adds `AC-C67`–`AC-C72` and labels them "unvalidated design targets." They cover `P1`–`P6` from the prior addendum. |
| `M3-02` | Applied | **Confirmed.** §7.1 gate, §7.1 decline semantics, `AC-C23` rewrite, `AC-C68`. The gate now precedes active state by construction. |
| `M3-03` | Applied | **Confirmed.** §9.3 principles plus `AC-C70`, including retrieval-time enforcement. |
| `M3-04` | Applied | **Confirmed, with a coverage narrowness** — see `M4-03`. |
| `M3-05` | Applied | **Confirmed.** §7.2 field expansion plus `AC-C72` survival and retrievability. |
| `M3-06` | Informational | **Correctly declined.** `AC-C66` remains a stable ID; no renumbering. |

## 4. Detailed probe — the three requested areas

### 4.1 Can a failed or interrupted confirmation leave partially applied state?

**Partly. The gate is specified for decline and for a failed apply, but not for interruption of the gate itself, and the commit boundary is ambiguous.**

§7.1's sequencing state is now correct: nothing active before confirmation, and decline is defined as leaving active world/scenario, participants, history, memory, and model context unchanged with no canonical transition or first message.

Interruption exposes two gaps.

**The commit boundary.** Step 4 says "On approval or when no confirmation is required, record the selection." Step 5 persists scenario identity and revision. Step 9 records the transition event. Step 10 requires the transition and first message to be persisted together "or report a recoverable partial failure." So the sequence contains at least two durable writes with a multi-step interval between them — the selection record at step 4 and the atomic pair at step 10 — and nothing states which one marks the transition as committed or which record the UI and model read as authoritative.

The text does establish that the system can distinguish non-active state, since step 3 permits a rejected attempt to be recorded "only as a non-active proposal." That distinction exists; it is simply never assigned to the step 4 record. A reader implementing this could reasonably make step 4's write canonical, in which case a crash before step 10 leaves a recorded selection with no transition and no first message — partially applied state that the text's own ordering permits.

**Restart during a pending confirmation.** The word "restart" appears in the amendment only for images (`AC-C01`), model selection (`DA-I10`), and general durability (§3). Nothing addresses a process that stops while the confirmation dialog is open. `AC-C68` forbids any change becoming active before required confirmation, so an implementation that resumed and auto-applied would violate a stated rule — but no criterion tests the restart path, which is the same category of gap `M3-01` identified: a rule with no test.

**Required (`M4-01`):** name the commit point explicitly — one durable applied marker written atomically with the transition and first message, with any record lacking that marker treated as a non-active proposal — and add one criterion that a restart during a pending or partially applied confirmation leaves either the complete pre-transition state or the complete post-transition state, never a hybrid, and does not re-apply without confirmation.

### 4.2 Is memory promotion atomic, auditable, and reversible?

**Auditable: yes. Atomic: not stated. Reversible: entirely unaddressed.**

**Auditable** is satisfied. §9.3 requires source scope, prior scope, and promotion reason, and `AC-C70` requires them when promoted. Recording prior scope is the right choice because it is the minimum data needed for any later decision about reversal.

**Atomic** is not addressed. Promotion writes at least three things — the entry's scope, the promotion record, and whatever retrieval view or index depends on scope — yet there is no coordinated-operation requirement and no recoverable-failure clause. This is an asymmetry worth noting: §7.1 step 10 imposes exactly that discipline on scenario transitions, and §9.5 imposes receipts on deletion, but promotion between the two receives neither.

**Reversible** is absent. Searching the amendment for demotion, revocation, or reversal produces no matches. `R-11` lists "Audit receipts and rollback" for deletion, and `DA-C12` covers tombstone and retention, but nothing governs a promotion being undone. Because `OD-04` is a least-privilege rule — the narrow scope is the safe default — the safe failure mode is being able to retract a promotion, and the amendment currently leaves that question to be answered silently during implementation.

The substantive problem is that reversal is not symmetric. Once an entry is character-persistent, it may already have been merged into other conversations' retrieval. Reversing the scope cannot unring that, and an implementation that tries will either fail or quietly do something destructive. The amendment should decide, not let this be decided in code.

**Required (`M4-02`):** add an atomicity clause to §9.3 mirroring §7.1 step 10, and require an explicit decision under `R-09`/`DA-C4` on whether demotion is supported. If supported, state that demotion is itself a recorded promotion-class change, that it stops future retrieval without erasing what other conversations already consumed, and add the matching criterion. If not supported, say so explicitly so the answer is not invented later.

### 4.3 Do the negative cross-chat tests cover direct records and derived memories, summaries, and retrieval?

**Direct records: covered. Derived memory: covered by a different criterion. Summaries, retrieval results, receipts, and new conversations: not covered.**

`AC-C71` tests that "Death, Step Out, return, or ghost participation in one conversation leaves the character's state in every other conversation unchanged," plus originating-chat history. That is the right negative shape — it tests the unaffected chat, not just the affected one, which is what makes it valuable.

Its weakness is enumeration. Compare `AC-C49`, which names eleven leakage channels for private data: memory, automatic lore, summaries, RAG, attachments, descriptions, transitions, model-input receipts, persistence, context assembly. `AC-C71` compresses everything into "the character's state." Four specific holes follow:

1. **Summaries.** A summary generated in another conversation must not acquire a lifecycle fact from the originating one. `AC-C53` covers private-only facts in public summaries; nothing covers cross-chat lifecycle facts.
2. **Retrieval results.** `AC-C55` constrains RAG by audience, not by conversation scope. A memory derived from a death event is neither a private-audience object nor a correction, so it falls between `AC-C70` and `AC-C71`.
3. **The gap between the two criteria.** `AC-C70` covers "a correction or derived memory originating from Chat." A lifecycle event is not a correction. `AC-C71` covers the state itself but not a memory derived from it. So "I died in that cave" recorded as a character memory is covered by neither.
4. **Newly created conversations.** `AC-C71` tests "every other conversation" — the ones that exist. A global status flag written to the character record but read only at conversation creation would leave every existing conversation unchanged and pass. That is the most likely way this requirement gets violated in practice, and the current criterion cannot catch it.

**Required (`M4-03`):** extend `AC-C71` or add `AC-C73` to enumerate the derived surfaces — memory, automatic lore, summaries, retrieval results, and model-input receipts — and to include a conversation created after the lifecycle event.

## 5. Required changes

| ID | Change | Scope |
|---|---|---|
| `M4-01` | Define the commit boundary in §7.1; add a restart/interruption criterion. | One clause, one criterion |
| `M4-02` | Add promotion atomicity to §9.3; require an explicit demotion decision; add criteria. | One clause, one register line, one or two criteria |
| `M4-03` | Enumerate derived surfaces and the new-conversation case in `AC-C71`. | One criterion |

All three are documentation changes and none reopens `M3-01`–`M3-05`.

## 6. Consistency checks requested

| Check | Result |
|---|---|
| `AC-C23` no longer implies application before confirmation | **Pass.** Rewritten to "After any required `OD-03` confirmation … A declined selection does not create a transition." |
| `R-09` prevented from weakening `OD-02`–`OD-04` | **Pass, and strong.** §12.2 states `R-09` "may specify the field taxonomy, storage, and promotion mechanics, but it may not relax `OD-02`, `OD-03`, or `OD-04`," and that any proposal broadening scope or weakening confirmation "must be reported as a conflict requiring a new owner decision, not resolved by research or implementation." §15 repeats the non-narrowing rule for `O-11`. |
| Revision 6 preserved at the recorded digest | **Pass.** Separate file, `ef2d5dd8…`, 71,428 B. |
| New criteria are design targets, not claims of behavior | **Pass.** §16 preamble and §16.8 both label them unvalidated design targets. |
| No research, runtime, migration, or implementation authorization implied | **Pass.** §20's prohibitions are unchanged; the delta adds no permission language; `DA-C4` still blocks implementation. |
| Baselines and prior review artifacts unchanged | **Pass.** All four digests match. |

## 7. Note on what this revision got right

Worth recording, since it is the reason only three narrow changes are required: the revision fixed causes rather than symptoms. `AC-C23`'s "immediately" was not left in place while a new criterion compensated for it. §9.3 put enforcement at retrieval, which is where scope rules actually fail. §7.2 expanded the record fields instead of adding a separate audit log that could drift from the transition. And `M3-06` was declined rather than "resolved," which is the discipline that keeps the ID space trustworthy.

## 8. Authorization boundary

This review authorizes nothing. Per the review request's stop condition: no `R-07`–`R-11` research, implementation, migration, runtime investigation, model loading, benchmarking, network testing, Android pairing, or production change follows from this document. Architecture approval and acceptance remain not passed. The package is a corrected requirements document awaiting a controlled revision, not a validated implementation design.
