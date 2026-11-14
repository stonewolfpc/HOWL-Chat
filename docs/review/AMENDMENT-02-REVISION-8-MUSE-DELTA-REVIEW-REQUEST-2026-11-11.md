# Amendment 02 — Revision 8 delta review request

**Date:** 2026-11-11  
**Reviewer:** Muse  
**Review type:** Focused documentation-delta review  
**Subject:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`, revision 8  
**Revision 8 SHA-256:** `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f`  
**Revision 8 size:** `82,682 B / 1,271 lines`  
**Prior reviewed identity:** revision 3, `fa766f76779e2517b31bd55da993f872525a52941b348fc8f78bbeb1d9382c44`  
**Preserved revision 7:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-7-2026-11-11.md`, `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1`

## 1. Request

Review the revision-8 delta that applies the three M4 findings from the revision-7 delta review. Determine whether it closes those findings without reopening `OD-02`, `OD-03`, or `OD-04`, changing the research scope, or implying implementation authorization.

This is not a request for architecture approval, acceptance, research execution, or runtime investigation.

## 2. Required corrections to inspect

| Finding | Required response |
|---|---|
| `M4-01` — commit boundary and restart | Confirm §7.1 distinguishes proposal, pending transition, committed transition, and rejected/interrupted transition; the durable applied marker is the sole commit point written atomically with the transition, first message, and scenario identity/revision; a record without the marker is non-active; restart does not infer approval or re-apply; `AC-C73` tests complete pre-state or complete post-state, never a hybrid. |
| `M4-02` — promotion atomicity and revocation | Confirm §9.3.1 makes promotion one atomic operation across scope, provenance, indexes, summaries, and caches; failure leaves the prior scope authoritative; `R-09`/`DA-C4` must explicitly decide what demotion is supported; demotion does not claim to erase information already consumed; unsupported demotion has a stated remedy; `AC-C74` tests the behavior. |
| `M4-03` — derived surfaces and new conversations | Confirm `AC-C71` covers canonical records, summaries, generated/consolidated and lifecycle-derived memories, automatic lore, retrieval indexes/results, model-input receipts, cached context, other persisted derivatives, and conversations created after the lifecycle event. |

## 3. Consistency checks

Confirm that:

- revision 7 is preserved at the recorded digest;
- the revision-8 change is limited to the three M4 corrections and their traceability;
- `M3-01`–`M3-06` remain correctly disposed;
- `R-09` may not weaken the owner decisions;
- the new criteria are design targets, not claims of implemented behavior;
- no research, runtime, migration, or implementation authorization is implied;
- the separate automated-lore research-scope question is not folded into this correction.

## 4. Required verdict

Return one of:

- **PASS** — the revision-8 delta closes `M4-01`–`M4-03` without weakening the owner policies or expanding scope.
- **PASS WITH REQUIRED CHANGES** — list exact changes; the delta is not closed until they are applied and rechecked.
- **FAIL** — identify the missing, weakened, or contradicted requirements.

A PASS is a documentation-delta result only. It is not architecture approval, acceptance, or research authorization.

## 5. Stop condition

After the verdict, stop. Do not begin R-07–R-11 research, implementation, migration, runtime investigation, model loading, benchmarking, network testing, or Android pairing. Those remain subject to separate owner authorization and the existing gates.
