# CORRECTIVE — Stage 3 corrective sequencing design

**Date:** 2026-11-07  
**Mode:** DOCUMENTATION ONLY  
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md`  
**Corrected proposal:** `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md`  
**Status:** Third documentation-only correction, including the D10 scope clarification and the focused-audit minor corrections. This record authorizes nothing.

## 1. Verified hashes

| Artifact | SHA-256 | Size |
|---|---|---|
| Corrective proposal before the third pass | `1fb449ba2baef7b12e42a7b712d2370ac27e6f4dddb9bcb48d6aac3cad87630d` | 710 lines / 65,213 bytes |
| Corrective proposal after the initial third pass | `7b73635d5b1dec52760122567ff3090544f009d071d7422d09c97c48a5738608` | 711 lines / 67,037 bytes |
| Corrective proposal after the D10 scope clarification | `911fbb83a7fe238d9a59c4a29157a5d9bdd8cbb6158481c934dbaacf46150a84` | 711 lines / 68,383 bytes |
| Corrective proposal after the focused-audit minors | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151` | 711 lines / 69,168 bytes |
| Focused D10-scope audit | `30ff59bfe89867824c0abed40c12e71a69e200b5a9c5beae067cf49a87f57e43` | 349 lines / 22,110 bytes |
| Adversarial design audit (round 3) | `6681e11b7eb770d59f3c79294206f05c8d3a3d9d4cfbbd038ef103fc8c48fd15` | 349 lines / 21,676 bytes |
| Adversarial design audit (round 2) | `a7384da9f3f38661fc28b3de48badb82245cfd41efb8e49964f17bcfaaf5a7e2` | 462 lines / 37,477 bytes |
| Adversarial design audit (round 1) | `71a17369e2389fae26e86783fc4df802ab91bad3daef512c0e7f743067e94ea2` | 249 lines / 25,962 bytes |
| Governing architecture | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 lines / 182,904 bytes |
| Stage 3 readiness investigation | `a63f9f32f4d1bd047746a7bbbf5f509588e13ebd4272cffeb23f61fcfa891c6c` | 586 lines / 51,248 bytes |

The round-3 audit reviewed the prior proposal at `1fb449ba…` and the prior record at `ccb28d57…`. It found M5, a major D9 scope defect, and the five minor findings m1–m5. The initial third pass repaired those findings. A subsequent D10 scope challenge is addressed below: D10 is explicitly labeled as two separately answerable direct branches, D10(a) and D10(b), with the A2-D9 supervisor path excluded. The focused D10-scope audit then found n1 and n2; both are corrected in §9. The r3 audit and the focused D10-scope audit remain untouched.

## 2. B1 — direct `Cancel` remains separate from the supervisor path

**Disposition: retained and strengthened.**

C1 keeps three separately implemented branches:

- **Supervisor cancellation path:** records `cancel_requested`; sends a cancel frame only when the capability is known present; when absent or unknown it skips the frame, synthesizes no `unsupported_capability`, and follows the existing grace, `KillTree`, `Wait()`, and `cancelled`/`timeout` precedence (A2-D9). The path uses a path-scoped internal operation so the exemption cannot leak into the exported method.
- **Direct `Cancel`:** remains a direct optional operation. If `howl.tts.cancel` is known absent, it returns `unsupported_capability` with no frame (A2-D1). If the snapshot is unknown, C1 freezes no result; D10(b) governs. The supervisor path remains outside D10.
- **No separate unsupported-capability frame or wire message** exists; capability absence is expressed, if at all, by the A2-D1 result code and never by a frame.

The `client.go` affected-files bullet now names the call site that must switch: **the existing cancellation watcher, which currently invokes exported `Cancel`, must call the internal path-scoped operation instead**. This closes m5 and prevents an implementation from adding the internal operation while leaving the watcher on the exported method.

C1's tests assert the supervisor path and both direct operations separately, including the new direct-`Cancel` unknown-snapshot structural test from m2. That test is D10(b), not an extension of the A2-D9 exemption.

## 3. B2 — unknown-snapshot unload does not reuse a frozen code

**Disposition: retained.**

C1 does not use `provider_not_loaded` for any unload branch. That code remains limited to its registry meaning, “Synthesis before successful load” (§17.9).

- Only the known-present and known-absent direct-unload branches are implemented.
- The unknown branch is structural only: no frame and no spawn, with no result code frozen.
- The no-spawn rule is attributed to §17.7.2's lazy-start row, not §17.19.
- The unresolved caller-facing result remains **D10(a)**, an owner decision with no selected answer and no new code.

D10(a) covers direct `unload`; D10(b) covers direct `Cancel` outside the supervisor path. The A2-D9 supervisor path is excluded. C1 cannot be accepted until both have concrete caller-facing dispositions; a scheduling deferral is not sufficient. Its acceptance text does not claim an unknown-snapshot code for either direct operation.

**Scope note added after the initial third pass.** D10 began as the unload-only unknown-snapshot result question. The direct-`Cancel` unknown branch is a separate direct-operation question: A2-D1 supplies the known-absent rule, while A2-D9 supplies the supervisor-path rule and does not govern an exported direct `Cancel`. The two branches therefore must not be collapsed into one automatic answer. The ledger now labels them D10(a) and D10(b); the owner may choose the same or different concrete result for each, and may defer scheduling the answer. Deferral alone resolves neither branch: C1 acceptance requires a concrete caller-facing result for each. This is an explicit widening of D10's label, not a resolution of either branch and not a change to C1's no-frame/no-spawn rule.

## 4. M1 — C2's `Available=false` projection remains complete

**Disposition: retained; C2 is clear of design findings.**

C2 reproduces A2-D5's full false set: `absent`, `starting`, `handshaking`, `stopping`, `crashed`, `incompatible`, `backoff`, and `unavailable` caused by a missing helper or invalid provider configuration.

Both legal `unavailable` cases remain `Available=true`. C2's matrix and tests include the seven aggregate states and both `unavailable` causes without moving transition ownership into C2.

## 5. M2 — C1 covers every A2-D3 semantic category

**Disposition: retained.**

C1 requires duplicate extension names and invalid handshake descriptors to classify as `protocol_incompatible`, not `protocol_violation`. Both are named in the tests, negative injection, and acceptance.

Wire capability order is not validated. Only normalization and `detail.missing` ordering remain; unordered wire arrays are not rejected. The earlier “stable capability-order validation” instruction is absent.

## 6. M3 — D1a gates C3, not C7

**Disposition: retained.**

D1a is a C3 prerequisite. C3 cannot add any loader seam until D1a is dispositioned. If D1a permits the seam, C3 introduces it; if D1a denies it, C3 is limited to state gating and an internal `Load` on the existing client, with no new adapter/session file.

§7, C3, C6/C7 prerequisites, §9.1, and §10 all record the gate before C3. C7 consumes an already-decided surface; it does not authorize the seam retroactively.

## 7. M4 and M5 — D9 is narrowed to observation/deadline reset

**Disposition: corrected; D9 remains open.**

M4 correctly recorded that idle-deadline reset behavior was not frozen. M5 correctly found that D9 had also asked whether a status request may respawn an idle-closed helper, even though §17.7.2 already makes an explicit status request from `absent` a lazy-start trigger.

The proposal now:

- states that D9 asks **only** whether status/catalog/health observation resets the 120 s idle deadline;
- records status-from-`absent` lazy start as contract-settled by §17.7.2, not as part of D9;
- records that suppressing poll-driven startup, if desired, is an application-layer polling policy—whether the shell issues status requests—not a C5 lifecycle exception;
- gates C5's observation-reset tests and acceptance on D9;
- verifies status-from-`absent` lazy start at C7's `ProviderStatus`, not as a D9 choice.

D9 is still an open product decision; this pass corrects its scope but does not decide it. C5 cannot be accepted until the narrowed D9 is dispositioned.

## 8. D7 — settled by §17.11.5

**Disposition: retained and closed.**

After a successful in-place recovery, readiness resumes; the restart ledger resets only after ten continuous minutes of stable readiness, with prior crashes counted until then. D7 is not an owner decision, has no C6 gate, and must not be reopened.

## 9. Round-3 and focused-audit minor findings

- **m1 — budget wording:** C5 acceptance now says the budget permits three restarts within ten minutes and a fourth attempt exhausts it, instead of “exactly at 3.” The existing three-restarts-then-unavailable test remains.
- **m2 — direct `Cancel` unknown snapshot:** C1 now has a structural test for direct `Cancel` with unknown capability knowledge: no frame, no spawn, and no asserted result code. D10(b) governs; this is separately answerable from D10(a)'s direct-`unload` branch, not an automatic same-code decision. The A2-D9-settled supervisor path is outside both.
- **m3 — §9.1 acceptance boundary:** the same-day C1/C2 bullet now says C1 cannot be accepted until D10(a) and D10(b) each have a concrete caller-facing disposition, and a deferral does not satisfy that gate, while C2 has no open design gate.
- **m4 — undefined frame term:** “unsupported-operation frame” has been removed. The proposal states that no separate unsupported-capability frame or wire message exists and that capability absence is expressed only by the A2-D1 result code.
- **m5 — watcher call site:** the C1 `client.go` bullet names the existing cancellation watcher as the call site that must switch from exported `Cancel` to the internal path-scoped operation.

**Focused D10-scope audit minors.**

- **n1 — §1 owner-decision count:** the summary now says five owner decisions (D1b, D2, D5, D9, D10), with D10 also a contract-surface gap. The former four-item count understated the gate set and is removed.
- **n2 — deferral is not a result:** every D10 acceptance site now requires a concrete caller-facing disposition for each branch and says a scheduling deferral alone does not satisfy the gate. The D10 ledger no longer lists deferral as a candidate behavioral result; it lists only concrete outcomes and states that deferral may schedule the owner's answer but cannot resolve it.

## 10. Prior corrections preserved

- F1's direct `unload`/cancellation split is preserved and strengthened by the named watcher call site.
- F2's two `unavailable` cases are preserved.
- F3's single-writer ownership of `Loaded` is preserved.
- F4's Windows real-child evidence rule is preserved.
- F6 is represented by the narrowed D9.
- F7's C4 intent diagnostic is preserved.
- F8's process-counted cap and process/job-handle registry are preserved.

No prior correction was narrowed or contradicted by this pass. D10 is explicitly widened from its unload-only origin to two labeled direct branches, D10(a) and D10(b); that widening is disclosed here and in the proposal, and neither branch is resolved. The focused-audit n1 count repair and n2 concrete-result gate change no decision outcome and do not authorize either branch.

## 11. Decision ledger after this pass

| ID | Disposition | Timing |
|---|---|---|
| D1a | Open contract-surface gap: §17.12.1 has no `load`, `unload`, or `cancel` entry point; a narrow internal seam needs an owner/adoption decision. | Before C3 introduces a seam; C3 is not authorizable until resolved. |
| D1b | Genuine product/architecture decision: explicit action and adapter lifetime. | Before C7 production reachability. |
| D2 | Genuine product/architecture decision: helper-backed activation without changing `TTSSynthesizer`; O-4 remains open. | Before C7 production activation. |
| D3 | Frozen cap; precise admission-rejection error surface remains unspecified. Do not invent or repurpose a code. | Before C7 cap-response acceptance. |
| D4 | Folded into D2/D5. | No separate gate. |
| D5 | Genuine product decision: renderer exposure and O-1 picker behavior. | Before UI integration. |
| D6 | Settled and applied: forced termination is `stopping` + grace → `crashed`; outcome remains `cancelled`/`timeout`; clean exit 0 → `absent`. | Applied in C4. |
| D7 | Settled by §17.11.5: ledger resets only after ten continuous minutes of stable readiness; prior crashes count until then. | Closed. |
| D8 | Deferred B12a/B12b contract issues. | Not resolved here. |
| D9 | Open product decision, narrowed: does status/catalog/health observation reset the 120 s idle deadline? Status-from-`absent` lazy start is settled by §17.7.2 and is not part of D9; poll suppression is an application-layer policy. | Before C5 acceptance. |
| D10 | Open owner decision, explicitly split: D10(a) unknown-snapshot direct `unload` result; D10(b) unknown-snapshot direct `Cancel` result outside the supervisor path. These are separately answerable and may receive the same or different concrete result. A2-D9 settles only the supervisor path, which is outside D10. No code selected for either; a concrete caller-facing result is required for each, and deferral alone is not a resolution. | Before C1 acceptance (both branches, concrete results). |

O-1–O-7 remain unchanged. C1 and C2 remain unauthorized, and no later stage is authorized by this record.

## 12. Sections changed in the proposal

- Header revision note; §1 D9/D10 scope and owner-decision count.
- §4.4 unsupported-capability frame wording.
- C1 `client.go` watcher call site and D10(a)/D10(b) direct branches; direct-`Cancel` unknown test; concrete D10 acceptance gate; frame-terminology wording; acceptance.
- C5 prerequisite, idle-reset rule, tests, affected `client.go` note, and budget acceptance.
- §9.1 C1/D10 acceptance boundary, concrete-result gate, and C5/D9 wording.
- §10 D9 introduction and ledger row; D10 scope, ledger row, and concrete-result gate.

The C1–C8 order, critical path, inert C1–C6 boundary, C7 default-off gate, C8 evidence rules, and native-provider exclusion are preserved.

## 13. Preservation and write set

- The governing architecture contract remains byte-identical at `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f`.
- The Stage 3 readiness investigation remains byte-identical at `a63f9f32f4d1bd047746a7bbbf5f509588e13ebd4272cffeb23f61fcfa891c6c`.
- The round-1 audit remains byte-identical at `71a17369e2389fae26e86783fc4df802ab91bad3daef512c0e7f743067e94ea2`; the round-2 audit remains `a7384da9…`; the round-3 audit remains untouched at `6681e11b7eb770d59f3c79294206f05c8d3a3d9d4cfbbd038ef103fc8c48fd15`. The focused D10-scope audit remains byte-identical at `30ff59bfe89867824c0abed40c12e71a69e200b5a9c5beae067cf49a87f57e43`.
- No Stage 1/Stage 2 report, correction, or audit was modified.
- No production source, test, fixture, protocol, interface, manager, pipeline, adapter, or native-provider file was modified.
- No API, UI, dependency, model, or asset was added.
- No build, test, replay, benchmark, process launch, download, commit, push, or stage was performed.

Exact changed-file inventory for this pass:

1. `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` — corrected in place.
2. `docs/review/CORRECTIVE-howl-tts-stage3-corrective-design-2026-11-07.md` — this correction record.

The README index was not modified because this pass was restricted to the proposal and correction record. Any index entry that still cites `1fb449ba…` or `ccb28d57…` is historical after this pass.

C1 and C2 remain unauthorized. C1 cannot be accepted until D10(a) and D10(b) each have a concrete caller-facing result; C2 has no open design gate; C3 is gated by D1a; C5 is gated by the narrowed D9; C7 still requires D1b, D2, and D3. Stage 3 implementation requires explicit, stage-specific authorization, and C8 requires separate runtime-execution authorization.
