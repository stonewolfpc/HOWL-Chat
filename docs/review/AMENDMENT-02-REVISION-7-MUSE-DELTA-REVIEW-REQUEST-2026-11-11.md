# Amendment 02 — Revision 7 delta review request

**Date:** 2026-11-11  
**Reviewer:** Muse  
**Review type:** Focused documentation-delta review  
**Subject:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`, revision 7  
**Revision 7 SHA-256:** `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1`  
**Revision 7 size:** `76,565 B / 1,241 lines`  
**Prior reviewed identity:** revision 3, `fa766f76779e2517b31bd55da993f872525a52941b348fc8f78bbeb1d9382c44`  
**Preserved revision 6:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-6-2026-11-11.md`, `ef2d5dd88fc30cf49424f6bcf20f50aca5c599ae6d9a8f76a33bc65673370433`

## 1. Request

Review the revision-7 delta that applies the OD-04 consistency audit. Determine whether it makes `OD-04` enforceable without weakening `OD-02` or `OD-03`, adding implementation authorization, or turning design questions into owner decisions.

This is not a request for architecture approval, implementation authorization, runtime investigation, or research execution.

## 2. Findings to inspect

| Finding | Required response |
|---|---|
| `M3-01` — no `OD-04` acceptance criterion | Confirm `AC-C67`–`AC-C72` cover preservation, confirmation ordering, scope visibility/promotion, memory scope, lifecycle isolation, and change history. |
| `M3-02` — transition has no confirmation gate or decline path | Confirm §7.1 now gates active state, first message, participant state, memory, and model context behind confirmation, and defines decline as no canonical transition or first message. |
| `M3-03` — Chat-typed corrections have no defined scope | Confirm §9.3 makes Chat-derived entries conversation-scoped by default and `AC-C70` tests non-retrieval in another chat. |
| `M3-04` — chat lifecycle lacks cross-chat isolation criteria | Confirm `AC-C71` tests death, Step Out, return, and ghost participation against another chat's state and preserves originating-chat history. |
| `M3-05` — change record is required but not tested | Confirm §7.2 and `AC-C72` require actor, reason, scope, previous value, resulting value, transition identity, and confirmation outcome, surviving reload and remaining retrievable. |
| `M3-06` — `AC-C66` sequence | Confirm it remains informational; ID stability is preserved. |

## 3. Additional consistency checks

Confirm that:

- `AC-C23` no longer implies a transition applies before confirmation.
- `R-09` is explicitly prevented from weakening `OD-02`, `OD-03`, or `OD-04`.
- revision 6 is preserved at the recorded digest and is not silently rewritten.
- the new criteria are design targets, not claims of implemented behavior.
- no research, runtime, migration, or implementation authorization is implied.

## 4. Required verdict

Return one of:

- **PASS** — the revision-7 delta closes the M3 requirements without weakening the owner policies.
- **PASS WITH REQUIRED CHANGES** — list exact changes; the delta is not closed until they are applied and rechecked.
- **FAIL** — identify the missing, weakened, or contradicted requirements.

A PASS is a documentation-delta result only. It is not architecture approval, acceptance, or research authorization.

## 5. Stop condition

After the verdict, stop. Do not begin R-07–R-11 research, implementation, migration, runtime investigation, model loading, benchmarking, network testing, or Android pairing. Those remain subject to separate owner authorization and the existing gates.
