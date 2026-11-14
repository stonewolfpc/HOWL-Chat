# Muse — Revision 8 delta review

**Date:** 2026-11-11
**Reviewer:** Muse
**Review type:** Focused documentation-delta review (read-only)
**Subject:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`, revision 8
**Request:** `docs/review/AMENDMENT-02-REVISION-8-MUSE-DELTA-REVIEW-REQUEST-2026-11-11.md`
**Status:** Complete; no document or source file modified

## 1. Verdict

**PASS WITH REQUIRED CHANGES**

Revision 8 closes `M4-01` through `M4-03` at the requirements level. §7.1 now separates proposal, pending, committed, and rejected/interrupted states and makes the applied marker the sole commit point; §9.3.1 introduces promotion atomicity, failure-to-prior-scope, and an explicit revocation decision; `AC-C71` now enumerates the derived surfaces and covers conversations created after the event. `AC-C73` and `AC-C74` were added rather than retrofitted onto existing IDs.

Four required changes remain, and they are of a different kind from the previous rounds. Revisions 4–8 have been closing gaps in coverage — rules without tests, tests without rules. The remaining items are places where the new mechanisms are **specified in a way that cannot be implemented as written, or specified without a home in the model that describes the data**. These matter because they will be resolved silently during implementation, and the resolutions will differ from reviewer to reviewer.

They are, briefly:

- the applied marker has no representation in §4's event model and no owner in §4.2;
- promotion atomicity is stated as physical rebuild of indexes, summaries, and caches, which is not achievable and conflates authorization with materialization;
- revocation is one word covering four different operations, and two of them are unaddressed;
- `AC-C71` tests stored state for invariance rather than testing retrieval and context assembly for leakage.

Details in §4.

This verdict is a documentation-delta result only. It is not architecture approval, acceptance, research authorization, or permission to begin any workstream.

## 2. Identity verification

| Artifact | Expected | Recomputed | Result |
|---|---|---|---|
| Amendment 02, revision 8 | `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f` | same | Match |
| Size / lines | `82,682 B / 1,271 lines` | same | Match |
| Sidecar | `sha256sum -c` | `OK` | Verified |
| Revision 7 archive | `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1` | same | Match |
| Revision 7 archive size | `76,565 B` | same | Match |
| Revision 7 archive sidecar | `sha256sum -c` | `OK` | Verified |
| Revision 6 archive | `ef2d5dd8…` / `71,428 B` | same | Unchanged |
| Revision-7 delta review | `8e7f248b07d4a6e5de6e46f8c509a2e685b58e0ba4181118a631a67507baf23c` | same | Unchanged |
| Revision-6 OD-04 addendum | `b7be1d46…` | same | Unchanged |
| Integration plan | `1e92a544…` | same | Unchanged |
| Amendment 01 | `f3a0f8b9…` | same | Unchanged |

The archive chain is intact and each revision is preserved as a separate file at its recorded digest. Nothing was rewritten in place.

## 3. Disposition of the M4 findings

| Finding | Claimed | My assessment |
|---|---|---|
| `M4-01` — commit boundary and restart | Applied | **Confirmed.** §7.1 four-state vocabulary, step 10 marker as sole commit point, steps 7–8 display only at commit, step 11 model context only after marker, step 12 restart recovery, `AC-C73` hybrid test. Substantive. |
| `M4-02` — promotion atomicity and revocation | Applied | **Confirmed at the rule level**, with implementability problems — see `M5-02`, `M5-03`. |
| `M4-03` — derived surfaces and new conversations | Applied | **Confirmed for enumeration**, but the test remains a state check — see `M5-04`. |

Two details in this revision deserve credit because they close failure paths the review did not name:

- Steps 7 and 8 were rewritten from "display" to "display … only as part of the commit in step 10." That removes the ordering hazard where the first message becomes visible before the transition is durable — a path that would have been reachable even with a correct marker rule.
- Step 12's "If the outcome is ambiguous, present the pending proposal for a new decision" is the right resolution for the genuinely ambiguous case, and it is stricter than inferring from the surviving record.

## 4. Detailed probe — the three requested challenges

### 4.1 Is the applied marker genuinely authoritative?

**Almost. The rule is correct and unusually complete; the mechanism is homeless.**

The text now satisfies the distinction you drew between persisted and committed. §7.1's final paragraph is explicit: "The UI, persistence layer, and model-context assembly must all read the applied marker as the authority; a pending record is never a half-transition." Each of your five checks maps to a sentence: UI at steps 7–8, model context at step 11, recovery at step 12, duplicate retries at `AC-C58` plus step 12's resume/discard rule, and hybrid state at `AC-C73`.

The problem is one level down. §7.1 now depends on four record states and a dedicated durable marker, but §4's canonical event model does not contain them:

- §4's `ConversationEvent` shape has `kind`, `actor`, `body`, `audience`, `sourceEventIds`, `modelInputReceiptIds`, and `createdAt`. It has no commit status, and no separate marker is shown.
- §4.1's seven invariants do not include "an uncommitted record is never projected, displayed, or read by context assembly." That invariant is now the load-bearing one, and it exists only in §7.1.
- §4.2's owner table assigns "scene state and transition events" to the scenario/conversation state authority, but assigns no owner to the applied marker, to pending records, or to pending first messages. The one record everything now reads as authority has no declared home.

This is not pedantry about a provisional schema — §4 explicitly labels its shape a hypothesis for `R-07`. It matters because two plausible implementations differ materially:

**Marker as a field on the transition event.** Commit is a single-record write; "atomically persist the transition, first message, scenario/revision, and marker" means one record plus its references, which is achievable in an ordinary transaction. Step 10 becomes a normal commit.

**Marker as a separate durable record.** Commit is a cross-record transaction. That is still achievable in a single store and not achievable across the stores §4.2 implies. The requirement "Either all of these become durable or none does" then forces a design constraint that no section states.

A second consequence: an abandoned pending first message written to the canonical event log is a canonical event that is neither visible nor committed. The term "non-active" in §7.1 has no counterpart in §4's vocabulary, so an implementer has to decide whether it is suppressed at projection, quarantined, or garbage-collected — and §4.1 invariant 4 ("A visible event has one canonical record") reads as if every canonical event is eligible for projection.

**Required (`M5-01`):** carry commit state in the provisional model — either as an explicit field on the transition event or as a named marker record — add the invariant that an uncommitted record is never projected, displayed, or readable by context assembly, assign ownership of the marker and pending records in §4.2, and state whether the four commit-set members must share a store or may span stores under an explicit two-phase discipline.

### 4.2 Does atomic promotion account for all derived surfaces?

**The surfaces are right; the atomicity is specified as physical and cannot be satisfied as written.**

§9.3.1's first sentence requires "one atomic, auditable operation covering the entry's scope change, its provenance record, and every affected retrieval view, index, summary, or cache."

Two problems, and the second is the serious one.

**Physical rebuild is not atomic and not always possible.** A summary is produced by a model. It cannot be regenerated inside a transaction, and if it must be, promotion becomes a long-running operation holding a write lock for the duration of an inference. For an on-device system whose whole premise is shared, constrained hardware, that is the worst possible shape. The same applies to any index with a non-trivial build step.

**The requirement conflates authorization with materialization.** Promotion does not change a summary's text. It changes who may retrieve that summary and the entry behind it. The correct atomic set is the *authorization state* — scope, provenance, and whatever visibility record retrieval consults. Materialized artifacts do not need rebuilding; they need re-gating. Writing "every affected … summary, or cache" into the atomic set makes the requirement both unimplementable and, worse, encourages an implementation that believes it has satisfied atomicity by rebuilding artifacts rather than by controlling access.

Your formulation is the correct one and the section is one sentence away from it: the second sentence already says "The broader scope is not authoritative until the complete promotion commits," which is an authorization statement. The fix is to move that to the front and demote materialization to an asynchronous obligation.

There is also a stale-artifact asymmetry worth naming. An artifact built under the *narrower* scope and served after promotion fails closed — it simply under-serves. An artifact built under the *broader* scope and served after demotion fails open, which is the dangerous direction. The section constrains the first implicitly and does not constrain the second at all.

**Required (`M5-02`):** restate promotion as an atomic *authorization* commit over scope, provenance, and visibility, with materialized artifacts re-gated and rebuildable asynchronously; require that no artifact is served under a scope broader than its own recorded authorization; and add the matching criterion. `AC-C74` currently tests "no partial retrieval exposure," which is the right outcome but does not force the right mechanism.

### 4.3 Is revocation one operation or several?

**Several, and the section collapses them.**

§9.3.1 handles two of the four distinctions you listed, and handles the hardest one well: revocation "stops future retrieval under the broader scope without erasing information already delivered to, or legitimately observed by, another conversation; it must not claim to recall what was already consumed." That is exactly right — it refuses to promise something no system can deliver.

The other two are unaddressed:

**Removing derived copies.** A consolidated memory that merged this entry now contains it. §9.4's shared-source rule governs *deletion* of a source message, not the demotion of an entry that has been merged forward. Nothing says whether demotion must split a consolidated entry, mark it as containing demoted material, or leave it intact.

**Correcting previously consolidated information.** If the promoted fact was wrong and is now retracted, "stop future retrieval" is not a correction — it is disappearance. A character may have been operating on it. The section provides no vocabulary for a retraction that differs from a scope change.

**Read-time enforcement.** "Stops future retrieval under the broader scope" is a requirement with no mechanism. Since §9.3.1 also makes caches part of the picture, the only enforcement that survives unbounded caches is an authorization check at read time against current scope. Anything that depends on cache invalidation will fail open eventually, and fail open silently.

**Required (`M5-03`):** name the four operations in §9.3.1 — revoke future access, remove derived copies, correct consolidated or previously circulated information, and acknowledge what was already consumed — require `R-09`/`DA-C4` to state, per class of derived information, which of the four a given revocation performs, and state that revocation is enforced at read/authorization time and does not depend on cache invalidation.

### 4.4 Does `AC-C71` test leakage or record classification?

**Record classification. The enumeration is excellent; the verb is wrong.**

The expanded criterion is a substantial improvement — it covers canonical character records, summaries, generated and consolidated memories, lifecycle-derived memories, automatic lore, retrieval indexes and results, model-input receipts, cached context, other persisted derivatives, and conversations created after the event. That closes the enumeration gap I raised, including the new-conversation case.

But the criterion's testable content is "leaves the character's state … unchanged" and "remains visible." Both are comparisons of stored state. Nothing requires the unaffected conversation to *attempt retrieval and context construction* and to be *observed failing to receive the lifecycle fact*.

You identified the exact failure this misses: a global status flag read only during new-chat creation leaves every stored record unchanged, and the criterion passes. The enumeration includes "cached context," but a flag consumed at prompt-assembly time from a global record is neither a cache nor a persisted derivative of the conversation — it is transient input.

The comparison is with `AC-C62`, which is written the right way: "Character A later receives D's response but not the image, a private description, or a summary revealing its contents." That is an observed non-receipt. `AC-C71` should be written in the same voice — a scenario with a forced attempt, not an invariant with a surface list.

**Required (`M5-04`):** rewrite `AC-C71` or add a criterion in `AC-C62`'s voice: in an existing conversation B and in a conversation C created after the lifecycle event, request a character response, deliberately force a retrieval that would match the lifecycle fact, generate a summary, and inspect the assembled model input and its receipt to demonstrate the fact and every derived form are absent.

## 5. Required changes

| ID | Change | Scope |
|---|---|---|
| `M5-01` | Place commit state in the model: field or marker record, non-projection invariant, §4.2 ownership, store-spanning rule. | §4, §4.1, §4.2 |
| `M5-02` | Separate atomic authorization commit from asynchronous materialization; forbid serving under a broader scope than recorded authorization; add criterion. | §9.3.1, §16.8 |
| `M5-03` | Name the four revocation operations; require the per-class decision; require read-time enforcement. | §9.3.1 |
| `M5-04` | Rewrite `AC-C71` as an observed non-receipt scenario with forced retrieval and prompt inspection. | §16.8 |

## 6. Consistency checks requested

| Check | Result |
|---|---|
| Revision 7 preserved at recorded digest | **Pass.** `4cd76064…`, 76,565 B, sidecar `OK`. |
| Delta limited to the three M4 corrections and traceability | **Pass.** Changes are §7.1, §9.3.1, `AC-C71`, `AC-C73`, `AC-C74`, and `CH-12`. `M3-01`–`M3-06` remain as previously disposed. |
| `M3-01`–`M3-06` remain correctly disposed | **Pass.** No renumbering; `AC-C66` still a stable non-contiguous ID. |
| `R-09` may not weaken the owner decisions | **Pass.** §12.2's prohibition is unchanged. |
| New criteria are design targets, not behavior claims | **Pass.** §16.8 states "These criteria make the owner decisions testable at the acceptance gate. They are unvalidated design targets." |
| No research, runtime, migration, or implementation authorization | **Pass.** §20 unchanged; §24 unchanged. |
| Automated-lore architectural question not folded in | **Pass.** The revision adds isolation obligations only. A search for retirement, overhaul, optionality, or architectural-role language returns nothing; `AC-C71` names automatic lore as a surface that must obey current isolation rules, which is the correct boundary. |

## 7. On what this chain does and does not establish

Eight revisions in, the document is now internally coherent about *requirements*. That is real progress and it is worth stating plainly what it is not: it establishes nothing about which mechanisms should exist.

No research has been executed. Every criterion in §16 remains an unvalidated design target. `R-07`–`R-11` remain question sets, not dossiers. The mechanism the last three revisions have been refining — an applied marker, an atomic authorization commit, read-time scope checks — are sound requirements that no evidence yet says are the right design for a small local model under a flushed context window. They are the shape that makes the *policy* enforceable; whether they are the shape that makes the *product* work is exactly what research is for.

That distinction is now load-bearing, because the revision chain is long enough to create its own momentum. A ninth revision that tightens wording further would add precision without adding knowledge.

On automated lore specifically: keeping its architectural future out of this amendment was right, and the reason you gave is the right reason. The isolation requirement is a constraint on whatever automated lore exists now; it is not an endorsement that it should exist in its current form. Whether it survives, gets overhauled, becomes optional, or is retired is a research question, and inheriting the limitations of the first memory experiment by inertia is precisely the failure mode a requirement document cannot prevent. A future automated-lore decision should be reached on evidence, not because the amendment described it as a surface.

## 8. Authorization boundary

Per the request's stop condition: no `R-07`–`R-11` research, implementation, migration, runtime investigation, model loading, benchmarking, network testing, Android pairing, or production change follows from this document. Architecture approval and acceptance remain not passed.
