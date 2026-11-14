# CROSS-CUTTING ENGINEERING REQUIREMENTS — Thread Safety, Hardware Safety, and Performance Preservation

**Revision:** r6 corrective  
**Date:** 2026-11-07  
**Status:** PROPOSED CORRECTIVE REVISION — NOT ADOPTED — NOT BINDING — A-1 SCOPE DECISION PENDING  
**Pre-review assessed predecessor:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r5-2026-11-07.md` (SHA-256 `93cfe1de5f362fba6840056aec70aa7e29763ce38827a0a16b66666cbdb4b4b4`, 686 lines / 88,436 bytes), preserved unchanged; r5 was not adopted and has no operative effect.  
**Audited lineage:** r4 (`52a550853291e6d4b0db8473e7f4c78185f1a7414d06ae207c5363b12e754429`, 603 / 68,813), r3 (`e45f1d898ab3d66016d2714472fd1023bcc93926a5e250a94bbb9072d9082d16`, 513 / 47,710), r2 (`c9bddd10608993ae73d86caa01a7adb6dd83220f98da8dab17e46c1aa3deb8ab`, 447 / 36,025), and the original recorded requirements (`ee38f0310042188a6821f8632e9e07bfa60dcac454a71ff2d972bb47439ac0e5`, 344 / 25,933), all preserved unchanged. The original remains in force under §2.1.  
**Review inputs:** Mimo's r4 audit `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r4-2026-11-07.md` (SHA-256 `eafaf6cc515fadaf80b34966ab92a8e5e5703db010f3e9bed1cc1517372aad41`, 258 lines / 30,900 bytes), Mimo's r3 audit, Mimo's r2 audit, and the 2026-11-07 pre-review assessment of r5. They are review evidence, not authority. The r5 pre-review is a document-level assessment, not an independent repository audit.  
**Pre-adoption gates:** Independent review of this r6 revision pending; owner adoption pending; **A-1 owner scope decision pending**.  
**Governing inputs:** Exact paths and hashes are recorded in §2.2; the adoption record must cite that table.  
**Applies to:** TTS Stage 3 corrective stages **C1–C8** and the eventual implementation.  
**Authority:** Additive engineering requirement proposal. It does not amend the governing contract, the proposed stage plan, the correction record, the owner-decision record, or any audit. It becomes binding only through the separate owner-adoption record in §2. Option B, if selected, requires a separate owner plan-amendment instrument outside this revision. It does not create owner authority; it records the owner as the stage acceptance authority in §7.2, subject to the separate-acceptor rule when the owner is the author, and limits budget review to conformance in §10.

This revision responds to the read-only findings F1–F10, the project-owner requirement to preserve normal-operation performance, Mimo's r2 audit findings A-1–A-11, Mimo's r3 audit findings R3-1–R3-12, Mimo's r4 audit findings R4-1–R4-11, and the three pre-review findings P1–P3. It is not an adoption, an authorization, a stage acceptance, or a decision on any pending owner question. It does not select an outcome for A-1.

## 1. Corrective basis

| Audit finding | Resolution in this revision |
|---|---|
| **F1 — authority** | §2 defines a separate, hash-bound owner-adoption mechanism and records the verified input hashes. It also states that the original recorded requirements remain in force until r6 adoption. |
| **F2 — C1 handshake contradiction** | C1 records the full contract start/restart map, not an `absent`-only rule. From `absent`, explicit TTS/synthesis, `ProviderCatalog`, or `ProviderStatus` may start. From `crashed`, the next explicit request (the contract also calls this the next explicit TTS request in §17.11.5), but not health/status/catalog observation, enters `backoff` and then `starting` only while the restart budget permits; exhaustion yields `unavailable`. From `backoff`, delay expiry enters `starting` with the same key. From `stopping`, only contract exit to `absent`; from `incompatible`, only the update-gated return to `absent`; neither permits a direct start. From `unavailable`, A2-D12 permits only an explicit same-identity `load` on a live helper whose helper-side state is `unloaded`; synthesis, catalog, voice-list, health, and status do not recover. `GetTTSVoiceCatalog` is never a start trigger. C1 owns the documented map; C5 owns executed restart/backoff; C7 owns executed start and reachability. See §8 C1, C5, C7. |
| **F3 — C6 overlap** | Replacement requires shutdown, termination if necessary, and observed `Wait()` before any new process starts; no concurrent old/new residency. See §8, C6. |
| **F4 — timeout/helper-exit races** | Timeout and unexpected helper exit are first-class TH-1.2 events and C4/C8 evidence. See §4.2, §8 C4, and §8 C8. |
| **F5 — 120 s idle grace** | The idle timer and helper may remain through the contract-defined 120 s grace; premature shutdown is not required. See §5.3 and §8 C5. |
| **F6 — app-wide I/O evidence** | §9 defines a mandatory I/O impact record for every stage and every touched application path. See R6-AC-7 and §8. |
| **F7 — measurable budgets** | §10 defines a resource-budget record; “bounded” requires a numeric value or explicit contract-provisional value, a measurement method, and an acceptance threshold. |
| **F8 — exception loophole** | §12 prohibits performance waivers of TH-1, HW-2, contract cleanup/lifecycle/data-integrity requirements, and process ordering. Any HW-1 budget exception requires a recorded owner decision with scope, expiry, and safety evidence. |
| **F9 — observability surface** | §11 makes bounded internal diagnostics the default. No renderer, public API, or status-shape change is authorized without separate contract/owner authority. |
| **F10 — C8 pass/fail** | C8 uses the contract, the original requirements, this revision, and, if the adoption record incorporates it, the proposed stage plan's evidence matrix. That matrix requires each row to name its environment and process boundary; r6 adds no CPU/GPU, throughput, latency, or workload claim. r6's lock, goroutine, process, handle, and temporary-file observations are additive TH-1 overlay checks, not claims about plan content. PR-1–PR-4 are not C8 acceptance conditions unless A-1 Option B is recorded; A-1 remains an open owner scope decision. See §6.4–§6.5. |
| **New requirement — performance preservation** | HW-1.8 retains normal-operation performance preservation. HW-2.4 defines the PR-1–PR-4 measurement framework, but assigns it no acceptance effect until A-1 is decided. |

No finding justifies abandoning the safety requirements. The predecessor remains the audited revision; this file is the proposed corrective revision for independent review.

### 1.1 Mimo r2 audit findings A-1–A-11

Mimo's r2 audit is review evidence, not authority. Its full artifact is `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r2-2026-11-07.md` (SHA-256 `567bc86c3d8b2a78b9687131d5937d77a9736393bed6a383dd0aba5ff948994b`, 221 lines / 29,130 bytes).

| Finding | Severity | r6 disposition |
|---|---|---|
| **A-1 — C8 performance-gate conflict** | Blocking | The PR-1–PR-4 C8 mandate is removed. A-1 is recorded as an open owner scope decision in §6.5; r6 does not choose between a separate performance gate and a formal C8 amendment. |
| **A-2 — all eight decisions required pending** | Major | §13 and §14 no longer make “all eight remain pending” an adoption condition. Adoption neither requires nor resolves any owner decision. |
| **A-3 — predecessor status unresolved** | Major | §2.1 states that the original recorded requirements remain in force until r6 adoption and that r6 supersedes them only where the texts differ after adoption. r2, r3, and r4 were audited but never adopted and have no operative effect. |
| **A-4 — no governing-input hashes** | Major | §2.2 records the exact r6 governing inputs and hashes, including r4 and its audit; the adoption record must cite that table. |
| **A-5 — handshake trigger unnamed** | Blocking for C1 | §8 C1 records the full state-scoped start/restart map: `absent` lazy start; `crashed` → `backoff` → `starting` only on the next explicit request and within budget; `backoff` delay expiry; `stopping`/`incompatible` only after the contract-defined return to `absent`; and `unavailable` no-inspection recovery with the A2-D12 explicit-load exception. `GetTTSVoiceCatalog` is never a trigger. C1 owns the map; C5 owns executed restart/backoff; C7 owns executed start/reachability. |
| **A-6 — inconsistent PR timing** | Major | §6.4 assigns PR-1–PR-4 no stage or C8 acceptance effect in r6 unless Option B is recorded. Stages record an R6-AC-8 impact statement only. |
| **A-7 — scope narrowed** | Minor | §3.1 restores “relies on … or is directly integrated with” and retains the anti-refactor guard. |
| **A-8 — exception can fail PR-3** | Minor | §12 states that no HW-1 exception may create a condition that fails a HW-2 pass condition; the normal-operation throttle prohibition is treated as part of HW-2.1 for that purpose. |
| **A-9 — PR-3 has no N/A or signal** | Minor | §6.4 adds a reviewed N/A rule for stages without adaptive control and requires a named contention signal and detection latency when controls exist. |
| **A-10 — R6-AC-7 N/A over-breadth** | Minor | R6-AC-7 requires N/A only for each in-scope path category not touched, by reference to the categories in §3.2. |
| **A-11 — circular approval field** | Minor | §10 defines the budget-conformance reviewer as a bounded, accountable function: appointment, independence, a duty to certify a conforming row, written reasons and exact unresolved conditions for a non-conforming determination, an escalation route for disputed source meaning, and replacement when no determination is issued. The reviewer cannot create a requirement by refusing to approve a conforming row and cannot accept a stage. §7.2 names the owner as the separate stage acceptance authority. |

### 1.2 Mimo r3 audit findings R3-1–R3-12

Mimo's r3 audit is review evidence, not authority. Its full artifact is `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r3-2026-11-07.md` (SHA-256 `c9804763e1cded31697a1ec05f3b9136a5dfab33dce5b5a4ebfb63f6035379a3`, 226 lines / 30,143 bytes). Its reported input hashes were rechecked in this workspace for this revision; its repository-inspection claims were not independently verified.

| Finding | Severity | r6 correction |
|---|---|---|
| **R3-1 — unbounded budget approver authority** | Blocking | §10 defines budget conformance review as a bounded, non-acceptance function. It may certify only limits within contract-fixed values, contract-provisional values with a named source, adopted requirements, and recorded owner decisions. It cannot approve a more permissive limit, resolve a pending decision, change scope, stage order, or architecture, or reject a stage. Withholding approval records an unresolved R6-AC-4 item. The two-helper-process cap is an explicit immutable maximum. Only the owner may authorize a requirement or scope change. |
| **R3-2 — §9.2 conjunctive status/catalog gate** | Blocking | §9.2 is disjunctive: a start is permitted when the contract permits it; an owner decision is required only where the contract is silent. The status/catalog lazy start is contract-owned under §17.7.2 and creates no new owner decision. §9.2 also distinguishes `ProviderCatalog` from `GetTTSVoiceCatalog`. |
| **R3-3 — design proposal called frozen/controlling** | Major | The design is a proposed stage plan and review input, not a source of requirements and not an adopted controlling instrument. Its stage criteria apply only where the adoption record expressly incorporates the plan under §2.1; otherwise they remain review inputs. No design text may narrow the contract, an original criterion, or this revision. A change to the plan is an owner decision recorded in the plan, not a contract amendment. |
| **R3-4 — AC identifier collision** | Major | Original AC-1–AC-7 remain textually and substantively unchanged. Every r6 overlay criterion is prefixed `R6-AC-`; bare `AC-n` always means the original. §7.3 reproduces the seven originals verbatim for verification, and §2.1 states that no overlay supersedes or narrows an original. |
| **R3-5 — state-blind start triggers** | Major | The start map is state-scoped and complete: `absent`; `crashed` → `backoff` → `starting`; `backoff` delay expiry; `stopping`/`incompatible` only after return to `absent`; and `unavailable` no-inspection recovery with the A2-D12 explicit-load exception. It does not prohibit the contract-authorized crash restart, does not authorize a direct start from `stopping` or `incompatible`, and does not treat status, health, or catalog observation as recovery from `unavailable`. |
| **R3-6 — C1/C7 evidence reallocation** | Major | C1 owns only the documented state-scoped start/restart map and structural no-frame/no-spawn checks. C5 owns executed crash-restart and backoff evidence. C7 owns executed start and reachability: status/catalog from `absent`, `ProviderCatalog` versus `GetTTSVoiceCatalog`, `unavailable` no-recovery, and explicit-load reachability. The proposed plan’s stage boundary is preserved; restart evidence is not moved to C7. |
| **R3-7 — wrong ProviderStatus/ProviderCatalog citation** | Minor | The trigger map cites both §17.7.2's lazy-start row and §17.12.1's `ProviderStatus`/`ProviderCatalog` rows. |
| **R3-8 — R6-AC-7 N/A cites no path categories** | Minor | R6-AC-7 points to the path categories enumerated in §3.2, not §9.1. |
| **R3-9 — C8 statements misdescribe the design** | Minor | F10 and the C8 row now distinguish the plan's actual evidence rule from r6's additive leak/handle/process checks. The unsupported design-content clause is removed. |
| **R3-10 — bare “catalog” ambiguity** | Minor | Every start rule distinguishes `ProviderCatalog` (may start/handshake/catalog, no load) from `GetTTSVoiceCatalog` (no helper start). |
| **R3-11 — unparseable approval sentence** | Minor | §10 replaces it with an explicit conformance certification: approval is not a measurement, does not substitute for the Measurement field, and does not authorize work beyond the approved limits. |
| **R3-12 — conditional native-provider authority** | Minor | Native-provider scope is not authorized. PR-1 records GPU/accelerator dimensions as N/A with review; no native-provider authorization class is implied. |

### 1.3 Mimo r4 audit findings R4-1–R4-11

Mimo's r4 audit is review evidence, not authority. Its full artifact is `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r4-2026-11-07.md` (SHA-256 `eafaf6cc515fadaf80b34966ab92a8e5e5703db010f3e9bed1cc1517372aad41`, 258 lines / 30,900 bytes). Its reported input hashes were rechecked in this workspace; its repository-inspection claims were not independently verified.

| Finding | Severity | r6 correction |
|---|---|---|
| **R4-1 — crash restart prohibited** | Blocking | §8 C1 replaces the `absent`-only rule with the complete contract map, including `crashed` → `backoff` → `starting` on the next explicit request within budget, and `backoff` delay expiry. C5 owns executed restart/backoff evidence; C7 owns executed start and reachability. |
| **R4-2 — unbounded reviewer withholding** | Major | §10 requires a written determination for every complete budget package, mandatory certification when it conforms, written element/source/evidence and cure conditions when it does not, an owner escalation route for disputed source meaning, and replacement if no determination is issued. The reviewer cannot create a requirement by withholding. |
| **R4-3 — no stage acceptance authority** | Major | §7.2 names the owner as the stage acceptance authority, acting after independent review and conformance findings. The author cannot accept; the reviewer cannot accept or reject. The matrix adds a separate acceptance-authority column. |
| **R4-4 — design role unresolved** | Major | §2.1 states one role: the design is a proposed execution plan, not a requirement source. Its criteria apply only if the adoption record incorporates the plan; otherwise they are review inputs. A plan change is an owner decision, not a contract amendment. |
| **R4-5 — original criteria narrowed** | Major | R6-AC-1, -2, -3, -5, -6, and -10 are non-normative pointers, not restatements. §7.3 reproduces AC-1–AC-7 verbatim. §2.1 states that the original governs on any difference and no overlay may narrow it. |
| **R4-6 — Option B adoption template contradiction** | Major | §2.2 makes the design non-amendment statement conditional; §2.3 adds design status and amendment-instrument fields; §2.5 records a drafting-time consistency check for A-1 Options A and B and for the open path. |
| **R4-7 — “likely N/A”** | Minor | C1 now requires an I/O impact record with explicit per-category N/A under §3.2. |
| **R4-8 — C1 start-test wording** | Minor | C1 owns no trigger-map start evidence; C7 owns executed start/reachability and C5 owns restart/backoff. C1's handshake tests, including real-process fake-helper handshakes, remain in C1 scope. |
| **R4-9 — C3 status-spawn ambiguity** | Minor | C3 now prohibits *unsolicited* spawn/load from inspection/status/startup and cross-refers the explicit `ProviderStatus` start to C7. |
| **R4-10 — O-1 omitted** | Minor | §13 records O-1 as a contract-deferred owner decision, distinguishes `ProviderCatalog` from `GetTTSVoiceCatalog`, and states that an O-1 answer authorizing picker startup is satisfied through `ProviderCatalog`. |
| **R4-11 — provisional-limit comparison undefined** | Minor | §10 makes the contract-provisional value the working ceiling until measured evidence is recorded; a more permissive value requires an express source range or a recorded owner decision. |

### 1.4 Pre-review findings P1–P3

The 2026-11-07 pre-review of r5 raised two authority issues and one lifecycle-wording issue. They are review evidence, not authority. r5 is preserved unchanged.

| Finding | Severity | Disposition in r6 |
|---|---|---|
| **P1 — Option B authority** | Blocking | The pre-review described C8's exclusion as contractual. The governing contract contains no C8 stage and no C8 performance clause; the C8 exclusion is in the proposed stage plan at C8 line 512. r5 correctly said the contract was not amended but did not identify C8 as plan-only or record the plan-amendment hash. r6 states that Option B is an owner amendment to the proposed stage plan under its amendment process, recorded in the adoption record with instrument, scope, and SHA-256. It is not a contract amendment. A contract-level C8 performance requirement would require a separate §17 contract amendment under A2-D11 and is not Option B. |
| **P2 — acceptance-authority exception** | Major | r5's matrix named the owner unconditionally while §7.2 requires an independent owner-designated acceptor when the owner is the author. r6 records the conditional in every matrix row: owner, or the §7.2 separate acceptor if the owner is the author. The separate acceptor performs the act; the owner's designation is not acceptance. |
| **P3 — restart trigger narrowed** | Major | r5 used a synthesis-or-load-only trigger. The contract says “next explicit request” (§17.7.2 line 1796) and “next explicit TTS request” (§17.11.5 line 2105). r6 restores that wording and keeps the plan's restriction that health/status/catalog observation does not restart. r6 does not narrow the contract's request class. |


## 2. Authority and owner adoption

### 2.1 Current status

This document is **PROPOSED — NOT BINDING**. Until the owner completes a separate adoption record, it may be used as a review checklist but cannot be used to accept or reject a stage, authorize implementation, or waive any contract requirement.

The original recorded requirements (`REQUIREMENTS-howl-thread-and-hardware-safety-2026-11-07.md`) remain in force, unamended, until this revision is adopted. Nothing in this revision suspends the original acceptance rule in the interim. r2, r3, r4, and r5 were audited or pre-review assessed but were not adopted and have no operative effect. Adoption of r6 supersedes the original only where the texts differ; where r6 is silent, the original remains in force. A pointer overlay (`R6-AC-1`, `-2`, `-3`, `-5`, `-6`, `-10`) does not supersede, narrow, or replace its original. The original text governs on any difference. No r6 overlay may waive or narrow an original criterion; only an explicit owner amendment to the original can change it.

The governing contract remains controlling. The design proposal is a **proposed execution plan**, not a source of requirements and not an adopted controlling instrument. Its stage criteria apply only to the extent the adoption record expressly incorporates the plan as the stage plan under §2.2; until then they are review inputs. No design text may narrow the contract, an original criterion, or this revision. A change to the plan or to stage allocation is an owner decision recorded through the plan’s amendment process, not a contract amendment. The correction record remains a review input. This document does not amend §17 or any adopted clarification. If a requirement conflicts with the contract, the conflict must be recorded and escalated; neither document may be silently overridden.

### 2.2 Separate adoption record required

A document cannot adopt itself, and inserting this file’s final hash into this file would change the hash. Adoption must therefore be recorded separately and must bind the exact artifact. The adoption record must contain:

- owner identity or accountable role, and the named stage acceptance authority under §7.2 if different from the owner;
- adoption date;
- exact path and SHA-256 of this file;
- adopted scope, including C1–C8 and the eventual implementation;
- the exact governing-input table below, with hashes rechecked at adoption;
- whether the proposed design is incorporated as the stage plan; if it is not, a statement that its criteria are review inputs only;
- confirmation that the governing contract and the correction record match the table and were not amended. The governing contract has no C8 stage or C8 performance clause; Option B does not amend it;
- the proposed stage plan's status: **not amended**, or **amended under §6.5 Option B**; for Option B, the record must include the plan-amendment instrument path, exact scope, and SHA-256, and design incorporation must be **Yes**;
- if the owner elects a contract-level C8 performance requirement, the separate §17 contract-amendment instrument and hash, recorded outside A-1 Option B; this revision does not authorize that amendment;
- the independent review record for r6 and its verified hash;
- a statement that the §7.1 requirements-authority matrix was checked and has no unresolved identifier collision or authority leakage, and that the four responsibilities—evidence production, conformance assessment, stage acceptance, and requirement/architecture change—are separated;
- a statement that the budget-conformance role in §10 is bounded to existing authority, that no owner authority has been delegated to it, and that the reviewer’s determination is not stage acceptance;
- the A-1 owner scope decision, if any performance gate is to have acceptance effect; if A-1 is unresolved, the adoption record must state that PR-1–PR-4 are inactive;
- any approved exceptions, with scope, duration, and safety evidence;
- an explicit statement that adoption does not authorize any stage or implementation.

Verified inputs at the time of drafting:

| Input | Path | SHA-256 | Lines | Bytes |
|---|---|---|---:|---:|
| Original recorded requirements | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-2026-11-07.md` | `ee38f0310042188a6821f8632e9e07bfa60dcac454a71ff2d972bb47439ac0e5` | 344 | 25,933 |
| Audited r2 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r2-2026-11-07.md` | `c9bddd10608993ae73d86caa01a7adb6dd83220f98da8dab17e46c1aa3deb8ab` | 447 | 36,025 |
| Mimo r2 audit | `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r2-2026-11-07.md` | `567bc86c3d8b2a78b9687131d5937d77a9736393bed6a383dd0aba5ff948994b` | 221 | 29,130 |
| Audited r3 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r3-2026-11-07.md` | `e45f1d898ab3d66016d2714472fd1023bcc93926a5e250a94bbb9072d9082d16` | 513 | 47,710 |
| Mimo r3 audit | `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r3-2026-11-07.md` | `c9804763e1cded31697a1ec05f3b9136a5dfab33dce5b5a4ebfb63f6035379a3` | 226 | 30,143 |
| Audited r4 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r4-2026-11-07.md` | `52a550853291e6d4b0db8473e7f4c78185f1a7414d06ae207c5363b12e754429` | 603 | 68,813 |
| Mimo r4 audit | `docs/review/AUDIT-howl-requirements-thread-hardware-safety-r4-2026-11-07.md` | `eafaf6cc515fadaf80b34966ab92a8e5e5703db010f3e9bed1cc1517372aad41` | 258 | 30,900 |
| Pre-review assessed r5 predecessor | `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r5-2026-11-07.md` | `93cfe1de5f362fba6840056aec70aa7e29763ce38827a0a16b66666cbdb4b4b4` | 686 | 88,436 |
| Governing contract | `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 | 182,904 |
| Proposed stage plan | `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151` | 711 | 69,168 |
| Correction record | `docs/review/CORRECTIVE-howl-tts-stage3-corrective-design-2026-11-07.md` | `8fe9356cc304f3135657d19a7f2ea6badd9cf2233585a03bd0bd5c4f6355e076` | 172 | 16,190 |
| Owner-decision instrument | `docs/review/OWNER-DECISIONS-howl-tts-stage3-2026-11-07.md` | `24d3344554e7ed6a465c8b7bdbca4461cda23f78353aa7626e3cf4ba7df122ba` | 610 | 45,891 |

A hash mismatch in any row reopens review. The table is not a substitute for the separate r6 hash: the adoption record must carry both.

### 2.3 Adoption block — left blank

**Owner:**  
**Stage acceptance authority (if different):**  
**Adoption date:**  
**Artifact path:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r6-2026-11-07.md`  
**Artifact SHA-256:**  
**Adopted scope:**  
**Governing-input table (§2.2) rechecked:** [ ] Yes  [ ] No  
**Requirements-authority matrix (§7.1) checked:** [ ] Yes  [ ] No  
**Independent r6 review record / hash:**  
**A-1 owner decision (if PR gate has acceptance effect):**  
**Design incorporation:** [ ] Incorporated as stage plan  [ ] Not incorporated; review input only  
**Proposed-stage-plan status:** [ ] Not amended  [ ] Amended under §6.5 Option B — instrument: / scope: / SHA-256:  
**Contract status:** [ ] Not amended (the contract has no C8 clause)  [ ] Separate §17 contract amendment recorded outside A-1 Option B — instrument/hash:  
**Approved exceptions:**  
**Decision:** [ ] Adopt  [ ] Reject  [ ] Return for correction

### 2.4 Pre-adoption rule

Owner adoption must not occur before the independent review of r6 is complete and recorded. A review finding that identifies a safety, authority, or contract contradiction reopens this revision rather than being treated as an adoption prerequisite waived by schedule. Mimo's r2, r3, and r4 audits are complete; none is the independent review of r6, and the r5 pre-review is not that independent review either.

A-1 must be decided before any PR-1–PR-4 result is given acceptance effect. If A-1 remains open, r6 may be considered only as a non-performance requirements revision, with PR-1–PR-4 expressly inactive; it cannot be adopted as a performance-gating instrument.

No adoption, review, or stage acceptance authorizes C1, C2, C3, C4, C5, C6, C7, C8, implementation, runtime execution, or benchmarking.

### 2.5 Adoption-path consistency check — A-1 Options A and B

This is a template test, not a selection of A-1. It verifies that the adoption record in §2.3 can express each permitted path without contradiction.

| A-1 path | Adoption-record entries | Consistency result |
|---|---|---|
| **Option A — separate performance gate** | A-1: separate gate; design incorporation: Yes/No as adopted; plan status: **Not amended**; contract: **Not amended**; C8: no PR condition; PR-1–PR-4: active only under the separate gate. | **Pass.** The plan non-amendment statement is true, the contract is not amended, and the gate is outside C8. |
| **Option B — formal C8 stage-scope amendment** | A-1: Option B; plan status: **Amended under §6.5 Option B — instrument/scope/SHA-256**; design incorporation: **Yes**; contract: **Not amended** (the contract has no C8 clause); C8: PR condition only as recorded in the plan amendment; adoption does not authorize measurement by itself. | **Pass.** The plan-amendment exception in §2.2 and the fields in §2.3 express the amendment. The plan-amendment instrument must exist before the record is complete. The contract is not amended because it has no C8 stage or C8 performance clause. |
| **A-1 open** | A-1: open; plan status: **Not amended** unless separately amended; contract: **Not amended**; PR-1–PR-4: inactive; C8: no PR condition; no performance claim. | **Pass.** This is the current drafting state; no measurement is authorized. |

Option B is not a contract amendment. The C8 exclusion is a proposed-stage-plan exclusion at C8 line 512, not a §17 contract clause. If Option B is selected, no adoption record may simultaneously say the plan was not amended, and a record missing the plan-amendment instrument, scope, or SHA-256 is incomplete. If the owner instead wants C8 performance to be a contract obligation, that requires a separate §17 contract amendment under A2-D11, recorded outside A-1 Option B; this revision does not authorize it. A-1 remains open at drafting, so no PR result has acceptance effect.
## 3. Scope

### 3.1 Stages and implementation

This revision applies to all eight corrective stages:

1. C1 — capability negotiation
2. C2 — status projection
3. C3 — lifecycle and load seam
4. C4 — exit and shutdown consistency
5. C5 — budget, backoff, and idle behavior
6. C6 — identity, recovery, and replacement
7. C7 — production adapter and explicit load
8. C8 — contract closure and runtime verification

It also applies to code changed by those stages and code those stages rely on for shared state, process ownership, scheduling, cleanup, or user-visible status, including code directly integrated with those paths.

“All code those stages rely on” is limited to code the stage changes, relies on, or directly integrates. It does not authorize an audit or retrofit of unrelated legacy code, third-party internals, or the entire application. Unrelated I/O debt discovered during the work is recorded separately under §9.4.

### 3.2 Application-wide I/O scope

HW-1 applies to every new or changed application path, not only the TTS helper. This includes settings, caches, model/asset preparation, status and health observation, renderer and picker paths, logging and diagnostics, helper lifecycle, and temporary/output files.

“Application-wide” means all paths the implementation changes, relies on, or directly integrates with. It does not mean that every pre-existing unrelated write in the application must be fixed by this program. The distinction must be recorded per stage in the I/O impact record (§9).

### 3.3 What this document does not do

This revision does not:

- amend the governing contract or its §17 rules;
- amend the recorded design proposal or correction record;
- amend the owner-decision record;
- resolve D10(a), D10(b), D1a, D9, D1b, D2, D5, or D3;
- decide A-1, either by creating a separate performance gate or by amending C8's scope;
- accept or authorize C1–C8;
- authorize implementation, build, test, launch, benchmark, runtime operation, commit, push, or deployment.

## 4. TH-1 — Thread safety

### 4.1 TH-1.1 Shared-state protection

All shared lifecycle, process, adapter, request, queue, cache, and status state must be protected against concurrent access and race conditions.

- Every shared mutable field must have an explicit synchronization or single-owner rule.
- Unsynchronized fast paths are not permitted merely because they are expected to be short.
- Snapshot and publication paths must not expose torn, partially updated, or stale-after-terminal state.
- A read path must not hold a lock across blocking I/O, process waits, channel operations, or unbounded work unless the design explicitly proves the blocking is required and bounded.
- Lock ordering and ownership must be documented where more than one lock can be held.

### 4.2 TH-1.2 Defined concurrency behavior

Cancellation, shutdown, timeout, load, unload, synthesis, and unexpected helper exit must each have defined concurrency behavior.

The design and tests must state, for every overlap that can occur:

- which operations may run concurrently;
- which operations are serialized, and at what boundary;
- which operation wins when requests race;
- whether a request is rejected, queued, superseded, or completed from an existing state;
- what terminal result is published, and that it is published at most once;
- how cancellation and timeout interact with startup, handshake, loading, active synthesis, and shutdown;
- how an unexpected `Wait()`/EOF exit interacts with an active request, cancellation, timeout, and cleanup;
- how a late terminal response is discarded or rejected after timeout, cancellation, or retirement;
- how a concurrent direct operation differs from the supervisor cancellation path where the contract distinguishes them.

A comment that a race is “unlikely” is not a concurrency specification.

### 4.3 TH-1.3 Liveness and containment

The implementation must have no known deadlock, unsafe lock ordering, uncontrolled background work, or orphaned workers or processes.

- No goroutine, timer, reader, writer, process, job object, temporary-file owner, or cleanup worker may start without a defined owner and termination path.
- Blocking operations must have bounded waits or an explicit cancellation path.
- Lock ordering must be acyclic and documented where more than one lock can be held.
- A failed operation must not leave a worker, reader, waiter, pipe, process, or temporary file orphaned.
- Shutdown must cause all owned background work to stop or be explicitly accounted for.
- A cancellation that cannot stop a worker must have a bounded escalation and cleanup path consistent with the contract.
- The contract-defined idle grace is not an orphan: a timer or helper waiting for the 120 s idle-close trigger is permitted until that trigger or an earlier authorized shutdown.

### 4.4 TH-1.4 Integrity under concurrency

Concurrent requests must not corrupt state or bypass resource limits.

- Admission and process-count checks must be atomic with respect to the resource they protect.
- Queue and concurrency bounds must not be bypassed by a second request, a retry, a status path, or a direct cancellation path.
- A stale completion must not overwrite a newer generation's state or terminal result.
- Duplicate terminal publication, lost terminal publication, send-on-closed-channel, double-close, and cleanup-before-kill ordering must be prevented or explicitly handled.
- Concurrent readers must observe a coherent snapshot, not an interleaving of old and new lifecycle fields.
- Timeout and unexpected exit must not race a success into acceptance after cancellation or retirement has won.

### 4.5 TH-1.5 Ownership and shutdown discipline

Shared resources must have one accountable owner and a deterministic release path.

- Process `Wait()` ownership must be single and explicit.
- Kill/wait/cleanup ordering must be deliberate and tested.
- A resource must not be released twice or by two independent owners.
- Shutdown must not race a new operation into starting work after the shutdown decision.
- Background work must be joined or proven terminated before the owner is released.

### 4.6 TH-1.6 Concurrency evidence

Concurrency claims require more than a successful functional run.

- A race detector or equivalent dynamic race evidence must be used where the language/toolchain supports it.
- Stress and overlap tests must exercise cancellation, shutdown, timeout, load/unload, synthesis, retries, replacement, and unexpected helper exit at the boundaries where races are possible.
- Tests must include timeout versus cancellation, timeout versus exit, cancellation versus exit, and late terminal versus retired request.
- Lock-order and ownership claims require static review in addition to tests.
- A test that exercises one request at a time does not establish thread safety.

### 4.7 TH-1.7 Applicability

Every stage must record a TH-1 applicability decision. “Not applicable” is permitted only when the stage has no shared state or concurrency surface, and the decision must be explicitly reviewed. A blanket “not applicable” statement is not evidence.

## 5. HW-1 — Hardware and resource safety

### 5.1 HW-1.1 Avoid unnecessary I/O

Avoid unnecessary disk writes and other major I/O across the application.

- Do not add write-on-read, write-per-poll, or background write cadence without a necessary product or durability reason.
- Do not reread the same model, asset, catalog, or configuration on every request when a bounded, valid cached view is sufficient.
- Do not flush logs or diagnostics per frame, per status poll, or per operation when a bounded state-change record is sufficient.
- Do not create, rewrite, or delete temporary/output files merely because an operation was observed rather than performed.
- “Major I/O” and “unnecessary” must be interpreted through the I/O impact record in §9, not by an undocumented reviewer intuition.

### 5.2 HW-1.2 Necessary writes are permitted

Model-loading writes and other genuinely necessary operations are permitted.

- A write is permitted when necessary to materialize or validate a model/asset, publish required output, preserve required settings, or maintain a bounded diagnostic record.
- Necessity must be justified at the call site and bounded in frequency and size.
- A necessary write must not become a repeated write on every request, retry, poll, or inspection.
- Where atomic publication is required by the contract or needed to prevent corruption, it must be used; atomicity is not unnecessary I/O.

### 5.3 HW-1.3 CPU, memory, and background activity

Avoid unnecessary memory allocations, CPU-intensive polling, redundant model discovery, and excessive background activity.

- Do not use a busy loop or tight polling loop in place of a timer, event, condition, or bounded wait.
- Do not perform full model/asset discovery when a cached, invalidated capability snapshot is sufficient.
- Do not allocate unbounded buffers, duplicate large payloads, or retain completed audio/model data without a defined lifetime.
- Do not leave completed background work running after its owner is finished.
- Do not add periodic work whose only purpose is to keep a path warm without an explicit product requirement and bounded cost.
- The contract-defined 120 s idle grace is an explicit exception: a helper or idle timer may remain until the idle-close trigger, and the timer must not be treated as an orphan or a violation of this rule.

### 5.4 HW-1.4 Bounded concurrency and consumption

Bound concurrency and resource consumption.

Every stage must declare and enforce limits for the resources it introduces or changes, including as applicable:

- worker/goroutine counts;
- queue depth and request concurrency;
- retry and restart counts;
- timer and background-task counts;
- process count;
- memory and buffer sizes;
- open files, pipes, handles, and job objects;
- temporary/output storage;
- log volume and retention.

“Unbounded” is not an acceptable default. A limit must be explicit, testable, and recorded in the resource-budget record in §10. A limit may use a contract-provisional value, but the source of that value must be named.

### 5.5 HW-1.5 Hardware-aware scheduling

Use hardware-aware scheduling and workload management so TTS does not monopolize CPU, memory, GPU, disk, or other resources.

- Concurrency must be appropriate to available CPU capacity and must leave headroom for the application and the user.
- A TTS workload must not assume a high-end machine, a dedicated GPU, or unlimited memory.
- Where platform scheduling or priority controls are used, they must be bounded, reversible, and must not degrade foreground responsiveness.
- Resource-intensive preparation or inference must not start from inspection, status, startup, or catalog paths unless the contract explicitly permits that start.
- The implementation must avoid a thundering herd of helpers, workers, timers, or retries when several requests arrive together.

### 5.6 HW-1.6 Cross-hardware safety

Prefer techniques that work safely across a range of hardware configurations, not just high-end machines.

- The design must state its assumptions about CPU cores, memory, disk, and GPU availability.
- Low-core and low-memory behavior must be considered, not only the developer workstation.
- A missing or weaker GPU must not turn into an unsafe retry loop, an unbounded fallback, or a silent correctness change.
- Degraded performance is acceptable only when behavior, cleanup, and lifecycle correctness remain intact.

### 5.7 HW-1.7 Observable contention and degradation

Resource contention and degradation must be observable and handled deliberately.

- The implementation must expose bounded, low-cost internal signals for queueing, rejection, retry/backoff, process count, and resource pressure where those conditions affect behavior.
- Observability must not violate HW-1.1 or HW-1.3; bounded counters and state changes are preferred over high-frequency logging.
- Degradation must have a defined response: wait within a bound, queue within a bound, reject with a contract-permitted result, or reduce nonessential work.
- Silent spinning, silent dropping, and unbounded retry are not deliberate handling.
- The observation surface follows §11: internal diagnostics by default; no public status, renderer, or API field is added without separate authority.

### 5.8 HW-1.8 Performance preservation during normal operation

Hardware protection must preserve available performance during normal operation.

- No arbitrary CPU or GPU utilization ceiling may be imposed during normal operation.
- No unnecessary throttling, sleep, delay, or concurrency reduction may be applied when no contention is measured.
- A contract-fixed process cap or a justified bounded concurrency limit is not an “arbitrary ceiling”; it must be recorded, justified, and must not be lowered further without contention evidence.
- Adaptive controls, where present, must activate in response to actual, measured contention, must be reversible, and must deactivate when contention clears. A stage with no adaptive control records that fact in its R6-AC-8 impact statement; if a performance gate is later authorized, PR-3 is N/A with review under §6.4.
- A static blanket throttle applied to normal operation is not an acceptable adaptive control.
- The overhead of safety instrumentation and verification must be measured separately from contention-induced degradation.
- Safety must never be sacrificed to improve benchmark results, and a benchmark result must never be used to justify weakening synchronization, cleanup, or lifecycle correctness.

## 6. HW-2 — Performance must not compromise safety

### 6.1 HW-2.1 Safety takes precedence

Never sacrifice synchronization, data integrity, process cleanup, or lifecycle correctness just to improve speed.

Performance work must not introduce:

- an unsynchronized fast path;
- a fire-and-forget worker without an owner or termination path;
- a detached process without kill/wait/cleanup ownership;
- an unbounded queue, retry, timer, or allocation;
- a skipped cleanup or best-effort release justified only by latency;
- a status or cache shortcut that can expose stale terminal state;
- a resource-limit bypass justified by expected load;
- an unmeasured permanent throttle justified only by “safety.”

### 6.2 HW-2.2 Deliberate degradation

Resource contention and degradation must be handled deliberately, not hidden.

- A slow path must remain a correct path.
- Backpressure, rejection, queueing, or reduced nonessential work must be explicit and bounded.
- A performance change must preserve the same terminal results, cleanup ordering, and ownership guarantees as the safe path.
- If the safe behavior cannot meet a performance target, the target is revised or the work is escalated; safety is not silently weakened.

### 6.3 HW-2.3 Performance evidence

A performance claim must include safety evidence.

- Before/after measurements must state the workload, hardware profile, and resource dimensions.
- A speedup measured only on a high-end machine does not establish cross-hardware safety.
- A performance optimization must include the concurrency, cleanup, and resource-limit tests affected by the change.
- A regression in safety evidence blocks acceptance regardless of the speed result.

### 6.4 HW-2.4 Performance-regression framework — acceptance assignment pending A-1

PR-1–PR-4 define a candidate performance-preservation measurement framework. Except for a recorded A-1 Option B amendment, they are not acceptance conditions for C8 or for any stage in this revision. The governing contract contains no C8 stage or C8 performance clause; C8's exclusion is in the proposed stage plan. Option B amends that plan, not the contract. The proposed stage plan's C8 scope is a review input unless the adoption record incorporates the plan; an Option B amendment changes the plan, not the requirements. r6 assigns PR evidence no acceptance effect because no performance claim or performance gate is authorized; the absence of that authority, not the proposal's status, controls. C8 does not make a performance claim, and no PR result may be used to accept or reject C8. The A-1 owner scope question is unresolved; §6.5 records the two options without selecting one. If Option B is selected, only the recorded plan amendment and adoption record can give PR-1–PR-4 C8 effect; otherwise this paragraph's no-effect default governs.

The framework remains available for a separately authorized performance-preservation gate if the owner selects that option. A stage that changes scheduling, concurrency, throttling, caching, or normal-operation work records an R6-AC-8 impact statement naming the PR dimensions its change affects. It does not run the full PR suite for stage acceptance. A stage may describe a PR subset in its R6-AC-8 impact statement, but it does not execute the suite or use results as acceptance evidence while A-1 is open.

| Check | Measurement | Pass condition |
|---|---|---|
| **PR-1 — Normal-operation overhead** | No injected contention; baseline versus candidate; latency (median and p95), throughput, CPU, memory, disk I/O, handles, and process count where applicable. Native-provider scope is not authorized; record GPU/accelerator dimensions as N/A with review. | Candidate does not exceed the separately authorized regression threshold; no arbitrary CPU/GPU ceiling or unnecessary throttle is active |
| **PR-2 — Contention degradation** | Same workload under a declared contention profile; baseline versus candidate | Degradation is bounded by the separately authorized threshold, and safety behavior remains correct |
| **PR-3 — Adaptive control** | Activation and deactivation under measured contention | Controls activate only after actual contention, are reversible, and deactivate when contention clears; no permanent normal-operation throttle. A stage with no adaptive control records PR-3 as N/A with review. Where controls exist, the contention signal must be named — queue depth, rejection count, or measured latency against the declared profile — and the detection latency must be stated so “after actual contention” is checkable. |
| **PR-4 — Safety regression** | Race, timeout, cancellation, exit, cleanup, and resource-limit checks under the same workload | Any safety regression blocks acceptance regardless of PR-1/PR-2 results |

If a separate performance gate is later authorized, the owner decision that authorizes the gate must set its numeric regression thresholds, contention profile, run count, and warm-up period. The budget conformance reviewer may only verify that the gate's recorded thresholds and profile fall within that owner decision; the reviewer may not create, widen, or waive a threshold. r6 does not invent these values. A report that presents measurements without a pass/fail comparison is not performance-gate evidence. No PR result may override PR-4 or any safety requirement.

### 6.5 A-1 open owner scope decision — performance gate

**Status: OPEN — OWNER DECISION REQUIRED. This section does not decide A-1.**

Mimo's r2 audit A-1 found that r2 made PR-1–PR-4 mandatory C8 acceptance evidence while the proposed stage plan excludes performance claims from C8. r6 removes that conflict by assigning no acceptance effect to PR-1–PR-4. It does not choose the owner-level scope question.

**Authority check.** The governing contract (`DESIGN-howl-tts-helper-architecture-2026-11-06.md`) contains no `C8` stage and no C8 performance clause; a case-sensitive scan returns no `C8` occurrence. C8 is stage 8 of the proposed stage plan, and its exclusion is at plan line 512: “No native provider, model download, licensing decision, or performance claim.” The contract's statement that “No value below is a performance measurement” (line 1844) governs provisional budget values, not a C8 clause. Accordingly, Option B does not amend the contract; it amends the proposed stage plan through the owner's plan-amendment process.

**Option A — separate performance-preservation gate (recommended for investigation, not selected).** PR-1–PR-4 run only under a separate, explicitly authorized gate with its own plan, owner-authorized thresholds, and acceptance record. C8 remains unchanged and excludes performance claims. The performance gate does not gate C8 contract closure. This preserves the plan boundary and gives the performance requirements a meaningful home.

**Option B — formal C8 stage-scope amendment to the proposed plan.** The owner explicitly authorizes an amendment to the proposed stage plan's C8 scope through the plan's amendment process and records in the adoption record the amendment instrument path, exact scope, SHA-256, and design incorporation. That amendment may make PR-1–PR-4 C8 acceptance conditions to the extent recorded. This is not a contract amendment and cannot be paired with an unqualified statement that the recorded proposed stage plan was not amended. If the owner instead wants C8 performance to be a contract obligation, that requires a separate §17 contract amendment under A2-D11, recorded outside A-1 Option B; this revision does not authorize it.

No third A-1 option is available: r6 must not leave the C8 exclusion unmentioned while giving PR evidence C8 effect. A contract-level C8 amendment is a separate instrument, not a third A-1 path.

Until the owner selects A or B, PR-1–PR-4 are inactive. They may be reviewed as a measurement design, but no PR result may accept, reject, block, or delay C8 or any stage. No performance measurement or benchmark is authorized.
## 7. Stage acceptance rules

The original recorded requirements define **AC-1–AC-7**. Those identifiers remain unchanged and independently operative. This revision adds overlay criteria with an explicit **R6-** prefix. A bare `AC-n` always means the original criterion; an `R6-AC-n` is an r6 overlay criterion. No unqualified `AC-n` in this revision redefines the original.

Every stage acceptance must include the applicable original AC-1–AC-7 and the applicable non-pointer r6 overlay. The pointer rows in §7.1 are traceability entries only; they add nothing and are satisfied by satisfying their originals. These rules are cumulative; a later stage does not excuse an earlier omission.

| Rule | Requirement |
|---|---|
| **R6-AC-1 — Concurrency model** | Non-normative pointer to original AC-1. It adds no requirement; the original text in §7.3 governs. |
| **R6-AC-2 — Lock and ownership review** | Non-normative pointer to original AC-2. It adds no requirement; the original text in §7.3 governs, including the static-review requirement. |
| **R6-AC-3 — Dynamic concurrency evidence** | Non-normative pointer to original AC-3. It adds no requirement; the original text in §7.3 governs, including repeatability. |
| **R6-AC-4 — Resource budget** | The stage records every applicable limit in the §10 budget record, with a numeric value or contract-provisional value, an authorized source, measurement method, and acceptance threshold. This expands original AC-4; it does not replace it. The §10 conformance determination is not stage acceptance. |
| **R6-AC-5 — Safety under load** | Non-normative pointer to original AC-5. It adds no requirement; the original text in §7.3 governs. |
| **R6-AC-6 — No hidden background work** | Non-normative pointer to original AC-6. It adds no requirement; the original text in §7.3 governs. The contract-defined idle grace is a bounded timer with an owner and termination path, not an exemption from AC-6. |
| **R6-AC-7 — Application I/O evidence** | Every stage records its I/O impact under §9, including explicit N/A for each in-scope path category it did not touch, per §3.2. |
| **R6-AC-8 — Performance preservation** | The stage records whether its changes affect normal-operation performance, naming the PR dimensions affected. C8 does not apply PR-1–PR-4 unless A-1 Option B is recorded; the performance gate otherwise remains unassigned pending A-1. |
| **R6-AC-9 — Applicability review** | A TH-1 or HW-1 “not applicable” determination requires an explicit reviewed entry; blanket N/A is not accepted. |
| **R6-AC-10 — No stage substitution** | Non-normative pointer to original AC-7. It adds no requirement; the original text in §7.3 governs. |

A stage that has functional tests but no applicable TH/HW evidence is not accepted. A stage with no declared resource budget is not accepted. A stage with no recorded acceptance decision by the §7.2 authority is not accepted.

### 7.1 Requirements-authority matrix

This matrix separates four responsibilities: **evidence production**, **conformance assessment**, **stage acceptance**, and **requirement/architecture change**. The evidence owner is the change author. The conformance assessor determines whether recorded evidence meets a criterion but does not accept or reject the stage, amend the criterion, or authorize a change. The stage acceptance authority is the owner under §7.2, **except** when the owner is also the change author: then it is the separate accountable acceptor designated by the owner under §7.2 and independent of the work. The matrix records this conditional in every row; the owner is the default actor, not an unconditional one. Requirement and architecture change authority remains with the owner through the appropriate recorded instrument. No conformance assessor may be the change author or the owner acting on a scope question. The budget-conformance reviewer's role is limited by §10. Stage-specific identifiers below are evidence-allocation labels; they do not create new compliance categories.

| Current identifier | Originating authority | Evidence owner | Conformance assessor | Stage acceptance authority | Requirement/architecture change authority |
|---|---|---|---|---|---|
| **AC-1 — Concurrency model (original)** | Original requirements AC-1 | Every stage C1–C8 | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original text may be changed only by recorded owner amendment |
| **AC-2 — Lock and ownership review (original)** | Original requirements AC-2 | Every stage C1–C8 | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original text may be changed only by recorded owner amendment |
| **AC-3 — Dynamic concurrency evidence (original)** | Original requirements AC-3 | Every stage C1–C8 | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original text may be changed only by recorded owner amendment |
| **AC-4 — Resource budget (original)** | Original requirements AC-4 | Every stage C1–C8 | Independent stage reviewer for conformance; budget conformance certification is separately bounded by §10 | Owner; if owner is author, §7.2 separate acceptor | Owner; original text may be changed only by recorded owner amendment |
| **AC-5 — Safety under load (original)** | Original requirements AC-5 | Every stage C1–C8 | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original text may be changed only by recorded owner amendment |
| **AC-6 — No hidden background work (original)** | Original requirements AC-6 | Every stage C1–C8 | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original text may be changed only by recorded owner amendment |
| **AC-7 — No stage substitution (original)** | Original requirements AC-7 | C8 only; cannot substitute for C1–C7 | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original text may be changed only by recorded owner amendment |
| **R6-AC-1 — Concurrency model (pointer)** | r6 §7; non-normative pointer to original AC-1; adds nothing | Every stage C1–C8 (via original AC-1) | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original AC-1 may be changed only by recorded owner amendment |
| **R6-AC-2 — Lock and ownership review (pointer)** | r6 §7; non-normative pointer to original AC-2; adds nothing | Every stage C1–C8 (via original AC-2) | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original AC-2 may be changed only by recorded owner amendment |
| **R6-AC-3 — Dynamic concurrency evidence (pointer)** | r6 §7; non-normative pointer to original AC-3; adds nothing | Every stage C1–C8 (via original AC-3) | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original AC-3 may be changed only by recorded owner amendment |
| **R6-AC-4 — Resource budget (expansion)** | r6 §7 and §10; expands original AC-4; original remains operative | Every stage C1–C8 | Budget conformance reviewer under §10; this is not stage acceptance | Owner; if owner is author, §7.2 separate acceptor | Owner; original AC-4 may be changed only by recorded owner amendment; budget exception under §12 |
| **R6-AC-5 — Safety under load (pointer)** | r6 §7; non-normative pointer to original AC-5; adds nothing | Every stage C1–C8 (via original AC-5) | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original AC-5 may be changed only by recorded owner amendment |
| **R6-AC-6 — No hidden background work (pointer)** | r6 §7; non-normative pointer to original AC-6; adds nothing | Every stage C1–C8 (via original AC-6) | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original AC-6 may be changed only by recorded owner amendment |
| **R6-AC-7 — Application I/O evidence (overlay)** | r6 §7 and §9; new overlay | Every stage C1–C8 | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner, by r6 adoption or recorded amendment |
| **R6-AC-8 — Performance preservation (overlay)** | r6 §5.8 and §6.4; new overlay | Every stage C1–C8; C8 carries no PR condition absent a recorded Option B plan amendment | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; A-1 owner decision required for any performance gate |
| **R6-AC-9 — Applicability review (overlay)** | r6 §4.7 and §5.7; new overlay | Every stage C1–C8 | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner, by r6 adoption or recorded amendment |
| **R6-AC-10 — No stage substitution (pointer)** | r6 §7; non-normative pointer to original AC-7; adds nothing | C8 only; cannot substitute for C1–C7 (via original AC-7) | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; original AC-7 may be changed only by recorded owner amendment |
| **C1-CAPABILITY** | r6 §8 C1; proposed stage plan (if incorporated under §2.1) C1 | C1 only; capability snapshot and classification | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; stage-allocation change requires owner decision and, if incorporated, plan amendment |
| **C1-START-MAP** | r6 §8 C1; proposed stage plan (if incorporated under §2.1) C1 | C1 only; documented assertion, no process start | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; stage-allocation change requires owner decision and, if incorporated, plan amendment |
| **C2-STATUS** | r6 §8 C2; proposed stage plan (if incorporated under §2.1) C2 | C2 only | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; stage-allocation change requires owner decision and, if incorporated, plan amendment |
| **C3-LIFECYCLE** | r6 §8 C3; proposed stage plan (if incorporated under §2.1) C3 | C3 only | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; stage-allocation change requires owner decision and, if incorporated, plan amendment |
| **C4-EXIT** | r6 §8 C4; proposed stage plan (if incorporated under §2.1) C4 | C4 only | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; stage-allocation change requires owner decision and, if incorporated, plan amendment |
| **C5-BUDGET** | r6 §8 C5; proposed stage plan (if incorporated under §2.1) C5 | C5 only; budget, restart/backoff, and concurrency evidence | Independent stage reviewer for stage conformance; budget conformance reviewer under §10 for budget rows | Owner; if owner is author, §7.2 separate acceptor; after budget conformance determination | Owner; reviewer may only certify conformance; owner alone may approve an exception or change a limit's authority |
| **C6-IDENTITY** | r6 §8 C6; proposed stage plan (if incorporated under §2.1) C6 | C6 only | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; stage-allocation change requires owner decision and, if incorporated, plan amendment |
| **C7-ADAPTER** | r6 §8 C7; proposed stage plan (if incorporated under §2.1) C7 | C7 only; adapter concurrency and default-off | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; stage-allocation change requires owner decision and, if incorporated, plan amendment |
| **C7-START-EVIDENCE** | r6 §8 C7; proposed stage plan (if incorporated under §2.1) C7; contract §17.7.2, §17.12.1, A2-D12 | C7 only; executed start and reachability | Independent stage reviewer; contract conformance only | Owner; if owner is author, §7.2 separate acceptor; after independent start/reachability review | Owner for requirements/scope; contract changes only through the contract process; no reviewer authority |
| **C8-RUNTIME** | r6 §8 C8; proposed stage plan (if incorporated under §2.1) C8 | C8 only; runtime evidence matrix | Independent stage reviewer; not the author | Owner; if owner is author, §7.2 separate acceptor | Owner; no performance claim or PR condition without A-1 |

Bare `AC-7` is always the original no-substitution rule. The I/O rule is `R6-AC-7`; a stage report that uses bare `AC-7` for I/O evidence is ambiguous and is not accepted as an r6 citation. A pointer row does not create an additional requirement and does not supersede its original. The matrix does not create a new stage owner or a new change authority.

### 7.2 Stage acceptance authority

Stage acceptance is an **owner act**. The owner is the stage acceptance authority for this program. This revision records that existing role; it does not create owner authority and does not transfer it to any reviewer.

The owner is the default stage acceptance authority. The matrix also records the §7.2 exception: when the owner is the change author, the owner-designated independent accountable acceptor performs the acceptance act. The owner’s designation is not the acceptance, and the author cannot accept their own stage. The applicable acceptance authority may accept a stage only when all of the following are recorded:

1. a complete stage evidence package, including the applicable original AC evidence and the non-pointer r6 overlay evidence;
2. an independent stage review and, for every budget row, a written §10 conformance determination;
3. all applicable criteria are satisfied, or a §12 exception is recorded within its scope;
4. no unresolved contract, safety, authority, or identifier contradiction remains; and
5. the acceptance record names the stage, artifact hashes, review record, decision, date, and any unresolved item disposition.

The evidence producer or change author may not accept their own stage. If the owner is also the change author, a separate accountable acceptor designated by the owner and independent of the work must accept; that acceptor is the acceptance authority for that row, not the owner acting as author. If no such acceptor exists, the stage remains unaccepted. A conformance assessor’s finding is not acceptance and is not rejection. The owner’s acceptance is not a conformance finding and cannot cure missing evidence. The owner may reject or return a stage for correction. Acceptance does not authorize implementation, runtime execution, or deployment; those require separate authorization.

### 7.3 Original criteria preservation

The seven original criteria remain textually and substantively operative. This table reproduces them for verification; the original file remains authoritative if any transcription differs. No r6 overlay, pointer, or stage table may narrow, waive, or replace them.

| Identifier | Original text (verbatim) |
|---|---|
| **AC-1 — Concurrency model** | The stage documents the shared state it touches, the operations that can overlap, the serialization boundary, and the winner/terminal rule for each relevant race. |
| **AC-2 — Lock and ownership review** | The stage records lock ordering, single-owner resources, blocking-under-lock exclusions, and background-work ownership. Static review is required where tests cannot prove the property. |
| **AC-3 — Dynamic concurrency evidence** | The stage runs race/stress/overlap tests for its affected operations, including cancellation, shutdown, load/unload, retry, and replacement where applicable. Evidence must be repeatable, not a single happy-path run. |
| **AC-4 — Resource budget** | The stage declares bounded limits for its workers, queues, retries, timers, processes, buffers, handles, I/O, and storage, and records the measurements or checks that support them. |
| **AC-5 — Safety under load** | The stage demonstrates that concurrent requests and degraded resources do not corrupt state, bypass limits, leak workers/processes, or lose/duplicate terminal results. |
| **AC-6 — No hidden background work** | Every started goroutine, timer, reader, process, and cleanup task has an owner, a bound, and a termination path. |
| **AC-7 — No stage substitution** | C8's final runtime verification may aggregate evidence, but it cannot substitute for the stage-level evidence C1–C7 must produce. |

## 8. C1–C8 acceptance overlay

The following table supplements the governing contract, the original requirements, and—only where the adoption record incorporates it under §2.1—the proposed stage plan. It does not replace any of them. The D10(a)/D10(b), D1a, D9, D1b, D2, D5, and D3 gates remain in force. It adds no performance acceptance condition by default: C8 makes no performance claim, and PR-1–PR-4 are inactive pending A-1. If A-1 Option B is recorded, C8's PR conditions are only those in the recorded plan amendment. C1 owns the documented start/restart map; C5 owns executed restart/backoff; C7 owns executed start and reachability. The plan’s stage allocation is preserved; the contract-authorized crash restart is not moved out of C5.

| Stage | TH-1 acceptance focus | HW-1 / HW-2 / performance focus | Required evidence |
|---|---|---|---|
| **C1 — Capability negotiation** | Capability snapshot publication is atomic and single-writer; concurrent negotiation, gating, and reads do not observe torn or stale fields; handshake completion and failure have one terminal outcome. | **State-scoped start/restart map.** From `absent`: explicit TTS/synthesis, `ProviderCatalog`, or `ProviderStatus` under §17.7.2. From `crashed`: the next explicit request (the contract also calls this the next explicit TTS request in §17.11.5)—not health, status, or catalog observation—enters `backoff`, then `starting` only while the restart budget permits; exhaustion yields `unavailable`. From `backoff`: delay expiry enters `starting` with the same key. From `stopping`: no direct start; only exit 0 to `absent`. From `incompatible`: no direct start; only the update-gated return to `absent`. From `unavailable`: no synthesis/catalog/voice-list/health/status recovery; only explicit same-identity `load` on a live helper whose helper-side state is `unloaded`, subject to D1a/D1b. `GetTTSVoiceCatalog` is never a trigger. No new process starts from any other live state; a `load` transition in an already-live helper is not a process start. No model load or disk write from negotiation or inspection. Parsing and descriptor storage are bounded. | A documented state-scoped trigger/start/restart map covering `absent`, `crashed`, `backoff`, `stopping`, `incompatible`, and `unavailable`; structural no-frame/no-spawn tests for direct `unload`, direct `Cancel`, and unknown-capability inspection; no-model-load/no-write evidence; allocation/limit review; I/O impact record with explicit per-category N/A under §3.2. C1 executes C1 handshake tests, including real-process fake-helper handshakes (using the plan's C1 tests when incorporated), but owns no trigger-map start evidence. C5 owns executed restart/backoff; C7 owns executed start and reachability. A map that starts from `unavailable` except the explicit-load case, treats `GetTTSVoiceCatalog` as a trigger, starts from `stopping`/`incompatible` before `absent`, restarts from `crashed` without an explicit request or budget, or omits the `crashed`/`backoff` edge fails C1. |
| **C2 — Status projection** | Multi-field status is a coherent snapshot; lifecycle transitions cannot produce mixed old/new fields; status reads do not hold locks across I/O or block transitions indefinitely. | Status is bounded-cost and performs no model discovery, file read, or high-frequency log write; buffers and counters are bounded; observability is internal §11 diagnostics only. | Concurrent status-vs-transition tests; no-torn-snapshot assertions; blocking-call review; I/O impact record; evidence that no public status/renderer/API field was added. |
| **C3 — Lifecycle and load seam** | Load/unload/cancel/synthesis have defined overlap and winner rules; one active load; lock order between adapter, client, and lifecycle state is documented; no orphaned worker on failure. | No **unsolicited** spawn or model load from inspection/status/startup. An explicit `ProviderStatus` request from `absent` is the C7 lazy-start evidence, not a C3 eager start. No duplicate load; waits are bounded; no write cadence added by lifecycle bookkeeping. | Overlap tests for load/unload/cancel; failure-path ownership checks; lock-order review; duplicate-load and orphan checks; I/O impact record; evidence that the only status-triggered start is the explicit C7 path. |
| **C4 — Exit and shutdown consistency** | Timeout, cancel, shutdown, kill, and wait races have one terminal result; unexpected helper exit is authoritative; single `Wait()` ownership holds; no double-close, send-on-closed, or cleanup-before-kill; shutdown prevents new work. | Grace and shutdown waits are bounded; no busy wait; no unnecessary file rewrite during cleanup; process and handle release is prompt. | Race tests for timeout/cancel/exit/shutdown vs active work; single-wait and double-terminal assertions; goroutine/process/handle leak checks; cleanup-ordering review; I/O impact record. |
| **C5 — Budget, backoff, and idle** | Restart ledger, backoff timers, idle deadline, and new requests have defined races; no lost wakeup, double close, timer leak, or deadlock between timer and request paths. | Timers/events replace CPU polling; restart count and backoff are bounded; no disk churn per tick; idle close occurs after the contract-defined 120 s grace and actually releases resources; the idle timer is not treated as an orphan. The restart path is contract-authorized: `crashed` plus the next explicit request enters `backoff`; 1/2/5 s delay then `starting` with the same key; maximum three consecutive restarts per key in 10 minutes; a fourth attempt exhausts to `unavailable`; reset after 10 minutes of stable readiness. No health/status/catalog observation restarts a helper. | Timer-vs-request overlap tests; **executed** `crashed` → `backoff` → `starting` evidence on an explicit request; 1/2/5 backoff checks; three-restart/10-minute budget and fourth-attempt exhaustion; stable-readiness reset; no health/status/catalog restart; idle-grace and idle-close resource check; race/stress evidence; I/O impact record. |
| **C6 — Identity, recovery, and replacement** | Identity transitions are serialized; replacement vs pending requests has a defined winner; stale completions cannot overwrite the new generation; cap admission cannot be bypassed. | The old process is shut down, killed if necessary, and `Wait()` is observed before any new process starts. No concurrent old/new process or model residency. No redundant model-load writes; replacement is bounded. Replacement is not a crash restart and is exempt from the crash-restart budget under the contract. | Concurrent identity/replacement tests; no-overlap and process-count assertions; stale-generation checks; resource release review; evidence that replacement is not counted as a crash restart; I/O impact record. |
| **C7 — Production adapter and explicit load** | Adapter methods are safe under concurrent calls; direct cancel/unload vs synthesis is defined; cap admission is atomic; no path bypasses the limit; explicit load cannot race shutdown. | Adapter is default-off; no startup load; catalog/discovery work is cached and bounded; no unnecessary renderer/status/settings/logging/disk I/O; queues and helpers are bounded. Executed start behavior is state-scoped: from `absent`, explicit `ProviderStatus`/`ProviderCatalog` may start/handshake/catalog and must not load; `GetTTSVoiceCatalog` must not start. From `unavailable`, no status/catalog/voice-list/health/synthesis request recovers; only explicit `load` on a live helper whose helper-side state is `unloaded` recovers, and only where D1a/D1b permit. The helper child-process cap is two; a third is refused before start, with the D3 error surface unresolved. C7 does not duplicate C5’s crash-restart execution evidence. | Executed start and reachability evidence: status/catalog from `absent`; `ProviderCatalog` versus `GetTTSVoiceCatalog`; no recovery from `unavailable`; explicit load only if D1a/D1b permit; actual child-process count capped at two; no load from status. Adapter concurrency tests; cap-admission race tests; default-off and no-eager-load checks; resource budget and I/O impact review covering every touched application path; caller-map test; cross-reference to C5 for `crashed`/`backoff` execution. |
| **C8 — Runtime verification** | Race/stress/soak coverage for cancellation, shutdown, timeout, load/unload, retry, replacement, unexpected exit, and cap admission; lock-order, goroutine, process, handle, and temporary-file leak checks; clean repeated runs. | C8 uses the contract, the original requirements, this revision, and—if incorporated—the proposed plan’s evidence matrix. The design rule that every row names its environment and process boundary applies. No performance claim, throughput/latency claim, GPU/accelerator claim, or PR-1–PR-4 condition is added absent a recorded Option B plan amendment. Lock-order, goroutine, process, handle, and temporary-file observations are additive TH-1 overlay checks, not claims that the plan already required them. | Pre-execution C8 evidence plan; commands, exit codes, logs, and hashes; every row names its environment and process boundary; Windows-specific rows close only through a Windows real-child run; no performance claim; PR results are not C8 evidence unless a recorded Option B plan amendment makes them so. |

## 9. Application-wide I/O evidence

### 9.1 Mandatory stage record

Every stage report must include an I/O impact record for every application path it changes, relies on, or directly integrates with. A stage with no such path must record an explicit N/A and have it reviewed; silence is not N/A.

Each record must state:

- path or component;
- operation (read, write, delete, scan, discover, log, or other);
- trigger and cadence;
- necessity and product/durability reason;
- bound on size, frequency, and retention;
- whether the write is model-loading, output publication, settings, cache, logging, or temporary-file handling;
- measurement or design proof;
- N/A rationale and reviewer, where applicable.

### 9.2 Application rules

- **Settings:** do not rewrite the settings file on read, status refresh, or an unchanged apply.
- **Caches:** do not rescan or rewrite a valid cache on every request; invalidate deliberately.
- **Logging:** do not emit per-frame, per-poll, or per-operation disk logs when a bounded state-change record is sufficient.
- **Status/health:** do not perform model discovery or file reads merely to answer a status request. An explicit status request from `absent` is the §17.7.2 lazy start; `ProviderStatus` may start/handshake but must not load. From `unavailable`, status and health do not recover a process. From `crashed` or `backoff`, status, health, and catalog do not restart; only the next explicit request under §17.7.2 and C5 does.
- **Renderer/catalog:** do not start a helper or load a model from an inspection path unless the contract permits it; an owner decision is required only where the contract is silent. The permitted adapter trigger is `ProviderCatalog` (contract §17.12.1); the renderer-facing `GetTTSVoiceCatalog` must not start a helper. Neither path loads a model.
- **Output/temporary files:** write only what is required, publish atomically where needed, and clean up deterministically.
- **Model loading:** materialization and validation writes are permitted when necessary; they must not become a per-request or per-retry rewrite.

### 9.3 Measurement

Where a stage introduces or changes I/O, the record must state how the frequency and size were measured or bounded. “No unnecessary I/O” is not demonstrated by an unmeasured assertion. If a separate performance gate is later authorized, it must include the I/O dimensions in PR-1/PR-2 where they are relevant. C8 itself makes no performance claim and runs no PR check.

### 9.4 Unrelated legacy I/O

An unrelated existing I/O problem discovered during this work is recorded as a separate finding. It does not expand the current stage and is not silently fixed outside scope. A stage may reference it but cannot use it as an excuse for adding new I/O.

## 10. Resource-budget record

Every stage must record each applicable resource in a table with these fields:

| Field | Requirement |
|---|---|
| Resource | Worker, queue, retry, timer, process, memory, buffer, file, handle, job object, storage, or log volume |
| Limit | A numeric value, or an explicit contract-provisional value with its source. The value may not be more permissive than the authorized source. For a contract-provisional source, the provisional value is the working ceiling until measured evidence is recorded and accepted; a longer timeout, larger buffer, or otherwise more permissive value requires a source that expressly authorizes that range or a recorded owner decision. |
| Authorized source | The exact contract section, adopted requirement, recorded owner decision, or contract-provisional source that fixes or authorizes the limit. A bare “design choice” is not an authorized source. |
| Basis | Why this limit is safe for the workload and hardware configuration. |
| Measurement | Test, trace, count, or design proof used to verify it. |
| Acceptance threshold | The pass/fail boundary for this stage. It may not create or widen a requirement. |
| Conformance determination | The budget conformance reviewer must issue a written determination that this row is within the authorized source and that its method and threshold are conforming, or must state the exact non-conforming element and the conditions to cure. This is not stage acceptance or rejection; the change author cannot determine conformance. The owner alone may approve an exception under §12. |

No “bounded” claim is accepted without a limit and a measurement. A limit may be provisional only where the contract labels it provisional; the source must be named. An arbitrarily high limit is not evidence of safety. The limit must not impose an arbitrary normal-operation CPU/GPU ceiling; justified concurrency and process caps remain permitted under HW-1.8.

A limit fixed by the contract, an adopted requirement, or a recorded owner decision is immutable at this gate. A budget may set a lower limit only where HW-1.8 permits and evidence supports it; it may not set a higher or more permissive limit. A missing, conflicting, ambiguous, or decision-dependent limit is an unresolved R6-AC-4 item, not reviewer discretion, and must be escalated to the owner. The contract's v1 hard cap of two helper processes is the concrete example: a proposed limit of three is not a budget within authority; it is a contract conflict. Two is the maximum, and a lower cap requires the evidence HW-1.8 requires.

**Budget conformance reviewer — role and appointment.** The role is created by this revision, not by the contract. The owner appoints a reviewer in the adoption record or the stage review charter, naming the person or accountable role, the appointment date, the independence attestation, the review scope, and the review period. The reviewer must be independent of the change author and must not be the owner acting on a scope question. The author cannot serve as reviewer. Appointment does not transfer owner authority, and a reviewer’s determination is not stage acceptance.

**Duty to determine.** For every complete budget package, the reviewer must issue a written determination within the review period fixed before the review begins:

- if the package conforms, the determination **must** be “conforms”; the reviewer may not withhold, delay, or condition a conforming determination on a change to scope, stage order, architecture, or a requirement; or
- if it does not conform, the determination must be “does not conform” and must name the exact limit, authorized source, missing or conflicting evidence, the criterion it fails, and the conditions required to cure it.

A failure to issue a determination within the review period is not a finding of non-conformance and is not a veto. The owner may appoint a replacement reviewer, who must issue a determination; the stage remains unaccepted until a determination or a recorded owner disposition exists.

**Scope of review.** The reviewer may only:

- compare each limit, authorized source, measurement method, and acceptance threshold with the plain text of the contract, an adopted requirement, a recorded owner decision, or a contract-provisional source;
- record a limit that is missing, conflicting, ambiguous, or dependent on a pending owner decision as an unresolved R6-AC-4 item;
- issue the written conformance determination described above.

The reviewer may not:

- approve a limit that is more permissive than the contract, an adopted requirement, or a recorded owner decision;
- create, amend, waive, or impose a new requirement, scope, stage order, contract value, or architecture;
- choose among competing readings of an ambiguous source; the reviewer must record a contract-interpretation question and escalate it;
- resolve or pre-empt a pending owner decision;
- condition approval on a change to scope, stage order, or architecture;
- approve a performance-gate threshold, profile, run count, or warm-up period that the owner has not authorized;
- accept or reject a stage, or substitute for the stage's acceptance authority;
- manufacture a new requirement by refusing to approve a conforming row.

**Escalation and disposition.** A disagreement about what a contract or source value means is a contract-interpretation question. The reviewer records the competing readings and the evidence and escalates to the owner; the reviewer does not resolve it. The owner may record an interpretation as a project decision, determine that the row conforms or does not conform, or require correction. An owner interpretation is not a contract amendment; if it would change the contract, the contract process must be used. The owner may not direct a “conforms” finding where evidence is missing or a limit is more permissive than its authorized source. A pending question remains an unresolved R6-AC-4 item; the stage is not accepted until it is dispositioned.

Approval certifies that the recorded limits and thresholds meet §10's conformance conditions. It is not a measurement, does not substitute for the Measurement field, and does not authorize work beyond the approved limits. The owner remains the only authority for requirement, scope, contract-value, architecture, exception, and stage-acceptance decisions, exercised through the appropriate recorded instrument. Stage acceptance is separately governed by §7.2.

If a separate performance gate is later authorized, the owner decision that authorizes it must set its thresholds, contention profile, run count, and warm-up period. The budget conformance reviewer may only verify that those values fall within that owner decision. r6 does not invent them, and no reviewer may create or widen them.

## 11. Observability

- Bounded internal counters, state changes, and diagnostics are the default observation surface.
- Do not emit per-frame, per-poll, or per-operation disk logs; prefer bounded state changes.
- Record queue depth, rejection, retry/backoff, process count, and resource-pressure indicators where those conditions affect behavior.
- A degradation response must be one of: wait within a bound, queue within a bound, reject with a contract-permitted result, or reduce nonessential work.
- Observability must not hold lifecycle locks across I/O or allocate unbounded per-request history.
- A health/status response must not become a hidden model-discovery or disk-scan operation.
- No renderer, public API, Wails binding, or status-shape change is authorized by this section. Any external exposure requires a separate contract/owner decision.

## 12. Conflict, exceptions, and waivers

- **Contract conflict:** if a TH/HW requirement cannot be met without violating the governing contract, stop and record the conflict. Do not silently weaken either document.
- **Stricter requirement:** a stricter TH/HW rule may be applied when compatible with the contract.
- **No safety waiver for performance:** TH-1, HW-2, synchronization, data integrity, process cleanup, lifecycle correctness, termination ordering, and contract-required limits may not be waived to improve latency, throughput, or benchmark results.
- **Budget exception:** a HW-1 budget exception is permitted only through an explicit owner decision recorded after adoption, with scope, duration, safety evidence, and a reason other than speed alone. It must not waive TH-1 or HW-2.
- **No exception may fail HW-2:** an HW-1 exception may not create a condition that fails a HW-2 pass condition. For this purpose, HW-1.8's prohibition on normal-operation throttling is treated as part of HW-2.1. An exception that would require a permanent normal-operation throttle is not available.
- **No self-approval:** the author of a change cannot approve its own exception and cannot serve as its budget conformance reviewer.
- **No review-created authority:** no stage review, budget-conformance certification, or adoption record may create owner authority, amend a limit, or resolve a pending decision. Only the owner may do so through the appropriate recorded instrument.
- **No implied authorization:** an exception, stage acceptance, review, or this document does not authorize implementation, C8, or deployment.

## 13. Pending owner decisions

This document records engineering requirements only. It resolves no owner question.

- At drafting, all eight Stage 3 owner decisions are pending: D10(a), D10(b), D1a, D9, D1b, D2, D5, and D3.
- This is a status report, not an adoption condition. Adoption neither requires any decision to remain pending nor treats an answered decision as authority for this revision.
- If any Stage 3 decision is answered before adoption, r6 does not incorporate that answer as authority; the stage gates remain as stated in the owner-decision instrument.
- O-1 is a contract-deferred owner decision, separate from the eight Stage 3 decisions. The contract leaves “whether the Character picker may start a helper” to O-1. r6 fixes the adapter-side rule: `ProviderCatalog` may start, `GetTTSVoiceCatalog` must not. An O-1 answer authorizing picker startup is satisfied through `ProviderCatalog` and does not make `GetTTSVoiceCatalog` a start trigger.
- A conversational preference for a D10(a) result is not a recorded owner decision unless completed in the owner-decision instrument; this document neither adopts nor resolves it.
- A-1 is a separate open owner scope decision. It is not one of the eight Stage 3 decisions and is not resolved here.
- C1 and C2 remain unauthorized; all later stages remain unauthorized; C8 still requires separate runtime-execution authorization.
- No PR-1–PR-4 result has acceptance effect while A-1 is open.

## 14. Adoption checklist

Before owner adoption, all of the following must be true:

- [ ] Independent review of r6 is complete and recorded with a verified hash.
- [ ] The review found no unresolved authority, contract, lifecycle, or safety contradiction.
- [ ] The owner-adoption record names this exact artifact hash and scope.
- [ ] The §2.2 governing-input hashes were rechecked and match, including r4 and the r4 audit.
- [ ] The §7.1 matrix has six columns and separates evidence production, conformance assessment, stage acceptance, and requirement/architecture change; every row records the §7.2 separate-acceptor condition when the owner is the author; bare `AC-n` is used only for the original criteria and `R6-AC-n` only for the overlay.
- [ ] The §7.2 acceptance authority is recorded; the evidence producer did not accept the stage; if the owner is the author, the matrix and the acceptance record name the separate accountable acceptor as the actor, not the owner.
- [ ] The §7.3 original AC-1–AC-7 text is present and matches the original; no pointer or overlay narrows AC-2’s static review or AC-3’s repeatability.
- [ ] The C1 map includes `absent`, `crashed`/`backoff`, `stopping`, `incompatible`, and `unavailable`; the `crashed` → `backoff` → `starting` restart is not prohibited, uses the contract’s “next explicit request” / “next explicit TTS request” wording, and does not narrow the contract's request class; C5 owns its execution evidence.
- [ ] C1 owns the map and structural no-spawn checks; C5 owns executed restart/backoff; C7 owns executed start and reachability. C1’s handshake tests, including real-process fake-helper handshakes, remain in scope.
- [ ] C3 prohibits only unsolicited inspection/status starts and cross-refers the explicit status start to C7.
- [ ] The design is either incorporated as the stage plan or expressly left as a review input; it is not described as a requirement source or a controlling adopted instrument.
- [ ] Every budget row names an authorized source; a provisional value is treated as the ceiling unless an express range or owner decision permits more; the budget conformance reviewer is independent, issued a written determination, did not act as owner, and did not create a requirement by withholding.
- [ ] The §2.5 adoption-path check passes for A-1 Options A and B and for the open path; Option B has a plan-amendment instrument, exact scope, and SHA-256; the contract is not amended for C8; and no unqualified plan non-amendment statement is used.
- [ ] §13 records O-1 and preserves the `ProviderCatalog` / `GetTTSVoiceCatalog` distinction.
- [ ] C8 contains no performance claim or PR-1–PR-4 acceptance condition unless a recorded Option B plan amendment applies.
- [ ] The governing contract remains unamended; it has no C8 clause; the proposed design is unamended except a recorded §6.5 Option B plan amendment; a contract-level C8 requirement is recorded only under a separate §17 amendment, not as Option B.
- [ ] This revision resolves no owner decision; the adoption record does not require or imply that any decision remains pending, and does not cite an answered decision as authority.
- [ ] A-1 is explicitly decided before any PR-1–PR-4 result is given acceptance effect; if A-1 is open, the adoption record states that PR-1–PR-4 are inactive.
- [ ] No exception is being used to waive TH-1, HW-2, cleanup, lifecycle, or data-integrity requirements.
- [ ] The adoption record states that adoption does not authorize any stage or implementation.

Until every item is complete, this revision is **PROPOSED — NOT BINDING**. A checked box is a claim requiring evidence, not a waiver.
