# Muse — Revision 9 delta review

**Date:** 2026-11-11
**Reviewer:** Muse
**Review type:** Focused documentation-delta review (read-only)
**Subject:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`, revision 9
**Request:** `docs/review/AMENDMENT-02-REVISION-9-MUSE-DELTA-REVIEW-REQUEST-2026-11-11.md`
**Status:** Complete; no document or source file modified

## 1. Verdict

**PASS WITH REQUIRED CHANGES**

Revision 9 closes `M5-01` through `M5-04`. The commit state is now in the canonical model with an owner; promotion is an authorization commit with materialization explicitly outside the transaction; revocation is four named operations; and `AC-C71` is now an observed non-receipt scenario rather than a state comparison.

Three required changes remain. All three are narrow, and two of them take the same form: the revision states a *universal property* that must hold across every participant or every retrieval path, but states it as an obligation rather than as a testable outcome. That is the pattern that produced the earlier "rule with no test" findings, and it is worth closing at the level of a criterion rather than by prescribing a mechanism the research has not chosen yet.

I want to be clear about what I am **not** asking for. I am not asking revision 9 to select a coordinator, a storage engine, or an authorization architecture. Those belong to `R-07`/`DA-C1` and `R-10`/`DA-C7`. What the amendment should do is state the property that any such mechanism must satisfy, in a form an acceptance test can falsify. Details in §4.

This is the point at which the revision chain has done its job. One more pass to apply these three, and then the next document should be a research dossier, not a revision 10. §7 explains why.

This verdict is a documentation-delta result only. It is not architecture approval, acceptance, research authorization, or permission to begin any workstream.

## 2. Identity verification

| Artifact | Expected | Recomputed | Result |
|---|---|---|---|
| Amendment 02, revision 9 | `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640` | same | Match |
| Size / lines | `88,499 B / 1,303 lines` | same | Match |
| Sidecar | `sha256sum -c` | `OK` | Verified |
| Revision 8 archive | `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f` | same | Match |
| Revision 8 archive size | `82,682 B` | same | Match |
| Revision 8 archive sidecar | `sha256sum -c` | `OK` | Verified |
| Revision 7 archive | `4cd76064…` | same | Unchanged |
| Revision 6 archive | `ef2d5dd8…` | same | Unchanged |
| Revision-8 delta review | `0c67713a62059db099add6b8891aa50e01ba99ca73298293b6941368b551e291` | same | Unchanged |
| Revision-7 delta review | `8e7f248b…` | same | Unchanged |
| Integration plan | `1e92a544…` | same | Unchanged |
| Amendment 01 | `f3a0f8b9…` | same | Unchanged |

## 3. Disposition of the M5 findings

| Finding | Claimed | My assessment |
|---|---|---|
| `M5-01` — commit state in the model | Applied | **Confirmed.** `commitState` and `appliedMarkerRef` in the §4 shape; invariant 8 forbids projection of uncommitted records; invariant 9 states the store/ two-phase rule; §4.2 assigns the marker and pending records. |
| `M5-02` — authorization versus materialization | Applied | **Confirmed, and well done.** §9.3.1 now names the three things that commit atomically and states that materialized artifacts "are not part of the physical transaction." "Cache invalidation is not an authorization boundary" is the sentence this section needed. |
| `M5-03` — four revocation operations | Applied | **Confirmed.** Four numbered outcomes, per-class decision required, read-time enforcement, and correction preserves "the historical fact that a character previously held the earlier belief." |
| `M5-04` — observable leakage test | Applied | **Confirmed in structure.** `AC-C71` is now a scenario with forced retrieval, summary generation, and payload inspection. One anchoring weakness remains — `M6-03`. |

Two choices in this revision are better than what I asked for:

- §4 states the marker "may be a field on the transition event or a named record," which keeps `M5-01` open for research without leaving the commit point undefined. That is the right altitude.
- `AC-C77` requires correction to preserve the historical belief while updating current truth. I asked for the distinction; the revision supplied the semantics, and "a character previously believed X" is materially different from "X is true."

## 4. Detailed probe — the three areas you named

### 4.1 Is the two-phase alternative implementable without ambiguous recovery states?

**As written, it is a phrase rather than a constraint — and the default ownership layout makes it the likely path, not the exotic one.**

§4 states: "If the transition, first message, scenario identity/revision, and marker may span stores, the design must use an explicit two-phase discipline with a durable prepared state and recovery; it must not perform an unstaged cross-store write." Invariant 9 repeats it. `AC-C75` tests that the commit set "shares a store or uses an explicit two-phase discipline with durable prepared state."

Three things are missing for this to be well-defined.

**Prepared-state completeness.** A two-phase commit is only safe if the prepared state in every participant is sufficient to complete the operation without further input. Here the prepared record must therefore carry the whole first message and the scenario identity/revision, not references that a later step resolves. Nothing states this. If prepared state is a reference, then a crash between prepare and commit forces recovery to *construct* data after the commit decision — which is exactly the ambiguous state you asked about.

**Durability in every participant.** "A durable prepared state" is singular. Two-phase commit requires each participant to have durably prepared before the decision is made durable. §4 does not say the prepared state must be durable in all participants, only that it must be durable.

**Who recovers.** §4.2 assigns the marker to the conversation/event authority and pending records to a combined conversation/scenario authority, but names no component responsible for driving a partially prepared commit to a single outcome. Step 12's "If the outcome is ambiguous, present the pending proposal for a new decision" handles the pending case well; it does not address the case where the marker is durable but a participant's write is not.

That last case is the one that matters. The amendment's rule is "a partial cross-store write is not a valid commit," which is correct and fails closed. But it means recovery must be able to *complete* a commit whose marker is durable. If the participant's prepared record is not sufficient, completion requires inventing data — and if recovery instead discards, the marker is durable while the state it certifies is not, which is the inverse failure.

**Why this is likely rather than theoretical:** §4.2 separates "Conversation events and ordering" from "Scene state and transition events." The first message is a conversation event; the scenario identity/revision is scene state. So the default owner layout puts two members of the commit set under different authorities, which is precisely the configuration that requires the two-phase path.

**Required (`M6-01`), stated at criterion altitude:** require an acceptance criterion that after any interruption the system resolves to exactly one of the two complete states without user input, that prepared state is never authoritative on its own, and that completing or discarding an interrupted commit is deterministic and produces identical results on repeated recovery. Do not prescribe the coordinator.

### 4.2 Does every retrieval path enforce current authorization, including stale summaries and caches?

**The stale-artifact rule is right. The universal quantifier has no mechanism, and one sentence compares against the wrong scope.**

§9.3.1 now says exactly what it should: no artifact may be served under a scope broader than its own recorded authorization; a narrower stale artifact may omit; a broader stale one must fail closed, be withheld, or be filtered; cache invalidation is not an authorization boundary. `AC-C76` tests read-time enforcement. This closes the substance of the finding.

Two residual issues, one normative and one structural.

**The normative sentence compares against recorded authorization.** "No materialized artifact may be served under a scope broader than its own recorded authorization" bounds the served scope by the *artifact's* record. If an entry is demoted, the artifact's recorded authorization is now stale and broader than current. Serving at the artifact's recorded scope is not "broader than its own recorded authorization," so that sentence alone permits the dangerous case. The following sentence ("a stale artifact that is too broad must fail closed") does cover it, and `AC-C76` covers it, so the intent is unambiguous — but the operative comparison should be stated as *current* authorization, because that is the one an implementer will code.

**The universal quantifier is distributed.** §9.3 requires that "retrieval must enforce scope before merging the entry into any character context." §10.5 lists what the backend must know. `AC-C49` enumerates channels. All of these are obligations on each path. Nothing requires that the paths share one authorization decision, so the property holds only as long as every future retrieval surface remembers to implement it. A relationship graph, a timeline index, or an "important moments" feature added in two years would silently not be covered.

This matters more than it would elsewhere because Howl is explicitly planning to grow beyond Chat into studios and broader UI unification. A distributed isolation rule decays with each new surface; a single decision point does not.

I am not asking the amendment to design that point — that is `R-10`/`DA-C7`. I am asking it to state the falsifiable property: a path that returns or injects an entry without consulting current authorization must fail an acceptance test, including a path added after the fact.

**Required (`M6-02`):** change the normative sentence to compare against current authorization, and add a criterion that every retrieval and context-assembly path consults one authorization decision against current scope, proven by a deliberately added bypass path that cannot serve content.

### 4.3 Does `AC-C71` verify what reaches the model, or inspect a receipt that may not represent the payload?

**It inspects both, which is the right instinct. Neither inspection is anchored outside the assembler, so a self-consistent omission passes.**

The rewritten criterion requires: existing Conversation B and new Conversation C, a forced retrieval that would match the lifecycle fact, a generated summary, and inspection of "the assembled payload, and receipt." That is a real improvement and it closes the state-comparison problem.

The residual risk is the one you named. Three points:

**The receipt cannot prove its own completeness.** §5.3 requires an "exact payload hash," but a hash binds a receipt to a byte range the assembler chooses. A receipt listing three events whose concatenation hashes to the recorded value is fully consistent while the request also carries a fourth block outside the hashed region. Integrity is proven; completeness is not.

**The payload inspection has no defined observation point.** If "the assembled payload" is obtained by serializing through the same assembly code, the test proves the assembler is self-consistent, not that the fact is absent from what was transmitted. The observation must be at the transport boundary — the bytes handed to the model backend — with the receipt hash recomputed from those bytes. Otherwise a single defect in the assembly path produces both a leak and an inspection that reports none.

**System framing is where this most likely hides.** `AC-C71` inspects "the response, summary, retrieved material, assembled payload, and receipt." A lifecycle fact written to the character record and rendered into the system prompt or character card appears in none of the enumerated event lists, and may not be enumerated in the receipt at all if the receipt covers only canonical events. The payload clause should be explicit that it includes system framing, character card content, and any instruction block.

**Required (`M6-03`):** anchor the criterion on the transmitted request bytes, require the receipt hash to cover the complete request including framing, and require recomputation at the transport boundary by a component that did not assemble the request.

## 5. Required changes

| ID | Change | Scope |
|---|---|---|
| `M6-01` | Criterion for deterministic single-outcome recovery, prepared state never authoritative, no user input required. | §16.8 |
| `M6-02` | Normative sentence to compare against current authorization; criterion that every path consults one current-authorization decision, with a bypass-path test. | §9.3.1, §16.8 |
| `M6-03` | Anchor `AC-C71` on transmitted bytes; receipt hashes the complete request including framing; recomputation outside the assembler. | §5.3, `AC-C71` |

## 6. Consistency checks requested

| Check | Result |
|---|---|
| Revision 8 preserved at recorded digest | **Pass.** `7f97d359…`, 82,682 B, sidecar `OK`. |
| Delta limited to the four M5 corrections | **Pass.** Changes are §4, §4.1, §4.2, §9.3.1, `AC-C71`, `AC-C74`, `AC-C75`–`AC-C77`, and `CH-13`. |
| `M4-01`–`M4-03` and `M3-01`–`M3-06` remain correctly disposed | **Pass.** No IDs reused or renumbered. |
| `R-09` may not weaken the owner decisions | **Pass.** §12.2 unchanged: "may not relax `OD-02`, `OD-03`, or `OD-04`." |
| No database schema or storage engine selected | **Pass.** The marker is explicitly permitted to be a field or a named record; §5.5 still states the amendment does not select a storage engine or schema. |
| Asynchronous rebuilding never creates a window of broader access | **Pass at the rule level.** §9.3.1 forbids serving a stale broader artifact; `M6-02` addresses the path-coverage side. |
| Cache invalidation is not treated as authorization | **Pass.** Stated explicitly in §9.3.1 and tested by `AC-C76`. |
| Correcting current knowledge does not rewrite historical belief | **Pass.** §9.3.1 outcome 3 and `AC-C77`. |
| Automated lore remains outside the architectural decision | **Pass.** It appears only as an isolation surface. |
| No research, runtime, migration, or implementation authorization | **Pass.** §20 unchanged. |

## 7. This should be the last revision

Nine revisions have taken the document from a requirements capture with seven gaps to a set of constraints where the remaining open questions are empirical rather than editorial.

The evidence for that is in the last three rounds. `M4-01`–`M4-03` were coverage defects: rules without tests. `M5-01`–`M5-04` were representability defects: mechanisms with no home in the model. What remains — whether the two-phase path or a single store is right, whether authorization is one decision point or a layered check, how a receipt can be proven complete against a payload — cannot be settled by more precise prose. Each requires knowing what the existing storage, the existing memory and lore stores, and the existing prompt assembler actually do, and what a small local model under a flushed context window can sustain.

Revision 10 that closes `M6-01`–`M6-03` should be the end of the chain. After it, the next artifact should be `R-07` evidence. Continuing to refine wording past that point would increase the document's apparent authority without increasing what is known — which is the specific risk you flagged, and it is real: a nine-revision document reads as settled in a way a nine-page requirements note does not.

Worth restating for the record, since this is the last delta review: no research has been executed. Every criterion in §16 remains an unvalidated design target. The applied marker, the authorization commit, the read-time scope check, and the transport-anchored receipt are constraints that make the *policies* enforceable. Whether they are the right architecture for Howl is a question this document is not capable of answering.

## 8. Authorization boundary

Per the request's stop condition: no `R-07`–`R-11` research, implementation, migration, runtime investigation, model loading, benchmarking, network testing, Android pairing, or production change follows from this document. Architecture approval and acceptance remain not passed. `OD-01`–`OD-04` remain settled; `O-01`–`O-10` remain open.
