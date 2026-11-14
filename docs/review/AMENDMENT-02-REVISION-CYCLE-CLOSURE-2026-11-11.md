# Amendment 02 — revision-cycle closure

**Date:** 2026-11-11  
**Status:** Requirements baseline closed at revision 9; no revision 10; research planning only  
**Scope:** Documentation and evidence-obligation handoff; no implementation, runtime, migration, model, benchmark, or research execution

## 1. Closure decision

The owner closes the Amendment 02 revision cycle at revision 9.

Muse's revision-9 review returned **PASS WITH REQUIRED CHANGES** for `M6-01`–`M6-03`. Those findings are not being converted into another amendment revision. They are recorded here as unresolved evidence obligations for the research process.

This decision does not approve the architecture, validate the requirements, authorize implementation, or authorize research execution. It ends the editorial correction cycle; it does not end the evidence work.

## 2. Preserved artifacts

| Artifact | Identity |
|---|---|
| Amendment 02, revision 9 | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`; `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640`; `88,499 B / 1,303 lines` |
| Revision 9 sidecar | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md.sha256`; `OK` |
| Muse revision-9 review | `docs/review/REVIEW-AMENDMENT-02-REVISION-9-MUSE-DELTA-2026-11-11.md`; `381fac1ce84e91f4e20d4649538598513f3f4de6e5bf0809bcb4d3574806a6e8`; `16,841 B / 146 lines` |
| Revision 8 archive | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-8-2026-11-11.md`; `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f` |
| Revision 7 archive | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-7-2026-11-11.md`; `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1` |
| Revision 6 archive | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-6-2026-11-11.md`; `ef2d5dd88fc30cf49424f6bcf20f50aca5c599ae6d9a8f76a33bc65673370433` |

Revision 9 remains the current requirements document. No revision 10 is created.

## 3. Evidence obligations carried forward

### `EO-01` — deterministic cross-store recovery

**Source:** `M6-01`  
**Research owner:** `R-07` / `DA-C1`, with `R-09` where transition state is involved.  
**Obligation:** Determine whether the commit set can share a store or needs a two-phase protocol. For any two-phase path, establish that prepared state is complete, durable in every participating store before the commit decision, never authoritative on its own, and sufficient to finish without reconstructing data. After any interruption, the system must resolve to exactly one complete state without user input; repeated recovery must be deterministic and idempotent.

**Evidence required:** actual storage ownership, event and scene-state persistence, transaction capabilities, prepared-state contents, recovery behavior, and failure injection. The coordinator, if any, is a research result, not a predetermined requirement.

### `EO-02` — current-authorization enforcement across all paths

**Source:** `M6-02`  
**Research owner:** `R-10` / `DA-C7`, with `R-07` for persistence and context assembly.  
**Obligation:** Determine how every retrieval and context-assembly path consults current authorization rather than an artifact's stale recorded scope. Include a deliberately added bypass path that must fail to serve content. Investigate whether a shared authorization decision point is appropriate; do not assume that every subsystem can safely implement its own check.

**Evidence required:** actual retrieval surfaces, cache and index behavior, prompt assembly, path coverage, bypass testing, and the operational cost of a shared decision point.

### `EO-03` — transport-boundary payload verification

**Source:** `M6-03`  
**Research owner:** `R-07` / `R-10.6`, with the model-input receipt requirements in `RQ-07.4`.  
**Obligation:** Anchor verification on the actual request bytes handed to the model backend, not a receipt produced by the assembler. The receipt hash must cover the complete request, including system prompts, character cards, instruction blocks, and other framing material. Recompute and verify at the transport boundary using a component that did not assemble the request. Exercise existing Conversation B and new Conversation C, including forced retrieval and summary generation.

**Evidence required:** the actual backend handoff, payload framing, receipt completeness, independent recomputation, and observed non-delivery of lifecycle facts and derived forms.

## 4. What this record does not establish

- It does not validate the proposed mechanisms.
- It does not prove that two-phase commit, read-time authorization, or transport-boundary verification are the best designs.
- It does not approve the architecture or pass acceptance.
- It does not authorize R-07–R-11 execution, implementation, migration, runtime investigation, model loading, benchmarking, network testing, or Android pairing.
- It does not settle `O-01`–`O-10` or the remaining `DA-C1`–`DA-C12` design questions.
- It does not fold the separate automated-lore architectural question into this amendment.

## 5. Research handoff

Research planning may proceed under the existing gates. The first evidence artifact should be `R-07`, because the storage, event-model, recovery, and payload questions block a credible evaluation of the rest.

Research must be allowed to conclude that:

- the proposed storage arrangement is wrong;
- two-phase commit is unnecessary under a better ownership model;
- a different memory architecture makes some of these mechanisms unnecessary;
- the existing prompt assembler requires substantial redesign;
- automated lore should be replaced, made optional, or retired.

The requirements state what Howl must guarantee. Research must still determine the best way to guarantee it.

## 6. Closure checks

- Revision 9 was not rewritten after review.
- The revision-9 sidecar verifies.
- Revision 8, revision 7, and revision 6 archives remain byte-identical.
- The Muse review and prior review-chain hashes remain unchanged.
- `M6-01`–`M6-03` are preserved here as `EO-01`–`EO-03`.
- No revision 10 exists.
- No research execution or implementation was performed.
