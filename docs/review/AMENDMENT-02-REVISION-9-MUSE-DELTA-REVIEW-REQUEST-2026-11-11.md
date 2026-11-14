# Amendment 02 — Revision 9 delta review request

**Date:** 2026-11-11  
**Reviewer:** Muse  
**Review type:** Focused documentation-delta review  
**Subject:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`, revision 9  
**Revision 9 SHA-256:** `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640`  
**Revision 9 size:** `88,499 B / 1,303 lines`  
**Prior reviewed identity:** revision 3, `fa766f76779e2517b31bd55da993f872525a52941b348fc8f78bbeb1d9382c44`  
**Preserved revision 8:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-8-2026-11-11.md`, `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f`

## 1. Request

Review the revision-9 delta that applies the four M5 findings from the revision-8 delta review. Determine whether it closes those findings without reopening `OD-02`, `OD-03`, or `OD-04`, selecting a database schema, changing the research scope, or implying implementation authorization.

This is not a request for architecture approval, acceptance, research execution, or runtime investigation.

## 2. Required corrections to inspect

| Finding | Required response |
|---|---|
| `M5-01` — commit-state ownership | Confirm §4 carries `commitState` and an applied-marker reference or named equivalent; §4.1 forbids projection of uncommitted records; §4.2 assigns the marker and pending records; and the model states that the commit set shares a store or uses an explicit two-phase discipline with durable prepared state. |
| `M5-02` — authorization versus materialization | Confirm §9.3.1 commits scope, provenance, and visibility atomically; summaries, indexes, and caches may rebuild asynchronously; no artifact may be served under a scope broader than its recorded authorization; stale narrower artifacts may omit, stale broader artifacts fail closed; and `AC-C74`/`AC-C76` test it. |
| `M5-03` — revocation operations | Confirm §9.3.1 distinguishes access revocation, derived-copy removal, correction, and consumption history; requires a per-class `R-09`/`DA-C4` decision; enforces at read/authorization time; and `AC-C77` preserves historical belief while updating current truth. |
| `M5-04` — observable leakage test | Confirm `AC-C71` runs in existing Conversation B and new Conversation C, forces retrieval that would match the lifecycle fact, generates a summary, assembles the actual model input, and inspects its receipt for the fact and derived forms. |

## 3. Consistency checks

Confirm that:

- revision 8 is preserved at the recorded digest;
- the revision-9 change is limited to the four M5 corrections and their traceability;
- `M4-01`–`M4-03` and `M3-01`–`M3-06` remain correctly disposed;
- `R-09` may not weaken the owner decisions;
- no database schema or storage engine is selected;
- asynchronous rebuilding never creates a window of broader access;
- cache invalidation is not treated as authorization;
- correcting current knowledge does not rewrite historical belief;
- automated lore remains outside the architectural decision;
- no research, runtime, migration, or implementation authorization is implied.

## 4. Required verdict

Return one of:

- **PASS** — the revision-9 delta closes `M5-01`–`M5-04` without weakening the owner policies or expanding scope.
- **PASS WITH REQUIRED CHANGES** — list exact changes; the delta is not closed until they are applied and rechecked.
- **FAIL** — identify the missing, weakened, or contradicted requirements.

A PASS is a documentation-delta result only. It is not architecture approval, acceptance, or research authorization.

## 5. Stop condition

After the verdict, stop. Do not begin R-07–R-11 research, implementation, migration, runtime investigation, model loading, benchmarking, network testing, or Android pairing. Those remain subject to separate owner authorization and the existing gates.
